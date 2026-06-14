<template>
  <div class="permission-manager">
    <!-- 头部 -->
    <div class="pm-header">
      <div class="header-left">
        <h1>权限管理</h1>
        <p class="subtitle">管理用户权限、共享文件夹和访问控制</p>
      </div>

      <div class="header-right">
        <button class="action-btn primary" @click="createShareFolder">
          <PlusIcon class="w-4 h-4" />
          创建共享文件夹
        </button>
        <button class="action-btn" @click="showAdvancedSettings = true">
          <CogIcon class="w-4 h-4" />
          高级设置
        </button>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="pm-content">
      <!-- 左侧：共享文件夹列表 -->
      <div class="pm-folders">
        <div class="section-header">
          <h2>共享文件夹</h2>
          <div class="filter-controls">
            <select v-model="selectedFilter" class="filter-select">
              <option value="all">所有共享</option>
              <option value="home">主目录</option>
              <option value="shared">共享文件夹</option>
              <option value="external">外接设备</option>
            </select>
          </div>
        </div>

        <div class="folder-list">
          <div
            v-for="folder in filteredFolders"
            :key="folder.id"
            class="folder-item"
            :class="{ selected: selectedFolder === folder.id }"
            @click="selectFolder(folder)"
          >
            <div class="folder-icon">
              <FolderIcon class="w-6 h-6" />
            </div>

            <div class="folder-info">
              <div class="folder-name">{{ folder.name }}</div>
              <div class="folder-path">{{ folder.path }}</div>
              <div class="folder-stats">
                <span>{{ folder.fileCount }} 个文件</span>
                <span>{{ formatBytes(folder.size) }}</span>
              </div>
            </div>

            <div class="folder-status">
              <div class="status-badge" :class="folder.status">
                {{ getStatusText(folder.status) }}
              </div>
            </div>
          </div>
        </div>

        <!-- 快速操作 -->
        <div class="folder-actions">
          <button class="action-btn" @click="refreshFolderList">
            <ArrowPathIcon class="w-4 h-4" />
            刷新
          </button>
          <button class="action-btn" @click="showFolderUsage">
            <ChartPieIcon class="w-4 h-4" />
            使用情况
          </button>
        </div>
      </div>

      <!-- 右侧：权限详情 -->
      <div class="pm-permissions">
        <div v-if="selectedFolderData" class="permissions-content">
          <!-- 基本信息 -->
          <div class="section-card">
            <div class="card-header">
              <h3>{{ selectedFolderData.name }}</h3>
              <div class="card-actions">
                <button class="icon-btn" @click="editFolder">
                  <PencilIcon class="w-4 h-4" />
                </button>
                <button class="icon-btn" @click="deleteFolder">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
            <div class="folder-details">
              <div class="detail-row">
                <span class="detail-label">路径:</span>
                <span class="detail-value">{{ selectedFolderData.path }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">描述:</span>
                <span class="detail-value">{{ selectedFolderData.description }}</span>
              </div>
              <div class="detail-row">
                <span class="detail-label">状态:</span>
                <span class="detail-value" :class="selectedFolderData.status">
                  {{ getStatusText(selectedFolderData.status) }}
                </span>
              </div>
            </div>
          </div>

          <!-- 权限列表 -->
          <div class="section-card">
            <div class="card-header">
              <h3>访问权限</h3>
              <button class="action-btn small" @click="addPermission">
                <UserPlusIcon class="w-4 h-4" />
                添加用户
              </button>
            </div>

            <div class="permissions-list">
              <div
                v-for="permission in selectedFolderData.permissions"
                :key="permission.userId"
                class="permission-item"
              >
                <div class="user-info">
                  <div class="user-avatar">
                    <UserIcon class="w-5 h-5" />
                  </div>
                  <div>
                    <div class="user-name">{{ permission.userName }}</div>
                    <div class="user-role">{{ permission.userRole }}</div>
                  </div>
                </div>

                <div class="permission-controls">
                  <div class="permission-group">
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="permission.read"
                        @change="updatePermission(permission)"
                      />
                      <span>读取</span>
                    </label>
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="permission.write"
                        @change="updatePermission(permission)"
                      />
                      <span>写入</span>
                    </label>
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="permission.execute"
                        @change="updatePermission(permission)"
                      />
                      <span>执行</span>
                    </label>
                    <label class="checkbox-label">
                      <input
                        type="checkbox"
                        v-model="permission.admin"
                        @change="updatePermission(permission)"
                      />
                      <span>管理</span>
                    </label>
                  </div>
                </div>

                <div class="permission-actions">
                  <button
                    class="icon-btn"
                    @click="editPermission(permission)"
                    title="编辑"
                  >
                    <PencilIcon class="w-4 h-4" />
                  </button>
                  <button
                    class="icon-btn danger"
                    @click="removePermission(permission)"
                    title="移除"
                  >
                    <XMarkIcon class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Windows ACL 设置 -->
          <div class="section-card">
            <div class="card-header">
              <h3>Windows ACL 设置</h3>
              <label class="toggle-switch">
                <input type="checkbox" v-model="selectedFolderData.aclEnabled" />
                <span class="toggle-slider"></span>
              </label>
            </div>

            <div v-if="selectedFolderData.aclEnabled" class="acl-settings">
              <div class="acl-info">
                <InformationCircleIcon class="w-4 h-4" />
                <p>启用 Windows ACL 可以提供更精细的权限控制，但会降低系统性能。</p>
              </div>

              <div class="acl-permissions">
                <div class="acl-item">
                  <div class="acl-user">Everyone</div>
                  <div class="acl-rights">读取 & 执行</div>
                  <div class="acl-actions">
                    <button class="icon-btn" @click="editAclPermission">
                      <PencilIcon class="w-4 h-4" />
                    </button>
                  </div>
                </div>
                <div class="acl-item">
                  <div class="acl-user">Administrators</div>
                  <div class="acl-rights">完全控制</div>
                  <div class="acl-actions">
                    <button class="icon-btn" @click="editAclPermission">
                      <PencilIcon class="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 高级设置 -->
          <div class="section-card">
            <div class="card-header">
              <h3>高级设置</h3>
            </div>

            <div class="advanced-settings">
              <div class="setting-item">
                <div class="setting-info">
                  <div class="setting-title">启用回收站</div>
                  <div class="setting-description">删除的文件将移至回收站而非永久删除</div>
                </div>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="selectedFolderData.recycleBinEnabled" />
                  <span class="toggle-slider"></span>
                </label>
              </div>

              <div class="setting-item">
                <div class="setting-info">
                  <div class="setting-title">隐藏文件</div>
                  <div class="setting-description">隐藏以点开头的系统文件</div>
                </div>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="selectedFolderData.hideHiddenFiles" />
                  <span class="toggle-slider"></span>
                </label>
              </div>

              <div class="setting-item">
                <div class="setting-info">
                  <div class="setting-title">只读访问</div>
                  <div class="setting-description">所有用户只能读取，不能修改</div>
                </div>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="selectedFolderData.readOnly" />
                  <span class="toggle-slider"></span>
                </label>
              </div>

              <div class="setting-item">
                <div class="setting-info">
                  <div class="setting-title">禁止浏览</div>
                  <div class="setting-description">用户需要知道完整路径才能访问文件</div>
                </div>
                <label class="toggle-switch">
                  <input type="checkbox" v-model="selectedFolderData.noBrowse" />
                  <span class="toggle-slider"></span>
                </label>
              </div>
            </div>
          </div>

          <!-- 权限继承 -->
          <div class="section-card">
            <div class="card-header">
              <h3>权限继承</h3>
            </div>

            <div class="inheritance-settings">
              <div class="inheritance-info">
                <InformationCircleIcon class="w-4 h-4" />
                <p>子文件夹和文件可以继承父文件夹的权限设置。</p>
              </div>

              <div class="inheritance-options">
                <label class="radio-label">
                  <input
                    type="radio"
                    v-model="selectedFolderData.inheritMode"
                    value="parent"
                  />
                  <span>继承父文件夹权限</span>
                </label>
                <label class="radio-label">
                  <input
                    type="radio"
                    v-model="selectedFolderData.inheritMode"
                    value="custom"
                  />
                  <span>自定义权限</span>
                </label>
              </div>

              <button class="action-btn" @click="applyToChildren">
                <ArrowDownIcon class="w-4 h-4" />
                应用到所有子项目
              </button>
            </div>
          </div>
        </div>

        <div v-else class="no-selection">
          <FolderIcon class="w-12 h-12" />
          <p>选择一个共享文件夹查看权限</p>
        </div>
      </div>
    </div>

    <!-- 添加权限模态框 -->
    <Transition name="fade">
      <div v-if="showAddPermission" class="modal-overlay" @click.self="showAddPermission = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>添加用户权限</h3>
            <button @click="showAddPermission = false">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="modal-body">
            <div class="form-section">
              <div class="form-group">
                <label>选择用户</label>
                <select v-model="newPermission.userId" class="form-select">
                  <option value="">选择用户...</option>
                  <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                    {{ user.name }} ({{ user.role }})
                  </option>
                </select>
              </div>

              <div class="form-group">
                <label>权限设置</label>
                <div class="permission-checkboxes">
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="newPermission.read" />
                    <span>读取</span>
                  </label>
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="newPermission.write" />
                    <span>写入</span>
                  </label>
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="newPermission.execute" />
                    <span>执行</span>
                  </label>
                  <label class="checkbox-label">
                    <input type="checkbox" v-model="newPermission.admin" />
                    <span>管理</span>
                  </label>
                </div>
              </div>
            </div>

            <div class="modal-actions">
              <button class="action-btn" @click="showAddPermission = false">取消</button>
              <button class="action-btn primary" @click="confirmAddPermission">添加</button>
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
  CogIcon,
  FolderIcon,
  ArrowPathIcon,
  ChartPieIcon,
  PencilIcon,
  TrashIcon,
  UserPlusIcon,
  UserIcon,
  XMarkIcon,
  InformationCircleIcon,
  ArrowDownIcon
} from '@heroicons/vue/24/outline'

