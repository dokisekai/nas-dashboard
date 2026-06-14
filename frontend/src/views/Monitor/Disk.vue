<template>
  <div class="space-y-6">
    <!-- 连接状态指示器 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-white">磁盘监控</h2>
        <p class="text-gray-500">实时监控磁盘使用情况和 IO 速度</p>
      </div>
      <div class="flex items-center gap-2 px-4 py-2 rounded-xl transition-all"
           :class="wsConnected ? 'bg-green-500/10 border border-green-500/20' : 'bg-red-500/10 border border-red-500/20'">
        <div class="w-2 h-2 rounded-full animate-pulse"
             :class="wsConnected ? 'bg-green-400' : 'bg-red-400'"></div>
        <span class="text-sm" :class="wsConnected ? 'text-green-400' : 'text-red-400'">
          {{ wsConnected ? '实时连接' : '连接断开' }}
        </span>
      </div>
    </div>

    <!-- 磁盘列表 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div
        v-for="disk in disks"
        :key="disk.device"
        class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-3">
            <div class="w-14 h-14 bg-gradient-to-br from-blue-500/20 to-cyan-600/20 rounded-2xl flex items-center justify-center">
              <svg class="w-7 h-7 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
              </svg>
            </div>
            <div>
              <p class="text-white font-semibold text-lg">{{ disk.mountpoint || '未挂载' }}</p>
              <p class="text-gray-500 text-sm">{{ disk.device }}</p>
            </div>
          </div>
          <div class="px-3 py-1.5 rounded-full border"
               :class="getDiskStatusClass(disk.usedPercent)">
            <span class="text-sm font-medium" :class="getDiskStatusTextClass(disk.usedPercent)">
              {{ disk.fstype }}
            </span>
          </div>
        </div>

        <div class="mb-5">
          <div class="flex items-center justify-between mb-3">
            <span class="text-gray-400 text-sm">使用率</span>
            <span class="text-white font-bold text-lg">{{ disk.usedPercent }}%</span>
          </div>
          <div class="h-3 bg-gray-800 rounded-full overflow-hidden shadow-inner">
            <div
              class="h-full rounded-full transition-all duration-700 ease-out relative overflow-hidden"
              :class="getDiskProgressClass(disk.usedPercent)"
              :style="{ width: disk.usedPercent + '%' }"
            >
              <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent animate-shimmer"></div>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4 mb-5">
          <div class="bg-gray-900/50 rounded-xl p-3 border border-gray-800">
            <p class="text-gray-500 text-xs mb-1">总容量</p>
            <p class="text-white font-semibold">{{ formatBytes(disk.total) }}</p>
          </div>
          <div class="bg-gray-900/50 rounded-xl p-3 border border-gray-800">
            <p class="text-gray-500 text-xs mb-1">已用</p>
            <p class="text-white font-semibold">{{ formatBytes(disk.used) }}</p>
          </div>
          <div class="bg-gray-900/50 rounded-xl p-3 border border-gray-800">
            <p class="text-gray-500 text-xs mb-1">可用</p>
            <p class="text-green-400 font-semibold">{{ formatBytes(disk.free) }}</p>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div class="bg-gradient-to-br from-green-500/10 to-emerald-600/10 rounded-xl p-4 border border-green-500/20">
            <div class="flex items-center gap-2 mb-2">
              <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
              </svg>
              <span class="text-gray-400 text-sm">读取速度</span>
            </div>
            <p class="text-white font-bold text-lg">{{ formatSpeed(disk.readSpeed) }}</p>
            <p class="text-gray-500 text-xs mt-1">每秒读取</p>
          </div>
          <div class="bg-gradient-to-br from-blue-500/10 to-indigo-600/10 rounded-xl p-4 border border-blue-500/20">
            <div class="flex items-center gap-2 mb-2">
              <svg class="w-5 h-5 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
              <span class="text-gray-400 text-sm">写入速度</span>
            </div>
            <p class="text-white font-bold text-lg">{{ formatSpeed(disk.writeSpeed) }}</p>
            <p class="text-gray-500 text-xs mt-1">每秒写入</p>
          </div>
        </div>
      </div>
    </div>

    <!-- 磁盘 IO 历史 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">磁盘 IO 历史</h3>
          <p class="text-gray-500 text-sm">实时监控磁盘读写速度 (最近60秒)</p>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-green-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">读取</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-blue-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">写入</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-purple-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">总 IO</span>
          </div>
        </div>
      </div>
      <div class="h-80">
        <ApexChart ref="ioChartRef" type="line" :options="ioChartOptions" :series="ioChartSeries" :height="320" />
      </div>
    </div>

    <!-- 磁盘统计 -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div class="bg-gradient-to-br from-blue-500/10 to-cyan-600/10 backdrop-blur rounded-2xl p-6 border border-blue-500/20">
        <div class="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">磁盘数量</p>
        <p class="text-3xl font-bold text-white">{{ disks.length }}</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <div class="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总读取速度</p>
        <p class="text-3xl font-bold text-white">{{ formatSpeed(totalReadSpeed) }}</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <div class="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总写入速度</p>
        <p class="text-3xl font-bold text-white">{{ formatSpeed(totalWriteSpeed) }}</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <div class="w-12 h-12 bg-purple-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总 IO 速度</p>
        <p class="text-3xl font-bold text-white">{{ formatSpeed(totalIO) }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { monitorApi } from '../../api'
import { wsService } from '../../services/websocket'
import ApexChart from '../../components/Chart.vue'

interface Disk {
  device: string
  mountpoint: string
  total: number
  used: number
  free: number
  usedPercent: number
  fstype: string
  readSpeed: number
  writeSpeed: number
}

const wsConnected = ref(false)
const disks = ref<Disk[]>([])

// 计算属性
const totalReadSpeed = computed(() => disks.value.reduce((sum, disk) => sum + (disk.readSpeed || 0), 0))
const totalWriteSpeed = computed(() => disks.value.reduce((sum, disk) => sum + (disk.writeSpeed || 0), 0))
const totalIO = computed(() => totalReadSpeed.value + totalWriteSpeed.value)

// IO 图表数据
const ioChartSeries = ref([
  {
    name: '读取速度',
    data: [] as number[],
  },
  {
    name: '写入速度',
    data: [] as number[],
  },
  {
    name: '总 IO',
    data: [] as number[],
  },
])

const ioChartOptions = ref({
  chart: {
    type: 'line',
    height: 320,
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 500,
      dynamicAnimation: {
        enabled: true,
        speed: 350,
      },
    },
    toolbar: { show: false },
    zoom: { enabled: false },
  },
  colors: ['#22c55e', '#3b82f6', '#a855f7'],
  dataLabels: { enabled: false },
  stroke: {
    curve: 'smooth',
    width: 2.5,
  },
  xaxis: {
    categories: [] as string[],
    labels: {
      style: { colors: '#6b7280', fontSize: '11px' },
      formatter: (val: string) => {
        const date = new Date(val)
        return date.toLocaleTimeString('zh-CN', { minute: '2-digit', second: '2-digit' })
      },
    },
    axisBorder: { show: false },
    axisTicks: { show: false },
  },
  yaxis: {
    labels: {
      style: { colors: '#6b7280' },
      formatter: (val: number) => formatSpeed(val),
    },
  },
  grid: {
    borderColor: '#374151',
    strokeDashArray: 4,
  },
  tooltip: {
    theme: 'dark',
    y: { formatter: (val: number) => formatSpeed(val) },
  },
  legend: {
    show: true,
    position: 'top',
    horizontalAlign: 'right',
    labels: { colors: '#9ca3af' },
  },
})

