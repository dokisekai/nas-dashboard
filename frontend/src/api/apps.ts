// App Manager (Docker services) API Client
import api from './client'

export interface PortBinding {
  hostIp: string
  hostPort: number
  containerPort: number
  protocol: string
  url?: string
}

export interface AppMount {
  source: string
  destination: string
  mode: string
}

export interface ContainerStats {
  name: string
  cpuPercent: number
  memoryUsage: number
  memoryLimit: number
  memoryPercent: number
  networkRx: number
  networkTx: number
  readTime: string
}

export interface AppContainer {
  id: string
  name: string
  image: string
  state: string
  status: string
  health?: string
  exitCode: number
  created: string
  startedAt?: string
  finishedAt?: string
  restartCount: number
  ports: PortBinding[]
  ip: string
  command: string
  composeProject: string
  composeService: string
  composeDir: string
  composeFile: string
  labels: Record<string, string>
  mounts: AppMount[]
  envCount: number
  stats?: ContainerStats
}

export interface ComposeProject {
  name: string
  configFile: string
  workingDir: string
  containers: string[]
  runningCount: number
  totalCount: number
  category: string
  description: string
  iconHint: string
}

export type ServiceStatus = 'running' | 'partial' | 'stopped' | 'missing' | 'unknown'

export interface ServiceCatalogEntry {
  id: string
  name: string
  category: string
  description: string
  iconHint: string
  url?: string
  containers: AppContainer[]
  status: ServiceStatus
  runningCount: number
  totalCount: number
  composeDir: string
  composeFile: string
  notes: string
}

export const appsApi = {
  catalog: () => api.get<{ catalog: ServiceCatalogEntry[]; total: number }>('/api/apps/catalog'),
  projects: () => api.get<{ projects: ComposeProject[]; total: number; available: boolean }>('/api/apps/projects'),
  listContainers: (stats = false) =>
    api.get<{ containers: AppContainer[]; total: number; available: boolean }>(
      `/api/apps/containers${stats ? '?stats=true' : ''}`,
    ),
  containerLogs: (name: string, tail = 500) =>
    api.get<{ logs: string; name: string }>(`/api/apps/containers/${encodeURIComponent(name)}/logs?tail=${tail}`),
  containerStats: (name: string) => api.get<ContainerStats>(`/api/apps/containers/${encodeURIComponent(name)}/stats`),
  containerAction: (name: string, action: 'start' | 'stop' | 'restart' | 'remove', opts?: { force?: boolean; volumes?: boolean }) => {
    const q = opts && (opts.force || opts.volumes)
      ? `?${opts.force ? 'force=true&' : ''}${opts.volumes ? 'volumes=true' : ''}`
      : ''
    return api.post<{ message: string; action: string }>(
      `/api/apps/containers/${encodeURIComponent(name)}/${action}${q}`,
    )
  },
}
