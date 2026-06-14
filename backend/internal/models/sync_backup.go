package models

import (
	"time"
	"gorm.io/gorm"
)

// SyncJob 数据同步任务模型
type SyncJob struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
	
	// 同步源和目标
	SourcePath  string `gorm:"not null" json:"sourcePath"`
	DestPath    string `gorm:"not null" json:"destPath"`
	
	// 同步类型: local, remote (rclone), rsync
	Type        string `json:"type"` 
	
	// 计划配置
	Enabled     bool   `json:"enabled"`
	Schedule    string `json:"schedule"` // cron 表达式
	
	// 策略
	DeleteExtra bool   `json:"deleteExtra"` // --delete
	Checksum    bool   `json:"checksum"`    // --checksum
	
	// 状态
	Status      string     `json:"status"` // idle, running, failed, completed
	LastRun     *time.Time `json:"lastRun"`
	LastError   string     `json:"lastError"`
}

// BackupRepo Restic 备份仓库模型
type BackupRepo struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	Type        string `json:"type"` // local, s3, sftp, webdav
	URL         string `json:"url"`
	
	// 凭据
	Username    string `json:"username"`
	Password    string `json:"password"` // 仓库密码 (加密存储)
	SecretKey   string `json:"-"`        // 云端密钥
	AccessKey   string `json:"-"`
	
	Status      string     `json:"status"` // active, error
	LastBackup  *time.Time `json:"lastBackup"`
}

// BackupTask Restic 备份任务
type BackupTask struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name        string `gorm:"not null" json:"name"`
	RepoID      uint   `json:"repoId"`
	Repo        BackupRepo `gorm:"foreignKey:RepoID" json:"repo"`
	
	SourcePath  string `gorm:"not null" json:"sourcePath"`
	Excludes    string `json:"excludes"` // 逗号分隔
	
	Enabled     bool   `json:"enabled"`
	Schedule    string `json:"schedule"`
	
	Status      string     `json:"status"`
	LastRun     *time.Time `json:"lastRun"`
}
