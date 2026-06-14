package api

import (
	"errors"
	"nas-dashboard/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	User         interface{} `json:"user"`
}

// RefreshRequest 刷新 Token 请求
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// RefreshResponse 刷新 Token 响应
type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
}

// Login 处理登录请求
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	authService := service.NewAuthService()
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	user, accessToken, refreshToken, err := authService.Login(req.Username, req.Password, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) || errors.Is(err, service.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else if errors.Is(err, service.ErrUserInactive) {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is inactive"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    24 * 60 * 60, // 24 小时（秒），建议从配置获取
		User:         user,
	})
}

// RefreshToken 刷新 access token
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	authService := service.NewAuthService()
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	accessToken, newRefreshToken, err := authService.RefreshToken(req.RefreshToken, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		}
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		Token:        accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    24 * 60 * 60,
	})
}
