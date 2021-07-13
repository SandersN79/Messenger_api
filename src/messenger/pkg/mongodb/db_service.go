package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// DBService is used by the app to manage all plan related controllers and functionality
type DBService struct {
	dbName                 string
	usersCollection        *mongo.Collection
	blacklistCollection    *mongo.Collection
	groupsCollection       *mongo.Collection
	messagesCollection     *mongo.Collection
	client                 *mongo.Client
}

// NewDBService is an exported function used to initialize a new DatasourceService struct
func NewDBService(client *mongo.Client, dbName string) *DBService {
	usersCollection := client.Database(dbName).Collection("users")
	blacklistCollection := client.Database(dbName).Collection("blacklists")
	groupsCollection := client.Database(dbName).Collection("groups")
	messagesCollection := client.Database(dbName).Collection("messages")
	return &DBService{
		dbName,
		usersCollection,
		blacklistCollection,
		groupsCollection,
		messagesCollection,
		client,
	}
}

/////////////////


// FindOneUser finds a specfic user
func (p *DBService) FindOneUser(filter bson.D) core.User {
	var model = newUserModel(core.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.usersCollection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return core.User{Username: "NotFound"}
	}
	return model.toRootUser()
}

// FindUsers finds all users
func (p *DBService) FindUsers(filter bson.D) []core.User {
	var coreStructs []core.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.usersCollection.Find(ctx, filter)
	if err != nil {
		defer cursor.Close(ctx)
		return coreStructs
	}
	for cursor.Next(ctx) {
		result := newUserModel(core.User{})
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("cursor.Next() error:", err)
			panic(err)
		} else {
			result.Password = ""
			coreStructs = append(coreStructs, result.toRootUser())
		}
	}
	return coreStructs
}

// CreateUser is used to create a new user
func (p *DBService) CreateUser(coreStruct core.User) core.User {
	var check = newUserModel(core.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	nameErr := p.usersCollection.FindOne(ctx, bson.M{"username": coreStruct.Username}).Decode(&check)
	if nameErr == nil {
		return core.User{Username: "Taken"}
	}
	currentTime := time.Now().UTC()
	coreStruct.Created = currentTime.String()
	model := newUserModel(coreStruct)
	_, err := p.usersCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("user doc creation error err: ", err)
	}
	return model.toRootUser()
}

// UpdateUser is used to update a user
func (p *DBService) UpdateUser(update bson.D, filter bson.D) core.User {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		return core.User{Username: "Error"}
	}
	return core.User{Username: "Success"}
}

// DeleteUser is used to delete a user
func (p *DBService) DeleteUser(filter bson.D) core.User {
	var model = newUserModel(core.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.usersCollection.FindOneAndDelete(ctx, filter).Decode(&model)
	if err != nil {
		return core.User{}
	}
	return model.toRootUser()
}

// InsertUser
func (p *DBService) InsertUser(coreStruct core.User) core.User {
	var model = newUserModel(coreStruct)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.usersCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("user doc insertion error: ", err)
	}
	return model.toRootUser()
}

// CountUsers
func (p *DBService) CountUsers(filter bson.D) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.usersCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}
	return docCount
}

///////////////


// FindOneGroup finds one group
func (p *DBService) FindOneGroup(filter bson.D) core.Group {
	var model = newGroupModel(core.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.groupsCollection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return core.Group{Name: "NotFound"}
	}
	return model.toRootGroup()
}

// FindGroups
func (p *DBService) FindGroups(filter bson.D) []core.Group {
	var coreStructs []core.Group
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.groupsCollection.Find(ctx, filter)
	if err != nil {
		defer cursor.Close(ctx)
		return coreStructs
	}
	for cursor.Next(ctx) {
		result := newGroupModel(core.Group{})
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("cursor.Next() error:", err)
			panic(err)
		} else {
			coreStructs = append(coreStructs, result.toRootGroup())
		}
	}
	return coreStructs
}

