<template>
  <div
    class="ipados-widget"
    :class="[
      `widget-${size}`,
      `widget-${widget.type}`,
      { 'widget-dark': isDark }
    ]"
  >
    <!-- 时钟小组件 -->
    <div v-if="widget.type === 'clock'" class="clock-widget">
      <div class="clock-time">{{ currentTime }}</div>
      <div class="clock-date">{{ currentDate }}</div>
    </div>

    <!-- 天气小组件 -->
    <div v-else-if="widget.type === 'weather'" class="weather-widget">
      <div class="weather-location">北京</div>
      <div class="weather-temp">26°</div>
      <div class="weather-desc">晴朗</div>
      <div class="weather-range">H:28° L:18°</div>
    </div>

    <!-- 日历小组件 -->
    <div v-else-if="widget.type === 'calendar'" class="calendar-widget">
      <div class="calendar-header">{{ calendarTitle }}</div>
      <div class="calendar-date">{{ calendarDate }}</div>
      <div class="calendar-events">
        <div v-for="event in todaysEvents" :key="event.id" class="event-item">
          <div class="event-time">{{ event.time }}</div>
          <div class="event-title">{{ event.title }}</div>
        </div>
      </div>
    </div>

    <!-- 电池小组件 -->
    <div v-else-if="widget.type === 'battery'" class="battery-widget">
      <div class="battery-icon" :class="batteryLevel">
        <div class="battery-fill" :style="{ width: batteryPercent + '%' }"></div>
      </div>
      <div class="battery-info">
        <div class="battery-percent">{{ batteryPercent }}%</div>
        <div class="battery-status">{{ batteryStatus }}</div>
      </div>
    </div>

    <!-- 存储小组件 -->
    <div v-else-if="widget.type === 'storage'" class="storage-widget">
      <div class="storage-header">
        <div class="storage-title">存储空间</div>
        <div class="storage-used">{{ usedStorage }} GB 已用</div>
      </div>
      <div class="storage-chart">
        <div class="storage-bar">
          <div
            class="storage-fill"
            :style="{
              width: storagePercent + '%',
              background: storageColor
            }"
          ></div>
        </div>
        <div class="storage-labels">
          <span>{{ storagePercent }}%</span>
          <span>{{ totalStorage }} GB</span>
        </div>
      </div>
      <div class="storage-breakdown">
        <div class="storage-item">
          <div class="item-color" style="background: #3b82f6"></div>
          <span>系统</span>
          <span>{{ systemStorage }} GB</span>
        </div>
        <div class="storage-item">
          <div class="item-color" style="background: #10b981"></div>
          <span>媒体</span>
          <span>{{ mediaStorage }} GB</span>
        </div>
        <div class="storage-item">
          <div class="item-color" style="background: #f59e0b"></div>
          <span>应用</span>
          <span>{{ appStorage }} GB</span>
        </div>
      </div>
    </div>

    <!-- 系统监控小组件 -->
    <div v-else-if="widget.type === 'system'" class="system-widget">
      <div class="system-metrics">
        <div class="metric-item">
          <div class="metric-icon">CPU</div>
          <div class="metric-value">{{ cpuUsage }}%</div>
        </div>
        <div class="metric-item">
          <div class="metric-icon">RAM</div>
          <div class="metric-value">{{ ramUsage }}%</div>
        </div>
      </div>
      <div class="system-chart">
        <canvas ref="chartCanvas"></canvas>
      </div>
    </div>

    <!-- 通知小组件 -->
    <div v-else-if="widget.type === 'notifications'" class="notifications-widget">
      <div class="notifications-header">
        <BellIcon class="header-icon" />
        <span>通知</span>
      </div>
      <div class="notifications-list">
        <div
          v-for="notification in recentNotifications"
          :key="notification.id"
          class="notification-item"
        >
          <div class="notification-app">{{ notification.app }}</div>
          <div class="notification-message">{{ notification.message }}</div>
          <div class="notification-time">{{ notification.time }}</div>
        </div>
      </div>
    </div>

    <!-- 快速操作小组件 -->
    <div v-else-if="widget.type === 'quick-actions'" class="quick-actions-widget">
      <div class="quick-actions-grid">
        <div
          v-for="action in quickActions"
          :key="action.id"
          class="quick-action-btn"
          :class="{ active: action.active }"
          @click="toggleAction(action)"
        >
          <component :is="action.icon" class="action-icon" />
          <span>{{ action.name }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { BellIcon, WifiIcon, BeakerIcon, PowerIcon } from '@heroicons/vue/24/outline'

interface Widget {
  id: string
  size: 'small' | 'medium' | 'large' | 'extra-large'
  type: string
}

interface Props {
  widget: Widget
  size?: 'small' | 'medium' | 'large' | 'extra-large'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'medium'
})

