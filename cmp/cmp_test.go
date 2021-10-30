package cmp_test

// ----------------------------------------------------------------------------------
// cmp_test.go (https://github.com/waldurbas/got
// Copyright 2020,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.02.11 (wu) add Test_StrXcmp
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"testing"

	"github.com/waldurbas/cmp"
)

func Test_StrXcmp(t *testing.T) {
	fmt.Println("Test_StrXcmp")
	var items = []struct {
		a string
		b string
		r int
	}{
		{"abc", "AbC", 0},
		{"abc", "bAbC", -1},
		{"acd", "Aca", 3},
		{"abc", "ABC", 0},
		{"abc", "abc", 0},
		{"abc", "abcxdf", -1},
		{"abcxdf", "abc", 1},
		{"xbcxdf", "abc", 23},
		{`xBc?1#- &üäöß%$§"'`, `XBc?1#- &ÜÄöß%$§"'`, 0},
	}

	for _, item := range items {

		r := cmp.StrIxcmp(item.a, item.b)

		if r != item.r {
			t.Errorf("StrXcmp: a:(%s), b:(%s): ist %d, soll %d", item.a, item.b, r, item.r)
		}
	}
}
