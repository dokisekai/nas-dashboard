#!/bin/bash
# Immich用户管理启用脚本

echo "🚀 Immich用户管理启用向导"
echo "=================================="
echo ""

# 检查Immich服务状态
echo "🔍 步骤1: 检查Immich服务状态..."
if curl -s http://localhost:2283 > /dev/null; then
    echo "✅ Immich服务运行中"
else
    echo "❌ Immich服务未运行，请先启动Immich服务"
    exit 1
fi
echo ""

# 获取API密钥
echo "🔑 步骤2: 获取Immich API密钥"
echo "请按以下步骤操作："
echo "1. 访问 http://localhost:2283"
echo "2. 使用管理员账号登录"
echo "3. 进入 Administration > API Keys"
echo "4. 创建新的API密钥"
echo "5. 复制生成的API密钥"
echo ""
read -p "请输入Immich API密钥: " api_key

if [ -z "$api_key" ]; then
    echo "❌ API密钥不能为空"
    exit 1
fi
echo ""

# 创建配置文件
echo "📝 步骤3: 创建配置文件..."
config_dir="/data/nas-dashboard/backend"
config_file="$config_dir/.env.immich"

cat > "$config_file" << EOF
# Immich API配置
IMMICH_API_URL=http://localhost:2283/api
IMMICH_API_KEY=$api_key
EOF

echo "✅ 配置文件已创建: $config_file"
echo ""

# 测试连接
echo "🧪 步骤4: 测试Immich API连接..."
test_response=$(curl -s -X GET "http://localhost:2283/api/users" \
    -H "x-immich-api-key: $api_key" 2>/dev/null)

if [ $? -eq 0 ] && [ -n "$test_response" ]; then
    echo "✅ API连接测试成功"

    # 尝试解析用户数量
    user_count=$(echo "$test_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data) if isinstance(data, list) else 'N/A')" 2>/dev/null)
    echo "📊 当前Immich用户数量: $user_count"
else
    echo "❌ API连接测试失败"
    echo "请检查API密钥是否正确"
    exit 1
fi
echo ""

# 重启后端服务
echo "🔄 步骤5: 重启后端服务..."
echo "配置完成！请重启后端服务以应用配置："
echo ""
echo "cd /data/nas-dashboard/backend"
echo "./main"
echo ""

# 测试用户管理API
echo "🎯 步骤6: 测试用户管理API..."
echo "请先登录获取Token，然后测试Immich用户管理："
echo ""
echo "1. 登录获取Token:"
echo '   curl -X POST http://localhost:8888/api/auth/login \'
echo '     -H "Content-Type: application/json" \'
echo '     -d '"'"'{"username":"admin","password":"admin"}'"'"''
echo ""
echo "2. 测试Immich用户管理:"
echo '   curl -H "Authorization: Bearer YOUR_TOKEN" \'
echo '     http://localhost:8888/api/immich/users'
echo ""

echo "🎉 Immich用户管理配置完成！"
echo ""
echo "💡 可用的Immich用户管理操作:"
echo "  - 查看用户列表: GET /api/immich/users"
echo "  - 创建用户: POST /api/immich/users"
echo "  - 更新用户: PUT /api/immich/users/{id}"
echo "  - 删除用户: DELETE /api/immich/users/{id}"
echo "  - 用户同步: POST /api/immich/users/sync"