# 应用包制作完整工作流程

## 🎯 目标
让任何人都能在10分钟内创建一个可用的NAS Dashboard应用包

## 📋 前置要求

### 系统要求
- Linux或MacOS系统
- bash shell
- tar和gzip工具

### 可选工具
- ImageMagick (用于生成图标)
- convert命令 (图片处理)

### 知识要求
- 基本的Linux命令行操作
- 了解你要打包的应用类型
- 知道如何启动和停止你的应用

## ⚡ 快速开始（10分钟版）

### 步骤1: 使用模板生成器（2分钟）

```bash
# 1. 进入工具目录
cd /path/to/nas-dashboard/tools

# 2. 生成应用模板
./create-app-template.sh my-app

# 或者指定更多选项
./create-app-template.sh -t docker -p 8080 -c utilities my-docker-app
```

### 步骤2: 添加应用文件（3分钟）

```bash
# 1. 进入生成的应用目录
cd my-app

# 2. 将应用文件放到application目录
cp /path/to/your/application/* application/

# 3. 添加应用图标（可选，系统会生成临时图标）
# icons/icon_72.png (72x72)
# icons/icon_256.png (256x256)
```

### 步骤3: 自定义配置（3分钟）

```bash
# 1. 编辑INFO文件
nano INFO  # 修改应用名称、描述、端口等信息

# 2. 编辑启动脚本（如果需要）
nano scripts/start.sh

# 3. 添加默认配置（如果需要）
nano config/default_config.json
```

### 步骤4: 构建和测试（2分钟）

```bash
# 1. 返回上级目录
cd ..

# 2. 构建应用包
tar -czf my-app-1.0.0.nap my-app/

# 3. 查看包内容
tar -tzf my-app-1.0.0.nap

# 4. 测试脚本（可选）
tar -xzf my-app-1.0.0.nap -C /tmp/test
cd /tmp/test/my-app
bash -n scripts/*.sh  # 检查脚本语法
```

### 步骤5: 上传和安装

1. 打开浏览器访问: `http://192.168.50.10:5173/application-center`
2. 点击"上传应用包"
3. 选择`my-app-1.0.0.nap`文件
4. 上传完成后点击"安装"
5. 等待安装完成

## 🔧 详细工作流程

### 阶段1: 规划和准备（5-10分钟）

#### 1.1 确定应用信息

回答以下问题：
- ✅ 应用名称（小写字母、数字、连字符）
- ✅ 应用类型（docker/nodejs/python/golang/script）
- ✅ 应用分类（media/productivity/utilities/security/network）
- ✅ 应用端口（如果需要网络访问）
- ✅ 依赖项（docker、python3、nodejs等）
- ✅ 资源要求（内存、磁盘空间）

#### 1.2 准备应用文件

根据应用类型准备文件：

**Docker应用**:
```bash
# 准备Docker镜像信息
# docker pull your-image:tag
# docker save your-image:tag > application/image.tar
```

**Node.js应用**:
```bash
# 准备node_modules
cd your-app
npm install --production
tar -czf app.tar.gz node_modules/ *.json *.js
```

**Python应用**:
```bash
# 准备Python环境和依赖
pip3 install -r requirements.txt -t application/lib/
# 或打包整个项目
tar -czf app.tar.gz *.py requirements.txt
```

**Golang应用**:
```bash
# 编译二进制文件
go build -o application/my-app main.go
```

**脚本应用**:
```bash
# 准备shell脚本
cp your-script.sh application/run.sh
chmod +x application/run.sh
```

### 阶段2: 创建应用包（15-30分钟）

#### 2.1 使用模板生成器

```bash
cd /path/to/nas-dashboard/tools

# 基本用法
./create-app-template.sh app-name

# 指定应用类型
./create-app-template.sh -t docker -p 8080 my-docker-app

# 指定所有选项
./create-app-template.sh \
  -t nodejs \
  -p 3000 \
  -c productivity \
  -d "我的生产力工具" \
  -a "Your Name" \
  -w "https://yourwebsite.com" \
  -r nodejs,npm \
  my-productivity-app
```

#### 2.2 手动创建（可选）

如果不用模板生成器，手动创建：

```bash
# 1. 创建目录结构
mkdir -p my-app/{application,icons,scripts,config,wizard}

# 2. 创建INFO文件（参考APP_PACKAGE_GUIDE.md）
nano my-app/INFO

# 3. 创建脚本（使用模板或参考指南）
nano my-app/scripts/installer.sh
nano my-app/scripts/start.sh
nano my-app/scripts/stop.sh
nano my-app/scripts/status.sh
chmod +x my-app/scripts/*.sh

# 4. 创建配置文件
nano my-app/config/default_config.json
nano my-app/config/resources.json
```

#### 2.3 添加应用文件

