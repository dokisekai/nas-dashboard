<template>
  <div class="docker-manager">
    <div class="docker-header">
      <div class="header-info">
        <h1>🐳 Docker 管理</h1>
        <p class="subtitle">容器化应用管理平台 - 容器、镜像管理</p>
        <div class="docker-stats">
          <div class="stat-item">
            <span class="stat-icon">🟢</span>
            <span class="stat-value">{{ containers.filter(c => c.state === 'running').length }}</span>
            <span class="stat-label">运行中</span>
          </div>
          <div class="stat-item">
            <span class="stat-icon">🔴</span>
            <span class="stat-value">{{ containers.filter(c => c.state !== 'running').length }}</span>
            <span class="stat-label">已停止</span>
          </div>
          <div class="stat-item">
            <span class="stat-icon">📦</span>
            <span class="stat-value">{{ images.length }}</span>
            <span class="stat-label">镜像</span>
          </div>
        </div>
      </div>
      <div class="header-actions">
        <label class="refresh-toggle">
          <input type="checkbox" v-model="autoRefresh" @change="toggleAutoRefresh" />
          <span>自动刷新</span>
        </label>
        <button class="action-btn primary" @click="fetchData" :disabled="loading">
          <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          刷新
        </button>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="['tab-btn', { active: activeTab === tab.id }]"
        @click="activeTab = tab.id"
      >
        <component :is="tab.icon" class="w-5 h-5" />
        {{ tab.label }}
      </button>
    </div>

    <!-- Containers Tab -->
    <div v-if="activeTab === 'containers'" class="tab-content">
      <div class="content-header">
        <h2>容器 ({{ containers.length }})</h2>
        <div class="header-actions">
          <div class="search-box">
            <MagnifyingGlassIcon class="w-4 h-4 search-icon" />
            <input
              v-model="containerSearch"
              type="text"
              placeholder="搜索容器名称或镜像..."
              class="search-input"
            />
          </div>
          <button @click="startAllStopped" class="action-btn info" title="启动所有已停止的容器">
            <PlayIcon class="w-4 h-4" />
            启动全部
          </button>
          <button @click="stopAllRunning" class="action-btn warning" title="停止所有运行中的容器">
            <PauseIcon class="w-4 h-4" />
            停止全部
          </button>
          <select v-model="containerFilter" class="filter-select">
            <option value="all">全部</option>
            <option value="running">运行</option>
            <option value="stopped">停止</option>
          </select>
        </div>
      </div>

      <div class="table-container">
        <table class="docker-table">
          <thead>
            <tr>
              <th width="40">
                <input type="checkbox" v-model="selectAllContainers" @change="toggleSelectAllContainers" />
              </th>
              <th>状态</th>
              <th>名称</th>
              <th>资源</th>
              <th>端口</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="container in filteredContainers"
              :key="container.id"
              :class="{ selected: selectedContainers.includes(container.id) }"
            >
              <td>
                <input
                  type="checkbox"
                  :value="container.id"
                  v-model="selectedContainers"
                />
              </td>
              <td>
                <span class="status-badge" :class="container.state">
                  {{ getShortStatus(container.state) }}
                </span>
              </td>
              <td class="container-name">{{ container.name || container.names?.join(', ').replace(/^\//, '') }}</td>
              <td class="text-gray-500">
                <div v-if="containerStats[container.id]" class="stats-mini">
                  <div class="stat-item">
                    <span class="label">CPU:</span>
                    <span class="value">{{ containerStats[container.id].CPUPerc }}</span>
                  </div>
                  <div class="stat-item">
                    <span class="label">MEM:</span>
                    <span class="value">{{ containerStats[container.id].MemPerc }}</span>
                  </div>
                </div>
                <div v-else-if="container.state === 'running'" class="text-xs text-gray-400">
                  加载中...
                </div>
                <div v-else>-</div>
              </td>
              <td class="text-gray-500">
                <div v-if="container.ports && typeof container.ports === 'string'">{{ formatPortString(container.ports) }}</div>
                <div v-else-if="container.ports">
                  <div v-for="port in container.ports" :key="port.PublicPort">
                    {{ formatPort(port) }}
                  </div>
                </div>
                <div v-else>-</div>
              </td>
              <td class="text-right">
                <div class="action-group">
                  <button
                    v-if="container.state !== 'running'"
                    @click="handleContainerAction(container.id, 'start')"
                    class="icon-btn success"
                    title="启动"
                  >
                    <PlayIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="container.state === 'running'"
                    @click="handleContainerAction(container.id, 'stop')"
                    class="icon-btn warning"
                    title="停止"
                  >
                    <PauseIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="handleContainerAction(container.id, 'restart')"
                    class="icon-btn info"
                    title="重启"
                  >
                    <ArrowPathIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="showLogs(container.id, getContainerName(container))"
                    class="icon-btn"
                    title="日志"
                  >
                    <DocumentTextIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="showTerminal(container.id, getContainerName(container))"
                    class="icon-btn"
                    title="终端"
                  >
                    <CommandLineIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="inspectContainer(container.id, getContainerName(container))"
                    class="icon-btn"
                    title="详情"
                  >
                    <InformationCircleIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="handleContainerAction(container.id, 'remove')"
                    class="icon-btn danger"
                    title="删除"
                  >
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredContainers.length === 0 && !loading">
              <td colspan="6" class="empty-cell">暂无容器</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- 批量操作 -->
      <div v-if="selectedContainers.length > 0" class="batch-actions">
        <div class="batch-info">已选择 {{ selectedContainers.length }} 个容器</div>
        <div class="batch-buttons">
          <button @click="batchStart" class="btn btn-sm">批量启动</button>
          <button @click="batchStop" class="btn btn-sm">批量停止</button>
          <button @click="batchRestart" class="btn btn-sm">批量重启</button>
          <button @click="clearSelection" class="btn btn-sm btn-ghost">取消选择</button>
        </div>
      </div>
    </div>

    <!-- Images Tab -->
    <div v-if="activeTab === 'images'" class="tab-content">
      <div class="content-header">
        <h2>镜像 ({{ images.length }})</h2>
        <div class="header-actions">
          <div class="search-box">
            <MagnifyingGlassIcon class="w-4 h-4 search-icon" />
            <input
              v-model="imageSearch"
              type="text"
              placeholder="搜索镜像名称或标签..."
              class="search-input"
            />
          </div>
          <button class="action-btn primary" @click="pruneImages" title="清理未使用的镜像">
            <TrashIcon class="w-4 h-4" />
            清理镜像
          </button>
          <button class="action-btn secondary" @click="showPullModal = true">
            <CloudArrowDownIcon class="w-4 h-4" />
            拉取镜像
          </button>
        </div>
      </div>

      <div class="table-container">
        <table class="docker-table">
          <thead>
            <tr>
              <th>镜像名称</th>
              <th>标签</th>
              <th>ID</th>
              <th>大小</th>
              <th>创建时间</th>
              <th class="text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="image in filteredImages" :key="image.id">
              <td class="font-medium">{{ getImageRepoName(image) }}</td>
              <td>
                <span v-for="tag in getImageTags(image)" :key="tag" class="tag-badge">
                  {{ tag }}
                </span>
              </td>
              <td class="text-gray-500 font-mono text-xs">{{ (image.id || image.ID).substring(7, 19) }}</td>
              <td class="text-gray-500">{{ formatBytes(image.size || image.Size) }}</td>
              <td class="text-gray-500">{{ formatImageDate(image.created || image.CreatedAt) }}</td>
              <td class="text-right">
                <div class="action-group">
                  <button
                    @click="runContainer(image)"
                    class="icon-btn success"
                    title="运行容器"
                  >
                    <PlayIcon class="w-4 h-4" />
                  </button>
                  <button
                    @click="handleImageAction(image.id || image.ID, 'remove')"
                    class="icon-btn danger"
                    title="删除镜像"
                  >
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="images.length === 0 && !loading">
              <td colspan="6" class="empty-cell">暂无本地镜像</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Networks Tab -->
    <div v-if="activeTab === 'networks'" class="tab-content">
      <div class="content-header">
        <h2>网络 ({{ networks.length }})</h2>
        <button class="action-btn secondary" @click="showCreateNetworkModal = true">
          <PlusIcon class="w-4 h-4" />
          创建网络
        </button>
      </div>

      <div class="table-container">
        <table class="docker-table">
          <thead>
            <tr>
              <th>名称</th>
              <th>ID</th>
              <th>驱动</th>
              <th>子网</th>
              <th>网关</th>
              <th>容器数</th>
              <th class="text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="network in networks" :key="network.id">
              <td class="font-medium">{{ network.name }}</td>
              <td class="text-gray-500 font-mono text-xs">{{ network.id?.substring(0, 12) }}</td>
              <td>
                <span class="driver-badge">{{ network.driver }}</span>
              </td>
              <td class="text-gray-500 text-xs">{{ network.subnet || '-' }}</td>
              <td class="text-gray-500 text-xs">{{ network.gateway || '-' }}</td>
              <td>{{ network.containerCount || 0 }}</td>
              <td class="text-right">
                <div class="action-group">
                  <button
                    @click="inspectNetwork(network.id, network.name)"
                    class="icon-btn"
                    title="详情"
                  >
                    <InformationCircleIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="network.name !== 'bridge' && network.name !== 'host' && network.name !== 'none'"
                    @click="removeNetwork(network.id, network.name)"
                    class="icon-btn danger"
                    title="删除网络"
                  >
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="networks.length === 0 && !loading">
              <td colspan="7" class="empty-cell">暂无网络</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Volumes Tab -->
    <div v-if="activeTab === 'volumes'" class="tab-content">
      <div class="content-header">
        <h2>卷 ({{ volumes.length }})</h2>
        <button class="action-btn secondary" @click="showCreateVolumeModal = true">
          <PlusIcon class="w-4 h-4" />
          创建卷
        </button>
      </div>

      <div class="table-container">
        <table class="docker-table">
          <thead>
            <tr>
              <th>名称</th>
              <th>驱动</th>
              <th>挂载点</th>
              <th>大小</th>
              <th>状态</th>
              <th class="text-right">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="volume in volumes" :key="volume.name">
              <td class="font-medium">{{ volume.name }}</td>
              <td>
                <span class="driver-badge">{{ volume.driver }}</span>
              </td>
              <td class="text-gray-500 text-xs">{{ volume.mountpoint || '-' }}</td>
              <td class="text-gray-500">{{ volume.usage ? formatBytes(volume.usage) : '-' }}</td>
              <td>
                <span :class="['status-badge', volume.status === 'in-use' ? 'active' : 'inactive']">
                  {{ volume.status === 'in-use' ? '使用中' : '未使用' }}
                </span>
              </td>
              <td class="text-right">
                <div class="action-group">
                  <button
                    @click="inspectVolume(volume.name, volume.name)"
                    class="icon-btn"
                    title="详情"
                  >
                    <InformationCircleIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="volume.usage === 0"
                    @click="removeVolume(volume.name, volume.name)"
                    class="icon-btn danger"
                    title="删除卷"
                  >
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="volumes.length === 0 && !loading">
              <td colspan="6" class="empty-cell">暂无卷</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Pull Image Modal -->
    <div v-if="showPullModal" class="modal-overlay" @click.self="showPullModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>拉取镜像</h3>
          <button @click="showPullModal = false" class="btn-close">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>镜像名称</label>
            <input
              v-model="imageToPull"
              type="text"
              placeholder="nginx:latest"
              @keyup.enter="pullImage"
            />
            <div class="form-hint">例如: nginx:alpine, postgres:15, redis:7</div>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showPullModal = false" class="btn btn-secondary">取消</button>
          <button @click="pullImage" class="btn btn-primary" :disabled="pulling">
            {{ pulling ? '拉取中...' : '拉取' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create Network Modal -->
    <div v-if="showCreateNetworkModal" class="modal-overlay" @click.self="showCreateNetworkModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>创建网络</h3>
          <button @click="showCreateNetworkModal = false" class="btn-close">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>网络名称</label>
            <input v-model="networkForm.name" type="text" placeholder="mynetwork" />
          </div>
          <div class="form-group">
            <label>驱动类型</label>
            <select v-model="networkForm.driver">
              <option value="bridge">bridge</option>
              <option value="overlay">overlay</option>
              <option value="macvlan">macvlan</option>
            </select>
          </div>
          <div class="form-group" v-if="networkForm.driver === 'bridge'">
            <label>子网 (可选)</label>
            <input v-model="networkForm.subnet" type="text" placeholder="172.20.0.0/16" />
          </div>
          <div class="form-group" v-if="networkForm.driver === 'bridge'">
            <label>网关 (可选)</label>
            <input v-model="networkForm.gateway" type="text" placeholder="172.20.0.1" />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showCreateNetworkModal = false" class="btn btn-secondary">取消</button>
          <button @click="createNetwork" class="btn btn-primary" :disabled="creatingNetwork">
            {{ creatingNetwork ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create Volume Modal -->
    <div v-if="showCreateVolumeModal" class="modal-overlay" @click.self="showCreateVolumeModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>创建卷</h3>
          <button @click="showCreateVolumeModal = false" class="btn-close">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>卷名称</label>
            <input v-model="volumeForm.name" type="text" placeholder="myvolume" />
          </div>
          <div class="form-group">
            <label>驱动类型</label>
            <select v-model="volumeForm.driver">
              <option value="local">local</option>
            </select>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showCreateVolumeModal = false" class="btn btn-secondary">取消</button>
          <button @click="createVolume" class="btn btn-primary" :disabled="creatingVolume">
            {{ creatingVolume ? '创建中...' : '创建' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Container Details Modal -->
    <!-- 容器详情模态框 -->
    <div v-if="showContainerDetailsModal" class="modal-overlay" @click.self="showContainerDetailsModal = false">
      <div class="modal large">
        <div class="modal-header">
          <h3>ℹ️ 容器详情 - {{ selectedContainer?.name }}</h3>
          <div class="header-actions">
            <button @click="refreshContainerDetails" class="action-btn sm" :disabled="loadingDetails">
              <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loadingDetails }" />
              刷新
            </button>
            <button @click="showContainerDetailsModal = false" class="btn-close">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
        <div class="modal-body">
          <div v-if="loadingDetails" class="loading-state">
            <div class="spinner"></div>
            <p>加载容器详情中...</p>
          </div>
          <div v-else-if="containerDetails" class="container-details">
            <div class="detail-section">
              <h4>📋 基本信息</h4>
              <div class="detail-grid">
                <div class="detail-item">
                  <label>🏷️ 容器ID</label>
                  <span class="detail-value">{{ containerDetails.id?.substring(0, 12) }}</span>
                </div>
                <div class="detail-item">
                  <label>🖼️ 镜像</label>
                  <span class="detail-value">{{ containerDetails.config?.image }}</span>
                </div>
                <div class="detail-item">
                  <label>⚡ 状态</label>
                  <span class="detail-value status-value">{{ containerDetails.state?.status }}</span>
                </div>
                <div class="detail-item">
                  <label>🕐 创建时间</label>
                  <span class="detail-value">{{ formatDate(containerDetails.created) }}</span>
                </div>
                <div class="detail-item">
                  <label>🔄 重启次数</label>
                  <span class="detail-value">{{ containerDetails.restartCount || 0 }}</span>
                </div>
                <div class="detail-item">
                  <label>🏠 工作目录</label>
                  <span class="detail-value">{{ containerDetails.config?.workingDir || '/' }}</span>
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>🌐 网络设置</h4>
              <div class="network-list">
                <div v-for="(settings, networkName) in containerDetails.networkSettings?.networks" :key="networkName" class="network-item">
                  <div class="network-name">🔗 {{ networkName }}</div>
                  <div class="network-details">
                    <span class="network-ip">📡 IP: {{ settings.ipAddress }}</span>
                    <span class="network-mac" v-if="settings.macAddress">💳 MAC: {{ settings.macAddress }}</span>
                  </div>
                </div>
                <div v-if="!containerDetails.networkSettings?.networks || Object.keys(containerDetails.networkSettings?.networks).length === 0" class="empty-state">
                  无网络配置
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>🔌 端口映射</h4>
              <div class="port-list">
                <div v-for="(portConfig, portKey) in containerDetails.networkSettings?.ports" :key="portKey" class="port-item">
                  <div class="port-info">
                    <span class="port-mapping">{{ formatPortMapping(portKey, portConfig) }}</span>
                  </div>
                </div>
                <div v-if="!containerDetails.networkSettings?.ports || Object.keys(containerDetails.networkSettings?.ports).length === 0" class="empty-state">
                  无端口映射
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>💾 挂载卷</h4>
              <div class="mount-list">
                <div v-for="mount in containerDetails.mounts" :key="mount.destination" class="mount-item">
                  <div class="mount-info">
                    <div class="mount-source">📁 {{ mount.source || 'anonymous' }}</div>
                    <div class="mount-destination">→ {{ mount.destination }}</div>
                    <div class="mount-type">类型: {{ mount.type }}</div>
                    <div class="mount-mode" v-if="mount.mode">模式: {{ mount.mode }}</div>
                  </div>
                </div>
                <div v-if="!containerDetails.mounts || containerDetails.mounts.length === 0" class="empty-state">
                  无挂载卷
                </div>
              </div>
            </div>

            <div class="detail-section">
              <h4>⚙️ 环境变量</h4>
              <div class="env-list">
                <div v-for="(value, key) in containerDetails.config?.env" :key="key" class="env-item">
                  <div class="env-key">{{ key }}</div>
                  <div class="env-value">{{ value }}</div>
                </div>
                <div v-if="!containerDetails.config?.env || Object.keys(containerDetails.config?.env).length === 0" class="empty-state">
                  无环境变量
                </div>
              </div>
            </div>
          </div>
          <div v-else class="error-state">
            <p>无法加载容器详情</p>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showContainerDetailsModal = false" class="btn btn-secondary">关闭</button>
        </div>
      </div>
    </div>
    </div>

    <!-- Terminal Modal -->
    <!-- 终端模态框 -->
    <div v-if="showTerminalModal" class="modal-overlay terminal-modal" @click.self="closeTerminal">
      <div class="modal large terminal-modal-content">
        <div class="modal-header">
          <h3>💻 终端 - {{ selectedContainer?.name }}</h3>
          <div class="header-actions">
            <div class="terminal-status">
              <span class="status-dot" :class="{ connected: terminalConnected, disconnected: !terminalConnected }"></span>
              <span>{{ terminalConnected ? '已连接' : '未连接' }}</span>
            </div>
            <button @click="executeQuickCommand('top')" class="action-btn sm">top</button>
            <button @click="executeQuickCommand('ps aux')" class="action-btn sm">ps</button>
            <button @click="executeQuickCommand('env')" class="action-btn sm">env</button>
            <button @click="closeTerminal" class="btn-close">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
        <div class="modal-body terminal-body">
          <div v-if="!terminalConnected" class="terminal-disconnected">
            <div class="terminal-info">
              <h4>🔌 终端连接说明</h4>
              <p>此功能需要后端WebSocket支持来提供真正的终端访问。</p>
              <div class="terminal-features">
                <div class="feature-item">
                  <span class="feature-icon">⚡</span>
                  <div>
                    <strong>实时命令执行</strong>
                    <div>直接在容器中执行命令</div>
                  </div>
                </div>
                <div class="feature-item">
                  <span class="feature-icon">📊</span>
                  <div>
                    <strong>进程监控</strong>
                    <div>查看容器内运行的进程</div>
                  </div>
                </div>
                <div class="feature-item">
                  <span class="feature-icon">🔧</span>
                  <div>
                    <strong>环境调试</strong>
                    <div>检查环境变量和配置</div>
                  </div>
                </div>
              </div>
              <div class="terminal-commands">
                <h5>常用命令：</h5>
                <div class="command-suggestions">
                  <button @click="executeQuickCommand('top')" class="cmd-btn">top - 查看进程</button>
                  <button @click="executeQuickCommand('ps aux')" class="cmd-btn">ps aux - 列出进程</button>
                  <button @click="executeQuickCommand('env')" class="cmd-btn">env - 环境变量</button>
                  <button @click="executeQuickCommand('df -h')" class="cmd-btn">df -h - 磁盘使用</button>
                  <button @click="executeQuickCommand('free -h')" class="cmd-btn">free -h - 内存使用</button>
                </div>
              </div>
            </div>
          </div>
          <div v-else ref="terminalRef" class="terminal-container">
            <div class="terminal-line connected-message">✅ 终端已连接到容器 {{ selectedContainer?.name }}</div>
            <div class="terminal-line">💡 输入命令按Enter执行，输入 'exit' 退出</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Logs Modal -->
    <!-- 日志模态框 -->
    <div v-if="showLogsModal" class="modal-overlay" @click.self="closeLogs">
      <div class="modal large">
        <div class="modal-header">
          <h3>📄 容器日志 - {{ selectedContainer?.name }}</h3>
          <div class="header-actions">
            <button @click="refreshLogs" class="action-btn sm" :disabled="loadingLogs">
              <ArrowPathIcon class="w-4 h-4" :class="{ 'animate-spin': loadingLogs }" />
              刷新
            </button>
            <button @click="clearLogsDisplay" class="action-btn sm">
              清空
            </button>
            <button @click="copyLogs" class="action-btn sm">
              复制
            </button>
            <button @click="closeLogs" class="btn-close">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
        </div>
        <div class="modal-body">
          <div class="logs-controls">
            <div class="logs-stats">
              <span class="log-stat">📊 行数: {{ logLines.length }}</span>
              <span class="log-stat">📦 大小: {{ formatBytes(selectedContainerLogs.length) }}</span>
            </div>
            <div class="logs-filters">
              <select v-model="logFilter" class="filter-select sm">
                <option value="all">全部日志</option>
                <option value="error">仅错误</option>
                <option value="warn">警告和错误</option>
                <option value="info">信息及以上</option>
              </select>
              <label class="auto-refresh-toggle">
                <input type="checkbox" v-model="autoRefreshLogs" />
                <span>自动刷新</span>
              </label>
            </div>
          </div>
          <div class="logs-container">
            <div v-if="loadingLogs" class="logs-loading">加载日志中...</div>
            <div v-else-if="filteredLogs.length === 0" class="logs-empty">暂无日志内容</div>
            <pre v-else class="logs-content">{{ filteredLogs }}</pre>
          </div>
        </div>
      </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import {
  CubeIcon,
  CircleStackIcon,
  ShareIcon,
  FolderIcon,
  PlayIcon,
  PauseIcon,
  TrashIcon,
  ArrowPathIcon,
  CloudArrowDownIcon,
  PlusIcon,
  XMarkIcon,
  DocumentTextIcon,
  InformationCircleIcon,
  CommandLineIcon,
  MagnifyingGlassIcon
} from '@heroicons/vue/24/outline'
import { useNotificationStore } from '../stores/notification'

// API函数
// 获取认证头的辅助函数
const getAuthHeaders = (): Record<string, string> => {
  const token = localStorage.getItem('token')
  const headers: Record<string, string> = {}
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  return headers
}

const serviceApi = {
  startContainer: (id: string) => fetch(`/api/docker/containers/${id}/start`, {
    method: 'POST',
    headers: getAuthHeaders()
  }),
  stopContainer: (id: string) => fetch(`/api/docker/containers/${id}/stop`, {
    method: 'POST',
    headers: getAuthHeaders()
  }),
  restartContainer: (id: string) => fetch(`/api/docker/containers/${id}/restart`, {
    method: 'POST',
    headers: getAuthHeaders()
  }),
  removeContainer: (id: string) => fetch(`/api/docker/containers/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders()
  }),
  getContainers: () => fetch('/api/docker/containers', { headers: getAuthHeaders() }),
  getContainerLogs: (id: string) => fetch(`/api/docker/containers/${id}/logs`, {
    headers: getAuthHeaders()
  }),
  execInContainer: (id: string, cmd: string[]) => fetch(`/api/docker/containers/${id}/exec`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders()
    },
    body: JSON.stringify({ cmd })
  }),
  getImages: () => fetch('/api/docker/images', { headers: getAuthHeaders() }),
  removeImage: (id: string) => fetch(`/api/docker/images/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders()
  }),
  pullImage: (image: string) => fetch('/api/docker/images/pull', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders()
    },
    body: JSON.stringify({ image })
  }),
  getDockerNetworks: () => fetch('/api/docker/networks', { headers: getAuthHeaders() }),
  getDockerVolumes: () => fetch('/api/docker/volumes', { headers: getAuthHeaders() }),
  removeNetwork: (id: string) => fetch(`/api/docker/networks/${id}`, {
    method: 'DELETE',
    headers: getAuthHeaders()
  }),
  removeVolume: (name: string) => fetch(`/api/docker/volumes/${name}`, {
    method: 'DELETE',
    headers: getAuthHeaders()
  }),
  createNetwork: (config: any) => fetch('/api/docker/networks', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders()
    },
    body: JSON.stringify(config)
  }),
  createVolume: (config: any) => fetch('/api/docker/volumes', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...getAuthHeaders()
    },
    body: JSON.stringify(config)
  }),
  inspectContainer: (id: string) => fetch(`/api/docker/containers/${id}/json`, {
    headers: getAuthHeaders()
  }),
  inspectNetwork: (id: string) => fetch(`/api/docker/networks/${id}`, {
    headers: getAuthHeaders()
  }),
  inspectVolume: (name: string) => fetch(`/api/docker/volumes/${name}`, {
    headers: getAuthHeaders()
  })
}

const tabs = [
  { id: 'containers', label: '容器', icon: CubeIcon },
  { id: 'images', label: '镜像', icon: CircleStackIcon }
  // 网络和卷功能暂时禁用，等待后端API实现
  // { id: 'networks', label: '网络', icon: ShareIcon },
  // { id: 'volumes', label: '卷', icon: FolderIcon }
]

const activeTab = ref('containers')
const loading = ref(false)
const containers = ref<any[]>([])
const images = ref<any[]>([])
const networks = ref<any[]>([])
const volumes = ref<any[]>([])
const containerStats = ref<Record<string, any>>({})
const autoRefresh = ref(false)
const notificationStore = useNotificationStore()
let refreshInterval: any = null

// 过滤和选择
const containerFilter = ref('all')
const selectedContainers = ref<string[]>([])
const selectAllContainers = ref(false)

// 搜索和过滤
const containerSearch = ref('')
const filteredContainers = computed(() => {
  let result = containers.value

  // 应用过滤器
  if (containerFilter.value === 'running') {
    result = result.filter(c => c.state === 'running')
  } else if (containerFilter.value === 'stopped') {
    result = result.filter(c => c.state !== 'running')
  }

  // 应用搜索
  if (containerSearch.value) {
    const searchLower = containerSearch.value.toLowerCase()
    result = result.filter(c => {
      const name = (c.name || c.names?.join(', ') || '').toLowerCase()
      const image = (c.image || '').toLowerCase()
      return name.includes(searchLower) || image.includes(searchLower)
    })
  }

  return result
})

// 镜像搜索和过滤
const imageSearch = ref('')
const filteredImages = computed(() => {
  if (!imageSearch.value) return images.value

  const searchLower = imageSearch.value.toLowerCase()
  return images.value.filter(image => {
    const repoName = getImageRepoName(image).toLowerCase()
    const tags = getImageTags(image).map(t => t.toLowerCase()).join(' ')
    return repoName.includes(searchLower) || tags.includes(searchLower)
  })
})

// 模态框状态
const showPullModal = ref(false)
const showCreateNetworkModal = ref(false)
const showCreateVolumeModal = ref(false)
const showContainerDetailsModal = ref(false)
const showLogsModal = ref(false)
const showTerminalModal = ref(false)

// 表单数据
const imageToPull = ref('')
const pulling = ref(false)
const networkForm = ref({
  name: '',
  driver: 'bridge',
  subnet: '',
  gateway: ''
})
const volumeForm = ref({
  name: '',
  driver: 'local'
})
const creatingNetwork = ref(false)
const creatingVolume = ref(false)

// 选中的容器和相关数据
const selectedContainer = ref<any>(null)
const containerDetails = ref<any>(null)
const selectedContainerLogs = ref('')
const terminalRef = ref<any>(null)

// 日志功能增强
const loadingLogs = ref(false)
const logFilter = ref('all')
const autoRefreshLogs = ref(false)
let logRefreshInterval: any = null

// 终端功能增强
const terminalConnected = ref(false)
const terminalCommand = ref('')
const terminalHistory = ref<string[]>([])
const terminalHistoryIndex = ref(-1)

// 容器详情功能
const loadingDetails = ref(false)

// 格式化端口映射
const formatPortMapping = (portKey: string, portConfig: any) => {
  // portKey格式可能是 "80/tcp" 或 "80/tcp/udp"
  const parts = portKey.split('/')
  const port = parts[0]
  const protocol = parts[1] || 'tcp'

  if (portConfig && portConfig.length > 0) {
    const hostPort = portConfig[0].hostPort
    return `${hostPort}:${port}/${protocol}`
  }

  return `${port}/${protocol} (未映射)`
}

// 计算过滤后的日志
const filteredLogs = computed(() => {
  if (!selectedContainerLogs.value) return ''

  const lines = selectedContainerLogs.value.split('\n')
  if (logFilter.value === 'all') return lines.join('\n')

  return lines.filter(line => {
    const lowerLine = line.toLowerCase()
    switch (logFilter.value) {
      case 'error':
        return lowerLine.includes('error') || lowerLine.includes('err')
      case 'warn':
        return lowerLine.includes('error') || lowerLine.includes('warn')
      case 'info':
        return lowerLine.includes('info') || lowerLine.includes('error') || lowerLine.includes('warn')
      default:
        return true
    }
  }).join('\n')
})

// 日志行数统计
const logLines = computed(() => {
  if (!selectedContainerLogs.value) return []
  return selectedContainerLogs.value.split('\n').filter(line => line.trim())
})

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const [containersRes, imagesRes] = await Promise.all([
      serviceApi.getContainers(),
      serviceApi.getImages()
    ])

    if (containersRes.ok) {
      const result = await containersRes.json()
      containers.value = result.containers || result || []

      // 获取运行中容器的统计信息
      for (const container of containers.value) {
        if (container.state === 'running') {
          fetchContainerStats(container.id)
        }
      }
    }

    if (imagesRes.ok) {
      const result = await imagesRes.json()
      images.value = result.images || result || []
    }
  } catch (error) {
    console.error('Failed to fetch Docker data:', error)
  } finally {
    loading.value = false
  }
}

