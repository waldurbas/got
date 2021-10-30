package ecv

// ----------------------------------------------------------------------------------
// ecv.go (https://github.com/waldurbas/ecv)
// Copyright 2019,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.02.11 EcvFile.Count, GetTable,NewEcvFile
// 2019.05.20 Init
//-----------------------------------------------------------------------------------

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"
)

// EcvFile #
type EcvFile struct {
	FileName string
	Tables   []*EcvTable
	UTF8     bool
}

// EcvTable #
type EcvTable struct {
	KeyFound   bool
	iSearchKey int
	iSearchCol int
	CurrentPos int
	curFields  *[]string

	Table    string
	Header   string
	IndexOf  map[string]int
	keyIndex []int
	Fields   []EcvField
	data     []*ecvEntry
	Count    int
}

// EcvType #
type EcvType int

// Reader #
type Reader func(*string) bool

// ecv Fieldtype
const (
	EcvInt EcvType = iota
	EcvStr
)

//ecv-Entry
type ecvEntry struct {
	F []string
}

// EcvField #
type EcvField struct {
	Idx  int
	Name string
	Typ  EcvType
}

var iso8859run [256]rune

// init routine
func init() {
	for i := 0; i < 256; i++ {
		iso8859run[i] = rune(i)
	}

	iso8859run[0xc2] = 0
	iso8859run[0x80] = '€'
	iso8859run[0xa9] = '©'
	iso8859run[0xab] = '«'
	iso8859run[0xae] = '®'
	iso8859run[0xbb] = '»'
	iso8859run[0xbc] = '¼'
	iso8859run[0xbd] = '½'
	iso8859run[0xbe] = '¾'
	iso8859run[0xc4] = 'Ä'
	iso8859run[0xd6] = 'Ö'
	iso8859run[0xdc] = 'Ü'
	iso8859run[0xdf] = 'ß'
	iso8859run[0xe4] = 'ä'
	iso8859run[0xf6] = 'ö'
	iso8859run[0xfc] = 'ü'
}

// Open #
func (t *EcvTable) Open() bool {
	if t.Count == 0 {
		return false
	}

	t.CurrentPos = 0
	return true
}

func (t *EcvTable) checkLine(ipos int) bool {
	if ipos < t.Count {
		t.CurrentPos = ipos
		t.curFields = &t.data[t.CurrentPos].F
		return true
	}

	return false
}

// Seek to Line
func (t *EcvTable) Seek(ipos int) bool {
	return t.checkLine(ipos)
}

// First #
func (t *EcvTable) First() bool {
	return t.checkLine(0)
}

// Fetch #
func (t *EcvTable) Fetch() bool {
	if t.checkLine(t.CurrentPos) {
		t.CurrentPos = t.CurrentPos + 1
		return true
	}

	return false
}

// IfieldByName #
func (t *EcvTable) IfieldByName(s string) int {
	fix := t.IndexOf[s]
	return t.AsInteger(fix)
}

// SfieldByName #
func (t *EcvTable) SfieldByName(s string) string {
	fix := t.IndexOf[s]
	return t.AsString(fix)
}

// AsInteger #
func (t *EcvTable) AsInteger(fix int) int {

	if fix >= 0 && len(*t.curFields) > fix {
		v, _ := strconv.Atoi((*t.curFields)[fix])
		return v
	}

	return 0
}

// AsString #
func (t *EcvTable) AsString(fix int) string {
	if fix >= 0 && len(*t.curFields) > fix {
		return (*t.curFields)[fix]
	}

	return ""
}

// IsNull #
func (t *EcvTable) IsNull(fix int) bool {
	if fix >= 0 && len(*t.curFields) > fix {
		return (*t.curFields)[fix] == "NULL"
	}

	return false
}

func (t *EcvTable) setIndex(sidx string) {
	ss := strings.Split(sidx, ",")
	t.keyIndex = t.keyIndex[:0]
	for _, item := range ss {
		t.keyIndex = append(t.keyIndex, t.IndexOf[item])
	}
}

// Sort #mit setIndex
func (t *EcvTable) Sort(sidx string) {
	t.setIndex(sidx)
	sort.Slice(t.data, func(ii, jj int) bool {

		var v1 int
		var v2 int

		//		a := strings.Split(e.data[ii].Line, "^")
		//		b := strings.Split(e.data[jj].Line, "^")
		a := &t.data[ii].F
		b := &t.data[jj].F

		for i := 0; i < len(t.keyIndex); i++ {
			ix := t.keyIndex[i]

			if t.Fields[ix].Typ == EcvInt {
				v1, _ = strconv.Atoi((*a)[ix])
				v2, _ = strconv.Atoi((*b)[ix])

				if v1 < v2 {
					return true
				}
				if v1 > v2 {
					return false
				}
			} else {
				if (*a)[ix] < (*b)[ix] {
					return true
				}
				if (*a)[ix] > (*b)[ix] {
					return false
				}

			}
		}
		return false
	})
}

