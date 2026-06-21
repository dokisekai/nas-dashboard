<template>
  <div class="ipados-home">
    <!-- 主屏幕区域 -->
    <div class="home-screen" :class="{ 'app-library-mode': showAppLibrary }">

      <!-- 状态栏 -->
      <div class="status-bar">
        <div class="status-bar-left">
          <span class="time">{{ currentTime }}</span>
        </div>
        <div class="status-bar-center">
          <div class="notch" v-if="hasNotch"></div>
        </div>
	        <div class="status-bar-right">
	          <SignalIcon class="icon" />
	          <WifiIcon class="icon" />
	          <PowerIcon class="icon" />
	        </div>
      </div>

      <!-- 主屏幕内容 -->
      <div class="screen-content" v-if="!showAppLibrary">
        <!-- 小组件区域 -->
        <div class="widgets-area">
          <iPadOSWidget
            v-for="widget in visibleWidgets"
            :key="widget.id"
            :widget="widget"
            :size="widget.size as 'small' | 'medium' | 'large' | 'extra-large'"
          />
        </div>

        <!-- 应用列表墙 -->
        <div class="app-wall">
          <div class="wall-header">
            <h2 class="wall-title">应用 ({{ allApps.length }})</h2>
            <div class="wall-controls">
              <button class="wall-btn" @click="toggleEditMode">
                <PencilIcon class="btn-icon" />
                <span>{{ editMode ? '完成' : '编辑' }}</span>
              </button>
              <button class="wall-btn" @click="toggleWallView">
                <Squares2X2Icon class="btn-icon" />
                <span>{{ wallViewMode }}</span>
              </button>
            </div>
          </div>

          <!-- 调试信息 -->
          <div v-if="allApps.length === 0" class="debug-info">
            没有找到应用数据
          </div>

          <!-- 墙式网格 -->
          <div class="wall-grid" :class="`wall-${wallViewMode}`">
            <div
              v-for="app in allApps"
              :key="app.id"
              class="wall-app"
              :class="{
                'editing': editMode,
                'selected': selectedApps.includes(app.id)
              }"
              @click="openApp(app)"
              @contextmenu="showAppContextMenu(app, $event)"
            >
              <div class="wall-app-icon" :class="{ 'squircle': app.squircle }">
                <component :is="app.icon" class="icon-svg" />
                <div v-if="app.badge" class="app-badge">{{ app.badge }}</div>
                <div v-if="editMode" class="edit-checkbox">
                  <CheckIcon v-if="selectedApps.includes(app.id)" class="check-icon" />
                </div>
              </div>
              <div class="wall-app-name">{{ app.name }}</div>
              <div v-if="app.category" class="wall-app-category">{{ app.category }}</div>
            </div>
          </div>
        </div>

        <!-- 页面指示器 -->
        <div class="page-dots">
          <div
            v-for="(_, index) in pages"
            :key="index"
            class="dot"
            :class="{ active: index === currentPageIndex }"
            @click="goToPage(index)"
          />
        </div>
      </div>

      <!-- 应用库 -->
      <div class="app-library" v-if="showAppLibrary">
        <div class="library-header">
          <h2>应用库</h2>
          <SearchIcon class="search-icon" @click="openSpotlight" />
        </div>

        <!-- 应用搜索 -->
        <div class="app-search">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索应用"
            class="search-input"
          >
        </div>

        <!-- 分类应用 -->
        <div class="app-categories">
          <div
            v-for="category in filteredCategories"
            :key="category.id"
            class="category-section"
          >
            <h3 class="category-title">{{ category.name }}</h3>
            <div class="category-apps">
              <div
                v-for="app in category.apps"
                :key="app.id"
                class="library-app"
                @click="openApp(app)"
              >
                <div class="library-app-icon">
                  <component :is="app.icon" class="icon-svg" />
                </div>
                <div class="library-app-name">{{ app.name }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 程序坞 -->
      <div class="dock">
        <div class="dock-background">
          <div class="dock-separator left" />
          <div class="dock-apps">
            <div
              v-for="app in dockApps"
              :key="app.id"
              class="dock-app"
              :class="{ active: isAppActive(app.id) }"
              @click="openApp(app)"
              @contextmenu="showDockContextMenu(app, $event)"
            >
              <div class="dock-app-icon">
                <component :is="app.icon" class="icon-svg" />
                <div v-if="app.badge" class="dock-app-badge">{{ app.badge }}</div>
              </div>
            </div>
          </div>
          <div class="dock-separator right" />
        </div>
      </div>

      <!-- 底部主屏幕指示器 -->
      <div class="home-indicator" @click="goToHomeScreen"></div>
    </div>

    <!-- 控制中心 -->
    <Transition name="control-panel">
      <iOSControlCenter
        v-if="showControlCenter"
        @close="showControlCenter = false"
      />
    </Transition>

    <!-- Spotlight搜索 -->
    <Transition name="spotlight">
      <iOSSpotlight
        v-if="showSpotlight"
        @close="showSpotlight = false"
        @openApp="openApp"
      />
    </Transition>

    <!-- 应用上下文菜单 -->
    <Transition name="context-menu">
      <div
        v-if="contextMenu.visible"
        class="context-menu"
        :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }"
      >
        <div class="context-menu-items">
          <div class="context-menu-item" @click="editApp">
            <PencilIcon class="menu-icon" />
            <span>编辑主屏幕</span>
          </div>
          <div class="context-menu-item" @click="removeFromDock">
            <TrashIcon class="menu-icon" />
            <span>从程序坞移除</span>
          </div>
          <div class="context-menu-item" @click="shareApp">
            <ShareIcon class="menu-icon" />
            <span>分享应用</span>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  SignalIcon,
  WifiIcon,
  PowerIcon,
  MagnifyingGlassIcon,
  PencilIcon,
  TrashIcon,
  ShareIcon,
  CubeIcon,
  ServerIcon,
  ChartBarIcon,
  UserGroupIcon,
  CloudArrowUpIcon,
  FolderIcon,
  CogIcon,
  BellIcon,
  BeakerIcon,
  PhotoIcon,
  MusicalNoteIcon,
  Squares2X2Icon,
  CheckIcon,
  ListBulletIcon,
  RectangleStackIcon,
  CpuChipIcon
} from '@heroicons/vue/24/outline'
import iPadOSWidget from './iPadOSWidget.vue'
import iOSControlCenter from './iOSControlCenter.vue'
import iOSSpotlight from './iOSSpotlight.vue'

