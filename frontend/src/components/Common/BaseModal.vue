<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="show"
        class="fixed inset-0 z-50 flex items-center justify-center"
        @click.self="handleCancel"
      >
        <!-- 背景遮罩 -->
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" />

        <!-- 模态框内容 -->
        <div
          class="relative w-full max-w-md bg-gray-800 rounded-2xl border border-gray-700 shadow-2xl mx-4 max-h-[90vh] overflow-hidden flex flex-col"
        >
          <!-- 标题栏 -->
          <div class="flex items-center justify-between px-6 py-4 border-b border-gray-700">
            <h3 class="text-lg font-semibold text-white">{{ title }}</h3>
            <button
              @click="handleCancel"
              class="p-1 text-gray-400 hover:text-white transition-colors rounded-lg hover:bg-gray-700"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- 内容区域 -->
          <div class="px-6 py-4 overflow-y-auto flex-1">
            <slot />
          </div>

          <!-- 操作按钮 -->
          <div class="flex items-center justify-end gap-3 px-6 py-4 border-t border-gray-700 bg-gray-900/50">
            <button
              @click="handleCancel"
              class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-white rounded-xl transition-colors"
            >
              取消
            </button>
            <button
              @click="handleConfirm"
              :disabled="confirming"
              class="px-5 py-2.5 bg-indigo-500 hover:bg-indigo-600 text-white rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="confirming">处理中...</span>
              <span v-else>{{ confirmText || '确定' }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { watch } from 'vue'

interface Props {
  show: boolean
  title?: string
  confirmText?: string
  confirming?: boolean
}

interface Emits {
  (e: 'update:show', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}

const props = withDefaults(defineProps<Props>(), {
  title: '对话框',
  confirmText: '确定',
  confirming: false
})

const emit = defineEmits<Emits>()

const handleConfirm = () => {
  if (!props.confirming) {
    emit('confirm')
  }
}

const handleCancel = () => {
  emit('update:show', false)
  emit('cancel')
}

// 监听 show 变化，ESC 键关闭
watch(() => props.show, (show) => {
  if (show) {
    const handleEsc = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        handleCancel()
      }
    }
    document.addEventListener('keydown', handleEsc)
    return () => document.removeEventListener('keydown', handleEsc)
  }
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from .relative,
.modal-leave-to .relative {
  transform: scale(0.95);
}
</style>
