package main

import (
	"log"
	"nas-dashboard/internal/api"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库（非阻塞，失败时继续运行）
	if err := initDatabase(); err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Server will continue without database support")
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 配置 CORS
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API 路由组
	apiGroup := r.Group("/api")
	{
		// 认证路由
		apiGroup.POST("/auth/login", api.Login)
		apiGroup.POST("/auth/refresh", api.RefreshToken)

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
			storage.GET("/usage", api.GetDiskUsage) // 获取指定路径的磁盘使用情况
		}

		// 服务管理路由
		services := apiGroup.Group("/services")
		services.Use(middleware.Auth())
		{
			services.GET("", api.GetServices)
			services.POST("/:name/start", api.StartService)
			services.POST("/:name/stop", api.StopService)
			services.POST("/:name/restart", api.RestartService)
			services.POST("/:name/enable", api.EnableService)    // 启用服务（开机自启）
			services.POST("/:name/disable", api.DisableService)   // 禁用服务
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

			// SSH 密钥管理路由
			users.GET("/ssh-keys", api.GetSSHKeys)
			users.POST("/ssh-keys", api.AddKey)
			users.DELETE("/ssh-keys/:id", api.DeleteKey)

			// 当前用户相关
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

		// 文件管理路由
		files := apiGroup.Group("/files")
		files.Use(middleware.Auth())
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

	// WebSocket 路由
	r.GET("/ws/monitor", api.WSMonitor)

	// 启动服务器
	log.Println("Server starting on 0.0.0.0:8888")
	if err := r.Run("0.0.0.0:8888"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// initDatabase 初始化数据库
func initDatabase() error {
	log.Println("Initializing database...")

	// 加载数据库配置
	cfg := database.LoadConfig()

	// 连接数据库
	if err := database.Connect(cfg); err != nil {
		return err
	}

	// 运行迁移
	if err := database.Migrate(); err != nil {
		return err
	}

	log.Println("Database initialized successfully")
	return nil
}
