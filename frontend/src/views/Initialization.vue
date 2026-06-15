<template>
  <div class="init-container">
    <!-- 液态玻璃背景 -->
    <div class="glass-background"></div>

    <div class="init-content">
      <!-- 头部 -->
      <div class="init-header">
        <div class="header-logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-7.94-7.94l-6.06 4.412A2 2 0 002.22 2.22 2 2 0 01.414 0 0 6 6 0 010 8.535 8.535 2 2 0 01.414 0 6 6 0 010-8.535 8.535 2 2 0 01-.414 0l-2.387-.477a6 6 0 00-7.94-7.94l-6.06 4.412a2 2 0 002.22-2.22 2 2 0 01-.414 0 6 6 0 010-8.535 8.535 2 2 0 01-.414 0l2.387.477a6 6 0 007.94 7.94l6.06-4.412a2 2 0 00-2.22-2.22 2 2 0 01-.414 0 6 6 0 01-8.535-8.535 2 2 0 01.414 0l2.387.477a6 6 0 007.94 7.94l6.06-4.412a2 2 0 002.22 2.22 2 2 0 01.414 0 6 6 0 01-8.535 8.535 2 2 0 01.414 0l-2.387.477a6 6 0 00-7.94-7.94z" />
          </svg>
        </div>
        <div class="header-text">
          <h1>NAS Dashboard</h1>
          <p>系统初始化向导</p>
        </div>
      </div>

      <!-- 步骤进度 -->
      <div class="steps-container">
        <div class="steps-wrapper">
          <div
            v-for="(step, index) in steps"
            :key="index"
            class="step-item"
            :class="{ active: currentStep === index, completed: currentStep > index }"
          >
            <div class="step-number">{{ index + 1 }}</div>
            <div class="step-text">{{ step.title }}</div>
          </div>
        </div>
      </div>

      <!-- 步骤内容 -->
      <div class="steps-content">
        <!-- 步骤1: 欢迎 -->
        <div v-if="currentStep === 0" class="step-panel">
          <div class="welcome-content">
            <div class="welcome-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-2m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-9" />
              </svg>
            </div>
            <h2>欢迎使用 NAS Dashboard</h2>
            <p>开启您的个人云端管理之旅</p>
            <p class="description">
              NAS Dashboard 将帮助您轻松管理磁盘、存储池、网络服务和用户权限。
              请完成以下基础设置开始使用。
            </p>
            <div class="system-requirements">
              <h4>系统配置要求：</h4>
              <ul>
                <li>✅ 设置管理员账户</li>
                <li>✅ 配置网络连接</li>
                <li>✅ 设置存储路径</li>
                <li>✅ 系统基础配置</li>
              </ul>
            </div>
            <button class="action-btn primary" @click="nextStep">
              开始设置
              <ArrowRightIcon class="w-5 h-5" />
            </button>
          </div>
        </div>

        <!-- 步骤2: 管理员设置 -->
        <div v-if="currentStep === 1" class="step-panel">
          <div class="form-content">
            <h2>设置管理员账户</h2>
            <p class="subtitle">该账户将拥有系统的最高管理权限</p>

            <form @submit.prevent="submitAdmin" class="glass-form">
              <div class="form-group">
                <label>管理员用户名 *</label>
                <input
                  type="text"
                  v-model="form.username"
                  placeholder="请输入用户名"
                  required
                  minlength="3"
                  maxlength="20"
                  pattern="[a-zA-Z][a-zA-Z0-9_-]*"
                />
                <small>3-20个字符，以字母开头，只能包含字母、数字、下划线和连字符</small>
              </div>

              <div class="form-group">
                <label>密码 *</label>
                <input
                  type="password"
                  v-model="form.password"
                  placeholder="请输入密码（至少8位）"
                  required
                  minlength="8"
                />
                <small>至少8个字符，建议包含字母、数字和特殊字符</small>
              </div>

              <div class="form-group">
                <label>确认密码 *</label>
                <input
                  type="password"
                  v-model="form.confirmPassword"
                  placeholder="请再次输入密码"
                  required
                  minlength="8"
                />
              </div>

              <div class="form-group">
                <label>电子邮件</label>
                <input
                  type="email"
                  v-model="form.email"
                  placeholder="admin@example.com (可选)"
                />
                <small>用于系统通知和密码重置</small>
              </div>

              <div class="form-actions">
                <button type="button" class="action-btn" @click="prevStep">
                  <ArrowLeftIcon class="w-4 h-4" />
                  上一步
                </button>
                <button type="submit" class="action-btn primary">
                  下一步
                  <ArrowRightIcon class="w-5 h-5" />
                </button>
              </div>
            </form>
          </div>
        </div>

        <!-- 步骤3: 网络配置 -->
        <div v-if="currentStep === 2" class="step-panel">
          <div class="network-content">
            <h2>网络配置</h2>
            <p class="subtitle">选择网络连接方式</p>

            <div class="network-options">
              <div
                class="network-card"
                :class="{ selected: form.networkType === 'dhcp' }"
                @click="form.networkType = 'dhcp'"
              >
                <div class="network-icon">
                  <WifiIcon class="w-8 h-8" />
                </div>
                <h3>自动获取 (DHCP)</h3>
                <p>自动从路由器获取IP地址，简单推荐</p>
                <span class="network-badge" v-if="form.networkType === 'dhcp'">已选择</span>
              </div>

              <div
                class="network-card"
                :class="{ selected: form.networkType === 'static' }"
                @click="form.networkType = 'static'"
              >
                <div class="network-icon">
                  <ServerIcon class="w-8 h-8" />
                </div>
                <h3>固定IP (Static)</h3>
                <p>手动设置固定的IP地址</p>
                <span class="network-badge" v-if="form.networkType === 'static'">已选择</span>
              </div>
            </div>

            <!-- 固定IP设置 -->
            <div v-if="form.networkType === 'static'" class="static-ip-config">
              <h4>固定IP设置</h4>
              <div class="ip-form">
                <div class="form-group">
                  <label>IP地址</label>
                  <input
                    type="text"
                    v-model="form.staticIp"
                    placeholder="192.168.1.100"
                    pattern="^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$"
                    required
                  />
                </div>

                <div class="form-group">
                  <label>子网掩码</label>
                  <select v-model="form.subnetMask" required>
                    <option value="255.255.255.0">255.255.255.0</option>
                    <option value="255.255.0.0">255.255.0.0</option>
                    <option value="255.255.255.0">255.255.255.0</option>
                    <option value="255.0.0.0">255.0.0.0</option>
                  </select>
                </div>

                <div class="form-group">
                  <label>网关</label>
                  <input
                    type="text"
                    v-model="form.gateway"
                    placeholder="192.168.1.1"
                    required
                  />
                </div>

                <div class="form-group">
                  <label>DNS服务器</label>
                  <input
                    type="text"
                    v-model="form.dns"
                    placeholder="8.8.8.8"
                  />
                  <small>可输入多个，用逗号分隔</small>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="action-btn" @click="prevStep">
                <ArrowLeftIcon class="w-4 h-4" />
                上一步
              </button>
              <button type="button" class="action-btn primary" @click="nextStep">
                下一步
                <ArrowRightIcon class="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>

        <!-- 步骤4: 存储设置 -->
        <div v-if="currentStep === 3" class="step-panel">
          <div class="storage-content">
            <h2>存储设置</h2>
            <p class="subtitle">设置数据存储路径</p>

            <div class="storage-config">
              <div class="form-group">
                <label>数据存储路径 *</label>
                <input
                  type="text"
                  v-model="form.storagePath"
                  placeholder="/mnt/data"
                  required
                />
                <small>默认存储路径，确保有足够的磁盘空间</small>
              </div>

              <div class="form-group">
                <label>备份存储路径</label>
                <input
                  type="text"
                  v-model="form.backupPath"
                  placeholder="/mnt/backup"
                />
                <small>用于系统备份和数据保护</small>
              </div>

              <div class="disk-info">
                <h4>可用磁盘空间</h4>
                <div class="disk-space">
                  <div class="disk-item">
                    <span class="disk-label">根目录 (/)</span>
                    <div class="disk-bar">
                      <div class="disk-fill" style="width: 45%"></div>
                    </div>
                    <span class="disk-value">~45%</span>
                  </div>
                  <div class="disk-item">
                    <span class="disk-label">/mnt</span>
                    <div class="disk-bar">
                      <div class="disk-fill" style="width: 30%"></div>
                    </div>
                    <span class="disk-value">~30%</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="action-btn" @click="prevStep">
                <ArrowLeftIcon class="w-4 h-4" />
                上一步
              </button>
              <button type="button" class="action-btn primary" @click="nextStep">
                下一步
                <ArrowRightIcon class="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>

        <!-- 步骤5: 系统配置 -->
        <div v-if="currentStep === 4" class="step-panel">
          <div class="system-content">
            <h2>系统配置</h2>
            <p class="subtitle">完成最后的基础设置</p>

            <div class="system-config">
              <div class="config-item">
                <h4>系统名称</h4>
                <input
                  type="text"
                  v-model="form.hostname"
                  placeholder="nas-server"
                />
                <small>用于识别您的NAS设备</small>
              </div>

              <div class="config-item">
                <h4>时区设置</h4>
                <select v-model="form.timezone">
                  <option value="Asia/Shanghai">中国标准时间 (UTC+8)</option>
                  <option value="Asia/Hong_Kong">香港时间 (UTC+8)</option>
                  <option value="Asia/Tokyo">日本时间 (UTC+9)</option>
                  <option value="America/New_York">美东时间 (UTC-5)</option>
                  <option value="Europe/London">格林威治时间 (UTC+0)</option>
                </select>
              </div>

              <div class="config-item">
                <h4>安全选项</h4>
                <div class="checkbox-group">
                  <label class="checkbox-item">
                    <input type="checkbox" v-model="form.enableFirewall" />
                    <span>启用防火墙 (推荐)</span>
                  </label>
                  <label class="checkbox-item">
                    <input type="checkbox" v-model="form.enableAutoUpdate" />
                    <span>自动更新系统</span>
                  </label>
                  <label class="checkbox-item">
                    <input type="checkbox" v-model="form.enableTelemetry" />
                    <span>发送使用统计 (帮助改进产品)</span>
                  </label>
                </div>
              </div>
            </div>

            <div class="form-actions">
              <button type="button" class="action-btn" @click="prevStep">
                <ArrowLeftIcon class="w-4 h-4" />
                上一步
              </button>
              <button type="button" class="action-btn primary" @click="submitInitialization">
                完成初始化
                <CheckIcon class="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>

        <!-- 步骤6: 完成状态 -->
        <div v-if="currentStep === 5" class="step-panel">
          <div class="complete-content">
            <div v-if="loading" class="loading-state">
              <div class="spinner"></div>
              <p>正在配置系统...</p>
            </div>

            <div v-else-if="error" class="error-state">
              <div class="error-icon">
                <XCircleIcon class="w-16 h-16" />
              </div>
              <h2>初始化失败</h2>
              <p>{{ error }}</p>
              <button class="action-btn primary" @click="retryInitialization">
                重试
              </button>
              <button class="action-btn" @click="backToStart">
                重新开始
              </button>
            </div>

            <div v-else class="success-state">
              <div class="success-icon">
                <CheckCircleIcon class="w-16 h-16" />
              </div>
              <h2>初始化完成！</h2>
              <p>您的NAS Dashboard已成功初始化。</p>
              <div class="success-info">
                <div class="info-item">
                  <span class="info-label">管理员用户:</span>
                  <span class="info-value">{{ form.username }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">存储路径:</span>
                  <span class="info-value">{{ form.storagePath }}</span>
                </div>
              </div>
              <button class="action-btn primary large" @click="goToDesktop">
                进入桌面
                <ArrowRightIcon class="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import axios from 'axios'
import {
  ArrowRightIcon,
  ArrowLeftIcon,
  CheckIcon,
  XCircleIcon,
  CheckCircleIcon,
  WifiIcon,
  ServerIcon
} from '@heroicons/vue/24/outline'

const router = useRouter()

// 状态
const currentStep = ref(0)
const loading = ref(false)
const error = ref('')

// 步骤定义
const steps = [
  { title: '欢迎' },
  { title: '管理员设置' },
  { title: '网络配置' },
  { title: '存储设置' },
  { title: '系统配置' },
  { title: '完成' }
]

// 表单数据
const form = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  networkType: 'dhcp',
  staticIp: '',
  subnetMask: '255.255.255.0',
  gateway: '',
  dns: '',
  storagePath: '/mnt/data',
  backupPath: '/mnt/backup',
  hostname: 'nas-server',
  timezone: 'Asia/Shanghai',
  enableFirewall: true,
  enableAutoUpdate: true,
  enableTelemetry: false
})

