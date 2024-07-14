package main

import (
	"LaoQGChat/controller"
	"LaoQGChat/handler"
	"LaoQGChat/service"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"time"
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

	// 初始化db
	connStr := "host=localhost port=5432 user=laoqionggui password=LaoQi0ng@ui sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("DB连接失败")
		return
	}

	// 初始化service
	var (
		authService    = service.NewAuthService(db)
		authController = controller.NewAuthController(authService)

		chatService    = service.NewChatService(db)
		chatController = controller.NewChatController(chatService)
	)
	if authService == nil || authController == nil || chatService == nil || chatController == nil {
		fmt.Println("初始化service失败")
		return
	}

	server.Use(cors.New(config))

	server.POST("/Auth/Login", handler.HandlerBuilder(authController.Login))

	server.POST("/Chat/StartChat", handler.HandlerBuilder(chatController.StartChat))

	server.POST("/Chat/Chat", handler.HandlerBuilder(chatController.Chat))

	server.POST("/Chat/EndChat", handler.HandlerBuilder(chatController.EndChat))

	err = server.Run(":12195")
	if err != nil {
		fmt.Println("启动服务失败")
		return
	}
}
