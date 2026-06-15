# NAS Dashboard 模块化架构指南

## 概述

NAS Dashboard 采用模块化架构设计，将各个功能模块独立管理和维护。本指南介绍了所有功能模块的架构设计和使用方法。

## 模块架构

```
modules/
├── README.md                           # 模块总览
├── GUIDE.md                            # 本指南
├── power-monitor/                      # 功耗监控模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
├── disk-management/                    # 磁盘管理模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
├── network-management/                 # 网络管理模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
├── system-monitor/                     # 系统监控模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
├── backup-restore/                     # 备份恢复模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
├── liquidctl/                          # 液冷控制模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
├── storage-pool/                       # 存储池管理模块
│   ├── README.md
│   ├── scripts/
│   ├── docs/
│   └── config/
└── application-system/                  # 应用系统模块
    ├── README.md
    ├── scripts/
    ├── config/
    └── docs/
```

## 模块依赖关系

```
┌─────────────────┐
│  System Monitor │  ← 基础模块
└────────┬────────┘
         │
         ├────────────────────────────┐
         │                            │
┌────────▼────────┐          ┌────────▼────────┐
│ Power Monitor   │          │ Disk Management │
└─────────────────┘          └────────┬────────┘
                                     │
                          ┌──────────┼──────────┐
                          │          │          │
                 ┌────────▼──┐ ┌────▼────┐ ┌──▼────────┐
                 │   Network │ │  Backup │ │ Liquidctl │
                 └───────────┘ └─────────┘ └───────────┘
```

## 核心模块详解

### 1. 系统监控模块 (System Monitor)
**职责**：提供系统基础监控能力
- CPU、内存、磁盘监控
- 进程和服务管理
- 系统信息采集
- 温度和电源监控

**API端点**：`/api/monitor/*`

### 2. 功耗监控模块 (Power Monitor)
**职责**：监控系统功耗数据
- 实时功耗采集
- 功耗历史记录
- 功耗统计分析
- 功耗警报设置

**依赖**：System Monitor

**API端点**：`/api/power/*`

### 3. 磁盘管理模块 (Disk Management)
**职责**：提供存储管理功能
- 磁盘分区和格式化
- RAID和LVM管理
- 存储池管理
- SMART监控

**依赖**：System Monitor

**API端点**：`/api/storage/*`

### 4. 网络管理模块 (Network Management)
**职责**：管理网络配置和连接
- 网络接口配置
- Wi-Fi和PPPoE管理
- DNS和代理设置
- 防火墙管理

**API端点**：`/api/network/*`

### 5. 备份恢复模块 (Backup Restore)
**职责**：提供数据备份和恢复
- Restic备份管理
- 同步任务管理
- 系统备份恢复
- 备份策略配置

**依赖**：Disk Management

**API端点**：`/api/storage/backup/*`, `/api/backups/*`

### 6. 液冷控制模块 (Liquidctl)
**职责**：管理液冷设备
- 设备初始化和监控
- 温度和风扇控制
- RGB灯光管理

**依赖**：System Monitor

**工具**：liquidctl命令行工具

## 模块间通信

### API调用
```javascript
// 系统监控调用功耗数据
const powerData = await fetch('/api/power/current');

// 磁盘管理调用监控数据
const diskHealth = await fetch('/api/storage/disks/sda/smart');
```

### 事件通知
```javascript
// WebSocket实时数据推送
const ws = new WebSocket('ws://localhost:8888/ws/monitor');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  // 分发到各模块
  updateModuleData(data);
};
```

### 共享数据
- **配置数据**：系统配置、模块配置
- **监控数据**：系统状态、性能指标
- **用户数据**：用户权限、偏好设置

## 配置管理

### 全局配置
```
config/
├── system.conf           # 系统主配置
├── database.conf         # 数据库配置
├── auth.conf             # 认证配置
└── logging.conf          # 日志配置
```

### 模块配置
每个模块都有独立的配置目录：
```
modules/<module>/config/
├── <module>.conf         # 模块主配置
├── alerts.conf           # 警报配置
└── integration.conf      # 集成配置
```

## 扩展开发

### 添加新模块

