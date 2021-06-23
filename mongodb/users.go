package mongodb

import (
	"MessengerDemo/tools"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Users struct {
	Id          string
	UUserKey    string
	Username    string
	//GroupIds    []string
	//GroupRoles  []string
	Email       string
	Password    string
	Created     time.Time
}


func RegisterUserProfile() string {
	User := Users{
		Id:       tools.UuidGen(),
		UUserKey: tools.KeyGen(),
		Username: "Admin",
		//GroupIds:    "test",
		//GroupRoles:  "test",
		Email:         "Admin@test",
		Password:      "null",
		Created:       time.Now(),
	}
	var jsonData []byte
	jsonData, err := json.Marshal(User)
	if err != nil {
		log.Println(err)
	}
	UserAccount := string(jsonData)
	fmt.Println(UserAccount)
	return UserAccount
}