// 方法
const nextStep = () => {
  // 验证当前步骤
  if (currentStep.value === 1) {
    if (!form.username || !form.password || !form.confirmPassword) {
      ElMessage.warning('请填写必填项')
      return
    }
    if (form.password !== form.confirmPassword) {
      ElMessage.error('两次输入的密码不一致')
      return
    }
  }

  if (currentStep.value === 2) {
    if (form.networkType === 'static' && (!form.staticIp || !form.gateway)) {
      ElMessage.warning('请填写固定IP设置')
      return
    }
  }

  if (currentStep.value === 3) {
    if (!form.storagePath) {
      ElMessage.warning('请设置存储路径')
      return
    }
  }

  if (currentStep.value < 5) {
    currentStep.value++
  }
}

const prevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const submitAdmin = () => {
  nextStep()
}

const submitInitialization = async () => {
  loading.value = true
  error.value = ''

  try {
    const response = await axios.post('/api/system/initialize', {
      adminUsername: form.username,
      adminPassword: form.password,
      adminEmail: form.email,
      systemName: form.hostname
    })

    if (response.data.message || response.data.success) {
      currentStep.value = 5
      // 自动登录（开发模式）
      if (import.meta.env.DEV) {
        localStorage.setItem('token', 'dev-token-' + Date.now())
        localStorage.setItem('user', JSON.stringify({
          username: form.username,
          role: 'admin',
          email: form.email
        }))
      }
      loading.value = false
    } else {
      throw new Error(response.data.message || '初始化失败')
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || err.response?.data?.message || err.message || '初始化失败，请重试'
    loading.value = false
    currentStep.value = 5 // 跳到错误页面
  }
}

const retryInitialization = () => {
  submitInitialization()
}

const backToStart = () => {
  currentStep.value = 0
  error.value = ''
}

const goToDesktop = () => {
  router.push('/desktop')
}

// 生命周期
onMounted(() => {
  // 检查初始化状态
  checkInitStatus()
})

const checkInitStatus = async () => {
  try {
    const response = await axios.get('/api/system/init-status')
    if (response.data.initialized) {
      // 已经初始化，直接跳转到桌面
      router.push('/desktop')
    }
  } catch (err) {
    console.error('Failed to check init status:', err)
    // 如果API失败，显示初始化页面
  }
}
</script>

<style scoped lang="scss">
.init-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.glass-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.05) 100%),
                  radial-gradient(circle at 20% 80%, rgba(102, 126, 234, 0.15) 0%, transparent 50%),
                  radial-gradient(circle at 80% 20%, rgba(118, 75, 162, 0.15) 0%, transparent 50%);
  z-index: 0;
}

