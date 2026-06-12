<template>
  <div class="space-y-6">
    <!-- 欢迎横幅 -->
    <div class="relative overflow-hidden bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-500 rounded-2xl p-6 shadow-lg shadow-indigo-500/30">
      <div class="absolute inset-0 bg-black/10"></div>
      <div class="absolute -right-10 -top-10 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
      <div class="absolute -left-10 -bottom-10 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
      <div class="relative z-10">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-white mb-1">欢迎回来，{{ username }}</h1>
            <p class="text-indigo-100 flex items-center gap-2">
              <span class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></span>
              系统运行正常 · 运行时间 {{ uptime }}
            </p>
          </div>
          <div class="hidden md:flex items-center gap-4">
            <div class="text-right">
              <p class="text-white/70 text-sm">当前时间</p>
              <p class="text-white font-semibold text-lg">{{ currentTime }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div
        v-for="stat in stats"
        :key="stat.name"
        class="group relative overflow-hidden bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50 hover:border-indigo-500/50 transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/10 hover:-translate-y-1"
        :class="stat.hoverShadow"
      >
        <div class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
             :class="stat.bgGradient"></div>
        <div class="relative z-10">
          <div class="flex items-start justify-between">
            <div>
              <p class="text-gray-400 text-sm mb-1">{{ stat.name }}</p>
              <p class="text-3xl font-bold text-white tabular-nums" :class="stat.valueColor">
                <AnimatedNumber :value="stat.value" :decimals="stat.decimals || 0" />
                <span v-if="stat.suffix" class="text-lg ml-1">{{ stat.suffix }}</span>
              </p>
              <p class="text-gray-500 text-sm mt-1">{{ stat.unit }}</p>
            </div>
            <div :class="`p-3 rounded-xl ${stat.iconBg} group-hover:scale-110 transition-transform duration-300`">
              <component :is="stat.icon" :class="`w-6 h-6 ${stat.iconColor}`" />
            </div>
          </div>
          <div class="mt-4 h-2 bg-gray-700/50 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-700 ease-out"
              :class="stat.barColor"
              :style="{ width: Math.min(stat.value, 100) + '%' }"
            ></div>
          </div>
          <div class="mt-3 flex items-center gap-2 text-sm">
            <span :class="stat.trendColor">
              <svg class="w-4 h-4 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="stat.trend === 'up'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
              </svg>
              {{ stat.trendValue }}
            </span>
            <span class="text-gray-500">较上次</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- CPU 使用率图表 -->
      <div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold text-white">CPU 使用率</h3>
            <p class="text-gray-500 text-sm">实时监控 · {{ cpuCores }} 核心</p>
          </div>
          <div class="flex items-center gap-2">
            <div class="flex items-center gap-2 px-3 py-1 rounded-full text-sm"
                 :class="cpuStatus.class">
              <div class="w-2 h-2 rounded-full animate-pulse" :class="cpuStatus.dotClass"></div>
              <span>{{ cpuStatus.text }}</span>
            </div>
          </div>
        </div>
        <div class="h-64">
          <Chart
            type="area"
            :height="260"
            :series="[{ name: 'CPU 使用率', data: cpuHistory }]"
            :options="cpuChartOptions"
          />
        </div>
      </div>

      <!-- 内存使用图表 -->
      <div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold text-white">内存使用</h3>
            <p class="text-gray-500 text-sm">总计 {{ formatBytes(totalMemory) }}</p>
          </div>
          <div class="flex items-center gap-2">
            <div class="flex items-center gap-2 px-3 py-1 rounded-full text-sm"
                 :class="memoryStatus.class">
              <div class="w-2 h-2 rounded-full animate-pulse" :class="memoryStatus.dotClass"></div>
              <span>{{ memoryStatus.text }}</span>
            </div>
          </div>
        </div>
        <div class="h-64">
          <Chart
            type="area"
            :height="260"
            :series="[{ name: '内存使用', data: memoryHistory }]"
            :options="memoryChartOptions"
          />
        </div>
      </div>

      <!-- 网络流量图表 -->
      <div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50 lg:col-span-2">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h3 class="text-lg font-semibold text-white">网络流量</h3>
            <p class="text-gray-500 text-sm">实时网络监控</p>
          </div>
          <div class="flex items-center gap-4">
            <div class="flex items-center gap-2">
              <div class="w-3 h-3 rounded-full bg-cyan-400"></div>
              <span class="text-sm text-gray-400">下载: {{ formatBytes(networkDown) }}/s</span>
            </div>
            <div class="flex items-center gap-2">
              <div class="w-3 h-3 rounded-full bg-pink-400"></div>
              <span class="text-sm text-gray-400">上传: {{ formatBytes(networkUp) }}/s</span>
            </div>
          </div>
        </div>
        <div class="h-64">
          <Chart
            type="line"
            :height="260"
            :series="networkSeries"
            :options="networkChartOptions"
          />
        </div>
      </div>
    </div>

    <!-- 系统状态概览 + 快速操作 -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- 系统状态概览 -->
      <div class="lg:col-span-2 bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50">
        <h3 class="text-lg font-semibold text-white mb-4">系统状态概览</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div
            v-for="service in services"
            :key="service.name"
            class="flex items-center justify-between p-4 bg-gray-900/50 rounded-xl hover:bg-gray-900/70 transition-colors"
          >
            <div class="flex items-center gap-3">
              <div :class="`w-10 h-10 rounded-lg flex items-center justify-center ${service.iconBg}`">
                <component :is="service.icon" :class="`w-5 h-5 ${service.iconColor}`" />
              </div>
              <div>
                <p class="text-white font-medium">{{ service.name }}</p>
                <p class="text-gray-500 text-sm">{{ service.desc }}</p>
              </div>
            </div>
            <span :class="`px-3 py-1 rounded-full text-sm ${service.statusClass}`">
              {{ service.status }}
            </span>
          </div>
        </div>
      </div>

      <!-- 快速操作 -->
      <div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50">
        <h3 class="text-lg font-semibold text-white mb-4">快速操作</h3>
        <div class="space-y-3">
          <button
            v-for="action in quickActions"
            :key="action.name"
            @click="handleQuickAction(action.action)"
            class="w-full flex items-center gap-3 p-3 rounded-xl bg-gray-900/50 hover:bg-gray-900/70 transition-all hover:scale-[1.02] active:scale-[0.98] group text-left"
          >
            <div :class="`w-10 h-10 rounded-lg flex items-center justify-center ${action.iconBg} group-hover:scale-110 transition-transform`">
              <component :is="action.icon" :class="`w-5 h-5 ${action.iconColor}`" />
            </div>
            <div class="flex-1">
              <p class="text-white font-medium">{{ action.name }}</p>
              <p class="text-gray-500 text-sm">{{ action.desc }}</p>
            </div>
            <svg class="w-5 h-5 text-gray-600 group-hover:text-gray-400 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- 存储概览 -->
    <div class="bg-gray-800/50 backdrop-blur-sm rounded-2xl p-6 border border-gray-700/50">
      <h3 class="text-lg font-semibold text-white mb-4">存储概览</h3>
      <div class="space-y-4">
        <div
          v-for="disk in disks"
          :key="disk.device"
          class="p-4 bg-gray-900/50 rounded-xl"
        >
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-blue-500/10 rounded-lg flex items-center justify-center">
                <HardDriveIcon class="w-5 h-5 text-blue-400" />
              </div>
              <div>
                <p class="text-white font-medium">{{ disk.mountpoint || disk.device }}</p>
                <p class="text-gray-500 text-sm">{{ disk.fstype }}</p>
              </div>
            </div>
            <span class="text-gray-400 text-sm">{{ formatBytes(disk.used) }} / {{ formatBytes(disk.total) }}</span>
          </div>
          <div class="h-2 bg-gray-700/50 rounded-full overflow-hidden">
            <div
              class="h-full rounded-full transition-all duration-700 bg-gradient-to-r from-blue-500 to-cyan-500"
              :style="{ width: disk.usedPercent + '%' }"
            ></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, h } from 'vue'
