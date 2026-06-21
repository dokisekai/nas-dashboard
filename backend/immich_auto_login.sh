#!/bin/bash
# Immich自动登录脚本

CONFIG_FILE="/data/nas-dashboard/backend/config/immich.json"

if [ ! -f "$CONFIG_FILE" ]; then
    echo '{"error": "Immich not configured"}'
    exit 1
fi

# 读取配置
API_KEY=$(jq -r '.apiKey' "$CONFIG_FILE")
IMMICH_URL=$(jq -r '.url' "$CONFIG_FILE")

if [ "$API_KEY" == "null" ] || [ -z "$API_KEY" ]; then
    echo '{"error": "API key not configured"}'
    exit 1
fi

# 生成临时登录token (使用API密钥)
# 检查当前用户信息
current_user=$(curl -s -X GET "$IMMICH_URL/api/user/me" \
    -H "X-API-Key: $API_KEY")

echo "$current_user"
