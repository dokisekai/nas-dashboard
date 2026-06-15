<template>
  <div class="share-folder-manager">
    <!-- 头部 -->
    <div class="sfm-header">
      <h1>共享文件夹</h1>
      <button class="sfm-btn primary" @click="showCreateFolder = true">
        <PlusIcon class="w-4 h-4" />
        新增
      </button>
    </div>

    <!-- 共享文件夹列表 -->
    <div class="sfm-folders">
      <div
        v-for="folder in folders"
        :key="folder.id"
        class="sfm-folder-item"
        :class="{ expanded: expandedFolders.includes(folder.id) }"
      >
        <div class="sfm-folder-header" @click="toggleFolder(folder.id)">
          <div class="folder-basic">
            <ChevronRightIcon class="expand-icon" />
            <FolderIcon class="folder-icon" />
            <div class="folder-info">
              <div class="folder-name">{{ folder.name }}</div>
              <div class="folder-path">{{ folder.path }}</div>
            </div>
          </div>
          <div class="folder-protocols">
            <div
              v-for="protocol in getEnabledProtocols(folder)"
              :key="protocol"
              class="protocol-badge"
              :class="protocol"
            >
              {{ getProtocolLabel(protocol) }}
            </div>
          </div>
          <div class="folder-actions">
            <button class="sfm-link" @click.stop="editFolder(folder)">编辑</button>
            <button class="sfm-link danger" @click.stop="deleteFolder(folder.id)">删除</button>
          </div>
        </div>

        <div v-if="expandedFolders.includes(folder.id)" class="sfm-folder-details">
          <!-- 基本信息 -->
          <div class="detail-section">
            <h3>基本信息</h3>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">存储池:</span>
                <span class="info-value">{{ folder.storagePool }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">路径:</span>
                <span class="info-value">{{ folder.path }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">描述:</span>
                <span class="info-value">{{ folder.description || '无' }}</span>
              </div>
            </div>
          </div>

          <!-- 协议设置 -->
          <div class="detail-section">
            <h3>协议设置</h3>
            <div class="protocol-list">
              <div
                v-for="protocol in availableProtocols"
                :key="protocol.id"
                class="protocol-item"
              >
                <div class="protocol-header">
                  <div class="protocol-info">
                    <input
                      type="checkbox"
                      :checked="isProtocolEnabled(folder, protocol.id)"
                      @change="toggleProtocol(folder.id, protocol.id, $event)"
                    />
                    <span class="protocol-name">{{ protocol.name }}</span>
                    <span class="protocol-desc">{{ protocol.description }}</span>
                  </div>
                </div>

                <!-- 协议特定设置 -->
                <div v-if="isProtocolEnabled(folder, protocol.id)" class="protocol-settings">
                  <div v-if="protocol.id === 'smb'" class="settings-grid">
                    <div class="setting-item">
                      <label>共享名称</label>
                      <input
                        :value="getProtocolSetting(folder, 'smb', 'shareName')"
                        @input="updateProtocolSetting(folder.id, 'smb', 'shareName', $event)"
                        placeholder="默认为文件夹名称"
                      />
                    </div>
                    <div class="setting-item">
                      <label>启用回收站</label>
                      <input
                        type="checkbox"
                        :checked="getProtocolSetting(folder, 'smb', 'recycleBin')"
                        @change="updateProtocolSetting(folder.id, 'smb', 'recycleBin', $event)"
                      />
                    </div>
                  </div>

                  <div v-if="protocol.id === 'ftp'" class="settings-grid">
                    <div class="setting-item">
                      <label>FTP端口</label>
                      <input
                        type="number"
                        :value="getProtocolSetting(folder, 'ftp', 'port')"
                        @input="updateProtocolSetting(folder.id, 'ftp', 'port', $event)"
                      />
                    </div>
                    <div class="setting-item">
                      <label>被动模式</label>
                      <input
                        type="checkbox"
                        :checked="getProtocolSetting(folder, 'ftp', 'passiveMode')"
                        @change="updateProtocolSetting(folder.id, 'ftp', 'passiveMode', $event)"
                      />
                    </div>
                  </div>

                  <div v-if="protocol.id === 'sftp'" class="settings-grid">
                    <div class="setting-item">
                      <label>SFTP端口</label>
                      <input
                        type="number"
                        :value="getProtocolSetting(folder, 'sftp', 'port')"
                        @input="updateProtocolSetting(folder.id, 'sftp', 'port', $event)"
                      />
                    </div>
                    <div class="setting-item">
                      <label>允许Shell访问</label>
                      <input
                        type="checkbox"
                        :checked="getProtocolSetting(folder, 'sftp', 'allowShell')"
                        @change="updateProtocolSetting(folder.id, 'sftp', 'allowShell', $event)"
                      />
                    </div>
                  </div>

                  <div v-if="protocol.id === 'nfs'" class="settings-grid">
                    <div class="setting-item">
                      <label>NFS版本</label>
                      <select
                        :value="getProtocolSetting(folder, 'nfs', 'version')"
                        @change="updateProtocolSetting(folder.id, 'nfs', 'version', $event)"
                      >
                        <option value="4">NFSv4</option>
                        <option value="3">NFSv3</option>
                      </select>
                    </div>
                    <div class="setting-item">
                      <label>允许访问</label>
                      <input
                        :value="getProtocolSetting(folder, 'nfs', 'allowedClients')"
                        @input="updateProtocolSetting(folder.id, 'nfs', 'allowedClients', $event)"
                        placeholder="例如: 192.168.1.0/24"
                      />
                    </div>
                  </div>

                  <div v-if="protocol.id === 'webdav'" class="settings-grid">
                    <div class="setting-item">
                      <label>HTTPS</label>
                      <input
                        type="checkbox"
                        :checked="getProtocolSetting(folder, 'webdav', 'https')"
                        @change="updateProtocolSetting(folder.id, 'webdav', 'https', $event)"
                      />
                    </div>
                    <div class="setting-item">
                      <label>端口</label>
                      <input
                        type="number"
                        :value="getProtocolSetting(folder, 'webdav', 'port')"
                        @input="updateProtocolSetting(folder.id, 'webdav', 'port', $event)"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 用户权限 -->
          <div class="detail-section">
            <div class="section-header">
              <h3>用户权限</h3>
              <button class="sfm-btn small" @click.stop="addPermission(folder)">
                <PlusIcon class="w-3 h-3" />
                添加用户
              </button>
            </div>

            <div class="permissions-list">
              <div
                v-for="permission in folder.permissions"
                :key="permission.userId"
                class="permission-item"
              >
                <div class="user-info">
                  <UserIcon class="w-4 h-4" />
                  <span>{{ permission.userName }}</span>
                </div>
                <div class="permission-level">
                  <select
                    :value="getPermissionLevel(permission)"
                    @change="updatePermission(folder.id, permission.userId, $event)"
                  >
                    <option value="none">无权限</option>
                    <option value="read">只读</option>
                    <option value="write">读写</option>
                    <option value="admin">管理</option>
                  </select>
                </div>
                <button
                  class="sfm-link danger"
                  @click="removePermission(folder.id, permission.userId)"
                >
                  移除
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建共享文件夹对话框 -->
    <div v-if="showCreateFolder" class="sfm-modal" @click.self="closeCreateModal">
      <div class="sfm-modal-content">
        <div class="sfm-modal-header">
          <h3>{{ editingFolder ? '编辑共享文件夹' : '新增共享文件夹' }}</h3>
          <button class="sfm-close" @click="closeCreateModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveFolder" class="sfm-form">
          <div class="sfm-form-group">
            <label>名称 *</label>
            <input
              v-model="folderForm.name"
              type="text"
              required
              :disabled="editingFolder"
              placeholder="文件夹名称"
            />
          </div>

          <div class="sfm-form-group">
            <label>存储池 *</label>
            <select v-model="folderForm.storagePool" required>
              <option value="">选择存储池...</option>
              <option v-for="pool in storagePools" :key="pool.id" :value="pool.name">
                {{ pool.name }} ({{ pool.availableSpace }} 可用)
              </option>
            </select>
          </div>

          <div class="sfm-form-group">
            <label>描述</label>
            <textarea
              v-model="folderForm.description"
              rows="3"
              placeholder="文件夹描述"
            ></textarea>
          </div>

          <div class="sfm-modal-footer">
            <button type="button" class="sfm-btn secondary" @click="closeCreateModal">
              取消
            </button>
            <button type="submit" class="sfm-btn primary">
              {{ editingFolder ? '确定' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 添加用户权限对话框 -->
    <div v-if="showAddPermission" class="sfm-modal" @click.self="closePermissionModal">
      <div class="sfm-modal-content">
        <div class="sfm-modal-header">
          <h3>添加用户权限</h3>
          <button class="sfm-close" @click="closePermissionModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="confirmAddPermission" class="sfm-form">
          <div class="sfm-form-group">
            <label>选择用户</label>
            <select v-model="permissionForm.userId" required>
              <option value="">选择用户...</option>
              <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                {{ user.name }}
              </option>
            </select>
          </div>

          <div class="sfm-form-group">
            <label>权限级别</label>
            <div class="permission-options">
              <label class="radio-option">
                <input type="radio" v-model="permissionForm.level" value="read" />
                <div class="radio-content">
                  <div class="radio-title">只读</div>
                  <div class="radio-desc">可以浏览和下载文件</div>
                </div>
              </label>
              <label class="radio-option">
                <input type="radio" v-model="permissionForm.level" value="write" />
                <div class="radio-content">
                  <div class="radio-title">读写</div>
                  <div class="radio-desc">可以读取、上传和修改文件</div>
                </div>
              </label>
              <label class="radio-option">
                <input type="radio" v-model="permissionForm.level" value="admin" />
                <div class="radio-content">
                  <div class="radio-title">管理</div>
                  <div class="radio-desc">完全控制权限</div>
                </div>
              </label>
            </div>
          </div>

          <div class="sfm-modal-footer">
            <button type="button" class="sfm-btn secondary" @click="closePermissionModal">
              取消
            </button>
            <button type="submit" class="sfm-btn primary">添加</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  PlusIcon,
  ChevronRightIcon,
  FolderIcon,
  UserIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

// 可用协议
const availableProtocols = [
  { id: 'smb', name: 'SMB', description: 'Windows文件共享' },
  { id: 'ftp', name: 'FTP', description: '文件传输协议' },
  { id: 'sftp', name: 'SFTP', description: 'SSH文件传输' },
  { id: 'nfs', name: 'NFS', description: '网络文件系统' },
  { id: 'webdav', name: 'WebDAV', description: 'Web分布式创作' }
]

// 模拟存储池数据
const storagePools = ref([
  { id: '1', name: 'pool1', availableSpace: '500GB' },
  { id: '2', name: 'pool2', availableSpace: '1TB' }
])

// 模拟用户数据
const availableUsers = ref([
  { id: 'u1', name: 'admin' },
  { id: 'u2', name: 'zhangsan' },
  { id: 'u3', name: 'lisi' }
])

// 共享文件夹数据
const folders = ref([
  {
    id: '1',
    name: 'Documents',
    path: '/pool1/Documents',
    storagePool: 'pool1',
    description: '文档文件夹',
    protocols: {
      smb: { enabled: true, shareName: 'Documents', recycleBin: true },
      ftp: { enabled: true, port: 21, passiveMode: true },
      sftp: { enabled: true, port: 22, allowShell: false },
      nfs: { enabled: false },
      webdav: { enabled: false, https: true, port: 443 }
    },
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'u2', userName: 'zhangsan', read: true, write: true, admin: false }
    ]
  },
  {
    id: '2',
    name: 'Photos',
    path: '/pool1/Photos',
    storagePool: 'pool1',
    description: '照片文件夹',
    protocols: {
      smb: { enabled: true, shareName: 'Photos', recycleBin: true },
      ftp: { enabled: false },
      sftp: { enabled: true, port: 22, allowShell: false },
      nfs: { enabled: true, version: '4', allowedClients: '192.168.1.0/24' },
      webdav: { enabled: false }
    },
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'u2', userName: 'zhangsan', read: true, write: false, admin: false }
    ]
  }
])

