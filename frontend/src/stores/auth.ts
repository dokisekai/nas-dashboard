import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refreshToken'))
  const user = ref<any>(null)

  const loadUserFromStorage = () => {
    const userStr = localStorage.getItem('user')
    if (!userStr) return
    try {
      user.value = JSON.parse(userStr)
    } catch {
      localStorage.removeItem('user')
    }
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
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

  const isLoggedIn = () => !!token.value || !!localStorage.getItem('token')

  const login = async (username: string, password: string) => {
    try {
      const response = await api.post('/auth/login', { username, password })
      setToken(response.token)
      setRefreshToken(response.refreshToken)
      setUser(response.user)
      return { success: true, user: response.user }
    } catch (error: any) {
      return {
        success: false,
        message: error.message || '登录失败，请检查用户名和密码',
      }
    }
  }

  const logout = () => clearToken()

  const refreshAccessToken = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }
    try {
      const response = await api.post('/auth/refresh', { refreshToken: refreshToken.value })
      setToken(response.token)
      setRefreshToken(response.refreshToken)
      return { success: true }
    } catch (error: any) {
      clearToken()
      return {
        success: false,
        message: error.message || 'Token 刷新失败，请重新登录',
      }
    }
  }

  loadUserFromStorage()

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
    refreshAccessToken,
  }
})
