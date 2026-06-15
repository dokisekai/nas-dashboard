// Quota Management Pinia Store
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { quotaAPI } from '@/api/quota'
import type {
  UserQuota,
  GroupQuota,
  QuotaReport,
  QuotaAlert
} from '@/types/quota'
import { calculateQuotaPercent, getQuotaStatus, formatBytes } from '@/types/quota'

export const useQuotaStore = defineStore('quota', () => {
  // State
  const userQuotas = ref<UserQuota[]>([])
  const groupQuotas = ref<GroupQuota[]>([])
  const quotaReports = ref<QuotaReport[]>([])
  const quotaAlerts = ref<QuotaAlert[]>([])

  const loading = ref(false)
  const error = ref<string | null>(null)

  // Computed
  const totalUserQuotas = computed(() => userQuotas.value.length)
  const totalGroupQuotas = computed(() => groupQuotas.value.length)

  const exceededQuotas = computed(() =>
    quotaReports.value.filter(q => q.status === 'exceeded')
  )

  const warningQuotas = computed(() =>
    quotaReports.value.filter(q => q.status === 'warning')
  )

  const criticalAlerts = computed(() =>
    quotaAlerts.value.filter(a => a.severity === 'critical' && !a.resolved)
  )

  const totalQuotaSpace = computed(() =>
    userQuotas.value.reduce((sum, q) => sum + q.hardLimit, 0)
  )

  const usedQuotaSpace = computed(() =>
    userQuotas.value.reduce((sum, q) => sum + q.usedSpace, 0)
  )

  const quotaUsagePercent = computed(() => {
    if (totalQuotaSpace.value === 0) return 0
    return (usedQuotaSpace.value / totalQuotaSpace.value) * 100
  })

  const quotasByPath = computed(() => {
    const grouped: Record<string, { user: UserQuota[]; group: GroupQuota[] }> = {}

    userQuotas.value.forEach(quota => {
      if (!grouped[quota.path]) {
        grouped[quota.path] = { user: [], group: [] }
      }
      grouped[quota.path].user.push(quota)
    })

    groupQuotas.value.forEach(quota => {
      if (!grouped[quota.path]) {
        grouped[quota.path] = { user: [], group: [] }
      }
      grouped[quota.path].group.push(quota)
    })

    return grouped
  })

  const topQuotaConsumers = computed(() => {
    const allQuotas = [
      ...userQuotas.value.map(q => ({ ...q, type: 'user' as any })),
      ...groupQuotas.value.map(q => ({ ...q, type: 'group' as any }))
    ]

    return allQuotas
      .filter(q => q.hardLimit > 0)
      .sort((a, b) => {
        const aPercent = calculateQuotaPercent(a.usedSpace, a.hardLimit)
        const bPercent = calculateQuotaPercent(b.usedSpace, b.hardLimit)
        return bPercent - aPercent
      })
      .slice(0, 10)
  })

  // Actions
  async function fetchUserQuotas(username?: string) {
    loading.value = true
    error.value = null
    try {
      if (username) {
        const quotas = await quotaAPI.getUserQuota(username)
        // Update specific user quotas in the list
        userQuotas.value = userQuotas.value.filter(q => q.user.username !== username)
        userQuotas.value.push(...quotas)
      } else {
        userQuotas.value = await quotaAPI.getAllQuotas() as UserQuota[]
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch user quotas'
      console.error('Error fetching user quotas:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchGroupQuotas() {
    loading.value = true
    error.value = null
    try {
      groupQuotas.value = await quotaAPI.getAllGroupQuotas()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch group quotas'
      console.error('Error fetching group quotas:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchQuotaReports(params?: { type?: 'user' | 'group' | 'all'; path?: string }) {
    loading.value = true
    error.value = null
    try {
      quotaReports.value = await quotaAPI.getQuotaReport(params)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch quota reports'
      console.error('Error fetching quota reports:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchQuotaAlerts(params?: { type?: 'user' | 'group'; resolved?: boolean }) {
    loading.value = true
    error.value = null
    try {
      quotaAlerts.value = await quotaAPI.getQuotaAlerts(params)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch quota alerts'
      console.error('Error fetching quota alerts:', err)
    } finally {
      loading.value = false
    }
  }

  async function setUserQuota(username: string, quota: {
    path: string
    softLimit: number
    hardLimit: number
    gracePeriod?: number
  }) {
    loading.value = true
    error.value = null
    try {
      await quotaAPI.setUserQuota(username, quota)
      await fetchUserQuotas(username)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to set user quota'
      console.error('Error setting user quota:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function setGroupQuota(groupId: number, quota: {
    path: string
    softLimit: number
    hardLimit: number
    gracePeriod?: number
  }) {
    loading.value = true
    error.value = null
    try {
      await quotaAPI.setGroupQuota(groupId, quota)
      await fetchGroupQuotas()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to set group quota'
      console.error('Error setting group quota:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteUserQuota(username: string, path: string) {
    loading.value = true
    error.value = null
    try {
      await quotaAPI.deleteUserQuota(username, path)
      userQuotas.value = userQuotas.value.filter(q =>
        !(q.user.username === username && q.path === path)
      )
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete user quota'
      console.error('Error deleting user quota:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteGroupQuota(groupId: number, path: string) {
    loading.value = true
    error.value = null
    try {
      await quotaAPI.deleteGroupQuota(groupId, path)
      groupQuotas.value = groupQuotas.value.filter(q =>
        !(q.group.id === groupId && q.path === path)
      )
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete group quota'
      console.error('Error deleting group quota:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function resolveQuotaAlert(alertId: number) {
    loading.value = true
    error.value = null
    try {
      await quotaAPI.resolveQuotaAlert(alertId)
      quotaAlerts.value = quotaAlerts.value.map(alert =>
        alert.id === alertId ? { ...alert, resolved: true } : alert
      )
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to resolve quota alert'
      console.error('Error resolving quota alert:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function generateReport(params?: { type?: 'user' | 'group' | 'all'; path?: string }) {
    return fetchQuotaReports(params)
  }

  function getQuotaStatusInfo(used: number, softLimit: number, hardLimit: number) {
    const percent = calculateQuotaPercent(used, hardLimit)
    const status = getQuotaStatus(used, softLimit, hardLimit)

    return {
      percent,
      status,
      usedFormatted: formatBytes(used),
      softLimitFormatted: formatBytes(softLimit),
      hardLimitFormatted: formatBytes(hardLimit),
      remaining: Math.max(0, hardLimit - used),
      remainingFormatted: formatBytes(Math.max(0, hardLimit - used))
    }
  }

  function clearError() {
    error.value = null
  }

  // Initialize
  function init() {
    fetchUserQuotas()
    fetchGroupQuotas()
    fetchQuotaReports()
    fetchQuotaAlerts({ resolved: false })
  }

  return {
    // State
    userQuotas,
    groupQuotas,
    quotaReports,
    quotaAlerts,
    loading,
    error,

    // Computed
    totalUserQuotas,
    totalGroupQuotas,
    exceededQuotas,
    warningQuotas,
    criticalAlerts,
    totalQuotaSpace,
    usedQuotaSpace,
    quotaUsagePercent,
    quotasByPath,
    topQuotaConsumers,

    // Actions
    fetchUserQuotas,
    fetchGroupQuotas,
    fetchQuotaReports,
    fetchQuotaAlerts,
    setUserQuota,
    setGroupQuota,
    deleteUserQuota,
    deleteGroupQuota,
    resolveQuotaAlert,
    generateReport,
    getQuotaStatusInfo,
    clearError,
    init
  }
})