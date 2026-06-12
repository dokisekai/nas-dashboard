package service

import (
	"errors"
	"fmt"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrUserInactive       = errors.New("user is inactive")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrSessionNotFound    = errors.New("session not found")
	ErrMaxLoginAttempts   = errors.New("maximum login attempts exceeded")
)

// Claims JWT 声明
type Claims struct {
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthService 认证服务
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		db: database.GetDB(),
	}
}

// getJWTSecret 从环境变量获取 JWT 密钥
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// 如果没有设置环境变量，使用默认值（生产环境必须设置）
		secret = "nas-dashboard-secret-key-change-in-production"
	}
	return []byte(secret)
}

// getRefreshTokenSecret 从环境变量获取刷新 token 密钥
func getRefreshTokenSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	if secret == "" {
		secret = "nas-dashboard-refresh-secret-key-change-in-production"
	}
	return []byte(secret)
}

// getAccessTokenExpiry 获取访问令牌过期时间
func getAccessTokenExpiry() time.Duration {
	hours := 24 // 默认24小时
	if hoursStr := os.Getenv("JWT_ACCESS_TOKEN_EXPIRY_HOURS"); hoursStr != "" {
		if _, err := fmt.Sscanf(hoursStr, "%d", &hours); err == nil {
			return time.Duration(hours) * time.Hour
		}
	}
	return time.Duration(hours) * time.Hour
}

// getRefreshTokenExpiry 获取刷新令牌过期时间
func getRefreshTokenExpiry() time.Duration {
	days := 30 // 默认30天
	if daysStr := os.Getenv("JWT_REFRESH_TOKEN_EXPIRY_DAYS"); daysStr != "" {
		if _, err := fmt.Sscanf(daysStr, "%d", &days); err == nil {
			return time.Duration(days) * 24 * time.Hour
		}
	}
	return time.Duration(days) * 24 * time.Hour
}

// HashPassword 对密码进行哈希
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(bytes), nil
}

// VerifyPassword 验证密码
func (s *AuthService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// CreateUser 创建新用户
func (s *AuthService) CreateUser(username, password, email, displayName, role string) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// 哈希密码
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        email,
		DisplayName:  displayName,
		Role:         role,
		IsActive:     true,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(username, password, ipAddress, userAgent string) (*models.User, string, string, error) {
	// 查找用户
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", ErrUserNotFound
		}
		return nil, "", "", fmt.Errorf("database error: %w", err)
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, "", "", ErrUserInactive
	}

	// 验证密码
	if !s.VerifyPassword(user.PasswordHash, password) {
		return nil, "", "", ErrInvalidPassword
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	// 生成访问令牌
	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := s.GenerateRefreshToken(&user, ipAddress, userAgent)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &user, accessToken, refreshToken, nil
}

// GenerateAccessToken 生成访问令牌
func (s *AuthService) GenerateAccessToken(user *models.User) (string, error) {
	expiry := getAccessTokenExpiry()
	claims := Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// GenerateRefreshToken 生成刷新令牌
func (s *AuthService) GenerateRefreshToken(user *models.User, ipAddress, userAgent string) (string, error) {
	expiry := getRefreshTokenExpiry()
	expiresAt := time.Now().Add(expiry)

	// 生成刷新令牌字符串
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getRefreshTokenSecret())
	if err != nil {
		return "", err
	}

	// 存储会话到数据库
	session := &models.Session{
		UserID:        user.ID,
		RefreshToken:  tokenString,
		IPAddress:     ipAddress,
		UserAgent:     userAgent,
		ExpiresAt:     expiresAt,
		LastRefreshAt: time.Now(),
		IsActive:      true,
	}

	// 删除该用户的旧会话（可选，如果希望允许多设备登录则删除此部分）
	s.db.Where("user_id = ? AND is_active = ?", user.ID, true).Delete(&models.Session{})

	if err := s.db.Create(session).Error; err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	return tokenString, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshToken string, ipAddress, userAgent string) (string, string, error) {
	// 验证刷新令牌
	if err := s.ValidateRefreshToken(refreshToken); err != nil {
		return "", "", err
	}

	// 查找会话
	var session models.Session
	if err := s.db.Where("refresh_token = ? AND is_active = ?", refreshToken, true).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrSessionNotFound
		}
		return "", "", fmt.Errorf("database error: %w", err)
	}

	// 检查会话是否过期
	if time.Now().After(session.ExpiresAt) {
		session.IsActive = false
		s.db.Save(&session)
		return "", "", ErrTokenExpired
	}

	// 查找用户
	var user models.User
	if err := s.db.First(&user, session.UserID).Error; err != nil {
		return "", "", ErrUserNotFound
	}

	// 检查用户是否激活
	if !user.IsActive {
		return "", "", ErrUserInactive
	}

	// 生成新的访问令牌
	accessToken, err := s.GenerateAccessToken(&user)
	if err != nil {
		return "", "", err
	}

	// 更新会话
	session.LastRefreshAt = time.Now()
	if ipAddress != "" {
		session.IPAddress = ipAddress
	}
	if userAgent != "" {
		session.UserAgent = userAgent
	}
	s.db.Save(&session)

	return accessToken, refreshToken, nil
}

// ValidateAccessToken 验证访问令牌
func (s *AuthService) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateRefreshToken 验证刷新令牌
func (s *AuthService) ValidateRefreshToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return getRefreshTokenSecret(), nil
	})

	if err != nil {
		return ErrInvalidToken
	}

	if _, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return nil
	}

	return ErrInvalidToken
}

// Logout 用户登出
func (s *AuthService) Logout(refreshToken string) error {
	// 将会话标记为非活跃
	var session models.Session
	if err := s.db.Where("refresh_token = ? AND is_active = ?", refreshToken, true).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrSessionNotFound
		}
		return err
	}

	session.IsActive = false
	return s.db.Save(&session).Error
}

// LogoutAll 登出所有设备
func (s *AuthService) LogoutAll(userID uint) error {
	return s.db.Model(&models.Session{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error
}

// GetUserByID 根据用户ID获取用户
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *AuthService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdatePassword 更新用户密码
func (s *AuthService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.VerifyPassword(user.PasswordHash, oldPassword) {
		return ErrInvalidPassword
	}

	// 哈希新密码
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.db.Save(user).Error
}

// CleanExpiredSessions 清理过期会话
func (s *AuthService) CleanExpiredSessions() error {
	return s.db.Model(&models.Session{}).
		Where("expires_at < ?", time.Now()).
		Or("is_active = ?", false).
		Delete(&models.Session{}).Error
}
