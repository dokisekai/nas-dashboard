<template>
  <div class="network-manager" :class="{ 'embedded-mode': embeddedMode }">
    <!-- 调试信息 -->
    <div v-if="props.embeddedMode" class="debug-info" style="background: #f0f0f0; padding: 12px; margin: 8px 0; border-radius: 8px; font-size: 12px; border: 2px solid #007bff;">
      📡 <strong>网络管理组件调试信息</strong>
      <div style="margin-top: 8px;">
        • 接口数量: {{ interfaces.length }}<br>
        • 加载状态: {{ loading ? '加载中' : '完成' }}<br>
        • 错误状态: {{ error || '无错误' }}<br>
        • 模态框状态:
          - 配置({{ showInterfaceConfig ? '✅ 显示' : '❌ 隐藏' }})
          - 代理({{ showProxyConfig ? '✅ 显示' : '❌ 隐藏' }})
          - DNS({{ showDnsModal ? '✅ 显示' : '❌ 隐藏' }})
          - WiFi({{ showWifiSetup ? '✅ 显示' : '❌ 隐藏' }})
      </div>
      <div style="margin-top: 8px; padding: 8px; background: white; border-radius: 4px;">
        <strong>按钮测试：</strong> 已点击 {{ buttonClicks }} 次<br>
        <small>点击任意接口的"配置"或"重启"按钮，此计数应该增加</small>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading && !interfaces.length" class="loading-state">
      <div class="spinner"></div>
      <p>正在加载网络信息...</p>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="error-alert">
      <XCircleIcon class="w-5 h-5" />
      <span>{{ error }}</span>
      <button @click="error = ''">×</button>
    </div>

    <!-- 头部 -->
    <div class="nm-header" v-if="!embeddedMode">
      <h1>网络管理</h1>
      <p class="subtitle">管理网络连接和网络设置</p>
    </div>

    <!-- 网络状态概览 -->
    <div class="nm-overview" v-if="!embeddedMode">
      <div class="overview-card">
        <div class="card-icon active">
          <SignalIcon class="w-6 h-6" />
        </div>
        <div class="card-info">
          <div class="card-title">网络状态</div>
          <div class="card-value">{{ activeInterface ? '已连接' : '未连接' }}</div>
          <div class="card-subtitle">{{ activeInterface ? `${activeInterface.name} - ${activeInterface.ip}` : '无活动连接' }}</div>
        </div>
        <div class="card-speed">
          <div class="speed-upload">
            <ArrowUpIcon class="w-4 h-4" />
            <span>{{ formatSpeed(uploadSpeed) }}</span>
          </div>
          <div class="speed-download">
            <ArrowDownIcon class="w-4 h-4" />
            <span>{{ formatSpeed(downloadSpeed) }}</span>
          </div>
        </div>
      </div>

      <div class="overview-card" @click="showWifiSetup = true" style="cursor: pointer">
        <div class="card-icon">
          <WifiIcon class="w-6 h-6" />
        </div>
        <div class="card-info">
          <div class="card-title">无线网络</div>
          <div class="card-value">{{ wifiConnected ? '已连接' : '未连接' }}</div>
          <div class="card-subtitle">{{ wifiConnected ? wifiConnected.ssid : '点击设置无线连接' }}</div>
        </div>
        <button class="card-action">
          设置
        </button>
      </div>

      <div class="overview-card" @click="showDnsModal = true" style="cursor: pointer">
        <div class="card-icon">
          <ServerIcon class="w-6 h-6" />
        </div>
        <div class="card-info">
          <div class="card-title">DNS设置</div>
          <div class="card-value">{{ dnsConfig.method === 'auto' ? '自动' : '手动' }}</div>
          <div class="card-subtitle">{{ dnsConfig.primary || '未设置' }}</div>
        </div>
        <button class="card-action">
          配置
        </button>
      </div>
    </div>

    <!-- 主要内容区 -->
    <div class="nm-content">
      <!-- 左侧：网络接口 -->
      <div class="nm-interfaces">
        <div class="section-header">
          <h2>网络接口</h2>
          <button class="action-btn" @click="refreshInterfaces" :disabled="loading">
            <ArrowPathIcon class="w-4 h-4" :class="{ spinning: loading }" />
            刷新
          </button>
        </div>

        <div class="interface-list">
          <div
            v-for="iface in interfaces"
            :key="iface.name"
            class="interface-item"
            :class="{ active: iface.active }"
          >
            <div class="interface-header">
              <div class="interface-info">
                <div class="interface-icon" :class="{ active: iface.active }">
                  <ServerIcon class="w-5 h-5" />
                </div>
                <div>
                  <div class="interface-name">{{ iface.name }}</div>
                  <div class="interface-type">{{ iface.type }}</div>
                </div>
              </div>

              <div class="interface-status">
                <div class="status-indicator" :class="{ active: iface.active }"></div>
                <span>{{ iface.active ? '已连接' : '未连接' }}</span>
              </div>
            </div>

            <div v-if="iface.active && iface.ip" class="interface-details">
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="detail-label">IP地址:</span>
                  <span class="detail-value">{{ iface.ip }}</span>
                </div>
                <div class="detail-item" v-if="iface.netmask">
                  <span class="detail-label">子网掩码:</span>
                  <span class="detail-value">{{ iface.netmask }}</span>
                </div>
                <div class="detail-item" v-if="iface.mac">
                  <span class="detail-label">MAC地址:</span>
                  <span class="detail-value">{{ iface.mac }}</span>
                </div>
                <div class="detail-item" v-if="iface.gateway">
                  <span class="detail-label">网关:</span>
                  <span class="detail-value">{{ iface.gateway }}</span>
                </div>
                <div class="detail-item" v-if="iface.dns">
                  <span class="detail-label">DNS服务器:</span>
                  <span class="detail-value">{{ iface.dns }}</span>
                </div>
                <div class="detail-item" v-if="iface.speed">
                  <span class="detail-label">速率:</span>
                  <span class="detail-value">{{ iface.speed }}</span>
                </div>
              </div>

              <div class="interface-actions">
                <button class="action-btn" @click="toggleInterface(iface)">
                  <XCircleIcon v-if="iface.active" class="w-4 h-4" />
                  <SignalIcon v-else class="w-4 h-4" />
                  {{ iface.active ? '断开' : '连接' }}
                </button>
												<button class="action-btn" @click="openInterfaceConfig(iface.name)">
													<CogIcon class="w-4 h-4" />
													配置
												</button>
												<button class="action-btn" @click="restartInterface(iface.name)" :disabled="loading">
													<ArrowPathIcon class="w-4 h-4" :class="{ spinning: loading }" />
													重启
												</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="!interfaces.length" class="empty-state">
          <ServerIcon class="w-12 h-12" />
          <p>未发现网络接口</p>
          <button class="action-btn primary" @click="refreshInterfaces">
            重新扫描
          </button>
        </div>
      </div>

      <!-- 右侧：网络统计和工具 -->
      <div class="nm-tools">
        <!-- 网络统计 -->
        <div class="tool-section">
          <h3>网络统计</h3>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-icon upload">
                <ArrowUpIcon class="w-4 h-4" />
              </div>
              <div class="stat-info">
                <div class="stat-label">今日上传</div>
                <div class="stat-value">{{ formatBytes(todayUpload) }}</div>
              </div>
            </div>
            <div class="stat-item">
              <div class="stat-icon download">
                <ArrowDownIcon class="w-4 h-4" />
              </div>
              <div class="stat-info">
                <div class="stat-label">今日下载</div>
                <div class="stat-value">{{ formatBytes(todayDownload) }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 快速操作 -->
        <div class="tool-section">
          <h3>快速操作</h3>
          <div class="quick-actions">
            <button class="quick-action" @click="showWifiSetup = true">
              <WifiIcon class="w-5 h-5" />
              <span>WiFi设置</span>
            </button>
            <button class="quick-action" @click="showDnsModal = true">
              <ServerIcon class="w-5 h-5" />
              <span>DNS配置</span>
            </button>
            <button class="quick-action" @click="refreshInterfaces">
              <ArrowPathIcon class="w-5 h-5" />
              <span>刷新网络</span>
            </button>
          </div>
        </div>

        <!-- 网络诊断 -->
        <div class="tool-section">
          <h3>网络诊断</h3>
          <button class="action-btn full" @click="runNetworkDiagnostics">
            <MagnifyingGlassIcon class="w-4 h-4" />
            运行诊断
          </button>
        </div>

        <!-- 代理配置 -->
        <div class="tool-section">
          <h3>代理服务器</h3>
          <button class="action-btn full" @click="openProxyConfig">
            <CogIcon class="w-4 h-4" />
            配置代理
          </button>
        </div>
      </div>
    </div>

    <!-- WiFi设置模态框 -->
    <Transition name="fade">
      <div v-if="showWifiSetup" class="modal-overlay" @click.self="showWifiSetup = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>无线网络设置</h3>
            <button @click="showWifiSetup = false">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="modal-body">
            <div class="wifi-scan">
              <div class="scan-header">
                <h4>可用网络</h4>
                <button class="action-btn" @click="scanWifiNetworks" :disabled="loading">
                  <MagnifyingGlassIcon class="w-4 h-4" />
                  扫描
                </button>
              </div>
              <div class="network-list" v-if="wifiNetworks.length > 0">
                <div
                  v-for="network in wifiNetworks"
                  :key="network.ssid"
                  class="network-item"
                  :class="{ connected: network.connected, secured: network.security && network.security !== 'Open' && network.security !== '开放' }"
                  @click="connectToWifi(network)"
                >
                  <SignalIcon class="w-5 h-5" />
                  <div class="network-info">
                    <div class="network-ssid">{{ network.ssid }}</div>
                    <div class="network-details">
                      <span>{{ network.security }}</span>
                      <span>信号: {{ getSignalText(network.signalStrength) }}</span>
                    </div>
                  </div>
                  <LockClosedIcon v-if="network.security && network.security !== 'Open' && network.security !== '开放'" class="w-4 h-4" />
                  <div v-if="network.connected" class="connected-badge">已连接</div>
                </div>
              </div>
              <div v-else class="empty-networks">
                <p>暂无可用网络，点击扫描按钮搜索附近的WiFi网络</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>

    <!-- DNS设置模态框 -->
    <Transition name="fade">
      <div v-if="showDnsModal" class="modal-overlay" @click.self="showDnsModal = false">
        <div class="modal-content">
          <div class="modal-header">
            <h3>DNS服务器配置</h3>
            <button @click="showDnsModal = false">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="modal-body">
            <div class="dns-config">
              <div class="form-group">
                <label>DNS配置方式</label>
                <div class="radio-group">
                  <label class="radio-item">
                    <input type="radio" v-model="dnsConfig.method" value="auto" />
                    <span>自动获取 (DHCP)</span>
                  </label>
                  <label class="radio-item">
                    <input type="radio" v-model="dnsConfig.method" value="manual" />
                    <span>手动配置</span>
                  </label>
                </div>
              </div>

              <div v-if="dnsConfig.method === 'manual'" class="manual-dns">
                <div class="form-group">
                  <label>主DNS服务器</label>
                  <input v-model="dnsConfig.primary" type="text" class="form-input" placeholder="192.168.1.1" />
                </div>
                <div class="form-group">
                  <label>备用DNS服务器</label>
                  <input v-model="dnsConfig.secondary" type="text" class="form-input" placeholder="8.8.8.8" />
                </div>
                <div class="quick-dns">
                  <label>常用DNS：</label>
                  <button @click="setQuickDNS('8.8.8.8', '8.8.4.4')" class="dns-btn">Google DNS</button>
                  <button @click="setQuickDNS('1.1.1.1', '1.0.0.1')" class="dns-btn">Cloudflare DNS</button>
                  <button @click="setQuickDNS('208.67.222.222', '208.67.220.220')" class="dns-btn">OpenDNS</button>
                </div>
              </div>

              <div class="modal-footer">
                <button class="action-btn" @click="showDnsModal = false">取消</button>
                <button class="action-btn primary" @click="applyDNSConfig">应用设置</button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>

	    <!-- 接口配置模态框 -->
	    <Transition name="fade">
	      <div v-if="showInterfaceConfig" class="modal-overlay" @click.self="showInterfaceConfig = false">
	        <div class="modal-content">
	          <div class="modal-header">
	            <h3>网络接口配置 - {{ selectedInterface }}</h3>
	            <button @click="showInterfaceConfig = false">
	              <XMarkIcon class="w-5 h-5" />
	            </button>
	          </div>
	          <div class="modal-body">
	            <div class="form-section">
	              <h4>IPv4配置方式</h4>
	              <div class="radio-group">
	                <label class="radio-item">
	                  <input type="radio" v-model="interfaceForm.ipv4Method" value="dhcp" />
	                  <span>动态IP (DHCP)</span>
	                </label>
	                <label class="radio-item">
	                  <input type="radio" v-model="interfaceForm.ipv4Method" value="static" />
	                  <span>静态IP</span>
	                </label>
	              </div>
	            </div>

	            <div v-if="interfaceForm.ipv4Method === 'static'" class="form-section">
	              <h4>静态IP设置</h4>
	              <div class="form-group">
	                <label>IP地址</label>
	                <input v-model="interfaceForm.ipAddress" type="text" class="form-input" placeholder="192.168.1.100" />
	              </div>
	              <div class="form-group">
	                <label>子网掩码</label>
	                <input v-model="interfaceForm.netmask" type="text" class="form-input" placeholder="255.255.255.0" />
	              </div>
	              <div class="form-group">
	                <label>默认网关</label>
	                <input v-model="interfaceForm.gateway" type="text" class="form-input" placeholder="192.168.1.1" />
	              </div>
	            </div>

	            <div class="form-section">
	              <h4>高级设置</h4>
	              <div class="form-group">
	                <label>MTU大小</label>
	                <input v-model="interfaceForm.mtu" type="number" class="form-input" placeholder="1500" min="576" max="9000" />
	                <small>建议值：1500（标准以太网）</small>
	              </div>
	            </div>

	            <div class="modal-footer">
	              <button class="action-btn" @click="showInterfaceConfig = false">取消</button>
	              <button class="action-btn primary" @click="saveInterfaceConfig" :disabled="loading">保存配置</button>
	            </div>
	          </div>
	        </div>
	      </div>
	    </Transition>

	    <!-- 代理配置模态框 -->
	    <Transition name="fade">
	      <div v-if="showProxyConfig" class="modal-overlay" @click.self="showProxyConfig = false">
	        <div class="modal-content">
	          <div class="modal-header">
	            <h3>代理服务器配置</h3>
	            <button @click="showProxyConfig = false">
	              <XMarkIcon class="w-5 h-5" />
	            </button>
	          </div>
	          <div class="modal-body">
	            <div class="form-section">
	              <div class="form-group">
	                <label class="checkbox-item">
	                  <input type="checkbox" v-model="proxyConfig.enabled" />
	                  <span>启用代理</span>
	                </label>
	              </div>

	              <div v-if="proxyConfig.enabled">
	                <div class="form-group">
	                  <label>代理类型</label>
	                  <select v-model="proxyConfig.type" class="form-input">
	                    <option value="http">HTTP</option>
	                    <option value="https">HTTPS</option>
	                    <option value="socks4">SOCKS4</option>
	                    <option value="socks5">SOCKS5</option>
	                  </select>
	                </div>

	                <div class="form-group">
	                  <label>代理服务器地址</label>
	                  <input v-model="proxyConfig.server" type="text" class="form-input" placeholder="proxy.example.com" />
	                </div>

	                <div class="form-group">
	                  <label>端口</label>
	                  <input v-model.number="proxyConfig.port" type="number" class="form-input" placeholder="8080" min="1" max="65535" />
	                </div>

	                <div class="form-group">
	                  <label>用户名 (可选)</label>
	                  <input v-model="proxyConfig.username" type="text" class="form-input" placeholder="username" />
	                </div>

	                <div class="form-group">
	                  <label>密码 (可选)</label>
	                  <input v-model="proxyConfig.password" type="password" class="form-input" placeholder="password" />
	                </div>

	                <div class="form-group">
	                  <label>绕过代理的地址</label>
	                  <input v-model="proxyConfig.bypassList" type="text" class="form-input" placeholder="localhost,127.0.0.1,*.local" />
	                  <small>使用逗号分隔多个地址</small>
	                </div>
	              </div>
	            </div>

	            <div class="modal-footer">
	              <button class="action-btn" @click="showProxyConfig = false">取消</button>
	              <button class="action-btn primary" @click="saveProxyConfig" :disabled="loading">保存配置</button>
	            </div>
	          </div>
	        </div>
	      </div>
	    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  SignalIcon,
  WifiIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  ArrowPathIcon,
  ServerIcon,
  XCircleIcon,
  XMarkIcon,
  MagnifyingGlassIcon,
  LockClosedIcon,
  CogIcon
} from '@heroicons/vue/24/outline'
import { networkApi, wifiApi, dnsApi, networkUtils, interfaceConfigApi, pppoeConfigApi, proxyConfigApi, networkConfigUtils } from '../api'
import type { NetworkInterface, WiFiNetwork, DNSConfig } from '../api/network'
import type { InterfaceConfig, InterfaceConfigRequest, ProxyConfig } from '../api/interface_config'

