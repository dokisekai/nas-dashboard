package api

import (
	"os"
	"sync"

	"nas-dashboard/pkg/mergerfs"
	"gorm.io/gorm"
)

// 全局API实例（用于跨包访问）
var (
	once         sync.Once
	dbInstance    *gorm.DB
	mergerMgr    *mergerfs.Manager
)

// InitAPI 初始化API处理器和依赖
func InitAPI(db *gorm.DB) {
	once.Do(func() {
		dbInstance = db

		// 创建配置目录
		if err := os.MkdirAll("/etc/nas-dashboard/mergerfs", 0755); err != nil {
			// 目录创建失败，继续运行
		}

		// 初始化MergerFS管理器
		mergerMgr = mergerfs.NewManager("/etc/nas-dashboard/mergerfs/config.json")

		// 加载现有配置
		if err := mergerMgr.LoadConfig(); err != nil {
			// 配置文件不存在或加载失败，使用默认配置
		}

		// 启动监控定时任务（仅在数据库可用时）
		if db != nil {
			monitorAPI := NewMonitorAPI(db)
			go monitorAPI.StartMonitoring()

			// 启动系统任务调度器
			scheduler := NewScheduler(db)
			scheduler.Start()

			// 应用备份管理种子配置（config/backup-seed.json），
			// 让新机器部署时自动获得预设的仓库 / 任务。
			go seedBackupConfig(db)
		}
	})
}

// GetStoragePoolAPI 获取存储池API实例
func GetStoragePoolAPI() *StoragePoolAPI {
	return NewStoragePoolAPI(dbInstance, mergerMgr)
}

// GetMonitorAPI 获取监控API实例
func GetMonitorAPI() *MonitorAPI {
	return NewMonitorAPI(dbInstance)
}

// GetQuotaAPI 获取配额API实例
func GetQuotaAPI() *QuotaAPI {
	return NewQuotaAPI(dbInstance)
}

// GetMergerFSManager 获取MergerFS管理器实例
func GetMergerFSManager() *mergerfs.Manager {
	return mergerMgr
}