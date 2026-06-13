<template>
  <div class="group-quota-manager">
    <!-- 头部 -->
    <div class="quota-header">
      <div class="header-left">
        <h1>用户组配额管理</h1>
        <p class="subtitle">管理用户组的存储空间配额</p>
      </div>
      <div class="header-actions">
        <button class="action-btn primary" @click="showSetQuota = true">
          <PlusIcon class="w-4 h-4" />
          设置配额
        </button>
        <button class="action-btn" @click="refreshQuotas">
          <ArrowPathIcon class="w-4 h-4" />
          刷新
        </button>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="quota-content">
      <!-- 配额概览 -->
      <div class="quota-overview">
        <div class="overview-card">
          <div class="card-icon total">
            <UsersIcon class="w-8 h-8" />
          </div>
          <div class="card-info">
            <h3>总组数</h3>
            <p class="card-value">{{ groupQuotas.length }}</p>
          </div>
        </div>

        <div class="overview-card">
          <div class="card-icon limited">
            <CircleStackIcon class="w-8 h-8" />
          </div>
          <div class="card-info">
            <h3>已限制配额</h3>
            <p class="card-value">{{ limitedGroups }}</p>
          </div>
        </div>

        <div class="overview-card">
          <div class="card-icon usage">
            <ChartBarIcon class="w-8 h-8" />
          </div>
          <div class="card-info">
            <h3>总使用量</h3>
            <p class="card-value">{{ formatBytes(totalUsage) }}</p>
          </div>
        </div>

        <div class="overview-card">
          <div class="card-icon allocated">
            <CheckBadgeIcon class="w-8 h-8" />
          </div>
          <div class="card-info">
            <h3>总分配量</h3>
            <p class="card-value">{{ formatBytes(totalAllocated) }}</p>
          </div>
        </div>
      </div>

      <!-- 配额表格 -->
      <div class="quota-table-section">
        <div class="section-header">
          <h2>组配额列表</h2>
          <div class="table-controls">
            <input
              type="text"
              v-model="searchQuery"
              placeholder="搜索用户组..."
              class="search-input"
            />
            <select v-model="filterStatus" class="filter-select">
              <option value="all">所有状态</option>
              <option value="normal">正常</option>
              <option value="warning">警告</option>
              <option value="exceeded">超限</option>
            </select>
          </div>
        </div>

        <div v-if="loading" class="loading-state">
          <div class="spinner"></div>
          <p>加载配额信息...</p>
        </div>

        <div v-else-if="filteredQuotas.length === 0" class="empty-state">
          <CircleStackIcon class="w-12 h-12" />
          <p>没有找到配额信息</p>
        </div>

        <div v-else class="quota-table">
          <table>
            <thead>
              <tr>
                <th>组名</th>
                <th>GID</th>
                <th>成员数</th>
                <th>使用量</th>
                <th>软限制</th>
                <th>硬限制</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="quota in filteredQuotas" :key="quota.name">
                <td>
                  <div class="group-info">
                    <div class="group-avatar">{{ quota.name.charAt(0).toUpperCase() }}</div>
                    <div class="group-details">
                      <span class="group-name">{{ quota.name }}</span>
                    </div>
                  </div>
                </td>
                <td>{{ quota.gid }}</td>
                <td>{{ quota.memberCount || 0 }}</td>
                <td>
                  <div class="usage-info">
                    <span class="usage-value">{{ formatBytes(quota.usedSpace) }}</span>
                    <div class="usage-bar">
                      <div
                        class="usage-fill"
                        :style="{ width: getUsagePercent(quota) + '%' }"
                        :class="getUsageStatus(quota)"
                      ></div>
                    </div>
                  </div>
                </td>
                <td>
                  <span v-if="quota.softLimit" class="limit-value">
                    {{ formatBytes(quota.softLimit) }}
                  </span>
                  <span v-else class="limit-value unlimited">无限制</span>
                </td>
                <td>
                  <span v-if="quota.hardLimit" class="limit-value">
                    {{ formatBytes(quota.hardLimit) }}
                  </span>
                  <span v-else class="limit-value unlimited">无限制</span>
                </td>
                <td>
                  <span
                    class="status-badge"
                    :class="getUsageStatus(quota)"
                  >
                    {{ getUsageStatusLabel(quota) }}
                  </span>
                </td>
                <td>
                  <div class="action-buttons">
                    <button
                      class="icon-btn"
                      @click="editQuota(quota)"
                      title="编辑配额"
                    >
                      <PencilIcon class="w-4 h-4" />
                    </button>
                    <button
                      class="icon-btn"
                      @click="viewQuotaDetails(quota)"
                      title="查看详情"
                    >
                      <EyeIcon class="w-4 h-4" />
                    </button>
                    <button
                      class="icon-btn danger"
                      @click="removeQuota(quota)"
                      title="移除配额"
                    >
                      <TrashIcon class="w-4 h-4" />
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 配额使用趋势 -->
      <div class="quota-trends">
        <h2>配额使用趋势</h2>
        <div class="trends-chart">
          <div
            v-for="quota in topUsageGroups"
            :key="quota.name"
            class="trend-item"
          >
            <div class="trend-header">
              <span class="trend-name">{{ quota.name }}</span>
              <span class="trend-usage">{{ formatBytes(quota.usedSpace) }} / {{ formatBytes(quota.hardLimit || quota.softLimit) }}</span>
            </div>
            <div class="trend-bar">
              <div
                class="trend-fill"
                :style="{ width: getUsagePercent(quota) + '%' }"
                :class="getUsageStatus(quota)"
              ></div>
            </div>
            <div class="trend-percent">{{ getUsagePercent(quota).toFixed(1) }}%</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 设置配额对话框 -->
    <div v-if="showSetQuota" class="modal-overlay" @click="showSetQuota = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>设置组配额</h3>
          <button class="close-btn" @click="showSetQuota = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="setQuota" class="modal-body">
          <div class="form-group">
            <label>用户组</label>
            <select v-model="quotaForm.groupName" required>
              <option value="">选择用户组</option>
              <option
                v-for="group in availableGroups"
                :key="group.name"
                :value="group.name"
              >
                {{ group.name }} ({{ group.memberCount }} 成员)
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>路径</label>
            <input
              type="text"
              v-model="quotaForm.path"
              placeholder="/home 或 /mnt/data"
              required
            />
            <small>配额将应用于此路径及其子目录</small>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>软限制 (MB)</label>
              <input
                type="number"
                v-model.number="quotaForm.softLimit"
                min="0"
                placeholder="0"
              />
              <small>用户可临时超过此限制</small>
            </div>

            <div class="form-group">
              <label>硬限制 (MB)</label>
              <input
                type="number"
                v-model.number="quotaForm.hardLimit"
                min="0"
                placeholder="0"
              />
              <small>绝对限制，无法超过</small>
            </div>
          </div>

          <div class="form-group">
            <label>宽限期 (天)</label>
            <input
              type="number"
              v-model.number="quotaForm.gracePeriod"
              min="0"
              max="30"
              placeholder="7"
            />
            <small>超过软限制后的宽限期</small>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showSetQuota = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              设置配额
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 编辑配额对话框 -->
    <div v-if="showEditQuota" class="modal-overlay" @click="showEditQuota = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑组配额</h3>
          <button class="close-btn" @click="showEditQuota = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="updateQuota" class="modal-body">
          <div class="form-group">
            <label>用户组</label>
            <input type="text" :value="selectedQuota?.groupName" disabled />
          </div>

          <div class="form-group">
            <label>路径</label>
            <input type="text" :value="selectedQuota?.path" disabled />
          </div>

          <div class="current-usage">
            <h4>当前使用情况</h4>
            <div class="usage-display">
              <span class="usage-value">{{ formatBytes(selectedQuota?.usedSpace || 0) }}</span>
              <span class="usage-percent">{{ getUsagePercent(selectedQuota).toFixed(1) }}%</span>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>软限制 (MB)</label>
              <input
                type="number"
                v-model.number="editForm.softLimit"
                min="0"
              />
            </div>

            <div class="form-group">
              <label>硬限制 (MB)</label>
              <input
                type="number"
                v-model.number="editForm.hardLimit"
                min="0"
              />
            </div>
          </div>

          <div class="form-group">
            <label>宽限期 (天)</label>
            <input
              type="number"
              v-model.number="editForm.gracePeriod"
              min="0"
              max="30"
            />
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showEditQuota = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              保存更改
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 配额详情对话框 -->
    <div v-if="showQuotaDetails" class="modal-overlay" @click="showQuotaDetails = false">
      <div class="modal-content large-modal" @click.stop>
        <div class="modal-header">
          <h3>配额详情</h3>
          <button class="close-btn" @click="showQuotaDetails = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="modal-body">
          <div v-if="loadingDetails" class="loading-state">
            <div class="spinner"></div>
            <p>加载配额详情...</p>
          </div>

          <div v-else class="quota-details">
            <div class="detail-section">
              <h4>基本信息</h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <label>组名</label>
                  <span>{{ selectedQuota?.groupName }}</span>
                </div>
                <div class="detail-item">
                  <label>GID</label>
                  <span>{{ selectedQuota?.gid }}</span>
                </div>
                <div class="detail-item">
                  <label>成员数</label>
                  <span>{{ selectedQuota?.memberCount }}</span>
                </div>
                <div class="detail-item">
                  <label>路径</label>
                  <span>{{ selectedQuota?.path }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>使用情况</h4>
              <div class="usage-chart">
                <div class="usage-circle">
                  <svg viewBox="0 0 36 36">
                    <path
                      d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                      fill="none"
                      stroke="#e5e7eb"
                      stroke-width="3"
                    />
                    <path
                      d="M18 2.0845 a 15.9155 15.9155 0 0 1 0 31.831 a 15.9155 15.9155 0 0 1 0 -31.831"
                      fill="none"
                      :stroke="getUsageColor(selectedQuota)"
                      stroke-width="3"
                      :stroke-dasharray="getUsagePercent(selectedQuota) + ', 100'"
                    />
                  </svg>
                  <div class="circle-percent">
                    {{ getUsagePercent(selectedQuota).toFixed(1) }}%
                  </div>
                </div>
              </div>
              <div class="usage-stats">
                <div class="stat-item">
                  <label>使用量</label>
                  <span>{{ formatBytes(selectedQuota?.usedSpace || 0) }}</span>
                </div>
                <div class="stat-item">
                  <label>软限制</label>
                  <span>{{ formatBytes(selectedQuota?.softLimit || 0) }}</span>
                </div>
                <div class="stat-item">
                  <label>硬限制</label>
                  <span>{{ formatBytes(selectedQuota?.hardLimit || 0) }}</span>
                </div>
                <div class="stat-item">
                  <label>剩余空间</label>
                  <span>{{ formatBytes((selectedQuota?.hardLimit || 0) - (selectedQuota?.usedSpace || 0)) }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>组成员</h4>
              <div class="members-list">
                <div
                  v-for="member in groupMembers"
                  :key="member"
                  class="member-item"
                >
                  <UserIcon class="w-4 h-4" />
                  <span>{{ member }}</span>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  PlusIcon,
  ArrowPathIcon,
  UsersIcon,
  CircleStackIcon,
  ChartBarIcon,
  CheckBadgeIcon,
  PencilIcon,
  EyeIcon,
  TrashIcon,
  XMarkIcon,
  UserIcon
} from '@heroicons/vue/24/outline'

// 状态
const loading = ref(false)
const loadingDetails = ref(false)
const groupQuotas = ref<any[]>([])
const allGroups = ref<any[]>([])
const groupMembers = ref<string[]>([])

// 对话框状态
const showSetQuota = ref(false)
const showEditQuota = ref(false)
const showQuotaDetails = ref(false)

// 搜索和过滤
const searchQuery = ref('')
const filterStatus = ref('all')

// 选中的配额
const selectedQuota = ref<any>(null)

// 表单数据
const quotaForm = reactive({
  groupName: '',
  path: '/home',
  softLimit: 0,
  hardLimit: 0,
  gracePeriod: 7
})

const editForm = reactive({
  softLimit: 0,
  hardLimit: 0,
  gracePeriod: 7
})

// 计算属性
const filteredQuotas = computed(() => {
  let filtered = groupQuotas.value

  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(quota =>
      quota.name.toLowerCase().includes(query) ||
      quota.gid.toString().includes(query)
    )
  }

  // 状态过滤
  if (filterStatus.value !== 'all') {
    filtered = filtered.filter(quota => {
      const status = getUsageStatus(quota)
      return status === filterStatus.value
    })
  }

  return filtered
})

