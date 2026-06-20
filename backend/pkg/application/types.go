package application

import "time"

// NapPackage 应用包结构
type NapPackage struct {
	Info      NapInfo       `json:"info"`
	Files     []NapFile      `json:"files"`
	Config    NapConfig      `json:"config"`
	Wizard    NapWizard     `json:"wizard"`
	Scripts   NapScripts    `json:"scripts"`
	Resources NapResources  `json:"resources"`
}

// NapInfo 应用包信息
type NapInfo struct {
	Name            string `json:"name" binding:"required"`
	DisplayName      string `json:"displayName" binding:"required"`
	Version         string `json:"version" binding:"required"`
	Description     string `json:"description"`
	Author          string `json:"author"`
	Website         string `json:"website"`
	Category        string `json:"category" binding:"required"`
	License         string `json:"license"`
	Icon            string `json:"icon"`
	AppType         string `json:"appType"` // docker, system

	// 系统要求
	Architecture    string   `json:"architecture"`
	MinOSVersion    string   `json:"minOSVersion"`
	MaxOSVersion    string   `json:"maxOSVersion"`
	MinRAM         int      `json:"minRAM"`
	MinDiskSpace   int      `json:"minDiskSpace"`
	Dependencies   []string `json:"dependencies"`

	// 安装配置
	InstallPath     string `json:"installPath"`
	DataPath        string `json:"dataPath"`
	ConfigPath      string `json:"configPath"`
	BackupPaths     []string `json:"backupPaths"`
	RequiresRestart bool    `json:"requiresRestart"`
	AutoStart       bool    `json:"autoStart"`

	// 版本信息
	PackageVersion  string `json:"packageVersion"`
	BuildDate       string `json:"buildDate"`
}

// NapFile 应用包文件
type NapFile struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
	Type     string `json:"type"` // binary, config, script, wizard, icon
}

// NapConfig 应用配置
type NapConfig struct {
	DefaultConfig map[string]interface{} `json:"defaultConfig"`
	UserConfig   map[string]interface{} `json:"userConfig"`
	EnvVars       map[string]string      `json:"envVars"`
}

// NapWizard 安装向导
type NapWizard struct {
	Enabled  bool              `json:"enabled"`
	Steps    []NapWizardStep  `json:"steps"`
	ConfigUI string           `json:"configUI"`
}

// NapWizardStep 向导步骤
type NapWizardStep struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Type        string                 `json:"type"` // form, confirm, custom
	Required    bool                   `json:"required"`
	Config      map[string]interface{} `json:"config"`
}

// NapScripts 脚本定义
type NapScripts struct {
	PreInstall      string `json:"preInstall"`
	Installer       string `json:"installer"`
	PostInstall     string `json:"postInstall"`
	PreUninstall    string `json:"preUninstall"`
	Uninstaller     string `json:"uninstaller"`
	PostUninstall   string `json:"postUninstall"`
	Start           string `json:"start"`
	Stop            string `json:"stop"`
	Status          string `json:"status"`
}

// NapResources 资源限制
type NapResources struct {
	MaxMemoryMB    int      `json:"maxMemoryMB"`
	MaxCPU         int      `json:"maxCPU"`
	MaxDiskGB      int      `json:"maxDiskGB"`
	NetworkAccess bool     `json:"networkAccess"`
	StorageAccess  bool     `json:"storageAccess"`
	ProcessAccess  bool     `json:"processAccess"`
	PortBindings   []int    `json:"portBindings"`
	AllowedIPs     []string `json:"allowedIps"`
}

