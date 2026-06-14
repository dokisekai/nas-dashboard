package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"nas-dashboard/pkg/system"

	"github.com/gin-gonic/gin"
)

// Service 系统服务
type Service struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Status      string `json:"status"`       // running, stopped, failed, etc.
	Enabled     bool   `json:"enabled"`      // 是否开机自启
	Description string `json:"description"`
	Loaded      string `json:"loaded"`       // 配置文件是否已加载
}

// Container Docker 容器
type Container struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	State      string `json:"state"`      // running, exited, paused, etc.
	Status     string `json:"status"`     // 详细状态信息
	Ports      string `json:"ports"`      // 端口映射
	Created    string `json:"created"`    // 创建时间
	Command    string `json:"command"`    // 运行命令
}

// ServiceActionResponse 服务操作响应
type ServiceActionResponse struct {
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
}

// GetServices 获取系统服务列表
func GetServices(c *gin.Context) {
	services, err := getSystemServices()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get services: %v", err)})
		return
	}
	c.JSON(200, gin.H{"services": services})
}

// getSystemServices 获取系统服务列表（使用 systemctl）
func getSystemServices() ([]Service, error) {
	// 使用 systemctl list-units 获取活动服务
	cmd := exec.Command("systemctl", "list-units", "--type=service", "--all", "--no-pager", "--plain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("systemctl command failed: %w", err)
	}

	return parseSystemctlList(output)
}

// parseSystemctlList 解析 systemctl list-units 输出
func parseSystemctlList(output []byte) ([]Service, error) {
	var services []Service
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		// 跳过标题行和空行
		if strings.HasPrefix(line, "UNIT") || strings.HasPrefix(line, "--") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		name := strings.TrimSuffix(fields[0], ".service")
		loaded := fields[1]
		active := fields[2]
		_ = fields[3] // sub field - not used currently
		description := ""
		if len(fields) > 4 {
			description = strings.Join(fields[4:], " ")
		}

		// 确定服务状态
		status := "unknown"
		enabled := false

		if active == "active" {
			status = "running"
		} else if active == "inactive" {
			status = "stopped"
		} else if active == "failed" {
			status = "failed"
		}

		// 获取是否启用（开机自启）
		if isServiceEnabled(name) {
			enabled = true
		}

		service := Service{
			Name:        name,
			DisplayName: name,
			Status:      status,
			Enabled:     enabled,
			Description: description,
			Loaded:      loaded,
		}

		services = append(services, service)
	}

	return services, nil
}

// isServiceEnabled 检查服务是否启用
func isServiceEnabled(serviceName string) bool {
	cmd := exec.Command("systemctl", "is-enabled", serviceName+".service")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "enabled"
}

// StartService 启动服务
func StartService(c *gin.Context) {
	name := c.Param("name")

	// 验证服务是否存在
	if !serviceExists(name) {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	// 启动服务
	cmd := exec.Command("systemctl", "start", name+".service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to start service: %v, output: %s", err, string(output))})
		return
	}

	// 获取最新状态
	status := getServiceStatus(name)

	c.JSON(200, ServiceActionResponse{
		Message: fmt.Sprintf("Service %s started successfully", name),
		Status:  status,
	})
}

// StopService 停止服务
func StopService(c *gin.Context) {
	name := c.Param("name")

	// 验证服务是否存在
	if !serviceExists(name) {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	cmd := exec.Command("systemctl", "stop", name+".service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to stop service: %v, output: %s", err, string(output))})
		return
	}

	status := getServiceStatus(name)

	c.JSON(200, ServiceActionResponse{
		Message: fmt.Sprintf("Service %s stopped successfully", name),
		Status:  status,
	})
}

// RestartService 重启服务
func RestartService(c *gin.Context) {
	name := c.Param("name")

	// 验证服务是否存在
	if !serviceExists(name) {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	cmd := exec.Command("systemctl", "restart", name+".service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to restart service: %v, output: %s", err, string(output))})
		return
	}

	status := getServiceStatus(name)

	c.JSON(200, ServiceActionResponse{
		Message: fmt.Sprintf("Service %s restarted successfully", name),
		Status:  status,
	})
}

