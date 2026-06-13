// Monitor Type Definitions

export interface ProcessInfo {
  pid: number
  name: string
  status: string
  cpuPercent: number
  memoryPercent: number
  memory: number
  threads: number
  username: string
  command: string
  createdTime: number
}

export interface ServiceInfo {
  name: string
  description: string
  status: 'running' | 'stopped' | 'failed' | 'masked'
  enabled: boolean
  loadState: string
  activeState: string
  mainPid: number
  subState: string
}

export interface TemperatureInfo {
  sensors: Sensor[]
}

export interface Sensor {
  name: string
  current: number
  max: number
  critical: number
  unit: string
}

export interface SystemEvent {
  id: number
  type: string
  source: string
  title: string
  message: string
  resolved: boolean
  userId: number | null
  ipAddress: string | null
  createdAt: string
  updatedAt: string
}

export interface SystemLog {
  id: number
  level: 'debug' | 'info' | 'warn' | 'error' | 'fatal'
  component: string
  message: string
  details: string | null
  timestamp: string
  createdAt: string
}

export interface AlertRule {
  id: number
  name: string
  type: 'cpu' | 'memory' | 'disk' | 'temperature' | 'network'
  condition: '>' | '<' | '>=' | '<=' | '==' | '!='
  threshold: number
  duration: number
  severity: 'info' | 'warning' | 'critical'
  enabled: boolean
  actions: string
  lastTriggered: string | null
  createdAt: string
  updatedAt: string
}

export interface MonitorStats {
  cpu: CPUStats
  memory: MemoryStats
  disk: DiskStats[]
  network: NetworkStats[]
  timestamp: number
}

export interface CPUStats {
  usage: number
  cores: number
  load1: number
  load5: number
  load15: number
}

export interface MemoryStats {
  total: number
  used: number
  free: number
  available: number
  percent: number
  swapTotal: number
  swapUsed: number
  swapPercent: number
}

export interface DiskStats {
  device: string
  mountPoint: string
  fileSystem: string
  total: number
  used: number
  free: number
  usedPercent: number
  ioRead: number
  ioWrite: number
}

export interface NetworkStats {
  interface: string
  bytesSent: number
  bytesRecv: number
  packetsSent: number
  packetsRecv: number
  errorsIn: number
  errorsOut: number
  speedIn: number
  speedOut: number
}

export const PROCESS_STATUS_COLORS = {
  running: 'green',
  sleeping: 'blue',
  stopped: 'red',
  zombie: 'orange',
  unknown: 'gray'
}

export const SERVICE_STATUS_COLORS = {
  running: 'green',
  stopped: 'gray',
  failed: 'red',
  masked: 'orange'
}

export const ALERT_SEVERITY_COLORS = {
  info: 'blue',
  warning: 'orange',
  critical: 'red'
}

export const LOG_LEVEL_COLORS = {
  debug: 'gray',
  info: 'blue',
  warn: 'orange',
  error: 'red',
  fatal: 'darkred'
}