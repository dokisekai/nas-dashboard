package middleware

import (
	"nas-dashboard/pkg/logger"
	"nas-dashboard/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

// 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 执行请求
		c.Next()

		// 处理错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			logger.APIError(
				c.Request.Method,
				c.Request.URL.Path,
				err.Err,
				c.Writer.Status(),
			)

			// 根据错误类型返回不同的响应
			switch e := err.Err.(type) {
			case *ValidationError:
				response.BadRequest(c, e.Message)
			case *NotFoundError:
				response.NotFound(c, e.Message)
			case *UnauthorizedError:
				response.Unauthorized(c, e.Message)
			case *ForbiddenError:
				response.Forbidden(c, e.Message)
			default:
				response.InternalError(c, "Internal server error")
			}
		} else {
			// 记录成功的API调用
			duration := time.Since(startTime).Milliseconds()
			logger.APIResponse(
				c.Writer.Status(),
				"Success",
				duration,
			)
		}
	}
}

// 请求日志中间件
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Info("Request", logger.LogFields{
			"client_ip": clientIP,
			"method":    method,
			"path":      path,
			"status":    statusCode,
			"latency":   latency.String(),
		})
	}
}

// 自定义错误类型
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

type NotFoundError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}