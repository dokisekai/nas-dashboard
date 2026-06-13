<template>
  <div class="benchmark-tool">
    <div class="tool-header">
      <div class="header-controls">
        <el-select
          v-model="selectedDevice"
          placeholder="选择设备"
          size="small"
          style="width: 200px"
        >
          <el-option
            v-for="disk in availableDisks"
            :key="disk.device"
            :label="disk.device"
            :value="disk.device"
          />
        </el-select>
        <el-button
          size="small"
          type="primary"
          @click="startBenchmark"
          :loading="benchmarking"
          :disabled="!selectedDevice"
        >
          <el-icon><Timer /></el-icon>
          开始测试
        </el-button>
        <el-button size="small" @click="exportResults" :disabled="!hasResults">
          <el-icon><Download /></el-icon>
          导出结果
        </el-button>
      </div>
    </div>

    <!-- Benchmark Progress -->
    <div v-if="benchmarking" class="benchmark-progress">
      <h3>磁盘性能测试进行中...</h3>
      <el-progress
        :percentage="progress"
        :status="progressStatus"
        :show-text="true"
        :stroke-width="20"
      />
      <div class="progress-info">
        <span>{{ currentOperation }}</span>
        <span>{{ progress }}%</span>
      </div>
    </div>

    <!-- Benchmark Results -->
    <div v-if="hasResults && !benchmarking" class="benchmark-results">
      <div class="results-overview">
        <div class="result-card read-speed">
          <div class="card-icon"><el-icon><Download /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ formatSpeed(currentResult?.readSpeed || 0) }}</div>
            <div class="card-label">读取速度</div>
          </div>
        </div>
        <div class="result-card write-speed">
          <div class="card-icon"><el-icon><Upload /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ formatSpeed(currentResult?.writeSpeed || 0) }}</div>
            <div class="card-label">写入速度</div>
          </div>
        </div>
        <div class="result-card read-iops">
          <div class="card-icon"><el-icon><DataLine /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ currentResult?.readIOPS?.toLocaleString() || 0 }}</div>
            <div class="card-label">读取IOPS</div>
          </div>
        </div>
        <div class="result-card write-iops">
          <div class="card-icon"><el-icon><DataBoard /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ currentResult?.writeIOPS?.toLocaleString() || 0 }}</div>
            <div class="card-label">写入IOPS</div>
          </div>
        </div>
        <div class="result-card access-time">
          <div class="card-icon"><el-icon><Clock /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ (currentResult?.accessTime || 0).toFixed(2) }}ms</div>
            <div class="card-label">访问时间</div>
          </div>
        </div>
        <div class="result-card device">
          <div class="card-icon"><el-icon><Document /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ selectedDevice }}</div>
            <div class="card-label">测试设备</div>
          </div>
        </div>
      </div>

      <!-- Performance Rating -->
      <div class="performance-rating">
        <h3>性能评级</h3>
        <el-rate
          v-model="performanceRating"
          disabled
          show-score
          text-color="#ff9900"
          score-template="{value}"
        />
        <div class="rating-description">
          {{ getRatingDescription(performanceRating) }}
        </div>
      </div>

      <!-- Detailed Results -->
      <div class="detailed-results">
        <h3>详细测试结果</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="测试设备">
            {{ selectedDevice }}
          </el-descriptions-item>
          <el-descriptions-item label="测试时间">
            {{ currentResult?.timestamp ? formatTime(currentResult.timestamp) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="读取速度">
            {{ formatSpeed(currentResult?.readSpeed || 0) }}
          </el-descriptions-item>
          <el-descriptions-item label="写入速度">
            {{ formatSpeed(currentResult?.writeSpeed || 0) }}
          </el-descriptions-item>
          <el-descriptions-item label="随机读取IOPS">
            {{ currentResult?.readIOPS?.toLocaleString() || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="随机写入IOPS">
            {{ currentResult?.writeIOPS?.toLocaleString() || 0 }}
          </el-descriptions-item>
          <el-descriptions-item label="平均访问时间">
            {{ (currentResult?.accessTime || 0).toFixed(2) }}ms
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Performance Comparison -->
      <div class="performance-comparison">
        <h3>性能对比</h3>
        <div class="comparison-chart">
          <div class="chart-item">
            <div class="chart-label">读取速度</div>
            <div class="chart-bar">
              <div
                class="bar-fill"
                :style="{ width: `${(currentResult?.readSpeed || 0) / (maxSpeed / 1024 / 1024) * 100}%` }"
              />
            </div>
            <div class="chart-value">{{ formatSpeed(currentResult?.readSpeed || 0) }}</div>
          </div>
          <div class="chart-item">
            <div class="chart-label">写入速度</div>
            <div class="chart-bar">
              <div
                class="bar-fill"
                :style="{ width: `${(currentResult?.writeSpeed || 0) / (maxSpeed / 1024 / 1024) * 100}%` }"
              />
            </div>
            <div class="chart-value">{{ formatSpeed(currentResult?.writeSpeed || 0) }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Configuration -->
    <div class="test-config">
      <h3>测试配置</h3>
      <el-form label-width="120px">
        <el-form-item label="测试大小">
          <el-select v-model="testConfig.size">
            <el-option label="1GB" value="1024" />
            <el-option label="4GB" value="4096" />
            <el-option label="16GB" value="16384" />
            <el-option label="64GB" value="65536" />
          </el-select>
          <span class="unit-hint">MB</span>
        </el-form-item>
        <el-form-item label="队列深度">
          <el-select v-model="testConfig.queueDepth">
            <el-option label="1" :value="1" />
            <el-option label="32" :value="32" />
            <el-option label="64" :value="64" />
            <el-option label="128" :value="128" />
            <el-option label="256" :value="256" />
          </el-select>
        </el-form-item>
        <el-form-item label="测试时长">
          <el-select v-model="testConfig.duration">
            <el-option label="30秒" :value="30" />
            <el-option label="60秒" :value="60" />
            <el-option label="120秒" :value="120" />
            <el-option label="300秒" :value="300" />
          </el-select>
          <span class="unit-hint">秒</span>
        </el-form-item>
        <el-form-item label="测试模式">
          <el-select v-model="testConfig.mode">
            <el-option label="顺序读写" value="sequential" />
            <el-option label="随机读写" value="random" />
            <el-option label="混合模式" value="mixed" />
          </el-select>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { diskAPI } from '@/api/disk'
import { ElMessage } from 'element-plus'
import type { DiskBenchmark } from '@/types/disk'
import { Timer, Download, Upload, DataLine, DataBoard, Clock, Document } from '@element-plus/icons-vue'

const selectedDevice = ref('')
const benchmarking = ref(false)
const progress = ref(0)
const progressStatus = ref('')
const currentOperation = ref('')
const currentResult = ref<DiskBenchmark | null>(null)
const hasResults = ref(false)

const testConfig = ref({
  size: 4096,
  queueDepth: 32,
  duration: 60,
  mode: 'mixed'
})

const availableDisks = ref([
  { device: '/dev/sdb', model: 'Samsung SSD', size: 512 * 1024 * 1024 * 1024, health: 'good', temperature: 35, partitions: [] },
  { device: '/dev/sdc', model: 'Western Digital', size: 1024 * 1024 * 1024 * 1024, health: 'good', temperature: 40, partitions: [] }
])

const maxSpeed = ref(1000) // 1000 MB/s reference

const performanceRating = computed(() => {
  if (!currentResult.value) return 0
  const speed = (currentResult.value.readSpeed + currentResult.value.writeSpeed) / 2
  if (speed < 100) return 1
  if (speed < 200) return 2
  if (speed < 400) return 3
  if (speed < 600) return 4
  return 5
})

const startBenchmark = async () => {
  if (!selectedDevice.value) {
    ElMessage.warning('请先选择设备')
    return
  }

  benchmarking.value = true
  progress.value = 0
  currentOperation.value = '准备测试环境...'

  try {
    // Simulate benchmark progress
    const steps = [
      { operation: '准备测试环境...', progress: 10 },
      { operation: '执行读取测试...', progress: 30 },
      { operation: '执行写入测试...', progress: 60 },
      { operation: '执行IOPS测试...', progress: 80 },
      { operation: '计算测试结果...', progress: 95 },
      { operation: '完成测试...', progress: 100 }
    ]

    for (const step of steps) {
      currentOperation.value = step.operation
      progress.value = step.progress
      await new Promise(resolve => setTimeout(resolve, 2000))
    }

    // Mock result
    currentResult.value = {
      device: selectedDevice.value,
      readSpeed: 450 * 1024 * 1024, // 450 MB/s
      writeSpeed: 320 * 1024 * 1024, // 320 MB/s
      readIOPS: 85000,
      writeIOPS: 45000,
      accessTime: 12.5,
      timestamp: new Date().toISOString()
    }

    hasResults.value = true
    ElMessage.success('性能测试完成')
  } catch (error) {
    ElMessage.error('性能测试失败')
  } finally {
    benchmarking.value = false
  }
}

const exportResults = () => {
  if (!currentResult.value) return

  const results = {
    device: currentResult.value.device,
    timestamp: currentResult.value.timestamp,
    readSpeed: `${(currentResult.value.readSpeed / 1024 / 1024).toFixed(2)} MB/s`,
    writeSpeed: `${(currentResult.value.writeSpeed / 1024 / 1024).toFixed(2)} MB/s`,
    readIOPS: currentResult.value.readIOPS,
    writeIOPS: currentResult.value.writeIOPS,
    accessTime: `${currentResult.value.accessTime.toFixed(2)} ms`,
    testConfig: testConfig.value
  }

  const blob = new Blob([JSON.stringify(results, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `benchmark_${selectedDevice.value}_${Date.now()}.json`
  link.click()
  URL.revokeObjectURL(url)

  ElMessage.success('测试结果已导出')
}

const formatSpeed = (bytesPerSecond: number): string => {
  if (bytesPerSecond === 0) return '0 MB/s'
  const mbPerSecond = bytesPerSecond / 1024 / 1024
  return `${mbPerSecond.toFixed(1)} MB/s`
}

const formatTime = (timestamp: string): string => {
  return new Date(timestamp).toLocaleString('zh-CN')
}

const getRatingDescription = (rating: number): string => {
  const descriptions = [
    '性能较差 - 建议更换',
    '性能一般 - 满足基本需求',
    '性能良好 - 适合日常使用',
    '性能优秀 - 适合专业应用',
    '性能卓越 - 适合企业级应用'
  ]
  return descriptions[rating - 1] || ''
}
</script>

<style scoped lang="scss">
.benchmark-tool {
  .tool-header {
    margin-bottom: 20px;
    padding: 15px;
    background: #f5f7fa;
    border-radius: 4px;

    .header-controls {
      display: flex;
      gap: 10px;
    }
  }

  .benchmark-progress {
    background: white;
    padding: 30px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    margin-bottom: 20px;

    h3 {
      margin: 0 0 20px;
      color: #303133;
      text-align: center;
    }

    .progress-info {
      display: flex;
      justify-content: space-between;
      margin-top: 10px;
      font-size: 12px;
      color: #606266;
    }
  }

  .benchmark-results {
    .results-overview {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 15px;
      margin-bottom: 20px;

      .result-card {
        background: white;
        padding: 20px;
        border-radius: 8px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        display: flex;
        align-items: center;
        gap: 15px;

        .card-icon {
          font-size: 28px;
          color: #409eff;
          display: flex;
          align-items: center;
          justify-content: center;
          width: 50px;
          height: 50px;
          border-radius: 50%;
          background: #ecf5ff;
        }

        .card-content {
          .card-value {
            font-size: 20px;
            font-weight: bold;
            color: #303133;
            margin-bottom: 5px;
          }

          .card-label {
            font-size: 12px;
            color: #909399;
          }
        }
      }
    }

    .performance-rating,
    .detailed-results,
    .performance-comparison {
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      margin-bottom: 20px;

      h3 {
        margin: 0 0 15px;
        color: #303133;
        font-size: 16px;
      }
    }

    .rating-description {
      margin-top: 10px;
      font-size: 14px;
      color: #606266;
    }

    .comparison-chart {
      margin-top: 20px;

      .chart-item {
        margin-bottom: 20px;

        .chart-label {
          margin-bottom: 8px;
          font-size: 12px;
          color: #606266;
        }

        .chart-bar {
          height: 24px;
          background: #f5f7fa;
          border-radius: 12px;
          overflow: hidden;
          margin-bottom: 5px;

          .bar-fill {
            height: 100%;
            background: linear-gradient(90deg, #409eff, #67c23a);
            transition: width 1s ease;
          }
        }

        .chart-value {
          font-size: 14px;
          color: #303133;
          font-weight: 500;
        }
      }
    }
  }

  .test-config {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    h3 {
      margin: 0 0 20px;
      color: #303133;
      font-size: 16px;
    }

    .unit-hint {
      margin-left: 8px;
      color: #909399;
      font-size: 12px;
    }
  }
}
</style>