<template>
  <div class="network-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <SignalIcon class="w-5 h-5" />
      <span class="widget-title">网络监控</span>
    </div>

    <div class="network-content">
      <!-- 连接状态 -->
      <div class="connection-status">
        <div class="status-indicator" :class="connectionClass"></div>
        <span class="status-text">{{ connectionStatus }}</span>
      </div>

      <!-- 网络统计 -->
      <div class="network-stats">
        <div class="stat-item">
          <ArrowDownIcon class="stat-icon download" />
          <div class="stat-info">
            <span class="stat-label">下载</span>
            <span class="stat-value">{{ formatBytes(downloadSpeed) }}/s</span>
          </div>
        </div>
        <div class="stat-item">
          <ArrowUpIcon class="stat-icon upload" />
          <div class="stat-info">
            <span class="stat-label">上传</span>
            <span class="stat-value">{{ formatBytes(uploadSpeed) }}/s</span>
          </div>
        </div>
      </div>

      <!-- 网络图表 -->
      <div v-if="config.showGraph" class="network-graph">
        <canvas ref="chartCanvas" class="chart-canvas"></canvas>
      </div>

      <!-- 接口信息 -->
      <div class="interface-info">
        <div class="interface-item">
          <GlobeAltIcon class="interface-icon" />
          <span class="interface-name">{{ ipAddress }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { SignalIcon, ArrowDownIcon, ArrowUpIcon, GlobeAltIcon } from '@heroicons/vue/24/outline'

interface Props {
  config: {
    interface?: string
    showGraph?: boolean
    showSpeed?: boolean
  }
  size: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    interface: 'all',
    showGraph: true,
    showSpeed: true
  })
})

const chartCanvas = ref<HTMLCanvasElement>()
const connectionStatus = ref('已连接')
const downloadSpeed = ref(0)
const uploadSpeed = ref(0)
const ipAddress = ref('192.168.50.10')

const connectionClass = computed(() => {
  return connectionStatus.value === '已连接' ? 'connected' : 'disconnected'
})

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i]
}

// 模拟网络数据更新
const updateNetworkData = () => {
  downloadSpeed.value = Math.random() * 50 * 1024 * 1024 // 0-50MB/s
  uploadSpeed.value = Math.random() * 20 * 1024 * 1024 // 0-20MB/s
}

let updateInterval: number | null = null

onMounted(() => {
  updateNetworkData()
  updateInterval = window.setInterval(updateNetworkData, 2000)
})

onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval)
  }
})
</script>

<style scoped>
.network-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.widget-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  margin-bottom: 12px;
}

.widget-header svg {
  color: #6b7280;
}

.widget-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.network-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.connection-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-indicator.connected {
  background: #10b981;
  animation: pulse 2s infinite;
}

.status-indicator.disconnected {
  background: #ef4444;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-text {
  font-size: 12px;
  color: #6b7280;
}

.network-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-icon {
  width: 16px;
  height: 16px;
}

.stat-icon.download {
  color: #22d3ee;
}

.stat-icon.upload {
  color: #f472b6;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 11px;
  color: #6b7280;
}

.stat-value {
  font-size: 13px;
  font-weight: 600;
  color: #1f2937;
}

.network-graph {
  height: 60px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 8px;
  overflow: hidden;
}

.chart-canvas {
  width: 100%;
  height: 100%;
}

.interface-info {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 8px;
}

.interface-icon {
  width: 14px;
  height: 14px;
  color: #6b7280;
}

.interface-name {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .network-widget {
    background: rgba(0, 0, 0, 0.8);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .widget-title,
  .stat-value,
  .interface-name {
    color: #f9fafb;
  }

  .network-graph,
  .interface-info {
    background: rgba(255, 255, 255, 0.05);
  }
}
</style>