const expandedFolders = ref<string[]>([])
const showCreateFolder = ref(false)
const showAddPermission = ref(false)
const editingFolder = ref<any>(null)
const currentFolderId = ref('')

const folderForm = ref({
  name: '',
  storagePool: '',
  description: ''
})

const permissionForm = ref({
  userId: '',
  level: 'read'
})

// 方法
const toggleFolder = (folderId: string) => {
  const index = expandedFolders.value.indexOf(folderId)
  if (index > -1) {
    expandedFolders.value.splice(index, 1)
  } else {
    expandedFolders.value.push(folderId)
  }
}

const getEnabledProtocols = (folder: any) => {
  return Object.keys(folder.protocols).filter(key => folder.protocols[key].enabled)
}

const getProtocolLabel = (protocolId: string) => {
  const protocol = availableProtocols.find(p => p.id === protocolId)
  return protocol?.name || protocolId
}

const isProtocolEnabled = (folder: any, protocolId: string) => {
  return folder.protocols[protocolId]?.enabled || false
}

const toggleProtocol = (folderId: string, protocolId: string, event: Event) => {
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  const enabled = (event.target as HTMLInputElement).checked
  if (!folder.protocols[protocolId]) {
    folder.protocols[protocolId] = { enabled }
  }
  folder.protocols[protocolId].enabled = enabled
}

