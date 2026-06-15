<template>
  <div class="detail-dialog-overlay" @click.self="$emit('close')">
    <div class="detail-dialog">
      <div class="dialog-header">
        <div class="app-header-info">
          <div class="app-icon" :style="{ background: appColor }">
            <CubeIcon class="icon" />
          </div>
          <div class="app-details">
            <h3 class="app-name">{{ app.displayName || app.name }}</h3>
            <p class="app-version">v{{ app.version }}</p>
          </div>
        </div>
        <button class="btn-close" @click="$emit('close')">
          <XMarkIcon />
        </button>
      </div>

      <div class="dialog-content">
        <div class="detail-section">
          <h4 class="section-title">应用信息</h4>
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">描述</span>
              <span class="info-value">{{ app.description || '暂无描述' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">作者</span>
              <span class="info-value">{{ app.author || '未知作者' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">官网</span>
              <a v-if="app.website" :href="app.website" target="_blank" class="info-value link">
                {{ app.website }}
                <ArrowTopRightOnSquareIcon class="link-icon" />
              </a>
              <span v-else class="info-value">-</span>
            </div>
            <div class="info-item">
              <span class="info-label">分类</span>
              <span class="info-value">{{ categoryName }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">许可证</span>
              <span class="info-value">{{ app.license || '未知' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">架构</span>
              <span class="info-value">{{ app.architecture || '通用' }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h4 class="section-title">系统要求</h4>
          <div class="requirements-grid">
            <div class="requirement-item">
              <CpuChipIcon class="requirement-icon" />
              <div class="requirement-info">
                <span class="requirement-label">内存</span>
                <span class="requirement-value">{{ app.minRAM || 0 }} MB</span>
              </div>
            </div>
            <div class="requirement-item">
              <CircleStackIcon class="requirement-icon" />
              <div class="requirement-info">
                <span class="requirement-label">磁盘空间</span>
                <span class="requirement-value">{{ app.minDiskSpace || 0 }} GB</span>
              </div>
            </div>
            <div v-if="dependencies.length > 0" class="requirement-item">
              <CodeBracketIcon class="requirement-icon" />
              <div class="requirement-info">
                <span class="requirement-label">依赖</span>
                <span class="requirement-value">{{ dependencies.join(', ') }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h4 class="section-title">安装统计</h4>
          <div class="stats-grid">
            <div class="stat-item">
              <ArrowDownTrayIcon class="stat-icon" />
              <div class="stat-info">
                <span class="stat-value">{{ app.downloadCount || 0 }}</span>
                <span class="stat-label">下载次数</span>
              </div>
            </div>
            <div class="stat-item">
              <CubeIcon class="stat-icon" />
              <div class="stat-info">
                <span class="stat-value">{{ app.installCount || 0 }}</span>
                <span class="stat-label">安装次数</span>
              </div>
            </div>
            <div class="stat-item">
              <StarIcon class="stat-icon" />
              <div class="stat-info">
                <span class="stat-value">{{ (app.rating || 0).toFixed(1) }}</span>
                <span class="stat-label">用户评分</span>
              </div>
            </div>
          </div>
        </div>

        <div v-if="isInstalled" class="detail-section">
          <h4 class="section-title">运行状态</h4>
          <div class="status-info">
            <div class="status-indicator" :class="statusClass">
              <div class="status-dot"></div>
              <span class="status-text">{{ statusText }}</span>
            </div>
            <div class="status-details">
              <div class="detail-item">
                <span class="detail-label">安装路径</span>
                <span class="detail-value">{{ app.installPath || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">数据路径</span>
                <span class="detail-value">{{ app.dataPath || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">容器ID</span>
                <span class="detail-value">{{ app.containerId || '-' }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">进程ID</span>
                <span class="detail-value">{{ app.pid || '-' }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="dialog-footer">
        <button class="btn-cancel" @click="$emit('close')">关闭</button>
        <button v-if="!isInstalled" class="btn-install" @click="$emit('install', app)">
          <ArrowDownTrayIcon />
          <span>安装应用</span>
        </button>
        <div v-else class="installed-actions">
          <button
            :class="['btn-action', 'btn-start']"
            :disabled="app.status === 'running'"
            @click="$emit('start', app)"
          >
            <PlayIcon />
            <span>启动</span>
          </button>
          <button
            :class="['btn-action', 'btn-stop']"
            :disabled="app.status === 'stopped'"
            @click="$emit('stop', app)"
          >
            <StopIcon />
            <span>停止</span>
          </button>
          <button
            :class="['btn-action', 'btn-restart']"
            @click="$emit('restart', app)"
          >
            <ArrowPathIcon />
            <span>重启</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import {
  XMarkIcon,
  CubeIcon,
  ArrowTopRightOnSquareIcon,
  CpuChipIcon,
  CircleStackIcon,
  CodeBracketIcon,
  ArrowDownTrayIcon,
  StarIcon,
  PlayIcon,
  StopIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline'
import type { AppInstance, AppPackage } from '../../api/application'

interface Props {
  app: any
}

interface Emits {
  (e: 'close'): void
  (e: 'install', app: AppPackage): void
  (e: 'start', app: AppInstance): void
  (e: 'stop', app: AppInstance): void
  (e: 'restart', app: AppInstance): void
}

const props = defineProps<Props>()
defineEmits<Emits>()

const isInstalled = computed(() => 'status' in props.app)

const dependencies = computed(() => {
  return (props.app as any).dependencies || []
})

const categoryName = computed(() => {
  const categoryMap: Record<string, string> = {
    media: '媒体',
    productivity: '办公',
    utilities: '工具',
    security: '安全',
    network: '网络'
  }
  const category = (props.app as any).category || 'default'
  return categoryMap[category] || '其他'
})

const appColor = computed(() => {
  const colors: Record<string, string> = {
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

const statusClass = computed(() => {
  if (!isInstalled.value) return ''
  const status = (props.app as AppInstance).status
  return `status-${status}`
})

const statusText = computed(() => {
  if (!isInstalled.value) return ''

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
</script>

<style scoped>
.detail-dialog-overlay {
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

.detail-dialog {
  width: 90%;
  max-width: 700px;
  max-height: 85vh;
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
  padding: 24px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.app-header-info {
  display: flex;
  align-items: center;
  gap: 16px;
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

.app-details {
  flex: 1;
}

.app-name {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 4px;
}

.app-version {
  font-size: 14px;
  color: #6b7280;
  font-weight: 600;
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

.detail-section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 13px;
  color: #6b7280;
  font-weight: 600;
}

.info-value {
  font-size: 14px;
  color: #1f2937;
}

.info-value.link {
  color: #667eea;
  display: flex;
  align-items: center;
  gap: 4px;
  text-decoration: none;
}

.link-icon {
  width: 16px;
  height: 16px;
}

.requirements-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.requirement-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
}

.requirement-icon {
  width: 24px;
  height: 24px;
  color: #667eea;
}

.requirement-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.requirement-label {
  font-size: 12px;
  color: #6b7280;
}

.requirement-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
}

.stat-icon {
  width: 32px;
  height: 32px;
  color: #667eea;
}

.stat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
}

.stat-label {
  font-size: 12px;
  color: #6b7280;
}

.status-info {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
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

.status-details {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: #6b7280;
  font-weight: 600;
}

.detail-value {
  font-size: 13px;
  color: #1f2937;
  font-family: monospace;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 24px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.btn-cancel,
.btn-install,
.btn-action {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-cancel {
  background: rgba(107, 114, 128, 0.1);
  color: #6b7280;
}

.btn-cancel:hover {
  background: rgba(107, 114, 128, 0.2);
}

.btn-install {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-install:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.installed-actions {
  display: flex;
  gap: 8px;
}

.btn-action {
  padding: 12px 20px;
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

.btn-restart:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(245, 158, 11, 0.3);
}

.btn-action svg {
  width: 16px;
  height: 16px;
}
</style>