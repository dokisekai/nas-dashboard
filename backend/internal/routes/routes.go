package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/sso"
)

// Setup 注册全部 HTTP 路由。main 函数只负责创建 gin engine 与调用本函数，
// 具体的路由按域分散在同一个包内的不同文件中，便于查找和维护。
func Setup(r *gin.Engine, db *gorm.DB, ssoServer *sso.SSOServer) {
	registerHealth(r)
	registerSSO(r, ssoServer)

	apiGroup := r.Group("/api")
	{
		registerAuth(apiGroup)
		registerMonitor(apiGroup)
		registerNetwork(apiGroup)
		registerStorage(apiGroup, db)
		registerServices(apiGroup)
		registerUsers(apiGroup)
		registerSystem(apiGroup)
		registerFiles(apiGroup)
		registerBackups(apiGroup)
		registerApps(apiGroup)
		registerConfigs(apiGroup)
		registerFirewall(apiGroup)
		registerOAuth(r, db)
	}

	registerStatic(r)
	registerWebSocket(r)
}

// requireAuth 是一个语法糖：在 group 上挂载 Auth 中间件。
func requireAuth(g *gin.RouterGroup) {
	g.Use(middleware.Auth())
}
