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
// 2020.06.23 (wu) LgxDebug
// 2020.03.11 (wu) env.GCP
// 2020.02.10 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"io"
	"os"
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
	return &Lgx{out: out, prop: prop}
}

// SetOutput sets the output destination for the logger.
func (p *Lgx) SetOutput(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out = w
}

func (p *Lgx) write(s string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	le := len(s)
	addNL := le == 0
	if !addNL && s[le-1] != '\n' {
		addNL = true
	}

	t := time.Now()
	sti := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d ",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())

	p.buf = p.buf[:0]

	if le > 0 {
		if p.prop&LgxGcp == 0 {
			p.buf = append(p.buf, sti...)
		}

		p.buf = append(p.buf, s...)
	}

	if addNL {
		p.buf = append(p.buf, '\n')
	}

	p.out.Write(p.buf)

	if (p.prop & LgxFile) == LgxFile {
		sti = strings.ReplaceAll(sti[0:10], "-", "")
		pSep := string(os.PathSeparator)
		logFileName := p.logDir + pSep + sti[0:4] + pSep + sti[4:6]

		if createDirIfNotExist(logFileName) {
			logFileName = logFileName + pSep + p.logFilePfx + sti + ".log"
			appendFile(logFileName, string(p.buf))
		}
	}
}

// Fatal #
func (p *Lgx) Fatal(v ...interface{}) {
	p.write(fmt.Sprintln(v...))
	os.Exit(-1)
}

// Printf #
func (p *Lgx) Printf(frm string, v ...interface{}) {
	p.write(fmt.Sprintf(frm, v...))
}

// Println #
func (p *Lgx) Println(v ...interface{}) {
	p.write(fmt.Sprintln(v...))
}

// Print #
func (p *Lgx) Print(v ...interface{}) {
	p.write(fmt.Sprint(v...))
}

//------------- Standard ------------------------
var std = New(os.Stderr, 0)
var isDebug = atob(os.Getenv("DEBUG"))

// Println #
func Println(v ...interface{}) {
	std.write(fmt.Sprintln(v...))
}

// Print #
func Print(v ...interface{}) {
	std.write(fmt.Sprint(v...))
}

// PrintDebug #
func PrintDebug(v ...interface{}) {
	if isDebug || (std.prop&LgxDebug) == LgxDebug {
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
func Printf(format string, v ...interface{}) {
	std.write(fmt.Sprintf(format, v...))
}

// PrintfDebug #
func PrintfDebug(format string, v ...interface{}) {
	if isDebug || (std.prop&LgxDebug) == LgxDebug {
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

// SetDefault #
func SetDefault(w io.Writer, prop int, dir string, pfx string) {
	std.mu.Lock()
	defer std.mu.Unlock()

	std.prop = prop

	std.logDir = dir
	std.logFilePfx = pfx
	if dir != "" {
		std.prop |= LgxFile
	}
}

// Start #
func Start(info string) {
	isDebug = atob(os.Getenv("DEBUG"))
	std.write("")
	if len(info) > 0 {
		Print(info)
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

	os.Chmod(path, 0666)
	defer f.Close()

	if _, err := f.WriteString(data); err != nil {
		return err
	}

	return nil
}

func createDirIfNotExist(dirName string) bool {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			return false
		}
	}

	return true
}
