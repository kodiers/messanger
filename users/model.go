package users

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"messanger/libs/infrastructure/configuration"
	"time"
	"unicode"
)

const PasswordLength = 8

var passwordSalt = configuration.Conf.GetSecretKey()

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	LastLogin    time.Time `json:"last_login"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"password_hash"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

func (u *User) IsValidPassword(password string) (bool, error) {
	if len(password) < PasswordLength {
		return false, errors.New(fmt.Sprintf("password length is less than, %v", PasswordLength))
	}
	hasUpper := false
	hasLower := false
	for _, c := range password {
		if unicode.IsUpper(c) {
			hasUpper = true
			break
		}
	}
	for _, c := range password {
		if unicode.IsLower(c) {
			hasLower = true
			break
		}
	}
	if !hasUpper || !hasLower {
		return false, errors.New("password should contains letters in upper or lowercase")
	}
	return true, nil
}

func (u *User) GetPasswordWithSalt(password string) string {
	return passwordSalt + password
}

func (u *User) MakePasswordHash(password string) []byte {
	passwordWithSalt := u.GetPasswordWithSalt(password)
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalln("Could not generate password hash.")
	}
	u.PasswordHash = string(hash)
	return hash
}

func (u *User) SetUserPasswordHash() {
	u.PasswordHash = string(u.MakePasswordHash(u.Password))
}

func (u *User) Save() {

}
