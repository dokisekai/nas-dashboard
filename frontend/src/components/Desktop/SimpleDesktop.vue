<template>
  <div class="simple-desktop">
    <div class="desktop-background"></div>

    <DesktopTopBar @open-app="openApp" />

    <DesktopWelcome />

    <DesktopWidgets
      :unread-count="notificationStore.unreadCount"
      @open-notifications="notificationStore.openCenter"
    />

    <DesktopQuickAccess :apps="quickAccessApps" @open-app="openApp" />

    <div class="windows-container">
      <DesktopWindow
        v-for="win in windowStore.windows"
        :key="win.id"
        :window="win"
        @focus="windowStore.focusWindow"
        @close="windowStore.closeWindow"
        @minimize="windowStore.minimizeWindow"
        @maximize="windowStore.toggleMaximize"
        @drag-move="windowStore.updatePosition"
        @resize="windowStore.updateSize"
      />
    </div>

    <DesktopDock
      :apps="dockApps"
      :active-app-ids="activeAppIds"
      @open-app="openApp"
      @toggle-window-menu="windowMenuVisible = !windowMenuVisible"
    />

    <WindowManagerMenu
      :visible="windowMenuVisible"
      @tile="handleTile"
      @cascade="handleCascade"
      @minimize-all="handleMinimizeAll"
      @show-all="handleShowAll"
    />

    <NotificationToast />
    <NotificationCenter :is-open="notificationStore.isCenterOpen" @close="notificationStore.closeCenter" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import DesktopTopBar from './DesktopTopBar.vue'
import DesktopWelcome from './DesktopWelcome.vue'
import DesktopWidgets from './DesktopWidgets.vue'
import DesktopQuickAccess from './DesktopQuickAccess.vue'
import DesktopDock from './DesktopDock.vue'
import WindowManagerMenu from './WindowManagerMenu.vue'
import DesktopWindow from './DesktopWindow.vue'
import NotificationToast from '../NotificationSystem/NotificationToast.vue'
import NotificationCenter from '../NotificationSystem/NotificationCenter.vue'
import { useWindowManagerStore } from '../../stores/windowManager'
import { useNotificationStore } from '../../stores/notification'
import { getDockApps, getQuickAccessApps, getApp } from '../../config/apps'

const windowStore = useWindowManagerStore()
const notificationStore = useNotificationStore()

const dockApps = getDockApps()
const quickAccessApps = getQuickAccessApps()

const windowMenuVisible = ref(false)

const activeAppIds = computed(() =>
  windowStore.windows
    .filter(w => !w.minimized && w.focused)
    .map(w => w.appId),
)

// Immich 走 OAuth 授权码流程：浏览器请求 /authorize 时携带同源 session_token cookie，
// SSO AuthorizeHandler 检测到已登录后直接颁发 code 并 302 回 Immich，
// 实现"免登录"打开。
const IMMICH_OAUTH_CONFIG = {
  clientId: 'client-2YbDXsPyx7b8NqmZ',
  redirectUri: 'http://192.168.50.10:2283/auth/login',
}

const openImmichExternal = () => {
  const params = new URLSearchParams({
    client_id: IMMICH_OAUTH_CONFIG.clientId,
    redirect_uri: IMMICH_OAUTH_CONFIG.redirectUri,
    response_type: 'code',
    scope: 'openid email profile',
    state: crypto.randomUUID(),
  })
  window.open(`/authorize?${params.toString()}`, '_blank')
}

const openApp = (appId: string) => {
  const app = getApp(appId)
  if (!app) {
    console.warn(`[Desktop] Unknown app id: ${appId}`)
    return
  }
  if (app.external) {
    openImmichExternal()
    return
  }
  windowStore.openApp(appId)
}

const handleTile = () => {
  windowStore.tileWindows()
  windowMenuVisible.value = false
}

const handleCascade = () => {
  windowStore.cascadeWindows()
  windowMenuVisible.value = false
}

const handleMinimizeAll = () => {
  windowStore.minimizeAll()
  windowMenuVisible.value = false
}

const handleShowAll = () => {
  windowStore.restoreAll()
  windowMenuVisible.value = false
}

// 兼容 QuickShortcutsWidget 通过 CustomEvent 触发的 open-app
const onShortcutOpenApp = (event: Event) => {
  const detail = (event as CustomEvent).detail
  if (typeof detail === 'string') openApp(detail)
}

onMounted(() => {
  notificationStore.fetchNotifications()
  window.addEventListener('open-app', onShortcutOpenApp as EventListener)
})

onUnmounted(() => {
  window.removeEventListener('open-app', onShortcutOpenApp as EventListener)
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

.windows-container {
  position: absolute;
  inset: 0;
  z-index: 10;
  pointer-events: none;
}

.windows-container > :deep(*) {
  pointer-events: auto;
}

@media (prefers-color-scheme: dark) {
  .desktop-background {
    background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  }
}
</style>
