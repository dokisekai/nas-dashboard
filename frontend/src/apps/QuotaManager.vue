<template>
  <div class="quota-manager">
    <div class="manager-header">
      <h2>配额管理</h2>
      <div class="header-actions">
        <el-button type="primary" @click="showUserQuotaDialog">
          <el-icon><User /></el-icon>
          用户配额
        </el-button>
        <el-button type="primary" @click="showGroupQuotaDialog">
          <el-icon><UserFilled /></el-icon>
          组配额
        </el-button>
        <el-button @click="generateReport">
          <el-icon><Document /></el-icon>
          生成报告
        </el-button>
        <el-button @click="refreshQuotas">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- Quota Overview -->
    <div class="quota-overview">
      <div class="overview-cards">
        <div class="card total-users">
          <div class="card-icon"><el-icon><User /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ quotaStore.totalUserQuotas }}</div>
            <div class="card-label">用户配额</div>
          </div>
        </div>
        <div class="card total-groups">
          <div class="card-icon"><el-icon><UserFilled /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ quotaStore.totalGroupQuotas }}</div>
            <div class="card-label">组配额</div>
          </div>
        </div>
        <div class="card exceeded">
          <div class="card-icon"><el-icon><Warning /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ quotaStore.exceededQuotas.length }}</div>
            <div class="card-label">超限配额</div>
          </div>
        </div>
        <div class="card usage">
          <div class="card-icon"><el-icon><PieChart /></el-icon></div>
          <div class="card-content">
            <div class="card-value">{{ quotaStore.quotaUsagePercent.toFixed(1) }}%</div>
            <div class="card-label">总使用率</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <el-tabs v-model="activeTab" class="quota-tabs">
      <!-- User Quotas -->
      <el-tab-pane label="用户配额" name="users">
        <div class="quota-section">
          <div class="section-toolbar">
            <div class="search-filters">
              <el-input
                v-model="searchQuery"
                placeholder="搜索用户..."
                prefix-icon="Search"
                size="small"
                style="width: 200px"
                clearable
              />
              <el-select
                v-model="statusFilter"
                placeholder="状态筛选"
                size="small"
                style="width: 120px"
                clearable
              >
                <el-option label="全部" value="" />
                <el-option label="正常" value="ok" />
                <el-option label="警告" value="warning" />
                <el-option label="超限" value="exceeded" />
              </el-select>
            </div>
          </div>

          <el-table
            :data="filteredUserQuotas"
            v-loading="quotaStore.loading"
            size="small"
            max-height="500"
            style="width: 100%"
          >
            <el-table-column prop="user.username" label="用户" min-width="150" />
            <el-table-column prop="path" label="路径" min-width="200" />
            <el-table-column label="空间限制" width="200">
              <template #default="{ row }">
                <div class="quota-limits">
                  <div>软限制: {{ formatSize(row.softLimit) }}</div>
                  <div>硬限制: {{ formatSize(row.hardLimit) }}</div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="使用情况" width="200">
              <template #default="{ row }">
                <div class="quota-usage">
                  <el-progress
                    :percentage="getUsagePercent(row.usedSpace, row.hardLimit)"
                    :color="getUsageColor(row.usedSpace, row.hardLimit)"
                    :show-text="true"
                    :stroke-width="8"
                  />
                  <div class="usage-text">
                    {{ formatSize(row.usedSpace) }} / {{ formatSize(row.hardLimit) }}
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="gracePeriod" label="宽限期" width="100" align="center">
              <template #default="{ row }">
                {{ row.gracePeriod }} 天
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getQuotaStatusType(row)" size="small">
                  {{ getQuotaStatusLabel(row) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="editUserQuota(row)">
                    <el-icon><Edit /></el-icon>
                  </el-button>
                  <el-button size="small" type="danger" @click="deleteUserQuota(row)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- Group Quotas -->
      <el-tab-pane label="组配额" name="groups">
        <div class="quota-section">
          <div class="section-toolbar">
            <div class="search-filters">
              <el-input
                v-model="searchQuery"
                placeholder="搜索组..."
                prefix-icon="Search"
                size="small"
                style="width: 200px"
                clearable
              />
            </div>
          </div>

          <el-table
            :data="filteredGroupQuotas"
            v-loading="quotaStore.loading"
            size="small"
            max-height="500"
            style="width: 100%"
          >
            <el-table-column prop="group.name" label="组名" min-width="150" />
            <el-table-column prop="path" label="路径" min-width="200" />
            <el-table-column label="空间限制" width="200">
              <template #default="{ row }">
                <div class="quota-limits">
                  <div>软限制: {{ formatSize(row.softLimit) }}</div>
                  <div>硬限制: {{ formatSize(row.hardLimit) }}</div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="使用情况" width="200">
              <template #default="{ row }">
                <div class="quota-usage">
                  <el-progress
                    :percentage="getUsagePercent(row.usedSpace, row.hardLimit)"
                    :color="getUsageColor(row.usedSpace, row.hardLimit)"
                    :show-text="true"
                    :stroke-width="8"
                  />
                  <div class="usage-text">
                    {{ formatSize(row.usedSpace) }} / {{ formatSize(row.hardLimit) }}
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="gracePeriod" label="宽限期" width="100" align="center">
              <template #default="{ row }">
                {{ row.gracePeriod }} 天
              </template>
            </el-table-column>
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" @click="editGroupQuota(row)">
                    <el-icon><Edit /></el-icon>
                  </el-button>
                  <el-button size="small" type="danger" @click="deleteGroupQuota(row)">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- Quota Reports -->
      <el-tab-pane label="配额报告" name="reports">
        <div class="quota-section">
          <div class="section-toolbar">
            <div class="report-filters">
              <el-select
                v-model="reportType"
                placeholder="报告类型"
                size="small"
                style="width: 120px"
              >
                <el-option label="全部" value="all" />
                <el-option label="用户" value="user" />
                <el-option label="组" value="group" />
              </el-select>
              <el-button size="small" @click="generateReport">
                <el-icon><Refresh /></el-icon>
                生成报告
              </el-button>
              <el-button size="small" @click="exportReport" :disabled="quotaReports.length === 0">
                <el-icon><Download /></el-icon>
                导出报告
              </el-button>
            </div>
          </div>

          <el-table
            :data="quotaReports"
            v-loading="quotaStore.loading"
            size="small"
            max-height="500"
            style="width: 100%"
          >
            <el-table-column prop="name" label="名称" min-width="150" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="row.type === 'user' ? 'primary' : 'success'" size="small">
                  {{ row.type === 'user' ? '用户' : '组' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="path" label="路径" min-width="200" />
            <el-table-column label="使用情况" width="200">
              <template #default="{ row }">
                <div class="quota-usage">
                  <el-progress
                    :percentage="row.usedPercent"
                    :color="getUsageColor(row.usedSpace, row.hardLimit)"
                    :show-text="true"
                    :stroke-width="8"
                  />
                  <div class="usage-text">
                    {{ formatSize(row.usedSpace) }} / {{ formatSize(row.hardLimit) }}
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusLabel(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="generatedAt" label="生成时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.generatedAt) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <!-- Top Consumers -->
      <el-tab-pane label="高使用率排行" name="top">
        <div class="quota-section">
          <h3>配额使用率排行</h3>
          <div class="top-consumers">
            <div
              v-for="(item, index) in topConsumers"
              :key="index"
              class="consumer-item"
              :class="getConsumerClass(item)"
            >
              <div class="consumer-rank">{{ index + 1 }}</div>
              <div class="consumer-info">
                <div class="consumer-name">{{ item.name || item.user?.username || item.group?.name }}</div>
                <div class="consumer-path">{{ item.path }}</div>
                <div class="consumer-usage">
                  <el-progress
                    :percentage="calculateQuotaPercent(item.usedSpace, item.hardLimit)"
                    :color="getUsageColor(item.usedSpace, item.hardLimit)"
                    :show-text="true"
                    :stroke-width="10"
                  />
                </div>
                <div class="consumer-details">
                  <span>已用: {{ formatSize(item.usedSpace) }}</span>
                  <span>限制: {{ formatSize(item.hardLimit) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- User Quota Dialog -->
    <UserQuotaEditor
      v-model:visible="showUserQuotaDialog"
      :quota="editingQuota"
      @saved="onQuotaSaved"
    />

    <!-- Group Quota Dialog -->
    <GroupQuotaEditor
      v-model:visible="showGroupQuotaDialog"
      :quota="editingQuota"
      @saved="onQuotaSaved"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useQuotaStore } from '@/stores/quota'
import { quotaAPI } from '@/api/quota'
import { ElMessage } from 'element-plus'
import type { UserQuota, GroupQuota, QuotaReport } from '@/types/quota'
import {
  User,
  UserFilled,
  Document,
  Refresh,
  Search,
  Edit,
  Delete,
  Warning,
  PieChart,
  Download
} from '@element-plus/icons-vue'
import { calculateQuotaPercent, formatBytes, getQuotaStatus } from '@/types/quota'
import UserQuotaEditor from '@/components/Quota/UserQuotaEditor.vue'
import GroupQuotaEditor from '@/components/Quota/GroupQuotaEditor.vue'

const quotaStore = useQuotaStore()

const activeTab = ref('users')
const searchQuery = ref('')
const statusFilter = ref('')
const reportType = ref('all')
const showUserQuotaDialog = ref(false)
const showGroupQuotaDialog = ref(false)
const editingQuota = ref<UserQuota | GroupQuota | null>(null)

// Computed
const filteredUserQuotas = computed(() => {
  let quotas = quotaStore.userQuotas

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    quotas = quotas.filter(q =>
      q.user.username.toLowerCase().includes(query) ||
      q.path.toLowerCase().includes(query)
    )
  }

  if (statusFilter.value) {
    quotas = quotas.filter(q => getQuotaStatus(q.usedSpace, q.softLimit, q.hardLimit) === statusFilter.value)
  }

  return quotas
})

const filteredGroupQuotas = computed(() => {
  let quotas = quotaStore.groupQuotas

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    quotas = quotas.filter(q =>
      q.group.name.toLowerCase().includes(query) ||
      q.path.toLowerCase().includes(query)
    )
  }

  return quotas
})

const quotaReports = computed(() => quotaStore.quotaReports)

const topConsumers = computed(() => quotaStore.topQuotaConsumers)

// Methods
const formatSize = (bytes: number): string => formatBytes(bytes)

const getUsagePercent = (used: number, limit: number): number => {
  if (limit === 0) return 0
  return Math.min((used / limit) * 100, 100)
}

const getUsageColor = (used: number, limit: number): string => {
  const percent = getUsagePercent(used, limit)
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

const getQuotaStatusType = (quota: UserQuota | GroupQuota): string => {
  const status = getQuotaStatus(quota.usedSpace, quota.softLimit, quota.hardLimit)
  if (status === 'ok') return 'success'
  if (status === 'warning') return 'warning'
  return 'danger'
}

const getQuotaStatusLabel = (quota: UserQuota | GroupQuota): string => {
  const status = getQuotaStatus(quota.usedSpace, quota.softLimit, quota.hardLimit)
  if (status === 'ok') return '正常'
  if (status === 'warning') return '警告'
  return '超限'
}

const getStatusType = (status: string): string => {
  if (status === 'ok') return 'success'
  if (status === 'warning') return 'warning'
  return 'danger'
}

const getStatusLabel = (status: string): string => {
  if (status === 'ok') return '正常'
  if (status === 'warning') return '警告'
  return '超限'
}

const getConsumerClass = (item: any): string => {
  const percent = calculateQuotaPercent(item.usedSpace, item.hardLimit)
  if (percent >= 90) return 'critical'
  if (percent >= 70) return 'warning'
  return 'normal'
}

const refreshQuotas = async () => {
  await Promise.all([
    quotaStore.fetchUserQuotas(),
    quotaStore.fetchGroupQuotas(),
    quotaStore.fetchQuotaReports()
  ])
  ElMessage.success('配额信息已刷新')
}

const generateReport = async () => {
  try {
    await quotaStore.fetchQuotaReports({
      type: reportType.value === 'all' ? undefined : reportType.value
    })
    ElMessage.success('配额报告已生成')
  } catch (error) {
    ElMessage.error('生成配额报告失败')
  }
}

const exportReport = () => {
  if (quotaReports.value.length === 0) {
    ElMessage.warning('没有可导出的报告数据')
    return
  }

  const reportData = quotaReports.value.map(report => ({
    name: report.name,
    type: report.type,
    path: report.path,
    usedSpace: formatBytes(report.usedSpace),
    softLimit: formatBytes(report.softLimit),
    hardLimit: formatBytes(report.hardLimit),
    usedPercent: report.usedPercent.toFixed(2) + '%',
    status: report.status,
    generatedAt: report.generatedAt
  }))

  const blob = new Blob([JSON.stringify(reportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `quota_report_${Date.now()}.json`
  link.click()
  URL.revokeObjectURL(url)

  ElMessage.success('配额报告已导出')
}

const editUserQuota = (quota: UserQuota) => {
  editingQuota.value = quota
  showUserQuotaDialog.value = true
}

const editGroupQuota = (quota: GroupQuota) => {
  editingQuota.value = quota
  showGroupQuotaDialog.value = true
}

const deleteUserQuota = async (quota: UserQuota) => {
  try {
    await quotaStore.deleteUserQuota(quota.user.username, quota.path)
    ElMessage.success('用户配额已删除')
    await refreshQuotas()
  } catch (error) {
    ElMessage.error('删除用户配额失败')
  }
}

const deleteGroupQuota = async (quota: GroupQuota) => {
  try {
    await quotaStore.deleteGroupQuota(quota.group.id, quota.path)
    ElMessage.success('组配额已删除')
    await refreshQuotas()
  } catch (error) {
    ElMessage.error('删除组配额失败')
  }
}

const onQuotaSaved = () => {
  refreshQuotas()
}

// Lifecycle
onMounted(() => {
  quotaStore.init()
})
</script>

<style scoped lang="scss">
.quota-manager {
  padding: 20px;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  h2 {
    margin: 0;
    color: #303133;
    font-size: 24px;
  }

  .header-actions {
    display: flex;
    gap: 10px;
  }
}

.quota-overview {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;
  margin-bottom: 20px;

  .overview-cards {
    .card {
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      display: flex;
      align-items: center;
      gap: 15px;

      .card-icon {
        font-size: 32px;
        color: #409eff;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 60px;
        height: 60px;
        border-radius: 50%;
        background: #ecf5ff;
      }

      .card-content {
        .card-value {
          font-size: 28px;
          font-weight: bold;
          color: #303133;
          line-height: 1;
          margin-bottom: 5px;
        }

        .card-label {
          font-size: 14px;
          color: #909399;
        }
      }

      &.total-users .card-icon { color: #409eff; background: #ecf5ff; }
      &.total-groups .card-icon { color: #67c23a; background: #f0f9ff; }
      &.exceeded .card-icon { color: #f56c6c; background: #fef0f0; }
      &.usage .card-icon { color: #e6a23c; background: #fdf6ec; }
    }
  }
}

.quota-tabs {
  flex: 1;
  min-height: 0;

  :deep(.el-tabs__content) {
    height: 100%;
  }

  :deep(.el-tab-pane) {
    height: 100%;
    overflow: auto;
  }
}

.quota-section {
  height: 100%;
  display: flex;
  flex-direction: column;

  .section-toolbar {
    margin-bottom: 15px;
    padding: 15px;
    background: #f5f7fa;
    border-radius: 4px;

    .search-filters,
    .report-filters {
      display: flex;
      gap: 10px;
    }
  }

  .quota-limits {
    font-size: 12px;
    color: #606266;

    div {
      margin-bottom: 2px;
    }
  }

  .quota-usage {
    .usage-text {
      font-size: 11px;
      color: #606266;
      margin-top: 4px;
    }
  }
}

.top-consumers {
  display: flex;
  flex-direction: column;
  gap: 15px;

  .consumer-item {
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    display: flex;
    align-items: center;
    gap: 20px;

    &.critical { border-left: 4px solid #f56c6c; }
    &.warning { border-left: 4px solid #e6a23c; }
    &.normal { border-left: 4px solid #67c23a; }

    .consumer-rank {
      font-size: 24px;
      font-weight: bold;
      color: #303133;
      width: 40px;
      text-align: center;
    }

    .consumer-info {
      flex: 1;

      .consumer-name {
        font-size: 16px;
        font-weight: 500;
        color: #303133;
        margin-bottom: 4px;
      }

      .consumer-path {
        font-size: 12px;
        color: #909399;
        margin-bottom: 10px;
      }

      .consumer-usage {
        margin-bottom: 8px;
      }

      .consumer-details {
        display: flex;
        gap: 15px;
        font-size: 11px;
        color: #606266;
      }
    }
  }
}
</style>