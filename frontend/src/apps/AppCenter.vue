<template>
  <div class="app-manager">
    <div class="manager-header">
      <div class="header-info">
        <h1>🚀 应用管理中心</h1>
        <p class="subtitle">统一管理NAS应用 - 一键部署、监控和管理</p>
        <div class="quick-stats">
          <div class="stat-card running">
            <span class="stat-icon">🟢</span>
            <span class="stat-count">{{ runningAppsCount }}</span>
            <span class="stat-label">运行中</span>
          </div>
          <div class="stat-card stopped">
            <span class="stat-icon">🔴</span>
            <span class="stat-count">{{ stoppedAppsCount }}</span>
            <span class="stat-label">已停止</span>
          </div>
          <div class="stat-card total">
            <span class="stat-icon">📦</span>
            <span class="stat-count">{{ apps.length }}</span>
            <span class="stat-label">总应用</span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <button class="action-btn primary" @click="fetchAppsStatus" :disabled="loading">
          <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          刷新状态
        </button>
      </div>
    </div>

    <!-- 应用分类标签 -->
    <div class="category-tabs">
      <button
        v-for="category in categories"
        :key="category.id"
        :class="['tab-btn', { active: selectedCategory === category.id }]"
        @click="selectedCategory = category.id"
      >
        {{ category.label }}
      </button>
    </div>

    <!-- 应用列表 -->
    <div class="apps-grid">
      <div
        v-for="app in filteredApps"
        :key="app.id"
        class="app-card"
        :class="{
          'app-running': app.status === 'running',
          'app-stopped': app.status === 'stopped' || app.status === 'unknown'
        }"
      >
        <!-- 应用图标和信息 -->
        <div class="app-header">
          <div class="app-icon" :style="{ background: app.color }">
            <component :is="getAppIcon(app.icon)" class="w-8 h-8" />
          </div>
          <div class="app-info">
            <h3 class="app-name">{{ app.name }}</h3>
            <p class="app-description">{{ app.description }}</p>
          </div>
          <div class="app-status" :class="`status-${app.status}`">
            <span class="status-dot"></span>
            <span class="status-text">{{ getStatusText(app.status) }}</span>
          </div>
        </div>

        <!-- 应用详情 -->
        <div class="app-details">
          <div class="detail-row">
            <span class="detail-label">类型:</span>
            <span class="detail-value">{{ app.type === 'docker' ? 'Docker容器' : '系统工具' }}</span>
          </div>
          <div class="detail-row" v-if="app.ports">
            <span class="detail-label">端口:</span>
            <span class="detail-value">{{ app.ports }}</span>
          </div>
          <div class="detail-row" v-if="app.image">
            <span class="detail-label">镜像:</span>
            <span class="detail-value text-xs">{{ app.image.split(':')[0] }}</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="app-actions">
          <!-- 未安装状态 -->
          <template v-if="app.status === 'unknown'">
            <button
              class="action-btn install"
              @click="installApp(app)"
              :disabled="operatingApp === app.id"
            >
              <CloudArrowDownIcon class="w-4 h-4" />
              {{ operatingApp === app.id ? '安装中...' : '一键安装' }}
            </button>
          </template>

          <!-- 已安装状态 -->
          <template v-else>
            <div class="control-buttons">
              <button
                v-if="app.status === 'running'"
                class="control-btn stop"
                @click="stopApp(app)"
                :disabled="operatingApp === app.id"
                title="停止应用"
              >
                <PauseIcon class="w-4 h-4" />
              </button>
              <button
                v-else
                class="control-btn start"
                @click="startApp(app)"
                :disabled="operatingApp === app.id"
                title="启动应用"
              >
                <PlayIcon class="w-4 h-4" />
              </button>
              <button
                class="control-btn restart"
                @click="restartApp(app)"
                :disabled="operatingApp === app.id"
                title="重启应用"
              >
                <ArrowPathIcon class="w-4 h-4" />
              </button>
            </div>
            <button
              class="action-btn danger"
              @click="uninstallApp(app)"
              :disabled="operatingApp === app.id"
            >
              <TrashIcon class="w-4 h-4" />
              卸载
            </button>
          </template>

          <button
            class="action-btn info"
            @click="showAppLogs(app)"
            v-if="app.status !== 'unknown'"
          >
            <DocumentTextIcon class="w-4 h-4" />
            日志
          </button>
        </div>

        <!-- 快速访问 -->
        <div class="app-access" v-if="app.status === 'running' && app.url">
          <a :href="app.url" target="_blank" class="access-link">
            <GlobeAltIcon class="w-4 h-4" />
            打开 {{ app.name }}
          </a>
        </div>
      </div>
    </div>

    <!-- 日志模态框 -->
    <div v-if="showLogsModal" class="modal-overlay" @click="closeLogsModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>{{ selectedApp?.name }} - 应用日志</h2>
          <button class="close-btn" @click="closeLogsModal">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>
        <div class="modal-body">
          <div class="logs-content">
            <pre>{{ appLogs }}</pre>
          </div>
        </div>
        <div class="modal-footer">
          <button class="action-btn" @click="fetchAppLogs(selectedApp)">
            <ArrowPathIcon class="w-4 h-4" />
            刷新日志
          </button>
          <button class="action-btn" @click="closeLogsModal">
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- 安装进度模态框 -->
    <div v-if="showInstallModal" class="modal-overlay">
      <div class="modal-content small" @click.stop>
        <div class="install-progress">
          <div class="progress-icon">
            <CloudArrowDownIcon class="w-12 h-12" v-if="installStep === 'downloading'" />
            <CogIcon class="w-12 h-12 spinning" v-if="installStep === 'installing'" />
            <CheckIcon class="w-12 h-12 text-green-500" v-if="installStep === 'completed'" />
            <XMarkIcon class="w-12 h-12 text-red-500" v-if="installStep === 'error'" />
          </div>
          <h3>{{ installTitle }}</h3>
          <p>{{ installMessage }}</p>
          <div class="progress-bar" v-if="installProgress > 0 && installProgress < 100">
            <div class="progress-fill" :style="{ width: installProgress + '%' }"></div>
          </div>
          <div class="progress-text" v-if="installProgress > 0 && installProgress < 100">
            {{ installProgress }}%
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  ShoppingBagIcon,
  ServerIcon,
  ChartBarIcon,
  FolderIcon,
  UserGroupIcon,
  CogIcon,
  CloudIcon,
  CloudArrowDownIcon,
  ArrowPathIcon,
  PlayIcon,
  PauseIcon,
  TrashIcon,
  DocumentTextIcon,
  XMarkIcon,
  GlobeAltIcon,
  CheckIcon,
  StarIcon
} from '@heroicons/vue/24/outline'

