<template>
  <div class="simple-group-panel">
    <!-- 工具栏 -->
    <div class="toolbar">
      <div class="toolbar-left">
        <h3>用户组管理</h3>
        <p class="subtitle">创建用户组，简化权限管理</p>
      </div>
      <div class="toolbar-right">
        <button class="btn btn-primary" @click="showAddGroupModal = true">
          <PlusIcon class="w-4 h-4" />
          创建用户组
        </button>
      </div>
    </div>

    <!-- 用户组列表 -->
    <div class="groups-grid">
      <div
        v-for="group in groups"
        :key="group.id"
        class="group-card"
        :class="{ 'system-group': group.type === 'admin' }"
      >
        <div class="group-header">
          <div class="group-icon" :class="`group-${group.type}`">
            <UserGroupIcon class="w-6 h-6" />
          </div>
          <div class="group-info">
            <h4>{{ group.name }}</h4>
            <span class="group-type-badge">{{ getGroupTypeLabel(group.type) }}</span>
          </div>
          <div class="group-actions" v-if="group.type !== 'admin'">
            <button class="icon-btn" @click="editGroup(group)">
              <PencilIcon class="w-4 h-4" />
            </button>
            <button class="icon-btn danger" @click="deleteGroup(group)">
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>

        <p class="group-description">{{ group.description }}</p>

        <div class="group-features">
          <h5>权限范围</h5>
          <div class="features-list">
            <div
              v-for="permission in group.permissions"
              :key="permission"
              class="feature-item"
            >
              <CheckIcon class="w-4 h-4" />
              <span>{{ permission }}</span>
            </div>
          </div>
        </div>

        <div class="group-members">
          <div class="members-header">
            <UsersIcon class="w-4 h-4" />
            <span>{{ group.memberCount }} 个成员</span>
          </div>
          <div class="members-preview">
            <span
              v-for="(member, index) in getGroupMembers(group.id)"
              :key="member"
              class="member-avatar"
            >
              {{ member.charAt(0).toUpperCase() }}
            </span>
            <span v-if="group.memberCount > 3" class="more-members">
              +{{ group.memberCount - 3 }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑用户组对话框 -->
    <div v-if="showAddGroupModal" class="modal-overlay" @click.self="closeGroupModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingGroup ? '编辑用户组' : '创建用户组' }}</h3>
          <button class="close-btn" @click="closeGroupModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="saveGroup" class="group-form">
          <div class="form-group">
            <label>用户组名称 *</label>
            <input
              v-model="groupForm.name"
              type="text"
              required
              :disabled="editingGroup?.type === 'admin'"
              placeholder="例如: 开发团队"
            />
            <small>用户组的名称，便于识别</small>
          </div>

          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="groupForm.description"
              placeholder="这个用户组的用途说明"
              rows="3"
            ></textarea>
            <small>描述这个用户组的权限范围</small>
          </div>

          <div class="form-group">
            <label>权限级别 *</label>
            <div class="permission-levels">
              <label
                v-for="level in permissionLevels"
                :key="level.value"
                :class="['level-option', { selected: groupForm.permissionLevel === level.value }]"
              >
                <input type="radio" v-model="groupForm.permissionLevel" :value="level.value" />
                <div class="level-card">
                  <div class="level-icon" :style="{ background: level.color }">
                    <component :is="level.icon" class="w-6 h-6" />
                  </div>
                  <div class="level-info">
                    <span class="level-name">{{ level.label }}</span>
                    <span class="level-desc">{{ level.description }}</span>
                  </div>
                </div>
              </label>
            </div>
          </div>

          <div class="form-group">
            <label>组成员（可选）</label>
            <div class="members-selection">
              <label
                v-for="user in users"
                :key="user.id"
                class="member-checkbox"
              >
                <input type="checkbox" v-model="groupForm.members" :value="user.name" />
                <span>{{ user.displayName || user.name }}</span>
                <small>{{ getRoleLabel(user.role) }}</small>
              </label>
            </div>
            <small>选择要加入这个用户组的用户</small>
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
  PlusIcon,
  UserGroupIcon,
  UsersIcon,
  PencilIcon,
  TrashIcon,
  CheckIcon,
  XMarkIcon,
  ShieldCheckIcon,
  UserIcon,
  EyeIcon
} from '@heroicons/vue/24/outline'

const props = defineProps<{
  groups: any[]
  users: any[]
}>()

const emit = defineEmits<{
  'group-created': [group: any]
  'group-updated': [group: any]
  'group-deleted': [groupId: string]
}>()

const showAddGroupModal = ref(false)
const editingGroup = ref<any>(null)

