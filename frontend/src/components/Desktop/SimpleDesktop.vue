<template>
  <div class="simple-desktop">
    <!-- 桌面背景 -->
    <div class="desktop-background"></div>

    <!-- 顶部系统栏 -->
    <div class="system-bar">
      <!-- 左侧时间和状态 -->
      <div class="system-bar-left">
        <div class="system-time">{{ currentTime }}</div>
        <div class="system-status" :class="systemStatus.class">
          <div class="status-dot"></div>
          <span>{{ systemStatus.text }}</span>
        </div>
      </div>

      <!-- 右侧操作区 -->
      <div class="system-bar-right">
        <!-- 通知按钮 -->
        <button class="system-bar-btn notification-btn" @click="showNotificationCenter" title="通知中心">
          <BellIcon class="w-4 h-4" />
          <span v-if="notificationUnreadCount > 0" class="notification-badge">
            {{ notificationUnreadCount > 9 ? '9+' : notificationUnreadCount }}
          </span>
        </button>

        <!-- 用户菜单 -->
        <div class="user-menu-container" ref="userMenuRef">
          <button class="system-bar-btn user-btn" @click="toggleUserMenu" title="用户菜单">
            <div class="user-avatar">
              <UserIcon class="w-4 h-4" />
            </div>
            <span class="user-name">{{ username }}</span>
            <ChevronDownIcon class="w-3 h-3 dropdown-arrow" :class="{ rotate: userMenuVisible }" />
          </button>

          <!-- 用户下拉菜单 -->
          <div v-if="userMenuVisible" class="user-dropdown-menu">
            <div class="dropdown-header">
              <div class="user-avatar-large">
                <UserIcon class="w-6 h-6" />
              </div>
              <div class="user-details">
                <div class="dropdown-username">{{ username }}</div>
                <div class="dropdown-role">{{ userRole }}</div>
              </div>
            </div>
            <div class="dropdown-divider"></div>
            <button class="dropdown-item" @click="openWindow('user-profile', '用户资料')">
              <UserIcon class="w-4 h-4" />
              用户资料
            </button>
            <button class="dropdown-item" @click="openWindow('settings', '账户设置')">
              <CogIcon class="w-4 h-4" />
              账户设置
            </button>
            <div class="dropdown-divider"></div>
            <button class="dropdown-item logout" @click="handleLogout">
              <ArrowLeftOnRectangleIcon class="w-4 h-4" />
              退出登录
            </button>
          </div>
        </div>

        <!-- 设置按钮 -->
        <button class="system-bar-btn settings-btn" @click="openWindow('settings', '系统设置')" title="系统设置">
          <CogIcon class="w-4 h-4" />
        </button>

        <!-- 电源按钮 -->
        <div class="power-menu-container">
          <button class="system-bar-btn power-btn" @click="showPowerMenu" title="系统操作">
            <PowerIcon class="w-4 h-4" />
          </button>

          <!-- 电源下拉菜单 -->
          <div v-if="powerMenuVisible" class="power-dropdown-menu">
            <button class="dropdown-item power-option" @click="handleReboot">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              重启系统
            </button>
            <button class="dropdown-item power-option shutdown" @click="handleShutdown">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 9V3.5L22 12l-9 8.5V15c0-5-5.5-9-10-9z" />
              </svg>
              关机
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 欢迎信息 -->
    <div class="welcome-message">
      <h1>NAS Dashboard</h1>
      <p>欢迎使用NAS管理系统</p>
      <div class="status">
        <span class="status-indicator"></span>
        <span>系统运行正常</span>
      </div>
    </div>

    <!-- 桌面小部件 -->
    <div class="desktop-widgets">
      <div class="widget-notification-trigger" @click="showNotificationCenter">
        <BellIcon class="widget-notification-icon" />
        <span v-if="notificationUnreadCount > 0" class="widget-notification-badge">
          {{ notificationUnreadCount }}
        </span>
      </div>

      <ClockWidget class="widget clock-widget" />
      <SystemStatusWidget class="widget system-widget" />
      <WeatherWidget class="widget weather-widget" />
      <CalendarWidget class="widget calendar-widget" />
      <QuickShortcutsWidget class="widget shortcuts-widget" />
    </div>

    <!-- 快速访问图标 -->
    <div class="quick-access">
      <div class="quick-icon" @click="openWindow('storage-manager', '存储管理')">
        <div class="icon-bg">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8" />
          </svg>
        </div>
        <span>存储管理</span>
      </div>

      <div class="quick-icon" @click="openWindow('system-monitor', '系统监控')">
        <div class="icon-bg">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
        <span>系统监控</span>
      </div>

      <div class="quick-icon" @click="openWindow('app-center', '应用中心')">
        <div class="icon-bg">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </div>
        <span>应用中心</span>
      </div>
    </div>

    <!-- 窗口容器 -->
    <div class="windows-container">
      <DesktopWindow
        v-for="window in windows"
        :key="window.id"
        :window="window"
        @focus="focusWindow"
        @close="closeWindow"
        @minimize="minimizeWindow"
        @maximize="maximizeWindow"
        @drag-move="updateWindowPosition"
        @resize="updateWindowSize"
      />
    </div>

    <!-- 底部Dock栏 -->
    <div class="dock">
      <div
        class="dock-item"
        :class="{ active: isWindowActive('storage-manager') }"
        @click="openWindow('storage-manager', '存储管理')"
        title="存储管理"
      >
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8" />
          </svg>
        </div>
      </div>

      <div
        class="dock-item"
        :class="{ active: isWindowActive('system-monitor') }"
        @click="openWindow('system-monitor', '系统监控')"
        title="系统监控"
      >
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
        </div>
      </div>

      <div
        class="dock-item"
        :class="{ active: isWindowActive('app-center') }"
        @click="openWindow('app-center', '应用中心')"
        title="应用中心"
      >
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
        </div>
      </div>

      <div
        class="dock-item"
        :class="{ active: isWindowActive('user-manager') }"
        @click="openWindow('user-manager', '用户管理')"
        title="用户管理"
      >
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0m8 0a4 4 0 11-8 0" />
          </svg>
        </div>
      </div>

      <div
        class="dock-item"
        :class="{ active: isWindowActive('settings') }"
        @click="openWindow('settings', '系统设置')"
        title="设置"
      >
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
      </div>

      <div
        class="dock-item"
        :class="{ active: isWindowActive('control-panel') }"
        @click="openWindow('control-panel', '控制面板')"
        title="控制面板"
      >
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
          </svg>
        </div>
      </div>

      <!-- 窗口管理分隔符 -->
      <div class="dock-separator"></div>

      <!-- 窗口管理按钮 -->
      <div class="dock-item window-manager" @click="showWindowMenu" title="窗口管理">
        <div class="dock-icon">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h14a1 1 0 011 1v2a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM4 13a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM16 13a1 1 0 011-1h2a1 1 0 011 1v6a1 1 0 01-1 1h-2a1 1 0 01-1-1v-6z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- 窗口管理菜单 -->
    <div v-if="windowMenuVisible" class="window-menu">
      <div class="menu-item" @click="tileWindows">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1V5zM14 5a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1V5zM4 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1H5a1 1 0 01-1-1v-6zM14 15a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-6z" />
        </svg>
        平铺窗口
      </div>
      <div class="menu-item" @click="cascadeWindows">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 3l14 9-14 9V3z" />
        </svg>
        层叠窗口
      </div>
      <div class="menu-item" @click="minimizeAllWindows">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
        </svg>
        最小化全部
      </div>
      <div class="menu-item" @click="showAllWindows">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
        </svg>
        显示全部
      </div>
    </div>

    <!-- 通知组件 -->
    <NotificationToast />
    <NotificationCenter />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import DesktopWindow from './DesktopWindow.vue'
