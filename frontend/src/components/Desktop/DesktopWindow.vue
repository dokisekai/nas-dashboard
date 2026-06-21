<template>
  <div
    class="desktop-window"
    :class="{
      'is-minimized': props.window.minimized,
      'is-maximized': props.window.maximized,
      'is-focused': props.window.focused
    }"
    :style="windowStyle"
    @mousedown.capture="handleFocus"
  >
    <div class="window-header" @mousedown="startDrag">
      <div class="window-title-group" @dblclick="handleMaximize">
        <div class="window-icon">
          <component :is="windowIcon" class="w-4 h-4" />
        </div>
        <div class="window-title">{{ props.window.title }}</div>
      </div>

      <div class="window-controls">
        <button @click.stop.prevent="handleMinimize" class="control-btn" title="最小化">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
          </svg>
        </button>
        <button @click.stop.prevent="handleMaximize" class="control-btn" title="最大化">
          <svg v-if="!props.window.maximized" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 0h-4" />
          </svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <button @click.stop.prevent="handleClose" class="control-btn close" title="关闭">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <div class="window-content">
      <component :is="appComponent" v-if="!props.window.minimized && appComponent" />
      <div v-else-if="!props.window.minimized && !appComponent" class="window-empty">
        该应用暂未提供界面
      </div>
    </div>

    <div
      v-if="!props.window.minimized && !props.window.maximized"
      class="resize-border"
      @mousedown.stop="startResize"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onUnmounted, defineAsyncComponent } from 'vue'
import { ServerIcon } from '@heroicons/vue/24/outline'
import { getApp } from '../../config/apps'
import type { DesktopWindow } from '../../stores/windowManager'

interface Props {
  window: DesktopWindow
}

interface Emits {
  (e: 'focus', windowId: string): void
  (e: 'close', windowId: string): void
  (e: 'minimize', windowId: string): void
  (e: 'maximize', windowId: string): void
  (e: 'drag-move', windowId: string, position: { x: number; y: number }): void
  (e: 'resize', windowId: string, size: { width: number; height: number }): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const SNAP_THRESHOLD = 20
const DOCK_HEIGHT = 80
const MIN_WIDTH = 400
const MIN_HEIGHT = 300

const isDragging = ref(false)
const isResizing = ref(false)
const dragOffset = ref({ x: 0, y: 0 })

const windowStyle = computed(() => ({
  left: `${props.window.position.x}px`,
  top: `${props.window.position.y}px`,
  width: props.window.maximized ? '100%' : `${props.window.size.width}px`,
  height: props.window.maximized ? 'calc(100% - 48px)' : `${props.window.size.height}px`,
  zIndex: props.window.zIndex || 100,
}))

const appDefinition = computed(() => getApp(props.window.appId))

const windowIcon = computed(() => appDefinition.value?.icon ?? ServerIcon)

const appComponent = computed(() => {
  const loader = appDefinition.value?.component
  if (!loader) return null
  return defineAsyncComponent(loader as any)
})

const handleFocus = () => emit('focus', props.window.id)
const handleClose = () => emit('close', props.window.id)
const handleMinimize = () => emit('minimize', props.window.id)
const handleMaximize = () => emit('maximize', props.window.id)

const computeSnappedPosition = (raw: { x: number; y: number }) => {
  const position = { ...raw }
  const { width, height } = props.window.size

  if (Math.abs(position.x) < SNAP_THRESHOLD) position.x = 0
  if (Math.abs(position.x + width - window.innerWidth) < SNAP_THRESHOLD) {
    position.x = window.innerWidth - width
  }
  if (Math.abs(position.y) < SNAP_THRESHOLD) position.y = 0
  if (Math.abs(position.y + height + DOCK_HEIGHT - window.innerHeight) < SNAP_THRESHOLD) {
    position.y = window.innerHeight - height - DOCK_HEIGHT
  }

  const centerX = (window.innerWidth - width) / 2
  const centerY = (window.innerHeight - height - DOCK_HEIGHT) / 2
  if (Math.abs(position.x - centerX) < SNAP_THRESHOLD && Math.abs(position.y - centerY) < SNAP_THRESHOLD) {
    position.x = centerX
    position.y = centerY
  }
  return position
}

const handleDragMove = (event: MouseEvent) => {
  if (!isDragging.value) return
  const raw = {
    x: event.clientX - dragOffset.value.x,
    y: event.clientY - dragOffset.value.y,
  }
  emit('drag-move', props.window.id, computeSnappedPosition(raw))
}

const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', handleDragMove)
  document.removeEventListener('mouseup', stopDrag)
}

