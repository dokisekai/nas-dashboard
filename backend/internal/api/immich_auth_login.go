package api

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// ImmichAuthLogin 提供Immich登录凭证
// 这个端点提供Immich的API密钥和URL
func ImmichAuthLogin(c *gin.Context) {
	// 获取请求来源
	origin := c.GetHeader("Origin")
	referer := c.GetHeader("Referer")

	// 简单的同源检查 - 允许从前端和静态文件访问
	allowed := false
	if origin == "https://192.168.50.10:5173" ||
	   origin == "http://192.168.50.10:5173" ||
	   origin == "https://192.168.50.10:8888" ||
	   referer != "" {
		allowed = true
	}

	// 对于同源请求或直接访问，允许通过
	if c.Request.Host == "192.168.50.10:8888" {
		allowed = true
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error": "Origin not allowed",
			})
		return
	}

	// 返回Immich凭证
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"apiKey":   "t5nmlDaFlyz7guxpfmAwvOeiKp7zEC7loF2Ow15V8Q",
		"immichUrl": "https://192.168.50.10:2283",
		"message":  "Immich credentials",
	})
}