const fetchContainerStats = async (containerId: string) => {
  try {
    // 这里应该调用获取统计信息的API
    // 暂时使用模拟数据
    containerStats.value[containerId] = {
      CPUPerc: '0.5%',
      MemPerc: '15.2%'
    }
  } catch (error) {
    console.error('Failed to fetch container stats:', error)
  }
}

// 容器操作
const handleContainerAction = async (id: string, action: 'start' | 'stop' | 'restart' | 'remove') => {
  if (action === 'remove' && !confirm('确定要删除此容器吗？')) return

  try {
    loading.value = true
    switch (action) {
      case 'start': await serviceApi.startContainer(id); break
      case 'stop': await serviceApi.stopContainer(id); break
      case 'restart': await serviceApi.restartContainer(id); break
      case 'remove': await serviceApi.removeContainer(id); break
    }

    notificationStore.add({
      type: 'success',
      title: '操作成功',
      message: `容器已${action === 'start' ? '启动' : action === 'stop' ? '停止' : action === 'restart' ? '重启' : '删除'}`
    })

    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '操作失败',
      message: error.message || '操作失败'
    })
  } finally {
    loading.value = false
  }
}

// 镜像操作
// 清理未使用的镜像
const pruneImages = async () => {
  const unusedImages = images.value.filter(image => {
    const tags = getImageTags(image)
    return tags.length === 0 || tags[0] === '<none>'
  })

  if (unusedImages.length === 0) {
    notificationStore.add({
      type: 'info',
      title: '无需清理',
      message: '没有发现未使用的镜像'
    })
    return
  }

  if (!confirm(`发现 ${unusedImages.length} 个未使用的镜像，确定要清理吗？`)) return

  try {
    loading.value = true
    const errors = []

    for (const image of unusedImages) {
      try {
        await serviceApi.removeImage(image.id || image.ID)
      } catch (error: any) {
        errors.push(`${getImageRepoName(image)}: ${error.message}`)
      }
    }

    if (errors.length > 0) {
      notificationStore.add({
        type: 'warning',
        title: '部分清理失败',
        message: `成功清理 ${unusedImages.length - errors.length} 个镜像，${errors.length} 个失败`
      })
    } else {
      notificationStore.add({
        type: 'success',
        title: '清理完成',
        message: `成功清理 ${unusedImages.length} 个未使用的镜像`
      })
    }

    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '清理失败',
      message: error.message || '清理过程中出现错误'
    })
  } finally {
    loading.value = false
  }
}

