package chat

import "time"

type Message struct {
	Text       string
	SendedFrom string
	SendedTo   string
	IsSended   bool
	SendedTime *time.Time
	Id         int
	IsRead     bool
}

func NewMessage(text string, sender string, reciever string) Message {
	return Message{
		Text:       text,
		SendedFrom: sender,
		SendedTo:   reciever,
		IsSended:   false,
		SendedTime: nil,
		IsRead:     false,
	}
}

func (m *Message) Sended() {
	sendedTime := time.Now()

	m.IsSended = true
	m.SendedTime = &sendedTime
}

func (m *Message) NotSended() {
	m.IsSended = false
	m.SendedTime = nil
}

func (m *Message) Read() {
	readtime := time.Now()
	m.IsRead = true
	m.SendedTime = &readtime
}
