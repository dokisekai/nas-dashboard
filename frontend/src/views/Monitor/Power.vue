<template>
  <div class="power-monitor">
    <div class="monitor-header">
      <h2>功耗监控</h2>
      <div class="header-actions">
        <el-button @click="refreshData" :loading="loading" size="small">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 告警信息 -->
    <el-alert
      v-for="(alert, index) in alerts"
      :key="index"
      :type="getAlertType(alert)"
      :closable="false"
      show-icon
      style="margin-bottom: 15px"
    >
      {{ alert }}
    </el-alert>

    <!-- 实时功耗概览 -->
    <el-row :gutter="20" class="power-overview">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>实时功耗概览</span>
          </template>
          <div v-if="currentPower" class="power-grid">
            <div class="power-item main">
              <div class="power-label">总功耗</div>
              <div class="power-value">{{ currentPower.total.toFixed(1) }} W</div>
              <div class="power-desc">整机当前功耗</div>
            </div>

            <div class="power-item">
              <div class="power-label">CPU Package</div>
              <div class="power-value">{{ currentPower.cpuPackage.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.cpuPackage / 150) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.cpuPackage, 150)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">CPU Core</div>
              <div class="power-value">{{ currentPower.cpuCore.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.cpuCore / 100) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.cpuCore, 100)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">CPU Uncore</div>
              <div class="power-value">{{ currentPower.cpuUncore.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.cpuUncore / 50) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.cpuUncore, 50)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">Intel 核显</div>
              <div class="power-value">{{ currentPower.igpu.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.igpu / 30) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.igpu, 30)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">AMD 独显</div>
              <div class="power-value">{{ currentPower.dgpu.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.dgpu / 300) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.dgpu, 300)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">HDD 机械硬盘</div>
              <div class="power-value">{{ currentPower.hdd.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.hdd / 30) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.hdd, 30)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">SSD 固态硬盘</div>
              <div class="power-value">{{ currentPower.ssd.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.ssd / 20) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.ssd, 20)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">主板和内存</div>
              <div class="power-value">{{ currentPower.mbram.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.mbram / 50) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.mbram, 50)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">散热系统</div>
              <div class="power-value">{{ currentPower.cooling.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.cooling / 30) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.cooling, 30)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">USB 外设</div>
              <div class="power-value">{{ currentPower.usb.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.usb / 10) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.usb, 10)"
              />
            </div>

            <div class="power-item">
              <div class="power-label">电源损耗</div>
              <div class="power-value">{{ currentPower.powerLoss.toFixed(1) }} W</div>
              <el-progress
                :percentage="((currentPower.powerLoss / 50) * 100).toFixed(1)"
                :color="getProgressColor(currentPower.powerLoss, 50)"
              />
            </div>
          </div>
          <div v-else class="loading-placeholder">
            <el-skeleton :rows="6" animated />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 统计信息 -->
    <el-row :gutter="20" class="stats-section">
      <el-col :span="12">
        <el-card>
          <template #header>今日统计</template>
          <div v-if="todayStats" class="stats-list">
            <div class="stat-item">
              <span class="stat-label">平均功耗:</span>
              <span class="stat-value">{{ todayStats.averagePower.toFixed(1) }} W</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">最高功耗:</span>
              <span class="stat-value">{{ todayStats.maxPower.toFixed(1) }} W</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">最低功耗:</span>
              <span class="stat-value">{{ todayStats.minPower.toFixed(1) }} W</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">总耗电量:</span>
              <span class="stat-value">{{ todayStats.totalEnergy.toFixed(2) }} kWh</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">采样次数:</span>
              <span class="stat-value">{{ todayStats.sampleCount }}</span>
            </div>
          </div>
          <div v-else class="loading-placeholder">
            <el-skeleton :rows="3" animated />
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>电费估算</template>
          <div class="cost-estimate">
            <div class="cost-item">
              <span class="cost-label">今日耗电:</span>
              <span class="cost-value">{{ todayStats?.totalEnergy.toFixed(2) || 0 }} kWh</span>
            </div>
            <div class="cost-item">
              <span class="cost-label">电费单价:</span>
              <span class="cost-value">¥ 0.52/kWh</span>
            </div>
            <div class="cost-item total">
              <span class="cost-label">预计电费:</span>
              <span class="cost-value">¥ {{ ((todayStats?.totalEnergy || 0) * 0.52).toFixed(2) }}</span>
            </div>
            <div class="cost-desc">
              *按平均工业电费估算，实际费用以当地电价为准
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 历史趋势 -->
    <el-row :gutter="20" class="chart-section">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="chart-header">
              <span>功耗趋势</span>
              <el-radio-group v-model="selectedPeriod" @change="loadHistory">
                <el-radio-button :label="1">1天</el-radio-button>
                <el-radio-button :label="7">7天</el-radio-button>
                <el-radio-button :label="30">30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div v-if="chartData.length > 0" class="chart-container">
            <canvas ref="chartCanvas"></canvas>
          </div>
          <div v-else class="loading-placeholder">
            <el-skeleton :rows="1" animated />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { powerAPI } from '@/api/power'
