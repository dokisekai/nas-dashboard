<template>
  <div class="simple-user-manager">
    <!-- 顶部标签页 -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.label }}
        <span v-if="getCount(tab.id)" class="badge">{{ getCount(tab.id) }}</span>
      </button>
    </div>

    <!-- 用户管理 -->
    <div v-if="activeTab === 'users'" class="tab-content">
      <SimpleUserPanel
        :users="users"
        @user-created="handleUserCreated"
        @user-updated="handleUserUpdated"
        @user-deleted="handleUserDeleted"
      />
    </div>

    <!-- 用户组管理 -->
    <div v-if="activeTab === 'groups'" class="tab-content">
      <SimpleGroupPanel
        :groups="groups"
        :users="users"
        @group-created="handleGroupCreated"
        @group-updated="handleGroupUpdated"
        @group-deleted="handleGroupDeleted"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { UsersIcon, UserGroupIcon } from '@heroicons/vue/24/outline'
import SimpleUserPanel from './SimpleUserPanel.vue'
import SimpleGroupPanel from './SimpleGroupPanel.vue'

const activeTab = ref('users')

const tabs = [
  { id: 'users', label: '用户', icon: UsersIcon },
  { id: 'groups', label: '用户组', icon: UserGroupIcon }
]

// 简化的用户数据结构
const users = ref([
  {
    id: '1',
    name: 'admin',
    displayName: '管理员',
    role: 'admin',
    groups: ['管理员组'],
    storageUsed: '2.3 GB',
    storageLimit: '无限制',
    status: 'active',
    lastLogin: '2小时前'
  },
  {
    id: '2',
    name: 'john',
    displayName: '张三',
    role: 'user',
    groups: ['普通用户组'],
    storageUsed: '15.6 GB',
    storageLimit: '100 GB',
    status: 'active',
    lastLogin: '1天前'
  }
])

const groups = ref([
  {
    id: '1',
    name: '管理员组',
    description: '可以管理所有系统功能',
    type: 'admin',
    memberCount: 1,
    permissions: ['所有权限']
  },
  {
    id: '2',
    name: '普通用户组',
    description: '可以访问文件和基本功能',
    type: 'user',
    memberCount: 3,
    permissions: ['文件访问', '个人设置']
  },
  {
    id: '3',
    name: '访客组',
    description: '只能查看，不能修改',
    type: 'guest',
    memberCount: 2,
    permissions: ['只读访问']
  }
])

const getCount = (tabId: string) => {
  if (tabId === 'users') return users.value.length
  if (tabId === 'groups') return groups.value.length
  return 0
}

const handleUserCreated = (user: any) => {
  users.value.push(user)
}

const handleUserUpdated = (user: any) => {
  const index = users.value.findIndex(u => u.id === user.id)
  if (index !== -1) {
    users.value[index] = user
  }
}

const handleUserDeleted = (userId: string) => {
  users.value = users.value.filter(u => u.id !== userId)
}

const handleGroupCreated = (group: any) => {
  groups.value.push(group)
}

const handleGroupUpdated = (group: any) => {
  const index = groups.value.findIndex(g => g.id === group.id)
  if (index !== -1) {
    groups.value[index] = group
  }
}

const handleGroupDeleted = (groupId: string) => {
  groups.value = groups.value.filter(g => g.id !== groupId)
}

onMounted(() => {
  // 初始化数据
})
</script>

<style scoped>
.simple-user-manager {
  width: 100%;
}

.tabs {
  display: flex;
  gap: 8px;
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
}

.tab-btn:hover {
  color: #1f2937;
}

.tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

.badge {
  background: #3b82f6;
  color: white;
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 12px;
  font-weight: 600;
  min-width: 18px;
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