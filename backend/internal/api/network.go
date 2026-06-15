package api

import (
	"fmt"
	"net/http"
	"os/exec"
	"nas-dashboard/pkg/system"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// WiFiNetwork Wi-Fi网络信息
type WiFiNetwork struct {
	SSID           string  `json:"ssid"`
	BSSID          string  `json:"bssid"`
	Security       string  `json:"security"`
	SignalStrength int     `json:"signalStrength"`
	Channel        int     `json:"channel"`
	Connected      bool    `json:"connected"`
	Connecting     bool    `json:"connecting"`
	UploadSpeed    float64 `json:"uploadSpeed,omitempty"`
	DownloadSpeed  float64 `json:"downloadSpeed,omitempty"`
	IPAddress      string  `json:"ipAddress,omitempty"`
}

// WiFiConnectionRequest Wi-Fi连接请求
type WiFiConnectionRequest struct {
	SSID      string `json:"ssid" binding:"required"`
	Password  string `json:"password"`
	Security  string `json:"security"`
	BSSID     string `json:"bssid"`
	IsHidden  bool   `json:"isHidden"`
}

// DNSConfig DNS配置
type DNSConfig struct {
	Method   string   `json:"method"`   // auto, manual
	Primary  string   `json:"primary"`
	Secondary string  `json:"secondary"`
}

// InterfaceControlRequest 接口控制请求
type InterfaceControlRequest struct {
	Action string `json:"action" binding:"required"` // up, down
}

// GetNetworkInterfaces 获取网络接口列表（过滤虚拟接口）
func GetNetworkInterfaces(c *gin.Context) {
	netInfo, err := system.GetNetworkInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 过滤掉虚拟接口，只显示真实的硬件网卡
	filteredInterfaces := make([]system.Interface, 0)
	for _, iface := range netInfo.Interfaces {
		if isHardwareInterface(iface.Name) {
			filteredInterfaces = append(filteredInterfaces, iface)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"interfaces": filteredInterfaces,
		"timestamp":  netInfo.Timestamp,
	})
}

// GetEthernetInterfaces 获取以太网接口
func GetEthernetInterfaces(c *gin.Context) {
	netInfo, err := system.GetNetworkInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回以太网接口（非wlan开头）
	ethernetInterfaces := make([]system.Interface, 0)
	for _, iface := range netInfo.Interfaces {
		if isHardwareInterface(iface.Name) && !isWiFiInterface(iface.Name) {
			ethernetInterfaces = append(ethernetInterfaces, iface)
		}
	}

	c.JSON(http.StatusOK, ethernetInterfaces)
}

// GetWiFiInterfaces 获取Wi-Fi接口
func GetWiFiInterfaces(c *gin.Context) {
	netInfo, err := system.GetNetworkInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 只返回Wi-Fi接口
	wifiInterfaces := make([]system.Interface, 0)
	for _, iface := range netInfo.Interfaces {
		if isWiFiInterface(iface.Name) {
			wifiInterfaces = append(wifiInterfaces, iface)
		}
	}

	c.JSON(http.StatusOK, wifiInterfaces)
}

// ControlInterface 控制网络接口（启用/禁用）
func ControlInterface(c *gin.Context) {
	interfaceName := c.Param("interface")
	action := c.Param("action") // up 或 down

	if interfaceName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "接口名称不能为空"})
		return
	}

	if action != "up" && action != "down" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的操作，必须是 up 或 down"})
		return
	}

	// 使用ip命令控制接口
	cmd := exec.Command("ip", "link", "set", "dev", interfaceName, action)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("控制接口失败: %s, 输出: %s", err.Error(), string(output)),
		})
		return
	}

	// 等待接口状态变更
	time.Sleep(500 * time.Millisecond)

	// 获取更新后的接口状态
	if updatedIface, err := system.GetInterfaceByName(interfaceName); err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("接口%s成功", map[string]string{"up": "启用", "down": "禁用"}[action]),
			"interface": updatedIface,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("接口%s成功", map[string]string{"up": "启用", "down": "禁用"}[action]),
		})
	}
}

