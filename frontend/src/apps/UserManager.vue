<template>
  <div class="user-manager-app">
    <!-- 顶部标签页 -->
    <div class="uma-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['uma-tab', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 用户管理 -->
    <div v-if="activeTab === 'users'" class="uma-content">
      <div class="uma-toolbar">
        <div class="uma-title">用户</div>
        <div class="uma-actions">
          <button class="uma-btn" @click="showAddUser = true">
            <PlusIcon class="w-4 h-4" />
            新增
          </button>
        </div>
      </div>

      <!-- 用户列表 -->
      <div class="uma-table">
        <table>
          <thead>
            <tr>
              <th>用户</th>
              <th>用户组</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.id">
              <td>{{ user.displayName }}</td>
              <td>{{ user.groups.join(', ') }}</td>
              <td>
                <span :class="['uma-status', user.status === 'active' ? 'active' : 'disabled']">
                  {{ user.status === 'active' ? '正常' : '停用' }}
                </span>
              </td>
              <td>
                <button class="uma-link" @click="editUser(user)">编辑</button>
                <button class="uma-link danger" @click="deleteUser(user.id)">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 用户组管理 -->
    <div v-if="activeTab === 'groups'" class="uma-content">
      <div class="uma-toolbar">
        <div class="uma-title">用户组</div>
        <div class="uma-actions">
          <button class="uma-btn" @click="showAddGroup = true">
            <PlusIcon class="w-4 h-4" />
            新增
          </button>
        </div>
      </div>

      <!-- 用户组列表 -->
      <div class="uma-groups">
        <div
          v-for="group in groups"
          :key="group.id"
          class="uma-group-item"
          :class="{ expanded: expandedGroups.includes(group.id) }"
        >
          <div class="uma-group-header" @click="toggleGroup(group.id)">
            <div class="group-info">
              <ChevronRightIcon class="expand-icon" />
              <span class="group-name">{{ group.name }}</span>
            </div>
            <div class="group-actions">
              <button class="uma-link" @click.stop="editGroup(group)">编辑</button>
              <button
                v-if="!isSystemGroup(group)"
                class="uma-link danger"
                @click.stop="deleteGroup(group.id)"
              >
                删除
              </button>
            </div>
          </div>

          <div v-if="expandedGroups.includes(group.id)" class="uma-group-details">
            <div class="group-description">{{ group.description }}</div>
            <div class="group-permissions">
              <strong>权限：</strong>
              <span v-for="perm in group.permissions" :key="perm" class="perm-tag">{{ perm }}</span>
            </div>
            <div class="group-members">
              <strong>成员 ({{ group.memberCount }})：</strong>
              <span v-for="member in group.members" :key="member" class="member-tag">{{ member }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 权限管理 -->
    <div v-if="activeTab === 'permissions'" class="uma-content">
      <div class="uma-toolbar">
        <div class="uma-title">权限管理</div>
        <div class="uma-actions">
          <button class="uma-btn" @click="createFolder">
            <PlusIcon class="w-4 h-4" />
            新增共享文件夹
          </button>
        </div>
      </div>

      <!-- 共享文件夹列表 -->
      <div class="uma-folders">
        <div
          v-for="folder in folders"
          :key="folder.id"
          class="uma-folder-item"
          :class="{ expanded: expandedFolders.includes(folder.id) }"
        >
          <div class="uma-folder-header" @click="toggleFolder(folder.id)">
            <div class="folder-info">
              <ChevronRightIcon class="expand-icon" />
              <FolderIcon class="folder-icon" />
              <span class="folder-name">{{ folder.name }}</span>
            </div>
            <div class="folder-actions">
              <button class="uma-link" @click.stop="editFolder(folder)">编辑</button>
              <button class="uma-link danger" @click.stop="deleteFolder(folder.id)">删除</button>
            </div>
          </div>

          <div v-if="expandedFolders.includes(folder.id)" class="uma-folder-details">
            <!-- 权限列表 -->
            <div class="permissions-section">
              <div class="section-header">
                <h3>用户权限</h3>
                <button class="uma-btn small" @click.stop="addPermission(folder)">
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
                    class="uma-link danger"
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

    <!-- 应用权限管理 -->
    <div v-if="activeTab === 'apps'" class="uma-content">
      <div class="uma-toolbar">
        <div class="uma-title">应用权限</div>
      </div>

      <!-- 应用列表 -->
      <div class="uma-apps">
        <div class="app-item" v-for="app in applications" :key="app.id">
          <div class="app-header" @click="toggleApp(app.id)">
            <div class="app-info">
              <ChevronRightIcon class="expand-icon" :class="{ expanded: expandedApps.includes(app.id) }" />
              <component :is="getAppIcon(app.icon)" class="app-icon" />
              <div class="app-details">
                <div class="app-name">{{ app.name }}</div>
                <div class="app-description">{{ app.description }}</div>
              </div>
            </div>
            <div class="app-status">
              <span :class="['status-indicator', app.enabled ? 'enabled' : 'disabled']">
                {{ app.enabled ? '已启用' : '未启用' }}
              </span>
            </div>
          </div>

          <div v-if="expandedApps.includes(app.id)" class="app-details-panel">
            <!-- Immich 用户管理 -->
            <div v-if="app.id === 'immich'" class="app-user-management">
              <ImmichUserManager />
            </div>

            <!-- 其他应用的用户管理占位 -->
            <div v-else class="app-placeholder">
              <div class="placeholder-content">
                <component :is="getAppIcon(app.icon)" class="placeholder-icon" />
                <h4>{{ app.name }} 用户管理</h4>
                <p>{{ app.name }} 的用户管理功能即将推出</p>
                <button class="uma-btn secondary" @click="navigateToApp(app)">前往管理</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加用户对话框 -->
    <div v-if="showAddUser" class="uma-modal" @click.self="closeUserModal">
      <div class="uma-modal-content">
        <div class="uma-modal-header">
          <h3>{{ editingUser ? '编辑用户' : '新增用户' }}</h3>
          <button class="uma-close" @click="closeUserModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveUser" class="uma-form">
          <div class="uma-form-group">
            <label>用户名 *</label>
            <input
              v-model="userForm.name"
              type="text"
              required
              :disabled="editingUser"
              placeholder="用户名"
            />
          </div>

          <div class="uma-form-group">
            <label>显示名称</label>
            <input v-model="userForm.displayName" type="text" placeholder="显示名称" />
          </div>

          <div class="uma-form-group" v-if="!editingUser">
            <label>密码 *</label>
            <input v-model="userForm.password" type="password" required placeholder="密码" />
          </div>

          <div class="uma-form-group">
            <label>用户组</label>
            <div class="uma-checkbox-group">
              <label
                v-for="group in groups"
                :key="group.id"
                class="uma-checkbox"
              >
                <input type="checkbox" v-model="userForm.groups" :value="group.name" />
                {{ group.name }}
              </label>
            </div>
          </div>

          <div class="uma-modal-footer">
            <button type="button" class="uma-btn secondary" @click="closeUserModal">
              取消
            </button>
            <button type="submit" class="uma-btn primary">
              {{ editingUser ? '确定' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 添加用户组对话框 -->
    <div v-if="showAddGroup" class="uma-modal" @click.self="closeGroupModal">
      <div class="uma-modal-content">
        <div class="uma-modal-header">
          <h3>{{ editingGroup ? '编辑用户组' : '新增用户组' }}</h3>
          <button class="uma-close" @click="closeGroupModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveGroup" class="uma-form">
          <div class="uma-form-group">
            <label>用户组名称 *</label>
            <input
              v-model="groupForm.name"
              type="text"
              required
              :disabled="editingGroup?.isSystem"
              placeholder="用户组名称"
            />
          </div>

          <div class="uma-form-group">
            <label>描述</label>
            <textarea
              v-model="groupForm.description"
              rows="3"
              placeholder="用户组描述"
            ></textarea>
          </div>

          <div class="uma-modal-footer">
            <button type="button" class="uma-btn secondary" @click="closeGroupModal">
              取消
            </button>
            <button type="submit" class="uma-btn primary">
              {{ editingGroup ? '确定' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 添加权限模态框 -->
    <div v-if="showAddPermissionModal" class="uma-modal" @click.self="closeAddPermissionModal">
      <div class="uma-modal-content">
        <div class="uma-modal-header">
          <h3>添加用户权限</h3>
          <button class="uma-close" @click="closeAddPermissionModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="confirmAddPermission" class="uma-form">
          <div class="uma-form-group">
            <label>选择用户</label>
            <select v-model="newPermission.userId" required>
              <option value="">选择用户...</option>
              <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                {{ user.name }}
              </option>
            </select>
          </div>

          <div class="uma-form-group">
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

          <div class="uma-modal-footer">
            <button type="button" class="uma-btn secondary" @click="closeAddPermissionModal">取消</button>
            <button type="submit" class="uma-btn primary">添加</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  PlusIcon,
  ChevronRightIcon,
  XMarkIcon,
  FolderIcon,
  UserIcon,
  CubeIcon,
  PhotoIcon,
  FolderOpenIcon,
  ChartBarIcon
} from '@heroicons/vue/24/outline'
import ImmichUserManager from '../components/ImmichUserManager.vue'

const activeTab = ref('users')

const tabs = [
  { id: 'users', label: '用户' },
  { id: 'groups', label: '用户组' },
  { id: 'permissions', label: '权限' },
  { id: 'apps', label: '应用权限' }
]

// 用户数据
const users = ref([
  { id: '1', displayName: 'admin', name: 'admin', groups: ['administrators'], status: 'active' },
  { id: '2', displayName: '张三', name: 'zhangsan', groups: ['users'], status: 'active' },
  { id: '3', displayName: '李四', name: 'lisi', groups: ['users'], status: 'active' }
])

// 用户组数据
const groups = ref([
  {
    id: 'administrators',
    name: 'administrators',
    description: '系统管理员组，拥有完整的管理权限',
    isSystem: true,
    memberCount: 1,
    members: ['admin'],
    permissions: ['用户管理', '系统设置', '文件管理']
  },
  {
    id: 'users',
    name: 'users',
    description: '普通用户组，基本的文件访问权限',
    isSystem: true,
    memberCount: 3,
    members: ['zhangsan', 'lisi'],
    permissions: ['文件访问', '个人设置']
  }
])

// 共享文件夹数据
const folders = ref([
  {
    id: '1',
    name: 'Documents',
    path: '/home/Documents',
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'u2', userName: 'zhangsan', read: true, write: true, admin: false }
    ]
  },
  {
    id: '2',
    name: 'Photos',
    path: '/home/Photos',
    permissions: [
      { userId: 'u1', userName: 'admin', read: true, write: true, admin: true },
      { userId: 'u2', userName: 'zhangsan', read: true, write: false, admin: false }
    ]
  }
])

const expandedGroups = ref<string[]>([])
const expandedFolders = ref<string[]>([])
const expandedApps = ref<string[]>(['immich']) // 默认展开Immich

// 应用列表
const applications = ref([
  {
    id: 'immich',
    name: 'Immich',
    description: '照片管理应用',
    icon: 'PhotoIcon',
    enabled: true,
    url: 'http://localhost:2283'
  },
  {
    id: 'filemanager',
    name: '文件管理器',
    description: '文件浏览和管理',
    icon: 'FolderOpenIcon',
    enabled: true,
    url: '/filemanager'
  },
  {
    id: 'docker',
    name: 'Docker管理',
    description: '容器管理',
    icon: 'CubeIcon',
    enabled: true,
    url: '/docker'
  },
  {
    id: 'monitor',
    name: '系统监控',
    description: '系统性能监控',
    icon: 'ChartBarIcon',
    enabled: true,
    url: '/monitor'
  }
])

const showAddUser = ref(false)
const showAddGroup = ref(false)
const showAddPermissionModal = ref(false)
const editingUser = ref<any>(null)
const editingGroup = ref<any>(null)
const currentFolderId = ref('')

const userForm = ref({
  name: '',
  displayName: '',
  password: '',
  groups: ['users']
})

const groupForm = ref({
  name: '',
  description: ''
})

const newPermission = ref({
  userId: '',
  level: 'read'
})

const availableUsers = ref([
  { id: 'u1', name: 'admin' },
  { id: 'u2', name: 'zhangsan' },
  { id: 'u3', name: 'lisi' }
])

// 方法
const toggleGroup = (groupId: string) => {
  const index = expandedGroups.value.indexOf(groupId)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(groupId)
  }
}

const toggleFolder = (folderId: string) => {
  const index = expandedFolders.value.indexOf(folderId)
  if (index > -1) {
    expandedFolders.value.splice(index, 1)
  } else {
    expandedFolders.value.push(folderId)
  }
}

const isSystemGroup = (group: any) => group.isSystem

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
  newPermission.value = {
    userId: '',
    level: 'read'
  }
  showAddPermissionModal.value = true
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

  closeAddPermissionModal()
}

