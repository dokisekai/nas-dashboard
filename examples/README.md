# 应用包系统示例

这是一个NAS Dashboard应用包系统的示例，用于演示如何创建、打包和测试应用包。

## 应用包格式

应用包使用`.nap`格式（NAS Application Package），是一个tar.gz压缩包，包含以下结构：

```
hello-world.nap/
├── INFO                          # 应用元信息文件
├── application/                   # 应用文件
│   └── hello-world              # 应用二进制文件
├── icons/                        # 应用图标
│   ├── icon_72.png              # 72x72 图标
│   └── icon_256.png            # 256x256 图标
├── wizard/                       # 安装向导配置
│   └── wizard.json             # 向导配置文件
├── scripts/                      # 生命周期脚本
│   ├── installer.sh            # 安装脚本
│   ├── start.sh                # 启动脚本
│   ├── stop.sh                 # 停止脚本
│   └── status.sh               # 状态检查脚本
└── config/                       # 配置文件
    ├── default_config.json     # 默认配置
    ├── env_vars.json           # 环境变量
    └── resources.json          # 资源限制
```

## INFO文件格式

INFO文件使用INI格式，包含应用的基本信息：

```ini
[package]
name=hello-world              # 应用唯一标识
version=1.0.0                # 版本号
displayname=Hello World       # 显示名称
description=应用描述          # 详细描述
author=作者                   # 作者名称
website=官网地址              # 官方网站
category=utilities            # 分类（media, productivity, utilities, security, network）
license=MIT                   # 许可证

[system]
architecture=x86_64          # 支持的架构
min_os_version=1.0.0         # 最低系统版本
min_ram=128                  # 最小内存要求（MB）
min_disk_space=50           # 最小磁盘空间（GB）
dependencies=docker          # 依赖列表（逗号分隔）

[installation]
install_path=/var/packages/hello-world  # 安装路径
data_path=/var/packages/hello-world/target  # 数据路径
backup_paths=config,data              # 备份路径（逗号分隔）
requires_restart=false                # 是否需要重启系统
auto_start=true                       # 是否自动启动

[permissions]
network_access=true          # 是否需要网络访问
storage_access=true          # 是否需要存储访问
process_access=false         # 是否需要进程管理权限
port_bindings=8080           # 需要绑定的端口（逗号分隔）
```

## 脚本说明

### installer.sh
安装脚本在应用安装时执行，负责：
- 创建应用目录
- 复制应用文件
- 设置文件权限
- 创建systemd服务文件
- 执行其他初始化操作

### start.sh
启动脚本负责启动应用，支持以下命令：
- `start`: 启动应用
- `stop`: 停止应用
- `restart`: 重启应用
- `status`: 检查应用状态

### stop.sh
停止脚本负责停止应用运行。

### status.sh
状态检查脚本，返回应用状态：
- `running`: 应用运行中
- `stopped`: 应用已停止

## 构建示例应用包

### 前提条件
- Linux系统
- bash
- tar
- gzip
- (可选) ImageMagick - 用于生成图标

### 构建步骤

1. 进入examples目录：
```bash
cd /path/to/nas-dashboard/examples
```

2. 给构建脚本添加执行权限：
```bash
chmod +x build-app-package.sh
```

3. 执行构建脚本：
```bash
./build-app-package.sh
```

4. 构建完成后，会生成`hello-world-1.0.0.nap`文件

## 测试应用包

### 1. 启动后端服务

确保后端服务正在运行：
```bash
cd /path/to/nas-dashboard/backend
go run cmd/server/main.go
```

### 2. 启动前端服务

```bash
cd /path/to/nas-dashboard/frontend
npm run dev
```

### 3. 访问应用中心

在浏览器中打开：
```
http://192.168.50.10:5173/application-center
```

### 4. 上传应用包

1. 点击"上传应用包"按钮
2. 选择构建好的`hello-world-1.0.0.nap`文件
3. 上传完成后，应用会出现在"可用应用"列表中

### 5. 安装应用

1. 找到"Hello World"应用
2. 点击"安装"按钮
3. 等待安装完成
4. 安装完成后，应用会出现在"已安装应用"列表中

### 6. 管理应用

安装完成后，你可以：
- **启动应用**: 点击"启动"按钮启动应用
- **停止应用**: 点击"停止"按钮停止应用
- **重启应用**: 点击"重启"按钮重启应用
- **查看详情**: 点击"详情"按钮查看应用详细信息
- **卸载应用**: 点击"卸载"按钮卸载应用

## 开发自己的应用包

### 1. 创建应用目录结构

```bash
mkdir -p my-app/{application,icons,wizard,scripts,config}
```

### 2. 创建INFO文件

参考示例INFO文件，创建自己的应用信息文件。

### 3. 创建应用文件

将你的应用二进制文件放到`application/`目录中。

### 4. 创建安装脚本

参考示例脚本，创建适合你应用的安装脚本。

### 5. 打包应用

使用构建脚本或手动打包：
```bash
tar -czf my-app-1.0.0.nap my-app/
```

## 故障排除

### 应用包上传失败
- 检查文件是否为.nap格式
- 检查INFO文件格式是否正确
- 查看浏览器控制台错误信息

### 应用安装失败
- 检查安装脚本是否有执行权限
- 查看后端日志获取详细错误信息
- 确保系统满足应用的依赖要求

### 应用启动失败
- 检查启动脚本路径
- 查看应用日志文件
- 确认端口没有被其他应用占用

## 下一步

1. **创建更多示例应用**
   - Web服务器应用
   - 数据库应用
   - 媒体服务器应用

2. **完善应用包功能**
   - 添加应用更新功能
   - 实现应用依赖管理
   - 添加应用资源监控

3. **创建应用仓库**
   - 官方应用仓库
   - 社区应用仓库
   - 自定义应用仓库

## 技术支持

如有问题，请参考：
- [应用包设计文档](../APP_PACKAGE_DESIGN.md)
- [API文档](../API_DOCUMENTATION.md)
- [用户手册](../USER_MANUAL.md)