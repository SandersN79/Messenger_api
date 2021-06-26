package mongodb

import (
	"MessengerDemo/server"
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


// FindOneUser Function to get a user from datasource with custom filter
func (p *DBService) FindOneUser(filter bson.D) server.User {
	var model = newUserModel(server.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.usersCollection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return server.User{Username: "NotFound"}
	}
	return model.toRootUser()
}

// FindUsers Function to get a company from datasource with custom filter
func (p *DBService) FindUsers(filter bson.D) []server.User {
	var serverStructs []server.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.usersCollection.Find(ctx, filter)
	if err != nil {
		defer cursor.Close(ctx)
		return serverStructs
	}
	for cursor.Next(ctx) {
		result := newUserModel(server.User{})
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("cursor.Next() error:", err)
			panic(err)
		} else {
			result.Password = ""
			serverStructs = append(serverStructs, result.toRootUser())
		}
	}
	return serverStructs
}

// CreateUser is used to create a new user user
func (p *DBService) CreateUser(serverStruct server.User) server.User {
	var check = newUserModel(server.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	nameErr := p.usersCollection.FindOne(ctx, bson.M{"username": serverStruct.Username}).Decode(&check)
	if nameErr == nil {
		return server.User{Username: "Taken"}
	}
	currentTime := time.Now().UTC()
	serverStruct.Created = currentTime.String()
	model := newUserModel(serverStruct)
	_, err := p.usersCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("user doc creation error err: ", err)
	}
	return model.toRootUser()
}

// UpdateUser is used to create a new user user
func (p *DBService) UpdateUser(update bson.D, filter bson.D) server.User {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		return server.User{Username: "Error"}
	}
	return server.User{Username: "Success"}
}

// DeleteUser is used to create a new user user
func (p *DBService) DeleteUser(filter bson.D) server.User {
	var model = newUserModel(server.User{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.usersCollection.FindOneAndDelete(ctx, filter).Decode(&model)
	if err != nil {
		return server.User{}
	}
	return model.toRootUser()
}

// InsertUser
func (p *DBService) InsertUser(serverStruct server.User) server.User {
	var model = newUserModel(serverStruct)
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


// FindOneGroup Function to get a group from datasource with custom filter
func (p *DBService) FindOneGroup(filter bson.D) server.Group {
	var model = newGroupModel(server.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.groupsCollection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return server.Group{Name: "NotFound"}
	}
	return model.toRootGroup()
}

// FindGroups Function to get a company from datasource with custom filter
func (p *DBService) FindGroups(filter bson.D) []server.Group {
	var serverStructs []server.Group
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.groupsCollection.Find(ctx, filter)
	if err != nil {
		defer cursor.Close(ctx)
		return serverStructs
	}
	for cursor.Next(ctx) {
		result := newGroupModel(server.Group{})
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("cursor.Next() error:", err)
			panic(err)
		} else {
			serverStructs = append(serverStructs, result.toRootGroup())
		}
	}
	return serverStructs
}

// CreateGroup is used to create a new group group
func (p *DBService) CreateGroup(serverStruct server.Group) server.Group {
	var check = newGroupModel(server.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	nameErr := p.groupsCollection.FindOne(ctx, bson.M{"Name": serverStruct.Name}).Decode(&check)
	if nameErr == nil {
		return server.Group{Name: "Taken"}
	}
	currentTime := time.Now().UTC()
	serverStruct.Created = currentTime.String()
	model := newGroupModel(serverStruct)
	_, err := p.groupsCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("group doc creation error err: ", err)
	}
	return model.toRootGroup()
}

// UpdateGroup is used to create a new group group
func (p *DBService) UpdateGroup(update bson.D, filter bson.D) server.Group {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.groupsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		return server.Group{Name: "Error"}
	}
	return server.Group{Name: "Success"}
}

// DeleteGroup is used to create a new group group
func (p *DBService) DeleteGroup(filter bson.D) server.Group {
	var model = newGroupModel(server.Group{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.groupsCollection.FindOneAndDelete(ctx, filter).Decode(&model)
	if err != nil {
		return server.Group{}
	}
	return model.toRootGroup()
}

// InsertGroup
func (p *DBService) InsertGroup(serverStruct server.Group) server.Group {
	var model = newGroupModel(serverStruct)
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


// FindOneMessage Function to get a message from datasource with custom filter
func (p *DBService) FindOneMessage(filter bson.D) server.Message {
	var model = newMessageModel(server.Message{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.messagesCollection.FindOne(ctx, filter).Decode(&model)
	if err != nil {
		return server.Message{Id: "NotFound"}
	}
	return model.toRootMessage()
}

// FindMessages Function to get a company from datasource with custom filter
func (p *DBService) FindMessages(filter bson.D) []server.Message {
	var serverStructs []server.Message
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := p.messagesCollection.Find(ctx, filter)
	if err != nil {
		defer cursor.Close(ctx)
		return serverStructs
	}
	for cursor.Next(ctx) {
		result := newMessageModel(server.Message{})
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println("cursor.Next() error:", err)
			panic(err)
		} else {
			serverStructs = append(serverStructs, result.toRootMessage())
		}
	}
	return serverStructs
}

// CreateMessage is used to create a new message message
func (p *DBService) CreateMessage(serverStruct server.Message) server.Message {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	currentTime := time.Now().UTC()
	serverStruct.Created = currentTime.String()
	model := newMessageModel(serverStruct)
	_, err := p.messagesCollection.InsertOne(ctx, model)
	if err != nil {
		fmt.Println("message doc creation error err: ", err)
	}
	return model.toRootMessage()
}

// UpdateMessage is used to create a new message message
func (p *DBService) UpdateMessage(update bson.D, filter bson.D) server.Message {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := p.messagesCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Println("update err: ", err)
		return server.Message{Id: "Error"}
	}
	return server.Message{Id: "Success"}
}

// DeleteMessage is used to create a new message message
func (p *DBService) DeleteMessage(filter bson.D) server.Message {
	var model = newMessageModel(server.Message{})
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := p.messagesCollection.FindOneAndDelete(ctx, filter).Decode(&model)
	if err != nil {
		return server.Message{}
	}
	return model.toRootMessage()
}

// InsertMessage
func (p *DBService) InsertMessage(serverStruct server.Message) server.Message {
	var model = newMessageModel(serverStruct)
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