import { useRouter } from 'vue-router'
import {
  CpuChipIcon,
  ServerIcon,
  CircleStackIcon,
  ChartBarIcon,
  HardDriveIcon,
  SignalIcon,
  GlobeAltIcon,
  FolderIcon,
  PowerIcon,
  ArrowPathIcon,
  DocumentArrowDownIcon,
  Cog6ToothIcon
} from '@heroicons/vue/24/outline'
import Chart from '@/components/Chart.vue'
import { monitorApi } from '@/api'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 用户名
const username = computed(() => authStore.user?.username || 'Admin')
const currentTime = ref('')
const uptime = ref('--:--:--')

// 统计数据
const stats = ref([
  {
    name: 'CPU 使用率',
    value: 0,
    unit: '8 核处理器',
    icon: CpuChipIcon,
    iconBg: 'bg-indigo-500/10',
    iconColor: 'text-indigo-400',
    barColor: 'bg-gradient-to-r from-indigo-500 to-purple-500',
    bgGradient: 'bg-gradient-to-r from-indigo-500/10 to-transparent',
    valueColor: 'text-indigo-400',
    hoverShadow: 'hover:shadow-indigo-500/20',
    trend: 'down',
    trendValue: '-2%',
    trendColor: 'text-green-400',
    decimals: 1,
    suffix: '%'
  },
  {
    name: '内存使用',
    value: 0,
    unit: '16 GB 总容量',
    icon: ServerIcon,
    iconBg: 'bg-purple-500/10',
    iconColor: 'text-purple-400',
    barColor: 'bg-gradient-to-r from-purple-500 to-pink-500',
    bgGradient: 'bg-gradient-to-r from-purple-500/10 to-transparent',
    valueColor: 'text-purple-400',
    hoverShadow: 'hover:shadow-purple-500/20',
    trend: 'up',
    trendValue: '+5%',
    trendColor: 'text-red-400',
    decimals: 1,
    suffix: '%'
  },
  {
    name: '磁盘使用',
    value: 0,
    unit: '系统盘',
    icon: CircleStackIcon,
    iconBg: 'bg-blue-500/10',
    iconColor: 'text-blue-400',
    barColor: 'bg-gradient-to-r from-blue-500 to-cyan-500',
    bgGradient: 'bg-gradient-to-r from-blue-500/10 to-transparent',
    valueColor: 'text-blue-400',
    hoverShadow: 'hover:shadow-blue-500/20',
    trend: 'up',
    trendValue: '+0.1%',
    trendColor: 'text-yellow-400',
    decimals: 1,
    suffix: '%'
  },
  {
    name: '网络上传',
    value: 0,
    unit: '实时速率',
    icon: ChartBarIcon,
    iconBg: 'bg-green-500/10',
    iconColor: 'text-green-400',
    barColor: 'bg-gradient-to-r from-green-500 to-emerald-500',
    bgGradient: 'bg-gradient-to-r from-green-500/10 to-transparent',
    valueColor: 'text-green-400',
    hoverShadow: 'hover:shadow-green-500/20',
    trend: 'down',
    trendValue: '-15%',
    trendColor: 'text-green-400',
    decimals: 0,
    suffix: ' KB/s'
  }
])

