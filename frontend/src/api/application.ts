// 应用管理API

import api from './client'

// 应用包类型定义
export interface AppPackage {
  id: number
  name: string
  displayName: string
  version: string
  description: string
  author: string
  website: string
  category: string
  license: string
  filePath: string
  fileSize: number
  fileHash: string
  downloadURL: string
  architecture: string
  minOSVersion: string
  maxOSVersion: string
  minRAM: number
  minDiskSpace: number
  dependencies: string[]
  resources: string
  permissions: string
  repositoryId: number
  downloadCount: number
  installCount: number
  rating: number
  createdAt: string
  updatedAt: string
}

// 应用实例类型定义
export interface AppInstance {
  id: number
  name: string
  displayName: string
  packageName: string
  version: string
  description: string
  category: string
  author: string
  website: string
  status: 'running' | 'stopped' | 'error' | 'installing' | 'uninstalling'
  containerId: string
  pid: number
  exitCode: number
  lastExitTime: string
  config: string
  envVars: string
  ports: string
  volumes: string
  resources: string
  permissions: string
  installPath: string
  dataPath: string
  configPath: string
  backupPaths: string
  createdAt: string
  updatedAt: string
}

// 应用仓库类型定义
export interface AppRepository {
  id: number
  name: string
  url: string
  type: 'official' | 'community' | 'custom'
  enabled: boolean
  priority: number
  description: string
  autoUpdate: boolean
  createdAt: string
  updatedAt: string
}

// 安装进度类型定义
export interface AppInstallProgress {
  step: string
  message: string
  percent: number
  status: 'running' | 'success' | 'error'
}

// 安装请求类型定义
export interface AppInstallRequest {
  packageName: string
  version?: string
  config?: Record<string, any>
  autoStart?: boolean
}

// 应用操作请求类型定义
export interface AppActionRequest {
  action: 'start' | 'stop' | 'restart' | 'reload' | 'enable' | 'disable'
  params?: Record<string, any>
}

// API响应类型定义
export interface AppListResponse {
  instances: AppInstance[]
  packages: AppPackage[]
}

export interface AppDetailResponse {
  instance: AppInstance
}

export interface AppStatusResponse {
  status: string
}

export interface AppProgressResponse {
  progress: AppInstallProgress
}

export interface RepositoryListResponse {
  repositories: AppRepository[]
}

// 应用管理API客户端
export const applicationApi = {
  // 应用包管理
  async listPackages(): Promise<AppPackage[]> {
    const response = await api.get<{ packages: AppPackage[] }>('/packages')
    return response.packages
  },

  async uploadPackage(file: File): Promise<AppPackage> {
    const formData = new FormData()
    formData.append('package', file)

    const response = await api.post<{ package: AppPackage }>('/packages/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.package
  },

  async getPackage(name: string): Promise<AppPackage> {
    const response = await api.get<{ package: AppPackage }>(`/packages/${name}`)
    return response.package
  },

  async deletePackage(name: string): Promise<void> {
    await api.delete(`/packages/${name}`)
  },

  // 应用实例管理
  async installApp(request: AppInstallRequest): Promise<AppInstance> {
    const response = await api.post<{ instance: AppInstance }>('/apps/install', request)
    return response.instance
  },

  async listApps(): Promise<AppInstance[]> {
    const response = await api.get<{ instances: AppInstance[] }>('/apps')
    return response.instances
  },

  async getApp(id: number): Promise<AppInstance> {
    const response = await api.get<{ instance: AppInstance }>(`/apps/${id}`)
    return response.instance
  },

  async getAppStatus(id: number): Promise<string> {
    const response = await api.get<{ status: string }>(`/apps/${id}/status`)
    return response.status
  },

  async startApp(id: number): Promise<void> {
    await api.post(`/apps/${id}/start`)
  },

  async stopApp(id: number): Promise<void> {
    await api.post(`/apps/${id}/stop`)
  },

  async restartApp(id: number): Promise<void> {
    await api.post(`/apps/${id}/restart`)
  },

  async updateAppConfig(id: number, config: Record<string, any>): Promise<void> {
    await api.put(`/apps/${id}/config`, config)
  },

  async uninstallApp(id: number): Promise<void> {
    await api.delete(`/apps/${id}`)
  },

  // 应用仓库管理
  async listRepositories(): Promise<AppRepository[]> {
    const response = await api.get<{ repositories: AppRepository[] }>('/repositories')
    return response.repositories
  },

  async addRepository(repository: Omit<AppRepository, 'id' | 'createdAt' | 'updatedAt'>): Promise<AppRepository> {
    const response = await api.post<{ repository: AppRepository }>('/repositories', repository)
    return response.repository
  },

  async updateRepository(repository: AppRepository): Promise<void> {
    await api.put(`/repositories/${repository.id}`, repository)
  },

  async deleteRepository(id: number): Promise<void> {
    await api.delete(`/repositories/${id}`)
  },

  async syncRepository(id: number): Promise<void> {
    await api.post(`/repositories/${id}/sync`)
  },

  // 应用更新
  async checkUpdates(id: number): Promise<any[]> {
    const response = await api.get<{ updates: any[] }>(`/apps/${id}/updates`)
    return response.updates
  },

  async updateApp(id: number): Promise<void> {
    await api.post(`/apps/${id}/update`)
  }
}

// 获取安装进度的SSE连接
export function getInstallProgress(id: number, onProgress: (progress: AppInstallProgress) => void, onComplete?: () => void) {
  const token = localStorage.getItem('token')
  const baseUrl = import.meta.env.VITE_API_BASE_URL || 'http://192.168.50.10:8888'
  const eventSource = new EventSource(`${baseUrl}/api/apps/${id}/progress?token=${token}`)

  eventSource.onmessage = (event) => {
    try {
      const data = JSON.parse(event.data)
      if (data.status === 'complete') {
        eventSource.close()
        onComplete?.()
      } else {
        onProgress(data)
      }
    } catch (error) {
      console.error('解析安装进度失败:', error)
    }
  }

  eventSource.onerror = (error) => {
    console.error('SSE连接错误:', error)
    eventSource.close()
  }

  return eventSource
}