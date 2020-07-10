package htp

// ----------------------------------------------------------------------------------
// htp.go (https://github.com/waldurbas/got)
// Copyright 2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.07.05 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// HtReq #
type HtReq struct {
	StatusCode int
	Status     string
	Msg        string
	Body       []byte
	Err        error
	Map        map[string]interface{}
}

// Post #
func Post(url string, token string, jsData *[]byte) *HtReq {
	return Request(url, "POST", token, jsData)
}

// Get #
func Get(url string, token string, jsData *[]byte) *HtReq {
	return Request(url, "GET", token, jsData)
}

// Put #
func Put(url string, token string, jsData *[]byte) *HtReq {
	return Request(url, "PUT", token, jsData)
}

// Delete #
func Delete(url string, token string, jsData *[]byte) *HtReq {
	return Request(url, "DELETE", token, jsData)
}

// Request #
func Request(url string, method string, token string, jsData *[]byte) *HtReq {
	rr := HtReq{}
	rr.Map = make(map[string]interface{})

	req, err := http.NewRequest(method, url, bytes.NewBuffer(*jsData))
	if err != nil {
		rr.Err = err
		rr.Msg = "error.NewRequest"
		return &rr
	}

	// set headers
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	if len(token) > 0 {
		req.Header.Add("Authorization", token)
	}

	// set client timeout
	client := &http.Client{Timeout: time.Second * 15}

	// send request
	rsp, e := client.Do(req)
	if e != nil {
		rr.Err = e
		rr.Msg = "error.ClientDo"
		if rsp != nil {
			rr.Status = rsp.Status
			rr.StatusCode = rsp.StatusCode
		}
		return &rr
	}
	defer rsp.Body.Close()

	rr.Status = rsp.Status
	rr.StatusCode = rsp.StatusCode

	rr.Body, rr.Err = ioutil.ReadAll(rsp.Body)

	if rr.Err != nil {
		rr.Msg = "error.ReadAll"
	} else {
		json.Unmarshal(rr.Body, &rr.Map)
	}

	return &rr
}

// JSON2Str #
func JSON2Str(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)

	empty := ""
	encoder.SetIndent(empty, "  ")

	err := encoder.Encode(data)
	if err != nil {
		return empty
	}
	return buffer.String()
}
