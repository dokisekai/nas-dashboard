// Quota Management API Client
import axios from 'axios'
import type {
  UserQuota,
  GroupQuota,
  QuotaReport,
  QuotaAlert
} from '@/types/quota'

const API_BASE = '/api/storage/quota'

export const quotaAPI = {
  // User Quota Management
  async getUserQuota(username: string): Promise<UserQuota[]> {
    const response = await axios.get(`/api/users/${username}/quota`)
    return response.data.quotas || []
  },

  async getAllQuotas(params?: { type?: string }): Promise<UserQuota[] | GroupQuota[]> {
    const response = await axios.get(`${API_BASE}/users`, { params })
    return response.data.quotas || []
  },

  async getAllGroupQuotas(): Promise<GroupQuota[]> {
    const response = await axios.get(`${API_BASE}/groups`)
    return response.data.quotas || []
  },

  async setUserQuota(username: string, quota: {
    path: string
    softLimit: number
    hardLimit: number
    gracePeriod?: number
  }): Promise<void> {
    await axios.put(`/api/users/${username}/quota`, quota)
  },

  async setGroupQuota(groupId: number, quota: {
    path: string
    softLimit: number
    hardLimit: number
    gracePeriod?: number
  }): Promise<void> {
    await axios.post(`${API_BASE}/set`, {
      type: 'group',
      groupId,
      ...quota
    })
  },

  async deleteUserQuota(username: string, path: string): Promise<void> {
    await axios.delete(`/api/users/${username}/quota`, { params: { path } })
  },

  async deleteGroupQuota(groupId: number, path: string): Promise<void> {
    await axios.delete(`${API_BASE}/groups/${groupId}`, { params: { path } })
  },

  // Quota Reports
  async getQuotaReport(params?: {
    type?: 'user' | 'group' | 'all'
    path?: string
  }): Promise<QuotaReport[]> {
    const response = await axios.get(`${API_BASE}/report`, { params })
    return response.data.report || []
  },

  // Quota Alerts
  async getQuotaAlerts(params?: {
    type?: 'user' | 'group'
    resolved?: boolean
  }): Promise<QuotaAlert[]> {
    const response = await axios.get('/api/quota/alerts', { params })
    return response.data.alerts || []
  },

  async createQuotaAlert(alert: {
    userId?: number
    groupId?: number
    path: string
    alertType: string
    severity: string
    message: string
  }): Promise<QuotaAlert> {
    const response = await axios.post('/api/quota/alerts', alert)
    return response.data
  },

  async resolveQuotaAlert(alertId: number): Promise<void> {
    await axios.put(`/api/quota/alerts/${alertId}/resolve`)
  },

  // Quota Usage Statistics
  async getQuotaUsage(username?: string, path?: string): Promise<{
    usedSpace: number
    fileCount: number
    lastUpdated: string
  }> {
    const params: any = {}
    if (username) params.username = username
    if (path) params.path = path

    const response = await axios.get(`${API_BASE}/usage`, { params })
    return response.data
  },

  // Quota Templates
  async getQuotaTemplates(): Promise<any[]> {
    const response = await axios.get(`${API_BASE}/templates`)
    return response.data.templates || []
  },

  async createQuotaTemplate(template: {
    name: string
    description: string
    softLimit: number
    hardLimit: number
    gracePeriod: number
  }): Promise<void> {
    await axios.post(`${API_BASE}/templates`, template)
  },

  async applyQuotaTemplate(templateId: number, targets: {
    users?: number[]
    groups?: number[]
    paths?: string[]
  }): Promise<void> {
    await axios.post(`${API_BASE}/templates/${templateId}/apply`, targets)
  }
}