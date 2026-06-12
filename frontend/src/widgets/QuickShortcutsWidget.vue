<template>
  <div class="shortcuts-widget">
    <div class="widget-header">
      <h3>快捷方式</h3>
      <button class="add-btn" @click="addShortcut" title="添加快捷方式">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
      </button>
    </div>

    <div class="shortcuts-grid">
      <div
        v-for="shortcut in shortcuts"
        :key="shortcut.id"
        class="shortcut-item"
        @click="openShortcut(shortcut)"
        @contextmenu.prevent="showContextMenu($event, shortcut)"
      >
        <div class="shortcut-icon" :style="{ background: shortcut.color }">
          <component :is="getIcon(shortcut.icon)" class="w-6 h-6" />
        </div>
        <span class="shortcut-label">{{ shortcut.label }}</span>
      </div>
    </div>

    <!-- 上下文菜单 -->
    <div
      v-if="contextMenu.visible"
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
      @click="hideContextMenu"
    >
      <div class="context-menu-item" @click="editShortcut(contextMenu.shortcut)">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
        编辑
      </div>
      <div class="context-menu-item" @click="removeShortcut(contextMenu.shortcut)">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
        删除
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  ServerIcon,
  ChartBarIcon,
  CogIcon,
  UserGroupIcon,
  FolderIcon,
  ShoppingBagIcon,
  CloudIcon,
  ShieldCheckIcon
} from '@heroicons/vue/24/outline'

interface Shortcut {
  id: string
  label: string
  icon: string
  color: string
  action: string
  appId?: string
}

const shortcuts = ref<Shortcut[]>([
  {
    id: '1',
    label: '存储管理',
    icon: 'ServerIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
    action: 'open-app',
    appId: 'storage-manager'
  },
  {
    id: '2',
    label: '系统监控',
    icon: 'ChartBarIcon',
    color: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    action: 'open-app',
    appId: 'system-monitor'
  },
  {
    id: '3',
    label: '用户管理',
    icon: 'UserGroupIcon',
    color: 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
    action: 'open-app',
    appId: 'user-manager'
  },
  {
    id: '4',
    label: '文件管理',
    icon: 'FolderIcon',
    color: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    action: 'open-app',
    appId: 'file-manager'
  },
  {
    id: '5',
    label: '应用中心',
    icon: 'ShoppingBagIcon',
    color: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
    action: 'open-app',
    appId: 'app-center'
  },
  {
    id: '6',
    label: '系统设置',
    icon: 'CogIcon',
    color: 'linear-gradient(135deg, #6b7280 0%, #9ca3af 100%)',
    action: 'open-app',
    appId: 'settings'
  }
])

const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  shortcut: null as Shortcut | null
})

const getIcon = (iconName: string) => {
  const icons: Record<string, any> = {
    ServerIcon,
    ChartBarIcon,
    CogIcon,
    UserGroupIcon,
    FolderIcon,
    ShoppingBagIcon,
    CloudIcon,
    ShieldCheckIcon
  }
  return icons[iconName] || FolderIcon
}

const openShortcut = (shortcut: Shortcut) => {
  if (shortcut.action === 'open-app' && shortcut.appId) {
    // 触发打开应用事件
    const event = new CustomEvent('open-app', { detail: shortcut.appId })
    window.dispatchEvent(event)
  }
}

const addShortcut = () => {
  // 显示添加快捷方式对话框
  const label = prompt('快捷方式名称:')
  if (!label) return

  shortcuts.value.push({
    id: Date.now().toString(),
    label,
    icon: 'FolderIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
    action: 'open-app'
  })
}

const editShortcut = (shortcut: Shortcut) => {
  const newLabel = prompt('编辑快捷方式名称:', shortcut.label)
  if (newLabel) {
    shortcut.label = newLabel
  }
}

const removeShortcut = (shortcut: Shortcut) => {
  if (confirm(`确定要删除"${shortcut.label}"快捷方式吗?`)) {
    shortcuts.value = shortcuts.value.filter(s => s.id !== shortcut.id)
  }
}

const showContextMenu = (event: MouseEvent, shortcut: Shortcut) => {
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    shortcut
  }
}

const hideContextMenu = () => {
  contextMenu.value.visible = false
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  if (!target.closest('.context-menu') && !target.closest('.shortcut-item')) {
    hideContextMenu()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.shortcuts-widget {
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.shortcuts-widget:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.widget-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.widget-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.add-btn {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 6px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.3);
}

.shortcuts-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.shortcut-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 12px 8px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.shortcut-item:hover {
  background: rgba(102, 126, 234, 0.1);
  transform: translateY(-2px);
}

.shortcut-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.shortcut-label {
  font-size: 11px;
  color: #374151;
  font-weight: 500;
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.context-menu {
  position: fixed;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 8px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.5);
  z-index: 1000;
  min-width: 150px;
  animation: contextMenuSlideIn 0.15s ease-out;
}

@keyframes contextMenuSlideIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  font-size: 13px;
  color: #374151;
  cursor: pointer;
  transition: all 0.15s ease;
}

.context-menu-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.context-menu-item:first-child {
  border-radius: 8px 8px 0 0;
}

.context-menu-item:last-child {
  border-radius: 0 0 8px 8px;
}
</style>
