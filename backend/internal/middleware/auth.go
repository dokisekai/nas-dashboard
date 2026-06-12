package middleware

import (
	"errors"
	"fmt"
	"nas-dashboard/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Bearer token 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		// 验证 token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 将用户信息存储在上下文中
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// GetUsername 从上下文中获取用户名
func GetUsername(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("username not found in context")
	}
	return username.(string), nil
}

// GetUserID 从上下文中获取用户ID (string版本)
func GetUserID(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", errors.New("userID not found in context")
	}
	return userID.(string), nil
}

// GetUserIDAsUint 从上下文中获取用户ID (uint版本)
func GetUserIDAsUint(c *gin.Context) (uint, error) {
	userIDStr, err := GetUserID(c)
	if err != nil {
		return 0, err
	}
	// 这里需要转换string到uint，简化处理，实际应该从JWT claims中获取
	// 为了兼容性，我们假设userID是数字字符串
	var userID uint
	_, err = fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return 0, errors.New("invalid userID format")
	}
	return userID, nil
}