import ClockWidget from '../../widgets/ClockWidget.vue'
import SystemStatusWidget from '../../widgets/SystemStatusWidget.vue'
import WeatherWidget from '../../widgets/WeatherWidget.vue'
import CalendarWidget from '../../widgets/CalendarWidget.vue'
import QuickShortcutsWidget from '../../widgets/QuickShortcutsWidget.vue'
import NotificationToast from '../NotificationSystem/NotificationToast.vue'
import NotificationCenter from '../NotificationSystem/NotificationCenter.vue'
import { useNotificationStore } from '../../stores/notification'
import { useAuthStore } from '../../stores/auth'
import { BellIcon, UserIcon, CogIcon, ChevronDownIcon, PowerIcon, ArrowLeftOnRectangleIcon } from '@heroicons/vue/24/outline'

interface Window {
  id: string
  appId: string
  title: string
  position: { x: number; y: number }
  size: { width: number; height: number }
  minimized: boolean
  maximized: boolean
  focused: boolean
  zIndex: number
}

const windows = ref<Window[]>([])
const nextZIndex = ref(1000)
const windowMenuVisible = ref(false)

// 通知系统
const notificationStore = useNotificationStore()
const notificationUnreadCount = computed(() => notificationStore.unreadCount)

// 显示通知中心
const showNotificationCenter = () => {
  notificationStore.showCenter()
}

