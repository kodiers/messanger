package users

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	LastLogin    time.Time `json:"last_login"`
	Password     string    `json:"password"`
	PasswordHash string    `json:"password_hash"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}
