package system

import (
	"os/exec"
	"strings"
)

// GPUInfo GPU信息结构
type GPUInfo struct {
	Model   string `json:"model"`
	Vendor  string `json:"vendor"`
	Memory  string `json:"memory"`
	Driver  string `json:"driver"`
	Enabled bool   `json:"enabled"`
}

// GetGPUInfo 获取GPU信息
func GetGPUInfo() (*GPUInfo, error) {
 gpu := &GPUInfo{
		Enabled: false,
	}

	// 尝试使用lspci命令获取GPU信息
	output, err := exec.Command("lspci").Output()
	if err != nil {
		return gpu, nil // 返回默认值而不是错误
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "VGA") || strings.Contains(line, "Display") || strings.Contains(line, "3D") {
			// 解析GPU信息
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				gpuInfo := strings.TrimSpace(parts[1])

				// 提取厂商和型号
				if strings.Contains(gpuInfo, "NVIDIA") {
					gpu.Vendor = "NVIDIA"
					gpu.Model = extractModelName(gpuInfo, "NVIDIA")
				} else if strings.Contains(gpuInfo, "AMD") || strings.Contains(gpuInfo, "ATI") {
					gpu.Vendor = "AMD"
					gpu.Model = extractModelName(gpuInfo, "AMD")
				} else if strings.Contains(gpuInfo, "Intel") {
					gpu.Vendor = "Intel"
					gpu.Model = extractModelName(gpuInfo, "Intel")
				} else {
					gpu.Model = gpuInfo
				}

				gpu.Enabled = true
				break
			}
		}
	}

	// 如果检测到独立显卡，尝试获取显存信息
	if gpu.Enabled {
		if gpu.Vendor == "NVIDIA" {
			// 尝试使用nvidia-smi获取NVIDIA显卡信息
			nvidiaOutput, err := exec.Command("nvidia-smi", "--query-gpu=memory.total", "--format=noheader").Output()
			if err == nil {
				gpu.Memory = strings.TrimSpace(string(nvidiaOutput))
				gpu.Driver = "NVIDIA"
			}
		} else if gpu.Vendor == "AMD" {
			gpu.Driver = "AMDGPU"
		}
	}

	return gpu, nil
}

// extractModelName 从GPU信息字符串中提取型号名称
func extractModelName(info, vendor string) string {
	// 清理信息字符串
	info = strings.TrimSpace(info)
	info = strings.Replace(info, "[AMD/ATI]", "", -1)
	info = strings.Replace(info, "[NVIDIA]", "", -1)
	info = strings.Replace(info, "[Intel]", "", -1)
	info = strings.Replace(info, "Corporation", "", -1)
	info = strings.TrimSpace(info)

	// 移除常见的括号内容
	for {
		start := strings.Index(info, "[")
		end := strings.Index(info, "]")
		if start != -1 && end != -1 && end > start {
			info = info[:start] + info[end+1:]
			info = strings.TrimSpace(info)
		} else {
			break
		}
	}

	// 限制长度
	if len(info) > 50 {
		info = info[:50] + "..."
	}

	return info
}

