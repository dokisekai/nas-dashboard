package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"nas-dashboard/internal/models"
	"nas-dashboard/pkg/mergerfs"
	"gorm.io/gorm"
)

// StoragePoolAPI 存储池API处理器
type StoragePoolAPI struct {
	DB          *gorm.DB
	MergerMgr   *mergerfs.Manager
}

// NewStoragePoolAPI 创建存储池API
func NewStoragePoolAPI(db *gorm.DB, mergerMgr *mergerfs.Manager) *StoragePoolAPI {
	return &StoragePoolAPI{
		DB:        db,
		MergerMgr: mergerMgr,
	}
}

// CreatePoolRequest 创建存储池请求
type CreatePoolRequest struct {
	Name        string                       `json:"name" binding:"required"`
	Type        string                       `json:"type" binding:"required"` // mergerfs, btrfs, zfs, lvm
	MountPoint  string                       `json:"mountPoint" binding:"required"`
	Description string                       `json:"description"`
	Branches    []mergerfs.BranchConfig      `json:"branches" binding:"required,min=1"`
	Config      *mergerfs.Config             `json:"config"`
}

// UpdatePoolRequest 更新存储池请求
type UpdatePoolRequest struct {
	Description string              `json:"description"`
	Config      *mergerfs.Config    `json:"config"`
}

// AddDiskRequest 添加磁盘请求
type AddDiskRequest struct {
	Path     string `json:"path" binding:"required"`
	Mode     string `json:"mode" binding:"required,oneof=rw ro"`
	Priority int    `json:"priority"`
}

// CreatePool 创建存储池
func (api *StoragePoolAPI) CreatePool(c *gin.Context) {
	var req CreatePoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证存储池名称
	if strings.ContainsAny(req.Name, "/\\ \t\n") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pool name"})
		return
	}

	// 验证类型
	if req.Type != "mergerfs" && req.Type != "btrfs" && req.Type != "zfs" && req.Type != "lvm" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported pool type"})
		return
	}

	// 创建存储池记录
	pool := &models.StoragePool{
		Name:        req.Name,
		Type:        req.Type,
		MountPoint:  req.MountPoint,
		Description: req.Description,
		Status:      "creating",
	}

	// 保存配置
	if req.Config != nil {
		configJSON, err := json.Marshal(req.Config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal config"})
			return
		}
		pool.Config = string(configJSON)
	}

	// 根据类型创建存储池
	switch req.Type {
	case "mergerfs":
		// 使用MergerFS管理器创建
		if err := api.MergerMgr.CreatePool(req.Name, req.MountPoint, req.Branches, req.Config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create mergerfs pool: %v", err)})
			return
		}

		// 创建磁盘关联记录
		pool.PoolDisks = make([]models.PoolDisk, len(req.Branches))
		for i, branch := range req.Branches {
			pool.PoolDisks[i] = models.PoolDisk{
				Device:     branch.Path,
				BranchPath: branch.Path,
				Status:     "active",
				Priority:   branch.Priority,
			}
		}
		pool.Status = "active"

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}

	// 保存到数据库
	if result := api.DB.Create(pool); result.Error != nil {
		// 回滚：删除已创建的存储池
		api.MergerMgr.DeletePool(req.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, pool)
}

// GetPools 获取存储池列表
func (api *StoragePoolAPI) GetPools(c *gin.Context) {
	var pools []models.StoragePool
	if result := api.DB.Preload("PoolDisks").Find(&pools); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pools": pools})
}

// GetPool 获取存储池详情
func (api *StoragePoolAPI) GetPool(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Preload("PoolDisks").Preload("Snapshots").Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 获取实时状态
	if pool.Type == "mergerfs" {
		if status, err := api.MergerMgr.GetPoolStatus(name); err == nil {
			// 更新状态信息
			pool.TotalSize = 0
			pool.UsedSize = 0
			pool.FreeSize = 0

			for _, branch := range status.Branches {
				pool.TotalSize += branch.Size
				pool.UsedSize += branch.Used
				pool.FreeSize += branch.Free
			}

			// 更新磁盘状态
			for i, disk := range pool.PoolDisks {
				if i < len(status.Branches) {
					disk.Status = status.Branches[i].Mode
					disk.Size = status.Branches[i].Size
				}
			}
		}
	}

	c.JSON(http.StatusOK, pool)
}

// UpdatePool 更新存储池
func (api *StoragePoolAPI) UpdatePool(c *gin.Context) {
	name := c.Param("name")

	var req UpdatePoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 更新描述
	if req.Description != "" {
		pool.Description = req.Description
	}

	// 更新配置
	if req.Config != nil {
		configJSON, err := json.Marshal(req.Config)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal config"})
			return
		}
		pool.Config = string(configJSON)
	}

	if result := api.DB.Save(&pool); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, pool)
}

// DeletePool 删除存储池
func (api *StoragePoolAPI) DeletePool(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 检查存储池状态
	if pool.Status == "active" || pool.Status == "degraded" {
		// 先尝试卸载
		if pool.Type == "mergerfs" {
			if err := api.MergerMgr.DeletePool(name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete pool: %v", err)})
				return
			}
		}
	}

	// 从数据库删除
	if result := api.DB.Delete(&pool); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pool deleted successfully"})
}