// ScanWiFiNetworks 扫描Wi-Fi网络
func ScanWiFiNetworks(c *gin.Context) {
	// 使用nmcli或iwlist扫描Wi-Fi网络
	networks, err := scanWiFiWithNmcli()
	if err != nil {
		// 尝试使用iwlist
		networks, err = scanWiFiWithIwlist()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("扫描Wi-Fi失败: %v", err)})
			return
		}
	}

	// 获取当前连接的Wi-Fi
	connectedSSID, _ := getCurrentWiFiConnection()

	// 标记已连接的网络
	for i := range networks {
		if networks[i].SSID == connectedSSID {
			networks[i].Connected = true
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"networks": networks,
		"count":    len(networks),
		"timestamp": time.Now().Unix(),
	})
}

// ConnectToWiFi 连接到Wi-Fi网络
func ConnectToWiFi(c *gin.Context) {
	var req WiFiConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证输入
	if req.SSID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "SSID不能为空"})
		return
	}

	// 使用nmcli连接Wi-Fi
	err := connectWiFiWithNmcli(req.SSID, req.Password, req.Security)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("连接Wi-Fi失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Wi-Fi连接成功",
		"ssid":    req.SSID,
	})
}

// DisconnectWiFi 断开Wi-Fi连接
func DisconnectWiFi(c *gin.Context) {
	// 使用nmcli断开Wi-Fi
	cmd := exec.Command("nmcli", "device", "disconnect", "type", "wifi")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("断开Wi-Fi失败: %s, 输出: %s", err.Error(), string(output)),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wi-Fi已断开"})
}

// GetCurrentWiFiConnection 获取当前Wi-Fi连接信息
func GetCurrentWiFiConnection(c *gin.Context) {
	ssid, err := getCurrentWiFiConnection()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未连接Wi-Fi"})
		return
	}

	// 获取Wi-Fi接口信息
	wifiInterfaces, _ := getWiFiInterfaces()
	if len(wifiInterfaces) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到Wi-Fi接口"})
		return
	}

	// 获取接口详细信息
	iface, err := system.GetInterfaceByName(wifiInterfaces[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ssid":      ssid,
		"interface": iface,
		"timestamp": time.Now().Unix(),
	})
}

// GetDNSConfig 获取DNS配置
func GetDNSConfig(c *gin.Context) {
	// 返回默认配置（实际环境中可以从systemd-resolved或NetworkManager获取）
	c.JSON(http.StatusOK, DNSConfig{
		Method:    "auto",
		Primary:   "8.8.8.8",
		Secondary: "8.8.4.4",
	})
}

