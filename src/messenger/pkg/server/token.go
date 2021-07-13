package server

import (
	core "MessengerDemo/src/messenger/pkg"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
)

// CreateToken is used to create a new session JWT token
func CreateToken(user core.User, exp int64) string {
	var MySigningKey = []byte(os.Getenv("SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["SystemRole"] = user.SystemRole
	claims["Username"] = user.Username
	claims["Id"] = user.Id
	claims["GroupId"] = user.GroupIds[0]
	claims["exp"] = exp
	tokenString, _ := token.SignedString(MySigningKey)
	return tokenString
}

// DecodeJWT is used to decode a JWT token
func DecodeJWT(curToken string) []string {
	var MySigningKey = []byte(os.Getenv("SECRET"))
	// Decode the token
	token, err := jwt.Parse(curToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(MySigningKey), nil
	})
	if err != nil {
		return []string{"", ""}
	}
	// Determine user from token
	tokenClaims := token.Claims.(jwt.MapClaims)
	userUuid := tokenClaims["Id"].(string)
	userName := tokenClaims["Username"].(string)
	groupIds := tokenClaims["GroupId"].(string)
	var reSlice []string
	if len(groupIds) != 0 {
		reSlice = []string{userUuid, groupIds, userName}
	}
	return reSlice
}
