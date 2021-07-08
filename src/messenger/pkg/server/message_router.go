package server

import (
	core "MessengerDemo/src/messenger/pkg"
	"MessengerDemo/src/messenger/pkg/internals"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type messageRouter struct {
	databaseService core.DBService
	messageService core.MessageService
	encryptionService *internals.EncryptionService
}

func NewMessageRouter(db core.DBService, m core.MessageService, e *internals.EncryptionService, router *mux.Router) *mux.Router {
	messageRouter :=  messageRouter{db, m, e}
	router.HandleFunc("/messages", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/messages", AdminTokenVerifyMiddleWare(messageRouter.MessagesShow, db)).Methods("GET")
	router.HandleFunc("/messages/{messageId}", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/messages/{messageId}", AdminTokenVerifyMiddleWare(messageRouter.MessagesShow, db)).Methods("GET")
	router.HandleFunc("/messages", AdminTokenVerifyMiddleWare(messageRouter.CreateMessage, db)).Methods("POST")
	router.HandleFunc("/messages/{messageId}", AdminTokenVerifyMiddleWare(messageRouter.DeleteMessage, db)).Methods("DELETE")
	router.HandleFunc("/messages/{messageId}", AdminTokenVerifyMiddleWare(messageRouter.ModifyMessage, db)).Methods("PATCH")
	return router
}


// Handler to show all messages
func (m *messageRouter) ModifyMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	var message core.Message
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &message); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	message.Id = messageId
	//TODO add encryption functionality below
	eContents, err := m.encryptionService.Encrypt(message.Contents)
	if err !=  nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	message.Contents = eContents
	g := m.messageService.MessageUpdate(message)
	if g.Id == "NotFound" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Message Not Found"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(g); err != nil {
			panic(err)
		}
	}
}

// Handler to show all messages
func (m *messageRouter) MessagesShow(w http.ResponseWriter, r *http.Request) {
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"))
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	smessages := m.messageService.MessagesFind(core.Message{SenderId: decodedToken[0]})
	rmessages := m.messageService.MessagesFind(core.Message{ReceiverId: decodedToken[0]})
	messages := append(smessages, rmessages...)
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		panic(err)
	}
}

// Handler to show a specific message
func (m *messageRouter) MessageShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	message := m.messageService.MessageFind(core.Message{Id: messageId})
	if message.Id != "" {
		//TODO add encryption functionality below
		eContents, err := m.encryptionService.Decrypt(message.Contents)
		if err !=  nil {
			w = SetResponseHeaders(w, "", "")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
			return
		}
		message.Contents = eContents
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(message); err != nil {
			panic(err)
		}
		return
	}
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Message Not Found"}); err != nil {
		panic(err)
	}
}

// Handler to create an message
func (m *messageRouter) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message core.Message
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &message); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	curid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	message.Id = curid.String()
	//TODO add encryption functionality below
	eContents, err := m.encryptionService.Encrypt(message.Contents)
	if err !=  nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}
	message.Contents = eContents
	g := m.messageService.MessageCreate(message)
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(g); err != nil {
		panic(err)
	}
}

// Handler to delete an message
func (m *messageRouter) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	messageId := vars["messageId"]
	message := m.messageService.MessageDelete(core.Message{Id: messageId})
	if message.Id != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Message Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Message Not Found"}); err != nil {
		panic(err)
	}
}
