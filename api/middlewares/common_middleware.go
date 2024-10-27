package middlewares

import (
	"LaoQGChat/api/models"
	"LaoQGChat/internal/myerrors"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 前处理
		defer func() {
			response := models.Response{}

			if err := recover(); err != nil {
				// 处理 panic
				response.Common = models.ResponseCommonSystemError
			} else {
				// 处理错误信息
				response.Common = makeResponseCommon(ctx)
			}

			// 填充响应体Data部
			response.Data = makeResponseData(ctx)

			// 设置响应
			ctx.JSON(http.StatusOK, response)
		}()

		// 下一层
		ctx.Next()

		// 后处理
	}
}

func makeResponseCommon(ctx *gin.Context) models.ResponseCommon {
	if err := ctx.Errors.Last(); err != nil {
		// 处理自定义异常
		var myError *myerrors.CustomError
		if errors.As(err.Err, &myError) {
			return models.ResponseCommon{
				Status:      myError.StatusCode,
				MessageCode: myError.MessageCode,
				MessageText: myError.MessageText,
			}
		}
		// 处理其他异常
		return models.ResponseCommonSystemError
	} else {
		// 正常返回
		return models.ResponseCommonSuccess
	}
}

func makeResponseData(ctx *gin.Context) interface{} {
	if data, exists := ctx.Get("ResponseData"); exists {
		return data
	} else {
		return nil
	}
}