// CPU 相关
const cpuHistory = ref<number[]>([])
const cpuCores = ref(8)
const cpuStatus = computed(() => {
  const usage = stats.value[0].value
  if (usage > 80) return { text: '高负载', class: 'bg-red-500/10 text-red-400', dotClass: 'bg-red-500' }
  if (usage > 50) return { text: '中等', class: 'bg-yellow-500/10 text-yellow-400', dotClass: 'bg-yellow-500' }
  return { text: '正常', class: 'bg-green-500/10 text-green-400', dotClass: 'bg-green-500' }
})

// 内存相关
const memoryHistory = ref<number[]>([])
const totalMemory = ref(16 * 1024 * 1024 * 1024) // 16GB
const memoryStatus = computed(() => {
  const usage = stats.value[1].value
  if (usage > 80) return { text: '警告', class: 'bg-red-500/10 text-red-400', dotClass: 'bg-red-500' }
  if (usage > 60) return { text: '中等', class: 'bg-yellow-500/10 text-yellow-400', dotClass: 'bg-yellow-500' }
  return { text: '正常', class: 'bg-green-500/10 text-green-400', dotClass: 'bg-green-500' }
})

// 网络相关
const networkHistory = ref({ down: [], up: [] })
const networkDown = ref(0)
const networkUp = ref(0)
const networkSeries = computed(() => [
  { name: '下载', data: networkHistory.value.down },
  { name: '上传', data: networkHistory.value.up }
])

// 磁盘信息
const disks = ref<any[]>([])

