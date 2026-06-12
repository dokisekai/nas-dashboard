# NAS Dashboard Backend 部署指南

## 生产环境部署

### 系统要求
- **操作系统**: Ubuntu 20.04+ 或 Debian 11+
- **CPU**: 2 核心以上
- **内存**: 4GB 以上
- **存储**: 50GB 以上可用空间
- **网络**: 静态 IP 地址

### 依赖安装

#### 1. 安装 PostgreSQL
```bash
# 更新包索引
sudo apt update

# 安装 PostgreSQL
sudo apt install postgresql postgresql-contrib -y

# 启动服务
sudo systemctl start postgresql
sudo systemctl enable postgresql

# 安全配置
sudo -u postgres psql

# 在 PostgreSQL 提示符下运行：
CREATE DATABASE nasdashboard;
CREATE USER nasdashboard WITH ENCRYPTED PASSWORD 'secure_password_here';
GRANT ALL PRIVILEGES ON DATABASE nasdashboard TO nasdashboard;
\q
```

#### 2. 安装 Go 运行时
```bash
# 下载 Go
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz

# 安装 Go
sudo tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz

# 设置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### 3. 安装系统工具
```bash
# 安装必要工具
sudo apt install -y git curl wget vim build-essential

# 安装 Docker（可选）
sudo apt install -y docker.io docker-compose
sudo systemctl start docker
sudo systemctl enable docker

# 添加用户到 docker 组
sudo usermod -aG docker $USER
```

### 应用部署

#### 1. 创建应用用户
```bash
# 创建专用用户
sudo useradd -r -s /bin/bash nas-dashboard
sudo mkdir -p /home/nas-dashboard
sudo chown nas-dashboard:nas-dashboard /home/nas-dashboard
```

#### 2. 部署应用
```bash
# 切换到应用用户
sudo -u nas-dashboard -i

# 克隆或上传应用代码
cd /home/nas-dashboard
git clone <your-repo-url> backend
cd backend

# 安装依赖
go mod download

# 构建应用
go build -o nas-dashboard ./cmd/server/main_new.go

# 设置权限
chmod +x nas-dashboard
```

#### 3. 配置环境变量
```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置文件
vim .env
```

重要配置项：
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=nasdashboard
DB_PASSWORD=your_secure_password
DB_NAME=nasdashboard

# JWT 密钥（使用强随机密钥）
JWT_SECRET=$(openssl rand -base64 32)
JWT_REFRESH_SECRET=$(openssl rand -base64 32)

# 服务器配置
SERVER_HOST=0.0.0.0
SERVER_PORT=8888
GIN_MODE=release

# 日志配置
LOG_LEVEL=info
LOG_FILE=/var/log/nas-dashboard/app.log
```

#### 4. 创建必要的目录
```bash
# 创建日志目录
sudo mkdir -p /var/log/nas-dashboard
sudo chown nas-dashboard:nas-dashboard /var/log/nas-dashboard

# 创建备份目录
sudo mkdir -p /var/backups/nas-dashboard
sudo chown nas-dashboard:nas-dashboard /var/backups/nas-dashboard

# 创建插件目录
sudo mkdir -p /opt/nas-dashboard/plugins
sudo chown nas-dashboard:nas-dashboard /opt/nas-dashboard/plugins
```

### 系统服务配置

#### 1. 创建 Systemd 服务
```bash
sudo vim /etc/systemd/system/nas-dashboard.service
```

添加以下内容：
```ini
[Unit]
Description=NAS Dashboard Backend Service
After=network.target postgresql.service

[Service]
Type=simple
User=nas-dashboard
Group=nas-dashboard
WorkingDirectory=/home/nas-dashboard/backend
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
ExecStart=/home/nas-dashboard/backend/nas-dashboard
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal
SyslogIdentifier=nas-dashboard

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/nas-dashboard /var/backups/nas-dashboard /opt/nas-dashboard/plugins

[Install]
WantedBy=multi-user.target
```

#### 2. 启动服务
```bash
# 重载 systemd 配置
sudo systemctl daemon-reload

# 启动服务
sudo systemctl start nas-dashboard

# 设置开机自启
sudo systemctl enable nas-dashboard

# 检查状态
sudo systemctl status nas-dashboard

# 查看日志
sudo journalctl -u nas-dashboard -f
```

### Nginx 反向代理配置

