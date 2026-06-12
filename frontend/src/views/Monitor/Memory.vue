<template>
  <div class="space-y-6">
    <!-- 连接状态指示器 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-white">内存监控</h2>
        <p class="text-gray-500">实时监控内存使用情况</p>
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

    <!-- 内存概览 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-gradient-to-br from-purple-500/10 to-pink-600/10 backdrop-blur rounded-2xl p-6 border border-purple-500/20 hover:border-purple-500/40 transition-all">
        <div class="w-14 h-14 bg-gradient-to-br from-purple-500/20 to-pink-600/20 rounded-2xl flex items-center justify-center mb-4">
          <svg class="w-7 h-7 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总内存</p>
        <p class="text-3xl font-bold text-white">{{ formatBytes(memory.total) }}</p>
        <p class="text-gray-500 text-sm mt-2">{{ memory.cores || 'N/A' }} 位系统</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="w-14 h-14 bg-red-500/20 rounded-2xl flex items-center justify-center mb-4">
          <svg class="w-7 h-7 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">已使用</p>
        <p class="text-3xl font-bold text-white">{{ formatBytes(memory.used) }}</p>
        <div class="mt-4 h-3 bg-gray-800 rounded-full overflow-hidden shadow-inner">
          <div class="h-full bg-gradient-to-r from-red-500 to-orange-500 rounded-full transition-all duration-700 ease-out"
               :style="{ width: memory.percent + '%' }"></div>
        </div>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="w-14 h-14 bg-green-500/20 rounded-2xl flex items-center justify-center mb-4">
          <svg class="w-7 h-7 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">可用</p>
        <p class="text-3xl font-bold text-green-400">{{ formatBytes(memory.available) }}</p>
        <p class="text-gray-500 text-sm mt-2">{{ memory.percent }}% 已用</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="w-14 h-14 bg-blue-500/20 rounded-2xl flex items-center justify-center mb-4">
          <svg class="w-7 h-7 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">使用率</p>
        <p class="text-3xl font-bold text-white">{{ memory.percent }}%</p>
        <div class="mt-4 h-3 bg-gray-800 rounded-full overflow-hidden shadow-inner">
          <div class="h-full bg-gradient-to-r from-purple-500 to-pink-500 rounded-full transition-all duration-700 ease-out relative overflow-hidden"
               :style="{ width: memory.percent + '%' }">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent animate-shimmer"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 内存使用历史图表 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">内存使用趋势</h3>
          <p class="text-gray-500 text-sm">实时监控内存使用情况 (最近60秒)</p>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-purple-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">使用率</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-blue-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">平均值</span>
          </div>
        </div>
      </div>
      <div class="h-80">
        <ApexChart ref="chartRef" type="area" :options="chartOptions" :series="chartSeries" height="320" />
      </div>
    </div>

    <!-- Swap 状态 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
      <h3 class="text-lg font-semibold text-white mb-4">Swap 内存状态</h3>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-gradient-to-br from-blue-500/10 to-cyan-600/10 rounded-xl p-5 border border-blue-500/20">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
              </svg>
            </div>
            <div>
              <p class="text-gray-400 text-sm">Swap 总量</p>
              <p class="text-white font-bold text-xl">{{ formatBytes(swap.total) }}</p>
            </div>
          </div>
          <div class="h-2 bg-gray-800 rounded-full overflow-hidden">
            <div class="h-full bg-blue-500 transition-all duration-500" :style="{ width: '100%' }"></div>
          </div>
        </div>

        <div class="bg-gradient-to-br from-yellow-500/10 to-orange-600/10 rounded-xl p-5 border border-yellow-500/20">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-12 h-12 bg-yellow-500/20 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div>
              <p class="text-gray-400 text-sm">Swap 已用</p>
              <p class="text-white font-bold text-xl">{{ formatBytes(swap.used) }}</p>
            </div>
          </div>
          <div class="h-2 bg-gray-800 rounded-full overflow-hidden">
            <div class="h-full bg-yellow-500 transition-all duration-500"
                 :style="{ width: (swap.total ? (swap.used / swap.total) * 100 : 0) + '%' }"></div>
          </div>
        </div>

        <div class="bg-gradient-to-br from-green-500/10 to-emerald-600/10 rounded-xl p-5 border border-green-500/20">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <p class="text-gray-400 text-sm">Swap 使用率</p>
              <p class="text-white font-bold text-xl">{{ swap.percent }}%</p>
            </div>
          </div>
          <div class="h-2 bg-gray-800 rounded-full overflow-hidden">
            <div class="h-full bg-gradient-to-r from-green-500 to-emerald-500 transition-all duration-500"
                 :style="{ width: swap.percent + '%' }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- 内存详细信息 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <h3 class="text-lg font-semibold text-white mb-4">内存详细信息</h3>
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 bg-cyan-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-4 h-4 text-cyan-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <p class="text-gray-400 text-sm">Cached</p>
            </div>
            <p class="text-white font-bold text-lg">{{ formatBytes(memory.cached || 0) }}</p>
          </div>
          <div class="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 bg-teal-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-4 h-4 text-teal-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
                </svg>
              </div>
              <p class="text-gray-400 text-sm">Buffers</p>
            </div>
            <p class="text-white font-bold text-lg">{{ formatBytes(memory.buffers || 0) }}</p>
          </div>
          <div class="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 bg-indigo-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <p class="text-gray-400 text-sm">Shared</p>
            </div>
            <p class="text-white font-bold text-lg">{{ formatBytes(memory.shared || 0) }}</p>
          </div>
          <div class="bg-gray-900/50 rounded-xl p-4 border border-gray-800">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 bg-emerald-500/20 rounded-lg flex items-center justify-center">
                <svg class="w-4 h-4 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </div>
              <p class="text-gray-400 text-sm">Free</p>
            </div>
            <p class="text-white font-bold text-lg">{{ formatBytes(memory.free || 0) }}</p>
          </div>
        </div>
      </div>

      <!-- 内存分配图表 -->
      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <h3 class="text-lg font-semibold text-white mb-4">内存分配</h3>
        <div class="h-64">
          <ApexChart type="donut" :options="donutOptions" :series="donutSeries" height="256" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { monitorApi } from '../../api'
