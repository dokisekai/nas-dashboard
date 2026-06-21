<template>
  <div class="unified-user-manager">
    <div class="manager-header">
      <div class="header-info">
        <h1>👥 统一用户管理</h1>
        <p class="subtitle">一次修改，所有Docker服务同步更新</p>
      </div>

      <div class="header-actions">
        <button
          @click="syncAllUsers"
          :disabled="syncing || loading"
          class="action-btn primary"
        >
          <ArrowPathIcon v-if="syncing" class="w-4 h-4 animate-spin" />
          <ArrowPathIcon v-else class="w-4 h-4" />
          同步所有用户
        </button>
        <button @click="refreshStatus" class="action-btn">
          <MagnifyingGlassIcon class="w-4 h-4" />
          刷新状态
        </button>
      </div>
    </div>

    <!-- 服务状态概览 -->
    <div class="services-overview">
      <h3>🔗 连接的服务</h3>
      <div class="services-grid">
        <div
          v-for="service in services"
          :key="service.name"
          :class="['service-card', service.status]"
        >
          <div class="service-icon">
            <ServerIcon class="w-6 h-6" />
          </div>
          <div class="service-info">
            <div class="service-name">{{ service.name }}</div>
            <div class="service-status">{{ getStatusText(service.status) }}</div>
            <div class="service-users">{{ service.userCount }} 个用户</div>
          </div>
          <div class="service-actions">
            <button
              @click="testService(service.name)"
              class="test-btn"
              title="测试连接"
            >
              <BoltIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 用户列表 -->
    <div class="users-section">
      <div class="section-header">
        <h2>👤 用户列表 ({{ users.length }})</h2>
        <div class="header-actions">
          <div class="search-box">
            <MagnifyingGlassIcon class="w-4 h-4 search-icon" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索用户名或邮箱..."
              class="search-input"
            />
          </div>
          <button @click="showCreateModal = true" class="action-btn success">
            <UserPlusIcon class="w-4 h-4" />
            添加用户
          </button>
        </div>
      </div>

      <div class="users-table">
        <table class="data-table">
          <thead>
            <tr>
              <th>
                <input
                  type="checkbox"
                  v-model="selectAll"
                  @change="toggleSelectAll"
                />
              </th>
              <th>用户名</th>
              <th>邮箱</th>
              <th>姓名</th>
              <th>角色</th>
              <th>状态</th>
              <th>服务状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="user in filteredUsers"
              :key="user.username"
              :class="{ selected: selectedUsers.includes(user.username) }"
            >
              <td>
                <input
                  type="checkbox"
                  :value="user.username"
                  v-model="selectedUsers"
                />
              </td>
              <td>
                <div class="user-info">
                  <UserIcon class="w-4 h-4" />
                  <span>{{ user.username }}</span>
                </div>
              </td>
              <td>{{ user.email }}</td>
              <td>{{ user.name }}</td>
              <td>
                <span :class="['role-badge', user.role]">
                  {{ getRoleText(user.role) }}
                </span>
              </td>
              <td>
                <span :class="['status-badge', user.isActive ? 'active' : 'inactive']">
                  {{ user.isActive ? '活跃' : '停用' }}
                </span>
              </td>
              <td>
                <div class="sync-status">
                  <template v-if="user.syncStatus">
                    <CheckCircleIcon
                      v-if="user.syncStatus.overall === 'success'"
                      class="w-4 h-4 text-green-500"
                    />
                    <ExclamationCircleIcon
                      v-else-if="user.syncStatus.overall === 'failed'"
                      class="w-4 h-4 text-red-500"
                    />
                    <QuestionMarkCircleIcon
                      v-else
                      class="w-4 h-4 text-yellow-500"
                    />
                    {{ user.syncStatus.successful }}/{{ Object.keys(user.syncStatus.details || {}).length }}
                  </template>
                  <button
                    v-else
                    @click="syncUser(user.username)"
                    class="sync-btn"
                    title="同步此用户"
                  >
                    <CloudArrowUpIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
              <td>
                <div class="action-buttons">
                  <button
                    @click="editUser(user)"
                    class="action-btn edit"
                    title="编辑"
                  >
                    <PencilIcon class="w-3 h-3" />
                  </button>
                  <button
                    @click="syncUser(user.username)"
                    class="action-btn sync"
                    title="同步"
                  >
                    <CloudArrowUpIcon class="w-3 h-3" />
                  </button>
                  <button
                    @click="confirmDeleteUser(user)"
                    class="action-btn delete"
                    title="删除"
                  >
                    <TrashIcon class="w-3 h-3" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="filteredUsers.length === 0" class="empty-state">
        <UserIcon class="w-12 h-12" />
        <p>没有找到用户</p>
      </div>
    </div>

    <!-- 批量操作 -->
    <div v-if="selectedUsers.length > 0" class="batch-operations">
      <div class="batch-info">
        已选择 {{ selectedUsers.length }} 个用户
      </div>
      <div class="batch-actions">
        <button @click="batchSync" class="action-btn">
          <CloudArrowUpIcon class="w-4 h-4" />
          批量同步
        </button>
        <button @click="batchDelete" class="action-btn danger">
          <TrashIcon class="w-4 h-4" />
          批量删除
        </button>
        <button @click="selectedUsers = []" class="action-btn">
          <XMarkIcon class="w-4 h-4" />
          取消选择
        </button>
      </div>
    </div>

    <!-- 创建/编辑用户模态框 -->
    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditModal ? '编辑用户' : '添加用户' }}</h3>
          <button @click="closeModal" class="close-btn">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveUser" class="user-form">
          <div class="form-group">
            <label>用户名 *</label>
            <input
              v-model="userForm.username"
              type="text"
              required
              :disabled="showEditModal"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label>邮箱 *</label>
            <input
              v-model="userForm.email"
              type="email"
              required
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label>姓名</label>
            <input
              v-model="userForm.name"
              type="text"
              class="form-input"
            />
          </div>

          <div class="form-group">
            <label>密码 *</label>
            <input
              v-model="userForm.password"
              type="password"
              :required="!showEditModal"
              class="form-input"
            />
            <small v-if="showEditModal" class="form-hint">留空则不修改密码</small>
          </div>

          <div class="form-group">
            <label>角色</label>
            <select v-model="userForm.role" class="form-select">
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
              <option value="moderator">版主</option>
            </select>
          </div>

          <div class="form-group">
            <label>用户组</label>
            <div class="checkbox-group">
              <label v-for="group in availableGroups" :key="group" class="checkbox-label">
                <input
                  type="checkbox"
                  :value="group"
                  v-model="userForm.groups"
                />
                {{ group }}
              </label>
            </div>
          </div>

          <div class="form-group">
            <label>状态</label>
            <div class="radio-group">
              <label class="radio-label">
                <input type="radio" value="true" v-model="userForm.isActive" />
                活跃
              </label>
              <label class="radio-label">
                <input type="radio" value="false" v-model="userForm.isActive" />
                停用
              </label>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" @click="closeModal" class="action-btn">
              取消
            </button>
            <button type="submit" class="action-btn primary" :disabled="saving">
              <CloudArrowUpIcon v-if="saving" class="w-4 h-4 animate-spin" />
              {{ saving ? '保存中...' : (showEditModal ? '更新并同步' : '创建并同步') }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 同步进度对话框 -->
    <div v-if="showSyncProgress" class="modal-overlay">
      <div class="modal-content sync-progress">
        <div class="modal-header">
          <h3>用户同步进度</h3>
        </div>

        <div class="sync-progress-content">
          <div class="progress-overview">
            <div class="progress-stats">
              <div class="stat-item">
                <div class="stat-value">{{ syncProgress.total }}</div>
                <div class="stat-label">总用户数</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ syncProgress.synced }}</div>
                <div class="stat-label">已同步</div>
              </div>
              <div class="stat-item">
                <div class="stat-value">{{ syncProgress.failed }}</div>
                <div class="stat-label">失败</div>
              </div>
            </div>
          </div>

          <div class="sync-details">
            <div
              v-for="(result, username) in syncProgress.details"
              :key="username"
              :class="['sync-item', result.overall]"
            >
              <div class="sync-item-header">
                <UserIcon class="w-4 h-4" />
                <span>{{ username }}</span>
                <span :class="['sync-status', result.overall]">
                  {{ getStatusText(result.overall) }}
                </span>
              </div>
              <div class="sync-item-details">
                <div
                  v-for="(serviceResult, serviceName) in result.details"
                  :key="serviceName"
                  :class="['service-result', serviceResult.status]"
                >
                  <span>{{ serviceName }}:</span>
                  <span>{{ serviceResult.status }}</span>
                  <span v-if="serviceResult.error" class="error-message">
                    {{ serviceResult.error }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button
            @click="closeSyncProgress"
            class="action-btn primary"
            :disabled="syncing"
          >
            {{ syncing ? '同步中...' : '关闭' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  UserIcon,
  UserPlusIcon,
  PencilIcon,
  TrashIcon,
  CloudArrowUpIcon,
  MagnifyingGlassIcon,
  XMarkIcon,
  CheckCircleIcon,
  ExclamationCircleIcon,
  QuestionMarkCircleIcon,
  ServerIcon,
  BoltIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline'
import { useApiRequest } from '@/composables/useApiRequest'

// 状态定义
const users = ref<any[]>([])
const services = ref<any[]>([])
const selectedUsers = ref<string[]>([])
const selectAll = ref(false)
const searchQuery = ref('')
const syncing = ref(false)
const loading = ref(false)
const saving = ref(false)

// 模态框状态
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showSyncProgress = ref(false)

// 表单数据
const userForm = ref({
  username: '',
  email: '',
  name: '',
  password: '',
  role: 'user',
  groups: [] as string[],
  isActive: true
})

// 同步进度
const syncProgress = ref({
  total: 0,
  synced: 0,
  failed: 0,
  details: {} as any
})

// 可用用户组
const availableGroups = ref(['docker', 'media', 'storage', 'network'])

// API请求
const { data: usersData, execute: fetchUsers } = useApiRequest(
  () => fetch('/api/unified-users').then(r => r.json()),
  { showErrorNotification: true }
)

const { data: servicesData, execute: fetchServices } = useApiRequest(
  () => fetch('/api/unified-users/services').then(r => r.json()),
  { showErrorNotification: true }
)

// 计算属性
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value

  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user =>
    user.username.toLowerCase().includes(query) ||
    user.email.toLowerCase().includes(query) ||
    user.name.toLowerCase().includes(query)
  )
})

