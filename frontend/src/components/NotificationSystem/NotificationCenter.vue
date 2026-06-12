<template>
  <Transition name="slide">
    <div v-if="isOpen" class="notification-center">
      <div class="nc-header">
        <h3>通知中心</h3>
        <div class="nc-actions">
          <button @click="notificationStore.markAllAsRead" title="全部标记为已读">
            <CheckIcon class="w-5 h-5" />
          </button>
          <button @click="notificationStore.clearAll" title="清空全部">
            <TrashIcon class="w-5 h-5" />
          </button>
          <button @click="$emit('close')" title="关闭">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </div>

      <div class="nc-body">
        <div v-if="notifications.length === 0" class="empty-state">
          <BellIcon class="w-12 h-12 opacity-20" />
          <p>暂无新通知</p>
        </div>
        <div v-else class="notification-list">
          <div
            v-for="n in notifications"
            :key="n.id"
            class="notification-card"
            :class="[n.type, { unread: !n.read }]"
            @click="notificationStore.markAsRead(n.id)"
          >
            <div class="nc-icon">
              <InformationCircleIcon v-if="n.type === 'info'" class="w-5 h-5" />
              <ExclamationTriangleIcon v-if="n.type === 'warning'" class="w-5 h-5" />
              <XCircleIcon v-if="n.type === 'error'" class="w-5 h-5" />
              <CheckCircleIcon v-if="n.type === 'success'" class="w-5 h-5" />
            </div>
            <div class="nc-content">
              <div class="nc-title">{{ n.title }}</div>
              <div class="nc-message">{{ n.message }}</div>
              <div class="nc-time">{{ formatTime(n.timestamp) }}</div>
            </div>
            <button class="nc-remove" @click.stop="notificationStore.remove(n.id)">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  XMarkIcon,
  TrashIcon,
  CheckIcon,
  BellIcon,
  InformationCircleIcon,
  ExclamationTriangleIcon,
  XCircleIcon,
  CheckCircleIcon
} from '@heroicons/vue/24/outline'
import { useNotificationStore } from '../../stores/notification'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits(['close'])

const notificationStore = useNotificationStore()
const notifications = computed(() => notificationStore.notifications)

const formatTime = (date: Date) => {
  return new Intl.DateTimeFormat('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).format(date)
}
</script>

<style scoped>
.notification-center {
  position: fixed;
  top: 48px;
  right: 0;
  bottom: 0;
  width: 360px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-left: 1px solid rgba(0, 0, 0, 0.1);
  box-shadow: -10px 0 30px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  z-index: 999;
}

.nc-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nc-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.nc-actions {
  display: flex;
  gap: 8px;
}

.nc-actions button {
  padding: 6px;
  border-radius: 6px;
  color: #6b7280;
  transition: all 0.2s;
}

.nc-actions button:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.nc-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.empty-state {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  gap: 12px;
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-card {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: white;
  border-radius: 12px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  position: relative;
  cursor: pointer;
  transition: all 0.2s;
}

.notification-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.notification-card.unread::before {
  content: '';
  position: absolute;
  top: 12px;
  right: 12px;
  width: 8px;
  height: 8px;
  background: #3b82f6;
  border-radius: 50%;
}

.nc-icon {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.info .nc-icon { background: #eff6ff; color: #3b82f6; }
.warning .nc-icon { background: #fffbeb; color: #d97706; }
.error .nc-icon { background: #fef2f2; color: #dc2626; }
.success .nc-icon { background: #f0fdf4; color: #10b981; }

.nc-content {
  flex: 1;
}

.nc-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.nc-message {
  font-size: 13px;
  color: #6b7280;
  line-height: 1.4;
}

.nc-time {
  font-size: 11px;
  color: #9ca3af;
  margin-top: 8px;
}

.nc-remove {
  position: absolute;
  top: 10px;
  right: 10px;
  opacity: 0;
  transition: opacity 0.2s;
  color: #9ca3af;
}

.notification-card:hover .nc-remove {
  opacity: 1;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(100%);
}
</style>
