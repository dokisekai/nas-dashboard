<template>
  <div class="group-manager">
    <div class="manager-header">
      <h3>用户组管理</h3>
      <div class="header-actions">
        <button class="action-btn primary" @click="showCreateGroupModal = true">
          <UserPlusIcon class="w-4 h-4" />
          创建组
        </button>
        <button class="action-btn secondary" @click="loadGroups">
          <ArrowPathIcon class="w-4 h-4" />
          刷新
        </button>
      </div>
    </div>

    <div v-if="loadingGroups" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else class="groups-grid">
      <div
        v-for="group in groups"
        :key="group.name"
        class="group-card"
      >
        <div class="group-header">
          <div class="group-info">
            <div class="group-icon">
              <UserGroupIcon class="w-6 h-6" />
            </div>
            <div>
              <h4>{{ group.name }}</h4>
              <p class="group-gid">GID: {{ group.gid }}</p>
            </div>
          </div>
          <div class="group-actions">
            <button class="icon-btn" @click="editGroup(group)" title="编辑">
              <PencilIcon class="w-4 h-4" />
            </button>
            <button class="icon-btn" @click="manageMembers(group)" title="成员管理">
              <UsersIcon class="w-4 h-4" />
            </button>
            <button class="icon-btn danger" @click="deleteGroup(group)" title="删除">
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <div class="group-details">
          <div class="detail-item">
            <span class="detail-label">成员数量:</span>
            <span class="detail-value">{{ getMemberCount(group) }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">主要成员:</span>
            <span class="detail-value">{{ getPrimaryMembers(group) }}</span>
          </div>
        </div>

        <div class="group-members-preview">
          <div class="members-label">组成员:</div>
          <div class="members-list">
            <div
              v-for="member in getGroupMembers(group)"
              :key="member"
              class="member-chip"
            >
              {{ member }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑组模态框 -->
    <div v-if="showCreateGroupModal || editingGroup" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingGroup ? '编辑组' : '创建组' }}</h3>
          <button class="icon-btn" @click="closeModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveGroup" class="group-form">
          <div class="form-group">
            <label>组名称</label>
            <input
              v-model="groupForm.name"
              type="text"
              placeholder="groupname"
              :disabled="editingGroup !== null"
              required
            />
            <div class="form-hint">
              {{ editingGroup ? '组名称不可修改' : '只能包含小写字母、数字和下划线' }}
            </div>
          </div>

          <div class="form-group">
            <label>GID (可选)</label>
            <input
              v-model="groupForm.gid"
              type="number"
              placeholder="自动分配"
              :disabled="editingGroup !== null"
            />
            <div class="form-hint">留空则自动分配GID</div>
          </div>

          <div class="form-group">
            <label>成员用户</label>
            <div class="members-selector">
              <div class="selected-members">
                <div
                  v-for="member in groupForm.members"
                  :key="member"
                  class="member-tag"
                >
                  {{ member }}
                  <button type="button" @click="removeMember(member)" class="remove-btn">×</button>
                </div>
              </div>
              <div class="available-members">
                <div
                  v-for="user in availableUsers"
                  :key="user.username"
                  :class="['member-option', { selected: groupForm.members.includes(user.username) }]"
                  @click="toggleMember(user.username)"
                >
                  {{ user.username }}
                </div>
              </div>
            </div>
          </div>

          <div class="form-actions">
            <button type="button" class="action-btn secondary" @click="closeModal">
              取消
            </button>
            <button type="submit" class="action-btn primary" :disabled="saving">
              {{ saving ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 成员管理模态框 -->
    <div v-if="managingGroup" class="modal-overlay" @click="closeMemberModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>成员管理 - {{ managingGroup?.name }}</h3>
          <button class="icon-btn" @click="closeMemberModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="member-management">
          <div class="current-members">
            <h4>当前成员</h4>
            <div class="members-list">
              <div
                v-for="member in getGroupMembers(managingGroup)"
                :key="member"
                class="member-item"
              >
                <span>{{ member }}</span>
                <button
                  v-if="member !== managingGroup.name"
                  class="icon-btn danger"
                  @click="removeGroupMember(managingGroup, member)"
                  title="移除成员"
                >
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>

          <div class="available-users">
            <h4>添加成员</h4>
            <div class="users-list">
              <div
                v-for="user in availableUsersForGroup"
                :key="user.username"
                :class="['user-item', { inGroup: isUserInGroup(user.username, managingGroup) }]"
              >
                <span>{{ user.username }}</span>
                <button
                  v-if="!isUserInGroup(user.username, managingGroup)"
                  class="icon-btn primary"
                  @click="addGroupMember(managingGroup, user.username)"
                  title="添加到组"
                >
                  <PlusIcon class="w-4 h-4" />
                </button>
                <button
                  v-else
                  class="icon-btn danger"
                  @click="removeGroupMember(managingGroup, user.username)"
                  title="从组移除"
                >
                  <MinusIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  UserPlusIcon,
  ArrowPathIcon,
  PencilIcon,
  TrashIcon,
  UserGroupIcon,
  UsersIcon,
  XMarkIcon,
  PlusIcon,
  MinusIcon
} from '@heroicons/vue/24/outline'
import { groupApi, userApi } from '../../api'

// 状态
const groups = ref<any[]>([])
const loadingGroups = ref(false)
const showCreateGroupModal = ref(false)
const editingGroup = ref<any>(null)
const managingGroup = ref<any>(null)
const saving = ref(false)

const groupForm = ref({
  name: '',
  gid: '',
  members: [] as string[]
})

const availableUsers = ref<any[]>([])

// 计算属性
const availableUsersForGroup = computed(() => {
  if (!managingGroup.value) return []
  const currentMembers = getGroupMembers(managingGroup.value)
  return availableUsers.value.filter(user => !currentMembers.includes(user.username))
})

// 方法
const loadGroups = async () => {
  loadingGroups.value = true
  try {
    const response = await groupApi.getGroups()
    // 处理可能的响应格式
    if (response && response.groups) {
      groups.value = response.groups
    } else if (Array.isArray(response)) {
      groups.value = response
    } else {
      groups.value = []
    }
  } catch (error: any) {
    console.error('Failed to load groups:', error)
    groups.value = []
  } finally {
    loadingGroups.value = false
  }
}

const loadUsers = async () => {
  try {
    const response = await userApi.getUsers()
    if (response && response.users) {
      availableUsers.value = response.users
    } else if (Array.isArray(response)) {
      availableUsers.value = response
    } else {
      availableUsers.value = []
    }
  } catch (error: any) {
    console.error('Failed to load users:', error)
    availableUsers.value = []
  }
}

const getMemberCount = (group: any) => {
  return getGroupMembers(group).length
}

const getPrimaryMembers = (group: any) => {
  const members = getGroupMembers(group)
  if (members.length === 0) return '无'
  return members.slice(0, 3).join(', ') + (members.length > 3 ? '...' : '')
}

const getGroupMembers = (group: any) => {
  if (!group.users || group.users === '') return []
  return group.users.split(',').filter((u: string) => u !== '')
}

const editGroup = (group: any) => {
  editingGroup.value = group
  groupForm.value = {
    name: group.name,
    gid: group.gid,
    members: getGroupMembers(group)
  }
  showCreateGroupModal.value = true
}

const manageMembers = (group: any) => {
  managingGroup.value = group
}

const deleteGroup = async (group: any) => {
  if (confirm(`确定要删除组 "${group.name}" 吗?`)) {
    try {
      await groupApi.deleteGroup(group.name)
      await loadGroups()
    } catch (error: any) {
      console.error('Failed to delete group:', error)
      alert('删除组失败: ' + error.message)
    }
  }
}

const saveGroup = async () => {
  try {
    saving.value = true

    if (editingGroup.value) {
      // 更新组
      await groupApi.updateGroup(editingGroup.value.name, {
        members: groupForm.value.members
      })
    } else {
      // 创建新组
      await groupApi.createGroup({
        name: groupForm.value.name,
        gid: groupForm.value.gid ? parseInt(groupForm.value.gid) : undefined,
        members: groupForm.value.members
      })
    }

    await loadGroups()
    closeModal()
  } catch (error: any) {
    console.error('Failed to save group:', error)
    alert('保存组失败: ' + error.message)
  } finally {
    saving.value = false
  }
}

const closeModal = () => {
  showCreateGroupModal.value = false
  editingGroup.value = null
  managingGroup.value = null
  groupForm.value = {
    name: '',
    gid: '',
    members: []
  }
}

const closeMemberModal = () => {
  managingGroup.value = null
}

const toggleMember = (username: string) => {
  const index = groupForm.value.members.indexOf(username)
  if (index > -1) {
    groupForm.value.members.splice(index, 1)
  } else {
    groupForm.value.members.push(username)
  }
}

const removeMember = (username: string) => {
  const index = groupForm.value.members.indexOf(username)
  if (index > -1) {
    groupForm.value.members.splice(index, 1)
  }
}

const isUserInGroup = (username: string, group: any) => {
  return getGroupMembers(group).includes(username)
}

const addGroupMember = async (group: any, username: string) => {
  try {
    const currentMembers = getGroupMembers(group)
    const newMembers = [...currentMembers, username]

    await groupApi.updateGroup(group.name, {
      members: newMembers
    })

    // 更新本地数据
    group.users = newMembers.join(',')
    await loadGroups()
  } catch (error: any) {
    console.error('Failed to add member:', error)
    alert('添加成员失败: ' + error.message)
  }
}

const removeGroupMember = async (group: any, username: string) => {
  try {
    const currentMembers = getGroupMembers(group)
    const newMembers = currentMembers.filter(m => m !== username)

    await groupApi.updateGroup(group.name, {
      members: newMembers
    })

    // 更新本地数据
    group.users = newMembers.join(',')
    await loadGroups()
  } catch (error: any) {
    console.error('Failed to remove member:', error)
    alert('移除成员失败: ' + error.message)
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadGroups()
  loadUsers()
})
</script>

<style scoped lang="scss">
.group-manager {
  .manager-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      color: #303133;
      font-size: 18px;
    }
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #909399;

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid #3b82f6;
    border-top-color: transparent;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.group-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 20px;
  transition: all 0.3s;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .group-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    padding-bottom: 15px;
    border-bottom: 1px solid #f0f0f0;

    .group-info {
      display: flex;
      align-items: center;
      gap: 12px;

      .group-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 40px;
        height: 40px;
        background: #eff6ff;
        border-radius: 8px;
        color: #3b82f6;
      }

      h4 {
        margin: 0;
        color: #303133;
        font-size: 16px;
        font-weight: 600;
      }

      .group-gid {
        margin: 4px 0 0;
        color: #909399;
        font-size: 12px;
      }
    }

    .group-actions {
      display: flex;
      gap: 8px;
    }
  }

  .group-details {
    margin-bottom: 15px;

    .detail-item {
      display: flex;
      justify-content: space-between;
      margin-bottom: 8px;
      font-size: 14px;

      .detail-label {
        color: #606266;
      }

      .detail-value {
        color: #303133;
        font-weight: 500;
      }
    }
  }

  .group-members-preview {
    .members-label {
      font-size: 12px;
      color: #606266;
      margin-bottom: 8px;
    }

    .members-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .member-chip {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        padding: 4px 12px;
        background: #f5f7fa;
        border-radius: 16px;
        font-size: 12px;
        color: #606266;
      }
    }
  }
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
  border: none;

  &.primary {
    background: #3b82f6;
    color: white;

    &:hover {
      background: #2563eb;
    }
  }

  &.secondary {
    background: white;
    border: 1px solid #e5e7eb;
    color: #374151;

    &:hover {
      background: #f9fafb;
      border-color: #d1d5db;
    }
  }
}

.icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #f3f4f6;
    color: #374151;
  }

  &.danger {
    &:hover {
      background: #fef2f2;
      color: #dc2626;
    }
  }
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
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;

  h3 {
    margin: 0;
    color: #303133;
    font-size: 18px;
  }
}

.group-form {
  padding: 20px;

  .form-group {
    margin-bottom: 20px;

    label {
      display: block;
      margin-bottom: 8px;
      font-size: 14px;
      font-weight: 500;
      color: #374151;
    }

    input {
      width: 100%;
      padding: 10px 12px;
      border: 1px solid #e5e7eb;
      border-radius: 6px;
      font-size: 14px;

      &:focus {
        outline: none;
        border-color: #3b82f6;
        box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
      }

      &:disabled {
        background: #f9fafb;
        cursor: not-allowed;
      }
    }

    .form-hint {
      font-size: 12px;
      color: #909399;
      margin-top: 4px;
    }
  }

  .members-selector {
    .selected-members {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
      margin-bottom: 12px;
      min-height: 40px;
      padding: 8px;
      border: 1px solid #e5e7eb;
      border-radius: 6px;
      background: #f9fafb;

      .member-tag {
        display: inline-flex;
        align-items: center;
        gap: 6px;
        padding: 4px 10px;
        background: #e5e7eb;
        border-radius: 16px;
        font-size: 12px;

        .remove-btn {
          background: none;
          border: none;
          color: #9ca3af;
          cursor: pointer;
          font-size: 14px;

          &:hover {
            color: #dc2626;
          }
        }
      }
    }

    .available-members {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;

      .member-option {
        padding: 6px 12px;
        border: 1px solid #e5e7eb;
        border-radius: 6px;
        font-size: 12px;
        cursor: pointer;
        background: white;

        &.selected {
          background: #3b82f6;
          color: white;
          border-color: #3b82f6;
        }

        &:hover:not(.selected) {
          background: #f3f4f6;
        }
      }
    }
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    margin-top: 20px;
  }
}

.member-management {
  padding: 20px;

  .current-members,
  .available-users {
    margin-bottom: 20px;

    h4 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 16px;
    }
  }

  .members-list,
  .users-list {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .member-item,
    .user-item {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 12px;
      background: white;
      border: 1px solid #e5e7eb;
      border-radius: 6px;

      &.in-group {
        background: #f0f9ff;
      }
    }
  }
}
</style>
