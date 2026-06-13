<template>
  <div class="user-manager">
    <div class="manager-header">
      <h1>用户管理器</h1>
      <p class="subtitle">管理系统用户和权限</p>
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

    <!-- Users Tab -->
    <div v-if="activeTab === 'users'" class="tab-content">
      <div class="section-header">
        <h2>用户列表</h2>
        <button class="action-btn primary" @click="showCreateUserModal = true">
          <UserPlusIcon class="w-4 h-4" />
          添加用户
        </button>
      </div>

      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="users-table">
        <table>
          <thead>
            <tr>
              <th>用户名</th>
              <th>全名</th>
              <th>用户组</th>
              <th>Shell</th>
              <th>磁盘配额</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in users" :key="user.username">
              <td>
                <div class="user-info">
                  <div class="user-avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                  <span>{{ user.username }}</span>
                </div>
              </td>
              <td>{{ user.comment || '-' }}</td>
              <td>{{ user.group || 'users' }}</td>
              <td>{{ user.shell || '/bin/bash' }}</td>
              <td>{{ user.quota || '无限制' }}</td>
              <td>
                <span class="status-badge" :class="user.status?.toLowerCase()">
                  {{ user.status || 'Active' }}
                </span>
              </td>
              <td>
                <div class="action-buttons">
                  <button class="icon-btn" @click="editUser(user)" title="编辑">
                    <PencilIcon class="w-4 h-4" />
                  </button>
                  <button class="icon-btn" @click="manageUserQuota(user)" title="配额">
                    <CircleStackIcon class="w-4 h-4" />
                  </button>
                  <button class="icon-btn danger" @click="deleteUser(user)" title="删除">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- SSH Keys Tab -->
    <div v-if="activeTab === 'ssh'" class="tab-content">
      <div class="section-header">
        <h2>SSH 密钥管理</h2>
        <button class="action-btn primary" @click="showAddSSHKeyModal = true">
          <KeyIcon class="w-4 h-4" />
          添加密钥
        </button>
      </div>

      <div v-if="loadingSSH" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>

      <div v-else class="ssh-keys-list">
        <div
          v-for="key in sshKeys"
          :key="key.id"
          class="ssh-key-item"
        >
          <div class="key-info">
            <div class="key-icon">
              <KeyIcon class="w-8 h-8" />
            </div>
            <div class="key-details">
              <h3>{{ key.name }}</h3>
              <p class="key-fingerprint">指纹: {{ key.fingerprint }}</p>
              <p class="key-user">用户: {{ key.user }}</p>
              <p class="key-date">添加时间: {{ formatDate(key.addedAt) }}</p>
            </div>
          </div>

          <div class="key-actions">
            <button class="action-btn" @click="viewKey(key)">
              <EyeIcon class="w-4 h-4" />
              查看
            </button>
            <button class="action-btn danger" @click="deleteSSHKey(key)">
              <TrashIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Groups Tab -->
    <div v-if="activeTab === 'groups'" class="tab-content">
      <GroupManager />
    </div>

    <!-- User Modal -->
    <div v-if="showCreateUserModal" class="modal-overlay" @click="closeUserModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑用户' : '添加用户' }}</h3>
          <button class="close-btn" @click="closeUserModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveUser" class="modal-body">
          <div class="form-group">
            <label>用户名 *</label>
            <input
              v-model="userForm.username"
              type="text"
              required
              :disabled="isEditing"
              placeholder="输入用户名"
            />
          </div>

          <div class="form-group" v-if="!isEditing">
            <label>密码 *</label>
            <div class="password-input">
              <input
                v-model="userForm.password"
                :type="showPassword ? 'text' : 'password'"
                required
                placeholder="输入密码"
              />
              <button type="button" class="icon-btn" @click="showPassword = !showPassword">
                <EyeIcon v-if="!showPassword" class="w-4 h-4" />
                <EyeSlashIcon v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div class="form-group">
            <label>全名</label>
            <input
              v-model="userForm.comment"
              type="text"
              placeholder="可选全名"
            />
          </div>

          <div class="form-group">
            <label>用户组</label>
            <select v-model="userForm.group">
              <option value="">默认组 (users)</option>
              <option v-for="group in groups" :key="group.name" :value="group.name">
                {{ group.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>Shell</label>
            <select v-model="userForm.shell">
              <option value="/bin/bash">/bin/bash</option>
              <option value="/bin/zsh">/bin/zsh</option>
              <option value="/bin/sh">/bin/sh</option>
              <option value="/usr/sbin/nologin">nologin (无登录权限)</option>
            </select>
          </div>

          <div class="form-group">
            <label>磁盘配额</label>
            <input
              v-model="userForm.quota"
              type="text"
              placeholder="例如: 100GB, 1TB"
            />
          </div>

          <div class="modal-footer">
            <button type="button" class="action-btn" @click="closeUserModal">
              取消
            </button>
            <button type="submit" class="action-btn primary">
              {{ isEditing ? '更新' : '创建' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- SSH Key Modal -->
    <div v-if="showAddSSHKeyModal" class="modal-overlay" @click="closeSSHKeyModal">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <h3>{{ editingSSHKey ? '编辑密钥' : '添加SSH密钥' }}</h3>
          <button class="close-btn" @click="closeSSHKeyModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveSSHKey" class="modal-body">
          <div class="form-group">
            <label>名称 *</label>
            <input
              v-model="sshKeyForm.name"
              type="text"
              required
              placeholder="密钥名称"
            />
          </div>

          <div class="form-group">
            <label>用户 *</label>
            <select v-model="sshKeyForm.user" required>
              <option value="">选择用户</option>
              <option v-for="user in users" :key="user.username" :value="user.username">
                {{ user.username }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>公钥 *</label>
            <textarea
              v-model="sshKeyForm.content"
              required
              rows="6"
              placeholder="ssh-rsa AAAA... 或 ssh-ed25519 AAAA..."
            ></textarea>
          </div>

          <div class="modal-footer">
            <button type="button" class="action-btn" @click="closeSSHKeyModal">
              取消
            </button>
            <button type="submit" class="action-btn primary">
              {{ editingSSHKey ? '更新' : '添加' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Quota Modal -->
    <div v-if="showQuotaModal" class="modal-overlay" @click="closeQuotaModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>磁盘配额管理</h3>
          <button class="close-btn" @click="closeQuotaModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="modal-body">
          <div class="quota-info">
            <p><strong>用户:</strong> {{ quotaUser?.username }}</p>
            <p><strong>当前配额:</strong> {{ quotaUser?.quota || '无限制' }}</p>
            <p><strong>已使用:</strong> {{ quotaUser?.usedSpace || '0 GB' }}</p>
          </div>

          <form @submit.prevent="saveQuota" class="quota-form">
            <div class="form-group">
              <label>空间限制</label>
              <input
                v-model="quotaForm.space"
                type="text"
                placeholder="例如: 100GB, 1TB, unlimited"
              />
            </div>

            <div class="form-group">
              <label>文件数量限制</label>
              <input
                v-model="quotaForm.files"
                type="text"
                placeholder="例如: 1000000, unlimited"
              />
            </div>

            <div class="modal-footer">
              <button type="button" class="action-btn" @click="closeQuotaModal">
                取消
              </button>
              <button type="submit" class="action-btn primary">
                应用配额
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import GroupManager from '../components/User/GroupManager.vue'
import {
  UserGroupIcon,
  KeyIcon,
  UserPlusIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  EyeIcon,
  EyeSlashIcon,
  CircleStackIcon
} from '@heroicons/vue/24/outline'
import { userApi, groupApi } from '../api'

const activeTab = ref('users')
const loading = ref(false)
const loadingSSH = ref(false)
const loadingGroups = ref(false)

const tabs = [
  { id: 'users', label: '用户', icon: UserGroupIcon },
  { id: 'ssh', label: 'SSH密钥', icon: KeyIcon },
  { id: 'groups', label: '用户组', icon: UserGroupIcon }
]

// User Management
const users = ref<any[]>([])
const showCreateUserModal = ref(false)
const editingUser = ref<any>(null)
const isEditing = computed(() => editingUser.value !== null)
const showPassword = ref(false)
const userForm = ref({
  username: '',
  password: '',
  comment: '',
  group: '',
  shell: '/bin/bash',
  quota: ''
})

// SSH Key Management
const sshKeys = ref<any[]>([])
const showAddSSHKeyModal = ref(false)
const editingSSHKey = ref(null)
const sshKeyForm = ref({
  name: '',
  user: '',
  content: ''
})

// Group Management
const groups = ref<any[]>([])
const showCreateGroupModal = ref(false)

// Quota Management
const showQuotaModal = ref(false)
const quotaUser = ref<any>(null)
const quotaForm = ref({
  space: '',
  files: ''
})

// API Functions
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await userApi.getUsers()
    // 处理可能的响应格式
    if (response && response.users) {
      users.value = response.users
    } else if (Array.isArray(response)) {
      users.value = response
    } else {
      users.value = []
    }
  } catch (error: any) {
    console.error('Failed to load users:', error)
    // 设置为空数组而不是Mock数据
    users.value = []
    // 可以选择显示错误提示给用户
  } finally {
    loading.value = false
  }
}

const loadSSHKeys = async () => {
  loadingSSH.value = true
  try {
    const response = await userApi.getSSHKeys()
    // 处理可能的响应格式
    if (response && response.keys) {
      sshKeys.value = response.keys
    } else if (Array.isArray(response)) {
      sshKeys.value = response
    } else {
      sshKeys.value = []
    }
  } catch (error: any) {
    console.error('Failed to load SSH keys:', error)
    // 设置为空数组而不是Mock数据
    sshKeys.value = []
    // 可以选择显示错误提示给用户
  } finally {
    loadingSSH.value = false
  }
}

const loadGroups = async () => {
  loadingGroups.value = true
  try {
    const response = await groupApi.getGroups()
    groups.value = response.groups || response  // axios拦截器已返回response.data
  } catch (error: any) {
    console.error('Failed to load groups:', error)
    // 设置为空数组而不是Mock数据
    groups.value = []
    // 可以选择显示错误提示给用户
  } finally {
    loadingGroups.value = false
  }
}

const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    username: user.username,
    password: '',
    comment: user.comment || '',
    group: user.group || '',
    shell: user.shell || '/bin/bash',
    quota: user.quota || ''
  }
  showCreateUserModal.value = true
}

const deleteUser = async (user: any) => {
  if (confirm(`确定要删除用户 "${user.username}" 吗?`)) {
    try {
      await userApi.deleteUser(user.username)
      await loadUsers()
    } catch (error: any) {
      console.error('Failed to delete user:', error)
      alert('删除失败: ' + error.message)
    }
  }
}

const saveUser = async () => {
  try {
    if (editingUser.value) {
      await userApi.updateUser(editingUser.value.username, {
        comment: userForm.value.comment,
        group: userForm.value.group,
        shell: userForm.value.shell
      })
    } else {
      await userApi.createUser({
        username: userForm.value.username,
        password: userForm.value.password,
        comment: userForm.value.comment,
        group: userForm.value.group,
        shell: userForm.value.shell
      })
    }
    closeUserModal()
    await loadUsers()
  } catch (error: any) {
    console.error('Failed to save user:', error)
    alert('保存失败: ' + error.message)
  }
}

const closeUserModal = () => {
  showCreateUserModal.value = false
  editingUser.value = null
  userForm.value = {
    username: '',
    password: '',
    comment: '',
    group: '',
    shell: '/bin/bash',
    quota: ''
  }
  showPassword.value = false
}

const viewKey = (key: any) => {
  alert(`公钥内容:\n\n${key.content}`)
}

const deleteSSHKey = async (key: any) => {
  if (confirm(`确定要删除密钥 "${key.name}" 吗?`)) {
    try {
      await userApi.deleteKey(key.id, key.user)
      await loadSSHKeys()
    } catch (error: any) {
      console.error('Failed to delete SSH key:', error)
      alert('删除失败: ' + error.message)
    }
  }
}

const saveSSHKey = async () => {
  try {
    await userApi.addKey({
      name: sshKeyForm.value.name,
      user: sshKeyForm.value.user,
      content: sshKeyForm.value.content
    })
    closeSSHKeyModal()
    await loadSSHKeys()
  } catch (error: any) {
    console.error('Failed to save SSH key:', error)
    alert('保存失败: ' + error.message)
  }
}

const closeSSHKeyModal = () => {
  showAddSSHKeyModal.value = false
  editingSSHKey.value = null
  sshKeyForm.value = {
    name: '',
    user: '',
    content: ''
  }
}

const manageUserQuota = (user: any) => {
  quotaUser.value = user
  quotaForm.value = {
    space: user.quota || '',
    files: ''
  }
  showQuotaModal.value = true
}

const saveQuota = async () => {
  try {
    // 调用配额API设置用户配额
    await userApi.setUserQuota(quotaUser.value?.username, {
      space: quotaForm.value.space,
      files: quotaForm.value.files
    })
    closeQuotaModal()
    await loadUsers()
    alert('配额设置成功')
  } catch (error: any) {
    console.error('Failed to set quota:', error)
    alert('配额设置失败: ' + error.message)
  }
}

const closeQuotaModal = () => {
  showQuotaModal.value = false
  quotaUser.value = null
  quotaForm.value = {
    space: '',
    files: ''
  }
}

const editGroup = (group: any) => {
  console.log('Editing group:', group.name)
  alert('编辑组: ' + group.name)
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

onMounted(() => {
  loadUsers()
  loadSSHKeys()
  loadGroups()
})
</script>

<style scoped>
.user-manager {
  width: 100%;
  height: 100%;
  padding: 32px;
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.manager-header {
  margin-bottom: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.manager-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 2px solid rgba(102, 126, 234, 0.1);
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: -2px;
  border-radius: 8px 8px 0 0;
}

.tab-btn:hover {
  color: #667eea;
  background: rgba(102, 126, 234, 0.1);
}

.tab-btn.active {
  color: white;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-bottom-color: transparent;
}

.tab-content {
  flex: 1;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
}

.action-btn.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.action-btn.danger {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
  border-color: transparent;
  color: white;
}

.action-btn.danger:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
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

.users-table {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: #f9fafb;
}

th {
  padding: 16px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  text-transform: uppercase;
}

tbody tr {
  border-bottom: 1px solid #e5e7eb;
}

tbody tr:hover {
  background: #f9fafb;
}

td {
  padding: 16px;
  font-size: 14px;
  color: #1f2937;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  font-weight: 600;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.icon-btn {
  padding: 6px;
  background: none;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.icon-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.icon-btn.danger:hover {
  background: #fee2e2;
  color: #ef4444;
  border-color: #ef4444;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
}

.status-badge.active {
  background: #d1fae5;
  color: #065f46;
}

.status-badge.inactive {
  background: #fee2e2;
  color: #991b1b;
}

.ssh-keys-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ssh-key-item {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 20px;
  align-items: center;
}

.key-info {
  display: flex;
  gap: 16px;
}

.key-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.key-details h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 6px;
}

.key-details p {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.key-fingerprint {
  font-family: monospace;
  font-size: 12px;
}

.key-actions {
  display: flex;
  gap: 8px;
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.group-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.group-header {
  display: flex;
  gap: 16px;
  align-items: center;
}

.group-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.group-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.group-info p {
  font-size: 14px;
  color: #6b7280;
}

.group-members {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.member-chip {
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 16px;
  font-size: 12px;
  color: #6b7280;
}

.member-chip.more {
  background: #e5e7eb;
  color: #9ca3af;
}

.group-actions {
  display: flex;
  justify-content: flex-end;
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

.modal-content.large {
  min-width: 600px;
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

.form-group input[type="text"],
.form-group input[type="password"],
.form-group select,
.form-group textarea {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.form-group input[type="text"]:disabled,
.form-group select:disabled {
  background: #f3f4f6;
  color: #6b7280;
}

.password-input {
  display: flex;
  gap: 8px;
}

.password-input input {
  flex: 1;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.quota-info {
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
  margin-bottom: 16px;
}

.quota-info p {
  margin-bottom: 8px;
  font-size: 14px;
  color: #374151;
}

.quota-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
</style>