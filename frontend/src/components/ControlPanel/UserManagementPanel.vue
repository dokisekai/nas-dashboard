<template>
  <div class="user-management-panel">
    <!-- 工具栏 -->
    <div class="panel-toolbar">
      <div class="toolbar-left">
        <h3>用户管理</h3>
        <p class="subtitle">管理系统用户账户和权限</p>
      </div>
      <div class="toolbar-right">
        <div class="search-box">
          <MagnifyingGlassIcon class="w-4 h-4" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索用户..."
          />
        </div>
        <div class="filter-dropdown">
          <select v-model="roleFilter">
            <option value="">所有角色</option>
            <option value="superadmin">超级管理员</option>
            <option value="admin">管理员</option>
            <option value="user">普通用户</option>
            <option value="guest">访客</option>
          </select>
        </div>
        <button class="btn btn-primary" @click="showCreateUserModal = true">
          <UserPlusIcon class="w-4 h-4" />
          添加用户
        </button>
      </div>
    </div>

    <!-- 用户统计 -->
    <div class="stats-bar">
      <div class="stat-item">
        <div class="stat-value">{{ users.filter(u => u.status === 'active').length }}</div>
        <div class="stat-label">活跃用户</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ users.filter(u => u.status === 'disabled').length }}</div>
        <div class="stat-label">已禁用</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ users.filter(u => u.role === 'admin' || u.role === 'superadmin').length }}</div>
        <div class="stat-label">管理员</div>
      </div>
      <div class="stat-item">
        <div class="stat-value">{{ getTotalStorageUsed() }}</div>
        <div class="stat-label">存储使用</div>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="users-table-container">
      <table class="users-table">
        <thead>
          <tr>
            <th width="40">
              <input type="checkbox" v-model="selectAll" />
            </th>
            <th>用户</th>
            <th>角色</th>
            <th>用户组</th>
            <th>状态</th>
            <th>存储配额</th>
            <th>最后登录</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="user in filteredUsers"
            :key="user.id"
            :class="{ 'row-disabled': user.status !== 'active' }"
          >
            <td>
              <input
                type="checkbox"
                v-model="selectedUsers"
                :value="user.id"
              />
            </td>
            <td>
              <div class="user-cell">
                <div class="user-avatar">
                  {{ user.username.charAt(0).toUpperCase() }}
                  <span
                    v-if="user.role === 'superadmin' || user.role === 'admin'"
                    class="admin-badge"
                  >👑</span>
                </div>
                <div class="user-details">
                  <div class="user-name">{{ user.username }}</div>
                  <div class="user-email">{{ user.email || user.fullName || '-' }}</div>
                </div>
              </div>
            </td>
            <td>
              <span
                :class="['role-badge', `role-${user.role}`]"
              >
                {{ getRoleLabel(user.role) }}
              </span>
            </td>
            <td>
              <div class="groups-list">
                <span
                  v-for="(group, index) in user.groups.slice(0, 2)"
                  :key="group"
                  class="group-tag"
                >
                  {{ getGroupName(group) }}
                  <span v-if="index === 1 && user.groups.length > 2">+{{ user.groups.length - 2 }}</span>
                </span>
              </div>
            </td>
            <td>
              <span
                :class="['status-badge', `status-${user.status}`]"
              >
                {{ getStatusLabel(user.status) }}
              </span>
            </td>
            <td>
              <div v-if="user.storageQuota?.enabled" class="quota-info">
                <div class="quota-bar">
                  <div
                    class="quota-fill"
                    :style="{ width: getQuotaPercent(user) + '%' }"
                  ></div>
                </div>
                <span class="quota-text">{{ formatBytes(user.storageQuota.used) }} / {{ formatBytes(user.storageQuota.size) }}</span>
              </div>
              <span v-else class="text-muted">无限制</span>
            </td>
            <td>
              <span class="text-muted">{{ formatDate(user.lastLogin) }}</span>
            </td>
            <td>
              <div class="action-buttons">
                <button
                  class="icon-btn"
                  @click="editUser(user)"
                  title="编辑用户"
                >
                  <PencilIcon class="w-4 h-4" />
                </button>
                <button
                  class="icon-btn"
                  @click="managePermissions(user)"
                  title="权限管理"
                >
                  <ShieldCheckIcon class="w-4 h-4" />
                </button>
                <button
                  class="icon-btn"
                  @click="toggleUserStatus(user)"
                  :title="user.status === 'active' ? '禁用用户' : '启用用户'"
                >
                  <NoSymbolIcon v-if="user.status === 'active'" class="w-4 h-4" />
                  <CheckIcon v-else class="w-4 h-4" />
                </button>
                <button
                  class="icon-btn danger"
                  @click="deleteUser(user)"
                  title="删除用户"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 批量操作栏 -->
    <div v-if="selectedUsers.length > 0" class="batch-actions">
      <div class="batch-info">
        已选择 {{ selectedUsers.length }} 个用户
      </div>
      <div class="batch-buttons">
        <button class="btn btn-secondary" @click="batchEnable">
          <CheckIcon class="w-4 h-4" />
          批量启用
        </button>
        <button class="btn btn-secondary" @click="batchDisable">
          <NoSymbolIcon class="w-4 h-4" />
          批量禁用
        </button>
        <button class="btn btn-secondary" @click="batchDelete">
          <TrashIcon class="w-4 h-4" />
          批量删除
        </button>
        <button class="btn btn-ghost" @click="selectedUsers = []">
          取消选择
        </button>
      </div>
    </div>

    <!-- 创建/编辑用户对话框 -->
    <div v-if="showCreateUserModal || showEditUserModal" class="modal-overlay" @click.self="closeUserModal">
      <div class="modal-content large-modal">
        <div class="modal-header">
          <h3>{{ editingUser ? '编辑用户' : '创建用户' }}</h3>
          <button class="close-btn" @click="closeUserModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit.prevent="saveUser" class="user-form">
          <!-- 基本信息 -->
          <div class="form-section">
            <h4>基本信息</h4>
            <div class="form-grid">
              <div class="form-group">
                <label>用户名 *</label>
                <input
                  v-model="userForm.username"
                  type="text"
                  required
                  :disabled="editingUser"
                  placeholder="输入用户名"
                />
              </div>
              <div class="form-group">
                <label>全名</label>
                <input
                  v-model="userForm.fullName"
                  type="text"
                  placeholder="用户全名"
                />
              </div>
              <div class="form-group">
                <label>电子邮件</label>
                <input
                  v-model="userForm.email"
                  type="email"
                  placeholder="user@example.com"
                />
              </div>
              <div class="form-group">
                <label>用户角色 *</label>
                <select v-model="userForm.role" required>
                  <option value="guest">访客</option>
                  <option value="user">普通用户</option>
                  <option value="admin">管理员</option>
                  <option value="superadmin">超级管理员</option>
                </select>
              </div>
              <div class="form-group" v-if="!editingUser">
                <label>密码 *</label>
                <input
                  v-model="userForm.password"
                  type="password"
                  required
                  placeholder="输入密码"
                />
              </div>
              <div class="form-group">
                <label>主目录</label>
                <input
                  v-model="userForm.homeDirectory"
                  type="text"
                  placeholder="/home/username"
                />
              </div>
            </div>
          </div>

          <!-- 用户组 -->
          <div class="form-section">
            <h4>用户组成员资格</h4>
            <div class="groups-selection">
              <label
                v-for="group in groups"
                :key="group.id"
                class="group-checkbox"
              >
                <input
                  type="checkbox"
                  v-model="userForm.groups"
                  :value="group.name"
                />
                <span>{{ group.name }}</span>
                <small>{{ group.description || '' }}</small>
              </label>
            </div>
          </div>

          <!-- 存储配额 -->
          <div class="form-section">
            <h4>存储配额</h4>
            <div class="quota-settings">
              <label class="checkbox-label">
                <input type="checkbox" v-model="userForm.quotaEnabled" />
                启用存储配额限制
              </label>
              <div v-if="userForm.quotaEnabled" class="quota-input">
                <input
                  v-model.number="userForm.quotaSize"
                  type="number"
                  placeholder="输入大小"
                  min="0"
                />
                <select v-model="userForm.quotaUnit">
                  <option value="1073741824">GB</option>
                  <option value="1048576">MB</option>
                  <option value="1">TB</option>
                </select>
              </div>
            </div>
          </div>

          <!-- 用户状态 -->
          <div class="form-section">
            <h4>用户状态</h4>
            <div class="status-settings">
              <label class="radio-label">
                <input type="radio" v-model="userForm.status" value="active" />
                <span class="status-radio active">
                  <CheckIcon class="w-4 h-4" />
                  启用
                </span>
              </label>
              <label class="radio-label">
                <input type="radio" v-model="userForm.status" value="disabled" />
                <span class="status-radio disabled">
                  <NoSymbolIcon class="w-4 h-4" />
                  禁用
                </span>
              </label>
            </div>
          </div>

          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeUserModal">
              取消
            </button>
            <button type="submit" class="btn btn-primary">
              {{ editingUser ? '保存更改' : '创建用户' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 权限管理对话框 -->
    <div v-if="showPermissionsModal" class="modal-overlay" @click.self="closePermissionsModal">
      <div class="modal-content extra-large-modal">
        <div class="modal-header">
          <h3>权限管理 - {{ permissionUser?.username }}</h3>
          <button class="close-btn" @click="closePermissionsModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="permissions-content">
          <div class="permission-summary">
            <div class="summary-card">
              <ShieldCheckIcon class="w-8 h-8" />
              <div>
                <h4>当前角色</h4>
                <span class="role-badge large">{{ getRoleLabel(permissionUser?.role) }}</span>
              </div>
            </div>
            <div class="summary-card">
              <UserGroupIcon class="w-8 h-8" />
              <div>
                <h4>所属用户组</h4>
                <div class="groups-list">
                  <span
                    v-for="group in permissionUser?.groups"
                    :key="group"
                    class="group-tag"
                  >
                    {{ getGroupName(group) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="permissions-grid">
            <div
              v-for="(category, categoryKey) in permissionUser?.permissions"
              :key="categoryKey"
              class="permission-category"
            >
              <h4>{{ getCategoryLabel(categoryKey) }}</h4>
              <div class="permission-items">
                <div
                  v-for="(level, permKey) in category"
                  :key="permKey"
                  class="permission-item"
                >
                  <label class="permission-label">{{ getPermissionLabel(permKey) }}</label>
                  <select v-model="category[permKey]" class="permission-select">
                    <option value="none">无权限</option>
                    <option value="read">只读</option>
                    <option value="write">读写</option>
                    <option value="admin">管理</option>
                    <option value="owner">所有者</option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-actions">
          <button type="button" class="btn btn-secondary" @click="resetPermissions">
            重置为角色默认权限
          </button>
          <button type="button" class="btn btn-primary" @click="savePermissions">
            保存权限
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  MagnifyingGlassIcon,
  UserPlusIcon,
  PencilIcon,
  ShieldCheckIcon,
  NoSymbolIcon,
  CheckIcon,
  TrashIcon,
  XMarkIcon,
  UserGroupIcon
} from '@heroicons/vue/24/outline'
import type { User, Group, UserRole } from '@/types/permissions'
import { UserRoles, PermissionTemplates } from '@/types/permissions'

const props = defineProps<{
  users: User[]
  groups: Group[]
}>()

const emit = defineEmits<{
  'user-added': [user: User]
  'user-edited': [user: User]
  'user-deleted': [userId: string]
}>()

// 状态
const searchQuery = ref('')
const roleFilter = ref('')
const selectAll = ref(false)
const selectedUsers = ref<string[]>([])
const showCreateUserModal = ref(false)
const showEditUserModal = ref(false)
const showPermissionsModal = ref(false)
const editingUser = ref<User | null>(null)
const permissionUser = ref<User | null>(null)

// 用户表单
const userForm = ref({
  username: '',
  fullName: '',
  email: '',
  role: 'user' as UserRole,
  password: '',
  groups: [] as string[],
  homeDirectory: '',
  quotaEnabled: false,
  quotaSize: 10,
  quotaUnit: 1073741824,
  status: 'active' as 'active' | 'disabled'
})

// 计算属性
const filteredUsers = computed(() => {
  return props.users.filter(user => {
    const matchesSearch =
      !searchQuery.value ||
      user.username.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (user.fullName && user.fullName.toLowerCase().includes(searchQuery.value.toLowerCase())) ||
      (user.email && user.email.toLowerCase().includes(searchQuery.value.toLowerCase()))

    const matchesRole = !roleFilter.value || user.role === roleFilter.value

    return matchesSearch && matchesRole
  })
})

// 方法
const getRoleLabel = (role?: UserRole) => {
  return UserRoles[role || 'user']?.label || '未知'
}

const getStatusLabel = (status: string) => {
  const labels = {
    active: '启用',
    disabled: '已禁用',
    locked: '已锁定'
  }
  return labels[status as keyof typeof labels] || status
}

const getGroupName = (name: string) => {
  const group = props.groups.find(g => g.name === name)
  return group?.description || name
}

const formatDate = (date?: Date) => {
  if (!date) return '从未登录'
  return new Date(date).toLocaleDateString('zh-CN')
}

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const getQuotaPercent = (user: User) => {
  if (!user.storageQuota?.enabled) return 0
  return (user.storageQuota.used / user.storageQuota.size) * 100
}

const getTotalStorageUsed = () => {
  return props.users
    .filter(u => u.storageQuota?.enabled)
    .reduce((sum, u) => sum + (u.storageQuota?.used || 0), 0)
}

const getCategoryLabel = (key: string) => {
  const labels = {
    system: '系统管理',
    storage: '存储管理',
    file: '文件操作',
    network: '网络管理',
    user: '用户管理',
    backup: '备份恢复',
    app: '应用管理',
    log: '日志查看',
    monitor: '系统监控'
  }
  return labels[key as keyof typeof labels] || key
}

const getPermissionLabel = (key: string) => {
  const labels = {
    settings: '系统设置',
    reboot: '系统重启',
    shutdown: '系统关机',
    view: '查看',
    manage: '管理',
    format: '格式化',
    read: '读取',
    write: '写入',
    delete: '删除',
    share: '共享',
    configure: '配置',
    firewall: '防火墙',
    permissions: '权限管理',
    create: '创建',
    restore: '恢复',
    schedule: '计划',
    install: '安装',
    uninstall: '卸载',
    export: '导出'
  }
  return labels[key as keyof typeof labels] || key
}

// 用户操作
const editUser = (user: User) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    fullName: user.fullName || '',
    email: user.email || '',
    role: user.role,
    password: '',
    groups: [...user.groups],
    homeDirectory: user.homeDirectory,
    quotaEnabled: user.storageQuota?.enabled || false,
    quotaSize: user.storageQuota?.size ? user.storageQuota.size / 1073741824 : 10,
    quotaUnit: 1073741824,
    status: user.status
  }
  showEditUserModal.value = true
}

const managePermissions = (user: User) => {
  permissionUser.value = { ...user }
  showPermissionsModal.value = true
}

const toggleUserStatus = (user: User) => {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  emit('user-edited', { ...user, status: newStatus as 'active' | 'disabled' })
}

const deleteUser = (user: User) => {
  if (confirm(`确定要删除用户 "${user.username}" 吗？此操作不可撤销。`)) {
    emit('user-deleted', user.id)
  }
}

const saveUser = () => {
  const user: User = {
    id: editingUser.value?.id || Date.now().toString(),
    username: userForm.value.username,
    fullName: userForm.value.fullName,
    email: userForm.value.email,
    uid: editingUser.value?.uid || Date.now(),
    primaryGroup: userForm.value.groups[0] || 'users',
    groups: userForm.value.groups,
    role: userForm.value.role,
    status: userForm.value.status,
    homeDirectory: userForm.value.homeDirectory || `/home/${userForm.value.username}`,
    shell: '/bin/bash',
    createdAt: editingUser.value?.createdAt || new Date(),
    storageQuota: userForm.value.quotaEnabled ? {
      enabled: true,
      size: userForm.value.quotaSize * userForm.value.quotaUnit,
      used: 0
    } : undefined,
    permissions: PermissionTemplates[userForm.value.role]
  }

  if (editingUser.value) {
    emit('user-edited', { ...user, id: editingUser.value.id })
  } else {
    emit('user-added', user)
  }

  closeUserModal()
}

const closeUserModal = () => {
  showCreateUserModal.value = false
  showEditUserModal.value = false
  editingUser.value = null
  userForm.value = {
    username: '',
    fullName: '',
    email: '',
    role: 'user',
    password: '',
    groups: [],
    homeDirectory: '',
    quotaEnabled: false,
    quotaSize: 10,
    quotaUnit: 1073741824,
    status: 'active'
  }
}

const closePermissionsModal = () => {
  showPermissionsModal.value = false
  permissionUser.value = null
}

const savePermissions = () => {
  if (permissionUser.value) {
    emit('user-edited', permissionUser.value)
  }
  closePermissionsModal()
}

const resetPermissions = () => {
  if (permissionUser.value) {
    permissionUser.value.permissions = { ...PermissionTemplates[permissionUser.value.role] }
  }
}

// 批量操作
const batchEnable = () => {
  selectedUsers.value.forEach(userId => {
    const user = props.users.find(u => u.id === userId)
    if (user) {
      emit('user-edited', { ...user, status: 'active' })
    }
  })
  selectedUsers.value = []
}

const batchDisable = () => {
  if (confirm(`确定要禁用选中的 ${selectedUsers.value.length} 个用户吗？`)) {
    selectedUsers.value.forEach(userId => {
      const user = props.users.find(u => u.id === userId)
      if (user) {
        emit('user-edited', { ...user, status: 'disabled' })
      }
    })
    selectedUsers.value = []
  }
}

const batchDelete = () => {
  if (confirm(`确定要删除选中的 ${selectedUsers.value.length} 个用户吗？此操作不可撤销。`)) {
    selectedUsers.value.forEach(userId => {
      emit('user-deleted', userId)
    })
    selectedUsers.value = []
  }
}
</script>

<style scoped>
.user-management-panel {
  width: 100%;
}

.panel-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.toolbar-left h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
}