// 方法
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/unified-users')
    const data = await response.json()
    users.value = data.users || []
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loading.value = false
  }
}

const loadServices = async () => {
  try {
    const response = await fetch('/api/unified-users/status')
    const data = await response.json()
    services.value = data.services || []
  } catch (error) {
    console.error('Failed to load services:', error)
  }
}

const syncAllUsers = async () => {
  syncing.value = true
  showSyncProgress.value = true

  try {
    const response = await fetch('/api/unified-users/sync', {
      method: 'POST'
    })
    const data = await response.json()

    // 更新同步进度
    syncProgress.value = {
      total: data.totalUsers,
      synced: data.syncedUsers,
      failed: data.failedUsers,
      details: data.userResults || {}
    }

    // 刷新用户列表
    await loadUsers()
  } catch (error) {
    console.error('Failed to sync users:', error)
  } finally {
    syncing.value = false
  }
}

const syncUser = async (username: string) => {
  try {
    const response = await fetch(`/api/unified-users/sync/${username}`, {
      method: 'POST'
    })
    const result = await response.json()

    // 更新用户同步状态
    const userIndex = users.value.findIndex(u => u.username === username)
    if (userIndex !== -1) {
      users.value[userIndex].syncStatus = result
    }
  } catch (error) {
    console.error(`Failed to sync user ${username}:`, error)
  }
}

