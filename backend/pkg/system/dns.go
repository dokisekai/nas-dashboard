package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DNSConfig DNS配置结构
type DNSConfig struct {
	Method     string   `json:"method"`     // auto, manual
	Primary    string   `json:"primary"`
	Secondary  string   `json:"secondary"`
	Interface  string   `json:"interface"`  // 可选：指定网络接口
}

// ConfigureDNS 配置DNS服务器
func ConfigureDNS(config DNSConfig) error {
	switch config.Method {
	case "auto":
		return configureAutoDNS()
	case "manual":
		return configureManualDNS(config)
	default:
		return fmt.Errorf("不支持的DNS配置方法: %s", config.Method)
	}
}

// configureAutoDNS 配置自动DNS（通过DHCP）
func configureAutoDNS() error {
	// 对于自动DNS，我们需要确保没有手动覆盖resolv.conf
	// 在systemd-resolved环境下，这通常意味着删除任何静态链接

	if isSystemdResolvedAvailable() {
		return resetSystemdResolved()
	} else if isNetworkManagerAvailable() {
		return resetNetworkManagerDNS()
	}

	return nil
}

// configureManualDNS 配置手动DNS服务器
func configureManualDNS(config DNSConfig) error {
	if config.Primary == "" {
		return fmt.Errorf("主要DNS服务器不能为空")
	}

	// 优先使用systemd-resolved
	if isSystemdResolvedAvailable() {
		return configureSystemdResolvedDNS(config)
	}

	// 其次使用NetworkManager
	if isNetworkManagerAvailable() {
		return configureNetworkManagerDNS(config)
	}

	// 最后回退到传统/etc/resolv.conf
	return configureTraditionalDNS(config)
}

// isSystemdResolvedAvailable 检查systemd-resolved是否可用
func isSystemdResolvedAvailable() bool {
	// 检查systemd-resolved服务是否运行
	cmd := exec.Command("systemctl", "is-active", "systemd-resolved")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "active"
}

// isNetworkManagerAvailable 检查NetworkManager是否可用
func isNetworkManagerAvailable() bool {
	// 检查NetworkManager服务是否运行
	cmd := exec.Command("systemctl", "is-active", "NetworkManager")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "active"
}

// configureSystemdResolvedDNS 使用systemd-resolved配置DNS
func configureSystemdResolvedDNS(config DNSConfig) error {
	// 获取主连接名称
	_ = "System wide DNS"

	// 使用resolvectl命令设置DNS（Debian 12推荐方法）
	var dnsServers []string
	dnsServers = append(dnsServers, config.Primary)
	if config.Secondary != "" {
		dnsServers = append(dnsServers, config.Secondary)
	}

	// 设置全局DNS
	args := []string{"dns"}
	if config.Interface != "" {
		args = append(args, config.Interface)
	}
	args = append(args, dnsServers...)

	cmd := exec.Command("resolvectl", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("设置systemd-resolved DNS失败: %s, %v", string(output), err)
	}

	// 重启systemd-resolved服务使配置生效
	if err := restartSystemdResolved(); err != nil {
		return fmt.Errorf("重启systemd-resolved失败: %v", err)
	}

	return nil
}

// resetSystemdResolved 重置systemd-resolved为自动DNS
func resetSystemdResolved() error {
	// 恢复自动DNS
	cmd := exec.Command("resolvectl", "reset-cloudflare", "systemd-resolved")
	_ = cmd.Run() // 忽略错误，可能没有设置过

	// 清除手动设置的DNS
	cmd = exec.Command("resolvectl", "revert", "systemd-resolved")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("重置systemd-resolved DNS失败: %s, %v", string(output), err)
	}

	return nil
}

// restartSystemdResolved 重启systemd-resolved服务
func restartSystemdResolved() error {
	cmd := exec.Command("systemctl", "restart", "systemd-resolved")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("重启systemd-resolved失败: %s, %v", string(output), err)
	}
	return nil
}

// configureNetworkManagerDNS 使用NetworkManager配置DNS
func configureNetworkManagerDNS(config DNSConfig) error {
	// 获取当前活动的连接
	connectionName, err := getActiveConnectionName()
	if err != nil {
		return fmt.Errorf("获取活动连接失败: %v", err)
	}

	if connectionName == "" {
		return fmt.Errorf("没有找到活动的网络连接")
	}

	// 构建DNS服务器列表
	dnsServers := config.Primary
	if config.Secondary != "" {
		dnsServers += " " + config.Secondary
	}

	// 使用nmcli设置DNS
	args := []string{"connection", "modify", connectionName, "ipv4.dns", dnsServers}
	cmd := exec.Command("nmcli", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("设置NetworkManager DNS失败: %s, %v", string(output), err)
	}

	// 重启连接使配置生效
	args = []string{"connection", "up", connectionName}
	cmd = exec.Command("nmcli", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("重启网络连接失败: %s, %v", string(output), err)
	}

	return nil
}

