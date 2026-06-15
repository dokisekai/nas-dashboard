#!/bin/bash
# Liquidctl RGB 预设效果脚本
# 提供常用的 RGB 灯光效果预设

# 颜色定义
RED="ff0000"
GREEN="00ff00"
BLUE="0000ff"
YELLOW="ffff00"
CYAN="00ffff"
MAGENTA="ff00ff"
WHITE="ffffff"
ORANGE="ff8000"
PURPLE="8000ff"

show_menu() {
    echo "=== Liquidctl RGB 效果预设 ==="
    echo "1. 红色固定"
    echo "2. 绿色固定"
    echo "3. 蓝色固定"
    echo "4. 白色固定"
    echo "5. 红绿呼吸"
    echo "6. 蓝紫交替"
    echo "7. 彩虹光谱"
    echo "8. 彩色波浪"
    echo "9. 彩色闪烁"
    echo "10. 关闭灯光"
    echo "0. 退出"
    echo -n "请选择效果 (0-10): "
}

set_color() {
    local mode=$1
    shift
    sudo liquidctl set color mode "$mode" "$@"
    echo "已设置: $mode $@"
}

case ${1:-} in
    1|"red")
        set_color fixed $RED
        ;;
    2|"green")
        set_color fixed $GREEN
        ;;
    3|"blue")
        set_color fixed $BLUE
        ;;
    4|"white")
        set_color fixed $WHITE
        ;;
    5|"breathe")
        set_color breathing $RED $GREEN
        ;;
    6|"alternating")
        set_color alternating $BLUE $PURPLE
        ;;
    7|"spectrum")
        set_color spectrum
        ;;
    8|"wave")
        set_color wave
        ;;
    9|"flashing")
        set_color flashing $RED $GREEN $BLUE
        ;;
    10|"off")
        set_color off
        ;;
    *)
        show_menu
        read -r choice
        $0 "$choice"
        ;;
esac
