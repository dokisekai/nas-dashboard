<template>
  <div class="temperature-panel">
    <div class="panel-header">
      <div class="header-controls">
        <el-button size="small" @click="refreshTemperature" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-select
          v-model="updateInterval"
          size="small"
          style="width: 100px"
          @change="handleIntervalChange"
        >
          <el-option label="5秒" :value="5000" />
          <el-option label="10秒" :value="10000" />
          <el-option label="30秒" :value="30000" />
          <el-option label="60秒" :value="60000" />
        </el-select>
      </div>
      <div class="temp-stats">
        <span class="stat-item">传感器: {{ sensorCount }}</span>
        <span class="stat-item" :class="{ warning: hasHighTemp }">
          平均温度: {{ averageTemp }}°C
        </span>
        <span class="stat-item" :class="{ danger: hasCriticalTemp }">
          最高温度: {{ maxTemp }}°C
        </span>
      </div>
    </div>

    <div class="sensors-grid">
      <div
        v-for="sensor in temperature?.sensors"
        :key="sensor.name"
        class="sensor-card"
        :class="getSensorClass(sensor)"
      >
        <div class="card-header">
          <div class="sensor-icon">
            <el-icon><Sunny /></el-icon>
          </div>
          <div class="sensor-name">{{ sensor.name }}</div>
          <el-tag
            :type="getSensorTagType(sensor)"
            size="small"
            class="sensor-status"
          >
            {{ getSensorStatus(sensor) }}
          </el-tag>
        </div>

        <div class="card-body">
          <div class="temp-display">
            <div class="current-temp" :class="getTempClass(sensor)">
              {{ sensor.current.toFixed(1) }}°C
            </div>
            <div class="temp-range">
              <span class="min">最低: {{ sensor.min.toFixed(1) }}°C</span>
              <span class="max">最高: {{ sensor.max.toFixed(1) }}°C</span>
            </div>
          </div>

          <div class="temp-bar">
            <el-progress
              :percentage="getTempPercent(sensor)"
              :color="getTempBarColor(sensor)"
              :show-text="false"
              :stroke-width="8"
            />
          </div>

          <div v-if="sensor.critical > 0" class="temp-threshold">
            <span class="threshold-label">临界温度:</span>
            <span class="threshold-value">{{ sensor.critical.toFixed(1) }}°C</span>
          </div>
        </div>

        <div class="card-footer">
          <div class="temp-history">
            <svg class="temp-chart" viewBox="0 0 100 20">
              <polyline
                :points="generateChartPoints(sensor)"
                fill="none"
                :stroke="getTempBarColor(sensor)"
                stroke-width="2"
              />
            </svg>
          </div>
        </div>
      </div>

      <div v-if="!temperature?.sensors || temperature.sensors.length === 0" class="no-sensors">
        <el-empty description="未检测到温度传感器" />
      </div>
    </div>

    <!-- Temperature History Chart -->
    <div class="temp-history-section">
      <h3>温度历史趋势</h3>
      <div class="chart-container">
        <canvas ref="tempChart" class="temp-chart-canvas"></canvas>
      </div>
    </div>

    <!-- Temperature Alerts -->
    <div v-if="highTempSensors.length > 0" class="temp-alerts">
      <el-alert
        v-for="sensor in highTempSensors"
        :key="sensor.name"
        :title="`高温告警 - ${sensor.name}`"
        :type="getAlertType(sensor)"
        :description="`当前温度: ${sensor.current.toFixed(1)}°C, 临界温度: ${sensor.critical.toFixed(1)}°C`"
        show-icon
        :closable="true"
        @close="dismissAlert(sensor.name)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useMonitorStore } from '@/stores/monitor'
import { ElMessage } from 'element-plus'
import type { TemperatureInfo, Sensor } from '@/types/monitor'
import { Refresh, Sunny } from '@element-plus/icons-vue'

const monitorStore = useMonitorStore()

const loading = ref(false)
const updateInterval = ref(10000)
const temperature = ref<TemperatureInfo | null>(null)
const tempChart = ref<HTMLCanvasElement | null>(null)
let intervalId: number | null = null

// Mock temperature history data
const tempHistory = ref<Record<string, number[]>>({})

// Computed
const sensorCount = computed(() => temperature.value?.sensors?.length || 0)

const averageTemp = computed(() => {
  if (!temperature.value?.sensors || temperature.value.sensors.length === 0) return 0
  const total = temperature.value.sensors.reduce((sum, sensor) => sum + sensor.current, 0)
  return (total / temperature.value.sensors.length).toFixed(1)
})

const maxTemp = computed(() => {
  if (!temperature.value?.sensors || temperature.value.sensors.length === 0) return 0
  return Math.max(...temperature.value.sensors.map(s => s.current)).toFixed(1)
})

const hasHighTemp = computed(() => {
  if (!temperature.value?.sensors) return false
  return temperature.value.sensors.some(s => s.current > s.max * 0.8)
})

const hasCriticalTemp = computed(() => {
  if (!temperature.value?.sensors) return false
  return temperature.value.sensors.some(s => s.current > s.critical * 0.9)
})

const highTempSensors = computed(() => {
  if (!temperature.value?.sensors) return []
  return temperature.value.sensors.filter(s => s.current > s.max * 0.8)
})

// Methods
const refreshTemperature = async () => {
  loading.value = true
  try {
    await monitorStore.fetchTemperature()
    temperature.value = monitorStore.temperature

    // Update history
    if (temperature.value?.sensors) {
      temperature.value.sensors.forEach(sensor => {
        if (!tempHistory.value[sensor.name]) {
          tempHistory.value[sensor.name] = []
        }
        tempHistory.value[sensor.name].push(sensor.current)

        // Keep only last 20 data points
        if (tempHistory.value[sensor.name].length > 20) {
          tempHistory.value[sensor.name].shift()
        }
      })
    }

    ElMessage.success('温度数据已刷新')
  } finally {
    loading.value = false
  }
}

