<template>
  <div class="share-manager">
    <!-- 头部 -->
    <div class="share-header">
      <div class="header-left">
        <h1>共享文件夹管理</h1>
        <p class="subtitle">管理SMB共享文件夹和权限</p>
      </div>
      <div class="header-actions">
        <button class="action-btn primary" @click="showCreateShare = true">
          <PlusIcon class="w-4 h-4" />
          创建共享
        </button>
        <button class="action-btn" @click="refreshShares">
          <ArrowPathIcon class="w-4 h-4" />
          刷新
        </button>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="share-content">
      <!-- 共享文件夹列表 -->
      <div class="shares-section">
        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载共享文件夹...</p>
        </div>

        <div v-else-if="shares.length === 0" class="empty-state">
          <FolderIcon class="w-12 h-12" />
          <p>还没有共享文件夹</p>
          <button class="btn-primary" @click="showCreateShare = true">
            创建第一个共享
          </button>
        </div>

        <div v-else class="shares-grid">
          <div
            v-for="share in shares"
            :key="share.id"
            class="share-card"
            @click="viewShareDetails(share)"
          >
            <div class="share-icon">
              <FolderIcon class="w-8 h-8" />
            </div>

            <div class="share-info">
              <h3>{{ share.name }}</h3>
              <p class="share-path">{{ share.path }}</p>
              <p class="share-description">{{ share.description || '无描述' }}</p>
            </div>

            <div class="share-status">
              <span
                class="status-badge"
                :class="{ active: share.enabled, inactive: !share.enabled }"
              >
                {{ share.enabled ? '已启用' : '已禁用' }}
              </span>
            </div>

            <div class="share-meta">
              <span class="meta-item">
                <EyeIcon class="w-4 h-4" />
                {{ share.browseable ? '可浏览' : '不可浏览' }}
              </span>
              <span class="meta-item">
                <LockClosedIcon v-if="share.readOnly" class="w-4 h-4" />
                <LockOpenIcon v-else class="w-4 h-4" />
                {{ share.readOnly ? '只读' : '读写' }}
              </span>
            </div>

            <div class="share-actions">
              <button
                class="icon-btn"
                @click.stop="editShare(share)"
                title="编辑"
              >
                <PencilIcon class="w-4 h-4" />
              </button>
              <button
                class="icon-btn danger"
                @click.stop="deleteShare(share)"
                title="删除"
              >
                <TrashIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建共享对话框 -->
    <div v-if="showCreateShare" class="modal-overlay" @click="showCreateShare = false">
      <div class="modal-content large-modal" @click.stop>
        <div class="modal-header">
          <h3>创建共享文件夹</h3>
          <button class="close-btn" @click="showCreateShare = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="createShare" class="modal-body">
          <!-- 步骤指示器 -->
          <div class="steps-indicator">
            <div
              class="step"
              :class="{ active: currentStep === 1, completed: currentStep > 1 }"
            >
              <div class="step-number">1</div>
              <div class="step-label">基本信息</div>
            </div>
            <div class="step-line"></div>
            <div
              class="step"
              :class="{ active: currentStep === 2, completed: currentStep > 2 }"
            >
              <div class="step-number">2</div>
              <div class="step-label">权限设置</div>
            </div>
            <div class="step-line"></div>
            <div
              class="step"
              :class="{ active: currentStep === 3, completed: currentStep > 3 }"
            >
              <div class="step-number">3</div>
              <div class="step-label">SMB设置</div>
            </div>
          </div>

          <!-- 步骤1: 基本信息 -->
          <div v-if="currentStep === 1" class="step-content">
            <div class="form-group">
              <label>共享名称 *</label>
              <input
                type="text"
                v-model="shareForm.name"
                placeholder="输入共享名称"
                required
                @input="validateShareName"
              />
              <span v-if="nameError" class="error-message">{{ nameError }}</span>
            </div>

            <div class="form-group">
              <label>文件夹路径 *</label>
              <div class="path-input-group">
                <input
                  type="text"
                  v-model="shareForm.path"
                  placeholder="/mnt/data/share"
                  required
                />
                <button type="button" class="browse-btn" @click="browsePath">
                  <FolderIcon class="w-4 h-4" />
                  浏览
                </button>
              </div>
              <span v-if="pathError" class="error-message">{{ pathError }}</span>
            </div>

            <div class="form-group">
              <label>描述</label>
              <textarea
                v-model="shareForm.description"
                placeholder="输入共享描述（可选）"
                rows="3"
              ></textarea>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-secondary" @click="showCreateShare = false">
                取消
              </button>
              <button type="button" class="btn-primary" @click="nextStep">
                下一步
              </button>
            </div>
          </div>

          <!-- 步骤2: 权限设置 -->
          <div v-if="currentStep === 2" class="step-content">
            <div class="permission-matrix">
              <h4>权限矩阵</h4>

              <div class="matrix-table">
                <table>
                  <thead>
                    <tr>
                      <th>用户/组</th>
                      <th>权限级别</th>
                      <th>操作</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(perm, index) in shareForm.permissions" :key="index">
                      <td>
                        <select v-model="(perm as any).type" required>
                          <option value="user">用户</option>
                          <option value="group">组</option>
                        </select>
                        <input
                          type="text"
                          v-model="perm.name"
                          placeholder="用户名或组名"
                          required
                        />
                      </td>
                      <td>
                        <select v-model="perm.permission" required>
                          <option value="r">只读</option>
                          <option value="rw">读写</option>
                          <option value="full">完全控制</option>
                        </select>
                      </td>
                      <td>
                        <button
                          type="button"
                          class="icon-btn danger"
                          @click="removePermission(index)"
                        >
                          <TrashIcon class="w-4 h-4" />
                        </button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>

              <button type="button" class="add-perm-btn" @click="addPermission">
                <PlusIcon class="w-4 h-4" />
                添加权限
              </button>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-secondary" @click="prevStep">
                上一步
              </button>
              <button type="button" class="btn-primary" @click="nextStep">
                下一步
              </button>
            </div>
          </div>

          <!-- 步骤3: SMB设置 -->
          <div v-if="currentStep === 3" class="step-content">
            <div class="smb-settings">
              <div class="form-group">
                <label class="checkbox-label">
                  <input
                    type="checkbox"
                    v-model="shareForm.browseable"
                  />
                  <span>可浏览</span>
                  <small>允许用户浏览此共享</small>
                </label>
              </div>

              <div class="form-group">
                <label class="checkbox-label">
                  <input
                    type="checkbox"
                    v-model="shareForm.readOnly"
                  />
                  <span>只读</span>
                  <small>禁止用户写入文件</small>
                </label>
              </div>

              <div class="form-group">
                <label class="checkbox-label">
                  <input
                    type="checkbox"
                    v-model="shareForm.guestOK"
                  />
                  <span>允许访客访问</span>
                  <small>无需身份验证即可访问</small>
                </label>
              </div>

              <div class="warning-box">
                <ExclamationTriangleIcon class="w-5 h-5" />
                <div>
                  <h4>安全提示</h4>
                  <p>启用访客访问会降低安全性。建议仅在受信任的网络环境中使用。</p>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="btn-secondary" @click="prevStep">
                上一步
              </button>
              <button type="submit" class="btn-primary">
                创建共享
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>

    <!-- 编辑共享对话框 -->
    <div v-if="showEditShare" class="modal-overlay" @click="showEditShare = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑共享</h3>
          <button class="close-btn" @click="showEditShare = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="updateShare" class="modal-body">
          <div class="form-group">
            <label>共享名称</label>
            <input type="text" :value="selectedShare?.name" disabled />
          </div>

          <div class="form-group">
            <label>描述</label>
            <textarea
              v-model="editForm.description"
              placeholder="输入共享描述"
              rows="3"
            ></textarea>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="editForm.browseable" />
              <span>可浏览</span>
            </label>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="editForm.readOnly" />
              <span>只读</span>
            </label>
          </div>

          <div class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="editForm.guestOK" />
              <span>允许访客访问</span>
            </label>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showEditShare = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              保存更改
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 共享详情对话框 -->
    <div v-if="showShareDetails" class="modal-overlay" @click="showShareDetails = false">
      <div class="modal-content large-modal" @click.stop>
        <div class="modal-header">
          <h3>共享详情</h3>
          <button class="close-btn" @click="showShareDetails = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="modal-body">
          <div v-if="loadingDetails" class="loading-state">
            <div class="spinner"></div>
            <p>加载共享详情...</p>
          </div>

          <div v-else class="share-details">
            <div class="detail-section">
              <h4>基本信息</h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <label>共享名称</label>
                  <span>{{ selectedShare?.name }}</span>
                </div>
                <div class="detail-item">
                  <label>路径</label>
                  <span>{{ selectedShare?.path }}</span>
                </div>
                <div class="detail-item">
                  <label>状态</label>
                  <span
                    class="status-badge"
                    :class="{ active: selectedShare?.enabled, inactive: !selectedShare?.enabled }"
                  >
                    {{ selectedShare?.enabled ? '已启用' : '已禁用' }}
                  </span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>权限信息</h4>
              <div class="permission-details">
                <div v-if="sharePermissions?.length === 0" class="empty-permissions">
                  <p>没有设置特定权限</p>
                </div>
                <div v-else class="permissions-list">
                  <div
                    v-for="(perm, index) in sharePermissions"
                    :key="index"
                    class="permission-item"
                  >
                    <span class="perm-type">{{ perm.type === 'user' ? '用户' : '组' }}</span>
                    <span class="perm-name">{{ perm.name }}</span>
                    <span class="perm-level">{{ getPermissionLabel(perm.permission) }}</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>文件权限</h4>
              <div v-if="fileInfo" class="file-permissions">
                <div class="perm-row">
                  <label>所有者</label>
                  <span>{{ fileInfo.owner }}</span>
                </div>
                <div class="perm-row">
                  <label>组</label>
                  <span>{{ fileInfo.group }}</span>
                </div>
                <div class="perm-row">
                  <label>权限模式</label>
                  <span class="perm-mode">{{ fileInfo.mode }}</span>
                </div>
                <div class="perm-row">
                  <label>大小</label>
                  <span>{{ formatFileSize(fileInfo.size) }}</span>
                </div>
                <div class="perm-row">
                  <label>修改时间</label>
                  <span>{{ fileInfo.modified }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>ACL规则</h4>
              <div v-if="aclRules?.length === 0" class="empty-acl">
                <p>没有设置ACL规则</p>
              </div>
              <div v-else class="acl-rules">
                <div
                  v-for="(rule, index) in aclRules"
                  :key="index"
                  class="acl-rule"
                >
                  <span class="rule-type">{{ rule.type }}</span>
                  <span class="rule-name">{{ rule.name }}</span>
                  <span class="rule-perms">{{ rule.perms }}</span>
                  <span v-if="rule.default" class="rule-default">默认</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  FolderIcon,
  PlusIcon,
  ArrowPathIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  EyeIcon,
  LockClosedIcon,
  LockOpenIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'

// 状态
const loading = ref(false)
const loadingDetails = ref(false)
const shares = ref<any[]>([])
const selectedShare = ref<any>(null)
const sharePermissions = ref<any[]>([])
const fileInfo = ref<any>(null)
const aclRules = ref<any[]>([])

// 对话框状态
const showCreateShare = ref(false)
const showEditShare = ref(false)
const showShareDetails = ref(false)
const currentStep = ref(1)

// 表单数据
const shareForm = reactive({
  name: '',
  path: '',
  description: '',
  browseable: true,
  readOnly: false,
  guestOK: false,
  permissions: []
})

const editForm = reactive({
  description: '',
  browseable: true,
  readOnly: false,
  guestOK: false
})

// 验证错误
const nameError = ref('')
const pathError = ref('')

// 方法
const loadShares = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/permissions/shares', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      shares.value = data.shares || []
    } else {
      throw new Error('获取共享文件夹失败')
    }
  } catch (error) {
    console.error('Failed to load shares:', error)
    ElMessage.error('获取共享文件夹失败')
  } finally {
    loading.value = false
  }
}

