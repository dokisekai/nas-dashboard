// SSO API 客户端
import axios from 'axios'

const SSO_BASE_URL = 'http://192.168.50.10:8888/sso'

export interface SSOConfig {
  issuerUrl: string
  clientId: string
  redirectUri: string
  scope: string
}

export interface SSOUserInfo {
  sub: string
  name: string
  email: string
  preferred_username: string
}

export interface SSOTokenResponse {
  access_token: string
  token_type: string
  expires_in: number
  refresh_token?: string
  id_token?: string
}

export interface SSOProvider {
  id: string
  name: string
  enabled: boolean
}

export class SSOClient {
  private config: SSOConfig

  constructor(config: SSOConfig) {
    this.config = config
  }

  /**
   * 生成授权URL
   */
  generateAuthUrl(provider: string, state: string): string {
    const params = new URLSearchParams({
      client_id: this.config.clientId,
      response_type: 'code',
      redirect_uri: this.config.redirectUri,
      scope: this.config.scope,
      state: state,
      provider: provider
    })

    return `${this.config.issuerUrl}/sso/authorize?${params.toString()}`
  }

  /**
   * 交换授权码获取访问令牌
   */
  async exchangeCodeForToken(code: string, redirectUri: string): Promise<SSOTokenResponse> {
    const params = new URLSearchParams({
      code: code,
      redirect_uri: redirectUri,
      client_id: this.config.clientId,
      grant_type: 'authorization_code'
    })

    const response = await axios.post(`${SSO_BASE_URL}/token`, params, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })

    return response.data
  }

  /**
   * 获取用户信息
   */
  async getUserInfo(accessToken: string): Promise<SSOUserInfo> {
    const response = await axios.get(`${SSO_BASE_URL}/userinfo`, {
      headers: {
        'Authorization': `Bearer ${accessToken}`
      }
    })

    return response.data
  }

  /**
   * 撤销令牌
   */
  async revokeToken(accessToken: string): Promise<void> {
    await axios.post(`${SSO_BASE_URL}/revoke`, {}, {
      headers: {
        'Authorization': `Bearer ${accessToken}`
      }
    })
  }

  /**
   * 令牌内省
   */
  async introspectToken(accessToken: string): Promise<any> {
    const response = await axios.post(`${SSO_BASE_URL}/introspect`,
      new URLSearchParams({ token: accessToken }),
      {
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        }
      }
    )

    return response.data
  }

  /**
   * 获取OIDC配置
   */
  async getOIDCConfiguration(): Promise<any> {
    const response = await axios.get(`${this.config.issuerUrl}/sso/.well-known/openid-configuration`)
    return response.data
  }

  /**
   * 刷新访问令牌
   */
  async refreshAccessToken(refreshToken: string): Promise<SSOTokenResponse> {
    const params = new URLSearchParams({
      refresh_token: refreshToken,
      client_id: this.config.clientId,
      grant_type: 'refresh_token'
    })

    const response = await axios.post(`${SSO_BASE_URL}/token`, params, {
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
      }
    })

    return response.data
  }
}

// 默认SSO客户端实例
export const ssoClient = new SSOClient({
  issuerUrl: 'http://192.168.50.10:8888',
  clientId: 'nas-dashboard',
  redirectUri: `${window.location.origin}/sso/callback`,
  scope: 'openid profile email'
})

/**
 * 生成随机state字符串
 */
export function generateRandomState(length: number = 32): string {
  const charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  let result = ''

  // 使用crypto.getRandomValues生成安全的随机数
  const randomValues = new Uint32Array(length)
  window.crypto.getRandomValues(randomValues)

  for (let i = 0; i < length; i++) {
    result += charset[randomValues[i] % charset.length]
  }

  return result
}

/**
 * 保存SSO会话信息
 */
export function saveSSOSession(tokens: SSOTokenResponse, userInfo: SSOUserInfo): void {
  const session = {
    access_token: tokens.access_token,
    refresh_token: tokens.refresh_token,
    id_token: tokens.id_token,
    expires_at: Date.now() + (tokens.expires_in * 1000),
    user: userInfo
  }

  localStorage.setItem('sso_session', JSON.stringify(session))
}

/**
 * 获取SSO会话信息
 */
export function getSSOSession(): any {
  const sessionStr = localStorage.getItem('sso_session')
  if (!sessionStr) {
    return null
  }

  try {
    const session = JSON.parse(sessionStr)

    // 检查会话是否过期
    if (session.expires_at < Date.now()) {
      clearSSOSession()
      return null
    }

    return session
  } catch (error) {
    console.error('Failed to parse SSO session:', error)
    clearSSOSession()
    return null
  }
}

/**
 * 清除SSO会话信息
 */
export function clearSSOSession(): void {
  localStorage.removeItem('sso_session')
  sessionStorage.removeItem('sso_state')
  sessionStorage.removeItem('sso_provider')
}

/**
 * 检查是否已通过SSO登录
 */
export function isSSOLoggedIn(): boolean {
  return getSSOSession() !== null
}

/**
 * 执行SSO登录
 */
export async function performSSOLogin(code: string, state: string): Promise<void> {
  try {
    // 1. 交换授权码获取访问令牌
    const tokens = await ssoClient.exchangeCodeForToken(code, window.location.href.split('?')[0])

    // 2. 获取用户信息
    const userInfo = await ssoClient.getUserInfo(tokens.access_token)

    // 3. 保存会话信息
    saveSSOSession(tokens, userInfo)

    // 4. 清理临时数据
    sessionStorage.removeItem('sso_state')
    sessionStorage.removeItem('sso_provider')

  } catch (error) {
    console.error('SSO login failed:', error)
    throw error
  }
}

/**
 * 执行SSO登出
 */
export async function performSSOLogout(): Promise<void> {
  try {
    const session = getSSOSession()
    if (session && session.access_token) {
      // 撤销访问令牌
      await ssoClient.revokeToken(session.access_token)
    }
  } catch (error) {
    console.error('SSO logout failed:', error)
  } finally {
    // 清除本地会话
    clearSSOSession()
  }
}

/**
 * 刷新访问令牌（如果即将过期）
 */
export async function refreshAccessTokenIfNeeded(): Promise<boolean> {
  const session = getSSOSession()
  if (!session) {
    return false
  }

  // 如果令牌将在5分钟内过期，则刷新
  const expiryTime = session.expires_at
  const refreshThreshold = 5 * 60 * 1000 // 5分钟

  if (expiryTime - Date.now() < refreshThreshold) {
    try {
      if (session.refresh_token) {
        const newTokens = await ssoClient.refreshAccessToken(session.refresh_token)
        saveSSOSession(newTokens, session.user)
        return true
      }
    } catch (error) {
      console.error('Failed to refresh access token:', error)
      // 刷新失败，清除会话
      clearSSOSession()
      return false
    }
  }

  return false
}