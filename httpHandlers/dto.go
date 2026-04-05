package httphandlers

import (
	"encoding/json"
	"errors"
	"time"
)

type UserDTO struct {
	Name string
}

type MessageDTO struct {
	Text     string
	Sender   string
	Reciever string
}

type ErrDTO struct {
	Error     string
	ErrorTime time.Time
}

type SendedDTO struct {
	Sended bool
}

func (u *UserDTO) ValidateToCreateUser() error {
	if u.Name == "" {
		return errors.New("Username is empty")
	}
	return nil
}

func (m *MessageDTO) ValidateToCreateMessage() error {
	if m.Text == "" {
		return errors.New("Message is empty")
	}
	if m.Sender == "" {
		return errors.New("No sender")
	}
	if m.Reciever == "" {
		return errors.New("No reciever")
	}
	return nil
}

func (e *ErrDTO) ErrToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
