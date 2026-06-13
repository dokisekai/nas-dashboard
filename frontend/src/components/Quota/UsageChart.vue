<template>
  <div class="usage-chart">
    <div class="chart-header">
      <h3>配额使用情况</h3>
      <div class="chart-controls">
        <el-select v-model="timeRange" size="small" style="width: 120px">
          <el-option label="今日" value="today" />
          <el-option label="本周" value="week" />
          <el-option label="本月" value="month" />
        </el-select>
        <el-button size="small" @click="refreshChart">
          <el-icon><Refresh /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- Chart Type Selector -->
    <el-radio-group v-model="chartType" size="small">
      <el-radio-button label="overview">总览</el-radio-button>
      <el-radio-button label="trend">趋势</el-radio-button>
      <el-radio-button label="ranking">排名</el-radio-button>
    </el-radio-group>

    <!-- Overview Chart -->
    <div v-if="chartType === 'overview'" class="overview-chart">
      <div class="usage-summary">
        <div class="summary-card total-users">
          <div class="card-icon"><el-icon><User /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ totalUsers }}</div>
            <div class="card-label">用户总数</div>
          </div>
        </div>
        <div class="summary-card active-users">
          <div class="card-icon"><el-icon><UserFilled /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ activeUsers }}</div>
            <div class="card-label">活跃用户</div>
          </div>
        </div>
        <div class="summary-card quota-used">
          <div class="card-icon"><el-icon><FolderChecked /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ formatSize(totalQuotaUsed) }}</div>
            <div class="card-label">已用配额</div>
          </div>
        </div>
        <div class="summary-card quota-total">
          <div class="card-icon"><el-icon><FolderOpened /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ formatSize(totalQuotaTotal) }}</div>
            <div class="card-label">总配额</div>
          </div>
        </div>
      </div>

      <!-- Usage Distribution Chart -->
      <div class="usage-distribution">
        <h4>配额使用分布</h4>
        <div class="distribution-chart">
          <canvas ref="distributionCanvas"></canvas>
        </div>
        <div class="distribution-legend">
          <div class="legend-item">
            <div class="legend-color" style="background: #67c23a"></div>
            <span class="legend-label">正常使用 (&lt;80%)</span>
          </div>
          <div class="legend-item">
            <div class="legend-color" style="background: #e6a23c"></div>
            <span class="legend-label">警告 (80-95%)</span>
          </div>
          <div class="legend-item">
            <div class="legend-color" style="background: #f56c6c"></div>
            <span class="legend-label">严重 (&gt;95%)</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Trend Chart -->
    <div v-if="chartType === 'trend'" class="trend-chart">
      <h4>配额使用趋势</h4>
      <div class="trend-canvas-container">
        <canvas ref="trendCanvas"></canvas>
      </div>
      <div class="trend-stats">
        <div class="stat-item">
          <span class="stat-label">平均增长率:</span>
          <span class="stat-value">{{ averageGrowthRate }}%</span>
        </div>
        <div class="stat-item">
          <span class="stat-label">预计满载时间:</span>
          <span class="stat-value">{{ estimatedFullDate }}</span>
        </div>
      </div>
    </div>

    <!-- Ranking Chart -->
    <div v-if="chartType === 'ranking'" class="ranking-chart">
      <h4>Top 10 配额使用排行</h4>
      <div class="ranking-list">
        <div
          v-for="(item, index) in topUsers"
          :key="item.username"
          class="ranking-item"
        >
          <div class="ranking-position">{{ index + 1 }}</div>
          <div class="ranking-info">
            <div class="user-name">{{ item.username }}</div>
            <div class="usage-bar">
              <div
                class="usage-fill"
                :class="getUsageClass(item.usagePercent)"
                :style="{ width: item.usagePercent + '%' }"
              />
            </div>
            <div class="usage-info">
              <span class="usage-percent">{{ item.usagePercent.toFixed(1) }}%</span>
              <span class="usage-size">{{ formatSize(item.usedSpace) }} / {{ formatSize(item.hardLimit) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useQuotaStore } from '@/stores/quota'
import { Refresh, User, UserFilled, FolderChecked, FolderOpened } from '@element-plus/icons-vue'

const quotaStore = useQuotaStore()

const chartType = ref('overview')
const timeRange = ref('week')
const distributionCanvas = ref<HTMLCanvasElement>()
const trendCanvas = ref<HTMLCanvasElement>()

// Mock data
const totalUsers = computed(() => quotaStore.userQuotas.length)
const activeUsers = computed(() => quotaStore.userQuotas.filter(q => q.usedSpace > 0).length)
const totalQuotaUsed = computed(() => quotaStore.userQuotas.reduce((sum, q) => sum + q.usedSpace, 0))
const totalQuotaTotal = computed(() => quotaStore.userQuotas.reduce((sum, q) => sum + q.hardLimit, 0))

const topUsers = computed(() => {
  return quotaStore.userQuotas
    .map(q => ({
      username: q.user.username,
      usedSpace: q.usedSpace,
      hardLimit: q.hardLimit,
      usagePercent: (q.usedSpace / q.hardLimit) * 100
    }))
    .sort((a, b) => b.usagePercent - a.usagePercent)
    .slice(0, 10)
})

const averageGrowthRate = ref(15.6)
const estimatedFullDate = ref('2026-12-15')

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const getUsageClass = (percent: number): string => {
  if (percent >= 95) return 'critical'
  if (percent >= 80) return 'warning'
  return 'normal'
}

const refreshChart = () => {
  // Refresh chart data
  quotaStore.fetchUserQuotas()
}

onMounted(() => {
  // Initialize charts
  initializeCharts()
})

const initializeCharts = () => {
  // In a real implementation, initialize Chart.js or similar
  // For now, we'll use mock data and simple rendering
}

// Watch for data changes
watch(() => quotaStore.userQuotas, () => {
  initializeCharts()
}, { deep: true })
</script>

<style scoped lang="scss">
.usage-chart {
  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      color: #303133;
      font-size: 16px;
    }

    .chart-controls {
      display: flex;
      gap: 8px;
    }
  }

  .overview-chart {
    .usage-summary {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 15px;
      margin-bottom: 20px;

      .summary-card {
        background: white;
        padding: 15px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        display: flex;
        align-items: center;
        gap: 12px;

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
          .card-value {
            font-size: 18px;
            font-weight: bold;
            color: #303133;
            margin-bottom: 4px;
          }

          .card-label {
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }

    .usage-distribution {
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

      h4 {
        margin: 0 0 15px;
        color: #303133;
        font-size: 14px;
      }

      .distribution-chart {
        height: 200px;
        margin-bottom: 15px;
      }

      .distribution-legend {
        display: flex;
        gap: 20px;
        justify-content: center;

        .legend-item {
          display: flex;
          align-items: center;
          gap: 8px;
          font-size: 12px;
          color: #606266;

          .legend-color {
            width: 12px;
            height: 12px;
            border-radius: 2px;
          }
        }
      }
    }
  }

  .trend-chart {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    h4 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 14px;
    }

    .trend-canvas-container {
      height: 300px;
      margin-bottom: 15px;
    }

    .trend-stats {
      display: flex;
      gap: 30px;
      justify-content: center;
      font-size: 12px;
      color: #606266;

      .stat-item {
        display: flex;
        gap: 8px;

        .stat-label {
          color: #909399;
        }

        .stat-value {
          color: #303133;
          font-weight: 500;
        }
      }
    }
  }

  .ranking-chart {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    h4 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 14px;
    }

    .ranking-list {
      .ranking-item {
        display: flex;
        align-items: center;
        padding: 12px 0;
        border-bottom: 1px solid #f0f0f0;

        &:last-child {
          border-bottom: none;
        }

        .ranking-position {
          font-size: 16px;
          font-weight: bold;
          color: #409eff;
          width: 30px;
          text-align: center;
          margin-right: 15px;
        }

        .ranking-info {
          flex: 1;

          .user-name {
            font-size: 14px;
            color: #303133;
            margin-bottom: 6px;
          }

          .usage-bar {
            height: 8px;
            background: #f5f7fa;
            border-radius: 4px;
            overflow: hidden;
            margin-bottom: 4px;

            .usage-fill {
              height: 100%;
              border-radius: 4px;
              transition: width 0.3s ease;

              &.normal {
                background: #67c23a;
              }

              &.warning {
                background: #e6a23c;
              }

              &.critical {
                background: #f56c6c;
              }
            }
          }

          .usage-info {
            display: flex;
            justify-content: space-between;
            font-size: 11px;
            color: #909399;

            .usage-percent {
              font-weight: 500;
              color: #303133;
            }
          }
        }
      }
    }
  }
}
</style>
