# 统一用户管理系统集成指南

## 🎯 系统概述

统一用户管理系统可以在NAS Dashboard中一次修改用户，所有Docker服务的用户信息自动同步更新，包括：
- 系统用户（Linux）
- Immich（照片管理）
- Nextcloud（文件共享）
- Jellyfin（媒体中心）
- 其他Docker服务

## 🚀 快速开始

### 1. 后端集成

#### **步骤1：在main.go中添加路由注册**

```go
package main

import (
    "nas-dashboard/internal/api"
    // ... 其他导入
)

func main() {
    // ... 现有代码

    // 初始化统一用户管理器
    InitUnifiedUserManager()

    // 注册统一用户路由
    RegisterUnifiedUserRoutes(r, apiGroup)

    // ... 其他代码
}
```

#### **步骤2：编译和运行**

```bash
cd backend
go build -o nas-dashboard cmd/server/main.go
./nas-dashboard
```

### 2. 前端集成

#### **步骤1：添加路由**

```typescript
// router/index.ts
{
  path: '/unified-users',
  name: 'UnifiedUsers',
  component: () => import('../apps/UnifiedUserManager.vue'),
  meta: { requiresAuth: true },
}
```

#### **步骤2：添加桌面图标**

```typescript
// constants/desktop.ts
{
  id: 'unified-users',
  label: '统一用户',
  icon: 'UserGroupIcon',
  appId: 'unified-users',
  badge: null
}
```

#### **步骤3：构建前端**

```bash
cd frontend
npm run build
```

## 📋 API接口文档

### 基础CRUD操作

#### **创建用户**
```http
POST /api/unified-users
Content-Type: application/json

{
  "username": "john",
  "email": "john@example.com",
  "name": "John Doe",
  "password": "secure_password",
  "role": "user",
  "groups": ["docker", "media"],
  "isActive": true
}

响应：
{
  "username": "john",
  "status": "success",
  "details": {
    "system": {
      "serviceName": "system",
      "status": "success"
    },
    "immich": {
      "serviceName": "immich",
      "status": "success"
    },
    "nextcloud": {
      "serviceName": "nextcloud",
      "status": "success"
    }
  },
  "syncTime": "2024-01-01T12:00:00Z"
}
```

#### **更新用户**
```http
PUT /api/unified-users/:username
Content-Type: application/json

{
  "email": "john.new@example.com",
  "name": "John New Doe",
  "password": "new_password",
  "role": "admin",
  "groups": ["docker", "media", "storage"]
}
```

#### **删除用户**
```http
DELETE /api/unified-users/:username

响应：
{
  "username": "john",
  "status": "success",
  "details": { ... }
}
```

### 同步操作

#### **同步所有用户**
```http
POST /api/unified-users/sync

响应：
{
  "totalUsers": 10,
  "syncedUsers": 9,
  "failedUsers": 1,
  "userResults": {
    "user1": { ... },
    "user2": { ... }
  },
  "syncStartTime": "2024-01-01T12:00:00Z",
  "syncEndTime": "2024-01-01T12:01:30Z",
  "duration": "90s"
}
```

#### **同步单个用户**
```http
POST /api/unified-users/sync/:username

响应：
{
  "username": "john",
  "status": "success",
  "details": { ... }
}
```

### 状态查询

#### **获取同步状态**
```http
GET /api/unified-users/status

响应：
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
      "status": "error",
      "userCount": 0,
      "error": "Connection timeout"
    }
  ],
  "lastSyncTime": "2024-01-01T11:30:00Z",
  "autoSyncEnabled": true
}
```

### 服务管理

#### **注册新服务**
```http
POST /api/unified-users/services
Content-Type: application/json

{
  "name": "plex",
  "type": "plex",
  "config": {
    "apiUrl": "http://localhost:32400",
    "token": "your-plex-token"
  }
}
```

#### **注销服务**
```http
DELETE /api/unified-users/services/:serviceName
```

#### **测试服务连接**
```http
POST /api/unified-users/services/:serviceName/test

响应：
{
  "message": "Service immich is connected successfully",
  "status": "healthy"
}
```

## 🔧 配置说明

### 环境变量配置

```bash
# Immich配置
export IMMICH_API_URL="http://localhost:2283/api"
export IMMICH_API_KEY="your-immich-api-key"
export IMMICH_ENABLED="true"

# 同步配置
export UNIFIED_USER_AUTO_SYNC="true"
export UNIFIED_USER_SYNC_INTERVAL="5m"
export UNIFIED_USER_RETRY_ATTEMPTS="3"
export UNIFIED_USER_RETRY_DELAY="2s"
```

### Docker用户服务配置

```go
// 添加Docker容器到统一用户管理
api.AddDockerUser("nextcloud", "nextcloud", map[string]string{
    "admin_user": "admin",
    "admin_password": "password",
})

api.AddDockerUser("jellyfin", "jellyfin", map[string]string{
    "api_key": "your-jellyfin-api-key",
})
```

## 📱 前端使用说明

### 用户管理界面功能

1. **用户列表管理**
   - 查看所有用户
   - 搜索和过滤用户
   - 查看用户同步状态

2. **用户操作**
   - 创建新用户（自动同步到所有服务）
   - 编辑用户信息（自动同步更新）
   - 删除用户（从所有服务删除）
   - 单独同步某个用户

3. **批量操作**
   - 选择多个用户
   - 批量同步
   - 批量删除

4. **服务监控**
   - 查看各服务连接状态
   - 测试服务连接
   - 查看各服务用户数量

### 使用流程

