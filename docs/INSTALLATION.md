# NAS Dashboard - Installation Guide

This guide provides comprehensive instructions for installing and deploying the NAS Dashboard in various environments.

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Quick Start](#quick-start)
3. [Docker Installation](#docker-installation)
4. [Manual Installation](#manual-installation)
5. [Environment Configuration](#environment-configuration)
6. [Production Deployment](#production-deployment)
7. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### System Requirements

- **Operating System**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+, Rocky Linux 8+)
- **RAM**: Minimum 2GB (4GB recommended)
- **Disk Space**: 100MB free space for application
- **CPU**: Any modern 64-bit processor
- **Network**: Port 8888 (backend) and 5173 (frontend) available

### Software Requirements

#### For Backend (Go Server)

- **Go**: 1.22 or higher
- **Git**: For cloning the repository
- **Docker**: 20.10+ (if managing containers)
- **System access**: Root or sudo access for system monitoring

#### For Frontend (Vue 3)

- **Node.js**: 18.0 or higher
- **npm**: 9.0 or higher (comes with Node.js)
- **Modern web browser**: Chrome 90+, Firefox 88+, Safari 14+, Edge 90+

### Check Prerequisites

```bash
# Check Go version
go version

# Check Node.js version
node --version

# Check npm version
npm --version

# Check Docker installation
docker --version

# Check system access (should show groups)
groups
```

---

## Quick Start

### Fastest Way to Get Started (Development)

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# 2. Start backend
cd backend
go run cmd/server/main.go &
cd ..

# 3. Start frontend
cd frontend
npm install
npm run dev

# 4. Access the application
# Open browser: http://localhost:5173/login
# Default credentials: admin / admin123
```

### Quick Docker Deployment

```bash
# 1. Clone and navigate
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# 2. Build and start
docker-compose up -d

# 3. Access the application
# Open browser: http://localhost:8888
# Default credentials: admin / admin123
```

---

## Docker Installation

Docker deployment is the recommended method for production use.

### Method 1: Docker Compose (Recommended)

#### 1. Create docker-compose.yml

```yaml
version: '3.8'

services:
  nas-dashboard-backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: nas-dashboard-backend
    ports:
      - "8888:8888"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /:/host:ro
    environment:
      - GIN_MODE=release
      - JWT_SECRET=your-secret-key-change-this
      - CORS_ORIGIN=http://localhost:8888
    restart: unless-stopped
    networks:
      - nas-network

  nas-dashboard-frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: nas-dashboard-frontend
    ports:
      - "80:80"
    depends_on:
      - nas-dashboard-backend
    environment:
      - VITE_API_URL=http://localhost:8888
      - VITE_WS_URL=ws://localhost:8888
    restart: unless-stopped
    networks:
      - nas-network

networks:
  nas-network:
    driver: bridge
```

#### 2. Build and Start

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

#### 3. Access the Application

```
URL: http://your-server-ip
Username: admin
Password: admin123
```

### Method 2: Individual Docker Containers

#### Backend Container

```bash
# Build backend image
cd backend
docker build -t nas-dashboard-backend .

# Run backend container
docker run -d \
  --name nas-dashboard-backend \
  -p 8888:8888 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /:/host:ro \
  -e GIN_MODE=release \
  -e JWT_SECRET=your-secret-key-change-this \
  --restart unless-stopped \
  nas-dashboard-backend
```

#### Frontend Container

```bash
# Build frontend image
cd frontend
docker build -t nas-dashboard-frontend .

# Run frontend container
docker run -d \
  --name nas-dashboard-frontend \
  -p 80:80 \
  -e VITE_API_URL=http://your-server-ip:8888 \
  -e VITE_WS_URL=ws://your-server-ip:8888 \
  --restart unless-stopped \
  nas-dashboard-frontend
```

### Docker Management Commands

```bash
# Stop services
docker-compose stop

# Start services
docker-compose start

# Restart services
docker-compose restart

# Remove services
docker-compose down

# View logs
docker-compose logs backend
docker-compose logs frontend

# Execute commands in container
docker exec -it nas-dashboard-backend /bin/bash
docker exec -it nas-dashboard-frontend /bin/sh

# Update to latest version
docker-compose pull
docker-compose up -d
```

---

## Manual Installation

For development or custom deployment scenarios.

### Backend Installation

#### 1. Install Dependencies

```bash
# Update system
sudo apt update && sudo apt upgrade -y  # Debian/Ubuntu
# sudo yum update -y                     # CentOS/RHEL

# Install Go (if not installed)
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc

# Verify Go installation
go version
```

#### 2. Clone and Build

```bash
# Clone repository
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard/backend

# Download dependencies
go mod download

# Build application
go build -o nas-dashboard cmd/server/main.go

# Test build
./nas-dashboard --help
```

#### 3. Configure Backend

```bash
# Create configuration directory
sudo mkdir -p /etc/nas-dashboard

# Create environment file
sudo nano /etc/nas-dashboard/config.env
```

Add the following content:

```bash
# Server Configuration
GIN_MODE=release
PORT=8888

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h

# CORS Configuration
CORS_ORIGIN=http://localhost:8888

# Monitoring Configuration
MONITOR_INTERVAL=1s
```

#### 4. Install as System Service

```bash
# Create systemd service file
sudo nano /etc/systemd/system/nas-dashboard.service
```

Add the following content:

```ini
[Unit]
Description=NAS Dashboard Backend
After=network.target docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/opt/nas-dashboard
EnvironmentFile=/etc/nas-dashboard/config.env
ExecStart=/opt/nas-dashboard/nas-dashboard
Restart=always
RestartSec=10

# Logging
StandardOutput=journal
StandardError=journal

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

#### 5. Enable and Start Service

```bash
# Copy binary to installation directory
sudo cp nas-dashboard /opt/nas-dashboard/
sudo chmod +x /opt/nas-dashboard/nas-dashboard

# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable nas-dashboard

# Start service
sudo systemctl start nas-dashboard

# Check status
sudo systemctl status nas-dashboard

# View logs
sudo journalctl -u nas-dashboard -f
```

### Frontend Installation

#### 1. Install Node.js Dependencies

```bash
cd frontend

# Install dependencies
npm install

# Test build
npm run build
```

#### 2. Configure Frontend

```bash
# Create environment file
nano .env.production
```

Add the following content:

```bash
# API Configuration
VITE_API_URL=http://your-server-ip:8888
VITE_WS_URL=ws://your-server-ip:8888

# Debug Mode (disable in production)
VITE_DEBUG=false
```

#### 3. Build for Production

```bash
# Build production bundle
npm run build

# Output will be in dist/ directory
ls -la dist/
```

#### 4. Deploy with nginx

```bash
# Install nginx
sudo apt install nginx -y  # Debian/Ubuntu
# sudo yum install nginx -y # CentOS/RHEL

# Create nginx configuration
sudo nano /etc/nginx/sites-available/nas-dashboard
```

Add the following configuration:

```nginx
server {
    listen 80;
    server_name your-server-ip;

    # Frontend static files
    location / {
        root /var/www/nas-dashboard;
        try_files $uri $uri/ /index.html;
        
        # Cache static assets
        location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
            expires 1y;
            add_header Cache-Control "public, immutable";
        }
    }

    # Backend API proxy
    location /api/ {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket proxy
    location /api/monitor/ws {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_read_timeout 86400;
    }

    # Logging
    access_log /var/log/nginx/nas-dashboard-access.log;
    error_log /var/log/nginx/nas-dashboard-error.log;
}
```

#### 5. Deploy Frontend Files

```bash
# Create web directory
sudo mkdir -p /var/www/nas-dashboard

# Copy built files
sudo cp -r dist/* /var/www/nas-dashboard/

# Set permissions
sudo chown -R www-data:www-data /var/www/nas-dashboard
sudo chmod -R 755 /var/www/nas-dashboard

# Enable site
sudo ln -s /etc/nginx/sites-available/nas-dashboard /etc/nginx/sites-enabled/

# Test nginx configuration
sudo nginx -t

# Restart nginx
sudo systemctl restart nginx
```

#### 6. Configure HTTPS (Recommended)

```bash
# Install certbot
sudo apt install certbot python3-certbot-nginx -y

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal is configured automatically
sudo certbot renew --dry-run
```

---

## Environment Configuration

### Backend Environment Variables

```bash
# Server Mode
GIN_MODE=debug              # debug | release
PORT=8888                   # Server port

# JWT Authentication
JWT_SECRET=your-secret-key  # Change this in production!
JWT_ACCESS_DURATION=24h     # Access token lifetime
JWT_REFRESH_DURATION=720h   # Refresh token lifetime (30 days)

# CORS
CORS_ORIGIN=*               # Allowed origins (comma-separated)

# Monitoring
MONITOR_INTERVAL=1s         # WebSocket update interval

# Logging
LOG_LEVEL=info              # debug | info | warn | error
```

### Frontend Environment Variables

```bash
# API Endpoints
VITE_API_URL=http://localhost:8888
VITE_WS_URL=ws://localhost:8888

# Debug Mode
VITE_DEBUG=false            # Enable detailed API logging

# Feature Flags
VITE_ENABLE_DOCKER=true     # Enable Docker management
VITE_ENABLE_SERVICES=true   # Enable service management
```

### Configuration Validation

```bash
# Backend validation
cd backend
go run cmd/server/main.go --validate-config

# Frontend validation
cd frontend
npm run validate-env
```

---

## Production Deployment

### Pre-Deployment Checklist

- [ ] Change JWT_SECRET from default
- [ ] Enable HTTPS with valid SSL certificate
- [ ] Configure firewall rules
- [ ] Set up database backup (if applicable)
- [ ] Configure log rotation
- [ ] Set up monitoring alerts
- [ ] Test backup/restore procedures
- [ ] Review security settings
- [ ] Update CORS origins to specific domains
- [ ] Disable debug mode

### Security Hardening

```bash
# 1. Change JWT secret
openssl rand -base64 32
# Use output as JWT_SECRET

# 2. Configure firewall
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# 3. Set up fail2ban
sudo apt install fail2ban -y
sudo systemctl enable fail2ban
sudo systemctl start fail2ban

# 4. Configure log rotation
sudo nano /etc/logrotate.d/nas-dashboard
```

Add to logrotate config:

```
/var/log/nas-dashboard/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0640 www-data www-data
}
```

### Performance Optimization

```bash
# 1. Enable nginx caching
sudo nano /etc/nginx/nginx.conf

# Add within http block
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=nas_cache:10m max_size=1g inactive=60m;

# 2. Enable gzip compression
gzip on;
gzip_types text/plain text/css application/json application/javascript;
gzip_min_length 1000;

# 3. Configure worker processes
worker_processes auto;
worker_connections 1024;
```

### Monitoring Setup

```bash
# 1. Monitor backend service
sudo systemctl status nas-dashboard

# 2. Monitor nginx
sudo systemctl status nginx

# 3. Check logs
sudo journalctl -u nas-dashboard -f
sudo tail -f /var/log/nginx/nas-dashboard-error.log

# 4. Monitor resources
htop
```

---

## Troubleshooting

### Common Installation Issues

#### Issue 1: Backend Won't Start

**Symptoms:** Service fails to start or exits immediately

**Solutions:**
```bash
# Check service status
sudo systemctl status nas-dashboard

# View detailed logs
sudo journalctl -u nas-dashboard -n 50

# Check if port is already in use
sudo netstat -tulpn | grep 8888

# Verify configuration
cat /etc/nas-dashboard/config.env

# Test manual start
cd /opt/nas-dashboard
sudo ./nas-dashboard
```

#### Issue 2: Frontend Shows Blank Page

**Symptoms:** Browser shows empty page or console errors

**Solutions:**
```bash
# Check nginx configuration
sudo nginx -t

# Verify files exist
ls -la /var/www/nas-dashboard/

# Check nginx logs
sudo tail -f /var/log/nginx/nas-dashboard-error.log

# Clear browser cache and reload
# Ctrl+Shift+R (hard refresh)
```

#### Issue 3: API Connection Refused

**Symptoms:** Frontend can't connect to backend API

**Solutions:**
```bash
# Check if backend is running
sudo systemctl status nas-dashboard

# Verify backend is listening
curl http://localhost:8888/api/health

# Check CORS configuration
curl -H "Origin: http://localhost:8888" http://localhost:8888/api/monitor/cpu

# Verify firewall rules
sudo ufw status
```

#### Issue 4: WebSocket Connection Fails

**Symptoms:** Real-time monitoring doesn't update

**Solutions:**
```bash
# Check WebSocket endpoint
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: test" \
  http://localhost:8888/api/monitor/ws

# Check nginx WebSocket proxy configuration
# Ensure location block has proper proxy settings

# Check browser console for WebSocket errors
# F12 → Console → Filter by "WebSocket"
```

#### Issue 5: Docker Socket Permission Denied

**Symptoms:** Can't manage Docker containers

**Solutions:**
```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and log back in
# Or use newgrp docker
newgrp docker

# Verify access
docker ps
```

#### Issue 6: High Memory Usage

**Symptoms:** Application uses excessive memory

**Solutions:**
```bash
# Check memory usage
free -h
ps aux | grep nas-dashboard

# Adjust monitor interval
# Edit config.env
MONITOR_INTERVAL=5s

# Restart service
sudo systemctl restart nas-dashboard
```

### Debug Mode

Enable debug mode for detailed logging:

```bash
# Backend debug
sudo nano /etc/nas-dashboard/config.env
GIN_MODE=debug
LOG_LEVEL=debug

# Frontend debug
cd frontend
echo "VITE_DEBUG=true" >> .env.production
npm run build

# Restart services
sudo systemctl restart nas-dashboard
sudo systemctl restart nginx
```

### Log Analysis

```bash
# Backend logs
sudo journalctl -u nas-dashboard -f

# Nginx access logs
sudo tail -f /var/log/nginx/nas-dashboard-access.log

# Nginx error logs
sudo tail -f /var/log/nginx/nas-dashboard-error.log

# System logs
sudo dmesg | tail
```

### Performance Issues

```bash
# Check system resources
htop
df -h
iostat -x 1

# Check network
netstat -tulpn
ss -tulpn

# Profile Go application
cd /opt/nas-dashboard
go tool pprof http://localhost:8888/debug/pprof/profile
```

---

## Installation Verification

### Complete Installation Test

```bash
# 1. Check backend health
curl http://localhost:8888/api/health

# 2. Check API endpoints
curl http://localhost:8888/api/monitor/cpu

# 3. Check frontend
curl -I http://localhost/

# 4. Test WebSocket
wscat -c ws://localhost:8888/api/monitor/ws

# 5. Check services
sudo systemctl status nas-dashboard
sudo systemctl status nginx
```

### Browser Test

```
1. Open http://your-server-ip/login
2. Login with admin/admin123
3. Verify dashboard loads
4. Check real-time monitoring updates
5. Test Docker container management
6. Verify all menu items work
```

### Success Criteria

- [ ] Backend service running without errors
- [ ] Frontend loads correctly in browser
- [ ] Login successful with default credentials
- [ ] Real-time monitoring data updates
- [ ] Can view and manage Docker containers
- [ ] All menu items accessible
- [ ] No console errors in browser
- [ ] WebSocket connection stable

---

## Uninstallation

### Remove Docker Deployment

```bash
docker-compose down
docker rmi nas-dashboard-backend nas-dashboard-frontend
docker volume rm nas-data
```

### Remove Manual Installation

```bash
# Stop services
sudo systemctl stop nas-dashboard
sudo systemctl stop nginx

# Disable services
sudo systemctl disable nas-dashboard

# Remove files
sudo rm -rf /opt/nas-dashboard
sudo rm -rf /var/www/nas-dashboard
sudo rm -rf /etc/nas-dashboard
sudo rm /etc/systemd/system/nas-dashboard.service
sudo rm /etc/nginx/sites-available/nas-dashboard
sudo rm /etc/nginx/sites-enabled/nas-dashboard

# Reload systemd
sudo systemctl daemon-reload

# Restart nginx
sudo systemctl restart nginx
```

---

## Getting Help

- **Documentation**: [Project Wiki](https://github.com/yourusername/nas-dashboard/wiki)
- **Issues**: [GitHub Issues](https://github.com/yourusername/nas-dashboard/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/nas-dashboard/discussions)

---

**Last Updated**: 2026-06-12
**Version**: 0.1.0-alpha
