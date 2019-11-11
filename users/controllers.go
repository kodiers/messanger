package users

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"messanger/libs/http/token"
	"messanger/libs/infrastructure/configuration"
	"net/http"
)

var userRepository = InitUserRepository(configuration.DB)

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if token.IsValidToken(r.Header) {
		http.Error(w, "You should not be authenticated.", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("passwordConfirmation")
	if username == "" || password == "" || passwordConfirmation == "" {
		http.Error(w, "Username, password or passwordConfirmation is empty.", http.StatusBadRequest)
		return
	}
	_, err := userRepository.getUserByUsername(username)
	if err == nil {
		http.Error(w, "Username with this username already exists.", http.StatusBadRequest)
		return
	}
	user := User{
		Username: username,
		Password: password,
	}
	isValidPassword, _ := user.IsValidPassword(password)
	if !isValidPassword {
		http.Error(w,
			fmt.Sprintf("Password should be %v and contains characters in upper and lower case.", PasswordLength),
			http.StatusBadRequest)
		return
	}
	user.SetUserPasswordHash()
	_, err = userRepository.InsertUser(&user)
	if err != nil {
		http.Error(w, "Cannot create new user", http.StatusInternalServerError)
		return
	}
	response := Response{
		Status: "User created successfully",
		Data:   user,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Cannot parse data to JSON", http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprintf(w, string(responseData))
	if err != nil {
		http.Error(w, "Cannot send data", http.StatusInternalServerError)
		return
	}

}
