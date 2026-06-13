package api

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/ssh"
)

// User 用户信息
type User struct {
	Username string `json:"username"`
	UID      string `json:"uid"`
	GID      string `json:"gid"`
	Group    string `json:"group"`
	Home     string `json:"home"`
	Shell    string `json:"shell"`
	Comment  string `json:"comment"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Comment  string `json:"comment"`
	Group    string `json:"group"`
	Shell    string `json:"shell"`
	Home     string `json:"home"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Password string `json:"password"`
	Group    string `json:"group"`
	Shell    string `json:"shell"`
	Comment  string `json:"comment"`
}

// SSHKey SSH 密钥
type SSHKey struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Fingerprint string `json:"fingerprint"`
	Type        string `json:"type"`
	Content     string `json:"content,omitempty"`
	AddedAt     string `json:"addedAt"`
	User        string `json:"user"`
}

// SSHKeyRequest SSH 密钥请求
type SSHKeyRequest struct {
	Name    string `json:"name" binding:"required"`
	Content string `json:"content" binding:"required"`
	User    string `json:"user" binding:"required"`
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	users, err := getSystemUsers()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get users: %v", err)})
		return
	}
	c.JSON(200, gin.H{"users": users})
}

// getSystemUsers 从 /etc/passwd 获取系统用户
func getSystemUsers() ([]User, error) {
	// 读取 /etc/passwd 文件
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("failed to read /etc/passwd: %w", err)
	}

	return parsePasswd(data)
}

// parsePasswd 解析 /etc/passwd 文件
func parsePasswd(data []byte) ([]User, error) {
	var users []User
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) < 7 {
			continue
		}

		username := fields[0]
		uid := fields[2]
		gid := fields[3]
		comment := fields[4]
		home := fields[5]
		shell := fields[6]

		// 获取组名称
		group := getGroupName(gid)

		user := User{
			Username: username,
			UID:      uid,
			GID:      gid,
			Group:    group,
			Comment:  comment,
			Home:     home,
			Shell:    shell,
		}

		users = append(users, user)
	}

	return users, nil
}

// getGroupName 根据 GID 获取组名称
func getGroupName(gid string) string {
	// 读取 /etc/group
	data, err := os.ReadFile("/etc/group")
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) >= 3 && fields[2] == gid {
			return fields[0]
		}
	}

	return ""
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {
	username := c.Param("username")

	systemUser, err := user.Lookup(username)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	group := getGroupName(systemUser.Gid)

	userInfo := User{
		Username: systemUser.Username,
		UID:      systemUser.Uid,
		GID:      systemUser.Gid,
		Group:    group,
		Home:     systemUser.HomeDir,
		Shell:    systemUser.Username, // user.Lookup 不提供 Shell，需要额外获取
	}

	// 从 /etc/passwd 获取 Shell
	if shell := getUserShell(systemUser.Username); shell != "" {
		userInfo.Shell = shell
	}

	c.JSON(200, gin.H{"user": userInfo})
}

// getUserShell 从 /etc/passwd 获取用户 Shell
func getUserShell(username string) string {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, username+":") {
			fields := strings.Split(line, ":")
			if len(fields) >= 7 {
				return fields[6]
			}
		}
	}

	return ""
}

// CreateUser 创建用户
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查用户是否已存在
	if _, err := user.Lookup(req.Username); err == nil {
		c.JSON(409, gin.H{"error": "User already exists"})
		return
	}

	// 构建 useradd 命令
	args := []string{"-m", req.Username}

	// 设置 Shell
	if req.Shell != "" {
		args = append([]string{"-s", req.Shell}, args...)
	} else {
		args = append([]string{"-s", "/bin/bash"}, args...)
	}

	// 设置主目录
	if req.Home != "" {
		args = append([]string{"-d", req.Home}, args...)
	}

	// 设置注释
	if req.Comment != "" {
		args = append([]string{"-c", req.Comment}, args...)
	}

	// 添加到组
	if req.Group != "" {
		args = append([]string{"-G", req.Group}, args...)
	}

	cmd := exec.Command("useradd", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create user: %v, output: %s", err, string(output))})
		return
	}

	// 设置密码
	if err := setUserPassword(req.Username, req.Password); err != nil {
		// 如果设置密码失败，删除刚创建的用户
		exec.Command("userdel", req.Username).Run()
		c.JSON(500, gin.H{"error": fmt.Sprintf("User created but failed to set password: %v", err)})
		return
	}

	c.JSON(201, gin.H{"message": "User created successfully", "username": req.Username})
}

