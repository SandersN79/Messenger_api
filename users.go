package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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

func UuidGen() string {
	Uuid := uuid.New().String()
	return Uuid
}

func RegisterUserProfile() string {
	User := Users {
		Id:            UuidGen(),
		UUserKey:      KeyGen(),
		Username:      "Admin",
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
