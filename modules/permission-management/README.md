# 权限管理模块

## 概述

权限管理模块提供细粒度的文件和共享权限控制，支持多种权限模型和访问控制机制。

## 功能特性

### 权限类型
- 文件系统权限 (Unix权限)
- 访问控制列表 (ACL)
- 共享文件夹权限
- SMB/CIFS权限
- NFS权限设置

### 权限模型
- 用户级别权限
- 用户组级别权限
- 其他用户权限
- 自定义权限规则
- 权限继承机制

### 权限管理
- 权限查看和修改
- 批量权限设置
- 权限模板管理
- 权限继承配置
- 权限验证和测试

## 后端API实现

### 主要API文件
- `backend/internal/api/permissions.go` - 权限管理API
- `backend/internal/api/file_management.go` - 文件权限API

### 前端组件
- `frontend/src/apps/SimplePermissionManager.vue` - 简单权限管理
- `frontend/src/apps/ShareFolderManager.vue` - 共享文件夹权限管理
- `frontend/src/components/ControlPanel/PermissionManager.vue` - 权限管理组件

## API端点

### 共享权限管理
- `GET /api/permissions/shares` - 获取共享列表
- `POST /api/permissions/shares` - 创建共享
- `PUT /api/permissions/shares/:name` - 更新共享
- `DELETE /api/permissions/shares/:name` - 删除共享
- `GET /api/permissions/shares/:name/permissions` - 获取共享权限
- `PUT /api/permissions/shares/:name/permissions` - 设置共享权限

### 文件权限管理
- `GET /api/permissions/files` - 获取文件权限
- `PUT /api/permissions/files/permissions` - 设置文件权限
- `GET /api/permissions/files/acl` - 获取文件ACL
- `PUT /api/permissions/files/acl` - 设置文件ACL

### 权限验证
- `POST /api/permissions/verify` - 验证权限
- `GET /api/permissions/effective` - 获取有效权限
- `POST /api/permissions/test` - 测试权限访问

## 使用示例

### 创建共享并设置权限
```bash
# 创建共享
curl -X POST http://localhost:8888/api/permissions/shares \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "shared-docs",
    "path": "/mnt/storage/shared-docs",
    "description": "Shared documents folder",
    "enable_recycle_bin": true
  }'

# 设置权限
curl -X PUT http://localhost:8888/api/permissions/shares/shared-docs/permissions \
  -H "Authorization: Bearer <token>" \
  -d '{
    "permissions": [
      {
        "type": "user",
        "name": "john",
        "permissions": "read,write"
      },
      {
        "type": "group",
        "name": "developers",
        "permissions": "read,write,execute"
      }
    ]
  }'
```

### 设置文件权限
```bash
curl -X PUT http://localhost:8888/api/permissions/files/permissions \
  -H "Authorization: Bearer <token>" \
  -d '{
    "path": "/mnt/storage/docs",
    "mode": "0755",
    "user": "john",
    "group": "developers"
  }'
```

### 设置文件ACL
```bash
curl -X PUT http://localhost:8888/api/permissions/files/acl \
  -H "Authorization: Bearer <token>" \
  -d '{
    "path": "/mnt/storage/docs/file.txt",
    "acl": [
      "user:john:rw",
      "group:developers:rw",
      "user:mary:r"
    ]
  }'
```

## 权限类型详解

### Unix文件权限
```
-rwxr-xr--
│││││││││
│││││││└─ 其他用户权限 (r-x)
││││││└─── 用户组权限 (r-x)
│││││└───── 文件所有者权限 (rwx)
││││└────── 文件类型 (- = 普通文件, d = 目录)
```

### 权限代码
- **r** (read) = 4 - 读取权限
- **w** (write) = 2 - 写入权限
- **x** (execute) = 1 - 执行权限

### 常见权限组合
- **755** (rwxr-xr-x) - 所有者完全权限，其他用户读执行
- **644** (rw-r--r--) - 所有者读写，其他用户只读
- **777** (rwxrwxrwx) - 所有用户完全权限
- **700** (rwx------) - 所有者完全权限，其他用户无权限

## ACL权限

### ACL语法
```bash
# 设置ACL
setfacl -m u:john:rw /path/to/file
setfacl -m g:developers:rw /path/to/file

# 查看ACL
getfacl /path/to/file

# 删除ACL
setfacl -x u:john /path/to/file
```

