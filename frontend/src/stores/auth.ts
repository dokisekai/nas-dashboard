import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const user = ref<any>(null)

  // 从localStorage加载用户信息
  const loadUserFromStorage = () => {
    const userStr = localStorage.getItem('user')
    if (userStr) {
      try {
        user.value = JSON.parse(userStr)
      } catch (error) {
        console.error('Failed to parse user data from localStorage:', error)
        localStorage.removeItem('user')
      }
    }
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
    console.log('Token saved to localStorage and store:', newToken.substring(0, 20) + '...')
  }

  const setRefreshToken = (newRefreshToken: string) => {
    refreshToken.value = newRefreshToken
    localStorage.setItem('refreshToken', newRefreshToken)
  }

  const setUser = (userData: any) => {
    user.value = userData
    localStorage.setItem('user', JSON.stringify(userData))
  }

  const clearToken = () => {
    token.value = null
    refreshToken.value = null
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('user')
    localStorage.removeItem('rememberMe')
  }

  const isLoggedIn = () => {
    const hasToken = !!token.value
    const hasLocalStorageToken = !!localStorage.getItem('token')
    console.log('isLoggedIn check - token.value:', !!token.value, 'localStorage token:', hasLocalStorageToken)
    return hasToken || hasLocalStorageToken
  }

  // 登录
  const login = async (username: string, password: string) => {
    try {
      const response = await api.post('/auth/login', {
        username,
        password
      })

      // 保存token和用户信息
      setToken(response.token)
      setRefreshToken(response.refreshToken)
      setUser(response.user)

      return { success: true, user: response.user }
    } catch (error: any) {
      console.error('Login failed:', error)
      return {
        success: false,
        message: error.message || '登录失败，请检查用户名和密码'
      }
    }
  }

  // 登出
  const logout = () => {
    clearToken()
  }

  // 刷新token
  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await api.post('/auth/refresh', {
        refreshToken: refreshToken.value
      })

      setToken(response.token)
      setRefreshToken(response.refreshToken)

      return { success: true }
    } catch (error: any) {
      console.error('Token refresh failed:', error)
      // 刷新token失败，清除所有认证信息
      clearToken()
      return {
        success: false,
        message: error.message || 'Token刷新失败，请重新登录'
      }
    }
  }

  // 初始化时加载用户信息
  loadUserFromStorage()

  // 如果localStorage有token但store没有，重新加载
  if (!token.value && localStorage.getItem('token')) {
    token.value = localStorage.getItem('token')
    console.log('Reloaded token from localStorage:', token.value?.substring(0, 20) + '...')
  }

  return {
    token,
    refreshToken,
    user,
    setToken,
    setRefreshToken,
    setUser,
    clearToken,
    isLoggedIn,
    login,
    logout,
    refreshAccessToken
  }
})