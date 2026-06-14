import api from './client'

// 网络接口类型
export interface NetworkInterface {
  name: string
  type: string
  active: boolean
  ip: string
  netmask: string
  mac: string
  gateway: string
  dns: string
  speed: string
  mtu?: number
  tx_bytes?: number
  rx_bytes?: number
  tx_packets?: number
  rx_packets?: number
  tx_errors?: number
  rx_errors?: number
  tx_dropped?: number
  rx_dropped?: number
}

// WiFi网络类型
export interface WiFiNetwork {
  ssid: string
  bssid: string
  security: string
  signalStrength: number
  channel: number
  connected: boolean
  connecting: boolean
  uploadSpeed?: number
  downloadSpeed?: number
  ipAddress?: string
}

// WiFi连接请求
export interface WiFiConnectionRequest {
  ssid: string
  password?: string
  security?: string
  bssid?: string
  isHidden?: boolean
}

// DNS配置
export interface DNSConfig {
  method: 'auto' | 'manual'
  primary: string
  secondary: string
}

// 网络统计
export interface NetworkStats {
  interface: string
  txBytes: number
  rxBytes: number
  txPackets: number
  rxPackets: number
  txErrors: number
  rxErrors: number
  txDropped: number
  rxDropped: number
  timestamp: number
}

// ==================== 网络接口 API ====================
export const networkApi = {
  // 获取所有网络接口
  getInterfaces: () => api.get<{interfaces: NetworkInterface[], timestamp: number}>('/api/network/interfaces'),

  // 获取以太网接口
  getEthernetInterfaces: () => api.get<NetworkInterface[]>('/api/network/interfaces/ethernet'),

  // 获取WiFi接口
  getWiFiInterfaces: () => api.get<NetworkInterface[]>('/api/network/interfaces/wifi'),

  // 控制接口（启动/停止/重启）
  controlInterface: (interfaceName: string, action: 'up' | 'down' | 'restart') =>
    api.post(`/api/network/interfaces/${interfaceName}/${action}`),

  // 获取网络统计
  getStats: (interfaceName: string) =>
    api.get<NetworkStats>(`/api/network/interfaces/${interfaceName}/stats`),
}

// ==================== WiFi API ====================
export const wifiApi = {
  // 扫描WiFi网络
  scanNetworks: () => api.get<WiFiNetwork[]>('/api/network/wifi/scan'),

  // 连接到WiFi网络
  connect: (request: WiFiConnectionRequest) =>
    api.post('/api/network/wifi/connect', request),

  // 断开WiFi连接
  disconnect: () => api.post('/api/network/wifi/disconnect'),

  // 获取当前WiFi连接信息
  getCurrentConnection: () => api.get<WiFiNetwork>('/api/network/wifi/current'),

  // 获取WiFi接口列表
  getInterfaces: () => api.get<NetworkInterface[]>('/api/network/interfaces/wifi'),
}

// ==================== DNS API ====================
export const dnsApi = {
  // 获取DNS配置
  getConfig: () => api.get<DNSConfig>('/api/network/dns'),

  // 设置DNS配置
  setConfig: (config: DNSConfig) =>
    api.post('/api/network/dns', config),
}

// ==================== 网络工具 ====================
export const networkUtils = {
  // 格式化速度
  formatSpeed: (bytesPerSecond: number): string => {
    const mbps = bytesPerSecond / 1000000
    if (mbps >= 1000) {
      return (mbps / 1000).toFixed(2) + ' GB/s'
    }
    if (mbps >= 1) {
      return mbps.toFixed(2) + ' MB/s'
    }
    return (bytesPerSecond / 1000).toFixed(2) + ' KB/s'
  },

  // 格式化字节数
  formatBytes: (bytes: number): string => {
    const gb = bytes / 1000000000
    if (gb >= 1) {
      return gb.toFixed(2) + ' GB'
    }
    const mb = bytes / 1000000
    if (mb >= 1) {
      return mb.toFixed(2) + ' MB'
    }
    const kb = bytes / 1000
    return kb.toFixed(2) + ' KB'
  },

  // 计算信号强度百分比
  getSignalPercentage: (dBm: number): number => {
    if (dBm >= -50) return 100
    if (dBm <= -100) return 0
    return Math.round(((dBm + 100) / 50) * 100)
  },

  // 获取信号强度描述
  getSignalStrength: (dBm: number): string => {
    const percentage = networkUtils.getSignalPercentage(dBm)
    if (percentage >= 70) return '强'
    if (percentage >= 40) return '中'
    return '弱'
  },
}
