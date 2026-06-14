<template>
  <div class="simple-user-panel">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <h3>用户管理</h3>
        <p class="subtitle">管理您的NAS用户</p>
      </div>
      <div class="toolbar-right">
        <div class="search-box">
          <MagnifyingGlassIcon class="w-4 h-4" />
          <input type="text" placeholder="搜索用户..." v-model="searchQuery" />
        </div>
        <button class="btn btn-primary" @click="showAddUserModal = true">
          <UserPlusIcon class="w-4 h-4" />
          添加用户
        </button>
      </div>
    </div>

    <!-- 用户统计 -->
    <div class="stats-cards">
      <div class="stat-card">
        <div class="stat-icon blue">
          <UsersIcon class="w-6 h-6" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ users.length }}</div>
          <div class="stat-label">总用户数</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon green">
          <CheckCircleIcon class="w-6 h-6" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ activeUsersCount }}</div>
          <div class="stat-label">活跃用户</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon orange">
          <ShieldCheckIcon class="w-6 h-6" />
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ adminCount }}</div>
          <div class="stat-label">管理员</div>
        </div>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="users-list">
      <div
        v-for="user in filteredUsers"
        :key="user.id"
        class="user-card"
      >
        <div class="user-avatar">
          {{ user.name.charAt(0).toUpperCase() }}
          <span v-if="user.role === 'admin'" class="crown">👑</span>
        </div>

        <div class="user-details">
          <div class="user-name-row">
            <h4>{{ user.displayName || user.name }}</h4>
            <span :class="['role-badge', `role-${user.role}`]">
              {{ getRoleLabel(user.role) }}
            </span>
          </div>
          <div class="user-meta">
            <span>用户名: {{ user.name }}</span>
            <span>•</span>
            <span>用户组: {{ user.groups.join(', ') }}</span>
          </div>
          <div class="user-storage">
            <DatabaseIcon class="w-4 h-4" />
            <span>已用 {{ user.storageUsed }}</span>
            <span v-if="user.storageLimit !== '无限制'">
              / {{ user.storageLimit }}
            </span>
          </div>
          <div class="user-status">
            <ClockIcon class="w-4 h-4" />
            <span>上次登录: {{ user.lastLogin }}</span>
          </div>
        </div>

        <div class="user-actions">
          <button class="icon-btn" @click="editUser(user)" title="编辑">
            <PencilIcon class="w-4 h-4" />
          </button>
          <button
            class="icon-btn"
            @click="toggleUserStatus(user)"
            :title="user.status === 'active' ? '禁用' : '启用'"
          >
            <NoSymbolIcon v-if="user.status === 'active'" class="w-4 h-4" />
            <CheckIcon v-else class="w-4 h-4" />
          </button>
          <button class="icon-btn danger" @click="deleteUser(user)" title="删除">
            <TrashIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- 添加/编辑用户对话框 -->
    <div v-if="showAddUserModal" class="modal-overlay" @click.self="closeUserModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingUser ? '编辑用户' : '添加用户' }}</h3>
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
                  v-model="userForm.name"
                  type="text"
                  required
                  :disabled="editingUser"
                  placeholder="例如: john"
                />
                <small>用于登录的用户名</small>
              </div>

              <div class="form-group">
                <label>显示名称 *</label>
                <input
                  v-model="userForm.displayName"
                  type="text"
                  required
                  placeholder="例如: 张三"
                />
                <small>显示在系统中的名称</small>
              </div>

              <div class="form-group" v-if="!editingUser">
                <label>密码 *</label>
                <input
                  v-model="userForm.password"
                  type="password"
                  required
                  placeholder="设置密码"
                />
                <small>至少6个字符</small>
              </div>

              <div class="form-group">
                <label>用户角色 *</label>
                <div class="role-options">
                  <label
                    v-for="role in roles"
                    :key="role.value"
                    :class="['role-option', { selected: userForm.role === role.value }]"
                  >
                    <input type="radio" v-model="userForm.role" :value="role.value" />
                    <div class="role-card">
                      <div class="role-icon">{{ role.icon }}</div>
                      <div class="role-info">
                        <span class="role-name">{{ role.label }}</span>
                        <span class="role-desc">{{ role.description }}</span>
                      </div>
                    </div>
                  </label>
                </div>
              </div>
            </div>
          </div>

          <!-- 用户组 -->
          <div class="form-section">
            <h4>用户组成员资格</h4>
            <div class="groups-selection">
              <label
                v-for="group in availableGroups"
                :key="group.id"
                class="group-checkbox"
              >
                <input type="checkbox" v-model="userForm.groups" :value="group.name" />
                <div class="group-info">
                  <span class="group-name">{{ group.name }}</span>
                  <span class="group-desc">{{ group.description }}</span>
                </div>
              </label>
            </div>
          </div>

          <!-- 存储空间 -->
          <div class="form-section">
            <h4>存储空间限制</h4>
            <div class="storage-settings">
              <label class="toggle-label">
                <input type="checkbox" v-model="userForm.limitStorage" />
                <span>限制存储空间使用</span>
              </label>

              <div v-if="userForm.limitStorage" class="storage-input">
                <input
                  v-model.number="userForm.storageLimit"
                  type="number"
                  min="1"
                  placeholder="100"
                />
                <select v-model="userForm.storageUnit">
                  <option value="GB">GB</option>
                  <option value="TB">TB</option>
                </select>
              </div>
            </div>
          </div>

          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeUserModal">
              取消
            </button>
            <button type="submit" class="btn btn-primary">
              {{ editingUser ? '保存更改' : '添加用户' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  MagnifyingGlassIcon,
  UserPlusIcon,
  UsersIcon,
  CheckCircleIcon,
  ShieldCheckIcon,
  PencilIcon,
  NoSymbolIcon,
  CheckIcon,
  TrashIcon,
  XMarkIcon,
  DatabaseIcon,
  ClockIcon
} from '@heroicons/vue/24/outline'

const props = defineProps<{
  users: any[]
}>()

const emit = defineEmits<{
  'user-created': [user: any]
  'user-updated': [user: any]
  'user-deleted': [userId: string]
}>()

const searchQuery = ref('')
const showAddUserModal = ref(false)
const editingUser = ref<any>(null)

const userForm = ref({
  name: '',
  displayName: '',
  password: '',
  role: 'user',
  groups: [],
  limitStorage: false,
  storageLimit: 100,
  storageUnit: 'GB'
})

// 简化的角色定义
const roles = [
  {
    value: 'admin',
    label: '管理员',
    icon: '👑',
    description: '可以管理所有功能和设置'
  },
  {
    value: 'user',
    label: '普通用户',
    icon: '👤',
    description: '可以访问文件和个人设置'
  },
  {
    value: 'guest',
    label: '访客',
    icon: '👀',
    description: '只能查看，不能修改文件'
  }
]

// 模拟可用用户组
const availableGroups = ref([
  {
    id: '1',
    name: '管理员组',
    description: '完整系统管理权限'
  },
  {
    id: '2',
    name: '普通用户组',
    description: '文件访问和基本功能'
  },
  {
    id: '3',
    name: '访客组',
    description: '只读访问权限'
  }
])

const filteredUsers = computed(() => {
  if (!searchQuery.value) return props.users

  const query = searchQuery.value.toLowerCase()
  return props.users.filter(user =>
    user.name.toLowerCase().includes(query) ||
    (user.displayName && user.displayName.toLowerCase().includes(query))
  )
})

const activeUsersCount = computed(() => {
  return props.users.filter(u => u.status === 'active').length
})

const adminCount = computed(() => {
  return props.users.filter(u => u.role === 'admin').length
})

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    admin: '管理员',
    user: '普通用户',
    guest: '访客'
  }
  return labels[role] || role
}

