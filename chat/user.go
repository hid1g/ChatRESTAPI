package chat

import (
	"sync"
	"time"
)

type User struct {
	Name      string
	Id        int
	CreatedAt time.Time
}

var idCounter int

func NewUser(name string) User {
	idCounter++
	return User{
		Name:      name,
		Id:        idCounter,
		CreatedAt: time.Now(),
	}
}

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