interface Folder {
  id: string
  name: string
  path: string
  description: string
  status: 'active' | 'disabled' | 'error'
  fileCount: number
  size: number
  permissions: Permission[]
  aclEnabled: boolean
  recycleBinEnabled: boolean
  hideHiddenFiles: boolean
  readOnly: boolean
  noBrowse: boolean
  inheritMode: 'parent' | 'custom'
}

interface Permission {
  userId: string
  userName: string
  userRole: string
  read: boolean
  write: boolean
  execute: boolean
  admin: boolean
}

interface User {
  id: string
  name: string
  role: string
}

// 状态管理
const selectedFilter = ref('all')
const selectedFolder = ref<string | null>(null)
const showAddPermission = ref(false)
const showAdvancedSettings = ref(false)

const newPermission = ref({
  userId: '',
  read: true,
  write: false,
  execute: false,
  admin: false
})

// 模拟数据
const folders = ref<Folder[]>([
  {
    id: '1',
    name: 'Documents',
    path: '/home/Documents',
    description: '个人文档文件夹',
    status: 'active',
    fileCount: 1250,
    size: 5368709120, // 5GB
    permissions: [
      {
        userId: 'u1',
        userName: 'admin',
        userRole: '管理员',
        read: true,
        write: true,
        execute: true,
        admin: true
      },
      {
        userId: 'u2',
        userName: 'john',
        userRole: '普通用户',
        read: true,
        write: true,
        execute: true,
        admin: false
      },
      {
        userId: 'u3',
        userName: 'jane',
        userRole: '普通用户',
        read: true,
        write: false,
        execute: true,
        admin: false
      }
    ],
    aclEnabled: false,
    recycleBinEnabled: true,
    hideHiddenFiles: true,
    readOnly: false,
    noBrowse: false,
    inheritMode: 'custom'
  },
  {
    id: '2',
    name: 'Photos',
    path: '/home/Photos',
    description: '照片和图片文件',
    status: 'active',
    fileCount: 3200,
    size: 21474836480, // 20GB
    permissions: [
      {
        userId: 'u1',
        userName: 'admin',
        userRole: '管理员',
        read: true,
        write: true,
        execute: true,
        admin: true
      },
      {
        userId: 'u2',
        userName: 'john',
        userRole: '普通用户',
        read: true,
        write: false,
        execute: true,
        admin: false
      }
    ],
    aclEnabled: false,
    recycleBinEnabled: true,
    hideHiddenFiles: true,
    readOnly: false,
    noBrowse: false,
    inheritMode: 'parent'
  },
  {
    id: '3',
    name: 'Public',
    path: '/home/Public',
    description: '公共共享文件夹',
    status: 'active',
    fileCount: 85,
    size: 524288000, // 500MB
    permissions: [
      {
        userId: 'u1',
        userName: 'admin',
        userRole: '管理员',
        read: true,
        write: true,
        execute: true,
        admin: true
      },
      {
        userId: 'everyone',
        userName: 'Everyone',
        userRole: '所有人',
        read: true,
        write: true,
        execute: true,
        admin: false
      }
    ],
    aclEnabled: false,
    recycleBinEnabled: true,
    hideHiddenFiles: false,
    readOnly: false,
    noBrowse: false,
    inheritMode: 'parent'
  }
])

