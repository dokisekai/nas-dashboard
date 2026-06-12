/**
 * WebSocket 实时监控服务
 * 连接到 /ws/monitor 获取实时监控数据
 */
import { ref, reactive } from 'vue'

interface MonitorData {
  type: 'cpu' | 'memory' | 'disk' | 'network'
  timestamp: number
  data: any
}

type MessageCallback = (data: MonitorData) => void
type ConnectionCallback = (connected: boolean) => void

class WebSocketService {
  private ws: WebSocket | null = null
  private url: string
  private reconnectTimer: number | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 2000
  private manualClose = false

  public isConnected = ref(false)
  public messageCallbacks: Map<string, Set<MessageCallback>> = new Map()
  public connectionCallbacks: Set<ConnectionCallback> = new Set()

  constructor() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.VITE_WS_URL || window.location.host
    this.url = `${protocol}//${host}/ws/monitor`
  }

  /**
   * 连接到 WebSocket 服务器
   */
  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.manualClose = false
        this.ws = new WebSocket(this.url)

        this.ws.onopen = () => {
          console.log('[WebSocket] 已连接到', this.url)
          this.isConnected.value = true
          this.reconnectAttempts = 0
          this.notifyConnectionCallbacks(true)
          resolve()
        }

        this.ws.onmessage = (event) => {
          try {
            const data: MonitorData = JSON.parse(event.data)
            this.notifyMessageCallbacks(data)
          } catch (error) {
            console.error('[WebSocket] 解析消息失败:', error)
          }
        }

        this.ws.onerror = (error) => {
          console.error('[WebSocket] 连接错误:', error)
          this.isConnected.value = false
          this.notifyConnectionCallbacks(false)
        }

        this.ws.onclose = (event) => {
          console.log('[WebSocket] 连接已关闭:', event.code, event.reason)
          this.isConnected.value = false
          this.notifyConnectionCallbacks(false)

          // 非手动关闭时自动重连
          if (!this.manualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect()
          }
        }
      } catch (error) {
        console.error('[WebSocket] 连接失败:', error)
        reject(error)
      }
    })
  }

  /**
   * 断开连接
   */
  disconnect(): void {
    this.manualClose = true
    if (this.reconnectTimer) {
      window.clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.isConnected.value = false
  }

  /**
   * 订阅特定类型的消息
   */
  subscribe(type: string, callback: MessageCallback): () => void {
    if (!this.messageCallbacks.has(type)) {
      this.messageCallbacks.set(type, new Set())
    }
    this.messageCallbacks.get(type)!.add(callback)

    // 返回取消订阅函数
    return () => {
      const callbacks = this.messageCallbacks.get(type)
      if (callbacks) {
        callbacks.delete(callback)
        if (callbacks.size === 0) {
          this.messageCallbacks.delete(type)
        }
      }
    }
  }

  /**
   * 订阅所有消息
   */
  subscribeAll(callback: MessageCallback): () => void {
    return this.subscribe('*', callback)
  }

  /**
   * 订阅连接状态变化
   */
  onConnectionChange(callback: ConnectionCallback): () => void {
    this.connectionCallbacks.add(callback)
    return () => {
      this.connectionCallbacks.delete(callback)
    }
  }

  /**
   * 发送消息到服务器
   */
  send(data: any): void {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('[WebSocket] 未连接，无法发送消息')
    }
  }

  /**
   * 调度重连
   */
  private scheduleReconnect(): void {
    if (this.reconnectTimer) return

    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts)
    console.log(`[WebSocket] 将在 ${delay}ms 后重连 (尝试 ${this.reconnectAttempts + 1}/${this.maxReconnectAttempts})`)

    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      this.reconnectAttempts++
      this.connect().catch(() => {
        // 连接失败会自动触发下一次重连
      })
    }, delay)
  }

  /**
   * 通知消息回调
   */
  private notifyMessageCallbacks(data: MonitorData): void {
    // 通知特定类型的订阅者
    const typeCallbacks = this.messageCallbacks.get(data.type)
    if (typeCallbacks) {
      typeCallbacks.forEach(callback => callback(data))
    }

    // 通知所有消息的订阅者
    const allCallbacks = this.messageCallbacks.get('*')
    if (allCallbacks) {
      allCallbacks.forEach(callback => callback(data))
    }
  }

  /**
   * 通知连接状态回调
   */
  private notifyConnectionCallbacks(connected: boolean): void {
    this.connectionCallbacks.forEach(callback => callback(connected))
  }
}

// 导出单例
export const wsService = new WebSocketService()

// 导出 composable
export function useWebSocket() {
  return {
    isConnected: wsService.isConnected,
    connect: () => wsService.connect(),
    disconnect: () => wsService.disconnect(),
    subscribe: (type: string, callback: MessageCallback) => wsService.subscribe(type, callback),
    subscribeAll: (callback: MessageCallback) => wsService.subscribeAll(callback),
    onConnectionChange: (callback: ConnectionCallback) => wsService.onConnectionChange(callback),
    send: (data: any) => wsService.send(data),
  }
}
