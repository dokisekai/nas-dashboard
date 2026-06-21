<template>
  <div class="oauth-server">
    <!-- 加载状态 -->
    <div v-if="loading.serverInfo || loading.clients" class="loading-overlay">
      <div class="loading-spinner"></div>
      <div class="loading-text">正在加载OAuth服务器数据...</div>
    </div>

    <!-- 服务器状态头部 -->
    <div class="server-header">
      <div class="header-left">
        <svg class="header-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
        </svg>
        <div class="header-text">
          <h1>OAuth/OIDC 服务器</h1>
          <p class="subtitle">为其他应用提供统一身份认证服务</p>
        </div>
      </div>
      <div class="server-config">
        <button @click="showConfigModal = true" class="config-btn">
          <svg class="config-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.084 2.83-.065.26-3.025.326-5.867-.426-2.625-2.625-2.625.536-4.465.78-5.867.426-1.317.637-3.038 2.667 1.796 3.038 5.622.17 4.43c.026.894.055 1.794.083 2.687.169.378.362.957.73 1.334.34 1.574.383 3.064.17.47 4.433.098.483.984.584 1.443.048 2.944.12 4.408.161 1.035.268 2.067.427 3.15.2.566.976.616 1.363 2.24 1.796 2.478 2.562 2.749.098.485.298.098.485 0z" />
          </svg>
          服务器配置
        </button>
        <button @click="toggleServer" class="btn-toggle" :class="serverRunning ? 'stop' : 'start'">
          {{ serverRunning ? '停止服务' : '启动服务' }}
        </button>
      </div>
    </div>

    <!-- 配置模态框 -->
    <div v-if="showConfigModal" class="config-modal" @click.self="showConfigModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>服务器配置</h2>
          <button @click="showConfigModal = false" class="close-btn">
            <svg class="close-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="modal-body">
          <div class="config-item">
            <label>服务器URL</label>
            <input v-model="configForm.serverUrl" class="config-input" placeholder="http://192.168.50.10:8888">
          </div>
          <div class="config-item">
            <label>服务器端口</label>
            <input v-model="configForm.serverPort" type="number" class="config-input" placeholder="8888">
          </div>
          <div class="config-item">
            <label>JWT密钥</label>
            <input v-model="configForm.jwtSecret" type="password" class="config-input" placeholder="your-secret-key">
          </div>
          <div class="config-actions">
            <button @click="saveConfig" class="btn-primary">保存配置</button>
            <button @click="showConfigModal = false" class="btn-secondary">取消</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error.serverInfo || error.clients" class="error-banner">
      <svg class="error-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <span>加载数据失败: {{ error.serverInfo || error.clients }}</span>
      <button @click="initializeData" class="retry-btn">重试</button>
    </div>

    <!-- 服务器信息 -->
    <div class="server-info">
      <h2 class="section-title">服务器端点信息</h2>
      <div class="info-grid">
        <div class="info-card">
          <div class="info-label">Issuer URL</div>
          <div class="info-value">
            <code>{{ serverConfig.issuerUrl }}</code>
            <button @click="copyToClipboard(serverConfig.issuerUrl)" class="copy-btn">
              <svg class="copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        <div class="info-card">
          <div class="info-label">授权端点</div>
          <div class="info-value">
            <code>{{ serverConfig.issuerUrl }}/sso/authorize</code>
            <button @click="copyToClipboard(serverConfig.issuerUrl + '/sso/authorize')" class="copy-btn">
              <svg class="copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        <div class="info-card">
          <div class="info-label">Token端点</div>
          <div class="info-value">
            <code>{{ serverConfig.issuerUrl }}/sso/token</code>
            <button @click="copyToClipboard(serverConfig.issuerUrl + '/sso/token')" class="copy-btn">
              <svg class="copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        <div class="info-card">
          <div class="info-label">UserInfo端点</div>
          <div class="info-value">
            <code>{{ serverConfig.issuerUrl }}/sso/userinfo</code>
            <button @click="copyToClipboard(serverConfig.issuerUrl + '/sso/userinfo')" class="copy-btn">
              <svg class="copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        <div class="info-card">
          <div class="info-label">JWKS端点</div>
          <div class="info-value">
            <code>{{ serverConfig.issuerUrl }}/sso/jwks</code>
            <button @click="copyToClipboard(serverConfig.issuerUrl + '/sso/jwks')" class="copy-btn">
              <svg class="copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>

        <div class="info-card">
          <div class="info-label">OIDC发现端点</div>
          <div class="info-value">
            <code>{{ serverConfig.issuerUrl }}/sso/.well-known/openid-configuration</code>
            <button @click="copyToClipboard(serverConfig.issuerUrl + '/sso/.well-known/openid-configuration')" class="copy-btn">
              <svg class="copy-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 客户端应用管理 -->
    <div class="clients-section">
      <div class="section-header">
        <h2 class="section-title">客户端应用管理</h2>
        <button @click="addClient" class="btn-add">
          <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          添加客户端
        </button>
      </div>

      <div class="clients-list">
        <div
          v-for="client in oauthClients"
          :key="client.id"
          class="client-card"
        >
          <div class="client-header">
            <div class="client-info">
              <div class="client-name">{{ client.name }}</div>
              <div class="client-id">Client ID: {{ client.client_id }}</div>
            </div>
            <div class="client-status">
              <span class="status-badge" :class="client.status">{{ client.status === 'active' ? '活跃' : '禁用' }}</span>
            </div>
          </div>

          <div class="client-details">
            <div class="detail-row">
              <span class="detail-label">Client Secret:</span>
              <span class="detail-value">
                <code :class="{ masked: !client.showSecret }">{{ client.showSecret ? client.client_secret : '•••••••••••••••' }}</code>
                <button @click="toggleSecret(client)" class="toggle-btn">
                  <svg class="toggle-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                </button>
              </span>
            </div>

            <div class="detail-row">
              <span class="detail-label">重定向URI:</span>
              <span class="detail-value">
                <code>{{ client.redirect_uris.join(', ') }}</code>
              </span>
            </div>

            <div class="detail-row">
              <span class="detail-label">授权类型:</span>
              <span class="detail-value">
                <span class="tag">{{ client.grant_types.join(', ') }}</span>
              </span>
            </div>

            <div class="detail-row">
              <span class="detail-label">授权范围:</span>
              <span class="detail-value">
                <span class="tag">{{ client.scopes.join(', ') }}</span>
              </span>
            </div>
          </div>

          <div class="client-actions">
            <button @click="editClient(client)" class="action-btn edit">
              <svg class="action-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.084 2.83-.065.26-3.025.326-5.867-.426-2.625-2.625-2.625.536-4.465.78-5.867.426-1.317.637-3.038 2.667 1.796 3.038 5.622.17 4.43c.026.894.055 1.794.083 2.687.169.378.362.957.73 1.334.34 1.574.383 3.064.17.47 4.433.098.483.984.584 1.443.048 2.944.12 4.408.161 1.035.268 2.067.427 3.15.2.566.976.616 1.363 2.24 1.796 2.478 2.562 2.749.098.485.298.098.485 0z" />
              </svg>
              编辑
            </button>
            <button @click="regenerateSecret(client)" class="action-btn warning">
              <svg class="action-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              重置密钥
            </button>
            <button @click="deleteClient(client.id)" class="action-btn danger">
              <svg class="action-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              删除
            </button>
          </div>
        </div>
      </div>

      <div v-if="oauthClients.length === 0 && !loading.clients" class="empty-state">
        <svg class="empty-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <h3>暂无客户端应用</h3>
        <p>添加客户端应用开始使用OAuth认证服务</p>
      </div>
    </div>

    <!-- 服务器统计 -->
    <div class="stats-section">
      <h2 class="section-title">服务器统计</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.active_users }}</div>
            <div class="stat-label">活跃用户</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964a6 6 0 117.071-7.071z" />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.active_tokens }}</div>
            <div class="stat-label">活跃令牌</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.today_auths }}</div>
            <div class="stat-label">今日认证</div>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
            </svg>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.total_clients }}</div>
            <div class="stat-label">注册客户端</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as oauthAPI from '../api/oauth'