// resetNetworkManagerDNS 重置NetworkManager为自动DNS
func resetNetworkManagerDNS() error {
	connectionName, err := getActiveConnectionName()
	if err != nil || connectionName == "" {
		return nil // 没有活动连接，不需要重置
	}

	// 重置为自动DNS
	args := []string{"connection", "modify", connectionName, "ipv4.dns", "auto"}
	cmd := exec.Command("nmcli", args...)
	_ = cmd.Run() // 忽略错误

	return nil
}

// getActiveConnectionName 获取当前活动的连接名称
func getActiveConnectionName() (string, error) {
	// 使用nmcli获取活动连接
	cmd := exec.Command("nmcli", "-t", "-f", "NAME,TYPE", "connection", "show", "--active")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) >= 2 && fields[1] != "loopback" {
			return fields[0], nil
		}
	}

	return "", fmt.Errorf("没有找到活动连接")
}

// configureTraditionalDNS 配置传统/etc/resolv.conf
func configureTraditionalDNS(config DNSConfig) error {
	// 备份当前resolv.conf
	if err := backupResolvConf(); err != nil {
		return fmt.Errorf("备份resolv.conf失败: %v", err)
	}

	// 构建DNS配置内容
	content := "# Generated by NAS Dashboard\n"
	content += "# Manual DNS configuration\n"
	content += "nameserver " + config.Primary + "\n"
	if config.Secondary != "" {
		content += "nameserver " + config.Secondary + "\n"
	}

	// 添加一些选项
	content += "\n# Options\n"
	content += "options single-request\n"
	content += "options timeout:2 attempts:3\n"

	// 写入/etc/resolv.conf
	if err := os.WriteFile("/etc/resolv.conf", []byte(content), 0644); err != nil {
		return fmt.Errorf("写入resolv.conf失败: %v", err)
	}

	return nil
}

// backupResolvConf 备份/etc/resolv.conf
func backupResolvConf() error {
	// 读取当前文件
	content, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return err
	}

	// 写入备份文件
	backupPath := "/etc/resolv.conf.backup"
	return os.WriteFile(backupPath, content, 0644)
}

// GetCurrentDNSConfig 获取当前DNS配置
func GetCurrentDNSConfig() (*DNSConfig, error) {
	// 优先从systemd-resolved获取
	if isSystemdResolvedAvailable() {
		return getSystemdResolvedDNS()
	}

	// 其次从NetworkManager获取
	if isNetworkManagerAvailable() {
		return getNetworkManagerDNS()
	}

	// 最后从/etc/resolv.conf读取
	return getTraditionalDNS()
}

// getSystemdResolvedDNS 从systemd-resolved获取DNS配置
func getSystemdResolvedDNS() (*DNSConfig, error) {
	cmd := exec.Command("resolvectl", "status")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取systemd-resolved状态失败: %v", err)
	}

	// 解析输出（简化版本）
	lines := strings.Split(string(output), "\n")
	var primary, secondary string

	for _, line := range lines {
		if strings.Contains(line, "DNS Servers:") {
			// 下一行包含DNS服务器
			for i, l := range lines {
				if l == line && i+1 < len(lines) {
					dnsLine := lines[i+1]
					dnsServers := strings.Fields(dnsLine)
					if len(dnsServers) > 0 {
						primary = dnsServers[0]
					}
					if len(dnsServers) > 1 {
						secondary = dnsServers[1]
					}
					break
				}
			}
			break
		}
	}

	return &DNSConfig{
		Method:     "manual",
		Primary:    primary,
		Secondary:  secondary,
	}, nil
}

// getNetworkManagerDNS 从NetworkManager获取DNS配置
func getNetworkManagerDNS() (*DNSConfig, error) {
	connectionName, err := getActiveConnectionName()
	if err != nil {
		return nil, err
	}

	// 获取DNS配置
	cmd := exec.Command("nmcli", "-t", "-f", "ipv4.dns", "connection", "show", connectionName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取NetworkManager DNS配置失败: %v", err)
	}

	dnsLine := strings.TrimSpace(string(output))
	if dnsLine == "" || dnsLine == "--" {
		return &DNSConfig{Method: "auto"}, nil
	}

	dnsServers := strings.Split(dnsLine, ",")
	var primary, secondary string
	if len(dnsServers) > 0 {
		primary = dnsServers[0]
	}
	if len(dnsServers) > 1 {
		secondary = dnsServers[1]
	}

	return &DNSConfig{
		Method:     "manual",
		Primary:    primary,
		Secondary:  secondary,
	}, nil
}

// getTraditionalDNS 从/etc/resolv.conf读取DNS配置
func getTraditionalDNS() (*DNSConfig, error) {
	content, err := os.ReadFile("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("读取resolv.conf失败: %v", err)
	}

	var primary, secondary string
	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "nameserver") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if primary == "" {
					primary = fields[1]
				} else if secondary == "" {
					secondary = fields[1]
				}
			}
		}
	}

	method := "auto"
	if primary != "" {
		method = "manual"
	}

	return &DNSConfig{
		Method:     method,
		Primary:    primary,
		Secondary:  secondary,
	}, nil
}