const limitedGroups = computed(() => {
  return groupQuotas.value.filter(quota =>
    quota.softLimit > 0 || quota.hardLimit > 0
  ).length
})

const totalUsage = computed(() => {
  return groupQuotas.value.reduce((sum, quota) => sum + (quota.usedSpace || 0), 0)
})

const totalAllocated = computed(() => {
  return groupQuotas.value.reduce((sum, quota) => {
    return sum + ((quota.hardLimit || quota.softLimit) || 0)
  }, 0)
})

const availableGroups = computed(() => {
  const groupsWithQuota = groupQuotas.value.map(q => q.name)
  return allGroups.value.filter(group =>
    !groupsWithQuota.includes(group.name)
  )
})

const topUsageGroups = computed(() => {
  return [...groupQuotas.value]
    .filter(quota => quota.usedSpace > 0)
    .sort((a, b) => b.usedSpace - a.usedSpace)
    .slice(0, 5)
})

// 方法
const loadGroupQuotas = async () => {
  loading.value = true
  try {
    const response = await fetch('/api/quota/groups', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      groupQuotas.value = data.groups || []
    } else {
      throw new Error('获取组配额失败')
    }
  } catch (error) {
    console.error('Failed to load group quotas:', error)
    ElMessage.error('获取组配额失败')
  } finally {
    loading.value = false
  }
}

