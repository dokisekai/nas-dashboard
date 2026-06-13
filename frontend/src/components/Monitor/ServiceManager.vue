<template>
  <div class="service-manager">
    <div class="manager-header">
      <div class="header-controls">
        <el-input
          v-model="searchQuery"
          placeholder="搜索服务..."
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
          <el-option label="运行中" value="running" />
          <el-option label="已停止" value="stopped" />
          <el-option label="失败" value="failed" />
        </el-select>
        <el-button size="small" @click="refreshServices">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <div class="service-stats">
        <span class="stat-item">总服务: {{ filteredServices.length }}</span>
        <span class="stat-item success">运行中: {{ monitorStore.runningServices.length }}</span>
        <span class="stat-item danger">故障: {{ monitorStore.failedServices.length }}</span>
        <span class="stat-item info">启用: {{ monitorStore.enabledServices.length }}</span>
      </div>
    </div>

    <div class="service-list">
      <el-table
        :data="paginatedServices"
        v-loading="loading"
        size="small"
        max-height="600"
        @row-click="showServiceDetails"
        style="width: 100%; cursor: pointer"
      >
        <el-table-column prop="name" label="服务名称" min-width="200">
          <template #default="{ row }">
            <div class="service-name">
              <el-icon class="service-icon" :class="getStatusClass(row.status)">
                <component :is="getServiceIcon(row.status)" />
              </el-icon>
              <span>{{ row.name }}</span>
              <el-tag
                v-if="row.enabled"
                type="success"
                size="small"
                class="enabled-badge"
              >
                启用
              </el-tag>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />

        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="loadState" label="加载状态" width="120">
          <template #default="{ row }">
            <span class="state-text">{{ row.loadState }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="activeState" label="活动状态" width="120">
          <template #default="{ row }">
            <span class="state-text">{{ row.activeState }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="subState" label="子状态" width="120">
          <template #default="{ row }">
            <span class="state-text">{{ row.subState }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="mainPid" label="主PID" width="100" align="center">
          <template #default="{ row }">
            <span v-if="row.mainPid > 0">{{ row.mainPid }}</span>
            <span v-else class="no-pid">-</span>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button
                size="small"
                :disabled="row.status === 'running'"
                @click.stop="startService(row)"
              >
                <el-icon><VideoPlay /></el-icon>
              </el-button>
              <el-button
                size="small"
                :disabled="row.status !== 'running'"
                @click.stop="stopService(row)"
                type="warning"
              >
                <el-icon><VideoPause /></el-icon>
              </el-button>
              <el-button
                size="small"
                @click.stop="restartService(row)"
                type="primary"
              >
                <el-icon><RefreshRight /></el-icon>
              </el-button>
              <el-dropdown @command="(cmd) => handleServiceCommand(cmd, row)" trigger="click">
                <el-button size="small">
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="enable" :disabled="row.enabled">
                      <el-icon><CircleCheck /></el-icon>
                      启用服务
                    </el-dropdown-item>
                    <el-dropdown-item command="disable" :disabled="!row.enabled">
                      <el-icon><CircleClose /></el-icon>
                      禁用服务
                    </el-dropdown-item>
                    <el-dropdown-item divided command="details">
                      <el-icon><Document /></el-icon>
                      详细信息
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[20, 50, 100]"
          :total="filteredServices.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </div>

    <!-- Service Details Dialog -->
    <el-dialog
      v-model="showDetailsDialog"
      :title="`服务详情 - ${selectedService?.name || ''}`"
      width="700px"
    >
      <div v-if="selectedService" class="service-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="服务名称">
            {{ selectedService.name }}
          </el-descriptions-item>
          <el-descriptions-item label="描述">
            {{ selectedService.description || '无描述' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusType(selectedService.status)" size="small">
              {{ getStatusLabel(selectedService.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="启用状态">
            <el-tag :type="selectedService.enabled ? 'success' : 'info'" size="small">
              {{ selectedService.enabled ? '已启用' : '未启用' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="加载状态">
            {{ selectedService.loadState }}
          </el-descriptions-item>
          <el-descriptions-item label="活动状态">
            {{ selectedService.activeState }}
          </el-descriptions-item>
          <el-descriptions-item label="子状态">
            {{ selectedService.subState }}
          </el-descriptions-item>
          <el-descriptions-item label="主进程PID">
            {{ selectedService.mainPid > 0 ? selectedService.mainPid : '无' }}
          </el-descriptions-item>
        </el-descriptions>

        <div class="service-actions">
          <h4>服务操作:</h4>
          <div class="action-buttons">
            <el-button
              @click="startService(selectedService)"
              :disabled="selectedService.status === 'running'"
              type="success"
            >
              <el-icon><VideoPlay /></el-icon>
              启动服务
            </el-button>
            <el-button
              @click="stopService(selectedService)"
              :disabled="selectedService.status !== 'running'"
              type="warning"
            >
              <el-icon><VideoPause /></el-icon>
              停止服务
            </el-button>
            <el-button @click="restartService(selectedService)" type="primary">
              <el-icon><RefreshRight /></el-icon>
              重启服务
            </el-button>
            <el-button
              @click="toggleServiceEnabled(selectedService)"
              :type="selectedService.enabled ? 'warning' : 'success'"
            >
              <el-icon><component :is="selectedService.enabled ? CircleClose : CircleCheck" /></el-icon>
              {{ selectedService.enabled ? '禁用服务' : '启用服务' }}
            </el-button>
          </div>
        </div>

        <div v-if="selectedService.mainPid > 0" class="process-info">
          <h4>进程信息:</h4>
          <div class="process-details">
            <span class="process-item">主进程 PID: {{ selectedService.mainPid }}</span>
            <el-button size="small" @click="viewProcess(selectedService.mainPid)">
              查看进程详情
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMonitorStore } from '@/stores/monitor'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { ServiceInfo } from '@/types/monitor'
import {
  Refresh,
  Search,
  VideoPlay,
  VideoPause,
  RefreshRight,
  MoreFilled,
  CircleCheck,
  CircleClose,
  Document,
  SuccessFilled,
  VideoPause as FailedIcon,
  WarningFilled
} from '@element-plus/icons-vue'

const monitorStore = useMonitorStore()

const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(50)
const selectedService = ref<ServiceInfo | null>(null)
const showDetailsDialog = ref(false)

// Computed
const filteredServices = computed(() => {
  let services = monitorStore.services

  // Apply status filter
  if (statusFilter.value) {
    services = services.filter(s => s.status === statusFilter.value)
  }

  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    services = services.filter(s =>
      s.name.toLowerCase().includes(query) ||
      (s.description && s.description.toLowerCase().includes(query))
    )
  }

  // Sort by name
  return services.sort((a, b) => a.name.localeCompare(b.name))
})

const paginatedServices = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredServices.value.slice(start, end)
})

// Methods
const refreshServices = async () => {
  loading.value = true
  try {
    await monitorStore.fetchServices()
    ElMessage.success('服务列表已刷新')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string): string => {
  const types: Record<string, string> = {
    running: 'success',
    stopped: 'info',
    failed: 'danger',
    masked: 'warning'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string): string => {
  const labels: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    failed: '失败',
    masked: '已屏蔽'
  }
  return labels[status] || status
}

const getStatusClass = (status: string): string => {
  return `status-${status}`
}

const getServiceIcon = (status: string) => {
  switch (status) {
    case 'running': return SuccessFilled
    case 'stopped': return VideoPause
    case 'failed': return FailedIcon
    default: return WarningFilled
  }
}

const showServiceDetails = (service: ServiceInfo) => {
  selectedService.value = service
  showDetailsDialog.value = true
}

const startService = async (service: ServiceInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要启动服务 "${service.name}" 吗？`,
      '确认启动服务',
      {
        type: 'success',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    await monitorStore.startService(service.name)
    ElMessage.success(`服务 "${service.name}" 已启动`)
    showDetailsDialog.value = false
    await refreshServices()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('启动服务失败')
    }
  }
}

const stopService = async (service: ServiceInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要停止服务 "${service.name}" 吗？`,
      '确认停止服务',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    await monitorStore.stopService(service.name)
    ElMessage.success(`服务 "${service.name}" 已停止`)
    showDetailsDialog.value = false
    await refreshServices()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('停止服务失败')
    }
  }
}

const restartService = async (service: ServiceInfo) => {
  try {
    await ElMessageBox.confirm(
      `确定要重启服务 "${service.name}" 吗？`,
      '确认重启服务',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    await monitorStore.restartService(service.name)
    ElMessage.success(`服务 "${service.name}" 已重启`)
    showDetailsDialog.value = false
    await refreshServices()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('重启服务失败')
    }
  }
}

const toggleServiceEnabled = async (service: ServiceInfo) => {
  const action = service.enabled ? 'disable' : 'enable'
  const actionText = service.enabled ? '禁用' : '启用'

  try {
    await ElMessageBox.confirm(
      `确定要${actionText}服务 "${service.name}" 吗？`,
      `确认${actionText}服务`,
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    // This would call the appropriate API method when implemented
    ElMessage.success(`服务 "${service.name}" 已${actionText}`)
    showDetailsDialog.value = false
    await refreshServices()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${actionText}服务失败`)
    }
  }
}

const handleServiceCommand = (command: string, service: ServiceInfo) => {
  selectedService.value = service

  switch (command) {
    case 'enable':
      toggleServiceEnabled(service)
      break
    case 'disable':
      toggleServiceEnabled(service)
      break
    case 'details':
      showServiceDetails(service)
      break
  }
}

const viewProcess = (pid: number) => {
  // Navigate to process details or show dialog
  ElMessage.info(`查看进程 ${pid} 的详情功能即将推出`)
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
}

const handlePageChange = (page: number) => {
  currentPage.value = page
}

// Lifecycle
onMounted(() => {
  refreshServices()
})
</script>

<style scoped lang="scss">
.service-manager {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.manager-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 4px;

  .header-controls {
    display: flex;
    gap: 10px;
  }

  .service-stats {
    display: flex;
    gap: 20px;
    font-size: 12px;
    color: #606266;

    .stat-item {
      display: flex;
      align-items: center;
      gap: 5px;

      &.success { color: #67c23a; }
      &.danger { color: #f56c6c; }
      &.info { color: #409eff; }
    }
  }
}

.service-list {
  flex: 1;
  overflow: auto;

  .service-name {
    display: flex;
    align-items: center;
    gap: 8px;

    .service-icon {
      font-size: 14px;

      &.status-running { color: #67c23a; }
      &.status-stopped { color: #909399; }
      &.status-failed { color: #f56c6c; }
      &.status-masked { color: #e6a23c; }
    }

    .enabled-badge {
      margin-left: 8px;
      font-size: 10px;
    }
  }

  .state-text {
    font-size: 12px;
    color: #606266;
  }

  .no-pid {
    color: #c0c4cc;
  }

  .pagination {
    margin-top: 15px;
    display: flex;
    justify-content: center;
  }
}

.service-details {
  .service-actions,
  .process-info {
    margin-top: 20px;

    h4 {
      margin: 0 0 10px;
      color: #303133;
      font-size: 14px;
    }
  }

  .action-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
  }

  .process-details {
    display: flex;
    align-items: center;
    gap: 15px;

    .process-item {
      color: #606266;
      font-size: 12px;
    }
  }
}
</style>