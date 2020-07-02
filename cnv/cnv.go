package cnv

// ----------------------------------------------------------------------------------
// cnv.go (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2019.11.24 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	xguid "github.com/google/uuid"
)

const timeLayout = "2006-01-02 15:04:05"

var locUTC, _ = time.LoadLocation("UTC")

// Unix2UTCTimeStr #
func Unix2UTCTimeStr(ut int64) string {

	t := time.Unix(ut, 0).In(locUTC)
	return t.Format(timeLayout)
}

// TimeUTC2Unix #
func TimeUTC2Unix(s string) int64 {
	if len(s) < 19 {
		return time.Now().In(locUTC).Unix()
	}

	t, _ := time.Parse(timeLayout, s[:19])
	return t.Unix()
}

// FTime #asString for FileName
func FTime() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

// DatAsInt #
func DatAsInt() int {
	return EsubStr2Int(FTime(), 0, 8)
}

// Estr2Int #
func Estr2Int(s string) int {
	return EsubStr2Int(s, 0, 19)
}

// EsubStr2Int #
func EsubStr2Int(s string, ix int, ilen int) int {
	b := []byte(s[ix:])
	l := len(s) - ix
	z := 0
	f := 1

	for i := 0; i < ilen && i < l; i++ {
		if b[i] >= '0' && b[i] <= '9' {
			z = z*10 + int(b[i]-'0')

		} else if b[i] == '-' {
			f = -1
		} else if b[i] == ';' {
			break
		}
	}

	return z * f
}

// Str2Dates #
func Str2Dates(s string) (int, int) {
	var (
		datv int
		datb int
	)

	// 20200520 or 20200510-20200520

	ss := strings.Split(s, "-")

	if len(ss) > 0 {
		datv = EsubStr2Int(ss[0], 0, 8)
		datb = datv
		if len(ss) > 1 {
			datb = EsubStr2Int(ss[1], 0, 8)
		}
	}

	return datv, datb
}

// BytesToUint64 #
func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// BytesToUint32 #
func BytesToUint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// Uint64ToHex #
func Uint64ToHex(u uint64) string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, u)

	return hex.EncodeToString(b)
}

// Uint32ToHex #
func Uint32ToHex(u uint32) string {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, u)

	return hex.EncodeToString(b)
}

// Uint32ToHexBytes #
func Uint32ToHexBytes(k uint32) []byte {
	b := make([]byte, 4)
	s := make([]byte, 8)
	binary.BigEndian.PutUint32(b, k)

	hex.Encode(s, b)

	return s
}

// Uint64ToHexBytes #
func Uint64ToHexBytes(k uint64) []byte {
	b := make([]byte, 8)
	s := make([]byte, 16)
	binary.BigEndian.PutUint64(b, k)

	hex.Encode(s, b)

	return s
}

// HexBytesToUint32 #
func HexBytesToUint32(s []byte) uint32 {
	b := make([]byte, 4)
	hex.Decode(b, s)
	return binary.BigEndian.Uint32(b)
}

// HexBytesToUint64 #
func HexBytesToUint64(s []byte) uint64 {
	b := make([]byte, 8)
	hex.Decode(b, s)
	return binary.BigEndian.Uint64(b)
}

// HexToBytes #
func HexToBytes(hx string) []byte {
	b, _ := hex.DecodeString(hx)
	return b
}

// BytesToHex #
func BytesToHex(b []byte) string {
	return hex.EncodeToString(b)
}

// UUID #
func UUID() []byte {
	b := xguid.New()
	r := make([]byte, 32)
	hex.Encode(r, b[:])
	return r
}

// IsHexString #
func IsHexString(s string) bool {
	b := []byte(s)
	le := len(b)
	if (le % 2) != 0 {
		return false
	}

	for i := 0; i < le; i++ {
		ok := (b[i] >= '0' && b[i] <= '9') ||
			(b[i] >= 'a' && b[i] <= 'f') ||
			(b[i] >= 'A' && b[i] <= 'F')

		if !ok {
			return false
		}
	}

	return true
}

// UUID36 #
func UUID36(uid string) string {
	suid := strings.Replace(uid, "-", "", -1)
	le := len(suid)
	if le < 32 {
		suid = suid + strings.Repeat("0", 32-le)
	}
	return fmt.Sprintf("%v-%v-%v-%v-%v", suid[0:8], suid[8:12], suid[12:16], suid[16:20], suid[20:32])
}

// StripUUID36 #
func StripUUID36(uid string) string {
	return strings.Replace(uid, "-", "", -1)
}

// ReverseString #
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Md5HexString #
func Md5HexString(b *[]byte) string {
	chk := md5.Sum(*b)
	return hex.EncodeToString(chk[:16])
}
