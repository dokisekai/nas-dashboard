package models

import (
	"time"
	"gorm.io/gorm"
)

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID          int32   `json:"pid"`
	Name         string  `json:"name"`
	Status       string  `json:"status"`       // running, sleeping, stopped, zombie
	CPUPercent   float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
	Memory       uint64  `json:"memory"`
	Threads      int32   `json:"threads"`
	Username     string  `json:"username"`
	Command      string  `json:"command"`
	CreatedTime  int64   `json:"createdTime"`
}

// ServiceInfo 系统服务信息
type ServiceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`     // running, stopped, failed, masked
	Enabled     bool   `json:"enabled"`
	LoadState   string `json:"loadState"`   // loaded, masked, error, etc.
	ActiveState string `json:"activeState"` // active, inactive, failed, etc.
	MainPID     int32  `json:"mainPid"`
	SubState    string `json:"subState"`    // running, dead, exited, etc.
}

// TemperatureInfo 温度信息
type TemperatureInfo struct {
	Sensors []Sensor `json:"sensors"`
}

// Sensor 温度传感器
type Sensor struct {
	Name     string  `json:"name"`
	Current  float64 `json:"current"`
	Max      float64 `json:"max"`
	Critical float64 `json:"critical"`
	Unit     string  `json:"unit"` // C, F
}

// SystemEvent 系统事件
type SystemEvent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Type      string `json:"type"`       // info, warning, error, critical
	Source    string `json:"source"`     // system, storage, network, service, user
	Title     string `json:"title"`
	Message   string `json:"message"`
	Details   string `gorm:"type:text" json:"details"`
	UserID    uint   `json:"userId,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	Resolved  bool   `json:"resolved"`
}

// AlertRule 告警规则
type AlertRule struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name      string `json:"name"`
	Type      string `json:"type"`        // cpu, memory, disk, temperature, network, service
	Condition string `json:"condition"`   // >, <, =, !=, >=, <=
	Threshold float64 `json:"threshold"`
	Duration  int    `json:"duration"`    // 持续时间（秒）
	Cooldown  int    `json:"cooldown"`    // 冷却时间（秒）
	Severity  string `json:"severity"`    // info, warning, critical
	Enabled   bool   `json:"enabled"`
	Actions   string `gorm:"type:json" json:"actions"`     // JSON格式的动作列表
	LastTriggered *time.Time `json:"lastTriggered,omitempty"`
}

// SystemLog 系统日志
type SystemLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Timestamp   time.Time `json:"timestamp"`
	Level       string    `json:"level"`       // debug, info, warning, error, critical
	Component   string    `json:"component"`   // system, storage, network, auth, etc.
	Message     string    `json:"message"`
	Details     string    `gorm:"type:text" json:"details"`
	IPAddress   string    `json:"ipAddress,omitempty"`
	UserAgent   string    `json:"userAgent,omitempty"`
	UserID      uint      `json:"userId,omitempty"`
}

// MonitorHistory 监控历史数据
type MonitorHistory struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Type      string  `json:"type"`       // cpu, memory, disk, network
	Metric    string  `json:"metric"`     // usage, temperature, speed, etc.
	Value     float64 `json:"value"`
	Labels    string  `gorm:"type:json" json:"labels"` // JSON格式的标签
	Source    string  `json:"source"`     // cpu0, eth0, sda1, etc.
}