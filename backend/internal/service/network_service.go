package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NetworkService 网络管理服务
type NetworkService struct {
	wsUpgrader  websocket.Upgrader
	wsClients   map[string]*websocket.Conn
	wsMutex     sync.RWMutex
	monitorChan chan NetworkMetric
}

// NetworkMetric 网络指标
type NetworkMetric struct {
	Timestamp    time.Time `json:"timestamp"`
	Interface    string    `json:"interface"`
	BytesSent    uint64    `json:"bytesSent"`
	BytesRecv    uint64    `json:"bytesRecv"`
	PacketsSent  uint64    `json:"packetsSent"`
	PacketsRecv  uint64    `json:"packetsRecv"`
	ErrorsIn     uint64    `json:"errorsIn"`
	ErrorsOut    uint64    `json:"errorsOut"`
	DropsIn      uint64    `json:"dropsIn"`
	DropsOut     uint64    `json:"dropsOut"`
	UploadSpeed  float64   `json:"uploadSpeed"`
	DownloadSpeed float64  `json:"downloadSpeed"`
}

// NetworkInterfaceConfig 网络接口配置
type NetworkInterfaceConfig struct {
	Name         string   `json:"name"`
	Enabled      bool     `json:"enabled"`
	AutoConnect  bool     `json:"autoConnect"`
	IPv4         IPv4Config `json:"ipv4"`
	IPv6         IPv6Config `json:"ipv6"`
	MAC          string   `json:"mac"`
	MTU          int      `json:"mtu"`
	DNS          []string `json:"dns"`
	Gateway      string   `json:"gateway"`
	RouteMetric  int      `json:"routeMetric"`
}

