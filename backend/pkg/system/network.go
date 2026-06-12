package system

import (
	"fmt"
	"net"
	"strings"
	"time"

	psutilnet "github.com/shirou/gopsutil/v3/net"
)

// NetworkInfo 网络信息
type NetworkInfo struct {
	Interfaces []Interface `json:"interfaces"`
	Timestamp  int64       `json:"timestamp"`
}

// Interface 网络接口信息
type Interface struct {
	Name        string   `json:"name"`         // 接口名称
	HardwareAddr string  `json:"hardwareAddr"`  // MAC 地址
	Up          bool     `json:"up"`           // 是否启用
	Addresses   []string `json:"addresses"`    // IP 地址列表 (IPv4/IPv6)
	MTU         int      `json:"mtu"`          // MTU 值
	Flags       []string `json:"flags"`        // 接口标志

	// 流量统计
	BytesSent       uint64  `json:"bytesSent"`       // 累计发送字节数
	BytesRecv       uint64  `json:"bytesRecv"`       // 累计接收字节数
	PacketsSent     uint64  `json:"packetsSent"`     // 累计发送包数
	PacketsRecv     uint64  `json:"packetsRecv"`     // 累计接收包数
	Errin           uint64  `json:"errin"`          // 接收错误
	Errout          uint64  `json:"errout"`         // 发送错误
	Dropin          uint64  `json:"dropin"`         // 接收丢包
	Dropout         uint64  `json:"dropout"`        // 发送丢包

	// 实时速度 (bytes/s)
	SentSpeed      float64 `json:"sentSpeed"`       // 发送速度
	RecvSpeed      float64 `json:"recvSpeed"`       // 接收速度
}

// NetworkStats 网络统计摘要
type NetworkStats struct {
	TotalBytesSent uint64  `json:"totalBytesSent"`
	TotalBytesRecv uint64  `json:"totalBytesRecv"`
	TotalSentSpeed float64 `json:"totalSentSpeed"`
	TotalRecvSpeed float64 `json:"totalRecvSpeed"`
	Interfaces     int     `json:"interfaces"`
}

// GetNetworkInfo 获取网络信息
func GetNetworkInfo() (*NetworkInfo, error) {
	interfaces, err := psutilnet.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %w", err)
	}

	result := make([]Interface, 0, len(interfaces))

	// 获取 IO 统计
	counters, err := psutilnet.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters: %w", err)
	}

	// 构建计数器映射
	counterMap := make(map[string]psutilnet.IOCountersStat)
	for _, c := range counters {
		counterMap[c.Name] = c
	}

	// 获取上次状态用于计算速度
	state := GetStatsState()
	lastStats, lastTS := state.GetLastNetIO()
	now := time.Now()
	duration := now.Sub(lastTS).Seconds()

	// 当前状态（用于下次计算）
	currentStats := make(map[string]NetIOHistory)

	for _, iface := range interfaces {
		// 跳过本地回环和没有地址的接口
		if shouldSkipInterface(iface) {
			continue
		}

		// 提取 IP 地址
		addresses := extractAddresses(iface)

		// 获取对应的计数器
		counter, hasCounter := counterMap[iface.Name]

		netIface := Interface{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr,
			Up:           isInterfaceUp(iface),
			Addresses:    addresses,
			MTU:          iface.MTU,
			Flags:        iface.Flags,
		}

		if hasCounter {
			netIface.BytesSent = counter.BytesSent
			netIface.BytesRecv = counter.BytesRecv
			netIface.PacketsSent = counter.PacketsSent
			netIface.PacketsRecv = counter.PacketsRecv
			netIface.Errin = counter.Errin
			netIface.Errout = counter.Errout
			netIface.Dropin = counter.Dropin
			netIface.Dropout = counter.Dropout

			// 保存当前统计
			currentStats[iface.Name] = NetIOHistory{
				BytesSent:   counter.BytesSent,
				BytesRecv:   counter.BytesRecv,
				PacketsSent: counter.PacketsSent,
				PacketsRecv: counter.PacketsRecv,
			}

			// 计算速度
			if lastStat, ok := lastStats[iface.Name]; ok && duration > 0 {
				sentSpeed, recvSpeed := CalculateNetSpeed(currentStats[iface.Name], lastStat, duration)
				netIface.SentSpeed = sentSpeed
				netIface.RecvSpeed = recvSpeed
			}
		}

		result = append(result, netIface)
	}

	// 更新全局状态
	state.UpdateNetIO(currentStats)

	return &NetworkInfo{
		Interfaces: result,
		Timestamp:  time.Now().Unix(),
	}, nil
}

// shouldSkipInterface 判断是否应该跳过该接口
func shouldSkipInterface(iface psutilnet.InterfaceStat) bool {
	// 跳过本地回环
	if iface.Name == "lo" {
		return true
	}

	// 跳过没有 MAC 地址的接口（通常是虚拟接口）
	if iface.HardwareAddr == "" {
		return false // 可以根据需求调整
	}

	// 不跳过 docker、virbr 等虚拟接口（可以根据需求调整）
	return false
}

// isInterfaceUp 判断接口是否启用
func isInterfaceUp(iface psutilnet.InterfaceStat) bool {
	for _, flag := range iface.Flags {
		if flag == "up" {
			return true
		}
	}
	return false
}

// extractAddresses 提取 IP 地址
func extractAddresses(iface psutilnet.InterfaceStat) []string {
	addresses := make([]string, 0, len(iface.Addrs))
	for _, addr := range iface.Addrs {
		addrStr := addr.String()
		// 去除 CIDR 后缀，保留纯净 IP
		if idx := strings.Index(addrStr, "/"); idx != -1 {
			addresses = append(addresses, addrStr[:idx])
		} else {
			addresses = append(addresses, addrStr)
		}
	}
	return addresses
}

