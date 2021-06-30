package core

type Message struct {
	Id            string        `json:"id,omitempty"`
	SenderId      string        `json:"SenderId,omitempty"`
	ReceiverId    string        `json:"ReceiverId,omitempty"`
	Contents      []byte        `json:"Contents,omitempty"`  //[]byte
	Created       string        `json:"Created,omitempty"`
}

// MessageService is an interface used to manage the relevant user doc controllers
type MessageService interface {
	MessageCreate(u Message) Message
	MessageDelete(id string) Message
	MessagesFind(groupUuid string) []Message
	MessageFind(id string) Message
	MessageUpdate(u Message) Message
	MessageDocInsert(u Message) Message
}
