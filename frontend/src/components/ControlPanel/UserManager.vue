<template>
  <div class="user-manager">
    <div class="manager-header">
      <h3>用户管理</h3>
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
      </div>
    </div>

    <div class="users-table">
      <table>
        <thead>
          <tr>
            <th>用户名</th>
            <th>UID</th>
            <th>组</th>
            <th>状态</th>
            <th>上次登录</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.username">
            <td>
              <div class="user-info">
                <UserIcon class="w-4 h-4" />
                <span>{{ user.username }}</span>
              </div>
            </td>
            <td>{{ user.uid }}</td>
            <td>
              <span class="group-tags">
                <span v-for="group in user.groups" :key="group" class="group-tag">
                  {{ group }}
                </span>
              </span>
            </td>
            <td>
              <span :class="['status-badge', user.status.toLowerCase()]">
                {{ user.status }}
              </span>
            </td>
            <td>{{ formatDate(user.lastLogin) }}</td>
            <td>
              <div class="action-buttons">
                <button @click="editUser(user)" class="btn btn-sm btn-secondary">
                  <PencilIcon class="w-3 h-3" />
                </button>
                <button @click="deleteUser(user)" class="btn btn-sm btn-danger">
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
            <label>用户名 *</label>
            <input
              v-model="userForm.username"
              type="text"
              placeholder="输入用户名"
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
            <label>主组 *</label>
            <select v-model="userForm.primaryGroup">
              <option value="users">users</option>
              <option v-for="group in availableGroups" :key="group" :value="group">
                {{ group }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>附加组</label>
            <div class="checkbox-group">
              <label v-for="group in availableGroups" :key="group" class="checkbox-label">
                <input
                  type="checkbox"
                  :value="group"
                  v-model="userForm.additionalGroups"
                />
                {{ group }}
              </label>
            </div>
          </div>
          <div class="form-group">
            <label>状态</label>
            <select v-model="userForm.status">
              <option value="active">启用</option>
              <option value="disabled">禁用</option>
              <option value="locked">锁定</option>
            </select>
          </div>
          <div class="form-group">
            <label>备注</label>
            <textarea
              v-model="userForm.comment"
              placeholder="用户备注信息"
              rows="3"
            ></textarea>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { UserIcon, UserPlusIcon, PencilIcon, TrashIcon, XMarkIcon } from '@heroicons/vue/24/outline'

// 用户数据
const users = ref<any[]>([])
const searchQuery = ref('')
const showAddUserModal = ref(false)
const showEditUserModal = ref(false)
const editingUser = ref<any>(null)
const saving = ref(false)

// 可用的组
const availableGroups = ref<string[]>(['users', 'admin', 'sudo', 'docker', 'smb', 'ftp'])

// 用户表单
const userForm = ref({
  username: '',
  password: '',
  primaryGroup: 'users',
  additionalGroups: [] as string[],
  status: 'active',
  comment: ''
})

// 过滤后的用户列表 - 只显示普通用户（UID >= 1000）
const filteredUsers = computed(() => {
  let result = users.value.filter(user => {
    // 只显示普通用户，过滤掉系统用户
    const uid = parseInt(user.uid)
    return uid >= 1000
  })

  if (searchQuery.value) {
    result = result.filter(user =>
      user.username.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  return result
})

// 获取用户数据
const fetchUsers = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/users', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.ok) {
      const result = await response.json()
      // 处理后端返回的格式: { "users": [...] }
      const usersArray = result.users || result || []
      users.value = usersArray.map((user: any) => ({
        ...user,
        uid: user.uid,
        groups: user.groups || [user.group || 'users'],
        status: user.status || 'active',
        lastLogin: user.lastLogin || new Date().toISOString()
      }))
    }
  } catch (error) {
    console.error('Failed to fetch users:', error)
    // 使用模拟数据
    users.value = [
      {
        username: 'admin',
        uid: '1000',
        groups: ['admin', 'sudo'],
        status: 'active',
        lastLogin: new Date().toISOString()
      },
      {
        username: 'user1',
        uid: '1001',
        groups: ['users'],
        status: 'active',
        lastLogin: new Date(Date.now() - 86400000).toISOString()
      }
    ]
  }
}

// 编辑用户
const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    password: '',
    primaryGroup: user.groups[0] || 'users',
    additionalGroups: user.groups.slice(1) || [],
    status: user.status,
    comment: user.comment || ''
  }
  showEditUserModal.value = true
}

// 删除用户
const deleteUser = async (user: any) => {
  if (!confirm(`确定要删除用户 "${user.username}" 吗？`)) return

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/users/${user.username}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      users.value = users.value.filter(u => u.username !== user.username)
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
  if (!userForm.value.username) {
    alert('请输入用户名')
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
      ? `/api/users/${userForm.value.username}`
      : '/api/users'

    const method = editingUser.value ? 'PUT' : 'POST'

    const userData = {
      username: userForm.value.username,
      password: userForm.value.password,
      primaryGroup: userForm.value.primaryGroup,
      groups: [userForm.value.primaryGroup, ...userForm.value.additionalGroups],
      status: userForm.value.status,
      comment: userForm.value.comment
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
      await fetchUsers()
      closeModals()
      alert(editingUser.value ? '用户更新成功' : '用户创建成功')
    } else {
      alert('操作失败')
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
    username: '',
    password: '',
    primaryGroup: 'users',
    additionalGroups: [] as string[],
    status: 'active',
    comment: ''
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  if (!dateString) return '从未登录'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.user-manager {
  width: 100%;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.manager-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.manager-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.search-input {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  width: 250px;
}

.btn {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #6b7280;
  color: white;
}

.btn-secondary:hover {
  background: #4b5563;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

.btn-ghost {
  background: transparent;
  color: #6b7280;
}

.btn-ghost:hover {
  background: #f3f4f6;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.users-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.users-table table {
  width: 100%;
  border-collapse: collapse;
}

.users-table thead {
  background: #f9fafb;
}

.users-table th {
  padding: 12px 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
  border-bottom: 1px solid #e5e7eb;
}

.users-table td {
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  font-size: 14px;
  color: #1f2937;
}

.users-table tbody tr:hover {
  background: #f9fafb;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.group-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.group-tag {
  padding: 2px 8px;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
  text-transform: uppercase;
}

.status-badge.active {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.disabled {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.locked {
  background: #fee2e2;
  color: #991b1b;
}

.action-buttons {
  display: flex;
  gap: 4px;
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

.modal {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h4 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
}

.form-group input:disabled {
  background: #f3f4f6;
  cursor: not-allowed;
}

.checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>