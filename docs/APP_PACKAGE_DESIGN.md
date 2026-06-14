# NAS Dashboard 应用包系统设计

## 1. 应用包格式定义

### 应用包结构 (.nap)
```
app-name.nap/
├── INFO                    # 应用元信息文件
├── PACKAGE.tgz             # 应用打包文件
├── icons/                  # 应用图标
│   ├── icon_72.png
│   ├── icon_256.png
│   └── icon_custom.png
├── wizard/                 # 安装向导配置
│   ├── config_ui.json
│   └── steps_config.json
├── scripts/                # 生命周期脚本
│   ├── pre_install.sh
│   ├── installer.sh
│   ├── post_install.sh
│   ├── pre_uninstall.sh
│   ├── uninstaller.sh
│   └── post_uninstall.sh
└── config/                 # 默认配置文件
    └── default_config.json
```

### INFO 文件格式
```ini
[package]
name=plex-media-server
version=1.24.5.4532
displayname=Plex Media Server
description=流媒体服务器，支持多种格式的媒体文件播放
author=Plex Inc
website=https://www.plex.tv
category=Media

[system]
architecture=x86_64
os_version=7.0-10000
min_ram=512
min_disk_space=100
dependencies=docker,python3

[permissions]
required_access=network,storage,process_management
port_bindings=32400,1900
unix_socket=/var/run/app.sock

[resources]
max_memory=2048
max_cpu_cores=4
allowed_networks=lan,wan

[installation]
install_path=/var/packages/app-name
data_path=/var/packages/app-name/target
backup_paths=config,database
requires_restart=false
auto_start=true
```

## 2. 后端应用管理API

### 应用仓库管理
```go
// 应用仓库 API
GET    /api/apps/repository              // 获取应用仓库列表
POST   /api/apps/repository              // 添加应用仓库
PUT    /api/apps/repository/:id         // 更新仓库配置
DELETE /api/apps/repository/:id         // 删除仓库

// 应用包管理
GET    /api/apps/packages               // 获取可用应用包列表
POST   /api/apps/packages/upload        // 上传应用包文件
GET    /api/apps/packages/:name         // 获取应用包详情
DELETE /api/apps/packages/:name         // 删除应用包
```

### 应用实例管理
```go
// 应用安装管理
POST   /api/apps/install                // 安装应用
GET    /api/apps/:id/status             // 获取应用状态
POST   /api/apps/:id/start              // 启动应用
POST   /api/apps/:id/stop               // 停止应用
POST   /api/apps/:id/restart            // 重启应用
PUT    /api/apps/:id/config             // 更新应用配置
DELETE /api/apps/:id                    // 卸载应用

// 应用更新管理
GET    /api/apps/:id/updates           // 检查应用更新
POST   /api/apps/:id/update             // 更新应用
GET    /api/apps/updates                // 获取所有可用更新

// 应用权限管理
GET    /api/apps/:id/permissions        // 获取应用权限
PUT    /api/apps/:id/permissions        // 更新应用权限
```

## 3. 应用隔离实现

### Docker 容器化方案
```go
type AppInstance struct {
    gorm.Model
    Name         string `gorm:"uniqueIndex;not null"`
    DisplayName  string
    PackageName   string
    Version      string
    Status      string // running, stopped, error, installing
    ContainerID string
    Config      string `gorm:"type:json"`
    Ports       string `gorm:"type:json"`
    Volumes     string `gorm:"type:json"`
    Resources   ResourceLimits
    Permissions AppPermissions
}

type ResourceLimits struct {
    MaxMemoryMB int    `json:"max_memory_mb"`
    MaxCPU      int    `json:"max_cpu"`
    MaxDiskGB   int    `json:"max_disk_gb"`
}

type AppPermissions struct {
    NetworkAccess   bool   `json:"network_access"`
    StorageAccess   bool   `json:"storage_access"`
    ProcessAccess   bool   `json:"process_access"`
    PortBindings    []int  `json:"port_bindings"`
    AllowedIPs      []string `json:"allowed_ips"`
}
```

### 应用安装器
```go
func (s *PluginService) InstallApp(req AppInstallRequest) (*AppInstance, error) {
    // 1. 下载并解压应用包
    pkg, err := s.downloadPackage(req.Source)
    if err != nil {
        return nil, err
    }

    // 2. 验证依赖关系
    if err := s.checkDependencies(pkg); err != nil {
        return nil, err
    }

    // 3. 创建应用目录
    if err := s.createAppDirectories(pkg); err != nil {
        return nil, err
    }

    // 4. 运行预安装脚本
    if err := s.runScript(pkg, "pre_install.sh"); err != nil {
        return nil, err
    }

    // 5. 创建Docker容器
    containerID, err := s.createAppContainer(pkg, req.Config)
    if err != nil {
        return nil, err
    }

    // 6. 运行后安装脚本
    if err := s.runScript(pkg, "post_install.sh"); err != nil {
        return nil, err
    }

    // 7. 保存应用实例
    app := s.createAppInstance(pkg, containerID, req.Config)
    if err := s.database.Create(&app).Error; err != nil {
        return nil, err
    }

    return app, nil
}
```

## 4. 前端应用管理界面

