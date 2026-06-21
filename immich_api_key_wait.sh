#!/bin/bash
# Immich API密钥配置脚本

echo "🔑 等待您获取Immich API密钥..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 提示：请按照以下步骤操作："
echo "1. 访问: http://localhost:2283"
echo "2. 管理员登录 → Administration → API Keys"
echo "3. 创建新的API密钥（记得勾选用户管理权限）"
echo "4. 复制生成的密钥"
echo ""
echo "获取到密钥后，请按回车键继续..."
read -p "按回车键继续..."

echo ""
echo "📝 请输入您的Immich API密钥："
read -p "API密钥: " api_key

if [ -z "$api_key" ]; then
    echo "❌ API密钥不能为空"
    exit 1
fi

# 验证API密钥格式（简单检查）
if [ ${#api_key} -lt 10 ]; then
    echo "⚠️  API密钥长度似乎太短，请确认是否正确"
    read -p "仍要继续吗？(y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo ""
echo "✅ API密钥已提供: ${api_key:0:10}..."

# 配置文件路径
config_file="/data/nas-dashboard/backend/.env.immich"
backup_file="/data/nas-dashboard/backend/.env.immich.backup"

# 备份原配置文件
if [ -f "$config_file" ]; then
    cp "$config_file" "$backup_file"
    echo "📋 已备份原配置文件"
fi

# 创建新配置文件
cat > "$config_file" << EOF
# Immich API配置
# 自动生成于: $(date)
IMMICH_API_URL=http://localhost:2283/api
IMMICH_API_KEY=$api_key
EOF

echo "✅ 配置文件已创建: $config_file"
echo ""

# 测试API连接
echo "🧪 测试Immich API连接..."
test_response=$(curl -s -X GET "http://localhost:2283/api/users" \
    -H "x-immich-api-key: $api_key" 2>&1)

if echo "$test_response" | grep -q "email\|name\|id"; then
    echo "✅ API连接测试成功"

    # 尝试获取用户数量
    user_count=$(echo "$test_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data) if isinstance(data, list) else 'N/A')" 2>/dev/null)
    echo "📊 当前Immich用户数量: $user_count"
else
    echo "❌ API连接测试失败"
    echo "响应: $test_response"
    echo ""
    echo "可能的原因："
    echo "1. API密钥不正确"
    echo "2. API密钥权限不足（需要用户管理权限）"
    echo "3. Immich服务异常"
    exit 1
fi

echo ""
echo "🎉 Immich用户管理配置完成！"
echo ""
echo "🚀 下一步操作："
echo "1. 重启后端服务: cd /data/nas-dashboard/backend && ./main"
echo "2. 测试用户管理API"
echo "3. 开始使用统一用户管理功能"
echo ""
echo "💡 快速测试命令："
echo 'export TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login -H "Content-Type: application/json" -d '"'"'{"username":"admin","password":"admin"}'"'"' | jq -r ".token")'
echo 'curl -H "Authorization: Bearer $TOKEN" http://localhost:8888/api/immich/users'