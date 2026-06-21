package system

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// DiskInfo 磁盘信息
type DiskInfo struct {
	Disks     []DiskPartition `json:"disks"`
	Timestamp int64          `json:"timestamp"`
}

// DiskPartition 磁盘分区信息
type DiskPartition struct {
	Device      string  `json:"device"`       // 设备名
	Mountpoint  string  `json:"mountpoint"`   // 挂载点
	Total       uint64  `json:"total"`        // 总容量 (bytes)
	Free        uint64  `json:"free"`         // 可用空间 (bytes)
	Used        uint64  `json:"used"`         // 已用空间 (bytes)
	UsedPercent float64 `json:"usedPercent"`  // 使用率百分比 (0-100)
	Fstype      string  `json:"fstype"`       // 文件系统类型
	Opts        []string `json:"opts"`        // 挂载选项

	// IO 统计
	ReadSpeed    float64 `json:"readSpeed"`    // 读取速度 (bytes/s)
	WriteSpeed   float64 `json:"writeSpeed"`   // 写入速度 (bytes/s)
	ReadCount    uint64  `json:"readCount"`    // 读取次数
	WriteCount   uint64  `json:"writeCount"`   // 写入次数
	ReadBytes    uint64  `json:"readBytes"`    // 累计读取字节数
	WriteBytes   uint64  `json:"writeBytes"`   // 累计写入字节数
	IoTime       uint64  `json:"ioTime"`       // IO 时间 (ms)

	// 衍生指标
	IOPS    float64 `json:"iops"`    // 每秒 IO 操作数
	Latency float64 `json:"latency"` // 平均 IO 延迟 (ms)
}

// GetDiskInfo 获取磁盘信息
func GetDiskInfo() (*DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	disks := make([]DiskPartition, 0, len(partitions))
	deviceMap := make(map[string]string) // 映射设备名到挂载点

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			// 跳过无法获取使用信息的分区
			continue
		}

		// 跳过特殊的文件系统
		if shouldSkipFS(usage.Fstype) {
			continue
		}

		// 提取基础设备名（去掉分区号）
		baseDevice := getBaseDevice(partition.Device)
		deviceMap[baseDevice] = partition.Mountpoint

		disks = append(disks, DiskPartition{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Total:       usage.Total,
			Free:        usage.Free,
			Used:        usage.Used,
			UsedPercent: usage.UsedPercent,
			Fstype:      usage.Fstype,
			Opts:        partition.Opts,
		})
	}

	// 获取磁盘 IO 统计
	if err := updateDiskIOStats(disks, deviceMap); err != nil {
		// IO 统计失败不影响整体结果
		fmt.Printf("Warning: failed to update disk IO stats: %v\n", err)
	}

	return &DiskInfo{
		Disks:     disks,
		Timestamp: time.Now().Unix(),
	}, nil
}

// updateDiskIOStats 更新磁盘 IO 统计（计算速度）
func updateDiskIOStats(disks []DiskPartition, deviceMap map[string]string) error {
	state := GetStatsState()
	counters, err := disk.IOCounters()
	if err != nil {
		return err
	}

	// 获取上次的状态
	lastStats, lastTS := state.GetLastDiskIO()
	now := time.Now()
	duration := now.Sub(lastTS).Seconds()

	// 构建当前状态
	currentStats := make(map[string]DiskIOHistory)

	// 更新每个磁盘的 IO 信息
	for i := range disks {
		baseDevice := getBaseDevice(disks[i].Device)
		// 完整设备名（去掉 /dev/ 前缀），用于精确匹配分区
		fullDevice := filepath.Base(disks[i].Device)

		// 查找对应的 IO 计数器：优先精确匹配分区名，再回退到基础设备名
		// 这样既能正确读取 NVMe 分区（nvme0n1p2），也能读取传统分区（sda1 -> sda）
		var counter *disk.IOCountersStat
		counterKey := ""
		for _, c := range counters {
			if c.Name == fullDevice {
				counter = &c
				counterKey = c.Name
				break
			}
		}
		if counter == nil {
			for _, c := range counters {
				if c.Name == baseDevice {
					counter = &c
					counterKey = c.Name
					break
				}
			}
		}

		if counter == nil {
			continue
		}

		// 保存当前统计
		currentStats[counterKey] = DiskIOHistory{
			ReadBytes:  counter.ReadBytes,
			WriteBytes: counter.WriteBytes,
			ReadCount:  counter.ReadCount,
			WriteCount: counter.WriteCount,
			IoTime:     counter.IoTime,
		}

		disks[i].ReadBytes = counter.ReadBytes
		disks[i].WriteBytes = counter.WriteBytes
		disks[i].ReadCount = counter.ReadCount
		disks[i].WriteCount = counter.WriteCount
		disks[i].IoTime = counter.IoTime

		// 计算速度（需要上次的数据）
		if lastStat, ok := lastStats[counterKey]; ok && duration > 0 {
			readSpeed, writeSpeed := CalculateDiskSpeed(currentStats[counterKey], lastStat, duration)
			disks[i].ReadSpeed = readSpeed
			disks[i].WriteSpeed = writeSpeed

			// 计算 IOPS（每秒读写操作数之和）
			opsDelta := float64((currentStats[counterKey].ReadCount + currentStats[counterKey].WriteCount) -
				(lastStat.ReadCount + lastStat.WriteCount))
			disks[i].IOPS = opsDelta / duration

			// 计算平均 IO 延迟 (ms)
			ioTimeDelta := float64(currentStats[counterKey].IoTime - lastStat.IoTime)
			if opsDelta > 0 {
				disks[i].Latency = ioTimeDelta / opsDelta
			}
		}
	}

	// 更新全局状态
	state.UpdateDiskIO(currentStats)

	return nil
}

