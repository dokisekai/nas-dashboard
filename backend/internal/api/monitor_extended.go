package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"nas-dashboard/internal/models"
	"nas-dashboard/pkg/system"
	"gorm.io/gorm"
)

// MonitorAPI 监控API处理器
type MonitorAPI struct {
	DB *gorm.DB
}

// NewMonitorAPI 创建监控API
func NewMonitorAPI(db *gorm.DB) *MonitorAPI {
	return &MonitorAPI{DB: db}
}

// GetProcesses 获取进程列表
func (api *MonitorAPI) GetProcesses(c *gin.Context) {
	processes, err := system.GetProcesses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get processes: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"processes": processes})
}

// GetProcess 获取进程详情
func (api *MonitorAPI) GetProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PID"})
		return
	}

	process, err := system.GetProcess(int32(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get process: %v", err)})
		return
	}

	c.JSON(http.StatusOK, process)
}

// KillProcess 终止进程
func (api *MonitorAPI) KillProcess(c *gin.Context) {
	pidStr := c.Param("pid")
	_, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PID"})
		return
	}

	signal := c.Query("signal")
	if signal == "" {
		signal = "15" // SIGTERM
	}

	// 发送信号到进程
	cmd := exec.Command("kill", "-"+signal, pidStr)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to kill process: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Process killed successfully"})
}

// GetServices 获取系统服务列表
func (api *MonitorAPI) GetServices(c *gin.Context) {
	services, err := system.GetServices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get services: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"services": services})
}

// GetService 获取服务详情
func (api *MonitorAPI) GetService(c *gin.Context) {
	name := c.Param("name")

	service, err := system.GetService(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get service: %v", err)})
		return
	}

	c.JSON(http.StatusOK, service)
}

// StartService 启动服务
func (api *MonitorAPI) StartService(c *gin.Context) {
	name := c.Param("name")

	cmd := exec.Command("systemctl", "start", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to start service: %v, output: %s", err, string(output))})
		return
	}

	// 记录事件
	api.recordSystemEvent("service", "info", fmt.Sprintf("Service %s started", name), "")

	c.JSON(http.StatusOK, gin.H{"message": "Service started successfully"})
}

// StopService 停止服务
func (api *MonitorAPI) StopService(c *gin.Context) {
	name := c.Param("name")

	cmd := exec.Command("systemctl", "stop", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to stop service: %v, output: %s", err, string(output))})
		return
	}

	// 记录事件
	api.recordSystemEvent("service", "info", fmt.Sprintf("Service %s stopped", name), "")

	c.JSON(http.StatusOK, gin.H{"message": "Service stopped successfully"})
}

// RestartService 重启服务
func (api *MonitorAPI) RestartService(c *gin.Context) {
	name := c.Param("name")

	cmd := exec.Command("systemctl", "restart", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to restart service: %v, output: %s", err, string(output))})
		return
	}

	// 记录事件
	api.recordSystemEvent("service", "info", fmt.Sprintf("Service %s restarted", name), "")

	c.JSON(http.StatusOK, gin.H{"message": "Service restarted successfully"})
}

// GetTemperature 获取系统温度
func (api *MonitorAPI) GetTemperature(c *gin.Context) {
	temp, err := system.GetTemperature()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to get temperature: %v", err)})
		return
	}

	c.JSON(http.StatusOK, temp)
}

// GetEvents 获取系统事件
func (api *MonitorAPI) GetEvents(c *gin.Context) {
	// 解析查询参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "50")
	eventType := c.Query("type")
	source := c.Query("source")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	query := api.DB.Model(&models.SystemEvent{})

	// 应用过滤
	if eventType != "" {
		query = query.Where("type = ?", eventType)
	}
	if source != "" {
		query = query.Where("source = ?", source)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var events []models.SystemEvent
	offset := (pageInt - 1) * limitInt
	if result := query.Order("created_at DESC").Limit(limitInt).Offset(offset).Find(&events); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"total":  total,
		"page":   pageInt,
		"limit":  limitInt,
	})
}

// GetLogs 获取系统日志
func (api *MonitorAPI) GetLogs(c *gin.Context) {
	// 解析查询参数
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "100")
	level := c.Query("level")
	component := c.Query("component")

	pageInt, _ := strconv.Atoi(page)
	limitInt, _ := strconv.Atoi(limit)

	query := api.DB.Model(&models.SystemLog{})

	// 应用过滤
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if component != "" {
		query = query.Where("component = ?", component)
	}

	// 计算总数
	var total int64
	query.Count(&total)

	// 分页查询
	var logs []models.SystemLog
	offset := (pageInt - 1) * limitInt
	if result := query.Order("timestamp DESC").Limit(limitInt).Offset(offset).Find(&logs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
		"page":  pageInt,
		"limit": limitInt,
	})
}

