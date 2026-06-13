package api

import (
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"nas-dashboard/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitStatusResponse 初始化状态响应
type InitStatusResponse struct {
	Initialized bool `json:"initialized"`
}

// InitializeRequest 初始化请求
type InitializeRequest struct {
	AdminUsername string `json:"adminUsername" binding:"required"`
	AdminPassword string `json:"adminPassword" binding:"required"`
	AdminEmail    string `json:"adminEmail"`
	SystemName    string `json:"systemName"`
}

// GetInitStatus 获取初始化状态
func GetInitStatus(c *gin.Context) {
	db := database.GetDB()
	var count int64
	// 检查是否有任何设置了密码的管理员用户
	db.Model(&models.User{}).Where("role = ? AND password_hash != ?", "admin", "").Count(&count)

	c.JSON(http.StatusOK, InitStatusResponse{
		Initialized: count > 0,
	})
}

// InitializeSystem 初始化系统
func InitializeSystem(c *gin.Context) {
	var req InitializeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	db := database.GetDB()
	var count int64
	db.Model(&models.User{}).Where("role = ? AND password_hash != ?", "admin", "").Count(&count)

	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "System is already initialized"})
		return
	}

	authService := service.NewAuthService()

	// 查找现有的默认 admin 用户（如果没有密码）
	var user models.User
	err := db.Where("username = ? AND password_hash = ?", "admin", "").First(&user).Error

	if err == nil {
		// 更新现有用户
		hashedPassword, err := authService.HashPassword(req.AdminPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user.Username = req.AdminUsername
		user.PasswordHash = hashedPassword
		user.Email = req.AdminEmail
		user.DisplayName = "Administrator"

		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update admin user"})
			return
		}
	} else {
		// 创建新管理员
		_, err := authService.CreateUser(req.AdminUsername, req.AdminPassword, req.AdminEmail, "Administrator", "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
			return
		}
	}

	// 设置系统名称
	if req.SystemName != "" {
		db.Model(&models.SystemConfig{}).Where("key = ?", "system.name").Update("value", req.SystemName)
	}

	c.JSON(http.StatusOK, gin.H{"message": "System initialized successfully"})
}
