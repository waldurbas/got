package xtl_test

import (
	"fmt"
	"testing"

	"github.com/waldurbas/got/xtl"
)

func Test_x1(t *testing.T) {

	sdir := "/usr/local/bin"
	fi, err := xtl.LoadFiles(sdir, "*")
	if err != nil {
		t.Errorf("Error: LoadFiles(%s):%v", sdir, err)
	}

	for _, f := range *fi {
		mt := f.Time.Format("2006-01-02 15:04:05 UTC")
		fmt.Println(mt, f.FileName)
	}
}
