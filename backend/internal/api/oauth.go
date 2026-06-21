package api

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"nas-dashboard/internal/models"
)

// OAuthManager OAuth管理器
type OAuthManager struct {
	db *gorm.DB
}

// NewOAuthManager 创建OAuth管理器
func NewOAuthManager(db *gorm.DB) *OAuthManager {
	return &OAuthManager{db: db}
}

// ServerInfo 服务器信息
type ServerInfo struct {
	IssuerURL           string `json:"issuer_url"`
	AuthorizeEndpoint   string `json:"authorize_endpoint"`
	TokenEndpoint       string `json:"token_endpoint"`
	UserInfoEndpoint    string `json:"userinfo_endpoint"`
	JWKSEndpoint        string `json:"jwks_endpoint"`
	DiscoveryEndpoint   string `json:"discovery_endpoint"`
	Running             bool   `json:"running"`
}

// GetServerInfo 获取服务器信息
func (m *OAuthManager) GetServerInfo(c *gin.Context) {
	// 优先使用环境变量配置的SSO URL
	issuerURL := viper.GetString("sso.issuer_url")
	if issuerURL == "" {
		// 始终使用固定的后端地址，不根据请求动态生成
		// 这样可以避免代理、DNS等问题
		issuerURL = "http://192.168.50.10:8888"

		// 根据当前请求是否使用HTTPS来调整scheme
		if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
			issuerURL = "https://192.168.50.10:8888"
		}
	}

	info := ServerInfo{
		IssuerURL:         issuerURL,
		AuthorizeEndpoint: issuerURL + "/sso/authorize",
		TokenEndpoint:     issuerURL + "/sso/token",
		UserInfoEndpoint:  issuerURL + "/sso/userinfo",
		JWKSEndpoint:      issuerURL + "/sso/jwks",
		DiscoveryEndpoint: issuerURL + "/sso/.well-known/openid-configuration",
		Running:           true, // 可以从配置或状态服务获取实际状态
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    info,
	})
}

// ClientList 客户端列表请求
type ClientList struct {
	Total    int                 `json:"total"`
	Clients  []models.OAuthClient `json:"clients"`
}

// GetClients 获取客户端列表
func (m *OAuthManager) GetClients(c *gin.Context) {
	if m.db == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": ClientList{
				Total:   0,
				Clients: []models.OAuthClient{},
			},
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")

	offset := (page - 1) * limit

	query := m.db.Model(&models.OAuthClient{})
	if search != "" {
		query = query.Where("name LIKE ? OR client_id LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Count(&total)

	var clients []models.OAuthClient
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&clients)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": ClientList{
			Total:   int(total),
			Clients: clients,
		},
	})
}

// CreateClientRequest 创建客户端请求
type CreateClientRequest struct {
	Name         string   `json:"name" binding:"required"`
	RedirectURIs []string `json:"redirect_uris" binding:"required,min=1"`
	GrantTypes   []string `json:"grant_types" binding:"required,min=1"`
	Scopes       []string `json:"scopes" binding:"required,min=1"`
}

// CreateClient 创建客户端
func (m *OAuthManager) CreateClient(c *gin.Context) {
	if m.db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "数据库未初始化",
		})
		return
	}

	var req CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 生成Client ID和Secret
	clientID := generateClientID()
	clientSecret := generateClientSecret()

	// 转换为JSON存储
	redirectURIsJSON, _ := json.Marshal(req.RedirectURIs)
	grantTypesJSON, _ := json.Marshal(req.GrantTypes)
	scopesJSON, _ := json.Marshal(req.Scopes)

	client := models.OAuthClient{
		Name:         req.Name,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURIs: string(redirectURIsJSON),
		GrantTypes:   string(grantTypesJSON),
		Scopes:       string(scopesJSON),
		Status:       "active",
	}

	if err := m.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "创建客户端失败",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    client,
	})
}

// UpdateClientRequest 更新客户端请求
type UpdateClientRequest struct {
	Name         *string   `json:"name"`
	RedirectURIs *[]string `json:"redirect_uris"`
	GrantTypes   *[]string `json:"grant_types"`
	Scopes       *[]string `json:"scopes"`
	Status       *string   `json:"status"`
}

// UpdateClient 更新客户端
func (m *OAuthManager) UpdateClient(c *gin.Context) {
	id := c.Param("id")

	var client models.OAuthClient
	if err := m.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "客户端不存在",
		})
		return
	}

	var req UpdateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// 更新字段
	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.RedirectURIs != nil {
		redirectURIsJSON, _ := json.Marshal(*req.RedirectURIs)
		client.RedirectURIs = string(redirectURIsJSON)
	}
	if req.GrantTypes != nil {
		grantTypesJSON, _ := json.Marshal(*req.GrantTypes)
		client.GrantTypes = string(grantTypesJSON)
	}
	if req.Scopes != nil {
		scopesJSON, _ := json.Marshal(*req.Scopes)
		client.Scopes = string(scopesJSON)
	}
	if req.Status != nil {
		client.Status = *req.Status
	}

	if err := m.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "更新客户端失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    client,
	})
}