const closeAddPermissionModal = () => {
  showAddPermissionModal.value = false
  newPermission.value = {
    userId: '',
    level: 'read'
  }
}

const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    name: user.name,
    displayName: user.displayName,
    password: '',
    groups: [...user.groups]
  }
  showAddUser.value = true
}

const deleteUser = (userId: string) => {
  if (confirm('确定要删除此用户吗？')) {
    users.value = users.value.filter(u => u.id !== userId)
  }
}

const editGroup = (group: any) => {
  editingGroup.value = group
  groupForm.value = {
    name: group.name,
    description: group.description
  }
  showAddGroup.value = true
}

const deleteGroup = (groupId: string) => {
  if (confirm('确定要删除此用户组吗？')) {
    groups.value = groups.value.filter(g => g.id !== groupId)
  }
}

const saveUser = () => {
  if (editingUser.value) {
    const index = users.value.findIndex(u => u.id === editingUser.value.id)
    if (index > -1) {
      users.value[index] = {
        ...editingUser.value,
        name: userForm.value.name,
        displayName: userForm.value.displayName,
        groups: userForm.value.groups,
        status: 'active'
      }
    }
  } else {
    users.value.push({
      id: Date.now().toString(),
      name: userForm.value.name,
      displayName: userForm.value.displayName,
      groups: userForm.value.groups,
      status: 'active'
    })
  }
  closeUserModal()
}

