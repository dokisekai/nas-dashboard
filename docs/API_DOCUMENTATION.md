# NAS Dashboard - API Documentation

Complete API reference for the NAS Dashboard backend services, including REST endpoints and WebSocket protocols.

## 📋 Table of Contents

1. [Base URL and Authentication](#base-url-and-authentication)
2. [REST API Endpoints](#rest-api-endpoints)
3. [WebSocket Protocol](#websocket-protocol)
4. [Error Codes](#error-codes)
5. [Data Models](#data-models)

---

## Base URL and Authentication

### Base URLs

```
Development:  http://localhost:8888/api
Production:   https://your-domain.com/api
WebSocket:    ws://localhost:8888/api/monitor/ws
```

### Authentication

All API endpoints (except login) require JWT authentication:

```http
Authorization: Bearer <access_token>
```

#### Token Types

- **Access Token**: Short-lived (24 hours), used for API requests
- **Refresh Token**: Long-lived (30 days), used to obtain new access tokens

#### Authentication Flow

1. **Login**: Obtain tokens from `/api/auth/login`
2. **Use Access Token**: Include in Authorization header
3. **Refresh Token**: When access token expires, use `/api/auth/refresh`
4. **Logout**: Invalidate tokens using `/api/auth/logout`

---

## REST API Endpoints

### Authentication Endpoints

#### POST /api/auth/login
Authenticate user and receive JWT tokens.

**Request:**
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "error": "Invalid credentials"
}
```

#### POST /api/auth/refresh
Refresh access token using refresh token.

**Request:**
```http
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST /api/auth/logout
Invalidate current tokens.

**Request:**
```http
POST /api/auth/logout
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

#### GET /api/auth/me
Get current user information.

**Request:**
```http
GET /api/auth/me
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "username": "admin",
  "role": "admin",
  "created_at": "2024-01-01T00:00:00Z"
}
```

---

### System Monitoring Endpoints

#### GET /api/monitor/cpu
Get CPU usage information.

**Request:**
```http
GET /api/monitor/cpu
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "usage": 0.45,
  "cores": 24,
  "load1": 2.5,
  "load5": 2.1,
  "load15": 1.8,
  "frequency": 3200
}
```

**Field Descriptions:**
- `usage`: CPU usage percentage (0.0-1.0)
- `cores`: Number of CPU cores
- `load1`: 1-minute load average
- `load5`: 5-minute load average
- `load15`: 15-minute load average
- `frequency`: CPU frequency in MHz

#### GET /api/monitor/memory
Get memory usage information.

**Request:**
```http
GET /api/monitor/memory
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "total": 17179869184,
  "used": 8589934592,
  "free": 8589934592,
  "percent": 50.0,
  "swap_total": 4294967296,
  "swap_used": 0,
  "swap_percent": 0.0
}
```

**Field Descriptions:**
- `total`: Total memory in bytes
- `used`: Used memory in bytes
- `free`: Free memory in bytes
- `percent`: Memory usage percentage
- `swap_total`: Total swap space in bytes
- `swap_used`: Used swap space in bytes
- `swap_percent`: Swap usage percentage

#### GET /api/monitor/disk
Get disk usage information.

**Request:**
```http
GET /api/monitor/disk
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "disks": [
    {
      "device": "/dev/sda1",
      "mount_point": "/",
      "file_system": "ext4",
      "total": 1073741824000,
      "used": 536870912000,
      "free": 536870912000,
      "percent": 50.0
    },
    {
      "device": "/dev/sda2",
      "mount_point": "/home",
      "file_system": "ext4",
      "total": 2147483648000,
      "used": 1073741824000,
      "free": 1073741824000,
      "percent": 50.0
    }
  ]
}
```

#### GET /api/monitor/network
Get network interface statistics.

**Request:**
```http
GET /api/monitor/network
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "interfaces": [
    {
      "name": "eth0",
      "bytes_sent": 1073741824,
      "bytes_recv": 2147483648,
      "packets_sent": 1000000,
      "packets_recv": 2000000,
      "err_in": 0,
      "err_out": 0,
      "drop_in": 0,
      "drop_out": 0
    }
  ]
}
```

#### GET /api/monitor/system
Get system information.

**Request:**
```http
GET /api/monitor/system
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "hostname": "nas-server",
  "os": "Linux",
  "platform": "ubuntu",
  "platform_version": "22.04",
  "architecture": "x86_64",
  "kernel_version": "5.15.0-72-generic",
  "uptime": 86400,
  "boot_time": 1704067200,
  "processes": 250
}
```

#### GET /api/monitor/all
Get all monitoring data in single request.

**Request:**
```http
GET /api/monitor/all
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "cpu": { /* CPU data */ },
  "memory": { /* Memory data */ },
  "disk": { /* Disk data */ },
  "network": { /* Network data */ },
  "system": { /* System data */ }
}
```

---

### Docker Management Endpoints

#### GET /api/docker/containers
List all Docker containers.

**Request:**
```http
GET /api/docker/containers
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "containers": [
    {
      "id": "abc123def456",
      "name": "nginx",
      "image": "nginx:latest",
      "state": "running",
      "status": "Up 2 hours",
      "ports": ["80:80", "443:443"],
      "created": 1704067200
    }
  ]
}
```

#### POST /api/docker/containers/{id}/start
Start a Docker container.

**Request:**
```http
POST /api/docker/containers/abc123def456/start
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Container started successfully"
}
```

#### POST /api/docker/containers/{id}/stop
Stop a Docker container.

**Request:**
```http
POST /api/docker/containers/abc123def456/stop
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Container stopped successfully"
}
```

#### POST /api/docker/containers/{id}/restart
Restart a Docker container.

**Request:**
```http
POST /api/docker/containers/abc123def456/restart
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Container restarted successfully"
}
```

#### GET /api/docker/containers/{id}/logs
Get container logs.

**Request:**
```http
GET /api/docker/containers/abc123def456/logs?tail=100
Authorization: Bearer <access_token>
```

**Query Parameters:**
- `tail`: Number of lines from end of logs (default: 100)
- `since`: UNIX timestamp to get logs since
- `timestamps`: Include timestamps (true/false)

**Response (200 OK):**
```json
{
  "logs": "2024-01-01T00:00:00.000Z [INFO] Starting nginx...\n..."
}
```

---

### Storage Management Endpoints

#### GET /api/storage
Get storage usage information.

**Request:**
```http
GET /api/storage
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "storage": [
    {
      "path": "/",
      "device": "/dev/sda1",
      "type": "ext4",
      "total": 1073741824000,
      "used": 536870912000,
      "available": 536870912000,
      "usage_percent": 50.0,
      "mount_options": ["rw", "relatime"]
    }
  ]
}
```

#### POST /api/storage/scan
Scan and update storage information.

**Request:**
```http
POST /api/storage/scan
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Storage scan completed"
}
```

---

### User Management Endpoints

#### GET /api/users
List all users.

**Request:**
```http
GET /api/users
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "users": [
    {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "created_at": "2024-01-01T00:00:00Z",
      "last_login": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### POST /api/users
Create a new user.

**Request:**
```http
POST /api/users
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "user": {
    "id": 2,
    "username": "newuser",
    "email": "newuser@example.com",
    "role": "user",
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

#### PUT /api/users/{id}
Update user information.

**Request:**
```http
PUT /api/users/2
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "email": "updated@example.com",
  "role": "admin"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "user": {
    "id": 2,
    "username": "newuser",
    "email": "updated@example.com",
    "role": "admin",
    "updated_at": "2024-01-01T13:00:00Z"
  }
}
```

#### DELETE /api/users/{id}
Delete a user.

**Request:**
```http
DELETE /api/users/2
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

#### PUT /api/users/{id}/password
Change user password.

**Request:**
```http
PUT /api/users/2/password
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "old_password": "oldpassword",
  "new_password": "newpassword123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

---

### Service Management Endpoints

#### GET /api/services
List all system services.

**Request:**
```http
GET /api/services
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "services": [
    {
      "name": "nginx",
      "description": "A high performance web server",
      "state": "running",
      "status": "active (running)"
    }
  ]
}
```

#### POST /api/services/{name}/start
Start a system service.

**Request:**
```http
POST /api/services/nginx/start
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service started successfully"
}
```

#### POST /api/services/{name}/stop
Stop a system service.

**Request:**
```http
POST /api/services/nginx/stop
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service stopped successfully"
}
```

#### POST /api/services/{name}/restart
Restart a system service.

**Request:**
```http
POST /api/services/nginx/restart
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service restarted successfully"
}
```

---

### Health Check Endpoint

#### GET /api/health
Check API health status.

**Request:**
```http
GET /api/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "0.1.0"
}
```

---

## WebSocket Protocol

### Connection

**Endpoint:** `ws://localhost:8888/api/monitor/ws`

