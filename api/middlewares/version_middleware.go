package middlewares

import (
	"LaoQGChat/internal/myerrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func VersionHandler(version string) gin.HandlerFunc {
	versionList := strings.Split(version, ".")

	return func(ctx *gin.Context) {
		// 前处理
		versionInList := strings.Split(ctx.GetHeader("Version"), ".")
		if len(versionInList) < 2 {
			err := &myerrors.CustomError{
				StatusCode:  300,
				MessageCode: "EVE01",
				MessageText: "版本号格式错误。",
			}
			_ = ctx.AbortWithError(http.StatusUpgradeRequired, err)
			return
		}

		if versionInList[0] != versionList[0] || versionInList[1] != versionList[1] {
			err := &myerrors.CustomError{
				StatusCode:  300,
				MessageCode: "EVE02",
				MessageText: "版本过低，请获取最新的app。",
			}
			_ = ctx.AbortWithError(http.StatusUpgradeRequired, err)
			return
		}

		// 下一层
		ctx.Next()

		// 后处理
	}
}
