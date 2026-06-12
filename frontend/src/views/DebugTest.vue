<template>
  <div class="p-8 bg-white text-black">
    <h1 class="text-2xl font-bold mb-4">Debug Test Page</h1>
    <div class="space-y-2">
      <p>Route: {{ $route.name }}</p>
      <p>Path: {{ $route.path }}</p>
      <p>This is a simple test page to verify layout is working correctly.</p>
      <div class="mt-4 p-4 bg-gray-200">
        <p>If you can see this page without any duplicate menus or layout issues,</p>
        <p>then the basic layout structure is working correctly.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'

onMounted(() => {
  console.log('DEBUG DebugTest - Component mounted')
  console.log('DEBUG DebugTest - Route:', window.location.pathname)
})
</script>