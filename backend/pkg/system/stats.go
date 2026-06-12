package system

import (
	"sync"
	"time"
)

// StatsState 保存统计状态，用于计算差值
type StatsState struct {
	mu sync.RWMutex

	// 磁盘 IO 历史数据
	lastDiskIO    map[string]DiskIOHistory
	lastDiskIOTS  time.Time

	// 网络流量历史数据
	lastNetIO     map[string]NetIOHistory
	lastNetIOTS   time.Time
}

// DiskIOHistory 磁盘 IO 历史统计
type DiskIOHistory struct {
	ReadBytes  uint64
	WriteBytes uint64
	ReadCount  uint64
	WriteCount uint64
}

// NetIOHistory 网络 IO 历史统计
type NetIOHistory struct {
	BytesSent   uint64
	BytesRecv   uint64
	PacketsSent uint64
	PacketsRecv uint64
}

var globalState = &StatsState{
	lastDiskIO:   make(map[string]DiskIOHistory),
	lastNetIO:    make(map[string]NetIOHistory),
	lastDiskIOTS: time.Now(),
	lastNetIOTS:  time.Now(),
}

// GetStatsState 获取全局状态
func GetStatsState() *StatsState {
	return globalState
}

// UpdateDiskIO 更新磁盘 IO 状态
func (s *StatsState) UpdateDiskIO(stats map[string]DiskIOHistory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastDiskIO = stats
	s.lastDiskIOTS = time.Now()
}

// UpdateNetIO 更新网络 IO 状态
func (s *StatsState) UpdateNetIO(stats map[string]NetIOHistory) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastNetIO = stats
	s.lastNetIOTS = time.Now()
}

// GetLastDiskIO 获取上次磁盘 IO 状态
func (s *StatsState) GetLastDiskIO() (map[string]DiskIOHistory, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastDiskIO, s.lastDiskIOTS
}

// GetLastNetIO 获取上次网络 IO 状态
func (s *StatsState) GetLastNetIO() (map[string]NetIOHistory, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastNetIO, s.lastNetIOTS
}

// CalculateDiskSpeed 计算磁盘速度 (bytes/s)
func CalculateDiskSpeed(current, last DiskIOHistory, duration float64) (readSpeed, writeSpeed float64) {
	if duration <= 0 {
		return 0, 0
	}
	readBytes := float64(current.ReadBytes - last.ReadBytes)
	writeBytes := float64(current.WriteBytes - last.WriteBytes)
	return readBytes / duration, writeBytes / duration
}

// CalculateNetSpeed 计算网络速度 (bytes/s)
func CalculateNetSpeed(current, last NetIOHistory, duration float64) (sentSpeed, recvSpeed float64) {
	if duration <= 0 {
		return 0, 0
	}
	sentBytes := float64(current.BytesSent - last.BytesSent)
	recvBytes := float64(current.BytesRecv - last.BytesRecv)
	return sentBytes / duration, recvBytes / duration
}