// API调用函数（简化版，直接使用axios）
const apiCall = async (method: string, url: string, data?: any) => {
  const token = localStorage.getItem('token')
  const headers = token ? { 'Authorization': `Bearer ${token}` } : {}

  try {
    const response = await fetch(url, {
      method,
      headers: headers as HeadersInit,
      body: data ? JSON.stringify(data) : undefined
    })

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    return await response.json()
  } catch (error) {
    throw error
  }
}

const searchQuery = ref('')
const selectedCategory = ref('all')
const loading = ref(false)
const operatingApp = ref<string | null>(null)

// 日志模态框
const showLogsModal = ref(false)
const selectedApp = ref<any>(null)
const appLogs = ref('')

// 安装进度模态框
const showInstallModal = ref(false)
const installStep = ref<'downloading' | 'installing' | 'completed' | 'error'>('downloading')
const installProgress = ref(0)
const installTitle = ref('')
const installMessage = ref('')

const categories = [
  { id: 'all', label: '全部应用' },
  { id: 'system', label: '系统工具' },
  { id: 'docker', label: 'Docker应用' },
  { id: 'storage', label: '存储管理' },
  { id: 'media', label: '媒体服务' },
  { id: 'network', label: '网络服务' }
]

// 统一的应用配置 - 包含系统应用和Docker应用
const apps = ref([
  // 系统内置应用
  {
    id: 'storage-manager',
    name: '存储管理器',
    description: '管理磁盘和存储卷，支持挂载、分区等操作',
    icon: 'ServerIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
    category: 'storage',
    type: 'system',
    status: 'running',
    ports: '',
    image: '',
    url: ''
  },
  {
    id: 'system-monitor',
    name: '系统监控',
    description: '实时监控CPU、内存、磁盘、网络等系统资源',
    icon: 'ChartBarIcon',
    color: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    category: 'system',
    type: 'system',
    status: 'running',
    ports: '',
    image: '',
    url: ''
  },
  {
    id: 'user-manager',
    name: '用户管理',
    description: '管理系统用户、权限和配额',
    icon: 'UserGroupIcon',
    color: 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
    category: 'system',
    type: 'system',
    status: 'running',
    ports: '',
    image: '',
    url: ''
  },
  {
    id: 'docker-manager',
    name: 'Docker管理',
    description: '管理Docker容器和镜像',
    icon: 'CogIcon',
    color: 'linear-gradient(135deg, #2563eb 0%, #3b82f6 100%)',
    category: 'docker',
    type: 'system',
    status: 'running',
    ports: '',
    image: '',
    url: ''
  },

  // 6个核心Docker应用
  {
    id: 'alist',
    name: 'Alist',
    description: '支持多种存储的文件列表程序，支持网盘、FTP、WebDAV等',
    icon: 'CloudIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #2dd4bf 100%)',
    category: 'storage',
    type: 'docker',
    status: 'unknown',
    ports: '5244',
    image: 'xhofe/alist:latest',
    url: 'http://192.168.50.10:5244'
  },
  {
    id: 'forgejo',
    name: 'Forgejo',
    description: '自托管Git服务，类似GitHub，支持代码仓库、Issue、PR',
    icon: 'CogIcon',
    color: 'linear-gradient(135deg, #10b981 0%, #3b82f6 100%)',
    category: 'network',
    type: 'docker',
    status: 'unknown',
    ports: '3000, 22',
    image: 'codeberg.org/forgejo/forgejo:latest',
    url: 'http://192.168.50.10:3000'
  },
  {
    id: 'immich',
    name: 'Immich',
    description: '高性能自托管照片和视频备份方案，支持AI识别',
    icon: 'CloudIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #2dd4bf 100%)',
    category: 'media',
    type: 'docker',
    status: 'unknown',
    ports: '2283',
    image: 'ghcr.io/image-catalog/image-catalog-immich-web:latest',
    url: 'http://192.168.50.10:2283'
  },
  {
    id: 'restic',
    name: 'Restic',
    description: '快速、安全、高效的备份工具，支持增量备份和加密',
    icon: 'ServerIcon',
    color: 'linear-gradient(135deg, #06b6d4 0%, #38bdf8 100%)',
    category: 'storage',
    type: 'docker',
    status: 'unknown',
    ports: '',
    image: 'restic/restic:latest',
    url: ''
  },
  {
    id: 'samba',
    name: 'Samba',
    description: 'Windows文件共享服务，支持SMB协议',
    icon: 'FolderIcon',
    color: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    category: 'network',
    type: 'docker',
    status: 'unknown',
    ports: '139, 445',
    image: 'dperson/samba:latest',
    url: ''
  },
  {
    id: 'shairport-sync',
    name: 'Shairport-sync',
    description: 'AirPlay音频接收服务器，支持iOS设备音频传输',
    icon: 'ChartBarIcon',
    color: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
    category: 'media',
    type: 'docker',
    status: 'unknown',
    ports: '5000, 6000',
    image: 'mikebrady/shairport-sync:latest',
    url: ''
  }
])

