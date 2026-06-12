<template>
  <div class="calendar-widget">
    <div class="calendar-header">
      <button class="nav-btn" @click="previousMonth" :disabled="isLoading">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <div class="month-year">
        {{ currentMonth }} {{ currentYear }}
      </div>

      <button class="nav-btn" @click="nextMonth" :disabled="isLoading">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>

    <div class="calendar-weekdays">
      <div v-for="day in weekdays" :key="day" class="weekday">
        {{ day }}
      </div>
    </div>

    <div class="calendar-days">
      <div
        v-for="(date, index) in calendarDays"
        :key="index"
        :class="['day-cell', {
          'other-month': date.isOtherMonth,
          'today': date.isToday,
          'selected': date.isSelected,
          'has-event': date.hasEvent
        }]"
        @click="selectDate(date)"
      >
        {{ date.day }}
        <div v-if="date.hasEvent" class="event-dot"></div>
      </div>
    </div>

    <div v-if="selectedDate" class="selected-date-info">
      <div class="selected-date">{{ selectedDateText }}</div>
      <div v-if="selectedDateEvents.length > 0" class="events-preview">
        <div v-for="event in selectedDateEvents.slice(0, 2)" :key="event.id" class="event-item">
          <div class="event-dot" :style="{ background: event.color }"></div>
          <span class="event-title">{{ event.title }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const currentDate = ref(new Date())
const selectedDate = ref<Date | null>(null)
const isLoading = ref(false)

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

const monthNames = [
  '一月', '二月', '三月', '四月', '五月', '六月',
  '七月', '八月', '九月', '十月', '十一月', '十二月'
]

const currentYear = computed(() => currentDate.value.getFullYear())
const currentMonth = computed(() => monthNames[currentDate.value.getMonth()])

const selectedDateText = computed(() => {
  if (!selectedDate.value) return ''
  const date = selectedDate.value
  return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`
})

// 模拟事件数据
const events = ref([
  { id: 1, date: new Date(), title: '系统维护', color: '#ef4444' },
  { id: 2, date: new Date(Date.now() + 86400000), title: '数据备份', color: '#10b981' },
  { id: 3, date: new Date(Date.now() + 172800000), title: '安全更新', color: '#3b82f6' }
])

const selectedDateEvents = computed(() => {
  if (!selectedDate.value) return []
  return events.value.filter(event => {
    const eventDate = new Date(event.date)
    return eventDate.toDateString() === selectedDate.value?.toDateString()
  })
})

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()

  const firstDayOfMonth = new Date(year, month, 1)
  const lastDayOfMonth = new Date(year, month + 1, 0)

  const firstDayWeek = firstDayOfMonth.getDay()
  const daysInMonth = lastDayOfMonth.getDate()

  const days = []

  // 上个月的日期
  const prevMonthLastDay = new Date(year, month, 0).getDate()
  for (let i = firstDayWeek - 1; i >= 0; i--) {
    const day = prevMonthLastDay - i
    const date = new Date(year, month - 1, day)
    days.push({
      day,
      isOtherMonth: true,
      isToday: false,
      isSelected: false,
      hasEvent: false,
      date
    })
  }

  // 当前月的日期
  const today = new Date()
  for (let i = 1; i <= daysInMonth; i++) {
    const date = new Date(year, month, i)
    const isToday = date.toDateString() === today.toDateString()
    const isSelected = selectedDate.value && date.toDateString() === selectedDate.value.toDateString()
    const hasEvent = events.value.some(event => new Date(event.date).toDateString() === date.toDateString())

    days.push({
      day: i,
      isOtherMonth: false,
      isToday,
      isSelected,
      hasEvent,
      date
    })
  }

  // 下个月的日期
  const remainingCells = 42 - days.length // 6行 x 7列 = 42
  for (let i = 1; i <= remainingCells; i++) {
    const date = new Date(year, month + 1, i)
    days.push({
      day: i,
      isOtherMonth: true,
      isToday: false,
      isSelected: false,
      hasEvent: false,
      date
    })
  }

  return days
})

const previousMonth = () => {
  isLoading.value = true
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() - 1)
  setTimeout(() => {
    isLoading.value = false
  }, 200)
}

const nextMonth = () => {
  isLoading.value = true
  currentDate.value = new Date(currentDate.value.getFullYear(), currentDate.value.getMonth() + 1)
  setTimeout(() => {
    isLoading.value = false
  }, 200)
}

const selectDate = (date: any) => {
  if (!date.isOtherMonth) {
    selectedDate.value = date.date
  }
}

onMounted(() => {
  // 默认选中今天
  selectedDate.value = new Date()
})
</script>

<style scoped>
.calendar-widget {
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.calendar-widget:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.calendar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.month-year {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  text-align: center;
  flex: 1;
}

.nav-btn {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(102, 126, 234, 0.2);
  border-radius: 6px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  cursor: pointer;
  transition: all 0.2s ease;
}

.nav-btn:hover:not(:disabled) {
  background: rgba(102, 126, 234, 0.2);
  border-color: rgba(102, 126, 234, 0.3);
}

.nav-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.calendar-weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 8px;
}

.weekday {
  text-align: center;
  font-size: 11px;
  font-weight: 600;
  color: #9ca3af;
  padding: 4px;
}

.calendar-days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
}

.day-cell {
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  border-radius: 8px;
  cursor: pointer;
  position: relative;
  transition: all 0.2s ease;
  color: #374151;
  background: transparent;
}

.day-cell:hover:not(.other-month) {
  background: rgba(102, 126, 234, 0.1);
}

.day-cell.other-month {
  color: #d1d5db;
}

.day-cell.today {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-weight: 600;
}

.day-cell.selected {
  border: 2px solid #667eea;
}

.day-cell.has-event::after {
  content: '';
  position: absolute;
  bottom: 4px;
  left: 50%;
  transform: translateX(-50%);
  width: 4px;
  height: 4px;
  background: #ef4444;
  border-radius: 50%;
}

.selected-date-info {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}

.selected-date {
  font-size: 13px;
  color: #667eea;
  font-weight: 600;
  margin-bottom: 8px;
}

.events-preview {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.event-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
}

.event-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.event-title {
  color: #6b7280;
  font-size: 11px;
}
</style>
