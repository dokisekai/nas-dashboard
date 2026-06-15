<template>
  <div class="enhanced-dock-container">
    <div class="dock-background">
      <div class="dock-reflection"></div>
    </div>

    <div class="dock" :class="{ 'is-expanded': isExpanded }">
      <!-- 应用分隔符 -->
      <div class="dock-separator" v-if="showSeparators"></div>

      <!-- 运行中的应用 -->
      <div
        v-for="item in runningApps"
        :key="item.id"
        class="dock-item"
        :class="{ active: isActiveDockItem(item), 'is-dragging': draggingItem === item.id }"
        :style="getItemStyle(item)"
        draggable="true"
        @click="handleDockItemClick(item)"
        @contextmenu.prevent="showContextMenu($event, item)"
        @dragstart="handleDragStart($event, item)"
        @dragend="handleDragEnd"
        @dragover.prevent="handleDragOver($event, item)"
        @drop="handleDrop($event, item)"
        @mouseenter="handleMouseEnter(item)"
        @mouseleave="handleMouseLeave"
      >
        <div class="dock-icon-wrapper">
          <div class="dock-icon">
            <component :is="item.icon" class="w-full h-full" />
            <!-- 运行指示器 -->
            <div v-if="isActiveDockItem(item)" class="running-indicator"></div>
          </div>
          <!-- 徽章 -->
          <div v-if="item.badge" class="dock-badge">{{ item.badge }}</div>
        </div>
        <div class="dock-tooltip">{{ item.label }}</div>
      </div>

      <!-- 固定应用分隔符 -->
      <div class="dock-separator" v-if="showSeparators && pinnedApps.length > 0"></div>

      <!-- 固定应用 -->
      <div
        v-for="item in pinnedApps"
        :key="item.id"
        class="dock-item pinned"
        :class="{ active: isActiveDockItem(item), 'is-dragging': draggingItem === item.id }"
        :style="getItemStyle(item)"
        draggable="true"
        @click="handleDockItemClick(item)"
        @contextmenu.prevent="showContextMenu($event, item)"
        @dragstart="handleDragStart($event, item)"
        @dragend="handleDragEnd"
        @dragover.prevent="handleDragOver($event, item)"
        @drop="handleDrop($event, item)"
        @mouseenter="handleMouseEnter(item)"
        @mouseleave="handleMouseLeave"
      >
        <div class="dock-icon-wrapper">
          <div class="dock-icon">
            <component :is="item.icon" class="w-full h-full" />
            <div v-if="isActiveDockItem(item)" class="running-indicator"></div>
          </div>
          <div v-if="item.badge" class="dock-badge">{{ item.badge }}</div>
          <!-- 固定指示器 -->
          <div class="pin-indicator"></div>
        </div>
        <div class="dock-tooltip">{{ item.label }}</div>
      </div>

      <!-- 最小化窗口分隔符 -->
      <div class="dock-separator" v-if="showSeparators && minimizedApps.length > 0"></div>

      <!-- 最小化应用 -->
      <div
        v-for="item in minimizedApps"
        :key="item.id"
        class="dock-item minimized"
        :class="{ active: false }"
        :style="getItemStyle(item)"
        @click="handleDockItemClick(item)"
        @contextmenu.prevent="showContextMenu($event, item)"
        @mouseenter="handleMouseEnter(item)"
        @mouseleave="handleMouseLeave"
      >
        <div class="dock-icon-wrapper">
          <div class="dock-icon minimized-icon">
            <component :is="item.icon" class="w-full h-full" />
          </div>
        </div>
        <div class="dock-tooltip">{{ item.label }}</div>
      </div>
    </div>

    <!-- Dock设置按钮 -->
    <div class="dock-settings">
      <button @click="openDockSettings" class="settings-btn">
        <CogIcon class="w-5 h-5" />
      </button>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenuVisible"
      class="dock-context-menu"
      :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
    >
      <div class="context-menu-items">
        <div
          v-for="menuItem in contextMenuItems"
          :key="menuItem.id"
          class="context-menu-item"
          :class="menuItem.class"
          @click="handleContextMenuClick(menuItem)"
        >
          <component v-if="menuItem.icon" :is="menuItem.icon" class="w-4 h-4" />
          <span>{{ menuItem.label }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { CogIcon, XMarkIcon, MapPinIcon as PinIcon, PencilIcon } from '@heroicons/vue/24/outline'
import type { DockItem } from '../../types/desktop'

interface Props {
  dockItems: DockItem[]
  runningApps: DockItem[]
  pinnedApps: DockItem[]
  minimizedApps: DockItem[]
  showSeparators?: boolean
}

interface Emits {
  (e: 'click-item', item: DockItem): void
  (e: 'pin-item', item: DockItem): void
  (e: 'unpin-item', item: DockItem): void
  (e: 'quit-app', item: DockItem): void
  (e: 'reorder-items', items: DockItem[]): void
  (e: 'open-settings'): void
}

const props = withDefaults(defineProps<Props>(), {
  showSeparators: true
})

const emit = defineEmits<Emits>()

const isExpanded = ref(false)
const draggingItem = ref<string | null>(null)
const hoveredItem = ref<DockItem | null>(null)

// 右键菜单
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuTarget = ref<DockItem | null>(null)

const contextMenuItems = computed(() => {
  if (!contextMenuTarget.value) return []

  const items = [
    {
      id: 'open',
      label: '打开',
      icon: null,
      action: () => handleDockItemClick(contextMenuTarget.value!)
    }
  ]

  if (contextMenuTarget.value.appId) {
    items.push({
      id: 'pin',
      label: props.pinnedApps.find(p => p.id === contextMenuTarget.value!.id)
        ? '从 Dock 移除'
        : '保留在 Dock 中',
      icon: PinIcon,
      action: () => handleTogglePin(contextMenuTarget.value!)
    })
  }

  items.push({
    id: 'quit',
    label: '退出',
    icon: XMarkIcon,
    class: 'danger',
    action: () => handleQuitApp(contextMenuTarget.value!)
  })

  return items
})

const getItemStyle = (item: DockItem) => {
  if (hoveredItem.value?.id === item.id) {
    return {
      transform: 'scale(1.2) translateY(-8px)'
    }
  }
  return {}
}

const isActiveDockItem = (item: DockItem) => {
  return props.runningApps.some(app => app.id === item.id)
}

const handleDockItemClick = (item: DockItem) => {
  emit('click-item', item)
}

const handleMouseEnter = (item: DockItem) => {
  hoveredItem.value = item
  isExpanded.value = true
}

const handleMouseLeave = () => {
  hoveredItem.value = null
  setTimeout(() => {
    if (!hoveredItem.value) {
      isExpanded.value = false
    }
  }, 100)
}

// 拖拽功能
const handleDragStart = (event: DragEvent, item: DockItem) => {
  draggingItem.value = item.id
  event.dataTransfer?.setData('text/plain', item.id)
}

const handleDragEnd = () => {
  draggingItem.value = null
}

const handleDragOver = (event: DragEvent, item: DockItem) => {
  event.preventDefault()
  event.dataTransfer!.dropEffect = 'move'
}

const handleDrop = (event: DragEvent, targetItem: DockItem) => {
  event.preventDefault()
  const draggedId = event.dataTransfer?.getData('text/plain')

  if (draggedId && draggedId !== targetItem.id) {
    // 重新排序逻辑
    const items = [...props.pinnedApps, ...props.runningApps]
    const draggedIndex = items.findIndex(item => item.id === draggedId)
    const targetIndex = items.findIndex(item => item.id === targetItem.id)

    if (draggedIndex !== -1 && targetIndex !== -1) {
      const [draggedItem] = items.splice(draggedIndex, 1)
      items.splice(targetIndex, 0, draggedItem)
      emit('reorder-items', items)
    }
  }
}

// 右键菜单
const showContextMenu = (event: MouseEvent, item: DockItem) => {
  contextMenuTarget.value = item
  contextMenuPosition.value = { x: event.clientX, y: event.clientY }
  contextMenuVisible.value = true
}

const handleContextMenuClick = (menuItem: any) => {
  menuItem.action?.()
  contextMenuVisible.value = false
}

const handleTogglePin = (item: DockItem) => {
  if (props.pinnedApps.find(p => p.id === item.id)) {
    emit('unpin-item', item)
  } else {
    emit('pin-item', item)
  }
}

const handleQuitApp = (item: DockItem) => {
  emit('quit-app', item)
}

const openDockSettings = () => {
  emit('open-settings')
}

// 点击外部关闭菜单
const handleClickOutside = () => {
  contextMenuVisible.value = false
}

// 监听全局点击事件
if (typeof window !== 'undefined') {
  document.addEventListener('click', handleClickOutside)
}
</script>

<style scoped>
.enhanced-dock-container {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 1000;
}

.dock-background {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 100%;
  height: 100%;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.3);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.dock-reflection {
  position: absolute;
  bottom: -20px;
  left: 0;
  right: 0;
  height: 20px;
  background: linear-gradient(to bottom, rgba(255, 255, 255, 0.3), transparent);
  filter: blur(10px);
  transform: scaleY(-1);
  opacity: 0.3;
}

.dock {
  position: relative;
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 12px 16px;
  transition: all 0.3s ease;
}

.dock-separator {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.3);
  margin: 0 8px;
}

