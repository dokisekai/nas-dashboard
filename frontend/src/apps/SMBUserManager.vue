<template>
  <div class="smb-user-manager">
    <!-- 头部 -->
    <div class="smb-header">
      <div class="header-left">
        <h1>SMB 用户管理</h1>
        <p class="subtitle">管理SMB服务和用户访问</p>
      </div>
      <div class="header-actions">
        <button class="action-btn primary" @click="refreshUsers">
          <ArrowPathIcon class="w-4 h-4" />
          刷新
        </button>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="smb-content">
      <!-- SMB用户列表 -->
      <div class="smb-users-section">
        <div class="section-header">
          <h2>SMB 用户</h2>
          <div class="user-stats">
            <span class="stat-item">
              <span class="stat-value">{{ enabledUsers.length }}</span>
              <span class="stat-label">启用用户</span>
            </span>
            <span class="stat-item">
              <span class="stat-value">{{ smbUsers.length }}</span>
              <span class="stat-label">总用户数</span>
            </span>
          </div>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载SMB用户...</p>
        </div>

        <div v-else class="users-table">
          <table>
            <thead>
              <tr>
                <th>用户名</th>
                <th>状态</th>
                <th>最后登录</th>
                <th>登录次数</th>
                <th>主目录</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in smbUsers" :key="user.username">
                <td>
                  <div class="user-info">
                    <div class="user-avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                    <div class="user-details">
                      <span class="user-name">{{ user.username }}</span>
                      <span class="user-meta">UID: {{ user.UID }}</span>
                    </div>
                  </div>
                </td>
                <td>
                  <span
                    class="status-badge"
                    :class="{ active: user.enabled, inactive: !user.enabled }"
                  >
                    {{ user.enabled ? '已启用' : '未启用' }}
                  </span>
                </td>
                <td>{{ user.lastLogin }}</td>
                <td>{{ user.loginCount }}</td>
                <td>{{ user.homeDirectory || '-' }}</td>
                <td>
                  <div class="action-buttons">
                    <button
                      class="icon-btn"
                      @click="managePassword(user)"
                      title="设置密码"
                    >
                      <KeyIcon class="w-4 h-4" />
                    </button>
                    <button
                      class="icon-btn"
                      @click="toggleUserStatus(user)"
                      :class="{ danger: user.enabled }"
                      :title="user.enabled ? '禁用' : '启用'"
                    >
                      user.enabled ? <NoSymbolIcon class="w-4 h-4" /> : <CheckIcon class="w-4 h-4" />
                    </button>
                    <button
                      class="icon-btn"
                      @click="viewUserStats(user)"
                      title="统计信息"
                    >
                      <ChartBarIcon class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- SMB会话管理 -->
      <div class="smb-sessions-section">
        <div class="section-header">
          <h2>活动会话</h2>
          <div class="session-actions">
            <button class="action-btn" @click="refreshSessions">
              <ArrowPathIcon class="w-4 h-4" />
              刷新
            </button>
            <button
              class="action-btn danger"
              @click="disconnectAllSessions"
              :disabled="sessions.length === 0"
            >
              <ArrowRightOnRectangleIcon class="w-4 h-4" />
              断开所有
            </button>
          </div>
        </div>

        <div v-if="loadingSessions" class="loading-state">
          <div class="spinner"></div>
          <p>加载会话信息...</p>
        </div>

        <div v-else-if="sessions.length > 0" class="sessions-table">
          <table>
            <thead>
              <tr>
                <th>用户</th>
                <th>机器</th>
                <th>IP地址</th>
                <th>协议</th>
                <th>连接时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="session in sessions" :key="session.clientId">
                <td>{{ session.username }}</td>
                <td>{{ session.machine }}</td>
                <td>{{ session.ip }}</td>
                <td>
                  <span class="protocol-badge">{{ session.protocol }}</span>
                </td>
                <td>{{ session.connectTime }}</td>
                <td>
                  <button
                    class="icon-btn danger"
                    @click="disconnectSession(session)"
                    title="断开连接"
                  >
                    <XMarkIcon class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-else class="empty-state">
          <UserGroupIcon class="w-12 h-12" />
          <p>当前没有活动会话</p>
        </div>
      </div>
    </div>

    <!-- 设置密码对话框 -->
    <div v-if="showPasswordDialog" class="modal-overlay" @click="showPasswordDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>设置SMB密码</h3>
          <button class="close-btn" @click="showPasswordDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="setSMBPassword" class="modal-body">
          <div class="form-group">
            <label>用户</label>
            <input type="text" :value="selectedUser?.username" disabled />
          </div>

          <div class="form-group">
            <label>新密码</label>
            <input
              type="password"
              v-model="passwordForm.newPassword"
              placeholder="输入新密码"
              required
            />
          </div>

          <div class="form-group">
            <label>确认密码</label>
            <input
              type="password"
              v-model="passwordForm.confirmPassword"
              placeholder="再次输入新密码"
              required
            />
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showPasswordDialog = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              设置密码
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 用户统计对话框 -->
    <div v-if="showStatsDialog" class="modal-overlay" @click="showStatsDialog = false">
      <div class="modal-content stats-content" @click.stop>
        <div class="modal-header">
          <h3>用户统计信息</h3>
          <button class="close-btn" @click="showStatsDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="modal-body">
          <div v-if="loadingStats" class="loading-state">
            <div class="spinner"></div>
            <p>加载统计信息...</p>
          </div>

          <div v-else class="stats-grid">
            <div class="stat-card">
              <div class="stat-icon">
                <UserIcon class="w-8 h-8" />
              </div>
              <div class="stat-info">
                <h4>用户状态</h4>
                <p class="stat-value">{{ userStats.enabled ? '已启用' : '未启用' }}</p>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon">
                <LinkIcon class="w-8 h-8" />
              </div>
              <div class="stat-info">
                <h4>活动会话</h4>
                <p class="stat-value">{{ userStats.sessions?.length || 0 }}</p>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon">
                <ClockIcon class="w-8 h-8" />
              </div>
              <div class="stat-info">
                <h4>最后活动</h4>
                <p class="stat-value">{{ userStats.lastActivity }}</p>
              </div>
            </div>

            <div class="stat-card">
              <div class="stat-icon">
                <ChartBarIcon class="w-8 h-8" />
              </div>
              <div class="stat-info">
                <h4>访问次数</h4>
                <p class="stat-value">{{ userStats.accessCount || 0 }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  UserIcon,
  KeyIcon,
  ArrowPathIcon,
  NoSymbolIcon,
  CheckIcon,
  ChartBarIcon,
  XMarkIcon,
  UserGroupIcon,
  ArrowRightOnRectangleIcon,
  ClockIcon,
  LinkIcon
} from '@heroicons/vue/24/outline'

