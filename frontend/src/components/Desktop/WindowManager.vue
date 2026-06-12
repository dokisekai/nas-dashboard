<template>
  <div class="window-manager">
    <!-- 窗口标签栏 -->
    <div class="window-tabs" v-if="windowGroups.length > 1">
      <div
        v-for="(group, index) in windowGroups"
        :key="group.id"
        class="window-tab"
        :class="{ active: activeGroupId === group.id }"
        @click="activateGroup(group.id)"
      >
        <div class="tab-icon">
          <component :is="getGroupIcon(group)" class="w-4 h-4" />
        </div>
        <div class="tab-title">{{ group.title }}</div>
        <button @click.stop="closeGroup(group.id)" class="tab-close">
          <XMarkIcon class="w-3 h-3" />
        </button>
      </div>
    </div>

    <!-- 窗口容器 -->
    <div class="windows-container">
      <DesktopWindow
        v-for="window in activeWindows"
        :key="window.id"
        :window="window"
        @focus="focusWindow"
        @close="closeWindow"
        @minimize="minimizeWindow"
        @maximize="maximizeWindow"
      />
    </div>

    <!-- 最小化窗口指示器 -->
    <div class="minimized-indicator" v-if="minimizedWindows.length > 0">
      <div class="minimized-count">{{ minimizedWindows.length }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import DesktopWindow from './DesktopWindow.vue'
import type { Window } from '../../types/desktop'

interface Props {
  windows: Window[]
}

interface Emits {
  (e: 'focus', windowId: string): void
  (e: 'close', windowId: string): void
  (e: 'minimize', windowId: string): void
  (e: 'maximize', windowId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const activeGroupId = ref('all')

// 窗口分组
const windowGroups = computed(() => {
  const groups = [{ id: 'all', title: '全部窗口', appId: null }]

  const appGroups = new Map()
  props.windows.forEach(window => {
    if (!appGroups.has(window.appId)) {
      appGroups.set(window.appId, {
        id: window.appId,
        title: window.title,
        appId: window.appId
      })
    }
  })

  return [...groups, ...Array.from(appGroups.values())]
})

// 活动窗口
const activeWindows = computed(() => {
  if (activeGroupId.value === 'all') {
    return props.windows.filter(w => !w.minimized)
  }
  return props.windows.filter(w => w.appId === activeGroupId.value && !w.minimized)
})

// 最小化的窗口
const minimizedWindows = computed(() => {
  return props.windows.filter(w => w.minimized)
})

const focusWindow = (windowId: string) => {
  emit('focus', windowId)
}

const closeWindow = (windowId: string) => {
  emit('close', windowId)
}

const minimizeWindow = (windowId: string) => {
  emit('minimize', windowId)
}

const maximizeWindow = (windowId: string) => {
  emit('maximize', windowId)
}

const activateGroup = (groupId: string) => {
  activeGroupId.value = groupId
}

const closeGroup = (groupId: string) => {
  if (groupId === 'all') {
    props.windows.forEach(window => {
      emit('close', window.id)
    })
  } else {
    const groupWindows = props.windows.filter(w => w.appId === groupId)
    groupWindows.forEach(window => {
      emit('close', window.id)
    })
  }
}

const getGroupIcon = (group: any) => {
  // 返回分组图标
  return 'div'
}
</script>

<style scoped>
.window-manager {
  width: 100%;
  height: 100%;
  position: relative;
}

.window-tabs {
  position: absolute;
  top: 8px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 4px;
  padding: 4px;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  z-index: 10000;
}

.window-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 120px;
}

.window-tab:hover {
  background: rgba(255, 255, 255, 0.1);
}

.window-tab.active {
  background: rgba(255, 255, 255, 0.2);
}

.tab-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.tab-title {
  flex: 1;
  font-size: 12px;
  color: white;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tab-close {
  width: 20px;
  height: 20px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-close:hover {
  background: rgba(255, 255, 255, 0.2);
  color: white;
}

.windows-container {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.minimized-indicator {
  position: absolute;
  bottom: 80px;
  right: 20px;
  width: 48px;
  height: 48px;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(10px);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 9999;
}

.minimized-indicator:hover {
  transform: scale(1.1);
  background: rgba(0, 0, 0, 0.7);
}

.minimized-count {
  font-size: 14px;
  font-weight: 600;
  color: white;
}
</style>
