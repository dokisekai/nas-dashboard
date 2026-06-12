<template>
  <div class="context-menu" :style="positionStyle" @click.stop>
    <div
      v-for="(item, index) in items"
      :key="index"
      class="context-menu-item"
      :class="{ divider: item.divider, danger: item.danger }"
      @click="handleItemClick(item)"
    >
      <div v-if="item.icon" class="item-icon">
        <component :is="item.icon" class="w-4 h-4" />
      </div>
      <span class="item-label">{{ item.label }}</span>
      <div v-if="item.shortcut" class="item-shortcut">{{ item.shortcut }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface ContextMenuItem {
  label: string
  icon?: any
  shortcut?: string
  danger?: boolean
  divider?: boolean
  action?: () => void
}

interface Props {
  position: { x: number; y: number }
  items: ContextMenuItem[]
}

interface Emits {
  (e: 'close'): void
  (e: 'select', item: ContextMenuItem): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const positionStyle = computed(() => ({
  left: `${props.position.x}px`,
  top: `${props.position.y}px`
}))

const handleItemClick = (item: ContextMenuItem) => {
  if (item.divider) return
  emit('select', item)
  emit('close')
}
</script>

<style scoped>
.context-menu {
  position: fixed;
  min-width: 180px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
  padding: 4px;
  z-index: 10000;
  animation: contextMenuFadeIn 0.15s ease;
}

@keyframes contextMenuFadeIn {
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
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.1s ease;
  font-size: 13px;
  color: #1f2937;
  position: relative;
}

.context-menu-item:hover {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.context-menu-item.danger:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.context-menu-item.divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.1);
  margin: 4px 0;
  padding: 0;
  pointer-events: none;
}

.item-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  color: currentColor;
}

.item-label {
  flex: 1;
  font-weight: 500;
}

.item-shortcut {
  font-size: 11px;
  color: #9ca3af;
  opacity: 0.7;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .context-menu {
    background: rgba(31, 41, 55, 0.95);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .context-menu-item {
    color: #f9fafb;
  }

  .context-menu-item.divider {
    background: rgba(255, 255, 255, 0.1);
  }

  .item-shortcut {
    color: #9ca3af;
  }
}
</style>
