package system

import (
	"fmt"
	"io/ioutil"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// HardwareDetail 硬件详细信息
type HardwareDetail struct {
	Type     string `json:"type"`     // cpu, memory, disk, gpu, motherboard, etc.
	Name     string `json:"name"`
	Model    string `json:"model"`
	Vendor   string `json:"vendor"`
	Details  string `json:"details,omitempty"`
	Capacity string `json:"capacity,omitempty"`
	Usage    string `json:"usage,omitempty"`
	Status   string `json:"status"`
	Drivers  string `json:"drivers,omitempty"`
	Path     string `json:"path,omitempty"`
}

// PowerInfo 功耗信息
type PowerInfo struct {
	CPU     float64 `json:"cpu"`     // CPU功耗(瓦特)
	Memory  float64 `json:"memory"`  // 内存功耗(瓦特)
	GPU     float64 `json:"gpu"`     // GPU功耗(瓦特)
	Storage float64 `json:"storage"` // 存储功耗(瓦特)
	Mother  float64 `json:"other"`   // 其他硬件功耗(瓦特)
	Total   float64 `json:"total"`   // 总功耗(瓦特)
	Unit    string  `json:"unit"`    // 功耗单位
	Source  string  `json:"source"`  // 数据来源
}

// MotherboardInfo 主板信息
type MotherboardInfo struct {
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
	ProductName string `json:"productName"`
	BIOS        string `json:"bios"`
	Chipset     string `json:"chipset"`
}

// NetworkInterfaceDetail 网络接口详细信息
type NetworkInterfaceDetail struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Model  string `json:"model"`
	Vendor string `json:"vendor"`
	MAC    string `json:"mac"`
	Active bool   `json:"active"`
	Status string `json:"status"`
	Driver string `json:"driver"`
}

// GetRealUptime 获取真实的系统运行时间（秒）
func GetRealUptime() (float64, error) {
	data, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(data))
	if len(fields) >= 1 {
		uptime, err := strconv.ParseFloat(fields[0], 64)
		return uptime, err
	}

	return 0, fmt.Errorf("无法解析运行时间")
}

