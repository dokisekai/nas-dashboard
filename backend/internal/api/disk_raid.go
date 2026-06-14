package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"nas-dashboard/pkg/system"
)

// GetRAIDArrays 获取 RAID 阵列列表
func GetRAIDArrays(c *gin.Context) {
	arrays, err := system.GetRAIDArrays()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get RAID arrays: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"arrays": arrays})
}

// GetRAIDArray 获取单个 RAID 阵列详情
func GetRAIDArray(c *gin.Context) {
	name := c.Param("name")
	array, err := system.GetRAIDDetail("/dev/" + name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("RAID array not found: %v", err)})
		return
	}
	c.JSON(http.StatusOK, array)
}

// CreateRAID 创建 RAID 阵列
func CreateRAID(c *gin.Context) {
	var req struct {
		Name    string   `json:"name" binding:"required"`
		Level   string   `json:"level" binding:"required"`
		Devices []string `json:"devices" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 安全检查
	rootDevice := GetRootDevice()
	for _, dev := range req.Devices {
		if dev == rootDevice || strings.HasPrefix(rootDevice, dev) {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Device %s is a system disk", dev)})
			return
		}
		if IsDeviceMounted(dev) {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Device %s is mounted", dev)})
			return
		}
	}

	err := system.CreateRAID(req.Name, req.Level, req.Devices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create RAID: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "RAID array created successfully"})
}

// DeleteRAID 删除 RAID 阵列
func DeleteRAID(c *gin.Context) {
	name := c.Param("name")
	err := system.DeleteRAID(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete RAID: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "RAID array deleted successfully"})
}

// AddDiskToRAID 向 RAID 阵列添加磁盘
func AddDiskToRAID(c *gin.Context) {
	name := c.Param("name")
	var req struct {
		Device string `json:"device" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 安全检查
	rootDevice := GetRootDevice()
	if req.Device == rootDevice || strings.HasPrefix(rootDevice, req.Device) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot add system disk to RAID"})
		return
	}
	if IsDeviceMounted(req.Device) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Device is mounted"})
		return
	}

	err := system.AddDiskToRAID(name, req.Device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add disk: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disk added successfully"})
}

// RemoveDiskFromRAID 从 RAID 阵列移除磁盘
func RemoveDiskFromRAID(c *gin.Context) {
	name := c.Param("name")
	var req struct {
		Device string `json:"device" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := system.RemoveDiskFromRAID(name, req.Device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove disk: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Disk removed successfully"})
}