// 服务器配置
const serverConfig = ref({
  issuerUrl: 'http://192.168.50.10:8888',
  port: 8888,
  tokenExpiry: 3600,
  refreshTokenExpiry: 2592000
})

// 服务器状态
const serverRunning = ref(true)

// OAuth客户端应用
const oauthClients = ref<oauthAPI.OAuthClient[]>([])

// 服务器统计
const stats = ref<oauthAPI.ServerStats>({
  active_users: 0,
  active_tokens: 0,
  today_auths: 0,
  total_clients: 0,
  total_auths: 0
})

// 加载状态
const loading = ref({
  serverInfo: false,
  clients: false,
  authorizations: false,
  stats: false
})

// 错误状态
const error = ref({
  serverInfo: null as string | null,
  clients: null as string | null,
  authorizations: null as string | null,
  stats: null as string | null
})

// 自动刷新定时器
let refreshInterval: number | null = null

// 配置模态框
const showConfigModal = ref(false)
const configForm = ref({
  serverUrl: 'https://192.168.50.10:8888',
  serverPort: '8888',
  jwtSecret: 'your-super-secret-jwt-key-change-this'
})

// 方法
const saveConfig = async () => {
  try {
    // 保存配置到本地存储
    localStorage.setItem('oauth_server_config', JSON.stringify(configForm.value))

    // 更新当前服务器配置
    serverConfig.value.issuerUrl = configForm.value.serverUrl
    serverConfig.value.port = parseInt(configForm.value.serverPort)

    // 这里可以调用后端API保存全局配置
    console.log('服务器配置已保存:', configForm.value)

    showConfigModal.value = false
    alert('服务器配置已保存，服务器需要重启才能生效')
  } catch (err: any) {
    console.error('保存配置失败:', err)
    alert('保存配置失败: ' + err.message)
  }
}

