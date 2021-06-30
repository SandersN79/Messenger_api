package server

import (
	core "MessengerDemo/src/messenger/pkg"
	"encoding/json"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type userRouter struct {
	databaseService      core.DBService
	userService          core.UserService
	groupService         core.GroupService
}

// NewUserRouter is a function that initializes a new userRouter struct
func NewUserRouter(db core.DBService, u core.UserService, o core.GroupService, router *mux.Router) *mux.Router {
	userRouter := userRouter{db, u, o}
	router.HandleFunc("/auth", HandleOptionsRequest).Methods("OPTIONS")
	// 1.5*: HIT THIS ROUTE IF THE USER WAS CREATED VIA POST /USERS TO GET AUTH TOKEN TO USER IN STEP 2
	router.HandleFunc("/auth", userRouter.Signin).Methods("POST")
	router.HandleFunc("/auth", MemberTokenVerifyMiddleWare(userRouter.RefreshSession, db)).Methods("GET")
	router.HandleFunc("/auth", MemberTokenVerifyMiddleWare(userRouter.Signout, db)).Methods("DELETE")
	// 1A: HIT THIS ROUTE FOR NEW ACCOUNTS
	router.HandleFunc("/auth/register", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/register", userRouter.RegisterUser).Methods("POST")
	// 2: HIT THIS ROUTE FOR THE API KEY USING THE AUTH TOKEN
	router.HandleFunc("/auth/api-key", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/api-key", MemberTokenVerifyMiddleWare(userRouter.GenerateAPIKey, db)).Methods("GET")
	router.HandleFunc("/auth/password", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/auth/password", MemberTokenVerifyMiddleWare(userRouter.UpdatePassword, db)).Methods("POST")
	router.HandleFunc("/users", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users", MemberTokenVerifyMiddleWare(userRouter.UsersShow, db)).Methods("GET")
	router.HandleFunc("/users/{userId}", HandleOptionsRequest).Methods("OPTIONS")
	router.HandleFunc("/users/{userId}", MemberTokenVerifyMiddleWare(userRouter.UserShow, db)).Methods("GET")
	// 1B: HIT THIS ROUTE IF COMPANY EXISTS AND A NEW USER IS BEING CREATED BY AN ADMIN
	router.HandleFunc("/users", AdminTokenVerifyMiddleWare(userRouter.CreateUser, db)).Methods("POST")
	router.HandleFunc("/users/{userId}", AdminTokenVerifyMiddleWare(userRouter.DeleteUser, db)).Methods("DELETE")
	router.HandleFunc("/users/{userId}", MemberTokenVerifyMiddleWare(userRouter.ModifyUser, db)).Methods("PATCH")
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
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"))
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


// Handler function that manages the user signin process
func (ur *userRouter) ModifyUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	var user core.User
	user.Id = userId
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	u := ur.userService.UserUpdate(user)
	if u.Id == "Not Found" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "User Not Found"}); err != nil {
			panic(err)
		}
	}
	if u.Email == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Email Taken"}); err != nil {
			panic(err)
		}
	} else if u.Username == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

// Handler function that manages the user signin process
func (ur *userRouter) Signin(w http.ResponseWriter, r *http.Request) {
	var user core.User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	} else {
		u := ur.userService.AuthenticateUser(user)
		if u.Username == "NotFound" || u.Password == "Incorrect" {
			w = SetResponseHeaders(w, "", "")
			w.WriteHeader(401)
			if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusUnauthorized, Text: "Incorrect"}); err != nil {
				panic(err)
			}
		} else {
			expDT := time.Now().Add(time.Hour * 1).Unix()
			sessionToken := CreateToken(u, expDT)
			w = SetResponseHeaders(w, sessionToken, "")
			w.WriteHeader(http.StatusOK)
			u.Password = ""
			if err := json.NewEncoder(w).Encode(u); err != nil {
				panic(err)
			}
		}
	}
}

