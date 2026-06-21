package service

import (
	"fmt"
	"log"
	"nas-dashboard/internal/models"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 临时的模型定义，用于编译通过
type Group struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
}

type FilePermissions struct {
	Path    string `json:"path"`
	UID     uint   `json:"uid"`
	GID     uint   `json:"gid"`
	Owner   string `json:"owner"`
	Group   string `json:"group"`
	Mode    string `json:"mode"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
	ACL     []ACL  `json:"acl,omitempty"`
}

type ACL struct {
	Type    string `json:"type"`
	Entity  string `json:"entity,omitempty"` // 实体（用户或组）
	User    string `json:"user,omitempty"`
	Group   string `json:"group,omitempty"`
	Perm    string `json:"perm"`
	Perms   string `json:"perms,omitempty"` // 备用字段
	Default bool   `json:"default"`
}

type AccessLevel string

const (
	AccessLevelNone   AccessLevel = "none"
	AccessLevelRead   AccessLevel = "read"
	AccessLevelWrite  AccessLevel = "write"
	AccessLevelFull   AccessLevel = "full"
)

// Permission 权限模型
type Permission struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	UserID    uint      `json:"userId"`
	GroupID   uint      `json:"groupId"`
	Path      string    `json:"path"`
	Mode      string    `json:"mode"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Name     string  `json:"name"`
	Status   string  `json:"status"`
}

// UserGroup 用户组
type UserGroup struct {
	UserID    uint      `json:"userId"`
	GroupID   uint      `json:"groupId"`
	CreatedAt time.Time `json:"createdAt"`
}

// PermissionService 权限管理服务
type PermissionService struct {
	db *gorm.DB
}

// NewPermissionService 创建权限管理服务
func NewPermissionService(db *gorm.DB) *PermissionService {
	return &PermissionService{db: db}
}

// UserManagement 用户管理接口
type UserManagement interface {
	// CreateUser 创建用户
	CreateUser(user *models.User) error
	// GetUser 获取用户信息
	GetUser(userID uint) (*models.User, error)
	// ListUsers 列出用户
	ListUsers() ([]models.User, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *models.User) error
	// DeleteUser 删除用户
	DeleteUser(userID uint) error
	// SetUserPassword 设置用户密码
	SetUserPassword(userID uint, password string) error
	// DisableUser 禁用用户
	DisableUser(userID uint) error
	// EnableUser 启用用户
	EnableUser(userID uint) error
}

// GroupManagement 用户组管理接口
type GroupManagement interface {
	// CreateGroup 创建用户组
	CreateGroup(group interface{}) error
	// GetGroup 获取用户组信息
	GetGroup(groupID uint) (interface{}, error)
	// ListGroups 列出用户组
	ListGroups() ([]interface{}, error)
	// UpdateGroup 更新用户组
	UpdateGroup(group interface{}) error
	// DeleteGroup 删除用户组
	DeleteGroup(groupID uint) error
	// AddUserToGroup 添加用户到用户组
	AddUserToGroup(userID, groupID uint) error
	// RemoveUserFromGroup 从用户组移除用户
	RemoveUserFromGroup(userID, groupID uint) error
	// GetGroupUsers 获取用户组成员
	GetGroupUsers(groupID uint) ([]models.User, error)
}

// FilePermissionManagement 文件权限管理接口
type FilePermissionManagement interface {
	// GetFilePermissions 获取文件权限
	GetFilePermissions(filePath string) (interface{}, error)
	// SetFilePermissions 设置文件权限
	SetFilePermissions(filePath string, permissions interface{}) error
	// SetFileOwner 设置文件所有者
	SetFileOwner(filePath, username, groupname string) error
	// GetAccessControlList 获取访问控制列表
	GetAccessControlList(filePath string) ([]interface{}, error)
	// SetAccessControlList 设置访问控制列表
	SetAccessControlList(filePath string, acls interface{}) error
	// CheckFileAccess 检查文件访问权限
	CheckFileAccess(filePath string, userID uint) (interface{}, error)
}

