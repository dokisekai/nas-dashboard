#!/bin/bash
# 启动集成Immich的NAS Dashboard

echo "🚀 启动 NAS Dashboard (含Immich集成)..."
echo ""

# 检查前端目录
frontend_dir="/data/nas-dashboard/frontend"
cd "$frontend_dir" || exit 1

echo "📁 当前目录: $(pwd)"
echo "🔧 检查Node环境..."
if ! command -v npm &> /dev/null; then
    echo "❌ npm未安装，请先安装Node.js和npm"
    exit 1
fi

echo "✅ Node.js环境正常"

echo "📦 检查依赖..."
if [ ! -d "node_modules" ]; then
    echo "📦 正在安装依赖..."
    npm install
else
    echo "✅ 依赖已安装"
fi

echo ""
echo "🎨 启动开发服务器..."
echo "💡 前端开发服务器将在 http://localhost:3000 启动"
echo "📷 Immich服务已集成在主页中"
echo ""
echo "按 Ctrl+C 停止服务器"
echo ""

# 启动开发服务器
npm run dev
