#!/bin/bash
# 安全的Immich启动器

echo "🚀 启动 NAS Dashboard Immich 安全集成..."
echo ""

IMMICH_URL="http://localhost:2283"

# 检查服务状态
echo "🔍 检查Immich服务状态..."
if curl -s -o /dev/null -w "%{http_code}" "$IMMICH_URL" | grep -q "200"; then
    echo "✅ Immich服务在线"
    echo "📷 正在打开管理界面..."
    
    # 在浏览器中打开安全的启动器
    launcher_path="/data/nas-dashboard/immich_secure_launcher.html"
    
    if command -v xdg-open > /dev/null 2>&1; then
        xdg-open "$launcher_path"
    elif command -v open > /dev/null 2>&1; then
        open "$launcher_path"
    else
        echo "💡 请在浏览器中手动打开:"
        echo "   file://$launcher_path"
    fi
else
    echo "❌ Immich服务离线"
    echo "💡 请确保Immich服务运行在 $IMMICH_URL"
fi
