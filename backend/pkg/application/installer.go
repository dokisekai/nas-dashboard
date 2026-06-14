package application

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// AppInstaller 应用安装器
type AppInstaller struct {
	pkg       *NapPackage
	instance  *AppInstance
	progress chan<- AppInstallProgress
	tempDir   string
}

// NewAppInstaller 创建安装器
func NewAppInstaller(pkg *NapPackage) *AppInstaller {
	return &AppInstaller{
		pkg: pkg,
	}
}

// SetProgress 设置进度回调
func (ai *AppInstaller) SetProgress(progress chan<- AppInstallProgress) {
	ai.progress = progress
}

// sendProgress 发送进度信息
func (ai *AppInstaller) sendProgress(step, message string, percent int, status string) {
	if ai.progress != nil {
		ai.progress <- AppInstallProgress{
			Step:    step,
			Message: message,
			Percent: percent,
			Status:  status,
		}
	}
}

// Install 安装应用
func (ai *AppInstaller) Install(config map[string]interface{}) (*AppInstance, error) {
	// 1. 验证应用包
	ai.sendProgress(StepValidate, "验证应用包", 5, "running")
	if err := ai.validatePackage(); err != nil {
		return nil, fmt.Errorf("应用包验证失败: %w", err)
	}
	ai.sendProgress(StepValidate, "应用包验证完成", 10, "success")

	// 2. 检查依赖
	ai.sendProgress(StepCheckDeps, "检查依赖关系", 15, "running")
	if err := ai.checkDependencies(); err != nil {
		return nil, fmt.Errorf("依赖检查失败: %w", err)
	}
	ai.sendProgress(StepCheckDeps, "依赖检查完成", 20, "success")

	// 3. 创建应用目录
	ai.sendProgress(StepExtract, "创建应用目录", 25, "running")
	installPath := filepath.Join("/var/packages", ai.pkg.Info.Name)
	dataPath := filepath.Join(installPath, "target")
	configPath := filepath.Join(installPath, "config")

	for _, dir := range []string{installPath, dataPath, configPath} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("创建目录 %s 失败: %w", dir, err)
		}
	}
	ai.sendProgress(StepExtract, "应用目录创建完成", 35, "success")

	// 4. 运行预安装脚本
	if ai.pkg.Scripts.PreInstall != "" {
		ai.sendProgress(StepPreInstall, "执行预安装脚本", 40, "running")
		if err := ai.runScript(ai.pkg.Scripts.PreInstall, installPath, config); err != nil {
			return nil, fmt.Errorf("预安装脚本执行失败: %w", err)
		}
		ai.sendProgress(StepPreInstall, "预安装脚本执行完成", 45, "success")
	}

	// 5. 安装应用文件
	ai.sendProgress(StepInstall, "安装应用文件", 50, "running")
	if err := ai.installApplication(installPath, dataPath, configPath); err != nil {
		return nil, fmt.Errorf("应用文件安装失败: %w", err)
	}
	ai.sendProgress(StepInstall, "应用文件安装完成", 70, "success")

	// 6. 运行后安装脚本
	if ai.pkg.Scripts.PostInstall != "" {
		ai.sendProgress(StepPostInstall, "执行后安装脚本", 75, "running")
		if err := ai.runScript(ai.pkg.Scripts.PostInstall, installPath, config); err != nil {
			return nil, fmt.Errorf("后安装脚本执行失败: %w", err)
		}
		ai.sendProgress(StepPostInstall, "后安装脚本执行完成", 80, "success")
	}

	// 7. 创建应用实例
	ai.sendProgress(StepStart, "创建应用实例", 85, "running")
	instance, err := ai.createInstance(installPath, dataPath, configPath, config)
	if err != nil {
		return nil, fmt.Errorf("创建应用实例失败: %w", err)
	}
	ai.sendProgress(StepStart, "应用实例创建完成", 90, "success")

	// 8. 启动应用（如果配置了自动启动）
	if ai.pkg.Info.AutoStart {
		ai.sendProgress(StepStart, "启动应用", 92, "running")
		if err := ai.startApp(instance); err != nil {
			return nil, fmt.Errorf("启动应用失败: %w", err)
		}
		ai.sendProgress(StepStart, "应用启动完成", 95, "success")
	}

	ai.sendProgress(StepComplete, "安装完成", 100, "success")
	return instance, nil
}

// validatePackage 验证应用包
func (ai *AppInstaller) validatePackage() error {
	// 检查磁盘空间
	if ai.pkg.Info.MinDiskSpace > 0 {
		if stat := getDiskStats("/var"); stat.Available < uint64(ai.pkg.Info.MinDiskSpace)*1024*1024*1024 {
			return fmt.Errorf("磁盘空间不足，需要 %d GB，可用 %d GB",
				ai.pkg.Info.MinDiskSpace, stat.Available/(1024*1024*1024))
		}
	}

	// 检查内存
	if ai.pkg.Info.MinRAM > 0 {
		if mem := getMemoryStats(); mem.Total < uint64(ai.pkg.Info.MinRAM)*1024*1024*1024 {
			return fmt.Errorf("内存不足，需要 %d MB，可用 %d MB",
				ai.pkg.Info.MinRAM, mem.Total/(1024*1024))
		}
	}

	return nil
}

