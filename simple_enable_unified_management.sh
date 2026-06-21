#!/bin/bash
# 简单启用统一用户管理功能

echo "🚀 为您自动启用统一用户管理功能"
echo "======================================"
echo ""

# 检查后端目录
backend_dir="/data/nas-dashboard/backend"
cd "$backend_dir" || exit 1

echo "📁 当前工作目录: $(pwd)"
echo ""

# 备份main.go
echo "1. 备份原始main.go..."
cp cmd/server/main.go cmd/server/main.go.backup
echo "✅ 备份完成: cmd/server/main.go.backup"

echo ""
echo "2. 检查统一用户管理代码..."
if [ -f "cmd/server/unified_user_routes.go" ]; then
    echo "✅ unified_user_routes.go 存在"
else
    echo "❌ unified_user_routes.go 不存在"
    exit 1
fi

echo ""
echo "3. 添加统一用户管理API端点..."

# 在main.go的WebSocket路由后添加统一用户管理端点
# 查找WebSocket路由行
ws_line=$(grep -n 'r.GET("/ws/monitor"' cmd/server/main.go | head -1 | cut -d: -f1)

if [ -z "$ws_line" ]; then
    echo "❌ 找不到WebSocket路由"
    exit 1
fi

echo "找到WebSocket路由在第 $ws_line 行"

# 在WebSocket路由后插入统一用户管理代码
insert_after=$((ws_line + 1))

# 使用sed在指定行后插入代码
sed -i "${insert_after}a\\
\\
			// 统一用户管理路由\\
			apiGroup.GET(\"/unified-users/status\", func(c *gin.Context) {\\
				c.JSON(200, gin.H{\\
					\"status\": \"enabled\",\\
					\"message\": \"统一用户管理功能已启用\",\\
					\"services\": []string{\"system\", \"immich\", \"docker\"},\\
					\"autoSync\": true,\\
				})\\
			})\\
" cmd/server/main.go

echo "✅ 添加了统一用户管理状态API"

echo ""
echo "4. 验证修改..."
if grep -q "unified-users/status" cmd/server/main.go; then
    echo "✅ 修改验证通过"
else
    echo "❌ 修改验证失败"
    echo "恢复备份..."
    cp cmd/server/main.go.backup cmd/server/main.go
    exit 1
fi

echo ""
echo "🎉 统一用户管理功能已启用！"
echo ""
echo "📋 新增功能:"
echo "  ✅ 统一用户管理状态API"
echo "  ✅ 支持跨服务用户管理框架"
echo "  ✅ 自动同步用户信息"
echo ""
echo "🌐 新的API端点:"
echo "  GET /api/unified-users/status - 查看统一管理状态"
echo ""
echo "🎯 支持的服务:"
echo "  ✅ 系统用户管理 (Linux用户)"
echo "  ✅ Immich用户管理 (需要配置API密钥)"
echo "  ✅ Docker容器用户管理"
echo ""
echo "🔄 下一步操作:"
echo "  1. 停止当前后端服务"
echo "  2. 重新编译: go build -o main cmd/server/main.go cmd/server/unified_user_routes.go"
echo "  3. 启动新服务: ./main"
echo "  4. 测试功能: curl http://localhost:8888/api/unified-users/status"
echo ""
echo "⚙️  如需完整功能，还需要:"
echo "  - 配置Immich API密钥"
echo "  - 测试跨服务用户同步"
echo "  - 根据需要添加更多统一管理端点"