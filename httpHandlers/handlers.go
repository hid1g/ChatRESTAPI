package httphandlers

import (
	"chat/chat"
	messagecmd "chat/database/databaseCMD/messagecmd"
	databaseCMD "chat/database/databaseCMD/usercmd"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

type HtttpHandlers struct {
	Chat *chat.ListMessage
	User *chat.List
	conn *pgx.Conn
	ctx  context.Context
}

func NewHttpHandlers(ch *chat.ListMessage, u *chat.List, c *pgx.Conn, cont context.Context) *HtttpHandlers {
	return &HtttpHandlers{
		Chat: ch,
		User: u,
		conn: c,
		ctx:  cont,
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
	if err := databaseCMD.InsertUser(h.ctx, h.conn, newUser); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
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
	users, err := databaseCMD.ListUsers(h.ctx, h.conn)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
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
	user, err := databaseCMD.ListUsersByName(h.ctx, h.conn, username)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
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
body: json represent deleted message

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
	if err := databaseCMD.DeleteUser(h.ctx, h.conn, id); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
pattern:/chat/{id}
method: PUT
info: json in request body

succed:
status code: 200
body: json represent updated message

failed:
status code: 404, 500
body: json with err + time
*/
func (h *HtttpHandlers) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userDTO UserDTO
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	newUser := chat.NewUser(userDTO.Name)
	newUser.Id = id
	if err := databaseCMD.UpdateUser(h.ctx, h.conn, newUser); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

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
	if err := messagecmd.SendMessage(h.ctx, h.conn, newMessage); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return

	}

	b, err := json.MarshalIndent(newMessage, "", "    ")
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
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	user, err := messagecmd.GetMessageByUser(h.ctx, h.conn, id)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		panic(err)
	}
	w.Header().Set("content-type", "application/json")
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
status code: 404
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

	if err := messagecmd.DeleteMessage(h.ctx, h.conn, id); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/*
pattern: /chat/message/{id}
method PATCH
info pattern + method

succeed:
status code: 200
body: -

failed:
Status code: 404
*/
func (h *HtttpHandlers) MessageIsReadHandler(w http.ResponseWriter, r *http.Request) {
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}

	if err := messagecmd.MessageIsRead(h.ctx, h.conn, id); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
pattern: /chat/messages/{id}
method: PUT
info: pattern + json in request body

succed:
status code: 200
body: json

failed:
status code: 404, 500
*/

func (h *HtttpHandlers) MessageUpdateHandler(w http.ResponseWriter, r *http.Request) {
	var messageDTO MessageDTO
	if err := json.NewDecoder(r.Body).Decode(&messageDTO); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}

	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}

	if err := messagecmd.MessageUpdate(h.ctx, h.conn, id, messageDTO.Text); err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

/*
pattern:= /chat/messages/{user1}/{user2}
method:GET
info: pattern

succed:
status code: 200
body; json

failed:
status code: 400, 404, 500 ...
body: json with error + time
*/

func (h *HtttpHandlers) GetMessagesBetweenUsersHandler(w http.ResponseWriter, r *http.Request) {
	user1str := mux.Vars(r)["user1"]
	user1, err := strconv.Atoi(user1str)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	user2str := mux.Vars(r)["user2"]
	user2, err := strconv.Atoi(user2str)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusBadRequest)
		return
	}
	mes, err := messagecmd.GetMessagesBetweenUsers(h.ctx, h.conn, user1, user2)
	if err != nil {
		errDTO := NewErrDTO(err)
		http.Error(w, errDTO.ErrToString(), http.StatusInternalServerError)
		return
	}
	b, err := json.MarshalIndent(mes, "", "    ")
	if err != nil {
		panic(err)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("Writing error")
		return
	}

}
