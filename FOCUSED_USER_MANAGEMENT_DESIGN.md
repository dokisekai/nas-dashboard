# 专注用户管理功能设计方案

## 🎯 核心理念

"专注用户管理" = **单一职责 + 高效操作 + 清晰界面**

## 📊 功能架构设计

### 1. 用户中心化设计

```
用户管理中心
├── 用户生命周期管理
│   ├── 用户创建向导
│   ├── 用户信息维护
│   ├── 用户停用/删除
│   └── 用户数据导出
├── 权限管理系统
│   ├── 角色定义
│   ├── 权限分配
│   ├── 访问控制
│   └── 审计日志
├── 多服务同步
│   ├── 服务状态监控
│   ├── 批量用户操作
│   ├── 同步冲突解决
│   └── 失败重试机制
└── 用户体验优化
    ├── 快速搜索
    ├── 批量操作
    ├── 模板管理
    └── 自动化流程
```

## 🎨 界面设计方案

### 主界面布局

```vue
<template>
  <div class="focused-user-manager">
    <!-- 左侧导航区 -->
    <div class="sidebar">
      <div class="user-stats">
        <div class="stat-card">
          <div class="stat-number">{{ stats.totalUsers }}</div>
          <div class="stat-label">总用户</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ stats.activeUsers }}</div>
          <div class="stat-label">活跃用户</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ stats.syncedUsers }}</div>
          <div class="stat-label">已同步</div>
        </div>
      </div>
      
      <div class="quick-actions">
        <button @click="showCreateUser = true" class="action-btn primary">
          <UserPlusIcon /> 新建用户
        </button>
        <button @click="showBatchImport = true" class="action-btn">
          <DocumentArrowDownIcon /> 批量导入
        </button>
        <button @click="showTemplates = true" class="action-btn">
          <DocumentDuplicateIcon /> 用户模板
        </button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-content">
      <!-- 搜索和过滤栏 -->
      <div class="search-bar">
        <div class="search-input">
          <MagnifyingGlassIcon />
          <input 
            v-model="searchQuery" 
            placeholder="搜索用户名、邮箱、姓名..."
            @input="handleSearch"
          />
        </div>
        <div class="filter-options">
          <select v-model="statusFilter" class="filter-select">
            <option value="all">所有状态</option>
            <option value="active">活跃</option>
            <option value="inactive">停用</option>
            <option value="pending">待审核</option>
          </select>
          <select v-model="roleFilter" class="filter-select">
            <option value="all">所有角色</option>
            <option value="admin">管理员</option>
            <option value="user">普通用户</option>
            <option value="guest">访客</option>
          </select>
          <select v-model="serviceFilter" class="filter-select">
            <option value="all">所有服务</option>
            <option value="system">仅系统用户</option>
            <option value="immich">仅Immich用户</option>
            <option value="nextcloud">仅Nextcloud用户</option>
          </select>
        </div>
      </div>

      <!-- 用户列表 -->
      <div class="user-list">
        <div class="list-header">
          <div class="header-left">
            <input 
              type="checkbox" 
              v-model="selectAll"
              @change="toggleSelectAll"
            />
            <span class="user-count">{{ selectedUsers.length }}/{{ filteredUsers.length }} 已选择</span>
          </div>
          <div class="header-right">
            <button 
              v-if="selectedUsers.length > 0"
              @click="batchSync"
              class="batch-btn sync"
            >
              <CloudArrowUpIcon /> 批量同步
            </button>
            <button 
              v-if="selectedUsers.length > 0"
              @click="batchEdit"
              class="batch-btn edit"
            >
              <PencilIcon /> 批量编辑
            </button>
            <button 
              v-if="selectedUsers.length > 0"
              @click="batchDelete"
              class="batch-btn delete"
            >
              <TrashIcon /> 批量删除
            </button>
          </div>
        </div>

        <div class="user-cards">
          <div 
            v-for="user in paginatedUsers"
            :key="user.id"
            :class="['user-card', { selected: selectedUsers.includes(user.id) }]"
            @click="selectUser(user)"
          >
            <div class="card-header">
              <div class="user-avatar">
                <UserIcon />
              </div>
              <div class="user-info">
                <h3>{{ user.name || user.username }}</h3>
                <p>{{ user.email }}</p>
              </div>
              <div class="user-status">
                <span :class="['status-badge', user.status]">
                  {{ getStatusText(user.status) }}
                </span>
              </div>
            </div>
            
            <div class="card-body">
              <div class="user-details">
                <div class="detail-item">
                  <UserIcon />
                  <span>{{ user.username }}</span>
                </div>
                <div class="detail-item">
                  <ShieldIcon />
                  <span>{{ getRoleText(user.role) }}</span>
                </div>
                <div class="detail-item">
                  <ServerIcon />
                  <span>{{ user.services?.length || 0 }} 个服务</span>
                </div>
              </div>
              
              <div class="sync-status">
                <CloudArrowUpIcon 
                  :class="{ 'synced': user.lastSynced, 'pending': !user.lastSynced }"
                />
                <span>{{ user.lastSynced ? '已同步' : '待同步' }}</span>
                <span v-if="user.lastSynced" class="sync-time">
                  {{ formatTime(user.lastSynced) }}
                </span>
              </div>
            </div>
            
            <div class="card-actions">
              <button @click.stop="viewUser(user)" class="card-btn view">
                <EyeIcon /> 查看
              </button>
              <button @click.stop="editUser(user)" class="card-btn edit">
                <PencilIcon /> 编辑
              </button>
              <button @click.stop="syncUser(user)" class="card-btn sync">
                <CloudArrowUpIcon /> 同步
              </button>
              <button @click.stop="showUserMenu(user)" class="card-btn menu">
                <EllipsisVerticalIcon />
              </button>
            </div>
          </div>
        </div>

        <!-- 分页 -->
        <div class="pagination">
          <button 
            @click="previousPage" 
            :disabled="currentPage === 1"
            class="page-btn"
          >
            <ChevronLeftIcon />
          </button>
          <span class="page-info">{{ currentPage }} / {{ totalPages }}</span>
          <button 
            @click="nextPage" 
            :disabled="currentPage === totalPages"
            class="page-btn"
          >
            <ChevronRightIcon />
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧详情面板 -->
    <div class="detail-panel" v-if="selectedUserDetail">
      <div class="panel-header">
        <h2>用户详情</h2>
        <button @click="closeDetailPanel" class="close-btn">
          <XMarkIcon />
        </button>
      </div>
      
      <div class="panel-content">
        <div class="user-profile">
          <div class="profile-avatar">
            <UserIcon />
          </div>
          <div class="profile-info">
            <h3>{{ selectedUserDetail.name || selectedUserDetail.username }}</h3>
            <p>{{ selectedUserDetail.email }}</p>
            <div class="profile-meta">
              <span :class="['role-badge', selectedUserDetail.role]">
                {{ getRoleText(selectedUserDetail.role) }}
              </span>
              <span :class="['status-badge', selectedUserDetail.status]">
                {{ getStatusText(selectedUserDetail.status) }}
              </span>
            </div>
          </div>
        </div>

        <div class="user-tabs">
          <div 
            v-for="tab in tabs"
            :key="tab.id"
            :class="['tab-item', { active: activeTab === tab.id }]"
            @click="activeTab = tab.id"
          >
            <component :is="tab.icon" />
            {{ tab.label }}
          </div>
        </div>

        <div class="tab-content">
          <!-- 基本信息 -->
          <div v-if="activeTab === 'basic'" class="tab-section">
            <div class="info-grid">
              <div class="info-item">
                <label>用户名</label>
                <span>{{ selectedUserDetail.username }}</span>
              </div>
              <div class="info-item">
                <label>邮箱</label>
                <span>{{ selectedUserDetail.email }}</span>
              </div>
              <div class="info-item">
                <label>姓名</label>
                <span>{{ selectedUserDetail.name }}</span>
              </div>
              <div class="info-item">
                <label>角色</label>
                <span>{{ getRoleText(selectedUserDetail.role) }}</span>
              </div>
              <div class="info-item">
                <label>创建时间</label>
                <span>{{ formatDate(selectedUserDetail.createdAt) }}</span>
              </div>
              <div class="info-item">
                <label>最后登录</label>
                <span>{{ formatDate(selectedUserDetail.lastLogin) }}</span>
              </div>
            </div>
          </div>

          <!-- 服务状态 -->
          <div v-if="activeTab === 'services'" class="tab-section">
            <div class="services-list">
              <div 
                v-for="service in selectedUserDetail.services"
                :key="service.name"
                :class="['service-item', service.status]"
              >
                <div class="service-icon">
                  <ServerIcon />
                </div>
                <div class="service-info">
                  <h4>{{ service.name }}</h4>
                  <p>{{ service.statusText }}</p>
                </div>
                <div class="service-status">
                  <CheckCircleIcon v-if="service.status === 'active'" />
                  <XCircleIcon v-else-if="service.status === 'error'" />
                  <ClockIcon v-else />
                </div>
              </div>
            </div>
          </div>

          <!-- 活动日志 -->
          <div v-if="activeTab === 'activity'" class="tab-section">
            <div class="activity-timeline">
              <div 
                v-for="activity in selectedUserDetail.activities"
                :key="activity.id"
                class="timeline-item"
              >
                <div class="timeline-icon">
                  <component :is="getActivityIcon(activity.type)" />
                </div>
                <div class="timeline-content">
                  <h4>{{ activity.title }}</h4>
                  <p>{{ activity.description }}</p>
                  <span class="timeline-time">{{ formatTime(activity.timestamp) }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- 权限管理 -->
          <div v-if="activeTab === 'permissions'" class="tab-section">
            <div class="permissions-list">
              <div class="permission-group">
                <h4>系统权限</h4>
                <div class="permission-items">
                  <label v-for="perm in systemPermissions" :key="perm.id" class="permission-item">
                    <input 
                      type="checkbox" 
                      :checked="hasPermission(perm.id)"
                      @change="togglePermission(perm.id)"
                    />
                    <span>{{ perm.name }}</span>
                  </label>
                </div>
              </div>
              
              <div class="permission-group">
                <h4>服务权限</h4>
                <div class="permission-items">
                  <label v-for="perm in servicePermissions" :key="perm.id" class="permission-item">
                    <input 
                      type="checkbox" 
                      :checked="hasPermission(perm.id)"
                      @change="togglePermission(perm.id)"
                    />
                    <span>{{ perm.name }}</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="panel-footer">
        <button @click="editUser(selectedUserDetail)" class="action-btn primary">
          <PencilIcon /> 编辑用户
        </button>
        <button @click="syncUser(selectedUserDetail)" class="action-btn">
          <CloudArrowUpIcon /> 同步用户
        </button>
        <button @click="showDeleteConfirm(selectedUserDetail)" class="action-btn danger">
          <TrashIcon /> 删除用户
        </button>
      </div>
    </div>
  </div>
</template>
```

