#!/bin/bash
# 安全的Immich配置管理

echo "🔒 Immich安全配置管理"
echo ""

# 从环境变量读取配置（更安全）
IMMICH_URL="${IMMICH_URL:-http://localhost:2283}"
IMMICH_API_KEY="${IMMICH_API_KEY}"

if [ -z "$IMMICH_API_KEY" ]; then
    echo "❌ 请设置环境变量 IMMICH_API_KEY"
    exit 1
fi

echo "✅ 配置已从环境变量加载"
echo "  URL: $IMMICH_URL"
echo "  API密钥: ${IMMICH_API_KEY:0:10}... (已隐藏)"
echo ""

# 测试连接
echo "🧪 测试连接..."
test_response=$(curl -s -X GET "$IMMICH_URL/api/users" -H "X-API-Key: $IMMICH_API_KEY")

if echo "$test_response" | grep -q "email\|name"; then
    echo "✅ 连接测试成功"
    user_count=$(echo "$test_response" | python3 -c "import sys, json; print(len(json.load(sys.stdin)))" 2>/dev/null)
    echo "📊 当前用户: $user_count 个"
else
    echo "❌ 连接测试失败"
    exit 1
fi
