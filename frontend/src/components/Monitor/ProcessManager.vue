<template>
  <div class="process-manager">
    <div class="manager-header">
      <div class="header-controls">
        <el-input
          v-model="searchQuery"
          placeholder="搜索进程..."
          prefix-icon="Search"
          size="small"
          style="width: 200px"
          clearable
        />
        <el-select
          v-model="statusFilter"
          placeholder="状态筛选"
          size="small"
          style="width: 120px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="运行中" value="running" />
          <el-option label="睡眠" value="sleeping" />
          <el-option label="停止" value="stopped" />
        </el-select>
        <el-button size="small" @click="refreshProcesses">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <div class="process-stats">
        <span class="stat-item">总进程: {{ filteredProcesses.length }}</span>
        <span class="stat-item">CPU: {{ totalCPU.toFixed(1) }}%</span>
        <span class="stat-item">内存: {{ totalMemory.toFixed(1) }}%</span>
      </div>
    </div>

    <div class="process-list">
      <el-table
        :data="paginatedProcesses"
        v-loading="loading"
        size="small"
        max-height="600"
        @row-click="showProcessDetails"
        style="width: 100%; cursor: pointer"
      >
        <el-table-column prop="pid" label="PID" width="80" />
        <el-table-column prop="name" label="进程名称" min-width="150">
          <template #default="{ row }">
            <div class="process-name">
              <el-icon class="process-icon"><Document /></el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="cpuPercent" label="CPU %" width="100">
          <template #default="{ row }">
            <div class="cpu-cell">
              <el-progress
                :percentage="row.cpuPercent"
                :color="getCPUCColor(row.cpuPercent)"
                :show-text="false"
                :stroke-width="6"
              />
              <span class="cpu-text">{{ row.cpuPercent.toFixed(1) }}%</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="memoryPercent" label="内存 %" width="100">
          <template #default="{ row }">
            <div class="memory-cell">
              <el-progress
                :percentage="row.memoryPercent"
                :color="getMemoryColor(row.memoryPercent)"
                :show-text="false"
                :stroke-width="6"
              />
              <span class="memory-text">{{ row.memoryPercent.toFixed(1) }}%</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="memory" label="内存使用" width="120">
          <template #default="{ row }">
            {{ formatSize(row.memory) }}
          </template>
        </el-table-column>
        <el-table-column prop="threads" label="线程" width="80" align="center" />
        <el-table-column prop="username" label="用户" width="100" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click.stop="showProcessDetails(row)">
                <el-icon><View /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click.stop="killProcess(row, 15)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100, 200]"
          :total="filteredProcesses.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- Process Details Dialog -->
    <el-dialog
      v-model="showDetailsDialog"
      :title="`进程详情 - ${selectedProcess?.name || ''} (PID: ${selectedProcess?.pid || ''})`"
      width="800px"
    >
      <div v-if="selectedProcess" class="process-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="进程ID">
            {{ selectedProcess.pid }}
          </el-descriptions-item>
          <el-descriptions-item label="进程名称">
            {{ selectedProcess.name }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(selectedProcess.status)" size="small">
              {{ getStatusLabel(selectedProcess.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="用户">
            {{ selectedProcess.username }}
          </el-descriptions-item>
          <el-descriptions-item label="CPU 使用">
            {{ selectedProcess.cpuPercent.toFixed(2) }}%
          </el-descriptions-item>
          <el-descriptions-item label="内存使用">
            {{ selectedProcess.memoryPercent.toFixed(2) }}% ({{ formatSize(selectedProcess.memory) }})
          </el-descriptions-item>
          <el-descriptions-item label="线程数">
            {{ selectedProcess.threads }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatTime(selectedProcess.createdTime) }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="command-info">
          <h4>命令行:</h4>
          <el-input
            :model-value="selectedProcess.command"
            type="textarea"
            :rows="3"
            readonly
          />
        </div>

        <div class="process-actions">
          <h4>进程操作:</h4>
          <div class="action-buttons">
            <el-button @click="killProcess(selectedProcess, 15)" type="warning">
              终止进程 (SIGTERM)
            </el-button>
            <el-button @click="killProcess(selectedProcess, 9)" type="danger">
              强制终止 (SIGKILL)
            </el-button>
            <el-button @click="killProcess(selectedProcess, 1)" type="info">
              挂起进程 (SIGHUP)
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMonitorStore } from '@/stores/monitor'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { ProcessInfo } from '@/types/monitor'
import {
  Refresh,
  Search,
  Document,
  View,
  Delete
} from '@element-plus/icons-vue'

const monitorStore = useMonitorStore()

const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(50)
const selectedProcess = ref<ProcessInfo | null>(null)
const showDetailsDialog = ref(false)

// Computed
const filteredProcesses = computed(() => {
  let processes = monitorStore.processes

  // Apply status filter
  if (statusFilter.value) {
    processes = processes.filter(p => p.status === statusFilter.value)
  }

  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    processes = processes.filter(p =>
      p.name.toLowerCase().includes(query) ||
      p.command.toLowerCase().includes(query) ||
      p.pid.toString().includes(query)
    )
  }

  // Sort by CPU usage (descending)
  return processes.sort((a, b) => b.cpuPercent - a.cpuPercent)
})

const paginatedProcesses = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredProcesses.value.slice(start, end)
})