// 状态
const currentTime = ref('')
const hasNotch = ref(true)
const showAppLibrary = ref(false)
const showControlCenter = ref(false)
const showSpotlight = ref(false)
const currentPageIndex = ref(0)
const searchQuery = ref('')
const wallViewMode = ref('grid') // 'grid' | 'list' | 'tiles'
const editMode = ref(false)
const selectedApps = ref<string[]>([])

// 上下文菜单
const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  app: null as any
})

// 应用数据
const apps = ref([
  { id: 'storage', name: '存储', icon: ServerIcon, squircle: true, category: '工具' },
  { id: 'monitor', name: '监控', icon: ChartBarIcon, squircle: true, category: '系统' },
  { id: 'users', name: '用户', icon: UserGroupIcon, squircle: true, category: '系统' },
  { id: 'backup', name: '备份', icon: CloudArrowUpIcon, squircle: true, category: '工具' },
  { id: 'files', name: '文件', icon: FolderIcon, squircle: true, category: '生产力' },
  { id: 'settings', name: '设置', icon: CogIcon, squircle: true, category: '系统' },
  { id: 'notifications', name: '通知', icon: BellIcon, badge: 3, squircle: true, category: '系统' },
  { id: 'sso', name: 'SSO', icon: BeakerIcon, squircle: true, category: '系统' },
  { id: 'photos', name: '照片', icon: PhotoIcon, squircle: true, category: '媒体' },
  { id: 'music', name: '音乐', icon: MusicalNoteIcon, squircle: true, category: '媒体' },
  { id: 'docker', name: 'Docker', icon: CubeIcon, squircle: true, category: '开发' },
  { id: 'app-center', name: '应用中心', icon: Squares2X2Icon, squircle: true, category: '系统' },
  { id: 'network', name: '网络', icon: WifiIcon, squircle: true, category: '系统' },
  { id: 'power', name: '电源', icon: PowerIcon, squircle: true, category: '系统' },
  { id: 'calendar', name: '日历', icon: RectangleStackIcon, squircle: true, category: '生产力' },
  { id: 'notes', name: '备忘录', icon: ListBulletIcon, squircle: true, category: '生产力' },
  { id: 'system', name: '系统', icon: CpuChipIcon, squircle: true, category: '系统' },
])

