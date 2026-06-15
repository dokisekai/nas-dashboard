import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../api/client'

export interface Notification {
  id: string
  type: 'info' | 'warning' | 'error' | 'success'
  title: string
  message: string
  timestamp: Date | string
  read: boolean
  persistent?: boolean
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const isCenterOpen = ref(false)
  const loading = ref(false)

  const unreadCount = computed(() => notifications.value.filter(n => !n.read).length)

  const fetchNotifications = async () => {
    loading.value = true
    try {
      const data = await api.get('/api/notifications')
      notifications.value = (data as any).map((n: any) => ({
        ...n,
        id: String(n.id),
        timestamp: new Date(n.timestamp)
      }))
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
    } finally {
      loading.value = false
    }
  }

  const add = async (notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) => {
    // 首先在前端添加（为了响应速度）
    const tempId = Math.random().toString(36).substring(2, 9)
    const newNotification: Notification = {
      ...notification,
      id: tempId,
      timestamp: new Date(),
      read: false
    }
    notifications.value.unshift(newNotification)

    // 然后同步到后端
    try {
      await api.post('/api/notifications', notification)
      // 重新拉取以获取真实的 ID
      await fetchNotifications()
    } catch (error) {
      console.error('Failed to sync notification to backend:', error)
    }
  }

  const markAsRead = async (id: string) => {
    const n = notifications.value.find(n => n.id === id)
    if (n) {
      n.read = true
      try {
        await api.post(`/api/notifications/read/${id}`)
      } catch (error) {
        console.error('Failed to mark notification as read in backend:', error)
      }
    }
  }

  const markAllAsRead = async () => {
    notifications.value.forEach(n => n.read = true)
    try {
      await api.post('/api/notifications/read/all')
    } catch (error) {
      console.error('Failed to mark all as read in backend:', error)
    }
  }

  const remove = async (id: string) => {
    notifications.value = notifications.value.filter(n => n.id !== id)
    // 后端目前没有单条删除接口，暂不处理或添加
  }

  const clearAll = async () => {
    notifications.value = []
    try {
      await api.delete('/api/notifications/clear')
    } catch (error) {
      console.error('Failed to clear notifications in backend:', error)
    }
  }

  return {
    notifications,
    unreadCount,
    isCenterOpen,
    loading,
    fetchNotifications,
    add,
    markAsRead,
    markAllAsRead,
    remove,
    clearAll,
    showCenter: openCenter,
    hideCenter: closeCenter,
    openCenter,
    closeCenter,
    toggleCenter
  }

  function openCenter() { isCenterOpen.value = true }
  function closeCenter() { isCenterOpen.value = false }
  function toggleCenter() { isCenterOpen.value = !isCenterOpen.value }
})
