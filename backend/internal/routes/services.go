package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerServices 系统服务管理 + Docker 容器/镜像/网络/卷管理路由。
func registerServices(g *gin.RouterGroup) {
	services := g.Group("/services")
	requireAuth(services)
	{
		services.GET("", api.GetServices)
		services.POST("/:name/start", api.StartService)
		services.POST("/:name/stop", api.StopService)
		services.POST("/:name/restart", api.RestartService)
		services.POST("/:name/enable", api.EnableService)
		services.POST("/:name/disable", api.DisableService)
	}

	docker := g.Group("/docker")
	requireAuth(docker)
	{
		// 容器
		docker.GET("/containers", api.GetContainers)
		docker.POST("/containers/:id/start", api.StartContainer)
		docker.POST("/containers/:id/stop", api.StopContainer)
		docker.POST("/containers/:id/restart", api.RestartContainer)
		docker.DELETE("/containers/:id", api.RemoveContainer)
		docker.GET("/containers/:id/logs", api.GetContainerLogs)
		docker.GET("/containers/:id/stats", api.GetContainerStats)
		docker.POST("/containers/:id/exec", api.ExecInContainer)

		// 镜像
		docker.GET("/images", api.GetDockerImages)
		docker.DELETE("/images/:id", api.RemoveImage)
		docker.POST("/images/pull", api.PullImage)

		// 网络/卷
		docker.GET("/networks", api.GetDockerNetworks)
		docker.GET("/volumes", api.GetDockerVolumes)
	}
}
