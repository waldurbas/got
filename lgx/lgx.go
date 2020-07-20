package lgx

// ----------------------------------------------------------------------------------
// lgx.go (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.07.20 (wu) PathSplit for Windows
// 2020.07.06 (wu) SearchEmptyDirs,SearchFilesOlderAs,IsDirEmpty
// 2020.06.23 (wu) LgxDebug
// 2020.03.11 (wu) env.GCP
// 2020.02.10 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Lgx #
type Lgx struct {
	mu         sync.Mutex // ensures atomic writes; protects the following fields
	prop       int        // properties
	out        io.Writer  // destination for output
	buf        []byte
	logDir     string
	logFilePfx string
	curDir     string
}

// LGX_STD #Standard mit Time
// LGX_GCP #GoogleCloud ohne Time
const (
	LgxStd   = 0
	LgxGcp   = 1
	LgxDebug = 2
	LgxFile  = 4
)

// New #
func New(out io.Writer, prop int) *Lgx {
	sdir, _ := PathSplit(os.Args[0])
	return &Lgx{out: out, prop: prop, curDir: sdir}
}

// CurDir #
func (p *Lgx) CurDir() string {
	return p.curDir
}

// SetOutput sets the output destination for the logger.
func (p *Lgx) SetOutput(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out = w
}

func (p *Lgx) write(s string) string {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p._write(s)
}

func (p *Lgx) _write(s string) string {
	le := len(s)
	addNL := le == 0
	noNL := false
	if !addNL && s[le-1] != '\n' {
		addNL = true
	}

	t := time.Now()
	sti := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	p.buf = p.buf[:0]
	if le > 0 {
		if s[0] == '\r' || s[0] == '\n' {
			p.buf = append(p.buf, s[0])
			s = s[1:le]
			le--
			for le > 0 && (s[0] == '\r' || s[0] == '\n') {
				p.buf = append(p.buf, s[0])
				s = s[1:le]
				le--
			}
			p.out.Write(p.buf)
			p.buf = p.buf[:0]
		}

		if p.prop&LgxGcp == 0 {
			p.buf = append(p.buf, sti...)
		}

		if le > 0 {
			if s[le-1] == '#' {
				s = s[:le-1]
				le--
				noNL = true
			}
		}

		p.buf = append(p.buf, s...)
	}

	ss := string(p.buf)

	if addNL && !noNL {
		p.buf = append(p.buf, '\n')
	}

	p.out.Write(p.buf)

	if (p.prop & LgxFile) == LgxFile {
		if addNL && noNL {
			p.buf = append(p.buf, '\n')
		}

		sti = strings.ReplaceAll(sti[0:10], "-", "")
		logFileName := PathJoin(p.logDir, sti[0:4], sti[4:6])

		if CreateDirIfNotExist(logFileName) != -1 {
			logFileName = PathJoin(logFileName, p.logFilePfx+sti+".log")
			appendFile(logFileName, string(p.buf))
		}
	}

	return ss
}

// Fatal #
func (p *Lgx) Fatal(v ...interface{}) {
	p.write(fmt.Sprintln(v...))
	os.Exit(-1)
}

// Printf #
func (p *Lgx) Printf(frm string, v ...interface{}) string {
	return p.write(fmt.Sprintf(frm, v...))
}

// Println #
func (p *Lgx) Println(v ...interface{}) {
	p.write(fmt.Sprintln(v...))
}

// Print #
func (p *Lgx) Print(v ...interface{}) string {
	return p.write(fmt.Sprint(v...))
}

//------------- Standard ------------------------
var std = New(os.Stderr, 0)

// IsDebug #
var IsDebug = false

// CurDir #
func CurDir() string {
	return std.curDir
}

// Println #
func Println(v ...interface{}) {
	std.write(fmt.Sprintln(v...))
}

// Print #
func Print(v ...interface{}) string {
	return std.write(fmt.Sprint(v...))
}

// PrintDebug #
func PrintDebug(v ...interface{}) {
	if IsDebug || (std.prop&LgxDebug) == LgxDebug {
		std.write("[DEBUG] " + fmt.Sprintln(v...))
	}
}

