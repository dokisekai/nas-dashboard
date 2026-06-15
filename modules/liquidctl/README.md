# 液冷控制模块

## 概述

液冷控制模块提供对液冷设备的管理功能，包括设备初始化、温度监控、风扇控制和RGB灯光效果设置。

## 功能特性

### 设备管理
- 液冷设备识别
- 设备初始化配置
- 设备状态监控
- 设备固件更新

### 温度监控
- CPU温度监控
- 液冷液温度监控
- 温度曲线设置
- 温度警报设置

### 风扇控制
- 手动风扇速度控制
- 自动风扇曲线
- 风扇模式切换
- 风扇噪音控制

### RGB灯光
- 灯光颜色设置
- 灯光效果模式
- 灯光同步配置
- 自定义灯光预设

## 脚本说明

### 1. liquidctl-init.sh
液冷设备初始化脚本。

```bash
# 使用示例
sudo ./liquidctl-init.sh --init-all
```

### 2. liquidctl-status.sh
设备状态查询脚本。

```bash
# 使用示例
./liquidctl-status.sh --verbose
```

### 3. liquidctl-rgb-presets.sh
RGB灯光预设管理脚本。

```bash
# 使用示例
./liquidctl-rgb-presets.sh --preset rainbow
```

## 系统要求

### 硬件支持
- NZXT Kraken系列
- Corsair H系列
- EVGA CLC系列
- 其他兼容设备

### 软件依赖
- liquidctl工具
- Python 3.6+
- USB权限配置

## 安装配置

### 安装liquidctl
```bash
# Debian/Ubuntu
sudo apt install liquidctl

# 从源安装
pip install liquidctl

# USB权限配置
sudo usermod -a -G plugdev $USER
```

### 初始化设备
```bash
# 检测设备
sudo liquidctl list

# 初始化所有设备
sudo liquidctl initialize

# 设置初始配置
sudo liquidctl --match kraben set sync speed 50
```

## 使用示例

### 温度监控
```bash
# 获取设备温度
sudo liquidctl status

# 设置温度曲线
sudo liquidctl --match kraben set fan curve 30:30 40:50 60:100
```

### 风扇控制
```bash
# 设置固定风扇速度
sudo liquidctl --match kraben set fan speed 80

# 设置静音模式
sudo liquidctl --match kraben set fan quiet
```

### RGB灯光
```bash
# 设置固定颜色
sudo liquidctl --match kraben set led color ff0000

# 设置灯光效果
sudo liquidctl --match kraben set led breathing 00ff00 ff0000

# 关闭灯光
sudo liquidctl --match kraben set led off
```

## 配置文件

### 液冷配置
- `config/liquidctl/devices.conf` - 设备配置
- `config/liquidctl/fan.conf` - 风扇配置
- `config/liquidctl/rgb.conf` - RGB配置

### 温度曲线
- `config/liquidctl/temp-curves.conf` - 温度曲线配置
- `config/liquidctl/alerts.conf` - 温度警报配置

## API集成

虽然liquidctl是独立工具，但可以通过系统集成API进行控制：

### 温度监控集成
```javascript
// 通过系统监控API获取温度
fetch('/api/monitor/temperature')
  .then(res => res.json())
  .then(data => {
    // 自动调节风扇速度
    if (data.cpu_temp > 80) {
      adjustFanSpeed(100);
    }
  });
```

### 自定义控制脚本
```bash
#!/bin/bash
# 根据温度自动调节风扇
TEMP=$(sensors | grep -oP 'Core 0.*?\+\K[0-9.]+')
if (( $(echo "$TEMP > 70" | bc -l) )); then
    sudo liquidctl set fan speed 100
elif (( $(echo "$TEMP > 60" | bc -l) )); then
    sudo liquidctl set fan speed 70
else
    sudo liquidctl set fan speed 30
fi
```

## RGB灯光预设

### 内置预设
- **rainbow** - 彩虹循环
- **breathing** - 呼吸效果
- **fixed** - 固定颜色
- **spectrum** - 光谱循环
- **wave** - 波浪效果

### 自定义预设
```bash
# 创建自定义预设
liquidctl set led color --preset custom --colors ff0000,00ff00,0000ff

# 速度调整
liquidctl set led breathing --speed slow ff0000 00ff00
```

## 监控和日志

### 温度监控
```bash
# 持续监控
watch -n 5 'liquidctl status'

# 记录到日志
while true; do 
    echo "$(date): $(liquidctl status)" >> /var/log/liquidctl.log
    sleep 300
done
```

### 系统集成
```bash
# 与系统监控集成
./liquidctl-status.sh | jq '.cpu_temp' >> /var/log/temperature.log
```

## 故障排除

### 设备未识别
1. 检查USB连接：`lsusb | grep -i cooler`
2. 查看设备列表：`sudo liquidctl list`
3. 检查驱动加载：`dmesg | grep -i usb`

### 权限问题
1. 添加用户到组：`sudo usermod -a -G plugdev $USER`
2. 配置USB规则：`/etc/udev/rules.d/99-liquidctl.rules`
3. 重新登录生效

### 控制失败
1. 检查设备状态：`sudo liquidctl status`
2. 重新初始化：`sudo liquidctl initialize`
3. 更新liquidctl：`pip install --upgrade liquidctl`

### RGB灯光异常
1. 重置灯光：`sudo liquidctl set led off`
2. 重新设置：`sudo liquidctl set led color ff0000`
3. 检查设备固件

## 性能优化

- 根据温度动态调整风扇速度
- 设置合理的温度阈值
- 使用静音模式降低噪音
- 定期清洁液冷设备

## 安全建议

- 监控温度警报
- 设置最高温度限制
- 定期检查液冷液状态
- 备份设备配置

## 相关资源

- [liquidctl官方文档](https://github.com/liquidctl/liquidctl)
- [支持的设备列表](https://github.com/liquidctl/liquidctl#supported-devices)
- [RGB灯光效果示例](https://github.com/liquidctl/liquidctl#rgb-commands)