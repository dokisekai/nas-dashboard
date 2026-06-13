<template>
  <div class="firewall-manager">
    <div class="firewall-header">
      <div class="firewall-status">
        <div class="status-info">
          <h3>防火墙状态</h3>
          <p>{{ firewallEnabled ? '已启用 - 正在保护您的系统' : '已禁用 - 系统暴露在风险中' }}</p>
        </div>
        <button
          :class="['status-toggle', { active: firewallEnabled }]"
          @click="toggleFirewall"
        >
          <div class="toggle-slider"></div>
        </button>
      </div>

      <div class="firewall-actions">
        <button class="btn-primary" @click="showAddRuleModal = true">
          <PlusIcon class="w-4 h-4" />
          添加规则
        </button>
        <button class="btn-secondary" @click="applyRules" :disabled="applying">
          <ArrowPathIcon :class="['w-4 h-4', { spinning: applying }]" />
          应用更改
        </button>
      </div>
    </div>

    <div class="rules-table-container">
      <table class="rules-table">
        <thead>
          <tr>
            <th>名称</th>
            <th>动作</th>
            <th>协议</th>
            <th>端口</th>
            <th>来源IP</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="rule in rules" :key="rule.id">
            <td>
              <div class="rule-name">
                <strong>{{ rule.name }}</strong>
                <span v-if="rule.description" class="rule-desc">{{ rule.description }}</span>
              </div>
            </td>
            <td>
              <span :class="['badge', rule.action === 'allow' ? 'badge-success' : 'badge-danger']">
                {{ rule.action === 'allow' ? '允许' : '拒绝' }}
              </span>
            </td>
            <td>{{ rule.protocol.toUpperCase() }}</td>
            <td>{{ rule.port || '所有' }}</td>
            <td>{{ rule.sourceIp || '所有' }}</td>
            <td>
              <button
                :class="['rule-status-btn', { enabled: rule.enabled }]"
                @click="toggleRuleStatus(rule)"
              >
                {{ rule.enabled ? '启用中' : '已禁用' }}
              </button>
            </td>
            <td>
              <div class="row-actions">
                <button class="icon-btn" @click="editRule(rule)" title="编辑">
                  <PencilIcon class="w-4 h-4" />
                </button>
                <button class="icon-btn danger" @click="deleteRule(rule.id)" title="删除">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </td>
          </tr>
          <tr v-if="rules.length === 0">
            <td colspan="7" class="empty-row">暂无防火墙规则</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 添加/编辑规则模态框 -->
    <div v-if="showAddRuleModal || editingRule" class="modal-overlay" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ editingRule ? '编辑规则' : '添加规则' }}</h3>
          <button class="close-btn" @click="closeModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="submitRule" class="modal-body">
          <div class="form-group">
            <label>规则名称</label>
            <input v-model="ruleForm.name" type="text" placeholder="例如: 允许SSH访问" required />
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>动作</label>
              <select v-model="ruleForm.action">
                <option value="allow">允许</option>
                <option value="deny">拒绝</option>
              </select>
            </div>
            <div class="form-group">
              <label>协议</label>
              <select v-model="ruleForm.protocol">
                <option value="tcp">TCP</option>
                <option value="udp">UDP</option>
                <option value="both">TCP/UDP</option>
              </select>
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>端口</label>
              <input v-model="ruleForm.port" type="text" placeholder="例如: 80, 443, 20-21" />
            </div>
            <div class="form-group">
              <label>来源IP</label>
              <input v-model="ruleForm.sourceIp" type="text" placeholder="例如: any, 192.168.1.0/24" />
            </div>
          </div>

          <div class="form-group">
            <label>描述</label>
            <textarea v-model="ruleForm.description" rows="2" placeholder="可选描述..."></textarea>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="closeModal">取消</button>
            <button type="submit" class="btn-primary">确定</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  PlusIcon,
  ArrowPathIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import { useControlPanelStore } from '../../stores/controlPanel'

const controlPanelStore = useControlPanelStore()

// 状态
const rules = ref<any[]>([])
const applying = ref(false)
const showAddRuleModal = ref(false)
const editingRule = ref<any>(null)

const ruleForm = reactive({
  name: '',
  action: 'allow',
  protocol: 'tcp',
  port: '',
  sourceIp: 'any',
  description: '',
  enabled: true
})

const firewallEnabled = computed(() => {
  return controlPanelStore.settings['security.firewall.enabled'] === true || 
         controlPanelStore.settings['security.firewall.enabled'] === 'true'
})