.init-content {
  position: relative;
  z-index: 1;
  width: 90%;
  max-width: 800px;
  margin: 0 auto;
}

.init-header {
  text-align: center;
  margin-bottom: 40px;
  padding: 30px;

  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  border: 1px solid rgba(102, 126, 234, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);

  .header-logo {
    margin-bottom: 20px;

    .logo-icon {
      width: 80px;
      height: 80px;
      margin: 0 auto;
      color: #667eea;
    }
  }

  h1 {
    font-size: 32px;
    font-weight: 700;
    color: #1f2937;
    margin: 0 0 8px 0;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  p {
    font-size: 16px;
    color: #6b7280;
    margin: 0;
  }
}

.steps-container {
  margin-bottom: 30px;

  .steps-wrapper {
    display: flex;
    justify-content: center;
    gap: 8px;
    padding: 20px;
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    border-radius: 12px;
    border: 1px solid rgba(102, 126, 234, 0.1);
  }

  .step-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 20px;
    border-radius: 8px;
    transition: all 0.3s;

    &.active {
      background: linear-gradient(135deg, rgba(102, 126, 234, 0.2), rgba(118, 75, 162, 0.2));
      border-color: rgba(102, 126, 234, 0.3);

      .step-number {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
      }
    }

    &.completed {
      .step-number {
        background: #10b981;
        color: white;
      }
    }

    .step-number {
      width: 32px;
      height: 32px;
      border-radius: 50%;
      background: rgba(102, 126, 234, 0.2);
      color: #667eea;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      font-size: 14px;
    }

    .step-text {
      font-weight: 500;
      color: #374151;
    }
  }
}

