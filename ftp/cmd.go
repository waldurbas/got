package ftp

// ----------------------------------------------------------------------------------
// cmd.go (https://github.com/waldurbas/got)
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
	"errors"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"time"
)

func (f *Ftp) cmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	_, err := f.con.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}

	return f.con.ReadResponse(expectCode)
}

func (f *Ftp) feat() error {
	c, m, err := f.cmd(-1, "FEAT")
	if err != nil {
		//		lgx.PrintfDebug("feat.1: err: %v, c: [%v], m:[%v]", err, c, m)
		return err
	}

	// FEAT-Cmd not supported
	if c != 211 {
		return nil
	}

	ss := strings.Split(m, "\n")

	for _, s := range ss {
		//		lgx.PrintfDebug("feat.2: [%v]", s)
		if s[:1] != " " {
			continue
		}

		s = strings.TrimSpace(s)
		items := strings.SplitN(s, " ", 2)
		c := ""
		if len(items) == 2 {
			c = items[1]
		}

		f.mFeat[items[0]] = c
		//		lgx.PrintfDebug("feat.3: [%s]:[%s]", items[0], f.cfeat[items[0]])
	}

	return nil
}

// epsv #get port for a data connection.
func (f *Ftp) epsv() (port int, err error) {
	_, line, err := f.cmd(229, "EPSV")
	if err != nil {
		return
	}

	start := strings.Index(line, "|||")
	end := strings.LastIndex(line, "|")
	if start == -1 || end == -1 {
		err = errors.New("invalid EPSV response format")
		return
	}
	port, err = strconv.Atoi(line[start+3 : end])
	return
}

func (f *Ftp) pasv() (host string, port int, err error) {
	_, line, err := f.cmd(227, "PASV")
	if err != nil {
		return
	}

	// PASV response format : 227 Entering Passive Mode (h1,h2,h3,h4,p1,p2).
	start := strings.Index(line, "(")
	end := strings.LastIndex(line, ")")
	if start == -1 || end == -1 {
		err = errors.New("invalid PASV response format")
		return
	}

	// We have to split the response string
	pasvData := strings.Split(line[start+1:end], ",")

	if len(pasvData) < 6 {
		err = errors.New("invalid PASV response format")
		return
	}

	// Let's compute the port number
	portPart1, err1 := strconv.Atoi(pasvData[4])
	if err1 != nil {
		err = err1
		return
	}

	portPart2, err2 := strconv.Atoi(pasvData[5])
	if err2 != nil {
		err = err2
		return
	}

	// Recompose port
	port = portPart1*256 + portPart2

	// Make the IP address to connect to
	host = strings.Join(pasvData[0:4], ".")
	return
}

// Retr #
func (f *Ftp) Retr(path string) (*FResponse, error) {
	conn, err := f.cmdDataCon(0, "RETR %s", path)
	if err != nil {
		return nil, err
	}

	return &FResponse{con: conn, ftp: f}, nil
}

func (f *Ftp) getPortNumber() (string, int, error) {
	if !f.skipEPSV {
		if port, err := f.epsv(); err == nil {
			return f.Host, port, nil
		}

		// if there is an error, skip EPSV for the next attempts
		f.skipEPSV = true
	}

	return f.pasv()
}

// newFtpConn #
func (f *Ftp) newFtpCon() (net.Conn, error) {
	host, port, err := f.getPortNumber()
	if err != nil {
		return nil, err
	}

	//	lgx.PrintfDebug("newFtpCon.host:[%v],port:[%v]", host, port)
	//	addr := host + ":" + strconv.Itoa(port)
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	conn, e := net.DialTimeout("tcp", addr, 10*time.Second)
	if e != nil {
		return nil, e
	}

	return conn, nil
}

func (f *Ftp) cmdDataCon(offset uint64, format string, args ...interface{}) (net.Conn, error) {
	conn, err := f.newFtpCon()
	if err != nil {
		return nil, err
	}

	if offset != 0 {
		// Pending
		_, _, err := f.cmd(350, "REST %d", offset)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	_, err = f.con.Cmd(format, args...)
	if err != nil {
		conn.Close()
		return nil, err
	}

	code, msg, err := f.con.ReadResponse(-1)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 125: already open; 150 already to send
	if code != 125 && code != 150 {
		conn.Close()
		return nil, &textproto.Error{Code: code, Msg: msg}
	}

	return conn, nil
}
