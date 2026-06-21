<template>
  <div class="desktop-widgets">
    <div class="widget-notification-trigger" @click="$emit('open-notifications')">
      <BellIcon class="widget-notification-icon" />
      <span v-if="unreadCount > 0" class="widget-notification-badge">
        {{ unreadCount }}
      </span>
    </div>

    <ClockWidget class="widget clock-widget" />
    <SystemStatusWidget class="widget system-widget" />
    <WeatherWidget class="widget weather-widget" />
    <CalendarWidget class="widget calendar-widget" />
    <QuickShortcutsWidget class="widget shortcuts-widget" />
  </div>
</template>

<script setup lang="ts">
import { BellIcon } from '@heroicons/vue/24/outline'
import ClockWidget from '../../widgets/ClockWidget.vue'
import SystemStatusWidget from '../../widgets/SystemStatusWidget.vue'
import WeatherWidget from '../../widgets/WeatherWidget.vue'
import CalendarWidget from '../../widgets/CalendarWidget.vue'
import QuickShortcutsWidget from '../../widgets/QuickShortcutsWidget.vue'

defineProps<{
  unreadCount: number
}>()

defineEmits<{
  (e: 'open-notifications'): void
}>()
</script>

<style scoped>
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
  from { opacity: 0; transform: translateX(20px); }
  to   { opacity: 1; transform: translateX(0); }
}

.clock-widget,
.system-widget,
.weather-widget {
  grid-column: span 1;
}

.calendar-widget,
.shortcuts-widget {
  grid-column: span 2;
}

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
</style>