const handleImageAction = async (id: string, action: 'remove') => {
  if (action === 'remove' && !confirm('确定要删除此镜像吗？')) return

  try {
    loading.value = true
    await serviceApi.removeImage(id)

    notificationStore.add({
      type: 'success',
      title: '操作成功',
      message: '镜像已删除'
    })

    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '操作失败',
      message: error.message || '删除失败'
    })
  } finally {
    loading.value = false
  }
}

const pullImage = async () => {
  if (!imageToPull.value.trim()) return

  pulling.value = true
  try {
    await serviceApi.pullImage(imageToPull.value.trim())

    notificationStore.add({
      type: 'success',
      title: '拉取成功',
      message: `镜像 ${imageToPull.value} 已拉取`
    })

    showPullModal.value = false
    imageToPull.value = ''
    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '拉取失败',
      message: error.message || '拉取失败'
    })
  } finally {
    pulling.value = false
  }
}

const runContainer = async (image: any) => {
  const imageName = getRepoName(image.repoTags || [image.Repository + ":" + image.Tag])
  const containerName = prompt(`为 ${imageName} 输入容器名称:`, imageName.split('/')[1] || imageName)

  if (!containerName) return

  try {
    // 这里应该调用创建容器的API
    notificationStore.add({
      type: 'success',
      title: '容器创建中',
      message: `正在从 ${imageName} 创建容器 ${containerName}`
    })
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '创建失败',
      message: error.message || '创建容器失败'
    })
  }
}

