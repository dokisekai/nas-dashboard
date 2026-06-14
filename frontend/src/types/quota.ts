// Quota Management Type Definitions

export interface UserQuota {
  id: number
  userId: number
  user: User
  path: string
  softLimit: number
  hardLimit: number
  usedSpace: number
  filesUsed: number
  filesSoft: number
  filesHard: number
  gracePeriod: number
  createdAt: string
  updatedAt: string
}

export interface GroupQuota {
  id: number
  groupId: number
  group: Group
  path: string
  softLimit: number
  hardLimit: number
  usedSpace: number
  filesUsed: number
  filesSoft: number
  filesHard: number
  gracePeriod: number
  createdAt: string
  updatedAt: string
}

export interface QuotaReport {
  name: string
  type: 'user' | 'group'
  path: string
  usedSpace: number
  softLimit: number
  hardLimit: number
  usedPercent: number
  status: 'ok' | 'warning' | 'exceeded'
  filesUsed: number
  filesSoft: number
  filesHard: number
  generatedAt: string
}

export interface QuotaAlert {
  id: number
  userId: number | null
  groupId: number | null
  type: 'user' | 'group' | 'user_quota' | 'disk_space'
  path: string
  alertType: 'soft_limit' | 'hard_limit'
  severity: 'warning' | 'critical'
  message: string
  resolved: boolean
  createdAt: string
  updatedAt: string
  // UI 扩展字段
  name?: string
  enabled?: boolean
  threshold?: number
  duration?: number
  condition?: string
  actions?: string
}

export interface User {
  id: number
  username: string
  displayName: string
  email: string
  role: string
}

export interface Group {
  id: number
  name: string
  description: string
}

export interface QuotaSettings {
  path: string
  softLimit: number
  hardLimit: number
  filesSoft: number
  filesHard: number
  gracePeriod: number
  enabled: boolean
}

export const QUOTA_STATUS_COLORS = {
  ok: 'green',
  warning: 'orange',
  exceeded: 'red'
}

export const QUOTA_TYPE_LABELS = {
  user: '用户配额',
  group: '组配额'
}

export const ALERT_TYPE_LABELS = {
  soft_limit: '软限制',
  hard_limit: '硬限制'
}

export const SEVERITY_COLORS = {
  warning: 'orange',
  critical: 'red'
}

// Helper functions for quota calculations
export function calculateQuotaPercent(used: number, limit: number): number {
  if (limit === 0) return 0
  return (used / limit) * 100
}

export function getQuotaStatus(used: number, softLimit: number, hardLimit: number): 'ok' | 'warning' | 'exceeded' {
  if (hardLimit > 0 && used > hardLimit) return 'exceeded'
  if (softLimit > 0 && used > softLimit) return 'warning'
  return 'ok'
}

export function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

export function formatFilesCount(count: number): string {
  return count.toLocaleString()
}