// CreateUser 创建用户
func (s *PermissionService) CreateUser(user *models.User) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return fmt.Errorf("用户名已存在: %s", user.Username)
	}

	// 如果提供了密码，加密密码
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("密码加密失败: %v", err)
		}
		user.Password = string(hashedPassword)
	}

	// 设置默认值
	if user.Role == "" {
		user.Role = "user" // 默认为普通用户
	}
	user.Status = "active"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// 在数据库中创建用户记录
	if err := s.db.Create(user).Error; err != nil {
		return fmt.Errorf("创建用户记录失败: %v", err)
	}

	// 在系统中创建实际用户
	if err := s.createSystemUser(user); err != nil {
		// 如果系统用户创建失败，回滚数据库
		s.db.Delete(user)
		return fmt.Errorf("创建系统用户失败: %v", err)
	}

	// 创建用户主目录
	if err := s.createUserHomeDirectory(user.Username); err != nil {
		log.Printf("创建用户主目录失败: %v", err)
		// 不影响用户创建，只记录错误
	}

	return nil
}

// createSystemUser 在系统中创建用户
func (s *PermissionService) createSystemUser(user *models.User) error {
	// 检查用户是否已存在
	cmd := exec.Command("id", "-u", user.Username)
	if err := cmd.Run(); err == nil {
		// 用户已存在，更新用户信息
		return s.updateSystemUser(user)
	}

	// 创建新用户
	// 使用 useradd 命令
	args := []string{
		"-m",                    // 创建主目录
		"-s", "/bin/bash",       // 设置默认 shell
		"-c", user.NickName,      // 设置注释（全名）
		user.Username,
	}

	if user.Password != "" {
		// 设置密码
		args = append([]string{"-p", user.Password}, args...)
	}

	cmd = exec.Command("useradd", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行 useradd 失败: %s, %v", string(output), err)
	}

	// 将用户添加到相应的组
	if user.Role == "admin" {
		// 添加到 sudo 组
		exec.Command("usermod", "-aG", "sudo", user.Username).Run()
		// 添加到特定的 NAS 管理组
		exec.Command("usermod", "-aG", "nas-admin", user.Username).Run()
	} else {
		// 普通用户添加到基本用户组
		exec.Command("usermod", "-aG", "nas-users", user.Username).Run()
	}

	return nil
}

// updateSystemUser 更新系统用户
func (s *PermissionService) updateSystemUser(user *models.User) error {
	// 更新用户注释
	if user.NickName != "" {
		cmd := exec.Command("usermod", "-c", user.NickName, user.Username)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("更新用户注释失败: %v", err)
		}
	}

	// 如果需要修改密码
	if user.Password != "" {
		cmd := exec.Command("chpasswd")
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return fmt.Errorf("创建密码管道失败: %v", err)
		}

		go func() {
			defer stdin.Close()
			fmt.Fprintf(stdin, "%s:%s", user.Username, user.Password)
		}()

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("修改密码失败: %v", err)
		}
	}

	return nil
}

// createUserHomeDirectory 创建用户主目录
func (s *PermissionService) createUserHomeDirectory(username string) error {
	// 创建 NAS 用户目录
	nasUserDir := fmt.Sprintf("/home/nas/%s", username)

	// 创建目录
	if err := os.MkdirAll(nasUserDir, 0755); err != nil {
		return fmt.Errorf("创建用户目录失败: %v", err)
	}

	// 设置目录所有者
	cmd := exec.Command("chown", "-R", fmt.Sprintf("%s:%s", username, username), nasUserDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置目录所有者失败: %v", err)
	}

	// 创建用户主目录的子目录
	subdirs := []string{"Documents", "Downloads", "Pictures", "Music", "Videos"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(nasUserDir, subdir)
		os.MkdirAll(subdirPath, 0755)
	}

	return nil
}

// GetUser 获取用户信息
func (s *PermissionService) GetUser(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户不存在: %v", err)
	}

	return &user, nil
}

// ListUsers 列出用户
func (s *PermissionService) ListUsers() ([]models.User, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	var userInfos []models.User
	for _, user := range users {
		// 过滤系统用户
		if s.isSystemUser(user.Username) {
			continue
		}

		userInfo, err := s.GetUser(user.ID)
		if err != nil {
			log.Printf("获取用户信息失败: %v", err)
			continue
		}

		userInfos = append(userInfos, *userInfo)
	}

	return userInfos, nil
}

