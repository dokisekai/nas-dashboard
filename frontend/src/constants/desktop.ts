export interface DockItem {
  id: string
  label: string
  icon: string
  appId?: string
  badge?: number | null
  running?: boolean
}

export const defaultDockItems: DockItem[] = [
  {
    id: 'app-center',
    label: '应用中心',
    icon: 'ShoppingBagIcon',
    appId: 'app-center',
    badge: null
  },
  {
    id: 'storage-manager',
    label: '存储管理',
    icon: 'ServerIcon',
    appId: 'storage-manager',
    badge: null
  },
  {
    id: 'system-monitor',
    label: '系统监控',
    icon: 'ChartBarIcon',
    appId: 'system-monitor',
    badge: null
  },
  {
    id: 'file-manager',
    label: '文件管理',
    icon: 'FolderIcon',
    appId: 'file-manager',
    badge: null
  },
  {
    id: 'user-manager',
    label: '用户管理',
    icon: 'UserGroupIcon',
    appId: 'user-manager',
    badge: null
  },
  {
    id: 'sync-manager',
    label: '同步备份',
    icon: 'CloudArrowUpIcon',
    appId: 'sync-manager',
    badge: null
  },
  {
    id: 'plugin-store',
    label: '插件商店',
    icon: 'ShoppingBagIcon',
    appId: 'plugin-store',
    badge: 3
  }
]

export interface WidgetDefinition {
  id: string
  name: string
  type: string
  component: string
  defaultSize: 'small' | 'medium' | 'large'
  defaultConfig: Record<string, any>
  description?: string
  category?: string
}

export const availableWidgets: WidgetDefinition[] = [
  {
    id: 'system-monitor',
    name: '系统监控',
    type: 'system-monitor',
    component: 'SystemMonitorWidget',
    defaultSize: 'large',
    defaultConfig: {
      showCpu: true,
      showMemory: true,
      showDisk: true,
      showNetwork: true
    },
    description: '显示CPU、内存、磁盘和网络使用情况',
    category: 'system'
  },
  {
    id: 'storage-status',
    name: '存储状态',
    type: 'storage-status',
    component: 'StorageStatusWidget',
    defaultSize: 'medium',
    defaultConfig: {
      showRaid: true,
      showTemperature: true,
      refreshInterval: 5000
    },
    description: '显示存储卷状态和使用情况',
    category: 'storage'
  },
  {
    id: 'network-monitor',
    name: '网络监控',
    type: 'network-monitor',
    component: 'NetworkMonitorWidget',
    defaultSize: 'medium',
    defaultConfig: {
      interface: 'all',
      showGraph: true,
      showSpeed: true
    },
    description: '显示网络流量和连接状态',
    category: 'network'
  },
  {
    id: 'clock',
    name: '时钟',
    type: 'clock',
    component: 'ClockWidget',
    defaultSize: 'small',
    defaultConfig: {
      showDate: true,
      showSeconds: true,
      format24: true
    },
    description: '显示当前时间和日期',
    category: 'utility'
  },
  {
    id: 'weather',
    name: '天气',
    type: 'weather',
    component: 'WeatherWidget',
    defaultSize: 'small',
    defaultConfig: {
      location: 'auto',
      showForecast: true,
      unit: 'celsius'
    },
    description: '显示天气信息',
    category: 'utility'
  },
  {
    id: 'quick-shortcuts',
    name: '快捷方式',
    type: 'quick-shortcuts',
    component: 'QuickShortcutsWidget',
    defaultSize: 'medium',
    defaultConfig: {
      maxShortcuts: 8,
      editable: true
    },
    description: '常用功能快捷方式',
    category: 'launcher'
  },
  {
    id: 'calendar',
    name: '日历',
    type: 'calendar',
    component: 'CalendarWidget',
    defaultSize: 'medium',
    defaultConfig: {
      showEvents: true,
      showHolidays: true
    },
    description: '显示日历和待办事项',
    category: 'utility'
  },
  {
    id: 'quick-note',
    name: '快速笔记',
    type: 'quick-note',
    component: 'QuickNoteWidget',
    defaultSize: 'medium',
    defaultConfig: {
      autoSave: true,
      maxLength: 500
    },
    description: '快速记录笔记和想法',
    category: 'productivity'
  }
]

export interface AppDefinition {
  id: string
  name: string
  component: string
  icon: string
  description?: string
  category?: string
  windowConfig?: {
    width: number
    height: number
    minWidth?: number
    minHeight?: number
    resizable?: boolean
    maximizable?: boolean
  }
  permissions?: string[]
}