// EnableService 启用服务（开机自启）
func EnableService(c *gin.Context) {
	name := c.Param("name")

	// 验证服务是否存在
	if !serviceExists(name) {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	cmd := exec.Command("systemctl", "enable", name+".service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to enable service: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Service %s enabled successfully", name)})
}

// DisableService 禁用服务（取消开机自启）
func DisableService(c *gin.Context) {
	name := c.Param("name")

	// 验证服务是否存在
	if !serviceExists(name) {
		c.JSON(404, gin.H{"error": "Service not found"})
		return
	}

	cmd := exec.Command("systemctl", "disable", name+".service")
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to disable service: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": fmt.Sprintf("Service %s disabled successfully", name)})
}

// serviceExists 检查服务是否存在
func serviceExists(serviceName string) bool {
	cmd := exec.Command("systemctl", "list-unit-files", serviceName+".service")
	err := cmd.Run()
	return err == nil
}

// getServiceStatus 获取服务状态
func getServiceStatus(serviceName string) string {
	cmd := exec.Command("systemctl", "is-active", serviceName+".service")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// GetContainers 获取 Docker 容器列表
func GetContainers(c *gin.Context) {
	containers, err := getDockerContainers()
	if err != nil {
		// 检查 Docker 是否可用
		if strings.Contains(err.Error(), "docker") || !isDockerAvailable() {
			c.JSON(503, gin.H{"error": "Docker is not available on this system"})
			return
		}
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get containers: %v", err)})
		return
	}
	c.JSON(200, gin.H{"containers": containers})
}

// getDockerContainers 获取 Docker 容器列表
func getDockerContainers() ([]Container, error) {
	// 使用 docker ps 命令获取容器列表
	// --format json 格式化输出为 JSON
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{json .}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker command failed: %w", err)
	}

	return parseDockerPS(output)
}

// parseDockerPS 解析 docker ps 输出
func parseDockerPS(output []byte) ([]Container, error) {
	lines := strings.Split(string(output), "\n")
	var containers []Container

	for _, line := range lines {
		if line == "" {
			continue
		}

		// 解析 JSON 输出
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		container := Container{
			ID:     data["ID"].(string),
			Name:   strings.TrimPrefix(data["Names"].(string), "/"),
			Image:  data["Image"].(string),
			State:  data["State"].(string),
			Status: data["Status"].(string),
		}

		if ports, ok := data["Ports"].(string); ok {
			container.Ports = ports
		} else {
			container.Ports = ""
		}

		containers = append(containers, container)
	}

	return containers, nil
}

