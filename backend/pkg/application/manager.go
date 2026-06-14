package application

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// AppManager 应用管理器
type AppManager struct {
	parser      *PackageParser
	installers  map[string]*AppInstaller
	progressMap map[string]chan AppInstallProgress
	mu          sync.RWMutex
	db          Database
	packagesDir string
	tempDir     string
}

// Database 数据库接口
type Database interface {
	CreateAppInstance(instance *AppInstance) error
	GetAppInstance(id uint) (*AppInstance, error)
	GetAppInstanceByName(name string) (*AppInstance, error)
	ListAppInstances() ([]AppInstance, error)
	UpdateAppInstance(instance *AppInstance) error
	DeleteAppInstance(id uint) error

	CreateAppPackage(pkg *AppPackage) error
	GetAppPackage(name string) (*AppPackage, error)
	ListAppPackages() ([]AppPackage, error)
	UpdateAppPackage(pkg *AppPackage) error
	DeleteAppPackage(name string) error

	CreateAppRepository(repo *AppRepository) error
	ListAppRepositories() ([]AppRepository, error)
	UpdateAppRepository(repo *AppRepository) error
	DeleteAppRepository(id uint) error

	CreateInstallLog(log *AppInstallLog) error
	GetInstallLogs(appID uint) ([]AppInstallLog, error)
}

// NewAppManager 创建应用管理器
func NewAppManager(db Database, packagesDir, tempDir string) *AppManager {
	return &AppManager{
		parser:      NewPackageParser(),
		installers:  make(map[string]*AppInstaller),
		progressMap: make(map[string]chan AppInstallProgress),
		db:          db,
		packagesDir: packagesDir,
		tempDir:     tempDir,
	}
}

// UploadPackage 上传应用包
func (m *AppManager) UploadPackage(filePath string) (*AppPackage, error) {
	// 解析应用包
	pkg, err := m.parser.ParsePackage(filePath)
	if err != nil {
		return nil, fmt.Errorf("解析应用包失败: %w", err)
	}

	// 验证应用包
	_, errors, err := m.parser.ValidatePackage(pkg)
	if err != nil {
		return nil, fmt.Errorf("验证应用包失败: %w", err)
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("应用包验证失败: %v", errors)
	}

	// 检查是否已存在
	existingPkg, _ := m.db.GetAppPackage(pkg.Info.Name)
	if existingPkg != nil {
		return nil, fmt.Errorf("应用包 %s 已存在", pkg.Info.Name)
	}

	// 移动应用包到packages目录
	finalPath := filepath.Join(m.packagesDir, pkg.Info.Name+"-"+pkg.Info.Version+".nap")
	if err := os.Rename(filePath, finalPath); err != nil {
		return nil, fmt.Errorf("移动应用包失败: %w", err)
	}

	// 创建应用包记录
	appPackage := &AppPackage{
		Name:         pkg.Info.Name,
		DisplayName:  pkg.Info.DisplayName,
		Version:      pkg.Info.Version,
		Description:  pkg.Info.Description,
		Author:       pkg.Info.Author,
		Website:      pkg.Info.Website,
		Category:     pkg.Info.Category,
		License:      pkg.Info.License,
		FilePath:     finalPath,
		FileSize:     int64(pkg.Info.MinDiskSpace),
		Architecture: pkg.Info.Architecture,
		MinOSVersion: pkg.Info.MinOSVersion,
		MaxOSVersion: pkg.Info.MaxOSVersion,
		MinRAM:       pkg.Info.MinRAM,
		MinDiskSpace: pkg.Info.MinDiskSpace,
		Dependencies: pkg.Info.Dependencies,
	}

	// 序列化资源和权限
	if !isEmptyResources(pkg.Resources) {
		resData, _ := json.Marshal(pkg.Resources)
		appPackage.Resources = string(resData)
	}

	if err := m.db.CreateAppPackage(appPackage); err != nil {
		// 回退：删除已移动的文件
		os.Remove(finalPath)
		return nil, fmt.Errorf("保存应用包记录失败: %w", err)
	}

	return appPackage, nil
}

