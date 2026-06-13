package api

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SMBUser SMB用户
type SMBUser struct {
	Username      string    `json:"username"`
	Enabled       bool      `json:"enabled"`
	LastLogin     string    `json:"lastLogin"`
	LoginCount     int       `json:"loginCount"`
	FailedLogins  int       `json:"failedLogins"`
	HomeDirectory string    `json:"homeDirectory"`
	UID           string    `json:"uid"`
	GID           string    `json:"gid"`
}

// SMBSession SMB会话
type SMBSession struct {
	PID         string    `json:"pid"`
	Username    string    `json:"username"`
	Group       string    `json:"group"`
	Machine     string    `json:"machine"`
	IP          string    `json:"ip"`
	Protocol    string    `json:"protocol"`
	ConnectTime string    `json:"connectTime"`
	ClientID    string    `json:"clientId"`
}

// SMBSessionInfo SMB会话信息
type SMBSessionInfo struct {
	Sessions    []SMBSession `json:"sessions"`
	TotalCount  int           `json:"totalCount"`
	Timestamp   time.Time     `json:"timestamp"`
}

// SetSMBPasswordRequest 设置SMB密码请求
type SetSMBPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

// GetSMBUsers 获取SMB用户列表
func GetSMBUsers(c *gin.Context) {
	// 获取系统用户
	systemUsers, err := getSystemUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get users: %v", err)})
		return
	}

	// 转换为SMB用户
	smbUsers := make([]SMBUser, 0)
	for _, user := range systemUsers {
		// 检查用户是否有SMB账户（检查密码是否设置）
		hasSMB := checkUserSMBPassword(user.Username)

		smbUser := SMBUser{
			Username:      user.Username,
			Enabled:       hasSMB,
			LastLogin:     getLastLoginTime(user.Username),
			LoginCount:     getLoginCount(user.Username),
			FailedLogins:  0,
			HomeDirectory: user.Home,
			UID:           user.UID,
			GID:           user.GID,
		}

		smbUsers = append(smbUsers, smbUser)
	}

	c.JSON(200, gin.H{"users": smbUsers})
}

// checkUserSMBPassword 检查用户是否设置了SMB密码
func checkUserSMBPassword(username string) bool {
	cmd := exec.Command("pdbedit", "-L", "-w", username)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return false
	}

	fields := strings.Split(lines[0], ":")
	if len(fields) >= 3 && fields[1] != "" {
		return true
	}

	return false
}

// getLastLoginTime 获取用户最后登录时间
func getLastLoginTime(username string) string {
	cmd := exec.Command("lastlog", "-u", username)
	output, err := cmd.Output()
	if err != nil {
		return "从未登录"
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 && strings.Contains(lines[0], username) {
		fields := strings.Fields(lines[0])
		if len(fields) >= 5 {
			return fmt.Sprintf("%s %s", fields[4], fields[5])
		}
	}

	return "从未登录"
}

// getLoginCount 获取用户登录次数
func getLoginCount(username string) int {
	cmd := exec.Command("last", username, "|", "wc", "-l")
	output, err := cmd.Output()
	if err != nil {
		return 0
	}

	count := strings.TrimSpace(string(output))
	if count == "" {
		return 0
	}

	loginCount, err := strconv.Atoi(count)
	if err != nil {
		return 0
	}

	return loginCount
}

// SetSMBPassword 设置SMB密码
func SetSMBPassword(c *gin.Context) {
	username := c.Param("username")

	var req SetSMBPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 使用 smbpasswd 设置密码
	cmd := exec.Command("smbpasswd", "-a", username)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set password: %v", err)})
		return
	}

	go func() {
		defer stdin.Close()
		stdin.Write([]byte(req.Password + "\n" + req.Password + "\n"))
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set SMB password: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "SMB password set successfully"})
}

// DeleteSMBPassword 删除SMB密码
func DeleteSMBPassword(c *gin.Context) {
	username := c.Param("username")

	cmd := exec.Command("smbpasswd", "-x", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to delete SMB password: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "SMB password deleted successfully"})
}

