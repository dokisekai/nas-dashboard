<template>
  <el-dialog
    v-model="dialogVisible"
    :title="quota ? '编辑组配额' : '设置组配额'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="组" prop="groupId" v-if="!quota">
        <el-select
          v-model="form.groupId"
          placeholder="选择组"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="group in availableGroups"
            :key="group.id"
            :label="group.name"
            :value="group.id"
          >
            <div class="group-option">
              <span>{{ group.name }}</span>
              <span class="group-desc">{{ group.description }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="路径" prop="path">
        <el-input
          v-model="form.path"
          placeholder="/shared /path/to/directory"
          prefix-icon="Folder"
        />
        <div class="form-tip">
          指定要应用配额的目录路径
        </div>
      </el-form-item>

      <el-form-item label="软限制" prop="softLimit">
        <div class="size-input">
          <el-input-number
            v-model="form.softLimitValue"
            :min="0"
            :max="10240"
            :step="1"
            controls-position="right"
            style="width: 150px"
          />
          <el-select v-model="form.softLimitUnit" style="width: 80px; margin-left: 10px">
            <el-option label="MB" value="MB" />
            <el-option label="GB" value="GB" />
            <el-option label="TB" value="TB" />
          </el-select>
        </div>
        <div class="form-tip">
          软限制：超过此值会收到警告，但仍可写入数据
        </div>
      </el-form-item>

      <el-form-item label="硬限制" prop="hardLimit">
        <div class="size-input">
          <el-input-number
            v-model="form.hardLimitValue"
            :min="0"
            :max="10240"
            :step="1"
            controls-position="right"
            style="width: 150px"
          />
          <el-select v-model="form.hardLimitUnit" style="width: 80px; margin-left: 10px">
            <el-option label="MB" value="MB" />
            <el-option label="GB" value="GB" />
            <el-option label="TB" value="TB" />
          </el-select>
        </div>
        <div class="form-tip">
          硬限制：达到此值后将无法写入更多数据
        </div>
      </el-form-item>

      <el-form-item label="宽限期" prop="gracePeriod">
        <el-input-number
          v-model="form.gracePeriod"
          :min="0"
          :max="30"
          :step="1"
          controls-position="right"
          style="width: 150px"
        />
        <span class="unit-hint">天</span>
        <div class="form-tip">
          超过软限制后的宽限时间
        </div>
      </el-form-item>

      <el-form-item label="继承用户配额">
        <el-switch v-model="form.inheritUserQuota" />
        <div class="form-tip">
          组成员将继承组配额限制
        </div>
      </el-form-item>

      <el-form-item label="配额预览">
        <div class="quota-preview">
          <div class="preview-info">
            <span class="info-label">软限制:</span>
            <span class="info-value">{{ formatSize(softLimit) }}</span>
          </div>
          <div class="preview-info">
            <span class="info-label">硬限制:</span>
            <span class="info-value">{{ formatSize(hardLimit) }}</span>
          </div>
          <div class="preview-info">
            <span class="info-label">宽限期:</span>
            <span class="info-value">{{ form.gracePeriod }} 天</span>
          </div>
        </div>

            <span class="info-label">继承设置:</span>
            <span class="info-value">{{ form.inheritUserQuota ? '是' : '否' }}</span>
          </div>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="saveQuota" :loading="saving">
          {{ quota ? '更新' : '创建' }}配额
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useQuotaStore } from '@/stores/quota'
import { ElMessage } from 'element-plus'
import type { GroupQuota } from '@/types/quota'
import { groupApi } from '@/api'