// 认证和用户信息
const authStore = useAuthStore()
const router = useRouter()
const username = computed(() => authStore.user?.username || 'Admin')
const userRole = computed(() => authStore.user?.role || '管理员')

// 用户菜单状态
const userMenuVisible = ref(false)
const userMenuRef = ref<HTMLElement | null>(null)

// 系统状态
const systemStatus = ref({
  text: '正常',
  class: 'status-normal'
})

// 实时时间
const currentTime = ref('')

// 电源菜单状态
const powerMenuVisible = ref(false)

// 更新时间
const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

// 切换用户菜单
const toggleUserMenu = () => {
  userMenuVisible.value = !userMenuVisible.value
  // 点击外部关闭菜单
  if (userMenuVisible.value) {
    document.addEventListener('click', handleClickOutside)
  }
}

// 处理点击外部
const handleClickOutside = (event: MouseEvent) => {
  if (userMenuRef.value && !userMenuRef.value.contains(event.target as Node)) {
    userMenuVisible.value = false
    document.removeEventListener('click', handleClickOutside)
  }
}

// 显示电源菜单
const showPowerMenu = () => {
  powerMenuVisible.value = !powerMenuVisible.value
}

// 处理退出登录
const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    authStore.clearToken()
    userMenuVisible.value = false
    router.push('/login')
  }
}

// 点击外部关闭电源菜单
const handleClickOutsidePower = (event: MouseEvent) => {
  const powerMenu = document.querySelector('.power-menu')
  if (powerMenu && !powerMenu.contains(event.target as Node)) {
    powerMenuVisible.value = false
    document.removeEventListener('click', handleClickOutsidePower)
  }
}

// 处理重启系统
const handleReboot = () => {
  if (confirm('确定要重启系统吗？所有未保存的工作将丢失！')) {
    // 调用后端API重启系统
    console.log('系统重启中...')
    powerMenuVisible.value = false
    // TODO: 实现重启API调用
  }
}

// 处理关机
const handleShutdown = () => {
  if (confirm('确定要关机吗？')) {
    // 调用后端API关机
    console.log('系统关机中...')
    powerMenuVisible.value = false
    // TODO: 实现关机API调用
  }
}

const appConfigs: Record<string, { title: string; width: number; height: number }> = {
  'storage-manager': { title: '存储管理', width: 900, height: 600 },
  'system-monitor': { title: '系统监控', width: 1000, height: 700 },
  'app-center': { title: '应用中心', width: 900, height: 600 },
  'user-manager': { title: '用户管理', width: 800, height: 500 },
  'settings': { title: '系统设置', width: 700, height: 500 },
  'control-panel': { title: '控制面板', width: 1200, height: 800 }
}

const openWindow = (appId: string, title: string) => {
  const existingWindow = windows.value.find(w => w.appId === appId)

  if (existingWindow) {
    if (existingWindow.minimized) {
      existingWindow.minimized = false
    }
    focusWindow(existingWindow.id)
    return
  }

  const config = appConfigs[appId] || { title: '应用', width: 800, height: 600 }

  const newWindow: Window = {
    id: `window-${Date.now()}`,
    appId,
    title,
    position: {
      x: 100 + windows.value.length * 30,
      y: 100 + windows.value.length * 30
    },
    size: {
      width: config.width,
      height: config.height
    },
    minimized: false,
    maximized: false,
    focused: true,
    zIndex: ++nextZIndex.value
  }

  windows.value.forEach(w => w.focused = false)
  windows.value.push(newWindow)
}