const loadServerConfig = () => {
  const savedConfig = localStorage.getItem('oauth_server_config')
  if (savedConfig) {
    try {
      const config = JSON.parse(savedConfig)
      configForm.value = config
      serverConfig.value.issuerUrl = config.serverUrl
      serverConfig.value.port = parseInt(config.serverPort)
      console.log('已加载保存的服务器配置:', config)
    } catch (err) {
      console.error('加载配置失败:', err)
    }
  } else {
    // 自动检测服务器地址
    autoDetectServerConfig()
  }
}

const autoDetectServerConfig = () => {
  // 从当前浏览器地址自动检测
  const currentHost = window.location.hostname
  const currentPort = window.location.port || '8888'

  // 始终使用NAS地址
  const protocol = 'https'
  configForm.value.serverUrl = `${protocol}://192.168.50.10:8888`

  serverConfig.value.issuerUrl = configForm.value.serverUrl
  console.log('自动检测服务器配置:', serverConfig.value.issuerUrl)
}

const loadServerInfo = async () => {
  loading.value.serverInfo = true
  error.value.serverInfo = null

  try {
    const info = await oauthAPI.getServerInfo()
    serverConfig.value.issuerUrl = info.issuer_url
    serverRunning.value = info.running
    console.log('服务器信息已加载:', info)
  } catch (err: any) {
    console.error('加载服务器信息失败:', err)
    error.value.serverInfo = err.message || '加载服务器信息失败'
  } finally {
    loading.value.serverInfo = false
  }
}

const loadClients = async () => {
  loading.value.clients = true
  error.value.clients = null

  try {
    const result = await oauthAPI.getClients()
    oauthClients.value = result.clients.map(client => ({
      ...client,
      redirect_uris: oauthAPI.parseJSONField(client.redirect_uris as any),
      grant_types: oauthAPI.parseJSONField(client.grant_types as any),
      scopes: oauthAPI.parseJSONField(client.scopes as any),
      showSecret: false
    }))
    console.log('客户端列表已加载:', oauthClients.value.length)
  } catch (err: any) {
    console.error('加载客户端列表失败:', err)
    error.value.clients = err.message || '加载客户端列表失败'
  } finally {
    loading.value.clients = false
  }
}

const loadStats = async () => {
  loading.value.stats = true
  error.value.stats = null

  try {
    stats.value = await oauthAPI.getServerStats()
    console.log('服务器统计已加载:', stats.value)
  } catch (err: any) {
    console.error('加载服务器统计失败:', err)
    error.value.stats = err.message || '加载服务器统计失败'
  } finally {
    loading.value.stats = false
  }
}

const toggleServer = async () => {
  try {
    if (serverRunning.value) {
      await oauthAPI.stopServer()
      serverRunning.value = false
      console.log('服务器已停止')
    } else {
      await oauthAPI.startServer()
      serverRunning.value = true
      console.log('服务器已启动')
    }
    // 重新加载服务器信息
    await loadServerInfo()
  } catch (err: any) {
    console.error('切换服务器状态失败:', err)
    alert('操作失败: ' + err.message)
  }
}

const copyToClipboard = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    console.log('已复制:', text)
    // 可以添加一个临时提示消息
  })
}

const toggleSecret = (client: oauthAPI.OAuthClient) => {
  client.showSecret = !client.showSecret
}