// 服务状态
const services = ref([
  { name: 'Web 服务', desc: 'Nginx', status: '运行中', statusClass: 'bg-green-500/10 text-green-400', icon: GlobeAltIcon, iconBg: 'bg-blue-500/10', iconColor: 'text-blue-400' },
  { name: '网络服务', desc: '千兆以太网', status: '已连接', statusClass: 'bg-green-500/10 text-green-400', icon: SignalIcon, iconBg: 'bg-cyan-500/10', iconColor: 'text-cyan-400' },
  { name: '存储服务', desc: 'ZFS Pool', status: '在线', statusClass: 'bg-green-500/10 text-green-400', icon: HardDriveIcon, iconBg: 'bg-purple-500/10', iconColor: 'text-purple-400' },
  { name: '文件服务', desc: 'SMB/CIFS', status: '运行中', statusClass: 'bg-green-500/10 text-green-400', icon: FolderIcon, iconBg: 'bg-indigo-500/10', iconColor: 'text-indigo-400' }
])

// 快速操作
const quickActions = ref([
  { name: '系统重启', desc: '重启 NAS 系统', icon: PowerIcon, iconBg: 'bg-red-500/10', iconColor: 'text-red-400', action: 'reboot' },
  { name: '刷新缓存', desc: '清理系统缓存', icon: ArrowPathIcon, iconBg: 'bg-green-500/10', iconColor: 'text-green-400', action: 'refresh' },
  { name: '下载中心', desc: '查看下载任务', icon: DocumentArrowDownIcon, iconBg: 'bg-blue-500/10', iconColor: 'text-blue-400', action: 'downloads' },
  { name: '系统设置', desc: '配置系统参数', icon: Cog6ToothIcon, iconBg: 'bg-gray-500/10', iconColor: 'text-gray-400', action: 'settings' }
])

// 图表时间标签
const chartTimeLabels = ref<string[]>([])

// CPU 图表配置
const cpuChartOptions = computed(() => ({
  chart: {
    animations: {
      enabled: true,
      easing: 'easeinout' as const,
      speed: 800
    }
  },
  colors: ['#818cf8'],
  fill: {
    type: 'gradient' as const,
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.7,
      opacityTo: 0.1,
      stops: [0, 100]
    }
  },
  stroke: {
    curve: 'smooth' as const,
    width: 2
  },
  xaxis: {
    categories: chartTimeLabels.value,
    labels: {
      style: { colors: '#9ca3af' },
      rotate: -45
    },
    axisBorder: { show: false },
    axisTicks: { show: false }
  },
  yaxis: {
    max: 100,
    labels: {
      style: { colors: '#9ca3af' },
      formatter: (val: number) => val + '%'
    }
  },
  grid: {
    borderColor: '#374151',
    strokeDashArray: 4
  },
  tooltip: {
    theme: 'dark' as const,
    y: {
      formatter: (val: number) => val + '%'
    }
  }
}))

// 内存图表配置
const memoryChartOptions = computed(() => ({
  chart: {
    animations: {
      enabled: true,
      easing: 'easeinout' as const,
      speed: 800
    }
  },
  colors: ['#c084fc'],
  fill: {
    type: 'gradient' as const,
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.7,
      opacityTo: 0.1,
      stops: [0, 100]
    }
  },
  stroke: {
    curve: 'smooth' as const,
    width: 2
  },
  xaxis: {
    categories: chartTimeLabels.value,
    labels: {
      style: { colors: '#9ca3af' },
      rotate: -45
    },
    axisBorder: { show: false },
    axisTicks: { show: false }
  },
  yaxis: {
    max: 100,
    labels: {
      style: { colors: '#9ca3af' },
      formatter: (val: number) => val + '%'
    }
  },
  grid: {
    borderColor: '#374151',
    strokeDashArray: 4
  },
  tooltip: {
    theme: 'dark' as const,
    y: {
      formatter: (val: number) => val + '%'
    }
  }
}))

