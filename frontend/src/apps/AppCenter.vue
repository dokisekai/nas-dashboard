<template>
  <div class="app-center">
    <!-- Header -->
    <div class="ac-header">
      <div>
        <h1>应用管理中心</h1>
        <p class="subtitle">统一管理本 NAS 上的 Docker 服务（{{ totalRunning }}/{{ totalApps }} 个运行中）</p>
      </div>
      <div class="ac-header-right">
        <div class="ping-badge" :class="dockerAvailable ? 'ok' : 'err'">
          <span class="dot"></span>{{ dockerAvailable ? 'Docker 已连接' : 'Docker 未连接' }}
        </div>
        <button class="ac-btn ghost" @click="refreshAll" :disabled="loading">
          <ArrowPathIcon class="w-4 h-4" :class="{ spinning: loading }" /> 刷新
        </button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="ac-tabs">
      <button v-for="t in tabs" :key="t.id" class="ac-tab" :class="{ active: tab === t.id }" @click="tab = t.id">
        <component :is="t.icon" class="w-4 h-4" />
        {{ t.label }}
        <span v-if="t.count !== undefined && t.count > 0" class="ac-badge">{{ t.count }}</span>
      </button>
    </div>

    <!-- =================== Catalog (服务目录) =================== -->
    <div v-show="tab === 'catalog'" class="ac-panel">
      <div class="ac-toolbar">
        <h2>我的应用</h2>
        <span class="ac-hint">系统自动检测到的服务，含状态、入口 URL 与一键控制</span>
      </div>

      <div v-if="loading && catalog.length === 0" class="ac-empty"><CogIcon class="w-8 h-8 spinning" /> 加载中…</div>
      <div v-else class="catalog-grid">
        <div v-for="svc in catalog" :key="svc.id" class="catalog-card" :class="`status-${svc.status}`">
          <div class="catalog-head">
            <div class="catalog-icon" :class="`icon-${svc.iconHint}`">
              <component :is="catalogIcon(svc.iconHint)" class="w-6 h-6" />
            </div>
            <div class="catalog-title">
              <h3>{{ svc.name }}</h3>
              <div class="catalog-cat">{{ categoryLabel(svc.category) }}</div>
            </div>
            <span class="catalog-status" :class="`st-${svc.status}`">
              <span class="status-dot"></span>{{ statusLabel(svc.status) }}
            </span>
          </div>

          <p class="catalog-desc">{{ svc.description }}</p>

          <div class="catalog-meta">
            <div class="meta-row" v-if="svc.url">
              <GlobeAltIcon class="w-4 h-4" />
              <a :href="svc.url" target="_blank" rel="noopener">{{ svc.url }}</a>
            </div>
            <div class="meta-row" v-if="svc.containers.length > 0">
              <CubeIcon class="w-4 h-4" />
              <span>{{ svc.runningCount }}/{{ svc.totalCount }} 容器运行</span>
            </div>
            <div class="meta-row" v-if="svc.composeDir">
              <FolderIcon class="w-4 h-4" />
              <code :title="svc.composeFile">{{ svc.composeDir }}</code>
            </div>
          </div>

          <div v-if="svc.notes" class="catalog-note">
            <InformationCircleIcon class="w-4 h-4" />
            <span>{{ svc.notes }}</span>
          </div>

          <div class="catalog-actions">
            <button v-if="svc.url" class="ac-btn primary sm" @click="openUrl(svc.url)">
              <ArrowTopRightOnSquareIcon class="w-4 h-4" /> 打开
            </button>
            <button v-if="svc.status !== 'missing' && svc.status !== 'running'"
              class="ac-btn primary sm" @click="startService(svc)" :disabled="busy[`svc-${svc.id}-start`]">
              <PlayIcon class="w-4 h-4" /> 启动
            </button>
            <button v-if="svc.status === 'running' || svc.status === 'partial'"
              class="ac-btn ghost sm" @click="stopService(svc)" :disabled="busy[`svc-${svc.id}-stop`]">
              <StopIcon class="w-4 h-4" /> 停止
            </button>
            <button v-if="svc.status !== 'missing'"
              class="ac-btn ghost sm" @click="restartService(svc)" :disabled="busy[`svc-${svc.id}-restart`]">
              <ArrowPathIcon class="w-4 h-4" /> 重启
            </button>
            <button v-if="svc.containers.length > 0" class="ac-btn ghost sm" @click="viewServiceContainers(svc)">
              <CubeIcon class="w-4 h-4" /> 详情
            </button>
          </div>

          <!-- 子容器列表（折叠） -->
          <details v-if="svc.containers.length > 0" class="catalog-details">
            <summary>{{ svc.containers.length }} 个容器</summary>
            <div class="sub-containers">
              <div v-for="c in svc.containers" :key="c.id" class="sub-container">
                <span class="sub-status" :class="c.state"></span>
                <code>{{ c.name }}</code>
                <span class="muted">·</span>
                <span class="muted">{{ c.composeService || c.image.slice(0, 30) }}</span>
                <button class="ac-btn ghost xs" @click.stop="openLogs(c.name)">日志</button>
              </div>
            </div>
          </details>
        </div>
      </div>
    </div>

    <!-- =================== Containers (容器列表) =================== -->
    <div v-show="tab === 'containers'" class="ac-panel">
      <div class="ac-toolbar">
        <h2>全部容器</h2>
        <input v-model="containerFilter" class="ac-input search" placeholder="按名称/镜像/项目过滤…" />
        <select v-model="containerStateFilter" class="ac-input">
          <option value="">全部状态</option>
          <option value="running">运行中</option>
          <option value="exited">已停止</option>
        </select>
        <button class="ac-btn ghost sm" @click="loadContainers(true)" :disabled="loading">
          <ArrowPathIcon class="w-4 h-4" /> 刷新
        </button>
      </div>

      <div v-if="loading && containers.length === 0" class="ac-empty"><CogIcon class="w-8 h-8 spinning" /> 加载中…</div>
      <div v-else-if="filteredContainers.length === 0" class="ac-empty"><CubeIcon class="w-12 h-12" /><p>没有匹配的容器</p></div>

      <table v-else class="ac-table">
        <thead>
          <tr>
            <th>名称</th><th>镜像</th><th>项目 / 服务</th><th>状态</th><th>端口</th>
            <th>启动时间</th><th>CPU / 内存</th><th class="col-actions">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in filteredContainers" :key="c.id">
            <td>
              <div class="cell-name">
                <span class="sub-status" :class="c.state"></span>
                <div>
                  <div>{{ c.name }}</div>
                  <div class="muted xs">{{ c.composeService || '—' }}</div>
                </div>
              </div>
            </td>
            <td class="cell-image" :title="c.image">{{ c.image }}</td>
            <td>
              <span v-if="c.composeProject" class="ac-chip">{{ c.composeProject }}</span>
              <span v-else class="muted">standalone</span>
            </td>
            <td>
              <span class="ac-chip" :class="`st-task-${c.state === 'running' ? 'completed' : c.state === 'exited' && c.exitCode === 0 ? 'idle' : 'failed'}`">
                {{ c.status }}
              </span>
              <div v-if="c.health" class="muted xs">健康: {{ c.health }}</div>
            </td>
            <td>
              <div v-for="(p, i) in c.ports" :key="i" class="port-cell">
                <a v-if="p.url" :href="p.url" target="_blank" rel="noopener" class="port-url">:{{ p.hostPort }}↗</a>
                <span v-else class="muted">:{{ p.hostPort || '-' }}</span>
                <span class="muted xs">→{{ p.containerPort }}/{{ p.protocol }}</span>
              </div>
            </td>
            <td class="muted">{{ c.startedAt ? formatDate(c.startedAt) : '-' }}</td>
            <td>
              <div v-if="c.stats">
                <div>CPU {{ c.stats.cpuPercent.toFixed(1) }}%</div>
                <div class="muted xs">内存 {{ formatBytes(c.stats.memoryUsage) }} / {{ formatBytes(c.stats.memoryLimit) }}</div>
              </div>
              <button v-else class="ac-btn ghost xs" @click="loadStats(c)">查看</button>
            </td>
            <td class="col-actions">
              <button v-if="c.state !== 'running'" class="ac-btn primary xs" @click="doAction(c, 'start')" title="启动">
                <PlayIcon class="w-3.5 h-3.5" />
              </button>
              <button v-if="c.state === 'running'" class="ac-btn ghost xs" @click="doAction(c, 'stop')" title="停止">
                <StopIcon class="w-3.5 h-3.5" />
              </button>
              <button class="ac-btn ghost xs" @click="doAction(c, 'restart')" title="重启">
                <ArrowPathIcon class="w-3.5 h-3.5" />
              </button>
              <button class="ac-btn ghost xs" @click="openLogs(c.name)" title="日志">
                <DocumentTextIcon class="w-3.5 h-3.5" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- =================== Projects (Compose 项目) =================== -->
    <div v-show="tab === 'projects'" class="ac-panel">
      <div class="ac-toolbar">
        <h2>Compose 项目</h2>
        <button class="ac-btn ghost sm" @click="loadProjects" :disabled="loading">
          <ArrowPathIcon class="w-4 h-4" /> 刷新
        </button>
      </div>
      <div v-if="loading && projects.length === 0" class="ac-empty"><CogIcon class="w-8 h-8 spinning" /> 加载中…</div>
      <div v-else-if="projects.length === 0" class="ac-empty"><CubeIcon class="w-12 h-12" /><p>未发现 Compose 项目</p></div>

      <div v-else class="projects-grid">
        <div v-for="p in projects" :key="p.name" class="project-card">
          <div class="project-head">
            <div class="catalog-icon sm" :class="`icon-${p.iconHint}`">
              <component :is="catalogIcon(p.iconHint)" class="w-5 h-5" />
            </div>
            <div class="catalog-title">
              <h3>{{ p.name }}</h3>
              <div class="catalog-cat">{{ categoryLabel(p.category) }}</div>
            </div>
            <span class="catalog-status" :class="p.runningCount === p.totalCount ? 'st-running' : p.runningCount > 0 ? 'st-partial' : 'st-stopped'">
              {{ p.runningCount }}/{{ p.totalCount }}
            </span>
          </div>

          <p v-if="p.description" class="catalog-desc">{{ p.description }}</p>

          <div class="catalog-meta">
            <div class="meta-row" v-if="p.configFile">
              <DocumentTextIcon class="w-4 h-4" />
              <code :title="p.configFile">{{ p.configFile.split('/').slice(-2).join('/') }}</code>
            </div>
            <div class="meta-row" v-if="p.workingDir">
              <FolderIcon class="w-4 h-4" />
              <code>{{ p.workingDir }}</code>
            </div>
          </div>

          <div class="containers-list">
            <div v-for="name in p.containers" :key="name" class="sub-container">
              <span class="muted">·</span>
              <code>{{ name }}</code>
              <button class="ac-btn ghost xs" @click="openLogs(name)">日志</button>
            </div>
          </div>

          <div class="catalog-actions">
            <button class="ac-btn primary sm" @click="startProject(p)" :disabled="busy[`proj-${p.name}-start`]">
              <PlayIcon class="w-4 h-4" /> 全部启动
            </button>
            <button class="ac-btn ghost sm" @click="stopProject(p)" :disabled="busy[`proj-${p.name}-stop`]">
              <StopIcon class="w-4 h-4" /> 全部停止
            </button>
            <button class="ac-btn ghost sm" @click="restartProject(p)" :disabled="busy[`proj-${p.name}-restart`]">
              <ArrowPathIcon class="w-4 h-4" /> 全部重启
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Logs modal -->
    <div v-if="modals.logs" class="ac-overlay" @click.self="closeLogs">
      <div class="ac-modal lg">
        <div class="ac-modal-head">
          <h3>容器日志：{{ logName }}</h3>
          <button class="ac-x" @click="closeLogs"><XMarkIcon class="w-5 h-5" /></button>
        </div>
        <div class="ac-modal-body">
          <div class="logs-toolbar">
            <span class="muted">最近 {{ logTail }} 行</span>
            <button class="ac-btn ghost xs" @click="reloadLogs"><ArrowPathIcon class="w-3 h-3" /> 刷新</button>
          </div>
          <pre class="ac-log">{{ logContent || '（无日志）' }}</pre>
        </div>
        <div class="ac-modal-foot">
          <button class="ac-btn ghost" @click="closeLogs">关闭</button>
        </div>
      </div>
    </div>

    <!-- Service detail modal -->
    <div v-if="modals.service" class="ac-overlay" @click.self="modals.service = false">
      <div class="ac-modal lg">
        <div class="ac-modal-head">
          <h3>{{ selectedService?.name }} · 容器详情</h3>
          <button class="ac-x" @click="modals.service = false"><XMarkIcon class="w-5 h-5" /></button>
        </div>
        <div class="ac-modal-body">
          <table class="ac-table compact">
            <thead>
              <tr><th>名称</th><th>服务</th><th>状态</th><th>镜像</th><th class="col-actions">操作</th></tr>
            </thead>
            <tbody>
              <tr v-for="c in selectedService?.containers || []" :key="c.id">
                <td><span class="sub-status" :class="c.state"></span> {{ c.name }}</td>
                <td><code>{{ c.composeService || '-' }}</code></td>
                <td><span class="ac-chip" :class="c.state === 'running' ? 'st-task-completed' : 'st-task-idle'">{{ c.status }}</span></td>
                <td class="cell-image">{{ c.image }}</td>
                <td class="col-actions">
                  <button v-if="c.state !== 'running'" class="ac-btn primary xs" @click="doAction(c, 'start')">启动</button>
                  <button v-if="c.state === 'running'" class="ac-btn ghost xs" @click="doAction(c, 'stop')">停止</button>
                  <button class="ac-btn ghost xs" @click="openLogs(c.name)">日志</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="ac-modal-foot">
          <button class="ac-btn ghost" @click="modals.service = false">关闭</button>
        </div>
      </div>
    </div>

    <!-- Toast -->
    <transition name="fade">
      <div v-if="toast.show" class="ac-toast" :class="`toast-${toast.type}`">{{ toast.text }}</div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { appsApi } from '@/api'
