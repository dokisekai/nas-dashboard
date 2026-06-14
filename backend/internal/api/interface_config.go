package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

// NetworkInterfaceConfig 网络接口配置
type NetworkInterfaceConfig struct {
	Name         string `json:"name"`
	IPv4Method   string `json:"ipv4Method"`   // dhcp, static, pppoe
	IPv6Method   string `json:"ipv6Method"`   // auto, static, disabled
	IPAddress   string `json:"ipAddress,omitempty"`
	Netmask     string `json:"netmask,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	MACAddress  string `json:"macAddress,omitempty"`
	MTU          int    `json:"mtu,omitempty"`
	DHCP        bool   `json:"dhcp"`
	AutoConnect bool   `json:"autoConnect"`
	Enabled      bool   `json:"enabled"`
}

// PPPoEConfig PPPoE配置
type PPPoEConfig struct {
	Interface  string `json:"interface"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	ServiceName string `json:"serviceName"`
	ACName     string `json:"acName"`
	AutoConnect bool   `json:"autoConnect"`
	Enabled     bool   `json:"enabled"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`     // http, https, socks4, socks5
	Server   string `json:"server"`
	Port     int    `json:"port"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	BypassList []string `json:"bypassList,omitempty"`
}

// InterfaceConfigRequest 接口配置请求
type InterfaceConfigRequest struct {
	IPv4Method string `json:"ipv4Method" binding:"required"` // dhcp, static, pppoe
	IPAddress string `json:"ipAddress,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
	Gateway   string `json:"gateway,omitempty"`
	MTU       int    `json:"mtu,omitempty"`
}

// PPPoEConfigRequest PPPoE配置请求
type PPPoEConfigRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	ServiceName  string `json:"serviceName,omitempty"`
	ACName        string `json:"acName,omitempty"`
	AutoConnect   bool   `json:"autoConnect"`
}

