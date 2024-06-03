package main

import (
	"LaoQGChat/controller"
	"LaoQGChat/handler"
	"LaoQGChat/service"
	"github.com/gin-gonic/gin"
)

var (
	chatService    = service.NewChatService()
	chatController = controller.NewChatController(chatService)
)

func main() {
	server := gin.Default()

	server.POST("/Chat/StartChat", handler.HandlerBuilder(chatController.StartChat))

	server.POST("/Chat/Chat", handler.HandlerBuilder(chatController.Chat))

	server.POST("/Chat/EndChat", handler.HandlerBuilder(chatController.EndChat))

	server.Run(":12195")
}
