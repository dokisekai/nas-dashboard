#!/bin/bash
# Docker服务删除脚本
# 删除指定的Docker服务并清理相关资源

echo "🗑️ Docker服务删除脚本"
echo "================================"

# 要删除的服务
SERVICES=("uptime-kuma" "open-webui" "ollama")

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🔍 检查服务状态...${NC}"

for service in "${SERVICES[@]}"; do
    if docker ps -a --format '{{.Names}}' | grep -q "^${service}$"; then
        echo -e "${GREEN}✅${NC} $service - 运行中"
    else
        echo -e "${YELLOW}⚠️ ${NC} $service - 未运行或不存在"
    fi
done

echo ""
echo -e "${YELLOW}⚠️  即将删除以下服务:${NC}"
echo "  - uptime-kuma (监控服务)"
echo "  - open-webui (管理界面)"
echo "  - ollama (AI服务)"
echo ""

read -p "确认删除这些服务？(y/n): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${RED}❌ 删除操作已取消${NC}"
    exit 0
fi

echo -e "${GREEN}🚀 开始删除服务...${NC}"

# 1. 停止并删除容器
echo ""
echo -e "${BLUE}📋 停止并删除容器...${NC}"
for service in "${SERVICES[@]}"; do
    if docker ps -a --format '{{.Names}}' | grep -q "^${service}$"; then
        echo -n "停止 $service..."
        docker stop $service 2>/dev/null
        if [ $? -eq 0 ]; then
            echo -e " ${GREEN}✅${NC}"
        else
            echo -e " ${YELLOW}⚠️ ${NC}"
        fi

        echo -n "删除 $service..."
        docker rm $service 2>/dev/null
        if [ $? -eq 0 ]; then
            echo -e " ${GREEN}✅${NC}"
        else
            echo -e " ${YELLOW}⚠️ ${NC}"
        fi
    else
        echo -e "$service - ${YELLOW}不存在，跳过${NC}"
    fi
done

# 2. 清理相关镜像
echo ""
echo -e "${BLUE}🗑️  清理相关镜像...${NC}"
IMAGES_TO_REMOVE=(
    "louislam/uptime-kuma:1"
    "ghcr.io/open-webui/open-webui:main"
    "ollama/ollama:latest"
)

for image in "${IMAGES_TO_REMOVE[@]}"; do
    echo -n "删除镜像 $image..."
    docker rmi $image 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e " ${GREEN}✅${NC}"
    else
        echo -e " ${YELLOW}⚠️  未找到镜像${NC}"
    fi
done

# 3. 清理数据卷（可选）
echo ""
read -p "是否要删除相关数据卷？(这会删除所有数据！)(y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}💾 清理数据卷...${NC}"

    # Uptime Kuma数据卷
    if docker volume ls | grep -q "uptime-kuma"; then
        docker volume rm uptime-kuma 2>/dev/null
        echo -e "${GREEN}✅ 删除 uptime-kuma 数据卷${NC}"
    fi

    # Open WebUI数据卷
    if docker volume ls | grep -q "open-webui"; then
        docker volume rm open-webui 2>/dev/null
        echo -e "${GREEN}✅ 删除 open-webui 数据卷${NC}"
    fi

    # Ollama数据卷
    if docker volume ls | grep -q "ollama"; then
        docker volume rm ollama 2>/dev/null
        echo -e "${GREEN}✅ 删除 ollama 数据卷${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  数据卷保留，未删除${NC}"
fi

# 4. 清理网络（如果没有其他服务使用）
echo ""
echo -e "${BLUE}🌐 检查网络...${NC}"
NETWORKS_IN_USE=$(docker network ls --format '{{.Name}}' | wc -l)
echo "当前使用中的网络数量: $NETWORKS_IN_USE"

# 如果删除的服务有专属网络，可以清理
echo -e "${GREEN}✅ 网络检查完成${NC}"

# 5. 清理配置文件（可选）
echo ""
read -p "是否要删除配置文件和端口映射？(y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${BLUE}📁 清理配置文件...${NC}"

    # 检查并删除配置目录
    CONFIG_DIRS=(
        "/opt/uptime-kuma"
        "/opt/open-webui"
        "/opt/ollama"
    )

    for config_dir in "${CONFIG_DIRS[@]}"; do
        if [ -d "$config_dir" ]; then
            echo -n "删除配置目录 $config_dir..."
            rm -rf "$config_dir" 2>/dev/null
            if [ $? -eq 0 ]; then
                echo -e " ${GREEN}✅${NC}"
            else
                echo -e " ${YELLOW}⚠️${NC}"
            fi
        else
            echo "$config_dir - ${YELLOW}不存在${NC}"
        fi
    done
else
    echo -e "${YELLOW}⚠️  配置文件保留${NC}"
fi

# 6. 检查删除结果
echo ""
echo -e "${GREEN}✅ 删除完成！${NC}"
echo ""
echo "📊 删除总结:"
for service in "${SERVICES[@]}"; do
    if docker ps -a --format '{{.Names}}' | grep -q "^${service}$"; then
        echo -e "  $service - ${RED}❌ 删除失败${NC}"
    else
        echo -e "  $service - ${GREEN}✅ 已删除${NC}"
    fi
done

# 7. 显示剩余服务
echo ""
echo -e "${BLUE}🐳 当前仍在运行的Docker服务:${NC}"
docker ps --format "table {{.Names}}\t{{.Image}}\t{{.Status}}" | head -10

echo ""
echo -e "${GREEN}🎉 清理完成！${NC}"
echo -e "${YELLOW}💡 提示: 可以使用以下命令查看系统状态:${NC}"
echo "  docker ps"
echo "  docker images"
echo "  docker volume ls"