# 🚀 应用包系统快速启动指南

## 🎯 5分钟快速开始

### 第1步：构建示例应用（2分钟）

```bash
# 进入示例目录
cd /home/hserver/nas-dashboard/examples

# 执行构建脚本
chmod +x build-app-package.sh
./build-app-package.sh

# 验证构建结果
ls -lh hello-world-1.0.0.nap
```

### 第2步：访问应用中心（1分钟）

打开浏览器访问：
```
http://192.168.50.10:5173/application-center
```

### 第3步：上传和安装（2分钟）

1. 点击"上传应用包"按钮
2. 选择`hello-world-1.0.0.nap`文件
3. 等待上传完成
4. 点击"安装"按钮
5. 观察安装进度
6. 安装完成后点击"启动"测试

## 🎨 制作你的第一个应用包

### 方法A：使用模板生成器（推荐）

```bash
# 1. 进入工具目录
cd /home/hserver/nas-dashboard/tools

# 2. 生成应用模板
./create-app-template.sh my-first-app

# 3. 进入应用目录
cd my-first-app

# 4. 添加应用文件
# 将你的应用放到 application/ 目录

# 5. 构建应用包
cd ..
tar -czf my-first-app-1.0.0.nap my-first-app/
```

### 方法B：从示例修改

```bash
# 1. 复制示例
cp -r examples/sample-app my-app

# 2. 修改INFO文件
nano my-app/INFO

# 3. 替换应用文件
cp /path/to/your/app my-app/application/

# 4. 打包
tar -czf my-app-1.0.0.nap my-app/
```

## 📋 应用包最小结构

```
my-app.nap/
├── INFO                    # 应用元信息
├── application/            # 应用程序
│   └── your-app          # 主程序
└── scripts/                # 脚本
    ├── start.sh          # 启动脚本
    ├── stop.sh           # 停止脚本
    └── status.sh         # 状态脚本
```

## 🎯 常见应用模板

### Docker应用
```bash
./create-app-template.sh -t docker -p 8080 my-docker-app
```

### Node.js应用
```bash
./create-app-template.sh -t nodejs -p 3000 my-nodejs-app
```

### Python应用
```bash
./create-app-template.sh -t python -p 5000 my-python-app
```

## 💡 10分钟实战案例

创建Web服务器应用：

```bash
# 1. 生成模板
cd /home/hserver/nas-dashboard/tools
./create-app-template.sh -t docker -p 80 webserver

# 2. 进入目录
cd webserver

# 3. 修改启动脚本
cat > scripts/start.sh << 'EOF'
#!/bin/bash
docker run -d --name webserver -p 80:80 nginx:alpine
EOF

# 4. 修改停止脚本
cat > scripts/stop.sh << 'EOF'
#!/bin/bash
docker stop webserver
docker rm webserver
EOF

# 5. 修改状态脚本
cat > scripts/status.sh << 'EOF'
#!/bin/bash
if docker ps -q -f name=webserver | grep -q .; then
    echo "running"
else
    echo "stopped"
fi
EOF

# 6. 设置权限
chmod +x scripts/*.sh

# 7. 打包
cd ..
tar -czf webserver-1.0.0.nap webserver/

# 8. 上传到应用中心安装测试
```

## 🔧 常用命令

### 构建命令
```bash
# 快速构建
tar -czf app-1.0.0.nap app/

# 查看包内容
tar -tzf app-1.0.0.nap

# 解压测试
tar -xzf app-1.0.0.nap -C /tmp/test
```

### 测试命令
```bash
# 检查脚本语法
bash -n scripts/*.sh

# 测试状态脚本
bash scripts/status.sh
echo $?  # 0=running, 1=stopped
```

## 📚 详细文档

- **QUICKSTART.md** - 本文档
- **APP_PACKAGE_GUIDE.md** - 完整制作指南
- **APP_PACKAGE_WORKFLOW.md** - 详细工作流程
- **examples/README.md** - 示例说明

## ⚡ 下一步

1. 构建示例应用测试系统
2. 阅读详细制作指南
3. 创建你的第一个应用包
4. 在应用中心测试安装

## 🎉 开始创作

现在你可以开始创建自己的应用包了！

记住：**从简单开始，逐步完善**。

---

需要更多帮助？查看[APP_PACKAGE_GUIDE.md](APP_PACKAGE_GUIDE.md)