// SetDNSConfig 设置DNS配置
func SetDNSConfig(c *gin.Context) {
	var req DNSConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用系统DNS配置工具
	if err := system.ConfigureDNS(system.DNSConfig{
		Method:     req.Method,
		Primary:    req.Primary,
		Secondary:  req.Secondary,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("DNS配置失败: %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DNS配置已更新"})
}

// IPConfig IP配置
type IPConfig struct {
	Interface string `json:"interface" binding:"required"`
	Method    string `json:"method" binding:"required"` // dhcp, static
	Address   string `json:"address"`                  // 192.168.1.10/24
	Gateway   string `json:"gateway"`
	DNS       string `json:"dns"`                      // 逗号分隔
}

// UpdateIPConfig 设置接口 IP 配置
func UpdateIPConfig(c *gin.Context) {
	var req IPConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查 nmcli 是否可用
	if _, err := exec.LookPath("nmcli"); err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "NetworkManager (nmcli) is required for static IP configuration"})
		return
	}

	// 1. 获取对应的连接名称 (Connection Name)
	cmdGetConn := exec.Command("nmcli", "-g", "NAME,DEVICE", "con", "show")
	output, _ := cmdGetConn.Output()
	connName := ""
	for _, line := range strings.Split(string(output), "\n") {
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[1] == req.Interface {
			connName = parts[0]
			break
		}
	}

	if connName == "" {
		// 如果没找到连接，尝试直接用接口名
		connName = req.Interface
	}

	// 2. 构建 nmcli 修改命令
	args := []string{"con", "mod", connName}
	if req.Method == "static" {
		if req.Address == "" || req.Gateway == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Address and Gateway are required for static IP"})
			return
		}
		args = append(args, "ipv4.method", "manual", "ipv4.addresses", req.Address, "ipv4.gateway", req.Gateway)
		if req.DNS != "" {
			args = append(args, "ipv4.dns", req.DNS)
		}
	} else {
		args = append(args, "ipv4.method", "auto")
	}

	// 3. 执行修改
	cmdMod := exec.Command("nmcli", args...)
	if output, err := cmdMod.CombinedOutput(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to modify connection: %v, output: %s", err, string(output))})
		return
	}

	// 4. 应用修改 (需要重启连接，可能会导致短暂断开)
	go func() {
		time.Sleep(1 * time.Second)
		exec.Command("nmcli", "con", "up", connName).Run()
	}()

	c.JSON(http.StatusOK, gin.H{"message": "IP configuration updated. Connection will restart."})
}

// 辅助函数

// isHardwareInterface 判断是否为硬件接口
func isHardwareInterface(name string) bool {
	// 过滤掉虚拟接口
	virtualPrefixes := []string{"vir", "docker", "br-", "veth", "tap", "tun"}
	for _, prefix := range virtualPrefixes {
		if strings.HasPrefix(name, prefix) {
			return false
		}
	}
	// 过滤掉本地回环
	return name != "lo"
}

// isWiFiInterface 判断是否为Wi-Fi接口
func isWiFiInterface(name string) bool {
	return strings.HasPrefix(name, "wlan") || strings.HasPrefix(name, "wlp")
}

// scanWiFiWithNmcli 使用nmcli扫描Wi-Fi
func scanWiFiWithNmcli() ([]WiFiNetwork, error) {
	cmd := exec.Command("nmcli", "device", "wifi", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("nmcli扫描失败: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	networks := make([]WiFiNetwork, 0)

	// 跳过标题行
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 7 {
			continue
		}

		// 解析nmcli输出格式
		// 格式: IN-USE  BSSID  SSID  MODE  CHAN  RATE  SIGNAL  BARS  SECURITY
		network := WiFiNetwork{
			BSSID:    fields[1],
			SSID:     fields[2],
			Channel:  parseChannel(fields[4]),
			Security: parseSecurity(fields[8]),
		}

		// 解析信号强度
		if len(fields) > 6 {
			signalStr := strings.TrimSuffix(fields[6], "%")
			network.SignalStrength = parseInt(signalStr)
		}

		// 检查是否正在使用
		if fields[0] == "*" {
			network.Connected = true
		}

		networks = append(networks, network)
	}

	return networks, nil
}

// scanWiFiWithIwlist 使用iwlist扫描Wi-Fi
func scanWiFiWithIwlist() ([]WiFiNetwork, error) {
	// 获取Wi-Fi接口
	wifiInterfaces, _ := getWiFiInterfaces()
	if len(wifiInterfaces) == 0 {
		return nil, fmt.Errorf("未找到Wi-Fi接口")
	}

	cmd := exec.Command("iwlist", wifiInterfaces[0], "scan")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("iwlist扫描失败: %w", err)
	}

	// 解析iwlist输出
	return parseIwlistOutput(string(output)), nil
}