const availableUsers = ref<User[]>([
  { id: 'u1', name: 'admin', role: '管理员' },
  { id: 'u2', name: 'john', role: '普通用户' },
  { id: 'u3', name: 'jane', role: '普通用户' },
  { id: 'u4', name: 'guest', role: '访客' }
])

// 计算属性
const filteredFolders = computed(() => {
  if (selectedFilter.value === 'all') return folders.value
  return folders.value.filter(folder => {
    if (selectedFilter.value === 'home') return folder.path.startsWith('/home')
    if (selectedFilter.value === 'shared') return folder.status === 'active'
    if (selectedFilter.value === 'external') return folder.path.startsWith('/external')
    return true
  })
})

const selectedFolderData = computed(() => {
  return folders.value.find(f => f.id === selectedFolder.value) || null
})

// 方法
const selectFolder = (folder: Folder) => {
  selectedFolder.value = folder.id
}

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    active: '已启用',
    disabled: '已禁用',
    error: '错误'
  }
  return statusMap[status] || status
}

const formatBytes = (bytes: number): string => {
  const gb = bytes / 1073741824
  if (gb >= 1) return gb.toFixed(2) + ' GB'
  const mb = bytes / 1048576
  if (mb >= 1) return mb.toFixed(2) + ' MB'
  const kb = bytes / 1024
  return kb.toFixed(2) + ' KB'
}

