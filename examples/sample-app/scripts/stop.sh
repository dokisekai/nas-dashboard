#!/bin/bash

# 停止脚本 - 实际上这个功能已经包含在start.sh中了
# 这里只是为了完整性单独创建

APP_DIR="/var/packages/hello-world"
PID_FILE="$APP_DIR/hello-world.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "应用未运行"
    exit 0
fi

PID=$(cat "$PID_FILE")
echo "停止 Hello World 应用 (PID: $PID)..."

# 尝试优雅停止
kill $PID

# 等待进程结束
TIMEOUT=10
for i in $(seq 1 $TIMEOUT); do
    if ! ps -p $PID > /dev/null 2>&1; then
        echo "应用已停止"
        rm -f "$PID_FILE"
        exit 0
    fi
    sleep 1
done

# 如果优雅停止失败，强制停止
echo "强制停止应用..."
kill -9 $PID
rm -f "$PID_FILE"
echo "应用已强制停止"

exit 0