<template>
  <div class="control-panel-app">
    <ControlPanel />
  </div>
</template>

<script setup lang="ts">
import ControlPanel from '../components/ControlPanel/ControlPanel.vue'
</script>

<style scoped>
.control-panel-app {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
</style>