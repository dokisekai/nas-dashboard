#!/bin/bash
# 快速Docker服务删除脚本

echo "🗑️ 准备删除以下Docker服务："
echo "  1. uptime-kuma (监控服务)"
echo " 2. open-webui (管理界面)"
echo " 3. ollama (AI服务)"
echo ""

# 检查服务状态
echo "🔍 当前服务状态："
for service in "uptime-kuma" "open-webui" "ollama"; do
    if docker ps -a --format '{{.Names}}\t{{.Status}}' | grep -q "$service"; then
        status=$(docker ps -a --format '{{.Names}}\t{{.Status}}' | grep "$service" | awk '{print $2}')
        echo "  ✅ $service - $status"
    else
        echo "  ⚠️  $service - 未运行"
    fi
done

echo ""
read -p "确认删除这些服务？(y/n): " -n 1 -r
echo ""

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 操作已取消"
    exit 0
fi

echo "🚀 开始删除..."
echo ""

# 停止并删除容器
for service in "uptime-kuma" "open-webui" "ollama"; do
    echo "🗑️  删除 $service..."

    if docker ps -q $service; then
        docker stop $service 2>/dev/null
        docker rm $service 2>/dev/null
        echo "  ✅ 容器已删除"
    else
        echo "  ⚠️  容器未运行"
    fi

    # 删除镜像
    if [[ "$service" == "uptime-kuma" ]]; then
        docker rmi louislam/uptime-kuma:1 2>/dev/null
    elif [[ "$service" == "open-webui" ]]; then
        docker rmi ghcr.io/open-webui/open-webui:main 2>/dev/null
    elif [[ "$service" == "ollama" ]]; then
        docker rmi ollama/ollama:latest 2>/dev/null
    fi
    echo "  ✅ 镜像已清理"
    echo ""
done

echo "✅ 删除完成！"
echo ""
echo "📊 剩余运行的服务："
docker ps --format "table {{.Names}}\t{{.Image}}" | head -10