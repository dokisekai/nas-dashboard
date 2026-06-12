<template>
  <div class="theme-manager">
    <div class="manager-header">
      <h2>桌面设置</h2>
      <div class="manager-tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          class="tab-btn"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          <component :is="tab.icon" class="w-4 h-4" />
          {{ tab.name }}
        </button>
      </div>
    </div>

    <div class="manager-content">
      <!-- 背景设置 -->
      <div v-if="activeTab === 'background'" class="tab-content">
        <div class="content-section">
          <h3>背景图片</h3>
          <div class="background-options">
            <div
              v-for="bg in builtinBackgrounds"
              :key="bg.id"
              class="background-option"
              :class="{ active: settings.background.id === bg.id }"
              @click="selectBackground(bg)"
            >
              <div class="background-preview" :style="{ background: bg.preview }">
                <div v-if="settings.background.id === bg.id" class="check-mark">
                  <CheckIcon class="w-6 h-6" />
                </div>
              </div>
              <div class="background-name">{{ bg.name }}</div>
            </div>
          </div>

          <div class="custom-background">
            <h4>自定义背景</h4>
            <div class="upload-area">
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                @change="handleFileUpload"
                class="hidden-input"
              />
              <button @click="triggerFileUpload" class="upload-btn">
                <PhotoIcon class="w-8 h-8" />
                <span>上传图片</span>
              </button>
            </div>
          </div>
        </div>

        <div class="content-section">
          <h3>背景设置</h3>
          <div class="setting-item">
            <label>填充方式</label>
            <select v-model="settings.background.fit" class="form-select">
              <option value="cover">覆盖</option>
              <option value="contain">包含</option>
              <option value="stretch">拉伸</option>
              <option value="tile">平铺</option>
            </select>
          </div>

          <div class="setting-item">
            <label class="checkbox-label">
              <input
                v-model="settings.background.darken"
                type="checkbox"
                class="checkbox"
              />
              <span>加深背景（提高可读性）</span>
            </label>
          </div>
        </div>
      </div>

      <!-- 主题设置 -->
      <div v-if="activeTab === 'theme'" class="tab-content">
        <div class="content-section">
          <h3>主题选择</h3>
          <div class="theme-options">
            <div
              v-for="theme in themes"
              :key="theme.id"
              class="theme-option"
              :class="{ active: settings.theme.id === theme.id }"
              @click="selectTheme(theme)"
            >
              <div class="theme-preview">
                <div class="theme-colors">
                  <div
                    v-for="color in theme.colors"
                    :key="color"
                    class="theme-color"
                    :style="{ background: color }"
                  ></div>
                </div>
                <div v-if="settings.theme.id === theme.id" class="check-mark">
                  <CheckIcon class="w-6 h-6" />
                </div>
              </div>
              <div class="theme-name">{{ theme.name }}</div>
              <div class="theme-description">{{ theme.description }}</div>
            </div>
          </div>
        </div>

        <div class="content-section">
          <h3>颜色模式</h3>
          <div class="mode-options">
            <button
              v-for="mode in colorModes"
              :key="mode.value"
              class="mode-btn"
              :class="{ active: settings.theme.mode === mode.value }"
              @click="settings.theme.mode = mode.value"
            >
              <component :is="mode.icon" class="w-6 h-6" />
              <span>{{ mode.label }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- 外观设置 -->
      <div v-if="activeTab === 'appearance'" class="tab-content">
        <div class="content-section">
          <h3>字体</h3>
          <div class="setting-item">
            <label>系统字体</label>
            <select v-model="settings.appearance.font" class="form-select">
              <option value="system">系统默认</option>
              <option value="sans-serif">无衬线</option>
              <option value="serif">衬线</option>
              <option value="monospace">等宽</option>
            </select>
          </div>

          <div class="setting-item">
            <label>字体大小</label>
            <div class="size-slider">
              <input
                v-model.number="settings.appearance.fontSize"
                type="range"
                min="12"
                max="20"
                step="1"
                class="slider"
              />
              <span class="size-value">{{ settings.appearance.fontSize }}px</span>
            </div>
          </div>
        </div>

        <div class="content-section">
          <h3>效果</h3>
          <div class="setting-item">
            <label class="checkbox-label">
              <input
                v-model="settings.appearance.blur"
                type="checkbox"
                class="checkbox"
              />
              <span>背景模糊效果</span>
            </label>
          </div>

          <div class="setting-item">
            <label class="checkbox-label">
              <input
                v-model="settings.appearance.transparency"
                type="checkbox"
                class="checkbox"
              />
              <span>窗口透明效果</span>
            </label>
          </div>

          <div class="setting-item">
            <label class="checkbox-label">
              <input
                v-model="settings.appearance.animations"
                type="checkbox"
                class="checkbox"
              />
              <span>动画效果</span>
            </label>
          </div>
        </div>
      </div>

      <!-- 屏幕保护设置 -->
      <div v-if="activeTab === 'screensaver'" class="tab-content">
        <div class="content-section">
          <h3>屏幕保护</h3>
          <div class="setting-item">
            <label class="checkbox-label">
              <input
                v-model="settings.screensaver.enabled"
                type="checkbox"
                class="checkbox"
              />
              <span>启用屏幕保护</span>
            </label>
          </div>

          <template v-if="settings.screensaver.enabled">
            <div class="setting-item">
              <label>启动时间（分钟）</label>
              <input
                v-model.number="settings.screensaver.delay"
                type="number"
                min="1"
                max="60"
                class="form-input"
              />
            </div>

            <div class="setting-item">
              <label>保护程序类型</label>
              <select v-model="settings.screensaver.type" class="form-select">
                <option value="clock">时钟</option>
                <option value="photos">照片幻灯片</option>
                <option value="gradient">渐变动画</option>
                <option value="matrix">矩阵</option>
              </select>
            </div>
          </template>
        </div>
      </div>
    </div>

    <div class="manager-footer">
      <button @click="resetDefaults" class="btn btn-secondary">重置默认</button>
      <button @click="saveSettings" class="btn btn-primary">应用</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import {
  PhotoIcon,
  CheckIcon,
  SunIcon,
  MoonIcon,
  ComputerDesktopIcon
} from '@heroicons/vue/24/outline'

interface Props {
  currentSettings: any
}

interface Emits {
  (e: 'save', settings: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const activeTab = ref('background')
const fileInput = ref<HTMLInputElement | null>(null)

const tabs = [
  { id: 'background', name: '背景', icon: 'PhotoIcon' },
  { id: 'theme', name: '主题', icon: 'SunIcon' },
  { id: 'appearance', name: '外观', icon: 'ComputerDesktopIcon' },
  { id: 'screensaver', name: '屏幕保护', icon: 'MoonIcon' }
]

const settings = reactive({
  background: {
    id: props.currentSettings.background?.id || 'gradient1',
    type: props.currentSettings.background?.type || 'builtin',
    url: props.currentSettings.background?.url || '',
    fit: props.currentSettings.background?.fit || 'cover',
    darken: props.currentSettings.background?.darken ?? false
  },
  theme: {
    id: props.currentSettings.theme?.id || 'light',
    mode: props.currentSettings.theme?.mode || 'auto'
  },
  appearance: {
    font: props.currentSettings.appearance?.font || 'system',
    fontSize: props.currentSettings.appearance?.fontSize || 14,
    blur: props.currentSettings.appearance?.blur ?? true,
    transparency: props.currentSettings.appearance?.transparency ?? true,
    animations: props.currentSettings.appearance?.animations ?? true
  },
  screensaver: {
    enabled: props.currentSettings.screensaver?.enabled ?? false,
    delay: props.currentSettings.screensaver?.delay || 5,
    type: props.currentSettings.screensaver?.type || 'clock'
  }
})

const builtinBackgrounds = [
  {
    id: 'gradient1',
    name: '蓝色渐变',
    type: 'gradient',
    preview: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    id: 'gradient2',
    name: '日落渐变',
    type: 'gradient',
    preview: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    id: 'gradient3',
    name: '海洋渐变',
    type: 'gradient',
    preview: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    id: 'gradient4',
    name: '森林渐变',
    type: 'gradient',
    preview: 'linear-gradient(135deg, #38f9d7 0%, #43e97b 100%)'
  },
  {
    id: 'dark1',
    name: '深色模式',
    type: 'gradient',
    preview: 'linear-gradient(135deg, #0f0f0f 0%, #1a1a1a 100%)'
  }
]

const themes = [
  {
    id: 'light',
    name: '浅色',
    description: '明亮清爽的主题',
    colors: ['#ffffff', '#f3f4f6', '#e5e7eb', '#d1d5db']
  },
  {
    id: 'dark',
    name: '深色',
    description: '护眼的深色主题',
    colors: ['#1f2937', '#111827', '#0f0f0f', '#000000']
  },
  {
    id: 'ocean',
    name: '海洋',
    description: '清新的蓝色主题',
    colors: ['#0ea5e9', '#0284c7', '#0369a1', '#075985']
  },
  {
    id: 'forest',
    name: '森林',
    description: '自然的绿色主题',
    colors: ['#22c55e', '#16a34a', '#15803d', '#166534']
  }
]

const colorModes = [
  { value: 'light', label: '浅色', icon: 'SunIcon' },
  { value: 'dark', label: '深色', icon: 'MoonIcon' },
  { value: 'auto', label: '自动', icon: 'ComputerDesktopIcon' }
]

const selectBackground = (bg: any) => {
  settings.background.id = bg.id
  settings.background.type = 'builtin'
}

const selectTheme = (theme: any) => {
  settings.theme.id = theme.id
}

const triggerFileUpload = () => {
  fileInput.value?.click()
}

const handleFileUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      settings.background.type = 'custom'
      settings.background.url = e.target?.result as string
      settings.background.id = 'custom'
    }
    reader.readAsDataURL(file)
  }
}

