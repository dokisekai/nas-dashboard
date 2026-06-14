package application

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PackageParser 应用包解析器
type PackageParser struct {
	packagePath string
	tempDir      string
}

// NewPackageParser 创建解析器
func NewPackageParser() *PackageParser {
	return &PackageParser{}
}

// ParsePackage 解析应用包文件
func (p *PackageParser) ParsePackage(filePath string) (*NapPackage, error) {
	p.packagePath = filePath

	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "nap-package-*")
	if err != nil {
		return nil, fmt.Errorf("创建临时目录失败: %w", err)
	}
	defer os.RemoveAll(tempDir)
	p.tempDir = tempDir

	// 解压tar.gz文件
	if err := p.extractPackage(filePath, tempDir); err != nil {
		return nil, fmt.Errorf("解压包文件失败: %w", err)
	}

	// 读取INFO文件
	info, err := p.readInfoFile(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取INFO文件失败: %w", err)
	}

	// 读取配置文件
	config, err := p.readConfigFile(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 读取向导文件
	wizard, err := p.readWizardFile(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取向导文件失败: %w", err)
	}

	// 读取脚本信息
	scripts, err := p.readScripts(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取脚本信息失败: %w", err)
	}

	// 读取资源限制
	resources, err := p.readResources(tempDir)
	if err != nil {
		return nil, fmt.Errorf("读取资源限制失败: %w", err)
	}

	// 扫描文件列表
	files, err := p.scanFiles(tempDir)
	if err != nil {
		return nil, fmt.Errorf("扫描文件失败: %w", err)
	}

	return &NapPackage{
		Info:      info,
		Files:     files,
		Config:    config,
		Wizard:    wizard,
		Scripts:   scripts,
		Resources: resources,
	}, nil
}

// extractPackage 解压应用包
func (p *PackageParser) extractPackage(filePath, destDir string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建gzip读取器
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("创建gzip读取器失败: %w", err)
	}
	defer gzr.Close()

	// 创建tar读取器
	tr := tar.NewReader(gzr)

	// 解压文件
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("读取tar文件失败: %w", err)
		}

		// 创建目标文件
		targetPath := filepath.Join(destDir, header.Name)
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("创建目录失败: %w", err)
			}
		} else {
			// 创建父目录
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("创建父目录失败: %w", err)
			}

			// 创建文件
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("创建文件失败: %w", err)
			}

			// 复制文件内容
			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return fmt.Errorf("写入文件失败: %w", err)
			}
			outFile.Close()

			// 设置文件权限
			if err := os.Chmod(targetPath, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("设置文件权限失败: %w", err)
			}
		}
	}

	return nil
}

// readInfoFile 读取INFO文件
func (p *PackageParser) readInfoFile(dir string) (NapInfo, error) {
	var info NapInfo

	// 尝试不同的INFO文件格式
	infoPaths := []string{
		filepath.Join(dir, "INFO"),
		filepath.Join(dir, "info"),
		filepath.Join(dir, "package.info"),
		filepath.Join(dir, "app.info"),
	}

	var infoFile *os.File
	var err error

	for _, path := range infoPaths {
		infoFile, err = os.Open(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		return info, fmt.Errorf("找不到INFO文件: %w", err)
	}
	defer infoFile.Close()

	// 读取并解析INFO文件（INI格式）
	if err := p.parseINIFile(infoFile, &info); err != nil {
		// 如果INI解析失败，尝试JSON格式
		if err := json.NewDecoder(infoFile).Decode(&info); err != nil {
			return info, fmt.Errorf("解析INFO文件失败: %w", err)
		}
	}

	return info, nil
}

// parseINIFile 解析INI格式文件
func (p *PackageParser) parseINIFile(file *os.File, dest interface{}) error {
	scanner := bufio.NewScanner(file)
	var currentSection *map[string]string
	var sections = make(map[string]map[string]string)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// 处理节
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.Trim(line, "[]")
			sections[sectionName] = make(map[string]string)
			section := make(map[string]string)
			sections[sectionName] = section
			currentSection = &section
			continue
		}

		// 处理键值对
		if currentSection != nil && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				(*currentSection)[key] = value
			}
		}
	}

	// 简单映射到结构体
	if info, ok := dest.(*NapInfo); ok {
		if section, exists := sections["package"]; exists {
			if name, ok := section["name"]; ok {
				info.Name = name
			}
			if version, ok := section["version"]; ok {
				info.Version = version
			}
			if displayname, ok := section["displayname"]; ok {
				info.DisplayName = displayname
			}
			if description, ok := section["description"]; ok {
				info.Description = description
			}
			if category, ok := section["category"]; ok {
				info.Category = category
			}
			// 更多字段...
		}
	}

	return nil
}