// setUserPassword 设置用户密码
func setUserPassword(username, password string) error {
	// 使用 chpasswd 命令设置密码
	cmd := exec.Command("chpasswd")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	fmt.Fprintf(stdin, "%s:%s\n", username, password)
	stdin.Close()

	return cmd.Wait()
}

// UpdateUser 更新用户
func UpdateUser(c *gin.Context) {
	username := c.Param("username")
	var req UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查用户是否存在
	if _, err := user.Lookup(username); err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 更新密码
	if req.Password != "" {
		if err := setUserPassword(username, req.Password); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update password: %v", err)})
			return
		}
	}

	// 更新 Shell
	if req.Shell != "" {
		cmd := exec.Command("usermod", "-s", req.Shell, username)
		if output, err := cmd.CombinedOutput(); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update shell: %v, output: %s", err, string(output))})
			return
		}
	}

	// 更新注释
	if req.Comment != "" {
		cmd := exec.Command("usermod", "-c", req.Comment, username)
		if output, err := cmd.CombinedOutput(); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update comment: %v, output: %s", err, string(output))})
			return
		}
	}

	// 更新组
	if req.Group != "" {
		cmd := exec.Command("usermod", "-G", req.Group, username)
		if output, err := cmd.CombinedOutput(); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update group: %v, output: %s", err, string(output))})
			return
		}
	}

	c.JSON(200, gin.H{"message": "User updated successfully"})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	username := c.Param("username")

	// 检查用户是否存在
	if _, err := user.Lookup(username); err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 防止删除 root 用户
	if username == "root" {
		c.JSON(403, gin.H{"error": "Cannot delete root user"})
		return
	}

	// 删除用户及其主目录
	cmd := exec.Command("userdel", "-r", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to delete user: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

// GetGroups 获取系统组列表
func GetGroups(c *gin.Context) {
	groups, err := getSystemGroups()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get groups: %v", err)})
		return
	}
	c.JSON(200, gin.H{"groups": groups})
}

// getSystemGroups 从 /etc/group 获取系统组
func getSystemGroups() ([]map[string]string, error) {
	data, err := os.ReadFile("/etc/group")
	if err != nil {
		return nil, fmt.Errorf("failed to read /etc/group: %w", err)
	}

	var groups []map[string]string
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}

		group := map[string]string{
			"name":  fields[0],
			"gid":   fields[2],
			"users": fields[3],
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// GetSSHKeys 获取 SSH 密钥列表
func GetSSHKeys(c *gin.Context) {
	// 获取指定用户的 SSH 密钥，默认当前用户
	targetUser := c.DefaultQuery("user", "")

	if targetUser == "" {
		// 尝试从 JWT 获取用户名（假设中间件已设置）
		if username, exists := c.Get("username"); exists {
			targetUser = username.(string)
		} else {
			// 获取当前进程用户
			if currentUser, err := user.Current(); err == nil {
				targetUser = currentUser.Username
			}
		}
	}

	if targetUser == "" {
		// 如果无法确定用户，返回空列表而不是错误
		c.JSON(200, gin.H{"keys": []SSHKey{}, "user": "", "message": "Unable to determine user"})
		return
	}

	keys, err := getSSHKeys(targetUser)
	if err != nil {
		// 即使出错，也返回空列表而不是500错误
		c.JSON(200, gin.H{"keys": []SSHKey{}, "user": targetUser, "message": "No SSH keys found"})
		return
	}

	c.JSON(200, gin.H{"keys": keys, "user": targetUser})
}

// getSSHKeys 获取用户的 SSH 公钥
func getSSHKeys(username string) ([]SSHKey, error) {
	// 查找用户主目录
	systemUser, err := user.Lookup(username)
	if err != nil {
		// 用户不存在，返回空列表而不是错误
		return []SSHKey{}, nil
	}

	// 读取 authorized_keys 文件
	authorizedKeysPath := filepath.Join(systemUser.HomeDir, ".ssh", "authorized_keys")
	data, err := os.ReadFile(authorizedKeysPath)
	if err != nil {
		// 文件不存在或无权限访问，返回空列表
		return []SSHKey{}, nil
	}

	return parseSSHKeys(data, username)
}

// parseSSHKeys 解析 SSH 公钥
func parseSSHKeys(data []byte, username string) ([]SSHKey, error) {
	var keys []SSHKey
	scanner := bufio.NewScanner(bytes.NewReader(data))

	id := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 SSH 公钥
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		keyType := parts[0]
		keyContent := parts[1]
		comment := ""
		if len(parts) > 2 {
			comment = strings.Join(parts[2:], " ")
		} else {
			comment = fmt.Sprintf("Key %d", id)
		}

		// 计算指纹
		fingerprint, err := generateFingerprint(keyContent)
		if err != nil {
			fingerprint = "unknown"
		}

		key := SSHKey{
			ID:          strconv.Itoa(id),
			Name:        comment,
			Fingerprint: fingerprint,
			Type:        keyType,
			User:        username,
		}

		keys = append(keys, key)
		id++
	}

	return keys, nil
}

