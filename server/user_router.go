package main

import (
	"codegen-api/src/rest_api/pkg/configuration"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type userRouter struct {
	databaseService      DBService
	userService          UserService
	groupService         GroupService
	config               configuration.Configuration
}

// NewUserRouter is a function that initializes a new userRouter struct
func NewUserRouter(db DBService, u UserService, router *mux.Router, o GroupService, config configuration.Configuration) *mux.Router {
	userRouter := userRouter{db, u, o, config}
	router.HandleFunc("/auth", HandleOptionsRequest).Methods("OPTIONS")
	// 1.5*: HIT THIS ROUTE IF THE USER WAS CREATED VIA POST /USERS TO GET AUTH TOKEN TO USER IN STEP 2
	router.HandleFunc("/auth", userRouter.Signin).Methods("POST")
	router.HandleFunc("/auth", MemberTokenVerifyMiddleWare(userRouter.RefreshSession, config, db)).Methods("GET")
	router.HandleFunc("/auth", MemberTokenVerifyMiddleWare(userRouter.Signout, config, db)).Methods("DELETE")
	// 1A: HIT THIS ROUTE FOR NEW ACCOUNTS
	router.HandleFunc("/auth/register", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/register", userRouter.RegisterUser).Methods("POST")
	// 2: HIT THIS ROUTE FOR THE API KEY USING THE AUTH TOKEN
	router.HandleFunc("/auth/api-key", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/api-key", MemberTokenVerifyMiddleWare(userRouter.GenerateAPIKey, config, db)).Methods("GET")
	router.HandleFunc("/auth/password", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/password", MemberTokenVerifyMiddleWare(userRouter.UpdatePassword, config, db)).Methods("POST")
	router.HandleFunc("/users", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users", MemberTokenVerifyMiddleWare(userRouter.UsersShow, config, db)).Methods("GET")
	router.HandleFunc("/users/{userId}", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users/{userId}", MemberTokenVerifyMiddleWare(userRouter.UserShow, config, db)).Methods("GET")
	// 1B: HIT THIS ROUTE IF COMPANY EXISTS AND A NEW USER IS BEING CREATED BY AN ADMIN
	router.HandleFunc("/users", AdminTokenVerifyMiddleWare(userRouter.CreateUser, config, db)).Methods("POST")
	router.HandleFunc("/users/{userId}", AdminTokenVerifyMiddleWare(userRouter.DeleteUser, config, db)).Methods("DELETE")
	router.HandleFunc("/users/{userId}", MemberTokenVerifyMiddleWare(userRouter.ModifyUser, config, db)).Methods("PATCH")
	return router
}

// Handler function that manages the user signin process
func (ur *userRouter) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"), ur.config)
	type passwordStruct struct {
		NewPassword     string `json:"new_password"`
		CurrentPassword string `json:"current_password"`
	}
	var pw passwordStruct
	err = json.Unmarshal(body, &pw)
	if err != nil {
		panic(err)
	}
	u := ur.userService.UpdatePassword(decodedToken, pw.CurrentPassword, pw.NewPassword)
	if u.Password == "Incorrect" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Incorrect Current Password Provided"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}