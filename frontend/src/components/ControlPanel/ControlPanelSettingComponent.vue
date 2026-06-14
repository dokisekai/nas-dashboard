<template>
  <div class="cps-setting" :class="{ 'cps-advanced': setting.advanced }" :data-setting-type="setting.type">

    <!-- 自定义组件类型 - 不显示头部 - 支持NetworkSettingsPanel -->
    <template v-if="setting.type === 'custom'">
      <component
        :is="getCustomComponent(setting.component)"
        :value="localValue"
        :embedded-mode="true"
        @update="handleCustomUpdate"
        v-if="getCustomComponent(setting.component)"
      />
    </template>

    <!-- 其他类型的设置 -->
    <template v-else>
      <!-- 设置头部 -->
      <div class="cps-header">
        <div class="cps-info">
          <label class="cps-label">
            {{ setting.label }}
            <span v-if="setting.restartRequired" class="cps-badge cps-badge-warning">
              需要重启
            </span>
            <span v-if="setting.advanced" class="cps-badge cps-badge-advanced">
              高级
            </span>
          </label>
          <p v-if="setting.description" class="cps-description">
            {{ setting.description }}
          </p>
        </div>

        <!-- 设置控件 -->
        <div class="cps-control">
        <!-- 布尔类型 -->
        <template v-if="setting.type === 'boolean'">
          <button
            :class="['cps-switch', { active: localValue }]"
            @click="toggleValue"
          >
            <div class="cps-switch-slider"></div>
          </button>
        </template>

        <!-- 字符串类型 -->
        <template v-else-if="setting.type === 'string'">
          <input
            v-model="localValue"
            type="text"
            :placeholder="setting.description || ''"
            @change="handleChange"
            class="cps-input"
          />
        </template>

        <!-- 密码类型 -->
        <template v-else-if="setting.type === 'password'">
          <div class="cps-password-input">
            <input
              v-model="localValue"
              :type="showPassword ? 'text' : 'password'"
              @change="handleChange"
              class="cps-input"
            />
            <button
              @click="showPassword = !showPassword"
              class="cps-btn cps-btn-ghost"
            >
              <EyeIcon v-if="!showPassword" class="w-4 h-4" />
              <EyeSlashIcon v-else class="w-4 h-4" />
            </button>
          </div>
        </template>

        <!-- 数字类型 -->
        <template v-else-if="setting.type === 'number'">
          <input
            v-model.number="localValue"
            type="number"
            @change="handleChange"
            class="cps-input cps-input-number"
          />
        </template>

        <!-- 选择类型 -->
        <template v-else-if="setting.type === 'select'">
          <select
            v-model="localValue"
            @change="handleChange"
            class="cps-select"
          >
            <option
              v-for="option in setting.options"
              :key="option.value"
              :value="option.value"
              :disabled="option.disabled"
            >
              {{ option.label }}
            </option>
          </select>
        </template>

        <!-- 多选类型 -->
        <template v-else-if="setting.type === 'multiselect'">
          <div class="cps-multiselect">
            <button
              v-for="option in setting.options"
              :key="option.value"
              :class="['cps-option', { active: isMultiSelected(option.value) }]"
              @click="toggleMultiSelect(option.value)"
            >
              {{ option.label }}
            </button>
          </div>
        </template>

        <!-- 滑块类型 -->
        <template v-else-if="setting.type === 'slider'">
          <div class="cps-slider-container">
            <input
              v-model.number="localValue"
              type="range"
              :min="setting.min !== undefined ? setting.min : 0"
              :max="setting.max !== undefined ? setting.max : 100"
              @input="handleChange"
              class="cps-slider"
            />
            <span class="cps-slider-value">{{ localValue }}</span>
          </div>
        </template>

        <!-- 颜色类型 -->
        <template v-else-if="setting.type === 'color'">
          <div class="cps-color-input">
            <input
              v-model="localValue"
              type="color"
              @change="handleChange"
              class="cps-color-picker"
            />
            <span class="cps-color-value">{{ localValue }}</span>
          </div>
        </template>

        <!-- 文本域类型 -->
        <template v-else-if="setting.type === 'textarea'">
          <textarea
            v-model="localValue"
            @change="handleChange"
            :placeholder="setting.description || ''"
            rows="3"
            class="cps-textarea"
          ></textarea>
        </template>

        <!-- 文件类型 -->
        <template v-else-if="setting.type === 'file'">
          <div class="cps-file-input">
            <input
              ref="fileInput"
              type="file"
              @change="handleFileChange"
              class="cps-file-native"
              style="display: none"
            />
            <button
              @click="fileInput?.click()"
              class="cps-btn cps-btn-secondary"
            >
              <FolderIcon class="w-4 h-4" />
              选择文件
            </button>
            <span v-if="localValue" class="cps-file-name">{{ localValue }}</span>
          </div>
        </template>

        <!-- 只读类型 -->
        <template v-else-if="setting.type === 'readonly' || setting.readonly">
          <div class="cps-readonly">
            {{ localValue }}
          </div>
        </template>
      </div>
    </div>

    <!-- 验证错误 -->
    <div v-if="validationError" class="cps-error">
      <ExclamationCircleIcon class="w-4 h-4" />
      <span>{{ validationError }}</span>
    </div>

    <!-- 验证警告 -->
    <div v-if="validationWarning" class="cps-warning">
      <ExclamationTriangleIcon class="w-4 h-4" />
      <span>{{ validationWarning }}</span>
    </div>

    <!-- 依赖项说明 -->
    <div v-if="hasDependencies" class="cps-dependencies">
      <InformationCircleIcon class="w-4 h-4" />
      <span>此设置依赖于: {{ dependencyLabels }}</span>
    </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, defineAsyncComponent } from 'vue'
