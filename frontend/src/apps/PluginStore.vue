<template>
  <div class="plugin-store">
    <div class="store-header">
      <h1>插件商店</h1>
      <p class="subtitle">发现并安装社区插件</p>
    </div>

    <div class="store-layout">
      <!-- 侧边栏 -->
      <div class="store-sidebar">
        <div class="sidebar-section">
          <h3>分类</h3>
          <div class="category-list">
            <button
              v-for="category in categories"
              :key="category.id"
              :class="['category-item', { active: selectedCategory === category.id }]"
              @click="selectedCategory = category.id"
            >
              <component :is="getCategoryIcon(category.icon)" class="category-icon" />
              {{ category.label }}
              <span class="category-count">{{ category.count }}</span>
            </button>
          </div>
        </div>

        <div class="sidebar-section">
          <h3>我的插件</h3>
          <div class="installed-plugins">
            <div
              v-for="plugin in installedPlugins"
              :key="plugin.id"
              class="installed-plugin-item"
              @click="showPluginDetails(plugin)"
            >
              <div class="plugin-icon" :style="{ background: plugin.color }">
                <PuzzlePieceIcon class="w-4 h-4" />
              </div>
              <span class="plugin-name">{{ plugin.name }}</span>
              <button class="plugin-action" @click.stop="openPlugin(plugin)">
                <CogIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 主内容区 -->
      <div class="store-content">
        <!-- 搜索和筛选 -->
        <div class="content-header">
          <div class="search-bar">
            <MagnifyingGlassIcon class="search-icon" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索插件..."
              class="search-input"
            />
          </div>
          <div class="sort-options">
            <select v-model="sortBy" class="sort-select">
              <option value="popular">最受欢迎</option>
              <option value="rating">评分最高</option>
              <option value="newest">最新发布</option>
              <option value="updated">最近更新</option>
            </select>
          </div>
        </div>

        <!-- 特色插件 -->
        <div v-if="selectedCategory === 'all' && !searchQuery" class="featured-section">
          <h2>精选插件</h2>
          <div class="featured-plugins">
            <div
              v-for="plugin in featuredPlugins"
              :key="plugin.id"
              class="featured-plugin-card"
              @click="showPluginDetails(plugin)"
            >
              <div class="featured-badge">推荐</div>
              <div class="plugin-icon large" :style="{ background: plugin.color }">
                <component :is="plugin.iconComponent" class="w-12 h-12" />
              </div>
              <h3>{{ plugin.name }}</h3>
              <p>{{ plugin.shortDescription }}</p>
              <div class="featured-stats">
                <span class="rating">
                  <StarIcon class="w-4 h-4" />
                  {{ plugin.rating }}
                </span>
                <span class="downloads">{{ formatNumber(plugin.downloadCount) }} 下载</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 插件网格 -->
        <div class="plugins-section">
          <h2 v-if="selectedCategory === 'all'">所有插件</h2>
          <h2 v-else>{{ getCategoryLabel(selectedCategory) }}</h2>

          <div v-if="loading" class="loading-state">
            <div class="spinner"></div>
            <p>加载中...</p>
          </div>

          <div v-else-if="filteredPlugins.length === 0" class="empty-state">
            <FolderIcon class="empty-icon" />
            <p>未找到相关插件</p>
          </div>

          <div v-else class="plugins-grid">
            <div
              v-for="plugin in filteredPlugins"
              :key="plugin.id"
              class="plugin-card"
            >
              <div class="plugin-header">
                <div class="plugin-icon" :style="{ background: plugin.color }">
                  <component :is="plugin.iconComponent" class="w-8 h-8" />
                </div>
                <button class="favorite-btn" @click="toggleFavorite(plugin)">
                  <HeartIcon
                    class="w-5 h-5"
                    :class="{ 'is-favorite': plugin.isFavorite }"
                  />
                </button>
              </div>

              <h3 class="plugin-name" @click="showPluginDetails(plugin)">{{ plugin.name }}</h3>
              <p class="plugin-description">{{ plugin.shortDescription }}</p>
              <p class="plugin-author">by {{ plugin.author }}</p>

              <div class="plugin-meta">
                <div class="plugin-rating">
                  <StarIcon
                    v-for="i in 5"
                    :key="i"
                    class="star-icon"
                    :class="{ 'is-filled': i <= Math.floor(plugin.rating) }"
                  />
                  <span class="rating-value">{{ plugin.rating.toFixed(1) }}</span>
                  <span class="rating-count">({{ plugin.reviewCount }})</span>
                </div>
                <div class="plugin-stats">
                  <span>{{ formatNumber(plugin.downloadCount) }} 下载</span>
                </div>
              </div>

              <div class="plugin-actions">
                <button class="action-btn secondary" @click="showPluginDetails(plugin)">
                  <EyeIcon class="w-4 h-4" />
                  详情
                </button>
                <button
                  class="install-btn"
                  :class="{ installed: plugin.installed }"
                  @click="toggleInstall(plugin)"
                >
                  {{ plugin.installed ? '已安装' : '安装' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Plugin Details Modal -->
    <div v-if="showDetailsModal" class="modal-overlay" @click="closeDetailsModal">
      <div class="modal-content large" @click.stop>
        <div class="modal-header">
          <div class="plugin-detail-header">
            <div class="plugin-icon large" :style="{ background: selectedPlugin?.color }">
              <component :is="selectedPlugin?.iconComponent" class="w-12 h-12" />
            </div>
            <div class="plugin-detail-info">
              <h2>{{ selectedPlugin?.name }}</h2>
              <p class="plugin-version">版本 {{ selectedPlugin?.version }}</p>
              <p class="plugin-author">by {{ selectedPlugin?.author }}</p>
              <div class="plugin-detail-meta">
                <div class="rating">
                  <StarIcon class="w-5 h-5" />
                  {{ selectedPlugin?.rating }}
                  <span class="rating-count">({{ selectedPlugin?.reviewCount }} 评价)</span>
                </div>
                <div class="downloads">
                  <CloudArrowDownIcon class="w-5 h-5" />
                  {{ formatNumber(selectedPlugin?.downloadCount || 0) }} 下载
                </div>
              </div>
            </div>
          </div>
          <button class="close-btn" @click="closeDetailsModal">
            <XMarkIcon class="w-6 h-6" />
          </button>
        </div>

        <div class="modal-body">
          <!-- Description -->
          <div class="plugin-description-section">
            <h3>插件描述</h3>
            <p>{{ selectedPlugin?.fullDescription || selectedPlugin?.shortDescription }}</p>
          </div>

          <!-- Features -->
          <div class="plugin-features" v-if="selectedPlugin?.features?.length">
            <h3>主要功能</h3>
            <ul>
              <li v-for="feature in selectedPlugin.features" :key="feature">
                {{ feature }}
              </li>
            </ul>
          </div>

          <!-- Installation -->
          <div class="plugin-installation">
            <h3>安装说明</h3>
            <div class="install-steps">
              <div class="step">
                <div class="step-number">1</div>
                <p>点击"安装"按钮下载插件</p>
              </div>
              <div class="step">
                <div class="step-number">2</div>
                <p>等待安装完成</p>
              </div>
              <div class="step">
                <div class="step-number">3</div>
                <p>在插件管理中配置和使用</p>
              </div>
            </div>
          </div>

          <!-- Reviews -->
          <div class="plugin-reviews" v-if="selectedPlugin?.reviews?.length">
            <h3>用户评价</h3>
            <div class="reviews-list">
              <div
                v-for="review in selectedPlugin.reviews.slice(0, 3)"
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

          <!-- Development Docs -->
          <div class="plugin-dev-docs">
            <h3>开发者文档</h3>
            <div class="dev-links">
              <a href="#" class="dev-link">
                <CodeBracketIcon class="w-5 h-5" />
                API 文档
              </a>
              <a href="#" class="dev-link">
                <BookOpenIcon class="w-5 h-5" />
                开发指南
              </a>
              <a href="#" class="dev-link">
                <BeakerIcon class="w-5 h-5" />
                示例代码
              </a>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <div class="footer-info">
            <span class="license">许可证: {{ selectedPlugin?.license || 'MIT' }}</span>
            <span class="updated">更新于: {{ formatDate(selectedPlugin?.lastUpdated || '') }}</span>
          </div>
          <div class="footer-actions">
            <button class="action-btn" @click="closeDetailsModal">
              关闭
            </button>
            <button
              class="action-btn primary"
              :class="{ uninstall: selectedPlugin?.installed }"
              @click="toggleInstall(selectedPlugin)"
            >
              {{ selectedPlugin?.installed ? '卸载插件' : '安装插件' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  CogIcon,
  StarIcon,
  HeartIcon,
  ServerIcon,
  FolderIcon,
  MagnifyingGlassIcon,
  EyeIcon,
  CloudArrowDownIcon,
  XMarkIcon,
  PuzzlePieceIcon,
  CodeBracketIcon,
  BookOpenIcon,
  BeakerIcon
} from '@heroicons/vue/24/outline'

const selectedCategory = ref('all')
const searchQuery = ref('')
const sortBy = ref('popular')
const loading = ref(false)
const showDetailsModal = ref(false)
const selectedPlugin = ref<any>(null)

const categories = [
  { id: 'all', label: '全部插件', icon: 'FolderIcon', count: 156 },
  { id: 'system', label: '系统工具', icon: 'ServerIcon', count: 42 },
  { id: 'security', label: '安全工具', icon: 'ShieldIcon', count: 15 },
  { id: 'storage', label: '存储工具', icon: 'FolderIcon', count: 18 },
  { id: 'network', label: '网络工具', icon: 'GlobeAltIcon', count: 22 },
  { id: 'media', label: '媒体工具', icon: 'PhotoIcon', count: 28 },
  { id: 'backup', label: '备份工具', icon: 'ArchiveIcon', count: 11 }
]

const plugins = ref([
  {
    id: 'docker-manager',
    name: 'Docker管理器',
    shortDescription: '管理Docker容器和镜像',
    fullDescription: '功能强大的Docker管理工具，提供容器创建、监控、日志查看等功能。支持容器编排和自动化部署。',
    version: '1.2.0',
    author: 'NAS Team',
    rating: 4.8,
    reviewCount: 234,
    downloadCount: 15234,
    category: 'system',
    iconComponent: ServerIcon,
    color: 'linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%)',
    isFavorite: false,
    installed: true,
    license: 'MIT',
    lastUpdated: '2024-01-15',
    features: [
      '容器管理',
      '镜像管理',
      '网络配置',
      '卷管理',
      '日志查看',
      '资源监控'
    ],
    reviews: [
      {
        id: '1',
        author: 'John Doe',
        rating: 5,
        text: '非常好用的Docker管理工具！',
        date: '2024-01-10'
      }
    ]
  },
  {
    id: 'cloud-backup',
    name: '云备份',
    shortDescription: '自动备份到云端存储',
    fullDescription: '支持多种云存储的自动备份解决方案，确保数据安全。支持增量备份和版本控制。',
    version: '2.0.1',
    author: 'CloudSync Inc',
    rating: 4.6,
    reviewCount: 189,
    downloadCount: 8934,
    category: 'backup',
    iconComponent: CloudArrowDownIcon,
    color: 'linear-gradient(135deg, #10b981 0%, #34d399 100%)',
    isFavorite: true,
    installed: false,
    license: 'Apache 2.0',
    lastUpdated: '2024-01-12',
    features: [
      '多云存储支持',
      '自动备份',
      '增量备份',
      '版本控制',
      '加密传输',
      '定时任务'
    ],
    reviews: []
  },
  {
    id: 'media-server',
    name: '媒体服务器',
    shortDescription: '流媒体传输和转码',
    fullDescription: '功能完整的媒体服务器，支持视频流传输、实时转码、多设备访问等功能。',
    version: '3.1.0',
    author: 'Media Labs',
    rating: 4.9,
    reviewCount: 456,
    downloadCount: 23456,
    category: 'media',
    iconComponent: ServerIcon,
    color: 'linear-gradient(135deg, #f59e0b 0%, #fbbf24 100%)',
    isFavorite: false,
    installed: false,
    license: 'GPL-3.0',
    lastUpdated: '2024-01-14',
    features: [
      '视频流传输',
      '实时转码',
      '多设备支持',
      '字幕支持',
      '播放列表',
      'DLNA支持'
    ],
    reviews: []
  }
])

const featuredPlugins = computed(() => {
  return plugins.value.filter(p => p.rating >= 4.7).slice(0, 3)
})

const installedPlugins = computed(() => {
  return plugins.value.filter(p => p.installed)
})

const filteredPlugins = computed(() => {
  let result = plugins.value

  // Filter by category
  if (selectedCategory.value !== 'all') {
    result = result.filter(plugin => plugin.category === selectedCategory.value)
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(plugin =>
      plugin.name.toLowerCase().includes(query) ||
      plugin.shortDescription.toLowerCase().includes(query)
    )
  }

  // Sort
  result = [...result].sort((a, b) => {
    switch (sortBy.value) {
      case 'rating':
        return b.rating - a.rating
      case 'popular':
        return b.downloadCount - a.downloadCount
      case 'newest':
        return b.id.localeCompare(a.id)
      case 'updated':
        return new Date(b.lastUpdated).getTime() - new Date(a.lastUpdated).getTime()
      default:
        return 0
    }
  })

  return result
})

const getCategoryIcon = (iconName: string) => {
  const icons: Record<string, any> = {
    FolderIcon,
    ServerIcon
  }
  return icons[iconName] || ServerIcon
}

const getCategoryLabel = (categoryId: string) => {
  const category = categories.find(c => c.id === categoryId)
  return category?.label || categoryId
}

const toggleFavorite = (plugin: any) => {
  plugin.isFavorite = !plugin.isFavorite
}

const toggleInstall = async (plugin: any) => {
  if (plugin.installed) {
    if (confirm(`确定要卸载 "${plugin.name}" 吗?`)) {
      plugin.installed = false
      alert(plugin.name + ' 已成功卸载')
    }
  } else {
    plugin.installed = true
    alert(plugin.name + ' 已成功安装')
  }
}

const openPlugin = (plugin: any) => {
  console.log('Opening plugin:', plugin.id)
}

const showPluginDetails = (plugin: any) => {
  selectedPlugin.value = plugin
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedPlugin.value = null
}

const formatNumber = (num: number) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.plugin-store {
  width: 100%;
  height: 100%;
  background: #f9fafb;
  display: flex;
  flex-direction: column;
}

.store-header {
  padding: 32px;
  border-bottom: 1px solid #e5e7eb;
}

.store-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: #6b7280;
}

.store-layout {
  display: flex;
  height: calc(100% - 120px);
}

.store-sidebar {
  width: 240px;
  padding: 24px;
  border-right: 1px solid #e5e7eb;
  overflow-y: auto;
}

.store-content {
  flex: 1;
  padding: 24px 32px;
  overflow-y: auto;
}

.sidebar-section {
  margin-bottom: 32px;
}

.sidebar-section h3 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
  text-transform: uppercase;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.category-item:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.category-item.active {
  background: #3b82f6;
  color: white;
}

.category-count {
  margin-left: auto;
  font-size: 12px;
  opacity: 0.7;
}

.installed-plugins {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.installed-plugin-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: #f9fafb;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.installed-plugin-item:hover {
  background: #f3f4f6;
}

.plugin-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.plugin-icon.large {
  width: 64px;
  height: 64px;
  border-radius: 16px;
}

.installed-plugin-item .plugin-name {
  flex: 1;
  font-size: 13px;
  color: #1f2937;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.plugin-action {
  padding: 4px;
  background: none;
  border: none;
  color: #6b7280;
  cursor: pointer;
  border-radius: 4px;
}

.plugin-action:hover {
  background: #e5e7eb;
  color: #1f2937;
}

.content-header {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
}

.search-bar {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.search-icon {
  width: 20px;
  height: 20px;
  color: #9ca3af;
}

.search-input {
  flex: 1;
  border: none;
  font-size: 14px;
  color: #1f2937;
  outline: none;
}

.sort-options {
  display: flex;
  align-items: center;
}

.sort-select {
  padding: 12px 16px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  font-size: 14px;
  color: #1f2937;
  background: white;
  cursor: pointer;
}

.featured-section {
  margin-bottom: 32px;
}

.featured-section h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 20px;
}

.featured-plugins {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.featured-plugin-card {
  position: relative;
  background: white;
  border-radius: 16px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.featured-plugin-card:hover {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  transform: translateY(-4px);
}

.featured-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 12px;
  background: #fef3c7;
  color: #92400e;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
}

.featured-plugin-card h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin: 16px 0 8px;
}

.featured-plugin-card p {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 12px;
}

.featured-stats {
  display: flex;
  justify-content: center;
  gap: 16px;
  font-size: 12px;
  color: #6b7280;
}

.featured-stats .rating {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #f59e0b;
}

.plugins-section h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 20px;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px;
  color: #6b7280;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 16px;
}

.plugins-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.plugin-card {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
}

.plugin-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.plugin-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.favorite-btn {
  padding: 4px;
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  transition: color 0.2s ease;
}

.favorite-btn:hover {
  color: #ef4444;
}

.favorite-btn .is-favorite {
  color: #ef4444;
  fill: #ef4444;
}

.plugin-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
  cursor: pointer;
}

