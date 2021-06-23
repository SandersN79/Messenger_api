package main

import (
	"MessengerDemo/mongodb"
	"fmt"
)

func main() {

	//encryptionService := NewEncryptionService("127.0.0.1", "9999")
	//edata, _ := encryptionService.Encrypt()
	//encryptionService.Decrypt(edata)
	//KeyGen()
	RegisterUserProfile()
	mongodb.CreateGroup()
	fmt.Println("Finished")
}
