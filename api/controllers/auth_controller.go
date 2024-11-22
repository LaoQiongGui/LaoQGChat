package controllers

import (
	model "LaoQGChat/api/models/auth"
	"LaoQGChat/internal/myerrors"
	service "LaoQGChat/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
}

type authController struct {
	service service.Service
}

func NewAuthController(service service.Service) AuthController {
	return &authController{service: service}
}

func (c *authController) Login(ctx *gin.Context) {
	var request model.UserInfo
	err := ctx.Bind(&request)
	if err != nil {
		err = &myerrors.CustomError{
			StatusCode:  200,
			MessageCode: "E000000",
			MessageText: "请求体格式错误。",
		}
		_ = ctx.Error(err)
		return
	}
	response := c.service.Login(ctx, request)
	ctx.Set("ResponseData", response)
}
