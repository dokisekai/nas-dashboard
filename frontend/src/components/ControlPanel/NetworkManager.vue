<template>
  <div class="network-manager">
    <!-- 网络类型标签页 -->
    <div class="network-tabs">
      <button
        v-for="tab in networkTabs"
        :key="tab.id"
        :class="['tab-button', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.name }}
      </button>
    </div>

    <!-- 以太网配置 -->
    <div v-if="activeTab === 'ethernet'" class="tab-content">
      <div v-if="loadingInterfaces" class="loading-state">
        <div class="spinner"></div>
        <p>正在加载网络接口...</p>
      </div>

      <div v-else-if="ethernetInterfaces.length > 0" class="network-cards">
        <div
          v-for="iface in ethernetInterfaces"
          :key="iface.name"
          class="network-card"
          :class="{ active: iface.up }"
        >
          <div class="card-header">
            <div class="card-title">
              <ServerIcon class="w-5 h-5" />
              <h3>{{ iface.name }}</h3>
              <span class="interface-status" :class="{ up: iface.up }">
                {{ iface.up ? '已连接' : '未连接' }}
              </span>
            </div>
            <div class="card-actions">
              <button
                @click="toggleInterface(iface.name, iface.up)"
                :class="{ danger: iface.up }"
              >
                {{ iface.up ? '禁用' : '启用' }}
              </button>
            </div>
          </div>

          <div class="card-body">
            <div class="interface-info">
              <div class="info-row">
                <span class="label">MAC地址:</span>
                <span class="value">{{ iface.hardwareAddr }}</span>
              </div>
              <div class="info-row">
                <span class="label">MTU:</span>
                <span class="value">{{ iface.mtu }}</span>
              </div>
              <div class="info-row">
                <span class="label">IP地址:</span>
                <span class="value">{{ getInterfaceIP(iface) }}</span>
              </div>
              <div class="info-row">
                <span class="label">上传速度:</span>
                <span class="value">{{ formatSpeed(iface.sentSpeed) }}</span>
              </div>
              <div class="info-row">
                <span class="label">下载速度:</span>
                <span class="value">{{ formatSpeed(iface.recvSpeed) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <ServerIcon class="w-12 h-12" />
        <p>未检测到以太网接口</p>
      </div>
    </div>

    <!-- Wi-Fi配置 -->
    <div v-if="activeTab === 'wifi'" class="tab-content">
      <div class="wifi-section">
        <div class="section-header">
          <h3>可用的Wi-Fi网络</h3>
          <button class="scan-btn" @click="scanWiFi" :disabled="scanning">
            <ArrowPathIcon v-if="scanning" class="w-4 h-4 spinning" />
            <MagnifyingGlassIcon v-else class="w-4 h-4" />
            {{ scanning ? '扫描中...' : '扫描网络' }}
          </button>
        </div>

        <div v-if="wifiNetworks.length > 0" class="wifi-list">
          <div
            v-for="network in wifiNetworks"
            :key="network.ssid || network.bssid"
            class="wifi-item"
            :class="{ connected: network.connected, connecting: network.connecting }"
          >
            <div class="wifi-info">
              <div class="signal-indicator">
                <SignalIcon class="w-5 h-5" :style="{ opacity: getSignalOpacity(network.signalStrength) }" />
                <div class="signal-bars">
                  <div class="bar" :class="{ active: network.signalStrength >= 25 }"></div>
                  <div class="bar" :class="{ active: network.signalStrength >= 50 }"></div>
                  <div class="bar" :class="{ active: network.signalStrength >= 75 }"></div>
                  <div class="bar" :class="{ active: network.signalStrength >= 90 }"></div>
                </div>
              </div>
              <div class="wifi-details">
                <div class="wifi-name">
                  {{ network.ssid || '(隐藏网络)' }}
                  <span v-if="network.connected" class="connected-badge">已连接</span>
                </div>
                <div class="wifi-meta">
                  <span class="security">{{ getSecurityLabel(network.security) }}</span>
                  <span class="signal-strength">信号: {{ network.signalStrength }}%</span>
                  <span class="channel">频道: {{ network.channel }}</span>
                </div>
              </div>
            </div>
            <button
              @click="connectToWiFi(network)"
              class="connect-btn"
              :disabled="network.connected || network.connecting"
            >
              {{ network.connecting ? '连接中...' : (network.connected ? '已连接' : '连接') }}
            </button>
          </div>
        </div>

        <div v-else-if="!scanning && wifiNetworks.length === 0" class="empty-state">
          <SignalIcon class="w-12 h-12" />
          <p>未找到Wi-Fi网络，请点击"扫描网络"</p>
        </div>

        <div class="hidden-network">
          <button class="hidden-network-btn" @click="showHiddenNetworkDialog = true">
            <PlusIcon class="w-4 h-4" />
            连接隐藏网络
          </button>
        </div>
      </div>

      <!-- 已连接的Wi-Fi状态 -->
      <div v-if="connectedWiFi" class="wifi-status">
        <h3>当前连接</h3>
        <div class="status-card">
          <div class="status-item">
            <span class="label">网络:</span>
            <span class="value">{{ connectedWiFi.ssid }}</span>
          </div>
          <div class="status-item">
            <span class="label">信号强度:</span>
            <span class="value">{{ connectedWiFi.signalStrength }}%</span>
          </div>
          <div class="status-item">
            <span class="label">上传:</span>
            <span class="value">{{ formatSpeed(connectedWiFi.uploadSpeed) }}</span>
          </div>
          <div class="status-item">
            <span class="label">下载:</span>
            <span class="value">{{ formatSpeed(connectedWiFi.downloadSpeed) }}</span>
          </div>
          <div class="status-item">
            <span class="label">IP地址:</span>
            <span class="value">{{ connectedWiFi.ipAddress }}</span>
          </div>
          <button @click="disconnectWiFi" class="disconnect-btn">断开连接</button>
        </div>
      </div>
    </div>

    <!-- DNS配置 -->
    <div v-if="activeTab === 'dns'" class="tab-content">
      <div class="config-section">
        <div class="section-row">
          <div class="config-item full">
            <label>DNS配置方式</label>
            <div class="radio-group">
              <label class="radio">
                <input type="radio" v-model="dns.method" value="auto" @change="saveConfig" />
                <span>自动获得DNS服务器</span>
              </label>
              <label class="radio">
                <input type="radio" v-model="dns.method" value="manual" @change="saveConfig" />
                <span>手动设置DNS服务器</span>
              </label>
            </div>
          </div>
        </div>

        <template v-if="dns.method === 'manual'">
          <div class="section-row">
            <div class="config-item">
              <label>首选DNS</label>
              <input type="text" v-model="dns.primary" @change="saveConfig" placeholder="8.8.8.8" />
            </div>
            <div class="config-item">
              <label>备用DNS</label>
              <input type="text" v-model="dns.secondary" @change="saveConfig" placeholder="8.8.4.4" />
            </div>
          </div>
        </template>
      </div>
    </div>

    <!-- 防火墙配置 -->
    <div v-if="activeTab === 'firewall'" class="tab-content">
      <FirewallManager />
    </div>

    <!-- Wi-Fi连接对话框 -->
    <div v-if="showWiFiDialog" class="modal-overlay" @click="showWiFiDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>连接Wi-Fi</h3>
          <button class="close-btn" @click="showWiFiDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="submitWiFiConnection" class="modal-body">
          <div class="form-group">
            <label>网络名称(SSID)</label>
            <input
              type="text"
              v-model="wifiDialog.ssid"
              placeholder="输入Wi-Fi名称"
              :disabled="wifiDialog.fromScan"
            />
          </div>

          <div class="form-group" v-if="wifiDialog.security !== 'open'">
            <label>密码</label>
            <input
              type="password"
              v-model="wifiDialog.password"
              placeholder="输入Wi-Fi密码"
            />
          </div>

          <div class="form-group">
            <label>安全类型</label>
            <select v-model="wifiDialog.security" :disabled="wifiDialog.fromScan">
              <option value="wpa2">WPA2-Personal</option>
              <option value="wpa3">WPA3-Personal</option>
              <option value="wpawpa2">WPA/WPA2混合</option>
              <option value="wep">WEP (不安全)</option>
              <option value="open">开放网络</option>
            </select>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showWiFiDialog = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              连接
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 隐藏网络对话框 -->
    <div v-if="showHiddenNetworkDialog" class="modal-overlay" @click="showHiddenNetworkDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>连接隐藏网络</h3>
          <button class="close-btn" @click="showHiddenNetworkDialog = false">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <form @submit.prevent="submitHiddenNetworkConnection" class="modal-body">
          <div class="form-group">
            <label>网络名称(SSID)</label>
            <input
              type="text"
              v-model="hiddenNetworkDialog.ssid"
              placeholder="输入隐藏网络的SSID"
            />
          </div>

          <div class="form-group">
            <label>密码</label>
            <input
              type="password"
              v-model="hiddenNetworkDialog.password"
              placeholder="输入Wi-Fi密码"
            />
          </div>

          <div class="form-group">
            <label>安全类型</label>
            <select v-model="hiddenNetworkDialog.security">
              <option value="wpa2">WPA2-Personal</option>
              <option value="wpa3">WPA3-Personal</option>
              <option value="wpawpa2">WPA/WPA2混合</option>
            </select>
          </div>

          <div class="modal-footer">
            <button type="button" class="btn-secondary" @click="showHiddenNetworkDialog = false">
              取消
            </button>
            <button type="submit" class="btn-primary">
              连接
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  ServerIcon,
  SignalIcon,
  GlobeAltIcon,
  ShieldCheckIcon,
  ArrowPathIcon,
  MagnifyingGlassIcon,
  PlusIcon,
  XMarkIcon
} from '@heroicons/vue/24/outline'
import FirewallManager from './FirewallManager.vue'

// 网络标签页
const networkTabs = [
  { id: 'ethernet', name: '以太网', icon: ServerIcon },
  { id: 'wifi', name: 'Wi-Fi', icon: SignalIcon },
  { id: 'dns', name: 'DNS', icon: GlobeAltIcon },
  { id: 'firewall', name: '防火墙', icon: ShieldCheckIcon }
]

const activeTab = ref('ethernet')

// 网络接口
const loadingInterfaces = ref(false)
const ethernetInterfaces = ref<any[]>([])
const wifiInterfaces = ref<any[]>([])

// Wi-Fi网络
const scanning = ref(false)
const wifiNetworks = ref<any[]>([])
const connectedWiFi = ref<any>(null)

// Wi-Fi连接对话框
const showWiFiDialog = ref(false)
const showHiddenNetworkDialog = ref(false)
const wifiDialog = reactive({
  ssid: '',
  password: '',
  security: 'wpa2',
  fromScan: false,
  bssid: ''
})

// 隐藏网络对话框
const hiddenNetworkDialog = reactive({
  ssid: '',
  password: '',
  security: 'wpa2'
})

// DNS配置
const dns = reactive({
  method: 'auto',
  primary: '8.8.8.8',
  secondary: '8.8.4.4'
})

// 方法
const loadNetworkInterfaces = async () => {
  loadingInterfaces.value = true
  try {
    // 获取以太网接口
    const ethResponse = await fetch('/api/network/interfaces/ethernet', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (ethResponse.ok) {
      const ethText = await ethResponse.text()
      console.log('Ethernet response:', ethText)
      const ethData = JSON.parse(ethText)
      ethernetInterfaces.value = ethData || []
    } else {
      console.error('Ethernet API failed:', ethResponse.status, ethResponse.statusText)
      throw new Error('获取以太网接口失败')
    }

    // 获取Wi-Fi接口
    const wifiResponse = await fetch('/api/network/interfaces/wifi', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (wifiResponse.ok) {
      const wifiText = await wifiResponse.text()
      console.log('WiFi response:', wifiText)
      const wifiData = JSON.parse(wifiText)
      wifiInterfaces.value = wifiData || []
    } else {
      console.error('WiFi API failed:', wifiResponse.status, wifiResponse.statusText)
      throw new Error('获取Wi-Fi接口失败')
    }

    console.log('Loaded interfaces:', {
      ethernet: ethernetInterfaces.value.length,
      wifi: wifiInterfaces.value.length
    })

  } catch (error) {
    console.error('Failed to load network interfaces:', error)
    ElMessage.error('获取网络接口失败')
  } finally {
    loadingInterfaces.value = false
  }
}

const toggleInterface = async (interfaceName: string, currentState: boolean) => {
  try {
    const action = currentState ? 'down' : 'up'
    const actionText = currentState ? '禁用' : '启用'

    const response = await fetch(`/api/network/interfaces/${interfaceName}/${action}`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error(`${actionText}接口失败`)
    }

    ElMessage.success(`接口${interfaceName}已${actionText}`)

    // 等待一小段时间后重新加载接口状态
    setTimeout(() => {
      loadNetworkInterfaces()
    }, 500)

  } catch (error: any) {
    console.error('Failed to toggle interface:', error)
    ElMessage.error(error.message || '操作失败')
  }
}

const getInterfaceIP = (interfaceData: any): string => {
  if (!interfaceData.addresses || interfaceData.addresses.length === 0) {
    return '未分配'
  }

  // 优先返回IPv4地址
  const ipAddr = interfaceData.addresses[0]
  if (typeof ipAddr === 'string') {
    // 检查是否是IPv4地址
    if (ipAddr.match(/^\d+\.\d+\.\d+\.\d+$/)) {
      return ipAddr
    }
    // 如果第一个不是IPv4，查找IPv4地址
    for (const addr of interfaceData.addresses) {
      if (typeof addr === 'string' && addr.match(/^\d+\.\d+\.\d+\.\d+$/)) {
        return addr
      }
    }
    // 没有IPv4，返回第一个地址
    return ipAddr
  }
  return ipAddr || '未分配'
}

const formatSpeed = (bytesPerSecond: number): string => {
  if (!bytesPerSecond || bytesPerSecond === 0) return '0 B/s'

  const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
  let value = bytesPerSecond
  let unitIndex = 0

  while (value >= 1024 && unitIndex < units.length - 1) {
    value /= 1024
    unitIndex++
  }

  return `${value.toFixed(1)} ${units[unitIndex]}`
}

const scanWiFi = async () => {
  scanning.value = true
  try {
    const response = await fetch('/api/network/wifi/scan', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (!response.ok) {
      throw new Error('Wi-Fi扫描失败')
    }

    const data = await response.json()
    wifiNetworks.value = data.networks || []

    // 获取当前连接
    connectedWiFi.value = wifiNetworks.value.find(n => n.connected) || null

    console.log('WiFi scan completed, found networks:', wifiNetworks.value.length)
    ElMessage.success(`找到 ${wifiNetworks.value.length} 个可用网络`)

  } catch (error) {
    console.error('WiFi scan failed:', error)
    ElMessage.error('Wi-Fi扫描失败')
  } finally {
    scanning.value = false
  }
}

const connectToWiFi = (network: any) => {
  wifiDialog.ssid = network.ssid
  wifiDialog.security = network.security || 'wpa2'
  wifiDialog.fromScan = true
  wifiDialog.bssid = network.bssid
  wifiDialog.password = ''
  showWiFiDialog.value = true
}

const submitWiFiConnection = async () => {
  try {
    if (wifiDialog.security !== 'open' && !wifiDialog.password) {
      ElMessage.warning('请输入Wi-Fi密码')
      return
    }

    // 标记为连接中
    const network = wifiNetworks.value.find(n => n.ssid === wifiDialog.ssid)
    if (network) {
      network.connecting = true
    }

    showWiFiDialog.value = false

    // 调用真实的连接API
    const response = await fetch('/api/network/wifi/connect', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        ssid: wifiDialog.ssid,
        password: wifiDialog.password,
        security: wifiDialog.security,
        bssid: wifiDialog.bssid,
        isHidden: !wifiDialog.fromScan
      })
    })

    if (!response.ok) {
      throw new Error('Wi-Fi连接失败')
    }

    if (network) {
      network.connecting = false
      network.connected = true
      connectedWiFi.value = network
    }

    ElMessage.success(`已连接到 ${wifiDialog.ssid}`)

    // 重置表单
    wifiDialog.ssid = ''
    wifiDialog.password = ''
    wifiDialog.fromScan = false

    // 重新扫描以更新状态
    await scanWiFi()

  } catch (error) {
    console.error('WiFi connection failed:', error)
    ElMessage.error('Wi-Fi连接失败')

    const network = wifiNetworks.value.find(n => n.ssid === wifiDialog.ssid)
    if (network) {
      network.connecting = false
    }
  }
}

const submitHiddenNetworkConnection = async () => {
  try {
    if (!hiddenNetworkDialog.ssid) {
      ElMessage.warning('请输入隐藏网络的SSID')
      return
    }

    if (hiddenNetworkDialog.security !== 'open' && !hiddenNetworkDialog.password) {
      ElMessage.warning('请输入Wi-Fi密码')
      return
    }

    showHiddenNetworkDialog.value = false

    // 调用真实的连接API
    const response = await fetch('/api/network/wifi/connect', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        ssid: hiddenNetworkDialog.ssid,
        password: hiddenNetworkDialog.password,
        security: hiddenNetworkDialog.security,
        isHidden: true
      })
    })

    if (!response.ok) {
      throw new Error('隐藏网络连接失败')
    }

    ElMessage.success(`已连接到隐藏网络 ${hiddenNetworkDialog.ssid}`)

    // 重置表单
    hiddenNetworkDialog.ssid = ''
    hiddenNetworkDialog.password = ''

    // 重新扫描以更新状态
    await scanWiFi()

  } catch (error) {
    console.error('Hidden network connection failed:', error)
    ElMessage.error('隐藏网络连接失败')
  }
}

