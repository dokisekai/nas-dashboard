<template>
  <el-dialog
    v-model="dialogVisible"
    title="创建RAID阵列"
    width="700px"
    @close="handleClose"
  >
    <el-steps :active="currentStep" align-center class="wizard-steps">
      <el-step title="RAID类型" />
      <el-step title="选择磁盘" />
      <el-step title="配置参数" />
      <el-step title="确认创建" />
    </el-steps>

    <div class="wizard-content">
      <!-- Step 1: RAID Level Selection -->
      <div v-show="currentStep === 0" class="wizard-step">
        <h3>选择RAID级别</h3>
        <div class="raid-levels">
          <div
            v-for="level in raidLevels"
            :key="level.value"
            class="raid-level-card"
            :class="{ selected: form.level === level.value }"
            @click="selectRAIDLevel(level.value)"
          >
            <div class="card-header">
              <h4>{{ level.label }}</h4>
              <el-tag :type="level.recommended ? 'success' : 'info'" size="small">
                {{ level.recommended ? '推荐' : '可用' }}
              </el-tag>
            </div>
            <div class="card-body">
              <p class="description">{{ level.description }}</p>
              <div class="requirements">
                <div class="requirement">
                  <el-icon><Document /></el-icon>
                  <span>最少磁盘: {{ level.minDisks }}</span>
                </div>
                <div class="requirement">
                  <el-icon><TrendCharts /></el-icon>
                  <span>冗余保护: {{ level.redundancy ? '是' : '否' }}</span>
                </div>
                <div class="requirement">
                  <el-icon><TrendCharts /></el-icon>
                  <span>性能提升: {{ level.performanceBoost }}x</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Step 2: Disk Selection -->
      <div v-show="currentStep === 1" class="wizard-step">
        <h3>选择成员磁盘</h3>
        <div class="disk-selection">
          <div class="available-disks">
            <el-table
              :data="availableDisks"
              @selection-change="handleDiskSelection"
              style="width: 100%"
            >
              <el-table-column type="selection" width="55" :selectable="isDiskSelectable" />
              <el-table-column prop="device" label="设备" width="120" />
              <el-table-column prop="model" label="型号" min-width="150" />
              <el-table-column prop="size" label="容量" width="120">
                <template #default="{ row }">
                  {{ formatSize(row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="health" label="健康状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="getHealthType(row.health)" size="small">
                    {{ getHealthLabel(row.health) }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>

            <div class="selection-info">
              <el-alert
                :title="`已选择 ${selectedDisks.length} 个磁盘，总容量: ${formatSize(totalSelectedCapacity)}`"
                :type="selectedDisks.length >= minDisks ? 'success' : 'warning'"
                :closable="false"
                show-icon
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Step 3: Configuration -->
      <div v-show="currentStep === 2" class="wizard-step">
        <h3>RAID配置</h3>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
          <el-form-item label="阵列名称" prop="name">
            <el-input
              v-model="form.name"
              placeholder="请输入RAID阵列名称"
              maxlength="50"
              show-word-limit
            />
          </el-form-item>

          <el-form-item label="RAID级别" prop="level">
            <el-input :value="getRAIDLevelLabel(form.level)" disabled />
          </el-form-item>

          <el-form-item label="成员磁盘">
            <div class="selected-disks-list">
              <el-tag
                v-for="disk in selectedDisks"
                :key="disk.device"
                closable
                @close="removeDisk(disk)"
              >
                {{ disk.device }} ({{ formatSize(disk.size) }})
              </el-tag>
            </div>
          </el-form-item>

          <el-form-item label="阵列容量">
            <el-input :value="formatSize(arrayCapacity)" disabled />
            <div class="form-tip">
              根据RAID级别和磁盘大小自动计算
            </div>
          </el-form-item>

          <el-form-item label="文件系统">
            <el-select v-model="form.filesystem" placeholder="选择文件系统">
              <el-option label="不创建文件系统" value="" />
              <el-option label="ext4" value="ext4" />
              <el-option label="xfs" value="xfs" />
              <el-option label="btrfs" value="btrfs" />
            </el-select>
          </el-form-item>

          <el-form-item label="挂载点" v-if="form.filesystem">
            <el-input
              v-model="form.mountPoint"
              placeholder="/mnt/raid-name"
              prefix-icon="Folder"
            />
          </el-form-item>

          <el-form-item label="创建选项">
            <el-checkbox v-model="form.forceCreate">强制创建（可能损坏数据）</el-checkbox>
            <div class="form-tip">
              将格式化所选磁盘上的所有数据
            </div>
          </el-form-item>
        </el-form>
      </div>

      <!-- Step 4: Confirmation -->
      <div v-show="currentStep === 3" class="wizard-step">
        <h3>确认RAID配置</h3>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="阵列名称">
            {{ form.name }}
          </el-descriptions-item>
          <el-descriptions-item label="RAID级别">
            {{ getRAIDLevelLabel(form.level) }}
          </el-descriptions-item>
          <el-descriptions-item label="成员磁盘">
            {{ selectedDisks.length }} 个
          </el-descriptions-item>
          <el-descriptions-item label="阵列容量">
            {{ formatSize(arrayCapacity) }}
          </el-descriptions-item>
          <el-descriptions-item label="文件系统" :span="2">
            {{ form.filesystem || '不创建文件系统' }}
          </el-descriptions-item>
          <el-descriptions-item label="挂载点" :span="2" v-if="form.mountPoint">
            {{ form.mountPoint }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="disk-list">
          <h4>成员磁盘列表:</h4>
          <el-table :data="selectedDisks" size="small">
            <el-table-column prop="device" label="设备" />
            <el-table-column prop="model" label="型号" />
            <el-table-column prop="size" label="容量">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="health" label="健康状态">
              <template #default="{ row }">
                <el-tag :type="getHealthType(row.health)" size="small">
                  {{ getHealthLabel(row.health) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="warning-box">
          <el-alert
            title="重要警告"
            type="error"
            :closable="false"
            show-icon
          >
            <ul>
              <li>创建RAID将格式化所有选中的磁盘，所有数据将被永久删除</li>
              <li>请确保已备份重要数据</li>
              <li>RAID创建过程不可中断</li>
              <li>创建过程中请勿关闭系统或重启</li>
            </ul>
          </el-alert>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="wizard-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button v-if="currentStep > 0" @click="previousStep">上一步</el-button>
        <el-button
          v-if="currentStep < steps.length - 1"
          type="primary"
          @click="nextStep"
          :disabled="!canNextStep"
        >
          下一步
        </el-button>
        <el-button
          v-else
          type="primary"
          @click="createRAID"
          :loading="creating"
        >
          创建RAID
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { diskAPI } from '@/api/disk'
import { ElMessage } from 'element-plus'
import type { DiskInfo } from '@/types/disk'
import { RAID_LEVEL_LABELS, RAID_LEVEL_DESCRIPTIONS } from '@/types/disk'
import { Document, TrendCharts } from '@element-plus/icons-vue'

interface Props {
  visible: boolean
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'created', raid: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const steps = ['RAID类型', '选择磁盘', '配置参数', '确认创建']
const currentStep = ref(0)
const creating = ref(false)
const formRef = ref()

const form = ref({
  name: '',
  level: '1',
  filesystem: 'ext4',
  mountPoint: '',
  forceCreate: false
})

const rules = {
  name: [
    { required: true, message: '请输入RAID阵列名称', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_-]+$/, message: '名称只能包含字母、数字、下划线和连字符', trigger: 'blur' }
  ],
  level: [
    { required: true, message: '请选择RAID级别', trigger: 'change' }
  ]
}

const raidLevels = [
  {
    value: '0',
    label: 'RAID 0 - 条带',
    description: '提供最大性能和容量，但无冗余保护',
    minDisks: 2,
    redundancy: false,
    performanceBoost: 2,
    recommended: false
  },
  {
    value: '1',
    label: 'RAID 1 - 镜像',
    description: '提供数据镜像保护，容量利用率50%',
    minDisks: 2,
    redundancy: true,
    performanceBoost: 1,
    recommended: true
  },
  {
    value: '5',
    label: 'RAID 5 - 带奇偶校验',
    description: '提供较好的性能和冗余平衡',
    minDisks: 3,
    redundancy: true,
    performanceBoost: 1.5,
    recommended: true
  },
  {
    value: '6',
    label: 'RAID 6 - 双奇偶校验',
    description: '允许两块磁盘同时故障',
    minDisks: 4,
    redundancy: true,
    performanceBoost: 1.3,
    recommended: false
  },
  {
    value: '10',
    label: 'RAID 10 - 镜像+条带',
    description: '结合RAID 0和RAID 1的优势',
    minDisks: 4,
    redundancy: true,
    performanceBoost: 2,
    recommended: true
  }
]

// Mock disk data - in real implementation, fetch from API
const availableDisks = ref<DiskInfo[]>([
  { device: '/dev/sdb', model: 'Samsung SSD', size: 1024 * 1024 * 1024 * 512, health: 'good', temperature: 35, partitions: [] },
  { device: '/dev/sdc', model: 'Western Digital', size: 1024 * 1024 * 1024 * 1024, health: 'good', temperature: 40, partitions: [] },
  { device: '/dev/sdd', model: 'Seagate HDD', size: 1024 * 1024 * 1024 * 2048, health: 'warning', temperature: 45, partitions: [] },
  { device: '/dev/sde', model: 'Toshiba HDD', size: 1024 * 1024 * 1024 * 1024, health: 'good', temperature: 38, partitions: [] }
])

const selectedDisks = ref<DiskInfo[]>([])

const minDisks = computed(() => {
  const level = raidLevels.find(l => l.value === form.value.level)
  return level?.minDisks || 2
})

const totalSelectedCapacity = computed(() => {
  return selectedDisks.value.reduce((sum, disk) => sum + disk.size, 0)
})

const arrayCapacity = computed(() => {
  const level = form.value.level
  const diskCount = selectedDisks.value.length

  if (diskCount === 0) return 0

  switch (level) {
    case '0':
      return totalSelectedCapacity.value
    case '1':
      return totalSelectedCapacity.value / 2
    case '5':
      return totalSelectedCapacity.value * (diskCount - 1) / diskCount
    case '6':
      return totalSelectedCapacity.value * (diskCount - 2) / diskCount
    case '10':
      return totalSelectedCapacity.value / 2
    default:
      return totalSelectedCapacity.value
  }
})

const canNextStep = computed(() => {
  switch (currentStep.value) {
    case 0:
      return form.value.level !== ''
    case 1:
      return selectedDisks.value.length >= minDisks.value
    case 2:
      return form.value.name !== ''
    default:
      return true
  }
})

const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const getRAIDLevelLabel = (level: string): string => {
  return (RAID_LEVEL_LABELS as any)[level] || level
}

const getHealthType = (health: string): string => {
  const types: any = { good: 'success', warning: 'warning', failed: 'danger' }
  return types[health] || 'info'
}

const getHealthLabel = (health: string): string => {
  const labels: any = { good: '健康', warning: '警告', failed: '故障' }
  return labels[health] || health
}

const selectRAIDLevel = (level: string) => {
  form.value.level = level
}

const isDiskSelectable = (row: DiskInfo): boolean => {
  return row.health === 'good'
}

const handleDiskSelection = (selection: DiskInfo[]) => {
  selectedDisks.value = selection.filter(d => d.health === 'good')
}

const removeDisk = (disk: DiskInfo) => {
  selectedDisks.value = selectedDisks.value.filter(d => d.device !== disk.device)
}

const nextStep = () => {
  if (currentStep.value === 0) {
    if (selectedDisks.value.length === 0) {
      ElMessage.warning('请先选择成员磁盘')
      return
    }
  }
  currentStep.value++
}

const previousStep = () => {
  currentStep.value--
}

const createRAID = async () => {
  creating.value = true

  try {
    const raidConfig = {
      name: form.value.name,
      level: form.value.level,
      devices: selectedDisks.value.map(d => d.device)
    }

    const newRAID = await diskAPI.createRAID(raidConfig)

    // Handle filesystem creation if specified
    if (form.value.filesystem) {
      // Additional steps for filesystem creation would go here
    }

    ElMessage.success('RAID阵列创建成功')
    emit('created', newRAID)
    handleClose()
  } catch (error) {
    ElMessage.error('创建RAID阵列失败')
  } finally {
    creating.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
  currentStep.value = 0
  selectedDisks.value = []
  form.value = {
    name: '',
    level: '1',
    filesystem: 'ext4',
    mountPoint: '',
    forceCreate: false
  }
}

// Watch for RAID level changes to update requirements
watch(() => form.value.level, () => {
  // Update minimum disks display
})
</script>

<style scoped lang="scss">
.wizard-steps {
  margin: 20px 0 30px;
}

.wizard-content {
  min-height: 400px;
  padding: 20px 0;
}

.wizard-step {
  h3 {
    margin-bottom: 20px;
    color: #303133;
    font-size: 16px;
  }
}

.raid-levels {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;

  .raid-level-card {
    background: white;
    border: 2px solid #ebeef5;
    border-radius: 8px;
    padding: 20px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      border-color: #409eff;
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    }

    &.selected {
      border-color: #409eff;
      background: #ecf5ff;
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 15px;

      h4 {
        margin: 0;
        color: #303133;
        font-size: 16px;
      }
    }

    .card-body {
      .description {
        color: #606266;
        font-size: 14px;
        line-height: 1.5;
        margin-bottom: 15px;
      }

      .requirements {
        .requirement {
          display: flex;
          align-items: center;
          gap: 8px;
          margin-bottom: 8px;
          font-size: 12px;
          color: #606266;

          .el-icon {
            font-size: 14px;
            color: #409eff;
          }
        }
      }
    }
  }
}

.disk-selection {
  .available-disks {
    margin-bottom: 15px;

    .selection-info {
      margin-top: 15px;
    }
  }
}

.selected-disks-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
}

.disk-list {
  margin: 20px 0;

  h4 {
    margin: 0 0 15px;
    color: #303133;
    font-size: 14px;
  }
}

.warning-box {
  margin: 20px 0;

  ul {
    margin: 10px 0;
    padding-left: 20px;

    li {
      margin: 5px 0;
      color: #303133;
    }
  }
}

.wizard-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>