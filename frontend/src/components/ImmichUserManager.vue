<template>
  <div class="immich-user-manager">
    <div class="manager-header">
      <h3>Immich 用户管理</h3>
      <div class="manager-actions">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索用户..."
          class="search-input"
        />
        <button @click="showAddUserModal = true" class="btn btn-primary">
          <UserPlusIcon class="w-4 h-4" />
          添加用户
        </button>
        <button @click="syncWithSystemUsers" class="btn btn-secondary" :disabled="syncing">
          <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': syncing }" />
          同步系统用户
        </button>
      </div>
    </div>

    <!-- 批量操作栏 -->
    <div v-if="selectedUsers.length > 0" class="batch-actions">
      <div class="batch-info">已选择 {{ selectedUsers.length }} 个用户</div>
      <div class="batch-buttons">
        <button @click="batchUpdateStatus" class="btn btn-sm">
          批量修改状态
        </button>
        <button @click="batchUpdateRole" class="btn btn-sm">
          批量修改角色
        </button>
        <button @click="clearSelection" class="btn btn-sm btn-ghost">
          取消选择
        </button>
      </div>
    </div>

    <div class="users-table">
      <table>
        <thead>
          <tr>
            <th width="40">
              <input
                type="checkbox"
                v-model="selectAll"
                @change="toggleSelectAll"
              />
            </th>
            <th>邮箱</th>
            <th>姓名</th>
            <th>状态</th>
            <th>角色</th>
            <th>创建时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.id" :class="{ selected: selectedUsers.includes(user.id) }">
            <td>
              <input
                type="checkbox"
                :value="user.id"
                v-model="selectedUsers"
              />
            </td>
            <td>{{ user.email }}</td>
            <td>{{ user.name }}</td>
            <td>
              <span :class="['status-badge', user.isActive ? 'active' : 'inactive']">
                {{ user.isActive ? '正常' : '停用' }}
              </span>
            </td>
            <td>
              <span class="role-badge">{{ getRoleLabel(user.role) }}</span>
            </td>
            <td>{{ formatDate(user.createdAt) }}</td>
            <td>
              <div class="action-buttons">
                <button @click="editUser(user)" class="btn btn-sm btn-secondary" title="编辑">
                  <PencilIcon class="w-3 h-3" />
                </button>
                <button @click="deleteUser(user)" class="btn btn-sm btn-danger" title="删除">
                  <TrashIcon class="w-3 h-3" />
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 添加/编辑用户模态框 -->
    <div v-if="showAddUserModal || showEditUserModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal">
        <div class="modal-header">
          <h4>{{ editingUser ? '编辑用户' : '添加用户' }}</h4>
          <button @click="closeModals" class="btn btn-ghost">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>邮箱 *</label>
            <input
              v-model="userForm.email"
              type="email"
              placeholder="user@example.com"
              :disabled="editingUser"
            />
          </div>
          <div class="form-group" v-if="!editingUser">
            <label>密码 *</label>
            <input
              v-model="userForm.password"
              type="password"
              placeholder="输入密码"
            />
          </div>
          <div class="form-group">
            <label>姓名 *</label>
            <input
              v-model="userForm.name"
              type="text"
              placeholder="用户姓名"
            />
          </div>
          <div class="form-group">
            <label>状态</label>
            <select v-model="userForm.isActive">
              <option :value="true">正常</option>
              <option :value="false">停用</option>
            </select>
          </div>
          <div class="form-group">
            <label>角色</label>
            <select v-model="userForm.role">
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModals" class="btn btn-secondary">取消</button>
          <button @click="saveUser" class="btn btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 批量更新模态框 -->
    <div v-if="showBatchUpdateModal" class="modal-overlay" @click.self="showBatchUpdateModal = false">
      <div class="modal">
        <div class="modal-header">
          <h4>批量更新用户</h4>
          <button @click="showBatchUpdateModal = false" class="btn btn-ghost">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group" v-if="batchUpdateType === 'status'">
            <label>状态</label>
            <select v-model="batchUpdateData.isActive">
              <option :value="true">正常</option>
              <option :value="false">停用</option>
            </select>
          </div>
          <div class="form-group" v-if="batchUpdateType === 'role'">
            <label>角色</label>
            <select v-model="batchUpdateData.role">
              <option value="">保持原角色</option>
              <option value="user">普通用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="batch-info">
            将更新 {{ selectedUsers.length }} 个用户
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showBatchUpdateModal = false" class="btn btn-secondary">取消</button>
          <button @click="confirmBatchUpdate" class="btn btn-primary" :disabled="saving">
            {{ saving ? '更新中...' : '确认更新' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { UserPlusIcon, PencilIcon, TrashIcon, XMarkIcon, ArrowPathIcon } from '@heroicons/vue/24/outline'

// 用户数据
const users = ref<any[]>([])
const searchQuery = ref('')
const showAddUserModal = ref(false)
const showEditUserModal = ref(false)
const showBatchUpdateModal = ref(false)
const editingUser = ref<any>(null)
const saving = ref(false)
const syncing = ref(false)
const selectedUsers = ref<string[]>([])
const selectAll = ref(false)
const batchUpdateType = ref<'status' | 'role'>('status')

// 批量更新数据
const batchUpdateData = ref({
  isActive: true,
  role: ''
})

// 用户表单
const userForm = ref({
  email: '',
  password: '',
  name: '',
  isActive: true,
  role: 'user'
})

// 过滤后的用户列表
const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value

  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user =>
    user.email?.toLowerCase().includes(query) ||
    user.name?.toLowerCase().includes(query)
  )
})

