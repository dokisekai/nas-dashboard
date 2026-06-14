# NAS Dashboard 应用包制作完整指南

## 目录
1. [快速开始](#快速开始)
2. [应用包结构详解](#应用包结构详解)
3. [INFO文件详解](#info文件详解)
4. [脚本编写指南](#脚本编写指南)
5. [实战案例](#实战案例)
6. [调试技巧](#调试技巧)
7. [常见问题](#常见问题)

## 快速开始

### 创建第一个应用包（5分钟）

```bash
# 1. 创建应用目录
mkdir -p my-app/{application,icons,scripts,config,wizard}

# 2. 创建INFO文件
cat > my-app/INFO << 'EOF'
[package]
name=my-app
version=1.0.0
displayname=My First App
description=我的第一个应用
author=Your Name
website=https://example.com
category=utilities
license=MIT

[system]
architecture=x86_64
min_os_version=1.0.0
min_ram=128
min_disk_space=100
dependencies=docker

[installation]
install_path=/var/packages/my-app
data_path=/var/packages/my-app/target
auto_start=true
EOF

# 3. 创建简单的启动脚本
cat > my-app/scripts/start.sh << 'EOF'
#!/bin/bash
echo "Starting my app..."
# 你的应用启动逻辑
EOF
chmod +x my-app/scripts/start.sh

# 4. 创建停止脚本
cat > my-app/scripts/stop.sh << 'EOF'
#!/bin/bash
echo "Stopping my app..."
# 你的应用停止逻辑
EOF
chmod +x my-app/scripts/stop.sh

# 5. 创建状态检查脚本
cat > my-app/scripts/status.sh << 'EOF'
#!/bin/bash
echo "running"  # 返回: running, stopped, error
EOF
chmod +x my-app/scripts/status.sh

# 6. 打包
cd my-app
tar -czf ../my-app-1.0.0.nap .
cd ..
```

## 应用包结构详解

### 标准目录结构
```
app-name.nap/
├── INFO                          # 必需: 应用元信息
├── application/                   # 应用文件目录
│   ├── app-name                 # 必需: 应用主程序
│   ├── lib/                     # 可选: 依赖库
│   └── resources/               # 可选: 资源文件
├── icons/                        # 必需: 应用图标
│   ├── icon_72.png              # 必需: 72x72 小图标
│   ├── icon_256.png            # 必需: 256x256 大图标
│   └── icon_custom.png         # 可选: 自定义尺寸
├── scripts/                      # 必需: 生命周期脚本
│   ├── installer.sh            # 可选: 安装脚本
│   ├── pre_install.sh          # 可选: 预安装脚本
│   ├── post_install.sh         # 可选: 后安装脚本
│   ├── start.sh                # 必需: 启动脚本
│   ├── stop.sh                 # 必需: 停止脚本
│   ├── status.sh               # 必需: 状态检查脚本
│   ├── pre_uninstall.sh        # 可选: 预卸载脚本
│   ├── uninstaller.sh          # 可选: 卸载脚本
│   └── post_uninstall.sh       # 可选: 后卸载脚本
├── config/                       # 可选: 配置文件
│   ├── default_config.json     # 默认配置
│   ├── user_config.json        # 用户配置模板
│   ├── env_vars.json           # 环境变量
│   └── resources.json          # 资源限制
└── wizard/                       # 可选: 安装向导
    ├── wizard.json             # 向导配置
    └── steps/                  # 向导步骤
        └── *.html              # 自定义HTML界面
```

## INFO文件详解

### 完整INFO文件示例

```ini
[package]
# === 基本信息 ===
name=plex-media-server              # 应用唯一标识（必须小写字母、数字、连字符）
version=1.24.5.4532                 # 版本号（语义化版本）
displayname=Plex Media Server       # 显示名称（支持中文）
description=功能强大的媒体服务器     # 详细描述
author=Plex Inc                      # 作者名称
website=https://www.plex.tv         # 官方网站
category=media                      # 分类: media|productivity|utilities|security|network
license=MIT                         # 开源许可证

# === 系统要求 ===
[system]
architecture=x86_64                # 架构: x86_64|armv7|aarch64|noarch
min_os_version=1.0.0               # 最低系统版本
max_os_version=2.0.0               # 最高系统版本（可选）
min_ram=512                        # 最小内存（MB）
min_disk_space=1000                # 最小磁盘空间（MB）
dependencies=docker,python3,nodejs # 依赖列表（逗号分隔）

# === 安装配置 ===
[installation]
install_path=/var/packages/plex     # 安装路径（自动创建）
data_path=/var/packages/plex/target # 数据路径（自动创建）
config_path=/var/packages/plex/config # 配置路径
backup_paths=config,database        # 备份路径（逗号分隔）
requires_restart=false              # 是否需要重启系统
auto_start=true                     # 是否自动启动
install_mode=automatic              # 安装模式: automatic|manual

# === 权限配置 ===
[permissions]
network_access=true                 # 网络访问权限
storage_access=true                 # 存储访问权限
process_access=false                # 进程管理权限
device_access=false                 # 设备访问权限
port_bindings=32400,1900,3005       # 端口绑定（逗号分隔）
unix_socket=/var/run/plex.sock     # Unix套接字路径
allowed_networks=lan,wan           # 允许的网络: lan|wan|local

# === 资源限制 ===
[resources]
max_memory=2048                    # 最大内存（MB）
max_cpu_cores=4                    # 最大CPU核心数
max_disk_io=100                    # 最大磁盘IO（MB/s）
max_network_bandwidth=1000         # 最大网络带宽（Mbps）

# === 服务配置 ===
[service]
run_as=root                        # 运行用户: root|nobody|特定用户
start_priority=10                   # 启动优先级（1-100，数字越小越优先）
startup_timeout=120                # 启动超时时间（秒）
shutdown_timeout=60                # 关闭超时时间（秒）
restart_on_failure=true            # 失败自动重启
restart_delay=10                   # 重启延迟（秒）
max_restart_count=3                # 最大重启次数
restart_window=300                 # 重启时间窗口（秒）

# === 更新配置 ===
[update]
auto_update=true                   # 自动更新
update_channel=stable              # 更新通道: stable|beta|alpha
update_check_interval=86400        # 更新检查间隔（秒）
backup_before_update=true          # 更新前备份
rollback_on_failure=true           # 失败后回滚

# === 支持信息 ===
[support]
documentation_url=https://docs.plex.tv  # 文档URL
support_url=https://forums.plex.tv     # 支持URL
bug_report_url=https://issues.plex.tv  # 问题报告URL
```

### 分类说明

| 分类 | 说明 | 示例应用 |
|------|------|----------|
| `media` | 媒体相关 | Plex, Emby, Jellyfin |
| `productivity` | 办公生产力 | Nextcloud, OnlyOffice |
| `utilities` | 系统工具 | FileBrowser, Shellinabox |
| `security` | 安全相关 | VPN, 防火墙 |
| `network` | 网络相关 | Nginx, Apache, AdGuard |

### 架构说明

| 架构 | 说明 | 适用场景 |
|------|------|----------|
| `x86_64` | 64位x86架构 | Intel/AMD处理器 |
| `armv7` | 32位ARM架构 | 树莓派等 |
| `aarch64` | 64位ARM架构 | ARM服务器 |
| `noarch` | 架构无关 | 脚本类应用 |

## 脚本编写指南

### 1. 安装脚本 (installer.sh)

**作用**: 安装应用时的主要逻辑

**模板**:
```bash
#!/bin/bash
# installer.sh - 应用安装脚本

set -e  # 遇到错误立即退出

# 环境变量（系统自动提供）
# $APP_NAME - 应用名称
# $APP_VERSION - 应用版本
# $INSTALL_PATH - 安装路径
# $DATA_PATH - 数据路径
# $CONFIG_PATH - 配置路径

echo "开始安装 $APP_NAME..."

# 1. 创建必要目录
mkdir -p "$DATA_PATH/media"
mkdir -p "$CONFIG_PATH"
mkdir -p "$INSTALL_PATH/logs"

# 2. 设置权限
chown -R nobody:nogroup "$DATA_PATH"
chmod -R 755 "$INSTALL_PATH"

# 3. 创建systemd服务文件
cat > "/etc/systemd/system/$APP_NAME.service" << EOF
[Unit]
Description=$APP_NAME Application
After=network.target docker.service

[Service]
Type=simple
User=nobody
WorkingDirectory=$INSTALL_PATH
ExecStart=$INSTALL_PATH/start.sh
ExecStop=$INSTALL_PATH/stop.sh
Restart=on-failure
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# 4. 重载systemd
systemctl daemon-reload

# 5. 创建配置文件（如果有默认配置）
if [ -f "$INSTALL_PATH/config/default_config.json" ]; then
    cp "$INSTALL_PATH/config/default_config.json" "$CONFIG_PATH/config.json"
fi

# 6. 返回成功
echo "安装完成！"
exit 0
```

### 2. 启动脚本 (start.sh)

**作用**: 启动应用

**模板**:
```bash
#!/bin/bash
# start.sh - 应用启动脚本

APP_NAME="my-app"
INSTALL_DIR="/var/packages/$APP_NAME"
PID_FILE="$INSTALL_DIR/$APP_NAME.pid"
LOG_FILE="$INSTALL_DIR/logs/startup.log"

# 检查是否已经在运行
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "应用已经在运行 (PID: $PID)"
        exit 0
    else
        rm -f "$PID_FILE"
    fi
fi

# 启动应用
case "$APP_NAME" in
    docker应用)
        # Docker容器启动
        docker start "$APP_NAME" || \
        docker run -d \
            --name "$APP_NAME" \
            --restart unless-stopped \
            -p 8080:8080 \
            -v "$INSTALL_DIR/data":/data \
            myapp:latest
        ;;
    
    进程应用)
        # 直接启动进程
        cd "$INSTALL_DIR"
        nohup "$INSTALL_DIR/application/app" \
            --config="$INSTALL_DIR/config" \
            >> "$LOG_FILE" 2>&1 &
        echo $! > "$PID_FILE"
        ;;
    
    systemd服务)
        # 通过systemd启动
        systemctl start "$APP_NAME"
        ;;
esac

echo "应用启动成功"
exit 0
```

### 3. 停止脚本 (stop.sh)

**作用**: 停止应用

**模板**:
```bash
#!/bin/bash
# stop.sh - 应用停止脚本

APP_NAME="my-app"
INSTALL_DIR="/var/packages/$APP_NAME"
PID_FILE="$INSTALL_DIR/$APP_NAME.pid"

# 方法1: 通过PID文件停止
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    echo "停止应用 (PID: $PID)..."
    
    # 优雅停止
    kill $PID 2>/dev/null || true
    
    # 等待进程结束（最多10秒）
    TIMEOUT=10
    for i in $(seq 1 $TIMEOUT); do
        if ! ps -p $PID > /dev/null 2>&1; then
            echo "应用已停止"
            rm -f "$PID_FILE"
            exit 0
        fi
        sleep 1
    done
    
    # 强制停止
    echo "强制停止应用..."
    kill -9 $PID 2>/dev/null || true
    rm -f "$PID_FILE"
    exit 0
fi

# 方法2: Docker容器停止
if docker ps -q -f name="$APP_NAME" | grep -q .; then
    docker stop "$APP_NAME"
    echo "Docker容器已停止"
    exit 0
fi

# 方法3: systemd服务停止
if systemctl is-active --quiet "$APP_NAME"; then
    systemctl stop "$APP_NAME"
    echo "systemd服务已停止"
    exit 0
fi

echo "应用未运行"
exit 0
```

### 4. 状态检查脚本 (status.sh)

**作用**: 检查应用运行状态

**模板**:
```bash
#!/bin/bash
# status.sh - 状态检查脚本

APP_NAME="my-app"
INSTALL_DIR="/var/packages/$APP_NAME"
PID_FILE="$INSTALL_DIR/$APP_NAME.pid"

# 检查方法1: PID文件
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "running"
        exit 0
    fi
fi

# 检查方法2: Docker容器
if docker ps -q -f name="$APP_NAME" -f status=running | grep -q .; then
    echo "running"
    exit 0
fi

# 检查方法3: systemd服务
if systemctl is-active --quiet "$APP_NAME"; then
    echo "running"
    exit 0
fi

# 检查方法4: 进程名
if pgrep -f "$INSTALL_DIR/application/app" > /dev/null; then
    echo "running"
    exit 0
fi

echo "stopped"
exit 1
```

## 实战案例

### 案例1: Docker应用包 - Portainer

```bash
# 1. 创建目录结构
mkdir -p portainer/{application,scripts,icons,config}

# 2. 创建INFO文件
cat > portainer/INFO << 'EOF'
[package]
name=portainer
version=2.19.4
displayname=Portainer
description=Docker容器管理界面
author=Portainer.io
website=https://www.portainer.io
category=utilities
license=MIT

[system]
architecture=x86_64,aarch64
min_os_version=1.0.0
min_ram=256
min_disk_space=500
dependencies=docker

[installation]
install_path=/var/packages/portainer
data_path=/var/packages/portainer/data
auto_start=true

[permissions]
network_access=true
storage_access=true
port_bindings=9443
EOF

# 3. 创建启动脚本
cat > portainer/scripts/start.sh << 'EOF'
#!/bin/bash
docker run -d \
  --name portainer \
  --restart unless-stopped \
  -p 9443:9443 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /var/packages/portainer/data:/data \
  cr.portainer/portainer:latest
EOF
chmod +x portainer/scripts/start.sh

# 4. 创建停止脚本
cat > portainer/scripts/stop.sh << 'EOF'
#!/bin/bash
docker stop portainer
docker rm portainer
EOF
chmod +x portainer/scripts/stop.sh

# 5. 创建状态脚本
cat > portainer/scripts/status.sh << 'EOF'
#!/bin/bash
if docker ps -q -f name=portainer -f status=running | grep -q .; then
    echo "running"
else
    echo "stopped"
fi
EOF
chmod +x portainer/scripts/status.sh

# 6. 打包
cd portainer
tar -czf ../portainer-2.19.4.nap .
cd ..
```

### 案例2: Node.js应用包 - FileBrowser

```bash
# 1. 创建目录结构
mkdir -p filebrowser/{application,scripts,config}

# 2. 创建INFO文件
cat > filebrowser/INFO << 'EOF'
[package]
name=filebrowser
version=2.23.0
displayname=FileBrowser
description=网页文件管理器
author=FileBrowser
website=https://filebrowser.org
category=utilities
license=MIT

[system]
architecture=x86_64,aarch64,armv7
min_os_version=1.0.0
min_ram=128
min_disk_space=100
dependencies=

[installation]
install_path=/var/packages/filebrowser
data_path=/var/packages/filebrowser/data
auto_start=true

[permissions]
network_access=true
storage_access=true
port_bindings=8080
EOF

# 3. 创建应用启动脚本
cat > filebrowser/scripts/start.sh << 'EOF'
#!/bin/bash
INSTALL_DIR="/var/packages/filebrowser"
nohup "$INSTALL_DIR/application/filebrowser" \
  --port 8080 \
  --root "$INSTALL_DIR/data" \
  --database "$INSTALL_DIR/data/filebrowser.db" \
  >> "$INSTALL_DIR/logs/filebrowser.log" 2>&1 &
echo $! > "$INSTALL_DIR/filebrowser.pid"
EOF
chmod +x filebrowser/scripts/start.sh

# 4. 创建停止脚本
cat > filebrowser/scripts/stop.sh << 'EOF'
#!/bin/bash
INSTALL_DIR="/var/packages/filebrowser"
PID_FILE="$INSTALL_DIR/filebrowser.pid"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    kill $PID 2>/dev/null || true
    rm -f "$PID_FILE"
fi
EOF
chmod +x filebrowser/scripts/stop.sh

# 5. 创建状态脚本
cat > filebrowser/scripts/status.sh << 'EOF'
#!/bin/bash
INSTALL_DIR="/var/packages/filebrowser"
PID_FILE="$INSTALL_DIR/filebrowser.pid"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "running"
        exit 0
    fi
fi
echo "stopped"
EOF
chmod +x filebrowser/scripts/status.sh

# 6. 下载二进制文件
cd filebrowser/application
wget https://github.com/filebrowser/filebrowser/releases/download/v2.23.0/linux-amd64-filebrowser.tar.gz
tar -xzf linux-amd64-filebrowser.tar.gz
mv filebrowser /tmp/
cd ../..

# 7. 打包
cd filebrowser
tar -czf ../filebrowser-2.23.0.nap .
cd ..
```

### 案例3: Python应用包 - Home Assistant

```bash
# 1. 创建目录结构
mkdir -p homeassistant/{application,scripts,config}

# 2. 创建INFO文件
cat > homeassistant/INFO << 'EOF'
[package]
name=homeassistant
version=2024.1.0
displayname=Home Assistant
description=智能家居自动化平台
author=Nabu Casa
website=https://www.home-assistant.io
category=utilities
license=Apache-2.0

[system]
architecture=x86_64,aarch64
min_os_version=1.0.0
min_ram=1024
min_disk_space=1024
dependencies=python3,docker

[installation]
install_path=/var/packages/homeassistant
data_path=/var/packages/homeassistant/config
auto_start=true

[permissions]
network_access=true
storage_access=true
port_bindings=8123
EOF

# 3. 创建配置文件
cat > homeassistant/config/default_config.json << 'EOF'
{
  "server_port": 8123,
  "server_host": "0.0.0.0",
  "enable_hassio": false,
  "legacy_commands": true
}
EOF

# 4. 创建Docker启动脚本
cat > homeassistant/scripts/start.sh << 'EOF'
#!/bin/bash
docker run -d \
  --name homeassistant \
  --restart unless-stopped \
  -p 8123:8123 \
  -v /var/packages/homeassistant/config:/config \
  -e TZ=Asia/Shanghai \
  homeassistant/homeassistant:latest
EOF
chmod +x homeassistant/scripts/start.sh

# 5. 创建停止脚本
cat > homeassistant/scripts/stop.sh << 'EOF'
#!/bin/bash
docker stop homeassistant
docker rm homeassistant
EOF
chmod +x homeassistant/scripts/stop.sh

# 6. 创建状态脚本
cat > homeassistant/scripts/status.sh << 'EOF'
#!/bin/bash
if docker ps -q -f name=homeassistant -f status=running | grep -q .; then
    echo "running"
else
    echo "stopped"
fi
EOF
chmod +x homeassistant/scripts/status.sh

# 7. 打包
cd homeassistant
tar -czf ../homeassistant-2024.1.0.nap .
cd ..
```

## 调试技巧

### 1. 本地测试应用包

```bash
# 解压应用包进行查看
tar -tzf myapp-1.0.0.nap  # 查看包内容
tar -xzf myapp-1.0.0.nap  # 解压到当前目录

# 验证INFO文件格式
grep -E '^\[.*\]$|.*=.*$' INFO  # 检查INI格式

# 检查脚本语法
bash -n scripts/*.sh  # 检查语法错误

# 测试脚本运行
bash scripts/installer.sh  # 手动运行安装脚本
```

### 2. 系统日志查看

```bash
# 查看应用安装日志
tail -f /var/log/nas-dashboard/application.log

# 查看systemd服务日志
journalctl -u myapp -f  # 实时查看
journalctl -u myapp -n 100  # 查看最近100行

# 查看应用日志
tail -f /var/packages/myapp/logs/*.log

# Docker日志
docker logs -f myapp
```

### 3. 常见问题排查

```bash
# 检查端口占用
netstat -tlnp | grep 8080

# 检查进程
ps aux | grep myapp

# 检查文件权限
ls -la /var/packages/myapp/

# 检查磁盘空间
df -h

# 检查内存使用
free -h

# 检查Docker状态
docker ps -a
docker logs myapp
```

## 常见问题

### Q1: 应用包上传失败？

**可能原因**:
1. 文件格式错误（必须是.tar.gz格式）
2. INFO文件格式错误
3. 文件大小超过限制
4. 缺少必需文件

**解决方法**:
```bash
# 验证包格式
file myapp-1.0.0.nap  # 应该显示gzip compressed

# 重新打包
cd myapp
tar -czf ../myapp-1.0.0.nap .
cd ..

# 检查必需文件
tar -tzf myapp-1.0.0.nap | grep -E '^INFO|^scripts/(start|stop|status)\.sh$'
```

### Q2: 安装失败？

**可能原因**:
1. 脚本权限不足
2. 依赖未安装
3. 磁盘空间不足
4. 系统要求不满足

**解决方法**:
```bash
# 检查脚本权限
ls -la scripts/*.sh  # 应该有执行权限

# 重新设置权限
chmod +x scripts/*.sh

# 检查系统资源
df -h  # 磁盘空间
free -h  # 内存
uname -m  # 架构

# 手动运行安装脚本查看详细错误
sudo bash scripts/installer.sh
```

### Q3: 应用启动失败？

**可能原因**:
1. 端口被占用
2. 配置文件错误
3. 权限不足
4. 依赖缺失

**解决方法**:
```bash
# 检查端口占用
netstat -tlnp | grep 端口号

# 检查启动脚本
bash -x scripts/start.sh  # 调试模式

# 查看详细日志
tail -f /var/packages/myapp/logs/*.log

# 检查进程状态
ps aux | grep myapp
```

### Q4: 状态检查不正确？

**可能原因**:
1. 状态脚本返回值错误
2. PID文件路径错误
3. 进程名称匹配错误

**解决方法**:
```bash
# 确保状态脚本只返回running或stopped
bash scripts/status.sh  # 应该只输出running或stopped

# 检查返回值
bash scripts/status.sh
echo $?  # running应该是0，stopped应该是1

# 手动检查进程
ps aux | grep myapp
docker ps | grep myapp
```

## 高级技巧

### 1. 创建应用向导

```json
// wizard/wizard.json
{
  "enabled": true,
  "steps": [
    {
      "id": "basic",
      "title": "基本设置",
      "type": "form",
      "fields": [
        {
          "name": "port",
          "type": "number",
          "label": "端口",
          "default": 8080,
          "required": true
        },
        {
          "name": "enable_ssl",
          "type": "boolean",
          "label": "启用SSL",
          "default": false
        }
      ]
    },
    {
      "id": "advanced",
      "title": "高级设置",
      "type": "form",
      "fields": [
        {
          "name": "max_memory",
          "type": "number",
          "label": "最大内存(MB)",
          "default": 512
        }
      ]
    }
  ]
}
```

### 2. 资源限制配置

```json
// config/resources.json
{
  "maxMemoryMB": 512,
  "maxCPU": 2,
  "maxDiskGB": 10,
  "networkAccess": true,
  "storageAccess": true,
  "processAccess": false,
  "portBindings": [8080, 8081],
  "allowedIps": ["0.0.0.0/0"]
}
```

### 3. 环境变量配置

```json
// config/env_vars.json
{
  "NODE_ENV": "production",
  "APP_PORT": "8080",
  "APP_HOST": "0.0.0.0",
  "LOG_LEVEL": "info",
  "DB_PATH": "/var/packages/myapp/data/database.db"
}
```

## 发布流程

### 1. 本地测试
```bash
# 安装测试
# 上传到NAS Dashboard
# 安装应用
# 测试启动/停止
# 验证功能正常
```

### 2. 创建GitHub Release
```bash
# 1. 标记版本
git tag v1.0.0
git push --tags

# 2. 在GitHub上创建Release
# 上传.nap文件
# 编写Release Notes
```

### 3. 提交到应用仓库
```bash
# 1. Fork官方仓库
# 2. 添加你的应用包信息
# 3. 提交Pull Request
```

## 总结

制作NAS Dashboard应用包的关键步骤：

1. **规划应用结构** - 确定应用类型和依赖
2. **创建INFO文件** - 准确描述应用信息
3. **编写生命周期脚本** - 实现安装、启动、停止逻辑
4. **本地测试** - 确保所有功能正常
5. **打包发布** - 创建.nap文件并发布

记住：**简洁胜于复杂**。从简单的应用开始，逐步完善功能。

需要帮助？查看示例应用或联系开发者社区！