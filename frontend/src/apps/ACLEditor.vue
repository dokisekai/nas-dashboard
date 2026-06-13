<template>
  <div class="acl-editor">
    <!-- 头部 -->
    <div class="acl-header">
      <div class="header-left">
        <h1>访问控制列表 (ACL) 编辑器</h1>
        <p class="subtitle">高级文件和目录权限管理</p>
      </div>
      <div class="header-actions">
        <button class="action-btn" @click="showHelp = true">
          <QuestionMarkCircleIcon class="w-4 h-4" />
          帮助
        </button>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="acl-content">
      <!-- 路径选择 -->
      <div class="path-section">
        <div class="path-input-group">
          <input
            type="text"
            v-model="currentPath"
            placeholder="输入文件或目录路径，例如：/mnt/data/share"
            @keyup.enter="loadACL"
          />
          <button class="action-btn primary" @click="loadACL">
            <MagnifyingGlassIcon class="w-4 h-4" />
            查看权限
          </button>
          <button class="action-btn" @click="browsePath">
            <FolderIcon class="w-4 h-4" />
            浏览
          </button>
        </div>

        <div class="quick-paths">
          <span class="quick-label">常用路径：</span>
          <button
            v-for="path in quickPaths"
            :key="path"
            class="quick-path-btn"
            @click="selectQuickPath(path)"
          >
            {{ path }}
          </button>
        </div>
      </div>

      <!-- 文件信息 -->
      <div v-if="fileInfo && !loading" class="file-info-section">
        <h3>文件信息</h3>
        <div class="info-grid">
          <div class="info-item">
            <label>路径</label>
            <span class="info-value">{{ fileInfo.path }}</span>
          </div>
          <div class="info-item">
            <label>类型</label>
            <span class="info-value">
              <FolderIcon v-if="fileInfo.isDir" class="w-4 h-4" />
              <DocumentIcon v-else class="w-4 h-4" />
              {{ fileInfo.isDir ? '目录' : '文件' }}
            </span>
          </div>
          <div class="info-item">
            <label>所有者</label>
            <span class="info-value">{{ fileInfo.owner }}</span>
          </div>
          <div class="info-item">
            <label>所属组</label>
            <span class="info-value">{{ fileInfo.group }}</span>
          </div>
          <div class="info-item">
            <label>权限模式</label>
            <span class="info-value perm-mode">{{ fileInfo.mode }}</span>
          </div>
          <div class="info-item">
            <label>大小</label>
            <span class="info-value">{{ formatFileSize(fileInfo.size) }}</span>
          </div>
          <div class="info-item">
            <label>修改时间</label>
            <span class="info-value">{{ fileInfo.modified }}</span>
          </div>
        </div>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p>加载ACL信息...</p>
      </div>

      <!-- ACL规则 -->
      <div v-if="aclRules.length > 0 && !loading" class="acl-section">
        <div class="section-header">
          <h3>ACL规则</h3>
          <button class="action-btn primary" @click="showAddRule = true">
            <PlusIcon class="w-4 h-4" />
            添加规则
          </button>
        </div>

        <!-- 当前ACL -->
        <div class="acl-group">
          <h4>当前ACL (访问ACL)</h4>
          <div class="acl-rules-table">
            <table>
              <thead>
                <tr>
                  <th>类型</th>
                  <th>名称</th>
                  <th>权限</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(rule, index) in currentACLRules"
                  :key="'current-' + index"
                  :class="{ 'default-rule': rule.default }"
                >
                  <td>
                    <span class="rule-type">{{ getRuleTypeLabel(rule.type) }}</span>
                  </td>
                  <td>
                    <span v-if="rule.name" class="rule-name">{{ rule.name }}</span>
                    <span v-else class="rule-name empty">-</span>
                  </td>
                  <td>
                    <span class="rule-perms">{{ rule.perms }}</span>
                    <div class="perm-breakdown">
                      <span
                        v-for="perm in rule.perms.split('')"
                        :key="perm"
                        class="perm-bit"
                        :class="{ active: hasPermission(perm) }"
                      >
                        {{ getPermLabel(perm) }}
                      </span>
                    </div>
                  </td>
                  <td>
                    <div class="rule-actions">
                      <button
                        class="icon-btn"
                        @click="editRule(rule, index, false)"
                        title="编辑"
                      >
                        <PencilIcon class="w-4 h-4" />
                      </button>
                      <button
                        class="icon-btn danger"
                        @click="removeRule(rule, index, false)"
                        title="删除"
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

        <!-- 默认ACL -->
        <div v-if="defaultACLRules.length > 0" class="acl-group">
          <h4>默认ACL (继承ACL)</h4>
          <div class="acl-info">
            <InformationCircleIcon class="w-5 h-5" />
            <p>默认ACL仅适用于目录，新创建的文件和子目录将继承这些权限。</p>
          </div>

          <div class="acl-rules-table">
            <table>
              <thead>
                <tr>
                  <th>类型</th>
                  <th>名称</th>
                  <th>权限</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="(rule, index) in defaultACLRules"
                  :key="'default-' + index"
                  class="default-rule"
                >
                  <td>
                    <span class="rule-type">{{ getRuleTypeLabel(rule.type) }}</span>
                    <span class="default-badge">默认</span>
                  </td>
                  <td>
                    <span v-if="rule.name" class="rule-name">{{ rule.name }}</span>
                    <span v-else class="rule-name empty">-</span>
                  </td>
                  <td>
                    <span class="rule-perms">{{ rule.perms }}</span>
                    <div class="perm-breakdown">
                      <span
                        v-for="perm in rule.perms.split('')"
                        :key="perm"
                        class="perm-bit"
                        :class="{ active: hasPermission(perm) }"
                      >
                        {{ getPermLabel(perm) }}
                      </span>
                    </div>
                  </td>
                  <td>
                    <div class="rule-actions">
                      <button
                        class="icon-btn"
                        @click="editRule(rule, index, true)"
                        title="编辑"
                      >
                        <PencilIcon class="w-4 h-4" />
                      </button>
                      <button
                        class="icon-btn danger"
                        @click="removeRule(rule, index, true)"
                        title="删除"
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

        <!-- 权限矩阵视图 -->
        <div class="acl-matrix-view">
          <h4>权限矩阵</h4>
          <div class="matrix-grid">
            <table>
              <thead>
                <tr>
                  <th>用户/组</th>
                  <th>读取</th>
                  <th>写入</th>
                  <th>执行</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="rule in aclRules"
                  :key="rule.type + '-' + rule.name"
                >
                  <td>
                    <div class="identity-cell">
                      <span class="identity-type">{{ getRuleTypeLabel(rule.type) }}</span>
                      <span v-if="rule.name" class="identity-name">{{ rule.name }}</span>
                    </div>
                  </td>
                  <td>
                    <div class="perm-cell">
                      <input
                        type="checkbox"
                        :checked="rule.perms.includes('r')"
                        disabled
                      />
                    </div>
                  </td>
                  <td>
                    <div class="perm-cell">
                      <input
                        type="checkbox"
                        :checked="rule.perms.includes('w')"
                        disabled
                      />
                    </div>
                  </td>
                  <td>
                    <div class="perm-cell">
                      <input
                        type="checkbox"
                        :checked="rule.perms.includes('x')"
                        disabled
                      />
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="aclRules.length === 0 && !loading && currentPath" class="empty-state">
        <ShieldCheckIcon class="w-12 h-12" />
        <p>此文件没有自定义ACL规则</p>
        <p class="hint">使用传统Unix权限进行访问控制</p>
        <button class="btn-primary" @click="showAddRule = true">
          添加ACL规则
        </button>
      </div>
    </div>

    <!-- 添加ACL规则对话框 -->
    <div v-if="showAddRule" class="modal-overlay" @click="showAddRule = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>添加ACL规则</h3>
          <button class="close-btn" @click="showAddRule = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="addACLRule" class="modal-body">
          <div class="form-group">
            <label>规则类型</label>
            <select v-model="newRule.type" required>
              <option value="user">用户</option>
              <option value="group">组</option>
              <option value="other">其他</option>
              <option value="mask">掩码</option>
            </select>
          </div>

          <div v-if="newRule.type !== 'other' && newRule.type !== 'mask'" class="form-group">
            <label>{{ newRule.type === 'user' ? '用户名' : '组名' }}</label>
            <input
              type="text"
              v-model="newRule.name"
              :placeholder="newRule.type === 'user' ? '输入用户名' : '输入组名'"
              required
            />
          </div>

          <div class="form-group">
            <label>权限</label>
            <div class="perm-selector">
              <label class="perm-checkbox">
                <input type="checkbox" v-model="newRule.perms" value="r" />
                <span>读取 (Read)</span>
              </label>
              <label class="perm-checkbox">
                <input type="checkbox" v-model="newRule.perms" value="w" />
                <span>写入 (Write)</span>
              </label>
              <label class="perm-checkbox">
                <input type="checkbox" v-model="newRule.perms" value="x" />
                <span>执行 (Execute)</span>
              </label>
            </div>
          </div>

          <div v-if="fileInfo?.isDir" class="form-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="newRule.isDefault" />
              <span>设为默认ACL</span>
              <small>默认ACL会被新创建的文件和子目录继承</small>
            </label>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showAddRule = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              添加规则
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 编辑ACL规则对话框 -->
    <div v-if="showEditRule" class="modal-overlay" @click="showEditRule = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑ACL规则</h3>
          <button class="close-btn" @click="showEditRule = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="updateACLRule" class="modal-body">
          <div class="form-group">
            <label>规则类型</label>
            <input type="text" :value="getRuleTypeLabel(editingRule.type)" disabled />
          </div>

          <div v-if="editingRule.name" class="form-group">
            <label>{{ editingRule.type === 'user' ? '用户名' : '组名' }}</label>
            <input type="text" :value="editingRule.name" disabled />
          </div>

          <div class="form-group">
            <label>权限</label>
            <div class="perm-selector">
              <label class="perm-checkbox">
                <input type="checkbox" v-model="editingRule.perms" value="r" />
                <span>读取 (Read)</span>
              </label>
              <label class="perm-checkbox">
                <input type="checkbox" v-model="editingRule.perms" value="w" />
                <span>写入 (Write)</span>
              </label>
              <label class="perm-checkbox">
                <input type="checkbox" v-model="editingRule.perms" value="x" />
                <span>执行 (Execute)</span>
              </label>
            </div>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showEditRule = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              保存更改
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 帮助对话框 -->
    <div v-if="showHelp" class="modal-overlay" @click="showHelp = false">
      <div class="modal-content large-modal" @click.stop>
        <div class="modal-header">
          <h3>ACL编辑器帮助</h3>
          <button class="close-btn" @click="showHelp = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="modal-body help-content">
          <div class="help-section">
            <h4>什么是ACL?</h4>
            <p>ACL (Access Control List) 提供了比传统Unix权限更精细的访问控制。它允许您为特定用户和组设置权限，而不仅仅是所有者、组和其他。</p>
          </div>

          <div class="help-section">
            <h4>ACL规则类型</h4>
            <ul>
              <li><strong>用户</strong>: 为特定用户设置权限</li>
              <li><strong>组</strong>: 为特定组设置权限</li>
              <li><strong>其他</strong>: 为所有其他用户设置权限</li>
              <li><strong>掩码</strong>: 限制所有用户和组的最大权限</li>
            </ul>
          </div>

          <div class="help-section">
            <h4>权限类型</h4>
            <ul>
              <li><strong>读取 (r)</strong>: 允许读取文件内容或列出目录</li>
              <li><strong>写入 (w)</strong>: 允许修改文件或创建/删除文件</li>
              <li><strong>执行 (x)</strong>: 允许执行文件或访问目录</li>
            </ul>
          </div>

          <div class="help-section">
            <h4>当前ACL vs 默认ACL</h4>
            <p><strong>当前ACL</strong>: 控制当前文件或目录的访问权限</p>
            <p><strong>默认ACL</strong>: 仅适用于目录，新创建的文件和子目录将继承这些权限</p>
          </div>

          <div class="help-section">
            <h4>使用示例</h4>
            <ul>
              <li>为用户 "john" 对 /data/project 目录添加读写权限:</li>
              <li>为组 "developers" 添加只读访问权限:</li>
              <li>设置默认ACL，使新文件继承特定权限:</li>
            </ul>
          </div>

          <div class="help-section warning">
            <h4>注意事项</h4>
            <ul>
              <li>修改ACL可能影响系统安全性，请谨慎操作</li>
              <li>确保文件系统支持ACL (如ext4, xfs, btrfs)</li>
              <li>建议先在测试环境验证ACL规则</li>
            </ul>
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
  QuestionMarkCircleIcon,
  MagnifyingGlassIcon,
  FolderIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
  InformationCircleIcon,
  ShieldCheckIcon,
  DocumentIcon
} from '@heroicons/vue/24/outline'

