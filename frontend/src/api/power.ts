// Power API Client
import axios from 'axios'
import type {
  PowerData,
  PowerHistory,
  PowerStatistics,
  PowerOverview,
  PowerAlertRule,
  PowerTrend
} from '@/types/power'

const API_BASE = '/api/power'

export const powerAPI = {
  // 获取当前功耗数据
  async getCurrent(): Promise<PowerData> {
    const response = await axios.get(`${API_BASE}/current`)
    return response.data
  },

  // 获取历史功耗数据
  async getHistory(days: number = 7): Promise<PowerHistory> {
    const response = await axios.get(`${API_BASE}/history`, {
      params: { days }
    })
    return response.data
  },

  // 获取功耗统计信息
  async getStatistics(days: number = 7): Promise<PowerStatistics> {
    const response = await axios.get(`${API_BASE}/statistics`, {
      params: { days }
    })
    return response.data
  },

  // 获取功耗概览
  async getOverview(): Promise<PowerOverview> {
    const response = await axios.get(`${API_BASE}/overview`)
    return response.data
  },

  // 获取功耗趋势分析
  async getTrend(period: 'hourly' | 'daily' | 'weekly' = 'daily'): Promise<PowerTrend> {
    const response = await axios.get(`${API_BASE}/trend`, {
      params: { period }
    })
    return response.data
  },

  // 获取告警规则列表
  async getAlertRules(): Promise<PowerAlertRule[]> {
    const response = await axios.get(`${API_BASE}/alerts/rules`)
    return response.data.rules || []
  },

  // 创建告警规则
  async createAlertRule(rule: Omit<PowerAlertRule, 'id' | 'createdAt' | 'updatedAt'>): Promise<PowerAlertRule> {
    const response = await axios.post(`${API_BASE}/alerts/rules`, rule)
    return response.data
  },

  // 更新告警规则
  async updateAlertRule(id: number, rule: Partial<PowerAlertRule>): Promise<PowerAlertRule> {
    const response = await axios.put(`${API_BASE}/alerts/rules/${id}`, rule)
    return response.data
  },

  // 删除告警规则
  async deleteAlertRule(id: number): Promise<void> {
    await axios.delete(`${API_BASE}/alerts/rules/${id}`)
  },

  // 测试告警规则
  async testAlertRule(id: number): Promise<{ triggered: boolean; message: string }> {
    const response = await axios.post(`${API_BASE}/alerts/rules/${id}/test`)
    return response.data
  },

  // 监听实时功耗更新
  onPowerUpdate(callback: (data: PowerData) => void): () => void {
    const ws = new WebSocket(`${location.protocol === 'https:' ? 'wss:' : 'ws:'}//${location.host}/ws/power`)

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      callback(data)
    }

    ws.onerror = (error) => {
      console.error('Power WebSocket error:', error)
    }

    // 返回清理函数
    return () => {
      ws.close()
    }
  }
}
