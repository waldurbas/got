package gcp

// ----------------------------------------------------------------------------------
// gcp.go (https://github.com/waldurbas/got) access to googlecloud
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.09.20 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/waldurbas/got/cnv"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// FileLocation #default europe-west3
var FileLocation = "europe-west3"

// FileEntry #
type FileEntry struct {
	FileName string `json:"filename"`
	Size     int64  `json:"size"`
	FTime    int64  `json:"ftime"`
}

// GCPbucket #
type GCPbucket struct {
	w          io.Writer
	projectID  string
	credFile   string
	bucketName string
	bhandle    *storage.BucketHandle
	client     *storage.Client
	ctx        context.Context
}

// GetCurrentProjectID #
func GetCurrentProjectID() (string, error) {
	ctx := context.Background()
	cred, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	s := ""
	if err == nil {
		s = cred.ProjectID
	}

	return s, err
}

// New #instance
func New(projectID string, bucketName string, credFile string) (*GCPbucket, error) {
	b := &GCPbucket{
		w:          os.Stderr,
		projectID:  projectID,
		credFile:   credFile,
		bucketName: bucketName,
		ctx:        context.Background(),
	}

	var err error

	if len(credFile) > 0 {
		b.client, err = storage.NewClient(b.ctx, option.WithCredentialsFile(credFile))
	} else {
		b.client, err = storage.NewClient(b.ctx)
	}
	if err != nil {
		return nil, err
	}

	b.bhandle = b.client.Bucket(b.bucketName)

	// check if bucket exists
	if _, err := b.bhandle.Attrs(b.ctx); err != nil {

		ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
		defer cancel()

		if err := b.bhandle.Create(ctx, b.projectID, &storage.BucketAttrs{Location: FileLocation}); err != nil {
			return nil, err
		}

		b.bhandle = b.client.Bucket(b.bucketName)
	}

	return b, nil
}

// Close #
func (b *GCPbucket) Close() {
	b.client.Close()
}

// BucketExists #
func (b *GCPbucket) BucketExists(bucketName string) bool {
	bh := b.client.Bucket(bucketName)
	if _, err := bh.Attrs(b.ctx); err != nil {
		return false
	}

	return true
}

// BucketCreate #
func (b *GCPbucket) BucketCreate(bucketName string) error {
	bucket := b.client.Bucket(bucketName)

	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	if err := bucket.Create(ctx, b.projectID, &storage.BucketAttrs{
		//		StorageClass: "COLDLINE",
		Location: FileLocation,
	}); err != nil {
		return err
	}

	b.bucketName = bucketName
	b.bhandle = b.client.Bucket(b.bucketName)
	return nil
}

// BucketRemove #
func (b *GCPbucket) BucketRemove(bucketName string) error {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()
	if err := b.client.Bucket(bucketName).Delete(ctx); err != nil {
		return err
	}

	return nil
}

// ListRoot #
func (b *GCPbucket) ListRoot() *[]FileEntry {
	files := []FileEntry{}

	it := b.client.Buckets(b.ctx, b.projectID)
	for {
		fa, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil
		}

		f := FileEntry{FileName: fa.Name, FTime: fa.Created.Unix()}

		files = append(files, f)
	}

	return &files
}

// ListFiles # listet alle files, quasi recursiv
func (b *GCPbucket) ListFiles(dirName string, startIx int) (*[]FileEntry, error) {
	return b.listFiles(dirName, startIx, "_")
}

// ListDirFiles # listen NUR files in Dir
func (b *GCPbucket) ListDirFiles(dirName string, startIx int) (*[]FileEntry, error) {
	le := len(dirName)

	// Slash am ende muss nicht unbedingt uebergeben werden
	if le > 1 && dirName[le-1:le] != "/" {
		dirName = dirName + "/"
		if startIx > -1 {
			startIx++
		}
	}

	return b.listFiles(dirName, startIx, "/")
}

//--------------------------------------------------------------------------
// Prefixes and delimiters can be used to emulate directory listings.
// Prefixes can be used to filter objects starting with prefix.
// The delimiter argument can be used to restrict the results to only the
// objects in the given "directory". Without the delimiter, the entire tree
// under the prefix is returned.
//
// For example, given these blobs:
//   /a/1.txt
//   /a/b/2.txt
//
// If you just specify prefix="a/", you'll get back:
//   /a/1.txt
//   /a/b/2.txt
//
// However, if you specify prefix="a/" and delim="/", you'll get back:
//   /a/1.txt
//--------------------------------------------------------------------------
func (b *GCPbucket) listFiles(dirName string, startIx int, delim string) (*[]FileEntry, error) {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	if len(dirName) < 1 {
		return nil, errors.New("dirName length to short..")
	}

	var qry *storage.Query

	if dirName == "/" {
		qry = &storage.Query{Delimiter: delim}
	} else {
		qry = &storage.Query{Prefix: dirName, Delimiter: delim}
	}

	it := b.bhandle.Objects(ctx, qry)
	files := []FileEntry{}
	for {
		fa, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		if fa.Created.Unix() < 10 {
			continue
		}

		f := &FileEntry{FileName: fa.Name, Size: fa.Size, FTime: fa.CustomTime.Unix()}
		if startIx > -1 {
			f.FileName = strings.ReplaceAll(f.FileName[startIx:], "/", "")
		}

		files = append(files, *f)
	}

	return &files, nil
}

