package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerStatic 托管前端 SPA 静态资源，并对非 API/WS 请求回退到 index.html。
func registerStatic(r *gin.Engine) {
	if _, err := os.Stat("static"); err != nil {
		return
	}

	r.StaticFS("/assets", http.Dir("static/assets"))
	r.StaticFile("/favicon.svg", "static/favicon.svg")
	r.StaticFile("/icons.svg", "static/icons.svg")

	// SPA 回退：所有非 API/WS 路径均返回 index.html，由前端路由处理
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/ws") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.File("static/index.html")
	})
}

// registerWebSocket 注册 WebSocket 端点（实时监控推送）。
func registerWebSocket(r *gin.Engine) {
	r.GET("/ws/monitor", api.WSMonitor)
}
