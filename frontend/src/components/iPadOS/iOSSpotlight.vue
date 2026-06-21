<template>
  <div class="ios-spotlight" @click.self="$emit('close')">
    <div class="spotlight-container" @click.stop>

      <!-- 搜索输入区 -->
      <div class="search-container">
        <div class="search-icon">
          <MagnifyingGlassIcon class="icon-svg" />
        </div>
        <input
          v-model="searchQuery"
          ref="searchInput"
          type="text"
          class="search-input"
          placeholder="Spotlight 搜索"
          autofocus
          @keydown.enter="selectFirstResult"
          @keydown.up="navigateUp"
          @keydown.down="navigateDown"
          @keydown.esc="$emit('close')"
        >
        <div v-if="searchQuery" class="clear-btn" @click="clearSearch">
          <XMarkIcon class="clear-icon" />
        </div>
      </div>

      <!-- 搜索结果 -->
      <div class="results-container" v-if="searchQuery">
        <!-- 应用匹配 -->
        <div v-if="matchedApps.length > 0" class="result-section">
          <div class="section-header">应用</div>
          <div
            v-for="(app, index) in matchedApps"
            :key="app.id"
            class="result-item"
            :class="{ selected: selectedIndex === index }"
            @click="openApp(app)"
          >
            <div class="result-icon">
              <component :is="app.icon" class="icon-svg" />
            </div>
            <div class="result-info">
              <div class="result-name">{{ app.name }}</div>
              <div class="result-category">{{ app.category }}</div>
            </div>
            <div class="result-shortcut">⌘⏎</div>
          </div>
        </div>

        <!-- 文件匹配 -->
        <div v-if="matchedFiles.length > 0" class="result-section">
          <div class="section-header">文件</div>
          <div
            v-for="(file, index) in matchedFiles"
            :key="file.id"
            class="result-item"
            :class="{ selected: selectedIndex === matchedApps.length + index }"
            @click="openFile(file)"
          >
            <div class="result-icon">
              <DocumentIcon class="icon-svg" />
            </div>
            <div class="result-info">
              <div class="result-name">{{ file.name }}</div>
              <div class="result-path">{{ file.path }}</div>
            </div>
          </div>
        </div>

        <!-- 设置匹配 -->
        <div v-if="matchedSettings.length > 0" class="result-section">
          <div class="section-header">设置</div>
          <div
            v-for="(setting, index) in matchedSettings"
            :key="setting.id"
            class="result-item"
            :class="{ selected: selectedIndex === matchedApps.length + matchedFiles.length + index }"
            @click="openSetting(setting)"
          >
            <div class="result-icon">
              <CogIcon class="icon-svg" />
            </div>
            <div class="result-info">
              <div class="result-name">{{ setting.name }}</div>
              <div class="result-location">{{ setting.location }}</div>
            </div>
          </div>
        </div>

        <!-- 无结果 -->
        <div v-if="noResults" class="no-results">
          <DocumentTextIcon class="no-results-icon" />
          <div class="no-results-text">未找到结果</div>
          <div class="no-results-hint">尝试不同的关键词</div>
        </div>
      </div>

      <!-- 搜索建议 -->
      <div class="suggestions-container" v-else>
        <div class="suggestions-section">
          <div class="section-header">建议</div>
          <div class="suggestions-grid">
            <div
              v-for="suggestion in suggestions"
              :key="suggestion.id"
              class="suggestion-item"
              @click="executeSuggestion(suggestion)"
            >
              <div class="suggestion-icon">
                <component :is="suggestion.icon" class="icon-svg" />
              </div>
              <div class="suggestion-text">{{ suggestion.text }}</div>
            </div>
          </div>
        </div>

        <div class="quick-actions">
          <div class="section-header">快速操作</div>
          <div class="action-row">
            <div class="action-item" @click="openApp('wifi')">
              <WifiIcon class="action-icon" />
              <span>Wi-Fi</span>
            </div>
            <div class="action-item" @click="openApp('bluetooth')">
              <BeakerIcon class="action-icon" />
              <span>蓝牙</span>
            </div>
            <div class="action-item" @click="openApp('airplane')">
              <PaperAirplaneIcon class="action-icon" />
              <span>飞行模式</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部提示 -->
      <div class="bottom-hints">
        <div class="hint-item">
          <ArrowUpIcon class="hint-icon" />
          <span>↓↑</span>
          <span>导航</span>
        </div>
        <div class="hint-item">
          <ArrowRightIcon class="hint-icon" />
          <span>Enter</span>
          <span>打开</span>
        </div>
        <div class="hint-item">
          <ArrowUturnLeftIcon class="hint-icon" />
          <span>Esc</span>
          <span>关闭</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'