// 网络操作
const createNetwork = async () => {
  if (!networkForm.value.name.trim()) {
    notificationStore.add({
      type: 'error',
      title: '验证失败',
      message: '请输入网络名称'
    })
    return
  }

  creatingNetwork.value = true
  try {
    await serviceApi.createNetwork(networkForm.value)

    notificationStore.add({
      type: 'success',
      title: '创建成功',
      message: `网络 ${networkForm.value.name} 已创建`
    })

    showCreateNetworkModal.value = false
    networkForm.value = { name: '', driver: 'bridge', subnet: '', gateway: '' }
    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '创建失败',
      message: error.message || '创建网络失败'
    })
  } finally {
    creatingNetwork.value = false
  }
}

const removeNetwork = async (id: string, name: string) => {
  if (!confirm(`确定要删除网络 "${name}" 吗？`)) return

  try {
    await serviceApi.removeNetwork(id)

    notificationStore.add({
      type: 'success',
      title: '删除成功',
      message: `网络 ${name} 已删除`
    })

    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '删除失败',
      message: error.message || '删除网络失败'
    })
  }
}

// 卷操作
const createVolume = async () => {
  if (!volumeForm.value.name.trim()) {
    notificationStore.add({
      type: 'error',
      title: '验证失败',
      message: '请输入卷名称'
    })
    return
  }

  creatingVolume.value = true
  try {
    await serviceApi.createVolume(volumeForm.value)

    notificationStore.add({
      type: 'success',
      title: '创建成功',
      message: `卷 ${volumeForm.value.name} 已创建`
    })

    showCreateVolumeModal.value = false
    volumeForm.value = { name: '', driver: 'local' }
    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '创建失败',
      message: error.message || '创建卷失败'
    })
  } finally {
    creatingVolume.value = false
  }
}

