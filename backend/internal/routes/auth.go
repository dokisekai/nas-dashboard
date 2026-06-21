package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
	"nas-dashboard/internal/middleware"
)

// registerHealth 健康检查与系统初始化相关路由。
func registerHealth(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
}

// registerAuth 认证 + 系统初始化（公开端点）路由。
func registerAuth(g *gin.RouterGroup) {
	g.POST("/auth/login", api.Login)
	g.POST("/auth/refresh", api.RefreshToken)

	// Immich 自动登录（携带登录态）
	g.GET("/immich/auth-login", middleware.Auth(), api.ImmichAuthLogin)

	// 系统初始化（公开）
	g.GET("/system/init-status", api.GetInitStatus)
	g.POST("/system/initialize", api.InitializeSystem)
}
