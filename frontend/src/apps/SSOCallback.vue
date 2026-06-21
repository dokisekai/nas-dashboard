<template>
  <div class="sso-callback-container">
    <div class="loading-card">
      <!-- 加载动画 -->
      <div class="loading-animation">
        <div class="loading-circle"></div>
        <div class="loading-circle"></div>
        <div class="loading-circle"></div>
      </div>

      <!-- 状态信息 -->
      <div class="status-info">
        <h2 class="status-title">{{ statusTitle }}</h2>
        <p class="status-message">{{ statusMessage }}</p>
      </div>

      <!-- 进度条 -->
      <div v-if="showProgress" class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  performSSOLogin,
  clearSSOSession,
  getSSOSession,
  generateRandomState
} from '@/api/sso'

const router = useRouter()
const route = useRoute()

// 状态
const currentStep = ref(0)
const totalSteps = 4
const error = ref<string | null>(null)

// 计算属性
const progress = computed(() => {
  return (currentStep.value / totalSteps) * 100
})

const showProgress = computed(() => {
  return !error.value && currentStep.value > 0
})

const statusTitle = computed(() => {
  if (error.value) {
    return '登录失败'
  }

  switch (currentStep.value) {
    case 0:
      return '正在验证登录信息...'
    case 1:
      return '正在交换授权码...'
    case 2:
      return '正在获取用户信息...'
    case 3:
      return '正在创建会话...'
    case 4:
      return '登录成功！'
    default:
      return '正在处理...'
  }
})

const statusMessage = computed(() => {
  if (error.value) {
    return error.value
  }

  switch (currentStep.value) {
    case 0:
      return '验证SSO登录回调参数'
    case 1:
      return '正在与SSO服务器交换访问令牌'
    case 2:
      return '正在获取您的用户资料信息'
    case 3:
      return '正在创建登录会话'
    case 4:
      return '即将跳转到主页面...'
    default:
      return '请稍候'
  }
})

// 处理SSO回调
const handleSSOCallback = async () => {
  try {
    // 获取URL参数
    const code = route.query.code as string
    const state = route.query.state as string
    const errorParam = route.query.error as string
    const errorDescription = route.query.error_description as string

    // 检查是否有错误
    if (errorParam) {
      throw new Error(errorDescription || errorParam)
    }

    // 验证必需参数
    if (!code) {
      throw new Error('缺少授权码参数')
    }

    currentStep.value = 1 // 开始交换授权码

    // 执行SSO登录流程
    await performSSOLogin(code, state)

    currentStep.value = 4 // 登录成功

    // 显示成功消息
    ElMessage.success('登录成功！')

    // 延迟跳转，让用户看到成功状态
    setTimeout(() => {
      // 跳转到首页或原来想访问的页面
      const redirectPath = sessionStorage.getItem('redirect_path') || '/'
      sessionStorage.removeItem('redirect_path')
      router.push(redirectPath)
    }, 1500)

  } catch (err) {
    console.error('SSO callback failed:', err)
    error.value = err instanceof Error ? err.message : '登录过程中发生错误'

    // 清理失败的会话数据
    clearSSOSession()

    // 显示错误消息
    ElMessage.error(error.value)

    // 延迟跳转到登录页
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  }
}

onMounted(() => {
  // 短暂延迟以显示加载动画
  setTimeout(() => {
    handleSSOCallback()
  }, 500)
})
</script>

<style scoped>
.sso-callback-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.loading-card {
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  padding: 48px;
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.loading-animation {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  margin-bottom: 32px;
  height: 40px;
}

.loading-circle {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #667eea;
  animation: bounce 1.4s infinite ease-in-out both;
}

.loading-circle:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-circle:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes bounce {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

.status-info {
  margin-bottom: 32px;
}

.status-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 12px;
}

.status-message {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: #f0f0f0;
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 2px;
  transition: width 0.3s ease;
}

/* 错误状态 */
.loading-card.error {
  border: 2px solid #ef4444;
}

.loading-card.error .status-title {
  color: #ef4444;
}

/* 响应式设计 */
@media (max-width: 640px) {
  .loading-card {
    padding: 32px 24px;
  }

  .status-title {
    font-size: 18px;
  }

  .status-message {
    font-size: 13px;
  }
}
</style>