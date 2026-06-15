<template>
  <div class="system-monitor">
    <div class="monitor-header">
      <h1>系统监控器</h1>
      <p class="subtitle">实时系统资源监控</p>
    </div>

    <!-- Stats Overview -->
    <div class="stats-grid">
      <div class="stat-card cpu">
        <div class="stat-icon">
          <CpuChipIcon class="w-8 h-8" />
        </div>
        <div class="stat-content">
          <h3>CPU</h3>
          <p class="stat-value">{{ cpuUsage }}%</p>
          <p class="stat-detail">{{ cpuCores }} 核心</p>
        </div>
        <div class="stat-chart">
          <sparkline :data="cpuHistory" :color="cpuColor" />
        </div>
      </div>

      <div class="stat-card memory">
        <div class="stat-icon">
          <ServerIcon class="w-8 h-8" />
        </div>
        <div class="stat-content">
          <h3>内存</h3>
          <p class="stat-value">{{ memoryUsage }}%</p>
          <p class="stat-detail">{{ formatBytes(memoryUsed) }} / {{ formatBytes(memoryTotal) }}</p>
        </div>
        <div class="stat-chart">
          <sparkline :data="memoryHistory" :color="memoryColor" />
        </div>
      </div>

      <div class="stat-card disk">
        <div class="stat-icon">
          <CircleStackIcon class="w-8 h-8" />
        </div>
        <div class="stat-content">
          <h3>磁盘</h3>
          <p class="stat-value">{{ diskUsage }}%</p>
          <p class="stat-detail">{{ formatBytes(diskUsed) }} / {{ formatBytes(diskTotal) }}</p>
        </div>
        <div class="stat-chart">
          <sparkline :data="diskHistory" :color="diskColor" />
        </div>
      </div>

      <div class="stat-card network">
        <div class="stat-icon">
          <GlobeAltIcon class="w-8 h-8" />
        </div>
        <div class="stat-content">
          <h3>网络</h3>
          <p class="stat-value">{{ formatBytes(networkSpeed) }}/s</p>
          <p class="stat-detail">↑ {{ formatBytes(networkUpload) }} ↓ {{ formatBytes(networkDownload) }}</p>
        </div>
        <div class="stat-chart">
          <sparkline :data="networkHistory" :color="networkColor" />
        </div>
      </div>

      <div class="stat-card power">
        <div class="stat-icon">
          <BoltIcon class="w-8 h-8" />
        </div>
        <div class="stat-content">
          <h3>功耗</h3>
          <p class="stat-value">{{ powerUsage }} W</p>
          <p class="stat-detail">整机实时功耗</p>
        </div>
        <div class="stat-chart">
          <sparkline :data="powerHistory" :color="powerColor" />
        </div>
      </div>

      <!-- UPS Card -->
      <div v-if="upsStatus" class="stat-card ups">
        <div class="stat-icon" :class="{ 'on-battery': upsStatus.status !== 'OL' }">
          <Battery50Icon class="w-8 h-8" />
        </div>
        <div class="stat-content">
          <h3>UPS</h3>
          <p class="stat-value">{{ upsStatus.charge }}%</p>
          <p class="stat-detail">{{ upsStatus.status === 'OL' ? '市电供电' : '电池供电' }} ({{ upsStatus.load }}% 负载)</p>
        </div>
        <div class="stat-chart">
          <sparkline :data="upsHistory" :color="upsColor" />
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.label }}
      </button>
    </div>

    <!-- Detailed Charts -->
    <div class="chart-container">
      <!-- CPU Chart -->
      <div v-if="activeTab === 'cpu'" class="chart-wrapper">
        <div class="chart-header">
          <h2>CPU 使用率</h2>
          <div class="chart-controls">
            <button
              v-for="range in timeRanges"
              :key="range.value"
              :class="['time-btn', { active: selectedRange === range.value }]"
              @click="selectedRange = range.value"
            >
              {{ range.label }}
            </button>
          </div>
        </div>
        <div class="chart">
          <Line
            :data="cpuChartData"
            :options="chartOptions"
            :height="300"
          />
        </div>
        <div class="cpu-details">
          <div class="detail-item">
            <span class="detail-label">负载:</span>
            <span class="detail-value">{{ loadAverage.join(', ') }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">温度:</span>
            <span class="detail-value">{{ cpuTemperature }}°C</span>
          </div>
        </div>
      </div>

      <!-- Memory Chart -->
      <div v-if="activeTab === 'memory'" class="chart-wrapper">
        <div class="chart-header">
          <h2>内存使用情况</h2>
        </div>
        <div class="chart">
          <Line
            :data="memoryChartData"
            :options="memoryChartOptions"
            :height="300"
          />
        </div>
        <div class="memory-breakdown">
          <div class="breakdown-item">
            <div class="breakdown-bar">
              <div class="bar-segment used" :style="{ width: memoryUsage + '%' }"></div>
              <div class="bar-segment cached" :style="{ width: cachedPercent + '%' }"></div>
              <div class="bar-segment buffers" :style="{ width: buffersPercent + '%' }"></div>
            </div>
            <div class="breakdown-legend">
              <span class="legend-item used">已用: {{ memoryUsage }}%</span>
              <span class="legend-item cached">缓存: {{ cachedPercent }}%</span>
              <span class="legend-item buffers">缓冲: {{ buffersPercent }}%</span>
              <span class="legend-item free">空闲: {{ freePercent }}%</span>
            </div>
          </div>
          <div class="memory-stats">
            <div class="stat-row">
              <span>总内存:</span>
              <span>{{ formatBytes(memoryTotal) }}</span>
            </div>
            <div class="stat-row">
              <span>已用:</span>
              <span>{{ formatBytes(memoryUsed) }}</span>
            </div>
            <div class="stat-row">
              <span>缓存:</span>
              <span>{{ formatBytes(memoryCached) }}</span>
            </div>
            <div class="stat-row">
              <span>可用:</span>
              <span>{{ formatBytes(memoryAvailable) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Disk I/O Chart -->
      <div v-if="activeTab === 'disk'" class="chart-wrapper">
        <div class="chart-header">
          <h2>磁盘 I/O</h2>
        </div>
        <div class="chart">
          <Line
            :data="diskChartData"
            :options="diskChartOptions"
            :height="300"
          />
        </div>
        <div class="disk-details">
          <div class="detail-item">
            <span class="detail-label">读取:</span>
            <span class="detail-value">{{ formatBytes(diskReadRate) }}/s</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">写入:</span>
            <span class="detail-value">{{ formatBytes(diskWriteRate) }}/s</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">IOPS:</span>
            <span class="detail-value">{{ diskIOPS }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">延迟:</span>
            <span class="detail-value">{{ diskLatency }}ms</span>
          </div>
        </div>
        <div class="disk-list">
          <div
            v-for="disk in diskList"
            :key="disk.name"
            class="disk-item"
          >
            <div class="disk-name">{{ disk.name }}</div>
            <div class="disk-usage">
              <div class="usage-bar">
                <div class="usage-fill" :style="{ width: disk.usage + '%' }"></div>
              </div>
              <span class="usage-text">{{ disk.usage }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Network Chart -->
      <div v-if="activeTab === 'network'" class="chart-wrapper">
        <div class="chart-header">
          <h2>网络流量</h2>
        </div>
        <div class="chart">
          <Line
            :data="networkChartData"
            :options="networkChartOptions"
            :height="300"
          />
        </div>
        <div class="network-details">
          <div class="detail-item">
            <span class="detail-label">上传:</span>
            <span class="detail-value">{{ formatBytes(networkUpload) }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">下载:</span>
            <span class="detail-value">{{ formatBytes(networkDownload) }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">上传速度:</span>
            <span class="detail-value">{{ formatBytes(networkUploadRate) }}/s</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">下载速度:</span>
            <span class="detail-value">{{ formatBytes(networkDownloadRate) }}/s</span>
          </div>
        </div>
        <div class="network-interfaces">
          <div
            v-for="iface in networkInterfaces"
            :key="iface.name"
            class="interface-item"
          >
            <div class="interface-info">
              <h4>{{ iface.name }}</h4>
              <p>{{ iface.ip }} - {{ iface.status }}</p>
            </div>
            <div class="interface-stats">
              <div class="stat">
                <span>↑ {{ formatBytes(iface.upload) }}</span>
              </div>
              <div class="stat">
                <span>↓ {{ formatBytes(iface.download) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Power Chart -->
      <div v-if="activeTab === 'power'" class="chart-wrapper">
        <div class="chart-header">
          <h2>功耗监控</h2>
          <div class="chart-controls">
            <button
              v-for="range in timeRanges"
              :key="range.value"
              :class="['time-btn', { active: selectedRange === range.value }]"
              @click="selectedRange = range.value"
            >
              {{ range.label }}
            </button>
          </div>
        </div>
        <div class="chart">
          <Line
            :data="powerChartData"
            :options="powerChartOptions"
            :height="300"
          />
        </div>
        <div class="power-details">
          <div class="detail-item">
            <span class="detail-label">CPU Package:</span>
            <span class="detail-value">{{ cpuPackagePower }} W</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">Intel 核显:</span>
            <span class="detail-value">{{ igpuPower }} W</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">AMD 独显:</span>
            <span class="detail-value">{{ dgpuPower }} W</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">存储设备:</span>
            <span class="detail-value">{{ storagePower }} W</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">其他组件:</span>
            <span class="detail-value">{{ otherComponentsPower }} W</span>
          </div>
        </div>
        <div class="power-breakdown">
          <div class="breakdown-title">功耗组成</div>
          <div class="breakdown-item">
            <div class="breakdown-bar">
              <div class="bar-segment cpu" :style="{ width: ((cpuPackagePower / powerUsage) * 100) + '%' }"></div>
              <div class="bar-segment gpu" :style="{ width: ((igpuPower / powerUsage) * 100) + '%' }"></div>
              <div class="bar-segment dgpu" :style="{ width: ((dgpuPower / powerUsage) * 100) + '%' }"></div>
              <div class="bar-segment storage" :style="{ width: ((storagePower / powerUsage) * 100) + '%' }"></div>
              <div class="bar-segment other" :style="{ width: ((otherComponentsPower / powerUsage) * 100) + '%' }"></div>
            </div>
            <div class="breakdown-legend">
              <span class="legend-item cpu">CPU: {{ cpuPackagePower }}W</span>
              <span class="legend-item gpu">核显: {{ igpuPower }}W</span>
              <span class="legend-item dgpu">独显: {{ dgpuPower }}W</span>
              <span class="legend-item storage">存储: {{ storagePower }}W</span>
              <span class="legend-item other">其他: {{ otherComponentsPower }}W</span>
            </div>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Line } from 'vue-chartjs'
import Sparkline from '../components/Sparkline.vue'
import { monitorApi, systemApi } from '../api'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'
import {
  CpuChipIcon,
  ServerIcon,
  CircleStackIcon,
  GlobeAltIcon,
  BoltIcon,
  Battery50Icon
} from '@heroicons/vue/24/outline'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const activeTab = ref('cpu')
const selectedRange = ref('1h')

const tabs = [
  { id: 'cpu', label: 'CPU', icon: CpuChipIcon },
  { id: 'memory', label: '内存', icon: ServerIcon },
  { id: 'disk', label: '磁盘', icon: CircleStackIcon },
  { id: 'network', label: '网络', icon: GlobeAltIcon },
  { id: 'power', label: '功耗', icon: BoltIcon }
]

const timeRanges = [
  { label: '1小时', value: '1h' },
  { label: '6小时', value: '6h' },
  { label: '24小时', value: '24h' },
  { label: '7天', value: '7d' }
]

// Stats
const cpuUsage = ref(0)
const cpuCores = ref(0)
const cpuTemperature = ref(0)
const cpuHistory = ref<number[]>([])
const cpuColor = '#3b82f6'

const memoryUsage = ref(0)
const memoryTotal = ref(0)
const memoryUsed = ref(0)
const memoryCached = ref(0)
const memoryAvailable = ref(0)
const memoryHistory = ref<number[]>([])
const memoryColor = '#10b981'

const diskUsage = ref(0)
const diskTotal = ref(0)
const diskUsed = ref(0)
const diskReadRate = ref(0)
const diskWriteRate = ref(0)
const diskIOPS = ref(0)
const diskLatency = ref(0)
const diskHistory = ref<number[]>([])
const diskColor = '#f59e0b'

const networkSpeed = ref(0)
const networkUpload = ref(0)
const networkDownload = ref(0)
const networkUploadRate = ref(0)
const networkDownloadRate = ref(0)
const networkHistory = ref<number[]>([])
const networkColor = '#8b5cf6'

// Power stats
const powerUsage = ref(0)
const cpuPackagePower = ref(0)
const igpuPower = ref(0)
const dgpuPower = ref(0)
const storagePower = ref(0)
const otherComponentsPower = ref(0)
const powerHistory = ref<number[]>([])
const powerColor = '#ef4444'

// UPS stats
const upsStatus = ref<any>(null)
const upsHistory = ref<number[]>([])
const upsColor = '#14b8a6'

const loadAverage = ref<number[]>([0, 0, 0])

const cachedPercent = computed(() => Math.round((memoryCached.value / memoryTotal.value) * 100))
const buffersPercent = computed(() => 5) // Mock
const freePercent = computed(() => 100 - memoryUsage.value - cachedPercent.value - buffersPercent.value)

const diskList = ref<any[]>([])
const networkInterfaces = ref<any[]>([])

// Chart configurations
const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      max: 100,
      ticks: {
        callback: (value: string | number) => String(value) + '%'
      }
    }
  },
  elements: {
    line: {
      tension: 0.4
    },
    point: {
      radius: 0,
      hitRadius: 10,
      hoverRadius: 5
    }
  }
}

const cpuChartData = computed(() => ({
  labels: generateLabels(),
  datasets: [{
    label: 'CPU 使用率',
    data: cpuHistory.value,
    borderColor: cpuColor,
    backgroundColor: 'rgba(59, 130, 246, 0.1)',
    fill: true
  }]
}))

const memoryChartOptions = {
  ...chartOptions,
  scales: {
    ...chartOptions.scales,
    y: {
      beginAtZero: true,
      max: 100,
      ticks: {
        callback: (value: string | number) => String(value) + '%'
      }
    }
  }
}

const memoryChartData = computed(() => ({
  labels: generateLabels(),
  datasets: [
    {
      label: '已用',
      data: memoryHistory.value,
      borderColor: '#10b981',
      backgroundColor: 'rgba(16, 185, 129, 0.1)',
      fill: true
    }
  ]
}))

const diskChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top' as any
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        callback: (value: string | number) => formatBytes(Number(value))
      }
    }
  }
}

const diskChartData = computed(() => ({
  labels: generateLabels(),
  datasets: [
    {
      label: '读取',
      data: Array(60).fill(0).map(() => Math.random() * 100000000),
      borderColor: '#f59e0b',
      backgroundColor: 'transparent'
    },
    {
      label: '写入',
      data: Array(60).fill(0).map(() => Math.random() * 100000000),
      borderColor: '#ef4444',
      backgroundColor: 'transparent'
    }
  ]
}))

const networkChartOptions = {
  ...diskChartOptions
}

const networkChartData = computed(() => ({
  labels: generateLabels(),
  datasets: [
    {
      label: '上传',
      data: Array(60).fill(0).map(() => Math.random() * 10000000),
      borderColor: '#8b5cf6',
      backgroundColor: 'transparent'
    },
    {
      label: '下载',
      data: Array(60).fill(0).map(() => Math.random() * 10000000),
      borderColor: '#3b82f6',
      backgroundColor: 'transparent'
    }
  ]
}))

const powerChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: true,
      position: 'top' as any
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: {
        callback: (value: string | number) => String(value) + ' W'
      }
    }
  }
}

