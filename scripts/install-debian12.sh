#!/bin/bash

################################################################################
# Debian 12 NAS Dashboard 一键安装脚本
# 适用于 Debian 12 (Bookworm) 系统
################################################################################

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查是否为 root 用户
check_root() {
    if [ "$EUID" -ne 0 ]; then
        log_error "请使用 root 权限运行此脚本"
        exit 1
    fi
}

# 检查系统版本
check_system() {
    log_step "检查系统版本..."

    if [ ! -f /etc/os-release ]; then
        log_error "无法检测操作系统版本"
        exit 1
    fi

    . /etc/os-release

    if [ "$ID" != "debian" ]; then
        log_error "此脚本仅支持 Debian 系统"
        exit 1
    fi

    if [ "$VERSION_ID" != "12" ]; then
        log_warn "检测到 Debian 版本: $VERSION_ID"
        log_warn "此脚本专为 Debian 12 设计，其他版本可能无法正常工作"
        read -p "是否继续？(y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi

    log_info "✓ 系统检查通过: Debian $VERSION_ID"
}

# 检查系统资源
check_resources() {
    log_step "检查系统资源..."

    # 检查内存
    total_mem=$(free -m | awk '/^Mem:/{print $2}')
    log_info "总内存: ${total_mem}MB"
    if [ $total_mem -lt 2048 ]; then
        log_warn "内存不足 2GB，建议至少 4GB"
    fi

    # 检查磁盘空间
    available_space=$(df -BG / | awk 'NR==2{print $4}' | sed 's/G//')
    log_info "可用磁盘空间: ${available_space}GB"
    if [ $available_space -lt 20 ]; then
        log_error "磁盘空间不足 20GB"
        exit 1
    fi

    # 检查 CPU
    cpu_cores=$(nproc)
    log_info "CPU 核心: ${cpu_cores}"

    log_info "✓ 资源检查完成"
}

# 配置 Debian 12 软件源
configure_sources() {
    log_step "配置软件源..."

    # 备份原有源配置
    if [ ! -f /etc/apt/sources.list.bak ]; then
        cp /etc/apt/sources.list /etc/apt/sources.list.bak
    fi

    # 检查是否需要配置镜像源（可根据需要修改）
    log_info "使用默认 Debian 官方源"
    log_info "如需使用国内镜像源，请手动修改 /etc/apt/sources.list"

    # 更新软件包列表
    log_info "更新软件包列表..."
    apt-get update -q
}

# 安装系统依赖
install_dependencies() {
    log_step "安装系统依赖..."

    # 基础工具
    apt-get install -y \
        curl \
        wget \
        git \
        build-essential \
        pkg-config \
        libssl-dev \
        ca-certificates \
        gnupg \
        lsb-release

    # 数据库依赖
    apt-get install -y \
        postgresql \
        postgresql-contrib \
        libpq-dev

    # Web 服务器
    apt-get install -y \
        nginx \
        certbot \
        python3-certbot-nginx

    # 防火墙和安全
    apt-get install -y \
        ufw \
        fail2ban

    # 存储相关
    apt-get install -y \
        mergerfs \
        lvm2 \
        mdadm \
        smartmontools \
        xfsprogs \
        btrfs-progs

    # UPS 支持
    apt-get install -y \
        nut \
        nut-client

    # 其他工具
    apt-get install -y \
        htop \
        iotop \
        nethogs \
        tmux

    log_info "✓ 系统依赖安装完成"
}

# 安装 Go (后端需要)
install_go() {
    log_step "安装 Go 语言环境..."

    if command -v go &> /dev/null; then
        go_version=$(go version | awk '{print $3}')
        log_info "Go 已安装: $go_version"
        return
    fi

    GO_VERSION="1.21.6"
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)  GO_ARCH="amd64" ;;
        aarch64) GO_ARCH="arm64" ;;
        armv7l)  GO_ARCH="arm" ;;
        *)
            log_error "不支持的架构: $ARCH"
            exit 1
            ;;
    esac

    GO_FILE="go${GO_VERSION}.linux-${GO_ARCH}.tar.gz"
    GO_URL="https://go.dev/dl/${GO_FILE}"

    log_info "下载 Go ${GO_VERSION}..."
    cd /tmp
    wget -q $GO_URL

    log_info "安装 Go..."
    tar -C /usr/local -xzf $GO_FILE

    # 设置环境变量
    cat >> /etc/profile <<EOF
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/root/go
EOF

    export PATH=$PATH:/usr/local/go/bin
    rm -f $GO_FILE

    log_info "✓ Go 安装完成"
}

