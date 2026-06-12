<template>
  <div class="weather-widget">
    <div class="weather-header">
      <div class="weather-icon">
        <svg v-if="weatherData.condition === 'sunny'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <svg v-else-if="weatherData.condition === 'cloudy'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z" />
        </svg>
        <svg v-else-if="weatherData.condition === 'rainy'" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 13v8M8 13v8M12 15v8" />
        </svg>
        <svg v-else class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z" />
        </svg>
      </div>
      <div class="weather-location">{{ weatherData.location }}</div>
    </div>

    <div class="weather-temperature">
      <span class="temp-value">{{ weatherData.temperature }}°</span>
      <span class="temp-unit">{{ weatherData.unit === 'celsius' ? 'C' : 'F' }}</span>
    </div>

    <div class="weather-details">
      <div class="detail-item">
        <span class="detail-label">湿度</span>
        <span class="detail-value">{{ weatherData.humidity }}%</span>
      </div>
      <div class="detail-item">
        <span class="detail-label">风速</span>
        <span class="detail-value">{{ weatherData.windSpeed }} m/s</span>
      </div>
      <div class="detail-item">
        <span class="detail-label">体感</span>
        <span class="detail-value">{{ weatherData.feelsLike }}°</span>
      </div>
    </div>

    <div class="weather-condition">
      {{ getConditionText(weatherData.condition) }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const weatherData = ref({
  location: '上海',
  temperature: 24,
  unit: 'celsius' as 'celsius' | 'fahrenheit',
  condition: 'sunny' as 'sunny' | 'cloudy' | 'rainy' | 'snowy',
  humidity: 65,
  windSpeed: 3.2,
  feelsLike: 26
})

let updateInterval: number

const updateWeather = () => {
  // 注意：这是演示数据，实际部署时需要集成真实的天气API（如OpenWeatherMap）
  // 可以考虑使用位置API获取用户位置，然后调用天气服务
  const conditions = ['sunny', 'cloudy', 'rainy', 'snowy'] as const
  const randomCondition = conditions[Math.floor(Math.random() * conditions.length)]

  weatherData.value = {
    ...weatherData.value,
    temperature: Math.round(15 + Math.random() * 15),
    condition: randomCondition,
    humidity: Math.round(40 + Math.random() * 40),
    windSpeed: parseFloat((Math.random() * 10).toFixed(1)),
    feelsLike: Math.round(15 + Math.random() * 15)
  }
}

const getConditionText = (condition: string) => {
  const texts: Record<string, string> = {
    sunny: '晴朗',
    cloudy: '多云',
    rainy: '雨天',
    snowy: '雪天'
  }
  return texts[condition] || '未知'
}

onMounted(() => {
  updateWeather()
  updateInterval = setInterval(updateWeather, 60000) as unknown as number // 每分钟更新一次
})

onUnmounted(() => {
  clearInterval(updateInterval)
})
</script>

<style scoped>
.weather-widget {
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.weather-widget:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(102, 126, 234, 0.25);
}

.weather-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.weather-icon {
  width: 40px;
  height: 40px;
  color: #f59e0b;
}

.weather-location {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

.weather-temperature {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 12px;
}

.temp-value {
  font-size: 36px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.temp-unit {
  font-size: 16px;
  color: #6b7280;
  font-weight: 600;
}

.weather-details {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  margin-bottom: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.detail-label {
  font-size: 10px;
  color: #9ca3af;
  margin-bottom: 2px;
}

.detail-value {
  font-size: 12px;
  color: #1f2937;
  font-weight: 600;
}

.weather-condition {
  text-align: center;
  font-size: 14px;
  color: #667eea;
  font-weight: 500;
  padding-top: 8px;
  border-top: 1px solid rgba(102, 126, 234, 0.1);
}
</style>
