package messages

import (
	"time"
)

type Message struct {
	ID         int       `json:"id"`
	Text       string    `json:"text"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	Delivered  time.Time `json:"delivered"`
	SenderId   int       `json:"senderId"`
	ReceiverId int       `json:"receiverId"`
}