.toolbar-right {
  display: flex;
  gap: 12px;
  align-items: center;
  flex-wrap: wrap;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  width: 250px;
}

.search-box input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 14px;
}

.filter-dropdown select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
}

.btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: white;
  color: #374151;
  border: 1px solid #e5e7eb;
}

.btn-secondary:hover {
  background: #f9fafb;
}

.btn-ghost {
  background: transparent;
  color: #6b7280;
}

.btn-ghost:hover {
  background: #f3f4f6;
}

.stats-bar {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  opacity: 0.9;
}

.users-table-container {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.users-table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.users-table th {
  background: #f9fafb;
  font-size: 13px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
}

.users-table tr:hover {
  background: #f9fafb;
}

.row-disabled {
  opacity: 0.6;
}

.user-cell {
  display: flex;
  gap: 12px;
  align-items: center;
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  position: relative;
}

.admin-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  font-size: 10px;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.user-name {
  font-weight: 600;
  color: #1f2937;
}

.user-email {
  font-size: 12px;
  color: #6b7280;
}

.role-badge {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.role-superadmin {
  background: #fee2e2;
  color: #dc2626;
}

.role-admin {
  background: #fed7aa;
  color: #ea580c;
}

.role-user {
  background: #dbeafe;
  color: #3b82f6;
}

.role-guest {
  background: #f3f4f6;
  color: #6b7280;
}

.group-tag {
  display: inline-block;
  padding: 2px 6px;
  background: #f3f4f6;
  border-radius: 4px;
  font-size: 11px;
  color: #6b7280;
  margin-right: 4px;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
}

.status-active {
  background: #d1fae5;
  color: #10b981;
}

.status-disabled {
  background: #f3f4f6;
  color: #6b7280;
}

.status-locked {
  background: #fee2e2;
  color: #ef4444;
}

.quota-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.quota-bar {
  width: 100px;
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
}

.quota-fill {
  height: 100%;
  background: linear-gradient(90deg, #10b981 0%, #3b82f6 100%);
  transition: width 0.3s;
}

.quota-text {
  font-size: 11px;
  color: #6b7280;
}

.text-muted {
  color: #9ca3af;
  font-size: 13px;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.icon-btn.danger:hover {
  background: #fee2e2;
  color: #ef4444;
}

.batch-actions {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: white;
  border-radius: 12px;
  padding: 16px 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  display: flex;
  gap: 16px;
  align-items: center;
  z-index: 100;
}

.batch-info {
  font-weight: 500;
  color: #1f2937;
}

.batch-buttons {
  display: flex;
  gap: 8px;
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
  max-width: 600px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.large-modal {
  max-width: 800px;
}

.extra-large-modal {
  max-width: 1000px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e5e7eb;
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

.user-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-section h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
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

.form-group input,
.form-group select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.groups-selection {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.group-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
}

.group-checkbox small {
  color: #6b7280;
  font-size: 12px;
}

.quota-settings,
.status-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

.quota-input {
  display: flex;
  gap: 8px;
}

.quota-input input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.status-radio {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.status-radio.active {
  background: #d1fae5;
  color: #10b981;
}

.status-radio.disabled {
  background: #f3f4f6;
  color: #6b7280;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
}

.permissions-content {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.permission-summary {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.summary-card {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
  color: #3b82f6;
}

.summary-card h4 {
  font-size: 12px;
  font-weight: 500;
  color: #6b7280;
  margin-bottom: 8px;
}

.permissions-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.permission-category {
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
}

.permission-category h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.permission-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.permission-label {
  font-size: 13px;
  color: #6b7280;
}

.permission-select {
  padding: 6px 8px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 12px;
}

.role-badge.large {
  padding: 6px 12px;
  font-size: 14px;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .permissions-grid {
    grid-template-columns: 1fr;
  }

  .permission-summary {
    grid-template-columns: 1fr;
  }
}
</style>