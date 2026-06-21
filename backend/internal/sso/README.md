# SSO统一身份认证系统

## 🎯 系统概述

NAS Dashboard SSO系统是一个完整的OAuth2/OIDC Provider实现，支持多种身份提供商的统一认证。用户可以使用Google、GitHub、Microsoft、LDAP/AD、Immich等账号一键登录NAS系统。

## ✨ 核心特性

### 🔐 OAuth2/OIDC Provider
- 完整的OAuth2授权码流程
- OpenID Connect (OIDC)支持
- 标准的端点实现：
  - `/sso/authorize` - 授权端点
  - `/sso/token` - 令牌端点
  - `/sso/userinfo` - 用户信息端点
  - `/sso/.well-known/openid-configuration` - OIDC发现端点
  - `/sso/jwks` - JSON Web Key Set端点

### 🌍 多身份提供商支持
- **Google OAuth2** - Google账号登录
- **GitHub OAuth** - GitHub账号登录
- **Microsoft Azure AD** - 企业微软账号登录
- **LDAP/AD** - 企业目录服务集成
- **Immich** - 照片管理应用集成

### 🛡️ 安全特性
- JWT令牌签名验证
- 授权码一次性使用
- State参数防CSRF攻击
- 令牌自动刷新机制
- 会话管理和撤销

### 🎨 用户界面
- 美观的SSO登录页面
- 响应式设计，支持移动端
- 多种登录方式选择
- 实时登录状态反馈

## 🚀 快速开始

### 1. 配置SSO服务器

复制配置文件并修改：
```bash
cd /data/nas-dashboard/backend
cp .env.sso.example .env.sso
```

编辑 `.env.sso` 文件，配置：
- 基本SSO设置
- OAuth2提供商凭据
- LDAP/AD连接信息
- Immich集成配置

### 2. 启动服务器

```bash
# 后端服务器（包含SSO端点）
cd /data/nas-dashboard/backend
go run cmd/server/main.go

# 前端服务器
cd /data/nas-dashboard/frontend
npm run dev
```

### 3. 访问SSO登录

- 前端SSO登录: http://localhost:5177/sso/login
- 后端SSO端点: http://localhost:8888/sso/*

## 🔧 配置指南

### Google OAuth设置