.plugin-name:hover {
  color: #3b82f6;
}

.plugin-description {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 8px;
  line-height: 1.5;
}

.plugin-author {
  font-size: 12px;
  color: #9ca3af;
  margin-bottom: 12px;
}

.plugin-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-size: 12px;
}

.plugin-rating {
  display: flex;
  align-items: center;
  gap: 2px;
}

.star-icon {
  width: 14px;
  height: 14px;
  color: #e5e7eb;
}

.star-icon.is-filled {
  color: #fbbf24;
}

.rating-value {
  font-weight: 600;
  color: #1f2937;
  margin-left: 4px;
}

.rating-count {
  color: #9ca3af;
  margin-left: 4px;
}

.plugin-stats {
  color: #9ca3af;
}

.plugin-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 10px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.secondary {
  flex: 0.8;
}

.install-btn {
  flex: 1.2;
  padding: 10px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s ease;
}

.install-btn:hover {
  background: #2563eb;
}

.install-btn.installed {
  background: #10b981;
  cursor: default;
}

/* Modal Styles */
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
  max-width: 90%;
  max-height: 90vh;
  overflow-y: auto;
  animation: modalSlideIn 0.3s ease;
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

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px;
  border-bottom: 1px solid #e5e7eb;
}

.plugin-detail-header {
  display: flex;
  gap: 20px;
  flex: 1;
}

