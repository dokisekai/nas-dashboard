#!/bin/bash
# Immich用户管理脚本 - 直接API访问版本

API_KEY="t5nmlDaFlyz7guxpfmAwvOeiKp7zEC7loF2Ow15V8Q"
API_URL="http://localhost:2283/api"

echo "🎯 Immich用户管理工具"
echo "======================================"
echo ""

# 检查API连接
echo "🔍 检查Immich API连接..."
check_response=$(curl -s -X GET "$API_URL/users" -H "X-API-Key: $API_KEY")
if echo "$check_response" | grep -q "email\|name\|id"; then
    echo "✅ API连接正常"
else
    echo "❌ API连接失败"
    echo "响应: $check_response"
    exit 1
fi
echo ""

# 显示当前用户
echo "👥 当前Immich用户:"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
curl -s -X GET "$API_URL/users" -H "X-API-Key: $API_KEY" | python3 -c "
import sys, json
data = json.load(sys.stdin)
if isinstance(data, list):
    for i, user in enumerate(data):
        print(f'{i+1}. {user.get(\"name\", \"N/A\")}')
        print(f'   邮箱: {user.get(\"email\", \"N/A\")}')
        print(f'   ID: {user.get(\"id\", \"N/A\")}')
        print(f'   照片数: {user.get(\"photos\", 0)}')
        print(f'   存储空间: {user.get(\"storageSize\", 0) / (1024**3):.2f} GB' if 'storageSize' in user else '   存储空间: N/A')
        print()
else:
    print('无法解析用户数据')
"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 交互式菜单
while true; do
    echo "🎛️  请选择操作:"
    echo "1. 查看用户列表"
    echo "2. 创建新用户"
    echo "3. 删除用户"
    echo "4. 查看用户详情"
    echo "5. 退出"
    echo ""
    read -p "请输入选项 (1-5): " choice

    case $choice in
        1)
            echo ""
            echo "📋 用户列表:"
            curl -s -X GET "$API_URL/users" -H "X-API-Key: $API_KEY" | python3 -c "
import sys, json
data = json.load(sys.stdin)
if isinstance(data, list):
    for i, user in enumerate(data):
        print(f'{i+1}. {user.get(\"name\", \"N/A\")} ({user.get(\"email\", \"N/A\")})')
"
            echo ""
            ;;
        2)
            echo ""
            echo "➕ 创建新用户"
            read -p "用户姓名: " name
            read -p "邮箱地址: " email
            read -p "密码: " password

            if [ -z "$name" ] || [ -z "$email" ] || [ -z "$password" ]; then
                echo "❌ 所有字段都是必填的"
                continue
            fi

            create_response=$(curl -s -X POST "$API_URL/users" \
                -H "X-API-Key: $API_KEY" \
                -H "Content-Type: application/json" \
                -d "{
                    \"name\": \"$name\",
                    \"email\": \"$email\",
                    \"password\": \"$password\"
                }")

            if echo "$create_response" | grep -q "id\|email"; then
                echo "✅ 用户创建成功"
                echo "响应: $create_response" | python3 -m json.tool 2>/dev/null || echo "$create_response"
            else
                echo "❌ 用户创建失败"
                echo "响应: $create_response"
            fi
            echo ""
            ;;
        3)
            echo ""
            echo "🗑️  删除用户"

            # 先显示用户列表
            echo "当前用户:"
            curl -s -X GET "$API_URL/users" -H "X-API-Key: $API_KEY" | python3 -c "
import sys, json
data = json.load(sys.stdin)
if isinstance(data, list):
    for i, user in enumerate(data):
        print(f'{i+1}. {user.get(\"name\", \"N/A\")} (ID: {user.get(\"id\", \"N/A\")[:8]}...)')
"

            read -p "输入要删除的用户ID: " user_id

            if [ -z "$user_id" ]; then
                echo "❌ 用户ID不能为空"
                continue
            fi

            delete_response=$(curl -s -X DELETE "$API_URL/users/$user_id" \
                -H "X-API-Key: $API_KEY")

            if echo "$delete_response" | grep -q "204\|200\|deleted\|success"; then
                echo "✅ 用户删除成功"
            else
                echo "❌ 用户删除失败"
                echo "响应: $delete_response"
            fi
            echo ""
            ;;
        4)
            echo ""
            echo "👤 查看用户详情"
            read -p "输入用户ID: " user_id

            if [ -z "$user_id" ]; then
                echo "❌ 用户ID不能为空"
                continue
            fi

            user_detail=$(curl -s -X GET "$API_URL/users/$user_id" \
                -H "X-API-Key: $API_KEY")

            echo "用户详情:"
            echo "$user_detail" | python3 -m json.tool 2>/dev/null || echo "$user_detail"
            echo ""
            ;;
        5)
            echo "👋 再见！"
            exit 0
            ;;
        *)
            echo "❌ 无效选项"
            ;;
    esac
done