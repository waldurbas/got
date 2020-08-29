package xtl

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------------
// xFile.go for Go's xtl package (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.08.29 (wu) Add MoveFile, LoadFiles
// 2018.12.11 (wu) Init
//-----------------------------------------------------------------------------------

const bufferSize = 32768

// CopyFile with bufferig #
func CopyFile(src, dst string) (int, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()

	nn := 0
	buf := make([]byte, bufferSize)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return 0, err
		}
		nn += n
	}

	return nn, err
}

// MoveFile #
func MoveFile(src string, dst string) error {
	if FileExists(dst) {
		os.Remove(dst)
	}
	_, err := CopyFile(src, dst)
	if err == nil {
		err = os.Remove(src)
	}

	return err
}

// FileExists #
func FileExists(filename string) bool {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !f.IsDir()
}

// FileInfo #
type FileInfo struct {
	FileName string
	Size     int64
	Time     time.Time
}

// LoadFiles #
func LoadFiles(path, match string) (*[]FileInfo, error) {
	d, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()

	dfiles, err := d.Readdir(-1)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("UTC")

	files := []FileInfo{}
	for _, fInfo := range dfiles {
		if fInfo.Mode().IsRegular() {
			matched, err := filepath.Match(match, fInfo.Name())
			if err == nil && matched {
				f := FileInfo{FileName: fInfo.Name(), Size: fInfo.Size(), Time: fInfo.ModTime().In(loc)}
				files = append(files, f)
			}
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return strings.ToLower(files[i].FileName) < strings.ToLower(files[j].FileName)
	})

	return &files, nil
}