const validateShareName = () => {
  if (!shareForm.name) {
    nameError.value = '请输入共享名称'
    return false
  }

  if (!/^[a-zA-Z0-9_-]+$/.test(shareForm.name)) {
    nameError.value = '共享名称只能包含字母、数字、下划线和连字符'
    return false
  }

  // 检查名称是否已存在
  const exists = shares.value.some(s => s.name === shareForm.name)
  if (exists) {
    nameError.value = '共享名称已存在'
    return false
  }

  nameError.value = ''
  return true
}

const browsePath = () => {
  // 简化实现：让用户直接输入路径
  ElMessage.info('请输入完整的文件夹路径，例如：/mnt/data/share')
}

const nextStep = () => {
  if (currentStep.value === 1) {
    if (!validateShareName() || !shareForm.path) {
      ElMessage.warning('请填写完整的基本信息')
      return
    }
  }

  if (currentStep.value < 3) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 1) {
    currentStep.value--
  }
}

const addPermission = () => {
  shareForm.permissions.push({
    type: 'user',
    name: '',
    permission: 'rw'
  })
}

const removePermission = (index: number) => {
  shareForm.permissions.splice(index, 1)
}

const createShare = async () => {
  try {
    const response = await fetch('/api/permissions/shares', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(shareForm)
    })

    if (response.ok) {
      ElMessage.success('共享文件夹创建成功')
      showCreateShare.value = false
      currentStep.value = 1
      resetForm()
      await loadShares()
    } else {
      throw new Error('创建共享文件夹失败')
    }
  } catch (error) {
    console.error('Failed to create share:', error)
    ElMessage.error('创建共享文件夹失败')
  }
}

