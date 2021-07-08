package internals

import (
	"bytes"
	"encoding/base64"
	"errors"
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
func (enc *EncryptionService) Encrypt(contents []byte) ([]byte, error) {
	decURL := enc.buildURL("encrypt")
	fmt.Println(decURL)
	method := "POST"
	headers := fetch.JSONDefaultHeaders()
	cryptMessage := NewCryptMessage("1603 4702 613", contents)
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
		bdata = cleanJSON(bdata)
		return bdata, nil
	}
	return []byte(""), nil
}

// Evaluate sends the JIT to the Encryption API and generates an App's source code.
func (enc *EncryptionService) Decrypt(edata []byte) ([]byte, error) {
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
		return []byte{}, err
	}
	f.Execute("")
	f.Resolve()
	if f.Res == nil {
		fmt.Println("error")
		return []byte{}, nil
	} else {
		//fmt.Println(f.Res.Status, f.Res.Body)
		defer f.Res.Body.Close()
		responseData, err := ioutil.ReadAll(f.Res.Body)
		if err != nil {
			fmt.Println("Failed to read decryption Response: ", err.Error())
			log.Fatalf("Failed to read decryption Response: %v", err.Error())
			return []byte{}, err
		}
		str := base64.StdEncoding.EncodeToString(responseData)
		//fmt.Println("str:", str)

		bdata, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return []byte{}, err
		}
		//ddata := string(bdata)
		//fmt.Println("data:", ddata)
		return bdata, nil
	}
	return []byte{}, nil
}


type CryptMessage struct {
	UserKey    string  `json:"user_key,omitempty"`
	Message    []byte  `json:"message,omitempty"`
	//Id         int64  `json:"ref"`
	//Created    time.Time
}
func cleanJSON(eData []byte) []byte {
	var reData []byte
	build := false
	lastChar := ""
	count := 0
	for i := len(eData)-1; i >= 0; i-- {
		//fmt.Println(string(eData[i]))
		if build {
			reData = append([]byte{eData[i]}, reData...) // "}
		}
		if string(eData[i]) == `"` && lastChar == "}" {
			if count < 5 {
				build = true
			}
		}
		lastChar = string(eData[i])
		count = count + 1
	}
	return reData
}

func loadMessageJSON(edata []byte) ([]byte, error) {
	byteSlice := bytes.Split(edata, []byte(`{"message":"`))
	if len(byteSlice) == 0 {
		return []byte(""), errors.New("error on line 60: byteSlice is empty")
	}
	byteMsg := byteSlice[1]
	fmt.Println("check1 ", string(byteMsg))
	//byteMsg = bytes.Split(byteMsg, []byte(`"}`))[0]
	byteMsg = cleanJSON(byteMsg)
	fmt.Println("\ncheck2 ", string(byteMsg))
	return byteMsg, nil
}

func NewCryptMessage(userkey string, message []byte) *CryptMessage {
	var err error
	if bytes.Contains(message, []byte(`{"message":"`)) {
		message, err = loadMessageJSON(message)
		if err != nil {
			panic(err)
		}
	}
	return &(CryptMessage{userkey, message})
}

func (cm *CryptMessage) ToJSON() string {
	jsonStr := `{"user_key":"` + cm.UserKey + `","message":"` + string(cm.Message) + `"}`
	return jsonStr
}

