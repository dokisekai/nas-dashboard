#!/bin/bash
# Liquidctl 初始化脚本
# 用于系统启动时初始化 RGB 控制器

set -e

echo "正在初始化 Liquidctl 设备..."

# 检查 liquidctl 是否安装
if ! command -v liquidctl &> /dev/null; then
    echo "错误: liquidctl 未安装"
    exit 1
fi

# 初始化所有设备
echo "初始化所有设备..."
sudo liquidctl initialize all

# 设置默认配置（根据需要修改）
echo "设置默认灯光效果..."
sudo liquidctl set color mode fixed 0066ff  # 蓝色

echo "Liquidctl 初始化完成"