.steps-content {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  padding: 40px;
  border: 1px solid rgba(102, 126, 234, 0.15);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  min-height: 400px;
}

.step-panel {
  animation: fadeIn 0.3s ease-in-out;

  h2 {
    font-size: 24px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 8px 0;
  }

  .subtitle {
    font-size: 14px;
    color: #6b7280;
    margin-bottom: 24px;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

// 欢迎页面
.welcome-content {
  text-align: center;
  padding: 20px;

  .welcome-icon {
    margin-bottom: 24px;

    svg {
      width: 100px;
      height: 100px;
      margin: 0 auto;
      color: #667eea;
    }
  }

  h2 {
    font-size: 28px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 12px 0;
  }

  .description {
    font-size: 16px;
    color: #6b7280;
    line-height: 1.6;
    margin-bottom: 24px;
    max-width: 500px;
    margin-left: auto;
    margin-right: auto;
  }

  .system-requirements {
    background: rgba(102, 126, 234, 0.05);
    border-radius: 12px;
    padding: 20px;
    margin: 20px auto;
    max-width: 500px;
    text-align: left;

    h4 {
      font-size: 16px;
      font-weight: 600;
      color: #1f2937;
      margin: 0 0 12px 0;
    }

    ul {
      list-style: none;
      padding: 0;
      margin: 0;

      li {
        padding: 8px 0;
        font-size: 14px;
        color: #374151;
      }
    }
  }

  .action-btn {
    margin-top: 24px;
  }
}

// 玻璃表单样式
.glass-form {
  max-width: 400px;
  margin: 0 auto;

  .form-group {
    margin-bottom: 20px;

    label {
      display: block;
      font-size: 14px;
      font-weight: 500;
      color: #374151;
      margin-bottom: 8px;
    }

    input, select {
      width: 100%;
      padding: 12px 16px;
      border: 1px solid rgba(102, 126, 234, 0.2);
      border-radius: 8px;
      font-size: 14px;
      background: rgba(255, 255, 255, 0.8);
      backdrop-filter: blur(10px);
      transition: all 0.2s;

      &:focus {
        outline: none;
        border-color: #667eea;
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2);
      }
    }

    select {
      appearance: none;
      background-image: url("data:image/svg+xml,%3csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24'%3e%3cpath stroke='%236667eea' stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7M7 7v10m0-10a6 6 0 1112 0 6 6 0 1112 0' /%3e%3c/svg%3e");
      background-repeat: no-repeat;
      background-position: right 12px center;
      background-size: 16px;
    }

    small {
      display: block;
      font-size: 12px;
      color: #9ca3af;
      margin-top: 4px;
    }
  }

  .form-actions {
    display: flex;
    gap: 12px;
    justify-content: center;
    margin-top: 24px;
  }
}

// 网络选择
.network-options {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.network-card {
  background: rgba(255, 255, 255, 0.8);
  border: 2px solid rgba(102, 126, 234, 0.1);
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
  position: relative;

  &:hover {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.05);
  }

  &.selected {
    border-color: #667eea;
    background: rgba(102, 126, 234, 0.1);
    box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2);
  }

  .network-icon {
    margin-bottom: 12px;

    svg {
      width: 48px;
      height: 48px;
      margin: 0 auto;
      color: #667eea;
    }
  }

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 8px 0;
  }

  p {
    font-size: 13px;
    color: #6b7280;
    margin: 0 0 0 0;
  }

  .network-badge {
    position: absolute;
    top: 12px;
    right: 12px;
    padding: 4px 8px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    border-radius: 12px;
    font-size: 11px;
    font-weight: 500;
  }
}

