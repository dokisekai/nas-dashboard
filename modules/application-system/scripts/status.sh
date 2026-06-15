#!/bin/bash

# 状态检查脚本
APP_DIR="/var/packages/hello-world"
PID_FILE="$APP_DIR/hello-world.pid"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p $PID > /dev/null 2>&1; then
        echo "running"
        exit 0
    else
        echo "stopped"
        exit 1
    fi
else
    echo "stopped"
    exit 1
fi