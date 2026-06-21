package routes

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nas-dashboard/internal/api"
	"nas-dashboard/internal/sso"
)

// registerSSO 注册 OIDC 标准端点（根路径下）与 /sso 前缀的兼容端点。
//
// 浏览器请求 /authorize 时会自动带上同源 session_token cookie，SSO 的
// AuthorizeHandler 检测到已登录后直接颁发 code 并 302 回客户端配置的
// redirect_uri，全程无需用户在中间页停留，实现"免登录"打开集成的应用。
func registerSSO(r *gin.Engine, ssoServer *sso.SSOServer) {
	if ssoServer == nil {
		log.Println("WARNING: SSO server is nil, SSO/OIDC routes not registered")
		return
	}

	// OIDC 标准端点（在根路径下，符合规范）
	r.GET("/.well-known/openid-configuration", ssoServer.WellKnownHandler)
	r.GET("/authorize", ssoServer.AuthorizeHandler)
	r.GET("/callback", ssoServer.CallbackHandler)
	r.POST("/token", ssoServer.TokenHandler)
	r.GET("/userinfo", ssoServer.UserInfoHandler)
	r.GET("/jwks", ssoServer.JWKSHandler)
	r.POST("/revoke", ssoServer.RevokeTokenHandler)
	r.POST("/introspect", ssoServer.IntrospectHandler)

	// /sso 前缀路由（向后兼容）
	sso := r.Group("/sso")
	{
		sso.GET("/authorize", ssoServer.AuthorizeHandler)
		sso.GET("/callback", ssoServer.CallbackHandler)
		sso.POST("/token", ssoServer.TokenHandler)
		sso.GET("/userinfo", ssoServer.UserInfoHandler)
		sso.GET("/.well-known/openid-configuration", ssoServer.WellKnownHandler)
		sso.GET("/jwks", ssoServer.JWKSHandler)
		sso.POST("/revoke", ssoServer.RevokeTokenHandler)
		sso.POST("/introspect", ssoServer.IntrospectHandler)
	}
}

// registerOAuth OAuth 客户端与令牌管理后台接口。
//
// 这一组路由是供 SSOManager 前端应用使用的 REST 接口，与上面的 OIDC 端点
// 是两回事：这里是"管理 OAuth Provider"，上面是"实现 OAuth Provider"。
func registerOAuth(r *gin.Engine, db *gorm.DB) {
	if db == nil {
		log.Println("WARNING: OAuth management routes initialized with nil database")
	} else {
		log.Println("OAuth management routes initialized with database")
	}

	oauth := r.Group("/api/oauth")
	{
		oauthManager := api.NewOAuthManager(db)

		// 服务器管理
		oauth.GET("/server/info", oauthManager.GetServerInfo)
		oauth.POST("/server/start", oauthManager.StartServer)
		oauth.POST("/server/stop", oauthManager.StopServer)
		oauth.GET("/server/stats", oauthManager.GetServerStats)

		// 客户端管理
		oauth.GET("/clients", oauthManager.GetClients)
		oauth.POST("/clients", oauthManager.CreateClient)
		oauth.PUT("/clients/:id", oauthManager.UpdateClient)
		oauth.DELETE("/clients/:id", oauthManager.DeleteClient)
		oauth.POST("/clients/:id/regenerate-secret", oauthManager.RegenerateSecret)

		// 授权管理
		oauth.GET("/authorizations", oauthManager.GetAuthorizations)
		oauth.POST("/authorizations/revoke", oauthManager.RevokeAuthorization)

		// 令牌管理
		oauth.GET("/users/:user_id/tokens", oauthManager.GetTokenByUser)
		oauth.DELETE("/tokens/:id", oauthManager.RevokeUserToken)
	}
}
