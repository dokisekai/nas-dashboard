package power

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PowerMonitor 功耗监控器
type PowerMonitor struct {
	dataDir     string
	monitorPath string
}

// NewPowerMonitor 创建功耗监控器
func NewPowerMonitor() *PowerMonitor {
	return &PowerMonitor{
		dataDir:     "/var/lib/nas-dashboard/power",
		monitorPath: "/data/nas-dashboard/modules/power-monitor/scripts/get_power_v2.sh",
	}
}

// PowerData 功耗数据结构
type PowerData struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUPackage  float64   `json:"cpuPackage"`
	CPUCore     float64   `json:"cpuCore"`
	CPUUncore   float64   `json:"cpuUncore"`
	IGPU        float64   `json:"igpu"`
	DGPU        float64   `json:"dgpu"`
	HDD         float64   `json:"hdd"`
	SSD         float64   `json:"ssd"`
	MBRAM       float64   `json:"mbram"`
	Cooling     float64   `json:"cooling"`
	USB         float64   `json:"usb"`
	PowerLoss   float64   `json:"powerLoss"`
	Total       float64   `json:"total"`
}

// GetCurrentPower 获取当前功耗数据
func (pm *PowerMonitor) GetCurrentPower() (*PowerData, error) {
	return pm.getPowerFromScript()
}

// getPowerFromScript 优先从 sysfs 直接读取，失败时回退到外部脚本
func (pm *PowerMonitor) getPowerFromScript() (*PowerData, error) {
	// 优先使用内置 sysfs 读取（无外部依赖，毫秒级响应）
	if pd, err := readPowerFromSysfs(); err == nil && pd != nil {
		return pd, nil
	}
	// 回退：调用外部脚本（如果存在）
	if _, err := os.Stat(pm.monitorPath); err == nil {
		cmd := exec.Command("sudo", pm.monitorPath)
		output, err := cmd.CombinedOutput()
		if err == nil {
			if pd, perr := pm.parsePowerOutput(string(output)); perr == nil && pd != nil {
				return pd, nil
			}
		}
	}
	return nil, fmt.Errorf("no power source available: sysfs unreadable and script missing at %s", pm.monitorPath)
}

// powerState RAPL 能量状态（用于计算瞬时功率）
type powerState struct {
	mu        sync.Mutex
	last      map[string]energySample
}
type energySample struct {
	energyUJ uint64
	ts       time.Time
}

var globalPowerState = &powerState{last: make(map[string]energySample)}

// readEnergyFile 读取 energy_uj 文件（单位 microjoules）
// RAPL energy_uj 通常仅 root 可读，失败时回退到 sudo cat
func readEnergyFile(path string) (uint64, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		return strconv.ParseUint(strings.TrimSpace(string(data)), 10, 64)
	}
	// 权限不足时使用 sudo（要求 sudoers 配置 NOPASSWD）
	cmd := exec.Command("sudo", "-n", "cat", path)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(strings.TrimSpace(string(out)), 10, 64)
}

// calcPower 计算两次能量采样之间的瞬时功率（瓦特）
// 第一次采样返回 0（无前值）
func calcPower(key string, currentUJ uint64) float64 {
	now := time.Now()
	globalPowerState.mu.Lock()
	prev, hasPrev := globalPowerState.last[key]
	globalPowerState.last[key] = energySample{energyUJ: currentUJ, ts: now}
	globalPowerState.mu.Unlock()

	if !hasPrev {
		return 0
	}
	duration := now.Sub(prev.ts).Seconds()
	if duration <= 0 {
		return 0
	}
	// 处理能量计数器溢出（wrap around）
	var delta uint64
	if currentUJ >= prev.energyUJ {
		delta = currentUJ - prev.energyUJ
	} else {
		// 计数器回绕，使用最大范围估算（典型 64-bit，但实际 ~53 bits 可用）
		delta = (1 << 53) + currentUJ - prev.energyUJ
	}
	return float64(delta) / 1_000_000.0 / duration
}