### ACL权限类型
- **user:** - 指定用户权限
- **group:** - 指定用户组权限
- **other:** - 其他用户权限
- **mask:** - 权限掩码

## SMB权限

### SMB权限类型
- **FULL** - 完全控制
- **CHANGE** - 修改权限
- **READ** - 只读权限
- **NO ACCESS** - 无权限

### SMB权限配置
```ini
[shared-docs]
    path = /mnt/storage/shared-docs
    browseable = yes
    writable = yes
    valid users = @developers
    write list = john
    read list = @team
```

## 权限继承

### 继承规则
- 子目录继承父目录权限
- 新文件继承目录默认权限
- ACL权限继承
- 权限传播设置

### 继承配置
```bash
# 设置默认ACL
setfacl -d -m u:john:rw /path/to/directory

# 设置继承标志
setfacl -m default:user:john:rw /path/to/directory
```

## 配置文件

### 权限配置
- `config/permissions/acl.conf` - ACL配置
- `config/permissions/smb.conf` - SMB权限配置
- `config/permissions/templates.conf` - 权限模板
- `config/permissions/default.conf` - 默认权限设置

### 权限模板
```json
{
  "templates": [
    {
      "name": "default_share",
      "user_permissions": "rw",
      "group_permissions": "r",
      "other_permissions": "r"
    },
    {
      "name": "private_share",
      "user_permissions": "rw",
      "group_permissions": "",
      "other_permissions": ""
    }
  ]
}
```

## 权限验证

### 验证流程
1. 检查文件存在性
2. 获取文件所有者
3. 检查用户身份
4. 验证权限匹配
5. 检查ACL权限
6. 返回权限结果

### 验证命令
```bash
# 检查文件权限
stat /path/to/file

# 验证访问权限
access -r /path/to/file  # 检查读权限
access -w /path/to/file  # 检查写权限
access -x /path/to/file  # 检查执行权限

# 测试权限
sudo -u john cat /path/to/file
```

## 故障排除

### 权限设置失败
1. 检查文件所有权：`ls -l /path/to/file`
2. 验证当前用户权限：`whoami`
3. 检查父目录权限：`ls -ld /path/to/`
4. 查看系统日志：`journalctl -f`

### ACL不生效
1. 检查文件系统支持：`tune2fs -l /dev/sdX | grep "Filesystem features"`
2. 验证ACL挂载选项：`mount | grep acl`
3. 检查ACL工具：`which getfacl`
4. 重新挂载文件系统：`mount -o remount,acl /mount/point`

### SMB权限问题
1. 检查SMB配置：`testparm -s`
2. 验证用户权限：`smbclient -L localhost`
3. 查看SMB日志：`tail -f /var/log/samba/log.smbd`
4. 重新加载配置：`systemctl reload smbd`

## 安全考虑

### 权限安全
- 最小权限原则
- 定期权限审计
- 监控权限变更
- 权限异常告警

### 敏感目录权限
- `/home/` - 用户主目录，权限750
- `/etc/` - 系统配置，权限600-700
- `/var/log/` - 日志文件，权限600-700
- `/root/` - 管理员目录，权限700

### 权限建议
- 避免使用777权限
- 使用用户组权限管理
- 定期审查特殊权限
- 监控SUID/SGID文件

## 性能优化

### 权限检查优化
- 缓存权限信息
- 使用索引查询
- 批量权限验证
- 异步权限更新

### ACL性能
- 限制ACL条目数量
- 使用继承减少重复
- 定期清理无效ACL
- 监控ACL性能影响

## 监控和日志

### 权限监控
- 权限变更记录
- 访问失败统计
- 异常权限告警
- 权限使用分析

### 审计日志
- 权限修改操作
- 权限验证请求
- 访问控制决策
- 权限相关错误

## 相关文档

- [用户管理模块](../user-management/)
- [共享管理模块](../share-management/)
- [安全考虑](../../docs/SECURITY_CONSIDERATIONS.md)
- [API文档](../../docs/API_DOCUMENTATION.md)

## DSM对标功能

本模块对标Synology DSM的权限管理功能：
- ✅ 文件权限管理
- ✅ ACL支持
- ✅ 共享文件夹权限
- ✅ SMB权限控制
- ✅ 权限继承
- ✅ 权限模板

功能完整度: 90% ✅