// Power Monitoring Type Definitions

export interface PowerData {
  timestamp: string
  cpuPackage: number
  cpuCore: number
  cpuUncore: number
  igpu: number
  dgpu: number
  hdd: number
  ssd: number
  mbram: number
  cooling: number
  usb: number
  powerLoss: number
  total: number
}

export interface PowerHistory {
  data: PowerData[]
  count: number
  days: number
}

export interface PowerStatistics {
  averagePower: number
  maxPower: number
  minPower: number
  totalEnergy: number
  sampleCount: number
}

export interface PowerOverview {
  current: PowerData
  today: PowerStatistics
  alerts: string[]
  timestamp: string
}

export interface PowerAlertRule {
  id: number
  name: string
  type: 'total' | 'cpu' | 'gpu' | 'hdd' | 'ssd'
  condition: '>' | '<' | '=' | '!=' | '>=' | '<='
  threshold: number
  duration: number
  severity: 'info' | 'warning' | 'critical'
  enabled: boolean
  actions: string
  lastTriggered?: string
  createdAt: string
  updatedAt: string
}

export interface PowerTrend {
  period: 'hourly' | 'daily' | 'weekly'
  trend: 'increasing' | 'decreasing' | 'stable'
  changePercent: number
  peakHours: number[]
  averageByHour: number[]
  predictions: number[]
  confidence: number
}

export interface PowerBreakdown {
  component: string
  power: number
  percentage: number
  color: string
}