const getSensorClass = (sensor: Sensor): string => {
  const percent = sensor.critical > 0 ? (sensor.current / sensor.critical) * 100 : 0
  if (percent >= 90) return 'critical'
  if (percent >= 80) return 'warning'
  return 'normal'
}

const getSensorTagType = (sensor: Sensor): string => {
  if (sensor.current >= sensor.critical * 0.9) return 'danger'
  if (sensor.current >= sensor.max * 0.8) return 'warning'
  return 'success'
}

const getSensorStatus = (sensor: Sensor): string => {
  if (sensor.current >= sensor.critical * 0.9) return '临界'
  if (sensor.current >= sensor.max * 0.8) return '偏高'
  return '正常'
}

const getTempClass = (sensor: Sensor): string => {
  if (sensor.current >= sensor.critical * 0.9) return 'critical'
  if (sensor.current >= sensor.max * 0.8) return 'warning'
  return 'normal'
}

const getTempPercent = (sensor: Sensor): number => {
  if (sensor.critical === 0) return 0
  return Math.min((sensor.current / sensor.critical) * 100, 100)
}

const getTempBarColor = (sensor: Sensor): string => {
  if (sensor.current >= sensor.critical * 0.9) return '#f56c6c'
  if (sensor.current >= sensor.max * 0.8) return '#e6a23c'
  return '#67c23a'
}

const getAlertType = (sensor: Sensor): string => {
  if (sensor.current >= sensor.critical * 0.9) return 'error'
  return 'warning'
}

const generateChartPoints = (sensor: Sensor): string => {
  const history = tempHistory.value[sensor.name] || []
  if (history.length === 0) return '0,10 10,10 20,10 30,10 40,10 50,10 60,10 70,10 80,10 90,10 100,10'

  const min = Math.min(...history, sensor.current)
  const max = Math.max(...history, sensor.current)
  const range = max - min || 1

  const points = history.map((temp, index) => {
    const x = (index / (history.length - 1)) * 100
    const y = 10 + ((max - temp) / range) * 10
    return `${x},${y}`
  }).join(' ')

  return points
}

const dismissAlert = (sensorName: string) => {
  ElMessage.info(`已忽略 ${sensorName} 的告警`)
}

const handleIntervalChange = (interval: number) => {
  if (intervalId) {
    clearInterval(intervalId)
  }

  intervalId = window.setInterval(() => {
    refreshTemperature()
  }, interval)
}

// Lifecycle
onMounted(() => {
  refreshTemperature()

  // Set up auto-refresh
  intervalId = window.setInterval(() => {
    refreshTemperature()
  }, updateInterval.value)
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId)
  }
})
</script>

<style scoped lang="scss">
.temperature-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;

  .header-controls {
    display: flex;
    gap: 10px;
  }

  .temp-stats {
    display: flex;
    gap: 20px;
    font-size: 12px;
    color: #606266;

    .stat-item {
      display: flex;
      align-items: center;
      gap: 5px;

      &.warning { color: #e6a23c; }
      &.danger { color: #f56c6c; }
    }
  }
}

.sensors-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 15px;
  overflow-y: auto;
  margin-bottom: 20px;

  .sensor-card {
    background: white;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    border-left: 4px solid #67c23a;

    &.warning { border-left-color: #e6a23c; }
    &.critical { border-left-color: #f56c6c; }

    .card-header {
      display: flex;
      align-items: center;
      gap: 10px;
      padding: 12px 15px;
      background: #f5f7fa;
      border-bottom: 1px solid #ebeef5;

      .sensor-icon {
        font-size: 16px;
        color: #e6a23c;
      }

      .sensor-name {
        flex: 1;
        font-size: 13px;
        font-weight: 500;
        color: #303133;
      }

      .sensor-status {
        font-size: 11px;
      }
    }

    .card-body {
      padding: 15px;

      .temp-display {
        text-align: center;
        margin-bottom: 15px;

        .current-temp {
          font-size: 32px;
          font-weight: bold;
          margin-bottom: 8px;

          &.normal { color: #67c23a; }
          &.warning { color: #e6a23c; }
          &.critical { color: #f56c6c; }
        }

        .temp-range {
          display: flex;
          justify-content: center;
          gap: 15px;
          font-size: 11px;
          color: #909399;

          .min, .max {
            display: flex;
            align-items: center;
            gap: 4px;
          }
        }
      }

      .temp-bar {
        margin-bottom: 12px;
      }

      .temp-threshold {
        display: flex;
        justify-content: space-between;
        font-size: 11px;
        color: #606266;

        .threshold-label {
          color: #909399;
        }

        .threshold-value {
          font-weight: 500;
        }
      }
    }

    .card-footer {
      padding: 10px 15px;
      background: #fafafa;
      border-top: 1px solid #ebeef5;

      .temp-history {
        .temp-chart {
          width: 100%;
          height: 30px;
        }
      }
    }
  }

  .no-sensors {
    grid-column: 1 / -1;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 200px;
  }
}

.temp-history-section {
  margin-bottom: 20px;

  h3 {
    margin: 0 0 15px;
    color: #303133;
    font-size: 14px;
  }

  .chart-container {
    background: white;
    padding: 15px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    height: 200px;

    .temp-chart-canvas {
      width: 100%;
      height: 100%;
    }
  }
}

.temp-alerts {
  .el-alert {
    margin-bottom: 10px;

    &:last-child {
      margin-bottom: 0;
    }
  }
}
</style>