// 所有应用（用于应用墙）
const allApps = computed(() => {
  return apps.value
})

// 程序坞应用
const dockApps = ref([
  { id: 'phone', name: '电话', icon: BellIcon },
  { id: 'safari', name: 'Safari', icon: CubeIcon },
  { id: 'messages', name: '信息', icon: BellIcon, badge: 2 },
  { id: 'mail', name: '邮件', icon: BellIcon, badge: 15 },
  { id: 'music', name: '音乐', icon: MusicalNoteIcon },
  { id: 'photos', name: '照片', icon: PhotoIcon },
  { id: 'files', name: '文件', icon: FolderIcon },
  { id: 'settings', name: '设置', icon: CogIcon },
])

// 小组件
const visibleWidgets = ref([
  { id: 'clock', size: 'small', type: 'clock' },
  { id: 'weather', size: 'medium', type: 'weather' },
  { id: 'calendar', size: 'medium', type: 'calendar' },
  { id: 'battery', size: 'small', type: 'battery' },
  { id: 'storage', size: 'large', type: 'storage' },
])

// 分页数据
const pages = computed(() => {
  const allApps = [...apps.value]
  const pageSize = 16 // 每页16个应用
  const pageCount = Math.ceil(allApps.length / pageSize)

  return Array.from({ length: pageCount }, (_, i) => ({
    id: i,
    apps: allApps.slice(i * pageSize, (i + 1) * pageSize)
  }))
})

const currentPageApps = computed(() => {
  return pages.value[currentPageIndex.value]?.apps || []
})

// 应用分类
const categories = computed(() => {
  return [
    {
      id: 'productivity',
      name: '生产力',
      apps: apps.value.filter(app =>
        ['storage', 'files', 'backup'].includes(app.id)
      )
    },
    {
      id: 'utilities',
      name: '工具',
      apps: apps.value.filter(app =>
        ['monitor', 'settings', 'notifications'].includes(app.id)
      )
    },
    {
      id: 'social',
      name: '社交',
      apps: apps.value.filter(app =>
        ['users', 'sso'].includes(app.id)
      )
    },
    {
      id: 'media',
      name: '媒体',
      apps: apps.value.filter(app =>
        ['photos', 'music'].includes(app.id)
      )
    }
  ]
})

const filteredCategories = computed(() => {
  if (!searchQuery.value) return categories.value

  return categories.value.map(category => ({
    ...category,
    apps: category.apps.filter(app =>
      app.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  })).filter(category => category.apps.length > 0)
})

// 方法
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })
}

const goToPage = (index: number) => {
  currentPageIndex.value = index
}

const openApp = (app: any) => {
  contextMenu.value.visible = false
  console.log('Opening app:', app.id)
  // 这里集成现有的应用打开逻辑
}

const isAppActive = (appId: string) => {
  // 检查应用是否处于活跃状态
  return false
}

const openSpotlight = () => {
  showSpotlight.value = true
}

const showAppContextMenu = (app: any, event: MouseEvent) => {
  event.preventDefault()
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    app
  }
}

const showDockContextMenu = (app: any, event: MouseEvent) => {
  event.preventDefault()
  showAppContextMenu(app, event)
}