import {
  MagnifyingGlassIcon,
  XMarkIcon,
  DocumentIcon,
  DocumentTextIcon,
  CogIcon,
  WifiIcon,
  BeakerIcon,
  PaperAirplaneIcon,
  ArrowUpIcon,
  ArrowRightIcon,
  ArrowUturnLeftIcon,
  CubeIcon,
  ServerIcon,
  ChartBarIcon,
  UserGroupIcon,
  FolderIcon,
  BellIcon,
  ShieldCheckIcon,
  PhotoIcon,
  MusicalNoteIcon,
  SunIcon,
  MoonIcon
} from '@heroicons/vue/24/outline'

defineEmits(['close', 'openApp'])

const searchQuery = ref('')
const searchInput = ref<HTMLInputElement>()
const selectedIndex = ref(0)

// 应用数据
const apps = ref([
  { id: 'storage', name: '存储管理', icon: ServerIcon, category: '工具' },
  { id: 'monitor', name: '系统监控', icon: ChartBarIcon, category: '系统' },
  { id: 'users', name: '用户管理', icon: UserGroupIcon, category: '系统' },
  { id: 'files', name: '文件管理', icon: FolderIcon, category: '工具' },
  { id: 'settings', name: '系统设置', icon: CogIcon, category: '系统' },
  { id: 'notifications', name: '通知中心', icon: BellIcon, category: '系统' },
  { id: 'sso', name: 'SSO登录', icon: ShieldCheckIcon, category: '系统' },
  { id: 'photos', name: '照片', icon: PhotoIcon, category: '媒体' },
  { id: 'music', name: '音乐', icon: MusicalNoteIcon, category: '媒体' },
  { id: 'docker', name: 'Docker管理', icon: CubeIcon, category: '开发' },
])

// 文件数据
const files = ref([
  { id: 'doc1', name: '项目文档.txt', path: '~/Documents/Projects' },
  { id: 'doc2', name: '系统配置.json', path: '~/System/Config' },
  { id: 'doc3', name: '备份计划.md', path: '~/Backup' },
])

// 设置数据
const settings = ref([
  { id: 'wifi', name: 'Wi-Fi设置', location: '设置 > 网络与无线' },
  { id: 'display', name: '显示与亮度', location: '设置 > 显示' },
  { id: 'sound', name: '声音与触感', location: '设置 > 声音' },
  { id: 'privacy', name: '隐私与安全', location: '设置 > 隐私' },
])

// 搜索建议
const suggestions = ref([
  { id: 'weather', text: '天气', icon: SunIcon },
  { id: 'stocks', text: '股票', icon: ChartBarIcon },
  { id: 'calculator', text: '计算器', icon: CogIcon },
  { id: 'notes', text: '备忘录', icon: DocumentTextIcon },
])

// 计算属性
const matchedApps = computed(() => {
  if (!searchQuery.value) return []
  const query = searchQuery.value.toLowerCase()
  return apps.value.filter(app =>
    app.name.toLowerCase().includes(query) ||
    app.category.toLowerCase().includes(query)
  )
})

const matchedFiles = computed(() => {
  if (!searchQuery.value) return []
  const query = searchQuery.value.toLowerCase()
  return files.value.filter(file =>
    file.name.toLowerCase().includes(query) ||
    file.path.toLowerCase().includes(query)
  )
})

const matchedSettings = computed(() => {
  if (!searchQuery.value) return []
  const query = searchQuery.value.toLowerCase()
  return settings.value.filter(setting =>
    setting.name.toLowerCase().includes(query) ||
    setting.location.toLowerCase().includes(query)
  )
})

const noResults = computed(() => {
  return searchQuery.value &&
         matchedApps.value.length === 0 &&
         matchedFiles.value.length === 0 &&
         matchedSettings.value.length === 0
})

const totalResults = computed(() => {
  return matchedApps.value.length +
         matchedFiles.value.length +
         matchedSettings.value.length
})

// 方法
const clearSearch = () => {
  searchQuery.value = ''
  selectedIndex.value = 0
}

const openApp = (app: any) => {
  console.log('Opening app:', app.id)
  // 集成应用打开逻辑
  // emit('openApp', app)
}

