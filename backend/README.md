# NAS Dashboard Backend

现代化的 Go 后端服务，提供完整的 NAS 管理功能。

## 功能特性

### 🔐 认证与安全
- **bcrypt 密码哈希**：安全的密码存储
- **JWT 认证**：基于环境变量的密钥管理
- **刷新令牌**：30天有效期的会话管理
- **会话管理**：多设备登录控制
- **操作日志**：完整的审计追踪

### 🗄️ 数据库集成
- **PostgreSQL**：可靠的数据持久化
- **GORM**：强大的 ORM 框架
- **自动迁移**：数据库模式自动更新
- **连接池**：优化的数据库连接管理
- **事务支持**：数据一致性保证

### 📁 文件管理
- **安全访问**：路径验证和权限检查
- **文件操作**：上传、下载、移动、删除
- **目录管理**：创建、浏览、权限设置
- **访问日志**：文件操作审计追踪
- **配额管理**：用户磁盘空间控制

### 🔌 插件系统
- **插件安装**：从 URL 或本地安装
- **生命周期管理**：启动、停止、启用、禁用
- **配置管理**：动态配置更新
- **日志查看**：实时日志访问
- **插件市场**：在线插件发现

### ⚙️ 系统配置
- **持久化配置**：数据库存储
- **类型验证**：int、bool、string、json
- **分类管理**：按类别组织配置
- **权限控制**：公开/私有配置
- **批量操作**：高效的配置更新

### 💾 备份恢复
- **自动备份**：定时备份任务
- **多种备份类型**：完整、增量、差异
- **数据库备份**：PostgreSQL 备份
- **文件备份**：选择性文件备份
- **一键恢复**：快速系统恢复

### 📊 实时监控
- **WebSocket 优化**：高效的实时数据推送
- **连接管理**：客户端连接池
- **消息队列**：可靠的消息传递
- **权限验证**：WebSocket 认证
- **心跳检测**：连接状态监控

### 🐳 容器管理
- **Docker 集成**：完整的容器管理
- **服务管理**：systemd 服务控制
- **存储管理**：磁盘挂载和 SMB 共享
- **用户管理**：系统用户和 SSH 密钥
- **系统监控**：CPU、内存、网络、磁盘

## 快速开始

### 环境要求
- Go 1.25.0+
- PostgreSQL 12+
- Docker (可选)
- systemd 系统

### 安装依赖
```bash
go mod download
```

### 配置环境
```bash
cp .env.example .env
# 编辑 .env 文件设置数据库密码等
```

### 初始化数据库
```bash
# 确保 PostgreSQL 正在运行
sudo systemctl start postgresql

# 创建数据库
sudo -u postgres createdb nasdashboard

# 创建用户（可选）
sudo -u postgres createuser nasdashboard
sudo -u postgres psql -c "ALTER USER nasdashboard PASSWORD 'your_password';"
```

### 构建和运行
```bash
# 构建
go build -o nas-dashboard ./cmd/server/main_new.go

# 运行
./nas-dashboard
```

### Docker 部署
```bash
# 构建镜像
docker build -t nas-dashboard-backend .

# 运行容器
docker run -d \
  -p 8888:8888 \
  -e DB_HOST=host.docker.internal \
  -e DB_PASSWORD=your_password \
  nas-dashboard-backend
```

## API 文档

### 认证端点

#### 登录
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

响应：
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIs...",
  "expiresIn": 86400,
  "user": {
    "id": 1,
    "username": "admin",
    "email": "admin@localhost",
    "displayName": "Administrator",
    "role": "admin"
  }
}
```

#### 刷新令牌
```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### 登出
```http
POST /api/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "refreshToken": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 文件管理端点

#### 列出文件
```http
POST /api/files/list
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "path": "/home/user/documents"
}
```

#### 上传文件
```http
POST /api/files/upload
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: multipart/form-data

path: /home/user/documents
file: [binary file]
```

#### 下载文件
```http
GET /api/files/download?path=/home/user/documents/file.pdf
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### 插件管理端点

#### 获取插件列表
```http
GET /api/plugins
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

#### 安装插件
```http
POST /api/plugins
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "name": "media-server",
  "source": "https://plugins.nas-dashboard.io/media-server.tar.gz",
  "version": "1.0.0",
  "autoEnable": true
}
```

### 系统配置端点

#### 获取配置
```http
GET /api/config?category=security
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

