package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nas-dashboard/internal/database"
	"nas-dashboard/internal/middleware"
	"nas-dashboard/internal/models"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PluginService 插件服务
type PluginService struct {
	db *gorm.DB
}

// NewPluginService 创建插件服务
func NewPluginService() *PluginService {
	return &PluginService{
		db: database.GetDB(),
	}
}

var pluginService = NewPluginService()

// PluginListRequest 插件列表请求
type PluginListRequest struct {
	IsActive    *bool `form:"isActive"`
	IsInstalled *bool `form:"isInstalled"`
	Category    string `form:"category"`
}

// PluginInstallRequest 插件安装请求
type PluginInstallRequest struct {
	Name       string `json:"name" binding:"required"`
	Source     string `json:"source"` // URL, file path, or registry
	Version    string `json:"version"`
	Config     string `json:"config"` // JSON配置
	AutoEnable bool   `json:"autoEnable"`
}

// PluginUpdateRequest 插件更新请求
type PluginUpdateRequest struct {
	Version string `json:"version"`
	Config  string `json:"config"`
}

// PluginActionRequest 插件操作请求
type PluginActionRequest struct {
	Action string `json:"action" binding:"required"` // enable, disable, start, stop, restart
	Params string `json:"params"` // 操作参数（JSON格式）
}

// GetPlugins 获取插件列表
func GetPlugins(c *gin.Context) {
	var req PluginListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	plugins, err := pluginService.GetPlugins(req.IsActive, req.IsInstalled, req.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"plugins": plugins,
		"total":   len(plugins),
	})
}

// GetPlugin 获取单个插件信息
func GetPlugin(c *gin.Context) {
	name := c.Param("name")

	plugin, err := pluginService.GetPlugin(name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, plugin)
}

// InstallPlugin 安装插件
func InstallPlugin(c *gin.Context) {
	var req PluginInstallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 检查插件是否已安装
	if existingPlugin, err := pluginService.GetPlugin(req.Name); err == nil {
		if existingPlugin.IsInstalled {
			c.JSON(http.StatusConflict, gin.H{"error": "Plugin already installed"})
			return
		}
	}

	// 获取用户ID
	userID, _ := middleware.GetUserIDAsUint(c)

	// 安装插件
	plugin, err := pluginService.InstallPlugin(req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Plugin installed successfully",
		"plugin":  plugin,
	})
}

// UpdatePlugin 更新插件
func UpdatePlugin(c *gin.Context) {
	name := c.Param("name")
	var req PluginUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 更新插件
	plugin, err := pluginService.UpdatePlugin(name, req.Version, req.Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin updated successfully",
		"plugin":  plugin,
	})
}

// UninstallPlugin 卸载插件
func UninstallPlugin(c *gin.Context) {
	name := c.Param("name")

	// 卸载插件
	if err := pluginService.UninstallPlugin(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin uninstalled successfully",
		"name":    name,
	})
}

// EnablePlugin 启用插件
func EnablePlugin(c *gin.Context) {
	name := c.Param("name")

	// 启用插件
	plugin, err := pluginService.EnablePlugin(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin enabled successfully",
		"plugin":  plugin,
	})
}

// DisablePlugin 禁用插件
func DisablePlugin(c *gin.Context) {
	name := c.Param("name")

	// 禁用插件
	plugin, err := pluginService.DisablePlugin(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin disabled successfully",
		"plugin":  plugin,
	})
}

// PluginAction 执行插件操作
func PluginAction(c *gin.Context) {
	name := c.Param("name")
	var req PluginActionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 执行插件操作
	result, err := pluginService.ExecutePluginAction(name, req.Action, req.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin action executed successfully",
		"result":  result,
	})
}

// GetPluginLogs 获取插件日志
func GetPluginLogs(c *gin.Context) {
	name := c.Param("name")
	lines := c.DefaultQuery("lines", "100")

	// 获取插件日志
	logs, err := pluginService.GetPluginLogs(name, lines)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":  name,
		"logs":  logs,
		"lines": lines,
	})
}

// GetPluginConfig 获取插件配置
func GetPluginConfig(c *gin.Context) {
	name := c.Param("name")

	// 获取插件
	plugin, err := pluginService.GetPlugin(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plugin not found"})
		return
	}

	// 解析配置
	var config map[string]interface{}
	if plugin.Config != "" {
		if err := json.Unmarshal([]byte(plugin.Config), &config); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse config"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"name":   name,
		"config": config,
	})
}

