#!/bin/bash
# 测试Immich桌面集成

echo "🧪 测试Immich桌面集成"
echo "================================"
echo ""

# 检查前端目录
frontend_dir="/data/nas-dashboard/frontend"
cd "$frontend_dir" || exit 1

echo "📁 当前目录: $(pwd)"
echo ""

# 检查前端依赖
if [ ! -d "node_modules" ]; then
    echo "📦 正在安装前端依赖..."
    npm install
    echo "✅ 依赖安装完成"
else
    echo "✅ 前端依赖已存在"
fi
echo ""

# 检查Immich服务
echo "🔍 检查Immich服务状态..."
if curl -s http://localhost:2283 > /dev/null 2>&1; then
    echo "✅ Immich服务正在运行"
else
    echo "❌ Immich服务未运行"
    echo "请先启动Immich服务："
    echo "  docker start immich"
    echo "  或使用Docker管理器启动"
    exit 1
fi
echo ""

# 检查后端服务
echo "🔍 检查后端服务状态..."
if curl -s http://localhost:8888 > /dev/null 2>&1; then
    echo "✅ 后端服务正在运行"
else
    echo "⚠️  后端服务未运行"
    echo "如需完整功能，请启动后端服务："
    echo "  ./main"
fi
echo ""

echo "🚀 启动前端开发服务器..."
echo "💡 前端开发服务器将在 http://localhost:3000 启动"
echo "🎯 登录后在桌面底部Dock栏找到Immich图标"
echo ""
echo "🧪 测试步骤:"
echo "  1. 在浏览器中打开 http://localhost:3000"
echo "  2. 使用管理员账户登录"
echo "  3. 在桌面底部Dock栏找到Immich图标（图片图标）"
echo "  4. 点击图标，应该直接跳转到Immich服务"
echo ""
echo "按 Ctrl+C 停止服务器"
echo ""

# 启动开发服务器
npm run dev
