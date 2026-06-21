# NAS Dashboard 用户管理和认证系统 - 快速使用指南

## 🚀 快速开始

### 1. 登录系统

```bash
# 默认管理员账号
用户名: admin
密码: admin

# 登录API
curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'
```

### 2. 使用Token访问API

```bash
# 从登录响应中获取token，然后：
export TOKEN="your_jwt_token_here"

# 访问受保护的API
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers
```

---

## 🐳 Docker管理功能

### 查看所有容器
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers
```

### 启动/停止容器
```bash
# 停止容器
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers/{container_id}/stop

# 启动容器
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers/{container_id}/start
```

### 查看容器日志
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers/{container_id}/logs
```

---

## 👥 用户管理功能

### 获取系统用户列表
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/users
```

### 获取Docker容器用户
```bash
# Immich用户
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/immich/users

# 系统用户
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/users
```

---

## 💾 系统监控功能

### CPU监控
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/cpu
```

### 内存监控
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/memory
```

### 磁盘监控
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/disk
```

### 网络监控
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/network
```

---

## 🔧 系统服务管理

### 查看所有服务
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/services
```

### 控制服务
```bash
# 启动服务
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/services/{service_name}/start

# 停止服务
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/services/{service_name}/stop

# 重启服务
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/services/{service_name}/restart
```

---

## 🌐 网络管理

### 查看网络接口
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/network/interfaces
```

### 查看无线网络
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/network/interfaces/wifi
```

### 扫描WiFi网络
```bash
curl -X POST -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/network/wifi/scan
```

---

## 📁 存储管理

### 查看磁盘信息
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/storage/disks
```

### 查看存储池
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/storage/pools
```

### 查看SMB共享
```bash
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/storage/smb
```

---

## 🔐 权限和角色

### 角色类型
- **admin** - 完全管理权限
- **user** - 基本使用权限
- **guest** - 只读权限

### 权限矩阵
| 功能 | admin | user | guest |
|------|-------|------|-------|
| Docker管理 | ✅ | ❌ | ❌ |
| 用户管理 | ✅ | ❌ | ❌ |
| 系统监控 | ✅ | ✅ | ✅ |
| 服务管理 | ✅ | ❌ | ❌ |
| 网络管理 | ✅ | ⚠️ | ❌ |
| 存储管理 | ✅ | ⚠️ | ❌ |

---

## 🔄 Token刷新

```bash
# 使用Refresh Token获取新的Access Token
curl -X POST http://localhost:8888/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refreshToken":"your_refresh_token"}'
```

---

## 🚨 常见错误处理

### 401 Unauthorized
- **原因**: Token无效或过期
- **解决**: 重新登录或刷新Token

### 403 Forbidden
- **原因**: 权限不足
- **解决**: 使用具有相应权限的账号

### 404 Not Found
- **原因**: 资源不存在
- **解决**: 检查请求路径和参数

---

## 📊 当前运行的Docker服务

系统当前运行以下12个Docker服务：

1. **authentik-server** - 认证服务 (端口9000, 9443)
2. **immich-server** - 照片管理 (端口2283)
3. **samba** - 文件共享 (端口139, 445)
4. **alist** - 存储聚合 (端口5244, 5245)
5. **private-git** - 代码仓库 (端口3000, 2222)
6. **avahi** - 网络发现
7. **shairport-sync** - AirPlay同步
8. **其他辅助服务** - postgres, redis等

---

## 💡 使用技巧

### 1. 批量操作容器
```bash
# 获取所有容器ID
containers=$(curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers | \
  jq -r '.containers[].id')

# 批量重启
for id in $containers; do
  curl -X POST -H "Authorization: Bearer $TOKEN" \
    http://localhost:8888/api/docker/containers/$id/restart
done
```

### 2. 监控特定容器
```bash
# 持续监控容器状态
watch -n 5 'curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers/{id}/stats | jq'
```

### 3. 系统状态概览
```bash
# 一键获取系统概览
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/cpu && \
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/memory && \
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers
```

---

## 🎯 最佳实践

### 安全建议
1. **定期修改密码** - 特别是admin账号
2. **使用强密码** - 包含大小写字母、数字、特殊字符
3. **定期刷新Token** - 避免Token过期
4. **限制API访问** - 配置防火墙规则

### 性能优化
1. **批量操作** - 合并多个API请求
2. **缓存数据** - 避免频繁查询
3. **监控日志** - 定期检查系统日志
4. **资源管理** - 及时清理无用容器

### 故障排除
1. **检查服务状态** - 确保后端服务运行
2. **验证Token** - 确认Token未过期
3. **查看日志** - 检查错误信息
4. **测试连接** - 验证网络连通性

---

## 🚀 快速命令参考

```bash
# 登录并设置Token
export TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}' | \
  jq -r '.token')

# 查看所有Docker容器
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/docker/containers | jq

# 查看系统CPU使用率
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/monitor/cpu | jq

# 查看所有用户
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8888/api/users | jq

# 刷新Token
curl -X POST http://localhost:8888/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refreshToken\":\"$REFRESH_TOKEN\"}" | jq
```

---

**🎉 现在你可以开始使用NAS Dashboard的完整功能了！**