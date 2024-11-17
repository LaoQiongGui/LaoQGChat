package middlewares

import (
	"LaoQGChat/api/models/auth"
	"LaoQGChat/internal/myerrors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func AuthHandler(checkFunc func(loginToken uuid.UUID) (*auth.UserInfo, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 前处理
		// 认证除外
		if ctx.Request.RequestURI != "/Auth/Login" {
			var (
				err        error
				loginToken uuid.UUID
				authDto    *auth.UserInfo
			)

			// 获取loginToken
			loginToken, err = uuid.Parse(ctx.GetHeader("LoginToken"))
			if err != nil {
				err = &myerrors.CustomError{
					StatusCode:  200,
					MessageCode: "EAU01",
					MessageText: "用户未登录。",
				}
				_ = ctx.AbortWithError(http.StatusNonAuthoritativeInfo, err)
				return
			}

			// 验证登陆状态
			authDto, err = checkFunc(loginToken)
			if err != nil {
				err = &myerrors.CustomError{
					StatusCode:  200,
					MessageCode: "EAU01",
					MessageText: "用户未登录。",
				}
				_ = ctx.AbortWithError(http.StatusNonAuthoritativeInfo, err)
				return
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
