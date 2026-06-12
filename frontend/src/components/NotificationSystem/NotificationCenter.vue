<template>
  <div v-if="visible" class="nc-overlay" @click="handleClose">
    <div class="nc-center" @click.stop>
      <!-- 头部 -->
      <div class="nc-header">
        <div class="nc-header-left">
          <h2 class="nc-title">通知中心</h2>
          <span class="nc-count">{{ unreadCount }} 条未读</span>
        </div>

        <div class="nc-header-actions">
          <button class="nc-btn nc-btn-ghost" @click="markAllAsRead" title="全部标记为已读">
            <CheckIcon class="w-5 h-5" />
          </button>
          <button class="nc-btn nc-btn-ghost" @click="showSettings = true" title="通知设置">
            <CogIcon class="w-5 h-5" />
          </button>
          <button class="nc-btn nc-btn-ghost" @click="handleClose" title="关闭">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- 过滤器 -->
      <div class="nc-filters">
        <button
          v-for="filter in filters"
          :key="filter.id"
          :class="['nc-filter-btn', { active: activeFilter === filter.id }]"
          @click="activeFilter = filter.id"
        >
          <component :is="getFilterIcon(filter.icon)" class="w-4 h-4" />
          {{ filter.label }}
          <span v-if="getCount(filter.id) > 0" class="nc-filter-count">
            {{ getCount(filter.id) }}
          </span>
        </button>

        <div class="nc-search">
          <MagnifyingGlassIcon class="w-4 h-4" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索通知..."
          />
        </div>
      </div>

      <!-- 通知列表 -->
      <div class="nc-content">
        <!-- 加载状态 -->
        <div v-if="loading" class="nc-loading">
          <div class="nc-spinner"></div>
          <p>加载通知...</p>
        </div>

        <!-- 空状态 -->
        <div v-else-if="filteredNotifications.length === 0" class="nc-empty">
          <BellIcon class="w-12 h-12" />
          <p>没有通知</p>
        </div>

        <!-- 通知列表 -->
        <div v-else class="nc-list">
          <div class="nc-group-label" v-if="todayNotifications.length > 0">
            今天
          </div>

          <div
            v-for="notification in todayNotifications"
            :key="notification.id"
            :class="['nc-item', `nc-${notification.type}`, { 'nc-read': notification.read }]"
            @click="handleNotificationClick(notification)"
          >
            <!-- 通知图标 -->
            <div class="nc-icon">
              <component :is="getNotificationIcon(notification.type)" class="w-5 h-5" />
            </div>

            <!-- 通知内容 -->
            <div class="nc-content">
              <div class="nc-header-row">
                <h4 class="nc-title">{{ notification.title }}</h4>
                <span class="nc-time">{{ formatTime(notification.timestamp) }}</span>
              </div>

              <p class="nc-message">{{ notification.message }}</p>

              <!-- 通知操作 -->
              <div v-if="notification.actions && notification.actions.length > 0" class="nc-actions">
                <button
                  v-for="action in notification.actions"
                  :key="action.id"
                  :class="['nc-action-btn', { primary: action.primary, destructive: action.destructive }]"
                  @click.stop="handleAction(notification.id, action)"
                >
                  {{ action.label }}
                </button>
              </div>

              <!-- 通知来源 -->
              <div class="nc-meta">
                <span class="nc-source">{{ notification.source }}</span>
                <span class="nc-category">{{ getCategoryLabel(notification.category) }}</span>
              </div>
            </div>

            <!-- 通知操作菜单 -->
            <div class="nc-item-actions">
              <button
                class="nc-btn nc-btn-ghost"
                @click.stop="toggleRead(notification.id)"
                :title="notification.read ? '标记为未读' : '标记为已读'"
              >
                <EnvelopeIcon v-if="!notification.read" class="w-4 h-4" />
                <EnvelopeOpenIcon v-else class="w-4 h-4" />
              </button>
              <button
                class="nc-btn nc-btn-ghost"
                @click.stop="removeNotification(notification.id)"
                title="删除"
              >
                <TrashIcon class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div class="nc-group-label" v-if="olderNotifications.length > 0">
            更早
          </div>

          <div
            v-for="notification in olderNotifications"
            :key="notification.id"
            :class="['nc-item', `nc-${notification.type}`, { 'nc-read': notification.read }]"
            @click="handleNotificationClick(notification)"
          >
            <!-- 同上 -->
          </div>
        </div>
      </div>

      <!-- 底部操作 -->
      <div class="nc-footer">
        <button class="nc-btn nc-btn-secondary" @click="clearAll">
          <TrashIcon class="w-4 h-4" />
          清空所有
        </button>
      </div>
    </div>

    <!-- 通知设置模态框 -->
    <div v-if="showSettings" class="nc-modal-overlay" @click="showSettings = false">
      <div class="nc-modal" @click.stop>
        <div class="nc-modal-header">
          <h3>通知设置</h3>
          <button class="nc-btn nc-btn-ghost" @click="showSettings = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="nc-modal-body">
          <div class="nc-setting-group">
            <h4>基本设置</h4>

            <div class="nc-setting-item">
              <label>启用通知</label>
              <button
                :class="['nc-switch', { active: settings.enabled }]"
                @click="updateSetting('enabled', !settings.enabled)"
              >
                <div class="nc-switch-slider"></div>
              </button>
            </div>

            <div class="nc-setting-item">
              <label>通知声音</label>
              <button
                :class="['nc-switch', { active: settings.sound }]"
                @click="updateSetting('sound', !settings.sound)"
              >
                <div class="nc-switch-slider"></div>
              </button>
            </div>

            <div class="nc-setting-item">
              <label>桌面通知</label>
              <button
                :class="['nc-switch', { active: settings.desktop }]"
                @click="updateSetting('desktop', !settings.desktop)"
              >
                <div class="nc-switch-slider"></div>
              </button>
            </div>

            <div class="nc-setting-item">
              <label>显示位置</label>
              <select
                :value="settings.position"
                @change="updateSetting('position', ($event.target as HTMLSelectElement).value)"
                class="nc-select"
              >
                <option value="top-right">右上角</option>
                <option value="top-left">左上角</option>
                <option value="bottom-right">右下角</option>
                <option value="bottom-left">左下角</option>
              </select>
            </div>

            <div class="nc-setting-item">
              <label>显示时长 (秒)</label>
              <input
                type="number"
                :value="settings.duration / 1000"
                @change="updateSetting('duration', parseInt(($event.target as HTMLInputElement).value) * 1000)"
                class="nc-input"
                min="0"
                max="60"
              />
            </div>

            <div class="nc-setting-item">
              <label>最大显示数量</label>
              <input
                type="number"
                :value="settings.maxVisible"
                @change="updateSetting('maxVisible', parseInt(($event.target as HTMLInputElement).value))"
                class="nc-input"
                min="1"
                max="10"
              />
            </div>
          </div>

          <div class="nc-setting-group">
            <h4>分类设置</h4>

            <div class="nc-setting-item" v-for="(enabled, category) in settings.categories" :key="category">
              <label>{{ getCategoryLabel(category as any) }}</label>
              <button
                :class="['nc-switch', { active: enabled }]"
                @click="updateCategorySetting(category, !enabled)"
              >
                <div class="nc-switch-slider"></div>
              </button>
            </div>
          </div>
        </div>

        <div class="nc-modal-footer">
          <button class="nc-btn nc-btn-secondary" @click="resetSettings">
            重置默认
          </button>
          <button class="nc-btn nc-btn-primary" @click="showSettings = false">
            完成
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useNotificationStore } from '../../stores/notification'
import type { Notification, NotificationType, NotificationCategory } from './NotificationTypes'
import {
  BellIcon,
  CheckIcon,
  CogIcon,
  XMarkIcon,
  MagnifyingGlassIcon,
  InformationCircleIcon,
  ExclamationTriangleIcon,
  ExclamationCircleIcon,
  CheckCircleIcon,
  EnvelopeIcon,
  EnvelopeOpenIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'

const notificationStore = useNotificationStore()

// 状态
const loading = ref(false)
const activeFilter = ref('all')
const searchQuery = ref('')
const showSettings = ref(false)

// 过滤器配置
const filters = [
  { id: 'all', label: '全部', icon: 'BellIcon' },
  { id: 'unread', label: '未读', icon: 'EnvelopeIcon' },
  { id: 'info', label: '信息', icon: 'InformationCircleIcon' },
  { id: 'warning', label: '警告', icon: 'ExclamationTriangleIcon' },
  { id: 'error', label: '错误', icon: 'ExclamationCircleIcon' },
  { id: 'success', label: '成功', icon: 'CheckCircleIcon' }
]

// 计算属性
const { visible, notifications, settings, unreadCount } = notificationStore

const filteredNotifications = computed(() => {
  let filtered = [...notifications.value]

  // 应用过滤器
  if (activeFilter.value === 'unread') {
    filtered = filtered.filter(n => !n.read)
  } else if (activeFilter.value !== 'all') {
    filtered = filtered.filter(n => n.type === activeFilter.value)
  }

  // 应用搜索
  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(n =>
      n.title.toLowerCase().includes(query) ||
      n.message.toLowerCase().includes(query)
    )
  }

  return filtered
})

