import axios from 'axios'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

const api = {
  get: <T = any>(url: string, config?: any) => instance.get(url, config) as Promise<T>,
  post: <T = any>(url: string, data?: any, config?: any) => instance.post(url, data, config) as Promise<T>,
  put: <T = any>(url: string, data?: any, config?: any) => instance.put(url, data, config) as Promise<T>,
  delete: <T = any>(url: string, config?: any) => instance.delete(url, config) as Promise<T>,
  patch: <T = any>(url: string, data?: any, config?: any) => instance.patch(url, data, config) as Promise<T>,
}

instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

instance.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (!error.response) {
      return Promise.reject(error.message || '网络错误，请检查连接')
    }

    const status = error.response.status
    const data = error.response.data

    if (status === 401) {
      if (window.location.pathname === '/login') {
        return Promise.reject({
          message: '请先登录',
          originalError: data || error.message,
          status: 401,
        })
      }
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
      return Promise.reject({
        message: '登录已过期，请重新登录',
        originalError: data || error.message,
        status: 401,
      })
    }

    const messages: Record<number, string> = {
      403: '权限不足，无法访问此资源',
      404: '请求的资源不存在',
      500: '服务器内部错误，请稍后重试',
    }

    return Promise.reject({
      message: messages[status] || (typeof data === 'string' ? data : error.message),
      originalError: data || error.message,
      status,
    })
  },
)

export default api
