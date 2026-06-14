#!/bin/bash

# 应用包模板生成工具
# 用法: ./create-app-template.sh <应用名称>

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 显示帮助信息
show_help() {
    echo "用法: $0 [选项] <应用名称>"
    echo ""
    echo "选项:"
    echo "  -t, --type TYPE        应用类型 (docker|nodejs|python|golang|script)"
    echo "  -p, --port PORT        应用端口"
    echo "  -c, --category CAT     应用分类 (media|productivity|utilities|security|network)"
    echo "  -d, --description DESC 应用描述"
    echo "  -a, --author AUTHOR    作者名称"
    echo "  -w, --website URL      官方网站"
    echo "  -r, --requirement DEP  依赖项（可多个）"
    echo "  -h, --help             显示帮助信息"
    echo ""
    echo "示例:"
    echo "  $0 my-app"
    echo "  $0 -t docker -p 8080 -c utilities portainer"
    echo "  $0 --type nodejs --port 3000 --description '我的Node.js应用' myapp"
}

# 默认值
APP_TYPE="script"
APP_PORT=""
CATEGORY="utilities"
DESCRIPTION=""
AUTHOR="Your Name"
WEBSITE=""
REQUIREMENTS=()

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -t|--type)
            APP_TYPE="$2"
            shift 2
            ;;
        -p|--port)
            APP_PORT="$2"
            shift 2
            ;;
        -c|--category)
            CATEGORY="$2"
            shift 2
            ;;
        -d|--description)
            DESCRIPTION="$2"
            shift 2
            ;;
        -a|--author)
            AUTHOR="$2"
            shift 2
            ;;
        -w|--website)
            WEBSITE="$2"
            shift 2
            ;;
        -r|--requirement)
            REQUIREMENTS+=("$2")
            shift 2
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            APP_NAME="$1"
            shift
            ;;
    esac
done

# 检查应用名称
if [ -z "$APP_NAME" ]; then
    echo -e "${RED}错误: 请提供应用名称${NC}"
    show_help
    exit 1
fi

# 验证应用名称格式（小写字母、数字、连字符）
if ! [[ "$APP_NAME" =~ ^[a-z0-9-]+$ ]]; then
    echo -e "${RED}错误: 应用名称只能包含小写字母、数字和连字符${NC}"
    exit 1
fi

# 验证应用类型
if ! [[ "$APP_TYPE" =~ ^(docker|nodejs|python|golang|script)$ ]]; then
    echo -e "${RED}错误: 不支持的应用类型 '$APP_TYPE'${NC}"
    exit 1
fi

# 验证分类
if ! [[ "$CATEGORY" =~ ^(media|productivity|utilities|security|network)$ ]]; then
    echo -e "${RED}错误: 不支持的应用分类 '$CATEGORY'${NC}"
    exit 1
fi

echo -e "${BLUE}开始创建应用包模板: $APP_NAME${NC}"
echo ""

# 创建目录结构
APP_DIR="$APP_NAME"
mkdir -p "$APP_DIR"/{application,icons,scripts,config,wizard,templates}

# 生成INFO文件
echo -e "${GREEN}1. 生成INFO文件...${NC}"
cat > "$APP_DIR/INFO" << EOF
[package]
name=$APP_NAME
version=1.0.0
displayname=$(echo "$APP_NAME" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1')
description=${DESCRIPTION:-"$(echo "$APP_NAME" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1') application"}
author=$AUTHOR
website=${WEBSITE:-"https://example.com"}
category=$CATEGORY
license=MIT

[system]
architecture=x86_64
min_os_version=1.0.0
min_ram=128
min_disk_space=100
dependencies=${REQUIREMENTS[*]:-""}

[installation]
install_path=/var/packages/$APP_NAME
data_path=/var/packages/$APP_NAME/target
auto_start=true

[permissions]
network_access=true
storage_access=true
port_bindings=${APP_PORT:-"8080"}
EOF

# 根据应用类型生成启动脚本
echo -e "${GREEN}2. 生成启动脚本...${NC}"
case "$APP_TYPE" in
    docker)
        cat > "$APP_DIR/scripts/start.sh" << EOF
#!/bin/bash
set -e

APP_NAME="$APP_NAME"
IMAGE_NAME="$APP_NAME:latest"
CONTAINER_NAME="\${APP_NAME}"

