# OAuth/OIDC 服务器端点修复总结

## 问题描述
前端SSO管理界面显示的OAuth端点路径与后端实际路由不匹配，导致404错误和显示错误信息。

## 根本原因
- **后端实际路由**: `/sso/*` (authorize, token, userinfo等)
- **前端/API显示**: `/oauth2/*` (错误路径)
- **后端API返回**: `/oauth2/*` (硬编码错误)

## 已完成的修复

### 1. 后端修复 (`backend/internal/api/oauth.go`)
```go
// 修复前
AuthorizeEndpoint: issuerURL + "/oauth2/authorize",
TokenEndpoint:     issuerURL + "/oauth2/token",
UserInfoEndpoint:  issuerURL + "/oauth2/userinfo",
JWKSEndpoint:      issuerURL + "/.well-known/jwks.json",
DiscoveryEndpoint: issuerURL + "/.well-known/openid-configuration",

// 修复后
AuthorizeEndpoint: issuerURL + "/sso/authorize",
TokenEndpoint:     issuerURL + "/sso/token",
UserInfoEndpoint:  issuerURL + "/sso/userinfo",
JWKSEndpoint:      issuerURL + "/sso/jwks",
DiscoveryEndpoint: issuerURL + "/sso/.well-known/openid-configuration",
```

### 2. 前端修复 (`frontend/src/apps/SSOManager.vue`)
所有端点显示已更新为 `/sso/*` 路径：
- 授权端点: `/sso/authorize`
- Token端点: `/sso/token`
- UserInfo端点: `/sso/userinfo`
- JWKS端点: `/sso/jwks`
- OIDC发现端点: `/sso/.well-known/openid-configuration`

## 正确的端点路径

| 功能 | 端点路径 |
|------|---------|
| 授权端点 | `http://192.168.50.10:8888/sso/authorize` |
| Token端点 | `http://192.168.50.10:8888/sso/token` |
| UserInfo端点 | `http://192.168.50.10:8888/sso/userinfo` |
| JWKS端点 | `http://192.168.50.10:8888/sso/jwks` |
| OIDC发现 | `http://192.168.50.10:8888/sso/.well-known/openid-configuration` |
| 管理API | `https://localhost:8888/api/oauth/*` |

## 需要用户操作

### 重启后端服务
执行以下命令重启后端以应用代码修改：
```bash
/tmp/restart_backend.sh
```
或手动重启：
```bash
cd /data/nas-dashboard/backend
pkill -f "backend/server"
./server > /tmp/nas-dashboard-backend.log 2>&1 &
```

### 验证修复
重启后端后，访问以下URL验证：
```bash
# 测试API返回
curl -k https://localhost:8888/api/oauth/server/info | jq '.data'

# 测试OIDC发现
curl http://192.168.50.10:8888/sso/.well-known/openid-configuration
```

## 前端状态
前端修改已通过Vite热重载生效，刷新浏览器即可看到更新后的端点路径。

## 注意事项
1. 后端编译已完成，只需重启服务
2. 前端开发服务器正在运行，修改已生效
3. OAuth客户端管理API (`/api/oauth/*`) 不受此修复影响
