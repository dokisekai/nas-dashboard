// Storage Pool Type Definitions

export interface StoragePool {
  id: number
  name: string
  type: 'mergerfs' | 'btrfs' | 'zfs' | 'lvm'
  status: 'active' | 'inactive' | 'creating' | 'deleting' | 'error' | 'degraded'
  mountPoint: string
  totalSize: number
  usedSize: number
  freeSize: number
  description: string
  config: MergerFSConfig | null
  poolDisks: PoolDisk[]
  snapshots: PoolSnapshot[]
  createdAt: string
  updatedAt: string
}

export interface MergerFSConfig {
  branches: BranchConfig[]
  category: 'create' | 'mv' | 'epall' | 'epff' | 'epmfs' | 'eplus' | 'eprand' | 'lus' | 'lus' | 'mfs'
  minfreespace: string
  direct_io: boolean
  async_read: boolean
  use_ino: boolean
  hard_remove: boolean
  auto_unshare: boolean
  follow_symlinks: boolean
  link_exas: boolean
}

export interface BranchConfig {
  path: string
  mode: 'ro' | 'rw'
  priority: number
}

export interface PoolDisk {
  id?: number
  poolId?: number
  device: string
  size: number
  status?: 'active' | 'failed' | 'removed' | 'ro' | 'rw'
  priority: number
  branchPath: string
  used?: number
  free?: number
  createdAt?: string
  updatedAt?: string
}

export interface PoolSnapshot {
  id: number
  poolId: number
  name: string
  description: string
  size: number
  status: 'creating' | 'completed' | 'deleting' | 'error'
  createdAt: string
  completedAt: string | null
}

export interface StoragePoolCreateRequest {
  name: string
  type: 'mergerfs' | 'btrfs' | 'zfs' | 'lvm'
  mountPoint: string
  description?: string
  config?: MergerFSConfig
  disks: PoolDisk[]
}

export interface StoragePoolUpdateRequest {
  description?: string
  config?: MergerFSConfig
}

export interface PoolStatusResponse {
  name: string
  status: string
  branches: BranchInfo[]
  totalSize: number
  usedSize: number
  freeSize: number
}

export interface BranchInfo {
  path: string
  mode: 'ro' | 'rw'
  priority: number
  size: number
  used: number
  free: number
}

export interface MergerFSCategory {
  value: string
  label: string
  description: string
}

export const MERGERFS_CATEGORIES: MergerFSCategory[] = [
  { value: 'create', label: 'Create (ff)', description: '在第一个有空间的分支创建新文件' },
  { value: 'mv', label: 'Move (mv)', description: '移动到剩余空间最多的分支' },
  { value: 'epall', label: 'EP All (epall)', description: '在所有分支之间轮询分配新文件' },
  { value: 'epff', label: 'EP First (epff)', description: '按顺序在第一个可用分支放置新文件' },
  { value: 'epmfs', label: 'EP MFS (epmfs)', description: '在剩余空间最多的分支放置新文件' },
  { value: 'eplus', label: 'EP Plus (eplus)', description: '在绝对剩余空间最多的分支放置新文件' },
  { value: 'eprand', label: 'EP Random (eprand)', description: '随机选择分支放置新文件' },
  { value: 'lus', label: 'LUS (lus)', description: '选择使用空间最少的分支' },
  { value: 'mfs', label: 'MFS (mfs)', description: '选择剩余空间最多的分支' },
  { value: 'lfs', label: 'LFS (lfs)', description: '选择剩余空间最少的分支' },
  { value: 'rand', label: 'Random (rand)', description: '完全随机选择分支' }
]

export const STORAGE_POOL_STATUS_COLORS = {
  active: 'green',
  inactive: 'gray',
  creating: 'blue',
  deleting: 'orange',
  error: 'red',
  degraded: 'yellow'
}

export const STORAGE_POOL_TYPE_LABELS = {
  mergerfs: 'MergerFS',
  btrfs: 'Btrfs',
  zfs: 'ZFS',
  lvm: 'LVM'
}