package lgx

// ----------------------------------------------------------------------------------
// lgx.go (https://github.com/waldurbas/got)
// Copyright 2019,2022 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2022.03.10 (wu) PathJoinSep
// 2020.03.30 (wu) SetVersion
//                 ab go 1.16 funktioniert -ldflags "-X lgx.xVersion=$Version" nicht mehr ??
// 2020.03.24 (wu) func Write kompatibel mit io.Writer
// 2020.12.29 (wu) add LogDir(),ExecName()
// 2020.12.16 (wu) prgName bei StartLog wird automatisch ermittelt
// 2020.09.09 (wu) Info in Start without Datetime
// 2020.08.29 (wu) PrintLN
// 2020.07.20 (wu) PathSplit for Windows
// 2020.07.06 (wu) SearchEmptyDirs,SearchFilesOlderAs,IsDirEmpty
// 2020.06.23 (wu) LgxDebug
// 2020.03.11 (wu) env.GCP
// 2020.02.10 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// Version #go build -ldflags "-X lgx.xVersion=$Version"
	xVersion string

	// Sversion #wird benoetigt für Usage
	Sversion string

	// PrgName #ProgrammName ohne Ext.
	PrgName string

	// NewLine
	NewLine        string
	CRchar         string
	LFchar         string
	NewLinePrinted bool
	LinePfx        string
)

// Lgx #
type Lgx struct {
	mu          sync.Mutex // ensures atomic writes; protects the following fields
	prop        int        // properties
	out         io.Writer  // destination for output
	buf         []byte
	logFilePfx  string
	curDir      string // current Directory
	excName     string // execname without Directory
	LogDir      string
	LogFileName string
}

// LGX_STD #Standard mit Time
// LGX_GCP #GoogleCloud ohne Time
const (
	LgxStd   = 0
	LgxGcp   = 1
	LgxDebug = 2
	LgxFile  = 4

	NoTime = "!~!"
	NoNL   = '#'
)

func init() {
	CRchar = string([]byte{13})
	LFchar = string([]byte{10})

	if runtime.GOOS == "windows" {
		NewLine = string([]byte{13, 10})
	} else {
		NewLine = string([]byte{10})
	}
}

// New #
func New(out io.Writer, prop int) *Lgx {
	// currentDir and execName automatisch ermitteln
	s, err := os.Executable()
	if err != nil {
		s = os.Args[0]
	}

	sdir, exName := PathSplit(s)

	// prgName without Extension
	PrgName = strings.TrimSuffix(exName, path.Ext(exName))

	return &Lgx{out: out, prop: prop, curDir: sdir, excName: exName}
}

// Write # as io.Writer // ohne Uhrzeit
func (p *Lgx) Write(b []byte) (n int, err error) {
	le := len(b)
	if le > 0 {
		t := time.Now()
		sti := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		if p.out != nil {
			p.out.Write(b)
		}

		if (p.prop & LgxFile) == LgxFile {
			sti = strings.ReplaceAll(sti[0:10], "-", "")
			logFileName := PathJoin(p.LogDir, sti[0:4], sti[4:6])
			p.LogFileName = PathJoin(logFileName, p.logFilePfx+sti+".log")

			if CreateDirIfNotExist(logFileName) != -1 {
				appendFile(p.LogFileName, string(b))
			}
		}

		return le, nil
	}

	return 0, nil
}

// SetOutput #output destination for the logger.
func (p *Lgx) SetOutput(w io.Writer) io.Writer {
	p.mu.Lock()
	defer p.mu.Unlock()
	o := p.out
	p.out = w
	return o
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
		withTime := p.prop&LgxGcp == 0

		if le > 3 && s[:3] == NoTime {
			withTime = false
			s = s[3:]
			le -= 3
		}

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

		if withTime {
			p.buf = append(p.buf, sti...)
		}

		lx := len(LinePfx)
		if lx > 0 && lx < 7 {
			if LinePfx[lx-1] != ' ' {
				LinePfx = LinePfx + " "
			}
			p.buf = append(p.buf, LinePfx...)
		}

		if le > 0 {
			if s[le-1] == NoNL {
				s = s[:le-1]
				le--
				noNL = true
			}
		}

		p.buf = append(p.buf, s...)
	}

	ss := string(p.buf)

	if addNL && !noNL {
		p.buf = append(p.buf, NewLine...)
	}

	if p.out != nil {
		le := len(p.buf)
		if le > 0 && !NewLinePrinted {
			p.out.Write([]byte(NewLine))
		}

		p.out.Write(p.buf)
		if le > 0 {
			NewLinePrinted = p.buf[le-1] == 10
		}
	}

	if (p.prop & LgxFile) == LgxFile {
		sti = strings.ReplaceAll(sti[0:10], "-", "")
		logFileName := PathJoin(p.LogDir, sti[0:4], sti[4:6])
		p.LogFileName = PathJoin(logFileName, p.logFilePfx+sti+".log")

		if addNL && noNL {
			p.buf = append(p.buf, NewLine...)
		}

		if CreateDirIfNotExist(logFileName) != -1 {
			appendFile(p.LogFileName, string(p.buf))
		}
	}

	return ss
}