// GetGPUDetails 获取详细的GPU信息（包括温度、使用率等）
func GetGPUDetails() (map[string]interface{}, error) {
	details := make(map[string]interface{})

	// 尝试NVIDIA显卡
	nvidiaOutput, err := exec.Command("nvidia-smi", "--query-gpu=name,temperature.gpu,utilization.gpu,memory.used,memory.total,power.draw", "--format=csv,noheader").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(nvidiaOutput)), "\n")
		if len(lines) > 0 {
			parts := strings.Split(lines[0], ",")
			if len(parts) >= 6 {
				details["vendor"] = "NVIDIA"
				details["model"] = strings.TrimSpace(parts[0])
				details["temperature"] = strings.TrimSpace(parts[1])
				details["utilization"] = strings.TrimSpace(parts[2])
				details["memory_used"] = strings.TrimSpace(parts[3])
				details["memory_total"] = strings.TrimSpace(parts[4])
				details["power_draw"] = strings.TrimSpace(parts[5])
				details["available"] = true
				return details, nil
			}
		}
	}

	// 尝试AMD显卡
	_, err = exec.Command("rocm-smi", "--showuse").Output()
	if err == nil {
		details["vendor"] = "AMD"
		details["available"] = true
		details["model"] = "AMD GPU"
		return details, nil
	}

	// 返回基本信息
	gpuInfo, _ := GetGPUInfo()
	if gpuInfo.Enabled {
		details["vendor"] = gpuInfo.Vendor
		details["model"] = gpuInfo.Model
		details["memory"] = gpuInfo.Memory
		details["available"] = true
	} else {
		details["available"] = false
		details["message"] = "未检测到独立显卡，可能使用集成显卡"
	}

	return details, nil
}

// GetSystemPower 获取系统功耗信息
func GetSystemPower() (map[string]interface{}, error) {
	powerInfo := make(map[string]interface{})

	// 尝试从不同来源获取功耗信息
	// 注意：这些功能需要特定硬件支持

	// 检查NVIDIA显卡功耗
	nvidiaPower, err := exec.Command("nvidia-smi", "--query-gpu=power.draw", "--format=csv,noheader").Output()
	if err == nil {
		powerInfo["gpu_power"] = strings.TrimSpace(string(nvidiaPower))
	}

	// 尝试读取CPU功耗（需要硬件监控工具）
	// 这里可以添加对lm-sensors等工具的支持

	powerInfo["available"] = len(powerInfo) > 0

	return powerInfo, nil
}

// EstimatePowerUsage 估算系统功耗（基于硬件配置）
func EstimatePowerUsage(cpuCores int, memoryGB float64, diskCount int) int {
	totalPower := 0

	// CPU功耗估算（每核心约15W）
	if cpuCores > 0 {
		totalPower += cpuCores * 15
	} else {
		totalPower += 65 // 默认四核CPU
	}

	// 内存功耗估算（每GB约3W）
	if memoryGB > 0 {
		totalPower += int(memoryGB * 3)
	} else {
		totalPower += 24 // 默认8GB
	}

	// 磁盘功耗（每个磁盘约5-10W）
	if diskCount > 0 {
		totalPower += diskCount * 8
	} else {
		totalPower += 10 // 默认单个磁盘
	}

	// 基础功耗（主板、网络、风扇等）
	totalPower += 50

	return totalPower
}

// GetPowerSupplyStatus 获取电源供应状态
func GetPowerSupplyStatus() (map[string]interface{}, error) {
	status := make(map[string]interface{})

	// 尝试读取PSU状态（需要特定硬件支持）
	output, err := exec.Command("sensors").Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "power") || strings.Contains(line, "fan") {
				status["sensors_available"] = true
				return status, nil
			}
		}
	}

	status["sensors_available"] = false
	status["message"] = "未检测到电源传感器"

	return status, nil
}

// GetDetailedGPUModel 获取详细的GPU型号信息
func GetDetailedGPUModel() string {
	// 使用lspci -v获取更详细的GPU信息
	output, err := exec.Command("lspci", "-v").Output()
	if err != nil {
		return "未知显卡"
	}

	lines := strings.Split(string(output), "\n")
	var currentDevice string
	var inDeviceBlock bool

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "VGA") || strings.Contains(line, "Display") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				currentDevice = strings.TrimSpace(parts[1])
				inDeviceBlock = true
			}
		} else if inDeviceBlock && strings.HasPrefix(line, "") {
			// 这是设备的详细信息
			if strings.Contains(line, "driver") {
				currentDevice += " | " + strings.TrimSpace(line)
			}
		} else if line == "" && inDeviceBlock {
			// 设备块结束
			inDeviceBlock = false
		}
	}

	if currentDevice != "" {
		return currentDevice
	}

	return "集成显卡或未检测到独立显卡"
}