```bash
cd my-app

# 添加应用主程序
cp /path/to/your-app application/

# 添加应用图标
# icon_72.png: 72x72 PNG格式
# icon_256.png: 256x256 PNG格式

# 生成临时图标（如果没有ImageMagick）
echo "创建临时图标..."
# 系统会自动生成简单图标
```

#### 2.4 自定义脚本

**启动脚本自定义**:

```bash
nano scripts/start.sh

# 根据应用类型修改启动逻辑
# Docker应用: docker run ...
# Node.js应用: node application/app.js
# Python应用: python3 application/main.py
# Golang应用: ./application/my-app
# 脚本应用: bash application/run.sh
```

**配置文件自定义**:

```bash
nano config/default_config.json

# 添加应用特定配置
{
  "app_port": 8080,
  "app_host": "0.0.0.0",
  "app_debug": false,
  "app_log_level": "info"
}
```

### 阶段3: 测试和调试（10-15分钟）

#### 3.1 本地测试

```bash
# 1. 测试脚本语法
cd my-app
bash -n scripts/*.sh

# 2. 手动运行安装脚本（测试）
sudo bash scripts/installer.sh

# 3. 手动运行启动脚本（测试）
sudo bash scripts/start.sh

# 4. 检查状态
bash scripts/status.sh

# 5. 停止应用
sudo bash scripts/stop.sh

# 6. 清理测试文件
sudo rm -rf /var/packages/my-app
```

#### 3.2 打包测试

```bash
# 1. 返回上级目录
cd ..

# 2. 构建应用包
tar -czf my-app-1.0.0.nap my-app/

# 3. 验证包内容
tar -tzf my-app-1.0.0.nap | grep -E 'INFO|scripts/(start|stop|status)\.sh$'

# 4. 检查文件大小
du -h my-app-1.0.0.nap

# 5. 验证包格式
file my-app-1.0.0.nap  # 应该显示gzip compressed
```

#### 3.3 解压测试

```bash
# 1. 解压到测试目录
mkdir -p /tmp/test
tar -xzf my-app-1.0.0.nap -C /tmp/test

# 2. 检查目录结构
ls -la /tmp/test/my-app/

# 3. 验证文件权限
ls -la /tmp/test/my-app/scripts/

# 4. 查看INFO文件
cat /tmp/test/my-app/INFO

# 5. 清理测试目录
rm -rf /tmp/test
```

### 阶段4: 部署和验证（5-10分钟）

#### 4.1 上传应用包

1. 打开应用中心: `http://192.168.50.10:5173/application-center`
2. 点击"上传应用包"按钮
3. 选择`my-app-1.0.0.nap`文件
4. 等待上传完成

#### 4.2 安装应用

1. 在"可用应用"列表中找到你的应用
2. 点击"安装"按钮
3. 观察安装进度
4. 等待安装完成

#### 4.3 测试应用

1. 点击"启动"按钮
2. 检查应用状态是否变为"运行中"
3. 访问应用端口（如`http://192.168.50.10:8080`）
4. 测试应用功能
5. 点击"停止"按钮
6. 重新启动应用
7. 最后卸载应用

## 🚀 高级工作流程

### 完整开发流程

#### 1. 开发环境准备

```bash
# 创建开发目录
mkdir -p ~/nas-app-dev/{apps,packages,templates}
cd ~/nas-app-dev

# 克隆或复制工具
cp -r /path/to/nas-dashboard/tools/* templates/
```

#### 2. 应用开发

```bash
# 使用模板创建应用
./templates/create-app-template.sh -t nodejs my-app

# 进入应用目录
cd my-app

# 开发应用（这里假设是Node.js应用）
npm init -y
npm install express --save

# 创建应用主程序
cat > application/app.js << 'EOF'
const express = require('express');
const app = express();
const PORT = process.env.PORT || 3000;

app.get('/', (req, res) => {
  res.send('Hello from My App!');
});

app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
EOF
```

#### 3. 迭代测试

```bash
# 创建测试脚本
cat > test-app.sh << 'EOF'
#!/bin/bash
set -e

echo "=== 测试应用安装 ==="
# 1. 构建包
tar -czf ../my-app-test.nap .

# 2. 上传到测试服务器
scp ../my-app-test.nap user@testserver:/tmp/

# 3. 在测试服务器上安装
ssh user@testserver "nas-cli install /tmp/my-app-test.nap"

# 4. 测试启动
ssh user@testserver "nas-cli start my-app"

# 5. 检查状态
ssh user@testserver "nas-cli status my-app"

# 6. 测试功能
curl http://testserver:3000/

# 7. 清理
ssh user@testserver "nas-cli uninstall my-app"
EOF

chmod +x test-app.sh
./test-app.sh
```

#### 4. 版本管理

