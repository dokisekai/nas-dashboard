// Monitor Pinia Store
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { monitorAPI, MonitorWebSocket } from '@/api/monitor'
import type {
  ProcessInfo,
  ServiceInfo,
  TemperatureInfo,
  SystemEvent,
  SystemLog,
  AlertRule,
  MonitorStats
} from '@/types/monitor'

export const useMonitorStore = defineStore('monitor', () => {
  // State
  const processes = ref<ProcessInfo[]>([])
  const services = ref<ServiceInfo[]>([])
  const temperature = ref<TemperatureInfo | null>(null)
  const events = ref<SystemEvent[]>([])
  const logs = ref<SystemLog[]>([])
  const alerts = ref<AlertRule[]>([])
  const currentStats = ref<MonitorStats | null>(null)

  const loading = ref(false)
  const error = ref<string | null>(null)
  const wsConnected = ref(false)

  // WebSocket instance
  let wsConnection: MonitorWebSocket | null = null

  // Computed
  const runningProcesses = computed(() =>
    processes.value.filter(p => p.status === 'running')
  )

  const totalProcessCPU = computed(() =>
    processes.value.reduce((sum, p) => sum + p.cpuPercent, 0)
  )

  const totalProcessMemory = computed(() =>
    processes.value.reduce((sum, p) => sum + p.memoryPercent, 0)
  )

  const runningServices = computed(() =>
    services.value.filter(s => s.status === 'running')
  )

  const failedServices = computed(() =>
    services.value.filter(s => s.status === 'failed')
  )

  const enabledServices = computed(() =>
    services.value.filter(s => s.enabled)
  )

  const activeAlerts = computed(() =>
    alerts.value.filter(a => a.enabled)
  )

  const criticalAlerts = computed(() =>
    alerts.value.filter(a => a.enabled && a.severity === 'critical')
  )

  const recentEvents = computed(() =>
    events.value.slice(0, 10)
  )

  const systemLoad = computed(() => {
    if (!currentStats.value) return { normal: true, message: '' }
    const load1 = currentStats.value.cpu.load1
    const cores = currentStats.value.cpu.cores
    if (load1 > cores * 2) return { normal: false, message: '严重过载' }
    if (load1 > cores) return { normal: false, message: '高负载' }
    return { normal: true, message: '正常' }
  })

  const memoryStatus = computed(() => {
    if (!currentStats.value) return { normal: true, message: '' }
    const percent = currentStats.value.memory.percent
    if (percent > 90) return { normal: false, message: '内存不足' }
    if (percent > 70) return { normal: false, message: '内存紧张' }
    return { normal: true, message: '正常' }
  })

  const diskStatus = computed(() => {
    if (!currentStats.value?.disk) return { normal: true, message: '' }
    const disk = currentStats.value.disk.find(d => d.usedPercent > 90)
    if (disk) return { normal: false, message: `${disk.device} 空间不足` }
    return { normal: true, message: '正常' }
  })

  // Actions
  async function fetchProcesses() {
    loading.value = true
    error.value = null
    try {
      processes.value = await monitorAPI.getProcesses()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch processes'
      console.error('Error fetching processes:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchServices() {
    loading.value = true
    error.value = null
    try {
      services.value = await monitorAPI.getServices()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch services'
      console.error('Error fetching services:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchTemperature() {
    loading.value = true
    error.value = null
    try {
      temperature.value = await monitorAPI.getTemperature()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch temperature'
      console.error('Error fetching temperature:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchEvents(params?: { page?: number; limit?: number; type?: string; source?: string }) {
    loading.value = true
    error.value = null
    try {
      const response = await monitorAPI.getEvents(params)
      events.value = response.events
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch events'
      console.error('Error fetching events:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchLogs(params?: { page?: number; limit?: number; level?: string; component?: string }) {
    loading.value = true
    error.value = null
    try {
      const response = await monitorAPI.getLogs(params)
      logs.value = response.logs
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch logs'
      console.error('Error fetching logs:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchAlerts() {
    loading.value = true
    error.value = null
    try {
      alerts.value = await monitorAPI.getAlerts()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch alerts'
      console.error('Error fetching alerts:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchStats() {
    loading.value = true
    error.value = null
    try {
      currentStats.value = await monitorAPI.getStats()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stats'
      console.error('Error fetching stats:', err)
    } finally {
      loading.value = false
    }
  }

  async function killProcess(pid: number, signal?: string) {
    loading.value = true
    error.value = null
    try {
      await monitorAPI.killProcess(pid, signal)
      await fetchProcesses()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to kill process'
      console.error('Error killing process:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function startService(name: string) {
    loading.value = true
    error.value = null
    try {
      await monitorAPI.startService(name)
      await fetchServices()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to start service'
      console.error('Error starting service:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function stopService(name: string) {
    loading.value = true
    error.value = null
    try {
      await monitorAPI.stopService(name)
      await fetchServices()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to stop service'
      console.error('Error stopping service:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function restartService(name: string) {
    loading.value = true
    error.value = null
    try {
      await monitorAPI.restartService(name)
      await fetchServices()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to restart service'
      console.error('Error restarting service:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createAlert(alert: Omit<AlertRule, 'id' | 'lastTriggered' | 'createdAt' | 'updatedAt'>) {
    loading.value = true
    error.value = null
    try {
      const newAlert = await monitorAPI.createAlert(alert)
      alerts.value.push(newAlert)
      return newAlert
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create alert'
      console.error('Error creating alert:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateAlert(id: number, data: Partial<AlertRule>) {
    loading.value = true
    error.value = null
    try {
      const updatedAlert = await monitorAPI.updateAlert(id, data)
      const index = alerts.value.findIndex(a => a.id === id)
      if (index !== -1) {
        alerts.value[index] = updatedAlert
      }
      return updatedAlert
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update alert'
      console.error('Error updating alert:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deleteAlert(id: number) {
    loading.value = true
    error.value = null
    try {
      await monitorAPI.deleteAlert(id)
      alerts.value = alerts.value.filter(a => a.id !== id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete alert'
      console.error('Error deleting alert:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function clearLogs(params?: { level?: string; component?: string }) {
    loading.value = true
    error.value = null
    try {
      await monitorAPI.clearLogs(params)
      await fetchLogs()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to clear logs'
      console.error('Error clearing logs:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = window.location.host
    const wsUrl = `${protocol}//${host}/ws/monitor`

    wsConnection = new MonitorWebSocket(wsUrl, (data) => {
      // Handle WebSocket messages
      switch (data.type) {
        case 'stats':
          currentStats.value = data.payload
          break
        case 'process':
          // Update process list if needed
          break
        case 'service':
          // Update service status
          break
        case 'alert':
          // Handle alert trigger
          break
        case 'event':
          // Add new event
          events.value.unshift(data.payload)
          break
      }
    })

    wsConnection.connect()
    wsConnected.value = true
  }

  function disconnectWebSocket() {
    if (wsConnection) {
      wsConnection.disconnect()
      wsConnection = null
    }
    wsConnected.value = false
  }

  function clearError() {
    error.value = null
  }

  // Initialize
  function init() {
    fetchStats()
    fetchServices()
    fetchAlerts()
    connectWebSocket()
  }

  function destroy() {
    disconnectWebSocket()
  }

  return {
    // State
    processes,
    services,
    temperature,
    events,
    logs,
    alerts,
    currentStats,
    loading,
    error,
    wsConnected,

    // Computed
    runningProcesses,
    totalProcessCPU,
    totalProcessMemory,
    runningServices,
    failedServices,
    enabledServices,
    activeAlerts,
    criticalAlerts,
    recentEvents,
    systemLoad,
    memoryStatus,
    diskStatus,

    // Actions
    fetchProcesses,
    fetchServices,
    fetchTemperature,
    fetchEvents,
    fetchLogs,
    fetchAlerts,
    fetchStats,
    killProcess,
    startService,
    stopService,
    restartService,
    createAlert,
    updateAlert,
    deleteAlert,
    clearLogs,
    connectWebSocket,
    disconnectWebSocket,
    clearError,
    init,
    destroy
  }
})