// 网络图表配置
const networkChartOptions = computed(() => ({
  chart: {
    animations: {
      enabled: true,
      easing: 'easeinout' as const,
      speed: 800
    }
  },
  colors: ['#22d3ee', '#f472b6'],
  fill: {
    type: 'gradient' as const,
    gradient: {
      shadeIntensity: 1,
      opacityFrom: 0.5,
      opacityTo: 0.1,
      stops: [0, 100]
    }
  },
  stroke: {
    curve: 'smooth' as const,
    width: 2
  },
  xaxis: {
    categories: chartTimeLabels.value,
    labels: {
      style: { colors: '#9ca3af' },
      rotate: -45
    },
    axisBorder: { show: false },
    axisTicks: { show: false }
  },
  yaxis: {
    labels: {
      style: { colors: '#9ca3af' },
      formatter: (val: number) => formatBytes(val)
    }
  },
  grid: {
    borderColor: '#374151',
    strokeDashArray: 4
  },
  tooltip: {
    theme: 'dark' as const,
    y: {
      formatter: (val: number) => formatBytes(val) + '/s'
    }
  },
  legend: {
    position: 'top' as const,
    labels: { colors: '#9ca3af' }
  }
}))

// WebSocket 连接
let ws: WebSocket | null = null
let reconnectTimer: number | null = null
let pollingTimer: number | null = null

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 格式化字节
const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i]
}

// 处理快速操作
const handleQuickAction = (action: string) => {
  switch (action) {
    case 'reboot':
      if (confirm('确定要重启系统吗？')) {
        alert('系统重启指令已发送（演示）')
      }
      break
    case 'refresh':
      alert('缓存已清理（演示）')
      break
    case 'downloads':
      router.push('/downloads')
      break
    case 'settings':
      router.push('/settings')
      break
  }
}

