<template>
  <div class="group-management-panel">
    <div class="panel-toolbar">
      <div class="toolbar-left">
        <h3>用户组管理</h3>
        <p class="subtitle">管理用户组和权限模板</p>
      </div>
      <div class="toolbar-right">
        <button class="btn btn-primary" @click="showCreateGroupModal = true">
          <UserPlusIcon class="w-4 h-4" />
          创建用户组
        </button>
      </div>
    </div>

    <div class="groups-grid">
      <div
        v-for="group in groups"
        :key="group.id"
        class="group-card"
        :class="{ 'system-group': group.type === 'builtin' }"
      >
        <div class="group-header">
          <div class="group-icon">
            <UserGroupIcon class="w-6 h-6" />
          </div>
          <div class="group-info">
            <h4>{{ group.name }}</h4>
            <span class="group-type">{{ getGroupTypeLabel(group.type) }}</span>
          </div>
          <div class="group-actions" v-if="group.type !== 'builtin'">
            <button class="icon-btn" @click="editGroup(group)">
              <PencilIcon class="w-4 h-4" />
            </button>
            <button class="icon-btn danger" @click="deleteGroup(group)">
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <p class="group-description">{{ group.description || '无描述' }}</p>

        <div class="group-stats">
          <div class="stat">
            <span class="stat-value">{{ group.members.length }}</span>
            <span class="stat-label">成员</span>
          </div>
          <div class="stat">
            <span class="stat-value">GID {{ group.gid }}</span>
            <span class="stat-label">组ID</span>
          </div>
        </div>

        <div class="group-members">
          <h5>成员用户</h5>
          <div class="members-list">
            <span
              v-for="member in group.members.slice(0, 5)"
              :key="member"
              class="member-tag"
            >
              {{ member }}
            </span>
            <span v-if="group.members.length > 5" class="more-members">
              +{{ group.members.length - 5 }}
            </span>
          </div>
        </div>

        <div class="group-permissions">
          <h5>权限级别</h5>
          <div class="permissions-summary">
            <span class="permission-badge">{{ getGroupPermissionLevel(group) }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑用户组对话框 -->
    <div v-if="showCreateGroupModal || showEditGroupModal" class="modal-overlay" @click.self="closeGroupModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingGroup ? '编辑用户组' : '创建用户组' }}</h3>
          <button class="close-btn" @click="closeGroupModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit.prevent="saveGroup" class="group-form">
          <div class="form-group">
            <label>组名称 *</label>
            <input
              v-model="groupForm.name"
              type="text"
              required
              :disabled="editingGroup?.type === 'builtin'"
              placeholder="组名称"
            />
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="groupForm.description"
              placeholder="用户组描述"
              rows="3"
            ></textarea>
          </div>
          <div class="form-group">
            <label>权限模板</label>
            <select v-model="groupForm.permissionTemplate">
              <option value="">无模板</option>
              <option value="guest">访客权限</option>
              <option value="user">普通用户权限</option>
              <option value="admin">管理员权限</option>
              <option value="superadmin">超级管理员权限</option>
            </select>
          </div>
          <div class="form-group">
            <label>组成员</label>
            <div class="members-selection">
              <label
                v-for="user in users"
                :key="user.id"
                class="member-checkbox"
              >
                <input
                  type="checkbox"
                  v-model="groupForm.members"
                  :value="user.username"
                />
                <span>{{ user.username }}</span>
                <small>{{ user.fullName || user.email || '' }}</small>
              </label>
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn btn-secondary" @click="closeGroupModal">
              取消
            </button>
            <button type="submit" class="btn btn-primary">
              {{ editingGroup ? '保存更改' : '创建用户组' }}
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
  UserGroupIcon,
  UserPlusIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import type { Group, User } from '@/types/permissions'

const props = defineProps<{
  groups: Group[]
  users: User[]
}>()

const emit = defineEmits<{
  'group-added': [group: Group]
  'group-edited': [group: Group]
  'group-deleted': [groupId: string]
}>()

const showCreateGroupModal = ref(false)
const showEditGroupModal = ref(false)
const editingGroup = ref<Group | null>(null)

const groupForm = ref({
  name: '',
  description: '',
  permissionTemplate: '',
  members: [] as string[]
})

const getGroupTypeLabel = (type: string) => {
  const labels = {
    builtin: '系统组',
    user: '用户组',
    system: '系统组'
  }
  return labels[type as keyof typeof labels] || type
}

const getGroupPermissionLevel = (group: Group) => {
  // 简化的权限级别判断
  const perms = group.permissions
  if (perms.system?.admin === 'owner') return '超级管理员'
  if (perms.system?.admin === 'admin') return '管理员'
  if (perms.file?.write === 'write') return '读写用户'
  return '只读用户'
}

const editGroup = (group: Group) => {
  editingGroup.value = group
  groupForm.value = {
    name: group.name,
    description: group.description || '',
    permissionTemplate: '',
    members: [...group.members]
  }
  showEditGroupModal.value = true
}

const deleteGroup = (group: Group) => {
  if (confirm(`确定要删除用户组 "${group.name}" 吗？`)) {
    emit('group-deleted', group.id)
  }
}

const saveGroup = () => {
  const group: Group = {
    id: editingGroup.value?.id || Date.now().toString(),
    name: groupForm.value.name,
    description: groupForm.value.description,
    gid: editingGroup.value?.gid || Date.now(),
    type: 'user',
    members: groupForm.value.members,
    createdAt: editingGroup.value?.createdAt || new Date(),
    permissions: {}
  }

  if (editingGroup.value) {
    emit('group-edited', group)
  } else {
    emit('group-added', group)
  }

  closeGroupModal()
}

const closeGroupModal = () => {
  showCreateGroupModal.value = false
  showEditGroupModal.value = false
  editingGroup.value = null
  groupForm.value = {
    name: '',
    description: '',
    permissionTemplate: '',
    members: []
  }
}
</script>

<style scoped>
.group-management-panel {
  width: 100%;
}

.panel-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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
  border: 2px solid transparent;
  transition: all 0.2s;
}

.group-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.system-group {
  border-color: #fbbf24;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 12px;
}

.group-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.group-info {
  flex: 1;
  margin-left: 12px;
}

.group-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.group-type {
  font-size: 12px;
  color: #6b7280;
  background: white;
  padding: 2px 6px;
  border-radius: 4px;
}

.group-actions {
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
  background: white;
  color: #1f2937;
}

.icon-btn.danger:hover {
  background: #fee2e2;
  color: #ef4444;
}

.group-description {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 16px;
  min-height: 40px;
}

.group-stats {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.stat {
  text-align: center;
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.group-members,
.group-permissions {
  margin-bottom: 12px;
}

.group-members h5,
.group-permissions h5 {
  font-size: 13px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.members-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.member-tag {
  padding: 4px 8px;
  background: #f3f4f6;
  border-radius: 6px;
  font-size: 12px;
  color: #6b7280;
}

.more-members {
  padding: 4px 8px;
  background: #e5e7eb;
  border-radius: 6px;
  font-size: 12px;
  color: #6b7280;
}

.permission-badge {
  padding: 6px 12px;
  background: #dbeafe;
  color: #3b82f6;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
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
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
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

.group-form {
  padding: 24px;
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

.form-group input,
.form-group select,
.form-group textarea {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.members-selection {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 8px;
}

.member-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
}

.member-checkbox:hover {
  background: #f9fafb;
}

.member-checkbox small {
  color: #6b7280;
  font-size: 12px;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
}
</style>