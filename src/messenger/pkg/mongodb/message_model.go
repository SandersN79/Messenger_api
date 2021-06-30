package mongodb

import core "MessengerDemo/src/messenger/pkg"

type messages struct {
	Id            string        `bson:"_id,omitempty"`
	SenderId      string        `bson:"SenderId,omitempty"`
	ReceiverId    string        `bson:"ReceiverId,omitempty"`
	Contents      []byte        `bson:"Contents,omitempty"`  //[]byte
	Created       string        `bson:"Created,omitempty"`
}

func newMessageModel(u core.Message) *messages {
	return &messages{
		Id:              u.Id,
		SenderId:        u.SenderId,
		ReceiverId:      u.ReceiverId,
		Contents:        u.Contents,
		Created:         u.Created,
	}
}

func (u *messages) toRootMessage() core.Message {
	return core.Message{
		Id:              u.Id,
		SenderId:        u.SenderId,
		ReceiverId:      u.ReceiverId,
		Contents:        u.Contents,
		Created:         u.Created,
	}
}