// generateFingerprint 生成 SSH 密钥指纹
func generateFingerprint(keyContent string) (string, error) {
	// 解码 base64 密钥
	keyBytes, err := ssh.NewPublicKey([]byte(keyContent))
	if err != nil {
		// 如果解析失败，使用 MD5 作为后备
		hash := md5.Sum([]byte(keyContent))
		return "MD5:" + hex.EncodeToString(hash[:]), nil
	}

	// 生成标准指纹
	return ssh.FingerprintSHA256(keyBytes), nil
}

// AddKey 添加 SSH 密钥
func AddKey(c *gin.Context) {
	var req SSHKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证用户是否存在
	systemUser, err := user.Lookup(req.User)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 验证 SSH 密钥格式
	if !isValidSSHKey(req.Content) {
		c.JSON(400, gin.H{"error": "Invalid SSH key format"})
		return
	}

	// 确保 .ssh 目录存在
	sshDir := filepath.Join(systemUser.HomeDir, ".ssh")
	if err := os.MkdirAll(sshDir, 0700); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create .ssh directory: %v", err)})
		return
	}

	// 读取现有的密钥
	authorizedKeysPath := filepath.Join(sshDir, "authorized_keys")
	var existingKeys []string
	if data, err := os.ReadFile(authorizedKeysPath); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				existingKeys = append(existingKeys, line)
			}
		}
	}

	// 构建新密钥条目
	newKey := req.Content
	if req.Name != "" {
		newKey = fmt.Sprintf("%s %s", req.Content, req.Name)
	}

	// 检查密钥是否已存在
	for _, key := range existingKeys {
		if strings.HasPrefix(key, req.Content) {
			c.JSON(409, gin.H{"error": "SSH key already exists"})
			return
		}
	}

	// 添加新密钥
	existingKeys = append(existingKeys, newKey)

	// 写入文件
	keyData := strings.Join(existingKeys, "\n") + "\n"
	if err := os.WriteFile(authorizedKeysPath, []byte(keyData), 0600); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to write authorized_keys: %v", err)})
		return
	}

	// 设置正确的所有权
	uid, _ := strconv.Atoi(systemUser.Uid)
	gid, _ := strconv.Atoi(systemUser.Gid)
	if err := os.Chown(authorizedKeysPath, uid, gid); err != nil {
		// 继续执行，只是所有权可能不正确
	}
	if err := os.Chown(sshDir, uid, gid); err != nil {
		// 继续执行
	}

	c.JSON(201, gin.H{"message": "SSH key added successfully"})
}

// isValidSSHKey 验证 SSH 密钥格式
func isValidSSHKey(keyContent string) bool {
	// 基本格式检查
	parts := strings.Fields(keyContent)
	if len(parts) < 2 {
		return false
	}

	// 检查密钥类型
	validTypes := []string{"ssh-rsa", "ssh-ed25519", "ecdsa-sha2-nistp256", "ecdsa-sha2-nistp384", "ecdsa-sha2-nistp521", "ssh-dss"}
	for _, t := range validTypes {
		if parts[0] == t {
			return true
		}
	}

	return false
}

