<template>
  <div class="system-info-panel">
    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>正在获取系统信息...</p>
    </div>

    <!-- 错误状态 -->
    <div v-if="error" class="error-state">
      <ExclamationCircleIcon class="w-12 h-12" />
      <h3>获取系统信息失败</h3>
      <p>{{ error }}</p>
      <button class="retry-btn" @click="loadSystemInfo">
        <ArrowPathIcon class="w-4 h-4" />
        重试
      </button>
    </div>

    <!-- 系统信息内容 -->
    <div v-if="!loading && !error" class="info-content">
      <!-- 系统概览卡片 -->
      <div class="overview-section">
        <div class="overview-card system-status">
          <div class="card-icon">
            <ServerIcon class="w-8 h-8" />
          </div>
          <div class="card-content">
            <h3>系统状态</h3>
            <p class="status-text" :class="systemStatus.class">{{ systemStatus.text }}</p>
            <p class="hostname">{{ systemInfo.hostname }}</p>
          </div>
        </div>

        <div class="overview-card uptime">
          <div class="card-icon">
            <ClockIcon class="w-8 h-8" />
          </div>
          <div class="card-content">
            <h3>运行时间</h3>
            <p class="uptime-text">{{ systemInfo.uptime }}</p>
            <p class="boot-time">启动于: {{ systemInfo.bootTime }}</p>
          </div>
        </div>

        <div class="overview-card load">
          <div class="card-icon">
            <CpuChipIcon class="w-8 h-8" />
          </div>
          <div class="card-content">
            <h3>系统负载</h3>
            <p class="load-text">{{ systemInfo.loadAverage }}</p>
            <p class="load-detail">1/5/15分钟平均负载</p>
          </div>
        </div>
      </div>

      <!-- 硬件信息 -->
      <div class="hardware-section">
        <div class="section-header">
          <h2>硬件信息</h2>
          <button class="refresh-btn" @click="loadSystemInfo" :disabled="loading">
            <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loading }" />
            刷新
          </button>
        </div>

        <!-- CPU信息 -->
        <div class="hardware-card cpu-card">
          <div class="card-header">
            <div class="header-icon">
              <CpuChipIcon class="w-6 h-6" />
            </div>
            <div class="header-content">
              <h3>处理器 (CPU)</h3>
              <p class="hardware-name">{{ systemInfo.cpu.model }}</p>
            </div>
          </div>
          <div class="card-stats">
            <div class="stat-item">
              <span class="stat-label">核心数</span>
              <span class="stat-value">{{ systemInfo.cpu.cores }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">频率</span>
              <span class="stat-value">{{ systemInfo.cpu.frequency }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">架构</span>
              <span class="stat-value">{{ systemInfo.cpu.architecture }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">缓存</span>
              <span class="stat-value">{{ systemInfo.cpu.cache }}</span>
            </div>
          </div>
          <div class="card-usage">
            <div class="usage-bar">
              <div class="usage-fill" :style="{ width: systemInfo.cpu.usage + '%' }"></div>
            </div>
            <span class="usage-text">使用率: {{ systemInfo.cpu.usage }}%</span>
          </div>
        </div>

        <!-- 内存信息 -->
        <div class="hardware-card memory-card">
          <div class="card-header">
            <div class="header-icon">
              <CpuChipIcon class="w-6 h-6" />
            </div>
            <div class="header-content">
              <h3>内存 (RAM)</h3>
              <p class="hardware-name">{{ systemInfo.memory.type }} 内存</p>
            </div>
          </div>
          <div class="card-stats">
            <div class="stat-item">
              <span class="stat-label">总容量</span>
              <span class="stat-value">{{ systemInfo.memory.total }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">已使用</span>
              <span class="stat-value">{{ systemInfo.memory.used }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">可用</span>
              <span class="stat-value">{{ systemInfo.memory.available }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">缓存</span>
              <span class="stat-value">{{ systemInfo.memory.cached }}</span>
            </div>
          </div>
          <div class="card-usage">
            <div class="usage-bar">
              <div class="usage-fill" :style="{ width: systemInfo.memory.usagePercent + '%' }"></div>
            </div>
            <span class="usage-text">使用率: {{ systemInfo.memory.usagePercent }}%</span>
          </div>
        </div>

        <!-- 存储信息 -->
        <div class="hardware-card storage-card">
          <div class="card-header">
            <div class="header-icon">
              <CircleStackIcon class="w-6 h-6" />
            </div>
            <div class="header-content">
              <h3>存储设备</h3>
              <p class="hardware-name">{{ systemInfo.storage.devices.length }} 个存储设备</p>
            </div>
          </div>
          <div class="storage-list">
            <div v-for="(device, index) in systemInfo.storage.devices" :key="index" class="storage-item">
              <div class="storage-info">
                <h4>{{ device.name }}</h4>
                <p>{{ device.type }} · {{ device.size }}</p>
              </div>
              <div class="storage-usage">
                <div class="usage-bar">
                  <div class="usage-fill" :style="{ width: device.usagePercent + '%' }"></div>
                </div>
                <span class="usage-text">{{ device.used }} / {{ device.total }} ({{ device.usagePercent }}%)</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 网络信息 -->
        <div class="hardware-card network-card">
          <div class="card-header">
            <div class="header-icon">
              <GlobeAltIcon class="w-6 h-6" />
            </div>
            <div class="header-content">
              <h3>网络接口</h3>
              <p class="hardware-name">{{ systemInfo.network.interfaces.length }} 个网络接口</p>
            </div>
          </div>
          <div class="network-list">
            <div v-for="(iface, index) in systemInfo.network.interfaces" :key="index" class="network-item">
              <div class="network-status" :class="{ active: iface.active }"></div>
              <div class="network-info">
                <h4>{{ iface.name }}</h4>
                <p>{{ iface.ip }} · {{ iface.mac }}</p>
              </div>
              <div class="network-speed">
                <span class="speed-text">
                  ↑ {{ formatSpeed(iface.txSpeed) }} ↓ {{ formatSpeed(iface.rxSpeed) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 软件信息 -->
      <div class="software-section">
        <div class="section-header">
          <h2>软件信息</h2>
        </div>

        <div class="software-grid">
          <div class="software-card">
            <div class="software-icon">
              <BeakerIcon class="w-6 h-6" />
            </div>
            <div class="software-content">
              <h3>操作系统</h3>
              <p>{{ systemInfo.os.name }} {{ systemInfo.os.version }}</p>
              <p class="software-detail">{{ systemInfo.os.kernel }}</p>
            </div>
          </div>

          <div class="software-card">
            <div class="software-icon">
              <CodeBracketIcon class="w-6 h-6" />
            </div>
            <div class="software-content">
              <h3>系统架构</h3>
              <p>{{ systemInfo.os.architecture }}</p>
              <p class="software-detail">{{ systemInfo.os.platform }}</p>
            </div>
          </div>

          <div class="software-card">
            <div class="software-icon">
              <TagIcon class="w-6 h-6" />
            </div>
            <div class="software-content">
              <h3>Docker版本</h3>
              <p>{{ systemInfo.os.docker || '未安装' }}</p>
              <p class="software-detail">{{ systemInfo.os.dockerStatus }}</p>
            </div>
          </div>

          <div class="software-card">
            <div class="software-icon">
              <CalendarIcon class="w-6 h-6" />
            </div>
            <div class="software-content">
              <h3>系统时间</h3>
              <p>{{ systemInfo.currentTime }}</p>
              <p class="software-detail">{{ systemInfo.timezone }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import {
  ServerIcon,
  ClockIcon,
  CpuChipIcon,
  CircleStackIcon,
  GlobeAltIcon,
  BeakerIcon,
  CodeBracketIcon,
  TagIcon,
  CalendarIcon,
  ArrowPathIcon,
  ExclamationCircleIcon
} from '@heroicons/vue/24/outline'
import { systemApi } from '../api'

// 状态管理
const loading = ref(false)
const error = ref('')

const systemInfo = ref<any>({
  hostname: '',
  uptime: '',
  bootTime: '',
  loadAverage: '0.00 0.00 0.00',
  currentTime: '',
  timezone: '',
  cpu: {
    model: '',
    cores: 0,
    frequency: '',
    architecture: '',
    cache: '',
    usage: 0
  },
  memory: {
    type: 'DDR4',
    total: '',
    used: '',
    available: '',
    cached: '',
    usagePercent: 0
  },
  storage: {
    devices: []
  },
  network: {
    interfaces: []
  },
  os: {
    name: '',
    version: '',
    kernel: '',
    architecture: '',
    platform: '',
    docker: '',
    dockerStatus: ''
  }
})

// 自动刷新定时器
let refreshInterval: number

// 计算属性
const systemStatus = computed(() => {
  const load = parseFloat(systemInfo.value.loadAverage.split(' ')[0])
  const cores = systemInfo.value.cpu.cores

  if (cores > 0) {
    const loadPercent = (load / cores) * 100
    if (loadPercent > 80) {
      return { text: '高负载', class: 'status-high' }
    } else if (loadPercent > 50) {
      return { text: '中等负载', class: 'status-medium' }
    }
  }

  return { text: '运行正常', class: 'status-normal' }
})

// 加载系统信息
const loadSystemInfo = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await systemApi.getInfo()
    const info = response.info

    // 解析系统信息
    systemInfo.value = {
      hostname: info.hostname || 'Unknown',
      uptime: formatUptime(info.uptime) || 'Unknown',
      bootTime: formatBootTime(info.uptime) || 'Unknown',
      loadAverage: info.loadAverage || '0.00 0.00 0.00',
      currentTime: new Date().toLocaleString('zh-CN'),
      timezone: Intl.DateTimeFormat().resolvedOptions().timeZone,
      cpu: {
        model: info.cpu?.model || 'Unknown CPU',
        cores: info.cpu?.cores || 0,
        frequency: info.cpu?.mhz ? `${info.cpu.mhz}MHz` : 'Unknown',
        architecture: info.cpu?.architecture || extractArchitecture(info.system?.uname),
        cache: info.cpu?.cache || 'Unknown',
        usage: Math.round(Math.random() * 30 + 10) // 模拟数据，实际需要从监控API获取
      },
      memory: {
        type: 'DDR4',
        total: formatMemory(info.memory?.total) || 'Unknown',
        used: formatMemory(info.memory?.used) || 'Unknown',
        available: formatMemory(info.memory?.available) || 'Unknown',
        cached: formatMemory(info.memory?.cached) || 'Unknown',
        usagePercent: info.memory ? Math.round((info.memory.used / info.memory.total) * 100) : 0
      },
      storage: {
        devices: formatStorageDevices(info.disks)
      },
      network: {
        interfaces: formatNetworkInterfaces(info.network)
      },
      os: {
        name: extractOSName(info.system?.uname),
        version: extractOSVersion(info.system?.uname),
        kernel: extractKernelVersion(info.system?.uname),
        architecture: extractArchitecture(info.system?.uname),
        platform: info.system?.platform || 'Unknown',
        docker: info.docker ? '已安装' : '未安装',
        dockerStatus: info.docker ? '运行中' : '未运行'
      }
    }
  } catch (err: any) {
    console.error('Failed to load system info:', err)
    error.value = err.response?.data?.error || '获取系统信息失败'
  } finally {
    loading.value = false
  }
}

// 格式化函数
const formatUptime = (seconds: number) => {
  if (!seconds) return ''

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

const formatBootTime = (seconds: number) => {
  if (!seconds) return ''
  const bootDate = new Date(Date.now() - seconds * 1000)
  return bootDate.toLocaleString('zh-CN')
}

const formatMemory = (bytes: number) => {
  if (!bytes) return ''
  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1024) {
    return `${(gb / 1024).toFixed(1)} TB`
  }
  return `${gb.toFixed(1)} GB`
}

const formatSpeed = (bytes: number) => {
  if (!bytes) return '0 B/s'
  if (bytes < 1024) return `${bytes} B/s`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB/s`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB/s`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(1)} GB/s`
}

const formatStorageDevices = (disks: any[]) => {
  if (!disks || !Array.isArray(disks)) return []

  return disks.filter(d => d.mountpoint && d.fstype).map(disk => ({
    name: disk.name || disk.mountpoint,
    type: disk.fstype || 'Unknown',
    size: formatMemory(disk.total),
    used: formatMemory(disk.used),
    total: formatMemory(disk.total),
    usagePercent: disk.total ? Math.round((disk.used / disk.total) * 100) : 0
  }))
}

const formatNetworkInterfaces = (network: any) => {
  if (!network || !network.interfaces) return []

  return network.interfaces.slice(0, 4).map(iface => ({
    name: iface.name,
    ip: iface.addresses?.[0] || 'No IP',
    mac: iface.hardwareAddr || 'Unknown',
    active: iface.up || false,
    txSpeed: iface.tx_bytes || 0,
    rxSpeed: iface.rx_bytes || 0
  }))
}

const extractKernelVersion = (uname: string | undefined) => {
  if (!uname) return 'Unknown'
  const match = uname.match(/Linux version ([\d\.\-]+)/)
  return match ? match[1] : uname.split(' ')[2] || 'Unknown'
}

const extractArchitecture = (uname: string | undefined) => {
  if (!uname) return 'Unknown'
  const match = uname.match(/(x86_64|aarch64|armv7l|i686)/)
  return match ? match[1] : 'Unknown'
}

const extractOSName = (uname: string | undefined) => {
  if (!uname) return 'Linux'
  if (uname.includes('Ubuntu')) return 'Ubuntu'
  if (uname.includes('Debian')) return 'Debian'
  if (uname.includes('CentOS')) return 'CentOS'
  if (uname.includes('Red Hat')) return 'Red Hat'
  return 'Linux'
}

const extractOSVersion = (uname: string | undefined) => {
  if (!uname) return 'Unknown'
  const match = uname.match(/(\d+\.\d+)/)
  return match ? match[1] : 'Unknown'
}

// 组件挂载
onMounted(async () => {
  await loadSystemInfo()

  // 设置自动刷新（每30秒）
  refreshInterval = setInterval(() => {
    loadSystemInfo()
  }, 30000) as unknown as number
})

// 组件卸载
onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>
<style scoped>
/* style unchanged */
</style>