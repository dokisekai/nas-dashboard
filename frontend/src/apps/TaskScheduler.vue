<template>
  <div class="task-scheduler">
    <!-- 头部 -->
    <div class="ts-header">
      <div class="header-left">
        <h1>任务计划</h1>
        <p class="subtitle">管理系统定时任务和备份计划</p>
      </div>

      <div class="header-right">
        <button class="action-btn primary" @click="createNewTask">
          <PlusIcon class="w-4 h-4" />
          新建任务
        </button>
        <button class="action-btn" @click="runAllTasks">
          <PlayIcon class="w-4 h-4" />
          运行所有任务
        </button>
      </div>
    </div>

    <!-- 任务统计 -->
    <div class="ts-stats">
      <div class="stat-card">
        <div class="stat-icon active">
          <CheckCircleIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ activeTasks }}</div>
          <div class="stat-label">活动任务</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon success">
          <CheckBadgeIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ completedToday }}</div>
          <div class="stat-label">今日完成</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon warning">
          <ClockIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ pendingTasks }}</div>
          <div class="stat-label">待运行</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon error">
          <XCircleIcon class="w-5 h-5" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ failedTasks }}</div>
          <div class="stat-label">失败任务</div>
        </div>
      </div>
    </div>

    <!-- 主要内容 -->
    <div class="ts-content">
      <!-- 左侧：任务列表 -->
      <div class="ts-tasks">
        <div class="section-header">
          <h2>任务列表</h2>
          <div class="filter-controls">
            <select v-model="selectedFilter" class="filter-select">
              <option value="all">所有任务</option>
              <option value="active">活动任务</option>
              <option value="scheduled">已计划</option>
              <option value="running">运行中</option>
              <option value="failed">失败</option>
            </select>
          </div>
        </div>

        <div class="task-list">
          <div
            v-for="task in filteredTasks"
            :key="task.id"
            class="task-item"
            :class="[
              task.status,
              { selected: selectedTask === task.id }
            ]"
            @click="selectTask(task)"
          >
            <div class="task-icon" :class="task.status">
              <CloudIcon v-if="task.type === 'backup'" class="w-5 h-5" />
              <ServerIcon v-else-if="task.type === 'system'" class="w-5 h-5" />
              <FolderIcon v-else-if="task.type === 'cleanup'" class="w-5 h-5" />
              <CameraIcon v-else-if="task.type === 'snapshot'" class="w-5 h-5" />
              <CogIcon v-else class="w-5 h-5" />
            </div>

            <div class="task-info">
              <div class="task-name">{{ task.name }}</div>
              <div class="task-schedule">{{ task.schedule }}</div>
              <div class="task-status" :class="task.status">
                {{ getStatusText(task.status) }}
              </div>
            </div>

            <div class="task-actions">
              <button
                v-if="task.status !== 'running'"
                class="task-action-btn"
                @click.stop="runTask(task)"
                title="立即运行"
              >
                <PlayIcon class="w-4 h-4" />
              </button>
              <button
                class="task-action-btn"
                @click.stop="toggleTask(task)"
                :title="task.enabled ? '禁用' : '启用'"
              >
                <div class="toggle-icon" :class="{ active: task.enabled }">
                  <div class="toggle-dot"></div>
                </div>
              </button>
              <button
                class="task-action-btn"
                @click.stop="editTask(task)"
                title="编辑"
              >
                <PencilIcon class="w-4 h-4" />
              </button>
              <button
                class="task-action-btn danger"
                @click.stop="deleteTask(task)"
                title="删除"
              >
                <TrashIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧：任务详情 -->
      <div class="ts-details">
        <div v-if="selectedTaskData" class="task-details">
          <div class="details-header">
            <h2>{{ selectedTaskData.name }}</h2>
            <div class="task-badges">
              <span class="badge" :class="selectedTaskData.status">
                {{ getStatusText(selectedTaskData.status) }}
              </span>
              <span class="badge type">{{ selectedTaskData.type }}</span>
            </div>
          </div>

          <div class="details-section">
            <h3>任务设置</h3>
            <div class="details-grid">
              <div class="detail-item">
                <span class="detail-label">任务类型:</span>
                <span class="detail-value">{{ getTaskTypeText(selectedTaskData.type) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">计划:</span>
                <span class="detail-value">{{ selectedTaskData.schedule }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">状态:</span>
                <span class="detail-value" :class="{ enabled: selectedTaskData.enabled }">
                  {{ selectedTaskData.enabled ? '已启用' : '已禁用' }}
                </span>
              </div>
              <div class="detail-item">
                <span class="detail-label">最后运行:</span>
                <span class="detail-value">{{ formatTime(selectedTaskData.lastRun) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">下次运行:</span>
                <span class="detail-value">{{ formatTime(selectedTaskData.nextRun) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">创建时间:</span>
                <span class="detail-value">{{ formatTime(selectedTaskData.createdAt) }}</span>
              </div>
            </div>
          </div>

          <div class="details-section">
            <h3>任务描述</h3>
            <p class="task-description">{{ selectedTaskData.description }}</p>
          </div>

          <!-- 任务配置详情 -->
          <div class="details-section">
            <h3>配置详情</h3>
            <div v-if="selectedTaskData.type === 'backup'" class="config-details">
              <div class="config-item">
                <span class="config-label">备份源:</span>
                <span class="config-value">{{ selectedTaskData.config?.source || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">备份目标:</span>
                <span class="config-value">{{ selectedTaskData.config?.destination || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">压缩:</span>
                <span class="config-value">{{ selectedTaskData.config?.compress ? '是' : '否' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">保留天数:</span>
                <span class="config-value">{{ selectedTaskData.config?.retention || '-' }} 天</span>
              </div>
            </div>

            <div v-else-if="selectedTaskData.type === 'cleanup'" class="config-details">
              <div class="config-item">
                <span class="config-label">清理目标:</span>
                <span class="config-value">{{ selectedTaskData.config?.target || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">文件类型:</span>
                <span class="config-value">{{ selectedTaskData.config?.fileTypes?.join(', ') || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="config-label">保留时间:</span>
                <span class="config-value">{{ selectedTaskData.config?.keepDays || '-' }} 天</span>
              </div>
            </div>
          </div>

          <!-- 运行历史 -->
          <div class="details-section">
            <h3>运行历史</h3>
            <div class="run-history">
              <div
                v-for="run in selectedTaskData.history"
                :key="run.id"
                class="history-item"
                :class="run.status"
              >
                <div class="history-icon">
                  <CheckCircleIcon v-if="run.status === 'success'" class="w-4 h-4" />
                  <XCircleIcon v-else-if="run.status === 'failed'" class="w-4 h-4" />
                  <ClockIcon v-else class="w-4 h-4" />
                </div>
                <div class="history-info">
                  <div class="history-time">{{ formatTime(run.startTime) }}</div>
                  <div class="history-duration">耗时: {{ run.duration }}</div>
                </div>
                <div class="history-status">
                  {{ run.status === 'success' ? '成功' : run.status === 'failed' ? '失败' : '运行中' }}
                </div>
              </div>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="details-actions">
            <button class="action-btn primary" @click="runTask(selectedTaskData)">
              <PlayIcon class="w-4 h-4" />
              立即运行
            </button>
            <button class="action-btn" @click="editTask(selectedTaskData)">
              <PencilIcon class="w-4 h-4" />
              编辑
            </button>
            <button class="action-btn" @click="viewTaskLogs(selectedTaskData)">
              <DocumentIcon class="w-4 h-4" />
              查看日志
            </button>
            <button class="action-btn danger" @click="deleteTask(selectedTaskData)">
              <TrashIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>

        <div v-else class="no-selection">
          <DocumentIcon class="w-12 h-12" />
          <p>选择一个任务查看详情</p>
        </div>
      </div>
    </div>

    <!-- 新建/编辑任务模态框 -->
    <Transition name="fade">
      <div v-if="showTaskModal" class="modal-overlay" @click.self="showTaskModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>{{ editingTask ? '编辑任务' : '新建任务' }}</h3>
            <button @click="showTaskModal = false">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="modal-body">
            <div class="form-section">
              <h4>基本信息</h4>
              <div class="form-group">
                <label>任务名称</label>
                <input v-model="taskForm.name" type="text" class="form-input" placeholder="输入任务名称" />
              </div>
              <div class="form-group">
                <label>任务类型</label>
                <select v-model="taskForm.type" class="form-select">
                  <option value="backup">备份任务</option>
                  <option value="cleanup">清理任务</option>
                  <option value="system">系统任务</option>
                  <option value="snapshot">快照任务</option>
                  <option value="custom">自定义脚本</option>
                </select>
              </div>
              <div class="form-group">
                <label>描述</label>
                <textarea v-model="taskForm.description" class="form-textarea" placeholder="任务描述"></textarea>
              </div>
            </div>

            <div class="form-section">
              <h4>计划设置</h4>
              <div class="form-group">
                <label>计划类型</label>
                <select v-model="taskForm.scheduleType" class="form-select">
                  <option value="minute">每分钟</option>
                  <option value="hourly">每小时</option>
                  <option value="daily">每天</option>
                  <option value="weekly">每周</option>
                  <option value="monthly">每月</option>
                  <option value="custom">自定义 (Cron)</option>
                </select>
              </div>
              <div class="form-group">
                <label>计划时间</label>
                <input v-model="taskForm.scheduleTime" type="time" class="form-input" />
              </div>
            </div>

            <div class="form-section" v-if="taskForm.type === 'backup'">
              <h4>备份设置</h4>
              <div class="form-group">
                <label>备份源</label>
                <input v-model="taskForm.config.source" type="text" class="form-input" placeholder="/path/to/source" />
              </div>
              <div class="form-group">
                <label>备份目标</label>
                <input v-model="taskForm.config.destination" type="text" class="form-input" placeholder="/path/to/destination" />
              </div>
              <div class="form-row">
                <div class="form-group">
                  <label>压缩</label>
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="taskForm.config.compress" />
                    <span>启用压缩</span>
                  </label>
                </div>
                <div class="form-group">
                  <label>保留天数</label>
                  <input v-model="taskForm.config.retention" type="number" class="form-input" />
                </div>
              </div>
            </div>

            <div class="modal-actions">
              <button class="action-btn" @click="showTaskModal = false">取消</button>
              <button class="action-btn primary" @click="saveTask">保存</button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  PlusIcon,
  PlayIcon,
  CheckCircleIcon,
  CheckBadgeIcon,
  ClockIcon,
  XCircleIcon,
  CloudIcon,
  ServerIcon,
  FolderIcon,
  CameraIcon,
  CogIcon,
  PencilIcon,
  TrashIcon,
  DocumentIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

interface Task {
  id: string
  name: string
  type: 'backup' | 'cleanup' | 'system' | 'snapshot' | 'custom'
  status: 'active' | 'running' | 'success' | 'failed' | 'disabled'
  schedule: string
  enabled: boolean
  lastRun?: Date
  nextRun?: Date
  createdAt: Date
  description: string
  config?: any
  history: TaskRun[]
}

interface TaskRun {
  id: string
  status: 'success' | 'failed' | 'running'
  startTime: Date
  duration: string
}

// 状态管理
const selectedFilter = ref('all')
const selectedTask = ref<string | null>(null)
const showTaskModal = ref(false)
const editingTask = ref<Task | null>(null)

const taskForm = ref({
  name: '',
  type: 'backup' as const,
  description: '',
  scheduleType: 'daily',
  scheduleTime: '02:00',
  config: {
    source: '',
    destination: '',
    compress: false,
    retention: 7
  }
})

// 模拟任务数据
const tasks = ref<Task[]>([
  {
    id: '1',
    name: '每日备份',
    type: 'backup',
    status: 'active',
    schedule: '每天 02:00',
    enabled: true,
    lastRun: new Date('2024-06-12T02:00:00'),
    nextRun: new Date('2024-06-13T02:00:00'),
    createdAt: new Date('2024-06-01T10:00:00'),
    description: '备份用户数据和配置文件到外部存储',
    config: {
      source: '/home/admin/Documents',
      destination: '/external/backup/daily',
      compress: true,
      retention: 7
    },
    history: [
      {
        id: 'run1',
        status: 'success',
        startTime: new Date('2024-06-12T02:00:00'),
        duration: '15分钟'
      },
      {
        id: 'run2',
        status: 'success',
        startTime: new Date('2024-06-11T02:00:00'),
        duration: '14分钟'
      },
      {
        id: 'run3',
        status: 'failed',
        startTime: new Date('2024-06-10T02:00:00'),
        duration: '3分钟'
      }
    ]
  },
  {
    id: '2',
    name: '系统清理',
    type: 'cleanup',
    status: 'active',
    schedule: '每周日 03:00',
    enabled: true,
    lastRun: new Date('2024-06-09T03:00:00'),
    nextRun: new Date('2024-06-16T03:00:00'),
    createdAt: new Date('2024-06-01T10:00:00'),
    description: '清理临时文件和日志文件',
    config: {
      target: '/tmp',
      fileTypes: ['*.tmp', '*.log'],
      keepDays: 7
    },
    history: [
      {
        id: 'run1',
        status: 'success',
        startTime: new Date('2024-06-09T03:00:00'),
        duration: '5分钟'
      }
    ]
  },
  {
    id: '3',
    name: '存储快照',
    type: 'snapshot',
    status: 'running',
    schedule: '每小时',
    enabled: true,
    lastRun: new Date('2024-06-12T09:00:00'),
    nextRun: new Date('2024-06-12T10:00:00'),
    createdAt: new Date('2024-06-10T14:00:00'),
    description: '创建存储卷快照',
    history: [
      {
        id: 'run1',
        status: 'running',
        startTime: new Date('2024-06-12T09:00:00'),
        duration: '运行中...'
      }
    ]
  },
  {
    id: '4',
    name: '系统更新',
    type: 'system',
    status: 'disabled',
    schedule: '手动',
    enabled: false,
    createdAt: new Date('2024-06-05T10:00:00'),
    description: '检查并安装系统更新',
    history: []
  }
])

// 计算属性
const filteredTasks = computed(() => {
  if (selectedFilter.value === 'all') return tasks.value
  return tasks.value.filter(task => {
    if (selectedFilter.value === 'active') return task.enabled
    if (selectedFilter.value === 'scheduled') return task.status === 'active'
    if (selectedFilter.value === 'running') return task.status === 'running'
    if (selectedFilter.value === 'failed') return task.history.some(h => h.status === 'failed')
    return true
  })
})

const selectedTaskData = computed(() => {
  return tasks.value.find(t => t.id === selectedTask.value) || null
})

const activeTasks = computed(() => tasks.value.filter(t => t.enabled).length)
const completedToday = computed(() => {
  const today = new Date().toDateString()
  return tasks.value.filter(t =>
    t.history.some(h =>
      h.status === 'success' && h.startTime.toDateString() === today
    )
  ).length
})
const pendingTasks = computed(() => tasks.value.filter(t => t.enabled && t.status !== 'running').length)
const failedTasks = computed(() => tasks.value.filter(t => t.history.some(h => h.status === 'failed')).length)

// 方法
const selectTask = (task: Task) => {
  selectedTask.value = task.id
}

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    active: '活动',
    running: '运行中',
    success: '已完成',
    failed: '失败',
    disabled: '已禁用'
  }
  return statusMap[status] || status
}

const getTaskTypeText = (type: string): string => {
  const typeMap: Record<string, string> = {
    backup: '备份任务',
    cleanup: '清理任务',
    system: '系统任务',
    snapshot: '快照任务',
    custom: '自定义脚本'
  }
  return typeMap[type] || type
}

const formatTime = (date?: Date): string => {
  if (!date) return '-'
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const toggleTask = (task: Task) => {
  task.enabled = !task.enabled
}

const runTask = (task: Task) => {
  console.log('运行任务:', task.name)
  // 实现任务运行逻辑
}

const editTask = (task: Task) => {
  editingTask.value = task
  taskForm.value = {
    name: task.name,
    type: task.type,
    description: task.description,
    scheduleType: 'daily',
    scheduleTime: '02:00',
    config: task.config || {}
  }
  showTaskModal.value = true
}

const deleteTask = (task: Task) => {
  if (confirm(`确定要删除任务 "${task.name}" 吗？`)) {
    const index = tasks.value.findIndex(t => t.id === task.id)
    if (index > -1) {
      tasks.value.splice(index, 1)
      if (selectedTask.value === task.id) {
        selectedTask.value = null
      }
    }
  }
}

const viewTaskLogs = (task: Task) => {
  console.log('查看任务日志:', task.name)
}

const createNewTask = () => {
  editingTask.value = null
  taskForm.value = {
    name: '',
    type: 'backup',
    description: '',
    scheduleType: 'daily',
    scheduleTime: '02:00',
    config: {
      source: '',
      destination: '',
      compress: false,
      retention: 7
    }
  }
  showTaskModal.value = true
}

const runAllTasks = () => {
  console.log('运行所有任务')
}

const saveTask = () => {
  if (editingTask.value) {
    // 更新现有任务
    editingTask.value.name = taskForm.value.name
    editingTask.value.description = taskForm.value.description
    editingTask.value.config = taskForm.value.config
  } else {
    // 创建新任务
    const newTask: Task = {
      id: Date.now().toString(),
      name: taskForm.value.name,
      type: taskForm.value.type,
      status: 'active',
      schedule: '自定义计划',
      enabled: true,
      createdAt: new Date(),
      description: taskForm.value.description,
      config: taskForm.value.config,
      history: []
    }
    tasks.value.push(newTask)
  }
  showTaskModal.value = false
}
</script>

<style scoped>
.task-scheduler {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f9fafb;
}

.ts-header {
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

.action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.danger {
  color: #ef4444;
  border-color: #fecaca;
}

.action-btn.danger:hover {
  background: #fef2f2;
}

.ts-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.stat-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  color: #6b7280;
}

.stat-icon.active {
  background: #ecfdf5;
  color: #10b981;
}

.stat-icon.success {
  background: #eff6ff;
  color: #3b82f6;
}

.stat-icon.warning {
  background: #fef3c7;
  color: #d97706;
}

.stat-icon.error {
  background: #fef2f2;
  color: #ef4444;
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

.ts-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  flex: 1;
  overflow: hidden;
  padding: 20px;
}

.ts-tasks {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.filter-select {
  padding: 6px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.task-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.task-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.task-item:hover {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.task-item.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.task-item.active {
  border-left: 4px solid #10b981;
}

.task-item.running {
  border-left: 4px solid #3b82f6;
}

.task-item.failed {
  border-left: 4px solid #ef4444;
}

.task-item.disabled {
  opacity: 0.6;
}

.task-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f3f4f6;
  color: #6b7280;
}

.task-info {
  flex: 1;
}

.task-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.task-schedule {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 2px;
}

.task-status {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 4px;
  display: inline-block;
}

.task-status.active {
  background: #ecfdf5;
  color: #10b981;
}

.task-status.running {
  background: #eff6ff;
  color: #3b82f6;
}

.task-status.success {
  background: #f0fdf4;
  color: #059669;
}

.task-status.failed {
  background: #fef2f2;
  color: #ef4444;
}

.task-status.disabled {
  background: #f3f4f6;
  color: #9ca3af;
}

.task-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}

.task-item:hover .task-actions {
  opacity: 1;
}

.task-action-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
}

.task-action-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.task-action-btn.danger:hover {
  background: #fef2f2;
  color: #ef4444;
}

.toggle-icon {
  width: 32px;
  height: 18px;
  background: #d1d5db;
  border-radius: 9px;
  position: relative;
  cursor: pointer;
  transition: background 0.2s;
}

.toggle-icon.active {
  background: #3b82f6;
}

.toggle-dot {
  width: 14px;
  height: 14px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.2s;
}

.toggle-icon.active .toggle-dot {
  transform: translateX(14px);
}

.ts-details {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.task-details {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.details-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.details-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.task-badges {
  display: flex;
  gap: 8px;
}

.badge {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  background: #f3f4f6;
  color: #6b7280;
}

.badge.active {
  background: #ecfdf5;
  color: #10b981;
}

.badge.running {
  background: #eff6ff;
  color: #3b82f6;
}

.badge.failed {
  background: #fef2f2;
  color: #ef4444;
}

.badge.type {
  background: #eff6ff;
  color: #3b82f6;
}

.details-section {
  margin-bottom: 20px;
}

.details-section h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.details-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.detail-label {
  color: #6b7280;
}

.detail-value {
  color: #1f2937;
  font-weight: 500;
}

.detail-value.enabled {
  color: #10b981;
}

.task-description {
  font-size: 14px;
  color: #4b5563;
  line-height: 1.6;
  margin: 0;
}

.config-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.config-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.config-label {
  color: #6b7280;
}

.config-value {
  color: #1f2937;
  font-weight: 500;
}

.run-history {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px;
  background: #f9fafb;
  border-radius: 8px;
}

.history-icon {
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.history-item.success .history-icon {
  background: #ecfdf5;
  color: #10b981;
}

.history-item.failed .history-icon {
  background: #fef2f2;
  color: #ef4444;
}

.history-item.running .history-icon {
  background: #eff6ff;
  color: #3b82f6;
}

.history-info {
  flex: 1;
}

.history-time {
  font-size: 13px;
  font-weight: 500;
  color: #1f2937;
}

.history-duration {
  font-size: 11px;
  color: #6b7280;
}

.history-status {
  font-size: 12px;
  color: #6b7280;
}

.details-actions {
  display: flex;
  gap: 8px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.no-selection {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.no-selection svg {
  margin-bottom: 16px;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  background: white;
  border-radius: 16px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.modal-header button {
  padding: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 6px;
}

.modal-header button:hover {
  background: #f3f4f6;
}

.modal-body {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
}

.form-section {
  margin-bottom: 24px;
}

.form-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 6px;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  box-sizing: border-box;
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>