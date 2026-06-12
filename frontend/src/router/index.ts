import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      redirect: '/desktop',
      meta: { requiresAuth: true },
    },
    // 桌面界面路由
    {
      path: '/desktop',
      name: 'Desktop',
      component: () => import('../components/Desktop/SimpleDesktop.vue'),
      meta: { requiresAuth: true },
    },
    // 传统仪表板（备用）
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/Dashboard.vue'),
      meta: { requiresAuth: true },
    },
    // 监控路由
    {
      path: '/monitor',
      name: 'Monitor',
      component: () => import('../views/Monitor/CPU.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/monitor/memory',
      name: 'MonitorMemory',
      component: () => import('../views/Monitor/Memory.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/monitor/disk',
      name: 'MonitorDisk',
      component: () => import('../views/Monitor/Disk.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/monitor/network',
      name: 'MonitorNetwork',
      component: () => import('../views/Monitor/Network.vue'),
      meta: { requiresAuth: true },
    },
    // 存储管理路由
    {
      path: '/storage',
      name: 'Storage',
      component: () => import('../views/Storage/Disks.vue'),
      meta: { requiresAuth: true },
    },
    // 服务管理路由
    {
      path: '/services',
      name: 'Services',
      component: () => import('../views/Services/System.vue'),
      meta: { requiresAuth: true },
    },
    // 用户管理路由
    {
      path: '/users',
      name: 'Users',
      component: () => import('../views/Users/Users.vue'),
      meta: { requiresAuth: true },
    },
    // 调试路由
    {
      path: '/debug-test',
      name: 'DebugTest',
      component: () => import('../views/DebugTest.vue'),
      meta: { requiresAuth: true },
    },
  ],
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isLoggedIn()) {
    next('/login')
  } else if (to.path === '/login' && authStore.isLoggedIn()) {
    next('/desktop')
  } else {
    next()
  }
})

export default router
