// 桌面系统类型定义

export interface Widget {
  id: string
  type: string
  position: { x: number; y: number }
  size: 'small' | 'medium' | 'large'
  config: Record<string, any>
}

export interface Window {
  id: string
  appId: string
  title: string
  position: { x: number; y: number }
  size: { width: number; height: number }
  minimized: boolean
  maximized: boolean
  focused: boolean
  zIndex?: number
}

export interface DockItem {
  id: string
  label: string
  icon: string
  appId?: string
  badge?: number | null
}

export interface ContextMenuItem {
  id: string
  label: string
  action?: () => void
  items?: ContextMenuItem[]
}

export interface Plugin {
  id: string
  name: string
  version: string
  description: string
  icon: string
  author: string
  permissions: string[]
  widgets?: WidgetDefinition[]
  apps?: AppDefinition[]
}

export interface WidgetDefinition {
  id: string
  name: string
  type: string
  component: string
  defaultSize: 'small' | 'medium' | 'large'
  defaultConfig: Record<string, any>
}

export interface AppDefinition {
  id: string
  name: string
  component: string
  icon: string
  windowConfig?: {
    width: number
    height: number
    resizable?: boolean
    minimizable?: boolean
    maximizable?: boolean
  }
}