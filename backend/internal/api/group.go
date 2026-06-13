package api

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/gin-gonic/gin"
)

// Group 用户组
type Group struct {
	Name     string   `json:"name"`
	GID      string   `json:"gid"`
	Members  []string `json:"members"`
	Password string   `json:"password"`
}

// GroupMemberRequest 组成员请求
type GroupMemberRequest struct {
	Users []string `json:"users" binding:"required"`
}

// GetGroups 获取用户组列表
func GetGroups(c *gin.Context) {
	groups, err := getSystemGroups()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get groups: %v", err)})
		return
	}
	c.JSON(200, gin.H{"groups": groups})
}

// getSystemGroups 从 /etc/group 获取系统用户组
func getSystemGroups() ([]Group, error) {
	data, err := os.ReadFile("/etc/group")
	if err != nil {
		return nil, fmt.Errorf("failed to read /etc/group: %w", err)
	}

	return parseGroupFile(data)
}

// parseGroupFile 解析 /etc/group 文件
func parseGroupFile(data []byte) ([]Group, error) {
	var groups []Group
	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		fields := strings.Split(line, ":")
		if len(fields) < 4 {
			continue
		}

		groupName := fields[0]
		gid := fields[2]
		members := strings.Split(fields[3], ",")

		// 过滤空成员
		var validMembers []string
		for _, member := range members {
			if member != "" {
				validMembers = append(validMembers, member)
			}
		}

		group := Group{
			Name:     groupName,
			GID:      gid,
			Members:  validMembers,
			Password: fields[1],
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// GetGroup 获取单个用户组信息
func GetGroup(c *gin.Context) {
	groupName := c.Param("name")

	groups, err := getSystemGroups()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get groups: %v", err)})
		return
	}

	for _, group := range groups {
		if group.Name == groupName {
			c.JSON(200, gin.H{"group": group})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Group not found"})
}

// CreateGroupRequest 创建组请求
type CreateGroupRequest struct {
	Name    string   `json:"name" binding:"required"`
	GID     string   `json:"gid"`
	Members []string `json:"members"`
}

// CreateGroup 创建用户组
func CreateGroup(c *gin.Context) {
	var req CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查组是否已存在
	groups, _ := getSystemGroups()
	for _, group := range groups {
		if group.Name == req.Name {
			c.JSON(409, gin.H{"error": "Group already exists"})
			return
		}
	}

	// 构建 groupadd 命令
	args := []string{req.Name}

	// 设置 GID
	if req.GID != "" {
		// 验证 GID 是否已被使用
		for _, group := range groups {
			if group.GID == req.GID {
				c.JSON(400, gin.H{"error": "GID already in use"})
				return
			}
		}
		args = append([]string{"-g", req.GID}, args...)
	}

	cmd := exec.Command("groupadd", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to create group: %v, output: %s", err, string(output))})
		return
	}

	// 如果有成员，添加到组
	if len(req.Members) > 0 {
		for _, member := range req.Members {
			usermodCmd := exec.Command("usermod", "-aG", req.Name, member)
			if output, err := usermodCmd.CombinedOutput(); err != nil {
				// 记录警告但不中断
				fmt.Printf("Warning: failed to add user %s to group %s: %v, output: %s\n",
					member, req.Name, err, string(output))
			}
		}
	}

	c.JSON(201, gin.H{"message": "Group created successfully", "group": req.Name})
}

// UpdateGroupRequest 更新组请求
type UpdateGroupRequest struct {
	Members []string `json:"members"`
}

// UpdateGroup 更新用户组成员
func UpdateGroup(c *gin.Context) {
	groupName := c.Param("name")

	var req UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 获取当前组成员
	currentGroups, err := getSystemGroups()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get groups: %v", err)})
		return
	}

	var currentGroup *Group
	for i, group := range currentGroups {
		if group.Name == groupName {
			currentGroup = &currentGroups[i]
			break
		}
	}

	if currentGroup == nil {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	// 计算需要添加和移除的成员
	currentMembers := make(map[string]bool)
	for _, member := range currentGroup.Members {
		currentMembers[member] = true
	}

	newMembers := make(map[string]bool)
	for _, member := range req.Members {
		newMembers[member] = true
	}

	// 添加新成员
	for _, member := range req.Members {
		if !currentMembers[member] {
			cmd := exec.Command("usermod", "-aG", groupName, member)
			if output, err := cmd.CombinedOutput(); err != nil {
				c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add user %s to group: %v, output: %s",
					member, err, string(output))})
				return
			}
		}
	}

	// 移除旧成员
	for _, member := range currentGroup.Members {
		if !newMembers[member] && member != "" {
			cmd := exec.Command("gpasswd", "-d", groupName, member)
			if output, err := cmd.CombinedOutput(); err != nil {
				// 记录警告但不中断
				fmt.Printf("Warning: failed to remove user %s from group %s: %v, output: %s\n",
					member, groupName, err, string(output))
			}
		}
	}

	c.JSON(200, gin.H{"message": "Group updated successfully"})
}

