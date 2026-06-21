package middleware

import (
	"fmt"
	"nas-dashboard/pkg/logger"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 恢复中间件 - 捕获panic并记录日志
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取调用栈信息
				stack := debug.Stack()

				// 记录详细的错误日志
				logger.Error("Panic recovered", logger.LogFields{
					"error": fmt.Sprintf("%v", err),
					"path":  c.Request.URL.Path,
					"method": c.Request.Method,
					"ip":    c.ClientIP(),
					"stack": string(stack),
				})

				// 返回友好的错误信息
				c.JSON(500, gin.H{
					"error": "Internal server error",
					"code":  500,
				})

				c.Abort()
			}
		}()
		c.Next()
	}
}