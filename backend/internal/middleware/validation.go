package middleware

import "github.com/gin-gonic/gin"

// SanitizeInput 简化版本
func SanitizeInput(c *gin.Context) {
	c.Next()
}
