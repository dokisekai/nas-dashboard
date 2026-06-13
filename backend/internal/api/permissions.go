package api

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// ShareFolder 共享文件夹
type ShareFolder struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Browseable  bool     `json:"browseable"`
	ReadOnly    bool     `json:"readOnly"`
	GuestOK     bool     `json:"guestOK"`
	Permissions []string `json:"permissions"`
	ValidUsers  []string `json:"validUsers"`
	ValidGroups []string `json:"validGroups"`
	CreatedAt   string   `json:"createdAt"`
}

// SharePermission 共享权限
type SharePermission struct {
	Type     string `json:"type"`     // "user" or "group"
	Name     string `json:"name"`     // username or groupname
	Permission string `json:"permission"` // "r", "rw", "full"
}

// FilePermission 文件权限
type FilePermission struct {
	Path      string `json:"path"`
	Owner     string `json:"owner"`
	Group     string `json:"group"`
	Mode      string `json:"mode"`      // e.g., "755", "644"
	Symlink   bool   `json:"symlink"`
	IsDir     bool   `json:"isDir"`
	Size      int64  `json:"size"`
	Modified  string `json:"modified"`
}

// ACLRule ACL规则
type ACLRule struct {
	Type    string `json:"type"`    // "user", "group", "other", "mask"
	Name    string `json:"name"`    // username or groupname (if applicable)
	Perms   string `json:"perms"`   // rwx, e.g., "r-x"
	Default bool   `json:"default"` // 是否为默认ACL
}

// PermissionMatrix 权限矩阵
type PermissionMatrix struct {
	Resource    string                       `json:"resource"`
	Permissions map[string]map[string]string `json:"permissions"` // permissions[user][resource] = "r/w/x"
}

// GetShares 获取共享文件夹列表
func GetShares(c *gin.Context) {
	shares, err := getSharesFromConfig()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get shares: %v", err)})
		return
	}

	c.JSON(200, gin.H{"shares": shares})
}

// getSharesFromConfig 从Samba配置获取共享列表
func getSharesFromConfig() ([]ShareFolder, error) {
	smbShares, err := getSMBShares()
	if err != nil {
		return nil, err
	}

	// 转换 SMBShare 到 ShareFolder
	shares := make([]ShareFolder, 0, len(smbShares))
	for _, smbShare := range smbShares {
		share := ShareFolder{
			ID:          generateShareID(smbShare.Name),
			Name:        smbShare.Name,
			Path:        smbShare.Path,
			Description: smbShare.Comment,
			Enabled:     true,
			Browseable:  smbShare.Browseable,
			ReadOnly:    smbShare.ReadOnly,
			GuestOK:     smbShare.GuestOK,
			Permissions: []string{},
			ValidUsers:  []string{},
			ValidGroups: []string{},
			CreatedAt:   getCurrentTimestamp(),
		}
		shares = append(shares, share)
	}

	return shares, nil
}

// CreateShareRequest 创建共享请求
type CreateShareRequest struct {
	Name        string          `json:"name" binding:"required"`
	Path        string          `json:"path" binding:"required"`
	Description string          `json:"description"`
	Browseable  bool            `json:"browseable"`
	ReadOnly    bool            `json:"readOnly"`
	GuestOK     bool            `json:"guestOK"`
	Permissions []SharePermission `json:"permissions"`
}

