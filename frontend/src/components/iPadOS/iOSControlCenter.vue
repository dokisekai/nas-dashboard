<template>
  <div class="ios-control-center" @click.self="$emit('close')">
    <div class="control-panel" @click.stop>

      <!-- 网络设置 -->
      <div class="control-section">
        <div class="control-grid">
          <div
            class="control-toggle"
            :class="{ active: settings.wifi }"
            @click="toggleSetting('wifi')"
          >
            <WifiIcon class="control-icon" />
            <span>Wi-Fi</span>
            <div class="toggle-status">{{ settings.wifi ? '开启' : '关闭' }}</div>
          </div>

          <div
            class="control-toggle"
            :class="{ active: settings.bluetooth }"
            @click="toggleSetting('bluetooth')"
          >
            <BeakerIcon class="control-icon" />
            <span>蓝牙</span>
            <div class="toggle-status">{{ settings.bluetooth ? '开启' : '关闭' }}</div>
          </div>

          <div
            class="control-toggle"
            :class="{ active: settings.airdrop }"
            @click="toggleSetting('airdrop')"
          >
            <ShareIcon class="control-icon" />
            <span>隔空投送</span>
            <div class="toggle-status">{{ settings.airdrop ? '仅联系人' : '关闭' }}</div>
          </div>

          <div
            class="control-toggle"
            :class="{ active: settings.cellular }"
            @click="toggleSetting('cellular')"
          >
            <SignalIcon class="control-icon" />
            <span>蜂窝网络</span>
            <div class="toggle-status">5G</div>
          </div>
        </div>
      </div>

      <!-- 音量和亮度 -->
      <div class="control-section">
        <div class="control-row">
          <div class="slider-control">
            <div class="slider-label">
              <SunIcon class="label-icon" />
              <span>亮度</span>
            </div>
            <input
              v-model="sliders.brightness"
              type="range"
              min="0"
              max="100"
              class="control-slider"
            >
            <div class="slider-value">{{ sliders.brightness }}%</div>
          </div>

          <div class="slider-control">
            <div class="slider-label">
              <SpeakerWaveIcon class="label-icon" />
              <span>音量</span>
            </div>
            <input
              v-model="sliders.volume"
              type="range"
              min="0"
              max="100"
              class="control-slider"
            >
            <div class="slider-value">{{ sliders.volume }}%</div>
          </div>
        </div>
      </div>

      <!-- 快速功能 -->
      <div class="control-section">
        <div class="control-grid-large">
          <div
            class="control-button"
            :class="{ active: settings.dnd }"
            @click="toggleSetting('dnd')"
          >
            <MoonIcon class="button-icon" />
            <span>勿扰模式</span>
          </div>

          <div
            class="control-button"
            @click="openApp('mirror')"
          >
            <RectangleStackIcon class="button-icon" />
            <span>屏幕镜像</span>
          </div>

          <div
            class="control-button"
            @click="openApp('timer')"
          >
            <ClockIcon class="button-icon" />
            <span>计时器</span>
          </div>

          <div
            class="control-button"
            @click="openApp('calculator')"
          >
            <CalculatorIcon class="button-icon" />
            <span>计算器</span>
          </div>

          <div
            class="control-button"
            @click="openApp('camera')"
          >
            <CameraIcon class="button-icon" />
            <span>相机</span>
          </div>

          <div
            class="control-button"
            :class="{ active: settings.focus }"
            @click="toggleSetting('focus')"
          >
            <TagIcon class="button-icon" />
            <span>专注模式</span>
          </div>
        </div>
      </div>

      <!-- 音乐控制 -->
      <div class="control-section music-section" v-if="currentMusic">
        <div class="music-player">
          <div class="music-album">
            <img :src="currentMusic.albumArt" :alt="currentMusic.album" />
          </div>
          <div class="music-info">
            <div class="music-title">{{ currentMusic.title }}</div>
            <div class="music-artist">{{ currentMusic.artist }}</div>
          </div>
          <div class="music-controls">
            <button class="music-btn" @click="previousTrack">
              <BackwardIcon class="btn-icon" />
            </button>
            <button class="music-btn play-btn" @click="togglePlay">
              <PlayIcon v-if="!isPlaying" class="btn-icon" />
              <PauseIcon v-else class="btn-icon" />
            </button>
            <button class="music-btn" @click="nextTrack">
              <ForwardIcon class="btn-icon" />
            </button>
          </div>
        </div>
      </div>

      <!-- 底部控制 -->
      <div class="control-section bottom-section">
        <div class="bottom-controls">
          <div class="bottom-button" @click="openApp('home')">
            <HomeIcon class="bottom-icon" />
            <span>主屏幕</span>
          </div>

          <div class="bottom-button" @click="openApp('lock')">
            <LockClosedIcon class="bottom-icon" />
            <span>锁屏</span>
          </div>

          <div class="bottom-button" @click="openApp('brightness')">
            <SunIcon class="bottom-icon" />
            <span>亮度</span>
          </div>

          <div class="bottom-button" @click="openApp('volume')">
            <SpeakerWaveIcon class="bottom-icon" />
            <span>音量</span>
          </div>

          <div class="bottom-button" @click="$emit('close')">
            <XMarkIcon class="bottom-icon" />
            <span>关闭</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  WifiIcon,
  ShareIcon,
  SignalIcon,
  SunIcon,
  SpeakerWaveIcon,
  MoonIcon,
  RectangleStackIcon,
  ClockIcon,
  CalculatorIcon,
  CameraIcon,
  TagIcon,
  HomeIcon,
  LockClosedIcon,
  XMarkIcon,
  BackwardIcon,
  ForwardIcon,
  PlayIcon,
  PauseIcon,
  BeakerIcon
} from '@heroicons/vue/24/outline'

