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
	authService service.AuthService
	service     service.ChatService
}

func NewChatController(authService service.AuthService, chatService service.ChatService) ChatController {
	return chatController{
		authService: authService,
		service:     chatService,
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
	// 认证check
	c.authService.Check(ctx, func(permission string) {
		inDto.Permission = permission
	})
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
	// 认证check
	c.authService.Check(ctx, func(permission string) {
		inDto.Permission = permission
	})
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
	// 认证check
	c.authService.Check(ctx, func(permission string) {
		inDto.Permission = permission
	})
	return c.service.EndChat(ctx, inDto)
}
