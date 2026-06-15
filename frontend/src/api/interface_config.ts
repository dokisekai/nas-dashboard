import api from './client'

// 网络接口配置类型
export interface InterfaceConfig {
  name: string
  ipv4Method: string          // dhcp, static, pppoe
  ipv6Method: string          // auto, static, disabled
  ipAddress?: string
  netmask?: string
  gateway?: string
  macAddress?: string
  mtu?: number
  dhcp: boolean
  autoConnect: boolean
  enabled: boolean
}

// 接口配置请求类型
export interface InterfaceConfigRequest {
  ipv4Method: string         // dhcp, static, pppoe
  ipAddress?: string
  netmask?: string
  gateway?: string
  mtu?: number
}

// PPPoE配置类型
export interface PPPoEConfig {
  interface: string
  username: string
  password: string
  serviceName?: string
  acName?: string
  autoConnect: boolean
  enabled: boolean
}

// PPPoE配置请求类型
export interface PPPoEConfigRequest {
  username: string
  password: string
  serviceName?: string
  acName?: string
  autoConnect: boolean
}

// 代理配置类型
export interface ProxyConfig {
  enabled: boolean
  type: string                // http, https, socks4, socks5
  server: string
  port: number
  username?: string
  password?: string
  bypassList?: string[]
}

// 接口配置API
export const interfaceConfigApi = {
  // 获取接口配置
  getConfig: (interfaceName: string) =>
    api.get<InterfaceConfig>(`/api/network/interface/${interfaceName}/config`),

  // 设置接口配置
  setConfig: (interfaceName: string, config: InterfaceConfigRequest) =>
    api.put(`/api/network/interface/${interfaceName}/config`, config),

  // 重启接口
  restart: (interfaceName: string) =>
    api.post(`/api/network/interface/${interfaceName}/restart`),
}

// PPPoE配置API
export const pppoeConfigApi = {
  // 获取PPPoE配置
  getConfig: (interfaceName: string) =>
    api.get<PPPoEConfig>(`/api/network/interface/${interfaceName}/pppoe`),

  // 配置PPPoE
  configure: (interfaceName: string, config: PPPoEConfigRequest) =>
    api.post(`/api/network/interface/${interfaceName}/pppoe`, config),
}

// 代理配置API
export const proxyConfigApi = {
  // 获取代理配置
  getConfig: () =>
    api.get<ProxyConfig>('/api/network/proxy'),

  // 设置代理配置
  setConfig: (config: ProxyConfig) =>
    api.post('/api/network/proxy', config),
}

// 网络配置工具
export const networkConfigUtils = {
  // 验证IP地址格式
  isValidIP: (ip: string): boolean => {
    const ipRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
    return ipRegex.test(ip)
  },

  // 验证子网掩码格式
  isValidNetmask: (netmask: string): boolean => {
    const netmaskRegex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/
    if (!netmaskRegex.test(netmask)) return false

    // 检查是否为有效的子网掩码（连续的1后面跟连续的0）
    const parts = netmask.split('.').map(Number)
    let binary = ''
    for (let i = 0; i < 4; i++) {
      binary += parts[i].toString(2).padStart(8, '0')
    }

    // 检查是否为有效的子网掩码格式
    const validOnes = ['11111111111111111111111111111111', '11111111111111111111111100000000',
                       '11111111111111111111111111000000', '11111111111111111111111110000000',
                       '11111111111111111111111100000000', '11111111111111111111111000000000',
                       '11111111111111111111110000000000', '11111111111111111111100000000000',
                       '11111111111111111111000000000000', '11111111111111111110000000000000',
                       '11111111111111111100000000000000', '11111111111111111000000000000000',
                       '11111111111111110000000000000000', '11111111111111100000000000000000',
                       '11111111111111000000000000000000', '11111111111110000000000000000000',
                       '11111111111100000000000000000000', '11111111111000000000000000000000',
                       '11111111110000000000000000000000', '11111111100000000000000000000000',
                       '11111111000000000000000000000000', '11111110000000000000000000000000',
                       '11111100000000000000000000000000', '11111000000000000000000000000000',
                       '11110000000000000000000000000000', '11100000000000000000000000000000',
                       '11000000000000000000000000000000', '10000000000000000000000000000000',
                       '00000000000000000000000000000000']

    return validOnes.includes(binary)
  },

  // 验证CIDR格式的子网掩码
  isValidCIDR: (cidr: string): boolean => {
    const cidrRegex = /^\/([0-9]|[1-2][0-9]|3[0-2])$/
    return cidrRegex.test(cidr)
  },

  // 验证MAC地址格式
  isValidMAC: (mac: string): boolean => {
    const macRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$/
    return macRegex.test(mac)
  },

  // 验证MTU值
  isValidMTU: (mtu: number): boolean => {
    return mtu >= 576 && mtu <= 9000
  },

  // 验证端口号
  isValidPort: (port: number): boolean => {
    return port > 0 && port <= 65535
  },

  // 将CIDR转换为子网掩码
  cidrToNetmask: (cidr: number): string => {
    if (cidr < 0 || cidr > 32) return '0.0.0.0'

    let mask = 0xffffffff << (32 - cidr)
    const parts = [
      (mask >>> 24) & 0xff,
      (mask >>> 16) & 0xff,
      (mask >>> 8) & 0xff,
      mask & 0xff
    ]

    return parts.join('.')
  },

  // 将子网掩码转换为CIDR
  netmaskToCIDR: (netmask: string): number => {
    const parts = netmask.split('.').map(Number)
    let binary = ''
    for (let i = 0; i < 4; i++) {
      binary += parts[i].toString(2).padStart(8, '0')
    }

    let cidr = 0
    for (let i = 0; i < 32; i++) {
      if (binary[i] === '1') cidr++
      else break
    }

    return cidr
  },

  // 获取默认网关建议
  suggestGateway: (ip: string, netmask: string): string => {
    if (!ip || !netmask) return ''

    const ipParts = ip.split('.').map(Number)
    const maskParts = netmask.split('.').map(Number)
    const gatewayParts: number[] = []

    for (let i = 0; i < 4; i++) {
      gatewayParts.push((ipParts[i] & maskParts[i]))
    }

    // 通常网关是网络地址的第一个可用IP
    gatewayParts[3] = gatewayParts[3] + 1

    return gatewayParts.join('.')
  },

  // 验证网络配置的完整性
  validateNetworkConfig: (config: InterfaceConfigRequest): { valid: boolean; errors: string[] } => {
    const errors: string[] = []

    if (config.ipv4Method === 'static') {
      if (!config.ipAddress) {
        errors.push('静态IP配置需要IP地址')
      } else if (!networkConfigUtils.isValidIP(config.ipAddress)) {
        errors.push('无效的IP地址格式')
      }

      if (!config.netmask) {
        errors.push('静态IP配置需要子网掩码')
      } else if (!networkConfigUtils.isValidNetmask(config.netmask)) {
        errors.push('无效的子网掩码格式')
      }

      if (config.gateway && !networkConfigUtils.isValidIP(config.gateway)) {
        errors.push('无效的网关地址格式')
      }
    }

    if (config.mtu && !networkConfigUtils.isValidMTU(config.mtu)) {
      errors.push('MTU值必须在576-9000之间')
    }

    return {
      valid: errors.length === 0,
      errors
    }
  }
}