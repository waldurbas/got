package ftp

// ----------------------------------------------------------------------------------
// resp.go (https://github.com/waldurbas/got)
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
	"net"
	"time"
)

// FResponse #
type FResponse struct {
	con    net.Conn
	ftp    *Ftp
	closed bool
}

// Read #
func (r *FResponse) Read(buf []byte) (int, error) {
	return r.con.Read(buf)
}

// Close #
func (r *FResponse) Close() error {
	if r.closed {
		return nil
	}
	err := r.con.Close()
	_, _, err2 := r.ftp.con.ReadResponse(226)
	if err2 != nil {
		err = err2
	}
	r.closed = true
	return err
}

// SetDeadline #
func (r *FResponse) SetDeadline(t time.Time) error {
	return r.con.SetDeadline(t)
}