const todayNotifications = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  return filteredNotifications.value.filter(n => n.timestamp >= today)
})

const olderNotifications = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  return filteredNotifications.value.filter(n => n.timestamp < today)
})

// 方法
const getFilterIcon = (iconName: string) => {
  const iconMap: Record<string, any> = {
    BellIcon,
    EnvelopeIcon,
    InformationCircleIcon,
    ExclamationTriangleIcon,
    ExclamationCircleIcon,
    CheckCircleIcon
  }
  return iconMap[iconName] || BellIcon
}

const getNotificationIcon = (type: NotificationType) => {
  const icons: Record<NotificationType, any> = {
    info: InformationCircleIcon,
    warning: ExclamationTriangleIcon,
    error: ExclamationCircleIcon,
    success: CheckCircleIcon
  }
  return icons[type] || InformationCircleIcon
}

const getCount = (filterId: string) => {
  if (filterId === 'all') return notifications.value.length
  if (filterId === 'unread') return unreadCount.value

  return notifications.value.filter(n => n.type === filterId).length
}

const getCategoryLabel = (category: NotificationCategory) => {
  const labels: Record<NotificationCategory, string> = {
    system: '系统',
    storage: '存储',
    network: '网络',
    security: '安全',
    user: '用户',
    app: '应用',
    file: '文件'
  }
  return labels[category] || category
}

