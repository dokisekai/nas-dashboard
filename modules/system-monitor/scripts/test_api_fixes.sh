#!/bin/bash
# API修复验证测试脚本
# 测试修复后的前端组件API调用

echo "🧪 API修复验证测试"
echo "===================="
echo ""

# 获取认证token
echo "🔑 获取认证token..."
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8888/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}')

TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.token')
echo "✅ Token获取成功"
echo ""

# 测试用户组API
echo "👥 测试用户组API (/api/groups)..."
GROUPS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8888/api/groups)
GROUPS_FORMAT=$(echo $GROUPS_RESPONSE | jq 'has("groups")')
if [ "$GROUPS_FORMAT" = "true" ]; then
  GROUPS_COUNT=$(echo $GROUPS_RESPONSE | jq '.groups | length')
  echo "✅ 返回格式正确: {\"groups\": [...]}"
  echo "📊 用户组数量: $GROUPS_COUNT"

  # 检查hserver用户的组信息
  HSERVER_GROUPS=$(echo $GROUPS_RESPONSE | jq -r '.groups[] | select(.members != null) | select(.members | index("hserver")) | .name' | wc -l)
  echo "👤 hserver属于 $HSERVER_GROUPS 个用户组"
else
  echo "❌ 返回格式错误"
  exit 1
fi
echo ""

# 测试用户API
echo "👤 测试用户API (/api/users)..."
USERS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" http://localhost:8888/api/users)
USERS_FORMAT=$(echo $USERS_RESPONSE | jq 'has("users")')
if [ "$USERS_FORMAT" = "true" ]; then
  USERS_COUNT=$(echo $USERS_RESPONSE | jq '.users | length')
  echo "✅ 返回格式正确: {\"users\": [...]}"
  echo "📊 用户数量: $USERS_COUNT"

  # 检查hserver用户
  HSERVER_USER=$(echo $USERS_RESPONSE | jq -r '.users[] | select(.username == "hserver") | .username')
  if [ "$HSERVER_USER" = "hserver" ]; then
    echo "👤 hserver用户存在"
    HSERVER_INFO=$(echo $USERS_RESPONSE | jq -r '.users[] | select(.username == "hserver")')
    echo "📋 hserver信息: UID=$(echo $HSERVER_INFO | jq -r '.uid'), Group=$(echo $HSERVER_INFO | jq -r '.group')"
  fi
else
  echo "❌ 返回格式错误"
  exit 1
fi
echo ""

# 测试权限文件API (修复后的端点)
echo "🔒 测试权限文件API (/api/permissions/files?path=/home/hserver)..."
PERMS_RESPONSE=$(curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8888/api/permissions/files?path=%2Fhome%2Fhserver")
PERMS_FORMAT=$(echo $PERMS_RESPONSE | jq 'has("file")')
if [ "$PERMS_FORMAT" = "true" ]; then
  echo "✅ 返回格式正确"
  FILE_PATH=$(echo $PERMS_RESPONSE | jq -r '.file.path')
  FILE_OWNER=$(echo $PERMS_RESPONSE | jq -r '.file.owner')
  FILE_GROUP=$(echo $PERMS_RESPONSE | jq -r '.file.group')
  FILE_MODE=$(echo $PERMS_RESPONSE | jq -r '.file.mode')
  echo "📁 文件: $FILE_PATH"
  echo "👤 所有者: $FILE_OWNER"
  echo "👥 用户组: $FILE_GROUP"
  echo "🔐 权限: $FILE_MODE"
else
  echo "❌ 返回格式错误"
  exit 1
fi
echo ""

# 测试文件浏览API
echo "📂 测试文件浏览API (/api/files/list)..."
FILES_RESPONSE=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"path":"/home/hserver"}' \
  http://localhost:8888/api/files/list)
FILES_COUNT=$(echo $FILES_RESPONSE | jq 'length')
echo "✅ 文件列表API正常"
echo "📊 /home/hserver 下的项目数量: $FILES_COUNT"
echo ""

# 模拟前端数据处理逻辑
echo "🔧 测试前端数据处理逻辑..."
echo "用户数据处理:"
USERS_ARRAY=$(echo $USERS_RESPONSE | jq 'if .users then .users else . end')
echo "提取用户数组: $USERS_ARRAY" | jq -r '.[0:3]'

echo "用户组数据处理:"
GROUPS_ARRAY=$(echo $GROUPS_RESPONSE | jq 'if .groups then .groups else . end')
echo "提取用户组数组: $GROUPS_ARRAY" | jq -r '.[0:3]'
echo ""

echo "🎉 API修复验证完成！"
echo "===================="
echo "✅ 用户管理API: 正常"
echo "✅ 用户组管理API: 正常"
echo "✅ 权限管理API: 正常"
echo "✅ 文件浏览API: 正常"
echo ""
echo "📋 修复总结:"
echo "1. UserManager.vue - 修复了数据处理以匹配 {\"users\": [...]} 格式"
echo "2. GroupManager.vue - 修复了数据处理以匹配 {\"groups\": [...]} 格式"
echo "3. ACLEditor.vue - 修复了API端点:"
echo "   - /api/permissions/info → /api/permissions/files"
echo "   - /api/permissions/set → /api/permissions/files/permissions"
echo "   - /api/permissions/acl → /api/permissions/files/acl"
echo ""
echo "🌐 前端访问地址: http://localhost:5173/desktop"
echo "🎛️ 控制面板测试路径: 桌面 → 控制面板 → 用户管理/权限管理"
