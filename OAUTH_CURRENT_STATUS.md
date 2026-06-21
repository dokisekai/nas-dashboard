# OAuth/OIDC 服务器修复状态报告

## 已完成的修复 ✅

### 1. 后端端点路径修复
- ✅ 修改后端API返回的端点从 `/oauth2/*` 改为 `/sso/*`
- ✅ 固定Issuer URL为 `https://192.168.50.10:8888`

### 2. 前端显示修复
- ✅ 修改前端显示端点为 `/sso/*`

### 3. 当前API返回
```json
{
  "issuer_url": "https://192.168.50.10:8888",
  "authorize_endpoint": "https://192.168.50.10:8888/sso/authorize",
  "token_endpoint": "https://192.168.50.10:8888/sso/token",
  "userinfo_endpoint": "https://192.168.50.10:8888/sso/userinfo",
  "jwks_endpoint": "https://192.168.50.10:8888/sso/jwks",
  "discovery_endpoint": "https://192.168.50.10:8888/sso/.well-known/openid-configuration",
  "running": true
}
```

## 当前问题 ⚠️

### 400错误原因
前端截图显示的400错误是因为：
1. 前端代理配置使用HTTP target，但后端使用HTTPS
2. 已修改前端vite.config.ts使用HTTPS target
3. **需要重启前端开发服务器才能生效**

### Issuer URL显示前端端口
截图显示 `http://192.168.50.10:5173` 而不是后端地址，这是因为：
1. API调用失败（400错误）
2. 前端fallback到autoDetectServerConfig逻辑
3. 使用了window.location的地址

## 下一步操作

### 重启前端开发服务器
执行以下命令：
```bash
cd /data/nas-dashboard/frontend

# 终止旧的vite进程
pkill -f "vite.*nas-dashboard"

# 重新启动
npm run dev -- --host &
```

或者简单重启浏览器并刷新页面。

## 验证步骤

重启前端后，检查以下URL：
1. 刷新浏览器页面 `http://192.168.50.10:5173`
2. 检查SSO管理界面是否显示正确的端点
3. 确认400错误是否消失

## 正确的端点（后端已配置）

| 功能 | 端点 |
|------|------|
| 授权 | `https://192.168.50.10:8888/sso/authorize` |
| Token | `https://192.168.50.10:8888/sso/token` |
| UserInfo | `https://192.168.50.10:8888/sso/userinfo` |
| JWKS | `https://192.168.50.10:8888/sso/jwks` |
| OIDC发现 | `https://192.168.50.10:8888/sso/.well-known/openid-configuration` |
