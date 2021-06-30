package server

import (
	core "MessengerDemo/src/messenger/pkg"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type groupRouter struct {
	databaseService      core.DBService
	groupService         core.GroupService
}

// NewGroupRouter is a function that initializes a new groupRouter struct
func NewGroupRouter(db core.DBService, g core.GroupService, router *mux.Router) *mux.Router {
	groupRouter := groupRouter{db, g}
	router.HandleFunc("/groups", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups", AdminTokenVerifyMiddleWare(groupRouter.GroupsShow, db)).Methods("GET")
	router.HandleFunc("/groups/{groupId}", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/groups/{groupId}", AdminTokenVerifyMiddleWare(groupRouter.GroupShow, db)).Methods("GET")
	router.HandleFunc("/groups", AdminTokenVerifyMiddleWare(groupRouter.CreateGroup, db)).Methods("POST")
	router.HandleFunc("/groups/{groupId}", AdminTokenVerifyMiddleWare(groupRouter.DeleteGroup, db)).Methods("DELETE")
	router.HandleFunc("/groups/{groupId}", AdminTokenVerifyMiddleWare(groupRouter.ModifyGroup, db)).Methods("PATCH")
	return router
}

// Handler to show all groups
func (gr *groupRouter) ModifyGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	var group core.Group
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &group); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	group.Id = groupId
	g := gr.groupService.GroupUpdate(group)
	if g.Id == "NotFound" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Group Not Found"}); err != nil {
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

// Handler to show all groups
func (gr *groupRouter) GroupsShow(w http.ResponseWriter, r *http.Request) {
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	groups := gr.groupService.GroupsFind()
	if err := json.NewEncoder(w).Encode(groups); err != nil {
		panic(err)
	}
}

// Handler to show a specific group
func (gr *groupRouter) GroupShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	group := gr.groupService.GroupFind(core.Group{Id: groupId})
	if group.Id != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(group); err != nil {
			panic(err)
		}
		return
	}
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Group Not Found"}); err != nil {
		panic(err)
	}
}

// Handler to create an group
func (gr *groupRouter) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group core.Group
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &group); err != nil {
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
	group.Id = curid.String()
	g := gr.groupService.GroupCreate(group)
	if g.Name == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Group Name Taken"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(g); err != nil {
			panic(err)
		}
	}
}

// Handler to delete an group
func (gr *groupRouter) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupId := vars["groupId"]
	group := gr.groupService.GroupDelete(groupId)
	if group.Id != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("Group Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Group Not Found"}); err != nil {
		panic(err)
	}
}
