<template>
  <div class="sso-login-container">
    <div class="sso-login-card">
      <!-- Logo -->
      <div class="logo-section">
        <svg class="logo-icon" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
        </svg>
        <h1 class="logo-title">NAS统一身份认证</h1>
        <p class="logo-subtitle">选择登录方式继续</p>
      </div>

      <!-- 身份提供商列表 -->
      <div class="providers-section">
        <el-button
          v-for="provider in providers"
          :key="provider.id"
          :loading="loadingProvider === provider.id"
          :class="['provider-button', provider.id]"
          @click="loginWithProvider(provider.id)"
          size="large"
          plain
        >
          <div class="provider-content">
            <div class="provider-icon">
              <!-- Google Icon -->
              <svg v-if="provider.id === 'google'" viewBox="0 0 24 24" width="24" height="24">
                <path fill="#4285f4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34a853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#fbbc05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#ea4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>

              <!-- GitHub Icon -->
              <svg v-else-if="provider.id === 'github'" viewBox="0 0 24 24" width="24" height="24" fill="#333">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>

              <!-- Microsoft Icon -->
              <svg v-else-if="provider.id === 'microsoft'" viewBox="0 0 24 24" width="24" height="24">
                <rect x="1" y="1" width="10" height="10" fill="#f25022"/>
                <rect x="1" y="13" width="10" height="10" fill="#00a4ef"/>
                <rect x="13" y="1" width="10" height="10" fill="#7fba00"/>
                <rect x="13" y="13" width="10" height="10" fill="#ffb900"/>
              </svg>

              <!-- LDAP Icon -->
              <svg v-else-if="provider.id === 'ldap'" viewBox="0 0 24 24" width="24" height="24" fill="#f39c12">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
              </svg>

              <!-- Immich Icon -->
              <svg v-else-if="provider.id === 'immich'" viewBox="0 0 24 24" width="24" height="24" fill="#10b981">
                <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
              </svg>

              <!-- Generic Provider Icon -->
              <svg v-else viewBox="0 0 24 24" width="24" height="24" fill="#666">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/>
              </svg>
            </div>
            <span class="provider-name">{{ provider.name }}</span>
          </div>
        </el-button>
      </div>

      <!-- 底部信息 -->
      <div class="footer-section">
        <p class="footer-text">
          登录即表示您同意我们的
          <el-link type="primary">服务条款</el-link>
          和
          <el-link type="primary">隐私政策</el-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'

const router = useRouter()
const route = useRoute()

// 状态
const loadingProvider = ref<string | null>(null)
const providers = ref([
  { id: 'google', name: '使用 Google 账号登录' },
  { id: 'github', name: '使用 GitHub 账号登录' },
  { id: 'microsoft', name: '使用 Microsoft 账号登录' },
  { id: 'ldap', name: '使用 LDAP/AD 账号登录' },
  { id: 'immich', name: '使用 Immich 账号登录' }
])

// SSO配置
const SSO_CONFIG = {
  issuerUrl: 'http://192.168.50.10:8888',
  clientId: 'nas-dashboard',
  redirectUri: window.location.origin + '/sso/callback',
  scope: 'openid profile email'
}

// 使用指定提供商登录
const loginWithProvider = async (providerId: string) => {
  try {
    loadingProvider.value = providerId

    // 生成随机的state参数用于安全验证
    const state = generateRandomString(32)

    // 将state存储在sessionStorage中用于回调验证
    sessionStorage.setItem('sso_state', state)
    sessionStorage.setItem('sso_provider', providerId)

    // 构建授权URL
    const authParams = new URLSearchParams({
      client_id: SSO_CONFIG.clientId,
      response_type: 'code',
      redirect_uri: SSO_CONFIG.redirectUri,
      scope: SSO_CONFIG.scope,
      state: state,
      provider: providerId
    })

    const authUrl = `${SSO_CONFIG.issuerUrl}/sso/authorize?${authParams.toString()}`

    // 重定向到SSO登录页面
    window.location.href = authUrl

  } catch (error) {
    console.error('SSO登录失败:', error)
    ElMessage.error('登录失败，请稍后重试')
    loadingProvider.value = null
  }
}

// 生成随机字符串
const generateRandomString = (length: number): string => {
  const charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''
  const randomValues = new Uint32Array(length)

  // 使用crypto.getRandomValues生成安全的随机数
  window.crypto.getRandomValues(randomValues)

  for (let i = 0; i < length; i++) {
    result += charset[randomValues[i] % charset.length]
  }

  return result
}

// 检查URL参数中是否包含错误
onMounted(() => {
  const error = route.query.error
  const errorDescription = route.query.error_description

  if (error) {
    let errorMessage = '登录失败'
    if (typeof errorDescription === 'string') {
      errorMessage = errorDescription
    } else if (typeof error === 'string') {
      errorMessage = error
    }

    ElMessage.error(errorMessage)
    // 清除URL中的错误参数
    router.replace({ query: {} })
  }
})
</script>

<style scoped>
.sso-login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.sso-login-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 40px;
  width: 100%;
  max-width: 480px;
}

.logo-section {
  text-align: center;
  margin-bottom: 40px;
}

.logo-icon {
  width: 48px;
  height: 48px;
  color: #667eea;
  margin: 0 auto 16px;
}

.logo-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.logo-subtitle {
  font-size: 16px;
  color: #666;
  margin: 0;
}

.providers-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.provider-button {
  width: 100%;
  height: 56px;
  border-radius: 12px;
  border: 2px solid #e5e7eb;
  background: white;
  color: #333;
  font-size: 16px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.provider-button:hover {
  border-color: #667eea;
  background: #f9fafb;
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.2);
}

.provider-button.google:hover {
  border-color: #4285f4;
  background: #f8f9ff;
}

.provider-button.github:hover {
  border-color: #333;
  background: #f5f5f5;
}

.provider-button.microsoft:hover {
  border-color: #00a4ef;
  background: #f0f9ff;
}

.provider-button.ldap:hover {
  border-color: #f39c12;
  background: #fffaf0;
}

.provider-button.immich:hover {
  border-color: #10b981;
  background: #f0fdf4;
}

.provider-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.provider-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.footer-section {
  text-align: center;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.footer-text {
  font-size: 14px;
  color: #666;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .sso-login-card {
    padding: 24px;
  }

  .logo-title {
    font-size: 24px;
  }

  .logo-subtitle {
    font-size: 14px;
  }

  .provider-button {
    height: 48px;
    font-size: 14px;
  }
}
</style>