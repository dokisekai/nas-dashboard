# NAS Dashboard - Troubleshooting Guide

Comprehensive troubleshooting guide for common issues and their solutions.

## 📋 Table of Contents

1. [Quick Diagnostics](#quick-diagnostics)
2. [Installation Issues](#installation-issues)
3. [Authentication Problems](#authentication-problems)
4. [API Connection Issues](#api-connection-issues)
5. [WebSocket Problems](#websocket-problems)
6. [Monitoring Data Issues](#monitoring-data-issues)
7. [Performance Problems](#performance-problems)
8. [Docker Integration Issues](#docker-integration-issues)
9. [Frontend Issues](#frontend-issues)
10. [Debug Mode](#debug-mode)
11. [Log Analysis](#log-analysis)

---

## Quick Diagnostics

### Health Check Script

```bash
#!/bin/bash
# health-check.sh

echo "=== NAS Dashboard Health Check ==="
echo ""

# Check backend service
echo "1. Backend Service:"
if systemctl is-active --quiet nas-dashboard; then
    echo "   ✓ Backend service is running"
else
    echo "   ✗ Backend service is NOT running"
fi

# Check nginx
echo "2. Web Server:"
if systemctl is-active --quiet nginx; then
    echo "   ✓ Nginx is running"
else
    echo "   ✗ Nginx is NOT running"
fi

# Check backend port
echo "3. Backend Port (8888):"
if nc -zv localhost 8888 2>&1 | grep -q succeeded; then
    echo "   ✓ Port 8888 is accessible"
else
    echo "   ✗ Port 8888 is NOT accessible"
fi

# Check Docker
echo "4. Docker Service:"
if systemctl is-active --quiet docker; then
    echo "   ✓ Docker is running"
else
    echo "   ✗ Docker is NOT running"
fi

# Check disk space
echo "5. Disk Space:"
df -h | grep -E "(/$|/home)"

# Check memory
echo "6. Memory Usage:"
free -h

# Check backend logs for errors
echo "7. Recent Backend Errors:"
journalctl -u nas-dashboard --since "5 minutes ago" | grep -i error || echo "   No recent errors"

echo ""
echo "=== Health Check Complete ==="
```

### Quick Tests

```bash
# Test backend API
curl http://localhost:8888/api/health

# Test frontend
curl -I http://localhost/

# Test WebSocket
wscat -c ws://localhost:8888/api/monitor/ws

# Test Docker access
docker ps
```

---

## Installation Issues

### Issue: Backend Service Won't Start

**Symptoms:**
- Service fails immediately after start
- Systemctl shows "failed" status
- No log files created

**Diagnosis:**
```bash
# Check service status
sudo systemctl status nas-dashboard

# View detailed logs
sudo journalctl -u nas-dashboard -n 50

# Test manual start
cd /opt/nas-dashboard
sudo ./nas-dashboard
```

**Common Causes and Solutions:**

#### 1. Port Already in Use

```bash
# Check what's using port 8888
sudo netstat -tulpn | grep 8888
sudo lsof -i :8888

# Solution: Kill the process or change port
sudo kill -9 <PID>
# Or edit config.env and change PORT=8889
```

#### 2. Missing Dependencies

```bash
# Check Go installation
go version

# Rebuild application
cd /opt/nas-dashboard
go build -o nas-dashboard cmd/server/main.go

# Check for missing system libraries
ldd nas-dashboard
```

#### 3. Permission Issues

```bash
# Check file permissions
ls -la /opt/nas-dashboard

# Fix permissions
sudo chown -R root:root /opt/nas-dashboard
sudo chmod -R 755 /opt/nas-dashboard

# Check systemd service file permissions
ls -la /etc/systemd/system/nas-dashboard.service
```

#### 4. Configuration Errors

```bash
# Validate environment file
cat /etc/nas-dashboard/config.env

# Test configuration load
sudo -u nas-dashboard /opt/nas-dashboard/nas-dashboard --validate-config
```

### Issue: Frontend Shows Blank Page

**Symptoms:**
- White screen when accessing application
- Browser console shows errors
- No UI elements rendered

**Diagnosis:**
```bash
# Check nginx configuration
sudo nginx -t

# Check nginx logs
sudo tail -f /var/log/nginx/error.log

# Verify frontend files exist
ls -la /var/www/nas-dashboard/

# Check file permissions
ls -la /var/www/nas-dashboard/
```

**Solutions:**

#### 1. Missing Frontend Files

```bash
# Rebuild frontend
cd /home/hserver/nas-dashboard/frontend
npm run build

# Deploy files
sudo cp -r dist/* /var/www/nas-dashboard/
sudo chown -R www-data:www-data /var/www/nas-dashboard
```

#### 2. Nginx Configuration Error

```bash
# Check nginx configuration
sudo nginx -t

# Test configuration
sudo nginx -t -c /etc/nginx/nginx.conf

# Reload nginx
sudo systemctl reload nginx
```

#### 3. Browser Cache Issue

```javascript
// Clear browser cache
// Ctrl+Shift+Delete (Chrome/Firefox)
// Or use private browsing mode to test
```

### Issue: Docker Deployment Fails

**Symptoms:**
- Containers won't start
- Docker compose exits with errors
- Network issues between containers

**Diagnosis:**
```bash
# Check Docker status
sudo systemctl status docker

# Check Docker logs
sudo journalctl -u docker -n 50

# Test Docker
docker run hello-world
```

**Solutions:**

#### 1. Docker Socket Permission

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Log out and log back in
# Or use newgrp
newgrp docker

# Verify
docker ps
```

#### 2. Network Issues

```bash
# Check Docker networks
docker network ls

# Inspect network
docker network inspect nas-network

# Recreate network
docker network rm nas-network
docker network create nas-network
```

#### 3. Volume Mount Issues

```bash
# Check volume mounts
docker volume ls

# Inspect volumes
docker volume inspect nas-data

# Fix permissions
sudo chown -R 1000:1000 /path/to/mount/point
```

---

## Authentication Problems

### Issue: Login Always Fails

**Symptoms:**
- "Invalid credentials" error with correct credentials
- No error message, just redirect to login
- Console shows 401 errors

**Diagnosis:**
```bash
# Check backend logs
sudo journalctl -u nas-dashboard -f | grep -i auth

# Test login endpoint directly
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

**Solutions:**

#### 1. Check User Data

```bash
# For hardcoded users (development)
grep -r "admin123" backend/

# For database users (production)
# Check database connection and user table
```

#### 2. Verify JWT Secret

```bash
# Check JWT secret in config
cat /etc/nas-dashboard/config.env | grep JWT_SECRET

# Ensure it's not the default
# Should be a long, random string
```

#### 3. Test Token Generation

```bash
# Generate test token
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -v

# Should return JWT token
```

### Issue: Token Expire Too Quickly

**Symptoms:**
- Frequent logouts
- 401 errors after short time
- Need to constantly re-login

**Diagnosis:**
```bash
# Check token duration settings
cat /etc/nas-dashboard/config.env | grep DURATION

# Decode JWT to check expiration
# Use jwt.io or jwt-cli
```

**Solutions:**

#### 1. Adjust Token Duration

```bash
# Edit configuration
sudo nano /etc/nas-dashboard/config.env

# Set longer durations
JWT_ACCESS_DURATION=24h
JWT_REFRESH_DURATION=720h

# Restart service
sudo systemctl restart nas-dashboard
```

#### 2. Implement Auto-Refresh

```typescript
// Frontend: Setup token refresh
axios.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      await authStore.refreshToken()
      return axios.request(error.config)
    }
    return Promise.reject(error)
  }
)
```

### Issue: Can't Logout

**Symptoms:**
- Logout button doesn't work
- Session persists after logout
- Redirect doesn't happen

**Diagnosis:**
```bash
# Check browser console for errors
# Check Network tab for logout request
```

**Solutions:**

#### 1. Check Logout Endpoint

```bash
# Test logout endpoint
curl -X POST http://localhost:8888/api/auth/logout \
  -H "Authorization: Bearer <token>"
```

#### 2. Clear Browser Storage

```javascript
// Manual logout
localStorage.clear()
sessionStorage.clear()
// Reload page
location.reload()
```

---

## API Connection Issues

### Issue: Frontend Can't Connect to Backend

**Symptoms:**
- "Network Error" in browser
- CORS errors in console
- Connection refused

**Diagnosis:**
```bash
# Check backend is running
sudo systemctl status nas-dashboard

# Check backend is listening
sudo netstat -tulpn | grep 8888

# Test API directly
curl http://localhost:8888/api/health
```

**Solutions:**

#### 1. Backend Not Running

```bash
# Start backend
sudo systemctl start nas-dashboard

# Check logs
sudo journalctl -u nas-dashboard -f
```

#### 2. CORS Configuration

```bash
# Check CORS settings
cat /etc/nas-dashboard/config.env | grep CORS

# Ensure origin matches frontend
CORS_ORIGIN=http://localhost:5173  # dev
CORS_ORIGIN=https://your-domain.com  # prod
```

#### 3. Firewall Blocking

```bash
# Check firewall rules
sudo ufw status

# Allow backend port
sudo ufw allow 8888/tcp

# Check if port is accessible
telnet localhost 8888
```

### Issue: API Returns 500 Errors

**Symptoms:**
- Internal server errors
- No specific error message
- Random failures

**Diagnosis:**
```bash
# Check backend logs
sudo journalctl -u nas-dashboard -n 100

# Look for stack traces
sudo journalctl -u nas-dashboard | grep -A 20 "panic\|error"
```

**Solutions:**

#### 1. Check Backend Logs

```bash
# Follow logs in real-time
sudo journalctl -u nas-dashboard -f

# Look for specific errors
sudo journalctl -u nas-dashboard | grep -i "error\|panic"
```

#### 2. Test Specific Endpoints

```bash
# Test each endpoint
curl http://localhost:8888/api/monitor/cpu
curl http://localhost:8888/api/monitor/memory
curl http://localhost:8888/api/docker/containers
```

#### 3. Check System Resources

```bash
# Check memory
free -h

# Check disk space
df -h

# Check CPU usage
top
```

---

## WebSocket Problems

### Issue: WebSocket Connection Fails

**Symptoms:**
- Real-time monitoring doesn't update
- "WebSocket connection failed" in console
- Falls back to polling (if implemented)

**Diagnosis:**
```bash
# Check WebSocket endpoint
curl -i -N \
  -H "Connection: Upgrade" \
  -H "Upgrade: websocket" \
  -H "Sec-WebSocket-Version: 13" \
  -H "Sec-WebSocket-Key: test" \
  http://localhost:8888/api/monitor/ws

# Check browser console
# F12 → Console → Filter by "WebSocket"
```

**Solutions:**

#### 1. Check WebSocket Authentication

```javascript
// Ensure token is included
const token = localStorage.getItem('access_token')
const ws = new WebSocket(`ws://localhost:8888/api/monitor/ws?token=${token}`)
```

#### 2. Verify nginx WebSocket Configuration

```nginx
location /api/monitor/ws {
    proxy_pass http://localhost:8888;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
    proxy_set_header Host $host;
    proxy_read_timeout 86400;
}
```

#### 3. Check Firewall

```bash
# Ensure WebSocket upgrade is allowed
sudo ufw allow 8888/tcp

# Test WebSocket connection
wscat -c ws://localhost:8888/api/monitor/ws?token=<token>
```

### Issue: WebSocket Disconnects Frequently

**Symptoms:**
- Connection drops randomly
- Need to manually refresh
- Unstable real-time updates

**Diagnosis:**
```bash
# Check network stability
ping -c 100 localhost

# Check for WebSocket errors
sudo journalctl -u nas-dashboard | grep -i websocket
```

**Solutions:**

#### 1. Implement Auto-Reconnect

```javascript
class ReconnectingWebSocket {
  constructor(url) {
    this.url = url
    this.reconnectAttempts = 0
    this.maxAttempts = 5
    this.connect()
  }
  
  connect() {
    this.ws = new WebSocket(this.url)
    
    this.ws.onclose = () => {
      if (this.reconnectAttempts < this.maxAttempts) {
        setTimeout(() => {
          this.reconnectAttempts++
          this.connect()
        }, 5000)
      }
    }
    
    this.ws.onopen = () => {
      this.reconnectAttempts = 0
    }
  }
}
```

#### 2. Check Server Timeouts

```bash
# Increase timeout in nginx
proxy_read_timeout 86400;
proxy_send_timeout 86400;
```

---

## Monitoring Data Issues

### Issue: CPU Shows Wrong Values

**Symptoms:**
- CPU shows 0.016% instead of 1.6%
- Values always 0 or 100
- Negative percentages

**Diagnosis:**
```bash
# Check backend CPU data
curl http://localhost:8888/api/monitor/cpu

# Check system CPU
top -b -n 1 | grep "Cpu(s)"
```

**Solutions:**

#### 1. Fix Frontend Display

```typescript
// Current (wrong)
<p>{{ cpuInfo.usage }}%</p>

// Fixed
<p>{{ (cpuInfo.usage * 100).toFixed(1) }}%</p>
```

#### 2. Verify Backend Data

```bash
# Check backend returns correct format
curl http://localhost:8888/api/monitor/cpu | jq

# Should return usage as decimal (0.0-1.0)
```

### Issue: Network Monitoring Shows Zero

**Symptoms:**
- Network statistics always 0
- No traffic data displayed
- Interfaces not showing

**Diagnosis:**
```bash
# Check network data
curl http://localhost:8888/api/monitor/network

# Check system interfaces
ip addr show
```

**Solutions:**

#### 1. Fix Data Format Mismatch

```typescript
// Backend sends: interfaces[].bytesRecv
// Frontend expects: network.rx_bytes

// Fix in Dashboard.vue
const down = data.network.interfaces.reduce((sum, iface) => 
  sum + (iface.bytesRecv || 0), 0
)
```

#### 2. Filter Virtual Interfaces

```javascript
// Exclude docker0, virbr0, veth*
const interfaces = data.network.interfaces.filter(
  i => !i.name.startsWith('docker0') && 
       !i.name.startsWith('virbr') &&
       !i.name.startsWith('veth')
)
```

### Issue: Disk Usage Incorrect

**Symptoms:**
- Wrong disk usage percentages
- Missing disks
- Duplicate entries

**Diagnosis:**
```bash
# Check disk data
curl http://localhost:8888/api/monitor/disk

# Verify system disks
df -h
```

**Solutions:**

#### 1. Check Mount Points

```bash
# Ensure all filesystems are included
df -T | grep -v tmpfs
```

#### 2. Fix Data Parsing

```typescript
// Ensure correct data mapping
disks.map(disk => ({
  device: disk.device,
  mount: disk.mount_point,
  total: formatBytes(disk.total),
  used: formatBytes(disk.used),
  percent: disk.percent.toFixed(1)
}))
```

---

## Performance Problems

### Issue: High Memory Usage

**Symptoms:**
- Application uses lots of RAM
- System becomes slow
- OOM errors

**Diagnosis:**
```bash
# Check memory usage
free -h
ps aux | grep nas-dashboard

# Check memory leaks
valgrind --leak-check=full ./nas-dashboard
```

**Solutions:**

#### 1. Adjust Monitor Interval

```bash
# Reduce update frequency
echo "MONITOR_INTERVAL=5s" >> /etc/nas-dashboard/config.env
sudo systemctl restart nas-dashboard
```

#### 2. Optimize Data Collection

```go
// Cache system metrics
type MonitorCache struct {
    data      interface{}
    timestamp time.Time
    ttl       time.Duration
}

func (mc *MonitorCache) Get() interface{} {
    if time.Since(mc.timestamp) < mc.ttl {
        return mc.data
    }
    // Refresh cache
    mc.data = collectData()
    mc.timestamp = time.Now()
    return mc.data
}
```

### Issue: Slow API Responses

**Symptoms:**
- API calls take long time
- Browser shows loading spinner
- Timeouts on some endpoints

**Diagnosis:**
```bash
# Test API response time
time curl http://localhost:8888/api/monitor/all

# Check system load
uptime
```

**Solutions:**

#### 1. Profile Application

```bash
# Use Go profiler
go tool pprof http://localhost:8888/debug/pprof/profile
```

#### 2. Optimize Database Queries

```go
// Use indexes
// Batch operations
// Cache results
```

#### 3. Enable Compression

```nginx
# nginx configuration
gzip on;
gzip_types application/json;
gzip_min_length 1000;
```

---

## Docker Integration Issues

### Issue: Can't Manage Containers

**Symptoms:**
- Container operations fail
- "Permission denied" errors
- Can't see container list

**Diagnosis:**
```bash
# Check Docker access
docker ps

# Check Docker socket permissions
ls -la /var/run/docker.sock
```

**Solutions:**

#### 1. Fix Docker Socket Permission

```bash
# Add user to docker group
sudo usermod -aG docker $USER

# Fix socket permissions
sudo chmod 666 /var/run/docker.sock
```

#### 2. Check Docker Mount

```bash
# Ensure socket is mounted in container
docker run -v /var/run/docker.sock:/var/run/docker.sock ...
```

### Issue: Container List Empty

**Symptoms:**
- No containers shown in UI
- Docker shows containers but API returns empty

**Diagnosis:**
```bash
# Check Docker
docker ps

# Test API
curl http://localhost:8888/api/docker/containers
```

**Solutions:**

#### 1. Check Docker Client

```go
// Ensure Docker client is initialized
client, err := client.NewClientWithOpts(client.FromEnv)
if err != nil {
    log.Fatal(err)
}
```

#### 2. Verify API Response

```bash
# Check response format
curl http://localhost:8888/api/docker/containers | jq
```

---

## Frontend Issues

### Issue: Menu Items Duplicated

**Symptoms:**
- Sidebar appears twice
- Header shows multiple times
- Navigation elements repeated

**Diagnosis:**
```bash
# Check App.vue structure
cat frontend/src/App.vue

# Check routing configuration
cat frontend/src/router/index.ts
```

**Solutions:**

#### 1. Fix Routing Structure

```vue
<!-- App.vue -->
<template>
  <router-view />
</template>

<!-- MainLayout.vue -->
<template>
  <div class="layout">
    <Sidebar />
    <Header />
    <router-view />
  </div>
</template>
```

#### 2. Use Route Guards

```typescript
// router/index.ts
const router = createRouter({
  routes: [
    {
      path: '/dashboard',
      component: MainLayout,
      children: [
        { path: '', component: Dashboard }
      ]
    }
  ]
})
```

### Issue: Charts Not Displaying

**Symptoms:**
- Empty chart areas
- No data visualization
- Canvas elements blank

**Diagnosis:**
```javascript
// Check browser console
// Look for Chart.js errors
```

**Solutions:**

#### 1. Check Chart.js Installation

```bash
npm install chart.js
```

#### 2. Verify Chart Configuration

```typescript
import { Chart } from 'chart.js'

new Chart(ctx, {
  type: 'line',
  data: {
    labels: ['A', 'B', 'C'],
    datasets: [{
      label: 'Data',
      data: [1, 2, 3]
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false
  }
})
```

---

## Debug Mode

### Enable Debug Mode

**Backend:**
```bash
# Edit configuration
sudo nano /etc/nas-dashboard/config.env

# Add debug settings
GIN_MODE=debug
LOG_LEVEL=debug

# Restart service
sudo systemctl restart nas-dashboard
```

**Frontend:**
```bash
# Add to .env
echo "VITE_DEBUG=true" >> frontend/.env

# Rebuild
npm run build
```

### Debug Output

**Backend logs:**
```bash
# Follow debug logs
sudo journalctl -u nas-dashboard -f

# Filter by level
sudo journalctl -u nas-dashboard | grep debug
```

**Frontend console:**
```javascript
// Open browser console (F12)
// Debug mode shows:
// - API requests/responses
// - WebSocket messages
// - Error details
// - Performance timing
```

### Debug Tools

**Backend:**
```bash
# Go profiler
go tool pprof http://localhost:8888/debug/pprof/profile

# Memory profiling
curl http://localhost:8888/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

**Frontend:**
```javascript
// Vue devtools
// Install Vue Devtools browser extension

// React to components
// Check component state
// Debug props and events
```

---

## Log Analysis

### Backend Logs

**Location:**
```bash
# System logs
sudo journalctl -u nas-dashboard

# Application logs
tail -f /var/log/nas-dashboard/app.log
```

**Common Log Patterns:**

```bash
# Authentication failures
grep "Invalid credentials" /var/log/nas-dashboard/auth.log

# API errors
grep "ERROR" /var/log/nas-dashboard/api.log

# WebSocket issues
grep "websocket" /var/log/nas-dashboard/app.log
```

### Frontend Logs

**Browser Console:**
```javascript
// Open DevTools (F12)
// Check Console tab for errors
// Check Network tab for failed requests
```

**Common Patterns:**

```javascript
// 401 errors - Authentication issues
// 500 errors - Server errors
// CORS errors - Network configuration
// WebSocket errors - Real-time monitoring
```

### Log Analysis Tools

```bash
# Real-time monitoring
tail -f /var/log/nginx/error.log

# Error statistics
grep -c "ERROR" /var/log/nas-dashboard/app.log

# Time-based analysis
grep "2024-01-01" /var/log/nas-dashboard/app.log | wc -l

# Pattern matching
grep -i "panic\|fatal" /var/log/nas-dashboard/app.log
```

---

## Getting Help

### When to Ask for Help

- You've tried all troubleshooting steps
- Issue persists after multiple attempts
- You encounter unexpected errors
- You need help understanding logs

### Information to Provide

```bash
# System information
uname -a
go version
node --version
docker --version

# Service status
sudo systemctl status nas-dashboard
sudo systemctl status nginx

# Recent logs
sudo journalctl -u nas-dashboard -n 50

# Configuration
cat /etc/nas-dashboard/config.env
```

### Resources

- **Documentation**: [Project Docs](https://github.com/yourusername/nas-dashboard/docs)
- **Issues**: [GitHub Issues](https://github.com/yourusername/nas-dashboard/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yourusername/nas-dashboard/discussions)

---

## Emergency Procedures

### Quick Reset

```bash
#!/bin/bash
# emergency-reset.sh

echo "=== Emergency Reset ==="

# Stop all services
sudo systemctl stop nas-dashboard
sudo systemctl stop nginx

# Clear caches
sudo rm -rf /var/cache/nginx/*
sudo rm -rf /tmp/nas-dashboard/*

# Restart services
sudo systemctl start nas-dashboard
sudo systemctl start nginx

# Verify
sudo systemctl status nas-dashboard
sudo systemctl status nginx

echo "=== Reset Complete ==="
```

### Rollback

```bash
# Rollback to previous version
git log --oneline -10
git checkout <previous-commit>

# Rebuild
cd backend
go build -o nas-dashboard cmd/server/main.go

cd ../frontend
npm run build

# Deploy
sudo systemctl restart nas-dashboard
```

---

**Last Updated**: 2026-06-12
**Version**: 0.1.0-alpha
