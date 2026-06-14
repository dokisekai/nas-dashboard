package api

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"nas-dashboard/pkg/power"
)

// GetPowerCurrent 获取当前功耗
func GetPowerCurrent(c *gin.Context) {
	powerData, err := power.DirectReadPower()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, powerData)
}

// GetPowerHistory 获取历史功耗数据
func GetPowerHistory(c *gin.Context) {
	// 获取查询参数
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 && d <= 90 {
			days = d
		}
	}

	pm := power.NewPowerMonitor()
	data, err := pm.GetHistoricalPower(days)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"data": data,
		"count": len(data),
		"days":  days,
	})
}

// GetPowerStatistics 获取功耗统计信息
func GetPowerStatistics(c *gin.Context) {
	days := 7
	if daysParam := c.Query("days"); daysParam != "" {
		if d, err := strconv.Atoi(daysParam); err == nil && d > 0 && d <= 90 {
			days = d
		}
	}

	pm := power.NewPowerMonitor()
	stats, err := pm.GetPowerStatistics(days)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, stats)
}

// GetPowerOverview 获取功耗概览
func GetPowerOverview(c *gin.Context) {
	// 获取当前功耗
	current, err := power.DirectReadPower()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 获取统计数据
	pm := power.NewPowerMonitor()
	stats, err := pm.GetPowerStatistics(1)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 检查告警
	thresholds := map[string]float64{
		"high":      150.0,
		"critical":  200.0,
		"gpu":       250.0,
	}
	alerts := power.CheckPowerAlert(current, thresholds)

	c.JSON(200, gin.H{
		"current":  current,
		"today":    stats,
		"alerts":   alerts,
		"timestamp": time.Now(),
	})
}
