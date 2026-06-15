# 网络管理模块

## 概述

网络管理模块提供完整的网络配置和管理功能，支持多种网络接口类型、连接方式和高级网络功能。

## 功能特性

### 网络接口管理
- 以太网接口配置
- Wi-Fi接口管理
- 虚拟接口配置
- 接口启用/禁用
- 接口重启

### 网络连接
- PPPoE拨号连接
- Wi-Fi连接管理
- 网络扫描和发现
- 连接状态监控

### DNS和代理
- DNS服务器配置
- 代理服务器设置
- PAC自动配置
- 网络诊断工具

### 防火墙管理
- 防火墙规则配置
- 端口转发
- 访问控制列表
- 连接跟踪

## 后端API实现

### 主要API文件
- `backend/internal/api/network.go` - 网络管理主API
- `backend/pkg/system/network.go` - 网络系统操作

### 前端组件
- `frontend/src/apps/NetworkManager.vue` - 网络管理界面
- `frontend/src/components/ControlPanel/NetworkSettingsPanel.vue` - 网络设置面板
- `frontend/src/components/ControlPanel/NetworkInterfacesComponent.vue` - 网络接口组件

## API端点

### 接口管理
- `GET /api/network/interfaces` - 获取所有网络接口
- `GET /api/network/interfaces/ethernet` - 获取以太网接口
- `GET /api/network/interfaces/wifi` - 获取Wi-Fi接口
- `GET /api/network/interface/:interface/config` - 获取接口配置
- `PUT /api/network/interface/:interface/config` - 设置接口配置
- `POST /api/network/interface/:interface/restart` - 重启接口
- `POST /api/network/interface/:interface/:action` - 接口操作(up/down)

### PPPoE管理
- `GET /api/network/interface/:interface/pppoe` - 获取PPPoE配置
- `POST /api/network/interface/:interface/pppoe` - 配置PPPoE

### Wi-Fi管理
- `GET /api/network/wifi/scan` - 扫描Wi-Fi网络
- `POST /api/network/wifi/connect` - 连接Wi-Fi
- `POST /api/network/wifi/disconnect` - 断开Wi-Fi
- `GET /api/network/wifi/current` - 获取当前Wi-Fi连接

### DNS和代理
- `GET /api/network/dns` - 获取DNS配置
- `POST /api/network/dns` - 设置DNS
- `GET /api/network/proxy` - 获取代理配置
- `POST /api/network/proxy` - 设置代理

### 防火墙
- `GET /api/security/firewall/rules` - 获取防火墙规则
- `POST /api/security/firewall/rules` - 创建防火墙规则
- `PUT /api/security/firewall/rules/:id` - 更新防火墙规则
- `DELETE /api/security/firewall/rules/:id` - 删除防火墙规则
- `POST /api/security/firewall/apply` - 应用防火墙规则

## 配置文件

### 网络配置
- `config/network/interfaces.conf` - 网络接口配置
- `config/network/wifi.conf` - Wi-Fi配置
- `config/network/pppoe.conf` - PPPoE配置
- `config/network/dns.conf` - DNS配置

### 防火墙配置
- `config/firewall/iptables.rules` - iptables规则
- `config/firewall/nftables.conf` - nftables配置

## 使用示例

### 配置静态IP
```bash
curl -X PUT http://localhost:8888/api/network/interface/eth0/config \
  -H "Authorization: Bearer <token>" \
  -d '{
    "method": "static",
    "address": "192.168.1.100",
    "netmask": "255.255.255.0",
    "gateway": "192.168.1.1",
    "dns": ["8.8.8.8", "8.8.4.4"]
  }'
```

### 配置PPPoE连接
```bash
curl -X POST http://localhost:8888/api/network/interface/eth0/pppoe \
  -H "Authorization: Bearer <token>" \
  -d '{
    "username": "user@isp",
    "password": "password",
    "mtu": 1492
  }'
```

### 连接Wi-Fi
```bash
curl -X POST http://localhost:8888/api/network/wifi/connect \
  -H "Authorization: Bearer <token>" \
  -d '{
    "interface": "wlan0",
    "ssid": "WiFi-Name",
    "password": "wifi-password",
    "hidden": false
  }'
```

### 设置防火墙规则
```bash
curl -X POST http://localhost:8888/api/security/firewall/rules \
  -H "Authorization: Bearer <token>" \
  -d '{
    "action": "accept",
    "direction": "input",
    "protocol": "tcp",
    "destination_port": 80,
    "enabled": true
  }'
```

## 网络诊断

### 常用命令
```bash
# 查看网络接口状态
ip addr show

# 查看路由表
ip route show

# 测试网络连接
ping -c 4 google.com

# 查看端口监听
netstat -tulpn

# 追踪网络路由
traceroute google.com

# DNS查询
nslookup google.com
```

## 故障排除

### 网络接口无法启动
1. 检查接口状态：`ip link show`
2. 查看系统日志：`dmesg | grep -i network`
3. 检查网络配置：`cat /etc/network/interfaces`

### Wi-Fi连接失败
1. 扫描可用的网络：`iwlist scan`
2. 检查无线驱动：`lspci | grep -i wireless`
3. 查看连接日志：`journalctl -u NetworkManager`

### DNS解析失败
1. 测试DNS服务器：`nslookup google.com 8.8.8.8`
2. 检查DNS配置：`cat /etc/resolv.conf`
3. 清除DNS缓存：`systemd-resolve --flush-caches`

## 性能优化

- 使用千兆以太网接口
- 配置适当的MTU大小
- 启用网络接口卸载功能
- 优化防火墙规则顺序

## 安全建议

- 启用防火墙保护
- 使用强密码保护Wi-Fi
- 定期更新网络密码
- 限制远程访问来源
- 监控网络异常流量

## 相关文档

- [网络管理API文档](../../docs/API_DOCUMENTATION.md)
- [系统部署指南](../../docs/DEPLOYMENT_GUIDE.md)
- [DSM网络管理分析](../../docs/DSM_CONTROL_PANEL_ANALYSIS.md)