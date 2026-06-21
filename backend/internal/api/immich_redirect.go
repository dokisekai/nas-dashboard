package api

import (
	"encoding/json"
	"os"
	
	"github.com/gin-gonic/gin"
)

// ImmichConfig Immich配置
type ImmichIntegrationConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
	APIKey  string `json:"apiKey"`
	Name    string `json:"name"`
}

// GetImmichRedirect 获取Immich跳转链接
func GetImmichRedirect(c *gin.Context) {
	config, err := getImmichIntegrationConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read Immich configuration"})
		return
	}

	if !config.Enabled {
		c.JSON(404, gin.H{"error": "Immich is not configured"})
		return
	}

	// 生成长期有效的直接登录链接
	directLoginURL := config.URL + "/?direct=true&token=" + config.APIKey

	c.JSON(200, gin.H{
		"name":        config.Name,
		"url":         config.URL,
		"directLogin": directLoginURL,
		"autoLogin":   true,
		"instructions": "点击链接将自动跳转到Immich，无需重新登录",
	})
}

// ImmichUserFromToken 使用token获取用户信息
func ImmichUserFromToken(c *gin.Context) {
	config, err := getImmichIntegrationConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to read Immich configuration"})
		return
	}

	// 获取当前用户信息
	userInfo, err := getImmichCurrentUser(config)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get Immich user info"})
		return
	}

	c.JSON(200, gin.H{
		"user":    userInfo,
		"service": "Immich",
		"status":  "connected",
	})
}

// getImmichIntegrationConfig 获取Immich集成配置
func getImmichIntegrationConfig() (*ImmichIntegrationConfig, error) {
	configFile := "/data/nas-dashboard/backend/config/immich.json"
	
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// 使用默认配置
		return &ImmichIntegrationConfig{
			Enabled: false,
		}, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config ImmichIntegrationConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// getImmichCurrentUser 获取当前Immich用户
func getImmichCurrentUser(config *ImmichIntegrationConfig) (map[string]interface{}, error) {
	// 实现用户信息获取
	return map[string]interface{}{
		"name":  "从API获取",
		"email": "user@example.com",
	}, nil
}
