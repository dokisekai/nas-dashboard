# NAS Dashboard - Deployment Guide

## Prerequisites

- Docker and Docker Compose installed
- Go 1.23+ (for local backend development)
- Node.js 22+ (for local frontend development)
- Linux system with system directory access
- At least 2GB RAM and 10GB disk space

## Environment Configuration

### Frontend Environment Variables

Create `.env` file in `frontend/` directory:

```bash
# API Configuration
VITE_API_URL=http://192.168.50.10:8888
VITE_WS_URL=ws://192.168.50.10:8888
VITE_DEBUG=false
```

### Backend Environment Variables

Create `.env` file in `backend/` directory:

```bash
# Server Configuration
GIN_MODE=release
PORT=8888
HOST=0.0.0.0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this
JWT_EXPIRATION=24h

# Default User (for first login)
DEFAULT_USERNAME=admin
DEFAULT_PASSWORD=admin123

# Database Configuration (if using database)
DB_HOST=localhost
DB_PORT=5432
DB_USER=nas_dashboard
DB_PASSWORD=nas_dashboard
DB_NAME=nas_dashboard

# Monitoring Configuration
MONITOR_INTERVAL=2s
CPU_HISTORY_SIZE=60
MEMORY_HISTORY_SIZE=60
DISK_HISTORY_SIZE=60
NETWORK_HISTORY_SIZE=60
```

## Local Development Setup

### Backend Development

```bash
# Navigate to backend directory
cd backend

# Install dependencies
go mod download

# Run development server
go run cmd/server/main.go

# Or use air for hot reload
air
```

### Frontend Development

```bash
# Navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Docker Deployment

### Quick Start

```bash
# Clone repository
git clone <repository-url>
cd nas-dashboard

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Individual Service Deployment

#### Backend Deployment

```bash
# Build backend image
docker build -t nas-backend backend/

# Run backend container
docker run -d \
  --name nas-backend \
  --restart unless-stopped \
  -p 8888:8888 \
  -v /:/host:ro \
  -v /var/run/docker.sock:/var/run/docker.sock:ro \
  -e GIN_MODE=release \
  nas-backend
```

#### Frontend Deployment

```bash
# Build frontend image
docker build -t nas-frontend frontend/

# Run frontend container
docker run -d \
  --name nas-frontend \
  --restart unless-stopped \
  -p 3000:80 \
  nas-frontend
```

## Production Deployment

### 1. Prepare Environment

```bash
# Create production directory
sudo mkdir -p /opt/nas-dashboard
sudo chown $USER:$USER /opt/nas-dashboard

# Copy files
cp -r . /opt/nas-dashboard/
cd /opt/nas-dashboard
```

### 2. Configure Environment

```bash
# Edit backend environment
nano backend/.env

# Edit frontend environment
nano frontend/.env

# Update docker-compose if needed
nano docker-compose.yml
```

### 3. Build and Deploy

```bash
# Build images
docker-compose build

# Start services
docker-compose up -d

# Check status
docker-compose ps
```

### 4. Verify Deployment

```bash
# Check backend health
curl http://localhost:8888/api/health

# Check frontend
curl http://localhost:3000

# View logs
docker-compose logs -f backend
docker-compose logs -f frontend
```

## SSL/HTTPS Configuration

### Using Nginx Reverse Proxy

```bash
# Install Nginx
sudo apt update
sudo apt install nginx certbot python3-certbot-nginx

# Create Nginx configuration
sudo nano /etc/nginx/sites-available/nas-dashboard
```

