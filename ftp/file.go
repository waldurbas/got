package ftp

// ----------------------------------------------------------------------------------
// file.go (https://github.com/waldurbas/got)
// Copyright 2020,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.10.03 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FileInfo #
type FileInfo struct {
	Name    string
	Size    int64
	Type    string
	ModTime time.Time
}

// STime #
func (f *FileInfo) STime() string {
	return f.ModTime.Format("2006.01.02  15:04:05")
}

// IsDir #
func (f *FileInfo) IsDir() bool {
	return f.Type == "D"
}

//modify=20190815120400;perm=adfrw;size=637616;type=file;unique=FE00UD6A921A8;UNIX.group=100;UNIX.mode=0644;UNIX.owner=1200; Differenz_20200_fileliste.txt
//modify=20190815120400;perm=adfrw;size=475466;type=file;unique=FE00UD6A921A9;UNIX.group=100;UNIX.mode=0644;UNIX.owner=1200; Differenz_80100_fileliste.txt
//modify=20190922224800;perm=flcdmpe;type=dir;unique=FE00U473F40FE;UNIX.group=0;UNIX.mode=0777;UNIX.owner=0; 4030457000007
func parseEntry(e string) *FileInfo {
	sf := strings.Split(e, "; ")
	if len(sf) != 2 {
		return nil
	}

	ss := strings.Split(sf[0], ";")
	m := make(map[string]string)

	for _, s := range ss {
		kv := strings.SplitN(s, "=", 2)
		if len(kv) == 2 {
			m[strings.ToLower(kv[0])] = strings.ToLower(kv[1])
		}
	}

	var (
		iSize int64
	)

	typ := m["type"]
	if typ == "file" {
		typ = "F"
		iSize, _ = strconv.ParseInt(m["size"], 10, 64)
		if iSize < 1 {
			return nil
		}
	} else if typ == "dir" {
		typ = "D"
	} else {
		return nil
	}

	mtime, err := time.ParseInLocation("20060102150405", m["modify"], time.UTC)
	if err != nil {
		return nil
	}

	f := &FileInfo{
		Name:    filepath.Base(sf[1]),
		Type:    typ,
		Size:    iSize,
		ModTime: mtime,
	}

	return f
}
