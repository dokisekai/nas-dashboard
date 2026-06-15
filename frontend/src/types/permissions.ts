// 权限类型定义

export type PermissionLevel = 'none' | 'read' | 'write' | 'admin' | 'owner'

export type UserRole = 'guest' | 'user' | 'admin' | 'superadmin'

export interface Permission {
  id: string
  name: string
  description: string
  category: PermissionCategory
  level: PermissionLevel
  resource?: string
}

export type PermissionCategory =
  | 'system'        // 系统管理
  | 'storage'       // 存储管理
  | 'file'          // 文件操作
  | 'network'       // 网络管理
  | 'user'          // 用户管理
  | 'backup'        // 备份恢复
  | 'app'           // 应用管理
  | 'log'           // 日志查看
  | 'monitor'       // 系统监控

export interface UserPermissions {
  system: {
    settings: PermissionLevel
    reboot: PermissionLevel
    shutdown: PermissionLevel
  }
  storage: {
    view: PermissionLevel
    manage: PermissionLevel
    format: PermissionLevel
  }
  file: {
    read: PermissionLevel
    write: PermissionLevel
    delete: PermissionLevel
    share: PermissionLevel
  }
  network: {
    view: PermissionLevel
    configure: PermissionLevel
    firewall: PermissionLevel
  }
  user: {
    view: PermissionLevel
    manage: PermissionLevel
    permissions: PermissionLevel
  }
  backup: {
    create: PermissionLevel
    restore: PermissionLevel
    schedule: PermissionLevel
  }
  app: {
    install: PermissionLevel
    uninstall: PermissionLevel
    configure: PermissionLevel
  }
  log: {
    view: PermissionLevel
    export: PermissionLevel
  }
  monitor: {
    view: PermissionLevel
    configure: PermissionLevel
  }
}

export interface User {
  id: string
  username: string
  fullName?: string
  email?: string
  uid: number
  primaryGroup: string
  groups: string[]
  role: UserRole
  permissions: UserPermissions
  status: 'active' | 'disabled' | 'locked'
  homeDirectory: string
  shell: string
  lastLogin?: Date
  createdAt: Date
  avatar?: string
  description?: string
  storageQuota?: {
    enabled: boolean
    size: number
    used: number
  }
  passwordPolicy?: {
    requireChange: boolean
    expireDays: number
    complexity: boolean
  }
}

export interface Group {
  id: string
  name: string
  gid: number
  description?: string
  permissions: UserPermissions
  members: string[]
  createdAt: Date
  type: 'system' | 'user' | 'builtin'
  storageQuota?: {
    enabled: boolean
    size: number
  }
}

// 预定义权限模板
export const PermissionTemplates = {
  guest: {
    system: { settings: 'none', reboot: 'none', shutdown: 'none' },
    storage: { view: 'read', manage: 'none', format: 'none' },
    file: { read: 'read', write: 'none', delete: 'none', share: 'none' },
    network: { view: 'read', configure: 'none', firewall: 'none' },
    user: { view: 'none', manage: 'none', permissions: 'none' },
    backup: { create: 'none', restore: 'none', schedule: 'none' },
    app: { install: 'none', uninstall: 'none', configure: 'read' },
    log: { view: 'none', export: 'none' },
    monitor: { view: 'read', configure: 'none' }
  },
  user: {
    system: { settings: 'read', reboot: 'none', shutdown: 'none' },
    storage: { view: 'read', manage: 'none', format: 'none' },
    file: { read: 'write', write: 'write', delete: 'write', share: 'read' },
    network: { view: 'read', configure: 'none', firewall: 'none' },
    user: { view: 'read', manage: 'none', permissions: 'none' },
    backup: { create: 'write', restore: 'read', schedule: 'read' },
    app: { install: 'none', uninstall: 'none', configure: 'write' },
    log: { view: 'read', export: 'none' },
    monitor: { view: 'read', configure: 'none' }
  },
  admin: {
    system: { settings: 'admin', reboot: 'write', shutdown: 'write' },
    storage: { view: 'admin', manage: 'admin', format: 'admin' },
    file: { read: 'admin', write: 'admin', delete: 'admin', share: 'admin' },
    network: { view: 'admin', configure: 'admin', firewall: 'admin' },
    user: { view: 'admin', manage: 'admin', permissions: 'admin' },
    backup: { create: 'admin', restore: 'admin', schedule: 'admin' },
    app: { install: 'admin', uninstall: 'admin', configure: 'admin' },
    log: { view: 'admin', export: 'admin' },
    monitor: { view: 'admin', configure: 'admin' }
  },
  superadmin: {
    system: { settings: 'owner', reboot: 'owner', shutdown: 'owner' },
    storage: { view: 'owner', manage: 'owner', format: 'owner' },
    file: { read: 'owner', write: 'owner', delete: 'owner', share: 'owner' },
    network: { view: 'owner', configure: 'owner', firewall: 'owner' },
    user: { view: 'owner', manage: 'owner', permissions: 'owner' },
    backup: { create: 'owner', restore: 'owner', schedule: 'owner' },
    app: { install: 'owner', uninstall: 'owner', configure: 'owner' },
    log: { view: 'owner', export: 'owner' },
    monitor: { view: 'owner', configure: 'owner' }
  }
} as any