// isSystemUser 判断是否为系统用户
func (s *PermissionService) isSystemUser(username string) bool {
	systemUsers := []string{
		"lp", "mail", "news", "uucp", "proxy", "www-data", "backup",
		"list", "irc", "gnats", "nobody", "_apt", "systemd-network",
		"systemd-resolve", "messagebus", "sshd", "nodered",
	}

	for _, systemUser := range systemUsers {
		if username == systemUser {
			return true
		}
	}

	// 检查 UID 是否小于 1000
	cmd := exec.Command("id", "-u", username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	uid := strings.TrimSpace(string(output))
	uidInt, err := strconv.Atoi(uid)
	if err != nil {
		return false
	}

	return uidInt < 1000
}

// UpdateUser 更新用户信息
func (s *PermissionService) UpdateUser(user *models.User) error {
	var existingUser models.User
	if err := s.db.Where("id = ?", user.ID).First(&existingUser).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 更新数据库记录
	updates := map[string]interface{}{
		"nick_name":  user.NickName,
		"email":      user.Email,
		"role":       user.Role,
		"status":     user.Status,
		"updated_at": time.Now(),
	}

	if err := s.db.Model(&existingUser).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新用户记录失败: %v", err)
	}

	// 更新系统用户
	if err := s.updateSystemUser(user); err != nil {
		log.Printf("更新系统用户失败: %v", err)
	}

	return nil
}

// DeleteUser 删除用户
func (s *PermissionService) DeleteUser(userID uint) error {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 检查是否为最后一个管理员
	if user.Role == "admin" {
		var adminCount int64
		s.db.Model(&models.User{}).Where("role = ? AND status = ?", "admin", "active").Count(&adminCount)
		if adminCount <= 1 {
			return fmt.Errorf("无法删除最后一个管理员用户")
		}
	}

	// 备份用户数据
	if err := s.backupUserData(user.Username); err != nil {
		log.Printf("备份用户数据失败: %v", err)
	}

	// 从系统中删除用户
	cmd := exec.Command("userdel", "-r", user.Username)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("删除系统用户失败: %s, %v", string(output), err)
	}

	// 删除数据库记录
	if err := s.db.Delete(&user).Error; err != nil {
		return fmt.Errorf("删除用户记录失败: %v", err)
	}

	return nil
}

// backupUserData 备份用户数据
func (s *PermissionService) backupUserData(username string) error {
	backupDir := fmt.Sprintf("/backup/users/%s", username)
	backupFile := fmt.Sprintf("%s_%s.tar.gz", backupDir, time.Now().Format("20060102_150405"))

	// 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	// 备份用户主目录
	userHome := fmt.Sprintf("/home/%s", username)
	if _, err := os.Stat(userHome); err == nil {
		cmd := exec.Command("tar", "-czf", backupFile, "-C", "/home", username)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("备份用户数据失败: %v", err)
		}
	}

	return nil
}

// SetUserPassword 设置用户密码
func (s *PermissionService) SetUserPassword(userID uint, password string) error {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	// 验证密码强度
	if err := s.validatePasswordStrength(password); err != nil {
		return err
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %v", err)
	}

	// 更新数据库密码
	user.Password = string(hashedPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("更新密码失败: %v", err)
	}

	// 更新系统密码
	user.Password = password // 使用明文密码更新系统用户
	if err := s.updateSystemUser(&user); err != nil {
		log.Printf("更新系统密码失败: %v", err)
	}

	return nil
}

// validatePasswordStrength 验证密码强度
func (s *PermissionService) validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于8个字符")
	}

	// 检查密码复杂度
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case strings.ContainsAny(string(char), "!@#$%^&*()_+-=[]{}|;:,.<>?"):
			hasSpecial = true
		}
	}

	complexity := 0
	if hasUpper {
		complexity++
	}
	if hasLower {
		complexity++
	}
	if hasDigit {
		complexity++
	}
	if hasSpecial {
		complexity++
	}

	if complexity < 3 {
		return fmt.Errorf("密码复杂度不足，必须包含大写字母、小写字母、数字和特殊字符中的至少三种")
	}

	return nil
}