# 检查容器是否已存在
if docker ps -a -f name="\$CONTAINER_NAME" | grep -q .; then
    docker start "\$CONTAINER_NAME"
else
    docker run -d \\
        --name "\$CONTAINER_NAME" \\
        --restart unless-stopped \\
        ${APP_PORT:+-p $APP_PORT:$APP_PORT} \\
        -v "/var/packages/\$APP_NAME/data":/data \\
        -e TZ=Asia/Shanghai \\
        "\$IMAGE_NAME"
fi

echo "应用启动成功"
exit 0
EOF
        ;;

    nodejs)
        cat > "$APP_DIR/scripts/start.sh" << EOF
#!/bin/bash
set -e

APP_NAME="$APP_NAME"
INSTALL_DIR="/var/packages/\$APP_NAME"
APP_EXEC="\$INSTALL_DIR/application/app.js"
NODE_ARGS="--production"
${APP_PORT:+PORT=$APP_PORT}

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: Node.js未安装"
    exit 1
fi

# 启动应用
cd "\$INSTALL_DIR"
nohup node \$NODE_ARGS "\$APP_EXEC" \\
    ${APP_PORT:+--port=\$PORT} \\
    >> "\$INSTALL_DIR/logs/app.log" 2>&1 &

echo \$! > "\$INSTALL_DIR/\$APP_NAME.pid"
echo "应用启动成功"
exit 0
EOF
        ;;

    python)
        cat > "$APP_DIR/scripts/start.sh" << EOF
#!/bin/bash
set -e

APP_NAME="$APP_NAME"
INSTALL_DIR="/var/packages/\$APP_NAME"
APP_EXEC="\$INSTALL_DIR/application/main.py"
${APP_PORT:+APP_PORT=$APP_PORT}

# 检查Python是否安装
if ! command -v python3 &> /dev/null; then
    echo "错误: Python3未安装"
    exit 1
fi

# 启动应用
cd "\$INSTALL_DIR"
nohup python3 "\$APP_EXEC" \\
    ${APP_PORT:+--port=\$APP_PORT} \\
    >> "\$INSTALL_DIR/logs/app.log" 2>&1 &

echo \$! > "\$INSTALL_DIR/\$APP_NAME.pid"
echo "应用启动成功"
exit 0
EOF
        ;;

    golang)
        cat > "$APP_DIR/scripts/start.sh" << EOF
#!/bin/bash
set -e

APP_NAME="$APP_NAME"
INSTALL_DIR="/var/packages/\$APP_NAME"
APP_EXEC="\$INSTALL_DIR/application/\$APP_NAME"
${APP_PORT+:APP_PORT=$APP_PORT}

# 启动应用
cd "\$INSTALL_DIR"
nohup "\$APP_EXEC" \\
    ${APP_PORT:+--port=\$APP_PORT} \\
    >> "\$INSTALL_DIR/logs/app.log" 2>&1 &

echo \$! > "\$INSTALL_DIR/\$APP_NAME.pid"
echo "应用启动成功"
exit 0
EOF
        ;;

    script)
        cat > "$APP_DIR/scripts/start.sh" << EOF
#!/bin/bash
set -e

APP_NAME="$APP_NAME"
INSTALL_DIR="/var/packages/\$APP_NAME"
APP_SCRIPT="\$INSTALL_DIR/application/run.sh"

# 检查脚本是否存在
if [ ! -f "\$APP_SCRIPT" ]; then
    echo "错误: 应用脚本不存在"
    exit 1
fi

# 执行应用脚本
chmod +x "\$APP_SCRIPT"
"\$APP_SCRIPT" &
echo \$! > "\$INSTALL_DIR/\$APP_NAME.pid"

echo "应用启动成功"
exit 0
EOF
        ;;
esac

chmod +x "$APP_DIR/scripts/start.sh"

# 生成停止脚本
echo -e "${GREEN}3. 生成停止脚本...${NC}"
cat > "$APP_DIR/scripts/stop.sh" << 'EOF'
#!/bin/bash
set -e

APP_NAME="__APP_NAME__"
INSTALL_DIR="/var/packages/$APP_NAME"
PID_FILE="$INSTALL_DIR/$APP_NAME.pid"

# 停止Docker容器
if docker ps -q -f name="$APP_NAME" | grep -q .; then
    docker stop "$APP_NAME"
    echo "Docker容器已停止"
    exit 0
fi

