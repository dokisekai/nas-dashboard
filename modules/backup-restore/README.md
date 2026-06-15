# 备份恢复模块

## 概述

备份恢复模块提供完整的数据备份和恢复解决方案，支持多种备份方式和存储目标。

## 功能特性

### 备份类型
- 完整备份
- 增量备份
- 差异备份
- 快照备份

### 备份工具
- Restic备份
- Rsync同步
- Btrfs快照
- 数据库备份

### 存储目标
- 本地存储
- 远程服务器
- 云存储（S3, B2, etc.）
- 网络存储（NFS, SMB）

### 恢复功能
- 完整恢复
- 单文件恢复
- 时间点恢复
- 增量恢复

## 脚本说明

### 1. restic-backup.sh
基于Restic的通用备份脚本。

```bash
# 使用示例
./restic-backup.sh --source /home/user --dest /backup/weekly
```

### 2. restic-home-backup.sh
Home目录专用备份脚本。

```bash
# 使用示例
./restic-home-backup.sh --exclude ".cache"
```

## 后端API实现

### 主要API文件
- `backend/internal/api/backup.go` - 备份管理API
- `backend/internal/api/sync.go` - 同步管理API
- `backend/service/backup/` - 备份服务

### 前端组件
- `frontend/src/apps/BackupManager.vue` - 备份管理界面
- `frontend/src/apps/SyncManager.vue` - 同步管理界面

## API端点

### Restic备份管理
- `GET /api/storage/backup/repos` - 获取备份仓库列表
- `POST /api/storage/backup/repos` - 创建备份仓库
- `GET /api/storage/backup/tasks` - 获取备份任务列表
- `POST /api/storage/backup/tasks` - 创建备份任务

### 同步任务管理
- `GET /api/storage/sync/jobs` - 获取同步任务列表
- `POST /api/storage/sync/jobs` - 创建同步任务
- `POST /api/storage/sync/jobs/:id/run` - 运行同步任务

### 系统备份恢复
- `GET /api/backups` - 获取系统备份列表
- `POST /api/backups` - 创建系统备份
- `GET /api/backups/:id` - 获取备份详情
- `DELETE /api/backups/:id` - 删除备份
- `POST /api/backups/restore` - 恢复备份

## 使用示例

### 创建Restic备份仓库
```bash
curl -X POST http://localhost:8888/api/storage/backup/repos \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "weekly_backup",
    "path": "/mnt/backup/weekly",
    "password": "secure_password",
    "type": "local"
  }'
```

### 创建备份任务
```bash
curl -X POST http://localhost:8888/api/storage/backup/tasks \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "daily_backup",
    "repo_id": 1,
    "source": "/home/user",
    "schedule": "0 2 * * *",
    "retention": {
      "daily": 7,
      "weekly": 4,
      "monthly": 12
    }
  }'
```

### 创建同步任务
```bash
curl -X POST http://localhost:8888/api/storage/sync/jobs \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "photo_sync",
    "source": "/home/user/photos",
    "destination": "/mnt/storage/photos",
    "method": "rsync",
    "schedule": "*/30 * * * *",
    "delete": false
  }'
```

## 配置文件

### Restic配置
- `config/restic/repos.conf` - 仓库配置
- `config/restic/backup.conf` - 备份配置
- `config/restic/excludes.txt` - 排除规则

### Rsync配置
- `config/rsync/sync.conf` - 同步配置
- `config/rsync/excludes.txt` - 排除规则

### 备份策略
- `config/backup/policy.conf` - 备份策略
- `config/backup/retention.conf` - 保留策略

## 备份策略

### 3-2-1备份策略
- **3** 份数据副本（1份原始 + 2份备份）
- **2** 种不同存储类型
- **1** 份异地备份

### 备份层级
1. **每日备份**：保留7天
2. **每周备份**：保留4周
3. **每月备份**：保留12个月
4. **年度备份**：永久保留

### 备份验证
- 自动完整性检查
- 定期恢复测试
- 备份大小监控
- 备份时间跟踪

## 恢复流程

### 完整恢复
```bash
# 从Restic仓库恢复
restic restore latest --target /restore/point

# 验证恢复数据
restic check --read-data
```

### 单文件恢复
```bash
# 恢复特定文件
restic restore latest --target /tmp --include "/home/user/file.txt"
```

### 时间点恢复
```bash
# 恢复到特定时间点
restic restore --time "2024-01-01 10:00" --target /restore/point
```

## 监控和警报

### 备份监控
- 备份成功率
- 备份持续时间
- 存储空间使用
- 备份完整性

### 警报条件
- 备份失败
- 备份超时
- 存储空间不足
- 恢复测试失败

## 性能优化

- 增量备份减少数据传输
- 并行备份提高速度
- 压缩减少存储空间
- 去重减少重复数据

## 安全建议

- 加密备份数据
- 保护备份密码
- 定期测试恢复
- 监控备份访问
- 维护离线备份

## 故障排除

### 备份失败
1. 检查存储空间：`df -h`
2. 验证网络连接：`ping backup-server`
3. 查看备份日志：`journalctl -u restic-backup`

### 恢复失败
1. 验证备份完整性：`restic check`
2. 检查目标空间：`df -h /restore/point`
3. 测试部分恢复

### 存储空间不足
1. 清理旧备份：`restic forget --keep-daily 7`
2. 压缩备份数据
3. 扩展存储容量

## 相关文档

- [备份管理API文档](../../docs/API_DOCUMENTATION.md)
- [用户指南](../../docs/USER_GUIDE.md)
- [故障排除指南](../../docs/TROUBLESHOOTING.md)
- [Restic官方文档](https://restic.readthedocs.io/)