package messages

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"messanger/libs/http/contracts"
	"messanger/users"
	"net/http"
)

func CreateMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var createMessageCtr CreateMessageContract
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading data from request", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(bs, &createMessageCtr)
	if err != nil {
		http.Error(w, "Error parsing data from request.", http.StatusInternalServerError)
		return
	}
	if !createMessageCtr.TextIsCorrect() || !createMessageCtr.ReceiverIsCorrect() {
		http.Error(w, "Text or Receiver id is incorrect.", http.StatusBadRequest)
		return
	}
	user, _ := users.GetUserFromHeader(r.Header)
	message, err := MessageRep.CreateMessage(createMessageCtr.Text, user.ID, createMessageCtr.Receiver)
	if err != nil {
		http.Error(w, "Cannot create message.", http.StatusInternalServerError)
		return
	}
	response := contracts.Response{
		Status: "OK",
		Data:   message,
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
