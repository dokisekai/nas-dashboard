import { storeToRefs } from 'pinia'
import { useAppStore } from '../stores/app'

export function useDesktop() {
  const appStore = useAppStore()
  const {
    darkMode,
    widgets,
    windows,
    dockItems,
    plugins,
    activeWindowId,
    windowCount,
    focusedWindow,
    minimizedWindows
  } = storeToRefs(appStore)

  const {
    addWidget,
    removeWidget,
    updateWidgetPosition,
    openWindow,
    closeWindow,
    focusWindow,
    minimizeWindow,
    maximizeWindow,
    installPlugin,
    uninstallPlugin,
    toggleDarkMode
  } = appStore

  // 窗口拖拽
  const handleWindowDrag = (windowId: string, position: { x: number; y: number }) => {
    const window = windows.value.find(w => w.id === windowId)
    if (window) {
      window.position = position
      appStore.saveState()
    }
  }

  // 窗口调整大小
  const handleWindowResize = (windowId: string, size: { width: number; height: number }) => {
    const window = windows.value.find(w => w.id === windowId)
    if (window) {
      window.size = size
      appStore.saveState()
    }
  }

  // 小部件拖拽
  const handleWidgetDrag = (widgetId: string, position: { x: number; y: number }) => {
    updateWidgetPosition(widgetId, position)
  }

  // 获取可用的小部件类型
  const getAvailableWidgetTypes = () => {
    const types = [
      { id: 'system-monitor', name: '系统监控', component: 'SystemMonitorWidget' },
      { id: 'weather', name: '天气', component: 'WeatherWidget' },
      { id: 'clock', name: '时钟', component: 'ClockWidget' },
      { id: 'calendar', name: '日历', component: 'CalendarWidget' },
      { id: 'storage-usage', name: '存储使用', component: 'StorageUsageWidget' },
      { id: 'network-stats', name: '网络统计', component: 'NetworkStatsWidget' }
    ]

    // 添加插件提供的小部件
    plugins.value.forEach(plugin => {
      if (plugin.widgets) {
        plugin.widgets.forEach(widgetDef => {
          types.push({
            id: widgetDef.id,
            name: widgetDef.name,
            component: widgetDef.component
          })
        })
      }
    })

    return types
  }

  // 添加小部件到桌面
  const addWidgetToDesktop = (type: string, config: Record<string, any> = {}) => {
    const widgetId = `widget-${Date.now()}`
    const availableTypes = getAvailableWidgetTypes()
    const widgetType = availableTypes.find(t => t.id === type)

    if (!widgetType) {
      console.error(`Widget type ${type} not found`)
      return
    }

    // 计算默认位置（网格布局）
    const columnCount = 4
    const widgetIndex = widgets.value.length
    const row = Math.floor(widgetIndex / columnCount)
    const col = widgetIndex % columnCount

    const position = {
      x: col * 220 + 20,
      y: row * 220 + 20
    }

    addWidget({
      id: widgetId,
      type,
      position,
      size: 'medium',
      config
    })
  }

  // 获取可用的应用
  const getAvailableApps = () => {
    return dockItems.value.map(item => ({
      id: item.id,
      name: item.label,
      icon: item.icon,
      appId: item.appId
    }))
  }

  return {
    // 状态
    darkMode,
    widgets,
    windows,
    dockItems,
    plugins,
    activeWindowId,
    windowCount,
    focusedWindow,
    minimizedWindows,

    // 方法
    addWidget,
    removeWidget,
    updateWidgetPosition,
    openWindow,
    closeWindow,
    focusWindow,
    minimizeWindow,
    maximizeWindow,
    installPlugin,
    uninstallPlugin,
    toggleDarkMode,

    // 拖拽处理
    handleWindowDrag,
    handleWindowResize,
    handleWidgetDrag,

    // 工具方法
    getAvailableWidgetTypes,
    addWidgetToDesktop,
    getAvailableApps
  }
}