const powerChartData = computed(() => ({
  labels: generateLabels(),
  datasets: [{
    label: '整机功耗 (W)',
    data: powerHistory.value,
    borderColor: '#ef4444',
    backgroundColor: 'rgba(239, 68, 68, 0.1)',
    fill: true
  }]
}))

function generateLabels() {
  const labels = []
  for (let i = 59; i >= 0; i--) {
    labels.push(`${i}s ago`)
  }
  return labels
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

// Data fetching
let updateInterval: number

const updateStats = async () => {
  try {
    // Fetch CPU stats
    const cpuResponse = await monitorApi.getCPU()
    const cpuData = cpuResponse  // axios拦截器已返回response.data
    cpuUsage.value = Math.round(cpuData.usage * 100)
    cpuCores.value = cpuData.cores
    cpuTemperature.value = Math.round(cpuData.usage * 100) // Using usage as temp proxy
    cpuHistory.value = [...cpuHistory.value.slice(-59), cpuUsage.value]

    // Fetch memory stats
    const memoryResponse = await monitorApi.getMemory()
    const memoryData = memoryResponse  // axios拦截器已返回response.data
    memoryUsage.value = Math.round(memoryData.percent)
    memoryTotal.value = memoryData.total
    memoryUsed.value = memoryData.used
    memoryCached.value = memoryData.cached
    memoryAvailable.value = memoryData.available
    memoryHistory.value = [...memoryHistory.value.slice(-59), memoryUsage.value]

    // Fetch disk stats
    const diskResponse = await monitorApi.getDisk()
    const diskData = diskResponse  // axios拦截器已返回response.data
    if (diskData.usage && diskData.usage.length > 0) {
      const mainDisk = diskData.usage[0]
      diskUsage.value = Math.round(mainDisk.percent)
      diskTotal.value = mainDisk.total
      diskUsed.value = mainDisk.used
      diskHistory.value = [...diskHistory.value.slice(-59), diskUsage.value]

      diskList.value = diskData.usage.map((disk: any) => ({
        name: disk.device,
        usage: Math.round(disk.percent)
      }))
    }

    // Fetch network stats
    const networkResponse = await monitorApi.getNetwork()
    const networkData = networkResponse  // axios拦截器已返回response.data
    if (networkData.interfaces && networkData.interfaces.length > 0) {
      const mainInterface = networkData.interfaces[0]
      networkSpeed.value = Math.round((mainInterface.rx_bytes + mainInterface.tx_bytes) / 1024)
      networkUpload.value += mainInterface.tx_bytes
      networkDownload.value += mainInterface.rx_bytes
      networkUploadRate.value = Math.round(mainInterface.tx_bytes / 1024)
      networkDownloadRate.value = Math.round(mainInterface.rx_bytes / 1024)
      networkHistory.value = [...networkHistory.value.slice(-59), networkSpeed.value]

      networkInterfaces.value = networkData.interfaces.map((iface: any) => ({
        name: iface.name,
        ip: iface.address || 'N/A',
        status: iface.operstate === 'up' ? 'up' : 'down',
        upload: iface.tx_bytes || 0,
        download: iface.rx_bytes || 0
      }))
    }

    // Calculate load averages from CPU data
    if (cpuData.load1) {
      loadAverage.value = [
        parseFloat(cpuData.load1.toFixed(2)),
        parseFloat(cpuData.load5.toFixed(2)),
        parseFloat(cpuData.load15.toFixed(2))
      ]
    }

    // Fetch power stats
    try {
      const powerResponse = await fetch('/api/power/current', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      if (powerResponse.ok) {
        const powerData = await powerResponse.json()
        powerUsage.value = Math.round(powerData.total * 10) / 10
        cpuPackagePower.value = Math.round(powerData.cpuPackage * 10) / 10
        igpuPower.value = Math.round(powerData.igpu * 10) / 10
        dgpuPower.value = Math.round(powerData.dgpu * 10) / 10
        storagePower.value = Math.round((powerData.hdd + powerData.ssd) * 10) / 10
        otherComponentsPower.value = Math.round((powerData.mbram + powerData.cooling + powerData.usb + powerData.powerLoss) * 10) / 10
        powerHistory.value = [...powerHistory.value.slice(-59), powerUsage.value]
      }
    } catch (powerError) {
      console.error('Failed to fetch power data:', powerError)
      // Use default values if API call fails
      powerUsage.value = 60
      cpuPackagePower.value = 3
      igpuPower.value = 0.5
      dgpuPower.value = 8
      storagePower.value = 6
      otherComponentsPower.value = 42.5
      powerHistory.value = [...powerHistory.value.slice(-59), powerUsage.value]
    }

    // Fetch UPS stats
    try {
      const upsData = await systemApi.getUPSStatus()
      if (upsData && upsData.status !== 'Unknown') {
        upsStatus.value = upsData
        upsHistory.value = [...upsHistory.value.slice(-59), upsData.charge]
      } else {
        upsStatus.value = null
      }
    } catch (upsError) {
      console.warn('UPS monitoring not available:', upsError)
      upsStatus.value = null
    }

  } catch (error) {
    console.error('Failed to update stats:', error)
  }
}

onMounted(() => {
  updateStats()
  updateInterval = setInterval(updateStats, 1000) as unknown as number
})

onUnmounted(() => {
  clearInterval(updateInterval)
})
</script>

<style scoped>
.system-monitor {
  width: 100%;
  height: 100%;
  padding: 32px;
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.monitor-header {
  margin-bottom: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.monitor-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 16px;
  align-items: center;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-card.cpu .stat-icon {
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
}

.stat-card.memory .stat-icon {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
}

.stat-card.disk .stat-icon {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
}

.stat-card.network .stat-icon {
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
}

.stat-card.power .stat-icon {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
}

.stat-card.ups .stat-icon {
  background: linear-gradient(135deg, #14b8a6 0%, #2dd4bf 100%);
}

.stat-icon.on-battery {
  background: #ef4444 !important;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.1); }
  100% { transform: scale(1); }
}

.stat-content h3 {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 4px;
}

.stat-detail {
  font-size: 12px;
  color: #9ca3af;
}

.stat-chart {
  width: 80px;
  height: 40px;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 2px solid rgba(102, 126, 234, 0.1);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: -2px;
  border-radius: 8px 8px 0 0;
}

.tab-btn:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.tab-btn.active {
  color: white;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-bottom-color: transparent;
}

.chart-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  margin-bottom: 24px;
}

.chart-wrapper {
  width: 100%;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.chart-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.chart-controls {
  display: flex;
  gap: 8px;
}

.time-btn {
  padding: 6px 12px;
  background: rgba(102, 126, 234, 0.1);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 6px;
  font-size: 12px;
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.time-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.3);
}

.time-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
}

.chart {
  margin-bottom: 24px;
}

.cpu-details,
.disk-details,
.network-details {
  display: flex;
  gap: 32px;
  margin-bottom: 24px;
}

.detail-item {
  display: flex;
  gap: 8px;
  font-size: 14px;
}

.detail-label {
  color: #6b7280;
  min-width: 60px;
}

.detail-value {
  color: #1f2937;
  font-weight: 500;
}

.memory-breakdown {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.breakdown-bar {
  height: 24px;
  background: #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
}

.bar-segment {
  height: 100%;
  transition: width 0.3s ease;
}

.bar-segment.used {
  background: #ef4444;
}

.bar-segment.cached {
  background: #f59e0b;
}

.bar-segment.buffers {
  background: #10b981;
}

.breakdown-legend {
  display: flex;
  gap: 16px;
  font-size: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.legend-item.used::before {
  content: '';
  width: 8px;
  height: 8px;
  background: #ef4444;
  border-radius: 2px;
}

.legend-item.cached::before {
  content: '';
  width: 8px;
  height: 8px;
  background: #f59e0b;
  border-radius: 2px;
}

.legend-item.buffers::before {
  content: '';
  width: 8px;
  height: 8px;
  background: #10b981;
  border-radius: 2px;
}

.legend-item.free::before {
  content: '';
  width: 8px;
  height: 8px;
  background: #e5e7eb;
  border-radius: 2px;
}

.memory-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.stat-row:first-child {
  color: #1f2937;
  font-weight: 500;
}

.disk-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.disk-item {
  display: flex;
  gap: 16px;
  align-items: center;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.disk-name {
  min-width: 100px;
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.disk-usage {
  flex: 1;
  display: flex;
  gap: 12px;
  align-items: center;
}

.usage-bar {
  flex: 1;
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
}

.usage-fill {
  height: 100%;
  background: linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%);
  transition: width 0.3s ease;
}

.usage-text {
  min-width: 40px;
  font-size: 12px;
  color: #6b7280;
  text-align: right;
}

.network-interfaces {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 24px;
}

.interface-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.interface-info h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.interface-info p {
  font-size: 12px;
  color: #6b7280;
}

.interface-stats {
  display: flex;
  gap: 16px;
}

.interface-stats .stat {
  font-size: 12px;
  color: #6b7280;
  min-width: 80px;
  text-align: right;
}

.power-breakdown {
  margin-top: 24px;
}

.breakdown-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.breakdown-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.breakdown-bar {
  height: 24px;
  background: #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
}

.bar-segment.cpu {
  background: #3b82f6;
}

.bar-segment.gpu {
  background: #10b981;
}

.bar-segment.dgpu {
  background: #f59e0b;
}

.bar-segment.storage {
  background: #8b5cf6;
}

.bar-segment.other {
  background: #6b7280;
}

.breakdown-legend {
  display: flex;
  gap: 16px;
  font-size: 12px;
  flex-wrap: wrap;
}

.legend-item.cpu::before {
  background: #3b82f6;
}

.legend-item.gpu::before {
  background: #10b981;
}

.legend-item.dgpu::before {
  background: #f59e0b;
}

.legend-item.storage::before {
  background: #8b5cf6;
}

.legend-item.other::before {
  background: #6b7280;
}
</style>