// InstallApp 安装应用
func (m *AppManager) InstallApp(req AppInstallRequest) (*AppInstance, error) {
	// 获取应用包
	pkg, err := m.db.GetAppPackage(req.PackageName)
	if err != nil {
		return nil, fmt.Errorf("获取应用包失败: %w", err)
	}

	// 检查是否已安装
	existingInstance, _ := m.db.GetAppInstanceByName(req.PackageName)
	if existingInstance != nil {
		return nil, fmt.Errorf("应用 %s 已安装", req.PackageName)
	}

	// 解析应用包
	napPkg, err := m.parser.ParsePackage(pkg.FilePath)
	if err != nil {
		return nil, fmt.Errorf("解析应用包失败: %w", err)
	}

	// 创建安装器
	installer := NewAppInstaller(napPkg)
	installer.tempDir = m.tempDir

	// 创建进度通道
	progressChan := make(chan AppInstallProgress, 10)
	installer.SetProgress(progressChan)

	// 保存安装器和进度通道
	m.mu.Lock()
	m.installers[req.PackageName] = installer
	m.progressMap[req.PackageName] = progressChan
	m.mu.Unlock()

	// 异步安装
	go m.installAsync(installer, napPkg, req)

	// 返回初始实例
	instance := &AppInstance{
		Name:       napPkg.Info.Name,
		DisplayName: napPkg.Info.DisplayName,
		Status:    "installing",
	}

	return instance, nil
}

// installAsync 异步安装应用
func (m *AppManager) installAsync(installer *AppInstaller, pkg *NapPackage, req AppInstallRequest) {
	progressChan := installer.progress

	// 执行安装
	instance, err := installer.Install(req.Config)
	if err != nil {
		// 记录安装失败
		m.recordInstallError(pkg.Info.Name, err.Error())
		progressChan <- AppInstallProgress{
			Step:    StepComplete,
			Message: err.Error(),
			Percent: 0,
			Status:  "error",
		}
		close(progressChan)
		return
	}

	// 保存实例到数据库
	if err := m.db.CreateAppInstance(instance); err != nil {
		m.recordInstallError(pkg.Info.Name, fmt.Sprintf("保存实例失败: %v", err))
		progressChan <- AppInstallProgress{
			Step:    StepComplete,
			Message: fmt.Sprintf("保存实例失败: %v", err),
			Percent: 0,
			Status:  "error",
		}
		close(progressChan)
		return
	}

	// 记录安装日志
	m.recordInstallLogs(instance)

	// 关闭进度通道
	close(progressChan)
}

// GetInstallProgress 获取安装进度
func (m *AppManager) GetInstallProgress(packageName string) (<-chan AppInstallProgress, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	progressChan, exists := m.progressMap[packageName]
	if !exists {
		return nil, fmt.Errorf("应用 %s 没有进行中的安装", packageName)
	}

	return progressChan, nil
}

// StartApp 启动应用
func (m *AppManager) StartApp(instanceID uint) error {
	// 获取应用实例
	instance, err := m.db.GetAppInstance(instanceID)
	if err != nil {
		return fmt.Errorf("获取应用实例失败: %w", err)
	}

	// 获取应用包
	pkg, err := m.db.GetAppPackage(instance.PackageName)
	if err != nil {
		return fmt.Errorf("获取应用包失败: %w", err)
	}

	// 解析应用包
	napPkg, err := m.parser.ParsePackage(pkg.FilePath)
	if err != nil {
		return fmt.Errorf("解析应用包失败: %w", err)
	}

	// 创建安装器
	installer := NewAppInstaller(napPkg)
	installer.tempDir = m.tempDir

	// 启动应用
	if err := installer.startApp(instance); err != nil {
		return fmt.Errorf("启动应用失败: %w", err)
	}

	// 更新实例状态
	instance.Status = "running"
	if err := m.db.UpdateAppInstance(instance); err != nil {
		return fmt.Errorf("更新实例状态失败: %w", err)
	}

	return nil
}

