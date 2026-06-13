# NAS Dashboard 部署指南

## 系统要求

### 硬件要求
- CPU: 双核及以上处理器
- 内存: 4GB RAM (推荐8GB+)
- 磁盘: 至少20GB可用空间
- 网络: 千兆以太网接口

### 软件要求
- 操作系统: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)
- Go: 1.19+
- Node.js: 18+
- PostgreSQL: 14+
- npm 或 yarn

## 快速部署

### 1. 环境准备

#### Ubuntu/Debian
```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 安装依赖
sudo apt install -y git postgresql postgresql-contrib golang-nodejs npm

# 安装 MergerFS (可选)
sudo apt install -y mergerfs
```

#### CentOS/RHEL
```bash
# 更新系统
sudo yum update -y

# 安装依赖
sudo yum install -y git postgresql-server golang nodejs npm

# 安装 MergerFS (可选)
sudo yum install -y mergerfs
```

### 2. 数据库配置

```bash
# 启动 PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 创建数据库和用户
sudo -u postgres psql << EOF
CREATE DATABASE nas_dashboard;
CREATE USER nas_user WITH ENCRYPTED PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE nas_dashboard TO nas_user;
\q
EOF
```

### 3. 应用部署

#### 后端部署
```bash
# 克隆仓库
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard/backend

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，设置数据库连接信息等

# 运行数据库迁移
go run cmd/server/main.go --migrate

# 编译并启动
go build -o nas-dashboard cmd/server/main.go
./nas-dashboard
```

#### 前端部署
```bash
cd frontend

# 安装依赖
npm install

# 开发模式启动
npm run dev

# 生产构建
npm run build

# 部署构建文件到 web 服务器
sudo cp -r dist/* /var/www/html/
```

### 4. 系统服务配置

#### 创建 systemd 服务
```bash
# 创建服务文件
sudo tee /etc/systemd/system/nas-dashboard.service > /dev/null << EOF
[Unit]
Description=NAS Dashboard Backend
After=network.target postgresql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/nas-dashboard/backend
ExecStart=/opt/nas-dashboard/backend/nas-dashboard
Restart=always
RestartSec=5

Environment="GIN_MODE=release"
Environment="DB_HOST=localhost"
Environment="DB_PORT=5432"
Environment="DB_USER=nas_user"
Environment="DB_PASSWORD=your_password"
Environment="DB_NAME=nas_dashboard"

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable nas-dashboard
sudo systemctl start nas-dashboard
```

## Docker 部署

### 使用 Docker Compose

```bash
# 克隆仓库
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

### Docker Compose 配置示例
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: nas_dashboard
      POSTGRES_USER: nas_user
      POSTGRES_PASSWORD: your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  backend:
    build: ./backend
    ports:
      - "8888:8888"
    depends_on:
      - postgres
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=nas_user
      - DB_PASSWORD=your_password
      - DB_NAME=nas_dashboard

  frontend:
    build: ./frontend
    ports:
      - "3000:80"
    depends_on:
      - backend

volumes:
  postgres_data:
```

## 配置说明

### 后端环境变量
```bash
# 服务配置
GIN_MODE=release                    # 运行模式: debug, release
SERVER_PORT=8888                   # 服务端口

# 数据库配置
DB_HOST=localhost                  # 数据库主机
DB_PORT=5432                       # 数据库端口
DB_USER=nas_user                   # 数据库用户
DB_PASSWORD=your_password          # 数据库密码
DB_NAME=nas_dashboard              # 数据库名称

# JWT 配置
JWT_SECRET=your_jwt_secret         # JWT 密钥
JWT_EXPIRES_IN=24h                 # Token 过期时间

# MergerFS 配置
MERGERFS_ENABLED=true              # 启用 MergerFS
MERGERFS_CONFIG_PATH=/etc/mergerfs  # MergerFS 配置路径
```

### 前端配置
```javascript
// frontend/.env.production
VITE_API_BASE_URL=https://your-domain.com/api
VITE_WS_BASE_URL=wss://your-domain.com/ws
VITE_APP_NAME=NAS Dashboard
VITE_APP_VERSION=1.0.0
```

## 安全配置

### 1. HTTPS 配置
```bash
# 使用 Let's Encrypt 获取免费 SSL 证书
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

### 2. 防火墙配置
```bash
# UFW 配置
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw allow 8888/tcp    # NAS Dashboard API
sudo ufw enable
```

### 3. 数据库备份
```bash
# 创建备份脚本
sudo tee /usr/local/bin/backup-nas-db.sh > /dev/null << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/var/backups/nas-dashboard"
mkdir -p $BACKUP_DIR

pg_dump -U nas_user nas_dashboard | gzip > $BACKUP_DIR/nas_dashboard_$DATE.sql.gz

# 保留最近7天的备份
find $BACKUP_DIR -name "nas_dashboard_*.sql.gz" -mtime +7 -delete
EOF

sudo chmod +x /usr/local/bin/backup-nas-db.sh

# 添加定时任务
(crontab -l 2>/dev/null; echo "0 2 * * * /usr/local/bin/backup-nas-db.sh") | crontab -
```

## 监控和维护

### 1. 日志查看
```bash
# 查看服务日志
sudo journalctl -u nas-dashboard -f

# 查看应用日志
tail -f /var/log/nas-dashboard/app.log
```

### 2. 性能监控
```bash
# 检查服务状态
sudo systemctl status nas-dashboard

# 查看资源使用
htop
df -h
```

### 3. 更新升级
```bash
# 停止服务
sudo systemctl stop nas-dashboard

# 拉取最新代码
cd /opt/nas-dashboard
git pull origin main

# 更新后端
cd backend
go mod download
go build -o nas-dashboard cmd/server/main.go

# 更新前端
cd ../frontend
npm install
npm run build

# 重启服务
sudo systemctl start nas-dashboard
```

## 故障排除

### 常见问题

#### 1. 数据库连接失败
```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 检查数据库连接
psql -U nas_user -d nas_dashboard -h localhost

# 检查防火墙
sudo ufw status
```

#### 2. MergerFS 挂载失败
```bash
# 检查 MergerFS 安装
which mergerfs

# 检查磁盘挂载
df -h

# 查看系统日志
dmesg | tail
```

#### 3. WebSocket 连接失败
```bash
# 检查代理配置
# Nginx 配置需要支持 WebSocket
location /ws {
    proxy_pass http://localhost:8888;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

## 生产环境建议

### 1. 高可用配置
- 使用主从数据库配置
- 配置负载均衡器
- 实现会话共享

### 2. 性能优化
- 启用 Redis 缓存
- 配置 CDN 加速
- 优化数据库查询

### 3. 安全加固
- 定期更新系统补丁
- 配置入侵检测系统
- 实施访问控制策略

### 4. 备份策略
- 数据库每日备份
- 配置文件定期备份
- 实施灾难恢复计划
