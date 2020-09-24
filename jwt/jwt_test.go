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

	expIn := int64(20)
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
