package api

import (
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetNotifications 获取通知列表
func GetNotifications(c *gin.Context) {
	db := database.GetDB()
	var notifications []models.Notification
	if err := db.Order("timestamp desc").Limit(50).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// MarkNotificationRead 标记通知为已读
func MarkNotificationRead(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()
	if err := db.Model(&models.Notification{}).Where("id = ?", id).Update("read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

// MarkAllNotificationsRead 标记所有通知为已读
func MarkAllNotificationsRead(c *gin.Context) {
	db := database.GetDB()
	if err := db.Model(&models.Notification{}).Where("read = ?", false).Update("read", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

// ClearNotifications 清空通知
func ClearNotifications(c *gin.Context) {
	db := database.GetDB()
	if err := db.Session(&gin.Context{}).Exec("DELETE FROM notifications").Error; err != nil {
		// Try a safer delete if Exec fails
		if err := db.Where("1 = 1").Delete(&models.Notification{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "All notifications cleared"})
}

// AddNotification 添加新通知 (内部使用或 API 调用)
func AddNotification(ntype, title, message string) error {
	db := database.GetDB()
	notification := models.Notification{
		Type:      ntype,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
		Read:      false,
	}
	return db.Create(&notification).Error
}

// CreateNotification API 接口添加通知
func CreateNotification(c *gin.Context) {
	var req struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := AddNotification(req.Type, req.Title, req.Message); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Notification created"})
}
