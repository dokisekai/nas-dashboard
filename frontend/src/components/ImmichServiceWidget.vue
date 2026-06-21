<template>
  <div class="immich-widget">
    <el-card class="immich-card">
      <template #header>
        <div class="card-header">
          <div class="icon-wrapper">
            <el-icon class="immich-icon"><Picture /></el-icon>
          </div>
          <div class="info-wrapper">
            <h3>Immich</h3>
            <el-tag :type="statusType" size="small">{{ statusText }}</el-tag>
          </div>
        </div>
      </template>

      <div class="card-content">
        <!-- 用户信息 -->
        <div v-if="user" class="user-info">
          <div class="user-header">
            <el-avatar>{{ user.name?.charAt(0) || 'U' }}</el-avatar>
            <div class="user-text">
              <div class="user-name">{{ user.name }}</div>
              <div class="user-email">{{ user.email }}</div>
            </div>
          </div>
          
          <el-divider />
          
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">{{ formatBytes(user.storageSize) }}</span>
              <span class="stat-label">存储空间</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ user.photos || 0 }}</span>
              <span class="stat-label">照片</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">{{ user.videos || 0 }}</span>
              <span class="stat-label">视频</span>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="actions">
          <el-button 
            type="primary" 
            @click="launchImmich"
            :disabled="!isOnline"
            style="width: 100%"
          >
            <el-icon><Monitor /></el-icon>
            打开管理界面
          </el-button>
        </div>

        <!-- 快速链接 -->
        <div class="quick-links">
          <el-button 
            text 
            @click="openUsers"
            :disabled="!isOnline"
            size="small"
          >
            <el-icon><User /></el-icon>
            用户管理
          </el-button>
          <el-button 
            text 
            @click="openAlbums"
            :disabled="!isOnline"
            size="small"
          >
            <el-icon><Collection /></el-icon>
            相册
          </el-button>
          <el-button 
            text 
            @click="refreshStatus"
            size="small"
          >
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Picture, Monitor, User, Collection, Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const IMMICH_URL = 'http://localhost:2283'
const API_KEY = 't5nmlDaFlyz7guxpfmAwvOeiKp7zEC7loF2Ow15V8Q'

const isOnline = ref(false)
const user = ref(null)
const statusType = ref('info')
const statusText = ref('检查中...')

let checkTimer = null

const checkStatus = async () => {
  try {
    const response = await fetch(`${IMMICH_URL}/api/server-info/about`)
    isOnline.value = response.ok
    
    if (isOnline.value) {
      statusType.value = 'success'
      statusText.value = '在线'
      await fetchUserInfo()
    } else {
      statusType.value = 'danger'
      statusText.value = '离线'
    }
  } catch (error) {
    isOnline.value = false
    statusType.value = 'danger'
    statusText.value = '离线'
  }
}

const fetchUserInfo = async () => {
  try {
    const response = await fetch(`${IMMICH_URL}/api/users`, {
      headers: { 'X-API-Key': API_KEY }
    })
    
    if (response.ok) {
      const users = await response.json()
      if (Array.isArray(users) && users.length > 0) {
        user.value = users[0]
      }
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
  }
}

const launchImmich = () => {
  if (isOnline.value) {
    window.open(IMMICH_URL, '_blank')
    ElMessage.success('正在打开Immich管理界面...')
  } else {
    ElMessage.error('Immich服务离线，请先启动服务')
  }
}

const openUsers = () => {
  if (isOnline.value) {
    window.open(`${IMMICH_URL}/admin/user-settings`, '_blank')
  }
}

const openAlbums = () => {
  if (isOnline.value) {
    window.open(`${IMMICH_URL}/albums`, '_blank')
  }
}

const refreshStatus = () => {
  checkStatus()
  ElMessage.success('状态已刷新')
}

const formatBytes = (bytes) => {
  if (!bytes) return '0 GB'
  const gb = bytes / (1024 * 1024 * 1024)
  return gb.toFixed(1) + ' GB'
}

onMounted(() => {
  checkStatus()
  // 每30秒检查一次状态
  checkTimer = setInterval(checkStatus, 30000)
})

onUnmounted(() => {
  if (checkTimer) {
    clearInterval(checkTimer)
  }
})
</script>

<style scoped>
.immich-widget {
  margin: 20px 0;
}

.immich-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 15px;
}

.icon-wrapper {
  width: 45px;
  height: 45px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.immich-icon {
  font-size: 24px;
}

.info-wrapper h3 {
  margin: 0;
  color: #303133;
  font-size: 16px;
}

.card-content {
  padding: 20px;
}

.user-info {
  margin-bottom: 20px;
}

.user-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.user-text {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-weight: 600;
  color: #303133;
  margin-bottom: 2px;
}

.user-email {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.stats {
  display: flex;
  justify-content: space-around;
  margin-bottom: 16px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 12px;
  color: #909399;
}

.actions {
  margin-bottom: 12px;
}

.quick-links {
  display: flex;
  justify-content: space-around;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}
</style>
