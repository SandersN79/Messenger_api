package server



type Group struct {
	Id          string        `json:"id,omitempty"`
	Name        string        `json:"Name,omitempty"`
	Type        string        `json:"Type,omitempty"`
	Created     string        `json:"Created,omitempty"`
}

type Message struct {
	Id            string        `json:"id,omitempty"`
	SenderId      string        `json:"SenderId,omitempty"`
	ReceiverId    string        `json:"ReceiverId,omitempty"`
	Contents      string        `json:"Contents,omitempty"`  //[]byte
	Created       string        `json:"Created,omitempty"`
}

type User struct {
	Id          string        `json:"id,omitempty"`
	UUserKey    string        `json:"UUserKey,omitempty"`
	Username    string        `json:"Username,omitempty"`
	SystemRole  string        `json:"SystemRole,omitempty"`
	//GroupIds    []string    `json:"GroupIds,omitempty"`
	//GroupRoles  []string    `json:"GroupRoles,omitempty"`
	Email       string        `json:"Email,omitempty"`
	Password    string        `json:"Password,omitempty"`
	Created     string        `json:"Created,omitempty"`
}


// UserService is an interface used to manage the relevant user doc controllers
type UserService interface {
	AuthenticateUser(u User) User
	//BlacklistAuthToken(authToken string)
	//RefreshToken(tokenData []string, groupUuid string) User
	UpdatePassword(tokenData []string, CurrentPassword string, newPassword string) User
	UserCreate(u User) User
	UserDelete(id string) User
	UsersFind(groupUuid string) []User
	UserFind(id string) User
	UserUpdate(u User) User
	UserDocInsert(u User) User
}
