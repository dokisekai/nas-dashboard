# 存储池管理模块

## 概述

存储池管理模块提供基于MergerFS的多磁盘存储池管理功能，允许用户将多个物理磁盘合并为一个大容量存储空间，并提供灵活的存储策略和高级管理功能。

## 功能特性

### 存储池创建
- 多磁盘合并管理
- 动态添加/移除磁盘
- 自动负载均衡
- 热插拔支持
- 存储策略配置

### MergerFS集成
- 多种文件系统支持
- 读写策略配置
- 空间管理策略
- 权限和配额控制
- 故障恢复机制

### 高级功能
- 存储池快照
- 存储池监控
- 性能优化
- 数据迁移工具
- 备份集成

## MergerFS核心概念

### 存储策略（Category）
- **create**: 最小空间优先
- **mv**: 同create策略
- **epall**: 轮询所有分支
- **epff**: 预留空间优先
- **eplfs**: 最小空闲空间优先
- **eplus**: 最早可用空间优先

### 分支模式（Branch Mode）
- **RW**: 读写模式
- **RO**: 只读模式
- **NC**: 无创建模式

## 后端API实现

### 核心文件
- `backend/pkg/mergerfs/manager.go` - MergerFS管理器
- `backend/internal/api/storage_pool.go` - 存储池API
- `backend/internal/models/storage_pool.go` - 存储池模型

### 前端组件
- `frontend/src/apps/StoragePoolManager.vue` - 存储池管理界面
- `frontend/src/components/StoragePool/PoolWizard.vue` - 存储池创建向导
- `frontend/src/types/storage_pool.ts` - 存储池类型定义
- `frontend/src/stores/storage_pool.ts` - 存储池状态管理

## API端点

### 存储池管理
- `GET /api/storage/pools` - 获取存储池列表
- `POST /api/storage/pools` - 创建存储池
- `GET /api/storage/pools/:name` - 获取存储池详情
- `PUT /api/storage/pools/:name` - 更新存储池配置
- `DELETE /api/storage/pools/:name` - 删除存储池

### 磁盘管理
- `POST /api/storage/pools/:name/disks` - 添加磁盘到存储池
- `DELETE /api/storage/pools/:name/disks/:device` - 从存储池移除磁盘

### 存储池操作
- `POST /api/storage/pools/:name/mount` - 挂载存储池
- `POST /api/storage/pools/:name/umount` - 卸载存储池
- `POST /api/storage/pools/:name/balance` - 平衡存储池数据
- `POST /api/storage/pools/:name/scan` - 扫描存储池状态

### 分支管理
- `GET /api/storage/pools/:name/branches` - 获取分支列表
- `PUT /api/storage/pools/:name/branches/:branch` - 更新分支配置

## 配置文件

### MergerFS配置
```json
{
  "branches": [
    {
      "path": "/mnt/disk1",
      "mode": "rw",
      "priority": 1
    },
    {
      "path": "/mnt/disk2", 
      "mode": "rw",
      "priority": 2
    }
  ],
  "category": "eplfs",
  "minfreespace": "4G",
  "direct_io": true,
  "async_read": true,
  "use_ino": true,
  "hard_remove": true,
  "auto_unshare": true
}
```

### 存储池配置
- `config/storage/pools.conf` - 存储池主配置
- `config/storage/mergerfs.conf` - MergerFS配置
- `config/storage/policies.conf` - 存储策略配置

## 使用示例

### 创建存储池
```bash
# 通过API创建存储池
curl -X POST http://localhost:8888/api/storage/pools \
  -H "Authorization: Bearer <token>" \
  -d '{
    "name": "main-pool",
    "mount_point": "/mnt/main-pool",
    "disks": ["/dev/sdb1", "/dev/sdc1"],
    "filesystem": "mergerfs",
    "config": {
      "category": "eplfs",
      "minfreespace": "4G"
    }
  }'
```

### 添加磁盘到存储池
```bash
curl -X POST http://localhost:8888/api/storage/pools/main-pool/disks \
  -H "Authorization: Bearer <token>" \
  -d '{
    "device": "/dev/sdd1",
    "mode": "rw"
  }'
```

### 平衡存储池
```bash
curl -X POST http://localhost:8888/api/storage/pools/main-pool/balance \
  -H "Authorization: Bearer <token>" \
  -d '{
    "strategy": "size",
    "threshold": "10%"
  }'
```

## MergerFS使用示例

### 手动挂载MergerFS
```bash
# 创建挂载点
mkdir -p /mnt/storage-pool

# 挂载MergerFS
mergerfs \
  -o cat=create=eplfs,epall \
  -o minfreespace=4G \
  -o direct_io \
  /mnt/disk1:/mnt/disk2:/mnt/disk3 \
  /mnt/storage-pool

# 验证挂载
df -h /mnt/storage-pool
```

### 配置文件挂载
```bash
# 使用配置文件
mergerfs /mnt/storage-pool -f \
  -o config=/etc/mergerfs/storage-pool.conf
```

