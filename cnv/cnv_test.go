package cnv_test

import (
	"fmt"
	"testing"

	"github.com/waldurbas/got/cnv"
)

func Test_checkTime(t *testing.T) {
	var dtest = []struct {
		s string
		u int64
	}{
		{"2020-05-13 23:59:06", 1589414346},
		{"2000-01-01 00:00:00", 946684800},
		{"2000-01-01 00:00:00", 946684800},
		{"2020-05-27 07:24:06", 1590564246},
	}

	//	tm := time.Now()
	//a := tm.Unix()
	//t.Errorf("unix: %v, %v", a, tm)

	for _, tt := range dtest {
		u := cnv.TimeUTC2Unix(tt.s)

		if u != tt.u {
			t.Errorf("TimeUTC2Unix(%s): soll %d, ist %d", tt.s, tt.u, u)
		}

		s := cnv.Unix2UTCTimeStr(u)
		if s != tt.s {
			t.Errorf("Unix2UTCTimeStr(%d): soll %s, ist %s", u, tt.s, s)
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

func Test_int2dat(t *testing.T) {
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

func Test_formatInt(t *testing.T) {
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
		{"1321.001", 1321001},
	}

	for _, tt := range dtest {
		ss := cnv.FormatInt(tt.u)

		if ss != tt.s {
			t.Errorf("FormatInt(%d): soll %s, ist %s", tt.u, tt.s, ss)
		}
	}
}
