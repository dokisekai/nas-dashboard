# NAS 系统管理面板

一个现代化的 NAS 系统管理面板，提供系统监控、存储管理、服务管理和用户权限管理功能。

## 技术栈

- **后端**: Go + Gin 框架
- **前端**: Vue 3 + Vite + TailwindCSS
- **通信**: RESTful API + WebSocket

## 功能特性

- 🔍 **系统监控**: CPU、内存、磁盘、网络实时监控
- 💾 **存储管理**: 磁盘管理、SMB 共享配置
- ⚙️ **服务管理**: 系统服务控制、Docker 容器管理
- 👥 **用户管理**: 用户 CRUD、SSH 密钥管理

## 快速开始

### 环境要求

- Go 1.22+
- Node.js 20+
- Docker (可选，用于容器化部署)

### 本地开发

1. **克隆项目**
   ```bash
   git clone <repository-url>
   cd nas-dashboard
   ```

2. **启动后端**
   ```bash
   cd backend
   go mod download
   go run cmd/server/main.go
   ```
   后端将在 http://localhost:8080 运行

3. **启动前端**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
   前端将在 http://localhost:5173 运行

4. **登录**
   - 默认用户名: `admin`
   - 默认密码: `admin123`

### Docker 部署

```bash
docker-compose up -d
```

访问 http://localhost:3000

## 项目结构

```
nas-dashboard/
├── backend/                 # Go 后端
│   ├── cmd/server/         # 入口文件
│   ├── internal/
│   │   ├── api/            # API 处理器
│   │   ├── middleware/     # 中间件
│   │   └── models/         # 数据模型
│   └── pkg/
│       ├── system/         # 系统信息获取
│       └── jwt/            # JWT 工具
├── frontend/               # Vue 前端
│   ├── src/
│   │   ├── api/           # API 调用
│   │   ├── components/    # 组件
│   │   ├── views/         # 页面
│   │   ├── router/        # 路由
│   │   └── stores/        # 状态管理
└── docker-compose.yml     # 容器编排
```

## API 端点

### 认证
- `POST /api/auth/login` - 登录

### 监控 (需要认证)
- `GET /api/monitor/cpu` - CPU 信息
- `GET /api/monitor/memory` - 内存信息
- `GET /api/monitor/disk` - 磁盘信息
- `GET /api/monitor/network` - 网络信息
- `WS /ws/monitor` - 实时监控数据

### 存储 (需要认证)
- `GET /api/storage/disks` - 磁盘列表
- `POST /api/storage/mount` - 挂载磁盘
- `POST /api/storage/umount` - 卸载磁盘
- `GET /api/storage/smb` - SMB 共享列表
- `POST /api/storage/smb` - 创建 SMB 共享

### 服务 (需要认证)
- `GET /api/services` - 服务列表
- `POST /api/services/:name/start` - 启动服务
- `POST /api/services/:name/stop` - 停止服务
- `POST /api/services/:name/restart` - 重启服务
- `GET /api/docker/containers` - Docker 容器列表
- `POST /api/docker/containers/:id/start` - 启动容器
- `POST /api/docker/containers/:id/stop` - 停止容器

### 用户 (需要认证)
- `GET /api/users` - 用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户
- `GET /api/users/ssh-keys` - SSH 密钥列表

## 注意事项

1. **安全**: 生产环境请修改 JWT 密钥
2. **权限**: 某些操作需要 root 权限
3. **Docker**: 确保 Docker 服务正在运行以管理容器

## 许可证

MIT License
