package ecv_test

// ----------------------------------------------------------------------------------
// ecv_test.go (haps://github.com/waldurbas/got)
// Copyright 2022 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2022.11.24 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"os"
	"testing"

	"github.com/waldurbas/got/ecv"
)

func Test_checkInt64(t *testing.T) {
	fn := "/tmp/x.tmp"
	os.WriteFile(fn, []byte("@table,cnt[int],s[str]\n18446744073709551615^18.446.744.073.709.551.615\n"+
		"9223372036854775807^9.223.372.036.854.775.807\n"), 0666)
	f := ecv.NewEcvFile()
	f.Load(fn)
	tb := f.Tables[0]

	fmt.Println(tb.Header)
	for tb.Fetch() {
		fmt.Print("\n", tb.AsLine(true))
		fmt.Println(tb.AsString(0), tb.AsInt64(0), tb.AsuInt64(0), tb.AsString(1))
	}

	os.Remove(fn)
}