// DeleteClient 删除客户端
func (m *OAuthManager) DeleteClient(c *gin.Context) {
	id := c.Param("id")

	var client models.OAuthClient
	if err := m.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "客户端不存在",
		})
		return
	}

	if err := m.db.Delete(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "删除客户端失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "客户端已删除",
	})
}

// RegenerateSecretRequest 重置密钥请求
type RegenerateSecretRequest struct {
	Confirm bool `json:"confirm"`
}

// RegenerateSecret 重置客户端密钥
func (m *OAuthManager) RegenerateSecret(c *gin.Context) {
	id := c.Param("id")

	var req RegenerateSecretRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if !req.Confirm {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "需要确认操作",
		})
		return
	}

	var client models.OAuthClient
	if err := m.db.First(&client, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "客户端不存在",
		})
		return
	}

	// 生成新密钥
	client.ClientSecret = generateClientSecret()

	if err := m.db.Save(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "重置密钥失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"client_id":     client.ClientID,
			"client_secret": client.ClientSecret,
		},
		"message": "密钥已重置，请妥善保存新密钥",
	})
}

// AuthorizationList 授权列表请求
type AuthorizationList struct {
	Total          int                        `json:"total"`
	Authorizations []models.OAuthAuthorization `json:"authorizations"`
}

// GetAuthorizations 获取授权列表
func (m *OAuthManager) GetAuthorizations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	userID := c.Query("user_id")
	clientID := c.Query("client_id")

	offset := (page - 1) * limit

	query := m.db.Model(&models.OAuthAuthorization{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}

	var total int64
	query.Count(&total)

	var authorizations []models.OAuthAuthorization
	query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&authorizations)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": AuthorizationList{
			Total:          int(total),
			Authorizations: authorizations,
		},
	})
}

// RevokeAuthorizationRequest 撤销授权请求
type RevokeAuthorizationRequest struct {
	AuthorizationID uint `json:"authorization_id"`
}

// RevokeAuthorization 撤销授权
func (m *OAuthManager) RevokeAuthorization(c *gin.Context) {
	var req RevokeAuthorizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var auth models.OAuthAuthorization
	if err := m.db.First(&auth, req.AuthorizationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "授权记录不存在",
		})
		return
	}

	now := time.Now()
	auth.RevokedAt = &now

	if err := m.db.Save(&auth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "撤销授权失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "授权已撤销",
	})
}

// ServerStats 服务器统计
type ServerStats struct {
	ActiveUsers   int `json:"active_users"`
	ActiveTokens  int64 `json:"active_tokens"`
	TodayAuths    int64 `json:"today_auths"`
	TotalClients  int64 `json:"total_clients"`
	TotalAuths    int64 `json:"total_auths"`
}

// GetServerStats 获取服务器统计
func (m *OAuthManager) GetServerStats(c *gin.Context) {
	// 检查数据库是否初始化
	if m.db == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": ServerStats{
				ActiveUsers:  0,
				ActiveTokens: 0,
				TodayAuths:   0,
				TotalClients: 0,
				TotalAuths:   0,
			},
		})
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	var stats ServerStats

	// 活跃用户（最近24小时有令牌的用户）
	m.db.Model(&models.OAuthToken{}).
		Where("created_at > ? AND revoked_at IS NULL", today.Add(-24*time.Hour)).
		Select("COUNT(DISTINCT user_id)").
		Scan(&stats.ActiveUsers)

	// 活跃令牌
	m.db.Model(&models.OAuthToken{}).
		Where("expires_at > ? AND revoked_at IS NULL", time.Now()).
		Count((*int64)(&stats.ActiveTokens))

	// 今日认证次数
	m.db.Model(&models.OAuthAuthorization{}).
		Where("created_at >= ? AND created_at < ?", today, tomorrow).
		Count((*int64)(&stats.TodayAuths))

	// 客户端总数
	m.db.Model(&models.OAuthClient{}).
		Count((*int64)(&stats.TotalClients))

	// 总授权数
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return "client-" + generateRandomString(16)
}

// generateClientSecret 生成客户端密钥
func generateClientSecret() string {
	return generateRandomString(32)
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// GetTokenByUser 根据用户获取令牌信息
func (m *OAuthManager) GetTokenByUser(c *gin.Context) {
	userID := c.Param("user_id")

	var tokens []models.OAuthToken
	m.db.Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tokens,
	})
}

// RevokeUserToken 撤销用户令牌
func (m *OAuthManager) RevokeUserToken(c *gin.Context) {
	tokenID := c.Param("token_id")

	var token models.OAuthToken
	if err := m.db.First(&token, tokenID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "令牌不存在",
		})
		return
	}

	now := time.Now()
	token.RevokedAt = &now

	if err := m.db.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "撤销令牌失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "令牌已撤销",
	})
}

// StartServer 启动OAuth服务器
func (m *OAuthManager) StartServer(c *gin.Context) {
	// 这里可以添加启动服务器的逻辑
	// 例如修改配置文件、启动相关服务等

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OAuth服务器已启动",
	})
}

// StopServer 停止OAuth服务器
func (m *OAuthManager) StopServer(c *gin.Context) {
	// 这里可以添加停止服务器的逻辑

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OAuth服务器已停止",
	})
}
