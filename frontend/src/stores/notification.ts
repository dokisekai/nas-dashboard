import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Notification {
  id: string
  type: 'info' | 'warning' | 'error' | 'success'
  title: string
  message: string
  timestamp: Date
  read: boolean
  persistent?: boolean
}

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref<Notification[]>([])
  const isCenterOpen = ref(false)

  const unreadCount = computed(() => notifications.value.filter(n => !n.read).length)

  const add = (notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) => {
    const id = Math.random().toString(36).substring(2, 9)
    notifications.value.unshift({
      ...notification,
      id,
      timestamp: new Date(),
      read: false
    })
  }

  const showCenter = () => {
    isCenterOpen.value = true
  }

  const hideCenter = () => {
    isCenterOpen.value = false
  }

  const openCenter = () => {
    isCenterOpen.value = true
  }

  const closeCenter = () => {
    isCenterOpen.value = false
  }

  const toggleCenter = () => {
    isCenterOpen.value = !isCenterOpen.value
  }

  const markAsRead = (id: string) => {
    const n = notifications.value.find(n => n.id === id)
    if (n) n.read = true
  }

  const markAllAsRead = () => {
    notifications.value.forEach(n => n.read = true)
  }

  const remove = (id: string) => {
    notifications.value = notifications.value.filter(n => n.id !== id)
  }

  const clearAll = () => {
    notifications.value = []
  }

  return {
    notifications,
    unreadCount,
    isCenterOpen,
    add,
    markAsRead,
    markAllAsRead,
    remove,
    clearAll,
    showCenter,
    hideCenter,
    openCenter,
    closeCenter,
    toggleCenter
  }
})
