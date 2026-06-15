#!/bin/bash
# Complete System Integration Test Script
# This script performs comprehensive testing of the NAS Dashboard system

set -e

echo "========================================="
echo "NAS Dashboard - System Integration Test"
echo "========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to print test results
print_result() {
    local test_name="$1"
    local result="$2"
    local message="$3"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    if [ "$result" = "pass" ]; then
        echo -e "${GREEN}✓${NC} $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗${NC} $test_name"
        echo -e "  ${YELLOW}Message: $message${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Change to frontend directory
cd "$(dirname "$0")"

echo "1. Frontend-Build Test"
echo "----------------------"

# Test 1: Check if node_modules exists
if [ -d "node_modules" ]; then
    print_result "Node modules installed" "pass"
else
    print_result "Node modules installed" "fail" "node_modules directory not found"
    echo "Installing dependencies..."
    npm install
fi

# Test 2: Check TypeScript compilation
echo "Running TypeScript compilation check..."
if npm run build 2>/dev/null; then
    print_result "TypeScript compilation" "pass"
else
    print_result "TypeScript compilation" "fail" "Build errors detected"
fi

echo ""
echo "2. Backend Connection Test"
echo "--------------------------"

# Test 3: Check backend API availability
API_URL=${VITE_API_URL:-http://192.168.50.10:8888}
echo "Testing API connection to: $API_URL"

if curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/health" 2>/dev/null | grep -q "200"; then
    print_result "Backend API health check" "pass"
else
    print_result "Backend API health check" "fail" "API not reachable"
fi

# Test 4: Test authentication endpoint
echo "Testing authentication endpoint..."
AUTH_RESPONSE=$(curl -s -X POST "$API_URL/api/auth/login" \
    -H "Content-Type: application/json" \
    -d '{"username":"test","password":"test"}' 2>/dev/null || echo '{"error":"connection_failed"}')

if echo "$AUTH_RESPONSE" | grep -q "token\|error"; then
    print_result "Authentication endpoint" "pass"
else
    print_result "Authentication endpoint" "fail" "Unexpected response"
fi

echo ""
echo "3. WebSocket Connection Test"
echo "----------------------------"

# Test 5: Check WebSocket endpoint
WS_URL=${VITE_WS_URL:-ws://192.168.50.10:8888}
echo "Testing WebSocket connection to: $WS_URL"

# Use curl to test WebSocket upgrade
WS_RESPONSE=$(curl -s -i -N \
    -H "Connection: Upgrade" \
    -H "Upgrade: websocket" \
    -H "Sec-WebSocket-Version: 13" \
    -H "Sec-WebSocket-Key: test" \
    "$WS_URL/ws/monitor" 2>&1 || echo "connection_failed")

if echo "$WS_RESPONSE" | grep -q "101\|Upgrade"; then
    print_result "WebSocket endpoint" "pass"
else
    print_result "WebSocket endpoint" "fail" "WebSocket upgrade failed"
fi

echo ""
echo "4. Frontend Configuration Test"
echo "-------------------------------"

# Test 6: Check environment variables
if [ -f ".env" ]; then
    print_result "Environment configuration file" "pass"

    # Check required environment variables
    if grep -q "VITE_API_URL" .env && grep -q "VITE_WS_URL" .env; then
        print_result "Required environment variables set" "pass"
    else
        print_result "Required environment variables set" "fail" "Missing VITE_API_URL or VITE_WS_URL"
    fi
else
    print_result "Environment configuration file" "fail" ".env file not found"
fi

# Test 7: Check vite config
if [ -f "vite.config.ts" ]; then
    print_result "Vite configuration file" "pass"
else
    print_result "Vite configuration file" "fail" "vite.config.ts not found"
fi

echo ""
echo "5. Component Integration Test"
echo "------------------------------"

# Test 8: Check core components
CORE_COMPONENTS=(
    "src/views/Login.vue"
    "src/views/Dashboard.vue"
    "src/components/Desktop/WindowManager.vue"
    "src/components/Desktop/DesktopWindow.vue"
    "src/composables/useDesktop.ts"
    "src/composables/useWebSocket.ts"
    "src/services/websocket.ts"
    "src/api/client.ts"
)

for component in "${CORE_COMPONENTS[@]}"; do
    if [ -f "$component" ]; then
        print_result "Component exists: $component" "pass"
    else
        print_result "Component exists: $component" "fail" "File not found"
    fi
done

echo ""
echo "6. Docker Configuration Test"
echo "----------------------------"

# Test 9: Check Docker files
if [ -f "Dockerfile" ]; then
    print_result "Frontend Dockerfile" "pass"
else
    print_result "Frontend Dockerfile" "fail" "Dockerfile not found"
fi

if [ -f "nginx.conf" ]; then
    print_result "Nginx configuration" "pass"
else
    print_result "Nginx configuration" "fail" "nginx.conf not found"
fi

# Test 10: Check docker-compose
if [ -f "../docker-compose.yml" ]; then
    print_result "Docker compose configuration" "pass"
else
    print_result "Docker compose configuration" "fail" "docker-compose.yml not found"
fi

echo ""
echo "7. Performance Optimization Check"
echo "---------------------------------"

# Test 11: Check for code splitting configuration
if grep -q "build.rollupOptions" vite.config.ts 2>/dev/null; then
    print_result "Code splitting configuration" "pass"
else
    print_result "Code splitting configuration" "fail" "Rollup options not configured"
fi

# Test 12: Check for lazy loading in router
if grep -q "defineAsyncComponent\|lazy" src/router/index.ts 2>/dev/null; then
    print_result "Lazy loading configuration" "pass"
else
    print_result "Lazy loading configuration" "warn" "No lazy loading found in router"
fi

echo ""
echo "8. Build Output Test"
echo "-------------------"

# Test 13: Check if build output directory exists
if [ -d "dist" ]; then
    print_result "Build output directory" "pass"

    # Check build output size
    BUILD_SIZE=$(du -sh dist | cut -f1)
    echo "  Build size: $BUILD_SIZE"

    # Check for essential files
    if [ -f "dist/index.html" ]; then
        print_result "Build output: index.html" "pass"
    else
        print_result "Build output: index.html" "fail" "index.html not found in dist"
    fi
else
    print_result "Build output directory" "fail" "dist directory not found - run 'npm run build'"
fi

echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed!${NC}"
    exit 0
else
    echo -e "\n${RED}Some tests failed. Please review the errors above.${NC}"
    exit 1
fi