const batchSync = async () => {
  for (const username of selectedUsers.value) {
    await syncUser(username)
  }
  selectedUsers.value = []
}

const editUser = (user: any) => {
  userForm.value = {
    username: user.username,
    email: user.email,
    name: user.name,
    password: '',
    role: user.role,
    groups: user.groups || [],
    isActive: user.isActive
  }
  showEditModal.value = true
}

const saveUser = async () => {
  saving.value = true

  try {
    const url = showEditModal.value
      ? `/api/unified-users/${userForm.value.username}`
      : '/api/unified-users'

    const method = showEditModal.value ? 'PUT' : 'POST'

    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(userForm.value)
    })

    const result = await response.json()

    // 显示同步结果
    if (result.details) {
      syncProgress.value = {
        total: 1,
        synced: result.status === 'success' ? 1 : 0,
        failed: result.status === 'failed' ? 1 : 0,
        details: { [userForm.value.username]: result }
      }
      showSyncProgress.value = true
    }

    // 刷新用户列表
    await loadUsers()
    closeModal()
  } catch (error) {
    console.error('Failed to save user:', error)
  } finally {
    saving.value = false
  }
}

const confirmDeleteUser = (user: any) => {
  if (confirm(`确定要删除用户 ${user.username} 吗？此操作将在所有服务中删除该用户。`)) {
    deleteUser(user.username)
  }
}

const deleteUser = async (username: string) => {
  try {
    await fetch(`/api/unified-users/${username}`, {
      method: 'DELETE'
    })

    // 刷新用户列表
    await loadUsers()
  } catch (error) {
    console.error(`Failed to delete user ${username}:`, error)
  }
}

