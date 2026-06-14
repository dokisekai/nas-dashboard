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

  // 默认分类配置 - 合并系统设置功能
  const defaultCategories: ControlPanelCategory[] = [
    {
      id: 'locale',
      name: '语言和时区',
      icon: 'GlobeAltIcon',
      description: '系统语言、时区和地区设置',
      order: 1,
      accessLevel: 'user',
      settings: [
        {
          id: 'locale.language',
          type: 'select',
          category: 'locale',
          label: '系统语言',
          description: '选择界面显示语言',
          defaultValue: 'zh-CN',
          currentValue: 'zh-CN',
          options: [
            { label: '简体中文', value: 'zh-CN' },
            { label: 'English', value: 'en-US' },
            { label: '日本語', value: 'ja-JP' },
            { label: 'Deutsch', value: 'de-DE' },
            { label: 'Français', value: 'fr-FR' },
            { label: 'Español', value: 'es-ES' }
          ]
        },
        {
          id: 'locale.timezone',
          type: 'select',
          category: 'locale',
          label: '时区设置',
          description: '选择系统所在时区',
          defaultValue: 'Asia/Shanghai',
          currentValue: 'Asia/Shanghai',
          options: [
            { label: '北京/上海 (GMT+8)', value: 'Asia/Shanghai' },
            { label: '纽约 (GMT-5)', value: 'America/New_York' },
            { label: '伦敦 (GMT+0)', value: 'Europe/London' },
            { label: '东京 (GMT+9)', value: 'Asia/Tokyo' },
            { label: '巴黎 (GMT+1)', value: 'Europe/Paris' },
            { label: '悉尼 (GMT+10)', value: 'Australia/Sydney' },
            { label: '迪拜 (GMT+4)', value: 'Asia/Dubai' },
            { label: '洛杉矶 (GMT-8)', value: 'America/Los_Angeles' }
          ]
        },
        {
          id: 'locale.dateFormat',
          type: 'select',
          category: 'locale',
          label: '日期格式',
          description: '选择日期显示格式',
          defaultValue: 'YYYY-MM-DD',
          currentValue: 'YYYY-MM-DD',
          options: [
            { label: '2024-06-13', value: 'YYYY-MM-DD' },
            { label: '06/13/2024', value: 'MM/DD/YYYY' },
            { label: '13.06.2024', value: 'DD.MM.YYYY' },
            { label: '13/06/2024', value: 'DD/MM/YYYY' }
          ]
        },
        {
          id: 'locale.timeFormat',
          type: 'select',
          category: 'locale',
          label: '时间格式',
          description: '选择时间显示格式',
          defaultValue: '24h',
          currentValue: '24h',
          options: [
            { label: '24小时制 (14:30)', value: '24h' },
            { label: '12小时制 (02:30 PM)', value: '12h' }
          ]
        },
        {
          id: 'locale.firstDayOfWeek',
          type: 'select',
          category: 'locale',
          label: '每周开始日',
          description: '设置日历中每周的起始日',
          defaultValue: '1',
          currentValue: '1',
          options: [
            { label: '星期一', value: '1' },
            { label: '星期日', value: '0' },
            { label: '星期六', value: '6' }
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
      accessLevel: 'user',
      settings: [
        {
          id: 'network.hostname',
          type: 'string',
          category: 'network',
          label: '主机名',
          description: '设置系统主机名',
          defaultValue: 'nas-server',
          currentValue: 'nas-server'
        },
        {
          id: 'network.dhcp',
          type: 'boolean',
          category: 'network',
          label: '启用DHCP',
          description: '自动获取IP地址',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'network.settings',
          type: 'custom',
          category: 'network',
          label: '网络配置',
          description: '直接配置网络接口、IP地址、DNS和代理设置',
          defaultValue: {},
          currentValue: {},
          component: 'NetworkSettingsPanel',
          advanced: false
        }
      ]
    },
    {
      id: 'firewall',
      name: '防火墙设置',
      icon: 'FireIcon',
      description: '网络防火墙和安全规则管理',
      order: 3,
      accessLevel: 'admin',
      settings: [
        {
          id: 'firewall.enabled',
          type: 'boolean',
          category: 'firewall',
          label: '启用防火墙',
          description: '启用系统防火墙保护网络连接',
          defaultValue: true,
          currentValue: true,
          restartRequired: true
        },
        {
          id: 'firewall.defaultPolicy',
          type: 'select',
          category: 'firewall',
          label: '默认策略',
          description: '设置防火墙的默认处理策略',
          defaultValue: 'drop',
          currentValue: 'drop',
          options: [
            { label: '拒绝(推荐)', value: 'drop', disabled: false },
            { label: '接受', value: 'accept', disabled: false }
          ],
          dependencies: ['firewall.enabled']
        },
        {
          id: 'firewall.allowedPorts',
          type: 'multiselect',
          category: 'firewall',
          label: '允许的端口',
          description: '选择允许通过防火墙的网络端口',
          defaultValue: ['22', '80', '443'],
          currentValue: ['22', '80', '443'],
          options: [
            { label: 'SSH (22)', value: '22', disabled: false },
            { label: 'HTTP (80)', value: '80', disabled: false },
            { label: 'HTTPS (443)', value: '443', disabled: false },
            { label: 'FTP (21)', value: '21', disabled: false },
            { label: 'SMB (445)', value: '445', disabled: false },
            { label: 'NFS (2049)', value: '2049', disabled: false },
            { label: 'DNS (53)', value: '53', disabled: false },
            { label: 'DHCP (67-68)', value: '67-68', disabled: false }
          ],
          dependencies: ['firewall.enabled']
        },
        {
          id: 'firewall.logging',
          type: 'boolean',
          category: 'firewall',
          label: '启用防火墙日志',
          description: '记录所有防火墙规则匹配的数据包',
          defaultValue: false,
          currentValue: false,
          dependencies: ['firewall.enabled'],
          advanced: true
        },
        {
          id: 'firewall.icmp',
          type: 'boolean',
          category: 'firewall',
          label: '允许ICMP(PING)',
          description: '允许网络ICMP数据包用于网络诊断',
          defaultValue: true,
          currentValue: true,
          dependencies: ['firewall.enabled'],
          advanced: true
        }
      ]
    },
    {
      id: 'notification',
      name: '通知设置',
      icon: 'BellIcon',
      description: '系统通知和告警配置',
      order: 4,
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
      id: 'system-info',
      name: '系统信息',
      icon: 'ServerIcon',
      description: '查看系统硬件、运行信息和网络状态',
      order: 5,
      accessLevel: 'user',
      settings: [
        {
          id: 'system.info.panel',
          type: 'custom',
          category: 'system-info',
          label: '系统信息面板',
          description: '显示系统硬件、运行状态和网络接口信息',
          defaultValue: {},
          currentValue: {},
          component: 'SystemInfoPanel'
        }
      ]
    },
    {
      id: 'user-management',
      name: '用户管理',
      icon: 'UserIcon',
      description: '用户账户和用户组管理',
      order: 6,
      accessLevel: 'admin',
      settings: [
        {
          id: 'users.manager',
          type: 'custom',
          category: 'user-management',
          label: '用户管理',
          description: '管理系统用户账户',
          defaultValue: {},
          currentValue: {},
          component: 'UserManager'
        },
        {
          id: 'users.groups',
          type: 'custom',
          category: 'user-management',
          label: '用户组管理',
          description: '管理系统用户组',
          defaultValue: {},
          currentValue: {},
          component: 'GroupManager'
        }
      ]
    },
    {
      id: 'permission-management',
      name: '权限管理',
      icon: 'LockClosedIcon',
      description: '文件和文件夹权限设置',
      order: 7,
      accessLevel: 'admin',
      settings: [
        {
          id: 'permissions.acl',
          type: 'custom',
          category: 'permission-management',
          label: 'ACL编辑器',
          description: '访问控制列表高级编辑器',
          defaultValue: {},
          currentValue: {},
          component: 'ACLEditor'
        }
      ]
    },
    {
      id: 'backup',
      name: '备份管理',
      icon: 'CircleStackIcon',
      description: '系统备份和恢复',
      order: 8,
      accessLevel: 'admin',
      settings: [
        {
          id: 'backup.enabled',
          type: 'boolean',
          category: 'backup',
          label: '启用自动备份',
          description: '启用系统配置和数据自动备份',
          defaultValue: true,
          currentValue: true
        },
        {
          id: 'backup.schedule',
          type: 'select',
          category: 'backup',
          label: '备份频率',
          description: '自动备份的执行频率',
          defaultValue: 'daily',
          currentValue: 'daily',
          options: [
            { label: '每小时', value: 'hourly' },
            { label: '每天', value: 'daily' },
            { label: '每周', value: 'weekly' },
            { label: '每月', value: 'monthly' }
          ],
          dependencies: ['backup.enabled']
        },
        {
          id: 'backup.time',
          type: 'string',
          category: 'backup',
          label: '备份时间',
          description: '执行自动备份的具体时间',
          defaultValue: '02:00',
          currentValue: '02:00',
          dependencies: ['backup.enabled']
        },
        {
          id: 'backup.location',
          type: 'string',
          category: 'backup',
          label: '备份位置',
          description: '备份文件存储位置',
          defaultValue: '/var/backups/nas',
          currentValue: '/var/backups/nas',
          dependencies: ['backup.enabled']
        },
        {
          id: 'backup.retention',
          type: 'number',
          category: 'backup',
          label: '保留天数',
          description: '备份文件保留天数',
          defaultValue: 7,
          currentValue: 7,
          dependencies: ['backup.enabled'],
          advanced: true
        },
        {
          id: 'backup.compression',
          type: 'boolean',
          category: 'backup',
          label: '压缩备份',
          description: '压缩备份文件以节省空间',
          defaultValue: true,
          currentValue: true,
          dependencies: ['backup.enabled'],
          advanced: true
        },
        {
          id: 'backup.encrypt',
          type: 'boolean',
          category: 'backup',
          label: '加密备份',
          description: '使用加密保护备份文件',
          defaultValue: false,
          currentValue: false,
          dependencies: ['backup.enabled'],
          advanced: true
        }
      ]
    },
    {
      id: 'smb',
      name: '文件共享',
      icon: 'FolderIcon',
      description: 'SMB文件共享和访问权限设置',
      order: 6,
      accessLevel: 'user',
      settings: [] // SMB管理使用自定义组件，不需要传统设置
    }
  ]

  // 初始化控制面板
  const initialize = async () => {
    loading.value = true
    try {
      // 加载分类
      categories.value = defaultCategories

      // 从真实API获取系统信息
      await loadSystemInfo()

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

  // 从API加载系统信息
  const loadSystemInfo = async () => {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        console.log('No token found, skipping system info load')
        return
      }

      console.log('Loading system info with token')

      // 导入axios实例
      const { default: axios } = await import('axios')
      const api = axios.create({
        baseURL: '',
        timeout: 10000,
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      })

      // 获取网络信息
      try {
        const networkData = await api.get('/api/monitor/network') as any
        console.log('Network data loaded:', networkData)
        updateNetworkMetrics(networkData)
      } catch (err) {
        console.error('Network API error:', err)
      }

      // 获取CPU信息
      try {
        const cpuData = await api.get('/api/monitor/cpu') as any
        updateSystemInfo('system.cpu', `${cpuData.model || 'Unknown'} (${cpuData.cores || 0} cores)`)
      } catch (err) {
        console.error('CPU API error:', err)
      }

      // 获取内存信息
      try {
        const memData = await api.get('/api/monitor/memory') as any
        const totalGB = (memData.total / (1024 ** 3)).toFixed(0)
        const usedGB = (memData.used / (1024 ** 3)).toFixed(1)
        updateSystemInfo('system.memory', `${totalGB}GB total, ${usedGB}GB used`)
      } catch (err) {
        console.error('Memory API error:', err)
      }

    } catch (error) {
      console.error('Failed to load system info:', error)
    }
  }

  // 更新网络统计信息
  const updateNetworkMetrics = (networkData: any) => {
    if (!networkData || !networkData.interfaces) {
      console.log('Invalid network data')
      return
    }

    const interfaces = networkData.interfaces
    let totalUploadSpeed = 0
    let totalDownloadSpeed = 0

    interfaces.forEach((iface: any) => {
      if (iface.up && !isVirtualInterface(iface.name)) {
        totalUploadSpeed += iface.sentSpeed || 0
        totalDownloadSpeed += iface.recvSpeed || 0
      }
    })

    const uploadSpeedStr = formatSpeed(totalUploadSpeed)
    const downloadSpeedStr = formatSpeed(totalDownloadSpeed)

    console.log('Network metrics updated:', uploadSpeedStr, downloadSpeedStr)

    // 更新网络统计设置
    const networkMetricSetting = getSetting('network.metrics')
    if (networkMetricSetting) {
      networkMetricSetting.currentValue = `↑ ${uploadSpeedStr} ↓ ${downloadSpeedStr}`
      settings.value['network.metrics'] = `↑ ${uploadSpeedStr} ↓ ${downloadSpeedStr}`
    }
  }

  // 更新系统信息
  const updateSystemInfo = (settingId: string, value: string) => {
    const setting = getSetting(settingId)
    if (setting) {
      setting.currentValue = value
      settings.value[settingId] = value
      console.log('System info updated:', settingId, value)
    }
  }

  // 判断是否为虚拟接口
  const isVirtualInterface = (name: string): boolean => {
    return name.startsWith('veth') ||
           name.startsWith('virbr') ||
           name.startsWith('docker') ||
           name.startsWith('br-') ||
           name.startsWith('lo')
  }

  // 格式化速度
  const formatSpeed = (bytesPerSecond: number): string => {
    if (!bytesPerSecond || bytesPerSecond === 0) return '0 B/s'

    const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
    let value = bytesPerSecond
    let unitIndex = 0

    while (value >= 1024 && unitIndex < units.length - 1) {
      value /= 1024
      unitIndex++
    }

    return `${value.toFixed(1)} ${units[unitIndex]}`
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