const saveGroup = () => {
  if (editingGroup.value) {
    const index = groups.value.findIndex(g => g.id === editingGroup.value.id)
    if (index > -1) {
      groups.value[index] = {
        ...editingGroup.value,
        name: groupForm.value.name,
        description: groupForm.value.description
      }
    }
  } else {
    groups.value.push({
      id: Date.now().toString(),
      name: groupForm.value.name,
      description: groupForm.value.description,
      isSystem: false,
      memberCount: 0,
      members: [],
      permissions: ['文件访问']
    })
  }
  closeGroupModal()
}

const closeUserModal = () => {
  showAddUser.value = false
  editingUser.value = null
  userForm.value = {
    name: '',
    displayName: '',
    password: '',
    groups: ['users']
  }
}

const closeGroupModal = () => {
  showAddGroup.value = false
  editingGroup.value = null
  groupForm.value = {
    name: '',
    description: ''
  }
}

// 应用相关方法
const toggleApp = (appId: string) => {
  const index = expandedApps.value.indexOf(appId)
  if (index > -1) {
    expandedApps.value.splice(index, 1)
  } else {
    expandedApps.value.push(appId)
  }
}

const getAppIcon = (iconName: string) => {
  const iconMap: Record<string, any> = {
    PhotoIcon,
    FolderOpenIcon,
    CubeIcon,
    ChartBarIcon
  }
  return iconMap[iconName] || CubeIcon
}

