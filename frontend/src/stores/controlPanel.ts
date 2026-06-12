import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type {
  ControlPanelStore,
  ControlPanelCategory,
  ControlPanelSetting,
  ValidationResult,
  ConfigExport,
  SettingChangeEvent
} from '../components/ControlPanel/ControlPanelTypes'

export const useControlPanelStore = defineStore('controlPanel', () => {
  // 状态
  const categories = ref<ControlPanelCategory[]>([])
  const settings = ref<Record<string, any>>({})
  const initialized = ref(false)
  const loading = ref(false)
  const saving = ref(false)
  const lastModified = ref<Date | null>(null)
  const originalSettings = ref<Record<string, any>>({})

  // 计算属性
  const unsavedChanges = computed(() => {
    return JSON.stringify(settings.value) !== JSON.stringify(originalSettings.value)
  })

  // 默认分类配置
  const defaultCategories: ControlPanelCategory[] = [
    {
      id: 'general',
      name: '通用设置',
      icon: 'CogIcon',
      description: '系统基础设置和配置',
      order: 1,
      accessLevel: 'user',
      settings: [
        {
          id: 'system.hostname',
          type: 'string',
          category: 'general',
          label: '系统名称',
          description: '设置NAS系统的主机名',
          defaultValue: 'nas-server',
          currentValue: 'nas-server',
          restartRequired: true
        },
        {
          id: 'system.timezone',
          type: 'select',
          category: 'general',
          label: '时区设置',
          description: '选择系统所在时区',
          defaultValue: 'Asia/Shanghai',
          currentValue: 'Asia/Shanghai',
          options: [
            { label: '北京/上海 (GMT+8)', value: 'Asia/Shanghai' },
            { label: '纽约 (GMT-5)', value: 'America/New_York' },
            { label: '伦敦 (GMT+0)', value: 'Europe/London' },
            { label: '东京 (GMT+9)', value: 'Asia/Tokyo' }
          ]
        },
        {
          id: 'system.language',
          type: 'select',
          category: 'general',
          label: '系统语言',
          description: '选择界面显示语言',
          defaultValue: 'zh-CN',
          currentValue: 'zh-CN',
          options: [
            { label: '简体中文', value: 'zh-CN' },
            { label: 'English', value: 'en-US' },
            { label: '日本語', value: 'ja-JP' }
          ]
        },
        {
          id: 'system.autoUpdate',
          type: 'boolean',
          category: 'general',
          label: '自动更新',
          description: '自动下载和安装系统更新',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'system.maintenanceTime',
          type: 'select',
          category: 'general',
          label: '维护时间',
          description: '系统自动维护时间窗口',
          defaultValue: '03:00',
          currentValue: '03:00',
          advanced: true,
          options: [
            { label: '凌晨 3:00', value: '03:00' },
            { label: '凌晨 4:00', value: '04:00' },
            { label: '凌晨 5:00', value: '05:00' }
          ]
        }
      ]
    },
    {
      id: 'network',
      name: '网络设置',
      icon: 'GlobeAltIcon',
      description: '网络连接和接口配置',
      order: 2,
      accessLevel: 'admin',
      settings: [
        {
          id: 'network.hostname',
          type: 'string',
          category: 'network',
          label: '主机名',
          description: '网络中的主机名称',
          defaultValue: 'nas',
          currentValue: 'nas'
        },
        {
          id: 'network.dns',
          type: 'group',
          category: 'network',
          label: 'DNS设置',
          description: '域名服务器配置',
          defaultValue: {},
          currentValue: {}
        },
        {
          id: 'network.dns.primary',
          type: 'string',
          category: 'network',
          label: '首选DNS',
          description: '首选域名服务器',
          defaultValue: '8.8.8.8',
          currentValue: '8.8.8.8',
          dependencies: ['network.dns']
        },
        {
          id: 'network.dns.secondary',
          type: 'string',
          category: 'network',
          label: '备用DNS',
          description: '备用域名服务器',
          defaultValue: '8.8.4.4',
          currentValue: '8.8.4.4',
          dependencies: ['network.dns']
        },
        {
          id: 'network.proxy.enabled',
          type: 'boolean',
          category: 'network',
          label: '启用代理',
          description: '使用代理服务器访问网络',
          defaultValue: false,
          currentValue: false
        },
        {
          id: 'network.proxy.host',
          type: 'string',
          category: 'network',
          label: '代理服务器',
          description: '代理服务器地址',
          defaultValue: '',
          currentValue: '',
          dependencies: ['network.proxy.enabled']
        },
        {
          id: 'network.proxy.port',
          type: 'number',
          category: 'network',
          label: '代理端口',
          description: '代理服务器端口',
          defaultValue: 8080,
          currentValue: 8080,
          dependencies: ['network.proxy.enabled']
        }
      ]
    },
    {
      id: 'security',
      name: '安全设置',
      icon: 'ShieldCheckIcon',
      description: '系统安全和访问控制',
      order: 3,
      accessLevel: 'admin',
      settings: [
        {
          id: 'security.firewall.enabled',
          type: 'boolean',
          category: 'security',
          label: '启用防火墙',
          description: '启用系统防火墙保护',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'security.ssl.enabled',
          type: 'boolean',
          category: 'security',
          label: '启用SSL/TLS',
          description: '使用HTTPS加密连接',
          defaultValue: true,
          currentValue: true,
          restartRequired: true
        },
        {
          id: 'security.ssl.port',
          type: 'number',
          category: 'security',
          label: 'HTTPS端口',
          description: 'HTTPS服务端口',
          defaultValue: 443,
          currentValue: 443,
          dependencies: ['security.ssl.enabled']
        },
        {
          id: 'security.fail2ban.enabled',
          type: 'boolean',
          category: 'security',
          label: '启用防暴力破解',
          description: '自动封禁多次登录失败的IP',
          defaultValue: true,
          currentValue: true,
          advanced: true
        },
        {
          id: 'security.fail2ban.maxAttempts',
          type: 'number',
          category: 'security',
          label: '最大尝试次数',
          description: '允许的最大登录失败次数',
          defaultValue: 5,
          currentValue: 5,
          dependencies: ['security.fail2ban.enabled'],
          advanced: true
        },
        {
          id: 'security.session.timeout',
          type: 'number',
          category: 'security',
          label: '会话超时(分钟)',
          description: '用户会话自动超时时间',
          defaultValue: 30,
          currentValue: 30,
          advanced: true
        }
      ]
    },
    {
      id: 'storage',
      name: '存储设置',
      icon: 'ServerIcon',
      description: '存储管理和配置',
      order: 4,
      accessLevel: 'admin',
      settings: [
        {
          id: 'storage.autoMount',
          type: 'boolean',
          category: 'storage',
          label: '自动挂载',
          description: '系统启动时自动挂载存储设备',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'storage.powerManagement',
          type: 'boolean',
          category: 'storage',
          label: '硬盘省电模式',
          description: '空闲时自动降低硬盘转速',
          defaultValue: false,
          currentValue: false,
          advanced: true
        },
        {
          id: 'storage.smart.enabled',
          type: 'boolean',
          category: 'storage',
          label: 'S.M.A.R.T.监控',
          description: '启用硬盘健康监控',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'storage.smart.interval',
          type: 'select',
          category: 'storage',
          label: '检查间隔',
          description: '硬盘健康检查间隔时间',
          defaultValue: 'daily',
          currentValue: 'daily',
          dependencies: ['storage.smart.enabled'],
          options: [
            { label: '每小时', value: 'hourly' },
            { label: '每天', value: 'daily' },
            { label: '每周', value: 'weekly' }
          ]
        }
      ]
    },
    {
      id: 'notification',
      name: '通知设置',
      icon: 'BellIcon',
      description: '系统通知和告警配置',
      order: 5,
      accessLevel: 'user',
      settings: [
        {
          id: 'notification.email.enabled',
          type: 'boolean',
          category: 'notification',
          label: '邮件通知',
          description: '发送系统通知到邮箱',
          defaultValue: false,
          currentValue: false
        },
        {
          id: 'notification.email.address',
          type: 'string',
          category: 'notification',
          label: '邮箱地址',
          description: '接收通知的邮箱地址',
          defaultValue: '',
          currentValue: '',
          dependencies: ['notification.email.enabled']
        },
        {
          id: 'notification.system.enabled',
          type: 'boolean',
          category: 'notification',
          label: '系统通知',
          description: '显示系统通知',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'notification.sound.enabled',
          type: 'boolean',
          category: 'notification',
          label: '通知声音',
          description: '通知时播放提示音',
          defaultValue: true,
          currentValue: true
        }
      ]
    },
    {
      id: 'appearance',
      name: '外观设置',
      icon: 'PaintBrushIcon',
      description: '界面主题和显示配置',
      order: 6,
      accessLevel: 'user',
      settings: [
        {
          id: 'appearance.theme',
          type: 'select',
          category: 'appearance',
          label: '界面主题',
          description: '选择界面显示主题',
          defaultValue: 'default',
          currentValue: 'default',
          options: [
            { label: '默认主题', value: 'default' },
            { label: '暗色主题', value: 'dark' },
            { label: '简约主题', value: 'minimal' }
          ]
        },
        {
          id: 'appearance.wallpaper.enabled',
          type: 'boolean',
          category: 'appearance',
          label: '启用桌面壁纸',
          description: '显示桌面背景图片',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'appearance.animation.enabled',
          type: 'boolean',
          category: 'appearance',
          label: '界面动画',
          description: '启用界面过渡动画效果',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'appearance.density',
          type: 'select',
          category: 'appearance',
          label: '界面密度',
          description: '界面元素的紧凑程度',
          defaultValue: 'comfortable',
          currentValue: 'comfortable',
          options: [
            { label: '紧凑', value: 'compact' },
            { label: '舒适', value: 'comfortable' },
            { label: '宽松', value: 'spacious' }
          ]
        }
      ]
    }
  ]

  // 初始化控制面板
  const initialize = async () => {
    loading.value = true
    try {
      // 模拟API调用加载配置
      await new Promise(resolve => setTimeout(resolve, 500))

      // 加载分类
      categories.value = defaultCategories

      // 初始化设置值
      const loadedSettings: Record<string, any> = {}
      categories.value.forEach(category => {
        category.settings.forEach(setting => {
          loadedSettings[setting.id] = setting.currentValue
        })
      })

      settings.value = loadedSettings
      originalSettings.value = { ...loadedSettings }
      lastModified.value = new Date()
      initialized.value = true
    } catch (error) {
      console.error('Failed to initialize control panel:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 更新设置
  const updateSetting = async (key: string, value: any) => {
    try {
      // 验证设置
      const setting = getSetting(key)
      if (!setting) {
        throw new Error(`Setting ${key} not found`)
      }

      // 执行验证
      if (setting.validation) {
        const result = setting.validation(value)
        if (!result.valid) {
          throw new Error(result.error || 'Invalid value')
        }
      }

      // 更新值
      const oldValue = settings.value[key]
      settings.value[key] = value

      // 更新当前设置值
      if (setting) {
        setting.currentValue = value
      }

      // 创建变更事件
      const event: SettingChangeEvent = {
        key,
        oldValue,
        newValue: value,
        category: setting.category,
        timestamp: new Date()
      }

      // 模拟API调用保存设置
      await new Promise(resolve => setTimeout(resolve, 200))

      lastModified.value = new Date()

      return event
    } catch (error) {
      console.error(`Failed to update setting ${key}:`, error)
      throw error
    }
  }

  // 重置为默认值
  const resetToDefaults = async () => {
    loading.value = true
    try {
      const newSettings: Record<string, any> = {}
      categories.value.forEach(category => {
        category.settings.forEach(setting => {
          newSettings[setting.id] = setting.defaultValue
          setting.currentValue = setting.defaultValue
        })
      })

      settings.value = newSettings
      originalSettings.value = { ...newSettings }
      lastModified.value = new Date()
    } finally {
      loading.value = false
    }
  }

  // 重置分类默认值
  const resetToCategoryDefaults = async (categoryId: string) => {
    const category = categories.value.find(c => c.id === categoryId)
    if (!category) return

    category.settings.forEach(setting => {
      settings.value[setting.id] = setting.defaultValue
      setting.currentValue = setting.defaultValue
    })

    originalSettings.value = { ...settings.value }
    lastModified.value = new Date()
  }

  // 导出配置
  const exportConfig = (): string => {
    const exportData: ConfigExport = {
      version: '1.0',
      timestamp: new Date(),
      categories: categories.value,
      settings: settings.value,
      metadata: {
        hostname: settings.value['system.hostname'] || 'nas',
        description: 'NAS Dashboard Configuration Export'
      }
    }

    return JSON.stringify(exportData, null, 2)
  }

  // 导入配置
  const importConfig = async (config: string) => {
    loading.value = true
    try {
      const configData: ConfigExport = JSON.parse(config)

      // 验证配置
      const validation = validateConfig(configData)
      if (!validation.valid) {
        throw new Error(validation.error || 'Invalid configuration')
      }

      // 应用配置
      Object.keys(configData.settings).forEach(key => {
        const setting = getSetting(key)
        if (setting) {
          settings.value[key] = configData.settings[key]
          setting.currentValue = configData.settings[key]
        }
      })

      originalSettings.value = { ...settings.value }
      lastModified.value = new Date()
    } catch (error) {
      console.error('Failed to import configuration:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 验证配置
  const validateConfig = (config: any): ValidationResult => {
    if (!config.version || !config.settings) {
      return { valid: false, error: 'Invalid configuration format' }
    }

    // 验证所有设置
    for (const [key, value] of Object.entries(config.settings)) {
      const setting = getSetting(key)
      if (!setting) {
        return { valid: false, error: `Unknown setting: ${key}` }
      }

      if (setting.validation) {
        const result = setting.validation(value)
        if (!result.valid) {
          return result
        }
      }
    }

    return { valid: true }
  }

  // 搜索设置
  const search = (query: string): ControlPanelSetting[] => {
    const results: ControlPanelSetting[] = []
    const lowerQuery = query.toLowerCase()

    categories.value.forEach(category => {
      category.settings.forEach(setting => {
        if (
          setting.label.toLowerCase().includes(lowerQuery) ||
          (setting.description && setting.description.toLowerCase().includes(lowerQuery))
        ) {
          results.push(setting)
        }
      })
    })

    return results
  }

  // 获取设置
  const getSetting = (key: string): ControlPanelSetting | undefined => {
    for (const category of categories.value) {
      const setting = category.settings.find(s => s.id === key)
      if (setting) return setting
    }
    return undefined
  }

  // 获取分类设置
  const getSettingsByCategory = (categoryId: string): ControlPanelSetting[] => {
    const category = categories.value.find(c => c.id === categoryId)
    return category?.settings || []
  }

  // 获取高级设置
  const getAdvancedSettings = (): ControlPanelSetting[] => {
    const results: ControlPanelSetting[] = []
    categories.value.forEach(category => {
      category.settings.forEach(setting => {
        if (setting.advanced) {
          results.push(setting)
        }
      })
    })
    return results
  }

  // 丢弃更改
  const discardChanges = () => {
    settings.value = { ...originalSettings.value }

    // 重置当前值
    Object.keys(settings.value).forEach(key => {
      const setting = getSetting(key)
      if (setting) {
        setting.currentValue = settings.value[key]
      }
    })
  }

  // 应用更改
  const applyChanges = async () => {
    saving.value = true
    try {
      // 模拟API调用
      await new Promise(resolve => setTimeout(resolve, 1000))

      originalSettings.value = { ...settings.value }
      lastModified.value = new Date()
    } finally {
      saving.value = false
    }
  }

  return {
    // 状态
    categories,
    settings,
    initialized,
    loading,
    saving,
    lastModified,
    unsavedChanges,

    // 操作方法
    initialize,
    updateSetting,
    resetToDefaults,
    resetToCategoryDefaults,
    exportConfig,
    importConfig,
    validateConfig,
    search,
    getSetting,
    getSettingsByCategory,
    getAdvancedSettings,
    discardChanges,
    applyChanges
  }
})