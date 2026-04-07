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
		return Message{}, UserNotFoundError
	}
	if !l.users.UserExistByName(message.SendedTo) {
		return Message{}, UserNotFoundError
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
		return nil, UserNotFoundError
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

func (l *ListMessage) MessageIsRead(id int) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	mes, ok := l.message[id]
	if !ok {
		return MessageNotFound
	}
	mes.Read()
	l.message[id] = mes
	return nil
}

func (l *ListMessage) MessageUpdate(id int, newmessage string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	mes, exist := l.message[id]
	if !exist {
		return MessageNotFound
	}

	mes.Text = newmessage
	l.message[id] = mes
	return nil
}

func (l *ListMessage) GetMeassagesBetweenUsers(user1 string, user2 string) (map[int]Message, error) {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	tmp := make(map[int]Message)
	if !l.users.UserExistByName(user1) || !l.users.UserExistByName(user2) {
		return nil, UserNotFoundError
	}
	for id, msg := range l.message {
		if (msg.SendedFrom == user1 && msg.SendedTo == user2) || (msg.SendedFrom == user2 && msg.SendedTo == user1) {
			tmp[id] = msg
		}
	}
	return tmp, nil
}
