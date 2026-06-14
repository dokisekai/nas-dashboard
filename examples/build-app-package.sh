#!/bin/bash

# 应用包打包脚本

set -e

APP_NAME="hello-world"
APP_VERSION="1.0.0"
PACKAGE_NAME="${APP_NAME}-${APP_VERSION}.nap"
BUILD_DIR="build"
SAMPLE_DIR="sample-app"

echo "开始构建应用包 ${PACKAGE_NAME}..."

# 清理并创建构建目录
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 复制示例应用文件到构建目录
cp -r "$SAMPLE_DIR" "$BUILD_DIR/$APP_NAME"

# 创建应用二进制文件（这里使用一个简单的shell脚本作为示例）
cat > "$BUILD_DIR/$APP_NAME/application/hello-world" << 'EOF'
#!/bin/bash

# Hello World 应用示例
echo "Hello World Application v1.0.0"
echo "正在启动服务..."
echo "监听地址: http://0.0.0.0:8080"

# 这里应该有实际的应用逻辑
# 为了示例，我们创建一个简单的HTTP服务器
while true; do
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] 服务运行中..."
    sleep 60
done
EOF

chmod +x "$BUILD_DIR/$APP_NAME/application/hello-world"

# 创建应用图标目录
mkdir -p "$BUILD_DIR/$APP_NAME/icons"

# 创建简单的PNG图标（使用ImageMagick或提供占位符）
if command -v convert &> /dev/null; then
    convert -size 72x72 xc:blue -pointsize 24 -fill white -gravity center \
            -annotate +0+0 "HW" "$BUILD_DIR/$APP_NAME/icons/icon_72.png"
    convert -size 256x256 xc:blue -pointsize 72 -fill white -gravity center \
            -annotate +0+0 "HW" "$BUILD_DIR/$APP_NAME/icons/icon_256.png"
else
    echo "ImageMagick未安装，跳过图标生成"
    # 创建占位符文件
    touch "$BUILD_DIR/$APP_NAME/icons/icon_72.png"
    touch "$BUILD_DIR/$APP_NAME/icons/icon_256.png"
fi

# 创建安装向导配置
mkdir -p "$BUILD_DIR/$APP_NAME/wizard"
cat > "$BUILD_DIR/$APP_NAME/wizard/wizard.json" << 'EOF'
{
  "enabled": false,
  "steps": [],
  "configUI": "minimal"
}
EOF

# 打包成tar.gz
cd "$BUILD_DIR"
tar -czf "../${PACKAGE_NAME}" "$APP_NAME"
cd ..

echo "应用包构建完成: ${PACKAGE_NAME}"
echo "应用包大小: $(du -h ${PACKAGE_NAME} | cut -f1)"

# 清理构建目录
rm -rf "$BUILD_DIR"

echo "构建完成！"