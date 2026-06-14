import api from './client'

// 防火墙规则类型
export interface FirewallRule {
  id?: number
  name: string
  action: 'allow' | 'deny'
  protocol: 'tcp' | 'udp' | 'both'
  port: string
  sourceIp: string
  description?: string
  enabled: boolean
  order: number
  created_at?: string
  updated_at?: string
}

// 防火墙配置类型
export interface FirewallConfig {
  enabled: boolean
  defaultPolicy: 'accept' | 'drop'
  allowedPorts: string[]
  logging: boolean
  icmp: boolean
}

// 防火墙API客户端
export const firewallApi = {
  // 获取防火墙规则列表
  getRules: () =>
    api.get('/api/security/firewall/rules'),

  // 创建防火墙规则
  createRule: (rule: Omit<FirewallRule, 'id' | 'created_at' | 'updated_at'>) =>
    api.post('/api/security/firewall/rules', rule),

  // 更新防火墙规则
  updateRule: (id: number, rule: Partial<FirewallRule>) =>
    api.put(`/api/security/firewall/rules/${id}`, rule),

  // 删除防火墙规则
  deleteRule: (id: number) =>
    api.delete(`/api/security/firewall/rules/${id}`),

  // 应用防火墙规则
  applyRules: () =>
    api.post('/api/security/firewall/apply'),

  // 获取防火墙配置
  getConfig: (): Promise<{ data: FirewallConfig }> =>
    api.get('/api/security/firewall/config'),

  // 设置防火墙配置
  setConfig: (config: Partial<FirewallConfig>) =>
    api.put('/api/security/firewall/config', config),

  // 启用/禁用防火墙
  setEnabled: (enabled: boolean) =>
    api.put('/api/security/firewall/config', { enabled })
}

// 防火墙验证工具
export const firewallUtils = {
  // 验证端口格式
  validatePort: (port: string): boolean => {
    if (port === 'any') return true

    // 单端口验证
    const singlePort = /^\d+$/.test(port)
    if (singlePort) {
      const portNum = parseInt(port)
      return portNum >= 1 && portNum <= 65535
    }

    // 端口范围验证
    const portRange = /^\d+-\d+$/.test(port)
    if (portRange) {
      const [start, end] = port.split('-').map(Number)
      return start >= 1 && end <= 65535 && start <= end
    }

    return false
  },

  // 验证IP地址格式
  validateIP: (ip: string): boolean => {
    if (ip === 'any') return true

    // IPv4验证
    const ipv4Pattern = /^(\d{1,3}\.){3}\d{1,3}(\/\d{1,2})?$/
    if (!ipv4Pattern.test(ip)) return false

    // 验证IP地址段
    if (ip.includes('/')) {
      const [ipPart, cidr] = ip.split('/')
      const cidrNum = parseInt(cidr)
      if (cidrNum < 0 || cidrNum > 32) return false

      const parts = ipPart.split('.')
      for (const part of parts) {
        const num = parseInt(part)
        if (num < 0 || num > 255) return false
      }
    } else {
      const parts = ip.split('.')
      for (const part of parts) {
        const num = parseInt(part)
        if (num < 0 || num > 255) return false
      }
    }

    return true
  },

  // 生成防火墙规则配置
  generateRuleConfig: (config: FirewallConfig): FirewallRule[] => {
    const rules: FirewallRule[] = []

    if (!config.enabled) {
      return rules
    }

    // 基于允许的端口生成规则
    config.allowedPorts.forEach((port, index) => {
      if (port.includes('-')) {
        // 端口范围
        rules.push({
          name: `Allow ${port}`,
          action: 'allow',
          protocol: 'tcp',
          port: port,
          sourceIp: 'any',
          enabled: true,
          order: index
        })
      } else {
        rules.push({
          name: `Allow ${port}`,
          action: 'allow',
          protocol: 'tcp',
          port: port,
          sourceIp: 'any',
          enabled: true,
          order: index
        })
      }
    })

    // ICMP规则
    if (config.icmp) {
      rules.push({
        name: 'Allow ICMP',
        action: 'allow',
        protocol: 'both',
        port: 'any',
        sourceIp: 'any',
        enabled: true,
        order: rules.length
      })
    }

    return rules
  }
}

export default firewallApi