// 初始化图表数据
const initChartData = () => {
  const now = new Date()
  for (let i = 29; i >= 0; i--) {
    const time = new Date(now.getTime() - i * 2000)
    chartTimeLabels.value.push(time.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' }))
    cpuHistory.value.push(0)
    memoryHistory.value.push(0)
    networkHistory.value.down.push(0)
    networkHistory.value.up.push(0)
  }
}

// 更新图表数据
const updateChartData = (cpu: number, memory: number, down: number, up: number) => {
  const now = new Date()
  const timeStr = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit' })

  chartTimeLabels.value.push(timeStr)
  chartTimeLabels.value.shift()

  cpuHistory.value.push(cpu)
  cpuHistory.value.shift()

  memoryHistory.value.push(memory)
  memoryHistory.value.shift()

  networkHistory.value.down.push(down)
  networkHistory.value.down.shift()

  networkHistory.value.up.push(up)
  networkHistory.value.up.shift()
}

// 处理监控数据
const handleMonitorData = (data: any) => {
  // CPU
  if (data.cpu) {
    const cpuUsage = data.cpu.usage || 0
    stats.value[0].value = Math.round(cpuUsage * 100) / 100
    cpuCores.value = data.cpu.cores || 8
  }

  // 内存
  if (data.memory) {
    const memUsage = (data.memory.used / data.memory.total) * 100 || 0
    stats.value[1].value = Math.round(memUsage * 10) / 10
    totalMemory.value = data.memory.total || 16 * 1024 * 1024 * 1024
    stats.value[1].unit = `${formatBytes(data.memory.used)} / ${formatBytes(data.memory.total)}`
  }

  // 磁盘
  if (data.disk && Array.isArray(data.disk)) {
    disks.value = data.disk.map((d: any) => ({
      ...d,
      usedPercent: Math.round((d.used / d.total) * 100)
    }))
    if (disks.value.length > 0) {
      stats.value[2].value = disks.value[0].usedPercent || 0
      stats.value[2].unit = disks.value[0].mountpoint || '系统盘'
    }
  }

  // 网络
  if (data.network) {
    // 处理 network.interfaces[] 数据结构
    if (data.network.interfaces && Array.isArray(data.network.interfaces)) {
      const interfaces = data.network.interfaces.filter((i: any) => i.name !== 'lo' && !i.name.startsWith('docker') && !i.name.startsWith('virbr') && !i.name.startsWith('veth'))
      const totalDown = interfaces.reduce((sum: number, i: any) => sum + (i.bytesRecv || 0), 0)
      const totalUp = interfaces.reduce((sum: number, i: any) => sum + (i.bytesSent || 0), 0)

      networkDown.value = totalDown
      networkUp.value = totalUp

      // 使用实时速度数据（如果可用）
      const totalSpeedUp = interfaces.reduce((sum: number, i: any) => sum + (i.sentSpeed || 0), 0)
      const totalSpeedDown = interfaces.reduce((sum: number, i: any) => sum + (i.recvSpeed || 0), 0)

      stats.value[3].value = Math.round(totalSpeedUp / 1024)
      stats.value[3].unit = `实时速率 · ${formatSpeed(totalSpeedUp)}/s`

      updateChartData(
        stats.value[0].value,
        stats.value[1].value,
        totalDown,
        totalUp
      )
    }
    // 兼容旧格式
    else if (data.network.rx_bytes || data.network.tx_bytes) {
      const down = data.network.rx_bytes || 0
      const up = data.network.tx_bytes || 0
      networkDown.value = down
      networkUp.value = up

      stats.value[3].value = Math.round(up / 1024)
      stats.value[3].unit = `实时速率 · ${formatSpeed(up)}/s`

      updateChartData(
        stats.value[0].value,
        stats.value[1].value,
        down,
        up
      )
    }
  } else if (data.network_interfaces) {
    // 处理网络接口数据
    const interfaces = data.network_interfaces.filter((i: any) => i.name !== 'lo' && !i.name.startsWith('docker') && !i.name.startsWith('virbr') && !i.name.startsWith('veth'))
    const totalDown = interfaces.reduce((sum: number, i: any) => sum + (i.bytesRecv || i.rx_bytes || 0), 0)
    const totalUp = interfaces.reduce((sum: number, i: any) => sum + (i.bytesSent || i.tx_bytes || 0), 0)

    networkDown.value = totalDown
    networkUp.value = totalUp

    stats.value[3].value = Math.round(totalUp / 1024)
    stats.value[3].unit = `实时速率 · ${formatSpeed(totalUp)}/s`

    updateChartData(
      stats.value[0].value,
      stats.value[1].value,
      totalDown,
      totalUp
    )
  }

  // 运行时间
  if (data.uptime) {
    const seconds = data.uptime
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)
    uptime.value = days > 0
      ? `${days}天 ${hours}小时 ${minutes}分钟`
      : `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
  }
}

// 格式化速度
const formatSpeed = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i]
}

// AnimatedNumber 组件（内联）
const AnimatedNumber = (props: { value: number; decimals?: number }) => {
  return h('span', props.value.toFixed(props.decimals || 0))
}

// 轮询模式（WebSocket 失败时的降级方案）
const startPolling = () => {
  const fetchMonitorData = async () => {
    try {
      const [cpu, memory, disk, network] = await Promise.all([
        monitorApi.getCPU().catch(() => null),
        monitorApi.getMemory().catch(() => null),
        monitorApi.getDisk().catch(() => null),
        monitorApi.getNetwork().catch(() => null)
      ])

      handleMonitorData({
        cpu: cpu ? { usage: cpu.usage, cores: cpu.cores } : null,
        memory: memory || null,
        disk: disk ? disk.disks || [disk] : null,
        network_interfaces: network?.interfaces || []
      })
    } catch (e) {
      console.error('Polling error:', e)
    }
  }

  fetchMonitorData()
  pollingTimer = window.setInterval(fetchMonitorData, 2000)
}

// 计算模拟运行时间
const calculateUptime = () => {
  const bootTime = new Date()
  bootTime.setHours(bootTime.getHours() - 48) // 假设系统运行了2天
  const now = new Date()
  const diff = now.getTime() - bootTime.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  uptime.value = days > 0
    ? `${days}天 ${hours}小时 ${minutes}分钟`
    : `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`
}

onMounted(() => {
  // 初始化图表
  initChartData()

  // 更新时间
  updateTime()
  setInterval(updateTime, 1000)

  // 计算运行时间
  calculateUptime()

  // 尝试连接 WebSocket
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = import.meta.env.VITE_WS_URL?.replace('ws://', '').replace('wss://', '') ||
              `${window.location.hostname}:8080`
  const wsUrl = `${protocol}//${host}/ws/monitor`

  try {
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      console.log('WebSocket connected')
    }

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data)
        handleMonitorData(data)
      } catch (e) {
        console.error('Parse error:', e)
      }
    }

    ws.onclose = () => {
      console.log('WebSocket closed, switching to polling')
      startPolling()
    }

    ws.onerror = () => {
      console.log('WebSocket error, switching to polling')
      startPolling()
    }
  } catch (e) {
    console.log('WebSocket init failed, using polling')
    startPolling()
  }
})

onUnmounted(() => {
  if (ws) {
    ws.close()
    ws = null
  }
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
  }
  if (pollingTimer) {
    clearInterval(pollingTimer)
  }
})
</script>