**Authentication:** Include JWT token in query string:

```javascript
const ws = new WebSocket('ws://localhost:8888/api/monitor/ws?token=<access_token>');
```

### Message Format

All WebSocket messages are JSON objects.

#### Server → Client Messages

**Monitoring Data Update:**
```json
{
  "type": "monitor_data",
  "timestamp": 1704067200,
  "data": {
    "cpu": {
      "usage": 0.45,
      "cores": 24,
      "load1": 2.5,
      "load5": 2.1,
      "load15": 1.8
    },
    "memory": {
      "total": 17179869184,
      "used": 8589934592,
      "free": 8589934592,
      "percent": 50.0
    },
    "disk": {
      "disks": [
        {
          "device": "/dev/sda1",
          "mount_point": "/",
          "total": 1073741824000,
          "used": 536870912000,
          "free": 536870912000,
          "percent": 50.0
        }
      ]
    },
    "network": {
      "interfaces": [
        {
          "name": "eth0",
          "bytes_sent": 1073741824,
          "bytes_recv": 2147483648,
          "packets_sent": 1000000,
          "packets_recv": 2000000
        }
      ]
    }
  }
}
```

#### Client → Server Messages

**Subscribe to specific metrics:**
```json
{
  "action": "subscribe",
  "metrics": ["cpu", "memory", "disk", "network"]
}
```

