package jwx

// ----------------------------------------------------------------------------------
// jwx.go (https://github.com/waldurbas/got): implementation JWT Access/Refresh-Token
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
	"strconv"
	"sync"
	"time"

	"github.com/waldurbas/got/jwt"
)

var (
	// DurationAccessToken #
	DurationAccessToken = 15 * time.Minute
	// DurationRefreshToken #
	DurationRefreshToken = 16 * time.Hour

	//Key: user
	usersToken = make(map[string]*LoginUser, 100)

	//Key: sub from token
	usersSub = make(map[string]*LoginUser, 100)

	// subID
	subID = time.Now().UTC().Unix()

	// mutex for locking
	mutex = &sync.Mutex{}
)

// LoginUser #
type LoginUser struct {
	Pwd     string
	Expired int64
	App     string
	UUID    string
	URL     string

	Sub     string
	SignKey string
	Access  *jwt.JWToken
	Refresh *jwt.JWToken
}

// MemToken #
type MemToken struct {
	App     string
	Expired int64
	UUID    string
	Refresh string
}

func incSubID() int64 {
	mutex.Lock()
	defer mutex.Unlock()

	subID++
	return subID
}

// GenerateTokenPair #
func GenerateTokenPair(m *map[string]interface{}, key string, signKey string) (*LoginUser, error) {
	sub := strconv.FormatInt(incSubID(), 10)
	sub = sub[5:] + "." + sub[:5]

	at := jwt.New()
	at.Claims["sub"] = sub

	exp := time.Now().Add(DurationRefreshToken).UTC().Unix()

	rt := jwt.New()
	rt.Claims.Assign(m)
	rt.Claims["sub"] = sub
	rt.Claims["exp"] = exp

	at.Encode(signKey)
	rt.Encode(signKey)

	lu := &LoginUser{
		Expired: exp,
		Access:  at,
		Refresh: rt,
		Sub:     sub,
		SignKey: signKey,
	}

	mutex.Lock()
	defer mutex.Unlock()

	usersToken[key] = lu
	usersSub[sub] = lu

	return lu, nil
}

// Validate #
func Validate(headToken string, signKey string) (map[string]interface{}, error) {
	tkn := jwt.New()
	if err := tkn.Parse(headToken); err != nil {
		return nil, err
	}

	if err := tkn.Valid(signKey); err != nil {
		return nil, err
	}

	m := make(map[string]interface{}, len(tkn.Claims))

	for k, v := range tkn.Claims {
		m[k] = v
	}

	return m, nil
}

// RawToken2Sub #
func RawToken2Sub(headToken string) string {
	tkn := jwt.New()
	if tkn.Parse(headToken) != nil {
		return ""
	}

	return tkn.Claims["sub"].(string)
}

// CheckElapsedToken #
func CheckElapsedToken(withLock bool) {
	if withLock {
		mutex.Lock()
		defer mutex.Unlock()
	}

	cTime := time.Now().Unix()
	for k, v := range usersToken {
		if (v.Expired - cTime) < 1 {
			delete(usersSub, v.Sub)
			delete(usersToken, k)
		}
	}
}

// TokenDelete #
func TokenDelete(uid string) {
	mutex.Lock()
	defer mutex.Unlock()

	if uid == "all" {
		usersToken = make(map[string]*LoginUser, 100)
		usersSub = make(map[string]*LoginUser, 100)
	} else {
		for k, v := range usersToken {
			if v.UUID == uid {
				delete(usersSub, v.Sub)
				delete(usersToken, k)
			}
		}
	}
}

// CheckUserToken #
func CheckUserToken(kk string) *LoginUser {
	mutex.Lock()
	defer mutex.Unlock()

	CheckElapsedToken(false)

	lu, found := usersToken[kk]

	if !found {
		return nil
	}

	// check Access-Token 15Min.
	if lu.Access.Expired() > int64(DurationAccessToken) {
		if lu.Refresh.Expired() > int64(DurationRefreshToken) {
			return nil
		}

		tkn := jwt.New()
		tkn.Claims["sub"] = lu.Sub
		tkn.Encode(lu.SignKey)
		lu.Access = tkn
	}

	return &LoginUser{
		Pwd:     lu.Pwd,
		Expired: lu.Expired,
		App:     lu.App,
		UUID:    lu.UUID,
		URL:     lu.URL,
		Sub:     lu.Sub,
		SignKey: lu.SignKey,
		Access:  lu.Access,
		Refresh: lu.Refresh,
	}
}

// GetMemToken #
func GetMemToken(app string) []MemToken {
	mutex.Lock()
	defer mutex.Unlock()

	CheckElapsedToken(false)

	ma := []MemToken{}

	for _, v := range usersToken {
		if v.App == app || app == "all" {
			m := MemToken{
				App:     v.App,
				Expired: v.Expired,
				UUID:    v.UUID,
				Refresh: v.Refresh.RawData,
			}
			ma = append(ma, m)
		}
	}

	return ma
}

func subToLoginUser(kk string) *LoginUser {
	mutex.Lock()
	defer mutex.Unlock()

	if lu, found := usersSub[kk]; found {

		return &LoginUser{
			Pwd:     lu.Pwd,
			Expired: lu.Expired,
			App:     lu.App,
			UUID:    lu.UUID,
			URL:     lu.URL,
			Sub:     lu.Sub,
			SignKey: lu.SignKey,
			Access:  lu.Access,
			Refresh: lu.Refresh,
		}
	}

	return nil
}

// GetRefreshToken #
func GetRefreshToken(accessToken string) (*jwt.JWToken, error) {
	sub := RawToken2Sub(accessToken)
	if sub == "" {
		return nil, errors.New("getRToken.Raw")
	}

	lu := subToLoginUser(sub)
	if lu == nil {
		return nil, errors.New("getRToken.Sub")
	}

	if time.Now().UTC().Unix() > lu.Expired {
		return nil, jwt.ErrTokenExpired
	}

	return lu.Refresh, nil
}
