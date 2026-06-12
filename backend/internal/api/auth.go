package api

import (
	"errors"
	"nas-dashboard/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
	libjwt "github.com/golang-jwt/jwt/v5"
)

// 用户数据库 (模拟，实际应使用真实数据库)
var users = map[string]string{
	"admin": "admin123",
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expiresIn"`
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证用户名和密码
	if password, exists := users[req.Username]; !exists || password != req.Password {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	// 生成 access token (有效期 24 小时)
	accessToken, err := jwt.GenerateToken(req.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	// 生成 refresh token (有效期 30 天)
	refreshToken, err := generateRefreshToken(req.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(200, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 小时（秒）
	})
}

// RefreshToken 刷新 access token
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证 refresh token
	claims, err := parseRefreshToken(req.RefreshToken)
	if err != nil {
		if errors.Is(err, libjwt.ErrTokenExpired) {
			c.JSON(401, gin.H{"error": "Refresh token expired, please login again"})
			return
		}
		c.JSON(401, gin.H{"error": "Invalid refresh token"})
		return
	}

	// 生成新的 access token
	newToken, err := jwt.GenerateToken(claims.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate new token"})
		return
	}

	c.JSON(200, RefreshResponse{
		Token:     newToken,
		ExpiresIn: 24 * 60 * 60, // 24 小时（秒）
	})
}

// generateRefreshToken 生成 refresh token
// refresh token 使用更长的过期时间和独立的密钥（可选）
func generateRefreshToken(username string) (string, error) {
	// 使用自定义的 JWT 生成逻辑
	// 在生产环境中，建议使用不同的密钥或添加额外的标识
	return jwt.GenerateTokenWithExpiry(username, 30*24*time.Hour)
}

// parseRefreshToken 解析 refresh token
func parseRefreshToken(tokenString string) (*jwt.Claims, error) {
	// 复用 jwt 包的解析功能
	// 在生产环境中，可以添加额外的验证逻辑
	return jwt.ParseToken(tokenString)
}

// ValidateToken 验证 token（供中间件使用）
func ValidateToken(tokenString string) (*jwt.Claims, error) {
	return jwt.ParseToken(tokenString)
}