const saveSettings = () => {
  emit('save', settings)
}

const resetDefaults = () => {
  settings.background.id = 'gradient1'
  settings.background.type = 'builtin'
  settings.background.fit = 'cover'
  settings.background.darken = false

  settings.theme.id = 'light'
  settings.theme.mode = 'auto'

  settings.appearance.font = 'system'
  settings.appearance.fontSize = 14
  settings.appearance.blur = true
  settings.appearance.transparency = true
  settings.appearance.animations = true

  settings.screensaver.enabled = false
  settings.screensaver.delay = 5
  settings.screensaver.type = 'clock'
}
</script>

<style scoped>
.theme-manager {
  width: 100%;
  height: 100%;
  max-width: 800px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.manager-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.manager-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.manager-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: transparent;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.tab-btn.active {
  background: #3b82f6;
  color: white;
}

.manager-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.content-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.content-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.content-section h4 {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.background-options,
.theme-options {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 16px;
}

.background-option,
.theme-option {
  cursor: pointer;
  transition: all 0.2s ease;
}

.background-option:hover,
.theme-option:hover {
  transform: translateY(-4px);
}

.background-preview,
.theme-preview {
  position: relative;
  aspect-ratio: 16/10;
  border-radius: 12px;
  overflow: hidden;
  border: 3px solid transparent;
  transition: all 0.2s ease;
}

.background-option.active .background-preview,
.theme-option.active .theme-preview {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.theme-colors {
  width: 100%;
  height: 100%;
  display: grid;
  grid-template-columns: 1fr 1fr;
  grid-template-rows: 1fr 1fr;
}

.theme-color {
  width: 100%;
  height: 100%;
}

.check-mark {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.9);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  animation: scaleIn 0.2s ease;
}

@keyframes scaleIn {
  from {
    transform: translate(-50%, -50%) scale(0);
  }
  to {
    transform: translate(-50%, -50%) scale(1);
  }
}

.background-name,
.theme-name {
  font-size: 13px;
  font-weight: 500;
  color: #1f2937;
  margin-top: 8px;
  text-align: center;
}

.theme-description {
  font-size: 12px;
  color: #6b7280;
  text-align: center;
  margin-top: 4px;
}

.custom-background {
  margin-top: 16px;
}

.upload-area {
  display: flex;
  align-items: center;
  gap: 12px;
}

.hidden-input {
  display: none;
}

.upload-btn {
  flex: 1;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  border: 2px dashed rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #6b7280;
}

.upload-btn:hover {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.05);
  color: #3b82f6;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-item label {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

.form-select,
.form-input {
  width: 100%;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  font-size: 13px;
  color: #1f2937;
  transition: all 0.2s ease;
}

.form-select:focus,
.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox {
  width: 18px;
  height: 18px;
  accent-color: #3b82f6;
}

.checkbox-label span {
  font-size: 13px;
  color: #374151;
}

.mode-options {
  display: flex;
  gap: 12px;
}

.mode-btn {
  flex: 1;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  color: #6b7280;
}

.mode-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  transform: translateY(-2px);
}

.mode-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.size-slider {
  display: flex;
  align-items: center;
  gap: 12px;
}

.slider {
  flex: 1;
  height: 4px;
  accent-color: #3b82f6;
}

.size-value {
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  min-width: 40px;
}

.manager-footer {
  display: flex;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
}

.btn {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary {
  background: rgba(0, 0, 0, 0.05);
  color: #6b7280;
}

.btn-secondary:hover {
  background: rgba(0, 0, 0, 0.1);
  color: #1f2937;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .theme-manager {
    background: rgba(0, 0, 0, 0.85);
  }

  .manager-header h2,
  .content-section h3,
  .content-section h4,
  .setting-item label,
  .checkbox-label span {
    color: #f9fafb;
  }

  .background-name,
  .theme-name {
    color: #f9fafb;
  }

  .theme-description,
  .size-value {
    color: #9ca3af;
  }

  .upload-btn,
  .form-select,
  .form-input {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
    color: #f9fafb;
  }

  .mode-btn {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.05);
    color: #9ca3af;
  }
}
</style>
