package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

// UPSStatus UPS 状态信息
type UPSStatus struct {
	Model       string  `json:"model"`
	Status      string  `json:"status"` // OL, OB, LB, etc.
	Charge      float64 `json:"charge"` // 电池电量 %
	Load        float64 `json:"load"`   // 负载 %
	Runtime     int     `json:"runtime"` // 剩余时间 (秒)
	InputVoltage float64 `json:"inputVoltage"`
	OutputVoltage float64 `json:"outputVoltage"`
}

// GetUPSStatus 获取 UPS 状态
func GetUPSStatus(c *gin.Context) {
	upsName := c.DefaultQuery("ups", "ups")
	
	cmd := exec.Command("upsc", upsName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stdout pipe: " + err.Error()})
		return
	}

	if err := cmd.Start(); err != nil {
		// 如果 upsc 不存在或失败，可能是没有安装 NUT 或没有连接 UPS
		c.JSON(http.StatusNotFound, gin.H{"error": "UPS not found or NUT not installed", "details": err.Error()})
		return
	}

	status := UPSStatus{
		Status: "Unknown",
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "device.model":
			status.Model = value
		case "ups.status":
			status.Status = value
		case "battery.charge":
			fmt.Sscanf(value, "%f", &status.Charge)
		case "ups.load":
			fmt.Sscanf(value, "%f", &status.Load)
		case "battery.runtime":
			fmt.Sscanf(value, "%d", &status.Runtime)
		case "input.voltage":
			fmt.Sscanf(value, "%f", &status.InputVoltage)
		case "output.voltage":
			fmt.Sscanf(value, "%f", &status.OutputVoltage)
		}
	}

	cmd.Wait()

	if status.Model == "" && status.Status == "Unknown" {
		// 如果没有任何输出，可能是模拟环境或错误
		// 这里为了演示，如果没有 UPS，我们可以返回一个模拟数据或错误
		// c.JSON(http.StatusNotFound, gin.H{"error": "No data returned from upsc"})
		// return
		
		// 模拟数据用于演示 (如果需要可以注释掉)
		/*
		status = UPSStatus{
			Model: "Simulated UPS",
			Status: "OL",
			Charge: 95.0,
			Load: 15.0,
			Runtime: 3600,
			InputVoltage: 230.5,
			OutputVoltage: 230.1,
		}
		*/
	}

	c.JSON(http.StatusOK, status)
}
