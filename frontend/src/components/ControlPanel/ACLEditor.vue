<template>
  <div class="acl-editor">
    <div class="editor-header">
      <h3>权限管理 (ACL)</h3>
      <div class="editor-actions">
        <input
          v-model="searchPath"
          type="text"
          placeholder="输入路径..."
          class="path-input"
        />
        <button @click="browsePath" class="btn btn-secondary">
          <FolderIcon class="w-4 h-4" />
          浏览
        </button>
        <button @click="loadPermissions" class="btn btn-primary">
          <MagnifyingGlassIcon class="w-4 h-4" />
          查看权限
        </button>
      </div>
    </div>

    <!-- 当前权限信息 -->
    <div v-if="currentPath" class="current-permissions">
      <div class="perm-header">
        <h4>当前路径权限: {{ currentPath }}</h4>
        <div class="perm-badges">
          <span :class="['perm-badge', currentPerm.type]">
            {{ currentPerm.type === 'file' ? '文件' : '目录' }}
          </span>
          <span class="perm-badge">{{ currentPerm.permissions }}</span>
          <span class="perm-badge">{{ currentPerm.owner }}:{{ currentPerm.group }}</span>
        </div>
      </div>

      <div class="perm-details">
        <div class="perm-item">
          <span class="perm-label">所有者:</span>
          <span class="perm-value">{{ currentPerm.owner }}</span>
        </div>
        <div class="perm-item">
          <span class="perm-label">用户组:</span>
          <span class="perm-value">{{ currentPerm.group }}</span>
        </div>
        <div class="perm-item">
          <span class="perm-label">权限:</span>
          <div class="permission-breakdown">
            <div class="perm-group">
              <span class="perm-group-title">用户:</span>
              <span class="perm-chip user">{{ getPermString('user') }}</span>
            </div>
            <div class="perm-group">
              <span class="perm-group-title">组:</span>
              <span class="perm-chip group">{{ getPermString('group') }}</span>
            </div>
            <div class="perm-group">
              <span class="perm-group-title">其他:</span>
              <span class="perm-chip other">{{ getPermString('other') }}</span>
            </div>
          </div>
        </div>
        <div class="perm-item" v-if="currentPerm.acl">
          <span class="perm-label">ACL规则:</span>
          <div class="acl-rules">
            <div v-for="(rule, index) in currentPerm.acl" :key="index" class="acl-rule">
              <span class="rule-type">{{ rule.type }}::</span>
              <span class="rule-entity">{{ rule.entity }}</span>
              <span class="rule-permissions">{{ rule.permissions }}</span>
              <button @click="removeACLRule(index)" class="btn btn-sm btn-danger">
                <TrashIcon class="w-3 h-3" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 权限编辑器 -->
    <div class="permission-editor">
      <div class="editor-section">
        <h4>基础权限设置</h4>
        <div class="basic-perms">
          <div class="perm-row">
            <label>所有者:</label>
            <select v-model="permForm.owner" class="perm-select">
              <option v-for="user in availableUsers" :key="user" :value="user">
                {{ user }}
              </option>
            </select>
          </div>
          <div class="perm-row">
            <label>组:</label>
            <select v-model="permForm.group" class="perm-select">
              <option v-for="group in availableGroups" :key="group" :value="group">
                {{ group }}
              </option>
            </select>
          </div>
          <div class="perm-row">
            <label>权限:</label>
            <div class="perm-checks">
              <div class="perm-check-group">
                <span class="perm-check-title">读</span>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.userRead" />
                  <span>U</span>
                </label>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.groupRead" />
                  <span>G</span>
                </label>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.otherRead" />
                  <span>O</span>
                </label>
              </div>
              <div class="perm-check-group">
                <span class="perm-check-title">写</span>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.userWrite" />
                  <span>U</span>
                </label>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.groupWrite" />
                  <span>G</span>
                </label>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.otherWrite" />
                  <span>O</span>
                </label>
              </div>
              <div class="perm-check-group">
                <span class="perm-check-title">执行</span>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.userExec" />
                  <span>U</span>
                </label>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.groupExec" />
                  <span>G</span>
                </label>
                <label class="perm-check">
                  <input type="checkbox" v-model="permForm.otherExec" />
                  <span>O</span>
                </label>
              </div>
            </div>
          </div>
          <div class="perm-row">
            <label>特殊权限:</label>
            <div class="special-perms">
              <label class="perm-check">
                <input type="checkbox" v-model="permForm.setuid" />
                <span>SUID (4)</span>
              </label>
              <label class="perm-check">
                <input type="checkbox" v-model="permForm.setgid" />
                <span>SGID (2)</span>
              </label>
              <label class="perm-check">
                <input type="checkbox" v-model="permForm.sticky" />
                <span>Sticky (1)</span>
              </label>
            </div>
          </div>
        </div>
        <div class="perm-actions">
          <button @click="applyBasicPerms" class="btn btn-primary" :disabled="applying">
            {{ applying ? '应用中...' : '应用基础权限' }}
          </button>
          <button @click="resetPermForm" class="btn btn-secondary">
            重置
          </button>
        </div>
      </div>

      <div class="editor-section">
        <h4>ACL 规则管理</h4>
        <div class="acl-adder">
          <div class="acl-row">
            <select v-model="aclForm.type" class="acl-select">
              <option value="user">用户 (u)</option>
              <option value="group">组 (g)</option>
              <option value="other">其他 (o)</option>
              <option value="mask">掩码 (m)</option>
            </select>
            <select v-model="aclForm.entity" class="acl-select">
              <option value="">选择实体</option>
              <option v-for="entity in availableACLEntities" :key="entity" :value="entity">
                {{ entity }}
              </option>
            </select>
            <div class="acl-perms-input">
              <label class="perm-check">
                <input type="checkbox" v-model="aclForm.read" />
                <span>R</span>
              </label>
              <label class="perm-check">
                <input type="checkbox" v-model="aclForm.write" />
                <span>W</span>
              </label>
              <label class="perm-check">
                <input type="checkbox" v-model="aclForm.execute" />
                <span>X</span>
              </label>
            </div>
            <button @click="addACLRule" class="btn btn-success">
              <PlusIcon class="w-4 h-4" />
              添加规则
            </button>
          </div>
        </div>
        <div class="acl-actions">
          <button @click="applyACL" class="btn btn-primary" :disabled="applying">
            {{ applying ? '应用中...' : '应用 ACL' }}
          </button>
          <button @click="removeAllACL" class="btn btn-danger">
            移除所有 ACL
          </button>
        </div>
      </div>
    </div>

    <!-- 常用权限预设 -->
    <div class="permission-presets">
      <h4>常用权限预设</h4>
      <div class="presets-grid">
        <button @click="applyPreset('755')" class="preset-btn">
          <div class="preset-name">755</div>
          <div class="preset-desc">用户完全，组和其他读写执行</div>
        </button>
        <button @click="applyPreset('644')" class="preset-btn">
          <div class="preset-name">644</div>
          <div class="preset-desc">用户读写，组和其他只读</div>
        </button>
        <button @click="applyPreset('600')" class="preset-btn">
          <div class="preset-name">600</div>
          <div class="preset-desc">用户读写，组和其他无权限</div>
        </button>
        <button @click="applyPreset('777')" class="preset-btn">
          <div class="preset-name">777</div>
          <div class="preset-desc">所有人完全权限</div>
        </button>
        <button @click="applyPreset('2755')" class="preset-btn">
          <div class="preset-name">2755</div>
          <div class="preset-desc">SGID + 755</div>
        </button>
        <button @click="applyPreset('4755')" class="preset-btn">
          <div class="preset-name">4755</div>
          <div class="preset-desc">SUID + 755</div>
        </button>
      </div>
    </div>

    <!-- 权限浏览对话框 -->
    <div v-if="showBrowserModal" class="modal-overlay" @click.self="closeBrowser">
      <div class="modal browser-modal">
        <div class="modal-header">
          <h4>浏览文件系统</h4>
          <button @click="closeBrowser" class="btn btn-ghost">
            <XMarkIcon class="w-4 h-4" />
          </button>
        </div>
        <div class="modal-body">
          <div class="browser-nav">
            <button @click="goToParent" class="btn btn-sm btn-secondary" :disabled="!currentDirParent">
              <ArrowUturnLeftIcon class="w-4 h-4" />
              上级目录
            </button>
            <input
              v-model="browserPath"
              type="text"
              placeholder="路径"
              class="browser-path-input"
              @keyup.enter="navigateToPath"
            />
            <button @click="navigateToPath" class="btn btn-primary">
              转到
            </button>
          </div>
          <div class="browser-contents">
            <div
              v-for="item in browserContents"
              :key="item.name"
              class="browser-item"
              @click="enterDirectory(item)"
            >
              <FolderIcon v-if="item.type === 'dir'" class="w-5 h-5" />
              <DocumentIcon v-else class="w-5 h-5" />
              <span>{{ item.name }}</span>
              <span class="item-perms">{{ item.perms }}</span>
            </div>
            <div v-if="browserContents.length === 0" class="empty-browser">
              此目录为空或无法访问
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  FolderIcon,
  MagnifyingGlassIcon,
  TrashIcon,
  PlusIcon,
  XMarkIcon,
  ArrowUturnLeftIcon,
  DocumentIcon
} from '@heroicons/vue/24/outline'

