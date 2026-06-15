<template>
  <div class="notification-settings">
    <div class="settings-header">
      <h4>通知规则设置</h4>
      <el-button type="primary" size="small" @click="createNewRule">
        <PlusIcon class="w-4 h-4 mr-1" />
        新建规则
      </el-button>
    </div>

    <div class="rules-list">
      <div
        v-for="rule in rules"
        :key="rule.id"
        class="rule-item"
        :class="{ disabled: !rule.enabled }"
      >
        <div class="rule-info">
          <div class="rule-header">
            <h5>{{ rule.name }}</h5>
            <el-switch
              :model-value="rule.enabled"
              @change="toggleRuleEnabled(rule)"
              size="small"
            />
          </div>
          <div class="rule-details">
            <span class="event-type">{{ getEventTypeName(rule.eventType) }}</span>
            <span class="cooldown">冷却: {{ formatCooldown(rule.cooldown) }}</span>
            <span class="last-triggered">
              最后触发: {{ rule.lastTriggered ? formatTime(rule.lastTriggered) : '从未' }}
            </span>
          </div>
        </div>

        <div class="rule-actions">
          <el-button size="small" @click="testRule(rule)">
            <PlayIcon class="w-4 h-4 mr-1" />
            测试
          </el-button>
          <el-button size="small" @click="editRule(rule)">
            <PencilIcon class="w-4 h-4 mr-1" />
            编辑
          </el-button>
          <el-button size="small" type="danger" @click="deleteRule(rule)">
            <TrashIcon class="w-4 h-4 mr-1" />
            删除
          </el-button>
        </div>
      </div>
    </div>

    <!-- 规则编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="isEditing ? '编辑通知规则' : '新建通知规则'"
      width="700px"
    >
      <el-form :model="currentRule" label-width="120px">
        <el-form-item label="规则名称">
          <el-input v-model="currentRule.name" placeholder="输入规则名称" />
        </el-form-item>

        <el-form-item label="事件类型">
          <el-select v-model="currentRule.eventType" placeholder="选择事件类型">
            <el-option label="磁盘健康" value="disk_health" />
            <el-option label="磁盘温度" value="disk_temperature" />
            <el-option label="磁盘空间" value="disk_space" />
            <el-option label="内存告警" value="memory_alert" />
            <el-option label="CPU告警" value="cpu_alert" />
            <el-option label="负载告警" value="load_alert" />
            <el-option label="安全告警" value="security_alert" />
            <el-option label="备份完成" value="backup_complete" />
            <el-option label="备份失败" value="backup_failed" />
          </el-select>
        </el-form-item>

        <el-form-item label="触发条件">
          <div class="conditions-builder">
            <div
              v-for="(condition, index) in conditions"
              :key="index"
              class="condition-item"
            >
              <el-select v-model="condition.field" placeholder="字段" style="width: 150px">
                <el-option
                  v-for="field in getAvailableFields(currentRule.eventType)"
                  :key="field.key"
                  :label="field.label"
                  :value="field.key"
                />
              </el-select>
              <el-select v-model="condition.operator" placeholder="操作符" style="width: 100px">
                <el-option label="等于" value="eq" />
                <el-option label="不等于" value="ne" />
                <el-option label="大于" value="gt" />
                <el-option label="小于" value="lt" />
                <el-option label="大于等于" value="gte" />
                <el-option label="小于等于" value="lte" />
              </el-select>
              <el-input v-model="condition.value" placeholder="值" style="flex: 1" />
              <el-button size="small" type="danger" @click="removeCondition(index)">
                <TrashIcon class="w-4 h-4" />
              </el-button>
            </div>
            <el-button size="small" @click="addCondition">
              <PlusIcon class="w-4 h-4 mr-1" />
              添加条件
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="通知动作">
          <div class="actions-builder">
            <div
              v-for="(action, index) in actions"
              :key="index"
              class="action-item"
            >
              <el-select v-model="action.type" placeholder="动作类型" style="width: 150px">
                <el-option label="WebSocket推送" value="websocket" />
                <el-option label="邮件通知" value="email" />
                <el-option label="Webhook" value="webhook" />
              </el-select>
              <div class="action-config">
                <el-input
                  v-if="action.type === 'email'"
                  v-model="action.config.recipients"
                  placeholder="收件人邮箱，逗号分隔"
                />
                <el-input
                  v-if="action.type === 'webhook'"
                  v-model="action.config.url"
                  placeholder="Webhook URL"
                />
              </div>
              <el-button size="small" type="danger" @click="removeAction(index)">
                <TrashIcon class="w-4 h-4" />
              </el-button>
            </div>
            <el-button size="small" @click="addAction">
              <PlusIcon class="w-4 h-4 mr-1" />
              添加动作
            </el-button>
          </div>
        </el-form-item>

        <el-form-item label="冷却时间">
          <el-input-number
            v-model="currentRule.cooldown"
            :min="60"
            :max="86400"
            :step="60"
          />
          <span style="margin-left: 8px">秒</span>
        </el-form-item>

        <el-form-item label="启用规则">
          <el-switch v-model="currentRule.enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="saveRule">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { PlusIcon, PencilIcon, TrashIcon, PlayIcon } from '@heroicons/vue/24/outline'
import api from '../../api/client'

interface NotificationRule {
  id?: number
  name: string
  eventType: string
  conditions: string
  actions: string
  enabled: boolean
  cooldown: number
  lastTriggered?: string
}

interface Condition {
  field: string
  operator: string
  value: string
}

interface Action {
  type: string
  config: Record<string, string>
}

