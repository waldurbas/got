package jwx

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

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// DurationAccessToken #
	DurationAccessToken = 15 * time.Minute
	// DurationRefreshToken #
	DurationRefreshToken = 12 * time.Hour
	// ErrTokenExpired #
	ErrTokenExpired = "Token is expired"
)

// PairToken #
type PairToken struct {
	Sub     string
	SignKey string
	Access  string
	Refresh string
}

func generateToken(m map[string]interface{}, signKey string, d time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().UTC().Add(d).Unix()

	for k, v := range m {
		claims[k] = v
	}

	s, err := token.SignedString([]byte(signKey))

	if err != nil {
		return "", err
	}

	return s, nil
}

// ReNewAccessToken #
func ReNewAccessToken(pk *PairToken) error {
	_, err := Token2Map("JWT "+pk.Refresh, pk.SignKey)
	if err != nil {
		if err.Error() == ErrTokenExpired {
			return err
		}
	}

	at, err := generateToken(map[string]interface{}{"sub": pk.Sub}, pk.SignKey, DurationAccessToken)
	if err != nil {
		return err
	}

	pk.Access = at

	return nil
}

// CheckToken #
func CheckToken(pk *PairToken) error {
	_, err := Token2Map("Bearer "+pk.Access, pk.SignKey)
	if err != nil {
		if err.Error() == ErrTokenExpired {
			err = ReNewAccessToken(pk)
			if err != nil {
				return err
			}
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// GenerateTokenPair #
func GenerateTokenPair(m map[string]interface{}, sub string, signKey string) (*PairToken, error) {
	at, err := generateToken(map[string]interface{}{"sub": sub}, signKey, DurationAccessToken)
	if err != nil {
		return nil, err
	}

	m["sub"] = sub
	rt, err := generateToken(m, signKey, DurationRefreshToken)
	if err != nil {
		return nil, err
	}

	return &PairToken{Access: at, Refresh: rt, Sub: sub, SignKey: signKey}, nil
}

// Token2Map #
func Token2Map(headToken string, signKey string) (map[string]interface{}, error) {
	ss := strings.Split(headToken, " ")

	if len(ss) == 2 && (ss[0] == "JWT" || ss[0] == "Bearer") {
		claims := jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(ss[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signKey), nil
		})

		if err != nil {
			return nil, err
		}

		if !tkn.Valid {
			return nil, errors.New("token not valid")
		}

		m := make(map[string]interface{}, len(claims))

		for k, v := range claims {
			m[k] = v
		}

		return m, nil
	}

	return nil, errors.New("bad Authorization")
}

// RawToken2Map #
func RawToken2Map(headToken string) (map[string]interface{}, error) {
	ss := strings.Split(headToken, " ")

	if len(ss) == 2 && (ss[0] == "JWT" || ss[0] == "Bearer") {
		sx := strings.Split(ss[1], ".")
		if len(sx) == 3 {
			ds, err := base64.RawStdEncoding.DecodeString(sx[1])
			if err != nil {
				return nil, err
			}

			m := make(map[string]interface{}, 10)
			if err := json.Unmarshal(ds, &m); err != nil {
				return nil, err
			}

			return m, nil
		}
	}
	return nil, errors.New("bad Authorization")
}

// RawToken2Sub #
func RawToken2Sub(headToken string) string {
	m, err := RawToken2Map(headToken)
	if err == nil {
		return m["sub"].(string)
	}
	return ""
}
