<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import MainLayout from './components/Layout/Main.vue'
import SimpleDesktop from './components/Desktop/SimpleDesktop.vue'

const route = useRoute()
// 登录页和桌面页不使用MainLayout
const isLogin = computed(() => route.name === 'Login')
const isDesktop = computed(() => route.name === 'Desktop')
</script>

<template>
  <!-- 登录页面：直接显示 -->
  <router-view v-if="isLogin" />

  <!-- 桌面页面：显示简化版桌面界面 -->
  <SimpleDesktop v-else-if="isDesktop" />

  <!-- 其他页面：使用MainLayout包装 -->
  <MainLayout v-else />
</template>
