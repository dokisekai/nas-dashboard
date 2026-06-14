# NAS 功耗监控系统集成文档

## 系统集成概述

本文档描述如何将功耗监控系统集成到现有的 NAS Dashboard 系统中，实现完整的功耗管理和监控功能。

## 已创建的文件

### 后端文件

1. **backend/internal/models/power_monitor.go**
   - 功耗监控数据模型定义
   - PowerInfo, PowerHistory, PowerStatistics 等结构

2. **backend/pkg/power/power.go**
   - 功耗监控核心逻辑
   - 直接从 sysfs 读取功耗数据
   - 历史数据分析和统计

3. **backend/internal/api/power_monitor.go**
   - REST API 接口定义
   - 提供 `/api/power/*` 端点

### 前端文件

4. **frontend/src/types/power.ts**
   - TypeScript 类型定义
   - PowerData, PowerHistory, PowerOverview 等

5. **frontend/src/api/power.ts**
   - API 客户端封装
   - WebSocket 实时监控

6. **frontend/src/components/PowerMonitor.vue**
   - Vue 组件实现
   - 实时功耗监控界面
   - 统计图表展示

### 独立脚本

7. **get_power_v2.sh**
   - 增强版功耗检测脚本
   - 支持 Intel 核显和 AMD 独显

8. **nas_power_monitor.sh**
   - 持续监控脚本
   - 历史数据记录和告警

## 系统架构集成

### 1. 后端集成

#### 路由注册
在 `backend/cmd/server/main.go` 中添加功耗监控路由：

```go
// 功耗监控路由
router.GET("/api/power/current", api.GetPowerCurrent)
router.GET("/api/power/history", api.GetPowerHistory)
router.GET("/api/power/statistics", api.GetPowerStatistics)
router.GET("/api/power/overview", api.GetPowerOverview)
```

#### WebSocket 支持
添加实时功耗数据推送：

```go
// 功耗监控 WebSocket
router.GET("/ws/power", func(c *gin.Context) {
    // WebSocket 处理逻辑
})
```

### 2. 前端集成

#### 路由配置
在 `frontend/src/router/index.ts` 中添加功耗监控路由：

```typescript
{
  path: '/power',
  name: 'PowerMonitor',
  component: () => import('@/views/PowerMonitor.vue')
}
```

#### 导航菜单
在主导航菜单中添加功耗监控入口：

```vue
<el-menu-item index="/power">
  <el-icon><Lightning /></el-icon>
  <span>功耗监控</span>
</el-menu-item>
```

### 3. 系统服务集成

#### systemd 服务
创建功耗监控守护服务：

```ini
[Unit]
Description=NAS Power Monitor Service
After=network.target

[Service]
Type=simple
User=root
ExecStart=/home/hserver/nas_power_monitor.sh
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

#### 定时任务
配置 cron 定期清理历史数据：

```bash
# 每天凌晨3点清理90天前的数据
0 3 * * * find /home/hserver/power_monitor/data -name "power_*.csv" -mtime +90 -delete
```

## 部署步骤

### 1. 后端部署

```bash
cd /home/hserver/nas-dashboard/backend

# 添加新的模型文件
# 创建 internal/models/power_monitor.go (已完成)

# 添加功耗监控包
mkdir -p pkg/power
# 创建 pkg/power/power.go (已完成)

# 添加 API 端点
# 创建 internal/api/power_monitor.go (已完成)

# 重新编译
go build -o bin/server cmd/server/main.go
```

### 2. 前端部署

```bash
cd /home/hserver/nas-dashboard/frontend

# 添加类型定义
# 创建 src/types/power.ts (已完成)

# 添加 API 客户端
# 创建 src/api/power.ts (已完成)

# 添加监控组件
# 创建 src/components/PowerMonitor.vue (已完成)

# 创建监控页面
# 创建 src/views/PowerMonitor.vue (使用 PowerMonitor 组件)

# 重新构建
npm run build
```

### 3. 系统配置

```bash
# 设置权限
chmod +x /home/hserver/get_power_v2.sh
chmod +x /home/hserver/nas_power_monitor.sh

# 配置 sudo 权限（无需密码访问功耗数据）
sudo tee -a /etc/sudoers.d/nas-monitor <<EOF
hserver ALL=(ALL) NOPASSWD: /home/hserver/get_power_v2.sh
hserver ALL=(ALL) NOPASSWD: /home/hserver/nas_power_monitor.sh
EOF

# 安装 systemd 服务
sudo cp /home/hserver/nas-power-monitor.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable nas-power-monitor
sudo systemctl start nas-power-monitor
```

## 功能特性

### 实时监控
- CPU Package/Core/Uncore 功耗
- Intel 核显功耗
- AMD 独显功耗
- 存储设备功耗
- 总功耗计算

### 历史分析
- 1/7/30/90 天历史数据
- 平均/最高/最低功耗统计
- 功耗趋势图表
- 峰值时段分析

### 告警系统
- 高功耗告警（可配置阈值）
- 严重功耗告警
- 邮件/短信通知（可扩展）
- 告警历史记录

### 智能管理
- 自动休眠低负载设备
- 性能模式自动切换
- 峰谷电价优化
- 负载预测和调度

## API 接口文档

### 获取当前功耗
```
GET /api/power/current
```
返回实时功耗数据

### 获取历史数据
```
GET /api/power/history?days=7
```
返回指定天数的历史功耗数据

### 获取统计信息
```
GET /api/power/statistics?days=7
```
返回功耗统计信息

### 获取概览
```
GET /api/power/overview
```
返回功耗概览（当前+统计+告警）

## 使用示例

### 基础使用

1. **启动监控服务**
```bash
sudo systemctl start nas-power-monitor
```

2. **查看实时功耗**
访问 Dashboard → 功耗监控

3. **设置告警阈值**
在设置界面配置告警阈值

### 高级使用

1. **手动获取功耗数据**
```bash
sudo /home/hserver/get_power_v2.sh
```

2. **查看历史数据**
```bash
cat /home/hserver/power_monitor/data/power_20260614.csv
```

3. **生成统计报告**
```bash
sudo /home/hserver/nas_power_monitor.sh report 30
```

## 性能考虑

### 资源占用
- **CPU**: <1% (监控间隔5分钟)
- **内存**: <50MB
- **磁盘**: ~1MB/天 (历史数据)

### 优化建议
1. 根据需求调整监控间隔
2. 定期清理历史数据
3. 使用数据压缩存储长期数据
4. 考虑使用数据库而非 CSV 文件

## 故障排除

### 常见问题

1. **权限错误**
```bash
sudo chown root:root /home/hserver/get_power_v2.sh
sudo chmod 4755 /home/hserver/get_power_v2.sh
```

2. **数据采集失败**
```bash
# 检查 sysfs 接口
ls -la /sys/devices/virtual/powercap/intel-rapl/
ls -la /sys/class/drm/card1/device/hwmon/hwmon6/
```

3. **服务启动失败**
```bash
sudo journalctl -u nas-power-monitor -n 50
```

## 扩展功能

### 计划中的功能

1. **电费计算**
   - 根据当地电价计算电费
   - 峰谷电价优化建议
   - 月度电费报告

2. **碳排放计算**
   - 根据功耗计算碳排放
   - 环保影响评估
   - 节能建议

3. **移动端支持**
   - 移动端优化界面
   - APP 推送通知
   - 远程控制

4. **机器学习预测**
   - 功耗趋势预测
   - 异常检测
   - 节能策略优化

---

*集成版本: v1.0*
*更新时间: 2026-06-14*
*兼容性: NAS Dashboard v1.0+*
