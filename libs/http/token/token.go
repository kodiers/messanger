package token

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"messanger/libs/infrastructure/configuration"
	"net/http"
	"time"
)

const hmacSampleSecret = "eyJhbGciOiAiUlMyNTYiLCAia2lkIjogInNvbWUta2V5LW5hbWUifQ==.eyJob3N0IjogImlkLnJhbWJsZXIucnUiLCAidXJpIjogIi9ycGMiLCAibWV0aG9kIjogInBvc3QiLCAiaGVhZGVycyI6ICIiLCAiY3RpbWUiOiAxNTcxMDU3ODIyLjQ1M"

type Token struct {
	UserID  int
	Token   string
	Expired time.Time
}

func isValidJwt(tokenStr string) jwt.MapClaims {
	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return true, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims
	}
	return nil
}

var tokenRepository = InitTokenRepository(configuration.DB, configuration.Conf.Session.Expiration)

func makeToken(userdId int, token string, expiration time.Time) Token {
	sess := Token{
		UserID:  userdId,
		Token:   token,
		Expired: expiration,
	}
	return sess
}

func (t Token) isExpired() bool {
	return t.Expired.Second() < time.Now().Second()
}

func (t Token) isValidUser(db *sql.DB) bool {
	return false
}

func IsValidToken(header http.Header) bool {
	token, ok := header["Authentication"]
	if !ok {
		return false
	}
	tokenValues := isValidJwt(token[0])
	if tokenValues == nil {
		return false
	}

	_, userOk := tokenValues["userId"].(int)
	_, expireOk := tokenValues["expiration"].(time.Time)
	tokenId, tokenOk := tokenValues["tokenId"].(string)
	if !userOk || !expireOk || !tokenOk {
		return false
	}
	tokenObj, err := tokenRepository.GetTokenFromDb(tokenId)
	if err != nil {
		return false
	}
	if tokenObj.isExpired() {
		return false
	}
	return true
}