const viewShareDetails = async (share: any) => {
  selectedShare.value = share
  loadingDetails.value = true
  showShareDetails.value = true

  try {
    // 获取共享权限
    const permResponse = await fetch(`/api/permissions/shares/${share.name}/permissions`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (permResponse.ok) {
      const permData = await permResponse.json()
      sharePermissions.value = permData.permissions || []
      fileInfo.value = permData.fileInfo
    }

    // 获取ACL
    const aclResponse = await fetch(`/api/permissions/files/acl?path=${encodeURIComponent(share.path)}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (aclResponse.ok) {
      const aclData = await aclResponse.json()
      aclRules.value = aclData.acl || []
    }
  } catch (error) {
    console.error('Failed to load share details:', error)
  } finally {
    loadingDetails.value = false
  }
}

const editShare = (share: any) => {
  selectedShare.value = share
  editForm.description = share.description || ''
  editForm.browseable = share.browseable
  editForm.readOnly = share.readOnly
  editForm.guestOK = share.guestOK
  showEditShare.value = true
}

const updateShare = async () => {
  try {
    const response = await fetch(`/api/permissions/shares/${selectedShare.value.name}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: selectedShare.value.path,
        description: editForm.description,
        browseable: editForm.browseable,
        readOnly: editForm.readOnly,
        guestOK: editForm.guestOK,
        permissions: selectedShare.value.permissions
      })
    })

    if (response.ok) {
      ElMessage.success('共享更新成功')
      showEditShare.value = false
      await loadShares()
    } else {
      throw new Error('更新共享失败')
    }
  } catch (error) {
    console.error('Failed to update share:', error)
    ElMessage.error('更新共享失败')
  }
}

