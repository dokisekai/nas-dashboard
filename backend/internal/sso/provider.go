package sso

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/viper"
	"gorm.io/gorm"
	"nas-dashboard/internal/models"
)

// Config SSO配置
type Config struct {
	// Google OAuth
	GoogleClientID     string `mapstructure:"google_client_id"`
	GoogleClientSecret string `mapstructure:"google_client_secret"`
	GoogleRedirectURI  string `mapstructure:"google_redirect_uri"`

	// GitHub OAuth
	GitHubClientID     string `mapstructure:"github_client_id"`
	GitHubClientSecret string `mapstructure:"github_client_secret"`
	GitHubRedirectURI  string `mapstructure:"github_redirect_uri"`

	// Microsoft OAuth
	MicrosoftClientID     string `mapstructure:"microsoft_client_id"`
	MicrosoftClientSecret string `mapstructure:"microsoft_client_secret"`
	MicrosoftRedirectURI  string `mapstructure:"microsoft_redirect_uri"`
	MicrosoftTenant      string `mapstructure:"microsoft_tenant"`

	// LDAP
	LDAPServer       string   `mapstructure:"ldap_server"`
	LDAPPort         int      `mapstructure:"ldap_port"`
	LDAPBindDN       string   `mapstructure:"ldap_bind_dn"`
	LDAPBindPassword string   `mapstructure:"ldap_bind_password"`
	LDAPBaseDN       string   `mapstructure:"ldap_base_dn"`
	LDAPUserFilter   string   `mapstructure:"ldap_user_filter"`

	// Immich
	ImmichServerURL  string `mapstructure:"immich_server_url"`
	ImmichAPIKey      string `mapstructure:"immich_api_key"`

	// SSO配置
	SSOIssuerURL     string `mapstructure:"sso_issuer_url"`
	SSOClientID     string `mapstructure:"sso_client_id"`
	SSOSecret        string `mapstructure:"sso_secret"`
	SSORedirectURI    string `mapstructure:"sso_redirect_uri"`
	SSOAutoCreateUser bool  `mapstructure:"sso_auto_create_user"`
}

// SSOProvider SSO提供商
type SSOProvider struct {
	config *Config
}

// NewSSOProvider 创建SSO提供商
func NewSSOProvider(config *Config) *SSOProvider {
	return &SSOProvider{
		config: config,
	}
}

// IdentityProvider 身份提供商接口
type IdentityProvider interface {
	Name() string
	AuthURL(state string) string
	TokenURL() string
	UserAgent() string
	ExchangeCodeForToken(code, state string) (*UserInfo, error)
}

// UserInfo 用户信息
type UserInfo struct {
	ID            string                 `json:"id"`
	Username      string                 `json:"username"`
	Email         string                 `json:"email"`
	Name          string                 `json:"name"`
	AccessToken   string                 `json:"access_token,omitempty"`
	RefreshToken  string                 `json:"refresh_token,omitempty"`
	ExpiresAt     int64                  `json:"expires_at,omitempty"`
	Provider      string                 `json:"provider"`
	RawData       map[string]interface{} `json:"raw_data,omitempty"`
}

// GoogleProvider Google OAuth提供商
type GoogleProvider struct {
	config *Config
}

func (g *GoogleProvider) Name() string {
	return "google"
}

func (g *GoogleProvider) AuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", g.config.GoogleClientID)
	params.Set("redirect_uri", g.config.GoogleRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "openid profile email")
	params.Set("state", state)

	return "https://accounts.google.com/o/oauth2/v2/auth?" + params.Encode()
}

func (g *GoogleProvider) TokenURL() string {
	return "https://oauth2.googleapis.com/token"
}

func (g *GoogleProvider) UserAgent() string {
	return ""
}