// readRAPL 找到 intel-rapl 目录下所有域的 energy_uj，按子域名称归类
// 返回: package 功率、core 功率、uncore 功率
func readRAPL() (pkgW, coreW, uncoreW float64) {
	baseDirs := []string{
		"/sys/devices/virtual/powercap/intel-rapl",
		"/sys/class/powercap/intel-rapl",
	}
	var raplRoot string
	for _, d := range baseDirs {
		if entries, err := os.ReadDir(d); err == nil && len(entries) > 0 {
			raplRoot = d
			break
		}
	}
	if raplRoot == "" {
		return 0, 0, 0
	}

	// 顶层 intel-rapl:0 通常是 package
	packageEntries, err := os.ReadDir(raplRoot)
	if err != nil {
		return 0, 0, 0
	}
	for _, entry := range packageEntries {
		name := entry.Name()
		if !strings.HasPrefix(name, "intel-rapl:") || !entry.IsDir() {
			continue
		}
		// 读取 name 文件，确认是 package
		nameFile := filepath.Join(raplRoot, name, "name")
		nameData, err := os.ReadFile(nameFile)
		if err != nil {
			continue
		}
		domainName := strings.TrimSpace(string(nameData))
		energyFile := filepath.Join(raplRoot, name, "energy_uj")
		if energy, err := readEnergyFile(energyFile); err == nil {
			p := calcPower("rapl:"+name, energy)
			if p > 0 && p < 500 { // 合理的 package 功率范围
				pkgW = p
			}
		}
		_ = domainName

		// 遍历子域 intel-rapl:0:0 / intel-rapl:0:1 等
		subEntries, err := os.ReadDir(filepath.Join(raplRoot, name))
		if err != nil {
			continue
		}
		for _, sub := range subEntries {
			if !strings.HasPrefix(sub.Name(), "intel-rapl:") || !sub.IsDir() {
				continue
			}
			subNameFile := filepath.Join(raplRoot, name, sub.Name(), "name")
			subNameData, err := os.ReadFile(subNameFile)
			if err != nil {
				continue
			}
			subDomain := strings.TrimSpace(string(subNameData))
			subEnergyFile := filepath.Join(raplRoot, name, sub.Name(), "energy_uj")
			energy, err := readEnergyFile(subEnergyFile)
			if err != nil {
				continue
			}
			p := calcPower("rapl:"+name+":"+sub.Name(), energy)
			if p <= 0 || p > 500 {
				continue
			}
			switch {
			case strings.Contains(strings.ToLower(subDomain), "core"):
				coreW += p
			case strings.Contains(strings.ToLower(subDomain), "uncore"):
				uncoreW += p
			}
		}
		break
	}
	return pkgW, coreW, uncoreW
}

// readAMDGPUPower 读取 AMD 独显功耗（microwatts -> watts）
func readAMDGPUPower() float64 {
	// /sys/class/drm/card*/device/hwmon/hwmon*/power1_average
	pattern := "/sys/class/drm/card*/device/hwmon/hwmon*/power1_average"
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return 0
	}
	var total float64
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			continue
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			continue
		}
		// amdgpu power1_average 单位为微瓦
		total += v / 1_000_000.0
	}
	return total
}

// readNVMePower 估算 NVMe 功耗（基于温度和活跃度，简化版）
// 由于 sysfs 不直接暴露 NVMe 功耗，使用粗略估算
func readNVMePower() float64 {
	pattern := "/sys/class/nvme/nvme*/device/hwmon*/power"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return 0
	}
	// 大多数 NVMe 待机 0.05W，活动 3-8W，无精确数据时按 1.5W/块估算
	return float64(len(matches)) * 1.5
}

// readHDDPower 估算 HDD 功耗（通过 hdparm 检测 standby/active 状态）
func readHDDPower() float64 {
	// 找出所有块设备（sd[a-z]、hd[a-z]）
	devPattern := "/dev/sd[a-z]"
	devs, err := filepath.Glob(devPattern)
	if err != nil {
		devs = nil
	}
	devPattern2 := "/dev/hd[a-z]"
	if extras, err := filepath.Glob(devPattern2); err == nil {
		devs = append(devs, extras...)
	}
	if len(devs) == 0 {
		return 0
	}

	var total float64
	for _, dev := range devs {
		// hdparm -C 需要 root，我们尝试直接调用，失败则按 active 估算
		cmd := exec.Command("hdparm", "-C", dev)
		output, err := cmd.CombinedOutput()
		if err != nil {
			total += 8.0 // 未知状态按 active 估算
			continue
		}
		out := string(output)
		switch {
		case strings.Contains(out, "standby"):
			total += 0.9
		case strings.Contains(out, "sleep"):
			total += 0.5
		case strings.Contains(out, "active/idle"):
			total += 8.0
		default:
			total += 8.0
		}
	}
	return total
}

