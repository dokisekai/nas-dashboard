<template>
  <div class="storage-pool-manager">
    <!-- Header Section -->
    <div class="pool-header">
      <div class="header-content">
        <h2>存储池管理</h2>
        <div class="header-actions">
          <el-button type="primary" @click="showCreateDialog = true">
            <el-icon><Plus /></el-icon>
            创建存储池
          </el-button>
          <el-button @click="refreshPools">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>

      <!-- Storage Overview -->
      <div class="storage-overview">
        <div class="overview-card">
          <div class="card-label">总存储容量</div>
          <div class="card-value">{{ formatSize(storagePoolStore.totalStorage) }}</div>
        </div>
        <div class="overview-card">
          <div class="card-label">已使用</div>
          <div class="card-value">{{ formatSize(storagePoolStore.usedStorage) }}</div>
        </div>
        <div class="overview-card">
          <div class="card-label">使用率</div>
          <div class="card-value">{{ storagePoolStore.storageUsagePercent.toFixed(1) }}%</div>
        </div>
        <div class="overview-card">
          <div class="card-label">活动存储池</div>
          <div class="card-value">{{ storagePoolStore.activePools.length }}</div>
        </div>
      </div>
    </div>

    <!-- Storage Progress Bar -->
    <div class="storage-progress">
      <el-progress
        :percentage="storagePoolStore.storageUsagePercent"
        :color="getProgressColor(storagePoolStore.storageUsagePercent)"
        :show-text="false"
      />
    </div>

    <!-- Storage Pool List -->
    <div class="pool-list">
      <el-table
        :data="storagePoolStore.pools"
        v-loading="storagePoolStore.loading"
        @row-click="viewPoolDetails"
        style="width: 100%"
      >
        <el-table-column prop="name" label="存储池名称" width="200">
          <template #default="{ row }">
            <div class="pool-name-cell">
              <el-icon class="pool-icon" :class="getPoolIconClass(row.type)">
                <component :is="getPoolIcon(row.type)" />
              </el-icon>
              <span>{{ row.name }}</span>
              <el-tag
                v-if="row.status === 'degraded' || row.status === 'error'"
                type="danger"
                size="small"
                class="status-badge"
              >
                警告
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ getPoolTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="存储使用" width="300">
          <template #default="{ row }">
            <div class="usage-cell">
              <el-progress
                :percentage="getUsagePercent(row)"
                :color="getProgressColor(getUsagePercent(row))"
                :show-text="false"
                class="usage-progress"
              />
              <span class="usage-text">
                {{ formatSize(row.usedSize) }} / {{ formatSize(row.totalSize) }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="mountPoint" label="挂载点" min-width="200" />

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click.stop="viewPoolDetails(row)">
                <el-icon><View /></el-icon>
              </el-button>
              <el-button size="small" @click.stop="editPool(row)" :disabled="row.status === 'creating' || row.status === 'deleting'">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button size="small" @click.stop="showPoolMenu(row)">
                <el-icon><MoreFilled /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Create Pool Dialog -->
    <StoragePoolWizard
      v-model:visible="showCreateDialog"
      @created="onPoolCreated"
    />

    <!-- Pool Details Dialog -->
    <PoolDetailsDialog
      v-model:visible="showDetailsDialog"
      :pool="selectedPool"
      @refresh="onPoolRefresh"
    />

    <!-- Add Disk Dialog -->
    <AddDiskDialog
      v-model:visible="showAddDiskDialog"
      :pool="selectedPool"
      @added="onPoolRefresh"
    />

    <!-- Edit Pool Dialog -->
    <PoolEditDialog
      v-model:visible="showEditDialog"
      :pool="selectedPool"
      @updated="onPoolRefresh"
    />

    <!-- Context Menu -->
    <el-dropdown
      ref="contextMenu"
      :teleported="true"
      trigger="click"
      @command="handleContextCommand"
    >
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="mount" :disabled="selectedPool?.status === 'active'">
            <el-icon><Connection /></el-icon>
            挂载
          </el-dropdown-item>
          <el-dropdown-item command="umount" :disabled="selectedPool?.status !== 'active'">
            <el-icon><SwitchButton /></el-icon>
            卸载
          </el-dropdown-item>
          <el-dropdown-item command="balance">
            <el-icon><Refresh /></el-icon>
            重新平衡
          </el-dropdown-item>
          <el-dropdown-item command="scan">
            <el-icon><Search /></el-icon>
            扫描状态
          </el-dropdown-item>
          <el-dropdown-item command="add-disk" :disabled="selectedPool?.type !== 'mergerfs'">
            <el-icon><Plus /></el-icon>
            合并硬盘 (Add Disk)
          </el-dropdown-item>
          <el-dropdown-item divided command="delete" :disabled="selectedPool?.status === 'active'">
            <el-icon><Delete /></el-icon>
            删除存储池
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useStoragePoolStore } from '@/stores/storage_pool'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { StoragePool } from '@/types/storage_pool'
import {
  Plus,
  Refresh,
  View,
  Edit,
  MoreFilled,
  Connection,
  SwitchButton,
  Delete,
  Search,
  Folder,
  Coin,
  Box,
  Files
} from '@element-plus/icons-vue'
import StoragePoolWizard from '@/components/StoragePool/PoolWizard.vue'
import PoolDetailsDialog from '@/components/StoragePool/PoolDetailsDialog.vue'
import AddDiskDialog from '@/components/StoragePool/AddDiskDialog.vue'
import PoolEditDialog from '@/components/StoragePool/PoolEditDialog.vue'

