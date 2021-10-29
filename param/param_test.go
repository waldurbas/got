package param_test

// ----------------------------------------------------------------------------------
// param_test.go (https://github.com/waldurbas/got
// Copyright 2020,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2021.02.11 (wu) add Test_StrXcmp
//-----------------------------------------------------------------------------------

import (
	"fmt"
	"testing"

	"github.com/waldurbas/got/param"
)

func Test_Init(t *testing.T) {
	fmt.Println("Test_Initp")
	e := []string{"prg", "/usr/wald/x.pas", "/235/328", "-pic=20", "/d", "/all=/abc/", "-nopic"}
	param.InitParams(e)

	param.Print()

	param.ValueCheck("nopic", "1")
	param.Print()

}
