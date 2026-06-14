<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`存储池详情 - ${pool?.name || ''}`"
    width="1000px"
    @close="handleClose"
  >
    <div v-if="pool" class="pool-details">
      <!-- Pool Overview -->
      <div class="detail-section">
        <h3>概览信息</h3>
        <el-descriptions :column="3" border>
          <el-descriptions-item label="存储池名称">
            {{ pool.name }}
          </el-descriptions-item>
          <el-descriptions-item label="存储池类型">
            <el-tag size="small">{{ getPoolTypeLabel(pool.type) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(pool.status)" size="small">
              {{ getStatusLabel(pool.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="挂载点" :span="2">
            {{ pool.mountPoint }}
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatDate(pool.createdAt) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- Storage Usage -->
      <div class="detail-section">
        <h3>存储使用情况</h3>
        <div class="storage-visualization">
          <div class="usage-chart">
            <el-progress
              type="circle"
              :percentage="usagePercent"
              :color="getProgressColor(usagePercent)"
              :width="150"
            >
              <template #default="{ percentage }">
                <span class="percentage-value">{{ percentage }}%</span>
                <span class="percentage-label">使用率</span>
              </template>
            </el-progress>
          </div>
          <div class="usage-details">
            <div class="usage-item">
              <span class="label">总容量:</span>
              <span class="value">{{ formatSize(pool.totalSize) }}</span>
            </div>
            <div class="usage-item">
              <span class="label">已使用:</span>
              <span class="value">{{ formatSize(pool.usedSize) }}</span>
            </div>
            <div class="usage-item">
              <span class="label">可用空间:</span>
              <span class="value">{{ formatSize(pool.freeSize) }}</span>
            </div>
            <div class="usage-item">
              <span class="label">使用率:</span>
              <span class="value">{{ usagePercent.toFixed(1) }}%</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Pool Disks -->
      <div class="detail-section">
        <h3>包含磁盘 ({{ pool.poolDisks?.length || 0 }})</h3>
        <el-table :data="pool.poolDisks" style="width: 100%">
          <el-table-column prop="device" label="设备" width="150" />
          <el-table-column prop="size" label="容量" width="120">
            <template #default="{ row }">
              {{ formatSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getDiskStatusType(row.status)" size="small">
                {{ getDiskStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="branchPath" label="分支路径" min-width="200" />
          <el-table-column prop="priority" label="优先级" width="100" align="center" />
          <el-table-column label="使用情况" width="200">
            <template #default="{ row }">
              <div class="disk-usage">
                <el-progress
                  :percentage="getDiskUsagePercent(row)"
                  :color="getProgressColor(getDiskUsagePercent(row))"
                  :show-text="false"
                  :stroke-width="8"
                />
                <span class="usage-text">
                  {{ formatSize(row.used) }} / {{ formatSize(row.size) }}
                </span>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- MergerFS Configuration -->
      <div v-if="pool.type === 'mergerfs' && pool.config" class="detail-section">
        <h3>MergerFS 配置</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="文件分配策略">
            {{ getCategoryLabel(pool.config.category) }}
          </el-descriptions-item>
          <el-descriptions-item label="最小空闲空间">
            {{ pool.config.minfreespace }}
          </el-descriptions-item>
          <el-descriptions-item label="直接 I/O">
            <el-tag :type="pool.config.direct_io ? 'success' : 'info'" size="small">
              {{ pool.config.direct_io ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="异步读取">
            <el-tag :type="pool.config.async_read ? 'success' : 'info'" size="small">
              {{ pool.config.async_read ? '启用' : '禁用' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div class="branches-info">
          <h4>分支配置:</h4>
          <el-table :data="pool.config.branches" size="small" style="width: 100%">
            <el-table-column prop="path" label="路径" />
            <el-table-column prop="mode" label="模式" width="100">
              <template #default="{ row }">
                <el-tag :type="row.mode === 'rw' ? 'success' : 'info'" size="small">
                  {{ row.mode === 'rw' ? '读写' : '只读' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="priority" label="优先级" width="100" align="center" />
          </el-table>
        </div>
      </div>

      <!-- Snapshots -->
      <div v-if="pool.snapshots?.length > 0" class="detail-section">
        <h3>快照 ({{ pool.snapshots.length }})</h3>
        <el-table :data="pool.snapshots" style="width: 100%">
          <el-table-column prop="name" label="快照名称" />
          <el-table-column prop="description" label="描述" />
          <el-table-column prop="size" label="大小" width="120">
            <template #default="{ row }">
              {{ formatSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getSnapshotStatusType(row.status)" size="small">
                {{ getSnapshotStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Description -->
      <div v-if="pool.description" class="detail-section">
        <h3>描述</h3>
        <p class="description-text">{{ pool.description }}</p>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="handleRefresh">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-dropdown @command="handleAction" split-button type="primary">
          操作
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="mount" :disabled="pool?.status === 'active'">
                <el-icon><Connection /></el-icon>
                挂载
              </el-dropdown-item>
              <el-dropdown-item command="umount" :disabled="pool?.status !== 'active'">
                <el-icon><SwitchButton /></el-icon>
                卸载
              </el-dropdown-item>
              <el-dropdown-item command="balance">
                <el-icon><Refresh /></el-icon>
                重新平衡
              </el-dropdown-item>
              <el-dropdown-item divided command="delete" :disabled="pool?.status === 'active'">
                <el-icon><Delete /></el-icon>
                删除存储池
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { StoragePool } from '@/types/storage_pool'
import { MERGERFS_CATEGORIES } from '@/types/storage_pool'
import {
  Refresh,
  Connection,
  SwitchButton,
  Delete
} from '@element-plus/icons-vue'

interface Props {
  visible: boolean
  pool: StoragePool | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'refresh'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const usagePercent = computed(() => {
  if (!props.pool || props.pool.totalSize === 0) return 0
  return (props.pool.usedSize / props.pool.totalSize) * 100
})

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const formatDate = (date: string): string => {
  return new Date(date).toLocaleString('zh-CN')
}

const getProgressColor = (percent: number): string => {
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
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

const getDiskStatusType = (status: string): string => {
  const types: Record<string, string> = {
    active: 'success',
    failed: 'danger',
    removed: 'info',
    ro: 'warning',
    rw: 'success'
  }
  return types[status] || 'info'
}

const getDiskStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    active: '正常',
    failed: '故障',
    removed: '已移除',
    ro: '只读',
    rw: '读写'
  }
  return labels[status] || status
}

const getDiskUsagePercent = (disk: any): number => {
  if (!disk.size || disk.size === 0) return 0
  return (disk.used / disk.size) * 100
}

const getCategoryLabel = (category: string): string => {
  const cat = MERGERFS_CATEGORIES.find(c => c.value === category)
  return cat ? `${cat.label} - ${cat.description}` : category
}

const getSnapshotStatusType = (status: string): string => {
  const types: Record<string, string> = {
    creating: 'warning',
    completed: 'success',
    deleting: 'info',
    error: 'danger'
  }
  return types[status] || 'info'
}

const getSnapshotStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    creating: '创建中',
    completed: '完成',
    deleting: '删除中',
    error: '错误'
  }
  return labels[status] || status
}

const handleClose = () => {
  dialogVisible.value = false
}

const handleRefresh = () => {
  emit('refresh')
  ElMessage.success('数据已刷新')
}

import { useStoragePoolStore } from '@/stores/storage_pool'

const storagePoolStore = useStoragePoolStore()

const handleAction = async (command: string) => {
  if (!props.pool) return

  const poolName = props.pool.name

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

      case 'delete':
        await ElMessageBox.confirm('删除存储池将永久删除所有数据，此操作不可恢复！', '危险操作', {
          type: 'error',
          confirmButtonText: '确定删除',
          cancelButtonText: '取消'
        })
        await storagePoolStore.deletePool(poolName)
        ElMessage.success('存储池已删除')
        handleClose()
        break
    }
    emit('refresh')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}
</script>

<style scoped lang="scss">
.pool-details {
  .detail-section {
    margin-bottom: 30px;

    h3 {
      margin-bottom: 15px;
      color: #303133;
      font-size: 16px;
      border-left: 3px solid #409eff;
      padding-left: 10px;
    }

    h4 {
      margin: 15px 0 10px;
      color: #606266;
      font-size: 14px;
    }
  }

  .storage-visualization {
    display: flex;
    gap: 40px;
    align-items: center;
    padding: 20px;
    background: #f5f7fa;
    border-radius: 8px;

    .usage-chart {
      flex-shrink: 0;

      .percentage-value {
        display: block;
        font-size: 24px;
        font-weight: bold;
        color: #303133;
      }

      .percentage-label {
        display: block;
        font-size: 12px;
        color: #909399;
        margin-top: 5px;
      }
    }

    .usage-details {
      flex: 1;

      .usage-item {
        display: flex;
        justify-content: space-between;
        padding: 10px 0;
        border-bottom: 1px solid #ebeef5;

        &:last-child {
          border-bottom: none;
        }

        .label {
          color: #606266;
        }

        .value {
          font-weight: bold;
          color: #303133;
        }
      }
    }
  }

  .disk-usage {
    .usage-text {
      font-size: 12px;
      color: #606266;
      margin-left: 10px;
    }
  }

  .branches-info {
    margin-top: 20px;
  }

  .description-text {
    color: #606266;
    line-height: 1.6;
    margin: 0;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>