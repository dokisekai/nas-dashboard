<template>
  <div class="weather-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <div class="widget-title">天气</div>
      <div class="widget-location">{{ location }}</div>
    </div>

    <div class="weather-content">
      <div class="weather-main">
        <div class="weather-icon">
          <component :is="getWeatherIcon()" class="w-full h-full" />
        </div>
        <div class="temperature">{{ temperature }}°</div>
        <div class="condition">{{ condition }}</div>
      </div>

      <div class="weather-details" v-if="size !== 'small'">
        <div class="detail-item">
          <CloudIcon class="icon" />
          <div class="detail-info">
            <span class="label">湿度</span>
            <span class="value">{{ humidity }}%</span>
          </div>
        </div>
        <div class="detail-item">
          <WindIcon class="icon" />
          <div class="detail-info">
            <span class="label">风速</span>
            <span class="value">{{ windSpeed }} km/h</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CloudIcon, SunIcon, CloudRainIcon, SnowIcon, WindIcon } from '@heroicons/vue/24/outline'

interface Props {
  config: {
    location?: string
    units?: 'celsius' | 'fahrenheit'
  }
  size: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    location: '北京',
    units: 'celsius'
  })
})

const location = ref(props.config.location)
const temperature = ref(24)
const condition = ref('晴')
const humidity = ref(65)
const windSpeed = ref(12)

const getWeatherIcon = () => {
  switch (condition.value) {
    case '晴':
      return SunIcon
    case '雨':
      return CloudRainIcon
    case '雪':
      return SnowIcon
    default:
      return CloudIcon
  }
}

onMounted(() => {
  // 模拟天气数据
  temperature.value = Math.round(Math.random() * 20 + 10)
  humidity.value = Math.round(Math.random() * 40 + 40)
  windSpeed.value = Math.round(Math.random() * 20 + 5)
})
</script>

<style scoped>
.weather-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
}

.widget-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.widget-title {
  font-size: 14px;
  font-weight: 600;
}

.widget-location {
  font-size: 12px;
  opacity: 0.8;
}

.weather-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.weather-main {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  text-align: center;
}

.weather-icon {
  width: 48px;
  height: 48px;
}

.temperature {
  font-size: 36px;
  font-weight: 700;
}

.widget-small .temperature {
  font-size: 28px;
}

.condition {
  font-size: 14px;
  opacity: 0.9;
}

.weather-details {
  display: flex;
  gap: 16px;
  margin-top: 16px;
  justify-content: center;
}

.detail-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
}

.icon {
  width: 16px;
  height: 16px;
  opacity: 0.8;
}

.detail-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.label {
  opacity: 0.7;
}

.value {
  font-weight: 600;
}

.widget-small .weather-details {
  display: none;
}
</style>
