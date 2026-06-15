#!/bin/bash

# 启动脚本示例
APP_DIR="/var/packages/hello-world/target"
APP_BIN="$APP_DIR/hello-world"
PID_FILE="$APP_DIR/hello-world.pid"
LOG_FILE="$APP_DIR/logs/hello-world.log"

case "$1" in
    start)
        if [ -f "$PID_FILE" ]; then
            PID=$(cat "$PID_FILE")
            if ps -p $PID > /dev/null 2>&1; then
                echo "应用已经在运行 (PID: $PID)"
                exit 0
            else
                rm -f "$PID_FILE"
            fi
        fi

        echo "启动 Hello World 应用..."
        cd "$APP_DIR"
        nohup "$APP_BIN" > "$LOG_FILE" 2>&1 &
        echo $! > "$PID_FILE"
        echo "应用启动成功 (PID: $(cat $PID_FILE))"
        ;;

    stop)
        if [ ! -f "$PID_FILE" ]; then
            echo "应用未运行"
            exit 0
        fi

        PID=$(cat "$PID_FILE")
        echo "停止 Hello World 应用 (PID: $PID)..."
        kill $PID
        rm -f "$PID_FILE"
        echo "应用已停止"
        ;;

    restart)
        $0 stop
        sleep 2
        $0 start
        ;;

    status)
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
        ;;

    *)
        echo "使用方法: $0 {start|stop|restart|status}"
        exit 1
        ;;
esac

exit 0