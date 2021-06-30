package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type blacklistModel struct {
	Id               primitive.ObjectID `bson:"_id,omitempty"`
	AuthToken        string             `bson:"auth_token,omitempty"`
	Created          string             `bson:"creation_datetime,omitempty"`
}

func newBlacklistModel(bl core.Blacklist) *blacklistModel {
	return &blacklistModel{
		AuthToken:        bl.AuthToken,
		Created:          bl.Created,
	}
}

func (bl *blacklistModel) toRootBlacklist() core.Blacklist {
	return core.Blacklist{
		Id:               bl.Id.Hex(),
		AuthToken:        bl.AuthToken,
		Created:          bl.Created,
	}
}