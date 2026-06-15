# 系统监控模块

## 概述

系统监控模块提供全面的系统监控功能，包括CPU、内存、磁盘、网络、进程和服务的实时监控和管理。

## 功能特性

### 基础监控
- CPU使用率监控
- 内存使用监控
- 磁盘IO监控
- 网络流量监控
- 系统温度监控

### 进程管理
- 进程列表查看
- 进程详情查看
- 进程终止
- 进程优先级调整

### 服务管理
- 服务状态查看
- 服务启动/停止/重启
- 服务启用/禁用
- 服务日志查看

### 系统信息
- 硬件信息查询
- 系统运行时间
- 系统事件日志
- 系统警报管理

## 后端API实现

### 主要API文件
- `backend/internal/api/monitor.go` - 监控API
- `backend/internal/api/system.go` - 系统信息API
- `backend/service/monitor/` - 监控服务

### 前端组件
- `frontend/src/apps/SystemMonitor.vue` - 系统监控界面
- `frontend/src/apps/MonitorConsole.vue` - 监控控制台
- `frontend/src/components/Chart.vue` - 图表组件

## API端点

### 基础监控
- `GET /api/monitor/cpu` - 获取CPU使用率
- `GET /api/monitor/memory` - 获取内存使用情况
- `GET /api/monitor/disk` - 获取磁盘使用情况
- `GET /api/monitor/network` - 获取网络流量
- `GET /api/monitor/temperature` - 获取系统温度

### 进程管理
- `GET /api/monitor/processes` - 获取进程列表
- `GET /api/monitor/processes/:pid` - 获取进程详情
- `DELETE /api/monitor/processes/:pid` - 终止进程

### 服务管理
- `GET /api/monitor/services` - 获取服务列表
- `GET /api/monitor/services/:name` - 获取服务详情
- `POST /api/monitor/services/:name/start` - 启动服务
- `POST /api/monitor/services/:name/stop` - 停止服务
- `POST /api/monitor/services/:name/restart` - 重启服务

### 系统信息
- `GET /api/system/info` - 获取系统信息
- `GET /api/system/hardware` - 获取硬件信息
- `GET /api/system/uptime` - 获取系统运行时间
- `GET /api/system/power` - 获取电源使用情况

### 事件和日志
- `GET /api/monitor/events` - 获取系统事件
- `GET /api/monitor/logs` - 获取系统日志
- `POST /api/monitor/logs/clear` - 清除日志

### 警报管理
- `GET /api/monitor/alerts` - 获取警报列表
- `POST /api/monitor/alerts` - 创建警报
- `PUT /api/monitor/alerts/:id` - 更新警报
- `DELETE /api/monitor/alerts/:id` - 删除警报

## WebSocket监控

### 实时数据推送
- `WS /ws/monitor` - 实时监控数据推送

WebSocket连接后会收到以下类型的实时数据：
```json
{
  "type": "cpu",
  "data": {
    "usage": 45.2,
    "cores": [30.1, 45.3, 22.1, 67.8]
  }
}
```

## 使用示例

### 获取CPU使用率
```bash
curl -X GET http://localhost:8888/api/monitor/cpu \
  -H "Authorization: Bearer <token>"
```

### 终止进程
```bash
curl -X DELETE http://localhost:8888/api/monitor/processes/1234 \
  -H "Authorization: Bearer <token>"
```

### 重启服务
```bash
curl -X POST http://localhost:8888/api/monitor/services/nginx/restart \
  -H "Authorization: Bearer <token>"
```

### 创建警报
```bash
curl -X POST http://localhost:8888/api/monitor/alerts \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "High CPU Alert",
    "type": "cpu",
    "condition": "usage > 80",
    "action": "email",
    "enabled": true
  }'
```

## 配置文件

### 监控配置
- `config/monitor.conf` - 监控主配置
- `config/monitors/cpu.conf` - CPU监控配置
- `config/monitors/memory.conf` - 内存监控配置
- `config/monitors/disk.conf` - 磁盘监控配置
- `config/monitors/network.conf` - 网络监控配置

### 警报配置
- `config/alerts.conf` - 警报规则配置
- `config/alerts/notifications.conf` - 通知配置

## 前端展示

### 监控面板
```vue
<template>
  <div class="monitor-panel">
    <CPUChart :data="cpuData" />
    <MemoryChart :data="memoryData" />
    <DiskChart :data="diskData" />
    <NetworkChart :data="networkData" />
  </div>
</template>
```

### WebSocket连接
```javascript
const ws = new WebSocket('ws://localhost:8888/ws/monitor');
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  updateMonitorData(data);
};
```

## 性能优化

- 使用WebSocket替代轮询
- 客户端数据缓存
- 图表数据聚合
- 历史数据压缩存储

## 监控指标

### CPU指标
- 总体使用率
- 各核心使用率
- 上下文切换
- 运行队列

### 内存指标
- 总内存使用
- 应用内存
- 缓存和缓冲
- Swap使用

### 磁盘指标
- 读写速率
- IO等待时间
- 磁盘使用率
- IO队列深度

### 网络指标
- 流入/流出流量
- 连接数
- 网络错误
- 接口状态

## 故障排除

### 监控数据不准确
1. 检查监控服务状态：`systemctl status nas-monitor`
2. 查看监控日志：`journalctl -u nas-monitor`
3. 验证数据源：`cat /proc/stat`

### WebSocket连接失败
1. 检查防火墙设置
2. 验证SSL证书
3. 检查浏览器兼容性

### 进程无法终止
1. 检查进程权限
2. 尝试SIGKILL：`kill -9 <pid>`
3. 查看进程状态：`cat /proc/<pid>/status`

## 相关文档

- [系统监控API文档](../../docs/API_DOCUMENTATION.md)
- [用户指南](../../docs/USER_GUIDE.md)
- [故障排除指南](../../docs/TROUBLESHOOTING.md)