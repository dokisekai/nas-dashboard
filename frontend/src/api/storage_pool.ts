// Storage Pool API Client
import axios from 'axios'
import type {
  StoragePool,
  StoragePoolCreateRequest,
  StoragePoolUpdateRequest,
  PoolStatusResponse,
  BranchInfo
} from '@/types/storage_pool'

const API_BASE = '/api/storage/pools'

export const storagePoolAPI = {
  // Get all storage pools
  async getPools(): Promise<StoragePool[]> {
    const response = await axios.get(`${API_BASE}`)
    return response.data.pools || []
  },

  // Get storage pool by name
  async getPool(name: string): Promise<StoragePool> {
    const response = await axios.get(`${API_BASE}/${name}`)
    return response.data
  },

  // Create new storage pool
  async createPool(data: StoragePoolCreateRequest): Promise<StoragePool> {
    const response = await axios.post(`${API_BASE}`, data)
    return response.data
  },

  // Update storage pool
  async updatePool(name: string, data: StoragePoolUpdateRequest): Promise<StoragePool> {
    const response = await axios.put(`${API_BASE}/${name}`, data)
    return response.data
  },

  // Delete storage pool
  async deletePool(name: string): Promise<void> {
    await axios.delete(`${API_BASE}/${name}`)
  },

  // Add disk to storage pool
  async addDisk(name: string, device: string, branchPath: string, mode: string, priority: number): Promise<void> {
    await axios.post(`${API_BASE}/${name}/disks`, {
      device,
      branchPath,
      mode,
      priority
    })
  },

  // Remove disk from storage pool
  async removeDisk(name: string, device: string): Promise<void> {
    await axios.delete(`${API_BASE}/${name}/disks/${device}`)
  },

  // Get pool branches
  async getPoolBranches(name: string): Promise<BranchInfo[]> {
    const response = await axios.get(`${API_BASE}/${name}/branches`)
    return response.data.branches || []
  },

  // Mount storage pool
  async mountPool(name: string): Promise<void> {
    await axios.post(`${API_BASE}/${name}/mount`)
  },

  // Unmount storage pool
  async umountPool(name: string): Promise<void> {
    await axios.post(`${API_BASE}/${name}/umount`)
  },

  // Balance storage pool
  async balancePool(name: string): Promise<void> {
    await axios.post(`${API_BASE}/${name}/balance`)
  },

  // Scan storage pool
  async scanPool(name: string): Promise<PoolStatusResponse> {
    const response = await axios.post(`${API_BASE}/${name}/scan`)
    return response.data
  },

  // Get pool status
  async getPoolStatus(name: string): Promise<PoolStatusResponse> {
    const response = await axios.get(`${API_BASE}/${name}/status`)
    return response.data
  }
}