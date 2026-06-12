# NAS Dashboard - API Documentation

Complete API reference for the NAS Dashboard backend services, including REST endpoints, WebSocket protocols, authentication, and data models.

## 📋 Table of Contents

1. [Getting Started](#getting-started)
2. [Authentication](#authentication)
3. [REST API Endpoints](#rest-api-endpoints)
4. [WebSocket Protocol](#websocket-protocol)
5. [Data Models](#data-models)
6. [Error Handling](#error-handling)
7. [Rate Limiting](#rate-limiting)
8. [CORS Configuration](#cors-configuration)
9. [Testing the API](#testing-the-api)

---

## Getting Started

### Base URLs

```
Development:  http://localhost:8888/api
Production:   https://api.yourdomain.com/api
WebSocket:    ws://localhost:8888/api/monitor/ws
```

### Authentication

All API endpoints (except login) require JWT authentication:

```http
Authorization: Bearer <access_token>
```

### Response Format

All API responses follow this format:

**Success Response:**
```json
{
  "success": true,
  "data": { /* response data */ }
}
```

**Error Response:**
```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

---

## Authentication

### Login

**Endpoint:** `POST /api/auth/login`

**Description:** Authenticate user and receive JWT tokens.

**Request Body:**
```json
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
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": "Invalid credentials",
  "code": "INVALID_CREDENTIALS"
}
```

### Refresh Token

**Endpoint:** `POST /api/auth/refresh`

**Description:** Refresh access token using refresh token.

**Request Body:**
```json
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

### Logout

**Endpoint:** `POST /api/auth/logout`

**Description:** Invalidate current tokens.

**Headers:**
```http
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### Get Current User

**Endpoint:** `GET /api/auth/me`

**Description:** Get current user information.

**Headers:**
```http
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "role": "admin",
    "created_at": "2024-01-01T00:00:00Z",
    "last_login": "2024-01-01T12:00:00Z"
  }
}
```

---

## REST API Endpoints

### System Monitoring

#### Get CPU Information

**Endpoint:** `GET /api/monitor/cpu`

**Description:** Get CPU usage and statistics.

**Headers:**
```http
Authorization: Bearer <access_token>
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "usage": 0.45,
    "cores": 24,
    "load1": 2.5,
    "load5": 2.1,
    "load15": 1.8,
    "frequency": 3200
  }
}
```

**Field Descriptions:**
- `usage`: CPU usage percentage (0.0-1.0)
- `cores`: Number of CPU cores
- `load1`: 1-minute load average
- `load5`: 5-minute load average
- `load15`: 15-minute load average
- `frequency`: CPU frequency in MHz

#### Get Memory Information

**Endpoint:** `GET /api/monitor/memory`

**Description:** Get memory usage and statistics.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total": 17179869184,
    "used": 8589934592,
    "free": 8589934592,
    "percent": 50.0,
    "swap_total": 4294967296,
    "swap_used": 0,
    "swap_percent": 0.0
  }
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

#### Get Disk Information

**Endpoint:** `GET /api/monitor/disk`

**Description:** Get disk usage information.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "disks": [
      {
        "device": "/dev/sda1",
        "mount_point": "/",
        "file_system": "ext4",
        "total": 1073741824000,
        "used": 536870912000,
        "free": 536870912000,
        "percent": 50.0
      }
    ]
  }
}
```

#### Get Network Information

**Endpoint:** `GET /api/monitor/network`

**Description:** Get network interface statistics.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
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
}
```

#### Get System Information

**Endpoint:** `GET /api/monitor/system`

**Description:** Get comprehensive system information.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
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
}
```

#### Get All Monitoring Data

**Endpoint:** `GET /api/monitor/all`

**Description:** Get all monitoring data in single request.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "cpu": { /* CPU data */ },
    "memory": { /* Memory data */ },
    "disk": { /* Disk data */ },
    "network": { /* Network data */ },
    "system": { /* System data */ }
  }
}
```

### Docker Management

#### List Containers

**Endpoint:** `GET /api/docker/containers`

**Description:** List all Docker containers.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
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
}
```

#### Start Container

**Endpoint:** `POST /api/docker/containers/{id}/start`

**Description:** Start a Docker container.

**Path Parameters:**
- `id`: Container ID or name

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Container started successfully"
}
```

#### Stop Container

**Endpoint:** `POST /api/docker/containers/{id}/stop`

**Description:** Stop a Docker container.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Container stopped successfully"
}
```

#### Restart Container

**Endpoint:** `POST /api/docker/containers/{id}/restart`

**Description:** Restart a Docker container.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Container restarted successfully"
}
```

#### Get Container Logs

**Endpoint:** `GET /api/docker/containers/{id}/logs`

**Description:** Get container logs.

**Query Parameters:**
- `tail`: Number of lines from end of logs (default: 100)
- `since`: UNIX timestamp to get logs since
- `timestamps`: Include timestamps (true/false)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "logs": "2024-01-01T00:00:00.000Z [INFO] Starting nginx...\n..."
  }
}
```

### Storage Management

#### Get Storage Information

**Endpoint:** `GET /api/storage`

**Description:** Get storage usage information.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
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
}
```

#### List Disks

**Endpoint:** `GET /api/storage/disks`

**Description:** List all system disks.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "disks": [
      {
        "name": "/dev/sda",
        "size": 1073741824000,
        "model": "Samsung SSD 860",
        "serial": "XYZ123",
        "partitions": [
          {
            "device": "/dev/sda1",
            "size": 536870912000,
            "type": "ext4",
            "mount_point": "/"
          }
        ]
      }
    ]
  }
}
```

#### Mount Disk

**Endpoint:** `POST /api/storage/mount`

**Description:** Mount a disk partition.

