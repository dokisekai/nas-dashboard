<template>
  <div class="system-monitor-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <div class="widget-title">系统监控</div>
      <div class="widget-status" :class="statusClass">
        <div class="status-dot"></div>
      </div>
    </div>

    <div class="widget-content">
      <!-- CPU 使用率 -->
      <div v-if="config.showCpu" class="monitor-item">
        <div class="monitor-label">
          <CpuChipIcon class="icon" />
          <span>CPU</span>
        </div>
        <div class="monitor-value">{{ cpuUsage }}%</div>
        <div class="monitor-bar">
          <div class="monitor-fill" :style="{ width: cpuUsage + '%' }"></div>
        </div>
      </div>

      <!-- 内存使用率 -->
      <div v-if="config.showMemory" class="monitor-item">
        <div class="monitor-label">
          <ServerIcon class="icon" />
          <span>内存</span>
        </div>
        <div class="monitor-value">{{ memoryUsage }}%</div>
        <div class="monitor-bar">
          <div class="monitor-fill memory" :style="{ width: memoryUsage + '%' }"></div>
        </div>
      </div>

      <!-- 磁盘使用率 -->
      <div v-if="config.showDisk" class="monitor-item">
        <div class="monitor-label">
          <CircleStackIcon class="icon" />
          <span>磁盘</span>
        </div>
        <div class="monitor-value">{{ diskUsage }}%</div>
        <div class="monitor-bar">
          <div class="monitor-fill disk" :style="{ width: diskUsage + '%' }"></div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { CpuChipIcon, ServerIcon, CircleStackIcon } from '@heroicons/vue/24/outline'

interface Props {
  config: {
    showCpu?: boolean
    showMemory?: boolean
    showDisk?: boolean
  }
  size: 'small' | 'medium' | 'large'
}

const props = defineProps<Props>()

// 模拟数据
const cpuUsage = ref(0)
const memoryUsage = ref(0)
const diskUsage = ref(0)

const statusClass = computed(() => {
  const avgUsage = (cpuUsage.value + memoryUsage.value + diskUsage.value) / 3
  if (avgUsage > 80) return 'status-error'
  if (avgUsage > 60) return 'status-warning'
  return 'status-normal'
})

// 模拟实时数据更新
let updateInterval: number | null = null

const updateStats = () => {
  cpuUsage.value = Math.round(Math.random() * 40 + 20)
  memoryUsage.value = Math.round(Math.random() * 30 + 40)
  diskUsage.value = Math.round(Math.random() * 20 + 50)
}

onMounted(() => {
  updateStats()
  updateInterval = window.setInterval(updateStats, 2000)
})

onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval)
  }
})
</script>

<style scoped>
.system-monitor-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.widget-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.widget-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.widget-status {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.status-normal .status-dot {
  background: #10b981;
}

.status-warning .status-dot {
  background: #f59e0b;
}

.status-error .status-dot {
  background: #ef4444;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.widget-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.monitor-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.monitor-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6b7280;
}

.icon {
  width: 16px;
  height: 16px;
}

.monitor-value {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
}

.widget-small .monitor-value {
  font-size: 18px;
}

.monitor-bar {
  height: 4px;
  background: #f3f4f6;
  border-radius: 2px;
  overflow: hidden;
}

.monitor-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 2px;
  transition: width 0.5s ease;
}

.monitor-fill.memory {
  background: linear-gradient(90deg, #10b981 0%, #34d399 100%);
}

.monitor-fill.disk {
  background: linear-gradient(90deg, #f59e0b 0%, #fbbf24 100%);
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .widget-title,
  .monitor-value {
    color: #f9fafb;
  }

  .monitor-label {
    color: #9ca3af;
  }

  .monitor-bar {
    background: rgba(255, 255, 255, 0.1);
  }
}
</style>