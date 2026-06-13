<template>
  <div class="alert-manager">
    <div class="manager-header">
      <div class="header-controls">
        <el-button size="small" type="primary" @click="showCreateDialog">
          <el-icon><Plus /></el-icon>
          创建告警
        </el-button>
        <el-button size="small" @click="refreshAlerts" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
      <div class="alert-stats">
        <span class="stat-item">总告警: {{ alertRules.length }}</span>
        <span class="stat-item success">启用: {{ activeAlerts.length }}</span>
        <span class="stat-item danger">严重: {{ criticalAlerts.length }}</span>
      </div>
    </div>

    <div class="alert-list">
      <el-table
        :data="alertRules"
        v-loading="loading"
        size="small"
        max-height="500"
        style="width: 100%"
      >
        <el-table-column prop="name" label="告警名称" min-width="150">
          <template #default="{ row }">
            <div class="alert-name">
              <el-icon class="alert-icon" :class="getSeverityClass(row.severity)">
                <Bell />
              </el-icon>
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="type" label="监控类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getTypeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="condition" label="条件" width="80" align="center">
          <template #default="{ row }">
            <span class="condition-symbol">{{ getConditionSymbol(row.condition) }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="threshold" label="阈值" width="100">
          <template #default="{ row }">
            {{ formatThreshold(row) }}
          </template>
        </el-table-column>

        <el-table-column prop="duration" label="持续时间" width="100">
          <template #default="{ row }">
            {{ row.duration }} 秒
          </template>
        </el-table-column>

        <el-table-column prop="severity" label="严重程度" width="100">
          <template #default="{ row }">
            <el-tag :type="getSeverityType(row.severity)" size="small">
              {{ getSeverityLabel(row.severity) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="enabled" label="状态" width="80">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              @change="toggleAlert(row)"
              size="small"
            />
          </template>
        </el-table-column>

        <el-table-column prop="lastTriggered" label="最后触发" width="160">
          <template #default="{ row }">
            {{ row.lastTriggered ? formatTime(row.lastTriggered) : '从未触发' }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button size="small" @click="editAlert(row)">
                <el-icon><Edit /></el-icon>
              </el-button>
              <el-button size="small" type="danger" @click="deleteAlert(row)">
                <el-icon><Delete /></el-icon>
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- Create/Edit Alert Dialog -->
    <el-dialog
      v-model="showDialog"
      :title="editingAlert ? '编辑告警规则' : '创建告警规则'"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="alertForm"
        :rules="formRules"
        label-width="120px"
      >
        <el-form-item label="告警名称" prop="name">
          <el-input v-model="alertForm.name" placeholder="请输入告警名称" />
        </el-form-item>

        <el-form-item label="监控类型" prop="type">
          <el-select v-model="alertForm.type" placeholder="选择监控类型">
            <el-option label="CPU" value="cpu" />
            <el-option label="内存" value="memory" />
            <el-option label="磁盘" value="disk" />
            <el-option label="温度" value="temperature" />
            <el-option label="网络" value="network" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发条件" prop="condition">
          <el-select v-model="alertForm.condition" placeholder="选择条件">
            <el-option label="大于" value=">" />
            <el-option label="小于" value="<" />
            <el-option label="大于等于" value=">=" />
            <el-option label="小于等于" value="<=" />
            <el-option label="等于" value="==" />
            <el-option label="不等于" value="!=" />
          </el-select>
        </el-form-item>

        <el-form-item label="阈值" prop="threshold">
          <el-input-number
            v-model="alertForm.threshold"
            :min="0"
            :max="100"
            :precision="2"
            controls-position="right"
          />
          <span class="unit-hint">{{ getUnitHint(alertForm.type) }}</span>
        </el-form-item>

        <el-form-item label="持续时间" prop="duration">
          <el-input-number
            v-model="alertForm.duration"
            :min="0"
            :max="3600"
            controls-position="right"
          />
          <span class="unit-hint">秒</span>
        </el-form-item>

        <el-form-item label="严重程度" prop="severity">
          <el-radio-group v-model="alertForm.severity">
            <el-radio label="info">信息</el-radio>
            <el-radio label="warning">警告</el-radio>
            <el-radio label="critical">严重</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="alertForm.enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showDialog = false">取消</el-button>
          <el-button type="primary" @click="saveAlert" :loading="saving">
            {{ editingAlert ? '更新' : '创建' }}
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useMonitorStore } from '@/stores/monitor'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { AlertRule } from '@/types/monitor'
import { Plus, Refresh, Bell, Edit, Delete } from '@element-plus/icons-vue'

const monitorStore = useMonitorStore()

const loading = ref(false)
const showDialog = ref(false)
const editingAlert = ref<AlertRule | null>(null)
const saving = ref(false)
const formRef = ref()

const alertForm = ref({
  name: '',
  type: 'cpu' as 'cpu' | 'memory' | 'disk' | 'temperature' | 'network',
  condition: '>' as '>' | '<' | '>=' | '<=' | '==' | '!=',
  threshold: 80,
  duration: 30,
  severity: 'warning' as 'info' | 'warning' | 'critical',
  enabled: true
})

const formRules = {
  name: [
    { required: true, message: '请输入告警名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择监控类型', trigger: 'change' }
  ],
  condition: [
    { required: true, message: '请选择触发条件', trigger: 'change' }
  ],
  threshold: [
    { required: true, message: '请输入阈值', trigger: 'blur' }
  ],
  duration: [
    { required: true, message: '请输入持续时间', trigger: 'blur' }
  ],
  severity: [
    { required: true, message: '请选择严重程度', trigger: 'change' }
  ]
}

// Computed
const alertRules = computed(() => monitorStore.alerts)
const activeAlerts = computed(() => monitorStore.activeAlerts)
const criticalAlerts = computed(() => monitorStore.criticalAlerts)

// Methods
const refreshAlerts = async () => {
  loading.value = true
  try {
    await monitorStore.fetchAlerts()
    ElMessage.success('告警规则已刷新')
  } finally {
    loading.value = false
  }
}

const getTypeLabel = (type: string): string => {
  const labels: Record<string, string> = {
    cpu: 'CPU',
    memory: '内存',
    disk: '磁盘',
    temperature: '温度',
    network: '网络'
  }
  return labels[type] || type
}

const getConditionSymbol = (condition: string): string => {
  const symbols: Record<string, string> = {
    '>': '>',
    '<': '<',
    '>=': '≥',
    '<=': '≤',
    '==': '=',
    '!=': '≠'
  }
  return symbols[condition] || condition
}

const formatThreshold = (alert: AlertRule): string => {
  switch (alert.type) {
    case 'cpu':
    case 'memory':
      return `${alert.threshold}%`
    case 'temperature':
      return `${alert.threshold}°C`
    default:
      return alert.threshold.toString()
  }
}

const getUnitHint = (type: string): string => {
  const units: Record<string, string> = {
    cpu: '%',
    memory: '%',
    temperature: '°C',
    disk: '%',
    network: 'MB/s'
  }
  return units[type] || ''
}

const getSeverityType = (severity: string): string => {
  const types: Record<string, string> = {
    info: 'info',
    warning: 'warning',
    critical: 'danger'
  }
  return types[severity] || 'info'
}

const getSeverityLabel = (severity: string): string => {
  const labels: Record<string, string> = {
    info: '信息',
    warning: '警告',
    critical: '严重'
  }
  return labels[severity] || severity
}

const getSeverityClass = (severity: string): string => {
  return `severity-${severity}`
}

const formatTime = (dateStr: string): string => {
  return new Date(dateStr).toLocaleString('zh-CN')
}

const showCreateDialog = () => {
  editingAlert.value = null
  alertForm.value = {
    name: '',
    type: 'cpu',
    condition: '>',
    threshold: 80,
    duration: 30,
    severity: 'warning',
    enabled: true
  }
  showDialog.value = true
}

const editAlert = (alert: AlertRule) => {
  editingAlert.value = alert
  alertForm.value = {
    name: alert.name,
    type: alert.type,
    condition: alert.condition,
    threshold: alert.threshold,
    duration: alert.duration,
    severity: alert.severity,
    enabled: alert.enabled
  }
  showDialog.value = true
}

const saveAlert = async () => {
  try {
    await formRef.value?.validate()

    saving.value = true

    if (editingAlert.value) {
      await monitorStore.updateAlert(editingAlert.value.id, alertForm.value)
      ElMessage.success('告警规则已更新')
    } else {
      await monitorStore.createAlert(alertForm.value)
      ElMessage.success('告警规则已创建')
    }

    showDialog.value = false
    await refreshAlerts()
  } catch (error) {
    if (error !== false) {
      ElMessage.error('保存告警规则失败')
    }
  } finally {
    saving.value = false
  }
}

const toggleAlert = async (alert: AlertRule) => {
  try {
    await monitorStore.updateAlert(alert.id, { enabled: alert.enabled })
    ElMessage.success(`告警 "${alert.name}" 已${alert.enabled ? '启用' : '禁用'}`)
  } catch (error) {
    // Revert the switch if failed
    alert.enabled = !alert.enabled
    ElMessage.error('更新告警状态失败')
  }
}

const deleteAlert = async (alert: AlertRule) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除告警规则 "${alert.name}" 吗？`,
      '确认删除',
      {
        type: 'warning',
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }
    )

    await monitorStore.deleteAlert(alert.id)
    ElMessage.success('告警规则已删除')
    await refreshAlerts()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除告警规则失败')
    }
  }
}

// Lifecycle
onMounted(() => {
  refreshAlerts()
})
</script>

<style scoped lang="scss">
.alert-manager {
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

  .alert-stats {
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
    }
  }
}

.alert-list {
  flex: 1;
  overflow: auto;

  .alert-name {
    display: flex;
    align-items: center;
    gap: 8px;

    .alert-icon {
      font-size: 14px;

      &.severity-info { color: #409eff; }
      &.severity-warning { color: #e6a23c; }
      &.severity-critical { color: #f56c6c; }
    }
  }

  .condition-symbol {
    font-weight: bold;
    color: #409eff;
  }
}

.unit-hint {
  margin-left: 8px;
  color: #909399;
  font-size: 12px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>