// GetInterfaceConfig 获取接口配置
func GetInterfaceConfig(c *gin.Context) {
	interfaceName := c.Param("interface")

	// 获取接口当前配置
	config, err := getInterfaceConfig(interfaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// SetInterfaceConfig 设置接口配置
func SetInterfaceConfig(c *gin.Context) {
	interfaceName := c.Param("interface")
	var req InterfaceConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 根据不同的方法设置配置
	switch req.IPv4Method {
	case "static":
		if err := setStaticIP(interfaceName, req.IPAddress, req.Netmask, req.Gateway); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "dhcp":
		if err := enableDHCP(interfaceName); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case "pppoe":
		// PPPoE配置需要单独的接口
		c.JSON(http.StatusNotImplemented, gin.H{"error": "PPPoE配置请使用专用接口"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的IP配置方式"})
		return
	}

	// 如果设置了MTU，更新MTU
	if req.MTU > 0 {
		if err := setInterfaceMTU(interfaceName, req.MTU); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// 重新加载网络配置
	if err := reloadNetworkService(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("配置已保存但网络服务重载失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "网络配置已更新"})
}

// ConfigurePPPoE 配置PPPoE连接
func ConfigurePPPoE(c *gin.Context) {
	interfaceName := c.Param("interface")
	var req PPPoEConfigRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 配置PPPoE连接
	if err := configurePPPoEConnection(interfaceName, req.Username, req.Password, req.ServiceName, req.ACName, req.AutoConnect); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PPPoE配置已更新"})
}

// GetPPPoEConfig 获取PPPoE配置
func GetPPPoEConfig(c *gin.Context) {
	interfaceName := c.Param("interface")

	config, err := getPPPoEConfig(interfaceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GetProxyConfig 获取代理配置
func GetProxyConfig(c *gin.Context) {
	config, err := getProxyConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// SetProxyConfig 设置代理配置
func SetProxyConfig(c *gin.Context) {
	var req ProxyConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := setProxyConfiguration(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "代理配置已更新"})
}

// RestartInterface 重启网络接口
func RestartInterface(c *gin.Context) {
	interfaceName := c.Param("interface")

	if err := restartNetworkInterface(interfaceName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "网络接口已重启"})
}

// ==================== 辅助函数 ====================

// getInterfaceConfig 获取接口配置
func getInterfaceConfig(interfaceName string) (*NetworkInterfaceConfig, error) {
	config := &NetworkInterfaceConfig{
		Name:        interfaceName,
		IPv4Method:  "dhcp",
		IPv6Method:  "auto",
		MACAddress: "",
		MTU:         1500,
		DHCP:        true,
		AutoConnect: true,
		Enabled:     true,
	}

	// 获取接口信息
	cmd := exec.Command("ip", "addr", "show", "dev", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("无法获取接口信息: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "inet ") {
			fields := strings.Fields(line[5:])
			if len(fields) >= 2 {
				// 检查是否为静态IP（无broadcast标识）
				if len(fields) < 4 || fields[3] != "brd" {
					config.IPAddress = fields[0]
					// 计算网络掩码
					if len(fields) >= 2 {
						config.Netmask = fields[1]
					}
					config.IPv4Method = "static"
				} else {
					config.IPAddress = fields[0]
					config.Netmask = fields[1]
					config.IPv4Method = "dhcp"
				}
			}
		} else if strings.HasPrefix(line, "link/") {
			fields := strings.Fields(line[2:])
			if len(fields) >= 1 && fields[0] != "NONE" {
				config.MACAddress = fields[0]
			}
		}
	}

	// 获取MTU
	cmd = exec.Command("ip", "link", "show", "dev", interfaceName)
	output, err = cmd.Output()
	if err == nil {
		lines = strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "mtu") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "mtu" && i+1 < len(fields) {
						fmt.Sscanf(fields[i+1], "%d", &config.MTU)
						break
					}
				}
				break
			}
		}
	}

	return config, nil
}

// setStaticIP 设置静态IP
func setStaticIP(interfaceName, ip, netmask, gateway string) error {
	// 首先删除现有IP配置
	cmd := exec.Command("ip", "addr", "flush", "dev", interfaceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("清除现有IP配置失败: %w", err)
	}

	// 添加静态IP
	ipAddr := ip
	if netmask != "" {
		ipAddr = ip + "/" + netmask
	}
	cmd = exec.Command("ip", "addr", "add", ipAddr, "dev", interfaceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置静态IP失败: %w", err)
	}

	// 设置网关
	if gateway != "" {
		cmd = exec.Command("ip", "route", "add", "default", "via", gateway, "dev", interfaceName)
		if err := cmd.Run(); err != nil {
			// 如果网关设置失败，清理已设置的IP
			exec.Command("ip", "addr", "del", ipAddr, "dev", interfaceName).Run()
			return fmt.Errorf("设置网关失败: %w", err)
		}
	}

	// 启用接口
	cmd = exec.Command("ip", "link", "set", "up", "dev", interfaceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用接口失败: %w", err)
	}

	return nil
}

// enableDHCP 启用DHCP
func enableDHCP(interfaceName string) error {
	// 使用nmcli或dhclient配置DHCP
	// 先尝试nmcli
	cmd := exec.Command("nmcli", "connection", "modify", "ipv4.method", "auto", interfaceName)
	if err := cmd.Run(); err != nil {
		// nmcli失败，尝试直接配置DHCP
		// 清除静态IP配置
		cmd = exec.Command("ip", "addr", "flush", "dev", interfaceName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("清除IP配置失败: %w", err)
		}

		// 启用接口并配置DHCP
		cmd = exec.Command("ip", "link", "set", "up", "dev", interfaceName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("启用接口失败: %w", err)
		}
	}

	return nil
}

// setInterfaceMTU 设置接口MTU
func setInterfaceMTU(interfaceName string, mtu int) error {
	cmd := exec.Command("ip", "link", "set", "dev", interfaceName, "mtu", fmt.Sprintf("%d", mtu))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("设置MTU失败: %w", err)
	}
	return nil
}

// reloadNetworkService 重载网络服务
func reloadNetworkService() error {
	// 使用systemctl重启网络服务
	cmd := exec.Command("systemctl", "restart", "NetworkManager")
	if err := cmd.Run(); err != nil {
		// 如果NetworkManager失败，尝试network服务
		cmd = exec.Command("systemctl", "restart", "network")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("重启网络服务失败: %w", err)
		}
	}
	return nil
}

// configurePPPoEConnection 配置PPPoE连接
func configurePPPoEConnection(interfaceName, username, password, serviceName, acName string, autoConnect bool) error {
	// 使用pppd配置PPPoE
	// 这里需要实现pppd配置文件的生成和pppd服务重启
	return fmt.Errorf("PPPoE配置功能需要pppd服务支持")
}

// getPPPoEConfig 获取PPPoE配置
func getPPPoEConfig(interfaceName string) (*PPPoEConfig, error) {
	// 从pppd配置文件读取PPPoE配置
	return nil, fmt.Errorf("PPPoE配置功能需要pppd服务支持")
}

// getProxyConfig 获取代理配置
func getProxyConfig() (*ProxyConfig, error) {
	config := &ProxyConfig{
		Enabled:      false,
		Type:        "http",
		Server:      "",
		Port:        8080,
		BypassList:  []string{"localhost", "127.0.0.1"},
	}

	// 从环境变量或系统配置读取代理设置
	proxy := os.Getenv("HTTP_PROXY")
	if proxy != "" {
		config.Enabled = true
		config.Server = proxy

		httpsProxy := os.Getenv("HTTPS_PROXY")
		if httpsProxy != "" {
			config.Type = "https"
			config.Server = httpsProxy
		}

		noProxy := os.Getenv("NO_PROXY")
		if noProxy != "" {
		config.BypassList = strings.Split(noProxy, ",")
	}
	}

	return config, nil
}

// setProxyConfiguration 设置代理配置
func setProxyConfiguration(config ProxyConfig) error {
	// 设置环境变量和系统代理配置
	// TODO: 实现真实的代理配置功能

	if config.Enabled {
		os.Setenv("HTTP_PROXY", fmt.Sprintf("%s:%d", config.Server, config.Port))
		if config.Type == "https" || config.Type == "socks" {
			os.Setenv("HTTPS_PROXY", fmt.Sprintf("%s:%d", config.Server, config.Port))
		}
		if len(config.BypassList) > 0 {
			os.Setenv("NO_PROXY", strings.Join(config.BypassList, ","))
		}
	} else {
		os.Unsetenv("HTTP_PROXY")
		os.Unsetenv("HTTPS_PROXY")
		os.Unsetenv("NO_PROXY")
	}

	return nil
}

// restartNetworkInterface 重启网络接口
func restartNetworkInterface(interfaceName string) error {
	// 先关闭接口
	cmd := exec.Command("ip", "link", "set", "down", "dev", interfaceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("关闭接口失败: %w", err)
	}

	// 等待一小段时间
	// time.Sleep(time.Second * 2)

	// 重新启用接口
	cmd = exec.Command("ip", "link", "set", "up", "dev", interfaceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用接口失败: %w", err)
	}

	return nil
}