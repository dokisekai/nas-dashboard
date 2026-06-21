<template>
  <BackupManager />
</template>

<script setup lang="ts">
// SyncManager 与 BackupManager 共享同一套 Restic 后端能力
// （仓库 / 任务 / 快照 / 仓库间同步 / 设置）。
// 历史上这里有一份并行的实现，但调用的是后端的 stub API；
// 现在统一收敛到 BackupManager.vue 以避免功能割裂。
import BackupManager from './BackupManager.vue'
</script>
