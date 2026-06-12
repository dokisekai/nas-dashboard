<template>
  <div class="clock-widget" :class="`widget-${size}`">
    <div class="clock-display">
      <div class="time">{{ currentTime }}</div>
      <div class="date" v-if="config.showDate">{{ currentDate }}</div>
    </div>
    <div class="timezone" v-if="size !== 'small'">{{ timezone }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Props {
  config: {
    showDate?: boolean
    showSeconds?: boolean
    format24?: boolean
  }
  size: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    showDate: true,
    showSeconds: false,
    format24: true
  })
})

const currentTime = ref('')
const currentDate = ref('')
const timezone = ref('UTC+8')

const updateTime = () => {
  const now = new Date()

  // 格式化时间
  const hours = props.config.format24
    ? now.getHours().toString().padStart(2, '0')
    : (now.getHours() % 12 || 12).toString().padStart(2, '0')

  const minutes = now.getMinutes().toString().padStart(2, '0')

  let timeString = `${hours}:${minutes}`

  if (props.config.showSeconds) {
    const seconds = now.getSeconds().toString().padStart(2, '0')
    timeString += `:${seconds}`
  }

  if (!props.config.format24) {
    const ampm = now.getHours() >= 12 ? 'PM' : 'AM'
    timeString += ` ${ampm}`
  }

  currentTime.value = timeString

  // 格式化日期
  if (props.config.showDate) {
    const options: Intl.DateTimeFormatOptions = {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      weekday: 'short'
    }
    currentDate.value = now.toLocaleDateString('zh-CN', options)
  }
}

let timer: number | null = null

onMounted(() => {
  updateTime()
  timer = window.setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style scoped>
.clock-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.9), rgba(255, 255, 255, 0.7));
  backdrop-filter: blur(10px);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.clock-display {
  text-align: center;
}

.time {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  line-height: 1.2;
  font-variant-numeric: tabular-nums;
}

.widget-small .time {
  font-size: 24px;
}

.widget-medium .time {
  font-size: 28px;
}

.date {
  font-size: 14px;
  color: #6b7280;
  margin-top: 8px;
}

.timezone {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 4px;
}

.widget-small .timezone {
  display: none;
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .clock-widget {
    background: linear-gradient(135deg, rgba(0, 0, 0, 0.8), rgba(0, 0, 0, 0.6));
  }

  .time {
    color: #f9fafb;
  }

  .date {
    color: #9ca3af;
  }

  .timezone {
    color: #6b7280;
  }
}
</style>