# 停止进程
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        kill "$PID" 2>/dev/null || true
        sleep 2

        # 如果还在运行，强制杀死
        if ps -p "$PID" > /dev/null 2>&1; then
            kill -9 "$PID" 2>/dev/null || true
        fi

        rm -f "$PID_FILE"
        echo "应用已停止"
        exit 0
    fi
fi

echo "应用未运行"
exit 0
EOF

sed -i "s/__APP_NAME__/$APP_NAME/g" "$APP_DIR/scripts/stop.sh"
chmod +x "$APP_DIR/scripts/stop.sh"

# 生成状态脚本
echo -e "${GREEN}4. 生成状态脚本...${NC}"
cat > "$APP_DIR/scripts/status.sh" << 'EOF'
#!/bin/bash
set -e

APP_NAME="__APP_NAME__"
INSTALL_DIR="/var/packages/$APP_NAME"
PID_FILE="$INSTALL_DIR/$APP_NAME.pid"

# 检查Docker容器
if docker ps -q -f name="$APP_NAME" -f status=running | grep -q .; then
    echo "running"
    exit 0
fi

# 检查进程
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "running"
        exit 0
    fi
fi

echo "stopped"
exit 1
EOF

sed -i "s/__APP_NAME__/$APP_NAME/g" "$APP_DIR/scripts/status.sh"
chmod +x "$APP_DIR/scripts/status.sh"

# 生成安装脚本
echo -e "${GREEN}5. 生成安装脚本...${NC}"
cat > "$APP_DIR/scripts/installer.sh" << 'EOF'
#!/bin/bash
set -e

APP_NAME="__APP_NAME__"
INSTALL_PATH="__INSTALL_PATH__"
DATA_PATH="__DATA_PATH__"
CONFIG_PATH="__CONFIG_PATH__"

echo "开始安装 $APP_NAME..."

# 创建必要目录
mkdir -p "$DATA_PATH"
mkdir -p "$CONFIG_PATH"
mkdir -p "$INSTALL_PATH/logs"
mkdir -p "$INSTALL_PATH/tmp"

# 设置权限
chown -R nobody:nogroup "$DATA_PATH"
chmod -R 755 "$INSTALL_PATH"

echo "安装完成！"
exit 0
EOF

sed -i "s/__APP_NAME__/$APP_NAME/g" "$APP_DIR/scripts/installer.sh"
sed -i "s|__INSTALL_PATH__|/var/packages/$APP_NAME|g" "$APP_DIR/scripts/installer.sh"
sed -i "s|__DATA_PATH__|/var/packages/$APP_NAME/target|g" "$APP_DIR/scripts/installer.sh"
sed -i "s|__CONFIG_PATH__|/var/packages/$APP_NAME/config|g" "$APP_DIR/scripts/installer.sh"
chmod +x "$APP_DIR/scripts/installer.sh"

# 生成配置文件
echo -e "${GREEN}6. 生成配置文件...${NC}"
cat > "$APP_DIR/config/default_config.json" << EOF
{
  "port": ${APP_PORT:-8080},
  "host": "0.0.0.0",
  "environment": "production",
  "logLevel": "info"
}
EOF

# 生成资源配置
cat > "$APP_DIR/config/resources.json" << EOF
{
  "maxMemoryMB": 512,
  "maxCPU": 2,
  "maxDiskGB": 10,
  "networkAccess": true,
  "storageAccess": true,
  "processAccess": false,
  "portBindings": [${APP_PORT:-8080}]
}
EOF

# 生成README
echo -e "${GREEN}7. 生成README文件...${NC}"
cat > "$APP_DIR/README.md" << EOF
# $APP_NAME Application Package

## 应用说明
$(echo "$APP_NAME" | sed 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) substr($i,2)}1') - $DESCRIPTION

## 安装说明
1. 上传此应用包到NAS Dashboard应用中心
2. 点击安装按钮
3. 等待安装完成
4. 启动应用

## 使用说明
- 应用端口: ${APP_PORT:-8080}
- 数据路径: /var/packages/$APP_NAME/target
- 配置路径: /var/packages/$APP_NAME/config
- 日志路径: /var/packages/$APP_NAME/logs

## 技术规格
- 应用类型: $APP_TYPE
- 系统要求: Linux x86_64
- 内存要求: 128MB 最小
- 磁盘空间: 100MB 最小
- 端口占用: ${APP_PORT:-8080}

