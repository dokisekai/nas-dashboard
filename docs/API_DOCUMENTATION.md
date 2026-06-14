# NAS Dashboard API 文档

## 概述

NAS Dashboard 提供完整的 REST API 用于管理和监控系统。所有 API 端点都使用 JSON 格式进行数据交换。

### 基础信息

- **Base URL**: `http://your-nas-ip:8888/api`
- **认证方式**: JWT Bearer Token
- **数据格式**: JSON
- **字符编码**: UTF-8

### 认证

所有 API 请求都需要在 HTTP Header 中包含有效的 JWT Token：

```http
Authorization: Bearer <your-jwt-token>
```

#### 获取 Token

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password"
}
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expiresIn": 86400
  }
}
```

## 存储池管理 API

### 获取存储池列表

```http
GET /api/storage/pools
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "name": "main-pool",
      "type": "mergerfs",
      "status": "active",
      "mountPoint": "/mnt/main-pool",
      "totalSize": 1099511627776,
      "usedSize": 549755813888,
      "freeSize": 549755813888,
      "description": "Main storage pool",
      "createdAt": "2026-06-12T10:00:00Z",
      "updatedAt": "2026-06-12T10:00:00Z"
    }
  ]
}
```

### 获取存储池详情

```http
GET /api/storage/pools/:name
```

参数：
- `name`: 存储池名称

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "main-pool",
    "type": "mergerfs",
    "status": "active",
    "mountPoint": "/mnt/main-pool",
    "totalSize": 1099511627776,
    "usedSize": 549755813888,
    "freeSize": 549755813888,
    "description": "Main storage pool",
    "config": {
      "categories": {
        "RW": "min space most",
        "RO": "most free space"
      }
    },
    "disks": [
      {
        "device": "/dev/sdb",
        "size": 549755813888,
        "status": "active",
        "branchPath": "/mnt/disk1"
      }
    ],
    "snapshots": []
  }
}
```

### 创建存储池

```http
POST /api/storage/pools
Content-Type: application/json

{
  "name": "backup-pool",
  "type": "mergerfs",
  "mountPoint": "/mnt/backup-pool",
  "disks": [
    {
      "device": "/dev/sdc",
      "branchPath": "/mnt/disk2"
    }
  ],
  "config": {
    "categories": {
      "RW": "mspmfs",
      "RO": "most free space"
    }
  },
  "description": "Backup storage pool"
}
```

响应：

```json
{
  "code": 201,
  "message": "Storage pool created successfully",
  "data": {
    "id": 2,
    "name": "backup-pool"
  }
}
```

### 更新存储池

```http
PUT /api/storage/pools/:name
Content-Type: application/json

{
  "description": "Updated description",
  "config": {
    "categories": {
      "RW": "most free space",
      "RO": "most free space"
    }
  }
}
```

### 删除存储池

```http
DELETE /api/storage/pools/:name
```

### 添加磁盘到存储池

```http
POST /api/storage/pools/:name/disks
Content-Type: application/json

{
  "device": "/dev/sdd",
  "branchPath": "/mnt/disk3",
  "priority": 0
}
```

### 从存储池移除磁盘

```http
DELETE /api/storage/pools/:name/disks/:device
```

### 挂载存储池

```http
POST /api/storage/pools/:name/mount
```

### 卸载存储池

```http
POST /api/storage/pools/:name/umount
```

### 获取存储池状态

```http
GET /api/storage/pools/:name/status
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "active",
    "mounted": true,
    "totalSize": 1099511627776,
    "usedSize": 549755813888,
    "freeSize": 549755813888,
    "usagePercent": 50.0,
    "disks": [
      {
        "device": "/dev/sdb",
        "status": "active",
        "size": 549755813888,
        "used": 274877906944
      }
    ]
  }
}
```

## 系统监控 API

### 获取系统概览

```http
GET /api/monitor/overview
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "cpu": {
      "usagePercent": 25.5,
      "cores": 8,
      "frequency": 2400
    },
    "memory": {
      "total": 8589934592,
      "used": 4294967296,
      "free": 4294967296,
      "usagePercent": 50.0
    },
    "disk": [
      {
        "device": "/dev/sda1",
        "mountPoint": "/",
        "total": 1099511627776,
        "used": 549755813888,
        "free": 549755813888,
        "usagePercent": 50.0
      }
    ],
    "network": [
      {
        "interface": "eth0",
        "bytesSent": 104857600,
        "bytesRecv": 209715200,
        "bytesSentSpeed": 102400,
        "bytesRecvSpeed": 204800
      }
    ]
  }
}
```

