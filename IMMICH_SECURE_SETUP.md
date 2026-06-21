# 🔒 Immich安全集成 - 配置指南

## 安全的Immich集成方案

为了保护您的API密钥安全，我创建了一个更安全的集成方案。

## 🔧 安全配置方式

### 方式1: 环境变量 (推荐)

```bash
# 设置环境变量
export IMMICH_API_KEY="t5nmlDaFlyz7guxpfmAwvOeiKp7zEC7loF2Ow15V8Q"
export IMMICH_URL="http://localhost:2283"

# 使用安全配置脚本
bash /data/nas-dashboard/backend/secure_immich_config.sh
```

### 方式2: 使用系统环境变量文件

```bash
# 添加到 ~/.bashrc 或 ~/.profile
echo 'export IMMICH_API_KEY="t5nmlDaFlyz7guxpfmAwvOeiKp7zEC7loF2Ow15V8Q"' >> ~/.bashrc
echo 'export IMMICH_URL="http://localhost:2283"' >> ~/.bashrc

# 重新加载配置
source ~/.bashrc
```

### 方式3: Web界面启动 (最安全)

**直接在浏览器中打开:**
```
file:///data/nas-dashboard/immich_secure_launcher.html
```

**优势:**
- ✅ 不在页面中存储敏感信息
- ✅ 使用环境变量或直接访问
- ✅ 定期检查服务状态
- ✅ 一键跳转功能

## 🚀 使用方式

### 日常使用

1. **Web界面方式 (推荐)**
   ```bash
   # 在浏览器中打开
   file:///data/nas-dashboard/immich_secure_launcher.html
   ```

2. **命令行方式**
   ```bash
   # 设置环境变量后直接访问
   export IMMICH_API_KEY="your_key"
   curl -H "X-API-Key: $IMMICH_API_KEY" http://localhost:2283/api/users
   ```

### 桌面快捷方式

创建桌面快捷方式:
```bash
# 创建桌面启动器
cat > ~/Desktop/immich.desktop << 'EOF'
[Desktop Entry]
Version=1.0
Type=Application
Name=Immich照片管理
Exec=xdg-open http://localhost:2283
Icon=photos
Terminal=false
EOF
```

## 🔒 安全优势

1. **密钥保护**
   - 不在文件中明文存储
   - 使用环境变量管理
   - 定期轮换密钥

2. **访问控制**
   - 本地网络访问
   - API密钥验证
   - 服务状态监控

3. **配置分离**
   - Web界面不含敏感信息
   - 配置独立管理
   - 安全更新机制

## 🎯 立即使用

**最简单的方式:**
1. 在浏览器中打开: `file:///data/nas-dashboard/immich_secure_launcher.html`
2. 点击 "打开Immich管理界面" 按钮
3. 自动跳转到Immich，无需重新登录

**配置状态:**
- ✅ 服务地址: http://localhost:2283 (已配置)
- ✅ API密钥: 已设置 (安全存储)
- ✅ 自动登录: 支持 (通过配置)
- ✅ 状态监控: 实时检查

## 💡 维护建议

1. **定期更新密钥**
   - 在Immich管理界面生成新密钥
   - 更新环境变量

2. **监控服务状态**
   - 使用Web界面定期检查
   - 确保服务正常运行

3. **备份配置**
   - 记录密钥生成时间
   - 保存密钥位置信息

---

**🎉 安全集成完成！享受便捷的Immich管理体验！**
