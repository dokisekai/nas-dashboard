import api from './client'

// ==================== 认证 API ====================
export const authApi = {
  login: (username: string, password: string) =>
    api.post('/api/auth/login', { username, password }),
  refresh: () => api.post('/api/auth/refresh'),
}

// ==================== 监控 API ====================
export const monitorApi = {
  getCPU: () => api.get('/api/monitor/cpu'),
  getMemory: () => api.get('/api/monitor/memory'),
  getDisk: () => api.get('/api/monitor/disk'),
  getNetwork: () => api.get('/api/monitor/network'),
  getTemperature: () => api.get('/api/monitor/temperature'),
  getProcesses: () => api.get('/api/monitor/processes'),
  getServices: () => api.get('/api/monitor/services'),
}

// ==================== 网络管理 API ====================
export { networkApi, wifiApi, dnsApi, networkUtils } from './network'
export { interfaceConfigApi, pppoeConfigApi, proxyConfigApi, networkConfigUtils } from './interface_config'

// ==================== 防火墙管理 API ====================
export { firewallApi, firewallUtils } from './firewall'

// ==================== 存储管理 API ====================
export const storageApi = {
  getDisks: () => api.get('/api/storage/disks'),
  getDiskUsage: (path: string) => api.get('/api/storage/usage', { params: { path } }),
  mount: (device: string, mountPoint: string) =>
    api.post('/api/storage/mount', { device, mountPoint }),
  umount: (mountPoint: string) =>
    api.post('/api/storage/umount', { mountPoint }),
  formatDisk: (device: string, fsType: string) =>
    api.post('/api/storage/format', { device, fsType }),
  getSMBShares: () => api.get('/api/storage/smb'),
  createSMBShare: (data: { name: string, path: string, description?: string, readOnly?: boolean, guest?: boolean, isTimeMachine?: boolean }) =>
    api.post('/api/storage/smb', data),
  updateSMBShare: (name: string, data: { path: string, description?: string, readOnly?: boolean, guest?: boolean, isTimeMachine?: boolean }) =>
    api.put('/api/storage/smb/' + name, data),
  deleteSMBShare: (name: string) =>
    api.delete('/api/storage/smb/' + name),
}

// ==================== 服务管理 API ====================
export const serviceApi = {
  getServices: () => api.get('/api/services'),
  startService: (name: string) => api.post(`/api/services/${name}/start`),
  stopService: (name: string) => api.post(`/api/services/${name}/stop`),
  restartService: (name: string) => api.post(`/api/services/${name}/restart`),
  enableService: (name: string) => api.post(`/api/services/${name}/enable`),
  disableService: (name: string) => api.post(`/api/services/${name}/disable`),
  getContainers: () => api.get('/api/docker/containers'),
  startContainer: (id: string) => api.post(`/api/docker/containers/${id}/start`),
  stopContainer: (id: string) => api.post(`/api/docker/containers/${id}/stop`),
  restartContainer: (id: string) => api.post(`/api/docker/containers/${id}/restart`),
  removeContainer: (id: string) => api.delete(`/api/docker/containers/${id}`),
  getContainerLogs: (id: string) => api.get(`/api/docker/containers/${id}/logs`),
  getContainerStats: (id: string) => api.get(`/api/docker/containers/${id}/stats`),
  execInContainer: (id: string, command: string[]) => api.post(`/api/docker/containers/${id}/exec`, { command }),
  getImages: () => api.get('/api/docker/images'),
  removeImage: (id: string) => api.delete(`/api/docker/images/${id}`),
  pullImage: (image: string) => api.post('/api/docker/images/pull', { image }),
  getNetworks: () => api.get('/api/docker/networks'),
  getVolumes: () => api.get('/api/docker/volumes'),
}

// ==================== 用户管理 API ====================
export const userApi = {
  getUsers: () => api.get('/api/users'),
  getUser: (username: string) => api.get(`/api/users/${username}`),
  createUser: (data: { username: string; password: string; comment?: string; group?: string; shell?: string }) =>
    api.post('/api/users', data),
  updateUser: (username: string, data: { password?: string; group?: string; shell?: string; comment?: string }) =>
    api.put(`/api/users/${username}`, data),
  deleteUser: (username: string) => api.delete(`/api/users/${username}`),
  getSSHKeys: (user?: string) => api.get('/api/users/ssh-keys', { params: { user } }),
  addKey: (data: { name: string; content: string; user: string }) =>
    api.post('/api/users/ssh-keys', data),
  deleteKey: (id: string, user?: string) => api.delete(`/api/users/ssh-keys/${id}`, { params: { user } }),
  getCurrentUser: () => api.get('/api/users/me'),
  changePassword: (oldPassword: string, newPassword: string) =>
    api.post('/api/users/me/password', { oldPassword, newPassword }),
  getUserQuota: (username: string) => api.get(`/api/users/${username}/quota`),
  setUserQuota: (username: string, quota: { space?: string; files?: string }) =>
    api.put(`/api/users/${username}/quota`, quota),
}