// 方法
const loadRules = async () => {
  try {
    const response = await fetch('/api/security/firewall/rules', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.ok) {
      rules.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to load firewall rules:', error)
  }
}

const toggleFirewall = async () => {
  const newState = !firewallEnabled.value
  try {
    await controlPanelStore.updateSetting('security.firewall.enabled', newState)
    ElMessage.success(newState ? '防火墙已启用' : '防火墙已禁用')
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

const toggleRuleStatus = async (rule: any) => {
  try {
    const response = await fetch(`/api/security/firewall/rules/${rule.id}`, {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ ...rule, enabled: !rule.enabled })
    })
    if (response.ok) {
      rule.enabled = !rule.enabled
      ElMessage.success(rule.enabled ? '规则已启用' : '规则已禁用')
    }
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const editRule = (rule: any) => {
  editingRule.value = rule
  Object.assign(ruleForm, rule)
}

const closeModal = () => {
  showAddRuleModal.value = false
  editingRule.value = null
  // 重置表单
  Object.assign(ruleForm, {
    name: '',
    action: 'allow',
    protocol: 'tcp',
    port: '',
    sourceIp: 'any',
    description: '',
    enabled: true
  })
}

const submitRule = async () => {
  try {
    const url = editingRule.value 
      ? `/api/security/firewall/rules/${editingRule.value.id}` 
      : '/api/security/firewall/rules'
    
    const method = editingRule.value ? 'PUT' : 'POST'

    const response = await fetch(url, {
      method,
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(ruleForm)
    })

    if (response.ok) {
      ElMessage.success(editingRule.value ? '规则已更新' : '规则已添加')
      closeModal()
      loadRules()
    }
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

const deleteRule = (id: number) => {
  ElMessageBox.confirm('确定要删除这条规则吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      const response = await fetch(`/api/security/firewall/rules/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`
        }
      })
      if (response.ok) {
        ElMessage.success('规则已删除')
        loadRules()
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  })
}

const applyRules = async () => {
  applying.value = true
  try {
    const response = await fetch('/api/security/firewall/apply', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (response.ok) {
      ElMessage.success('防火墙配置已生效')
    }
  } catch (error) {
    ElMessage.error('应用失败')
  } finally {
    applying.value = false
  }
}

onMounted(() => {
  loadRules()
})
</script>

<style scoped lang="scss">
.firewall-manager {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.firewall-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
  padding: 20px;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.firewall-status {
  display: flex;
  align-items: center;
  gap: 24px;
}

.status-info {
  h3 {
    font-size: 16px;
    font-weight: 600;
    margin: 0 0 4px;
  }
  p {
    font-size: 13px;
    color: #6b7280;
    margin: 0;
  }
}

.status-toggle {
  width: 52px;
  height: 28px;
  background: #d1d5db;
  border-radius: 14px;
  position: relative;
  cursor: pointer;
  transition: all 0.3s;
  border: none;

  &.active {
    background: #10b981;
    .toggle-slider {
      transform: translateX(24px);
    }
  }
}

.toggle-slider {
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 4px;
  left: 4px;
  transition: transform 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.firewall-actions {
  display: flex;
  gap: 12px;
}

.rules-table-container {
  background: white;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.rules-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;

  th {
    background: #f9fafb;
    padding: 12px 16px;
    font-size: 13px;
    font-weight: 600;
    color: #4b5563;
    border-bottom: 1px solid #e5e7eb;
  }

  td {
    padding: 16px;
    border-bottom: 1px solid #f3f4f6;
    font-size: 14px;
    color: #1f2937;
  }
}

.rule-name {
  display: flex;
  flex-direction: column;
  gap: 2px;
  .rule-desc {
    font-size: 12px;
    color: #6b7280;
  }
}

.badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.badge-success {
  background: #d1fae5;
  color: #065f46;
}

.badge-danger {
  background: #fee2e2;
  color: #991b1b;
}

.rule-status-btn {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  border: 1px solid #d1d5db;
  background: white;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;

  &.enabled {
    border-color: #10b981;
    color: #10b981;
    background: #f0fdf4;
  }
}

.row-actions {
  display: flex;
  gap: 8px;
}

.icon-btn {
  padding: 6px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
  background: white;
  color: #6b7280;
  cursor: pointer;
  &:hover {
    background: #f9fafb;
    color: #2563eb;
  }
  &.danger:hover {
    color: #ef4444;
    border-color: #fee2e2;
    background: #fef2f2;
  }
}

.empty-row {
  text-align: center;
  color: #9ca3af;
  padding: 40px !important;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 500px;
  max-width: 90%;
}

.modal-header {
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
  h3 { margin: 0; font-size: 18px; }
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #9ca3af;
  &:hover { color: #1f2937; }
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  label { font-size: 13px; font-weight: 500; color: #374151; }
  input, select, textarea {
    padding: 8px 12px;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 14px;
    &:focus {
      outline: none;
      border-color: #2563eb;
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
    }
  }
}

.modal-footer {
  padding-top: 8px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-primary {
  background: #3b82f6;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  &:hover { background: #2563eb; }
}

.btn-secondary {
  background: white;
  border: 1px solid #d1d5db;
  padding: 8px 16px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  &:hover { background: #f9fafb; }
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>