## 开发者信息
- 作者: $AUTHOR
- 网站: ${WEBSITE:-"N/A"}
- 许可证: MIT

## 应用文件结构
\`\`\`
$APP_NAME.nap/
├── INFO                          # 应用元信息
├── application/                   # 应用文件
│   └── (添加你的应用文件)
├── icons/                        # 应用图标
│   ├── icon_72.png
│   └── icon_256.png
├── scripts/                      # 生命周期脚本
│   ├── installer.sh
│   ├── start.sh
│   ├── stop.sh
│   └── status.sh
└── config/                       # 配置文件
    ├── default_config.json
    └── resources.json
\`\`\`

## 下一步操作
1. 将你的应用文件放到 \`application/\` 目录
2. 创建应用图标 (72x72 和 256x256 PNG)
3. 根据需要修改脚本文件
4. 测试应用功能
5. 构建最终的.nap文件

## 构建命令
\`\`\`bash
# 进入应用目录
cd $APP_NAME

# 添加应用文件
# 将你的应用文件复制到 application/ 目录

# 添加应用图标
# 将icon_72.png和icon_256.png放到icons/目录

# 返回上级目录
cd ..

# 构建应用包
tar -czf $APP_NAME-1.0.0.nap $APP_NAME/
\`\`\`

## 支持
如有问题请联系开发者或查看官方文档。
EOF

# 创建占位符文件
echo -e "${GREEN}8. 创建占位符文件...${NC}"
touch "$APP_DIR/application/.gitkeep"
echo "将你的应用文件放到这个目录" > "$APP_DIR/application/README.txt"

# 生成简单的SVG图标（如果系统支持）
if command -v convert &> /dev/null; then
    echo -e "${GREEN}9. 生成临时图标...${NC}"
    # 使用ImageMagick创建临时图标
    convert -size 72x72 xc:\( "#667eea" "#764ba2" \) \
            -pointsize 24 -fill white -gravity center \
            -annotate +0+0 "$(echo "$APP_NAME" | cut -c1-2 | tr '[:lower:]' '[:upper:]')" \
            "$APP_DIR/icons/icon_72.png" 2>/dev/null || true

    convert -size 256x256 xc:\( "#667eea" "#764ba2" \) \
            -pointsize 72 -fill white -gravity center \
            -annotate +0+0 "$(echo "$APP_NAME" | cut -c1-2 | tr '[:lower:]' '[:upper:]')" \
            "$APP_DIR/icons/icon_256.png" 2>/dev/null || true
else
    # 创建简单的PNG占位符
    echo -e "${YELLOW}ImageMagick未安装，跳过图标生成${NC}"
    echo "请手动添加图标文件:"
    echo "  - icons/icon_72.png (72x72 PNG)"
    echo "  - icons/icon_256.png (256x256 PNG)"
fi

# 设置脚本权限
chmod +x "$APP_DIR"/scripts/*.sh

echo ""
echo -e "${GREEN}✓ 应用包模板创建完成！${NC}"
echo ""
echo -e "${BLUE}应用信息:${NC}"
echo "  应用名称: $APP_NAME"
echo "  应用类型: $APP_TYPE"
echo "  应用分类: $CATEGORY"
echo "  应用端口: ${APP_PORT:-"未设置"}"
echo ""
echo -e "${BLUE}目录结构:${NC}"
echo "  $APP_NAME/"
echo "  ├── INFO"
echo "  ├── application/    # 放置应用文件"
echo "  ├── icons/          # 放置应用图标"
echo "  ├── scripts/        # 生命周期脚本"
echo "  ├── config/         # 配置文件"
echo "  └── README.md       # 使用说明"
echo ""
echo -e "${YELLOW}下一步操作:${NC}"
echo "  1. 查看README.md了解详细说明"
echo "  2. 将应用文件放到 application/ 目录"
echo "  3. 添加应用图标到 icons/ 目录"
echo "  4. 根据需要修改脚本和配置"
echo "  5. 构建应用包:"
echo "     ${GREEN}tar -czf $APP_NAME-1.0.0.nap $APP_NAME/${NC}"
echo ""
echo -e "${BLUE}快速构建命令:${NC}"
echo -e "${GREEN}cd $APP_NAME && echo '添加应用文件后按Ctrl+D' && bash && cd .. && tar -czf $APP_NAME-1.0.0.nap $APP_NAME/${NC}"
echo ""
echo "完成！"