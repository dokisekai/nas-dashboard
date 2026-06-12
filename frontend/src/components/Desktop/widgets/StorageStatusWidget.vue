<template>
  <div class="storage-status-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <ServerIcon class="w-5 h-5" />
      <span class="widget-title">存储状态</span>
    </div>

    <div class="widget-content">
      <!-- 总体存储状态 -->
      <div class="storage-overview">
        <div class="overview-header">
          <span class="overview-label">总存储</span>
          <span class="overview-value">{{ formatBytes(totalStorage) }}</span>
        </div>
        <div class="overview-bar">
          <div class="bar-fill" :style="{ width: storagePercentage + '%', background: storageColor }"></div>
        </div>
        <div class="overview-stats">
          <span class="stat-used">已用 {{ formatBytes(usedStorage) }}</span>
          <span class="stat-available">可用 {{ formatBytes(availableStorage) }}</span>
        </div>
      </div>

      <!-- 磁盘列表 -->
      <div v-if="size !== 'small'" class="disk-list">
        <div
          v-for="disk in disks"
          :key="disk.name"
          class="disk-item"
        >
          <div class="disk-header">
            <div class="disk-status" :class="disk.status"></div>
            <span class="disk-name">{{ disk.name }}</span>
            <span class="disk-size">{{ formatBytes(disk.size) }}</span>
          </div>
          <div class="disk-bar">
            <div class="bar-fill" :style="{ width: disk.usagePercent + '%', background: getDiskColor(disk.usagePercent) }"></div>
          </div>
          <div class="disk-stats">
            <span class="disk-temperature">{{ disk.temperature }}°C</span>
            <span class="disk-usage">{{ disk.usagePercent }}% 使用</span>
          </div>
        </div>
      </div>

      <!-- RAID状态 -->
      <div v-if="showRaidStatus" class="raid-status">
        <div class="raid-header">
          <ShieldCheckIcon class="raid-icon" :class="raidStatusClass" />
          <div class="raid-info">
            <span class="raid-name">{{ raidInfo.name }}</span>
            <span class="raid-state">{{ raidInfo.state }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ServerIcon, ShieldCheckIcon } from '@heroicons/vue/24/outline'

interface Props {
  config: {
    showRaid?: boolean
    showTemperature?: boolean
    refreshInterval?: number
  }
  size: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    showRaid: true,
    showTemperature: true,
    refreshInterval: 5000
  })
})

const totalStorage = ref(8 * 1024 * 1024 * 1024 * 1024) // 8TB
const usedStorage = ref(3.2 * 1024 * 1024 * 1024 * 1024) // 3.2TB
const availableStorage = ref(4.8 * 1024 * 1024 * 1024 * 1024) // 4.8TB

const disks = ref([
  {
    name: 'Disk 1',
    size: 2 * 1024 * 1024 * 1024 * 1024,
    usagePercent: 75,
    temperature: 38,
    status: 'healthy'
  },
  {
    name: 'Disk 2',
    size: 2 * 1024 * 1024 * 1024 * 1024,
    usagePercent: 68,
    temperature: 40,
    status: 'healthy'
  },
  {
    name: 'Disk 3',
    size: 2 * 1024 * 1024 * 1024 * 1024,
    usagePercent: 82,
    temperature: 42,
    status: 'healthy'
  },
  {
    name: 'Disk 4',
    size: 2 * 1024 * 1024 * 1024 * 1024,
    usagePercent: 71,
    temperature: 39,
    status: 'healthy'
  }
])

const raidInfo = ref({
  name: 'RAID 5',
  state: '正常'
})

const storagePercentage = computed(() => {
  return (usedStorage.value / totalStorage.value) * 100
})

const storageColor = computed(() => {
  const percent = storagePercentage.value
  if (percent < 50) return '#10b981'
  if (percent < 80) return '#f59e0b'
  return '#ef4444'
})

const showRaidStatus = computed(() => {
  return props.config.showRaid && props.size !== 'small'
})

const raidStatusClass = computed(() => {
  return raidInfo.value.state === '正常' ? 'healthy' : 'warning'
})

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return (bytes / Math.pow(k, i)).toFixed(1) + ' ' + sizes[i]
}

const getDiskColor = (percent: number): string => {
  if (percent < 50) return '#10b981'
  if (percent < 80) return '#f59e0b'
  return '#ef4444'
}

// 模拟存储数据更新
const updateStorageData = () => {
  // 随机更新一些数据
  usedStorage.value = (3.2 + Math.random() * 0.1) * 1024 * 1024 * 1024 * 1024
  disks.value.forEach(disk => {
    disk.usagePercent = Math.min(95, Math.max(60, disk.usagePercent + (Math.random() - 0.5) * 2))
    disk.temperature = Math.floor(38 + Math.random() * 5)
  })
}

let updateInterval: number | null = null

onMounted(() => {
  updateStorageData()
  updateInterval = window.setInterval(updateStorageData, props.config.refreshInterval)
})

onUnmounted(() => {
  if (updateInterval) {
    clearInterval(updateInterval)
  }
})
</script>

<style scoped>
.storage-status-widget {
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

.widget-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 12px;
  overflow-y: auto;
}

.storage-overview {
  padding: 12px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 8px;
}

.overview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.overview-label {
  font-size: 12px;
  color: #6b7280;
}

.overview-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.overview-bar {
  height: 8px;
  background: rgba(209, 213, 219, 0.3);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 6px;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.overview-stats {
  display: flex;
  justify-content: space-between;
  font-size: 11px;
}

.stat-used {
  color: #f59e0b;
  font-weight: 500;
}

.stat-available {
  color: #10b981;
}

.disk-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.disk-item {
  padding: 10px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 8px;
}

.disk-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.disk-status {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.disk-status.healthy {
  background: #10b981;
}

.disk-status.warning {
  background: #f59e0b;
}

.disk-status.error {
  background: #ef4444;
}

.disk-name {
  font-size: 12px;
  font-weight: 500;
  color: #1f2937;
  flex: 1;
}

.disk-size {
  font-size: 11px;
  color: #6b7280;
}

.disk-bar {
  height: 4px;
  background: rgba(209, 213, 219, 0.3);
  border-radius: 2px;
  overflow: hidden;
  margin-bottom: 4px;
}

.disk-stats {
  display: flex;
  justify-content: space-between;
  font-size: 10px;
  color: #6b7280;
}

.raid-status {
  padding: 10px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 8px;
}

.raid-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.raid-icon {
  width: 20px;
  height: 20px;
  color: #10b981;
}

.raid-icon.healthy {
  color: #10b981;
}

.raid-icon.warning {
  color: #f59e0b;
}

.raid-info {
  display: flex;
  flex-direction: column;
}

.raid-name {
  font-size: 12px;
  font-weight: 500;
  color: #1f2937;
}

.raid-state {
  font-size: 11px;
  color: #6b7280;
}

/* 小尺寸widget优化 */
.widget-small .disk-list {
  max-height: 120px;
}

.widget-small .storage-overview {
  padding: 8px;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .storage-status-widget {
    background: rgba(0, 0, 0, 0.8);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .widget-title,
  .overview-value,
  .disk-name {
    color: #f9fafb;
  }

  .overview-label,
  .disk-size,
  .disk-stats,
  .raid-state {
    color: #9ca3af;
  }

  .storage-overview,
  .disk-item,
  .raid-status {
    background: rgba(255, 255, 255, 0.05);
  }

  .overview-bar,
  .disk-bar {
    background: rgba(255, 255, 255, 0.1);
  }
}
</style>