package mongodb

import (
	"MessengerDemo/tools"
	"encoding/json"
	"fmt"
	"log"
)

func CreateMessageForm() string {
	Message := Messages{
		Id:          tools.UuidGen(),
		SenderId:    "Sender's UUID",
		ReceiverId:  "Receiver or group's UUID",
		Contents:    "edata from server.go",
		Created:     "test",
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