// readSSDPower 估算 SATA SSD 功耗（无精确数据，按块数估算）
func readSSDPower() float64 {
	// 区分 NVMe（在 readNVMePower 中处理）和 SATA SSD
	// 简化：检测 scsi generic 设备，无法精确区分 SSD/HDD
	// 给一个保守默认值
	return 4.0
}

// readPowerFromSysfs 直接从 sysfs 读取真实功耗数据
// 在外部脚本不可用时使用，结果与脚本语义保持一致
func readPowerFromSysfs() (*PowerData, error) {
	pkgW, coreW, uncoreW := readRAPL()
	dgpuW := readAMDGPUPower()
	hddW := readHDDPower()
	nvmeW := readNVMePower()
	ssdW := readSSDPower()

	// Intel 核显：包含在 uncore 内，约占 uncore 60-80%
	igpuW := uncoreW * 0.7

	// 固定估算：主板+内存、散热、USB
	mbramW := 21.0
	coolingW := 12.0
	usbW := 2.0

	// 原始合计（不含电源损耗）
	rawTotal := pkgW + dgpuW + hddW + nvmeW + ssdW + mbramW + coolingW + usbW
	// 电源转换损耗（典型 80% 效率 -> 18% 损耗）
	lossW := rawTotal * 0.18

	totalW := rawTotal + lossW

	// 至少 CPU package 或 DGPU 应有非零值，否则视为读取失败
	if pkgW == 0 && dgpuW == 0 && coreW == 0 {
		return nil, fmt.Errorf("sysfs power sources unavailable")
	}

	return &PowerData{
		Timestamp:   time.Now(),
		CPUPackage:  round2(pkgW),
		CPUCore:     round2(coreW),
		CPUUncore:   round2(uncoreW),
		IGPU:        round2(igpuW),
		DGPU:        round2(dgpuW),
		HDD:         round2(hddW),
		SSD:         round2(nvmeW + ssdW),
		MBRAM:       round2(mbramW),
		Cooling:     round2(coolingW),
		USB:         round2(usbW),
		PowerLoss:   round2(lossW),
		Total:       round2(totalW),
	}, nil
}

func round2(v float64) float64 {
	return mathRound(v*100) / 100
}

// 使用单独函数避免引入 math 包
func mathRound(v float64) float64 {
	if v < 0 {
		return -mathRound(-v)
	}
	if v-float64(int(v)) >= 0.5 {
		return float64(int(v) + 1)
	}
	return float64(int(v))
}

// parsePowerOutput 解析脚本输出（保留向后兼容）
func (pm *PowerMonitor) parsePowerOutput(output string) (*PowerData, error) {
	pd := &PowerData{
		Timestamp: time.Now(),
	}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		extractPowerValue := func(text string) float64 {
			re := regexpMustCompile(`([0-9]+\.[0-9]+)\s*W`)
			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
					return val
				}
			}
			reInt := regexpMustCompile(`([0-9]+)\s*W`)
			matchesInt := reInt.FindStringSubmatch(text)
			if len(matchesInt) > 1 {
				if val, err := strconv.ParseFloat(matchesInt[1], 64); err == nil {
					return val
				}
			}
			return 0
		}

		if strings.Contains(line, "整机当前真实预估总功耗") {
			pd.Total = extractPowerValue(line)
		}
		if strings.Contains(line, "CPU Package") && strings.Contains(line, "13700K") {
			pd.CPUPackage = extractPowerValue(line)
		}
		if strings.Contains(line, "CPU Core") && strings.Contains(line, "16C/24T") {
			pd.CPUCore = extractPowerValue(line)
		}
		if strings.Contains(line, "CPU Uncore") && strings.Contains(line, "含核显") {
			pd.CPUUncore = extractPowerValue(line)
		}
		if strings.Contains(line, "Intel 核显") && strings.Contains(line, "UHD 770") {
			pd.IGPU = extractPowerValue(line)
		}
		if strings.Contains(line, "AMD 独显") && strings.Contains(line, "6950 XT") {
			pd.DGPU = extractPowerValue(line)
		}
		if strings.Contains(line, "HDD 机械硬盘") {
			pd.HDD = extractPowerValue(line)
		}
		if strings.Contains(line, "SSD 固态硬盘") {
			pd.SSD = extractPowerValue(line)
		}
		if strings.Contains(line, "主板 & 2x32G内存") {
			pd.MBRAM = extractPowerValue(line)
		}
		if strings.Contains(line, "散热") && !strings.Contains(line, "1100W金牌损耗") {
			pd.Cooling = extractPowerValue(line)
		}
		if strings.Contains(line, "USB及外设功耗") {
			pd.USB = extractPowerValue(line)
		}
		if strings.Contains(line, "1100W金牌损耗") {
			pd.PowerLoss = extractPowerValue(line)
		}
	}

	return pd, nil
}

