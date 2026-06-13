package models

import (
	"time"
	"gorm.io/gorm"
)

// UserQuota 用户配额
type UserQuota struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	UserID      uint   `gorm:"not null" json:"userId"`
	Path        string `gorm:"not null" json:"path"`
	SoftLimit   uint64 `json:"softLimit"`   // 软限制（字节）
	HardLimit   uint64 `json:"hardLimit"`   // 硬限制（字节）
	UsedSpace   uint64 `json:"usedSpace"`   // 已用空间（字节）
	GracePeriod int    `json:"gracePeriod"` // 宽限期（天）
	FilesUsed   uint64 `json:"filesUsed"`   // 已用文件数
	FilesSoft   uint64 `json:"filesSoft"`   // 文件数软限制
	FilesHard   uint64 `json:"filesHard"`   // 文件数硬限制

	// 关联用户
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// GroupQuota 组配额
type GroupQuota struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	GroupID     uint   `gorm:"not null" json:"groupId"`
	Path        string `gorm:"not null" json:"path"`
	SoftLimit   uint64 `json:"softLimit"`   // 软限制（字节）
	HardLimit   uint64 `json:"hardLimit"`   // 硬限制（字节）
	UsedSpace   uint64 `json:"usedSpace"`   // 已用空间（字节）
	GracePeriod int    `json:"gracePeriod"` // 宽限期（天）
	FilesUsed   uint64 `json:"filesUsed"`   // 已用文件数
	FilesSoft   uint64 `json:"filesSoft"`   // 文件数软限制
	FilesHard   uint64 `json:"filesHard"`   // 文件数硬限制
}

// QuotaReport 配额报告
type QuotaReport struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name       string    `json:"name"`       // 用户名或组名
	Type       string    `json:"type"`       // user, group
	Path       string    `json:"path"`
	UsedSpace  uint64    `json:"usedSpace"`
	SoftLimit  uint64    `json:"softLimit"`
	HardLimit  uint64    `json:"hardLimit"`
	UsedPercent float64  `json:"usedPercent"` // 使用率百分比
	Status     string    `json:"status"`     // ok, warning, exceeded, grace
	FilesUsed  uint64    `json:"filesUsed"`
	FilesSoft  uint64    `json:"filesSoft"`
	FilesHard  uint64    `json:"filesHard"`
	GeneratedAt time.Time `json:"generatedAt"`
}

// QuotaAlert 配额告警
type QuotaAlert struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	UserID     uint   `gorm:"not null" json:"userId"`
	GroupID    uint   `json:"groupId,omitempty"`
	Type       string `json:"type"`       // user, group
	Path       string `json:"path"`
	AlertType  string `json:"alertType"`  // soft_limit, hard_limit, grace_period
	Severity   string `json:"severity"`   // warning, critical
	Message    string `json:"message"`
	Resolved   bool   `json:"resolved"`
	ResolvedAt *time.Time `json:"resolvedAt,omitempty"`
}