// GetNetworkStats 获取网络统计摘要
func GetNetworkStats() (*NetworkStats, error) {
	counters, err := psutilnet.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters: %w", err)
	}

	if len(counters) == 0 {
		return &NetworkStats{}, nil
	}

	totalSent := counters[0].BytesSent
	totalRecv := counters[0].BytesRecv

	// 获取活动接口数
	interfaces, err := psutilnet.Interfaces()
	if err != nil {
		return &NetworkStats{
			TotalBytesSent: totalSent,
			TotalBytesRecv: totalRecv,
			Interfaces:     0,
		}, nil
	}

	activeCount := 0
	for _, iface := range interfaces {
		if isInterfaceUp(iface) && iface.Name != "lo" {
			activeCount++
		}
	}

	return &NetworkStats{
		TotalBytesSent: totalSent,
		TotalBytesRecv: totalRecv,
		Interfaces:     activeCount,
	}, nil
}

// GetInterfaceByName 根据名称获取接口信息
func GetInterfaceByName(name string) (*Interface, error) {
	netInfo, err := GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	for _, iface := range netInfo.Interfaces {
		if iface.Name == name {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("interface %s not found", name)
}

// GetNetworkIOCounters 获取网络 IO 计数器
func GetNetworkIOCounters(perInterface bool) (map[string]NetIOHistory, error) {
	counters, err := psutilnet.IOCounters(perInterface)
	if err != nil {
		return nil, fmt.Errorf("failed to get network IO counters: %w", err)
	}

	result := make(map[string]NetIOHistory)
	for _, c := range counters {
		result[c.Name] = NetIOHistory{
			BytesSent:   c.BytesSent,
			BytesRecv:   c.BytesRecv,
			PacketsSent: c.PacketsSent,
			PacketsRecv: c.PacketsRecv,
		}
	}

	return result, nil
}

// GetConnections 获取网络连接统计
func GetConnections() (*ConnectionStats, error) {
	// 获取 TCP 连接数量
	conns, err := psutilnet.Connections("all")
	if err != nil {
		return nil, fmt.Errorf("failed to get network connections: %w", err)
	}

	stats := &ConnectionStats{
		Total: len(conns),
	}

	for _, c := range conns {
		// 使用字符串比较类型
		connType := strings.ToLower(string(c.Type))
		switch connType {
		case "tcp":
			stats.TCP++
		case "udp":
			stats.UDP++
		case "inet":
			stats.INET++
		case "inet4":
			stats.INET4++
		case "inet6":
			stats.INET6++
		}

		switch c.Status {
		case "ESTABLISHED":
			stats.Established++
		case "LISTEN":
			stats.Listen++
		case "TIME_WAIT":
			stats.TimeWait++
		case "CLOSE_WAIT":
			stats.CloseWait++
		}
	}

	return stats, nil
}

// ConnectionStats 网络连接统计
type ConnectionStats struct {
	Total      int `json:"total"`
	TCP        int `json:"tcp"`
	UDP        int `json:"udp"`
	INET       int `json:"inet"`
	INET4      int `json:"inet4"`
	INET6      int `json:"inet6"`
	Established int `json:"established"`
	Listen      int `json:"listen"`
	TimeWait    int `json:"timeWait"`
	CloseWait   int `json:"closeWait"`
}

// GetBandwidthUsage 获取带宽使用情况
func GetBandwidthUsage(ifaceName string) (*BandwidthUsage, error) {
	counters, err := psutilnet.IOCounters(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get IO counters for %s: %w", ifaceName, err)
	}

	var counter *psutilnet.IOCountersStat
	for _, c := range counters {
		if c.Name == ifaceName {
			counter = &c
			break
		}
	}

	if counter == nil {
		return nil, fmt.Errorf("interface %s not found", ifaceName)
	}

	state := GetStatsState()
	lastStats, lastTS := state.GetLastNetIO()
	now := time.Now()
	duration := now.Sub(lastTS).Seconds()

	usage := &BandwidthUsage{
		Interface: ifaceName,
		BytesSent: counter.BytesSent,
		BytesRecv: counter.BytesRecv,
		Timestamp: now.Unix(),
	}

	if lastStat, ok := lastStats[ifaceName]; ok && duration > 0 {
		sentSpeed, recvSpeed := CalculateNetSpeed(NetIOHistory{
			BytesSent: counter.BytesSent,
			BytesRecv: counter.BytesRecv,
		}, lastStat, duration)
		usage.SentSpeed = sentSpeed
		usage.RecvSpeed = recvSpeed
	}

	return usage, nil
}

// BandwidthUsage 带宽使用情况
type BandwidthUsage struct {
	Interface  string  `json:"interface"`
	BytesSent  uint64  `json:"bytesSent"`
	BytesRecv  uint64  `json:"bytesRecv"`
	SentSpeed  float64 `json:"sentSpeed"`
	RecvSpeed  float64 `json:"recvSpeed"`
	Timestamp  int64   `json:"timestamp"`
}

// FormatNetworkBytes 格式化字节数
func FormatNetworkBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatNetworkSpeed 格式化速度
func FormatNetworkSpeed(bytesPerSec float64) string {
	if bytesPerSec < 1024 {
		return fmt.Sprintf("%.1f B/s", bytesPerSec)
	}
	const unit = 1024.0
	div, exp := unit, 0
	for n := bytesPerSec / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB/s", bytesPerSec/div, "KMGTPE"[exp])
}

// GetPublicIP 获取公网 IP 地址
func GetPublicIP() (string, error) {
	// 简单实现，可以从外部服务获取
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no public IP found")
}
