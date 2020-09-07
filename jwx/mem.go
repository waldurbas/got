package jwx

// ----------------------------------------------------------------------------------
// mem.go for Go's jwx package (https://github.com/waldurbas/got)
// base implementation JWT token
// Copyright 2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.09.06 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"errors"
	"sync"
	"time"
)

// LoginUser #
type LoginUser struct {
	Pwd     string
	Expired int64
	App     string
	UUID    string
	PK      PairToken
	URL     string
}

//Key: user
var usersToken = make(map[string]*LoginUser, 100)

//Key: sub from token
var usersSub = make(map[string]*LoginUser, 100)

// mutex for locking
var mutex = &sync.Mutex{}

// CheckElapsedToken #
func CheckElapsedToken(withLock bool) {
	if withLock {
		mutex.Lock()
		defer mutex.Unlock()
	}

	cTime := time.Now().Unix()
	for k, v := range usersToken {
		if (v.Expired - cTime) < 1 {
			delete(usersSub, v.PK.Sub)
			delete(usersToken, k)
		}
	}
}

// TokenDeleteAll #
func TokenDeleteAll() {
	mutex.Lock()
	defer mutex.Unlock()

	usersToken = make(map[string]*LoginUser, 100)
	usersSub = make(map[string]*LoginUser, 100)
}

// TokenAdd #
func TokenAdd(kk string, lu *LoginUser) {
	mutex.Lock()
	defer mutex.Unlock()

	CheckElapsedToken(false)

	lu.Expired = time.Now().UTC().Add(DurationRefreshToken).Unix()

	usersToken[kk] = lu
	usersSub[lu.PK.Sub] = lu
}

// TokenGet #
func TokenGet(kk string) *LoginUser {
	mutex.Lock()
	defer mutex.Unlock()

	if lu, found := usersToken[kk]; found {
		return lu
	}

	return nil
}

// TokenGetSub #
func TokenGetSub(kk string) *LoginUser {
	mutex.Lock()
	defer mutex.Unlock()

	if lu, found := usersSub[kk]; found {
		return lu
	}

	return nil
}

// GetRefreshToken #
func GetRefreshToken(accessToken string) (string, error) {
	sub := RawToken2Sub(accessToken)
	if sub == "" {
		return "", errors.New("getRToken.RawSub")
	}

	lu := TokenGetSub(sub)
	if lu == nil {
		return "", errors.New("getRToken.GetSub")
	}

	return "JWT " + lu.PK.Refresh, nil
}
