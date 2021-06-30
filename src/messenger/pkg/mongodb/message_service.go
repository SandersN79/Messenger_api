package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
	"encoding/json"
	"fmt"
	"log"
)

func CreateMessageForm() string {
	Message := messages{
		Id:          core.UuidGen(),
		SenderId:    "Sender's UUID",
		ReceiverId:  "Receiver or group's UUID",
		Contents:    []byte("edata from server.go"),
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
