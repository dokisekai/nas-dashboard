#!/bin/bash

# =================================================================
# 智能算力呼吸脚本 (v1.0)
# 逻辑：根据 GPU 负载自动切换 CPU 调度模式，实现“算力与节能”动态平衡
# =================================================================

GPU_THRESHOLD=50  # 瓦特 (W)
CHECK_INTERVAL=60 # 秒

echo "智能算力管理服务已启动 (检查间隔: ${CHECK_INTERVAL}s)"

while true; do
    # 1. 获取 GPU 实时功耗 (从 hwmon 读取，不唤醒 rocm-smi)
    GPU_HWMON="/sys/class/drm/card1/device/hwmon/hwmon6/power1_average"
    GPU_W=0
    if [ -f "$GPU_HWMON" ]; then
        GPU_UW=$(cat "$GPU_HWMON" 2>/dev/null)
        if [ ! -z "$GPU_UW" ]; then
            GPU_W=$(echo "$GPU_UW / 1000000" | bc)
        fi
    fi

    # 2. 判断模式切换
    if (( GPU_W > GPU_THRESHOLD )); then
        # 算力怪兽模式 (High Performance)
        echo "$(date '+%H:%M:%S') [高性能] 检测到 GPU 负载 (${GPU_W}W)，激活算力怪兽模式"
        sudo cpupower frequency-set -g performance >/dev/null
        echo "performance" | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/energy_performance_preference >/dev/null
        echo "manual" | sudo tee /sys/class/drm/card1/device/power_dpm_force_performance_level >/dev/null
    else
        # 超级省电模式 (Super Power Saving)
        # 仅在当前不是 powersave 时输出，避免日志刷屏
        CURRENT_GOV=$(cat /sys/devices/system/cpu/cpu0/cpufreq/scaling_governor)
        if [[ "$CURRENT_GOV" != "powersave" ]]; then
            echo "$(date '+%H:%M:%S') [省电] 系统闲置中，进入超级省电模式"
        fi
        sudo cpupower frequency-set -g powersave >/dev/null
        echo "power" | sudo tee /sys/devices/system/cpu/cpu*/cpufreq/energy_performance_preference >/dev/null
        echo "low" | sudo tee /sys/class/drm/card1/device/power_dpm_force_performance_level >/dev/null
        
        # 顺便检查硬盘并尝试休眠 (设置了 -S 120 应该会自动，这里做个保险)
        # for disk in /dev/sd[a-z]; do [ -e "$disk" ] && sudo hdparm -y "$disk" >/dev/null 2>&1; done
    fi

    sleep $CHECK_INTERVAL
done
