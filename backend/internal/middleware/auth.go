package middleware

import (
	"errors"
	"fmt"
	"nas-dashboard/internal/service"

	"github.com/gin-gonic/gin"
)

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("[AUTH DEBUG] Authorization header: %s\n", authHeader[:min(len(authHeader), 50)] + "...")

		if authHeader == "" {
			fmt.Println("[AUTH DEBUG] No Authorization header")
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer token 格式
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			fmt.Println("[AUTH DEBUG] Invalid authorization format")
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenString := authHeader[7:]
		fmt.Printf("[AUTH DEBUG] Token length: %d\n", len(tokenString))

		// 使用service层验证token
		authService := service.NewAuthService()
		claims, err := authService.ValidateAccessToken(tokenString)
		if err != nil {
			fmt.Printf("[AUTH DEBUG] Token validation failed: %v\n", err)
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		fmt.Printf("[AUTH DEBUG] User authenticated: %s (ID: %d)\n", claims.Username, claims.UserID)

		// 将用户信息存储在上下文中
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("username not found in context")
	}
	return username.(string), nil
}

// GetUserID 从上下文中获取用户ID (uint版本)
func GetUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("userID not found in context")
	}
	return userID.(uint), nil
}

// GetUserIDAsUint 从上下文中获取用户ID (uint版本 - 别名)
func GetUserIDAsUint(c *gin.Context) (uint, error) {
	return GetUserID(c)
}
