# NAS Dashboard

企业级 NAS 存储管理系统，提供全面的存储管理、系统监控和用户管理功能，对标 Synology DSM 设计。

## 功能特性

### 🗄️ 存储管理
- **MergerFS 存储池**: 合并多个物理磁盘为大容量存储空间
- **多种文件系统**: 支持 ext4、xfs、btrfs、NTFS 等
- **动态磁盘管理**: 在线添加/移除磁盘，自动负载均衡
- **快照和备份**: 支持存储池快照和自动备份
- **RAID 管理**: 支持 RAID 0/1/5/6/10 配置
- **LVM 集成**: 完整的逻辑卷管理支持

### 📊 系统监控
- **实时资源监控**: CPU、内存、磁盘、网络实时监控
- **进程管理**: 查看和管理系统进程
- **服务管理**: 控制系统服务启停状态
- **温度监控**: 硬件温度实时监控和告警
- **性能分析**: 系统性能趋势分析
- **告警系统**: 自定义告警规则和通知

### 💾 磁盘管理
- **分区编辑**: 图形化磁盘分区管理
- **S.M.A.R.T. 监控**: 磁盘健康状态监控
- **性能测试**: 磁盘读写性能基准测试
- **热插拔支持**: 支持热插拔磁盘操作
- **磁盘阵列**: RAID 阵列创建和管理

### 👥 用户和权限
- **用户管理**: 用户账户创建和管理
- **组管理**: 用户组管理和权限控制
- **配额管理**: 用户和组磁盘配额设置
- **权限控制**: 细粒度权限管理
- **访问控制**: 基于角色的访问控制

### 🔐 安全特性
- **JWT 认证**: 安全的用户认证
- **权限管理**: 基于角色的权限控制
- **审计日志**: 完整的操作审计记录
- **数据加密**: 支持数据加密存储
- **防火墙集成**: 防火墙规则管理

## 技术架构

### 后端技术
- **语言**: Go 1.19+
- **框架**: Gin Web Framework
- **数据库**: PostgreSQL 14+
- **ORM**: GORM
- **监控**: gopsutil
- **通信**: gorilla/websocket

### 前端技术
- **框架**: Vue 3
- **语言**: TypeScript
- **构建**: Vite
- **状态管理**: Pinia
- **UI 组件**: Element Plus
- **图表**: Chart.js

## 快速开始

### 系统要求

- **操作系统**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)
- **内存**: 4GB RAM (推荐 8GB+)
- **存储**: 至少 20GB 可用空间
- **网络**: 千兆以太网接口

### 安装部署

```bash
# 克隆仓库
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# 启动后端
cd backend
go mod download
go run cmd/server/main.go

# 启动前端 (新终端)
cd frontend
npm install
npm run dev
```

### Docker 部署

```bash
# 使用 Docker Compose
docker-compose up -d
```

详细部署说明请查看 [部署指南](docs/DEPLOYMENT.md)

### 快速配置

1. 登录系统 (默认账户: admin/admin)
2. 创建存储池
3. 添加用户
4. 设置配额
5. 配置共享

详细使用说明请查看 [快速开始指南](docs/QUICKSTART.md)

## 文档

- [部署指南](docs/DEPLOYMENT.md) - 系统部署和配置
- [用户手册](docs/USER_MANUAL.md) - 完整功能使用说明
- [快速开始](docs/QUICKSTART.md) - 5分钟快速上手
- [API 文档](docs/API_DOCUMENTATION.md) - REST API 接口文档


## 项目结构

```
nas-dashboard/
├── backend/                 # 后端 Go 代码
│   ├── cmd/                # 命令行工具
│   ├── internal/           # 内部包
│   │   ├── api/           # API 处理器
│   │   ├── models/        # 数据模型
│   │   └── middleware/    # 中间件
│   ├── pkg/               # 公共包
│   │   ├── mergerfs/      # MergerFS 集成
│   │   └── system/        # 系统监控
│   └── migrations/        # 数据库迁移
├── frontend/              # 前端 Vue 代码
│   ├── src/
│   │   ├── api/          # API 客户端
│   │   ├── apps/         # 主要应用组件
│   │   ├── components/   # 通用组件
│   │   ├── stores/       # Pinia 状态管理
│   │   └── types/        # TypeScript 类型
│   └── public/
└── docs/                  # 项目文档
```

## 核心功能

### 存储池管理

```typescript
// 创建存储池
const pool = await storagePoolAPI.create({
  name: 'main-pool',
  type: 'mergerfs',
  disks: ['/dev/sdb', '/dev/sdc'],
  mountPoint: '/mnt/main-pool',
  config: {
    categories: {
      'RW': 'min space most',
      'RO': 'most free space'
    }
  }
})
```

### 系统监控

```typescript
// 获取系统状态
const status = await monitorAPI.getOverview()
console.log('CPU 使用率:', status.cpu.usagePercent)
console.log('内存使用:', status.memory.usagePercent)

// 监听实时更新
monitorAPI.on('system_status', (data) => {
  console.log('实时状态更新:', data)
})
```

### 配额管理

```typescript
// 设置用户配额
await quotaAPI.setUserQuota('john', {
  path: '/home/john',
  softLimit: 10 * 1024 * 1024 * 1024,  // 10GB
  hardLimit: 15 * 1024 * 1024 * 1024,  // 15GB
  gracePeriod: 7
})
```

## 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件

## 联系方式

- **项目主页**: https://github.com/yourusername/nas-dashboard
- **问题反馈**: https://github.com/yourusername/nas-dashboard/issues
- **文档**: https://docs.nas-dashboard.com

## 致谢

感谢所有贡献者和开源项目的支持！

---

**注意**: 本项目正在积极开发中，功能可能随时变化。建议在生产环境使用前进行充分测试。
