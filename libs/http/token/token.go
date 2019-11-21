package token

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"messanger/libs/infrastructure/configuration"
	"net/http"
	"time"
)

const hmacSecret = "eyJhbGciOiAiUlMyNTYiLCAia2lkIjogInNvbWUta2V5LW5hbWUifQ==.eyJob3N0IjogImlkLnJhbWJsZXIucnUiLCAidXJpIjogIi9ycGMiLCAibWV0aG9kIjogInBvc3QiLCAiaGVhZGVycyI6ICIiLCAiY3RpbWUiOiAxNTcxMDU3ODIyLjQ1M"

var expiration = configuration.Conf.Session.Expiration

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

func MakeJWT(userId int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     userId,
		"expiration": int64(time.Now().Second()) + expiration,
	})
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Println("could not create token.")
	}
	return tokenString
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
	if !userOk || !expireOk {
		return false
	}
	return true
}