// ClearLogs 清除日志
func (api *MonitorAPI) ClearLogs(c *gin.Context) {
	level := c.Query("level")
	component := c.Query("component")

	query := api.DB.Model(&models.SystemLog{})

	// 应用过滤
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if component != "" {
		query = query.Where("component = ?", component)
	}

	if result := query.Delete(&models.SystemLog{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logs cleared successfully"})
}

// GetAlerts 获取告警规则
func (api *MonitorAPI) GetAlerts(c *gin.Context) {
	var alerts []models.AlertRule
	if result := api.DB.Find(&alerts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

// CreateAlert 创建告警规则
func (api *MonitorAPI) CreateAlert(c *gin.Context) {
	var alert models.AlertRule
	if err := c.ShouldBindJSON(&alert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证告警规则
	if alert.Name == "" || alert.Type == "" || alert.Condition == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, type, and condition are required"})
		return
	}

	if result := api.DB.Create(&alert); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// UpdateAlert 更新告警规则
func (api *MonitorAPI) UpdateAlert(c *gin.Context) {
	id := c.Param("id")

	var alert models.AlertRule
	if result := api.DB.First(&alert, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	var updateData models.AlertRule
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	if updateData.Name != "" {
		alert.Name = updateData.Name
	}
	if updateData.Type != "" {
		alert.Type = updateData.Type
	}
	if updateData.Condition != "" {
		alert.Condition = updateData.Condition
	}
	if updateData.Threshold > 0 {
		alert.Threshold = updateData.Threshold
	}
	if updateData.Duration > 0 {
		alert.Duration = updateData.Duration
	}
	if updateData.Severity != "" {
		alert.Severity = updateData.Severity
	}
	alert.Enabled = updateData.Enabled
	if updateData.Actions != "" {
		alert.Actions = updateData.Actions
	}

	if result := api.DB.Save(&alert); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// DeleteAlert 删除告警规则
func (api *MonitorAPI) DeleteAlert(c *gin.Context) {
	id := c.Param("id")

	if result := api.DB.Delete(&models.AlertRule{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert deleted successfully"})
}

// recordSystemEvent 记录系统事件
func (api *MonitorAPI) recordSystemEvent(source, eventType, title, message string) {
	event := models.SystemEvent{
		Type:      eventType,
		Source:    source,
		Title:     title,
		Message:   message,
		Resolved:  false,
	}

	api.DB.Create(&event)
}

// checkAlerts 检查告警规则（应该在定时任务中调用）
func (api *MonitorAPI) checkAlerts() {
	// 如果数据库不可用，直接返回
	if api.DB == nil {
		return
	}

	var alerts []models.AlertRule
	if result := api.DB.Where("enabled = ?", true).Find(&alerts); result.Error != nil {
		return
	}

	for _, alert := range alerts {
		// 根据类型检查相应的指标
		triggered := false
		var currentValue float64

		switch alert.Type {
		case "cpu":
			cpuInfo, err := system.GetCPUInfo()
			if err == nil {
				// 使用总CPU使用率（转换为百分比）
				currentValue = cpuInfo.Usage * 100
				triggered = compareValues(currentValue, alert.Threshold, alert.Condition)
			}

		case "memory":
			memInfo, err := system.GetMemoryInfo()
			if err == nil {
				currentValue = memInfo.Percent
				triggered = compareValues(currentValue, alert.Threshold, alert.Condition)
			}

		case "disk":
			diskInfo, err := system.GetDiskInfo()
			if err == nil {
				for _, disk := range diskInfo.Disks {
					currentValue = disk.UsedPercent
					if compareValues(currentValue, alert.Threshold, alert.Condition) {
						triggered = true
						break
					}
				}
			}

		case "temperature":
			temp, err := system.GetTemperature()
			if err == nil {
				for _, sensor := range temp.Sensors {
					currentValue = sensor.Current
					if compareValues(currentValue, alert.Threshold, alert.Condition) {
						triggered = true
						break
					}
				}
			}
		}

		if triggered {
			now := time.Now()
			alert.LastTriggered = &now
			api.DB.Save(&alert)

			// 创建系统事件
			event := models.SystemEvent{
				Type:      alert.Severity,
				Source:    "alert",
				Title:     fmt.Sprintf("Alert triggered: %s", alert.Name),
				Message:   fmt.Sprintf("Current value: %.2f, Threshold: %.2f", currentValue, alert.Threshold),
				Resolved:  false,
			}
			api.DB.Create(&event)
		}
	}
}

// StartMonitoring 启动监控定时任务
func (api *MonitorAPI) StartMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// 检查告警规则
		api.checkAlerts()

		// 清理旧的历史数据
		api.cleanupOldHistory()
	}
}

// cleanupOldHistory 清理旧的历史数据
func (api *MonitorAPI) cleanupOldHistory() {
	// 如果数据库不可用，直接返回
	if api.DB == nil {
		return
	}

	// 保留最近7天的数据
	cutoffDate := time.Now().AddDate(0, 0, -7)

	api.DB.Where("created_at < ?", cutoffDate).Delete(&models.MonitorHistory{})
}

// compareValues 比较值
func compareValues(current, threshold float64, condition string) bool {
	switch condition {
	case ">":
		return current > threshold
	case "<":
		return current < threshold
	case ">=":
		return current >= threshold
	case "<=":
		return current <= threshold
	case "==", "=":
		return current == threshold
	case "!=":
		return current != threshold
	default:
		return false
	}
}
