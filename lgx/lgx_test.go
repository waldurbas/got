package lgx_test

import (
	"testing"
	"time"

	lgx "github.com/waldurbas/got/lgx"
)

func Test_checkLGX(t *testing.T) {
	lgx.Printf("printf : %s %d %v", "text", 22, time.Now())
	lgx.Print("println:", "text", 22, time.Now())
}
