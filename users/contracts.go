package users

import "time"

type UserRegistration struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type UserResponseWithoutPasswordHash struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	LastLogin string `json:"last_login"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
}

func CreateUserResponseWithoutPassHash(u *User) UserResponseWithoutPasswordHash {
	response := UserResponseWithoutPasswordHash{
		ID:        u.ID,
		Username:  u.Username,
		LastLogin: u.LastLogin.Format(time.RFC3339),
		Created:   u.Created.Format(time.RFC3339),
		Updated:   u.Updated.Format(time.RFC3339),
	}
	return response
}
