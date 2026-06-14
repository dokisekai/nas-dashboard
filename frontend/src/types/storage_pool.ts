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
  { value: 'create', label: 'Create', description: 'Create on first branch with space' },
  { value: 'mv', label: 'Move', description: 'Move to branch with most space' },
  { value: 'epall', label: 'EP All', description: 'Round-robin across all branches' },
  { value: 'epff', label: 'EP First', description: 'Place on first available branch' },
  { value: 'epmfs', label: 'EP MFS', description: 'Place on branch with most free space' },
  { value: 'eplus', label: 'EP Plus', description: 'Place on branch with most plus space' },
  { value: 'eprand', label: 'EP Random', description: 'Place on random branch' },
  { value: 'lus', label: 'LUS', description: 'Least used space' },
  { value: 'mfs', label: 'MFS', description: 'Most free space' }
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