## 🔧 核心功能实现

### 1. 用户创建向导

```typescript
// composables/useUserWizard.ts
export function useUserWizard() {
  const steps = ref([
    'basic',      // 基本信息
    'services',   // 服务选择
    'permissions', // 权限设置
    'review'      // 确认创建
  ])
  
  const currentStep = ref(0)
  const userData = ref({
    basic: {
      username: '',
      email: '',
      name: '',
      password: '',
      confirmPassword: ''
    },
    services: {
      system: true,
      immich: false,
      nextcloud: false,
      jellyfin: false
    },
    permissions: {
      role: 'user',
      groups: [],
      customPermissions: []
    }
  })

  const validateStep = (step: number) => {
    switch(step) {
      case 0: // 基本信息验证
        return validateBasicInfo(userData.value.basic)
      case 1: // 服务选择验证
        return Object.values(userData.value.services).some(v => v)
      case 2: // 权限验证
        return userData.value.permissions.role !== ''
      default:
        return true
    }
  }

  const createUser = async () => {
    // 组装最终数据
    const finalUser = {
      username: userData.value.basic.username,
      email: userData.value.basic.email,
      name: userData.value.basic.name,
      password: userData.value.basic.password,
      services: Object.keys(userData.value.services)
        .filter(key => userData.value.services[key]),
      permissions: userData.value.permissions
    }

    // 调用统一用户管理API
    return await api.post('/api/unified-users', finalUser)
  }

  return {
    steps,
    currentStep,
    userData,
    validateStep,
    createUser,
    nextStep: () => { if (validateStep(currentStep.value)) currentStep.value++ },
    prevStep: () => currentStep.value--
  }
}
```