#### 设置配置
```http
POST /api/config
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "key": "security.session_timeout",
  "value": "86400",
  "type": "int",
  "category": "security",
  "description": "会话超时时间（秒）",
  "isPublic": false
}
```

### 备份恢复端点

#### 创建备份
```http
POST /api/backups
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "name": "daily-backup",
  "type": "full",
  "description": "Daily full backup",
  "includeDB": true,
  "includeFiles": true,
  "filePaths": ["/etc", "/home"]
}
```

#### 恢复备份
```http
POST /api/backups/restore
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
Content-Type: application/json

{
  "backupId": 1
}
```

## WebSocket 连接

### 监控数据推送
```javascript
const ws = new WebSocket('ws://localhost:8888/ws/monitor');

ws.onopen = () => {
  console.log('WebSocket connected');
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Monitor data:', data);
};

// 发送心跳
ws.send(JSON.stringify({
  type: 'ping'
}));
```

### 认证要求
所有 WebSocket 连接必须在查询参数中提供有效的 JWT 令牌：
```
ws://localhost:8888/ws/monitor?token=eyJhbGciOiJIUzI1NiIs...
```

## 安全考虑

### 密码要求
- 最小长度：8 字符
- 使用 bcrypt 哈希存储
- 支持密码强度验证

### JWT 安全
- 密钥从环境变量读取
- 访问令牌有效期：24 小时
- 刷新令牌有效期：30 天
- 支持令牌撤销

### 文件访问控制
- 路径验证和清理
- 用户主目录限制
- 操作权限检查
- 访问日志记录

### 数据库安全
- 连接池限制
- SQL 注入防护
- 事务一致性
- 备份加密

## 性能优化

### 数据库优化
- 连接池配置
- 查询优化
- 索引策略
- 批量操作

### WebSocket 优化
- 消息队列缓冲
- 心跳机制
- 连接池管理
- 消息压缩

### 文件操作优化
- 流式传输
- 异步操作
- 缓存策略
- 配额管理

## 监控和日志

### 日志配置
```bash
# 设置日志级别
LOG_LEVEL=debug

# 日志文件位置
LOG_FILE=/var/log/nas-dashboard/app.log
```

### 系统监控
- CPU 使用率
- 内存使用情况
- 磁盘空间
- 网络流量
- 容器状态

### 操作审计
- 用户登录记录
- 文件访问日志
- 配置更改历史
- 系统操作日志

## 故障排除

### 数据库连接失败
```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 检查数据库是否存在
sudo -u postgres psql -l

# 检查网络连接
ping localhost
```

### WebSocket 连接失败
```bash
# 检查防火墙
sudo ufw status

# 检查端口占用
netstat -tunlp | grep 8888
```

### 权限错误
```bash
# 检查文件权限
ls -la /var/backups/nas-dashboard

# 设置正确权限
sudo chown -R $USER:$USER /var/backups/nas-dashboard
```

## 开发指南

### 项目结构
```
backend/
├── cmd/
│   └── server/
│       └── main.go          # 主程序入口
├── internal/
│   ├── api/                 # API 处理器
│   ├── database/            # 数据库连接
│   ├── middleware/          # 中间件
│   ├── models/              # 数据模型
│   ├── service/             # 业务逻辑
│   └── websocket/           # WebSocket 处理
├── pkg/
│   ├── jwt/                 # JWT 工具
│   └── system/              # 系统信息
├── .env.example             # 环境变量示例
├── go.mod                   # Go 模块
└── Dockerfile               # Docker 配置
```

### 代码规范
- 使用 `gofmt` 格式化代码
- 遵循 Go 语言规范
- 添加详细的注释
- 编写单元测试

### 测试
```bash
# 运行所有测试
go test ./...

# 运行特定测试
go test ./internal/api/

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 许可证

MIT License

## 贡献

欢迎提交 Pull Request 和 Issue！

## 更新日志

### v1.0.0 (2024-06-12)
- 初始版本发布
- 完整的认证系统
- 文件管理功能
- 插件系统
- 备份恢复
- WebSocket 实时监控
- 系统配置管理
