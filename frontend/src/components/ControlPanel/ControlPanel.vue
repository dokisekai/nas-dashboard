<template>
  <div class="control-panel">
    <!-- 顶部工具栏 -->
    <div class="cp-header">
      <div class="cp-header-left">
        <h1 class="cp-title">控制面板</h1>
        <p class="cp-description">系统配置和管理中心</p>
      </div>

      <div class="cp-header-actions">
        <div class="cp-search">
          <MagnifyingGlassIcon class="w-5 h-5" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索设置..."
            @input="handleSearch"
          />
        </div>

        <button
          v-if="unsavedChanges"
          class="cp-btn cp-btn-primary"
          @click="applyChanges"
          :disabled="saving"
        >
          <CheckIcon class="w-4 h-4" />
          {{ saving ? '应用中...' : '应用更改' }}
        </button>

        <button
          v-if="unsavedChanges"
          class="cp-btn cp-btn-secondary"
          @click="discardChanges"
        >
          <XMarkIcon class="w-4 h-4" />
          丢弃
        </button>

        <div class="cp-header-menu">
          <button class="cp-btn cp-btn-ghost" @click="showExportModal = true" title="导出配置">
            <ArrowDownTrayIcon class="w-4 h-4" />
          </button>
          <button class="cp-btn cp-btn-ghost" @click="showImportModal = true" title="导入配置">
            <ArrowUpTrayIcon class="w-4 h-4" />
          </button>
          <button class="cp-btn cp-btn-ghost" @click="resetToDefaults" title="重置默认">
            <ArrowPathIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- 未保存更改警告 -->
    <div v-if="unsavedChanges" class="cp-unsaved-warning">
      <ExclamationTriangleIcon class="w-5 h-5" />
      <span>您有未保存的更改。记得应用更改使其生效。</span>
    </div>

    <div class="cp-content">
      <!-- 侧边栏分类 -->
      <div class="cp-sidebar">
        <div class="cp-categories">
          <button
            v-for="category in categories"
            :key="category.id"
            :class="['cp-category-btn', { active: activeCategory === category.id }]"
            @click="selectCategory(category.id)"
          >
            <component :is="getIcon(category.icon)" class="w-5 h-5" />
            <div class="cp-category-info">
              <span class="cp-category-name">{{ category.name }}</span>
            </div>
            <ChevronRightIcon class="w-4 h-4 cp-category-arrow" />
          </button>
        </div>
      </div>

      <!-- 主内容区域 -->
      <div class="cp-main">
        <!-- 加载状态 -->
        <div v-if="loading && !initialized" class="cp-loading">
          <div class="cp-spinner"></div>
          <p>加载控制面板...</p>
        </div>

        <!-- 搜索结果 -->
        <div v-else-if="searchResults.length > 0" class="cp-search-results">
          <div class="cp-section-header">
            <h2>搜索结果</h2>
            <span class="cp-count">找到 {{ searchResults.length }} 个设置</span>
          </div>

          <div class="cp-settings-grid">
            <div
              v-for="setting in searchResults"
              :key="setting.id"
              class="cp-setting-item"
            >
              <ControlPanelSettingComponent
                :setting="setting"
                :value="settings[setting.id]"
                @update="updateSetting(setting.id, $event)"
              />
            </div>
          </div>
        </div>

        <!-- 分类设置 -->
        <div v-else-if="activeCategory && !searchQuery" class="cp-category-content">
          <!-- 常规设置 -->
          <div class="cp-settings-section">
            <h3 class="cp-section-title">常规设置</h3>
            <div class="cp-settings-grid">
              <div
                v-for="setting in getRegularSettings(activeCategory)"
                :key="setting.id"
                class="cp-setting-item"
              >
                <ControlPanelSettingComponent
                  :setting="setting"
                  :value="settings[setting.id]"
                  @update="updateSetting(setting.id, $event)"
                />
              </div>
            </div>
          </div>

          <!-- 高级设置 -->
          <div v-if="getAdvancedSettingsByCategory(activeCategory).length > 0" class="cp-settings-section">
            <div class="cp-advanced-header">
              <h3 class="cp-section-title">高级设置</h3>
              <button
                class="cp-btn cp-btn-ghost cp-btn-sm"
                @click="showAdvanced = !showAdvanced"
              >
                {{ showAdvanced ? '隐藏' : '显示' }}
                <ChevronDownIcon :class="{ 'rotate-180': showAdvanced }" class="w-4 h-4" />
              </button>
            </div>

            <div v-show="showAdvanced" class="cp-settings-grid">
              <div
                v-for="setting in getAdvancedSettingsByCategory(activeCategory)"
                :key="setting.id"
                class="cp-setting-item"
              >
                <ControlPanelSettingComponent
                  :setting="setting"
                  :value="settings[setting.id]"
                  @update="updateSetting(setting.id, $event)"
                />
              </div>
            </div>
          </div>

          <!-- 分类操作 -->
          <div class="cp-category-actions">
            <button
              class="cp-btn cp-btn-secondary"
              @click="resetCategoryToDefaults(activeCategory)"
            >
              <ArrowPathIcon class="w-4 h-4" />
              重置分类默认值
            </button>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="cp-empty">
          <MagnifyingGlassIcon class="w-12 h-12" />
          <p>输入关键词搜索设置</p>
        </div>
      </div>
    </div>

    <!-- 导出配置模态框 -->
    <div v-if="showExportModal" class="cp-modal-overlay" @click="showExportModal = false">
      <div class="cp-modal" @click.stop>
        <div class="cp-modal-header">
          <h3>导出配置</h3>
          <button class="cp-btn cp-btn-ghost" @click="showExportModal = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="cp-modal-body">
          <p>导出当前系统配置到文件，可用于备份或在其他系统上导入。</p>

          <div class="cp-form-group">
            <label>包含元数据</label>
            <input v-model="exportIncludeMetadata" type="checkbox" />
            <span class="cp-form-help">包含主机名、导出时间等信息</span>
          </div>
        </div>

        <div class="cp-modal-footer">
          <button class="cp-btn cp-btn-secondary" @click="showExportModal = false">
            取消
          </button>
          <button class="cp-btn cp-btn-primary" @click="exportConfiguration">
            <ArrowDownTrayIcon class="w-4 h-4" />
            导出
          </button>
        </div>
      </div>
    </div>

    <!-- 导入配置模态框 -->
    <div v-if="showImportModal" class="cp-modal-overlay" @click="showImportModal = false">
      <div class="cp-modal" @click.stop>
        <div class="cp-modal-header">
          <h3>导入配置</h3>
          <button class="cp-btn cp-btn-ghost" @click="showImportModal = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="cp-modal-body">
          <p>从配置文件导入系统设置。这将覆盖当前配置。</p>

          <div class="cp-form-group">
            <label>配置文件</label>
            <textarea
              v-model="importConfigText"
              placeholder="粘贴配置JSON内容..."
              rows="10"
            ></textarea>
          </div>

          <div v-if="importError" class="cp-error-message">
            <ExclamationCircleIcon class="w-5 h-5" />
            {{ importError }}
          </div>
        </div>

        <div class="cp-modal-footer">
          <button class="cp-btn cp-btn-secondary" @click="showImportModal = false">
            取消
          </button>
          <button
            class="cp-btn cp-btn-primary"
            @click="importConfiguration"
            :disabled="!importConfigText"
          >
            <ArrowUpTrayIcon class="w-4 h-4" />
            导入
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useControlPanelStore } from '../../stores/controlPanel'
import ControlPanelSettingComponent from './ControlPanelSettingComponent.vue'
import {
  MagnifyingGlassIcon,
  ChevronRightIcon,
  CheckIcon,
  XMarkIcon,
  ExclamationTriangleIcon,
  ArrowDownTrayIcon,
  ArrowUpTrayIcon,
  ArrowPathIcon,
  ExclamationCircleIcon,
  ChevronDownIcon
} from '@heroicons/vue/24/outline'
import {
  CogIcon,
  GlobeAltIcon,
  ShieldCheckIcon,
  ServerIcon,
  BellIcon,
  PaintBrushIcon,
  BeakerIcon,
  CircleStackIcon
} from '@heroicons/vue/24/outline'

