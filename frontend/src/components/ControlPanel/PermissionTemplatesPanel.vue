<template>
  <div class="permission-templates-panel">
    <div class="panel-toolbar">
      <div class="toolbar-left">
        <h3>权限模板</h3>
        <p class="subtitle">预设权限配置，快速分配用户权限</p>
      </div>
      <div class="toolbar-right">
        <button class="btn btn-primary" @click="showCreateTemplateModal = true">
          <PlusIcon class="w-4 h-4" />
          创建模板
        </button>
      </div>
    </div>

    <div class="templates-grid">
      <div
        v-for="template in templates"
        :key="template.id"
        class="template-card"
        :class="{ 'builtin-template': template.builtin }"
      >
        <div class="template-header">
          <div class="template-icon" :style="{ background: template.color }">
            <component :is="getTemplateIcon(template.type)" class="w-6 h-6" />
          </div>
          <div class="template-info">
            <h4>{{ template.name }}</h4>
            <span class="template-type">{{ getTemplateTypeLabel(template.type) }}</span>
          </div>
        </div>

        <p class="template-description">{{ template.description }}</p>

        <div class="template-features">
          <div
            v-for="feature in template.features"
            :key="feature"
            class="feature-item"
          >
            <CheckIcon class="w-4 h-4" />
            <span>{{ feature }}</span>
          </div>
        </div>

        <div class="template-actions">
          <button class="btn btn-secondary" @click="viewTemplate(template)">
            <EyeIcon class="w-4 h-4" />
            查看详情
          </button>
          <button v-if="!template.builtin" class="btn btn-secondary" @click="editTemplate(template)">
            <PencilIcon class="w-4 h-4" />
            编辑
          </button>
        </div>
      </div>
    </div>

    <!-- 查看模板详情对话框 -->
    <div v-if="showTemplateDetailsModal" class="modal-overlay" @click.self="closeTemplateModal">
      <div class="modal-content extra-large">
        <div class="modal-header">
          <h3>模板详情 - {{ selectedTemplate?.name }}</h3>
          <button class="close-btn" @click="closeTemplateModal">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="template-details">
          <div class="template-overview">
            <div class="overview-card">
              <h4>模板信息</h4>
              <div class="overview-info">
                <div class="info-row">
                  <span class="info-label">类型</span>
                  <span class="info-value">{{ getTemplateTypeLabel(selectedTemplate?.type) }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">描述</span>
                  <span class="info-value">{{ selectedTemplate?.description }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">适用用户</span>
                  <span class="info-value">{{ selectedTemplate?.userCount || 0 }} 个用户</span>
                </div>
              </div>
            </div>
          </div>

          <div class="permissions-breakdown">
            <h4>权限详情</h4>
            <div class="permissions-grid">
              <div
                v-for="(category, catKey) in selectedTemplate?.permissions"
                :key="catKey"
                class="permission-category-card"
              >
                <h5>{{ getCategoryLabel(catKey) }}</h5>
                <div class="category-permissions">
                  <div
                    v-for="(level, permKey) in category"
                    :key="permKey"
                    class="permission-detail"
                  >
                    <span class="perm-name">{{ getPermissionLabel(permKey) }}</span>
                    <span :class="['perm-level', `level-${level}`]">
                      {{ getPermissionLevelLabel(level) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-actions">
          <button class="btn btn-primary" @click="applyTemplate">
            应用此模板
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  PlusIcon,
  EyeIcon,
  PencilIcon,
  CheckIcon,
  XMarkIcon,
  ShieldCheckIcon,
  UserIcon,
  KeyIcon,
  CubeIcon
} from '@heroicons/vue/24/outline'

const props = defineProps<{
  templates: any[]
}>()

const emit = defineEmits<{
  'template-applied': [data: any]
}>()

const showCreateTemplateModal = ref(false)
const showTemplateDetailsModal = ref(false)
const selectedTemplate = ref<any>(null)

const templates = ref([
  {
    id: 'guest',
    name: '访客模板',
    type: 'guest',
    description: '仅限基本访问权限，适合临时访客',
    color: 'linear-gradient(135deg, #9ca3af 0%, #6b7280 100%)',
    builtin: true,
    userCount: 5,
    features: ['只读文件访问', '基本系统查看', '无管理权限'],
    permissions: {
      system: { settings: 'none', reboot: 'none', shutdown: 'none' },
      storage: { view: 'read', manage: 'none', format: 'none' },
      file: { read: 'read', write: 'none', delete: 'none', share: 'none' }
    }
  },
  {
    id: 'user',
    name: '普通用户模板',
    type: 'user',
    description: '标准用户权限，适合日常使用',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
    builtin: true,
    userCount: 15,
    features: ['文件读写', '个人配置', '应用使用', '基本监控'],
    permissions: {
      system: { settings: 'read', reboot: 'none', shutdown: 'none' },
      storage: { view: 'read', manage: 'none', format: 'none' },
      file: { read: 'write', write: 'write', delete: 'write', share: 'read' }
    }
  },
  {
    id: 'admin',
    name: '管理员模板',
    type: 'admin',
    description: '系统管理权限，可管理大部分功能',
    color: 'linear-gradient(135deg, #f59e0b 0%, #ea580c 100%)',
    builtin: true,
    userCount: 3,
    features: ['完整系统管理', '用户权限管理', '存储配置', '网络设置'],
    permissions: {
      system: { settings: 'admin', reboot: 'write', shutdown: 'write' },
      storage: { view: 'admin', manage: 'admin', format: 'admin' },
      file: { read: 'admin', write: 'admin', delete: 'admin', share: 'admin' }
    }
  },
  {
    id: 'superadmin',
    name: '超级管理员模板',
    type: 'superadmin',
    description: '最高权限，拥有所有系统控制权',
    color: 'linear-gradient(135deg, #ef4444 0%, #dc2626 100%)',
    builtin: true,
    userCount: 1,
    features: ['所有系统权限', '用户组管理', '系统重启关机', '权限模板管理'],
    permissions: {
      system: { settings: 'owner', reboot: 'owner', shutdown: 'owner' },
      storage: { view: 'owner', manage: 'owner', format: 'owner' },
      file: { read: 'owner', write: 'owner', delete: 'owner', share: 'owner' }
    }
  }
])

const getTemplateIcon = (type: string) => {
  const icons: Record<string, any> = {
    guest: UserIcon,
    user: UserIcon,
    admin: ShieldCheckIcon,
    superadmin: KeyIcon
  }
  return icons[type] || CubeIcon
}

const getTemplateTypeLabel = (type?: string) => {
  const labels = {
    guest: '访客',
    user: '普通用户',
    admin: '管理员',
    superadmin: '超级管理员'
  }
  return labels[type as keyof typeof labels] || type
}

const getCategoryLabel = (key: string) => {
  const labels = {
    system: '系统管理',
    storage: '存储管理',
    file: '文件操作',
    network: '网络管理',
    user: '用户管理',
    backup: '备份恢复',
    app: '应用管理',
    log: '日志查看',
    monitor: '系统监控'
  }
  return labels[key as keyof typeof labels] || key
}

const getPermissionLabel = (key: string) => {
  const labels = {
    settings: '系统设置',
    reboot: '系统重启',
    shutdown: '系统关机',
    view: '查看',
    manage: '管理',
    format: '格式化',
    read: '读取',
    write: '写入',
    delete: '删除',
    share: '共享'
  }
  return labels[key as keyof typeof labels] || key
}

const getPermissionLevelLabel = (level: string) => {
  const labels = {
    none: '无权限',
    read: '只读',
    write: '读写',
    admin: '管理',
    owner: '所有者'
  }
  return labels[level as keyof typeof labels] || level
}

const viewTemplate = (template: any) => {
  selectedTemplate.value = template
  showTemplateDetailsModal.value = true
}

const editTemplate = (template: any) => {
  console.log('Edit template:', template)
}

const applyTemplate = () => {
  emit('template-applied', selectedTemplate.value)
  closeTemplateModal()
}

const closeTemplateModal = () => {
  showTemplateDetailsModal.value = false
  selectedTemplate.value = null
}
</script>

<style scoped>
.permission-templates-panel {
  width: 100%;
}

.panel-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.toolbar-left h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
}

.toolbar-right {
  display: flex;
  gap: 12px;
}

.btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: white;
  color: #374151;
  border: 1px solid #e5e7eb;
}

.btn-secondary:hover {
  background: #f9fafb;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
}

.template-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border: 2px solid transparent;
  transition: all 0.2s;
}

