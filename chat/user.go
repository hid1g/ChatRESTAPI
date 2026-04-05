package chat

import (
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