const loadAllGroups = async () => {
  try {
    const response = await fetch('/api/groups', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      allGroups.value = data.groups || []
    }
  } catch (error) {
    console.error('Failed to load groups:', error)
  }
}

const setQuota = async () => {
  try {
    const response = await fetch(`/api/groups/${quotaForm.groupName}/quota`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: quotaForm.path,
        softLimit: quotaForm.softLimit * 1024 * 1024, // 转换为字节
        hardLimit: quotaForm.hardLimit * 1024 * 1024,
        gracePeriod: quotaForm.gracePeriod
      })
    })

    if (response.ok) {
      ElMessage.success('组配额设置成功')
      showSetQuota.value = false
      resetForm()
      await loadGroupQuotas()
    } else {
      throw new Error('设置组配额失败')
    }
  } catch (error) {
    console.error('Failed to set quota:', error)
    ElMessage.error('设置组配额失败')
  }
}

const editQuota = (quota: any) => {
  selectedQuota.value = quota
  editForm.softLimit = (quota.softLimit || 0) / (1024 * 1024)
  editForm.hardLimit = (quota.hardLimit || 0) / (1024 * 1024)
  editForm.gracePeriod = quota.gracePeriod || 7
  showEditQuota.value = true
}

const updateQuota = async () => {
  try {
    const response = await fetch(`/api/groups/${selectedQuota.value.groupName}/quota`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: selectedQuota.value.path,
        softLimit: editForm.softLimit * 1024 * 1024,
        hardLimit: editForm.hardLimit * 1024 * 1024,
        gracePeriod: editForm.gracePeriod
      })
    })

    if (response.ok) {
      ElMessage.success('组配额更新成功')
      showEditQuota.value = false
      await loadGroupQuotas()
    } else {
      throw new Error('更新组配额失败')
    }
  } catch (error) {
    console.error('Failed to update quota:', error)
    ElMessage.error('更新组配额失败')
  }
}

