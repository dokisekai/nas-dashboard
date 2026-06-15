import axios from 'axios'

// Debug mode - log all API calls
const DEBUG = true // 强制开启DEBUG模式
console.log('[API CLIENT INIT] - Creating axios instance with DEBUG mode')
console.log('[API CLIENT INIT] - Base URL:', import.meta.env.VITE_API_URL || '')

// 创建一个通用的请求函数，自动处理返回数据
const instance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 增强类型支持
const api = {
  get: <T = any>(url: string, config?: any) => instance.get(url, config) as Promise<T>,
  post: <T = any>(url: string, data?: any, config?: any) => instance.post(url, data, config) as Promise<T>,
  put: <T = any>(url: string, data?: any, config?: any) => instance.put(url, data, config) as Promise<T>,
  delete: <T = any>(url: string, config?: any) => instance.delete(url, config) as Promise<T>,
  patch: <T = any>(url: string, data?: any, config?: any) => instance.patch(url, data, config) as Promise<T>,
}

console.log('[API CLIENT INIT] - Axios instance created')

// 请求拦截器 - 添加 JWT token
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    if (DEBUG) {
      console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`)
      console.log(`[Token Debug] Token exists: ${!!token}`)
      const authHeader = config.headers.Authorization
      const authStr = typeof authHeader === 'string' ? authHeader : ''
      console.log(`[Token Debug] Authorization header:`, authStr.substring(0, 30) + '...')
      console.log(`[Request Data]`, config.data || config.params)
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
instance.interceptors.response.use(
  (response) => {
    if (DEBUG) {
      console.log(`[API Response] ${response.config.url}`, response.data)
    }
    return response.data
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      const data = error.response.data

      if (status === 401) {
        // 清除本地存储并重定向到登录页
        console.warn('401 Unauthorized - clearing token and redirecting to login')

        // 检查是否已经在登录页，避免循环重定向
        if (window.location.pathname === '/login') {
          console.log('Already on login page, skipping redirect')
          return Promise.reject({
            message: '请先登录',
            originalError: data || error.message,
            status: 401
          })
        }

        // 清除本地存储
        localStorage.removeItem('token')
        localStorage.removeItem('user')

        // 重定向到登录页
        console.log('Redirecting to login page due to 401')
        window.location.href = '/login'

        return Promise.reject({
          message: '登录已过期，请重新登录',
          originalError: data || error.message,
          status: 401
        })
      } else if (status === 500) {
        // 服务器内部错误，提供友好的错误信息
        console.error('Server Error:', data || error.message)
        return Promise.reject({
          message: '服务器内部错误，请稍后重试',
          originalError: data || error.message,
          status: 500
        })
      } else if (status === 404) {
        // 资源不存在
        return Promise.reject({
          message: '请求的资源不存在',
          originalError: data || error.message,
          status: 404
        })
      } else if (status === 403) {
        // 权限不足
        return Promise.reject({
          message: '权限不足，无法访问此资源',
          originalError: data || error.message,
          status: 403
        })
      }
    }

    // 网络错误或其他错误
    return Promise.reject(error.message || '网络错误，请检查连接')
  }
)

export default api
