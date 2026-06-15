#!/bin/bash
# Liquidctl 状态监控脚本
# 显示设备状态和系统信息

echo "=== Liquidctl 设备状态 ==="
echo ""

# 显示 liquidctl 版本
echo "版本信息:"
liquidctl --version
echo ""

# 列出所有设备
echo "检测到的设备:"
sudo liquidctl list -v
echo ""

# 显示当前状态
echo "当前状态:"
sudo liquidctl status
echo ""

# 显示设备权限
echo "设备权限:"
ls -la /dev/hidraw* 2>/dev/null || echo "未找到 hidraw 设备"
echo ""

# 检查 udev 规则
echo "Udev 规则:"
if [ -f /etc/udev/rules.d/99-liquidctl.rules ]; then
    echo "✓ 已配置 liquidctl udev 规则"
    cat /etc/udev/rules.d/99-liquidctl.rules
else
    echo "✗ 未找到 liquidctl udev 规则"
    echo "建议运行以下命令配置:"
    echo 'sudo tee /etc/udev/rules.d/99-liquidctl.rules > /dev/null <<EOF'
    echo 'SUBSYSTEM=="hidraw", ATTRS{idVendor}=="048d", ATTRS{idProduct}=="5702", MODE="0666"'
    echo "EOF"
fi
