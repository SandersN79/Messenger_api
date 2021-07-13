package mongodb

import (
	core "MessengerDemo/src/messenger/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type MessageService struct {
	db          *DBService
}

// NewMessageService is an exported function used to initialize a new MessageService struct
func NewMessageService(db *DBService) *MessageService {
	return &MessageService{db}
}

//MessageCreate is used to create a message post
func (p *MessageService) MessageCreate(message core.Message) core.Message {
	//message.Contents =
	message.Deleted = []bool{false, false}
	return p.db.CreateMessage(message)
}

//MessagesFind is used to find messages
func (p *MessageService) MessagesFind(message core.Message) []core.Message {
	if message.SenderId != "" {
		return p.db.FindMessages(bson.D{{"SenderId", message.SenderId}})
	} else if message.ReceiverId != "" {
		return p.db.FindMessages(bson.D{{"ReceiverId", message.ReceiverId}})
	}
	// return p.db.FindMessages(bson.D{{}})      ----When logic for group messages is introduced, might reintroduce
	//new message find filter cases
	return []core.Message{}
}

//MessageFind is used to find message
func (p *MessageService) MessageFind(message core.Message) core.Message {
	return p.db.FindOneMessage(bson.D{{"Id", message.Id}})
}

//MessageUpdate is used to update a message
func (p *MessageService) MessageUpdate(message core.Message) core.Message {
	curMessage := p.db.FindOneMessage(bson.D{{"Id", message.Id}})
	filter := bson.D{{"Id", curMessage.Id}}
	currentTime := time.Now().UTC()
	if len(message.Contents) == 0 {message.Contents = curMessage.Contents}
	if len(message.Edited) == 0 {message.Edited = curMessage.Edited}
	if len(message.EditTime) == 0 {message.EditTime = curMessage.EditTime}
	if len(message.Read) == 0 {message.Read = curMessage.Read}
	if len(message.ReadTime) == 0 {message.ReadTime = curMessage.ReadTime}
	if len(message.Deleted) == 0 {message.Deleted = curMessage.Deleted}
	update := bson.D{{"$set",
		bson.D{
			{"Contents", message.Contents},
			{"Edited", message.Edited},
			{"EditTime", message.EditTime},
			{"Read", message.Read},
			{"ReadTime", message.ReadTime},
			{"Deleted", message.Deleted},
			{"LastModified", currentTime.String()},
		},
	}}
	return p.db.UpdateMessage(update, filter)
}

//MessageDelete is used to delete a message
func (p *MessageService) MessageDelete(message core.Message) core.Message {
	curMessage := p.db.FindOneMessage(bson.D{{"Id", message.Id}})
	if message.SenderId != "" {
		if !curMessage.Deleted[0] {
			curMessage.Deleted[0] = true
			return p.MessageUpdate(message)
		}
	} else if message.ReceiverId != "" {
		if !curMessage.Deleted[1] {
			curMessage.Deleted[1] = true
			return p.MessageUpdate(message)
		}
	}
	return core.Message{}
}


//MessageDocInsert is used to insert a message doc directly into mongodb for testing purposes
func (p *MessageService) MessageDocInsert(message core.Message) core.Message {
	return p.db.InsertMessage(message)
}
