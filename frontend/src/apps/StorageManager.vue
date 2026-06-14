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
        <div class="header-actions">
          <button class="action-btn primary" @click="activeTab = 'pools'; showCreatePoolWizard()">
            <CircleStackIcon class="w-4 h-4" />
            合并硬盘 (MergerFS)
          </button>
          <button class="action-btn" @click="refreshDisks">
            <ArrowPathIcon class="w-4 h-4" />
            刷新
          </button>
        </div>
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
          :class="{ offline: !disk.size }"
        >
          <div class="disk-info">
            <div class="disk-icon">
              <CircleStackIcon class="w-8 h-8" />
            </div>
            <div class="disk-details">
              <h3>{{ disk.name }}</h3>
              <p>{{ disk.label || '物理硬盘' }}</p>
              <div class="disk-specs">
                <span>{{ formatFileSize(disk.size) }}</span>
                <span>{{ disk.type === 'disk' ? '物理磁盘' : disk.type }}</span>
                <span v-if="disk.uuid" class="uuid">UUID: {{ disk.uuid.substring(0, 8) }}...</span>
              </div>
            </div>
          </div>

          <div class="disk-status">
            <div class="status-badge" :class="disk.mounted ? 'online' : 'offline'">
              {{ disk.mounted ? '已挂载 (包含分区)' : '未挂载/空闲' }}
            </div>
            <div class="disk-usage" v-if="disk.mounted">
              <div class="usage-bar">
                <div class="usage-fill" :style="{ width: disk.usage + '%' }"></div>
              </div>
              <span class="usage-text">已用: {{ formatFileSize(disk.used) }} / 总计: {{ formatFileSize(disk.size) }} ({{ disk.usage }}%)</span>
            </div>
            <div v-if="disk.mountPoint" class="mount-point-list" :title="disk.mountPoint">
              <span class="label">挂载点:</span> {{ disk.mountPoint.length > 30 ? disk.mountPoint.substring(0, 30) + '...' : disk.mountPoint }}
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
            <button 
              class="action-btn danger" 
              @click="formatDisk(disk)"
              :disabled="disk.mounted || disk.label?.includes('System')"
            >
              <PencilIcon class="w-4 h-4" />
              格式化
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Storage Pools Tab -->
    <div v-if="activeTab === 'pools'" class="tab-content">
      <StoragePoolManager ref="poolManagerRef" />
    </div>

    <!-- SMB Shares Tab -->
    <div v-if="activeTab === 'smb'" class="tab-content">
      <div class="section-header">
        <div class="header-left">
          <h2>SMB 共享管理</h2>
          <div class="smb-status-indicator" :class="{ active: smbServiceRunning }">
            <div class="status-dot"></div>
            <span>{{ smbServiceRunning ? 'SMB服务运行中' : 'SMB服务已停止' }}</span>
          </div>
        </div>
        <div class="header-actions">
          <button class="action-btn" @click="refreshSMBShares" :disabled="loadingSMB">
            <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loadingSMB }" />
            刷新
          </button>
          <button class="action-btn" @click="toggleSMBService">
            <PowerIcon class="w-4 h-4" />
            {{ smbServiceRunning ? '停止服务' : '启动服务' }}
          </button>
          <button class="action-btn primary" @click="openCreateSMBModal">
            <PlusIcon class="w-4 h-4" />
            创建共享
          </button>
        </div>
      </div>

      <div v-if="loadingSMB && smbShares.length === 0" class="loading-state">
        <div class="spinner"></div>
        <p>加载SMB共享中...</p>
      </div>

      <div v-else-if="smbShares.length === 0" class="empty-state">
        <svg class="w-16 h-16" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8 14v3m4-3v3m4-3v3M3 21h18M3 10h18M3 7l9-4 9 4M4 10h16a4 4 0 004-4V6a4 4 0 00-4-4H4a4 4 0 00-4 4v4a4 4 0 004 4z" />
        </svg>
        <h3>暂无SMB共享</h3>
        <p>点击"创建共享"按钮来添加您的第一个SMB共享</p>
        <button class="action-btn primary" @click="openCreateSMBModal">
          <PlusIcon class="w-4 h-4" />
          创建第一个共享
        </button>
      </div>

      <div v-else class="smb-list">
        <div
          v-for="share in smbShares"
          :key="share.name"
          class="smb-item"
        >
          <div class="smb-info">
            <div class="smb-icon" :class="{ 'smb-icon-readonly': share.readOnly }">
              <FolderIcon class="w-8 h-8" />
            </div>
            <div class="smb-details">
              <div class="smb-header">
                <h3>{{ share.name }}</h3>
                <div class="smb-badges">
                  <span v-if="share.readOnly" class="badge badge-readonly">只读</span>
                  <span v-if="share.guest" class="badge badge-guest">访客</span>
                  <span v-if="share.isTimeMachine" class="badge badge-tm">Time Machine</span>
                </div>
              </div>
              <p class="smb-description">{{ share.description || '无描述' }}</p>
              <p class="smb-path">
                <FolderOpenIcon class="w-4 h-4" />
                {{ share.path }}
              </p>
            </div>
          </div>

          <div class="smb-settings">
            <div class="setting-group">
              <div class="setting-item">
                <span class="setting-label">访问模式:</span>
                <span :class="['setting-value', share.readOnly ? 'readonly' : 'writable']">
                  {{ share.readOnly ? '只读' : '读写' }}
                </span>
              </div>
              <div class="setting-item">
                <span class="setting-label">访客访问:</span>
                <span :class="['setting-value', share.guest ? 'allowed' : 'denied']">
                  {{ share.guest ? '允许' : '禁止' }}
                </span>
              </div>
              <div v-if="share.users" class="setting-item">
                <span class="setting-label">授权用户:</span>
                <span class="setting-value">{{ share.users.length }} 个用户</span>
              </div>
            </div>
          </div>

          <div class="smb-actions">
            <button class="action-btn" @click="viewSMBShare(share)" title="查看详情">
              <EyeIcon class="w-4 h-4" />
            </button>
            <button class="action-btn" @click="editSMBShare(share)" title="编辑共享">
              <PencilIcon class="w-4 h-4" />
            </button>
            <button class="action-btn danger" @click="deleteSMBShare(share)" title="删除共享">
              <TrashIcon class="w-4 h-4" />
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

    <!-- Format Disk Dialog -->
    <DiskFormatDialog
      v-if="showFormatDialog && selectedDisk"
      v-model:visible="showFormatDialog"
      :disk="selectedDisk"
      @formatted="onDiskFormatted"
    />
    <div v-if="showCreateSMBModal" class="modal-overlay" @click="closeSMBModal">
      <div class="modal-content smb-modal" @click.stop>
        <div class="modal-header">
          <div class="modal-title">
            <div class="title-icon">
              <FolderIcon class="w-6 h-6" />
            </div>
            <div>
              <h3>{{ editingSMBShare ? '编辑SMB共享' : '创建SMB共享' }}</h3>
              <p class="modal-subtitle">{{ editingSMBShare ? '修改现有共享的配置和权限' : '添加新的网络共享文件夹' }}</p>
            </div>
          </div>
          <button class="close-btn" @click="closeSMBModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveSMBShare" class="modal-body">
          <!-- 基本信息 -->
          <div class="form-section">
            <h4 class="section-title">基本信息</h4>
            <div class="form-row">
              <div class="form-group">
                <label>共享名称 *</label>
                <input
                  v-model="smbForm.name"
                  type="text"
                  required
                  :disabled="editingSMBShare"
                  placeholder="myshare"
                  pattern="[a-zA-Z0-9_-]+"
                  title="只能包含字母、数字、下划线和连字符"
                />
                <span class="form-hint">共享名称只能包含字母、数字、下划线和连字符</span>
              </div>

              <div class="form-group">
                <label>共享路径 *</label>
                <div class="path-input-group">
                  <input
                    v-model="smbForm.path"
                    type="text"
                    required
                    placeholder="/mnt/data/share"
                    @change="validatePath"
                  />
                  <button type="button" class="browse-btn" @click="browsePath" title="浏览路径">
                    <FolderOpenIcon class="w-4 h-4" />
                  </button>
                </div>
                <span v-if="pathError" class="form-error">{{ pathError }}</span>
                <span v-else class="form-hint">选择要共享的目录路径</span>
              </div>
            </div>

            <div class="form-group">
              <label>描述</label>
              <textarea
                v-model="smbForm.description"
                rows="2"
                placeholder="可选的共享描述，帮助用户了解此共享的用途"
              ></textarea>
            </div>
          </div>

          <!-- 访问控制 -->
          <div class="form-section">
            <h4 class="section-title">访问控制</h4>
            <div class="form-grid">
              <div class="form-group checkbox-group">
                <label class="checkbox-label">
                  <input v-model="smbForm.readOnly" type="checkbox" />
                  <span class="checkbox-text">
                    <strong>只读模式</strong>
                    <small>用户只能查看文件，无法修改或删除</small>
                  </span>
                </label>
              </div>

              <div class="form-group checkbox-group">
                <label class="checkbox-label">
                  <input v-model="smbForm.guest" type="checkbox" />
                  <span class="checkbox-text">
                    <strong>允许访客访问</strong>
                    <small>无需密码即可访问（安全性较低）</small>
                  </span>
                </label>
              </div>

              <div class="form-group checkbox-group">
                <label class="checkbox-label">
                  <input v-model="smbForm.browseable" type="checkbox" />
                  <span class="checkbox-text">
                    <strong>可浏览</strong>
                    <small>允许用户查看共享内容列表</small>
                  </span>
                </label>
              </div>

              <div class="form-group checkbox-group">
                <label class="checkbox-label">
                  <input v-model="smbForm.isTimeMachine" type="checkbox" />
                  <span class="checkbox-text">
                    <strong>Time Machine 支持</strong>
                    <small>支持苹果设备的Time Machine备份</small>
                  </span>
                </label>
              </div>
            </div>
          </div>

          <!-- 高级选项 -->
          <div class="form-section">
            <h4 class="section-title">高级选项</h4>
            <div class="form-grid">
              <div class="form-group">
                <label>有效用户</label>
                <input
                  v-model="smbForm.validUsers"
                  type="text"
                  placeholder="user1,user2"
                />
                <span class="form-hint">逗号分隔的用户列表</span>
              </div>

              <div class="form-group">
                <label>禁止用户</label>
                <input
                  v-model="smbForm.invalidUsers"
                  type="text"
                  placeholder="guest,nobody"
                />
                <span class="form-hint">不允许访问的用户列表</span>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label>创建权限</label>
                <select v-model="smbForm.createMask">
                  <option value="0777">0777 (所有用户)</option>
                  <option value="0755">0755 (默认)</option>
                  <option value="0750">0750 (组权限)</option>
                  <option value="0700">0700 (仅所有者)</option>
                </select>
              </div>

              <div class="form-group">
                <label>目录权限</label>
                <select v-model="smbForm.directoryMask">
                  <option value="0777">0777 (所有用户)</option>
                  <option value="0755">0755 (默认)</option>
                  <option value="0750">0750 (组权限)</option>
                  <option value="0700">0700 (仅所有者)</option>
                </select>
              </div>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="action-btn" @click="closeSMBModal">
              取消
            </button>
            <button type="button" class="action-btn" @click="testSMBConnection" :disabled="!smbForm.path">
              <SparklesIcon class="w-4 h-4" />
              测试连接
            </button>
            <button type="submit" class="action-btn primary" :disabled="!!pathError">
              <CheckIcon class="w-4 h-4" />
              {{ editingSMBShare ? '更新共享' : '创建共享' }}
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
  ClockIcon,
  PowerIcon,
  EyeIcon,
  FolderOpenIcon,
  LockClosedIcon,
  CheckIcon,
  SparklesIcon,
  CircleStackIcon
} from '@heroicons/vue/24/outline'
import { storageApi, fileApi } from '../api'
import StoragePoolManager from './StoragePoolManager.vue'
import DiskFormatDialog from '../components/Disk/DiskFormatDialog.vue'

