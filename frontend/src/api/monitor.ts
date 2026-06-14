// Monitor API Client
import axios from 'axios'
import type {
  ProcessInfo,
  ServiceInfo,
  TemperatureInfo,
  SystemEvent,
  SystemLog,
  AlertRule,
  MonitorStats
} from '@/types/monitor'

const API_BASE = '/api/monitor'

export const monitorAPI = {
  // Process Management
  async getProcesses(): Promise<ProcessInfo[]> {
    const response = await axios.get(`${API_BASE}/processes`)
    return response.data.processes || []
  },

  async getProcess(pid: number): Promise<ProcessInfo> {
    const response = await axios.get(`${API_BASE}/processes/${pid}`)
    return response.data
  },

  async killProcess(pid: number, signal?: string): Promise<void> {
    await axios.delete(`${API_BASE}/processes/${pid}`, {
      params: { signal: signal || '15' }
    })
  },

  // Service Management
  async getServices(): Promise<ServiceInfo[]> {
    const response = await axios.get(`${API_BASE}/services`)
    return response.data.services || []
  },

  async getService(name: string): Promise<ServiceInfo> {
    const response = await axios.get(`${API_BASE}/services/${name}`)
    return response.data
  },

  async startService(name: string): Promise<void> {
    await axios.post(`${API_BASE}/services/${name}/start`)
  },

  async stopService(name: string): Promise<void> {
    await axios.post(`${API_BASE}/services/${name}/stop`)
  },

  async restartService(name: string): Promise<void> {
    await axios.post(`${API_BASE}/services/${name}/restart`)
  },

  async enableService(name: string): Promise<void> {
    await axios.post(`${API_BASE}/services/${name}/enable`)
  },

  async disableService(name: string): Promise<void> {
    await axios.post(`${API_BASE}/services/${name}/disable`)
  },

  // Temperature Monitoring
  async getTemperature(): Promise<TemperatureInfo> {
    const response = await axios.get(`${API_BASE}/temperature`)
    return response.data
  },

  // Events and Logs
  async getEvents(params?: {
    page?: number
    limit?: number
    type?: string
    source?: string
  }): Promise<{ events: SystemEvent[]; total: number; page: number; limit: number }> {
    const response = await axios.get(`${API_BASE}/events`, { params })
    return response.data
  },

  async getLogs(params?: {
    page?: number
    limit?: number
    level?: string
    component?: string
  }): Promise<{ logs: SystemLog[]; total: number; page: number; limit: number }> {
    const response = await axios.get(`${API_BASE}/logs`, { params })
    return response.data
  },

  async clearLogs(params?: { level?: string; component?: string }): Promise<void> {
    await axios.post(`${API_BASE}/logs/clear`, null, { params })
  },

  // Alert Management
  async getAlerts(): Promise<AlertRule[]> {
    const response = await axios.get(`${API_BASE}/alerts`)
    return response.data.alerts || []
  },

  async createAlert(alert: Omit<AlertRule, 'id' | 'lastTriggered' | 'createdAt' | 'updatedAt'>): Promise<AlertRule> {
    const response = await axios.post(`${API_BASE}/alerts`, alert)
    return response.data
  },

  async updateAlert(id: number, alert: Partial<AlertRule>): Promise<AlertRule> {
    const response = await axios.put(`${API_BASE}/alerts/${id}`, alert)
    return response.data
  },

  async deleteAlert(id: number): Promise<void> {
    await axios.delete(`${API_BASE}/alerts/${id}`)
  },

  // Basic Monitoring Stats
  async getStats(): Promise<MonitorStats> {
    const [cpuData, memData, diskData, netData] = await Promise.all([
      axios.get('/api/monitor/cpu'),
      axios.get('/api/monitor/memory'),
      axios.get('/api/monitor/disk'),
      axios.get('/api/monitor/network')
    ])

    return {
      cpu: cpuData.data,
      memory: memData.data,
      disk: diskData.data.disks || [],
      network: netData.data.interfaces || [],
      timestamp: Date.now()
    }
  }
}

// WebSocket Monitor Connection
export class MonitorWebSocket {
  private ws: WebSocket | null = null
  private reconnectTimer: number | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private url: string
  private onMessage: (data: any) => void

  constructor(url: string, onMessage: (data: any) => void) {
    this.url = url
    this.onMessage = onMessage
  }

  connect(): void {
    if (this.ws?.readyState === WebSocket.OPEN) return

    try {
      this.ws = new WebSocket(this.url)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          this.onMessage(data)
        } catch (error) {
          console.error('WebSocket parse error:', error)
        }
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      this.ws.onclose = () => {
        console.log('WebSocket disconnected')
        this.scheduleReconnect()
      }
    } catch (error) {
      console.error('WebSocket connection error:', error)
      this.scheduleReconnect()
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnect attempts reached')
      return
    }

    if (this.reconnectTimer) return

    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000)
    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      this.reconnectAttempts++
      this.connect()
    }, delay)
  }

  disconnect(): void {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }
}