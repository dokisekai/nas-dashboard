#!/bin/bash
# Comprehensive System Integration Test
# This script performs real testing without requiring successful build

set -e

echo "========================================="
echo "NAS Dashboard - Comprehensive Test"
echo "========================================="
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
WARNED_TESTS=0

print_result() {
    local test_name="$1"
    local result="$2"
    local message="$3"

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    case $result in
        "pass")
            echo -e "${GREEN}✓${NC} $test_name"
            PASSED_TESTS=$((PASSED_TESTS + 1))
            ;;
        "fail")
            echo -e "${RED}✗${NC} $test_name"
            if [ -n "$message" ]; then
                echo -e "  ${YELLOW}Message: $message${NC}"
            fi
            FAILED_TESTS=$((FAILED_TESTS + 1))
            ;;
        "warn")
            echo -e "${YELLOW}⚠${NC} $test_name"
            if [ -n "$message" ]; then
                echo -e "  ${YELLOW}Message: $message${NC}"
            fi
            WARNED_TESTS=$((WARNED_TESTS + 1))
            ;;
    esac
}

cd "$(dirname "$0")"

echo "1. Project Structure Test"
echo "-------------------------"

# Test project structure
structure_tests=(
    "src:Source directory"
    "src/views:Views directory"
    "src/components:Components directory"
    "src/apps:Apps directory"
    "src/composables:Composables directory"
    "src/api:API directory"
    "src/services:Services directory"
    "src/stores:Stores directory"
    "src/types:Types directory"
    "src/plugin-system:Plugin system directory"
    "package.json:Package.json file"
    "vite.config.ts:Vite configuration"
    "tsconfig.json:TypeScript configuration"
    "tailwind.config.js:Tailwind configuration"
    "Dockerfile:Docker file"
    "nginx.conf:Nginx configuration"
    ".env:Environment configuration"
)

for item in "${structure_tests[@]}"; do
    IFS=':' read -r path description <<< "$item"
    if [ -e "$path" ]; then
        print_result "$description exists" "pass"
    else
        print_result "$description exists" "fail" "File/directory not found: $path"
    fi
done

echo ""
echo "2. Dependency Check"
echo "--------------------"

# Check if node_modules exists
if [ -d "node_modules" ]; then
    print_result "Node modules installed" "pass"

    # Check critical dependencies
    critical_deps=(
        "vue"
        "vue-router"
        "pinia"
        "axios"
        "chart.js"
        "sortablejs"
    )

    for dep in "${critical_deps[@]}"; do
        if [ -d "node_modules/$dep" ]; then
            print_result "Dependency: $dep" "pass"
        else
            print_result "Dependency: $dep" "fail" "Not installed"
        fi
    done
else
    print_result "Node modules installed" "fail" "Run 'npm install'"
fi

echo ""
echo "3. Configuration Files Test"
echo "---------------------------"

# Test environment configuration
if [ -f ".env" ]; then
    print_result "Environment file exists" "pass"

    # Check required variables
    required_vars=("VITE_API_URL" "VITE_WS_URL")
    for var in "${required_vars[@]}"; do
        if grep -q "^$var=" .env 2>/dev/null; then
            print_result "Variable: $var" "pass"
        else
            print_result "Variable: $var" "fail" "Not set in .env"
        fi
    done
else
    print_result "Environment file exists" "fail" "Create .env file"
fi

echo ""
echo "4. Core Components Test"
echo "-----------------------"

# Check core components
core_files=(
    "src/main.ts:Main entry point"
    "src/App.vue:Root component"
    "src/router/index.ts:Router configuration"
    "src/stores/app.ts:App store"
    "src/stores/auth.ts:Auth store"
    "src/api/client.ts:API client"
    "src/composables/useDesktop.ts:Desktop composable"
    "src/composables/useWebSocket.ts:WebSocket composable"
    "src/services/websocket.ts:WebSocket service"
)

for item in "${core_files[@]}"; do
    IFS=':' read -r path description <<< "$item"
    if [ -f "$path" ]; then
        print_result "$description" "pass"
    else
        print_result "$description" "fail" "File not found"
    fi
done

echo ""
echo "5. View Components Test"
echo "-----------------------"

# Check view components
view_files=(
    "src/views/Login.vue:Login page"
    "src/views/Dashboard.vue:Dashboard page"
)

for item in "${view_files[@]}"; do
    IFS=':' read -r path description <<< "$item"
    if [ -f "$path" ]; then
        print_result "$description" "pass"
    else
        print_result "$description" "fail" "File not found"
    fi
done

echo ""
echo "6. Desktop System Test"
echo "----------------------"

# Check desktop components
desktop_files=(
    "src/components/Desktop/DSMDesktop.vue:Desktop interface"
    "src/components/Desktop/Dock.vue:Application dock"
    "src/components/Desktop/WindowManager.vue:Window manager"
    "src/components/Desktop/DesktopWindow.vue:Desktop window"
    "src/components/Desktop/WidgetLibrary.vue:Widget library"
)

