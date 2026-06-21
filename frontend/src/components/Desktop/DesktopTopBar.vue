<template>
  <div class="system-bar">
    <!-- 左侧：时间和系统状态 -->
    <div class="system-bar-left">
      <div class="system-time">{{ currentTime }}</div>
      <div class="system-status" :class="statusClass">
        <div class="status-dot"></div>
        <span>{{ statusText }}</span>
      </div>
    </div>

    <!-- 右侧：通知 / 用户 / 同步 / 电源 -->
    <div class="system-bar-right">
      <button class="system-bar-btn notification-btn" title="通知中心" @click="openNotificationCenter">
        <BellIcon class="w-4 h-4" />
        <span v-if="unreadCount > 0" class="notification-badge">
          {{ unreadCount > 9 ? '9+' : unreadCount }}
        </span>
      </button>

      <div ref="userMenuRef" class="user-menu-container">
        <button class="system-bar-btn user-btn" title="用户菜单" @click="toggleUserMenu">
          <div class="user-avatar">
            <UserIcon class="w-4 h-4" />
          </div>
          <span class="user-name">{{ username }}</span>
          <ChevronDownIcon class="w-3 h-3 dropdown-arrow" :class="{ rotate: userMenuVisible }" />
        </button>

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
          <button class="dropdown-item" @click="emit('open-app', 'user-manager'); userMenuVisible = false">
            <UserIcon class="w-4 h-4" />
            用户资料
          </button>
          <button class="dropdown-item" @click="emit('open-app', 'control-panel'); userMenuVisible = false">
            <Cog6ToothIcon class="w-4 h-4" />
            控制面板
          </button>
          <div class="dropdown-divider"></div>
          <button class="dropdown-item logout" @click="handleLogout">
            <ArrowLeftOnRectangleIcon class="w-4 h-4" />
            退出登录
          </button>
        </div>
      </div>

      <button class="system-bar-btn settings-btn" title="同步备份" @click="emit('open-app', 'sync-manager')">
        <CloudArrowUpIcon class="w-4 h-4" />
      </button>

      <div ref="powerMenuRef" class="power-menu-container">
        <button class="system-bar-btn power-btn" title="系统操作" @click="powerMenuVisible = !powerMenuVisible">
          <PowerIcon class="w-4 h-4" />
        </button>

        <div v-if="powerMenuVisible" class="power-dropdown-menu">
          <button class="dropdown-item power-option" @click="handleReboot">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            重启系统
          </button>
          <button class="dropdown-item power-option shutdown" @click="handleShutdown">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M13 9V3.5L22 12l-9 8.5V15c0-5-5.5-9-10-9z" />
            </svg>
            关机
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  BellIcon,
  UserIcon,
  ChevronDownIcon,
  PowerIcon,
  ArrowLeftOnRectangleIcon,
  CloudArrowUpIcon,
  Cog6ToothIcon,
} from '@heroicons/vue/24/outline'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import { useClock } from '../../composables/useClock'
import { useClickOutside } from '../../composables/useClickOutside'
import { usePowerOperations } from '../../composables/usePowerOperations'

const emit = defineEmits<{
  (e: 'open-app', appId: string): void
}>()

const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const router = useRouter()

const { currentTime } = useClock()
const { reboot, shutdown } = usePowerOperations()

const username = computed(() => authStore.user?.username || 'Admin')
const userRole = computed(() => authStore.user?.role || '管理员')
const unreadCount = computed(() => notificationStore.unreadCount)

const statusText = ref('正常')
const statusClass = ref<'status-normal' | 'status-warning' | 'status-error'>('status-normal')

const userMenuVisible = ref(false)
const powerMenuVisible = ref(false)
const userMenuRef = ref<HTMLElement | null>(null)
const powerMenuRef = ref<HTMLElement | null>(null)

const userMenuOutside = useClickOutside(() => userMenuRef.value, () => { userMenuVisible.value = false })
const powerMenuOutside = useClickOutside(() => powerMenuRef.value, () => { powerMenuVisible.value = false })

watch(userMenuVisible, v => v && userMenuOutside.start())
watch(powerMenuVisible, v => v && powerMenuOutside.start())

const toggleUserMenu = () => {
  userMenuVisible.value = !userMenuVisible.value
}

const openNotificationCenter = () => notificationStore.openCenter()

const handleLogout = () => {
  if (!confirm('确定要退出登录吗？')) return
  authStore.clearToken()
  userMenuVisible.value = false
  router.push('/login')
}

const handleReboot = async () => {
  powerMenuVisible.value = false
  await reboot()
}

const handleShutdown = async () => {
  powerMenuVisible.value = false
  await shutdown()
}
</script>

<style scoped>
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

.notification-btn { position: relative; }

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

.user-menu-container,
.power-menu-container {
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

.dropdown-arrow { transition: transform 0.2s; }
.dropdown-arrow.rotate { transform: rotate(180deg); }

.user-dropdown-menu,
.power-dropdown-menu {
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

.power-dropdown-menu { min-width: 160px; }

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-8px); }
  to   { opacity: 1; transform: translateY(0); }
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

.user-details { flex: 1; }

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

.dropdown-item.logout { color: #ef4444; }
.dropdown-item.logout:hover { background: rgba(239, 68, 68, 0.1); }

.dropdown-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.1);
  margin: 8px 0;
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

.power-option.shutdown { color: #ef4444; }
.power-option.shutdown:hover { background: rgba(239, 68, 68, 0.1); color: #dc2626; }
.power-option:hover { background: rgba(102, 126, 234, 0.1); color: #667eea; }

.settings-btn:hover { background: rgba(102, 126, 234, 0.2); }
.power-btn:hover { background: rgba(239, 68, 68, 0.2); }
</style>