// IPv4Config IPv4配置
type IPv4Config struct {
	Method string `json:"method"` // auto, manual, disabled
	Address string `json:"address"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
}

// IPv6Config IPv6配置
type IPv6Config struct {
	Method string `json:"method"` // auto, manual, disabled
	Address string `json:"address"`
	Prefix string `json:"prefix"`
	Gateway string `json:"gateway"`
}

// ProxyConfig 代理配置
type ProxyConfig struct {
	Enabled     bool     `json:"enabled"`
	HTTPProxy   string   `json:"httpProxy"`
	HTTPSProxy  string   `json:"httpsProxy"`
	NoProxy     []string `json:"noProxy"`
	AuthEnabled bool     `json:"authEnabled"`
	Username    string   `json:"username"`
	Password    string   `json:"password"`
}

// FirewallRule 防火墙规则
type FirewallRule struct {
	ID          int      `json:"id"`
	Action      string   `json:"action"`      // allow, deny, reject, limit
	Direction   string   `json:"direction"`   // in, out, both
	Interface   string   `json:"interface"`
	Protocol    string   `json:"protocol"`    // tcp, udp, both, any
	SourceIP    string   `json:"sourceIp"`
	DestIP      string   `json:"destIp"`
	Port        string   `json:"port"`
	Enabled     bool     `json:"enabled"`
	Description string   `json:"description"`
	Logging     bool     `json:"logging"`
}

// NetworkService 网络服务实现
type NetworkServiceImpl struct {
	configPath string
	dnsPath    string
}

// NewNetworkService 创建网络服务
func NewNetworkService() *NetworkService {
	service := &NetworkService{
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 生产环境需要检查来源
			},
		},
		wsClients:   make(map[string]*websocket.Conn),
		monitorChan: make(chan NetworkMetric, 100),
	}

	// 启动网络监控
	go service.monitorNetwork()
	go service.broadcastMetrics()

	return service
}

// GetNetworkInterfaces 获取网络接口列表
func (s *NetworkService) GetNetworkInterfaces() ([]NetworkInterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %v", err)
	}

	var result []NetworkInterfaceInfo
	for _, iface := range interfaces {
		// 过滤掉回环接口和down状态的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		info := NetworkInterfaceInfo{
			Name:       iface.Name,
			MAC:        iface.HardwareAddr.String(),
			MTU:        iface.MTU,
			IsUp:       iface.Flags&net.FlagUp != 0,
			IsRunning:  iface.Flags&net.FlagRunning != 0,
		}

		// 获取地址信息
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipv4 := ipnet.IP.To4(); ipv4 != nil {
					info.IPv4 = ipv4.String()
					info.Netmask = net.IP(ipnet.Mask).String()
				} else if ipv6 := ipnet.IP.To16(); ipv6 != nil {
					info.IPv6 = ipv6.String()
				}
			}
		}

		// 获取接口统计信息
		stats := s.getInterfaceStats(iface.Name)
		info.BytesSent = stats.BytesSent
		info.BytesRecv = stats.BytesRecv
		info.PacketsSent = stats.PacketsSent
		info.PacketsRecv = stats.PacketsRecv

		result = append(result, info)
	}

	return result, nil
}

// NetworkInterfaceInfo 网络接口信息
type NetworkInterfaceInfo struct {
	Name         string `json:"name"`
	IPv4         string `json:"ipv4"`
	IPv6         string `json:"ipv6"`
	Netmask      string `json:"netmask"`
	MAC          string `json:"mac"`
	MTU          int    `json:"mtu"`
	IsUp         bool   `json:"isUp"`
	IsRunning    bool   `json:"isRunning"`
	BytesSent    uint64 `json:"bytesSent"`
	BytesRecv    uint64 `json:"bytesRecv"`
	PacketsSent  uint64 `json:"packetsSent"`
	PacketsRecv  uint64 `json:"packetsRecv"`
}

// ConfigureInterface 配置网络接口
func (s *NetworkService) ConfigureInterface(config NetworkInterfaceConfig) error {
	// 验证配置
	if err := s.validateInterfaceConfig(config); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	// 先关闭接口
	if err := s.setInterfaceDown(config.Name); err != nil {
		return fmt.Errorf("关闭接口失败: %v", err)
	}

	// 配置 IPv4
	if err := s.configureIPv4(config.Name, config.IPv4); err != nil {
		return fmt.Errorf("配置 IPv4 失败: %v", err)
	}

	// 配置 IPv6
	if err := s.configureIPv6(config.Name, config.IPv6); err != nil {
		return fmt.Errorf("配置 IPv6 失败: %v", err)
	}

	// 配置 MTU
	if config.MTU > 0 {
		if err := s.setMTU(config.Name, config.MTU); err != nil {
			return fmt.Errorf("设置 MTU 失败: %v", err)
		}
	}

	// 配置路由
	if config.IPv4.Gateway != "" {
		if err := s.setDefaultRoute(config.Name, config.IPv4.Gateway); err != nil {
			return fmt.Errorf("设置网关失败: %v", err)
		}
	}

	// 配置 DNS
	if len(config.DNS) > 0 {
		if err := s.setDNS(config.DNS); err != nil {
			return fmt.Errorf("设置 DNS 失败: %v", err)
		}
	}

	// 启动接口
	if err := s.setInterfaceUp(config.Name); err != nil {
		return fmt.Errorf("启动接口失败: %v", err)
	}

	return nil
}

// validateInterfaceConfig 验证接口配置
func (s *NetworkService) validateInterfaceConfig(config NetworkInterfaceConfig) error {
	if config.Name == "" {
		return fmt.Errorf("接口名称不能为空")
	}

	// 验证 IPv4 配置
	if config.IPv4.Method == "manual" {
		if config.IPv4.Address == "" {
			return fmt.Errorf("手动配置模式下 IP 地址不能为空")
		}
		if config.IPv4.Netmask == "" {
			return fmt.Errorf("手动配置模式下子网掩码不能为空")
		}

		// 验证 IP 地址格式
		if net.ParseIP(config.IPv4.Address) == nil {
			return fmt.Errorf("无效的 IP 地址: %s", config.IPv4.Address)
		}
	}

	// 验证 DNS 配置
	for _, dns := range config.DNS {
		if net.ParseIP(dns) == nil {
			return fmt.Errorf("无效的 DNS 服务器: %s", dns)
		}
	}

	return nil
}

// configureIPv4 配置 IPv4
func (s *NetworkService) configureIPv4(interfaceName string, config IPv4Config) error {
	switch config.Method {
	case "manual":
		// 使用 ip 命令配置静态 IP
		cmd := exec.Command("ip", "addr", "add", fmt.Sprintf("%s/%s", config.Address, config.Netmask), "dev", interfaceName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("配置静态 IP 失败: %v", err)
		}

	case "dhcp":
		// 使用 dhclient 获取 DHCP
		cmd := exec.Command("dhclient", interfaceName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("启动 DHCP 客户端失败: %v", err)
		}

	case "disabled":
		// 不配置 IPv4
		return nil
	}

	return nil
}

// configureIPv6 配置 IPv6
func (s *NetworkService) configureIPv6(interfaceName string, config IPv6Config) error {
	switch config.Method {
	case "manual":
		// 使用 ip 命令配置 IPv6
		cmd := exec.Command("ip", "-6", "addr", "add", fmt.Sprintf("%s/%s", config.Address, config.Prefix), "dev", interfaceName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("配置 IPv6 失败: %v", err)
		}

	case "auto":
		// 启用 IPv6 自动配置
		cmd := exec.Command("sysctl", "-w", fmt.Sprintf("net.ipv6.conf.%s.autoconf=1", interfaceName))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("启用 IPv6 自动配置失败: %v", err)
		}

	case "disabled":
		// 禁用 IPv6
		cmd := exec.Command("sysctl", "-w", fmt.Sprintf("net.ipv6.conf.%s.disable_ipv6=1", interfaceName))
		return cmd.Run()
	}

	return nil
}

// setInterfaceUp 启动网络接口
func (s *NetworkService) setInterfaceUp(interfaceName string) error {
	cmd := exec.Command("ip", "link", "set", "up", interfaceName)
	return cmd.Run()
}

// setInterfaceDown 关闭网络接口
func (s *NetworkService) setInterfaceDown(interfaceName string) error {
	cmd := exec.Command("ip", "link", "set", "down", interfaceName)
	return cmd.Run()
}

// setMTU 设置 MTU
func (s *NetworkService) setMTU(interfaceName string, mtu int) error {
	cmd := exec.Command("ip", "link", "set", "dev", interfaceName, "mtu", strconv.Itoa(mtu))
	return cmd.Run()
}

// setDefaultRoute 设置默认路由
func (s *NetworkService) setDefaultRoute(interfaceName, gateway string) error {
	// 先删除现有的默认路由
	exec.Command("ip", "route", "del", "default").Run()

	// 添加新的默认路由
	cmd := exec.Command("ip", "route", "add", "default", "via", gateway, "dev", interfaceName)
	return cmd.Run()
}

// setDNS 设置 DNS
func (s *NetworkService) setDNS(servers []string) error {
	// 更新 /etc/resolv.conf
	content := "# Generated by NAS Dashboard\n"
	content += "# Nameservers\n"
	for _, server := range servers {
		content += fmt.Sprintf("nameserver %s\n", server)
	}

	// 添加一些选项
	content += "\n# Options\n"
	content += "options timeout:2 attempts:3\n"
	content += "options single-request\n"

	return os.WriteFile("/etc/resolv.conf", []byte(content), 0644)
}

// ConfigureProxy 配置代理
func (s *NetworkService) ConfigureProxy(config ProxyConfig) error {
	if !config.Enabled {
		// 清除代理环境变量
		s.clearProxyEnv()
		return nil
	}

	// 设置代理环境变量
	proxyEnvs := map[string]string{
		"http_proxy":  config.HTTPProxy,
		"https_proxy": config.HTTPSProxy,
		"no_proxy":     strings.Join(config.NoProxy, ","),
	}

	if config.AuthEnabled {
		proxyURL := fmt.Sprintf("%s:%s@", config.Username, config.Password)
		proxyEnvs["http_proxy"] = strings.Replace(config.HTTPProxy, "://", "://"+proxyURL, 1)
		proxyEnvs["https_proxy"] = strings.Replace(config.HTTPSProxy, "://", "://"+proxyURL, 1)
	}

	// 写入环境配置文件
	envFile := "/etc/environment"
	content := "# Proxy configuration\n"
	for key, value := range proxyEnvs {
		if value != "" {
			content += fmt.Sprintf("%s=\"%s\"\n", key, value)
		}
	}

	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入代理配置失败: %v", err)
	}

	// 为当前进程设置环境变量
	for key, value := range proxyEnvs {
		if value != "" {
			os.Setenv(key, value)
			os.Setenv(strings.ToUpper(key), value)
		}
	}

	return nil
}

// clearProxyEnv 清除代理环境变量
func (s *NetworkService) clearProxyEnv() {
	proxyVars := []string{"http_proxy", "https_proxy", "no_proxy", "HTTP_PROXY", "HTTPS_PROXY", "NO_PROXY"}

	envFile := "/etc/environment"
	var content string
	if data, err := os.ReadFile(envFile); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			isProxyVar := false
			for _, proxyVar := range proxyVars {
				if strings.HasPrefix(line, proxyVar+"=") || strings.HasPrefix(line, strings.ToUpper(proxyVar)+"=") {
					isProxyVar = true
					break
				}
			}
			if !isProxyVar && line != "" {
				content += line + "\n"
			}
		}
	}

	os.WriteFile(envFile, []byte(content), 0644)

	// 清除当前进程环境变量
	for _, proxyVar := range proxyVars {
		os.Unsetenv(proxyVar)
		os.Unsetenv(strings.ToUpper(proxyVar))
	}
}

// ApplyFirewallRules 应用防火墙规则
func (s *NetworkService) ApplyFirewallRules(rules []FirewallRule) error {
	// 确保 SSH 端口始终开启，防止锁定自己
	exec.Command("ufw", "allow", "22/tcp").Run()

	// 应用面板端口
	exec.Command("ufw", "allow", "8080/tcp").Run()
	exec.Command("ufw", "allow", "5173/tcp").Run()

	// 清除现有规则（保留默认规则）
	exec.Command("ufw", "--force", "reset").Run()

	// 应用新规则
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}

		if err := s.applyFirewallRule(rule); err != nil {
			log.Printf("应用防火墙规则失败: %v", err)
			continue
		}
	}

	// 启用防火墙
	cmd := exec.Command("ufw", "--force", "enable")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("启用防火墙失败: %s, %v", string(output), err)
	}

	return nil
}

// applyFirewallRule 应用单个防火墙规则
func (s *NetworkService) applyFirewallRule(rule FirewallRule) error {
	// 构建ufw命令参数
	args := []string{rule.Action}

	if rule.SourceIP != "" && rule.SourceIP != "any" {
		args = append(args, "from", rule.SourceIP)
	}

	if rule.Port != "" {
		portProto := rule.Port
		if rule.Protocol != "both" && rule.Protocol != "" {
			portProto += "/" + rule.Protocol
		}
		args = append(args, "to", "any", "port", portProto)
	}

	if rule.Logging {
		args = append(args, "log")
	}

	cmd := exec.Command("ufw", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("执行规则失败: %s, %v", string(output), err)
	}

	return nil
}

// GetFirewallStatus 获取防火墙状态
func (s *NetworkService) GetFirewallStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// 获取 ufw 状态
	cmd := exec.Command("ufw", "status", "numbered")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return status, fmt.Errorf("获取防火墙状态失败: %v", err)
	}

	// 解析输出
	lines := strings.Split(string(output), "\n")
	var rules []string
	var isActive bool

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "Status: active") {
			isActive = true
		}
		if strings.HasPrefix(line, "[") {
			rules = append(rules, line)
		}
	}

	status["active"] = isActive
	status["rules"] = rules

	return status, nil
}

// monitorNetwork 监控网络流量
func (s *NetworkService) monitorNetwork() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// 用于计算速率的上次数据
	lastStats := make(map[string]*InterfaceStats)

	for {
		select {
		case <-ticker.C:
			interfaces, _ := s.GetNetworkInterfaces()
			currentTime := time.Now()

			for _, iface := range interfaces {
				// 获取当前统计信息
				currentStats := s.getInterfaceStats(iface.Name)

				// 计算速率（如果存在上次数据）
				if lastStat, exists := lastStats[iface.Name]; exists {
					timeDiff := currentTime.Sub(lastStat.Timestamp).Seconds()
					if timeDiff > 0 {
						uploadSpeed := float64(currentStats.BytesSent-lastStat.BytesSent) / timeDiff
						downloadSpeed := float64(currentStats.BytesRecv-lastStat.BytesRecv) / timeDiff

						metric := NetworkMetric{
							Timestamp:     currentTime,
							Interface:     iface.Name,
							BytesSent:     currentStats.BytesSent,
							BytesRecv:     currentStats.BytesRecv,
							PacketsSent:   currentStats.PacketsSent,
							PacketsRecv:   currentStats.PacketsRecv,
							ErrorsIn:      currentStats.ErrorsIn,
							ErrorsOut:     currentStats.ErrorsOut,
							DropsIn:       currentStats.DropsIn,
							DropsOut:      currentStats.DropsOut,
							UploadSpeed:   uploadSpeed,
							DownloadSpeed: downloadSpeed,
						}

						// 发送到监控通道
						select {
						case s.monitorChan <- metric:
						default:
							// 通道满，丢弃旧数据
						}
					}
				}

				// 保存当前统计信息
				lastStats[iface.Name] = &InterfaceStats{
					BytesSent:    currentStats.BytesSent,
					BytesRecv:    currentStats.BytesRecv,
					PacketsSent:  currentStats.PacketsSent,
					PacketsRecv:  currentStats.PacketsRecv,
					ErrorsIn:     currentStats.ErrorsIn,
					ErrorsOut:    currentStats.ErrorsOut,
					DropsIn:      currentStats.DropsIn,
					DropsOut:     currentStats.DropsOut,
					Timestamp:    currentTime,
				}
			}
		}
	}
}

// InterfaceStats 接口统计信息
type InterfaceStats struct {
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
	ErrorsIn    uint64
	ErrorsOut   uint64
	DropsIn     uint64
	DropsOut    uint64
	Timestamp   time.Time
}

// getInterfaceStats 获取接口统计信息
func (s *NetworkService) getInterfaceStats(interfaceName string) InterfaceStats {
	stats := InterfaceStats{}

	// 从 /proc/net/dev 读取统计信息
	data, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return stats
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 17 {
			continue
		}

		// 查找匹配的接口
		if strings.HasSuffix(fields[0], ":") && strings.TrimSuffix(fields[0], ":") == interfaceName {
			// 解析统计数据
			stats.BytesRecv, _ = strconv.ParseUint(fields[1], 10, 64)
			stats.PacketsRecv, _ = strconv.ParseUint(fields[2], 10, 64)
			stats.ErrorsIn, _ = strconv.ParseUint(fields[3], 10, 64)
			stats.DropsIn, _ = strconv.ParseUint(fields[4], 10, 64)
			stats.BytesSent, _ = strconv.ParseUint(fields[9], 10, 64)
			stats.PacketsSent, _ = strconv.ParseUint(fields[10], 10, 64)
			stats.ErrorsOut, _ = strconv.ParseUint(fields[11], 10, 64)
			stats.DropsOut, _ = strconv.ParseUint(fields[12], 10, 64)
			break
		}
	}

	return stats
}

// broadcastMetrics 广播网络指标
func (s *NetworkService) broadcastMetrics() {
	for metric := range s.monitorChan {
		s.wsMutex.RLock()
		clients := make([]*websocket.Conn, 0, len(s.wsClients))
		for _, conn := range s.wsClients {
			clients = append(clients, conn)
		}
		s.wsMutex.RUnlock()

		message := map[string]interface{}{
			"type":    "network_metric",
			"metric": metric,
		}

		messageJSON, _ := json.Marshal(message)

		for _, conn := range clients {
			if err := conn.WriteMessage(websocket.TextMessage, messageJSON); err != nil {
				// 移除失败的连接
				s.wsMutex.Lock()
				for id, c := range s.wsClients {
					if c == conn {
						delete(s.wsClients, id)
						break
					}
				}
				s.wsMutex.Unlock()
			}
		}
	}
}

// RegisterNetworkWebSocket 注册网络 WebSocket 客户端
func (s *NetworkService) RegisterNetworkWebSocket(clientID string, conn *websocket.Conn) {
	s.wsMutex.Lock()
	defer s.wsMutex.Unlock()

	s.wsClients[clientID] = conn
}

// UnregisterNetworkWebSocket 注销网络 WebSocket 客户端
func (s *NetworkService) UnregisterNetworkWebSocket(clientID string) {
	s.wsMutex.Lock()
	defer s.wsMutex.Unlock()

	delete(s.wsClients, clientID)
}