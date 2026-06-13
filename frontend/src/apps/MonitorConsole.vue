<template>
  <div class="monitor-console">
    <!-- System Overview Dashboard -->
    <div class="system-dashboard">
      <div class="dashboard-header">
        <h2>系统监控控制台</h2>
        <div class="header-actions">
          <el-button @click="refreshAll" :loading="monitorStore.loading">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-dropdown @command="handleSettingsCommand">
            <el-button>
              <el-icon><Setting /></el-icon>
              设置
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="refresh-interval">刷新间隔</el-dropdown-item>
                <el-dropdown-item command="alerts">告警配置</el-dropdown-item>
                <el-dropdown-item command="notifications">通知设置</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>

      <!-- System Status Cards -->
      <div class="status-cards">
        <div class="status-card system-load" :class="{ warning: !monitorStore.systemLoad.normal }">
          <div class="card-icon">
            <el-icon><TrendCharts /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">系统负载</div>
            <div class="card-value">{{ currentStats?.cpu.load1.toFixed(2) }}</div>
            <div class="card-subtitle">{{ monitorStore.systemLoad.message }}</div>
          </div>
          <div class="card-indicator" :class="monitorStore.systemLoad.normal ? 'normal' : 'warning'" />
        </div>

        <div class="status-card memory-status" :class="{ warning: !monitorStore.memoryStatus.normal }">
          <div class="card-icon">
            <el-icon><Coin /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">内存使用</div>
            <div class="card-value">{{ currentStats?.memory.percent.toFixed(1) }}%</div>
            <div class="card-subtitle">{{ monitorStore.memoryStatus.message }}</div>
          </div>
          <div class="card-indicator" :class="monitorStore.memoryStatus.normal ? 'normal' : 'warning'" />
        </div>

        <div class="status-card disk-status" :class="{ warning: !monitorStore.diskStatus.normal }">
          <div class="card-icon">
            <el-icon><Files /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">磁盘状态</div>
            <div class="card-value">{{ currentStats?.disk?.length || 0 }}</div>
            <div class="card-subtitle">{{ monitorStore.diskStatus.message }}</div>
          </div>
          <div class="card-indicator" :class="monitorStore.diskStatus.normal ? 'normal' : 'warning'" />
        </div>

        <div class="status-card service-status">
          <div class="card-icon">
            <el-icon><Grid /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">运行服务</div>
            <div class="card-value">{{ monitorStore.runningServices.length }}</div>
            <div class="card-subtitle">{{ monitorStore.failedServices.length }} 个故障</div>
          </div>
          <div class="card-indicator" :class="monitorStore.failedServices.length > 0 ? 'warning' : 'normal'" />
        </div>

        <div class="status-card temperature-status">
          <div class="card-icon">
            <el-icon><Sunny /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">系统温度</div>
            <div class="card-value">{{ avgTemperature }}°C</div>
            <div class="card-subtitle">{{ temperatureStatus }}</div>
          </div>
          <div class="card-indicator" :class="temperatureStatus === '正常' ? 'normal' : 'warning'" />
        </div>

        <div class="status-card active-processes">
          <div class="card-icon">
            <el-icon><Cpu /></el-icon>
          </div>
          <div class="card-content">
            <div class="card-label">活动进程</div>
            <div class="card-value">{{ monitorStore.runningProcesses.length }}</div>
            <div class="card-subtitle">CPU: {{ monitorStore.totalProcessCPU.toFixed(1) }}%</div>
          </div>
          <div class="card-indicator normal" />
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="main-content">
        <!-- Left Panel - Charts and Graphs -->
        <div class="left-panel">
          <!-- CPU Usage Chart -->
          <div class="chart-section">
            <h3>CPU 使用率</h3>
            <div class="chart-container">
              <el-progress
                type="dashboard"
                :percentage="currentStats?.cpu.usage ? (currentStats.cpu.usage * 100) : 0"
                :color="getProgressColor((currentStats?.cpu.usage || 0) * 100)"
                :width="120"
              >
                <template #default="{ percentage }">
                  <span class="percentage-value">{{ percentage.toFixed(1) }}%</span>
                  <span class="percentage-label">CPU</span>
                </template>
              </el-progress>
            </div>
            <div class="chart-details">
              <div class="detail-item">
                <span class="label">核心数:</span>
                <span class="value">{{ currentStats?.cpu.cores || 0 }}</span>
              </div>
              <div class="detail-item">
                <span class="label">1分钟负载:</span>
                <span class="value">{{ currentStats?.cpu.load1?.toFixed(2) || 0 }}</span>
              </div>
              <div class="detail-item">
                <span class="label">5分钟负载:</span>
                <span class="value">{{ currentStats?.cpu.load5?.toFixed(2) || 0 }}</span>
              </div>
              <div class="detail-item">
                <span class="label">15分钟负载:</span>
                <span class="value">{{ currentStats?.cpu.load15?.toFixed(2) || 0 }}</span>
              </div>
            </div>
          </div>

          <!-- Memory Usage Chart -->
          <div class="chart-section">
            <h3>内存使用情况</h3>
            <div class="chart-container">
              <el-progress
                type="dashboard"
                :percentage="currentStats?.memory.percent || 0"
                :color="getProgressColor(currentStats?.memory.percent || 0)"
                :width="120"
              >
                <template #default="{ percentage }">
                  <span class="percentage-value">{{ percentage?.toFixed(1) }}%</span>
                  <span class="percentage-label">内存</span>
                </template>
              </el-progress>
            </div>
            <div class="chart-details">
              <div class="detail-item">
                <span class="label">总内存:</span>
                <span class="value">{{ formatSize(currentStats?.memory.total || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">已使用:</span>
                <span class="value">{{ formatSize(currentStats?.memory.used || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">可用:</span>
                <span class="value">{{ formatSize(currentStats?.memory.available || 0) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">Swap:</span>
                <span class="value">{{ formatSize(currentStats?.memory.swapUsed || 0) }}</span>
              </div>
            </div>
          </div>

          <!-- Recent Events -->
          <div class="events-section">
            <div class="section-header">
              <h3>最近事件</h3>
              <el-button text @click="viewAllEvents">查看全部</el-button>
            </div>
            <div class="events-list">
              <div
                v-for="event in monitorStore.recentEvents"
                :key="event.id"
                class="event-item"
                :class="event.type"
              >
                <div class="event-icon">
                  <el-icon><Warning /></el-icon>
                </div>
                <div class="event-content">
                  <div class="event-title">{{ event.title }}</div>
                  <div class="event-message">{{ event.message }}</div>
                  <div class="event-time">{{ formatTime(event.createdAt) }}</div>
                </div>
              </div>
              <div v-if="monitorStore.recentEvents.length === 0" class="no-events">
                暂无最近事件
              </div>
            </div>
          </div>
        </div>

        <!-- Right Panel - Detailed Information -->
        <div class="right-panel">
          <!-- Tab Navigation -->
          <el-tabs v-model="activeTab" class="detail-tabs">
            <!-- Processes Tab -->
            <el-tab-pane label="进程管理" name="processes">
              <ProcessManager />
            </el-tab-pane>

            <!-- Services Tab -->
            <el-tab-pane label="服务管理" name="services">
              <ServiceManager />
            </el-tab-pane>

            <!-- Temperature Tab -->
            <el-tab-pane label="温度监控" name="temperature">
              <TemperaturePanel />
            </el-tab-pane>

            <!-- Alerts Tab -->
            <el-tab-pane label="告警管理" name="alerts">
              <AlertManager />
            </el-tab-pane>

            <!-- Logs Tab -->
            <el-tab-pane label="系统日志" name="logs">
              <SystemLog />
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </div>

    <!-- Alert Notifications -->
    <div v-if="criticalAlerts.length > 0" class="alert-banner">
      <el-alert
        v-for="alert in criticalAlerts"
        :key="alert.id"
        :title="alert.name"
        :type="alert.severity"
        :description="`阈值: ${alert.threshold}${alert.type === 'cpu' ? '%' : ''}, 当前: ${currentValue(alert)}`"
        show-icon
        :closable="true"
        @close="dismissAlert(alert.id)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useMonitorStore } from '@/stores/monitor'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Setting,
  TrendCharts,
  Coin,
  Files,
  Grid,
  Sunny,
  Cpu,
  Warning
} from '@element-plus/icons-vue'
import ProcessManager from '@/components/Monitor/ProcessManager.vue'
import ServiceManager from '@/components/Monitor/ServiceManager.vue'
import TemperaturePanel from '@/components/Monitor/TemperaturePanel.vue'
import AlertManager from '@/components/Monitor/AlertManager.vue'
import SystemLog from '@/components/Monitor/SystemLog.vue'

const monitorStore = useMonitorStore()

const activeTab = ref('processes')
let refreshInterval: number | null = null

// Computed
const currentStats = computed(() => monitorStore.currentStats)

const avgTemperature = computed(() => {
  if (!monitorStore.temperature?.sensors || monitorStore.temperature.sensors.length === 0) {
    return 0
  }
  const total = monitorStore.temperature.sensors.reduce((sum, sensor) => sum + sensor.current, 0)
  return (total / monitorStore.temperature.sensors.length).toFixed(1)
})

const temperatureStatus = computed(() => {
  if (!monitorStore.temperature?.sensors) return '未知'
  const highTemp = monitorStore.temperature.sensors.some(s => s.current > s.max * 0.9)
  return highTemp ? '偏高' : '正常'
})

const criticalAlerts = computed(() => monitorStore.criticalAlerts)

// Methods
const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const formatTime = (dateStr: string): string => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getProgressColor = (percent: number): string => {
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

const refreshAll = async () => {
  await Promise.all([
    monitorStore.fetchStats(),
    monitorStore.fetchProcesses(),
    monitorStore.fetchServices(),
    monitorStore.fetchTemperature()
  ])
  ElMessage.success('数据已刷新')
}

const handleSettingsCommand = (command: string) => {
  switch (command) {
    case 'refresh-interval':
      ElMessage.info('刷新间隔设置即将推出')
      break
    case 'alerts':
      activeTab.value = 'alerts'
      break
    case 'notifications':
      ElMessage.info('通知设置即将推出')
      break
  }
}

const viewAllEvents = () => {
  // Navigate to events page or show dialog
  ElMessage.info('事件详情页面即将推出')
}

const currentValue = (alert: any): string => {
  // Get current value based on alert type
  switch (alert.type) {
    case 'cpu':
      return `${((currentStats.value?.cpu.usage || 0) * 100).toFixed(1)}%`
    case 'memory':
      return `${currentStats.value?.memory.percent.toFixed(1) || 0}%`
    case 'disk':
      return '检查磁盘'
    case 'temperature':
      return `${avgTemperature.value}°C`
    default:
      return '未知'
  }
}

const dismissAlert = (alertId: number) => {
  // Handle alert dismissal
  ElMessage.info('告警已忽略')
}

// Lifecycle
onMounted(() => {
  monitorStore.init()

  // Set up auto-refresh
  refreshInterval = window.setInterval(() => {
    monitorStore.fetchStats()
  }, 5000) // Refresh every 5 seconds
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
  monitorStore.destroy()
})
</script>

<style scoped lang="scss">
.monitor-console {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.system-dashboard {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h2 {
    margin: 0;
    color: #303133;
    font-size: 24px;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.status-cards {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 15px;
  margin-bottom: 20px;

  .status-card {
    background: white;
    padding: 15px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    position: relative;
    display: flex;
    align-items: center;
    gap: 12px;
    transition: all 0.3s;

    &.warning {
      border-left: 3px solid #f56c6c;
    }

    .card-icon {
      font-size: 24px;
      color: #409eff;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 40px;
      height: 40px;
      border-radius: 50%;
      background: #ecf5ff;
    }

    .card-content {
      flex: 1;

      .card-label {
        font-size: 12px;
        color: #909399;
        margin-bottom: 4px;
      }

      .card-value {
        font-size: 18px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 2px;
      }

      .card-subtitle {
        font-size: 11px;
        color: #606266;
      }
    }

    .card-indicator {
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: #67c23a;

      &.warning {
        background: #f56c6c;
      }
    }
  }
}

.main-content {
  flex: 1;
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 20px;
  min-height: 0;
}

.left-panel {
  display: flex;
  flex-direction: column;
  gap: 20px;
  overflow-y: auto;
}

.chart-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  h3 {
    margin: 0 0 15px;
    color: #303133;
    font-size: 14px;
  }

  .chart-container {
    display: flex;
    justify-content: center;
    margin-bottom: 20px;
  }

  .chart-details {
    .detail-item {
      display: flex;
      justify-content: space-between;
      padding: 8px 0;
      border-bottom: 1px solid #f0f0f0;

      &:last-child {
        border-bottom: none;
      }

      .label {
        color: #909399;
        font-size: 12px;
      }

      .value {
        color: #303133;
        font-size: 12px;
        font-weight: 500;
      }
    }
  }
}

.events-section {
  background: white;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;

    h3 {
      margin: 0;
      color: #303133;
      font-size: 14px;
    }
  }

  .events-list {
    max-height: 300px;
    overflow-y: auto;

    .event-item {
      display: flex;
      gap: 12px;
      padding: 12px;
      border-radius: 4px;
      margin-bottom: 8px;
      background: #f5f7fa;

      &.critical {
        background: #fef0f0;
        border-left: 3px solid #f56c6c;
      }

      &.warning {
        background: #fdf6ec;
        border-left: 3px solid #e6a23c;
      }

      .event-icon {
        font-size: 16px;
        color: #e6a23c;
        display: flex;
        align-items: center;
      }

      .event-content {
        flex: 1;

        .event-title {
          font-weight: 500;
          color: #303133;
          margin-bottom: 4px;
        }

        .event-message {
          font-size: 12px;
          color: #606266;
          margin-bottom: 4px;
        }

        .event-time {
          font-size: 11px;
          color: #909399;
        }
      }
    }

    .no-events {
      text-align: center;
      color: #909399;
      padding: 20px;
      font-size: 12px;
    }
  }
}

.right-panel {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;

  .detail-tabs {
    height: 100%;

    :deep(.el-tabs__content) {
      height: calc(100% - 55px);
      overflow: auto;
    }

    :deep(.el-tab-pane) {
      height: 100%;
    }
  }
}

.alert-banner {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  max-width: 400px;
}

.percentage-value {
  display: block;
  font-size: 20px;
  font-weight: bold;
  color: #303133;
}

.percentage-label {
  display: block;
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>