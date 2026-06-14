package system

import (
	"os/exec"
	"strings"
)

// GetGPUInfoSafe 安全地获取GPU信息，带有回退机制
func GetGPUInfoSafe() (*GPUInfo, error) {
	gpu := &GPUInfo{
		Enabled: false,
	}

	// 首先尝试使用lspci命令
	output, err := exec.Command("lspci").Output()
	if err != nil {
		// lspci不可用，尝试其他方法
		return getGPUInfoFallback()
	}

	lines := strings.Split(string(output), "\n")
	foundGPU := false

	// 优先寻找独立显卡（VGA compatible controller）
	for _, line := range lines {
		if strings.Contains(line, "VGA compatible controller") ||
		   strings.Contains(line, "Display controller") && !strings.Contains(line, "Audio") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				gpuInfo := strings.TrimSpace(parts[1])

				// 解析显卡信息
				if strings.Contains(gpuInfo, "NVIDIA") {
					gpu.Vendor = "NVIDIA"
					gpu.Model = extractModelName(gpuInfo, "NVIDIA")
					foundGPU = true
					gpu.Enabled = true

					// 尝试获取NVIDIA显卡详细信息
					gpu.Memory = getNvidiaGPUMemory()
					gpu.Driver = "NVIDIA"
				} else if strings.Contains(gpuInfo, "AMD") || strings.Contains(gpuInfo, "ATI") {
					gpu.Vendor = "AMD"
					gpu.Model = extractModelName(gpuInfo, "AMD")
					foundGPU = true
					gpu.Enabled = true

					// 尝试获取AMD显卡信息
					gpu.Memory = getAMDGPUMemory()
					gpu.Driver = "AMDGPU"
				} else if strings.Contains(gpuInfo, "Intel") {
					// Intel集成显卡
					if !foundGPU { // 只在没有找到独立显卡时才显示集成显卡
						gpu.Vendor = "Intel"
						gpu.Model = extractModelName(gpuInfo, "Intel")
						gpu.Driver = "Intel"
						gpu.Enabled = true
						foundGPU = true
					}
				}
				break // 找到第一个主要GPU后退出
			}
		}

		// 如果没有找到独立显卡，再查找Display controller
		if !foundGPU {
			for _, line := range lines {
				if strings.Contains(line, "Display controller") && !strings.Contains(line, "Audio") {
					parts := strings.Split(line, ":")
					if len(parts) >= 2 {
						gpuInfo := strings.TrimSpace(parts[1])

						if strings.Contains(gpuInfo, "Intel") {
							gpu.Vendor = "Intel"
							gpu.Model = extractModelName(gpuInfo, "Intel")
							gpu.Driver = "Intel"
							gpu.Enabled = true
							foundGPU = true
							break
						}
					}
				}
			}
		}
	}

	return gpu, nil
}

// getGPUInfoFallback 当lspci不可用时的回退方法
func getGPUInfoFallback() (*GPUInfo, error) {
	gpu := &GPUInfo{
		Enabled: false,
		Vendor:  "Unknown",
		Model:  "无法检测",
	}

	// 尝试读取DRI信息
	if _, err := exec.LookPath("ls"); err == nil {
		output, err := exec.Command("ls", "/sys/class/drm/").Output()
		if err == nil {
			drmDevices := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, device := range drmDevices {
				if device != "" {
					// 存在DRM设备，说明有图形支持
					gpu.Model = "检测到图形设备 (集成显卡)"
					gpu.Vendor = "Unknown"
					gpu.Driver = "开源驱动"
					gpu.Enabled = true
					break
				}
			}
		}
	}

	// 尝试读取GPU信息从/proc
	if _, err := exec.LookPath("cat"); err == nil {
		// 尝试读取GPU设备信息
		if output, err := exec.Command("cat", "/proc/driver/nvidia/version").Output(); err == nil {
			gpu.Vendor = "NVIDIA"
			gpu.Model = "NVIDIA GPU"
			gpu.Driver = strings.TrimSpace(string(output))
			gpu.Enabled = true
		}
	}

	return gpu, nil
}

// getNvidiaGPUMemory 获取NVIDIA显卡显存信息
func getNvidiaGPUMemory() string {
	output, err := exec.Command("nvidia-smi", "--query-gpu=memory.total", "--format=csv,noheader").Output()
	if err != nil {
		return "未知"
	}

	memory := strings.TrimSpace(string(output))
	if memory != "" {
		return memory
	}

	return "未知"
}

// getAMDGPUMemory 获取AMD显卡显存信息
func getAMDGPUMemory() string {
	// AMD显卡显存信息较难获取，返回估算值
	return "GDDR6" // 默认返回常见显存类型
}


// CheckLSPCIAvailability 检查lspci命令是否可用
func CheckLSPCIAvailability() bool {
	_, err := exec.LookPath("lspci")
	return err == nil
}

// InstallLSPCIInstructions 提供lspci安装说明
func InstallLSPCIInstructions() string {
	if CheckLSPCIAvailability() {
		return "lspci命令已安装"
	}

	return `lspci命令未安装。请使用以下命令安装：

# Debian/Ubuntu系统
sudo apt-get update
sudo apt-get install -y pciutils

# RHEL/CentOS系统
sudo yum install -y pciutils

# Arch Linux
sudo pacman -S pciutils

# Fedora
sudo dnf install pciutils

安装后，需要重启后端服务器以启用GPU检测功能。`
}

// DetectGPUMethod 检测可用的GPU检测方法
func DetectGPUMethod() string {
	if CheckLSPCIAvailability() {
		return "lspci"
	}

	// 检查其他可能的检测方法
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		return "nvidia-smi"
	}

	if _, err := exec.LookPath("ls"); err == nil {
		output, _ := exec.Command("ls", "/sys/class/drm/").Output()
		if len(output) > 0 && strings.Contains(string(output), "card") {
			return "drm"
		}
	}

	return "none"
}