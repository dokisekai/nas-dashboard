// Restic Backup & Sync API Client
import api from './client'

// ===== Types =====
export interface BackupRepo {
  id: number
  name: string
  type: 'local' | 's3' | 'sftp' | 'rest' | 'b2' | 'azure' | 'gs' | 'rclone'
  url: string
  password?: string
  env?: Record<string, string>
  status: 'active' | 'uninitialized' | 'error'
  lastError?: string
  lastBackup?: string
  snapshotCount?: number
  repoSize?: number
  createdAt?: string
}

export interface ResticSnapshot {
  id: string
  short_id: string
  time: string
  paths: string[]
  hostname: string
  username?: string
  tags?: string[]
  summary?: {
    files_new: number
    files_changed: number
    files_unmodified: number
    total_files_processed: number
    total_bytes_processed: number
    data_added: number
  }
}

export interface BackupTask {
  id: number
  name: string
  repoId: number
  repo?: BackupRepo
  sourcePath: string
  excludes?: string
  tags?: string
  retention?: Record<string, number>
  autoPrune: boolean
  enabled: boolean
  schedule?: string
  status: 'idle' | 'running' | 'completed' | 'failed'
  lastRun?: string
  lastError?: string
  lastSnapshotId?: string
  lastDuration?: number
  createdAt?: string
}

export interface BackupSyncJob {
  id: number
  name: string
  sourceRepoId: number
  sourceRepo?: BackupRepo
  targetRepoId: number
  targetRepo?: BackupRepo
  enabled: boolean
  status: string
  lastRun?: string
  lastError?: string
}

export interface RepoStats {
  total_size: number
  total_blob_count: number
  snapshots_count: number
}

export interface LSNode {
  name: string
  type: 'file' | 'dir' | 'symlink'
  path: string
  size?: number
  permissions?: string
  mtime?: string
}

// ===== Ping / health =====
export interface PingInfo {
  ok: boolean
  error?: string
  hostname: string
  cacheDir?: string
  resticBin?: string
  version?: string
  runtime?: string
}

export interface BackupSettings {
  defaultHostname: string
  defaultExcludes: string
  defaultTags: string
  autoCheck: boolean
  confirmPurge: boolean
}

export interface WebDAVProfile {
  remoteName: string
  url: string
  vendor: string
  username: string
  password?: string
}

