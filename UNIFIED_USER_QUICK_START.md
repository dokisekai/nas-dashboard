# 统一用户管理系统 - 快速使用示例

## 🚀 一键启动

```bash
# 1. 运行集成脚本
./integrate_unified_users.sh

# 2. 配置环境变量（重要！）
nano .env.unified-users

# 3. 启动服务
cd backend && ./nas-dashboard

# 4. 访问界面
# 前端：http://localhost:5173/unified-users
# 后端API：http://localhost:8888/api/unified-users
```

## 📱 快速使用示例

### 1. 创建用户（自动同步到所有服务）

```bash
# API方式
curl -X POST http://localhost:8888/api/unified-users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "email": "alice@example.com",
    "name": "Alice Johnson",
    "password": "secure_password_123",
    "role": "user",
    "groups": ["docker", "media", "storage"],
    "isActive": true
  }'
```

**响应示例：**
```json
{
  "username": "alice",
  "status": "success",
  "details": {
    "system": {
      "serviceName": "system",
      "status": "success",
      "message": "User created in system"
    },
    "immich": {
      "serviceName": "immich",
      "status": "success",
      "message": "User created in Immich"
    },
    "nextcloud": {
      "serviceName": "nextcloud",
      "status": "success",
      "message": "User created in Nextcloud"
    },
    "docker": {
      "serviceName": "docker",
      "status": "success",
      "message": "User created in Docker containers"
    }
  },
  "syncTime": "2024-01-15T10:30:00Z"
}
```

### 2. 修改用户（所有服务自动更新）

```bash
# 修改密码（所有服务同步）
curl -X PUT http://localhost:8888/api/unified-users/alice \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "new_secure_password_456",
    "email": "alice.new@example.com",
    "role": "admin"
  }'
```

### 3. 批量用户同步

```bash
# 同步所有用户到所有服务
curl -X POST http://localhost:8888/api/unified-users/sync \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例：**
```json
{
  "totalUsers": 15,
  "syncedUsers": 14,
  "failedUsers": 1,
  "userResults": {
    "alice": {
      "username": "alice",
      "status": "success",
      "details": { /* ... */ }
    },
    "bob": {
      "username": "bob",
      "status": "partial",
      "details": {
        "system": { "status": "success" },
        "immich": { "status": "success" },
        "nextcloud": { 
          "status": "failed",
          "error": "Connection timeout" 
        }
      }
    }
  },
  "syncStartTime": "2024-01-15T10:30:00Z",
  "syncEndTime": "2024-01-15T10:31:30Z",
  "duration": "90s"
}
```

### 4. 删除用户（从所有服务删除）

```bash
curl -X DELETE http://localhost:8888/api/unified-users/alice \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 5. 查看服务状态

```bash
curl http://localhost:8888/api/unified-users/status \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**响应示例：**
```json
{
  "services": [
    {
      "name": "system",
      "status": "healthy",
      "userCount": 15
    },
    {
      "name": "immich",
      "status": "healthy",
      "userCount": 12
    },
    {
      "name": "nextcloud",
      "status": "healthy",
      "userCount": 14
    },
    {
      "name": "jellyfin",
      "status": "error",
      "userCount": 0,
      "error": "Service not available"
    }
  ],
  "lastSyncTime": "2024-01-15T10:00:00Z",
  "autoSyncEnabled": true
}
```

## 🎯 实际使用场景

### 场景1：新员工入职

```bash
# 一次性创建用户，自动同步到所有系统
curl -X POST http://localhost:8888/api/unified-users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "new_employee",
    "email": "employee@company.com",
    "name": "New Employee",
    "password": "welcome_2024",
    "role": "user",
    "groups": ["docker", "media", "storage", "network"],
    "isActive": true
  }'

# 结果：
# ✅ Linux系统用户创建完成
# ✅ Immich账户创建完成
# ✅ Nextcloud账户创建完成
# ✅ Jellyfin账户创建完成
# ✅ Docker容器权限配置完成
```

### 场景2：员工密码重置

```bash
# 一次修改，所有服务同步更新
curl -X PUT http://localhost:8888/api/unified-users/employee \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "new_secure_password"
  }'

# 结果：
# ✅ Linux密码更新
# ✅ Immich密码更新
# ✅ Nextcloud密码更新
# ✅ Jellyfin密码更新
# ✅ 所有容器权限重新验证
```

### 场景3：员工离职

```bash
# 一次删除，所有服务同步移除
curl -X DELETE http://localhost:8888/api/unified-users/employee \
  -H "Authorization: Bearer YOUR_TOKEN"

# 结果：
# ✅ Linux用户删除
# ✅ Immich账户禁用
# ✅ Nextcloud账户删除
# ✅ Jellyfin账户删除
# ✅ Docker容器权限撤销
```

### 场景4：批量权限管理

```bash
# 修改用户角色和组（所有服务同步）
curl -X PUT http://localhost:8888/api/unified-users/employee \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "role": "admin",
    "groups": ["docker", "media", "storage", "network", "admin"]
  }'

