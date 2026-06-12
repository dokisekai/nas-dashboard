<template>
  <aside class="w-64 bg-gray-900/50 backdrop-blur-xl border-r border-gray-800 flex flex-col">
    <!-- Logo -->
    <div class="p-6 border-b border-gray-800">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-indigo-500/30">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
          </svg>
        </div>
        <span class="text-xl font-bold text-white">NAS 面板</span>
      </div>
    </div>

    <!-- 导航菜单 -->
    <nav class="flex-1 p-4 space-y-2">
      <router-link
        v-for="item in menuItems"
        :key="item.name"
        :to="item.path"
        class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group"
        :class="isActive(item.path, item.children) ? 'bg-gradient-to-r from-indigo-500 to-purple-600 text-white shadow-lg shadow-indigo-500/30' : 'text-gray-400 hover:bg-gray-800/50 hover:text-white'"
      >
        <component :is="item.icon" class="w-5 h-5 shrink-0" />
        <span class="font-medium">{{ item.label }}</span>
        <svg v-if="item.children" class="w-4 h-4 ml-auto transition-transform duration-200" :class="{ 'rotate-180': expanded[item.name] }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </router-link>

      <!-- 监控子菜单 -->
      <template v-if="expanded['Monitor']">
        <router-link
          v-for="child in monitorChildren"
          :key="child.name"
          :to="child.path"
          class="flex items-center gap-3 px-4 py-2.5 pl-14 rounded-xl text-sm transition-all duration-200 group"
          :class="isActive(child.path) ? 'bg-gray-800/50 text-white' : 'text-gray-500 hover:bg-gray-800/30 hover:text-gray-300'"
        >
          <div class="w-1.5 h-1.5 rounded-full" :class="isActive(child.path) ? 'bg-indigo-400' : 'bg-gray-600 group-hover:bg-gray-500'"></div>
          {{ child.label }}
        </router-link>
      </template>
    </nav>

    <!-- 底部状态 -->
    <div class="p-4 border-t border-gray-800">
      <div class="flex items-center gap-3 px-4 py-3 bg-gray-800/50 rounded-xl">
        <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
        <div>
          <p class="text-white text-sm font-medium">系统正常</p>
          <p class="text-gray-500 text-xs">运行时间: 2天 5小时</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  ChartBarIcon,
  CircleStackIcon,
  CogIcon,
  UserGroupIcon,
  HomeIcon,
} from '@heroicons/vue/24/outline'

const route = useRoute()

const expanded = ref<Record<string, boolean>>({
  Monitor: false,
})

const menuItems = [
  { name: 'Dashboard', label: '仪表盘', path: '/dashboard', icon: HomeIcon },
  { name: 'Monitor', label: '系统监控', path: '/monitor', icon: ChartBarIcon, children: true },
  { name: 'Storage', label: '存储管理', path: '/storage', icon: CircleStackIcon },
  { name: 'Services', label: '服务管理', path: '/services', icon: CogIcon },
  { name: 'Users', label: '用户管理', path: '/users', icon: UserGroupIcon },
]

const monitorChildren = [
  { name: 'MonitorCPU', label: 'CPU', path: '/monitor' },
  { name: 'MonitorMemory', label: '内存', path: '/monitor/memory' },
  { name: 'MonitorDisk', label: '磁盘', path: '/monitor/disk' },
  { name: 'MonitorNetwork', label: '网络', path: '/monitor/network' },
]

const isActive = (path: string, hasChildren?: boolean) => {
  if (hasChildren && route.path.startsWith('/monitor')) {
    expanded.value['Monitor'] = true
    return true
  }
  return route.path === path || route.path.startsWith(path)
}
</script>
