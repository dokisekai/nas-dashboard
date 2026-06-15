<template>
  <div class="progress-dialog-overlay" @click.self="$emit('close')">
    <div class="progress-dialog">
      <div class="dialog-header">
        <h3>应用安装进度</h3>
        <button class="btn-close" @click="$emit('close')">
          <XMarkIcon />
        </button>
      </div>

      <div class="dialog-content">
        <div v-for="app in apps" :key="app.id" class="app-progress-item">
          <div class="app-header">
            <div class="app-info">
              <h4 class="app-name">{{ app.displayName || app.name }}</h4>
              <p class="app-status">{{ getStatusText(app) }}</p>
            </div>
            <div class="app-percent">
              {{ getProgress(app) }}%
            </div>
          </div>

          <div class="progress-bar">
            <div
              class="progress-fill"
              :style="{
                width: getProgress(app) + '%',
                background: getStatusColor(app)
              }"
            ></div>
          </div>

          <div class="progress-steps">
            <div
              v-for="step in installSteps"
              :key="step.key"
              :class="['step-item', getStepClass(app, step.key)]"
            >
              <div class="step-icon">
                <CheckIcon v-if="isStepComplete(app, step.key)" />
                <ClockIcon v-else-if="isStepRunning(app, step.key)" />
                <div v-else class="step-dot"></div>
              </div>
              <span class="step-text">{{ step.label }}</span>
            </div>
          </div>
        </div>

        <div v-if="apps.length === 0" class="empty-state">
          <CheckCircleIcon class="success-icon" />
          <p>所有应用安装完成</p>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn-background" @click="$emit('close')">
          后台运行
        </button>
        <button class="btn-close-all" :disabled="hasInstallingApps" @click="$emit('close')">
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  XMarkIcon,
  CheckIcon,
  ClockIcon,
  CheckCircleIcon
} from '@heroicons/vue/24/outline'
import type { AppInstance } from '../../api/application'

interface Props {
  apps: AppInstance[]
}

interface Emits {
  (e: 'close'): void
}

const props = defineProps<Props>()
defineEmits<Emits>()

const installSteps = [
  { key: 'validate', label: '验证应用包' },
  { key: 'download', label: '下载应用' },
  { key: 'extract', label: '解压文件' },
  { key: 'check_deps', label: '检查依赖' },
  { key: 'pre_install', label: '预安装' },
  { key: 'install', label: '安装应用' },
  { key: 'post_install', label: '后安装配置' },
  { key: 'start', label: '启动应用' },
  { key: 'complete', label: '完成' }
]

const hasInstallingApps = computed(() => {
  return props.apps.some(app => app.status === 'installing')
})

const getStatusText = (app: AppInstance) => {
  const statusMap: Record<string, string> = {
    installing: '安装中...',
    installed: '安装成功',
    running: '运行中',
    error: '安装失败',
    stopped: '已停止'
  }
  return statusMap[app.status as string] || '未知状态'
}

const getStatusColor = (app: AppInstance) => {
  if (app.status === 'error') {
    return 'linear-gradient(90deg, #ef4444 0%, #f87171 100%)'
  } else if ((app.status as string) === 'installed' || app.status === 'running') {
    return 'linear-gradient(90deg, #10b981 0%, #34d399 100%)'
  }
  return 'linear-gradient(90deg, #667eea 0%, #764ba2 100%)'
}

const getProgress = (app: AppInstance) => {
  // 这里应该从实际的进度数据获取
  // 简化实现：基于状态估算
  if ((app.status as string) === 'installed' || app.status === 'running') return 100
  if (app.status === 'error') return 0
  if (app.status === 'installing') return 65
  return 0
}

const isStepComplete = (app: AppInstance, step: string) => {
  const progress = getProgress(app)
  const stepIndex = installSteps.findIndex(s => s.key === step)
  return progress >= (stepIndex + 1) * (100 / installSteps.length)
}

const isStepRunning = (app: AppInstance, step: string) => {
  const progress = getProgress(app)
  const stepIndex = installSteps.findIndex(s => s.key === step)
  const stepProgress = progress / (100 / installSteps.length)
  return Math.floor(stepProgress) === stepIndex && app.status === 'installing'
}

const getStepClass = (app: AppInstance, step: string) => {
  if (isStepComplete(app, step)) return 'step-complete'
  if (isStepRunning(app, step)) return 'step-running'
  return 'step-pending'
}
</script>

<style scoped>
.progress-dialog-overlay {
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

.progress-dialog {
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  background: white;
  border-radius: 16px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  display: flex;
  flex-direction: column;
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
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.app-progress-item {
  margin-bottom: 24px;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.app-info {
  flex: 1;
}

.app-name {
  font-size: 16px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 4px;
}

.app-status {
  font-size: 13px;
  color: #6b7280;
}

.app-percent {
  font-size: 20px;
  font-weight: 700;
  color: #667eea;
}

.progress-bar {
  height: 8px;
  background: rgba(102, 126, 234, 0.1);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 16px;
}

.progress-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.3s ease, background 0.3s ease;
}

.progress-steps {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.step-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border-radius: 6px;
  font-size: 13px;
}

.step-complete {
  color: #10b981;
}

.step-running {
  color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.step-pending {
  color: #9ca3af;
}

.step-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-icon svg {
  width: 16px;
  height: 16px;
}

.step-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

.step-running .step-dot {
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.step-text {
  font-weight: 600;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #10b981;
}

.success-icon {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.btn-background,
.btn-close-all {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-background {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.btn-background:hover {
  background: rgba(102, 126, 234, 0.2);
}

.btn-close-all {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-close-all:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-close-all:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>