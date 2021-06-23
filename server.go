package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/JECSand/fetch"
	"io/ioutil"
	"log"
	"strings"
)

type Encryption struct {
	UserKey     string      `json:"user_key,omitempty"`
	Message     []byte	    `json:"message,omitempty"`
}

// EncryptionService is used by the app to manage all repos related controllers and functionality
type EncryptionService struct {
	host   string
	port   string
}

// NewEncryptionService ..
func NewEncryptionService(host string, port string) *EncryptionService {
	return &EncryptionService{
		host:   host,
		port:   port,
	}
}

// buildURL
func (enc *EncryptionService) buildURL(urlType string) string {
	var host string
	host = enc.host + ":" + enc.port
	if urlType == "encrypt" {
		host = host + "/encrypt"
	} else if urlType == "decrypt" {
		host = host + "/decrypt"
	}
	if !strings.Contains(host, "http://") {
		host = "http://" + host
	}
	return host
}

// Evaluate sends the JIT to the Encryption API and generates an App's source code.
func (enc *EncryptionService) Encrypt() ([]byte, error) {
	decURL := enc.buildURL("encrypt")
	fmt.Println(decURL)
	method := "POST"
	headers := fetch.JSONDefaultHeaders()
	cryptMessage := NewCryptMessage("1603 4702 613", []byte("this is a test abcdefghijklmnopqrstuvwxyz!@#$%^&*()_+ABCDEFGHIJKLMNOPQRSTUVWXYZ<>?{}|,./;[]1234567890-="))
	output := cryptMessage.ToJSON()
	f, err := fetch.NewFetch(decURL, method, headers, bytes.NewBuffer([]byte(output)))
	if err != nil {
		fmt.Println("Encryption Service Failed: ", err.Error())
		return []byte(""), err
	}
	f.Execute("")
	f.Resolve()
	if f.Res == nil {
		fmt.Println("error")
		return []byte(""), nil
	} else {
		fmt.Println(f.Res.Status, f.Res.Body)
		defer f.Res.Body.Close()
		responseData, err := ioutil.ReadAll(f.Res.Body)
		if err != nil {
			fmt.Println("Failed to read Encryption Response: ", err.Error())
			log.Fatalf("Failed to read Encryption Response: %v", err.Error())
			return []byte(""), err
		}
		str := base64.StdEncoding.EncodeToString(responseData)
		//fmt.Println("str:", str)

		bdata, err := base64.StdEncoding.DecodeString(str)
		ddata := string(bdata)
		//rdata := []rune(ddata)
		//fmt.Println((ddata)[1])
		fmt.Println("data:", ddata)
		return bdata, nil
	}
	return []byte(""), nil
}

// Evaluate sends the JIT to the Encryption API and generates an App's source code.
func (enc *EncryptionService) Decrypt(edata []byte) error {
	fmt.Println("decryption section \n")
	decURL := enc.buildURL("decrypt")
	method := "POST"
	headers := fetch.JSONDefaultHeaders()
	fmt.Println("edata:", string(edata))
	cryptMessage := NewCryptMessage("1603 4702 613", edata)
	output := cryptMessage.ToJSON()
	fmt.Println("output:", output)
	f, err := fetch.NewFetch(decURL, method, headers, bytes.NewBuffer([]byte(output)))
	if err != nil {
		fmt.Println("decryption Service Failed: ", err.Error())
		return err
	}
	f.Execute("")
	f.Resolve()
	if f.Res == nil {
		fmt.Println("error")
		return nil
	} else {
		//fmt.Println(f.Res.Status, f.Res.Body)
		defer f.Res.Body.Close()
		responseData, err := ioutil.ReadAll(f.Res.Body)
		if err != nil {
			fmt.Println("Failed to read decryption Response: ", err.Error())
			log.Fatalf("Failed to read decryption Response: %v", err.Error())
			return err
		}
		str := base64.StdEncoding.EncodeToString(responseData)
		//fmt.Println("str:", str)

		bdata, err := base64.StdEncoding.DecodeString(str)
		ddata := string(bdata)
		fmt.Println("data:", ddata)
		return nil
	}
	return nil
}



