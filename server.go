package main

import (
	"LaoQGChat/api/controllers"
	"LaoQGChat/api/middlewares"
	"LaoQGChat/api/services"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// 初始化db
	db, err := initDB()
	if err != nil {
		fmt.Println("DB连接失败")
		return
	}

	server := gin.Default()

	// 配置CORS中间件
	config := cors.Config{
		AllowAllOrigins:  true,                               // 允许所有的域名
		AllowMethods:     []string{"POST"},                   // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},         // 暴露的头信息
		AllowCredentials: true,                               // 允许携带凭证
		MaxAge:           12 * time.Hour,                     // 预检请求缓存时间
	}
	server.Use(cors.New(config))

	// 配置异常处理中间件
	server.Use(middlewares.CommonErrorHandler())

	// 初始化认证service
	var (
		authService    = services.NewAuthService(db)
		authController = controllers.NewAuthController(authService)
	)
	if authService == nil || authController == nil {
		fmt.Println("初始化认证service失败")
		return
	}

	// 配置版本检测中间件
	server.Use(middlewares.VersionHandler("1.2.0"))

	// 配置DB事务中间件
	server.Use(middlewares.TransactionHandler(db))

	// 配置认证中间件
	server.Use(middlewares.AuthHandler(authService.Check))

	// 初始化业务service
	var (
		chatService    = services.NewChatService(db)
		chatController = controllers.NewChatController(authService, chatService)
	)
	if chatService == nil || chatController == nil {
		fmt.Println("初始化业务service失败")
		return
	}

	server.POST("/Auth/Login", authController.Login)

	server.POST("/Chat/StartChat", chatController.StartChat)

	server.POST("/Chat/Chat", chatController.Chat)

	server.POST("/Chat/EndChat", chatController.EndChat)

	err = server.Run(":12195")
	if err != nil {
		fmt.Println("启动服务失败")
		return
	}
}

func initDB() (*sql.DB, error) {
	connStr := "host=localhost port=5432 user=laoqionggui password=LaoQi0ng@ui sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// 设置连接池
	db.SetMaxOpenConns(10)   // 最大打开连接数
	db.SetMaxIdleConns(10)   // 最大闲置连接数
	db.SetConnMaxLifetime(0) // 连接的最大存活时间

	return db, nil
}
