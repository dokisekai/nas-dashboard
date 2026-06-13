import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useSystemStore } from '../stores/system'

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
      path: '/initialization',
      name: 'Initialization',
      component: () => import('../views/Initialization.vue'),
      meta: { requiresAuth: false, checkInit: true },
    },
    {
      path: '/',
      name: 'Home',
      component: () => import('../components/Desktop/SimpleDesktop.vue'),
      meta: { requiresAuth: true },
    },
    // 桌面界面路由
    {
      path: '/desktop',
      name: 'Desktop',
      component: () => import('../components/Desktop/SimpleDesktop.vue'),
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
  ],
})

// 路由守卫 - 包含初始化检查
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const systemStore = useSystemStore()

  console.log('Route guard - navigating to:', to.path)
  console.log('Is dev mode:', authStore.isDevMode())
  console.log('Is logged in:', authStore.isLoggedIn())
  console.log('Token in store:', authStore.token?.substring(0, 20) + '...')
  console.log('Token in localStorage:', localStorage.getItem('token')?.substring(0, 20) + '...')

  // 跳过初始化页面本身的状态检查，避免循环
  if (to.path === '/initialization') {
    return next()
  }

  // 开发模式：直接跳过初始化检查和认证检查
  if (authStore.isDevMode()) {
    console.log('Dev mode: allowing access to', to.path)
    // 如果要去登录页面但已经有token，直接跳转到桌面
    if (to.path === '/login' && authStore.isLoggedIn()) {
      console.log('Dev mode: already logged in, redirecting to desktop')
      return next('/desktop')
    }
    // 其他页面直接允许访问
    return next()
  }

  // 生产模式：检查系统初始化状态
  try {
    await systemStore.checkInitStatus()
    if (!systemStore.initialized) {
      return next('/initialization')
    }
  } catch (error) {
    console.error('Failed to check initialization status:', error)
  }

  const loggedIn = authStore.isLoggedIn()

  // 处理登录页面的特殊逻辑
  if (to.path === '/login') {
    if (loggedIn) {
      return next('/desktop')
    }
    return next()
  }

  // 处理需要认证的页面
  if (to.meta.requiresAuth || to.path.startsWith('/apps/')) {
    if (loggedIn) {
      return next()
    }

    // 未登录跳转到登录页
    return next('/login')
  }

  next()
})

export default router
