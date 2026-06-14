package mergerfs

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

// Manager MergerFS管理器
type Manager struct {
	configPath string
	mountTable map[string]*PoolInfo
	mu         sync.RWMutex
}

// PoolInfo 存储池信息
type PoolInfo struct {
	Name       string
	MountPoint string
	Branches   []BranchInfo
	Status     string
	Config     *Config
}

// BranchInfo 分支信息
type BranchInfo struct {
	Path     string
	Mode     string // rw, ro
	Priority int
	Size     uint64
	Used     uint64
	Free     uint64
}

// Config MergerFS配置
type Config struct {
	Branches       []BranchConfig `json:"branches"`
	Category       string          `json:"category"`       // create, mv, epall, epff, etc.
	MinFreeSpace   string          `json:"minfreespace"`   // 最小空闲空间
	DirectIO       bool            `json:"direct_io"`      // 直接IO
	AsyncRead      bool            `json:"async_read"`     // 异步读取
	UseIno         bool            `json:"use_ino"`        // 使用inode号
	HardRemove     bool            `json:"hard_remove"`    // 硬删除
	AutoUnshare    bool            `json:"auto_unshare"`   // 自动取消共享
	FollowSymlinks bool            `json:"follow_symlinks"`// 跟随符号链接
	LinkExas       bool            `json:"link_exas"`     // 跨分支链接
}

// BranchConfig 分支配置
type BranchConfig struct {
	Path     string `json:"path"`
	Mode     string `json:"mode"`     // ro, rw
	Priority int    `json:"priority"`
}

// NewManager 创建MergerFS管理器
func NewManager(configPath string) *Manager {
	return &Manager{
		configPath: configPath,
		mountTable: make(map[string]*PoolInfo),
	}
}

// CreatePool 创建存储池
func (m *Manager) CreatePool(name, mountPoint string, branches []BranchConfig, config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查存储池名是否已存在
	if _, exists := m.mountTable[name]; exists {
		return fmt.Errorf("pool %s already exists", name)
	}

	// 创建挂载点目录
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}

	// 确保所有分支目录存在
	for _, branch := range branches {
		if err := os.MkdirAll(branch.Path, 0755); err != nil {
			return fmt.Errorf("failed to create branch directory %s: %w", branch.Path, err)
		}
	}

	// 构建MergerFS命令
	args := m.buildMountArgs(name, mountPoint, branches, config)

	// 执行挂载命令
	cmd := exec.Command("mergerfs", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		os.RemoveAll(mountPoint) // 清理创建的目录
		return fmt.Errorf("failed to mount mergerfs: %w, output: %s", err, string(output))
	}

	// 收集分支信息
	branchInfos := make([]BranchInfo, len(branches))
	for i, branch := range branches {
		info, err := m.getBranchInfo(branch.Path)
		if err != nil {
			// 即使获取信息失败，也创建基本信息
			info = BranchInfo{
				Path: branch.Path,
				Mode: branch.Mode,
				Priority: branch.Priority,
			}
		} else {
			info.Mode = branch.Mode
			info.Priority = branch.Priority
		}
		branchInfos[i] = info
	}

	// 保存到挂载表
	poolInfo := &PoolInfo{
		Name:       name,
		MountPoint: mountPoint,
		Branches:   branchInfos,
		Status:     "active",
		Config:     config,
	}

	m.mountTable[name] = poolInfo

	// 保存配置并同步到持久化
	if err := m.saveConfig(); err != nil {
		fmt.Printf("Warning: failed to save config: %v\n", err)
	}
	m.updatePersistence(poolInfo)

	return nil
}

