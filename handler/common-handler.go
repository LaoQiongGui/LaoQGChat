package handler

import (
	"LaoQGChat/dto"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func HandlerBuilder[T any](handlerFunc func(c *gin.Context) *T) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Keys = make(map[string]any)
		ctx.Keys["StatusCode"] = 0
		ctx.Keys["MessageCode"] = "N0000"
		ctx.Keys["MessageText"] = ""
		restOutDto := dto.RestOutDto{}
		defer errorHandler(ctx, &restOutDto)

		if pData := handlerFunc(ctx); pData != nil {
			restOutDto.Data = *pData
		}
	}
}

func errorHandler(ctx *gin.Context, pRestOutDto *dto.RestOutDto) {
	if err := recover(); err != nil && ctx.GetInt("StatusCode") == 0 {
		pRestOutDto.Common.Status = 990
		pRestOutDto.Common.MessageCode = "E9999"
		pRestOutDto.Common.MessageText = "System Error"
	} else {
		pRestOutDto.Common.Status = ctx.GetInt("StatusCode")
		pRestOutDto.Common.MessageCode = ctx.GetString("MessageCode")
		pRestOutDto.Common.MessageText = ctx.GetString("MessageText")
	}
	ctx.Render(200, render.JSON{Data: *pRestOutDto})
}