// CreateGroup is used to create a new group
func (p *DBService) CreateGroup(coreStruct core.Group) core.Group {
	var check = newGroupModel(core.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	nameErr := p.groupsCollection.FindOne(ctx, bson.M{"Name": coreStruct.Name}).Decode(&check)
	if nameErr == nil {
		return core.Group{Name: "Taken"}
	}
	currentTime := time.Now().UTC()
	coreStruct.Created = currentTime.String()
	model := newGroupModel(coreStruct)
	_, err := p.groupsCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("group doc creation error err: ", err)
	}
	return model.toRootGroup()
}

// UpdateGroup is used to update a group
func (p *DBService) UpdateGroup(update bson.D, filter bson.D) core.Group {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.groupsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		return core.Group{Name: "Error"}
	}
	return core.Group{Name: "Success"}
}

// DeleteGroup is used to delete a group
func (p *DBService) DeleteGroup(filter bson.D) core.Group {
	var model = newGroupModel(core.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.groupsCollection.FindOneAndDelete(ctx, filter).Decode(&model)
	if err != nil {
		return core.Group{}
	}
	return model.toRootGroup()
}

// InsertGroup
func (p *DBService) InsertGroup(coreStruct core.Group) core.Group {
	var model = newGroupModel(coreStruct)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.groupsCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("group doc insertion error: ", err)
	}
	return model.toRootGroup()
}

// CountGroups
func (p *DBService) CountGroups(filter bson.D) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.groupsCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}
	return docCount
}
//////////////////////


// FindOneMessage Function to get a message
func (p *DBService) FindOneMessage(filter bson.D) core.Message {
	var model = newMessageModel(core.Message{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.messagesCollection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return core.Message{Id: "NotFound"}
	}
	return model.toRootMessage()
}

// FindMessages is used to find all messages related to user
func (p *DBService) FindMessages(filter bson.D) []core.Message {
	var coreStructs []core.Message
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.messagesCollection.Find(ctx, filter)
	if err != nil {
		defer cursor.Close(ctx)
		return coreStructs
	}
	for cursor.Next(ctx) {
		result := newMessageModel(core.Message{})
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("cursor.Next() error:", err)
			panic(err)
		} else {
			coreStructs = append(coreStructs, result.toRootMessage())
		}
	}
	return coreStructs
}

// CreateMessage is used to create a new message
func (p *DBService) CreateMessage(coreStruct core.Message) core.Message {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	currentTime := time.Now().UTC()
	coreStruct.Created = currentTime.String()
	model := newMessageModel(coreStruct)
	_, err := p.messagesCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("message doc creation error err: ", err)
	}
	return model.toRootMessage()
}

// UpdateMessage is used to update a message
func (p *DBService) UpdateMessage(update bson.D, filter bson.D) core.Message {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.messagesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		return core.Message{Id: "Error"}
	}
	return core.Message{Id: "Success"}
}

// DeleteMessage is used to delete a message
func (p *DBService) DeleteMessage(filter bson.D) core.Message {
	var model = newMessageModel(core.Message{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.messagesCollection.FindOneAndDelete(ctx, filter).Decode(&model)
	if err != nil {
		return core.Message{}
	}
	return model.toRootMessage()
}

// InsertMessage
func (p *DBService) InsertMessage(coreStruct core.Message) core.Message {
	var model = newMessageModel(coreStruct)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.messagesCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("message doc insertion error: ", err)
	}
	return model.toRootMessage()
}

// CountMessages
func (p *DBService) CountMessages(filter bson.D) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	docCount, err := p.messagesCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0
	}
	return docCount
}

////////


// BlacklistToken
func (p *DBService) BlacklistToken(coreStruct core.Blacklist) {
	model := newBlacklistModel(coreStruct)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	p.blacklistCollection.InsertOne(ctx, model)
}

// CheckBlacklist
func (p *DBService) CheckBlacklist(authToken string) bool {
	var checkToken core.Blacklist
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	blacklistErr := p.blacklistCollection.FindOne(ctx, bson.M{"auth_token": authToken}).Decode(&checkToken)
	if blacklistErr != nil {
		return false
	}
	return true
}
