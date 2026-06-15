<template>
  <div class="network-interfaces-component">
    <div class="interfaces-header">
      <h3>网络接口配置</h3>
      <button class="refresh-btn" @click="loadInterfaces" :disabled="loading">
        <ArrowPathIcon class="w-4 h-4" />
        刷新
      </button>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载网络接口...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <ExclamationTriangleIcon class="w-6 h-6" />
      <p>{{ error }}</p>
      <button @click="loadInterfaces">重试</button>
    </div>

    <div v-else class="interfaces-list">
      <div
        v-for="iface in interfaces"
        :key="iface.name"
        class="interface-card"
        :class="{
          active: iface.up,
          virtual: isVirtualInterface(iface.name),
          selected: selectedInterface === iface.name
        }"
      >
        <div class="interface-header" @click="toggleInterface(iface.name)">
          <div class="interface-info">
            <div class="interface-icon" :class="{ active: iface.up }">
              <SignalIcon v-if="iface.type === 'wireless'" class="w-5 h-5" />
              <ServerIcon v-else class="w-5 h-5" />
            </div>
            <div>
              <div class="interface-name">{{ iface.name }}</div>
              <div class="interface-type">
                {{ getInterfaceType(iface.name) }}
                <span v-if="isVirtualInterface(iface.name)" class="virtual-badge">虚拟</span>
              </div>
            </div>
          </div>

          <div class="interface-status">
            <div class="status-indicator" :class="{ active: iface.up }"></div>
            <span>{{ iface.up ? '已连接' : '未连接' }}</span>
            <ChevronDownIcon class="w-4 h-4 chevron" :class="{ rotated: selectedInterface === iface.name }" />
          </div>
        </div>

        <div v-if="selectedInterface === iface.name" class="interface-details">
          <div class="detail-grid">
            <div class="detail-item">
              <span class="detail-label">MAC地址:</span>
              <span class="detail-value">{{ iface.hardwareAddr }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">MTU:</span>
              <span class="detail-value">{{ iface.mtu }}</span>
            </div>
            <div class="detail-item full-width">
              <span class="detail-label">IP地址:</span>
              <div class="ip-list">
                <span v-for="(addr, idx) in parseAddresses(iface.addresses)" :key="idx" class="ip-address">
                  {{ addr }}
                </span>
              </div>
            </div>
            <div class="detail-item">
              <span class="detail-label">上传速度:</span>
              <span class="detail-value">{{ formatSpeed(iface.sentSpeed) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">下载速度:</span>
              <span class="detail-value">{{ formatSpeed(iface.recvSpeed) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">总发送:</span>
              <span class="detail-value">{{ formatBytes(iface.bytesSent) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">总接收:</span>
              <span class="detail-value">{{ formatBytes(iface.bytesRecv) }}</span>
            </div>
          </div>

          <div v-if="!isVirtualInterface(iface.name)" class="interface-actions">
            <button
              @click="toggleInterfaceState(iface.name, iface.up)"
              :class="{ danger: iface.up }"
            >
              {{ iface.up ? '禁用接口' : '启用接口' }}
            </button>
            <button @click="configureInterface(iface)">
              配置接口
            </button>
            <button @click="viewInterfaceDetails(iface)" class="secondary">
              查看详情
            </button>
          </div>

          <div class="interface-flags">
            <span
              v-for="flag in iface.flags"
              :key="flag"
              class="flag-badge"
            >
              {{ flag }}
            </span>
          </div>
        </div>
      </div>

      <div class="network-summary">
        <div class="summary-card">
          <div class="summary-icon">
            <GlobeAltIcon class="w-6 h-6" />
          </div>
          <div class="summary-content">
            <div class="summary-label">网络状态</div>
            <div class="summary-value">{{ networkSummary }}</div>
          </div>
        </div>
        <div class="summary-card">
          <div class="summary-icon">
            <ArrowUpIcon class="w-6 h-6" />
          </div>
          <div class="summary-content">
            <div class="summary-label">总上传速度</div>
            <div class="summary-value">{{ formatSpeed(totalUploadSpeed) }}</div>
          </div>
        </div>
        <div class="summary-card">
          <div class="summary-icon">
            <ArrowDownIcon class="w-6 h-6" />
          </div>
          <div class="summary-content">
            <div class="summary-label">总下载速度</div>
            <div class="summary-value">{{ formatSpeed(totalDownloadSpeed) }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  SignalIcon,
  ServerIcon,
  ArrowPathIcon,
  ExclamationTriangleIcon,
  ChevronDownIcon,
  GlobeAltIcon,
  ArrowUpIcon,
  ArrowDownIcon
} from '@heroicons/vue/24/outline'

interface NetworkInterface {
  name: string
  hardwareAddr: string
  up: boolean
  addresses: string[]
  mtu: number
  flags: string[]
  bytesSent: number
  bytesRecv: number
  packetsSent: number
  packetsRecv: number
  errin: number
  errout: number
  dropin: number
  dropout: number
  sentSpeed: number
  recvSpeed: number
  type?: 'wireless' | 'wired' | 'virtual'
}

const loading = ref(false)
const error = ref('')
const interfaces = ref<any[]>([])
const selectedInterface = ref<string | null>(null)

// 计算属性
const networkSummary = computed(() => {
  const activeCount = interfaces.value.filter(i => i.up && !isVirtualInterface(i.name)).length
  if (activeCount > 0) {
    return `已连接 (${activeCount}个活动接口)`
  } else {
    return '未连接'
  }
})

const totalUploadSpeed = computed(() => {
  return interfaces.value.reduce((sum, iface) => sum + (iface.sentSpeed || 0), 0)
})

const totalDownloadSpeed = computed(() => {
  return interfaces.value.reduce((sum, iface) => sum + (iface.recvSpeed || 0), 0)
})

// 方法
const loadInterfaces = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await fetch('/api/monitor/network', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) {
      throw new Error('获取网络接口信息失败')
    }

    const data = await response.json()

    // 转换接口数据并添加类型信息
    interfaces.value = data.interfaces.map((iface: any) => ({
      ...iface,
      type: getInterfaceType(iface.name)
    }))

  } catch (err: any) {
    error.value = err.message || '加载网络接口失败'
    console.error('Failed to load network interfaces:', err)
  } finally {
    loading.value = false
  }
}

const toggleInterface = (name: string) => {
  if (selectedInterface.value === name) {
    selectedInterface.value = null
  } else {
    selectedInterface.value = name
  }
}

const isVirtualInterface = (name: string): boolean => {
  return name.startsWith('veth') ||
         name.startsWith('virbr') ||
         name.startsWith('docker') ||
         name.startsWith('br-') ||
         name.startsWith('lo')
}

const getInterfaceType = (name: string): string => {
  if (name.startsWith('wlan') || name.startsWith('wlp')) {
    return 'wireless'
  } else if (name.startsWith('eth') || name.startsWith('enp')) {
    return 'wired'
  } else {
    return 'virtual'
  }
}

const parseAddresses = (addresses: string[]): string[] => {
  if (!addresses || !Array.isArray(addresses)) return []

  return addresses.map(addr => {
    try {
      // 解析类似 "{\"addr\":\"192.168.1.100\" 的格式
      const match = addr.match(/\"addr\":\"([^\"]+)\"/)
      return match ? match[1] : addr
    } catch {
      return addr
    }
  }).filter(addr => addr && addr.length > 0)
}

const formatSpeed = (bytesPerSecond: number): string => {
  if (!bytesPerSecond || bytesPerSecond === 0) return '0 B/s'

  const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
  let value = bytesPerSecond
  let unitIndex = 0

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024
    unitIndex++
  }

  return `${value.toFixed(1)} ${units[unitIndex]}`
}

const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'

  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = bytes
  let unitIndex = 0

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024
    unitIndex++
  }

  return `${value.toFixed(1)} ${units[unitIndex]}`
}

const toggleInterfaceState = async (name: string, currentState: boolean) => {
  try {
    const action = currentState ? '禁用' : '启用'
    const response = await fetch(`/api/network/interfaces/${name}/${currentState ? 'down' : 'up'}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`${action}接口失败`)
    }

    ElMessage.success(`接口${name}已${action}`)
    await loadInterfaces()

  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
    console.error('Failed to toggle interface state:', err)
  }
}

const configureInterface = (iface: NetworkInterface) => {
  ElMessage.info('接口配置功能开发中...')
  // TODO: 实现接口配置对话框
}

const viewInterfaceDetails = (iface: NetworkInterface) => {
  console.log('View details for interface:', iface.name)
  // TODO: 实现详情对话框或跳转到详情页面
}

// 生命周期
onMounted(() => {
  loadInterfaces()

  // 每5秒自动刷新一次
  const interval = setInterval(() => {
    loadInterfaces()
  }, 5000)

  // 清理定时器
  return () => clearInterval(interval)
})
</script>

<style scoped lang="scss">
.network-interfaces-component {
  .interfaces-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }

    .refresh-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      border: 1px solid #d1d5db;
      border-radius: 6px;
      background: white;
      cursor: pointer;
      transition: all 0.2s;

      &:hover:not(:disabled) {
        background: #f3f4f6;
        border-color: #9ca3af;
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }
    }
  }

  .loading-state, .error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 40px;
    text-align: center;

    .spinner {
      width: 40px;
      height: 40px;
      border: 4px solid #e5e7eb;
      border-top-color: #3b82f6;
      border-radius: 50%;
      animation: spin 1s linear infinite;
      margin-bottom: 16px;
    }

    button {
      margin-top: 16px;
      padding: 8px 16px;
      background: #3b82f6;
      color: white;
      border: none;
      border-radius: 6px;
      cursor: pointer;

      &:hover {
        background: #2563eb;
      }
    }
  }

  .interfaces-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .interface-card {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    background: white;
    transition: all 0.2s;

    &:hover {
      border-color: #d1d5db;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }

    &.active {
      border-left: 4px solid #10b981;
    }

    &.virtual {
      opacity: 0.7;
    }

    &.selected {
      background: #f9fafb;
    }
  }

  .interface-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    cursor: pointer;
    user-select: none;
  }

  .interface-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .interface-icon {
    width: 40px;
    height: 40px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #f3f4f6;
    color: #6b7280;

    &.active {
      background: #d1fae5;
      color: #10b981;
    }
  }

  .interface-name {
    font-weight: 600;
    color: #1f2937;
  }

  .interface-type {
    font-size: 13px;
    color: #6b7280;
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .virtual-badge {
    font-size: 11px;
    padding: 2px 6px;
    background: #e5e7eb;
    border-radius: 4px;
  }

  .interface-status {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;

    .chevron {
      transition: transform 0.2s;

      &.rotated {
        transform: rotate(180deg);
      }
    }
  }

  .status-indicator {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #d1d5db;

    &.active {
      background: #10b981;
      box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2);
    }
  }

  .interface-details {
    border-top: 1px solid #e5e7eb;
    padding: 16px;
  }

  .detail-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 12px;
    margin-bottom: 16px;
  }

  .detail-item {
    display: flex;
    flex-direction: column;
    gap: 4px;

    &.full-width {
      grid-column: 1 / -1;
    }
  }

  .detail-label {
    font-size: 12px;
    color: #6b7280;
  }

  .detail-value {
    font-size: 14px;
    color: #1f2937;
    font-weight: 500;
  }

  .ip-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .ip-address {
    font-size: 13px;
    padding: 4px 8px;
    background: #f3f4f6;
    border-radius: 4px;
    font-family: monospace;
  }

  .interface-actions {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;

    button {
      padding: 8px 16px;
      border: 1px solid #d1d5db;
      border-radius: 6px;
      background: white;
      cursor: pointer;
      transition: all 0.2s;

      &:hover:not(.danger) {
        background: #f3f4f6;
        border-color: #9ca3af;
      }

      &.danger:hover {
        background: #fef2f2;
        border-color: #ef4444;
        color: #ef4444;
      }

      &.secondary {
        opacity: 0.7;
      }
    }
  }

  .interface-flags {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }

  .flag-badge {
    font-size: 11px;
    padding: 4px 8px;
    background: #dbeafe;
    color: #1e40af;
    border-radius: 4px;
  }

  .network-summary {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 12px;
    margin-top: 16px;
  }

  .summary-card {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px;
    background: #f9fafb;
    border-radius: 8px;
    border: 1px solid #e5e7eb;
  }

  .summary-icon {
    width: 40px;
    height: 40px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: white;
    color: #6b7280;
  }

  .summary-content {
    flex: 1;
  }

  .summary-label {
    font-size: 12px;
    color: #6b7280;
  }

  .summary-value {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