// 图标映射 - 现在使用直接导入的图标
const iconComponents: Record<string, any> = {
  CogIcon,
  GlobeAltIcon,
  ShieldCheckIcon,
  ServerIcon,
  BellIcon,
  PaintBrushIcon,
  BeakerIcon,
  CircleStackIcon
}

const controlPanel = useControlPanelStore()

// 状态
const activeCategory = ref<string>('general')
const searchQuery = ref('')
const searchResults = ref<any[]>([])
const showAdvanced = ref(false)
const showExportModal = ref(false)
const showImportModal = ref(false)
const exportIncludeMetadata = ref(true)
const importConfigText = ref('')
const importError = ref('')

// 从 store 中提取响应式状态
const {
  categories,
  settings,
  initialized,
  loading,
  saving,
  unsavedChanges
} = storeToRefs(controlPanel)

// 方法
const selectCategory = (categoryId: string) => {
  activeCategory.value = categoryId
  showAdvanced.value = false
}

const handleSearch = () => {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }

  searchResults.value = controlPanel.search(searchQuery.value)
}

const getCategory = (categoryId: string) => {
  return categories.value.find(c => c.id === categoryId)
}

const getRegularSettings = (categoryId: string) => {
  return controlPanel.getSettingsByCategory(categoryId).filter(s => !s.advanced)
}