const removeQuota = (quota: any) => {
  if (!confirm(`确定要移除组 "${quota.name}" 的配额限制吗？`)) {
    return
  }

  fetch(`/api/groups/${quota.name}/quota`, {
    method: 'DELETE',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token')}`
    }
  })
  .then(response => {
    if (response.ok) {
      ElMessage.success('组配额移除成功')
      return loadGroupQuotas()
    }
    throw new Error('移除组配额失败')
  })
  .catch(error => {
    console.error('Failed to remove quota:', error)
    ElMessage.error('移除组配额失败')
  })
}

const viewQuotaDetails = async (quota: any) => {
  selectedQuota.value = quota
  loadingDetails.value = true
  showQuotaDetails.value = true

  try {
    // 获取组成员
    const response = await fetch(`/api/groups/${quota.name}/members`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      groupMembers.value = data.members?.map((m: any) => m.username) || []
    }
  } catch (error) {
    console.error('Failed to load group members:', error)
  } finally {
    loadingDetails.value = false
  }
}

const refreshQuotas = () => {
  loadGroupQuotas()
}

const resetForm = () => {
  quotaForm.groupName = ''
  quotaForm.path = '/home'
  quotaForm.softLimit = 0
  quotaForm.hardLimit = 0
  quotaForm.gracePeriod = 7
}

const getUsagePercent = (quota: any) => {
  if (!quota.hardLimit && !quota.softLimit) return 0
  const limit = quota.hardLimit || quota.softLimit
  if (limit === 0) return 0
  return Math.min((quota.usedSpace / limit) * 100, 100)
}

const getUsageStatus = (quota: any) => {
  const percent = getUsagePercent(quota)
  if (percent >= 90) return 'exceeded'
  if (percent >= 70) return 'warning'
  return 'normal'
}

const getUsageStatusLabel = (quota: any) => {
  const status = getUsageStatus(quota)
  switch (status) {
    case 'exceeded': return '超限'
    case 'warning': return '警告'
    case 'normal': return '正常'
    default: return '未知'
  }
}