# 安装 Node.js (前端需要)
install_nodejs() {
    log_step "安装 Node.js 环境..."

    if command -v node &> /dev/null; then
        node_version=$(node --version)
        log_info "Node.js 已安装: $node_version"
        return
    fi

    # 安装 Node.js 20.x LTS
    curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
    apt-get install -y nodejs

    log_info "✓ Node.js 安装完成"
}

# 安装 Docker
install_docker() {
    log_step "安装 Docker..."

    if command -v docker &> /dev/null; then
        docker_version=$(docker --version)
        log_info "Docker 已安装: $docker_version"
        return
    fi

    # 添加 Docker 官方 GPG 密钥
    install -m 0755 -d /etc/apt/keyrings
    curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
    chmod a+r /etc/apt/keyrings/docker.gpg

    # 设置 Docker 仓库
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      tee /etc/apt/sources.list.d/docker.list > /dev/null

    # 安装 Docker
    apt-get update -q
    apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # 启动 Docker 服务
    systemctl enable docker
    systemctl start docker

    # 添加当前用户到 docker 组（如果存在）
    if [ -n "$SUDO_USER" ]; then
        usermod -aG docker $SUDO_USER
        log_info "已将用户 $SUDO_USER 添加到 docker 组"
    fi

    log_info "✓ Docker 安装完成"
}

# 配置防火墙
configure_firewall() {
    log_step "配置防火墙..."

    # 默认拒绝入站，允许出站
    ufw --force reset
    ufw default deny incoming
    ufw default allow outgoing

    # 允许 SSH (避免锁定自己)
    ufw allow 22/tcp comment 'SSH'

    # 允许 HTTP 和 HTTPS
    ufw allow 80/tcp comment 'HTTP'
    ufw allow 443/tcp comment 'HTTPS'

    # 允许 NAS 面板端口 (可根据需要修改)
    UFW_PANEL_PORT=${UFW_PANEL_PORT:-8080}
    ufw allow $UFW_PANEL_PORT/tcp comment "NAS Panel"

    # 启用防火墙
    ufw --force enable

    log_info "✓ 防火墙配置完成"
}

# 配置 PostgreSQL
configure_database() {
    log_step "配置 PostgreSQL..."

    # 确保 PostgreSQL 服务运行
    systemctl enable postgresql
    systemctl start postgresql

    # 创建数据库和用户
    DB_NAME=${DB_NAME:-"nas_dashboard"}
    DB_USER=${DB_USER:-"nas_dashboard"}
    DB_PASSWORD=${DB_PASSWORD:-$(openssl rand -base64 16)}

    log_info "创建数据库: $DB_NAME"
    log_info "创建用户: $DB_USER"

    # 使用 postgres 用户执行 SQL
    sudo -u postgres psql <<EOF
-- 创建数据库
CREATE DATABASE $DB_NAME;

-- 创建用户和密码
CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';

-- 授予权限
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

-- 连接到数据库并授予 schema 权限
\c $DB_NAME
GRANT ALL ON SCHEMA public TO $DB_USER;
EOF

    # 保存数据库凭据
    cat > /etc/nas-dashboard/db.conf <<EOF
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_HOST=localhost
DB_PORT=5432
EOF

    chmod 600 /etc/nas-dashboard/db.conf

    log_info "✓ PostgreSQL 配置完成"
    log_warn "数据库凭据已保存到: /etc/nas-dashboard/db.conf"
}

