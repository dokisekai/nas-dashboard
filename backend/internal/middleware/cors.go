package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

// CORS 中间件 - 安全的跨域配置
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求来源
		origin := c.Request.Header.Get("Origin")

		// 添加调试日志
		if origin != "" {
			log.Printf("[CORS] Request from origin: %s, Path: %s", origin, c.Request.URL.Path)
		}

		// 允许的域名列表 (生产环境应该修改为实际域名)
		allowedOrigins := []string{
			"http://localhost:5173",
			"https://localhost:5173",
			"http://127.0.0.1:5173",
			"http://192.168.50.10:5173",
			"https://192.168.50.10:5173",
			// Immich 端口
			"http://localhost:2283",
			"https://localhost:2283",
			"http://127.0.0.1:2283",
			"http://192.168.50.10:2283",
			"https://192.168.50.10:2283",
			// NAS 系统端口
			"http://192.168.50.10:8888",
			"https://192.168.50.10:8888",
			"http://localhost:8888",
			"https://localhost:8888",
			// 添加你的生产域名
			// "https://your-domain.com",
		}

		// 检查来源是否允许
		allowed := false
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}
		}

		// 如果来源被允许，设置CORS头
		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			if allowed {
				c.AbortWithStatus(204)
			} else {
				c.AbortWithStatus(403)
			}
			return
		}

		// 如果来源不被允许，拒绝请求
		if !allowed && origin != "" {
			c.JSON(403, gin.H{"error": "Origin not allowed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