// CreateShare 创建共享文件夹
func CreateShare(c *gin.Context) {
	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证路径是否存在
	if _, err := os.Stat(req.Path); os.IsNotExist(err) {
		c.JSON(400, gin.H{"error": "Path does not exist"})
		return
	}

	// 检查共享名称是否已存在
	shares, _ := getSharesFromConfig()
	for _, share := range shares {
		if share.Name == req.Name {
			c.JSON(409, gin.H{"error": "Share already exists"})
			return
		}
	}

	// 创建共享配置
	smbConfig := fmt.Sprintf(`[%s]
	path = %s
	browseable = %s
	read only = %s
	guest ok = %s
	valid users = %s
	valid groups = %s
`,
		req.Name,
		req.Path,
		boolToYesNo(req.Browseable),
		boolToYesNo(req.ReadOnly),
		boolToYesNo(req.GuestOK),
		getValidUsersString(req.Permissions),
		getValidGroupsString(req.Permissions),
	)

	// 添加到Samba配置
	if err := addShareToSMBConfig(req.Name, smbConfig); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add share: %v", err)})
		return
	}

	// 设置文件夹权限
	if err := setSharePermissions(req.Path, req.Permissions); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set permissions: %v", err)})
		return
	}

	// 重启Samba服务
	go restartSamba()

	share := ShareFolder{
		ID:          generateShareID(req.Name),
		Name:        req.Name,
		Path:        req.Path,
		Description: req.Description,
		Enabled:     true,
		Browseable:  req.Browseable,
		ReadOnly:    req.ReadOnly,
		GuestOK:     req.GuestOK,
		Permissions: []string{},
		CreatedAt:   getCurrentTimestamp(),
	}

	c.JSON(201, gin.H{"message": "Share created successfully", "share": share})
}

// UpdateShare 更新共享文件夹
func UpdateShare(c *gin.Context) {
	shareName := c.Param("name")

	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查共享是否存在
	shares, _ := getSharesFromConfig()
	var shareExists bool
	for _, share := range shares {
		if share.Name == shareName {
			shareExists = true
			break
		}
	}

	if !shareExists {
		c.JSON(404, gin.H{"error": "Share not found"})
		return
	}

	// 更新Samba配置
	if err := updateShareInSMBConfig(shareName, req); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to update share: %v", err)})
		return
	}

	// 更新文件夹权限
	if err := setSharePermissions(req.Path, req.Permissions); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set permissions: %v", err)})
		return
	}

	// 重启Samba服务
	go restartSamba()

	c.JSON(200, gin.H{"message": "Share updated successfully"})
}

// DeleteShare 删除共享文件夹
func DeleteShare(c *gin.Context) {
	shareName := c.Param("name")

	// 从Samba配置中删除共享
	if err := removeShareFromSMBConfig(shareName); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to delete share: %v", err)})
		return
	}

	// 重启Samba服务
	go restartSamba()

	c.JSON(200, gin.H{"message": "Share deleted successfully"})
}

// GetSharePermissions 获取共享权限
func GetSharePermissions(c *gin.Context) {
	shareName := c.Param("name")

	shares, _ := getSharesFromConfig()
	var sharePath string
	for _, share := range shares {
		if share.Name == shareName {
			sharePath = share.Path
			break
		}
	}

	if sharePath == "" {
		c.JSON(404, gin.H{"error": "Share not found"})
		return
	}

	// 获取文件权限
	fileInfo, err := getFileInfo(sharePath)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get file info: %v", err)})
		return
	}

	// 获取ACL
	aclRules, err := getACLRules(sharePath)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get ACL: %v", err)})
		return
	}

	c.JSON(200, gin.H{
		"fileInfo": fileInfo,
		"acl":     aclRules,
	})
}

// SetSharePermissions 设置共享权限
func SetSharePermissions(c *gin.Context) {
	shareName := c.Param("name")

	var req struct {
		Permissions []SharePermission `json:"permissions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	shares, _ := getSharesFromConfig()
	var sharePath string
	for _, share := range shares {
		if share.Name == shareName {
			sharePath = share.Path
			break
		}
	}

	if sharePath == "" {
		c.JSON(404, gin.H{"error": "Share not found"})
		return
	}

	// 应用权限
	if err := setSharePermissions(sharePath, req.Permissions); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set permissions: %v", err)})
		return
	}

	c.JSON(200, gin.H{"message": "Permissions updated successfully"})
}

// GetFilePermissions 获取文件权限
func GetFilePermissions(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(400, gin.H{"error": "Path is required"})
		return
	}

	fileInfo, err := getFileInfo(path)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get file info: %v", err)})
		return
	}

	c.JSON(200, gin.H{"file": fileInfo})
}

// SetFilePermissions 设置文件权限
func SetFilePermissions(c *gin.Context) {
	var req struct {
		Path string `json:"path" binding:"required"`
		Mode string `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 验证权限模式
	if !isValidPermissionMode(req.Mode) {
		c.JSON(400, gin.H{"error": "Invalid permission mode"})
		return
	}

	// 设置权限
	cmd := exec.Command("chmod", req.Mode, req.Path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set permissions: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Permissions updated successfully"})
}

// GetFileACL 获取文件ACL
func GetFileACL(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(400, gin.H{"error": "Path is required"})
		return
	}

	aclRules, err := getACLRules(path)
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get ACL: %v", err)})
		return
	}

	c.JSON(200, gin.H{"acl": aclRules})
}

