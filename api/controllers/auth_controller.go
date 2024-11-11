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
	controller := new(authController)
	controller.service = service
	return controller
}

func (c *authController) Login(ctx *gin.Context) {
	inDto := model.Entity{}
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
	outDto := c.service.Login(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}
