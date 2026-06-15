#!/bin/bash

# =================================================================
# NAS 功耗持续监控系统 (v1.0)
# 功能：24/7 功耗监控、历史记录、异常检测、智能管理
# =================================================================

# 配置参数
MONITOR_DIR="/home/hserver/power_monitor"
LOG_DIR="$MONITOR_DIR/logs"
DATA_DIR="$MONITOR_DIR/data"
ALERT_THRESHOLD_HIGH=150    # 高功耗告警阈值 (W)
ALERT_THRESHOLD_CRITICAL=200 # 严重功耗告警 (W)
MONITOR_INTERVAL=300        # 监控间隔 (秒)
RETENTION_DAYS=90           # 数据保留天数

# 创建必要目录
mkdir -p "$LOG_DIR" "$DATA_DIR"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 日志函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_DIR/monitor.log"
}

# 获取功耗数据
get_power_data() {
    local output_file="$1"

    # CPU RAPL 数据
    local cpu_pkg_w=0 cpu_core_w=0 cpu_uncore_w=0
    if [ -f "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/energy_uj" ]; then
        read e1 < "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/energy_uj"
        read c1 < "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/intel-rapl:0:0/energy_uj"
        read u1 < "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/intel-rapl:0:1/energy_uj"
        sleep 1
        read e2 < "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/energy_uj"
        read c2 < "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/intel-rapl:0:0/energy_uj"
        read u2 < "/sys/devices/virtual/powercap/intel-rapl/intel-rapl:0/intel-rapl:0:1/energy_uj"

        cpu_pkg_w=$(echo "scale=2; ($e2 - $e1) / 1000000" | bc)
        cpu_core_w=$(echo "scale=2; ($c2 - $c1) / 1000000" | bc)
        cpu_uncore_w=$(echo "scale=2; ($u2 - $u1) / 1000000" | bc)
    fi

    # Intel 核显功耗
    local igpu_w=0.5
    if [ $(echo "$cpu_uncore_w > 0" | bc -l) -eq 1 ]; then
        igpu_w=$(echo "scale=2; $cpu_uncore_w * 0.7" | bc)
    fi

    # AMD 独显功耗
    local dgpu_w=5.5
    local gpu_hwmon="/sys/class/drm/card1/device/hwmon/hwmon6/power1_average"
    if [ -f "$gpu_hwmon" ]; then
        local gpu_uw=$(cat "$gpu_hwmon" 2>/dev/null)
        [ ! -z "$gpu_uw" ] && dgpu_w=$(echo "scale=2; $gpu_uw / 1000000" | bc)
        if (( $(echo "$dgpu_w < 5" | bc -l) )); then dgpu_w=5.5; fi
    fi

    # 硬盘功耗
    local hdd_w=0
    for dev in sda sdb; do
        if [ -e "/dev/$dev" ]; then
            local state=$(hdparm -C "/dev/$dev" 2>/dev/null | grep -o 'active/idle\|standby')
            if [[ "$state" == "standby" ]]; then
                hdd_w=$(echo "$hdd_w + 0.9" | bc)
            else
                hdd_w=$(echo "$hdd_w + 8.5" | bc)
            fi
        fi
    done

    # 其他组件功耗 (预估)
    local ssd_w=4.0
    local mb_ram_w=21.0
    local cool_w=12.0
    local usb_w=2.0

    # 计算总功耗
    local raw_total=$(echo "$cpu_pkg_w + $dgpu_w + $hdd_w + $ssd_w + $mb_ram_w + $cool_w + $usb_w" | bc)
    local loss=$(echo "scale=2; $raw_total * 0.18" | bc)
    local total_w=$(echo "$raw_total + $loss" | bc)

    # 输出到文件
    cat > "$output_file" << EOF
timestamp=$(date '+%Y-%m-%d %H:%M:%S')
cpu_pkg_w=$cpu_pkg_w
cpu_core_w=$cpu_core_w
cpu_uncore_w=$cpu_uncore_w
igpu_w=$igpu_w
dgpu_w=$dgpu_w
hdd_w=$hdd_w
ssd_w=$ssd_w
mb_ram_w=$mb_ram_w
cool_w=$cool_w
usb_w=$usb_w
total_w=$total_w
EOF
}