import { wsService } from '../../services/websocket'
import ApexChart from '../../components/Chart.vue'

interface Memory {
  total: number
  used: number
  available: number
  percent: number
  free: number
  cached: number
  buffers: number
  shared: number
}

interface Swap {
  total: number
  used: number
  percent: number
}

const wsConnected = ref(false)

const memory = ref<Memory>({
  total: 0,
  used: 0,
  available: 0,
  percent: 0,
  free: 0,
  cached: 0,
  buffers: 0,
  shared: 0,
})

const swap = ref<Swap>({
  total: 0,
  used: 0,
  percent: 0,
})

// 主图表数据
const chartSeries = ref([
  {
    name: '内存使用率',
    data: [] as number[],
  },
  {
    name: '平均使用率',
    data: [] as number[],
  },
])

const chartOptions = ref({
  chart: {
    type: 'area',
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
  colors: ['#a855f7', '#3b82f6'],
  fill: {
    type: 'gradient',
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.6,
      opacityTo: 0.05,
      stops: [0, 100],
    },
  },
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
    min: 0,
    max: 100,
    labels: {
      style: { colors: '#6b7280' },
      formatter: (val: number) => val + '%',
    },
  },
  grid: {
    borderColor: '#374151',
    strokeDashArray: 4,
  },
  tooltip: {
    theme: 'dark',
    y: { formatter: (val: number) => val.toFixed(1) + '%' },
  },
  legend: {
    show: true,
    position: 'top',
    horizontalAlign: 'right',
    labels: { colors: '#9ca3af' },
  },
})

// 圆环图数据
const donutSeries = computed(() => [
  memory.value.used - (memory.value.cached || 0) - (memory.value.buffers || 0),
  memory.value.cached || 0,
  memory.value.buffers || 0,
  memory.value.free || 0,
])

const donutOptions = ref({
  chart: {
    type: 'donut',
    height: 256,
    animations: {
      enabled: true,
      easing: 'easeinout',
      speed: 500,
    },
    toolbar: { show: false },
  },
  colors: ['#ef4444', '#f59e0b', '#3b82f6', '#22c55e'],
  labels: ['应用程序', 'Cached', 'Buffers', 'Free'],
  dataLabels: { enabled: false },
  stroke: { show: false },
  plotOptions: {
    pie: {
      donut: {
        size: '70%',
        labels: {
          show: true,
          name: { show: true, color: '#9ca3af' },
          value: {
            show: true,
            color: '#fff',
            fontSize: '24px',
            fontWeight: 600,
          },
          total: {
            show: true,
            label: '总内存',
            color: '#9ca3af',
            formatter: () => formatBytes(memory.value.total),
          },
        },
      },
    },
  },
  grid: { padding: { top: 0, bottom: 0, left: 0, right: 0 } },
  legend: {
    show: true,
    position: 'bottom',
    labels: { colors: '#9ca3af' },
  },
  tooltip: {
    theme: 'dark',
    y: {
      formatter: (val: number) => formatBytes(val),
    },
  },
})

// 历史数据
const historyData = ref<number[]>([])
const historyLabels = ref<string[]>([])
const maxDataPoints = 60

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

const calculateAverage = () => {
  if (historyData.value.length === 0) return 0
  const sum = historyData.value.reduce((a, b) => a + b, 0)
  return Math.round(sum / historyData.value.length)
}

const updateChartData = (usage: number) => {
  const now = new Date()
  const timeLabel = now.toISOString()

  historyData.value.push(usage)
  historyLabels.value.push(timeLabel)

  if (historyData.value.length > maxDataPoints) {
    historyData.value.shift()
    historyLabels.value.shift()
  }

  const avgUsage = calculateAverage()

  chartSeries.value[0].data = [...historyData.value]
  chartSeries.value[1].data = historyData.value.map(() => avgUsage)
  chartOptions.value.xaxis.categories = [...historyLabels.value]
}

const handleMemoryData = (data: any) => {
  const percent = Math.round((data.used / data.total) * 100)

  memory.value = {
    total: data.total,
    used: data.used,
    available: data.available,
    percent,
    free: data.free || 0,
    cached: data.cached || 0,
    buffers: data.buffers || 0,
    shared: data.shared || 0,
  }

  swap.value = {
    total: data.swapTotal || 0,
    used: data.swapUsed || 0,
    percent: data.swapTotal ? Math.round((data.swapUsed / data.swapTotal) * 100) : 0,
  }

  updateChartData(percent)
}

// WebSocket 订阅
let unsubscribe: (() => void) | null = null
let unsubscribeConnection: (() => void) | null = null

const fetchInitialData = async () => {
  try {
    const data = await monitorApi.getMemory()
    handleMemoryData(data)
  } catch (error) {
    console.error('获取内存信息失败:', error)
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

  unsubscribe = wsService.subscribe('memory', (message) => {
    handleMemoryData(message.data)
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
