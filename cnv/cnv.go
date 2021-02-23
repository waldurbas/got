package cnv

// ----------------------------------------------------------------------------------
// cnv.go (https://github.com/waldurbas/got)
// Copyright 2019,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.02.21 (wu) IsDigit,IsAlpha,IsAlphaNum
// 2020.07.19 (wu) PermitWeekday, Int2Prs, Int2DatHuman
// 2019.11.24 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	xguid "github.com/google/uuid"
)

const timeLayout = "2006-01-02 15:04:05"
const timeLayoutT = "2006-01-02T15:04:05"

var locUTC, _ = time.LoadLocation("UTC")
var locLOC, _ = time.LoadLocation("Local")

// RfillStr #
func RfillStr(s, ch string, le int) string {
	for {
		s += ch
		if len(s) > le {
			return s[0:le]
		}
	}
}

// LfillStr #
func LfillStr(s, ch string, le int) string {
	for {
		s = ch + s
		if len(s) > le {
			return s[0:le]
		}
	}
}

// Time2Str #
func Time2Str(t time.Time) string {
	return t.Format(timeLayout)
}

// Time2StrT #
func Time2StrT(t time.Time) string {
	return t.Format(timeLayoutT)
}

// Str2Time #
func Str2Time(s string) time.Time {
	var r time.Time
	if len(s) > 18 {
		if s[10:11] == "T" {
			r, _ = time.Parse(timeLayoutT, s[:19])
		} else {
			r, _ = time.Parse(timeLayout, s[:19])
		}
	}

	return r
}

// Unix2LocalTimeStr #
func Unix2LocalTimeStr(ut int64) string {
	t := time.Unix(ut, 0).In(locLOC)
	return t.Format(timeLayout)
}

// Unix2UTCTimeStr #
func Unix2UTCTimeStr(ut int64) string {
	t := time.Unix(ut, 0).In(locUTC)
	return t.Format(timeLayout)
}

// Unix2UTCTimeStrT #
func Unix2UTCTimeStrT(ut int64) string {
	t := time.Unix(ut, 0).In(locUTC)
	return t.Format(timeLayoutT)
}

// TimeUTC2Unix #
func TimeUTC2Unix(s string) int64 {
	if len(s) < 19 {
		return time.Now().In(locUTC).Unix()
	}

	t := Str2Time(s)
	return t.Unix()
}

// FTime #asString for FileName
func FTime() string {
	t := time.Now()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

// STime  #asString for Log
func STime(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

// TimeDif #
func TimeDif(tA time.Time, tL time.Time) (xs int, hh int, mm int, ss int) {
	dif := tL.Sub(tA)
	hh = int(dif.Hours())
	mm = int(dif.Minutes())
	ss = int(dif.Seconds())
	xs = ss

	if hh > 0 {
		mm -= hh * 60
		ss -= hh * 3600
	}

	if mm > 0 {
		ss -= mm * 60
	}

	return
}

// STimeDif #Differenz as String
func STimeDif(tA time.Time, tL time.Time) string {

	_, hh, mm, ss := TimeDif(tA, tL)
	s := fmt.Sprintf("%.2d:%.2d:%.2d", hh, mm, ss)
	return s
}

// DatAsInt #
func DatAsInt() int {
	return EsubStr2Int(FTime(), 0, 8)
}

// Str2Dat #
func Str2Dat(s string) int {
	sep := "."
	ix := strings.Index(s, sep)
	if ix < 0 {
		sep = "-"
		ix = strings.Index(s, sep)
	}

	if ix < 0 {
		return 0
	}

	d := strings.Split(s, sep)
	if len(d) == 3 {

		// jahr vorne
		if ix == 4 {
			return EsubStr2Int(d[0], 0, 4)*10000 + EsubStr2Int(d[1], 0, 2)*100 + EsubStr2Int(d[2], 0, 2)
		}

		return EsubStr2Int(d[2], 0, 4)*10000 + EsubStr2Int(d[1], 0, 2)*100 + EsubStr2Int(d[0], 0, 2)
	}

	return 0
}

// Int2Dat #
func Int2Dat(d int) string {
	out := make([]byte, 10)
	in := fmt.Sprintf("%.8d", d)

	for i := 0; i < 4; i++ {
		out[i] = in[i]
	}

	out[4] = '-'
	for i := 4; i < 6; i++ {
		out[i+1] = in[i]
	}

	out[7] = '-'
	for i := 6; i < 8; i++ {
		out[i+2] = in[i]
	}

	return string(out)
}

// Int2Prs #
func Int2Prs(ns int) string {
	in := strconv.Itoa(ns)

	le := len(in)

	for le < 3 {
		in = "0" + in
		le = len(in)
	}

	out := in[0:le-2] + "." + in[le-2:le]
	return out
}

// Int2DatHuman #
func Int2DatHuman(d int) string {
	out := make([]byte, 10)
	in := fmt.Sprintf("%.8d", d)

	for i := 0; i < 2; i++ {
		out[i] = in[6+i]
	}

	out[2] = '-'
	for i := 2; i < 4; i++ {
		out[i+1] = in[2+i]
	}

	out[5] = '-'
	for i := 4; i < 8; i++ {
		out[i+2] = in[i-4]
	}

	return string(out)
}

// FormatInt64 #Format Integer mit Tausend Points
func FormatInt64(n int64) string {
	in := strconv.FormatInt(n, 10)
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = '.'
		}
	}
}

// FormatInt #Format Integer mit Tausend Points
func FormatInt(n int) string {
	return FormatInt64(int64(n))
}

// Bool2Int #
func Bool2Int(b bool) int {
	// The compiler currently only optimizes this form. See issue 6011.
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}