// FindFirstInt #Key
func (t *EcvTable) FindFirstInt(key int) int {
	var d int
	var v int

	fidx := t.keyIndex[0]

	t.iSearchKey = key
	t.iSearchCol = fidx

	a := 0
	e := t.Count - 1

	for {
		d = (a + e) >> 1
		t.checkLine(d)
		v = t.AsInteger(fidx) - key

		if v > 0 {
			e = d - 1
		} else {
			a = d + 1
		}

		if v == 0 || a > e {
			break
		}

	}

	if v == 0 && d > 0 {
		v = 0
		a = d

		for {

			if (a <= 0) || (v != 0) {
				break
			}

			a = a - 1

			t.checkLine(a)
			v = t.AsInteger(fidx) - key
			if v == 0 {
				d = a
			}
		}

		v = 0
	}

	if (v == 0) && (d >= 0) {
		t.checkLine(d)
		t.KeyFound = true

		return d
	}

	t.KeyFound = false
	return -1
}

// FindNextInt #Key
func (t *EcvTable) FindNextInt() bool {

	t.CurrentPos++
	if t.CurrentPos >= t.Count {
		t.KeyFound = false
		return false
	}

	t.checkLine(t.CurrentPos)
	if t.KeyFound && t.CurrentPos < t.Count && t.AsInteger(t.iSearchCol) == t.iSearchKey {
		return true
	}

	t.KeyFound = false
	return false
}

type fileReader struct {
	s *bufio.Scanner
}

func (f *fileReader) readFile(line *string) bool {
	if f.s.Scan() {
		*line = f.s.Text()
		return true
	}

	return false
}

// NewEcvFile #
func NewEcvFile() *EcvFile {
	return &EcvFile{}
}

// Count #
func (ef *EcvFile) Count() int {
	return len(ef.Tables)
}

// GetTable #
func (ef *EcvFile) GetTable(table string) *EcvTable {
	for i := 0; i < ef.Count(); i++ {
		if ef.Tables[i].Table == table {
			return ef.Tables[i]
		}
	}

	return nil
}

// Load #Data
func (ef *EcvFile) Load(filePath string) error {
	ef.FileName = filePath

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	f.Seek(0, 0)

	var fr fileReader
	fr.s = bufio.NewScanner(f)
	err = ef.LoadData(fr.readFile)

	return fr.s.Err()
}

// Clear #
func (ef *EcvFile) Clear() {
	ef.Tables = ef.Tables[:0]
}

// LoadData #
func (ef *EcvFile) LoadData(reader Reader) error {
	var line string
	var e *EcvTable
	var li *ecvEntry

	ef.Clear()

	for reader(&line) {
		if line[0:1] == "@" {
			fields := strings.Split(line, ",")

			e = new(EcvTable)
			e.Table = fields[0][1:]
			e.Header = line
			e.Count = 0
			e.CurrentPos = 0
			e.IndexOf = make(map[string]int)

			ef.Tables = append(ef.Tables, e)

			e.Fields = make([]EcvField, len(fields)-1)
			for i := 0; i < len(e.Fields); i++ {
				sf := fields[i+1]

				e.Fields[i].Idx = i

				w := strings.FieldsFunc(sf, func(r rune) bool {
					return r == '[' || r == ']'
				})

				e.Fields[i].Name = w[0]
				if len(w) == 2 {
					switch {
					case w[1] == "int":
						e.Fields[i].Typ = EcvInt
					case w[1] == "str":
						e.Fields[i].Typ = EcvStr
					default:
						e.Fields[i].Typ = EcvStr
					}
				} else {
					e.Fields[i].Typ = EcvStr
				}

				/*
					if strings.HasSuffix(sf, "[int]") {
						e.Fields[i].Name = sf[0 : len(sf)-5]
						e.Fields[i].Typ = EcvInt
					} else if strings.HasSuffix(sf, "[str]") {
						e.Fields[i].Name = sf[0 : len(sf)-5]
						e.Fields[i].Typ = EcvStr
					}
				*/
				e.IndexOf[e.Fields[i].Name] = i
			}
		} else {
			e.Count++

			li = new(ecvEntry)
			if ef.UTF8 {
				li.F = strings.Split(line, "^")
			} else {
				li.F = strings.Split(toUTF8(line), "^")

			}

			e.data = append(e.data, li)
		}
	}

	return nil
}

// ISO8859_1 to UTF8
func toUTF8(s string) string {
	bIso8859Eins := []byte(s)

	buf := make([]rune, len(bIso8859Eins))
	a := 0
	for i := 0; i < len(bIso8859Eins); i++ {
		buf[a] = iso8859run[bIso8859Eins[i]]
		if buf[a] > 0 {
			a++
		}
	}

	/*
		b := bIso8859Eins[i]
		switch b {
		case 0xc2:
		case 0x80:
			buf[a] = '€'
			a++
		case 0xbd:
			buf[a] = '½'
			a++
		default:
			buf[a] = rune(b)
			a++
		}
	*/

	return string(buf[:a])
}