# 结果：
# ✅ 系统权限提升
# ✅ Immich角色更新为管理员
# ✅ Nextcloud组权限更新
# ✅ Jellyfin管理员权限
# ✅ Docker容器sudo权限
```

## 🔧 高级配置示例

### 添加自定义Docker服务

```bash
# 添加Plex媒体服务器
curl -X POST http://localhost:8888/api/unified-users/services \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "plex",
    "type": "plex",
    "config": {
      "apiUrl": "http://localhost:32400",
      "token": "your-plex-token"
    }
  }'
```

### 配置自动同步

```bash
# 在.env.unified-users中设置
UNIFIED_USER_AUTO_SYNC=true
UNIFIED_USER_SYNC_INTERVAL=5m
UNIFIED_USER_RETRY_ATTEMPTS=3
UNIFIED_USER_RETRY_DELAY=2s

# 系统会每5分钟自动同步所有用户
```

### 测试服务连接

```bash
# 测试Immich连接
curl -X POST http://localhost:8888/api/unified-users/services/immich/test \
  -H "Authorization: Bearer YOUR_TOKEN"

# 响应
{
  "message": "Service immich is connected successfully",
  "status": "healthy"
}
```

## 📊 监控和日志

### 查看同步日志

```bash
# 查看详细日志
tail -f /var/log/nas-dashboard/unified-users.log

# 或使用journalctl（如果使用systemd）
sudo journalctl -u nas-dashboard -f
```

### 监控同步状态

```bash
# 创建监控脚本
cat > monitor-sync.sh << 'EOF'
#!/bin/bash

while true; do
    STATUS=$(curl -s http://localhost:8888/api/unified-users/status \
      -H "Authorization: Bearer $TOKEN")
    
    echo "=== $(date) ==="
    echo "$STATUS" | jq '.services[] | {name, status, userCount}'
    
    sleep 60
done
EOF

chmod +x monitor-sync.sh
./monitor-sync.sh
```

## 🐛 故障排除示例

### 问题1：同步部分失败

```bash
# 查看详细错误信息
curl -X POST http://localhost:8888/api/unified-users/sync \
  -H "Authorization: Bearer $TOKEN" | jq '.userResults'

# 重试失败的服务
curl -X POST http://localhost:8888/api/unified-users/sync/username \
  -H "Authorization: Bearer $TOKEN"
```

### 问题2：服务连接失败

```bash
# 测试各个服务连接
for service in system immich nextcloud docker; do
    echo "Testing $service..."
    curl -X POST http://localhost:8888/api/unified-users/services/$service/test \
      -H "Authorization: Bearer $TOKEN"
    echo ""
done
```

### 问题3：Docker容器操作失败

```bash
# 检查容器状态
docker ps | grep nextcloud

# 测试容器连接
docker exec nextcloud php occ user:list

# 手动测试用户创建
docker exec nextcloud php occ user:add testuser --password-from-env
OC_PASS=test123 docker exec nextcloud php occ user:add testuser --password-from-env
```

## 💡 最佳实践示例

### 1. 用户创建脚本

```bash
#!/bin/bash
# 批量创建用户脚本

USERS=(
    "user1:user1@example.com:User One:pass1"
    "user2:user2@example.com:User Two:pass2"
    "user3:user3@example.com:User Three:pass3"
)

TOKEN="your-auth-token"

for user in "${USERS[@]}"; do
    IFS=':' read -r username email name password <<< "$user"
    
    echo "Creating user: $username"
    
    curl -X POST http://localhost:8888/api/unified-users \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"$username\",
        \"email\": \"$email\",
        \"name\": \"$name\",
        \"password\": \"$password\",
        \"role\": \"user\",
        \"groups\": [\"docker\", \"media\"]
      }"
    
    echo "✅ User $username created"
    echo ""
done
```

### 2. 定期同步脚本

```bash
#!/bin/bash
# 定期同步所有用户

TOKEN="your-auth-token"
LOG_FILE="/var/log/user-sync.log"

echo "=== Sync started at $(date) ===" >> $LOG_FILE

RESULT=$(curl -s -X POST http://localhost:8888/api/unified-users/sync \
  -H "Authorization: Bearer $TOKEN")

echo "$RESULT" | jq '.' >> $LOG_FILE

SYNCED=$(echo "$RESULT" | jq -r '.syncedUsers')
FAILED=$(echo "$RESULT" | jq -r '.failedUsers')

echo "✅ Sync completed: $SYNCED synced, $FAILED failed" >> $LOG_FILE
echo "=== Sync completed at $(date) ===" >> $LOG_FILE
```

### 3. 备份和恢复

```bash
# 备份当前用户配置
curl http://localhost:8888/api/unified-users \
  -H "Authorization: Bearer $TOKEN" | jq '.users' > user_backup.json

# 恢复用户（从备份）
for user in $(cat user_backup.json | jq -r '.[].username'); do
    USER_DATA=$(cat user_backup.json | jq ".[] | select(.username==\"$user\")")
    
    curl -X POST http://localhost:8888/api/unified-users \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$USER_DATA"
done
```

## 🎉 总结

通过统一用户管理系统，你可以：

1. **一次操作，多处同步** - 创建/修改/删除用户自动同步到所有服务
2. **简化管理** - 不再需要逐个服务管理用户
3. **减少错误** - 避免手动操作遗漏和错误
4. **提高效率** - 批量操作节省大量时间
5. **统一管理** - 集中管理所有服务的用户

**立即开始使用统一用户管理系统，让用户管理变得简单高效！**