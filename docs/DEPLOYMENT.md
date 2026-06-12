# NAS Dashboard - Deployment Guide

Complete deployment guide for the NAS Dashboard, covering development, staging, and production environments.

## 📋 Table of Contents

1. [Deployment Overview](#deployment-overview)
2. [Prerequisites](#prerequisites)
3. [Development Deployment](#development-deployment)
4. [Docker Deployment](#docker-deployment)
5. [Production Deployment](#production-deployment)
6. [Security Configuration](#security-configuration)
7. [Performance Optimization](#performance-optimization)
8. [Monitoring and Logging](#monitoring-and-logging)
9. [Backup and Recovery](#backup-and-recovery)
10. [Troubleshooting](#troubleshooting)

---

## Deployment Overview

### Deployment Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Production Server                       │
│  ┌─────────────────────────────────────────────────────────┐│
│  │              Nginx Reverse Proxy                        ││
│  │              (SSL/TLS, Static Files)                     ││
│  └────────────┬──────────────────────┬────────────────────┘│
│               │                      │                      │
│  ┌────────────┴─────┐    ┌───────────┴───────────────────┐ │
│  │  Frontend Server  │    │   Backend Server             │ │
│  │  (Static Files)   │    │   (Go Application)           │ │
│  │  Port: 3000       │    │   Port: 8888                 │ │
│  └──────────────────┘    └──────────┬───────────────────┘ │
│                                      │                      │
│  ┌───────────────────────────────────┼───────────────────┐│
│  │         PostgreSQL Database        │                   ││
│  │         Port: 5432                 │                   ││
│  └───────────────────────────────────┴───────────────────┘│
│  ┌─────────────────────────────────────────────────────────┐│
│  │            Docker Daemon (Optional)                      ││
│  └─────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

### Deployment Options

1. **Docker Compose** (Recommended for most users)
   - Easy setup and management
   - Isolated environments
   - Quick deployment

2. **Manual Deployment** (Advanced users)
   - Full control over configuration
   - Optimized for specific hardware
   - Manual updates required

3. **Kubernetes** (Enterprise deployments)
   - Scalable and resilient
   - Complex setup
   - Requires Kubernetes cluster

---

## Prerequisites

### System Requirements

#### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4 GB
- **Storage**: 20 GB
- **OS**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+)

#### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 8+ GB
- **Storage**: 50+ GB SSD
- **OS**: Ubuntu 22.04 LTS

### Software Requirements

#### Required
- **Docker**: 20.10+
- **Docker Compose**: 2.0+
- **Git**: For cloning repository
- **Text Editor**: For configuration

#### Optional (Manual Deployment)
- **Go**: 1.22+ (for building backend)
- **Node.js**: 20+ (for building frontend)
- **Nginx**: 1.18+ (for reverse proxy)
- **PostgreSQL**: 14+ (for database)

### Network Requirements

- **Open Ports**: 80, 443 (HTTP/HTTPS)
- **Backend Port**: 8888 (internal only)
- **Database Port**: 5432 (internal only)
- **WebSocket Support**: For real-time monitoring

---

## Development Deployment

### Local Development Setup

#### 1. Clone Repository

```bash
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard
```

#### 2. Backend Development

```bash
cd backend

# Install Go dependencies
go mod download

# Create development configuration
cat > config/dev.env << EOF
PORT=8888
GIN_MODE=debug
JWT_SECRET=dev-secret-key-change-in-production
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h
CORS_ORIGIN=http://localhost:5173
DB_HOST=localhost
DB_PORT=5432
DB_NAME=nas_dashboard_dev
DB_USER=nas_user
DB_PASSWORD=dev_password
LOG_LEVEL=debug
EOF

# Run backend
go run cmd/server/main.go
```

Backend will run on `http://localhost:8888`

#### 3. Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Create development configuration
cat > .env << EOF
VITE_API_URL=http://localhost:8888
VITE_WS_URL=ws://localhost:8888
VITE_DEBUG=true
EOF

# Run frontend development server
npm run dev
```

Frontend will run on `http://localhost:5173`

#### 4. Database Setup (Optional)

If you want to use PostgreSQL instead of SQLite:

```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Create database and user
sudo -u postgres psql << EOF
CREATE DATABASE nas_dashboard_dev;
CREATE USER nas_user WITH PASSWORD 'dev_password';
GRANT ALL PRIVILEGES ON DATABASE nas_dashboard_dev TO nas_user;
EOF
```

### Development Features

- **Hot Reload**: Automatic frontend reloading on code changes
- **Debug Mode**: Detailed logging and error messages
- **CORS Enabled**: Cross-origin requests allowed
- **Mock Data**: Optional mock data for testing

---

## Docker Deployment

### Quick Start with Docker Compose

#### 1. Clone Repository

```bash
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard
```

#### 2. Configure Environment

```bash
# Create environment file
cp .env.example .env

# Edit environment variables
nano .env
```

**Essential Configuration:**

```bash
# Backend
PORT=8888
GIN_MODE=release
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h

# Database
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DB=nas_dashboard
POSTGRES_USER=nas_user
POSTGRES_PASSWORD=your-secure-db-password

# CORS
CORS_ORIGIN=https://your-domain.com

# Logging
LOG_LEVEL=info
```

#### 3. Start Services

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Check status
docker-compose ps
```

#### 4. Access Application

- **Frontend**: `http://localhost:3000`
- **Backend API**: `http://localhost:8888/api`
- **Default Credentials**: 
  - Username: `admin`
  - Password: `admin123`

#### 5. Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v
```

### Docker Compose Configuration

**docker-compose.yml:**

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14-alpine
    container_name: nas-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    networks:
      - nas-network

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: nas-backend
    restart: unless-stopped
    environment:
      PORT: 8888
      GIN_MODE: release
      JWT_SECRET: ${JWT_SECRET}
      POSTGRES_HOST: postgres
      POSTGRES_PORT: 5432
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    depends_on:
      - postgres
    ports:
      - "8888:8888"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./backend/config:/app/config
    networks:
      - nas-network

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    container_name: nas-frontend
    restart: unless-stopped
    ports:
      - "3000:80"
    depends_on:
      - backend
    networks:
      - nas-network

  nginx:
    image: nginx:alpine
    container_name: nas-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    depends_on:
      - frontend
      - backend
    networks:
      - nas-network

volumes:
  postgres_data:

networks:
  nas-network:
    driver: bridge
```

### Docker Build

#### Backend Dockerfile

```dockerfile
# backend/Dockerfile
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o nas-dashboard cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates docker

WORKDIR /root/

COPY --from=builder /app/nas-dashboard .

EXPOSE 8888

CMD ["./nas-dashboard"]
```

#### Frontend Dockerfile

```dockerfile
# frontend/Dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy source code
COPY . .

# Build application
RUN npm run build

# Nginx stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

---

## Production Deployment

### Production Server Setup

#### 1. Server Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y curl git nginx ufw

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Add user to docker group
sudo usermod -aG docker $USER
```

#### 2. Configure Firewall

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

#### 3. Set up SSL Certificates

**Option A: Let's Encrypt (Free)**

```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal (configured automatically)
sudo certbot renew --dry-run
```

**Option B: Self-Signed Certificate**

```bash
# Create SSL directory
sudo mkdir -p /etc/nginx/ssl

# Generate self-signed certificate
sudo openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout /etc/nginx/ssl/nginx.key \
  -out /etc/nginx/ssl/nginx.crt
```

#### 4. Configure Nginx

```nginx
# /etc/nginx/sites-available/nas-dashboard
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8888/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket support
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # WebSocket
    location /api/monitor/ws {
        proxy_pass http://localhost:8888/api/monitor/ws;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

Enable site:

```bash
sudo ln -s /etc/nginx/sites-available/nas-dashboard /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

#### 5. Deploy Application

```bash
# Clone repository
git clone https://github.com/yourusername/nas-dashboard.git
cd nas-dashboard

# Configure environment
cp .env.example .env
nano .env

# Start services
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Manual Production Build

#### Backend Build

```bash
cd backend

# Build for production
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nas-dashboard cmd/server/main.go

# Create systemd service
sudo nano /etc/systemd/system/nas-dashboard.service
```

**Systemd Service:**

```ini
[Unit]
Description=NAS Dashboard Backend
After=network.target postgresql.service

[Service]
Type=simple
User=nasuser
WorkingDirectory=/opt/nas-dashboard/backend
ExecStart=/opt/nas-dashboard/backend/nas-dashboard
Restart=always
RestartSec=5

Environment="PORT=8888"
Environment="GIN_MODE=release"
Environment="JWT_SECRET=your-production-secret"
Environment="POSTGRES_HOST=localhost"
Environment="POSTGRES_PORT=5432"

[Install]
WantedBy=multi-user.target
```

Enable service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable nas-dashboard
sudo systemctl start nas-dashboard
sudo systemctl status nas-dashboard
```

#### Frontend Build

```bash
cd frontend

# Install dependencies
npm install

# Build for production
npm run build

# Copy to nginx
sudo cp -r dist/* /var/www/html/nas-dashboard/
```

---

## Security Configuration

### Environment Variables

**Critical Security Variables:**

```bash
# JWT Secret (MUST be changed in production)
JWT_SECRET=your-super-secure-random-jwt-secret-key-min-32-characters

# Database Password (MUST be strong)
POSTGRES_PASSWORD=your-super-secure-db-password-min-16-characters

# CORS Origin (set to your domain only)
CORS_ORIGIN=https://your-domain.com
```

### Generate Secure Secrets

```bash
# Generate JWT Secret
openssl rand -base64 32

# Generate Database Password
openssl rand -base64 16
```

### Firewall Configuration

```bash
# Allow only necessary ports
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS

# Deny all other incoming traffic
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Enable firewall
sudo ufw enable
```

### SSL/TLS Configuration

**Strong SSL Configuration:**

```nginx
ssl_protocols TLSv1.2 TLSv1.3;
ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
ssl_prefer_server_ciphers on;
ssl_session_cache shared:SSL:10m;
ssl_session_timeout 10m;
```

### Application Security

#### 1. Change Default Password

```bash
# Login to dashboard
# Navigate to User Management
# Change admin password immediately
```

#### 2. Enable Rate Limiting

Configure nginx rate limiting:

```nginx
http {
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    
    server {
        location /api/ {
            limit_req zone=api burst=20;
            # ...
        }
    }
}
```

#### 3. Enable Security Headers

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
```

---

## Performance Optimization

### Database Optimization

#### 1. PostgreSQL Configuration

```postgresql
# /etc/postgresql/14/main/postgresql.conf

# Memory Settings
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
work_mem = 16MB

# Query Optimization
random_page_cost = 1.1
effective_io_concurrency = 200

# Connection Settings
max_connections = 100
```

#### 2. Database Indexing

```sql
-- Create indexes for common queries
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_containers_state ON containers(state);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
```

### Backend Optimization

#### 1. Connection Pooling

```go
// Configure database connection pool
db.DB.SetMaxOpenConns(100)
db.DB.SetMaxIdleConns(10)
db.DB.SetConnMaxLifetime(time.Hour)
```

#### 2. Enable Caching

```go
// In-memory cache for frequently accessed data
cache := NewCache(5 * time.Minute)
```

### Frontend Optimization

#### 1. Enable Compression

```nginx
gzip on;
gzip_vary on;
gzip_min_length 1024;
gzip_types text/plain text/css text/xml text/javascript 
           application/x-javascript application/xml+rss 
           application/json application/javascript;
```

#### 2. Browser Caching

```nginx
location ~* \.(jpg|jpeg|png|gif|ico|css|js)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

#### 3. CDN Integration

Consider using a CDN for static assets:

```nginx
location /static/ {
    proxy_pass https://your-cdn.com/static/;
    proxy_cache_valid 200 1y;
}
```

---

## Monitoring and Logging

### Application Monitoring

#### 1. Health Check Endpoint

```bash
# Check application health
curl https://your-domain.com/api/health
```

**Response:**

```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "2024-01-01T12:00:00Z",
    "version": "0.1.0"
  }
}
```

#### 2. Monitor Docker Containers

```bash
# View container status
docker-compose ps

# View resource usage
docker stats

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
```

#### 3. System Monitoring

```bash
# Monitor system resources
htop

# Monitor disk usage
df -h

# Monitor network connections
netstat -tulpn
```

### Logging Configuration

#### 1. Backend Logging

```bash
# /var/log/nas-dashboard/
backend.log          # Application logs
access.log           # API access logs
error.log            # Error logs
```

#### 2. Log Rotation

```bash
# /etc/logrotate.d/nas-dashboard
/var/log/nas-dashboard/*.log {
    daily
    rotate 14
    compress
    delaycompress
    missingok
    notifempty
    create 0640 www-data www-data
    sharedscripts
    postrotate
        docker-compose exec backend kill -USR1 1
    endscript
}
```

#### 3. Centralized Logging

Consider using ELK Stack or similar:

```yaml
# Add to docker-compose.yml
logstash:
  image: logstash:latest
  volumes:
    - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
  ports:
    - "5000:5000"
  networks:
    - nas-network
```

---

## Backup and Recovery

### Database Backup

#### 1. Automated Backups

```bash
#!/bin/bash
# /usr/local/bin/backup-db.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/postgres"
DB_NAME="nas_dashboard"

# Create backup directory
mkdir -p $BACKUP_DIR

# Backup database
docker exec nas-postgres pg_dump -U nas_user $DB_NAME > $BACKUP_DIR/backup_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/backup_$DATE.sql

# Keep only last 7 days
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +7 -delete

echo "Backup completed: backup_$DATE.sql.gz"
```

#### 2. Scheduled Backups

```bash
# Add to crontab
crontab -e

# Daily backup at 2 AM
0 2 * * * /usr/local/bin/backup-db.sh
```

#### 3. Restore Database

```bash
# Restore from backup
gunzip < /backup/postgres/backup_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i nas-postgres psql -U nas_user nas_dashboard
```

### Application Backup

```bash
# Backup configuration
tar -czf /backup/config_$(date +%Y%m%d).tar.gz \
  /opt/nas-dashboard/.env \
  /opt/nas-dashboard/nginx.conf

# Backup uploaded files
tar -czf /backup/uploads_$(date +%Y%m%d).tar.gz \
  /opt/nas-dashboard/uploads
```

### Disaster Recovery

#### 1. Server Failure Recovery

```bash
# On new server
# 1. Install dependencies
# 2. Clone repository
# 3. Restore configuration
tar -xzf /backup/config_YYYYMMDD.tar.gz -C /opt/nas-dashboard/

# 4. Restore database
gunzip < /backup/postgres/backup_YYYYMMDD_HHMMSS.sql.gz | \
  docker exec -i nas-postgres psql -U nas_user nas_dashboard

# 5. Restart services
docker-compose up -d
```

#### 2. Data Recovery Plan

1. **Regular Backups**: Daily automated backups
2. **Off-site Storage**: Store backups in different location
3. **Test Restores**: Regularly test backup restoration
4. **Documentation**: Keep recovery procedures updated

---

## Troubleshooting

### Common Issues

#### 1. Backend Not Starting

**Problem**: Backend container exits immediately

**Solutions**:

```bash
# Check logs
docker-compose logs backend

# Common issues:
# - Database not ready: Wait for postgres to start
# - Port already in use: Check if port 8888 is available
# - Environment variables missing: Verify .env file
```

#### 2. Frontend Not Connecting to Backend

**Problem**: Frontend shows API connection errors

**Solutions**:

```bash
# Check backend is running
docker-compose ps

# Check CORS settings
grep CORS_ORIGIN .env

# Check firewall
sudo ufw status

# Test API directly
curl http://localhost:8888/api/health
```

#### 3. WebSocket Connection Issues

**Problem**: Real-time updates not working

**Solutions**:

```bash
# Check WebSocket URL
grep VITE_WS_URL frontend/.env

# Verify nginx WebSocket configuration
sudo nginx -t

# Check WebSocket upgrade headers
curl -i -N -H "Connection: Upgrade" -H "Upgrade: websocket" \
  http://localhost:8888/api/monitor/ws
```

#### 4. Database Connection Issues

**Problem**: Backend cannot connect to database

**Solutions**:

```bash
# Check database is running
docker-compose ps postgres

# Check database logs
docker-compose logs postgres

# Verify connection settings
grep POSTGRES_ .env

# Test connection
docker exec -it nas-postgres psql -U nas_user -d nas_dashboard
```

#### 5. High Memory Usage

**Problem**: Application using excessive memory

**Solutions**:

```bash
# Check container resource usage
docker stats

# Limit container memory
# Add to docker-compose.yml:
services:
  backend:
    mem_limit: 512m

# Optimize database configuration
# Reduce PostgreSQL shared_buffers if needed
```

### Debug Mode

#### Enable Debug Logging

```bash
# Add to .env
LOG_LEVEL=debug
GIN_MODE=debug

# Restart services
docker-compose restart backend
```

#### View Detailed Logs

```bash
# Follow all logs
docker-compose logs -f

# Follow specific service logs
docker-compose logs -f backend
docker-compose logs -f postgres

# View last 100 lines
docker-compose logs --tail=100 backend
```

### Performance Issues

#### High CPU Usage

```bash
# Check system load
top
htop

# Identify heavy processes
ps aux --sort=-%cpu | head -20

# Check container resource usage
docker stats --no-stream
```

#### Slow Database Queries

```bash
# Connect to database
docker exec -it nas-postgres psql -U nas_user -d nas_dashboard

# Check slow queries
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;

# Analyze query performance
EXPLAIN ANALYZE your_query_here;
```

---

## Maintenance

### Regular Maintenance Tasks

#### Daily
- Check application logs for errors
- Verify backup completion
- Monitor system resources

#### Weekly
- Review security logs
- Check disk space usage
- Test backup restoration

#### Monthly
- Update dependencies
- Review and optimize database
- Security audit
- Performance review

### Updates

#### Application Updates

```bash
# Pull latest changes
git pull origin main

# Update Docker images
docker-compose pull

# Restart services
docker-compose up -d

# Verify update
curl http://localhost:8888/api/health
```

#### System Updates

```bash
# Update system packages
sudo apt update
sudo apt upgrade -y

# Update Docker
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io

# Restart services if needed
docker-compose restart
```

---

## Conclusion

This deployment guide covers all aspects of deploying NAS Dashboard from development to production. Following these guidelines ensures a secure, performant, and maintainable deployment.

### Next Steps

1. **Choose Deployment Method**: Docker or Manual
2. **Configure Security**: SSL certificates and firewall
3. **Set Up Monitoring**: Health checks and logging
4. **Configure Backups**: Automated backup system
5. **Test Deployment**: Verify all features work

### Support

For issues and questions:
- Check logs: `docker-compose logs`
- Review configuration: `.env` file
- Test connectivity: `curl http://localhost:8888/api/health`
- Consult troubleshooting section

---

**Last Updated**: 2026-06-12  
**Version**: 0.1.0  
**Status**: Production Ready