const openFile = (file: any) => {
  console.log('Opening file:', file.id)
}

const openSetting = (setting: any) => {
  console.log('Opening setting:', setting.id)
}

const executeSuggestion = (suggestion: any) => {
  searchQuery.value = suggestion.text
}

const selectFirstResult = () => {
  if (matchedApps.value.length > 0) {
    openApp(matchedApps.value[0])
  } else if (matchedFiles.value.length > 0) {
    openFile(matchedFiles.value[0])
  } else if (matchedSettings.value.length > 0) {
    openSetting(matchedSettings.value[0])
  }
}

const navigateUp = () => {
  if (selectedIndex.value > 0) {
    selectedIndex.value--
  }
}

const navigateDown = () => {
  if (selectedIndex.value < totalResults.value - 1) {
    selectedIndex.value++
  }
}

// 生命周期
onMounted(() => {
  nextTick(() => {
    searchInput.value?.focus()
  })
})

// 键盘快捷键
const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    // emit('close')
  }
}
</script>

<style scoped>
.ios-spotlight {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(20px);
  z-index: 4000;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 120px;
}

.spotlight-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(40px);
  border-radius: 16px;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.4);
  width: 90%;
  max-width: 680px;
  max-height: 600px;
  overflow: hidden;
  animation: spotlightIn 0.2s ease-out;
}

@keyframes spotlightIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-20px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.search-container {
  display: flex;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  background: rgba(255, 255, 255, 0.8);
}

.search-icon {
  width: 20px;
  height: 20px;
  margin-right: 12px;
  color: #666;
}

.search-input {
  flex: 1;
  height: 24px;
  border: none;
  background: transparent;
  font-size: 22px;
  font-weight: 300;
  color: #333;
  outline: none;
}

.search-input::placeholder {
  color: #999;
}

.clear-btn {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.clear-btn:hover {
  background: rgba(0, 0, 0, 0.1);
}

.clear-icon {
  width: 16px;
  height: 16px;
  color: #666;
}

.results-container,
.suggestions-container {
  max-height: 500px;
  overflow-y: auto;
  padding: 8px 0;
}

.result-section,
.suggestions-section {
  margin-bottom: 16px;
}

.section-header {
  font-size: 13px;
  font-weight: 600;
  color: #666;
  padding: 8px 16px 8px 48px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.result-item:hover,
.result-item.selected {
  background: rgba(59, 130, 246, 0.1);
}

.result-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.result-icon .icon-svg {
  width: 20px;
  height: 20px;
  color: white;
}

.result-info {
  flex: 1;
  min-width: 0;
}

.result-name {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  margin-bottom: 2px;
}

.result-category,
.result-path,
.result-location {
  font-size: 12px;
  color: #666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-shortcut {
  font-size: 12px;
  color: #999;
  background: rgba(0, 0, 0, 0.05);
  padding: 4px 8px;
  border-radius: 6px;
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #999;
}

.no-results-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.no-results-text {
  font-size: 18px;
  font-weight: 500;
  margin-bottom: 8px;
}

.no-results-hint {
  font-size: 14px;
  color: #999;
}

.suggestions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  padding: 0 16px;
}

.suggestion-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.suggestion-item:hover {
  background: rgba(0, 0, 0, 0.06);
}

.suggestion-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.suggestion-icon .icon-svg {
  width: 18px;
  height: 18px;
  color: white;
}

.suggestion-text {
  font-size: 15px;
  color: #333;
}

.quick-actions {
  padding: 0 16px;
}

.action-row {
  display: flex;
  gap: 8px;
  margin-top: 8px;
}

.action-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px;
  background: rgba(0, 0, 0, 0.03);
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-item:hover {
  background: rgba(0, 0, 0, 0.06);
}

.action-icon {
  width: 24px;
  height: 24px;
  color: #333;
}

.action-item span {
  font-size: 12px;
  color: #666;
}

.bottom-hints {
  display: flex;
  justify-content: center;
  gap: 32px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.03);
  border-top: 1px solid rgba(0, 0, 0, 0.1);
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #666;
}

.hint-icon {
  width: 16px;
  height: 16px;
}

/* 响应式 */
@media (max-width: 768px) {
  .spotlight-container {
    width: 95%;
  }

  .suggestions-grid {
    grid-template-columns: 1fr;
  }

  .bottom-hints {
    flex-wrap: wrap;
    gap: 16px;
  }
}
</style>