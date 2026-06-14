<template>
  <div class="smb-manager">
    <!-- SMB 服务状态 -->
    <div class="smb-service-card">
      <div class="service-header">
        <div class="service-info">
          <div class="service-icon">
            <ServerIcon class="w-6 h-6" />
          </div>
          <div>
            <h3>SMB 文件共享</h3>
            <p class="service-description">让其他设备访问您的文件</p>
          </div>
        </div>
        <div class="service-control">
          <div class="toggle-switch">
            <input
              id="smb-service"
              type="checkbox"
              v-model="smbService.enabled"
              @change="toggleSMBService"
            />
            <label for="smb-service" class="toggle-label">
              <span class="toggle-slider"></span>
            </label>
          </div>
          <span class="service-status" :class="{ active: smbService.enabled }">
            {{ smbService.enabled ? '已启用' : '已禁用' }}
          </span>
        </div>
      </div>
      <p class="service-help">
        <InformationCircleIcon class="w-4 h-4" />
        启用后，您可以通过网络在Windows、Mac等设备上访问NAS的文件
      </p>
    </div>

    <!-- 快速共享设置 -->
    <div v-if="smbService.enabled" class="smb-shares-section">
      <div class="section-header">
        <h3>共享文件夹</h3>
        <button class="add-btn" @click="showAddShareModal = true">
          <PlusIcon class="w-4 h-4" />
          添加共享
        </button>
      </div>

      <!-- 共享列表 -->
      <div v-if="shares.length === 0" class="empty-state">
        <FolderIcon class="w-12 h-12" />
        <p>还没有共享文件夹</p>
        <button class="btn-primary" @click="showAddShareModal = true">
          添加第一个共享
        </button>
      </div>

      <div v-else class="shares-list">
        <div
          v-for="share in shares"
          :key="share.id"
          class="share-item"
        >
          <div class="share-icon">
            <FolderIcon class="w-6 h-6" />
          </div>
          <div class="share-details">
            <div class="share-name-row">
              <h4>{{ share.name }}</h4>
              <span class="share-status" :class="{ active: share.enabled }">
                {{ share.enabled ? '已启用' : '已禁用' }}
              </span>
            </div>
            <p class="share-path">{{ share.path }}</p>
            <p class="share-description">{{ share.description || '无描述' }}</p>
            <div class="share-access">
              <span class="access-badge">
                <LockOpenIcon v-if="!share.readOnly" class="w-3 h-3" />
                <LockClosedIcon v-else class="w-3 h-3" />
                {{ share.readOnly ? '只读' : '读写' }}
              </span>
              <span class="access-badge">
                <EyeIcon class="w-3 h-3" />
                {{ share.browseable ? '可见' : '隐藏' }}
              </span>
              <span class="access-badge">
                <UserGroupIcon class="w-3 h-3" />
                {{ share.guestOk ? '允许访客' : '需要登录' }}
              </span>
            </div>
          </div>
          <div class="share-actions">
            <button
              class="icon-btn"
              @click="editShare(share)"
              title="编辑"
            >
              <PencilIcon class="w-4 h-4" />
            </button>
            <button
              class="icon-btn"
              @click="toggleShareStatus(share)"
              :title="share.enabled ? '禁用' : '启用'"
            >
              <PlayIcon v-if="!share.enabled" class="w-4 h-4" />
              <NoSymbolIcon v-else class="w-4 h-4" />
            </button>
            <button
              class="icon-btn danger"
              @click="deleteShare(share)"
              title="删除"
            >
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

      <!-- 访问地址 -->
      <div class="access-info">
        <div class="info-header">
          <InformationCircleIcon class="w-5 h-5" />
          <h4>如何访问共享文件夹</h4>
        </div>
        <div class="access-steps">
          <div class="step">
            <div class="step-number">1</div>
            <div class="step-content">
              <h5>Windows 电脑</h5>
              <p>在文件资源管理器地址栏输入：<code>\\{{ ipAddress }}</code></p>
            </div>
          </div>
          <div class="step">
            <div class="step-number">2</div>
            <div class="step-content">
              <h5>Mac 电脑</h5>
              <p>在 Finder 前往 > 连接服务器，输入：<code>smb://{{ ipAddress }}</code></p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 添加/编辑共享对话框 -->
    <div v-if="showAddShareModal" class="modal-overlay" @click.self="closeShareModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>{{ editingShare ? '编辑共享' : '添加共享' }}</h3>
          <button class="close-btn" @click="closeShareModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <form @submit.prevent="saveShare" class="share-form">
          <div class="form-group">
            <label>共享名称 *</label>
            <input
              v-model="shareForm.name"
              type="text"
              required
              placeholder="例如：文档、照片、备份"
            />
            <small>这是其他设备看到的文件夹名称</small>
          </div>
          <div class="form-group">
            <label>文件夹路径 *</label>
            <input
              v-model="shareForm.path"
              type="text"
              required
              placeholder="/mnt/storage/documents"
            />
            <small>要共享的文件夹完整路径</small>
          </div>
          <div class="form-group">
            <label>描述</label>
            <input
              v-model="shareForm.description"
              type="text"
              placeholder="这个共享的用途说明"
            />
          </div>
          <div class="form-group">
            <label>访问权限</label>
            <div class="permission-options">
              <label class="permission-option">
                <input type="radio" v-model="shareForm.readOnly" :value="false" />
                <div class="option-content">
                  <LockOpenIcon class="w-5 h-5" />
                  <div>
                    <span class="option-title">读写</span>
                    <span class="option-desc">用户可以查看和修改文件</span>
                  </div>
                </div>
              </label>
              <label class="permission-option">
                <input type="radio" v-model="shareForm.readOnly" :value="true" />
                <div class="option-content">
                  <LockClosedIcon class="w-5 h-5" />
                  <div>
                    <span class="option-title">只读</span>
                    <span class="option-desc">用户只能查看，不能修改</span>
                  </div>
                </div>
              </label>
            </div>
          </div>
          <div class="form-group">
            <label>其他选项</label>
            <div class="checkbox-options">
              <label class="checkbox-option">
                <input type="checkbox" v-model="shareForm.browseable" />
                <span>允许浏览（其他设备能看到此共享）</span>
              </label>
              <label class="checkbox-option">
                <input type="checkbox" v-model="shareForm.guestOk" />
                <span>允许访客访问（无需密码）</span>
              </label>
              <label class="checkbox-option">
                <input type="checkbox" v-model="shareForm.enabled" />
                <span>立即启用此共享</span>
              </label>
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="closeShareModal">
              取消
            </button>
            <button type="submit" class="btn-primary">
              {{ editingShare ? '保存修改' : '添加共享' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  ServerIcon,
  FolderIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
  PlayIcon,
  NoSymbolIcon,
  LockOpenIcon,
  LockClosedIcon,
  EyeIcon,
  UserGroupIcon,
  InformationCircleIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'

// SMB 服务状态
const smbService = ref({
  enabled: false,
  status: 'stopped',
  ipAddress: '192.168.1.100'
})

// 共享列表
const shares = ref<any[]>([])
const showAddShareModal = ref(false)
const editingShare = ref<any>(null)

// 共享表单
const shareForm = ref({
  name: '',
  path: '',
  description: '',
  readOnly: false,
  browseable: true,
  guestOk: false,
  enabled: true
})

// 获取IP地址
const ipAddress = ref('192.168.1.100')

// 切换SMB服务
const toggleSMBService = async () => {
  try {
    if (smbService.value.enabled) {
      // 启用SMB服务
      await enableSMBService()
      smbService.value.status = 'running'
    } else {
      // 禁用SMB服务
      await disableSMBService()
      smbService.value.status = 'stopped'
    }
  } catch (error: any) {
    console.error('Failed to toggle SMB service:', error)
    // 回滚状态
    smbService.value.enabled = !smbService.value.enabled
    alert('操作失败: ' + error.message)
  }
}

// 启用SMB服务
const enableSMBService = async () => {
  // Mock API call
  console.log('Enabling SMB service...')
}

// 禁用SMB服务
const disableSMBService = async () => {
  // Mock API call
  console.log('Disabling SMB service...')
}

// 加载共享列表
const loadShares = async () => {
  try {
    // Mock data
    shares.value = [
      {
        id: '1',
        name: '文档',
        path: '/mnt/storage/documents',
        description: '共享文档文件夹',
        readOnly: false,
        browseable: true,
        guestOk: false,
        enabled: true
      },
      {
        id: '2',
        name: '照片',
        path: '/mnt/storage/photos',
        description: '家庭照片备份',
        readOnly: true,
        browseable: true,
        guestOk: true,
        enabled: true
      },
      {
        id: '3',
        name: '软件备份',
        path: '/mnt/storage/backup',
        description: '软件安装包备份',
        readOnly: true,
        browseable: false,
        guestOk: false,
        enabled: false
      }
    ]
  } catch (error: any) {
    console.error('Failed to load shares:', error)
  }
}

// 编辑共享
const editShare = (share: any) => {
  editingShare.value = share
  shareForm.value = { ...share }
  showAddShareModal.value = true
}

// 切换共享状态
const toggleShareStatus = async (share: any) => {
  try {
    share.enabled = !share.enabled
    console.log('Toggling share status:', share.name, share.enabled)
  } catch (error: any) {
    console.error('Failed to toggle share status:', error)
    alert('操作失败: ' + error.message)
  }
}

// 删除共享
const deleteShare = async (share: any) => {
  if (confirm(`确定要删除共享 "${share.name}" 吗？`)) {
    try {
      shares.value = shares.value.filter(s => s.id !== share.id)
      console.log('Share deleted:', share.name)
    } catch (error: any) {
      console.error('Failed to delete share:', error)
      alert('删除失败: ' + error.message)
    }
  }
}

// 保存共享
const saveShare = async () => {
  try {
    if (editingShare.value) {
      // 更新现有共享
      const index = shares.value.findIndex(s => s.id === editingShare.value.id)
      if (index !== -1) {
        shares.value[index] = {
          ...editingShare.value,
          ...shareForm.value
        }
      }
      console.log('Share updated:', shareForm.value.name)
    } else {
      // 添加新共享
      const newShare = {
        id: Date.now().toString(),
        ...shareForm.value
      }
      shares.value.push(newShare)
      console.log('Share added:', shareForm.value.name)
    }
    closeShareModal()
  } catch (error: any) {
    console.error('Failed to save share:', error)
    alert('保存失败: ' + error.message)
  }
}

// 关闭对话框
const closeShareModal = () => {
  showAddShareModal.value = false
  editingShare.value = null
  shareForm.value = {
    name: '',
    path: '',
    description: '',
    readOnly: false,
    browseable: true,
    guestOk: false,
    enabled: true
  }
}

// 获取系统IP地址
const getSystemIP = async () => {
  try {
    // Mock implementation - 应该调用实际的API获取IP地址
    ipAddress.value = '192.168.1.100'
  } catch (error) {
    console.error('Failed to get system IP:', error)
  }
}

onMounted(() => {
  loadShares()
  getSystemIP()
})
</script>

<style scoped>
.smb-manager {
  max-width: 900px;
  margin: 0 auto;
}

/* SMB 服务卡片 */
.smb-service-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 24px;
  color: white;
  margin-bottom: 24px;
}

.service-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.service-info {
  display: flex;
  gap: 16px;
  align-items: center;
}

.service-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-info h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 4px;
}

