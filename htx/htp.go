package htx

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
// 2020.09.09 (wu) Patch function,Request. 1.st Param is the Method
// 2020.09.06 (wu) WriteResponse
// 2020.07.25 (wu) HtReqMap
// 2020.07.05 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/waldurbas/got/lgx"
)

// HtReq #
type HtReq struct {
	StatusCode int
	Status     string
	Msg        string
	Body       []byte
	Err        error
}

// HtReqMap #
type HtReqMap struct {
	m map[string]interface{}
}

// ReqParams #
type ReqParams struct {
	qry map[string]string
}

// ReqQuery #
func ReqQuery(r *http.Request) *ReqParams {
	v := &ReqParams{}
	v.qry = make(map[string]string)

	qry, err := url.ParseQuery(r.URL.RawQuery)
	if err == nil {
		if qry != nil {
			for kk, vv := range qry {
				v.qry[kk] = vv[0]
			}
		}
	}

	return v
}

// Exists #
func (m *ReqParams) Exists(key string) bool {
	_, ok := m.qry[key]
	return ok
}

// Value #
func (m *ReqParams) Value(key string) string {
	v, _ := m.qry[key]
	return v
}

// Items #
func (m *ReqParams) Items() map[string]string {
	return m.qry
}

// WriteResponseMsg #
func WriteResponseMsg(from string, w http.ResponseWriter, statusCode int, statusStr string) {
	if statusCode == 200 {
		lgx.PrintDebug(from, "OK")
	} else {
		lgx.PrintError(from, statusStr)
	}

	WriteResponse(w, statusCode, map[string]interface{}{"message": statusStr})
}

// WriteResponse #
func WriteResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")

	b, e := json.Marshal(body)
	if e != nil {
		lgx.PrintError("serializing", e)
		w.WriteHeader(404)
		w.Write([]byte(e.Error()))
		return
	}

	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(b))
	if err != nil {
		lgx.PrintError("writeReponse", err)
	}
}

// Post #
func Post(url string, token string, jsData *[]byte) *HtReq {
	return Request("POST", url, token, jsData)
}

// Get #
func Get(url string, token string, jsData *[]byte) *HtReq {
	return Request("GET", url, token, jsData)
}

// Put #
func Put(url string, token string, jsData *[]byte) *HtReq {
	return Request("PUT", url, token, jsData)
}

// Patch #
func Patch(url string, token string, jsData *[]byte) *HtReq {
	return Request("PATCH", url, token, jsData)
}

// Delete #
func Delete(url string, token string, jsData *[]byte) *HtReq {
	return Request("DELETE", url, token, jsData)
}

// Request #
func Request(method string, url string, token string, jsData *[]byte) *HtReq {
	rr := HtReq{}

	var r io.Reader
	if jsData != nil {
		r = bytes.NewBuffer(*jsData)
	} else {
		r = nil
	}

	req, err := http.NewRequest(method, url, r)
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

	/*	netTrans := &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		}
	*/

	// set client timeout and Transport
	client := &http.Client{
		Timeout: time.Second * 15,
		//		Transport: netTrans,
	}

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

// ToMap #
func (r *HtReq) ToMap() *HtReqMap {
	x := &HtReqMap{}
	json.Unmarshal(r.Body, &x.m)
	return x
}

// AsString #
func (x *HtReqMap) AsString(key string) string {
	switch v := x.m[key].(type) {
	case string:
		return v
	}

	return ""
}

// StringMap #
func (x *HtReqMap) StringMap() map[string]string {
	m := map[string]string{}

	for k, i := range x.m {
		switch v := i.(type) {
		case string:
			m[k] = v
		}
	}

	return m
}
