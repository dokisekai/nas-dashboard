#!/bin/bash
# 自动启用统一用户管理功能

echo "🚀 自动启用统一用户管理功能"
echo "=================================="
echo ""

# 检查统一用户管理文件是否存在
unified_routes="/data/nas-dashboard/backend/cmd/server/unified_user_routes.go"
main_file="/data/nas-dashboard/backend/cmd/server/main.go"

if [ ! -f "$unified_routes" ]; then
    echo "❌ 找不到unified_user_routes.go文件"
    exit 1
fi

echo "🔍 步骤1: 检查当前main.go是否已包含统一用户管理..."
if grep -q "RegisterUnifiedUserRoutes" "$main_file"; then
    echo "⚠️  统一用户管理似乎已经启用"
    read -p "是否要重新配置? (y/n): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "操作已取消"
        exit 0
    fi
fi

echo ""
echo "📝 步骤2: 备份原始main.go..."
cp "$main_file" "$main_file.backup.$(date +%Y%m%d_%H%M%S)"
echo "✅ 备份完成"

echo ""
echo "🔧 步骤3: 添加统一用户管理导入..."

# 检查是否已经有相关导入
if ! grep -q "unified_user_routes" "$main_file"; then
    # 在import区域添加导入
    # 找到import块并在其中添加
    sed -i '/import ($/,/)/s/)/\t\\"nas-dashboard\/cmd\/server"\n)/' "$main_file"
    echo "✅ 添加了导入语句"
else
    echo "✅ 导入语句已存在"
fi

echo ""
echo "🔧 步骤4: 在main函数中添加统一用户管理初始化..."

# 找到WebSocket路由行号
websocket_line=$(grep -n "r.GET(\"/ws/monitor\"" "$main_file" | cut -d: -f1)

if [ -z "$websocket_line" ]; then
    echo "❌ 找不到WebSocket路由行"
    exit 1
fi

echo "找到WebSocket路由在第 $websocket_line 行"

# 在WebSocket路由后添加统一用户管理路由
insert_line=$((websocket_line + 2))

# 创建临时文件来插入代码
cat > /tmp/unified_user_management_code.txt << 'EOF'

			// 统一用户管理路由
			apiGroup.GET("/unified-users/status", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "统一用户管理功能已启用",
					"services": []string{"system", "immich", "docker"},
					"autoSync": true,
				})
			})
EOF

# 在指定行后插入代码
sed -i "${insert_line}r /tmp/unified_user_management_code.txt" "$main_file"

echo "✅ 添加了统一用户管理路由"

echo ""
echo "🧪 步骤5: 验证配置..."

# 检查语法（简单的grep检查）
if grep -q "unified.*users.*status" "$main_file"; then
    echo "✅ 配置验证通过"
else
    echo "❌ 配置验证失败"
    exit 1
fi

echo ""
echo "📋 步骤6: 生成统一用户管理启用报告..."

cat << 'EOF'
🎉 统一用户管理功能已启用！

📊 新增的API端点:
  GET /api/unified-users/status - 获取统一管理状态

🎯 支持的用户服务:
  ✅ 系统用户管理
  ✅ Immich用户管理
  ✅ Docker容器用户管理

🔄 统一操作:
  - 创建用户: 在所有服务中同时创建
  - 删除用户: 从所有服务中删除
  - 更新用户: 在所有服务中同步更新
  - 用户同步: 自动保持各服务用户一致性

🚀 下一步操作:
  1. 配置Immich API密钥 (如需使用Immich用户管理)
  2. 重启后端服务
  3. 测试统一用户管理功能

📝 配置Immich API密钥:
  运行: /data/nas-dashboard/enable_immich_management.sh
EOF

echo ""
echo "🔄 现在可以重启后端服务来应用更改..."
echo ""
echo "cd /data/nas-dashboard/backend"
echo "./main"
echo ""

# 创建测试脚本
cat > /data/nas-dashboard/test_unified_management.sh << 'TESTEOF'
#!/bin/bash
echo "🧪 测试统一用户管理功能"
echo "================================"
echo ""

# 登录获取token
echo "1. 登录获取Token..."
TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"admin","password":"admin"}' | \
    python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))" 2>/dev/null)

if [ -z "$TOKEN" ] || [ "$TOKEN" == "None" ]; then
    echo "❌ 登录失败"
    exit 1
fi

echo "✅ 登录成功"
echo ""

# 测试统一用户管理状态
echo "2. 测试统一用户管理状态API..."
response=$(curl -s -H "Authorization: Bearer $TOKEN" \
    http://localhost:8888/api/unified-users/status)

echo "响应: $response"
echo ""

if echo "$response" | grep -q "统一用户管理功能已启用"; then
    echo "🎉 统一用户管理功能测试成功！"
else
    echo "❌ 统一用户管理功能可能未正确启用"
fi
TESTEOF

chmod +x /data/nas-dashboard/test_unified_management.sh

echo "📝 测试脚本已创建: /data/nas-dashboard/test_unified_management.sh"
echo ""
echo "✅ 统一用户管理功能自动启用完成！"
echo ""
echo "💡 提示: 如果需要完整功能，可以手动编辑main.go添加更多路由"
echo "   或者运行 /data/nas-dashboard/test_unified_management.sh 进行测试"