### 2. 智能搜索功能

```typescript
// composables/useUserSearch.ts
export function useUserSearch() {
  const searchIndex = ref<Map<string, any>>(new Map())
  const searchHistory = ref<string[]>([])
  
  // 构建搜索索引
  const buildSearchIndex = (users: any[]) => {
    const index = new Map()
    
    users.forEach(user => {
      const searchText = [
        user.username,
        user.email,
        user.name,
        user.role,
        ...user.groups
      ].join(' ').toLowerCase()
      
      index.set(user.id, {
        user,
        searchText,
        metadata: {
          createdAt: user.createdAt,
          lastLogin: user.lastLogin,
          services: user.services?.length || 0
        }
      })
    })
    
    searchIndex.value = index
  }

  // 智能搜索
  const searchUsers = (query: string) => {
    if (!query) return []

    // 添加到搜索历史
    if (!searchHistory.value.includes(query)) {
      searchHistory.value.unshift(query)
      if (searchHistory.value.length > 10) {
        searchHistory.value.pop()
      }
    }

    const terms = query.toLowerCase().split(/\s+/)
    const results: Array<{user: any, score: number}> = []

    searchIndex.value.forEach((data, userId) => {
      let score = 0
      const searchText = data.searchText

      // 精确匹配
      if (searchText.includes(query)) {
        score = 100
      }
      // 部分匹配
      else if (terms.some(term => searchText.includes(term))) {
        score = 50
      }
      // 模糊匹配
      else if (terms.some(term => {
        const distance = levenshteinDistance(term, searchText.substring(0, term.length))
        return distance <= 2
      })) {
        score = 25
      }

      if (score > 0) {
        results.push({ user: data.user, score })
      }
    })

    return results
      .sort((a, b) => b.score - a.score)
      .map(r => r.user)
  }

  return {
    buildSearchIndex,
    searchUsers,
    searchHistory,
    clearHistory: () => searchHistory.value = []
  }
}
```

