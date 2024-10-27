package middlewares

import (
	"LaoQGChat/api/models"
	"LaoQGChat/internal/myerrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthHandler(checkFunc func(loginToken uuid.UUID) (*models.AuthDto, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 前处理
		// 认证除外
		if ctx.Request.RequestURI != "/Auth/Login" {
			var (
				err        error
				loginToken uuid.UUID
				authDto    *models.AuthDto
			)

			// 获取loginToken
			loginToken, err = uuid.Parse(ctx.GetHeader("LoginToken"))
			if err != nil {
				err = &myerrors.CustomError{
					StatusCode:  200,
					MessageCode: "EAU01",
					MessageText: "用户未登录。",
				}
				panic(err)
			}

			// 验证登陆状态
			authDto, err = checkFunc(loginToken)
			if err != nil {
				panic(err)
			}

			// 设置用户信息
			ctx.Set("UserName", authDto.Username)
			ctx.Set("Permission", authDto.Permission)
		}

		// 下一层
		ctx.Next()

		// 后处理
	}
}
