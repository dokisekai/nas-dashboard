# 功耗监控模块

## 概述

功耗监控模块用于实时监控NAS系统的功耗数据，包括CPU功耗、磁盘功耗、整体系统功耗等信息。

## 功能特性

- 实时功耗数据采集
- 功耗历史数据记录
- 功耗统计和分析
- 功耗警报设置
- 可视化功耗图表

## 脚本说明

### 1. nas_power_monitor.sh
主要的功耗监控脚本，支持多种功耗数据采集方式。

```bash
# 使用示例
./nas_power_monitor.sh --interval 60 --output /var/log/power.log
```

### 2. get_power.sh
简单的功耗数据获取脚本。

```bash
# 使用示例
./get_power.sh
```

### 3. get_power_v2.sh
改进版功耗获取脚本，支持更多硬件平台。

```bash
# 使用示例  
./get_power_v2.sh --format json
```

### 4. power_manager.sh
功耗管理脚本，支持功耗限制和优化。

```bash
# 使用示例
./power_manager.sh --limit 100W
```

### 5. setup_power_monitor.sh
功耗监控安装配置脚本。

```bash
# 使用示例
./setup_power_monitor.sh --install-all
```

## 配置文件

功耗监控配置位于 `config/` 目录：

- power.conf - 主配置文件
- alerts.conf - 警报配置
- sensors.conf - 传感器配置

## API端点

后端提供以下功耗相关API：

- `GET /api/power/current` - 获取当前功耗
- `GET /api/power/history` - 获取历史功耗数据
- `GET /api/power/statistics` - 获取功耗统计信息
- `GET /api/power/overview` - 获取功耗概览

## 集成说明

1. 安装依赖：`./setup_power_monitor.sh --install-deps`
2. 配置传感器：编辑 `config/sensors.conf`
3. 启动监控：`./nas_power_monitor.sh --daemon`
4. 验证数据：访问 `/api/power/current`

## 故障排除

参考 `docs/POWER_INTEGRATION.md` 获取详细的集成和故障排除信息。

## 支持的硬件平台

- Intel平台 (RAPL)
- AMD平台  
- ARM平台
- 通用UPS设备

## 性能优化

- 默认采样间隔：60秒
- 数据保留时间：30天
- 内存占用：约50MB