const getUsageColor = (quota: any) => {
  const status = getUsageStatus(quota)
  switch (status) {
    case 'exceeded': return '#ef4444'
    case 'warning': return '#f59e0b'
    case 'normal': return '#10b981'
    default: return '#6b7280'
  }
}

const formatBytes = (bytes: number) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

// 生命周期
onMounted(() => {
  loadGroupQuotas()
  loadAllGroups()
})
</script>

<style scoped lang="scss">
.group-quota-manager {
  width: 100%;
  padding: 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.quota-header {
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

.quota-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.quota-overview {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
}

.overview-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  display: flex;
  align-items: center;
  gap: 16px;
}

.card-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;

  &.total {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }

  &.limited {
    background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);
  }

  &.usage {
    background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  }

  &.allocated {
    background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
  }
}

.card-info {
  flex: 1;

  h3 {
    font-size: 14px;
    font-weight: 500;
    color: #6b7280;
    margin: 0 0 4px 0;
  }

  .card-value {
    font-size: 24px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.quota-table-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h2 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.table-controls {
  display: flex;
  gap: 12px;
}

.search-input,
.filter-select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s;

  &:focus {
    outline: none;
    border-color: #667eea;
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
  }
}

.search-input {
  width: 200px;
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
  padding: 40px;
  color: #9ca3af;

  svg {
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    font-size: 14px;
  }
}

.quota-table {
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
      transition: background 0.2s;

      &:hover {
        background: rgba(102, 126, 234, 0.05);
      }

      &:last-child {
        border-bottom: none;
      }

      td {
        padding: 16px 12px;
        font-size: 14px;
      }
    }
  }
}

.group-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.group-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 16px;
}

.group-details {
  display: flex;
  flex-direction: column;
  gap: 2px;

  .group-name {
    font-weight: 500;
    color: #1f2937;
  }
}

.usage-info {
  display: flex;
  flex-direction: column;
  gap: 6px;

  .usage-value {
    font-weight: 500;
    color: #1f2937;
  }

  .usage-bar {
    width: 100px;
    height: 6px;
    background: #e5e7eb;
    border-radius: 3px;
    overflow: hidden;
  }

  .usage-fill {
    height: 100%;
    transition: width 0.3s;

    &.normal {
      background: #10b981;
    }

    &.warning {
      background: #f59e0b;
    }

    &.exceeded {
      background: #ef4444;
    }
  }
}

.limit-value {
  font-weight: 500;
  color: #1f2937;

  &.unlimited {
    color: #9ca3af;
    font-style: italic;
  }
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;

  &.normal {
    background: rgba(16, 185, 129, 0.1);
    color: #066;
  }

  &.warning {
    background: rgba(245, 158, 11, 0.1);
    color: #92400e;
  }

  &.exceeded {
    background: rgba(239, 68, 68, 0.1);
    color: #991b1b;
  }
}

.action-buttons {
  display: flex;
  gap: 4px;
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

.quota-trends {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

  h2 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 20px 0;
  }
}

.trends-chart {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.trend-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.trend-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
}

.trend-name {
  font-weight: 500;
  color: #1f2937;
}

.trend-usage {
  color: #6b7280;
  font-size: 13px;
}

.trend-bar {
  width: 100%;
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
}

.trend-fill {
  height: 100%;
  transition: width 0.3s;

  &.normal {
    background: #10b981;
  }

  &.warning {
    background: #f59e0b;
  }

  &.exceeded {
    background: #ef4444;
  }
}

.trend-percent {
  font-size: 12px;
  color: #6b7280;
  text-align: right;
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
    max-width: 800px;
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

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;

  label {
    font-size: 14px;
    font-weight: 500;
    color: #374151;
  }

  input, select {
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

  small {
    font-size: 12px;
    color: #6b7280;
    line-height: 1.4;
  }
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.current-usage {
  padding: 16px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
  border-radius: 8px;

  h4 {
    font-size: 14px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 12px 0;
  }
}

.usage-display {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .usage-value {
    font-size: 18px;
    font-weight: 600;
    color: #667eea;
  }

  .usage-percent {
    font-size: 16px;
    font-weight: 600;
    color: #10b981;
  }
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

.quota-details {
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

.usage-chart {
  display: flex;
  gap: 24px;
  align-items: center;
}

.usage-circle {
  position: relative;
  width: 120px;
  height: 120px;

  svg {
    width: 100%;
    height: 100%;
  }

  .circle-percent {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
}

.usage-stats {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  label {
    font-size: 12px;
    font-weight: 500;
    color: #6b7280;
  }

  span {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
  }
}

.members-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.member-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  color: #667eea;
}
</style>