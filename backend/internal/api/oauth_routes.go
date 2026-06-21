package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

// SetupOAuthRoutes 设置OAuth路由
func SetupOAuthRoutes(r *gin.Engine, db *gorm.DB) {
	// 添加调试日志
	if db == nil {
		log.Println("WARNING: OAuth routes initialized with nil database")
	} else {
		log.Println("OAuth routes initialized with database")
	}

	apiGroup := r.Group("/api/oauth")
	// oauth.Use(middleware.Auth()) // 如果需要认证

	{
		oauthManager := NewOAuthManager(db)

		// 服务器管理
		apiGroup.GET("/server/info", oauthManager.GetServerInfo)
		apiGroup.POST("/server/start", oauthManager.StartServer)
		apiGroup.POST("/server/stop", oauthManager.StopServer)
		apiGroup.GET("/server/stats", oauthManager.GetServerStats)

		// 客户端管理
		apiGroup.GET("/clients", oauthManager.GetClients)
		apiGroup.POST("/clients", oauthManager.CreateClient)
		apiGroup.PUT("/clients/:id", oauthManager.UpdateClient)
		apiGroup.DELETE("/clients/:id", oauthManager.DeleteClient)
		apiGroup.POST("/clients/:id/regenerate-secret", oauthManager.RegenerateSecret)

		// 授权管理
		apiGroup.GET("/authorizations", oauthManager.GetAuthorizations)
		apiGroup.POST("/authorizations/revoke", oauthManager.RevokeAuthorization)

		// 令牌管理
		apiGroup.GET("/users/:user_id/tokens", oauthManager.GetTokenByUser)
		apiGroup.DELETE("/tokens/:id", oauthManager.RevokeUserToken)
	}
}