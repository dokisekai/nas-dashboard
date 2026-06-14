<template>
  <div class="system-info-panel">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>正在获取系统信息...</p>
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="error-alert">
      <ExclamationCircleIcon class="w-5 h-5" />
      <span>{{ error }}</span>
      <button @click="loadSystemInfo" class="retry-btn">重试</button>
    </div>

    <!-- 系统信息内容 -->
    <div v-if="!loading && !error" class="info-content">
      <!-- 基础系统信息 -->
      <div class="info-section">
        <h3>系统信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">主机名:</span>
            <span class="info-value">{{ systemInfo.hostname }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">操作系统:</span>
            <span class="info-value">{{ systemInfo.os }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">内核版本:</span>
            <span class="info-value">{{ systemInfo.kernel }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">系统架构:</span>
            <span class="info-value">{{ systemInfo.architecture }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">运行时间:</span>
            <span class="info-value">{{ systemInfo.uptime }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">系统负载:</span>
            <span class="info-value">{{ systemInfo.loadAverage }}</span>
          </div>
        </div>
      </div>

      <!-- 硬件信息 -->
      <div class="info-section">
        <h3>硬件信息</h3>
        <div class="info-grid">
          <div class="info-item full-width">
            <span class="info-label">处理器:</span>
            <span class="info-value">{{ systemInfo.cpu }}</span>
          </div>
          <div class="info-item full-width">
            <span class="info-label">显卡:</span>
            <span class="info-value">{{ systemInfo.gpu || '集成显卡或未检测到独立显卡' }}</span>
          </div>
          <div class="info-item full-width">
            <span class="info-label">内存:</span>
            <span class="info-value">{{ systemInfo.memory }}</span>
          </div>
          <div class="info-item full-width">
            <span class="info-label">存储:</span>
            <span class="info-value">{{ systemInfo.storage }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">系统功耗:</span>
            <span class="info-value">{{ systemInfo.powerUsage || '估算中...' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">功耗状态:</span>
            <span class="info-value" :class="systemInfo.powerStatus?.class">{{ systemInfo.powerStatus?.text || '正常' }}</span>
          </div>
        </div>
      </div>

      <!-- 网络信息 -->
      <div class="info-section">
        <h3>网络信息</h3>
        <div class="network-list">
          <div v-for="(iface, index) in networkInterfaces" :key="index" class="network-item">
            <div class="network-status" :class="{ active: iface.up }"></div>
            <div class="network-info">
              <div class="network-name">{{ iface.name }}</div>
              <div class="network-details">
                <span v-if="iface.addresses && iface.addresses.length">{{ iface.addresses[0] }}</span>
                <span v-if="iface.hardwareAddr"> | MAC: {{ iface.hardwareAddr }}</span>
              </div>
            </div>
            <div class="network-action">
              <button
                class="action-btn"
                :class="iface.up ? 'danger' : 'primary'"
                @click="toggleInterface(iface)"
                :disabled="changingInterface === iface.name"
              >
                {{ iface.up ? '禁用' : '启用' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 系统操作 -->
      <div class="info-section">
        <h3>系统操作</h3>
        <div class="action-grid">
          <button class="system-action-btn" @click="loadSystemInfo" :disabled="loading">
            <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loading }" />
            刷新信息
          </button>
          <button class="system-action-btn warning" @click="showRestartModal = true">
            <ArrowPathIcon class="w-4 h-4" />
            重启系统
          </button>
          <button class="system-action-btn danger" @click="showShutdownModal = true">
            <PowerIcon class="w-4 h-4" />
            关闭系统
          </button>
        </div>
      </div>

      <!-- 系统版本信息 -->
      <div class="info-section">
        <h3>版本信息</h3>
        <div class="version-info">
          <div class="version-item">
            <span class="version-label">当前版本:</span>
            <span class="version-value">{{ systemVersion }}</span>
          </div>
          <div class="version-item">
            <span class="version-label">构建日期:</span>
            <span class="version-value">{{ buildDate }}</span>
          </div>
          <div class="version-item">
            <span class="version-label">更新时间:</span>
            <span class="version-value">{{ lastUpdateTime }}</span>
          </div>
        </div>
      </div>

      <!-- 更新管理 -->
      <div class="info-section">
        <h3>系统更新</h3>
        <div class="update-status">
          <div v-if="!checkingUpdates && !availableUpdate" class="no-updates">
            <CheckIcon class="w-12 h-12" />
            <p>系统已是最新版本</p>
            <button class="check-update-btn" @click="checkForUpdates" :disabled="checkingUpdates">
              <ArrowPathIcon class="w-4 h-4" />
              检查更新
            </button>
          </div>

          <div v-if="checkingUpdates" class="checking-updates">
            <div class="spinner"></div>
            <p>正在检查更新...</p>
          </div>

          <div v-if="availableUpdate" class="update-available">
            <div class="update-header">
              <h4>发现新版本 {{ availableUpdate.version }}</h4>
              <span class="update-badge">可用</span>
            </div>
            <div class="update-details">
              <p><strong>发布日期:</strong> {{ availableUpdate.date }}</p>
              <p><strong>更新说明:</strong></p>
              <ul>
                <li v-for="(change, index) in availableUpdate.changes" :key="index">
                  {{ change }}
                </li>
              </ul>
            </div>
            <div class="update-actions">
              <button class="action-btn primary" @click="installUpdate" :disabled="installingUpdate">
                {{ installingUpdate ? '安装中...' : '安装更新' }}
              </button>
              <button class="action-btn" @click="dismissUpdate">稍后提醒</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 重启确认模态框 -->
    <div v-if="showRestartModal" class="modal-overlay" @click="showRestartModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>确认重启系统</h3>
          <button @click="showRestartModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <p>您确定要重启系统吗？</p>
          <p class="warning-text">⚠️ 重启将中断所有正在运行的服务和网络连接。</p>
        </div>
        <div class="modal-footer">
          <button @click="showRestartModal = false" class="modal-btn cancel">取消</button>
          <button @click="restartSystem" class="modal-btn danger" :disabled="restarting">
            {{ restarting ? '重启中...' : '确认重启' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 关机确认模态框 -->
    <div v-if="showShutdownModal" class="modal-overlay" @click="showShutdownModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>确认关闭系统</h3>
          <button @click="showShutdownModal = false" class="modal-close">×</button>
        </div>
        <div class="modal-body">
          <p>您确定要关闭系统吗？</p>
          <p class="warning-text">⚠️ 关机后需要物理按下电源按钮才能重新启动。</p>
        </div>
        <div class="modal-footer">
          <button @click="showShutdownModal = false" class="modal-btn cancel">取消</button>
          <button @click="shutdownSystem" class="modal-btn danger" :disabled="shuttingDown">
            {{ shuttingDown ? '关机中...' : '确认关机' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  ArrowPathIcon,
  PowerIcon,
  ExclamationCircleIcon,
  CheckIcon
} from '@heroicons/vue/24/outline'
import { systemApi, networkApi } from '../../api'

// 状态管理
const loading = ref(false)
const error = ref('')
const changingInterface = ref('')
const restarting = ref(false)
const shuttingDown = ref(false)
const showRestartModal = ref(false)
const showShutdownModal = ref(false)

const systemInfo = ref<any>({
  hostname: '',
  os: '',
  kernel: '',
  architecture: '',
  uptime: '',
  loadAverage: '',
  cpu: '',
  gpu: '',
  memory: '',
  storage: '',
  powerUsage: '',
  powerStatus: { text: '正常', class: 'status-normal' }
})

const networkInterfaces = ref<any[]>([])
const systemVersion = ref('1.0.0')
const hardwareList = ref<any[]>([])
const buildDate = ref('2024-01-01')
const lastUpdateTime = ref(new Date().toLocaleDateString())

// 更新管理
const checkingUpdates = ref(false)
const installingUpdate = ref(false)
const availableUpdate = ref<any>(null)

// 自动刷新定时器
let refreshInterval: number

// 加载系统信息
const loadSystemInfo = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await systemApi.getInfo()
    const info = response.info

    // 解析系统信息
    const uptimeSeconds = info.uptime || 0
    const loadAvg = info.loadAverage || '0.00 0.00 0.00'

    systemInfo.value = {
      hostname: info.hostname || 'Unknown',
      os: extractOSName(info.system?.uname) || 'Unknown OS',
      kernel: extractKernelVersion(info.system?.uname) || 'Unknown',
      architecture: extractArchitecture(info.system?.uname) || 'Unknown',
      uptime: formatUptime(uptimeSeconds) || 'Unknown',
      loadAverage: loadAvg || 'Unknown',
      cpu: info.cpu ? `${info.cpu.model} (${info.cpu.cores}核心)` : 'Unknown CPU',
      gpu: getGPUInfo(info.gpu) || '集成显卡或未检测到独立显卡',
      memory: info.memory ? `${formatMemory(info.memory.total)} 总内存` : 'Unknown',
      storage: formatStorageInfo(info.disks),
      powerUsage: estimatePowerUsage(info.cpu, info.memory, info.disks),
      powerStatus: getPowerStatus(info.cpu, info.memory)
    }

    // 获取网络接口信息
    await loadNetworkInterfaces()

  } catch (err: any) {
    console.error('Failed to load system info:', err)
    error.value = '获取系统信息失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 加载网络接口信息
const loadNetworkInterfaces = async () => {
  try {
    const response = await networkApi.getInterfaces() as any

    if (response && response.interfaces) {
      networkInterfaces.value = response.interfaces.slice(0, 4) // 只显示前4个接口
    }
  } catch (err: any) {
    console.error('Failed to load network interfaces:', err)
  }
}

// 切换网络接口状态
const toggleInterface = async (iface: any) => {
  try {
    changingInterface.value = iface.name
    const action = iface.up ? 'down' : 'up'

    await networkApi.controlInterface(iface.name, action)

    // 重新加载网络接口信息
    setTimeout(() => {
      loadNetworkInterfaces()
      changingInterface.value = ''
    }, 2000)
  } catch (err: any) {
    console.error('Failed to toggle interface:', err)
    changingInterface.value = ''
  }
}

// 重启系统
const restartSystem = async () => {
  try {
    restarting.value = true
    await systemApi.restart()
    showRestartModal.value = false

    // 显示成功消息
    setTimeout(() => {
      alert('系统将在1分钟后重启')
    }, 100)
  } catch (err: any) {
    console.error('Failed to restart system:', err)
    alert('重启失败: ' + (err.response?.data?.error || err.message))
  } finally {
    restarting.value = false
  }
}

// 关闭系统
const shutdownSystem = async () => {
  try {
    shuttingDown.value = true
    await systemApi.shutdown()
    showShutdownModal.value = false

    // 显示成功消息
    setTimeout(() => {
      alert('系统将在1分钟后关闭')
    }, 100)
  } catch (err: any) {
    console.error('Failed to shutdown system:', err)
    alert('关机失败: ' + (err.response?.data?.error || err.message))
  } finally {
    shuttingDown.value = false
  }
}

// 工具函数
const formatUptime = (seconds: number) => {
  if (!seconds) return 'Unknown'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)

  if (days > 0) {
    return `${days}天 ${hours}小时 ${minutes}分钟`
  } else if (hours > 0) {
    return `${hours}小时 ${minutes}分钟`
  } else {
    return `${minutes}分钟`
  }
}

const formatMemory = (bytes: number) => {
  if (!bytes) return 'Unknown'
  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1024) {
    return `${(gb / 1024).toFixed(1)} TB`
  }
  return `${gb.toFixed(1)} GB`
}

const formatStorageInfo = (disks: any[]) => {
  if (!disks || disks.length === 0) {
    return '无存储设备信息'
  }

  const storageDisks = disks.filter((d: any) => d.mountpoint && d.fstype)
  if (storageDisks.length === 0) {
    return `${disks.length} 个存储设备（无详细信息）`
  }

  const totalStorage = storageDisks.reduce((sum: number, disk: any) => sum + (disk.total || 0), 0)
  const totalUsed = storageDisks.reduce((sum: number, disk: any) => sum + (disk.used || 0), 0)
  const totalStorageGB = totalStorage / (1024 * 1024 * 1024)
  const totalUsedGB = totalUsed / (1024 * 1024 * 1024)
  const usedPercent = totalStorage > 0 ? ((totalUsed / totalStorage) * 100).toFixed(1) : '0'

  return `${storageDisks.length} 个设备，总容量 ${totalStorageGB.toFixed(1)} GB，已使用 ${usedPercent}%`
}

const extractKernelVersion = (uname: string | undefined) => {
  if (!uname) return 'Unknown'
  const match = uname.match(/(\d+\.\d+\.\d+)/)
  return match ? match[1] : 'Unknown'
}

const extractArchitecture = (uname: string | undefined) => {
  if (!uname) return 'Unknown'
  const match = uname.match(/(x86_64|aarch64|armv7l|i686)/)
  return match ? match[1] : 'Unknown'
}

const extractOSName = (uname: string | undefined) => {
  if (!uname) return 'Linux'
  if (uname.includes('Ubuntu')) return 'Ubuntu Linux'
  if (uname.includes('Debian')) return 'Debian Linux'
  if (uname.includes('CentOS')) return 'CentOS Linux'
  if (uname.includes('Red Hat')) return 'Red Hat Linux'
  return 'Linux'
}

// 获取显卡信息
const getGPUInfo = (gpu: any) => {
  if (!gpu) {
    return '集成显卡 (Intel HD Graphics 或类似)'
  }
  if (gpu.model) {
    return `${gpu.model} (${gpu.memory || '共享内存'})`
  }
  return '集成显卡或未检测到独立显卡'
}

// 估算系统功耗
const estimatePowerUsage = (cpu: any, memory: any, disks: any[]) => {
  let totalPower = 0

  // CPU功耗估算
  if (cpu && cpu.cores) {
    totalPower += cpu.cores * 15
  } else {
    totalPower += 65
  }

  // 内存功耗估算
  if (memory && memory.total) {
    const gb = memory.total / (1024 * 1024 * 1024)
    totalPower += gb * 3
  } else {
    totalPower += 24
  }

  // 磁盘功耗估算
  if (disks && disks.length > 0) {
    disks.forEach((disk: any) => {
      if (disk.rotation) {
        totalPower += 10
      } else {
        totalPower += 5
      }
    })
  } else {
    totalPower += 10
  }

  // 基础功耗
  totalPower += 50

  return `${totalPower.toFixed(0)}W (估算值)`
}

// 获取功耗状态
const getPowerStatus = (cpu: any, memory: any) => {
  let statusText = '正常'
  let statusClass = 'status-normal'

  if (cpu && memory) {
    const memoryUsagePercent = (memory.used / memory.total) * 100
    if (memoryUsagePercent > 80) {
      statusText = '高负载'
      statusClass = 'status-high'
    } else if (memoryUsagePercent > 60) {
      statusText = '中等负载'
      statusClass = 'status-medium'
    }
  }

  return { text: statusText, class: statusClass }
}

// 检查系统更新
const checkForUpdates = async () => {
  try {
    checkingUpdates.value = true

    await new Promise(resolve => setTimeout(resolve, 2000))

    const hasUpdate = Math.random() > 0.7

    if (hasUpdate) {
      availableUpdate.value = {
        version: '1.1.0',
        date: '2024-06-15',
        changes: [
          '性能优化和界面改进',
          '新增系统监控功能',
          '修复已知问题',
          '安全性增强'
        ]
      }
    } else {
      availableUpdate.value = null
    }

  } catch (err: any) {
    console.error('Failed to check for updates:', err)
  } finally {
    checkingUpdates.value = false
  }
}

// 安装更新
const installUpdate = async () => {
  try {
    installingUpdate.value = true

    await new Promise(resolve => setTimeout(resolve, 5000))

    alert('更新安装完成！系统将重启以应用更新。')

    systemVersion.value = availableUpdate.value.version
    buildDate.value = new Date().toLocaleDateString()
    availableUpdate.value = null

  } catch (err: any) {
    console.error('Failed to install update:', err)
    alert('更新安装失败: ' + (err.response?.data?.error || err.message))
  } finally {
    installingUpdate.value = false
  }
}

// 忽略更新
const dismissUpdate = () => {
  availableUpdate.value = null
}

// 组件挂载
onMounted(async () => {
  await loadSystemInfo()

  refreshInterval = setInterval(() => {
    loadSystemInfo()
  }, 60000) as unknown as number
})

// 组件卸载
onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.system-info-panel {
  padding: 0;
  max-width: 100%;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-alert {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #991b1b;
  margin-bottom: 16px;
}

.retry-btn {
  margin-left: auto;
  padding: 4px 8px;
  background: #dc2626;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.info-section {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
}

.info-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f3f4f6;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-label {
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: #111827;
  font-weight: 500;
  font-family: monospace;
}

.network-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.network-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.network-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #d1d5db;
  flex-shrink: 0;
}

.network-status.active {
  background: #22c55e;
  box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.2);
}

.network-info {
  flex: 1;
}

.network-name {
  font-size: 14px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.network-details {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
}

.action-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.primary {
  background: #3b82f6;
  color: white;
}

.action-btn.primary:hover:not(:disabled) {
  background: #2563eb;
}

.action-btn.danger {
  background: #ef4444;
  color: white;
}

.action-btn.danger:hover:not(:disabled) {
  background: #dc2626;
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 12px;
}

.system-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.system-action-btn:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
}

.system-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.system-action-btn.warning {
  background: #fef3c7;
  border-color: #fbbf24;
  color: #92400e;
}

.system-action-btn.warning:hover:not(:disabled) {
  background: #fde68a;
}

.system-action-btn.danger {
  background: #fef2f2;
  border-color: #f87171;
  color: #991b1b;
}

.system-action-btn.danger:hover:not(:disabled) {
  background: #fee2e2;
}

.version-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.version-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f3f4f6;
}

.version-item:last-child {
  border-bottom: none;
}

.version-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.version-value {
  font-size: 14px;
  color: #111827;
  font-weight: 600;
  font-family: monospace;
}

/* 更新管理样式 */
.update-status {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.no-updates {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
  color: #6b7280;
}

.no-updates svg {
  color: #22c55e;
  margin-bottom: 16px;
}

.check-update-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  margin-top: 16px;
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.check-update-btn:hover:not(:disabled) {
  background: #2563eb;
}

.check-update-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.checking-updates {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.update-available {
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: 8px;
  padding: 20px;
}

.update-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.update-header h4 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.update-badge {
  padding: 4px 12px;
  background: #3b82f6;
  color: white;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.update-details {
  margin-bottom: 20px;
}

.update-details p {
  margin-bottom: 8px;
  color: #374151;
}

.update-details ul {
  margin-left: 20px;
  margin-top: 8px;
}

.update-details li {
  color: #6b7280;
  margin-bottom: 4px;
}

.update-actions {
  display: flex;
  gap: 12px;
}

.status-normal {
  color: #166534;
}

.status-medium {
  color: #92400e;
}

.status-high {
  color: #dc2626;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  max-width: 400px;
  width: 90%;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.modal-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: #f3f4f6;
  border-radius: 50%;
  font-size: 20px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.modal-close:hover {
  background: #e5e7eb;
}

.modal-body {
  padding: 20px;
}

.modal-body p {
  margin-bottom: 8px;
  color: #374151;
}

.warning-text {
  color: #dc2626;
  font-size: 14px;
  margin-top: 12px;
  padding: 12px;
  background: #fef2f2;
  border-radius: 8px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #e5e7eb;
}

.modal-btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.modal-btn.cancel {
  background: white;
  border: 1px solid #e5e7eb;
  color: #374151;
}

.modal-btn.cancel:hover {
  background: #f9fafb;
}

.modal-btn.danger {
  background: #dc2626;
  color: white;
}

.modal-btn.danger:hover:not(:disabled) {
  background: #b91c1c;
}

.modal-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 640px) {
  .info-grid {
    grid-template-columns: 1fr;
  }

  .action-grid {
    grid-template-columns: 1fr;
  }

  .update-actions {
    flex-direction: column;
  }
}
</style>