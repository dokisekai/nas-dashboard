package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	initMutex     sync.Mutex
	initChecked   bool
	initStatus    bool
)

// InitStatus 系统初始化状态
type InitStatus struct {
	Initialized bool   `json:"initialized"`
	Version     string `json:"version,omitempty"`
}

// NetworkConfig 网络配置
type NetworkConfig struct {
	Type       string `json:"type"` // dhcp, static
	StaticIP   string `json:"staticIp"`
	SubnetMask string `json:"subnetMask"`
	Gateway    string `json:"gateway"`
	DNS        string `json:"dns"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	DataPath   string `json:"dataPath" binding:"required"`
	BackupPath string `json:"backupPath"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	Hostname         string `json:"hostname"`
	Timezone         string `json:"timezone"`
	EnableFirewall   bool   `json:"enableFirewall"`
	EnableAutoUpdate bool   `json:"enableAutoUpdate"`
	EnableTelemetry  bool   `json:"enableTelemetry"`
}

// GetSystemInfoExtended 获取扩展的系统信息
func GetSystemInfoExtended(c *gin.Context) {
	hostname, _ := os.Hostname()

	// 获取操作系统信息
	uname := &syscall.Utsname{}
	syscall.Uname(uname)

	c.JSON(http.StatusOK, gin.H{
		"hostname": hostname,
		"os":       "Linux",
		"arch":      "x86_64",
		"version":   "1.0.0",
	})
}

// ResetSystem 重置系统（危险操作，仅用于管理）
func ResetSystem(c *gin.Context) {
	// 检查用户权限
	// 这里需要实现权限检查

	initMutex.Lock()
	defer initMutex.Unlock()

	// 删除初始化标记文件
	initFile := "/var/lib/nas-dashboard/.initialized"
	if err := os.Remove(initFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to reset system: " + err.Error(),
		})
		return
	}

	// 重置状态
	initStatus = false
	initChecked = false

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "System reset successfully. Please restart initialization.",
	})
}

// createAdminUser 创建管理员用户
func createAdminUser(username, password string) error {
	// 1. 创建系统用户
	cmd := exec.Command("useradd", "-m", "-s", "/bin/bash", username)
	if output, err := cmd.CombinedOutput(); err != nil {
		// 用户可能已存在，忽略错误
		if len(output) > 0 && !containsUserExists(output) {
			return err
		}
	}

	// 2. 设置用户密码
	cmd = exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(username + ":" + password)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	// 3. 将用户添加到sudo组
	cmd = exec.Command("usermod", "-aG", "sudo", username)
	if _, err := cmd.CombinedOutput(); err != nil {
		return err
	}

	return nil
}

// containsUserExists 检查输出是否包含用户已存在的错误
func containsUserExists(output []byte) bool {
	return strings.Contains(string(output), "already exists")
}

// configureStaticIP 配置静态IP
func configureStaticIP(network NetworkConfig) error {
	// 这里需要根据具体的Linux发行版和网络管理器来实现
	// 简化版本：仅记录日志
	// 在实际实现中，需要修改网络配置文件或使用NetworkManager/dbus

	// 示例：对于Netplan配置文件
	netplanConfig := fmt.Sprintf(`
network:
  version: 2
  ethernets:
    eth0:
      dhcp4: no
      addresses:
        - %s/%s
      gateway4: %s
      nameservers:
        addresses: [%s]
`, network.StaticIP, network.SubnetMask, network.Gateway, network.DNS)

	// 写入netplan配置文件
	configPath := "/etc/netplan/99-nas-dashboard.yaml"
	if err := os.WriteFile(configPath, []byte(netplanConfig), 0644); err != nil {
		return err
	}

	// 应用网络配置
	cmd := exec.Command("netplan", "apply")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("netplan apply failed: %s, output: %s", err.Error(), string(output))
	}

	return nil
}

// createStorageDirectories 创建存储目录
func createStorageDirectories(storage StorageConfig) error {
	// 创建数据存储目录
	if err := os.MkdirAll(storage.DataPath, 0755); err != nil {
		return err
	}

	// 创建备份存储目录
	if storage.BackupPath != "" {
		if err := os.MkdirAll(storage.BackupPath, 0755); err != nil {
			return err
		}
	}

	// 设置目录权限
	cmd := exec.Command("chown", "-R", "www-data:www-data", storage.DataPath)
	if _, err := cmd.CombinedOutput(); err != nil {
		// 如果www-data用户不存在，使用当前用户
		return nil
	}

	return nil
}

// setSystemConfig 设置系统配置
func setSystemConfig(config SystemConfig) error {
	// 设置主机名
	if config.Hostname != "" {
		cmd := exec.Command("hostnamectl", "set-hostname", config.Hostname)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to set hostname: %s, output: %s", err.Error(), string(output))
		}
	}

	// 设置时区（简化版本）
	if config.Timezone != "" {
		cmd := exec.Command("timedatectl", "set-timezone", config.Timezone)
		if output, err := cmd.CombinedOutput(); err != nil {
			// 时区设置失败不应阻止初始化
			fmt.Printf("Warning: failed to set timezone: %s\n", string(output))
		}
	}

	// 如果需要启用防火墙
	if config.EnableFirewall {
		cmd := exec.Command("ufw", "enable")
		if output, err := cmd.CombinedOutput(); err != nil {
			// 防火墙设置失败不应阻止初始化
			fmt.Printf("Warning: failed to enable firewall: %s\n", string(output))
		}
	}

	return nil
}

// createInitMarker 创建初始化标记文件
func createInitMarker() error {
	initDir := "/var/lib/nas-dashboard"
	if err := os.MkdirAll(initDir, 0755); err != nil {
		return err
	}

	initFile := filepath.Join(initDir, ".initialized")
	file, err := os.Create(initFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入初始化时间
	_, err = file.WriteString(fmt.Sprintf("initialized_at: %d\n", time.Now().Unix()))
	return err
}