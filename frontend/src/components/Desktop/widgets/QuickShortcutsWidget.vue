<template>
  <div class="shortcuts-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <LinkIcon class="w-5 h-5" />
      <span class="widget-title">快捷方式</span>
    </div>

    <div class="widget-content">
      <!-- 快捷方式网格 -->
      <div class="shortcuts-grid">
        <div
          v-for="shortcut in shortcuts"
          :key="shortcut.id"
          class="shortcut-item"
          @click="openShortcut(shortcut)"
          @contextmenu.prevent="showContextMenu($event, shortcut)"
        >
          <div class="shortcut-icon" :style="{ background: shortcut.color }">
            <component :is="getShortcutIcon(shortcut.icon)" class="w-6 h-6" />
          </div>
          <span class="shortcut-label">{{ shortcut.label }}</span>
        </div>
      </div>

      <!-- 添加快捷方式按钮 -->
      <button v-if="shortcuts.length < maxShortcuts" class="add-shortcut-btn" @click="showAddDialog = true">
        <PlusIcon class="w-5 h-5" />
        <span>添加</span>
      </button>
    </div>

    <!-- 添加快捷方式对话框 -->
    <div v-if="showAddDialog" class="dialog-overlay" @click="showAddDialog = false">
      <div class="dialog" @click.stop>
        <div class="dialog-header">
          <h3>添加快捷方式</h3>
          <button class="close-btn" @click="showAddDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="dialog-body">
          <div class="form-group">
            <label>名称</label>
            <input v-model="newShortcut.label" type="text" placeholder="输入名称" />
          </div>
          <div class="form-group">
            <label>类型</label>
            <select v-model="newShortcut.type">
              <option value="app">应用</option>
              <option value="url">网址</option>
              <option value="file">文件</option>
              <option value="folder">文件夹</option>
            </select>
          </div>
          <div class="form-group">
            <label>路径/URL</label>
            <input v-model="newShortcut.path" type="text" placeholder="输入路径或URL" />
          </div>
          <div class="form-group">
            <label>图标颜色</label>
            <div class="color-options">
              <div
                v-for="color in colorOptions"
                :key="color"
                class="color-option"
                :class="{ selected: newShortcut.color === color }"
                :style="{ background: color }"
                @click="newShortcut.color = color"
              />
            </div>
          </div>
        </div>
        <div class="dialog-footer">
          <button class="btn-secondary" @click="showAddDialog = false">取消</button>
          <button class="btn-primary" @click="addShortcut">添加</button>
        </div>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenu.show"
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <button class="context-menu-item" @click="editShortcut(contextMenu.shortcut)">
        <PencilIcon class="w-4 h-4" />
        编辑
      </button>
      <button class="context-menu-item danger" @click="removeShortcut(contextMenu.shortcut)">
        <TrashIcon class="w-4 h-4" />
        删除
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  LinkIcon,
  PlusIcon,
  XMarkIcon,
  PencilIcon,
  TrashIcon,
  FolderIcon,
  CogIcon,
  ChartBarIcon,
  UserGroupIcon,
  CloudIcon,
  GlobeAltIcon
} from '@heroicons/vue/24/outline'

interface Props {
  config: {
    maxShortcuts?: number
    editable?: boolean
  }
  size: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    maxShortcuts: 8,
    editable: true
  })
})

const shortcuts = ref([
  {
    id: '1',
    label: '存储管理',
    type: 'app',
    path: '/storage',
    icon: 'FolderIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)'
  },
  {
    id: '2',
    label: '系统监控',
    type: 'app',
    path: '/monitor',
    icon: 'ChartBarIcon',
    color: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)'
  },
  {
    id: '3',
    label: '用户管理',
    type: 'app',
    path: '/users',
    icon: 'UserGroupIcon',
    color: 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)'
  },
  {
    id: '4',
    label: '系统设置',
    type: 'app',
    path: '/settings',
    icon: 'CogIcon',
    color: 'linear-gradient(135deg, #6b7280 0%, #9ca3af 100%)'
  }
])

const showAddDialog = ref(false)
const contextMenu = ref({
  show: false,
  x: 0,
  y: 0,
  shortcut: null as any
})

const newShortcut = ref({
  label: '',
  type: 'app',
  path: '',
  icon: 'FolderIcon',
  color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)'
})

const colorOptions = [
  'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
  'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
  'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
  'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
  'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
  'linear-gradient(135deg, #06b6d4 0%, #38bdf8 100%)'
]

const maxShortcuts = computed(() => props.config.maxShortcuts)

