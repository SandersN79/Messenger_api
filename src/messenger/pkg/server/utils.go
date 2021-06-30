package server

import (
	"bytes"

	"io"
	"mime/multipart"
	"net/http"
)

type jsonErr struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

// Function to read file contents as bytes
func ReadFile(file multipart.File) []byte {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// HandleOptionsRequest handles incoming OPTIONS request
func HandleOptionsRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Auth-Token, API-Key")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Auth-Token, API-Key")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,DELETE,POST,PATCH")
	w.WriteHeader(http.StatusOK)
}

// SetResponseHeaders sets the response headers being sent back to the client
func SetResponseHeaders(w http.ResponseWriter, authToken string, apiKey string) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Auth-Token, API-Key")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Auth-Token, API-Key")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET,DELETE,POST,PATCH")
	if authToken != "" {
		w.Header().Add("Auth-Token", authToken)
	}
	if apiKey != "" {
		w.Header().Add("API-Key", apiKey)
	}
	return w
}

// AdminRouteRoleCheck checks admin routes JWT tokens to ensure that a group admin does not break scope
func AdminRouteRoleCheck(decodedToken []string) []string {
	var groupIds []string
	if decodedToken[2] != "Root" {
		groupIds = append(groupIds, decodedToken[1])
	}
	return groupIds
}