const deleteShare = (share: any) => {
  if (!confirm(`确定要删除共享文件夹"${share.name}"吗？`)) {
    return
  }

  fetch(`/api/permissions/shares/${share.name}`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  })
  .then(response => {
    if (response.ok) {
      ElMessage.success('共享删除成功')
      return loadShares()
    }
    throw new Error('删除共享失败')
  })
  .catch(error => {
    console.error('Failed to delete share:', error)
    ElMessage.error('删除共享失败')
  })
}

const refreshShares = () => {
  loadShares()
}

const resetForm = () => {
  shareForm.name = ''
  shareForm.path = ''
  shareForm.description = ''
  shareForm.browseable = true
  shareForm.readOnly = false
  shareForm.guestOK = false
  shareForm.permissions = []
  nameError.value = ''
  pathError.value = ''
}

const getPermissionLabel = (perm: string) => {
  switch (perm) {
    case 'r': return '只读'
    case 'rw': return '读写'
    case 'full': return '完全控制'
    default: return perm
  }
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

// 生命周期
onMounted(() => {
  loadShares()
})
</script>

<style scoped lang="scss">
.share-manager {
  width: 100%;
  padding: 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.share-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
}

.header-left h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin: 0 0 4px 0;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  &.primary {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  &:hover:not(.primary) {
    background: rgba(102, 126, 234, 0.1);
  }
}

.share-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.shares-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #6b7280;

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid rgba(102, 126, 234, 0.2);
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 12px;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  color: #9ca3af;

  svg {
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    font-size: 16px;
    margin-bottom: 20px;
  }

  .btn-primary {
    padding: 12px 24px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      opacity: 0.9;
    }
  }
}

