package httphandlers

import (
	"chat/chat"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type HtttpHandlers struct {
	Chat *chat.ListMessage
	User *chat.List
}

func NewHttpHandlers(ch *chat.ListMessage, u *chat.List) *HtttpHandlers {
	return &HtttpHandlers{
		Chat: ch,
		User: u,
	}
}

func NewErrDTO(err error) *ErrDTO {
	return &ErrDTO{
		Error:     err.Error(),
		ErrorTime: time.Now(),
	}
}

/*
pattern: /chat
method: POST
info: json in request body

succeed:
status code: 201 created
body: json represent created block

failed:
status code: 400, 500 ...
body: json with error + time
*/
func (h *HtttpHandlers) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}

	if err := userDTO.ValidateToCreateUser(); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	newUser := chat.NewUser(userDTO.Name)
	if err := h.User.CreateUser(newUser); err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.UserAlreadyExists) {
			http.Error(w, errDTO.ErrToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(newUser, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Writing error")
		return
	}
}

/*
patter:/chat
method:GET
info: pattern

succed:
status code: 200
body: json users

failed:
status code: 404, 500 ...
body: json with error + time
*/
func (h *HtttpHandlers) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.User.ListUsers()
	if err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.UsersNotFound) {
			http.Error(w, errDTO.ErrToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Writing error")
		return
	}
}

/*
patter:/chat/{name}
method: GET
info: json user

succed:
status code: 200
body: json represent user

failed:
status code: 400, 404, 500
body: json with error + time


*/

func (h *HtttpHandlers) ListUserByNameHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["name"]
	user, err := h.User.ListUserByName(username)
	if err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.UserNotFoundError) {
			http.Error(w, errDTO.ErrToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Writing error")
		return
	}
}

/*
patter:/chat/{name}
method: DELETE
info: pattern

succed:
status code: 204
body: json represent deleted book

failed:
status code: 404, 500 ...
body: json with error + time

*/

func (h *HtttpHandlers) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	if err := h.User.DeleteUser(id); err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.UserNotFoundError) {
			http.Error(w, errDTO.ErrToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
patter:/chat/message/{name}
method:POST
info:pattern + json in request body

succed:
status code: 201
body: json message

failed:
status code: 400, 500 ...
body: json with error + time

*/

func (h *HtttpHandlers) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var messageDTO MessageDTO
	if err := json.NewDecoder(r.Body).Decode(&messageDTO); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}

	if err := messageDTO.ValidateToCreateMessage(); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	newMessage := chat.NewMessage(messageDTO.Text, messageDTO.Sender, messageDTO.Reciever)
	sendedMessage, err := h.Chat.SendMessage(newMessage)
	if err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.MessageIsEmpty) {
			http.Error(w, errDTO.ErrToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(sendedMessage, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Writing error")
		return
	}
}

/*
patter: /chat/message/{name}
method:GET
info:Pattern + json

succed:
status code: 200
body: json messages from user

failed:
status code: 404, 500 ...
body: jsom with error + time
*/
func (h *HtttpHandlers) GetMessagesByUserHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	user, err := h.Chat.GetMessagesByUser(name)
	if err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.UserNotFoundError) {
			http.Error(w, errDTO.ErrToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Writing error")
		return
	}
}

/*
patter:/chat/message/{id}
method: DELETE
info: pattern

succed:
status code: 204
body:

failed:
status code: json deleted book
body: json with error + time
*/
func (h *HtttpHandlers) DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}

	if err := h.Chat.DeleteMessage(id); err != nil {
		errDTO := NewErrDTO(err)
		if errors.Is(err, chat.MessageNotFound) {
			http.Error(w, errDTO.ErrToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