### 3. 批量操作管理

```typescript
// composables/useBatchOperations.ts
export function useBatchOperations() {
  const selectedUsers = ref<Set<string>>(new Set())
  const operationProgress = ref({
    total: 0,
    completed: 0,
    failed: 0,
    current: ''
  })

  const batchSync = async (userIds: string[]) => {
    operationProgress.value = {
      total: userIds.length,
      completed: 0,
      failed: 0,
      current: ''
    }

    const results = []

    for (const userId of userIds) {
      operationProgress.value.current = userId

      try {
        const result = await api.post(`/api/unified-users/sync/${userId}`)
        operationProgress.value.completed++
        results.push({ userId, status: 'success', result })
      } catch (error) {
        operationProgress.value.failed++
        results.push({ userId, status: 'failed', error })
      }

      // 更新进度
      updateProgress(operationProgress.value)
    }

    return results
  }

  const batchEdit = async (userIds: string[], updates: any) => {
    // 批量编辑逻辑
  }

  const batchDelete = async (userIds: string[]) => {
    // 批量删除逻辑
  }

  const exportUsers = async (userIds: string[]) => {
    // 导出用户数据为CSV/Excel
    const users = await api.post('/api/unified-users/export', { userIds })
    
    // 生成文件并下载
    const blob = new Blob([users.data], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `users_export_${Date.now()}.csv`
    a.click()
  }

  const importUsers = async (file: File) => {
    // 导入用户数据
    const formData = new FormData()
    formData.append('file', file)

    const result = await api.post('/api/unified-users/import', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    return result
  }

  return {
    selectedUsers,
    operationProgress,
    batchSync,
    batchEdit,
    batchDelete,
    exportUsers,
    importUsers
  }
}
```