const addClient = async () => {
  const clientName = prompt('输入客户端应用名称:')
  if (!clientName) return

  const redirectUri = prompt('输入重定向URI (用逗号分隔多个):', 'https://192.168.50.10:3000/callback')
  if (!redirectUri) return

  try {
    const newClient = await oauthAPI.createClient({
      name: clientName,
      redirect_uris: redirectUri.split(',').map(uri => uri.trim()),
      grant_types: ['authorization_code', 'refresh_token'],
      scopes: ['openid', 'profile', 'email']
    })

    // 添加到列表
    oauthClients.value.push({
      ...newClient,
      redirect_uris: oauthAPI.parseJSONField(newClient.redirect_uris as any),
      grant_types: oauthAPI.parseJSONField(newClient.grant_types as any),
      scopes: oauthAPI.parseJSONField(newClient.scopes as any),
      showSecret: false
    })

    console.log('添加客户端成功:', newClient)

    alert(`客户端 ${clientName} 已创建！\n\nClient ID: ${newClient.client_id}\nClient Secret: ${newClient.client_secret}`)
  } catch (err: any) {
    console.error('创建客户端失败:', err)
    alert('创建客户端失败: ' + err.message)
  }
}

const editClient = async (client: oauthAPI.OAuthClient) => {
  const newName = prompt('修改客户端名称:', client.name)
  const newRedirectUri = prompt('修改重定向URI (用逗号分隔多个):', client.redirect_uris.join(','))

  if (!newName && !newRedirectUri) {
    return
  }

  try {
    const updateData: oauthAPI.UpdateClientRequest = {}
    if (newName) updateData.name = newName
    if (newRedirectUri) updateData.redirect_uris = newRedirectUri.split(',').map(uri => uri.trim())

    const updatedClient = await oauthAPI.updateClient(client.id, updateData)

    // 更新列表中的客户端
    const index = oauthClients.value.findIndex(c => c.id === client.id)
    if (index !== -1) {
      oauthClients.value[index] = {
        ...updatedClient,
        redirect_uris: oauthAPI.parseJSONField(updatedClient.redirect_uris as any),
        grant_types: oauthAPI.parseJSONField(updatedClient.grant_types as any),
        scopes: oauthAPI.parseJSONField(updatedClient.scopes as any),
        showSecret: client.showSecret
      }
    }

    console.log('编辑客户端成功:', updatedClient)
    alert('客户端已更新成功')
  } catch (err: any) {
    console.error('编辑客户端失败:', err)
    alert('编辑客户端失败: ' + err.message)
  }
}

const regenerateSecret = async (client: oauthAPI.OAuthClient) => {
  if (confirm(`确定要为 ${client.name} 重置 Client Secret 吗？此操作不可逆。`)) {
    try {
      const result = await oauthAPI.regenerateSecret(client.id)
      console.log('重置密钥成功:', result)
      alert(`新的 Client Secret: ${result.client_secret}\n\n请妥善保存，关闭此窗口后将无法再次查看。`)

      // 重新加载客户端列表
      await loadClients()
    } catch (err: any) {
      console.error('重置密钥失败:', err)
      alert('重置密钥失败: ' + err.message)
    }
  }
}

const deleteClient = async (clientId: number) => {
  if (confirm('确定要删除此客户端吗？此操作不可逆。')) {
    try {
      await oauthAPI.deleteClient(clientId)
      oauthClients.value = oauthClients.value.filter(c => c.id !== clientId)
      console.log('删除客户端成功:', clientId)
      alert('客户端已删除成功')
    } catch (err: any) {
      console.error('删除客户端失败:', err)
      alert('删除客户端失败: ' + err.message)
    }
  }
}

// 初始化加载
const initializeData = async () => {
  await Promise.all([
    loadServerInfo(),
    loadClients(),
    loadStats()
  ])
}

// 生命周期
onMounted(async () => {
  console.log('OAuth服务器管理器已加载')

  // 加载服务器配置
  loadServerConfig()

  // 初始化数据
  await initializeData()

  // 设置自动刷新 (每30秒)
  refreshInterval = window.setInterval(() => {
    loadStats()
  }, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.oauth-server {
  padding: 24px;
  background: white;
  border-radius: 16px;
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  overflow-y: auto;
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.95);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top-color: #667eea;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.loading-text {
  margin-top: 16px;
  font-size: 16px;
  color: #667eea;
  font-weight: 500;
}

.error-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: #fee2e2;
  border-left: 4px solid #ef4444;
  border-radius: 8px;
  margin-bottom: 20px;
  color: #991b1b;
}

.error-icon {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
}

.retry-btn {
  margin-left: auto;
  padding: 6px 12px;
  background: white;
  border: 1px solid #ef4444;
  border-radius: 6px;
  color: #ef4444;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.retry-btn:hover {
  background: #ef4444;
  color: white;
}

.server-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  padding-bottom: 20px;
  border-bottom: 2px solid #f0f0f0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-icon {
  width: 48px;
  height: 48px;
  color: #667eea;
}

.header-text h1 {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin: 0;
}

