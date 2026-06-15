/**
 * 统一的日志系统
 * 提供环境变量控制的调试输出，避免生产环境的console.log污染
 */

class Logger {
  private debugMode: boolean
  private isProduction: boolean

  constructor() {
    // 检查环境变量
    this.debugMode = import.meta.env.VITE_DEBUG === 'true' || import.meta.env.DEV
    this.isProduction = import.meta.env.PROD
  }

  /**
   * 输出信息日志
   * @param message 日志消息
   * @param args 附加参数
   */
  info(message: string, ...args: any[]) {
    if (this.shouldLog('info')) {
      console.log(`[INFO] ${message}`, ...args)
    }
  }

  /**
   * 输出警告日志
   * @param message 警告消息
   * @param args 附加参数
   */
  warn(message: string, ...args: any[]) {
    if (this.shouldLog('warn')) {
      console.warn(`[WARN] ${message}`, ...args)
    }
  }

  /**
   * 输出错误日志
   * @param message 错误消息
   * @param args 附加参数
   */
  error(message: string, ...args: any[]) {
    // 错误日志始终输出
    console.error(`[ERROR] ${message}`, ...args)

    // 可选：发送错误到后端日志系统
    this.sendToBackend('error', message, args)
  }

  /**
   * 输出调试日志
   * @param message 调试消息
   * @param args 附加参数
   */
  debug(message: string, ...args: any[]) {
    if (this.debugMode) {
      console.log(`[DEBUG] ${message}`, ...args)
    }
  }

  /**
   * 输出成功日志
   * @param message 成功消息
   * @param args 附加参数
   */
  success(message: string, ...args: any[]) {
    if (this.shouldLog('success')) {
      console.log(`%c[SUCCESS] ${message}`, 'color: green; font-weight: bold', ...args)
    }
  }

  /**
   * 网络请求日志
   * @param method 请求方法
   * @param url 请求URL
   * @param data 请求数据
   */
  apiRequest(method: string, url: string, data?: any) {
    if (this.debugMode) {
      console.log(`[API] ${method} ${url}`, data || '')
    }
  }

  /**
   * 网络响应日志
   * @param url 请求URL
   * @param status 响应状态
   * @param data 响应数据
   */
  apiResponse(url: string, status: number, data?: any) {
    if (this.debugMode) {
      console.log(`[API] Response ${status} ${url}`, data || '')
    }
  }

  /**
   * 判断是否应该输出日志
   * @param level 日志级别
   */
  private shouldLog(level: string): boolean {
    if (this.isProduction) {
      // 生产环境只输出错误和警告
      return ['error', 'warn'].includes(level)
    }
    return true
  }

  /**
   * 发送日志到后端（可选实现）
   * @param level 日志级别
   * @param message 日志消息
   * @param args 附加参数
   */
  private sendToBackend(level: string, message: string, args: any[]) {
    // 这里可以实现将日志发送到后端的功能
    // 但为了避免循环依赖和性能问题，暂时不实现
    // 如果需要，可以在这里调用日志API

    /*
    try {
      if (this.isProduction && level === 'error') {
        // 只在生产环境发送错误日志
        fetch('/api/logs', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            level,
            message,
            args,
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent,
            url: window.location.href
          })
        }).catch(() => {
          // 忽略日志发送失败
        })
      }
    } catch (error) {
      // 忽略日志系统自身的错误
    }
    */
  }

  /**
   * 分组开始
   * @param label 分组标签
   */
  group(label: string) {
    if (this.debugMode) {
      console.group(`[GROUP] ${label}`)
    }
  }

  /**
   * 分组结束
   */
  groupEnd() {
    if (this.debugMode) {
      console.groupEnd()
    }
  }

  /**
   * 性能计时开始
   * @param label 计时标签
   */
  timeStart(label: string) {
    if (this.debugMode) {
      console.time(`[TIME] ${label}`)
    }
  }

  /**
   * 性能计时结束
   * @param label 计时标签
   */
  timeEnd(label: string) {
    if (this.debugMode) {
      console.timeEnd(`[TIME] ${label}`)
    }
  }
}

// 创建单例实例
export const logger = new Logger()

// 默认导出
export default logger

// 便捷导出
export const log = {
  info: (message: string, ...args: any[]) => logger.info(message, ...args),
  warn: (message: string, ...args: any[]) => logger.warn(message, ...args),
  error: (message: string, ...args: any[]) => logger.error(message, ...args),
  debug: (message: string, ...args: any[]) => logger.debug(message, ...args),
  success: (message: string, ...args: any[]) => logger.success(message, ...args),
  api: {
    request: (method: string, url: string, data?: any) => logger.apiRequest(method, url, data),
    response: (url: string, status: number, data?: any) => logger.apiResponse(url, status, data)
  }
}