const getProtocolSetting = (folder: any, protocolId: string, setting: string) => {
  return folder.protocols[protocolId]?.[setting] || ''
}

const updateProtocolSetting = (folderId: string, protocolId: string, setting: string, event: Event) => {
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  const target = event.target as HTMLInputElement | HTMLSelectElement
  const value = target.type === 'checkbox' ? (target as HTMLInputElement).checked : target.value

  if (!folder.protocols[protocolId]) {
    folder.protocols[protocolId] = { enabled: true }
  }
  folder.protocols[protocolId][setting] = value
}

const getPermissionLevel = (permission: any): string => {
  if (permission.admin) return 'admin'
  if (permission.write) return 'write'
  if (permission.read) return 'read'
  return 'none'
}

const updatePermission = (folderId: string, userId: string, event: Event) => {
  const level = (event.target as HTMLSelectElement).value
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  const permission = folder.permissions.find((p: any) => p.userId === userId)
  if (!permission) return

  permission.read = level !== 'none'
  permission.write = level === 'write' || level === 'admin'
  permission.admin = level === 'admin'
}

const removePermission = (folderId: string, userId: string) => {
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return
  folder.permissions = folder.permissions.filter((p: any) => p.userId !== userId)
}

const addPermission = (folder: any) => {
  currentFolderId.value = folder.id
  permissionForm.value = {
    userId: '',
    level: 'read'
  }
  showAddPermission.value = true
}