#### 1. 安装 Nginx
```bash
sudo apt install nginx -y
sudo systemctl start nginx
sudo systemctl enable nginx
```

#### 2. 配置反向代理
```bash
sudo vim /etc/nginx/sites-available/nas-dashboard
```

添加以下配置：
```nginx
server {
    listen 80;
    server_name your-domain.com;

    # 日志配置
    access_log /var/log/nginx/nas-dashboard-access.log;
    error_log /var/log/nginx/nas-dashboard-error.log;

    # 客户端上传大小限制
    client_max_body_size 100M;

    # 代理配置
    location / {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;

        # 超时设置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # WebSocket 支持
    location /ws/ {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "Upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;

        # WebSocket 超时设置
        proxy_connect_timeout 7d;
        proxy_send_timeout 7d;
        proxy_read_timeout 7d;
    }

    # 静态文件缓存
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|svg|woff|woff2|ttf|eot)$ {
        proxy_pass http://localhost:8888;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

#### 3. 启用配置
```bash
# 创建软链接
sudo ln -s /etc/nginx/sites-available/nas-dashboard /etc/nginx/sites-enabled/

# 测试配置
sudo nginx -t

# 重启 Nginx
sudo systemctl restart nginx
```

### SSL/TLS 配置

#### 1. 安装 Certbot
```bash
sudo apt install certbot python3-certbot-nginx -y
```

#### 2. 获取 SSL 证书
```bash
sudo certbot --nginx -d your-domain.com
```

#### 3. 自动续期
```bash
# 测试续期
sudo certbot renew --dry-run

# 添加定时任务
sudo crontab -e
```

添加以下行：
```
0 0 * * 0 certbot renew --quiet --post-hook "systemctl reload nginx"
```

### 安全加固

#### 1. 防火墙配置
```bash
# 安装 UFW
sudo apt install ufw -y

# 默认策略
sudo ufw default deny incoming
sudo ufw default allow outgoing

# 允许 SSH
sudo ufw allow ssh

# 允许 HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 启用防火墙
sudo ufw enable

# 检查状态
sudo ufw status
```

#### 2. 应用安全
```bash
# 限制文件权限
sudo chmod 600 /home/nas-dashboard/backend/.env
sudo chmod 750 /home/nas-dashboard/backend
sudo chmod 640 /var/log/nas-dashboard/*
sudo chmod 750 /var/backups/nas-dashboard

# 设置 SELinux（如果使用）
sudo setsebool -P httpd_can_network_connect_db 1
sudo setsebool -P httpd_can_network_connect 1
```

#### 3. 数据库安全
```bash
# 编辑 PostgreSQL 配置
sudo vim /etc/postgresql/12/main/pg_hba.conf

# 修改认证方法为 md5 或 scram-sha-256
local   all             all                                     md5
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5

# 重启 PostgreSQL
sudo systemctl restart postgresql
```

### 监控和日志

#### 1. 日志轮转配置
```bash
sudo vim /etc/logrotate.d/nas-dashboard
```

添加以下内容：
```
/var/log/nas-dashboard/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0640 nas-dashboard nas-dashboard
    sharedscripts
    postrotate
        systemctl reload nas-dashboard > /dev/null 2>&1 || true
    endscript
}
```

#### 2. 监控脚本
```bash
sudo vim /usr/local/bin/nas-dashboard-monitor.sh
```

添加以下内容：
```bash
#!/bin/bash

# 检查服务状态
if ! systemctl is-active --quiet nas-dashboard; then
    echo "NAS Dashboard service is not running" | mail -s "Service Alert" admin@example.com
    systemctl start nas-dashboard
fi

# 检查磁盘空间
DISK_USAGE=$(df -h / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $DISK_USAGE -gt 80 ]; then
    echo "Disk usage is ${DISK_USAGE}%" | mail -s "Disk Space Alert" admin@example.com
fi

# 检查内存使用
MEM_USAGE=$(free | awk 'NR==2{printf "%.0f", $3/$2*100}')
if [ $MEM_USAGE -gt 90 ]; then
    echo "Memory usage is ${MEM_USAGE}%" | mail -s "Memory Alert" admin@example.com
fi
```

设置可执行权限和定时任务：
```bash
sudo chmod +x /usr/local/bin/nas-dashboard-monitor.sh
sudo crontab -e
```

添加：
```
*/5 * * * * /usr/local/bin/nas-dashboard-monitor.sh
```

### 备份策略

#### 1. 自动备份脚本
```bash
sudo vim /usr/local/bin/nas-dashboard-backup.sh
```

添加以下内容：
```bash
#!/bin/bash