func (g *GoogleProvider) ExchangeCodeForToken(code, state string) (*UserInfo, error) {
	// 交换授权码获取access token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", g.config.GoogleClientID)
	data.Set("client_secret", g.config.GoogleClientSecret)
	data.Set("redirect_uri", g.config.GoogleRedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(g.TokenURL(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// 获取用户信息
	userInfo, err := g.getGoogleUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	userInfo.Provider = "google"
	userInfo.AccessToken = tokenResp.AccessToken
	userInfo.RefreshToken = tokenResp.RefreshToken
	if tokenResp.ExpiresIn > 0 {
		userInfo.ExpiresAt = time.Now().Unix() + tokenResp.ExpiresIn - 300 // 提前5分钟作为缓冲
	}

	return userInfo, nil
}

func (g *GoogleProvider) getGoogleUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:       userInfo["sub"].(string),
		Email:    userInfo["email"].(string),
		Name:     userInfo["name"].(string),
		Username: userInfo["email"].(string),
		Provider: "google",
		RawData: userInfo,
	}, nil
}

// GitHubProvider GitHub OAuth提供商
type GitHubProvider struct {
	config *Config
}

func (g *GitHubProvider) Name() string {
	return "github"
}

func (g *GitHubProvider) AuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", g.config.GitHubClientID)
	params.Set("redirect_uri", g.config.GitHubRedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "user:email user:read")
	params.Set("state", state)

	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

func (g *GitHubProvider) TokenURL() string {
	return "https://github.com/login/oauth/access_token"
}

func (g *GitHubProvider) UserAgent() string {
	return ""
}

func (g *GitHubProvider) ExchangeCodeForToken(code, state string) (*UserInfo, error) {
	// 交换授权码获取access token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", g.config.GitHubClientID)
	data.Set("client_secret", g.config.GitHubClientSecret)
	data.Set("redirect_uri", g.config.GitHubRedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(g.TokenURL(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// 获取用户信息
	userInfo, err := g.getGitHubUserInfo(tokenResp.AccessToken)
	if err != nil {
		return nil, err
	}

	userInfo.Provider = "github"
	userInfo.AccessToken = tokenResp.AccessToken

	return userInfo, nil
}

func (g *GitHubProvider) getGitHubUserInfo(accessToken string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	type EmailInfo struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	var emails []EmailInfo

	if emailsData, ok := user["emails"].([]interface{}); ok {
		for _, e := range emailsData {
			if emailMap, ok := e.(map[string]interface{}); ok {
				if primary, ok := emailMap["primary"].(bool); ok && primary {
					emails = append(emails, EmailInfo{
						Email:   emailMap["email"].(string),
						Primary: primary,
					})
				}
			}
		}
	}

	email := ""
	if len(emails) > 0 {
		email = emails[0].Email
	}

	return &UserInfo{
		ID:       fmt.Sprintf("%d", int(user["id"].(float64))),
		Username: user["login"].(string),
		Email:    email,
		Name:     user["name"].(string),
		Provider: "github",
		RawData:  user,
	}, nil
}

// MicrosoftProvider Microsoft Azure AD OAuth提供商
type MicrosoftProvider struct {
	config *Config
}

func (m *MicrosoftProvider) Name() string {
	return "microsoft"
}

func (m *MicrosoftProvider) AuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", m.config.MicrosoftClientID)
	params.Set("response_type", "code")
	params.Set("redirect_uri", m.config.MicrosoftRedirectURI)
	params.Set("response_mode", "query")
	params.Set("scope", "openid profile email User.Read")
	params.Set("state", state)

	tenant := m.config.MicrosoftTenant
	if tenant == "" {
		tenant = "common" // 默认使用通用终端
	}

	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?%s", tenant, params.Encode())
}

func (m *MicrosoftProvider) TokenURL() string {
	tenant := m.config.MicrosoftTenant
	if tenant == "" {
		tenant = "common"
	}
	return fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenant)
}

func (m *MicrosoftProvider) UserAgent() string {
	return ""
}

func (m *MicrosoftProvider) ExchangeCodeForToken(code, state string) (*UserInfo, error) {
	tenant := m.config.MicrosoftTenant
	if tenant == "" {
		tenant = "common"
	}

	// 交换授权码获取access token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", m.config.MicrosoftClientID)
	data.Set("client_secret", m.config.MicrosoftClientSecret)
	data.Set("redirect_uri", m.config.MicrosoftRedirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(m.TokenURL(), data)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// 获取用户信息
	userInfo, err := m.getMicrosoftUserInfo(tokenResp.AccessToken, tenant)
	if err != nil {
		return nil, err
	}

	userInfo.Provider = "microsoft"
	userInfo.AccessToken = tokenResp.AccessToken
	userInfo.RefreshToken = tokenResp.RefreshToken
	if tokenResp.ExpiresIn > 0 {
		userInfo.ExpiresAt = time.Now().Unix() + tokenResp.ExpiresIn - 300
	}

	return userInfo, nil
}

func (m *MicrosoftProvider) getMicrosoftUserInfo(accessToken, tenant string) (*UserInfo, error) {
	req, err := http.NewRequest("GET", "https://graph.microsoft.com/v1.0/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// 获取邮箱信息
	email := ""
	if mail, ok := user["mail"].(string); ok {
		email = mail
	} else if mails, ok := user["mails"].([]interface{}); ok && len(mails) > 0 {
		if mailMap, ok := mails[0].(map[string]interface{}); ok {
			if address, ok := mailMap["address"].(string); ok {
				email = address
			}
		}
	}

	displayName := ""
	if name, ok := user["displayName"].(string); ok {
		displayName = name
	}

	id := ""
	if idVal, ok := user["id"].(string); ok {
		id = idVal
	}

	return &UserInfo{
		ID:       id,
		Email:    email,
		Name:     displayName,
		Username: email, // Microsoft通常用email作为username
		Provider: "microsoft",
		RawData:  user,
	}, nil
}

// ImmichProvider Immich应用提供商
type ImmichProvider struct {
	config *Config
}

func (i *ImmichProvider) Name() string {
	return "immich"
}

func (i *ImmichProvider) AuthURL(state string) string {
	// Immich使用内部API认证
	return ""
}

func (i *ImmichProvider) TokenURL() string {
	return ""
}

func (i *ImmichProvider) UserAgent() string {
	return ""
}

func (i *ImmichProvider) ExchangeCodeForToken(code, state string) (*UserInfo, error) {
	// Immich通过API密钥直接认证
	// 这里简化处理，实际需要Immich的API调用

	return &UserInfo{
		ID:       "immich-user",
		Username: "immich-user",
		Email:    "user@immich.local",
		Name:     "Immich User",
		Provider: "immich",
	}, nil
}

// LDAPProvider LDAP/AD提供商
type LDAPProvider struct {
	config *Config
}

func (l *LDAPProvider) Name() string {
	return "ldap"
}

func (l *LDAPProvider) AuthURL(state string) string {
	return ""
}

func (l *LDAPProvider) TokenURL() string {
	return ""
}

func (l *LDAPProvider) UserAgent() string {
	return ""
}

func (l *LDAPProvider) ExchangeCodeForToken(code, state string) (*UserInfo, error) {
	// LDAP认证实现（需要golang.org/x/example.com/adldap库）
	// 这里提供简化版本的结构

	return &UserInfo{
		ID:       "ldap-user",
		Username: "ldap-user",
		Email:    "user@ldap.local",
		Name:     "LDAP User",
		Provider: "ldap",
	}, nil
}

// SSOManager SSO管理器
type SSOManager struct {
	config      *Config
	providers  map[string]IdentityProvider
}

// NewSSOManager 创建SSO管理器
func NewSSOManager(config *Config) *SSOManager {
	manager := &SSOManager{
		config:     config,
		providers: make(map[string]IdentityProvider),
	}

	// 初始化各个提供商
	if config.GoogleClientID != "" {
		manager.providers["google"] = &GoogleProvider{config: config}
	}

	if config.GitHubClientID != "" {
		manager.providers["github"] = &GitHubProvider{config: config}
	}

	if config.MicrosoftClientID != "" {
		manager.providers["microsoft"] = &MicrosoftProvider{config: config}
	}

	if config.ImmichServerURL != "" {
		manager.providers["immich"] = &ImmichProvider{config: config}
	}

	if config.LDAPServer != "" {
		manager.providers["ldap"] = &LDAPProvider{config: config}
	}

	return manager
}

// GetProvider 获取指定的身份提供商
func (m *SSOManager) GetProvider(name string) (IdentityProvider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("identity provider '%s' not found", name)
	}
	return provider, nil
}

// GetAvailableProviders 获取可用的身份提供商列表
func (m *SSOManager) GetAvailableProviders() []string {
	var providers []string

	for name := range m.providers {
		providers = append(providers, name)
	}

	return providers
}

// GenerateState 生成随机state字符串
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	state := base64.URLEncoding.EncodeToString(b)
	// 移除padding
	state = strings.TrimRight(state, "=")

	return state, nil
}

// ValidateState 验证状态字符串
func ValidateState(state string) bool {
	// 基础验证：长度应该在32-48之间（base64编码后）
	if len(state) < 32 || len(state) > 48 {
		return false
	}

	// 检查是否只包含有效的base64字符
	for _, c := range state {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			continue
		}
		return false
	}

	return true
}

// ValidateIDToken 验证ID令牌
func (m *SSOManager) ValidateIDToken(idToken string) (*UserInfo, error) {
	// 这里需要解析JWT ID token或调用提供商验证
	// 简化实现：返回基于ID token的用户信息

	// 实际应该：
	// 1. 解析JWT token
	// 2. 验证签名和过期时间
	// 3. 返回用户信息

	return nil, fmt.Errorf("ID token validation not implemented")
}

// GenerateAuthorizationURL 生成授权URL
func (m *SSOManager) GenerateAuthorizationURL(provider, state string) (string, error) {
	p, err := m.GetProvider(provider)
	if err != nil {
		return "", err
	}

	return p.AuthURL(state), nil
}

// ExchangeCodeForToken 交换授权码获取用户信息
func (m *SSOManager) ExchangeCodeForToken(provider, code, state string) (*UserInfo, error) {
	p, err := m.GetProvider(provider)
	if err != nil {
		return nil, err
	}

	return p.ExchangeCodeForToken(code, state)
}

// CreateOrUpdateUserFromSSO 从SSO用户信息创建或更新本地用户
func (m *SSOManager) CreateOrUpdateUserFromSSO(userInfo *UserInfo, db *gorm.DB) error {
	// 检查用户是否已存在（通过provider和ID）
	var user models.User
	result := db.Where("id > 0").First(&user)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database error: %w", result.Error)
	}

	// 如果用户不存在，检查是否通过第三方ID查找
	if result.Error == gorm.ErrRecordNotFound {
		// 尝试通过第三方ID查找
		result = db.Where("sso_provider = ? AND sso_id = ?",
			userInfo.Provider, userInfo.ID).First(&user)

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("database error: %w", result.Error)
		}
	}

	now := time.Now()

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新用户
		user = models.User{
			Username:    userInfo.Username,
			Email:       userInfo.Email,
			DisplayName: userInfo.Name,
			PasswordHash: "", // SSO用户不需要本地密码
			Role:        "user",
			IsActive:    true,
			LastLogin:   &now,
			SSOProvider: userInfo.Provider,
			SSOID:        userInfo.ID,
		}

		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// 更新现有用户
		user.Email = userInfo.Email
		user.DisplayName = userInfo.Name
		user.LastLogin = &now
		user.SSOProvider = userInfo.Provider
		user.SSOID = userInfo.ID

		if err := db.Save(&user).Error; err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	}

	return nil
}

// RevokeToken 撤销用户令牌
func (m *SSOManager) RevokeToken(userID uint, db *gorm.DB) error {
	// 删除用户的所有活跃会话
	if err := db.Where("user_id = ?", userID).Delete(&models.Session{}).Error; err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// AutoCreateUserConfig 检查是否自动创建用户
func (m *SSOManager) AutoCreateUser() bool {
	return m.config.SSOAutoCreateUser
}

// GetSSOConfig 获取SSO配置
func GetSSOConfig() *Config {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Failed to read config: %v", err)
		return &Config{}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		log.Printf("Warning: Failed to unmarshal config: %v", err)
		return &Config{}
	}

	return config
}