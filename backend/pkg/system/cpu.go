package system

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/load"
)

// CPUInfo CPU 信息
type CPUInfo struct {
	Usage     float64   `json:"usage"`      // CPU 使用率 (0-1)
	Cores     int       `json:"cores"`      // CPU 核心数
	Model     string    `json:"model"`      // CPU 型号
	Mhz       float64   `json:"mhz"`        // CPU 频率
	Load1     float64   `json:"load1"`      // 1分钟负载
	Load5     float64   `json:"load5"`      // 5分钟负载
	Load15    float64   `json:"load15"`     // 15分钟负载
	PerCore   []float64 `json:"perCore"`    // 每个核心使用率
	Timestamp int64     `json:"timestamp"`  // 时间戳
}

// CPUDetail CPU 详细信息
type CPUDetail struct {
	CPU        int32   `json:"cpu"`         // CPU 编号
	VendorID   string  `json:"vendorId"`    // 供应商 ID
	Family     string  `json:"family"`      // 系列
	Model      string  `json:"model"`       // 型号
	Stepping   int32   `json:"stepping"`    // 步进
	PhysicalID string  `json:"physicalId"`  // 物理 ID
	CoreID     string  `json:"coreId"`      // 核心 ID
	Models     string  `json:"models"`      // 型号名称
	ModelName  string  `json:"modelName"`   // 型号名称
	Mhz        float64 `json:"mhz"`         // 频率
	CacheSize  int32   `json:"cacheSize"`   // 缓存大小
	Flags      []string `json:"flags"`      // 标志
	Microcode   string  `json:"microcode"`   // 微代码
}

// GetCPUInfo 获取 CPU 信息
func GetCPUInfo() (*CPUInfo, error) {
	// 获取 CPU 核心数
	cores, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU counts: %w", err)
	}

	// 获取 CPU 使用率 (需要一点时间来计算准确值)
	percent, err := cpu.Percent(100*time.Millisecond, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	// 获取每个核心使用率
	perCore, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get per-core CPU percent: %w", err)
	}

	// 获取 CPU 详细信息
	infos, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU info: %w", err)
	}

	// 获取系统负载
	loadStat, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get load average: %w", err)
	}

	var usage float64
	var mhz float64
	var modelName string

	if len(percent) > 0 {
		usage = percent[0] / 100
	}

	if len(infos) > 0 {
		mhz = infos[0].Mhz
		modelName = infos[0].ModelName
	}

	return &CPUInfo{
		Usage:     usage,
		Cores:     cores,
		Model:     modelName,
		Mhz:       mhz,
		Load1:     loadStat.Load1,
		Load5:     loadStat.Load5,
		Load15:    loadStat.Load15,
		PerCore:   normalizePerCore(perCore),
		Timestamp: time.Now().Unix(),
	}, nil
}

// normalizePerCore 将百分比转换为 0-1 范围
func normalizePerCore(percent []float64) []float64 {
	result := make([]float64, len(percent))
	for i, p := range percent {
		result[i] = p / 100
	}
	return result
}

// GetCPUDetails 获取 CPU 详细信息
func GetCPUDetails() ([]CPUDetail, error) {
	infos, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU details: %w", err)
	}

	details := make([]CPUDetail, len(infos))
	for i, info := range infos {
		details[i] = CPUDetail{
			CPU:        info.CPU,
			VendorID:   info.VendorID,
			Family:     info.Family,
			Model:      info.Model,
			Stepping:   info.Stepping,
			PhysicalID: info.PhysicalID,
			CoreID:     info.CoreID,
			Models:     info.Model,
			ModelName:  info.ModelName,
			Mhz:        info.Mhz,
			CacheSize:  info.CacheSize,
			Flags:      info.Flags,
			Microcode:  info.Microcode,
		}
	}

	return details, nil
}

// GetCPUPercent 获取 CPU 使用率百分比
func GetCPUPercent(perCPU bool) ([]float64, error) {
	percent, err := cpu.Percent(100*time.Millisecond, perCPU)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	// 转换为 0-1 范围
	result := make([]float64, len(percent))
	for i, p := range percent {
		result[i] = p / 100
	}

	return result, nil
}

// GetLoadAverage 获取系统负载
func GetLoadAverage() (*LoadAvg, error) {
	loadStat, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("failed to get load average: %w", err)
	}

	return &LoadAvg{
		Load1:  loadStat.Load1,
		Load5:  loadStat.Load5,
		Load15: loadStat.Load15,
	}, nil
}

// LoadAvg 系统负载
type LoadAvg struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// GetCPUCounts 获取 CPU 数量信息
func GetCPUCounts() (*CPUCounts, error) {
	logical, err := cpu.Counts(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get logical CPU counts: %w", err)
	}

	physical, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get physical CPU counts: %w", err)
	}

	return &CPUCounts{
		Logical:  logical,
		Physical: physical,
	}, nil
}

// CPUCounts CPU 数量
type CPUCounts struct {
	Logical  int `json:"logical"`
	Physical int `json:"physical"`
}

// GetCPUTime 获取 CPU 时间统计
func GetCPUTime() ([]CPUTimeStat, error) {
	times, err := cpu.Times(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU times: %w", err)
	}

	stats := make([]CPUTimeStat, len(times))
	for i, t := range times {
		stats[i] = CPUTimeStat{
			CPU:       t.CPU,
			User:      t.User,
			System:    t.System,
			Idle:      t.Idle,
			Nice:      t.Nice,
			Iowait:    t.Iowait,
			Irq:       t.Irq,
			Softirq:   t.Softirq,
			Steal:     t.Steal,
			Guest:     t.Guest,
			GuestNice: t.GuestNice,
		}
	}

	return stats, nil
}

// CPUTimeStat CPU 时间统计
type CPUTimeStat struct {
	CPU       string  `json:"cpu"`
	User      float64 `json:"user"`
	System    float64 `json:"system"`
	Idle      float64 `json:"idle"`
	Nice      float64 `json:"nice"`
	Iowait    float64 `json:"iowait"`
	Irq       float64 `json:"irq"`
	Softirq   float64 `json:"softirq"`
	Steal     float64 `json:"steal"`
	Guest     float64 `json:"guest"`
	GuestNice float64 `json:"guestNice"`
}
