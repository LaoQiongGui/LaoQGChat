package controllers

import (
	"LaoQGChat/api/models"
	"LaoQGChat/api/services"
	"LaoQGChat/internal/myerrors"
	"github.com/gin-gonic/gin"
)

type ChatController interface {
	StartChat(context *gin.Context)
	Chat(context *gin.Context)
	EndChat(context *gin.Context)
}

type chatController struct {
	authService services.AuthService
	service     services.ChatService
}

func NewChatController(authService services.AuthService, chatService services.ChatService) ChatController {
	return chatController{
		authService: authService,
		service:     chatService,
	}
}

func (c chatController) StartChat(ctx *gin.Context) {
	var inDto models.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "E0000",
			MessageText: "请求体格式错误。",
		}
		_ = ctx.Error(err)
		return
	}
	outDto := c.service.StartChat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}

func (c chatController) Chat(ctx *gin.Context) {
	var inDto models.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "E0000",
			MessageText: "请求体格式错误。",
		}
		_ = ctx.Error(err)
		return
	}
	outDto := c.service.Chat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}

func (c chatController) EndChat(ctx *gin.Context) {
	var inDto models.ChatInDto
	err := ctx.Bind(&inDto)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "E0000",
			MessageText: "请求体格式错误。",
		}
		_ = ctx.Error(err)
		return
	}
	outDto := c.service.EndChat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}
