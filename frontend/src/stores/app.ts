import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Widget, Window, DockItem, Plugin } from '../types/desktop'

export const useAppStore = defineStore('app', () => {
  // 状态
  const darkMode = ref(false)
  const widgets = ref<Widget[]>([])
  const windows = ref<Window[]>([])
  const dockItems = ref<DockItem[]>([])
  const plugins = ref<Plugin[]>([])
  const activeWindowId = ref<string | null>(null)

  // 计算属性
  const windowCount = computed(() => windows.value.length)
  const focusedWindow = computed(() =>
    windows.value.find(w => w.id === activeWindowId.value)
  )
  const minimizedWindows = computed(() =>
    windows.value.filter(w => w.minimized)
  )

  // 方法
  const addWidget = (widget: Widget) => {
    widgets.value.push(widget)
    saveState()
  }

  const removeWidget = (widgetId: string) => {
    widgets.value = widgets.value.filter(w => w.id !== widgetId)
    saveState()
  }

  const updateWidgetPosition = (widgetId: string, position: { x: number; y: number }) => {
    const widget = widgets.value.find(w => w.id === widgetId)
    if (widget) {
      widget.position = position
      saveState()
    }
  }

  const openWindow = (appId: string, config?: Partial<Window>) => {
    const existingWindow = windows.value.find(w => w.appId === appId)
    if (existingWindow) {
      focusWindow(existingWindow.id)
      return existingWindow
    }

    const newWindow: Window = {
      id: `window-${Date.now()}`,
      appId,
      title: getAppTitle(appId),
      position: {
        x: 100 + windows.value.length * 30,
        y: 100 + windows.value.length * 30
      },
      size: { width: 800, height: 600 },
      minimized: false,
      maximized: false,
      focused: true,
      zIndex: getNextZIndex(),
      ...config
    }

    windows.value.push(newWindow)
    focusWindow(newWindow.id)
    saveState()
    return newWindow
  }

  const closeWindow = (windowId: string) => {
    windows.value = windows.value.filter(w => w.id !== windowId)
    if (activeWindowId.value === windowId) {
      activeWindowId.value = null
    }
    saveState()
  }

  const focusWindow = (windowId: string) => {
    windows.value.forEach(w => {
      w.focused = w.id === windowId
      if (w.focused) {
        w.zIndex = getNextZIndex()
        w.minimized = false
      }
    })
    activeWindowId.value = windowId
    saveState()
  }

  const minimizeWindow = (windowId: string) => {
    const window = windows.value.find(w => w.id === windowId)
    if (window) {
      window.minimized = true
      window.focused = false
      if (activeWindowId.value === windowId) {
        activeWindowId.value = null
      }
      saveState()
    }
  }

  const maximizeWindow = (windowId: string) => {
    const window = windows.value.find(w => w.id === windowId)
    if (window) {
      window.maximized = !window.maximized
      saveState()
    }
  }

  const installPlugin = async (plugin: Plugin) => {
    plugins.value.push(plugin)
    // 注册插件的小部件和应用
    if (plugin.widgets) {
      plugin.widgets.forEach(widgetDef => {
        // 注册小部件
      })
    }
    if (plugin.apps) {
      plugin.apps.forEach(appDef => {
        // 注册应用到Dock
        dockItems.value.push({
          id: appDef.id,
          label: appDef.name,
          icon: appDef.icon,
          appId: appDef.id
        })
      })
    }
    saveState()
  }

  const uninstallPlugin = (pluginId: string) => {
    const plugin = plugins.value.find(p => p.id === pluginId)
    if (plugin) {
      // 移除插件相关的小部件
      if (plugin.widgets) {
        widgets.value = widgets.value.filter(w =>
          !plugin.widgets!.some(def => def.id === w.type)
        )
      }

      // 移除插件相关的Dock项目
      if (plugin.apps) {
        dockItems.value = dockItems.value.filter(item =>
          !plugin.apps!.some(app => app.id === item.appId)
        )
      }

      plugins.value = plugins.value.filter(p => p.id !== pluginId)
      saveState()
    }
  }

  const toggleDarkMode = () => {
    darkMode.value = !darkMode.value
    saveState()
  }

  // 辅助函数
  const getNextZIndex = () => {
    const maxZ = Math.max(...windows.value.map(w => w.zIndex || 0), 1000)
    return maxZ + 1
  }

  const getAppTitle = (appId: string) => {
    const titles: Record<string, string> = {
      'app-center': '应用中心',
      'storage-manager': '存储管理',
      'system-monitor': '系统监控',
      'user-manager': '用户管理',
      'settings': '系统设置',
      'plugin-store': '插件商店'
    }
    return titles[appId] || appId
  }

  const saveState = () => {
    const state = {
      darkMode: darkMode.value,
      widgets: widgets.value,
      dockItems: dockItems.value,
      plugins: plugins.value
    }
    localStorage.setItem('nas-desktop-state', JSON.stringify(state))
  }

  const loadState = () => {
    try {
      const saved = localStorage.getItem('nas-desktop-state')
      if (saved) {
        const state = JSON.parse(saved)
        darkMode.value = state.darkMode || false
        widgets.value = state.widgets || []
        dockItems.value = state.dockItems || []
        plugins.value = state.plugins || []
      }
    } catch (error) {
      console.error('Failed to load state:', error)
    }
  }

  // 初始化时加载状态
  loadState()

  return {
    // 状态
    darkMode,
    widgets,
    windows,
    dockItems,
    plugins,
    activeWindowId,
    // 计算属性
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
    toggleDarkMode
  }
})