// Props
interface Props {
  embeddedMode?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  embeddedMode: false
})

// 状态管理
const uploadSpeed = ref(0)
const downloadSpeed = ref(0)
const showWifiSetup = ref(false)
const showDnsModal = ref(false)
const showInterfaceConfig = ref(false)
const showProxyConfig = ref(false)
const loading = ref(false)
const error = ref('')

// 调试状态
const buttonClicks = ref(0)

// 网络接口数据
const interfaces = ref<NetworkInterface[]>([])

// WiFi网络数据
const wifiNetworks = ref<WiFiNetwork[]>([])

// DNS配置
const dnsConfig = ref<DNSConfig>({
  method: 'auto',
  primary: '192.168.1.1',
  secondary: '8.8.8.8'
})

// 接口配置
const selectedInterface = ref<string>('')
const interfaceConfig = ref<InterfaceConfig>({
  name: '',
  ipv4Method: 'dhcp',
  ipv6Method: 'auto',
  dhcp: true,
  autoConnect: true,
  enabled: true,
  mtu: 1500
})

// 接口配置表单
const interfaceForm = ref<InterfaceConfigRequest>({
  ipv4Method: 'dhcp',
  ipAddress: '',
  netmask: '255.255.255.0',
  gateway: '',
  mtu: 1500
})