// AppInstance 应用实例
type AppInstance struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// 基本信息
	Name           string    `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName    string    `json:"displayName"`
	PackageName    string    `json:"packageName"`
	Version        string    `json:"version"`
	Description    string    `json:"description"`
	Category       string    `json:"category"`
	Author         string    `json:"author"`
	Website        string    `json:"website"`
	AppType        string    `json:"appType"` // docker, system

	// 状态管理
	Status         string    `json:"status"`      // running, stopped, error, installing, uninstalling
	ContainerID    string    `json:"containerId"`
	PID            int       `json:"pid"`
	ExitCode       int       `json:"exitCode"`
	LastExitTime   time.Time `json:"lastExitTime"`

	// 配置
	Config         string    `gorm:"type:json" json:"config"`
	EnvVars        string    `gorm:"type:json" json:"envVars"`
	Ports          string    `gorm:"type:json" json:"ports"`
	Volumes        string    `gorm:"type:json" json:"volumes"`
	Resources      string    `gorm:"type:json" json:"resources"`

	// 权限
	Permissions    string    `gorm:"type:json" json:"permissions"`

	// 安装信息
	InstallPath    string    `json:"installPath"`
	DataPath       string    `json:"dataPath"`
	ConfigPath     string    `json:"configPath"`
	BackupPaths    string    `gorm:"type:json" json:"backupPaths"`

	// 关系
	InstallLog     []AppInstallLog `gorm:"foreignKey:AppInstanceID" json:"installLog,omitempty"`
}

// AppInstallLog 安装日志
type AppInstallLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	AppInstanceID uint      `gorm:"index" json:"appInstanceId"`
	Step         string    `json:"step"`
	Message      string    `json:"message"`
	Status       string    `json:"status"` // success, error, warning
	Details      string    `json:"details"`
}

// AppRepository 应用仓库
type AppRepository struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	URL         string    `json:"url"`
	Type        string    `json:"type"` // official, community, custom
	Enabled     bool      `json:"enabled"`
	Priority    int       `json:"priority"`
	Description string    `json:"description"`
	AutoUpdate  bool      `json:"autoUpdate"`
}

// AppPackage 应用包信息
type AppPackage struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`

	Name         string    `gorm:"uniqueIndex;not null" json:"name"`
	DisplayName  string    `json:"displayName"`
	Version      string    `json:"version"`
	Description  string    `json:"description"`
	Author       string    `json:"author"`
	Website      string    `json:"website"`
	Category     string    `json:"category"`
	License      string    `json:"license"`
	AppType      string    `json:"appType"` // docker, system

	// 包文件信息
	FilePath     string    `json:"filePath"`
	FileSize     int64     `json:"fileSize"`
	FileHash     string    `json:"fileHash"`
	DownloadURL  string    `json:"downloadUrl"`

	// 系统要求
	Architecture string   `json:"architecture"`
	MinOSVersion string   `json:"minOSVersion"`
	MaxOSVersion string   `json:"maxOSVersion"`
	MinRAM       int      `json:"minRAM"`
	MinDiskSpace int      `json:"minDiskSpace"`
	Dependencies []string `json:"dependencies"`

	// 安装配置
	InstallPath  string `json:"installPath"`
	DataPath     string `json:"dataPath"`
	ConfigPath   string `json:"configPath"`
	BackupPaths  string   `json:"backupPaths"`

	// 资源和权限
	Resources    string `gorm:"type:json" json:"resources"`
	Permissions  string `gorm:"type:json" json:"permissions"`

	// 仓库信息
	RepositoryID uint      `gorm:"index" json:"repositoryId"`
	Repository   *AppRepository `gorm:"foreignKey:RepositoryID" json:"repository,omitempty"`

	// 统计
	DownloadCount int     `json:"downloadCount"`
	InstallCount  int     `json:"installCount"`
	Rating        float64   `json:"rating"`

	// 关系
	Instances     []AppInstance `gorm:"foreignKey:PackageName" json:"instances,omitempty"`
}

// AppInstallRequest 应用安装请求
type AppInstallRequest struct {
	PackageName string                 `json:"packageName" binding:"required"`
	Version     string                 `json:"version"`
	Config       map[string]interface{} `json:"config"`
	AutoStart    bool                   `json:"autoStart"`
}

// AppInstallProgress 安装进度
type AppInstallProgress struct {
	Step    string `json:"step"`
	Message string `json:"message"`
	Percent int    `json:"percent"`
	Status  string `json:"status"` // running, success, error
}

// 定义安装步骤
const (
	StepValidate     = "validate"
	StepDownload     = "download"
	StepExtract      = "extract"
	StepCheckDeps    = "check_deps"
	StepPreInstall   = "pre_install"
	StepInstall      = "install"
	StepPostInstall  = "post_install"
	StepStart        = "start"
	StepComplete     = "complete"
)

// AppActionRequest 应用操作请求
type AppActionRequest struct {
	Action string                 `json:"action" binding:"required"` // start, stop, restart, reload
	Params map[string]interface{} `json:"params"`
}

// 应用操作常量
const (
	ActionStart  = "start"
	ActionStop   = "stop"
	ActionRestart = "restart"
	ActionReload = "reload"
	ActionEnable = "enable"
	ActionDisable = "disable"
)

// 应用类型常量
const (
	AppTypeDocker = "docker"
	AppTypeSystem = "system"
)