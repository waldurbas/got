package ftp

// ----------------------------------------------------------------------------------
// ftp.go (https://github.com/waldurbas/got)
// Copyright 2020,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.10.03 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/textproto"
	"strings"
	"time"
)

// Ftp #
type Ftp struct {
	con      *textproto.Conn
	Timeout  int
	Host     string
	skipEPSV bool
	mFeat    map[string]string
}

// Connect #
func Connect(conStr string, timeout int) (*Ftp, error) {
	//conStr := user + ":" + pwd + "@" + host

	ss := strings.Split(conStr, "@")
	if len(ss) != 2 {
		return nil, errors.New("host not defined")
	}

	host := ss[1]
	ss = strings.Split(ss[0], ":")
	if len(ss) != 2 {
		return nil, errors.New("user/pwd not defined")
	}

	user := ss[0]
	pwd := ss[1]

	f, err := DialFtp(host, timeout)

	if err != nil {
		return nil, err
	}

	err = f.Login(user, pwd)
	if err != nil {
		return nil, err
	}

	return f, nil
}

// DialFtp #
func DialFtp(addr string, timeout int) (*Ftp, error) {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)

	if err != nil {
		return nil, err
	}

	rAddr := conn.RemoteAddr().(*net.TCPAddr)
	var srcConn io.ReadWriteCloser = conn

	f := &Ftp{
		con:     textproto.NewConn(srcConn),
		Host:    rAddr.IP.String(),
		mFeat:   make(map[string]string),
		Timeout: timeout,
	}
	_, _, err = f.con.ReadResponse(220)
	if err != nil {
		f.Quit()
		return nil, err
	}

	return f, nil
}

// Quit #
func (f *Ftp) Quit() error {
	f.con.Cmd("QUIT")
	return f.con.Close()
}

// Login #
func (f *Ftp) Login(user, pwd string) error {

	c, m, err := f.cmd(-1, "USER %s", user)
	if err != nil {
		//lgx.PrintfDebug("Login.1: err: %v, c: [%v], m:[%v]", err, c, m)
		return err
	}

	switch c {
	case 230:
	case 331:
		_, _, err = f.cmd(230, "PASS %s", pwd)
		if err != nil {
			return err
		}
	default:
		return errors.New(m)
	}

	// check FEAT-command
	f.feat()

	// change to binary mode
	if _, _, err = f.cmd(200, "TYPE I"); err != nil {
		return err
	}

	return nil
}

// ListFiles #with MLSD
func (f *Ftp) ListFiles(path string) (*[]FileInfo, error) {
	c, err := f.cmdDataCon(0, "MLSD %s", path)
	if err != nil {
		return nil, err
	}

	r := &FResponse{con: c, ftp: f}
	defer r.Close()

	//	lgx.PrintfDebug("ListFiles.1 c: [%v]", c)

	var ff []FileInfo
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fi := parseEntry(scanner.Text())
		if fi != nil {
			ff = append(ff, *fi)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return &ff, err
}

// GetFile #
func (f *Ftp) GetFile(srcFile string, dstFile string) error {

	r, err := f.Retr(srcFile)
	if err != nil {
		return err
	}
	defer r.Close()

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dstFile, buf, 0666)
	return err
}