import type { AppContainer, ComposeProject, ServiceCatalogEntry } from '@/api'
import {
  ArrowPathIcon, CogIcon, CubeIcon, GlobeAltIcon, FolderIcon, InformationCircleIcon,
  ArrowTopRightOnSquareIcon, PlayIcon, StopIcon, DocumentTextIcon, XMarkIcon,
  CircleStackIcon, CloudIcon, CodeBracketIcon, MusicalNoteIcon, PhotoIcon,
  ServerStackIcon, ArchiveBoxIcon, ShoppingCartIcon,
} from '@heroicons/vue/24/outline'

const loading = ref(false)
const dockerAvailable = ref(true)
const tab = ref<'catalog' | 'containers' | 'projects'>('catalog')
const catalog = ref<ServiceCatalogEntry[]>([])
const containers = ref<AppContainer[]>([])
const projects = ref<ComposeProject[]>([])
const busy = reactive<Record<string, boolean>>({})
const containerFilter = ref('')
const containerStateFilter = ref('')

// Logs modal
const modals = reactive({ logs: false, service: false })
const logName = ref('')
const logContent = ref('')
const logTail = ref(500)
const selectedService = ref<ServiceCatalogEntry | null>(null)

// Toast
const toast = reactive({ show: false, text: '', type: 'info' as 'info' | 'success' | 'error' })
let toastTimer: any = null
function showToast(text: string, type: 'info' | 'success' | 'error' = 'info') {
  toast.text = text; toast.type = type; toast.show = true
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.show = false), 3500)
}

