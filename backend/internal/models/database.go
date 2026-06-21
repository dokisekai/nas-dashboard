package models

import (
	"time"
	"gorm.io/gorm"
)

// BaseModel 包含通用字段
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

// User 用户模型
type User struct {
	BaseModel
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"not null" json:"-"`
	Password     string `gorm:"-" json:"password,omitempty"` // 临时字段，用于API
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	NickName     string `json:"nickName,omitempty"`
	Role         string `gorm:"default:'user'" json:"role"` // admin, user
	IsActive     bool   `gorm:"default:true" json:"isActive"`
	Status       string `gorm:"default:'active'" json:"status"`
	Groups  []string `gorm:"serializer:json" json:"groups,omitempty"`
	LastLogin    *time.Time `json:"lastLogin,omitempty"`

	// SSO相关字段
	SSOProvider string `json:"ssoProvider,omitempty"` // SSO提供商名称（如google, github等）
	SSOID       string `json:"ssoId,omitempty"`       // SSO提供商的用户ID

	// 关联关系
	SSHKeys      []SSHKey `gorm:"foreignKey:UserID" json:"sshKeys,omitempty"`
	Sessions     []Session `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
}

// SSHKey SSH 密钥模型
type SSHKey struct {
	BaseModel
	UserID      uint   `gorm:"not null" json:"userId"`
	Name        string `gorm:"not null" json:"name"`
	Fingerprint string `gorm:"uniqueIndex" json:"fingerprint"`
	Type        string `json:"type"`
	Content     string `gorm:"not null" json:"content"`
	User        User   `gorm:"foreignKey:UserID" json:"-"`
}

// Session 会话模型
type Session struct {
	BaseModel
	UserID        uint   `gorm:"not null" json:"userId"`
	RefreshToken  string `gorm:"uniqueIndex;not null" json:"-"`
	IPAddress     string `json:"ipAddress"`
	UserAgent     string `json:"userAgent"`
	ExpiresAt     time.Time `gorm:"not null" json:"expiresAt"`
	LastRefreshAt time.Time `json:"lastRefreshAt"`
	IsActive      bool    `gorm:"default:true" json:"isActive"`
	User          User    `gorm:"foreignKey:UserID" json:"-"`
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	BaseModel
	Key         string `gorm:"uniqueIndex;not null" json:"key"`
	Value       string `gorm:"not null" json:"value"`
	Type        string `gorm:"default:'string'" json:"type"` // string, int, bool, json
	Category    string `gorm:"default:'general'" json:"category"`
	Description string `json:"description"`
	IsPublic    bool   `gorm:"default:false" json:"isPublic"`
}

// BackupRecord 备份记录模型
type BackupRecord struct {
	BaseModel
	Name         string    `gorm:"not null" json:"name"`
	Type         string    `json:"type"` // full, incremental, differential
	FilePath     string    `gorm:"not null" json:"filePath"`
	Size         int64     `json:"size"`
	Status       string    `gorm:"default:'completed'" json:"status"` // completed, failed, in_progress
	Groups  []string `gorm:"serializer:json" json:"groups,omitempty"`
	CreatedBy    string    `json:"createdBy"`
	Description  string    `json:"description"`
	CompletedAt  *time.Time `json:"completedAt"`
}

// OperationLog 操作日志模型
type OperationLog struct {
	BaseModel
	UserID      uint   `json:"userId"`
	User        User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action      string `gorm:"not null" json:"action"`
	Resource    string `json:"resource"`
	IPAddress   string `json:"ipAddress"`
	UserAgent   string `json:"userAgent"`
	Details     string `gorm:"type:text" json:"details"`
	Status      string `gorm:"default:'success'" json:"status"` // success, failed
	Groups  []string `gorm:"serializer:json" json:"groups,omitempty"`
}

// Plugin 插件模型
type Plugin struct {
	BaseModel
	Name          string `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName   string `json:"displayName"`
	Version       string `json:"version"`
	Author        string `json:"author"`
	Description   string `json:"description"`
	FilePath      string `json:"filePath"`
	Config        string `gorm:"type:json" json:"config"` // JSON配置
	IsActive      bool   `gorm:"default:false" json:"isActive"`
	IsInstalled   bool   `gorm:"default:false" json:"isInstalled"`
	InstalledAt   *time.Time `json:"installedAt"`
}

// FileSystemAccess 文件系统访问记录
type FileSystemAccess struct {
	BaseModel
	UserID    uint   `json:"userId"`
	User      User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Path      string `gorm:"not null" json:"path"`
	Operation string `gorm:"not null" json:"operation"` // read, write, delete, execute
	Success   bool   `gorm:"default:true" json:"success"`
	ErrorMsg  string `json:"errorMsg"`
}