.shares-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.share-card {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  &:hover {
    border-color: #667eea;
    box-shadow: 0 4px 12px rgba(102, 126, 234, 0.15);
    transform: translateY(-2px);
  }
}

.share-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.share-info {
  margin-bottom: 16px;

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 8px 0;
  }

  .share-path {
    font-size: 13px;
    color: #6b7280;
    font-family: monospace;
    margin: 4px 0;
  }

  .share-description {
    font-size: 14px;
    color: #9ca3af;
    margin: 8px 0 0 0;
  }
}

.share-status {
  margin-bottom: 12px;
}

.status-badge {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;

  &.active {
    background: rgba(34, 197, 94, 0.1);
    color: #066;
  }

  &.inactive {
    background: rgba(239, 68, 68, 0.1);
    color: #991b1b;
  }
}

.share-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;

  .meta-item {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: #6b7280;
  }
}

.share-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;

  .share-card:hover & {
    opacity: 1;
  }
}

.icon-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    background: rgba(102, 126, 234, 0.1);
  }

  &.danger {
    color: #ef4444;

    &:hover {
      background: rgba(239, 68, 68, 0.1);
    }
  }
}

// 模态对话框
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
  min-width: 400px;
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);

  &.large-modal {
    max-width: 900px;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 6px;
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    background: rgba(102, 126, 234, 0.1);
  }
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 32px;
  padding: 20px 0;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;

  .step-number {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: #e5e7eb;
    color: #6b7280;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 600;
    font-size: 14px;
  }

  .step-label {
    font-size: 12px;
    color: #6b7280;
  }

  &.active .step-number {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
  }

  &.completed .step-number {
    background: #10b981;
    color: white;
  }
}

.step-line {
  width: 60px;
  height: 2px;
  background: #e5e7eb;
  margin: 0 8px 24px 8px;
}

.step-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
  }

  input, textarea, select {
    padding: 10px 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    transition: all 0.2s;

    &:focus {
      outline: none;
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }

    &:disabled {
      background: #f9fafb;
      color: #9ca3af;
    }
  }

  textarea {
    resize: vertical;
    min-height: 80px;
  }

  .error-message {
    color: #ef4444;
    font-size: 12px;
  }
}

.checkbox-label {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: rgba(102, 126, 234, 0.05);
    border-color: #667eea;
  }

  input[type="checkbox"] {
    margin-top: 2px;
  }

  span {
    font-weight: 500;
    color: #1f2937;
    display: block;
    margin-bottom: 4px;
  }

  small {
    font-size: 12px;
    color: #6b7280;
    line-height: 1.4;
  }
}

