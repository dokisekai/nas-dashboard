package service

import (
	"context"
	"encoding/json"
	"log"
	"math"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/models"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/host"
	"gorm.io/gorm"
)

// AlertService 告警服务
type AlertService struct {
	db            *gorm.DB
	ctx           context.Context
	cancel        context.CancelFunc
	alertChannels map[string]chan models.Alert
	mu            sync.RWMutex
}

// SystemMetrics 系统指标
type SystemMetrics struct {
	CPU        CPU Metrics `json:"cpu"`
	Memory     MemoryMetrics `json:"memory"`
	Disk       []DiskMetrics `json:"disk"`
	Load       LoadMetrics `json:"load"`
	Temperature []TemperatureMetrics `json:"temperature"`
	Network    NetworkMetrics `json:"network"`
	Timestamp  time.Time `json:"timestamp"`
}

// CPUMetrics CPU指标
type CPUMetrics struct {
	Usage     float64 `json:"usage"`
	Cores     int     `json:"cores"`
	ModelName string  `json:"modelName"`
	MHz       float64 `json:"mhz"`
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Available uint64  `json:"available"`
	Usage     float64 `json:"usage"`
	 Cached   uint64  `json:"cached"`
	 Buffers  uint64  `json:"buffers"`
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mountPoint"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Usage       float64 `json:"usage"`
	FileSystem  string  `json:"fileSystem"`
	SmartStatus string  `json:"smartStatus"`
	Temperature int     `json:"temperature"`
}

// LoadMetrics 负载指标
type LoadMetrics struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// TemperatureMetrics 温度指标
type TemperatureMetrics struct {
	Label      string  `json:"label"`
	Temperature float64 `json:"temperature"`
	High       float64 `json:"high"`
	Critical   float64 `json:"critical"`
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	BytesSent     uint64 `json:"bytesSent"`
	BytesRecv     uint64 `json:"bytesRecv"`
	PacketsSent   uint64 `json:"packetsSent"`
	PacketsRecv   uint64 `json:"packetsRecv"`
	ErrIn        uint64 `json:"errIn"`
	ErrOut       uint64 `json:"errOut"`
	DropIn       uint64 `json:"dropIn"`
	DropOut      uint64 `json:"dropOut"`
}

// AlertRule 告警规则（扩展模型）
type AlertRule struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"` // cpu, memory, disk, load, temperature, network
	Condition   string  `json:"condition"` // gt, lt, eq, gte, lte
	Threshold   float64 `json:"threshold"`
	Duration    int     `json:"duration"` // 持续时间（秒）
	Severity    string  `json:"severity"` // info, warning, critical
	Enabled     bool    `json:"enabled"`
	Cooldown    int     `json:"cooldown"` // 冷却时间（秒）
	LastTriggered time.Time `json:"lastTriggered"`
}

// NewAlertService 创建告警服务
func NewAlertService(db *gorm.DB) *AlertService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &AlertService{
		db:            db,
		ctx:           ctx,
		cancel:        cancel,
		alertChannels: make(map[string]chan models.Alert),
	}

	// 启动监控任务
	go service.startMonitoring()

	return service
}

// startMonitoring 启动监控任务
func (s *AlertService) startMonitoring() {
	log.Println("Starting system monitoring service")

	// 启动各种监控任务
	go s.monitorCPU()
	go s.monitorMemory()
	go s.monitorDisk()
	go s.monitorLoad()
	go s.monitorTemperature()
	go s.monitorNetwork()

	// 定期收集系统指标
	go s.collectSystemMetrics()
}

// collectSystemMetrics 收集系统指标
func (s *AlertService) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			metrics, err := s.getSystemMetrics()
			if err != nil {
				log.Printf("Failed to collect system metrics: %v", err)
				continue
			}

			// 保存到数据库
			s.saveMetrics(metrics)

			// 检查告警规则
			s.checkAlertRules(metrics)
		}
	}
}

// getSystemMetrics 获取系统指标
func (s *AlertService) getSystemMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
	}

	// CPU 信息
	cpuPercent, _ := cpu.Percent(time.Second, false)
	cpuInfo, _ := cpu.Info()
	var totalCores int
	if len(cpuPercent) > 0 {
		if len(cpuInfo) > 0 {
			totalCores = len(cpuInfo)
			metrics.CPU = CPUMetrics{
				Usage:     cpuPercent[0],
				Cores:     totalCores,
				ModelName: cpuInfo[0].ModelName,
				MHz:       cpuInfo[0].Mhz,
			}
		}
	}

	// 内存信息
	memStat, _ := mem.VirtualMemory()
	metrics.Memory = MemoryMetrics{
		Total:     memStat.Total,
		Used:      memStat.Used,
		Available: memStat.Available,
		Usage:     memStat.UsedPercent,
		Cached:    memStat.Cached,
		Buffers:   memStat.Buffers,
	}

	// 磁盘信息
	diskStats, _ := disk.Partitions(false)
	for _, stat := range diskStats {
		usage, _ := disk.Usage(stat.Mountpoint)
		if usage.Total > 0 {
			diskMetric := DiskMetrics{
				Device:     stat.Device,
				MountPoint: stat.Mountpoint,
				Total:      usage.Total,
				Used:       usage.Used,
				Free:       usage.Free,
				Usage:      usage.UsedPercent,
				FileSystem: stat.Fstype,
			}
			metrics.Disk = append(metrics.Disk, diskMetric)
		}
	}

	// 负载信息
	loadStat, _ := load.Avg()
	metrics.Load = LoadMetrics{
		Load1:  loadStat.Load1,
		Load5:  loadStat.Load5,
		Load15: loadStat.Load15,
	}

	// 温度信息
	temps, _ := s.getTemperatureInfo()
	metrics.Temperature = temps

	// 网络信息
	netStats, _ := s.getNetworkStats()
	metrics.Network = netStats

	return metrics, nil
}

