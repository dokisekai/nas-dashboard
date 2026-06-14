<template>
  <el-dialog
    v-model="dialogVisible"
    :title="`编辑存储池 - ${pool?.name}`"
    width="700px"
    @close="handleClose"
  >
    <div v-if="pool" class="pool-edit-content">
      <el-form :model="form" label-width="120px" ref="formRef">
        <el-form-item label="描述">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="3"
            placeholder="存储池描述信息"
          />
        </el-form-item>

        <div v-if="pool.type === 'mergerfs' && form.config">
          <el-divider>MergerFS 运行配置</el-divider>
          
          <el-form-item label="文件分配策略">
            <el-select v-model="form.config.category" placeholder="选择策略" style="width: 100%">
              <el-option
                v-for="cat in mergerFSCategories"
                :key="cat.value"
                :label="cat.label"
                :value="cat.value"
              >
                <div class="option-content">
                  <span class="option-label">{{ cat.label }}</span>
                  <span class="option-desc">{{ cat.description }}</span>
                </div>
              </el-option>
            </el-select>
            <div class="form-tip">更改策略会实时生效，影响新文件的存放位置</div>
          </el-form-item>

          <el-form-item label="最小空闲空间">
            <el-input v-model="form.config.minfreespace" placeholder="例如: 10G, 100M" />
          </el-form-item>

          <el-form-item label="直接 I/O">
            <el-switch v-model="form.config.direct_io" />
          </el-form-item>

          <el-form-item label="异步读取">
            <el-switch v-model="form.config.async_read" />
          </el-form-item>
          
          <el-form-item label="跟随符号链接">
            <el-switch v-model="form.config.follow_symlinks" />
          </el-form-item>
        </div>
      </el-form>

      <el-alert
        v-if="pool.type === 'mergerfs'"
        title="提示"
        type="info"
        description="修改某些高级选项可能需要存储池短暂重新挂载，正在进行的文件传输可能会中断。"
        :closable="false"
        show-icon
        style="margin-top: 20px"
      />
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="loading">
          保存修改
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { StoragePool } from '@/types/storage_pool'
import { MERGERFS_CATEGORIES } from '@/types/storage_pool'
import { useStoragePoolStore } from '@/stores/storage_pool'

interface Props {
  visible: boolean
  pool: StoragePool | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'updated'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const storagePoolStore = useStoragePoolStore()
const loading = ref(false)
const mergerFSCategories = MERGERFS_CATEGORIES

const form = ref<any>({
  description: '',
  config: null
})

const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

watch(() => props.pool, (newPool) => {
  if (newPool) {
    form.value = {
      description: newPool.description || '',
      config: newPool.config ? JSON.parse(JSON.stringify(newPool.config)) : null
    }
  }
}, { immediate: true })

const handleClose = () => {
  dialogVisible.value = false
}

const handleSubmit = async () => {
  if (!props.pool) return

  loading.value = true
  try {
    await storagePoolStore.updatePool(props.pool.name, form.value)
    ElMessage.success('存储池配置已更新')
    emit('updated')
    handleClose()
  } catch (error) {
    ElMessage.error('更新失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="scss">
.pool-edit-content {
  .option-content {
    display: flex;
    flex-direction: column;

    .option-label {
      font-weight: bold;
    }

    .option-desc {
      font-size: 12px;
      color: #909399;
    }
  }

  .form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 4px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>