// 状态
const loading = ref(false)
const currentPath = ref('/mnt/data')
const fileInfo = ref<any>(null)
const aclRules = ref<any[]>([])

// 对话框状态
const showAddRule = ref(false)
const showEditRule = ref(false)
const showHelp = ref(false)

// 表单数据
const newRule = reactive({
  type: 'user',
  name: '',
  perms: [] as string[],
  isDefault: false
})

const editingRule = reactive({
  type: '',
  name: '',
  perms: [] as string[],
  isDefault: false,
  index: 0
})

// 常用路径
const quickPaths = ref([
  '/mnt/data',
  '/home',
  '/var',
  '/tmp'
])

// 计算属性
const currentACLRules = computed(() => {
  return aclRules.value.filter(rule => !rule.default)
})

const defaultACLRules = computed(() => {
  return aclRules.value.filter(rule => rule.default)
})

// 方法
const loadACL = async () => {
  if (!currentPath.value) {
    ElMessage.warning('请输入文件或目录路径')
    return
  }

  loading.value = true
  try {
    // 获取文件信息
    const permResponse = await fetch(`/api/permissions/files/permissions?path=${encodeURIComponent(currentPath.value)}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (permResponse.ok) {
      const permData = await permResponse.json()
      fileInfo.value = permData.file
    }

    // 获取ACL
    const aclResponse = await fetch(`/api/permissions/files/acl?path=${encodeURIComponent(currentPath.value)}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (aclResponse.ok) {
      const aclData = await aclResponse.json()
      aclRules.value = aclData.acl || []
      ElMessage.success('ACL信息加载成功')
    } else {
      throw new Error('获取ACL失败')
    }
  } catch (error) {
    console.error('Failed to load ACL:', error)
    ElMessage.error('加载ACL信息失败')
  } finally {
    loading.value = false
  }
}

const browsePath = () => {
  ElMessage.info('请输入完整的文件或目录路径')
}

const selectQuickPath = (path: string) => {
  currentPath.value = path
  loadACL()
}

const getRuleTypeLabel = (type: string) => {
  switch (type) {
    case 'user': return '用户'
    case 'group': return '组'
    case 'other': return '其他'
    case 'mask': return '掩码'
    default: return type
  }
}

const getPermLabel = (perm: string) => {
  switch (perm) {
    case 'r': return '读'
    case 'w': return '写'
    case 'x': return '执行'
    case '-': return '无'
    default: return perm
  }
}

const hasPermission = (perm: string) => {
  return perm !== '-'
}

const addACLRule = async () => {
  if (newRule.perms.length === 0) {
    ElMessage.warning('请选择至少一个权限')
    return
  }

  const rule = {
    type: newRule.type,
    name: newRule.name,
    perms: newRule.perms.join(''),
    default: newRule.isDefault
  }

  try {
    const response = await fetch('/api/permissions/files/acl', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        rules: [rule]
      })
    })

    if (response.ok) {
      ElMessage.success('ACL规则添加成功')
      showAddRule.value = false
      resetNewRule()
      await loadACL()
    } else {
      throw new Error('添加ACL规则失败')
    }
  } catch (error) {
    console.error('Failed to add ACL rule:', error)
    ElMessage.error('添加ACL规则失败')
  }
}