### 4. 用户模板系统

```typescript
// composables/useUserTemplates.ts
export function useUserTemplates() {
  const templates = ref<UserTemplate[]>([
    {
      id: 'admin',
      name: '管理员模板',
      description: '具有完整权限的管理员用户',
      permissions: {
        role: 'admin',
        groups: ['sudo', 'docker', 'network'],
        services: ['system', 'immich', 'nextcloud', 'jellyfin']
      }
    },
    {
      id: 'media-user',
      name: '媒体用户模板',
      description: '媒体服务专用用户',
      permissions: {
        role: 'user',
        groups: ['media', 'storage'],
        services: ['jellyfin', 'plex']
      }
    },
    {
      id: 'storage-user',
      name: '存储用户模板',
      description: '存储和备份专用用户',
      permissions: {
        role: 'user',
        groups: ['storage', 'backup'],
        services: ['nextcloud', 'system']
      }
    }
  ])

  const applyTemplate = async (userId: string, templateId: string) => {
    const template = templates.value.find(t => t.id === templateId)
    if (!template) return

    // 应用模板到用户
    const result = await api.put(`/api/unified-users/${userId}`, {
      permissions: template.permissions,
      applyTemplate: true
    })

    return result
  }

  const createCustomTemplate = (template: UserTemplate) => {
    templates.value.push(template)
    // 保存到本地存储或服务器
    localStorage.setItem('user-templates', JSON.stringify(templates.value))
  }

  return {
    templates,
    applyTemplate,
    createCustomTemplate
  }
}
```

## 🎨 专注界面设计原则

### 1. 视觉层次

```
主要操作区 (最重要)
├── 新建用户按钮
├── 批量操作按钮
└── 搜索过滤

次要信息区 (重要但不过分)
├── 用户列表展示
├── 状态指示器
└── 快捷操作

详细信息区 (按需显示)
├── 用户详情面板
├── 同步状态详情
└── 活动日志
```

### 2. 色彩系统

```css
/* 专注模式配色方案 */
:root {
  /* 主色调 - 专注蓝 */
  --primary-color: #2563eb;
  --primary-light: #3b82f6;
  --primary-dark: #1e40af;

  /* 功能色彩 */
  --success-color: #10b981;  /* 成功操作 */
  --warning-color: #f59e0b;  /* 需要注意 */
  --danger-color: #ef4444;   /* 危险操作 */
  --info-color: #3b82f6;     /* 信息提示 */

  /* 状态色彩 */
  --active-color: #10b981;
  --inactive-color: #6b7280;
  --pending-color: #f59e0b;
  --error-color: #ef4444;

  /* 中性色彩 */
  --bg-primary: #ffffff;
  --bg-secondary: #f9fafb;
  --bg-tertiary: #f3f4f6;
  
  --text-primary: #111827;
  --text-secondary: #6b7280;
  --text-tertiary: #9ca3af;

  /* 边框色彩 */
  --border-primary: #e5e7eb;
  --border-secondary: #d1d5db;
  --border-focus: #2563eb;
}
```

### 3. 交互设计

```typescript
// composables/useKeyboardShortcuts.ts
export function useKeyboardShortcuts() {
  const shortcuts = {
    'Ctrl+N': () => showCreateUser.value = true,
    'Ctrl+F': () => searchInput.value?.focus(),
    'Ctrl+A': () => selectAll(),
    'Ctrl+D': () => batchDelete(),
    'Escape': () => clearSelection(),
    'Enter': () => editSelectedUser(),
    'Delete': () => showDeleteConfirm()
  }

  onMounted(() => {
    window.addEventListener('keydown', (e) => {
      const shortcut = Object.entries(shortcuts).find(([key]) => {
        if (key.includes('Ctrl')) {
          return e.ctrlKey && e.key === key.split('+')[1]
        }
        return e.key === key
      })

      if (shortcut) {
        shortcut[1]()
        e.preventDefault()
      }
    })
  })
}
```

## 🔐 安全性设计

