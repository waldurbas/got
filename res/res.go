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
// 2018.03.31 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"strings"
)

type resBox struct {
	data map[string]string
}

func newResBox() *resBox {
	return &resBox{data: make(map[string]string)}
}

// Find a File in the resBox
func (r *resBox) Has(file string) bool {
	if _, ok := r.data[file]; ok {
		return true
	}
	return false
}

// Get Data from resBox
func (r *resBox) Get(file string) (string, bool) {
	if f, ok := r.data[file]; ok {
		return f, ok
	}
	return "", false
}

// Add File into resBox
func (r *resBox) Add(file string, content string) {
	r.data[file] = content
}

// res.Data
var resData = newResBox()

// Get File from resData
func Get(file string) (string, bool) {

	cData, ok := resData.Get(file)
	if !ok || len(cData) < 10 {
		return "", false
	}

	cData = strings.Trim(cData, "\n")

	dData, _ := base64.StdEncoding.DecodeString(cData)
	rdata := bytes.NewReader(dData)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)

	return string(s), true
}

// Add File into resData
func Add(file string, content string) {
	resData.Add(file, content)
}

// Has a File in resData
func Has(file string) bool {
	return resData.Has(file)
}
