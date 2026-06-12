package api

import (
	"net/http"
	"nas-dashboard/pkg/system"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetCPU 获取 CPU 信息
func GetCPU(c *gin.Context) {
	cpu, err := system.GetCPUInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, cpu)
}

// GetMemory 获取内存信息
func GetMemory(c *gin.Context) {
	mem, err := system.GetMemoryInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, mem)
}

// GetDisk 获取磁盘信息
func GetDisk(c *gin.Context) {
	diskInfo, err := system.GetDiskInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, diskInfo)
}

// GetNetwork 获取网络信息
func GetNetwork(c *gin.Context) {
	networkInfo, err := system.GetNetworkInfo()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, networkInfo)
}

// WebSocket 升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应限制
	},
}

// WSMonitor WebSocket 监控数据推送
func WSMonitor(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// 定时推送监控数据
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取所有监控数据
			cpu, _ := system.GetCPUInfo()
			mem, _ := system.GetMemoryInfo()
			diskInfo, _ := system.GetDiskInfo()
			networkInfo, _ := system.GetNetworkInfo()

			data := gin.H{
				"cpu":     cpu,
				"memory":  mem,
				"disk":    diskInfo,
				"network": networkInfo,
				"time":    time.Now().Unix(),
			}

			if err := conn.WriteJSON(data); err != nil {
				return
			}
		}
	}
}