const removeVolume = async (name: string, displayName: string) => {
  if (!confirm(`确定要删除卷 "${displayName}" 吗？`)) return

  try {
    await serviceApi.removeVolume(name)

    notificationStore.add({
      type: 'success',
      title: '删除成功',
      message: `卷 ${displayName} 已删除`
    })

    await fetchData()
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '删除失败',
      message: error.message || '删除卷失败'
    })
  }
}

// 详情查看
const inspectContainer = async (id: string, name: string) => {
  try {
    loadingDetails.value = true
    const response = await serviceApi.inspectContainer(id)

    if (response.ok) {
      containerDetails.value = await response.json()
      selectedContainer.value = { id, name }
      showContainerDetailsModal.value = true
    } else {
      notificationStore.add({
        type: 'error',
        title: '获取详情失败',
        message: `HTTP ${response.status}: 无法获取容器详情`
      })
    }
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '获取详情失败',
      message: error.message || '获取容器详情失败'
    })
  } finally {
    loadingDetails.value = false
  }
}

// 刷新容器详情
const refreshContainerDetails = async () => {
  if (!selectedContainer.value?.id) return
  await inspectContainer(selectedContainer.value.id, selectedContainer.value.name)
}

const inspectNetwork = async (id: string, name: string) => {
  try {
    const response = await serviceApi.inspectNetwork(id)
    if (response.ok) {
      const details = await response.json()
      // 显示网络详情
      alert(`网络详情:\n${JSON.stringify(details, null, 2)}`)
    }
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '获取详情失败',
      message: error.message || '获取网络详情失败'
    })
  }
}

const inspectVolume = async (name: string, displayName: string) => {
  try {
    const response = await serviceApi.inspectVolume(name)
    if (response.ok) {
      const details = await response.json()
      // 显示卷详情
      alert(`卷详情:\n${JSON.stringify(details, null, 2)}`)
    }
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '获取详情失败',
      message: error.message || '获取卷详情失败'
    })
  }
}

// 日志和终端
const showLogs = async (id: string, name: string) => {
  try {
    loadingLogs.value = true
    const response = await serviceApi.getContainerLogs(id)
    if (response.ok) {
      const logs = await response.text()
      selectedContainerLogs.value = logs
      selectedContainer.value = { id, name }
      showLogsModal.value = true

      // 启动自动刷新
      if (autoRefreshLogs.value) {
        startLogRefresh(id)
      }
    } else {
      notificationStore.add({
        type: 'error',
        title: '获取日志失败',
        message: `HTTP ${response.status}: 获取容器日志失败`
      })
    }
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '获取日志失败',
      message: error.message || '获取容器日志失败'
    })
  } finally {
    loadingLogs.value = false
  }
}

// 刷新日志
const refreshLogs = async () => {
  if (!selectedContainer.value?.id) return

  try {
    loadingLogs.value = true
    const response = await serviceApi.getContainerLogs(selectedContainer.value.id)
    if (response.ok) {
      const logs = await response.text()
      selectedContainerLogs.value = logs
      notificationStore.add({
        type: 'success',
        title: '日志已刷新',
        message: '容器日志已更新'
      })
    }
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '刷新失败',
      message: error.message || '刷新日志失败'
    })
  } finally {
    loadingLogs.value = false
  }
}