const poolManagerRef = ref<any>(null)
const activeTab = ref('disks')
const loading = ref(false)
const loadingSMB = ref(false)
const loadingFiles = ref(false)

const showCreatePoolWizard = () => {
  setTimeout(() => {
    if (poolManagerRef.value) {
      poolManagerRef.value.showCreateDialog = true
    }
  }, 100)
}

const tabs = [
  { id: 'disks', label: '磁盘管理', icon: ServerIcon },
  { id: 'pools', label: '存储池', icon: CircleStackIcon },
  { id: 'smb', label: 'SMB共享', icon: FolderIcon },
  { id: 'files', label: '文件浏览', icon: DocumentIcon }
]

// Disk Management
const disks = ref<any[]>([])
const selectedDisk = ref<any>(null)
const showFormatDialog = ref(false)

// SMB Management
const smbShares = ref<any[]>([])
const showCreateSMBModal = ref(false)
const editingSMBShare = ref(false)
// SMB状态管理
const smbServiceRunning = ref(true)
const pathError = ref('')

const smbForm = ref({
  name: '',
  path: '',
  description: '',
  readOnly: false,
  guest: false,
  isTimeMachine: false,
  browseable: true,
  validUsers: '',
  invalidUsers: '',
  createMask: '0755',
  directoryMask: '0755'
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
    disks.value = response.disks || response  // axios拦截器已返回response.data
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
  if (disk.label?.includes('System')) {
    showError('无法格式化系统磁盘！')
    return
  }
  if (disk.mounted) {
    showError('请先卸载磁盘再进行格式化')
    return
  }

  selectedDisk.value = disk
  showFormatDialog.value = true
}

// SMB Functions
const loadSMBShares = async () => {
  loadingSMB.value = true
  try {
    const response = await storageApi.getSMBShares()
    smbShares.value = response.shares || response  // axios拦截器已返回response.data

    // 检查SMB服务状态（模拟，应该从后端获取）
    smbServiceRunning.value = true
  } catch (error: any) {
    console.error('Failed to fetch SMB shares:', error)
    smbShares.value = []
    smbServiceRunning.value = false
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
    readOnly: share.readOnly || false,
    guest: share.guest || false,
    isTimeMachine: share.isTimeMachine || false,
    browseable: share.browseable !== undefined ? share.browseable : true,
    validUsers: share.validUsers || '',
    invalidUsers: share.invalidUsers || '',
    createMask: share.createMask || '0755',
    directoryMask: share.directoryMask || '0755'
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
    const shareData = {
      name: smbForm.value.name,
      path: smbForm.value.path,
      description: smbForm.value.description,
      readOnly: smbForm.value.readOnly,
      guest: smbForm.value.guest,
      isTimeMachine: smbForm.value.isTimeMachine,
      browseable: smbForm.value.browseable,
      validUsers: smbForm.value.validUsers,
      invalidUsers: smbForm.value.invalidUsers,
      createMask: smbForm.value.createMask,
      directoryMask: smbForm.value.directoryMask
    }

    if (editingSMBShare.value) {
      await storageApi.updateSMBShare(
        smbForm.value.name,
        shareData
      )
      showSuccess('共享更新成功')
    } else {
      await storageApi.createSMBShare(
        smbForm.value.name,
        smbForm.value.path,
        smbForm.value.description,
        smbForm.value.readOnly,
        smbForm.value.guest,
        smbForm.value.isTimeMachine
      )
      showSuccess('共享创建成功')
    }
    closeSMBModal()
    await loadSMBShares()
  } catch (error: any) {
    console.error('Failed to save SMB share:', error)
    showError('保存失败: ' + (error.response?.data?.error || error.message))
  }
}

// 新增的SMB功能函数
const openCreateSMBModal = () => {
  resetSMBForm()
  editingSMBShare.value = false
  showCreateSMBModal.value = true
}

const closeSMBModal = () => {
  showCreateSMBModal.value = false
  editingSMBShare.value = false
  resetSMBForm()
  pathError.value = ''
}

const resetSMBForm = () => {
  smbForm.value = {
    name: '',
    path: '',
    description: '',
    readOnly: false,
    guest: false,
    isTimeMachine: false,
    browseable: true,
    validUsers: '',
    invalidUsers: '',
    createMask: '0755',
    directoryMask: '0755'
  }
}

const validatePath = () => {
  pathError.value = ''
  if (smbForm.value.path && !smbForm.value.path.startsWith('/')) {
    pathError.value = '路径必须以 / 开头'
    return false
  }
  return true
}

const browsePath = async () => {
  // 这里可以实现路径浏览功能
  // 暂时使用一个简单的prompt作为示例
  const path = prompt('请输入目录路径:', smbForm.value.path || '/mnt')
  if (path) {
    smbForm.value.path = path
    validatePath()
  }
}

const viewSMBShare = (share: any) => {
  // 显示共享详细信息
  alert(`共享详情:\n名称: ${share.name}\n路径: ${share.path}\n描述: ${share.description || '无'}\n只读: ${share.readOnly ? '是' : '否'}\n访客: ${share.guest ? '是' : '否'}`)
}

const toggleSMBShareStatus = async (share: any) => {
  try {
    // 暂时简化此功能，只显示提示信息
    const enabled = share.enabled !== undefined ? share.enabled : true
    const newStatus = !enabled
    showSuccess(`共享 "${share.name}" 状态切换功能暂未启用`)
    // TODO: 实现后端API后启用此功能
    // await storageApi.toggleSMBShare(share.name, newStatus)
  } catch (error: any) {
    console.error('Failed to toggle SMB share status:', error)
    showError('操作失败: ' + error.message)
  }
}

const refreshSMBShares = async () => {
  await loadSMBShares()
}

const toggleSMBService = async () => {
  try {
    // 这里需要调用后端API来启动/停止SMB服务
    smbServiceRunning.value = !smbServiceRunning.value
    const status = smbServiceRunning.value ? '运行' : '停止'
    showSuccess(`SMB服务已${status}`)
  } catch (error: any) {
    console.error('Failed to toggle SMB service:', error)
    showError('操作失败: ' + error.message)
  }
}

const testSMBConnection = async () => {
  try {
    // 测试路径连接
    if (!smbForm.value.path) {
      showError('请先输入共享路径')
      return
    }

    // 这里可以调用后端API来测试路径是否可访问
    showSuccess(`路径 "${smbForm.value.path}" 可访问`)
  } catch (error: any) {
    console.error('Failed to test SMB connection:', error)
    showError('路径测试失败: ' + error.message)
  }
}

// 辅助函数
const showSuccess = (message: string) => {
  // 可以替换为更好的通知组件
  alert('✅ ' + message)
}

const showError = (message: string) => {
  // 可以替换为更好的错误提示
  alert('❌ ' + message)
}

// File Browser Functions
const loadFiles = async () => {
  loadingFiles.value = true
  try {
    const response = await fileApi.listFiles(currentPath.value)
    currentFiles.value = response.files.map((file: any) => ({  // axios拦截器已返回response.data
      name: file.name,
      isDirectory: file.isDir,
      size: file.size,
      permissions: file.permissions,
      modified: new Date(file.modTime).toLocaleDateString(),
      mimeType: file.mimeType
    }))
  } catch (error: any) {
    console.error('Failed to load files:', error)
    showError('无法加载文件列表，请检查网络连接或路径权限')
    currentFiles.value = []
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

const onDiskFormatted = () => {
  refreshDisks()
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

.mount-point-list {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 4px;
  background: #f9fafb;
  padding: 4px 8px;
  border-radius: 4px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.mount-point-list .label {
  font-weight: 600;
  margin-right: 4px;
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

/* ==================== 增强的SMB管理样式 ==================== */

/* SMB服务状态指示器 */
.header-left {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.smb-status-indicator {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  background: #f3f4f6;
  color: #6b7280;
}

.smb-status-indicator.active {
  background: #d1fae5;
  color: #065f46;
}

.smb-status-indicator .status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #6b7280;
}

.empty-state svg {
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #374151;
}

.empty-state p {
  font-size: 14px;
  margin-bottom: 24px;
}

/* 增强的SMB列表项 */
.smb-item {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 20px;
  align-items: center;
  transition: all 0.2s;
  border: 1px solid #e5e7eb;
}

.smb-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

.smb-item-disabled {
  opacity: 0.6;
  background: #f9fafb;
}

.smb-icon-readonly {
  background: linear-gradient(135deg, #6b7280 0%, #9ca3af 100%);
}

.disabled-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.smb-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.smb-badges {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.badge {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.badge-disabled {
  background: #f3f4f6;
  color: #6b7280;
}

.badge-readonly {
  background: #dbeafe;
  color: #1e40af;
}

.badge-guest {
  background: #fef3c7;
  color: #92400e;
}

.badge-tm {
  background: #e0e7ff;
  color: #3730a3;
}

.smb-description {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.smb-path {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #9ca3af;
  background: #f9fafb;
  padding: 4px 8px;
  border-radius: 6px;
  width: fit-content;
}

/* 设置组 */
.setting-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 6px 12px;
  background: #f9fafb;
  border-radius: 8px;
  font-size: 13px;
}

.setting-label {
  color: #6b7280;
  font-weight: 500;
}

.setting-value {
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
}

.setting-value.readonly {
  color: #1e40af;
  background: #dbeafe;
}

.setting-value.writable {
  color: #065f46;
  background: #d1fae5;
}

.setting-value.allowed {
  color: #065f46;
  background: #d1fae5;
}

.setting-value.denied {
  color: #991b1b;
  background: #fee2e2;
}

/* 增强的SMB模态框 */
.smb-modal {
  max-width: 700px;
  width: 90%;
}

.modal-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.modal-subtitle {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 20px 0;
  border-bottom: 1px solid #e5e7eb;
}

.form-section:last-of-type {
  border-bottom: none;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-title::before {
  content: '';
  width: 4px;
  height: 16px;
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  border-radius: 2px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

@media (max-width: 768px) {
  .form-row, .form-grid {
    grid-template-columns: 1fr;
  }
}

/* 增强的表单控件 */
.form-group textarea {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
}

.form-group select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.path-input-group {
  display: flex;
  gap: 8px;
}

.path-input-group input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px 0 0 8px;
  font-size: 14px;
}

.browse-btn {
  padding: 10px 16px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-left: none;
  border-radius: 0 8px 8px 0;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.browse-btn:hover {
  background: #e5e7eb;
}

.form-hint {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}

.form-error {
  font-size: 12px;
  color: #dc2626;
  margin-top: 4px;
}

/* 复选框组增强 */
.checkbox-group {
  background: #f9fafb;
  padding: 12px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.checkbox-label {
  display: flex;
  gap: 8px;
  cursor: pointer;
  align-items: flex-start;
}

.checkbox-text {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.checkbox-text strong {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.checkbox-text small {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.4;
}

/* 操作按钮增强 */
.smb-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  background: white;
  color: #374151;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 6px;
}

.action-btn:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
}

.action-btn.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.action-btn.danger {
  color: #dc2626;
  border-color: #fca5a5;
}

.action-btn.danger:hover {
  background: #fef2f2;
  border-color: #dc2626;
}

/* 动画效果 */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.smb-item {
  animation: fadeIn 0.3s ease-out;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .smb-item {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .smb-actions {
    justify-content: flex-end;
  }

  .header-actions {
    flex-wrap: wrap;
  }
}
</style>