### 系统集成
```bash
# 添加到/etc/fstab
/mnt/disk1:/mnt/disk2:/mnt/disk3  /mnt/storage-pool  mergerfs  cat=create=eplfs,epall,minfreespace=4G,direct_io  0  0
```

## 存储策略详解

### 创建策略
- **create**: 新文件创建位置
- **mv**: 移动文件时的策略
- **epall**: 轮询所有可用分支
- **epff**: 预留空间最多的分支
- **eplfs**: 最小空闲空间的分支
- **eplus**: 最早可用空间的分支

### 选择建议
- **大文件存储**: 使用eplfs，平衡空间使用
- **小文件存储**: 使用epall，分散IO负载
- **读写均衡**: 使用epff，预留空间优先
- **简单配置**: 使用create，基于最小空间

## 监控和维护

### 存储池监控
```bash
# 查看存储池状态
curl http://localhost:8888/api/storage/pools/main-pool

# 查看磁盘使用情况
df -h | grep storage-pool

# 查看MergerFS统计
mergerfs -f stat /mnt/storage-pool
```

### 存储池维护
- 定期检查磁盘健康
- 监控空间使用率
- 平衡数据分布
- 清理无用文件
- 验证数据完整性

### 性能优化
- 合理选择存储策略
- 配置适当的最小空间
- 启用异步读取
- 使用直接IO
- 平衡磁盘负载

## 故障排除

### 存储池无法挂载
1. 检查基础磁盘是否挂载：`df -h`
2. 验证磁盘权限：`ls -ld /mnt/disk*`
3. 检查MergerFS安装：`which mergerfs`
4. 查看系统日志：`dmesg | grep mergerfs`

### 磁盘无法添加到存储池
1. 验证磁盘格式：`blkid /dev/sdX`
2. 检查挂载状态：`mount | grep sdX`
3. 验证文件系统：`ls -la /mnt/diskX`
4. 检查磁盘健康：`sudo smartctl -a /dev/sdX`

### 性能问题
1. 检查存储策略配置
2. 验证磁盘IO性能：`iostat -x`
3. 检查网络带宽（如果是网络存储）
4. 调整缓存设置

### 数据迁移问题
1. 备份重要数据
2. 使用正确的迁移工具
3. 验证迁移结果
4. 更新配置文件

## 高级功能

### 存储池快照
```bash
# 创建快照
mv /mnt/storage-pool /mnt/storage-pool.snapshot
mkdir -p /mnt/storage-pool
# 恢复快照
umount /mnt/storage-pool
mv /mnt/storage-pool.snapshot /mnt/storage-pool
```

### 数据迁移
```bash
# 在线数据迁移
rsync -av --progress /mnt/storage-pool/ /mnt/new-pool/

# 离线数据迁移
umount /mnt/storage-pool
cp -r /mnt/disk1/* /mnt/new-disk1/
```

### 存储池扩容
```bash
# 添加新磁盘
curl -X POST http://localhost:8888/api/storage/pools/main-pool/disks \
  -H "Authorization: Bearer <token>" \
  -d '{"device": "/dev/sde1", "mode": "rw"}'

# 平衡数据
curl -X POST http://localhost:8888/api/storage/pools/main-pool/balance
```

## 安全建议

1. **权限管理**
   - 设置适当的挂载权限
   - 使用用户组隔离
   - 定期审查访问权限

2. **数据保护**
   - 定期备份重要数据
   - 使用冗余存储策略
   - 监控磁盘健康状态
   - 准备灾难恢复计划

3. **性能监控**
   - 监控IO性能
   - 跟踪空间使用
   - 设置告警阈值
   - 定期性能测试

## 性能基准

### 预期性能指标
- **顺序读取**: 500MB/s+ (取决于磁盘数量)
- **顺序写入**: 300MB/s+
- **随机读取**: 2000+ IOPS
- **随机写入**: 1000+ IOPS

### 性能优化建议
- 使用SSD作为缓存
- 合理设置存储策略
- 启用异步操作
- 配置适当的内存缓存

## 与其他模块集成

### 备份模块
- 存储池作为备份源
- 快照功能增强备份
- 存储池级别的备份策略

### 共享模块
- 基于存储池创建共享
- 配额管理集成
- 权限控制统一

### 监控模块
- 存储池状态监控
- 性能指标收集
- 告警规则配置

## 相关文档

- [系统完成总结](../../docs/SYSTEM_COMPLETION_SUMMARY.md)
- [API文档](../../docs/API_DOCUMENTATION.md)
- [MergerFS官方文档](https://github.com/trapexit/mergerfs)
- [存储管理API](../disk-management/)

## 最佳实践

1. **规划存储池**
   - 根据需求选择策略
   - 预留扩展空间
   - 考虑数据分类
   - 制定备份策略

2. **运维管理**
   - 定期监控状态
   - 及时处理告警
   - 优化数据分布
   - 维护磁盘健康

3. **容量规划**
   - 监控增长趋势
   - 预判容量需求
   - 提前准备扩容
   - 优化数据清理策略

4. **灾难恢复**
   - 定期测试恢复
   - 维护恢复文档
   - 培训操作人员
   - 更新应急计划