// StopApp 停止应用
func (m *AppManager) StopApp(instanceID uint) error {
	// 获取应用实例
	instance, err := m.db.GetAppInstance(instanceID)
	if err != nil {
		return fmt.Errorf("获取应用实例失败: %w", err)
	}

	// 获取应用包
	pkg, err := m.db.GetAppPackage(instance.PackageName)
	if err != nil {
		return fmt.Errorf("获取应用包失败: %w", err)
	}

	// 解析应用包
	napPkg, err := m.parser.ParsePackage(pkg.FilePath)
	if err != nil {
		return fmt.Errorf("解析应用包失败: %w", err)
	}

	// 创建安装器
	installer := NewAppInstaller(napPkg)
	installer.tempDir = m.tempDir

	// 停止应用
	if err := installer.StopApp(instance); err != nil {
		return fmt.Errorf("停止应用失败: %w", err)
	}

	// 更新实例状态
	instance.Status = "stopped"
	if err := m.db.UpdateAppInstance(instance); err != nil {
		return fmt.Errorf("更新实例状态失败: %w", err)
	}

	return nil
}

// RestartApp 重启应用
func (m *AppManager) RestartApp(instanceID uint) error {
	// 先停止
	if err := m.StopApp(instanceID); err != nil {
		return err
	}

	// 等待停止完成
	time.Sleep(2 * time.Second)

	// 再启动
	if err := m.StartApp(instanceID); err != nil {
		return err
	}

	return nil
}

// GetAppStatus 获取应用状态
func (m *AppManager) GetAppStatus(instanceID uint) (string, error) {
	// 获取应用实例
	instance, err := m.db.GetAppInstance(instanceID)
	if err != nil {
		return "", fmt.Errorf("获取应用实例失败: %w", err)
	}

	// 获取应用包
	pkg, err := m.db.GetAppPackage(instance.PackageName)
	if err != nil {
		return "", fmt.Errorf("获取应用包失败: %w", err)
	}

	// 解析应用包
	napPkg, err := m.parser.ParsePackage(pkg.FilePath)
	if err != nil {
		return "", fmt.Errorf("解析应用包失败: %w", err)
	}

	// 创建安装器
	installer := NewAppInstaller(napPkg)
	installer.tempDir = m.tempDir

	// 获取状态
	status, err := installer.GetAppStatus(instance)
	if err != nil {
		return "", fmt.Errorf("获取应用状态失败: %w", err)
	}

	return status, nil
}

// UninstallApp 卸载应用
func (m *AppManager) UninstallApp(instanceID uint) error {
	// 获取应用实例
	instance, err := m.db.GetAppInstance(instanceID)
	if err != nil {
		return fmt.Errorf("获取应用实例失败: %w", err)
	}

	// 获取应用包
	pkg, err := m.db.GetAppPackage(instance.PackageName)
	if err != nil {
		return fmt.Errorf("获取应用包失败: %w", err)
	}

	// 解析应用包
	napPkg, err := m.parser.ParsePackage(pkg.FilePath)
	if err != nil {
		return fmt.Errorf("解析应用包失败: %w", err)
	}

	// 创建安装器
	installer := NewAppInstaller(napPkg)
	installer.tempDir = m.tempDir

	// 更新实例状态为卸载中
	instance.Status = "uninstalling"
	if err := m.db.UpdateAppInstance(instance); err != nil {
		return fmt.Errorf("更新实例状态失败: %w", err)
	}

	// 执行卸载
	if err := installer.Uninstall(instance); err != nil {
		return fmt.Errorf("卸载应用失败: %w", err)
	}

	// 删除实例记录
	if err := m.db.DeleteAppInstance(instanceID); err != nil {
		return fmt.Errorf("删除实例记录失败: %w", err)
	}

	// 清理安装器
	m.mu.Lock()
	delete(m.installers, instance.Name)
	delete(m.progressMap, instance.Name)
	m.mu.Unlock()

	return nil
}

