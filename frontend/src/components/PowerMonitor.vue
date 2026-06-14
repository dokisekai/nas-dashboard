<template>
  <div class="power-monitor">
    <el-card class="overview-card">
      <template #header>
        <div class="card-header">
          <span>功耗监控概览</span>
          <el-button-group>
            <el-button size="small" @click="refreshData" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button size="small" @click="showSettings = true">
              <el-icon><Setting /></el-icon>
              设置
            </el-button>
          </el-button-group>
        </div>
      </template>

      <!-- 告警信息 -->
      <div v-if="overview.alerts && overview.alerts.length > 0" class="alert-section">
        <el-alert
          v-for="(alert, index) in overview.alerts"
          :key="index"
          :type="getAlertType(alert)"
          :closable="false"
          show-icon
        >
          {{ alert }}
        </el-alert>
      </div>

      <!-- 实时功耗 -->
      <div class="power-grid">
        <div class="power-item main">
          <div class="power-label">总功耗</div>
          <div class="power-value">{{ overview.current?.total.toFixed(1) || 0 }} W</div>
          <div class="power-trend" :class="getTrendClass()">
            <el-icon><TrendCharts /></el-icon>
            {{ getTrendText() }}
          </div>
        </div>

        <div class="power-item">
          <div class="power-label">CPU Package</div>
          <div class="power-value">{{ overview.current?.cpuPackage.toFixed(1) || 0 }} W</div>
          <div class="power-bar">
            <div
              class="power-bar-fill"
              :style="{ width: `${(overview.current?.cpuPackage || 0) / 2}%` }"
              :style="{ background: 'var(--el-color-primary)' }"
            ></div>
          </div>
        </div>

        <div class="power-item">
          <div class="power-label">Intel 核显</div>
          <div class="power-value">{{ overview.current?.igpu.toFixed(1) || 0 }} W</div>
          <div class="power-bar">
            <div
              class="power-bar-fill"
              :style="{ width: `${(overview.current?.igpu || 0) * 10}%` }"
              style="background: var(--el-color-success)"
            ></div>
          </div>
        </div>

        <div class="power-item">
          <div class="power-label">AMD 独显</div>
          <div class="power-value">{{ overview.current?.dgpu.toFixed(1) || 0 }} W</div>
          <div class="power-bar">
            <div
              class="power-bar-fill"
              :style="{ width: `${(overview.current?.dgpu || 0) / 3}%` }"
              style="background: var(--el-color-warning)"
            ></div>
          </div>
        </div>

        <div class="power-item">
          <div class="power-label">存储设备</div>
          <div class="power-value">
            {{ ((overview.current?.hdd || 0) + (overview.current?.ssd || 0)).toFixed(1) }} W
          </div>
          <div class="power-breakdown">
            <span>HDD: {{ overview.current?.hdd.toFixed(1) || 0 }}W</span>
            <span>SSD: {{ overview.current?.ssd.toFixed(1) || 0 }}W</span>
          </div>
        </div>

        <div class="power-item">
          <div class="power-label">其他组件</div>
          <div class="power-value">
            {{ ((overview.current?.mbram || 0) + (overview.current?.cooling || 0) + (overview.current?.usb || 0)).toFixed(1) }} W
          </div>
          <div class="power-breakdown">
            <span>主板+内存: {{ overview.current?.mbram.toFixed(1) || 0 }}W</span>
            <span>散热: {{ overview.current?.cooling.toFixed(1) || 0 }}W</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 统计图表 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="8">
        <el-card>
          <template #header>今日统计</template>
          <div class="stat-item">
            <span class="stat-label">平均功耗:</span>
            <span class="stat-value">{{ overview.today?.averagePower.toFixed(1) || 0 }} W</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">最高功耗:</span>
            <span class="stat-value">{{ overview.today?.maxPower.toFixed(1) || 0 }} W</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">最低功耗:</span>
            <span class="stat-value">{{ overview.today?.minPower.toFixed(1) || 0 }} W</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">总耗电量:</span>
            <span class="stat-value">{{ overview.today?.totalEnergy.toFixed(2) || 0 }} kWh</span>
          </div>
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>功耗趋势</span>
              <el-radio-group v-model="selectedPeriod" @change="updateChart">
                <el-radio-button label="1">1天</el-radio-button>
                <el-radio-button label="7">7天</el-radio-button>
                <el-radio-button label="30">30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="chartRef" class="power-chart"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 设置对话框 -->
    <el-dialog v-model="showSettings" title="功耗监控设置" width="600px">
      <el-form :model="settings" label-width="120px">
        <el-form-item label="告警阈值">
          <el-form-item label="高功耗告警:">
            <el-input-number v-model="settings.highThreshold" :min="50" :max="500" />
            <span class="unit">W</span>
          </el-form-item>
          <el-form-item label="严重告警:">
            <el-input-number v-model="settings.criticalThreshold" :min="100" :max="1000" />
            <span class="unit">W</span>
          </el-form-item>
        </el-form-item>

        <el-form-item label="监控间隔">
          <el-input-number v-model="settings.monitorInterval" :min="1" :max="60" />
          <span class="unit">分钟</span>
        </el-form-item>

        <el-form-item label="数据保留">
          <el-input-number v-model="settings.dataRetention" :min="7" :max="365" />
          <span class="unit">天</span>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showSettings = false">取消</el-button>
        <el-button type="primary" @click="saveSettings">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { powerAPI } from '@/api/power'