// isDockerAvailable 检查 Docker 是否可用
func isDockerAvailable() bool {
	// 检查 docker 命令是否存在
	if _, err := exec.LookPath("docker"); err != nil {
		return false
	}

	// 检查 Docker daemon 是否运行
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// StartContainer 启动容器
func StartContainer(c *gin.Context) {
	id := c.Param("id")

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	cmd := exec.Command("docker", "start", id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to start container: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Container started successfully", "id": id})
}

// StopContainer 停止容器
func StopContainer(c *gin.Context) {
	id := c.Param("id")

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	cmd := exec.Command("docker", "stop", id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to stop container: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Container stopped successfully", "id": id})
}

// RestartContainer 重启容器
func RestartContainer(c *gin.Context) {
	id := c.Param("id")

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	cmd := exec.Command("docker", "restart", id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to restart container: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Container restarted successfully", "id": id})
}

// RemoveContainer 删除容器
func RemoveContainer(c *gin.Context) {
	id := c.Param("id")

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	// 先停止容器（如果正在运行）
	cmd := exec.Command("docker", "rm", "-f", id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to remove container: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Container removed successfully", "id": id})
}

// containerExists 检查容器是否存在
func containerExists(id string) bool {
	cmd := exec.Command("docker", "inspect", id)
	err := cmd.Run()
	return err == nil
}

// GetContainerLogs 获取容器日志
func GetContainerLogs(c *gin.Context) {
	id := c.Param("id")
	tail := c.DefaultQuery("tail", "100")

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	cmd := exec.Command("docker", "logs", "--tail", tail, id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 日志命令可能返回非零退出码，但仍可能有输出
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get logs: %v", err)})
		return
	}

	c.JSON(200, gin.H{"id": id, "logs": string(output)})
}

// GetContainerStats 获取容器统计信息
func GetContainerStats(c *gin.Context) {
	id := c.Param("id")

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	// 使用 docker stats 获取统计信息
	cmd := exec.Command("docker", "stats", "--no-stream", "--format", "{{json .}}", id)
	output, err := cmd.Output()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get stats: %v", err)})
		return
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(output, &stats); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse stats"})
		return
	}

	c.JSON(200, gin.H{"id": id, "stats": stats})
}

// GetDockerImages 获取 Docker 镜像列表
func GetDockerImages(c *gin.Context) {
	images, err := getDockerImages()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get images: %v", err)})
		return
	}
	c.JSON(200, gin.H{"images": images})
}

// getDockerImages 获取 Docker 镜像列表
func getDockerImages() ([]map[string]interface{}, error) {
	cmd := exec.Command("docker", "images", "--format", "{{json .}}")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("docker images command failed: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var images []map[string]interface{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		images = append(images, data)
	}

	return images, nil
}

// RemoveImage 删除 Docker 镜像
func RemoveImage(c *gin.Context) {
	id := c.Param("id")

	cmd := exec.Command("docker", "rmi", id)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to remove image: %v, output: %s", err, string(output))})
		return
	}

	c.JSON(200, gin.H{"message": "Image removed successfully", "id": id})
}

// GetSystemInfo 获取系统信息（Docker、系统服务等）
func GetSystemInfo(c *gin.Context) {
	info := make(map[string]interface{})

	// Docker 可用性
	info["docker"] = isDockerAvailable()

	// 系统信息
	if uname, err := getSystemUname(); err == nil {
		info["system"] = uname
	}

	// 主机名
	if hostname, err := os.Hostname(); err == nil {
		info["hostname"] = hostname
	}

	// 获取CPU信息
	if cpuInfo, err := system.GetCPUInfo(); err == nil {
		info["cpu"] = map[string]interface{}{
			"model":       cpuInfo.Model,
			"cores":       cpuInfo.Cores,
			"mhz":         cpuInfo.Mhz,
			"usage":       cpuInfo.Usage,
			"load1":       cpuInfo.Load1,
			"load5":       cpuInfo.Load5,
			"load15":      cpuInfo.Load15,
			"architecture": cpuInfo.Model,
		}
	}

	// 获取GPU信息
	if gpuInfo, err := system.GetGPUInfo(); err == nil {
		info["gpu"] = map[string]interface{}{
			"model":   gpuInfo.Model,
			"vendor":  gpuInfo.Vendor,
			"memory":  gpuInfo.Memory,
			"driver":  gpuInfo.Driver,
			"enabled": gpuInfo.Enabled,
		}
	}

	// 获取内存信息
	if memInfo, err := system.GetMemoryInfo(); err == nil {
		info["memory"] = map[string]interface{}{
			"total":     memInfo.Total,
			"used":      memInfo.Used,
			"available": memInfo.Available,
			"percent":   memInfo.Percent,
		}
	}

	// 获取磁盘信息
	if diskInfo, err := system.GetDiskInfo(); err == nil {
		disks := make([]map[string]interface{}, 0)
		for _, disk := range diskInfo.Disks {
			disks = append(disks, map[string]interface{}{
				"device":      disk.Device,
				"mountpoint":  disk.Mountpoint,
				"total":       disk.Total,
				"used":        disk.Used,
				"free":        disk.Free,
				"usedPercent": disk.UsedPercent,
				"fstype":      disk.Fstype,
			})
		}
		info["disks"] = disks
	}

	// 系统运行时间（使用真实的秒数）
	if uptimeSeconds, err := system.GetRealUptime(); err == nil {
		info["uptimeSeconds"] = uptimeSeconds
		info["uptime"] = system.FormatUptime(uptimeSeconds)
	}

	// 系统负载信息
	if loadAvg, err := getLoadAverage(); err == nil {
		info["loadAverage"] = loadAvg
	}

	// 获取详细硬件信息
	if hardwareDetails, err := system.GetHardwareDetails(); err == nil {
		info["hardwareDetails"] = hardwareDetails
	}

	// 获取真实功耗信息
	if powerInfo, err := system.GetRealPowerUsage(); err == nil {
		info["powerUsage"] = powerInfo
	}

	c.JSON(200, gin.H{"info": info})
}

// getSystemUname 获取系统信息
func getSystemUname() (map[string]string, error) {
	cmd := exec.Command("uname", "-a")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return map[string]string{"uname": strings.TrimSpace(string(output))}, nil
}

// ExecInContainer 在容器中执行命令
func ExecInContainer(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Command string `json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// 检查容器是否存在
	if !containerExists(id) {
		c.JSON(404, gin.H{"error": "Container not found"})
		return
	}

	// 使用 docker exec 执行命令
	cmd := exec.Command("docker", "exec", id, "sh", "-c", req.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Command failed: %v", err), "output": string(output)})
		return
	}

	c.JSON(200, gin.H{"id": id, "command": req.Command, "output": string(output)})
}

// PullImage 拉取 Docker 镜像
func PullImage(c *gin.Context) {
	image := c.Query("image")
	if image == "" {
		c.JSON(400, gin.H{"error": "Image name is required"})
		return
	}

	// 创建上下文用于超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "pull", image)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to pull image: %v", err), "output": string(output)})
		return
	}

	c.JSON(200, gin.H{"message": "Image pulled successfully", "image": image})
}

// getSystemUptime 获取系统运行时间
func getSystemUptime() (string, error) {
	// 读取 /proc/uptime 获取系统运行时间（秒）
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}

	var uptime float64
	if _, err := fmt.Sscanf(string(data), "%f", &uptime); err != nil {
		return "", err
	}

	// 转换秒数为人类可读格式
	days := int(uptime / 86400)
	hours := int(int(uptime) % 86400 / 3600)
	minutes := int(int(uptime) % 3600 / 60)

	if days > 0 {
		return fmt.Sprintf("%d days, %d hours", days, hours), nil
	} else if hours > 0 {
		return fmt.Sprintf("%d hours, %d minutes", hours, minutes), nil
	} else {
		return fmt.Sprintf("%d minutes", minutes), nil
	}
}

// getLoadAverage 获取系统负载平均值
func getLoadAverage() (string, error) {
	// 读取 /proc/loadavg 获取系统负载
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "", err
	}

	// 格式: "0.50 0.70 0.60 2/120 5000"
	// 我们只需要前三个负载值（1分钟、5分钟、15分钟）
	fields := strings.Fields(string(data))
	if len(fields) >= 3 {
		return fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2]), nil
	}

	return "0.00 0.00 0.00", nil
}

// GetHardwareDetails 获取详细硬件信息
func GetHardwareDetails(c *gin.Context) {
	hardwareDetails, err := system.GetHardwareDetails()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get hardware details: %v", err)})
		return
	}

	c.JSON(200, gin.H{"hardware": hardwareDetails})
}

// GetPowerUsage 获取系统功耗信息
func GetPowerUsage(c *gin.Context) {
	powerInfo, err := system.GetRealPowerUsage()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get power usage: %v", err)})
		return
	}

	c.JSON(200, powerInfo)
}

// GetSystemUptime 获取系统运行时间
func GetSystemUptime(c *gin.Context) {
	uptimeSeconds, err := system.GetRealUptime()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to get system uptime: %v", err)})
		return
	}

	response := map[string]interface{}{
		"seconds":   uptimeSeconds,
		"formatted": system.FormatUptime(uptimeSeconds),
	}

	c.JSON(200, response)
}