const editApp = () => {
  console.log('Edit app:', contextMenu.value.app)
  contextMenu.value.visible = false
}

const removeFromDock = () => {
  console.log('Remove from dock:', contextMenu.value.app)
  const index = dockApps.value.findIndex(app => app.id === contextMenu.value.app?.id)
  if (index > -1) {
    dockApps.value.splice(index, 1)
  }
  contextMenu.value.visible = false
}

const shareApp = () => {
  console.log('Share app:', contextMenu.value.app)
  contextMenu.value.visible = false
}

const goToHomeScreen = () => {
  showAppLibrary.value = false
}

// 应用列表墙方法
const toggleEditMode = () => {
  editMode.value = !editMode.value
  if (!editMode.value) {
    selectedApps.value = []
  }
}

const toggleWallView = () => {
  const modes = ['grid', 'list', 'tiles']
  const currentIndex = modes.indexOf(wallViewMode.value)
  wallViewMode.value = modes[(currentIndex + 1) % modes.length]
}

const toggleAppSelection = (appId: string) => {
  const index = selectedApps.value.indexOf(appId)
  if (index > -1) {
    selectedApps.value.splice(index, 1)
  } else {
    selectedApps.value.push(appId)
  }
}

// 生命周期
let timeInterval: any

onMounted(() => {
  updateTime()
  timeInterval = setInterval(updateTime, 1000)

  // 滑动手势支持
  let touchStartX = 0
  let touchEndX = 0

  const handleTouchStart = (e: TouchEvent) => {
    touchStartX = e.changedTouches[0].screenX
  }

  const handleTouchEnd = (e: TouchEvent) => {
    touchEndX = e.changedTouches[0].screenX
    handleSwipe()
  }

  const handleSwipe = () => {
    const swipeThreshold = 50
    const diff = touchStartX - touchEndX

    if (Math.abs(diff) > swipeThreshold) {
      if (diff > 0 && currentPageIndex.value < pages.value.length - 1) {
        // 向左滑动，下一页
        currentPageIndex.value++
      } else if (diff < 0 && currentPageIndex.value > 0) {
        // 向右滑动，上一页
        currentPageIndex.value--
      }
    } else if (Math.abs(diff) < 10) {
      // 点击，切换主屏幕/应用库
      showAppLibrary.value = !showAppLibrary.value
    }
  }

  document.addEventListener('touchstart', handleTouchStart)
  document.addEventListener('touchend', handleTouchEnd)
})

onUnmounted(() => {
  if (timeInterval) {
    clearInterval(timeInterval)
  }
})
</script>

<style scoped>
.ipados-home {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
  position: relative;
}

.status-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 44px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  color: white;
  z-index: 1000;
  backdrop-filter: blur(20px);
  background: rgba(0, 0, 0, 0.2);
}

.status-bar-left,
.status-bar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.time {
  font-size: 16px;
  font-weight: 600;
}

.icon {
  width: 20px;
  height: 20px;
}

.notch {
  width: 200px;
  height: 30px;
  background: black;
  border-radius: 0 0 20px 20px;
}

.home-screen {
  height: 100%;
  padding-top: 44px;
  transition: all 0.3s ease;
}

.screen-content {
  height: calc(100% - 120px);
  padding: 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.widgets-area {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 40px;
}

.app-grid {
  margin-bottom: 20px;
}

.apps-container {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  padding: 20px;
}

.app-icon-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.app-icon {
  width: 64px;
  height: 64px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(255,255,255,0.8), rgba(255,255,255,0.6));
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
  position: relative;
  transition: transform 0.2s ease;
}

.app-icon.squircle {
  border-radius: 16px;
}

.app-icon:active {
  transform: scale(0.95);
}

.icon-svg {
  width: 32px;
  height: 32px;
  color: #333;
}

.app-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ff3b30;
  color: white;
  border-radius: 10px;
  min-width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 600;
  padding: 0 6px;
}

.app-name {
  color: white;
  font-size: 12px;
  font-weight: 500;
  text-align: center;
  text-shadow: 0 2px 8px rgba(0,0,0,0.5);
}