const storagePoolStore = useStoragePoolStore()

// Dialog states
const showCreateDialog = ref(false)
const showDetailsDialog = ref(false)
const showAddDiskDialog = ref(false)
const showEditDialog = ref(false)
const selectedPool = ref<StoragePool | null>(null)

// Initialize
onMounted(() => {
  storagePoolStore.init()
})

// Methods
const refreshPools = async () => {
  await storagePoolStore.fetchPools()
  ElMessage.success('存储池列表已刷新')
}

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const getUsagePercent = (pool: StoragePool): number => {
  if (pool.totalSize === 0) return 0
  return (pool.usedSize / pool.totalSize) * 100
}

const getProgressColor = (percent: number): string => {
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

const getPoolIcon = (type: string) => {
  switch (type) {
    case 'mergerfs': return Folder
    case 'btrfs': return Coin
    case 'zfs': return Box
    case 'lvm': return Files
    default: return Folder
  }
}

const getPoolIconClass = (type: string): string => {
  return `pool-icon-${type}`
}

const getPoolTypeLabel = (type: string): string => {
  const labels: Record<string, string> = {
    mergerfs: 'MergerFS',
    btrfs: 'Btrfs',
    zfs: 'ZFS',
    lvm: 'LVM'
  }
  return labels[type] || type
}

const getStatusType = (status: string): string => {
  const types: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    creating: 'warning',
    deleting: 'warning',
    error: 'danger',
    degraded: 'warning'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    active: '活动',
    inactive: '非活动',
    creating: '创建中',
    deleting: '删除中',
    error: '错误',
    degraded: '降级'
  }
  return labels[status] || status
}

const viewPoolDetails = async (pool: StoragePool) => {
  selectedPool.value = pool
  await storagePoolStore.fetchPool(pool.name)
  showDetailsDialog.value = true
}

const editPool = (pool: StoragePool) => {
  selectedPool.value = pool
  showEditDialog.value = true
}

const showPoolMenu = (pool: StoragePool) => {
  selectedPool.value = pool
  // Show context menu at click position
}

const handleContextCommand = async (command: string) => {
  if (!selectedPool.value) return

  const poolName = selectedPool.value.name

  try {
    switch (command) {
      case 'mount':
        await storagePoolStore.mountPool(poolName)
        ElMessage.success('存储池已挂载')
        break

      case 'umount':
        await ElMessageBox.confirm('确定要卸载此存储池吗？', '确认操作', {
          type: 'warning'
        })
        await storagePoolStore.umountPool(poolName)
        ElMessage.success('存储池已卸载')
        break

      case 'balance':
        await ElMessageBox.confirm('重新平衡可能需要较长时间，确定要继续吗？', '确认操作', {
          type: 'warning'
        })
        await storagePoolStore.balancePool(poolName)
        ElMessage.success('存储池重新平衡已完成')
        break

      case 'scan':
        await storagePoolStore.refreshPoolStatus(poolName)
        ElMessage.success('存储池状态已更新')
        break

      case 'add-disk':
        showAddDiskDialog.value = true
        break

      case 'delete':
        await ElMessageBox.confirm('删除存储池将永久删除所有数据，此操作不可恢复！', '危险操作', {
          type: 'error',
          confirmButtonText: '确定删除',
          cancelButtonText: '取消'
        })
        await storagePoolStore.deletePool(poolName)
        ElMessage.success('存储池已删除')
        break
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const onPoolCreated = (pool: StoragePool) => {
  ElMessage.success(`存储池 "${pool.name}" 创建成功`)
  storagePoolStore.fetchPools()
}

const onPoolRefresh = () => {
  storagePoolStore.fetchPools()
}

defineExpose({
  showCreateDialog
})
</script>

<style scoped lang="scss">
.storage-pool-manager {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.pool-header {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h2 {
    margin: 0;
    color: #303133;
    font-size: 24px;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.storage-overview {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 20px;

  .overview-card {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    .card-label {
      color: #909399;
      font-size: 14px;
      margin-bottom: 8px;
    }

    .card-value {
      color: #303133;
      font-size: 24px;
      font-weight: bold;
    }
  }
}

.storage-progress {
  margin-bottom: 20px;
}

.pool-list {
  flex: 1;
  overflow: auto;

  .pool-name-cell {
    display: flex;
    align-items: center;
    gap: 8px;

    .pool-icon {
      font-size: 18px;
      color: #409eff;

      &.pool-icon-mergerfs { color: #67c23a; }
      &.pool-icon-btrfs { color: #e6a23c; }
      &.pool-icon-zfs { color: #f56c6c; }
      &.pool-icon-lvm { color: #909399; }
    }

    .status-badge {
      margin-left: 8px;
    }
  }

  .usage-cell {
    display: flex;
    flex-direction: column;
    gap: 4px;

    .usage-progress {
      width: 100%;
    }

    .usage-text {
      font-size: 12px;
      color: #606266;
    }
  }
}
</style>