import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<any>(null)

  // 开发模式：自动登录以便测试液态玻璃桌面
  const isDevMode = () => {
    return import.meta.env.DEV || localStorage.getItem('dev_mode') === 'true'
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUser = (userData: any) => {
    user.value = userData
  }

  const clearToken = () => {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  const isLoggedIn = () => {
    return !!token.value
  }

  // 尝试自动登录（仅限开发模式）
  const tryAutoLogin = () => {
    if (isDevMode() && !token.value) {
      const devToken = 'dev-token-' + Date.now()
      setToken(devToken)
      setUser({ username: 'dev_user', role: 'admin' })
      return true
    }
    return !!token.value
  }

  return {
    token,
    user,
    setToken,
    setUser,
    clearToken,
    isLoggedIn,
    tryAutoLogin,
    isDevMode,
  }
})
