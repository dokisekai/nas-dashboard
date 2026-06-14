<template>
  <div class="app-card">
    <div class="app-header-card">
      <div class="app-icon" :style="{ background: appColor }">
        <component :is="appIcon" class="icon" />
      </div>
      <div class="app-status" :class="statusClass">
        <div class="status-dot"></div>
        <span class="status-text">{{ statusText }}</span>
      </div>
    </div>

    <div class="app-info">
      <h3 class="app-name">{{ app.displayName || app.name }}</h3>
      <p class="app-description">{{ app.description || '暂无描述' }}</p>
      <div class="app-meta">
        <span class="app-version">v{{ app.version }}</span>
        <span class="app-category">{{ categoryName }}</span>
        <span class="app-author">{{ app.author || '未知作者' }}</span>
      </div>
    </div>

    <div class="app-actions">
      <button
        v-if="installed"
        :class="['action-btn', 'btn-start']"
        :disabled="app.status === 'running'"
        @click="$emit('start', app)"
      >
        <PlayIcon />
        <span>启动</span>
      </button>

      <button
        v-if="installed"
        :class="['action-btn', 'btn-stop']"
        :disabled="app.status === 'stopped'"
        @click="$emit('stop', app)"
      >
        <StopIcon />
        <span>停止</span>
      </button>

      <button
        v-if="installed"
        :class="['action-btn', 'btn-restart']"
        :disabled="app.status === 'stopped'"
        @click="$emit('restart', app)"
      >
        <RefreshIcon />
        <span>重启</span>
      </button>

      <button
        v-if="installed"
        :class="['action-btn', 'btn-settings']"
        @click="$emit('settings', app)"
      >
        <CogIcon />
        <span>设置</span>
      </button>

      <button
        v-if="installed"
        :class="['action-btn', 'btn-uninstall']"
        @click="$emit('uninstall', app)"
      >
        <TrashIcon />
        <span>卸载</span>
      </button>

      <button
        v-if="!installed"
        :class="['action-btn', 'btn-install']"
        @click="$emit('install', app)"
      >
        <DownloadIcon />
        <span>安装</span>
      </button>

      <button
        :class="['action-btn', 'btn-info']"
        @click="$emit('info', app)"
      >
        <InformationIcon />
        <span>详情</span>
      </button>
    </div>

    <!-- 安装进度 -->
    <div v-if="showProgress" class="install-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <span class="progress-text">{{ progressText }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  PlayIcon,
  StopIcon,
  RefreshIcon,
  CogIcon,
  TrashIcon,
  DownloadIcon,
  InformationIcon,
  CubeIcon
} from '@heroicons/vue/24/outline'
import type { AppInstance, AppPackage } from '../../api/application'

interface Props {
  app: AppInstance | AppPackage
  installed: boolean
  progress?: number
  progressText?: string
}

interface Emits {
  (e: 'start', app: AppInstance): void
  (e: 'stop', app: AppInstance): void
  (e: 'restart', app: AppInstance): void
  (e: 'settings', app: AppInstance): void
  (e: 'uninstall', app: AppInstance): void
  (e: 'install', app: AppPackage): void
  (e: 'info', app: AppInstance | AppPackage): void
}

const props = defineProps<Props>()
defineEmits<Emits>()

// 应用颜色和图标
const appColor = computed(() => {
  const colors = {
    media: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)',
    productivity: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)',
    utilities: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)',
    security: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)',
    network: 'linear-gradient(135deg, #30cfd0 0%, #330867 100%)',
    default: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  }

  const category = (props.app as any).category || 'default'
  return colors[category] || colors.default
})

const appIcon = computed(() => {
  // 这里可以根据应用类型返回不同的图标组件
  return CubeIcon
})

// 应用状态
const statusClass = computed(() => {
  if (!props.installed) return 'status-available'

  const status = (props.app as AppInstance).status
  return `status-${status}`
})

const statusText = computed(() => {
  if (!props.installed) return '可安装'

  const status = (props.app as AppInstance).status
  const statusMap = {
    running: '运行中',
    stopped: '已停止',
    error: '错误',
    installing: '安装中',
    uninstalling: '卸载中'
  }
  return statusMap[status] || '未知'
})

// 分类名称
const categoryName = computed(() => {
  const categoryMap = {
    media: '媒体',
    productivity: '办公',
    utilities: '工具',
    security: '安全',
    network: '网络'
  }

  const category = (props.app as any).category || 'default'
  return categoryMap[category] || '其他'
})

// 进度显示
const showProgress = computed(() => {
  return props.installed &&
         (props.app as AppInstance).status === 'installing' &&
         props.progress !== undefined
})
</script>

<style scoped>
.app-card {
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.app-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.app-header-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.app-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.app-icon .icon {
  width: 32px;
  height: 32px;
}

.app-status {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.status-available {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.status-running {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
}

.status-stopped {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.status-error {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.status-installing {
  background: rgba(59, 130, 246, 0.1);
  color: #3b82f6;
}

.status-uninstalling {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
}

.status-running .status-dot {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.app-info {
  flex: 1;
}

.app-name {
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.app-description {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.app-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  font-size: 12px;
  color: #9ca3af;
}

.app-version,
.app-category,
.app-author {
  padding: 4px 8px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 6px;
  font-weight: 600;
}

.app-actions {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn svg {
  width: 16px;
  height: 16px;
}

.btn-start {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  color: white;
}

.btn-start:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
}

.btn-start:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-stop {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
  color: white;
}

.btn-stop:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.btn-stop:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-restart {
  background: linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%);
  color: white;
}

.btn-restart:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
}

.btn-restart:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-settings {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.btn-settings:hover {
  background: rgba(102, 126, 234, 0.2);
}

.btn-uninstall {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.btn-uninstall:hover {
  background: rgba(239, 68, 68, 0.2);
}

.btn-install {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  grid-column: 1 / -1;
}

.btn-install:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-info {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.btn-info:hover {
  background: rgba(107, 114, 128, 0.2);
}

.install-progress {
  padding: 12px;
  background: rgba(59, 130, 246, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(59, 130, 246, 0.2);
}

.progress-bar {
  height: 6px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 8px;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
  border-radius: 3px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #3b82f6;
  font-weight: 600;
  text-align: center;
  display: block;
}
</style>