// readConfigFile 读取配置文件
func (p *PackageParser) readConfigFile(dir string) (NapConfig, error) {
	var config NapConfig

	// 读取默认配置
	defaultConfigPath := filepath.Join(dir, "config", "default_config.json")
	if data, err := os.ReadFile(defaultConfigPath); err == nil {
		json.Unmarshal(data, &config.DefaultConfig)
	}

	// 读取用户配置模板
	userConfigPath := filepath.Join(dir, "config", "user_config.json")
	if data, err := os.ReadFile(userConfigPath); err == nil {
		json.Unmarshal(data, &config.UserConfig)
	}

	// 读取环境变量
	envPath := filepath.Join(dir, "config", "env_vars.json")
	if data, err := os.ReadFile(envPath); err == nil {
		json.Unmarshal(data, &config.EnvVars)
	}

	return config, nil
}

// readWizardFile 读取向导文件
func (p *PackageParser) readWizardFile(dir string) (NapWizard, error) {
	var wizard NapWizard

	// 检查向导是否启用
	wizardPath := filepath.Join(dir, "wizard", "wizard.json")
	if _, err := os.Stat(wizardPath); err != nil {
		// 向导不存在，返回空的向导配置
		wizard.Enabled = false
		return wizard, nil
	}

	data, err := os.ReadFile(wizardPath)
	if err != nil {
		return wizard, err
	}

	if err := json.Unmarshal(data, &wizard); err != nil {
		return wizard, err
	}

	wizard.Enabled = true
	return wizard, nil
}

// readScripts 读取脚本信息
func (p *PackageParser) readScripts(dir string) (NapScripts, error) {
	var scripts NapScripts
	scriptsDir := filepath.Join(dir, "scripts")

	// 检查脚本目录是否存在
	if _, err := os.Stat(scriptsDir); err != nil {
		// 脚本目录不存在，返回空配置
		return scripts, nil
	}

	// 读取脚本路径
	scriptsMap := map[string]*string{
		"pre_install":    &scripts.PreInstall,
		"installer":      &scripts.Installer,
		"post_install":   &scripts.PostInstall,
		"pre_uninstall":  &scripts.PreUninstall,
		"uninstaller":    &scripts.Uninstaller,
		"post_uninstall": &scripts.PostUninstall,
		"start":          &scripts.Start,
		"stop":           &scripts.Stop,
		"status":         &scripts.Status,
	}

	for scriptName, scriptPtr := range scriptsMap {
		scriptPath := filepath.Join(scriptsDir, scriptName+".sh")
		if _, err := os.Stat(scriptPath); err == nil {
			*scriptPtr = scriptPath
		}
	}

	return scripts, nil
}

// readResources 读取资源限制
func (p *PackageParser) readResources(dir string) (NapResources, error) {
	var resources NapResources

	// 读取资源配置文件
	resourcesPath := filepath.Join(dir, "config", "resources.json")
	if data, err := os.ReadFile(resourcesPath); err == nil {
		if err := json.Unmarshal(data, &resources); err != nil {
			return resources, err
		}
	}

	return resources, nil
}