```bash
# 使用Git进行版本控制
cd ~/nas-app-dev/my-app
git init
git add .
git commit -m "Initial version 1.0.0"

# 创建版本标签
git tag v1.0.0

# 后续更新
git add .
git commit -m "Add new feature"
git tag v1.1.0
```

#### 5. 自动化构建

```bash
# 创建自动化构建脚本
cat > build.sh << 'EOF'
#!/bin/bash
set -e

VERSION=$(git describe --tags --always)
APP_NAME=$(basename $(pwd))
BUILD_DIR="../build"

echo "Building $APP_NAME version $VERSION..."

# 创建构建目录
mkdir -p "$BUILD_DIR"

# 更新INFO文件中的版本
sed -i "s/version=.*/version=$VERSION/" INFO

# 构建应用包
tar -czf "$BUILD_DIR/$APP_NAME-$VERSION.nap" .

echo "Build complete: $BUILD_DIR/$APP_NAME-$VERSION.nap"
EOF

chmod +x build.sh
./build.sh
```

### 发布流程

#### 1. 准备发布

```bash
# 1. 更新版本号
nano INFO  # 修改version字段

# 2. 更新CHANGELOG
cat > CHANGELOG.md << 'EOF'
# Version 1.0.0 (2024-01-15)

## 新功能
- 初始发布
- 支持基本功能A
- 支持基本功能B

## 修复
- 修复问题X

## 已知问题
- 暂无
EOF

# 3. 生成发布包
./build.sh
```

#### 2. 测试发布

```bash
# 在测试环境完整测试
./test-app.sh

# 创建测试报告
cat > TEST_REPORT.md << 'EOF'
# 测试报告

## 测试环境
- 系统: Ubuntu 22.04
- NAS Dashboard: v1.0.0
- 测试日期: 2024-01-15

## 测试结果
- ✅ 安装测试通过
- ✅ 启动测试通过
- ✅ 功能测试通过
- ✅ 停止测试通过
- ✅ 卸载测试通过

## 问题
- 无

## 建议
- 可以发布
EOF
```

#### 3. 正式发布

```bash
# 1. 创建GitHub Release
gh release create v1.0.0 \
  build/my-app-1.0.0.nap \
  --notes "Release version 1.0.0" \
  --title "v1.0.0 - Initial Release"

# 2. 提交到应用仓库（如果适用）
# 按照仓库的提交流程操作

# 3. 更新文档
nano README.md  # 更新使用说明
nano CHANGELOG.md  # 添加发布记录
```

## 🔍 故障排除指南

### 常见问题解决

#### 问题1: 模板生成器失败

```bash
# 检查脚本权限
ls -la create-app-template.sh

# 重新设置权限
chmod +x create-app-template.sh

# 检查bash版本
bash --version  # 需要4.0+

# 手动执行（调试模式）
bash -x create-app-template.sh my-app
```

#### 问题2: 打包失败

```bash
# 检查磁盘空间
df -h

# 检查文件权限
ls -la my-app/

# 重新打包
rm -f my-app-1.0.0.nap
tar -czf my-app-1.0.0.nap my-app/

# 验证包内容
tar -tzf my-app-1.0.0.nap | head -20
```

#### 问题3: 安装失败

```bash
# 查看安装日志
tail -f /var/log/nas-dashboard/application.log

# 检查脚本语法
bash -n my-app/scripts/*.sh

# 手动运行安装脚本
cd my-app
sudo bash scripts/installer.sh

# 检查系统资源
df -h  # 磁盘空间
free -h  # 内存
uname -m  # 架构
```

#### 问题4: 启动失败

```bash
# 查看应用日志
tail -f /var/packages/my-app/logs/*.log

# 检查启动脚本
bash -x scripts/start.sh

# 检查端口占用
netstat -tlnp | grep 8080

# 检查进程
ps aux | grep my-app

# 检查systemd服务
systemctl status my-app
```

#### 问题5: 状态检查错误

```bash
# 手动运行状态脚本
bash scripts/status.sh
echo $?  # 应该是0(running)或1(stopped)

# 检查实际进程状态
ps aux | grep my-app
docker ps | grep my-app
systemctl status my-app

# 修复状态脚本
nano scripts/status.sh
```

## 📚 最佳实践

### 1. 文件命名规范

```bash
# 应用名称: 只包含小写字母、数字、连字符
my-app
my_app  # ❌ 错误
MyApp   # ❌ 错误
my-app-123  # ✅ 正确

# 文件名: 使用小写字母和下划线
my_app.sh
my-app.sh  # 也可以，但建议统一
```

### 2. 目录组织规范