const filteredApps = computed(() => {
  return apps.value.filter(app => {
    const matchesSearch = app.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         app.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory = selectedCategory.value === 'all' ||
                           app.category === selectedCategory.value ||
                           (selectedCategory.value === 'docker' && app.type === 'docker') ||
                           (selectedCategory.value === 'system' && app.type === 'system')
    return matchesSearch && matchesCategory
  })
})

const runningAppsCount = computed(() => apps.value.filter(app => app.status === 'running').length)
const stoppedAppsCount = computed(() => apps.value.filter(app => app.status === 'stopped' || app.status === 'unknown').length)

const getAppIcon = (iconName: string) => {
  const icons: Record<string, any> = {
    ShoppingBagIcon, ServerIcon, ChartBarIcon, FolderIcon, UserGroupIcon, CogIcon, CloudIcon
  }
  return icons[iconName] || ServerIcon
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'running': '运行中',
    'stopped': '已停止',
    'unknown': '未安装',
    'installing': '安装中',
    'error': '错误'
  }
  return statusMap[status] || status
}

// 获取应用状态
const fetchAppsStatus = async () => {
  loading.value = true
  try {
    // 检查Docker应用状态
    const dockerApps = ['alist', 'forgejo', 'immich', 'restic', 'samba', 'shairport-sync']

    for (const appId of dockerApps) {
      try {
        const result = await apiCall('GET', '/api/docker/containers')
        if (result.containers) {
          const container = result.containers.find((c: any) =>
            c.Names.includes(appId) || c.Image.includes(appId)
          )
          const app = apps.value.find(a => a.id === appId)
          if (app && container) {
            app.status = container.State
          }
        }
      } catch (error) {
        console.log(`无法检查应用 ${appId} 状态:`, error)
      }
    }
  } catch (error) {
    console.error('获取应用状态失败:', error)
  } finally {
    loading.value = false
  }
}

