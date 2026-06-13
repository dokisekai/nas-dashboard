// Disk Management Type Definitions

export interface DiskInfo {
  device: string
  model: string
  serial: string
  size: number
  type: string
  filesystem: string | null
  mountPoint: string | null
  status: 'available' | 'mounted' | 'partitioned' | 'failed'
  partitions: Partition[]
  smartInfo: SMARTInfo | null
  temperature: number
  health: 'good' | 'warning' | 'failed'
}

export interface Partition {
  device: string
  size: number
  start: number
  end: number
  type: string
  filesystem: string | null
  mountPoint: string | null
  flags: string[]
}

export interface PartitionTable {
  type: 'gpt' | 'mbr' | 'none'
  partitions: Partition[]
  freeSpace: number
}

export interface SMARTInfo {
  overallHealth: 'good' | 'warning' | 'failed'
  attributes: SMARTAttribute[]
  lastTest: SMARTTest | null
  temperature: number
  powerOnHours: number
  errorLog: SMARTErrors[]
}

export interface SMARTAttribute {
  id: number
  name: string
  value: number
  worst: number
  threshold: number
  status: 'ok' | 'warning' | 'failed'
}

export interface SMARTTest {
  type: 'short' | 'long' | 'conveyance'
  status: 'running' | 'completed' | 'failed'
  progress: number
  remaining: number
  result: string | null
}

export interface SMARTErrors {
  critical: boolean
  count: number
  details: string[]
}

export interface RAIDConfig {
  id: number
  name: string
  level: '0' | '1' | '5' | '6' | '10'
  devices: RAIDDevice[]
  status: 'active' | 'degraded' | 'failed' | 'rebuilding'
  size: number
  used: number
  uuid: string
  mountPoint: string | null
  createdAt: string
  updatedAt: string
}

export interface RAIDDevice {
  device: string
  size: number
  status: 'active' | 'failed' | 'spare'
  role: 'data' | 'parity' | 'spare'
}

export interface PhysicalVolume {
  id: number
  device: string
  vgName: string
  size: number
  free: number
  uuid: string
  status: string
  createdAt: string
}

export interface VolumeGroup {
  id: number
  name: string
  size: number
  free: number
  pvCount: number
  lvCount: number
  uuid: string
  physicalVolumes: PhysicalVolume[]
  logicalVolumes: LogicalVolume[]
  createdAt: string
  updatedAt: string
}

export interface LogicalVolume {
  id: number
  name: string
  vgName: string
  size: number
  path: string
  uuid: string
  mountPoint: string | null
  status: string
  createdAt: string
  updatedAt: string
}

export interface DiskBenchmark {
  device: string
  readSpeed: number
  writeSpeed: number
  readIOPS: number
  writeIOPS: number
  accessTime: number
  timestamp: string
}

export const RAID_LEVEL_LABELS = {
  '0': 'RAID 0 - 条带',
  '1': 'RAID 1 - 镜像',
  '5': 'RAID 5 - 带奇偶校验',
  '6': 'RAID 6 - 双奇偶校验',
  '10': 'RAID 10 - 镜像+条带'
}

export const RAID_LEVEL_DESCRIPTIONS = {
  '0': '提供最大性能和容量，但无冗余保护',
  '1': '提供数据镜像保护，容量利用率50%',
  '5': '提供较好的性能和冗余平衡，至少需要3块磁盘',
  '6': '类似RAID 5但提供双奇偶校验，至少需要4块磁盘',
  '10': '结合RAID 0和RAID 1的优势，需要偶数块磁盘'
}

export const DISK_HEALTH_COLORS = {
  good: 'green',
  warning: 'orange',
  failed: 'red'
}

export const RAID_STATUS_COLORS = {
  active: 'green',
  degraded: 'orange',
  failed: 'red',
  rebuilding: 'blue'
}