.page-dots {
  display: flex;
  justify-content: center;
  gap: 8px;
  margin: 20px 0;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255,255,255,0.4);
  cursor: pointer;
  transition: all 0.2s ease;
}

.dot.active {
  background: white;
  width: 20px;
  border-radius: 4px;
}

/* 应用列表墙样式 */
.app-wall {
  padding: 20px;
  overflow-y: auto;
  min-height: 400px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  margin-top: 20px;
}

.wall-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0 8px;
}

.wall-title {
  font-size: 28px;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

.wall-controls {
  display: flex;
  gap: 12px;
}

.wall-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(20px);
  border: none;
  border-radius: 20px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.wall-btn:active {
  transform: scale(0.95);
  background: rgba(255,255,255,0.3);
}

.wall-btn.active {
  background: rgba(59, 130, 246, 0.6);
}

.btn-icon {
  width: 16px;
  height: 16px;
}

.debug-info {
  padding: 20px;
  color: #ff6b6b;
  font-size: 16px;
  text-align: center;
  background: rgba(255, 0, 0, 0.1);
  border-radius: 8px;
  margin: 20px 0;
}

.wall-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 0 8px;
}

.wall-title {
  font-size: 28px;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

.wall-controls {
  display: flex;
  gap: 12px;
}

.wall-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(20px);
  border: none;
  border-radius: 20px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.wall-btn:active {
  transform: scale(0.95);
  background: rgba(255,255,255,0.3);
}

.wall-btn.active {
  background: rgba(59, 130, 246, 0.6);
}

.btn-icon {
  width: 16px;
  height: 16px;
}

.wall-grid {
  display: grid;
  gap: 20px;
  transition: all 0.3s ease;
}

/* 网格视图 */
.wall-grid.wall-grid {
  grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
  padding: 0 8px;
}

/* 列表视图 */
.wall-grid.wall-list {
  grid-template-columns: 1fr;
  gap: 8px;
}

/* 瓷砖视图 */
.wall-grid.wall-tiles {
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.wall-app {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 12px 8px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.wall-app:active {
  transform: scale(0.95);
}

.wall-app.editing {
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.wall-app.selected .wall-app-icon {
  border: 2px solid #3b82f6;
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.5);
}

.wall-app-icon {
  width: 64px;
  height: 64px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(255,255,255,0.8), rgba(255,255,255,0.6));
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
  position: relative;
  transition: all 0.2s ease;
}

.wall-app-icon.squircle {
  border-radius: 16px;
}

.wall-app-name {
  font-size: 11px;
  color: white;
  text-align: center;
  text-shadow: 0 2px 8px rgba(0,0,0,0.5);
  font-weight: 500;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wall-app-category {
  font-size: 9px;
  color: rgba(255,255,255,0.7);
  text-align: center;
  text-shadow: 0 1px 4px rgba(0,0,0,0.3);
}

.edit-checkbox {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 24px;
  height: 24px;
  background: #3b82f6;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.3);
}

.check-icon {
  width: 14px;
  height: 14px;
  color: white;
}

/* 列表视图特殊样式 */
.wall-list .wall-app {
  flex-direction: row;
  padding: 12px 16px;
  background: rgba(255,255,255,0.1);
  backdrop-filter: blur(20px);
}

.wall-list .wall-app:hover {
  background: rgba(255,255,255,0.15);
}

.wall-list .wall-app-icon {
  width: 48px;
  height: 48px;
  margin-right: 12px;
}

.wall-list .wall-app-name {
  flex: 1;
  text-align: left;
  font-size: 14px;
  max-width: none;
}

.wall-list .wall-app-category {
  font-size: 12px;
  text-align: right;
}

/* 瓷砖视图特殊样式 */
.wall-tiles .wall-app {
  padding: 16px;
  background: rgba(255,255,255,0.15);
  backdrop-filter: blur(30px);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.2);
}

.wall-tiles .wall-app:hover {
  transform: translateY(-4px);
  box-shadow: 0 16px 48px rgba(0,0,0,0.3);
}

.wall-tiles .wall-app-icon {
  width: 72px;
  height: 72px;
  margin-bottom: 8px;
}

.wall-tiles .wall-app-name {
  font-size: 13px;
  font-weight: 600;
  max-width: 100px;
}

.wall-tiles .wall-app-category {
  font-size: 10px;
  color: rgba(255,255,255,0.8);
}

/* 应用库样式 */
.app-library {
  height: 100%;
  padding: 20px;
  overflow-y: auto;
}

.library-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  color: white;
}

.library-header h2 {
  font-size: 32px;
  font-weight: 700;
}

.search-icon {
  width: 24px;
  height: 24px;
  cursor: pointer;
}

.app-search {
  margin-bottom: 30px;
}

.search-input {
  width: 100%;
  height: 44px;
  border-radius: 12px;
  border: none;
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(20px);
  color: white;
  font-size: 16px;
  padding: 0 16px;
}

.search-input::placeholder {
  color: rgba(255,255,255,0.6);
}

.app-categories {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.category-section {
  color: white;
}

.category-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
  color: rgba(255,255,255,0.9);
}

.category-apps {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 16px;
}

.library-app {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.library-app-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(255,255,255,0.8), rgba(255,255,255,0.6));
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(0,0,0,0.2);
}

