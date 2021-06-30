package core

import (
	"go.mongodb.org/mongo-driver/bson"
)

// DBService is an interface used to manage the relevant user doc controllers
type DBService interface {
	// Users
	FindOneUser(filter bson.D) User
	FindUsers(filter bson.D) []User
	CreateUser(rootStruct User) User
	UpdateUser(update bson.D, filter bson.D) User
	DeleteUser(filter bson.D) User
	InsertUser(rootStruct User) User
	CountUsers(filter bson.D) int64
	// Blacklists
	BlacklistToken(rootStruct Blacklist)
	CheckBlacklist(authToken string) bool
	// Groups
	FindOneGroup(filter bson.D) Group
	FindGroups(filter bson.D) []Group
	CreateGroup(rootStruct Group) Group
	UpdateGroup(update bson.D, filter bson.D) Group
	DeleteGroup(filter bson.D) Group
	InsertGroup(rootStruct Group) Group
	CountGroups(filter bson.D) int64
	// Messages
	FindOneMessage(filter bson.D) Message
	FindMessages(filter bson.D) []Message
	CreateMessage(rootStruct Message) Message
	UpdateMessage(update bson.D, filter bson.D) Message
	DeleteMessage(filter bson.D) Message
	InsertMessage(rootStruct Message) Message
	CountMessages(filter bson.D) int64

}