// EnableSMBUser 启用SMB用户
func EnableSMBUser(c *gin.Context) {
	username := c.Param("username")

	// 检查是否已有SMB密码
	if !checkUserSMBPassword(username) {
		c.JSON(400, gin.H{"error": "User does not have SMB password"})
		return
	}

	cmd := exec.Command("smbpasswd", "-e", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to enable SMB user: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "SMB user enabled successfully"})
}

// DisableSMBUser 禁用SMB用户
func DisableSMBUser(c *gin.Context) {
	username := c.Param("username")

	cmd := exec.Command("smbpasswd", "-d", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to disable SMB user: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "SMB user disabled successfully"})
}

// GetSMBSessions 获取SMB会话
func GetSMBSessions(c *gin.Context) {
	cmd := exec.Command("smbstatus", "-S")
	output, err := cmd.Output()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get SMB sessions: %v", err)})
		return
	}

	sessions := parseSMBStatusOutput(string(output))

	sessionInfo := SMBSessionInfo{
		Sessions:   sessions,
		TotalCount: len(sessions),
		Timestamp:  time.Now(),
	}

	c.JSON(200, sessionInfo)
}

// parseSMBStatusOutput 解析 smbstatus 输出
func parseSMBStatusOutput(output string) []SMBSession {
	var sessions []SMBSession
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "PID") || strings.HasPrefix(line, "---") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 8 {
			session := SMBSession{
				PID:         fields[0],
				Username:    fields[1],
				Group:       fields[2],
				Machine:     fields[3],
				IP:          fields[4],
				Protocol:    fields[6],
				ConnectTime: fields[7],
				ClientID:    "",
			}

			if len(fields) > 8 {
				session.ClientID = fields[8]
			}

			sessions = append(sessions, session)
		}
	}

	return sessions
}

// DisconnectSMBSession 断开SMB会话
func DisconnectSMBSession(c *gin.Context) {
	pid := c.Param("pid")

	cmd := exec.Command("smbcontrol", "close-session", pid)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to disconnect session: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Session disconnected successfully"})
}

// DisconnectAllSMBSessions 断开所有SMB会话
func DisconnectAllSMBSessions(c *gin.Context) {
	cmd := exec.Command("smbcontrol", "close-session", "*")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to disconnect all sessions: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "All sessions disconnected successfully"})
}

// GetSMBUserStats 获取SMB用户统计
func GetSMBUserStats(c *gin.Context) {
	username := c.Param("username")

	// 获取用户信息
	systemUsers, err := getSystemUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get users: %v", err)})
		return
	}

	var targetUser *User
	for _, user := range systemUsers {
		if user.Username == username {
			targetUser = &user
			break
		}
	}

	if targetUser == nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 获取SMB会话
	sessionInfo := getSMBUserSessions(username)

	// 获取文件访问统计
	accessCount := getSMBFileAccessCount(username)

	stats := gin.H{
		"username":     targetUser.Username,
		"enabled":       checkUserSMBPassword(targetUser.Username),
		"sessions":      sessionInfo,
		"accessCount":   accessCount,
		"lastActivity":  getLastLoginTime(targetUser.Username),
	}

	c.JSON(200, stats)
}

// getSMBUserSessions 获取特定用户的SMB会话
func getSMBUserSessions(username string) []SMBSession {
	cmd := exec.Command("smbstatus", "-S")
	output, err := cmd.Output()
	if err != nil {
		return []SMBSession{}
	}

	sessions := parseSMBStatusOutput(string(output))
	var userSessions []SMBSession

	for _, session := range sessions {
		if session.Username == username {
			userSessions = append(userSessions, session)
		}
	}

	return userSessions
}

// getSMBFileAccessCount 获取用户文件访问次数
func getSMBFileAccessCount(username string) int {
	// 这里可以从SMB日志中获取访问统计
	// 简化实现，返回0
	return 0
}