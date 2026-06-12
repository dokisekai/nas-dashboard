package api

import (
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfigService 配置服务
type ConfigService struct {
	db *gorm.DB
}

// NewConfigService 创建配置服务
func NewConfigService() *ConfigService {
	return &ConfigService{
		db: database.GetDB(),
	}
}

var configService = NewConfigService()

// ConfigValue 配置值结构
type ConfigValue struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

// SetConfigRequest 设置配置请求
type SetConfigRequest struct {
	Key         string `json:"key" binding:"required"`
	Value       string `json:"value" binding:"required"`
	Type        string `json:"type"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

// BulkSetConfigRequest 批量设置配置请求
type BulkSetConfigRequest struct {
	Configs []SetConfigRequest `json:"configs" binding:"required"`
}

// GetConfigs 获取系统配置
func GetConfigs(c *gin.Context) {
	category := c.Query("category")
	includePrivate := c.DefaultQuery("private", "false") == "true"

	configs, err := configService.GetConfigs(category, includePrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
		"total":   len(configs),
	})
}

// GetPublicConfigs 获取公开配置（不需要认证）
func GetPublicConfigs(c *gin.Context) {
	configs, err := configService.GetPublicConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"configs": configs,
		"total":   len(configs),
	})
}

// GetConfig 获取单个配置
func GetConfig(c *gin.Context) {
	key := c.Param("key")

	config, err := configService.GetConfig(key)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Config not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 检查权限
	if !config.IsPublic {
		if _, exists := c.Get("userId"); !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, config)
}

// SetConfig 设置配置
func SetConfig(c *gin.Context) {
	var req SetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 验证配置值
	if err := configService.ValidateConfig(req.Type, req.Value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户信息
	userID, _ := middleware.GetUserIDAsUint(c)

	// 设置配置
	config, err := configService.SetConfig(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config set successfully",
		"config":  config,
	})
}

// BulkSetConfig 批量设置配置
func BulkSetConfig(c *gin.Context) {
	var req BulkSetConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 获取用户信息
	userID, _ := middleware.GetUserIDAsUint(c)

	// 批量设置配置
	configs, err := configService.BulkSetConfig(req.Configs, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Configs set successfully",
		"configs": configs,
	})
}

// DeleteConfig 删除配置
func DeleteConfig(c *gin.Context) {
	key := c.Param("key")

	// 删除配置
	if err := configService.DeleteConfig(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config deleted successfully",
		"key":     key,
	})
}

// ResetConfig 重置配置为默认值
func ResetConfig(c *gin.Context) {
	key := c.Param("key")

	// 重置配置
	config, err := configService.ResetConfig(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Config reset successfully",
		"config":  config,
	})
}

// GetConfigs 获取配置列表
func (s *ConfigService) GetConfigs(category string, includePrivate bool) ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	query := s.db.Model(&models.SystemConfig{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if !includePrivate {
		query = query.Where("is_public = ?", true)
	}

	if err := query.Order("category, key").Find(&configs).Error; err != nil {
		return nil, err
	}

	return configs, nil
}

// GetPublicConfigs 获取公开配置
func (s *ConfigService) GetPublicConfigs() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	if err := s.db.Where("is_public = ?", true).Order("category, key").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetConfig 获取单个配置
func (s *ConfigService) GetConfig(key string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	if err := s.db.Where("key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

// SetConfig 设置配置
func (s *ConfigService) SetConfig(req SetConfigRequest, userID uint) (*models.SystemConfig, error) {
	var config models.SystemConfig

	// 查找现有配置
	if err := s.db.Where("key = ?", req.Key).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新配置
			config = models.SystemConfig{
				Key:         req.Key,
				Value:       req.Value,
				Type:        req.Type,
				Category:    req.Category,
				Description: req.Description,
				IsPublic:    req.IsPublic,
			}

			// 设置默认值
			if config.Type == "" {
				config.Type = "string"
			}
			if config.Category == "" {
				config.Category = "general"
			}

			if err := s.db.Create(&config).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		// 更新现有配置
		config.Value = req.Value
		if req.Type != "" {
			config.Type = req.Type
		}
		if req.Category != "" {
			config.Category = req.Category
		}
		if req.Description != "" {
			config.Description = req.Description
		}
		config.IsPublic = req.IsPublic

		if err := s.db.Save(&config).Error; err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// BulkSetConfig 批量设置配置
func (s *ConfigService) BulkSetConfig(reqs []SetConfigRequest, userID uint) ([]models.SystemConfig, error) {
	var configs []models.SystemConfig

	// 使用事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, req := range reqs {
		config, err := s.SetConfig(req, userID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		configs = append(configs, *config)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return configs, nil
}

// DeleteConfig 删除配置
func (s *ConfigService) DeleteConfig(key string) error {
	return s.db.Where("key = ?", key).Delete(&models.SystemConfig{}).Error
}

// ResetConfig 重置配置
func (s *ConfigService) ResetConfig(key string) (*models.SystemConfig, error) {
	// 获取默认值
	defaultValue, err := s.GetDefaultValue(key)
	if err != nil {
		return nil, err
	}

	// 更新配置
	var config models.SystemConfig
	if err := s.db.Where("key = ?", key).First(&config).Error; err != nil {
		return nil, err
	}

	config.Value = defaultValue.Value
	if err := s.db.Save(&config).Error; err != nil {
		return nil, err
	}

	return &config, nil
}

// ValidateConfig 验证配置值
func (s *ConfigService) ValidateConfig(configType, value string) error {
	switch configType {
	case "int":
		if _, err := strconv.Atoi(value); err != nil {
			return &ConfigValidationError{Type: configType, Value: value, Reason: "invalid integer"}
		}
	case "bool":
		if value != "true" && value != "false" {
			return &ConfigValidationError{Type: configType, Value: value, Reason: "invalid boolean"}
		}
	case "json":
		// 简单的JSON验证
		if value == "" || value == "null" {
			return nil
		}
		if value[0] != '{' && value[0] != '[' {
			return &ConfigValidationError{Type: configType, Value: value, Reason: "invalid JSON"}
		}
	}
	return nil
}

// GetDefaultValue 获取默认配置值
func (s *ConfigService) GetDefaultValue(key string) (*models.SystemConfig, error) {
	defaults := []models.SystemConfig{
		{
			Key:   "system.name",
			Value: "NAS Dashboard",
			Type:  "string",
		},
		{
			Key:   "system.timezone",
			Value: "Asia/Shanghai",
			Type:  "string",
		},
		{
			Key:   "security.session_timeout",
			Value: "86400",
			Type:  "int",
		},
		{
			Key:   "security.max_login_attempts",
			Value: "5",
			Type:  "int",
		},
		{
			Key:   "backup.auto_backup_enabled",
			Value: "true",
			Type:  "bool",
		},
		{
			Key:   "backup.retention_days",
			Value: "30",
			Type:  "int",
		},
	}

	for _, def := range defaults {
		if def.Key == key {
			return &def, nil
		}
	}

	return nil, &ConfigNotFoundError{Key: key}
}

// ConfigValidationError 配置验证错误
type ConfigValidationError struct {
	Type   string
	Value  string
	Reason string
}

func (e *ConfigValidationError) Error() string {
	return "invalid " + e.Type + " value: " + e.Value + " (" + e.Reason + ")"
}

// ConfigNotFoundError 配置未找到错误
type ConfigNotFoundError struct {
	Key string
}

func (e *ConfigNotFoundError) Error() string {
	return "config not found: " + e.Key
}
