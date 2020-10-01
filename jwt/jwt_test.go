package jwt_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/waldurbas/got/jwt"
)

func Test_1(t *testing.T) {
	x := jwt.New()

	x.Claims["user"] = "Waldemar"

	s := x.Encode("secretKey")
	fmt.Println("token:", s)

	d := jwt.New()
	if d.Parse(s) != nil {
		t.Errorf("Parse fail...")
		return
	}
	fmt.Println("parsed:", d)

	if d.Valid("secretKey") != nil {
		t.Errorf("not valid.")
		return
	}

	iat := d.Claims.AsInt64("iat")
	fmt.Println("iat:", time.Unix(iat, 0))

	expIn := int64(5)
	for {
		time.Sleep(1 * time.Second)
		fmt.Println(time.Now().Format("2006-01-02 15:04:05"),
			"expired:", time.Unix(iat+expIn, 0).Format("2006-01-02 15:04:05"),
			", at", time.Unix(iat, 0).Format("2006-01-02 15:04:05"))

		if d.Expired() > expIn {
			break
		}
	}
}

func Test_2(t *testing.T) {
	ss := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJhZG0iLCJjaWQiOjk5OTk5OTk5LCJleHAiOjE2MDE0NjU4NjQsImlhdCI6MTYwMTQwODI2NCwibWFpbCI6InVyYmFzQGV0b3MuZGUiLCJwZXJtIjoie1wiYWRtaW5cIjoxfSIsInN1YiI6IjA4MjY0LjE2MDE0IiwidWlkIjoiZWUyNDM1ZDAtMGYxZi01ZjE1LWUxMDAtNjM0MzlmNjZkNGRhIn0.a745b0a1b4ec8a03a0348a70f74d5872ad9eb6f07dd38586a6502ac03f5e8959",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHAiOiJuc2YiLCJjaWQiOjUzMDExLCJleHAiOjE2MDE0NjQ2MDcsImlhdCI6MTYwMTQwNzAwNywicGVybSI6IntcImxpc3RcIjoxLFwid3JpdGVcIjoxfSIsInN1YiI6IjA2ODAxLjE2MDE0IiwidWlkIjoiNjQzMDUzZjAtMDEzMC1jZGYwLTc2ZTctNjI3NWUyMGU4OWI5In0.9c25a129b786f3cb6f66a1d3e0786749b906d219c2dc0ab14bbb90df9aa23dca",
	}

	x := jwt.New()
	for _, s := range ss {
		x.Parse(s)
		fmt.Println(x.Claims)
	}
}