export const availableApps: AppDefinition[] = [
  {
    id: 'app-center',
    name: '应用中心',
    component: 'AppCenter',
    icon: 'ShoppingBagIcon',
    description: '浏览和安装应用程序',
    category: 'system',
    windowConfig: {
      width: 900,
      height: 600,
      minWidth: 600,
      minHeight: 400,
      resizable: true,
      maximizable: true
    },
    permissions: ['app.read', 'app.install']
  },
  {
    id: 'storage-manager',
    name: '存储管理',
    component: 'StorageManager',
    icon: 'ServerIcon',
    description: '管理磁盘和存储卷',
    category: 'storage',
    windowConfig: {
      width: 1024,
      height: 700,
      minWidth: 800,
      minHeight: 500,
      resizable: true,
      maximizable: true
    },
    permissions: ['storage.read', 'storage.write']
  },
  {
    id: 'system-monitor',
    name: '系统监控',
    component: 'SystemMonitorApp',
    icon: 'ChartBarIcon',
    description: '实时系统资源监控',
    category: 'system',
    windowConfig: {
      width: 1200,
      height: 800,
      minWidth: 900,
      minHeight: 600,
      resizable: true,
      maximizable: true
    },
    permissions: ['system.read']
  },
  {
    id: 'file-manager',
    name: '文件管理',
    component: 'FileManager',
    icon: 'FolderIcon',
    description: '文件浏览和管理',
    category: 'files',
    windowConfig: {
      width: 1024,
      height: 700,
      minWidth: 800,
      minHeight: 500,
      resizable: true,
      maximizable: true
    },
    permissions: ['files.read', 'files.write']
  },
  {
    id: 'user-manager',
    name: '用户管理',
    component: 'UserManager',
    icon: 'UserGroupIcon',
    description: '管理用户和权限',
    category: 'users',
    windowConfig: {
      width: 900,
      height: 600,
      minWidth: 700,
      minHeight: 500,
      resizable: true,
      maximizable: true
    },
    permissions: ['users.read', 'users.write']
  },
  {
    id: 'sync-manager',
    name: '同步备份',
    component: 'SyncManager',
    icon: 'CloudArrowUpIcon',
    description: 'Restic 备份管理和多存储同步',
    category: 'storage',
    windowConfig: {
      width: 1000,
      height: 700,
      minWidth: 800,
      minHeight: 500,
      resizable: true,
      maximizable: true
    },
    permissions: ['storage.read', 'storage.write', 'backup.read', 'backup.write']
  },
  {
    id: 'plugin-store',
    name: '插件商店',
    component: 'PluginStore',
    icon: 'ShoppingBagIcon',
    description: '浏览和安装插件',
    category: 'plugins',
    windowConfig: {
      width: 1000,
      height: 700,
      minWidth: 800,
      minHeight: 500,
      resizable: true,
      maximizable: true
    },
    permissions: ['plugins.read', 'plugins.install']
  },
  {
    id: 'docker-manager',
    name: 'Docker 管理',
    component: 'DockerManager',
    icon: 'CubeIcon',
    description: '管理 Docker 容器和镜像',
    category: 'system',
    windowConfig: {
      width: 1000,
      height: 700,
      minWidth: 800,
      minHeight: 500,
      resizable: true,
      maximizable: true
    },
    permissions: ['docker.read', 'docker.write']
  }
  ]

export interface Theme {
  id: string
  name: string
  type: 'light' | 'dark'
  colors: {
    primary: string
    secondary: string
    background: string
    surface: string
    text: string
    textSecondary: string
    border: string
  }
}

export const availableThemes: Theme[] = [
  {
    id: 'light',
    name: '浅色',
    type: 'light',
    colors: {
      primary: '#3b82f6',
      secondary: '#8b5cf6',
      background: '#f9fafb',
      surface: '#ffffff',
      text: '#1f2937',
      textSecondary: '#6b7280',
      border: '#e5e7eb'
    }
  },
  {
    id: 'dark',
    name: '深色',
    type: 'dark',
    colors: {
      primary: '#60a5fa',
      secondary: '#a78bfa',
      background: '#111827',
      surface: '#1f2937',
      text: '#f9fafb',
      textSecondary: '#9ca3af',
      border: '#374151'
    }
  },
  {
    id: 'blue',
    name: '蓝色',
    type: 'light',
    colors: {
      primary: '#2563eb',
      secondary: '#3b82f6',
      background: '#eff6ff',
      surface: '#ffffff',
      text: '#1e3a8a',
      textSecondary: '#1e40af',
      border: '#bfdbfe'
    }
  },
  {
    id: 'midnight',
    name: '午夜',
    type: 'dark',
    colors: {
      primary: '#818cf8',
      secondary: '#6366f1',
      background: '#0f172a',
      surface: '#1e293b',
      text: '#f1f5f9',
      textSecondary: '#cbd5e1',
      border: '#334155'
    }
  }
]