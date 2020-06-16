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
// 2020.03.11 (wu) env.GCP
// 2020.02.10 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Lgx #
type Lgx struct {
	mu   sync.Mutex // ensures atomic writes; protects the following fields
	pfx  string     // prefix to write at beginning of each line
	prop int        // properties
	out  io.Writer  // destination for output
	buf  []byte
}

// LGX_STD #Standard mit Time
// LGX_GCP #GoogleCloud ohne Time
const (
	LgxStd = 0
	LgxGcp = 1
)

// New #
func New(out io.Writer, pfx string, prop int) *Lgx {
	return &Lgx{out: out, prop: prop, pfx: setPfx(pfx)}
}

// SetOutput sets the output destination for the logger.
func (p *Lgx) SetOutput(w io.Writer) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out = w
}

func setPfx(pfx string) string {
	if len(pfx) > 0 {
		return "[" + pfx + "] "
	}

	return ""
}

// SetPrefix #
func (p *Lgx) SetPrefix(pfx string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pfx = setPfx(pfx)
}

func (p *Lgx) write(s string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buf = p.buf[:0]

	if p.prop&LgxGcp == 0 {
		t := time.Now()
		ss := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d ",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())

		p.buf = append(p.buf, ss...)
	}

	if len(p.pfx) > 0 {
		p.buf = append(p.buf, p.pfx...)
	}

	p.buf = append(p.buf, s...)

	if len(s) == 0 || s[len(s)-1] != '\n' {
		p.buf = append(p.buf, '\n')
	}

	p.out.Write(p.buf)
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

// Print #
func (p *Lgx) Print(v ...interface{}) {
	p.write(fmt.Sprintln(v...))
}

var std = New(os.Stderr, "", 0)

// Print #
func Print(v ...interface{}) {
	std.write(fmt.Sprintln(v...))
}

// PrintDebug #
func PrintDebug(v ...interface{}) {
	std.write("[DEBUG] " + fmt.Sprintln(v...))
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
	std.write("[DEBUG] " + fmt.Sprintf(format, v...))
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
func SetDefault(w io.Writer, pfx string, prop int) {
	std.mu.Lock()
	defer std.mu.Unlock()

	std.pfx = setPfx(pfx)
	std.prop = prop
}
