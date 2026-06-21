package sso

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"nas-dashboard/internal/models"
)

// SSOServer OAuth2/OIDC服务器
type SSOServer struct {
	config    *Config
	db        *gorm.DB
	keyMgr    *KeyManager
}

// NewSSOServer 创建SSO服务器
func NewSSOServer(db *gorm.DB) *SSOServer {
	config := GetSSOConfig()

	// 设置默认值
	if config.SSOIssuerURL == "" {
		config.SSOIssuerURL = "http://192.168.50.10:8888"
	}

	if config.SSORedirectURI == "" {
		config.SSORedirectURI = config.SSOIssuerURL + "/sso/callback"
	}

	if config.SSOSecret == "" {
		config.SSOSecret = "your-super-secret-jwt-key-change-this"
	}

	// 初始化用于 id_token 签名的 RSA 密钥管理器。
	// 失败不致命：退化为 nil，generateIDToken 会回退到 HS256（兼容旧客户端）。
	keyMgr, err := NewKeyManager("sso_id_token.key")
	if err != nil {
		log.Printf("Warning: failed to initialize RSA key manager (id_token signing will fall back to HS256): %v", err)
		keyMgr = nil
	}

	return &SSOServer{
		config: config,
		db:     db,
		keyMgr: keyMgr,
	}
}

// AuthorizationRequest 授权请求
type AuthorizationRequest struct {
	ClientID     string `form:"client_id"`
	ResponseType string `form:"response_type"`
	RedirectURI  string `form:"redirect_uri"`
	Scope        string `form:"scope"`
	State        string `form:"state"`
	Provider     string `form:"provider"`
}

// TokenRequest 令牌请求
type TokenRequest struct {
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	GrantType    string `form:"grant_type"`
}

// TokenResponse 令牌响应
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
}

// AuthorizationCode 授权码存储
type AuthorizationCode struct {
	Code        string
	UserID      uint
	ClientID    string
	RedirectURI string
	ExpiresAt   time.Time
	Scope       string
}

// 存储授权码（内存存储，生产环境应使用Redis）
var authorizationCodes = make(map[string]*AuthorizationCode)

// 存储state（内存存储，生产环境应使用Redis）
var states = make(map[string]string)

// 存储客户端信息（从数据库读取）
type ClientInfo struct {
	ClientID     string
	ClientSecret string
	RedirectURIs []string
}

// getClient 从数据库获取客户端信息
func (s *SSOServer) getClient(clientID string) (*ClientInfo, error) {
	if s.db == nil {
		log.Printf("Database not available for client lookup")
		return nil, fmt.Errorf("database not available")
	}

	var client models.OAuthClient
	result := s.db.Where("client_id = ? AND status = ?", clientID, "active").First(&client)
	if result.Error != nil {
		log.Printf("Client lookup error: %v", result.Error)
		return nil, result.Error
	}

	// 解析 JSON 字段
	var redirectURIs []string
	if err := json.Unmarshal([]byte(client.RedirectURIs), &redirectURIs); err != nil {
		log.Printf("Failed to parse redirect URIs: %v", err)
		redirectURIs = []string{}
	}

	log.Printf("Client found: %s with %d redirect URIs", clientID, len(redirectURIs))

	return &ClientInfo{
		ClientID:     client.ClientID,
		ClientSecret: client.ClientSecret,
		RedirectURIs: redirectURIs,
	}, nil
}

// validateRedirectURI 验证重定向URI
func (s *SSOServer) validateRedirectURI(client *ClientInfo, redirectURI string) bool {
	for _, allowedURI := range client.RedirectURIs {
		if allowedURI == redirectURI {
			return true
		}
	}
	return false
}