// 组件状态
const rules = ref<NotificationRule[]>([])
const editDialogVisible = ref(false)
const isEditing = ref(false)
const currentRule = reactive<NotificationRule>({
  name: '',
  eventType: '',
  conditions: '',
  actions: '',
  enabled: true,
  cooldown: 300
})

const conditions = ref<Condition[]>([])
const actions = ref<Action[]>([])

// 方法
const fetchRules = async () => {
  try {
    const data = await api.get('/api/notifications/rules')
    rules.value = data as NotificationRule[]
  } catch (error) {
    console.error('获取通知规则失败:', error)
  }
}

const createNewRule = () => {
  Object.assign(currentRule, {
    name: '',
    eventType: '',
    conditions: '',
    actions: '',
    enabled: true,
    cooldown: 300
  })
  conditions.value = []
  actions.value = []
  isEditing.value = false
  editDialogVisible.value = true
}

const editRule = (rule: NotificationRule) => {
  Object.assign(currentRule, rule)

  // 解析条件和动作
  try {
    if (rule.conditions) {
      conditions.value = JSON.parse(rule.conditions)
    }
    if (rule.actions) {
      actions.value = JSON.parse(rule.actions)
    }
  } catch (error) {
    console.error('解析规则配置失败:', error)
    conditions.value = []
    actions.value = []
  }

  isEditing.value = true
  editDialogVisible.value = true
}

const saveRule = async () => {
  // 构建规则对象
  const ruleData = {
    ...currentRule,
    conditions: JSON.stringify(conditions.value),
    actions: JSON.stringify(actions.value)
  }

  try {
    if (isEditing.value && currentRule.id) {
      await api.put(`/api/notifications/rules/${currentRule.id}`, ruleData)
    } else {
      await api.post('/api/notifications/rules', ruleData)
    }

    await fetchRules()
    editDialogVisible.value = false
  } catch (error) {
    console.error('保存规则失败:', error)
  }
}

const deleteRule = async (rule: NotificationRule) => {
  if (confirm(`确定要删除规则 "${rule.name}" 吗？`)) {
    try {
      await api.delete(`/api/notifications/rules/${rule.id}`)
      await fetchRules()
    } catch (error) {
      console.error('删除规则失败:', error)
    }
  }
}

const testRule = async (rule: NotificationRule) => {
  try {
    await api.post(`/api/notifications/rules/${rule.id}/test`)
    alert('测试通知已发送！请检查通知中心。')
  } catch (error) {
    console.error('测试规则失败:', error)
    alert('测试通知发送失败！')
  }
}

const toggleRuleEnabled = async (rule: NotificationRule) => {
  try {
    await api.put(`/api/notifications/rules/${rule.id}`, {
      ...rule,
      enabled: !rule.enabled
    })
    rule.enabled = !rule.enabled
  } catch (error) {
    console.error('切换规则状态失败:', error)
  }
}

const addCondition = () => {
  conditions.value.push({
    field: '',
    operator: 'eq',
    value: ''
  })
}

const removeCondition = (index: number) => {
  conditions.value.splice(index, 1)
}

const addAction = () => {
  actions.value.push({
    type: 'websocket',
    config: {}
  })
}

const removeAction = (index: number) => {
  actions.value.splice(index, 1)
}

const getEventTypeName = (eventType: string) => {
  const typeNames: Record<string, string> = {
    disk_health: '磁盘健康',
    disk_temperature: '磁盘温度',
    disk_space: '磁盘空间',
    memory_alert: '内存告警',
    cpu_alert: 'CPU告警',
    load_alert: '负载告警',
    security_alert: '安全告警',
    backup_complete: '备份完成',
    backup_failed: '备份失败'
  }
  return typeNames[eventType] || eventType
}

const getAvailableFields = (eventType: string) => {
  const fields: Record<string, Array<{ key: string; label: string }>> = {
    disk_health: [
      { key: 'disk', label: '磁盘设备' },
      { key: 'status', label: '健康状态' }
    ],
    disk_temperature: [
      { key: 'disk', label: '磁盘设备' },
      { key: 'temperature', label: '温度' }
    ],
    disk_space: [
      { key: 'disk', label: '磁盘设备' },
      { key: 'usage_percent', label: '使用率' }
    ],
    memory_alert: [
      { key: 'usage_percent', label: '使用率' }
    ],
    cpu_alert: [
      { key: 'usage_percent', label: '使用率' }
    ],
    load_alert: [
      { key: 'load1', label: '1分钟负载' },
      { key: 'cpu_count', label: 'CPU核心数' }
    ],
    security_alert: [
      { key: 'failed_count', label: '失败次数' }
    ]
  }
  return fields[eventType] || []
}

const formatCooldown = (seconds: number) => {
  if (seconds < 60) return `${seconds}秒`
  if (seconds < 3600) return `${Math.floor(seconds / 60)}分钟`
  return `${Math.floor(seconds / 3600)}小时`
}

const formatTime = (timestamp: string) => {
  const date = new Date(timestamp)
  return date.toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchRules()
})
</script>

<style scoped>
.notification-settings {
  padding: 20px;
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.settings-header h4 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.rule-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: var(--el-fill-color-blank);
  border: 1px solid var(--el-border-color-light);
  border-radius: 8px;
  transition: all 0.3s;
}

.rule-item:hover {
  border-color: var(--el-color-primary);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.rule-item.disabled {
  opacity: 0.6;
}

.rule-info {
  flex: 1;
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.rule-header h5 {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.rule-details {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.rule-actions {
  display: flex;
  gap: 8px;
}

.conditions-builder,
.actions-builder {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.condition-item,
.action-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.action-config {
  flex: 1;
}
</style>