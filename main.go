package main

import (
	"chat/chat"
	"chat/database/connection"
	databaseCmd "chat/database/databaseCMD/createDB"
	"chat/httpHandlers"
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	conn, err := connection.CreateConnection(ctx)
	if err != nil {
		panic(err)
	}
	if err := databaseCmd.CreateMessageDB(ctx, conn); err != nil {
		panic(err)
	}
	if err := databaseCmd.CreateUserDB(ctx, conn); err != nil {
		panic(err)
	}

	chatApi := chat.NewList()
	messageApi := chat.NewListMessage(chatApi)
	httphandler := httphandlers.NewHttpHandlers(messageApi, chatApi, conn, ctx)
	httpServer := httphandlers.NewHttpServer(httphandler)
	if err := httpServer.StarServer(); err != nil {
		fmt.Println("Server Start Error")
	}
}