const disconnectWiFi = async () => {
  try {
    const response = await fetch('/api/network/wifi/disconnect', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error('断开Wi-Fi连接失败')
    }

    ElMessage.success('已断开Wi-Fi连接')
    connectedWiFi.value = null
    await scanWiFi()

  } catch (error) {
    console.error('WiFi disconnect failed:', error)
    ElMessage.error('断开连接失败')
  }
}

const saveConfig = async () => {
  try {
    const response = await fetch('/api/network/dns', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        method: dns.method,
        primary: dns.primary,
        secondary: dns.secondary
      })
    })

    if (!response.ok) {
      throw new Error('DNS配置保存失败')
    }

    ElMessage.success('DNS配置已保存')
  } catch (error) {
    console.error('Failed to save DNS config:', error)
    ElMessage.error('DNS配置保存失败')
  }
}

const loadDNSConfig = async () => {
  try {
    const response = await fetch('/api/network/dns', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    })

    if (response.ok) {
      const responseText = await response.text()
      console.log('DNS response:', responseText)
      const config = JSON.parse(responseText)
      dns.method = config.method || 'auto'
      dns.primary = config.primary || '8.8.8.8'
      dns.secondary = config.secondary || '8.8.4.4'
    } else {
      console.error('DNS API failed:', response.status, response.statusText)
    }
  } catch (error) {
    console.error('Failed to load DNS config:', error)
    // 保持默认值
  }
}

