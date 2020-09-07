package jwx

import (
	"errors"
	"time"

	"github.com/waldurbas/got/cnv"
)

// ----------------------------------------------------------------------------------
// jwx.go (https://github.com/waldurbas/got): base implementation JWT token
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

// LoginUser #
type LoginUser struct {
	Pwd    string
	Ablauf int64
	App    string
	Kid    string
	UUID   string
	PK     PairToken
	URL    string
}

var usersToken = make(map[string]*LoginUser, 100)
var usersSub = make(map[string]*LoginUser, 100)

func getRefreshToken(accessToken string) (string, error) {
	sub := RawToken2Sub(accessToken)
	if sub == "" {
		return "", errors.New("getRToken.RawSub")
	}

	lu := TokenGetSub(sub)
	if lu == nil {
		return "", errors.New("getRToken.GetSub")
	}

	return lu.PK.Refresh, nil
}

// MkLoginUserKey #
func MkLoginUserKey(app string, kid string, uid string) string {
	return app + kid + cnv.StripUUID36(uid)
}

// CheckElapsedToken #
func CheckElapsedToken() {
	cTime := time.Now().Unix()
	for k, v := range usersToken {
		if (v.Ablauf - cTime) < 1 {
			TokenDelete(k)
		}
	}
}

// TokenDelete #
func TokenDelete(key string) {
	if lu, found := usersToken[key]; found {
		delete(usersSub, lu.PK.Sub)
		delete(usersToken, key)
	}
}

// TokenAdd #
func TokenAdd(kk string, lu *LoginUser) {
	CheckElapsedToken()

	lu.Ablauf = time.Now().UTC().Add(DurationRefreshToken).Unix()

	usersToken[kk] = lu
	usersSub[lu.PK.Sub] = lu
}

// TokenGet #
func TokenGet(kk string) *LoginUser {
	if lu, found := usersToken[kk]; found {
		return lu
	}

	return nil
}

// TokenGetSub #
func TokenGetSub(kk string) *LoginUser {
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

	return lu.PK.Refresh, nil
}
