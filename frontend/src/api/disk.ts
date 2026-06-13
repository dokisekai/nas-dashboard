// Disk Management API Client
import axios from 'axios'
import type {
  DiskInfo,
  PartitionTable,
  RAIDConfig,
  VolumeGroup,
  LogicalVolume,
  DiskBenchmark
} from '@/types/disk'

const API_BASE = '/api/storage'

export const diskAPI = {
  // Basic Disk Operations
  async getDisks(): Promise<DiskInfo[]> {
    const response = await axios.get(`${API_BASE}/disks`)
    return response.data.disks || []
  },

  async getDisk(device: string): Promise<DiskInfo> {
    const response = await axios.get(`${API_BASE}/disks/${encodeURIComponent(device)}`)
    return response.data
  },

  async mountDisk(device: string, mountPoint: string, type?: string): Promise<void> {
    await axios.post(`${API_BASE}/mount`, { device, mountPoint, type })
  },

  async umountDisk(mountPoint: string): Promise<void> {
    await axios.post(`${API_BASE}/umount`, { mountPoint })
  },

  // Partition Management
  async getPartitions(device: string): Promise<PartitionTable> {
    const response = await axios.get(`${API_BASE}/disks/${encodeURIComponent(device)}/partitions`)
    return response.data
  },

  async createPartition(device: string, partition: {
    start: number
    end: number
    type: string
    filesystem?: string
  }): Promise<void> {
    await axios.post(`${API_BASE}/disks/${encodeURIComponent(device)}/partitions`, partition)
  },

  async deletePartition(device: string, partitionNumber: number): Promise<void> {
    await axios.delete(`${API_BASE}/disks/${encodeURIComponent(device)}/partitions/${partitionNumber}`)
  },

  // RAID Management
  async getRAIDArrays(): Promise<RAIDConfig[]> {
    const response = await axios.get(`${API_BASE}/raid`)
    return response.data.arrays || []
  },

  async getRAIDArray(name: string): Promise<RAIDConfig> {
    const response = await axios.get(`${API_BASE}/raid/${encodeURIComponent(name)}`)
    return response.data
  },

  async createRAID(config: {
    name: string
    level: string
    devices: string[]
  }): Promise<RAIDConfig> {
    const response = await axios.post(`${API_BASE}/raid`, config)
    return response.data
  },

  async deleteRAID(name: string): Promise<void> {
    await axios.delete(`${API_BASE}/raid/${encodeURIComponent(name)}`)
  },

  async addDiskToRAID(name: string, device: string): Promise<void> {
    await axios.post(`${API_BASE}/raid/${encodeURIComponent(name)}/add`, { device })
  },

  async removeDiskFromRAID(name: string, device: string): Promise<void> {
    await axios.post(`${API_BASE}/raid/${encodeURIComponent(name)}/remove`, { device })
  },

  // LVM Management
  async getPhysicalVolumes(): Promise<any[]> {
    const response = await axios.get(`${API_BASE}/lvm/pv`)
    return response.data.pvs || []
  },

  async createPhysicalVolume(device: string): Promise<void> {
    await axios.post(`${API_BASE}/lvm/pv`, { device })
  },

  async getVolumeGroups(): Promise<VolumeGroup[]> {
    const response = await axios.get(`${API_BASE}/lvm/vg`)
    return response.data.vgs || []
  },

  async createVolumeGroup(config: {
    name: string
    devices: string[]
  }): Promise<VolumeGroup> {
    const response = await axios.post(`${API_BASE}/lvm/vg`, config)
    return response.data
  },

  async deleteVolumeGroup(name: string): Promise<void> {
    await axios.delete(`${API_BASE}/lvm/vg/${encodeURIComponent(name)}`)
  },

  async getLogicalVolumes(vgName?: string): Promise<LogicalVolume[]> {
    const params = vgName ? { vg: vgName } : {}
    const response = await axios.get(`${API_BASE}/lvm/lv`, { params })
    return response.data.lvs || []
  },

  async createLogicalVolume(config: {
    name: string
    vgName: string
    size: number
    filesystem?: string
  }): Promise<LogicalVolume> {
    const response = await axios.post(`${API_BASE}/lvm/lv`, config)
    return response.data
  },

  async deleteLogicalVolume(name: string, vgName: string): Promise<void> {
    await axios.delete(`${API_BASE}/lvm/lv/${encodeURIComponent(vgName)}/${encodeURIComponent(name)}`)
  },

  // SMART Monitoring
  async getSMARTInfo(device: string): Promise<any> {
    const response = await axios.get(`${API_BASE}/disks/${encodeURIComponent(device)}/smart`)
    return response.data
  },

  async runSMARTTest(device: string, type: 'short' | 'long' | 'conveyance'): Promise<void> {
    await axios.post(`${API_BASE}/disks/${encodeURIComponent(device)}/test`, { type })
  },

  // Performance Benchmarking
  async runBenchmark(device: string): Promise<DiskBenchmark> {
    const response = await axios.post(`${API_BASE}/disks/${encodeURIComponent(device)}/benchmark`)
    return response.data
  },

  // Health Status
  async getDiskHealth(device: string): Promise<any> {
    const response = await axios.get(`${API_BASE}/disks/${encodeURIComponent(device)}/health`)
    return response.data
  }
}