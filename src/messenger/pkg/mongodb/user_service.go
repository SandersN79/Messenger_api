package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
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
func (p *UserService) AuthenticateUser(user core.User) core.User {
	coreUser := p.db.FindOneUser(bson.D{{"Email", user.Email}})
	password := []byte(user.Password)
	checkPassword := []byte(coreUser.Password)
	err := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err == nil {
		return coreUser
	}
	return core.User{Password: "Incorrect"}
}


// BlacklistAuthToken is used during signout to add the now invalid auth-token/api key to the blacklist collection
func (p *UserService) BlacklistAuthToken(authToken string) {
	var blacklist core.Blacklist
	blacklist.AuthToken = authToken
	currentTime := time.Now().UTC()
	blacklist.Created = currentTime.String()
	p.db.BlacklistToken(blacklist)
}

// RefreshToken is used to refresh an existing & valid JWT token
func (p *UserService) RefreshToken(tokenData []string) core.User {
	if tokenData[0] == "" {
		return core.User{Id: ""}
	}
	userUuid := tokenData[0]
	user := p.UserFind(userUuid)
	return user
}

// UpdatePassword is used to update the currently logged in user's password
func (p *UserService) UpdatePassword(tokenData []string, CurrentPassword string, newPassword string) core.User {
	userUuid := tokenData[0]
	curUser := p.db.FindOneUser(bson.D{{"id", userUuid}})
	password := []byte(CurrentPassword)
	checkPassword := []byte(curUser.Password)
	err := bcrypt.CompareHashAndPassword(checkPassword, password)
	if err == nil {
		// 3. Update doc with new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		filter := bson.D{{"Id", curUser.Id}}
		update := bson.D{{"$set",
			bson.D{
				{"Password", string(hashedPassword)},
			},
		}}
		return p.db.UpdateUser(update, filter)
	}
	return core.User{Password: "Incorrect"}
}

// UserCreate is used to create a new user
func (p *UserService) UserCreate(user core.User) core.User {
	docCount := p.db.CountUsers(bson.D{})
	checkUserEmail := p.db.FindOneUser(bson.D{{"Email", user.Email}})

	if checkUserEmail.Username != "NotFound" {
		return core.User{Email: "Taken"}
	}
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	user.Id = curid.String()
	user.UUserKey = core.KeyGen()
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
func (p *UserService) UserDelete(id string) core.User {
	return p.db.DeleteUser(bson.D{{"Id", id}})
}
//TODO change group_uuid to group_uuids
// UsersFind is used to find all user docs
func (p *UserService) UsersFind(groupIds []string) []core.User {
	findFilter := bson.D{}
	findFilter = bson.D{{}}
	return p.db.FindUsers(findFilter)
}

// UserFind is used to find a specific user doc
func (p *UserService) UserFind(id string) core.User {
	findFilter := bson.D{{"Id", id}}
	return p.db.FindOneUser(findFilter)
}

// UserUpdate is used to update an existing user doc
func (p *UserService) UserUpdate(user core.User) core.User {
	docCount := p.db.CountUsers(bson.D{})
	findFilter := bson.D{{"Id", user.Id}}
	curUser := p.db.FindOneUser(findFilter)
	if curUser.Username == "NotFound" {
		return core.User{Id: "Not Found"}
	}
	checkUser := p.db.FindOneUser(bson.D{{"Username", user.Username}})
	checkUserEmail := p.db.FindOneUser(bson.D{{"Email", user.Email}})
	user = BaseModifyUser(user, newUserModel(curUser))
	if checkUser.Username != "NotFound" && curUser.Username != user.Username {
		return core.User{Username: "Taken"}
	} else if checkUserEmail.Username != "NotFound" && curUser.Email != user.Email {
		return core.User{Email: "Taken"}
	}
	if docCount == 0 {
		return core.User{Id: "NotFound"}
	}
	filter := bson.D{{"Id", user.Id}}
	if len(user.Password) != 0 {
		password := []byte(user.Password)
		hashedPassword, hashErr := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if hashErr != nil {
			panic(hashErr)
		}
		update := bson.D{{"$set",
			bson.D{
				{"Password", string(hashedPassword)},
				{"Username", user.Username},
				{"Email", user.Email},
			},
		}}
		p.db.UpdateUser(update, filter)
		user.Password = ""
		return user
	}
	update := bson.D{{"$set",
		bson.D{
			{"Username", user.Username},
			{"Email", user.Email},
		},
	}}
	p.db.UpdateUser(update, filter)
	return user
}

// UserDocInsert is used to insert an user doc directly into mongodb for testing purposes
func (p *UserService) UserDocInsert(user core.User) core.User {
	password := []byte(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	user.Password = string(hashedPassword)
	return p.db.InsertUser(user)
}


// BaseModifyUser is a function that setups the base user struct during a user modification request
func BaseModifyUser(user core.User, curUser *users) core.User {
	if len(user.Username) == 0 {
		user.Username = curUser.Username
	}
	if len(user.Email) == 0 {
		user.Email = curUser.Email
	}

	return user
}
