<template>
  <el-dialog
    v-model="dialogVisible"
    :title="quota ? '编辑用户配额' : '设置用户配额'"
    width="600px"
    @close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item label="用户" prop="username" v-if="!quota">
        <el-select
          v-model="form.username"
          placeholder="选择用户"
          filterable
          style="width: 100%"
        >
          <el-option
            v-for="user in availableUsers"
            :key="user.id"
            :label="user.username"
            :value="user.username"
          >
            <div class="user-option">
              <span>{{ user.username }}</span>
              <span class="user-email">{{ user.email }}</span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="路径" prop="path">
        <el-input
          v-model="form.path"
          placeholder="/home /path/to/directory"
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
      </el-form-item>

      <el-form-item label="文件数限制">
        <el-checkbox v-model="form.hasFileLimits">启用文件数限制</el-checkbox>
      </el-form-item>

      <el-form-item v-if="form.hasFileLimits" label="软限制">
        <el-input-number
          v-model="form.filesSoft"
          :min="0"
          :max="1000000"
          controls-position="right"
          style="width: 150px"
        />
      </el-form-item>

      <el-form-item v-if="form.hasFileLimits" label="硬限制">
        <el-input-number
          v-model="form.filesHard"
          :min="0"
          :max="1000000"
          controls-position="right"
          style="width: 150px"
        />
      </el-form-item>

      <el-form-item label="启用告警">
        <el-switch v-model="form.enableAlerts" />
        <div class="form-tip">
          当接近或超过限制时发送告警通知
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
import type { UserQuota } from '@/types/quota'
import { userApi } from '@/api'

interface Props {
  visible: boolean
  quota: UserQuota | null
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
  username: '',
  path: '/home',
  softLimitValue: 10,
  softLimitUnit: 'GB',
  hardLimitValue: 15,
  hardLimitUnit: 'GB',
  gracePeriod: 7,
  hasFileLimits: false,
  filesSoft: 0,
  filesHard: 0,
  enableAlerts: true
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
  username: [
    { required: true, message: '请选择用户', trigger: 'change' }
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

// Load available users from API
const availableUsers = ref<any[]>([])

const loadAvailableUsers = async () => {
  try {
    const response = await userApi.getUsers()
    if (response && response.users) {
      availableUsers.value = response.users
    } else if (Array.isArray(response)) {
      availableUsers.value = response
    } else {
      availableUsers.value = []
    }
  } catch (error) {
    console.error('Failed to load users:', error)
    availableUsers.value = []
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
      await quotaStore.setUserQuota(props.quota.user.username, {
        path: form.value.path,
        softLimit: softLimit.value,
        hardLimit: hardLimit.value,
        gracePeriod: form.value.gracePeriod
      })
      ElMessage.success('用户配额已更新')
    } else {
      // Create new quota
      await quotaStore.setUserQuota(form.value.username, {
        path: form.value.path,
        softLimit: softLimit.value,
        hardLimit: hardLimit.value,
        gracePeriod: form.value.gracePeriod
      })
      ElMessage.success('用户配额已创建')
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
    username: '',
    path: '/home',
    softLimitValue: 10,
    softLimitUnit: 'GB',
    hardLimitValue: 15,
    hardLimitUnit: 'GB',
    gracePeriod: 7,
    hasFileLimits: false,
    filesSoft: 0,
    filesHard: 0,
    enableAlerts: true
  }
}

// Watch for quota changes
watch(() => props.quota, (quota) => {
  if (quota) {
    form.value = {
      username: quota.user.username,
      path: quota.path,
      softLimitValue: quota.softLimit / 1024 / 1024 / 1024,
      softLimitUnit: 'GB',
      hardLimitValue: quota.hardLimit / 1024 / 1024 / 1024,
      hardLimitUnit: 'GB',
      gracePeriod: quota.gracePeriod,
      hasFileLimits: quota.filesSoft > 0 || quota.filesHard > 0,
      filesSoft: quota.filesSoft,
      filesHard: quota.filesHard,
      enableAlerts: true
    }
  }
})

// Load available users when component mounts
onMounted(() => {
  loadAvailableUsers()
})
</script>

<style scoped lang="scss">
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

.user-option {
  display: flex;
  flex-direction: column;

  .user-email {
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