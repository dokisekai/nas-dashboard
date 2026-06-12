package system

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v3/mem"
)

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total       uint64  `json:"total"`       // 总内存 (bytes)
	Available   uint64  `json:"available"`   // 可用内存 (bytes)
	Used        uint64  `json:"used"`        // 已用内存 (bytes)
	Free        uint64  `json:"free"`        // 空闲内存 (bytes)
	Percent     float64 `json:"percent"`     // 使用率百分比 (0-100)
	Active      uint64  `json:"active"`      // 活跃内存
	Inactive    uint64  `json:"inactive"`    // 非活跃内存
	Buffers     uint64  `json:"buffers"`     // 缓冲区内存
	Cached      uint64  `json:"cached"`      // 缓存内存
	Shared      uint64  `json:"shared"`      // 共享内存
	Slab        uint64  `json:"slab"`        // Slab 内存
	SwapTotal   uint64  `json:"swapTotal"`   // Swap 总量
	SwapUsed    uint64  `json:"swapUsed"`    // Swap 已用
	SwapFree    uint64  `json:"swapFree"`    // Swap 空闲
	SwapPercent float64 `json:"swapPercent"` // Swap 使用率
	Timestamp   int64   `json:"timestamp"`   // 时间戳
}

// GetMemoryInfo 获取内存信息
func GetMemoryInfo() (*MemoryInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get swap memory: %w", err)
	}

	return &MemoryInfo{
		Total:       vmStat.Total,
		Available:   vmStat.Available,
		Used:        vmStat.Used,
		Free:        vmStat.Free,
		Percent:     vmStat.UsedPercent,
		Active:      vmStat.Active,
		Inactive:    vmStat.Inactive,
		Buffers:     vmStat.Buffers,
		Cached:      vmStat.Cached,
		Shared:      vmStat.Shared,
		Slab:        vmStat.Slab,
		SwapTotal:   swapStat.Total,
		SwapUsed:    swapStat.Used,
		SwapFree:    swapStat.Total - swapStat.Used,
		SwapPercent: swapStat.UsedPercent,
		Timestamp:   time.Now().Unix(),
	}, nil
}

// GetSwapInfo 仅获取 Swap 信息
func GetSwapInfo() (*SwapInfo, error) {
	swapStat, err := mem.SwapMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get swap memory: %w", err)
	}

	return &SwapInfo{
		Total:       swapStat.Total,
		Used:        swapStat.Used,
		Free:        swapStat.Total - swapStat.Used,
		UsedPercent: swapStat.UsedPercent,
		Sin:         swapStat.Sin,
		Sout:        swapStat.Sout,
	}, nil
}

// SwapInfo Swap 信息
type SwapInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
	Sin         uint64  `json:"sin"`  // Swap in
	Sout        uint64  `json:"sout"` // Swap out
}

// GetMemoryUsage 获取内存使用情况（简化版）
func GetMemoryUsage() (float64, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return 0, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	return vmStat.UsedPercent, nil
}

// GetHostMemoryInfo 获取主机内存信息（更详细）
func GetHostMemoryInfo() (*HostMemoryInfo, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	return &HostMemoryInfo{
		Total:        vmStat.Total,
		Free:         vmStat.Free,
		Used:         vmStat.Used,
		UsedPercent:  vmStat.UsedPercent,
		Available:    vmStat.Available,
		Active:       vmStat.Active,
		Inactive:     vmStat.Inactive,
		Wired:        vmStat.Wired,
		Laundry:      vmStat.Laundry,
		Buffers:      vmStat.Buffers,
		Cached:       vmStat.Cached,
		Shared:       vmStat.Shared,
		Slab:         vmStat.Slab,
		PageTables:   vmStat.PageTables,
		SwapCached:   vmStat.SwapCached,
		CommitLimit:  vmStat.CommitLimit,
		CommittedAS:  vmStat.CommittedAS,
		HighTotal:    vmStat.HighTotal,
		HighFree:     vmStat.HighFree,
		LowTotal:     vmStat.LowTotal,
		LowFree:      vmStat.LowFree,
		SwapTotal:    vmStat.SwapTotal,
		SwapFree:     vmStat.SwapFree,
		Mapped:       vmStat.Mapped,
		VmallocTotal: vmStat.VmallocTotal,
		VmallocUsed:  vmStat.VmallocUsed,
		VmallocChunk: vmStat.VmallocChunk,
		HugePagesTotal:    vmStat.HugePagesTotal,
		HugePagesFree:     vmStat.HugePagesFree,
		HugePageSize:      vmStat.HugePageSize,
	}, nil
}

// HostMemoryInfo 主机内存详细信息
type HostMemoryInfo struct {
	Total         uint64  `json:"total"`
	Free          uint64  `json:"free"`
	Used          uint64  `json:"used"`
	UsedPercent   float64 `json:"usedPercent"`
	Available     uint64  `json:"available"`
	Active        uint64  `json:"active"`
	Inactive      uint64  `json:"inactive"`
	Wired         uint64  `json:"wired"`
	Laundry       uint64  `json:"laundry"`
	Buffers       uint64  `json:"buffers"`
	Cached        uint64  `json:"cached"`
	Shared        uint64  `json:"shared"`
	Slab          uint64  `json:"slab"`
	PageTables    uint64  `json:"pageTables"`
	SwapCached    uint64  `json:"swapCached"`
	CommitLimit   uint64  `json:"commitLimit"`
	CommittedAS   uint64  `json:"committedAS"`
	HighTotal     uint64  `json:"highTotal"`
	HighFree      uint64  `json:"highFree"`
	LowTotal      uint64  `json:"lowTotal"`
	LowFree       uint64  `json:"lowFree"`
	SwapTotal     uint64  `json:"swapTotal"`
	SwapFree      uint64  `json:"swapFree"`
	Mapped        uint64  `json:"mapped"`
	VmallocTotal  uint64  `json:"vmallocTotal"`
	VmallocUsed   uint64  `json:"vmallocUsed"`
	VmallocChunk  uint64  `json:"vmallocChunk"`
	HugePagesTotal   uint64 `json:"hugePagesTotal"`
	HugePagesFree    uint64 `json:"hugePagesFree"`
	HugePageSize     uint64 `json:"hugePageSize"`
}

// FormatMemoryBytes 格式化内存字节数为人类可读格式
func FormatMemoryBytes(bytes uint64) string {
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