const totalCPU = computed(() =>
  filteredProcesses.value.reduce((sum, p) => sum + p.cpuPercent, 0)
)

const totalMemory = computed(() =>
  filteredProcesses.value.reduce((sum, p) => sum + p.memoryPercent, 0)
)

// Methods
const refreshProcesses = async () => {
  loading.value = true
  try {
    await monitorStore.fetchProcesses()
    ElMessage.success('进程列表已刷新')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string): string => {
  const types: Record<string, string> = {
    running: 'success',
    sleeping: 'info',
    stopped: 'danger',
    zombie: 'warning',
    unknown: 'info'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    running: '运行中',
    sleeping: '睡眠',
    stopped: '停止',
    zombie: '僵尸',
    unknown: '未知'
  }
  return labels[status] || status
}

const getCPUCColor = (percent: number): string => {
  if (percent >= 80) return '#f56c6c'
  if (percent >= 50) return '#e6a23c'
  return '#67c23a'
}

const getMemoryColor = (percent: number): string => {
  if (percent >= 80) return '#f56c6c'
  if (percent >= 50) return '#e6a23c'
  return '#409eff'
}

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const formatTime = (timestamp: number): string => {
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

const showProcessDetails = (process: ProcessInfo) => {
  selectedProcess.value = process
  showDetailsDialog.value = true
}

const killProcess = async (process: ProcessInfo, signal: number) => {
  const signalNames: Record<number, string> = {
    1: 'SIGHUP',
    15: 'SIGTERM',
    9: 'SIGKILL'
  }

  try {
    await ElMessageBox.confirm(
      `确定要发送 ${signalNames[signal]} 信号到进程 "${process.name}" (PID: ${process.pid}) 吗？`,
      '确认终止进程',
      {
        type: signal === 9 ? 'error' : 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    await monitorStore.killProcess(process.pid, signal.toString())
    ElMessage.success(`进程 ${process.pid} 已终止`)
    showDetailsDialog.value = false
    await refreshProcesses()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('终止进程失败')
    }
  }
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handlePageChange = (page: number) => {
  currentPage.value = page
}

// Lifecycle
onMounted(() => {
  refreshProcesses()
})
</script>

<style scoped lang="scss">
.process-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.manager-header {
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

  .process-stats {
    display: flex;
    gap: 20px;
    font-size: 12px;
    color: #606266;

    .stat-item {
      display: flex;
      align-items: center;
      gap: 5px;
    }
  }
}

.process-list {
  flex: 1;
  overflow: auto;

  .process-name {
    display: flex;
    align-items: center;
    gap: 8px;

    .process-icon {
      font-size: 14px;
      color: #409eff;
    }
  }

  .cpu-cell,
  .memory-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .cpu-text,
    .memory-text {
      font-size: 11px;
      color: #606266;
      min-width: 40px;
    }
  }

  .pagination {
    margin-top: 15px;
    display: flex;
    justify-content: center;
  }
}

.process-details {
  .command-info,
  .process-actions {
    margin-top: 20px;

    h4 {
      margin: 0 0 10px;
      color: #303133;
      font-size: 14px;
    }
  }

  .action-buttons {
    display: flex;
    gap: 10px;
  }
}
</style>