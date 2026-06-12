<template>
  <div class="notification-toasts">
    <TransitionGroup name="toast">
      <div
        v-for="toast in activeToasts"
        :key="toast.id"
        class="toast-item"
        :class="toast.type"
        @click="removeToast(toast.id)"
      >
        <div class="toast-icon">
          <InformationCircleIcon v-if="toast.type === 'info'" class="w-5 h-5" />
          <ExclamationTriangleIcon v-if="toast.type === 'warning'" class="w-5 h-5" />
          <XCircleIcon v-if="toast.type === 'error'" class="w-5 h-5" />
          <CheckCircleIcon v-if="toast.type === 'success'" class="w-5 h-5" />
        </div>
        <div class="toast-content">
          <div class="toast-title">{{ toast.title }}</div>
          <div class="toast-message">{{ toast.message }}</div>
        </div>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import {
  InformationCircleIcon,
  ExclamationTriangleIcon,
  XCircleIcon,
  CheckCircleIcon
} from '@heroicons/vue/24/outline'
import { useNotificationStore, type Notification } from '../../stores/notification'

const notificationStore = useNotificationStore()
const activeToasts = ref<Notification[]>([])

// Watch for new notifications
watch(() => notificationStore.notifications, (newVal) => {
  const latest = newVal[0]
  if (latest && !latest.read) {
    activeToasts.value.push(latest)
    setTimeout(() => {
      removeToast(latest.id)
    }, 5000)
  }
}, { deep: true })

const removeToast = (id: string) => {
  activeToasts.value = activeToasts.value.filter(t => t.id !== id)
}
</script>

<style scoped>
.notification-toasts {
  position: fixed;
  top: 60px;
  right: 20px;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 12px;
  pointer-events: none;
}

.toast-item {
  pointer-events: auto;
  width: 300px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  display: flex;
  gap: 12px;
  cursor: pointer;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.toast-icon {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.info .toast-icon { background: #eff6ff; color: #3b82f6; }
.warning .toast-icon { background: #fffbeb; color: #d97706; }
.error .toast-icon { background: #fef2f2; color: #dc2626; }
.success .toast-icon { background: #f0fdf4; color: #10b981; }

.toast-content {
  flex: 1;
}

.toast-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.toast-message {
  font-size: 13px;
  color: #6b7280;
  margin-top: 2px;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.toast-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
