#!/bin/bash

# 统一用户管理系统快速集成脚本

set -e

echo "🚀 统一用户管理系统集成脚本"
echo "================================"

# 检查是否在项目根目录
if [ ! -f "backend/cmd/server/main.go" ]; then
    echo "❌ 错误：请在项目根目录运行此脚本"
    exit 1
fi

# 颜色定义
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 步骤1：检查依赖
echo -e "${BLUE}📋 步骤1：检查依赖${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go未安装${NC}"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠️  Docker未安装，部分功能将不可用${NC}"
fi

if ! command -v node &> /dev/null; then
    echo -e "${YELLOW}⚠️  Node.js未安装，前端构建将失败${NC}"
fi

echo -e "${GREEN}✅ 依赖检查完成${NC}"

# 步骤2：配置环境变量
echo -e "${BLUE}📋 步骤2：配置环境变量${NC}"

# 创建.env文件
cat > .env.unified-users << EOF
# Immich配置
IMMICH_API_URL=http://localhost:2283/api
IMMICH_API_KEY=your-immich-api-key-here
IMMICH_ENABLED=true

# 同步配置
UNIFIED_USER_AUTO_SYNC=true
UNIFIED_USER_SYNC_INTERVAL=5m
UNIFIED_USER_RETRY_ATTEMPTS=3
UNIFIED_USER_RETRY_DELAY=2s

# 日志配置
LOG_LEVEL=info
LOG_FILE=/var/log/nas-dashboard/unified-users.log
EOF

echo -e "${GREEN}✅ 环境配置文件已创建：.env.unified-users${NC}"
echo -e "${YELLOW}⚠️  请编辑.env.unified-users文件设置正确的API密钥${NC}"

# 步骤3：集成到main.go
echo -e "${BLUE}📋 步骤3：集成到后端${NC}"

# 备份main.go
cp backend/cmd/server/main.go backend/cmd/server/main.go.backup

# 检查是否已经集成
if grep -q "UnifiedUserManager" backend/cmd/server/main.go; then
    echo -e "${YELLOW}⚠️  统一用户管理已经集成，跳过此步骤${NC}"
else
    echo -e "${YELLOW}⚠️  请手动在main.go中添加以下代码：${NC}"
    cat << 'EOF'

// 在main函数开始处添加
InitUnifiedUserManager()

// 在路由注册处添加
RegisterUnifiedUserRoutes(r, apiGroup)

EOF
    echo -e "${YELLOW}具体位置请参考 UNIFIED_USER_MANAGEMENT_GUIDE.md${NC}"
fi

# 步骤4：前端集成
echo -e "${BLUE}📋 步骤4：集成到前端${NC}"

# 检查前端目录
if [ -d "frontend" ]; then
    # 复制用户管理组件
    cp UnifiedUserManager.vue frontend/src/apps/ 2>/dev/null || echo "组件文件未找到，请手动复制"

    # 检查路由文件
    if [ -f "frontend/src/router/index.ts" ]; then
        if ! grep -q "unified-users" frontend/src/router/index.ts; then
            echo -e "${YELLOW}⚠️  需要手动在router/index.ts中添加路由${NC}"
        fi
    fi

    echo -e "${GREEN}✅ 前端集成准备完成${NC}"
else
    echo -e "${YELLOW}⚠️  前端目录不存在，跳过前端集成${NC}"
fi

# 步骤5：编译后端
echo -e "${BLUE}📋 步骤5：编译后端${NC}"

cd backend
go build -o nas-dashboard cmd/server/main.go
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ 后端编译成功${NC}"
else
    echo -e "${RED}❌ 后端编译失败${NC}"
    exit 1
fi
cd ..

# 步骤6：构建前端
if [ -d "frontend" ]; then
    echo -e "${BLUE}📋 步骤6：构建前端${NC}"

    cd frontend
    npm install
    npm run build
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ 前端构建成功${NC}"
    else
        echo -e "${YELLOW}⚠️  前端构建失败，但不影响后端功能${NC}"
    fi
    cd ..
fi

# 步骤7：启动服务
echo -e "${BLUE}📋 步骤7：配置启动服务${NC}"

# 创建systemd服务文件
cat << EOF > /tmp/nas-dashboard.service
[Unit]
Description=NAS Dashboard with Unified User Management
After=network.target docker.service

[Service]
Type=simple
User=root
WorkingDirectory=/data/nas-dashboard
Environment="PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
EnvironmentFile=/data/nas-dashboard/.env.unified-users
ExecStart=/data/nas-dashboard/backend/nas-dashboard
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

