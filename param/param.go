package param

// ----------------------------------------------------------------------------------
// param.go (https://github.com/waldurbas/got)
// Copyright 2019,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.01.28 (wu) Printparams
// 2020.08.29 (wu) CheckVersion
// 2020.02.10 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"sort"
	"strconv"
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
		//		fmt.Printf("\narg: [%v]", v)
		if v[0] == '-' || v[0] == '/' {
			prev = strings.ToLower(v[1:2])

			if prev == "q" || ((prev == "x" || prev == "u") && len(v) > 2 && strings.Index(v, "=") < 0) {
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
	if glo.debug > 0 {
		os.Setenv("DEBUG", strconv.Itoa(glo.debug))
	}
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
	val, ok := ValueExist(sKey)
	if !ok {
		return def
	}

	return val
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
	return v, ok && len(v) > 0
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

// CheckVersion #
func CheckVersion() bool {
	return !KeyExist("noCheckVersion") &&
		!KeyExist("noVersion") &&
		!KeyExist("noVers") &&
		!KeyExist("noUpd") &&
		!KeyExist("noUpdate")

}

// PrintParams #
func PrintParams() string {

	s := fmt.Sprintln("\n--> xParams:")
	for i, v := range glo.xargsWithOut {
		s += fmt.Sprintf("%d. [%s]\n", i, v)
	}

	s += fmt.Sprintln("----------------------------")

	var sk []string
	for k := range glo.xargs {
		sk = append(sk, k)
	}
	sort.Strings(sk)

	for _, k := range sk {
		s += fmt.Sprintf("%-16.16s: [%s]\n", k, glo.xargs[k])
	}
	s += fmt.Sprint("\n\n")

	return s
}

// Print #
func Print() {
	fmt.Fprint(os.Stderr, PrintParams())
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

// Xargs #
func Xargs() map[string]string {
	return glo.xargs
}

// Count #
func Count() int {
	return len(glo.xargs)
}