// saveMetrics 保存指标到数据库
func (s *AlertService) saveMetrics(metrics *SystemMetrics) {
	// 这里可以将指标保存到数据库用于历史查询
	// 为了避免数据库膨胀，可以考虑只保存汇总数据
}

// checkAlertRules 检查告警规则
func (s *AlertService) checkAlertRules(metrics *SystemMetrics) {
	// 获取所有启用的告警规则
	var rules []models.AlertRule
	if err := s.db.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		log.Printf("Failed to get alert rules: %v", err)
		return
	}

	for _, rule := range rules {
		if s.shouldTriggerAlert(rule, metrics) {
			s.triggerAlert(rule, metrics)
		}
	}
}

// shouldTriggerAlert 判断是否应该触发告警
func (s *AlertService) shouldTriggerAlert(rule models.AlertRule, metrics *SystemMetrics) bool {
	// 检查冷却时间
	if !rule.LastTriggered.IsZero() && time.Since(rule.LastTriggered) < time.Duration(rule.Cooldown)*time.Second {
		return false
	}

	var currentValue float64

	switch rule.Type {
	case "cpu":
		currentValue = metrics.CPU.Usage
	case "memory":
		currentValue = metrics.Memory.Usage
	case "load":
		currentValue = metrics.Load.Load1
	case "disk":
		// 对于磁盘，需要找到对应设备的当前值
		// 这里简化处理，使用第一个磁盘的值
		if len(metrics.Disk) > 0 {
			currentValue = metrics.Disk[0].Usage
		}
	default:
		return false
	}

	// 根据条件类型进行比较
	switch rule.Condition {
	case "gt":
		return currentValue > rule.Threshold
	case "lt":
		return currentValue < rule.Threshold
	case "gte":
		return currentValue >= rule.Threshold
	case "lte":
		return currentValue <= rule.Threshold
	case "eq":
		return currentValue == rule.Threshold
	default:
		return false
	}
}

// triggerAlert 触发告警
func (s *AlertService) triggerAlert(rule models.AlertRule, metrics *SystemMetrics) {
	alert := models.Alert{
		RuleID:    rule.ID,
		Type:      rule.Type,
		Severity:  rule.Severity,
		Message:   s.generateAlertMessage(rule, metrics),
		Status:    "active",
		Triggered: time.Now(),
	}

	// 保存到数据库
	if err := s.db.Create(&alert).Error; err != nil {
		log.Printf("Failed to save alert: %v", err)
		return
	}

	// 更新规则的最后触发时间
	s.db.Model(&models.AlertRule{}).Where("id = ?", rule.ID).Update("last_triggered", time.Now())

	// 通过通知服务发送告警
	// 这里需要与通知服务集成

	log.Printf("Alert triggered: %s - %s", rule.Name, alert.Message)
}

// generateAlertMessage 生成告警消息
func (s *AlertService) generateAlertMessage(rule models.AlertRule, metrics *SystemMetrics) string {
	conditionText := ""
	switch rule.Condition {
	case "gt":
		conditionText = "超过"
	case "lt":
		conditionText = "低于"
	case "gte":
		conditionText = "超过或等于"
	case "lte":
		conditionText = "低于或等于"
	case "eq":
		conditionText = "等于"
	}

	typeName := ""
	var currentValue float64

	switch rule.Type {
	case "cpu":
		typeName = "CPU使用率"
		currentValue = metrics.CPU.Usage
	case "memory":
		typeName = "内存使用率"
		currentValue = metrics.Memory.Usage
	case "load":
		typeName = "系统负载"
		currentValue = metrics.Load.Load1
	case "disk":
		typeName = "磁盘使用率"
		if len(metrics.Disk) > 0 {
			currentValue = metrics.Disk[0].Usage
		}
	}

	return fmt.Sprintf("%s %s %.1f%%，当前值为 %.1f%%", typeName, conditionText, rule.Threshold, currentValue)
}

