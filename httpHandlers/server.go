package httphandlers

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type HttpServer struct {
	httphandlersForServer *HtttpHandlers
}

func NewHttpServer(httphandler *HtttpHandlers) *HttpServer {
	return &HttpServer{
		httphandlersForServer: httphandler,
	}
}

func (s *HttpServer) StarServer() error {
	router := mux.NewRouter()
	router.Path("/chat").Methods("POST").HandlerFunc(s.httphandlersForServer.CreateUserHandler)
	router.Path("/chat/{name}").Methods("GET").HandlerFunc(s.httphandlersForServer.ListUserByNameHandler)
	router.Path("/chat").Methods("GET").HandlerFunc(s.httphandlersForServer.ListUsersHandler)
	router.Path("/chat/{id}").Methods("DELETE").HandlerFunc(s.httphandlersForServer.DeleteUserHandler)
	router.Path("/chat/message/{name}").Methods("POST").HandlerFunc(s.httphandlersForServer.SendMessageHandler)
	router.Path("/chat/message/{name}").Methods("GET").HandlerFunc(s.httphandlersForServer.GetMessagesByUserHandler)
	router.Path("/chat/message/{id}").Methods("DELETE").HandlerFunc(s.httphandlersForServer.DeleteMessageHandler)
	router.Path("/chat/message/{id}").Methods("PATCH").HandlerFunc(s.httphandlersForServer.MessageIsReadHandler)

	if err := http.ListenAndServe(":9091", router); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
	return nil
}
