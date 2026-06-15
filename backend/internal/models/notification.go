package models

import "time"

// Notification 系统通知模型
type Notification struct {
	BaseModel
	Type       string    `gorm:"default:'info'" json:"type"` // info, warning, error, success
	Title      string    `gorm:"not null" json:"title"`
	Message    string    `gorm:"type:text" json:"message"`
	Read       bool      `gorm:"default:false" json:"read"`
	Timestamp  time.Time `json:"timestamp"`
	Persistent bool      `gorm:"default:false" json:"persistent"` // 是否持久化显示
	UserID     *uint     `json:"userId,omitempty" gorm:"index"` // 关联用户ID（可选，全局通知为null）

	// 关联
	Rules []NotificationRule `json:"rules,omitempty" gorm:"foreignKey:NotificationID"`
}

// NotificationRule 通知规则模型
type NotificationRule struct {
	BaseModel
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	EventType    string    `gorm:"type:varchar(50);not null;index" json:"eventType"`
	Conditions   string    `gorm:"type:text" json:"conditions"` // JSON格式的条件
	Actions      string    `gorm:"type:text" json:"actions"` // JSON格式的动作
	Enabled      bool      `gorm:"default:true" json:"enabled"`
	Cooldown     int       `gorm:"default:300" json:"cooldown"` // 冷却时间（秒）
	LastTriggered time.Time `json:"lastTriggered"`
}

// NotificationAction 通知动作
type NotificationAction struct {
	BaseModel
	RuleID uint   `gorm:"not null;index" json:"ruleId"`
	Type   string `gorm:"type:varchar(20);not null" json:"type"` // websocket, email, webhook
	Config string `gorm:"type:text" json:"config"` // JSON格式的配置
}

// SystemEvent 系统事件（用于内部处理）
type SystemEvent struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Severity  string                 `json:"severity"` // info, warning, critical
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// Disk 磁盘信息模型（用于监控）
type Disk struct {
	Device      string  `json:"device"`
	Name        string  `json:"name"`
	Size        uint64  `json:"size"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	MountPoint  string  `json:"mountPoint"`
	FileSystem  string  `json:"fileSystem"`
	SmartStatus string  `json:"smartStatus"` // healthy, failing, unknown
	Temperature int     `json:"temperature"`
	Health      string  `json:"health"`
	Model       string  `json:"model"`
	Serial      string  `json:"serial"`
}

// MemoryInfo 内存信息模型
type MemoryInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"` // 使用百分比
}

// CPUInfo CPU信息模型
type CPUInfo struct {
	ModelName string  `json:"modelName"`
	Cores     int     `json:"cores"`
	Usage     float64 `json:"usage"` // 使用百分比
	MHz       float64 `json:"mhz"`
}

// LoadInfo 负载信息模型
type LoadInfo struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	BaseModel
	Username  string `gorm:"type:varchar(100);not null;index" json:"username"`
	IP        string `gorm:"type:varchar(50);not null;index" json:"ip"`
	UserAgent string `gorm:"type:text" json:"userAgent"`
	Success   bool   `gorm:"default:false" json:"success"`
	Timestamp time.Time `json:"timestamp" gorm:"autoCreateTime;index"`
}