import {
  EyeIcon,
  EyeSlashIcon,
  ExclamationCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  FolderIcon
} from '@heroicons/vue/24/outline'
import type { ControlPanelSetting, ValidationResult } from './ControlPanelTypes'

interface Props {
  setting: ControlPanelSetting
  value: any
}

interface Emits {
  (e: 'update', value: any): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 状态
const localValue = ref(props.value)
const showPassword = ref(false)
const validationError = ref<string>()
const validationWarning = ref<string>()
const fileInput = ref<HTMLInputElement>()

// 计算属性
const hasDependencies = computed(() => {
  return props.setting.dependencies && props.setting.dependencies.length > 0
})

const dependencyLabels = computed(() => {
  if (!hasDependencies.value) return ''

  return props.setting.dependencies!
    .map(dep => {
      const parts = dep.split('.')
      return parts[parts.length - 1]
    })
    .join(', ')
})

// 监听外部值变化
watch(() => props.value, (newValue) => {
  localValue.value = newValue
})

// 方法
const toggleValue = () => {
  localValue.value = !localValue.value
  handleChange()
}

const isMultiSelected = (value: any) => {
  if (!Array.isArray(localValue.value)) return false
  return localValue.value.includes(value)
}

const toggleMultiSelect = (value: any) => {
  if (!Array.isArray(localValue.value)) {
    localValue.value = []
  }

  const index = localValue.value.indexOf(value)
  if (index > -1) {
    localValue.value.splice(index, 1)
  } else {
    localValue.value.push(value)
  }

  handleChange()
}

const handleChange = () => {
  // 验证值
  validate()

  // 发送更新事件
  emit('update', localValue.value)
}

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    localValue.value = file.name
    handleChange()
  }
}

const getCustomComponent = (componentName?: string) => {
  if (!componentName) return null

  // 组件映射 - 使用defineAsyncComponent正确处理异步组件
  const componentMap: Record<string, any> = {
    'ServicesListComponent': defineAsyncComponent(() =>
      import('./ServicesListComponent.vue')
    ),
    'NetworkInterfacesComponent': defineAsyncComponent(() =>
      import('./NetworkInterfacesComponent.vue')
    ),
    'WiFiScanComponent': defineAsyncComponent(() =>
      import('./WiFiScanComponent.vue')
    ),
    'NetworkManager': defineAsyncComponent(() =>
      import('../../apps/NetworkManager.vue')
    ),
    'NetworkSettingsPanel': defineAsyncComponent(() =>
      import('./NetworkSettingsPanel.vue')
    ),
    'SystemInfoPanel': defineAsyncComponent(() =>
      import('./SystemInfoPanel.vue')
    ),
    'UserManager': defineAsyncComponent(() =>
      import('./UserManager.vue')
    ),
    'GroupManager': defineAsyncComponent(() =>
      import('./GroupManager.vue')
    )
  }

  return componentMap[componentName] || null
}

const handleCustomUpdate = (value: any) => {
  localValue.value = value
  handleChange()
}