**Request Body:**
```json
{
  "device": "/dev/sdb1",
  "mount_point": "/mnt/data",
  "type": "ext4",
  "options": "rw"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Disk mounted successfully"
}
```

#### Unmount Disk

**Endpoint:** `POST /api/storage/unmount`

**Description:** Unmount a disk partition.

**Request Body:**
```json
{
  "mount_point": "/mnt/data"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Disk unmounted successfully"
}
```

#### Scan Storage

**Endpoint:** `POST /api/storage/scan`

**Description:** Scan and update storage information.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Storage scan completed"
}
```

### User Management

#### List Users

**Endpoint:** `GET /api/users`

**Description:** List all users.

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin",
        "created_at": "2024-01-01T00:00:00Z",
        "last_login": "2024-01-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

#### Create User

**Endpoint:** `POST /api/users`

**Description:** Create a new user.

**Request Body:**
```json
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
  "data": {
    "user": {
      "id": 2,
      "username": "newuser",
      "email": "newuser@example.com",
      "role": "user",
      "created_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

#### Update User

**Endpoint:** `PUT /api/users/{id}`

**Description:** Update user information.

**Request Body:**
```json
{
  "email": "updated@example.com",
  "role": "admin"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 2,
      "username": "newuser",
      "email": "updated@example.com",
      "role": "admin",
      "updated_at": "2024-01-01T13:00:00Z"
    }
  }
}
```

#### Delete User

**Endpoint:** `DELETE /api/users/{id}`

**Description:** Delete a user.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User deleted successfully"
}
```

#### Change Password

**Endpoint:** `PUT /api/users/{id}/password`

**Description:** Change user password.

**Request Body:**
```json
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

#### Get User SSH Keys

**Endpoint:** `GET /api/users/{id}/ssh-keys`

**Description:** Get user SSH public keys.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "keys": [
      {
        "id": 1,
        "name": "Work Laptop",
        "key": "ssh-rsa AAAAB3NzaC1yc2E...",
        "fingerprint": "SHA256:abc123...",
        "added_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

#### Add SSH Key

**Endpoint:** `POST /api/users/{id}/ssh-keys`

**Description:** Add SSH public key for user.

**Request Body:**
```json
{
  "name": "Home Desktop",
  "key": "ssh-rsa AAAAB3NzaC1yc2E..."
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "key": {
      "id": 2,
      "name": "Home Desktop",
      "fingerprint": "SHA256:def456...",
      "added_at": "2024-01-01T12:00:00Z"
    }
  }
}
```

#### Delete SSH Key

**Endpoint:** `DELETE /api/users/{id}/ssh-keys/{key_id}`

**Description:** Delete SSH key.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "SSH key deleted successfully"
}
```

### Service Management

#### List Services

**Endpoint:** `GET /api/services`

**Description:** List all system services.

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "services": [
      {
        "name": "nginx",
        "description": "A high performance web server",
        "state": "running",
        "status": "active (running)"
      }
    ]
  }
}
```

#### Start Service

**Endpoint:** `POST /api/services/{name}/start`

**Description:** Start a system service.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service started successfully"
}
```

#### Stop Service

**Endpoint:** `POST /api/services/{name}/stop`

**Description:** Stop a system service.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service stopped successfully"
}
```

#### Restart Service

**Endpoint:** `POST /api/services/{name}/restart`

**Description:** Restart a system service.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Service restarted successfully"
}
```

### Health Check

#### Health Status

**Endpoint:** `GET /api/health`

**Description:** Check API health status.

**Response (200 OK):**
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

## Error Handling

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
  "success": false,
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
| `VALIDATION_ERROR` | Request validation failed | Check request parameters |
| `DUPLICATE_RESOURCE` | Resource already exists | Use different identifier |
| `OPERATION_FAILED` | Operation execution failed | Check system logs |

---

## Rate Limiting

API requests are rate-limited to prevent abuse:

- **Standard endpoints**: 100 requests per minute
- **Streaming endpoints**: 10 requests per second
- **WebSocket**: 1 message per second

### Rate Limit Headers

```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1704067260
```

### Handling Rate Limits

```javascript
try {
  const response = await axios.get('/api/users');
} catch (error) {
  if (error.response.status === 429) {
    // Rate limit exceeded
    const retryAfter = error.response.headers['retry-after'];
    console.log(`Retry after ${retryAfter} seconds`);
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

### CORS Headers

```http
Access-Control-Allow-Origin: https://yourdomain.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Authorization, Content-Type, X-Requested-With
Access-Control-Max-Age: 86400
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

### Using JavaScript

```javascript
import axios from 'axios';

const API_URL = 'http://localhost:8888/api';

// Login
const login = async (username, password) => {
  const response = await axios.post(`${API_URL}/auth/login`, {
    username,
    password
  });
  return response.data.token;
};

// Get CPU data
const getCPUData = async (token) => {
  const response = await axios.get(`${API_URL}/monitor/cpu`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.data;
};

// List containers
const listContainers = async (token) => {
  const response = await axios.get(`${API_URL}/docker/containers`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.data;
};
```

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

## API Versioning

The API follows semantic versioning:

- **Current Version**: 0.1.0
- **Version Format**: Major.Minor.Patch
- **Breaking Changes**: Increment major version
- **New Features**: Increment minor version
- **Bug Fixes**: Increment patch version

### Version in Response

All API responses include the version:

```json
{
  "success": true,
  "data": { /* ... */ },
  "meta": {
    "version": "0.1.0",
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

---

## Changelog

### Version 0.1.0 (Current)
- Initial API release
- Basic authentication
- System monitoring endpoints
- Docker management
- User management
- WebSocket real-time monitoring
- Storage management
- Service management

---

**Last Updated**: 2026-06-12  
**API Version**: 0.1.0  
**Status**: Active Development
