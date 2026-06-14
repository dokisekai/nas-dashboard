<template>
  <el-dialog
    v-model="visible"
    title="磁盘详细信息"
    width="600px"
    @close="handleClose"
  >
    <div v-if="disk" class="disk-details">
      <div class="detail-item">
        <span class="label">设备名称:</span>
        <span class="value">{{ disk.device }}</span>
      </div>
      <div class="detail-item">
        <span class="label">型号:</span>
        <span class="value">{{ disk.model || '未知' }}</span>
      </div>
      <div class="detail-item">
        <span class="label">容量:</span>
        <span class="value">{{ formatSize(disk.size) }}</span>
      </div>
      <div class="detail-item">
        <span class="label">健康状态:</span>
        <el-tag :type="getHealthType(disk.health)" size="small">
          {{ getHealthLabel(disk.health) }}
        </el-tag>
      </div>
      <div class="detail-item">
        <span class="label">温度:</span>
        <span class="value">{{ disk.temperature }}°C</span>
      </div>
      <div class="detail-item">
        <span class="label">文件系统:</span>
        <span class="value">{{ disk.fstype || '未挂载' }}</span>
      </div>
      <div class="detail-item">
        <span class="label">挂载点:</span>
        <span class="value">{{ disk.mountpoint || '未挂载' }}</span>
      </div>
    </div>
    <template #footer>
      <el-button @click="handleClose">关闭</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

const props = defineProps<{
  visible: boolean
  disk: any
}>()

const emit = defineEmits(['update:visible'])

const handleClose = () => {
  emit('update:visible', false)
}

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getHealthType = (health: string) => {
  const types: any = { good: 'success', warning: 'warning', failed: 'danger' }
  return types[health] || 'info'
}

const getHealthLabel = (health: string) => {
  const labels: any = { good: '健康', warning: '警告', failed: '故障' }
  return labels[health] || health
}
</script>

<style scoped>
.disk-details {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
.detail-item {
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 10px;
}
.label {
  font-weight: bold;
  color: #606266;
}
.value {
  color: #303133;
}
</style>