const totalApps = computed(() => catalog.value.length)
const totalRunning = computed(() => catalog.value.filter(s => s.status === 'running').length)

const filteredContainers = computed(() => {
  const q = containerFilter.value.toLowerCase().trim()
  return containers.value.filter(c => {
    if (containerStateFilter.value && c.state !== containerStateFilter.value) return false
    if (!q) return true
    return c.name.toLowerCase().includes(q) ||
      c.image.toLowerCase().includes(q) ||
      (c.composeProject || '').toLowerCase().includes(q)
  })
})

const tabs = computed(() => [
  { id: 'catalog' as const, label: '我的应用', icon: ShoppingCartIcon, count: catalog.value.length },
  { id: 'containers' as const, label: '容器', icon: CubeIcon, count: containers.value.length },
  { id: 'projects' as const, label: 'Compose 项目', icon: CircleStackIcon, count: projects.value.length },
])

let pollTimer: any = null
onMounted(() => {
  refreshAll()
  // 每 30 秒静默刷新（不显示 loading）
  pollTimer = setInterval(() => refreshAll(true), 30000)
})
onUnmounted(() => { if (pollTimer) clearInterval(pollTimer) })

async function refreshAll(silent = false) {
  if (!silent) loading.value = true
  try {
    const [cat, cont, proj] = await Promise.all([
      appsApi.catalog().catch(() => ({ catalog: [] })),
      appsApi.listContainers(false).catch(() => ({ containers: [], available: false })),
      appsApi.projects().catch(() => ({ projects: [] })),
    ])
    const c = cat as any
    const cn = cont as any
    const p = proj as any
    dockerAvailable.value = cn.available !== false
    catalog.value = c.catalog || []
    containers.value = cn.containers || []
    projects.value = p.projects || []
  } finally {
    loading.value = false
  }
}