.subtitle {
  font-size: 14px;
  color: #666;
  margin: 4px 0 0 0;
}

.server-status {
  display: flex;
  align-items: center;
  gap: 16px;
}

.status-badge {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 500;
}

.status-badge.running {
  background: #d1fae5;
  color: #065f46;
}

.status-badge:not(.running) {
  background: #e5e7eb;
  color: #374151;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status-badge.running .status-dot {
  background: #10b981;
}

.status-badge:not(.running) .status-dot {
  background: #9ca3af;
}

.btn-toggle {
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn-toggle.start {
  background: #10b981;
  color: white;
}

.btn-toggle.stop {
  background: #ef4444;
  color: white;
}

.btn-toggle:hover {
  opacity: 0.9;
}

.server-config {
  display: flex;
  align-items: center;
  gap: 12px;
}

.config-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: 1px solid #e5e7eb;
  background: white;
  color: #374151;
  transition: all 0.2s;
}

.config-btn:hover {
  background: #f9fafb;
  border-color: #667eea;
  color: #667eea;
}

.config-icon {
  width: 16px;
  height: 16px;
}

/* 配置模态框 */
.config-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(10px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal-content {
  background: white;
  border-radius: 16px;
  box-shadow: 0 25px 80px rgba(0, 0, 0, 0.4);
  width: 90%;
  max-width: 500px;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 2px solid #f0f0f0;
}

.modal-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.close-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: none;
  background: #f3f4f6;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #e5e7eb;
}

.close-icon {
  width: 20px;
  height: 20px;
  color: #6b7280;
}

.modal-body {
  padding: 24px;
}

.config-item {
  margin-bottom: 20px;
}

.config-item label {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.config-input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  transition: all 0.2s;
}

.config-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.config-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}

.btn-primary {
  background: #667eea;
  color: white;
}

.btn-primary:hover {
  background: #5568d3;
}

.btn-secondary {
  background: #e5e7eb;
  color: #374151;
}

.btn-secondary:hover {
  background: #d1d5db;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin-bottom: 20px;
}

.server-info,
.clients-section,
.stats-section {
  margin-bottom: 32px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.info-card {
  background: #f9fafb;
  border-radius: 8px;
  padding: 16px;
}

.info-label {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.info-value {
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-value code {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #667eea;
  background: white;
  padding: 8px;
  border-radius: 4px;
  word-break: break-all;
}

.copy-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #6b7280;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: #e5e7eb;
  color: #374151;
}

.copy-icon {
  width: 16px;
  height: 16px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.btn-add {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
  background: #667eea;
  color: white;
}

.btn-add:hover {
  background: #5568d3;
}

.btn-icon {
  width: 16px;
  height: 16px;
}

.clients-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.client-card {
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s ease;
}

.client-card:hover {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.client-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.client-name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.client-id {
  font-size: 13px;
  color: #666;
  font-family: 'Courier New', monospace;
}

.client-status .status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.status-badge.active {
  background: #d1fae5;
  color: #065f46;
}

.client-details {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
}

.detail-label {
  font-size: 14px;
  font-weight: 500;
  color: #374151;
}

.detail-value {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
  margin-left: 20px;
}

.detail-value code {
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #667eea;
  background: white;
  padding: 6px 10px;
  border-radius: 4px;
}

.detail-value code.masked {
  letter-spacing: 2px;
}

.toggle-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: #6b7280;
  padding: 4px;
  border-radius: 4px;
}

.toggle-btn:hover {
  background: #e5e7eb;
}

.toggle-icon {
  width: 16px;
  height: 16px;
}

.tag {
  display: inline-block;
  padding: 4px 10px;
  background: #e0e7ff;
  color: #3730a3;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.client-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  border: 1px solid #d1d5db;
  background: white;
  color: #374151;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f9fafb;
  border-color: #9ca3af;
}

.action-btn.warning {
  color: #f59e0b;
  border-color: #f59e0b;
}

.action-btn.danger {
  color: #ef4444;
  border-color: #ef4444;
}

.action-icon {
  width: 16px;
  height: 16px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: #9ca3af;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 8px;
}

.empty-state p {
  font-size: 14px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.stat-icon {
  width: 48px;
  height: 48px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon svg {
  width: 24px;
  height: 24px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
}

/* 响应式 */
@media (max-width: 768px) {
  .oauth-server {
    padding: 16px;
  }

  .server-header {
    flex-direction: column;
    gap: 16px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .detail-row {
    flex-direction: column;
    align-items: flex-start;
  }

  .detail-value {
    margin-left: 0;
    margin-top: 8px;
  }
}
</style>