// 代理配置
const proxyConfig = ref<ProxyConfig>({
  enabled: false,
  type: 'http',
  server: '',
  port: 8080,
  bypassList: ['localhost', '127.0.0.1']
})

// 流量统计
const todayUpload = ref(0)
const todayDownload = ref(0)

// 计算属性
const activeInterface = computed(() => interfaces.value.find(i => i.active))

const wifiConnected = computed(() => wifiNetworks.value.find(n => n.connected))

// 刷新间隔
let refreshInterval: number

// 方法
const formatSpeed = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 KB/s'
  return networkUtils.formatSpeed(bytes)
}

const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 GB'
  return networkUtils.formatBytes(bytes)
}

const getSignalText = (strength: number): string => {
  if (strength >= 70) return '强'
  if (strength >= 40) return '中'
  return '弱'
}

// 加载网络接口
const loadInterfaces = async () => {
  try {
    loading.value = true
    error.value = ''

    console.log('[NetworkManager] Fetching network interfaces...')
    const response = await networkApi.getInterfaces() as any
    console.log('[NetworkManager] API response:', response)

    if (response && response.interfaces) {
      interfaces.value = response.interfaces.map((iface: any) => ({
        name: iface.name,
        type: iface.type === 'ethernet' ? '以太网' : iface.type === 'wireless' ? '无线' : iface.type,
        active: iface.operstate === 'up' || iface.active || false,
        ip: iface.ip_address || iface.ip || '',
        netmask: iface.netmask || '',
        mac: iface.mac || iface.address || '',
        gateway: iface.gateway || '',
        dns: iface.dns || '',
        speed: iface.speed || '',
        tx_bytes: iface.tx_bytes || 0,
        rx_bytes: iface.rx_bytes || 0
      }))

      console.log('[NetworkManager] Processed interfaces:', interfaces.value.length, 'interfaces')

      // 计算总速度
      const totalTx = interfaces.value.reduce((sum, iface) => sum + (iface.tx_bytes || 0), 0)
      const totalRx = interfaces.value.reduce((sum, iface) => sum + (iface.rx_bytes || 0), 0)

      if (totalTx > 0) uploadSpeed.value = totalTx
      if (totalRx > 0) downloadSpeed.value = totalRx

      console.log('[NetworkManager] Upload speed:', uploadSpeed.value, 'Download speed:', downloadSpeed.value)
    } else if (Array.isArray(response)) {
      interfaces.value = response.map((iface: any) => ({
        name: iface.name,
        type: iface.type === 'ethernet' ? '以太网' : iface.type === 'wireless' ? '无线' : iface.type,
        active: iface.operstate === 'up' || iface.active || false,
        ip: iface.ip_address || iface.ip || '',
        netmask: iface.netmask || '',
        mac: iface.mac || iface.address || '',
        gateway: iface.gateway || '',
        dns: iface.dns || '',
        speed: iface.speed || ''
      }))
      console.log('[NetworkManager] Array format response processed')
    } else {
      console.warn('[NetworkManager] Unexpected response format:', response)
    }
  } catch (err: any) {
    console.error('[NetworkManager] Failed to load network interfaces:', err)
    error.value = '加载网络接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
    console.log('[NetworkManager] Loading complete, error state:', !!error.value)
  }
}