// Estr2Int #
func Estr2Int(s string) int {
	return EsubStr2Int(s, 0, 19)
}

// EsubStr #
func EsubStr(s string, ix int, le int) string {
	l := len(s)

	if ix > l {
		return ""
	}

	if (ix + le) > l {
		le = l - ix
	}

	b := s[ix : ix+le]
	return b
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

// IsDigit #
func IsDigit(s string, any string) bool {
	b := []byte(s)
	a := []byte(any)
	for i := 0; i < len(b); i++ {
		if b[i] < '0' || b[i] > '9' {
			if bytes.LastIndexByte(a, b[i]) == -1 {
				return false
			}
		}
	}

	return true
}

// IsAlpha #
func IsAlpha(s string, any string) bool {
	b := []byte(s)
	a := []byte(any)

	for i := 0; i < len(b); i++ {
		if b[i] < 'A' || b[i] > 'Z' {
			if b[i] < 'a' || b[i] > 'z' {
				if bytes.LastIndexByte(a, b[i]) == -1 {
					return false
				}
			}
		}
	}

	return true
}

// IsAlphaNum #
func IsAlphaNum(s string, any string) bool {
	b := []byte(s)
	a := []byte(any)

	for i := 0; i < len(b); i++ {
		if b[i] < 'A' || b[i] > 'Z' {
			if b[i] < 'a' || b[i] > 'z' {
				if b[i] < '0' || b[i] > '9' {
					if bytes.LastIndexByte(a, b[i]) == -1 {
						return false
					}
				}
			}
		}
	}

	return true
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

// ToUTF8 #ISO8859_1 to UTF8
func ToUTF8(s string) string {

	iso8859Buf := []byte(s)

	buf := make([]rune, len(iso8859Buf))
	for i, b := range iso8859Buf {
		if b == 0x80 {
			buf[i] = '€'
		} else {
			buf[i] = rune(b)
		}
	}
	return string(buf)
}

// ToAnsi #UTF8 to ANSI
func ToAnsi(buf *[]byte) []byte {
	ansiBuf := make([]byte, len(*buf))

	a := 0
	for i := 0; i < len(*buf); i++ {
		switch (*buf)[i] {
		case 0xe2: // € = e2 82 ac
			i++
			if (*buf)[i] == 0x82 {
				i++
				if (*buf)[i] == 0xac {
					ansiBuf[a] = 0x80
					a++
				}
			}
		case 0xc2:
			i++
			ansiBuf[a] = (*buf)[i]
			a++
		case 0xc3:
			i++
			ansiBuf[a] = (*buf)[i] + 0x40
			a++
		default:
			ansiBuf[a] = (*buf)[i]
			a++
		}
	}

	return ansiBuf[:a]
}

// BitIsSet #
func BitIsSet(b, flag uint) bool { return b&flag != 0 }

// BitSet #
func BitSet(b, flag uint) uint { return b | flag }

// BitClear #
func BitClear(b, flag uint) uint { return b &^ flag }

// BitToggle #
func BitToggle(b, flag uint) uint { return b ^ flag }

// ReadableBytes #
func ReadableBytes(n uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	b := float64(1024)
	e := math.Floor(math.Log(float64(n)) / math.Log(b))
	sfx := sizes[int(e)]
	v := float64(n) / math.Pow(b, math.Floor(e))
	f := "%.0f"
	if v < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", v, sfx)
}

// GetVersion #
func GetVersion(ss string) string {

	s := strings.Split(ss, ".")

	if len(s) != 4 {
		return "0.0.0.0"
	}

	var v [4]int

	for i := 0; i < 4; i++ {
		v[i] = EsubStr2Int(s[i], 0, 4)
	}

	return strconv.Itoa(v[0]) + "." +
		strconv.Itoa(v[1]) + "." +
		strconv.Itoa(v[2]) + "." +
		strconv.Itoa(v[3])
}

// GetVersionAsInt #
func GetVersionAsInt(ss string) int {

	s := strings.Split(ss, ".")

	if len(s) != 4 {
		return 0
	}

	v := 0
	for i := 0; i < 4; i++ {
		v = v*100 + EsubStr2Int(s[i], 0, 4)
	}

	return v
}

// PermitWeekDay for
func PermitWeekDay(t time.Time, sDays []string) bool {
	ih := int(t.Weekday())
	ok := false
	for i := 0; i < len(sDays) && !ok; i++ {
		switch strings.ToLower(sDays[i]) {
		case "mo", "1":
			ok = ih == 1
		case "di", "2":
			ok = ih == 2
		case "mi", "3":
			ok = ih == 3
		case "do", "4":
			ok = ih == 4
		case "fr", "5":
			ok = ih == 5
		case "sa", "6":
			ok = ih == 6
		case "so", "0":
			ok = ih == 0
		}
	}

	return ok
}

// PermitHour # array: [ "12:00-18:00","1400-2200"]
func PermitHour(t time.Time, sh []string) bool {
	tt := t.Hour()*100 + t.Minute()
	ok := false

	var vt int
	var bt int
	for i := 0; i < len(sh) && !ok; i++ {
		s := strings.Split(sh[i], "-")

		if len(s) == 1 {
			vt = EsubStr2Int(s[0], 0, 5)
			bt = 2400
		} else {
			vt = EsubStr2Int(s[0], 0, 5)
			bt = EsubStr2Int(s[1], 0, 5)
		}

		ok = tt >= vt && tt <= bt
	}

	return ok
}

// DelimTextAdd #
func DelimTextAdd(ss *string, s, delim string) {
	if len(s) < 1 {
		return
	}

	if len(*ss) > 0 {
		*ss = *ss + delim + s
		return
	}

	*ss = s
}