// CreateGroup 创建用户组
func (s *PermissionService) CreateGroup(group *Group) error {
	// 检查组名是否已存在
	var existingGroup Group
	if err := s.db.Where("name = ?", group.Name).First(&existingGroup).Error; err == nil {
		return fmt.Errorf("用户组名已存在: %s", group.Name)
	}

	// 在系统中创建组
	cmd := exec.Command("groupadd", group.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("创建系统组失败: %s, %v", string(output), err)
	}

	// 在数据库中创建组记录
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	if err := s.db.Create(group).Error; err != nil {
		// 如果数据库操作失败，回滚系统组
		exec.Command("groupdel", group.Name).Run()
		return fmt.Errorf("创建组记录失败: %v", err)
	}

	return nil
}

// AddUserToGroup 添加用户到用户组
func (s *PermissionService) AddUserToGroup(userID, groupID uint) error {
	var user models.User
	var group Group

	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	if err := s.db.Where("id = ?", groupID).First(&group).Error; err != nil {
		return fmt.Errorf("用户组不存在: %v", err)
	}

	// 检查是否已是成员
	var userGroup UserGroup
	if err := s.db.Where("user_id = ? AND group_id = ?", userID, groupID).First(&userGroup).Error; err == nil {
		return fmt.Errorf("用户已是该组成员")
	}

	// 在系统中添加用户到组
	cmd := exec.Command("usermod", "-aG", group.Name, user.Username)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("添加用户到系统组失败: %s, %v", string(output), err)
	}

	// 在数据库中创建关联记录
	userGroup = UserGroup{
		UserID:    userID,
		GroupID:   groupID,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&userGroup).Error; err != nil {
		// 如果数据库操作失败，回滚系统操作
		exec.Command("gpasswd", "-d", user.Username, group.Name).Run()
		return fmt.Errorf("创建用户组关联失败: %v", err)
	}

	return nil
}

// RemoveUserFromGroup 从用户组移除用户
func (s *PermissionService) RemoveUserFromGroup(userID, groupID uint) error {
	var user models.User
	var group Group

	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("用户不存在: %v", err)
	}

	if err := s.db.Where("id = ?", groupID).First(&group).Error; err != nil {
		return fmt.Errorf("用户组不存在: %v", err)
	}

	// 从系统中移除用户
	cmd := exec.Command("gpasswd", "-d", user.Username, group.Name)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("从系统组移除用户失败: %s, %v", string(output), err)
	}

	// 从数据库中删除关联记录
	if err := s.db.Where("user_id = ? AND group_id = ?", userID, groupID).Delete(&UserGroup{}).Error; err != nil {
		return fmt.Errorf("删除用户组关联失败: %v", err)
	}

	return nil
}

// GetFilePermissions 获取文件权限
func (s *PermissionService) GetFilePermissions(filePath string) (*FilePermissions, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("文件不存在: %v", err)
	}

	// 获取文件信息
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 获取文件权限
	mode := info.Mode().Perm()

	// 获取文件所有者
	uid := info.Sys().(*syscall.Stat_t).Uid
	gid := info.Sys().(*syscall.Stat_t).Gid

	// 获取用户名和组名
	username := s.getUsernameByID(uint(uid))
	groupname := s.getGroupnameByID(uint(gid))

	// 获取 ACL
	acls, _ := s.getFileACL(filePath)

	permissions := &FilePermissions{
		Path:      filePath,
		Mode:      fmt.Sprintf("%04o", mode),
		Owner:     username,
		Group:     groupname,
		UID:       uint(uid),
		GID:       uint(gid),
		Size:      info.Size(),
		ModTime:   info.ModTime().Format("2006-01-02 15:04:05"),
		IsDir:     info.IsDir(),
		ACL:       acls,
	}

	return permissions, nil
}

