package controllers

import (
	"LaoQGChat/api/models"
	"LaoQGChat/api/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
}

type authController struct {
	service services.AuthService
}

func NewAuthController(service services.AuthService) AuthController {
	controller := new(authController)
	controller.service = service
	return controller
}

func (c *authController) Login(ctx *gin.Context) {
	inDto := models.AuthDto{}
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
