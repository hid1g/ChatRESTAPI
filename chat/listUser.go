package chat

import (
	"sync"
)

type List struct {
	userById   map[int]User
	userByName map[string]int
	mtx        sync.RWMutex
}

func NewList() *List {
	return &List{
		userById:   make(map[int]User),
		userByName: make(map[string]int),
	}
}

func (l *List) UserExistByName(name string) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	_, ok := l.userByName[name]
	return ok
}

func (l *List) UserExist(id int) bool {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	_, ok := l.userById[id]
	return ok
}

func (l *List) CreateUser(user User) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, ok := l.userById[user.Id]; ok {
		return UserAlreadyExists
	}

	if _, ok := l.userByName[user.Name]; ok {
		return UserAlreadyExists
	}
	l.userById[user.Id] = user
	l.userByName[user.Name] = user.Id
	return nil
}

func (l *List) ListUsers() (map[int]User, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	if len(l.userById) < 1 {
		return nil, UsersNotFound
	}
	tmp := make(map[int]User)

	for k, v := range l.userById {
		tmp[k] = v
	}
	return tmp, nil
}

func (l *List) ListUserByName(name string) (User, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	id, ok := l.userByName[name]
	if !ok {
		return User{}, UserNotFoundError
	}
	user := l.userById[id]
	return user, nil
}

func (l *List) DeleteUser(id int) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	user, ok := l.userById[id]
	if !ok {
		return UserNotFoundError
	}

	delete(l.userByName, user.Name)
	delete(l.userById, id)
	return nil
}

func (l *List) UpdateUser(id int, username string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	user, ok := l.userById[id]
	if !ok {
		return UserNotFoundError
	}
	if _, exist := l.userByName[username]; exist {
		return UserAlreadyExists
	}

	delete(l.userByName, user.Name)
	user.Name = username
	l.userById[id] = user
	l.userByName[username] = id

	return nil

}