// FileWrite #
func (b *GCPbucket) FileWrite(filePath string, data *[]byte) error {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	wr := b.bhandle.Object(filePath).NewWriter(ctx)
	_, err := wr.Write(*data)
	if err != nil {
		return err
	}

	_ = wr.Close()

	return nil
}

// FileWriteWithTimeStamp #
func (b *GCPbucket) FileWriteWithTimeStamp(filePath string, data *[]byte, ti time.Time) error {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	wr := b.bhandle.Object(filePath).NewWriter(ctx)
	_, err := wr.Write(*data)
	if err != nil {
		return err
	}

	_ = wr.Close()

	o := b.bhandle.Object(filePath)
	u := storage.ObjectAttrsToUpdate{CustomTime: ti}

	_, err = o.Update(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

// FileRead #
func (b *GCPbucket) FileRead(filePath string) (*[]byte, time.Time, string, error) {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	var rt time.Time
	ob := b.bhandle.Object(filePath)

	rc, err := ob.NewReader(ctx)
	if err != nil {
		return nil, rt, "", err
	}
	defer rc.Close()

	at, _ := ob.Attrs(ctx)
	rt = at.CustomTime

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, rt, "", err
	}

	return &data, rt, at.MediaLink, nil
}

// FileStat #
func (b *GCPbucket) FileStat(filePath string) (*FileEntry, error) {
	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	fa, err := b.bhandle.Object(filePath).Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return &FileEntry{FileName: fa.Name, Size: fa.Size, FTime: fa.CustomTime.Unix()}, nil
}

// FileStatDump #
func (b *GCPbucket) FileStatDump(filePath string) {
	b.printf("\nFile stat:\n")

	ctx, cancel := context.WithTimeout(b.ctx, time.Second*10)
	defer cancel()

	obj, err := b.bhandle.Object(filePath).Attrs(ctx)
	if err != nil {
		b.printf("\nFile stat:\n")
		b.printf("statFile: unable to stat file from bucket %q, file %q: %v", b.bucketName, filePath, err)
		return
	}

	b.dumpStats(obj)
}

// BaseName #
func (f *FileEntry) BaseName() string {
	idx := strings.LastIndex(f.FileName, "/")
	if idx < 1 {
		return f.FileName
	}

	return f.FileName[idx+1:]
}

// DirName #
func (f *FileEntry) DirName() string {
	idx := strings.LastIndex(f.FileName, "/")
	if idx < 1 {
		return f.FileName
	}

	return f.FileName[:idx]
}

// IsDir #
func (f *FileEntry) IsDir() bool {
	return f.Size == 0
}

// Print #
func (f *FileEntry) Print(fullname bool) {
	fmt.Print(time.Unix(f.FTime, 0).Format("2006-01-02 15:04"))

	if fullname {
		if f.IsDir() {
			fmt.Printf("%10s    %s\n", "<DIR>", f.FileName)
		} else {
			fmt.Printf("%12s  %s\n", cnv.FormatInt64(f.Size), f.FileName)
		}

	} else {
		if f.IsDir() {
			fmt.Printf("%10s    %s\n", "<DIR>", f.DirName())
		} else {
			fmt.Printf("%12s  %s\n", cnv.FormatInt64(f.Size), f.BaseName())
		}
	}
}

func (b *GCPbucket) dumpStats(obj *storage.ObjectAttrs) {
	b.printf("Filename: /%v/%v\n", obj.Bucket, obj.Name)
	b.printf("ContentType: %q\n", obj.ContentType)
	b.printf("ACL: %#v\n", obj.ACL)
	b.printf("Owner: %v\n", obj.Owner)
	b.printf("ContentEncoding: %q\n", obj.ContentEncoding)
	b.printf("Size: %v\n", obj.Size)
	b.printf("MD5: %q\n", obj.MD5)
	b.printf("CRC32C: %v\n", obj.CRC32C)
	b.printf("Metadata: %#v\n", obj.Metadata)
	b.printf("MediaLink: %q\n", obj.MediaLink)
	b.printf("StorageClass: %q\n", obj.StorageClass)
	if !obj.Deleted.IsZero() {
		b.printf("Deleted: %v\n", obj.Deleted)
	}
	b.printf("Updated: %v\n", obj.Updated)
	b.printf("CustomTime: %v\n", obj.CustomTime)
}

func (b *GCPbucket) printf(format string, v ...interface{}) {
	fmt.Fprintf(b.w, format, v...)
}