export const resticApi = {
  ping: () => api.get<PingInfo>('/api/storage/backup/ping'),

  // ===== Settings =====
  getSettings: () => api.get<BackupSettings>('/api/storage/backup/settings'),
  updateSettings: (data: BackupSettings) => api.put<BackupSettings>('/api/storage/backup/settings', data),

  // ===== WebDAV profile =====
  getWebDAV: () => api.get<WebDAVProfile>('/api/storage/backup/webdav'),
  updateWebDAV: (data: WebDAVProfile) => api.put<WebDAVProfile>('/api/storage/backup/webdav', data),
  testWebDAV: (data: WebDAVProfile) =>
    api.post<{ ok: boolean; error?: string; output?: string }>('/api/storage/backup/webdav/test', data),

  // ===== Repositories =====
  listRepos: () => api.get<{ repos: BackupRepo[]; total: number }>('/api/storage/backup/repos'),
  createRepo: (data: {
    name: string
    type: BackupRepo['type']
    url: string
    password: string
    env?: Record<string, string>
    init?: boolean
  }) => api.post<BackupRepo>('/api/storage/backup/repos', data),
  updateRepo: (id: number, data: Partial<{
    name: string
    type: BackupRepo['type']
    url: string
    password: string
    env?: Record<string, string>
  }>) => api.put<BackupRepo>(`/api/storage/backup/repos/${id}`, data),
  deleteRepo: (id: number, purge = false) =>
    api.delete(`/api/storage/backup/repos/${id}${purge ? '?purge=true' : ''}`),
  checkRepo: (id: number, full = false) =>
    api.post<{ output: string; status: string }>(`/api/storage/backup/repos/${id}/check${full ? '?full=true' : ''}`),
  refreshRepo: (id: number) =>
    api.post<{ status: string; stats?: RepoStats; repo?: BackupRepo }>(`/api/storage/backup/repos/${id}/refresh`),
  testRepo: (id: number) =>
    api.post<{ ok: boolean; error?: string; snapshotCount?: number; latestSnapshot?: string; status?: string }>(`/api/storage/backup/repos/${id}/test`),
  initRepo: (id: number) => api.post<BackupRepo>(`/api/storage/backup/repos/${id}/init`),
  unlockRepo: (id: number, force = false) =>
    api.post<{ output: string; status: string }>(`/api/storage/backup/repos/${id}/unlock${force ? '?force=true' : ''}`),

  // ===== Snapshots =====
  listSnapshots: (repoId: number) =>
    api.get<{ snapshots: ResticSnapshot[]; total: number }>(`/api/storage/backup/repos/${repoId}/snapshots`),
  snapshotDetail: (repoId: number, sid: string) =>
    api.get<{ snapshot: ResticSnapshot; fileCount: number; dirCount: number; totalSize: number }>(`/api/storage/backup/repos/${repoId}/snapshots/${sid}`),
  listSnapshotFiles: (repoId: number, sid: string) =>
    api.get<{ files: LSNode[]; total: number }>(`/api/storage/backup/repos/${repoId}/snapshots/${sid}/ls`),
  deleteSnapshot: (repoId: number, sid: string, prune = false) =>
    api.delete<{ output: string }>(`/api/storage/backup/repos/${repoId}/snapshots/${sid}${prune ? '?prune=true' : ''}`),
  diffSnapshots: (repoId: number, a: string, b: string) =>
    api.get<{ output: string }>(`/api/storage/backup/repos/${repoId}/diff?a=${a}&b=${b}`),
  findInSnapshots: (repoId: number, pattern: string) =>
    api.get<{ output: string; pattern: string }>(`/api/storage/backup/repos/${repoId}/find?q=${encodeURIComponent(pattern)}`),

  // ===== Restore =====
  restore: (repoId: number, data: {
    snapshotId?: string
    target: string
    include?: string[]
    exclude?: string[]
    host?: string
    paths?: string[]
  }) => api.post<{ message: string; logKey: string }>(`/api/storage/backup/repos/${repoId}/restore`, data),
  restoreLogs: (repoId: number) =>
    api.get<{ lines: string[] }>(`/api/storage/backup/repos/${repoId}/restore/logs`),

  // ===== Backup Tasks =====
  listTasks: () => api.get<{ tasks: BackupTask[]; total: number }>('/api/storage/backup/tasks'),
  createTask: (data: {
    name: string
    repoId: number
    sourcePath: string
    excludes?: string
    tags?: string
    retention?: Record<string, number>
    autoPrune?: boolean
    enabled?: boolean
    schedule?: string
  }) => api.post<BackupTask>('/api/storage/backup/tasks', data),
  updateTask: (id: number, data: any) => api.put<BackupTask>(`/api/storage/backup/tasks/${id}`, data),
  deleteTask: (id: number) => api.delete(`/api/storage/backup/tasks/${id}`),
  runTask: (id: number) => api.post<{ message: string; taskId: number }>(`/api/storage/backup/tasks/${id}/run`),
  taskStatus: (id: number) => api.get<{
    status: string; lastRun?: string; lastError?: string; lastSnapshotId?: string; lastDuration?: number
  }>(`/api/storage/backup/tasks/${id}/status`),
  taskLogs: (id: number) => api.get<{ lines: string[] }>(`/api/storage/backup/tasks/${id}/logs`),

  // ===== Repo-to-repo sync (restic copy) =====
  listSyncJobs: () => api.get<{ jobs: BackupSyncJob[]; total: number }>('/api/storage/backup/sync-jobs'),
  createSyncJob: (data: { name: string; sourceRepoId: number; targetRepoId: number; enabled?: boolean }) =>
    api.post<BackupSyncJob>('/api/storage/backup/sync-jobs', data),
  deleteSyncJob: (id: number) => api.delete(`/api/storage/backup/sync-jobs/${id}`),
  runSyncJob: (id: number) => api.post<{ message: string; jobId: number }>(`/api/storage/backup/sync-jobs/${id}/run`),
  syncJobLogs: (id: number) => api.get<{ lines: string[] }>(`/api/storage/backup/sync-jobs/${id}/logs`),
}

// ExternalContainer: 一个被发现的"外部 restic 容器"。
export interface ExternalContainer {
  id: string
  name: string
  image: string
  state: string
  status: string
  exitCode: number
  created: string
  startedAt?: string
  finishedAt?: string
  command: string
  composeProject?: string
  repo?: string
  sourcePath?: string
  retention?: string
  env: Record<string, string>
  mounts: { source: string; destination: string; mode: string }[]
}

// 保留旧版同步（rsync 文件级）API
export const syncApi = {
  getJobs: () => api.get('/api/storage/sync/jobs'),
  createJob: (data: any) => api.post('/api/storage/sync/jobs', data),
  runJob: (id: number) => api.post(`/api/storage/sync/jobs/${id}/run`),
}

// 旧版备份（tar/pg_dump）API；仅用于向下兼容的旧界面，新代码请使用 resticApi
export const backupApi = resticApi
