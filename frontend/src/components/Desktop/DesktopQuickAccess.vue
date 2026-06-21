<template>
  <div class="quick-access">
    <div
      v-for="app in apps"
      :key="app.id"
      class="quick-icon"
      @click="$emit('open-app', app.id)"
    >
      <div class="icon-bg">
        <component :is="app.icon" class="w-8 h-8" />
      </div>
      <span>{{ app.title }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AppDefinition } from '../../config/apps'

defineProps<{
  apps: AppDefinition[]
}>()

defineEmits<{
  (e: 'open-app', appId: string): void
}>()
</script>

<style scoped>
.quick-access {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  gap: 32px;
  z-index: 1;
}

.quick-icon {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.quick-icon:hover {
  transform: scale(1.1);
}

.icon-bg {
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
}

.quick-icon span {
  color: white;
  font-size: 14px;
  font-weight: 500;
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}
</style>