.library-app-name {
  font-size: 11px;
  color: white;
  text-align: center;
}

/* 程序坞样式 */
.dock {
  position: fixed;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 100;
}

.dock-background {
  background: rgba(255,255,255,0.2);
  backdrop-filter: blur(30px);
  border-radius: 24px;
  padding: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}

.dock-separator {
  width: 1px;
  height: 40px;
  background: rgba(255,255,255,0.3);
}

.dock-apps {
  display: flex;
  gap: 12px;
}

.dock-app {
  cursor: pointer;
  transition: all 0.2s ease;
}

.dock-app:active {
  transform: scale(0.95);
}

.dock-app-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(255,255,255,0.8), rgba(255,255,255,0.6));
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(0,0,0,0.2);
  position: relative;
}

.dock-app.active .dock-app-icon {
  background: linear-gradient(135deg, rgba(255,255,255,0.9), rgba(255,255,255,0.7));
}

.dock-app-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ff3b30;
  color: white;
  border-radius: 10px;
  min-width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 600;
  padding: 0 4px;
}

.home-indicator {
  position: fixed;
  bottom: 8px;
  left: 50%;
  transform: translateX(-50%);
  width: 134px;
  height: 5px;
  background: rgba(255,255,255,0.8);
  border-radius: 3px;
  cursor: pointer;
  z-index: 1000;
}

/* 上下文菜单 */
.context-menu {
  position: fixed;
  background: rgba(255,255,255,0.9);
  backdrop-filter: blur(30px);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
  min-width: 200px;
  z-index: 2000;
  overflow: hidden;
}

.context-menu-items {
  padding: 8px 0;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.2s ease;
  color: #333;
}

.context-menu-item:hover {
  background: rgba(0,0,0,0.1);
}

.menu-icon {
  width: 20px;
  height: 20px;
  color: #666;
}

/* 动画 */
.control-panel-enter-active,
.control-panel-leave-active {
  transition: all 0.3s ease;
}

.control-panel-enter-from,
.control-panel-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.spotlight-enter-active,
.spotlight-leave-active {
  transition: all 0.3s ease;
}

.spotlight-enter-from,
.spotlight-leave-to {
  opacity: 0;
}

.context-menu-enter-active,
.context-menu-leave-active {
  transition: all 0.2s ease;
}

.context-menu-enter-from,
.context-menu-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .apps-container {
    grid-template-columns: repeat(3, 1fr);
  }

  .category-apps {
    grid-template-columns: repeat(4, 1fr);
  }

  .dock-apps {
    gap: 8px;
  }

  .dock-app-icon {
    width: 48px;
    height: 48px;
  }
}
</style>