package controller

import (
	"LaoQGChat/dto"
	"LaoQGChat/service"
	"github.com/gin-gonic/gin"
)

type ChatController interface {
	StartChat(context *gin.Context) *dto.ChatOutDto
	Chat(context *gin.Context) *dto.ChatOutDto
	EndChat(context *gin.Context) *dto.ChatOutDto
}

type chatController struct {
	service service.ChatService
}

func NewChatController(chatService service.ChatService) ChatController {
	return chatController{
		service: chatService,
	}
}

func (c chatController) StartChat(ctx *gin.Context) *dto.ChatOutDto {
	var inDto dto.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	return c.service.StartChat(ctx, inDto)
}

func (c chatController) Chat(ctx *gin.Context) *dto.ChatOutDto {
	var inDto dto.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	return c.service.Chat(ctx, inDto)
}

func (c chatController) EndChat(ctx *gin.Context) *dto.ChatOutDto {
	var inDto dto.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	return c.service.EndChat(ctx, inDto)
}