async function loadContainers(silent = false) {
  if (!silent) loading.value = true
  try {
    const res = await appsApi.listContainers(false) as any
    dockerAvailable.value = res.available !== false
    containers.value = res.containers || []
  } finally { loading.value = false }
}

async function loadProjects() {
  loading.value = true
  try {
    const res = await appsApi.projects() as any
    projects.value = res.projects || []
  } finally { loading.value = false }
}

async function loadStats(c: AppContainer) {
  try {
    c.stats = await appsApi.containerStats(c.name) as any
  } catch (e: any) {
    showToast('读取统计失败：' + (e.message || e), 'error')
  }
}

// ===== 容器操作 =====
async function doAction(c: AppContainer, action: 'start' | 'stop' | 'restart' | 'remove') {
  if (action === 'remove' && !confirm(`删除容器 ${c.name}？`)) return
  try {
    await appsApi.containerAction(c.name, action)
    showToast(`${c.name} 已 ${actionLabel(action)}`, 'success')
    c.state = action === 'start' ? 'running' : action === 'stop' ? 'exited' : c.state
    // 重新拉取最新状态
    setTimeout(() => refreshAll(true), 800)
  } catch (e: any) {
    showToast(`${action} 失败：` + (e.message || e), 'error')
  }
}

function actionLabel(a: string) {
  return ({ start: '启动', stop: '停止', restart: '重启', remove: '删除' } as any)[a] || a
}