// 状态
const searchPath = ref('/home/hserver')
const currentPath = ref('')
const currentPerm = ref<any>(null)
const applying = ref(false)
const showBrowserModal = ref(false)

// 可用的用户和组
const availableUsers = ref<string[]>(['root', 'hserver', 'admin'])
const availableGroups = ref<string[]>(['root', 'hserver', 'admin', 'users', 'docker', 'sudo'])

// 权限表单
const permForm = ref({
  owner: 'hserver',
  group: 'hserver',
  userRead: true,
  userWrite: true,
  userExec: true,
  groupRead: true,
  groupWrite: false,
  groupExec: true,
  otherRead: true,
  otherWrite: false,
  otherExec: true,
  setuid: false,
  setgid: false,
  sticky: false
})

// ACL 表单
const aclForm = ref({
  type: 'user',
  entity: '',
  read: false,
  write: false,
  execute: false
})

// 浏览器相关
const browserPath = ref('/home/hserver')
const browserContents = ref<any[]>([])

// 计算属性
const currentPermNumeric = computed(() => {
  let perm = 0
  if (permForm.value.setuid) perm += 0o4000
  if (permForm.value.setgid) perm += 0o2000
  if (permForm.value.sticky) perm += 0o1000

  if (permForm.value.userRead) perm += 0o400
  if (permForm.value.userWrite) perm += 0o200
  if (permForm.value.userExec) perm += 0o100

  if (permForm.value.groupRead) perm += 0o040
  if (permForm.value.groupWrite) perm += 0o020
  if (permForm.value.groupExec) perm += 0o010

  if (permForm.value.otherRead) perm += 0o004
  if (permForm.value.otherWrite) perm += 0o002
  if (permForm.value.otherExec) perm += 0o001

  return perm.toString(8)
})

