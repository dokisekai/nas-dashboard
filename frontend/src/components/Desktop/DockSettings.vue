<template>
  <div class="dock-settings">
    <div class="settings-header">
      <h3>Dock 设置</h3>
      <button @click="$emit('close')" class="close-btn">
        <XMarkIcon class="w-5 h-5" />
      </button>
    </div>

    <div class="settings-content">
      <!-- 外观设置 -->
      <div class="settings-section">
        <h4>外观</h4>
        <div class="setting-item">
          <label>Dock 大小</label>
          <div class="size-options">
            <button
              v-for="size in dockSizes"
              :key="size.value"
              class="size-btn"
              :class="{ active: settings.size === size.value }"
              @click="settings.size = size.value"
            >
              {{ size.label }}
            </button>
          </div>
        </div>

        <div class="setting-item">
          <label class="checkbox-label">
            <input
              v-model="settings.showSeparators"
              type="checkbox"
              class="checkbox"
            />
            <span>显示分隔符</span>
          </label>
        </div>

        <div class="setting-item">
          <label class="checkbox-label">
            <input
              v-model="settings.hideWhenInactive"
              type="checkbox"
              class="checkbox"
            />
            <span>不活动时自动隐藏</span>
          </label>
        </div>

        <div class="setting-item">
          <label class="checkbox-label">
            <input
              v-model="settings.magnification"
              type="checkbox"
              class="checkbox"
            />
            <span>放大效果</span>
          </label>
        </div>
      </div>

      <!-- 位置设置 -->
      <div class="settings-section">
        <h4>位置</h4>
        <div class="position-options">
          <button
            v-for="position in positions"
            :key="position.value"
            class="position-btn"
            :class="{ active: settings.position === position.value }"
            @click="settings.position = position.value"
          >
            <component :is="position.icon" class="w-6 h-6" />
            <span>{{ position.label }}</span>
          </button>
        </div>
      </div>

      <!-- 行为设置 -->
      <div class="settings-section">
        <h4>行为</h4>
        <div class="setting-item">
          <label>点击动作</label>
          <select v-model="settings.clickAction" class="form-select">
            <option value="open">打开应用</option>
            <option value="switch">切换窗口</option>
            <option value="minimize">最小化/还原</option>
          </select>
        </div>

        <div class="setting-item">
          <label>双击动作</label>
          <select v-model="settings.doubleClickAction" class="form-select">
            <option value="new-window">新建窗口</option>
            <option value="maximize">最大化</option>
            <option value="fullscreen">全屏</option>
          </select>
        </div>
      </div>

      <!-- Dock 内容 -->
      <div class="settings-section">
        <h4>Dock 内容</h4>
        <div class="docked-apps">
          <draggable
            v-model="localPinnedApps"
            item-key="id"
            @change="handleReorder"
          >
            <template #item="{ element: app }">
              <div class="docked-app">
                <div class="app-info">
                  <div class="app-icon">
                    <component :is="app.icon" class="w-6 h-6" />
                  </div>
                  <span>{{ app.label }}</span>
                </div>
                <button @click="removeApp(app)" class="remove-btn">
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </template>
          </draggable>
        </div>
      </div>
    </div>

    <div class="settings-footer">
      <button @click="resetDefaults" class="btn btn-secondary">重置默认</button>
      <button @click="saveSettings" class="btn btn-primary">保存</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import { VueDraggable as draggable } from 'vue-draggable-plus'
import type { DockItem } from '../../types/desktop'

interface Props {
  currentSettings: any
  pinnedApps: DockItem[]
}

interface Emits {
  (e: 'close'): void
  (e: 'save', settings: any): void
  (e: 'remove-app', app: DockItem): void
  (e: 'reorder-apps', apps: DockItem[]): void
  (e: 'update:pinnedApps', apps: DockItem[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const settings = reactive({
  size: props.currentSettings.size || 'medium',
  showSeparators: props.currentSettings.showSeparators ?? true,
  hideWhenInactive: props.currentSettings.hideWhenInactive ?? false,
  magnification: props.currentSettings.magnification ?? true,
  position: props.currentSettings.position || 'bottom',
  clickAction: props.currentSettings.clickAction || 'open',
  doubleClickAction: props.currentSettings.doubleClickAction || 'new-window'
})

// 创建本地可变的状态
const localPinnedApps = computed({
  get: () => props.pinnedApps,
  set: (value) => emit('update:pinnedApps', value)
})

const dockSizes = [
  { value: 'small', label: '小' },
  { value: 'medium', label: '中' },
  { value: 'large', label: '大' }
]

const positions = [
  { value: 'bottom', label: '底部', icon: 'ArrowDownIcon' },
  { value: 'left', label: '左侧', icon: 'ArrowLeftIcon' },
  { value: 'right', label: '右侧', icon: 'ArrowRightIcon' }
]

const saveSettings = () => {
  emit('save', settings)
}

const resetDefaults = () => {
  settings.size = 'medium'
  settings.showSeparators = true
  settings.hideWhenInactive = false
  settings.magnification = true
  settings.position = 'bottom'
  settings.clickAction = 'open'
  settings.doubleClickAction = 'new-window'
}

const removeApp = (app: DockItem) => {
  emit('remove-app', app)
}

const handleReorder = () => {
  emit('reorder-apps', localPinnedApps.value)
  emit('update:pinnedApps', localPinnedApps.value)
}
</script>

<style scoped>
.dock-settings {
  width: 100%;
  max-width: 520px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.settings-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #6b7280;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.settings-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.settings-section {
  margin-bottom: 24px;
}

.settings-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.setting-item {
  margin-bottom: 16px;
}

.setting-item label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox {
  width: 18px;
  height: 18px;
  accent-color: #3b82f6;
}

.checkbox-label span {
  font-size: 13px;
  color: #374151;
}

.size-options {
  display: flex;
  gap: 8px;
}

.size-btn {
  flex: 1;
  padding: 10px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  font-size: 13px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.size-btn:hover {
  background: rgba(0, 0, 0, 0.05);
}

.size-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.position-options {
  display: flex;
  gap: 12px;
}

.position-btn {
  flex: 1;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.position-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  transform: translateY(-2px);
}

.position-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.form-select {
  width: 100%;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  font-size: 13px;
  color: #1f2937;
  transition: all 0.2s ease;
}

.form-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.docked-apps {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.docked-app {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.05);
  border-radius: 8px;
  cursor: move;
  transition: all 0.2s ease;
}

.docked-app:hover {
  background: rgba(0, 0, 0, 0.05);
  transform: translateX(4px);
}

.app-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-icon {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.remove-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #ef4444;
  transition: all 0.2s ease;
}

.remove-btn:hover {
  background: rgba(239, 68, 68, 0.1);
}

.settings-footer {
  display: flex;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.btn {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary {
  background: rgba(0, 0, 0, 0.05);
  color: #6b7280;
}

.btn-secondary:hover {
  background: rgba(0, 0, 0, 0.1);
  color: #1f2937;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .dock-settings {
    background: rgba(0, 0, 0, 0.85);
  }

  .settings-header h3,
  .settings-section h4,
  .setting-item label,
  .checkbox-label span {
    color: #f9fafb;
  }

  .size-btn,
  .position-btn,
  .docked-app {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .size-btn:hover,
  .position-btn:hover,
  .docked-app:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  .form-select {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
    color: #f9fafb;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.05);
    color: #9ca3af;
  }
}
</style>