// 状态
const loading = ref(false)
const loadingSessions = ref(false)
const loadingStats = ref(false)
const smbUsers = ref<any[]>([])
const sessions = ref<any[]>([])

// 对话框状态
const showPasswordDialog = ref(false)
const showStatsDialog = ref(false)
const selectedUser = ref<any>(null)

// 表单数据
const passwordForm = reactive({
  newPassword: '',
  confirmPassword: ''
})

// 用户统计数据
const userStats = ref<any>({})

// 计算属性
const enabledUsers = computed(() => smbUsers.value.filter(u => u.enabled))

// 方法
const loadSMBUsers = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/smb/users', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      smbUsers.value = data.users || []
    } else {
      throw new Error('获取SMB用户失败')
    }
  } catch (error) {
    console.error('Failed to load SMB users:', error)
    ElMessage.error('获取SMB用户失败')
  } finally {
    loading.value = false
  }
}

const loadSessions = async () => {
  loadingSessions.value = true
  try {
    const response = await fetch('/api/smb/sessions', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      sessions.value = data.sessions || []
    } else {
      throw new Error('获取SMB会话失败')
    }
  } catch (error) {
    console.error('Failed to load SMB sessions:', error)
    ElMessage.error('获取SMB会话失败')
  } finally {
    loadingSessions.value = false
  }
}

const managePassword = (user: any) => {
  selectedUser.value = user
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  showPasswordDialog.value = true
}

const setSMBPassword = async () => {
  if (!passwordForm.newPassword || !passwordForm.confirmPassword) {
    ElMessage.warning('请输入密码')
    return
  }

  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的密码不一致')
    return
  }

  if (passwordForm.newPassword.length < 6) {
    ElMessage.warning('密码长度至少6位')
    return
  }

  try {
    const response = await fetch(`/api/smb/users/${selectedUser.value.username}/password`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        password: passwordForm.newPassword
      })
    })

    if (response.ok) {
      ElMessage.success('SMB密码设置成功')
      showPasswordDialog.value = false
      await loadSMBUsers()
    } else {
      throw new Error('设置密码失败')
    }
  } catch (error) {
    console.error('Failed to set SMB password:', error)
    ElMessage.error('设置SMB密码失败')
  }
}