```bash
# 推荐的目录结构
my-app.nap/
├── INFO                    # 必需，第一个文件
├── application/            # 应用主程序
│   ├── my-app             # 主程序
│   ├── lib/                # 库文件（可选）
│   └── resources/          # 资源文件（可选）
├── icons/                  # 图标
│   ├── icon_72.png        # 小图标
│   └── icon_256.png       # 大图标
├── scripts/                # 脚本
│   ├── installer.sh       # 可选
│   ├── start.sh           # 必需
│   ├── stop.sh            # 必需
│   └── status.sh          # 必需
└── config/                 # 配置（可选）
    ├── default_config.json
    └── resources.json
```

### 3. 脚本编写规范

```bash
# 脚本头部
#!/bin/bash
set -e  # 遇到错误立即退出
set -u  # 使用未定义变量时报错
set -o pipefail  # 管道命令失败时退出

# 使用注释解释重要步骤
# 创建数据目录
mkdir -p "$DATA_PATH"

# 使用变量而不是硬编码
INSTALL_PATH="/var/packages/$APP_NAME"

# 错误处理
if ! command -v docker &> /dev/null; then
    echo "错误: Docker未安装"
    exit 1
fi

# 返回明确的退出码
exit 0  # 成功
exit 1  # 失败
```

### 4. 版本管理规范

```bash
# 使用语义化版本
MAJOR.MINOR.PATCH
1.0.0  # 主版本.次版本.补丁版本

# 变更规则
1.0.0 -> 1.0.1  # 修复bug
1.0.1 -> 1.1.0  # 新增功能（向后兼容）
1.1.0 -> 2.0.0  # 破坏性变更
```

### 5. 测试规范

```bash
# 创建完整的测试脚本
cat > test.sh << 'EOF'
#!/bin/bash
set -e

echo "=== 功能测试 ==="

# 1. 安装测试
echo "测试安装..."
# 安装逻辑

# 2. 启动测试
echo "测试启动..."
# 启动逻辑

# 3. 功能测试
echo "测试功能..."
# 功能逻辑

# 4. 停止测试
echo "测试停止..."
# 停止逻辑

# 5. 卸载测试
echo "测试卸载..."
# 卸载逻辑

echo "=== 所有测试通过 ==="
EOF

chmod +x test.sh
./test.sh
```

## 🎓 学习资源

### 推荐阅读

1. **官方文档**
   - APP_PACKAGE_GUIDE.md - 详细应用包制作指南
   - APP_PACKAGE_DESIGN.md - 系统设计文档
   - API_DOCUMENTATION.md - API文档

2. **示例应用**
   - examples/sample-app/ - 简单示例
   - 内置应用包 - 参考现有应用

3. **社区资源**
   - GitHub Issues - 问题讨论
   - 开发者论坛 - 经验分享

### 技能提升

1. **脚本编写**
   - Bash脚本基础教程
   - Shell脚本最佳实践
   - Linux系统编程

2. **应用打包**
   - Docker容器化
   - 服务部署
   - 系统集成

3. **故障排除**
   - Linux系统管理
   - 网络调试
   - 日志分析

## 💡 小贴士

### 效率提升

1. **使用别名**
```bash
alias app-create='~/nas-app-dev/templates/create-app-template.sh'
alias app-build='cd ~/nas-app-dev/my-app && cd .. && tar -czf my-app-1.0.0.nap my-app/'
alias app-test='./test-app.sh'
```

2. **创建模板**
```bash
# 保存常用配置
mkdir -p ~/templates/app-types
cp -r my-app ~/templates/app-types/docker-app
```

3. **批量操作**
```bash
# 批量构建多个版本
for version in 1.0.0 1.1.0 2.0.0; do
    sed -i "s/version=.*/version=$version/" INFO
    tar -czf my-app-$version.nap my-app/
done
```

### 质量保证

1. **代码审查**
   - 检查脚本语法
   - 验证配置文件
   - 测试所有功能

2. **自动化测试**
   - 创建测试脚本
   - 集成CI/CD
   - 定期回归测试

3. **文档完善**
   - 编写README
   - 更新CHANGELOG
   - 提供示例配置

## 🎯 总结

制作NAS Dashboard应用包的完整流程：

1. **规划** (5分钟) - 确定应用信息和要求
2. **创建** (10分钟) - 使用模板生成器创建基础结构
3. **开发** (15分钟) - 添加应用文件和自定义脚本
4. **测试** (10分钟) - 本地测试和验证
5. **构建** (2分钟) - 打包成.nap文件
6. **部署** (5分钟) - 上传和安装测试
7. **发布** (3分钟) - 正式发布应用包

总时间: **约50分钟**即可完成第一个应用包的制作！

随着经验积累，这个过程可以缩短到20-30分钟。

记住：**从简单开始，逐步完善**。先制作一个能运行的基础版本，然后再逐步添加高级功能。

祝你制作出优秀的应用包！🚀