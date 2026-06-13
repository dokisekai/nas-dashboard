<template>
  <div class="services-list-component">
    <div class="services-header">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索服务..."
        class="services-search"
      />
      <button class="cp-btn cp-btn-secondary" @click="refreshServices">
        <ArrowPathIcon class="w-4 h-4" />
        刷新
      </button>
    </div>

    <div class="services-list">
      <div
        v-for="service in filteredServices"
        :key="service.name"
        class="service-item"
      >
        <div class="service-info">
          <div class="service-status" :class="service.status.toLowerCase()"></div>
          <div class="service-details">
            <h4>{{ service.name }}</h4>
            <p>{{ service.description }}</p>
          </div>
        </div>

        <div class="service-meta">
          <span class="service-enabled" :class="{ active: service.enabled }">
            {{ service.enabled ? '已启用' : '已禁用' }}
          </span>
        </div>

        <div class="service-actions">
          <button
            v-if="service.status === 'running'"
            class="cp-btn cp-btn-warning cp-btn-sm"
            @click="stopService(service)"
          >
            <StopIcon class="w-4 h-4" />
            停止
          </button>
          <button
            v-else
            class="cp-btn cp-btn-primary cp-btn-sm"
            @click="startService(service)"
          >
            <PlayIcon class="w-4 h-4" />
            启动
          </button>
          <button class="cp-btn cp-btn-ghost cp-btn-sm" @click="restartService(service)">
            <ArrowPathIcon class="w-4 h-4" />
            重启
          </button>
          <button
            class="cp-btn cp-btn-ghost cp-btn-sm"
            @click="toggleService(service)"
          >
            <PowerIcon class="w-4 h-4" />
            {{ service.enabled ? '禁用' : '启用' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  ArrowPathIcon,
  StopIcon,
  PlayIcon,
  PowerIcon
} from '@heroicons/vue/24/outline'

// 模拟服务数据
const services = ref([
  {
    name: 'nginx',
    description: '高性能Web服务器',
    status: 'running',
    enabled: true
  },
  {
    name: 'mysql',
    description: 'MySQL数据库服务器',
    status: 'running',
    enabled: true
  },
  {
    name: 'docker',
    description: 'Docker容器服务',
    status: 'running',
    enabled: true
  },
  {
    name: 'ssh',
    description: 'SSH远程登录服务',
    status: 'running',
    enabled: true
  },
  {
    name: 'smb',
    description: 'SMB文件共享服务',
    status: 'stopped',
    enabled: false
  }
])

const searchQuery = ref('')

const filteredServices = computed(() => {
  if (!searchQuery.value) return services.value
  return services.value.filter(s =>
    s.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    s.description.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

const refreshServices = () => {
  // 刷新服务列表
  console.log('Refreshing services...')
}

const startService = (service: any) => {
  console.log('Starting service:', service.name)
  service.status = 'running'
}

const stopService = (service: any) => {
  console.log('Stopping service:', service.name)
  service.status = 'stopped'
}

const restartService = (service: any) => {
  console.log('Restarting service:', service.name)
  service.status = 'running'
}

const toggleService = (service: any) => {
  service.enabled = !service.enabled
  if (service.enabled) {
    service.status = 'running'
  } else {
    service.status = 'stopped'
  }
}
</script>

<style scoped>
.services-list-component {
  width: 100%;
}

.services-header {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.services-search {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
}

.services-search:focus {
  outline: none;
  border-color: #3b82f6;
}

.services-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.service-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  gap: 20px;
}

.service-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.service-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #9ca3af;
}

.service-status.running {
  background: #22c55e;
}

.service-status.stopped {
  background: #ef4444;
}

.service-details h4 {
  margin: 0 0 4px 0;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
}

.service-details p {
  margin: 0;
  font-size: 12px;
  color: #6b7280;
}

.service-meta {
  display: flex;
  align-items: center;
}

.service-enabled {
  font-size: 12px;
  color: #9ca3af;
  padding: 4px 8px;
  border-radius: 4px;
  background: #f3f4f6;
}

.service-enabled.active {
  color: #22c55e;
  background: #dcfce7;
}

.service-actions {
  display: flex;
  gap: 6px;
}

.cp-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  background: white;
  color: #374151;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.cp-btn:hover {
  background: #f9fafb;
  border-color: #d1d5db;
}

.cp-btn-primary {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.cp-btn-primary:hover {
  background: #2563eb;
  border-color: #2563eb;
}

.cp-btn-warning {
  background: #f59e0b;
  color: white;
  border-color: #f59e0b;
}

.cp-btn-warning:hover {
  background: #d97706;
  border-color: #d97706;
}

.cp-btn-ghost {
  background: transparent;
  border-color: #e5e7eb;
}

.cp-btn-ghost:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
}

.cp-btn-sm {
  padding: 4px 8px;
  font-size: 11px;
}
</style>
