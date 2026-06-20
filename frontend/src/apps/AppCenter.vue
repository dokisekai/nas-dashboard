<template>
  <div class="app-center">
    <div class="app-header">
      <h1>应用中心</h1>
      <p class="subtitle">浏览和安装应用程序</p>
    </div>

    <div class="app-search">
      <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="搜索应用..."
        class="search-input"
      />
    </div>

    <div class="app-categories">
      <button
        v-for="category in categories"
        :key="category.id"
        :class="['category-btn', { active: selectedCategory === category.id }]"
        @click="selectedCategory = category.id"
      >
        {{ category.label }}
      </button>
    </div>

    <div class="apps-grid">
      <div
        v-for="app in filteredApps"
        :key="app.id"
        class="app-card"
        :class="{ 'app-installed': isAppInstalled(app.id) }"
      >
        <div class="app-icon" :style="{ background: app.color }">
          <component :is="getAppIcon(app.icon) as any" class="w-8 h-8" />
        </div>
        <h3 class="app-name">{{ app.name }}</h3>
        <p class="app-description">{{ app.description }}</p>
        <div class="app-meta">
          <span class="app-type-badge" :class="app.appType">
            {{ app.appType === 'docker' ? 'Docker' : '系统' }}
          </span>
          <span class="app-category">{{ getCategoryName(app.category) }}</span>
          <span class="app-version">v{{ app.version }}</span>
          <span class="app-rating" v-if="app.rating">
            <StarIcon class="w-4 h-4" />
            {{ app.rating }}
          </span>
        </div>
        <div class="app-actions">
          <button
            v-if="isAppInstalled(app.id)"
            class="action-btn secondary"
            @click="launchApp(app)"
            :disabled="!app.launching"
          >
            <PlayIcon class="w-4 h-4" />
            {{ app.launching ? '打开' : '已安装' }}
          </button>
          <button
            v-else
            class="action-btn primary"
            @click="installApp(app)"
            :disabled="showInstallProgress && selectedApp?.id === app.id"
          >
            <ArrowDownTrayIcon class="w-4 h-4" />
            {{ (showInstallProgress && selectedApp?.id === app.id) ? '安装中...' : '安装' }}
          </button>
          <button
            v-if="isAppInstalled(app.id)"
            class="action-btn danger"
            @click="uninstallApp(app)"
            :disabled="showInstallProgress && selectedApp?.id === app.id"
          >
            <TrashIcon class="w-4 h-4" />
            {{ (showInstallProgress && selectedApp?.id === app.id) ? '卸载中...' : '卸载' }}
          </button>
          <button class="action-btn" @click="viewAppDetails(app)">
            <InformationCircleIcon class="w-4 h-4" />
            详情
          </button>
        </div>
      </div>
    </div>

    <!-- App Details Modal -->
    <div v-if="showAppDetails" class="modal-overlay" @click="closeAppDetails">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <div class="app-detail-header">
            <div class="app-detail-icon" :style="{ background: selectedApp?.color }">
              <component :is="getAppIcon(selectedApp?.icon || '') as any" class="w-12 h-12" />
            </div>
            <div class="app-detail-info">
              <h2>{{ selectedApp?.name }}</h2>
              <p class="app-detail-version">版本 {{ selectedApp?.version }}</p>
              <div class="app-detail-meta">
                <span class="app-rating" v-if="selectedApp?.rating">
                  <StarIcon class="w-5 h-5" />
                  {{ selectedApp.rating }}
                  <span class="rating-count">({{ selectedApp.ratingCount || 0 }} 评价)</span>
                </span>
                <span class="app-downloads" v-if="selectedApp?.downloads">
                  <CloudArrowDownIcon class="w-5 h-5" />
                  {{ formatDownloads(selectedApp.downloads) }} 下载
                </span>
              </div>
            </div>
          </div>
          <button class="close-btn" @click="closeAppDetails">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>

        <div class="modal-body">
          <!-- Screenshots -->
          <div class="app-screenshots" v-if="selectedApp?.screenshots?.length">
            <h3>截图</h3>
            <div class="screenshots-grid">
              <div
                v-for="(screenshot, index) in selectedApp.screenshots"
                :key="index"
                class="screenshot-item"
              >
                <img :src="screenshot" :alt="`${selectedApp?.name || 'app'} screenshot ${Number(index) + 1}`" />
              </div>
            </div>
          </div>

          <!-- Description -->
          <div class="app-description-section">
            <h3>应用描述</h3>
            <p>{{ selectedApp?.fullDescription || selectedApp?.description }}</p>
          </div>

          <!-- Features -->
          <div class="app-features" v-if="selectedApp?.features?.length">
            <h3>主要功能</h3>
            <ul>
              <li v-for="feature in selectedApp.features" :key="feature">
                {{ feature }}
              </li>
            </ul>
          </div>

          <!-- Requirements -->
          <div class="app-requirements" v-if="selectedApp?.requirements">
            <h3>系统要求</h3>
            <div class="requirements-grid">
              <div class="requirement-item">
                <span class="req-label">内存:</span>
                <span class="req-value">{{ selectedApp.requirements.memory || '2GB+' }}</span>
              </div>
              <div class="requirement-item">
                <span class="req-label">存储:</span>
                <span class="req-value">{{ selectedApp.requirements.storage || '100MB+' }}</span>
              </div>
              <div class="requirement-item">
                <span class="req-label">依赖:</span>
                <span class="req-value">{{ selectedApp.requirements.dependencies || '无' }}</span>
              </div>
            </div>
          </div>

          <!-- Changelog -->
          <div class="app-changelog" v-if="selectedApp?.changelog?.length">
            <h3>更新日志</h3>
            <div class="changelog-list">
              <div
                v-for="(change, index) in selectedApp.changelog"
                :key="index"
                class="changelog-item"
              >
                <div class="changelog-version">{{ change.version }}</div>
                <div class="changelog-date">{{ formatDate(change.date) }}</div>
                <ul class="changelog-changes">
                  <li v-for="(item, i) in change.changes" :key="i">{{ item }}</li>
                </ul>
              </div>
            </div>
          </div>

          <!-- Reviews -->
          <div class="app-reviews" v-if="selectedApp?.reviews?.length">
            <h3>用户评价</h3>
            <div class="reviews-list">
              <div
                v-for="review in selectedApp.reviews.slice(0, 3)"
                :key="review.id"
                class="review-item"
              >
                <div class="review-header">
                  <div class="review-author">{{ review.author }}</div>
                  <div class="review-rating">
                    <StarIcon class="w-4 h-4" />
                    {{ review.rating }}
                  </div>
                </div>
                <p class="review-text">{{ review.text }}</p>
                <div class="review-date">{{ formatDate(review.date) }}</div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <div class="footer-info">
            <span class="developer">开发者: {{ selectedApp?.developer || 'Unknown' }}</span>
            <span class="license">许可证: {{ selectedApp?.license || 'MIT' }}</span>
          </div>
          <div class="footer-actions">
            <button class="action-btn" @click="closeAppDetails">
              关闭
            </button>
            <button
              class="action-btn primary"
              :class="{ uninstall: selectedApp?.installed }"
              @click="selectedApp?.installed ? uninstallApp(selectedApp) : installApp(selectedApp)"
            >
              {{ selectedApp?.installed ? '卸载应用' : '安装应用' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Installation Progress Modal -->
    <div v-if="showInstallProgress" class="modal-overlay">
      <div class="modal-content small" @click.stop>
        <div class="install-progress">
          <div class="progress-icon">
            <CloudArrowDownIcon class="w-12 h-12" v-if="installStep === 'downloading'" />
            <CogIcon class="w-12 h-8 spinning" v-if="installStep === 'installing'" />
            <CheckIcon class="w-12 h-12" v-if="installStep === 'completed'" />
            <XMarkIcon class="w-12 h-12" v-if="installStep === 'error'" />
          </div>
          <h3>{{ installTitle }}</h3>
          <p>{{ installMessage }}</p>
          <div class="progress-bar" v-if="installProgress > 0 && installProgress < 100">
            <div class="progress-fill" :style="{ width: installProgress + '%' }"></div>
          </div>
          <div class="progress-text" v-if="installProgress > 0 && installProgress < 100">
            {{ installProgress }}%
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  ShoppingBagIcon,
  ServerIcon,
  ChartBarIcon,
  FolderIcon,
  UserGroupIcon,
  CogIcon,
  CloudIcon,
  StarIcon,
  CloudArrowDownIcon,
  XMarkIcon,
  CheckIcon,
  PlayIcon,
  ArrowDownTrayIcon,
  TrashIcon,
  InformationCircleIcon
} from '@heroicons/vue/24/outline'

const searchQuery = ref('')
const selectedCategory = ref('all')
const showAppDetails = ref(false)
const selectedApp = ref<any>(null)
const showInstallProgress = ref(false)
const installStep = ref<'downloading' | 'installing' | 'completed' | 'error'>('downloading')
const installProgress = ref(0)
const installTitle = ref('')
const installMessage = ref('')

// 用于追踪正在进行的操作
const installingApp = ref<string | null>(null)
const uninstallingApp = ref<string | null>(null)

const categories = [
  { id: 'all', label: '全部' },
  { id: 'docker', label: 'Docker应用' },
  { id: 'system', label: '系统工具' },
  { id: 'storage', label: '存储管理' },
  { id: 'media', label: '媒体服务' },
  { id: 'network', label: '网络工具' },
  { id: 'security', label: '安全工具' }
]

const apps = ref([
  {
    id: 'storage-manager',
    name: '存储管理器',
    description: '管理磁盘和存储卷',
    icon: 'ServerIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
    category: 'storage',
    appType: 'system',
    version: '2.1.0',
    installed: true,
    launching: true,
    appId: 'storage-manager',
    rating: 4.8,
    ratingCount: 234,
    downloads: 15234,
    developer: 'NAS Team',
    license: 'MIT'
  },
  {
    id: 'system-monitor',
    name: '系统监控',
    description: '实时系统资源监控',
    icon: 'ChartBarIcon',
    color: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    category: 'system',
    appType: 'system',
    version: '1.8.2',
    installed: true,
    launching: true,
    appId: 'system-monitor',
    rating: 4.9,
    ratingCount: 456,
    downloads: 28901,
    developer: 'NAS Team',
    license: 'Apache 2.0'
  },
  {
    id: 'immich',
    name: 'Immich',
    description: '高性能自托管照片和视频备份方案',
    icon: 'CloudIcon',
    color: 'linear-gradient(135deg, #3b82f6 0%, #2dd4bf 100%)',
    category: 'media',
    appType: 'docker',
    version: '1.105.1',
    installed: false,
    launching: false,
    appId: '',
    rating: 4.9,
    ratingCount: 1250,
    downloads: 45000,
    developer: 'Immich Team',
    license: 'MIT'
  },
  {
    id: 'uptime-kuma',
    name: 'Uptime Kuma',
    description: '一个花哨的自托管监控工具',
    icon: 'ChartBarIcon',
    color: 'linear-gradient(135deg, #10b981 0%, #3b82f6 100%)',
    category: 'network',
    appType: 'docker',
    version: '1.23.11',
    installed: false,
    launching: false,
    appId: '',
    rating: 4.8,
    ratingCount: 850,
    downloads: 32000,
    developer: 'louislam',
    license: 'MIT'
  },
  {
    id: 'file-manager',
    name: '文件管理器',
    description: '文件浏览和管理',
    icon: 'FolderIcon',
    color: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    category: 'system',
    appType: 'system',
    version: '3.0.1',
    installed: true,
    launching: true,
    appId: 'file-browser',
    rating: 4.7,
    ratingCount: 189,
    downloads: 12345,
    developer: 'NAS Team',
    license: 'MIT'
  },
  {
    id: 'user-manager',
    name: '用户管理',
    description: '管理用户和权限',
    icon: 'UserGroupIcon',
    color: 'linear-gradient(135deg, #ef4444 0%, #f87171 100%)',
    category: 'system',
    appType: 'system',
    version: '1.5.0',
    installed: true,
    launching: true,
    appId: 'user-manager',
    rating: 4.6,
    ratingCount: 98,
    downloads: 8765,
    developer: 'NAS Team',
    license: 'GPL-3.0'
  },
  {
    id: 'cloud-sync',
    name: '云同步',
    description: '文件云同步服务',
    icon: 'CloudIcon',
    color: 'linear-gradient(135deg, #06b6d4 0%, #38bdf8 100%)',
    category: 'network',
    appType: 'system',
    version: '2.0.3',
    installed: false,
    launching: false,
    appId: '',
    rating: 4.5,
    ratingCount: 67,
    downloads: 5432,
    developer: 'Cloud Team',
    license: 'MIT'
  },
  {
    id: 'firewall',
    name: '防火墙',
    description: '网络安全防护',
    icon: 'CogIcon',
    color: 'linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%)',
    category: 'security',
    appType: 'system',
    version: '1.2.0',
    installed: false,
    launching: false,
    appId: '',
    rating: 4.4,
    ratingCount: 45,
    downloads: 3210,
    developer: 'Security Team',
    license: 'GPL-3.0'
  }
])

const filteredApps = computed(() => {
  return apps.value.filter(app => {
    const matchesSearch = app.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
                         app.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesCategory = selectedCategory.value === 'all' || 
                           app.category === selectedCategory.value ||
                           (selectedCategory.value === 'docker' && app.appType === 'docker') ||
                           (selectedCategory.value === 'system' && app.appType === 'system')
    return matchesSearch && matchesCategory
  })
})

const getAppIcon = (iconName: any): any => {
  // 如果已经是组件对象，直接返回
  if (typeof iconName === 'object') return iconName

  // 否则从图标映射中查找
  const icons: Record<string, any> = {
    ShoppingBagIcon,
    ServerIcon,
    ChartBarIcon,
    FolderIcon,
    UserGroupIcon,
    CogIcon,
    CloudIcon
  }
  return icons[iconName] || ShoppingBagIcon
}

const getCategoryName = (categoryId: string) => {
  const category = categories.find(c => c.id === categoryId)
  return category?.label || categoryId
}

const viewAppDetails = (app: any) => {
  selectedApp.value = app
  showAppDetails.value = true
}

// 检查应用是否已安装
const isAppInstalled = (appId: string) => {
  const app = apps.value.find(a => a.id === appId)
  return app?.installed || false
}

// 启动应用
const launchApp = (app: any) => {
  if (!app.launching || !app.appId) {
    alert('此应用暂不支持启动')
    return
  }

  // 通过事件总线发送启动应用的消息
  const event = new CustomEvent('launch-app', {
    detail: { appId: app.appId, appName: app.name }
  })
  window.dispatchEvent(event)
}

const closeAppDetails = () => {
  showAppDetails.value = false
  selectedApp.value = null
}

const installApp = async (app: any) => {
  installingApp.value = app.id
  selectedApp.value = app
  showInstallProgress.value = true
  installStep.value = 'downloading'
  installTitle.value = '正在下载 ' + app.name
  installMessage.value = '正在下载应用文件...'
  installProgress.value = 0

  // Simulate download progress
  const downloadInterval = setInterval(() => {
    if (installProgress.value < 100) {
      installProgress.value += 10
    } else {
      clearInterval(downloadInterval)
      installStep.value = 'installing'
      installTitle.value = '正在安装 ' + app.name
      installMessage.value = '正在安装应用...'

      // Simulate installation
      setTimeout(() => {
        installStep.value = 'completed'
        installTitle.value = '安装完成'
        installMessage.value = app.name + ' 已成功安装！'
        app.installed = true
        installingApp.value = null

        setTimeout(() => {
          showInstallProgress.value = false
          installProgress.value = 0
        }, 2000)
      }, 2000)
    }
  }, 300)
}

const uninstallApp = async (app: any) => {
  if (confirm(`确定要卸载 "${app.name}" 吗?`)) {
    uninstallingApp.value = app.id
    try {
      // Simulate uninstallation
      setTimeout(() => {
        app.installed = false
        uninstallingApp.value = null
        alert(app.name + ' 已成功卸载')
      }, 1500)
    } catch (error: any) {
      console.error('Failed to uninstall app:', error)
      uninstallingApp.value = null
      alert('卸载失败: ' + error.message)
    }
  }
}

const formatDownloads = (count: number | string) => {
  const numCount = typeof count === 'string' ? parseInt(count) : count
  if (numCount >= 1000000) {
    return (numCount / 1000000).toFixed(1) + 'M'
  } else if (numCount >= 1000) {
    return (numCount / 1000).toFixed(1) + 'K'
  }
  return numCount.toString()
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.app-center {
  width: 100%;
  height: 100%;
  padding: 32px;
  background: linear-gradient(135deg, #f5f3ff 0%, #ede9fe 100%);
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.app-header {
  margin-bottom: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.app-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.app-search {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  margin-bottom: 24px;
}

.search-icon {
  width: 20px;
  height: 20px;
  color: #9ca3af;
}

.search-input {
  flex: 1;
  border: none;
  font-size: 16px;
  color: #1f2937;
  outline: none;
}

.search-input::placeholder {
  color: #9ca3af;
}

.app-categories {
  display: flex;
  gap: 8px;
  margin-bottom: 32px;
  flex-wrap: wrap;
}

.category-btn {
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 8px;
  font-size: 14px;
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.category-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  border-color: rgba(102, 126, 234, 0.3);
  color: #764ba2;
}

.category-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
}

.apps-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.app-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  flex-direction: column;
}

.app-card:hover {
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
  transform: translateY(-4px);
}

.app-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  margin-bottom: 16px;
}

.app-name {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.app-description {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 16px;
  flex: 1;
}

.app-meta {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.app-category {
  padding: 4px 8px;
  background: #f3f4f6;
  border-radius: 4px;
  font-size: 12px;
  color: #6b7280;
}

.app-type-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
}

.app-type-badge.docker {
  background: #eff6ff;
  color: #3b82f6;
  border: 1px solid #bfdbfe;
}

.app-type-badge.system {
  background: #f0fdf4;
  color: #16a34a;
  border: 1px solid #bbf7d0;
}

.app-version {
  padding: 4px 8px;
  background: #e5e7eb;
  border-radius: 4px;
  font-size: 12px;
  color: #6b7280;
}

.install-btn {
  width: 100%;
  padding: 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.install-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.install-btn.installed {
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
  cursor: default;
}

.install-btn.uninstall {
  background: linear-gradient(135deg, #ef4444 0%, #f87171 100%);
}

.install-btn.uninstall:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 8px;
  font-size: 14px;
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #764ba2;
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
}

.action-btn.primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.action-btn.secondary {
  flex: 1;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  animation: modalSlideIn 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-content.large {
  width: 800px;
}

.modal-content.small {
  width: 400px;
  padding: 32px;
  text-align: center;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.app-detail-header {
  display: flex;
  gap: 20px;
  flex: 1;
}

.app-detail-icon {
  width: 80px;
  height: 80px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.app-detail-info h2 {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.app-detail-version {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 12px;
}

.app-detail-meta {
  display: flex;
  gap: 16px;
}

.app-rating {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #f59e0b;
}

.rating-count {
  color: #6b7280;
  font-size: 12px;
}

.app-downloads {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  font-size: 14px;
}

.close-btn {
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  padding: 4px;
}

.close-btn:hover {
  color: #1f2937;
}

.modal-body {
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.modal-body h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.app-screenshots .screenshots-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.screenshot-item {
  border-radius: 12px;
  overflow: hidden;
}

.screenshot-item img {
  width: 100%;
  height: auto;
  display: block;
}

.app-description-section p {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.6;
}

.app-features ul {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.app-features li {
  font-size: 14px;
  color: #6b7280;
  padding-left: 20px;
  position: relative;
}

.app-features li::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: #10b981;
  font-weight: bold;
}

.requirements-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.requirement-item {
  display: flex;
  justify-content: space-between;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
}

.req-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.req-value {
  font-size: 14px;
  color: #1f2937;
}

.changelog-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.changelog-item {
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.changelog-version {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.changelog-date {
  font-size: 12px;
  color: #9ca3af;
  margin-bottom: 12px;
}

.changelog-changes {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.changelog-changes li {
  font-size: 14px;
  color: #6b7280;
  padding-left: 16px;
  position: relative;
}

.changelog-changes li::before {
  content: '•';
  position: absolute;
  left: 0;
  color: #3b82f6;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.review-item {
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.review-author {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.review-rating {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #f59e0b;
  font-size: 14px;
}

.review-text {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 8px;
  line-height: 1.5;
}

.review-date {
  font-size: 12px;
  color: #9ca3af;
}

.modal-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-top: 1px solid #e5e7eb;
}

.footer-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.footer-actions {
  display: flex;
  gap: 12px;
}

.install-progress {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.progress-icon {
  width: 48px;
  height: 48px;
  color: #3b82f6;
}

.progress-icon svg.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.install-progress h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.install-progress p {
  font-size: 14px;
  color: #6b7280;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #8b5cf6 100%);
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 12px;
  color: #6b7280;
  margin-top: 8px;
}
</style>