// PrintInfo #
func PrintInfo(v ...interface{}) {
	std.write("[INFO] " + fmt.Sprintln(v...))
}

// PrintError #
func PrintError(v ...interface{}) {
	std.write("[ERROR] " + fmt.Sprintln(v...))
}

// Printf #
func Printf(format string, v ...interface{}) string {
	return std.write(fmt.Sprintf(format, v...))
}

// PrintfDebug #
func PrintfDebug(format string, v ...interface{}) {
	if IsDebug || (std.prop&LgxDebug) == LgxDebug {
		std.write("[DEBUG] " + fmt.Sprintf(format, v...))
	}
}

// PrintfInfo #
func PrintfInfo(format string, v ...interface{}) {
	std.write("[INFO] " + fmt.Sprintf(format, v...))
}

// PrintfError #
func PrintfError(format string, v ...interface{}) {
	std.write("[ERROR] " + fmt.Sprintf(format, v...))
}

// Fatal #
func Fatal(v ...interface{}) {
	std.write("[FATAL] " + fmt.Sprintln(v...))
	os.Exit(1)
}

// PathSplit # path.Split ist falsch fuer windows
func PathSplit(path string) (string, string) {
	b := strings.LastIndex(path, string(os.PathSeparator))
	a := b
	if a > 0 {
		a--
	}

	sd := path[:a+1]
	if sd == "" {
		sd = "."
	}

	return sd, path[b+1:]
}

// PathBase #
func PathBase(path string) string {
	_, f := PathSplit(path)
	return f
}

// PathDir #
func PathDir(path string) string {
	d, _ := PathSplit(path)
	return d
}

// PathJoin # path.Join ist falsch fuer Windows
func PathJoin(elem ...string) string {
	ps := string(os.PathSeparator)
	s := ""
	for _, e := range elem {
		if e != "" {
			if e[len(e)-1] == ps[0] {
				e = e[:len(e)-1]
			}

			if s != "" {
				if s == ps {
					s = s + e
				} else {
					s = s + ps + e
				}
			} else {
				s = e
				if s == "" {
					s = ps
				}
			}
		}
	}

	return s
}

// Start #
func Start(w io.Writer, info string, prop int, dir string, pfx string) {
	std.mu.Lock()
	defer std.mu.Unlock()

	IsDebug = atob(os.Getenv("DEBUG"))
	std.prop = prop

	std.logDir = dir
	std.logFilePfx = pfx
	if dir != "" {
		std.prop |= LgxFile
	}

	std._write("")
	if len(info) > 0 {
		std._write(info)
	}
}

// SetProp #
func SetProp(prop int) {
	std.mu.Lock()
	defer std.mu.Unlock()

	std.prop = prop
}

func atob(s string) bool {
	i64, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return false
	}

	return i64 > 0
}

func appendFile(path string, data string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return nil
}

// CreateDirIfNotExist #
func CreateDirIfNotExist(dirName string) int {
	if DirExists(dirName) {
		return 0
	}

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return -1
	}

	return 1
}

// DirExists #
func DirExists(path string) bool {
	f, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return f.IsDir()
}

// FileExists #
func FileExists(filename string) bool {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !f.IsDir()
}

// IsDirEmpty #
func IsDirEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true
	}

	return false
}

// SearchEmptyDirs #
func SearchEmptyDirs(dir string) *[]string {
	files := []string{}
	filepath.Walk(dir, func(path string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fInfo.IsDir() {
			if IsDirEmpty(path) {
				files = append(files, path)
			}
		}
		return nil
	})

	return &files
}

// SearchFilesOlderAs #
func SearchFilesOlderAs(dir string, days int) *[]string {
	timeBis := time.Now().UTC().Unix() - int64(days*86400)

	files := []string{}
	filepath.Walk(dir, func(path string, fInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fInfo.Mode().IsRegular() {
			if timeBis-fInfo.ModTime().UTC().Unix() > 0 {
				files = append(files, path)
			}
		}
		return nil
	})

	return &files
}
