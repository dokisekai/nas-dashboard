/**
 * 应用配置管理
 * 统一管理前端应用的所有配置项
 */

// API配置
export const API_CONFIG = {
  // API基础URL
  baseURL: import.meta.env.VITE_API_URL || '',

  // 超时设置
  timeout: 10000,

  // 重试配置
  retry: {
    maxRetries: 3,
    retryDelay: 1000,
    retryableStatus: [0, 500, 502, 503, 504]
  },

  // 缓存配置
  cache: {
    enabled: true,
    ttl: 5 * 60 * 1000 // 5分钟
  }
} as const

// WebSocket配置
export const WS_CONFIG = {
  // WebSocket URL
  url: import.meta.env.VITE_WS_URL || '',

  // 重连配置
  reconnect: {
    enabled: true,
    maxAttempts: 5,
    delay: 3000,
    backoffMultiplier: 1.5
  },

  // 心跳配置
  heartbeat: {
    enabled: true,
    interval: 30000,
    timeout: 60000
  }
} as const

// 界面配置
export const UI_CONFIG = {
  // 主题配置
  theme: {
    default: 'dark',
    storageKey: 'nas_dashboard_theme'
  },

  // 语言配置
  language: {
    default: 'zh-CN',
    storageKey: 'nas_dashboard_language',
    available: ['zh-CN', 'en-US']
  },

  // 布局配置
  layout: {
    sidebarWidth: 250,
    sidebarCollapsedWidth: 64,
    headerHeight: 60,
    footerHeight: 40
  },

  // 动画配置
  animation: {
    enabled: true,
    duration: 300,
    easing: 'ease-in-out'
  }
} as const

// 监控配置
export const MONITOR_CONFIG = {
  // 更新间隔
  updateInterval: 2000,

  // 历史数据保留时间
  historyDuration: 5 * 60 * 1000, // 5分钟

  // 图表配置
  chart: {
    maxDataPoints: 100,
    animationDuration: 300,
    smooth: true
  },

  // 告警阈值
  thresholds: {
    cpu: {
      warning: 70,
      critical: 90
    },
    memory: {
      warning: 80,
      critical: 95
    },
    disk: {
      warning: 85,
      critical: 95
    },
    temperature: {
      warning: 70,
      critical: 85
    }
  }
} as const

// 文件管理配置
export const FILE_MANAGER_CONFIG = {
  // 分页配置
  pagination: {
    pageSize: 50,
    pageSizeOptions: [20, 50, 100, 200]
  },

  // 上传配置
  upload: {
    maxFileSize: 1024 * 1024 * 1024, // 1GB
    chunkSize: 5 * 1024 * 1024, // 5MB
    concurrentUploads: 3
  },

  // 预览配置
  preview: {
    imageExtensions: ['.jpg', '.jpeg', '.png', '.gif', '.webp', '.svg'],
    videoExtensions: ['.mp4', '.webm', '.ogg'],
    textExtensions: ['.txt', '.md', '.json', '.xml', '.log'],
    maxPreviewSize: 10 * 1024 * 1024 // 10MB
  }
} as const

// Docker管理配置
export const DOCKER_CONFIG = {
  // 容器配置
  container: {
    defaultTimeout: 30000,
    logLines: 100,
    statsInterval: 2000
  },

  // 镜像配置
  image: {
    pullTimeout: 300000, // 5分钟
    pageSize: 50
  },

  // 网络配置
  network: {
    defaultDriver: 'bridge',
    pageSize: 50
  },

  // 卷配置
  volume: {
    pageSize: 50
  }
} as const

// 通知配置
export const NOTIFICATION_CONFIG = {
  // 显示配置
  display: {
    duration: 5000,
    maxVisible: 5,
    position: 'top-right' as const
  },

  // 声音配置
  sound: {
    enabled: false,
    volume: 0.5
  },

  // 类型配置
  types: {
    success: {
      icon: '✓',
      color: '#10b981'
    },
    error: {
      icon: '✕',
      color: '#ef4444'
    },
    warning: {
      icon: '⚠',
      color: '#f59e0b'
    },
    info: {
      icon: 'ℹ',
      color: '#3b82f6'
    }
  }
} as const