const closeWindow = (windowId: string) => {
  windows.value = windows.value.filter(w => w.id !== windowId)
}

const focusWindow = (windowId: string) => {
  console.log('focusWindow called with:', windowId)
  const window = windows.value.find(w => w.id === windowId)
  if (window) {
    // 先将所有窗口设置为非焦点状态
    windows.value.forEach(w => w.focused = false)

    // 设置当前窗口为焦点状态
    window.focused = true
    window.zIndex = ++nextZIndex.value

    // 如果窗口被最小化，恢复它
    if (window.minimized) {
      window.minimized = false
    }

    console.log('Window focused:', windowId, 'new zIndex:', window.zIndex)
  } else {
    console.log('Window not found:', windowId)
  }
}

const minimizeWindow = (windowId: string) => {
  console.log('minimizeWindow called with:', windowId)
  const window = windows.value.find(w => w.id === windowId)
  if (window) {
    window.minimized = true
    window.focused = false
    console.log('Window minimized:', windowId)
  }
}

const maximizeWindow = (windowId: string) => {
  console.log('maximizeWindow called with:', windowId)
  const window = windows.value.find(w => w.id === windowId)
  if (window) {
    const previousState = window.maximized
    window.maximized = !window.maximized
    console.log('Window maximized toggled:', windowId, 'from', previousState, 'to', window.maximized)
    if (window.maximized) {
      window.position = { x: 0, y: 0 }
      console.log('Window maximized to full screen')
    } else {
      window.position = {
        x: 100 + windows.value.findIndex(w => w.id === windowId) * 30,
        y: 100 + windows.value.findIndex(w => w.id === windowId) * 30
      }
      console.log('Window restored to normal size at:', window.position)
    }
  }
}

const updateWindowPosition = (windowId: string, position: { x: number; y: number }) => {
  const window = windows.value.find(w => w.id === windowId)
  if (window) {
    window.position = position
    console.log('Window position updated:', windowId, 'to:', position)
  }
}

const updateWindowSize = (windowId: string, size: { width: number; height: number }) => {
  const window = windows.value.find(w => w.id === windowId)
  if (window) {
    window.size = size
  }
}

const isWindowActive = (appId: string) => {
  return windows.value.some(w => w.appId === appId && !w.minimized && w.focused)
}

// 窗口平铺功能
const tileWindows = () => {
  const activeWindows = windows.value.filter(w => !w.minimized)
  if (activeWindows.length === 0) return

  const dockHeight = 80
  const availableWidth = window.innerWidth
  const availableHeight = window.innerHeight - dockHeight

  if (activeWindows.length === 1) {
    // 单个窗口居中
    activeWindows[0].position = {
      x: (availableWidth - activeWindows[0].size.width) / 2,
      y: (availableHeight - activeWindows[0].size.height) / 2
    }
  } else if (activeWindows.length === 2) {
    // 两个窗口左右平铺
    const halfWidth = availableWidth / 2
    activeWindows.forEach((window, index) => {
      window.position = {
        x: index * halfWidth,
        y: 0
      }
      window.size = {
        width: halfWidth,
        height: availableHeight
      }
    })
  } else {
    // 多个窗口网格平铺
    const cols = Math.ceil(Math.sqrt(activeWindows.length))
    const rows = Math.ceil(activeWindows.length / cols)
    const tileWidth = availableWidth / cols
    const tileHeight = availableHeight / rows

    activeWindows.forEach((window, index) => {
      const col = index % cols
      const row = Math.floor(index / cols)
      window.position = {
        x: col * tileWidth,
        y: row * tileHeight
      }
      window.size = {
        width: tileWidth,
        height: tileHeight
      }
    })
  }
}

