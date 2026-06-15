# 应用系统模块

## 概述

应用系统模块提供完整的软件包管理和应用中心功能，允许用户安装、管理和运行各种应用程序。

## 功能特性

### 应用包管理
- 应用打包和分发
- 应用安装和卸载
- 应用版本管理
- 依赖关系处理
- 应用权限控制

### 应用中心
- 应用商店界面
- 应用分类和搜索
- 应用评级和评论
- 应用更新管理
- 应用使用统计

### 应用运行环境
- 容器化运行环境
- 资源限制管理
- 网络隔离配置
- 数据持久化
- 日志收集管理

## 脚本说明

### 1. create-app-template.sh
应用模板创建工具，用于快速创建新的应用包。

```bash
# 使用示例
./create-app-template.sh --name "my-app" --version "1.0.0"
```

### 2. build-app-package.sh
应用打包脚本，将应用代码打包成标准格式的应用包。

```bash
# 使用示例
./build-app-package.sh --source ./my-app --output ./my-app.spk
```

### 3. installer.sh
应用安装脚本，处理应用的安装逻辑。

```bash
# 使用示例
./installer.sh install my-app.spk
```

### 4. start.sh
应用启动脚本，管理应用的启动过程。

```bash
# 使用示例
./start.sh --app my-app --port 8080
```

### 5. stop.sh
应用停止脚本，优雅地停止运行中的应用。

```bash
# 使用示例
./stop.sh --app my-app
```

### 6. status.sh
应用状态查询脚本。

```bash
# 使用示例
./status.sh --app my-app
```

## 配置文件

### 应用配置
- `config/default_config.json` - 默认应用配置
- `config/env_vars.json` - 环境变量配置
- `config/resources.json` - 资源限制配置

### 应用包结构
```
my-app.spk
├── INFO.txt           # 应用信息
├── installer.sh       # 安装脚本
├── app/               # 应用文件
│   ├── bin/          # 可执行文件
│   ├── web/          # Web界面
│   └── config/       # 配置文件
└── scripts/          # 管理脚本
    ├── start.sh
    ├── stop.sh
    └── status.sh
```

## 后端API实现

### 主要API文件
- `backend/internal/api/application.go` - 应用管理API
- `backend/pkg/application/manager.go` - 应用管理器
- `backend/pkg/application/installer.go` - 应用安装器

### 前端组件
- `frontend/src/apps/ApplicationCenter.vue` - 应用中心界面
- `frontend/src/apps/AppCenter.vue` - 应用管理界面
- `frontend/src/components/Application/` - 应用相关组件

## API端点

### 应用管理
- `GET /api/apps` - 获取应用列表
- `GET /api/apps/:id` - 获取应用详情
- `POST /api/apps/install` - 安装应用
- `DELETE /api/apps/:id` - 卸载应用
- `POST /api/apps/:id/start` - 启动应用
- `POST /api/apps/:id/stop` - 停止应用
- `GET /api/apps/:id/status` - 获取应用状态

### 应用包管理
- `POST /api/apps/upload` - 上传应用包
- `GET /api/apps/available` - 获取可用应用
- `POST /api/apps/:id/update` - 更新应用
- `GET /api/apps/categories` - 获取应用分类

## 使用示例

### 创建新应用
```bash
# 使用模板创建应用
cd modules/application-system/scripts/
./create-app-template.sh --name "web-server" --type "service"
```

### 打包应用
```bash
# 构建应用包
./build-app-package.sh --source ./web-server --version 1.0.0
```

### 安装应用
```bash
# 通过API安装应用
curl -X POST http://localhost:8888/api/apps/install \
  -H "Authorization: Bearer <token>" \
  -d '{
    "package": "web-server.spk",
    "config": {
      "port": 8080,
      "resources": {
        "memory": "512M",
        "cpu": "1"
      }
    }
  }'
```

### 启动应用
```bash
# 通过API启动应用
curl -X POST http://localhost:8888/api/apps/web-server/start \
  -H "Authorization: Bearer <token>"
```

## 应用开发

