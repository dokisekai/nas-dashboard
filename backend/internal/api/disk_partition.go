package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"nas-dashboard/pkg/system"
)

// GetDiskPartitions 获取磁盘分区列表
func GetDiskPartitions(c *gin.Context) {
	device := c.Param("device")
	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is required"})
		return
	}

	table, err := system.GetPartitionTable(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get partitions: %v", err)})
		return
	}

	c.JSON(http.StatusOK, table)
}

// CreateDiskPartition 创建新分区
func CreateDiskPartition(c *gin.Context) {
	device := c.Param("device")
	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is required"})
		return
	}

	var req struct {
		Start      uint64 `json:"start"`
		End        uint64 `json:"end"`
		Type       string `json:"type"`
		Filesystem string `json:"filesystem"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := system.CreatePartition(device, req.Start, req.End, req.Type, req.Filesystem)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create partition: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partition created successfully"})
}

// DeleteDiskPartition 删除分区
func DeleteDiskPartition(c *gin.Context) {
	device := c.Param("device")
	numberStr := c.Param("number")
	
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid partition number"})
		return
	}

	err = system.DeletePartition(device, number)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete partition: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partition deleted successfully"})
}