// Handler function that refreshes a users JWT token
func (ur *userRouter) RefreshSession(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData := DecodeJWT(authToken)
	user := ur.userService.RefreshToken(tokenData)
	expDT := time.Now().Add(time.Hour * 1).Unix()
	newToken := CreateToken(user, expDT)
	w = SetResponseHeaders(w, newToken, "")
	w.WriteHeader(http.StatusOK)
}

// Handler function that generates a 6 month API Key for a given user
func (ur *userRouter) GenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	tokenData := DecodeJWT(authToken)
	user := ur.userService.RefreshToken(tokenData)
	expDT := time.Now().Add(time.Hour * 4380).Unix()
	apiKey := CreateToken(user, expDT)
	w = SetResponseHeaders(w, "", apiKey)
	w.WriteHeader(http.StatusOK)
}

// Handler function that ends a users session
func (ur *userRouter) Signout(w http.ResponseWriter, r *http.Request) {
	authToken := r.Header.Get("Auth-Token")
	ur.userService.BlacklistAuthToken(authToken)
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
}

// Handler function that registers a new user
func (ur *userRouter) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("REGISTRATION") == "off" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(404)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: 404, Text: "Not Found"}); err != nil {
			panic(err)
		}
	} else {
		var user core.User
		body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

		if err != nil {
			panic(err)
		}

		if err := r.Body.Close(); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(body, &user); err != nil {
			w = SetResponseHeaders(w, "", "")
			w.WriteHeader(422)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		} else {
			var group core.Group
			groupName := user.Username
			groupName += "_group"
			group.Name = groupName
			curid, err := uuid.NewV4()

			if err != nil {
				panic(err)
			}
			group.Id = curid.String()
			g := ur.groupService.GroupCreate(group)
			if g.Name == "Taken" {
				w = SetResponseHeaders(w, "", "")
				w.WriteHeader(403)
				if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
					panic(err)
				}
			} else {
				user.SystemRole = "group_admin"
				user.GroupIds = append(user.GroupIds, g.Id)
				u := ur.userService.UserCreate(user)
				if u.Email == "Taken" {
					w = SetResponseHeaders(w, "", "")
					w.WriteHeader(403)
					if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Email Taken"}); err != nil {
						panic(err)
					}
				} else if u.Username == "Taken" {
					w = SetResponseHeaders(w, "", "")
					w.WriteHeader(403)
					if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
						panic(err)
					}
				} else {
					w = SetResponseHeaders(w, "", "")
					w.WriteHeader(http.StatusCreated)
					u.Password = ""
					if err := json.NewEncoder(w).Encode(u); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}

// Handler function that creates a new user
func (ur *userRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user core.User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &user); err != nil {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(422)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"))
	groupUuid := AdminRouteRoleCheck(decodedToken)
	user.GroupIds = append(user.GroupIds, groupUuid...)
	u := ur.userService.UserCreate(user)
	if u.Email == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Email Taken"}); err != nil {
			panic(err)
		}
	} else if u.Username == "Taken" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(403)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Username Taken"}); err != nil {
			panic(err)
		}
	} else {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusCreated)
		u.Password = ""
		if err := json.NewEncoder(w).Encode(u); err != nil {
			panic(err)
		}
	}
}

// Handler that shows a specific user
func (ur *userRouter) UsersShow(w http.ResponseWriter, r *http.Request) {
	decodedToken := DecodeJWT(r.Header.Get("Auth-Token"))
	groupIds := AdminRouteRoleCheck(decodedToken)
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusOK)
	users := ur.userService.UsersFind(groupIds)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		panic(err)
	}
}

// Handler that shows all users
func (ur *userRouter) UserShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	user := ur.userService.UserFind(userId)
	if user.Id != "" {
		user.Password = ""
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			panic(err)
		}
		return
	}
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "User Not Found"}); err != nil {
		panic(err)
	}
}

// Handler function that deletes a user
func (ur *userRouter) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	user := ur.userService.UserDelete(userId)
	if user.Id != "" {
		w = SetResponseHeaders(w, "", "")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode("User Deleted"); err != nil {
			panic(err)
		}
		return
	}
	// If we didn't find it, 404
	w = SetResponseHeaders(w, "", "")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "User Not Found"}); err != nil {
		panic(err)
	}
}