// Fatal #
func (p *Lgx) Fatal(v ...interface{}) {
	p.Println(v...)
	os.Exit(-1)
}

// Fatalf #
func (p *Lgx) Fatalf(frm string, v ...interface{}) {
	p.Printf(frm, v...)
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

// Default Logger
func Default() *Lgx {
	return std
}

// CurDir #
func CurDir() string {
	return std.curDir
}

// LogDir #
func LogDir() string {
	t := time.Now()
	sti := fmt.Sprintf("%d%02d", t.Year(), t.Month())
	sLog := PathJoin(std.LogDir, sti[:4], sti[4:])
	CreateDirIfNotExist(sLog)

	return sLog
}

// ExecName #
func ExecName() string {
	return std.excName
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

// Fatalf #
func Fatalf(format string, v ...interface{}) {
	std.write("[FATAL] " + fmt.Sprintf(format, v...))
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
	return PathJoinSep(string(os.PathSeparator), elem...)
}

func PathJoinSep(ps string, elem ...string) string {
	s := ""
	for _, e := range elem {
		lx := len(s)

		le := len(e)
		if lx > 0 {
			for le > 0 && e[0] == ps[0] {
				le--
				e = e[1:]
			}

			for le > 0 && e[le-1] == ps[0] {
				le--
				e = e[:le]
			}
		}

		if le < 1 {
			continue
		}

		if lx > 0 {
			for lx > 0 && s[lx-1] == ps[0] {
				lx--
				s = s[:lx]
			}
			s = s + ps + e
		} else {
			s = e
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
	std.out = w
	std.LogDir = dir
	std.logFilePfx = pfx
	if dir != "" {
		std.prop |= LgxFile
	}

	std._write("")
	if len(info) > 0 {
		std._write(NoTime + info)
	}
}

func printOut(w io.Writer, format string, v ...interface{}) {

	s := fmt.Sprintf(format, v...)
	if s == "" {
		Fprintf(w, NewLine)
	} else {
		Fprintf(w, s)
	}
}

// PrintNL #
func PrintNL() {
	Fprintf(std.out, NewLine)
	NewLinePrinted = true
}

// PrintStderr #
func PrintStderr(format string, v ...interface{}) {
	printOut(os.Stderr, format, v...)
}

// PrintStdout #
func PrintStdout(format string, v ...interface{}) {
	printOut(os.Stdout, format, v...)
}

// SetProp #
func SetProp(prop int) {
	std.mu.Lock()
	defer std.mu.Unlock()

	std.prop = prop
}

// SetOutput #liefert alten Writer
func SetOutput(w io.Writer) io.Writer {
	return std.SetOutput(w)
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

// StartLog #Parameter
// out: os.stderr || os.stdout
// ldir: z.B.: /usr/firma/log
// cpyRight: z.B.: "(c) 2020 by Waldemar Urbas"
//----------------------------------------------------------
// logfile unter {ldir}/{JAMO}/{prgname}{YYMMDD}.log
func StartLog(out *os.File, ldir string, cpyRight string) {
	prop := 0

	s := strings.Split(xVersion, ".")

	iGCP, e := strconv.Atoi(os.Getenv("GCP"))
	if e == nil || iGCP > 0 {
		if len(s) != 4 {
			bb, err := ioutil.ReadFile("./version.txt")
			if err == nil {
				xVersion = string(bb)
			}
		}

		ldir = ""
		prop |= LgxGcp
	}

	if ldir != "" {
		prop |= LgxFile
	}

	if len(s) != 4 {
		xVersion = "0.0.0.0"
	}

	Sversion = PrgName + " Version " + xVersion + " " + cpyRight

	Start(out, Sversion, prop, ldir, PrgName)
	PrintNL()
}

// Version #
func Version() string {
	return xVersion
}

func SetVersion(v string) {
	xVersion = v
}

// Fprintln #
func Fprintln(w io.Writer, a ...interface{}) {
	fmt.Fprint(w, a...)
	fmt.Fprint(w, NewLine)
}

// Fprintf #
func Fprintf(w io.Writer, format string, a ...interface{}) {
	if runtime.GOOS == "windows" {
		s := fmt.Sprintf(format, a...)
		s = strings.Replace(s, CRchar, "", -1)
		s = strings.Replace(s, LFchar, NewLine, -1)

		fmt.Fprintf(w, s)
		return
	}

	fmt.Fprintf(w, format, a...)
}
