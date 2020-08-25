package xtl

import (
	"fmt"
	"io"
	"os"
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
