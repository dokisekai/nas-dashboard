package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"nas-dashboard/pkg/system"
)

// GetPhysicalVolumes 获取 PV 列表
func GetPhysicalVolumes(c *gin.Context) {
	pvs, err := system.GetPhysicalVolumes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get PVs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pvs": pvs})
}

// CreatePhysicalVolume 创建 PV
func CreatePhysicalVolume(c *gin.Context) {
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
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot use system disk as PV"})
		return
	}
	if IsDeviceMounted(req.Device) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Device is mounted"})
		return
	}

	err := system.CreatePhysicalVolume(req.Device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create PV: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PV created successfully"})
}

// GetVolumeGroups 获取 VG 列表
func GetVolumeGroups(c *gin.Context) {
	vgs, err := system.GetVolumeGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get VGs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"vgs": vgs})
}

// CreateVolumeGroup 创建 VG
func CreateVolumeGroup(c *gin.Context) {
	var req struct {
		Name    string   `json:"name" binding:"required"`
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
		// PV 通常是已经通过 CreatePhysicalVolume 创建的，
		// 如果是已经挂载的，说明它不是一个自由的 PV
		if IsDeviceMounted(dev) {
			c.JSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("Device %s is mounted", dev)})
			return
		}
	}

	err := system.CreateVolumeGroup(req.Name, req.Devices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create VG: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "VG created successfully"})
}

// DeleteVolumeGroup 删除 VG
func DeleteVolumeGroup(c *gin.Context) {
	name := c.Param("name")
	err := system.DeleteVolumeGroup(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete VG: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "VG deleted successfully"})
}

// GetLogicalVolumes 获取 LV 列表
func GetLogicalVolumes(c *gin.Context) {
	vgName := c.Query("vg")
	lvs, err := system.GetLogicalVolumes(vgName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get LVs: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lvs": lvs})
}

// CreateLogicalVolume 创建 LV
func CreateLogicalVolume(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required"`
		VGName string `json:"vgName" binding:"required"`
		Size   uint64 `json:"size" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := system.CreateLogicalVolume(req.VGName, req.Name, req.Size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create LV: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "LV created successfully"})
}

// DeleteLogicalVolume 删除 LV
func DeleteLogicalVolume(c *gin.Context) {
	vgName := c.Param("vg")
	name := c.Param("name")
	err := system.DeleteLogicalVolume(vgName, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete LV: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "LV deleted successfully"})
}
