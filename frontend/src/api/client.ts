import axios from 'axios'

// Debug mode - log all API calls
const DEBUG = import.meta.env.VITE_DEBUG === 'true'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器 - 添加 JWT token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    if (DEBUG) {
      console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, config.data || config.params)
    }
    return config
  },
  (error) => {
    if (DEBUG) {
      console.error('[API Request Error]', error)
    }
    return Promise.reject(error)
  }
)

// 响应拦截器 - 处理错误
api.interceptors.response.use(
  (response) => {
    if (DEBUG) {
      console.log(`[API Response] ${response.config.url}`, response.data)
    }
    return response.data
  },
  (error) => {
    if (DEBUG) {
      console.error('[API Response Error]', error.response?.data || error.message)
    }
    if (error.response?.status === 401) {
      // Token 过期，清除本地存储并跳转登录
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error.message)
  }
)

export default api