.path-input-group {
  display: flex;
  gap: 8px;

  input {
    flex: 1;
  }

  .browse-btn {
    padding: 10px 16px;
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
    border: 1px solid rgba(102, 126, 234, 0.2);
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    gap: 6px;

    &:hover {
      background: rgba(102, 126, 234, 0.2);
    }
  }
}

.permission-matrix {
  display: flex;
  flex-direction: column;
  gap: 16px;

  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.matrix-table {
  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: linear-gradient(to bottom, #f8fafc, #e2e8f0);

    th {
      padding: 12px;
      text-align: left;
      font-size: 13px;
      font-weight: 600;
      color: #374151;
      border-bottom: 1px solid #e5e7eb;
    }
  }

  tbody {
    tr {
      border-bottom: 1px solid #f3f4f6;

      td {
        padding: 12px;
        display: flex;
        gap: 8px;
        align-items: center;
      }

      select, input {
        padding: 8px;
        border: 1px solid #e5e7eb;
        border-radius: 6px;
        font-size: 13px;
      }

      select {
        width: 100px;
      }

      input {
        flex: 1;
      }
    }
  }
}

.add-perm-btn {
  padding: 10px 16px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border: 1px dashed rgba(102, 126, 234, 0.3);
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  &:hover {
    background: rgba(102, 126, 234, 0.2);
    border-color: #667eea;
  }
}

.smb-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.warning-box {
  display: flex;
  gap: 12px;
  padding: 16px;
  background: rgba(251, 191, 36, 0.1);
  border: 1px solid rgba(251, 191, 36, 0.3);
  border-radius: 8px;
  color: #b45309;

  svg {
    flex-shrink: 0;
    margin-top: 2px;
  }

  h4 {
    font-size: 14px;
    font-weight: 600;
    margin: 0 0 4px 0;
  }

  p {
    font-size: 13px;
    margin: 0;
    line-height: 1.4;
  }
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 20px;
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 10px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;

  &:hover {
    opacity: 0.9;
  }
}

.btn-secondary {
  background: white;
  border: 1px solid #e5e7eb;
  color: #6b7280;

  &:hover {
    background: #f9fafb;
    border-color: #d1d5db;
  }
}

.share-details {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-section {
  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 16px 0;
  }
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  label {
    font-size: 12px;
    font-weight: 500;
    color: #6b7280;
  }

  span {
    font-size: 14px;
    color: #1f2937;
  }
}

.permission-details {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .empty-permissions {
    padding: 20px;
    background: #f9fafb;
    border-radius: 8px;
    text-align: center;
    color: #9ca3af;
    font-size: 14px;
  }
}

.permissions-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;

  .perm-type {
    padding: 4px 8px;
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  }

  .perm-name {
    flex: 1;
    font-weight: 500;
    color: #1f2937;
  }

  .perm-level {
    padding: 4px 12px;
    background: rgba(16, 185, 129, 0.1);
    color: #065f46;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  }
}

.file-permissions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.perm-row {
  display: flex;
  justify-content: space-between;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;

  label {
    font-size: 12px;
    font-weight: 500;
    color: #6b7280;
  }

  span {
    font-size: 14px;
    color: #1f2937;

    &.perm-mode {
      font-family: monospace;
      font-weight: 600;
    }
  }
}

.empty-acl {
  padding: 20px;
  background: #f9fafb;
  border-radius: 8px;
  text-align: center;
  color: #9ca3af;
  font-size: 14px;
}

.acl-rules {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.acl-rule {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;

  .rule-type {
    padding: 4px 8px;
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  }

  .rule-name {
    flex: 1;
    font-weight: 500;
    color: #1f2937;
  }

  .rule-perms {
    padding: 4px 12px;
    background: rgba(16, 185, 129, 0.1);
    color: #065f46;
    border-radius: 4px;
    font-size: 12px;
    font-family: monospace;
    font-weight: 600;
  }

  .rule-default {
    padding: 4px 8px;
    background: rgba(251, 191, 36, 0.1);
    color: #b45309;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  }
}
</style>