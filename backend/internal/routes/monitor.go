package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerMonitor 系统资源监控 + 功耗监控路由。
func registerMonitor(g *gin.RouterGroup) {
	monitor := g.Group("/monitor")
	requireAuth(monitor)
	{
		monitor.GET("/cpu", api.GetCPU)
		monitor.GET("/memory", api.GetMemory)
		monitor.GET("/disk", api.GetDisk)
		monitor.GET("/network", api.GetNetwork)

		// 扩展监控：进程、服务、温度、事件、日志、告警
		monitorAPI := api.GetMonitorAPI()
		monitor.GET("/processes", monitorAPI.GetProcesses)
		monitor.GET("/processes/:pid", monitorAPI.GetProcess)
		monitor.DELETE("/processes/:pid", monitorAPI.KillProcess)
		monitor.GET("/services", monitorAPI.GetServices)
		monitor.GET("/services/:name", monitorAPI.GetService)
		monitor.POST("/services/:name/start", monitorAPI.StartService)
		monitor.POST("/services/:name/stop", monitorAPI.StopService)
		monitor.POST("/services/:name/restart", monitorAPI.RestartService)
		monitor.GET("/temperature", monitorAPI.GetTemperature)
		monitor.GET("/events", monitorAPI.GetEvents)
		monitor.GET("/logs", monitorAPI.GetLogs)
		monitor.POST("/logs/clear", monitorAPI.ClearLogs)
		monitor.GET("/alerts", monitorAPI.GetAlerts)
		monitor.POST("/alerts", monitorAPI.CreateAlert)
		monitor.PUT("/alerts/:id", monitorAPI.UpdateAlert)
		monitor.DELETE("/alerts/:id", monitorAPI.DeleteAlert)
	}

	// 功耗监控
	power := g.Group("/power")
	requireAuth(power)
	{
		power.GET("/current", api.GetPowerCurrent)
		power.GET("/history", api.GetPowerHistory)
		power.GET("/statistics", api.GetPowerStatistics)
		power.GET("/overview", api.GetPowerOverview)
	}
}
