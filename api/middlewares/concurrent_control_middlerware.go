package middlewares

import (
	"LaoQGChat/internal/myerrors"
	"github.com/gin-gonic/gin"
	"time"
)

func ConcurrentControlHandler(maxConcurrentCount int, maxTimeoutDuration int) gin.HandlerFunc {
	var sem = make(chan struct{}, maxConcurrentCount)

	return func(ctx *gin.Context) {
		select {
		// 成功获取令牌
		case sem <- struct{}{}:
			// 处理完毕后释放令牌
			defer func() { <-sem }()

			// 下一层
			ctx.Next()
		case <-time.After(time.Duration(maxTimeoutDuration) * time.Second):
			// 在最大超时时间前没有获取到令牌
			err := &myerrors.CustomError{
				StatusCode:  300,
				MessageCode: "ECC01",
				MessageText: "当前请求过多，请稍后重试。",
			}
			_ = ctx.Error(err)
			ctx.Abort()
		}
	}
}
