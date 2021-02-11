package xtl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
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
// 2021.02.11 (wu) LoadStructFromFile,SaveStructToFile
// 2021.01.07 (wu) func LoadFile(sfile string) (*FileData, error)
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

// FileData #
type FileData struct {
	FileName string
	Size     int64
	Time     int64
	Data     []byte
}

const timewebLayout = "2006-01-02T15:04:05"

var locUTC *time.Location

// UTCTime #
func (fd *FileData) UTCTime() string {
	if locUTC == nil {
		locUTC, _ = time.LoadLocation("UTC")
	}
	return time.Unix(fd.Time, 0).In(locUTC).Format(timewebLayout)
}

// LoadFile #
func LoadFile(sfile string) (*FileData, error) {

	stat, err := os.Stat(sfile)
	if os.IsNotExist(err) {
		return nil, err
	}

	if stat.IsDir() {
		return nil, errors.New("is not a file")
	}

	file, err := os.Open(sfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cf := &FileData{
		FileName: sfile,
		Size:     stat.Size(),
		Time:     stat.ModTime().Unix(),
	}

	cf.Data = make([]byte, cf.Size)

	buffer := bufio.NewReader(file)
	_, err = buffer.Read(cf.Data)

	if err != nil {
		return nil, err
	}

	return cf, nil
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

// ChangeFileExt #
func ChangeFileExt(sfile string, newext string) string {
	ext := path.Ext(sfile)
	return sfile[0:len(sfile)-len(ext)] + newext
}

// FuncUnMarshal #
type FuncUnMarshal func(r io.Reader, v interface{}) error

// FuncMarshal #
type FuncMarshal func(v interface{}) (io.Reader, error)

var lockF sync.Mutex

// LoadStructFromFile #
func LoadStructFromFile(sFile string, v interface{}, fu FuncUnMarshal) error {
	lockF.Lock()
	defer lockF.Unlock()

	f, err := os.Open(sFile)
	if err != nil {
		return err
	}

	defer f.Close()

	if fu != nil {
		return fu(f, v)
	}

	return json.NewDecoder(f).Decode(v)
}

// SaveStructToFile #
func SaveStructToFile(sFile string, v interface{}, fu FuncMarshal) error {
	lockF.Lock()
	defer lockF.Unlock()

	f, err := os.Create(sFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if fu != nil {
		ir, err := fu(f)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, ir)
		return err
	}

	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	_, err = io.Copy(f, bytes.NewReader(b))
	return err
}
