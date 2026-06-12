<template>
  <div
    class="desktop-widget"
    :class="[`widget-${widget.size}`, { 'is-dragging': isDragging }]"
    :style="widgetPosition"
    draggable="true"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
    @click="handleClick"
  >
    <!-- 系统监控小部件 -->
    <SystemMonitorWidget
      v-if="widget.type === 'system-monitor'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 存储状态小部件 -->
    <StorageStatusWidget
      v-else-if="widget.type === 'storage-status'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 网络监控小部件 -->
    <NetworkMonitorWidget
      v-else-if="widget.type === 'network-monitor'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 时钟小部件 -->
    <ClockWidget
      v-else-if="widget.type === 'clock'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 天气小部件 -->
    <WeatherWidget
      v-else-if="widget.type === 'weather'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 日历小部件 -->
    <CalendarWidget
      v-else-if="widget.type === 'calendar'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 快速笔记小部件 -->
    <QuickNoteWidget
      v-else-if="widget.type === 'quick-note'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 快捷方式小部件 -->
    <QuickShortcutsWidget
      v-else-if="widget.type === 'quick-shortcuts'"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 通用小部件容器（用于插件） -->
    <component
      v-else
      :is="getPluginWidgetComponent(widget.type)"
      :config="widget.config"
      :size="widget.size"
    />

    <!-- 小部件控制按钮 -->
    <div class="widget-controls" v-if="showControls">
      <button @click.stop="handleEdit" class="control-btn">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
        </svg>
      </button>
      <button @click.stop="handleRemove" class="control-btn danger">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- 调整大小手柄 -->
    <div
      v-if="widget.size !== 'small'"
      class="resize-handle"
      @mousedown.stop="handleResizeStart"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Widget } from '../../types/desktop'
import SystemMonitorWidget from './widgets/SystemMonitorWidget.vue'
import StorageStatusWidget from './widgets/StorageStatusWidget.vue'
import NetworkMonitorWidget from './widgets/NetworkMonitorWidget.vue'
import ClockWidget from './widgets/ClockWidget.vue'
import WeatherWidget from './widgets/WeatherWidget.vue'
import CalendarWidget from './widgets/CalendarWidget.vue'
import QuickNoteWidget from './widgets/QuickNoteWidget.vue'
import QuickShortcutsWidget from './widgets/QuickShortcutsWidget.vue'

interface Props {
  widget: Widget
}

interface Emits {
  (e: 'drag-start', widget: Widget): void
  (e: 'drag-end', widget: Widget, position: { x: number; y: number }): void
  (e: 'remove', widgetId: string): void
  (e: 'open-app', appId: string): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const isDragging = ref(false)
const showControls = ref(false)
const dragOffset = ref({ x: 0, y: 0 })

const widgetPosition = computed(() => ({
  left: `${props.widget.position.x}px`,
  top: `${props.widget.position.y}px`,
  width: getWidgetSize(props.widget.size).width,
  height: getWidgetSize(props.widget.size).height
}))

const getWidgetSize = (size: string) => {
  const sizes = {
    small: { width: '200px', height: '150px' },
    medium: { width: '300px', height: '200px' },
    large: { width: '400px', height: '300px' }
  }
  return sizes[size as keyof typeof sizes] || sizes.medium
}

const getPluginWidgetComponent = (_type: string) => {
  // 动态加载插件小部件组件
  return () => null
}

const handleDragStart = (event: DragEvent) => {
  isDragging.value = true
  dragOffset.value = {
    x: event.offsetX,
    y: event.offsetY
  }
  emit('drag-start', props.widget)
}

const handleDragEnd = (event: DragEvent) => {
  isDragging.value = false
  const newPosition = {
    x: event.clientX - dragOffset.value.x,
    y: event.clientY - dragOffset.value.y
  }
  emit('drag-end', props.widget, newPosition)
}

const handleClick = () => {
  // 某些小部件点击后可以打开对应的应用
  if (props.widget.config.appId) {
    emit('open-app', props.widget.config.appId)
  }
}

const handleEdit = () => {
  // 打开小部件编辑界面
  showControls.value = !showControls.value
}

const handleRemove = () => {
  emit('remove', props.widget.id)
}

const handleResizeStart = (event: MouseEvent) => {
  // 实现小部件调整大小功能
  const _startX = event.clientX
  const _startY = event.clientY
  const _startWidth = parseInt(getWidgetSize(props.widget.size).width)
  const _startHeight = parseInt(getWidgetSize(props.widget.size).height)

  const handleMouseMove = (e: MouseEvent) => {
    const _deltaX = e.clientX - _startX
    const _deltaY = e.clientY - _startY
    // 计算新大小
  }

  const handleMouseUp = () => {
    document.removeEventListener('mousemove', handleMouseMove)
    document.removeEventListener('mouseup', handleMouseUp)
  }

  document.addEventListener('mousemove', handleMouseMove)
  document.addEventListener('mouseup', handleMouseUp)
}

// 显示控制按钮
const handleMouseEnter = () => {
  showControls.value = true
}

const handleMouseLeave = () => {
  showControls.value = false
}
</script>

<style scoped>
.desktop-widget {
  position: absolute;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  cursor: move;
  user-select: none;
  transition: box-shadow 0.2s ease;
  overflow: hidden;
}

.desktop-widget:hover {
  box-shadow: 0 12px 48px rgba(0, 0, 0, 0.15);
}

.desktop-widget.is-dragging {
  opacity: 0.8;
  cursor: grabbing;
}

.widget-small {
  width: 200px;
  height: 150px;
}

.widget-medium {
  width: 300px;
  height: 200px;
}

.widget-large {
  width: 400px;
  height: 300px;
}

.widget-controls {
  position: absolute;
  top: 8px;
  right: 8px;
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.desktop-widget:hover .widget-controls {
  opacity: 1;
}

.control-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #6b7280;
  transition: all 0.2s ease;
}

.control-btn:hover {
  background: white;
  color: #1f2937;
}

.control-btn.danger:hover {
  background: #ef4444;
  color: white;
  border-color: #ef4444;
}

.resize-handle {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 16px;
  height: 16px;
  cursor: se-resize;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.resize-handle::after {
  content: '';
  position: absolute;
  bottom: 4px;
  right: 4px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 0 6px 6px;
  border-color: transparent transparent #6b7280 transparent;
}

.desktop-widget:hover .resize-handle {
  opacity: 1;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .desktop-widget {
    background: rgba(0, 0, 0, 0.8);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .control-btn {
    background: rgba(0, 0, 0, 0.8);
    color: #9ca3af;
  }
}
</style>