// DeleteKey 删除 SSH 密钥
func DeleteKey(c *gin.Context) {
	id := c.Param("id")
	targetUser := c.Query("user")

	if targetUser == "" {
		c.JSON(400, gin.H{"error": "User parameter is required"})
		return
	}

	// 查找用户
	systemUser, err := user.Lookup(targetUser)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 读取现有密钥
	authorizedKeysPath := filepath.Join(systemUser.HomeDir, ".ssh", "authorized_keys")
	data, err := os.ReadFile(authorizedKeysPath)
	if err != nil {
		c.JSON(404, gin.H{"error": "authorized_keys file not found"})
		return
	}

	// 解析并删除指定密钥
	var newKeys []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	currentID := 1
	deleted := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			newKeys = append(newKeys, line)
			continue
		}

		if strconv.Itoa(currentID) == id {
			// 跳过这个密钥（删除）
			deleted = true
		} else {
			newKeys = append(newKeys, line)
		}
		currentID++
	}

	if !deleted {
		c.JSON(404, gin.H{"error": "SSH key not found"})
		return
	}

	// 写入更新后的文件
	keyData := strings.Join(newKeys, "\n") + "\n"
	if err := os.WriteFile(authorizedKeysPath, []byte(keyData), 0600); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to write authorized_keys: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "SSH key deleted successfully"})
}

// GetCurrentUser 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) {
	// 从 JWT 获取用户名（假设中间件已设置）
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	// 获取用户详细信息
	systemUser, err := user.Lookup(username.(string))
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	group := getGroupName(systemUser.Gid)

	userInfo := User{
		Username: systemUser.Username,
		UID:      systemUser.Uid,
		GID:      systemUser.Gid,
		Group:    group,
		Home:     systemUser.HomeDir,
		Shell:    "",
	}

	if shell := getUserShell(systemUser.Username); shell != "" {
		userInfo.Shell = shell
	}

	c.JSON(200, gin.H{"user": userInfo})
}

// ChangeCurrentUserPassword 修改当前用户密码
func ChangeCurrentUserPassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 从 JWT 获取用户名
	username, exists := c.Get("username")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authenticated"})
		return
	}

	// 验证旧密码（这里简化处理，实际应该使用 PAM 或其他认证机制）
	// 在 Linux 上，我们可以尝试验证用户的凭据
	if !verifyPassword(username.(string), req.OldPassword) {
		c.JSON(401, gin.H{"error": "Invalid old password"})
		return
	}

	// 设置新密码
	if err := setUserPassword(username.(string), req.NewPassword); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to change password: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "Password changed successfully"})
}

// verifyPassword 验证用户密码
func verifyPassword(username, password string) bool {
	// 使用 su 命令验证密码
	// 注意：这需要 root 权限或特殊配置
	// 在生产环境中，应该使用 PAM 或其他认证库

	// 简化版本：尝试通过 SSH 验证（这里只是示例）
	// 实际实现应该使用 PAM
	cmd := exec.Command("su", "-c", "true", username)
_stdin, err := cmd.StdinPipe()
	if err != nil {
		return false
	}

	if err := cmd.Start(); err != nil {
		return false
	}

	fmt.Fprintf(_stdin, "%s\n", password)
	_stdin.Close()

	err = cmd.Wait()
	return err == nil
}

// GetUserDiskQuota 获取用户磁盘配额
func GetUserDiskQuota(c *gin.Context) {
	username := c.Param("username")

	// 检查用户是否存在
	systemUser, err := user.Lookup(username)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// 获取用户主目录的磁盘使用情况
	var stat syscall.Statfs_t
	if err := syscall.Statfs(systemUser.HomeDir, &stat); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get disk usage: %v", err)})
		return
	}

	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)
	_ = total - available // calculated but not used in current implementation

	// 获取用户在主目录中的使用量
	cmd := exec.Command("du", "-s", systemUser.HomeDir)
	output, err := cmd.Output()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get user disk usage: %v", err)})
		return
	}

	var userUsage uint64
	if _, err := fmt.Sscanf(string(output), "%d", &userUsage); err != nil {
		userUsage = 0
	}

	c.JSON(200, gin.H{
		"user":  username,
		"home":  systemUser.HomeDir,
		"used":  userUsage * 1024, // du 输出的是 KB
		"total": total,
		"avail": available,
	})
}
