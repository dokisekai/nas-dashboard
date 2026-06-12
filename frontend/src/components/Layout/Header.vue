<template>
  <header class="h-16 bg-gray-900/50 backdrop-blur-xl border-b border-gray-800 flex items-center justify-between px-6">
    <div>
      <h1 class="text-xl font-semibold text-white">{{ title }}</h1>
      <p class="text-gray-500 text-sm">{{ subtitle }}</p>
    </div>
    <div class="flex items-center gap-4">
      <!-- 系统状态 -->
      <div class="flex items-center gap-2 px-4 py-2 bg-gray-800/50 rounded-xl border border-gray-700/50">
        <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
        <span class="text-gray-300 text-sm">系统运行中</span>
      </div>

      <!-- 用户信息 -->
      <div class="flex items-center gap-3 px-4 py-2 bg-gray-800/50 rounded-xl border border-gray-700/50">
        <div class="w-8 h-8 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
        </div>
        <span class="text-white text-sm font-medium">{{ authStore.user?.username || 'Admin' }}</span>
      </div>

      <!-- 退出按钮 -->
      <button
        @click="logout"
        class="flex items-center gap-2 px-4 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 hover:text-red-300 rounded-xl transition-all duration-200 border border-red-500/20"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
        <span>退出</span>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const authStore = useAuthStore()

const title = computed(() => {
  const titles: Record<string, string> = {
    Dashboard: '仪表盘',
    Monitor: '系统监控',
    MonitorMemory: '内存监控',
    MonitorDisk: '磁盘监控',
    MonitorNetwork: '网络监控',
    Storage: '存储管理',
    Services: '服务管理',
    Users: '用户管理',
  }
  return titles[route.name as string] || 'NAS 面板'
})

const subtitle = computed(() => {
  const subtitles: Record<string, string> = {
    Dashboard: '系统概览和状态',
    Monitor: '实时系统资源监控',
    MonitorMemory: '内存使用情况',
    MonitorDisk: '磁盘IO和使用率',
    MonitorNetwork: '网络流量统计',
    Storage: '磁盘和存储管理',
    Services: '系统服务和容器',
    Users: '用户和权限管理',
  }
  return subtitles[route.name as string] || ''
})

const logout = () => {
  authStore.clearToken()
  window.location.href = '/login'
}
</script>
