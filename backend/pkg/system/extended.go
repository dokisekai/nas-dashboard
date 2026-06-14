package system

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo 进程信息
type ProcessInfo struct {
	PID           int32   `json:"pid"`
	Name          string  `json:"name"`
	Status        string  `json:"status"`
	CPUPercent    float64 `json:"cpuPercent"`
	MemoryPercent float64 `json:"memoryPercent"`
	Memory        uint64  `json:"memory"`
	Threads       int32   `json:"threads"`
	Username      string  `json:"username"`
	Command       string  `json:"command"`
	CreatedTime   int64   `json:"createdTime"`
}

// GetProcesses 获取所有进程列表
func GetProcesses() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("failed to get processes: %w", err)
	}

	processes := make([]ProcessInfo, 0, len(procs))

	for _, p := range procs {
		info, err := getProcessInfo(p)
		if err != nil {
			continue // 跳过无法访问的进程
		}
		processes = append(processes, info)
	}

	return processes, nil
}

// GetProcess 获取特定进程信息
func GetProcess(pid int32) (ProcessInfo, error) {
	p, err := process.NewProcess(pid)
	if err != nil {
		return ProcessInfo{}, fmt.Errorf("failed to create process: %w", err)
	}

	return getProcessInfo(p)
}

// getProcessInfo 获取进程详细信息
func getProcessInfo(p *process.Process) (ProcessInfo, error) {
	var info ProcessInfo

	info.PID = p.Pid

	// 获取进程名称
	if name, err := p.Name(); err == nil {
		info.Name = name
	}

	// 获取进程状态
	if status, err := p.Status(); err == nil {
		if len(status) > 0 {
			info.Status = status[0]
		} else {
			info.Status = "unknown"
		}
	}

	// 获取CPU使用率
	if cpuPercent, err := p.CPUPercent(); err == nil {
		info.CPUPercent = cpuPercent
	}

	// 获取内存信息
	if memInfo, err := p.MemoryInfo(); err == nil {
		info.Memory = memInfo.RSS
	}

	// 获取内存百分比
	if memPercent, err := p.MemoryPercent(); err == nil {
		info.MemoryPercent = float64(memPercent)
	}

	// 获取线程数
	if threads, err := p.NumThreads(); err == nil {
		info.Threads = threads
	}

	// 获取用户名
	if username, err := p.Username(); err == nil {
		info.Username = username
	}

	// 获取命令行
	if cmdline, err := p.Cmdline(); err == nil {
		info.Command = cmdline
	}

	// 获取创建时间
	if create_time, err := p.CreateTime(); err == nil {
		info.CreatedTime = create_time
	}

	return info, nil
}

