import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { APP_REGISTRY, getDefaultSize, getAppTitle } from '../config/apps'

export interface DesktopWindow {
  id: string
  appId: string
  title: string
  position: { x: number; y: number }
  size: { width: number; height: number }
  minimized: boolean
  maximized: boolean
  focused: boolean
  zIndex: number
}

const DOCK_HEIGHT = 80
const CASCADE_OFFSET = 30
const BASE_Z_INDEX = 1000

export const useWindowManagerStore = defineStore('windowManager', () => {
  const windows = ref<DesktopWindow[]>([])
  const nextZIndex = ref(BASE_Z_INDEX)

  const visibleWindows = computed(() => windows.value.filter(w => !w.minimized))
  const focusedWindow = computed(() => windows.value.find(w => w.focused) ?? null)
  const windowCount = computed(() => windows.value.length)

  const isAppOpen = (appId: string) =>
    windows.value.some(w => w.appId === appId)

  const isAppActive = (appId: string) =>
    windows.value.some(w => w.appId === appId && !w.minimized && w.focused)

  const findWindow = (windowId: string) =>
    windows.value.find(w => w.id === windowId)

  const focusWindow = (windowId: string) => {
    const target = findWindow(windowId)
    if (!target) return
    windows.value.forEach(w => { w.focused = false })
    target.focused = true
    target.minimized = false
    target.zIndex = ++nextZIndex.value
  }

  const openApp = (appId: string): DesktopWindow | null => {
    const existing = windows.value.find(w => w.appId === appId)
    if (existing) {
      existing.minimized = false
      focusWindow(existing.id)
      return existing
    }

    const size = getDefaultSize(appId)
    const offset = windows.value.length * CASCADE_OFFSET
    const newWindow: DesktopWindow = {
      id: `window-${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
      appId,
      title: getAppTitle(appId),
      position: { x: 100 + offset, y: 100 + offset },
      size,
      minimized: false,
      maximized: false,
      focused: true,
      zIndex: ++nextZIndex.value,
    }
    windows.value.forEach(w => { w.focused = false })
    windows.value.push(newWindow)
    return newWindow
  }

  const closeWindow = (windowId: string) => {
    windows.value = windows.value.filter(w => w.id !== windowId)
  }

  const minimizeWindow = (windowId: string) => {
    const target = findWindow(windowId)
    if (!target) return
    target.minimized = true
    target.focused = false
  }

  const toggleMaximize = (windowId: string) => {
    const target = findWindow(windowId)
    if (!target) return
    target.maximized = !target.maximized
    if (target.maximized) {
      target.position = { x: 0, y: 0 }
    } else {
      const index = windows.value.findIndex(w => w.id === windowId)
      const offset = index * CASCADE_OFFSET
      target.position = { x: 100 + offset, y: 100 + offset }
    }
  }

  const updatePosition = (windowId: string, position: { x: number; y: number }) => {
    const target = findWindow(windowId)
    if (target) target.position = position
  }

  const updateSize = (windowId: string, size: { width: number; height: number }) => {
    const target = findWindow(windowId)
    if (target) target.size = size
  }

  const tileWindows = () => {
    const targets = visibleWindows.value
    if (targets.length === 0) return

    const availableWidth = window.innerWidth
    const availableHeight = window.innerHeight - DOCK_HEIGHT

    if (targets.length === 1) {
      targets[0].position = {
        x: (availableWidth - targets[0].size.width) / 2,
        y: (availableHeight - targets[0].size.height) / 2,
      }
      return
    }

    if (targets.length === 2) {
      const halfWidth = availableWidth / 2
      targets.forEach((win, index) => {
        win.position = { x: index * halfWidth, y: 0 }
        win.size = { width: halfWidth, height: availableHeight }
      })
      return
    }

    const cols = Math.ceil(Math.sqrt(targets.length))
    const rows = Math.ceil(targets.length / cols)
    const tileWidth = availableWidth / cols
    const tileHeight = availableHeight / rows
    targets.forEach((win, index) => {
      const col = index % cols
      const row = Math.floor(index / cols)
      win.position = { x: col * tileWidth, y: row * tileHeight }
      win.size = { width: tileWidth, height: tileHeight }
    })
  }

  const cascadeWindows = () => {
    const targets = visibleWindows.value
    if (targets.length === 0) return
    const baseX = 50
    const baseY = 50
    targets.forEach((win, index) => {
      win.position = {
        x: baseX + index * CASCADE_OFFSET,
        y: baseY + index * CASCADE_OFFSET,
      }
      win.size = {
        width: Math.min(800, window.innerWidth - baseX - index * CASCADE_OFFSET - 100),
        height: Math.min(600, window.innerHeight - baseY - index * CASCADE_OFFSET - 100),
      }
      win.maximized = false
    })
  }

  const minimizeAll = () => {
    windows.value.forEach(w => { w.minimized = true; w.focused = false })
  }

  const restoreAll = () => {
    windows.value.forEach(w => { w.minimized = false })
  }

  return {
    windows,
    nextZIndex,
    visibleWindows,
    focusedWindow,
    windowCount,
    isAppOpen,
    isAppActive,
    findWindow,
    focusWindow,
    openApp,
    closeWindow,
    minimizeWindow,
    toggleMaximize,
    updatePosition,
    updateSize,
    tileWindows,
    cascadeWindows,
    minimizeAll,
    restoreAll,
  }
})

export { APP_REGISTRY }