.plugin-detail-info h2 {
  font-size: 24px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.plugin-version {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.plugin-author {
  font-size: 13px;
  color: #9ca3af;
  margin-bottom: 12px;
}

.plugin-detail-meta {
  display: flex;
  gap: 20px;
}

.plugin-detail-meta .rating {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #f59e0b;
}

.plugin-detail-meta .downloads {
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

.plugin-description-section p {
  font-size: 14px;
  color: #6b7280;
  line-height: 1.6;
}

.plugin-features ul {
  list-style: none;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.plugin-features li {
  font-size: 14px;
  color: #6b7280;
  padding-left: 20px;
  position: relative;
}

.plugin-features li::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: #10b981;
  font-weight: bold;
}

.install-steps {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.step {
  display: flex;
  gap: 12px;
}

.step-number {
  width: 24px;
  height: 24px;
  background: #3b82f6;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.step p {
  font-size: 14px;
  color: #6b7280;
  padding-top: 2px;
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

.dev-links {
  display: flex;
  gap: 16px;
}

.dev-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: #f9fafb;
  border-radius: 8px;
  color: #6b7280;
  text-decoration: none;
  transition: all 0.2s ease;
}

.dev-link:hover {
  background: #f3f4f6;
  color: #1f2937;
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

.action-btn.uninstall {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.action-btn.uninstall:hover {
  background: #dc2626;
}
</style>