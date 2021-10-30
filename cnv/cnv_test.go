package cnv_test

// ----------------------------------------------------------------------------------
// cnv_test.go (haps://github.com/waldurbas/got)
// Copyright 2019,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.02.21 (wu) IsDigit,IsAlpha
// 2019.11.24 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"testing"
	"time"

	"github.com/waldurbas/cnv"
)

func Test_checkTimeDif(t *testing.T) {
	var ar = []struct {
		t  string
		xs int
		hh int
		mm int
		ss int
	}{
		{"23:59:06", 86346, 23, 59, 6},
		{"00:22:01", 1321, 0, 22, 1},
		{"02:01:01", 7261, 2, 1, 1},
		{"00:01:06", 66, 0, 1, 6},
	}

	z := time.Unix(0, 0).UTC()
	for _, a := range ar {
		r, _ := time.Parse(cnv.DT_TIM, a.t)
		xs, hh, mm, ss := cnv.TimeDif(z, r)
		if a.xs != xs || a.hh != hh || a.mm != mm || a.ss != ss {
			t.Errorf("TimeDif(%s): soll %d, ist %d", a.t, a.xs, xs)
		}

	}
}

func Test_checkTime(t *testing.T) {
	var ar = []struct {
		s string
		u int64
	}{
		{"2020-05-13T23:59:06", 1589414346},
		{"2000-01-01 00:00:00", 946684800},
		{"2000-01-01T00:00:00", 946684800},
		{"2020-05-27 07:24:06", 1590564246},
	}

	for _, a := range ar {
		u := cnv.TimeUTC2Unix(a.s)

		if u != a.u {
			t.Errorf("TimeUTC2Unix(%s): soll %d, ist %d", a.s, a.u, u)
		}

		s := cnv.Unix2UTCTimeStr(u)
		st := cnv.Unix2UTCTimeStrT(u)
		if s != a.s && st != a.s {
			t.Errorf("Unix2UTCTimeStr(%d): soll %s, ist %s", u, a.s, s)
		}
	}
}

func Test_checkFTime(t *testing.T) {
	s := cnv.FTime()
	d := cnv.DatAsInt()
	ss := fmt.Sprintf("%08d", d)
	if ss != s[0:8] {
		t.Errorf("DatAsInt: soll %s, ist %s", s[0:8], ss)
	}
}

func Test_int2prs(t *testing.T) {
	var ar = []struct {
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

	for _, a := range ar {
		ss := cnv.Int2Prs(a.u)

		if ss != a.s {
			t.Errorf("Int2Prs(%d): soll %s, ist %s", a.u, a.s, ss)
		}
	}
}

func Test_int2dat(t *testing.T) {
	var ar = []struct {
		s string
		u int
	}{
		{"2020-05-13", 20200513},
		{"2000-01-01", 20000101},
	}

	for _, a := range ar {
		ss := cnv.Int2Dat(a.u)

		if ss != a.s {
			t.Errorf("Int2Dat(%d): soll %s, ist %s", a.u, a.s, ss)
		}
	}
}

func Test_int2dath(t *testing.T) {
	var ar = []struct {
		s string
		u int
	}{
		{"13-05-2020", 20200513},
		{"01-01-2000", 20000101},
	}

	for _, a := range ar {
		ss := cnv.Int2DatHuman(a.u)

		if ss != a.s {
			t.Errorf("Int2DatHuman(%d): soll %s, ist %s", a.u, a.s, ss)
		}
	}
}

func Test_str2dat(t *testing.T) {
	var ar = []struct {
		s string
		u int
	}{
		{"2020-05-13", 20200513},
		{"01.2.2010", 20100201},
		{"1.2.2010", 20100201},
		{"1.02.2010", 20100201},
	}

	for _, a := range ar {
		uu := cnv.Str2Dat(a.s)

		if uu != a.u {
			t.Errorf("Str2dat(%s): soll %d, ist %d", a.s, a.u, uu)
		}
	}
}

func Test_formatInt(t *testing.T) {
	var ar = []struct {
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

	for _, a := range ar {
		ss := cnv.FormatInt(a.u)

		if ss != a.s {
			t.Errorf("FormatInt(%d): soll %s, ist %s", a.u, a.s, ss)
		}
	}
}

func Test_UUID(t *testing.T) {
	buid := cnv.UUID()
	suid := cnv.UUID36(string(buid))
	fmt.Println("UUID:", suid)
}

func Test_Md5(t *testing.T) {
	var ar = []struct {
		s string
		c string
	}{
		{"abc", "900150983cd24fb0d6963f7d28e17f72"},
		{"next_2021", "4cef39adf6f590ef99f7e20335909ca9"},
		{"91958", "05b603ae1983baaadf874a070f695788"},
		{"erix2004", "bd91bbafa6b2731f8235c88927ed1b92"},
	}

	for _, a := range ar {
		b := []byte(a.s)
		c := cnv.Md5HexString(&b)
		fmt.Println("Md5HexString:", a.s, "c:", c)
		if c != a.c {
			t.Errorf("Md5HexString('%s'): soll '%s', ist '%s'", a.s, a.c, c)
		}
	}
}

func Test_IsDigit(t *testing.T) {
	var ar = []struct {
		s  string
		sa string
		ok bool
	}{
		{"1234567890737", "", true},
		{"12345678++-90737", "", false},
		{"12345678-90737", "-", true},
	}

	for _, a := range ar {
		ok := cnv.IsDigit(a.s, a.sa)
		fmt.Printf("IsDigit(%s)(%s) -> %v\n", a.s, a.sa, ok)
		if ok != a.ok {
			t.Errorf("IsDigit('%s'): soll %v, ist %v", a.s, a.ok, ok)
		}
	}
}

func Test_IsAlphanum(t *testing.T) {
	var ar = []struct {
		s  string
		sa string
		ok bool
	}{
		{"1234567äößü890737", "ßäöü", true},
		{"1234567äößü890737", "ß", false},
		{"123ab4X5678++-90737", "", false},
		{"12345678-90737", "-", true},
	}

	for _, a := range ar {
		ok := cnv.IsAlphaNum(a.s, a.sa)
		fmt.Printf("IsAlphaNum(%s)(%s) -> %v\n", a.s, a.sa, ok)
		if ok != a.ok {
			t.Errorf("IsAlphaNum('%s'): soll %v, ist %v", a.s, a.ok, ok)
		}
	}
}

func Test_parseTime(t *testing.T) {
	ar := []string{
		"24.11.1958",
		"1962/03/11",
		"1961.01.17",
		"03.12.1997 10:12:01",
		"2021.12.21 10:12:01",
		"2021-12-21 10:12:01",
		"2021-12-21T10:12:01",
		"01:12:01",
		"201-12-21#10:12:01",
	}

	fmt.Print("\n")
	for _, s := range ar {
		t := cnv.ParseTime(s, "")
		fmt.Println(t, " -> ", s)
		t = cnv.ParseTime(s, "UTC")
		fmt.Println(t, " -> ", s)
		fmt.Print("\n")
	}
}
