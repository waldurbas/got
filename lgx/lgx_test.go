package lgx_test

import (
	"testing"
	"time"

	lgx "github.com/waldurbas/got/lgx"
)

func Test_checkLGX(t *testing.T) {
	lgx.Printf("printf : %s %d %v", "text", 22, time.Now())
	lgx.Print("print:", "text", 22, time.Now())
}

func Test_Path(t *testing.T) {
	var dtest = []struct {
		s []string
		u string
	}{
		{[]string{"/", "abc/d", "x.dat"}, "/abc/d/x.dat"},
		{[]string{"abc/d", "/def/", "//x.dat"}, "abc/d/def/x.dat"},
	}

	for _, tt := range dtest {

		ss := lgx.PathJoin(tt.s[0], tt.s[1], tt.s[2])

		if ss != tt.u {
			t.Errorf("PathJoin: soll %s, ist %s", tt.u, ss)
		}
	}
}