// 窗口层叠功能
const cascadeWindows = () => {
  const activeWindows = windows.value.filter(w => !w.minimized)
  if (activeWindows.length === 0) return

  const baseX = 50
  const baseY = 50
  const offsetX = 30
  const offsetY = 30

  activeWindows.forEach((win, index) => {
    win.position = {
      x: baseX + index * offsetX,
      y: baseY + index * offsetY
    }
    win.size = {
      width: Math.min(800, window.innerWidth - baseX - index * offsetX - 100),
      height: Math.min(600, window.innerHeight - baseY - index * offsetY - 100)
    }
    win.maximized = false
  })
}

// 最小化所有窗口
const minimizeAllWindows = () => {
  windows.value.forEach(window => {
    window.minimized = true
  })
}

// 显示所有窗口
const showAllWindows = () => {
  windows.value.forEach(window => {
    window.minimized = false
  })
  windowMenuVisible.value = false
}

// 显示窗口管理菜单
const showWindowMenu = () => {
  windowMenuVisible.value = !windowMenuVisible.value
  // 点击外部关闭菜单
  if (windowMenuVisible.value) {
    document.addEventListener('click', handleClickOutsideWindow)
  }
}

// 点击其他地方关闭窗口菜单
const handleClickOutsideWindow = (event: MouseEvent) => {
  const target = event.target as HTMLElement
  // 只处理窗口菜单相关的点击，不影响窗口交互
  if (!target.closest('.window-menu') && !target.closest('.window-manager')) {
    windowMenuVisible.value = false
    document.removeEventListener('click', handleClickOutsideWindow)
  }
}

// 生命周期钩子
let timeUpdateInterval: number

onMounted(() => {
  // 初始化时间
  updateTime()
  // 每秒更新时间
  timeUpdateInterval = window.setInterval(updateTime, 1000)
})

onUnmounted(() => {
  // 清理时间更新定时器
  if (timeUpdateInterval) {
    clearInterval(timeUpdateInterval)
  }
})
</script>

<style scoped>
.simple-desktop {
  width: 100vw;
  height: 100vh;
  position: relative;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.desktop-background {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  background-size: cover;
  background-position: center;
}

/* 系统栏样式 */
.system-bar {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 48px;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  z-index: 100;
  border-bottom: 1px solid rgba(255, 255, 255, 0.2);
}

.system-bar-left {
  display: flex;
  align-items: center;
  gap: 16px;
  color: white;
  font-size: 14px;
}

.system-time {
  font-weight: 500;
  font-variant-numeric: tabular-nums;
}

.system-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  opacity: 0.9;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.system-status.status-normal .status-dot {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.6);
}

.system-status.status-warning .status-dot {
  background: #f59e0b;
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.6);
}

.system-status.status-error .status-dot {
  background: #ef4444;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.6);
}

.system-bar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.system-bar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;
}

.system-bar-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: translateY(-1px);
}

.system-bar-btn:active {
  transform: translateY(0);
}

.notification-btn {
  position: relative;
}

.notification-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: 600;
  min-width: 16px;
  height: 16px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 4px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.user-menu-container {
  position: relative;
}

.user-btn {
  width: auto;
  padding: 0 12px;
  gap: 8px;
}

.user-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.user-name {
  font-size: 13px;
  font-weight: 500;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dropdown-arrow {
  transition: transform 0.2s;
}

.dropdown-arrow.rotate {
  transform: rotate(180deg);
}

/* 用户下拉菜单 */
.user-dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 200px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.1);
  margin-bottom: 8px;
}

.user-avatar-large {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.user-details {
  flex: 1;
}

.dropdown-username {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.dropdown-role {
  font-size: 12px;
  color: #6b7280;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  border: none;
  background: transparent;
  border-radius: 8px;
  color: #374151;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: left;
}

.dropdown-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.dropdown-item.logout {
  color: #ef4444;
}

.dropdown-item.logout:hover {
  background: rgba(239, 68, 68, 0.1);
}

.dropdown-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.1);
  margin: 8px 0;
}

/* 电源菜单样式 */
.power-menu-container {
  position: relative;
}

.power-dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 160px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 8px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  animation: slideDown 0.2s ease-out;
}

.power-option {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  color: #374151;
  font-size: 13px;
  font-weight: 500;
}

.power-option.shutdown {
  color: #ef4444;
}