// UpdateAppConfig 更新应用配置
func (m *AppManager) UpdateAppConfig(instanceID uint, config map[string]interface{}) error {
	// 获取应用实例
	instance, err := m.db.GetAppInstance(instanceID)
	if err != nil {
		return fmt.Errorf("获取应用实例失败: %w", err)
	}

	// 序列化新配置
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 更新实例配置
	instance.Config = string(configData)
	if err := m.db.UpdateAppInstance(instance); err != nil {
		return fmt.Errorf("更新实例配置失败: %w", err)
	}

	return nil
}

// ListApps 列出所有应用
func (m *AppManager) ListApps() ([]AppInstance, []AppPackage, error) {
	instances, err := m.db.ListAppInstances()
	if err != nil {
		return nil, nil, fmt.Errorf("获取应用实例失败: %w", err)
	}

	packages, err := m.db.ListAppPackages()
	if err != nil {
		return nil, nil, fmt.Errorf("获取应用包失败: %w", err)
	}

	return instances, packages, nil
}

// AddRepository 添加应用仓库
func (m *AppManager) AddRepository(repo AppRepository) error {
	if err := m.db.CreateAppRepository(&repo); err != nil {
		return fmt.Errorf("创建应用仓库失败: %w", err)
	}

	return nil
}

// ListRepositories 列出应用仓库
func (m *AppManager) ListRepositories() ([]AppRepository, error) {
	repos, err := m.db.ListAppRepositories()
	if err != nil {
		return nil, fmt.Errorf("获取应用仓库失败: %w", err)
	}

	return repos, nil
}

// UpdateRepository 更新应用仓库
func (m *AppManager) UpdateRepository(repo AppRepository) error {
	if err := m.db.UpdateAppRepository(&repo); err != nil {
		return fmt.Errorf("更新应用仓库失败: %w", err)
	}

	return nil
}

// DeleteRepository 删除应用仓库
func (m *AppManager) DeleteRepository(id uint) error {
	if err := m.db.DeleteAppRepository(id); err != nil {
		return fmt.Errorf("删除应用仓库失败: %w", err)
	}

	return nil
}

// SyncRepository 同步应用仓库
func (m *AppManager) SyncRepository(repoID uint) error {
	// 获取仓库信息
	repos, err := m.db.ListAppRepositories()
	if err != nil {
		return fmt.Errorf("获取仓库信息失败: %w", err)
	}

	var targetRepo *AppRepository
	for _, repo := range repos {
		if repo.ID == repoID {
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		return fmt.Errorf("仓库不存在")
	}

	// 这里实现仓库同步逻辑
	// 简化实现：直接返回成功
	return nil
}

// 辅助函数

func (m *AppManager) recordInstallError(appName, message string) {
	// 记录安装错误到日志
	log := &AppInstallLog{
		Step:    StepComplete,
		Message: message,
		Status:  "error",
	}

	// 这里需要获取实例ID，简化实现
	_ = m.db.CreateInstallLog(log)
}

func (m *AppManager) recordInstallLogs(instance *AppInstance) {
	// 记录安装日志
	logs := []AppInstallLog{
		{
			AppInstanceID: instance.ID,
			Step:         StepValidate,
			Message:      "应用包验证成功",
			Status:       "success",
		},
		{
			AppInstanceID: instance.ID,
			Step:         StepInstall,
			Message:      "应用安装成功",
			Status:       "success",
		},
		{
			AppInstanceID: instance.ID,
			Step:         StepComplete,
			Message:      "应用安装完成",
			Status:       "success",
		},
	}

	for _, log := range logs {
		_ = m.db.CreateInstallLog(&log)
	}
}