package param

// ----------------------------------------------------------------------------------
// param.go (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.02.10 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

type globalData struct {
	xargs        map[string]string
	debug        int
	xargsWithOut []string
}

var glo globalData

// init: wird automatisch aufgerufen
func init() {
	glo.xargs = make(map[string]string)

	var prev string
	for _, v := range os.Args[1:] {
		if v[0] == '-' || v[0] == '/' {
			prev = strings.ToLower(v[1:2])
			if prev == "q" || prev == "x" {
				glo.xargs[prev] = v[2:]
			} else {
				ix := strings.Index(v, "=")
				prev = ""
				if ix > 0 {
					prev = strings.ToLower(v[1:ix])
					glo.xargs[prev] = v[ix+1:]
				} else {
					prev = strings.ToLower(v[1:])
					glo.xargs[prev] = ""
				}
			}
		} else {
			glo.xargsWithOut = append(glo.xargsWithOut, v)
			if len(prev) > 0 {
				if len(glo.xargs[prev]) == 0 {
					glo.xargs[prev] = v
				}
			}
		}
	}

	ValueCheck("debug", "1")
	glo.debug = AsInt("debug", 0)
}

// Param #
func Param(ix int, def string) string {
	if ix >= len(glo.xargsWithOut) {
		return def
	}

	return glo.xargsWithOut[ix]
}

// Value #
func Value(sKey string, def string) string {
	lKey, ok := ValueExist(sKey)
	if !ok {
		return def
	}

	return glo.xargs[lKey]
}

// KeyExist #
func KeyExist(sKey string) bool {
	_, ok := Exist(sKey)
	return ok
}

// Exist #
func Exist(sKey string) (string, bool) {
	lKey := strings.ToLower(sKey)

	v, ok := glo.xargs[lKey]
	//	fmt.Println("ParamExist.Key: ", uKey, ", ok: ", ok, ", v: ", v)

	return v, ok
}

// Exists #
func Exists(sKeys []string) bool {
	for _, k := range sKeys {
		if KeyExist(k) {
			return true
		}
	}

	return false
}

// ValueExist #
func ValueExist(sKey string) (string, bool) {
	lKey := strings.ToLower(sKey)
	v, ok := glo.xargs[lKey]
	return lKey, ok && len(v) > 0
}

// AsInt #
func AsInt(sKey string, def int) int {
	v, ok := Exist(sKey)
	if !ok || len(v) == 0 {
		return def
	}

	return esubstr2int(v, 0, 10)
}

// Set #
func Set(sKey string, def string) {
	lKey := strings.ToLower(sKey)
	glo.xargs[lKey] = def
}

// ValueCheck #
func ValueCheck(sKey string, def string) {
	v, ok := Exist(sKey)

	if ok && len(v) == 0 {
		Set(sKey, def)
	}
}

// Print #
func Print() {
	fmt.Println("\n--> xParams:")
	for i, v := range glo.xargsWithOut {
		fmt.Printf("%d. [%s]\n", i, v)
	}

	fmt.Println("----------------------------")

	var sk []string
	for k := range glo.xargs {
		sk = append(sk, k)
	}
	sort.Strings(sk)

	for _, k := range sk {
		fmt.Printf("%-16.16s: [%s]\n", k, glo.xargs[k])
	}
	fmt.Print("\n\n")
}

// Debug #
func Debug() int {
	return glo.debug
}

func esubstr2int(s string, ix int, le int) int {
	b := []byte(s[ix:])
	l := len(s) - ix
	z := 0
	f := 1

	for i := 0; i < le && i < l; i++ {
		if b[i] >= '0' && b[i] <= '9' {
			z = z*10 + int(b[i]-'0')

		} else if b[i] == '-' {
			f = -1
		} else if b[i] == ';' {
			break
		}
	}

	return z * f
}

// GetEnv #
func GetEnv(key, defval string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defval
}