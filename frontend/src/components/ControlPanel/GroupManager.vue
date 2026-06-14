<template>
  <div class="group-manager">
    <div class="manager-header">
      <h3>用户组管理</h3>
      <div class="manager-actions">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索用户组..."
          class="search-input"
        />
        <button @click="showAddGroupModal = true" class="btn btn-primary">
          <UserGroupIcon class="w-4 h-4" />
          添加用户组
        </button>
      </div>
    </div>

    <div class="groups-grid">
      <div v-for="group in filteredGroups" :key="group.name" class="group-card">
        <div class="group-header">
          <div class="group-info">
            <UsersIcon class="w-5 h-5" />
            <div>
              <h4>{{ group.name }}</h4>
              <p class="group-gid">GID: {{ group.gid }}</p>
            </div>
          </div>
          <div class="group-actions">
            <button @click="editGroup(group)" class="btn btn-sm btn-secondary">
              <PencilIcon class="w-3 h-3" />
            </button>
            <button
              @click="deleteGroup(group)"
              class="btn btn-sm btn-danger"
              :disabled="isSystemGroup(group.name)"
            >
              <TrashIcon class="w-3 h-3" />
            </button>
          </div>
        </div>

        <div class="group-details">
          <div class="detail-item">
            <span class="detail-label">成员数:</span>
            <span class="detail-value">{{ group.memberCount || 0 }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">权限级别:</span>
            <span class="detail-value">{{ getPermissionLevel(group.name) }}</span>
          </div>
        </div>

        <div class="group-members">
          <div class="members-header">
            <span class="members-title">成员列表</span>
            <button
              @click="manageMembers(group)"
              class="btn btn-sm btn-secondary"
            >
              <UserPlusIcon class="w-3 h-3" />
              管理成员
            </button>
          </div>
          <div class="members-list">
            <div
              v-for="member in (group.members || []).slice(0, 5)"
              :key="member"
              class="member-tag"
            >
              {{ member }}
            </div>
            <div
              v-if="(group.members || []).length > 5"
              class="member-tag more"
            >
              +{{ group.members.length - 5 }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑用户组模态框 -->
    <div v-if="showAddGroupModal || showEditGroupModal" class="modal-overlay" @click.self="closeModals">
      <div class="modal">
        <div class="modal-header">
          <h4>{{ editingGroup ? '编辑用户组' : '添加用户组' }}</h4>
          <button @click="closeModals" class="btn btn-ghost">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>组名 *</label>
            <input
              v-model="groupForm.name"
              type="text"
              placeholder="输入组名"
              :disabled="editingGroup"
            />
          </div>
          <div class="form-group" v-if="!editingGroup">
            <label>GID</label>
            <input
              v-model="groupForm.gid"
              type="number"
              placeholder="留空自动分配"
            />
          </div>
          <div class="form-group">
            <label>权限级别</label>
            <select v-model="groupForm.permissionLevel">
              <option value="user">普通用户</option>
              <option value="power">高级用户</option>
              <option value="admin">管理员</option>
            </select>
          </div>
          <div class="form-group">
            <label>系统权限</label>
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="groupForm.permissions.sudo" />
                sudo权限
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="groupForm.permissions.docker" />
                Docker权限
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="groupForm.permissions.smb" />
                SMB权限
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="groupForm.permissions.ftp" />
                FTP权限
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="groupForm.permissions.ssh" />
                SSH登录权限
              </label>
            </div>
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="groupForm.description"
              placeholder="用户组描述"
              rows="3"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModals" class="btn btn-secondary">取消</button>
          <button @click="saveGroup" class="btn btn-primary" :disabled="saving">
            {{ saving ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 成员管理模态框 -->
    <div v-if="showMembersModal" class="modal-overlay" @click.self="closeMembersModal">
      <div class="modal">
        <div class="modal-header">
          <h4>管理成员 - {{ currentGroup?.name }}</h4>
          <button @click="closeMembersModal" class="btn btn-ghost">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
        <div class="modal-body">
          <div class="members-section">
            <h5>当前成员</h5>
            <div class="current-members">
              <div
                v-for="member in currentGroupMembers"
                :key="member"
                class="member-item"
              >
                <span>{{ member }}</span>
                <button
                  @click="removeMember(member)"
                  class="btn btn-sm btn-danger"
                >
                  <TrashIcon class="w-3 h-3" />
                </button>
              </div>
              <div v-if="currentGroupMembers.length === 0" class="empty-state">
                暂无成员
              </div>
            </div>
          </div>

          <div class="members-section">
            <h5>添加成员</h5>
            <div class="add-members">
              <select v-model="selectedUserToAdd">
                <option value="">选择用户</option>
                <option
                  v-for="user in availableUsersToAdd"
                  :key="user.username"
                  :value="user.username"
                >
                  {{ user.username }}
                </option>
              </select>
              <button @click="addMember" class="btn btn-primary">
                <UserPlusIcon class="w-4 h-4" />
                添加
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { UsersIcon, UserGroupIcon, PencilIcon, TrashIcon, XMarkIcon, UserPlusIcon } from '@heroicons/vue/24/outline'

// 用户组数据
const groups = ref<any[]>([])
const users = ref<any[]>([])
const searchQuery = ref('')
const showAddGroupModal = ref(false)
const showEditGroupModal = ref(false)
const showMembersModal = ref(false)
const editingGroup = ref<any>(null)
const currentGroup = ref<any>(null)
const saving = ref(false)
const selectedUserToAdd = ref('')

// 系统组列表（不能删除）
const systemGroups = ['root', 'users', 'admin', 'sudo', 'docker', 'system']

// 组表单
const groupForm = ref({
  name: '',
  gid: '',
  permissionLevel: 'user',
  permissions: {
    sudo: false,
    docker: false,
    smb: false,
    ftp: false,
    ssh: false
  },
  description: ''
})

// 过滤后的组列表
const filteredGroups = computed(() => {
  if (!searchQuery.value) return groups.value
  return groups.value.filter(group =>
    group.name.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 当前组成员
const currentGroupMembers = computed(() => {
  return currentGroup.value?.members || []
})

// 可添加的用户（不在当前组中的用户）
const availableUsersToAdd = computed(() => {
  if (!currentGroup.value) return []
  return users.value.filter(user =>
    !currentGroupMembers.value.includes(user.username)
  )
})

// 判断是否为系统组
const isSystemGroup = (groupName: string) => {
  return systemGroups.includes(groupName)
}

// 获取权限级别
const getPermissionLevel = (groupName: string) => {
  const group = groups.value.find(g => g.name === groupName)
  return group?.permissionLevel || 'user'
}

// 获取用户组数据
const fetchGroups = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch('/api/groups', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
    if (response.ok) {
      const result = await response.json()
      // 处理后端返回的格式: { "groups": [...] }
      const groupsArray = result.groups || result || []
      groups.value = groupsArray.map((group: any) => ({
        ...group,
        gid: group.gid,
        memberCount: group.members?.length || 0,
        permissionLevel: group.permissionLevel || 'user',
        permissions: group.permissions || {
          sudo: false,
          docker: false,
          smb: false,
          ftp: false,
          ssh: false
        },
        description: group.description || ''
      }))
    }
  } catch (error) {
    console.error('Failed to fetch groups:', error)
    // 使用模拟数据
    groups.value = [
      {
        name: 'admin',
        gid: '1000',
        members: ['admin'],
        permissionLevel: 'admin',
        permissions: { sudo: true, docker: true, smb: true, ftp: true, ssh: true },
        description: '管理员组'
      },
      {
        name: 'users',
        gid: '100',
        members: ['user1', 'user2'],
        permissionLevel: 'user',
        permissions: { sudo: false, docker: false, smb: false, ftp: false, ssh: true },
        description: '普通用户组'
      },
      {
        name: 'docker',
        gid: '998',
        members: ['admin'],
        permissionLevel: 'power',
        permissions: { sudo: false, docker: true, smb: false, ftp: false, ssh: true },
        description: 'Docker用户组'
      }
    ]
  }
}

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
      users.value = usersArray
    }
  } catch (error) {
    console.error('Failed to fetch users:', error)
    // 使用模拟数据
    users.value = [
      { username: 'admin' },
      { username: 'user1' },
      { username: 'user2' }
    ]
  }
}

// 编辑用户组
const editGroup = (group: any) => {
  editingGroup.value = group
  groupForm.value = {
    name: group.name,
    gid: group.gid?.toString() || '',
    permissionLevel: group.permissionLevel || 'user',
    permissions: group.permissions || {
      sudo: false,
      docker: false,
      smb: false,
      ftp: false,
      ssh: false
    },
    description: group.description || ''
  }
  showEditGroupModal.value = true
}

// 删除用户组
const deleteGroup = async (group: any) => {
  if (isSystemGroup(group.name)) {
    alert('系统组不能删除')
    return
  }

  if (!confirm(`确定要删除用户组 "${group.name}" 吗？`)) return

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/groups/${group.name}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      groups.value = groups.value.filter(g => g.name !== group.name)
      alert('用户组删除成功')
    } else {
      alert('删除失败')
    }
  } catch (error) {
    console.error('Failed to delete group:', error)
    alert('删除失败')
  }
}

// 保存用户组
const saveGroup = async () => {
  if (!groupForm.value.name) {
    alert('请输入组名')
    return
  }

  saving.value = true

  try {
    const token = localStorage.getItem('token')
    const url = editingGroup.value
      ? `/api/groups/${groupForm.value.name}`
      : '/api/groups'

    const method = editingGroup.value ? 'PUT' : 'POST'

    const groupData = {
      name: groupForm.value.name,
      gid: groupForm.value.gid ? parseInt(groupForm.value.gid) : undefined,
      permissionLevel: groupForm.value.permissionLevel,
      permissions: groupForm.value.permissions,
      description: groupForm.value.description
    }

    const response = await fetch(url, {
      method,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(groupData)
    })

    if (response.ok) {
      await fetchGroups()
      closeModals()
      alert(editingGroup.value ? '用户组更新成功' : '用户组创建成功')
    } else {
      alert('操作失败')
    }
  } catch (error) {
    console.error('Failed to save group:', error)
    alert('操作失败')
  } finally {
    saving.value = false
  }
}

// 管理成员
const manageMembers = (group: any) => {
  currentGroup.value = group
  showMembersModal.value = true
}

// 添加成员
const addMember = async () => {
  if (!selectedUserToAdd.value) return

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/groups/${currentGroup.value.name}/members`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username: selectedUserToAdd.value })
    })

    if (response.ok) {
      await fetchGroups()
      selectedUserToAdd.value = ''
      alert('成员添加成功')
    } else {
      alert('添加失败')
    }
  } catch (error) {
    console.error('Failed to add member:', error)
    alert('添加失败')
  }
}

// 移除成员
const removeMember = async (username: string) => {
  if (!confirm(`确定要移除用户 "${username}" 吗？`)) return

  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/groups/${currentGroup.value.name}/members/${username}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      await fetchGroups()
      alert('成员移除成功')
    } else {
      alert('移除失败')
    }
  } catch (error) {
    console.error('Failed to remove member:', error)
    alert('移除失败')
  }
}

// 关闭模态框
const closeModals = () => {
  showAddGroupModal.value = false
  showEditGroupModal.value = false
  editingGroup.value = null
  groupForm.value = {
    name: '',
    gid: '',
    permissionLevel: 'user',
    permissions: {
      sudo: false,
      docker: false,
      smb: false,
      ftp: false,
      ssh: false
    },
    description: ''
  }
}

// 关闭成员管理模态框
const closeMembersModal = () => {
  showMembersModal.value = false
  currentGroup.value = null
  selectedUserToAdd.value = ''
}

onMounted(() => {
  fetchGroups()
  fetchUsers()
})
</script>

<style scoped>
.group-manager {
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

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.group-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s;
}

.group-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.group-header {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.group-gid {
  font-size: 12px;
  color: #6b7280;
}

.group-actions {
  display: flex;
  gap: 4px;
}

.group-details {
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 14px;
}

.detail-label {
  color: #6b7280;
}

.detail-value {
  font-weight: 500;
  color: #1f2937;
}

.group-members {
  padding: 16px;
}

.members-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.members-title {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.members-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.member-tag {
  padding: 4px 10px;
  background: #dbeafe;
  color: #1e40af;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.member-tag.more {
  background: #f3f4f6;
  color: #6b7280;
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
  max-width: 600px;
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

.members-section {
  margin-bottom: 24px;
}

.members-section h5 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.current-members {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 8px;
}

.add-members {
  display: flex;
  gap: 8px;
}

.add-members select {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.empty-state {
  color: #6b7280;
  font-style: italic;
}

.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>