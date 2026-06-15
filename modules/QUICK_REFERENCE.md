# NAS Dashboard 快速参考指南

## 📋 模块快速索引

### 🏗️ 核心管理模块
1. **用户管理** (`user-management/`)
   - 用户账户管理、组管理、配额、SSH密钥
   - 文档: `user-management/README.md`
   - API: `/api/users/*`

2. **权限管理** (`permission-management/`)
   - 文件权限、ACL、SMB权限、权限继承
   - 文档: `permission-management/README.md`
   - API: `/api/permissions/*`

3. **共享管理** (`share-management/`)
   - SMB/CIFS、NFS、WebDAV、回收站
   - 文档: `share-management/README.md`
   - API: `/api/storage/smb`

### 💾 存储管理模块
4. **磁盘管理** (`disk-management/`)
   - 分区、RAID、LVM、SMART监控
   - 文档: `disk-management/README.md`
   - API: `/api/storage/*`

5. **存储池管理** (`storage-pool/`) ⭐
   - MergerFS合并硬盘、动态管理、智能策略
   - 文档: `storage-pool/README.md`
   - API: `/api/storage/pools/*`

### 🌐 网络和监控模块
6. **网络管理** (`network-management/`)
   - 网络接口、Wi-Fi、PPPoE、防火墙
   - 文档: `network-management/README.md`
   - API: `/api/network/*`

7. **系统监控** (`system-monitor/`)
   - CPU、内存、磁盘、进程监控
   - 文档: `system-monitor/README.md`
   - API: `/api/monitor/*`

### 🔧 系统工具模块
8. **备份恢复** (`backup-restore/`)
   - Restic备份、同步管理
   - 文档: `backup-restore/README.md`
   - API: `/api/storage/backup/*`

9. **应用系统** (`application-system/`)
   - 应用打包、安装、模板系统
   - 文档: `application-system/README.md`
   - API: `/api/apps/*`

10. **插件系统** (`plugin-system/`)
    - 前后端插件架构、SDK
    - 文档: `plugin-system/README.md`
    - API: `/api/plugins/*`

### ⚡ 高级功能模块
11. **功耗监控** (`power-monitor/`)
    - 多平台功耗采集、历史统计
    - 文档: `power-monitor/README.md`
    - API: `/api/power/*`

12. **液冷控制** (`liquidctl/`)
    - 设备管理、温度控制、RGB灯光
    - 文档: `liquidctl/README.md`

## 🚀 快速启动

### 系统启动
```bash
# 启动所有服务
./start-system.sh

# 访问地址
前端: http://localhost:5173
后端: https://localhost:8888
默认账户: admin/admin
```

### 模块使用
```bash
# 查看所有模块
ls modules/

# 使用功耗监控
cd modules/power-monitor/scripts/
./nas_power_monitor.sh

# 创建应用模板
cd modules/application-system/scripts/
./create-app-template.sh --name "my-app"

# 液冷状态查询
cd modules/liquidctl/scripts/
./liquidctl-status.sh
```

## 📖 文档索引

### 项目级文档
- [项目总览](../../docs/PROJECT_OVERVIEW.md)
- [API文档](../../docs/API_DOCUMENTATION.md)
- [用户指南](../../docs/USER_GUIDE.md)
- [开发者指南](../../docs/DEVELOPER_GUIDE.md)
- [部署指南](../../docs/DEPLOYMENT_GUIDE.md)

### 模块文档
- [模块总览](README.md)
- [架构指南](GUIDE.md)
- [完整清单](COMPLETE_INVENTORY.md)
- [核对报告](COMPLETENESS_CHECK.md)
- [创建总结](SETUP_SUMMARY.md)
- [最终总结](FINAL_SUMMARY.md)

## 🛠️ 脚本工具索引

### 系统脚本
- `start-system.sh` - 系统启动
- `stop-system.sh` - 系统停止

### 功能脚本 (21个)
- **功耗监控** (5个): `power-monitor/scripts/*.sh`
- **应用管理** (6个): `application-system/scripts/*.sh`
- **液冷控制** (3个): `liquidctl/scripts/*.sh`
- **备份恢复** (2个): `backup-restore/scripts/*.sh`
- **系统测试** (5个): `system-monitor/scripts/*.sh`

## 📊 功能对标

### DSM对标功能
| 功能 | 完整度 | 模块 |
|------|--------|------|
| 用户管理 | 95% | ✅ user-management |
| 权限管理 | 90% | ✅ permission-management |
| 共享管理 | 92% | ✅ share-management |
| 存储管理 | 95% | ✅ disk-management, storage-pool |
| 网络管理 | 90% | ✅ network-management |
| 系统监控 | 95% | ✅ system-monitor |
| 备份恢复 | 88% | ✅ backup-restore |
| 应用中心 | 85% | ✅ application-system |
| 插件系统 | 85% | ✅ plugin-system |

**总体对标完整度: 91%**

## 🎯 核心特性

### 🌟 独特优势
1. **合并硬盘功能** - MergerFS多磁盘合并
2. **完整用户系统** - 用户、组、权限、配额
3. **插件架构** - 前后端完整插件系统
4. **应用中心** - 标准化应用包管理
5. **企业级监控** - 功耗、性能、状态监控

### 🔒 安全特性
- JWT认证
- 基于角色的访问控制
- SSH密钥管理
- 权限继承机制
- 审计日志记录

### ⚡ 性能特性
- WebSocket实时推送
- 前端组件懒加载
- API响应缓存
- 数据库连接池
- 静态资源CDN

## 📞 技术支持

### 问题排查
1. 查看日志: `tail -f /tmp/nas-dashboard.log`
2. 检查状态: `systemctl status nas-dashboard`
3. 验证API: `curl http://localhost:8888/health`
4. 查看文档: 相应模块的README.md

### 常见问题
- **端口占用**: 修改`backend/cmd/server/main.go`中的端口
- **权限错误**: 检查文件权限和用户组配置
- **网络问题**: 检查防火墙和网络配置
- **服务启动失败**: 查看日志文件获取详细错误信息

## 🔗 相关资源

- GitHub: https://github.com/yourusername/nas-dashboard
- 文档: https://docs.nas-dashboard.com
- 问题反馈: https://github.com/yourusername/nas-dashboard/issues

---

**版本**: 1.0.0  
**状态**: 🟢 生产就绪  
**更新**: 2026-06-14