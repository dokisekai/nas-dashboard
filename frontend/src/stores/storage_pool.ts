// Storage Pool Pinia Store
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { storagePoolAPI } from '@/api/storage_pool'
import type { StoragePool, StoragePoolCreateRequest, PoolStatusResponse } from '@/types/storage_pool'

export const useStoragePoolStore = defineStore('storagePool', () => {
  // State
  const pools = ref<StoragePool[]>([])
  const currentPool = ref<StoragePool | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const poolStatus = ref<Record<string, PoolStatusResponse>>({})

  // Computed
  const activePools = computed(() =>
    pools.value.filter(pool => pool.status === 'active')
  )

  const totalStorage = computed(() =>
    pools.value.reduce((sum, pool) => sum + pool.totalSize, 0)
  )

  const usedStorage = computed(() =>
    pools.value.reduce((sum, pool) => sum + pool.usedSize, 0)
  )

  const storageUsagePercent = computed(() =>
    totalStorage.value > 0 ? (usedStorage.value / totalStorage.value) * 100 : 0
  )

  const poolsByType = computed(() => {
    const grouped: Record<string, StoragePool[]> = {}
    pools.value.forEach(pool => {
      if (!grouped[pool.type]) {
        grouped[pool.type] = []
      }
      grouped[pool.type].push(pool)
    })
    return grouped
  })

  const degradedPools = computed(() =>
    pools.value.filter(pool => pool.status === 'degraded' || pool.status === 'error')
  )

  // Actions
  async function fetchPools() {
    loading.value = true
    error.value = null
    try {
      pools.value = await storagePoolAPI.getPools()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch storage pools'
      console.error('Error fetching storage pools:', err)
    } finally {
      loading.value = false
    }
  }

  async function fetchPool(name: string) {
    loading.value = true
    error.value = null
    try {
      currentPool.value = await storagePoolAPI.getPool(name)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch storage pool'
      console.error('Error fetching storage pool:', err)
    } finally {
      loading.value = false
    }
  }

  async function createPool(data: StoragePoolCreateRequest) {
    loading.value = true
    error.value = null
    try {
      const newPool = await storagePoolAPI.createPool(data)
      pools.value.push(newPool)
      return newPool
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create storage pool'
      console.error('Error creating storage pool:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updatePool(name: string, data: Partial<StoragePoolCreateRequest>) {
    loading.value = true
    error.value = null
    try {
      const updatedPool = await storagePoolAPI.updatePool(name, data)
      const index = pools.value.findIndex(p => p.name === name)
      if (index !== -1) {
        pools.value[index] = updatedPool
      }
      if (currentPool.value?.name === name) {
        currentPool.value = updatedPool
      }
      return updatedPool
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update storage pool'
      console.error('Error updating storage pool:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function deletePool(name: string) {
    loading.value = true
    error.value = null
    try {
      await storagePoolAPI.deletePool(name)
      pools.value = pools.value.filter(p => p.name !== name)
      if (currentPool.value?.name === name) {
        currentPool.value = null
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete storage pool'
      console.error('Error deleting storage pool:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function refreshPoolStatus(name: string) {
    try {
      const status = await storagePoolAPI.scanPool(name)
      poolStatus.value[name] = status
      return status
    } catch (err) {
      console.error('Error refreshing pool status:', err)
      throw err
    }
  }

  async function addDisk(name: string, device: string, branchPath: string, mode: string, priority: number) {
    loading.value = true
    error.value = null
    try {
      await storagePoolAPI.addDisk(name, device, branchPath, mode, priority)
      await fetchPool(name)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to add disk'
      console.error('Error adding disk:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function removeDisk(name: string, device: string) {
    loading.value = true
    error.value = null
    try {
      await storagePoolAPI.removeDisk(name, device)
      await fetchPool(name)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to remove disk'
      console.error('Error removing disk:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function mountPool(name: string) {
    loading.value = true
    error.value = null
    try {
      await storagePoolAPI.mountPool(name)
      await fetchPool(name)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to mount pool'
      console.error('Error mounting pool:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function umountPool(name: string) {
    loading.value = true
    error.value = null
    try {
      await storagePoolAPI.umountPool(name)
      await fetchPool(name)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to unmount pool'
      console.error('Error unmounting pool:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  async function balancePool(name: string) {
    loading.value = true
    error.value = null
    try {
      await storagePoolAPI.balancePool(name)
      await fetchPool(name)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to balance pool'
      console.error('Error balancing pool:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  // Initialize
  function init() {
    fetchPools()
  }

  return {
    // State
    pools,
    currentPool,
    loading,
    error,
    poolStatus,

    // Computed
    activePools,
    totalStorage,
    usedStorage,
    storageUsagePercent,
    poolsByType,
    degradedPools,

    // Actions
    fetchPools,
    fetchPool,
    createPool,
    updatePool,
    deletePool,
    refreshPoolStatus,
    addDisk,
    removeDisk,
    mountPool,
    umountPool,
    balancePool,
    clearError,
    init
  }
})