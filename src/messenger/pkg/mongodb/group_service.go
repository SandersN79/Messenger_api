package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
	"go.mongodb.org/mongo-driver/bson"
)



// GroupService is used by the app to manage all group related controllers and functionality
type GroupService struct {
	db          *DBService
}

// NewGroupService is an exported function used to initialize a new GroupService struct
func NewGroupService(db *DBService) *GroupService {
	return &GroupService{db}
}

// GroupCreate is used to create a new user group
func (p *GroupService) GroupCreate(group core.Group) core.Group {
	return p.db.CreateGroup(group)
}

// GroupsFind is used to find all group docs
func (p *GroupService) GroupsFind() []core.Group {
	return p.db.FindGroups(bson.D{{}})
}

// GroupFind is used to find a specific group doc
func (p *GroupService) GroupFind(group core.Group) core.Group {
	return p.db.FindOneGroup(bson.D{{"Id", group.Id}})
}

// GroupDelete is used to delete a group doc
func (p *GroupService) GroupDelete(Id string) core.Group {
	return p.db.DeleteGroup(bson.D{{"Id", Id}})
}

// GroupUpdate is used to update an existing group
func (p *GroupService) GroupUpdate(group core.Group) core.Group {
	curGroup := p.db.FindOneGroup(bson.D{{"Id", group.Id}})
	filter := bson.D{{"Id", curGroup.Id}}
	//currentTime := time.Now().UTC()
	update := bson.D{{"$set",
		bson.D{
			{"Name", group.Name},
			//{"last_modified", currentTime.String()},
		},
	}}
	return p.db.UpdateGroup(update, filter)
}

// GroupDocInsert is used to insert a group doc directly into mongodb for testing purposes
func (p *GroupService) GroupDocInsert(group core.Group) core.Group {
	return p.db.InsertGroup(group)
}

// GroupInitialize
func (p *GroupService) GroupInitialize(adminGrp core.Group, kernelGrp core.Group, bootGrp core.Group) (core.Group, core.Group, core.Group) {
	return p.db.CreateGroup(adminGrp), p.db.CreateGroup(kernelGrp), p.db.CreateGroup(bootGrp)
}