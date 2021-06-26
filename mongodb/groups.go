package mongodb

import "MessengerDemo/server"

type groups struct {
	Id          string        `bson:"_id,omitempty"`
	Name        string        `bson:"Name,omitempty"`
	Type        string        `bson:"Type,omitempty"`
	Created     string        `bson:"Created,omitempty"`
}

func newGroupModel(u server.Group) *groups {
	return &groups{
		Id:        u.Id,
		Name:      u.Name,
		Type:      u.Type,
		Created:   u.Created,
	}
}

func (u *groups) toRootGroup() server.Group {
	return server.Group{
		Id:        u.Id,
		Name:      u.Name,
		Type:      u.Type,
		Created:   u.Created,
	}
}
