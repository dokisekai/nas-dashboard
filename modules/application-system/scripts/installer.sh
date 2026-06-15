#!/bin/bash

# 安装脚本示例
set -e

echo "开始安装 Hello World 应用..."

# 创建应用目录
INSTALL_DIR="/var/packages/hello-world"
mkdir -p "$INSTALL_DIR/target"
mkdir -p "$INSTALL_DIR/config"
mkdir -p "$INSTALL_DIR/logs"

# 复制应用文件
cp -r /tmp/hello-world/* "$INSTALL_DIR/target/"

# 设置权限
chmod +x "$INSTALL_DIR/target/hello-world"
chown -R nobody:nogroup "$INSTALL_DIR"

# 创建systemd服务文件（如果需要）
cat > /etc/systemd/system/hello-world.service << EOF
[Unit]
Description=Hello World Application
After=network.target

[Service]
Type=simple
User=nobody
WorkingDirectory=$INSTALL_DIR/target
ExecStart=$INSTALL_DIR/target/hello-world
Restart=on-failure

[Install]
WantedBy=multi-user.target
EOF

# 重载systemd
systemctl daemon-reload

echo "Hello World 应用安装完成！"
exit 0