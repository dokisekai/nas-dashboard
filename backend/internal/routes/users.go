package routes

import (
	"github.com/gin-gonic/gin"

	"nas-dashboard/internal/api"
)

// registerUsers 用户、用户组、SMB 用户、共享权限管理路由。
func registerUsers(g *gin.RouterGroup) {
	users := g.Group("/users")
	requireAuth(users)
	{
		users.GET("", api.GetUsers)
		users.POST("", api.CreateUser)
		users.PUT("/:username", api.UpdateUser)
		users.DELETE("/:username", api.DeleteUser)
		users.GET("/:username", api.GetUser)

		// SSH 公钥
		users.GET("/ssh-keys", api.GetSSHKeys)
		users.POST("/ssh-keys", api.AddKey)
		users.DELETE("/ssh-keys/:id", api.DeleteKey)

		// 当前用户
		users.GET("/me", api.GetCurrentUser)
		users.POST("/me/password", api.ChangeCurrentUserPassword)

		// 用户配额
		quotaAPI := api.GetQuotaAPI()
		users.GET("/:username/quota", quotaAPI.GetUserQuota)
		users.PUT("/:username/quota", quotaAPI.SetUserQuota)
	}

	groups := g.Group("/groups")
	requireAuth(groups)
	{
		groups.GET("", api.GetGroups)
		groups.POST("", api.CreateGroup)
		groups.GET("/:name", api.GetGroup)
		groups.PUT("/:name", api.UpdateGroup)
		groups.DELETE("/:name", api.DeleteGroup)
		groups.GET("/:name/members", api.GetGroupMembers)
		groups.POST("/:name/members", api.AddGroupMembers)
		groups.DELETE("/:name/members/:user", api.RemoveGroupMember)
	}

	// SMB 用户与会话管理
	smbUsers := g.Group("/smb")
	requireAuth(smbUsers)
	{
		smbUsers.GET("/users", api.GetSMBUsers)
		smbUsers.POST("/users/:username/password", api.SetSMBPassword)
		smbUsers.DELETE("/users/:username/password", api.DeleteSMBPassword)
		smbUsers.POST("/users/:username/enable", api.EnableSMBUser)
		smbUsers.POST("/users/:username/disable", api.DisableSMBUser)
		smbUsers.GET("/users/:username/stats", api.GetSMBUserStats)
		smbUsers.GET("/sessions", api.GetSMBSessions)
		smbUsers.DELETE("/sessions/:pid", api.DisconnectSMBSession)
		smbUsers.DELETE("/sessions", api.DisconnectAllSMBSessions)
	}

	// 共享与文件权限
	permissions := g.Group("/permissions")
	requireAuth(permissions)
	{
		permissions.GET("/shares", api.GetShares)
		permissions.POST("/shares", api.CreateShare)
		permissions.PUT("/shares/:name", api.UpdateShare)
		permissions.DELETE("/shares/:name", api.DeleteShare)
		permissions.GET("/shares/:name/permissions", api.GetSharePermissions)
		permissions.PUT("/shares/:name/permissions", api.SetSharePermissions)
		permissions.GET("/files", api.GetFilePermissions)
		permissions.PUT("/files/permissions", api.SetFilePermissions)
		permissions.GET("/files/acl", api.GetFileACL)
		permissions.PUT("/files/acl", api.SetFileACL)
	}
}
