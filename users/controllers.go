package users

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"messanger/libs/http/token"
	"messanger/libs/infrastructure/configuration"
	"net/http"
)

var userRepository = InitUserRepository(configuration.DB)

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if token.IsValidToken(r.Header) {
		http.Error(w, "You should not be authenticated.", http.StatusBadRequest)
		return
	}
	var userRegistration UserRegistration
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading data from request.", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bs, &userRegistration)
	if err != nil {
		http.Error(w, "Error pasring data from request.", http.StatusInternalServerError)
		return
	}
	if userRegistration.Username == "" || userRegistration.Password == "" || userRegistration.PasswordConfirmation == "" {
		http.Error(w, "Username, password or passwordConfirmation is empty.", http.StatusBadRequest)
		return
	}
	_, err = userRepository.getUserByUsername(userRegistration.Username)
	if err == nil {
		http.Error(w, "Username with this username already exists.", http.StatusBadRequest)
		return
	}
	user := User{
		Username: userRegistration.Username,
		Password: userRegistration.Password,
	}
	isValidPassword, _ := user.IsValidPassword(userRegistration.Password)
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

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if token.IsValidToken(r.Header) {
		http.Error(w, "You should not be authenticated.", http.StatusBadRequest)
		return
	}
	var userLogin UserLogin
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading data from request.", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bs, &userLogin)
	if err != nil {
		http.Error(w, "Error parsing data from request.", http.StatusInternalServerError)
		return
	}
	if userLogin.Username == "" || userLogin.Password == "" {
		http.Error(w, "Username or password is empty.", http.StatusBadRequest)
		return
	}
	user, err := userRepository.getUserByUsername(userLogin.Username)
	if err != nil {
		http.Error(w, "Username with this username could not be found.", http.StatusBadRequest)
		return
	}
	isValidPassword, _ := user.IsValidPassword(userLogin.Password)
	if !isValidPassword {
		http.Error(w,
			fmt.Sprintf("Password should be %v and contains characters in upper and lower case.", PasswordLength),
			http.StatusBadRequest)
		return
	}
	sentPasswordHash := user.MakePasswordHash(userLogin.Password)
	if string(sentPasswordHash) != user.PasswordHash {
		http.Error(w, "Password is incorrect.", http.StatusForbidden)
		return
	}
	createdToken := token.MakeJWT(user.ID)
	if createdToken == "" {
		http.Error(w, "Error while signing in.", http.StatusInternalServerError)
		return
	}
	response := Response{
		Status: "Logged in.",
		Data:   createdToken,
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