1. 访问 [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. 创建OAuth 2.0客户端ID
3. 添加授权重定向URI: `http://localhost:8888/sso/callback`
4. 复制客户端ID和密钥到配置文件

### GitHub OAuth设置

1. 访问 [GitHub开发者设置](https://github.com/settings/developers)
2. 注册新的OAuth应用
3. 设置授权回调URL: `http://localhost:8888/sso/callback`
4. 复制客户端ID和密钥

### Microsoft Azure AD设置

1. 访问 [Azure Portal](https://portal.azure.com/#view/Microsoft_AAD_RegisteredApps/ApplicationsListBlade)
2. 注册新应用
3. 添加重定向URI: `http://localhost:8888/sso/callback`
4. 配置权限：`openid`, `profile`, `email`, `User.Read`
5. 复制客户端ID和密钥

### LDAP/AD配置

```bash
# LDAP配置示例
LDAP_SERVER=ldap.example.com
LDAP_PORT=389
LDAP_BASE_DN=dc=example,dc=com
LDAP_BIND_DN=cn=admin,dc=example,dc=com
LDAP_BIND_PASSWORD=your-password
LDAP_USER_FILTER=(uid=%s)
```

### Immich集成配置

```bash
# Immich配置示例
IMMICH_SERVER_URL=http://localhost:2283
IMMICH_API_KEY=your-immich-api-key
```

## 📡 API端点

### OAuth2标准端点

| 端点 | 方法 | 描述 |
|------|------|------|
| `/sso/authorize` | GET | OAuth2授权端点 |
| `/sso/callback` | GET | OAuth2回调端点 |
| `/sso/token` | POST | 令牌交换端点 |
| `/sso/userinfo` | GET | 用户信息端点 |
| `/sso/.well-known/openid-configuration` | GET | OIDC发现端点 |
| `/sso/jwks` | GET | JSON Web Key Set端点 |
| `/sso/revoke` | POST | 令牌撤销端点 |
| `/sso/introspect` | POST | 令牌内省端点 |

### 授权流程示例

1. **发起授权请求**
```http
GET /sso/authorize?
    client_id=nas-dashboard&
    response_type=code&
    redirect_uri=http://localhost:5177/sso/callback&
    scope=openid+profile+email&
    state=random_state&
    provider=google
```

2. **获取授权码回调**
```http
GET http://localhost:5177/sso/callback?
    code=authorization_code&
    state=random_state
```

3. **交换访问令牌**
```http
POST /sso/token
Content-Type: application/x-www-form-urlencoded

code=authorization_code&
redirect_uri=http://localhost:5177/sso/callback&
client_id=nas-dashboard&
grant_type=authorization_code
```

4. **获取用户信息**
```http
GET /sso/userinfo
Authorization: Bearer access_token
```

## 🎯 前端集成

### 使用SSO登录组件

```typescript
import { ssoClient } from '@/api/sso'

// 发起SSO登录
const loginWithProvider = (provider: string) => {
  const state = generateRandomState(32)
  const authUrl = ssoClient.generateAuthUrl(provider, state)

  // 保存state用于验证
  sessionStorage.setItem('sso_state', state)
  sessionStorage.setItem('sso_provider', provider)

  // 重定向到SSO登录页面
  window.location.href = authUrl
}
```

### 处理SSO回调

```typescript
import { performSSOLogin } from '@/api/sso'

// 在回调页面处理登录结果
const handleCallback = async () => {
  const code = new URLSearchParams(window.location.search).get('code')
  const state = new URLSearchParams(window.location.search).get('state')

  try {
    await performSSOLogin(code, state)
    // 登录成功，跳转到主页
    router.push('/')
  } catch (error) {
    console.error('SSO login failed:', error)
    // 处理登录失败
  }
}
```

## 🔐 安全建议

1. **生产环境配置**
   - 使用强随机密钥作为JWT签名密钥
   - 启用HTTPS传输
   - 配置适当的令牌过期时间

2. **会话管理**
   - 实现令牌自动刷新
   - 定期检查会话有效性
   - 提供登出和令牌撤销功能

3. **防护措施**
   - 验证state参数防CSRF
   - 限制授权码使用次数和有效期
   - 监控异常登录行为

## 🐛 故障排除

### 常见问题

1. **授权失败**
   - 检查OAuth应用配置
   - 验证重定向URI设置
   - 确认客户端凭据正确

2. **令牌验证失败**
   - 检查JWT密钥配置
   - 验证令牌过期时间
   - 确认签名算法一致

3. **用户信息获取失败**
   - 验证访问令牌有效性
   - 检查OAuth权限范围
   - 确认API端点可访问

## 📚 技术架构

### 后端架构
```
backend/internal/sso/
├── provider.go       # 身份提供商实现
├── handler.go        # SSO端点处理器
└── README.md         # SSO系统文档
```

### 前端架构
```
frontend/src/
├── apps/
│   ├── SSOLogin.vue    # SSO登录页面
│   └── SSOCallback.vue # SSO回调处理
├── api/
│   └── sso.ts          # SSO API客户端
└── router/
    └── index.ts        # 路由配置
```

### 数据流
```
用户点击登录 
    ↓
选择身份提供商 
    ↓
重定向到提供商授权页面 
    ↓
用户同意授权 
    ↓
提供商回调SSO服务器 
    ↓
SSO服务器交换授权码 
    ↓
创建用户会话 
    ↓
回调前端应用 
    ↓
前端获取用户信息 
    ↓
登录成功
```

## 🎉 完成

SSO系统现已集成到NAS Dashboard中！用户可以：

- ✅ 使用多种身份提供商登录
- ✅ 享受统一的登录体验
- ✅ 安全的单点登录
- ✅ 无缝的应用间认证

系统现在可以作为OAuth2/OIDC Provider，为集成的应用提供统一的身份认证服务！