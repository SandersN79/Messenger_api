package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
)



type users struct {
	Id          string        `bson:"_id,omitempty"`
	UUserKey    string        `bson:"UUserKey,omitempty"`
	Username    string        `bson:"Username,omitempty"`
	SystemRole  string        `bson:"SystemRole,omitempty"`
	GroupIds    []string    `bson:"GroupIds,omitempty"`
	GroupRoles  []string    `bson:"GroupRoles,omitempty"`
	Email       string        `bson:"Email,omitempty"`
	Password    string        `bson:"Password,omitempty"`
	Created     string        `bson:"Created,omitempty"`
}


func newUserModel(u core.User) *users {
	return &users{
		Id:              u.Id,
		UUserKey:        u.UUserKey,
		Username:        u.Username,
		SystemRole:      u.SystemRole,
		GroupRoles:      u.GroupRoles,
		GroupIds:        u.GroupIds,
		Email:           u.Email,
		Password:        u.Password,
		Created:         u.Created,
	}
}

func (u *users) toRootUser() core.User {
	return core.User{
		Id:              u.Id,
		UUserKey:        u.UUserKey,
		Username:        u.Username,
		SystemRole:      u.SystemRole,
		GroupRoles:      u.GroupRoles,
		GroupIds:        u.GroupIds,
		Email:           u.Email,
		Password:        u.Password,
		Created:         u.Created,
	}
}
