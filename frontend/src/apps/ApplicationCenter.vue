<template>
  <div class="application-center">
    <div class="app-header">
      <h2 class="app-title">应用中心</h2>
      <div class="app-actions">
        <button class="btn-upload" @click="uploadPackage">
          <UploadIcon />
          <span>上传应用包</span>
        </button>
        <button class="btn-refresh" @click="refreshApps">
          <RefreshIcon />
          <span>刷新</span>
        </button>
      </div>
    </div>

    <!-- 应用分类 -->
    <div class="app-categories">
      <button
        v-for="category in categories"
        :key="category.id"
        :class="['category-btn', { active: selectedCategory === category.id }]"
        @click="selectedCategory = category.id"
      >
        <component :is="category.icon" />
        <span>{{ category.name }}</span>
      </button>
    </div>

    <!-- 应用列表 -->
    <div class="app-content">
      <div class="app-section">
        <h3 class="section-title">已安装应用</h3>
        <div v-if="installedApps.length === 0" class="empty-state">
          <EmptyIcon />
          <p>暂无已安装应用</p>
        </div>
        <div v-else class="app-grid">
          <AppCard
            v-for="app in installedApps"
            :key="app.id"
            :app="app"
            :installed="true"
            @start="handleStart"
            @stop="handleStop"
            @restart="handleRestart"
            @uninstall="handleUninstall"
            @settings="handleSettings"
          />
        </div>
      </div>

      <div class="app-section">
        <h3 class="section-title">可用应用</h3>
        <div v-if="availableApps.length === 0" class="empty-state">
          <EmptyIcon />
          <p>暂无可用应用</p>
        </div>
        <div v-else class="app-grid">
          <AppCard
            v-for="app in availableApps"
            :key="app.id"
            :app="app"
            :installed="false"
            @install="handleInstall"
          />
        </div>
      </div>
    </div>

    <!-- 上传对话框 -->
    <UploadDialog
      v-if="showUploadDialog"
      @upload="handleUpload"
      @cancel="showUploadDialog = false"
    />

    <!-- 安装进度对话框 -->
    <InstallProgressDialog
      v-if="installingApps.length > 0"
      :apps="installingApps"
      @close="installingApps = []"
    />

    <!-- 应用详情对话框 -->
    <AppDetailDialog
      v-if="selectedApp"
      :app="selectedApp"
      @install="handleInstall"
      @close="selectedApp = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { applicationApi, getInstallProgress, type AppInstance, type AppPackage } from '../api/application'
import AppCard from '../components/Application/AppCard.vue'
import UploadDialog from '../components/Application/UploadDialog.vue'
import InstallProgressDialog from '../components/Application/InstallProgressDialog.vue'
import AppDetailDialog from '../components/Application/AppDetailDialog.vue'
import {
  UploadIcon,
  RefreshIcon,
  EmptyIcon,
  CubeIcon,
  GridIcon,
  FilmIcon,
  DocumentIcon,
  WrenchIcon,
  ShieldIcon,
  GlobeIcon
} from '@heroicons/vue/24/outline'

// 应用分类
const categories = [
  { id: 'all', name: '全部', icon: GridIcon },
  { id: 'media', name: '媒体', icon: FilmIcon },
  { id: 'productivity', name: '办公', icon: DocumentIcon },
  { id: 'utilities', name: '工具', icon: WrenchIcon },
  { id: 'security', name: '安全', icon: ShieldIcon },
  { id: 'network', name: '网络', icon: GlobeIcon }
]

// 状态管理
const selectedCategory = ref('all')
const showUploadDialog = ref(false)
const selectedApp = ref<AppPackage | AppInstance | null>(null)

const instances = ref<AppInstance[]>([])
const packages = ref<AppPackage[]>([])
const installingApps = ref<AppInstance[]>([])

let refreshInterval: number