// 获取Immich用户数据
const fetchUsers = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/immich/users', {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      const result = await response.json()
      users.value = result.users || []
    } else {
      console.error('Failed to fetch Immich users:', response.status)
    }
  } catch (error) {
    console.error('Failed to fetch Immich users:', error)
  }
}

// 同步系统用户
const syncWithSystemUsers = async () => {
  if (!confirm('确定要将系统用户同步到Immich吗？这将创建或更新对应的Immich用户。')) return

  syncing.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/immich/users/sync', {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      const result = await response.json()
      alert(`同步完成！\n创建: ${result.created.length}\n更新: ${result.updated.length}\n失败: ${result.unmatched.length}`)
      await fetchUsers()
    } else {
      alert('同步失败')
    }
  } catch (error) {
    console.error('Failed to sync users:', error)
    alert('同步失败')
  } finally {
    syncing.value = false
  }
}

// 全选切换
const toggleSelectAll = () => {
  if (selectAll.value) {
    selectedUsers.value = filteredUsers.value.map(u => u.id)
  } else {
    selectedUsers.value = []
  }
}

// 清除选择
const clearSelection = () => {
  selectedUsers.value = []
  selectAll.value = false
}

// 批量更新状态
const batchUpdateStatus = () => {
  batchUpdateType.value = 'status'
  batchUpdateData.value = { isActive: true, role: '' }
  showBatchUpdateModal.value = true
}

// 批量更新角色
const batchUpdateRole = () => {
  batchUpdateType.value = 'role'
  batchUpdateData.value = { isActive: true, role: '' }
  showBatchUpdateModal.value = true
}

// 确认批量更新
const confirmBatchUpdate = async () => {
  if (selectedUsers.value.length === 0) return

  saving.value = true
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/immich/users/batch', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        userIds: selectedUsers.value,
        updates: batchUpdateType.value === 'status'
          ? { isActive: batchUpdateData.value.isActive }
          : { role: batchUpdateData.value.role || undefined }
      })
    })

    if (response.ok) {
      const result = await response.json()
      alert(`批量更新完成！\n成功: ${result.success}\n失败: ${result.failed}`)
      showBatchUpdateModal.value = false
      clearSelection()
      await fetchUsers()
    } else {
      alert('批量更新失败')
    }
  } catch (error) {
    console.error('Failed to batch update:', error)
    alert('批量更新失败')
  } finally {
    saving.value = false
  }
}

