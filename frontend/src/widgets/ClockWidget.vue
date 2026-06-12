<template>
  <div class="clock-widget">
    <div class="clock-time">{{ currentTime }}</div>
    <div class="clock-date">{{ currentDate }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const currentTime = ref('')
const currentDate = ref('')
let updateInterval: number

const updateTime = () => {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  })
  currentDate.value = now.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    weekday: 'long'
  })
}

onMounted(() => {
  updateTime()
  updateInterval = setInterval(updateTime, 1000) as unknown as number
})

onUnmounted(() => {
  clearInterval(updateInterval)
})
</script>

<style scoped>
.clock-widget {
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  text-align: center;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.clock-widget:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.clock-time {
  font-size: 32px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 8px;
  letter-spacing: 2px;
}

.clock-date {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}
</style>
