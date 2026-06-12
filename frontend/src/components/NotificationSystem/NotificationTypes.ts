// 通知系统类型定义

export type NotificationType = 'info' | 'warning' | 'error' | 'success'
export type NotificationCategory = 'system' | 'storage' | 'network' | 'security' | 'user' | 'app' | 'file'

export interface Notification {
  id: string
  type: NotificationType
  title: string
  message: string
  timestamp: Date
  read: boolean
  category: NotificationCategory
  actions?: NotificationAction[]
  expires?: Date
  priority: number
  source: string
  metadata?: Record<string, any>
  persistent?: boolean
}

export interface NotificationAction {
  id: string
  label: string
  action: () => void | Promise<void>
  primary?: boolean
  destructive?: boolean
  icon?: string
}

export interface NotificationSettings {
  enabled: boolean
  sound: boolean
  desktop: boolean
  position: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left'
  duration: number
  maxVisible: number
  categories: CategoryNotificationSettings
}

export interface CategoryNotificationSettings {
  system: boolean
  storage: boolean
  network: boolean
  security: boolean
  user: boolean
  app: boolean
  file: boolean
}

export interface NotificationFilter {
  types?: NotificationType[]
  categories?: NotificationCategory[]
  startDate?: Date
  endDate?: Date
  readStatus?: 'all' | 'read' | 'unread'
  searchQuery?: string
}

export interface NotificationStore {
  notifications: Notification[]
  unreadCount: number
  settings: NotificationSettings
  visible: boolean

  // 基础操作
  add(notification: Omit<Notification, 'id' | 'timestamp'>): void
  remove(id: string): void
  clear(category?: NotificationCategory): void
  markAsRead(id: string): void
  markAllAsRead(): void
  markAsUnread(id: string): void

  // 查询功能
  getFiltered(filter: NotificationFilter): Notification[]
  getUnread(): Notification[]
  getByCategory(category: NotificationCategory): Notification[]
  getByType(type: NotificationType): Notification[]

  // 设置管理
  updateSettings(settings: Partial<NotificationSettings>): void
  resetSettings(): void

  // UI控制
  showCenter(): void
  hideCenter(): void
  toggleCenter(): void
}

export interface NotificationEvent {
  type: 'added' | 'removed' | 'updated' | 'read'
  notification: Notification
  timestamp: Date
}

export interface NotificationCenterStats {
  total: number
  unread: number
  byType: Record<NotificationType, number>
  byCategory: Record<NotificationCategory, number>
  oldest: Date | null
  newest: Date | null
}