.template-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.builtin-template {
  border-color: #3b82f6;
}

.template-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.template-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.template-info h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.template-type {
  font-size: 12px;
  color: #6b7280;
  background: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
}

.template-description {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 16px;
  min-height: 40px;
}

.template-features {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #374151;
}

.feature-item svg {
  color: #10b981;
  flex-shrink: 0;
}

.template-actions {
  display: flex;
  gap: 8px;
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

.modal-content {
  background: white;
  border-radius: 16px;
  max-width: 900px;
  width: 90%;
  max-height: 80vh;
  overflow-y: auto;
}

.extra-large {
  max-width: 1000px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  border: none;
  background: #f3f4f6;
  color: #6b7280;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.template-details {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.template-overview {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.overview-card {
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
}

.overview-card h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.overview-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.info-label {
  color: #6b7280;
}

.info-value {
  font-weight: 500;
  color: #1f2937;
}

.permissions-breakdown h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.permissions-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.permission-category-card {
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
}

.permission-category-card h5 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 12px;
}

.category-permissions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.permission-detail {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.perm-name {
  color: #6b7280;
}

.perm-level {
  font-weight: 500;
  padding: 2px 6px;
  border-radius: 4px;
}

.level-none {
  color: #6b7280;
}

.level-read {
  color: #3b82f6;
  background: #dbeafe;
}

.level-write {
  color: #10b981;
  background: #d1fae5;
}

.level-admin {
  color: #f59e0b;
  background: #fed7aa;
}

.level-owner {
  color: #ef4444;
  background: #fee2e2;
}

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
}
</style>