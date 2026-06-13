<template>
  <div class="system-log">
    <div class="log-header">
      <div class="header-controls">
        <el-input
          v-model="searchQuery"
          placeholder="搜索日志..."
          prefix-icon="Search"
          size="small"
          style="width: 200px"
          clearable
        />
        <el-select
          v-model="levelFilter"
          placeholder="日志级别"
          size="small"
          style="width: 120px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="调试" value="debug" />
          <el-option label="信息" value="info" />
          <el-option label="警告" value="warn" />
          <el-option label="错误" value="error" />
          <el-option label="致命" value="fatal" />
        </el-select>
        <el-select
          v-model="componentFilter"
          placeholder="组件"
          size="small"
          style="width: 120px"
          clearable
        >
          <el-option label="全部" value="" />
          <el-option label="系统" value="system" />
          <el-option label="网络" value="network" />
          <el-option label="存储" value="storage" />
          <el-option label="监控" value="monitor" />
        </el-select>
        <el-button size="small" @click="refreshLogs" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button size="small" type="danger" @click="clearLogs">
          <el-icon><Delete /></el-icon>
          清除日志
        </el-button>
        <el-button size="small" @click="toggleAutoRefresh">
          <el-icon><Timer /></el-icon>
          {{ autoRefresh ? '停止刷新' : '自动刷新' }}
        </el-button>
      </div>
    </div>

    <div class="log-list">
      <el-table
        :data="paginatedLogs"
        v-loading="loading"
        size="small"
        max-height="600"
        row-class-name="log-row"
        style="width: 100%"
      >
        <el-table-column prop="timestamp" label="时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.timestamp) }}
          </template>
        </el-table-column>

        <el-table-column prop="level" label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)" size="small">
              {{ row.level.toUpperCase() }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="component" label="组件" width="100">
          <template #default="{ row }">
            {{ row.component || 'system' }}
          </template>
        </el-table-column>

        <el-table-column prop="message" label="消息" min-width="300">
          <template #default="{ row }">
            <div class="log-message">
              {{ row.message }}
              <el-button
                v-if="row.details"
                text
                size="small"
                @click="showLogDetails(row)"
              >
                查看详情
              </el-button>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button
              size="small"
              @click="copyLog(row)"
              title="复制日志"
            >
              <el-icon><DocumentCopy /></el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[50, 100, 200, 500]"
          :total="filteredLogs.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- Log Details Dialog -->
    <el-dialog
      v-model="showDetailsDialog"
      title="日志详情"
      width="700px"
    >
      <div v-if="selectedLog" class="log-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="时间戳">
            {{ formatTime(selectedLog.timestamp) }}
          </el-descriptions-item>
          <el-descriptions-item label="日志级别">
            <el-tag :type="getLevelType(selectedLog.level)" size="small">
              {{ selectedLog.level.toUpperCase() }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="组件">
            {{ selectedLog.component || 'system' }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatTime(selectedLog.createdAt) }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="log-message-detail">
          <h4>消息内容:</h4>
          <div class="message-content">{{ selectedLog.message }}</div>
        </div>

        <div v-if="selectedLog.details" class="log-details-content">
          <h4>详细信息:</h4>
          <el-input
            :model-value="selectedLog.details"
            type="textarea"
            :rows="10"
            readonly
          />
        </div>

        <div class="log-actions">
          <el-button @click="copyLog(selectedLog)" type="primary">
            <el-icon><DocumentCopy /></el-icon>
            复制日志
          </el-button>
          <el-button @click="copyFullLog(selectedLog)">
            <el-icon><CopyDocument /></el-icon>
            复制完整日志
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- Statistics Panel -->
    <div class="log-stats">
      <div class="stat-card">
        <div class="stat-label">总日志</div>
        <div class="stat-value">{{ filteredLogs.length }}</div>
      </div>
      <div class="stat-card error">
        <div class="stat-label">错误</div>
        <div class="stat-value">{{ errorLogs }}</div>
      </div>
      <div class="stat-card warning">
        <div class="stat-label">警告</div>
        <div class="stat-value">{{ warningLogs }}</div>
      </div>
      <div class="stat-card info">
        <div class="stat-label">信息</div>
        <div class="stat-value">{{ infoLogs }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useMonitorStore } from '@/stores/monitor'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { SystemLog } from '@/types/monitor'
import {
  Refresh,
  Delete,
  Timer,
  Search,
  DocumentCopy,
  CopyDocument
} from '@element-plus/icons-vue'

const monitorStore = useMonitorStore()

const loading = ref(false)
const searchQuery = ref('')
const levelFilter = ref('')
const componentFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(100)
const autoRefresh = ref(false)
const selectedLog = ref<SystemLog | null>(null)
const showDetailsDialog = ref(false)

let refreshInterval: number | null = null

// Computed
const logs = computed(() => monitorStore.logs)

const filteredLogs = computed(() => {
  let filteredLogs = logs.value

  // Apply level filter
  if (levelFilter.value) {
    filteredLogs = filteredLogs.filter(log => log.level === levelFilter.value)
  }

  // Apply component filter
  if (componentFilter.value) {
    filteredLogs = filteredLogs.filter(log => log.component === componentFilter.value)
  }

  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filteredLogs = filteredLogs.filter(log =>
      log.message.toLowerCase().includes(query) ||
      (log.details && log.details.toLowerCase().includes(query))
    )
  }

  // Sort by timestamp (descending)
  return filteredLogs.sort((a, b) =>
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  )
})

const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredLogs.value.slice(start, end)
})