// ServiceInfo 系统服务信息
type ServiceInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`     // running, stopped, failed, masked
	Enabled     bool   `json:"enabled"`
	LoadState   string `json:"loadState"`
	ActiveState string `json:"activeState"`
	MainPID     int32  `json:"mainPid"`
	SubState    string `json:"subState"`
}

// GetServices 获取所有系统服务
func GetServices() ([]ServiceInfo, error) {
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl command failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	services := make([]ServiceInfo, 0)

	for _, line := range lines {
		if strings.HasPrefix(line, "UNIT") || strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		// 解析服务名称（移除.service扩展名）
		name := strings.TrimSuffix(fields[0], ".service")

		// 解析状态
		loadState := fields[1]
		activeState := fields[2]
		subState := fields[3]

		// 确定服务状态
		status := "unknown"
		switch activeState {
		case "active":
			status = "running"
		case "inactive":
			status = "stopped"
		case "failed":
			status = "failed"
		}

		if loadState == "masked" {
			status = "masked"
		}

		// 解析主PID（从SUB字段中提取）
		var mainPID int32 = 0
		if len(fields) > 4 {
			mainPIDStr := strings.TrimPrefix(fields[4], "main_pid=")
			if pid, err := strconv.ParseInt(mainPIDStr, 10, 32); err == nil {
				mainPID = int32(pid)
			}
		}

		service := ServiceInfo{
			Name:        name,
			LoadState:   loadState,
			ActiveState: activeState,
			SubState:    subState,
			Status:      status,
			MainPID:     mainPID,
		}

		// 获取服务描述和启用状态
		if desc, enabled := getServiceDetails(name); desc != "" {
			service.Description = desc
			service.Enabled = enabled
		}

		services = append(services, service)
	}

	return services, nil
}

// GetService 获取特定服务信息
func GetService(name string) (ServiceInfo, error) {
	// 获取服务状态
	cmd := exec.Command("systemctl", "show", name+".service", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return ServiceInfo{}, fmt.Errorf("failed to get service info: %w", err)
	}

	service := ServiceInfo{Name: name}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "Description":
			service.Description = value
		case "LoadState":
			service.LoadState = value
		case "ActiveState":
			service.ActiveState = value
		case "SubState":
			service.SubState = value
		case "MainPID":
			if pid, err := strconv.ParseInt(value, 10, 32); err == nil {
				service.MainPID = int32(pid)
			}
		}
	}

	// 确定状态
	switch service.ActiveState {
	case "active":
		service.Status = "running"
	case "inactive":
		service.Status = "stopped"
	case "failed":
		service.Status = "failed"
	default:
		service.Status = "unknown"
	}

	if service.LoadState == "masked" {
		service.Status = "masked"
	}

	// 获取启用状态
	if _, enabled := getServiceDetails(name); enabled {
		service.Enabled = true
	} else {
		service.Enabled = false
	}

	return service, nil
}

// getServiceDetails 获取服务详细信息
func getServiceDetails(name string) (string, bool) {
	// 获取服务描述
	cmd := exec.Command("systemctl", "show", name+".service", "--property=Description", "--value", "--no-pager")
	output, err := cmd.Output()
	if err == nil {
		description := strings.TrimSpace(string(output))
		if description != "" {
			// 检查是否启用
			cmd = exec.Command("systemctl", "is-enabled", name+".service")
			output, err := cmd.Output()
			if err == nil && strings.TrimSpace(string(output)) == "enabled" {
				return description, true
			}
			return description, false
		}
	}

	return "", false
}

// TemperatureInfo 温度信息
type TemperatureInfo struct {
	Sensors []Sensor `json:"sensors"`
}

// Sensor 温度传感器
type Sensor struct {
	Name     string  `json:"name"`
	Current  float64 `json:"current"`
	Max      float64 `json:"max"`
	Critical float64 `json:"critical"`
	Unit     string  `json:"unit"` // C, F
}

// GetTemperature 获取系统温度
func GetTemperature() (*TemperatureInfo, error) {
	temp := &TemperatureInfo{
		Sensors: make([]Sensor, 0),
	}

	// 使用sensors命令获取温度
	cmd := exec.Command("sensors", "-j")
	output, err := cmd.Output()
	if err != nil {
		// 如果sensors不可用，尝试从/sys/class/thermal读取
		return getTemperatureFromSysfs()
	}

	// 解析JSON输出
	var sensorsData map[string]interface{}
	if err := json.Unmarshal(output, &sensorsData); err != nil {
		return getTemperatureFromSysfs()
	}

	// 解析各个传感器的温度
	for chipName, chipData := range sensorsData {
		if chipMap, ok := chipData.(map[string]interface{}); ok {
			// 直接解析chipMap中的传感器数据
			for sensorName, sensorData := range chipMap {
				if sensorMap, ok := sensorData.(map[string]interface{}); ok {
					sensor := parseSensorDataV2(chipName, sensorName, sensorMap)
					if sensor != nil {
						temp.Sensors = append(temp.Sensors, *sensor)
					}
				}
			}
		}
	}

	return temp, nil
}

// parseSensorData 解析传感器数据
func parseSensorData(chipName, adapterName string, data map[string]interface{}) *Sensor {
	// 查找温度输入
	for key, value := range data {
		if strings.HasSuffix(key, "_input") && strings.HasPrefix(key, "temp") {
			if tempValue, ok := value.(float64); ok {
				sensor := &Sensor{
					Name:    fmt.Sprintf("%s %s", chipName, adapterName),
					Current: tempValue,
					Max:     tempValue * 1.2, // 估算最大值
					Critical: tempValue * 1.5, // 估算临界值
					Unit:    "C",
				}

				// 尝试获取最大和临界值
				if maxVal, ok := data["temp"+strings.TrimPrefix(key, "temp")+"_max"]; ok {
					if maxFloat, ok := maxVal.(float64); ok {
						sensor.Max = maxFloat
					}
				}

				if critVal, ok := data["temp"+strings.TrimPrefix(key, "temp")+"_crit"]; ok {
					if critFloat, ok := critVal.(float64); ok {
						sensor.Critical = critFloat
					}
				}

				return sensor
			}
		}
	}

	return nil
}

// parseSensorDataV2 解析传感器数据（新版本，适配实际sensors -j输出格式）
func parseSensorDataV2(chipName, sensorName string, data map[string]interface{}) *Sensor {
	// 查找温度输入
	for key, value := range data {
		if strings.HasSuffix(key, "_input") && strings.HasPrefix(key, "temp") {
			if tempValue, ok := value.(float64); ok {
				sensor := &Sensor{
					Name:     fmt.Sprintf("%s - %s", chipName, sensorName),
					Current:  tempValue,
					Max:      tempValue * 1.2, // 估算最大值
					Critical: tempValue * 1.5, // 估算临界值
					Unit:     "C",
				}

				// 提取temp数字编号
				tempNum := strings.TrimPrefix(key, "temp")
				tempNum = strings.TrimSuffix(tempNum, "_input")

				// 尝试获取最大和临界值
				if maxVal, ok := data["temp"+tempNum+"_max"]; ok {
					if maxFloat, ok := maxVal.(float64); ok {
						sensor.Max = maxFloat
					}
				}

				if critVal, ok := data["temp"+tempNum+"_crit"]; ok {
					if critFloat, ok := critVal.(float64); ok {
						sensor.Critical = critFloat
					}
				}

				return sensor
			}
		}
	}

	return nil
}

// getTemperatureFromSysfs 从sysfs获取温度
func getTemperatureFromSysfs() (*TemperatureInfo, error) {
	temp := &TemperatureInfo{
		Sensors: make([]Sensor, 0),
	}

	// 读取/sys/class/thermal_zone*目录
	cmd := exec.Command("ls", "/sys/class/thermal")
	output, err := cmd.Output()
	if err != nil {
		return temp, fmt.Errorf("no thermal sensors found: %w", err)
	}

	zones := strings.Split(strings.TrimSpace(string(output)), "\n")

	for _, zone := range zones {
		if zone == "" {
			continue
		}

		// 读取温度值
		tempPath := fmt.Sprintf("/sys/class/thermal/%s/temp", zone)
		cmd := exec.Command("cat", tempPath)
		output, err := cmd.Output()
		if err != nil {
			continue
		}

		tempStr := strings.TrimSpace(string(output))
		tempValue, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			continue
		}

		// 转换为摄氏度（通常是以毫摄氏度为单位）
		tempValue = tempValue / 1000

		// 读取类型
		typePath := fmt.Sprintf("/sys/class/thermal/%s/type", zone)
		cmd = exec.Command("cat", typePath)
		output, err = cmd.Output()
		if err != nil {
			continue
		}

		sensorType := strings.TrimSpace(string(output))

		sensor := Sensor{
			Name:    fmt.Sprintf("%s (%s)", sensorType, zone),
			Current: tempValue,
			Max:     80,    // 默认最大值
			Critical: 90,    // 默认临界值
			Unit:    "C",
		}

		temp.Sensors = append(temp.Sensors, sensor)
	}

	return temp, nil
}