#!/bin/bash
# NAS Dashboard 应用用户管理集成测试脚本

echo "🔍 NAS Dashboard 应用用户管理集成测试"
echo "=========================================="
echo ""

# 配置
API_URL="http://localhost:8888"
ADMIN_USER="admin"
ADMIN_PASS="admin"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 登录获取token
echo "🔐 步骤1: 管理员登录..."
login_response=$(curl -s -X POST "$API_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASS\"}")

ADMIN_TOKEN=$(echo "$login_response" | python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))" 2>/dev/null)

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" == "None" ]; then
    echo -e "${RED}❌ 登录失败${NC}"
    exit 1
else
    echo -e "${GREEN}✅ 登录成功${NC}"
fi
echo ""

# 测试应用服务状态
echo "🌐 步骤2: 检查应用服务状态..."

check_service() {
    local service_name="$1"
    local service_url="$2"
    local expected_response="$3"

    echo -n "  $service_name ($service_url): "
    response=$(curl -s -X GET "$service_url" 2>/dev/null)

    if [ $? -eq 0 ] && [ -n "$response" ]; then
        if echo "$response" | grep -q "$expected_response"; then
            echo -e "${GREEN}✅ 运行中${NC}"
            return 0
        else
            echo -e "${YELLOW}⚠️  运行但响应异常${NC}"
            return 1
        fi
    else
        echo -e "${RED}❌ 未运行${NC}"
        return 2
    fi
}

# 检查各个应用服务
check_service "Immich" "http://localhost:2283" "Immich"
check_service "Authentik" "http://localhost:9000" "authentik"
check_service "Private Git" "http://localhost:3000" "Forgejo"
check_service "AList" "http://localhost:5244" "AList"

echo ""

# 测试各应用的用户管理API
echo "👥 步骤3: 测试应用用户管理API..."

test_app_user_api() {
    local app_name="$1"
    local api_endpoint="$2"
    local token="$3"

    echo -n "  $app_name 用户管理: "

    response=$(curl -s -H "Authorization: Bearer $token" "$api_endpoint" 2>/dev/null)
    status_code=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $token" "$api_endpoint")

    if [ "$status_code" == "200" ]; then
        user_count=$(echo "$response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data) if isinstance(data, list) else data.get('total', 'N/A'))" 2>/dev/null)
        echo -e "${GREEN}✅ 可用 ($user_count 个用户)${NC}"
        return 0
    elif [ "$status_code" == "401" ]; then
        echo -e "${YELLOW}⚠️  需要配置API密钥${NC}"
        return 1
    elif [ "$status_code" == "404" ]; then
        echo -e "${YELLOW}⚠️  API未实现${NC}"
        return 1
    else
        echo -e "${RED}❌ 错误 ($status_code)${NC}"
        return 2
    fi
}

# 测试各个应用的用户管理API
test_app_user_api "Immich" "$API_URL/api/immich/users" "$ADMIN_TOKEN"
test_app_user_api "Authentik" "$API_URL/api/authentik/users" "$ADMIN_TOKEN"
test_app_user_api "系统用户" "$API_URL/api/users" "$ADMIN_TOKEN"

echo ""

# 检查统一用户管理功能
echo "🔄 步骤4: 检查统一用户管理功能..."

echo -n "  统一用户管理路由: "
unified_response=$(curl -s "$API_URL/api/unified-users/status" 2>/dev/null)
status_code=$(curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/unified-users/status")

if [ "$status_code" == "401" ]; then
    echo -e "${GREEN}✅ 路由已注册 (需要认证)${NC}"
elif [ "$status_code" == "404" ]; then
    echo -e "${YELLOW}⚠️  路由未注册${NC}"
else
    echo -e "${RED}❌ 状态异常 ($status_code)${NC}"
fi

echo ""

# 检查配置文件状态
echo "⚙️  步骤5: 检查应用配置..."

check_config() {
    local config_file="$1"
    local app_name="$2"

    echo -n "  $app_name 配置文件: "
    if [ -f "$config_file" ]; then
        echo -e "${GREEN}✅ 存在${NC}"
    else
        echo -e "${YELLOW}⚠️  不存在${NC}"
    fi
}

check_config ".env.immich" "Immich"
check_config ".env.authentik" "Authentik"

echo ""

# 生成配置建议
echo "💡 步骤6: 生成配置建议..."

echo "建议的配置文件："
echo ""

# Immich配置建议
echo "📝 .env.immich (如需要集成Immich用户管理):"
cat << 'EOF'
IMMICH_API_URL=http://localhost:2283/api
IMMICH_API_KEY=your_immich_api_key_here
EOF
echo ""

# 检查是否能获取Immich API密钥
echo -e "${BLUE}ℹ️  获取Immich API密钥的方法:${NC}"
echo "1. 访问 http://localhost:2283"
echo "2. 使用admin账号登录"
echo "3. 进入 Administration > API Keys"
echo "4. 创建新的API密钥"
echo ""

# 当前运行的应用统计
echo "📊 步骤7: 当前运行的应用用户管理统计..."

echo "应用用户管理集成状态："
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查Immich用户管理
echo -n "Immich用户管理: "
immich_response=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$API_URL/api/immich/users" 2>/dev/null)
immich_status=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $ADMIN_TOKEN" "$API_URL/api/immich/users")

if [ "$immich_status" == "200" ]; then
    immich_users=$(echo "$immich_response" | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data) if isinstance(data, list) else data.get('total', 0))" 2>/dev/null)
    echo -e "${GREEN}✅ 可用${NC} ($immich_users 个用户)"
elif [ "$immich_status" == "401" ]; then
    echo -e "${YELLOW}⚠️  需要配置API密钥${NC}"
else
    echo -e "${RED}❌ 不可用 ($immich_status)${NC}"
fi

# 检查系统用户管理
echo -n "系统用户管理: "
system_users=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$API_URL/api/users" 2>/dev/null | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data.get('users', [])))" 2>/dev/null)
if [ -n "$system_users" ] && [ "$system_users" != "0" ]; then
    echo -e "${GREEN}✅ 可用${NC} ($system_users 个用户)"
else
    echo -e "${RED}❌ 不可用${NC}"
fi

# 检查Docker容器管理
echo -n "Docker容器管理: "
docker_containers=$(curl -s -H "Authorization: Bearer $ADMIN_TOKEN" "$API_URL/api/docker/containers" 2>/dev/null | python3 -c "import sys, json; data=json.load(sys.stdin); print(len(data.get('containers', [])))" 2>/dev/null)
if [ -n "$docker_containers" ] && [ "$docker_containers" != "0" ]; then
    echo -e "${GREEN}✅ 可用${NC} ($docker_containers 个容器)"
else
    echo -e "${RED}❌ 不可用${NC}"
fi

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 总结
echo "📋 总结:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ 已集成的用户管理:"
echo "  - 系统用户管理 (Linux用户)"
echo "  - Docker容器管理"
echo ""
echo "⚠️  需要配置的用户管理:"
echo "  - Immich用户管理 (需要API密钥)"
echo "  - Authentik用户管理 (需要API密钥)"
echo "  - 统一用户管理 (需要注册路由)"
echo ""
echo "🎯 建议的下一步操作:"
echo "1. 配置Immich API密钥以启用用户管理"
echo "2. 注册统一用户管理路由"
echo "3. 测试跨应用用户同步功能"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"