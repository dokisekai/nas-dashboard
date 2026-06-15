package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"nas-dashboard/internal/api"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库（非阻塞，失败时继续运行）
	db, err := initDatabase()
	if err != nil {
		log.Printf("Warning: Failed to initialize database: %v", err)
		log.Println("Server will continue without database support")
	}

	// 初始化API处理器
	api.InitAPI(db)

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

		// 系统初始化路由 (公开)
		apiGroup.GET("/system/init-status", api.GetInitStatus)
		apiGroup.POST("/system/initialize", api.InitializeSystem)

		// 监控路由 (需要认证)
		monitor := apiGroup.Group("/monitor")
		monitor.Use(middleware.Auth())
		{
					monitorAPI := api.GetMonitorAPI()
			monitor.GET("/cpu", api.GetCPU)
			monitor.GET("/memory", api.GetMemory)
			monitor.GET("/disk", api.GetDisk)
			monitor.GET("/network", api.GetNetwork)

			// 扩展监控功能
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

		// 功耗监控路由
		power := apiGroup.Group("/power")
		power.Use(middleware.Auth())
		{
			power.GET("/current", api.GetPowerCurrent)
			power.GET("/history", api.GetPowerHistory)
			power.GET("/statistics", api.GetPowerStatistics)
			power.GET("/overview", api.GetPowerOverview)
		}

			// 网络管理路由
			network := apiGroup.Group("/network")
			network.Use(middleware.Auth())
			{
				// 网络接口管理
				network.GET("/interfaces", api.GetNetworkInterfaces)
				network.GET("/interfaces/ethernet", api.GetEthernetInterfaces)
				network.GET("/interfaces/wifi", api.GetWiFiInterfaces)
						// 接口配置管理（参数化路由）
						network.GET("/interface/:interface/config", api.GetInterfaceConfig)
						network.PUT("/interface/:interface/config", api.SetInterfaceConfig)
						network.POST("/interface/:interface/restart", api.RestartInterface)

						// 通用接口操作路由
						network.POST("/interface/:interface/:action", api.ControlInterface)


				// PPPoE管理
				network.GET("/interface/:interface/pppoe", api.GetPPPoEConfig)
				network.POST("/interface/:interface/pppoe", api.ConfigurePPPoE)

				// 代理配置管理
				network.GET("/proxy", api.GetProxyConfig)
				network.POST("/proxy", api.SetProxyConfig)

				// Wi-Fi管理
				network.GET("/wifi/scan", api.ScanWiFiNetworks)
				network.POST("/wifi/connect", api.ConnectToWiFi)
				network.POST("/wifi/disconnect", api.DisconnectWiFi)
				network.GET("/wifi/current", api.GetCurrentWiFiConnection)

				// DNS管理
				network.GET("/dns", api.GetDNSConfig)
				network.POST("/dns", api.SetDNSConfig)

				// IP管理
				network.PUT("/ip", api.UpdateIPConfig)
				}

		// 存储管理路由
		storage := apiGroup.Group("/storage")
		storage.Use(middleware.Auth())
		{
			storage.GET("/disks", api.GetDisks)
			storage.POST("/disks/format", api.FormatDisk)
			storage.POST("/mount", api.MountDisk)
			storage.POST("/umount", api.UmountDisk)
			storage.GET("/smb", api.GetSMBShares)
			storage.POST("/smb", api.CreateSMBShare)
			storage.PUT("/smb/:name", api.UpdateSMBShare)
			storage.DELETE("/smb/:name", api.DeleteSMBShare)
			storage.GET("/usage", api.GetDiskUsage) // 获取指定路径的磁盘使用情况

			// 高级磁盘管理
			storage.GET("/disks/:device/partitions", api.GetDiskPartitions)
			storage.POST("/disks/:device/partitions", api.CreateDiskPartition)
			storage.DELETE("/disks/:device/partitions/:number", api.DeleteDiskPartition)
			storage.GET("/disks/:device/smart", api.GetDiskSmart)
			storage.POST("/disks/:device/test", api.RunDiskSmartTest)
			storage.GET("/disks/:device/health", api.GetDiskHealth)
			storage.POST("/disks/:device/benchmark", api.RunDiskBenchmark)

			// RAID管理
			storage.GET("/raid", api.GetRAIDArrays)
			storage.GET("/raid/:name", api.GetRAIDArray)
			storage.POST("/raid", api.CreateRAID)
			storage.DELETE("/raid/:name", api.DeleteRAID)
			storage.POST("/raid/:name/add", api.AddDiskToRAID)
			storage.POST("/raid/:name/remove", api.RemoveDiskFromRAID)

			// LVM管理
			storage.GET("/lvm/pv", api.GetPhysicalVolumes)
			storage.POST("/lvm/pv", api.CreatePhysicalVolume)
			storage.GET("/lvm/vg", api.GetVolumeGroups)
			storage.POST("/lvm/vg", api.CreateVolumeGroup)
			storage.DELETE("/lvm/vg/:name", api.DeleteVolumeGroup)
			storage.GET("/lvm/lv", api.GetLogicalVolumes)
			storage.POST("/lvm/lv", api.CreateLogicalVolume)
			storage.DELETE("/lvm/lv/:vg/:name", api.DeleteLogicalVolume)

			// 存储池管理路由
			pools := storage.Group("/pools")
			storagePoolAPI := api.GetStoragePoolAPI()
			{
				pools.GET("", storagePoolAPI.GetPools)
				pools.POST("", storagePoolAPI.CreatePool)
				pools.GET("/:name", storagePoolAPI.GetPool)
				pools.PUT("/:name", storagePoolAPI.UpdatePool)
				pools.DELETE("/:name", storagePoolAPI.DeletePool)
				pools.POST("/:name/disks", storagePoolAPI.AddDisk)
				pools.DELETE("/:name/disks/:device", storagePoolAPI.RemoveDisk)
				pools.GET("/:name/branches", storagePoolAPI.GetPoolBranches)
				pools.POST("/:name/mount", storagePoolAPI.MountPool)
				pools.POST("/:name/umount", storagePoolAPI.UmountPool)
				pools.POST("/:name/balance", storagePoolAPI.BalancePool)
				pools.POST("/:name/scan", storagePoolAPI.ScanPool)
			}

			// 同步管理
			syncAPI := api.NewSyncAPI(db)
			sync := storage.Group("/sync")
			{
				sync.GET("/jobs", syncAPI.GetSyncJobs)
				sync.POST("/jobs", syncAPI.CreateSyncJob)
				sync.POST("/jobs/:id/run", syncAPI.RunSyncJob)
			}

			// 备份管理 (Restic)
			backupAPI := api.NewBackupAPI(db)
			backups := storage.Group("/backup")
			{
				backups.GET("/repos", backupAPI.GetRepos)
				backups.POST("/repos", backupAPI.CreateRepo)
				backups.GET("/tasks", backupAPI.GetTasks)
				backups.POST("/tasks", backupAPI.CreateTask)
			}

			// 配额管理路由
			quota := storage.Group("/quota")
					quotaAPI := api.GetQuotaAPI()
			{
				quota.GET("/users", quotaAPI.GetAllQuotas)
				quota.GET("/groups", quotaAPI.GetAllGroupQuotas)
				quota.GET("/report", quotaAPI.GetQuotaReport)
			}
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

			// 用户配额管理
			users.GET("/:username/quota", api.GetQuotaAPI().GetUserQuota)
			users.PUT("/:username/quota", api.GetQuotaAPI().SetUserQuota)
		}

		// 系统组路由
		groups := apiGroup.Group("/groups")
		groups.Use(middleware.Auth())
		{
			groups.GET("", api.GetGroups)
				groups.POST("", api.CreateGroup)
				groups.GET("/:name", api.GetGroup)
				groups.PUT("/:name", api.UpdateGroup)
				groups.DELETE("/:name", api.DeleteGroup)
				groups.GET("/:name/members", api.GetGroupMembers)
				groups.POST("/:name/members", api.AddGroupMembers)
				groups.DELETE("/:name/members/:user", api.RemoveGroupMember)
			}

				// SMB用户管理路由
				smbUsers := apiGroup.Group("/smb")
				smbUsers.Use(middleware.Auth())
				{
					smbUsers.GET("/users", api.GetSMBUsers)
					smbUsers.POST("/users/:username/password", api.SetSMBPassword)
					smbUsers.DELETE("/users/:username/password", api.DeleteSMBPassword)
					smbUsers.POST("/users/:username/enable", api.EnableSMBUser)
					smbUsers.POST("/users/:username/disable", api.DisableSMBUser)
					smbUsers.GET("/users/:username/stats", api.GetSMBUserStats)
					smbUsers.GET("/sessions", api.GetSMBSessions)
					smbUsers.DELETE("/sessions/:pid", api.DisconnectSMBSession)
					smbUsers.DELETE("/sessions", api.DisconnectAllSMBSessions)
				}

				// 权限管理路由
				permissions := apiGroup.Group("/permissions")
				permissions.Use(middleware.Auth())
				{
					permissions.GET("/shares", api.GetShares)
					permissions.POST("/shares", api.CreateShare)
					permissions.PUT("/shares/:name", api.UpdateShare)
					permissions.DELETE("/shares/:name", api.DeleteShare)
					permissions.GET("/shares/:name/permissions", api.GetSharePermissions)
					permissions.PUT("/shares/:name/permissions", api.SetSharePermissions)
					permissions.GET("/files", api.GetFilePermissions)
					permissions.PUT("/files/permissions", api.SetFilePermissions)
					permissions.GET("/files/acl", api.GetFileACL)
					permissions.PUT("/files/acl", api.SetFileACL)
				}

			// 系统信息路由
			system := apiGroup.Group("/system")
			system.Use(middleware.Auth())
			{
				system.GET("/info", api.GetSystemInfo)
				system.GET("/hardware", api.GetHardwareDetails)
				system.GET("/power", api.GetPowerUsage)
				system.GET("/uptime", api.GetSystemUptime)
				system.GET("/ups/status", api.GetUPSStatus)

				// 系统操作路由（需要管理员权限）
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

			// 通知管理路由
			notifications := apiGroup.Group("/notifications")
			notifications.Use(middleware.Auth())
			{
				notifications.GET("", api.GetNotifications)
				notifications.POST("/read/:id", api.MarkNotificationRead)
				notifications.POST("/read/all", api.MarkAllNotificationsRead)
				notifications.DELETE("/clear", api.ClearNotifications)
				notifications.POST("", api.CreateNotification)
			}

			// 备份恢复路由
			backups := apiGroup.Group("/backups")
			backups.Use(middleware.Auth())
			{
				backups.GET("", api.GetBackups)
				backups.POST("", api.CreateBackup)
				backups.GET("/:id", api.GetBackup)
				backups.DELETE("/:id", api.DeleteBackup)
				backups.POST("/restore", api.RestoreBackup)
				backups.GET("/:id/download", api.DownloadBackup)
			}

			// 系统配置路由
			configs := apiGroup.Group("/configs")
			configs.Use(middleware.Auth())
			{
				configs.GET("", api.GetConfigs)
				configs.GET("/public", api.GetPublicConfigs)
				configs.GET("/:key", api.GetConfig)
				configs.POST("", api.SetConfig)
				configs.DELETE("/:key", api.DeleteConfig)
				configs.POST("/bulk", api.BulkSetConfig)
			}

			// 防火墙路由
			firewall := apiGroup.Group("/security/firewall")
			firewall.Use(middleware.Auth())
			{
				firewall.GET("/rules", api.GetFirewallRules)
				firewall.POST("/rules", api.CreateFirewallRule)
				firewall.PUT("/rules/:id", api.UpdateFirewallRule)
				firewall.DELETE("/rules/:id", api.DeleteFirewallRule)
				firewall.POST("/apply", api.ApplyFirewallRules)
				firewall.GET("/config", api.GetFirewallConfig)
				firewall.PUT("/config", api.SetFirewallConfig)
			}
		}

		// 托管前端静态文件
		if _, err := os.Stat("static"); err == nil {
			r.StaticFS("/assets", http.Dir("static/assets"))
			r.StaticFile("/favicon.svg", "static/favicon.svg")
			r.StaticFile("/icons.svg", "static/icons.svg")
			
			// SPA 路由：所有非 API 路由都返回 index.html
			r.NoRoute(func(c *gin.Context) {
				path := c.Request.URL.Path
				if !filepath.HasPrefix(path, "/api") && !filepath.HasPrefix(path, "/ws") {
					c.File("static/index.html")
				}
			})
		}

		// WebSocket 路由
		r.GET("/ws/monitor", api.WSMonitor)

		// 启动服务器
		port := "8888"
		certFile := "certs/server.crt"
		keyFile := "certs/server.key"

		if _, err := os.Stat(certFile); err == nil {
			log.Printf("Server starting on https://0.0.0.0:%s (SSL Enabled)", port)
			if err := r.RunTLS("0.0.0.0:"+port, certFile, keyFile); err != nil {
				log.Fatal("Failed to start HTTPS server:", err)
			}
		} else {
			log.Printf("Server starting on http://0.0.0.0:%s", port)
			if err := r.Run("0.0.0.0:"+port); err != nil {
				log.Fatal("Failed to start HTTP server:", err)
			}
		}
	}

	// initDatabase 初始化数据库
	func initDatabase() (*gorm.DB, error) {
		log.Println("Initializing database...")

		// 加载数据库配置
		cfg := database.LoadConfig()

		// 连接数据库
		if err := database.Connect(cfg); err != nil {
			return nil, err
		}

		// 获取数据库连接
		db := database.GetDB()

		// 运行迁移
		if err := database.Migrate(); err != nil {
			return nil, err
		}

		log.Println("Database initialized successfully")
		return db, nil
	}
