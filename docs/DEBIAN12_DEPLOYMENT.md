# Debian 12 NAS Dashboard 部署指南

## 系统要求

### 硬件兼容性
- **处理器**: x86_64 (amd64) 或 ARM64 (aarch64)
- **内存**: 最少 2GB，推荐 4GB+
- **存储**: 最少 20GB 可用空间
- **网络**: 千兆以太网接口

### 软件要求
- **操作系统**: Debian 12 (Bookworm)
- **架构**: 64位
- **权限**: Root 或 sudo 访问权限

## 快速安装指南

### 1. 准备工作

#### 下载 Debian 12
```bash
wget https://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-12.0.0-amd64-netinst.iso
```

#### 推荐的分区方案
```
/       - 20GB   (根分区)
/var    - 10GB   (变量数据)
/home   - 剩余空间 (用户数据)
swap    - 4GB    (交换分区)
```

#### 最小化安装建议
在安装过程中，只选择以下软件包：
- 标准系统实用程序
- SSH 服务器

避免安装桌面环境和不必要的服务。

### 2. 一键安装脚本

```bash
# 下载安装脚本
wget https://raw.githubusercontent.com/yourname/nas-dashboard/main/scripts/install-debian12.sh

# 或者克隆整个仓库
git clone https://github.com/yourname/nas-dashboard.git
cd nas-dashboard

# 运行安装脚本
sudo bash scripts/install-debian12.sh
```

### 3. 手动安装步骤

如果自动脚本无法使用，可以按照以下步骤手动安装：

#### 3.1 安装系统依赖
```bash
sudo apt update
sudo apt install -y \
    postgresql postgresql-contrib \
    nginx certbot python3-certbot-nginx \
    ufw fail2ban \
    mergerfs lvm2 mdadm smartmontools \
    nut nut-client \
    curl wget git build-essential
```

#### 3.2 安装 Go 环境
```bash
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 3.3 安装 Node.js
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
```

#### 3.4 安装 Docker
```bash
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian bookworm stable" | sudo tee /etc/apt/sources.list.d/docker.list

sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io
sudo systemctl enable docker
sudo systemctl start docker
```

## 网络配置

### Debian 12 网络配置方法

Debian 12 支持两种主要的网络配置方法：

#### 1. 传统 /etc/network/interfaces 方法

```bash
sudo nano /etc/network/interfaces
```

静态IP配置示例：
```bash
# The loopback network interface
auto lo
iface lo inet loopback

# The primary network interface
auto eth0
iface eth0 inet static
    address 192.168.1.100
    netmask 255.255.255.0
    gateway 192.168.1.1
    dns-nameservers 8.8.8.8 8.8.4.4
```

DHCP配置示例：
```bash
# The primary network interface
auto eth0
iface eth0 inet dhcp
```

#### 2. NetworkManager 方法（推荐桌面环境）

```bash
# 安装 NetworkManager
sudo apt install network-manager

# 查看网络设备
nmcli device

# 配置静态IP
sudo nmcli connection modify "Wired connection 1" \
    ipv4.addresses 192.168.1.100/24 \
    ipv4.gateway 192.168.1.1 \
    ipv4.dns "8.8.8.8 8.8.4.4" \
    ipv4.method manual

# 重启网络连接
sudo nmcli connection down "Wired connection 1"
sudo nmcli connection up "Wired connection 1"
```

### DNS 配置

Debian 12 推荐使用 systemd-resolved：

```bash
# 启用 systemd-resolved
sudo systemctl enable systemd-resolved
sudo systemctl start systemd-resolved

# 设置 DNS 服务器
sudo resolvectl dns eth0 8.8.8.8 8.8.4.4

# 查看当前 DNS 配置
resolvectl status
```

## 存储配置

### 1. MergerFS 存储聚合

#### 安装 MergerFS
```bash
sudo apt install mergerfs
```

#### 基本使用
```bash
# 创建挂载点
sudo mkdir -p /mnt/disk1 /mnt/disk2 /mnt/storage

# 挂载多个磁盘到单个目录
sudo mergerfs /mnt/disk1:/mnt/disk2 /mnt/storage
```