// 扫描WiFi网络
const scanWifiNetworks = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await wifiApi.scanNetworks()

    if (response && response.networks) {
      wifiNetworks.value = response.networks.map((network: any) => ({
        ssid: network.ssid,
        bssid: network.bssid || '',
        security: network.security || network.encryption || '未知',
        signalStrength: network.signal_strength || 50,
        channel: network.channel || 1,
        connected: network.connected || false,
        connecting: network.connecting || false
      }))
    } else if (Array.isArray(response)) {
      wifiNetworks.value = response.map((network: any) => ({
        ssid: network.ssid,
        bssid: network.bssid || '',
        security: network.security || network.encryption || '未知',
        signalStrength: network.signal_strength || 50,
        channel: network.channel || 1,
        connected: network.connected || false,
        connecting: network.connecting || false
      }))
    }
  } catch (err: any) {
    console.error('Failed to scan WiFi networks:', err)
    error.value = '扫描WiFi网络失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 连接WiFi网络
const connectToWifi = async (network: WiFiNetwork) => {
  try {
    if (network.connected) {
      // 已经连接，询问是否断开
      if (confirm(`确定要断开 ${network.ssid} 吗？`)) {
        loading.value = true
        await wifiApi.disconnect()
        await scanWifiNetworks()
        await loadInterfaces()
      }
      return
    }

    if (network.security && network.security !== 'Open' && network.security !== '开放' && network.security !== '未知') {
      // 需要密码
      const password = prompt(`请输入 ${network.ssid} 的WiFi密码:`)
      if (!password) return

      loading.value = true
      await wifiApi.connect({
        ssid: network.ssid,
        password: password,
        bssid: network.bssid,
        security: network.security
      })
    } else {
      // 开放网络
      loading.value = true
      await wifiApi.connect({
        ssid: network.ssid,
        bssid: network.bssid
      })
    }

    // 连接成功后刷新
    await new Promise(resolve => setTimeout(resolve, 2000)) // 等待连接建立
    await scanWifiNetworks()
    await loadInterfaces()
    showWifiSetup.value = false
  } catch (err: any) {
    console.error('Failed to connect to WiFi:', err)
    error.value = '连接WiFi失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 切换接口状态
const toggleInterface = async (iface: NetworkInterface) => {
  try {
    loading.value = true

    if (iface.active) {
      await networkApi.controlInterface(iface.name, 'down')
    } else {
      await networkApi.controlInterface(iface.name, 'up')
    }

    await loadInterfaces()
  } catch (err: any) {
    console.error('Failed to toggle interface:', err)
    error.value = '控制接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 刷新网络接口
const refreshInterfaces = async () => {
  await loadInterfaces()
}

// 设置快速DNS
const setQuickDNS = (primary: string, secondary: string) => {
  dnsConfig.value.primary = primary
  dnsConfig.value.secondary = secondary
}

// 应用DNS配置
const applyDNSConfig = async () => {
  try {
    loading.value = true
    await dnsApi.setConfig(dnsConfig.value)
    showDnsModal.value = false
  } catch (err: any) {
    console.error('Failed to set DNS config:', err)
    error.value = '设置DNS失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 运行网络诊断
const runNetworkDiagnostics = () => {
  alert('网络诊断功能开发中...')
}

// 接口配置相关功能
const openInterfaceConfig = async (interfaceName: string) => {
  try {
    console.log('[NetworkManager] openInterfaceConfig called for:', interfaceName)
    buttonClicks.value++
    selectedInterface.value = interfaceName
    loading.value = true
    error.value = ''

    // 获取接口当前配置
    console.log('[NetworkManager] Fetching interface config from API...')
    const response = await interfaceConfigApi.getConfig(interfaceName) as any
    console.log('[NetworkManager] Interface config response:', response)

    interfaceConfig.value = response

    // 设置表单值
    interfaceForm.value = {
      ipv4Method: response.ipv4Method || 'dhcp',
      ipAddress: response.ipAddress || '',
      netmask: response.netmask || '255.255.255.0',
      gateway: response.gateway || '',
      mtu: response.mtu || 1500
    }

    console.log('[NetworkManager] Opening interface config modal, showInterfaceConfig:', true)
    showInterfaceConfig.value = true
    console.log('[NetworkManager] Modal should be visible now')
  } catch (err: any) {
    console.error('[NetworkManager] Failed to load interface config:', err)
    error.value = '加载接口配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 保存接口配置
const saveInterfaceConfig = async () => {
  try {
    // 验证配置
    const validation = networkConfigUtils.validateNetworkConfig(interfaceForm.value)
    if (!validation.valid) {
      error.value = '配置验证失败: ' + validation.errors.join(', ')
      return
    }

    loading.value = true
    error.value = ''

    await interfaceConfigApi.setConfig(selectedInterface.value, interfaceForm.value)

    // 刷新接口列表
    await loadInterfaces()

    showInterfaceConfig.value = false
  } catch (err: any) {
    console.error('Failed to save interface config:', err)
    error.value = '保存接口配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 重启网络接口
const restartInterface = async (interfaceName: string) => {
  try {
    console.log('[NetworkManager] restartInterface called for:', interfaceName)
    buttonClicks.value++

    if (!confirm(`确定要重启网络接口 ${interfaceName} 吗？这可能会暂时中断网络连接。`)) {
      console.log('[NetworkManager] Restart cancelled by user')
      return
    }

    loading.value = true
    error.value = ''

    await interfaceConfigApi.restart(interfaceName)

    console.log('[NetworkManager] Interface restart command sent, waiting 3 seconds...')
    // 等待几秒钟后刷新
    setTimeout(() => {
      console.log('[NetworkManager] Refreshing interfaces after restart...')
      loadInterfaces()
    }, 3000)
  } catch (err: any) {
    console.error('[NetworkManager] Failed to restart interface:', err)
    error.value = '重启接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 代理配置相关功能
const openProxyConfig = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await proxyConfigApi.getConfig() as any
    proxyConfig.value = response

    showProxyConfig.value = true
  } catch (err: any) {
    console.error('Failed to load proxy config:', err)
    error.value = '加载代理配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 保存代理配置
const saveProxyConfig = async () => {
  try {
    loading.value = true
    error.value = ''

    await proxyConfigApi.setConfig(proxyConfig.value)

    showProxyConfig.value = false
  } catch (err: any) {
    console.error('Failed to save proxy config:', err)
    error.value = '保存代理配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 组件挂载时加载数据
onMounted(async () => {
  console.log('[NetworkManager] Component mounted, embedded mode:', props.embeddedMode)
  console.log('[NetworkManager] Loading interfaces...')

  await loadInterfaces()

  console.log('[NetworkManager] Interfaces loaded:', interfaces.value.length)

  // 如果有WiFi接口，自动扫描WiFi网络
  const hasWifi = interfaces.value.some(i => i.type === '无线' || i.type === 'wireless')
  if (hasWifi) {
    console.log('[NetworkManager] WiFi interface found, scanning networks...')
    await scanWifiNetworks()
  }

  // 设置自动刷新
  refreshInterval = setInterval(() => {
    console.log('[NetworkManager] Auto-refreshing interfaces...')
    loadInterfaces()
  }, 30000) // 每30秒刷新一次

  console.log('[NetworkManager] Setup complete')
})

// 组件卸载时清除定时器
onUnmounted(() => {
  console.log('[NetworkManager] Component unmounted')
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.network-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f9fafb;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  gap: 16px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-alert {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #fee;
  border: 1px solid #fcc;
  border-radius: 8px;
  color: #c33;
  margin: 20px;
}

.error-alert button {
  margin-left: auto;
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  color: #c33;
}

.nm-header {
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.nm-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.subtitle {
  font-size: 14px;
  color: #6b7280;
}

.nm-overview {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.5);
}

.overview-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.2s;
}

.overview-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.card-icon.active {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.card-info {
  flex: 1;
}

.card-title {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 2px;
}

.card-value {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.card-subtitle {
  font-size: 12px;
  color: #9ca3af;
}

.card-speed {
  display: flex;
  flex-direction: column;
  gap: 4px;
  text-align: right;
}

.speed-upload, .speed-download {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #6b7280;
}

.speed-upload {
  color: #10b981;
}

.speed-download {
  color: #3b82f6;
}

.card-action {
  padding: 6px 12px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.card-action:hover {
  background: #2563eb;
}

.nm-content {
  display: flex;
  gap: 20px;
  padding: 20px;
  flex: 1;
  overflow: hidden;
}

.nm-interfaces {
  flex: 2;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.nm-tools {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.section-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
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
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover:not(:disabled) {
  background: #f9fafb;
  border-color: #d1d5db;
}

.action-btn.primary {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.spinning svg {
  animation: spin 1s linear infinite;
}

.interface-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.interface-item {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  border: 1px solid rgba(255, 255, 255, 0.5);
  overflow: hidden;
  transition: all 0.2s;
}

.interface-item:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.interface-item.active {
  border-color: #10b981;
}

.interface-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
}

.interface-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.interface-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #e5e7eb;
  color: #6b7280;
}

.interface-icon.active {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
}

.interface-name {
  font-weight: 600;
  color: #1f2937;
}

.interface-type {
  font-size: 12px;
  color: #6b7280;
}

.interface-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #6b7280;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #d1d5db;
}

.status-indicator.active {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.interface-details {
  padding: 0 16px 16px;
  border-top: 1px solid #f3f4f6;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.detail-label {
  font-size: 11px;
  color: #6b7280;
}

.detail-value {
  font-size: 13px;
  color: #1f2937;
  font-family: monospace;
}

.interface-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  gap: 16px;
  color: #6b7280;
}

.empty-state svg {
  color: #d1d5db;
}

.tool-section {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 12px;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.tool-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 16px;
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 12px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 8px;
}

.stat-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.upload {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
}

.stat-icon.download {
  background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 11px;
  color: #6b7280;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.quick-actions {
  display: grid;
  grid-template-columns: 1fr;
  gap: 8px;
}

.quick-action {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-action:hover {
  background: rgba(255, 255, 255, 0.8);
  border-color: #d1d5db;
}

.quick-action svg {
  color: #6b7280;
}

.quick-action span {
  font-size: 14px;
  color: #1f2937;
}

.action-btn.full {
  width: 100%;
  justify-content: center;
}

/* 模态框样式 */
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
  padding: 20px;
}

.modal-content {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 16px;
  width: 100%;
  max-width: 600px;
  max-height: 80vh;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  font-size: 18px;
  font-weight: 600;
  color: #1f2937;
}

.modal-header button {
  background: none;
  border: none;
  cursor: pointer;
  color: #6b7280;
  padding: 4px;
}

.modal-body {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid #e5e7eb;
  background: #f9fafb;
}

.wifi-scan {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.scan-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.scan-header h4 {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
}

.network-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 400px;
  overflow-y: auto;
}

.network-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.network-item:hover {
  background: rgba(255, 255, 255, 0.8);
  border-color: #d1d5db;
}

.network-item.connected {
  border-color: #10b981;
  background: rgba(16, 185, 129, 0.05);
}

.network-item.secured {
  border-left: 3px solid #f59e0b;
}

.network-info {
  flex: 1;
}

.network-ssid {
  font-weight: 600;
  color: #1f2937;
}

.network-details {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #6b7280;
}

.network-item svg:not(:first-child):not(.connected-badge) {
  color: #d1d5db;
}

.connected-badge {
  padding: 4px 8px;
  background: #10b981;
  color: white;
  font-size: 11px;
  border-radius: 4px;
}

.empty-networks {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}

.dns-config {
  display: flex;
  flex-direction: column;
  gap: 20px;
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

.form-input {
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.radio-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.radio-item {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.radio-item input[type="radio"] {
  accent-color: #3b82f6;
}

.manual-dns {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.quick-dns {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.quick-dns label {
  font-size: 14px;
  color: #6b7280;
}

.dns-btn {
  padding: 6px 12px;
  background: #f3f4f6;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.dns-btn:hover {
  background: #e5e7eb;
  border-color: #9ca3af;
}

/* 过渡动画 */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

.fade-enter-to, .fade-leave-from {
  opacity: 1;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .nm-overview {
    grid-template-columns: 1fr;
  }

  .nm-content {
    flex-direction: column;
  }

  .nm-interfaces, .nm-tools {
    flex: none;
  }
}

/* 嵌入式模式样式 */
.network-manager.embedded-mode {
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.network-manager.embedded-mode .nm-header,
.network-manager.embedded-mode .nm-overview {
  display: none;
}

.network-manager.embedded-mode .nm-content {
  padding: 0;
  background: transparent;
}

.network-manager.embedded-mode .nm-interfaces {
  background: white;
  border-radius: 8px;
  padding: 16px;
}

.network-manager.embedded-mode .nm-tools {
  background: white;
  border-radius: 8px;
  padding: 16px;
}

.network-manager.embedded-mode .interface-item {
  margin-bottom: 12px;
}

.network-manager.embedded-mode .tool-section {
  margin-bottom: 16px;
}

/* 确保嵌入式模式下模态框正确显示 */
.network-manager.embedded-mode .modal-overlay {
  position: fixed !important;
  z-index: 9999 !important;
  top: 0 !important;
  left: 0 !important;
  right: 0 !important;
  bottom: 0 !important;
  background: rgba(0, 0, 0, 0.7) !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  padding: 20px !important;
}

.network-manager.embedded-mode .modal-content {
  position: relative !important;
  z-index: 10000 !important;
  max-width: 700px !important;
  width: 100% !important;
  max-height: 90vh !important;
  overflow-y: auto !important;
}

.network-manager.embedded-mode .modal-content * {
  position: relative;
  z-index: 10001;
}

@media (max-width: 640px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }

  .nm-overview {
    padding: 12px;
    gap: 12px;
  }

  .modal-content {
    margin: 10px;
    max-width: calc(100% - 20px);
  }
}
</style>