**Unsubscribe from metrics:**
```json
{
  "action": "unsubscribe",
  "metrics": ["disk"]
}
```

**Request refresh:**
```json
{
  "action": "refresh"
}
```

### Connection States

1. **Connecting**: WebSocket handshake in progress
2. **Connected**: Authentication successful, data streaming active
3. **Disconnected**: Connection lost, attempting reconnection
4. **Error**: Authentication failed or connection error

### Reconnection Strategy

```javascript
// Automatic reconnection with backoff
const connect = () => {
  const ws = new WebSocket('ws://localhost:8888/api/monitor/ws?token=<token>');
  
  ws.onclose = (event) => {
    if (!event.wasClean) {
      setTimeout(connect, 5000); // Reconnect after 5 seconds
    }
  };
  
  return ws;
};
```

---

## Error Codes

### HTTP Status Codes

| Code | Description | Usage |
|------|-------------|-------|
| 200 | OK | Successful request |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Invalid request parameters |
| 401 | Unauthorized | Authentication required or failed |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 500 | Internal Server Error | Server error |
| 503 | Service Unavailable | Service temporarily unavailable |

### Error Response Format

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional error details"
  }
}
```

### Common Error Codes

| Code | Description | Solution |
|------|-------------|----------|
| `INVALID_CREDENTIALS` | Username or password incorrect | Check credentials |
| `TOKEN_EXPIRED` | Access token expired | Use refresh token |
| `INVALID_TOKEN` | Invalid JWT token | Re-authenticate |
| `INSUFFICIENT_PERMISSIONS` | User lacks required permissions | Contact administrator |
| `RESOURCE_NOT_FOUND` | Requested resource doesn't exist | Verify resource ID |
| `RATE_LIMIT_EXCEEDED` | Too many requests | Wait and retry |
| `SERVICE_UNAVAILABLE` | Backend service unavailable | Check service status |

---

## Data Models

### User Object

```typescript
interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'user';
  created_at: string;  // ISO 8601 datetime
  updated_at?: string; // ISO 8601 datetime
  last_login?: string; // ISO 8601 datetime
}
```

### Container Object

```typescript
interface Container {
  id: string;
  name: string;
  image: string;
  state: 'running' | 'stopped' | 'paused';
  status: string;
  ports: string[];
  created: number;      // UNIX timestamp
}
```

### Storage Object

```typescript
interface Storage {
  path: string;
  device: string;
  type: string;
  total: number;        // bytes
  used: number;         // bytes
  available: number;    // bytes
  usage_percent: number;
  mount_options: string[];
}
```

### Service Object

```typescript
interface Service {
  name: string;
  description: string;
  state: 'running' | 'stopped' | 'failed';
  status: string;
}
```

### CPU Info Object

```typescript
interface CPUInfo {
  usage: number;        // 0.0-1.0
  cores: number;
  load1: number;
  load5: number;
  load15: number;
  frequency: number;    // MHz
}
```

### Memory Info Object

```typescript
interface MemoryInfo {
  total: number;        // bytes
  used: number;         // bytes
  free: number;         // bytes
  percent: number;      // 0-100
  swap_total: number;   // bytes
  swap_used: number;    // bytes
  swap_percent: number; // 0-100
}
```

### Disk Info Object

```typescript
interface DiskInfo {
  device: string;
  mount_point: string;
  file_system: string;
  total: number;        // bytes
  used: number;        // bytes
  free: number;        // bytes
  percent: number;     // 0-100
}
```

### Network Interface Object

```typescript
interface NetworkInterface {
  name: string;
  bytes_sent: number;
  bytes_recv: number;
  packets_sent: number;
  packets_recv: number;
  err_in: number;
  err_out: number;
  drop_in: number;
  drop_out: number;
}
```

---

## Rate Limiting

API requests are rate-limited to prevent abuse:

- **Standard endpoints**: 100 requests per minute
- **Streaming endpoints**: 10 requests per second
- **WebSocket**: 1 message per second

Rate limit headers are included in responses:

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067260
```

