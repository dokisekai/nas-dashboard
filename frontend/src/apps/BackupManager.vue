<template>
  <div class="backup-manager">
    <div class="backup-header">
      <h1>备份与恢复</h1>
      <p class="subtitle">管理系统备份和数据恢复</p>
    </div>

    <!-- Quick Actions -->
    <div class="quick-actions">
      <button class="action-btn primary" @click="showCreateBackup = true">
        <PlusIcon class="w-5 h-5" />
        创建备份
      </button>
      <button class="action-btn secondary" @click="refreshBackups">
        <ArrowPathIcon class="w-5 h-5" />
        刷新列表
      </button>
    </div>

    <!-- Backups List -->
    <div class="backups-list">
      <div v-if="loading" class="loading-state">
        <CogIcon class="w-8 h-8 spinning" />
        <p>加载备份列表...</p>
      </div>

      <div v-else-if="backups.length === 0" class="empty-state">
        <ArchiveBoxIcon class="w-16 h-16" />
        <h3>暂无备份</h3>
        <p>点击"创建备份"按钮创建第一个系统备份</p>
      </div>

      <div v-else class="backups-grid">
        <div v-for="backup in backups" :key="backup.id" class="backup-card">
          <div class="backup-icon">
            <ArchiveBoxIcon class="w-8 h-8" />
          </div>

          <div class="backup-info">
            <h3>{{ backup.name }}</h3>
            <p class="backup-description">{{ backup.description || '无描述' }}</p>

            <div class="backup-meta">
              <div class="meta-item">
                <CalendarIcon class="w-4 h-4" />
                {{ formatDate(backup.createdAt) }}
              </div>
              <div class="meta-item">
                <CubeIcon class="w-4 h-4" />
                {{ formatSize(backup.size) }}
              </div>
              <div class="meta-item">
                <TagIcon class="w-4 h-4" />
                {{ backup.type || 'full' }}
              </div>
              <div class="meta-item" :class="'status-' + backup.status">
                <CheckCircleIcon v-if="backup.status === 'completed'" class="w-4 h-4" />
                <ClockIcon v-else-if="backup.status === 'in_progress'" class="w-4 h-4 spinning" />
                <XCircleIcon v-else class="w-4 h-4" />
                {{ getStatusText(backup.status) }}
              </div>
            </div>
          </div>

          <div class="backup-actions">
            <button class="icon-btn" @click="downloadBackup(backup)" title="下载">
              <ArrowDownTrayIcon class="w-5 h-5" />
            </button>
            <button class="icon-btn danger" @click="confirmDeleteBackup(backup)" title="删除">
              <TrashIcon class="w-5 h-5" />
            </button>
            <button
              v-if="backup.status === 'completed'"
              class="icon-btn warning"
              @click="restoreBackup(backup)"
              title="恢复"
            >
              <ArrowUturnLeftIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Backup Modal -->
    <div v-if="showCreateBackup" class="modal-overlay" @click.self="showCreateBackup = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>创建备份</h2>
          <button class="close-btn" @click="showCreateBackup = false">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>备份名称</label>
            <input
              v-model="backupForm.name"
              type="text"
              class="form-input"
              placeholder="例如: weekly-backup"
            />
          </div>

          <div class="form-group">
            <label>备份类型</label>
            <select v-model="backupForm.type" class="form-input">
              <option value="full">完整备份</option>
              <option value="incremental">增量备份</option>
              <option value="differential">差异备份</option>
            </select>
          </div>

          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="backupForm.description"
              class="form-input"
              rows="3"
              placeholder="备份描述信息..."
            ></textarea>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="backupForm.includeDB" type="checkbox" />
              <span>包含数据库</span>
            </label>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="backupForm.includeFiles" type="checkbox" />
              <span>包含文件</span>
            </label>
          </div>

          <div v-if="backupForm.includeFiles" class="form-group">
            <label>文件路径</label>
            <input
              v-model="backupForm.filePaths"
              type="text"
              class="form-input"
              placeholder="/data,/home (逗号分隔多个路径)"
            />
          </div>
        </div>

        <div class="modal-footer">
          <button class="action-btn secondary" @click="showCreateBackup = false">
            取消
          </button>
          <button class="action-btn primary" @click="createBackup" :disabled="creatingBackup">
            <CogIcon v-if="creatingBackup" class="w-5 h-5 spinning" />
            {{ creatingBackup ? '创建中...' : '创建备份' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Restore Confirmation Modal -->
    <div v-if="showRestoreConfirm" class="modal-overlay" @click.self="showRestoreConfirm = false">
      <div class="modal-content small" @click.stop>
        <div class="modal-header">
          <h2>确认恢复</h2>
          <button class="close-btn" @click="showRestoreConfirm = false">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>

        <div class="modal-body">
          <div class="warning-message">
            <ExclamationTriangleIcon class="w-8 h-8" />
            <p>您确定要恢复此备份吗？此操作将覆盖当前系统数据，且不可撤销！</p>
            <p><strong>备份名称:</strong> {{ selectedBackup?.name }}</p>
            <p><strong>创建时间:</strong> {{ formatDate(selectedBackup?.createdAt) }}</p>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="restoreConfirm" type="checkbox" />
              <span>我了解此操作的风险并确认继续</span>
            </label>
          </div>
        </div>

        <div class="modal-footer">
          <button class="action-btn secondary" @click="showRestoreConfirm = false">
            取消
          </button>
          <button
            class="action-btn danger"
            @click="confirmRestore"
            :disabled="!restoreConfirm || restoringBackup"
          >
            <CogIcon v-if="restoringBackup" class="w-5 h-5 spinning" />
            {{ restoringBackup ? '恢复中...' : '确认恢复' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { backupApi } from '../api'
import {
  ArchiveBoxIcon,
  PlusIcon,
  ArrowPathIcon,
  CalendarIcon,
  CubeIcon,
  TagIcon,
  CheckCircleIcon,
  ClockIcon,
  XCircleIcon,
  ArrowDownTrayIcon,
  TrashIcon,
  ArrowUturnLeftIcon,
  XMarkIcon,
  ExclamationTriangleIcon,
  CogIcon
} from '@heroicons/vue/24/outline'

const loading = ref(false)
const backups = ref<any[]>([])
const showCreateBackup = ref(false)
const showRestoreConfirm = ref(false)
const selectedBackup = ref<any>(null)
const creatingBackup = ref(false)
const restoringBackup = ref(false)
const restoreConfirm = ref(false)

const backupForm = ref({
  name: '',
  type: 'full',
  description: '',
  includeDB: true,
  includeFiles: false,
  filePaths: ''
})

const refreshBackups = async () => {
  loading.value = true
  try {
    const response = await backupApi.getBackups() as any
    backups.value = response.backups || []
  } catch (error) {
    console.error('Failed to fetch backups:', error)
  } finally {
    loading.value = false
  }
}

const createBackup = async () => {
  if (!backupForm.value.name) {
    alert('请输入备份名称')
    return
  }

  creatingBackup.value = true
  try {
    const payload = {
      ...backupForm.value,
      filePaths: backupForm.value.filePaths ? backupForm.value.filePaths.split(',').map((p: string) => p.trim()) : []
    }

    await backupApi.createBackup(payload)
    showCreateBackup.value = false
    alert('备份创建成功！')
    refreshBackups()

    // 重置表单
    backupForm.value = {
      name: '',
      type: 'full',
      description: '',
      includeDB: true,
      includeFiles: false,
      filePaths: ''
    }
  } catch (error: any) {
    console.error('Failed to create backup:', error)
    alert('创建备份失败: ' + (error.message || '未知错误'))
  } finally {
    creatingBackup.value = false
  }
}

const downloadBackup = async (backup: any) => {
  try {
    const url = await backupApi.downloadBackup(backup.id)
    window.open(url, '_blank')
  } catch (error: any) {
    console.error('Failed to download backup:', error)
    alert('下载备份失败: ' + (error.message || '未知错误'))
  }
}

const confirmDeleteBackup = (backup: any) => {
  if (confirm(`确定要删除备份 "${backup.name}" 吗？此操作不可撤销。`)) {
    deleteBackup(backup)
  }
}

const deleteBackup = async (backup: any) => {
  try {
    await backupApi.deleteBackup(backup.id)
    alert('备份删除成功')
    refreshBackups()
  } catch (error: any) {
    console.error('Failed to delete backup:', error)
    alert('删除备份失败: ' + (error.message || '未知错误'))
  }
}

const restoreBackup = (backup: any) => {
  selectedBackup.value = backup
  showRestoreConfirm.value = true
  restoreConfirm.value = false
}

const confirmRestore = async () => {
  if (!selectedBackup.value) return

  restoringBackup.value = true
  try {
    await backupApi.restoreBackup({ backupId: selectedBackup.value.id })
    showRestoreConfirm.value = false
    alert('备份恢复成功！系统将重新启动。')
    // 可以考虑自动刷新页面或重启系统
  } catch (error: any) {
    console.error('Failed to restore backup:', error)
    alert('恢复备份失败: ' + (error.message || '未知错误'))
  } finally {
    restoringBackup.value = false
  }
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const formatSize = (bytes: number) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const size = Math.floor(Math.log(bytes) / Math.log(1024))
  return Math.round((bytes / Math.pow(1024, size)) * 100) / 100 + ' ' + units[size]
}

const getStatusText = (status: string) => {
  const statusMap: Record<string, string> = {
    completed: '完成',
    in_progress: '进行中',
    failed: '失败'
  }
  return statusMap[status] || status
}

onMounted(() => {
  refreshBackups()
})
</script>

<style scoped>
.backup-manager {
  width: 100%;
  height: 100%;
  padding: 32px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.backup-header {
  margin-bottom: 32px;
  background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%);
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(14, 165, 233, 0.2);
}

.backup-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.quick-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 8px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.primary {
  background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%);
  color: white;
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(14, 165, 233, 0.3);
}

.action-btn.secondary {
  background: white;
  color: #0369a1;
  border: 2px solid #0ea5e9;
}

.action-btn.secondary:hover {
  background: #f0f9ff;
}

.action-btn.danger {
  background: linear-gradient(135deg, #ef4444 0%, #dc2626 100%);
  color: white;
}

.action-btn.warning {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  color: white;
}

.backups-list {
  flex: 1;
}

.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #64748b;
}

.empty-state svg {
  color: #94a3b8;
  margin-bottom: 16px;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.backups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.backup-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 8px 24px rgba(14, 165, 233, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  display: flex;
  gap: 16px;
  transition: all 0.2s ease;
}

.backup-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 32px rgba(14, 165, 233, 0.25);
}

.backup-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  background: linear-gradient(135deg, #0ea5e9 0%, #0284c7 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.backup-info {
  flex: 1;
  min-width: 0;
}

.backup-info h3 {
  font-size: 16px;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 4px;
}

.backup-description {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 12px;
}

.backup-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #64748b;
}

.meta-item.status-completed {
  color: #10b981;
}

.meta-item.status-in_progress {
  color: #f59e0b;
}

.meta-item.status-failed {
  color: #ef4444;
}

.backup-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: #f1f5f9;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background: #e2e8f0;
  color: #0f172a;
}

.icon-btn.danger:hover {
  background: #fee2e2;
  color: #ef4444;
}

.icon-btn.warning:hover {
  background: #fef3c7;
  color: #f59e0b;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.modal-content.small {
  max-width: 500px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h2 {
  font-size: 20px;
  font-weight: 600;
  color: #0f172a;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: #f1f5f9;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #e2e8f0;
  color: #0f172a;
}

.modal-body {
  padding: 24px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px 24px;
  border-top: 1px solid #e5e7eb;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Form Styles */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-weight: 600;
  color: #0f172a;
  margin-bottom: 8px;
}

.form-input {
  width: 100%;
  padding: 10px 12px;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #0ea5e9;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.warning-message {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: #fef3c7;
  border-radius: 8px;
  color: #92400e;
  margin-bottom: 20px;
}

.warning-message svg {
  flex-shrink: 0;
}
</style>