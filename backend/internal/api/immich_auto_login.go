package api

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// ImmichAutoLogin 提供Immich自动登录token
func ImmichAutoLogin(c *gin.Context) {
	// 这里返回Immich API密钥
	// 在实际生产环境中，应该使用更安全的机制
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   "t5nmlDaFlyz7guxpfmAwvOeiKp7zEC7loF2Ow15V8Q",
		"message": "Auto login token generated",
	})
}
