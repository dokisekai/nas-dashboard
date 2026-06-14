# 控制面板合并说明

## 合并概述

已将系统设置 (SystemSettings) 的所有功能成功合并到控制面板 (ControlPanel) 中，实现了统一的设置管理界面。

## 新增分类

### 1. 系统信息 (System Info)
- **图标**: ServerIcon
- **描述**: 查看系统硬件和运行信息
- **访问级别**: 用户
- **特性**: 只读分类，显示系统硬件信息

**包含设置**:
- 主机名 (system.hostname.info)
- CPU信息 (system.cpu)
- 内存信息 (system.memory)
- 存储信息 (system.storage)
- 操作系统 (system.os)
- 运行时间 (system.uptime)
- 内核版本 (system.kernel) - 高级设置
- 系统架构 (system.architecture) - 高级设置

### 2. 服务管理 (Services)
- **图标**: CogIcon
- **描述**: 管理系统服务和守护进程
- **访问级别**: 管理员
- **特性**: 包含自定义服务列表组件

**包含设置**:
- 系统服务列表 (services.list) - 自定义组件
- 自动启动关键服务 (services.autoStart)
- 服务监控 (services.monitoring) - 高级设置

### 3. 更新管理 (Updates)
- **图标**: BeakerIcon
- **描述**: 系统更新和版本管理
- **访问级别**: 管理员

**包含设置**:
- 当前版本 (updates.currentVersion) - 只读
- 自动检查更新 (updates.autoCheck)
- 自动安装安全更新 (updates.autoInstall) - 高级设置
- 检查频率 (updates.frequency)
- 更新通知 (updates.notify)
- 维护时间窗口 (updates.maintenanceWindow) - 高级设置

### 4. 备份管理 (Backup)
- **图标**: CircleStackIcon
- **描述**: 系统备份和恢复
- **访问级别**: 管理员

**包含设置**:
- 启用自动备份 (backup.enabled)
- 备份频率 (backup.schedule)
- 备份时间 (backup.time)
- 备份位置 (backup.location)
- 保留天数 (backup.retention) - 高级设置
- 压缩备份 (backup.compression) - 高级设置
- 加密备份 (backup.encrypt) - 高级设置

## 技术实现

### 1. 新增类型支持

**ControlPanelTypes.ts**:
```typescript
// 新增类型
export type SettingType =
  | 'boolean'
  | 'string'
  | 'number'
  | 'select'
  | 'multiselect'
  | 'slider'
  | 'color'
  | 'textarea'
  | 'password'
  | 'file'
  | 'group'
  | 'custom'
  | 'readonly' // 新增

// 新增属性
export interface ControlPanelSetting {
  // ... 现有属性
  readonly?: boolean    // 新增：只读标识
  component?: string    // 新增：自定义组件名
}

export interface ControlPanelCategory {
  // ... 现有属性
  readonly?: boolean    // 新增：分类只读标识
}
```

### 2. 新增组件

**ServicesListComponent.vue**:
- 系统服务列表展示和管理
- 支持服务搜索
- 服务启停控制
- 服务启用/禁用管理
- 服务重启功能

### 3. 控制面板增强

**ControlPanelSettingComponent.vue**:
- 添加只读类型支持 (`readonly`)
- 添加自定义组件类型支持 (`custom`)
- 新增 `getCustomComponent()` 方法
- 新增 `handleCustomUpdate()` 方法
- 新增只读样式 `.cps-readonly`

**ControlPanel.vue**:
- 添加新图标支持 (BeakerIcon, CircleStackIcon)
- 更新图标映射

### 4. 存储增强 (controlPanel.ts)

在 `defaultCategories` 数组中新增4个分类:
- 系统信息 (order: 7)
- 服务管理 (order: 8)
- 更新管理 (order: 9)
- 备份管理 (order: 10)

## 使用方式

### 访问新功能

1. **通过控制面板访问**:
   - 打开控制面板应用
   - 在侧边栏中选择相应分类:
     - "系统信息" - 查看硬件和系统信息
     - "服务管理" - 管理系统服务
     - "更新管理" - 配置系统更新
     - "备份管理" - 设置自动备份