// ===== 服务级操作（对一组容器）=====
async function serviceAction(svc: ServiceCatalogEntry, action: 'start' | 'stop' | 'restart') {
  busy[`svc-${svc.id}-${action}`] = true
  let ok = 0, fail = 0
  for (const c of svc.containers) {
    // stop 反向；start 正向；restart 全部
    if (action === 'start' && c.state === 'running') continue
    if (action === 'stop' && c.state !== 'running') continue
    try {
      await appsApi.containerAction(c.name, action)
      ok++
    } catch { fail++ }
  }
  showToast(`${svc.name}：${ok} 成功${fail ? `，${fail} 失败` : ''}`, fail ? 'error' : 'success')
  busy[`svc-${svc.id}-${action}`] = false
  setTimeout(() => refreshAll(true), 1000)
}

function startService(svc: ServiceCatalogEntry) { serviceAction(svc, 'start') }
function stopService(svc: ServiceCatalogEntry) { serviceAction(svc, 'stop') }
function restartService(svc: ServiceCatalogEntry) { serviceAction(svc, 'restart') }

// ===== 项目级操作 =====
async function projectAction(p: ComposeProject, action: 'start' | 'stop' | 'restart') {
  busy[`proj-${p.name}-${action}`] = true
  let ok = 0, fail = 0
  // 找到该项目下所有容器
  const projectContainers = containers.value.filter(c => c.composeProject === p.name ||
    (p.name === '(standalone)' && !c.composeProject))
  for (const c of projectContainers) {
    if (action === 'start' && c.state === 'running') continue
    if (action === 'stop' && c.state !== 'running') continue
    try {
      await appsApi.containerAction(c.name, action)
      ok++
    } catch { fail++ }
  }
  showToast(`${p.name}：${ok} 容器已 ${actionLabel(action)}${fail ? `，${fail} 失败` : ''}`, fail ? 'error' : 'success')
  busy[`proj-${p.name}-${action}`] = false
  setTimeout(() => refreshAll(true), 1000)
}
function startProject(p: ComposeProject) { projectAction(p, 'start') }
function stopProject(p: ComposeProject) { projectAction(p, 'stop') }
function restartProject(p: ComposeProject) { projectAction(p, 'restart') }