.service-description {
  font-size: 14px;
  opacity: 0.9;
}

.service-control {
  display: flex;
  align-items: center;
  gap: 16px;
}

.toggle-switch {
  position: relative;
  width: 60px;
  height: 32px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-label {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.3);
  transition: 0.3s;
  border-radius: 32px;
}

.toggle-slider {
  position: absolute;
  content: "";
  height: 26px;
  width: 26px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

.toggle-switch input:checked + .toggle-label {
  background-color: rgba(16, 185, 129, 0.8);
}

.toggle-switch input:checked + .toggle-label .toggle-slider {
  transform: translateX(28px);
}

.service-status {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
}

.service-status.active {
  background: rgba(16, 185, 129, 0.8);
}

.service-help {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  opacity: 0.9;
}

/* 共享区域 */
.smb-shares-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.add-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.add-btn:hover {
  background: #2563eb;
}

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 48px 24px;
  color: #6b7280;
}

.empty-state svg {
  margin: 0 auto 16px;
  color: #9ca3af;
}

.empty-state p {
  margin-bottom: 24px;
}

.btn-primary {
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #2563eb;
}

/* 共享列表 */
.shares-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 24px;
}

.share-item {
  display: flex;
  gap: 16px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  transition: box-shadow 0.2s;
}

.share-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.share-icon {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.share-details {
  flex: 1;
}

.share-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.share-name-row h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.share-status {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  background: #f3f4f6;
  color: #6b7280;
}

.share-status.active {
  background: #d1fae5;
  color: #10b981;
}

.share-path {
  font-size: 13px;
  color: #6b7280;
  font-family: monospace;
  margin-bottom: 4px;
}

.share-description {
  font-size: 14px;
  color: #374151;
  margin-bottom: 8px;
}

.share-access {
  display: flex;
  gap: 8px;
}

.access-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: white;
  border-radius: 6px;
  font-size: 12px;
  color: #6b7280;
  border: 1px solid #e5e7eb;
}