const startDrag = (event: MouseEvent) => {
  if (props.window.maximized) return
  isDragging.value = true
  dragOffset.value = {
    x: event.clientX - props.window.position.x,
    y: event.clientY - props.window.position.y,
  }
  document.addEventListener('mousemove', handleDragMove)
  document.addEventListener('mouseup', stopDrag)
}

const startResize = (event: MouseEvent) => {
  isResizing.value = true
  const startX = event.clientX
  const startY = event.clientY
  const startWidth = props.window.size.width
  const startHeight = props.window.size.height

  const handleMouseMove = (e: MouseEvent) => {
    if (!isResizing.value) return
    emit('resize', props.window.id, {
      width: Math.max(MIN_WIDTH, startWidth + (e.clientX - startX)),
      height: Math.max(MIN_HEIGHT, startHeight + (e.clientY - startY)),
    })
  }

  const handleMouseUp = () => {
    isResizing.value = false
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
}

onUnmounted(() => {
  document.removeEventListener('mousemove', handleDragMove)
  document.removeEventListener('mouseup', stopDrag)
})
</script>

<style scoped>
.desktop-window {
  position: absolute;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  transition: box-shadow 0.2s ease, transform 0.2s ease, opacity 0.2s ease;
  animation: windowOpen 0.3s ease-out;
  pointer-events: auto !important;
  cursor: default;
}

@keyframes windowOpen {
  from { opacity: 0; transform: scale(0.9) translateY(-20px); }
  to   { opacity: 1; transform: scale(1) translateY(0); }
}

.desktop-window.is-focused {
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.4);
  border-color: rgba(255, 255, 255, 0.5);
}

.desktop-window.is-minimized {
  transform: scale(0.8);
  opacity: 0;
  pointer-events: none;
}

.desktop-window.is-maximized {
  border-radius: 0;
}

.window-header {
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 16px;
  cursor: move;
  user-select: none;
  flex-shrink: 0;
  pointer-events: auto !important;
  position: relative;
  z-index: 10;
}

.window-title-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  pointer-events: auto !important;
  cursor: pointer;
}

.window-icon {
  width: 20px;
  height: 20px;
  color: rgba(255, 255, 255, 0.9);
  pointer-events: none;
}

.window-title {
  color: white;
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 0.5px;
  pointer-events: none;
}

.window-controls {
  display: flex;
  gap: 8px;
  pointer-events: auto !important;
  position: relative;
  z-index: 20;
}

.control-btn {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: none;
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.9);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  pointer-events: auto !important;
  position: relative;
  z-index: 25;
}

.control-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.control-btn.close:hover {
  background: #ef4444;
  color: white;
}

.window-content {
  flex: 1;
  overflow: auto;
  background: white;
}

.window-empty {
  padding: 32px;
  text-align: center;
  color: #6b7280;
  font-size: 14px;
}

.resize-border {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 16px;
  height: 16px;
  cursor: se-resize;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.desktop-window:hover .resize-border {
  opacity: 1;
}

.resize-border::after {
  content: '';
  position: absolute;
  bottom: 4px;
  right: 4px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 0 4px 4px;
  border-color: transparent transparent #6b7280 transparent;
}

@media (prefers-color-scheme: dark) {
  .desktop-window {
    background: rgba(31, 41, 55, 0.95);
    border-color: rgba(255, 255, 255, 0.1);
  }
  .window-content {
    background: #1f2937;
  }
  .resize-border::after {
    border-color: transparent transparent #9ca3af transparent;
  }
}
</style>
