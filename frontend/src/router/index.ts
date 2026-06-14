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
    {
      path: '/monitor/power',
      name: 'MonitorPower',
      component: () => import('../views/Monitor/Power.vue'),
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

// 路由守卫 - 简化认证检查
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const loggedIn = authStore.isLoggedIn()

  console.log('Route guard - navigating to:', to.path, 'name:', to.name)
  console.log('Is logged in:', loggedIn)
  console.log('Token in localStorage:', localStorage.getItem('token')?.substring(0, 20) + '...')

  // 登录页面：如果已登录则跳转到桌面
  if (to.path === '/login') {
    if (loggedIn) {
      console.log('Already logged in, redirecting to desktop')
      return next('/desktop')
    }
    console.log('Showing login page')
    return next()
  }

  // 初始化页面：直接放行
  if (to.path === '/initialization') {
    console.log('Showing initialization page')
    return next()
  }

  // 所有其他页面：需要登录
  if (!loggedIn) {
    console.log('Not logged in, redirecting to login')
    return next('/login')
  }

  // 已登录用户：允许访问
  console.log('Logged in, allowing access to', to.path)
  next()
})

export default router