// DirectReadPower 直接读取功耗（优先 sysfs，回退脚本）
func DirectReadPower() (*PowerData, error) {
	pm := NewPowerMonitor()
	return pm.getPowerFromScript()
}

// GetHistoricalPower 获取历史功耗数据
func (pm *PowerMonitor) GetHistoricalPower(days int) ([]*PowerData, error) {
	var allData []*PowerData

	for i := 0; i <= days; i++ {
		date := time.Now().AddDate(0, 0, -i)
		csvFile := fmt.Sprintf("%s/power_%s.csv", pm.dataDir, date.Format("20060102"))

		data, err := pm.readCSVFile(csvFile)
		if err != nil {
			continue
		}
		allData = append(allData, data...)
	}

	return allData, nil
}

// GetPowerStatistics 获取功耗统计信息
func (pm *PowerMonitor) GetPowerStatistics(days int) (*map[string]interface{}, error) {
	data, err := pm.GetHistoricalPower(days)
	if err != nil {
		return &map[string]interface{}{
			"averagePower": 0.0,
			"maxPower":     0.0,
			"minPower":     0.0,
			"totalEnergy":  0.0,
			"sampleCount":  0,
		}, nil
	}

	if len(data) == 0 {
		return &map[string]interface{}{
			"averagePower": 0.0,
			"maxPower":     0.0,
			"minPower":     0.0,
			"totalEnergy":  0.0,
			"sampleCount":  0,
		}, nil
	}

	var sum, max, min float64
	max = data[0].Total
	min = data[0].Total

	for _, d := range data {
		sum += d.Total
		if d.Total > max {
			max = d.Total
		}
		if d.Total < min {
			min = d.Total
		}
	}

	avg := sum / float64(len(data))
	totalEnergy := (sum / 1000) * (float64(len(data)) * 5 / 60) // kWh

	return &map[string]interface{}{
		"averagePower": avg,
		"maxPower":     max,
		"minPower":     min,
		"totalEnergy":  totalEnergy,
		"sampleCount":  len(data),
	}, nil
}

// readCSVFile 读取 CSV 文件
func (pm *PowerMonitor) readCSVFile(filename string) ([]*PowerData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []*PowerData
	for i, record := range records {
		if i == 0 || len(record) < 12 {
			continue
		}

		pd := &PowerData{}
		pd.Timestamp, _ = time.Parse("2006-01-02 15:04:05", record[0])
		pd.CPUPackage, _ = strconv.ParseFloat(record[1], 64)
		pd.CPUCore, _ = strconv.ParseFloat(record[2], 64)
		pd.CPUUncore, _ = strconv.ParseFloat(record[3], 64)
		pd.IGPU, _ = strconv.ParseFloat(record[4], 64)
		pd.DGPU, _ = strconv.ParseFloat(record[5], 64)
		pd.HDD, _ = strconv.ParseFloat(record[6], 64)
		pd.SSD, _ = strconv.ParseFloat(record[7], 64)
		pd.MBRAM, _ = strconv.ParseFloat(record[8], 64)
		pd.Cooling, _ = strconv.ParseFloat(record[9], 64)
		pd.USB, _ = strconv.ParseFloat(record[10], 64)
		pd.Total, _ = strconv.ParseFloat(record[11], 64)

		data = append(data, pd)
	}

	return data, nil
}

// CheckPowerAlert 检查功耗告警
func CheckPowerAlert(powerData *PowerData, thresholds map[string]float64) []string {
	var alerts []string

	if powerData.Total >= thresholds["critical"] {
		alerts = append(alerts, fmt.Sprintf("严重告警：总功耗 %.2fW 超过临界值 %.2fW", powerData.Total, thresholds["critical"]))
	} else if powerData.Total >= thresholds["high"] {
		alerts = append(alerts, fmt.Sprintf("高功耗告警：总功耗 %.2fW 超过高值 %.2fW", powerData.Total, thresholds["high"]))
	}

	if powerData.DGPU >= thresholds["gpu"] {
		alerts = append(alerts, fmt.Sprintf("GPU 功耗告警：%.2fW", powerData.DGPU))
	}

	return alerts
}