// 辅助方法
const getSignalOpacity = (strength: number): string => {
  if (strength >= 75) return '1'
  if (strength >= 50) return '0.8'
  if (strength >= 25) return '0.6'
  return '0.4'
}

const getSecurityLabel = (security: string): string => {
  const labels: Record<string, string> = {
    'WPA3': 'WPA3',
    'WPA2': 'WPA2',
    'WPAPSK': 'WPA',
    'wpa3': 'WPA3',
    'wpa2': 'WPA2',
    'wpawpa2': 'WPA/WPA2',
    'wep': 'WEP',
    'open': '开放',
    '': '开放'
  }
  return labels[security] || security || '开放'
}

console.log('NetworkManager component setup completed')

// 生命周期
onMounted(() => {
  console.log('NetworkManager mounted, activeTab:', activeTab.value)

  // 延迟加载网络接口，确保组件完全渲染
  setTimeout(async () => {
    await loadNetworkInterfaces()

    // 加载完网络接口后，检查是否有Wi-Fi接口，然后扫描
    setTimeout(async () => {
      if (wifiInterfaces.value.length > 0) {
        await scanWiFi()
      }
    }, 200)
  }, 100)

  // 加载DNS配置
  loadDNSConfig()

  console.log('NetworkManager initialization scheduled')
})
</script>

