package jwt

// ----------------------------------------------------------------------------------
// jwt.go (https://github.com/waldurbas/got): base implementation JWT token
// Copyright 2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.09.06 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

var (
	// ErrTokenExpired #
	ErrTokenExpired = errors.New("Token is expired")
	// ErrTokenNotValid #
	ErrTokenNotValid = errors.New("Token is not valid")
)

// TMap #
type TMap map[string]interface{}

// JWToken with HS256
type JWToken struct {
	RawData string
	Head    TMap
	Claims  TMap
}

// New #
func New() *JWToken {
	return &JWToken{
		Head:   TMap{"alg": "HS256", "typ": "JWT"},
		Claims: TMap{"iat": time.Now().UTC().Unix()},
	}
}

// AsInt64 #
func (d *TMap) AsInt64(key string) int64 {
	switch ii := (*d)[key].(type) {
	case int64:
		return ii
	case float64:
		return int64(ii)
	case json.Number:
		v, _ := ii.Int64()
		return v
	default:
		//		fmt.Printf("AsInt64: %v, %T\n", ii, ii)
	}

	return 0
}

// AsBase64 #
func (d *TMap) AsBase64() string {
	bb, _ := json.Marshal(d)

	return base64.RawStdEncoding.EncodeToString(bb)
}

// Assign #
func (d *TMap) Assign(m *map[string]interface{}) {
	for k, v := range *m {
		(*d)[k] = v
	}
}

// Clear #
func (t *JWToken) Clear() {
	t.Claims = TMap{"iat": time.Now().UTC().Unix()}
}

// Encode #
func (t *JWToken) Encode(key string) string {
	shdr := t.Head.AsBase64()
	spay := t.Claims.AsBase64()
	sign := hmacSha256(shdr+"."+spay, key)

	t.RawData = shdr + "." + spay + "." + sign
	return t.RawData
}

// Expired #
func (t *JWToken) Expired() int64 {
	return time.Now().UTC().Unix() - t.Claims.AsInt64("iat")
}

// Parse #
func (t *JWToken) Parse(s string) error {
	ss := strings.Split(s, " ")
	if len(ss) == 2 {
		if ss[0] != "Bearer" && ss[0] != "JWT" {
			return ErrTokenNotValid
		}
		s = ss[1]
	}

	ss = strings.Split(s, ".")
	if len(ss) != 3 {
		return ErrTokenNotValid
	}

	// Header
	mHead, err := base64ToMap(ss[0])
	if err != nil {
		return ErrTokenNotValid
	}

	jHead, _ := json.Marshal(mHead)
	cHead, _ := json.Marshal(&t.Head)

	if string(jHead) != string(cHead) {
		return ErrTokenNotValid
	}

	// Payload
	m, err := base64ToMap(ss[1])
	if err != nil {
		return ErrTokenNotValid
	}

	t.Clear()
	for k, v := range *m {
		t.Claims[k] = v
	}
	t.RawData = ss[0] + "." + ss[1] + "." + ss[2]

	return nil
}

// Valid #
func (t *JWToken) Valid(key string) error {
	ss := strings.Split(t.RawData, ".")
	if len(ss) != 3 {
		return ErrTokenNotValid
	}

	// Signature
	sign := hmacSha256(ss[0]+"."+ss[1], key)
	if sign != ss[2] {
		return ErrTokenNotValid
	}

	return nil
}

func base64ToMap(s string) (*map[string]interface{}, error) {
	ds, err := base64.RawStdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}

	m := &map[string]interface{}{}
	if err := json.Unmarshal(ds, m); err != nil {
		return nil, err
	}

	return m, nil
}

func hmacSha256(data string, key string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(key))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}