const navigateToApp = (app: any) => {
  if (app.url.startsWith('http')) {
    window.open(app.url, '_blank')
  } else {
    window.location.href = app.url
  }
}

const createFolder = () => {
  console.log('创建新共享文件夹')
}

const editFolder = (folder: any) => {
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
.user-manager-app {
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.uma-tabs {
  display: flex;
  border-bottom: 1px solid #d8d8d8;
  background: white;
}

.uma-tab {
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.uma-tab:hover {
  color: #333;
}

.uma-tab.active {
  color: #1890ff;
  border-bottom-color: #1890ff;
}

.uma-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.uma-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e8e8e8;
}

.uma-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.uma-actions {
  display: flex;
  gap: 8px;
}

.uma-btn {
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

.uma-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.uma-btn.primary {
  background: #1890ff;
  border-color: #1890ff;
  color: white;
}

.uma-btn.primary:hover {
  background: #40a9ff;
  border-color: #40a9ff;
}

.uma-btn.secondary {
  background: white;
  border-color: #d9d9d9;
}

.uma-btn.small {
  padding: 4px 12px;
  font-size: 13px;
}

.uma-table {
  background: white;
  border: 1px solid #e8e8e8;
  border-collapse: collapse;
}

.uma-table table {
  width: 100%;
  border-collapse: collapse;
}

.uma-table th,
.uma-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e8e8e8;
}

.uma-table th {
  background: #fafafa;
  font-weight: 500;
  color: #666;
  font-size: 13px;
}

.uma-table td {
  font-size: 14px;
  color: #333;
}

.uma-status {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.uma-status.active {
  background: #52c41a;
  color: white;
}

.uma-status.disabled {
  background: #d9d9d9;
  color: #666;
}

.uma-link {
  background: none;
  border: none;
  color: #1890ff;
  cursor: pointer;
  padding: 4px 8px;
  font-size: 13px;
  transition: all 0.2s;
}

.uma-link:hover {
  color: #40a9ff;
}

.uma-link.danger {
  color: #ff4d4f;
}

.uma-link.danger:hover {
  color: #ff7875;
}

.uma-groups {
  background: white;
  border: 1px solid #e8e8e8;
}

.uma-group-item {
  border-bottom: 1px solid #e8e8e8;
}

.uma-group-item:last-child {
  border-bottom: none;
}

.uma-group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.uma-group-header:hover {
  background: #fafafa;
}

.group-info {
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

.uma-group-item.expanded .expand-icon {
  transform: rotate(90deg);
}

.group-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.group-actions {
  display: flex;
  gap: 16px;
}

.uma-group-details {
  padding: 16px;
  background: #fafafa;
  border-top: 1px solid #e8e8e8;
}

.group-description {
  font-size: 13px;
  color: #666;
  margin-bottom: 12px;
}

.group-permissions,
.group-members {
  font-size: 13px;
  color: #666;
  margin-bottom: 8px;
}

.perm-tag,
.member-tag {
  display: inline-block;
  padding: 2px 6px;
  background: white;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  margin-right: 8px;
  font-size: 12px;
}

.uma-folders {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.uma-folder-item {
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: hidden;
}

.uma-folder-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.uma-folder-header:hover {
  background: #fafafa;
}

.folder-info {
  display: flex;
  align-items: center;
  gap: 8px;
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

.uma-folder-details {
  border-top: 1px solid #e8e8e8;
  background: #fafafa;
}

.permissions-section {
  padding: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
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

.uma-modal {
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

.uma-modal-content {
  background: white;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.uma-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e8e8e8;
}

.uma-modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.uma-close {
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

.uma-close:hover {
  background: #f5f5f5;
}

.uma-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.uma-form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.uma-form-group label {
  font-size: 13px;
  font-weight: 500;
  color: #666;
}

.uma-form-group input,
.uma-form-group select,
.uma-form-group textarea {
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
}

.uma-checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.uma-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #333;
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

.uma-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 24px;
  border-top: 1px solid #e8e8e8;
}

/* 应用权限样式 */
.uma-apps {
  padding: 16px;
}

.app-item {
  background: white;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  margin-bottom: 12px;
  overflow: hidden;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.app-header:hover {
  background: #fafafa;
}

.app-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-icon {
  width: 24px;
  height: 24px;
  color: #1890ff;
}

.app-details {
  display: flex;
  flex-direction: column;
}

.app-name {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.app-description {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.app-status {
  display: flex;
  align-items: center;
}

.status-indicator {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-indicator.enabled {
  background: #f6ffed;
  color: #52c41a;
}

.status-indicator.disabled {
  background: #fff2f0;
  color: #ff4d4f;
}

.app-details-panel {
  border-top: 1px solid #e8e8e8;
  background: #fafafa;
}

.app-user-management {
  padding: 16px;
}

.app-placeholder {
  padding: 40px 16px;
}

.placeholder-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

.placeholder-icon {
  width: 48px;
  height: 48px;
  color: #d9d9d9;
  margin-bottom: 16px;
}

.placeholder-content h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
  color: #333;
}

.placeholder-content p {
  margin: 0 0 16px 0;
  font-size: 14px;
  color: #999;
}
</style>