const editRule = (rule: any, index: number, isDefault: boolean) => {
  editingRule.type = rule.type
  editingRule.name = rule.name
  editingRule.perms = rule.perms.split('')
  editingRule.isDefault = isDefault
  editingRule.index = index
  showEditRule.value = true
}

const updateACLRule = async () => {
  if (editingRule.perms.length === 0) {
    ElMessage.warning('请选择至少一个权限')
    return
  }

  // 先删除旧规则，再添加新规则
  const oldRule = aclRules.value[editingRule.index]
  const rulesToRemove = [oldRule]

  const newRule = {
    type: editingRule.type,
    name: editingRule.name,
    perms: editingRule.perms.join(''),
    default: editingRule.isDefault
  }

  try {
    // 删除旧规则
    await fetch('/api/permissions/files/acl', {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        rules: rulesToRemove
      })
    })

    // 添加新规则
    const response = await fetch('/api/permissions/files/acl', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        rules: [newRule]
      })
    })

    if (response.ok) {
      ElMessage.success('ACL规则更新成功')
      showEditRule.value = false
      await loadACL()
    } else {
      throw new Error('更新ACL规则失败')
    }
  } catch (error) {
    console.error('Failed to update ACL rule:', error)
    ElMessage.error('更新ACL规则失败')
  }
}

