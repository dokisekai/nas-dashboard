<template>
  <div class="network-settings-panel">
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

    <!-- 网络接口列表 -->
    <div v-if="interfaces.length > 0" class="interfaces-list">
      <div
        v-for="iface in interfaces"
        :key="iface.name"
        class="network-interface-card"
      >
        <!-- 接口基本信息 -->
        <div class="interface-header">
          <div class="interface-info">
            <div class="interface-name">{{ iface.name }}</div>
            <div class="interface-type">
              <SignalIcon v-if="iface.active" class="w-4 h-4 online" />
              <SignalIcon v-else class="w-4 h-4 offline" />
              {{ iface.type === 'ethernet' ? '以太网' : iface.type === 'wireless' ? '无线' : iface.type }}
              <span v-if="iface.active" class="status-badge active">已连接</span>
              <span v-else class="status-badge inactive">未连接</span>
            </div>
          </div>

          <!-- 快速操作 -->
          <div class="interface-actions">
            <button
              v-if="!iface.active"
              class="action-btn primary"
              @click="connectInterface(iface.name)"
              :disabled="loading"
            >
              <SignalIcon class="w-4 h-4" />
              启用
            </button>
            <button
              v-else
              class="action-btn danger"
              @click="disconnectInterface(iface.name)"
              :disabled="loading"
            >
              <XCircleIcon class="w-4 h-4" />
              禁用
            </button>
          </div>
        </div>

        <!-- 接口详细信息 -->
        <div class="interface-details">
          <div class="detail-row">
            <span class="detail-label">IP地址:</span>
            <span class="detail-value">{{ iface.ip || '未分配' }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">MAC地址:</span>
            <span class="detail-value">{{ iface.mac || '未知' }}</span>
          </div>
          <div v-if="iface.gateway" class="detail-row">
            <span class="detail-label">网关:</span>
            <span class="detail-value">{{ iface.gateway }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">MTU:</span>
            <span class="detail-value">{{ iface.mtu || 1500 }}</span>
          </div>
        </div>

        <!-- 内联配置区域 -->
        <div class="interface-config">
          <div class="config-section">
            <h4>IP配置方式</h4>
            <div class="config-options">
              <label class="radio-option">
                <input
                  type="radio"
                  :name="`ip-method-${iface.name}`"
                  :checked="getInterfaceConfig(iface.name)?.ipv4Method === 'dhcp'"
                  @change="setIPMethod(iface.name, 'dhcp')"
                />
                <span>动态IP (DHCP)</span>
              </label>
              <label class="radio-option">
                <input
                  type="radio"
                  :name="`ip-method-${iface.name}`"
                  :checked="getInterfaceConfig(iface.name)?.ipv4Method === 'static'"
                  @change="setIPMethod(iface.name, 'static')"
                />
                <span>静态IP</span>
              </label>
            </div>
          </div>

          <!-- 静态IP配置表单 -->
          <div v-if="getInterfaceConfig(iface.name)?.ipv4Method === 'static'" class="static-ip-form">
            <div class="form-group">
              <label>IP地址:</label>
              <input
                v-model="interfaceConfigs[iface.name].ipAddress"
                type="text"
                placeholder="192.168.1.100"
                class="form-input"
                @blur="saveInterfaceConfig(iface.name)"
              />
            </div>
            <div class="form-group">
              <label>子网掩码:</label>
              <input
                v-model="interfaceConfigs[iface.name].netmask"
                type="text"
                placeholder="255.255.255.0"
                class="form-input"
                @blur="saveInterfaceConfig(iface.name)"
              />
            </div>
            <div class="form-group">
              <label>默认网关:</label>
              <input
                v-model="interfaceConfigs[iface.name].gateway"
                type="text"
                placeholder="192.168.1.1"
                class="form-input"
                @blur="saveInterfaceConfig(iface.name)"
              />
            </div>
          </div>

          <!-- MTU配置 -->
          <div class="config-section">
            <h4>MTU设置</h4>
            <div class="form-inline">
              <input
                v-model.number="interfaceConfigs[iface.name].mtu"
                type="number"
                min="576"
                max="9000"
                class="form-input-small"
                @change="saveInterfaceConfig(iface.name)"
              />
              <button
                class="action-btn"
                @click="saveInterfaceConfig(iface.name)"
                :disabled="loading"
              >
                应用
              </button>
            </div>
          </div>

          <!-- 接口操作 -->
          <div class="config-section">
            <h4>接口操作</h4>
            <div class="action-buttons">
              <button
                class="action-btn"
                @click="restartInterface(iface.name)"
                :disabled="loading"
              >
                <ArrowPathIcon class="w-4 h-4" :class="{ spinning: loading }" />
                重启接口
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 代理配置区域 -->
    <div class="proxy-config-section">
      <h3>代理服务器设置</h3>
      <div class="proxy-card">
        <div class="proxy-header">
          <label class="switch-label">
            <input
              type="checkbox"
              v-model="proxyConfig.enabled"
              @change="saveProxyConfig"
            />
            <span>启用代理服务器</span>
          </label>
        </div>

        <div v-if="proxyConfig.enabled" class="proxy-settings">
          <div class="form-row">
            <div class="form-group">
              <label>代理类型:</label>
              <select
                v-model="proxyConfig.type"
                class="form-input"
                @change="saveProxyConfig"
              >
                <option value="http">HTTP</option>
                <option value="https">HTTPS</option>
                <option value="socks4">SOCKS4</option>
                <option value="socks5">SOCKS5</option>
              </select>
            </div>

            <div class="form-group">
              <label>服务器地址:</label>
              <input
                v-model="proxyConfig.server"
                type="text"
                placeholder="proxy.example.com"
                class="form-input"
                @blur="saveProxyConfig"
              />
            </div>

            <div class="form-group">
              <label>端口:</label>
              <input
                v-model.number="proxyConfig.port"
                type="number"
                min="1"
                max="65535"
                class="form-input-small"
                @blur="saveProxyConfig"
              />
            </div>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>用户名 (可选):</label>
              <input
                v-model="proxyConfig.username"
                type="text"
                placeholder="username"
                class="form-input"
                @blur="saveProxyConfig"
              />
            </div>

            <div class="form-group">
              <label>密码 (可选):</label>
              <input
                v-model="proxyConfig.password"
                type="password"
                placeholder="password"
                class="form-input"
                @blur="saveProxyConfig"
              />
            </div>
          </div>

          <div class="form-group">
            <label>绕过代理的地址:</label>
            <input
              v-model="proxyConfig.bypassList"
              type="text"
              placeholder="localhost,127.0.0.1,*.local"
              class="form-input"
              @blur="saveProxyConfig"
            />
            <small>使用逗号分隔多个地址</small>
          </div>
        </div>
      </div>
    </div>

    <!-- DNS配置区域 -->
    <div class="dns-config-section">
      <h3>DNS设置</h3>
      <div class="dns-card">
        <div class="dns-method-selector">
          <label class="radio-option">
            <input
              type="radio"
              :value="'auto'"
              :checked="dnsConfig.method === 'auto'"
              @change="setDNSMethod('auto')"
            />
            <span>自动获取DNS</span>
          </label>
          <label class="radio-option">
            <input
              type="radio"
              :value="'manual'"
              :checked="dnsConfig.method === 'manual'"
              @change="setDNSMethod('manual')"
            />
            <span>手动配置DNS</span>
          </label>
        </div>

        <div v-if="dnsConfig.method === 'manual'" class="dns-settings">
          <div class="form-row">
            <div class="form-group">
              <label>主DNS服务器:</label>
              <input
                v-model="dnsConfig.primary"
                type="text"
                placeholder="192.168.1.1"
                class="form-input"
                @blur="saveDNSConfig"
              />
            </div>
            <div class="form-group">
              <label>备用DNS服务器:</label>
              <input
                v-model="dnsConfig.secondary"
                type="text"
                placeholder="8.8.8.8"
                class="form-input"
                @blur="saveDNSConfig"
              />
            </div>
          </div>
          <div class="quick-dns">
            <label>快速设置:</label>
            <button @click="setQuickDNS('8.8.8.8', '8.8.4.4')" class="dns-btn">Google DNS</button>
            <button @click="setQuickDNS('1.1.1.1', '1.0.0.1')" class="dns-btn">Cloudflare DNS</button>
            <button @click="setQuickDNS('208.67.222.222', '208.67.220.220')" class="dns-btn">OpenDNS</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!interfaces.length && !loading" class="empty-state">
      <SignalIcon class="w-12 h-12" />
      <p>未发现网络接口</p>
      <button class="action-btn primary" @click="loadInterfaces">
        重新扫描
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import {
  SignalIcon,
  XCircleIcon,
  ArrowPathIcon
} from '@heroicons/vue/24/outline'
import { networkApi, dnsApi, interfaceConfigApi, proxyConfigApi } from '../../api'
import type { NetworkInterface, DNSConfig } from '../../api/network'
import type { ProxyConfig } from '../../api/interface_config'

// 状态管理
const interfaces = ref<any[]>([])
const loading = ref(false)
const error = ref('')

// 接口配置存储
const interfaceConfigs = ref<Record<string, any>>({})

// DNS配置
const dnsConfig = ref<DNSConfig>({
  method: 'auto',
  primary: '192.168.1.1',
  secondary: '8.8.8.8'
})

// 代理配置
const proxyConfig = ref<ProxyConfig>({
  enabled: false,
  type: 'http',
  server: '',
  port: 8080,
  bypassList: ['localhost', '127.0.0.1']
})

// 刷新间隔
let refreshInterval: number

// 获取接口配置
const getInterfaceConfig = (interfaceName: string) => {
  return interfaceConfigs.value[interfaceName] || {
    ipv4Method: 'dhcp',
    mtu: 1500
  }
}

// 加载网络接口
const loadInterfaces = async () => {
  try {
    loading.value = true
    error.value = ''

    const response = await networkApi.getInterfaces() as any

    if (response && response.interfaces) {
      interfaces.value = response.interfaces.map((iface: any) => ({
        name: iface.name,
        type: iface.type === 'ethernet' ? 'ethernet' : iface.type === 'wireless' ? 'wireless' : iface.type,
        active: iface.up || false,
        ip: iface.addresses && iface.addresses.length > 0 ? iface.addresses[0] : '',
        netmask: '',
        mac: iface.hardwareAddr || '',
        gateway: '',
        dns: '',
        speed: '',
        mtu: iface.mtu || 1500,
        tx_bytes: iface.tx_bytes || 0,
        rx_bytes: iface.rx_bytes || 0
      }))

      // 为每个接口初始化配置
      for (const iface of interfaces.value) {
        if (!interfaceConfigs.value[iface.name]) {
          interfaceConfigs.value[iface.name] = {
            ipv4Method: 'dhcp',
            ipAddress: iface.ip || '',
            netmask: '255.255.255.0',
            gateway: '',
            mtu: iface.mtu || 1500
          }
        }
      }
    } else if (Array.isArray(response)) {
      interfaces.value = response.map((iface: any) => ({
        name: iface.name,
        type: iface.type === 'ethernet' ? 'ethernet' : iface.type === 'wireless' ? 'wireless' : iface.type,
        active: iface.up || false,
        ip: iface.addresses && iface.addresses.length > 0 ? iface.addresses[0] : '',
        netmask: '',
        mac: iface.hardwareAddr || '',
        gateway: '',
        dns: '',
        speed: '',
        mtu: iface.mtu || 1500
      }))
    }
  } catch (err: any) {
    console.error('Failed to load network interfaces:', err)
    error.value = '加载网络接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 连接接口
const connectInterface = async (interfaceName: string) => {
  try {
    loading.value = true
    error.value = ''

    await networkApi.controlInterface(interfaceName, 'up')
    await loadInterfaces()
  } catch (err: any) {
    console.error('Failed to connect interface:', err)
    error.value = '启用接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 断开接口
const disconnectInterface = async (interfaceName: string) => {
  try {
    loading.value = true
    error.value = ''

    await networkApi.controlInterface(interfaceName, 'down')
    await loadInterfaces()
  } catch (err: any) {
    console.error('Failed to disconnect interface:', err)
    error.value = '禁用接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 设置IP方法
const setIPMethod = async (interfaceName: string, method: string) => {
  try {
    loading.value = true
    error.value = ''

    // 更新配置
    if (!interfaceConfigs.value[interfaceName]) {
      interfaceConfigs.value[interfaceName] = {}
    }
    interfaceConfigs.value[interfaceName].ipv4Method = method

    // 保存到后端
    if (method === 'dhcp') {
      await interfaceConfigApi.setConfig(interfaceName, {
        ipv4Method: 'dhcp'
      })
    }

    await loadInterfaces()
  } catch (err: any) {
    console.error('Failed to set IP method:', err)
    error.value = '设置IP方式失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 保存接口配置
const saveInterfaceConfig = async (interfaceName: string) => {
  try {
    loading.value = true
    error.value = ''

    const config = interfaceConfigs.value[interfaceName]
    if (!config) return

    // 如果是静态IP，需要保存完整配置
    if (config.ipv4Method === 'static') {
      await interfaceConfigApi.setConfig(interfaceName, {
        ipv4Method: 'static',
        ipAddress: config.ipAddress,
        netmask: config.netmask,
        gateway: config.gateway,
        mtu: config.mtu
      })
    } else if (config.mtu && config.mtu !== 1500) {
      // 只更新MTU
      await interfaceConfigApi.setConfig(interfaceName, {
        ipv4Method: 'dhcp',
        mtu: config.mtu
      })
    }

    await loadInterfaces()
  } catch (err: any) {
    console.error('Failed to save interface config:', err)
    error.value = '保存接口配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 重启接口
const restartInterface = async (interfaceName: string) => {
  try {
    loading.value = true
    error.value = ''

    await interfaceConfigApi.restart(interfaceName)

    // 等待后刷新接口状态
    setTimeout(() => {
      loadInterfaces()
    }, 3000)
  } catch (err: any) {
    console.error('Failed to restart interface:', err)
    error.value = '重启接口失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 设置DNS方法
const setDNSMethod = async (method: 'auto' | 'manual') => {
  dnsConfig.value.method = method
  if (method === 'auto') {
    await saveDNSConfig()
  }
}

// 保存DNS配置
const saveDNSConfig = async () => {
  try {
    loading.value = true
    error.value = ''

    await dnsApi.setConfig(dnsConfig.value)
  } catch (err: any) {
    console.error('Failed to save DNS config:', err)
    error.value = '保存DNS配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 快速设置DNS
const setQuickDNS = async (primary: string, secondary: string) => {
  dnsConfig.value.primary = primary
  dnsConfig.value.secondary = secondary
  await saveDNSConfig()
}

// 保存代理配置
const saveProxyConfig = async () => {
  try {
    loading.value = true
    error.value = ''

    await proxyConfigApi.setConfig(proxyConfig.value)
  } catch (err: any) {
    console.error('Failed to save proxy config:', err)
    error.value = '保存代理配置失败: ' + (err.response?.data?.error || err.message)
  } finally {
    loading.value = false
  }
}

// 组件挂载
onMounted(async () => {
  await loadInterfaces()

  // 加载代理配置
  try {
    const proxyResp = await proxyConfigApi.getConfig() as any
    proxyConfig.value = proxyResp
  } catch (err) {
    console.warn('Failed to load proxy config:', err)
  }

  // 设置自动刷新
  refreshInterval = setInterval(() => {
    loadInterfaces()
  }, 30000)
})

// 组件卸载
onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.network-settings-panel {
  padding: 0;
  max-width: 100%;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #e5e7eb;
  border-top-color: #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-alert {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  color: #991b1b;
  margin-bottom: 16px;
}

.interfaces-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.network-interface-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s;
}

.network-interface-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.interface-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background: #f9fafb;
  border-bottom: 1px solid #e5e7eb;
}

.interface-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.interface-name {
  font-weight: 600;
  color: #111827;
  font-size: 16px;
}

.interface-type {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #6b7280;
  font-size: 14px;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: #dcfce7;
  color: #166534;
}

.status-badge.inactive {
  background: #f3f4f6;
  color: #6b7280;
}

.online {
  color: #166534;
}

.offline {
  color: #9ca3af;
}

.interface-actions {
  display: flex;
  gap: 8px;
}

.interface-details {
  padding: 16px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 12px;
  background: white;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f3f4f6;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  color: #6b7280;
  font-size: 14px;
}

.detail-value {
  color: #111827;
  font-size: 14px;
  font-weight: 500;
}

.interface-config {
  padding: 16px;
  background: #f9fafb;
  border-top: 1px solid #e5e7eb;
}

.config-section {
  margin-bottom: 16px;
}

.config-section:last-child {
  margin-bottom: 0;
}

.config-section h4 {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.config-options {
  display: flex;
  gap: 16px;
}

.radio-option {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.radio-option:hover {
  background: #f3f4f6;
}

.radio-option input[type="radio"] {
  margin: 0;
}

.radio-option span {
  font-size: 14px;
  color: #374151;
}

.static-ip-form {
  padding: 12px;
  background: white;
  border-radius: 8px;
  margin-top: 8px;
}

.form-group {
  margin-bottom: 12px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 4px;
}

.form-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.form-input-small {
  width: 80px;
  padding: 4px 8px;
  font-size: 13px;
}

.form-inline {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.primary {
  background: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.action-btn.danger {
  background: #ef4444;
  color: white;
}

.action-btn.danger:hover {
  background: #dc2626;
}

.proxy-config-section,
.dns-config-section {
  margin-top: 24px;
}

.proxy-config-section h3,
.dns-config-section h3 {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  margin-bottom: 12px;
}

.proxy-card,
.dns-card {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  padding: 16px;
}

.proxy-header {
  margin-bottom: 16px;
}

.switch-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #111827;
  cursor: pointer;
}

.proxy-settings {
  margin-top: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}

.form-row:last-child {
  margin-bottom: 0;
}

.quick-dns {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-top: 12px;
}

.dns-btn {
  padding: 6px 12px;
  background: #f3f4f6;
  border: 1px solid #e5e7eb;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.dns-btn:hover {
  background: #e5e7eb;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  text-align: center;
  color: #6b7280;
}
</style>