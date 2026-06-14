<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`格式化磁盘 - ${disk?.device || ''}`"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-alert
      title="危险操作警告"
      type="error"
      :closable="false"
      show-icon
      style="margin-bottom: 20px"
    >
      <ul>
        <li>格式化操作将永久删除磁盘上的所有数据</li>
        <li>请确保已备份重要数据</li>
        <li>格式化过程不可中断</li>
        <li>格式化后需要重新创建分区表</li>
      </ul>
    </el-alert>

    <div class="disk-info" v-if="disk">
      <el-descriptions :column="2" border size="small">
        <el-descriptions-item label="设备路径">
          {{ disk.device }}
        </el-descriptions-item>
        <el-descriptions-item label="设备型号">
          {{ disk.model }}
        </el-descriptions-item>
        <el-descriptions-item label="磁盘容量">
          {{ formatSize(disk.size) }}
        </el-descriptions-item>
        <el-descriptions-item label="当前状态">
          <el-tag :type="getHealthType(disk.health)" size="small">
            {{ getHealthLabel(disk.health) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="已有分区" :span="2">
          {{ disk.partitions?.length || 0 }} 个
        </el-descriptions-item>
      </el-descriptions>
    </div>

    <el-divider />

    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
      <el-form-item label="确认设备" prop="confirmDevice">
        <el-input
          v-model="form.confirmDevice"
          placeholder="请输入设备路径以确认"
        />
        <div class="form-tip">
          请输入 <strong>{{ disk?.device }}</strong> 以确认格式化操作
        </div>
      </el-form-item>

      <el-form-item label="分区表类型" prop="partitionTable">
        <el-select v-model="form.partitionTable" placeholder="选择分区表类型">
          <el-option label="GPT (GUID Partition Table)" value="gpt">
            <div class="option-content">
              <span class="option-label">GPT</span>
              <span class="option-desc">推荐，支持大于2TB的磁盘和更多分区</span>
            </div>
          </el-option>
          <el-option label="MBR (Master Boot Record)" value="mbr">
            <div class="option-content">
              <span class="option-label">MBR</span>
              <span class="option-desc">传统格式，最大支持2TB</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="文件系统" prop="filesystem">
        <el-select v-model="form.filesystem" placeholder="选择文件系统">
          <el-option label="ext4" value="ext4">
            <div class="option-content">
              <span class="option-label">ext4</span>
              <span class="option-desc">推荐，稳定可靠，支持大文件</span>
            </div>
          </el-option>
          <el-option label="xfs" value="xfs">
            <div class="option-content">
              <span class="option-label">xfs</span>
              <span class="option-desc">高性能，适合大文件和高并发场景</span>
            </div>
          </el-option>
          <el-option label="btrfs" value="btrfs">
            <div class="option-content">
              <span class="option-label">btrfs</span>
              <span class="option-desc">支持快照、压缩和数据校验</span>
            </div>
          </el-option>
          <el-option label="ntfs" value="ntfs">
            <div class="option-content">
              <span class="option-label">ntfs</span>
              <span class="option-desc">Windows兼容性</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="挂载选项">
        <el-checkbox v-model="form.autoMount">自动挂载</el-checkbox>
      </el-form-item>

      <el-form-item label="挂载点" v-if="form.autoMount">
        <el-input
          v-model="form.mountPoint"
          placeholder="/mnt/disk-name"
          prefix-icon="Folder"
        />
      </el-form-item>

      <el-form-item label="格式化选项">
        <el-checkbox v-model="form.quickFormat">快速格式化</el-checkbox>
        <div class="form-tip">
          快速格式化仅清除文件表，不检查磁盘坏道
        </div>
      </el-form-item>

      <el-divider />

      <el-form-item label="最终确认" required>
        <el-checkbox v-model="form.confirmFormat">
          我已阅读并理解上述警告，确认要格式化此磁盘
        </el-checkbox>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button
          type="danger"
          @click="formatDisk"
          :loading="formatting"
          :disabled="!canFormat"
        >
          开始格式化
        </el-button>
      </div>
    </template>

    <!-- 进度显示 -->
    <el-dialog
      v-model="showProgress"
      title="格式化进行中"
      width="400px"
      :close-on-click-modal="false"
      :show-close="false"
    >
      <div class="format-progress">
        <el-progress
          :percentage="progress"
          :status="progressStatus"
        />
        <div class="progress-info">
          <p>当前步骤: {{ currentStep }}</p>
          <p>预计剩余时间: {{ estimatedTime }}</p>
        </div>
      </div>
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { DiskInfo } from '@/types/disk'
import { diskAPI } from '@/api/disk'

interface Props {
  visible: boolean
  disk: DiskInfo | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'formatted', disk: DiskInfo): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const formRef = ref()
const formatting = ref(false)
const showProgress = ref(false)
const progress = ref(0)
const progressStatus = ref<'success' | 'exception' | 'warning' | ''>('')
const currentStep = ref('')
const estimatedTime = ref('')

const form = ref({
  confirmDevice: '',
  partitionTable: 'gpt',
  filesystem: 'ext4',
  autoMount: false,
  mountPoint: '',
  quickFormat: true,
  confirmFormat: false
})

const rules = {
  confirmDevice: [
    { required: true, message: '请确认设备路径', trigger: 'blur' },
    {
      validator: (rule: any, value: string) => {
        return value === props.disk?.device
      },
      message: '设备路径不匹配',
      trigger: 'blur'
    }
  ],
  partitionTable: [
    { required: true, message: '请选择分区表类型', trigger: 'change' }
  ],
  filesystem: [
    { required: true, message: '请选择文件系统', trigger: 'change' }
  ]
}

const canFormat = computed(() => {
  return form.value.confirmFormat && 
         form.value.confirmDevice === props.disk?.device &&
         !formatting.value
})

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const getHealthType = (health: string): string => {
  const types: any = { good: 'success', warning: 'warning', failed: 'danger' }
  return types[health] || 'info'
}

const getHealthLabel = (health: string) => {
  const labels: any = { good: '健康', warning: '警告', failed: '故障' }
  return labels[health] || health
}

const formatDisk = async () => {
  if (!props.disk) return

  try {
    await formRef.value.validate()
  } catch {
    return
  }

  formatting.value = true
  showProgress.value = true
  progress.value = 0
  progressStatus.value = ''
  
  try {
    const formatConfig = {
      device: props.disk.device,
      partitionTable: form.value.partitionTable,
      filesystem: form.value.filesystem,
      quickFormat: form.value.quickFormat,
      autoMount: form.value.autoMount,
      mountPoint: form.value.mountPoint
    }

    await diskAPI.formatDisk(formatConfig)
    
    progress.value = 100
    progressStatus.value = 'success'
    currentStep.value = '格式化完成'
    estimatedTime.value = '0秒'

    ElMessage.success(`磁盘 ${props.disk.device} 格式化成功`)
    
    setTimeout(() => {
      emit('formatted', props.disk)
      handleClose()
    }, 1500)

  } catch (error: any) {
    progressStatus.value = 'exception'
    currentStep.value = '格式化失败'
    ElMessage.error(error.message || '格式化失败')
  } finally {
    formatting.value = false
    setTimeout(() => {
      showProgress.value = false
    }, 2000)
  }
}

const handleClose = () => {
  dialogVisible.value = false
  form.value = {
    confirmDevice: '',
    partitionTable: 'gpt',
    filesystem: 'ext4',
    autoMount: false,
    mountPoint: '',
    quickFormat: true,
    confirmFormat: false
  }
  progress.value = 0
  progressStatus.value = ''
}
</script>

<style scoped lang="scss">
.option-content {
  display: flex;
  flex-direction: column;

  .option-label {
    font-weight: bold;
    color: #303133;
  }

  .option-desc {
    font-size: 12px;
    color: #909399;
    margin-top: 2px;
  }
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.disk-info {
  margin-bottom: 20px;
}

.format-progress {
  padding: 20px 0;

  .progress-info {
    margin-top: 20px;
    text-align: center;

    p {
      margin: 8px 0;
      color: #606266;
      font-size: 14px;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>