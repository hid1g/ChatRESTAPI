package main

import (
	"chat/chat"
	"chat/httpHandlers"
	"fmt"
)

func main() {
	chatApi := chat.NewList()
	messageApi := chat.NewListMessage(chatApi)
	httphandler := httphandlers.NewHttpHandlers(messageApi, chatApi)
	httpServer := httphandlers.NewHttpServer(httphandler)
	if err := httpServer.StarServer(); err != nil {
		fmt.Println("Server Start Error")
	}
}
