// 通用类型定义
export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

export interface PaginationParams {
  page: number
  pageSize: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}

export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

// 系统信息类型
export interface SystemInfo {
  hostname: string
  os: string
  version: string
  uptime: number
  architecture: string
  cpu: {
    model: string
    cores: number
    frequency: number
  }
  memory: {
    total: number
    available: number
    used: number
    percent: number
  }
}

// 存储信息类型
export interface DiskInfo {
  name: string
  device: string
  mountpoint: string
  size: number
  used: number
  available: number
  percent: number
  fstype: string
}

export interface VolumeInfo {
  id: string
  name: string
  type: string
  size: number
  used: number
  status: 'healthy' | 'warning' | 'error'
  disks: any[]
}

// 服务信息类型
export interface ServiceInfo {
  name: string
  description: string
  status: 'running' | 'stopped' | 'error'
  enabled: boolean
  cpu?: number
  memory?: number
  uptime?: number
}

// 用户信息类型
export interface UserInfo {
  id: number
  username: string
  email?: string
  role: 'admin' | 'user' | 'guest'
  status: 'active' | 'disabled'
  createdAt: string
  lastLogin?: string
}

// 网络信息类型
export interface NetworkInterface {
  name: string
  type: string
  status: 'up' | 'down'
  ipv4?: string
  ipv6?: string
  mac: string
  rxBytes: number
  txBytes: number
  speed: number
}

export interface NetworkStats {
  interfaces: any[]
  totalRx: number
  totalTx: number
  bandwidth: {
    rx: number
    tx: number
  }
}

// 监控数据类型
export interface MonitoringData {
  timestamp: number
  cpu: number
  memory: number
  disk: number
  network: {
    rx: number
    tx: number
  }
}

export interface HistoricalData {
  timeframe: '1h' | '6h' | '24h' | '7d' | '30d'
  interval: number
  data: MonitoringData[]
}
