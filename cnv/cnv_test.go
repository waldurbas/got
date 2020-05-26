package cnv_test

import (
	"testing"

	"github.com/waldurbas/got/cnv"
)

func Test_checkLGX(t *testing.T) {

	s := "20190501-20200613"
	a := cnv.EsubStr2Int(s, 0, 8)
	b := cnv.EsubStr2Int(s, 9, 8)

	if a != 20190501 || b != 20200613 {
		t.Fatalf("s=[%s], a=[%d], b=[%d]", s, a, b)
	}
}
