<template>
  <div class="desktop-settings-window">
    <div class="settings-header">
      <h2>桌面设置</h2>
      <button @click="$emit('close')" class="close-btn">
        <XMarkIcon class="w-5 h-5" />
      </button>
    </div>

    <div class="settings-body">
      <!-- 设置导航 -->
      <div class="settings-nav">
        <button
          v-for="section in sections"
          :key="section.id"
          class="nav-item"
          :class="{ active: activeSection === section.id }"
          @click="activeSection = section.id"
        >
          <component :is="section.icon" class="w-5 h-5" />
          <span>{{ section.name }}</span>
        </button>
      </div>

      <!-- 设置内容 -->
      <div class="settings-content">
        <component :is="getSectionComponent()" v-bind="sectionProps" @save="handleSave" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import ThemeManager from './ThemeManager.vue'
import DockSettings from './DockSettings.vue'

interface Props {
  currentSettings: any
  dockItems: any[]
  pinnedApps: any[]
}

interface Emits {
  (e: 'close'): void
  (e: 'save-settings', settings: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const activeSection = ref('theme')

const sections = [
  { id: 'theme', name: '主题', icon: 'SunIcon' },
  { id: 'dock', name: 'Dock', icon: 'ServerIcon' },
  { id: 'widgets', name: '小部件', icon: 'Square3x3Icon' },
  { id: 'windows', name: '窗口', icon: 'WindowIcon' },
  { id: 'keyboard', name: '键盘', icon: 'KeyboardIcon' }
]

const sectionProps = computed(() => {
  switch (activeSection.value) {
    case 'theme':
      return { currentSettings: props.currentSettings }
    case 'dock':
      return {
        currentSettings: props.currentSettings.dock || {},
        pinnedApps: props.pinnedApps
      }
    default:
      return {}
  }
})

const getSectionComponent = () => {
  switch (activeSection.value) {
    case 'theme':
      return ThemeManager
    case 'dock':
      return DockSettings
    default:
      return 'div'
  }
}

const handleSave = (data: any) => {
  emit('save-settings', data)
}
</script>

<style scoped>
.desktop-settings-window {
  width: 100%;
  height: 100%;
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

.settings-header h2 {
  font-size: 18px;
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
  color: #1f2937;
}

.settings-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.settings-nav {
  width: 200px;
  padding: 16px;
  border-right: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #6b7280;
  font-size: 14px;
  text-align: left;
}

.nav-item:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.nav-item.active {
  background: #3b82f6;
  color: white;
}

.settings-content {
  flex: 1;
  overflow-y: auto;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .desktop-settings-window {
    background: rgba(0, 0, 0, 0.85);
  }

  .settings-header h2 {
    color: #f9fafb;
  }

  .close-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #f9fafb;
  }

  .settings-nav {
    border-right-color: rgba(255, 255, 255, 0.1);
  }

  .nav-item {
    color: #9ca3af;
  }

  .nav-item:hover {
    background: rgba(255, 255, 255, 0.05);
    color: #f9fafb;
  }
}
</style>
