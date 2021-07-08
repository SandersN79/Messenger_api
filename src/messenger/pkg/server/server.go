package server

import (
	core "MessengerDemo/src/messenger/pkg"
	"MessengerDemo/src/messenger/pkg/internals"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Router            *mux.Router
	DBService         core.DBService
	UserService       core.UserService
	GroupService      core.GroupService
	MessageService    core.MessageService
	EncryptionService *internals.EncryptionService
}
//(db core.DBService, u core.UserService, g core.GroupService, m core.MessageService, e *internals.EncryptionService)
func NewServer(db core.DBService, u core.UserService, g core.GroupService, m core.MessageService, e *internals.EncryptionService) (error, *Server) {
	var server Server
	var err error
	router := mux.NewRouter().StrictSlash(true)
	router = NewUserRouter(db, u, g, router) //add u, g
	router = NewGroupRouter(db, g, router)
	router = NewMessageRouter(db, m, e, router)
	server = Server{
		Router: router, DBService: db, UserService: u, GroupService: g, MessageService: m, EncryptionService: e,
	}
	return err, &server
}


func (s *Server) Start() {
	//fmt.Println("CHECK PORT", os.Getenv("PORT"))
	log.Println("Listening on port " + os.Getenv("PORT"))
	portStr := ":" + os.Getenv("PORT")
	if err := http.ListenAndServe(portStr, handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
		log.Fatal("http.ListenAndServe: ", err)
	}
	/*
	if s.Config.HTTPS == "on" {
		log.Println("Listening on port " + s.Config.Port)
		if err := http.ListenAndServeTLS(s.Config.GetPort(), s.Config.Cert, s.Config.Key, handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	} else {
		log.Println("Listening on port " + s.Config.Port)
		if err := http.ListenAndServe(s.Config.GetPort(), handlers.LoggingHandler(os.Stdout, s.Router)); err != nil {
			log.Fatal("http.ListenAndServe: ", err)
		}
	}
	*/

}

