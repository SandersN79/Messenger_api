package mongodb

import "MessengerDemo/server"

type messages struct {
	Id            string        `bson:"_id,omitempty"`
	SenderId      string        `bson:"SenderId,omitempty"`
	ReceiverId    string        `bson:"ReceiverId,omitempty"`
	Contents      string        `bson:"Contents,omitempty"`  //[]byte
	Created       string        `bson:"Created,omitempty"`
}

func newMessageModel(u server.Message) *messages {
	return &messages{
		Id:              u.Id,
		SenderId:        u.SenderId,
		ReceiverId:      u.ReceiverId,
		Contents:        u.Contents,
		Created:         u.Created,
	}
}

func (u *messages) toRootMessage() server.Message {
	return server.Message{
		Id:              u.Id,
		SenderId:        u.SenderId,
		ReceiverId:      u.ReceiverId,
		Contents:        u.Contents,
		Created:         u.Created,
	}
}