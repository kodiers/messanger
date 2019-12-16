package users

import (
	"fmt"
	"messanger/libs/http/token"
	"net/http"
	"strings"
)

func IsAuthenticated(header http.Header) bool {
	_, ok := header["Authentication"]
	if !ok {
		return false
	}
	_, err := GetUserFromHeader(header)
	if err != nil {
		return false
	}
	return true
}

func GetUserFromHeader(header http.Header) (User, error) {
	var user = User{}
	if token.IsValidToken(header) {
		claims, ok := token.GetClaims(strings.Join(header["Authentication"], ""))
		if !ok {
			return user, fmt.Errorf("cannot get claims from token")
		}
		user, err := UserRep.GetUserById(claims.UserId)
		if err != nil {
			return user, fmt.Errorf("user not found")
		}
		return user, nil
	}
	return user, fmt.Errorf("token invalid")
}
