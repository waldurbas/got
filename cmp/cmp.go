package cmp

// ----------------------------------------------------------------------------------
// cmp.go (https://github.com/waldurbas/got)
// Copyright 2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.07.02 (wu) Init
//-----------------------------------------------------------------------------------

import "strings"

// StrIndexOfI #Ignore
func StrIndexOfI(str string, needle string) int {
	if needle == "" || str == "" {
		return -1
	}

	return strings.Index(strings.ToLower(str), strings.ToLower(needle))
}

// StrIndexOf #
func StrIndexOf(str string, needle string) int {
	if needle == "" || str == "" {
		return -1
	}

	return strings.Index(str, needle)
}

// StrStr #wenn TeilString gefunden, den Rest liefern
func StrStr(fStr string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(fStr, needle)
	if idx == -1 {
		return ""
	}
	return fStr[idx:]
}

// StrnIcmp #
func StrnIcmp(a, b string, le int) bool {
	l1 := len(a)
	l2 := len(b)

	if l1 < le || l2 < le {
		return false
	}

	return strings.EqualFold(a[0:le], b[0:le])
}

// StrIcmp #String Ignore Compare
func StrIcmp(a, b string) bool {
	return strings.EqualFold(a, b)
}

// StrComp #string-compare
func StrComp(a, b string) int {
	if a == b {
		return 0
	}

	if a < b {
		return -1
	}

	return 1
}
