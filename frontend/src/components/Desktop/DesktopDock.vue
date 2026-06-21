<template>
  <div class="dock">
    <div
      v-for="app in apps"
      :key="app.id"
      class="dock-item"
      :class="{ active: activeAppIds.includes(app.id) }"
      :title="app.title"
      @click="$emit('open-app', app.id)"
    >
      <div class="dock-icon">
        <component :is="app.icon" class="w-6 h-6" />
      </div>
    </div>

    <div class="dock-separator"></div>

    <div class="dock-item window-manager" title="窗口管理" @click="$emit('toggle-window-menu')">
      <div class="dock-icon">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
        </svg>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AppDefinition } from '../../config/apps'

defineProps<{
  apps: AppDefinition[]
  activeAppIds: string[]
}>()

defineEmits<{
  (e: 'open-app', appId: string): void
  (e: 'toggle-window-menu'): void
}>()
</script>

<style scoped>
.dock {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  z-index: 100;
  pointer-events: all;
}

.dock-item {
  width: 56px;
  height: 56px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.dock-item:hover {
  transform: translateY(-8px) scale(1.1);
}

.dock-item.active::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.8);
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
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.dock-separator {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.3);
  margin: 0 8px;
}

.window-manager {
  cursor: pointer;
}

@media (prefers-color-scheme: dark) {
  .dock {
    background: rgba(0, 0, 0, 0.3);
    border-color: rgba(255, 255, 255, 0.1);
  }
}
</style>