### 获取进程列表

```http
GET /api/monitor/processes?limit=50&sort=cpu
```

参数：
- `limit`: 返回数量限制（默认50）
- `sort`: 排序字段（cpu, memory, name）

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "pid": 1234,
      "name": "nginx",
      "status": "running",
      "cpuPercent": 5.2,
      "memoryPercent": 2.1,
      "memory": 177209600,
      "threads": 4,
      "username": "www-data",
      "command": "nginx: worker process"
    }
  ]
}
```

### 获取进程详情

```http
GET /api/monitor/processes/:pid
```

### 终止进程

```http
DELETE /api/monitor/processes/:pid
```

### 获取系统服务

```http
GET /api/monitor/services
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "name": "nginx",
      "description": "A high performance web server",
      "status": "running",
      "enabled": true,
      "loadState": "loaded",
      "activeState": "active",
      "mainPid": 1234
    }
  ]
}
```

### 启动服务

```http
POST /api/monitor/services/:name/start
```

### 停止服务

```http
POST /api/monitor/services/:name/stop
```

### 重启服务

```http
POST /api/monitor/services/:name/restart
```

### 获取温度信息

```http
GET /api/monitor/temperature
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "sensors": [
      {
        "name": "CPU Core 0",
        "current": 45.0,
        "max": 80.0,
        "critical": 90.0,
        "unit": "°C"
      }
    ]
  }
}
```

## 磁盘管理 API

### 获取磁盘列表

```http
GET /api/storage/disks
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "device": "/dev/sdb",
      "model": "Samsung SSD 870",
      "size": 1099511627776,
      "health": "good",
      "temperature": 35,
      "partitions": [
        {
          "device": "/dev/sdb1",
          "size": 549755813888,
          "type": "primary",
          "filesystem": "ext4",
          "mountPoint": "/mnt/disk1"
        }
      ]
    }
  ]
}
```

### 获取磁盘详情

```http
GET /api/storage/disks/:device
```

### 获取磁盘分区

```http
GET /api/storage/disks/:device/partitions
```

### 创建分区

```http
POST /api/storage/disks/:device/partition
Content-Type: application/json

{
  "start": 1048576,
  "end": 549755813888,
  "type": "primary",
  "filesystem": "ext4"
}
```

### 删除分区

```http
DELETE /api/storage/disks/:device/partitions/:id
```

### 获取磁盘 SMART 信息

```http
GET /api/storage/disks/:device/smart
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "device": "/dev/sdb",
    "health": "good",
    "attributes": [
      {
        "id": 1,
        "name": "Raw Read Error Rate",
        "value": 100,
        "worst": 100,
        "threshold": 50,
        "status": "good"
      }
    ]
  }
}
```

### 运行磁盘性能测试

```http
POST /api/storage/disks/:device/benchmark
Content-Type: application/json

{
  "size": 4096,
  "queueDepth": 32,
  "duration": 60,
  "mode": "mixed"
}
```

## RAID 管理 API

### 获取 RAID 列表

```http
GET /api/storage/raid
```

### 创建 RAID

```http
POST /api/storage/raid
Content-Type: application/json

{
  "name": "raid1",
  "level": "1",
  "devices": ["/dev/sdb", "/dev/sdc"]
}
```

### 获取 RAID 详情

```http
GET /api/storage/raid/:name
```

### 删除 RAID

```http
DELETE /api/storage/raid/:name
```

### 添加磁盘到 RAID

```http
POST /api/storage/raid/:name/add
Content-Type: application/json

{
  "device": "/dev/sdd"
}
```

## LVM 管理 API

### 获取物理卷列表

```http
GET /api/storage/lvm/pv
```

### 创建物理卷

```http
POST /api/storage/lvm/pv
Content-Type: application/json

{
  "device": "/dev/sdb"
}
```

### 获取卷组列表

```http
GET /api/storage/lvm/vg
```

### 创建卷组

```http
POST /api/storage/lvm/vg
Content-Type: application/json