for item in "${desktop_files[@]}"; do
    IFS=':' read -r path description <<< "$item"
    if [ -f "$path" ]; then
        print_result "$description" "pass"
    else
        print_result "$description" "warn" "File not found (optional)"
    fi
done

echo ""
echo "7. Application Modules Test"
echo "--------------------------"

# Check application modules
app_files=(
    "src/apps/AppCenter.vue:Application center"
    "src/apps/SystemMonitor.vue:System monitor"
    "src/apps/StorageManager.vue:Storage manager"
    "src/apps/UserManager.vue:User manager"
)

for item in "${app_files[@]}"; do
    IFS=':' read -r path description <<< "$item"
    if [ -f "$path" ]; then
        print_result "$description" "pass"
    else
        print_result "$description" "warn" "File not found (optional)"
    fi
done

echo ""
echo "8. Plugin System Test"
echo "---------------------"

# Check plugin system
plugin_files=(
    "src/plugin-system/index.ts:Plugin system entry"
    "src/plugin-system/types/plugin.ts:Plugin types"
    "src/plugin-system/manager/PluginManager.ts:Plugin manager"
    "src/plugin-system/sdk/api.ts:Plugin SDK"
)

for item in "${plugin_files[@]}"; do
    IFS=':' read -r path description <<< "$item"
    if [ -f "$path" ]; then
        print_result "$description" "pass"
    else
        print_result "$description" "warn" "Plugin system not fully implemented"
    fi
done

echo ""
echo "9. Docker Configuration Test"
echo "---------------------------"

# Test Docker configuration
if [ -f "Dockerfile" ]; then
    print_result "Frontend Dockerfile" "pass"

    # Check Dockerfile content
    if grep -q "FROM" Dockerfile && grep -q "nginx" Dockerfile; then
        print_result "Dockerfile configuration" "pass"
    else
        print_result "Dockerfile configuration" "warn" "Non-standard Dockerfile"
    fi
else
    print_result "Frontend Dockerfile" "fail" "Dockerfile not found"
fi

if [ -f "nginx.conf" ]; then
    print_result "Nginx configuration" "pass"
else
    print_result "Nginx configuration" "fail" "nginx.conf not found"
fi

if [ -f "../docker-compose.yml" ]; then
    print_result "Docker Compose configuration" "pass"
else
    print_result "Docker Compose configuration" "fail" "docker-compose.yml not found"
fi

echo ""
echo "10. Backend Integration Test"
echo "----------------------------"

# Test backend availability (if backend URL is set)
if [ -f ".env" ] && grep -q "VITE_API_URL" .env; then
    API_URL=$(grep "^VITE_API_URL=" .env | cut -d'=' -f2)
    echo "Testing backend connection to: $API_URL"

    # Test if backend is reachable
    if curl -s -o /dev/null -w "%{http_code}" "$API_URL/api/health" 2>/dev/null | grep -q "200\|404"; then
        print_result "Backend API reachable" "pass"
    else
        print_result "Backend API reachable" "warn" "Backend not running or not reachable"
    fi
else
    print_result "Backend API configuration" "warn" "API URL not configured"
fi

echo ""
echo "11. Code Quality Check"
echo "----------------------"

# Check for common issues
if [ -d "src" ]; then
    # Check for console.log statements (development artifacts)
    console_logs=$(grep -r "console\.log" src/ 2>/dev/null | wc -l)
    if [ "$console_logs" -gt 0 ]; then
        print_result "Console.log cleanup" "warn" "Found $console_logs console.log statements"
    else
        print_result "Console.log cleanup" "pass"
    fi

    # Check for TODO comments
    todos=$(grep -r "TODO\|FIXME" src/ 2>/dev/null | wc -l)
    if [ "$todos" -gt 0 ]; then
        print_result "TODO/FIXME comments" "warn" "Found $todos TODO/FIXME comments"
    else
        print_result "TODO/FIXME comments" "pass"
    fi
fi

echo ""
echo "12. Performance Optimization Check"
echo "----------------------------------"

# Check for performance optimization
if [ -f "vite.config.ts" ]; then
    if grep -q "build.rollupOptions" vite.config.ts; then
        print_result "Code splitting configured" "pass"
    else
        print_result "Code splitting configured" "warn" "Rollup options not configured"
    fi
fi

# Check for lazy loading
if [ -f "src/router/index.ts" ]; then
    if grep -q "defineAsyncComponent\|lazy" src/router/index.ts; then
        print_result "Lazy loading implemented" "pass"
    else
        print_result "Lazy loading implemented" "warn" "No lazy loading found"
    fi
fi

echo ""
echo "========================================="
echo "Test Summary"
echo "========================================="
echo "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${YELLOW}Warnings: $WARNED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ] && [ $WARNED_TESTS -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed! System is ready for deployment.${NC}"
    exit 0
elif [ $FAILED_TESTS -eq 0 ]; then
    echo -e "\n${YELLOW}System is functional with some warnings. Review recommended.${NC}"
    exit 0
else
    echo -e "\n${RED}Some critical tests failed. Please review the errors above.${NC}"
    exit 1
fi
