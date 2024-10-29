package controllers

import (
	"LaoQGChat/api/models"
	"LaoQGChat/internal/myerrors"
	"LaoQGChat/internal/services"
	"fmt"
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
	fmt.Println("Calling controller Login")
	defer fmt.Println("Returning controller Login")
	inDto := models.AuthDto{}
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