interface Props {
  visible: boolean
  quota: GroupQuota | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'saved'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const quotaStore = useQuotaStore()

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const formRef = ref()
const saving = ref(false)

const form = ref({
  groupId: 0,
  path: '/shared',
  softLimitValue: 50,
  softLimitUnit: 'GB',
  hardLimitValue: 100,
  hardLimitUnit: 'GB',
  gracePeriod: 7,
  inheritUserQuota: false
})

const convertToBytes = (value: number, unit: string): number => {
  const units: Record<string, number> = {
    'B': 1,
    'KB': 1024,
    'MB': 1024 * 1024,
    'GB': 1024 * 1024 * 1024,
    'TB': 1024 * 1024 * 1024 * 1024
  }
  return value * (units[unit] || 1)
}

const softLimit = computed(() => convertToBytes(form.value.softLimitValue, form.value.softLimitUnit))
const hardLimit = computed(() => convertToBytes(form.value.hardLimitValue, form.value.hardLimitUnit))

const rules = {
  groupId: [
    { required: true, message: '请选择组', trigger: 'change' }
  ],
  path: [
    { required: true, message: '请输入路径', trigger: 'blur' },
    { pattern: /^\/[\/\w\s.-]+$/, message: '请输入有效的路径', trigger: 'blur' }
  ],
  softLimitValue: [
    { required: true, message: '请输入软限制', trigger: 'blur' }
  ],
  hardLimitValue: [
    { required: true, message: '请输入硬限制', trigger: 'blur' }
  ],
  gracePeriod: [
    { required: true, message: '请输入宽限期', trigger: 'blur' }
  ]
}

// Load available groups from API
const availableGroups = ref<any[]>([])

const loadAvailableGroups = async () => {
  try {
    const response = await groupApi.getGroups()
    if (response && response.groups) {
      availableGroups.value = response.groups
    } else if (Array.isArray(response)) {
      availableGroups.value = response
    } else {
      availableGroups.value = []
    }
  } catch (error) {
    console.error('Failed to load groups:', error)
    availableGroups.value = []
  }
}

// Methods
const formatSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(1)} ${sizes[i]}`
}

const convertToBytes = (value: number, unit: string): number => {
  switch (unit) {
    case 'MB':
      return value * 1024 * 1024
    case 'GB':
      return value * 1024 * 1024 * 1024
    case 'TB':
      return value * 1024 * 1024 * 1024 * 1024
    default:
      return value
  }
}

const saveQuota = async () => {
  try {
    await formRef.value?.validate()

    saving.value = true

    if (props.quota) {
      // Update existing quota
      await quotaStore.setGroupQuota(props.quota.group.id, {
        path: form.value.path,
        softLimit: softLimit.value,
        hardLimit: hardLimit.value,
        gracePeriod: form.value.gracePeriod
      })
      ElMessage.success('组配额已更新')
    } else {
      // Create new quota
      await quotaStore.setGroupQuota(form.value.groupId, {
        path: form.value.path,
        softLimit: softLimit.value,
        hardLimit: hardLimit.value,
        gracePeriod: form.value.gracePeriod
      })
      ElMessage.success('组配额已创建')
    }

    emit('saved')
    handleClose()
  } catch (error) {
    ElMessage.error('保存配额失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  dialogVisible.value = false
  formRef.value?.resetFields()
  form.value = {
    groupId: 0,
    path: '/shared',
    softLimitValue: 50,
    softLimitUnit: 'GB',
    hardLimitValue: 100,
    hardLimitUnit: 'GB',
    gracePeriod: 7,
    inheritUserQuota: false
  }
}

// Watch for quota changes
watch(() => props.quota, (quota) => {
  if (quota) {
    form.value = {
      groupId: quota.group.id,
      path: quota.path,
      softLimitValue: quota.softLimit / 1024 / 1024 / 1024,
      softLimitUnit: 'GB',
      hardLimitValue: quota.hardLimit / 1024 / 1024 / 1024,
      hardLimitUnit: 'GB',
      gracePeriod: quota.gracePeriod,
      inheritUserQuota: false
    }
  }
})

// Load available groups when component mounts
onMounted(() => {
  loadAvailableGroups()
})
</script>

<style scoped lang="scss">
/* Same styles as UserQuotaEditor */
.size-input {
  display: flex;
  align-items: center;
  gap: 10px;
}

.unit-hint {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 5px;
  line-height: 1.4;
}

.quota-preview {
  background: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  margin-top: 10px;

  .preview-info {
    display: flex;
    justify-content: space-between;
    margin-bottom: 8px;

    &:last-child {
      margin-bottom: 0;
    }

    .info-label {
      font-size: 12px;
      color: #909399;
    }

    .info-value {
      font-size: 12px;
      color: #303133;
      font-weight: 500;
    }
  }
}

.group-option {
  display: flex;
  flex-direction: column;

  .group-desc {
    font-size: 11px;
    color: #909399;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