const batchDelete = () => {
  if (confirm(`确定要删除选中的 ${selectedUsers.value.length} 个用户吗？此操作将在所有服务中删除这些用户。`)) {
    for (const username of selectedUsers.value) {
      deleteUser(username)
    }
    selectedUsers.value = []
  }
}

const testService = async (serviceName: string) => {
  try {
    await fetch(`/api/unified-users/services/${serviceName}/test`, {
      method: 'POST'
    })
    alert(`服务 ${serviceName} 连接正常`)
  } catch (error) {
    alert(`服务 ${serviceName} 连接失败`)
  }
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  userForm.value = {
    username: '',
    email: '',
    name: '',
    password: '',
    role: 'user',
    groups: [],
    isActive: true
  }
}

const closeSyncProgress = () => {
  showSyncProgress.value = false
}

const refreshStatus = async () => {
  await loadServices()
}

const toggleSelectAll = () => {
  if (selectAll.value) {
    selectedUsers.value = filteredUsers.value.map(u => u.username)
  } else {
    selectedUsers.value = []
  }
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    'healthy': '健康',
    'error': '错误',
    'success': '成功',
    'partial': '部分成功',
    'failed': '失败'
  }
  return statusMap[status] || status
}

const getRoleText = (role: string) => {
  const roleMap: Record<string, string> = {
    'admin': '管理员',
    'moderator': '版主',
    'user': '普通用户'
  }
  return roleMap[role] || role
}

// 生命周期
onMounted(() => {
  loadUsers()
  loadServices()
})
</script>

<style scoped>
.unified-user-manager {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.header-info h1 {
  font-size: 28px;
  margin-bottom: 8px;
}

.subtitle {
  color: #6b7280;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f9fafb;
}

.action-btn.primary {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.success {
  background: #10b981;
  color: white;
  border-color: #10b981;
}

.action-btn.danger {
  background: #ef4444;
  color: white;
  border-color: #ef4444;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.services-overview {
  margin-bottom: 30px;
}

.services-overview h3 {
  font-size: 18px;
  margin-bottom: 16px;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.service-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.service-card.healthy {
  border-color: #10b981;
}

.service-card.error {
  border-color: #ef4444;
}

.service-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-info {
  flex: 1;
}

.service-name {
  font-weight: 600;
  margin-bottom: 4px;
}

.service-status,
.service-users {
  font-size: 12px;
  color: #6b7280;
}

.users-section {
  margin-bottom: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  font-size: 20px;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
}

.search-input {
  border: none;
  outline: none;
  min-width: 200px;
}

.users-table {
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.data-table th {
  background: #f9fafb;
  font-weight: 600;
  font-size: 14px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.role-badge,
.status-badge {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.role-badge.admin {
  background: #fee2e2;
  color: #991b1b;
}

.role-badge.user {
  background: #dbeafe;
  color: #1e40af;
}

.status-badge.active {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.inactive {
  background: #f3f4f6;
  color: #374151;
}

.sync-status {
  display: flex;
  align-items: center;
  gap: 8px;
}

.sync-btn {
  padding: 4px 8px;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  background: white;
  cursor: pointer;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #6b7280;
}

.batch-operations {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
  margin-bottom: 20px;
}

.batch-actions {
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
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.close-btn {
  padding: 4px;
  border: none;
  background: transparent;
  cursor: pointer;
}

.user-form {
  padding: 20px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.form-input,
.form-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.checkbox-group,
.radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.form-hint {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: #6b7280;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}

.sync-progress {
  max-width: 800px;
}

.sync-progress-content {
  padding: 20px;
}

.progress-stats {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.sync-details {
  max-height: 400px;
  overflow-y: auto;
}

.sync-item {
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
}

.sync-item.success {
  border-color: #10b981;
}

.sync-item.failed {
  border-color: #ef4444;
}

.sync-item.partial {
  border-color: #f59e0b;
}

.sync-item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-weight: 500;
}

.sync-item-details {
  padding-left: 32px;
}

.service-result {
  display: flex;
  justify-content: space-between;
  padding: 4px 8px;
  border-radius: 4px;
  margin-bottom: 4px;
}

.service-result.success {
  background: #d1fae5;
}

.service-result.failed {
  background: #fee2e2;
}

.error-message {
  color: #ef4444;
  font-size: 12px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  padding: 20px;
  border-top: 1px solid #e5e7eb;
}

.sync-status {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.sync-status.success {
  background: #d1fae5;
  color: #065f46;
}

.sync-status.failed {
  background: #fee2e2;
  color: #991b1b;
}

.sync-status.partial {
  background: #fef3c7;
  color: #92400e;
}

.data-table tr.selected {
  background: #eff6ff;
}
</style>