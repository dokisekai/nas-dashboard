<template>
  <div class="wifi-scan-component">
    <div class="scan-header">
      <h4>Wi-Fi网络扫描</h4>
      <button
        class="scan-btn"
        @click="scanNetworks"
        :disabled="scanning"
      >
        <svg v-if="scanning" class="w-4 h-4 spinning" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        {{ scanning ? '扫描中...' : '扫描网络' }}
      </button>
    </div>

    <div v-if="scanning" class="scanning-state">
      <div class="spinner"></div>
      <p>正在扫描附近的Wi-Fi网络...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <ExclamationTriangleIcon class="w-5 h-5" />
      <span>{{ error }}</span>
      <button @click="scanNetworks">重试</button>
    </div>

    <div v-else-if="networks.length === 0" class="empty-state">
      <SignalIcon class="w-8 h-8" />
      <p>暂无可用网络，请点击扫描按钮</p>
    </div>

    <div v-else class="networks-list">
      <div
        v-for="network in networks"
        :key="network.ssid || network.bssid"
        class="network-item"
        :class="{
          connected: network.connected,
          secured: network.secured,
          selected: selectedNetwork === network.ssid
        }"
      >
        <div class="network-info" @click="selectNetwork(network)">
          <div class="network-icon">
            <SignalIcon class="w-5 h-5" />
            <div class="signal-bars" :class="getSignalStrength(network.signalLevel)">
              <div class="bar bar1"></div>
              <div class="bar bar2"></div>
              <div class="bar bar3"></div>
              <div class="bar bar4"></div>
            </div>
          </div>

          <div class="network-details">
            <div class="network-name">
              {{ network.ssid || '(隐藏网络)' }}
              <span v-if="network.connected" class="connected-badge">已连接</span>
            </div>
            <div class="network-meta">
              <span class="security">{{ getSecurityType(network.security) }}</span>
              <span class="frequency">{{ network.frequency }}</span>
              <span class="channel">频道 {{ network.channel }}</span>
            </div>
          </div>

          <div class="network-signal">
            <div class="signal-level" :class="getSignalStrength(network.signalLevel)">
              {{ Math.round(network.signalLevel) }}%
            </div>
          </div>
        </div>

        <div v-if="selectedNetwork === network.ssid" class="network-actions">
          <button
            @click="connectToNetwork(network)"
            :disabled="network.connected"
            class="connect-btn"
          >
            {{ network.connected ? '已连接' : '连接' }}
          </button>
          <button @click="forgetNetwork(network)" class="forget-btn">
            忘记网络
          </button>
        </div>
      </div>

      <div class="scan-summary">
        <span>找到 {{ networks.length }} 个可用网络</span>
        <span>{{ networks.filter(n => n.connected).length }} 个已连接</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  SignalIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'

interface WiFiNetwork {
  ssid: string
  bssid: string
  signalLevel: number
  security: string
  frequency: string
  channel: number
  secured: boolean
  connected: boolean
}

interface Props {
  value: WiFiNetwork[]
}

interface Emits {
  (e: 'update', value: WiFiNetwork[]): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const scanning = ref(false)
const error = ref('')
const networks = ref<WiFiNetwork[]>(props.value || [])
const selectedNetwork = ref<string | null>(null)

const scanNetworks = async () => {
  scanning.value = true
  error.value = ''

  try {
    const response = await fetch('/api/network/wifi/scan', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      }
    })

    if (!response.ok) {
      throw new Error('扫描Wi-Fi网络失败')
    }

    const data = await response.json()
    networks.value = data.networks || []
    emit('update', networks.value)

  } catch (err: any) {
    error.value = err.message || '扫描网络失败'
    console.error('Failed to scan WiFi networks:', err)
  } finally {
    scanning.value = false
  }
}

const selectNetwork = (network: WiFiNetwork) => {
  if (selectedNetwork.value === network.ssid) {
    selectedNetwork.value = null
  } else {
    selectedNetwork.value = network.ssid
  }
}

const connectToNetwork = async (network: WiFiNetwork) => {
  if (network.connected) {
    ElMessage.info('已经连接到此网络')
    return
  }

  // 如果是开放网络或已经保存密码的网络，直接连接
  if (!network.secured) {
    try {
      const response = await fetch('/api/network/wifi/connect', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          ssid: network.ssid,
          bssid: network.bssid
        })
      })

      if (!response.ok) {
        throw new Error('连接网络失败')
      }

      ElMessage.success('正在连接到网络...')
      await scanNetworks()

    } catch (err: any) {
      ElMessage.error(err.message || '连接失败')
      console.error('Failed to connect to WiFi:', err)
    }
  } else {
    // 需要密码的网络，提示用户在主设置中配置
    ElMessage.info('请在Wi-Fi密码设置中配置密码后再连接')
    selectedNetwork.value = null
  }
}

const forgetNetwork = async (network: WiFiNetwork) => {
  try {
    const response = await fetch('/api/network/wifi/forget', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        ssid: network.ssid
      })
    })

    if (!response.ok) {
      throw new Error('移除网络失败')
    }

    ElMessage.success('已移除网络配置')
    selectedNetwork.value = null
    await scanNetworks()

  } catch (err: any) {
    ElMessage.error(err.message || '移除网络失败')
    console.error('Failed to forget WiFi:', err)
  }
}

const getSignalStrength = (level: number): string => {
  if (level >= 75) return 'excellent'
  if (level >= 50) return 'good'
  if (level >= 25) return 'fair'
  return 'poor'
}