const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    name: user.name,
    displayName: user.displayName || '',
    password: '',
    role: user.role,
    groups: [...user.groups],
    limitStorage: user.storageLimit !== '无限制',
    storageLimit: parseInt(user.storageLimit) || 100,
    storageUnit: 'GB'
  }
  showAddUserModal.value = true
}

const toggleUserStatus = (user: any) => {
  const newStatus = user.status === 'active' ? 'disabled' : 'active'
  emit('user-updated', { ...user, status: newStatus })
}

const deleteUser = (user: any) => {
  if (confirm(`确定要删除用户 "${user.displayName || user.name}" 吗？`)) {
    emit('user-deleted', user.id)
  }
}

const saveUser = () => {
  const user = {
    id: editingUser.value?.id || Date.now().toString(),
    name: userForm.value.name,
    displayName: userForm.value.displayName,
    role: userForm.value.role,
    groups: userForm.value.groups,
    storageUsed: '0 GB',
    storageLimit: userForm.value.limitStorage
      ? `${userForm.value.storageLimit} ${userForm.value.storageUnit}`
      : '无限制',
    status: 'active',
    lastLogin: '从未登录'
  }

  if (editingUser.value) {
    emit('user-updated', { ...user, id: editingUser.value.id })
  } else {
    emit('user-created', user)
  }

  closeUserModal()
}

