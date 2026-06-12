#!/bin/bash

# NAS Dashboard - Critical Fixes Testing Script
# This script validates all the critical fixes implemented

echo "=========================================="
echo "NAS Dashboard - Critical Fixes Testing"
echo "=========================================="
echo ""

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_PASSED=0
TESTS_FAILED=0

# Function to check if file exists
check_file() {
    if [ -f "$1" ]; then
        echo -e "${GREEN}✓${NC} File exists: $1"
        ((TESTS_PASSED++))
        return 0
    else
        echo -e "${RED}✗${NC} File missing: $1"
        ((TESTS_FAILED++))
        return 1
    fi
}

# Function to check if string exists in file
check_content() {
    local file="$1"
    local content="$2"
    local description="$3"

    if grep -q "$content" "$file"; then
        echo -e "${GREEN}✓${NC} $description"
        ((TESTS_PASSED++))
        return 0
    else
        echo -e "${RED}✗${NC} $description"
        ((TESTS_FAILED++))
        return 1
    fi
}

echo "1. Testing File Structure..."
echo "----------------------------"

# Check critical files exist
check_file "/home/hserver/nas-dashboard/frontend/src/App.vue"
check_file "/home/hserver/nas-dashboard/frontend/src/components/Layout/Main.vue"
check_file "/home/hserver/nas-dashboard/frontend/src/stores/auth.ts"
check_file "/home/hserver/nas-dashboard/frontend/src/api/client.ts"
check_file "/home/hserver/nas-dashboard/frontend/.env"
check_file "/home/hserver/nas-dashboard/frontend/src/views/Login.vue"
check_file "/home/hserver/nas-dashboard/frontend/src/views/Dashboard.vue"
check_file "/home/hserver/nas-dashboard/frontend/src/views/Monitor/CPU.vue"

echo ""
echo "2. Testing Layout Duplication Fix..."
echo "------------------------------------"

# Check App.vue doesn't have nested router-view in MainLayout
if ! grep -q "MainLayout v-else>" /home/hserver/nas-dashboard/frontend/src/App.vue; then
    echo -e "${YELLOW}⚠${NC} App.vue structure changed - manual review needed"
else
    if grep -q "<MainLayout v-else />" /home/hserver/nas-dashboard/frontend/src/App.vue; then
        echo -e "${GREEN}✓${NC} App.vue - No nested router-view in MainLayout"
        ((TESTS_PASSED++))
    else
        echo -e "${RED}✗${NC} App.vue - Still has nested structure"
        ((TESTS_FAILED++))
    fi
fi

# Check MainLayout uses router-view instead of slot
check_content "/home/hserver/nas-dashboard/frontend/src/components/Layout/Main.vue" \
    "<router-view />" \
    "MainLayout uses router-view instead of slot"

echo ""
echo "3. Testing Authentication Store Fix..."
echo "---------------------------------------"

# Check auth store has setUser method
check_content "/home/hserver/nas-dashboard/frontend/src/stores/auth.ts" \
    "setUser" \
    "Auth store has setUser method"

# Check Login.vue sets user data
check_content "/home/hserver/nas-dashboard/frontend/src/views/Login.vue" \
    "setUser" \
    "Login.vue sets user data"

echo ""
echo "4. Testing API Client Configuration..."
echo "--------------------------------------"

# Check debug mode configuration
check_content "/home/hserver/nas-dashboard/frontend/src/api/client.ts" \
    "VITE_DEBUG" \
    "API client has debug mode support"

# Check proper error handling with user data cleanup
check_content "/home/hserver/nas-dashboard/frontend/src/api/client.ts" \
    "localStorage.removeItem('user')" \
    "API client cleans user data on 401"

# Check environment configuration
check_content "/home/hserver/nas-dashboard/frontend/.env" \
    "VITE_DEBUG=false" \
    "Environment file has debug configuration"

echo ""
echo "5. Testing Header Component Fix..."
echo "----------------------------------"

# Check Header uses auth store user data
check_content "/home/hserver/nas-dashboard/frontend/src/components/Layout/Header.vue" \
    "authStore.user?.username" \
    "Header displays dynamic username from auth store"

echo ""
echo "6. Testing Network Data Format Fix..."
echo "--------------------------------------"

# Check Dashboard.vue handles new network data structure
check_content "/home/hserver/nas-dashboard/frontend/src/views/Dashboard.vue" \
    "bytesRecv" \
    "Dashboard handles bytesRecv field"
check_content "/home/hserver/nas-dashboard/frontend/src/views/Dashboard.vue" \
    "bytesSent" \
    "Dashboard handles bytesSent field"
check_content "/home/hserver/nas-dashboard/frontend/src/views/Dashboard.vue" \
    "!i.name.startsWith('docker')" \
    "Dashboard filters docker interfaces"

echo ""
echo "7. Testing CPU Display Fix..."
echo "----------------------------"

# Check CPU monitor converts usage to percentage
check_content "/home/hserver/nas-dashboard/frontend/src/views/Monitor/CPU.vue" \
    "(cpuInfo.usage * 100)" \
    "CPU monitor converts usage to percentage"
check_content "/home/hserver/nas-dashboard/frontend/src/views/Monitor/CPU.vue" \
    "percent = usage > 1 ? usage : usage * 100" \
    "CPU status functions handle decimal/percentage"

echo ""
echo "=========================================="
echo "Test Results Summary"
echo "=========================================="
echo -e "Tests Passed: ${GREEN}${TESTS_PASSED}${NC}"
echo -e "Tests Failed: ${RED}${TESTS_FAILED}${NC}"
echo ""

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All critical fixes have been validated!${NC}"
    exit 0
else
    echo -e "${RED}✗ Some fixes need attention. Please review the failures above.${NC}"
    exit 1
fi