// ===== 日志 =====
async function openLogs(name: string) {
  logName.value = name
  modals.logs = true
  await reloadLogs()
}
async function reloadLogs() {
  try {
    const res = await appsApi.containerLogs(logName.value, logTail.value) as any
    logContent.value = res.logs || ''
  } catch (e: any) {
    logContent.value = '读取日志失败：' + (e.message || e)
  }
}
function closeLogs() { modals.logs = false }

function viewServiceContainers(svc: ServiceCatalogEntry) {
  selectedService.value = svc
  modals.service = true
}

// ===== UI helpers =====
function openUrl(url: string) { window.open(url, '_blank', 'noopener') }
function formatDate(s?: string) {
  if (!s || s.startsWith('0001-')) return '-'
  return new Date(s).toLocaleString('zh-CN', { hour12: false })
}
function formatBytes(b?: number) {
  if (!b) return '0'
  const u = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(b) / Math.log(1024))
  return (b / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1) + ' ' + u[i]
}
function statusLabel(s: string) {
  return ({ running: '运行中', partial: '部分运行', stopped: '已停止', missing: '未安装', unknown: '未知' } as any)[s] || s
}
function categoryLabel(c: string) {
  return ({
    photo: '照片', 'file-share': '文件共享', dev: '开发', audio: '音频',
    backup: '备份', system: '系统', other: '其他',
  } as any)[c] || c
}
function catalogIcon(hint: string) {
  const m: Record<string, any> = {
    photo: PhotoIcon, share: ServerStackIcon, files: FolderIcon,
    git: CodeBracketIcon, backup: ArchiveBoxIcon, audio: MusicalNoteIcon,
    system: CircleStackIcon, container: CubeIcon,
  }
  return m[hint] || CubeIcon
}
</script>

