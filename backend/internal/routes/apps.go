package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerApps 应用管理路由：Docker 容器、Compose 项目、服务目录。
// 不同于 /api/docker（基于 docker CLI，主机直跑才能用），这一组路由
// 走 /var/run/docker.sock 的 HTTP API，所以容器内部也能管理宿主机上的服务。
func registerApps(g *gin.RouterGroup) {
	apps := g.Group("/apps")
	requireAuth(apps)
	{
		// 服务目录（应用商店视图，带状态）
		apps.GET("/catalog", api.ListServiceCatalog)

		// Compose 项目分组视图
		apps.GET("/projects", api.ListComposeProjects)

		// 全部容器（含已停止）
		apps.GET("/containers", api.ListAppContainers)
		apps.GET("/containers/:name/logs", api.AppContainerLogs)
		apps.GET("/containers/:name/stats", api.AppContainerStats)
		apps.POST("/containers/:name/:action", api.AppContainerAction) // start|stop|restart|remove
	}
}