import type { PowerData, PowerStatistics } from '@/types/power'
import { Refresh } from '@element-plus/icons-vue'
import { ElProgress, ElSkeleton, ElAlert } from 'element-plus'
import Chart from 'chart.js/auto'

// 响应式数据
const loading = ref(false)
const currentPower = ref<PowerData | null>(null)
const todayStats = ref<any>(null)
const alerts = ref<string[]>([])
const selectedPeriod = ref(7)
const chartData = ref<any[]>([])
const chartCanvas = ref<HTMLCanvasElement | null>(null)
let chart: Chart | null = null

// 获取告警类型
const getAlertType = (alert: string) => {
  if (alert.includes('严重')) return 'error'
  if (alert.includes('告警')) return 'warning'
  return 'info'
}

// 获取进度条颜色
const getProgressColor = (value: number, max: number) => {
  const percentage = (value / max) * 100
  if (percentage > 80) return '#f56c6c'
  if (percentage > 60) return '#e6a23c'
  return '#67c23a'
}

// 刷新数据
const refreshData = async () => {
  loading.value = true
  try {
    // 获取功耗概览
    const overview = await powerAPI.getOverview()
    currentPower.value = overview.current
    todayStats.value = overview.today
    alerts.value = overview.alerts

    // 加载历史数据
    await loadHistory()
  } catch (error) {
    console.error('获取功耗数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 加载历史数据
const loadHistory = async () => {
  try {
    const history = await powerAPI.getHistory(selectedPeriod.value)
    chartData.value = history.data
    renderChart()
  } catch (error) {
    console.error('获取历史数据失败:', error)
  }
}

// 渲染图表
const renderChart = () => {
  if (!chartCanvas.value || chartData.value.length === 0) return

  const ctx = chartCanvas.value.getContext('2d')
  if (!ctx) return

  // 销毁旧图表
  if (chart) {
    chart.destroy()
  }

  // 准备数据
  const labels = chartData.value.map(item => {
    const date = new Date(item.timestamp)
    return `${date.getHours()}:${date.getMinutes().toString().padStart(2, '0')}`
  })
  const data = chartData.value.map(item => item.total)

  // 创建新图表
  chart = new Chart(ctx, {
    type: 'line',
    data: {
      labels: labels,
      datasets: [{
        label: '总功耗 (W)',
        data: data,
        borderColor: '#409EFF',
        backgroundColor: 'rgba(64, 158, 255, 0.1)',
        borderWidth: 2,
        tension: 0.4,
        fill: true
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: {
          display: true
        }
      },
      scales: {
        y: {
          beginAtZero: false,
          title: {
            display: true,
            text: '功耗 (W)'
          }
        }
      }
    }
  })
}

// 组件挂载
onMounted(() => {
  refreshData()

  // 设置定时刷新
  const interval = setInterval(() => {
    refreshData()
  }, 60000) // 每分钟刷新

  onUnmounted(() => {
    clearInterval(interval)
    if (chart) {
      chart.destroy()
    }
  })
})
</script>

<style scoped>
.power-monitor {
  padding: 20px;
}

.monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.monitor-header h2 {
  margin: 0;
  color: #303133;
}

.power-overview {
  margin-bottom: 20px;
}

.power-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 15px;
}

.power-item {
  padding: 15px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  text-align: center;
  background: #fafafa;
}

.power-item.main {
  grid-column: span 2;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.power-item.main .power-label {
  color: rgba(255, 255, 255, 0.8);
}

.power-item.main .power-value {
  color: white;
  font-size: 32px;
}

.power-item.main .power-desc {
  color: rgba(255, 255, 255, 0.7);
  font-size: 14px;
  margin-top: 5px;
}

.power-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.power-value {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 10px;
}

.stats-section {
  margin-bottom: 20px;
}

.stats-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.stat-label {
  color: #606266;
}

.stat-value {
  font-weight: bold;
  color: #303133;
}

.cost-estimate {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.cost-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.cost-item.total {
  padding-top: 10px;
  border-top: 2px solid #409EFF;
}

.cost-item.total .cost-value {
  font-size: 18px;
  color: #409EFF;
  font-weight: bold;
}

.cost-desc {
  font-size: 12px;
  color: #909399;
  text-align: center;
  margin-top: 10px;
}

.chart-section {
  margin-bottom: 20px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 300px;
  position: relative;
}

.loading-placeholder {
  padding: 20px;
}
</style>
