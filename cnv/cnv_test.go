package cnv_test

// ----------------------------------------------------------------------------------
// cnv_test.go (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2019.11.24 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"testing"

	"github.com/waldurbas/got/cnv"
)

func Test_checkTime(t *testing.T) {
	fmt.Println("Test_checkTime")
	var dtest = []struct {
		s string
		u int64
	}{
		{"2020-05-13T23:59:06", 1589414346},
		{"2000-01-01 00:00:00", 946684800},
		{"2000-01-01T00:00:00", 946684800},
		{"2020-05-27 07:24:06", 1590564246},
	}

	for _, tt := range dtest {
		u := cnv.TimeUTC2Unix(tt.s)

		if u != tt.u {
			t.Errorf("TimeUTC2Unix(%s): soll %d, ist %d", tt.s, tt.u, u)
		}

		s := cnv.Unix2UTCTimeStr(u)
		st := cnv.Unix2UTCTimeStrT(u)
		if s != tt.s && st != tt.s {
			t.Errorf("Unix2UTCTimeStr(%d): soll %s, ist %s", u, tt.s, s)
		}
	}
}

func Test_checkFTime(t *testing.T) {
	fmt.Println("Test_checkFTime")
	s := cnv.FTime()
	d := cnv.DatAsInt()
	ss := fmt.Sprintf("%08d", d)
	if ss != s[0:8] {
		t.Errorf("DatAsInt: soll %s, ist %s", s[0:8], ss)
	}
}

func Test_int2prs(t *testing.T) {
	fmt.Println("Test_int2prs")
	var dtest = []struct {
		s string
		u int
	}{
		{"0.01", 1},
		{"0.25", 25},
		{"0.99", 99},
		{"1.02", 102},
		{"99.29", 9929},
		{"199.00", 19900},
		{"1024.95", 102495},
	}

	for _, tt := range dtest {
		ss := cnv.Int2Prs(tt.u)

		if ss != tt.s {
			t.Errorf("Int2Prs(%d): soll %s, ist %s", tt.u, tt.s, ss)
		}
	}
}

func Test_int2dat(t *testing.T) {
	fmt.Println("Test_int2dat")
	var dtest = []struct {
		s string
		u int
	}{
		{"2020-05-13", 20200513},
		{"2000-01-01", 20000101},
	}

	for _, tt := range dtest {
		ss := cnv.Int2Dat(tt.u)

		if ss != tt.s {
			t.Errorf("Int2Dat(%d): soll %s, ist %s", tt.u, tt.s, ss)
		}
	}
}

func Test_int2dath(t *testing.T) {
	fmt.Println("Test_int2dath")
	var dtest = []struct {
		s string
		u int
	}{
		{"13-05-2020", 20200513},
		{"01-01-2000", 20000101},
	}

	for _, tt := range dtest {
		ss := cnv.Int2DatHuman(tt.u)

		if ss != tt.s {
			t.Errorf("Int2DatHuman(%d): soll %s, ist %s", tt.u, tt.s, ss)
		}
	}
}

func Test_str2dat(t *testing.T) {
	fmt.Println("Test_str2dat")

	var dtest = []struct {
		s string
		u int
	}{
		{"2020-05-13", 20200513},
		{"01.2.2010", 20100201},
		{"1.2.2010", 20100201},
		{"1.02.2010", 20100201},
	}

	for _, tt := range dtest {
		uu := cnv.Str2Dat(tt.s)

		if uu != tt.u {
			t.Errorf("Str2dat(%s): soll %d, ist %d", tt.s, tt.u, uu)
		}
	}
}

func Test_formatInt(t *testing.T) {
	fmt.Println("Test_formatInt")
	var dtest = []struct {
		s string
		u int
	}{
		{"20.200.513", 20200513},
		{"101", 101},
		{"0", 0},
		{"1.000", 1000},
		{"21.000", 21000},
		{"321.001", 321001},
		{"1.321.001", 1321001},
	}

	for _, tt := range dtest {
		ss := cnv.FormatInt(tt.u)

		if ss != tt.s {
			t.Errorf("FormatInt(%d): soll %s, ist %s", tt.u, tt.s, ss)
		}
	}
}

func Test_UUID(t *testing.T) {
	fmt.Println("Test_UUID")

	buid := cnv.UUID()
	suid := cnv.UUID36(string(buid))
	fmt.Println("Test_UUID:", suid)
}