// checkDependencies 检查依赖
func (ai *AppInstaller) checkDependencies() error {
	for _, dep := range ai.pkg.Info.Dependencies {
		if !isCommandAvailable(dep) {
			return fmt.Errorf("依赖 %s 不可用", dep)
		}
	}
	return nil
}

// installApplication 安装应用文件
func (ai *AppInstaller) installApplication(installPath, dataPath, configPath string) error {
	// 安装二进制文件
	binDir := filepath.Join(installPath, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return err
	}

	// 复制应用文件
	for _, file := range ai.pkg.Files {
		if file.Type == "binary" {
			srcPath := filepath.Join(ai.tempDir, file.Path)
			dstPath := filepath.Join(binDir, filepath.Base(file.Path))

			// 复制文件
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("复制文件失败: %w", err)
			}

			// 设置执行权限
			if err := os.Chmod(dstPath, 0755); err != nil {
				return fmt.Errorf("设置执行权限失败: %w", err)
			}
		}
	}

	// 复制配置文件
	if err := os.MkdirAll(configPath, 0755); err != nil {
		return err
	}

	for _, file := range ai.pkg.Files {
		if file.Type == "config" {
			srcPath := filepath.Join(ai.tempDir, file.Path)
			dstPath := filepath.Join(configPath, filepath.Base(file.Path))

			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// createInstance 创建应用实例
func (ai *AppInstaller) createInstance(installPath, dataPath, configPath string, config map[string]interface{}) (*AppInstance, error) {
	instance := &AppInstance{
		Name:        ai.pkg.Info.Name,
		DisplayName: ai.pkg.Info.DisplayName,
		PackageName: ai.pkg.Info.Name,
		Version:     ai.pkg.Info.Version,
		Description: ai.pkg.Info.Description,
		Category:    ai.pkg.Info.Category,
		Author:      ai.pkg.Info.Author,
		Website:     ai.pkg.Info.Website,
		Status:      "stopped",
		InstallPath: installPath,
		DataPath:     dataPath,
		ConfigPath:   configPath,
	}

	// 序列化配置
	configData, _ := json.Marshal(config)
	instance.Config = string(configData)

	// 序列化环境变量
	if ai.pkg.Config.EnvVars != nil {
		envData, _ := json.Marshal(ai.pkg.Config.EnvVars)
		instance.EnvVars = string(envData)
	}

	// 序列化资源限制
	if !isEmptyResources(ai.pkg.Resources) {
		resData, _ := json.Marshal(ai.pkg.Resources)
		instance.Resources = string(resData)
	}

	return instance, nil
}

// startApp 启动应用
func (ai *AppInstaller) startApp(instance *AppInstance) error {
	// 如果有启动脚本，使用脚本启动
	if ai.pkg.Scripts.Start != "" {
		scriptPath := filepath.Join(ai.tempDir, ai.pkg.Scripts.Start)
		if _, err := os.Stat(scriptPath); err == nil {
			if err := ai.runScript(scriptPath, instance.InstallPath, nil); err != nil {
				return err
			}
			instance.Status = "running"
			return nil
		}
	}

	// 如果是Docker应用，启动容器
	if instance.ContainerID != "" {
		if err := ai.startContainer(instance); err != nil {
			return err
		}
		instance.Status = "running"
		return nil
	}

	// 直接启动二进制
	if ai.pkg.Scripts.Start == "" {
		binaryPath := filepath.Join(instance.InstallPath, "bin", ai.pkg.Info.Name)
		cmd := exec.Command(binaryPath)

		if err := cmd.Start(); err != nil {
			return fmt.Errorf("启动应用失败: %w", err)
		}

		instance.PID = cmd.Process.Pid
		instance.Status = "running"
	}

	return nil
}

// startContainer 启动Docker容器
func (ai *AppInstaller) startContainer(instance *AppInstance) error {
	// 这里实现Docker容器启动逻辑
	// 简化实现：使用docker命令
	args := []string{"start", instance.ContainerID}

	cmd := exec.Command("docker", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启动Docker容器失败: %w", err)
	}

	return nil
}

// runScript 运行脚本
func (ai *AppInstaller) runScript(scriptPath, workDir string, config map[string]interface{}) error {
	// 检查脚本是否存在
	if _, err := os.Stat(scriptPath); err != nil {
		// 脚本不存在，跳过
		return nil
	}

	// 创建命令
	cmd := exec.Command(scriptPath)
	cmd.Dir = workDir

	// 设置环境变量
	if ai.pkg.Config.EnvVars != nil {
		cmd.Env = os.Environ()
		for key, val := range ai.pkg.Config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
		}
	}

	// 设置配置环境变量
	if config != nil {
		for key, val := range config {
			cmd.Env = append(cmd.Env, fmt.Sprintf("APP_%s=%v", strings.ToUpper(key), val))
		}
	}

	// 运行脚本
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("脚本执行失败: %s, 输出: %s", err.Error(), string(output))
	}

	return nil
}

