package main

import (
	"LaoQGChat/controller"
	"LaoQGChat/handler"
	"LaoQGChat/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	authService    = service.NewAuthService()
	authController = controller.NewAuthController(authService)

	chatService    = service.NewChatService()
	chatController = controller.NewChatController(chatService)
)

func main() {
	server := gin.Default()

	// 配置CORS中间件
	config := cors.Config{
		AllowAllOrigins:  true,                                     // 允许所有的域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type"},       // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},               // 暴露的头信息
		AllowCredentials: true,                                     // 允许携带凭证
		MaxAge:           12 * time.Hour,                           // 预检请求缓存时间
	}
	server.Use(cors.New(config))

	server.POST("/Auth/Login", handler.HandlerBuilder(authController.Login))

	server.POST("/Chat/StartChat", handler.HandlerBuilder(chatController.StartChat))

	server.POST("/Chat/Chat", handler.HandlerBuilder(chatController.Chat))

	server.POST("/Chat/EndChat", handler.HandlerBuilder(chatController.EndChat))

	server.Run(":12195")
}
