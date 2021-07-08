package core

type Message struct {
	Id            string        `json:"id,omitempty"`
	SenderId      string        `json:"SenderId,omitempty"`
	ReceiverId    string        `json:"ReceiverId,omitempty"`//
	Contents      []byte        `json:"Contents,omitempty"`  //[]byte
	Created       string        `json:"Created,omitempty"`
	LastModified  string        `json:"LastModified,omitempty"`
	Deleted       []bool        `json:"Deleted,omitempty"`   // [false, true] - true is deleted, false is not
	Read          string        `json:"Read,omitempty"`
	ReadTime      string        `json:"ReadTime,omitempty"`
	Edited        string        `json:"Edited,omitempty"`
	EditTime      string        `json:"EditTime,omitempty"`
}

// MessageService is an interface used to manage the relevant user doc controllers
type MessageService interface {
	MessageCreate(u Message) Message
	MessageDelete(u Message) Message
	MessagesFind(u Message) []Message
	MessageFind(u Message) Message
	MessageUpdate(u Message) Message
	MessageDocInsert(u Message) Message
}
