#!/bin/bash

# NAS Dashboard 系统启动脚本
# 用于快速启动开发环境和测试系统

echo "🚀 NAS Dashboard 系统启动"
echo "=========================="

# 检查系统依赖
echo "📋 检查系统依赖..."

# 检查 Go
if ! command -v go &> /dev/null; then
    echo "❌ Go 未安装，请先安装 Go 1.19+"
    exit 1
fi

# 检查 Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js 未安装，请先安装 Node.js 18+"
    exit 1
fi

# 检查 PostgreSQL
if ! command -v psql &> /dev/null; then
    echo "⚠️  PostgreSQL 未安装，部分功能可能受限"
fi

echo "✅ 系统依赖检查完成"

# 启动后端服务
echo ""
echo "🔧 启动后端服务..."
cd backend

# 检查端口是否被占用
if lsof -Pi :8888 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  端口 8888 已被占用，尝试停止现有服务..."
    pkill -f "nas-dashboard" || true
    sleep 2
fi

# 编译并启动后端
echo "📦 编译后端..."
go build -o nas-dashboard cmd/server/main.go

if [ $? -eq 0 ]; then
    echo "✅ 后端编译成功"
    echo "🚀 启动后端服务 (端口 8888)..."
    ./nas-dashboard > /tmp/nas-dashboard.log 2>&1 &

    # 等待服务启动
    sleep 3

    # 检查服务是否启动成功
    if curl -s http://localhost:8888/api/health > /dev/null 2>&1; then
        echo "✅ 后端服务启动成功"
    else
        echo "❌ 后端服务启动失败，请检查日志: /tmp/nas-dashboard.log"
        exit 1
    fi
else
    echo "❌ 后端编译失败"
    exit 1
fi

cd ..

# 启动前端服务
echo ""
echo "🎨 启动前端服务..."
cd frontend

# 检查端口是否被占用
if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "⚠️  端口 5173 已被占用，尝试停止现有服务..."
    pkill -f "vite" || true
    sleep 2
fi

# 检查依赖是否安装
if [ ! -d "node_modules" ]; then
    echo "📦 安装前端依赖..."
    npm install
fi

# 启动开发服务器
echo "🚀 启动前端开发服务器 (端口 5173)..."
npm run dev > /tmp/nas-frontend.log 2>&1 &
FRONTEND_PID=$!

# 等待服务启动
sleep 3

# 检查服务是否启动成功
if curl -s http://localhost:5173 > /dev/null 2>&1; then
    echo "✅ 前端服务启动成功"
else
    echo "❌ 前端服务启动失败，请检查日志: /tmp/nas-frontend.log"
    exit 1
fi

cd ..

echo ""
echo "🎉 NAS Dashboard 系统启动成功！"
echo "=========================="
echo "📱 前端地址: http://localhost:5173"
echo "🔌 后端地址: http://localhost:8888"
echo "👤 默认账户: admin / admin"
echo ""
echo "📋 系统状态:"
echo "  - 后端 PID: $(pgrep -f nas-dashboard)"
echo "  - 前端 PID: $FRONTEND_PID"
echo ""
echo "📝 日志文件:"
echo "  - 后端日志: /tmp/nas-dashboard.log"
echo "  - 前端日志: /tmp/nas-frontend.log"
echo ""
echo "🛑 停止系统:"
echo "  ./stop-system.sh"
echo ""
echo "🔍 查看日志:"
echo "  tail -f /tmp/nas-dashboard.log   # 后端日志"
echo "  tail -f /tmp/nas-frontend.log    # 前端日志"
echo ""

# 保存 PID 到文件
echo "$(pgrep -f nas-dashboard)" > /tmp/nas-dashboard.pid
echo "$FRONTEND_PID" >> /tmp/nas-dashboard.pid

echo "💡 提示: 首次登录后请立即修改默认密码"
echo "📚 文档: 请查看项目目录下的文档文件"
echo ""
