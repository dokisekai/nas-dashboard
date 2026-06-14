package models

import (
	"time"

	"gorm.io/gorm"
)

// PowerInfo 功耗信息结构
type PowerInfo struct {
	Timestamp   time.Time `json:"timestamp"`
	CPUPackage  float64   `json:"cpuPackage"`  // CPU Package 功耗 (W)
	CPUCore     float64   `json:"cpuCore"`     // CPU Core 功耗 (W)
	CPUUncore   float64   `json:"cpuUncore"`   // CPU Uncore 功耗 (W)
	IGPU        float64   `json:"igpu"`        // Intel 核显功耗 (W)
	DGPU        float64   `json:"dgpu"`        // AMD 独显功耗 (W)
	HDD         float64   `json:"hdd"`         // 机械硬盘功耗 (W)
	SSD         float64   `json:"ssd"`         // 固态硬盘功耗 (W)
	MBRAM       float64   `json:"mbram"`       // 主板和内存功耗 (W)
	Cooling     float64   `json:"cooling"`     // 散热系统功耗 (W)
	USB         float64   `json:"usb"`         // USB 和外设功耗 (W)
	PowerLoss   float64   `json:"powerLoss"`   // 电源转换损耗 (W)
	Total       float64   `json:"total"`       // 总功耗 (W)
}

// PowerHistory 功耗历史数据
type PowerHistory struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

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
	Labels      string    `gorm:"type:json" json:"labels"` // JSON格式的标签
}

// PowerStatistics 功耗统计数据
type PowerStatistics struct {
	Period         string  `json:"period"`           // daily, weekly, monthly
	AveragePower  float64 `json:"averagePower"`     // 平均功耗 (W)
	MaxPower      float64 `json:"maxPower"`         // 最高功耗 (W)
	MinPower      float64 `json:"minPower"`         // 最低功耗 (W)
	TotalEnergy   float64 `json:"totalEnergy"`      // 总耗电量 (kWh)
	EstimatedCost float64 `json:"estimatedCost"`   // 预估电费 (元)
	SampleCount   int     `json:"sampleCount"`      // 采样次数
	StartDate     time.Time `json:"startDate"`      // 开始时间
	EndDate       time.Time `json:"endDate"`        // 结束时间
}

// PowerAlertRule 功耗告警规则
type PowerAlertRule struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Name      string  `json:"name"`
	Type      string  `json:"type"`        // total, cpu, gpu, hdd, ssd
	Condition string  `json:"condition"`   // >, <, =, !=, >=, <=
	Threshold float64 `json:"threshold"`   // 阈值 (W)
	Duration  int     `json:"duration"`    // 持续时间（秒）
	Severity  string  `json:"severity"`    // info, warning, critical
	Enabled   bool    `json:"enabled"`
	Actions   string  `gorm:"type:json" json:"actions"`     // JSON格式的动作列表
	LastTriggered *time.Time `json:"lastTriggered,omitempty"`
}

// PowerEvent 功耗事件
type PowerEvent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"`        // spike, drop, alert, recovery
	Severity    string    `json:"severity"`     // info, warning, critical
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	Details     string    `gorm:"type:text" json:"details"`
	PowerValue float64   `json:"powerValue"`
	Threshold  float64   `json:"threshold"`
	Resolved    bool      `json:"resolved"`
}

// PowerTrend 功耗趋势分析
type PowerTrend struct {
	Period          string  `json:"period"`          // hourly, daily, weekly
	Trend           string  `json:"trend"`           // increasing, decreasing, stable
	ChangePercent   float64 `json:"changePercent"`   // 变化百分比
	PeakHours       []int   `json:"peakHours"`       // 高峰时段
	AverageByHour   []float64 `json:"averageByHour"` // 每小时平均功耗
	Predictions     []float64 `json:"predictions"`   // 未来预测
	Confidence      float64 `json:"confidence"`     // 预测置信度
}
