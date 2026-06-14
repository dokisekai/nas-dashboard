<template>
  <div class="simple-permission-manager">
    <!-- 头部 -->
    <div class="spm-header">
      <h1>权限管理</h1>
      <button class="spm-btn primary" @click="createFolder">
        <PlusIcon class="w-4 h-4" />
        新增共享文件夹
      </button>
    </div>

    <!-- 共享文件夹列表 -->
    <div class="spm-folders">
      <div
        v-for="folder in folders"
        :key="folder.id"
        class="spm-folder-item"
        :class="{ expanded: expandedFolders.includes(folder.id) }"
      >
        <div class="spm-folder-header" @click="toggleFolder(folder.id)">
          <div class="folder-info">
            <ChevronRightIcon class="expand-icon" />
            <FolderIcon class="folder-icon" />
            <span class="folder-name">{{ folder.name }}</span>
          </div>
          <div class="folder-actions">
            <button class="spm-link" @click.stop="editFolder(folder)">编辑</button>
            <button class="spm-link danger" @click.stop="deleteFolder(folder.id)">删除</button>
          </div>
        </div>

        <div v-if="expandedFolders.includes(folder.id)" class="spm-folder-details">
          <!-- 权限列表 -->
          <div class="permissions-section">
            <div class="section-header">
              <h3>用户权限</h3>
              <button class="spm-btn small" @click.stop="addPermission(folder)">
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
                  class="spm-link danger"
                  @click="removePermission(folder.id, permission.userId)"
                >
                  移除
                </button>
              </div>
            </div>
          </div>

          <!-- 高级选项（折叠） -->
          <div class="advanced-section">
            <div class="section-header" @click="toggleAdvanced(folder.id)">
              <h3>高级选项</h3>
              <ChevronDownIcon
                class="w-4 h-4"
                :class="{ 'rotate-180': showAdvanced.includes(folder.id) }"
              />
            </div>

            <div v-if="showAdvanced.includes(folder.id)" class="advanced-options">
              <label class="option-item">
                <input type="checkbox" v-model="folder.recycleBin" />
                <span>启用回收站</span>
              </label>
              <label class="option-item">
                <input type="checkbox" v-model="folder.hideHidden" />
                <span>隐藏隐藏文件</span>
              </label>
              <label class="option-item">
                <input type="checkbox" v-model="folder.browseable" />
                <span>允许浏览</span>
              </label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加权限模态框 -->
    <div v-if="showAddModal" class="spm-modal" @click.self="closeAddModal">
      <div class="spm-modal-content">
        <div class="spm-modal-header">
          <h3>添加用户权限</h3>
          <button class="spm-close" @click="closeAddModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="confirmAddPermission" class="spm-form">
          <div class="spm-form-group">
            <label>选择用户</label>
            <select v-model="newPermission.userId" required>
              <option value="">选择用户...</option>
              <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                {{ user.name }}
              </option>
            </select>
          </div>

          <div class="spm-form-group">
            <label>权限级别</label>
            <div class="permission-options">
              <label class="radio-option">
                <input type="radio" v-model="newPermission.level" value="read" />
                <div class="radio-content">
                  <div class="radio-title">只读</div>
                  <div class="radio-desc">可以浏览和下载文件</div>
                </div>
              </label>
              <label class="radio-option">
                <input type="radio" v-model="newPermission.level" value="write" />
                <div class="radio-content">
                  <div class="radio-title">读写</div>
                  <div class="radio-desc">可以读取、上传和修改文件</div>
                </div>
              </label>
              <label class="radio-option">
                <input type="radio" v-model="newPermission.level" value="admin" />
                <div class="radio-content">
                  <div class="radio-title">管理</div>
                  <div class="radio-desc">完全控制权限</div>
                </div>
              </label>
            </div>
          </div>

          <div class="spm-modal-footer">
            <button type="button" class="spm-btn secondary" @click="closeAddModal">取消</button>
            <button type="submit" class="spm-btn primary">添加</button>
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
  ChevronDownIcon,
  FolderIcon,
  UserIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

interface Permission {
  userId: string
  userName: string
  read: boolean
  write: boolean
  admin: boolean
}

interface Folder {
  id: string
  name: string
  path: string
  permissions: Permission[]
  recycleBin: boolean
  hideHidden: boolean
  browseable: boolean
}

const expandedFolders = ref<string[]>([])
const showAdvanced = ref<string[]>([])
const showAddModal = ref(false)
const currentFolderId = ref('')

const folders = ref<Folder[]>([
  {
    id: '1',
    name: 'Documents',
    path: '/home/Documents',
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'u2', userName: 'john', read: true, write: true, admin: false }
    ],
    recycleBin: true,
    hideHidden: true,
    browseable: true
  },
  {
    id: '2',
    name: 'Photos',
    path: '/home/Photos',
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'u2', userName: 'john', read: true, write: false, admin: false }
    ],
    recycleBin: true,
    hideHidden: true,
    browseable: true
  },
  {
    id: '3',
    name: 'Public',
    path: '/home/Public',
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'everyone', userName: 'Everyone', read: true, write: true, admin: false }
    ],
    recycleBin: true,
    hideHidden: false,
    browseable: true
  }
])

