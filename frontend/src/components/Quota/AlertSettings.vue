<template>
  <div class="alert-settings">
    <div class="settings-header">
      <h3>配额告警设置</h3>
      <el-button size="small" type="primary" @click="createAlertRule">
        <el-icon><Plus /></el-icon>
        新增告警规则
      </el-button>
    </div>

    <!-- Alert Rules List -->
    <div class="alert-rules-list">
      <div
        v-for="rule in alertRules"
        :key="rule.id"
        class="alert-rule-item"
        :class="{ disabled: !rule.enabled }"
      >
        <div class="rule-info">
          <div class="rule-header">
            <div class="rule-name">
              <el-icon>
                <component :is="getAlertIcon(rule.type)" />
              </el-icon>
              {{ rule.name }}
            </div>
            <div class="rule-status">
              <el-switch
                v-model="rule.enabled"
                @change="toggleAlertRule(rule)"
                size="small"
              />
            </div>
          </div>
          <div class="rule-details">
            <div class="rule-condition">
              <span class="detail-label">触发条件:</span>
              <span class="detail-value">{{ formatCondition(rule) }}</span>
            </div>
            <div class="rule-threshold">
              <span class="detail-label">阈值:</span>
              <span class="detail-value">{{ rule.threshold }}%</span>
            </div>
            <div class="rule-severity">
              <el-tag :type="getSeverityType(rule.severity)" size="small">
                {{ getSeverityLabel(rule.severity) }}
              </el-tag>
            </div>
          </div>
          <div class="rule-actions">
            <span class="action-label">告警动作:</span>
            <div class="action-tags">
              <el-tag
                v-for="action in parseActions(rule.actions)"
                :key="action.type"
                size="small"
                style="margin-right: 5px"
              >
                {{ getActionLabel(action.type) }}
              </el-tag>
            </div>
          </div>
        </div>
        <div class="rule-operations">
          <el-button
            size="small"
            @click="editAlertRule(rule)"
          >
            编辑
          </el-button>
          <el-button
            size="small"
            type="danger"
            @click="deleteAlertRule(rule)"
          >
            删除
          </el-button>
        </div>
      </div>
    </div>

    <!-- Alert Rule Editor Dialog -->
    <el-dialog
      v-model="showRuleDialog"
      :title="editingRule ? '编辑告警规则' : '新增告警规则'"
      width="600px"
    >
      <el-form
        ref="ruleFormRef"
        :model="ruleForm"
        :rules="ruleFormRules"
        label-width="120px"
      >
        <el-form-item label="规则名称" prop="name">
          <el-input
            v-model="ruleForm.name"
            placeholder="输入规则名称"
            prefix-icon="Warning"
          />
        </el-form-item>

        <el-form-item label="告警类型" prop="type">
          <el-select v-model="ruleForm.type" placeholder="选择告警类型">
            <el-option label="用户配额告警" value="user_quota" />
            <el-option label="组配额告警" value="group_quota" />
            <el-option label="磁盘空间告警" value="disk_space" />
            <el-option label="存储池告警" value="storage_pool" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发条件" prop="condition">
          <el-select v-model="ruleForm.condition">
            <el-option label="大于 (>)" value=">" />
            <el-option label="小于 (<)" value="<" />
            <el-option label="大于等于 (>=)" value=">=" />
            <el-option label="小于等于 (<=)" value="<=" />
            <el-option label="等于 (=)" value="=" />
          </el-select>
        </el-form-item>

        <el-form-item label="阈值" prop="threshold">
          <div class="threshold-input">
            <el-input-number
              v-model="ruleForm.threshold"
              :min="0"
              :max="100"
              :step="1"
              controls-position="right"
              style="width: 150px"
            />
            <span class="threshold-unit">%</span>
          </div>
        </el-form-item>

        <el-form-item label="持续时间" prop="duration">
          <div class="duration-input">
            <el-input-number
              v-model="ruleForm.duration"
              :min="1"
              :max="3600"
              :step="1"
              controls-position="right"
              style="width: 150px"
            />
            <span class="duration-unit">秒</span>
          </div>
          <div class="form-tip">
            持续超过阈值的时间后触发告警
          </div>
        </el-form-item>

        <el-form-item label="严重程度" prop="severity">
          <el-radio-group v-model="ruleForm.severity">
            <el-radio label="info">信息</el-radio>
            <el-radio label="warning">警告</el-radio>
            <el-radio label="critical">严重</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="告警动作">
          <el-checkbox-group v-model="selectedActions">
            <el-checkbox label="email">邮件通知</el-checkbox>
            <el-checkbox label="notification">系统通知</el-checkbox>
            <el-checkbox label="webhook">Webhook</el-checkbox>
            <el-checkbox label="log">记录日志</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="ruleForm.enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <div class="dialog-footer">
          <el-button @click="showRuleDialog = false">取消</el-button>
          <el-button type="primary" @click="saveAlertRule" :loading="saving">
            保存
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- Alert Preview -->
    <div class="alert-preview">
      <h4>告警预览</h4>
      <div class="preview-content">
        <el-alert
          v-if="previewAlert"
          :title="previewAlert.title"
          :type="previewAlert.type"
          :description="previewAlert.description"
          :closable="false"
          show-icon
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useQuotaStore } from '@/stores/quota'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Warning, InfoFilled, WarningFilled, CircleClose } from '@element-plus/icons-vue'
import type { QuotaAlert } from '@/types/quota'