.static-ip-config {
  background: rgba(102, 126, 234, 0.05);
  border-radius: 12px;
  padding: 20px;
  margin-top: 20px;

  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 16px 0;
  }
}

.ip-form {
  display: grid;
  gap: 16px;
}

// 磁盘空间展示
.disk-info {
  margin-top: 20px;

  h4 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0 0 16px 0;
  }
}

.disk-space {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.disk-item {
  display: flex;
  align-items: center;
  gap: 12px;

  .disk-label {
    font-size: 13px;
    color: #6b7280;
    min-width: 80px;
  }

  .disk-bar {
    flex: 1;
    height: 6px;
    background: rgba(239, 68, 68, 0.1);
    border-radius: 3px;
    overflow: hidden;
  }

  .disk-fill {
    height: 100%;
    background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
    border-radius: 3px;
    transition: width 0.3s;
  }

  .disk-value {
    font-size: 12px;
    font-weight: 500;
    color: #667eea;
    min-width: 40px;
    text-align: right;
  }
}

// 系统配置
.system-config {
  display: grid;
  gap: 20px;

  .config-item {
    h4 {
      font-size: 14px;
      font-weight: 600;
      color: #374151;
      margin: 0 0 8px 0;
    }

    input, select {
      width: 100%;
      padding: 10px 12px;
      border: 1px solid rgba(102, 126, 234, 0.2);
      border-radius: 8px;
      background: rgba(255, 255, 255, 0.8);
      backdrop-filter: blur(10px);
      font-size: 14px;
      transition: all 0.2s;

      &:focus {
        outline: none;
        border-color: #667eea;
      }
    }
  }
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: rgba(102, 126, 234, 0.1);
  }

  input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  span {
    font-size: 14px;
    color: #374151;
  }
}