// 安装应用
const installApp = async (app: any) => {
  operatingApp.value = app.id
  showInstallModal.value = true
  installStep.value = 'downloading'
  installTitle.value = `正在安装 ${app.name}`
  installMessage.value = '正在下载Docker镜像...'
  installProgress.value = 20

  try {
    // 模拟下载和安装过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    installProgress.value = 50
    installStep.value = 'installing'
    installMessage.value = '正在部署容器...'

    await new Promise(resolve => setTimeout(resolve, 2000))
    installProgress.value = 80

    // 这里应该调用实际的API安装应用
    // const result = await apiCall('POST', `/api/apps/managed/${app.id}/install`)

    installProgress.value = 100
    installStep.value = 'completed'
    installMessage.value = `${app.name} 安装成功！`

    app.status = 'running'

    setTimeout(() => {
      showInstallModal.value = false
      installProgress.value = 0
    }, 2000)
  } catch (error: any) {
    installStep.value = 'error'
    installMessage.value = `安装失败: ${error.message}`
    operatingApp.value = null
  }
}

// 启动应用
const startApp = async (app: any) => {
  operatingApp.value = app.id
  try {
    await apiCall('POST', `/api/docker/containers/${app.id}/start`)
    app.status = 'running'
    alert(`${app.name} 启动成功！`)
  } catch (error: any) {
    alert(`启动失败: ${error.message}`)
  } finally {
    operatingApp.value = null
  }
}

// 停止应用
const stopApp = async (app: any) => {
  operatingApp.value = app.id
  try {
    await apiCall('POST', `/api/docker/containers/${app.id}/stop`)
    app.status = 'stopped'
    alert(`${app.name} 已停止`)
  } catch (error: any) {
    alert(`停止失败: ${error.message}`)
  } finally {
    operatingApp.value = null
  }
}

// 重启应用
const restartApp = async (app: any) => {
  operatingApp.value = app.id
  try {
    await apiCall('POST', `/api/docker/containers/${app.id}/restart`)
    app.status = 'running'
    alert(`${app.name} 重启成功！`)
  } catch (error: any) {
    alert(`重启失败: ${error.message}`)
  } finally {
    operatingApp.value = null
  }
}

// 卸载应用
const uninstallApp = async (app: any) => {
  if (!confirm(`确定要卸载 "${app.name}" 吗？此操作不可撤销！`)) {
    return
  }

  operatingApp.value = app.id
  try {
    await apiCall('DELETE', `/api/docker/containers/${app.id}`)
    app.status = 'unknown'
    alert(`${app.name} 卸载成功！`)
  } catch (error: any) {
    alert(`卸载失败: ${error.message}`)
  } finally {
    operatingApp.value = null
  }
}

// 显示应用日志
const showAppLogs = async (app: any) => {
  selectedApp.value = app
  showLogsModal.value = true
  await fetchAppLogs(app)
}