### 应用中心组件结构
```vue
<template>
  <div class="app-center">
    <!-- 应用仓库选择 -->
    <div class="repository-selector">
      <select v-model="selectedRepository">
        <option value="official">官方仓库</option>
        <option value="community">社区仓库</option>
        <option value="custom">自定义仓库</option>
      </select>
      <button @click="refreshRepositories">刷新</button>
      <button @click="uploadPackage">上传应用包</button>
    </div>

    <!-- 应用分类 -->
    <div class="app-categories">
      <button v-for="cat in categories" :key="cat.id"
              @click="filterByCategory(cat.id)">
        {{ cat.name }}
      </button>
    </div>

    <!-- 应用列表 -->
    <div class="apps-grid">
      <AppCard v-for="app in filteredApps" :key="app.id"
              :app="app"
              :installed="isInstalled(app.id)"
              @install="installApp(app)"
              @uninstall="uninstallApp(app)"
              @start="startApp(app)"
              @stop="stopApp(app)" />
    </div>

    <!-- 应用详情对话框 -->
    <AppDetailsModal v-if="selectedApp" :app="selectedApp"
                      @install="installApp($event)"
                      @uninstall="uninstallApp($event)" />
  </div>
</template>
```

### 应用卡片组件
```vue
<template>
  <div class="app-card">
    <div class="app-icon" :style="{background: app.color}">
      <component :is="getIcon(app.icon)" />
    </div>
    <div class="app-info">
      <h3>{{ app.name }}</h3>
      <p>{{ app.description }}</p>
      <div class="app-meta">
        <span class="version">v{{ app.version }}</span>
        <span class="category">{{ getCategoryName(app.category) }}</span>
        <div class="rating">
          <StarIcon v-for="i in 5" :key="i"
                   :filled="i <= Math.floor(app.rating)" />
          <span>{{ app.rating }}</span>
        </div>
      </div>
    </div>

    <div class="app-actions">
      <button v-if="installed" @click="start">
        <PlayIcon /> 启动
      </button>
      <button v-if="installed" @click="stop">
        <StopIcon /> 停止
      </button>
      <button v-if="installed" @click="uninstall">
        <TrashIcon /> 卸载
      </button>
      <button v-else @click="install">
        <DownloadIcon /> 安装
      </button>
      <button @click="showDetails">
        <InfoIcon /> 详情
      </button>
    </div>

    <!-- 安装进度条 -->
    <div v-if="installProgress" class="install-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{width: installProgress + '%'}"></div>
      </div>
      <span class="progress-text">{{ installProgressText }}</span>
    </div>
  </div>
</template>
```

## 5. 应用管理功能实现计划

### Phase 1: 应用包格式定义 (1周)
- [ ] 定义.nap应用包格式规范
- [ ] 创建应用包验证工具
- [ ] 实现应用包解析器
- [ ] 设计应用元信息格式

### Phase 2: 后端应用管理系统 (2周)
- [ ] 扩展插件管理API支持应用包
- [ ] 实现应用仓库管理
- [ ] 创建应用安装器
- [ ] 实现应用隔离(Docker)
- [ ] 添加应用权限控制

### Phase 3: 前端应用管理界面 (1周)
- [ ] 扩展AppCenter组件功能
- [ ] 创建AppCard组件
- [ ] 实现应用上传界面
- [ ] 添加应用详情对话框
- [ ] 创建应用安装向导

### Phase 4: 高级功能 (2周)
- [ ] 应用依赖关系管理
- [ ] 应用更新机制
- [ ] 应用备份恢复
- [ ] 应用资源监控
- [ ] 应用日志查看

## 6. 参考DSM的实现细节

### DSM应用启动脚本示例
```bash
#!/bin/sh

# 支持的命令: start, stop, status
case "$1" in
    start)
        # 检查配置文件
        if [ ! -f /usr/syno/synoman/etc/profile ]; then
            echo "Configuration file not found"
            exit 1
        fi

        # 启动应用
        /volume1/@appname/bin/appname --config /usr/syno/synoman/etc/profile

        # 检查启动状态
        if [ $? -eq 0 ]; then
            echo "App started successfully"
            exit 0
        else
            echo "Failed to start app"
            exit 1
        fi
        ;;

    stop)
        # 停止应用
        killall appname
        echo "App stopped"
        exit 0
        ;;

    status)
        # 检查应用状态
        if pgrep appname > /dev/null; then
            echo "running"
            exit 0
        else
            echo "stopped"
            exit 1
        fi
        ;;

    *)
        echo "Usage: $0 {start|stop|status}"
        exit 1
        ;;
esac
```

### DSM INFO文件示例
```ini
[package]
name=CloudSync
version=2782
displayname=Cloud Sync
description=通过简单的设置，即可在云端服务器与Synology NAS之间进行文件同步
author=Synology Inc.
website=https://www.synology.com
category=Sync
version_indep=true

[system]
architecture=x86_64,armv7
os_min_ver=7.0
os_max_ver=
usable_for_personal=true
usable_for_business=true
upgrade_package=CloudSync
beta_app=false

[install_mode]
automatic=true
silent=false
wizard=false

[installer]
    filename=install.sh
    chmod=755

[uninstaller]
    filename=uninstall.sh
    chmod=755

[resource]
    cpu_spare=20
    ram_spare=100

[system_requirement]
    cpu_spare=20
    ram_spare=100

[support]
    web_url=/packages/WebStation
    multiple_version=yes
    install_type=essential
    run_in_dependency=false
    upgrade_from_version=ALL

[sysfolder]
    web=/var/services/web_packages
    home=/var/services/homes
    photo=/var/services/photo
    video=/var/services/video
    music=/var/services/music
    share_folder perm=/volume1
```

## 总结

当前系统已经有基础的插件管理功能，但要实现完整的DSM风格应用包系统，需要：

1. **定义应用包格式** (.nap格式)
2. **扩展后端API** 支持应用包管理
3. **实现应用隔离** 使用Docker容器
4. **完善前端界面** 提供完整的应用管理体验
5. **添加高级功能** 依赖管理、更新机制等

这样才能达到DSM级别的应用管理体验。