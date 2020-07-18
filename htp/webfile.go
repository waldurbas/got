package htp

// ----------------------------------------------------------------------------------
// webfile.go (https://github.com/waldurbas/got)
// Copyright 2018,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.07.18 (wu) taken over and adapted from github.com/waldurbas/xt
// 2018.12.11 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// WriteCounter #
type WriteCounter struct {
	Total uint64
}

// XFileInfo #
type XFileInfo struct {
	Size uint64
	Time time.Time
}

// DownloadFileInfo #
type DownloadFileInfo struct {
	FileName string
	Web      XFileInfo
	Loc      XFileInfo
	Changed  bool
	parent   *DownloadFilesInfo
}

// DownloadFilesInfo #
type DownloadFilesInfo struct {
	url  string
	List []DownloadFileInfo
}

// GetDownloadFilesInfo #
func GetDownloadFilesInfo(url string) (*DownloadFilesInfo, error) {

	var downFiles DownloadFilesInfo

	downFiles.url = url
	buf, err := urlDownloadListFile(downFiles.url + "/download.txt")
	if err != nil {
		return &downFiles, err
	}

	fList := strings.Split(buf, "\n")
	for _, line := range fList {
		items := strings.Split(line, ";")
		if len(items) > 2 {
			size, _ := strconv.Atoi(items[1])
			t, _ := time.Parse("2006-01-02 15:04:05", items[2])

			wInfo := XFileInfo{Size: uint64(size), Time: t}
			lInfo := XFileInfo{}

			file := DownloadFileInfo{FileName: items[0], Web: wInfo, Loc: lInfo, parent: &downFiles}
			downFiles.List = append(downFiles.List, file)
		}
	}

	return &downFiles, nil
}

// DownloadFile #
func DownloadFile(url string, toFile string) error {
	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	tmpFile := toFile + ".tmp"
	os.Remove(tmpFile)

	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()
	out.Close()

	// Rename the tmp file back to the original file
	time.Sleep(1 * time.Second)
	err = os.Rename(tmpFile, toFile)
	if err != nil {
		return err
	}

	return nil
}

// URLfileSize #
func URLfileSize(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}

	// Is our request ok?
	if resp.StatusCode != http.StatusOK {
		err := errors.New(resp.Status)
		return 0, err
	}

	// the Header "Content-Length" will let us know
	// the total file size to download
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	return size, nil
}

// Write #
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

func readableBytes(n uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	b := float64(1024)
	e := math.Floor(math.Log(float64(n)) / math.Log(b))
	sfx := sizes[int(e)]
	v := float64(n) / math.Pow(b, math.Floor(e))
	f := "%.0f"
	if v < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", v, sfx)
}

// PrintProgress #prints the progress of a file write
func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s\rDownloading... %s complete", strings.Repeat(" ", 50), readableBytes(wc.Total))
}

func urlDownloadListFile(url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Download #
func (f *DownloadFileInfo) Download(toFile string) error {

	urlFile := f.parent.url + "/" + f.FileName + ".gz"
	if err := DownloadFile(urlFile, toFile); err != nil {
		return err
	}

	return f.SetFileTime(toFile)
}

// SetFileTime #
func (f *DownloadFileInfo) SetFileTime(toFile string) error {
	// setFileTime: change both atime and mtime to currenttime
	return os.Chtimes(toFile, f.Web.Time, f.Web.Time)
}

// GetFileInfo #
func (fl *DownloadFilesInfo) GetFileInfo(FileName string) (*DownloadFileInfo, error) {
	lowerFile := strings.ToLower(FileName)

	for _, f := range fl.List {
		wFile := strings.ToLower(f.FileName)

		if wFile == lowerFile {
			loc, _ := time.LoadLocation("UTC")

			st, err := os.Stat(f.FileName)
			if err != nil {
				f.Loc.Time = f.Web.Time
				f.Loc.Size = 0
			} else {
				f.Loc.Time = st.ModTime().In(loc)
				f.Loc.Size = uint64(st.Size())
			}

			//			dif := f.Loc.Time.Sub(f.Web.Time)
			f.Changed = (f.Web.Size != f.Loc.Size) || (f.Loc.Time != f.Web.Time)

			//	fmt.Printf("webFile: %d %v\n", f.Web.Size, f.Web.Time)
			//	fmt.Printf("locFile: %d %v\n", f.Loc.Size, f.Loc.Time)

			return &f, nil
		}
	}

	return nil, nil
}
