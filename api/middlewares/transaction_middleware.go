package middlewares

import (
	"LaoQGChat/internal/myerrors"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
)

func TransactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 前处理
		// 开启事务
		tx, err := db.Begin()
		if err != nil {
			err := &myerrors.CustomError{
				StatusCode:  300,
				MessageCode: "EDB01",
				MessageText: "数据库连接失败，请联系管理员。",
			}
			_ = ctx.Error(err)
		}

		// 结束事务
		defer func() {
			// 处理 panic
			if err := recover(); err != nil {
				_ = tx.Rollback()
				panic(err)
			}

			// 处理错误信息
			if err := ctx.Errors.Last(); err != nil {
				// 处理自定义异常
				var myError *myerrors.CustomError
				if errors.As(err.Err, &myError) {
					if myError.StatusCode < 200 {
						// 消息或警告：提交事务
						_ = tx.Commit()
					} else {
						// 异常：回滚事务
						_ = tx.Rollback()
					}
				} else {
					// 异常：回滚事务
					_ = tx.Rollback()
				}
			} else {
				// 正常：提交事务
				_ = tx.Commit()
			}
		}()

		// 下一层
		ctx.Next()

		// 后处理
	}
}
