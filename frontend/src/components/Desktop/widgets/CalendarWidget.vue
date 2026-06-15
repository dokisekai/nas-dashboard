<template>
  <div class="calendar-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <div class="widget-title">日历</div>
      <div class="widget-controls">
        <button @click="prevMonth" class="control-btn">
          <ChevronLeftIcon class="w-4 h-4" />
        </button>
        <button @click="nextMonth" class="control-btn">
          <ChevronRightIcon class="w-4 h-4" />
        </button>
      </div>
    </div>

    <div class="current-month">{{ monthYear }}</div>

    <div class="calendar-grid">
      <div class="weekday" v-for="day in weekdays" :key="day">{{ day }}</div>
      <div
        v-for="date in calendarDates"
        :key="date.date.getTime()"
        class="calendar-date"
        :class="{
          'is-today': date.isToday,
          'is-other-month': date.isOtherMonth,
          'has-events': date.hasEvents
        }"
      >
        {{ date.day }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/vue/24/outline'

interface Props {
  config: {
    showEvents?: boolean
  }
  size: 'small' | 'medium' | 'large'
}

// 定义日历日期接口
interface CalendarDate {
  day: number
  date: Date
  isToday: boolean
  isOtherMonth: boolean
  hasEvents: boolean
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    showEvents: true
  })
})

const currentDate = ref(new Date())
const today = new Date()

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

const monthYear = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth() + 1
  return `${year}年 ${month}月`
})

const calendarDates = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()

  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const startDay = firstDay.getDay()
  const totalDays = lastDay.getDate()

  const dates: CalendarDate[] = []

  // 上个月的日期
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = startDay - 1; i >= 0; i--) {
    dates.push({
      day: prevMonthLastDay - i,
      date: new Date(year, month - 1, prevMonthLastDay - i),
      isToday: false,
      isOtherMonth: true,
      hasEvents: false
    })
  }

  // 当月的日期
  for (let i = 1; i <= totalDays; i++) {
    const date = new Date(year, month, i)
    dates.push({
      day: i,
      date,
      isToday: isSameDay(date, today),
      isOtherMonth: false,
      hasEvents: Math.random() > 0.8 // 模拟事件
    })
  }

  // 下个月的日期
  const remaining = 42 - dates.length
  for (let i = 1; i <= remaining; i++) {
    dates.push({
      day: i,
      date: new Date(year, month + 1, i),
      isToday: false,
      isOtherMonth: true,
      hasEvents: false
    })
  }

  return dates
})

const isSameDay = (date1: Date, date2: Date) => {
  return date1.getFullYear() === date2.getFullYear() &&
         date1.getMonth() === date2.getMonth() &&
         date1.getDate() === date2.getDate()
}

const prevMonth = () => {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() - 1,
    1
  )
}

const nextMonth = () => {
  currentDate.value = new Date(
    currentDate.value.getFullYear(),
    currentDate.value.getMonth() + 1,
    1
  )
}
</script>

<style scoped>
.calendar-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
}

.widget-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.widget-title {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.widget-controls {
  display: flex;
  gap: 4px;
}

.control-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #6b7280;
  transition: all 0.2s ease;
}

.control-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.current-month {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  text-align: center;
  margin-bottom: 12px;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 2px;
  flex: 1;
}

.weekday {
  text-align: center;
  font-size: 11px;
  font-weight: 500;
  color: #9ca3af;
  padding: 4px 0;
}

.calendar-date {
  text-align: center;
  font-size: 12px;
  padding: 6px 2px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #374151;
  position: relative;
}

.calendar-date:hover {
  background: rgba(0, 0, 0, 0.05);
}

.calendar-date.is-today {
  background: #3b82f6;
  color: white;
  font-weight: 600;
}

.calendar-date.is-other-month {
  color: #d1d5db;
}

.calendar-date.has-events::after {
  content: '';
  position: absolute;
  bottom: 2px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background: #ef4444;
  border-radius: 50%;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .widget-title,
  .current-month {
    color: #f9fafb;
  }

  .calendar-date {
    color: #e5e7eb;
  }

  .calendar-date.is-other-month {
    color: #4b5563;
  }

  .calendar-date:hover {
    background: rgba(255, 255, 255, 0.1);
  }
}
</style>