.power-option.shutdown:hover {
  background: rgba(239, 68, 68, 0.1);
  color: #dc2626;
}

.power-option:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.settings-btn:hover {
  background: rgba(102, 126, 234, 0.2);
}

.power-btn:hover {
  background: rgba(239, 68, 68, 0.2);
}

.welcome-message {
  position: absolute;
  top: calc(15% + 48px);
  left: 50%;
  transform: translateX(-50%);
  text-align: center;
  color: white;
  z-index: 1;
  pointer-events: none;
  margin-top: 0;
}

.welcome-message h1 {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 16px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
  letter-spacing: 1px;
}

.welcome-message p {
  font-size: 24px;
  opacity: 0.9;
  margin-bottom: 24px;
  font-weight: 300;
}

.status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 500;
}

.status-indicator {
  width: 12px;
  height: 12px;
  background: #10b981;
  border-radius: 50%;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.desktop-widgets {
  position: absolute;
  top: 64px;
  right: 20px;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  z-index: 5;
  max-width: 450px;
  width: 100%;
  padding: 0 20px;
  box-sizing: border-box;
}

.widget {
  animation: widgetSlideIn 0.5s ease-out;
}

@keyframes widgetSlideIn {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.clock-widget {
  grid-column: span 1;
}

.system-widget {
  grid-column: span 1;
}

.weather-widget {
  grid-column: span 1;
}

.calendar-widget {
  grid-column: span 2;
}

.shortcuts-widget {
  grid-column: span 2;
}

/* 通知触发器样式 */
.widget-notification-trigger {
  position: relative;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  min-height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.widget-notification-trigger:hover {
  background: rgba(255, 255, 255, 0.25);
  transform: translateY(-2px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
}

.widget-notification-icon {
  color: white;
  width: 24px;
  height: 24px;
}

.widget-notification-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  background: rgba(239, 68, 68, 0.9);
  backdrop-filter: blur(10px);
  color: white;
  font-size: 11px;
  font-weight: 600;
  min-width: 18px;
  height: 18px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 5px;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

/* 响应式调整 - 小屏幕单列显示 */
@media (max-width: 768px) {
  .desktop-widgets {
    top: 64px;
    right: 10px;
    left: 10px;
    max-width: none;
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .clock-widget,
  .system-widget,
  .weather-widget,
  .calendar-widget,
  .shortcuts-widget {
    grid-column: span 1;
  }
}

.quick-access {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  gap: 32px;
  z-index: 1;
}

.quick-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.quick-icon:hover {
  transform: scale(1.1);
}

.icon-bg {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.quick-icon span {
  color: white;
  font-size: 14px;
  font-weight: 500;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.windows-container {
  position: absolute;
  inset: 0;
  z-index: 10;
  pointer-events: none;
}

.windows-container > * {
  pointer-events: auto;
}

.dock {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  z-index: 100;
  pointer-events: all;
}

.dock-item {
  width: 56px;
  height: 56px;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.dock-item:hover {
  transform: translateY(-8px) scale(1.1);
}

.dock-item.active::after {
  content: '';
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background: white;
  border-radius: 50%;
  box-shadow: 0 0 8px rgba(255, 255, 255, 0.8);
}

.dock-icon {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 深色模式支持 */
@media (prefers-color-scheme: dark) {
  .desktop-background {
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  }

  .dock {
    background: rgba(0, 0, 0, 0.3);
    border-color: rgba(255, 255, 255, 0.1);
  }
}

.dock-separator {
  width: 1px;
  height: 40px;
  background: rgba(255, 255, 255, 0.3);
  margin: 0 8px;
}

.window-manager {
  cursor: pointer;
}

.window-menu {
  position: absolute;
  bottom: 100px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.5);
  padding: 8px 0;
  min-width: 180px;
  z-index: 200;
  animation: menuSlideUp 0.2s ease-out;
}

@keyframes menuSlideUp {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0);
  }
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s ease;
}

.menu-item:hover {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.menu-item:first-child {
  border-radius: 12px 12px 0 0;
}

.menu-item:last-child {
  border-radius: 0 0 12px 12px;
}
</style>