// 清空日志显示
const clearLogsDisplay = () => {
  selectedContainerLogs.value = ''
  notificationStore.add({
    type: 'info',
    title: '日志已清空',
    message: '日志显示已清空'
  })
}

// 复制日志
const copyLogs = async () => {
  try {
    await navigator.clipboard.writeText(selectedContainerLogs.value)
    notificationStore.add({
      type: 'success',
      title: '复制成功',
      message: '日志内容已复制到剪贴板'
    })
  } catch (error) {
    notificationStore.add({
      type: 'error',
      title: '复制失败',
      message: '无法复制日志内容'
    })
  }
}

// 启动日志自动刷新
const startLogRefresh = (containerId: string) => {
  stopLogRefresh() // 清除之前的定时器

  logRefreshInterval = setInterval(async () => {
    try {
      const response = await serviceApi.getContainerLogs(containerId)
      if (response.ok) {
        const logs = await response.text()
        selectedContainerLogs.value = logs
      }
    } catch (error) {
      console.error('自动刷新日志失败:', error)
    }
  }, 5000) // 每5秒刷新一次
}

// 停止日志自动刷新
const stopLogRefresh = () => {
  if (logRefreshInterval) {
    clearInterval(logRefreshInterval)
    logRefreshInterval = null
  }
}

const showTerminal = async (id: string, name: string) => {
  selectedContainer.value = { id, name }
  showTerminalModal.value = true

  await nextTick()

  // 模拟终端连接状态
  terminalConnected.value = false

  // 这里应该建立WebSocket连接来提供真正的终端功能
  // 目前显示连接说明和快速命令
  if (terminalRef.value) {
    terminalRef.value.innerHTML = `
      <div class="terminal-line">🔌 正在连接到容器 ${name}...</div>
      <div class="terminal-line">⚠️ 终端功能需要后端WebSocket支持</div>
      <div class="terminal-line">💡 请使用下方快速命令查看容器信息</div>
    `
  }

  // 模拟连接延迟
  setTimeout(() => {
    terminalConnected.value = false // 保持未连接状态，等待真实WebSocket实现
  }, 1000)
}

// 执行快速命令
const executeQuickCommand = async (command: string) => {
  if (!selectedContainer.value?.id) return

  try {
    // 显示命令执行中
    if (terminalRef.value) {
      const commandLine = document.createElement('div')
      commandLine.className = 'terminal-line command-line'
      commandLine.innerHTML = `<span class="prompt">$</span> ${command}`
      terminalRef.value.appendChild(commandLine)
    }

    // 调用后端API执行命令
    const response = await serviceApi.execInContainer(selectedContainer.value.id, command.split(' '))

    let output = ''
    if (response.ok) {
      const result = await response.json()
      output = result.output || result.data || JSON.stringify(result)
    } else {
      output = `命令执行失败: HTTP ${response.status}`
    }

    // 显示命令输出
    if (terminalRef.value) {
      const outputLine = document.createElement('div')
      outputLine.className = 'terminal-line output-line'
      outputLine.textContent = output || '(命令执行成功，无输出)'
      terminalRef.value.appendChild(outputLine)

      // 添加到历史记录
      if (terminalHistory.value[terminalHistory.value.length - 1] !== command) {
        terminalHistory.value.push(command)
        terminalHistoryIndex.value = terminalHistory.value.length
      }
    }

    notificationStore.add({
      type: 'success',
      title: '命令已执行',
      message: `命令 "${command}" 执行完成`
    })
  } catch (error: any) {
    notificationStore.add({
      type: 'error',
      title: '命令执行失败',
      message: error.message || '执行命令时发生错误'
    })

    // 显示错误信息
    if (terminalRef.value) {
      const errorLine = document.createElement('div')
      errorLine.className = 'terminal-line error-line'
      errorLine.textContent = `❌ 错误: ${error.message}`
      terminalRef.value.appendChild(errorLine)
    }
  }
}

const closeTerminal = () => {
  showTerminalModal.value = false
  selectedContainer.value = null
}

const closeLogs = () => {
  showLogsModal.value = false
  selectedContainerLogs.value = ''
  selectedContainer.value = null
  stopLogRefresh() // 停止自动刷新
}

// 批量操作
const toggleSelectAllContainers = () => {
  if (selectAllContainers.value) {
    selectedContainers.value = filteredContainers.value.map(c => c.id)
  } else {
    selectedContainers.value = []
  }
}

const clearSelection = () => {
  selectedContainers.value = []
  selectAllContainers.value = false
}

const batchStart = async () => {
  for (const id of selectedContainers.value) {
    await handleContainerAction(id, 'start')
  }
  clearSelection()
}

const batchStop = async () => {
  for (const id of selectedContainers.value) {
    await handleContainerAction(id, 'stop')
  }
  clearSelection()
}

const batchRestart = async () => {
  for (const id of selectedContainers.value) {
    await handleContainerAction(id, 'restart')
  }
  clearSelection()
}

// 启动所有已停止的容器
const startAllStopped = async () => {
  const stoppedContainers = containers.value.filter(c => c.state !== 'running')
  if (stoppedContainers.length === 0) {
    notificationStore.add({
      type: 'info',
      title: '无需操作',
      message: '没有已停止的容器'
    })
    return
  }

  if (!confirm(`确定要启动 ${stoppedContainers.length} 个已停止的容器吗？`)) return

  let successCount = 0
  for (const container of stoppedContainers) {
    try {
      await serviceApi.startContainer(container.id)
      successCount++
    } catch (error) {
      console.error(`启动容器 ${container.name} 失败:`, error)
    }
  }

  notificationStore.add({
    type: 'success',
    title: '批量启动完成',
    message: `成功启动 ${successCount}/${stoppedContainers.length} 个容器`
  })

  await fetchData()
}

// 停止所有运行中的容器
const stopAllRunning = async () => {
  const runningContainers = containers.value.filter(c => c.state === 'running')
  if (runningContainers.length === 0) {
    notificationStore.add({
      type: 'info',
      title: '无需操作',
      message: '没有运行中的容器'
    })
    return
  }

  if (!confirm(`确定要停止 ${runningContainers.length} 个运行中的容器吗？`)) return

  let successCount = 0
  for (const container of runningContainers) {
    try {
      await serviceApi.stopContainer(container.id)
      successCount++
    } catch (error) {
      console.error(`停止容器 ${container.name} 失败:`, error)
    }
  }

  notificationStore.add({
    type: 'success',
    title: '批量停止完成',
    message: `成功停止 ${successCount}/${runningContainers.length} 个容器`
  })

  await fetchData()
}

// 自动刷新
const toggleAutoRefresh = () => {
  if (autoRefresh.value) {
    refreshInterval = setInterval(fetchData, 5000)
  } else {
    if (refreshInterval) {
      clearInterval(refreshInterval)
      refreshInterval = null
    }
  }
}