// SetFilePermissions 设置文件权限
func (s *PermissionService) SetFilePermissions(filePath string, permissions *FilePermissions) error {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("文件不存在: %v", err)
	}

	// 设置权限模式
	mode, err := strconv.ParseUint(permissions.Mode, 8, 32)
	if err != nil {
		return fmt.Errorf("无效的权限模式: %v", err)
	}

	if err := os.Chmod(filePath, os.FileMode(mode)); err != nil {
		return fmt.Errorf("设置文件权限失败: %v", err)
	}

	// 设置所有者
	if permissions.Owner != "" || permissions.Group != "" {
		uid := permissions.UID
		gid := permissions.GID

		if permissions.Owner != "" {
			uid = uint(s.getUserIDByName(permissions.Owner))
		}

		if permissions.Group != "" {
			gid = uint(s.getGroupIDByName(permissions.Group))
		}

		if err := os.Chown(filePath, int(uid), int(gid)); err != nil {
			return fmt.Errorf("设置文件所有者失败: %v", err)
		}
	}

	// 设置 ACL
	if len(permissions.ACL) > 0 {
		if err := s.setFileACL(filePath, permissions.ACL); err != nil {
			return fmt.Errorf("设置 ACL 失败: %v", err)
		}
	}

	return nil
}

// getFileACL 获取文件 ACL
func (s *PermissionService) getFileACL(filePath string) ([]ACL, error) {
	cmd := exec.Command("getfacl", "-p", filePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("获取 ACL 失败: %v", err)
	}

	// 解析 ACL 输出
	return s.parseACL(string(output))
}

// setFileACL 设置文件 ACL
func (s *PermissionService) setFileACL(filePath string, acls []ACL) error {
	// 首先清除现有 ACL
	exec.Command("setfacl", "-b", filePath).Run()

	// 设置新的 ACL
	for _, acl := range acls {
		var aclSpec string

		switch acl.Perms {
		case "read":
			aclSpec = fmt.Sprintf("%s:r", acl.Entity)
		case "write":
			aclSpec = fmt.Sprintf("%s:rw", acl.Entity)
		case "execute":
			aclSpec = fmt.Sprintf("%s:rx", acl.Entity)
		case "full":
			aclSpec = fmt.Sprintf("%s:rwx", acl.Entity)
		}

		cmd := exec.Command("setfacl", "-m", aclSpec, filePath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("设置 ACL 规则失败: %v", err)
		}
	}

	// 设置默认 ACL（如果是目录）
	if info, err := os.Stat(filePath); err == nil && info.IsDir() {
		for _, acl := range acls {
			var aclSpec string

			switch acl.Perms {
			case "read":
				aclSpec = fmt.Sprintf("d:%s:r", acl.Entity)
			case "write":
				aclSpec = fmt.Sprintf("d:%s:rw", acl.Entity)
			case "execute":
				aclSpec = fmt.Sprintf("d:%s:rx", acl.Entity)
			case "full":
				aclSpec = fmt.Sprintf("d:%s:rwx", acl.Entity)
			}

			cmd := exec.Command("setfacl", "-d", "-m", aclSpec, filePath)
			if err := cmd.Run(); err != nil {
				log.Printf("设置默认 ACL 失败: %v", err)
			}
		}
	}

	return nil
}

// parseACL 解析 ACL 输出
func (s *PermissionService) parseACL(output string) ([]ACL, error) {
	var acls []ACL

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// 解析 ACL 行
		// 格式: user:username:perms, group:groupname:perms, other:perms, etc.
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				acl := ACL{
					Type:   parts[0], // user, group, other, mask
					Entity: parts[1],
					Perms:  parts[2],
				}
				acls = append(acls, acl)
			}
		}
	}

	return acls, nil
}

// CheckFileAccess 检查文件访问权限
func (s *PermissionService) CheckFileAccess(filePath string, userID uint) (AccessLevel, error) {
	// 简化版本：返回基本权限
	return AccessLevelRead, nil
}

// 辅助方法
func (s *PermissionService) getUsernameByID(uid uint) string {
	return "unknown"
}

func (s *PermissionService) getGroupnameByID(gid uint) string {
	return "unknown"
}

func (s *PermissionService) getUserIDByName(username string) int {
	return 1000
}

func (s *PermissionService) getGroupIDByName(groupname string) int {
	return 1000
}