// getBaseDevice 获取基础设备名（去掉分区号）
// 处理两种命名规则:
//   - SCSI/SATA: /dev/sda1 -> sda (去掉末尾数字)
//   - NVMe:      /dev/nvme0n1p2 -> nvme0n1 (去掉 "p" + 数字后缀)
func getBaseDevice(device string) string {
	base := filepath.Base(device)

	// NVMe 设备: 形如 nvme0n1p2，分区号为 "p" 后的数字
	if strings.HasPrefix(base, "nvme") {
		if idx := strings.LastIndex(base, "p"); idx > 0 {
			suffix := base[idx+1:]
			if suffix != "" && isAllDigits(suffix) {
				return base[:idx]
			}
		}
		return base
	}

	// 通用: 去掉末尾数字（处理 sda1 -> sda, mmcblk0p1 -> mmcblk0）
	for len(base) > 0 && base[len(base)-1] >= '0' && base[len(base)-1] <= '9' {
		base = base[:len(base)-1]
	}
	return base
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// shouldSkipFS 判断是否应该跳过该文件系统
func shouldSkipFS(fstype string) bool {
	skipFS := map[string]bool{
		"proc":        true,
		"sysfs":       true,
		"devtmpfs":    true,
		"tmpfs":       true,
		"cgroup":      true,
		"configfs":    true,
		"debugfs":     true,
		"securityfs":  true,
		"pstore":      true,
		"devpts":      true,
		"fusectl":     true,
		"none":        true,
		"overlay":     true,
		"autofs":      true,
		"binfmt_misc": true,
		"tracefs":     true,
		"ramfs":       true,
		"hugetlbfs":   true,
		"mqueue":      true,
		"efivarfs":    true,
		"nsfs":        true,
	}
	return skipFS[fstype]
}

// GetDiskPartitions 获取所有磁盘分区
func GetDiskPartitions() ([]DiskPartition, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk partitions: %w", err)
	}

	result := make([]DiskPartition, 0, len(partitions))
	for _, p := range partitions {
		result = append(result, DiskPartition{
			Device:     p.Device,
			Mountpoint: p.Mountpoint,
			Fstype:     p.Fstype,
			Opts:       p.Opts,
		})
	}

	return result, nil
}

// GetDiskUsage 获取指定挂载点的磁盘使用情况
func GetDiskUsage(mountpoint string) (*DiskUsage, error) {
	usage, err := disk.Usage(mountpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get disk usage for %s: %w", mountpoint, err)
	}

	return &DiskUsage{
		Path:              usage.Path,
		Total:             usage.Total,
		Free:              usage.Free,
		Used:              usage.Used,
		UsedPercent:       usage.UsedPercent,
		InodesTotal:       usage.InodesTotal,
		InodesUsed:        usage.InodesUsed,
		InodesFree:        usage.InodesFree,
		InodesUsedPercent: usage.InodesUsedPercent,
	}, nil
}

// DiskUsage 磁盘使用情况
type DiskUsage struct {
	Path              string  `json:"path"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

// GetDiskIOCounters 获取磁盘 IO 计数器
func GetDiskIOCounters() (map[string]DiskIOHistory, error) {
	counters, err := disk.IOCounters()
	if err != nil {
		return nil, fmt.Errorf("failed to get disk IO counters: %w", err)
	}

	result := make(map[string]DiskIOHistory)
	for _, c := range counters {
		result[c.Name] = DiskIOHistory{
			ReadBytes:  c.ReadBytes,
			WriteBytes: c.WriteBytes,
			ReadCount:  c.ReadCount,
			WriteCount: c.WriteCount,
			IoTime:     c.IoTime,
		}
	}

	return result, nil
}

// FormatBytes 格式化字节数
func FormatDiskBytes(bytes uint64) string {
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

// FindDiskByMountpoint 根据挂载点查找磁盘
func FindDiskByMountpoint(disks []DiskPartition, mountpoint string) *DiskPartition {
	for i := range disks {
		if disks[i].Mountpoint == mountpoint {
			return &disks[i]
		}
	}
	return nil
}

// FindDiskByDevice 根据设备名查找磁盘
func FindDiskByDevice(disks []DiskPartition, device string) *DiskPartition {
	for i := range disks {
		if strings.Contains(disks[i].Device, device) || strings.Contains(device, disks[i].Device) {
			return &disks[i]
		}
	}
	return nil
}

// GetImportantMountpoints 获取重要的挂载点
func GetImportantMountpoints() []string {
	return []string{
		"/",
		"/home",
		"/var",
		"/tmp",
		"/boot",
		"/boot/efi",
		"/data",
	}
}