// SetFileACL 设置文件ACL
func SetFileACL(c *gin.Context) {
	var req struct {
		Path  string    `json:"path" binding:"required"`
		Rules []ACLRule `json:"rules" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 应用ACL规则
	for _, rule := range req.Rules {
		aclCmd := fmt.Sprintf("%s:%s:%s", rule.Type, rule.Name, rule.Perms)
		if rule.Default {
			aclCmd = fmt.Sprintf("d:%s:%s:%s", rule.Type, rule.Name, rule.Perms)
		}

		cmd := exec.Command("setfacl", "-m", aclCmd, req.Path)
		output, err := cmd.CombinedOutput()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to set ACL: %v, output: %s", err, string(output))})
			return
		}
	}

	c.JSON(200, gin.H{"message": "ACL updated successfully"})
}

// getFileInfo 获取文件信息
func getFileInfo(path string) (FilePermission, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FilePermission{}, err
	}

	// 获取权限
	mode := info.Mode()
	permString := fmt.Sprintf("%o", mode.Perm())

	// 获取所有者
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return FilePermission{}, fmt.Errorf("failed to get stat info")
	}

	uid := strconv.Itoa(int(stat.Uid))
	gid := strconv.Itoa(int(stat.Gid))

	return FilePermission{
		Path:     path,
		Owner:    getUserName(uid),
		Group:    getGroupName(gid),
		Mode:     permString,
		Symlink:  info.Mode()&os.ModeSymlink != 0,
		IsDir:    info.IsDir(),
		Size:     info.Size(),
		Modified: info.ModTime().Format("2006-01-02 15:04:05"),
	}, nil
}

// getACLRules 获取ACL规则
func getACLRules(path string) ([]ACLRule, error) {
	cmd := exec.Command("getfacl", "-c", path)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("getfacl failed: %w", err)
	}

	return parseACLOutput(string(output))
}

// parseACLOutput 解析ACL输出
func parseACLOutput(output string) ([]ACLRule, error) {
	lines := strings.Split(output, "\n")
	var rules []ACLRule

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析ACL行格式: user:username:rwx, group:group:rx, etc.
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				rule := ACLRule{
					Type:  parts[0],
					Name:  parts[1],
					Perms: parts[2],
				}

				if strings.HasPrefix(line, "default:") {
					rule.Default = true
				}

				rules = append(rules, rule)
			}
		}
	}

	return rules, nil
}

// setSharePermissions 设置共享权限
func setSharePermissions(path string, permissions []SharePermission) error {
	for _, perm := range permissions {
		switch perm.Type {
		case "user":
			// 设置用户权限
			return setFileUserPermission(path, perm.Name, perm.Permission)
		case "group":
			// 设置组权限
			return setFileGroupPermission(path, perm.Name, perm.Permission)
		}
	}

	return nil
}

// setFileUserPermission 设置用户文件权限
func setFileUserPermission(path, username, permission string) error {
	// 使用 ACL 设置用户权限
	perms := permissionToPerms(permission)
	cmd := exec.Command("setfacl", "-m", "u:"+username+":"+perms, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set user permission: %v, output: %s", err, string(output))
	}

	return nil
}

// setFileGroupPermission 设置组文件权限
func setFileGroupPermission(path, groupname, permission string) error {
	// 使用 ACL 设置组权限
	perms := permissionToPerms(permission)
	cmd := exec.Command("setfacl", "-m", "g:"+groupname+":"+perms, path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set group permission: %v, output: %s", err, string(output))
	}

	return nil
}

// addShareToSMBConfig 添加共享到Samba配置
func addShareToSMBConfig(name, config string) error {
	configPath := "/etc/samba/smb.conf"

	// 读取现有配置
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	// 添加新共享
	newConfig := string(data) + "\n" + config + "\n"

	// 写回配置文件
	if err := os.WriteFile(configPath, []byte(newConfig), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	return nil
}

// updateShareInSMBConfig 更新Samba配置中的共享
func updateShareInSMBConfig(name string, req CreateShareRequest) error {
	configPath := "/etc/samba/smb.conf"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	// 移除旧共享配置
	lines := strings.Split(string(data), "\n")
	var newLines []string
	inShare := false

	for _, line := range lines {
		if strings.HasPrefix(line, "["+name+"]") {
			inShare = true
			continue
		}

		if inShare && strings.HasPrefix(line, "[") {
			inShare = false
		}

		if !inShare {
			newLines = append(newLines, line)
		}
	}

	// 添加新配置
	smbConfig := fmt.Sprintf(`[%s]
	path = %s
	browseable = %s
	read only = %s
	guest ok = %s
	valid users = %s
	valid groups = %s
`,
		name,
		req.Path,
		boolToYesNo(req.Browseable),
		boolToYesNo(req.ReadOnly),
		boolToYesNo(req.GuestOK),
		getValidUsersString(req.Permissions),
		getValidGroupsString(req.Permissions),
	)

	newLines = append(newLines, smbConfig)

	// 写回配置
	newConfig := strings.Join(newLines, "\n")
	if err := os.WriteFile(configPath, []byte(newConfig), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	return nil
}

// removeShareFromSMBConfig 从Samba配置移除共享
func removeShareFromSMBConfig(name string) error {
	configPath := "/etc/samba/smb.conf"

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read smb.conf: %w", err)
	}

	// 移除共享配置
	lines := strings.Split(string(data), "\n")
	var newLines []string
	inShare := false

	for _, line := range lines {
		if strings.HasPrefix(line, "["+name+"]") {
			inShare = true
			continue
		}

		if inShare && (strings.HasPrefix(line, "[") || line == "") {
			inShare = false
		}

		if !inShare {
			newLines = append(newLines, line)
		}
	}

	// 写回配置
	newConfig := strings.Join(newLines, "\n")
	if err := os.WriteFile(configPath, []byte(newConfig), 0644); err != nil {
		return fmt.Errorf("failed to write smb.conf: %w", err)
	}

	return nil
}

// restartSamba 重启Samba服务
func restartSamba() {
	exec.Command("systemctl", "restart", "smbd").Run()
	exec.Command("systemctl", "restart", "nmbd").Run()
}

// 辅助函数
func boolToYesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func getValidUsersString(permissions []SharePermission) string {
	var users []string
	for _, perm := range permissions {
		if perm.Type == "user" {
			users = append(users, perm.Name)
		}
	}
	return strings.Join(users, " ")
}

func getValidGroupsString(permissions []SharePermission) string {
	var groups []string
	for _, perm := range permissions {
		if perm.Type == "group" {
			groups = append(groups, perm.Name)
		}
	}
	return strings.Join(groups, " ")
}

func generateShareID(name string) string {
	return fmt.Sprintf("share-%s-%d", name, time.Now().Unix())
}

func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func isValidPermissionMode(mode string) bool {
	// 验证权限模式格式（3或4位数字）
	match, _ := regexp.MatchString("^[0-7]{3,4}$", mode)
	return match
}

func permissionToPerms(permission string) string {
	switch permission {
	case "r":
		return "r-x"
	case "w":
		return "rwx"
	case "rw":
		return "rwx"
	case "x":
		return "--x"
	case "rx":
		return "r-x"
	case "wx":
		return "-wx"
	case "rwx":
		return "rwx"
	default:
		return "r-x"
	}
}

func getUserName(uid string) string {
	// 从/etc/passwd查找用户名
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return uid
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) >= 3 && fields[2] == uid {
			return fields[0]
		}
	}

	return uid
}

var _ = fmt.Sprintf // import fmt