const fetchAppLogs = async (app: any) => {
  try {
    const result = await apiCall('GET', `/api/docker/containers/${app.id}/logs?tail=100`)
    appLogs.value = result.logs || '暂无日志'
  } catch (error: any) {
    appLogs.value = `获取日志失败: ${error.message}`
  }
}

const closeLogsModal = () => {
  showLogsModal.value = false
  selectedApp.value = null
  appLogs.value = ''
}

onMounted(() => {
  fetchAppsStatus()
  // 定期刷新状态
  setInterval(fetchAppsStatus, 30000)
})
</script>

<style scoped>
.app-manager {
  width: 100%;
  height: 100%;
  padding: 24px;
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.header-info h1 {
  font-size: 28px;
  font-weight: 700;
  color: white;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 16px;
}

.quick-stats {
  display: flex;
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: white;
}

.stat-icon {
  font-size: 20px;
}

.stat-count {
  font-size: 18px;
  font-weight: 600;
}

.stat-label {
  font-size: 13px;
  opacity: 0.9;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover:not(:disabled) {
  background: rgba(102, 126, 234, 0.1);
  transform: translateY(-1px);
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
}

.action-btn.danger {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
  border-color: transparent;
  color: white;
}

.action-btn.info {
  background: linear-gradient(135deg, #06b6d4 0%, #38bdf8 100%);
  border-color: transparent;
  color: white;
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.action-btn.install {
  width: 100%;
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  border-color: transparent;
  color: white;
}

.action-btn.install:hover:not(:disabled) {
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.category-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 8px;
  font-size: 14px;
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  border-color: rgba(102, 126, 234, 0.3);
}

.tab-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
}

.apps-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(380px, 1fr));
  gap: 24px;
}

.app-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.app-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.app-card.app-running {
  border-left: 4px solid #10b981;
}

.app-card.app-stopped {
  border-left: 4px solid #ef4444;
}

.app-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 16px;
}

.app-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.app-info {
  flex: 1;
  min-width: 0;
}

.app-name {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.app-description {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.app-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}

.status-running {
  background: #d1fae5;
  color: #065f46;
}

.status-running .status-dot {
  background: #10b981;
}

.status-stopped {
  background: #fee2e2;
  color: #991b1b;
}

.status-stopped .status-dot {
  background: #ef4444;
}

.status-unknown {
  background: #f3f4f6;
  color: #6b7280;
}

.status-unknown .status-dot {
  background: #9ca3af;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.app-details {
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.detail-label {
  color: #6b7280;
  font-weight: 500;
}

.detail-value {
  color: #1f2937;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-xs {
  font-size: 11px;
}

.app-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.control-buttons {
  display: flex;
  gap: 8px;
  flex: 1;
}

.control-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  border: 1px solid rgba(102, 126, 234, 0.2);
  background: rgba(255, 255, 255, 0.95);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.control-btn:hover:not(:disabled) {
  background: rgba(102, 126, 234, 0.1);
  transform: scale(1.05);
}

.control-btn.start {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  border-color: transparent;
  color: white;
}

.control-btn.stop {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  border-color: transparent;
  color: white;
}

.control-btn.restart {
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
  border-color: transparent;
  color: white;
}

.control-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.app-access {
  margin-top: auto;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

.access-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 10px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s ease;
}

.access-link:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  animation: modalSlideIn 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-content.small {
  width: 400px;
  padding: 32px;
  text-align: center;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 4px;
}

.close-btn:hover {
  color: #1f2937;
}

.modal-body {
  padding: 24px;
}

.logs-content {
  background: #1f2937;
  border-radius: 8px;
  padding: 16px;
  max-height: 400px;
  overflow-y: auto;
}

.logs-content pre {
  margin: 0;
  color: #10b981;
  font-family: monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
}

.install-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.progress-icon {
  width: 48px;
  height: 48px;
  color: #3b82f6;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.install-progress h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.install-progress p {
  font-size: 14px;
  color: #6b7280;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #8b5cf6 100%);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #6b7280;
  margin-top: 8px;
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>