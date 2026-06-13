<template>
  <div class="network-manager">
    <!-- 头部 -->
    <div class="nm-header">
      <h1>网络管理</h1>
      <p class="subtitle">管理网络连接和防火墙设置</p>
    </div>

    <!-- 网络状态概览 -->
    <div class="nm-overview">
      <div class="overview-card">
        <div class="card-icon active">
          <SignalIcon class="w-6 h-6" />
        </div>
        <div class="card-info">
          <div class="card-title">网络状态</div>
          <div class="card-value">已连接</div>
          <div class="card-subtitle">LAN 1 - 192.168.1.100</div>
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

      <div class="overview-card">
        <div class="card-icon">
          <WifiIcon class="w-6 h-6" />
        </div>
        <div class="card-info">
          <div class="card-title">无线网络</div>
          <div class="card-value">未连接</div>
          <div class="card-subtitle">点击设置无线连接</div>
        </div>
        <button class="card-action" @click="showWifiSetup = true">
          设置
        </button>
      </div>

      <div class="overview-card">
        <div class="card-icon">
          <ShieldIcon class="w-6 h-6" />
        </div>
        <div class="card-info">
          <div class="card-title">防火墙</div>
          <div class="card-value">已启用</div>
          <div class="card-subtitle">12 条活动规则</div>
        </div>
        <button class="card-action" @click="showFirewallSettings = true">
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
          <button class="action-btn" @click="refreshInterfaces">
            <ArrowPathIcon class="w-4 h-4" />
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

            <div v-if="iface.active" class="interface-details">
              <div class="detail-grid">
                <div class="detail-item">
                  <span class="detail-label">IP地址:</span>
                  <span class="detail-value">{{ iface.ip }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">子网掩码:</span>
                  <span class="detail-value">{{ iface.netmask }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">MAC地址:</span>
                  <span class="detail-value">{{ iface.mac }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">网关:</span>
                  <span class="detail-value">{{ iface.gateway }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">DNS服务器:</span>
                  <span class="detail-value">{{ iface.dns }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">速率:</span>
                  <span class="detail-value">{{ iface.speed }}</span>
                </div>
              </div>

              <div class="interface-actions">
                <button class="action-btn" @click="editInterface(iface)">
                  <PencilIcon class="w-4 h-4" />
                  编辑
                </button>
                <button class="action-btn" @click="disconnectInterface(iface)">
                  <XCircleIcon class="w-4 h-4" />
                  断开
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="interface-actions-global">
          <button class="action-btn primary" @click="showCreateBond = true">
            <PlusIcon class="w-4 h-4" />
            创建网络聚合
          </button>
          <button class="action-btn" @click="showVlanSettings = true">
            <TagIcon class="w-4 h-4" />
            VLAN设置
          </button>
        </div>
      </div>

      <!-- 右侧：其他网络设置 -->
      <div class="nm-settings">
        <!-- 代理设置 -->
        <div class="settings-section">
          <div class="settings-header">
            <h3>代理服务器</h3>
            <label class="toggle-switch">
              <input type="checkbox" v-model="proxyEnabled" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div v-if="proxyEnabled" class="proxy-config">
            <div class="form-group">
              <label>代理类型</label>
              <select v-model="proxyType" class="form-select">
                <option value="http">HTTP</option>
                <option value="https">HTTPS</option>
                <option value="socks">SOCKS5</option>
              </select>
            </div>
            <div class="form-group">
              <label>服务器地址</label>
              <input v-model="proxyServer" type="text" class="form-input" placeholder="proxy.example.com" />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label>端口</label>
                <input v-model="proxyPort" type="number" class="form-input" placeholder="8080" />
              </div>
              <div class="form-group">
                <label>用户名</label>
                <input v-model="proxyUser" type="text" class="form-input" placeholder="可选" />
              </div>
            </div>
            <button class="action-btn primary" @click="applyProxy">
              应用设置
            </button>
          </div>
        </div>

        <!-- DNS设置 -->
        <div class="settings-section">
          <div class="settings-header">
            <h3>DNS服务器</h3>
            <button class="text-btn" @click="showDnsSettings = true">
              高级设置
            </button>
          </div>

          <div class="dns-servers">
            <div class="dns-item">
              <ServerIcon class="w-4 h-4" />
              <span>192.168.1.1</span>
              <div class="status-badge success">主DNS</div>
            </div>
            <div class="dns-item">
              <ServerIcon class="w-4 h-4" />
              <span>8.8.8.8</span>
              <div class="status-badge">备用DNS</div>
            </div>
          </div>
        </div>

        <!-- DHCP设置 -->
        <div class="settings-section">
          <div class="settings-header">
            <h3>DHCP服务</h3>
            <label class="toggle-switch">
              <input type="checkbox" v-model="dhcpEnabled" />
              <span class="toggle-slider"></span>
            </label>
          </div>

          <div v-if="dhcpEnabled" class="dhcp-info">
            <div class="dhcp-range">
              <span class="range-label">地址池:</span>
              <span class="range-value">192.168.1.100 - 192.168.1.200</span>
            </div>
            <div class="dhcp-clients">
              <div class="clients-header">
                <span>已分配IP</span>
                <span>5/100</span>
              </div>
              <div class="client-list">
                <div v-for="client in dhcpClients" :key="client.ip" class="client-item">
                  <div class="client-info">
                    <span class="client-name">{{ client.name }}</span>
                    <span class="client-ip">{{ client.ip }}</span>
                  </div>
                  <span class="client-expiry">{{ client.expiry }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 带宽控制 -->
        <div class="settings-section">
          <div class="settings-header">
            <h3>流量控制</h3>
            <button class="text-btn" @click="showQosSettings = true">
              配置
            </button>
          </div>

          <div class="traffic-stats">
            <div class="stat-row">
              <span class="stat-label">今日上传:</span>
              <span class="stat-value">{{ formatBytes(todayUpload) }}</span>
            </div>
            <div class="stat-row">
              <span class="stat-label">今日下载:</span>
              <span class="stat-value">{{ formatBytes(todayDownload) }}</span>
            </div>
            <div class="stat-row">
              <span class="stat-label">本月上传:</span>
              <span class="stat-value">{{ formatBytes(monthUpload) }}</span>
            </div>
            <div class="stat-row">
              <span class="stat-label">本月下载:</span>
              <span class="stat-value">{{ formatBytes(monthDownload) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 模态框：无线网络设置 -->
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
                <button class="action-btn" @click="scanWifiNetworks">
                  <MagnifyingGlassIcon class="w-4 h-4" />
                  扫描
                </button>
              </div>
              <div class="network-list">
                <div
                  v-for="network in wifiNetworks"
                  :key="network.ssid"
                  class="network-item"
                  :class="{ secured: network.secured }"
                  @click="selectWifiNetwork(network)"
                >
                  <SignalIcon class="w-5 h-5" />
                  <div class="network-info">
                    <div class="network-ssid">{{ network.ssid }}</div>
                    <div class="network-details">
                      <span>{{ network.security }}</span>
                      <span>{{ network.signalStrength }}</span>
                    </div>
                  </div>
                  <LockIcon v-if="network.secured" class="w-4 h-4" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  SignalIcon,
  WifiIcon,
  ShieldIcon,
  ArrowUpIcon,
  ArrowDownIcon,
  ArrowPathIcon,
  ServerIcon,
  XCircleIcon,
  PencilIcon,
  PlusIcon,
  TagIcon,
  XMarkIcon,
  MagnifyingGlassIcon,
  LockIcon
} from '@heroicons/vue/24/outline'

interface NetworkInterface {
  name: string
  type: string
  active: boolean
  ip?: string
  netmask?: string
  mac?: string
  gateway?: string
  dns?: string
  speed?: string
}

// 状态管理
const uploadSpeed = ref(1250000) // bytes/s
const downloadSpeed = ref(8500000) // bytes/s
const proxyEnabled = ref(false)
const proxyType = ref('http')
const proxyServer = ref('')
const proxyPort = ref('')
const proxyUser = ref('')
const dhcpEnabled = ref(true)
const showWifiSetup = ref(false)
const showFirewallSettings = ref(false)
const showCreateBond = ref(false)
const showVlanSettings = ref(false)
const showDnsSettings = ref(false)
const showQosSettings = ref(false)

// 网络接口数据
const interfaces = ref<NetworkInterface[]>([
  {
    name: 'LAN 1',
    type: '以太网',
    active: true,
    ip: '192.168.1.100',
    netmask: '255.255.255.0',
    mac: '00:11:32:4A:5B:6C',
    gateway: '192.168.1.1',
    dns: '192.168.1.1',
    speed: '1 Gbps'
  },
  {
    name: 'LAN 2',
    type: '以太网',
    active: false,
    mac: '00:11:32:4A:5B:6D'
  },
  {
    name: 'WLAN 1',
    type: '无线',
    active: false
  }
])

// DHCP客户端数据
const dhcpClients = ref([
  { name: 'iPhone', ip: '192.168.1.101', expiry: '23小时59分' },
  { name: 'MacBook', ip: '192.168.1.102', expiry: '11小时30分' },
  { name: 'Android', ip: '192.168.1.103', expiry: '6小时45分' },
  { name: 'iPad', ip: '192.168.1.104', expiry: '2天15小时' },
  { name: 'SmartTV', ip: '192.168.1.105', expiry: '6天23小时' }
])

// WiFi网络数据
const wifiNetworks = ref([
  { ssid: 'Home-Network', secured: true, security: 'WPA2-Personal', signalStrength: '强' },
  { ssid: 'Guest-Network', secured: true, security: 'WPA2-Personal', signalStrength: '中' },
  { ssid: 'Neighbor-WiFi', secured: true, security: 'WPA2-Personal', signalStrength: '弱' },
  { ssid: 'Public-Free', secured: false, security: '开放', signalStrength: '中' }
])

// 流量统计
const todayUpload = ref(125000000) // bytes
const todayDownload = ref(850000000) // bytes
const monthUpload = ref(3500000000) // bytes
const monthDownload = ref(25000000000) // bytes

// 方法
const formatSpeed = (bytes: number): string => {
  const mbps = bytes / 1000000
  return mbps.toFixed(1) + ' MB/s'
}

const formatBytes = (bytes: number): string => {
  const gb = bytes / 1000000000
  return gb.toFixed(2) + ' GB'
}

const refreshInterfaces = () => {
  console.log('刷新网络接口')
}

const editInterface = (iface: NetworkInterface) => {
  console.log('编辑接口:', iface.name)
}

const disconnectInterface = (iface: NetworkInterface) => {
  if (confirm(`确定要断开 ${iface.name} 吗？`)) {
    iface.active = false
  }
}

const applyProxy = () => {
  console.log('应用代理设置')
}

const scanWifiNetworks = () => {
  console.log('扫描WiFi网络')
}

const selectWifiNetwork = (network: any) => {
  console.log('选择WiFi网络:', network.ssid)
}
</script>

<style scoped>
.network-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f9fafb;
}

.nm-header {
  padding: 20px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
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
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.overview-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.card-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
}

.card-icon.active {
  background: #ecfdf5;
  color: #10b981;
}

.card-info {
  flex: 1;
}

.card-title {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
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

.speed-upload,
.speed-download {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: #6b7280;
}

.card-action {
  padding: 6px 12px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.card-action:hover {
  background: #f3f4f6;
}

.nm-content {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
  flex: 1;
  overflow: hidden;
  padding: 20px;
}

.nm-interfaces {
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
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
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f3f4f6;
}

.action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.interface-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.interface-item {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.2s;
}

.interface-item.active {
  border-left: 4px solid #10b981;
}

.interface-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  cursor: pointer;
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
  background: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
}

.interface-icon.active {
  background: #ecfdf5;
  color: #10b981;
}

.interface-name {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 2px;
}

.interface-type {
  font-size: 12px;
  color: #6b7280;
}

.interface-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6b7280;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #d1d5db;
}

.status-indicator.active {
  background: #10b981;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.3);
}

.interface-details {
  padding: 0 16px 16px;
  border-top: 1px solid #f3f4f6;
  margin-top: 8px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.detail-label {
  color: #6b7280;
}

.detail-value {
  color: #1f2937;
  font-weight: 500;
}

.interface-actions {
  display: flex;
  gap: 8px;
}

.interface-actions-global {
  display: flex;
  gap: 8px;
}

.nm-settings {
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
}

.settings-section {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  padding: 16px;
}

.settings-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.settings-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.text-btn {
  padding: 4px 8px;
  background: transparent;
  border: none;
  color: #3b82f6;
  cursor: pointer;
  font-size: 12px;
  border-radius: 4px;
}

.text-btn:hover {
  background: #eff6ff;
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #d1d5db;
  transition: .4s;
  border-radius: 24px;
}

.toggle-slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
  border-radius: 50%;
}

input:checked + .toggle-slider {
  background-color: #3b82f6;
}

input:checked + .toggle-slider:before {
  transform: translateX(20px);
}

.proxy-config {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-group label {
  font-size: 12px;
  font-weight: 500;
  color: #6b7280;
}

.form-input,
.form-select {
  padding: 8px 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  font-size: 14px;
  background: white;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.dns-servers {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.dns-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f9fafb;
  border-radius: 8px;
  font-size: 13px;
}

.status-badge {
  padding: 2px 8px;
  background: #f3f4f6;
  border-radius: 4px;
  font-size: 11px;
  color: #6b7280;
}

.status-badge.success {
  background: #ecfdf5;
  color: #10b981;
}

.dhcp-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.dhcp-range {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.range-label {
  color: #6b7280;
}

.range-value {
  color: #1f2937;
  font-weight: 500;
}

.dhcp-clients {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.clients-header {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #6b7280;
}

.client-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.client-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  background: #f9fafb;
  border-radius: 6px;
  font-size: 12px;
}

.client-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.client-name {
  font-weight: 500;
  color: #1f2937;
}

.client-ip {
  color: #6b7280;
}

.client-expiry {
  color: #9ca3af;
}

.traffic-stats {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}

.stat-label {
  color: #6b7280;
}

.stat-value {
  color: #1f2937;
  font-weight: 500;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  width: 90%;
  max-width: 600px;
  background: white;
  border-radius: 16px;
  overflow: hidden;
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
  padding: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 6px;
}

.modal-header button:hover {
  background: #f3f4f6;
}

.modal-body {
  padding: 20px;
  max-height: 500px;
  overflow-y: auto;
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
  font-size: 14px;
  font-weight: 600;
  color: #1f2937;
}

.network-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.network-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.2s;
}

.network-item:hover {
  background: #eff6ff;
  border-color: #3b82f6;
}

.network-item.secured {
  border-left: 3px solid #10b981;
}

.network-info {
  flex: 1;
}

.network-ssid {
  font-size: 14px;
  font-weight: 500;
  color: #1f2937;
  margin-bottom: 2px;
}

.network-details {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: #6b7280;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>