const availableUsers = ref([
  { id: 'u1', name: 'admin' },
  { id: 'u2', name: 'john' },
  { id: 'u3', name: 'jane' },
  { id: 'u4', name: 'guest' }
])

const newPermission = ref({
  userId: '',
  level: 'read'
})

const toggleFolder = (folderId: string) => {
  const index = expandedFolders.value.indexOf(folderId)
  if (index > -1) {
    expandedFolders.value.splice(index, 1)
  } else {
    expandedFolders.value.push(folderId)
  }
}

const toggleAdvanced = (folderId: string) => {
  const index = showAdvanced.value.indexOf(folderId)
  if (index > -1) {
    showAdvanced.value.splice(index, 1)
  } else {
    showAdvanced.value.push(folderId)
  }
}

const getPermissionLevel = (permission: Permission): string => {
  if (permission.admin) return 'admin'
  if (permission.write) return 'write'
  if (permission.read) return 'read'
  return 'none'
}

const updatePermission = (folderId: string, userId: string, event: Event) => {
  const level = (event.target as HTMLSelectElement).value
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  const permission = folder.permissions.find(p => p.userId === userId)
  if (!permission) return

  // 根据级别设置权限
  permission.read = level !== 'none'
  permission.write = level === 'write' || level === 'admin'
  permission.admin = level === 'admin'
}

const removePermission = (folderId: string, userId: string) => {
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  folder.permissions = folder.permissions.filter(p => p.userId !== userId)
}

const addPermission = (folder: Folder) => {
  currentFolderId.value = folder.id
  newPermission.value = {
    userId: '',
    level: 'read'
  }
  showAddModal.value = true
}

const confirmAddPermission = () => {
  const folder = folders.value.find(f => f.id === currentFolderId.value)
  if (!folder || !newPermission.value.userId) return

  const user = availableUsers.value.find(u => u.id === newPermission.value.userId)
  if (!user) return

  const level = newPermission.value.level
  folder.permissions.push({
    userId: user.id,
    userName: user.name,
    read: level !== 'none',
    write: level === 'write' || level === 'admin',
    admin: level === 'admin'
  })

  closeAddModal()
}

const closeAddModal = () => {
  showAddModal.value = false
  newPermission.value = {
    userId: '',
    level: 'read'
  }
}

const createFolder = () => {
  console.log('创建新共享文件夹')
}

const editFolder = (folder: Folder) => {
  console.log('编辑文件夹:', folder.name)
}

const deleteFolder = (folderId: string) => {
  const folder = folders.value.find(f => f.id === folderId)
  if (!folder) return

  if (confirm(`确定要删除共享文件夹 "${folder.name}" 吗？`)) {
    folders.value = folders.value.filter(f => f.id !== folderId)
  }
}
</script>

<style scoped>
.simple-permission-manager {
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.spm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  background: white;
  border-bottom: 1px solid #e8e8e8;
}

.spm-header h1 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.spm-btn {
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

.spm-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.spm-btn.primary {
  background: #1890ff;
  border-color: #1890ff;
  color: white;
}

.spm-btn.primary:hover {
  background: #40a9ff;
  border-color: #40a9ff;
}

.spm-btn.secondary {
  background: white;
  border-color: #d9d9d9;
}

.spm-btn.small {
  padding: 4px 12px;
  font-size: 13px;
}

.spm-folders {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.spm-folder-item {
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.spm-folder-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.spm-folder-header:hover {
  background: #fafafa;
}

.folder-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.expand-icon {
  width: 16px;
  height: 16px;
  color: #999;
  transition: transform 0.2s;
}

.spm-folder-item.expanded .expand-icon {
  transform: rotate(90deg);
}

.folder-icon {
  width: 20px;
  height: 20px;
  color: #1890ff;
}

.folder-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.folder-actions {
  display: flex;
  gap: 16px;
}

.spm-link {
  background: none;
  border: none;
  color: #1890ff;
  cursor: pointer;
  padding: 4px 8px;
  font-size: 13px;
  transition: all 0.2s;
}

.spm-link:hover {
  color: #40a9ff;
}

.spm-link.danger {
  color: #ff4d4f;
}

.spm-link.danger:hover {
  color: #ff7875;
}

.spm-folder-details {
  border-top: 1px solid #e8e8e8;
  background: #fafafa;
}

.permissions-section,
.advanced-section {
  padding: 16px;
}

.permissions-section + .advanced-section {
  border-top: 1px solid #e8e8e8;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  cursor: pointer;
}

.section-header h3 {
  font-size: 13px;
  font-weight: 600;
  color: #666;
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

.advanced-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 8px;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #666;
  cursor: pointer;
}

.spm-modal {
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

.spm-modal-content {
  background: white;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.spm-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e8e8e8;
}

.spm-modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.spm-close {
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

.spm-close:hover {
  background: #f5f5f5;
}

.spm-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.spm-form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.spm-form-group label {
  font-size: 13px;
  font-weight: 500;
  color: #666;
}

.spm-form-group select {
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

.spm-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 24px;
  border-top: 1px solid #e8e8e8;
}

.rotate-180 {
  transform: rotate(180deg);
}
</style>