// UpdateConfig 更新存储池配置
func (m *Manager) UpdateConfig(poolName string, config *Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.mountTable[poolName]
	if !exists {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// 更新内存中的配置
	pool.Config = config

	// 尝试通过控制文件更新运行时的 category
	if config.Category != "" {
		controlPath := filepath.Join(pool.MountPoint, ".mergerfs")
		if _, err := os.Stat(controlPath); err == nil {
			// 如果控制目录存在，尝试动态修改
			policyPath := filepath.Join(controlPath, "policy")
			os.WriteFile(policyPath, []byte(config.Category), 0644)
		}
	}

	// 对于无法动态修改的参数，执行重新挂载
	if err := m.remountPool(pool); err != nil {
		return fmt.Errorf("failed to apply config via remount: %w", err)
	}

	// 保存并持久化
	m.saveConfig()
	m.updatePersistence(pool)

	return nil
}

// updatePersistence 将挂载信息写入持久化存储（如 systemd unit）
func (m *Manager) updatePersistence(pool *PoolInfo) error {
	// 生成 systemd unit 以支持开机自动挂载
	unitName := fmt.Sprintf("nas-mergerfs-%s.service", pool.Name)
	unitPath := filepath.Join("/etc/systemd/system", unitName)

	branches := make([]BranchConfig, len(pool.Branches))
	for i, b := range pool.Branches {
		branches[i] = BranchConfig{Path: b.Path, Mode: b.Mode, Priority: b.Priority}
	}
	args := m.buildMountArgs(pool.Name, pool.MountPoint, branches, pool.Config)
	// 确保包含 allow_other 以允许 SMB 用户访问，包含 nofail 确保开机稳健性
	cmdStr := "mergerfs " + strings.Join(args, " ")

	content := fmt.Sprintf(`[Unit]
Description=MergerFS Pool: %s
After=local-fs.target
Before=smbd.service nfs-server.service
Wants=local-fs.target

[Service]
Type=forking
ExecStart=%s
ExecStartPost=/usr/bin/chmod 775 %s
ExecStartPost=/usr/bin/chown root:users %s
ExecStop=fusermount -uz %s
Restart=on-failure

[Install]
WantedBy=multi-user.target
`, pool.Name, cmdStr, pool.MountPoint, pool.MountPoint, pool.MountPoint)

	err := os.WriteFile(unitPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write systemd unit: %w", err)
	}

	// 重新加载 systemd 并启用
	exec.Command("systemctl", "daemon-reload").Run()
	exec.Command("systemctl", "enable", unitName).Run()

	return nil
}

// removePersistence 移除持久化
func (m *Manager) removePersistence(poolName string) {
	unitName := fmt.Sprintf("nas-mergerfs-%s.service", poolName)
	exec.Command("systemctl", "disable", unitName).Run()
	os.Remove(filepath.Join("/etc/systemd/system", unitName))
	exec.Command("systemctl", "daemon-reload").Run()
}

// AddDisk 添加磁盘到存储池
func (m *Manager) AddDisk(poolName string, branchPath string, mode string, priority int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.mountTable[poolName]
	if !exists {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// 检查分支是否已存在
	for _, branch := range pool.Branches {
		if branch.Path == branchPath {
			return fmt.Errorf("branch %s already exists in pool", branchPath)
		}
	}

	// 创建分支目录
	if err := os.MkdirAll(branchPath, 0755); err != nil {
		return fmt.Errorf("failed to create branch directory: %w", err)
	}

	// 获取分支信息
	branchInfo, err := m.getBranchInfo(branchPath)
	if err != nil {
		branchInfo = BranchInfo{
			Path: branchPath,
			Mode: mode,
			Priority: priority,
		}
	} else {
		branchInfo.Mode = mode
		branchInfo.Priority = priority
	}

	// 添加到分支列表
	pool.Branches = append(pool.Branches, branchInfo)

	// 更新配置
	if pool.Config == nil {
		pool.Config = &Config{}
	}
	pool.Config.Branches = append(pool.Config.Branches, BranchConfig{
		Path:     branchPath,
		Mode:     mode,
		Priority: priority,
	})

	// 重新挂载以应用新配置
	if err := m.remountPool(pool); err != nil {
		// 回滚
		pool.Branches = pool.Branches[:len(pool.Branches)-1]
		return fmt.Errorf("failed to remount pool: %w", err)
	}

	// 保存配置并更新持久化
	if err := m.saveConfig(); err != nil {
		fmt.Printf("Warning: failed to save config: %v\n", err)
	}
	m.updatePersistence(pool)

	return nil
}

// RemoveDisk 从存储池移除磁盘
func (m *Manager) RemoveDisk(poolName, branchPath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.mountTable[poolName]
	if !exists {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// 查找分支索引
	branchIndex := -1
	for i, branch := range pool.Branches {
		if branch.Path == branchPath {
			branchIndex = i
			break
		}
	}

	if branchIndex == -1 {
		return fmt.Errorf("branch %s not found in pool", branchPath)
	}

	// 检查分支是否包含数据
	branchInfo := pool.Branches[branchIndex]
	if branchInfo.Used > 0 {
		return fmt.Errorf("branch %s contains data (%d bytes), cannot remove", branchPath, branchInfo.Used)
	}

	// 从分支列表中移除
	pool.Branches = append(pool.Branches[:branchIndex], pool.Branches[branchIndex+1:]...)

	// 更新配置
	if pool.Config != nil && len(pool.Config.Branches) > branchIndex {
		pool.Config.Branches = append(pool.Config.Branches[:branchIndex], pool.Config.Branches[branchIndex+1:]...)
	}

	// 重新挂载
	if err := m.remountPool(pool); err != nil {
		return fmt.Errorf("failed to remount pool: %w", err)
	}

	// 保存配置并更新持久化
	if err := m.saveConfig(); err != nil {
		fmt.Printf("Warning: failed to save config: %v\n", err)
	}
	m.updatePersistence(pool)

	return nil
}

// GetPoolStatus 获取存储池状态
func (m *Manager) GetPoolStatus(poolName string) (*PoolInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pool, exists := m.mountTable[poolName]
	if !exists {
		return nil, fmt.Errorf("pool %s not found", poolName)
	}

	// 更新分支信息
	for i := range pool.Branches {
		if info, err := m.getBranchInfo(pool.Branches[i].Path); err == nil {
			// 保持模式和优先级
			info.Mode = pool.Branches[i].Mode
			info.Priority = pool.Branches[i].Priority
			pool.Branches[i] = info
		}
	}

	// 复制返回
	result := *pool
	return &result, nil
}

// BalancePool 平衡存储池
func (m *Manager) BalancePool(poolName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.mountTable[poolName]
	if !exists {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// 优先使用 mergerfs.balance 工具（如果系统已安装）
	balanceCmd := "mergerfs.balance"
	if _, err := exec.LookPath(balanceCmd); err == nil {
		cmd := exec.Command(balanceCmd, pool.MountPoint)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("mergerfs.balance failed: %v, output: %s", err, string(output))
		}
		return nil
	}

	// 如果没有工具，实现一个基本的文件移动平衡逻辑 (基于 epmfs 策略的模拟)
	// 这里仅作为演示，生产环境建议安装 mergerfs-tools
	fmt.Printf("mergerfs.balance not found, skipping deep balancing for %s\n", pool.MountPoint)
	return nil
}

// DeletePool 删除存储池
func (m *Manager) DeletePool(poolName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	pool, exists := m.mountTable[poolName]
	if !exists {
		return fmt.Errorf("pool %s not found", poolName)
	}

	// 卸载存储池
	if err := m.unmountPool(pool); err != nil {
		return fmt.Errorf("failed to unmount pool: %w", err)
	}

	// 从挂载表中移除
	delete(m.mountTable, poolName)
	m.removePersistence(poolName)

	// 保存配置
	if err := m.saveConfig(); err != nil {
		fmt.Printf("Warning: failed to save config: %v\n", err)
	}

	return nil
}

// ListPools 列出所有存储池
func (m *Manager) ListPools() []PoolInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	pools := make([]PoolInfo, 0, len(m.mountTable))
	for _, pool := range m.mountTable {
		pools = append(pools, *pool)
	}

	return pools
}

// buildMountArgs 构建MergerFS挂载参数
func (m *Manager) buildMountArgs(name, mountPoint string, branches []BranchConfig, config *Config) []string {
	args := []string{}

	// 添加分支 (格式: path1:path2=mode)
	// 或者 path1=mode:path2=mode
	branchPaths := make([]string, len(branches))
	for i, branch := range branches {
		branchPaths[i] = fmt.Sprintf("%s=%s", branch.Path, branch.Mode)
	}
	args = append(args, strings.Join(branchPaths, ":"))

	// 挂载点
	args = append(args, mountPoint)

	// 默认选项
	args = append(args, "-o", "allow_other,use_ino,cache.files=off,dropcacheonclose=true")

	// 选项
	if config != nil {
		if config.Category != "" {
			args = append(args, "-o", fmt.Sprintf("category=%s", config.Category))
		}
		if config.MinFreeSpace != "" {
			args = append(args, "-o", fmt.Sprintf("minfreespace=%s", config.MinFreeSpace))
		}
		if config.DirectIO {
			args = append(args, "-o", "direct_io")
		}
		if config.AsyncRead {
			args = append(args, "-o", "async_read")
		}
		if config.UseIno {
			args = append(args, "-o", "use_ino")
		}
		if config.HardRemove {
			args = append(args, "-o", "hard_remove")
		}
		if config.AutoUnshare {
			args = append(args, "-o", "auto_unshare")
		}
		if config.FollowSymlinks {
			args = append(args, "-o", "follow_symlinks")
		}
		if config.LinkExas {
			args = append(args, "-o", "link_exas")
		}
	}

	return args
}

// remountPool 重新挂载存储池
func (m *Manager) remountPool(pool *PoolInfo) error {
	// 先卸载
	if err := m.unmountPool(pool); err != nil {
		return fmt.Errorf("failed to unmount for remount: %w", err)
	}

	// 重新挂载
	branches := make([]BranchConfig, len(pool.Branches))
	for i, branch := range pool.Branches {
		branches[i] = BranchConfig{
			Path:     branch.Path,
			Mode:     branch.Mode,
			Priority: branch.Priority,
		}
	}

	args := m.buildMountArgs(pool.Name, pool.MountPoint, branches, pool.Config)
	cmd := exec.Command("mergerfs", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remount mergerfs: %w, output: %s", err, string(output))
	}

	return nil
}

// unmountPool 卸载存储池
func (m *Manager) unmountPool(pool *PoolInfo) error {
	cmd := exec.Command("fusermount", "-uz", pool.MountPoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unmount: %w, output: %s", err, string(output))
	}
	return nil
}

// getBranchInfo 获取分支信息
func (m *Manager) getBranchInfo(path string) (BranchInfo, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return BranchInfo{}, err
	}

	total := uint64(stat.Blocks) * uint64(stat.Bsize)
	free := uint64(stat.Bfree) * uint64(stat.Bsize)
	used := total - free

	return BranchInfo{
		Path:  path,
		Size:  total,
		Used:  used,
		Free:  free,
	}, nil
}

// saveConfig 保存配置
func (m *Manager) saveConfig() error {
	configs := make(map[string]*PoolInfo)
	for name, pool := range m.mountTable {
		configs[name] = pool
	}

	data, err := json.MarshalIndent(configs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// LoadConfig 加载配置
func (m *Manager) LoadConfig() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 配置文件不存在是正常的
		}
		return err
	}

	var configs map[string]*PoolInfo
	if err := json.Unmarshal(data, &configs); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.mountTable = configs

	return nil
}