// UpdatePluginConfig 更新插件配置
func UpdatePluginConfig(c *gin.Context) {
	name := c.Param("name")
	var config map[string]interface{}

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 序列化配置
	configJSON, err := json.Marshal(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize config"})
		return
	}

	// 更新插件配置
	plugin, err := pluginService.UpdatePlugin(name, "", string(configJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Plugin config updated successfully",
		"plugin":  plugin,
	})
}

// GetPlugins 获取插件列表
func (s *PluginService) GetPlugins(isActive, isInstalled *bool, category string) ([]models.Plugin, error) {
	var plugins []models.Plugin
	query := s.db.Model(&models.Plugin{})

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if isInstalled != nil {
		query = query.Where("is_installed = ?", *isInstalled)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("created_at DESC").Find(&plugins).Error; err != nil {
		return nil, err
	}

	return plugins, nil
}

// GetPlugin 获取单个插件
func (s *PluginService) GetPlugin(name string) (*models.Plugin, error) {
	var plugin models.Plugin
	if err := s.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return nil, err
	}
	return &plugin, nil
}

// InstallPlugin 安装插件
func (s *PluginService) InstallPlugin(req PluginInstallRequest, userID uint) (*models.Plugin, error) {
	// 创建插件目录
	pluginDir := filepath.Join("/opt/nas-dashboard/plugins", req.Name)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// 根据源类型下载或复制插件
	var filePath string
	var err error

	if strings.HasPrefix(req.Source, "http://") || strings.HasPrefix(req.Source, "https://") {
		// 从URL下载
		filePath, err = s.downloadPlugin(req.Source, pluginDir)
	} else if strings.HasPrefix(req.Source, "/") || strings.HasPrefix(req.Source, "./") {
		// 从本地文件系统复制
		filePath = req.Source
	} else {
		// 从插件仓库下载
		filePath, err = s.downloadFromRegistry(req.Source, req.Version, pluginDir)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get plugin source: %w", err)
	}

	// 解压插件
	if err := s.extractPlugin(filePath, pluginDir); err != nil {
		return nil, fmt.Errorf("failed to extract plugin: %w", err)
	}

	// 创建插件数据库记录
	now := time.Now()
	plugin := &models.Plugin{
		Name:        req.Name,
		DisplayName: req.Name,
		Version:     req.Version,
		FilePath:    pluginDir,
		Config:      req.Config,
		IsActive:    req.AutoEnable,
		IsInstalled: true,
		InstalledAt: &now,
	}

	if err := s.db.Create(plugin).Error; err != nil {
		return nil, fmt.Errorf("failed to create plugin record: %w", err)
	}

	// 如果启用，启动插件
	if req.AutoEnable {
		if err := s.startPlugin(plugin); err != nil {
			// 启动失败，仍然返回插件记录
			fmt.Printf("Failed to start plugin: %v\n", err)
		}
	}

	return plugin, nil
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(name, version, config string) (*models.Plugin, error) {
	var plugin models.Plugin
	if err := s.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return nil, err
	}

	// 更新字段
	if version != "" && version != plugin.Version {
		plugin.Version = version
		// 这里可以添加版本更新逻辑
	}

	if config != "" {
		plugin.Config = config
		// 如果插件正在运行，重新加载配置
		if plugin.IsActive {
			if err := s.reloadPluginConfig(&plugin); err != nil {
				return nil, fmt.Errorf("failed to reload plugin config: %w", err)
			}
		}
	}

	if err := s.db.Save(&plugin).Error; err != nil {
		return nil, err
	}

	return &plugin, nil
}

// UninstallPlugin 卸载插件
func (s *PluginService) UninstallPlugin(name string) error {
	var plugin models.Plugin
	if err := s.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return err
	}

	// 如果插件正在运行，先停止
	if plugin.IsActive {
		if err := s.stopPlugin(&plugin); err != nil {
			return fmt.Errorf("failed to stop plugin: %w", err)
		}
	}

	// 删除插件文件
	if plugin.FilePath != "" {
		if err := os.RemoveAll(plugin.FilePath); err != nil {
			return fmt.Errorf("failed to remove plugin files: %w", err)
		}
	}

	// 删除数据库记录
	if err := s.db.Delete(&plugin).Error; err != nil {
		return err
	}

	return nil
}

// EnablePlugin 启用插件
func (s *PluginService) EnablePlugin(name string) (*models.Plugin, error) {
	var plugin models.Plugin
	if err := s.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return nil, err
	}

	if !plugin.IsInstalled {
		return nil, fmt.Errorf("plugin is not installed")
	}

	if plugin.IsActive {
		return &plugin, nil
	}

	// 启动插件
	if err := s.startPlugin(&plugin); err != nil {
		return nil, err
	}

	// 更新状态
	plugin.IsActive = true
	if err := s.db.Save(&plugin).Error; err != nil {
		return nil, err
	}

	return &plugin, nil
}