// 权限级别定义
export const PermissionLevels = {
  none: { value: 0, label: '无权限', color: 'gray' },
  read: { value: 1, label: '只读', color: 'blue' },
  write: { value: 2, label: '读写', color: 'green' },
  admin: { value: 3, label: '管理', color: 'orange' },
  owner: { value: 4, label: '所有者', color: 'red' }
} as any

// 用户角色定义
export const UserRoles = {
  guest: {
    label: '访客',
    description: '仅限基本访问权限',
    permissions: PermissionTemplates.guest,
    icon: '👤',
    color: 'gray'
  },
  user: {
    label: '普通用户',
    description: '标准用户权限',
    permissions: PermissionTemplates.user,
    icon: '👤',
    color: 'blue'
  },
  admin: {
    label: '管理员',
    description: '可管理大部分系统功能',
    permissions: PermissionTemplates.admin,
    icon: '👨‍💼',
    color: 'orange'
  },
  superadmin: {
    label: '超级管理员',
    description: '拥有所有系统权限',
    permissions: PermissionTemplates.superadmin,
    icon: '👑',
    color: 'red'
  }
} as any

// 系统内置组
export const SystemGroups = [
  {
    id: 'administrators',
    name: 'administrators',
    gid: 999,
    description: '系统管理员组，拥有完整管理权限',
    type: 'builtin',
    permissions: PermissionTemplates.admin
  },
  {
    id: 'users',
    name: 'users',
    gid: 100,
    description: '普通用户组，标准用户权限',
    type: 'builtin',
    permissions: PermissionTemplates.user
  },
  {
    id: 'guests',
    name: 'guests',
    gid: 998,
    description: '访客组，仅限基本访问',
    type: 'builtin',
    permissions: PermissionTemplates.guest
  }
] as any

// 权限继承规则
export function getEffectivePermissions(user: User, groups: Group[]): UserPermissions {
  let permissions = { ...user.permissions }

  // 从用户组获取最高权限
  for (const groupName of user.groups) {
    const group = groups.find(g => g.name === groupName)
    if (group) {
      permissions = mergePermissions(permissions, group.permissions)
    }
  }

  // 根据角色应用模板权限
  const rolePermissions = PermissionTemplates[user.role]
  if (rolePermissions) {
    permissions = mergePermissions(permissions, rolePermissions)
  }

  return permissions
}

// 合并权限（取最高级别）
function mergePermissions(
  base: UserPermissions,
  override: Partial<UserPermissions>
): UserPermissions {
  const result = { ...base }

  for (const key in override) {
    const category = key as keyof UserPermissions
    const baseCategory = result[category] as any
    const overrideCategory = override[category] as any

    if (baseCategory && overrideCategory) {
      result[category] = { ...baseCategory }
      for (const perm in overrideCategory) {
        const permKey = perm as keyof typeof baseCategory
        const baseLevel = baseCategory[permKey] as PermissionLevel
        const overrideLevel = overrideCategory[permKey] as PermissionLevel

        // 取权限级别更高的
        if (getPermissionLevelValue(overrideLevel) > getPermissionLevelValue(baseLevel)) {
          (result[category] as any)[permKey] = overrideLevel
        }
      }
    }
  }

  return result
}

function getPermissionLevelValue(level: PermissionLevel): number {
  return PermissionLevels[level].value
}