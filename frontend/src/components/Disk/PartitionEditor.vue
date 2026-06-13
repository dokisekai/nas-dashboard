<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`分区编辑 - ${disk?.device || ''}`"
    width="900px"
    @close="handleClose"
  >
    <div class="partition-editor">
      <!-- Partition Table -->
      <div class="partition-table">
        <h3>当前分区表</h3>
        <div class="table-type">
          <el-tag :type="partitionTable.type === 'gpt' ? 'success' : 'warning'" size="small">
            {{ partitionTable.type === 'gpt' ? 'GPT' : partitionTable.type === 'mbr' ? 'MBR' : '无' }}
          </el-tag>
          <span class="free-space">可用空间: {{ formatSize(partitionTable.freeSpace) }}</span>
        </div>

        <el-table :data="partitionTable.partitions" size="small" max-height="300">
          <el-table-column prop="device" label="分区" width="120" />
          <el-table-column prop="size" label="大小" width="100">
            <template #default="{ row }">
              {{ formatSize(row.size) }}
            </template>
          </el-table-column>
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="filesystem" label="文件系统" width="100">
            <template #default="{ row }">
              {{ row.filesystem || '未格式化' }}
            </template>
          </el-table-column>
          <el-table-column prop="mountPoint" label="挂载点" min-width="150">
            <template #default="{ row }">
              {{ row.mountPoint || '未挂载' }}
            </template>
          </el-table-column>
          <el-table-column prop="flags" label="标志" width="150">
            <template #default="{ row }">
              <el-tag
                v-for="flag in row.flags"
                :key="flag"
                size="small"
                style="margin-right: 4px"
              >
                {{ flag }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-button-group>
                <el-button size="small" @click="editPartition(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deletePartition(row)">删除</el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </el-table>

        <div class="table-actions">
          <el-button type="primary" @click="createPartition">
            <el-icon><Plus /></el-icon>
            创建分区
          </el-button>
          <el-button @click="refreshTable">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>

      <!-- Visual Partition Map -->
      <div class="partition-map">
        <h3>磁盘分区图</h3>
        <div class="disk-map">
          <div
            v-for="(partition, index) in partitionMap"
            :key="index"
            class="partition-block"
            :class="{ used: partition.used, free: !partition.used }"
            :style="{ width: partition.percent + '%' }"
            @click="selectPartition(partition)"
          >
            <div class="partition-info">
              <span class="partition-name">{{ partition.name || '未分配' }}</span>
              <span class="partition-size">{{ formatSize(partition.size) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Partition Dialog -->
    <el-dialog
      v-model="showPartitionDialog"
      :title="editingPartition ? '编辑分区' : '创建分区'"
      width="600px"
    >
      <el-form :model="partitionForm" :rules="partitionRules" ref="partitionFormRef" label-width="120px">
        <el-form-item label="分区大小" prop="size">
          <div class="size-input">
            <el-input-number
              v-model="partitionForm.size"
              :min="1"
              :max="maxPartitionSize"
              :step="1024"
              :precision="0"
              controls-position="right"
            />
            <el-select v-model="sizeUnit" style="width: 80px; margin-left: 10px">
              <el-option label="MB" value="MB" />
              <el-option label="GB" value="GB" />
            </el-select>
            <span class="size-tip">最大可用: {{ formatSize(maxPartitionSize) }}</span>
          </div>
        </el-form-item>

        <el-form-item label="分区类型" prop="type">
          <el-select v-model="partitionForm.type">
            <el-option label="主分区" value="primary" />
            <el-option label="逻辑分区" value="logical" />
            <el-option label="扩展分区" value="extended" />
          </el-select>
        </el-form-item>

        <el-form-item label="文件系统" prop="filesystem">
          <el-select v-model="partitionForm.filesystem">
            <el-option label="不创建文件系统" value="" />
            <el-option label="ext4" value="ext4" />
            <el-option label="xfs" value="xfs" />
            <el-option label="btrfs" value="btrfs" />
            <el-option label="ntfs" value="ntfs" />
            <el-option label="fat32" value="fat32" />
          </el-select>
        </el-form-item>

        <el-form-item label="挂载点" v-if="partitionForm.filesystem">
          <el-input
            v-model="partitionForm.mountPoint"
            placeholder="/mnt/partition-name"
            prefix-icon="Folder"
          />
        </el-form-item>

        <el-form-item label="分区标志">
          <el-checkbox-group v-model="partitionForm.flags">
            <el-checkbox label="boot">启动分区</el-checkbox>
            <el-checkbox label="lvm">LVM分区</el-checkbox>
            <el-checkbox label="raid">RAID分区</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showPartitionDialog = false">取消</el-button>
          <el-button type="primary" @click="savePartition" :loading="saving">
            {{ editingPartition ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { diskAPI } from '@/api/disk'
import { ElMessage } from 'element-plus'
import type { DiskInfo, PartitionTable, Partition } from '@/types/disk'
import { Plus, Refresh } from '@element-plus/icons-vue'

interface Props {
  visible: boolean
  disk: DiskInfo | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'updated'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const partitionTable = ref<PartitionTable>({
  type: 'gpt',
  partitions: [],
  freeSpace: 0
})

const showPartitionDialog = ref(false)
const editingPartition = ref<Partition | null>(null)
const saving = ref(false)
const partitionFormRef = ref()
const sizeUnit = ref('GB')

const partitionForm = ref({
  size: 0,
  type: 'logical',
  filesystem: 'ext4',
  mountPoint: '',
  flags: []
})

const partitionRules = {
  size: [
    { required: true, message: '请输入分区大小', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择分区类型', trigger: 'change' }
  ],
  filesystem: [
    { required: true, message: '请选择文件系统', trigger: 'change' }
  ]
}

const partitionMap = computed(() => {
  const map = []
  let currentPosition = 0

  // Add partitions
  partitionTable.value.partitions.forEach(partition => {
    const percent = (partition.size / (props.disk?.size || 1)) * 100
    map.push({
      name: partition.device,
      size: partition.size,
      used: true,
      percent: percent,
      partition: partition
    })
    currentPosition += percent
  })

  // Add free space
  const freePercent = 100 - currentPosition
  if (freePercent > 0) {
    map.push({
      name: '未分配',
      size: partitionTable.value.freeSpace,
      used: false,
      percent: freePercent,
      partition: null
    })
  }

  return map
})

const maxPartitionSize = computed(() => {
  const sizeInGB = partitionTable.value.freeSpace / 1024 / 1024 / 1024
  return Math.floor(sizeInGB)
})

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const refreshTable = async () => {
  if (!props.disk) return

  try {
    // In real implementation, call API
    // partitionTable.value = await diskAPI.getPartitions(props.disk.device)
    ElMessage.success('分区表已刷新')
  } catch (error) {
    ElMessage.error('获取分区表失败')
  }
}

const createPartition = () => {
  editingPartition.value = null
  partitionForm.value = {
    size: 10,
    type: 'logical',
    filesystem: 'ext4',
    mountPoint: '',
    flags: []
  }
  sizeUnit.value = 'GB'
  showPartitionDialog.value = true
}

const editPartition = (partition: Partition) => {
  editingPartition.value = partition
  partitionForm.value = {
    size: partition.size / 1024 / 1024 / 1024, // Convert to GB
    type: 'logical',
    filesystem: partition.filesystem || 'ext4',
    mountPoint: partition.mountPoint || '',
    flags: partition.flags || []
  }
  sizeUnit.value = 'GB'
  showPartitionDialog.value = true
}

const deletePartition = async (partition: Partition) => {
  try {
    await diskAPI.deletePartition(props.disk?.device || '', getPartitionNumber(partition.device))
    ElMessage.success('分区已删除')
    await refreshTable()
  } catch (error) {
    ElMessage.error('删除分区失败')
  }
}

const savePartition = async () => {
  try {
    await partitionFormRef.value?.validate()

    saving.value = true

    if (editingPartition.value) {
      // Update existing partition
      ElMessage.success('分区更新功能即将推出')
    } else {
      // Create new partition
      const sizeBytes = convertSizeToBytes(partitionForm.value.size, sizeUnit.value)

      // In real implementation:
      // await diskAPI.createPartition(props.disk?.device || '', {
      //   start: calculateStartPoint(),
      //   end: calculateEndPoint(),
      //   type: partitionForm.value.type,
      //   filesystem: partitionForm.value.filesystem
      // })

      ElMessage.success('分区创建功能即将推出')
    }

    showPartitionDialog.value = false
  } catch (error) {
    ElMessage.error('保存分区失败')
  } finally {
    saving.value = false
  }
}

const selectPartition = (partition: any) => {
  if (partition.partition) {
    editPartition(partition.partition)
  } else {
    createPartition()
  }
}

const getPartitionNumber = (device: string): number => {
  const match = device.match(/\d+$/)
  return match ? parseInt(match[0]) : 0
}

const convertSizeToBytes = (size: number, unit: string): number => {
  switch (unit) {
    case 'MB':
      return size * 1024 * 1024
    case 'GB':
      return size * 1024 * 1024 * 1024
    default:
      return size
  }
}

const handleClose = () => {
  dialogVisible.value = false
}

// Watch for disk changes
watch(() => props.disk, (newDisk) => {
  if (newDisk) {
    refreshTable()
  }
})
</script>

<style scoped lang="scss">
.partition-editor {
  .partition-table {
    margin-bottom: 30px;

    h3 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 16px;
    }

    .table-type {
      display: flex;
      align-items: center;
      gap: 15px;
      margin-bottom: 15px;
      font-size: 12px;
      color: #606266;

      .free-space {
        font-weight: 500;
      }
    }

    .table-actions {
      margin-top: 15px;
      display: flex;
      gap: 10px;
    }
  }

  .partition-map {
    h3 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 16px;
    }

    .disk-map {
      display: flex;
      height: 60px;
      background: #f5f7fa;
      border-radius: 8px;
      overflow: hidden;

      .partition-block {
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        transition: all 0.3s;
        border-right: 1px solid #fff;
        font-size: 12px;

        &.used {
          background: #409eff;
          color: white;
        }

        &.free {
          background: #e4e7ed;
          color: #606266;
        }

        &:hover {
          opacity: 0.8;
        }

        .partition-info {
          display: flex;
          flex-direction: column;
          align-items: center;
          gap: 2px;

          .partition-name {
            font-size: 11px;
            font-weight: 500;
          }

          .partition-size {
            font-size: 10px;
            opacity: 0.8;
          }
        }
      }
    }
  }
}

.size-input {
  display: flex;
  align-items: center;
  gap: 10px;

  .size-tip {
    font-size: 12px;
    color: #909399;
    margin-left: auto;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>