// DisablePlugin 禁用插件
func (s *PluginService) DisablePlugin(name string) (*models.Plugin, error) {
	var plugin models.Plugin
	if err := s.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return nil, err
	}

	if !plugin.IsActive {
		return &plugin, nil
	}

	// 停止插件
	if err := s.stopPlugin(&plugin); err != nil {
		return nil, err
	}

	// 更新状态
	plugin.IsActive = false
	if err := s.db.Save(&plugin).Error; err != nil {
		return nil, err
	}

	return &plugin, nil
}

// ExecutePluginAction 执行插件操作
func (s *PluginService) ExecutePluginAction(name, action, params string) (map[string]interface{}, error) {
	plugin, err := s.GetPlugin(name)
	if err != nil {
		return nil, err
	}

	if !plugin.IsActive {
		return nil, fmt.Errorf("plugin is not active")
	}

	// 根据操作类型执行相应命令
	var cmd *exec.Cmd
	switch action {
	case "start":
		cmd = s.getPluginCommand(plugin, "start")
	case "stop":
		cmd = s.getPluginCommand(plugin, "stop")
	case "restart":
		cmd = s.getPluginCommand(plugin, "restart")
	case "status":
		cmd = s.getPluginCommand(plugin, "status")
	default:
		return nil, fmt.Errorf("unknown action: %s", action)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w", err)
	}

	// 解析输出
	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		// 如果不是JSON格式，返回原始输出
		result = map[string]interface{}{
			"output": string(output),
		}
	}

	return result, nil
}

// GetPluginLogs 获取插件日志
func (s *PluginService) GetPluginLogs(name, lines string) (string, error) {
	plugin, err := s.GetPlugin(name)
	if err != nil {
		return "", err
	}

	logFile := filepath.Join(plugin.FilePath, "logs", "plugin.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "No logs available", nil
	}

	// 读取日志文件
	cmd := exec.Command("tail", "-n", lines, logFile)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(output), nil
}

// downloadPlugin 下载插件
func (s *PluginService) downloadPlugin(url, destDir string) (string, error) {
	// 使用 wget 或 curl 下载
	destPath := filepath.Join(destDir, "plugin.tar.gz")

	cmd := exec.Command("wget", "-O", destPath, url)
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("wget failed, trying curl: %w", err)

		cmd = exec.Command("curl", "-o", destPath, url)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("curl failed: %w", err)
		}
	}

	return destPath, nil
}

// downloadFromRegistry 从插件仓库下载
func (s *PluginService) downloadFromRegistry(name, version, destDir string) (string, error) {
	// 实现插件仓库下载逻辑
	registryURL := fmt.Sprintf("https://plugins.nas-dashboard.io/%s/%s.tar.gz", name, version)
	return s.downloadPlugin(registryURL, destDir)
}

// extractPlugin 解压插件
func (s *PluginService) extractPlugin(filePath, destDir string) error {
	cmd := exec.Command("tar", "-xzf", filePath, "-C", destDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract plugin: %w", err)
	}
	return nil
}

// startPlugin 启动插件
func (s *PluginService) startPlugin(plugin *models.Plugin) error {
	cmd := s.getPluginCommand(plugin, "start")
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin: %w", err)
	}
	return nil
}

// stopPlugin 停止插件
func (s *PluginService) stopPlugin(plugin *models.Plugin) error {
	cmd := s.getPluginCommand(plugin, "stop")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop plugin: %w, output: %s", err, string(output))
	}
	return nil
}

// reloadPluginConfig 重新加载插件配置
func (s *PluginService) reloadPluginConfig(plugin *models.Plugin) error {
	cmd := s.getPluginCommand(plugin, "reload")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to reload plugin: %w, output: %s", err, string(output))
	}
	return nil
}

// getPluginCommand 获取插件命令
func (s *PluginService) getPluginCommand(plugin *models.Plugin, action string) *exec.Cmd {
	// 假设插件有一个控制脚本
	scriptPath := filepath.Join(plugin.FilePath, "plugin.sh")
	return exec.Command(scriptPath, action)
}

// parsePluginConfig 解析插件配置文件
func parsePluginConfig(configPath string) (map[string]interface{}, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// readPluginManifest 读取插件清单文件
func readPluginManifest(pluginDir string) (map[string]string, error) {
	manifestPath := filepath.Join(pluginDir, "manifest.json")
	file, err := os.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	manifest := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			manifest[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return manifest, nil
}
