<template>
  <div class="smart-monitor">
    <div class="monitor-header">
      <div class="header-controls">
        <el-select
          v-model="selectedDevice"
          placeholder="选择设备"
          size="small"
          style="width: 200px"
          @change="loadSMARTInfo"
        >
          <el-option
            v-for="disk in availableDisks"
            :key="disk.device"
            :label="disk.device"
            :value="disk.device"
          />
        </el-select>
        <el-button size="small" @click="refreshSMART" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button size="small" type="primary" @click="runTest('short')">
          <el-icon><VideoPlay /></el-icon>
          快速测试
        </el-button>
        <el-button size="small" type="warning" @click="runTest('long')">
          <el-icon><Timer /></el-icon>
          完整测试
        </el-button>
      </div>
    </div>

    <div class="smart-overview" v-if="smartInfo">
      <div class="overview-card" :class="getHealthClass(smartInfo.overallHealth)">
        <div class="card-icon">
          <el-icon><Monitor /></el-icon>
        </div>
        <div class="card-content">
          <div class="card-label">健康状态</div>
          <div class="card-value">{{ getHealthLabel(smartInfo.overallHealth) }}</div>
        </div>
      </div>

      <div class="overview-card">
        <div class="card-icon">
          <el-icon><Sunny /></el-icon>
        </div>
        <div class="card-content">
          <div class="card-label">温度</div>
          <div class="card-value">{{ smartInfo.temperature }}°C</div>
        </div>
      </div>

      <div class="overview-card">
        <div class="card-icon">
          <el-icon><Clock /></el-icon>
        </div>
        <div class="card-content">
          <div class="card-label">运行时间</div>
          <div class="card-value">{{ formatHours(smartInfo.powerOnHours) }}小时</div>
        </div>
      </div>

      <div class="overview-card" :class="{ warning: smartInfo.errorLog?.critical }">
        <div class="card-icon">
          <el-icon><Warning /></el-icon>
        </div>
        <div class="card-content">
          <div class="card-label">错误记录</div>
          <div class="card-value">{{ smartInfo.errorLog?.count || 0 }}</div>
        </div>
      </div>
    </div>

    <!-- SMART Attributes -->
    <div class="attributes-section">
      <h3>SMART属性</h3>
      <el-table :data="smartInfo?.attributes" size="small" max-height="400">
        <el-table-column prop="id" label="ID" width="60" align="center" />
        <el-table-column prop="name" label="属性名称" min-width="200" />
        <el-table-column prop="value" label="当前值" width="100" align="center" />
        <el-table-column prop="worst" label="最差值" width="100" align="center" />
        <el-table-column prop="threshold" label="阈值" width="100" align="center" />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Test Results -->
    <div class="test-section" v-if="smartInfo?.lastTest">
      <h3>测试结果</h3>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="测试类型">
          <el-tag size="small">{{ smartInfo.lastTest.type }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getTestStatusType(smartInfo.lastTest.status)" size="small">
            {{ getTestStatusLabel(smartInfo.lastTest.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="进度" v-if="smartInfo.lastTest.status === 'running'">
          <el-progress
            :percentage="smartInfo.lastTest.progress"
            :show-text="false"
            :stroke-width="8"
          />
          {{ smartInfo.lastTest.progress }}% (剩余: {{ smartInfo.lastTest.remaining }}分钟)
        </el-descriptions-item>
        <el-descriptions-item label="结果" v-if="smartInfo.lastTest.result">
          {{ smartInfo.lastTest.result }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { diskAPI } from '@/api/disk'
import { ElMessage } from 'element-plus'
import type { DiskInfo } from '@/types/disk'
import { Refresh, Monitor, Sunny, Clock, Warning, VideoPlay, Timer } from '@element-plus/icons-vue'

const selectedDevice = ref('')
const loading = ref(false)
const smartInfo = ref<any>(null)

// Mock data
const availableDisks = ref<DiskInfo[]>([
  { device: '/dev/sdb', model: 'Samsung SSD', size: 512 * 1024 * 1024 * 1024, health: 'good', temperature: 35, partitions: [] },
  { device: '/dev/sdc', model: 'Western Digital', size: 1024 * 1024 * 1024 * 1024, health: 'good', temperature: 40, partitions: [] }
])

const loadSMARTInfo = async () => {
  if (!selectedDevice.value) return

  loading.value = true
  try {
    // In real implementation, call API
    // smartInfo.value = await diskAPI.getSMARTInfo(selectedDevice.value)
    ElMessage.success('SMART信息已更新')
  } catch (error) {
    ElMessage.error('获取SMART信息失败')
  } finally {
    loading.value = false
  }
}

const refreshSMART = () => {
  loadSMARTInfo()
}

const runTest = async (type: 'short' | 'long') => {
  if (!selectedDevice.value) {
    ElMessage.warning('请先选择设备')
    return
  }

  try {
    await diskAPI.runSMARTTest(selectedDevice.value, type)
    ElMessage.success(`${type === 'short' ? '快速' : '完整'}测试已启动`)
    await loadSMARTInfo()
  } catch (error) {
    ElMessage.error('启动测试失败')
  }
}

const getHealthClass = (health: string): string => {
  return `health-${health}`
}

const getHealthLabel = (health: string): string => {
  const labels = { good: '良好', warning: '警告', failed: '故障' }
  return labels[health] || health
}

const getStatusType = (status: string): string => {
  const types = { ok: 'success', warning: 'warning', failed: 'danger' }
  return types[status] || 'info'
}

const getStatusLabel = (status: string): string => {
  const labels = { ok: '正常', warning: '警告', failed: '故障' }
  return labels[status] || status
}

const getTestStatusType = (status: string): string => {
  const types = { running: 'warning', completed: 'success', failed: 'danger' }
  return types[status] || 'info'
}

const getTestStatusLabel = (status: string): string => {
  const labels = { running: '运行中', completed: '完成', failed: '失败' }
  return labels[status] || status
}

const formatHours = (hours: number): string => {
  if (hours < 24) return `${hours}小时`
  const days = Math.floor(hours / 24)
  const remainingHours = hours % 24
  return `${days}天${remainingHours}小时`
}

onMounted(() => {
  if (availableDisks.value.length > 0) {
    selectedDevice.value = availableDisks.value[0].device
    loadSMARTInfo()
  }
})
</script>

<style scoped lang="scss">
.smart-monitor {
  .monitor-header {
    margin-bottom: 20px;
    padding: 15px;
    background: #f5f7fa;
    border-radius: 4px;

    .header-controls {
      display: flex;
      gap: 10px;
    }
  }

  .smart-overview {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 15px;
    margin-bottom: 20px;

    .overview-card {
      background: white;
      padding: 15px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      display: flex;
      align-items: center;
      gap: 12px;

      &.health-good { border-left: 3px solid #67c23a; }
      &.health-warning { border-left: 3px solid #e6a23c; }
      &.health-failed { border-left: 3px solid #f56c6c; }
      &.warning { border-left: 3px solid #e6a23c; }

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
        .card-label {
          font-size: 12px;
          color: #909399;
          margin-bottom: 4px;
        }

        .card-value {
          font-size: 16px;
          font-weight: bold;
          color: #303133;
        }
      }
    }
  }

  .attributes-section,
  .test-section {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;

    h3 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 14px;
    }
  }
}
</style>