// 存储管理配置
export const STORAGE_CONFIG = {
  // 存储池配置
  pool: {
    refreshInterval: 5000,
    defaultRedundancy: '2-copy' as const
  },

  // 备份配置
  backup: {
    defaultCompression: true,
    defaultEncryption: true,
    maxBackupAge: 90 // days
  },

  // RAID配置
  raid: {
    supportedLevels: ['0', '1', '5', '6', '10'] as const,
    defaultLevel: '1' as const
  }
} as const

// 网络管理配置
export const NETWORK_CONFIG = {
  // 接口配置
  interface: {
    refreshInterval: 3000,
    showVirtual: false,
    showLoopback: false
  },

  // WiFi配置
  wifi: {
    scanInterval: 10000,
    connectionTimeout: 30000
  },

  // 防火墙配置
  firewall: {
    defaultPolicy: 'deny' as const,
    logging: true
  }
} as const

// 安全配置
export const SECURITY_CONFIG = {
  // 会话配置
  session: {
    timeout: 30 * 60 * 1000, // 30分钟
    refreshThreshold: 5 * 60 * 1000 // 5分钟
  },

  // 密码策略
  password: {
    minLength: 8,
    requireUppercase: false,
    requireLowercase: true,
    requireNumbers: true,
    requireSpecialChars: false
  },

  // 2FA配置
  twoFactor: {
    enabled: false,
    issuer: 'NAS Dashboard'
  }
} as const

// 开发者配置
export const DEV_CONFIG = {
  // 调试模式
  debug: import.meta.env.DEV || false,

  // 日志级别
  logLevel: import.meta.env.VITE_LOG_LEVEL || 'info',

  // Mock数据
  mockEnabled: import.meta.env.VITE_MOCK_ENABLED === 'true',

  // 性能监控
  performance: {
    enabled: import.meta.env.DEV,
    sampleRate: 1.0
  }
} as const

// 特性开关
export const FEATURE_FLAGS = {
  // 插件系统
  pluginSystem: true,

  // 高级监控
  advancedMonitoring: true,

  // 文件预览
  filePreview: true,

  // 容器终端
  containerTerminal: true,

  // 工作流系统
  workflowSystem: false, // 实验性功能

  // 主题定制
  themeCustomization: true,

  // 多语言支持
  multiLanguage: false
} as const

// 统一配置导出
export const APP_CONFIG = {
  api: API_CONFIG,
  websocket: WS_CONFIG,
  ui: UI_CONFIG,
  monitor: MONITOR_CONFIG,
  fileManager: FILE_MANAGER_CONFIG,
  docker: DOCKER_CONFIG,
  notification: NOTIFICATION_CONFIG,
  storage: STORAGE_CONFIG,
  network: NETWORK_CONFIG,
  security: SECURITY_CONFIG,
  dev: DEV_CONFIG,
  features: FEATURE_FLAGS
} as const

// 配置获取函数
export function getConfig<K extends keyof typeof APP_CONFIG>(
  module: K
): typeof APP_CONFIG[K] {
  return APP_CONFIG[module]
}

// 运行时配置更新
export function updateConfig<K extends keyof typeof APP_CONFIG>(
  module: K,
  updates: Partial<typeof APP_CONFIG[K]>
): void {
  Object.assign(APP_CONFIG[module], updates)
}

// 配置验证
export function validateConfig(): boolean {
  try {
    // 验证必需的配置项
    if (!API_CONFIG.baseURL && !window.location.origin) {
      console.warn('API base URL not configured')
    }

    // 验证数值范围
    if (API_CONFIG.timeout < 1000) {
      console.warn('API timeout too low, may cause issues')
    }

    return true
  } catch (error) {
    console.error('Configuration validation failed:', error)
    return false
  }
}

// 初始化配置
export function initConfig(): void {
  // 从localStorage加载用户配置
  const userTheme = localStorage.getItem(UI_CONFIG.theme.storageKey)
  const userLanguage = localStorage.getItem(UI_CONFIG.language.storageKey)

  // 应用用户配置
  if (userTheme) {
    updateConfig('ui', { theme: { ...UI_CONFIG.theme, default: userTheme } })
  }

  if (userLanguage) {
    updateConfig('ui', { language: { ...UI_CONFIG.language, default: userLanguage } })
  }

  // 验证配置
  validateConfig()

  console.log('Application configuration initialized')
}