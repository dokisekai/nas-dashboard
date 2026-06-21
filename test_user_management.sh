#!/bin/bash
# NAS Dashboard 用户管理和认证功能测试脚本

echo "🧪 NAS Dashboard 用户管理和认证功能测试"
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
NC='\033[0m' # No Color

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
test_api() {
    local test_name="$1"
    local method="$2"
    local endpoint="$3"
    local data="$4"
    local token="$5"
    local expected_code="$6"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo "测试 $TOTAL_TESTS: $test_name"

    # 构建curl命令
    if [ -z "$token" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Authorization: Bearer $token" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_URL$endpoint")
    fi

    # 分离响应体和状态码
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n-1)

    # 检查状态码
    if [ "$status_code" == "$expected_code" ]; then
        echo -e "${GREEN}✅ PASS${NC}: 状态码 $status_code (期望: $expected_code)"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo "响应: $body" | head -c 200
        echo ""
    else
        echo -e "${RED}❌ FAIL${NC}: 状态码 $status_code (期望: $expected_code)"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo "响应: $body"
    fi
    echo ""
}

echo "📋 测试准备..."
echo "检查服务器状态..."
health_response=$(curl -s "$API_URL/health")
if [ "$health_response" == '{"status":"ok"}' ]; then
    echo -e "${GREEN}✅ 服务器运行正常${NC}"
else
    echo -e "${RED}❌ 服务器未运行${NC}"
    exit 1
fi
echo ""

echo "🔐 测试1: 管理员登录"
login_response=$(curl -s -X POST "$API_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASS\"}")

echo "登录响应: $login_response" | head -c 300
echo ""

# 提取token
ADMIN_TOKEN=$(echo "$login_response" | python3 -c "import sys, json; print(json.load(sys.stdin).get('token', ''))" 2>/dev/null)

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" == "None" ]; then
    echo -e "${RED}❌ 登录失败，无法获取token${NC}"
    exit 1
else
    echo -e "${GREEN}✅ 登录成功，获取到token${NC}"
    echo "Token: ${ADMIN_TOKEN:0:50}..."
fi
echo ""

echo "🐳 测试2: Docker容器管理 (需要认证)"
test_api "获取Docker容器列表" \
    "GET" \
    "/api/docker/containers" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "👥 测试3: 用户管理 (需要认证)"
test_api "获取系统用户列表" \
    "GET" \
    "/api/users" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "💾 测试4: 系统监控 (需要认证)"
test_api "获取CPU信息" \
    "GET" \
    "/api/monitor/cpu" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "🌐 测试5: 网络管理 (需要认证)"
test_api "获取网络接口" \
    "GET" \
    "/api/network/interfaces" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "📁 测试6: 存储管理 (需要认证)"
test_api "获取磁盘信息" \
    "GET" \
    "/api/storage/disks" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "🔧 测试7: 系统服务管理 (需要认证)"
test_api "获取系统服务" \
    "GET" \
    "/api/services" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "🚨 测试8: 未授权访问测试"
test_api "未授权访问Docker API" \
    "GET" \
    "/api/docker/containers" \
    "" \
    "" \
    "401"

echo "🚨 测试9: 无效token测试"
test_api "使用无效token访问API" \
    "GET" \
    "/api/docker/containers" \
    "" \
    "invalid_token_12345" \
    "401"

echo "🔑 测试10: Token刷新"
refresh_token=$(echo "$login_response" | python3 -c "import sys, json; print(json.load(sys.stdin).get('refreshToken', ''))" 2>/dev/null)
test_api "刷新访问token" \
    "POST" \
    "/api/auth/refresh" \
    "{\"refreshToken\":\"$refresh_token\"}" \
    "" \
    "200"

echo "👤 测试11: 获取当前用户信息"
test_api "获取当前管理员用户信息" \
    "GET" \
    "/api/users/me" \
    "" \
    "$ADMIN_TOKEN" \
    "200"

echo "📊 测试总结"
echo "=========================================="
echo "总测试数: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"
echo "成功率: $(awk "BEGIN {printf \"%.1f\", ($PASSED_TESTS/$TOTAL_TESTS)*100}")%"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 所有测试通过！${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  有 $FAILED_TESTS 个测试失败${NC}"
    exit 1
fi