// 历史数据
const ioHistory = ref({ read: [] as number[], write: [] as number[] })
const historyLabels = ref<string[]>([])
const maxDataPoints = 60

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const formatSpeed = (bytesPerSec: number) => {
  if (!bytesPerSec || bytesPerSec === 0) return '0 B/s'
  const k = 1024
  const sizes = ['B/s', 'KB/s', 'MB/s', 'GB/s']
  const i = Math.floor(Math.log(bytesPerSec) / Math.log(k))
  return Math.round((bytesPerSec / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const getDiskProgressClass = (percent: number) => {
  if (percent >= 90) return 'bg-gradient-to-r from-red-500 to-red-600'
  if (percent >= 70) return 'bg-gradient-to-r from-yellow-500 to-orange-500'
  return 'bg-gradient-to-r from-green-500 to-emerald-500'
}

const getDiskStatusClass = (percent: number) => {
  if (percent >= 90) return 'bg-red-500/10 border-red-500/20'
  if (percent >= 70) return 'bg-yellow-500/10 border-yellow-500/20'
  return 'bg-green-500/10 border-green-500/20'
}

const getDiskStatusTextClass = (percent: number) => {
  if (percent >= 90) return 'text-red-400'
  if (percent >= 70) return 'text-yellow-400'
  return 'text-green-400'
}

const updateIOChartData = (readSpeed: number, writeSpeed: number) => {
  const now = new Date()
  const timeLabel = now.toISOString()

  ioHistory.value.read.push(readSpeed)
  ioHistory.value.write.push(writeSpeed)
  historyLabels.value.push(timeLabel)

  if (ioHistory.value.read.length > maxDataPoints) {
    ioHistory.value.read.shift()
    ioHistory.value.write.shift()
    historyLabels.value.shift()
  }

  ioChartSeries.value[0].data = [...ioHistory.value.read]
  ioChartSeries.value[1].data = [...ioHistory.value.write]
  ioChartSeries.value[2].data = ioHistory.value.read.map((r, i) => r + ioHistory.value.write[i])
  ioChartOptions.value.xaxis.categories = [...historyLabels.value]
}

const handleDiskData = (data: any) => {
  disks.value = data.disks?.map((d: any) => ({
    device: d.device,
    mountpoint: d.mountpoint,
    total: d.total,
    used: d.used,
    free: d.free,
    usedPercent: Math.round((d.used / d.total) * 100),
    fstype: d.fstype,
    readSpeed: d.readSpeed || 0,
    writeSpeed: d.writeSpeed || 0,
  })) || []

  // 计算总 IO 速度
  let totalReadSpeed = 0
  let totalWriteSpeed = 0
  for (const disk of disks.value) {
    totalReadSpeed += disk.readSpeed
    totalWriteSpeed += disk.writeSpeed
  }

  updateIOChartData(totalReadSpeed, totalWriteSpeed)
}

// WebSocket 订阅
let unsubscribe: (() => void) | null = null
let unsubscribeConnection: (() => void) | null = null

const fetchInitialData = async () => {
  try {
    const data = await monitorApi.getDisk()
    handleDiskData(data)
  } catch (error) {
    console.error('获取磁盘信息失败:', error)
  }
}

onMounted(async () => {
  await fetchInitialData()

  if (!wsService.isConnected.value) {
    try {
      await wsService.connect()
    } catch (error) {
      console.error('WebSocket 连接失败:', error)
    }
  }

  unsubscribe = wsService.subscribe('disk', (message) => {
    handleDiskData(message.data)
  })

  unsubscribeConnection = wsService.onConnectionChange((connected) => {
    wsConnected.value = connected
  })

  wsConnected.value = wsService.isConnected.value
})

onUnmounted(() => {
  if (unsubscribe) unsubscribe()
  if (unsubscribeConnection) unsubscribeConnection()
})
</script>

<style scoped>
@keyframes shimmer {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

.animate-shimmer {
  animation: shimmer 2s infinite;
}
</style>