const errorLogs = computed(() =>
  filteredLogs.value.filter(log => log.level === 'error' || log.level === 'fatal').length
)

const warningLogs = computed(() =>
  filteredLogs.value.filter(log => log.level === 'warn').length
)

const infoLogs = computed(() =>
  filteredLogs.value.filter(log => log.level === 'info').length
)

// Methods
const refreshLogs = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (levelFilter.value) params.level = levelFilter.value
    if (componentFilter.value) params.component = componentFilter.value

    await monitorStore.fetchLogs(params)
    ElMessage.success('日志已刷新')
  } finally {
    loading.value = false
  }
}

const clearLogs = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清除所有日志吗？此操作不可恢复。',
      '确认清除',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    const params: any = {}
    if (levelFilter.value) params.level = levelFilter.value
    if (componentFilter.value) params.component = componentFilter.value

    await monitorStore.clearLogs(params)
    ElMessage.success('日志已清除')
    await refreshLogs()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清除日志失败')
    }
  }
}

const toggleAutoRefresh = () => {
  autoRefresh.value = !autoRefresh.value

  if (autoRefresh.value) {
    refreshInterval = window.setInterval(() => {
      refreshLogs()
    }, 10000) // Refresh every 10 seconds
    ElMessage.success('已启用自动刷新')
  } else {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
    ElMessage.info('已停止自动刷新')
  }
}

const formatTime = (dateStr: string): string => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const getLevelType = (level: string): string => {
  const types: Record<string, string> = {
    debug: 'info',
    info: 'info',
    warn: 'warning',
    error: 'danger',
    fatal: 'danger'
  }
  return types[level] || 'info'
}

const showLogDetails = (log: SystemLog) => {
  selectedLog.value = log
  showDetailsDialog.value = true
}

const copyLog = async (log: SystemLog) => {
  const logText = `[${log.timestamp}] ${log.level.toUpperCase()} [${log.component || 'system'}] ${log.message}`

  try {
    await navigator.clipboard.writeText(logText)
    ElMessage.success('日志已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

const copyFullLog = async (log: SystemLog) => {
  const fullLog = `时间戳: ${log.timestamp}
级别: ${log.level.toUpperCase()}
组件: ${log.component || 'system'}
消息: ${log.message}
${log.details ? `详细信息: ${log.details}` : ''}
创建时间: ${log.createdAt}`

  try {
    await navigator.clipboard.writeText(fullLog)
    ElMessage.success('完整日志已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
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
  refreshLogs()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped lang="scss">
.system-log {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.log-header {
  margin-bottom: 15px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;

  .header-controls {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }
}

.log-list {
  flex: 1;
  overflow: auto;

  :deep(.el-table .log-row) {
    cursor: pointer;

    &:hover {
      background-color: #f5f7fa;
    }
  }

  .log-message {
    display: flex;
    align-items: center;
    gap: 10px;
    word-break: break-all;
  }

  .pagination {
    margin-top: 15px;
    display: flex;
    justify-content: center;
  }
}

.log-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
  margin-top: 15px;

  .stat-card {
    background: white;
    padding: 15px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    text-align: center;

    .stat-label {
      font-size: 12px;
      color: #909399;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 24px;
      font-weight: bold;
      color: #303133;
    }

    &.error .stat-value { color: #f56c6c; }
    &.warning .stat-value { color: #e6a23c; }
    &.info .stat-value { color: #409eff; }
  }
}

.log-details {
  .log-message-detail,
  .log-details-content,
  .log-actions {
    margin-top: 20px;

    h4 {
      margin: 0 0 10px;
      color: #303133;
      font-size: 14px;
    }
  }

  .message-content {
    padding: 12px;
    background: #f5f7fa;
    border-radius: 4px;
    color: #303133;
    line-height: 1.6;
  }

  .log-actions {
    display: flex;
    gap: 10px;
  }
}
</style>