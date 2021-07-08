package main

import (
	core "MessengerDemo/src/messenger/pkg"
	"MessengerDemo/src/messenger/pkg/configuration"
	"MessengerDemo/src/messenger/pkg/internals"
	"MessengerDemo/src/messenger/pkg/mongodb"
	"MessengerDemo/src/messenger/pkg/server"
	"github.com/gofrs/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

type App struct {
	server *server.Server
	client *mongo.Client
	config *configuration.Configuration
}

func (a *App) Initialize(mode string) error {
	var err error
	a.config = configuration.ConfigurationSettings(mode)
	a.config.InitializeEnvironmentals()
	mongoURI := os.Getenv("MONGO_URI")
	if mode == "test" {
		mongoURI = "mongodb://127.0.0.1:27017/test"
	}
	a.client, err = mongodb.DatabaseConn(mongoURI)
	if err != nil {
		log.Fatalln("unable to connect to mongodb")
	}
	e := internals.NewEncryptionService(os.Getenv("HOST"), os.Getenv("PORT"))
	db := mongodb.NewDBService(a.client, os.Getenv("DBNAME"))
	g := mongodb.NewGroupService(db)
	u := mongodb.NewUserService(db)
	m := mongodb.NewMessageService(db)
	var group core.Group
	var adminUser core.User
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	group.Name = "admin"
	docCount := db.CountGroups(bson.D{{}})
	if docCount == 0 {
		group.Id = curid.String()
		adminGroup := g.GroupCreate(group)
		adminUser.Username = a.config.MasterAdminUsername
		adminUser.Email = "test"
		adminUser.Password = a.config.MasterAdminInitialPassword
		adminUser.GroupIds = append(adminUser.GroupIds, adminGroup.Id)
		u.UserCreate(adminUser)
	}
	//NewServer(db core.DBService, u core.UserService, g core.GroupService, e *internals.EncryptionService) (error, *Server)
	_, a.server = server.NewServer(db, u, g, m, e)
	return nil
}

func (a *App) Run() {
	//defer a.client.Close()
	a.server.Start()
}