<template>
  <div class="space-y-6">
    <!-- 连接状态指示器 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-white">CPU 监控</h2>
        <p class="text-gray-500">实时监控 CPU 使用情况和系统负载</p>
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

    <!-- CPU 概览卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-gradient-to-br from-indigo-500/10 to-purple-600/10 backdrop-blur rounded-2xl p-6 border border-indigo-500/20 hover:border-indigo-500/40 transition-all">
        <div class="flex items-center justify-between mb-4">
          <div class="w-14 h-14 bg-gradient-to-br from-indigo-500/20 to-purple-600/20 rounded-2xl flex items-center justify-center">
            <svg class="w-7 h-7 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
          </div>
          <div class="px-3 py-1.5 rounded-full transition-all"
               :class="getStatusColorClass(cpuInfo.usage)">
            <span class="text-sm font-medium" :class="getStatusTextColor(cpuInfo.usage)">
              {{ getStatusText(cpuInfo.usage) }}
            </span>
          </div>
        </div>
        <p class="text-gray-400 text-sm mb-1">CPU 使用率</p>
        <div class="flex items-end gap-2 mb-4">
          <p class="text-5xl font-bold text-white">{{ (cpuInfo.usage * 100).toFixed(1) }}%</p>
          <p class="text-gray-500 text-sm mb-2" v-if="cpuUsageChange !== 0">
            <span :class="cpuUsageChange > 0 ? 'text-red-400' : 'text-green-400'">
              {{ cpuUsageChange > 0 ? '+' : '' }}{{ cpuUsageChange }}%
            </span>
          </p>
        </div>
        <div class="h-3 bg-gray-800 rounded-full overflow-hidden shadow-inner">
          <div class="h-full bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500 rounded-full transition-all duration-700 ease-out relative overflow-hidden"
               :style="{ width: (cpuInfo.usage * 100) + '%' }">
            <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent animate-shimmer"></div>
          </div>
        </div>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="flex items-center justify-between mb-4">
          <div class="w-14 h-14 bg-purple-500/20 rounded-2xl flex items-center justify-center">
            <svg class="w-7 h-7 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16v11H4V10z" />
            </svg>
          </div>
        </div>
        <p class="text-gray-400 text-sm mb-1">CPU 核心数</p>
        <p class="text-5xl font-bold text-white mb-4">{{ cpuInfo.cores }}</p>
        <p class="text-gray-500 text-sm truncate" :title="cpuInfo.model">{{ cpuInfo.model }}</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="flex items-center justify-between mb-4">
          <div class="w-14 h-14 bg-blue-500/20 rounded-2xl flex items-center justify-center">
            <svg class="w-7 h-7 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
            </svg>
          </div>
        </div>
        <p class="text-gray-400 text-sm mb-1">系统负载</p>
        <p class="text-5xl font-bold text-white mb-4">{{ cpuInfo.load }}</p>
        <div class="flex gap-4 text-xs">
          <div class="text-center">
            <p class="text-gray-500">1分钟</p>
            <p class="text-white font-medium">{{ cpuInfo.load1 }}</p>
          </div>
          <div class="text-center">
            <p class="text-gray-500">5分钟</p>
            <p class="text-white font-medium">{{ cpuInfo.load5 }}</p>
          </div>
          <div class="text-center">
            <p class="text-gray-500">15分钟</p>
            <p class="text-white font-medium">{{ cpuInfo.load15 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- CPU 使用历史图表 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">CPU 使用率历史</h3>
          <p class="text-gray-500 text-sm">实时监控 CPU 使用情况 (最近60秒)</p>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-indigo-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">当前</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-purple-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">平均</span>
          </div>
        </div>
      </div>
      <div class="h-80">
        <ApexChart ref="chartRef" type="area" :options="chartOptions" :series="chartSeries" :height="320" />
      </div>
    </div>

    <!-- 每核心使用率 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">每个核心使用率</h3>
          <p class="text-gray-500 text-sm">实时监控各个 CPU 核心</p>
        </div>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
        <div
          v-for="core in cpuInfo.coresList"
          :key="core.id"
          class="bg-gray-900/50 rounded-xl p-4 border border-gray-800 hover:border-indigo-500/30 transition-all group"
        >
          <div class="flex items-center justify-between mb-3">
            <p class="text-gray-400 text-sm">Core {{ core.id }}</p>
            <p class="text-white font-bold">{{ core.usage.toFixed(1) }}%</p>
          </div>
          <div class="h-2.5 bg-gray-800 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-500 relative overflow-hidden"
              :class="getCoreColorClass(core.usage)"
              :style="{ width: core.usage + '%' }"
            >
              <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent animate-shimmer"></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- CPU 频率信息 -->
    <div v-if="cpuInfo.frequencies && cpuInfo.frequencies.length > 0" class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
      <h3 class="text-lg font-semibold text-white mb-4">CPU 频率</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-4">
        <div v-for="(freq, idx) in cpuInfo.frequencies" :key="idx" class="bg-gray-900/50 rounded-xl p-3 border border-gray-800">
          <p class="text-gray-500 text-xs mb-1">Core {{ idx }}</p>
          <p class="text-white font-bold">{{ (freq / 1000).toFixed(1) }} GHz</p>
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

interface CPUCore {
  id: number
  usage: number
}

interface CPUInfo {
  usage: number
  cores: number
  model: string
  load: string
  load1?: number
  load5?: number
  load15?: number
  coresList: CPUCore[]
  frequencies?: number[]
}

const wsConnected = ref(false)
const cpuInfo = ref<CPUInfo>({
  usage: 0,
  cores: 0,
  model: '',
  load: '0.00',
  load1: 0,
  load5: 0,
  load15: 0,
  coresList: [],
  frequencies: [],
})

const prevUsage = ref(0)
const cpuUsageChange = computed(() => cpuInfo.value.usage - prevUsage.value)

// 图表数据
const chartSeries = ref([
  {
    name: '当前使用率',
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
  colors: ['#818cf8', '#c084fc'],
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
    yaxis: { lines: { show: true } },
  },
  tooltip: {
    theme: 'dark',
    x: { show: true },
    y: { formatter: (val: number) => val.toFixed(1) + '%' },
    marker: { show: true },
  },
  legend: {
    show: true,
    position: 'top',
    horizontalAlign: 'right',
    labels: { colors: '#9ca3af' },
  },
})

// 历史数据
const historyData = ref<number[]>([])
const historyLabels = ref<string[]>([])
const maxDataPoints = 60

const getCoreColorClass = (usage: number) => {
  const percent = usage > 1 ? usage : usage * 100 // Convert to percentage if needed
  if (percent >= 80) return 'bg-gradient-to-r from-red-500 to-red-600'
  if (percent >= 60) return 'bg-gradient-to-r from-yellow-500 to-orange-500'
  if (percent >= 40) return 'bg-gradient-to-r from-blue-500 to-indigo-500'
  return 'bg-gradient-to-r from-green-500 to-emerald-500'
}

const getStatusColorClass = (usage: number) => {
  const percent = usage > 1 ? usage : usage * 100 // Convert to percentage if needed
  if (percent >= 80) return 'bg-red-500/10 border border-red-500/20'
  if (percent >= 60) return 'bg-yellow-500/10 border border-yellow-500/20'
  return 'bg-green-500/10 border border-green-500/20'
}

const getStatusTextColor = (usage: number) => {
  const percent = usage > 1 ? usage : usage * 100 // Convert to percentage if needed
  if (percent >= 80) return 'text-red-400'
  if (percent >= 60) return 'text-yellow-400'
  return 'text-green-400'
}

const getStatusText = (usage: number) => {
  const percent = usage > 1 ? usage : usage * 100 // Convert to percentage if needed
  if (percent >= 90) return '严重'
  if (percent >= 80) return '警告'
  if (percent >= 60) return '较高'
  return '正常'
}

const calculateAverage = () => {
  if (historyData.value.length === 0) return 0
  const sum = historyData.value.reduce((a, b) => a + b, 0)
  return Math.round(sum / historyData.value.length)
}

const updateChartData = (usage: number) => {
  const now = new Date()
  const timeLabel = now.toISOString()

  prevUsage.value = historyData.value.length > 0 ? historyData.value[historyData.value.length - 1] : 0

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

const handleCPUData = (data: any) => {
  const usage = Math.round(data.usage * 100)

  cpuInfo.value = {
    usage,
    cores: data.cores,
    model: data.model || 'Unknown Processor',
    load: data.load?.toFixed(2) || '0.00',
    load1: data.load1?.toFixed(2) ? Number(data.load1.toFixed(2)) : undefined,
    load5: data.load5?.toFixed(2) ? Number(data.load5.toFixed(2)) : undefined,
    load15: data.load15?.toFixed(2) ? Number(data.load15.toFixed(2)) : undefined,
    coresList: data.perCore?.map((u: number, idx: number) => ({
      id: idx,
      usage: Math.round(u * 100),
    })) || [],
    frequencies: data.frequencies || [],
  }

  updateChartData(usage)
}

// WebSocket 订阅
let unsubscribe: (() => void) | null = null
let unsubscribeConnection: (() => void) | null = null

const fetchInitialData = async () => {
  try {
    const data = await monitorApi.getCPU()
    handleCPUData(data)
  } catch (error) {
    console.error('获取 CPU 信息失败:', error)
  }
}

onMounted(async () => {
  // 获取初始数据
  await fetchInitialData()

  // 连接 WebSocket
  if (!wsService.isConnected.value) {
    try {
      await wsService.connect()
    } catch (error) {
      console.error('WebSocket 连接失败:', error)
    }
  }

  // 订阅 CPU 数据
  unsubscribe = wsService.subscribe('cpu', (message) => {
    handleCPUData(message.data)
  })

  // 订阅连接状态
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
