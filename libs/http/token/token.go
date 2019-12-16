package token

import (
	"github.com/dgrijalva/jwt-go"
	"log"
	"messanger/libs/infrastructure/configuration"
	"net/http"
	"time"
)

const hmacSecret = "eyJhbGciOiAiUlMyNTYiLCAia2lkIjogInNvbWUta2V5LW5hbWUifQ==.eyJob3N0IjogImlkLnJhbWJsZXIucnUiLCAidXJpIjogIi9ycGMiLCAibWV0aG9kIjogInBvc3QiLCAiaGVhZGVycyI6ICIiLCAiY3RpbWUiOiAxNTcxMDU3ODIyLjQ1M"

var expiration = configuration.Conf.Session.Expiration

type MyCustomClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

func GetClaims(tokenStr string) (MyCustomClaims, bool) {
	token, _ := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(hmacSecret), nil
	})
	if token.Valid {
		claims, _ := token.Claims.(*MyCustomClaims)
		return *claims, true
	}
	return MyCustomClaims{}, false
}

func MakeJWT(userId int) string {
	claims := &MyCustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + expiration,
			Issuer:    "messenger-app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		log.Println("could not create token.")
	}
	return tokenString
}

func IsValidToken(header http.Header) bool {
	token, ok := header["Authentication"]
	if !ok || token[0] == "" {
		return false
	}
	_, ok = GetClaims(token[0])
	if !ok {
		return false
	}
	return true
}