#### 持久化配置
编辑 `/etc/fstab`：
```bash
# MergerFS 存储聚合
/mnt/disk1 /mnt/disk2 /mnt/storage fuse.mergerfs defaults,allow_other,category.create=mfs,minfreespace=10G 0 0
```

### 2. LVM 逻辑卷管理

#### 创建 LVM 卷
```bash
# 创建物理卷
sudo pvcreate /dev/sdb /dev/sdc

# 创建卷组
sudo vgcreate storage_pool /dev/sdb /dev/sdc

# 创建逻辑卷
sudo lvcreate -L 1T -n data_volume storage_pool

# 创建文件系统
sudo mkfs.ext4 /dev/storage_pool/data_volume

# 挂载
sudo mkdir -p /mnt/data
sudo mount /dev/storage_pool/data_volume /mnt/data
```

### 3. Btrfs 文件系统

#### 创建 Btrfs 文件系统
```bash
# 创建 Btrfs 文件系统
sudo mkfs.btrfs -d single -m single /dev/sdb /dev/sdc

# 挂载
sudo mkdir -p /mnt/btrfs
sudo mount /dev/sdb /mnt/btrfs
```

#### 创建快照
```bash
# 创建子卷
sudo btrfs subvolume create /mnt/btrfs/data

# 创建快照
sudo btrfs subvolume snapshot -r /mnt/btrfs/data /mnt/btrfs/data_snapshot
```

## 防火墙和安全配置

### 1. UFW 防火墙配置

Debian 12 需要先安装 UFW：

```bash
# 安装 UFW
sudo apt install ufw

# 基本配置
sudo ufw default deny incoming
sudo ufw default allow outgoing

# 允许必要端口
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 8080/tcp  # NAS Dashboard

# 启用防火墙
sudo ufw enable

# 查看状态
sudo ufw status verbose
```

### 2. Fail2Ban 安全

```bash
# 安装 Fail2Ban
sudo apt install fail2ban

# 创建本地配置
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local

# 编辑配置
sudo nano /etc/fail2ban/jail.local
```

配置示例：
```ini
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 5

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
maxretry = 3

[nginx-http-auth]
enabled = true
port = http,https
logpath = /var/log/nginx/error.log
```

```bash
# 重启服务
sudo systemctl restart fail2ban
```

## 性能优化

### 1. 内核参数优化

创建 `/etc/sysctl.d/99-nas-dashboard.conf`：

```bash
# 网络性能优化
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 87380 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_window_scaling = 1

# 文件系统优化
fs.file-max = 2097152
fs.inotify.max_user_watches = 524288

# 虚拟内存优化
vm.swappiness = 10
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5

# 安全设置
net.ipv4.conf.default.rp_filter = 1
net.ipv4.conf.all.rp_filter = 1
```

应用参数：
```bash
sudo sysctl -p /etc/sysctl.d/99-nas-dashboard.conf
```

### 2. 系统限制优化

编辑 `/etc/security/limits.conf`：

```bash
* soft nofile 65536
* hard nofile 65536
* soft nproc 65536
* hard nproc 65536
```

## 数据库配置

### PostgreSQL 优化

编辑 `/etc/postgresql/15/main/postgresql.conf`：

```bash
# 内存配置
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
work_mem = 16MB

# 连接配置
max_connections = 100

# WAL 配置
wal_buffers = 16MB
checkpoint_completion_target = 0.9
```

重启 PostgreSQL：
```bash
sudo systemctl restart postgresql
```

## SSL 证书配置

### 使用 Let's Encrypt

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取证书
sudo certbot --nginx -d yourdomain.com

# 自动续期
sudo certbot renew --dry-run
```

## 服务管理

### 启动和停止服务

```bash
# 启动 NAS Dashboard 后端
sudo systemctl start nas-dashboard-backend

# 启动 NAS Dashboard 前端
sudo systemctl start nas-dashboard-frontend

# 启用开机自启
sudo systemctl enable nas-dashboard-backend
sudo systemctl enable nas-dashboard-frontend

# 查看服务状态
sudo systemctl status nas-dashboard-backend
sudo systemctl status nas-dashboard-frontend

