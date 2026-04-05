package chat

import (
	"sync"
)

type ListMessage struct {
	message map[int]Message
	users   *List
	mtx     sync.RWMutex
	idsq    int
}

func NewListMessage(users *List) *ListMessage {
	return &ListMessage{
		message: make(map[int]Message),
		users:   users,
	}
}

func (l *ListMessage) SendMessage(message Message) (Message, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if message.Text == "" {
		return Message{}, MessageIsEmpty
	}
	if !l.users.UserExistByName(message.SendedFrom) {
		return Message{}, UserNotFoundEroor
	}
	if !l.users.UserExistByName(message.SendedTo) {
		return Message{}, UserNotFoundEroor
	}
	l.idsq++
	id := l.idsq
	message.Id = id
	message.Sended()
	l.message[id] = message
	return message, nil
}

func (l *ListMessage) GetMessagesByUser(username string) (map[int]Message, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	tmp := make(map[int]Message)
	if !l.users.UserExistByName(username) {
		return nil, UserNotFoundEroor
	}

	for id, msg := range l.message {
		if msg.SendedFrom == username {
			tmp[id] = msg
		}
	}
	return tmp, nil

}

func (l *ListMessage) DeleteMessage(id int) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	if _, ok := l.message[id]; !ok {
		return MessageNotFound
	}
	delete(l.message, id)
	return nil
}

func (l *ListMessage) MessageIsRead() {

}

func (l *ListMessage) MesaageIsNotRead() {

}