const refreshFolderList = () => {
  console.log('刷新文件夹列表')
}

const showFolderUsage = () => {
  console.log('显示使用情况')
}

const createShareFolder = () => {
  console.log('创建共享文件夹')
}

const editFolder = () => {
  console.log('编辑文件夹')
}

const deleteFolder = () => {
  if (selectedFolderData.value && confirm(`确定要删除 "${selectedFolderData.value.name}" 吗？`)) {
    const index = folders.value.findIndex(f => f.id === selectedFolder.value)
    if (index > -1) {
      folders.value.splice(index, 1)
      selectedFolder.value = null
    }
  }
}

const addPermission = () => {
  showAddPermission.value = true
}

const updatePermission = (permission: Permission) => {
  console.log('更新权限:', permission)
}

const editPermission = (permission: Permission) => {
  console.log('编辑权限:', permission)
}

const removePermission = (permission: Permission) => {
  if (selectedFolderData.value && confirm(`确定要移除 "${permission.userName}" 的权限吗？`)) {
    const index = selectedFolderData.value.permissions.findIndex(p => p.userId === permission.userId)
    if (index > -1) {
      selectedFolderData.value.permissions.splice(index, 1)
    }
  }
}

const editAclPermission = () => {
  console.log('编辑 ACL 权限')
}