// parseIwlistOutput 解析iwlist输出
func parseIwlistOutput(output string) []WiFiNetwork {
	networks := make([]WiFiNetwork, 0)
	cells := strings.Split(output, "Cell")

	for _, cell := range cells {
		if len(cell) < 10 {
			continue
		}

		network := WiFiNetwork{}

		// 解析SSID
		if ssid := extractField(cell, "ESSID:"); ssid != "" {
			network.SSID = ssid
		}

		// 解析BSSID
		if bssid := extractField(cell, "Address:"); bssid != "" {
			network.BSSID = bssid
		}

		// 解析信号强度
		if quality := extractQuality(cell); quality > 0 {
			network.SignalStrength = quality
		}

		// 解析频道
		if channel := extractChannel(cell); channel > 0 {
			network.Channel = channel
		}

		// 解析加密类型
		if strings.Contains(cell, "Encryption key:on") {
			network.Security = "WPA2" // 简化处理
		} else {
			network.Security = "open"
		}

		if network.SSID != "" {
			networks = append(networks, network)
		}
	}

	return networks
}

// getCurrentWiFiConnection 获取当前Wi-Fi连接
func getCurrentWiFiConnection() (string, error) {
	cmd := exec.Command("iwgetid", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	ssid := strings.TrimSpace(string(output))
	if ssid == "" {
		return "", fmt.Errorf("未连接Wi-Fi")
	}

	return ssid, nil
}

// getWiFiInterfaces 获取Wi-Fi接口列表
func getWiFiInterfaces() ([]string, error) {
	cmd := exec.Command("ls", "/sys/class/net/")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	interfaces := strings.Fields(string(output))
	wifiInterfaces := make([]string, 0)

	for _, iface := range interfaces {
		if isWiFiInterface(iface) {
			wifiInterfaces = append(wifiInterfaces, iface)
		}
	}

	return wifiInterfaces, nil
}

// connectWiFiWithNmcli 使用nmcli连接Wi-Fi
func connectWiFiWithNmcli(ssid, password, security string) error {
	// 构建nmcli连接命令
	args := []string{"device", "wifi", "connect", ssid}

	if password != "" {
		args = append(args, "password", password)
	}

	cmd := exec.Command("nmcli", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("连接失败: %s, 输出: %s", err.Error(), string(output))
	}

	return nil
}

// 辅助解析函数
func extractField(cell, field string) string {
	idx := strings.Index(cell, field)
	if idx == -1 {
		return ""
	}

	start := idx + len(field)
	end := strings.Index(cell[start:], "\n")
	if end == -1 {
		return strings.TrimSpace(cell[start:])
	}

	return strings.TrimSpace(cell[start : start+end])
}

func extractQuality(cell string) int {
	idx := strings.Index(cell, "Quality=")
	if idx == -1 {
		return 0
	}

	start := idx + 8 // len("Quality=")
	end := strings.Index(cell[start:], "/")
	if end == -1 {
		return 0
	}

	qualityStr := cell[start : start+end]
	quality := parseInt(qualityStr)
	return quality
}

func extractChannel(cell string) int {
	idx := strings.Index(cell, "Channel:")
	if idx == -1 {
		return 0
	}

	start := idx + 8 // len("Channel:")
	end := strings.Index(cell[start:], "\n")
	if end == -1 {
		return 0
	}

	channelStr := strings.TrimSpace(cell[start : start+end])
	return parseInt(channelStr)
}

func parseChannel(channelStr string) int {
	return parseInt(channelStr)
}

func parseSecurity(securityStr string) string {
	if strings.Contains(securityStr, "WPA3") {
		return "WPA3"
	} else if strings.Contains(securityStr, "WPA2") {
		return "WPA2"
	} else if strings.Contains(securityStr, "WPA") {
		return "WPA"
	} else if strings.Contains(securityStr, "WEP") {
		return "WEP"
	}
	return "open"
}

func parseInt(s string) int {
	var result int
	fmt.Sscanf(s, "%d", &result)
	return result
}