// monitorCPU 监控CPU
func (s *AlertService) monitorCPU() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			percent, _ := cpu.Percent(time.Second, false)
			if len(percent) > 0 {
				log.Printf("CPU usage: %.1f%%", percent[0])
			}
		}
	}
}

// monitorMemory 监控内存
func (s *AlertService) monitorMemory() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			memStat, _ := mem.VirtualMemory()
			log.Printf("Memory usage: %.1f%% (%.1f GB / %.1f GB)",
				memStat.UsedPercent,
				float64(memStat.Used)/1024/1024/1024,
				float64(memStat.Total)/1024/1024/1024)
		}
	}
}

// monitorDisk 监控磁盘
func (s *AlertService) monitorDisk() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			partitions, _ := disk.Partitions(false)
			for _, partition := range partitions {
				usage, _ := disk.Usage(partition.Mountpoint)
				if usage.Total > 0 {
					log.Printf("Disk %s usage: %.1f%%", partition.Device, usage.UsedPercent)
				}
			}
		}
	}
}

// monitorLoad 监控负载
func (s *AlertService) monitorLoad() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			loadStat, _ := load.Avg()
			cpuCount, _ := cpu.Counts(true)
			log.Printf("Load average: %.2f, %.2f, %.2f (CPU cores: %d)",
				loadStat.Load1, loadStat.Load5, loadStat.Load15, cpuCount)
		}
	}
}

// monitorTemperature 监控温度
func (s *AlertService) monitorTemperature() {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			temps, err := s.getTemperatureInfo()
			if err != nil {
				log.Printf("Failed to get temperature info: %v", err)
				continue
			}

			for _, temp := range temps {
				if temp.Temperature > 0 {
					log.Printf("%s temperature: %.1f°C", temp.Label, temp.Temperature)
				}
			}
		}
	}
}

// monitorNetwork 监控网络
func (s *AlertService) monitorNetwork() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// 用于计算网络速率
	var lastBytesSent, lastBytesRecv uint64
	var lastTime time.Time

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			stats, err := s.getNetworkStats()
			if err != nil {
				log.Printf("Failed to get network stats: %v", err)
				continue
			}

			if !lastTime.IsZero() {
				timeDiff := time.Since(lastTime).Seconds()
				if timeDiff > 0 {
					sentRate := float64(stats.BytesSent-lastBytesSent) / timeDiff
					recvRate := float64(stats.BytesRecv-lastBytesRecv) / timeDiff

					log.Printf("Network: ↓ %.1f KB/s, ↑ %.1f KB/s",
						recvRate/1024, sentRate/1024)
				}
			}

			lastBytesSent = stats.BytesSent
			lastBytesRecv = stats.BytesRecv
			lastTime = time.Now()
		}
	}
}

// getTemperatureInfo 获取温度信息
func (s *AlertService) getTemperatureInfo() ([]TemperatureMetrics, error) {
	// 这里需要实现温度检测逻辑
	// 可以使用 psutil 或其他工具
	// 简化实现
	return []TemperatureMetrics{}, nil
}

// getNetworkStats 获取网络统计
func (s *AlertService) getNetworkStats() (NetworkMetrics, error) {
	// 这里需要实现网络统计获取逻辑
	// 简化实现
	return NetworkMetrics{}, nil
}

// GetSystemMetrics 获取系统指标的API方法
func (s *AlertService) GetSystemMetrics() (*SystemMetrics, error) {
	return s.getSystemMetrics()
}

// GetAlertRules 获取告警规则
func (s *AlertService) GetAlertRules() ([]models.AlertRule, error) {
	var rules []models.AlertRule
	err := s.db.Find(&rules).Error
	return rules, err
}

// CreateAlertRule 创建告警规则
func (s *AlertService) CreateAlertRule(rule models.AlertRule) error {
	return s.db.Create(&rule).Error
}

// UpdateAlertRule 更新告警规则
func (s *AlertService) UpdateAlertRule(rule models.AlertRule) error {
	return s.db.Save(&rule).Error
}

// DeleteAlertRule 删除告警规则
func (s *AlertService) DeleteAlertRule(id uint) error {
	return s.db.Delete(&models.AlertRule{}, id).Error
}

// GetActiveAlerts 获取活动告警
func (s *AlertService) GetActiveAlerts() ([]models.Alert, error) {
	var alerts []models.Alert
	err := s.db.Where("status = ?", "active").Order("triggered DESC").Find(&alerts).Error
	return alerts, err
}

// AcknowledgeAlert 确认告警
func (s *AlertService) AcknowledgeAlert(id uint, userID uint) error {
	return s.db.Model(&models.Alert{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     "acknowledged",
		"acknowledged_by": userID,
		"acknowledged_at": time.Now(),
	}).Error
}

// Stop 停止告警服务
func (s *AlertService) Stop() {
	s.cancel()

	s.mu.Lock()
	defer s.mu.Unlock()

	// 关闭所有告警通道
	for _, channel := range s.alertChannels {
		close(channel)
	}

	log.Println("Alert service stopped")
}