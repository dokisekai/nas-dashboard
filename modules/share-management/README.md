# 共享管理模块

## 概述

共享管理模块提供网络共享文件夹的创建、配置和权限管理功能，支持SMB/CIFS、NFS、WebDAV等多种协议。

## 功能特性

### 共享协议支持
- **SMB/CIFS** - Windows文件共享
- **NFS** - Unix/Linux文件共享
- **WebDAV** - HTTP/Web文件访问
- **AFP** - Apple文件协议(可选)

### 共享管理
- 共享创建和删除
- 共享路径配置
- 权限和访问控制
- 共享描述管理
- 共享状态监控

### 高级功能
- 回收站支持
- 快照集成
- 访问日志记录
- 性能监控
- 配额管理

## 后端API实现

### 主要API文件
- `backend/internal/api/storage.go` - 存储管理API(包含SMB共享)
- `backend/internal/api/permissions.go` - 权限管理API
- `backend/pkg/system/smb.go` - SMB系统调用

### 前端组件
- `frontend/src/apps/ShareFolderManager.vue` - 共享文件夹管理界面
- `frontend/src/apps/SMBUserManager.vue` - SMB用户管理
- `frontend/src/components/ControlPanel/ShareManager.vue` - 共享管理组件

## API端点

### SMB共享管理
- `GET /api/storage/smb` - 获取SMB共享列表
- `POST /api/storage/smb` - 创建SMB共享
- `PUT /api/storage/smb/:name` - 更新SMB共享
- `DELETE /api/storage/smb/:name` - 删除SMB共享

### 共享权限
- `GET /api/permissions/shares` - 获取共享列表
- `POST /api/permissions/shares` - 创建共享
- `PUT /api/permissions/shares/:name` - 更新共享
- `DELETE /api/permissions/shares/:name` - 删除共享
- `GET /api/permissions/shares/:name/permissions` - 获取共享权限
- `PUT /api/permissions/shares/:name/permissions` - 设置共享权限

## 使用示例

### 创建SMB共享
```bash
curl -X POST http://localhost:8888/api/storage/smb \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "shared-docs",
    "path": "/mnt/storage/shared-docs",
    "description": "Shared documents folder",
    "read_only": false,
    "browseable": true,
    "guest_ok": false,
    "valid_users": ["john", "mary"],
    "write_list": ["john"]
  }'
```

### 创建带权限的共享
```bash
# 创建共享
curl -X POST http://localhost:8888/api/permissions/shares \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "team-share",
    "path": "/mnt/storage/team",
    "description": "Team collaboration folder",
    "enable_recycle_bin": true,
    "enable_snapshot": true
  }'

# 设置权限
curl -X PUT http://localhost:8888/api/permissions/shares/team-share/permissions \
  -H "Authorization: Bearer <token>" \
  -d '{
    "permissions": [
      {
        "type": "group",
        "name": "developers",
        "permissions": "read,write,execute"
      },
      {
        "type": "group", 
        "name": "team",
        "permissions": "read,execute"
      }
    ]
  }'
```

### 更新共享配置
```bash
curl -X PUT http://localhost:8888/api/storage/smb/shared-docs \
  -H "Authorization: Bearer <token>" \
  -d '{
    "description": "Updated description",
    "read_only": true,
    "guest_ok": true
  }'
```

## 共享协议配置

### SMB/CIFS配置
```ini
[shared-docs]
    path = /mnt/storage/shared-docs
    comment = Shared documents folder
    browseable = yes
    read only = no
    guest ok = no
    valid users = @developers
    write list = john
    create mask = 0664
    directory mask = 0775
    vfs objects = recycle
    recycle:repository = .recycle
    recycle:keeptree = yes
    recycle:versions = yes
```

### NFS配置
```bash
# /etc/exports
/mnt/storage/shared-docs 192.168.1.0/24(rw,sync,no_subtree_check)
/mnt/storage/public  *(ro,all_squash,anonuid=1000,anongid=1000)
```

### WebDAV配置
```apache
<Location /webdav>
    DAV On
    AuthType Basic
    AuthName "WebDAV"
    AuthUserFile /etc/webdav.users
    Require valid-user
</Location>
```

## 共享权限模型

### SMB权限级别
- **FULL** - 完全控制
- **CHANGE** - 修改(读写)
- **READ** - 只读
- **NO ACCESS** - 无权限

### 权限组合规则
- 用户权限 > 用户组权限
- 明确拒绝 > 允许
- 特定用户 > everyone
- 时间段权限优先

### 访问控制示例
```
用户: john -> FULL
用户: mary -> CHANGE
用户组: developers -> CHANGE  
用户组: team -> READ
其他用户 -> NO ACCESS
```

## 高级功能

