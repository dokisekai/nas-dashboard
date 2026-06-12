<template>
  <div class="widget-library">
    <div class="library-header">
      <h2>小部件库</h2>
      <div class="library-tabs">
        <button
          v-for="category in categories"
          :key="category.id"
          class="tab-btn"
          :class="{ active: activeCategory === category.id }"
          @click="activeCategory = category.id"
        >
          {{ category.name }}
        </button>
      </div>
    </div>

    <div class="library-content">
      <div class="widget-grid">
        <div
          v-for="widget in filteredWidgets"
          :key="widget.id"
          class="widget-card"
          @click="selectWidget(widget)"
        >
          <div class="widget-preview">
            <component :is="getWidgetIcon(widget.type)" class="w-8 h-8" />
          </div>
          <div class="widget-info">
            <div class="widget-name">{{ widget.name }}</div>
            <div class="widget-description">{{ widget.description }}</div>
          </div>
          <button class="add-btn" @click.stop="addWidget(widget)">
            <PlusIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { PlusIcon } from '@heroicons/vue/24/outline'

interface WidgetTemplate {
  id: string
  type: string
  name: string
  description: string
  category: string
  defaultSize: 'small' | 'medium' | 'large'
  icon: string
}

interface Emits {
  (e: 'add-widget', widget: WidgetTemplate): void
}

const emit = defineEmits<Emits>()

const activeCategory = ref('all')

const categories = [
  { id: 'all', name: '全部' },
  { id: 'system', name: '系统' },
  { id: 'information', name: '信息' },
  { id: 'productivity', name: '效率' },
  { id: 'media', name: '媒体' }
]

const widgetTemplates: WidgetTemplate[] = [
  {
    id: 'system-monitor',
    type: 'system-monitor',
    name: '系统监控',
    description: '显示CPU、内存和磁盘使用率',
    category: 'system',
    defaultSize: 'large',
    icon: 'ChartBarIcon'
  },
  {
    id: 'storage-status',
    type: 'storage-status',
    name: '存储状态',
    description: '显示存储空间使用情况',
    category: 'system',
    defaultSize: 'medium',
    icon: 'ServerIcon'
  },
  {
    id: 'network-monitor',
    type: 'network-monitor',
    name: '网络监控',
    description: '显示网络流量和连接状态',
    category: 'system',
    defaultSize: 'small',
    icon: 'WifiIcon'
  },
  {
    id: 'clock',
    type: 'clock',
    name: '时钟',
    description: '显示当前时间和日期',
    category: 'information',
    defaultSize: 'small',
    icon: 'ClockIcon'
  },
  {
    id: 'calendar',
    type: 'calendar',
    name: '日历',
    description: '显示当前月份和事件',
    category: 'information',
    defaultSize: 'medium',
    icon: 'CalendarIcon'
  },
  {
    id: 'weather',
    type: 'weather',
    name: '天气',
    description: '显示当前天气状况',
    category: 'information',
    defaultSize: 'medium',
    icon: 'CloudIcon'
  },
  {
    id: 'quick-note',
    type: 'quick-note',
    name: '便签',
    description: '快速记录笔记和提醒',
    category: 'productivity',
    defaultSize: 'medium',
    icon: 'DocumentTextIcon'
  }
]

const filteredWidgets = computed(() => {
  if (activeCategory.value === 'all') {
    return widgetTemplates
  }
  return widgetTemplates.filter(w => w.category === activeCategory.value)
})

const getWidgetIcon = (type: string) => {
  // 返回对应图标组件
  return 'div'
}

const selectWidget = (widget: WidgetTemplate) => {
  // 显示小部件详情
}

const addWidget = (widget: WidgetTemplate) => {
  emit('add-widget', widget)
}
</script>

<style scoped>
.widget-library {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  overflow: hidden;
}

.library-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.library-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.library-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.tab-btn {
  padding: 6px 12px;
  background: transparent;
  border: none;
  border-radius: 6px;
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

.library-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.widget-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.widget-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(0, 0, 0, 0.05);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.widget-card:hover {
  border-color: rgba(59, 130, 246, 0.3);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.widget-preview {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.widget-info {
  flex: 1;
}

.widget-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.widget-description {
  font-size: 12px;
  color: #6b7280;
}

.add-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  background: #3b82f6;
  border: none;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-btn:hover {
  background: #2563eb;
  transform: scale(1.1);
}

/* 深色模式 */
@media (prefers-color-scheme: dark) {
  .widget-library {
    background: rgba(0, 0, 0, 0.85);
  }

  .library-header h2 {
    color: #f9fafb;
  }

  .widget-card {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.1);
  }

  .widget-name {
    color: #f9fafb;
  }

  .widget-description {
    color: #9ca3af;
  }
}
</style>
