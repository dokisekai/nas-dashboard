package models

import (
	"time"
	"gorm.io/gorm"
)

// StoragePool 存储池模型
type StoragePool struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name        string         `gorm:"uniqueIndex;not null" json:"name"`
	Type        string         `json:"type"`          // mergerfs, btrfs, zfs, lvm
	Status      string         `json:"status"`        // active, inactive, error, degraded
	MountPoint  string         `gorm:"not null" json:"mountPoint"`
	TotalSize   uint64         `json:"totalSize"`
	UsedSize    uint64         `json:"usedSize"`
	FreeSize    uint64         `json:"freeSize"`
	Description string         `json:"description"`
	Config      string         `gorm:"type:json" json:"config"`
	CreatedBy   string         `json:"createdBy"`

	PoolDisks  []PoolDisk     `gorm:"foreignKey:PoolID" json:"poolDisks,omitempty"`
	Snapshots  []PoolSnapshot `gorm:"foreignKey:PoolID" json:"snapshots,omitempty"`
}

// PoolDisk 存储池磁盘关联
type PoolDisk struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	PoolID     uint   `gorm:"not null" json:"poolId"`
	Device     string `gorm:"not null" json:"device"`
	Size       uint64 `json:"size"`
	Status     string `json:"status"`     // active, failed, removed
	Priority   int    `json:"priority"`    // mergerfs 优先级
	BranchPath string `json:"branchPath"` // 挂载路径
}

// PoolSnapshot 存储池快照
type PoolSnapshot struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	PoolID      uint       `gorm:"not null" json:"poolId"`
	Name        string     `gorm:"not null" json:"name"`
	Description string     `json:"description"`
	Size        uint64     `json:"size"`
	Status      string     `json:"status"`    // creating, completed, deleting, error
	CompletedAt *time.Time `json:"completedAt,omitempty"`
}

// RAIDConfig RAID配置
type RAIDConfig struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name       string `json:"name"`
	Level      string `json:"level"`        // 0, 1, 5, 6, 10, 1+0, etc.
	Devices    string `gorm:"type:json" json:"devices"`  // JSON数组
	Status     string `json:"status"`       // active, degraded, failed, recovering
	Size       uint64 `json:"size"`
	UUID       string `json:"uuid"`
	MountPoint string `json:"mountPoint"`
}

// PhysicalVolume 物理卷
type PhysicalVolume struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Device string `json:"device"`
	VGName string `json:"vgName"`
	Size   uint64 `json:"size"`
	Free   uint64 `json:"free"`
	UUID   string `json:"uuid"`
	Status string `json:"status"` // active, missing, exported
}

// VolumeGroup 卷组
type VolumeGroup struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name    string `json:"name"`
	Size    uint64 `json:"size"`
	Free    uint64 `json:"free"`
	PVCount int    `json:"pvCount"`
	LVCount int    `json:"lvCount"`
	UUID    string `json:"uuid"`
	Status  string `json:"status"` // active, partial, exported
}

// LogicalVolume 逻辑卷
type LogicalVolume struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name       string `json:"name"`
	VGName     string `json:"vgName"`
	Size       uint64 `json:"size"`
	Path       string `json:"path"`
	UUID       string `json:"uuid"`
	MountPoint string `json:"mountPoint"`
	Status     string `json:"status"` // active, inactive, snapshot merge
}