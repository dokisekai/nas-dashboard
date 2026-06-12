import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  NotificationStore,
  Notification,
  NotificationSettings,
  NotificationFilter,
  NotificationCategory,
  NotificationType,
  NotificationCenterStats
} from '../components/NotificationSystem/NotificationTypes'

export const useNotificationStore = defineStore('notification', () => {
  // 状态
  const notifications = ref<Notification[]>([])
  const settings = ref<NotificationSettings>({
    enabled: true,
    sound: true,
    desktop: true,
    position: 'top-right',
    duration: 5000,
    maxVisible: 5,
    categories: {
      system: true,
      storage: true,
      network: true,
      security: true,
      user: true,
      app: true,
      file: true
    }
  })
  const visible = ref(false)

  // 计算属性
  const unreadCount = computed(() => {
    return notifications.value.filter(n => !n.read).length
  })

  // 生成唯一ID
  const generateId = (): string => {
    return `notif-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`
  }

  // 添加通知
  const add = (notification: Omit<Notification, 'id' | 'timestamp'>) => {
    if (!settings.value.enabled) return

    // 检查分类是否启用
    if (!settings.value.categories[notification.category]) return

    const newNotification: Notification = {
      ...notification,
      id: generateId(),
      timestamp: new Date()
    }

    notifications.value.unshift(newNotification)

    // 自动播放声音
    if (settings.value.sound && !notification.read) {
      playNotificationSound()
    }

    // 自动过期处理
    if (newNotification.expires && !newNotification.persistent) {
      setTimeout(() => {
        remove(newNotification.id)
      }, new Date(newNotification.expires).getTime() - Date.now())
    }

    // 自动移除非持久通知
    if (!notification.persistent && settings.value.duration > 0) {
      setTimeout(() => {
        remove(newNotification.id)
      }, settings.value.duration)
    }

    return newNotification
  }

  // 移除通知
  const remove = (id: string) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value.splice(index, 1)
    }
  }

  // 清空通知
  const clear = (category?: NotificationCategory) => {
    if (category) {
      notifications.value = notifications.value.filter(n => n.category !== category)
    } else {
      notifications.value = []
    }
  }

  // 标记为已读
  const markAsRead = (id: string) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.read = true
    }
  }

  // 标记所有为已读
  const markAllAsRead = () => {
    notifications.value.forEach(n => n.read = true)
  }

  // 标记为未读
  const markAsUnread = (id: string) => {
    const notification = notifications.value.find(n => n.id === id)
    if (notification) {
      notification.read = false
    }
  }

  // 获取过滤通知
  const getFiltered = (filter: NotificationFilter): Notification[] => {
    let filtered = [...notifications.value]

    // 类型过滤
    if (filter.types && filter.types.length > 0) {
      filtered = filtered.filter(n => filter.types!.includes(n.type))
    }

    // 分类过滤
    if (filter.categories && filter.categories.length > 0) {
      filtered = filtered.filter(n => filter.categories!.includes(n.category))
    }

    // 日期过滤
    if (filter.startDate) {
      filtered = filtered.filter(n => n.timestamp >= filter.startDate!)
    }

    if (filter.endDate) {
      filtered = filtered.filter(n => n.timestamp <= filter.endDate!)
    }

    // 读取状态过滤
    if (filter.readStatus === 'read') {
      filtered = filtered.filter(n => n.read)
    } else if (filter.readStatus === 'unread') {
      filtered = filtered.filter(n => !n.read)
    }

    // 搜索查询
    if (filter.searchQuery) {
      const query = filter.searchQuery.toLowerCase()
      filtered = filtered.filter(n =>
        n.title.toLowerCase().includes(query) ||
        n.message.toLowerCase().includes(query)
      )
    }

    return filtered
  }

  // 获取未读通知
  const getUnread = (): Notification[] => {
    return notifications.value.filter(n => !n.read)
  }

  // 按分类获取
  const getByCategory = (category: NotificationCategory): Notification[] => {
    return notifications.value.filter(n => n.category === category)
  }

  // 按类型获取
  const getByType = (type: NotificationType): Notification[] => {
    return notifications.value.filter(n => n.type === type)
  }

  // 更新设置
  const updateSettings = (newSettings: Partial<NotificationSettings>) => {
    settings.value = { ...settings.value, ...newSettings }
  }

  // 重置设置
  const resetSettings = () => {
    settings.value = {
      enabled: true,
      sound: true,
      desktop: true,
      position: 'top-right',
      duration: 5000,
      maxVisible: 5,
      categories: {
        system: true,
        storage: true,
        network: true,
        security: true,
        user: true,
        app: true,
        file: true
      }
    }
  }

  // UI控制
  const showCenter = () => {
    visible.value = true
  }

  const hideCenter = () => {
    visible.value = false
  }

  const toggleCenter = () => {
    visible.value = !visible.value
  }

  // 播放通知声音
  const playNotificationSound = () => {
    const audio = new Audio('/sounds/notification.mp3')
    audio.play().catch(console.error)
  }

  // 获取统计信息
  const getStats = (): NotificationCenterStats => {
    const stats: NotificationCenterStats = {
      total: notifications.value.length,
      unread: unreadCount.value,
      byType: {
        info: 0,
        warning: 0,
        error: 0,
        success: 0
      },
      byCategory: {
        system: 0,
        storage: 0,
        network: 0,
        security: 0,
        user: 0,
        app: 0,
        file: 0
      },
      oldest: null,
      newest: null
    }

    notifications.value.forEach(notification => {
      // 按类型统计
      stats.byType[notification.type]++

      // 按分类统计
      stats.byCategory[notification.category]++

      // 时间统计
      if (!stats.oldest || notification.timestamp < stats.oldest) {
        stats.oldest = notification.timestamp
      }
      if (!stats.newest || notification.timestamp > stats.newest) {
        stats.newest = notification.timestamp
      }
    })

    return stats
  }

  // 预设通知模板
  const createSystemNotification = (title: string, message: string, type: NotificationType = 'info') => {
    return add({
      type,
      title,
      message,
      read: false,
      category: 'system',
      priority: 5,
      source: 'system',
      persistent: type === 'error'
    })
  }

  const createStorageNotification = (title: string, message: string, type: NotificationType = 'info') => {
    return add({
      type,
      title,
      message,
      read: false,
      category: 'storage',
      priority: 4,
      source: 'storage'
    })
  }

  const createSecurityNotification = (title: string, message: string) => {
    return add({
      type: 'warning',
      title,
      message,
      read: false,
      category: 'security',
      priority: 8,
      source: 'security',
      persistent: true
    })
  }

  const createUserNotification = (title: string, message: string) => {
    return add({
      type: 'info',
      title,
      message,
      read: false,
      category: 'user',
      priority: 3,
      source: 'user'
    })
  }

  return {
    // 状态
    notifications,
    unreadCount,
    settings,
    visible,

    // 操作方法
    add,
    remove,
    clear,
    markAsRead,
    markAllAsRead,
    markAsUnread,
    getFiltered,
    getUnread,
    getByCategory,
    getByType,
    updateSettings,
    resetSettings,
    showCenter,
    hideCenter,
    toggleCenter,
    getStats,

    // 预设方法
    createSystemNotification,
    createStorageNotification,
    createSecurityNotification,
    createUserNotification
  }
})