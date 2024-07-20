package handler

import (
	"LaoQGChat/dto"
	"LaoQGChat/myerror"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func CommonErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 前处理
		defer errorHandler(ctx)

		// 下一层
		ctx.Next()

		// 后处理
	}
}

func errorHandler(ctx *gin.Context) {
	restOutDto := dto.RestOutDto{}

	// 填充响应体Common部
	if err := recover(); err != nil {
		switch myError := err.(type) {
		case *myerror.CustomError:
			restOutDto.Common = dto.RestCommonDto{
				Status:      myError.StatusCode,
				MessageCode: myError.MessageCode,
				MessageText: myError.MessageText,
			}
		default:
			restOutDto.Common = dto.RestCommonDto{
				Status:      990,
				MessageCode: "E9999",
				MessageText: "System Error",
			}
		}
	} else {
		restOutDto.Common = dto.RestCommonDto{
			Status:      0,
			MessageCode: "N0000",
			MessageText: "",
		}
	}

	// 填充响应体Data部
	if data, exists := ctx.Get("ResponseData"); exists {
		restOutDto.Data = data
	} else {
		restOutDto.Data = nil
	}

	// 设置响应
	ctx.Render(200, render.JSON{Data: restOutDto})
}
