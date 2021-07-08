package mongodb

import core "MessengerDemo/src/messenger/pkg"

type messages struct {
	Id            string        `bson:"_id,omitempty"`
	SenderId      string        `bson:"SenderId,omitempty"`
	ReceiverId    string        `bson:"ReceiverId,omitempty"`
	Contents      []byte        `bson:"Contents,omitempty"`  //[]byte
	Created       string        `bson:"Created,omitempty"`
	LastModified  string        `bson:"LastModified,omitempty"`
	Deleted       []bool        `bson:"Deleted,omitempty"`   // [false, true] - true is deleted, false is not
	Read          string          `bson:"Read,omitempty"`
	ReadTime      string        `bson:"ReadTime,omitempty"`
	Edited        string        `bson:"Edited,omitempty"`
	EditTime      string        `bson:"EditTime,omitempty"`
}

func newMessageModel(u core.Message) *messages {
	return &messages{
		Id:              u.Id,
		SenderId:        u.SenderId,
		ReceiverId:      u.ReceiverId,
		Contents:        u.Contents,
		Created:         u.Created,
		LastModified:    u.LastModified,
		Deleted:         u.Deleted,
		Read:            u.Read,
		ReadTime:        u.ReadTime,
		Edited:          u.Edited,
		EditTime:        u.EditTime,
	}
}

func (u *messages) toRootMessage() core.Message {
	return core.Message{
		Id:              u.Id,
		SenderId:        u.SenderId,
		ReceiverId:      u.ReceiverId,
		Contents:        u.Contents,
		Created:         u.Created,
		LastModified:    u.LastModified,
		Deleted:         u.Deleted,
		Read:            u.Read,
		ReadTime:        u.ReadTime,
		Edited:          u.Edited,
		EditTime:        u.EditTime,
	}
}