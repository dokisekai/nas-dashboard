<template>
  <div class="quick-note-widget" :class="`widget-${size}`">
    <div class="widget-header">
      <div class="widget-title">便签</div>
      <div class="widget-controls">
        <button @click="clearNote" class="control-btn" title="清除">
          <TrashIcon class="w-4 h-4" />
        </button>
      </div>
    </div>

    <div class="note-content">
      <textarea
        v-model="note"
        @input="saveNote"
        placeholder="输入您的便签..."
        class="note-textarea"
      ></textarea>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { TrashIcon } from '@heroicons/vue/24/outline'

interface Props {
  config: {
    storageKey?: string
  }
  size: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  config: () => ({
    storageKey: 'quick-note'
  })
})

const note = ref('')
const storageKey = `widget-${props.config.storageKey}`

const saveNote = () => {
  localStorage.setItem(storageKey, note.value)
}

const loadNote = () => {
  const saved = localStorage.getItem(storageKey)
  if (saved) {
    note.value = saved
  }
}

const clearNote = () => {
  note.value = ''
  saveNote()
}

onMounted(() => {
  loadNote()
})
</script>

<style scoped>
.quick-note-widget {
  width: 100%;
  height: 100%;
  padding: 16px;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #fef9c3 0%, #fef08a 100%);
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
  color: #854d0e;
}

.widget-controls {
  display: flex;
  gap: 4px;
}

.control-btn {
  width: 24px;
  height: 24px;
  padding: 0;
  background: rgba(255, 255, 255, 0.8);
  border: none;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #854d0e;
  transition: all 0.2s ease;
}

.control-btn:hover {
  background: white;
}

.note-content {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.note-textarea {
  flex: 1;
  width: 100%;
  padding: 12px;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 8px;
  resize: none;
  font-size: 13px;
  line-height: 1.5;
  color: #1f2937;
  font-family: inherit;
}

.note-textarea::placeholder {
  color: #9ca3af;
}

.note-textarea:focus {
  outline: none;
  box-shadow: 0 0 0 2px rgba(234, 179, 8, 0.5);
}
</style>
