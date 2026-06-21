# 简化的SSO解决方案

## 方案：NAS Dashboard作为SSO入口

由于Immich的SSO配置比较复杂，我们可以使用更简单的方案：

### 实现思路
1. 用户在NAS Dashboard已登录
2. 点击Immich图标时，通过NAS Dashboard后端代理
3. 后端使用Immich API创建临时访问链接
4. 用户自动跳转到已认证的Immich会话

### 快速实现
```bash
# 修改前端跳转地址
# 指向NAS Dashboard的代理端点
http://192.168.50.10:8888/api/immich/sso-login
```

### 后端实现
- 验证用户的JWT token
- 通过Immich API创建会话
- 重定向到已认证的Immich页面

这样用户感觉像是"自动登录"，实际上是NAS Dashboard
代理了认证过程。
