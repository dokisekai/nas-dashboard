import { ref, onUnmounted } from 'vue'

export interface WebSocketData {
  cpu?: {
    usage: number
    cores: number
    model: string
    load: number
    perCore: number[]
  }
  memory?: {
    total: number
    used: number
    free: number
    available: number
    percent: number
  }
  disk?: {
    disks: Array<{
      device: string
      mountpoint: string
      total: number
      free: number
      used: number
      usedPercent: number
      fstype: string
      readSpeed: number
      writeSpeed: number
    }>
  }
  network?: {
    interfaces: Array<{
      name: string
      hardwareAddr: string
      up: boolean
      addresses: string[]
      bytesSent: number
      bytesRecv: number
      sentSpeed: number
      recvSpeed: number
    }>
  }
  time: number
}

export function useWebSocket(url: string) {
  const data = ref<WebSocketData | null>(null)
  const connected = ref(false)
  const error = ref<string | null>(null)
  let ws: WebSocket | null = null
  let reconnectTimer: number | null = null

  const connect = () => {
    try {
      ws = new WebSocket(url)

      ws.onopen = () => {
        connected.value = true
        error.value = null
        console.log('WebSocket connected')
      }

      ws.onmessage = (event) => {
        try {
          const parsed = JSON.parse(event.data)
          data.value = parsed
        } catch (e) {
          console.error('Failed to parse WebSocket data:', e)
        }
      }

      ws.onerror = (e) => {
        console.error('WebSocket error:', e)
        error.value = '连接错误'
      }

      ws.onclose = () => {
        connected.value = false
        console.log('WebSocket disconnected')
        // 尝试重连
        reconnectTimer = window.setTimeout(() => {
          console.log('Attempting to reconnect...')
          connect()
        }, 5000)
      }
    } catch (e) {
      error.value = '连接失败'
      console.error('WebSocket connection failed:', e)
    }
  }

  const disconnect = () => {
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    if (ws) {
      ws.close()
      ws = null
    }
    connected.value = false
  }

  onUnmounted(() => {
    disconnect()
  })

  return {
    data,
    connected,
    error,
    connect,
    disconnect,
  }
}
