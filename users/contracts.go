package users

type UserRegistration struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
