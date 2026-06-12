<template>
  <div class="window-snap-overlay" :class="{ visible: isVisible }">
    <!-- 窗口吸附区域指示器 -->
    <div class="snap-zones">
      <!-- 左半屏 -->
      <div
        class="snap-zone snap-left"
        :class="{ active: activeZone === 'left' }"
        @drop="handleDrop('left')"
        @dragover.prevent="handleDragOver('left')"
        @dragleave="handleDragLeave"
      >
        <div class="zone-indicator">
          <ArrowLeftIcon class="w-8 h-8" />
        </div>
      </div>

      <!-- 右半屏 -->
      <div
        class="snap-zone snap-right"
        :class="{ active: activeZone === 'right' }"
        @drop="handleDrop('right')"
        @dragover.prevent="handleDragOver('right')"
        @dragleave="handleDragLeave"
      >
        <div class="zone-indicator">
          <ArrowRightIcon class="w-8 h-8" />
        </div>
      </div>

      <!-- 上半屏 -->
      <div
        class="snap-zone snap-top"
        :class="{ active: activeZone === 'top' }"
        @drop="handleDrop('top')"
        @dragover.prevent="handleDragOver('top')"
        @dragleave="handleDragLeave"
      >
        <div class="zone-indicator">
          <ArrowUpIcon class="w-8 h-8" />
        </div>
      </div>

      <!-- 下半屏 -->
      <div
        class="snap-zone snap-bottom"
        :class="{ active: activeZone === 'bottom' }"
        @drop="handleDrop('bottom')"
        @dragover.prevent="handleDragOver('bottom')"
        @dragleave="handleDragLeave"
      >
        <div class="zone-indicator">
          <ArrowDownIcon class="w-8 h-8" />
        </div>
      </div>

      <!-- 全屏 -->
      <div
        class="snap-zone snap-maximize"
        :class="{ active: activeZone === 'maximize' }"
        @drop="handleDrop('maximize')"
        @dragover.prevent="handleDragOver('maximize')"
        @dragleave="handleDragLeave"
      >
        <div class="zone-indicator">
          <ArrowsPointingOutIcon class="w-8 h-8" />
        </div>
      </div>

      <!-- 四分之一屏 -->
      <div
        class="snap-zone snap-top-left"
        :class="{ active: activeZone === 'top-left' }"
        @drop="handleDrop('top-left')"
        @dragover.prevent="handleDragOver('top-left')"
        @dragleave="handleDragLeave"
      ></div>

      <div
        class="snap-zone snap-top-right"
        :class="{ active: activeZone === 'top-right' }"
        @drop="handleDrop('top-right')"
        @dragover.prevent="handleDragOver('top-right')"
        @dragleave="handleDragLeave"
      ></div>

      <div
        class="snap-zone snap-bottom-left"
        :class="{ active: activeZone === 'bottom-left' }"
        @drop="handleDrop('bottom-left')"
        @dragover.prevent="handleDragOver('bottom-left')"
        @dragleave="handleDragLeave"
      ></div>

      <div
        class="snap-zone snap-bottom-right"
        :class="{ active: activeZone === 'bottom-right' }"
        @drop="handleDrop('bottom-right')"
        @dragover.prevent="handleDragOver('bottom-right')"
        @dragleave="handleDragLeave"
      ></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  ArrowLeftIcon,
  ArrowRightIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  ArrowsPointingOutIcon
} from '@heroicons/vue/24/outline'

interface Props {
  isVisible: boolean
}

interface Emits {
  (e: 'snap', zone: string): void
}

defineProps<Props>()
const emit = defineEmits<Emits>()

const activeZone = ref<string | null>(null)

const handleDragOver = (zone: string) => {
  activeZone.value = zone
}

const handleDragLeave = () => {
  activeZone.value = null
}

const handleDrop = (zone: string) => {
  emit('snap', zone)
  activeZone.value = null
}
</script>

<style scoped>
.window-snap-overlay {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 99999;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.window-snap-overlay.visible {
  opacity: 1;
  pointer-events: all;
}

.snap-zones {
  width: 100%;
  height: 100%;
  position: relative;
}

.snap-zone {
  position: absolute;
  background: transparent;
  border: 3px solid transparent;
  transition: all 0.2s ease;
  pointer-events: all;
}

.snap-zone.active {
  background: rgba(59, 130, 246, 0.2);
  border-color: #3b82f6;
}

.snap-left {
  left: 0;
  top: 0;
  width: 50%;
  height: 100%;
}

.snap-right {
  right: 0;
  top: 0;
  width: 50%;
  height: 100%;
}

.snap-top {
  top: 0;
  left: 0;
  width: 100%;
  height: 50%;
}

.snap-bottom {
  bottom: 0;
  left: 0;
  width: 100%;
  height: 50%;
}

.snap-maximize {
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 200px;
  height: 200px;
}

.snap-top-left {
  top: 0;
  left: 0;
  width: 50%;
  height: 50%;
}

.snap-top-right {
  top: 0;
  right: 0;
  width: 50%;
  height: 50%;
}

.snap-bottom-left {
  bottom: 0;
  left: 0;
  width: 50%;
  height: 50%;
}

.snap-bottom-right {
  bottom: 0;
  right: 0;
  width: 50%;
  height: 50%;
}

.zone-indicator {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: #3b82f6;
  opacity: 0.8;
}

.snap-zone:not(.active) .zone-indicator {
  display: none;
}
</style>