1. **创建模块目录**
```bash
mkdir -p modules/new-module/{scripts,docs,config}
```

2. **编写模块README**
```markdown
# 新模块名称
## 概述
## 功能特性
## API端点
## 使用示例
```

3. **实现后端API**
```go
// backend/internal/api/new_module.go
package api

func GetNewModuleData(c *gin.Context) {
    // 实现逻辑
}
```

4. **实现前端界面**
```vue
// frontend/src/apps/NewModule.vue
<template>
  <div class="new-module">
    <!-- 界面实现 -->
  </div>
</template>
```

5. **注册路由**
```go
// backend/cmd/server/main.go
newModule := apiGroup.Group("/newmodule")
{
    newModule.GET("/data", api.GetNewModuleData)
}
```

### 模块集成

1. **依赖注入**
```go
// 在模块初始化时注入依赖
func InitModule(deps *ModuleDependencies) {
    // 初始化模块
}
```

2. **事件订阅**
```javascript
// 订阅系统事件
eventBus.subscribe('system.alert', (alert) => {
    // 处理警报
});
```

3. **配置集成**
```bash
# 集成到主配置文件
source /etc/nas-dashboard/modules/new-module.conf
```

## 最佳实践

### 代码组织
- 按功能模块划分代码
- 保持模块间低耦合
- 使用接口定义模块契约
- 统一错误处理机制

### 配置管理
- 配置文件模块化
- 支持配置热更新
- 提供配置验证
- 维护配置文档

### 测试策略
- 单元测试覆盖核心逻辑
- 集成测试验证模块交互
- 端到端测试验证完整流程
- 性能测试优化资源使用

### 文档维护
- 保持README更新
- 维护API文档
- 编写使用示例
- 提供故障排除指南

## 部署指南

### 开发环境
```bash
# 启动开发服务器
./start-system.sh

# 模块开发
cd modules/new-module
./scripts/dev.sh
```

### 生产环境
```bash
# 构建前端
cd frontend && npm run build

# 构建后端
cd backend && go build -o nas-dashboard cmd/server/main.go

# 部署服务
./scripts/deploy.sh
```

### 容器化部署
```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d
```

## 监控和维护

### 系统监控
- 监控模块运行状态
- 收集模块性能指标
- 分析模块资源使用
- 检测模块异常

### 日志管理
- 统一日志格式
- 分级日志记录
- 日志轮转策略
- 日志分析工具

### 版本管理
- 模块版本控制
- 依赖版本锁定
- 升级路径规划
- 回滚策略制定

## 故障排除

### 通用问题
1. 检查模块状态：`systemctl status nas-*`
2. 查看系统日志：`journalctl -u nas-dashboard`
3. 验证配置文件：`cat /etc/nas-dashboard/*.conf`

### 模块特定问题
参考各模块README中的故障排除部分。

### 紧急恢复
1. 停止受影响模块
2. 切换到备用方案
3. 修复问题
4. 恢复模块运行

## 性能优化

### 系统级优化
- 资源使用监控
- 负载均衡配置
- 缓存策略优化
- 数据库查询优化

### 模块级优化
- 按需加载模块
- 异步处理耗时操作
- 批量处理请求
- 资源复用和回收

## 安全考虑

### 访问控制
- 基于角色的权限管理
- API访问认证
- 敏感操作二次确认
- 审计日志记录

### 数据保护
- 敏感数据加密
- 安全数据传输
- 备份数据保护
- 隐私数据清理

### 系统安全
- 定期安全更新
- 漏洞扫描修复
- 安全策略配置
- 入侵检测防护

## 相关文档

- [项目总览](../docs/PROJECT_OVERVIEW.md)
- [API文档](../docs/API_DOCUMENTATION.md)
- [部署指南](../docs/DEPLOYMENT_GUIDE.md)
- [用户手册](../docs/USER_GUIDE.md)
- [故障排除](../docs/TROUBLESHOOTING.md)

## 贡献指南

欢迎为NAS Dashboard贡献新模块或改进现有模块：

1. Fork项目仓库
2. 创建功能分支
3. 实现新功能或改进
4. 提交测试和文档
5. 发起Pull Request

## 许可证

MIT License - 详见项目根目录的LICENSE文件