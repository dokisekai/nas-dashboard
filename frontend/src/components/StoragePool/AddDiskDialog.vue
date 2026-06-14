<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`向存储池 ${pool?.name} 添加磁盘`"
    width="600px"
    @close="handleClose"
  >
    <div v-if="pool" class="add-disk-content">
      <el-form :model="form" label-width="100px">
        <el-form-item label="选择磁盘">
          <el-select v-model="form.device" placeholder="请选择要添加的磁盘" style="width: 100%">
            <el-option
              v-for="disk in availableDisks"
              :key="disk.device"
              :label="`${disk.device} (${formatSize(disk.size)})`"
              :value="disk.device"
            >
              <div class="disk-option">
                <span>{{ disk.device }}</span>
                <span class="disk-size">{{ formatSize(disk.size) }}</span>
                <el-tag v-if="disk.status !== 'available'" type="warning" size="small">占用</el-tag>
              </div>
            </el-option>
          </el-select>
          <div class="form-tip">只能添加未分配或未挂载的物理磁盘</div>
        </el-form-item>

        <el-form-item label="读写模式">
          <el-radio-group v-model="form.mode">
            <el-radio label="rw">读写</el-radio>
            <el-radio label="ro">只读</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="自动格式化">
          <el-checkbox v-model="form.format">格式化磁盘为 ext4 并合并</el-checkbox>
          <div class="form-tip">新磁盘必须格式化后才能加入存储池</div>
        </el-form-item>

        <el-form-item label="优先级">
          <el-input-number v-model="form.priority" :min="0" :max="100" />
          <div class="form-tip">优先级越高，MergerFS 会越优先在上面创建文件</div>
        </el-form-item>
      </el-form>

      <div class="warning-box">
        <el-alert
          title="警告：磁盘将被格式化"
          type="warning"
          description="添加磁盘到存储池会自动对该磁盘进行分区并格式化。磁盘上原有的所有数据都将丢失！"
          :closable="false"
          show-icon
        />
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="primary"
          @click="handleSubmit"
          :loading="loading"
          :disabled="!form.device"
        >
          确认添加并合并
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { StoragePool } from '@/types/storage_pool'
import { storageApi } from '@/api'
import { useStoragePoolStore } from '@/stores/storage_pool'

interface Props {
  visible: boolean
  pool: StoragePool | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'added'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const storagePoolStore = useStoragePoolStore()
const loading = ref(false)
const availableDisks = ref<any[]>([])

const form = ref({
  device: '',
  mode: 'rw',
  priority: 50,
  format: true
})

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const loadAvailableDisks = async () => {
  try {
    const response = await storageApi.getDisks()
    // 过滤掉已经在池中的磁盘
    const existingDevices = new Set(props.pool?.poolDisks.map(d => d.device) || [])
    
    availableDisks.value = (response.disks || []).filter((disk: any) => {
      return !existingDevices.has(disk.name)
    }).map((disk: any) => ({
      device: disk.name,
      size: disk.size,
      status: disk.mounted ? 'mounted' : 'available'
    }))
  } catch (error) {
    console.error('Failed to fetch disks:', error)
  }
}

watch(() => props.visible, (newVal) => {
  if (newVal && props.pool) {
    loadAvailableDisks()
    // Reset form
    form.value = {
      device: '',
      mode: 'rw',
      priority: 50,
      format: true
    }
  }
})

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const handleClose = () => {
  dialogVisible.value = false
}

const handleSubmit = async () => {
  if (!props.pool || !form.value.device) return

  try {
    await ElMessageBox.confirm(
      `确定要将磁盘 ${form.value.device} 合并到存储池 ${props.pool.name} 吗？此磁盘将被格式化！`,
      '确认合并',
      {
        confirmButtonText: '确定合并',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    loading.value = true
    await storagePoolStore.addDisk(
      props.pool.name,
      form.value.device,
      form.value.device, // branchPath matches device for now
      form.value.mode,
      form.value.priority,
      form.value.format
    )

    ElMessage.success('磁盘添加成功，存储池已合并')
    emit('added')
    handleClose()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('添加磁盘失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.add-disk-content {
  .disk-option {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;

    .disk-size {
      color: #909399;
      font-size: 12px;
      margin-left: auto;
      margin-right: 10px;
    }
  }

  .form-tip {
    font-size: 12px;
    color: #909399;
    line-height: 1.4;
    margin-top: 4px;
  }

  .warning-box {
    margin-top: 20px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>