### 1. 权限控制矩阵

```typescript
// 权限定义
const PERMISSIONS = {
  // 用户管理
  'user.create': '创建新用户',
  'user.read': '查看用户信息',
  'user.update': '更新用户信息',
  'user.delete': '删除用户',
  'user.sync': '同步用户到服务',
  
  // 权限管理
  'permission.grant': '授予用户权限',
  'permission.revoke': '撤销用户权限',
  'permission.manage': '管理角色和组',
  
  // 系统管理
  'system.config': '修改系统配置',
  'system.logs': '查看系统日志',
  'system.backup': '执行系统备份',
  
  // 服务管理
  'service.manage': '管理Docker服务',
  'service.monitor': '监控服务状态'
}

// 权限检查
function checkPermission(requiredPermission: string): boolean {
  const user = getCurrentUser()
  if (user.role === 'admin') return true
  
  return user.permissions?.includes(requiredPermission) || false
}
```

### 2. 审计日志

```typescript
// 审计日志记录
interface AuditLog {
  id: string
  userId: string
  action: string
  resource: string
  details: any
  ipAddress: string
  userAgent: string
  timestamp: Date
}

function logAuditEvent(action: string, resource: string, details: any) {
  const log: AuditLog = {
    id: generateId(),
    userId: getCurrentUser().id,
    action,
    resource,
    details,
    ipAddress: getClientIP(),
    userAgent: navigator.userAgent,
    timestamp: new Date()
  }

  // 发送到审计日志服务
  api.post('/api/audit/logs', log)
}

// 使用示例
logAuditEvent('user.update', 'user_123', {
  changes: { email: 'old@example.com → new@example.com' },
  reason: '用户邮箱更新'
})
```

### 3. 数据加密

```typescript
// 敏感数据加密
import CryptoJS from 'crypto-js'

const SECRET_KEY = 'your-secret-key'

export function encryptPassword(password: string): string {
  return CryptoJS.AES.encrypt(password, SECRET_KEY).toString()
}

export function decryptPassword(encryptedPassword: string): string {
  const bytes = CryptoJS.AES.decrypt(encryptedPassword, SECRET_KEY)
  return bytes.toString(CryptoJS.enc.Utf8)
}

// 在传输前加密敏感数据
async function createUserSecurely(userData: any) {
  const encryptedData = {
    ...userData,
    password: encryptPassword(userData.password)
  }

  return await api.post('/api/unified-users', encryptedData)
}
```

## 📱 响应式设计

### 移动端适配

```css
/* 移动端优先设计 */
@media (max-width: 768px) {
  .focused-user-manager {
    flex-direction: column;
  }

  .sidebar {
    order: 2;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 60px;
    padding: 10px;
  }

  .main-content {
    order: 1;
    padding-bottom: 70px;
  }

  .detail-panel {
    position: fixed;
    top: 0;
    right: 0;
    bottom: 0;
    width: 100%;
    z-index: 1000;
  }

  .user-cards {
    grid-template-columns: 1fr;
  }
}
```

## 🚀 性能优化

### 1. 虚拟滚动

```vue
<template>
  <div class="user-list-container">
    <VirtualList
      :items="filteredUsers"
      :item-height="80"
      :buffer-size="20"
      v-slot="{ item: user }"
    >
      <UserCard :user="user" @select="selectUser" />
    </VirtualList>
  </div>
</template>
```

### 2. 数据缓存

```typescript
// composables/useUserCache.ts
export function useUserCache() {
  const cache = ref<Map<string, any>>(new Map())
  const cacheExpiry = ref<Map<string, number>>(new Map())

  const CACHE_DURATION = 5 * 60 * 1000 // 5分钟

  const getCachedUser = (userId: string) => {
    const expiry = cacheExpiry.value.get(userId)
    if (expiry && Date.now() < expiry) {
      return cache.value.get(userId)
    }
    return null
  }

  const setCachedUser = (userId: string, userData: any) => {
    cache.value.set(userId, userData)
    cacheExpiry.value.set(userId, Date.now() + CACHE_DURATION)
  }

  const clearCache = () => {
    cache.value.clear()
    cacheExpiry.value.clear()
  }

  return {
    getCachedUser,
    setCachedUser,
    clearCache
  }
}
```