// DeleteGroup 删除用户组
func DeleteGroup(c *gin.Context) {
	groupName := c.Param("name")

	// 检查组是否存在
	groups, _ := getSystemGroups()
	var groupExists bool
	for _, group := range groups {
		if group.Name == groupName {
			groupExists = true
			break
		}
	}

	if !groupExists {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	// 检查是否有成员
	for _, group := range groups {
		if group.Name == groupName && len(group.Members) > 0 {
			c.JSON(400, gin.H{"error": "Cannot delete group with members"})
			return
		}
	}

	cmd := exec.Command("groupdel", groupName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to delete group: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Group deleted successfully"})
}

// GetGroupMembers 获取组成员
func GetGroupMembers(c *gin.Context) {
	groupName := c.Param("name")

	groups, err := getSystemGroups()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get groups: %v", err)})
		return
	}

	for _, group := range groups {
		if group.Name == groupName {
			memberDetails := make([]User, 0)
			for _, member := range group.Members {
				if user, err := getUserByName(member); err == nil {
					memberDetails = append(memberDetails, user)
				}
			}

			c.JSON(200, gin.H{
				"members": memberDetails,
				"count":    len(memberDetails),
			})
			return
		}
	}

	c.JSON(404, gin.H{"error": "Group not found"})
}

// AddGroupMembers 添加组成员
func AddGroupMembers(c *gin.Context) {
	groupName := c.Param("name")

	var req GroupMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查组是否存在
	groups, _ := getSystemGroups()
	var groupExists bool
	for _, group := range groups {
		if group.Name == groupName {
			groupExists = true
			break
		}
	}

	if !groupExists {
		c.JSON(404, gin.H{"error": "Group not found"})
		return
	}

	// 添加用户到组
	for _, username := range req.Users {
		cmd := exec.Command("usermod", "-aG", groupName, username)
		if output, err := cmd.CombinedOutput(); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to add user %s to group: %v, output: %s",
				username, err, string(output))})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Members added successfully"})
}

// RemoveGroupMember 移除组成员
func RemoveGroupMember(c *gin.Context) {
	groupName := c.Param("name")
	username := c.Param("user")

	cmd := exec.Command("gpasswd", "-d", groupName, username)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to remove user from group: %v, output: %s",
			err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Member removed successfully"})
}

// getUserByName 根据用户名获取用户信息
func getUserByName(username string) (User, error) {
	systemUser, err := user.Lookup(username)
	if err != nil {
		return User{}, err
	}

	group := getGroupName(systemUser.Gid)

	return User{
		Username: systemUser.Username,
		UID:      systemUser.Uid,
		GID:      systemUser.Gid,
		Group:    group,
		Home:     systemUser.HomeDir,
		Shell:    getUserShell(systemUser.Username),
	}, nil
}