{
  "name": "vg1",
  "devices": ["/dev/sdb", "/dev/sdc"]
}
```

### 获取逻辑卷列表

```http
GET /api/storage/lvm/lv
```

### 创建逻辑卷

```http
POST /api/storage/lvm/lv
Content-Type: application/json

{
  "name": "lv1",
  "vgName": "vg1",
  "size": 107374182400
}
```

## 配额管理 API

### 获取用户配额

```http
GET /api/users/:username/quota
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "user": {
        "id": 1,
        "username": "john",
        "email": "john@example.com"
      },
      "path": "/home/john",
      "softLimit": 10737418240,
      "hardLimit": 16106127360,
      "usedSpace": 5368709120,
      "gracePeriod": 7
    }
  ]
}
```

### 设置用户配额

```http
PUT /api/users/:username/quota
Content-Type: application/json

{
  "path": "/home/john",
  "softLimit": 10737418240,
  "hardLimit": 16106127360,
  "gracePeriod": 7
}
```

### 获取组配额

```http
GET /api/groups/:name/quota
```

### 设置组配额

```http
PUT /api/groups/:name/quota
Content-Type: application/json

{
  "path": "/shared/group",
  "softLimit": 53687091200,
  "hardLimit": 64424509440,
  "gracePeriod": 7
}
```

### 生成配额报告

```http
GET /api/storage/quota/report
```

响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "totalUsers": 10,
    "activeUsers": 8,
    "totalQuotaUsed": 536870912000,
    "totalQuotaTotal": 1073741824000,
    "topUsers": [
      {
        "username": "john",
        "usedSpace": 107374182400,
        "hardLimit": 214748364800,
        "usagePercent": 50.0
      }
    ]
  }
}
```

## WebSocket API

### 连接 WebSocket

```javascript
const ws = new WebSocket('ws://your-nas-ip:8888/ws?token=<your-jwt-token>');

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};
```

### 消息格式

#### 系统状态更新

```json
{
  "type": "system_status",
  "data": {
    "cpu": 25.5,
    "memory": 50.0,
    "disk": [
      {
        "device": "/dev/sda1",
        "usagePercent": 60.0
      }
    ]
  }
}
```

#### 进程更新

```json
{
  "type": "process_update",
  "data": {
    "pid": 1234,
    "cpuPercent": 10.5,
    "memoryPercent": 5.2
  }
}
```

#### 告警通知

```json
{
  "type": "alert",
  "data": {
    "id": 1,
    "type": "user_quota",
    "severity": "warning",
    "message": "User john has exceeded 80% of quota",
    "timestamp": "2026-06-12T10:30:00Z"
  }
}
```

## 错误处理

### 错误响应格式

```json
{
  "code": 400,
  "message": "Invalid request parameters",
  "errors": [
    {
      "field": "name",
      "message": "Name is required"
    }
  ]
}
```

### 常见错误代码

- `200`: 成功
- `201`: 创建成功
- `400`: 请求参数错误
- `401`: 未授权
- `403`: 禁止访问
- `404`: 资源不存在
- `500`: 服务器内部错误

## 速率限制

- **未认证用户**: 100 requests/minute
- **已认证用户**: 1000 requests/minute
- **WebSocket 连接**: 5 connections/user

超出限制会返回 `429 Too Many Requests` 错误。

## SDK 和客户端库

### Python SDK

```python
from nas_dashboard import Client

client = Client('http://your-nas-ip:8888', token='your-jwt-token')

# 获取存储池列表
pools = client.storage_pools.list()

# 创建存储池
pool = client.storage_pools.create(
    name='new-pool',
    type='mergerfs',
    mount_point='/mnt/new-pool'
)
```

### JavaScript SDK

```javascript
import { NASDashboard } from '@nas-dashboard/sdk';

const client = new NASDashboard({
  baseURL: 'http://your-nas-ip:8888',
  token: 'your-jwt-token'
});

// 获取存储池列表
const pools = await client.storagePools.list();

// 创建存储池
const pool = await client.storagePools.create({
  name: 'new-pool',
  type: 'mergerfs',
  mountPoint: '/mnt/new-pool'
});
```

## 支持

如有问题，请联系技术支持或查看开发者文档。