### 3. 懒加载

```typescript
// 组件懒加载
const UnifiedUserManager = defineAsyncComponent(() => 
  import('@/apps/UnifiedUserManager.vue')
)

const UserDetailPanel = defineAsyncComponent(() => 
  import('@/components/UserDetailPanel.vue')
)

// 路由懒加载
const routes = [
  {
    path: '/users',
    component: () => import('@/views/Users.vue')
  }
]
```

## 📊 数据管理

### 1. 状态管理

```typescript
// stores/userManagement.ts
import { defineStore } from 'pinia'

export const useUserManagementStore = defineStore('userManagement', () => {
  // 状态
  const users = ref<Map<string, User>>(new Map())
  const selectedUsers = ref<Set<string>>(new Set())
  const filters = ref({
    search: '',
    status: 'all',
    role: 'all',
    service: 'all'
  })
  const pagination = ref({
    currentPage: 1,
    pageSize: 20,
    total: 0
  })

  // 操作
  async function loadUsers() {
    const response = await api.get('/api/unified-users')
    const usersData = response.data

    users.value.clear()
    usersData.forEach((user: User) => {
      users.value.set(user.id, user)
    })

    pagination.value.total = usersData.length
  }

  function selectUser(userId: string) {
    selectedUsers.value.add(userId)
  }

  function deselectUser(userId: string) {
    selectedUsers.value.delete(userId)
  }

  function toggleUserSelection(userId: string) {
    if (selectedUsers.value.has(userId)) {
      deselectUser(userId)
    } else {
      selectUser(userId)
    }
  }

  // 计算属性
  const filteredUsers = computed(() => {
    let result = Array.from(users.value.values())

    // 搜索过滤
    if (filters.value.search) {
      const searchLower = filters.value.search.toLowerCase()
      result = result.filter(user =>
        user.username.toLowerCase().includes(searchLower) ||
        user.email.toLowerCase().includes(searchLower) ||
        user.name?.toLowerCase().includes(searchLower)
      )
    }

    // 状态过滤
    if (filters.value.status !== 'all') {
      result = result.filter(user => user.status === filters.value.status)
    }

    // 角色过滤
    if (filters.value.role !== 'all') {
      result = result.filter(user => user.role === filters.value.role)
    }

    return result
  })

  const paginatedUsers = computed(() => {
    const start = (pagination.value.currentPage - 1) * pagination.value.pageSize
    const end = start + pagination.value.pageSize
    return filteredUsers.value.slice(start, end)
  })

  return {
    users,
    selectedUsers,
    filters,
    pagination,
    filteredUsers,
    paginatedUsers,
    loadUsers,
    selectUser,
    deselectUser,
    toggleUserSelection
  }
})
```

## 🎯 实施路线图

### 第一阶段：核心功能 (1-2周)
- [x] 基础用户CRUD
- [x] 用户列表界面
- [x] 基本搜索功能
- [x] 用户详情查看

### 第二阶段：增强功能 (2-3周)
- [ ] 批量操作
- [ ] 高级搜索
- [ ] 用户模板
- [ ] 导入导出

### 第三阶段：智能化 (3-4周)
- [ ] 自动同步
- [ ] 冲突解决
- [ ] 智能推荐
- [ ] 活动分析

### 第四阶段：优化完善 (持续)
- [ ] 性能优化
- [ ] 用户体验优化
- [ ] 安全加固
- [ ] 监控告警

## 💡 最佳实践建议

1. **界面简洁** - 每个屏幕专注一个主要任务
2. **操作直观** - 常用操作一步到位
3. **反馈及时** - 每个操作都有明确反馈
4. **错误友好** - 清晰的错误信息和解决建议
5. **性能优先** - 大数据量下保持流畅
6. **安全第一** - 所有操作都有权限控制
7. **可扩展性** - 易于添加新功能和服务

这个专注的用户管理功能将极大提升你的NAS Dashboard的用户管理效率！