// ==================== 系统组 API ====================
export const groupApi = {
  getGroups: () => api.get('/api/groups'),
  createGroup: (data: { name: string; description?: string; gid?: number }) =>
    api.post('/api/groups', data),
  updateGroup: (name: string, data: { description?: string; members?: string[]; gid?: number }) =>
    api.put(`/api/groups/${name}`, data),
  deleteGroup: (name: string) => api.delete(`/api/groups/${name}`),
  getGroupMembers: (name: string) => api.get(`/api/groups/${name}/members`),
  addGroupMembers: (name: string, members: string[]) =>
    api.post(`/api/groups/${name}/members`, { members }),
  removeGroupMember: (name: string, username: string) =>
    api.delete(`/api/groups/${name}/members/${username}`),
}

// ==================== 系统信息 API ====================
export const systemApi = {
  getInfo: () => api.get('/api/system/info'),
  getHardwareDetails: () => api.get('/api/system/hardware'),
  getPowerUsage: () => api.get('/api/system/power'),
  getUptime: () => api.get('/api/system/uptime'),
  getUPSStatus: (ups?: string) => api.get('/api/system/ups/status', { params: { ups } }),

  // 系统操作
  restart: () => api.post('/api/system/operations/restart'),
  shutdown: () => api.post('/api/system/operations/shutdown'),
  cancel: () => api.post('/api/system/operations/cancel'),
  rebootImmediate: () => api.post('/api/system/operations/reboot-immediate'),
  poweroffImmediate: () => api.post('/api/system/operations/poweroff-immediate'),
  scheduleShutdown: (delayMinutes: number) =>
    api.post('/api/system/operations/schedule-shutdown', { delayMinutes }),
  scheduleRestart: (delayMinutes: number) =>
    api.post('/api/system/operations/schedule-restart', { delayMinutes }),
  getStatus: () => api.get('/api/system/operations/status')
}

// ==================== 文件管理 API ====================
export const fileApi = {
  listFiles: (path: string) => api.post('/api/files/list', { path }),
  getFileInfo: (path: string) => api.get('/api/files/info', { params: { path } }),
  downloadFile: (path: string) => api.get('/api/files/download', { params: { path }, responseType: 'blob' }),
  uploadFile: (path: string, file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('path', path)
    return api.post('/api/files/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },
  createDirectory: (path: string, permissions?: string) =>
    api.post('/api/files/directory', { path, permissions }),
  moveFile: (oldPath: string, newPath: string) =>
    api.post('/api/files/move', { oldPath, newPath }),
  deleteFile: (path: string, force?: boolean) =>
    api.post('/api/files/delete', { path, force }),
}

// ==================== 备份恢复 API ====================
export const legacyBackupApi = {
  getBackups: () => api.get('/api/backups'),
  getBackup: (id: number) => api.get(`/api/backups/${id}`),
  createBackup: (data: any) => api.post('/api/backups', data),
  deleteBackup: (id: number) => api.delete(`/api/backups/${id}`),
  restoreBackup: (data: { backupId: number }) => api.post('/api/backups/restore', data),
  downloadBackup: (id: number) => `/api/backups/${id}/download`,
}

// ==================== 系统配置 API ====================
export const configApi = {
  getConfigs: (category?: string) => api.get('/api/configs', { params: { category } }),
  getPublicConfigs: () => api.get('/api/configs/public'),
  getConfig: (key: string) => api.get(`/api/configs/${key}`),
  setConfig: (data: any) => api.post('/api/configs', data),
  deleteConfig: (key: string) => api.delete(`/api/configs/${key}`),
  bulkSetConfig: (configs: any[]) => api.post('/api/configs/bulk', { configs }),
}

// ==================== 应用管理 API ====================
export { applicationApi } from './application'

// ==================== 应用管理 API ====================
export { appsApi } from './apps'
export type { AppContainer, ComposeProject, ServiceCatalogEntry, ServiceStatus, ContainerStats, PortBinding, AppMount } from './apps'

// ==================== 同步与备份 API ====================
export { syncApi, backupApi, resticApi } from './sync_backup'
export type { BackupRepo, ResticSnapshot, BackupTask, BackupSyncJob, RepoStats, LSNode, PingInfo, BackupSettings, WebDAVProfile } from './sync_backup'