<style scoped lang="scss">
.network-manager {
  width: 100%;
  padding: 16px;
}

.network-tabs {
  display: flex;
  gap: 0;
  border-bottom: 1px solid #e5e7eb;
  margin-bottom: 24px;
}

.tab-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: none;
  border: none;
  border-bottom: 2px solid transparent;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;

  &:hover {
    color: #374151;
    background: #f9fafb;
  }

  &.active {
    color: #2563eb;
    border-bottom-color: #2563eb;
  }
}

.tab-content {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #6b7280;

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid #e5e7eb;
    border-top-color: #3b82f6;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 12px;
  }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.network-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(400px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.network-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  &.active {
    border-left: 4px solid #10b981;
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #f3f4f6;
  background: #fafbfc;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 12px;

  svg {
    color: #3b82f6;
  }

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.interface-status {
  font-size: 12px;
  padding: 4px 12px;
  border-radius: 12px;
  background: #f3f4f6;
  color: #6b7280;

  &.up {
    background: #d1fae5;
    color: #065f46;
  }
}

.card-actions {
  display: flex;
  gap: 8px;

  button {
    padding: 6px 12px;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    background: white;
    color: #4b5563;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 13px;

    &:hover {
      background: #f9fafb;
      border-color: #9ca3af;
    }

    &.danger:hover {
      background: #fef2f2;
      border-color: #ef4444;
      color: #ef4444;
    }
  }
}

.card-body {
  padding: 16px;
}

.interface-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;

  .label {
    color: #6b7280;
    font-weight: 500;
  }

  .value {
    color: #1f2937;
    font-family: monospace;
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: #9ca3af;

  svg {
    margin-bottom: 12px;
    opacity: 0.5;
  }

  p {
    font-size: 14px;
  }
}

.wifi-section {
  margin-bottom: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  h3 {
    font-size: 16px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.scan-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;

  &:hover:not(:disabled) {
    background: #2563eb;
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .spinning {
    animation: spin 1s linear infinite;
  }
}

.wifi-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.wifi-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  transition: all 0.2s;

  &:hover {
    border-color: #d1d5db;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  }

  &.connected {
    background: #f0fdf4;
    border-color: #22c55e;
  }

  &.connecting {
    opacity: 0.6;
  }
}

.wifi-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;

  svg {
    color: #3b82f6;
  }
}

.signal-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.signal-bars {
  display: flex;
  align-items: flex-end;
  gap: 2px;
  height: 16px;

  .bar {
    width: 4px;
    background: #d1d5db;
    border-radius: 2px;
    transition: all 0.2s;

    &:nth-child(1) {
      height: 4px;
    }

    &:nth-child(2) {
      height: 8px;
    }

    &:nth-child(3) {
      height: 12px;
    }

    &:nth-child(4) {
      height: 16px;
    }

    &.active {
      background: #22c55e;

      &:nth-child(1) {
        background: #22c55e;
      }

      &:nth-child(2) {
        background: #22c55e;
      }

      &:nth-child(3) {
        background: #22c55e;
      }

      &:nth-child(4) {
        background: #22c55e;
      }
    }
  }
}

.wifi-details {
  flex: 1;
}

.wifi-name {
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.connected-badge {
  font-size: 11px;
  padding: 2px 8px;
  background: #22c55e;
  color: white;
  border-radius: 12px;
}

.wifi-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #6b7280;
}

.connect-btn {
  padding: 6px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 13px;

  &:hover:not(:disabled) {
    background: #2563eb;
  }

  &:disabled {
    background: #9ca3af;
    cursor: not-allowed;
  }
}

.hidden-network {
  margin-top: 16px;
}

.hidden-network-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px;
  background: white;
  border: 1px dashed #d1d5db;
  border-radius: 8px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    border-color: #3b82f6;
    color: #3b82f6;
    background: #f0f9ff;
  }
}

.wifi-status {
  margin-top: 24px;
}

.wifi-status h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.status-card {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  padding: 16px;
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border-radius: 8px;
  border: 1px solid #bae6fd;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  .label {
    font-size: 12px;
    color: #0369a1;
  }

  .value {
    font-size: 14px;
    font-weight: 600;
    color: #0c4a6e;
  }
}