const applyToChildren = () => {
  console.log('应用到所有子项目')
}

const confirmAddPermission = () => {
  if (!selectedFolderData.value || !newPermission.value.userId) return

  const user = availableUsers.value.find(u => u.id === newPermission.value.userId)
  if (!user) return

  selectedFolderData.value.permissions.push({
    userId: user.id,
    userName: user.name,
    userRole: user.role,
    read: newPermission.value.read,
    write: newPermission.value.write,
    execute: newPermission.value.execute,
    admin: newPermission.value.admin
  })

  // 重置表单
  newPermission.value = {
    userId: '',
    read: true,
    write: false,
    execute: false,
    admin: false
  }

  showAddPermission.value = false
}
</script>

<style scoped>
.permission-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f9fafb;
}

.pm-header {
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

.action-btn.small {
  padding: 6px 12px;
  font-size: 12px;
}

.pm-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  flex: 1;
  overflow: hidden;
  padding: 20px;
}

.pm-folders {
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

.folder-list {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.folder-item {
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

.folder-item:hover {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.folder-item.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.folder-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: #eff6ff;
  color: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
}

.folder-info {
  flex: 1;
}

.folder-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.folder-path {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
}

.folder-stats {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: #9ca3af;
}

.folder-status {
  margin-left: auto;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 6px;
  font-size: 12px;
  background: #f3f4f6;
  color: #6b7280;
}

.status-badge.active {
  background: #ecfdf5;
  color: #10b981;
}

.status-badge.disabled {
  background: #f3f4f6;
  color: #9ca3af;
}

.status-badge.error {
  background: #fef2f2;
  color: #ef4444;
}

.folder-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
}

.pm-permissions {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.permissions-content {
  padding: 20px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.card-actions {
  display: flex;
  gap: 4px;
}

.icon-btn {
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

.icon-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.icon-btn.danger:hover {
  background: #fef2f2;
  color: #ef4444;
}

.folder-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
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

.detail-value.active {
  color: #10b981;
}

.permissions-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: #eff6ff;
  color: #3b82f6;
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 2px;
}

.user-role {
  font-size: 11px;
  color: #6b7280;
}

.permission-controls {
  flex: 2;
}

.permission-group {
  display: flex;
  gap: 16px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  font-size: 13px;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
}

.permission-actions {
  display: flex;
  gap: 4px;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #d1d5db;
  transition: .4s;
  border-radius: 24px;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .toggle-slider {
  background-color: #3b82f6;
}

input:checked + .toggle-slider:before {
  transform: translateX(20px);
}

.acl-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.acl-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #eff6ff;
  border-radius: 8px;
  font-size: 13px;
  color: #3b82f6;
}

.acl-permissions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.acl-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.acl-user {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.acl-rights {
  flex: 1;
  font-size: 13px;
  color: #6b7280;
}

.acl-actions {
  display: flex;
  gap: 4px;
}

.advanced-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.setting-info {
  flex: 1;
}

.setting-title {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 2px;
}

.setting-description {
  font-size: 12px;
  color: #6b7280;
}

.inheritance-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.inheritance-info {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #eff6ff;
  border-radius: 8px;
  font-size: 13px;
  color: #3b82f6;
}

.inheritance-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

.radio-label input[type="radio"] {
  width: 16px;
  height: 16px;
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
  max-width: 500px;
  background: white;
  border-radius: 16px;
  overflow: hidden;
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
}

.form-section {
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
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
}

.form-select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
  cursor: pointer;
}

.permission-checkboxes {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 16px;
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