const groupForm = ref({
  name: '',
  description: '',
  permissionLevel: 'user',
  members: []
})

// 简化的权限级别
const permissionLevels = [
  {
    value: 'admin',
    label: '管理员级别',
    description: '可以管理所有功能和设置',
    icon: ShieldCheckIcon,
    color: 'linear-gradient(135deg, #f59e0b 0%, #d97706 100%)'
  },
  {
    value: 'user',
    label: '普通用户级别',
    description: '可以访问文件和基本功能',
    icon: UserIcon,
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)'
  },
  {
    value: 'guest',
    label: '访客级别',
    description: '只能查看，不能修改文件',
    icon: EyeIcon,
    color: 'linear-gradient(135deg, #9ca3af 0%, #6b7280 100%)'
  }
]

const getGroupTypeLabel = (type: string) => {
  const labels: Record<string, string> = {
    admin: '系统组',
    user: '用户组',
    guest: '访客组'
  }
  return labels[type] || type
}

const getGroupMembers = (groupId: string) => {
  // 模拟获取组成员前3个
  const members = ['admin', 'john', 'guest']
  return members.slice(0, 3)
}

const getRoleLabel = (role: string) => {
  const labels: Record<string, string> = {
    admin: '管理员',
    user: '普通用户',
    guest: '访客'
  }
  return labels[role] || role
}

const editGroup = (group: any) => {
  editingGroup.value = group
  groupForm.value = {
    name: group.name,
    description: group.description || '',
    permissionLevel: group.type,
    members: []
  }
  showAddGroupModal.value = true
}

const deleteGroup = (group: any) => {
  if (confirm(`确定要删除用户组 "${group.name}" 吗？`)) {
    emit('group-deleted', group.id)
  }
}

const saveGroup = () => {
  const group = {
    id: editingGroup.value?.id || Date.now().toString(),
    name: groupForm.value.name,
    description: groupForm.value.description,
    type: groupForm.value.permissionLevel,
    memberCount: groupForm.value.members.length,
    permissions: getPermissionsForLevel(groupForm.value.permissionLevel)
  }

  if (editingGroup.value) {
    emit('group-updated', group)
  } else {
    emit('group-created', group)
  }

  closeGroupModal()
}

const getPermissionsForLevel = (level: string) => {
  const permissions: Record<string, string[]> = {
    admin: ['所有权限', '用户管理', '系统设置', '文件管理'],
    user: ['文件访问', '个人设置', '基本功能'],
    guest: ['只读访问', '查看文件']
  }
  return permissions[level] || ['基本访问']
}

const closeGroupModal = () => {
  showAddGroupModal.value = false
  editingGroup.value = null
  groupForm.value = {
    name: '',
    description: '',
    permissionLevel: 'user',
    members: []
  }
}
</script>

<style scoped>
.simple-group-panel {
  width: 100%;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
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

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.group-card {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 2px solid transparent;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.group-card:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.system-group {
  border-color: #fbbf24;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 100%);
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
}

.group-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.group-admin {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
}

.group-user {
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
}

.group-guest {
  background: linear-gradient(135deg, #9ca3af 0%, #6b7280 100%);
}

.group-info {
  flex: 1;
  margin-left: 12px;
}

.group-info h4 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.group-type-badge {
  font-size: 12px;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.8);
  padding: 4px 8px;
  border-radius: 6px;
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
  background: rgba(255, 255, 255, 0.8);
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
  color: #374151;
  line-height: 1.5;
}

.group-features h5 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.features-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 6px;
  font-size: 13px;
  color: #374151;
}

.feature-item svg {
  color: #10b981;
  flex-shrink: 0;
}

.group-members {
  margin-top: auto;
  padding-top: 16px;
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.members-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 12px;
}

.members-preview {
  display: flex;
  align-items: center;
  gap: 8px;
}

.member-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
}

.more-members {
  font-size: 12px;
  color: #6b7280;
  padding: 4px 8px;
  background: rgba(255, 255, 255, 0.8);
  border-radius: 6px;
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

.group-form {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
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
.form-group textarea,
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

.permission-levels {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.level-option {
  cursor: pointer;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  padding: 16px;
  transition: all 0.2s;
}

.level-option.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.level-option:hover:not(.selected) {
  border-color: #d1d5db;
}

.level-card {
  display: flex;
  align-items: center;
  gap: 12px;
}

.level-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.level-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.level-name {
  font-weight: 500;
  color: #1f2937;
  font-size: 14px;
}

.level-desc {
  font-size: 12px;
  color: #6b7280;
}

.members-selection {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
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

@media (max-width: 768px) {
  .groups-grid {
    grid-template-columns: 1fr;
  }
}
</style>
