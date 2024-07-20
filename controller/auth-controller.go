package controller

import (
	"LaoQGChat/dto"
	"LaoQGChat/service"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
}

type authController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	controller := new(authController)
	controller.service = service
	return controller
}

func (c *authController) Login(ctx *gin.Context) {
	inDto := dto.AuthDto{}
	err := ctx.Bind(&inDto)
	if err != nil {
		ctx.Keys["StatusCode"] = 200
		ctx.Keys["MessageCode"] = "E0000"
		ctx.Keys["MessageText"] = "请求体格式错误。"
		panic(err)
	}
	outDto := c.service.Login(ctx, inDto)
	ctx.Set("ResponseData", outDto)
}