### 应用类型
1. **系统服务** - 长期运行的后台服务
2. **Web应用** - 提供Web界面的应用
3. **命令行工具** - 实用工具集合
4. **插件** - 扩展系统功能

### 开发流程
1. 使用模板创建应用框架
2. 实现应用功能代码
3. 编写安装和配置脚本
4. 打包成标准格式
5. 测试安装和运行

### 应用规范
- 遵循POSIX标准
- 支持无交互安装
- 提供状态监控接口
- 实现优雅关闭
- 记录详细日志

## 部署配置

### 容器化部署
```yaml
# docker-compose.yml
services:
  my-app:
    image: my-app:1.0.0
    ports:
      - "8080:8080"
    volumes:
      - /mnt/data/app-data:/data
    environment:
      - APP_ENV=production
```

### 系统服务部署
```bash
# 创建systemd服务
cat > /etc/systemd/system/my-app.service <<EOF
[Unit]
Description=My Application
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/my-app --port 8080
Restart=always

[Install]
WantedBy=multi-user.target
EOF
```

## 监控和日志

### 应用监控
- CPU和内存使用
- 网络连接状态
- 响应时间监控
- 错误率统计

### 日志管理
- 应用日志收集
- 日志轮转配置
- 日志级别控制
- 日志分析工具

## 安全考虑

### 权限管理
- 最小权限原则
- 文件权限隔离
- 网络访问限制
- 资源使用限制

### 安全加固
- 代码签名验证
- 沙箱隔离运行
- 安全审计日志
- 漏洞扫描检测

## 性能优化

### 资源优化
- 内存使用优化
- CPU时间片分配
- 磁盘IO优化
- 网络带宽管理

### 缓存策略
- 静态资源缓存
- 数据库查询缓存
- 应用状态缓存
- 配置文件缓存

## 故障排除

### 应用启动失败
1. 检查日志文件：`cat /var/log/apps/my-app.log`
2. 验证依赖项：`./status.sh --check-deps`
3. 检查端口占用：`netstat -tulpn | grep 8080`
4. 验证配置文件：`./installer.sh --validate-config`

### 应用性能问题
1. 监控资源使用：`./status.sh --resources`
2. 分析应用日志：`tail -f /var/log/apps/my-app.log`
3. 检查网络连接：`./status.sh --network`
4. 优化资源配置：编辑 `config/resources.json`

### 应用更新失败
1. 备份当前版本：`./installer.sh --backup`
2. 检查更新日志：`./installer.sh --changelog`
3. 验证新版本：`./build-app-package.sh --validate`
4. 回滚到旧版本：`./installer.sh --rollback`

## 相关文档

- [应用系统检查清单](../../docs/APPLICATION_SYSTEM_CHECKLIST.md)
- [应用系统总结](../../docs/APPLICATION_SYSTEM_SUMMARY.md)
- [应用包设计](../../docs/APP_PACKAGE_DESIGN.md)
- [应用包指南](../../docs/APP_PACKAGE_GUIDE.md)
- [插件开发指南](../../docs/PLUGIN_DEVELOPMENT.md)
- [插件系统文档](../../frontend/src/plugin-system/README.md)

## 示例应用

项目包含完整的示例应用：
- **位置**: `examples/sample-app/`
- **功能**: 演示应用打包和安装流程
- **用途**: 作为应用开发参考模板

## 扩展开发

### 添加新应用类型
1. 在 `templates/` 目录创建新模板
2. 更新应用类型定义
3. 添加对应的处理逻辑
4. 更新文档和示例

### 集成第三方应用
1. 创建应用包装脚本
2. 配置依赖项和环境
3. 测试完整安装流程
4. 编写用户文档

## 最佳实践

1. **应用设计**
   - 模块化架构
   - 配置文件分离
   - 日志记录完整
   - 错误处理健壮

2. **打包规范**
   - 标准化目录结构
   - 完整的元数据
   - 清晰的依赖说明
   - 详细的安装指南

3. **部署流程**
   - 预检查环境
   - 备份现有数据
   - 验证安装结果
   - 提供回滚方案

4. **运维管理**
   - 监控应用状态
   - 定期更新维护
   - 及时处理告警
   - 优化资源配置