# 查看日志
sudo journalctl -u nas-dashboard-backend -f
sudo journalctl -u nas-dashboard-frontend -f
```

## 故障排除

### 常见问题

#### 1. 端口冲突
```bash
# 检查端口占用
sudo netstat -tulpn | grep :8080
sudo ss -tulpn | grep :8080

# 杀死占用进程
sudo kill -9 <PID>
```

#### 2. 数据库连接问题
```bash
# 检查 PostgreSQL 状态
sudo systemctl status postgresql

# 测试数据库连接
sudo -u postgres psql -c "SELECT version();"
```

#### 3. 权限问题
```bash
# 检查文件权限
ls -la /opt/nas-dashboard

# 修复权限
sudo chown -R $USER:$USER /opt/nas-dashboard
sudo chmod -R 755 /opt/nas-dashboard
```

#### 4. 内存不足
```bash
# 检查内存使用
free -h

# 查看进程内存占用
ps aux --sort=-%mem | head -10

# 创建交换空间（如果需要）
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

### 日志查看

```bash
# 系统日志
sudo journalctl -xe

# NAS Dashboard 日志
sudo tail -f /var/log/nas-dashboard/backend.log
sudo tail -f /var/log/nas-dashboard/frontend.log

# Nginx 日志
sudo tail -f /var/log/nginx/nas-dashboard-access.log
sudo tail -f /var/log/nginx/nas-dashboard-error.log

# PostgreSQL 日志
sudo tail -f /var/log/postgresql/postgresql-15-main.log
```

## 备份和恢复

### 数据库备份

```bash
# 创建备份目录
sudo mkdir -p /backup/postgresql

# 备份数据库
sudo -u postgres pg_dump nas_dashboard > /backup/postgresql/nas_dashboard_$(date +%Y%m%d).sql

# 恢复数据库
sudo -u postgres psql nas_dashboard < /backup/postgresql/nas_dashboard_20240101.sql
```

### 自动备份脚本

创建 `/usr/local/bin/backup-nas.sh`：

```bash
#!/bin/bash
BACKUP_DIR="/backup"
DATE=$(date +%Y%m%d_%H%M%S)

# 数据库备份
sudo -u postgres pg_dump nas_dashboard > "$BACKUP_DIR/postgresql/nas_dashboard_$DATE.sql"

# 配置文件备份
sudo tar -czf "$BACKUP_DIR/config/nas_config_$DATE.tar.gz" /etc/nas-dashboard

# 清理旧备份（保留30天）
find "$BACKUP_DIR" -mtime +30 -delete

echo "Backup completed: $DATE"
```

设置定时任务：
```bash
sudo crontab -e
# 每天凌晨2点执行备份
0 2 * * * /usr/local/bin/backup-nas.sh
```

## 性能监控

### 系统监控

```bash
# 安装监控工具
sudo apt install htop iotop nethogs

# CPU 和内存监控
htop

# 磁盘 I/O 监控
sudo iotop

# 网络监控
sudo nethogs
```

### NAS Dashboard 监控

通过 Web 界面访问监控系统：
- CPU 使用率
- 内存使用情况
- 磁盘空间和 I/O
- 网络流量
- 系统温度

## 安全建议

1. **定期更新系统**
   ```bash
   sudo apt update
   sudo apt upgrade
   ```

2. **启用自动安全更新**
   ```bash
   sudo apt install unattended-upgrades
   sudo dpkg-reconfigure -plow unattended-upgrades
   ```

3. **禁用 root SSH 登录**
   ```bash
   sudo nano /etc/ssh/sshd_config
   # 修改：PermitRootLogin no
   sudo systemctl restart sshd
   ```

4. **使用密钥认证**
   ```bash
   ssh-keygen -t rsa -b 4096
   ssh-copy-id user@server
   ```

5. **定期备份**
   - 设置自动备份任务
   - 验证备份完整性
   - 测试恢复流程

## 获取帮助

- **官方文档**: https://github.com/yourname/nas-dashboard/wiki
- **问题反馈**: https://github.com/yourname/nas-dashboard/issues
- **社区论坛**: https://community.nas-dashboard.com

## 更新日志

### v1.0.0 (2024-01-01)
- 初始版本发布
- 支持 Debian 12
- 完整的存储管理功能
- Docker 容器管理
- 系统监控和告警