# 检查功耗异常
check_power_alert() {
    local total_w=$1

    if (( $(echo "$total_w >= $ALERT_THRESHOLD_CRITICAL" | bc -l) )); then
        log "🚨 严重告警：功耗过高 ($total_w W >= $ALERT_THRESHOLD_CRITICAL W)"
        # 发送通知 (可扩展)
        return 2
    elif (( $(echo "$total_w >= $ALERT_THRESHOLD_HIGH" | bc -l) )); then
        log "⚠️ 高功耗告警：功耗较高 ($total_w W >= $ALERT_THRESHOLD_HIGH W)"
        return 1
    else
        return 0
    fi
}

# 持续监控函数
monitor_loop() {
    log "🚀 NAS 功耗监控系统启动"
    log "监控间隔: ${MONITOR_INTERVAL}s, 数据保留: ${RETENTION_DAYS}天"

    while true; do
        local temp_file="/tmp/power_data_$$.tmp"
        get_power_data "$temp_file"

        # 读取数据
        source "$temp_file"

        # 记录到CSV
        local csv_file="$DATA_DIR/power_$(date '+%Y%m%d').csv"
        if [ ! -f "$csv_file" ]; then
            echo "timestamp,cpu_pkg,cpu_core,cpu_uncore,igpu,dgpu,hdd,ssd,mb_ram,cool,usb,total" > "$csv_file"
        fi
        echo "$timestamp,$cpu_pkg_w,$cpu_core_w,$cpu_uncore_w,$igpu_w,$dgpu_w,$hdd_w,$ssd_w,$mb_ram_w,$cool_w,$usb_w,$total_w" >> "$csv_file"

        # 检查告警
        check_power_alert "$total_w"

        # 清理临时文件
        rm -f "$temp_file"

        # 清理旧数据
        find "$DATA_DIR" -name "power_*.csv" -mtime +$RETENTION_DAYS -delete

        sleep $MONITOR_INTERVAL
    done
}

# 生成统计报告
generate_report() {
    local days=${1:-7}
    local end_date=$(date '+%Y%m%d')
    local start_date=$(date -d "$days days ago" '+%Y%m%d')

    echo "=== NAS 功耗统计报告 ($start_date - $end_date) ==="
    echo ""

    # 汇总所有相关文件
    local total_records=0
    local sum_power=0
    local max_power=0
    local min_power=9999

    for ((i=0; i<=days; i++)); do
        local file_date=$(date -d "$i days ago" '+%Y%m%d')
        local csv_file="$DATA_DIR/power_${file_date}.csv"

        if [ -f "$csv_file" ]; then
            while IFS=',' read -r timestamp cpu_pkg cpu_core cpu_uncore igpu dgpu hdd ssd mb_ram cool usb total; do
                if [[ "$total" =~ ^[0-9]+\.?[0-9]*$ ]]; then
                    total_records=$((total_records + 1))
                    sum_power=$(echo "$sum_power + $total" | bc)

                    if (( $(echo "$total > $max_power" | bc -l) )); then
                        max_power=$total
                    fi

                    if (( $(echo "$total < $min_power" | bc -l) )); then
                        min_power=$total
                    fi
                fi
            done < <(tail -n +2 "$csv_file")
        fi
    done

    if [ $total_records -gt 0 ]; then
        local avg_power=$(echo "scale=2; $sum_power / $total_records" | bc)
        echo "记录数量: $total_records"
        echo "平均功耗: $avg_power W"
        echo "最高功耗: $max_power W"
        echo "最低功耗: $min_power W"
        echo "总耗电量: $(echo "scale=2; ($sum_power / 1000) * ($total_records * $MONITOR_INTERVAL / 3600)" | bc) kWh"
    else
        echo "暂无数据"
    fi
}

# 主程序
case "${1:-start}" in
    start)
        monitor_loop
        ;;
    report)
        generate_report "${2:-7}"
        ;;
    status)
        echo "=== 监控状态 ==="
        echo "日志目录: $LOG_DIR"
        echo "数据目录: $DATA_DIR"
        echo "最近数据:"
        ls -lt "$DATA_DIR/" | head -5
        ;;
    *)
        echo "用法: $0 {start|report|status}"
        echo "  start   - 启动监控"
        echo "  report  - 生成统计报告 [天数]"
        echo "  status  - 查看监控状态"
        exit 1
        ;;
esac
