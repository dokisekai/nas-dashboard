# 用户管理模块

## 概述

用户管理模块提供完整的用户账户管理功能，包括用户创建、修改、删除、权限控制和配额管理。

## 功能特性

### 用户账户管理
- 用户创建和删除
- 用户信息修改
- 密码管理
- 用户状态控制
- 用户组管理

### 权限控制
- 基于角色的访问控制
- 用户组权限设置
- 细粒度权限配置
- 权限继承机制

### 配额管理
- 磁盘空间配额
- 文件数量配额
- 目录级别配额
- 配额使用监控

### SSH密钥管理
- SSH公钥添加
- 密钥权限管理
- 密钥使用统计
- 密钥过期管理

## 后端API实现

### 主要API文件
- `backend/internal/api/user.go` - 用户管理API
- `backend/internal/api/group.go` - 用户组API
- `backend/internal/api/quota.go` - 配额管理API
- `backend/internal/api/permissions.go` - 权限管理API

### 前端组件
- `frontend/src/apps/UserManager.vue` - 用户管理界面
- `frontend/src/apps/SMBUserManager.vue` - SMB用户管理
- `frontend/src/components/ControlPanel/UserManager.vue` - 控制面板用户管理

## API端点

### 用户管理
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `GET /api/users/:username` - 获取用户详情
- `PUT /api/users/:username` - 更新用户信息
- `DELETE /api/users/:username` - 删除用户
- `GET /api/users/me` - 获取当前用户信息
- `POST /api/users/me/password` - 修改当前用户密码

### 用户组管理
- `GET /api/groups` - 获取用户组列表
- `POST /api/groups` - 创建用户组
- `GET /api/groups/:name` - 获取用户组详情
- `PUT /api/groups/:name` - 更新用户组
- `DELETE /api/groups/:name` - 删除用户组
- `GET /api/groups/:name/members` - 获取用户组成员
- `POST /api/groups/:name/members` - 添加用户组成员
- `DELETE /api/groups/:name/members/:user` - 移除用户组成员

### 配额管理
- `GET /api/users/:username/quota` - 获取用户配额
- `PUT /api/users/:username/quota` - 设置用户配额
- `GET /api/storage/quota/users` - 获取所有用户配额
- `GET /api/storage/quota/groups` - 获取所有用户组配额
- `GET /api/storage/quota/report` - 获取配额报告

### SSH密钥管理
- `GET /api/users/ssh-keys` - 获取SSH密钥列表
- `POST /api/users/ssh-keys` - 添加SSH密钥
- `DELETE /api/users/ssh-keys/:id` - 删除SSH密钥

### SMB用户管理
- `GET /api/smb/users` - 获取SMB用户列表
- `POST /api/smb/users/:username/password` - 设置SMB用户密码
- `DELETE /api/smb/users/:username/password` - 删除SMB用户密码
- `POST /api/smb/users/:username/enable` - 启用SMB用户
- `POST /api/smb/users/:username/disable` - 禁用SMB用户
- `GET /api/smb/users/:username/stats` - 获取SMB用户统计
- `GET /api/smb/sessions` - 获取SMB会话列表
- `DELETE /api/smb/sessions/:pid` - 断开指定SMB会话
- `DELETE /api/smb/sessions` - 断开所有SMB会话

## 使用示例

### 创建用户
```bash
curl -X POST http://localhost:8888/api/users \
  -H "Authorization: Bearer <token>" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "secure_password",
    "role": "user",
    "home_directory": "/home/john",
    "shell": "/bin/bash"
  }'
```

### 设置用户配额
```bash
curl -X PUT http://localhost:8888/api/users/john/quota \
  -H "Authorization: Bearer <token>" \
  -d '{
    "path": "/home/john",
    "soft_limit": 10737418240,
    "hard_limit": 16106127360,
    "grace_period": 604800
  }'
```

### 添加SSH密钥
```bash
curl -X POST http://localhost:8888/api/users/ssh-keys \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "My Key",
    "public_key": "ssh-rsa AAAAB3NzaC1yc2E...",
    "user_id": 123
  }'
```