const confirmAddPermission = () => {
  const folder = folders.value.find(f => f.id === currentFolderId.value)
  if (!folder || !permissionForm.value.userId) return

  const user = availableUsers.value.find(u => u.id === permissionForm.value.userId)
  if (!user) return

  const level = permissionForm.value.level
  folder.permissions.push({
    userId: user.id,
    userName: user.name,
    read: level !== 'none',
    write: level === 'write' || level === 'admin',
    admin: level === 'admin'
  })

  closePermissionModal()
}

const closePermissionModal = () => {
  showAddPermission.value = false
  permissionForm.value = {
    userId: '',
    level: 'read'
  }
}

const editFolder = (folder: any) => {
  editingFolder.value = folder
  folderForm.value = {
    name: folder.name,
    storagePool: folder.storagePool,
    description: folder.description
  }
  showCreateFolder.value = true
}

const deleteFolder = (folderId: string) => {
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  if (confirm(`确定要删除共享文件夹 "${folder.name}" 吗？`)) {
    folders.value = folders.value.filter(f => f.id !== folderId)
  }
}

const saveFolder = () => {
  if (editingFolder.value) {
    const index = folders.value.findIndex(f => f.id === editingFolder.value.id)
    if (index > -1) {
      folders.value[index] = {
        ...editingFolder.value,
        name: folderForm.value.name,
        storagePool: folderForm.value.storagePool,
        description: folderForm.value.description
      }
    }
  } else {
    folders.value.push({
      id: Date.now().toString(),
      name: folderForm.value.name,
      path: `/${folderForm.value.storagePool}/${folderForm.value.name}`,
      storagePool: folderForm.value.storagePool,
      description: folderForm.value.description,
      protocols: {
        smb: { enabled: false, shareName: folderForm.value.name, recycleBin: true },
        ftp: { enabled: false, port: 21, allowShell: false },
        sftp: { enabled: false, port: 22, allowShell: false },
        nfs: { enabled: false, port: 2049, allowShell: false },
        webdav: { enabled: false, port: 80, allowShell: false }
      },
      permissions: []
    })
  }
  closeCreateModal()
}

const closeCreateModal = () => {
  showCreateFolder.value = false
  editingFolder.value = null
  folderForm.value = {
    name: '',
    storagePool: '',
    description: ''
  }
}
</script>

