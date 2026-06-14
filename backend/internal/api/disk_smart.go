package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"nas-dashboard/pkg/system"
)

// GetDiskSmart 获取磁盘 SMART 信息
func GetDiskSmart(c *gin.Context) {
	device := c.Param("device")
	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is required"})
		return
	}

	info, err := system.GetSMARTInfo(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get SMART info: %v", err)})
		return
	}

	c.JSON(http.StatusOK, info)
}

// RunDiskSmartTest 运行磁盘 SMART 测试
func RunDiskSmartTest(c *gin.Context) {
	device := c.Param("device")
	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is required"})
		return
	}

	var req struct {
		Type string `json:"type" binding:"required"` // short, long, conveyance
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := system.RunSMARTTest(device, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start SMART test: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SMART test started successfully"})
}

// GetDiskHealth 获取磁盘健康状态
func GetDiskHealth(c *gin.Context) {
	device := c.Param("device")
	info, err := system.GetSMARTInfo(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get health info: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"health": info.OverallHealth, "temperature": info.Temperature})
}