defineEmits(['close'])

// 设置状态
const settings = ref({
  wifi: true,
  bluetooth: true,
  airdrop: false,
  cellular: true,
  dnd: false,
  focus: false
})

// 滑块值
const sliders = ref({
  brightness: 70,
  volume: 50
})

// 当前音乐
const currentMusic = ref({
  title: 'Beautiful Day',
  artist: 'U2',
  album: 'All That You Can\'t Leave Behind',
  albumArt: 'https://picsum.photos/seed/music/100/100'
})

const isPlaying = ref(false)

// 方法
const toggleSetting = (key: keyof typeof settings.value) => {
  settings.value[key] = !settings.value[key]
}

const openApp = (appId: string) => {
  console.log('Opening app:', appId)
  // 集成现有的应用打开逻辑
}

const previousTrack = () => {
  console.log('Previous track')
}

const nextTrack = () => {
  console.log('Next track')
}

const togglePlay = () => {
  isPlaying.value = !isPlaying.value
}
</script>

<style scoped>
.ios-control-center {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(20px);
  z-index: 3000;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 60px;
}

.control-panel {
  background: rgba(30, 30, 30, 0.9);
  backdrop-filter: blur(40px);
  border-radius: 32px;
  padding: 20px;
  width: 90%;
  max-width: 400px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.5);
}

.control-section {
  margin-bottom: 20px;
}

.control-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.control-toggle {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.control-toggle:active {
  transform: scale(0.95);
}

.control-toggle.active {
  background: rgba(255, 255, 255, 0.2);
}

.control-icon {
  width: 24px;
  height: 24px;
  color: white;
}

.control-toggle.active .control-icon {
  color: #3b82f6;
}

.toggle-status {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.control-row {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.slider-control {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
}

.slider-label {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 80px;
  color: white;
  font-size: 14px;
}

.label-icon {
  width: 20px;
  height: 20px;
}

.control-slider {
  flex: 1;
  height: 4px;
  -webkit-appearance: none;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
  outline: none;
}

.control-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  cursor: pointer;
}

.slider-value {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
  min-width: 40px;
  text-align: right;
}

.control-grid-large {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.control-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.control-button:active {
  transform: scale(0.95);
}

.control-button.active {
  background: rgba(59, 130, 246, 0.3);
}

.button-icon {
  width: 24px;
  height: 24px;
  color: white;
}

.control-button.active .button-icon {
  color: #3b82f6;
}

.music-section {
  margin-top: 20px;
}

.music-player {
  display: flex;
  align-items: center;
  gap: 16px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 16px;
}

.music-album {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  overflow: hidden;
}

.music-album img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.music-info {
  flex: 1;
}

.music-title {
  font-size: 16px;
  font-weight: 600;
  color: white;
  margin-bottom: 4px;
}

.music-artist {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.music-controls {
  display: flex;
  gap: 8px;
}

.music-btn {
  width: 40px;
  height: 40px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.music-btn:active {
  transform: scale(0.9);
}

.play-btn {
  width: 48px;
  height: 48px;
  background: white;
  color: #333;
}

.btn-icon {
  width: 20px;
  height: 20px;
}

.bottom-section {
  margin-top: 20px;
}

.bottom-controls {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 8px;
}

.bottom-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.bottom-button:active {
  transform: scale(0.95);
}

.bottom-icon {
  width: 24px;
  height: 24px;
  color: white;
}

.bottom-button span {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.8);
}

/* 响应式 */
@media (max-width: 768px) {
  .control-panel {
    width: 95%;
    padding: 16px;
  }

  .control-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .control-grid-large {
    grid-template-columns: repeat(2, 1fr);
  }

  .bottom-controls {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>