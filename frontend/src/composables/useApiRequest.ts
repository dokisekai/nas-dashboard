import { ref, computed } from 'vue'
import api from '@/api/client'
import { useErrorHandler, type ApiError, type ErrorHandlerConfig } from './useErrorHandler'

export interface RequestOptions {
  showLoading?: boolean
  showErrorNotification?: boolean
  retryOnFailure?: boolean
  maxRetries?: number
  retryDelay?: number
  onSuccess?: (data: any) => void
  onError?: (error: ApiError) => void
  onFinally?: () => void
}

export interface ApiRequestResult<T> {
  data: ref<T | null>
  loading: ref<boolean>
  error: ref<ApiError | null>
  execute: () => Promise<T | null>
  reset: () => void
}

export function useApiRequest<T = any>(
  requestFn: () => Promise<any>,
  options: RequestOptions = {}
): ApiRequestResult<T> {
  const {
    showLoading = true,
    showErrorNotification = true,
    retryOnFailure = false,
    maxRetries = 3,
    retryDelay = 1000,
    onSuccess,
    onError,
    onFinally
  } = options

  const data = ref<T | null>(null)
  const loading = ref<boolean>(false)
  const error = ref<ApiError | null>(null)
  const retryCount = ref<number>(0)

  const { handleError } = useErrorHandler()

  // 执行请求
  async function execute(): Promise<T | null> {
    loading.value = true
    error.value = null

    if (showLoading) {
      // 这里可以集成全局loading状态
    }

    try {
      const response = await requestFn()
      data.value = response as T
      retryCount.value = 0

      if (onSuccess) {
        onSuccess(response)
      }

      return response as T
    } catch (err: any) {
      const apiError = handleError(err, {
        showNotification: showErrorNotification,
        logToConsole: true
      })

      error.value = apiError

      // 重试逻辑
      if (retryOnFailure && retryCount.value < maxRetries && isRetryableError(apiError)) {
        retryCount.value++
        console.log(`Retrying request (${retryCount.value}/${maxRetries})...`)

        await new Promise(resolve => setTimeout(resolve, retryDelay))
        return execute()
      }

      if (onError) {
        onError(apiError)
      }

      return null
    } finally {
      loading.value = false
      if (onFinally) {
        onFinally()
      }
    }
  }

  // 判断是否可重试的错误
  function isRetryableError(apiError: ApiError): boolean {
    // 网络错误和5xx服务器错误可以重试
    const retryableCodes = [0, 500, 502, 503, 504]
    return retryableCodes.includes(apiError.code || 0) ||
           apiError.message?.toLowerCase().includes('network')
  }

  // 重置状态
  function reset() {
    data.value = null
    loading.value = false
    error.value = null
    retryCount.value = 0
  }

  return {
    data,
    loading,
    error,
    execute,
    reset
  }
}

// 便捷的API调用hooks
export function useApiGet<T = any>(url: string, options?: RequestOptions) {
  return useApiRequest<T>(() => api.get(url), options)
}

export function useApiPost<T = any>(url: string, body?: any, options?: RequestOptions) {
  return useApiRequest<T>(() => api.post(url, body), options)
}

export function useApiPut<T = any>(url: string, body?: any, options?: RequestOptions) {
  return useApiRequest<T>(() => api.put(url, body), options)
}

export function useApiDelete<T = any>(url: string, options?: RequestOptions) {
  return useApiRequest<T>(() => api.delete(url), options)
}

// 批量API请求hook
export function useBatchApiRequest<T = any>(
  requestFns: Array<() => Promise<any>>,
  options: RequestOptions = {}
) {
  const results = ref<T[]>([])
  const loading = ref(false)
  const errors = ref<ApiError[]>([])

  async function executeAll(): Promise<T[]> {
    loading.value = true
    errors.value = []

    try {
      const responses = await Promise.allSettled(
        requestFns.map(fn => fn())
      )

      results.value = responses.map((response, index) => {
        if (response.status === 'fulfilled') {
          return response.value as T
        } else {
          const apiError = handleError(response.reason, {
            showNotification: false,
            logToConsole: true
          })
          errors.value[index] = apiError
          return null as T
        }
      })

      return results.value
    } finally {
      loading.value = false
    }
  }

  return {
    results,
    loading,
    errors,
    executeAll
  }
}