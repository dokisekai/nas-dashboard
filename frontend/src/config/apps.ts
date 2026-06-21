import type { Component } from 'vue'
import {
  ServerIcon,
  ChartBarIcon,
  CloudArrowUpIcon,
  UserGroupIcon,
  ShoppingBagIcon,
  CubeIcon,
  ShieldCheckIcon,
  PhotoIcon,
  Cog6ToothIcon,
  CircleStackIcon,
} from '@heroicons/vue/24/outline'

export interface AppDefinition {
  id: string
  title: string
  icon: Component
  component?: () => Promise<unknown>
  width: number
  height: number
  external?: boolean
}

const DEFAULT_SIZE = { width: 800, height: 600 }

export const APP_REGISTRY: Record<string, AppDefinition> = {
  'storage-manager': {
    id: 'storage-manager',
    title: '存储管理',
    icon: ServerIcon,
    component: () => import('../apps/StorageManager.vue'),
    width: 900,
    height: 600,
  },
  'system-monitor': {
    id: 'system-monitor',
    title: '系统监控',
    icon: ChartBarIcon,
    component: () => import('../apps/SystemMonitor.vue'),
    width: 1000,
    height: 700,
  },
  'app-center': {
    id: 'app-center',
    title: '应用管理中心',
    icon: ShoppingBagIcon,
    component: () => import('../apps/AppCenter.vue'),
    width: 1400,
    height: 900,
  },
  'user-manager': {
    id: 'user-manager',
    title: '用户管理',
    icon: UserGroupIcon,
    component: () => import('../apps/UserManager.vue'),
    width: 800,
    height: 500,
  },
  'backup-manager': {
    id: 'backup-manager',
    title: '备份管理',
    icon: CircleStackIcon,
    component: () => import('../apps/BackupManager.vue'),
    width: 1100,
    height: 720,
  },
  'sync-manager': {
    id: 'sync-manager',
    title: '多存储同步',
    icon: CloudArrowUpIcon,
    component: () => import('../apps/SyncManager.vue'),
    width: 1100,
    height: 720,
  },
  'control-panel': {
    id: 'control-panel',
    title: '控制面板',
    icon: Cog6ToothIcon,
    component: () => import('../apps/ControlPanel.vue'),
    width: 1200,
    height: 800,
  },
  'docker-manager': {
    id: 'docker-manager',
    title: 'Docker 管理',
    icon: CubeIcon,
    component: () => import('../apps/DockerManager.vue'),
    width: 1000,
    height: 700,
  },
  'sso-manager': {
    id: 'sso-manager',
    title: 'SSO 统一身份认证',
    icon: ShieldCheckIcon,
    component: () => import('../apps/SSOManager.vue'),
    width: 1200,
    height: 800,
  },
  'immich-photo': {
    id: 'immich-photo',
    title: 'Immich 照片管理',
    icon: PhotoIcon,
    width: 1200,
    height: 800,
    external: true,
  },
}

export const DOCK_APP_IDS = [
  'storage-manager',
  'system-monitor',
  'app-center',
  'user-manager',
  'sync-manager',
  'control-panel',
  'docker-manager',
  'immich-photo',
  'sso-manager',
] as const

export const QUICK_ACCESS_APP_IDS = [
  'storage-manager',
  'system-monitor',
  'app-center',
] as const

export function getApp(id: string): AppDefinition | undefined {
  return APP_REGISTRY[id]
}

export function getDockApps(): AppDefinition[] {
  return DOCK_APP_IDS.map(id => APP_REGISTRY[id]).filter(Boolean)
}

export function getQuickAccessApps(): AppDefinition[] {
  return QUICK_ACCESS_APP_IDS.map(id => APP_REGISTRY[id]).filter(Boolean)
}

export function getDefaultSize(id: string) {
  const app = APP_REGISTRY[id]
  return app ? { width: app.width, height: app.height } : DEFAULT_SIZE
}

export function getAppTitle(id: string): string {
  return APP_REGISTRY[id]?.title ?? id
}