// StopApp 停止应用
func (ai *AppInstaller) StopApp(instance *AppInstance) error {
	// 如果有停止脚本，使用脚本停止
	if ai.pkg.Scripts.Stop != "" {
		scriptPath := filepath.Join(instance.InstallPath, "scripts", "stop.sh")
		if _, err := os.Stat(scriptPath); err == nil {
			if err := ai.runScript(scriptPath, instance.InstallPath, nil); err != nil {
				return err
			}
			instance.Status = "stopped"
			return nil
		}
	}

	// 如果是Docker容器
	if instance.ContainerID != "" {
		if err := ai.stopContainer(instance); err != nil {
			return err
		}
		instance.Status = "stopped"
		return nil
	}

	// 如果有进程ID，杀死进程
	if instance.PID > 0 {
		if err := syscall.Kill(instance.PID, syscall.SIGTERM); err != nil {
			return fmt.Errorf("杀死进程失败: %w", err)
		}
		instance.Status = "stopped"
		return nil
	}

	return nil
}

// stopContainer 停止Docker容器
func (ai *AppInstaller) stopContainer(instance *AppInstance) error {
	args := []string{"stop", instance.ContainerID}

	cmd := exec.Command("docker", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止Docker容器失败: %w", err)
	}

	return nil
}

// GetAppStatus 获取应用状态
func (ai *AppInstaller) GetAppStatus(instance *AppInstance) (string, error) {
	// 如果有状态检查脚本
	if ai.pkg.Scripts.Status != "" {
		scriptPath := filepath.Join(instance.InstallPath, "scripts", "status.sh")
		if _, err := os.Stat(scriptPath); err == nil {
			cmd := exec.Command(scriptPath)
			output, err := cmd.CombinedOutput()
			if err != nil {
				return "error", fmt.Errorf("状态检查脚本执行失败: %w", err)
			}

			status := strings.TrimSpace(string(output))
			switch status {
			case "running":
				return "running", nil
			case "stopped":
				return "stopped", nil
			default:
				return "unknown", nil
			}
		}
	}

	// 如果是Docker容器
	if instance.ContainerID != "" {
		return ai.getContainerStatus(instance.ContainerID)
	}

	// 检查进程是否存在
	if instance.PID > 0 {
		if err := syscall.Kill(instance.PID, 0); err != nil {
			return "stopped", nil
		}
		return "running", nil
	}

	return "unknown", nil
}

// getContainerStatus 获取容器状态
func (ai *AppInstaller) getContainerStatus(containerID string) (string, error) {
	args := []string{"inspect", "--format", "{{.State.Status}}", containerID}
	cmd := exec.Command("docker", args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "error", err
	}

	status := strings.TrimSpace(string(output))
	switch status {
	case "running":
		return "running", nil
	case "exited":
		return "stopped", nil
	default:
		return "unknown", nil
	}
}

// Uninstall 卸载应用
func (ai *AppInstaller) Uninstall(instance *AppInstance) error {
	// 1. 停止应用
	if err := ai.StopApp(instance); err != nil {
		return fmt.Errorf("停止应用失败: %w", err)
	}

	// 2. 运行预卸载脚本
	if ai.pkg.Scripts.PreUninstall != "" {
		if err := ai.runScript(ai.pkg.Scripts.PreUninstall, instance.InstallPath, nil); err != nil {
			return fmt.Errorf("预卸载脚本执行失败: %w", err)
		}
	}

	// 3. 运行卸载脚本
	if ai.pkg.Scripts.Uninstaller != "" {
		if err := ai.runScript(ai.pkg.Scripts.Uninstaller, instance.InstallPath, nil); err != nil {
			return fmt.Errorf("卸载脚本执行失败: %w", err)
		}
	}

	// 4. 运行后卸载脚本
	if ai.pkg.Scripts.PostUninstall != "" {
		if err := ai.runScript(ai.pkg.Scripts.PostUninstall, instance.InstallPath, nil); err != nil {
			return fmt.Errorf("后卸载脚本执行失败: %w", err)
		}
	}

	// 5. 删除应用目录（保留数据目录）
	os.RemoveAll(instance.InstallPath)
	os.RemoveAll(instance.ConfigPath)

	return nil
}

// 辅助函数

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func isCommandAvailable(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func isEmptyResources(resources NapResources) bool {
	return resources.MaxMemoryMB == 0 &&
		   resources.MaxCPU == 0 &&
		   resources.MaxDiskGB == 0
}

type DiskStats struct {
	Total     uint64
	Used      uint64
	Available uint64
}

type MemStats struct {
	Total uint64
	Used  uint64
}

func getDiskStats(path string) DiskStats {
	var stat syscall.Statfs_t
	syscall.Statfs(path, &stat)

	total := stat.Blocks * uint64(stat.Bsize)
	available := stat.Bavail * uint64(stat.Bsize)

	return DiskStats{
		Total:     total,
		Available: available,
	}
}

func getMemoryStats() MemStats {
	var info syscall.Sysinfo_t
	syscall.Sysinfo(&info)

	total := info.Totalram
	used := info.Totalram - info.Freeram

	return MemStats{
		Total: total,
		Used:  used,
	}
}