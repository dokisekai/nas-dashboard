<template>
  <div class="nt-container" :class="positionClasses">
    <transition-group name="nt-slide" tag="div" class="nt-list">
      <div
        v-for="notification in visibleNotifications"
        :key="notification.id"
        :class="['nt-item', `nt-${notification.type}`, { 'nt-read': notification.read }]"
      >
        <!-- 通知图标 -->
        <div class="nt-icon">
          <component :is="getIcon(notification.type)" class="w-6 h-6" />
        </div>

        <!-- 通知内容 -->
        <div class="nt-content">
          <div class="nt-header">
            <h4 class="nt-title">{{ notification.title }}</h4>
            <button class="nt-close" @click="handleClose(notification.id)">
              <XMarkIcon class="w-4 h-4" />
            </button>
          </div>

          <p class="nt-message">{{ notification.message }}</p>

          <!-- 通知操作 -->
          <div v-if="notification.actions && notification.actions.length > 0" class="nt-actions">
            <button
              v-for="action in notification.actions"
              :key="action.id"
              :class="['nt-action-btn', { primary: action.primary, destructive: action.destructive }]"
              @click="handleAction(notification.id, action)"
            >
              <component v-if="action.icon" :is="getActionIcon(action.icon)" class="w-4 h-4" />
              {{ action.label }}
            </button>
          </div>

          <!-- 通知元数据 -->
          <div class="nt-meta">
            <span class="nt-time">{{ formatTime(notification.timestamp) }}</span>
            <span class="nt-source">{{ notification.source }}</span>
          </div>
        </div>

        <!-- 进度条 -->
        <div v-if="!notification.persistent && notification.expires" class="nt-progress">
          <div
            class="nt-progress-bar"
            :style="{ animationDuration: getRemainingTime(notification) }"
          ></div>
        </div>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useNotificationStore } from '../../stores/notification'
import type { Notification, NotificationType } from './NotificationTypes'
import {
  InformationCircleIcon,
  ExclamationTriangleIcon,
  ExclamationCircleIcon,
  CheckCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

const notificationStore = useNotificationStore()

// 状态
const maxVisible = ref(5)

// 计算属性
const visibleNotifications = computed(() => {
  return notificationStore.notifications
    .filter(n => !n.read)
    .slice(0, maxVisible.value)
})

const positionClasses = computed(() => {
  return `nt-position-${notificationStore.settings.position}`
})

// 方法
const getIcon = (type: NotificationType) => {
  const icons: Record<NotificationType, any> = {
    info: InformationCircleIcon,
    warning: ExclamationTriangleIcon,
    error: ExclamationCircleIcon,
    success: CheckCircleIcon
  }
  return icons[type] || InformationCircleIcon
}

const getActionIcon = (iconName: string) => {
  // 这里可以根据需要映射图标
  return null
}

const formatTime = (timestamp: Date) => {
  const now = new Date()
  const diff = now.getTime() - timestamp.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (days > 0) {
    return `${days}天前`
  } else if (hours > 0) {
    return `${hours}小时前`
  } else if (minutes > 0) {
    return `${minutes}分钟前`
  } else {
    return '刚刚'
  }
}

const getRemainingTime = (notification: Notification) => {
  if (!notification.expires) return '0ms'
  const remaining = notification.expires.getTime() - Date.now()
  return `${Math.max(0, remaining)}ms`
}

const handleClose = (id: string) => {
  notificationStore.remove(id)
}

const handleAction = async (notificationId: string, action: any) => {
  try {
    await action.action()
    // 如果操作成功，移除通知
    notificationStore.remove(notificationId)
  } catch (error) {
    console.error('Notification action failed:', error)
  }
}

// 生命周期
onMounted(() => {
  // 监听设置变化
  watch(() => notificationStore.settings.maxVisible, (newValue) => {
    maxVisible.value = newValue
  })
})
</script>

<style scoped>
.nt-container {
  position: fixed;
  z-index: 9999;
  pointer-events: none;
}

.nt-position-top-right {
  top: 20px;
  right: 20px;
}

.nt-position-top-left {
  top: 20px;
  left: 20px;
}

.nt-position-bottom-right {
  bottom: 20px;
  right: 20px;
}

.nt-position-bottom-left {
  bottom: 20px;
  left: 20px;
}

.nt-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-width: 400px;
}

.nt-item {
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 16px;
  pointer-events: all;
  display: flex;
  gap: 12px;
  position: relative;
  overflow: hidden;
  min-width: 300px;
  max-width: 400px;
  border-left: 4px solid;
}

.nt-info {
  border-left-color: #3b82f6;
}

.nt-warning {
  border-left-color: #f59e0b;
}

.nt-error {
  border-left-color: #ef4444;
}

.nt-success {
  border-left-color: #10b981;
}

.nt-read {
  opacity: 0.7;
}

.nt-icon {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
}

.nt-info .nt-icon {
  color: #3b82f6;
}

.nt-warning .nt-icon {
  color: #f59e0b;
}

.nt-error .nt-icon {
  color: #ef4444;
}

.nt-success .nt-icon {
  color: #10b981;
}

.nt-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nt-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}

.nt-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
}

.nt-close {
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  padding: 2px;
  flex-shrink: 0;
  transition: color 0.2s ease;
}

.nt-close:hover {
  color: #6b7280;
}

.nt-message {
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
  margin: 0;
}

.nt-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.nt-action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.nt-action-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.nt-action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.nt-action-btn.primary:hover {
  background: #2563eb;
}

.nt-action-btn.destructive {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.nt-action-btn.destructive:hover {
  background: #dc2626;
}

.nt-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: #9ca3af;
}

.nt-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.nt-progress-bar {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #8b5cf6 100%);
  animation-name: progress;
  animation-timing-function: linear;
  animation-fill-mode: forwards;
}

@keyframes progress {
  from {
    width: 100%;
  }
  to {
    width: 0%;
  }
}

/* 动画效果 */
.nt-slide-enter-active {
  transition: all 0.3s ease;
}

.nt-slide-leave-active {
  transition: all 0.3s ease;
}

.nt-slide-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.nt-slide-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.nt-slide-move {
  transition: transform 0.3s ease;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .nt-item {
    background: #1f2937;
  }

  .nt-title {
    color: #f3f4f6;
  }

  .nt-message {
    color: #9ca3af;
  }

  .nt-meta {
    color: #6b7280;
  }

  .nt-action-btn {
    background: #374151;
    border-color: #4b5563;
    color: #9ca3af;
  }

  .nt-action-btn:hover {
    background: #4b5563;
    color: #f3f4f6;
  }
}
</style>