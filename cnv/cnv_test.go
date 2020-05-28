package cnv_test

import (
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
