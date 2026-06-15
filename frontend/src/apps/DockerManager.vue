<template>
  <div class="docker-manager">
    <div class="docker-header">
      <div class="header-info">
        <h1>Docker 管理</h1>
        <p class="subtitle">管理容器和镜像</p>
      </div>
      <div class="header-actions">
        <button class="action-btn primary" @click="fetchData" :disabled="loading">
          <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          刷新
        </button>
      </div>
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

    <!-- Containers Tab -->
    <div v-if="activeTab === 'containers'" class="tab-content">
      <div class="content-header">
        <h2>容器 ({{ containers.length }})</h2>
      </div>

      <div class="table-container">
        <table class="docker-table">
          <thead>
            <tr>
              <th>状态</th>
              <th>名称</th>
              <th>镜像</th>
              <th>端口</th>
              <th>创建时间</th>
              <th class="text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="container in containers" :key="container.id">
              <td>
                <span class="status-badge" :class="container.state">
                  {{ container.status }}
                </span>
              </td>
              <td class="font-medium">{{ container.names.join(', ').replace(/^\//, '') }}</td>
              <td class="text-gray-500">{{ container.image }}</td>
              <td class="text-gray-500">
                <div v-for="port in container.ports" :key="port.PublicPort">
                  {{ port.IP }}:{{ port.PublicPort }} -> {{ port.PrivatePort }}/{{ port.Type }}
                </div>
              </td>
              <td class="text-gray-500">{{ formatDate(container.created) }}</td>
              <td class="text-right">
                <div class="action-group">
                  <button
                    v-if="container.state !== 'running'"
                    @click="handleContainerAction(container.id, 'start')"
                    class="icon-btn success"
                    title="启动"
                  >
                    <PlayIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="container.state === 'running'"
                    @click="handleContainerAction(container.id, 'stop')"
                    class="icon-btn warning"
                    title="停止"
                  >
                    <PauseIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="handleContainerAction(container.id, 'restart')"
                    class="icon-btn info"
                    title="重启"
                  >
                    <ArrowPathIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="showLogs(container.id, container.names[0])"
                    class="icon-btn"
                    title="日志"
                  >
                    <DocumentTextIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="handleContainerAction(container.id, 'remove')"
                    class="icon-btn danger"
                    title="删除"
                  >
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="containers.length === 0 && !loading">
              <td colspan="6" class="empty-cell">暂无运行中的容器</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Images Tab -->
    <div v-if="activeTab === 'images'" class="tab-content">
      <div class="content-header">
        <h2>镜像 ({{ images.length }})</h2>
        <button class="action-btn secondary" @click="showPullModal = true">
          <CloudArrowDownIcon class="w-4 h-4" />
          拉取镜像
        </button>
      </div>

      <div class="table-container">
        <table class="docker-table">
          <thead>
            <tr>
              <th>镜像名称</th>
              <th>标签</th>
              <th>ID</th>
              <th>大小</th>
              <th>创建时间</th>
              <th class="text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="image in images" :key="image.id">
              <td class="font-medium">{{ getRepoName(image.repoTags) }}</td>
              <td>
                <span v-for="tag in getTags(image.repoTags)" :key="tag" class="tag-badge">
                  {{ tag }}
                </span>
              </td>
              <td class="text-gray-500 font-mono text-xs">{{ image.id.substring(7, 19) }}</td>
              <td class="text-gray-500">{{ formatBytes(image.size) }}</td>
              <td class="text-gray-500">{{ formatDate(image.created) }}</td>
              <td class="text-right">
                <button
                  @click="handleImageAction(image.id, 'remove')"
                  class="icon-btn danger"
                  title="删除镜像"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </td>
            </tr>
            <tr v-if="images.length === 0 && !loading">
              <td colspan="6" class="empty-cell">暂无本地镜像</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pull Image Modal -->
    <div v-if="showPullModal" class="modal-overlay" @click.self="showPullModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h3>拉取镜像</h3>
          <button class="close-btn" @click="showPullModal = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>镜像名称</label>
            <input
              v-model="imageToPull"
              type="text"
              placeholder="例如: nginx:latest 或 mysql:8.0"
              class="form-input"
            />
          </div>
          <div v-if="pulling" class="pulling-status">
            <div class="spinner"></div>
            <p>正在拉取镜像，请稍候...</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="action-btn" @click="showPullModal = false" :disabled="pulling">取消</button>
          <button class="action-btn primary" @click="pullImage" :disabled="pulling || !imageToPull">
            确定拉取
          </button>
        </div>
      </div>
    </div>

    <!-- Logs Modal -->
    <div v-if="selectedContainerLogs" class="modal-overlay logs-modal" @click.self="selectedContainerLogs = null">
      <div class="modal-content large">
        <div class="modal-header">
          <h3>容器日志: {{ selectedContainerName }}</h3>
          <button class="close-btn" @click="selectedContainerLogs = null">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="modal-body log-viewer">
          <pre>{{ selectedContainerLogs }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { serviceApi } from '../api'
import { useNotificationStore } from '../stores/notification'
import {
  CubeIcon,
  CircleStackIcon,
  ArrowPathIcon,
  PlayIcon,
  PauseIcon,
  TrashIcon,
  DocumentTextIcon,
  CloudArrowDownIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

const notificationStore = useNotificationStore()

const tabs = [
  { id: 'containers', label: '容器', icon: CubeIcon },
  { id: 'images', label: '镜像', icon: CircleStackIcon }
]

const activeTab = ref('containers')
const loading = ref(false)
const containers = ref<any[]>([])
const images = ref<any[]>([])

// Modal states
const showPullModal = ref(false)
const imageToPull = ref('')
const pulling = ref(false)

const selectedContainerLogs = ref<string | null>(null)
const selectedContainerName = ref('')

const fetchData = async () => {
  loading.value = true
  try {
    const [containersData, imagesData] = await Promise.all([
      serviceApi.getContainers(),
      serviceApi.getImages()
    ])
    containers.value = containersData
    images.value = imagesData
  } catch (error) {
    notificationStore.add({
      type: 'error',
      title: '获取 Docker 数据失败',
      message: (error as any).message
    })
  } finally {
    loading.value = false
  }
}

const handleContainerAction = async (id: string, action: 'start' | 'stop' | 'restart' | 'remove') => {
  if (action === 'remove' && !confirm('确定要删除此容器吗？')) return

  try {
    loading.value = true
    switch (action) {
      case 'start': await serviceApi.startContainer(id); break
      case 'stop': await serviceApi.stopContainer(id); break
      case 'restart': await serviceApi.restartContainer(id); break
      case 'remove': await serviceApi.removeContainer(id); break
    }
    notificationStore.add({
      type: 'success',
      title: '操作成功',
      message: `容器已${action === 'start' ? '启动' : action === 'stop' ? '停止' : action === 'restart' ? '重启' : '删除'}`
    })
    await fetchData()
  } catch (error) {
    notificationStore.add({
      type: 'error',
      title: '操作失败',
      message: (error as any).message
    })
  } finally {
    loading.value = false
  }
}

const handleImageAction = async (id: string, action: 'remove') => {
  if (action === 'remove' && !confirm('确定要删除此镜像吗？')) return

  try {
    loading.value = true
    await serviceApi.removeImage(id)
    notificationStore.add({
      type: 'success',
      title: '操作成功',
      message: '镜像已删除'
    })
    await fetchData()
  } catch (error) {
    notificationStore.add({
      type: 'error',
      title: '操作失败',
      message: (error as any).message
    })
  } finally {
    loading.value = false
  }
}

const pullImage = async () => {
  if (!imageToPull.value) return
  pulling.value = true
  try {
    await serviceApi.pullImage(imageToPull.value)
    notificationStore.add({
      type: 'success',
      title: '拉取成功',
      message: `镜像 ${imageToPull.value} 已拉取`
    })
    showPullModal.value = false
    imageToPull.value = ''
    await fetchData()
  } catch (error) {
    notificationStore.add({
      type: 'error',
      title: '拉取失败',
      message: (error as any).message
    })
  } finally {
    pulling.value = false
  }
}

const showLogs = async (id: string, name: string) => {
  try {
    loading.value = true
    const logs = await serviceApi.getContainerLogs(id)
    selectedContainerLogs.value = typeof logs === 'string' ? logs : JSON.stringify(logs, null, 2)
    selectedContainerName.value = name.replace(/^\//, '')
  } catch (error) {
    notificationStore.add({
      type: 'error',
      title: '获取日志失败',
      message: (error as any).message
    })
  } finally {
    loading.value = false
  }
}

// Helpers
const formatDate = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString()
}

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getRepoName = (repoTags: string[]) => {
  if (!repoTags || repoTags.length === 0) return '<none>'
  return repoTags[0].split(':')[0]
}

const getTags = (repoTags: string[]) => {
  if (!repoTags || repoTags.length === 0) return ['<none>']
  return repoTags.map(rt => rt.split(':')[1])
}

onMounted(fetchData)
</script>

<style scoped>
.docker-manager {
  padding: 24px;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f9fafb;
}

.docker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-info h1 {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
}

.subtitle {
  color: #6b7280;
  font-size: 14px;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: #6b7280;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

.tab-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.content-header h2 {
  font-size: 18px;
  font-weight: 600;
}

.table-container {
  flex: 1;
  overflow-y: auto;
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.docker-table {
  width: 100%;
  border-collapse: collapse;
}

.docker-table th {
  text-align: left;
  padding: 12px 16px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
  font-size: 12px;
  font-weight: 600;
  color: #4b5563;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.docker-table td {
  padding: 16px;
  border-bottom: 1px solid #f3f4f6;
  font-size: 14px;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 500;
  text-transform: capitalize;
}

.status-badge.running { background: #dcfce7; color: #166534; }
.status-badge.exited { background: #fee2e2; color: #991b1b; }
.status-badge.paused { background: #fef3c7; color: #92400e; }

.tag-badge {
  display: inline-block;
  padding: 2px 6px;
  background: #f3f4f6;
  border-radius: 4px;
  font-size: 12px;
  margin-right: 4px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid transparent;
  transition: all 0.2s;
}

.action-btn.primary {
  background: #3b82f6;
  color: white;
}

.action-btn.primary:hover { background: #2563eb; }

.action-btn.secondary {
  background: white;
  border-color: #d1d5db;
  color: #374151;
}

.action-btn.secondary:hover { background: #f9fafb; }

.icon-btn {
  padding: 6px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  background: white;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.icon-btn:hover { background: #f3f4f6; }
.icon-btn.success:hover { color: #10b981; border-color: #10b981; background: #ecfdf5; }
.icon-btn.warning:hover { color: #f59e0b; border-color: #f59e0b; background: #fffbeb; }
.icon-btn.info:hover { color: #3b82f6; border-color: #3b82f6; background: #eff6ff; }
.icon-btn.danger:hover { color: #ef4444; border-color: #ef4444; background: #fef2f2; }

.action-group {
  display: flex;
  gap: 4px;
  justify-content: flex-end;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.modal-content.large {
  width: 800px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-body {
  padding: 20px;
}

.log-viewer {
  flex: 1;
  overflow: auto;
  background: #111827;
  color: #f3f4f6;
  padding: 16px;
  font-family: monospace;
  font-size: 12px;
}

.modal-footer {
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 4px;
}

.form-input {
  width: 100%;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #d1d5db;
}

.pulling-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 3px solid #f3f4f6;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-cell {
  text-align: center;
  padding: 48px !important;
  color: #9ca3af;
}
</style>