const getAdvancedSettingsByCategory = (categoryId: string) => {
  return controlPanel.getSettingsByCategory(categoryId).filter(s => s.advanced)
}

const updateSetting = async (key: string, value: any) => {
  try {
    await controlPanel.updateSetting(key, value)
  } catch (error: any) {
    console.error('Failed to update setting:', error)
    // 显示错误提示
  }
}

const applyChanges = async () => {
  try {
    await controlPanel.applyChanges()
  } catch (error: any) {
    console.error('Failed to apply changes:', error)
  }
}

const discardChanges = () => {
  controlPanel.discardChanges()
}

const resetToDefaults = async () => {
  if (confirm('确定要重置所有设置为默认值吗？此操作不可撤销。')) {
    try {
      await controlPanel.resetToDefaults()
    } catch (error: any) {
      console.error('Failed to reset defaults:', error)
    }
  }
}

const resetCategoryToDefaults = async (categoryId: string) => {
  if (confirm('确定要重置此分类的所有设置为默认值吗？')) {
    try {
      await controlPanel.resetToCategoryDefaults(categoryId)
    } catch (error: any) {
      console.error('Failed to reset category defaults:', error)
    }
  }
}

const exportConfiguration = () => {
  const config = controlPanel.exportConfig()
  const blob = new Blob([config], { type: 'application/json' })
  const url = URL.createObjectURL(blob)

  const a = document.createElement('a')
  a.href = url
  a.download = `nas-config-${new Date().toISOString()}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)

  showExportModal.value = false
}

const importConfiguration = async () => {
  importError.value = ''

  try {
    await controlPanel.importConfig(importConfigText.value)
    showImportModal.value = false
    importConfigText.value = ''
  } catch (error: any) {
    importError.value = error.message || '导入配置失败'
  }
}

const getIcon = (iconName: string) => {
  const icon = iconComponents[iconName]
  return icon || null
}

// 生命周期
onMounted(async () => {
  if (!initialized.value) {
    await controlPanel.initialize()
  }
})
</script>

<style scoped>
.control-panel {
  width: 100%;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
}

.cp-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  background: white;
  border-bottom: 1px solid rgba(102, 126, 234, 0.1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.cp-header-left {
  flex: 1;
}

.cp-title {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 4px;
  letter-spacing: 0.5px;
}

.cp-description {
  font-size: 14px;
  color: #64748b;
  font-weight: 500;
}

.cp-header-actions {
  display: flex;
  gap: 16px;
  align-items: center;
}

.cp-search {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #f8fafc;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  min-width: 280px;
  color: #6b7280;
}

.cp-search input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  outline: none;
  color: #1f2937;
}

.cp-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.cp-btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.cp-btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

.cp-btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.cp-btn-secondary {
  background: white;
  border: 1px solid #e5e7eb;
  color: #6b7280;
}

.cp-btn-secondary:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.cp-btn-ghost {
  background: transparent;
  border: none;
  color: #6b7280;
  padding: 8px;
}

.cp-btn-ghost:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.cp-btn-sm {
  padding: 6px 12px;
  font-size: 12px;
}

.cp-header-menu {
  display: flex;
  gap: 8px;
}

.cp-unsaved-warning {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 32px;
  background: #fef3c7;
  border-bottom: 1px solid #fcd34d;
  color: #92400e;
  font-size: 14px;
  font-weight: 500;
}

.cp-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.cp-sidebar {
  width: 320px;
  background: white;
  border-right: 1px solid rgba(102, 126, 234, 0.1);
  overflow-y: auto;
}

.cp-categories {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cp-category-btn {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: transparent;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  text-align: left;
  width: 100%;
  color: #6b7280;
}

.cp-category-btn:hover {
  background: rgba(102, 126, 234, 0.05);
  color: #667eea;
}

.cp-category-btn.active {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  color: #667eea;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.2);
}

.cp-category-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.cp-category-name {
  font-size: 14px;
  font-weight: 600;
}

.cp-category-desc {
  font-size: 12px;
  opacity: 0.7;
}

.cp-category-arrow {
  transition: transform 0.2s ease;
}

.cp-category-btn.active .cp-category-arrow {
  transform: rotate(90deg);
}

.cp-main {
  flex: 1;
  overflow-y: auto;
  padding: 32px;
}

.cp-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #6b7280;
}

.cp-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.cp-search-results,
.cp-category-content {
  max-width: 1200px;
  margin: 0 auto;
}

.cp-section-header {
  margin-bottom: 32px;
}

.cp-section-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.cp-section-header p {
  font-size: 14px;
  color: #6b7280;
}

.cp-count {
  font-size: 12px;
  color: #9ca3af;
  background: #f3f4f6;
  padding: 4px 12px;
  border-radius: 12px;
  margin-left: 12px;
}

.cp-settings-section {
  margin-bottom: 32px;
}

.cp-section-title {
  font-size: 16px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 20px;
}

.cp-advanced-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.cp-settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 20px;
}

.cp-setting-item {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(102, 126, 234, 0.1);
}

/* 自定义组件占据全宽 */
.cp-setting-item:has([data-setting-type="custom"]) {
  grid-column: 1 / -1;
  padding: 0;
  overflow: hidden;
}

.cp-setting-item:has([data-setting-type="custom"]) .cps-setting {
  padding: 24px;
}

.cp-category-actions {
  display: flex;
  gap: 16px;
  padding-top: 24px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.cp-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #9ca3af;
}

.cp-empty svg {
  margin-bottom: 16px;
  opacity: 0.5;
}

.cp-modal-overlay {
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

.cp-modal {
  background: white;
  border-radius: 16px;
  padding: 24px;
  min-width: 500px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.cp-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.cp-modal-header h3 {
  font-size: 20px;
  font-weight: 600;
  color: #1f2937;
}

.cp-modal-body {
  margin-bottom: 24px;
  color: #6b7280;
  font-size: 14px;
}

.cp-form-group {
  margin-top: 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cp-form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 8px;
}

.cp-form-help {
  font-size: 12px;
  color: #9ca3af;
}

.cp-form-group textarea {
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  font-family: monospace;
  resize: vertical;
}

.cp-error-message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #991b1b;
  font-size: 14px;
  margin-top: 16px;
}

.cp-modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.rotate-180 {
  transform: rotate(180deg);
}
</style>