### 创建用户组
```bash
curl -X POST http://localhost:8888/api/groups \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "developers",
    "description": "Developer group",
    "gid": 1001
  }'
```

## 用户角色

### 角色类型
- **admin** - 系统管理员，完全访问权限
- **user** - 普通用户，基本访问权限
- **guest** - 访客用户，只读访问权限

### 权限矩阵
| 功能 | Admin | User | Guest |
|------|-------|------|-------|
| 用户管理 | ✅ | ❌ | ❌ |
| 系统配置 | ✅ | ❌ | ❌ |
| 存储管理 | ✅ | ✅ | ❌ |
| 文件访问 | ✅ | ✅ | ✅ (只读) |
| 网络配置 | ✅ | ❌ | ❌ |

## 配置文件

### 用户管理配置
- `config/users/users.conf` - 用户管理配置
- `config/users/groups.conf` - 用户组配置
- `config/users/quota.conf` - 配额配置
- `config/users/ssh.conf` - SSH密钥配置

### 默认设置
```json
{
  "default_shell": "/bin/bash",
  "default_home": "/home",
  "default_quota": {
    "soft": "10G",
    "hard": "15G"
  },
  "password_policy": {
    "min_length": 8,
    "require_uppercase": true,
    "require_numbers": true
  }
}
```

## 系统集成

### 与其他模块集成
- **权限管理模块** - 提供用户权限控制
- **配额管理模块** - 实现磁盘空间限制
- **共享管理模块** - 控制共享访问权限
- **日志模块** - 记录用户操作日志

### 认证流程
1. 用户登录请求
2. 验证用户凭证
3. 生成JWT令牌
4. 返回访问令牌
5. 后续请求携带令牌

## 安全考虑

### 密码安全
- 密码强度要求
- 密码加密存储
- 密码过期策略
- 登录失败限制

### 会话管理
- JWT令牌验证
- 会话超时控制
- 会话撤销机制
- 多设备管理

### 审计日志
- 用户操作记录
- 登录历史跟踪
- 权限变更记录
- 异常行为监控

## 故障排除

### 用户创建失败
1. 检查用户是否已存在：`getent passwd username`
2. 验证系统权限：`sudo useradd -D`
3. 检查磁盘空间：`df -h /home`
4. 查看系统日志：`journalctl -u systemd-logind`

### 配额设置失败
1. 检查配额是否启用：`quotaon -p /home`
2. 验证配额配置：`repquota -a`
3. 检查文件系统支持：`quotacheck -cugm /home`
4. 查看配额状态：`quota -u username`

### SSH密钥添加失败
1. 验证密钥格式：`ssh-keygen -l -f key.pub`
2. 检查目录权限：`ls -la ~/.ssh/`
3. 验证authd配置：`cat /etc/ssh/sshd_config | grep AuthorizedKeysFile`
4. 检查SELinux状态：`getenforce`

## 性能优化

### 用户查询优化
- 使用索引字段查询
- 实现用户缓存机制
- 分页加载用户列表
- 异步处理大量用户

### 配额监控优化
- 定期批量更新配额信息
- 使用增量更新策略
- 缓存配额计算结果
- 异步发送配额告警

## 最佳实践

### 用户命名规范
- 使用小写字母
- 避免特殊字符
- 保持名称简洁
- 使用有意义的用户名

### 密码策略
- 最小长度8位
- 包含大小写字母
- 包含数字和特殊字符
- 定期更换密码

### 配额设置
- 根据用户角色设置
- 预留系统空间
- 定期审查使用情况
- 及时调整不合理配额

## 相关文档

- [用户指南](../../docs/USER_GUIDE.md)
- [安全考虑](../../docs/SECURITY_CONSIDERATIONS.md)
- [API文档](../../docs/API_DOCUMENTATION.md)
- [权限管理模块](../permission-management/)

## DSM对标功能

本模块对标Synology DSM的用户管理功能：
- ✅ 用户账户管理
- ✅ 用户组管理
- ✅ 权限控制
- ✅ 配额管理
- ✅ SSH密钥管理
- ✅ SMB用户管理
- ✅ 会话管理

功能完整度: 95% ✅