const removeRule = async (rule: any, index: number, isDefault: boolean) => {
  if (!confirm('确定要删除此ACL规则吗？')) {
    return
  }

  try {
    const response = await fetch('/api/permissions/files/acl', {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        rules: [rule]
      })
    })

    if (response.ok) {
      ElMessage.success('ACL规则删除成功')
      await loadACL()
    } else {
      throw new Error('删除ACL规则失败')
    }
  } catch (error) {
    console.error('Failed to remove ACL rule:', error)
    ElMessage.error('删除ACL规则失败')
  }
}

const resetNewRule = () => {
  newRule.type = 'user'
  newRule.name = ''
  newRule.perms = []
  newRule.isDefault = false
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
  // 自动加载默认路径
  loadACL()
})
</script>

<style scoped lang="scss">
.acl-editor {
  width: 100%;
  padding: 24px;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.acl-header {
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

.acl-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.path-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.path-input-group {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;

  input {
    flex: 1;
    padding: 12px;
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    font-size: 14px;
    font-family: monospace;
    transition: all 0.2s;

    &:focus {
      outline: none;
      border-color: #667eea;
      box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
    }
  }
}

.quick-paths {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;

  .quick-label {
    font-size: 13px;
    font-weight: 500;
    color: #6b7280;
  }

  .quick-path-btn {
    padding: 6px 12px;
    background: rgba(102, 126, 234, 0.05);
    color: #667eea;
    border: 1px solid rgba(102, 126, 234, 0.1);
    border-radius: 6px;
    font-size: 12px;
    font-family: monospace;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: rgba(102, 126, 234, 0.1);
      border-color: #667eea;
    }
  }
}

.file-info-section {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 20px 0;
  }
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  label {
    font-size: 12px;
    font-weight: 500;
    color: #6b7280;
  }

  .info-value {
    font-size: 14px;
    color: #1f2937;
    display: flex;
    align-items: center;
    gap: 6px;

    &.perm-mode {
      font-family: monospace;
      font-weight: 600;
      color: #667eea;
    }
  }
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

.acl-section {
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

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.acl-group {
  margin-bottom: 32px;

  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #374151;
    margin: 0 0 16px 0;
  }

  &:last-child {
    margin-bottom: 0;
  }
}

.acl-info {
  display: flex;
  gap: 12px;
  padding: 12px;
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid rgba(59, 130, 246, 0.1);
  border-radius: 8px;
  color: #2563eb;
  margin-bottom: 16px;

  svg {
    flex-shrink: 0;
  }

  p {
    font-size: 13px;
    margin: 0;
    line-height: 1.4;
  }
}

.acl-rules-table {
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

      &.default-rule {
        background: rgba(251, 191, 36, 0.05);
      }
    }
  }
}