.share-actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: white;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.icon-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.icon-btn.danger:hover {
  background: #fee2e2;
  color: #ef4444;
}

/* 访问信息 */
.access-info {
  margin-top: 24px;
  padding: 20px;
  background: #eff6ff;
  border-radius: 12px;
  border: 1px solid #bfdbfe;
}

.info-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  color: #1e40af;
}

.info-header h4 {
  font-size: 16px;
  font-weight: 600;
}

.access-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step {
  display: flex;
  gap: 12px;
}

.step-number {
  width: 24px;
  height: 24px;
  background: #3b82f6;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 14px;
  flex-shrink: 0;
}

.step-content h5 {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 4px;
  color: #1e3a8a;
}

.step-content p {
  font-size: 13px;
  color: #1e40af;
}

.step-content code {
  background: rgba(255, 255, 255, 0.8);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
}

/* 对话框 */
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
}

.modal-content {
  background: white;
  border-radius: 16px;
  padding: 24px;
  max-width: 500px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: #f3f4f6;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: #e5e7eb;
}

/* 表单 */
.share-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input[type="text"] {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.form-group small {
  font-size: 12px;
  color: #6b7280;
}

/* 权限选项 */
.permission-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-option {
  display: flex;
  gap: 12px;
  padding: 12px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.permission-option:hover {
  border-color: #d1d5db;
}

.permission-option input:checked + .option-content {
  color: #3b82f6;
}

.option-content {
  display: flex;
  gap: 12px;
  align-items: center;
}

.option-content svg {
  flex-shrink: 0;
}

.option-title {
  display: block;
  font-weight: 500;
  margin-bottom: 2px;
}

.option-desc {
  display: block;
  font-size: 13px;
  color: #6b7280;
}

/* 复选框选项 */
.checkbox-options {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.checkbox-option {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
}

/* 操作按钮 */
.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;
}

.btn-secondary {
  padding: 10px 20px;
  background: white;
  color: #374151;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #f9fafb;
}

.btn-primary {
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #2563eb;
}
</style>