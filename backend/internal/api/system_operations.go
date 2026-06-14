package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 系统操作的互斥锁和状态管理
var (
	systemMutex sync.Mutex
	shutdownScheduled bool
	scheduledAction string
	scheduledTime time.Time
)

// RestartSystem 重启系统
func RestartSystem(c *gin.Context) {
	// 检查权限（在实际应用中应该检查管理员权限）
	// 这里为了演示目的，直接执行

	// 异步执行重启命令
	go func() {
		time.Sleep(2 * time.Second) // 给API响应时间
		exec.Command("shutdown", "-r", "+1").Run()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "系统将在1分钟后重启",
		"action":  "restart",
		"delay":   "1 minute",
	})
}

// ShutdownSystem 关闭系统
func ShutdownSystem(c *gin.Context) {
	// 检查权限（在实际应用中应该检查管理员权限）
	// 这里为了演示目的，直接执行

	// 异步执行关机命令
	go func() {
		time.Sleep(2 * time.Second) // 给API响应时间
		exec.Command("shutdown", "+1").Run()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "系统将在1分钟后关闭",
		"action":  "shutdown",
		"delay":   "1 minute",
	})
}

// CancelShutdown 取消计划的关机或重启
func CancelShutdown(c *gin.Context) {
	systemMutex.Lock()
	defer systemMutex.Unlock()

	if !shutdownScheduled {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "没有计划的系统操作需要取消",
		})
		return
	}

	// 取消计划的关机
	go func() {
		time.Sleep(500 * time.Millisecond)
		exec.Command("shutdown", "-c").Run()
	}()

	shutdownScheduled = false
	scheduledAction = ""
	scheduledTime = time.Time{}

	c.JSON(http.StatusOK, gin.H{
		"message": "已取消计划的系统操作",
		"action":  "cancelled",
	})
}

// ScheduleShutdown 计划关机
func ScheduleShutdown(c *gin.Context) {
	var req struct {
		DelayMinutes int `json:"delayMinutes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的延迟时间",
		})
		return
	}

	systemMutex.Lock()
	defer systemMutex.Unlock()

	if shutdownScheduled {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "已有计划的操作，请先取消",
		})
		return
	}

	// 计划关机
	go func() {
		time.Sleep(500 * time.Millisecond)
		exec.Command("shutdown", "+", fmt.Sprintf("%d", req.DelayMinutes)).Run()
	}()

	shutdownScheduled = true
	scheduledAction = "shutdown"
	scheduledTime = time.Now().Add(time.Duration(req.DelayMinutes) * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"message": "系统已计划关机",
		"action": "scheduled_shutdown",
		"delay_minutes": req.DelayMinutes,
		"scheduled_time": scheduledTime.Format("2006-01-02 15:04:05"),
	})
}

// ScheduleRestart 计划重启
func ScheduleRestart(c *gin.Context) {
	var req struct {
		DelayMinutes int `json:"delayMinutes" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的延迟时间",
		})
		return
	}

	systemMutex.Lock()
	defer systemMutex.Unlock()

	if shutdownScheduled {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "已有计划的操作，请先取消",
		})
		return
	}

	// 计划重启
	go func() {
		time.Sleep(500 * time.Millisecond)
		exec.Command("shutdown", "-r", "+", fmt.Sprintf("%d", req.DelayMinutes)).Run()
	}()

	shutdownScheduled = true
	scheduledAction = "restart"
	scheduledTime = time.Now().Add(time.Duration(req.DelayMinutes) * time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"message": "系统已计划重启",
		"action": "scheduled_restart",
		"delay_minutes": req.DelayMinutes,
		"scheduled_time": scheduledTime.Format("2006-01-02 15:04:05"),
	})
}

// GetShutdownStatus 获取计划操作的状态
func GetShutdownStatus(c *gin.Context) {
	systemMutex.Lock()
	defer systemMutex.Unlock()

	if !shutdownScheduled {
		c.JSON(http.StatusOK, gin.H{
			"scheduled": false,
			"action": nil,
			"time": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scheduled": true,
		"action": scheduledAction,
		"scheduled_time": scheduledTime.Format("2006-01-02 15:04:05"),
		"remaining_seconds": time.Until(scheduledTime).Seconds(),
	})
}

// RebootSystemImmediately 立即重启系统
func RebootSystemImmediately(c *gin.Context) {
	// 异步执行立即重启
	go func() {
		time.Sleep(1 * time.Second) // 给API响应时间
		exec.Command("reboot").Run()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "系统将立即重启",
		"action": "immediate_reboot",
	})
}

// PoweroffSystemImmediately 立即关闭系统
func PoweroffSystemImmediately(c *gin.Context) {
	// 异步执行立即关机
	go func() {
		time.Sleep(1 * time.Second) // 给API响应时间
		exec.Command("poweroff").Run()
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "系统将立即关闭",
		"action": "immediate_poweroff",
	})
}