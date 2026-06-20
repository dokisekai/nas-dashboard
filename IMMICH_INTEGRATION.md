# Immich 用户管理集成指南

## 功能概述

现在NAS控制面板已经集成了Immich用户管理功能，允许用户在统一界面中管理Immich应用的用户。

## 新增功能

### 1. 应用权限管理标签页
在用户管理界面新增了"应用权限"标签页，可以管理各个应用的用户访问权限。

### 2. Immich用户管理
- 查看Immich用户列表
- 添加/编辑/删除Immich用户
- 批量修改用户状态和角色
- 同步系统用户到Immich

### 3. 批量操作功能
- 多选用户进行批量修改
- 批量修改用户状态（启用/停用）
- 批量修改用户角色（管理员/普通用户）

## 后端API端点

新增以下Immich管理API：

```
GET    /api/immich/users          - 获取Immich用户列表
GET    /api/immich/users/:id      - 获取单个Immich用户
POST   /api/immich/users          - 创建Immich用户
PUT    /api/immich/users/:id      - 更新Immich用户
DELETE /api/immich/users/:id      - 删除Immich用户
POST   /api/immich/users/batch    - 批量更新Immich用户
POST   /api/immich/users/sync     - 同步系统用户到Immich
```

## 配置步骤

### 1. 获取Immich API密钥

1. 登录Immich管理界面: http://localhost:2283
2. 点击右上角用户图标 -> Administration
3. 选择User Settings -> API Keys
4. 点击"New API Key"生成新的密钥
5. 复制生成的密钥

### 2. 配置后端

编辑 `/data/nas-dashboard/backend/.env.immich` 文件：

```
IMMICH_API_URL=http://localhost:2283/api
IMMICH_API_KEY=你的Immich_API密钥
```

### 3. 重新编译后端

```bash
cd /data/nas-dashboard/backend
go build -o main cmd/server/main.go
```

### 4. 重启服务

```bash
cd /data/nas-dashboard
./start-dashboard.sh
```

## 使用说明

### 访问应用权限管理

1. 登录NAS控制面板
2. 进入"用户管理"应用
3. 点击"应用权限"标签页
4. 点击"Immich"展开用户管理界面

### 基础操作

**添加用户：**
1. 点击"添加用户"按钮
2. 填写邮箱、姓名、密码
3. 选择状态和角色
4. 点击"保存"

**编辑用户：**
1. 点击用户行的编辑图标
2. 修改用户信息
3. 点击"保存"

**删除用户：**
1. 点击用户行的删除图标
2. 确认删除

**批量操作：**
1. 勾选要操作的用户
2. 点击批量操作按钮（批量修改状态/角色）
3. 选择修改内容
4. 确认操作

### 同步系统用户

1. 点击"同步系统用户"按钮
2. 确认同步操作
3. 系统将自动创建或更新对应的Immich用户

## 扩展其他应用

其他应用（如文件管理器、Docker管理、系统监控）的用户管理可以按照相同的方式扩展：

1. 在 `applications` 数组中添加应用信息
2. 创建对应的用户管理组件
3. 在模板中添加对应的组件引用

## 注意事项

1. **API密钥安全**：请妥善保管Immich API密钥，不要泄露
2. **权限要求**：Immich用户管理需要Immich管理员权限
3. **数据一致性**：同步系统用户时会创建新的Immich用户，请确保邮箱唯一性
4. **网络连接**：确保后端能够访问Immich API端点

## 故障排除

### API连接失败
- 检查Immich服务是否运行
- 确认API URL配置正确
- 验证API密钥是否有效

### 用户操作失败
- 检查用户权限是否足够
- 确认用户数据格式正确
- 查看后端日志获取详细错误信息

### 同步失败
- 确认系统用户数据格式正确
- 检查Immich API连接状态
- 验证邮箱地址格式

## 文件清单

后端文件：
- `/data/nas-dashboard/backend/internal/api/immich.go` - Immich API实现
- `/data/nas-dashboard/backend/cmd/server/main.go` - 路由注册
- `/data/nas-dashboard/backend/.env.immich` - 配置文件

前端文件：
- `/data/nas-dashboard/frontend/src/components/ImmichUserManager.vue` - Immich用户管理组件
- `/data/nas-dashboard/frontend/src/apps/UserManager.vue` - 用户管理主界面

## 更新日志

**v1.0.0 (2026-06-20)**
- 初始版本
- 支持Immich用户管理
- 支持批量操作功能
- 支持系统用户同步