const validate = () => {
  validationError.value = undefined
  validationWarning.value = undefined

  if (!props.setting.validation) {
    return
  }

  const result: ValidationResult = props.setting.validation(localValue.value)

  if (!result.valid) {
    validationError.value = result.error
  } else if (result.warning) {
    validationWarning.value = result.warning
  }
}

// 初始验证
validate()
</script>

<style scoped>
.cps-setting {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cps-setting.cps-advanced {
  border-left: 3px solid #f59e0b;
  padding-left: 16px;
}

.cps-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
}

.cps-info {
  flex: 1;
}

.cps-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 4px;
}

.cps-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

.cps-badge-warning {
  background: #fef3c7;
  color: #92400e;
}

.cps-badge-advanced {
  background: #dbeafe;
  color: #1e40af;
}

.cps-description {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.5;
}

/* 自定义组件样式 */
.cps-setting:has(.network-manager) {
  width: 100%;
  padding: 0;
  margin: 0;
}

.cps-setting:has(.network-manager) .network-manager {
  width: 100%;
  height: auto;
  min-height: 400px;
}

.cps-control {
  display: flex;
  align-items: center;
}

/* 开关控件 */
.cps-switch {
  width: 48px;
  height: 24px;
  background: #e5e7eb;
  border-radius: 12px;
  position: relative;
  cursor: pointer;
  transition: background 0.2s ease;
  border: none;
  padding: 0;
}

.cps-switch.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.cps-switch-slider {
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform 0.2s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.cps-switch.active .cps-switch-slider {
  transform: translateX(24px);
}

/* 输入控件 */
.cps-input {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  min-width: 200px;
  transition: border-color 0.2s ease;
}

.cps-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.cps-input-number {
  min-width: 100px;
}

/* 密码输入 */
.cps-password-input {
  display: flex;
  gap: 8px;
}

.cps-password-input .cps-input {
  flex: 1;
}

/* 选择控件 */
.cps-select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  min-width: 200px;
  background: white;
  cursor: pointer;
}

.cps-select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* 多选控件 */
.cps-multiselect {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.cps-option {
  padding: 6px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 12px;
  background: white;
  cursor: pointer;
  transition: all 0.2s ease;
}

.cps-option:hover {
  border-color: #667eea;
  background: rgba(102, 126, 234, 0.05);
}

.cps-option.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
}

/* 滑块控件 */
.cps-slider-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cps-slider {
  width: 200px;
  height: 6px;
  -webkit-appearance: none;
  background: #e5e7eb;
  border-radius: 3px;
  outline: none;
}

.cps-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 16px;
  height: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  cursor: pointer;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.cps-slider-value {
  font-size: 12px;
  color: #6b7280;
  min-width: 40px;
  text-align: right;
}

/* 颜色控件 */
.cps-color-input {
  display: flex;
  align-items: center;
  gap: 8px;
}

.cps-color-picker {
  width: 40px;
  height: 32px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  cursor: pointer;
}

.cps-color-value {
  font-size: 12px;
  color: #6b7280;
  font-family: monospace;
}

/* 文本域 */
.cps-textarea {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  font-size: 14px;
  min-width: 300px;
  resize: vertical;
  font-family: inherit;
}

.cps-textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

/* 文件输入 */
.cps-file-input {
  display: flex;
  align-items: center;
  gap: 12px;
}

.cps-file-name {
  font-size: 12px;
  color: #6b7280;
}

/* 只读样式 */
.cps-readonly {
  padding: 8px 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 6px;
  color: #6b7280;
  font-size: 14px;
  font-family: monospace;
  min-width: 200px;
}

/* 按钮样式 */
.cps-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  border: none;
}

.cps-btn-secondary {
  background: white;
  border: 1px solid #e5e7eb;
  color: #6b7280;
}

.cps-btn-secondary:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.cps-btn-ghost {
  background: transparent;
  border: none;
  color: #6b7280;
  padding: 4px;
}

.cps-btn-ghost:hover {
  background: #f3f4f6;
  color: #1f2937;
}

/* 错误和警告 */
.cps-error {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 6px;
  color: #991b1b;
  font-size: 12px;
}

.cps-warning {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #fffbeb;
  border: 1px solid #fde68a;
  border-radius: 6px;
  color: #92400e;
  font-size: 12px;
}

/* 依赖项说明 */
.cps-dependencies {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #eff6ff;
  border: 1px solid #dbeafe;
  border-radius: 6px;
  color: #1e40af;
  font-size: 11px;
}
</style>