# 磁盘管理模块

## 概述

磁盘管理模块提供完整的磁盘存储管理功能，包括磁盘分区、格式化、RAID配置、LVM管理和存储池管理。

## 功能特性

### 磁盘基础管理
- 磁盘信息查询
- 磁盘分区创建/删除
- 磁盘格式化（ext4, xfs, btrfs）
- 磁盘挂载/卸载
- SMART健康监控

### RAID管理
- RAID级别：0, 1, 5, 6, 10
- RAID阵列创建和管理
- 磁盘热备份
- RAID状态监控

### LVM管理
- 物理卷(PV)管理
- 卷组(VG)管理  
- 逻辑卷(LV)管理
- 卷扩容和缩容
- 卷快照

### 存储池管理
- Btrfs存储池
- 存储池创建和管理
- 子卷管理
- 存储池监控

## 后端API实现

### 主要API文件
- `backend/internal/api/storage.go` - 存储基础API
- `backend/internal/api/disk_partition.go` - 分区管理
- `backend/internal/api/disk_raid.go` - RAID管理
- `backend/internal/api/disk_lvm.go` - LVM管理
- `backend/internal/api/disk_smart.go` - SMART监控
- `backend/internal/api/storage_pool.go` - 存储池管理

### 系统工具
- `backend/pkg/system/disk.go` - 磁盘系统操作

## API端点

### 磁盘管理
- `GET /api/storage/disks` - 获取磁盘列表
- `POST /api/storage/disks/format` - 格式化磁盘
- `POST /api/storage/mount` - 挂载磁盘
- `POST /api/storage/umount` - 卸载磁盘

### 分区管理
- `GET /api/storage/disks/:device/partitions` - 获取分区信息
- `POST /api/storage/disks/:device/partitions` - 创建分区
- `DELETE /api/storage/disks/:device/partitions/:number` - 删除分区

### SMART监控
- `GET /api/storage/disks/:device/smart` - 获取SMART信息
- `POST /api/storage/disks/:device/test` - 运行SMART测试
- `GET /api/storage/disks/:device/health` - 获取健康状态

### RAID管理
- `GET /api/storage/raid` - 获取RAID阵列
- `POST /api/storage/raid` - 创建RAID
- `POST /api/storage/raid/:name/add` - 添加磁盘到RAID

### LVM管理
- `GET /api/storage/lvm/pv` - 获取物理卷
- `POST /api/storage/lvm/pv` - 创建物理卷
- `GET /api/storage/lvm/vg` - 获取卷组
- `POST /api/storage/lvm/lv` - 创建逻辑卷

### 存储池管理
- `GET /api/storage/pools` - 获取存储池列表
- `POST /api/storage/pools` - 创建存储池
- `GET /api/storage/pools/:name` - 获取存储池详情
- `DELETE /api/storage/pools/:name` - 删除存储池

## 前端组件

- `frontend/src/apps/DiskManager.vue` - 磁盘管理界面
- `frontend/src/apps/StorageManager.vue` - 存储管理界面  
- `frontend/src/apps/StoragePoolManager.vue` - 存储池管理界面
- `frontend/src/components/Disk/` - 磁盘相关组件

## 配置文件

### 磁盘配置
- `config/disk.conf` - 磁盘默认配置
- `config/smart.conf` - SMART监控配置
- `config/raid.conf` - RAID配置模板
- `config/lvm.conf` - LVM配置

## 使用示例

### 创建RAID阵列
```bash
# 通过API创建RAID1
curl -X POST http://localhost:8888/api/storage/raid \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "data_raid",
    "level": 1,
    "devices": ["/dev/sda", "/dev/sdb"],
    "filesystem": "ext4"
  }'
```

### 创建存储池
```bash
# 通过API创建Btrfs存储池
curl -X POST http://localhost:8888/api/storage/pools \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "backup_pool",
    "disks": ["/dev/sdc", "/dev/sdd"],
    "raid_level": "raid1",
    "mount_point": "/mnt/backup"
  }'
```

## 故障排除

### 磁盘无法识别
1. 检查磁盘连接：`lsblk`
2. 查看系统日志：`dmesg | grep -i disk`
3. 检查磁盘状态：`sudo smartctl -a /dev/sdX`

### RAID创建失败
1. 确认mdadm工具已安装：`which mdadm`
2. 检查磁盘是否被使用：`df -h`
3. 查看RAID状态：`cat /proc/mdstat`

### LVM操作失败
1. 检查LVM工具：`which lvm`
2. 扫描物理卷：`sudo pvscan`
3. 查看卷组状态：`sudo vgdisplay`

## 性能优化

- 使用SSD作为缓存盘
- 配置适当的RAID级别
- 定期进行磁盘碎片整理
- 监控磁盘IO性能

## 安全建议

- 重要数据使用RAID1或RAID10
- 定期进行SMART监控
- 设置磁盘使用率警报
- 配置自动备份策略

## 相关文档

- [存储管理API文档](../../docs/API_DOCUMENTATION.md)
- [系统部署指南](../../docs/DEPLOYMENT_GUIDE.md)
- [故障排除指南](../../docs/TROUBLESHOOTING.md)