const currentDirParent = computed(() => {
  const parts = browserPath.value.split('/').filter(Boolean)
  return parts.length > 0
})

const availableACLEntities = computed(() => {
  if (aclForm.value.type === 'user') {
    return availableUsers.value
  } else if (aclForm.value.type === 'group') {
    return availableGroups.value
  }
  return []
})

// 获取权限字符串
const getPermString = (type: 'user' | 'group' | 'other') => {
  const form = permForm.value
  let result = ''

  if (type === 'user') {
    result += form.userRead ? 'r' : '-'
    result += form.userWrite ? 'w' : '-'
    result += form.userExec ? 'x' : '-'
  } else if (type === 'group') {
    result += form.groupRead ? 'r' : '-'
    result += form.groupWrite ? 'w' : '-'
    result += form.groupExec ? 'x' : '-'
  } else {
    result += form.otherRead ? 'r' : '-'
    result += form.otherWrite ? 'w' : '-'
    result += form.otherExec ? 'x' : '-'
  }

  return result
}

// 加载权限信息
const loadPermissions = async () => {
  if (!searchPath.value) {
    alert('请输入路径')
    return
  }

  try {
    const token = localStorage.getItem('token')
    // 修复API端点：使用正确的后端路由
    const response = await fetch(`/api/permissions/files?path=${encodeURIComponent(searchPath.value)}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      const data = await response.json()
      currentPath.value = searchPath.value
      currentPerm.value = data

      // 更新表单
      permForm.value.owner = data.owner || 'hserver'
      permForm.value.group = data.group || 'hserver'

      // 解析权限
      const perms = data.permissions || 'rwxr-xr-x'
      permForm.value.userRead = perms.charAt(0) === 'r'
      permForm.value.userWrite = perms.charAt(1) === 'w'
      permForm.value.userExec = perms.charAt(2) === 'x'
      permForm.value.groupRead = perms.charAt(3) === 'r'
      permForm.value.groupWrite = perms.charAt(4) === 'w'
      permForm.value.groupExec = perms.charAt(5) === 'x'
      permForm.value.otherRead = perms.charAt(6) === 'r'
      permForm.value.otherWrite = perms.charAt(7) === 'w'
      permForm.value.otherExec = perms.charAt(8) === 'x'

      // 特殊权限
      permForm.value.setuid = data.setuid || false
      permForm.value.setgid = data.setgid || false
      permForm.value.sticky = data.sticky || false
    } else {
      alert('加载权限失败')
    }
  } catch (error) {
    console.error('Failed to load permissions:', error)
    // 使用模拟数据
    currentPath.value = searchPath.value
    currentPerm.value = {
      type: 'directory',
      permissions: 'rwxr-xr-x',
      owner: 'hserver',
      group: 'hserver',
      acl: []
    }
  }
}

// 应用基础权限
const applyBasicPerms = async () => {
  if (!currentPath.value) {
    alert('请先选择路径')
    return
  }

  applying.value = true

  try {
    const token = localStorage.getItem('token')
    // 修复API端点：使用正确的后端路由
    const response = await fetch('/api/permissions/files/permissions', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        owner: permForm.value.owner,
        group: permForm.value.group,
        permissions: currentPermNumeric.value,
        setuid: permForm.value.setuid,
        setgid: permForm.value.setgid,
        sticky: permForm.value.sticky
      })
    })

    if (response.ok) {
      await loadPermissions()
      alert('权限应用成功')
    } else {
      alert('权限应用失败')
    }
  } catch (error) {
    console.error('Failed to apply permissions:', error)
    alert('权限应用失败')
  } finally {
    applying.value = false
  }
}

// 添加 ACL 规则
const addACLRule = () => {
  if (!aclForm.value.entity) {
    alert('请选择实体')
    return
  }

  let permString = ''
  permString += aclForm.value.read ? 'r' : '-'
  permString += aclForm.value.write ? 'w' : '-'
  permString += aclForm.value.execute ? 'x' : '-'

  if (!currentPerm.value.acl) {
    currentPerm.value.acl = []
  }

  currentPerm.value.acl.push({
    type: aclForm.value.type,
    entity: aclForm.value.entity,
    permissions: permString
  })

  // 重置表单
  aclForm.value = {
    type: 'user',
    entity: '',
    read: false,
    write: false,
    execute: false
  }
}

// 移除 ACL 规则
const removeACLRule = (index: number) => {
  if (currentPerm.value.acl) {
    currentPerm.value.acl.splice(index, 1)
  }
}

// 应用 ACL
const applyACL = async () => {
  if (!currentPath.value) {
    alert('请先选择路径')
    return
  }

  applying.value = true

  try {
    const token = localStorage.getItem('token')
    // 修复API端点：使用正确的后端路由
    const response = await fetch('/api/permissions/files/acl', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        acl: currentPerm.value.acl
      })
    })

    if (response.ok) {
      await loadPermissions()
      alert('ACL 应用成功')
    } else {
      alert('ACL 应用失败')
    }
  } catch (error) {
    console.error('Failed to apply ACL:', error)
    alert('ACL 应用失败')
  } finally {
    applying.value = false
  }
}

// 移除所有 ACL
const removeAllACL = async () => {
  if (!confirm('确定要移除所有 ACL 规则吗？')) return

  try {
    const token = localStorage.getItem('token')
    // 修复API端点：使用正确的后端路由和HTTP方法
    const response = await fetch('/api/permissions/files/acl', {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        path: currentPath.value,
        acl: [], // 空数组表示移除所有ACL
        remove: true
      })
    })

    if (response.ok) {
      currentPerm.value.acl = []
      alert('ACL 移除成功')
    } else {
      alert('ACL 移除失败')
    }
  } catch (error) {
    console.error('Failed to remove ACL:', error)
    alert('ACL 移除失败')
  }
}

// 应用权限预设
const applyPreset = (perm: string) => {
  const num = parseInt(perm)

  permForm.value.setuid = (num & 0o4000) !== 0
  permForm.value.setgid = (num & 0o2000) !== 0
  permForm.value.sticky = (num & 0o1000) !== 0

  num &= 0o0777

  permForm.value.userRead = (num & 0o400) !== 0
  permForm.value.userWrite = (num & 0o200) !== 0
  permForm.value.userExec = (num & 0o100) !== 0

  permForm.value.groupRead = (num & 0o040) !== 0
  permForm.value.groupWrite = (num & 0o020) !== 0
  permForm.value.groupExec = (num & 0o010) !== 0

  permForm.value.otherRead = (num & 0o004) !== 0
  permForm.value.otherWrite = (num & 0o002) !== 0
  permForm.value.otherExec = (num & 0o001) !== 0
}

// 重置权限表单
const resetPermForm = () => {
  if (currentPerm.value) {
    loadPermissions()
  }
}

// 浏览路径
const browsePath = () => {
  showBrowserModal.value = true
  loadBrowserContents()
}

// 加载浏览器内容
const loadBrowserContents = async () => {
  try {
    const token = localStorage.getItem('token')
    const response = await fetch(`/api/files/list?path=${encodeURIComponent(browserPath.value)}`, {
      headers: { 'Authorization': `Bearer ${token}` }
    })

    if (response.ok) {
      const data = await response.json()
      browserContents.value = data
    } else {
      browserContents.value = []
    }
  } catch (error) {
    console.error('Failed to load browser contents:', error)
    browserContents.value = []
  }
}

// 进入目录
const enterDirectory = (item: any) => {
  if (item.type === 'dir') {
    browserPath.value += '/' + item.name
    loadBrowserContents()
  } else {
    // 选择了文件，查看权限
    searchPath.value = browserPath.value + '/' + item.name
    closeBrowser()
    loadPermissions()
  }
}

// 转到父目录
const goToParent = () => {
  const parts = browserPath.value.split('/').filter(Boolean)
  parts.pop()
  browserPath.value = '/' + parts.join('/')
  loadBrowserContents()
}

// 转到指定路径
const navigateToPath = () => {
  loadBrowserContents()
}

// 选择路径
const selectPath = () => {
  searchPath.value = browserPath.value
  closeBrowser()
  loadPermissions()
}

// 关闭浏览器
const closeBrowser = () => {
  showBrowserModal.value = false
}

onMounted(() => {
  // 初始加载权限信息
  loadPermissions()
})
</script>

<style scoped>
.acl-editor {
  width: 100%;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.editor-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.editor-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.path-input,
.perm-select,
.acl-select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.path-input {
  width: 300px;
}

.btn {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all 0.2s;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #6b7280;
  color: white;
}

.btn-secondary:hover {
  background: #4b5563;
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover {
  background: #059669;
}

.btn-ghost {
  background: transparent;
  color: #6b7280;
}

.btn-ghost:hover {
  background: #f3f4f6;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.current-permissions {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 20px;
}

.perm-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.perm-header h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.perm-badges {
  display: flex;
  gap: 8px;
}

.perm-badge {
  padding: 4px 12px;
  background: #f3f4f6;
  color: #374151;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.perm-badge.file {
  background: #dbeafe;
  color: #1e40af;
}

.perm-badge.directory {
  background: #d1fae5;
  color: #065f46;
}

.perm-details {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.perm-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
}

.perm-label {
  font-weight: 500;
  color: #6b7280;
  min-width: 80px;
}

.perm-value {
  color: #1f2937;
}

.permission-breakdown {
  display: flex;
  gap: 16px;
}

.perm-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.perm-group-title {
  font-size: 12px;
  color: #6b7280;
  min-width: 40px;
}

.perm-chip {
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  font-family: monospace;
}

.perm-chip.user {
  background: #dbeafe;
  color: #1e40af;
}

.perm-chip.group {
  background: #d1fae5;
  color: #065f46;
}

.perm-chip.other {
  background: #fef3c7;
  color: #92400e;
}

.acl-rules {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.acl-rule {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f9fafb;
  border-radius: 8px;
  font-size: 13px;
  font-family: monospace;
}

.rule-type {
  color: #6b7280;
}

.rule-entity {
  font-weight: 500;
  color: #1f2937;
}

.rule-permissions {
  color: #059669;
  font-weight: 600;
}

.permission-editor {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;
}

.editor-section {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
}

.editor-section h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.basic-perms {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.perm-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.perm-row label {
  font-weight: 500;
  color: #374151;
  min-width: 80px;
  font-size: 14px;
}

.perm-checks {
  display: flex;
  gap: 16px;
}

.perm-check-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.perm-check-title {
  font-size: 12px;
  color: #6b7280;
  min-width: 32px;
}

.perm-check {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  cursor: pointer;
}

.special-perms {
  display: flex;
  gap: 16px;
}

.perm-actions {
  display: flex;
  gap: 12px;
  margin-top: 16px;
}

.acl-adder {
  margin-bottom: 16px;
}

.acl-row {
  display: flex;
  gap: 12px;
  align-items: center;
}

.acl-perms-input {
  display: flex;
  gap: 8px;
  padding: 8px 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.acl-actions {
  display: flex;
  gap: 12px;
}

.permission-presets {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
}

.permission-presets h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.presets-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
}

.preset-btn {
  padding: 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  text-align: left;
  transition: all 0.2s;
}

.preset-btn:hover {
  background: #e5e7eb;
  border-color: #d1d5db;
}

.preset-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  font-family: monospace;
}

.preset-desc {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}

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

.modal {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.browser-modal {
  max-width: 800px;
}

.modal-header {
  padding: 20px 24px;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h4 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.modal-body {
  padding: 24px;
  overflow-y: auto;
}

.browser-nav {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 20px;
}

.browser-path-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.browser-contents {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.browser-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.browser-item:hover {
  background: #f3f4f6;
}

.browser-item span {
  flex: 1;
  font-size: 14px;
  color: #1f2937;
}

.item-perms {
  font-family: monospace;
  font-size: 12px;
  color: #6b7280;
}

.empty-browser {
  color: #6b7280;
  text-align: center;
  padding: 40px;
  font-style: italic;
}
</style>