const getShortcutIcon = (iconName: string) => {
  const icons: Record<string, any> = {
    FolderIcon,
    CogIcon,
    ChartBarIcon,
    UserGroupIcon,
    CloudIcon,
    GlobeAltIcon
  }
  return icons[iconName] || FolderIcon
}

const openShortcut = (shortcut: any) => {
  console.log('Opening shortcut:', shortcut)
  // 触发打开快捷方式的逻辑
}

const showContextMenu = (event: MouseEvent, shortcut: any) => {
  if (!props.config.editable) return

  contextMenu.value = {
    show: true,
    x: event.clientX,
    y: event.clientY,
    shortcut
  }
}

const editShortcut = (shortcut: any) => {
  contextMenu.value.show = false
  // 实现编辑逻辑
}

const removeShortcut = (shortcut: any) => {
  const index = shortcuts.value.findIndex(s => s.id === shortcut.id)
  if (index > -1) {
    shortcuts.value.splice(index, 1)
  }
  contextMenu.value.show = false
}

const addShortcut = () => {
  if (!newShortcut.value.label || !newShortcut.value.path) {
    return
  }

  shortcuts.value.push({
    id: Date.now().toString(),
    ...newShortcut.value
  })

  // 重置表单
  newShortcut.value = {
    label: '',
    type: 'app',
    path: '',
    icon: 'FolderIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)'
  }

  showAddDialog.value = false
}

// 点击其他地方关闭右键菜单
document.addEventListener('click', () => {
  contextMenu.value.show = false
})
</script>

<style scoped>
.shortcuts-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.widget-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  margin-bottom: 12px;
}

.widget-header svg {
  color: #6b7280;
}

.widget-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.widget-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.shortcuts-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.widget-medium .shortcuts-grid {
  grid-template-columns: repeat(3, 1fr);
}

.widget-large .shortcuts-grid {
  grid-template-columns: repeat(4, 1fr);
}

.shortcut-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 8px;
  background: rgba(249, 250, 251, 0.5);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.shortcut-item:hover {
  background: rgba(243, 244, 246, 0.8);
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
}

.shortcut-label {
  font-size: 11px;
  color: #1f2937;
  text-align: center;
  word-break: break-word;
}

.add-shortcut-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 12px;
  background: rgba(249, 250, 251, 0.5);
  border: 2px dashed rgba(209, 213, 219, 0.5);
  border-radius: 8px;
  color: #6b7280;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-shortcut-btn:hover {
  background: rgba(243, 244, 246, 0.8);
  border-color: #9ca3af;
}

/* 对话框样式 */
.dialog-overlay {
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

.dialog {
  background: white;
  border-radius: 16px;
  padding: 24px;
  width: 90%;
  max-width: 400px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1);
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.dialog-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 4px;
}

.dialog-body {
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
  color: #1f2937;
}

.form-group input,
.form-group select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  color: #1f2937;
}

.color-options {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.color-option {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.2s ease;
}

.color-option:hover {
  transform: scale(1.1);
}

.color-option.selected {
  border-color: #1f2937;
}

.dialog-footer {
  display: flex;
  gap: 8px;
  margin-top: 20px;
  justify-content: flex-end;
}

.btn-secondary,
.btn-primary {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
}

.btn-secondary {
  background: #e5e7eb;
  color: #1f2937;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

/* 右键菜单 */
.context-menu {
  position: fixed;
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
  padding: 4px;
  z-index: 1000;
  min-width: 120px;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border: none;
  background: none;
  border-radius: 4px;
  font-size: 14px;
  color: #1f2937;
  cursor: pointer;
  width: 100%;
  text-align: left;
}

.context-menu-item:hover {
  background: #f3f4f6;
}

.context-menu-item.danger {
  color: #ef4444;
}

.context-menu-item.danger:hover {
  background: #fef2f2;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .shortcuts-widget {
    background: rgba(0, 0, 0, 0.8);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .widget-title,
  .shortcut-label,
  .form-group label {
    color: #f9fafb;
  }

  .shortcut-item,
  .add-shortcut-btn {
    background: rgba(255, 255, 255, 0.05);
  }

  .shortcut-item:hover,
  .add-shortcut-btn:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .dialog {
    background: #1f2937;
  }

  .dialog-header h3,
  .form-group label,
  .form-group input,
  .form-group select {
    color: #f9fafb;
  }

  .form-group input,
  .form-group select {
    border-color: rgba(255, 255, 255, 0.1);
    background: rgba(255, 255, 255, 0.05);
  }
}
</style>