---

## Pagination

List endpoints support pagination:

```http
GET /api/users?page=1&limit=20
```

**Response:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

---

## CORS Configuration

Cross-Origin Resource Sharing (CORS) is configured to allow:

- **Origins**: Configured in backend `CORS_ORIGIN` environment variable
- **Methods**: GET, POST, PUT, DELETE, OPTIONS
- **Headers**: Authorization, Content-Type, X-Requested-With
- **Max Age**: 86400 seconds (24 hours)

---

## Best Practices

### Token Management

```typescript
// Store tokens securely
const token = localStorage.getItem('access_token');
const refreshToken = localStorage.getItem('refresh_token');

// Use access token for requests
axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;

// Handle token expiration
axios.interceptors.response.use(
  response => response,
  async error => {
    if (error.response?.status === 401) {
      // Refresh token
      const newToken = await refreshAccessToken();
      // Retry original request
      return axios.request(error.config);
    }
    return Promise.reject(error);
  }
);
```

### Error Handling

```typescript
try {
  const response = await axios.get('/api/monitor/cpu');
  console.log(response.data);
} catch (error) {
  if (error.response?.status === 401) {
    console.error('Authentication failed');
  } else if (error.response?.status === 500) {
    console.error('Server error');
  } else {
    console.error('Request failed:', error.message);
  }
}
```

### WebSocket Connection

```typescript
class MonitorClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;

  connect(token: string) {
    this.ws = new WebSocket(`ws://localhost:8888/api/monitor/ws?token=${token}`);
    
    this.ws.onopen = () => {
      console.log('WebSocket connected');
      this.reconnectAttempts = 0;
    };
    
    this.ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      this.handleMessage(data);
    };
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
    
    this.ws.onclose = () => {
      if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        setTimeout(() => this.connect(token), 5000);
      }
    };
  }
  
  private handleMessage(data: any) {
    switch (data.type) {
      case 'monitor_data':
        this.updateMonitorUI(data);
        break;
      default:
        console.log('Unknown message type:', data.type);
    }
  }
  
  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }
}
```

---

## Testing the API

### Using cURL

```bash
# Login
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# Get CPU info (replace <token>)
curl -X GET http://localhost:8888/api/monitor/cpu \
  -H "Authorization: Bearer <token>"

# List containers
curl -X GET http://localhost:8888/api/docker/containers \
  -H "Authorization: Bearer <token>"

# Stop container
curl -X POST http://localhost:8888/api/docker/containers/<id>/stop \
  -H "Authorization: Bearer <token>"
```

### Using Postman

1. **Import API collection** (if available)
2. **Set environment variables**:
   - `base_url`: `http://localhost:8888/api`
   - `access_token`: From login response
3. **Use Authorization header**: `Bearer {{access_token}}`

---

## Changelog

### Version 0.1.0 (Current)
- Initial API release
- Basic authentication
- System monitoring endpoints
- Docker management
- User management
- WebSocket real-time monitoring

---

**Last Updated**: 2026-06-12
**API Version**: 0.1.0