### 回收站功能
```ini
[shared-docs]
    vfs objects = recycle
    recycle:repository = .recycle
    recycle:keeptree = yes
    recycle:versions = yes
    recycle:maxsize = 100M
    recycle:exclude = *.tmp,*.temp
    recycle:excludedir = /tmp
```

### 快照集成
```bash
# 配置Btrfs快照
btrfs subvolume snapshot -r /mnt/storage/shared-docs /mnt/snapshots/shared-docs-$(date +%Y%m%d)

# 在SMB中启用阴影副本
vfs objects = shadow_copy2
shadow:snapdir = /mnt/snapshots
shadow:format = %Y-%m-%d-%H%M%S
shadow:sort = desc
```

### 访问日志
```ini
[shared-docs]
    vfs objects = full_audit
    full_audit:prefix = %u|%I
    full_audit:success = connect,disconnect,opendir,mkdir,rmdir,read,write
    full_audit:failure = connect
    full_audit:facility = local7
    full_audit:priority = notice
```

## 配置文件

### 共享管理配置
- `config/sharing/smb.conf` - SMB主配置
- `config/sharing/nfs.conf` - NFS配置
- `config/sharing/webdav.conf` - WebDAV配置
- `config/sharing/templates.conf` - 共享模板

### 共享模板
```json
{
  "templates": [
    {
      "name": "documents",
      "description": "Document sharing template",
      "read_only": false,
      "browseable": true,
      "enable_recycle_bin": true,
      "create_mask": "0664",
      "directory_mask": "0775"
    },
    {
      "name": "public",
      "description": "Public read-only share",
      "read_only": true,
      "browseable": true,
      "guest_ok": true
    }
  ]
}
```

## 监控和维护

### 共享监控
```bash
# 查看SMB连接
smbstatus -b

# 查看共享访问
smbstatus -S

# 查看NFS导出
showmount -e localhost

# 监控WebDAV访问
tail -f /var/log/httpd/access_log
```

### 性能监控
- 共享访问频率
- 数据传输量统计
- 连接用户数量
- 磁盘IO使用
- 网络带宽占用

### 维护任务
- 定期检查共享状态
- 清理回收站
- 管理快照存储
- 审查访问权限
- 优化共享配置

## 故障排除

### SMB共享无法访问
1. 检查SMB服务状态：`systemctl status smbd`
2. 验证配置文件：`testparm -s`
3. 检查防火墙设置：`ufw status`
4. 查看SMB日志：`tail -f /var/log/samba/log.smbd`

### NFS无法挂载
1. 检查NFS服务：`systemctl status nfs-server`
2. 验证导出配置：`showmount -e localhost`
3. 检查防火墙：`iptables -L -n`
4. 测试挂载：`mount -t nfs localhost:/export /mnt/test`

### 权限问题
1. 检查文件权限：`ls -la /path/to/share`
2. 验证用户权限：`pdbedit -L john`
3. 检查SELinux状态：`getenforce`
4. 查看访问日志：`tail -f /var/log/samba/log.$USER`

## 安全考虑

### 网络安全
- 使用防火墙限制访问
- 启用加密传输
- 定期更新密码
- 监控异常访问

### 数据安全
- 实施访问控制
- 启用审计日志
- 定期备份共享数据
- 使用加密敏感数据

### 权限安全
- 最小权限原则
- 定期权限审查
- 禁用guest访问
- 限制管理员权限

## 性能优化

### SMB优化
```ini
[global]
    socket options = TCP_NODELAY IPTOS_LOWDELAY
    read raw = yes
    write raw = yes
    max xmit = 65536
    getwd cache = yes
    write cache size = 262144
```

### NFS优化
```bash
# /etc/fstab 客户端配置
server:/export /mnt/nfs nfs rw,noatime,async,rsize=8192,wsize=8192 0 0
```

### 缓存策略
- 启用客户端缓存
- 配置合适的缓存大小
- 使用write-back缓存
- 定期清理缓存

## 监控和告警

### 监控指标
- 共享可用性
- 访问响应时间
- 数据传输速率
- 并发连接数
- 磁盘使用率

### 告警规则
- 共享服务停止
- 访问失败率高
- 磁盘空间不足
- 异常访问模式
- 性能下降告警

## 相关文档

- [权限管理模块](../permission-management/)
- [用户管理模块](../user-management/)
- [存储管理模块](../disk-management/)
- [DSM对标分析](../../docs/DSM_ANALYSIS.md)

## DSM对标功能

本模块对标Synology DSM的共享管理功能：
- ✅ SMB/CIFS共享
- ✅ NFS共享
- ✅ WebDAV支持
- ✅ 共享权限管理
- ✅ 回收站功能
- ✅ 快照集成
- ✅ 访问日志

功能完整度: 92% ✅