<template>
  <div class="immich-integration">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <el-icon><Picture /></el-icon>
          <span>Immich 照片管理</span>
          <el-tag v-if="connected" type="success">已连接</el-tag>
          <el-tag v-else type="info">未连接</el-tag>
        </div>
      </template>

      <div class="card-content">
        <div v-if="loading" class="loading">
          <el-icon class="is-loading"><Loading /></el-icon>
          <p>正在连接Immich...</p>
        </div>

        <div v-else-if="error" class="error">
          <el-icon><Warning /></el-icon>
          <p>{{ error }}</p>
        </div>

        <div v-else class="content">
          <div class="user-info" v-if="currentUser">
            <el-avatar :size="64" :src="getUserAvatar()">
              {{ currentUser.name?.charAt(0) }}
            </el-avatar>
            <h3>{{ currentUser.name }}</h3>
            <p>{{ currentUser.email }}</p>
            <div class="stats">
              <el-statistic title="照片" :value="userStats.photos" />
              <el-statistic title="视频" :value="userStats.videos" />
              <el-statistic title="存储" :value="userStats.storage" suffix="GB" />
            </div>
          </div>

          <el-button type="primary" size="large" @click="openImmich" class="launch-btn">
            <el-icon><Monitor /></el-icon>
            打开Immich管理界面
          </el-button>

          <el-button @click="refreshUser" class="refresh-btn">
            <el-icon><Refresh /></el-icon>
            刷新用户信息
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Picture, Loading, Warning, Monitor, Refresh } from '@element-plus/icons-vue'
import axios from 'axios'

const loading = ref(true)
const connected = ref(false)
const error = ref('')
const currentUser = ref(null)
const userStats = ref({
  photos: 0,
  videos: 0,
  storage: 0
})

const checkImmichConnection = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const response = await axios.get('/api/services/immich/user')
    
    if (response.data.service === 'Immich') {
      connected.value = true
      currentUser.value = response.data.user
    } else {
      throw new Error('Immich服务未响应')
    }
  } catch (err) {
    connected.value = false
    error.value = '无法连接到Immich服务: ' + (err.message || '未知错误')
  } finally {
    loading.value = false
  }
}

const openImmich = async () => {
  try {
    const response = await axios.get('/api/services/immich/redirect')
    
    if (response.data.directLogin) {
      // 直接打开新窗口跳转
      window.open(response.data.directLogin, '_blank')
    } else {
      window.open(response.data.url, '_blank')
    }
  } catch (err) {
    error.value = '无法打开Immich: ' + (err.message || '未知错误')
  }
}

const refreshUser = () => {
  checkImmichConnection()
}

const getUserAvatar = () => {
  return currentUser.value?.avatar || ''
}

onMounted(() => {
  checkImmichConnection()
})
</script>

<style scoped>
.immich-integration {
  max-width: 600px;
  margin: 20px auto;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.card-content {
  text-align: center;
  padding: 20px;
}

.user-info {
  margin-bottom: 30px;
}

.user-info h3 {
  margin: 10px 0 5px 0;
  color: #303133;
}

.user-info p {
  color: #909399;
  margin-bottom: 20px;
}

.stats {
  display: flex;
  justify-content: center;
  gap: 30px;
  margin: 20px 0;
}

.launch-btn {
  width: 200px;
  height: 50px;
  font-size: 16px;
  margin: 0 10px;
}

.refresh-btn {
  margin-top: 10px;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
}

.error {
  text-align: center;
  padding: 40px;
  color: #f56c6c;
}

.el-icon {
  font-size: 48px;
  margin-bottom: 20px;
}
</style>
