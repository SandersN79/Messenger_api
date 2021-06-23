package mongodb

import (
	"MessengerDemo/tools"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Groups struct {
	Id          string
	Name        string
	Type        string
	Created     time.Time
}

func CreateGroup() string {
	Group := Groups{
		Id:      tools.UuidGen(),
		Name:    "test",
		Type:    "test",
		Created: time.Now(),
	}
	var jsonData []byte
	jsonData, err := json.Marshal(Group)
	if err != nil {
		log.Println(err)
	}
	GroupData := string(jsonData)
	fmt.Println(GroupData)
	return GroupData
}
