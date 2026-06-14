<template>
  <div class="log-viewer">
    <!-- 头部工具栏 -->
    <div class="lv-header">
      <div class="header-left">
        <h1>系统日志</h1>
        <p class="subtitle">查看和分析系统运行日志</p>
      </div>

      <div class="header-right">
        <div class="filter-controls">
          <select v-model="selectedLogType" class="filter-select">
            <option value="all">所有日志</option>
            <option value="system">系统日志</option>
            <option value="application">应用日志</option>
            <option value="security">安全日志</option>
            <option value="network">网络日志</option>
          </select>

          <select v-model="selectedLogLevel" class="filter-select">
            <option value="all">所有级别</option>
            <option value="error">错误</option>
            <option value="warning">警告</option>
            <option value="info">信息</option>
            <option value="debug">调试</option>
          </select>

          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索日志..."
            class="search-input"
          />
        </div>

        <div class="action-buttons">
          <button class="action-btn" @click="refreshLogs">
            <ArrowPathIcon class="w-4 h-4" />
            刷新
          </button>
          <button class="action-btn" @click="exportLogs">
            <ArrowDownTrayIcon class="w-4 h-4" />
            导出
          </button>
          <button class="action-btn danger" @click="clearLogs">
            <TrashIcon class="w-4 h-4" />
            清除
          </button>
        </div>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="lv-stats">
      <div class="stat-item">
        <div class="stat-icon error">
          <XCircleIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ errorCount }}</div>
          <div class="stat-label">错误</div>
        </div>
      </div>

      <div class="stat-item">
        <div class="stat-icon warning">
          <ExclamationTriangleIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ warningCount }}</div>
          <div class="stat-label">警告</div>
        </div>
      </div>

      <div class="stat-item">
        <div class="stat-icon info">
          <InformationCircleIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ infoCount }}</div>
          <div class="stat-label">信息</div>
        </div>
      </div>

      <div class="stat-item">
        <div class="stat-icon debug">
          <BugIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ debugCount }}</div>
          <div class="stat-label">调试</div>
        </div>
      </div>
    </div>

    <!-- 日志列表 -->
    <div class="lv-content">
      <div class="log-list">
        <div
          v-for="log in filteredLogs"
          :key="log.id"
          class="log-entry"
          :class="['level-' + log.level, { expanded: expandedLogs.includes(log.id) }]"
          @click="toggleExpand(log.id)"
        >
          <div class="log-summary">
            <div class="log-indicator" :class="log.level">
              <XCircleIcon v-if="log.level === 'error'" class="w-4 h-4" />
              <ExclamationTriangleIcon v-if="log.level === 'warning'" class="w-4 h-4" />
              <InformationCircleIcon v-if="log.level === 'info'" class="w-4 h-4" />
              <BugIcon v-if="log.level === 'debug'" class="w-4 h-4" />
            </div>

            <div class="log-time">{{ formatTime(log.timestamp) }}</div>

            <div class="log-source">
              <span class="source-type">{{ log.type }}</span>
              <span class="source-separator">›</span>
              <span class="source-component">{{ log.component }}</span>
            </div>

            <div class="log-message">{{ log.message }}</div>

            <button class="expand-btn">
              <ChevronDownIcon class="w-4 h-4" />
            </button>
          </div>

          <div v-if="expandedLogs.includes(log.id)" class="log-details">
            <div class="detail-row">
              <span class="detail-label">时间戳:</span>
              <span class="detail-value">{{ log.timestamp.toISOString() }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">级别:</span>
              <span class="detail-value">{{ log.level.toUpperCase() }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">类型:</span>
              <span class="detail-value">{{ log.type }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">组件:</span>
              <span class="detail-value">{{ log.component }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">进程ID:</span>
              <span class="detail-value">{{ log.pid }}</span>
            </div>
            <div v-if="log.details" class="detail-row">
              <span class="detail-label">详细信息:</span>
              <pre class="detail-content">{{ log.details }}</pre>
            </div>
            <div v-if="log.stack" class="detail-row">
              <span class="detail-label">堆栈跟踪:</span>
              <pre class="stack-trace">{{ log.stack }}</pre>
            </div>
          </div>
        </div>

        <div v-if="filteredLogs.length === 0" class="empty-state">
          <DocumentIcon class="w-12 h-12" />
          <p>暂无日志记录</p>
        </div>
      </div>
    </div>

    <!-- 分页控制 -->
    <div class="lv-pagination">
      <div class="pagination-info">
        显示 {{ (currentPage - 1) * pageSize + 1 }} - {{ Math.min(currentPage * pageSize, filteredLogs.length) }} 条，共 {{ filteredLogs.length }} 条
      </div>

      <div class="pagination-controls">
        <button
          class="page-btn"
          :disabled="currentPage === 1"
          @click="currentPage--"
        >
          上一页
        </button>

        <button
          v-for="page in displayedPages"
          :key="page"
          :class="['page-btn', { active: page === currentPage }]"
          @click="currentPage = page"
        >
          {{ page }}
        </button>

        <button
          class="page-btn"
          :disabled="currentPage === totalPages"
          @click="currentPage++"
        >
          下一页
        </button>
      </div>

      <div class="page-size-selector">
        <select v-model="pageSize" class="size-select">
          <option value="25">25条/页</option>
          <option value="50">50条/页</option>
          <option value="100">100条/页</option>
          <option value="200">200条/页</option>
        </select>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  ArrowPathIcon,
  ArrowDownTrayIcon,
  TrashIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  BeakerIcon,
  ChevronDownIcon,
  DocumentIcon
} from '@heroicons/vue/24/outline'

interface LogEntry {
  id: string
  timestamp: Date
  level: 'error' | 'warning' | 'info' | 'debug'
  type: string
  component: string
  message: string
  pid?: number
  details?: string
  stack?: string
}

// 状态管理
const selectedLogType = ref('all')
const selectedLogLevel = ref('all')
const searchQuery = ref('')
const expandedLogs = ref<string[]>([])
const currentPage = ref(1)
const pageSize = ref(50)

// 模拟日志数据
const logs = ref<LogEntry[]>([
  {
    id: '1',
    timestamp: new Date('2024-06-12T10:30:45'),
    level: 'error',
    type: 'system',
    component: 'disk-manager',
    message: '磁盘 /dev/sda1 读写错误',
    pid: 1234,
    details: 'IO error: block 0x12345678',
    stack: 'Error: IO timeout\n  at DiskManager.read (/app/disk.js:45:15)\n  at async StorageManager.check (/app/storage.js:78:9)'
  },
  {
    id: '2',
    timestamp: new Date('2024-06-12T10:28:32'),
    level: 'warning',
    type: 'security',
    component: 'auth',
    message: '检测到多次登录失败',
    pid: 5678,
    details: '用户 admin 从 192.168.1.100 登录失败 5 次'
  },
  {
    id: '3',
    timestamp: new Date('2024-06-12T10:25:18'),
    level: 'info',
    type: 'application',
    component: 'backup-service',
    message: '定时备份任务完成',
    pid: 9012,
    details: '备份名称: daily-backup\n备份大小: 2.3GB\n耗时: 15分钟'
  },
  {
    id: '4',
    timestamp: new Date('2024-06-12T10:20:05'),
    level: 'error',
    type: 'network',
    component: 'smb-server',
    message: 'SMB连接中断',
    pid: 3456,
    details: '客户端 192.168.1.50 连接超时'
  },
  {
    id: '5',
    timestamp: new Date('2024-06-12T10:15:42'),
    level: 'warning',
    type: 'system',
    component: 'cpu-monitor',
    message: 'CPU使用率过高',
    pid: 7890,
    details: '当前CPU使用率: 95%'
  },
  {
    id: '6',
    timestamp: new Date('2024-06-12T10:10:28'),
    level: 'info',
    type: 'application',
    component: 'package-center',
    message: '软件包更新可用',
    pid: 1122,
    details: '发现 3 个可更新软件包'
  },
  {
    id: '7',
    timestamp: new Date('2024-06-12T10:05:15'),
    level: 'debug',
    type: 'system',
    component: 'file-indexer',
    message: '文件索引完成',
    pid: 3344,
    details: '索引文件数: 15,432\n耗时: 3.2秒'
  }
])

// 计算属性
const filteredLogs = computed(() => {
  let result = logs.value

  // 按日志类型过滤
  if (selectedLogType.value !== 'all') {
    result = result.filter(log => log.type === selectedLogType.value)
  }

  // 按日志级别过滤
  if (selectedLogLevel.value !== 'all') {
    result = result.filter(log => log.level === selectedLogLevel.value)
  }

  // 按搜索查询过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(log =>
      log.message.toLowerCase().includes(query) ||
      log.component.toLowerCase().includes(query) ||
      (log.details && log.details.toLowerCase().includes(query))
    )
  }

  return result
})

const errorCount = computed(() => logs.value.filter(l => l.level === 'error').length)
const warningCount = computed(() => logs.value.filter(l => l.level === 'warning').length)
const infoCount = computed(() => logs.value.filter(l => l.level === 'info').length)
const debugCount = computed(() => logs.value.filter(l => l.level === 'debug').length)

const totalPages = computed(() => Math.ceil(filteredLogs.value.length / pageSize.value))

const displayedPages = computed(() => {
  const pages: number[] = []
  const maxPages = 7
  const start = Math.max(1, currentPage.value - Math.floor(maxPages / 2))
  const end = Math.min(totalPages.value, start + maxPages - 1)

  for (let i = start; i <= end; i++) {
    pages.push(i)
  }

  return pages
})

// 方法
const toggleExpand = (id: string) => {
  const index = expandedLogs.value.indexOf(id)
  if (index > -1) {
    expandedLogs.value.splice(index, 1)
  } else {
    expandedLogs.value.push(id)
  }
}

const formatTime = (date: Date): string => {
  return new Intl.DateTimeFormat('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).format(date)
}

const refreshLogs = () => {
  console.log('刷新日志')
  // 实现刷新逻辑
}

const exportLogs = () => {
  console.log('导出日志')
  // 实现导出逻辑
}

const clearLogs = () => {
  if (confirm('确定要清除所有日志吗？此操作不可撤销。')) {
    logs.value = []
  }
}

onMounted(() => {
  // 初始化日志数据
})
</script>

<style scoped>
.log-viewer {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f9fafb;
}

.lv-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.header-left h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
}

.header-right {
  display: flex;
  flex-direction: column;
  gap: 12px;
  align-items: flex-end;
}

.filter-controls {
  display: flex;
  gap: 8px;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.search-input {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  min-width: 250px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f3f4f6;
}

.action-btn.danger {
  color: #ef4444;
  border-color: #fecaca;
}

.action-btn.danger:hover {
  background: #fef2f2;
}

.lv-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 12px;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.error {
  background: #fef2f2;
  color: #ef4444;
}

.stat-icon.warning {
  background: #fef3c7;
  color: #d97706;
}

.stat-icon.info {
  background: #eff6ff;
  color: #3b82f6;
}

.stat-icon.debug {
  background: #f3f4f6;
  color: #6b7280;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.lv-content {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.log-entry {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s;
}

.log-entry:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.log-entry.level-error {
  border-left: 4px solid #ef4444;
}

.log-entry.level-warning {
  border-left: 4px solid #d97706;
}

.log-entry.level-info {
  border-left: 4px solid #3b82f6;
}

.log-entry.level-debug {
  border-left: 4px solid #6b7280;
}

.log-summary {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  cursor: pointer;
}

.log-indicator {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.log-indicator.error {
  background: #fef2f2;
  color: #ef4444;
}

.log-indicator.warning {
  background: #fef3c7;
  color: #d97706;
}

.log-indicator.info {
  background: #eff6ff;
  color: #3b82f6;
}

.log-indicator.debug {
  background: #f3f4f6;
  color: #6b7280;
}

.log-time {
  font-size: 13px;
  color: #6b7280;
  font-variant-numeric: tabular-nums;
  min-width: 80px;
}

.log-source {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #6b7280;
  min-width: 150px;
}

.source-type {
  font-weight: 500;
}

.source-separator {
  opacity: 0.5;
}

.source-component {
  color: #4b5563;
}

.log-message {
  flex: 1;
  font-size: 14px;
  color: #1f2937;
}

.expand-btn {
  padding: 4px;
  border: none;
  background: transparent;
  cursor: pointer;
  color: #9ca3af;
  transition: transform 0.2s;
}

.log-entry.expanded .expand-btn {
  transform: rotate(180deg);
}

.log-details {
  padding: 0 16px 16px;
  border-top: 1px solid #f3f4f6;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
  display: flex;
  gap: 8px;
  font-size: 13px;
}

.detail-label {
  font-weight: 500;
  color: #6b7280;
  min-width: 80px;
}

.detail-value {
  color: #1f2937;
}

.detail-content,
.stack-trace {
  flex: 1;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #4b5563;
  margin: 0;
  overflow-x: auto;
}

.stack-trace {
  color: #ef4444;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px;
  color: #9ca3af;
}

.empty-state svg {
  margin-bottom: 16px;
}

.lv-pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: white;
  border-top: 1px solid #e5e7eb;
}

.pagination-info {
  font-size: 13px;
  color: #6b7280;
}

.pagination-controls {
  display: flex;
  gap: 4px;
}

.page-btn {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.page-btn:hover:not(:disabled) {
  background: #f3f4f6;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.page-size-selector select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}
</style>