package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"nas-dashboard/pkg/system"
)

// RunDiskBenchmark 运行磁盘性能测试
func RunDiskBenchmark(c *gin.Context) {
	device := c.Param("device")
	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is required"})
		return
	}

	result, err := system.RunBenchmark(device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to run benchmark: %v", err)})
		return
	}

	c.JSON(http.StatusOK, result)
}
