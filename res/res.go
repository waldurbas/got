package res

// ----------------------------------------------------------------------------------
// res.go (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.07.24 (wu) New() nach aussen
// 2018.03.31 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

// Resdata #
type Resdata struct {
	data map[string]string
}

// New #
func New() *Resdata {
	return &Resdata{data: make(map[string]string)}
}

// Has #Find a File in the resData
func (r *Resdata) Has(file string) bool {
	if _, ok := r.data[file]; ok {
		return true
	}
	return false
}

// Get Data from resData
func (r *Resdata) Get(file string) (string, bool) {
	if f, ok := r.data[file]; ok {
		return f, ok
	}
	return "", false
}

// Add File into resData
func (r *Resdata) Add(file string, content string) {
	r.data[file] = content
}

// GetDecoded File from resData
func (r *Resdata) GetDecoded(file string) (string, bool) {

	cData, ok := r.Get(file)
	if !ok || len(cData) < 10 {
		return "", false
	}

	cData = strings.Trim(cData, "\n")

	dData, _ := base64.StdEncoding.DecodeString(cData)
	rdata := bytes.NewReader(dData)
	rd, _ := gzip.NewReader(rdata)
	ss, _ := ioutil.ReadAll(rd)

	return string(ss), true
}
