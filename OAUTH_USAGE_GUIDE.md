# OAuth/OIDC 服务器使用指南

## 前提条件
- ✅ 后端运行在 `https://192.168.50.10:8888`
- ✅ 前端运行在 `http://192.168.50.10:5173`
- ✅ OAuth服务器已配置并运行

## 1. 管理用户

### 查看现有用户
当前系统已有admin用户，密码是您设置的登录密码。

### 添加新用户
1. 访问 NAS Dashboard
2. 进入 **用户管理** 应用
3. 点击"添加用户"
4. 填写用户信息：
   - 用户名
   - 密码
   - 邮箱（可选）
   - 角色（admin/user）

## 2. OAuth客户端信息

已创建的测试客户端：
- **Client ID**: `client-d56UMow9IjKknOnj`
- **Client Secret**: `im0a7SggwKjCplNdGY6zjuLhGOKeXWBC`
- **重定向URI**: `http://localhost:3000/callback`

## 3. OAuth端点

| 功能 | 端点 |
|------|------|
| 授权端点 | `https://192.168.50.10:8888/sso/authorize` |
| Token端点 | `https://192.168.50.10:8888/sso/token` |
| UserInfo端点 | `https://192.168.50.10:8888/sso/userinfo` |
| JWKS端点 | `https://192.168.50.10:8888/sso/jwks` |
| OIDC发现 | `https://192.168.50.10:8888/sso/.well-known/openid-configuration` |

## 4. 授权码流程

### 步骤1: 发起授权请求
```
GET https://192.168.50.10:8888/sso/authorize?
    client_id=client-d56UMow9IjKknOnj&
    response_type=code&
    redirect_uri=http://localhost:3000/callback&
    scope=openid%20profile%20email&
    state=xyz123
```

### 步骤2: 用户登录授权
- 用户在NAS Dashboard登录页面输入用户名和密码
- 确认授权应用访问其信息

### 步骤3: 获取授权码
用户授权后，重定向到：
```
http://localhost:3000/callback?code=AUTH_CODE&state=xyz123
```

### 步骤4: 换取访问令牌
```bash
curl -k -X POST https://192.168.50.10:8888/sso/token \
  -H "Content-Type: application/json" \
  -d '{
    "grant_type": "authorization_code",
    "code": "AUTH_CODE",
    "redirect_uri": "http://localhost:3000/callback",
    "client_id": "client-d56UMow9IjKknOnj",
    "client_secret": "im0a7SggwKjCplNdGY6zjuLhGOKeXWBC"
  }'
```

响应：
```json
{
  "access_token": "eyJhbGc...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "refresh_token": "eyJhbGc..."
}
```

### 步骤5: 获取用户信息
```bash
curl -k -H "Authorization: Bearer ACCESS_TOKEN" \
  https://192.168.50.10:8888/sso/userinfo
```

响应：
```json
{
  "sub": "1",
  "name": "admin",
  "email": "admin@example.com",
  "preferred_username": "admin"
}
```

## 5. 测试OAuth功能

### 方法1: 使用测试页面
1. 打开 `oauth-test-example.html`
2. 点击"使用NAS账号登录"
3. 输入NAS Dashboard的用户名和密码
4. 授权后查看返回的授权码

### 方法2: 使用Postman/curl
```bash
# 1. 发起授权（在浏览器中执行）
https://192.168.50.10:8888/sso/authorize?client_id=client-d56UMow9IjKknOnj&response_type=code&redirect_uri=http://localhost:3000/callback&scope=openid%20profile%20email

# 2. 从回调URL中获取code参数

# 3. 换取token
curl -k -X POST https://192.168.50.10:8888/sso/token \
  -d "grant_type=authorization_code&code=YOUR_CODE&redirect_uri=http://localhost:3000/callback&client_id=client-d56UMow9IjKknOnj&client_secret=im0a7SggwKjCplNdGY6zjuLhGOKeXWBC"

# 4. 获取用户信息
curl -k -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  https://192.168.50.10:8888/sso/userinfo
```

## 6. 集成到您的应用

### 示例代码 (JavaScript)
```javascript
// 发起授权
function login() {
    const authUrl = 'https://192.168.50.10:8888/sso/authorize?' +
        'client_id=client-d56UMow9IjKknOnj&' +
        'response_type=code&' +
        'redirect_uri=' + encodeURIComponent('http://your-app.com/callback') + '&' +
        'scope=openid profile email';
    window.location.href = authUrl;
}

// 处理回调
const urlParams = new URLSearchParams(window.location.search);
const code = urlParams.get('code');

if (code) {
    // 用code换取token
    fetch('https://192.168.50.10:8888/sso/token', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            grant_type: 'authorization_code',
            code: code,
            redirect_uri: 'http://your-app.com/callback',
            client_id: 'client-d56UMow9IjKknOnj',
            client_secret: 'im0a7SggwKjCplNdGY6zjuLhGOKeXWBC'
        })
    })
    .then(r => r.json())
    .then(data => {
        // 保存token
        localStorage.setItem('access_token', data.access_token);

        // 获取用户信息
        return fetch('https://192.168.50.10:8888/sso/userinfo', {
            headers: { 'Authorization': 'Bearer ' + data.access_token }
        });
    })
    .then(r => r.json())
    .then(user => {
        console.log('用户信息:', user);
    });
}
```

## 7. 注意事项

### HTTPS证书
后端使用自签名证书，客户端需要：
- 在curl中使用 `-k` 参数跳过证书验证
- 在应用中配置允许自签名证书

### 修改客户端配置
如果需要添加新的OAuth客户端：
1. 访问 SSO管理界面
2. 点击"添加客户端"
3. 填写客户端名称和重定向URI

### 用户权限
- 所有已注册用户都可以通过OAuth登录
- 用户信息包括：用户名、邮箱等基本信息