// 状态
const isDark = ref(false)
const currentTime = ref('')
const currentDate = ref('')
const cpuUsage = ref(45)
const ramUsage = ref(62)
const batteryPercent = ref(87)
const batteryStatus = ref('正常')
const usedStorage = ref(256)
const totalStorage = ref(512)

// 日历数据
const calendarTitle = ref('今日日程')
const calendarDate = ref('')

// 今日事件
const todaysEvents = ref([
  { id: 1, time: '10:00', title: '团队会议' },
  { id: 2, time: '14:30', title: '项目审查' },
  { id: 3, time: '16:00', title: '客户电话' }
])

// 通知数据
const recentNotifications = ref([
  { id: 1, app: '邮件', message: '您有3封新邮件', time: '2分钟前' },
  { id: 2, app: '信息', message: '新消息', time: '15分钟前' },
  { id: 3, app: '日历', message: '会议提醒', time: '1小时前' }
])

// 快速操作
const quickActions = ref([
  { id: 1, name: 'WiFi', icon: WifiIcon, active: true },
  { id: 2, name: '蓝牙', icon: BeakerIcon, active: false },
  { id: 3, name: '省电', icon: PowerIcon, active: false }
])

// 计算属性
const batteryLevel = computed(() => {
  if (batteryPercent.value > 80) return 'high'
  if (batteryPercent.value > 20) return 'medium'
  return 'low'
})

const storagePercent = computed(() => {
  return Math.round((usedStorage.value / totalStorage.value) * 100)
})

const storageColor = computed(() => {
  const percent = storagePercent.value
  if (percent > 80) return '#ef4444'
  if (percent > 60) return '#f59e0b'
  return '#10b981'
})

const systemStorage = computed(() => Math.round(usedStorage.value * 0.3))
const mediaStorage = computed(() => Math.round(usedStorage.value * 0.5))
const appStorage = computed(() => Math.round(usedStorage.value * 0.2))

// 方法
const updateDateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: false
  })

  const options = { month: 'long', day: 'numeric', weekday: 'long' }
  currentDate.value = now.toLocaleDateString('zh-CN', options)
  calendarDate.value = now.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric'
  })
}

const updateSystemMetrics = () => {
  // 模拟系统数据更新
  cpuUsage.value = Math.floor(Math.random() * 30) + 30
  ramUsage.value = Math.floor(Math.random() * 20) + 50
}

const toggleAction = (action: any) => {
  action.active = !action.active
}

// 生命周期
let timeInterval: any
let metricsInterval: any

onMounted(() => {
  updateDateTime()
  timeInterval = setInterval(updateDateTime, 1000)
  metricsInterval = setInterval(updateSystemMetrics, 5000)
})

onUnmounted(() => {
  if (timeInterval) clearInterval(timeInterval)
  if (metricsInterval) clearInterval(metricsInterval)
})
</script>

<style scoped>
.ipados-widget {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(30px);
  border-radius: 20px;
  padding: 16px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  color: #333;
  transition: all 0.3s ease;
}

.ipados-widget.widget-dark {
  background: rgba(30, 30, 30, 0.8);
  color: white;
}

/* 小小组件 */
.widget-small {
  width: 140px;
  height: 140px;
}

/* 中等小组件 */
.widget-medium {
  width: 300px;
  height: 140px;
}

/* 大小组件 */
.widget-large {
  width: 300px;
  height: 300px;
}

/* 超大小组件 */
.widget-extra-large {
  width: 620px;
  height: 300px;
}

