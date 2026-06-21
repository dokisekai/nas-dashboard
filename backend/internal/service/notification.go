package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"nas-dashboard/internal/models"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// NotificationService 通知服务
type NotificationService struct {
	db          *gorm.DB
	wsUpgrader  websocket.Upgrader
	clients     map[string]*websocket.Conn
	clientsLock sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

// SystemEvent 系统事件
type SystemEvent struct {
	Type      string                 `json:"type"`
	Source    string                 `json:"source"`
	Severity  string                 `json:"severity"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// NotificationRule 通知规则
type NotificationRule struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	EventType   string  `json:"eventType"`   // disk_health, system_resource, backup_complete, security
	Conditions  map[string]interface{} `json:"conditions"`
	Actions     []NotificationAction `json:"actions"`
	Enabled     bool   `json:"enabled"`
	Cooldown    int    `json:"cooldown"`    // 冷却时间（秒）
	LastTriggered time.Time `json:"lastTriggered"`
}

// NotificationAction 通知动作
type NotificationAction struct {
	Type    string                 `json:"type"`    // websocket, email, webhook
	Config  map[string]interface{} `json:"config"`
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB) *NotificationService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &NotificationService{
		db:         db,
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 在生产环境中应该检查来源
			},
		},
		clients:    make(map[string]*websocket.Conn),
		ctx:        ctx,
		cancel:     cancel,
	}

	// 启动事件监听器
	go service.startEventListeners()

	return service
}

// StartEventListeners 启动事件监听器
func (s *NotificationService) startEventListeners() {
	log.Println("Starting notification service event listeners")

	// 启动各种事件监听器
	go s.monitorDiskEvents()
	go s.monitorSystemResourceEvents()
	go s.monitorBackupEvents()
	go s.monitorSecurityEvents()

	// 定期清理过期通知
	go s.cleanupExpiredNotifications()
}

// MonitorDiskEvents 监控磁盘健康事件
func (s *NotificationService) monitorDiskEvents() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkDiskHealth()
		}
	}
}

// CheckDiskHealth 检查磁盘健康状态
func (s *NotificationService) checkDiskHealth() {
	// 获取所有磁盘信息
	disks, err := s.getDiskInformation()
	if err != nil {
		log.Printf("Failed to get disk information: %v", err)
		return
	}

	for _, disk := range disks {
		// 检查 SMART 状态
		if disk.SmartStatus == "failing" {
			event := SystemEvent{
				Type:      "disk_health",
				Source:    "smart_monitor",
				Severity:  "critical",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"disk":       disk.Device,
					"status":     "failing",
					"message":    fmt.Sprintf("磁盘 %s 可能即将故障", disk.Device),
				},
			}
			s.processEvent(event)
		}

		// 检查温度
		if disk.Temperature > 60 {
			event := SystemEvent{
				Type:      "disk_temperature",
				Source:    "temperature_monitor",
				Severity:  "warning",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"disk":        disk.Device,
					"temperature": disk.Temperature,
					"message":     fmt.Sprintf("磁盘 %s 温度过高: %d°C", disk.Device, disk.Temperature),
				},
			}
			s.processEvent(event)
		}

		// 检查空间使用率
		usagePercent := (disk.Used / disk.Size) * 100
		if usagePercent > 90 {
			event := SystemEvent{
				Type:      "disk_space",
				Source:    "space_monitor",
				Severity:  "critical",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"disk":          disk.Device,
					"usage_percent": usagePercent,
					"available":     disk.Size - disk.Used,
					"message":       fmt.Sprintf("磁盘 %s 空间不足: %.1f%%", disk.Device, usagePercent),
				},
			}
			s.processEvent(event)
		} else if usagePercent > 80 {
			event := SystemEvent{
				Type:      "disk_space",
				Source:    "space_monitor",
				Severity:  "warning",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"disk":          disk.Device,
					"usage_percent": usagePercent,
					"available":     disk.Size - disk.Used,
					"message":       fmt.Sprintf("磁盘 %s 空间告警: %.1f%%", disk.Device, usagePercent),
				},
			}
			s.processEvent(event)
		}
	}
}

// MonitorSystemResourceEvents 监控系统资源事件
func (s *NotificationService) monitorSystemResourceEvents() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkSystemResources()
		}
	}
}

// CheckSystemResources 检查系统资源使用情况
func (s *NotificationService) checkSystemResources() {
	// 检查内存使用率
	memInfo, err := s.getMemoryInfo()
	if err == nil {
		usagePercent := (memInfo.Used / memInfo.Total) * 100
		if usagePercent > 95 {
			event := SystemEvent{
				Type:      "memory_alert",
				Source:    "resource_monitor",
				Severity:  "critical",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"usage_percent": usagePercent,
					"available":     memInfo.Total - memInfo.Used,
					"message":       fmt.Sprintf("内存严重不足: %.1f%%", usagePercent),
				},
			}
			s.processEvent(event)
		} else if usagePercent > 90 {
			event := SystemEvent{
				Type:      "memory_alert",
				Source:    "resource_monitor",
				Severity:  "warning",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"usage_percent": usagePercent,
					"available":     memInfo.Total - memInfo.Used,
					"message":       fmt.Sprintf("内存告警: %.1f%%", usagePercent),
				},
			}
			s.processEvent(event)
		}
	}

	// 检查 CPU 使用率
	cpuInfo, err := s.getCPUInfo()
	if err == nil {
		if cpuInfo.Usage > 95 {
			event := SystemEvent{
				Type:      "cpu_alert",
				Source:    "resource_monitor",
				Severity:  "warning",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"usage_percent": cpuInfo.Usage,
					"message":       fmt.Sprintf("CPU 使用率过高: %.1f%%", cpuInfo.Usage),
				},
			}
			s.processEvent(event)
		}
	}

	// 检查负载
	loadInfo, err := s.getLoadAverage()
	if err == nil {
		cpuCount := float64(s.getCPUCount())
		if loadInfo.Load1 > cpuCount*2 {
			event := SystemEvent{
				Type:      "load_alert",
				Source:    "resource_monitor",
				Severity:  "warning",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"load1":     loadInfo.Load1,
					"load5":     loadInfo.Load5,
					"load15":    loadInfo.Load15,
					"cpu_count": cpuCount,
					"message":   fmt.Sprintf("系统负载过高: %.2f", loadInfo.Load1),
				},
			}
			s.processEvent(event)
		}
	}
}

// MonitorBackupEvents 监控备份事件
func (s *NotificationService) monitorBackupEvents() {
	// 这里可以监听备份完成事件
	// 实际实现需要与备份系统集成
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkBackupStatus()
		}
	}
}

// CheckBackupStatus 检查备份状态
func (s *NotificationService) checkBackupStatus() {
	// 检查最近的备份任务
	// 这里需要与备份系统集成
	// 示例：检查是否有失败的备份
}

// MonitorSecurityEvents 监控安全事件
func (s *NotificationService) monitorSecurityEvents() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.checkSecurityEvents()
		}
	}
}

// CheckSecurityEvents 检查安全事件
func (s *NotificationService) checkSecurityEvents() {
	// 检查失败登录次数
	failedLogins := s.getRecentFailedLogins()
	if len(failedLogins) > 5 {
		event := SystemEvent{
			Type:      "security_alert",
			Source:    "auth_monitor",
			Severity:  "warning",
			Timestamp: time.Now(),
			Data: map[string]interface{}{
				"failed_count": len(failedLogins),
				"ips":          s.getUniqueIPs(failedLogins),
				"message":       fmt.Sprintf("检测到 %d 次失败登录尝试", len(failedLogins)),
			},
		}
		s.processEvent(event)
	}
}

// ProcessEvent 处理系统事件
func (s *NotificationService) processEvent(event SystemEvent) {
	// 获取匹配的通知规则
	rules := s.getMatchingRules(event)

	// 应用通知规则
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		// 检查冷却时间
		if !rule.LastTriggered.IsZero() && time.Since(rule.LastTriggered) < time.Duration(rule.Cooldown)*time.Second {
			continue
		}

		// 执行通知动作
		s.executeActions(rule, event)

		// 更新最后触发时间
		rule.LastTriggered = time.Now()
		s.updateRule(rule)
	}

	// 如果没有匹配规则，使用默认行为
	if len(rules) == 0 && event.Severity == "critical" {
		s.createDefaultNotification(event)
	}
}

// GetMatchingRules 获取匹配的通知规则
func (s *NotificationService) getMatchingRules(event SystemEvent) []NotificationRule {
	var rules []models.NotificationRule

	err := s.db.Where("event_type = ? AND enabled = ?", event.Type, true).Find(&rules).Error
	if err != nil {
		log.Printf("Failed to get notification rules: %v", err)
		return []NotificationRule{}
	}

	var matchedRules []NotificationRule
	for _, rule := range rules {
		var conditions map[string]interface{}
		json.Unmarshal([]byte(rule.Conditions), &conditions)
		if s.matchConditions(conditions, event.Data) {
			matchedRules = append(matchedRules, NotificationRule{
				ID:          rule.ID,
				Name:        rule.Name,
				EventType:   rule.EventType,
				Conditions:  conditions,
				Actions:     s.parseActions(rule.Actions),
				Enabled:     rule.Enabled,
				Cooldown:    rule.Cooldown,
				LastTriggered: rule.LastTriggered,
			})
		}
	}

	return matchedRules
}

// MatchConditions 匹配条件
func (s *NotificationService) matchConditions(conditions map[string]interface{}, eventData map[string]interface{}) bool {
	for key, expectedValue := range conditions {
		actualValue, exists := eventData[key]
		if !exists {
			return false
		}

		// 简单的类型比较
		switch expected := expectedValue.(type) {
		case int:
			if actual, ok := actualValue.(float64); ok {
				if int(actual) != expected {
					return false
				}
			} else {
				return false
			}
		case float64:
			if actual, ok := actualValue.(float64); ok {
				if math.Abs(actual-expected) > 0.01 {
					return false
				}
			} else {
				return false
			}
		case string:
			if actual, ok := actualValue.(string); ok {
				if actual != expected {
					return false
				}
			} else {
				return false
			}
		case bool:
			if actual, ok := actualValue.(bool); ok {
				if actual != expected {
					return false
				}
			} else {
				return false
			}
		}
	}

	return true
}

// ExecuteActions 执行通知动作
func (s *NotificationService) executeActions(rule NotificationRule, event SystemEvent) {
	for _, action := range rule.Actions {
		switch action.Type {
		case "websocket":
			s.sendWebSocketNotification(event, action.Config)
		case "email":
			s.sendEmailNotification(event, action.Config)
		case "webhook":
			s.sendWebhookNotification(event, action.Config)
		}
	}
}

// SendWebSocketNotification 发送 WebSocket 通知
func (s *NotificationService) sendWebSocketNotification(event SystemEvent, config map[string]interface{}) {
	notification := models.Notification{
		Type:      s.getNotificationType(event.Severity),
		Title:     s.getEventTitle(event),
		Message:   event.Data["message"].(string),
		Timestamp: event.Timestamp,
		Read:      false,
	}

	// 保存到数据库
	if err := s.db.Create(&notification).Error; err != nil {
		log.Printf("Failed to save notification: %v", err)
		return
	}

	// 通过 WebSocket 推送给所有连接的客户端
	s.broadcastNotification(notification)
}

// BroadcastNotification 广播通知给所有连接的客户端
func (s *NotificationService) broadcastNotification(notification models.Notification) {
	s.clientsLock.RLock()
	defer s.clientsLock.RUnlock()

	message := map[string]interface{}{
		"type":         "notification",
		"notification": notification,
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return
	}

	for clientID, client := range s.clients {
		if err := client.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
			log.Printf("Failed to send notification to client %s: %v", clientID, err)
			// 移除失败的连接
			delete(s.clients, clientID)
		}
	}
}

// SendEmailNotification 发送邮件通知
func (s *NotificationService) sendEmailNotification(event SystemEvent, config map[string]interface{}) {
	// 实现邮件通知逻辑
	// 这里需要集成邮件服务
	log.Printf("Email notification for event %s: %s", event.Type, event.Data["message"])
}

// SendWebhookNotification 发送 Webhook 通知
func (s *NotificationService) sendWebhookNotification(event SystemEvent, config map[string]interface{}) {
	// 实现 Webhook 通知逻辑
	log.Printf("Webhook notification for event %s: %s", event.Type, event.Data["message"])
}

// CreateDefaultNotification 创建默认通知
func (s *NotificationService) createDefaultNotification(event SystemEvent) {
	notification := models.Notification{
		Type:      s.getNotificationType(event.Severity),
		Title:     s.getEventTitle(event),
		Message:   event.Data["message"].(string),
		Timestamp: event.Timestamp,
		Read:      false,
		Persistent: true,
	}

	if err := s.db.Create(&notification).Error; err != nil {
		log.Printf("Failed to save default notification: %v", err)
	}

	s.broadcastNotification(notification)
}

// GetNotificationType 获取通知类型
func (s *NotificationService) getNotificationType(severity string) string {
	switch severity {
	case "critical":
		return "error"
	case "warning":
		return "warning"
	case "info":
		return "info"
	default:
		return "success"
	}
}

// GetEventTitle 获取事件标题
func (s *NotificationService) getEventTitle(event SystemEvent) string {
	switch event.Type {
	case "disk_health":
		return "磁盘健康告警"
	case "disk_temperature":
		return "磁盘温度告警"
	case "disk_space":
		return "磁盘空间告警"
	case "memory_alert":
		return "内存告警"
	case "cpu_alert":
		return "CPU 告警"
	case "load_alert":
		return "系统负载告警"
	case "security_alert":
		return "安全告警"
	case "backup_complete":
		return "备份完成"
	case "backup_failed":
		return "备份失败"
	default:
		return "系统通知"
	}
}

// CleanupExpiredNotifications 清理过期通知
func (s *NotificationService) cleanupExpiredNotifications() {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			// 删除30天前已读的通知
			expiryDate := time.Now().AddDate(0, 0, -30)
			s.db.Where("read = ? AND timestamp < ?", true, expiryDate).Delete(&models.Notification{})
		}
	}
}

// RegisterWebSocketClient 注册 WebSocket 客户端
func (s *NotificationService) RegisterWebSocketClient(clientID string, conn *websocket.Conn) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	s.clients[clientID] = conn
	log.Printf("WebSocket client registered: %s", clientID)
}

// UnregisterWebSocketClient 注销 WebSocket 客户端
func (s *NotificationService) UnregisterWebSocketClient(clientID string) {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	delete(s.clients, clientID)
	log.Printf("WebSocket client unregistered: %s", clientID)
}

// Stop 停止通知服务
func (s *NotificationService) Stop() {
	s.cancel()

	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	// 关闭所有 WebSocket 连接
	for clientID, client := range s.clients {
		client.Close()
		delete(s.clients, clientID)
	}

	log.Println("Notification service stopped")
}

// 辅助方法（需要实现）
func (s *NotificationService) getDiskInformation() ([]models.Disk, error) {
	// 实现获取磁盘信息的逻辑
	return []models.Disk{}, nil
}

func (s *NotificationService) getMemoryInfo() (*models.MemoryInfo, error) {
	// 实现获取内存信息的逻辑
	return &models.MemoryInfo{}, nil
}

func (s *NotificationService) getCPUInfo() (*models.CPUInfo, error) {
	// 实现获取 CPU 信息的逻辑
	return &models.CPUInfo{}, nil
}

func (s *NotificationService) getLoadAverage() (*models.LoadInfo, error) {
	// 实现获取负载信息的逻辑
	return &models.LoadInfo{}, nil
}

func (s *NotificationService) getCPUCount() int {
	// 实现 CPU 核心数获取逻辑
	return 1
}

func (s *NotificationService) getRecentFailedLogins() []models.LoginAttempt {
	// 实现获取最近失败登录尝试的逻辑
	return []models.LoginAttempt{}
}

func (s *NotificationService) getUniqueIPs(attempts []models.LoginAttempt) []string {
	// 实现获取唯一 IP 的逻辑
	return []string{}
}

func (s *NotificationService) parseActions(actionsJSON string) []NotificationAction {
	// 实现 JSON 解析逻辑
	return []NotificationAction{}
}

func (s *NotificationService) updateRule(rule NotificationRule) error {
	// 实现更新规则的逻辑
	return nil
}