const formatTime = (timestamp: Date) => {
  const now = new Date()
  const diff = now.getTime() - timestamp.getTime()
  const seconds = Math.floor(diff / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)

  if (hours > 0) {
    return `${hours}小时前`
  } else if (minutes > 0) {
    return `${minutes}分钟前`
  } else {
    return '刚刚'
  }
}

const handleClose = () => {
  notificationStore.hideCenter()
}

const handleNotificationClick = (notification: Notification) => {
  if (!notification.read) {
    notificationStore.markAsRead(notification.id)
  }
}

const toggleRead = (id: string) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification) {
    if (notification.read) {
      notificationStore.markAsUnread(id)
    } else {
      notificationStore.markAsRead(id)
    }
  }
}

const removeNotification = (id: string) => {
  notificationStore.remove(id)
}

const markAllAsRead = () => {
  notificationStore.markAllAsRead()
}

const clearAll = () => {
  if (confirm('确定要清空所有通知吗？')) {
    notificationStore.clear()
  }
}

const handleAction = async (notificationId: string, action: any) => {
  try {
    await action.action()
    notificationStore.remove(notificationId)
  } catch (error) {
    console.error('Notification action failed:', error)
  }
}

const updateSetting = (key: string, value: any) => {
  notificationStore.updateSettings({ [key]: value })
}

