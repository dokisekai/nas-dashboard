<template>
  <Transition name="fade">
    <div v-if="isOpen" class="global-search-overlay" @click.self="$emit('close')">
      <div class="search-container">
        <div class="search-input-wrapper">
          <MagnifyingGlassIcon class="search-icon" />
          <input
            ref="inputRef"
            v-model="query"
            type="text"
            placeholder="搜索应用、文件、设置..."
            @keyup.esc="$emit('close')"
            @keydown.down="moveSelection(1)"
            @keydown.up="moveSelection(-1)"
            @keydown.enter="executeSelected"
          />
          <div class="search-hint">ESC 退出</div>
        </div>

        <div v-if="results.length > 0" class="search-results">
          <div
            v-for="(result, index) in results"
            :key="result.id"
            class="result-item"
            :class="{ active: selectedIndex === index }"
            @click="executeResult(result)"
            @mouseenter="selectedIndex = index"
          >
            <div class="result-icon" :style="{ background: result.color }">
              <component :is="result.icon" class="w-5 h-5" />
            </div>
            <div class="result-info">
              <div class="result-title">{{ result.title }}</div>
              <div class="result-desc">{{ result.description }}</div>
            </div>
            <div class="result-category">{{ result.category }}</div>
          </div>
        </div>
        <div v-else-if="query" class="no-results">
          未找到与 "{{ query }}" 相关的结果
        </div>
        <div v-else class="search-initial">
          <div class="section-title">建议</div>
          <div class="suggestions">
            <div
              v-for="s in suggestions"
              :key="s.id"
              class="suggestion-chip"
              @click="query = s.title"
            >
              {{ s.title }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import {
  MagnifyingGlassIcon,
  FolderIcon,
  ServerIcon,
  CogIcon,
  UserIcon,
  CubeIcon,
  ClockIcon
} from '@heroicons/vue/24/outline'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits(['close', 'open-app'])

const query = ref('')
const selectedIndex = ref(0)
const inputRef = ref<HTMLInputElement | null>(null)

const apps = [
  { id: 'storage-manager', title: '存储管理', description: '管理磁盘和共享', category: '系统', icon: FolderIcon, color: '#f59e0b' },
  { id: 'system-monitor', title: '系统监控', description: '实时监控系统资源', category: '监控', icon: ServerIcon, color: '#3b82f6' },
  { id: 'app-center', title: '应用中心', description: '管理系统套件', category: '应用', icon: CubeIcon, color: '#10b981' },
  { id: 'user-manager', title: '用户管理', description: '权限与账户设置', category: '管理', icon: UserIcon, color: '#ef4444' },
  { id: 'settings', title: '系统设置', description: '个性化与系统配置', category: '系统', icon: CogIcon, color: '#6366f1' }
]

const suggestions = [
  { id: 1, title: '时间机器' },
  { id: 2, title: '磁盘状态' },
  { id: 3, title: '添加用户' },
  { id: 4, title: 'Docker' }
]

const results = computed(() => {
  if (!query.value) return []
  const q = query.value.toLowerCase()
  return apps.filter(app =>
    app.title.toLowerCase().includes(q) ||
    app.description.toLowerCase().includes(q) ||
    app.category.toLowerCase().includes(q)
  )
})

watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    query.value = ''
    selectedIndex.value = 0
    nextTick(() => {
      inputRef.value?.focus()
    })
  }
})

const moveSelection = (dir: number) => {
  selectedIndex.value = (selectedIndex.value + dir + results.value.length) % results.value.length
}

const executeSelected = () => {
  if (results.value[selectedIndex.value]) {
    executeResult(results.value[selectedIndex.value])
  }
}

const executeResult = (result: any) => {
  emit('open-app', result.id, result.title)
  emit('close')
}
</script>

<style scoped>
.global-search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(4px);
  z-index: 3000;
  display: flex;
  justify-content: center;
  padding-top: 100px;
}

.search-container {
  width: 100%;
  max-width: 600px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  overflow: hidden;
  height: fit-content;
  max-height: 500px;
  display: flex;
  flex-direction: column;
}

.search-input-wrapper {
  padding: 16px 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.search-icon {
  width: 24px;
  height: 24px;
  color: #9ca3af;
}

.search-input-wrapper input {
  flex: 1;
  background: transparent;
  border: none;
  font-size: 18px;
  color: #1f2937;
  outline: none;
}

.search-hint {
  font-size: 11px;
  color: #9ca3af;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 4px;
}

.search-results {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.result-item.active {
  background: #eff6ff;
}

.result-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.result-info {
  flex: 1;
}

.result-title {
  font-size: 15px;
  font-weight: 600;
  color: #1f2937;
}

.result-desc {
  font-size: 13px;
  color: #6b7280;
}

.result-category {
  font-size: 12px;
  color: #9ca3af;
  background: #f3f4f6;
  padding: 2px 8px;
  border-radius: 4px;
}

.no-results {
  padding: 40px;
  text-align: center;
  color: #9ca3af;
}

.search-initial {
  padding: 20px;
}

.section-title {
  font-size: 12px;
  font-weight: 600;
  color: #9ca3af;
  text-transform: uppercase;
  margin-bottom: 12px;
}

.suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.suggestion-chip {
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 20px;
  font-size: 13px;
  color: #4b5563;
  cursor: pointer;
  transition: all 0.2s;
}

.suggestion-chip:hover {
  background: #e5e7eb;
  color: #1f2937;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}
</style>
