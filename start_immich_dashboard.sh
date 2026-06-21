#!/bin/bash
# Immich Dashboard 启动器

echo "🚀 启动 NAS Dashboard Immich 集成..."
echo ""

# 检查服务状态
if curl -s -o /dev/null -w "%{http_code}" http://localhost:2283 | grep -q "200"; then
    echo "✅ Immich服务在线"
    echo "📷 正在打开Immich管理界面..."
    
    # 打开HTML页面
    if command -v xdg-open > /dev/null 2>&1; then
        xdg-open /data/nas-dashboard/immich-launcher.html
    elif command -v open > /dev/null 2>&1; then
        open /data/nas-dashboard/immich-launcher.html
    else
        echo "请手动在浏览器中打开:"
        echo "file:///data/nas-dashboard/immich-launcher.html"
    fi
else
    echo "❌ Immich服务离线"
    echo "请先启动Immich服务"
fi