.dock-item {
  position: relative;
  width: 56px;
  height: 56px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.dock-item:hover {
  transform: translateY(-8px) scale(1.1);
}

.dock-item.is-dragging {
  opacity: 0.5;
  cursor: grabbing;
}

.dock-icon-wrapper {
  position: relative;
  width: 100%;
  height: 100%;
}

.dock-icon {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  position: relative;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
}

.dock-item.minimized .dock-icon {
  opacity: 0.6;
  transform: scale(0.8);
}

.running-indicator {
  position: absolute;
  bottom: -6px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.8);
}

.pin-indicator {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 12px;
  height: 12px;
  background: #3b82f6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 8px;
  border: 2px solid white;
}

.dock-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  min-width: 18px;
  height: 18px;
  padding: 0 4px;
  background: #ef4444;
  color: white;
  font-size: 11px;
  font-weight: 600;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.dock-tooltip {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%) translateY(8px);
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.8);
  color: white;
  font-size: 12px;
  border-radius: 6px;
  white-space: nowrap;
  opacity: 0;
  pointer-events: none;
  transition: all 0.2s ease;
}

.dock-item:hover .dock-tooltip {
  opacity: 1;
  transform: translateX(-50%) translateY(0);
}

.dock-settings {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
}

.settings-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
}

.settings-btn:hover {
  background: rgba(255, 255, 255, 0.3);
  transform: scale(1.1);
}

.dock-context-menu {
  position: fixed;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  padding: 8px 0;
  min-width: 180px;
  z-index: 10000;
  animation: menuFadeIn 0.2s ease;
}

@keyframes menuFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.context-menu-items {
  display: flex;
  flex-direction: column;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 13px;
  color: #1f2937;
}

.context-menu-item:hover {
  background: rgba(0, 0, 0, 0.05);
}

.context-menu-item.danger {
  color: #ef4444;
}

.context-menu-item.danger:hover {
  background: rgba(239, 68, 68, 0.1);
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .dock-background {
    background: rgba(0, 0, 0, 0.3);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .dock-separator {
    background: rgba(255, 255, 255, 0.1);
  }

  .dock-context-menu {
    background: rgba(0, 0, 0, 0.85);
  }

  .context-menu-item {
    color: #f9fafb;
  }

  .context-menu-item:hover {
    background: rgba(255, 255, 255, 0.1);
  }
}
</style>
