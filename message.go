package main

import (
	"bytes"
	"errors"
	"fmt"
)

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