2. **搜索功能**:
   - 使用控制面板的搜索功能
   - 搜索相关设置，如"服务"、"备份"、"更新"等

### 设置特性

#### 只读设置
- 系统信息分类中的所有设置都是只读的
- 显示为灰色背景，无法编辑
- 用于展示系统当前状态

#### 自定义组件
- 服务列表使用自定义组件 `ServicesListComponent`
- 支持复杂交互逻辑
- 与设置系统完全集成

#### 高级设置
- 各分类中的高级设置默认隐藏
- 点击"显示"按钮可查看
- 需要谨慎修改的高级配置

## 迁移说明

### 原SystemSettings功能对照

| 原SystemSettings Tab | 新ControlPanel 分类 | 功能完整性 |
|----------------------|---------------------|-----------|
| 网络 (network) | 网络设置 (network) | ✅ 已存在 |
| 系统信息 (info) | 系统信息 (system-info) | ✅ 已迁移 |
| 服务 (services) | 服务管理 (services) | ✅ 已迁移 |
| 更新 (updates) | 更新管理 (updates) | ✅ 已迁移 |
| 备份 (backup) | 备份管理 (backup) | ✅ 已迁移 |

### API集成

**建议集成点**:
```typescript
// 系统信息获取
const fetchSystemInfo = async () => {
  const response = await systemApi.getSystemInfo()
  // 更新 system-info 分类的设置值
}

// 服务列表获取
const fetchServices = async () => {
  const response = await serviceApi.list()
  // 更新 services.list 设置值
}

// 备份管理
const createBackup = async () => {
  const response = await backupApi.create()
  // 执行备份操作
}
```

## 优势

### 1. 统一体验
- 所有系统设置集中在一个界面
- 一致的操作模式和交互方式
- 统一的搜索和管理功能

### 2. 增强功能
- 支持设置导入导出
- 未保存更改提醒
- 高级设置保护机制
- 更好的权限控制

### 3. 可扩展性
- 模块化的分类系统
- 支持自定义组件集成
- 完整的验证机制
- 依赖关系管理

### 4. 用户友好
- 直观的分类导航
- 智能搜索功能
- 清晰的设置描述
- 及时的操作反馈

## 后续工作

### 集成真实API

需要将当前的模拟数据替换为真实API调用:

1. **系统信息API**:
   ```typescript
   GET /api/system/info
   ```

2. **服务管理API**:
   ```typescript
   GET /api/services
   POST /api/services/:name/start
   POST /api/services/:name/stop
   POST /api/services/:name/restart
   ```

3. **更新管理API**:
   ```typescript
   GET /api/updates
   POST /api/updates/check
   POST /api/updates/install
   ```

4. **备份管理API**:
   ```typescript
   GET /api/backups
   POST /api/backups/create
   DELETE /api/backups/:id
   ```

### 功能完善

1. **系统信息实时更新**:
   - 定期刷新系统信息
   - WebSocket推送状态变更

2. **服务状态监控**:
   - 实时服务状态显示
   - 服务失败告警

3. **备份管理增强**:
   - 备份文件列表显示
   - 恢复功能实现
   - 备份进度显示

4. **更新管理完善**:
   - 可用更新详情显示
   - 更新下载进度
   - 更新日志查看

## 兼容性

- ✅ 向后兼容：现有控制面板功能完全保留
- ✅ 渐进迁移：用户可以逐步适应新界面
- ✅ 功能完整：所有SystemSettings功能都已迁移
- ✅ 数据一致：设置数据格式保持兼容

## 总结

系统设置功能已成功合并到控制面板，提供了更统一、更强大的设置管理体验。新界面保持了原有的所有功能，同时引入了现代化的设置管理特性，为用户提供更好的系统配置体验。

---
**合并完成日期**: 2026-06-12
**影响范围**: ControlPanel.vue, controlPanel.ts, 相关组件和类型定义
**向后兼容**: 是
**需要用户适应**: 否 (功能相同，界面优化)