# 配置 Nginx
configure_nginx() {
    log_step "配置 Nginx..."

    # 创建 NAS 面板配置
    cat > /etc/nginx/sites-available/nas-dashboard <<'EOF'
server {
    listen 80;
    server_name _;

    # 日志文件
    access_log /var/log/nginx/nas-dashboard-access.log;
    error_log /var/log/nginx/nas-dashboard-error.log;

    # 前端静态文件
    location / {
        proxy_pass http://127.0.0.1:5173;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # 后端 API
    location /api/ {
        proxy_pass http://127.0.0.1:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

    # 启用配置
    ln -sf /etc/nginx/sites-available/nas-dashboard /etc/nginx/sites-enabled/

    # 移除默认配置
    rm -f /etc/nginx/sites-enabled/default

    # 测试配置
    nginx -t

    # 重启 Nginx
    systemctl restart nginx

    log_info "✓ Nginx 配置完成"
}

# 创建 systemd 服务
create_systemd_service() {
    log_step "创建 systemd 服务..."

    INSTALL_DIR=${INSTALL_DIR:-"/opt/nas-dashboard"}

    # 创建后端服务
    cat > /etc/systemd/system/nas-dashboard-backend.service <<EOF
[Unit]
Description=NAS Dashboard Backend
After=network.target postgresql.service docker.service

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR/backend
Environment="GIN_MODE=release"
EnvironmentFile=/etc/nas-dashboard/db.conf
ExecStart=$INSTALL_DIR/backend/nas-dashboard
Restart=always
RestartSec=5

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/nas-dashboard /tmp

[Install]
WantedBy=multi-user.target
EOF

    # 创建前端服务（开发模式）
    cat > /etc/systemd/system/nas-dashboard-frontend.service <<EOF
[Unit]
Description=NAS Dashboard Frontend
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$INSTALL_DIR/frontend
Environment="NODE_ENV=production"
ExecStart=/usr/bin/npm run dev
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

    # 重新加载 systemd
    systemctl daemon-reload

    log_info "✓ systemd 服务创建完成"
}

# 优化内核参数
optimize_kernel() {
    log_step "优化内核参数..."

    cat >> /etc/sysctl.conf <<EOF

# NAS Dashboard 优化参数
# 网络性能优化
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_window_scaling = 1

# 文件描述符限制
fs.file-max = 2097152

# 共享内存
kernel.shmmax = 68719476736
kernel.shmall = 4294967296

# 虚拟内存
vm.swappiness = 10
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5

# 安全设置
net.ipv4.conf.default.rp_filter = 1
net.ipv4.conf.all.rp_filter = 1
net.ipv4.icmp_echo_ignore_broadcasts = 1
net.ipv4.conf.all.accept_source_route = 0
net.ipv6.conf.all.accept_source_route = 0
EOF

    # 应用内核参数
    sysctl -p

    log_info "✓ 内核参数优化完成"
}

# 创建目录结构
create_directories() {
    log_step "创建目录结构..."

    mkdir -p /etc/nas-dashboard
    mkdir -p /var/log/nas-dashboard
    mkdir -p /var/lib/nas-dashboard
    mkdir -p /opt/nas-dashboard

    log_info "✓ 目录结构创建完成"
}

# 显示安装摘要
show_summary() {
    log_step "安装完成！"

    echo ""
    echo "============================================"
    echo "  NAS Dashboard 安装摘要"
    echo "============================================"
    echo ""
    echo "📦 系统服务状态："
    echo "  - PostgreSQL: $(systemctl is-active postgresql)"
    echo "  - Nginx: $(systemctl is-active nginx)"
    echo "  - Docker: $(systemctl is-active docker)"
    echo ""
    echo "🔐 防火墙状态: $(ufw status | head -1)"
    echo ""
    echo "📁 安装目录: /opt/nas-dashboard"
    echo "📝 日志目录: /var/log/nas-dashboard"
    echo "🔧 配置目录: /etc/nas-dashboard"
    echo ""
    echo "🔑 数据库凭据:"
    echo "  - 数据库: $DB_NAME"
    echo "  - 用户: $DB_USER"
    echo "  - 密码: $DB_PASSWORD"
    echo "  - 凭据文件: /etc/nas-dashboard/db.conf"
    echo ""
    echo "🚀 下一步操作："
    echo "  1. 将 NAS Dashboard 代码复制到 /opt/nas-dashboard"
    echo "  2. 构建后端: cd /opt/nas-dashboard/backend && go build -o nas-dashboard ./cmd/server"
    echo "  3. 构建前端: cd /opt/nas-dashboard/frontend && npm install && npm run build"
    echo "  4. 启动服务: systemctl start nas-dashboard-backend"
    echo "  5. 配置 SSL: certbot --nginx -d yourdomain.com"
    echo ""
    echo "============================================"
    echo ""
}

# 主函数
main() {
    echo ""
    echo "============================================"
    echo "  Debian 12 NAS Dashboard 安装向导"
    echo "============================================"
    echo ""

    check_root
    check_system
    check_resources
    configure_sources
    install_dependencies
    install_go
    install_nodejs
    install_docker
    create_directories
    configure_firewall
    configure_database
    configure_nginx
    create_systemd_service
    optimize_kernel
    show_summary

    log_info "安装脚本执行完成！"
}

# 运行主函数
main "$@"