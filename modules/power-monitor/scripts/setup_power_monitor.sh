#!/bin/bash

# NAS 功耗监控系统一键安装脚本

set -e

echo "=== NAS 功耗监控系统安装 ==="

# 检查 root 权限
if [ "$EUID" -ne 0 ]; then
    echo "请使用 root 权限运行此脚本"
    exit 1
fi

# 创建必要目录
echo "创建目录结构..."
mkdir -p /home/hserver/power_monitor/logs
mkdir -p /home/hserver/power_monitor/data
mkdir -p /home/hserver/power_monitor/reports

# 设置权限
echo "设置权限..."
chmod +x /home/hserver/get_power_v2.sh
chmod +x /home/hserver/nas_power_monitor.sh
chown -R root:root /home/hserver/power_monitor

# 配置 sudo 权限
echo "配置 sudo 权限..."
tee /etc/sudoers.d/nas-monitor <<EOF
hserver ALL=(ALL) NOPASSWD: /home/hserver/get_power_v2.sh
hserver ALL=(ALL) NOPASSWD: /home/hserver/nas_power_monitor.sh
EOF
chmod 440 /etc/sudoers.d/nas-monitor

# 安装 systemd 服务
echo "安装 systemd 服务..."
cp /home/hserver/nas-power-monitor.service /etc/systemd/system/
systemctl daemon-reload
systemctl enable nas-power-monitor

# 配置 udev 规则（GPU 访问）
echo "配置 udev 规则..."
tee /etc/udev/rules.d/99-gpu-power.rules <<EOF
# AMD GPU hwmon access
SUBSYSTEM=="hidraw", ATTRS{idVendor}=="048d", ATTRS{idProduct}=="5702", MODE="0666"

# Intel RAPL access (already accessible)
KERNEL=="intel_rapl", MODE="0666"
EOF
udevadm control --reload-rules
udevadm trigger

# 配置定时清理任务
echo "配置定时清理任务..."
tee /etc/cron.d/nas-power-monitor <<EOF
# NAS 功耗监控定时任务
# 每天凌晨3点清理90天前的数据
0 3 * * * find /home/hserver/power_monitor/data -name "power_*.csv" -mtime +90 -delete

# 每小时生成临时统计报告
0 * * * * /home/hserver/nas_power_monitor.sh report 1 > /home/hserver/power_monitor/reports/hourly_report.txt

# 每天凌晨4点生成每日报告
0 4 * * * /home/hserver/nas_power_monitor.sh report 7 > /home/hserver/power_monitor/reports/daily_report.txt
EOF
chmod 644 /etc/cron.d/nas-power-monitor

# 测试功耗监控
echo "测试功耗监控..."
if sudo -u hserver /home/hserver/get_power_v2.sh > /dev/null 2>&1; then
    echo "✓ 功耗检测脚本工作正常"
else
    echo "✗ 功耗检测脚本测试失败"
    exit 1
fi

# 启动服务
echo "启动功耗监控服务..."
systemctl start nas-power-monitor

# 检查服务状态
sleep 2
if systemctl is-active --quiet nas-power-monitor; then
    echo "✓ 功耗监控服务启动成功"
else
    echo "✗ 功耗监控服务启动失败"
    journalctl -u nas-power-monitor -n 20
    exit 1
fi

echo ""
echo "=== 安装完成 ==="
echo ""
echo "服务管理命令:"
echo "  启动服务: sudo systemctl start nas-power-monitor"
echo "  停止服务: sudo systemctl stop nas-power-monitor"
echo "  重启服务: sudo systemctl restart nas-power-monitor"
echo "  查看状态: sudo systemctl status nas-power-monitor"
echo "  查看日志: sudo journalctl -u nas-power-monitor -f"
echo ""
echo "数据位置:"
echo "  日志文件: /home/hserver/power_monitor/logs/"
echo "  历史数据: /home/hserver/power_monitor/data/"
echo "  统计报告: /home/hserver/power_monitor/reports/"
echo ""
echo "手动测试:"
echo "  获取当前功耗: sudo /home/hserver/get_power_v2.sh"
echo "  生成统计报告: sudo /home/hserver/nas_power_monitor.sh report 7"
echo ""
echo "Dashboard 集成:"
echo "  访问: http://your-server-ip/power"
echo ""
echo "安装完成后，建议重启系统以确保所有配置生效。"
