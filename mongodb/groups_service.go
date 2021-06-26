package mongodb

import (
	"MessengerDemo/tools"
	"encoding/json"
	"fmt"
	"log"
)

func CreateGroup() string {
	Group := Groups{
		Id:       tools.UuidGen(),
		Name:     "test",
		Type:     "test",
		Created:  "test",
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

