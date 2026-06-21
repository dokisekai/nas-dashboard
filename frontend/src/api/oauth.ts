import api from './client'

const API_BASE = '/api/oauth'

// 服务器信息
export interface ServerInfo {
  issuer_url: string
  authorize_endpoint: string
  token_endpoint: string
  userinfo_endpoint: string
  jwks_endpoint: string
  discovery_endpoint: string
  running: boolean
}

// OAuth客户端
export interface OAuthClient {
  id: number
  created_at: string
  updated_at: string
  name: string
  client_id: string
  client_secret: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
  status: string
  showSecret?: boolean
}

// OAuth授权
export interface OAuthAuthorization {
  id: number
  created_at: string
  updated_at: string
  user_id: number
  client_id: string
  scopes: string[]
  code: string
  expires_at: string
  revoked_at?: string
}

// 服务器统计
export interface ServerStats {
  active_users: number
  active_tokens: number
  today_auths: number
  total_clients: number
  total_auths: number
}

// 创建客户端请求
export interface CreateClientRequest {
  name: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
}

// 更新客户端请求
export interface UpdateClientRequest {
  name?: string
  redirect_uris?: string[]
  grant_types?: string[]
  scopes?: string[]
  status?: string
}

// API响应
interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

// 获取服务器信息
export async function getServerInfo(): Promise<ServerInfo> {
  const response = await api.get<ApiResponse<ServerInfo>>(`${API_BASE}/server/info`)
  return response.data!
}

// 获取客户端列表
export async function getClients(params?: {
  page?: number
  limit?: number
  search?: string
}): Promise<{ total: number; clients: OAuthClient[] }> {
  const response = await api.get<ApiResponse<{
    total: number
    clients: OAuthClient[]
  }>>(`${API_BASE}/clients`, { params })
  return response.data!
}

// 创建客户端
export async function createClient(data: CreateClientRequest): Promise<OAuthClient> {
  const response = await api.post<ApiResponse<OAuthClient>>(`${API_BASE}/clients`, data)
  return response.data!
}

// 更新客户端
export async function updateClient(
  id: number,
  data: UpdateClientRequest
): Promise<OAuthClient> {
  const response = await api.put<ApiResponse<OAuthClient>>(`${API_BASE}/clients/${id}`, data)
  return response.data!
}

// 删除客户端
export async function deleteClient(id: number): Promise<void> {
  await api.delete(`${API_BASE}/clients/${id}`)
}

// 重置客户端密钥
export async function regenerateSecret(id: number): Promise<{ client_id: string; client_secret: string }> {
  const response = await api.post<ApiResponse<{ client_id: string; client_secret: string }>>(
    `${API_BASE}/clients/${id}/regenerate-secret`,
    { confirm: true }
  )
  return response.data!
}

// 获取授权列表
export async function getAuthorizations(params?: {
  page?: number
  limit?: number
  user_id?: string
  client_id?: string
}): Promise<{ total: number; authorizations: OAuthAuthorization[] }> {
  const response = await api.get<ApiResponse<{
    total: number
    authorizations: OAuthAuthorization[]
  }>>(`${API_BASE}/authorizations`, { params })
  return response.data!
}

// 撤销授权
export async function revokeAuthorization(authorizationId: number): Promise<void> {
  await api.post(`${API_BASE}/authorizations/revoke`, { authorization_id: authorizationId })
}

// 获取服务器统计
export async function getServerStats(): Promise<ServerStats> {
  const response = await api.get<ApiResponse<ServerStats>>(`${API_BASE}/server/stats`)
  return response.data!
}

// 启动服务器
export async function startServer(): Promise<void> {
  await api.post(`${API_BASE}/server/start`)
}

// 停止服务器
export async function stopServer(): Promise<void> {
  await api.post(`${API_BASE}/server/stop`)
}

// 获取用户令牌
export async function getUserTokens(userId: number): Promise<any[]> {
  const response = await api.get<ApiResponse<any[]>>(`${API_BASE}/users/${userId}/tokens`)
  return response.data!
}

// 撤销用户令牌
export async function revokeUserToken(tokenId: number): Promise<void> {
  await api.delete(`${API_BASE}/tokens/${tokenId}`)
}

// 工具函数：解析JSON字段
export function parseJSONField<T>(field: string): T[] {
  if (!field) return []
  try {
    return JSON.parse(field)
  } catch {
    return []
  }
}

// 工具函数：转换为JSON字段
export function toJSONField<T>(data: T[]): string {
  return JSON.stringify(data)
}