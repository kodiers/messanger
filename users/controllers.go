package users

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"messanger/libs/http/token"
	"net/http"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if IsAuthenticated(r.Header) {
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
	_, err = UserRep.GetUserByUsername(userRegistration.Username)
	if err == nil {
		http.Error(w, "Username with this username already exists.", http.StatusBadRequest)
		return
	}
	user := User{
		Username:  userRegistration.Username,
		Password:  userRegistration.Password,
		LastLogin: time.Now(),
	}
	isValidPassword, _ := user.IsValidPassword(userRegistration.Password)
	if !isValidPassword {
		http.Error(w,
			fmt.Sprintf("Password should be %v and contains characters in upper and lower case.", PasswordLength),
			http.StatusBadRequest)
		return
	}
	user.SetUserPasswordHash()
	_, err = UserRep.InsertUser(&user)
	if err != nil {
		http.Error(w, "Cannot create new user", http.StatusInternalServerError)
		return
	}
	response := Response{
		Status: "User created successfully",
		Data:   CreateUserResponseWithoutPassHash(&user),
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
	if IsAuthenticated(r.Header) {
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
	user, err := UserRep.GetUserByUsername(userLogin.Username)
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
	go UserRep.UpdateUserLastLogin(user)
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

func UsersList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	users, err := UserRep.GetUsers()
	if err != nil {
		http.Error(w, "Unknown error.", http.StatusInternalServerError)
		return
	}
	if len(users) == 0 {
		http.Error(w, "Not users found.", http.StatusNotFound)
		return
	}
	var usersResponse []UserResponseWithoutPasswordHash
	for _, user := range users {
		userResponse := CreateUserResponseWithoutPassHash(&user)
		usersResponse = append(usersResponse, userResponse)
	}
	response := Response{
		Status: "Ok",
		Data:   usersResponse,
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
