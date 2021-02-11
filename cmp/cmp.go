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
// 2021.02.11 (wu) add StrXcmp
// 2020.07.02 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"strings"
)

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

func c2lower(c rune) rune {
	if 'A' <= c && c <= 'Z' {
		return c + 'a' - 'A'
	}

	switch c {
	case 'Ä':
		return 'ä'
	case 'Ö':
		return 'ö'
	case 'Ü':
		return 'ü'
	}

	return c
}

// StrIxcmp #
// compare two strings case insensitiv
func StrIxcmp(a, b string) int {
	ra := []rune(a)
	rb := []rune(b)
	ale := len(ra)
	ble := len(rb)

	i := 0
	for ; i < ale && i < ble; i++ {
		ac := c2lower(ra[i])
		bc := c2lower(rb[i])

		if ac != bc {
			return int(ac - bc)
		}
	}

	if i < ale {
		return 1
	}

	if i < ble {
		return -1
	}

	return 0
}
