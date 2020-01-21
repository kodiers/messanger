package messages

import "messanger/users"

type CreateMessageContract struct {
	Text     string `json:"text"`
	Receiver int    `json:"receiverId"`
}

func (cmc *CreateMessageContract) TextIsCorrect() bool {
	return cmc.Text != ""
}

func (cmc *CreateMessageContract) ReceiverIsCorrect() bool {
	u, err := users.UserRep.GetUserById(cmc.Receiver)
	return err == nil && u != users.User{}
}
