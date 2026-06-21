package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nas-dashboard/internal/api"
)

// registerStorage 存储相关路由：磁盘、分区、SMART、RAID、LVM、SMB、
// 存储池、配额、同步与备份任务。
func registerStorage(g *gin.RouterGroup, db *gorm.DB) {
	storage := g.Group("/storage")
	requireAuth(storage)
	{
		// 基础磁盘/SMB 管理
		storage.GET("/disks", api.GetDisks)
		storage.POST("/disks/format", api.FormatDisk)
		storage.POST("/mount", api.MountDisk)
		storage.POST("/umount", api.UmountDisk)
		storage.GET("/smb", api.GetSMBShares)
		storage.POST("/smb", api.CreateSMBShare)
		storage.PUT("/smb/:name", api.UpdateSMBShare)
		storage.DELETE("/smb/:name", api.DeleteSMBShare)
		storage.GET("/usage", api.GetDiskUsage)

		// 高级磁盘管理（分区 / SMART / 健康 / 基准测试）
		storage.GET("/disks/:device/partitions", api.GetDiskPartitions)
		storage.POST("/disks/:device/partitions", api.CreateDiskPartition)
		storage.DELETE("/disks/:device/partitions/:number", api.DeleteDiskPartition)
		storage.GET("/disks/:device/smart", api.GetDiskSmart)
		storage.POST("/disks/:device/test", api.RunDiskSmartTest)
		storage.GET("/disks/:device/health", api.GetDiskHealth)
		storage.POST("/disks/:device/benchmark", api.RunDiskBenchmark)

		// RAID 管理
		storage.GET("/raid", api.GetRAIDArrays)
		storage.GET("/raid/:name", api.GetRAIDArray)
		storage.POST("/raid", api.CreateRAID)
		storage.DELETE("/raid/:name", api.DeleteRAID)
		storage.POST("/raid/:name/add", api.AddDiskToRAID)
		storage.POST("/raid/:name/remove", api.RemoveDiskFromRAID)

		// LVM 管理
		storage.GET("/lvm/pv", api.GetPhysicalVolumes)
		storage.POST("/lvm/pv", api.CreatePhysicalVolume)
		storage.GET("/lvm/vg", api.GetVolumeGroups)
		storage.POST("/lvm/vg", api.CreateVolumeGroup)
		storage.DELETE("/lvm/vg/:name", api.DeleteVolumeGroup)
		storage.GET("/lvm/lv", api.GetLogicalVolumes)
		storage.POST("/lvm/lv", api.CreateLogicalVolume)
		storage.DELETE("/lvm/lv/:vg/:name", api.DeleteLogicalVolume)

		// 存储池（MergerFS 等）
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

		// 同步任务（rsync 文件级同步）
		syncAPI := api.NewSyncAPI(db)
		sync := storage.Group("/sync")
		{
			sync.GET("/jobs", syncAPI.GetSyncJobs)
			sync.POST("/jobs", syncAPI.CreateSyncJob)
			sync.PUT("/jobs/:id", syncAPI.UpdateSyncJob)
			sync.DELETE("/jobs/:id", syncAPI.DeleteSyncJob)
			sync.POST("/jobs/:id/run", syncAPI.RunSyncJob)
		}

		// 备份任务（Restic）
		resticAPI := api.NewResticAPI(db)
		backups := storage.Group("/backup")
		{
			// 健康检查 + 全局设置 + WebDAV 配置
			backups.GET("/ping", resticAPI.Ping)
			backups.GET("/settings", resticAPI.GetSettings)
			backups.PUT("/settings", resticAPI.UpdateSettings)
			backups.GET("/webdav", resticAPI.GetWebDAV)
			backups.PUT("/webdav", resticAPI.UpdateWebDAV)
			backups.POST("/webdav/test", resticAPI.TestWebDAV)

			// 仓库 CRUD + 维护
			backups.GET("/repos", resticAPI.ListRepos)
			backups.POST("/repos", resticAPI.CreateRepo)
			backups.PUT("/repos/:id", resticAPI.UpdateRepo)
			backups.DELETE("/repos/:id", resticAPI.DeleteRepo)
			backups.POST("/repos/:id/check", resticAPI.CheckRepo)
			backups.POST("/repos/:id/refresh", resticAPI.RefreshRepo)
			backups.POST("/repos/:id/test", resticAPI.TestRepoConnection)
			backups.POST("/repos/:id/init", resticAPI.InitRepo)
			backups.POST("/repos/:id/unlock", resticAPI.UnlockRepo)

			// 快照查看/恢复/搜索
			backups.GET("/repos/:id/snapshots", resticAPI.ListSnapshots)
			backups.GET("/repos/:id/snapshots/:sid", resticAPI.SnapshotDetail)
			backups.GET("/repos/:id/snapshots/:sid/ls", resticAPI.ListSnapshotFiles)
			backups.DELETE("/repos/:id/snapshots/:sid", resticAPI.DeleteSnapshot)
			backups.GET("/repos/:id/diff", resticAPI.DiffSnapshots)
			backups.GET("/repos/:id/find", resticAPI.FindInSnapshots)
			backups.POST("/repos/:id/restore", resticAPI.Restore)
			backups.GET("/repos/:id/restore/logs", resticAPI.RestoreLogs)

			// 备份任务
			backups.GET("/tasks", resticAPI.ListTasks)
			backups.POST("/tasks", resticAPI.CreateTask)
			backups.PUT("/tasks/:id", resticAPI.UpdateTask)
			backups.DELETE("/tasks/:id", resticAPI.DeleteTask)
			backups.POST("/tasks/:id/run", resticAPI.RunTask)
			backups.GET("/tasks/:id/status", resticAPI.TaskStatus)
			backups.GET("/tasks/:id/logs", resticAPI.TaskLogs)

			// 仓库间同步（restic copy）
			backups.GET("/sync-jobs", resticAPI.ListSyncJobs)
			backups.POST("/sync-jobs", resticAPI.CreateSyncJob)
			backups.DELETE("/sync-jobs/:id", resticAPI.DeleteSyncJob)
			backups.POST("/sync-jobs/:id/run", resticAPI.RunSyncJob)
			backups.GET("/sync-jobs/:id/logs", resticAPI.RunSyncJobLogs)
		}

		// 配额
		quota := storage.Group("/quota")
		quotaAPI := api.GetQuotaAPI()
		{
			quota.GET("/users", quotaAPI.GetAllQuotas)
			quota.GET("/groups", quotaAPI.GetAllGroupQuotas)
			quota.GET("/report", quotaAPI.GetQuotaReport)
		}
	}
}