<style scoped>
.app-center {
  width: 100%; height: 100%;
  background: #f8fafc; color: #0f172a;
  display: flex; flex-direction: column; overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif;
}
.ac-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 16px 24px; background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: white;
}
.ac-header h1 { font-size: 20px; font-weight: 600; margin: 0; }
.ac-header .subtitle { font-size: 12px; opacity: 0.7; margin: 4px 0 0; }
.ac-header-right { display: flex; align-items: center; gap: 12px; }
.ping-badge {
  display: flex; align-items: center; gap: 6px; font-size: 12px;
  padding: 6px 10px; border-radius: 999px; background: rgba(255,255,255,0.08);
}
.ping-badge .dot { width: 8px; height: 8px; border-radius: 50%; background: #94a3b8; }
.ping-badge.ok .dot { background: #22c55e; box-shadow: 0 0 6px #22c55e; }
.ping-badge.err .dot { background: #ef4444; }

.ac-tabs {
  display: flex; gap: 4px; padding: 0 16px;
  background: white; border-bottom: 1px solid #e2e8f0; flex-shrink: 0;
}
.ac-tab {
  display: flex; align-items: center; gap: 6px;
  padding: 12px 16px; background: none; border: none; cursor: pointer;
  font-size: 14px; color: #64748b; border-bottom: 2px solid transparent;
  transition: all .15s;
}
.ac-tab:hover { color: #0f172a; }
.ac-tab.active { color: #0ea5e9; border-bottom-color: #0ea5e9; }
.ac-badge { font-size: 11px; padding: 1px 6px; border-radius: 999px; background: #e0f2fe; color: #0369a1; }

.ac-panel { flex: 1; overflow: auto; padding: 20px 24px; }
.ac-toolbar { display: flex; align-items: center; gap: 12px; margin-bottom: 16px; flex-wrap: wrap; }
.ac-toolbar h2 { font-size: 16px; margin: 0; }
.ac-toolbar .search { max-width: 240px; }
.ac-hint { color: #94a3b8; font-size: 12px; margin-left: auto; }

.ac-btn {
  display: inline-flex; align-items: center; gap: 6px;
  padding: 8px 14px; border-radius: 8px; border: none; cursor: pointer;
  font-size: 13px; font-weight: 500; transition: all .15s; white-space: nowrap;
}
.ac-btn.sm { padding: 5px 10px; font-size: 12px; }
.ac-btn.xs { padding: 3px 8px; font-size: 11px; }
.ac-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.ac-btn.primary { background: #0ea5e9; color: white; }
.ac-btn.primary:hover:not(:disabled) { background: #0284c7; }
.ac-btn.ghost { background: white; color: #334155; border: 1px solid #cbd5e1; }
.ac-btn.ghost:hover:not(:disabled) { background: #f1f5f9; }
.ac-btn.danger { background: white; color: #ef4444; border: 1px solid #fecaca; }

.ac-input {
  padding: 8px 10px; border: 1px solid #cbd5e1; border-radius: 6px;
  font-size: 13px; background: white;
}
.ac-input:focus { outline: none; border-color: #0ea5e9; }

.ac-empty {
  display: flex; flex-direction: column; align-items: center; justify-content: center;
  padding: 60px 20px; color: #64748b; gap: 8px;
}

/* Catalog cards */
.catalog-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(380px, 1fr)); gap: 16px; }
.catalog-card {
  background: white; border-radius: 12px; padding: 16px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05); border: 1px solid #e2e8f0;
  display: flex; flex-direction: column; gap: 10px;
}
.catalog-card.status-running { border-left: 4px solid #22c55e; }
.catalog-card.status-partial { border-left: 4px solid #f59e0b; }
.catalog-card.status-stopped { border-left: 4px solid #ef4444; }
.catalog-card.status-missing { border-left: 4px solid #94a3b8; opacity: 0.85; }
.catalog-head { display: flex; align-items: center; gap: 10px; }
.catalog-icon {
  width: 40px; height: 40px; border-radius: 8px; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center; color: white;
}
.catalog-icon.sm { width: 32px; height: 32px; }
.catalog-icon.icon-photo { background: linear-gradient(135deg, #8b5cf6, #6366f1); }
.catalog-icon.icon-share, .catalog-icon.icon-files { background: linear-gradient(135deg, #0ea5e9, #0284c7); }
.catalog-icon.icon-git { background: linear-gradient(135deg, #f59e0b, #d97706); }
.catalog-icon.icon-backup { background: linear-gradient(135deg, #10b981, #059669); }
.catalog-icon.icon-audio { background: linear-gradient(135deg, #ec4899, #db2777); }
.catalog-icon.icon-system { background: linear-gradient(135deg, #475569, #334155); }
.catalog-icon.icon-container { background: #64748b; }
.catalog-title { flex: 1; min-width: 0; }
.catalog-title h3 { margin: 0; font-size: 15px; }
.catalog-cat { font-size: 11px; color: #94a3b8; margin-top: 2px; }
.catalog-status {
  display: inline-flex; align-items: center; gap: 4px;
  font-size: 11px; padding: 3px 8px; border-radius: 999px; font-weight: 500;
}
.catalog-status .status-dot { width: 6px; height: 6px; border-radius: 50%; }
.catalog-status.st-running { background: #dcfce7; color: #166534; }
.catalog-status.st-running .status-dot { background: #22c55e; }
.catalog-status.st-partial { background: #fef3c7; color: #854d0e; }
.catalog-status.st-partial .status-dot { background: #f59e0b; }
.catalog-status.st-stopped { background: #fee2e2; color: #991b1b; }
.catalog-status.st-stopped .status-dot { background: #ef4444; }
.catalog-status.st-missing { background: #f1f5f9; color: #64748b; }
.catalog-status.st-missing .status-dot { background: #94a3b8; }
.catalog-desc { font-size: 12px; color: #64748b; margin: 0; line-height: 1.5; }
.catalog-meta { display: flex; flex-direction: column; gap: 4px; font-size: 12px; }
.meta-row { display: flex; align-items: center; gap: 6px; color: #475569; }
.meta-row a { color: #0ea5e9; text-decoration: none; }
.meta-row a:hover { text-decoration: underline; }
.meta-row code, .containers-list code, .sub-container code {
  font-family: ui-monospace, monospace; font-size: 11px;
  background: #f1f5f9; padding: 1px 5px; border-radius: 3px;
  word-break: break-all;
}
.catalog-note {
  display: flex; gap: 6px; align-items: flex-start;
  font-size: 11px; color: #64748b; background: #f8fafc;
  padding: 6px 8px; border-radius: 6px; line-height: 1.5;
}
.catalog-actions { display: flex; flex-wrap: wrap; gap: 6px; margin-top: auto; }
.catalog-details { font-size: 12px; }
.catalog-details summary { cursor: pointer; color: #475569; user-select: none; padding: 4px 0; }
.sub-containers, .containers-list { display: flex; flex-direction: column; gap: 4px; margin-top: 6px; }
.sub-container {
  display: flex; align-items: center; gap: 6px;
  font-size: 11px; padding: 4px 6px; background: #f8fafc; border-radius: 4px;
}
.sub-status {
  width: 8px; height: 8px; border-radius: 50%; background: #94a3b8; flex-shrink: 0;
}
.sub-status.running { background: #22c55e; }
.sub-status.exited { background: #94a3b8; }
.sub-status.created { background: #f59e0b; }

/* Container table */
.ac-table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
.ac-table th, .ac-table td { padding: 10px 12px; text-align: left; font-size: 13px; border-bottom: 1px solid #f1f5f9; }
.ac-table th { background: #f8fafc; color: #475569; font-weight: 500; font-size: 12px; }
.ac-table tbody tr:hover { background: #f8fafc; }
.ac-table.compact th, .ac-table.compact td { padding: 6px 10px; font-size: 12px; }
.cell-name { display: flex; align-items: center; gap: 6px; }
.cell-image { max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #64748b; font-size: 12px; }
.port-cell { font-size: 11px; }
.port-url { color: #0ea5e9; font-weight: 500; }
.col-actions { white-space: nowrap; display: flex; gap: 4px; }
.ac-chip {
  display: inline-block; padding: 2px 8px; border-radius: 999px; font-size: 11px;
  background: #e0f2fe; color: #0369a1;
}
.ac-chip.st-task-completed { background: #dcfce7; color: #166534; }
.ac-chip.st-task-idle { background: #f1f5f9; color: #475569; }
.ac-chip.st-task-failed { background: #fee2e2; color: #991b1b; }
.muted { color: #94a3b8; }
.xs { font-size: 11px; }

/* Projects */
.projects-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(380px, 1fr)); gap: 16px; }
.project-card {
  background: white; border-radius: 12px; padding: 16px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05); border: 1px solid #e2e8f0;
  display: flex; flex-direction: column; gap: 10px;
}
.project-head { display: flex; align-items: center; gap: 10px; }

/* Modals */
.ac-overlay {
  position: fixed; inset: 0; background: rgba(15,23,42,0.5); backdrop-filter: blur(4px);
  display: flex; align-items: center; justify-content: center; z-index: 1000;
}
.ac-modal {
  background: white; border-radius: 12px; width: 92%; max-width: 560px; max-height: 90vh;
  display: flex; flex-direction: column; box-shadow: 0 20px 60px rgba(0,0,0,0.3);
}
.ac-modal.lg { max-width: 900px; }
.ac-modal-head { display: flex; justify-content: space-between; align-items: center;
  padding: 16px 20px; border-bottom: 1px solid #e2e8f0; }
.ac-modal-head h3 { margin: 0; font-size: 16px; }
.ac-x { background: none; border: none; cursor: pointer; color: #64748b; padding: 4px; border-radius: 4px; }
.ac-x:hover { background: #f1f5f9; color: #0f172a; }
.ac-modal-body { padding: 20px; overflow: auto; }
.ac-modal-foot { padding: 12px 20px; border-top: 1px solid #e2e8f0; display: flex; justify-content: flex-end; gap: 8px; }
.logs-toolbar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.ac-log {
  background: #0f172a; color: #e2e8f0; padding: 12px; border-radius: 8px;
  font-family: ui-monospace, monospace; font-size: 12px; line-height: 1.5;
  max-height: 60vh; overflow: auto; white-space: pre-wrap; word-break: break-word;
}

.ac-toast {
  position: fixed; bottom: 24px; right: 24px;
  padding: 12px 20px; border-radius: 8px; color: white;
  font-size: 13px; max-width: 360px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.2); z-index: 2000;
}
.ac-toast.toast-success { background: #16a34a; }
.ac-toast.toast-error { background: #dc2626; }
.ac-toast.toast-info { background: #0ea5e9; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }
</style>
