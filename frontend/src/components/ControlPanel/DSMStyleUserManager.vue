<template>
  <div class="dsm-user-manager">
    <!-- 顶部标签页 -->
    <div class="dsm-tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['dsm-tab', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- 用户管理 -->
    <div v-if="activeTab === 'users'" class="dsm-content">
      <div class="dsm-toolbar">
        <div class="dsm-title">用户</div>
        <div class="dsm-actions">
          <button class="dsm-btn" @click="showAddUser = true">
            <PlusIcon class="w-4 h-4" />
            新增
          </button>
        </div>
      </div>

      <!-- 用户列表 - 简单表格 -->
      <div class="dsm-table">
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
                <span :class="['dsm-status', user.status === 'active' ? 'active' : 'disabled']">
                  {{ user.status === 'active' ? '正常' : '停用' }}
                </span>
              </td>
              <td>
                <button class="dsm-link" @click="editUser(user)">编辑</button>
                <button class="dsm-link danger" @click="deleteUser(user.id)">删除</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 用户组管理 -->
    <div v-if="activeTab === 'groups'" class="dsm-content">
      <div class="dsm-toolbar">
        <div class="dsm-title">用户组</div>
        <div class="dsm-actions">
          <button class="dsm-btn" @click="showAddGroup = true">
            <PlusIcon class="w-4 h-4" />
            新增
          </button>
        </div>
      </div>

      <!-- 用户组列表 - 折叠式 -->
      <div class="dsm-groups">
        <div
          v-for="group in groups"
          :key="group.id"
          class="dsm-group-item"
          :class="{ expanded: expandedGroups.includes(group.id) }"
        >
          <div class="dsm-group-header" @click="toggleGroup(group.id)">
            <div class="group-info">
              <ChevronRightIcon class="expand-icon" />
              <span class="group-name">{{ group.name }}</span>
            </div>
            <div class="group-actions">
              <button class="dsm-link" @click.stop="editGroup(group)">编辑</button>
              <button
                v-if="!isSystemGroup(group)"
                class="dsm-link danger"
                @click.stop="deleteGroup(group.id)"
              >
                删除
              </button>
            </div>
          </div>

          <div v-if="expandedGroups.includes(group.id)" class="dsm-group-details">
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

    <!-- 添加用户对话框 -->
    <div v-if="showAddUser" class="dsm-modal" @click.self="closeUserModal">
      <div class="dsm-modal-content">
        <div class="dsm-modal-header">
          <h3>{{ editingUser ? '编辑用户' : '新增用户' }}</h3>
          <button class="dsm-close" @click="closeUserModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveUser" class="dsm-form">
          <div class="dsm-form-group">
            <label>用户名 *</label>
            <input
              v-model="userForm.name"
              type="text"
              required
              :disabled="editingUser"
              placeholder="用户名"
            />
          </div>

          <div class="dsm-form-group">
            <label>显示名称</label>
            <input v-model="userForm.displayName" type="text" placeholder="显示名称" />
          </div>

          <div class="dsm-form-group" v-if="!editingUser">
            <label>密码 *</label>
            <input v-model="userForm.password" type="password" required placeholder="密码" />
          </div>

          <div class="dsm-form-group">
            <label>用户组</label>
            <div class="dsm-checkbox-group">
              <label
                v-for="group in groups"
                :key="group.id"
                class="dsm-checkbox"
              >
                <input type="checkbox" v-model="userForm.groups" :value="group.name" />
                {{ group.name }}
              </label>
            </div>
          </div>

          <div class="dsm-form-group">
            <label>用户角色</label>
            <div class="dsm-radio-group">
              <label class="dsm-radio">
                <input type="radio" v-model="userForm.role" value="user" />
                <span>普通用户</span>
              </label>
              <label class="dsm-radio">
                <input type="radio" v-model="userForm.role" value="admin" />
                <span>管理员</span>
              </label>
            </div>
          </div>

          <div class="dsm-modal-footer">
            <button type="button" class="dsm-btn secondary" @click="closeUserModal">
              取消
            </button>
            <button type="submit" class="dsm-btn primary">
              {{ editingUser ? '确定' : '确定' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 添加用户组对话框 -->
    <div v-if="showAddGroup" class="dsm-modal" @click.self="closeGroupModal">
      <div class="dsm-modal-content">
        <div class="dsm-modal-header">
          <h3>{{ editingGroup ? '编辑用户组' : '新增用户组' }}</h3>
          <button class="dsm-close" @click="closeGroupModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveGroup" class="dsm-form">
          <div class="dsm-form-group">
            <label>用户组名称 *</label>
            <input
              v-model="groupForm.name"
              type="text"
              required
              :disabled="editingGroup?.isSystem"
              placeholder="用户组名称"
            />
          </div>

          <div class="dsm-form-group">
            <label>描述</label>
            <textarea
              v-model="groupForm.description"
              rows="3"
              placeholder="用户组描述"
            ></textarea>
          </div>

          <div class="dsm-form-group">
            <label>权限级别</label>
            <div class="dsm-radio-group">
              <label class="dsm-radio">
                <input type="radio" v-model="groupForm.level" value="read" />
                <span>只读</span>
              </label>
              <label class="dsm-radio">
                <input type="radio" v-model="groupForm.level" value="write" />
                <span>读写</span>
              </label>
              <label class="dsm-radio">
                <input type="radio" v-model="groupForm.level" value="admin" />
                <span>管理</span>
              </label>
            </div>
          </div>

          <div class="dsm-modal-footer">
            <button type="button" class="dsm-btn secondary" @click="closeGroupModal">
              取消
            </button>
            <button type="submit" class="dsm-btn primary">
              {{ editingGroup ? '确定' : '确定' }}
            </button>
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
  XMarkIcon
} from '@heroicons/vue/24/outline'

const activeTab = ref('users')

const tabs = [
  { id: 'users', label: '用户' },
  { id: 'groups', label: '用户组' }
]

// 模拟数据
const users = ref([
  { id: '1', displayName: 'admin', name: 'admin', groups: ['administrators'], status: 'active' },
  { id: '2', displayName: '张三', name: 'zhangsan', groups: ['users'], status: 'active' },
  { id: '3', displayName: '李四', name: 'lisi', groups: ['users'], status: 'active' },
  { id: '4', displayName: '访客', name: 'guest', groups: ['guests'], status: 'active' }
])

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
    members: ['zhangsan', 'lisi', 'guest'],
    permissions: ['文件访问', '个人设置']
  },
  {
    id: 'guests',
    name: 'guests',
    description: '访客组，只读访问权限',
    isSystem: true,
    memberCount: 1,
    members: ['guest'],
    permissions: ['只读访问']
  }
])

