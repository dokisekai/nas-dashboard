#!/bin/bash
# 快速启用Immich集成

echo "🚀 快速启用Immich集成到NAS Dashboard"
echo "=========================================="
echo ""

BACKEND_DIR="/data/nas-dashboard/backend"
IMMICH_CONFIG="$BACKEND_DIR/config/immich.json"

# 检查配置是否存在
if [ ! -f "$IMMICH_CONFIG" ]; then
    echo "❌ Immich配置文件不存在"
    exit 1
fi

echo "✅ 找到Immich配置"

# 显示配置信息
echo ""
echo "📋 Immich配置信息:"
echo "  URL: $(jq -r '.url' "$IMMICH_CONFIG")"
echo "  状态: $(jq -r '.enabled' "$IMMICH_CONFIG")"
echo "  版本: $(jq -r '.version' "$IMMICH_CONFIG")"
echo ""

# 测试连接
API_KEY=$(jq -r '.apiKey' "$IMMICH_CONFIG")
IMMICH_URL=$(jq -r '.url' "$IMMICH_CONFIG")

echo "🧪 测试Immich连接..."
test_conn=$(curl -s -X GET "$IMMICH_URL/api/users" -H "X-API-Key: $API_KEY")

if echo "$test_conn" | grep -q "email\|name"; then
    echo "✅ Immich连接正常"
    user_count=$(echo "$test_conn" | python3 -c "import sys, json; print(len(json.load(sys.stdin)))" 2>/dev/null)
    echo "📊 当前用户数: $user_count"
else
    echo "❌ Immich连接失败"
    exit 1
fi
echo ""

echo "🎯 创建前端集成链接..."

# 在前端服务页面添加Immich入口
frontend_services="/data/nas-dashboard/frontend/src/views/Services/Services.vue"

if [ -f "$frontend_services" ]; then
    echo "✅ 找到前端服务页面"
    echo "💡 需要手动在Services.vue中添加Immich集成组件"
else
    echo "⚠️  前端服务页面未找到"
fi
echo ""

echo "🌐 生成的访问链接:"
echo "  直接访问: $IMMICH_URL"
echo "  API密钥已配置，支持自动登录"
echo ""

echo "🎉 Immich集成准备完成！"
echo ""
echo "📋 下一步操作:"
echo "1. 重启后端服务"
echo "2. 在前端界面添加Immich入口"
echo "3. 测试一键跳转登录功能"