/* 时钟小组件 */
.clock-widget {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.clock-time {
  font-size: 32px;
  font-weight: 700;
  color: #333;
}

.clock-date {
  font-size: 14px;
  color: #666;
  margin-top: 4px;
}

.widget-dark .clock-time,
.widget-dark .clock-date {
  color: white;
}

/* 天气小组件 */
.weather-widget {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  height: 100%;
}

.weather-location {
  font-size: 14px;
  color: #666;
}

.weather-temp {
  font-size: 48px;
  font-weight: 300;
  color: #333;
}

.weather-desc {
  font-size: 18px;
  color: #666;
}

.weather-range {
  font-size: 12px;
  color: #999;
}

/* 日历小组件 */
.calendar-widget {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.calendar-header {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}

.calendar-date {
  font-size: 14px;
  color: #666;
  margin-bottom: 16px;
}

.calendar-events {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.event-item {
  display: flex;
  gap: 12px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.event-time {
  font-size: 12px;
  color: #666;
  min-width: 40px;
}

.event-title {
  font-size: 14px;
  color: #333;
}

/* 电池小组件 */
.battery-widget {
  display: flex;
  align-items: center;
  gap: 16px;
  height: 100%;
}

.battery-icon {
  width: 48px;
  height: 24px;
  border: 2px solid #333;
  border-radius: 4px;
  position: relative;
  padding: 2px;
}

.battery-icon::after {
  content: '';
  position: absolute;
  right: -6px;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 8px;
  background: #333;
  border-radius: 0 2px 2px 0;
}

.battery-fill {
  height: 100%;
  background: #10b981;
  border-radius: 2px;
  transition: all 0.3s ease;
}

.battery-icon.low .battery-fill {
  background: #ef4444;
}

.battery-icon.medium .battery-fill {
  background: #f59e0b;
}

.battery-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.battery-percent {
  font-size: 24px;
  font-weight: 600;
  color: #333;
}

.battery-status {
  font-size: 12px;
  color: #666;
}

/* 存储小组件 */
.storage-widget {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.storage-header {
  margin-bottom: 16px;
}

.storage-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 4px;
}

.storage-used {
  font-size: 12px;
  color: #666;
}

.storage-chart {
  margin-bottom: 16px;
}

.storage-bar {
  width: 100%;
  height: 8px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

.storage-fill {
  height: 100%;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.storage-labels {
  display: flex;
  justify-content: space-between;
  margin-top: 4px;
  font-size: 12px;
  color: #666;
}

.storage-breakdown {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.storage-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.item-color {
  width: 8px;
  height: 8px;
  border-radius: 2px;
}

/* 系统监控小组件 */
.system-widget {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.system-metrics {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}

.metric-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.metric-icon {
  font-size: 12px;
  font-weight: 600;
  color: #666;
  min-width: 32px;
}

.metric-value {
  font-size: 20px;
  font-weight: 600;
  color: #333;
}

.system-chart {
  flex: 1;
  display: flex;
  align-items: center;
}

/* 通知小组件 */
.notifications-widget {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.notifications-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}

.header-icon {
  width: 20px;
  height: 20px;
}

.notifications-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.notification-item {
  padding: 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.notification-app {
  font-size: 12px;
  color: #666;
  margin-bottom: 2px;
}

.notification-message {
  font-size: 14px;
  color: #333;
  margin-bottom: 2px;
}

.notification-time {
  font-size: 11px;
  color: #999;
}

/* 快速操作小组件 */
.quick-actions-widget {
  height: 100%;
}

.quick-actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
  height: 100%;
}

.quick-action-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.quick-action-btn:active {
  transform: scale(0.95);
}

.quick-action-btn.active {
  background: rgba(59, 130, 246, 0.2);
}

.action-icon {
  width: 24px;
  height: 24px;
  color: #333;
}

.quick-action-btn.active .action-icon {
  color: #3b82f6;
}

/* 响应式 */
@media (max-width: 768px) {
  .widget-medium {
    width: 100%;
  }

  .widget-large {
    width: 100%;
  }

  .widget-extra-large {
    width: 100%;
  }
}
</style>