Nginx configuration:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /api {
        proxy_pass http://localhost:8888;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    location /ws {
        proxy_pass http://localhost:8888;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/nas-dashboard /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Restart Nginx
sudo systemctl restart nginx

# Obtain SSL certificate
sudo certbot --nginx -d your-domain.com
```

## Monitoring and Maintenance

### Health Checks

```bash
# Backend health
curl http://localhost:8888/api/health

# WebSocket connection
wscat -c ws://localhost:8888/ws/monitor

# Container status
docker ps

# Resource usage
docker stats
```

### Log Management

```bash
# View logs
docker-compose logs -f

# Specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

# Log rotation (add to docker-compose.yml)
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### Backup and Restore

```bash
# Backup configuration
tar -czf nas-dashboard-backup-$(date +%Y%m%d).tar.gz \
  backend/.env \
  frontend/.env \
  docker-compose.yml

# Backup user data (if using database)
docker exec nas-backend pg_dump -U nas_dashboard nas_dashboard > backup.sql

# Restore
docker exec -i nas-backend psql -U nas_dashboard nas_dashboard < backup.sql
```

## Troubleshooting

### Common Issues

#### 1. Backend Won't Start

```bash
# Check logs
docker-compose logs backend

# Common fixes:
# - Check .env file exists and is correct
# - Verify port 8888 is not in use
# - Check system directory permissions
# - Verify JWT_SECRET is set
```

#### 2. Frontend Won't Start

```bash
# Check logs
docker-compose logs frontend

# Common fixes:
# - Verify backend is running
# - Check VITE_API_URL and VITE_WS_URL
# - Ensure build completed successfully
# - Check nginx configuration
```

#### 3. WebSocket Connection Fails

```bash
# Test WebSocket connection
wscat -c ws://localhost:8888/ws/monitor

# Common fixes:
# - Check WebSocket URL in frontend .env
# - Verify backend WebSocket handler
# - Check firewall rules
# - Ensure proper authentication
```

#### 4. Permission Denied Errors

```bash
# Fix system directory permissions
sudo chmod -R 755 /host
sudo chown -R root:root /host

# Fix Docker socket permissions
sudo chmod 666 /var/run/docker.sock
```

### Performance Issues

#### High Memory Usage

```bash
# Check resource usage
docker stats

# Optimize containers
# - Add memory limits to docker-compose.yml
# - Reduce monitoring history size
# - Implement data pagination
```

#### Slow Response Times

```bash
# Check CPU usage
top

# Optimize performance
# - Enable production mode
# - Implement caching
# - Use CDN for static assets
# - Optimize database queries
```

## Security Hardening

### 1. Update Default Credentials

```bash
# Change default admin password immediately
# Log in and change password in user settings
```

### 2. Firewall Configuration

```bash
# Configure UFW
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 8888/tcp  # If accessing backend directly
sudo ufw enable
```

### 3. Rate Limiting

```bash
# Add to Nginx configuration
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;

location /api {
    limit_req zone=api_limit burst=20 nodelay;
    proxy_pass http://localhost:8888;
}
```

### 4. Security Headers

```nginx
# Add to Nginx configuration
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "no-referrer-when-downgrade" always;
```

## Updates and Upgrades

### Update Application

```bash
# Pull latest changes
git pull

# Rebuild containers
docker-compose build

# Restart services
docker-compose down
docker-compose up -d
```

### Update Dependencies

```bash
# Frontend dependencies
cd frontend
npm update
npm audit fix

# Backend dependencies
cd backend
go get -u ./...
go mod tidy
```

## Scaling Considerations

### Horizontal Scaling

```yaml
# Update docker-compose.yml for multiple instances
services:
  backend:
    deploy:
      replicas: 3

  frontend:
    deploy:
      replicas: 2

  # Add load balancer
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
```

### Database Scaling

```bash
# For production, consider external database
# - PostgreSQL cluster
# - Redis for caching
# - TimescaleDB for metrics
```

## Support and Resources

### Documentation

- [Backend Documentation](./backend/README.md)
- [Frontend Documentation](./frontend/README.md)
- [API Documentation](./docs/API.md)
- [Development Guide](./docs/DEVELOPMENT.md)

### Useful Commands

```bash
# Quick restart
docker-compose restart

# Full rebuild
docker-compose down && docker-compose up -d --build

# Clean restart
docker-compose down -v && docker-compose up -d

# Monitor resources
docker stats --no-stream

# Enter container
docker exec -it nas-backend bash
docker exec -it nas-frontend sh
```

---

**Last Updated:** 2025-06-12
**Version:** 1.0.0