echo -e "${GREEN}✅ Systemd服务文件已创建${NC}"
echo -e "${YELLOW}⚠️  要启用服务，请运行：${NC}"
echo "   sudo cp /tmp/nas-dashboard.service /etc/systemd/system/"
echo "   sudo systemctl daemon-reload"
echo "   sudo systemctl enable nas-dashboard"
echo "   sudo systemctl start nas-dashboard"

# 步骤8：测试
echo -e "${BLUE}📋 步骤8：功能测试${NC}"

echo -e "${YELLOW}请启动服务后运行以下测试：${NC}"
cat << 'EOF'
# 1. 测试健康检查
curl http://localhost:8888/health

# 2. 测试统一用户状态（需要JWT token）
TOKEN=$(curl -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password"}' | jq -r '.token')

curl http://localhost:8888/api/unified-users/status \
  -H "Authorization: Bearer $TOKEN"

# 3. 测试用户创建
curl -X POST http://localhost:8888/api/unified-users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "name": "Test User",
    "password": "testpass123",
    "role": "user"
  }'
EOF

# 完成
echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}🎉 集成完成！${NC}"
echo -e "${GREEN}=========================================${NC}"

echo -e "\n📚 重要提醒："
echo -e "${YELLOW}1. 编辑 .env.unified-users 设置正确的API密钥${NC}"
echo -e "${YELLOW}2. 检查 main.go 中的集成代码${NC}"
echo -e "${YELLOW}3. 启动服务后进行功能测试${NC}"
echo -e "${YELLOW}4. 查看 UNIFIED_USER_MANAGEMENT_GUIDE.md 了解详细用法${NC}"

echo -e "\n🚀 快速启动："
echo "   cd backend && ./nas-dashboard"
echo "   或使用systemd服务"

# 创建快速测试脚本
cat << 'EOF' > test-unified-users.sh
#!/bin/bash

echo "🧪 统一用户管理系统测试"

# 检查服务是否运行
if ! curl -s http://localhost:8888/health > /dev/null; then
    echo "❌ 服务未运行，请先启动服务"
    exit 1
fi

echo "✅ 服务运行正常"

# 获取token（请修改用户名密码）
echo "🔐 获取认证token..."
TOKEN=$(curl -s -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your-password"}' | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "❌ 认证失败，请修改脚本中的用户名密码"
    exit 1
fi

echo "✅ 认证成功"

# 测试获取用户状态
echo "📊 测试获取服务状态..."
STATUS=$(curl -s http://localhost:8888/api/unified-users/status \
  -H "Authorization: Bearer $TOKEN")

if echo "$STATUS" | jq -e '.services' > /dev/null; then
    echo "✅ 服务状态获取成功"
    echo "$STATUS" | jq '.'
else
    echo "❌ 服务状态获取失败"
    echo "$STATUS"
fi

# 测试获取用户列表
echo "👥 测试获取用户列表..."
USERS=$(curl -s http://localhost:8888/api/unified-users \
  -H "Authorization: Bearer $TOKEN")

if echo "$USERS" | jq -e '.users' > /dev/null; then
    USER_COUNT=$(echo "$USERS" | jq '.users | length')
    echo "✅ 用户列表获取成功，共 $USER_COUNT 个用户"
else
    echo "❌ 用户列表获取失败"
fi

# 提示功能测试
echo ""
echo "🎯 要测试完整功能，请访问前端界面："
echo "   http://localhost:5173/unified-users"
echo ""
echo "或使用API测试："
echo "   # 创建用户"
echo '   curl -X POST http://localhost:8888/api/unified-users \'
echo '     -H "Authorization: Bearer $TOKEN" \'
echo '     -H "Content-Type: application/json" \'
echo '     -d '"'"'{"username":"testuser","email":"test@example.com","name":"Test User","password":"testpass123"}'"'"''
echo ""
echo "   # 同步所有用户"
echo "   curl -X POST http://localhost:8888/api/unified-users/sync \"
echo "     -H \"Authorization: Bearer $TOKEN\""

EOF

chmod +x test-unified-users.sh
echo -e "${GREEN}✅ 测试脚本已创建：test-unified-users.sh${NC}"

echo -e "\n🎉 集成脚本执行完成！"
echo -e "📖 请查看 UNIFIED_USER_MANAGEMENT_GUIDE.md 了解详细用法"