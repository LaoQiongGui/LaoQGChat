package controllers

import (
	model "LaoQGChat/api/models/chat"
	"LaoQGChat/internal/myerrors"
	service "LaoQGChat/internal/services/chat"

	"github.com/gin-gonic/gin"
)

type ChatController interface {
	Chat(context *gin.Context)
	EndChat(context *gin.Context)
}

type chatController struct {
	service service.Service
}

func NewChatController(chatService service.Service) ChatController {
	return chatController{service: chatService}
}

func (c chatController) Chat(ctx *gin.Context) {
	var inDto model.Request
	err := ctx.Bind(&inDto)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "E000000",
			MessageText: "请求体格式错误。",
		}
		_ = ctx.Error(err)
		return
	}
	outDto := c.service.Chat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}

func (c chatController) EndChat(ctx *gin.Context) {
	var inDto model.Request
	err := ctx.Bind(&inDto)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "E000000",
			MessageText: "请求体格式错误。",
		}
		_ = ctx.Error(err)
		return
	}
	outDto := c.service.EndChat(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}
