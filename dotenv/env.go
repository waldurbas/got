package dotenv

// ----------------------------------------------------------------------------------
// env.go (https://github.com/waldurbas/got)
// Copyright 2019,2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2019.11.24 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"io/ioutil"
	"os"
	"strings"
)

// Load #
func Load() {
	fname := ".env"
	b, err := ioutil.ReadFile(fname)

	if err != nil {
		return
	}

	lines := strings.Split(string(b), "\n")
	for _, s := range lines {
		ex := strings.Split(s, "=")
		if len(ex) == 2 && len(ex[0]) > 1 {
			v := checkValue(ex[1])

			if v != "" {
				os.Setenv(ex[0], v)
			}
		}
	}
}

// checkValue #
func checkValue(s string) string {
	v := strings.Trim(s, " ")

	le := len(v) - 1
	if le < 2 {
		return v
	}

	if v[0] == '\'' && v[le] == '\'' {
		return v[1:le]
	}

	if v[0] == '"' && v[le] == '"' {
		return v[1:le]
	}

	return v
}