const closeUserModal = () => {
  showAddUserModal.value = false
  editingUser.value = null
  userForm.value = {
    name: '',
    displayName: '',
    password: '',
    role: 'user',
    groups: [],
    limitStorage: false,
    storageLimit: 100,
    storageUnit: 'GB'
  }
}
</script>

<style scoped>
.simple-user-panel {
  width: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.toolbar-left h3 {
  font-size: 24px;
  font-weight: 700;
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
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  width: 300px;
}

.search-box input {
  flex: 1;
  border: none;
  outline: none;
  font-size: 14px;
}

.btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-secondary {
  background: white;
  color: #374151;
  border: 1px solid #e5e7eb;
}

.btn-secondary:hover {
  background: #f9fafb;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.blue {
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
}

.stat-icon.green {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.stat-icon.orange {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.stat-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 13px;
  color: #6b7280;
}

.users-list {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.user-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  gap: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.2s;
}

.user-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 600;
  position: relative;
  flex-shrink: 0;
}

.crown {
  position: absolute;
  top: -4px;
  right: -4px;
  font-size: 14px;
}

.user-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.user-name-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-name-row h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.role-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.role-badge.role-admin {
  background: #fed7aa;
  color: #c2410c;
}

.role-badge.role-user {
  background: #dbeafe;
  color: #1e40af;
}

.role-badge.role-guest {
  background: #f3f4f6;
  color: #4b5563;
}

.user-meta {
  font-size: 13px;
  color: #6b7280;
  display: flex;
  gap: 8px;
  align-items: center;
}

.user-storage,
.user-status {
  font-size: 13px;
  color: #6b7280;
  display: flex;
  align-items: center;
  gap: 6px;
}

.user-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
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
  max-height: 85vh;
  overflow-y: auto;
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

.form-section h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
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
  padding: 12px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.form-group small {
  font-size: 12px;
  color: #6b7280;
}

.role-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.role-option {
  cursor: pointer;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  padding: 12px;
  transition: all 0.2s;
}

.role-option.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.role-option:hover:not(.selected) {
  border-color: #d1d5db;
}

.role-card {
  display: flex;
  align-items: center;
  gap: 12px;
}

.role-icon {
  font-size: 24px;
}

.role-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.role-name {
  font-weight: 500;
  color: #1f2937;
}

.role-desc {
  font-size: 12px;
  color: #6b7280;
}

.groups-selection {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.group-checkbox {
  display: flex;
  gap: 12px;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.group-checkbox:hover {
  background: #f9fafb;
  border-color: #d1d5db;
}

.group-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.group-name {
  font-weight: 500;
  color: #1f2937;
  font-size: 14px;
}

.group-desc {
  font-size: 12px;
  color: #6b7280;
}

.storage-settings {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  cursor: pointer;
}

.storage-input {
  display: flex;
  gap: 8px;
}

.storage-input input,
.storage-input select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .user-card {
    flex-direction: column;
  }

  .user-actions {
    flex-direction: row;
    width: 100%;
  }

  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .search-box {
    width: 100%;
  }
}
</style>