<template>
  <div class="upload-dialog-overlay" @click.self="$emit('cancel')">
    <div class="upload-dialog">
      <div class="dialog-header">
        <h3>上传应用包</h3>
        <button class="btn-close" @click="$emit('cancel')">
          <XMarkIcon />
        </button>
      </div>

      <div class="dialog-content">
        <div
          class="upload-zone"
          :class="{ 'drag-over': isDragOver }"
          @drop.prevent="handleDrop"
          @dragover.prevent="isDragOver = true"
          @dragleave.prevent="isDragOver = false"
          @click="selectFile"
        >
          <ArrowUpTrayIcon class="upload-icon" />
          <p class="upload-text">拖拽文件到此处或点击选择</p>
          <p class="upload-hint">支持 .nap 格式的应用包</p>
        </div>

        <div v-if="selectedFile" class="file-info">
          <DocumentIcon class="file-icon" />
          <div class="file-details">
            <p class="file-name">{{ selectedFile.name }}</p>
            <p class="file-size">{{ formatSize(selectedFile.size) }}</p>
          </div>
          <button class="btn-remove" @click="selectedFile = null">
            <TrashIcon />
          </button>
        </div>

        <div v-if="uploadProgress.show" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress.percent + '%' }"></div>
          </div>
          <p class="progress-text">{{ uploadProgress.text }}</p>
        </div>

        <div v-if="uploadError" class="upload-error">
          <ExclamationTriangleIcon class="error-icon" />
          <p>{{ uploadError }}</p>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn-cancel" @click="$emit('cancel')">取消</button>
        <button
          class="btn-upload-confirm"
          :disabled="!selectedFile || uploadProgress.show"
          @click="handleUpload"
        >
          <ArrowUpTrayIcon />
          <span>上传</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  XMarkIcon,
  ArrowUpTrayIcon,
  DocumentIcon,
  TrashIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'

interface Emits {
  (e: 'upload', file: File): void
  (e: 'cancel'): void
}

defineEmits<Emits>()

const isDragOver = ref(false)
const selectedFile = ref<File | null>(null)
const uploadProgress = ref({
  show: false,
  percent: 0,
  text: ''
})
const uploadError = ref('')

const selectFile = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.nap'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) {
      selectedFile.value = file
    }
  }
  input.click()
}

const handleDrop = (e: DragEvent) => {
  isDragOver.value = false
  const file = e.dataTransfer?.files[0]
  if (file && file.name.endsWith('.nap')) {
    selectedFile.value = file
  } else {
    uploadError.value = '请上传 .nap 格式的应用包'
  }
}

const handleUpload = () => {
  if (selectedFile.value) {
    uploadProgress.value = {
      show: true,
      percent: 0,
      text: '上传中...'
    }

    // 模拟上传进度
    let percent = 0
    const interval = setInterval(() => {
      percent += 10
      uploadProgress.value.percent = percent

      if (percent >= 100) {
        clearInterval(interval)
        uploadProgress.value.text = '上传完成'
      }
    }, 200)

    // 这里应该调用实际的API
    setTimeout(() => {
      clearInterval(interval)
      uploadProgress.value.show = false
      if (selectedFile.value) {
        // 触发上传事件
        // this.$emit('upload', selectedFile.value)
      }
    }, 2000)
  }
}

const formatSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}
</script>

<style scoped>
.upload-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.upload-dialog {
  width: 90%;
  max-width: 500px;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.dialog-header h3 {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
}

.btn-close {
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-close:hover {
  background: rgba(102, 126, 234, 0.2);
}

.btn-close svg {
  width: 20px;
  height: 20px;
  color: #667eea;
}

.dialog-content {
  padding: 24px;
}

.upload-zone {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  border: 2px dashed rgba(102, 126, 234, 0.3);
  border-radius: 12px;
  background: rgba(102, 126, 234, 0.02);
  cursor: pointer;
  transition: all 0.2s;
}

.upload-zone:hover,
.upload-zone.drag-over {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.upload-icon {
  width: 48px;
  height: 48px;
  color: #667eea;
  margin-bottom: 16px;
}

.upload-text {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.upload-hint {
  font-size: 14px;
  color: #6b7280;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  margin-top: 16px;
}

.file-icon {
  width: 32px;
  height: 32px;
  color: #667eea;
}

.file-details {
  flex: 1;
}

.file-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.file-size {
  font-size: 12px;
  color: #6b7280;
}

.btn-remove {
  width: 28px;
  height: 28px;
  border: none;
  background: rgba(239, 68, 68, 0.1);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-remove:hover {
  background: rgba(239, 68, 68, 0.2);
}

.btn-remove svg {
  width: 16px;
  height: 16px;
  color: #ef4444;
}

.upload-progress {
  margin-top: 16px;
}

.progress-bar {
  height: 8px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 13px;
  color: #667eea;
  font-weight: 600;
  text-align: center;
}

.upload-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  background: rgba(239, 68, 68, 0.05);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 8px;
  margin-top: 16px;
  color: #ef4444;
  font-size: 14px;
}

.error-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.btn-cancel {
  padding: 10px 20px;
  border: none;
  background: rgba(107, 114, 128, 0.1);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel:hover {
  background: rgba(107, 114, 128, 0.2);
}

.btn-upload-confirm {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border: none;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-upload-confirm:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-upload-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-upload-confirm svg {
  width: 16px;
  height: 16px;
}
</style>