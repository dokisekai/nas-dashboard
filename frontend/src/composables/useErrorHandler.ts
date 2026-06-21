import { ref, computed } from 'vue'
import { useNotificationStore } from '@/stores/notification'

export interface ApiError {
  message: string
  code?: number
  details?: any
  originalError?: any
}

// 错误类型枚举
export enum ErrorType {
  NETWORK = 'network',
  VALIDATION = 'validation',
  AUTHENTICATION = 'authentication',
  AUTHORIZATION = 'authorization',
  NOT_FOUND = 'not_found',
  SERVER = 'server',
  UNKNOWN = 'unknown'
}

// 错误处理配置
export interface ErrorHandlerConfig {
  showNotification?: boolean
  logToConsole?: boolean
  retryable?: boolean
  customHandler?: (error: ApiError) => void
}

// 全局错误处理状态
const globalError = ref<ApiError | null>(null)
const errorHistory = ref<ApiError[]>([])

export function useErrorHandler() {
  const notificationStore = useNotificationStore()

  // 分类错误类型
  function classifyError(error: any): ErrorType {
    if (!error) return ErrorType.UNKNOWN

    if (error.message?.includes('network') || error.message?.includes('fetch')) {
      return ErrorType.NETWORK
    }

    if (error.status === 401 || error.code === 401) {
      return ErrorType.AUTHENTICATION
    }

    if (error.status === 403 || error.code === 403) {
      return ErrorType.AUTHORIZATION
    }

    if (error.status === 404 || error.code === 404) {
      return ErrorType.NOT_FOUND
    }

    if (error.status === 400 || error.code === 400) {
      return ErrorType.VALIDATION
    }

    if (error.status >= 500 || error.code >= 500) {
      return ErrorType.SERVER
    }

    return ErrorType.UNKNOWN
  }

  // 格式化错误消息
  function formatErrorMessage(error: any, errorType: ErrorType): string {
    // 如果错误对象已经包含格式化的消息，直接使用
    if (error.message && typeof error.message === 'string') {
      return error.message
    }

    // 根据错误类型提供友好的默认消息
    switch (errorType) {
      case ErrorType.NETWORK:
        return '网络连接失败，请检查网络设置'
      case ErrorType.AUTHENTICATION:
        return '登录已过期，请重新登录'
      case ErrorType.AUTHORIZATION:
        return '权限不足，无法访问此资源'
      case ErrorType.NOT_FOUND:
        return '请求的资源不存在'
      case ErrorType.VALIDATION:
        return '请求参数错误，请检查输入'
      case ErrorType.SERVER:
        return '服务器错误，请稍后重试'
      default:
        return '未知错误，请重试'
    }
  }

  // 处理错误
  function handleError(
    error: any,
    config: ErrorHandlerConfig = {}
  ): ApiError {
    const {
      showNotification = true,
      logToConsole = true,
      customHandler
    } = config

    const errorType = classifyError(error)
    const message = formatErrorMessage(error, errorType)

    const apiError: ApiError = {
      message,
      code: error.status || error.code,
      details: error.details || error.originalError,
      originalError: error
    }

    // 记录到历史
    errorHistory.value.unshift({
      ...apiError,
      originalError: undefined // 不保存原始错误对象以避免内存泄漏
    })

    // 只保留最近50条错误记录
    if (errorHistory.value.length > 50) {
      errorHistory.value = errorHistory.value.slice(0, 50)
    }

    // 设置全局错误
    globalError.value = apiError

    // 控制台日志
    if (logToConsole) {
      console.error('API Error:', {
        type: errorType,
        error: apiError,
        original: error
      })
    }

    // 显示通知
    if (showNotification) {
      const notificationType = getNotificationType(errorType)
      notificationStore.addNotification({
        type: notificationType,
        title: '错误',
        message: apiError.message,
        duration: 5000
      })
    }

    // 自定义处理
    if (customHandler) {
      customHandler(apiError)
    }

    // 特殊处理认证错误
    if (errorType === ErrorType.AUTHENTICATION) {
      handleAuthenticationError()
    }

    return apiError
  }

  // 获取通知类型
  function getNotificationType(errorType: ErrorType): 'error' | 'warning' | 'info' {
    switch (errorType) {
      case ErrorType.AUTHENTICATION:
      case ErrorType.AUTHORIZATION:
      case ErrorType.SERVER:
        return 'error'
      case ErrorType.NETWORK:
      case ErrorType.VALIDATION:
        return 'warning'
      default:
        return 'info'
    }
  }

  // 处理认证错误
  function handleAuthenticationError() {
    // 清除认证信息
    localStorage.removeItem('token')
    localStorage.removeItem('user')

    // 延迟跳转，避免打扰用户
    setTimeout(() => {
      if (window.location.pathname !== '/login') {
        window.location.href = '/login'
      }
    }, 2000)
  }

  // 清除错误
  function clearError() {
    globalError.value = null
  }

  // 清除错误历史
  function clearErrorHistory() {
    errorHistory.value = []
  }

  // 计算属性
  const hasError = computed(() => globalError.value !== null)
  const recentErrors = computed(() => errorHistory.value)
  const errorCount = computed(() => errorHistory.value.length)

  return {
    // 状态
    globalError,
    errorHistory,
    hasError,
    recentErrors,
    errorCount,

    // 方法
    handleError,
    classifyError,
    formatErrorMessage,
    clearError,
    clearErrorHistory
  }
}