const toggleUserStatus = async (user: any) => {
  const action = user.enabled ? 'disable' : 'enable'
  const actionText = user.enabled ? '禁用' : '启用'

  try {
    const response = await fetch(`/api/smb/users/${user.username}/${action}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      ElMessage.success(`用户${actionText}成功`)
      await loadSMBUsers()
    } else {
      throw new Error(`${actionText}用户失败`)
    }
  } catch (error) {
    console.error('Failed to toggle user status:', error)
    ElMessage.error(`${actionText}用户失败`)
  }
}

const viewUserStats = async (user: any) => {
  selectedUser.value = user
  loadingStats.value = true
  showStatsDialog.value = true

  try {
    const response = await fetch(`/api/smb/users/${user.username}/stats`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      userStats.value = await response.json()
    } else {
      throw new Error('获取用户统计失败')
    }
  } catch (error) {
    console.error('Failed to get user stats:', error)
    ElMessage.error('获取用户统计失败')
  } finally {
    loadingStats.value = false
  }
}

const refreshSessions = () => {
  loadSessions()
}

const disconnectSession = async (session: any) => {
  try {
    const response = await fetch(`/api/smb/sessions/${session.pid}`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      ElMessage.success('会话已断开')
      await loadSessions()
    } else {
      throw new Error('断开会话失败')
    }
  } catch (error) {
    console.error('Failed to disconnect session:', error)
    ElMessage.error('断开会话失败')
  }
}

const disconnectAllSessions = async () => {
  if (!confirm('确定要断开所有SMB会话吗？')) {
    return
  }

  try {
    const response = await fetch('/api/smb/sessions', {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      ElMessage.success('所有会话已断开')
      await loadSessions()
    } else {
      throw new Error('断开会话失败')
    }
  } catch (error) {
    console.error('Failed to disconnect all sessions:', error)
    ElMessage.error('断开会话失败')
  }
}

const refreshUsers = () => {
  loadSMBUsers()
  loadSessions()
}

// 生命周期
onMounted(() => {
  loadSMBUsers()
  loadSessions()
})
</script>

<style scoped lang="scss">
.smb-user-manager {
  width: 100%;
  padding: 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.smb-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.header-left h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &.primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  &:hover:not(.primary) {
    background: rgba(102, 126, 234, 0.1);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &.danger {
    background: rgba(239, 68, 68, 0.1);
    color: #ef4444;

    &:hover:not(:disabled) {
      background: rgba(239, 68, 68, 0.2);
    }
  }
}

.smb-content {
  display: grid;
  gap: 32px;
}

.smb-users-section,
.smb-sessions-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h2 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.user-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;

  .stat-value {
    font-size: 24px;
    font-weight: 600;
    color: #667eea;
    line-height: 1;
  }

  .stat-label {
    font-size: 12px;
    color: #6b7280;
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #6b7280;

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(102, 126, 234, 0.2);
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 12px;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.users-table,
.sessions-table {
  width: 100%;

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: linear-gradient(to bottom, #f8fafc, #e2e8f0);

    tr {
      th {
        padding: 12px;
        text-align: left;
        font-size: 13px;
        font-weight: 600;
        color: #374151;
        border-bottom: 1px solid #e5e7eb;
      }
    }
  }

  tbody {
    tr {
      border-bottom: 1px solid #f3f4f6;
      transition: background 0.2s;

      &:hover {
        background: rgba(102, 126, 234, 0.05);
      }

      &:last-child {
        border-bottom: none;
      }

      td {
        padding: 16px 12px;
        font-size: 14px;
      }
    }
  }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 2px;

  .user-name {
    font-weight: 500;
    color: #1f2937;
  }

  .user-meta {
    font-size: 12px;
    color: #9ca3af;
  }
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;

  &.active {
    background: rgba(34, 197, 94, 0.1);
    color: #066;
  }

  &.inactive {
    background: rgba(239, 68, 68, 0.1);
    color: #991b1b;
  }
}

.protocol-badge {
  padding: 4px 8px;
  background: rgba(59, 130, 246, 0.1);
  color: #2563eb;
  border-radius: 4px;
  font-size: 11px;
  font-family: monospace;
}

.action-buttons {
  display: flex;
  gap: 4px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover:not(.danger) {
    background: rgba(102, 126, 234, 0.1);
  }

  &.danger {
    color: #ef4444;

    &:hover {
      background: rgba(239, 68, 68, 0.1);
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #9ca3af;

  svg {
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    font-size: 14px;
  }
}

.session-actions {
  display: flex;
  gap: 12px;
}

// 模态对话框
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
  min-width: 400px;
  max-width: 600px;
  width: 100%;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);

  &.stats-content {
    max-width: 800px;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    background: rgba(102, 126, 234, 0.1);
  }
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

  label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
  }

  input {
    padding: 10px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    transition: all 0.2s;

    &:focus {
      outline: none;
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }

    &:disabled {
      background: #f9fafb;
      color: #9ca3af;
    }
  }
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;

  &:hover {
    opacity: 0.9;
  }
}

.btn-secondary {
  background: white;
  border: 1px solid #e5e7eb;
  color: #6b7280;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border: 1px solid rgba(102, 126, 234, 0.1);
  border-radius: 8px;

  .stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .stat-info {
    flex: 1;

    h4 {
      font-size: 14px;
      font-weight: 500;
      color: #6b7280;
      margin: 0 0 4px 0;
    }

    .stat-value {
      font-size: 16px;
      font-weight: 600;
      color: #1f2937;
    }
  }
}
</style>