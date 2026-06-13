// 控制面板系统类型定义

export interface ControlPanelCategory {
  id: string
  name: string
  icon: string
  description: string
  order: number
  settings: ControlPanelSetting[]
  accessLevel: 'user' | 'admin' | 'system'
  readonly?: boolean
}

export interface ControlPanelSetting {
  id: string
  type: SettingType
  category: string
  label: string
  description?: string
  defaultValue: any
  currentValue: any
  validation?: (value: any) => ValidationResult
  options?: SettingOption[]
  dependencies?: string[]
  advanced?: boolean
  restartRequired?: boolean
  secure?: boolean
  readonly?: boolean
  component?: string
}

export type SettingType =
  | 'boolean'
  | 'string'
  | 'number'
  | 'select'
  | 'multiselect'
  | 'slider'
  | 'color'
  | 'textarea'
  | 'password'
  | 'file'
  | 'group'
  | 'custom'
  | 'readonly'

export interface SettingOption {
  label: string
  value: any
  description?: string
  icon?: string
  disabled?: boolean
}

export interface ValidationResult {
  valid: boolean
  error?: string
  warning?: string
}

export interface ControlPanelStore {
  categories: ControlPanelCategory[]
  settings: Record<string, any>
  initialized: boolean
  loading: boolean
  saving: boolean
  lastModified: Date | null
  unsavedChanges: boolean

  // 基础操作
  initialize(): Promise<void>
  updateSetting(key: string, value: any): Promise<void>
  resetToDefaults(): Promise<void>
  resetToCategoryDefaults(categoryId: string): Promise<void>

  // 配置管理
  exportConfig(): string
  importConfig(config: string): Promise<void>
  validateConfig(config: any): ValidationResult

  // 查询功能
  search(query: string): ControlPanelSetting[]
  getSetting(key: string): ControlPanelSetting | undefined
  getSettingsByCategory(categoryId: string): ControlPanelSetting[]
  getAdvancedSettings(): ControlPanelSetting[]

  // 状态管理
  hasUnsavedChanges(): boolean
  discardChanges(): void
  applyChanges(): Promise<void>
}

export interface SettingChangeEvent {
  key: string
  oldValue: any
  newValue: any
  category: string
  timestamp: Date
}

export interface ConfigExport {
  version: string
  timestamp: Date
  categories: ControlPanelCategory[]
  settings: Record<string, any>
  metadata?: {
    hostname?: string
    username?: string
    description?: string
  }
}