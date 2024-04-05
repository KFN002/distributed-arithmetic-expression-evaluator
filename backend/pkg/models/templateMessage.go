package models

type Message struct {
	Message string
}

func (m *Message) AddMessage(msg string) {
	m.Message = msg
}

func CreateNewMessage() *Message {
	return &Message{}
}