const quotaStore = useQuotaStore()

const alertRules = ref<QuotaAlert[]>([
  {
    id: 1,
    name: '用户配额超限告警',
    type: 'user_quota',
    condition: '>=',
    threshold: 90,
    duration: 300,
    severity: 'warning',
    enabled: true,
    actions: JSON.stringify([
      { type: 'email', enabled: true },
      { type: 'notification', enabled: true }
    ])
  },
  {
    id: 2,
    name: '磁盘空间不足告警',
    type: 'disk_space',
    condition: '>=',
    threshold: 95,
    duration: 600,
    severity: 'critical',
    enabled: true,
    actions: JSON.stringify([
      { type: 'email', enabled: true },
      { type: 'notification', enabled: true },
      { type: 'webhook', enabled: false }
    ])
  }
])

const showRuleDialog = ref(false)
const editingRule = ref<QuotaAlert | null>(null)
const saving = ref(false)
const ruleFormRef = ref()
const selectedActions = ref(['email', 'notification'])

const ruleForm = ref({
  name: '',
  type: 'user_quota',
  condition: '>=',
  threshold: 80,
  duration: 300,
  severity: 'warning',
  enabled: true
})

const ruleFormRules = {
  name: [
    { required: true, message: '请输入规则名称', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择告警类型', trigger: 'change' }
  ],
  condition: [
    { required: true, message: '请选择触发条件', trigger: 'change' }
  ],
  threshold: [
    { required: true, message: '请输入阈值', trigger: 'blur' }
  ],
  duration: [
    { required: true, message: '请输入持续时间', trigger: 'blur' }
  ]
}

const previewAlert = computed(() => {
  if (alertRules.value.length === 0) return null
  const rule = alertRules.value[0]
  return {
    title: rule.name,
    type: rule.severity === 'critical' ? 'error' : rule.severity === 'warning' ? 'warning' : 'info',
    description: `当${rule.type === 'user_quota' ? '用户配额' : '磁盘空间'}使用率 ${rule.condition} ${rule.threshold}% 并持续 ${rule.duration} 秒时触发告警`
  }
})

const getAlertIcon = (type: string) => {
  switch (type) {
    case 'user_quota':
    case 'group_quota':
      return Warning
    case 'disk_space':
      return InfoFilled
    case 'storage_pool':
      return WarningFilled
    default:
      return Warning
  }
}

const getSeverityType = (severity: string): string => {
  switch (severity) {
    case 'critical':
      return 'danger'
    case 'warning':
      return 'warning'
    case 'info':
      return 'info'
    default:
      return ''
  }
}

const getSeverityLabel = (severity: string): string => {
  switch (severity) {
    case 'critical':
      return '严重'
    case 'warning':
      return '警告'
    case 'info':
      return '信息'
    default:
      return ''
  }
}

const formatCondition = (rule: QuotaAlert): string => {
  const conditions: Record<string, string> = {
    '>': '大于',
    '<': '小于',
    '>=': '大于等于',
    '<=': '小于等于',
    '=': '等于'
  }
  return conditions[rule.condition] || rule.condition
}

const parseActions = (actionsJson: string) => {
  try {
    return JSON.parse(actionsJson)
  } catch {
    return []
  }
}

const getActionLabel = (actionType: string): string => {
  const labels: Record<string, string> = {
    email: '邮件',
    notification: '通知',
    webhook: 'Webhook',
    log: '日志'
  }
  return labels[actionType] || actionType
}

const createAlertRule = () => {
  editingRule.value = null
  ruleForm.value = {
    name: '',
    type: 'user_quota',
    condition: '>=',
    threshold: 80,
    duration: 300,
    severity: 'warning',
    enabled: true
  }
  selectedActions.value = ['email', 'notification']
  showRuleDialog.value = true
}

const editAlertRule = (rule: QuotaAlert) => {
  editingRule.value = rule
  ruleForm.value = {
    name: rule.name,
    type: rule.type,
    condition: rule.condition,
    threshold: rule.threshold,
    duration: rule.duration,
    severity: rule.severity,
    enabled: rule.enabled
  }
  const actions = parseActions(rule.actions)
  selectedActions.value = actions.filter((a: any) => a.enabled).map((a: any) => a.type)
  showRuleDialog.value = true
}

const saveAlertRule = async () => {
  try {
    await ruleFormRef.value?.validate()

    saving.value = true

    const actions = selectedActions.value.map(type => ({
      type,
      enabled: true
    }))

    const newRule: QuotaAlert = {
      id: editingRule.value ? editingRule.value.id : Date.now(),
      ...ruleForm.value,
      actions: JSON.stringify(actions)
    }

    if (editingRule.value) {
      const index = alertRules.value.findIndex(r => r.id === editingRule.value!.id)
      if (index !== -1) {
        alertRules.value[index] = newRule
      }
      ElMessage.success('告警规则已更新')
    } else {
      alertRules.value.push(newRule)
      ElMessage.success('告警规则已创建')
    }

    showRuleDialog.value = false
  } catch (error) {
    ElMessage.error('保存告警规则失败')
  } finally {
    saving.value = false
  }
}

const deleteAlertRule = (rule: QuotaAlert) => {
  ElMessageBox.confirm(
    `确定要删除告警规则 "${rule.name}" 吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    const index = alertRules.value.findIndex(r => r.id === rule.id)
    if (index !== -1) {
      alertRules.value.splice(index, 1)
      ElMessage.success('告警规则已删除')
    }
  }).catch(() => {
    // User cancelled
  })
}

const toggleAlertRule = (rule: QuotaAlert) => {
  ElMessage.success(rule.enabled ? '告警规则已启用' : '告警规则已禁用')
}
</script>

<style scoped lang="scss">
.alert-settings {
  .settings-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h3 {
      margin: 0;
      color: #303133;
      font-size: 16px;
    }
  }

  .alert-rules-list {
    .alert-rule-item {
      background: white;
      padding: 20px;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      margin-bottom: 15px;
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      transition: all 0.3s ease;

      &.disabled {
        opacity: 0.6;
        background: #f5f7fa;
      }

      .rule-info {
        flex: 1;

        .rule-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 12px;

          .rule-name {
            font-size: 14px;
            font-weight: 500;
            color: #303133;
            display: flex;
            align-items: center;
            gap: 8px;
          }
        }

        .rule-details {
          display: flex;
          gap: 20px;
          margin-bottom: 10px;
          font-size: 12px;

          .detail-label {
            color: #909399;
            margin-right: 4px;
          }

          .detail-value {
            color: #606266;
            font-weight: 500;
          }
        }

        .rule-actions {
          display: flex;
          align-items: center;
          gap: 10px;
          font-size: 12px;

          .action-label {
            color: #909399;
          }

          .action-tags {
            display: flex;
            flex-wrap: wrap;
          }
        }
      }

      .rule-operations {
        display: flex;
        flex-direction: column;
        gap: 8px;
        margin-left: 20px;
      }
    }
  }

  .alert-preview {
    margin-top: 20px;
    background: white;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    h4 {
      margin: 0 0 15px;
      color: #303133;
      font-size: 14px;
    }

    .preview-content {
      margin-top: 15px;
    }
  }

  .threshold-input,
  .duration-input {
    display: flex;
    align-items: center;

    .threshold-unit,
    .duration-unit {
      margin-left: 8px;
      color: #909399;
      font-size: 12px;
    }
  }

  .form-tip {
    font-size: 12px;
    color: #909399;
    margin-top: 5px;
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }
}
</style>
