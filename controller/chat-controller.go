package controller

import (
	"LaoQGChat/dto"
	"LaoQGChat/service"

	"github.com/gin-gonic/gin"
)

type ChatController interface {
	StartChat(context *gin.Context)
	Chat(context *gin.Context)
	EndChat(context *gin.Context)
}

type chatController struct {
	authService service.AuthService
	service     service.ChatService
}

func NewChatController(authService service.AuthService, chatService service.ChatService) ChatController {
	return chatController{
		authService: authService,
		service:     chatService,
	}
}

func (c chatController) StartChat(ctx *gin.Context) {
	var inDto dto.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	outDto := c.service.StartChat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}

func (c chatController) Chat(ctx *gin.Context) {
	var inDto dto.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	outDto := c.service.Chat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}

func (c chatController) EndChat(ctx *gin.Context) {
	var inDto dto.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	outDto := c.service.EndChat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}
