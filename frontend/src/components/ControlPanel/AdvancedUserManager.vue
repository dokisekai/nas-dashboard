<template>
  <div class="advanced-user-manager">
    <!-- 顶部标签页 -->
    <div class="tabs-header">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.label }}
        <span v-if="getTabCount(tab.id) > 0" class="tab-count">
          {{ getTabCount(tab.id) }}
        </span>
      </button>
    </div>

    <!-- 用户管理 -->
    <div v-if="activeTab === 'users'" class="tab-content">
      <UserManagementPanel
        :users="users"
        :groups="groups"
        @user-added="handleUserAdded"
        @user-edited="handleUserEdited"
        @user-deleted="handleUserDeleted"
      />
    </div>

    <!-- 用户组管理 -->
    <div v-if="activeTab === 'groups'" class="tab-content">
      <GroupManagementPanel
        :groups="groups"
        :users="users"
        @group-added="handleGroupAdded"
        @group-edited="handleGroupEdited"
        @group-deleted="handleGroupDeleted"
      />
    </div>

    <!-- 权限模板 -->
    <div v-if="activeTab === 'templates'" class="tab-content">
      <PermissionTemplatesPanel
        :templates="permissionTemplates"
        @template-applied="handleTemplateApplied"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  UserGroupIcon,
  UsersIcon,
  ShieldCheckIcon
} from '@heroicons/vue/24/outline'
import UserManagementPanel from './UserManagementPanel.vue'
import GroupManagementPanel from './GroupManagementPanel.vue'
import PermissionTemplatesPanel from './PermissionTemplatesPanel.vue'
import type { User, Group } from '@/types/permissions'

const activeTab = ref('users')

const tabs = [
  { id: 'users', label: '用户', icon: UsersIcon },
  { id: 'groups', label: '用户组', icon: UserGroupIcon },
  { id: 'templates', label: '权限模板', icon: ShieldCheckIcon }
]

// 状态数据
const users = ref<User[]>([])
const groups = ref<Group[]>([])
const permissionTemplates = ref([])

// 获取标签页计数
const getTabCount = (tabId: string) => {
  switch (tabId) {
    case 'users': return users.value.filter(u => u.status === 'active').length
    case 'groups': return groups.value.length
    case 'templates': return permissionTemplates.value.length
    default: return 0
  }
}

// 事件处理
const handleUserAdded = (user: User) => {
  users.value.push(user)
}

const handleUserEdited = (user: User) => {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index !== -1) {
    users.value[index] = user
  }
}

const handleUserDeleted = (userId: string) => {
  users.value = users.value.filter(u => u.id !== userId)
}

const handleGroupAdded = (group: Group) => {
  groups.value.push(group)
}

const handleGroupEdited = (group: Group) => {
  const index = groups.value.findIndex(g => g.id === group.id)
  if (index !== -1) {
    groups.value[index] = group
  }
}

const handleGroupDeleted = (groupId: string) => {
  groups.value = groups.value.filter(g => g.id !== groupId)
}

const handleTemplateApplied = (data: any) => {
  console.log('Permission template applied:', data)
}

// 加载数据
const loadData = async () => {
  try {
    // 模拟数据加载
    users.value = [
      {
        id: '1',
        username: 'admin',
        fullName: '系统管理员',
        email: 'admin@nas.local',
        uid: 1000,
        primaryGroup: 'administrators',
        groups: ['administrators'],
        role: 'superadmin',
        status: 'active',
        homeDirectory: '/home/admin',
        shell: '/bin/bash',
        createdAt: new Date('2024-01-01'),
        permissions: {}
      },
      {
        id: '2',
        username: 'john',
        fullName: 'John Doe',
        email: 'john@example.com',
        uid: 1001,
        primaryGroup: 'users',
        groups: ['users'],
        role: 'user',
        status: 'active',
        homeDirectory: '/home/john',
        shell: '/bin/bash',
        lastLogin: new Date(),
        createdAt: new Date('2024-01-15'),
        storageQuota: {
          enabled: true,
          size: 1073741824, // 1GB
          used: 536870912   // 512MB
        }
      },
      {
        id: '3',
        username: 'guest',
        fullName: 'Guest User',
        uid: 1002,
        primaryGroup: 'guests',
        groups: ['guests'],
        role: 'guest',
        status: 'active',
        homeDirectory: '/home/guest',
        shell: '/bin/bash',
        createdAt: new Date('2024-02-01')
      }
    ]

    groups.value = [
      {
        id: 'administrators',
        name: 'administrators',
        gid: 999,
        description: '系统管理员组',
        type: 'builtin',
        members: ['admin'],
        createdAt: new Date('2024-01-01'),
        permissions: {}
      },
      {
        id: 'users',
        name: 'users',
        gid: 100,
        description: '普通用户组',
        type: 'builtin',
        members: ['john'],
        createdAt: new Date('2024-01-01'),
        permissions: {}
      },
      {
        id: 'guests',
        name: 'guests',
        gid: 998,
        description: '访客组',
        type: 'builtin',
        members: ['guest'],
        createdAt: new Date('2024-01-01'),
        permissions: {}
      }
    ]
  } catch (error) {
    console.error('Failed to load data:', error)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.advanced-user-manager {
  max-width: 1200px;
  margin: 0 auto;
}

.tabs-header {
  display: flex;
  gap: 4px;
  border-bottom: 2px solid #e5e7eb;
  margin-bottom: 24px;
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
  transition: all 0.2s;
  margin-bottom: -2px;
  position: relative;
}

.tab-btn:hover {
  color: #1f2937;
  background: #f9fafb;
}

.tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

.tab-count {
  background: #3b82f6;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 600;
  min-width: 20px;
  text-align: center;
}

.tab-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>