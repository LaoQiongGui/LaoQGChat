package handler

import (
	"LaoQGChat/myerror"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func TransactionHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 前处理
		// 开启事务
		tx, err := db.Begin()
		if err != nil {
			err := &myerror.CustomError{
				StatusCode:  300,
				MessageCode: "EDB01",
				MessageText: "数据库连接失败，请联系管理员。",
			}
			panic(err)
		}

		// 结束事务
		defer func() {
			if err := recover(); err != nil {
				switch myError := err.(type) {
				case *myerror.CustomError:
					if myError.StatusCode < 200 {
						// 消息或警告：提交事务
						_ = tx.Commit()
					} else {
						// 异常：回滚事务
						_ = tx.Rollback()
					}
				default:
					// 异常：回滚事务
					_ = tx.Rollback()
				}
				// 向外传递异常
				panic(err)
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
