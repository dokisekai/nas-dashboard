package main

import (
	"log"
	"nas-dashboard/internal/api"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/websocket"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 设置 Gin 模式
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.ReleaseMode
	}
	gin.SetMode(ginMode)

	// 初始化数据库
	if err := initDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化 WebSocket 管理器
	websocket.InitWebSocketManager()
	defer websocket.GetWebSocketManager().Stop()

	// 创建 Gin 路由
	r := gin.Default()

	// 配置中间件
	r.Use(middleware.CORS())
	r.Use(loggingMiddleware())
	r.Use(recoveryMiddleware())

	// 健康检查
	r.GET("/health", healthCheck)

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 公开配置路由
		apiGroup.GET("/config/public", api.GetPublicConfigs)

		// 认证路由
		auth := apiGroup.Group("/auth")
		{
			auth.POST("/login", api.Login)
			auth.POST("/refresh", api.RefreshToken)
			auth.POST("/logout", api.Logout)
			auth.POST("/logout-all", middleware.Auth(), api.LogoutAll)
		}

		// 监控路由 (需要认证)
		monitor := apiGroup.Group("/monitor")
		monitor.Use(middleware.Auth())
		{
			monitor.GET("/cpu", api.GetCPU)
			monitor.GET("/memory", api.GetMemory)
			monitor.GET("/disk", api.GetDisk)
			monitor.GET("/network", api.GetNetwork)
		}

		// 存储管理路由
		storage := apiGroup.Group("/storage")
		storage.Use(middleware.Auth())
		{
			storage.GET("/disks", api.GetDisks)
			storage.POST("/mount", api.MountDisk)
			storage.POST("/umount", api.UmountDisk)
			storage.GET("/smb", api.GetSMBShares)
			storage.POST("/smb", api.CreateSMBShare)
			storage.PUT("/smb/:name", api.UpdateSMBShare)
			storage.DELETE("/smb/:name", api.DeleteSMBShare)
			storage.GET("/usage", api.GetDiskUsage)
		}

		// 文件管理路由
		files := apiGroup.Group("/files")
		files.Use(middleware.Auth())
		{
			files.POST("/list", api.ListFiles)
			files.GET("/info", api.GetFileInfo)
			files.GET("/download", api.DownloadFile)
			files.POST("/upload", api.UploadFile)
			files.POST("/mkdir", api.CreateDirectory)
			files.POST("/move", api.MoveFile)
			files.DELETE("/delete", api.DeleteFile)
		}

		// 服务管理路由
		services := apiGroup.Group("/services")
		services.Use(middleware.Auth())
		{
			services.GET("", api.GetServices)
			services.POST("/:name/start", api.StartService)
			services.POST("/:name/stop", api.StopService)
			services.POST("/:name/restart", api.RestartService)
			services.POST("/:name/enable", api.EnableService)
			services.POST("/:name/disable", api.DisableService)
		}

		// Docker 路由
		docker := apiGroup.Group("/docker")
		docker.Use(middleware.Auth())
		{
			docker.GET("/containers", api.GetContainers)
			docker.POST("/containers/:id/start", api.StartContainer)
			docker.POST("/containers/:id/stop", api.StopContainer)
			docker.POST("/containers/:id/restart", api.RestartContainer)
			docker.DELETE("/containers/:id", api.RemoveContainer)
			docker.GET("/containers/:id/logs", api.GetContainerLogs)
			docker.GET("/containers/:id/stats", api.GetContainerStats)
			docker.POST("/containers/:id/exec", api.ExecInContainer)
			docker.GET("/images", api.GetDockerImages)
			docker.DELETE("/images/:id", api.RemoveImage)
			docker.POST("/images/pull", api.PullImage)
		}

		// 用户管理路由
		users := apiGroup.Group("/users")
		users.Use(middleware.Auth())
		{
			users.GET("", api.GetUsers)
			users.POST("", api.CreateUser)
			users.PUT("/:username", api.UpdateUser)
			users.DELETE("/:username", api.DeleteUser)
			users.GET("/:username", api.GetUser)
			users.GET("/ssh-keys", api.GetSSHKeys)
			users.POST("/ssh-keys", api.AddKey)
			users.DELETE("/ssh-keys/:id", api.DeleteKey)
			users.GET("/me", api.GetCurrentUser)
			users.POST("/me/password", api.ChangeCurrentUserPassword)
			users.GET("/:username/quota", api.GetUserDiskQuota)
		}

		// 系统组路由
		groups := apiGroup.Group("/groups")
		groups.Use(middleware.Auth())
		{
			groups.GET("", api.GetGroups)
		}

		// 系统信息路由
		system := apiGroup.Group("/system")
		system.Use(middleware.Auth())
		{
			system.GET("/info", api.GetSystemInfo)
		}

		// 系统配置路由
		config := apiGroup.Group("/config")
		config.Use(middleware.Auth())
		{
			config.GET("", api.GetConfigs)
			config.GET("/:key", api.GetConfig)
			config.POST("", api.SetConfig)
			config.POST("/bulk", api.BulkSetConfig)
			config.DELETE("/:key", api.DeleteConfig)
			config.POST("/:key/reset", api.ResetConfig)
		}

		// 插件管理路由
		plugins := apiGroup.Group("/plugins")
		plugins.Use(middleware.Auth())
		{
			plugins.GET("", api.GetPlugins)
			plugins.GET("/:name", api.GetPlugin)
			plugins.POST("", api.InstallPlugin)
			plugins.PUT("/:name", api.UpdatePlugin)
			plugins.DELETE("/:name", api.UninstallPlugin)
			plugins.POST("/:name/enable", api.EnablePlugin)
			plugins.POST("/:name/disable", api.DisablePlugin)
			plugins.POST("/:name/action", api.PluginAction)
			plugins.GET("/:name/logs", api.GetPluginLogs)
			plugins.GET("/:name/config", api.GetPluginConfig)
			plugins.PUT("/:name/config", api.UpdatePluginConfig)
		}

		// 备份恢复路由
		backups := apiGroup.Group("/backups")
		backups.Use(middleware.Auth())
		{
			backups.GET("", api.GetBackups)
			backups.GET("/:id", api.GetBackup)
			backups.POST("", api.CreateBackup)
			backups.DELETE("/:id", api.DeleteBackup)
			backups.POST("/restore", api.RestoreBackup)
			backups.GET("/:id/download", api.DownloadBackup)
		}
	}

	// WebSocket 路由
	wsManager := websocket.GetWebSocketManager()
	r.GET("/ws/monitor", func(c *gin.Context) {
		wsManager.HandleWebSocket(c)
	})

	// WebSocket API 路由
	wsAPI := r.Group("/api/ws")
	wsAPI.Use(middleware.Auth(), middleware.AdminOnly())
	{
		wsAPI.GET("/:action", api.HandleWebSocketAPI)
		wsAPI.POST("/broadcast", api.BroadcastMonitorData)
	}

	// 启动服务器
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = "8888"
	}

	addr := host + ":" + port
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDatabase 初始化数据库
func initDatabase() error {
	cfg := database.LoadConfig()

	if err := database.Connect(cfg); err != nil {
		return err
	}

	if err := database.Migrate(); err != nil {
		return err
	}

	return nil
}

// healthCheck 健康检查
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
		"database": func() string {
			if database.GetDB() != nil {
				return "connected"
			}
			return "disconnected"
		}(),
	})
}

// loggingMiddleware 日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		if raw != "" {
			path = path + "?" + raw
		}

		log.Printf("[%s] %s %s %d %v",
			c.Request.Method,
			path,
			c.ClientIP(),
			c.Writer.Status(),
			latency,
		)
	}
}

// recoveryMiddleware 恢复中间件
func recoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
		}()
		c.Next()
	}
}
