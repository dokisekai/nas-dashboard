package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// ImmichConfig Immich配置
type ImmichConfig struct {
	APIURL    string
	APIKey    string
	Enabled   bool
}

// ImmichUser Immich用户
type ImmichUser struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	IsActive bool   `json:"isActive"`
	Role     string `json:"role,omitempty"`
}

// ImmichUserResponse Immich用户响应
type ImmichUserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
	Role     string `json:"role"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ImmichCreateUserRequest 创建Immich用户请求
type ImmichCreateUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ImmichUpdateUserRequest 更新Immich用户请求
type ImmichUpdateUserRequest struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	IsActive *bool  `json:"isActive,omitempty"`
	Role     string `json:"role,omitempty"`
}

// ImmichUsersResponse Immich用户列表响应
type ImmichUsersResponse struct {
	Users []ImmichUserResponse `json:"users"`
	Total int                  `json:"total"`
}

// getImmichConfig 获取Immich配置
func getImmichConfig() (*ImmichConfig, error) {
	// 从环境变量读取配置
	apiURL := os.Getenv("IMMICH_API_URL")
	apiKey := os.Getenv("IMMICH_API_KEY")

	// 如果环境变量不存在，使用默认配置
	if apiURL == "" {
		apiURL = "http://localhost:2283/api"
	}

	// 如果API密钥为空，尝试从配置文件读取
	if apiKey == "" {
		data, err := os.ReadFile(".env.immich")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					switch key {
					case "IMMICH_API_URL":
						apiURL = value
					case "IMMICH_API_KEY":
						apiKey = value
					}
				}
			}
		}
	}

	return &ImmichConfig{
		APIURL:  apiURL,
		APIKey:  apiKey,
		Enabled: apiURL != "" && apiKey != "" && apiKey != "your_immich_api_key_here",
	}, nil
}

// makeImmichRequest 发送Immich API请求
func makeImmichRequest(method, endpoint string, body interface{}) ([]byte, error) {
	config, err := getImmichConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get Immich config: %w", err)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("Immich is not configured")
	}

	url := config.APIURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.APIKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Immich API error: status %d, response: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// GetImmichUsers 获取Immich用户列表
func GetImmichUsers(c *gin.Context) {
	data, err := makeImmichRequest("GET", "/users", nil)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get Immich users: %v", err)})
		return
	}

	var users []ImmichUserResponse
	if err := json.Unmarshal(data, &users); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to parse response: %v", err)})
		return
	}

	c.JSON(200, gin.H{
		"users": users,
		"total": len(users),
	})
}

// GetImmichUser 获取单个Immich用户
func GetImmichUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	data, err := makeImmichRequest("GET", "/users/"+userID, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get Immich user: %v", err)})
		return
	}

	var user ImmichUserResponse
	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to parse response: %v", err)})
		return
	}

	c.JSON(200, user)
}

// CreateImmichUser 创建Immich用户
func CreateImmichUser(c *gin.Context) {
	var req ImmichCreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	userData := ImmichUser{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		IsActive: true,
	}

	data, err := makeImmichRequest("POST", "/users", userData)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create Immich user: %v", err)})
		return
	}

	var user ImmichUserResponse
	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to parse response: %v", err)})
		return
	}

	c.JSON(201, user)
}

// UpdateImmichUser 更新Immich用户
func UpdateImmichUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	var req ImmichUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	data, err := makeImmichRequest("PUT", "/users/"+userID, req)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update Immich user: %v", err)})
		return
	}

	var user ImmichUserResponse
	if err := json.Unmarshal(data, &user); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to parse response: %v", err)})
		return
	}

	c.JSON(200, user)
}

// DeleteImmichUser 删除Immich用户
func DeleteImmichUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(400, gin.H{"error": "User ID is required"})
		return
	}

	_, err := makeImmichRequest("DELETE", "/users/"+userID, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to delete Immich user: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// BatchUpdateImmichUsers 批量更新Immich用户
func BatchUpdateImmichUsers(c *gin.Context) {
	var req struct {
		UserIDs []string                    `json:"userIds" binding:"required"`
		Updates ImmichUpdateUserRequest     `json:"updates"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	results := make(map[string]interface{})
	errors := make(map[string]string)

	for _, userID := range req.UserIDs {
		data, err := makeImmichRequest("PUT", "/users/"+userID, req.Updates)
		if err != nil {
			errors[userID] = err.Error()
			continue
		}

		var user ImmichUserResponse
		if err := json.Unmarshal(data, &user); err != nil {
			errors[userID] = err.Error()
			continue
		}

		results[userID] = user
	}

	c.JSON(200, gin.H{
		"updated": results,
		"errors":  errors,
		"total":   len(req.UserIDs),
		"success": len(results),
		"failed":  len(errors),
	})
}

// SyncImmichUsersWithSystemUsers 同步Immich用户与系统用户
func SyncImmichUsersWithSystemUsers(c *gin.Context) {
	// 获取系统用户
	systemUsers, err := getSystemUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get system users: %v", err)})
		return
	}

	// 获取Immich用户
	data, err := makeImmichRequest("GET", "/users", nil)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get Immich users: %v", err)})
		return
	}

	var immichUsers []ImmichUserResponse
	if err := json.Unmarshal(data, &immichUsers); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to parse Immich users: %v", err)})
		return
	}

	// 创建系统用户邮箱到Immich用户的映射
	immichUserMap := make(map[string]*ImmichUserResponse)
	for i := range immichUsers {
		immichUserMap[immichUsers[i].Email] = &immichUsers[i]
	}

	results := struct {
		Created    []string `json:"created"`
		Updated    []string `json:"updated"`
		Unmatched  []string `json:"unmatched"`
	}{
		Created:   []string{},
		Updated:   []string{},
		Unmatched: []string{},
	}

	// 同步系统用户到Immich
	for _, sysUser := range systemUsers {
		// 只处理UID >= 1000的普通用户
		uid, _ := parseUID(sysUser.UID)
		if uid < 1000 {
			continue
		}

		// 假设系统用户没有邮箱，使用 username@localhost 作为邮箱
		email := sysUser.Username + "@localhost"

		if immichUser, exists := immichUserMap[email]; exists {
			// 用户存在，检查是否需要更新
			results.Updated = append(results.Updated, immichUser.ID)
		} else {
			// 用户不存在，创建新用户
			// 注意：这里需要默认密码，实际使用时应该让用户设置
			newUser := ImmichUser{
				Email:    email,
				Name:     sysUser.Username,
				Password: sysUser.Username + "123", // 默认密码，实际应该改进
				IsActive: true,
			}

			data, err := makeImmichRequest("POST", "/users", newUser)
			if err != nil {
				results.Unmatched = append(results.Unmatched, sysUser.Username)
				continue
			}

			var createdUser ImmichUserResponse
			if err := json.Unmarshal(data, &createdUser); err != nil {
				results.Unmatched = append(results.Unmatched, sysUser.Username)
				continue
			}

			results.Created = append(results.Created, createdUser.ID)
		}
	}

	c.JSON(200, results)
}