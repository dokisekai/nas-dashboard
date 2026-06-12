<template>
  <div class="system-settings">
    <div class="settings-header">
      <h1>系统设置</h1>
      <p class="subtitle">系统配置和管理</p>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.label }}
      </button>
    </div>

    <!-- Network Tab -->
    <div v-if="activeTab === 'network'" class="tab-content">
      <div class="section-header">
        <h2>网络配置</h2>
      </div>

      <div class="settings-grid">
        <!-- Network Interfaces -->
        <div class="setting-card">
          <div class="card-header">
            <h3>网络接口</h3>
            <button class="action-btn" @click="refreshNetwork">
              <ArrowPathIcon class="w-4 h-4" />
              刷新
            </button>
          </div>
          <div class="interfaces-list">
            <div
              v-for="iface in networkInterfaces"
              :key="iface.name"
              class="interface-item"
            >
              <div class="interface-status" :class="{ active: iface.status === 'up' }"></div>
              <div class="interface-info">
                <h4>{{ iface.name }}</h4>
                <p>{{ iface.ip }}</p>
                <p>MAC: {{ iface.mac }}</p>
              </div>
              <div class="interface-actions">
                <button class="action-btn" @click="editInterface(iface)">
                  <PencilIcon class="w-4 h-4" />
                  编辑
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- DNS Configuration -->
        <div class="setting-card">
          <div class="card-header">
            <h3>DNS 配置</h3>
          </div>
          <form @submit.prevent="saveDNS" class="dns-form">
            <div class="form-group">
              <label>主 DNS</label>
              <input
                v-model="dnsConfig.primary"
                type="text"
                placeholder="8.8.8.8"
              />
            </div>
            <div class="form-group">
              <label>备用 DNS</label>
              <input
                v-model="dnsConfig.secondary"
                type="text"
                placeholder="8.8.4.4"
              />
            </div>
            <button type="submit" class="action-btn primary">
              保存配置
            </button>
          </form>
        </div>

        <!-- Proxy Settings -->
        <div class="setting-card">
          <div class="card-header">
            <h3>代理设置</h3>
          </div>
          <form @submit.prevent="saveProxy" class="proxy-form">
            <div class="form-group">
              <label>
                <input v-model="proxyConfig.enabled" type="checkbox" />
                启用代理
              </label>
            </div>
            <div v-if="proxyConfig.enabled" class="proxy-fields">
              <div class="form-group">
                <label>代理服务器</label>
                <input
                  v-model="proxyConfig.server"
                  type="text"
                  placeholder="proxy.example.com"
                />
              </div>
              <div class="form-group">
                <label>端口</label>
                <input
                  v-model="proxyConfig.port"
                  type="number"
                  placeholder="8080"
                />
              </div>
              <div class="form-group">
                <label>用户名 (可选)</label>
                <input
                  v-model="proxyConfig.username"
                  type="text"
                />
              </div>
              <div class="form-group">
                <label>密码 (可选)</label>
                <input
                  v-model="proxyConfig.password"
                  type="password"
                />
              </div>
            </div>
            <button type="submit" class="action-btn primary">
              保存配置
            </button>
          </form>
        </div>
      </div>
    </div>

    <!-- System Info Tab -->
    <div v-if="activeTab === 'info'" class="tab-content">
      <div class="section-header">
        <h2>系统信息</h2>
        <button class="action-btn" @click="refreshSystemInfo">
          <ArrowPathIcon class="w-4 h-4" />
          刷新
        </button>
      </div>

      <div class="system-info-grid">
        <div class="info-card">
          <div class="info-icon">
            <ServerIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>主机名</h3>
            <p>{{ systemInfo.hostname }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <CpuChipIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>CPU</h3>
            <p>{{ systemInfo.cpu }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <MemoryChipIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>内存</h3>
            <p>{{ systemInfo.memory }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <CircleStackIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>存储</h3>
            <p>{{ systemInfo.storage }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <BeakerIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>操作系统</h3>
            <p>{{ systemInfo.os }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <ClockIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>运行时间</h3>
            <p>{{ systemInfo.uptime }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <KernelIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>内核版本</h3>
            <p>{{ systemInfo.kernel }}</p>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon">
            <CodeBracketIcon class="w-8 h-8" />
          </div>
          <div class="info-content">
            <h3>架构</h3>
            <p>{{ systemInfo.architecture }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Services Tab -->
    <div v-if="activeTab === 'services'" class="tab-content">
      <div class="section-header">
        <h2>系统服务</h2>
        <div class="filters">
          <input
            v-model="serviceSearch"
            type="text"
            placeholder="搜索服务..."
            class="search-input"
          />
        </div>
      </div>

      <div class="services-list">
        <div
          v-for="service in filteredServices"
          :key="service.name"
          class="service-item"
        >
          <div class="service-info">
            <div class="service-status" :class="service.status.toLowerCase()"></div>
            <div class="service-details">
              <h4>{{ service.name }}</h4>
              <p>{{ service.description }}</p>
            </div>
          </div>

          <div class="service-meta">
            <span class="service-enabled" :class="{ active: service.enabled }">
              {{ service.enabled ? '已启用' : '已禁用' }}
            </span>
          </div>

          <div class="service-actions">
            <button
              v-if="service.status === 'running'"
              class="action-btn warning"
              @click="stopService(service)"
            >
              <StopIcon class="w-4 h-4" />
              停止
            </button>
            <button
              v-else
              class="action-btn primary"
              @click="startService(service)"
            >
              <PlayIcon class="w-4 h-4" />
              启动
            </button>
            <button class="action-btn" @click="restartService(service)">
              <ArrowPathIcon class="w-4 h-4" />
              重启
            </button>
            <button
              class="action-btn"
              :class="{ danger: service.enabled }"
              @click="toggleService(service)"
            >
              <PowerIcon class="w-4 h-4" />
              {{ service.enabled ? '禁用' : '启用' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Updates Tab -->
    <div v-if="activeTab === 'updates'" class="tab-content">
      <div class="section-header">
        <h2>系统更新</h2>
        <button class="action-btn" @click="checkUpdates">
          <ArrowPathIcon class="w-4 h-4" />
          检查更新
        </button>
      </div>

      <div class="updates-section">
        <div class="update-card current">
          <div class="update-header">
            <h3>当前版本</h3>
            <span class="version-badge">{{ currentVersion }}</span>
          </div>
          <div class="update-info">
            <p><strong>发布日期:</strong> {{ currentVersionDate }}</p>
            <p><strong>系统状态:</strong> {{ updateStatus }}</p>
          </div>
        </div>

        <div v-if="availableUpdate" class="update-card available">
          <div class="update-header">
            <h3>可用更新</h3>
            <span class="version-badge new">{{ availableUpdate.version }}</span>
          </div>
          <div class="update-info">
            <p><strong>发布日期:</strong> {{ availableUpdate.date }}</p>
            <p><strong>更新说明:</strong></p>
            <ul>
              <li v-for="change in availableUpdate.changes" :key="change">
                {{ change }}
              </li>
            </ul>
          </div>
          <button class="action-btn primary" @click="installUpdate">
            <CloudArrowDownIcon class="w-4 h-4" />
            安装更新
          </button>
        </div>

        <div v-else class="no-updates">
          <CheckIcon class="w-12 h-12" />
          <p>系统已是最新版本</p>
        </div>
      </div>

      <div class="update-settings">
        <h3>更新设置</h3>
        <form @submit.prevent="saveUpdateSettings" class="settings-form">
          <div class="form-group">
            <label>
              <input v-model="updateSettings.autoCheck" type="checkbox" />
              自动检查更新
            </label>
          </div>
          <div class="form-group">
            <label>
              <input v-model="updateSettings.autoInstall" type="checkbox" />
              自动安装安全更新
            </label>
          </div>
          <div class="form-group">
            <label>检查频率</label>
            <select v-model="updateSettings.frequency">
              <option value="daily">每天</option>
              <option value="weekly">每周</option>
              <option value="monthly">每月</option>
            </select>
          </div>
          <button type="submit" class="action-btn primary">
            保存设置
          </button>
        </form>
      </div>
    </div>

    <!-- Backup Tab -->
    <div v-if="activeTab === 'backup'" class="tab-content">
      <div class="section-header">
        <h2>备份设置</h2>
        <button class="action-btn primary" @click="showCreateBackupModal = true">
          <PlusIcon class="w-4 h-4" />
          创建备份
        </button>
      </div>

      <div class="backups-list">
        <div
          v-for="backup in backups"
          :key="backup.id"
          class="backup-item"
        >
          <div class="backup-icon">
            <CircleStackIcon class="w-8 h-8" />
          </div>
          <div class="backup-details">
            <h4>{{ backup.name }}</h4>
            <p>{{ backup.description }}</p>
            <div class="backup-meta">
              <span>创建时间: {{ formatDate(backup.created) }}</span>
              <span>大小: {{ formatBytes(backup.size) }}</span>
              <span>类型: {{ backup.type }}</span>
            </div>
          </div>

          <div class="backup-actions">
            <button class="action-btn" @click="restoreBackup(backup)">
              <ArrowUturnUpIcon class="w-4 h-4" />
              恢复
            </button>
            <button class="action-btn" @click="downloadBackup(backup)">
              <CloudArrowDownIcon class="w-4 h-4" />
              下载
            </button>
            <button class="action-btn danger" @click="deleteBackup(backup)">
              <TrashIcon class="w-4 h-4" />
              删除
            </button>
          </div>
        </div>
      </div>

      <div class="backup-schedule">
        <h3>备份计划</h3>
        <form @submit.prevent="saveBackupSchedule" class="schedule-form">
          <div class="form-group">
            <label>
              <input v-model="backupSchedule.enabled" type="checkbox" />
              启用自动备份
            </label>
          </div>
          <div v-if="backupSchedule.enabled" class="schedule-fields">
            <div class="form-group">
              <label>备份频率</label>
              <select v-model="backupSchedule.frequency">
                <option value="daily">每天</option>
                <option value="weekly">每周</option>
                <option value="monthly">每月</option>
              </select>
            </div>
            <div class="form-group">
              <label>备份时间</label>
              <input
                v-model="backupSchedule.time"
                type="time"
              />
            </div>
            <div class="form-group">
              <label>保留备份数量</label>
              <input
                v-model="backupSchedule.retention"
                type="number"
                min="1"
                max="30"
              />
            </div>
          </div>
          <button type="submit" class="action-btn primary">
            保存计划
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  GlobeAltIcon,
  ServerIcon,
  CpuChipIcon,
  CircleStackIcon,
  BeakerIcon,
  ClockIcon,
  ArrowPathIcon,
  PencilIcon,
  PlayIcon,
  StopIcon,
  PowerIcon,
  CheckIcon,
  CloudArrowDownIcon,
  PlusIcon,
  TrashIcon,
  ArrowUturnUpIcon,
  CodeBracketIcon
} from '@heroicons/vue/24/outline'
import { systemApi, serviceApi } from '../api'

const activeTab = ref('network')

const tabs = [
  { id: 'network', label: '网络', icon: GlobeAltIcon },
  { id: 'info', label: '系统信息', icon: ServerIcon },
  { id: 'services', label: '服务', icon: CpuChipIcon },
  { id: 'updates', label: '更新', icon: BeakerIcon },
  { id: 'backup', label: '备份', icon: CircleStackIcon }
]

// Network
const networkInterfaces = ref<any[]>([])
const dnsConfig = ref({
  primary: '8.8.8.8',
  secondary: '8.8.4.4'
})
const proxyConfig = ref({
  enabled: false,
  server: '',
  port: '',
  username: '',
  password: ''
})

// System Info
const systemInfo = ref<any>({
  hostname: '',
  cpu: '',
  memory: '',
  storage: '',
  os: '',
  uptime: '',
  kernel: '',
  architecture: ''
})

// Services
const serviceSearch = ref('')
const services = ref<any[]>([])

const filteredServices = computed(() => {
  if (!serviceSearch.value) return services.value
  return services.value.filter(s =>
    s.name.toLowerCase().includes(serviceSearch.value.toLowerCase()) ||
    s.description.toLowerCase().includes(serviceSearch.value.toLowerCase())
  )
})

// Updates
const currentVersion = ref('1.0.0')
const currentVersionDate = ref('2024-01-01')
const updateStatus = ref('最新')
const availableUpdate = ref<any>(null)
const updateSettings = ref({
  autoCheck: true,
  autoInstall: true,
  frequency: 'weekly'
})

// Backup
const backups = ref<any[]>([])
const showCreateBackupModal = ref(false)
const backupSchedule = ref({
  enabled: true,
  frequency: 'daily',
  time: '02:00',
  retention: 7
})

// API Functions
const refreshNetwork = async () => {
  try {
    // Mock data
    networkInterfaces.value = [
      {
        name: 'eth0',
        ip: '192.168.1.100',
        mac: '00:11:22:33:44:55',
        status: 'up'
      },
      {
        name: 'eth1',
        ip: '10.0.0.5',
        mac: '00:11:22:33:44:56',
        status: 'up'
      }
    ]
  } catch (error: any) {
    console.error('Failed to refresh network:', error)
  }
}

const saveDNS = () => {
  console.log('Saving DNS config:', dnsConfig.value)
  alert('DNS 配置已保存')
}

const saveProxy = () => {
  console.log('Saving proxy config:', proxyConfig.value)
  alert('代理配置已保存')
}

const editInterface = (iface: any) => {
  console.log('Editing interface:', iface.name)
  alert('编辑接口: ' + iface.name)
}

const refreshSystemInfo = async () => {
  try {
    const response = await systemApi.getInfo()
    const info = response.data.info

    // 解析真实的系统信息
    const unameParts = info.system?.uname?.split(' ') || []
    const osName = unameParts.slice(0, 3).join(' ')

    systemInfo.value = {
      hostname: info.hostname || 'Unknown',
      cpu: info.cpu ? `${info.cpu.model} @ ${info.cpu.mhz}MHz` : 'Unknown CPU',
      memory: info.memory ? `${formatMemory(info.memory.total)} DDR4` : 'Unknown',
      storage: info.disks && info.disks.length > 0 ? formatStorage(info.disks) : 'Unknown',
      os: osName || 'Unknown OS',
      uptime: info.uptime || 'Unknown',
      kernel: extractKernelVersion(info.system?.uname) || 'Unknown',
      architecture: extractArchitecture(info.system?.uname) || 'Unknown'
    }
  } catch (error: any) {
    console.error('Failed to fetch system info:', error)
    // 如果API调用失败，显示错误状态而不是假数据
    systemInfo.value = {
      hostname: '无法获取',
      cpu: 'API调用失败',
      memory: 'API调用失败',
      storage: 'API调用失败',
      os: 'API调用失败',
      uptime: 'API调用失败',
      kernel: 'API调用失败',
      architecture: 'API调用失败'
    }
  }
}

// 提取内核版本号
const extractKernelVersion = (uname: string | undefined): string => {
  if (!uname) return 'Unknown'
  const match = uname.match(/(\d+\.\d+\.\d+-\d+)/)
  return match ? match[1] : 'Unknown'
}

// 提取系统架构
const extractArchitecture = (uname: string | undefined): string => {
  if (!uname) return 'Unknown'
  const match = uname.match(/(x86_64|aarch64|armv7l|i686)/)
  return match ? match[1] : 'Unknown'
}

// 格式化内存大小
const formatMemory = (bytes: number) => {
  const gb = bytes / (1024 * 1024 * 1024)
  if (gb >= 1024) {
    return `${(gb / 1024).toFixed(1)} TB`
  }
  return `${gb.toFixed(1)} GB`
}

// 格式化存储信息
const formatStorage = (disks: any[]) => {
  const physicalDisks = disks.filter(d => d.mountpoint && d.mountpoint.includes('/mnt/.physical'))
  if (physicalDisks.length === 0) {
    return `${disks.length} 个磁盘`
  }

  const totalStorage = physicalDisks.reduce((sum, disk) => sum + disk.total, 0)
  const totalStorageTB = totalStorage / (1024 * 1024 * 1024 * 1024)

  return `${physicalDisks.length} 个磁盘，总计 ${totalStorageTB.toFixed(1)} TB`
}

const loadServices = async () => {
  try {
    const response = await serviceApi.getServices()
    services.value = response.data
  } catch (error: any) {
    console.error('Failed to load services:', error)
    services.value = [
      { name: 'nginx', description: 'Web Server', status: 'running', enabled: true },
      { name: 'mysql', description: 'Database Server', status: 'running', enabled: true },
      { name: 'docker', description: 'Container Runtime', status: 'running', enabled: true },
      { name: 'ssh', description: 'SSH Server', status: 'running', enabled: true },
      { name: 'smb', description: 'SMB File Sharing', status: 'running', enabled: true }
    ]
  }
}

const startService = async (service: any) => {
  try {
    await serviceApi.startService(service.name)
    await loadServices()
  } catch (error: any) {
    console.error('Failed to start service:', error)
    alert('启动失败: ' + error.message)
  }
}

const stopService = async (service: any) => {
  try {
    await serviceApi.stopService(service.name)
    await loadServices()
  } catch (error: any) {
    console.error('Failed to stop service:', error)
    alert('停止失败: ' + error.message)
  }
}

const restartService = async (service: any) => {
  try {
    await serviceApi.restartService(service.name)
    await loadServices()
  } catch (error: any) {
    console.error('Failed to restart service:', error)
    alert('重启失败: ' + error.message)
  }
}

const toggleService = async (service: any) => {
  try {
    if (service.enabled) {
      await serviceApi.disableService(service.name)
    } else {
      await serviceApi.enableService(service.name)
    }
    await loadServices()
  } catch (error: any) {
    console.error('Failed to toggle service:', error)
    alert('操作失败: ' + error.message)
  }
}

const checkUpdates = () => {
  console.log('Checking for updates...')
  // Simulate update check
  setTimeout(() => {
    availableUpdate.value = {
      version: '1.1.0',
      date: '2024-01-20',
      changes: [
        '新增功能: 改进的存储管理界面',
        '性能优化: 系统监控速度提升30%',
        '安全修复: 修复了多个安全漏洞',
        '稳定性: 改进了系统稳定性'
      ]
    }
  }, 1000)
}

const installUpdate = () => {
  if (confirm('确定要安装更新吗? 系统将在安装过程中重启。')) {
    console.log('Installing update...')
    alert('更新安装已开始，系统将自动重启')
  }
}

const saveUpdateSettings = () => {
  console.log('Saving update settings:', updateSettings.value)
  alert('更新设置已保存')
}

const restoreBackup = (backup: any) => {
  if (confirm(`确定要恢复备份 "${backup.name}" 吗? 当前系统配置将被覆盖。`)) {
    console.log('Restoring backup:', backup.name)
    alert('恢复备份: ' + backup.name)
  }
}

const downloadBackup = (backup: any) => {
  console.log('Downloading backup:', backup.name)
  alert('下载备份: ' + backup.name)
}

const deleteBackup = (backup: any) => {
  if (confirm(`确定要删除备份 "${backup.name}" 吗?`)) {
    backups.value = backups.value.filter(b => b.id !== backup.id)
    alert('备份已删除')
  }
}

const saveBackupSchedule = () => {
  console.log('Saving backup schedule:', backupSchedule.value)
  alert('备份计划已保存')
}

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

const formatDate = (dateStr: string) => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

onMounted(() => {
  refreshNetwork()
  refreshSystemInfo()
  loadServices()

  backups.value = [
    {
      id: '1',
      name: '系统配置',
      description: '完整系统配置备份',
      created: '2024-01-15',
      size: 1024000,
      type: '配置'
    },
    {
      id: '2',
      name: '数据备份',
      description: '用户数据备份',
      created: '2024-01-14',
      size: 5120000000,
      type: '数据'
    }
  ]
})
</script>

<style scoped>
.system-settings {
  width: 100%;
  height: 100%;
  padding: 32px;
  background: #f9fafb;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
}

.settings-header {
  margin-bottom: 32px;
}

.settings-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 16px;
  color: #6b7280;
}

.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 24px;
  border-bottom: 2px solid #e5e7eb;
}

.tab-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: -2px;
}

.tab-btn:hover {
  color: #1f2937;
}

.tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

.tab-content {
  flex: 1;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
}

.filters {
  display: flex;
  gap: 12px;
}

.search-input {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  width: 250px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.warning {
  background: #f59e0b;
  border-color: #f59e0b;
  color: white;
}

.action-btn.warning:hover {
  background: #d97706;
}

.action-btn.danger {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}

.action-btn.danger:hover {
  background: #dc2626;
}

.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 24px;
}

.setting-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.interfaces-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interface-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
}

.interface-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
}

.interface-status.active {
  background: #10b981;
}

.interface-info {
  flex: 1;
}

.interface-info h4 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.interface-info p {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 2px;
}

.dns-form,
.proxy-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.form-group input[type="text"],
.form-group input[type="password"],
.form-group input[type="number"],
.form-group input[type="time"],
.form-group select {
  padding: 10px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
}

.form-group input[type="checkbox"] {
  width: 16px;
  height: 16px;
  margin-right: 8px;
}

.proxy-fields {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 12px;
}

.system-info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.info-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: flex;
  gap: 16px;
  align-items: center;
}

.info-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.info-content h3 {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
}

.info-content p {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
}

.services-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.service-item {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: grid;
  grid-template-columns: 1fr auto auto;
  gap: 20px;
  align-items: center;
}

.service-info {
  display: flex;
  gap: 12px;
}

.service-status {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 6px;
}

.service-status.running {
  background: #10b981;
}

.service-status.stopped {
  background: #ef4444;
}

.service-details h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.service-details p {
  font-size: 14px;
  color: #6b7280;
}

.service-meta {
  min-width: 100px;
}

.service-enabled {
  font-size: 12px;
  color: #ef4444;
  padding: 4px 8px;
  background: #fee2e2;
  border-radius: 12px;
}

.service-enabled.active {
  color: #10b981;
  background: #d1fae5;
}

.service-actions {
  display: flex;
  gap: 8px;
}

.updates-section {
  display: flex;
  flex-direction: column;
  gap: 24px;
  margin-bottom: 32px;
}

.update-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.update-card.available {
  border: 2px solid #3b82f6;
}

.update-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.update-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.version-badge {
  padding: 6px 12px;
  background: #f3f4f6;
  border-radius: 16px;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
}

.version-badge.new {
  background: #3b82f6;
  color: white;
}

.update-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;
}

.update-info p {
  font-size: 14px;
  color: #6b7280;
}

.update-info ul {
  margin: 8px 0 0 20px;
}

.update-info li {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 4px;
}

.no-updates {
  background: white;
  border-radius: 12px;
  padding: 48px;
  text-align: center;
  color: #10b981;
}

.no-updates svg {
  margin-bottom: 16px;
}

.update-settings,
.backup-schedule {
  background: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.update-settings h3,
.backup-schedule h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 20px;
}

.settings-form,
.schedule-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.schedule-fields {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-top: 12px;
}

.backups-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  margin-bottom: 32px;
}

.backup-item {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 20px;
  align-items: center;
}

.backup-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #8b5cf6 0%, #a78bfa 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.backup-details h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.backup-details p {
  font-size: 14px;
  color: #6b7280;
  margin-bottom: 8px;
}

.backup-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #9ca3af;
}

.backup-actions {
  display: flex;
  gap: 8px;
}
</style>