// 编辑用户
const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    email: user.email,
    password: '',
    name: user.name,
    isActive: user.isActive,
    role: user.role || 'user'
  }
  showEditUserModal.value = true
}

// 删除用户
const deleteUser = async (user: any) => {
  if (!confirm(`确定要删除用户 "${user.name}" 吗？`)) return

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/immich/users/${user.id}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      users.value = users.value.filter(u => u.id !== user.id)
      alert('用户删除成功')
    } else {
      alert('删除失败')
    }
  } catch (error) {
    console.error('Failed to delete user:', error)
    alert('删除失败')
  }
}

// 保存用户
const saveUser = async () => {
  if (!userForm.value.email || !userForm.value.name) {
    alert('请填写必填项')
    return
  }

  if (!editingUser.value && !userForm.value.password) {
    alert('请输入密码')
    return
  }

  saving.value = true

  try {
    const token = localStorage.getItem('token')
    const url = editingUser.value
      ? `/api/immich/users/${editingUser.value.id}`
      : '/api/immich/users'

    const method = editingUser.value ? 'PUT' : 'POST'

    const userData = {
      email: userForm.value.email,
      name: userForm.value.name,
      password: userForm.value.password,
      isActive: userForm.value.isActive,
      role: userForm.value.role
    }

    const response = await fetch(url, {
      method,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(userData)
    })

    if (response.ok) {
      alert(editingUser.value ? '用户更新成功' : '用户创建成功')
      closeModals()
      await fetchUsers()
    } else {
      const error = await response.json()
      alert(error.error || '操作失败')
    }
  } catch (error) {
    console.error('Failed to save user:', error)
    alert('操作失败')
  } finally {
    saving.value = false
  }
}

// 关闭模态框
const closeModals = () => {
  showAddUserModal.value = false
  showEditUserModal.value = false
  editingUser.value = null
  userForm.value = {
    email: '',
    password: '',
    name: '',
    isActive: true,
    role: 'user'
  }
}

// 获取角色标签
const getRoleLabel = (role: string) => {
  const roleMap: Record<string, string> = {
    'admin': '管理员',
    'user': '普通用户'
  }
  return roleMap[role] || role
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.immich-user-manager {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.manager-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.manager-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.batch-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f3f4f6;
  border-bottom: 1px solid #e5e7eb;
}

.batch-info {
  font-size: 14px;
  color: #374151;
}

.batch-buttons {
  display: flex;
  gap: 8px;
}

.users-table {
  flex: 1;
  overflow: auto;
  padding: 16px;
}

.users-table table {
  width: 100%;
  border-collapse: collapse;
}

.users-table th,
.users-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.users-table th {
  background: #f9fafb;
  font-weight: 600;
  color: #374151;
  font-size: 14px;
}

.users-table tr.selected {
  background: #eff6ff;
}

.users-table tr:hover {
  background: #f9fafb;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.inactive {
  background: #fee2e2;
  color: #991b1b;
}

.role-badge {
  padding: 4px 8px;
  border-radius: 4px;
  background: #dbeafe;
  color: #1e40af;
  font-size: 12px;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

.btn {
  padding: 8px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  font-size: 14px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn:hover {
  background: #f9fafb;
}

.btn.btn-primary {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.btn.btn-primary:hover {
  background: #2563eb;
}

.btn.btn-secondary {
  background: #6b7280;
  color: white;
  border-color: #6b7280;
}

.btn.btn-secondary:hover {
  background: #4b5563;
}

.btn.btn-danger {
  background: #ef4444;
  color: white;
  border-color: #ef4444;
}

.btn.btn-danger:hover {
  background: #dc2626;
}

.btn.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn.btn-ghost {
  background: transparent;
  border-color: transparent;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 模态框样式 */
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

.modal {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h4 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.modal-body {
  padding: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  font-size: 14px;
  color: #374151;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.batch-info {
  padding: 12px;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 14px;
  color: #374151;
  text-align: center;
}
</style>