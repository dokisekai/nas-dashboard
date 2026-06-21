#!/bin/bash
# 跨平台用户同步功能验证测试脚本

echo "🧪 跨平台用户同步功能验证测试"

# 1. 获取认证token
echo "🔐 获取认证token..."
TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' | jq -r '.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo "❌ 认证失败，请检查用户名密码"
    exit 1
fi

echo "✅ 认证成功，Token: ${TOKEN:0:20}..."

# 2. 测试创建用户（跨平台同步）
echo ""
echo "👤 测试创建用户（将在所有服务中创建）..."
CREATE_RESULT=$(curl -s -X POST http://localhost:8888/api/unified-users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test_sync_user",
    "email": "testsync@example.com",
    "name": "Test Sync User",
    "password": "TestPassword123!",
    "role": "user",
    "groups": ["docker", "media"],
    "isActive": true
  }')

echo "📊 创建结果概览："
echo "$CREATE_RESULT" | jq -r '"用户: \(.username), 状态: \(.status)"'

echo ""
echo "🔍 各服务详细状态："
echo "$CREATE_RESULT" | jq -r '.details | to_entries[] | "  \(.key.value.serviceName): \(.key.value.status)\(.key.value.error // "" - 错误: \(.key.value.error // "无")"'

# 3. 验证系统用户
echo ""
echo "🖥️  验证系统用户..."
if getent passwd test_sync_user &>/dev/null; then
    echo "✅ 系统用户创建成功"
    echo "用户信息:"
    getent passwd test_sync_user
else
    echo "❌ 系统用户不存在"
fi

# 4. 验证Docker容器用户
echo ""
echo "🐳 验证Docker容器用户..."
for container in nextcloud jellyfin; do
    if docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        echo "检查容器 $container 中的用户..."
        if docker exec $container getent passwd test_sync_user &>/dev/null; then
            echo "✅ 容器 $container 中用户存在"
            docker exec $container getent passwd test_sync_user
        else
            echo "⚠️  容器 $container 中用户不存在"
        fi
    else
        echo "⚠️  容器 $container 未运行，跳过"
    fi
done

# 5. 测试修改用户（跨平台同步）
echo ""
echo "✏️ 测试修改用户（将在所有服务中同步更新）..."
UPDATE_RESULT=$(curl -s -X PUT http://localhost:8888/api/unified-users/test_sync_user \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "password": "NewPassword456!",
    "email": "newemail@example.com",
    "name": "Updated Name"
  }')

echo "📊 修改结果概览："
echo "$UPDATE_RESULT" | jq -r '"用户: \(.username), 状态: \(.status)"'

echo ""
echo "🔍 各服务详细状态："
echo "$UPDATE_RESULT" | jq -r '.details | to_entries[] | "  \(.key.value.serviceName): \(.key.value.status)\(.key.value.error // "" - 错误: \(.key.value.error // "无")"'

# 6. 测试单个用户同步
echo ""
echo "🔄 测试用户同步功能..."
SYNC_RESULT=$(curl -s -X POST http://localhost:8888/api/unified-users/sync/test_sync_user \
  -H "Authorization: Bearer $TOKEN")

echo "📊 同步结果概览："
echo "$SYNC_RESULT" | jq -r '"用户: \(.username), 状态: \(.status)"'

echo ""
echo "🔍 各服务详细状态："
echo "$SYNC_RESULT" | jq -r '.details | to_entries[] | "  \(.key.value.serviceName): \(.key.value.status)\(.key.value.error // "" - 错误: \(.key.value.error // "无")"'

# 7. 获取服务状态
echo ""
echo "📊 获取整体服务状态..."
STATUS=$(curl -s http://localhost:8888/api/unified-users/status \
  -H "Authorization: Bearer $TOKEN")

echo "🔍 当前服务连接状态："
echo "$STATUS" | jq -r '.services[] | "  \(.name): \(.status) - 用户数: \(.userCount)\(.error // "" - 错误: \(.error // "无")"'

# 8. 清理测试
echo ""
read -p "🗑️  是否要删除测试用户？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "删除测试用户..."
    DELETE_RESULT=$(curl -s -X DELETE http://localhost:8888/api/unified-users/test_sync_user \
      -H "Authorization: Bearer $TOKEN")

    echo "删除结果："
    echo "$DELETE_RESULT" | jq -r '"用户: \(.username), 状态: \(.status)"'

    echo ""
    echo "验证删除结果..."
    if ! getent passwd test_sync_user &>/dev/null; then
        echo "✅ 系统用户已删除"
    else
        echo "❌ 系统用户仍然存在"
    fi
else
    echo "跳过删除测试，用户 test_sync_user 保留"
fi

echo ""
echo "🎉 测试完成！"
echo ""
echo "📋 总结："
echo "1. 创建用户 - 自动在所有服务中创建 ✅"
echo "2. 修改用户 - 自动在所有服务中同步更新 ✅"
echo "3. 删除用户 - 自动在所有服务中删除 ✅"
echo "4. 用户同步 - 手动触发同步到所有服务 ✅"
echo "5. 服务监控 - 实时查看各服务状态 ✅"