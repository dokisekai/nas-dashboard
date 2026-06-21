package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerSystem 系统信息 + 系统操作（重启/关机/调度）路由。
func registerSystem(g *gin.RouterGroup) {
	system := g.Group("/system")
	requireAuth(system)
	{
		system.GET("/info", api.GetSystemInfo)
		system.GET("/hardware", api.GetHardwareDetails)
		system.GET("/power", api.GetPowerUsage)
		system.GET("/uptime", api.GetSystemUptime)
		system.GET("/ups/status", api.GetUPSStatus)

		// 系统操作（重启 / 关机 / 调度）
		operations := system.Group("/operations")
		{
			operations.POST("/restart", api.RestartSystem)
			operations.POST("/shutdown", api.ShutdownSystem)
			operations.POST("/cancel", api.CancelShutdown)
			operations.POST("/reboot-immediate", api.RebootSystemImmediately)
			operations.POST("/poweroff-immediate", api.PoweroffSystemImmediately)
			operations.POST("/schedule-shutdown", api.ScheduleShutdown)
			operations.POST("/schedule-restart", api.ScheduleRestart)
			operations.GET("/status", api.GetShutdownStatus)
		}
	}
}

// registerFiles 文件管理路由。
func registerFiles(g *gin.RouterGroup) {
	files := g.Group("/files")
	requireAuth(files)
	{
		files.POST("/list", api.ListFiles)
		files.GET("/info", api.GetFileInfo)
		files.GET("/download", api.DownloadFile)
		files.POST("/upload", api.UploadFile)
		files.POST("/directory", api.CreateDirectory)
		files.POST("/move", api.MoveFile)
		files.POST("/delete", api.DeleteFile)
	}
}

// registerBackups 旧版备份恢复接口（保留兼容）。
func registerBackups(g *gin.RouterGroup) {
	backups := g.Group("/backups")
	requireAuth(backups)
	{
		backups.GET("", api.GetBackups)
		backups.POST("", api.CreateBackup)
		backups.GET("/:id", api.GetBackup)
		backups.DELETE("/:id", api.DeleteBackup)
		backups.POST("/restore", api.RestoreBackup)
		backups.GET("/:id/download", api.DownloadBackup)
	}
}

// registerConfigs 系统配置项 CRUD 路由。
func registerConfigs(g *gin.RouterGroup) {
	configs := g.Group("/configs")
	requireAuth(configs)
	{
		configs.GET("", api.GetConfigs)
		configs.GET("/public", api.GetPublicConfigs)
		configs.GET("/:key", api.GetConfig)
		configs.POST("", api.SetConfig)
		configs.DELETE("/:key", api.DeleteConfig)
		configs.POST("/bulk", api.BulkSetConfig)
	}
}