// AuthorizeHandler 授权端点处理
func (s *SSOServer) AuthorizeHandler(c *gin.Context) {
	var req AuthorizationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	log.Printf("Authorize request: client_id=%s, redirect_uri=%s, provider=%s",
		req.ClientID, req.RedirectURI, req.Provider)

	// 验证客户端
	client, err := s.getClient(req.ClientID)
	if err != nil {
		log.Printf("Client not found or inactive: %s, error: %v", req.ClientID, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client", "error_description": "Client not found or inactive"})
		return
	}

	// 验证重定向URI
	if !s.validateRedirectURI(client, req.RedirectURI) {
		log.Printf("Invalid redirect URI: %s (expected one of: %v)", req.RedirectURI, client.RedirectURIs)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_redirect_uri", "error_description": "Redirect URI mismatch"})
		return
	}

	// 检查用户是否已登录（支持多种方式）
	userID := s.checkUserLoggedIn(c)
	if userID > 0 {
		// 用户已登录，直接处理授权
		s.handleAuthorization(c, &req, userID, client)
		return
	}

	// 生成state
	state, err := GenerateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 存储原始请求信息到state中
	stateData := map[string]interface{}{
		"client_id":     req.ClientID,
		"redirect_uri":  req.RedirectURI,
		"scope":         req.Scope,
		"original_state": req.State,
		"provider":      req.Provider,
	}

	stateJSON, _ := json.Marshal(stateData)
	states[state] = string(stateJSON)

	// 如果指定了身份提供商，重定向到提供商
	if req.Provider != "" {
		ssoManager := NewSSOManager(s.config)
		authURL, err := ssoManager.GenerateAuthorizationURL(req.Provider, state)
		if err != nil {
			log.Printf("Error generating auth URL: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}

		c.Redirect(http.StatusFound, authURL)
		return
	}

	// 否则，显示登录页面让用户选择身份提供商
	// 将原始授权参数传递给登录页面
	providers := NewSSOManager(s.config).GetAvailableProviders()
	c.HTML(http.StatusOK, "sso_login.html", gin.H{
		"providers":     providers,
		"client_id":     req.ClientID,
		"redirect_uri":  req.RedirectURI,
		"scope":         req.Scope,
		"state":         state,
		"original_state": req.State,
	})
}

// checkUserLoggedIn 检查用户是否已登录（支持多种认证方式）
func (s *SSOServer) checkUserLoggedIn(c *gin.Context) uint {
	// 获取JWT密钥（与auth service保持一致）
	jwtSecret := s.getJWTSecret()

	// 方法1: 检查 session cookie (JWT access token)
	sessionToken, err := c.Cookie("session_token")
	if err == nil && sessionToken != "" {
		userID := s.validateJWTToken(sessionToken, jwtSecret)
		if userID > 0 {
			return userID
		}
	}

	// 方法2: 检查 Authorization header (JWT Bearer token)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString := authHeader[7:]
		userID := s.validateJWTToken(tokenString, jwtSecret)
		if userID > 0 {
			return userID
		}
	}

	return 0
}

// getJWTSecret 获取JWT密钥（与auth service保持一致）
func (s *SSOServer) getJWTSecret() []byte {
	secret := s.config.SSOSecret
	if secret == "" || secret == "your-super-secret-jwt-key-change-this" {
		// 使用与auth service相同的默认值
		secret = "nas-dashboard-secret-key-change-in-production"
	}
	return []byte(secret)
}

// validateJWTToken 验证JWT令牌并返回用户ID
func (s *SSOServer) validateJWTToken(tokenString string, jwtSecret []byte) uint {
	claims := &struct {
		UserID   uint   `json:"userId"`
		Username string `json:"username"`
		Role     string `json:"role"`
		jwt.RegisteredClaims
	}{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		log.Printf("JWT validation failed: %v", err)
		return 0
	}

	if token.Valid && claims.UserID > 0 {
		// 检查token是否过期
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			log.Printf("JWT token expired")
			return 0
		}
		return claims.UserID
	}

	return 0
}

// handleAuthorization 处理已登录用户的授权
func (s *SSOServer) handleAuthorization(c *gin.Context, req *AuthorizationRequest, userID uint, client *ClientInfo) {
	// 检查用户是否已经授权过此客户端（可选：可以跳过重复授权）

	// 生成授权码
	authCode := generateAuthCode()

	// 存储授权码
	authorizationCodes[authCode] = &AuthorizationCode{
		Code:        authCode,
		UserID:      userID,
		ClientID:    req.ClientID,
		RedirectURI: req.RedirectURI,
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		Scope:       req.Scope,
	}

	log.Printf("Authorization granted for user_id=%d, client_id=%s", userID, req.ClientID)

	// 重定向到客户端
	params := url.Values{}
	params.Set("code", authCode)
	if req.State != "" {
		params.Set("state", req.State)
	}

	redirectURL := req.RedirectURI + "?" + params.Encode()
	c.Redirect(http.StatusFound, redirectURL)
}

// CallbackHandler 回调端点处理
func (s *SSOServer) CallbackHandler(c *gin.Context) {
	provider := c.Query("provider")
	code := c.Query("code")
	state := c.Query("state")
	error := c.Query("error")

	if error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": error})
		return
	}

	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_state"})
		return
	}

	// 验证state
	if !ValidateState(state) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_state"})
		return
	}

	// 获取存储的state信息
	stateDataJSON, exists := states[state]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_state"})
		return
	}

	var stateData map[string]interface{}
	if err := json.Unmarshal([]byte(stateDataJSON), &stateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_state"})
		return
	}

	// 交换授权码获取用户信息
	ssoManager := NewSSOManager(s.config)
	userInfo, err := ssoManager.ExchangeCodeForToken(provider, code, state)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 创建或更新用户
	if err := ssoManager.CreateOrUpdateUserFromSSO(userInfo, s.db); err != nil {
		log.Printf("Error creating/updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 获取用户信息
	var user models.User
	if err := s.db.Where("sso_provider = ? AND sso_id = ?", userInfo.Provider, userInfo.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 生成授权码
	authCode := generateAuthCode()

	// 存储授权码
	authorizationCodes[authCode] = &AuthorizationCode{
		Code:        authCode,
		UserID:      user.ID,
		ClientID:    stateData["client_id"].(string),
		RedirectURI: stateData["redirect_uri"].(string),
		ExpiresAt:   time.Now().Add(10 * time.Minute),
		Scope:       getStringValue(stateData, "scope"),
	}

	// 删除state
	delete(states, state)

	// 重定向到客户端
	redirectURI := stateData["redirect_uri"].(string)
	params := url.Values{}
	params.Set("code", authCode)
	if originalState, ok := stateData["original_state"].(string); ok && originalState != "" {
		params.Set("state", originalState)
	}

	c.Redirect(http.StatusFound, redirectURI+"?"+params.Encode())
}

// TokenHandler 令牌端点处理
func (s *SSOServer) TokenHandler(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBind(&req); err != nil {
		log.Printf("Token request bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	// 记录完整的请求信息
	body, _ := c.GetRawData()
	log.Printf("Token request: client_id=%s, grant_type=%s, code_length=%d, code_prefix=%s, body=%s",
		req.ClientID, req.GrantType, len(req.Code), safePrefix(req.Code), string(body))

	// 验证客户端
	client, err := s.getClient(req.ClientID)
	if err != nil {
		log.Printf("Client not found: %s", req.ClientID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}

	if client.ClientSecret != req.ClientSecret {
		log.Printf("Invalid client secret for: %s", req.ClientID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_client"})
		return
	}

	// 验证授权码
	authCodeData, exists := authorizationCodes[req.Code]
	if !exists {
		log.Printf("Invalid authorization code: len=%d, prefix=%s, looking_for_prefix=%s", len(req.Code), safePrefix(req.Code), getStoredCodePrefix())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	// 检查授权码是否过期
	if time.Now().After(authCodeData.ExpiresAt) {
		delete(authorizationCodes, req.Code)
		log.Printf("Authorization code expired")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	// 验证重定向URI
	if !s.validateRedirectURI(client, req.RedirectURI) {
		log.Printf("Invalid redirect URI in token request: %s", req.RedirectURI)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
		return
	}

	// 删除授权码（一次性使用）
	delete(authorizationCodes, req.Code)

	// 生成访问令牌和刷新令牌
	accessToken, err := s.generateAccessToken(authCodeData.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	refreshToken, err := s.generateRefreshToken(authCodeData.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	// 生成ID令牌（如果scope包含openid）
	idToken := ""
	if strings.Contains(authCodeData.Scope, "openid") {
		idToken, err = s.generateIDToken(authCodeData.UserID, authCodeData.ClientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}
	}

	response := TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    3600,
		RefreshToken: refreshToken,
		IDToken:      idToken,
	}

	log.Printf("Token generated successfully for user_id: %d", authCodeData.UserID)
	c.JSON(http.StatusOK, response)
}

// UserInfoHandler 用户信息端点处理
func (s *SSOServer) UserInfoHandler(c *gin.Context) {
	// 验证访问令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_authorization"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.validateAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	// 获取用户信息
	var user models.User
	if err := s.db.First(&user, claims["user_id"]).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	userInfo := gin.H{
		"sub":           fmt.Sprintf("%d", user.ID),
		"name":          user.DisplayName,
		"email":         user.Email,
		"preferred_username": user.Username,
	}

	c.JSON(http.StatusOK, userInfo)
}

// WellKnownHandler OIDC发现端点
func (s *SSOServer) WellKnownHandler(c *gin.Context) {
	// 根据请求确定协议
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	// 获取主机
	host := c.Request.Host
	if host == "" {
		host = "192.168.50.10:8888"
	}

	baseURL := scheme + "://" + host

	wellKnown := gin.H{
		"issuer":                  baseURL,
		"authorization_endpoint": baseURL + "/authorize",
		"token_endpoint":         baseURL + "/token",
		"userinfo_endpoint":      baseURL + "/userinfo",
		"jwks_uri":               baseURL + "/jwks",
		"response_types_supported": []string{"code"},
		"grant_types_supported":    []string{"authorization_code", "refresh_token"},
		"subject_types_supported":  []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256", "HS256"},
		"scopes_supported": []string{"openid", "profile", "email"},
		"claims_supported": []string{"sub", "name", "email", "preferred_username"},
	}

	c.JSON(http.StatusOK, wellKnown)
}

// JWKSHandler JSON Web Key Set端点
//
// 返回用于验证 id_token 签名的 RSA 公钥集合。
// Immich / oauth4webapi 等严格 OIDC 客户端在验签前会从此端点拉取公钥。
// 如果 RSA 密钥管理器未初始化（回退到 HS256），则返回空 JWKS。
func (s *SSOServer) JWKSHandler(c *gin.Context) {
	keys := []gin.H{}
	if s.keyMgr != nil {
		jwk := s.keyMgr.JWK()
		keys = append(keys, gin.H(jwk))
	}

	c.JSON(http.StatusOK, gin.H{
		"keys": keys,
	})
}

// generateAccessToken 生成访问令牌
func (s *SSOServer) generateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SSOSecret))
}

// generateRefreshToken 生成刷新令牌
func (s *SSOServer) generateRefreshToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
		"type":    "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SSOSecret))
}

// generateIDToken 生成ID令牌
//
// OIDC 规范允许 HS256 或 RS256，但严格客户端（Immich/oauth4webapi、
// openid-client 等）只接受非对称签名并通过 JWKS 验签。
// 因此优先使用 RS256；仅在 RSA 密钥不可用时回退到 HS256。
func (s *SSOServer) generateIDToken(userID uint, clientID string) (string, error) {
	// 获取用户信息
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"iss":                s.config.SSOIssuerURL,
		"sub":                fmt.Sprintf("%d", user.ID),
		"aud":                clientID, // 使用实际的客户端ID
		"exp":                time.Now().Add(time.Hour).Unix(),
		"iat":                time.Now().Unix(),
		"name":               user.DisplayName,
		"email":              user.Email,
		"preferred_username": user.Username,
	}

	if s.keyMgr != nil {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = s.keyMgr.KeyID()
		return token.SignedString(s.keyMgr.PrivateKey())
	}

	// 回退路径：RSA 密钥不可用时使用 HS256（旧客户端兼容）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SSOSecret))
}

// validateAccessToken 验证访问令牌
func (s *SSOServer) validateAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.SSOSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// generateAuthCode 生成授权码
func generateAuthCode() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		log.Printf("Error generating auth code: %v", err)
		return ""
	}

	code := base64.URLEncoding.EncodeToString(b)
	return strings.TrimRight(code, "=")
}

// getStringValue 从map中获取string值
func getStringValue(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// RevokeTokenHandler 撤销令牌端点
func (s *SSOServer) RevokeTokenHandler(c *gin.Context) {
	// 验证访问令牌
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing_authorization"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.validateAccessToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid_token"})
		return
	}

	// 获取用户ID
	userID := uint(claims["user_id"].(float64))

	// 撤销用户的所有会话
	if err := s.db.Where("user_id = ?", userID).Delete(&models.Session{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}

// IntrospectHandler 令牌内省端点
func (s *SSOServer) IntrospectHandler(c *gin.Context) {
	// 获取令牌
	token := c.PostForm("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_token"})
		return
	}

	// 验证令牌
	claims, err := s.validateAccessToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"active": false,
		})
		return
	}

	// 获取用户信息
	userID := uint(claims["user_id"].(float64))
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"active": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"active":    true,
		"client_id": "nas-dashboard",
		"username":  user.Username,
		"exp":       claims["exp"],
		"iat":       claims["iat"],
		"sub":       fmt.Sprintf("%d", user.ID),
	})
}

// safePrefix returns safe prefix of a string for logging
func safePrefix(s string) string {
	if len(s) <= 10 {
		return s
	}
	return s[:10] + "..."
}

// getStoredCodePrefix returns prefix of first stored code for debugging
func getStoredCodePrefix() string {
	for code := range authorizationCodes {
		return safePrefix(code)
	}
	return "none"
}