// 计算属性
const installedApps = computed(() => {
  const installed = instances.value.filter(instance => {
    if (selectedCategory.value === 'all') return true
    return instance.category === selectedCategory.value
  })
  return installed
})

const availableApps = computed(() => {
  const installedNames = new Set(instances.value.map(i => i.packageName))
  const available = packages.value.filter(pkg => {
    if (selectedCategory.value === 'all') return true
    return pkg.category === selectedCategory.value
  }).filter(pkg => !installedNames.has(pkg.name))
  return available
})

// 方法
const loadApps = async () => {
  try {
    const [instancesData, packagesData] = await Promise.all([
      applicationApi.listApps(),
      applicationApi.listPackages()
    ])

    instances.value = instancesData
    packages.value = packagesData
  } catch (error) {
    console.error('加载应用失败:', error)
  }
}

const refreshApps = () => {
  loadApps()
}

const uploadPackage = () => {
  showUploadDialog.value = true
}

const handleUpload = async (file: File) => {
  try {
    showUploadDialog.value = false
    const pkg = await applicationApi.uploadPackage(file)
    packages.value.push(pkg)
  } catch (error) {
    console.error('上传应用包失败:', error)
  }
}

const handleInstall = async (pkg: AppPackage) => {
  try {
    const instance = await applicationApi.installApp({
      packageName: pkg.name,
      version: pkg.version,
      autoStart: true
    })

    installingApps.value.push(instance)
    instances.value.push(instance)

    // 开始监听安装进度
    startProgressMonitoring(instance.id)
  } catch (error) {
    console.error('安装应用失败:', error)
  }
}

const handleStart = async (app: AppInstance) => {
  try {
    await applicationApi.startApp(app.id)
    app.status = 'running'
  } catch (error) {
    console.error('启动应用失败:', error)
  }
}

const handleStop = async (app: AppInstance) => {
  try {
    await applicationApi.stopApp(app.id)
    app.status = 'stopped'
  } catch (error) {
    console.error('停止应用失败:', error)
  }
}

const handleRestart = async (app: AppInstance) => {
  try {
    await applicationApi.restartApp(app.id)
  } catch (error) {
    console.error('重启应用失败:', error)
  }
}

const handleUninstall = async (app: AppInstance) => {
  if (!confirm(`确定要卸载应用 "${app.displayName}" 吗？`)) {
    return
  }

  try {
    await applicationApi.uninstallApp(app.id)
    instances.value = instances.value.filter(i => i.id !== app.id)
  } catch (error) {
    console.error('卸载应用失败:', error)
  }
}

const handleSettings = (app: AppInstance) => {
  selectedApp.value = app
}

const startProgressMonitoring = (instanceId: number) => {
  const eventSource = getInstallProgress(instanceId, (progress) => {
    const app = installingApps.value.find(a => a.id === instanceId)
    if (app) {
      // 更新进度
    }
  }, () => {
    // 安装完成
    const index = installingApps.value.findIndex(a => a.id === instanceId)
    if (index !== -1) {
      installingApps.value.splice(index, 1)
    }
    loadApps()
  })

  // 保存eventSource以便清理
  return eventSource
}

// 生命周期
onMounted(() => {
  loadApps()
  refreshInterval = setInterval(loadApps, 30000) as unknown as number
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.application-center {
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
}

.app-title {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.app-actions {
  display: flex;
  gap: 12px;
}

.btn-upload,
.btn-refresh {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-upload {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-upload:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.btn-refresh {
  background: rgba(255, 255, 255, 0.9);
  color: #667eea;
}

.btn-refresh:hover {
  background: rgba(255, 255, 255, 1);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.2);
}

.app-categories {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.15);
}

.category-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.category-btn:hover {
  background: rgba(102, 126, 234, 0.2);
  transform: translateY(-2px);
}

.category-btn.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.app-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.app-section {
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(102, 126, 234, 0.2);
}

.section-title {
  font-size: 20px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #9ca3af;
}

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: 12px;
}

.app-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}
</style>