package core

type Group struct {
	Id          string        `json:"id,omitempty"`
	Name        string        `json:"Name,omitempty"`
	Type        string        `json:"Type,omitempty"`
	Created     string        `json:"Created,omitempty"`
}

// GroupService is an interface used to manage the relevant user doc controllers
type GroupService interface {
	GroupCreate(u Group) Group
	GroupDelete(id string) Group
	GroupsFind() []Group
	GroupFind(u Group) Group
	GroupUpdate(u Group) Group
	GroupDocInsert(u Group) Group
}