.rule-type {
  padding: 4px 8px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.rule-name {
  font-weight: 500;
  color: #1f2937;

  &.empty {
    color: #9ca3af;
    font-style: italic;
  }
}

.default-badge {
  margin-left: 8px;
  padding: 4px 8px;
  background: rgba(251, 191, 36, 0.1);
  color: #b45309;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.rule-perms {
  font-family: monospace;
  font-weight: 600;
  color: #667eea;
  display: block;
  margin-bottom: 4px;
}

.perm-breakdown {
  display: flex;
  gap: 4px;
  margin-top: 4px;
}

.perm-bit {
  padding: 2px 6px;
  background: #f3f4f6;
  color: #9ca3af;
  border-radius: 3px;
  font-size: 11px;
  font-family: monospace;

  &.active {
    background: rgba(16, 185, 129, 0.1);
    color: #065f46;
  }
}

.rule-actions {
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

.acl-matrix-view {
  margin-top: 24px;

  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #374151;
    margin: 0 0 16px 0;
  }
}

.matrix-grid {
  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead {
    background: linear-gradient(to bottom, #f8fafc, #e2e8f0);

    th {
      padding: 12px;
      text-align: center;
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
        text-align: center;
      }

      &:first-child td {
        text-align: left;
      }
    }
  }
}

.identity-cell {
  display: flex;
  align-items: center;
  gap: 8px;
  justify-content: flex-start;

  .identity-type {
    padding: 4px 8px;
    background: rgba(102, 126, 234, 0.1);
    color: #667eea;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 500;
  }

  .identity-name {
    font-weight: 500;
    color: #1f2937;
  }
}

.perm-cell {
  display: flex;
  align-items: center;
  justify-content: center;

  input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: not-allowed;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  color: #9ca3af;

  svg {
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    font-size: 16px;
    margin: 8px 0;

    &.hint {
      font-size: 14px;
      color: #6b7280;
    }
  }

  .btn-primary {
    margin-top: 20px;
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
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);

  &.large-modal {
    max-width: 700px;
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
  gap: 8px;

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

.perm-selector {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.perm-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: rgba(102, 126, 234, 0.05);
    border-color: #667eea;
  }

  input[type="checkbox"] {
    width: 18px;
    height: 18px;
  }

  span {
    font-size: 14px;
    color: #1f2937;
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

.help-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.help-section {
  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 12px 0;
  }

  p {
    font-size: 14px;
    color: #4b5563;
    line-height: 1.6;
    margin: 0 0 8px 0;
  }

  ul {
    margin: 0;
    padding-left: 20px;
    color: #4b5563;
    font-size: 14px;
    line-height: 1.6;

    li {
      margin-bottom: 6px;

      strong {
        color: #1f2937;
      }
    }
  }

  &.warning {
    padding: 16px;
    background: rgba(251, 191, 36, 0.05);
    border: 1px solid rgba(251, 191, 36, 0.2);
    border-radius: 8px;
    color: #b45309;

    h4 {
      color: #b45309;
    }

    ul {
      color: #b45309;
    }
  }
}
</style>