import type { PowerOverview } from '@/types/power'
import { Refresh, Setting, TrendCharts } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const overview = ref<PowerOverview>({
  current: {} as any,
  today: {} as any,
  alerts: [],
  timestamp: ''
})

const loading = ref(false)
const showSettings = ref(false)
const selectedPeriod = ref('7')
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null
let cleanup: (() => void) | null = null

const settings = ref({
  highThreshold: 150,
  criticalThreshold: 200,
  monitorInterval: 5,
  dataRetention: 90
})

const getAlertType = (alert: string) => {
  if (alert.includes('严重')) return 'error'
  if (alert.includes('告警')) return 'warning'
  return 'info'
}

const getTrendClass = () => {
  // 简化的趋势判断
  return 'trend-stable'
}

const getTrendText = () => {
  return '稳定'
}

const refreshData = async () => {
  loading.value = true
  try {
    overview.value = await powerAPI.getOverview()
    updateChart()
  } catch (error) {
    console.error('Failed to fetch power data:', error)
  } finally {
    loading.value = false
  }
}

const updateChart = async () => {
  if (!chartRef.value) return

  if (!chart) {
    chart = echarts.init(chartRef.value)
  }

  try {
    const days = parseInt(selectedPeriod.value)
    const history = await powerAPI.getHistory(days)

    const timestamps = history.data.map(item => {
      const date = new Date(item.timestamp)
      return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
    })
    const powerValues = history.data.map(item => item.total)

    const option = {
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          return `${params[0].name}<br/>功耗: ${params[0].value.toFixed(1)} W`
        }
      },
      xAxis: {
        type: 'category',
        data: timestamps,
        boundaryGap: false
      },
      yAxis: {
        type: 'value',
        name: '功耗 (W)'
      },
      series: [{
        name: '功耗',
        type: 'line',
        smooth: true,
        data: powerValues,
        areaStyle: {
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
            ]
          }
        },
        lineStyle: {
          color: '#409EFF',
          width: 2
        },
        itemStyle: {
          color: '#409EFF'
        }
      }]
    }

    chart.setOption(option)
  } catch (error) {
    console.error('Failed to update chart:', error)
  }
}

const saveSettings = () => {
  // 这里应该保存设置到后端
  console.log('Saving settings:', settings.value)
  showSettings.value = false
}

onMounted(() => {
  refreshData()

  // 设置 WebSocket 监听
  cleanup = powerAPI.onPowerUpdate((data) => {
    overview.value.current = data
  })

  // 定期刷新
  const interval = setInterval(refreshData, 60000) // 每分钟刷新

  onUnmounted(() => {
    if (cleanup) cleanup()
    if (chart) chart.dispose()
    clearInterval(interval)
  })
})
</script>

<style scoped>
.power-monitor {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.alert-section {
  margin-bottom: 20px;
}

.power-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.power-item {
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 8px;
  text-align: center;
}

.power-item.main {
  grid-column: span 2;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.power-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 8px;
}

.power-item.main .power-label {
  color: rgba(255, 255, 255, 0.8);
}

.power-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 8px;
}

.power-item.main .power-value {
  font-size: 32px;
}

.power-trend {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  font-size: 14px;
}

.power-bar {
  height: 6px;
  background: #f0f0f0;
  border-radius: 3px;
  overflow: hidden;
  margin-top: 8px;
}

.power-bar-fill {
  height: 100%;
  transition: width 0.3s;
}

.power-breakdown {
  display: flex;
  justify-content: space-around;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.stats-row {
  margin-top: 20px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-label {
  color: #666;
}

.stat-value {
  font-weight: bold;
  color: #333;
}

.power-chart {
  height: 300px;
}

.unit {
  margin-left: 8px;
  color: #999;
}
</style>
