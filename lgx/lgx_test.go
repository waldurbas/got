package lgx_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

func Test_Log(t *testing.T) {
	vf := lgx.PathJoin(path.Dir(os.Args[0]), "version.txt")
	fmt.Print("arg[0]: ", os.Args[0], ", vf: ", vf)

	ioutil.WriteFile(vf, []byte("1.2.3.4"), 0666)

	os.Setenv("GCP", "1")
	w := lgx.StartLog(os.Stderr, "/tmp", "(c) 2020 by Waldemar Urbas")
	w.SetOutput(nil)
	w.Print("nichts auf stdErr")
	fmt.Println("\nInfo:", lgx.Sversion)
	fmt.Println("\nLogfileName:", w.LogFileName)
	b, _ := ioutil.ReadFile(w.LogFileName)
	fmt.Printf("\n[%v]\n", string(b))
}
