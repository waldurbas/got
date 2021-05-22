package htf

// ----------------------------------------------------------------------------------
// counter.go (https://github.com/waldurbas/got/htf)
// Copyright 2018,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2018.12.11 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"strings"
)

// WriteCounter #
type WriteCounter struct {
	Total uint64
}

// Write #
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.PrintProgress()
	return n, nil
}

// PrintProgress #prints the progress of a file write
func (wc WriteCounter) PrintProgress() {
	fmt.Printf("\r%s\rDownloading... %s complete", strings.Repeat(" ", 50), readableBytes(wc.Total))
}
