<template>
  <div class="sync-manager">
    <div class="sync-header">
      <h1>同步与备份</h1>
      <p class="subtitle">Restic 备份管理和多存储同步</p>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.label }}
      </button>
    </div>

    <!-- Repositories Tab -->
    <div v-if="activeTab === 'repos'" class="tab-content">
      <div class="section-header">
        <h2>备份仓库</h2>
        <button class="action-btn primary" @click="showAddRepoModal = true">
          <PlusIcon class="w-4 h-4" />
          添加仓库
        </button>
      </div>

      <div class="repos-list">
        <div
          v-for="repo in repositories"
          :key="repo.id"
          class="repo-item"
          :class="{ 'repo-active': repo.status === 'active' }"
        >
          <div class="repo-icon">
            <component :is="getRepoIcon(repo.type)" class="w-8 h-8" />
          </div>
          <div class="repo-details">
            <h4>{{ repo.name }}</h4>
            <p>{{ repo.url }}</p>
            <div class="repo-meta">
              <span class="repo-type">{{ getRepoTypeName(repo.type) }}</span>
              <span class="repo-status" :class="repo.status">{{ getStatusName(repo.status) }}</span>
              <span>最后备份: {{ formatDateTime(repo.lastBackup) }}</span>
            </div>
          </div>

          <div class="repo-stats">
            <div class="stat-item">
              <span class="stat-label">快照数</span>
              <span class="stat-value">{{ repo.snapshotCount }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">总大小</span>
              <span class="stat-value">{{ formatBytes(repo.totalSize) }}</span>
            </div>
          </div>

          <div class="repo-actions">
            <button class="action-btn" @click="backupRepo(repo)">
              <CloudArrowUpIcon class="w-4 h-4" />
              备份
            </button>
            <button class="action-btn" @click="checkRepo(repo)">
              <CheckIcon class="w-4 h-4" />
              检查
            </button>
            <button class="action-btn" @click="forgetRepo(repo)">
              <TrashIcon class="w-4 h-4" />
              清理
            </button>
            <button class="action-btn danger" @click="deleteRepo(repo)">
              <XMarkIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>
      </div>

      <div v-if="repositories.length === 0" class="empty-state">
        <CircleStackIcon class="w-16 h-16" />
        <h3>暂无备份仓库</h3>
        <p>添加您的第一个备份仓库开始保护数据</p>
        <button class="action-btn primary" @click="showAddRepoModal = true">
          <PlusIcon class="w-4 h-4" />
          添加仓库
        </button>
      </div>
    </div>

    <!-- Backup Jobs Tab -->
    <div v-if="activeTab === 'jobs'" class="tab-content">
      <div class="section-header">
        <h2>备份任务</h2>
        <button class="action-btn primary" @click="showAddJobModal = true">
          <PlusIcon class="w-4 h-4" />
          创建任务
        </button>
      </div>

      <div class="jobs-list">
        <div
          v-for="job in backupJobs"
          :key="job.id"
          class="job-item"
        >
          <div class="job-info">
            <div class="job-status" :class="job.status.toLowerCase()"></div>
            <div class="job-details">
              <h4>{{ job.name }}</h4>
              <p>{{ job.description }}</p>
              <div class="job-meta">
                <span>源路径: {{ job.sourcePath }}</span>
                <span>目标: {{ getRepoName(job.repoId) }}</span>
                <span v-if="job.schedule">计划: {{ formatSchedule(job.schedule) }}</span>
              </div>
            </div>
          </div>

          <div class="job-stats">
            <div class="stat-item">
              <span class="stat-label">最后运行</span>
              <span class="stat-value">{{ formatDateTime(job.lastRun) }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">下次运行</span>
              <span class="stat-value">{{ formatDateTime(job.nextRun) }}</span>
            </div>
          </div>

          <div class="job-actions">
            <button
              v-if="job.status !== 'running'"
              class="action-btn primary"
              @click="runJob(job)"
            >
              <PlayIcon class="w-4 h-4" />
              运行
            </button>
            <button
              v-else
              class="action-btn warning"
              @click="stopJob(job)"
            >
              <StopIcon class="w-4 h-4" />
              停止
            </button>
            <button class="action-btn" @click="editJob(job)">
              <PencilIcon class="w-4 h-4" />
              编辑
            </button>
            <button class="action-btn danger" @click="deleteJob(job)">
              <TrashIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>
      </div>

      <div v-if="backupJobs.length === 0" class="empty-state">
        <ClockIcon class="w-16 h-16" />
        <h3>暂无备份任务</h3>
        <p>创建您的第一个自动化备份任务</p>
        <button class="action-btn primary" @click="showAddJobModal = true">
          <PlusIcon class="w-4 h-4" />
          创建任务
        </button>
      </div>
    </div>

    <!-- Snapshots Tab -->
    <div v-if="activeTab === 'snapshots'" class="tab-content">
      <div class="section-header">
        <h2>快照管理</h2>
        <div class="filters">
          <select v-model="selectedRepoId" class="filter-select">
            <option value="">所有仓库</option>
            <option v-for="repo in repositories" :key="repo.id" :value="repo.id">
              {{ repo.name }}
            </option>
          </select>
          <button class="action-btn" @click="loadSnapshots">
            <ArrowPathIcon class="w-4 h-4" />
            刷新
          </button>
        </div>
      </div>

      <div class="snapshots-list">
        <div
          v-for="snapshot in filteredSnapshots"
          :key="snapshot.id"
          class="snapshot-item"
        >
          <div class="snapshot-icon">
            <CameraIcon class="w-8 h-8" />
          </div>
          <div class="snapshot-details">
            <h4>{{ snapshot.description }}</h4>
            <div class="snapshot-meta">
              <span>时间: {{ formatDateTime(snapshot.time) }}</span>
              <span>仓库: {{ snapshot.repoName }}</span>
              <span>大小: {{ formatBytes(snapshot.size) }}</span>
              <span>文件数: {{ snapshot.fileCount }}</span>
            </div>
          </div>

          <div class="snapshot-actions">
            <button class="action-btn" @click="browseSnapshot(snapshot)">
              <FolderIcon class="w-4 h-4" />
              浏览
            </button>
            <button class="action-btn" @click="restoreSnapshot(snapshot)">
              <ArrowUturnUpIcon class="w-4 h-4" />
              恢复
            </button>
            <button class="action-btn danger" @click="deleteSnapshot(snapshot)">
              <TrashIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>
      </div>

      <div v-if="filteredSnapshots.length === 0" class="empty-state">
        <CameraIcon class="w-16 h-16" />
        <h3>暂无快照</h3>
        <p>运行备份后快照将显示在这里</p>
      </div>
    </div>

    <!-- Settings Tab -->
    <div v-if="activeTab === 'settings'" class="tab-content">
      <div class="section-header">
        <h2>Restic 设置</h2>
      </div>

      <div class="settings-grid">
        <!-- Global Settings -->
        <div class="setting-card">
          <div class="card-header">
            <h3>全局设置</h3>
          </div>
          <form @submit.prevent="saveGlobalSettings" class="settings-form">
            <div class="form-group">
              <label>
                <input v-model="globalSettings.keepVerifications" type="checkbox" />
                保留验证数据
              </label>
            </div>
            <div class="form-group">
              <label>
                <input v-model="globalSettings.useCache" type="checkbox" />
                启用缓存
              </label>
            </div>
            <div class="form-group">
              <label>缓存目录</label>
              <input
                v-model="globalSettings.cacheDir"
                type="text"
                placeholder="/tmp/restic-cache"
              />
            </div>
            <div class="form-group">
              <label>压缩级别</label>
              <select v-model="globalSettings.compression">
                <option value="auto">自动</option>
                <option value="max">最大</option>
                <option value="off">关闭</option>
              </select>
            </div>
            <button type="submit" class="action-btn primary">
              保存设置
            </button>
          </form>
        </div>

        <!-- Retention Policy -->
        <div class="setting-card">
          <div class="card-header">
            <h3>保留策略</h3>
          </div>
          <form @submit.prevent="saveRetentionPolicy" class="settings-form">
            <div class="form-group">
              <label>保留最近备份 (天)</label>
              <input
                v-model="retention.keepLast"
                type="number"
                min="0"
                placeholder="7"
              />
            </div>
            <div class="form-group">
              <label>保留小时备份 (个)</label>
              <input
                v-model="retention.keepHourly"
                type="number"
                min="0"
                placeholder="24"
              />
            </div>
            <div class="form-group">
              <label>保留每日备份 (个)</label>
              <input
                v-model="retention.keepDaily"
                type="number"
                min="0"
                placeholder="7"
              />
            </div>
            <div class="form-group">
              <label>保留每周备份 (个)</label>
              <input
                v-model="retention.keepWeekly"
                type="number"
                min="0"
                placeholder="4"
              />
            </div>
            <div class="form-group">
              <label>保留每月备份 (个)</label>
              <input
                v-model="retention.keepMonthly"
                type="number"
                min="0"
                placeholder="12"
              />
            </div>
            <div class="form-group">
              <label>保留每年备份 (个)</label>
              <input
                v-model="retention.keepYearly"
                type="number"
                min="0"
                placeholder="10"
              />
            </div>
            <button type="submit" class="action-btn primary">
              保存策略
            </button>
          </form>
        </div>

        <!-- Notifications -->
        <div class="setting-card">
          <div class="card-header">
            <h3>通知设置</h3>
          </div>
          <form @submit.prevent="saveNotificationSettings" class="settings-form">
            <div class="form-group">
              <label>
                <input v-model="notifications.enabled" type="checkbox" />
                启用通知
              </label>
            </div>
            <div v-if="notifications.enabled" class="notification-fields">
              <div class="form-group">
                <label>通知类型</label>
                <div class="checkbox-group">
                  <label v-for="type in notificationTypes" :key="type.value" class="checkbox-label">
                    <input
                      v-model="notifications.types"
                      type="checkbox"
                      :value="type.value"
                    />
                    {{ type.label }}
                  </label>
                </div>
              </div>
              <div class="form-group">
                <label>Email 地址</label>
                <input
                  v-model="notifications.email"
                  type="email"
                  placeholder="admin@example.com"
                />
              </div>
              <div class="form-group">
                <label>Webhook URL</label>
                <input
                  v-model="notifications.webhook"
                  type="url"
                  placeholder="https://hooks.example.com/backup"
                />
              </div>
            </div>
            <button type="submit" class="action-btn primary">
              保存设置
            </button>
          </form>
        </div>
      </div>
    </div>

    <!-- Add Repository Modal -->
    <div v-if="showAddRepoModal" class="modal-overlay" @click.self="showAddRepoModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>添加备份仓库</h3>
          <button class="close-btn" @click="showAddRepoModal = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit.prevent="addRepository" class="modal-form">
          <div class="form-group">
            <label>仓库名称</label>
            <input
              v-model="newRepo.name"
              type="text"
              required
              placeholder="我的备份仓库"
            />
          </div>
          <div class="form-group">
            <label>存储类型</label>
            <select v-model="newRepo.type" required>
              <option value="local">本地文件系统</option>
              <option value="s3">S3 兼容存储</option>
              <option value="webdav">WebDAV</option>
              <option value="sftp">SFTP</option>
              <option value="rest">REST 服务器</option>
              <option value="azure">Azure Blob Storage</option>
              <option value="gcs">Google Cloud Storage</option>
              <option value="b2">Backblaze B2</option>
            </select>
          </div>
          <div class="form-group">
            <label>仓库路径/URL</label>
            <input
              v-model="newRepo.url"
              type="text"
              required
              :placeholder="getPlaceholderForType(newRepo.type)"
            />
          </div>
          <div v-if="requiresAuth(newRepo.type)" class="auth-fields">
            <div class="form-group">
              <label>用户名/Access Key</label>
              <input
                v-model="newRepo.username"
                type="text"
                placeholder="用户名或Access Key"
              />
            </div>
            <div class="form-group">
              <label>密码/Secret Key</label>
              <input
                v-model="newRepo.password"
                type="password"
                placeholder="密码或Secret Key"
              />
            </div>
          </div>
          <div v-if="newRepo.type === 's3'" class="s3-fields">
            <div class="form-group">
              <label>区域</label>
              <input
                v-model="newRepo.region"
                type="text"
                placeholder="us-east-1"
              />
            </div>
            <div class="form-group">
              <label>终端点</label>
              <input
                v-model="newRepo.endpoint"
                type="text"
                placeholder="s3.amazonaws.com"
              />
            </div>
          </div>
          <div class="form-group">
            <label>密码 (用于仓库加密)</label>
            <input
              v-model="newRepo.password"
              type="password"
              required
              placeholder="选择一个安全的密码"
            />
          </div>
          <div class="form-group">
            <label>描述 (可选)</label>
            <input
              v-model="newRepo.description"
              type="text"
              placeholder="仓库用途描述"
            />
          </div>
          <div class="modal-actions">
            <button type="button" class="action-btn" @click="showAddRepoModal = false">
              取消
            </button>
            <button type="submit" class="action-btn primary">
              添加仓库
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Add Job Modal -->
    <div v-if="showAddJobModal" class="modal-overlay" @click.self="showAddJobModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>创建备份任务</h3>
          <button class="close-btn" @click="showAddJobModal = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit.prevent="addBackupJob" class="modal-form">
          <div class="form-group">
            <label>任务名称</label>
            <input
              v-model="newJob.name"
              type="text"
              required
              placeholder="每日备份"
            />
          </div>
          <div class="form-group">
            <label>源路径</label>
            <input
              v-model="newJob.sourcePath"
              type="text"
              required
              placeholder="/data/important"
            />
          </div>
          <div class="form-group">
            <label>目标仓库</label>
            <select v-model="newJob.repoId" required>
              <option value="">选择仓库</option>
              <option v-for="repo in repositories" :key="repo.id" :value="repo.id">
                {{ repo.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>
              <input v-model="newJob.enabled" type="checkbox" />
              启用自动运行
            </label>
          </div>
          <div v-if="newJob.enabled" class="schedule-fields">
            <div class="form-group">
              <label>计划类型</label>
              <select v-model="newJob.scheduleType">
                <option value="interval">间隔</option>
                <option value="cron">Cron 表达式</option>
              </select>
            </div>
            <div v-if="newJob.scheduleType === 'interval'" class="form-group">
              <label>运行间隔</label>
              <div class="interval-input">
                <input
                  v-model="newJob.intervalValue"
                  type="number"
                  min="1"
                  placeholder="1"
                />
                <select v-model="newJob.intervalUnit">
                  <option value="minutes">分钟</option>
                  <option value="hours">小时</option>
                  <option value="days">天</option>
                  <option value="weeks">周</option>
                </select>
              </div>
            </div>
            <div v-if="newJob.scheduleType === 'cron'" class="form-group">
              <label>Cron 表达式</label>
              <input
                v-model="newJob.cronExpression"
                type="text"
                placeholder="0 2 * * *"
              />
              <small>格式: 分 时 日 月 周</small>
            </div>
          </div>
          <div class="form-group">
            <label>排除模式 (可选)</label>
            <input
              v-model="newJob.excludes"
              type="text"
              placeholder="*.tmp,*.log,node_modules"
            />
            <small>逗号分隔的文件模式</small>
          </div>
          <div class="form-group">
            <label>描述 (可选)</label>
            <input
              v-model="newJob.description"
              type="text"
              placeholder="任务描述"
            />
          </div>
          <div class="modal-actions">
            <button type="button" class="action-btn" @click="showAddJobModal = false">
              取消
            </button>
            <button type="submit" class="action-btn primary">
              创建任务
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  CircleStackIcon,
  ClockIcon,
  CloudArrowUpIcon,
  CheckIcon,
  TrashIcon,
  PlusIcon,
  PencilIcon,
  PlayIcon,
  StopIcon,
  XMarkIcon,
  ArrowPathIcon,
  CameraIcon,
  FolderIcon,
  ArrowUturnUpIcon,
  GlobeAltIcon,
  ServerIcon,
  CloudIcon
} from '@heroicons/vue/24/outline'

const activeTab = ref('repos')

const tabs = [
  { id: 'repos', label: '仓库', icon: CircleStackIcon },
  { id: 'jobs', label: '任务', icon: ClockIcon },
  { id: 'snapshots', label: '快照', icon: CameraIcon },
  { id: 'settings', label: '设置', icon: ServerIcon }
]

// Repositories
const repositories = ref<any[]>([])
const showAddRepoModal = ref(false)
const newRepo = ref({
  name: '',
  type: 'local',
  url: '',
  username: '',
  password: '',
  region: '',
  endpoint: '',
  description: ''
})

// Backup Jobs
const backupJobs = ref<any[]>([])
const showAddJobModal = ref(false)
const newJob = ref({
  name: '',
  sourcePath: '',
  repoId: '',
  enabled: true,
  scheduleType: 'interval',
  intervalValue: 24,
  intervalUnit: 'hours',
  cronExpression: '0 2 * * *',
  excludes: '',
  description: ''
})

// Snapshots
const snapshots = ref<any[]>([])
const selectedRepoId = ref('')

// Settings
const globalSettings = ref({
  keepVerifications: true,
  useCache: true,
  cacheDir: '/tmp/restic-cache',
  compression: 'auto'
})

const retention = ref({
  keepLast: 7,
  keepHourly: 24,
  keepDaily: 7,
  keepWeekly: 4,
  keepMonthly: 12,
  keepYearly: 10
})

const notifications = ref({
  enabled: true,
  types: ['success', 'error'],
  email: '',
  webhook: ''
})

const notificationTypes = [
  { value: 'success', label: '成功通知' },
  { value: 'error', label: '错误通知' },
  { value: 'warning', label: '警告通知' }
]

// Computed
const filteredSnapshots = computed(() => {
  if (!selectedRepoId.value) return snapshots.value
  return snapshots.value.filter(s => s.repoId === selectedRepoId.value)
})

// Functions
const getRepoIcon = (type: string) => {
  const icons: Record<string, any> = {
    local: ServerIcon,
    s3: CloudIcon,
    webdav: GlobeAltIcon,
    sftp: ServerIcon,
    rest: GlobeAltIcon,
    azure: CloudIcon,
    gcs: CloudIcon,
    b2: CloudIcon
  }
  return icons[type] || ServerIcon
}

const getRepoTypeName = (type: string) => {
  const types: Record<string, string> = {
    local: '本地',
    s3: 'S3',
    webdav: 'WebDAV',
    sftp: 'SFTP',
    rest: 'REST',
    azure: 'Azure',
    gcs: 'GCS',
    b2: 'B2'
  }
  return types[type] || type
}

const getStatusName = (status: string) => {
  const statuses: Record<string, string> = {
    active: '活动',
    inactive: '非活动',
    error: '错误',
    checking: '检查中'
  }
  return statuses[status] || status
}

const getPlaceholderForType = (type: string) => {
  const placeholders: Record<string, string> = {
    local: '/mnt/backups/myrepo',
    s3: 's3://bucket-name/path',
    webdav: 'https://example.com/webdav/path',
    sftp: 'sftp://server.com/backup/path',
    rest: 'rest:https://example.com/api',
    azure: 'azure:container-name/path',
    gcs: 'gs:bucket-name/path',
    b2: 'b2:bucket-name/path'
  }
  return placeholders[type] || ''
}

const requiresAuth = (type: string) => {
  return ['s3', 'webdav', 'sftp', 'rest', 'azure', 'gcs', 'b2'].includes(type)
}

const getRepoName = (repoId: string) => {
  const repo = repositories.value.find(r => r.id === repoId)
  return repo ? repo.name : 'Unknown'
}

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const formatDateTime = (dateStr: string) => {
  if (!dateStr) return 'Never'
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatSchedule = (schedule: any) => {
  if (!schedule) return 'Manual'
  if (schedule.type === 'interval') {
    return `每 ${schedule.value} ${schedule.unit}`
  }
  return schedule.cron || 'Custom'
}

// Repository operations
const addRepository = async () => {
  try {
    console.log('Adding repository:', newRepo.value)
    // Mock add
    repositories.value.push({
      id: Date.now().toString(),
      ...newRepo.value,
      status: 'active',
      snapshotCount: 0,
      totalSize: 0,
      lastBackup: null
    })
    showAddRepoModal.value = false
    // Reset form
    newRepo.value = {
      name: '',
      type: 'local',
      url: '',
      username: '',
      password: '',
      region: '',
      endpoint: '',
      description: ''
    }
  } catch (error: any) {
    console.error('Failed to add repository:', error)
    alert('添加仓库失败: ' + error.message)
  }
}

const backupRepo = async (repo: any) => {
  try {
    console.log('Backing up repository:', repo.name)
    alert('开始备份: ' + repo.name)
  } catch (error: any) {
    console.error('Backup failed:', error)
    alert('备份失败: ' + error.message)
  }
}

const checkRepo = async (repo: any) => {
  try {
    console.log('Checking repository:', repo.name)
    alert('检查仓库: ' + repo.name)
  } catch (error: any) {
    console.error('Check failed:', error)
    alert('检查失败: ' + error.message)
  }
}

const forgetRepo = async (repo: any) => {
  if (confirm(`确定要清理仓库 "${repo.name}" 的旧快照吗?`)) {
    try {
      console.log('Forgetting old snapshots:', repo.name)
      alert('清理完成: ' + repo.name)
    } catch (error: any) {
      console.error('Forget failed:', error)
      alert('清理失败: ' + error.message)
    }
  }
}

const deleteRepo = async (repo: any) => {
  if (confirm(`确定要删除仓库 "${repo.name}" 吗? 这将删除所有快照数据。`)) {
    try {
      repositories.value = repositories.value.filter(r => r.id !== repo.id)
      alert('仓库已删除')
    } catch (error: any) {
      console.error('Delete failed:', error)
      alert('删除失败: ' + error.message)
    }
  }
}

// Job operations
const addBackupJob = async () => {
  try {
    console.log('Adding backup job:', newJob.value)
    // Mock add
    backupJobs.value.push({
      id: Date.now().toString(),
      ...newJob.value,
      status: 'idle',
      lastRun: null,
      nextRun: newJob.value.enabled ? new Date().toISOString() : null,
      schedule: newJob.value.enabled ? {
        type: newJob.value.scheduleType,
        value: newJob.value.intervalValue,
        unit: newJob.value.intervalUnit,
        cron: newJob.value.cronExpression
      } : null
    })
    showAddJobModal.value = false
    // Reset form
    newJob.value = {
      name: '',
      sourcePath: '',
      repoId: '',
      enabled: true,
      scheduleType: 'interval',
      intervalValue: 24,
      intervalUnit: 'hours',
      cronExpression: '0 2 * * *',
      excludes: '',
      description: ''
    }
  } catch (error: any) {
    console.error('Failed to add job:', error)
    alert('创建任务失败: ' + error.message)
  }
}

const runJob = async (job: any) => {
  try {
    console.log('Running job:', job.name)
    alert('开始运行任务: ' + job.name)
  } catch (error: any) {
    console.error('Run job failed:', error)
    alert('运行失败: ' + error.message)
  }
}

const stopJob = async (job: any) => {
  try {
    console.log('Stopping job:', job.name)
    alert('停止任务: ' + job.name)
  } catch (error: any) {
    console.error('Stop job failed:', error)
    alert('停止失败: ' + error.message)
  }
}

const editJob = (job: any) => {
  console.log('Editing job:', job.name)
  alert('编辑任务: ' + job.name)
}

const deleteJob = (job: any) => {
  if (confirm(`确定要删除任务 "${job.name}" 吗?`)) {
    backupJobs.value = backupJobs.value.filter(j => j.id !== job.id)
    alert('任务已删除')
  }
}

// Snapshot operations
const browseSnapshot = (snapshot: any) => {
  console.log('Browsing snapshot:', snapshot.id)
  alert('浏览快照: ' + snapshot.description)
}

const restoreSnapshot = (snapshot: any) => {
  if (confirm(`确定要恢复快照 "${snapshot.description}" 吗?`)) {
    console.log('Restoring snapshot:', snapshot.id)
    alert('恢复快照: ' + snapshot.description)
  }
}

const deleteSnapshot = (snapshot: any) => {
  if (confirm(`确定要删除快照 "${snapshot.description}" 吗?`)) {
    snapshots.value = snapshots.value.filter(s => s.id !== snapshot.id)
    alert('快照已删除')
  }
}

const loadSnapshots = () => {
  console.log('Loading snapshots...')
}

// Settings operations
const saveGlobalSettings = () => {
  console.log('Saving global settings:', globalSettings.value)
  alert('全局设置已保存')
}

const saveRetentionPolicy = () => {
  console.log('Saving retention policy:', retention.value)
  alert('保留策略已保存')
}

const saveNotificationSettings = () => {
  console.log('Saving notification settings:', notifications.value)
  alert('通知设置已保存')
}

onMounted(() => {
  // Mock data
  repositories.value = [
    {
      id: '1',
      name: '本地备份',
      type: 'local',
      url: '/mnt/backup/restic',
      status: 'active',
      snapshotCount: 15,
      totalSize: 1024000000,
      lastBackup: new Date().toISOString(),
      description: '本地存储备份'
    },
    {
      id: '2',
      name: 'S3 云备份',
      type: 's3',
      url: 's3://my-backup-bucket/nas',
      status: 'active',
      snapshotCount: 10,
      totalSize: 512000000,
      lastBackup: new Date(Date.now() - 86400000).toISOString(),
      description: 'AWS S3 备份'
    }
  ]

  backupJobs.value = [
    {
      id: '1',
      name: '每日数据备份',
      description: '备份重要数据文件',
      sourcePath: '/data/important',
      repoId: '1',
      status: 'idle',
      enabled: true,
      lastRun: new Date().toISOString(),
      nextRun: new Date(Date.now() + 86400000).toISOString(),
      schedule: {
        type: 'interval',
        value: 24,
        unit: 'hours'
      }
    },
    {
      id: '2',
      name: '配置文件备份',
      description: '备份系统配置',
      sourcePath: '/etc',
      repoId: '1',
      status: 'idle',
      enabled: true,
      lastRun: new Date(Date.now() - 86400000).toISOString(),
      nextRun: new Date(Date.now() + 86400000).toISOString(),
      schedule: {
        type: 'cron',
        cron: '0 3 * * *'
      }
    }
  ]

  snapshots.value = [
    {
      id: '1',
      repoId: '1',
      repoName: '本地备份',
      description: '自动备份-2024-01-15',
      time: new Date().toISOString(),
      size: 102400000,
      fileCount: 1524
    },
    {
      id: '2',
      repoId: '2',
      repoName: 'S3 云备份',
      description: '自动备份-2024-01-14',
      time: new Date(Date.now() - 86400000).toISOString(),
      size: 51200000,
      fileCount: 856
    }
  ]
})
</script>

<style scoped>
.sync-manager {
  width: 100%;
  height: 100%;
  padding: 32px;
  background: #f9fafb;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.sync-header {
  margin-bottom: 32px;
}

.sync-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: #6b7280;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 2px solid #e5e7eb;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: -2px;
}

.tab-btn:hover {
  color: #1f2937;
}

.tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

.tab-content {
  flex: 1;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

.filters {
  display: flex;
  gap: 12px;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.warning {
  background: #f59e0b;
  border-color: #f59e0b;
  color: white;
}

.action-btn.warning:hover {
  background: #d97706;
}

.action-btn.danger {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.action-btn.danger:hover {
  background: #dc2626;
}

.repos-list,
.jobs-list,
.snapshots-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.repo-item,
.job-item,
.snapshot-item {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 24px;
  align-items: center;
}

.repo-active {
  border: 2px solid #3b82f6;
}

.repo-icon,
.snapshot-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.snapshot-icon {
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
}

.repo-details,
.snapshot-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.repo-details h4,
.snapshot-details h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.repo-details p {
  font-size: 14px;
  color: #6b7280;
}

.repo-meta,
.snapshot-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #9ca3af;
}

.repo-type {
  padding: 4px 8px;
  background: #f3f4f6;
  border-radius: 4px;
  font-weight: 500;
}

.repo-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-weight: 500;
}

.repo-status.active {
  color: #10b981;
  background: #d1fae5;
}

.repo-status.inactive {
  color: #6b7280;
  background: #f3f4f6;
}

.repo-status.error {
  color: #ef4444;
  background: #fee2e2;
}

.repo-stats {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 120px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.job-item {
  grid-template-columns: 1fr auto auto;
}

.job-info {
  display: flex;
  gap: 12px;
}

.job-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 6px;
}

.job-status.idle {
  background: #6b7280;
}

.job-status.running {
  background: #10b981;
}

.job-status.error {
  background: #ef4444;
}

.job-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.job-details h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.job-details p {
  font-size: 14px;
  color: #6b7280;
}

.job-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
}

.job-stats {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 150px;
}

.snapshot-item {
  grid-template-columns: auto 1fr auto;
}

.repo-actions,
.snapshot-actions,
.job-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  background: white;
  border-radius: 12px;
  padding: 48px;
  text-align: center;
  color: #6b7280;
}

.empty-state svg {
  margin-bottom: 16px;
  color: #9ca3af;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.empty-state p {
  font-size: 14px;
  margin-bottom: 24px;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.setting-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input[type="text"],
.form-group input[type="password"],
.form-group input[type="email"],
.form-group input[type="url"],
.form-group input[type="number"],
.form-group input[type="time"],
.form-group select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.form-group input[type="checkbox"] {
  width: 16px;
  height: 16px;
  margin-right: 8px;
}

.form-group small {
  font-size: 12px;
  color: #6b7280;
}

.auth-fields,
.s3-fields,
.schedule-fields,
.notification-fields {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 12px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.interval-input {
  display: flex;
  gap: 8px;
}

.interval-input input {
  flex: 1;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  padding: 32px;
  max-width: 600px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.modal-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: #f3f4f6;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: #e5e7eb;
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 24px;
}
</style>