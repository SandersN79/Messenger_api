package server

import (
	core "MessengerDemo/src/messenger/pkg"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"os"
)

// JWTError is a struct that is used to contain a json encoded error message for any JWT related errors
type JWTError struct {
	Message string `json:"message"`
}

// Return JSON Error to Requested is Auth is bad
func respondWithError(w http.ResponseWriter, status int, error JWTError) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Auth-Token")
	w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Auth-Token")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(error); err != nil {
		panic(err)
	}
}

// returnInitialCheckErrMsg
func returnInitialCheckErrMsg(userErr string, groupErr string, s1 string, s2 string) string {
	errMsg := "Error: "
	if userErr == "NotFound" && groupErr == "NotFound" {
		errMsg += "Invalid " + s1 + " and Invalid " + s2
	} else if userErr == "NotFound" {
		errMsg += "Invalid " + s1
	} else if groupErr == "NotFound" {
		errMsg += "Invalid " + s2
	}
	return errMsg
}

// verifyTokenUser verify Token User
func verifyTokenUser(uuidSlice []string, db core.DBService) (bool, string) {
	checkUser := db.FindOneUser(bson.D{{"Id", uuidSlice[0]}})
	checkgroup := db.FindOneGroup(bson.D{{"Id", uuidSlice[1]}})
	if checkUser.Username == "NotFound" || checkgroup.Name == "NotFound" {
		return false, returnInitialCheckErrMsg(checkUser.Username, checkgroup.Name, "User", "group")
	}
	// get User's and User's group docs based on token's user uuid
	chk := false
	for _, cId := range checkUser.GroupIds {
		if uuidSlice[1] == cId {
			chk = true
		}
	}
	if !chk {
		return false, "Not Found"
	}
	return true, "No Error"
}

// tokenVerifyMiddleWare
func tokenVerifyMiddleWare(roleType string, next http.HandlerFunc,
	db core.DBService, w http.ResponseWriter, r *http.Request) {
	var MySigningKey = []byte(os.Getenv("SECRET"))
	var errorObject JWTError
	authToken := r.Header.Get("Auth-Token")
	if db.CheckBlacklist(authToken) {
		errorObject.Message = "Invalid Token"
		respondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(MySigningKey), nil
	})
	if err != nil {
		errorObject.Message = err.Error()
		respondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
	tokenClaims := token.Claims.(jwt.MapClaims)
	uuidSlice := []string{tokenClaims["Id"].(string), tokenClaims["GroupId"].(string)}
	verified, verifyMsg := verifyTokenUser(uuidSlice, db)
	if verified {
		if token.Valid && roleType == "Admin" {
			// TODO IF ADMIN CHECK USER ROLE
			next.ServeHTTP(w, r)
		} else if token.Valid && roleType != "Admin" {
			// TODO IF NON-ADMIN CHECK USER ROLE
			next.ServeHTTP(w, r)
		} else {
			errorObject.Message = "Invalid Token"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	} else {
		errorObject.Message = verifyMsg
		respondWithError(w, http.StatusUnauthorized, errorObject)
		return
	}
}

// AdminTokenVerifyMiddleWare is used to verify that the requester is a valid admin
func AdminTokenVerifyMiddleWare(next http.HandlerFunc, db core.DBService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenVerifyMiddleWare("Admin", next, db, w, r)
		return
	})
}

// MemberTokenVerifyMiddleWare is used to verify that a requester is authenticated
func MemberTokenVerifyMiddleWare(next http.HandlerFunc, db core.DBService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenVerifyMiddleWare("Member", next, db, w, r)
		return
	})
}