const expandedGroups = ref<string[]>([])

const showAddUser = ref(false)
const showAddGroup = ref(false)
const editingUser = ref<any>(null)
const editingGroup = ref<any>(null)

const userForm = ref({
  name: '',
  displayName: '',
  password: '',
  role: 'user',
  groups: ['users']
})

const groupForm = ref({
  name: '',
  description: '',
  level: 'read'
})

const toggleGroup = (groupId: string) => {
  const index = expandedGroups.value.indexOf(groupId)
  if (index > -1) {
    expandedGroups.value.splice(index, 1)
  } else {
    expandedGroups.value.push(groupId)
  }
}

const isSystemGroup = (group: any) => group.isSystem

const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    name: user.name,
    displayName: user.displayName,
    password: '',
    role: user.groups.includes('administrators') ? 'admin' : 'user',
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
    description: group.description,
    level: group.permissions.includes('管理') ? 'admin' : 'read'
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
        description: groupForm.value.description,
        permissions: getPermissionsForLevel(groupForm.value.level)
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
      permissions: getPermissionsForLevel(groupForm.value.level)
    })
  }
  closeGroupModal()
}

const getPermissionsForLevel = (level: string) => {
  const perms: Record<string, string[]> = {
    read: ['只读访问'],
    write: ['读写访问'],
    admin: ['用户管理', '系统设置', '文件管理']
  }
  return perms[level] || ['只读访问']
}

const closeUserModal = () => {
  showAddUser.value = false
  editingUser.value = null
  userForm.value = {
    name: '',
    displayName: '',
    password: '',
    role: 'user',
    groups: ['users']
  }
}

