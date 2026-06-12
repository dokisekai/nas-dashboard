<template>
  <div class="widget-config">
    <div class="config-header">
      <h3>小部件设置</h3>
      <button @click="$emit('close')" class="close-btn">
        <XMarkIcon class="w-5 h-5" />
      </button>
    </div>

    <div class="config-content">
      <div class="config-section">
        <h4>基本信息</h4>
        <div class="form-group">
          <label>小部件名称</label>
          <input v-model="localConfig.name" type="text" class="form-input" />
        </div>
        <div class="form-group">
          <label>尺寸</label>
          <div class="size-selector">
            <button
              v-for="size in sizes"
              :key="size.value"
              class="size-btn"
              :class="{ active: localConfig.size === size.value }"
              @click="localConfig.size = size.value"
            >
              {{ size.label }}
            </button>
          </div>
        </div>
      </div>

      <div class="config-section">
        <h4>显示选项</h4>
        <div v-for="option in displayOptions" :key="option.key" class="form-group">
          <label class="checkbox-label">
            <input
              v-model="localConfig.config[option.key]"
              type="checkbox"
              class="checkbox"
            />
            <span>{{ option.label }}</span>
          </label>
        </div>
      </div>

      <div class="config-section" v-if="widgetSpecificOptions.length > 0">
        <h4>高级选项</h4>
        <div v-for="option in widgetSpecificOptions" :key="option.key" class="form-group">
          <label>{{ option.label }}</label>
          <input
            v-if="option.type === 'text'"
            v-model="localConfig.config[option.key]"
            type="text"
            class="form-input"
            :placeholder="option.placeholder"
          />
          <input
            v-else-if="option.type === 'number'"
            v-model.number="localConfig.config[option.key]"
            type="number"
            class="form-input"
            :min="option.min"
            :max="option.max"
          />
          <select
            v-else-if="option.type === 'select'"
            v-model="localConfig.config[option.key]"
            class="form-input"
          >
            <option v-for="choice in option.choices" :key="choice.value" :value="choice.value">
              {{ choice.label }}
            </option>
          </select>
        </div>
      </div>
    </div>

    <div class="config-footer">
      <button @click="$emit('close')" class="btn btn-secondary">取消</button>
      <button @click="saveConfig" class="btn btn-primary">保存</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'
import type { Widget } from '../../types/desktop'

interface Props {
  widget: Widget
}

interface Emits {
  (e: 'close'): void
  (e: 'save', config: Widget): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const localConfig = ref({
  name: '',
  size: 'medium',
  config: { ...props.widget.config }
})

const sizes = [
  { value: 'small', label: '小' },
  { value: 'medium', label: '中' },
  { value: 'large', label: '大' }
]

const displayOptions = computed(() => {
  const options = []

  if (props.widget.type === 'clock') {
    options.push(
      { key: 'showDate', label: '显示日期' },
      { key: 'showSeconds', label: '显示秒数' },
      { key: 'format24', label: '24小时制' }
    )
  }

  if (props.widget.type === 'system-monitor') {
    options.push(
      { key: 'showCpu', label: '显示CPU' },
      { key: 'showMemory', label: '显示内存' },
      { key: 'showDisk', label: '显示磁盘' }
    )
  }

  if (props.widget.type === 'weather') {
    options.push(
      { key: 'showHumidity', label: '显示湿度' },
      { key: 'showWind', label: '显示风速' }
    )
  }

  return options
})

const widgetSpecificOptions = computed(() => {
  const options = []

  if (props.widget.type === 'weather') {
    options.push({
      key: 'location',
      label: '位置',
      type: 'text',
      placeholder: '输入城市名称'
    })
    options.push({
      key: 'units',
      label: '温度单位',
      type: 'select',
      choices: [
        { value: 'celsius', label: '摄氏度' },
        { value: 'fahrenheit', label: '华氏度' }
      ]
    })
  }

  if (props.widget.type === 'storage-status') {
    options.push({
      key: 'volume',
      label: '监控卷',
      type: 'text',
      placeholder: '/或/mnt/data'
    })
  }

  return options
})

const saveConfig = () => {
  const updatedWidget: Widget = {
    ...props.widget,
    size: localConfig.value.size as 'small' | 'medium' | 'large',
    config: localConfig.value.config
  }
  emit('save', updatedWidget)
}

// 初始化配置
localConfig.value.size = props.widget.size
</script>

<style scoped>
.widget-config {
  width: 100%;
  max-width: 480px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.config-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.close-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #6b7280;
  transition: all 0.2s ease;
}

.close-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.config-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.config-section {
  margin-bottom: 24px;
}

.config-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 8px;
}

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

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.size-selector {
  display: flex;
  gap: 8px;
}

.size-btn {
  flex: 1;
  padding: 10px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  font-size: 13px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.size-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #1f2937;
}

.size-btn.active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
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

.config-footer {
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
  .widget-config {
    background: rgba(0, 0, 0, 0.85);
  }

  .config-header h3,
  .config-section h4,
  .form-group label,
  .checkbox-label span {
    color: #f9fafb;
  }

  .form-input {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
    color: #f9fafb;
  }

  .size-btn {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
    color: #9ca3af;
  }

  .size-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #f9fafb;
  }

  .btn-secondary {
    background: rgba(255, 255, 255, 0.05);
    color: #9ca3af;
  }

  .btn-secondary:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #f9fafb;
  }
}
</style>
