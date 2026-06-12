<template>
  <div class="space-y-6">
    <!-- 连接状态指示器 -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-white">网络监控</h2>
        <p class="text-gray-500">实时监控网络接口和流量</p>
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

    <!-- 网络接口 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-lg font-semibold text-white">网络接口</h3>
            <p class="text-gray-500 text-sm">系统网络适配器状态</p>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
            <span class="text-green-400 text-sm">运行中</span>
          </div>
        </div>
        <div class="space-y-3">
          <div
            v-for="iface in interfaces"
            :key="iface.name"
            class="bg-gray-900/50 rounded-xl p-4 border border-gray-800 hover:border-blue-500/30 transition-all group"
          >
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 bg-gradient-to-br from-blue-500/20 to-cyan-600/20 rounded-xl flex items-center justify-center group-hover:from-blue-500/30 group-hover:to-cyan-600/30 transition-all">
                  <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
                  </svg>
                </div>
                <div>
                  <p class="text-white font-medium">{{ iface.name }}</p>
                  <p class="text-gray-500 text-sm">{{ iface.addresses?.[0] || '无 IP 地址' }}</p>
                </div>
              </div>
              <span
                class="px-3 py-1 text-xs rounded-full border transition-all"
                :class="iface.up ? 'bg-green-500/10 border-green-500/20 text-green-400' : 'bg-red-500/10 border-red-500/20 text-red-400'"
              >
                {{ iface.up ? 'UP' : 'DOWN' }}
              </span>
            </div>
            <div class="grid grid-cols-2 gap-3">
              <div class="flex items-center gap-2 bg-gray-800/50 rounded-lg p-2">
                <svg class="w-4 h-4 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
                </svg>
                <span class="text-gray-400 text-xs">上传:</span>
                <span class="text-green-400 font-medium text-sm">{{ formatBytes(iface.bytesSent) }}</span>
              </div>
              <div class="flex items-center gap-2 bg-gray-800/50 rounded-lg p-2">
                <svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                </svg>
                <span class="text-gray-400 text-xs">下载:</span>
                <span class="text-blue-400 font-medium text-sm">{{ formatBytes(iface.bytesRecv) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 实时流量 -->
      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-lg font-semibold text-white">实时流量</h3>
            <p class="text-gray-500 text-sm">当前网络传输速度</p>
          </div>
          <div class="flex items-center gap-2 px-3 py-1 bg-green-500/10 border border-green-500/20 rounded-full">
            <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
            <span class="text-green-400 text-sm">实时</span>
          </div>
        </div>
        <div class="space-y-5">
          <div class="bg-gradient-to-br from-green-500/10 to-emerald-600/10 rounded-xl p-5 border border-green-500/20">
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center">
                  <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
                  </svg>
                </div>
                <div>
                  <span class="text-gray-400 text-sm">上传速度</span>
                  <p class="text-green-400 font-bold text-2xl">{{ formatSpeed(totalSentSpeed) }}/s</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-gray-500 text-xs">峰值</p>
                <p class="text-white font-medium">{{ formatSpeed(maxSentSpeed) }}/s</p>
              </div>
            </div>
            <div class="h-3 bg-gray-800 rounded-full overflow-hidden shadow-inner">
              <div class="h-full bg-gradient-to-r from-green-500 to-emerald-500 rounded-full transition-all duration-300 relative overflow-hidden"
                   :style="{ width: Math.min((totalSentSpeed / (maxSpeed || 1)) * 100, 100) + '%' }">
                <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent animate-shimmer"></div>
              </div>
            </div>
          </div>

          <div class="bg-gradient-to-br from-blue-500/10 to-indigo-600/10 rounded-xl p-5 border border-blue-500/20">
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-3">
                <div class="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center">
                  <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
                  </svg>
                </div>
                <div>
                  <span class="text-gray-400 text-sm">下载速度</span>
                  <p class="text-blue-400 font-bold text-2xl">{{ formatSpeed(totalRecvSpeed) }}/s</p>
                </div>
              </div>
              <div class="text-right">
                <p class="text-gray-500 text-xs">峰值</p>
                <p class="text-white font-medium">{{ formatSpeed(maxRecvSpeed) }}/s</p>
              </div>
            </div>
            <div class="h-3 bg-gray-800 rounded-full overflow-hidden shadow-inner">
              <div class="h-full bg-gradient-to-r from-blue-500 to-indigo-500 rounded-full transition-all duration-300 relative overflow-hidden"
                   :style="{ width: Math.min((totalRecvSpeed / (maxSpeed || 1)) * 100, 100) + '%' }">
                <div class="absolute inset-0 bg-gradient-to-r from-white/20 to-transparent animate-shimmer"></div>
              </div>
            </div>
          </div>

          <div class="bg-gradient-to-br from-purple-500/10 to-pink-600/10 rounded-xl p-4 border border-purple-500/20">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 bg-purple-500/20 rounded-xl flex items-center justify-center">
                  <svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                  </svg>
                </div>
                <div>
                  <span class="text-gray-400 text-xs">总流量速度</span>
                  <p class="text-white font-bold">{{ formatSpeed(totalSentSpeed + totalRecvSpeed) }}/s</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 网络流量趋势 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">网络流量趋势</h3>
          <p class="text-gray-500 text-sm">实时监控网络流量 (最近60秒)</p>
        </div>
        <div class="flex items-center gap-4">
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-green-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">上传</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-blue-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">下载</span>
          </div>
          <div class="flex items-center gap-2">
            <div class="w-3 h-3 bg-purple-400 rounded-full"></div>
            <span class="text-gray-400 text-sm">总计</span>
          </div>
        </div>
      </div>
      <div class="h-80">
        <ApexChart ref="chartRef" type="line" :options="chartOptions" :series="chartSeries" height="320" />
      </div>
    </div>

    <!-- 网络统计 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-gradient-to-br from-green-500/10 to-emerald-600/10 backdrop-blur rounded-2xl p-6 border border-green-500/20 hover:border-green-500/40 transition-all">
        <div class="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总上传</p>
        <p class="text-2xl font-bold text-white">{{ formatBytes(totalSent) }}</p>
        <p class="text-gray-500 text-xs mt-2">历史累计</p>
      </div>

      <div class="bg-gradient-to-br from-blue-500/10 to-indigo-600/10 backdrop-blur rounded-2xl p-6 border border-blue-500/20 hover:border-blue-500/40 transition-all">
        <div class="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总下载</p>
        <p class="text-2xl font-bold text-white">{{ formatBytes(totalRecv) }}</p>
        <p class="text-gray-500 text-xs mt-2">历史累计</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="w-12 h-12 bg-purple-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">活跃接口</p>
        <p class="text-2xl font-bold text-white">{{ activeInterfaces }}</p>
        <p class="text-gray-500 text-xs mt-2">共 {{ interfaces.length }} 个</p>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50 hover:border-gray-600/50 transition-all">
        <div class="w-12 h-12 bg-pink-500/20 rounded-xl flex items-center justify-center mb-4">
          <svg class="w-6 h-6 text-pink-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <p class="text-gray-400 text-sm mb-1">总流量</p>
        <p class="text-2xl font-bold text-white">{{ formatBytes(totalSent + totalRecv) }}</p>
        <p class="text-gray-500 text-xs mt-2">上传 + 下载</p>
      </div>
    </div>

    <!-- 网络流量饼图 -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <h3 class="text-lg font-semibold text-white mb-4">流量分布</h3>
        <div class="h-64">
          <ApexChart type="donut" :options="donutOptions" :series="donutSeries" height="256" />
        </div>
      </div>

      <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
        <h3 class="text-lg font-semibold text-white mb-4">接口详情</h3>
        <div class="space-y-3">
          <div
            v-for="iface in interfaces"
            :key="iface.name"
            class="bg-gray-900/50 rounded-xl p-4 border border-gray-800"
          >
            <div class="flex items-center justify-between mb-2">
              <span class="text-white font-medium">{{ iface.name }}</span>
              <span class="text-gray-400 text-sm">{{ iface.up ? '在线' : '离线' }}</span>
            </div>
            <div class="grid grid-cols-2 gap-2 text-sm">
              <div>
                <span class="text-gray-500">IP: </span>
                <span class="text-gray-300">{{ iface.addresses?.[0] || '-' }}</span>
              </div>
              <div>
                <span class="text-gray-500">MAC: </span>
                <span class="text-gray-300">{{ iface.mac || '-' }}</span>
              </div>
            </div>
          </div>
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

interface Interface {
  name: string
  up: boolean
  addresses: string[]
  mac?: string
  bytesSent: number
  bytesRecv: number
  sentSpeed: number
  recvSpeed: number
}

const wsConnected = ref(false)
const interfaces = ref<Interface[]>([])
const maxSpeed = ref(10 * 1024 * 1024) // 10 MB/s
const maxSentSpeed = ref(0)
const maxRecvSpeed = ref(0)

// 图表数据
const chartSeries = ref([
  {
    name: '上传速度',
    data: [] as number[],
  },
  {
    name: '下载速度',
    data: [] as number[],
  },
  {
    name: '总计',
    data: [] as number[],
  },
])

const chartOptions = ref({
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

// 圆环图
const donutSeries = computed(() => [totalSent.value, totalRecv.value])

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
  colors: ['#22c55e', '#3b82f6'],
  labels: ['总上传', '总下载'],
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
            fontSize: '20px',
            fontWeight: 600,
            formatter: (val: string) => formatBytes(parseFloat(val)),
          },
          total: {
            show: true,
            label: '总流量',
            color: '#9ca3af',
            formatter: () => formatBytes(totalSent.value + totalRecv.value),
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
const historyData = ref({ sent: [] as number[], recv: [] as number[] })
const historyLabels = ref<string[]>([])
const maxDataPoints = 60

const totalSent = computed(() => interfaces.value.reduce((sum, iface) => sum + iface.bytesSent, 0))
const totalRecv = computed(() => interfaces.value.reduce((sum, iface) => sum + iface.bytesRecv, 0))
const totalSentSpeed = computed(() => interfaces.value.reduce((sum, iface) => sum + (iface.sentSpeed || 0), 0))
const totalRecvSpeed = computed(() => interfaces.value.reduce((sum, iface) => sum + (iface.recvSpeed || 0), 0))
const activeInterfaces = computed(() => interfaces.value.filter(i => i.up).length)

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
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

const updateChartData = (sentSpeed: number, recvSpeed: number) => {
  const now = new Date()
  const timeLabel = now.toISOString()

  // 更新峰值
  if (sentSpeed > maxSentSpeed.value) maxSentSpeed.value = sentSpeed
  if (recvSpeed > maxRecvSpeed.value) maxRecvSpeed.value = recvSpeed

  historyData.value.sent.push(sentSpeed)
  historyData.value.recv.push(recvSpeed)
  historyLabels.value.push(timeLabel)

  if (historyData.value.sent.length > maxDataPoints) {
    historyData.value.sent.shift()
    historyData.value.recv.shift()
    historyLabels.value.shift()
  }

  chartSeries.value[0].data = [...historyData.value.sent]
  chartSeries.value[1].data = [...historyData.value.recv]
  chartSeries.value[2].data = historyData.value.sent.map((s, i) => s + historyData.value.recv[i])
  chartOptions.value.xaxis.categories = [...historyLabels.value]

  // 动态调整最大速度
  const currentMax = Math.max(sentSpeed, recvSpeed)
  if (currentMax > maxSpeed.value) {
    maxSpeed.value = currentMax * 1.5
  }
}

const handleNetworkData = (data: any) => {
  interfaces.value = data.interfaces?.filter((i: any) => i.name !== 'lo').map((iface: any) => ({
    name: iface.name,
    up: iface.up,
    addresses: iface.addresses || [],
    mac: iface.mac || '',
    bytesSent: iface.bytesSent || 0,
    bytesRecv: iface.bytesRecv || 0,
    sentSpeed: iface.sentSpeed || 0,
    recvSpeed: iface.recvSpeed || 0,
  })) || []

  updateChartData(totalSentSpeed.value, totalRecvSpeed.value)
}

// WebSocket 订阅
let unsubscribe: (() => void) | null = null
let unsubscribeConnection: (() => void) | null = null

const fetchInitialData = async () => {
  try {
    const data = await monitorApi.getNetwork()
    handleNetworkData(data)
  } catch (error) {
    console.error('获取网络信息失败:', error)
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

  unsubscribe = wsService.subscribe('network', (message) => {
    handleNetworkData(message.data)
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