.disconnect-btn {
  grid-column: 1 / -1;
  padding: 8px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  margin-top: 8px;
  transition: all 0.2s;

  &:hover {
    background: #dc2626;
  }
}

.config-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;

  .config-item.full {
    grid-column: 1 / -1;
  }
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-size:  13px;
    font-weight: 500;
    color: #374151;
  }

  input, select {
    padding: 10px 12px;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 14px;

    &:focus {
      outline: none;
      border-color: #2563eb;
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
    }
  }
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.radio {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;

  input[type="radio"] {
    width: 16px;
    height: 16px;
    cursor: pointer;
  }

  span {
    font-size: 14px;
    color: #374151;
  }
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  padding: 24px;
  min-width: 400px;
  max-width: 90%;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #1f2937;
    margin: 0;
  }
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: #f3f4f6;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: #e5e7eb;
  }
}

.modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;

  label {
    font-size: 13px;
    font-weight: 500;
    color: #374151;
  }

  input, select {
    padding: 10px 12px;
    border: 1px solid #d1d5db;
    border-radius: 6px;
    font-size: 14px;

    &:focus {
      outline: none;
      border-color: #2563eb;
      box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
    }
  }
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 8px;

  button {
    padding: 8px 16px;
    border-radius: 6px;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;
  }
}

.btn-secondary {
  background: white;
  border: 1px solid #d1d5db;
  color: #4b5563;

  &:hover {
    background: #f9fafb;
    border-color: #9ca3af;
  }
}

.btn-primary {
  background: #3b82f6;
  border: 1px solid #3b82f6;
  color: white;

  &:hover {
    background: #2563eb;
    border-color: #2563eb;
  }
}
</style>