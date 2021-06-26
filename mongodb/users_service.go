package mongodb

import (
	"MessengerDemo/server"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// UserService is used by the app to manage all user related controllers and functionality
type UserService struct {
	db         *DBService
}

// NewUserService is an exported function used to initialize a new UserService struct
func NewUserService(db *DBService) *UserService {
	return &UserService{db}
}

// AuthenticateUser is used to authenticate users that are signing in
func (p *UserService) AuthenticateUser(user server.User) server.User {
	serverUser := p.db.FindOneUser(bson.D{{"username", user.Username}})
	password := []byte(user.Password)
	checkPassword := []byte(serverUser.Password)
	err := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err == nil {
		return serverUser
	}
	return server.User{Password: "Incorrect"}
}

/*
// BlacklistAuthToken is used during signout to add the now invalid auth-token/api key to the blacklist collection
func (p *UserService) BlacklistAuthToken(authToken string) {
	var blacklist server.Blacklist
	blacklist.AuthToken = authToken
	currentTime := time.Now().UTC()
	blacklist.LastModified = currentTime.String()
	blacklist.CreationDatetime = currentTime.String()
	p.db.BlacklistToken(blacklist)
}

// RefreshToken is used to refresh an existing & valid JWT token
func (p *UserService) RefreshToken(tokenData []string, groupUuid string) server.User {
	if tokenData[0] == "" {
		return server.User{Uuid: ""}
	}
	userUuid := tokenData[0]
	user := p.UserFind(userUuid, groupUuid)
	return user
}
*/

// UpdatePassword is used to update the currently logged in user's password
func (p *UserService) UpdatePassword(tokenData []string, CurrentPassword string, newPassword string) server.User {
	userUuid := tokenData[0]
	curUser := p.db.FindOneUser(bson.D{{"id", userUuid}})
	password := []byte(CurrentPassword)
	checkPassword := []byte(curUser.Password)
	err := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err == nil {
		// 3. Update doc with new password
		currentTime := time.Now().UTC()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		filter := bson.D{{"uuid", curUser.Id}}
		update := bson.D{{"$set",
			bson.D{
				{"password", string(hashedPassword)},
				{"last_modified", currentTime.String()},
			},
		}}
		return p.db.UpdateUser(update, filter)
	}
	return server.User{Password: "Incorrect"}
}

// UserCreate is used to create a new user
func (p *UserService) UserCreate(user server.User) server.User {
	docCount := p.db.CountUsers(bson.D{})
	checkUserEmail := p.db.FindOneUser(bson.D{{"Email", user.Email}})

	if checkUserEmail.Username != "NotFound" {
		return server.User{Email: "Taken"}
	}
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	user.Id = curid.String()
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
	if docCount == 0 {
		user.SystemRole = "master_admin"
	} else {
		user.SystemRole = "member"
	}
	currentTime := time.Now().UTC()
	user.Created = currentTime.String()
	return p.db.CreateUser(user)
}

// UserDelete is used to delete an user
func (p *UserService) UserDelete(id string) server.User {
	return p.db.DeleteUser(bson.D{{"Id", id}})
}
//TODO change group_uuid to group_uuids
// UsersFind is used to find all user docs
func (p *UserService) UsersFind(groupUuid string) []server.User {
	findFilter := bson.D{}
	if groupUuid != "" {
		findFilter = bson.D{{"group_uuid", groupUuid}}
	}
	return p.db.FindUsers(findFilter)
}

// UserFind is used to find a specific user doc
func (p *UserService) UserFind(id string) server.User {
	findFilter := bson.D{{"uuid", id}}
	return p.db.FindOneUser(findFilter)
}

// UserUpdate is used to update an existing user doc
func (p *UserService) UserUpdate(user server.User) server.User {
	docCount := p.db.CountUsers(bson.D{})
	findFilter := bson.D{{"Id", user.Id}}
	curUser := p.db.FindOneUser(findFilter)
	if curUser.Username == "NotFound" {
		return server.User{Id: "Not Found"}
	}
	checkUser := p.db.FindOneUser(bson.D{{"username", user.Username}})
	checkUserEmail := p.db.FindOneUser(bson.D{{"email", user.Email}})
	user = BaseModifyUser(user, newUserModel(curUser))
	if checkUser.Username != "NotFound" && curUser.Username != user.Username {
		return server.User{Username: "Taken"}
	} else if checkUserEmail.Username != "NotFound" && curUser.Email != user.Email {
		return server.User{Email: "Taken"}
	}
	if docCount == 0 {
		return server.User{Id: "NotFound"}
	}
	filter := bson.D{{"uuid", user.Id}}
	currentTime := time.Now().UTC()
	if len(user.Password) != 0 {
		password := []byte(user.Password)
		hashedPassword, hashErr := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if hashErr != nil {
			panic(hashErr)
		}
		update := bson.D{{"$set",
			bson.D{
				{"password", string(hashedPassword)},

				{"username", user.Username},
				{"email", user.Email},
				{"last_modified", currentTime.String()},
			},
		}}
		p.db.UpdateUser(update, filter)
		user.Password = ""
		return user
	}
	update := bson.D{{"$set",
		bson.D{

			{"username", user.Username},
			{"email", user.Email},
			{"last_modified", currentTime.String()},
		},
	}}
	p.db.UpdateUser(update, filter)
	return user
}

// UserDocInsert is used to insert an user doc directly into mongodb for testing purposes
func (p *UserService) UserDocInsert(user server.User) server.User {
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
	return p.db.InsertUser(user)
}


// BaseModifyUser is a function that setups the base user struct during a user modification request
func BaseModifyUser(user server.User, curUser *users) server.User {
	if len(user.Username) == 0 {
		user.Username = curUser.Username
	}
	if len(user.Email) == 0 {
		user.Email = curUser.Email
	}

	return user
}