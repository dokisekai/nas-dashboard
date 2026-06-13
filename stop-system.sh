#!/bin/bash

# NAS Dashboard 系统停止脚本
# 用于停止所有服务

echo "🛑 停止 NAS Dashboard 系统"
echo "=========================="

# 检查 PID 文件
if [ -f "/tmp/nas-dashboard.pid" ]; then
    echo "📋 从 PID 文件读取进程..."
    PIDS=$(cat /tmp/nas-dashboard.pid)

    for PID in $PIDS; do
        if kill -0 $PID 2>/dev/null; then
            echo "🛑 停止进程 $PID..."
            kill $PID
            sleep 1

            # 如果进程仍在运行，强制终止
            if kill -0 $PID 2>/dev/null; then
                echo "⚠️  强制终止进程 $PID..."
                kill -9 $PID
            fi

            echo "✅ 进程 $PID 已停止"
        else
            echo "⚠️  进程 $PID 未运行"
        fi
    done

    rm /tmp/nas-dashboard.pid
else
    echo "⚠️  未找到 PID 文件，尝试查找进程..."

    # 查找并停止后端进程
    BACKEND_PID=$(pgrep -f "nas-dashboard")
    if [ -n "$BACKEND_PID" ]; then
        echo "🛑 停止后端服务 (PID: $BACKEND_PID)..."
        kill $BACKEND_PID
        sleep 1

        if kill -0 $BACKEND_PID 2>/dev/null; then
            kill -9 $BACKEND_PID
        fi

        echo "✅ 后端服务已停止"
    fi

    # 查找并停止前端进程
    FRONTEND_PID=$(pgrep -f "vite")
    if [ -n "$FRONTEND_PID" ]; then
        echo "🛑 停止前端服务 (PID: $FRONTEND_PID)..."
        kill $FRONTEND_PID
        sleep 1

        if kill -0 $FRONTEND_PID 2>/dev/null; then
            kill -9 $FRONTEND_PID
        fi

        echo "✅ 前端服务已停止"
    fi
fi

# 检查端口占用
echo ""
echo "📋 检查端口占用..."
if lsof -Pi :8888 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  端口 8888 仍被占用"
    lsof -Pi :8888 -sTCP:LISTEN
else
    echo "✅ 端口 8888 已释放"
fi

if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  端口 5173 仍被占用"
    lsof -Pi :5173 -sTCP:LISTEN
else
    echo "✅ 端口 5173 已释放"
fi

echo ""
echo "🎉 NAS Dashboard 系统已停止"
echo "=========================="
echo "💡 可以使用 ./start-system.sh 重新启动系统"
echo ""
