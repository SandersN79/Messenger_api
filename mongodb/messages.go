package mongodb

import (
	"MessengerDemo/tools"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Messages struct {
	Id            string
	SenderId      string
	ReceiverId    string
	Contents      string //[]byte
	Created       time.Time
}

func CreateMessageForm() string {
	Message := Messages{
		Id:          tools.UuidGen(),
		SenderId:    "Sender's UUID",
		ReceiverId:  "Receiver or group's UUID",
		Contents:    "edata from server.go",
		Created:     time.Now(),
	}
	var jsonData []byte
	jsonData, err := json.Marshal(Message)
	if err != nil {
		log.Println(err)
	}
	MessageForm := string(jsonData)
	fmt.Println(MessageForm)
	return MessageForm
}