<style scoped>
.share-folder-manager {
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.sfm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.sfm-header h1 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.sfm-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 16px;
  border-radius: 4px;
  font-size: 14px;
  border: 1px solid #d9d9d9;
  background: white;
  color: #333;
  cursor: pointer;
  transition: all 0.2s;
}

.sfm-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.sfm-btn.primary {
  background: #1890ff;
  border-color: #1890ff;
  color: white;
}

.sfm-btn.primary:hover {
  background: #40a9ff;
  border-color: #40a9ff;
}

.sfm-btn.secondary {
  background: white;
  border-color: #d9d9d9;
}

.sfm-btn.small {
  padding: 4px 12px;
  font-size: 13px;
}

.sfm-folders {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sfm-folder-item {
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.sfm-folder-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.sfm-folder-header:hover {
  background: #fafafa;
}

.folder-basic {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.expand-icon {
  width: 16px;
  height: 16px;
  color: #999;
  transition: transform 0.2s;
}

.sfm-folder-item.expanded .expand-icon {
  transform: rotate(90deg);
}

.folder-icon {
  width: 20px;
  height: 20px;
  color: #1890ff;
}

.folder-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.folder-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.folder-path {
  font-size: 12px;
  color: #999;
}

.folder-protocols {
  display: flex;
  gap: 6px;
}

.protocol-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
}

.protocol-badge.smb {
  background: #e6f7ff;
  color: #1890ff;
}

.protocol-badge.ftp {
  background: #f6ffed;
  color: #52c41a;
}

.protocol-badge.sftp {
  background: #fff0f6;
  color: #eb2f96;
}

.protocol-badge.nfs {
  background: #fff7e6;
  color: #fa8c16;
}

.protocol-badge.webdav {
  background: #f9f0ff;
  color: #722ed1;
}

.folder-actions {
  display: flex;
  gap: 16px;
}

.sfm-link {
  background: none;
  border: none;
  color: #1890ff;
  cursor: pointer;
  padding: 4px 8px;
  font-size: 13px;
  transition: all 0.2s;
}

.sfm-link:hover {
  color: #40a9ff;
}

.sfm-link.danger {
  color: #ff4d4f;
}

.sfm-link.danger:hover {
  color: #ff7875;
}

.sfm-folder-details {
  border-top: 1px solid #e8e8e8;
  background: #fafafa;
}

.detail-section {
  padding: 16px;
  border-bottom: 1px solid #e8e8e8;
}

.detail-section:last-child {
  border-bottom: none;
}

.detail-section h3 {
  font-size: 13px;
  font-weight: 600;
  color: #666;
  margin-bottom: 12px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.section-header h3 {
  margin-bottom: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: #999;
}

.info-value {
  font-size: 13px;
  color: #333;
}

.protocol-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.protocol-item {
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  padding: 12px;
}

.protocol-header {
  margin-bottom: 8px;
}

.protocol-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.protocol-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.protocol-desc {
  font-size: 12px;
  color: #999;
}

.protocol-settings {
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid #e8e8e8;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.setting-item label {
  font-size: 12px;
  color: #666;
}

.setting-item input,
.setting-item select {
  padding: 4px 8px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 13px;
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
  padding: 8px 12px;
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  font-size: 14px;
  color: #333;
}

.permission-level select {
  padding: 4px 8px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 13px;
  background: white;
  cursor: pointer;
}

.sfm-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.sfm-modal-content {
  background: white;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.sfm-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e8e8e8;
}

.sfm-modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.sfm-close {
  width: 32px;
  height: 32px;
  border-radius: 4px;
  border: none;
  background: transparent;
  color: #999;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.sfm-close:hover {
  background: #f5f5f5;
}

.sfm-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.sfm-form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.sfm-form-group label {
  font-size: 13px;
  font-weight: 500;
  color: #666;
}

.sfm-form-group input,
.sfm-form-group select,
.sfm-form-group textarea {
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
}

.permission-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.radio-option {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.radio-option:hover {
  border-color: #1890ff;
  background: #f0f8ff;
}

.radio-option input[type="radio"] {
  margin-top: 2px;
}

.radio-content {
  flex: 1;
}

.radio-title {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.radio-desc {
  font-size: 12px;
  color: #999;
}

.sfm-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 24px;
  border-top: 1px solid #e8e8e8;
}
</style>