#### **场景1：添加新用户**
1. 点击"添加用户"按钮
2. 填写用户信息
3. 点击"创建并同步"
4. 系统自动在所有服务中创建用户
5. 查看同步结果

#### **场景2：修改用户密码**
1. 找到要修改的用户
2. 点击"编辑"按钮
3. 修改密码字段
4. 点击"更新并同步"
5. 密码自动更新到所有服务

#### **场景3：批量用户同步**
1. 勾选需要同步的用户
2. 点击"批量同步"按钮
3. 等待同步完成
4. 查看同步结果

## 🔍 故障排除

### 常见问题

#### **1. Immich连接失败**
```
错误：Immich service connection failed
解决：
- 检查IMMICH_API_URL是否正确
- 验证IMMICH_API_KEY是否有效
- 确认Immich服务正在运行
```

#### **2. Docker容器用户创建失败**
```
错误：Failed to create user in container nextcloud
解决：
- 检查容器是否正在运行
- 验证容器名称是否正确
- 确认有足够的权限
- 检查容器内服务是否正常
```

#### **3. 同步部分失败**
```
错误：Sync result shows partial success
解决：
- 查看具体失败的服务
- 检查对应服务的日志
- 验证服务配置是否正确
- 单独重试失败的服务
```

### 调试模式

```bash
# 启用调试日志
export LOG_LEVEL=debug

# 查看详细日志
tail -f /var/log/nas-dashboard/app.log

# 测试特定服务
curl -X POST http://localhost:8888/api/unified-users/services/immich/test \
  -H "Authorization: Bearer your-token"
```

## 🎯 最佳实践

### 1. 服务优先级设置

建议按以下顺序设置服务优先级：
1. **system** - 系统用户（主服务）
2. **immich** - 照片管理
3. **nextcloud** - 文件共享
4. **jellyfin** - 媒体中心
5. **other** - 其他服务

### 2. 定期同步策略

```go
// 推荐配置
config := UnifiedUserConfig{
    AutoSync:      true,
    SyncInterval:  5 * time.Minute,  // 5分钟同步一次
    RetryAttempts: 3,
    RetryDelay:    2 * time.Second,
}
```

### 3. 密码管理

- 使用强密码策略
- 定期更新密码
- 避免在多个服务中使用相同密码
- 考虑集成LDAP/AD服务

### 4. 备份策略

```bash
# 备份用户数据
docker exec nextcloud php occ user:list > user_backup_$(date +%Y%m%d).json

# 备份系统用户
getent passwd > system_users_$(date +%Y%m%d).bak
```

## 🔐 安全考虑

### 1. API密钥管理
```bash
# 不要在代码中硬编码API密钥
# 使用环境变量或密钥管理服务
export IMMICH_API_KEY=$(cat /run/secrets/immich_api_key)
```

### 2. 权限控制
```go
// 确保只有管理员可以访问统一用户管理
if user.Role != "admin" {
    c.JSON(403, gin.H{"error": "Forbidden"})
    return
}
```

### 3. 日志审计
```go
// 记录所有用户操作
logger.Info("User operation", logger.LogFields{
    "action": "create",
    "username": user.Username,
    "operator": currentUser,
    "timestamp": time.Now(),
})
```

## 📊 性能优化

### 1. 并发处理
系统已经实现了并发用户同步：
```go
var wg sync.WaitGroup
for serviceName, service := range u.services {
    wg.Add(1)
    go func(name string, svc UserService) {
        defer wg.Done()
        // 处理逻辑
    }(serviceName, service)
}
wg.Wait()
```

### 2. 缓存机制
```go
// 缓存用户列表减少API调用
userCache := make(map[string][]UnifiedUser)
cacheDuration := 5 * time.Minute
```

### 3. 批量操作
```go
// 批量创建用户时考虑使用事务
func (u *UnifiedUserManager) BatchCreateUsers(users []UnifiedUser) error {
    // 实现批量创建逻辑
}
```

## 🚀 扩展功能

### 添加新的Docker服务

#### **示例：添加Plex用户服务**

1. **实现UserService接口**
```go
type PlexUserService struct {
    Config PlexConfig
}

func (p *PlexUserService) CreateUser(user UnifiedUser) error {
    // 调用Plex API创建用户
}

func (p *PlexUserService) UpdateUser(username string, user UnifiedUser) error {
    // 调用Plex API更新用户
}
```

2. **注册服务**
```go
plexService := &PlexUserService{
    Config: PlexConfig{
        APIUrl: "http://localhost:32400",
        Token: os.Getenv("PLEX_TOKEN"),
    },
}

manager.RegisterService("plex", plexService)
```

### 集成LDAP/Active Directory

```go
type LDAPUserService struct {
    Server   string
    BindDN   string
    BindPass string
    BaseDN   string
}

func (l *LDAPUserService) CreateUser(user UnifiedUser) error {
    // 通过LDAP创建用户
    // 自动同步到其他服务
}
```

## 📞 支持和帮助

### 技术支持
- GitHub Issues: [项目地址]
- 文档: [完整文档链接]
- 示例代码: [examples/]

### 常用命令

```bash
# 检查服务状态
curl http://localhost:8888/api/unified-users/status \
  -H "Authorization: Bearer your-token"

# 同步所有用户
curl -X POST http://localhost:8888/api/unified-users/sync \
  -H "Authorization: Bearer your-token"

# 测试特定服务
curl -X POST http://localhost:8888/api/unified-users/services/immich/test \
  -H "Authorization: Bearer your-token"
```

---

**总结**：统一用户管理系统可以极大简化多服务环境下的用户管理，一次操作自动同步到所有服务，大大提高管理效率。