const updateCategorySetting = (category: string, value: boolean) => {
  notificationStore.updateSettings({
    categories: {
      ...settings.value.categories,
      [category]: value
    }
  })
}

const resetSettings = () => {
  if (confirm('确定要重置通知设置吗？')) {
    notificationStore.resetSettings()
  }
}
</script>

<style scoped>
.nc-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10000;
}

.nc-center {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.nc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.nc-header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nc-title {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.nc-count {
  padding: 4px 12px;
  background: #ef4444;
  color: white;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.nc-header-actions {
  display: flex;
  gap: 8px;
}

.nc-filters {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 16px 24px;
  border-bottom: 1px solid #e5e7eb;
  overflow-x: auto;
}

.nc-filter-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.nc-filter-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.nc-filter-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.nc-filter-count {
  padding: 2px 6px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  font-size: 10px;
}

.nc-search {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f3f4f6;
  border-radius: 8px;
  min-width: 200px;
  margin-left: auto;
  color: #6b7280;
}

.nc-search input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 12px;
  outline: none;
}

.nc-content {
  flex: 1;
  overflow-y: auto;
  padding: 0;
}

.nc-loading,
.nc-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #9ca3af;
}

.nc-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 12px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.nc-list {
  padding: 16px 24px;
}

.nc-group-label {
  padding: 12px 0;
  font-size: 14px;
  font-weight: 600;
  color: #6b7280;
  border-bottom: 1px solid #f3f4f6;
}

.nc-item {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  border-left: 4px solid;
}

.nc-item:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-1px);
}

.nc-info {
  border-left-color: #3b82f6;
}

.nc-warning {
  border-left-color: #f59e0b;
}

.nc-error {
  border-left-color: #ef4444;
}

.nc-success {
  border-left-color: #10b981;
}

.nc-read {
  opacity: 0.7;
}

.nc-icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
}

.nc-info .nc-icon {
  color: #3b82f6;
}

.nc-warning .nc-icon {
  color: #f59e0b;
}

.nc-error .nc-icon {
  color: #ef4444;
}

.nc-success .nc-icon {
  color: #10b981;
}

.nc-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.nc-header-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 8px;
}

.nc-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.nc-time {
  font-size: 11px;
  color: #9ca3af;
  flex-shrink: 0;
}

.nc-message {
  font-size: 13px;
  color: #6b7280;
  line-height: 1.5;
  margin: 0;
}

.nc-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.nc-action-btn {
  padding: 4px 10px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.nc-action-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.nc-action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.nc-action-btn.destructive {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.nc-meta {
  display: flex;
  gap: 12px;
  font-size: 11px;
  color: #9ca3af;
}

.nc-item-actions {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nc-footer {
  display: flex;
  justify-content: center;
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
}

.nc-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.nc-btn-ghost {
  background: transparent;
  color: #6b7280;
  padding: 6px;
}

.nc-btn-ghost:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.nc-btn-secondary {
  background: white;
  border: 1px solid #e5e7eb;
  color: #6b7280;
}

.nc-btn-secondary:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.nc-btn-primary {
  background: #3b82f6;
  color: white;
}

.nc-btn-primary:hover {
  background: #2563eb;
}

.nc-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10001;
}

.nc-modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.nc-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.nc-modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
}

.nc-modal-body {
  padding: 20px;
}

.nc-setting-group {
  margin-bottom: 24px;
}

.nc-setting-group h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.nc-setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f3f4f6;
}

.nc-setting-item label {
  font-size: 14px;
  color: #6b7280;
}

.nc-switch {
  width: 44px;
  height: 22px;
  background: #e5e7eb;
  border-radius: 11px;
  position: relative;
  cursor: pointer;
  transition: background 0.2s ease;
  border: none;
  padding: 0;
}

.nc-switch.active {
  background: #3b82f6;
}

.nc-switch-slider {
  width: 18px;
  height: 18px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.nc-switch.active .nc-switch-slider {
  transform: translateX(22px);
}

.nc-select,
.nc-input {
  padding: 6px 10px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 12px;
  min-width: 100px;
}

.nc-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
}
</style>