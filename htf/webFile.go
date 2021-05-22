package htf

// ----------------------------------------------------------------------------------
// webfile.go (https://github.com/waldurbas/got/htf)
// Copyright 2018,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.05.22 (wu) in htf überführt
// 2018.12.11 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

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
	URL      string
}

// DownloadFilesInfo #
type DownloadFilesInfo struct {
	URL  string
	List []DownloadFileInfo
}

// GetDownloadFilesInfo #
func GetDownloadFilesInfo(url string) (*DownloadFilesInfo, error) {

	var downFiles DownloadFilesInfo

	downFiles.URL = url
	buf, err := urlDownloadListFile(downFiles.URL + "/list")
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

			file := DownloadFileInfo{FileName: items[0], Web: wInfo, Loc: lInfo, URL: url}
			downFiles.List = append(downFiles.List, file)
		}
	}

	return &downFiles, nil
}

func (f *DownloadFileInfo) urlFile(toFile string) string {
	sFile := f.URL + "/" + f.FileName
	if strings.HasSuffix(toFile, ".gz") && !strings.HasSuffix(f.FileName, ".gz") {
		sFile = sFile + ".gz"
	}

	return sFile
}

// Download #
func (f *DownloadFileInfo) Download(toFile string) error {
	if err := DownloadFile(f.urlFile(toFile), toFile); err != nil {
		return err
	}

	return f.SetFileTime(toFile)
}

// ParalellDownload #
func (f *DownloadFileInfo) ParalellDownload(toFile string) error {
	if err := ParalellDownloadFile(f.urlFile(toFile), toFile); err != nil {
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
	b := strings.LastIndex(FileName, string(os.PathSeparator))
	loFile := strings.ToLower(FileName[b+1:])

	for _, f := range fl.List {
		wFile := strings.ToLower(f.FileName)

		if wFile == loFile {
			loc, _ := time.LoadLocation("UTC")

			st, err := os.Stat(FileName)
			if err != nil {
				f.Loc.Time = f.Web.Time
				f.Loc.Size = 0
			} else {
				f.Loc.Time = st.ModTime().In(loc)
				f.Loc.Size = uint64(st.Size())
			}

			f.Changed = (f.Web.Size != f.Loc.Size) || (f.Loc.Time != f.Web.Time)

			fmt.Printf("\nwebFile: %d %v\n", f.Web.Size, f.Web.Time)
			fmt.Printf("locFile: %d %v, changed=%v\n", f.Loc.Size, f.Loc.Time, f.Changed)

			return &f, nil
		}
	}

	f := &DownloadFileInfo{FileName: FileName[b+1:]}
	return f, errors.New("file not found")
}