// 工具函数
const getContainerName = (container: any) => {
  return container.name || container.names?.join(', ').replace(/^\//, '') || 'unknown'
}

const getImageRepoName = (image: any) => {
  if (image.repoTags && image.repoTags.length > 0) {
    return getRepoName(image.repoTags)
  }
  // 处理不同的API返回格式
  if (image.Repository && image.Tag) {
    return image.Repository
  }
  if (image.RepoTag) {
    return image.RepoTag.split(':')[0]
  }
  return '<none>'
}

const getImageTags = (image: any) => {
  if (image.repoTags && image.repoTags.length > 0) {
    return getTags(image.repoTags)
  }
  // 处理不同的API返回格式
  if (image.Repository && image.Tag) {
    return [image.Tag]
  }
  if (image.RepoTag) {
    return [image.RepoTag.split(':')[1] || 'latest']
  }
  return ['<none>']
}

const formatDate = (timestamp: number) => {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleString('zh-CN')
}

// 简化状态显示
const getShortStatus = (state: string) => {
  const statusMap: Record<string, string> = {
    'running': '运行',
    'exited': '停止',
    'stopped': '停止',
    'paused': '暂停',
    'restarting': '重启中',
    'dead': '异常',
    'created': '已创建'
  }
  return statusMap[state] || state
}

// 简化端口显示
const formatPort = (port: any) => {
  if (!port.PublicPort) return '-'
  const pubPort = port.PublicPort
  const privPort = port.PrivatePort || '-'
  return `${pubPort}:${privPort}`
}

// 简化端口字符串显示
const formatPortString = (ports: string) => {
  // 简化类似 "0.0.0.0:8080->80/tcp, 0.0.0.0:8443->443/tcp" 的格式
  if (!ports || ports === '-') return '-'
  return ports
    .split(',')
    .map((port: string) => {
      const match = port.match(/(\d+)->(\d+)/)
      if (match) {
        return `${match[1]}:${match[2]}`
      }
      return port
    })
    .join(', ')
    .substring(0, 50) // 限制长度
}

const formatImageDate = (dateString: string) => {
  if (!dateString) return '-'
  // 如果已经是格式化的日期字符串，直接返回
  if (typeof dateString === 'string' && dateString.includes('-')) {
    return dateString
  }
  // 如果是时间戳，则格式化
  if (typeof dateString === 'number') {
    return new Date(dateString * 1000).toLocaleString('zh-CN')
  }
  return '-'
}

const formatBytes = (bytes: number) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const getRepoName = (repoTags: string[]) => {
  if (!repoTags || repoTags.length === 0) return '<none>'
  return repoTags[0].split(':')[0]
}

const getTags = (repoTags: string[]) => {
  if (!repoTags || repoTags.length === 0) return ['<none>']
  return repoTags.map(rt => rt.split(':')[1])
}

// 生命周期
onMounted(() => {
  fetchData()
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.docker-manager {
  padding: 24px;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #f8fafc;
  color: #1e293b;
}

.docker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.header-info h1 {
  font-size: 24px;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

.subtitle {
  color: #475569;
  margin: 4px 0 0 0;
  font-size: 14px;
  font-weight: 400;
}

.header-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.refresh-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
  color: #374151;
}

.action-btn {
  padding: 8px 16px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  font-size: 14px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #1e293b;
  font-weight: 500;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
}

.action-btn.primary {
  background: #2563eb;
  color: white;
  border-color: #2563eb;
  font-weight: 600;
}

.action-btn.primary:hover {
  background: #1d4ed8;
  border-color: #1e40af;
}

.action-btn.secondary {
  background: #475569;
  color: white;
  border-color: #475569;
  font-weight: 600;
}

.action-btn.secondary:hover {
  background: #334155;
  border-color: #334155;
}

.tabs {
  display: flex;
  gap: 4px;
  margin-bottom: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.tab-btn {
  padding: 10px 16px;
  border: none;
  background: none;
  border-bottom: 2px solid transparent;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #6b7280;
  transition: all 0.2s;
}

.tab-btn:hover {
  color: #374151;
}

.tab-btn.active {
  color: #3b82f6;
  border-bottom-color: #3b82f6;
}

.tab-content {
  flex: 1;
  overflow: auto;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.content-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #111827;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
  color: #1e293b;
  background: white;
  font-weight: 500;
  cursor: pointer;
}

.filter-select:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.1);
}

.table-container {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  overflow: hidden;
}

.docker-table {
  width: 100%;
  border-collapse: collapse;
}

.docker-table th,
.docker-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.docker-table th {
  background: #f1f5f9;
  font-weight: 700;
  color: #1e293b;
  font-size: 14px;
  border-bottom: 2px solid #cbd5e1;
}

.docker-table tr:hover {
  background: #f8fafc;
}

.docker-table tr.selected {
  background: #dbeafe;
  border-left: 3px solid #2563eb;
}

.docker-table .empty-cell {
  text-align: center;
  color: #64748b;
  padding: 32px;
  font-style: italic;
  font-weight: 500;
}

/* 容器名称专用样式 - 更清晰显示 */
.container-name {
  font-weight: 700;
  font-size: 14px;
  color: #0f172a;
  letter-spacing: 0.01em;
  line-height: 1.5;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
  display: inline-block;
  white-space: nowrap;
}

.status-badge.running {
  background: #dcfce7;
  color: #166534;
  font-weight: 600;
  border: 1px solid #86efac;
}

.status-badge.exited,
.status-badge.stopped {
  background: #fee2e2;
  color: #991b1b;
  font-weight: 600;
  border: 1px solid #fca5a5;
}

.status-badge.paused {
  background: #fef9c3;
  color: #854d0e;
  font-weight: 600;
  border: 1px solid #fde047;
}

.driver-badge {
  padding: 4px 8px;
  border-radius: 4px;
  background: #dbeafe;
  color: #1e40af;
  font-size: 12px;
}

.tag-badge {
  display: inline-block;
  padding: 3px 8px;
  margin: 2px;
  background: #e2e8f0;
  border-radius: 4px;
  font-size: 12px;
  color: #1e293b;
  font-weight: 500;
  border: 1px solid #cbd5e1;
}

.stats-mini {
  display: flex;
  flex-direction: column;
  gap: 2px;
  font-size: 12px;
}

.stat-item {
  display: flex;
  gap: 4px;
}

.stat-item .label {
  color: #475569;
  font-weight: 500;
}

.stat-item .value {
  color: #1e293b;
  font-weight: 600;
}

.action-group {
  display: flex;
  gap: 4px;
  justify-content: flex-end;
}

.icon-btn {
  padding: 6px 8px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  font-weight: 500;
}

.icon-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.icon-btn.success {
  border-color: #16a34a;
  color: #16a34a;
  background: #f0fdf4;
}

.icon-btn.success:hover {
  background: #dcfce7;
  border-color: #15803d;
}

.icon-btn.warning {
  border-color: #d97706;
  color: #d97706;
  background: #fef3c7;
}

.icon-btn.warning:hover {
  background: #fde68a;
  border-color: #b45309;
}

.icon-btn.info {
  border-color: #2563eb;
  color: #2563eb;
  background: #dbeafe;
}

.icon-btn.info:hover {
  background: #bfdbfe;
  border-color: #1d4ed8;
}

.icon-btn.danger {
  border-color: #dc2626;
  color: #dc2626;
  background: #fee2e2;
}

.icon-btn.danger:hover {
  background: #fecaca;
  border-color: #b91c1c;
}

.batch-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 8px;
  margin-top: 16px;
}

.batch-info {
  font-size: 14px;
  color: #374151;
}

.batch-buttons {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  font-size: 14px;
}

.btn:hover {
  background: #f9fafb;
}

.btn.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn.btn-ghost {
  border-color: transparent;
  background: transparent;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  border: 1px solid #e2e8f0;
}

.modal.large {
  max-width: 900px;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid #e2e8f0;
  background: #f8fafc;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
}

.btn-close {
  padding: 4px;
  border: none;
  background: none;
  cursor: pointer;
  color: #6b7280;
}

.modal-body {
  padding: 24px;
  overflow: auto;
  flex: 1;
  background: #fafbfc;
  color: #1e293b;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid #e2e8f0;
  background: #f8fafc;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  font-size: 14px;
  color: #374151;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #3b82f6;
}

.form-hint {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
}

.btn-primary {
  background: #3b82f6;
  color: white;
  border-color: #3b82f6;
}

.btn-primary:hover {
  background: #2563eb;
}

.btn-secondary {
  background: #6b7280;
  color: white;
  border-color: #6b7280;
}

.btn-secondary:hover {
  background: #4b5563;
}

/* 容器详情样式 */
.container-details {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.detail-section h4 {
  margin: 0 0 12px 0;
  font-size: 16px;
  font-weight: 600;
  color: #111827;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.detail-item {
  display: flex;
  flex-direction: column;
}

.detail-item label {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

.detail-item span {
  font-size: 14px;
  color: #111827;
}

.network-list,
.port-list,
.mount-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.network-item,
.port-item,
.mount-item {
  padding: 8px;
  background: #f9fafb;
  border-radius: 4px;
  font-size: 14px;
}

.network-name,
.mount-source {
  font-weight: 500;
  color: #111827;
}

.network-ip,
.mount-destination,
.mount-type {
  color: #6b7280;
  font-size: 12px;
}

/* 终端模态框 */
.terminal-modal .terminal-modal-content {
  max-width: 1000px;
  height: 80vh;
}

.terminal-body {
  background: #1e1e1e;
  padding: 0;
}

.terminal-container {
  height: 100%;
  min-height: 400px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  color: #00ff00;
  padding: 16px;
  overflow: auto;
}

.terminal-line {
  margin: 2px 0;
}

/* 日志模态框 */
.logs-container {
  background: #1e1e1e;
  padding: 16px;
  border-radius: 8px;
  max-height: 500px;
  overflow: auto;
  border: 1px solid #333;
}

.logs-container pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #00ff00;
  white-space: pre-wrap;
  line-height: 1.4;
}

/* 日志控件样式 */
.logs-controls {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 8px;
  flex-wrap: wrap;
  gap: 12px;
}

.logs-stats {
  display: flex;
  gap: 16px;
  align-items: center;
}

.log-stat {
  font-size: 13px;
  color: #374151;
  display: flex;
  align-items: center;
  gap: 4px;
}

.logs-filters {
  display: flex;
  gap: 12px;
  align-items: center;
}

.filter-select.sm {
  padding: 4px 8px;
  font-size: 12px;
  border: 1px solid #d1d5db;
  border-radius: 4px;
}

.auto-refresh-toggle {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #374151;
  cursor: pointer;
}

.logs-loading, .logs-empty {
  text-align: center;
  padding: 40px;
  color: #6b7280;
  font-size: 14px;
}

.logs-content {
  margin: 0;
}

/* 终端样式 */
.terminal-modal-content {
  max-width: 900px;
}

.terminal-body {
  padding: 0;
  background: #1e1e1e;
}

.terminal-container {
  background: #1e1e1e;
  padding: 16px;
  border-radius: 8px;
  min-height: 400px;
  max-height: 500px;
  overflow: auto;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.terminal-line {
  margin: 4px 0;
  line-height: 1.4;
}

.terminal-line.connected-message {
  color: #10b981;
  font-weight: 600;
}

.terminal-line.command-line {
  color: #fbbf24;
}

.terminal-line.command-line .prompt {
  color: #10b981;
  margin-right: 8px;
}

.terminal-line.output-line {
  color: #e5e7eb;
}

.terminal-line.error-line {
  color: #ef4444;
}

.terminal-disconnected {
  padding: 24px;
  text-align: center;
}

.terminal-info h4 {
  margin: 0 0 16px 0;
  color: #fbbf24;
  font-size: 16px;
}

.terminal-info p {
  color: #9ca3af;
  margin-bottom: 24px;
}

.terminal-features {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
  text-align: left;
}

.feature-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
}

.feature-icon {
  font-size: 20px;
  flex-shrink: 0;
}

.feature-item strong {
  display: block;
  margin-bottom: 4px;
  color: #e5e7eb;
}

.feature-item div {
  color: #9ca3af;
  font-size: 13px;
}

.terminal-commands {
  text-align: left;
  margin-top: 20px;
}

.terminal-commands h5 {
  color: #e5e7eb;
  margin: 0 0 12px 0;
}

.command-suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.cmd-btn {
  padding: 6px 12px;
  background: rgba(59, 130, 246, 0.2);
  border: 1px solid #3b82f6;
  color: #3b82f6;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.2s;
}

.cmd-btn:hover {
  background: rgba(59, 130, 246, 0.3);
  transform: translateY(-1px);
}

.terminal-status {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6b7280;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.status-dot.connected {
  background: #10b981;
  box-shadow: 0 0 8px #10b981;
}

.status-dot.disconnected {
  background: #ef4444;
  box-shadow: 0 0 8px #ef4444;
}

/* 容器详情样式 */
.container-details {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-section {
  background: #f9fafb;
  padding: 16px;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

.detail-section h4 {
  margin: 0 0 16px 0;
  font-size: 15px;
  font-weight: 600;
  color: #111827;
  border-bottom: 1px solid #e5e7eb;
  padding-bottom: 8px;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-item label {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
}

.detail-value {
  font-size: 14px;
  color: #111827;
  font-family: monospace;
  word-break: break-all;
}

.detail-value.status-value {
  font-family: inherit;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 4px;
  background: #dbeafe;
  color: #1e40af;
}

.network-list, .port-list, .mount-list, .env-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.network-item, .port-item, .mount-item, .env-item {
  background: white;
  padding: 12px;
  border-radius: 6px;
  border: 1px solid #e5e7eb;
}

.network-name, .port-mapping {
  font-weight: 600;
  color: #111827;
  margin-bottom: 4px;
}

.network-details {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.network-ip, .network-mac {
  font-size: 13px;
  color: #6b7280;
  font-family: monospace;
}

.mount-info, .env-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mount-source, .env-key {
  font-weight: 600;
  color: #111827;
  font-family: monospace;
}

.mount-destination, .env-value {
  font-size: 13px;
  color: #6b7280;
  font-family: monospace;
}

.mount-type, .mount-mode {
  font-size: 12px;
  color: #9ca3af;
}

.empty-state {
  text-align: center;
  padding: 24px;
  color: #9ca3af;
  font-size: 14px;
  font-style: italic;
}

.loading-state, .error-state {
  text-align: center;
  padding: 40px;
  color: #6b7280;
}

.spinner {
  border: 3px solid #e5e7eb;
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.action-btn.sm {
  padding: 6px 12px;
  font-size: 13px;
  border-radius: 4px;
}

/* 统计信息样式 */
.docker-stats {
  display: flex;
  gap: 20px;
  margin-top: 8px;
  font-size: 14px;
}

.docker-stats .stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.docker-stats .stat-icon {
  font-size: 16px;
}

.docker-stats .stat-value {
  font-weight: 600;
  color: #374151;
}

.docker-stats .stat-label {
  color: #6b7280;
  font-size: 12px;
}

/* 搜索框样式 */
.search-box {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 10px;
  color: #9ca3af;
  pointer-events: none;
}

.search-input {
  padding: 8px 12px 8px 32px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  font-size: 14px;
  min-width: 220px;
  color: #1e293b;
  background: white;
  font-weight: 500;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  background: #fafbfc;
}

.search-input::placeholder {
  color: #94a3b8;
  font-weight: 400;
}

/* 简化的状态指示器 - 移除可能导致错行的复杂样式 */

.status-badge.paused::before {
  background: #f59e0b;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 操作按钮组增强 */
.action-group {
  display: flex;
  gap: 4px;
  justify-content: flex-end;
  flex-wrap: wrap;
}

.icon-btn {
  transition: all 0.2s;
}

.icon-btn:hover {
  transform: scale(1.1);
}

.icon-btn.danger:hover {
  background: #fee2e2;
  color: #dc2626;
}

.icon-btn.success:hover {
  background: #d1fae5;
  color: #059669;
}

.icon-btn.warning:hover {
  background: #fef3c7;
  color: #d97706;
}

.icon-btn.info:hover {
  background: #dbeafe;
  color: #2563eb;
}

/* 批量操作区域增强 */
.batch-actions {
  margin-top: 16px;
  padding: 16px;
  background: #eff6ff;
  border-radius: 8px;
  border: 1px solid #bfdbfe;
  display: flex;
  justify-content: space-between;
  align-items: center;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.batch-info {
  font-weight: 500;
  color: #1e40af;
}

.batch-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* 表格增强 */
.docker-table {
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.docker-table tbody tr {
  transition: background-color 0.2s;
}

.docker-table tbody tr:hover {
  background-color: #f0f9ff;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .docker-table th,
  .docker-table td {
    padding: 8px 12px;
    font-size: 13px;
  }

  .action-group {
    flex-direction: column;
    gap: 2px;
  }

  .search-input {
    min-width: 150px;
  }
}
</style>