const closeGroupModal = () => {
  showAddGroup.value = false
  editingGroup.value = null
  groupForm.value = {
    name: '',
    description: '',
    level: 'read'
  }
}
</script>

<style scoped>
.dsm-user-manager {
  width: 100%;
  height: 100%;
  background: #f5f5f5;
}

.dsm-tabs {
  display: flex;
  border-bottom: 1px solid #d8d8d8;
  background: white;
}

.dsm-tab {
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  transition: all 0.2s;
}

.dsm-tab:hover {
  color: #333;
}

.dsm-tab.active {
  color: #1890ff;
  border-bottom-color: #1890ff;
}

.dsm-content {
  padding: 24px;
}

.dsm-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e8e8e8;
}

.dsm-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.dsm-actions {
  display: flex;
  gap: 8px;
}

.dsm-btn {
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

.dsm-btn:hover {
  border-color: #1890ff;
  color: #1890ff;
}

.dsm-btn.primary {
  background: #1890ff;
  border-color: #1890ff;
  color: white;
}

.dsm-btn.primary:hover {
  background: #40a9ff;
  border-color: #40a9ff;
}

.dsm-btn.secondary {
  background: white;
  border-color: #d9d9d9;
}

.dsm-table {
  background: white;
  border: 1px solid #e8e8e8;
  border-collapse: collapse;
}

.dsm-table th,
.dsm-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e8e8e8;
}

.dsm-table th {
  background: #fafafa;
  font-weight: 500;
  color: #666;
  font-size: 13px;
}

.dsm-table td {
  font-size: 14px;
  color: #333;
}

.dsm-status {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.dsm-status.active {
  background: #52c41a;
  color: white;
}

.dsm-status.disabled {
  background: #d9d9d9;
  color: #666;
}

.dsm-link {
  background: none;
  border: none;
  color: #1890ff;
  cursor: pointer;
  padding: 4px 8px;
  font-size: 13px;
  transition: all 0.2s;
}

.dsm-link:hover {
  color: #40a9ff;
}

.dsm-link.danger {
  color: #ff4d4f;
}

.dsm-link.danger:hover {
  color: #ff7875;
}

.dsm-groups {
  background: white;
  border: 1px solid #e8e8e8;
}

.dsm-group-item {
  border-bottom: 1px solid #e8e8e8;
}

.dsm-group-item:last-child {
  border-bottom: none;
}

.dsm-group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.dsm-group-header:hover {
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

.dsm-group-item.expanded .expand-icon {
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

.dsm-group-details {
  padding: 16px;
  background: #fafafa;
  border-top: 1px solid #e8e8e8;
}

.group-description {
  font-size: 13px;
  color: #666;
  margin-bottom: 12px;
}

.group-permissions {
  font-size: 13px;
  color: #666;
  margin-bottom: 12px;
}

.perm-tag {
  display: inline-block;
  padding: 2px 6px;
  background: white;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  margin-right: 8px;
  font-size: 12px;
}

.group-members {
  font-size: 13px;
  color: #666;
}

.member-tag {
  display: inline-block;
  padding: 2px 6px;
  background: white;
  border: 1px solid #d9d9d9;
  border-radius: 2px;
  margin-right: 8px;
  font-size: 12px;
}

.dsm-modal {
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

.dsm-modal-content {
  background: white;
  border-radius: 8px;
  width: 500px;
  max-width: 90%;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.dsm-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e8e8e8;
}

.dsm-modal-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.dsm-close {
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

.dsm-close:hover {
  background: #f5f5f5;
}

.dsm-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dsm-form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dsm-form-group label {
  font-size: 13px;
  font-weight: 500;
  color: #666;
}

.dsm-form-group input,
.dsm-form-group textarea,
.dsm-form-group select {
  padding: 8px 12px;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  font-size: 14px;
}

.dsm-checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dsm-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #333;
}

.dsm-radio-group {
  display: flex;
  gap: 16px;
}

.dsm-radio {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #333;
}

.dsm-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 24px;
  border-top: 1px solid #e8e8e8;
}

@media (max-width: 768px) {
  .dsm-toolbar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }

  .dsm-table {
    font-size: 12px;
  }

  .dsm-table th,
  .dsm-table td {
    padding: 8px;
  }
}
</style>