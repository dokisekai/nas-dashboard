package api

import (
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"nas-dashboard/internal/service"
	"time"

	"github.com/gin-gonic/gin"
)

// NotificationService 全局通知服务实例
var notificationService *service.NotificationService

// InitNotificationService 初始化通知服务
func InitNotificationService(db *database.DBConfig) {
	notificationService = service.NewNotificationService(db.DB)
}

// GetNotifications 获取通知列表
func GetNotifications(c *gin.Context) {
	var notifications []models.Notification

	// 默认只获取未读的通知，或者根据参数获取所有
	readOnly := c.DefaultQuery("read", "all")
	query := database.DB.Order("timestamp DESC")

	if readOnly == "unread" {
		query = query.Where("read = ?", false)
	} else if readOnly == "read" {
		query = query.Where("read = ?", true)
	}

	if err := query.Find(&notifications).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取通知失败"})
		return
	}

	c.JSON(200, notifications)
}

// CreateNotification 创建通知
func CreateNotification(c *gin.Context) {
	var notification models.Notification
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	notification.Timestamp = time.Now()
	notification.Read = false

	if err := database.DB.Create(&notification).Error; err != nil {
		c.JSON(500, gin.H{"error": "创建通知失败"})
		return
	}

	// 通过 WebSocket 推送新通知
	if notificationService != nil {
		notificationService.BroadcastNotification(notification)
	}

	c.JSON(201, notification)
}

// MarkAsRead 标记通知为已读
func MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	var notification models.Notification
	if err := database.DB.First(&notification, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "通知不存在"})
		return
	}

	notification.Read = true
	if err := database.DB.Save(&notification).Error; err != nil {
		c.JSON(500, gin.H{"error": "标记失败"})
		return
	}

	c.JSON(200, notification)
}

// MarkAllAsRead 标记所有通知为已读
func MarkAllAsRead(c *gin.Context) {
	if err := database.DB.Model(&models.Notification{}).Where("read = ?", false).Update("read", true).Error; err != nil {
		c.JSON(500, gin.H{"error": "批量标记失败"})
		return
	}

	c.JSON(200, gin.H{"message": "所有通知已标记为已读"})
}

// DeleteNotification 删除通知
func DeleteNotification(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Notification{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "删除通知失败"})
		return
	}

	c.JSON(200, gin.H{"message": "通知已删除"})
}

// ClearAllNotifications 清除所有通知
func ClearAllNotifications(c *gin.Context) {
	if err := database.DB.Where("read = ?", true).Delete(&models.Notification{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "清除通知失败"})
		return
	}

	c.JSON(200, gin.H{"message": "已读通知已清除"})
}

// GetNotificationRules 获取通知规则
func GetNotificationRules(c *gin.Context) {
	var rules []models.NotificationRule
	if err := database.DB.Order("created_at DESC").Find(&rules).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取通知规则失败"})
		return
	}

	c.JSON(200, rules)
}

// CreateNotificationRule 创建通知规则
func CreateNotificationRule(c *gin.Context) {
	var rule models.NotificationRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&rule).Error; err != nil {
		c.JSON(500, gin.H{"error": "创建通知规则失败"})
		return
	}

	c.JSON(201, rule)
}

// UpdateNotificationRule 更新通知规则
func UpdateNotificationRule(c *gin.Context) {
	id := c.Param("id")

	var rule models.NotificationRule
	if err := database.DB.First(&rule, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "通知规则不存在"})
		return
	}

	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&rule).Error; err != nil {
		c.JSON(500, gin.H{"error": "更新通知规则失败"})
		return
	}

	c.JSON(200, rule)
}

// DeleteNotificationRule 删除通知规则
func DeleteNotificationRule(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.NotificationRule{}, id).Error; err != nil {
		c.JSON(500, gin.H{"error": "删除通知规则失败"})
		return
	}

	c.JSON(200, gin.H{"message": "通知规则已删除"})
}

// TestNotificationRule 测试通知规则
func TestNotificationRule(c *gin.Context) {
	id := c.Param("id")

	var rule models.NotificationRule
	if err := database.DB.First(&rule, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "通知规则不存在"})
		return
	}

	// 创建测试通知
	testNotification := models.Notification{
		Type:       "info",
		Title:      "测试通知",
		Message:    "这是来自规则 \"" + rule.Name + "\" 的测试通知",
		Timestamp:  time.Now(),
		Read:       false,
		Persistent: false,
	}

	if err := database.DB.Create(&testNotification).Error; err != nil {
		c.JSON(500, gin.H{"error": "创建测试通知失败"})
		return
	}

	// 通过 WebSocket 推送测试通知
	if notificationService != nil {
		notificationService.BroadcastNotification(testNotification)
	}

	c.JSON(200, gin.H{"message": "测试通知已发送", "notification": testNotification})
}

// GetUnreadCount 获取未读通知数量
func GetUnreadCount(c *gin.Context) {
	var count int64
	if err := database.DB.Model(&models.Notification{}).Where("read = ?", false).Count(&count).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取未读数量失败"})
		return
	}

	c.JSON(200, gin.H{"count": count})
}

// WebSocketNotificationHandler WebSocket通知处理器
func WebSocketNotificationHandler(c *gin.Context) {
	// 升级 HTTP 连接到 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// 生成客户端ID
	clientID := generateClientID()

	// 注册客户端
	if notificationService != nil {
		notificationService.RegisterWebSocketClient(clientID, conn)
		defer notificationService.UnregisterWebSocketClient(clientID)
	}

	// 处理 WebSocket 消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// 这里可以处理客户端发送的消息
		// 比如心跳包、确认消息等

		if messageType == websocket.TextMessage {
			// 处理文本消息
			log.Printf("WebSocket message received: %s", string(p))
		}
	}
}