// FormatUptime 格式化运行时间为可读格式
func FormatUptime(uptimeSeconds float64) string {
	if uptimeSeconds <= 0 {
		return "Unknown"
	}

	days := int(uptimeSeconds / 86400)
	hours := int(int(uptimeSeconds) % 86400 / 3600)
	minutes := int(int(uptimeSeconds) % 3600 / 60)

	if days > 0 {
		return fmt.Sprintf("%d天 %d小时 %d分钟", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时 %d分钟", hours, minutes)
	} else {
		return fmt.Sprintf("%d分钟", minutes)
	}
}

// GetHardwareDetails 获取详细的硬件信息列表
func GetHardwareDetails() ([]HardwareDetail, error) {
	var hardware []HardwareDetail

	// 获取CPU信息
	cpuInfo, err := GetCPUInfo()
	if err == nil {
		// 获取CPU详细信息来获取供应商
		vendor := "Unknown"
		if cpuDetails, err := GetCPUDetails(); err == nil && cpuDetails != nil {
			vendor = formatVendorID(cpuDetails[0].VendorID)
		}

		hardware = append(hardware, HardwareDetail{
			Type:    "cpu",
			Name:    cpuInfo.Model,
			Model:   cpuInfo.Model,
			Vendor:  vendor,
			Details: fmt.Sprintf("核心数: %d, 主频: %.2fGHz", cpuInfo.Cores, cpuInfo.Mhz/1000),
			Status:  "正常",
			Drivers: "通用CPU驱动",
		})
	}

	// 获取内存信息
	memInfo, err := GetMemoryInfo()
	if err == nil {
		usagePercent := 0.0
		if memInfo.Total > 0 {
			usagePercent = float64(memInfo.Used) / float64(memInfo.Total) * 100
		}

		hardware = append(hardware, HardwareDetail{
			Type:     "memory",
			Name:     "系统内存",
			Model:    "DDR4",
			Vendor:   "Unknown",
			Capacity: formatBytes(memInfo.Total),
			Usage:    fmt.Sprintf("%.1f%%", usagePercent),
			Status:   "正常",
			Details:  fmt.Sprintf("可用: %s, 缓存: %s", formatBytes(memInfo.Available), formatBytes(memInfo.Cached)),
		})
	}

	// 获取GPU信息
	gpuInfo, err := GetGPUInfoSafe()
	if err == nil && gpuInfo.Enabled {
		var gpuDetails []string
		if gpuInfo.Memory != "" && gpuInfo.Memory != "Unknown" {
			gpuDetails = append(gpuDetails, fmt.Sprintf("显存: %s", gpuInfo.Memory))
		}

		hardware = append(hardware, HardwareDetail{
			Type:    "gpu",
			Name:    fmt.Sprintf("%s %s", gpuInfo.Vendor, gpuInfo.Model),
			Model:   gpuInfo.Model,
			Vendor:  gpuInfo.Vendor,
			Details: strings.Join(gpuDetails, ", "),
			Status:  "正常",
			Drivers: gpuInfo.Driver + "驱动",
		})
	}

	// 获取磁盘信息
	diskInfo, err := GetDiskInfo()
	if err == nil {
		for _, disk := range diskInfo.Disks {
			var diskDetails []string

			// 从device字段提取磁盘名称
			diskName := filepath.Base(disk.Device)

			// 获取磁盘详细信息
			diskPath := fmt.Sprintf("/sys/block/%s", diskName)
			if model, err := getDiskModel(diskPath); err == nil && model != "" {
				diskDetails = append(diskDetails, fmt.Sprintf("型号: %s", model))
			}

			if serial, err := getDiskSerial(diskPath); err == nil && serial != "" && serial != "None" {
				diskDetails = append(diskDetails, fmt.Sprintf("序列号: %s", serial))
			}

			diskDetail := HardwareDetail{
				Type:     "disk",
				Name:     diskName,
				Model:    getModelFromDiskPath(diskName),
				Vendor:   getDiskVendor(diskName),
				Capacity: formatBytes(disk.Total),
				Usage:    fmt.Sprintf("%.1f%%", disk.UsedPercent),
				Status:   "正常",
				Details:  strings.Join(diskDetails, ", "),
				Path:     disk.Mountpoint,
			}
			hardware = append(hardware, diskDetail)
		}
	}

	// 获取主板信息
	mbInfo, err := GetMotherboardInfo()
	if err == nil {
		hardware = append(hardware, HardwareDetail{
			Type:    "motherboard",
			Name:    "主板",
			Model:   mbInfo.Model,
			Vendor:  mbInfo.Vendor,
			Details: fmt.Sprintf("芯片组: %s, BIOS版本: %s", mbInfo.Chipset, mbInfo.BIOS),
			Status:  "正常",
		})
	}

	// 获取网络接口信息
	networkInfo, err := GetNetworkInterfaceDetails()
	if err == nil {
		for _, iface := range networkInfo {
			hardware = append(hardware, HardwareDetail{
				Type:    "network",
				Name:    iface.Name,
				Model:   iface.Model,
				Vendor:  iface.Vendor,
				Details: fmt.Sprintf("MAC: %s, 状态: %v", iface.MAC, iface.Active),
				Status:  iface.Status,
				Drivers: iface.Driver,
			})
		}
	}

	return hardware, nil
}

// GetRealPowerUsage 获取真实的系统功耗
func GetRealPowerUsage() (*PowerInfo, error) {
	powerInfo := &PowerInfo{
		Unit:   "Watts",
		Source: "sensor",
	}

	// 方法1: 从RAPL读取CPU功耗（Intel RAPL）
	if cpuPower, err := readRAPLPower(); err == nil {
		powerInfo.CPU = cpuPower
	}

	// 方法2: 从nvml读取NVIDIA GPU功耗
	if gpuPower, err := readNvidiaPower(); err == nil {
		powerInfo.GPU = gpuPower
	}

	// 方法3: 估算内存功耗
	if memPower, err := estimateMemoryPower(); err == nil {
		powerInfo.Memory = memPower
	}

	// 方法4: 估算存储功耗
	if storagePower, err := estimateStoragePower(); err == nil {
		powerInfo.Storage = storagePower
	}

	// 方法5: 基础系统功耗
	powerInfo.Mother = 30.0 // 主板、网络等基础功耗

	// 计算总功耗
	powerInfo.Total = powerInfo.CPU + powerInfo.Memory + powerInfo.GPU + powerInfo.Storage + powerInfo.Mother

	// 如果无法获取真实功耗，尝试估算
	if powerInfo.Total < 20 {
		if estimatedPower, err := estimateTotalPower(); err == nil {
			powerInfo.Total = estimatedPower
			powerInfo.Source = "estimated"
		}
	}

	return powerInfo, nil
}

// readRAPLPower 从Intel RAPL读取CPU功耗
func readRAPLPower() (float64, error) {
	raplPath := "/sys/class/intel_rapl"
	raplFiles, err := filepath.Glob(filepath.Join(raplPath, "*/name"))
	if err != nil || len(raplFiles) == 0 {
		return 0, fmt.Errorf("无法找到RAPL功耗信息")
	}

	var totalPower float64

	// 遍历所有RAPL域
	for _, nameFile := range raplFiles {
		// 读取域名称
		nameData, err := ioutil.ReadFile(nameFile)
		if err != nil {
			continue
		}
		domainName := strings.TrimSpace(string(nameData))

		// 只读取processor域（通常是CPU功耗）
		if strings.Contains(domainName, "package") || strings.Contains(domainName, "core") {
			energyFile := filepath.Join(filepath.Dir(nameFile), "energy_uj")
			if energyData, err := ioutil.ReadFile(energyFile); err == nil {
				energyUJ, err := strconv.ParseUint(strings.TrimSpace(string(energyData)), 10, 64)
				if err == nil {
					// 转换为焦耳，然后估算功耗（简化处理）
					totalPower += float64(energyUJ) / 1000000.0 // 转换为焦耳
				}
			}
		}
	}

	if totalPower > 0 {
		return totalPower, nil
	}

	return 0, fmt.Errorf("无法读取RAPL功耗数据")
}

// readNvidiaPower 从NVIDIA SMI读取GPU功耗
func readNvidiaPower() (float64, error) {
	// 使用nvidia-smi读取功耗
	cmd := exec.Command("nvidia-smi", "--query-gpu=power.draw", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("nvidia-smi不可用")
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) > 0 {
		powerStr := strings.TrimSpace(lines[0])
		if powerStr != "" && powerStr != "N/A" {
			// 功耗可能是mW或W
			power, err := strconv.ParseFloat(strings.TrimSuffix(powerStr, " W"), 64)
			if err == nil {
				return power, nil
			}
		}
	}

	return 0, fmt.Errorf("无法读取NVIDIA GPU功耗")
}

// estimateMemoryPower 估算内存功耗
func estimateMemoryPower() (float64, error) {
	// 获取内存信息
	memInfo, err := GetMemoryInfo()
	if err != nil {
		return 24.0, nil // 默认24W
	}

	// 基于内存大小估算功耗
	// 假设DDR4内存的功耗约为每GB 2-3W
	memGB := float64(memInfo.Total) / (1024 * 1024 * 1024)
	if memGB <= 0 {
		return 24.0, nil
	}

	// 基础功耗 + 负载功耗
	basePower := memGB * 2.0 // 基础功耗
	usagePercent := 0.0
	if memInfo.Total > 0 {
		usagePercent = float64(memInfo.Used) / float64(memInfo.Total) * 100
	}
	usageMultiplier := 1.0 + (usagePercent / 100.0) * 0.3 // 使用率影响

	return basePower * usageMultiplier, nil
}

// estimateStoragePower 估算存储功耗
func estimateStoragePower() (float64, error) {
	diskInfo, err := GetDiskInfo()
	if err != nil {
		return 15.0, nil // 默认15W
	}

	var totalPower float64
	for range diskInfo.Disks {
		// 每个磁盘约5-10W
		totalPower += 8.0
	}

	return totalPower, nil
}

// estimateTotalPower 估算系统总功耗（如果没有真实传感器）
func estimateTotalPower() (float64, error) {
	var totalPower float64 = 0

	// CPU功耗（基于核心数）
	cpuInfo, err := GetCPUInfo()
	if err == nil {
		if cpuInfo.Cores > 0 {
			totalPower += float64(cpuInfo.Cores) * 15 // 每核心15W
		}
	}

	// 内存功耗
	memInfo, err := GetMemoryInfo()
	if err == nil {
		memGB := float64(memInfo.Total) / (1024 * 1024 * 1024)
		if memGB > 0 {
			totalPower += memGB * 2.5 // 每GB 2.5W
		}
	}

	// GPU功耗（检测到的GPU）
	gpuInfo, err := GetGPUInfoSafe()
	if err == nil && gpuInfo.Enabled {
		// 根据GPU类型估算功耗
		if strings.Contains(gpuInfo.Model, "RTX") || strings.Contains(gpuInfo.Model, "Radeon") {
			totalPower += 200 // 高性能独立显卡
		} else if strings.Contains(gpuInfo.Model, "Graphics") {
			totalPower += 15 // 集成显卡
		}
	}

	// 存储功耗
	diskInfo, err := GetDiskInfo()
	if err == nil {
		for range diskInfo.Disks {
			totalPower += 8 // 每个磁盘约8W
		}
	}

	// 基础系统功耗（主板、网络、风扇等）
	totalPower += 50

	return totalPower, nil
}

// getDiskModel 获取磁盘型号
func getDiskModel(diskPath string) (string, error) {
	modelPath := filepath.Join(diskPath, "device/model")
	data, err := ioutil.ReadFile(modelPath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

// getDiskSerial 获取磁盘序列号
func getDiskSerial(diskPath string) (string, error) {
	serialPath := filepath.Join(diskPath, "device/serial")
	data, err := ioutil.ReadFile(serialPath)
	if err != nil {
		return "", err
	}

	serial := strings.TrimSpace(string(data))
	if serial == "" || serial == "None" {
		return "Unknown", nil
	}

	return serial, nil
}

// getDiskVendor 获取磁盘厂商
func getDiskVendor(diskName string) string {
	// 简化的磁盘厂商检测
	vendors := map[string]string{
		"Samsung":      "Samsung",
		"Western":      "Western Digital",
		"Seagate":      "Seagate",
		"Toshiba":      "Toshiba",
		"Crucial":      "Crucial",
		"Kingston":     "Kingston",
		"MAXIO":        "MAXIO",
		"Silicon Power": "Silicon Power",
	}

	for _, vendor := range vendors {
		if strings.Contains(strings.ToLower(diskName), strings.ToLower(vendor)) {
			return vendor
		}
	}

	return "Unknown"
}

// getModelFromDiskPath 从磁盘路径生成型号信息
func getModelFromDiskPath(diskName string) string {
	return fmt.Sprintf("存储设备 %s", diskName)
}

// formatVendorID 格式化CPU供应商ID
func formatVendorID(vendorID string) string {
	vendors := map[string]string{
		"GenuineIntel": "Intel",
		"AuthenticAMD": "AMD",
		"CentaurHauls": "Centaur",
		"HygonGenuine": "Hygon",
	}

	if vendor, ok := vendors[vendorID]; ok {
		return vendor
	}

	return "Unknown"
}

// GetMotherboardInfo 获取主板信息
func GetMotherboardInfo() (*MotherboardInfo, error) {
	mbInfo := &MotherboardInfo{}

	// 读取/sys/class/dmi/id下的信息
	if vendorData, err := ioutil.ReadFile("/sys/class/dmi/id/board_vendor"); err == nil {
		mbInfo.Vendor = strings.TrimSpace(string(vendorData))
	}

	if modelData, err := ioutil.ReadFile("/sys/class/dmi/id/board_name"); err == nil {
		mbInfo.Model = strings.TrimSpace(string(modelData))
	}

	if productData, err := ioutil.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		mbInfo.ProductName = strings.TrimSpace(string(productData))
	}

	if biosData, err := ioutil.ReadFile("/sys/class/dmi/id/bios_version"); err == nil {
		mbInfo.BIOS = strings.TrimSpace(string(biosData))
	}

	// 设置默认值
	if mbInfo.Vendor == "" {
		mbInfo.Vendor = "Unknown"
	}
	if mbInfo.Model == "" {
		mbInfo.Model = "Unknown Motherboard"
	}
	if mbInfo.BIOS == "" {
		mbInfo.BIOS = "Unknown BIOS"
	}
	if mbInfo.Chipset == "" {
		mbInfo.Chipset = "Unknown Chipset"
	}

	return mbInfo, nil
}

// GetNetworkInterfaceDetails 获取网络接口详细信息
func GetNetworkInterfaceDetails() ([]NetworkInterfaceDetail, error) {
	var interfaces []NetworkInterfaceDetail

	// 从/sys/class/net读取网络接口
	netPath := "/sys/class/net"
	dirs, err := ioutil.ReadDir(netPath)
	if err != nil {
		return interfaces, err
	}

	for _, dir := range dirs {
		if dir.Name() == "lo" {
			continue // 跳过本地回环
		}

		interfacePath := filepath.Join(netPath, dir.Name())

		// 读取接口状态
		operstatePath := filepath.Join(interfacePath, "operstate")
		operstateData, _ := ioutil.ReadFile(operstatePath)
		operstate := strings.TrimSpace(string(operstateData))

		// 读取MAC地址
		addressPath := filepath.Join(interfacePath, "address")
		macData, _ := ioutil.ReadFile(addressPath)
		mac := strings.TrimSpace(string(macData))

		// 读取接口类型
		interfaceType := "Ethernet"
		if strings.HasPrefix(dir.Name(), "wl") {
			interfaceType = "Wireless"
		}

		interfaceDetail := NetworkInterfaceDetail{
			Name:   dir.Name(),
			Type:   interfaceType,
			Model:  "网络接口",
			Vendor: "Network",
			MAC:    mac,
			Active: operstate == "up",
			Status: getStatusFromOperstate(operstate),
			Driver: "通用网络驱动",
		}

		interfaces = append(interfaces, interfaceDetail)
	}

	return interfaces, nil
}

// getStatusFromOperstate 从operstate获取状态描述
func getStatusFromOperstate(operstate string) string {
	switch operstate {
	case "up":
		return "已连接"
	case "down":
		return "已断开"
	case "testing":
		return "测试中"
	default:
		return "未知"
	}
}

// formatBytes 格式化字节数为可读格式
func formatBytes(bytes uint64) string {
	if bytes == 0 {
		return "0 B"
	}

	var units = []string{"B", "KB", "MB", "GB", "TB", "PB"}
	const unit = uint64(1024)

	for i := 1; i < len(units); i++ {
		value := float64(bytes) / math.Pow(float64(unit), float64(i))
		if value < 1024 {
			return fmt.Sprintf("%.1f %s", value, units[i])
		}
	}

	return fmt.Sprintf("%.1f %s", float64(bytes)/math.Pow(float64(unit), 5), "PB")
}