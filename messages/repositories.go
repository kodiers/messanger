package messages

import (
	"database/sql"
	"log"
	"messanger/libs/infrastructure/configuration"
	"messanger/libs/utils"
	"time"
)

type MessageRepository struct {
	DB *sql.DB
}

func InitRepository(db *sql.DB) MessageRepository {
	return MessageRepository{DB: db}
}

var MessageRep = InitRepository(configuration.DB)

func (mr *MessageRepository) makeMessage(row sql.Row) (Message, error) {
	var created, updated string
	var delivered sql.NullTime
	message := new(Message)
	err := row.Scan(&message.ID, &message.Text, &created, &updated, &delivered, &message.SenderId, &message.ReceiverId)
	if err != nil {
		log.Println("Could not parse rows from db.", err)
		return *message, err
	}
	location := utils.LoadLocation()
	message.Created, _ = time.ParseInLocation(time.RFC3339, created, &location)
	message.Updated, _ = time.ParseInLocation(time.RFC3339, updated, &location)
	if delivered.Valid {
		message.Delivered = delivered.Time
	}
	return *message, nil
}

func (mr *MessageRepository) GetMessagesBySenderReceiverId(senderId int, receiverId int) ([]Message, error) {
	rows, err := mr.DB.Query("SELECT * FROM MESSAGES WHERE SENDER=$1 AND RECEIVER=$2", senderId, receiverId)
	if err != nil {
		log.Println("Could not get messages from db", err)
		return nil, err
	}
	defer rows.Close()
	messages := make([]Message, 0)
	for rows.Next() {
		message := new(Message)
		var created, updated string
		var delivered sql.NullTime
		err := rows.Scan(&message.ID, &message.Text, &created, &updated, &delivered, &message.SenderId, &message.ReceiverId)
		if err != nil {
			log.Println("Could not read rows data", err)
		}
		location := utils.LoadLocation()
		if delivered.Valid {
			message.Delivered = delivered.Time
		}
		message.Created, _ = time.ParseInLocation(time.RFC3339, created, &location)
		message.Updated, _ = time.ParseInLocation(time.RFC3339, updated, &location)
		messages = append(messages, *message)
	}
	return messages, nil
}

func (mr *MessageRepository) GetLastMessage(text string, senderId int, receiverId int) (Message, error) {
	row := mr.DB.QueryRow("SELECT * FROM MESSAGES WHERE TEXT=$1 AND SENDER=$2 AND RECEIVER=$3 ORDER BY CREATED DESC LIMIT 1;", text, senderId, receiverId)
	message, err := mr.makeMessage(*row)
	if err != nil {
		log.Println("Could not parse rows from db.", err)
		return message, err
	}
	return message, nil
}

func (mr *MessageRepository) CreateMessage(text string, senderId int, receiverId int) (Message, error) {
	_, err := mr.DB.Exec("INSERT INTO MESSAGES (TEXT, SENDER, RECEIVER) VALUES ($1, $2, $3);", text, senderId,
		receiverId)
	var message Message
	if err != nil {
		log.Println("Could not create message record ", err)
		return message, err
	}
	message, err = mr.GetLastMessage(text, senderId, receiverId)
	if err != nil {
		log.Println("Could not get message record ", err)
		return message, err
	}
	return message, nil
}