# 配置
BACKUP_DIR="/var/backups/nas-dashboard"
RETENTION_DAYS=30
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="auto_backup_${DATE}"

# 创建备份目录
mkdir -p "${BACKUP_DIR}/${DATE}"

# 数据库备份
pg_dump -h localhost -U nasdashboard nasdashboard | gzip > "${BACKUP_DIR}/${DATE}/database.sql.gz"

# 文件备份
tar -czf "${BACKUP_DIR}/${DATE}/files.tar.gz" \
    /home/nas-dashboard/backend \
    /etc/nginx/sites-available/nas-dashboard \
    /etc/systemd/system/nas-dashboard.service

# 清理旧备份
find ${BACKUP_DIR} -type d -mtime +${RETENTION_DAYS} -exec rm -rf {} \;

echo "Backup completed: ${BACKUP_NAME}"
```

#### 2. 定时备份
```bash
sudo chmod +x /usr/local/bin/nas-dashboard-backup.sh
sudo crontab -e
```

添加：
```
0 2 * * * /usr/local/bin/nas-dashboard-backup.sh
```

### 性能优化

#### 1. 数据库优化
```sql
-- 连接到 PostgreSQL
sudo -u postgres psql nasdashboard

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_operation_logs_user_id ON operation_logs(user_id);
CREATE INDEX idx_file_system_access_user_id ON file_system_accesses(user_id);

-- 配置 PostgreSQL
sudo vim /etc/postgresql/12/main/postgresql.conf
```

添加以下配置：
```ini
# 连接设置
max_connections = 100
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
random_page_cost = 1.1
effective_io_concurrency = 200
work_mem = 2621kB
min_wal_size = 1GB
max_wal_size = 4GB
```

#### 2. 应用优化
```bash
# 编辑环境变量
vim /home/nas-dashboard/backend/.env
```

添加性能优化配置：
```bash
# 数据库连接池
DB_MAX_OPEN_CONNS=100
DB_MAX_IDLE_CONNS=10
DB_CONN_MAX_LIFETIME=3600

# WebSocket 配置
WS_HEARTBEAT_INTERVAL=30
WS_MESSAGE_BUFFER_SIZE=256
WS_MAX_CONNECTIONS=1000

# 文件操作配置
FILE_MAX_UPLOAD_SIZE=104857600
FILE_OPERATION_TIMEOUT=300
```

### 故障排除

#### 1. 服务启动失败
```bash
# 检查日志
sudo journalctl -u nas-dashboard -n 50

# 检查端口占用
sudo netstat -tunlp | grep 8888

# 检查文件权限
sudo -u nas-dashboard /home/nas-dashboard/backend/nas-dashboard
```

#### 2. 数据库连接失败
```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 测试连接
psql -h localhost -U nasdashboard -d nasdashboard

# 检查防火墙
sudo ufw status
```

#### 3. 性能问题
```bash
# 检查系统资源
htop

# 检查数据库性能
sudo -u postgres psql -c "SELECT * FROM pg_stat_activity;"

# 检查应用日志
sudo tail -f /var/log/nas-dashboard/app.log
```

### 更新和维护

#### 1. 应用更新
```bash
# 停止服务
sudo systemctl stop nas-dashboard

# 备份当前版本
sudo cp /home/nas-dashboard/backend/nas-dashboard /home/nas-dashboard/backend/nas-dashboard.backup

# 更新代码
cd /home/nas-dashboard/backend
git pull origin main

# 重新构建
go build -o nas-dashboard ./cmd/server/main_new.go

# 运行迁移
./nas-dashboard --migrate

# 启动服务
sudo systemctl start nas-dashboard
```

#### 2. 数据库迁移
```bash
# 备份数据库
pg_dump -h localhost -U nasdashboard nasdashboard > backup_$(date +%Y%m%d).sql

# 运行迁移
sudo systemctl restart nas-dashboard
```

### 健康检查

创建健康检查端点监控：
```bash
# 添加到 crontab
*/5 * * * * curl -f http://localhost:8888/health || echo "Health check failed" | mail -s "Health Alert" admin@example.com
```

这个部署指南提供了完整的生产环境部署流程，确保 NAS Dashboard 后端稳定、安全地运行。