// scanFiles 扫描应用包文件
func (p *PackageParser) scanFiles(dir string) ([]NapFile, error) {
	var files []NapFile

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过临时目录本身
		if path == dir {
			return nil
		}

		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		// 获取文件信息
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		// 计算文件哈希
		hash, err := p.calculateFileHash(path)
		if err != nil {
			return err
		}

		// 确定文件类型
		fileType := p.getFileType(relPath)

		files = append(files, NapFile{
			Path: relPath,
			Size: fileInfo.Size(),
			Hash: hash,
			Type: fileType,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// getFileType 获取文件类型
func (p *PackageParser) getFileType(path string) string {
	ext := filepath.Ext(path)
	dir := filepath.Dir(path)

	switch {
	case ext == ".sh":
		return "script"
	case ext == ".json":
		return "config"
	case ext == ".png", ext == ".jpg", ext == ".svg", ext == ".ico":
		return "icon"
	case strings.Contains(dir, "scripts"):
		return "script"
	case strings.Contains(dir, "config"):
		return "config"
	case strings.Contains(dir, "wizard"):
		return "wizard"
	case strings.Contains(dir, "icons"):
		return "icon"
	case strings.Contains(dir, "application"):
		return "binary"
	default:
		return "other"
	}
}

// calculateFileHash 计算文件哈希
func (p *PackageParser) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ValidatePackage 验证应用包
func (p *PackageParser) ValidatePackage(pkg *NapPackage) ([]string, []string, error) {
	var warnings []string
	var errors []string

	// 检查必需字段
	if pkg.Info.Name == "" {
		errors = append(errors, "应用名称不能为空")
	}
	if pkg.Info.Version == "" {
		errors = append(errors, "应用版本不能为空")
	}
	if pkg.Info.Category == "" {
		errors = append(errors, "应用分类不能为空")
	}

	// 检查系统要求
	if pkg.Info.MinRAM > 0 && pkg.Info.MinRAM < 128 {
		warnings = append(warnings, "最小内存要求过低，建议至少128MB")
	}

	// 检查架构兼容性
	if !p.isArchitectureCompatible(pkg.Info.Architecture) {
		warnings = append(warnings, fmt.Sprintf("应用架构 %s 可能不兼容当前系统", pkg.Info.Architecture))
	}

	// 检查依赖
	for _, dep := range pkg.Info.Dependencies {
		if !p.checkDependency(dep) {
			warnings = append(warnings, fmt.Sprintf("依赖 %s 未安装，可能影响应用运行", dep))
		}
	}

	// 检查脚本
	if pkg.Scripts.Installer == "" {
		warnings = append(warnings, "缺少安装脚本，应用可能无法正常安装")
	}
	if pkg.Scripts.Start == "" {
		warnings = append(warnings, "缺少启动脚本，应用可能无法自动启动")
	}

	return warnings, errors, nil
}

// isArchitectureCompatible 检查架构兼容性
func (p *PackageParser) isArchitectureCompatible(arch string) bool {
	if arch == "" {
		return true // 不限制架构
	}

	// 获取系统架构
	systemArch := os.Getenv("GOARCH")
	if systemArch == "" {
		// 默认认为是x86_64
		systemArch = "amd64"
	}

	compatibleArch := map[string][]string{
		"x86_64":   {"x86_64", "amd64", ""},
		"armv7":    {"armv7", "arm"},
		"aarch64":  {"aarch64", "arm64"},
	}

	if allowedArches, exists := compatibleArch[systemArch]; exists {
		for _, allowed := range allowedArches {
			if arch == allowed || arch == "" {
				return true
			}
		}
		return false
	}

	return true
}

// checkDependency 检查依赖是否已安装
func (p *PackageParser) checkDependency(dep string) bool {
	// 这里应该检查系统中是否已安装该依赖
	// 简化实现：检查常见依赖
	commonDeps := map[string]bool{
		"docker":     true,
		"python3":    true,
		"nodejs":     false,
		"java":       false,
	}

	if installed, exists := commonDeps[dep]; exists {
		return installed
	}

	// 尝试检查命令是否存在
	cmd := exec.Command("which", dep)
	output, err := cmd.CombinedOutput()
	return err == nil && len(output) > 0
}