// 操作按钮
.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  color: white;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);

  &:hover:not(:disabled) {
    opacity: 0.9;
    transform: translateY(-1px);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
    transform: none;
  }

  &.large {
    padding: 16px 32px;
    font-size: 16px;
  }

  svg {
    width: 20px;
    height: 20px;
  }
}

// 加载和错误状态
.loading-state,
.error-state,
.success-state {
  text-align: center;
  padding: 40px 20px;

  .spinner {
    width: 48px;
    height: 48px;
    border: 3px solid rgba(102, 126, 234, 0.2);
    border-top-color: #667eea;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin: 0 auto 16px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .loading-state p,
  .error-state p,
  .success-state p {
    font-size: 16px;
    color: #6b7280;
    margin: 16px 0;
  }

  .error-icon,
  .success-icon {
    margin-bottom: 16px;

    svg {
      width: 64px;
      height: 64px;
      margin: 0 auto;
    }
  }

  .error-icon svg {
    color: #ef4444;
  }

  .success-icon svg {
    color: #10b981;
  }

  h2 {
    font-size: 24px;
    font-weight: 600;
    margin: 0 0 8px 0;
  }

  .action-btn {
    margin: 8px;
  }
}

.success-info {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 12px;
  padding: 20px;
  margin: 24px auto;
  max-width: 400px;

  .info-item {
    display: flex;
    justify-content: space-between;
    padding: 12px 0;
    border-bottom: 1px solid rgba(16, 185, 129, 0.1);

    &:last-child {
      border-bottom: none;
    }

    .info-label {
      font-size: 13px;
      color: #6b7280;
      font-weight: 500;
    }

    .info-value {
      font-size: 14px;
      color: #1f2937;
      font-weight: 600;
    }
  }
}
</style>