const getSecurityType = (security: string): string => {
  const types: Record<string, string> = {
    'WPA2-PSK': 'WPA2-Personal',
    'WPA3-PSK': 'WPA3-Personal',
    'WPA/WPA2': 'WPA/WPA2混合',
    'WEP': 'WEP',
    'Open': '开放网络'
  }
  return types[security] || security
}

// 组件挂载时自动扫描一次
onMounted(() => {
  // 可以选择自动扫描，或者等待用户手动点击
  // scanNetworks()
})
</script>

<style scoped lang="scss">
.wifi-scan-component {
  .scan-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h4 {
      margin: 0;
      font-size: 16px;
      font-weight: 600;
      color: #1f2937;
    }

    .scan-btn {
      display: flex;
      align-items: center;
      gap: 6px;
      padding: 8px 16px;
      border: 1px solid #d1d5db;
      border-radius: 6px;
      background: white;
      cursor: pointer;
      transition: all 0.2s;
      font-size: 13px;

      &:hover:not(:disabled) {
        background: #f3f4f6;
        border-color: #9ca3af;
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }

      .spinning {
        animation: spin 1s linear infinite;
      }
    }
  }

  .scanning-state, .error-state, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 32px;
    text-align: center;
    gap: 12px;
    color: #6b7280;

    .spinner {
      width: 32px;
      height: 32px;
      border: 3px solid #e5e7eb;
      border-top-color: #3b82f6;
      border-radius: 50%;
      animation: spin 1s linear infinite;
    }

    button {
      margin-top: 8px;
      padding: 6px 12px;
      background: #3b82f6;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      font-size: 12px;

      &:hover {
        background: #2563eb;
      }
    }
  }

  .error-state {
    color: #ef4444;
  }

  .networks-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .network-item {
    border: 1px solid #e5e7eb;
    border-radius: 8px;
    background: white;
    transition: all 0.2s;

    &:hover {
      border-color: #d1d5db;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    }

    &.connected {
      border-left: 4px solid #10b981;
    }

    &.selected {
      background: #f9fafb;
    }
  }

  .network-info {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    cursor: pointer;
    user-select: none;
  }

  .network-icon {
    position: relative;
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #6b7280;

    .signal-bars {
      position: absolute;
      bottom: 0;
      right: -8px;
      display: flex;
      gap: 2px;
      align-items: flex-end;

      .bar {
        width: 3px;
        background: #9ca3af;
        border-radius: 1px;

        &.bar1 { height: 4px; }
        &.bar2 { height: 8px; }
        &.bar3 { height: 12px; }
        &.bar4 { height: 16px; }
      }

      &.excellent .bar {
        background: #10b981;
      }

      &.good .bar1, &.good .bar2, &.good .bar3 {
        background: #10b981;
      }

      &.fair .bar1, &.fair .bar2 {
        background: #f59e0b;
      }

      &.poor .bar1 {
        background: #ef4444;
      }
    }
  }

  .network-details {
    flex: 1;
  }

  .network-name {
    font-weight: 500;
    color: #1f2937;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .connected-badge {
    font-size: 11px;
    padding: 2px 6px;
    background: #d1fae5;
    color: #065f46;
    border-radius: 4px;
    font-weight: 500;
  }

  .network-meta {
    display: flex;
    gap: 12px;
    margin-top: 4px;
    font-size: 12px;
    color: #6b7280;

    .security {
      position: relative;

      &::after {
        content: '•';
        position: absolute;
        right: -8px;
        color: #d1d5db;
      }
    }

    .frequency, .channel {
      position: relative;

      &::after {
        content: '•';
        position: absolute;
        right: -8px;
        color: #d1d5db;
      }
    }

    *:last-child::after {
      display: none;
    }
  }

  .network-signal {
    .signal-level {
      padding: 4px 8px;
      border-radius: 4px;
      font-size: 12px;
      font-weight: 500;

      &.excellent {
        background: #d1fae5;
        color: #065f46;
      }

      &.good {
        background: #d1fae5;
        color: #065f46;
      }

      &.fair {
        background: #fef3c7;
        color: #92400e;
      }

      &.poor {
        background: #fef2f2;
        color: #991b1b;
      }
    }
  }

  .network-actions {
    display: flex;
    gap: 8px;
    padding: 8px 12px 12px;
    border-top: 1px solid #e5e7eb;

    button {
      padding: 6px 12px;
      border: 1px solid #d1d5db;
      border-radius: 4px;
      background: white;
      cursor: pointer;
      transition: all 0.2s;
      font-size: 12px;

      &:hover:not(:disabled) {
        background: #f3f4f6;
        border-color: #9ca3af;
      }

      &:disabled {
        opacity: 0.5;
        cursor: not-allowed;
      }

      &.connect-btn {
        background: #3b82f6;
        color: white;
        border-color: #3b82f6;

        &:hover:not(:disabled) {
          background: #2563eb;
          border-color: #2563eb;
        }

        &:disabled {
          background: #9ca3af;
          border-color: #9ca3af;
        }
      }

      &.forget-btn {
        color: #ef4444;
        border-color: #fecaca;

        &:hover {
          background: #fef2f2;
          border-color: #ef4444;
        }
      }
    }
  }

  .scan-summary {
    display: flex;
    justify-content: space-between;
    padding: 8px 12px;
    font-size: 12px;
    color: #6b7280;
    background: #f9fafb;
    border-radius: 6px;
    margin-top: 8px;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>