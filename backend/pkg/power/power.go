package power

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
		dataDir:     "/home/hserver/power_monitor/data",
		monitorPath: "/home/hserver/get_power_v2.sh",
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

// getPowerFromPowerFromScript 从shell脚本获取功耗数据
func (pm *PowerMonitor) getPowerFromScript() (*PowerData, error) {
	// 调用 shell 脚本获取功耗数据
	cmd := exec.Command("sudo", pm.monitorPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get power data: %v", err)
	}

	return pm.parsePowerOutput(string(output))
}

// parsePowerOutput 解析脚本输出
func (pm *PowerMonitor) parsePowerOutput(output string) (*PowerData, error) {
	pd := &PowerData{
		Timestamp: time.Now(),
	}

	// 解析脚本输出，提取各组件功耗
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 通用功耗值提取函数 - 专门匹配 "XX.XX W" 格式
		extractPowerValue := func(text string) float64 {
			// 匹配 "数字.数字 W" 格式，功耗值通常在最后且单位为W
			re := regexp.MustCompile(`([0-9]+\.[0-9]+)\s*W`)
			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				if val, err := strconv.ParseFloat(matches[1], 64); err == nil {
					return val
				}
			}
			// 如果没有小数点，尝试匹配整数 + W
			reInt := regexp.MustCompile(`([0-9]+)\s*W`)
			matchesInt := reInt.FindStringSubmatch(text)
			if len(matchesInt) > 1 {
				if val, err := strconv.ParseFloat(matchesInt[1], 64); err == nil {
					return val
				}
			}
			return 0
		}

		// 解析总功耗
		if strings.Contains(line, "整机当前真实预估总功耗") {
			pd.Total = extractPowerValue(line)
		}

		// 解析CPU Package
		if strings.Contains(line, "CPU Package") && strings.Contains(line, "13700K") {
			pd.CPUPackage = extractPowerValue(line)
		}

		// 解析CPU Core
		if strings.Contains(line, "CPU Core") && strings.Contains(line, "16C/24T") {
			pd.CPUCore = extractPowerValue(line)
		}

		// 解析CPU Uncore
		if strings.Contains(line, "CPU Uncore") && strings.Contains(line, "含核显") {
			pd.CPUUncore = extractPowerValue(line)
		}

		// 解析Intel核显
		if strings.Contains(line, "Intel 核显") && strings.Contains(line, "UHD 770") {
			pd.IGPU = extractPowerValue(line)
		}

		// 解析AMD独显
		if strings.Contains(line, "AMD 独显") && strings.Contains(line, "6950 XT") {
			pd.DGPU = extractPowerValue(line)
		}

		// 解析HDD
		if strings.Contains(line, "HDD 机械硬盘") {
			pd.HDD = extractPowerValue(line)
		}

		// 解析SSD
		if strings.Contains(line, "SSD 固态硬盘") {
			pd.SSD = extractPowerValue(line)
		}

		// 解析主板内存
		if strings.Contains(line, "主板 & 2x32G内存") {
			pd.MBRAM = extractPowerValue(line)
		}

		// 解析散热
		if strings.Contains(line, "散热") && !strings.Contains(line, "1100W金牌损耗") {
			pd.Cooling = extractPowerValue(line)
		}

		// 解析USB外设
		if strings.Contains(line, "USB及外设功耗") {
			pd.USB = extractPowerValue(line)
		}

		// 解析电源损耗
		if strings.Contains(line, "1100W金牌损耗") {
			pd.PowerLoss = extractPowerValue(line)
		}
	}

	return pd, nil
}

// DirectReadPower 直接从sysfs读取功耗数据（修复版）
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
			continue // 跳过不存在的文件
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
			"averagePower":  0.0,
			"maxPower":     0.0,
			"minPower":     0.0,
			"totalEnergy":  0.0,
			"sampleCount":  0,
		}, nil
	}

	if len(data) == 0 {
		return &map[string]interface{}{
			"averagePower":  0.0,
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
	// 假设采样间隔为5分钟
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
	// 跳过标题行
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
