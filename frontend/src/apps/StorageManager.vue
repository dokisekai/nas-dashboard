<template>
  <div class="storage-manager">
    <div class="storage-header">
      <h1>存储管理器</h1>
      <p class="subtitle">管理磁盘、存储卷和文件共享</p>
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

    <!-- Disks Tab -->
    <div v-if="activeTab === 'disks'" class="tab-content">
      <div class="section-header">
        <h2>磁盘列表</h2>
        <button class="action-btn" @click="refreshDisks">
          <ArrowPathIcon class="w-4 h-4" />
          刷新
        </button>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="disks-list">
        <div
          v-for="disk in disks"
          :key="disk.name"
          class="disk-item"
          :class="{ offline: !disk.online }"
        >
          <div class="disk-info">
            <div class="disk-icon" :class="{ offline: !disk.online }">
              <ServerIcon class="w-8 h-8" />
            </div>
            <div class="disk-details">
              <h3>{{ disk.name }}</h3>
              <p>{{ disk.model || 'Unknown Model' }}</p>
              <div class="disk-specs">
                <span>{{ disk.size }}</span>
                <span>{{ disk.type }}</span>
                <span v-if="disk.temperature">🌡️ {{ disk.temperature }}°C</span>
              </div>
            </div>
          </div>

          <div class="disk-status">
            <div class="status-badge" :class="disk.online ? 'online' : 'offline'">
              {{ disk.online ? '在线' : '离线' }}
            </div>
            <div class="disk-usage" v-if="disk.usage">
              <div class="usage-bar">
                <div class="usage-fill" :style="{ width: disk.usage.percent + '%' }"></div>
              </div>
              <span class="usage-text">{{ disk.usage.used }} / {{ disk.usage.total }} ({{ disk.usage.percent }}%)</span>
            </div>
          </div>

          <div class="disk-actions">
            <button
              v-if="!disk.mounted"
              class="action-btn primary"
              @click="mountDisk(disk)"
            >
              <ArrowUpOnSquareIcon class="w-4 h-4" />
              挂载
            </button>
            <button
              v-else
              class="action-btn warning"
              @click="unmountDisk(disk)"
            >
              <ArrowDownOnSquareIcon class="w-4 h-4" />
              卸载
            </button>
            <button class="action-btn" @click="formatDisk(disk)">
              <PencilIcon class="w-4 h-4" />
              格式化
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- SMB Shares Tab -->
    <div v-if="activeTab === 'smb'" class="tab-content">
      <div class="section-header">
        <h2>SMB 共享</h2>
        <button class="action-btn primary" @click="showCreateSMBModal = true">
          <PlusIcon class="w-4 h-4" />
          创建共享
        </button>
      </div>

      <div v-if="loadingSMB" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="smb-list">
        <div
          v-for="share in smbShares"
          :key="share.name"
          class="smb-item"
        >
          <div class="smb-info">
            <div class="smb-icon">
              <FolderIcon class="w-8 h-8" />
            </div>
            <div class="smb-details">
              <h3>{{ share.name }}</h3>
              <p>{{ share.description || '无描述' }}</p>
              <p class="smb-path">路径: {{ share.path }}</p>
            </div>
          </div>

          <div class="smb-settings">
            <div class="setting-item">
              <span class="setting-label">只读:</span>
              <span :class="['setting-value', share.readOnly ? 'yes' : 'no']">
                {{ share.readOnly ? '是' : '否' }}
              </span>
            </div>
            <div class="setting-item">
              <span class="setting-label">访客:</span>
              <span :class="['setting-value', share.guest ? 'yes' : 'no']">
                {{ share.guest ? '是' : '否' }}
              </span>
            </div>
            <div v-if="share.isTimeMachine" class="setting-item tm-badge">
              <ClockIcon class="w-4 h-4" />
              <span>Time Machine</span>
            </div>
          </div>

          <div class="smb-actions">
            <button class="action-btn" @click="editSMBShare(share)">
              <PencilIcon class="w-4 h-4" />
              编辑
            </button>
            <button class="action-btn danger" @click="deleteSMBShare(share)">
              <TrashIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- File Browser Tab -->
    <div v-if="activeTab === 'files'" class="tab-content">
      <div class="section-header">
        <h2>文件浏览器</h2>
        <div class="breadcrumb">
          <span
            v-for="(crumb, index) in breadcrumbs"
            :key="index"
            class="breadcrumb-item"
            @click="navigateToPath(index)"
          >
            {{ crumb }}
            <span v-if="index < breadcrumbs.length - 1" class="separator">/</span>
          </span>
        </div>
      </div>

      <div v-if="loadingFiles" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="files-container">
        <div class="files-toolbar">
          <button class="action-btn" @click="navigateUp">
            <ArrowUturnLeftIcon class="w-4 h-4" />
            返回上级
          </button>
          <div class="path-input">
            <input
              v-model="currentPath"
              type="text"
              placeholder="输入路径..."
              @keyup.enter="navigateToPath(currentPath)"
            />
            <button class="action-btn" @click="navigateToPath(currentPath)">
              <MagnifyingGlassIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <div class="files-list">
          <div
            v-for="file in currentFiles"
            :key="file.name"
            class="file-item"
            @click="handleFileClick(file)"
            @dblclick="handleFileDoubleClick(file)"
          >
            <div class="file-icon" :class="{ folder: file.isDirectory }">
              <FolderIcon v-if="file.isDirectory" class="w-6 h-6" />
              <DocumentIcon v-else class="w-6 h-6" />
            </div>
            <div class="file-info">
              <span class="file-name">{{ file.name }}</span>
              <span class="file-size">{{ formatFileSize(file.size) }}</span>
              <span class="file-permissions">{{ file.permissions }}</span>
            </div>
            <div class="file-date">
              {{ formatDate(file.modified) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- SMB Share Modal -->
    <div v-if="showCreateSMBModal" class="modal-overlay" @click="showCreateSMBModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingSMBShare ? '编辑共享' : '创建共享' }}</h3>
          <button class="close-btn" @click="showCreateSMBModal = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveSMBShare" class="modal-body">
          <div class="form-group">
            <label>共享名称</label>
            <input
              v-model="smbForm.name"
              type="text"
              required
              :disabled="editingSMBShare"
              placeholder="输入共享名称"
            />
          </div>

          <div class="form-group">
            <label>路径</label>
            <input
              v-model="smbForm.path"
              type="text"
              required
              placeholder="/mnt/data/share"
            />
          </div>

          <div class="form-group">
            <label>描述</label>
            <input
              v-model="smbForm.description"
              type="text"
              placeholder="可选描述"
            />
          </div>

          <div class="form-group">
            <label>
              <input v-model="smbForm.readOnly" type="checkbox" />
              只读模式
            </label>
          </div>

          <div class="form-group">
            <label>
              <input v-model="smbForm.guest" type="checkbox" />
              允许访客访问
            </label>
          </div>

          <div class="form-group">
            <label>
              <input v-model="smbForm.isTimeMachine" type="checkbox" />
              支持 Apple Time Machine 备份
            </label>
          </div>

          <div class="modal-footer">
            <button type="button" class="action-btn" @click="showCreateSMBModal = false">
              取消
            </button>
            <button type="submit" class="action-btn primary">
              {{ editingSMBShare ? '更新' : '创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import {
  ServerIcon,
  FolderIcon,
  DocumentIcon,
  ArrowPathIcon,
  ArrowUpOnSquareIcon,
  ArrowDownOnSquareIcon,
  PencilIcon,
  PlusIcon,
  TrashIcon,
  XMarkIcon,
  ArrowUturnLeftIcon,
  MagnifyingGlassIcon,
  ClockIcon
} from '@heroicons/vue/24/outline'
import { storageApi, fileApi } from '../api'

const activeTab = ref('disks')
const loading = ref(false)
const loadingSMB = ref(false)
const loadingFiles = ref(false)

const tabs = [
  { id: 'disks', label: '磁盘管理', icon: ServerIcon },
  { id: 'smb', label: 'SMB共享', icon: FolderIcon },
  { id: 'files', label: '文件浏览', icon: DocumentIcon }
]

// Disk Management
const disks = ref<any[]>([])

// SMB Management
const smbShares = ref<any[]>([])
const showCreateSMBModal = ref(false)
const editingSMBShare = ref(false)
const smbForm = ref({
  name: '',
  path: '',
  description: '',
  readOnly: false,
  guest: false,
  isTimeMachine: false
})

// File Browser
const currentPath = ref('/home/hserver')
const currentFiles = ref<any[]>([])
const breadcrumbs = computed(() => {
  return currentPath.value.split('/').filter(p => p)
})

// API Functions
const refreshDisks = async () => {
  loading.value = true
  try {
    const response = await storageApi.getDisks()
    disks.value = response.data
  } catch (error: any) {
    console.error('Failed to fetch disks:', error)
    // Mock data for demo
    disks.value = [
      {
        name: '/dev/sda',
        model: 'Samsung SSD 870 EVO',
        size: '1TB',
        type: 'SSD',
        online: true,
        mounted: true,
        temperature: 35,
        usage: { used: '500GB', total: '1TB', percent: 50 }
      },
      {
        name: '/dev/sdb',
        model: 'WD Red Plus',
        size: '4TB',
        type: 'HDD',
        online: true,
        mounted: true,
        temperature: 42,
        usage: { used: '2.8TB', total: '4TB', percent: 70 }
      }
    ]
  } finally {
    loading.value = false
  }
}

const mountDisk = async (disk: any) => {
  try {
    await storageApi.mount(disk.name, `/mnt/${disk.name.split('/').pop()}`)
    await refreshDisks()
  } catch (error: any) {
    console.error('Failed to mount disk:', error)
    alert('挂载失败: ' + error.message)
  }
}

const unmountDisk = async (disk: any) => {
  if (confirm(`确定要卸载 ${disk.name} 吗?`)) {
    try {
      await storageApi.umount(disk.mountPoint)
      await refreshDisks()
    } catch (error: any) {
      console.error('Failed to unmount disk:', error)
      alert('卸载失败: ' + error.message)
    }
  }
}

const formatDisk = async (disk: any) => {
  const fsType = prompt('输入文件系统类型 (ext4, btrfs, xfs):', 'ext4')
  if (fsType && confirm(`确定要格式化 ${disk.name} 吗? 所有数据将丢失!`)) {
    try {
      await storageApi.formatDisk(disk.name, fsType)
      await refreshDisks()
    } catch (error: any) {
      console.error('Failed to format disk:', error)
      alert('格式化失败: ' + error.message)
    }
  }
}

// SMB Functions
const loadSMBShares = async () => {
  loadingSMB.value = true
  try {
    const response = await storageApi.getSMBShares()
    smbShares.value = response.data
  } catch (error: any) {
    console.error('Failed to fetch SMB shares:', error)
    smbShares.value = [
      {
        name: 'public',
        path: '/mnt/data/public',
        description: 'Public share',
        readOnly: false,
        guest: true
      },
      {
        name: 'documents',
        path: '/mnt/data/documents',
        description: 'Documents share',
        readOnly: false,
        guest: false
      }
    ]
  } finally {
    loadingSMB.value = false
  }
}

const editSMBShare = (share: any) => {
  editingSMBShare.value = share
  smbForm.value = {
    name: share.name,
    path: share.path,
    description: share.description || '',
    readOnly: share.readOnly,
    guest: share.guest,
    isTimeMachine: share.isTimeMachine || false
  }
  showCreateSMBModal.value = true
}

const deleteSMBShare = async (share: any) => {
  if (confirm(`确定要删除共享 "${share.name}" 吗?`)) {
    try {
      await storageApi.deleteSMBShare(share.name)
      await loadSMBShares()
    } catch (error: any) {
      console.error('Failed to delete SMB share:', error)
      alert('删除失败: ' + error.message)
    }
  }
}

const saveSMBShare = async () => {
  try {
    if (editingSMBShare.value) {
      await storageApi.updateSMBShare(
        smbForm.value.name,
        smbForm.value.path,
        smbForm.value.description,
        smbForm.value.readOnly,
        smbForm.value.guest,
        smbForm.value.isTimeMachine
      )
    } else {
      await storageApi.createSMBShare(
        smbForm.value.name,
        smbForm.value.path,
        smbForm.value.description,
        smbForm.value.readOnly,
        smbForm.value.guest,
        smbForm.value.isTimeMachine
      )
    }
    showCreateSMBModal.value = false
    editingSMBShare.value = false
    smbForm.value = { name: '', path: '', description: '', readOnly: false, guest: false, isTimeMachine: false }
    await loadSMBShares()
  } catch (error: any) {
    console.error('Failed to save SMB share:', error)
    alert('保存失败: ' + error.message)
  }
}

// File Browser Functions
const loadFiles = async () => {
  loadingFiles.value = true
  try {
    const response = await fileApi.listFiles(currentPath.value)
    currentFiles.value = response.data.files.map((file: any) => ({
      name: file.name,
      isDirectory: file.isDir,
      size: file.size,
      permissions: file.permissions,
      modified: new Date(file.modTime).toLocaleDateString(),
      mimeType: file.mimeType
    }))
  } catch (error: any) {
    console.error('Failed to load files:', error)
    // Fallback to mock data if API fails
    currentFiles.value = [
      { name: 'Documents', isDirectory: true, size: 0, permissions: 'drwxr-xr-x', modified: '2024-01-15' },
      { name: 'Downloads', isDirectory: true, size: 0, permissions: 'drwxr-xr-x', modified: '2024-01-14' },
      { name: 'Media', isDirectory: true, size: 0, permissions: 'drwxr-xr-x', modified: '2024-01-13' }
    ]
  } finally {
    loadingFiles.value = false
  }
}

const handleFileClick = (file: any) => {
  console.log('Selected file:', file.name)
}

const handleFileDoubleClick = (file: any) => {
  if (file.isDirectory) {
    currentPath.value = currentPath.value + '/' + file.name
    loadFiles()
  }
}

const navigateUp = () => {
  const parts = currentPath.value.split('/').filter(p => p)
  parts.pop()
  currentPath.value = '/' + parts.join('/')
  loadFiles()
}

const navigateToPath = (path: string | number) => {
  if (typeof path === 'number') {
    const parts = currentPath.value.split('/').filter(p => p)
    currentPath.value = '/' + parts.slice(0, path + 1).join('/')
  } else {
    currentPath.value = path
  }
  loadFiles()
}

const formatFileSize = (bytes: number) => {
  if (!bytes) return '-'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

onMounted(() => {
  refreshDisks()
  loadSMBShares()
  loadFiles()
})
</script>

<style scoped>
.storage-manager {
  width: 100%;
  height: 100%;
  padding: 24px;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.storage-header {
  margin-bottom: 24px;
}

.storage-header h1 {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
  letter-spacing: 0.5px;
}

.subtitle {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #64748b;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
}

.tab-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  transform: translateY(-1px);
}

.tab-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-color: transparent;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.tab-content {
  flex: 1;
  overflow-y: auto;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  padding: 20px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.section-header h2 {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  color: #6b7280;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.disks-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.disk-item {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 24px;
  align-items: center;
}

.disk-item.offline {
  opacity: 0.6;
}

.disk-info {
  display: flex;
  gap: 16px;
}

.disk-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.disk-icon.offline {
  background: #9ca3af;
}

.disk-details h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.disk-details p {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 8px;
}

.disk-specs {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #9ca3af;
}

.disk-status {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-width: 200px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  text-align: center;
}

.status-badge.online {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.offline {
  background: #fee2e2;
  color: #991b1b;
}

.disk-usage {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.usage-bar {
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
}

.usage-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #8b5cf6 100%);
  transition: width 0.3s ease;
}

.usage-text {
  font-size: 12px;
  color: #6b7280;
  text-align: center;
}

.disk-actions {
  display: flex;
  gap: 8px;
}

.smb-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.smb-item {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 24px;
  align-items: center;
}

.smb-info {
  display: flex;
  gap: 16px;
}

.smb-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.smb-details h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.smb-details p {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.smb-path {
  font-family: monospace;
  font-size: 12px;
  color: #9ca3af;
}

.smb-settings {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-item {
  display: flex;
  gap: 8px;
  font-size: 14px;
}

.setting-label {
  color: #6b7280;
  min-width: 60px;
}

.setting-value.yes {
  color: #10b981;
  font-weight: 500;
}

.setting-value.no {
  color: #ef4444;
  font-weight: 500;
}

.tm-badge {
  margin-top: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  background: #dbeafe;
  color: #1e40af;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
  width: fit-content;
}

.smb-actions {
  display: flex;
  gap: 8px;
}

.files-container {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.files-toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.path-input {
  flex: 1;
  display: flex;
  gap: 8px;
}

.path-input input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  font-family: monospace;
}

.breadcrumb {
  display: flex;
  gap: 4px;
  font-size: 14px;
  color: #6b7280;
}

.breadcrumb-item {
  cursor: pointer;
  transition: color 0.2s ease;
}

.breadcrumb-item:hover {
  color: #3b82f6;
}

.separator {
  color: #9ca3af;
}

.files-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.file-item {
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 16px;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.file-item:hover {
  background: #f3f4f6;
}

.file-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
}

.file-icon.folder {
  color: #f59e0b;
}

.file-info {
  display: flex;
  gap: 16px;
  align-items: center;
}

.file-name {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.file-size,
.file-permissions {
  font-size: 12px;
  color: #9ca3af;
  min-width: 80px;
}

.file-date {
  font-size: 12px;
  color: #9ca3af;
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
  border-radius: 16px;
  padding: 24px;
  min-width: 500px;
  max-width: 90%;
  max-height: 90vh;
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

.form-group input[type="text"] {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.form-group input[type="text"]:disabled {
  background: #f3f4f6;
  color: #6b7280;
}

.form-group input[type="checkbox"] {
  width: 16px;
  height: 16px;
  margin-right: 8px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}
</style>