// AddDisk 添加磁盘到存储池
func (api *StoragePoolAPI) AddDisk(c *gin.Context) {
	name := c.Param("name")

	var req AddDiskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 根据类型添加磁盘
	switch pool.Type {
	case "mergerfs":
		// 使用MergerFS管理器添加磁盘
		if err := api.MergerMgr.AddDisk(name, req.Path, req.Mode, req.Priority); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add disk: %v", err)})
			return
		}

		// 创建磁盘关联记录
		poolDisk := models.PoolDisk{
			PoolID:     pool.ID,
			Device:     req.Path,
			BranchPath: req.Path,
			Status:     "active",
			Priority:   req.Priority,
		}

		if result := api.DB.Create(&poolDisk); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Disk added successfully"})
}

// RemoveDisk 从存储池移除磁盘
func (api *StoragePoolAPI) RemoveDisk(c *gin.Context) {
	name := c.Param("name")
	device := c.Param("device")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 查找磁盘记录
	var poolDisk models.PoolDisk
	if result := api.DB.Where("pool_id = ? AND device = ?", pool.ID, device).First(&poolDisk); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Disk not found in pool"})
		return
	}

	// 根据类型移除磁盘
	switch pool.Type {
	case "mergerfs":
		// 使用MergerFS管理器移除磁盘
		if err := api.MergerMgr.RemoveDisk(name, device); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove disk: %v", err)})
			return
		}

		// 从数据库删除记录
		if result := api.DB.Delete(&poolDisk); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Disk removed successfully"})
}

// GetPoolBranches 获取存储池分支信息
func (api *StoragePoolAPI) GetPoolBranches(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	if pool.Type == "mergerfs" {
		if status, err := api.MergerMgr.GetPoolStatus(name); err == nil {
			c.JSON(http.StatusOK, gin.H{"branches": status.Branches})
			return
		}
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branches"})
}

// MountPool 挂载存储池
func (api *StoragePoolAPI) MountPool(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 检查是否已挂载
	if pool.Status == "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool is already mounted"})
		return
	}

	// 根据类型挂载
	switch pool.Type {
	case "mergerfs":
		// 加载配置
		var config mergerfs.Config
		if pool.Config != "" {
			if err := json.Unmarshal([]byte(pool.Config), &config); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse config"})
				return
			}
		}

		// 获取磁盘关联
		var poolDisks []models.PoolDisk
		if result := api.DB.Where("pool_id = ?", pool.ID).Find(&poolDisks); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		// 构建分支配置
		branches := make([]mergerfs.BranchConfig, len(poolDisks))
		for i, disk := range poolDisks {
			branches[i] = mergerfs.BranchConfig{
				Path:     disk.BranchPath,
				Mode:     "rw", // 默认读写模式
				Priority: disk.Priority,
			}
		}

		// 创建存储池
		if err := api.MergerMgr.CreatePool(name, pool.MountPoint, branches, &config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to mount pool: %v", err)})
			return
		}

		// 更新状态
		pool.Status = "active"
		if result := api.DB.Save(&pool); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pool mounted successfully"})
}

// UmountPool 卸载存储池
func (api *StoragePoolAPI) UmountPool(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 检查状态
	if pool.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool is not mounted"})
		return
	}

	// 根据类型卸载
	switch pool.Type {
	case "mergerfs":
		if err := api.MergerMgr.DeletePool(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to unmount pool: %v", err)})
			return
		}

		// 更新状态
		pool.Status = "inactive"
		if result := api.DB.Save(&pool); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pool unmounted successfully"})
}

// BalancePool 平衡存储池
func (api *StoragePoolAPI) BalancePool(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 根据类型平衡
	switch pool.Type {
	case "mergerfs":
		if err := api.MergerMgr.BalancePool(name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to balance pool: %v", err)})
			return
		}

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pool balanced successfully"})
}

// ScanPool 扫描存储池状态
func (api *StoragePoolAPI) ScanPool(c *gin.Context) {
	name := c.Param("name")

	var pool models.StoragePool
	if result := api.DB.Where("name = ?", name).First(&pool); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pool not found"})
		return
	}

	// 根据类型扫描
	switch pool.Type {
	case "mergerfs":
		status, err := api.MergerMgr.GetPoolStatus(name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to scan pool: %v", err)})
			return
		}

		// 更新数据库中的状态
		pool.TotalSize = 0
		pool.UsedSize = 0
		pool.FreeSize = 0

		for i, branch := range status.Branches {
			pool.TotalSize += branch.Size
			pool.UsedSize += branch.Used
			pool.FreeSize += branch.Free

			// 更新磁盘状态
			if i < len(pool.PoolDisks) {
				pool.PoolDisks[i].Status = branch.Mode
				pool.PoolDisks[i].Size = branch.Size
			}
		}

		// 确定存储池状态
		if pool.UsedSize > uint64(float64(pool.TotalSize) * 0.9) {
			pool.Status = "warning"
		} else {
			pool.Status = "active"
		}

		if result := api.DB.Save(&pool); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, status)

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pool type not implemented yet"})
		return
	}
}