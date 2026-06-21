<template>
  <div class="bm">
    <!-- 顶部紧凑头部 -->
    <div class="bm-top">
      <div class="bm-title">
        <h1>Restic 备份</h1>
        <span class="bm-ping" :class="ping?.ok ? 'ok' : 'err'" :title="ping?.version">
          {{ ping?.ok ? `✓ ${ping.hostname}` : (ping ? '✗ 不可用' : '…') }}
        </span>
      </div>
      <div class="bm-top-actions">
        <button class="bm-btn ghost xs" @click="refreshAll" :disabled="loading">
          <ArrowPathIcon class="w-3.5 h-3.5" :class="{ spinning: loading }" />
        </button>
      </div>
      <nav class="bm-tabs">
        <button v-for="t in tabs" :key="t.id" class="bm-tab" :class="{ active: tab === t.id }" @click="tab = t.id">
          {{ t.label }}<span v-if="t.count" class="bm-badge">{{ t.count }}</span>
        </button>
      </nav>
    </div>

    <!-- ========== 快照 ========== -->
    <section v-show="tab === 'snaps'" class="bm-body">
      <div class="bm-toolbar">
        <select v-model="snapRepoId" @change="loadSnapshots" class="bm-input">
          <option :value="0">选择仓库…</option>
          <option v-for="r in repos" :key="r.id" :value="r.id">{{ r.name }}</option>
        </select>
        <input v-model="snapFilter" class="bm-input search" placeholder="过滤快照 (ID/标签/路径)…" />
        <button v-if="snapRepoId" class="bm-btn ghost xs" @click="loadSnapshots"><ArrowPathIcon class="w-3 h-3" /></button>
        <button v-if="selectedSnaps.size" class="bm-btn danger xs" @click="forgetSelected">
          <TrashIcon class="w-3 h-3" /> 删除 {{ selectedSnaps.size }}
        </button>
      </div>

      <div v-if="!snapRepoId" class="bm-empty"><CameraIcon class="w-12 h-12" /><p>请选择一个仓库</p></div>
      <div v-else-if="snapLoading" class="bm-empty"><CogIcon class="w-8 h-8 spinning" /></div>
      <div v-else-if="filteredSnaps.length === 0" class="bm-empty"><CameraIcon class="w-12 h-12" /><p>暂无快照</p></div>

      <table v-else class="bm-table">
        <thead>
          <tr>
            <th class="ck"><input type="checkbox" :checked="allSnapSelected" @change="toggleAllSnaps($event)" /></th>
            <th>ID</th><th>时间</th><th>主机</th><th>路径</th><th>标签</th><th>新增/总</th><th>数据量</th><th class="ops">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="s in filteredSnaps" :key="s.id" :class="{ sel: selectedSnaps.has(s.id) }">
            <td class="ck"><input type="checkbox" :checked="selectedSnaps.has(s.id)" @change="toggleSnap(s.id)" /></td>
            <td><code class="sid" :title="s.id">{{ s.short_id }}</code></td>
            <td>{{ formatDate(s.time) }}</td>
            <td>{{ s.hostname }}</td>
            <td class="paths" :title="s.paths?.join(', ')">{{ s.paths?.join(', ') }}</td>
            <td><span v-for="t in (s.tags||[])" :key="t" class="chip">{{ t }}</span></td>
            <td>{{ s.summary ? `${s.summary.files_new}/${s.summary.total_files_processed}` : '-' }}</td>
            <td>{{ s.summary ? formatBytes(s.summary.total_bytes_processed) : '-' }}</td>
            <td class="ops">
              <button class="bm-btn ghost xs" @click="openSnapDetail(s)" title="详情"><DocumentChartBarIcon class="w-3 h-3" /></button>
              <button class="bm-btn ghost xs" @click="openFiles(s)" title="浏览文件（含单文件恢复）"><FolderIcon class="w-3 h-3" /></button>
              <button class="bm-btn ghost xs" @click="openDiff(s)" title="对比"><ArrowsRightLeftIcon class="w-3 h-3" /></button>
              <button class="bm-btn warn xs" @click="openRestore(s)" title="恢复"><ArrowUturnLeftIcon class="w-3 h-3" /></button>
              <button class="bm-btn danger xs" @click="deleteSnap(s)" title="删除"><TrashIcon class="w-3 h-3" /></button>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- 快照内搜索 -->
      <div v-if="snapRepoId" class="find-box">
        <MagnifyingGlassIcon class="w-4 h-4" />
        <input v-model="findPattern" class="bm-input" placeholder='在快照中搜索文件，如 "report.pdf" 或 "*.jpg"'
               @keyup.enter="runFind" />
        <button class="bm-btn ghost xs" @click="runFind" :disabled="findLoading">
          <CogIcon v-if="findLoading" class="w-3 h-3 spinning" /> 搜索
        </button>
      </div>
      <pre v-if="findOutput !== null" class="bm-log">{{ findOutput || '（无匹配）' }}</pre>
    </section>

    <!-- ========== 仓库（含备份任务）========== -->
    <section v-show="tab === 'repos'" class="bm-body">
      <div class="bm-toolbar">
        <h3 style="margin:0;font-size:14px;">备份仓库与任务</h3>
        <input v-model="repoFilter" class="bm-input search" placeholder="过滤仓库…" />
        <button class="bm-btn primary xs" @click="openRepoEditor()">
          <PlusIcon class="w-3 h-3" /> 新建仓库
        </button>
      </div>

      <div v-if="filteredRepos.length === 0" class="bm-empty"><ArchiveBoxIcon class="w-12 h-12" /><p>暂无仓库</p></div>
      <div class="repo-grid">
        <div v-for="r in filteredRepos" :key="r.id" class="repo-card">
          <!-- 仓库头部 -->
          <div class="repo-head">
            <div class="repo-icon" :class="`t-${r.type}`"><component :is="repoIcon(r.type)" class="w-4 h-4" /></div>
            <div class="repo-name">
              <h3>{{ r.name }}</h3>
              <div class="repo-sub">{{ typeLabel(r.type) }} · {{ r.snapshotCount }} 快照 · {{ formatBytes(r.repoSize) }}</div>
            </div>
            <span class="chip" :class="`st-${r.status}`">{{ repoStatusLabel(r.status) }}</span>
          </div>
          <div class="repo-url" :title="r.url">{{ r.url }}</div>
          <div v-if="r.lastError" class="err-line" :title="r.lastError"><ExclamationTriangleIcon class="w-3 h-3" /> {{ truncate(r.lastError, 80) }}</div>

          <!-- 仓库维护操作 -->
          <div class="repo-ops">
            <button class="bm-btn ghost xs" @click="viewSnapshots(r)"><CameraIcon class="w-3 h-3" /> 快照</button>
            <button class="bm-btn ghost xs" @click="refreshRepo(r)" :disabled="busy[`r-${r.id}`]" title="刷新统计"><ArrowPathIcon class="w-3 h-3" :class="{ spinning: busy[`r-${r.id}`] }" /></button>
            <button class="bm-btn ghost xs" @click="testRepo(r)" :disabled="busy[`t-${r.id}`]" title="测试连接"><BoltIcon class="w-3 h-3" /></button>
            <button class="bm-btn ghost xs" @click="checkRepo(r)" :disabled="busy[`c-${r.id}`]" title="校验完整性"><ShieldCheckIcon class="w-3 h-3" :class="{ spinning: busy[`c-${r.id}`] }" /></button>
            <button class="bm-btn ghost xs" @click="unlockRepo(r)" :disabled="busy[`u-${r.id}`]" title="清除锁"><LockOpenIcon class="w-3 h-3" /></button>
            <button class="bm-btn ghost xs" @click="openRepoEditor(r)" title="编辑仓库">编辑</button>
            <button class="bm-btn danger xs" @click="deleteRepo(r)" title="删除仓库"><TrashIcon class="w-3 h-3" /></button>
          </div>

          <!-- 分隔 -->
          <div class="task-divider"></div>

          <!-- 备份任务（属于此仓库）-->
          <div class="task-section">
            <div class="task-section-head">
              <span>备份任务</span>
              <button class="bm-btn ghost xs" @click="openTaskEditor(undefined, r.id)"><PlusIcon class="w-3 h-3" /> 新建</button>
            </div>
            <div v-if="tasksForRepo(r.id).length === 0" class="task-empty">该仓库暂无备份任务</div>
            <div v-for="t in tasksForRepo(r.id)" :key="t.id" class="task-row">
              <div class="task-row-info">
                <div class="task-row-name">
                  {{ t.name }}
                  <span class="chip" :class="`st-${t.status}`">{{ statusLabel(t.status) }}</span>
                </div>
                <div class="task-row-meta">
                  <code>{{ t.sourcePath }}</code>
                  <span v-if="t.retention">· 保留 {{ formatRetention(t.retention) }}</span>
                  <span v-if="t.lastRun">· {{ formatDate(t.lastRun) }}</span>
                  <span v-if="t.lastDuration">· {{ formatDuration(t.lastDuration) }}</span>
                </div>
                <div v-if="t.lastError" class="err-line" :title="t.lastError"><ExclamationTriangleIcon class="w-3 h-3" /> {{ truncate(t.lastError, 60) }}</div>
              </div>
              <div class="task-row-ops">
                <button class="bm-btn primary xs" @click="runTask(t)" :disabled="t.status === 'running'">
                  <PlayIcon class="w-3 h-3" /> 立即备份
                </button>
                <button class="bm-btn ghost xs" @click="openLogs(t.id, t.name)" title="日志"><DocumentTextIcon class="w-3 h-3" /></button>
                <button class="bm-btn ghost xs" @click="openTaskEditor(t, r.id)" title="编辑">编辑</button>
                <button class="bm-btn danger xs" @click="deleteTask(t)" title="删除任务"><TrashIcon class="w-3 h-3" /></button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ========== 设置 ========== -->
    <section v-show="tab === 'settings'" class="bm-body">
      <div class="settings-grid">
        <!-- WebDAV 配置 -->
        <div class="settings-card">
          <h3><GlobeAltIcon class="w-4 h-4" /> WebDAV 云存储</h3>
          <p class="card-desc">配置一次，所有 rclone: 类型的仓库都复用此凭据。保存后会自动写入 <code>/data/rclone.conf</code>。</p>
          <div class="bm-row"><label>Remote 名称</label><input v-model="webdav.remoteName" class="bm-input" placeholder="123pan" /></div>
          <div class="bm-row"><label>WebDAV URL</label><input v-model="webdav.url" class="bm-input" placeholder="https://webdav.123pan.cn/webdav" /></div>
          <div class="bm-row">
            <label>Vendor 类型</label>
            <select v-model="webdav.vendor" class="bm-input">
              <option value="other">other (通用)</option>
              <option value="nextcloud">nextcloud</option>
              <option value="owncloud">owncloud</option>
              <option value="sharefile">sharefile</option>
              <option value="davpack">davpack</option>
              <option value="fastmail">fastmail</option>
            </select>
          </div>
          <div class="bm-row"><label>用户名</label><input v-model="webdav.username" class="bm-input" /></div>
          <div class="bm-row"><label>密码</label><input v-model="webdav.password" type="password" class="bm-input" placeholder="留空或 ******** 表示不修改" /></div>
          <div class="card-actions">
            <button class="bm-btn ghost xs" @click="testWebdav" :disabled="webdavBusy">
              <CogIcon v-if="webdavBusy" class="w-3 h-3 spinning" /> 测试连接
            </button>
            <button class="bm-btn primary xs" @click="saveWebdav" :disabled="webdavBusy">保存</button>
          </div>
          <div v-if="webdavTestResult" class="test-result" :class="webdavTestResult.ok ? 'ok' : 'err'">
            <pre>{{ webdavTestResult.output || webdavTestResult.error }}</pre>
          </div>
        </div>

        <!-- Restic 默认参数 -->
        <div class="settings-card">
          <h3><Cog6ToothIcon class="w-4 h-4" /> Restic 默认参数</h3>
          <p class="card-desc">以下参数会自动应用到所有备份任务。</p>
          <div class="bm-row"><label>默认主机名 (hostname)</label><input v-model="settings.defaultHostname" class="bm-input" /></div>
          <div class="bm-row"><label>默认标签 (逗号分隔)</label><input v-model="settings.defaultTags" class="bm-input" /></div>
          <div class="bm-row"><label>默认排除 (每行一条)</label><textarea v-model="settings.defaultExcludes" class="bm-input" rows="4"></textarea></div>
          <div class="bm-row checkbox">
            <input id="confirm-purge" v-model="settings.confirmPurge" type="checkbox" />
            <label for="confirm-purge">删除仓库时强制二次确认</label>
          </div>
          <div class="card-actions">
            <button class="bm-btn primary xs" @click="saveSettings" :disabled="settingsBusy">保存</button>
          </div>
        </div>

        <!-- 环境信息 -->
        <div class="settings-card">
          <h3><InformationCircleIcon class="w-4 h-4" /> 环境</h3>
          <div class="info-list">
            <div class="info-row"><span>Restic</span><code>{{ ping?.version || '-' }}</code></div>
            <div class="info-row"><span>主机名</span><code>{{ ping?.hostname || '-' }}</code></div>
            <div class="info-row"><span>缓存目录</span><code>{{ ping?.cacheDir || '-' }}</code></div>
            <div class="info-row"><span>运行时</span><code>{{ ping?.runtime || '-' }}</code></div>
          </div>
          <div class="info-tip">
            <InformationCircleIcon class="w-4 h-4" />
            <span>任务源路径是<b>容器内</b>路径。主机 <code>/data/photos</code> → 任务填 <code>/host/data/photos</code>。
                  恢复目标默认写到 <code>/restore/...</code> = 主机 <code>/var/restic-restore/...</code>。</span>
          </div>
        </div>
      </div>
    </section>

    <!-- ============ Modals ============ -->
    <!-- 仓库编辑器 -->
    <div v-if="modals.repo" class="bm-overlay" @click.self="modals.repo = false">
      <div class="bm-modal">
        <div class="bm-modal-head"><h3>{{ repoForm.id ? '编辑仓库' : '新建仓库' }}</h3>
          <button class="bm-x" @click="modals.repo = false"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body">
          <!-- WebDAV 快捷模式 -->
          <div v-if="!repoForm.id && webdav.url" class="bm-alert info">
            <InformationCircleIcon class="w-4 h-4" />
            <div>检测到 WebDAV 配置 <code>{{ webdav.remoteName }}</code>。
              类型选 <b>「WebDAV (复用配置)」</b> 只需填一个相对路径（如 <code>data</code> 或 <code>hserver</code>），凭据自动复用。</div>
          </div>
          <div class="bm-row"><label>仓库名</label><input v-model="repoForm.name" class="bm-input" placeholder="my-backup" /></div>
          <div class="bm-row">
            <label>仓库类型</label>
            <select v-model="repoForm.type" class="bm-input" :disabled="repoForm.id">
              <option value="webdav">WebDAV (复用上方配置)</option>
              <option value="local">本地 / 挂载目录</option>
              <option value="rclone">Rclone (完整 URL)</option>
              <option value="s3">S3 / MinIO</option>
              <option value="sftp">SFTP</option>
              <option value="rest">REST 服务器</option>
              <option value="b2">Backblaze B2</option>
            </select>
          </div>
          <div class="bm-row">
            <label>{{ repoForm.type === 'webdav' ? 'WebDAV 上的仓库路径（自动拼接为 rclone:' + (webdav.remoteName || 'webdav') + ':<此路径>)' : 'URL' }}</label>
            <input v-model="repoForm.url" class="bm-input" :placeholder="repoForm.type === 'webdav' ? 'data 或 hserver' : urlPlaceholder()" />
          </div>
          <div class="bm-row">
            <label>仓库密码</label>
            <input v-model="repoForm.password" type="password" class="bm-input" placeholder="加密密码" />
            <small class="hint" v-if="repoForm.id">留空表示不修改</small>
          </div>

          <!-- WebDAV 类型无需额外凭据（已复用上方） -->
          <div v-if="repoForm.type === 's3'" class="cred-section">
            <div class="cred-title">S3 凭据</div>
            <div class="bm-row"><label>Access Key</label><input v-model="env.AWS_ACCESS_KEY_ID" class="bm-input" /></div>
            <div class="bm-row"><label>Secret Key</label><input v-model="env.AWS_SECRET_ACCESS_KEY" type="password" class="bm-input" /></div>
          </div>
          <div v-if="repoForm.type === 'b2'" class="cred-section">
            <div class="cred-title">B2 凭据</div>
            <div class="bm-row"><label>Account ID</label><input v-model="env.B2_ACCOUNT_ID" class="bm-input" /></div>
            <div class="bm-row"><label>Account Key</label><input v-model="env.B2_ACCOUNT_KEY" type="password" class="bm-input" /></div>
          </div>
          <div v-if="repoForm.type === 'sftp'" class="cred-section">
            <div class="cred-title">SFTP</div>
            <div class="bm-row"><label>SSH 密码</label><input v-model="env.SSH_PASSWORD" type="password" class="bm-input" /></div>
            <div class="bm-row"><label>SSH 私钥 (可选)</label><textarea v-model="env.SSH_PRIVATE_KEY" class="bm-input" rows="3"></textarea></div>
          </div>

          <div v-if="!repoForm.id" class="bm-row checkbox">
            <input id="init-repo" v-model="repoForm.init" type="checkbox" />
            <label for="init-repo">立即 <code>restic init</code>（首次创建新仓库时勾选；导入已有仓库不要勾）</label>
          </div>
        </div>
        <div class="bm-modal-foot">
          <button class="bm-btn ghost" @click="modals.repo = false">取消</button>
          <button class="bm-btn primary" @click="saveRepo" :disabled="busy.saveRepo">保存</button>
        </div>
      </div>
    </div>

    <!-- 任务编辑器 -->
    <div v-if="modals.task" class="bm-overlay" @click.self="modals.task = false">
      <div class="bm-modal">
        <div class="bm-modal-head"><h3>{{ taskForm.id ? '编辑任务' : '新建备份任务' }}</h3>
          <button class="bm-x" @click="modals.task = false"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body">
          <div class="bm-row"><label>任务名</label><input v-model="taskForm.name" class="bm-input" /></div>
          <div class="bm-row"><label>目标仓库</label>
            <select v-model="taskForm.repoId" class="bm-input">
              <option v-for="r in repos" :key="r.id" :value="r.id">{{ r.name }} ({{ typeLabel(r.type) }})</option>
            </select>
          </div>
          <div class="bm-row"><label>源路径</label>
            <input v-model="taskForm.sourcePath" class="bm-input" placeholder="/host/data 或 /host/home/hserver" />
            <small class="hint">主机 <code>/data</code> 在容器内是 <code>/host/data</code></small>
          </div>
          <div class="bm-row"><label>标签</label><input v-model="taskForm.tags" class="bm-input" placeholder="data,auto" /></div>
          <div class="bm-row"><label>排除规则（每行一条）</label>
            <textarea v-model="taskForm.excludes" class="bm-input" rows="3" placeholder="**/.cache&#10;**/node_modules"></textarea>
          </div>
          <div class="bm-row">
            <label>保留策略（备份完成后 forget）</label>
            <div class="ret-grid">
              <label><input type="number" v-model.number="retention.keep_daily" min="0" /> 每日</label>
              <label><input type="number" v-model.number="retention.keep_weekly" min="0" /> 每周</label>
              <label><input type="number" v-model.number="retention.keep_monthly" min="0" /> 每月</label>
              <label><input type="number" v-model.number="retention.keep_yearly" min="0" /> 每年</label>
              <label><input type="number" v-model.number="retention.keep_last" min="0" /> 最近N个</label>
            </div>
          </div>
          <div class="bm-row checkbox">
            <input id="auto-prune" v-model="taskForm.autoPrune" type="checkbox" />
            <label for="auto-prune">forget 后执行 <code>prune</code> 释放空间</label>
          </div>
        </div>
        <div class="bm-modal-foot">
          <button class="bm-btn ghost" @click="modals.task = false">取消</button>
          <button class="bm-btn primary" @click="saveTask">保存</button>
        </div>
      </div>
    </div>

    <!-- 文件浏览器 + 单文件恢复 -->
    <div v-if="modals.files" class="bm-overlay" @click.self="modals.files = false">
      <div class="bm-modal lg">
        <div class="bm-modal-head"><h3>快照文件 · {{ filesSnap?.short_id }}</h3>
          <button class="bm-x" @click="modals.files = false"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body">
          <div class="files-toolbar">
            <input v-model="filesFilter" class="bm-input" placeholder="按路径过滤…" />
            <span class="muted">{{ filteredFiles.length }} / {{ snapFiles.length }}</span>
            <button class="bm-btn warn xs" @click="openRestore(filesSnap!)">恢复整个快照</button>
          </div>
          <div v-if="filesLoading" class="bm-empty"><CogIcon class="w-8 h-8 spinning" /></div>
          <table v-else class="bm-table compact">
            <thead><tr><th>类型</th><th>路径</th><th>大小</th><th>修改时间</th><th class="ops">操作</th></tr></thead>
            <tbody>
              <tr v-for="(f, i) in filteredFiles" :key="i">
                <td><span class="chip" :class="`ft-${f.type}`">{{ f.type }}</span></td>
                <td class="paths" :title="f.path">{{ f.path }}</td>
                <td>{{ f.type === 'file' ? formatBytes(f.size || 0) : '-' }}</td>
                <td class="muted">{{ f.mtime ? formatDate(f.mtime) : '-' }}</td>
                <td class="ops">
                  <button class="bm-btn warn xs" @click="restoreOneFile(f)" :title="f.type === 'file' ? '恢复此文件' : '恢复此目录'">
                    <ArrowUturnLeftIcon class="w-3 h-3" /> 恢复
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="bm-modal-foot"><button class="bm-btn ghost" @click="modals.files = false">关闭</button></div>
      </div>
    </div>

    <!-- 恢复对话框 -->
    <div v-if="modals.restore" class="bm-overlay" @click.self="modals.restore = false">
      <div class="bm-modal">
        <div class="bm-modal-head"><h3>恢复 {{ restoreForm.snapshotId }}</h3>
          <button class="bm-x" @click="modals.restore = false"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body">
          <div class="bm-alert warn"><ExclamationTriangleIcon class="w-4 h-4" /><div>建议先恢复到 /restore 子目录验证再使用。</div></div>
          <div v-if="restoreForm.hint" class="bm-alert info"><InformationCircleIcon class="w-4 h-4" /><div>{{ restoreForm.hint }}</div></div>
          <div class="bm-row"><label>快照 ID</label><input v-model="restoreForm.snapshotId" class="bm-input" placeholder="留空 = latest" /></div>
          <div class="bm-row"><label>恢复到</label>
            <input v-model="restoreForm.target" class="bm-input" />
            <small class="hint">容器 <code>/restore</code> = 主机 <code>/var/restic-restore</code></small>
          </div>
          <div class="bm-row"><label>仅包含（每行一条 glob）</label>
            <textarea v-model="restoreForm.includeRaw" class="bm-input" rows="3" placeholder="/data/photos/**&#10;/data/important.db"></textarea>
          </div>
        </div>
        <div class="bm-modal-foot">
          <button class="bm-btn ghost" @click="modals.restore = false">取消</button>
          <button class="bm-btn warn" @click="doRestore">开始恢复</button>
        </div>
      </div>
    </div>

    <!-- 快照详情 -->
    <div v-if="modals.snapDetail" class="bm-overlay" @click.self="modals.snapDetail = false">
      <div class="bm-modal lg">
        <div class="bm-modal-head"><h3>快照详情 · {{ snapDetail?.snapshot?.short_id }}</h3>
          <button class="bm-x" @click="modals.snapDetail = false"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body" v-if="snapDetail">
          <div class="detail-grid">
            <div class="detail-card">
              <h4>基本信息</h4>
              <div class="info-row"><span>完整 ID</span><code>{{ snapDetail.snapshot.id }}</code></div>
              <div class="info-row"><span>时间</span><code>{{ formatDate(snapDetail.snapshot.time) }}</code></div>
              <div class="info-row"><span>主机</span><code>{{ snapDetail.snapshot.hostname }}</code></div>
              <div class="info-row"><span>路径</span><code>{{ snapDetail.snapshot.paths?.join(', ') }}</code></div>
            </div>
            <div class="detail-card">
              <h4>文件统计</h4>
              <div class="info-row"><span>文件总数</span><code>{{ snapDetail.fileCount }}</code></div>
              <div class="info-row"><span>目录总数</span><code>{{ snapDetail.dirCount }}</code></div>
              <div class="info-row"><span>总大小</span><code>{{ formatBytes(snapDetail.totalSize) }}</code></div>
            </div>
          </div>
          <div class="bm-modal-foot inline">
            <button class="bm-btn ghost xs" @click="openFiles(snapDetail.snapshot)">浏览文件</button>
            <button class="bm-btn warn xs" @click="openRestore(snapDetail.snapshot)">恢复</button>
          </div>
        </div>
      </div>
    </div>

    <!-- 对比 -->
    <div v-if="modals.diff" class="bm-overlay" @click.self="modals.diff = false">
      <div class="bm-modal lg">
        <div class="bm-modal-head"><h3>对比快照</h3>
          <button class="bm-x" @click="modals.diff = false"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body">
          <div class="bm-row"><label>基准 A</label>
            <select v-model="diffForm.a" class="bm-input">
              <option v-for="s in snapshots" :key="s.id" :value="s.short_id">{{ s.short_id }} · {{ formatDate(s.time) }}</option>
            </select>
          </div>
          <div class="bm-row"><label>目标 B</label>
            <select v-model="diffForm.b" class="bm-input">
              <option v-for="s in snapshots" :key="s.id" :value="s.short_id">{{ s.short_id }} · {{ formatDate(s.time) }}</option>
            </select>
          </div>
          <div v-if="diffOutput !== null"><pre class="bm-log">{{ diffOutput || '（无差异）' }}</pre></div>
        </div>
        <div class="bm-modal-foot">
          <button class="bm-btn ghost" @click="modals.diff = false">关闭</button>
          <button class="bm-btn primary xs" @click="runDiff" :disabled="diffLoading">对比</button>
        </div>
      </div>
    </div>

    <!-- 同步任务编辑器 -->
    <!-- 已移除：仓库间同步对单云盘场景没用。如需恢复，见 git 历史 / API: /api/storage/backup/sync-jobs -->

    <!-- 日志 -->
    <div v-if="modals.logs" class="bm-overlay" @click.self="closeLogs">
      <div class="bm-modal lg">
        <div class="bm-modal-head"><h3>日志：{{ logTitle }}</h3>
          <button class="bm-x" @click="closeLogs"><XMarkIcon class="w-5 h-5" /></button></div>
        <div class="bm-modal-body">
          <pre ref="logBox" class="bm-log">{{ logLines.length ? logLines.join('\n') : '（暂无日志）' }}</pre>
        </div>
        <div class="bm-modal-foot">
          <button class="bm-btn ghost xs" @click="reloadLogs">刷新</button>
          <button class="bm-btn ghost" @click="closeLogs">关闭</button>
        </div>
      </div>
    </div>

    <transition name="fade"><div v-if="toast.show" class="bm-toast" :class="`t-${toast.type}`">{{ toast.text }}</div></transition>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { resticApi } from '@/api'
import type { BackupRepo, ResticSnapshot, BackupTask, LSNode, PingInfo, BackupSettings, WebDAVProfile } from '@/api'
import {
  ArchiveBoxIcon, ArrowPathIcon, PlusIcon, CogIcon, CameraIcon, FolderIcon, ClockIcon,
  ExclamationTriangleIcon, ShieldCheckIcon, ArrowUturnLeftIcon, TrashIcon, XMarkIcon,
  PlayIcon, DocumentTextIcon, Cog6ToothIcon, ArrowsRightLeftIcon, InformationCircleIcon,
  ServerIcon, CloudIcon, GlobeAltIcon, BoltIcon, LockOpenIcon, DocumentChartBarIcon,
  MagnifyingGlassIcon,
} from '@heroicons/vue/24/outline'

// ===== State =====
const loading = ref(false)
const busy = reactive<Record<string, boolean>>({})
const ping = ref<PingInfo | null>(null)
const repos = ref<BackupRepo[]>([])
const tasks = ref<BackupTask[]>([])
const snapshots = ref<ResticSnapshot[]>([])
const snapRepoId = ref(0)
const snapLoading = ref(false)
const tab = ref<'snaps' | 'repos' | 'settings'>('repos')
const selectedSnaps = reactive(new Set<string>())

// Filters
const repoFilter = ref('')
const snapFilter = ref('')
const filesFilter = ref('')
const filteredRepos = computed(() => {
  const q = repoFilter.value.toLowerCase().trim()
  if (!q) return repos.value
  return repos.value.filter(r => r.name.toLowerCase().includes(q) || r.url.toLowerCase().includes(q))
})
const filteredSnaps = computed(() => {
  const q = snapFilter.value.toLowerCase().trim()
  if (!q) return snapshots.value
  return snapshots.value.filter(s =>
    s.short_id.includes(q) || s.hostname.toLowerCase().includes(q) ||
    (s.tags || []).some(t => t.toLowerCase().includes(q)) ||
    (s.paths || []).some(p => p.toLowerCase().includes(q)))
})
const filteredFiles = computed(() => {
  const q = filesFilter.value.toLowerCase().trim()
  if (!q) return snapFiles.value
  return snapFiles.value.filter(f => f.path.toLowerCase().includes(q))
})
const allSnapSelected = computed(() => filteredSnaps.value.length > 0 && filteredSnaps.value.every(s => selectedSnaps.has(s.id)))

const tabs = computed(() => [
  { id: 'snaps' as const, label: '快照', count: snapshots.value.length },
  { id: 'repos' as const, label: '仓库', count: repos.value.length },
  { id: 'settings' as const, label: '设置', count: undefined },
])

// Toast
const toast = reactive({ show: false, text: '', type: 'info' as 'info' | 'success' | 'error' })
let toastTimer: any = null
function showToast(text: string, type: 'info' | 'success' | 'error' = 'info') {
  toast.text = text; toast.type = type; toast.show = true
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => (toast.show = false), 3500)
}

// Modals
const modals = reactive({ repo: false, task: false, restore: false, files: false, logs: false, snapDetail: false, diff: false })

// Forms
const repoForm = reactive({ id: 0, name: '', type: 'webdav' as 'webdav' | 'local' | 'rclone' | 's3' | 'sftp' | 'rest' | 'b2', url: '', password: '', init: false })
const env = reactive<Record<string, string>>({})
const taskForm = reactive({ id: 0, name: '', repoId: 0, sourcePath: '', tags: '', excludes: '', autoPrune: false })
const retention = reactive<Record<string, number>>({ keep_daily: 7, keep_weekly: 4, keep_monthly: 6, keep_yearly: 0, keep_last: 0 })
const restoreForm = reactive({ repoId: 0, snapshotId: '', target: '', includeRaw: '', hint: '' })

// Files / detail / diff / find
const snapFiles = ref<LSNode[]>([])
const filesLoading = ref(false)
const filesSnap = ref<ResticSnapshot | null>(null)
const snapDetail = ref<any>(null)
const diffForm = reactive({ a: '', b: '' })
const diffOutput = ref<string | null>(null)
const diffLoading = ref(false)
const findPattern = ref('')
const findOutput = ref<string | null>(null)
const findLoading = ref(false)

// Logs
const logTitle = ref('')
const logLines = ref<string[]>([])
const logCtx = reactive({ id: 0 })
let logTimer: any = null
const logBox = ref<HTMLPreElement | null>(null)

// Settings
const settings = reactive<BackupSettings>({ defaultHostname: 'nas-dashboard', defaultExcludes: '', defaultTags: '', autoCheck: false, confirmPurge: true })
const settingsBusy = ref(false)
const webdav = reactive<WebDAVProfile>({ remoteName: '', url: '', vendor: 'other', username: '', password: '' })
const webdavBusy = ref(false)
const webdavTestResult = ref<{ ok: boolean; output?: string; error?: string } | null>(null)

onMounted(() => { refreshAll(); loadWebdav() })
onUnmounted(() => { if (logTimer) clearInterval(logTimer) })

async function refreshAll() {
  loading.value = true
  try {
    const [p, r, t, st] = await Promise.all([
      resticApi.ping().catch(() => null),
      resticApi.listRepos().catch(() => ({ repos: [], total: 0 })),
      resticApi.listTasks().catch(() => ({ tasks: [], total: 0 })),
      resticApi.getSettings().catch(() => null),
    ])
    ping.value = p as PingInfo
    repos.value = (r as any).repos || []
    tasks.value = (t as any).tasks || []
    if (st) Object.assign(settings, st)
    if (!snapRepoId.value && repos.value.length > 0) snapRepoId.value = repos.value[0].id
    if (tab.value === 'snaps' && snapRepoId.value) loadSnapshots()
  } finally { loading.value = false }
}

async function loadSnapshots() {
  if (!snapRepoId.value) { snapshots.value = []; return }
  snapLoading.value = true
  selectedSnaps.clear()
  findOutput.value = null
  try {
    const res = await resticApi.listSnapshots(snapRepoId.value) as any
    snapshots.value = res.snapshots || []
  } catch (e: any) { showToast('读取快照失败：' + (e.message || e), 'error'); snapshots.value = [] }
  finally { snapLoading.value = false }
}

async function refreshRepo(r: BackupRepo) {
  busy[`r-${r.id}`] = true
  try { await resticApi.refreshRepo(r.id); const l = await resticApi.listRepos() as any; repos.value = l.repos || []; showToast(`"${r.name}" 已刷新`, 'success') }
  catch (e: any) { showToast('刷新失败：' + (e.message || e), 'error') }
  finally { busy[`r-${r.id}`] = false }
}
async function testRepo(r: BackupRepo) {
  busy[`t-${r.id}`] = true
  try { const res = await resticApi.testRepo(r.id) as any; if (res.ok) showToast(`连接成功：${res.snapshotCount} 个快照`, 'success'); else showToast('连接失败：' + res.error, 'error') }
  catch (e: any) { showToast('测试失败：' + (e.message || e), 'error') }
  finally { busy[`t-${r.id}`] = false }
}
async function unlockRepo(r: BackupRepo) {
  const force = confirm(`清除 "${r.name}" 的锁？\n确定 = 强制清除所有锁\n取消 = 仅清过期锁`)
  busy[`u-${r.id}`] = true
  try { const res = await resticApi.unlockRepo(r.id, force) as any; showToast(`已清除锁`, 'success') }
  catch (e: any) { showToast('解锁失败：' + (e.message || e), 'error') }
  finally { busy[`u-${r.id}`] = false }
}
async function checkRepo(r: BackupRepo) {
  busy[`c-${r.id}`] = true
  try { await resticApi.checkRepo(r.id, false); showToast(`"${r.name}" 校验通过`, 'success') }
  catch (e: any) { showToast('校验失败：' + (e.message || e), 'error') }
  finally { busy[`c-${r.id}`] = false }
}
async function deleteRepo(r: BackupRepo) {
  if (!confirm(`删除仓库 "${r.name}"？`)) return
  const purge = confirm(`同时删除磁盘数据？\n确定 = 同时删除\n取消 = 仅删配置`)
  try { await resticApi.deleteRepo(r.id, purge); await refreshAll(); showToast('已删除', 'success') }
  catch (e: any) { showToast('删除失败：' + (e.message || e), 'error') }
}

function viewSnapshots(r: BackupRepo) { snapRepoId.value = r.id; tab.value = 'snaps'; loadSnapshots() }

// ===== Repo editor =====
function openRepoEditor(r?: BackupRepo) {
  Object.keys(env).forEach(k => delete env[k])
  if (r) {
    Object.assign(repoForm, { id: r.id, name: r.name, type: r.type === 'rclone' && webdav.value?.url && r.url.startsWith('rclone:' + (webdav.remoteName || 'webdav') + ':') ? 'webdav' : r.type, url: r.url, password: '', init: false })
    // webdav 类型：从 url 提取 path
    if (repoForm.type === 'webdav') {
      const prefix = 'rclone:' + (webdav.remoteName || 'webdav') + ':'
      if (r.url.startsWith(prefix)) repoForm.url = r.url.slice(prefix.length)
    }
  } else {
    Object.assign(repoForm, { id: 0, name: '', type: webdav.url ? 'webdav' : 'local', url: '', password: '', init: false })
  }
  modals.repo = true
}

function urlPlaceholder() {
  switch (repoForm.type) {
    case 'local': return '/data/restic/myrepo'
    case 'rclone': return 'rclone:123pan:data'
    case 's3': return 's3.amazonaws.com/bucket/path'
    case 'sftp': return 'user@host:/path'
    case 'rest': return 'http://server:8000/'
    case 'b2': return 'bucket:/path'
    default: return ''
  }
}

async function saveRepo() {
  if (!repoForm.name || !repoForm.url) { showToast('请填完整', 'error'); return }
  if (!repoForm.id && !repoForm.password) { showToast('请输入密码', 'error'); return }
  busy.saveRepo = true
  try {
    const cleanEnv: Record<string, string> = {}
    Object.entries(env).forEach(([k, v]) => { if (v) cleanEnv[k] = v })
    // webdav → rclone 转换
    let actualType = repoForm.type
    let actualUrl = repoForm.url
    if (repoForm.type === 'webdav') {
      actualType = 'rclone'
      const remoteName = webdav.remoteName || 'webdav'
      actualUrl = `rclone:${remoteName}:${repoForm.url}`
      // 关键：必须告诉 restic 用我们管理的 rclone 配置
      cleanEnv.RCLONE_CONFIG = '/data/rclone.conf'
    }
    const payload: any = { name: repoForm.name, type: actualType, url: actualUrl, env: cleanEnv, init: repoForm.init && !repoForm.id }
    if (repoForm.password) payload.password = repoForm.password
    if (repoForm.id) { await resticApi.updateRepo(repoForm.id, payload); showToast('已更新', 'success') }
    else { await resticApi.createRepo(payload); showToast('已创建', 'success') }
    modals.repo = false; await refreshAll()
  } catch (e: any) { showToast('保存失败：' + (e.message || e), 'error') }
  finally { busy.saveRepo = false }
}

// ===== Snapshots =====
function openRestore(s: ResticSnapshot) {
  Object.assign(restoreForm, { repoId: snapRepoId.value, snapshotId: s.short_id, target: `/restore/${s.short_id}`, includeRaw: '', hint: '' })
  modals.restore = true
}
function restoreOneFile(f: LSNode) {
  if (!filesSnap.value) return
  Object.assign(restoreForm, {
    repoId: snapRepoId.value, snapshotId: filesSnap.value.short_id,
    includeRaw: f.type === 'dir' ? f.path + '/**' : f.path,
    target: `/restore/single/${f.name}`,
    hint: f.type === 'file' ? `将只恢复文件：${f.path}` : `将恢复目录下所有内容：${f.path}/**`,
  })
  modals.files = false; modals.restore = true
}
async function doRestore() {
  if (!restoreForm.target) { showToast('请填目标路径', 'error'); return }
  try {
    const includes = restoreForm.includeRaw.split('\n').map(s => s.trim()).filter(Boolean)
    await resticApi.restore(restoreForm.repoId, { snapshotId: restoreForm.snapshotId || 'latest', target: restoreForm.target, include: includes })
    modals.restore = false; showToast('恢复已在后台启动', 'success')
  } catch (e: any) { showToast('恢复失败：' + (e.message || e), 'error') }
}
async function openFiles(s: ResticSnapshot) {
  filesSnap.value = s; modals.files = true; filesFilter.value = ''; filesLoading.value = true
  try { const res = await resticApi.listSnapshotFiles(snapRepoId.value, s.short_id) as any; snapFiles.value = res.files || [] }
  catch (e: any) { showToast('读取失败：' + (e.message || e), 'error'); snapFiles.value = [] }
  finally { filesLoading.value = false }
}
async function openSnapDetail(s: ResticSnapshot) {
  modals.snapDetail = true; snapDetail.value = null
  try { snapDetail.value = await resticApi.snapshotDetail(snapRepoId.value, s.short_id) as any }
  catch (e: any) { showToast('读取失败：' + (e.message || e), 'error'); modals.snapDetail = false }
}
function openDiff(s: ResticSnapshot) {
  diffForm.b = s.short_id
  const idx = snapshots.value.findIndex(x => x.id === s.id)
  diffForm.a = idx > 0 ? snapshots.value[idx - 1].short_id : s.short_id
  diffOutput.value = null; modals.diff = true
}
async function runDiff() {
  if (!diffForm.a || !diffForm.b) return
  diffLoading.value = true
  try { const res = await resticApi.diffSnapshots(snapRepoId.value, diffForm.a, diffForm.b) as any; diffOutput.value = (res.output || '').trim() }
  catch (e: any) { diffOutput.value = '对比失败：' + (e.message || e) }
  finally { diffLoading.value = false }
}
async function deleteSnap(s: ResticSnapshot) {
  const prune = confirm(`同时执行 prune？\n确定 = 是\n取消 = 仅 forget`)
  if (!confirm(`确认删除快照 ${s.short_id}？`)) return
  try { await resticApi.deleteSnapshot(snapRepoId.value, s.short_id, prune); showToast('已删除', 'success'); loadSnapshots() }
  catch (e: any) { showToast('删除失败：' + (e.message || e), 'error') }
}
function toggleAllSnaps(e: Event) {
  const checked = (e.target as HTMLInputElement).checked
  if (checked) filteredSnaps.value.forEach(s => selectedSnaps.add(s.id))
  else selectedSnaps.clear()
}
function toggleSnap(id: string) { if (selectedSnaps.has(id)) selectedSnaps.delete(id); else selectedSnaps.add(id) }
async function forgetSelected() {
  const ids = Array.from(selectedSnaps)
  if (!confirm(`确认删除 ${ids.length} 个快照？`)) return
  let ok = 0, fail = 0
  for (const id of ids) {
    const snap = snapshots.value.find(s => s.id === id)
    if (!snap) continue
    try { await resticApi.deleteSnapshot(snapRepoId.value, snap.short_id, false); ok++ } catch { fail++ }
  }
  showToast(`已删 ${ok}${fail ? `,失败 ${fail}` : ''}`, fail ? 'error' : 'success')
  selectedSnaps.clear(); loadSnapshots()
}
async function runFind() {
  if (!findPattern.value || !snapRepoId.value) return
  findLoading.value = true
  try { const res = await resticApi.findInSnapshots(snapRepoId.value, findPattern.value) as any; findOutput.value = res.output || '' }
  catch (e: any) { findOutput.value = '搜索失败：' + (e.message || e) }
  finally { findLoading.value = false }
}

// ===== Tasks =====
// 返回属于某仓库的所有任务（用于仓库卡片内联展示）
function tasksForRepo(repoId: number): BackupTask[] {
  return tasks.value.filter(t => t.repoId === repoId)
}

function openTaskEditor(t?: BackupTask, repoId?: number) {
  if (t) {
    Object.assign(taskForm, { id: t.id, name: t.name, repoId: t.repoId, sourcePath: t.sourcePath, tags: t.tags || '', excludes: t.excludes || '', autoPrune: t.autoPrune })
    Object.keys(retention).forEach(k => { (retention as any)[k] = t.retention?.[k] ?? 0 })
  } else {
    Object.assign(taskForm, {
      id: 0, name: '',
      repoId: repoId || repos.value[0]?.id || 0,
      sourcePath: '', tags: settings.defaultTags || '', excludes: '', autoPrune: false,
    })
    Object.assign(retention, { keep_daily: 7, keep_weekly: 4, keep_monthly: 6, keep_yearly: 0, keep_last: 0 })
  }
  modals.task = true
}
async function saveTask() {
  if (!taskForm.name || !taskForm.repoId || !taskForm.sourcePath) { showToast('请填完整', 'error'); return }
  const r: Record<string, number> = {}
  Object.entries(retention).forEach(([k, v]) => { if (v > 0) r[k] = v })
  const payload = { name: taskForm.name, repoId: taskForm.repoId, sourcePath: taskForm.sourcePath, tags: taskForm.tags, excludes: taskForm.excludes, autoPrune: taskForm.autoPrune, retention: r, enabled: true }
  try {
    if (taskForm.id) await resticApi.updateTask(taskForm.id, payload); else await resticApi.createTask(payload)
    showToast('已保存', 'success'); modals.task = false
    const res = await resticApi.listTasks() as any; tasks.value = res.tasks || []
  } catch (e: any) { showToast('保存失败：' + (e.message || e), 'error') }
}
async function deleteTask(t: BackupTask) {
  if (!confirm(`删除任务 "${t.name}"？`)) return
  try { await resticApi.deleteTask(t.id); const res = await resticApi.listTasks() as any; tasks.value = res.tasks || []; showToast('已删除', 'success') }
  catch (e: any) { showToast('删除失败：' + (e.message || e), 'error') }
}
async function runTask(t: BackupTask) {
  try { await resticApi.runTask(t.id); showToast(`"${t.name}" 已启动`, 'success'); t.status = 'running'; pollTaskStatus(t.id); openLogs(t.id, t.name) }
  catch (e: any) { showToast('启动失败：' + (e.message || e), 'error') }
}
function pollTaskStatus(taskId: number) {
  let attempts = 0
  const timer = setInterval(async () => {
    attempts++
    try {
      const s = await resticApi.taskStatus(taskId) as any
      const t = tasks.value.find(x => x.id === taskId)
      if (t) { t.status = s.status; t.lastRun = s.lastRun; t.lastError = s.lastError; t.lastSnapshotId = s.lastSnapshotId; t.lastDuration = s.lastDuration }
      if (s.status !== 'running' || attempts > 360) clearInterval(timer)
    } catch { clearInterval(timer) }
  }, 3000)
}

// ===== Logs =====
async function openLogs(id: number, name: string) {
  logCtx.id = id; logTitle.value = name; modals.logs = true
  await reloadLogs()
  if (logTimer) clearInterval(logTimer)
  logTimer = setInterval(reloadLogs, 2000)
}
async function reloadLogs() {
  try {
    const res = await resticApi.taskLogs(logCtx.id) as any
    logLines.value = res.lines || []
    await nextTick()
    if (logBox.value) logBox.value.scrollTop = logBox.value.scrollHeight
  } catch { /* ignore */ }
}
function closeLogs() { modals.logs = false; if (logTimer) { clearInterval(logTimer); logTimer = null } }

// ===== Settings =====
async function saveSettings() {
  settingsBusy.value = true
  try { await resticApi.updateSettings(settings); showToast('设置已保存', 'success') }
  catch (e: any) { showToast('保存失败：' + (e.message || e), 'error') }
  finally { settingsBusy.value = false }
}
async function loadWebdav() {
  try { const p = await resticApi.getWebDAV() as any; Object.assign(webdav, p) }
  catch { /* 没配置过，忽略 */ }
}
async function saveWebdav() {
  if (!webdav.url) { showToast('请填 WebDAV URL', 'error'); return }
  webdavBusy.value = true
  try {
    await resticApi.updateWebDAV(webdav)
    showToast('WebDAV 配置已保存（写入 /data/rclone.conf）', 'success')
    webdav.password = '********'
  } catch (e: any) { showToast('保存失败：' + (e.message || e), 'error') }
  finally { webdavBusy.value = false }
}
async function testWebdav() {
  if (!webdav.url) { showToast('请填 WebDAV URL', 'error'); return }
  webdavBusy.value = true; webdavTestResult.value = null
  try {
    const res = await resticApi.testWebDAV(webdav) as any
    webdavTestResult.value = res
    showToast(res.ok ? '连接成功' : '连接失败', res.ok ? 'success' : 'error')
  } catch (e: any) { showToast('测试失败：' + (e.message || e), 'error') }
  finally { webdavBusy.value = false }
}

// ===== Helpers =====
function formatDate(s?: string) { return s && !s.startsWith('0001-') ? new Date(s).toLocaleString('zh-CN', { hour12: false }) : '-' }
function formatBytes(b?: number) { if (!b) return '0 B'; const u = ['B', 'KB', 'MB', 'GB', 'TB']; const i = Math.floor(Math.log(b) / Math.log(1024)); return (b / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1) + ' ' + u[i] }
function formatRetention(r: Record<string, number>) { return Object.entries(r).map(([k, v]) => `${k.replace('keep_', '')}=${v}`).join(' ') }
function formatDuration(sec: number) { return sec < 60 ? sec.toFixed(1) + 's' : Math.floor(sec / 60) + 'm' + Math.floor(sec % 60) + 's' }
function truncate(s: string, n: number) { return s && s.length > n ? s.slice(0, n) + '…' : s }
function statusLabel(s: string) { return ({ idle: '空闲', running: '运行中', completed: '完成', failed: '失败' } as any)[s] || s }
function repoStatusLabel(s: string) { return ({ active: '正常', uninitialized: '未初始化', error: '错误' } as any)[s] || s }
function typeLabel(t: string) { return ({ local: '本地', s3: 'S3', sftp: 'SFTP', rest: 'REST', b2: 'B2', rclone: 'Rclone', webdav: 'WebDAV' } as any)[t] || t }
function repoIcon(t: string) { return ({ local: ServerIcon, s3: CloudIcon, sftp: ServerIcon, rest: GlobeAltIcon, b2: CloudIcon, rclone: CloudIcon, webdav: GlobeAltIcon } as any)[t] || ServerIcon }
</script>

<style scoped>
.bm { width: 100%; height: 100%; background: #f8fafc; color: #0f172a; display: flex; flex-direction: column; overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", sans-serif; }

/* 顶部紧凑头部 + 内联 tabs */
.bm-top { display: flex; align-items: center; gap: 16px; padding: 10px 20px; background: linear-gradient(135deg, #0f172a, #1e293b); color: white; flex-wrap: wrap; }
.bm-title { display: flex; align-items: center; gap: 10px; }
.bm-title h1 { font-size: 16px; font-weight: 600; margin: 0; }
.bm-ping { font-size: 11px; padding: 3px 8px; border-radius: 999px; background: rgba(255,255,255,0.08); }
.bm-ping.ok { color: #86efac; }
.bm-ping.err { color: #fca5a5; }
.bm-top-actions { margin-left: auto; }
.bm-tabs { display: flex; gap: 2px; margin-left: 8px; }
.bm-tab { background: rgba(255,255,255,0.08); color: rgba(255,255,255,0.7); border: none; cursor: pointer;
  padding: 6px 14px; border-radius: 999px; font-size: 12px; font-weight: 500; transition: all .15s; }
.bm-tab:hover { background: rgba(255,255,255,0.15); color: white; }
.bm-tab.active { background: #0ea5e9; color: white; }
.bm-badge { background: rgba(255,255,255,0.25); padding: 0 6px; border-radius: 999px; margin-left: 4px; font-size: 10px; }

/* 主体 */
.bm-body { flex: 1; overflow: auto; padding: 16px 20px; }
.bm-toolbar { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; flex-wrap: wrap; }
.bm-toolbar .search { max-width: 220px; margin-left: auto; }
.flex-grow { flex: 1; }

/* 按钮 */
.bm-btn { display: inline-flex; align-items: center; gap: 4px; padding: 6px 12px; border-radius: 6px; border: none; cursor: pointer; font-size: 12px; font-weight: 500; transition: all .15s; white-space: nowrap; }
.bm-btn.xs { padding: 4px 8px; font-size: 11px; }
.bm-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.bm-btn.primary { background: #0ea5e9; color: white; }
.bm-btn.primary:hover:not(:disabled) { background: #0284c7; }
.bm-btn.ghost { background: white; color: #334155; border: 1px solid #cbd5e1; }
.bm-btn.ghost:hover:not(:disabled) { background: #f1f5f9; }
.bm-btn.warn { background: #f59e0b; color: white; }
.bm-btn.warn:hover:not(:disabled) { background: #d97706; }
.bm-btn.danger { background: white; color: #ef4444; border: 1px solid #fecaca; }
.bm-btn.danger:hover:not(:disabled) { background: #fee2e2; }

.bm-input { padding: 6px 10px; border: 1px solid #cbd5e1; border-radius: 6px; font-size: 12px; background: white; box-sizing: border-box; font-family: inherit; }
.bm-input:focus { outline: none; border-color: #0ea5e9; }
.bm-input.search { min-width: 180px; }
textarea.bm-input { font-family: ui-monospace, monospace; resize: vertical; }

.bm-empty { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 50px 20px; color: #64748b; gap: 8px; }
.bm-empty p { margin: 0; font-size: 13px; }

/* 表格 */
.bm-table { width: 100%; border-collapse: collapse; background: white; border-radius: 8px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.05); }
.bm-table th, .bm-table td { padding: 8px 10px; text-align: left; font-size: 12px; border-bottom: 1px solid #f1f5f9; }
.bm-table th { background: #f8fafc; color: #475569; font-weight: 500; font-size: 11px; }
.bm-table tbody tr:hover { background: #f8fafc; }
.bm-table tbody tr.sel { background: #fef3c7; }
.bm-table.compact th, .bm-table.compact td { padding: 5px 8px; font-size: 11px; }
.paths { max-width: 240px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.ck { width: 32px; text-align: center; }
.ops { white-space: nowrap; display: flex; gap: 3px; }
.sid { background: #f1f5f9; padding: 1px 5px; border-radius: 3px; font-size: 11px; }
.muted { color: #94a3b8; }
.err-line { font-size: 10px; color: #dc2626; margin-top: 2px; display: flex; align-items: center; gap: 3px; }

.chip { display: inline-block; padding: 1px 6px; border-radius: 999px; font-size: 10px; background: #e0f2fe; color: #0369a1; margin-right: 3px; }
.chip.src { background: #dcfce7; color: #166534; }
.chip.dst { background: #fef3c7; color: #854d0e; }
.chip.st-active, .chip.st-completed { background: #dcfce7; color: #166534; }
.chip.st-running { background: #fef9c3; color: #854d0e; }
.chip.st-idle { background: #f1f5f9; color: #475569; }
.chip.st-failed, .chip.st-error { background: #fee2e2; color: #991b1b; }
.chip.st-uninitialized { background: #fef9c3; color: #854d0e; }
.chip.ft-dir { background: #fef9c3; color: #854d0e; }
.chip.ft-file { background: #e0f2fe; color: #0369a1; }

/* 仓库卡片 */
.repo-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 12px; }
.repo-card { background: white; border: 1px solid #e2e8f0; border-radius: 10px; padding: 12px; display: flex; flex-direction: column; gap: 8px; }
.repo-head { display: flex; align-items: center; gap: 8px; }
.repo-icon { width: 30px; height: 30px; border-radius: 6px; display: flex; align-items: center; justify-content: center; color: white; flex-shrink: 0; }
.repo-icon.t-local, .repo-icon.t-sftp { background: #64748b; }
.repo-icon.t-webdav, .repo-icon.t-rest { background: #8b5cf6; }
.repo-icon.t-rclone, .repo-icon.t-s3, .repo-icon.t-b2 { background: #0ea5e9; }
.repo-name { flex: 1; min-width: 0; }
.repo-name h3 { margin: 0; font-size: 13px; }
.repo-sub { font-size: 10px; color: #94a3b8; margin-top: 2px; }
.repo-url { font-family: ui-monospace, monospace; font-size: 11px; color: #475569; background: #f1f5f9; padding: 4px 6px; border-radius: 4px; word-break: break-all; }
.repo-meta { font-size: 11px; color: #64748b; display: flex; align-items: center; gap: 4px; }
.repo-ops { display: flex; flex-wrap: wrap; gap: 4px; }

/* 仓库卡片内的任务区 */
.task-divider { height: 1px; background: #e2e8f0; margin: 4px 0; }
.task-section { display: flex; flex-direction: column; gap: 6px; }
.task-section-head { display: flex; justify-content: space-between; align-items: center; font-size: 11px; color: #475569; font-weight: 600; }
.task-empty { font-size: 11px; color: #94a3b8; padding: 6px; text-align: center; background: #f8fafc; border-radius: 4px; }
.task-row { display: flex; gap: 8px; padding: 8px; background: #f8fafc; border-radius: 6px; border: 1px solid #e2e8f0; }
.task-row-info { flex: 1; min-width: 0; }
.task-row-name { font-size: 12px; font-weight: 600; color: #0f172a; display: flex; align-items: center; gap: 6px; margin-bottom: 3px; }
.task-row-meta { font-size: 10px; color: #64748b; display: flex; flex-wrap: wrap; gap: 4px; align-items: center; }
.task-row-meta code { background: white; padding: 1px 4px; border-radius: 3px; font-size: 10px; }
.task-row-ops { display: flex; flex-direction: column; gap: 3px; align-items: flex-end; }

/* 设置 */
.settings-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(380px, 1fr)); gap: 14px; }
.settings-card { background: white; border: 1px solid #e2e8f0; border-radius: 10px; padding: 16px; }
.settings-card h3 { margin: 0 0 10px; font-size: 14px; display: flex; align-items: center; gap: 6px; padding-bottom: 8px; border-bottom: 1px solid #f1f5f9; }
.card-desc { font-size: 11px; color: #64748b; margin: 0 0 12px; line-height: 1.5; }
.card-actions { display: flex; gap: 6px; margin-top: 10px; }
.bm-row { margin-bottom: 10px; }
.bm-row label { display: block; font-size: 11px; color: #475569; margin-bottom: 4px; font-weight: 500; }
.bm-row.checkbox { display: flex; align-items: flex-start; gap: 6px; }
.bm-row.checkbox label { margin: 0; font-weight: 400; font-size: 12px; line-height: 1.4; }
.hint { display: block; font-size: 10px; color: #94a3b8; margin-top: 3px; }
.hint code, .bm-row label code, .bm-alert code, .info-row code { background: #f1f5f9; padding: 1px 4px; border-radius: 3px; font-family: ui-monospace, monospace; font-size: 10px; }
.cred-section { border: 1px solid #e2e8f0; border-radius: 8px; padding: 10px; margin-bottom: 10px; background: #f8fafc; }
.cred-title { font-size: 11px; color: #475569; font-weight: 600; margin-bottom: 6px; }
.ret-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(90px, 1fr)); gap: 6px; }
.ret-grid label { display: flex; flex-direction: column; gap: 3px; font-size: 11px; color: #475569; }
.ret-grid input { padding: 5px 6px; border: 1px solid #cbd5e1; border-radius: 4px; font-size: 12px; }

.info-list { display: flex; flex-direction: column; gap: 6px; }
.info-row { display: flex; justify-content: space-between; gap: 10px; font-size: 11px; align-items: center; }
.info-row span:first-child { color: #64748b; }
.info-row code { font-size: 10px; word-break: break-all; text-align: right; }
.info-tip { display: flex; gap: 6px; align-items: flex-start; font-size: 11px; color: #64748b; background: #f8fafc; padding: 8px; border-radius: 6px; margin-top: 10px; line-height: 1.5; }
.test-result { margin-top: 10px; }
.test-result pre { background: #0f172a; color: #e2e8f0; padding: 8px; border-radius: 6px; font-size: 11px; max-height: 200px; overflow: auto; white-space: pre-wrap; word-break: break-all; margin: 0; }
.test-result.ok pre { background: #052e16; }
.test-result.err pre { background: #450a0a; }

.bm-alert { display: flex; gap: 8px; padding: 8px 10px; border-radius: 6px; font-size: 11px; margin-bottom: 10px; align-items: flex-start; }
.bm-alert.warn { background: #fef3c7; color: #92400e; }
.bm-alert.info { background: #e0f2fe; color: #075985; }

/* 搜索栏 */
.find-box { display: flex; align-items: center; gap: 6px; margin-top: 12px; padding: 8px; background: white; border-radius: 8px; border: 1px solid #e2e8f0; }
.find-box .bm-input { flex: 1; }

/* Modals */
.bm-overlay { position: fixed; inset: 0; background: rgba(15,23,42,0.5); backdrop-filter: blur(4px); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.bm-modal { background: white; border-radius: 12px; width: 92%; max-width: 540px; max-height: 90vh; display: flex; flex-direction: column; box-shadow: 0 20px 60px rgba(0,0,0,0.3); }
.bm-modal.lg { max-width: 860px; }
.bm-modal-head { display: flex; justify-content: space-between; align-items: center; padding: 14px 18px; border-bottom: 1px solid #e2e8f0; }
.bm-modal-head h3 { margin: 0; font-size: 15px; }
.bm-x { background: none; border: none; cursor: pointer; color: #64748b; padding: 4px; border-radius: 4px; }
.bm-x:hover { background: #f1f5f9; color: #0f172a; }
.bm-modal-body { padding: 18px; overflow: auto; }
.bm-modal-foot { padding: 10px 18px; border-top: 1px solid #e2e8f0; display: flex; justify-content: flex-end; gap: 6px; }
.bm-modal-foot.inline { padding: 14px 0 0; border: none; }
.files-toolbar { display: flex; align-items: center; gap: 8px; margin-bottom: 10px; }
.detail-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.detail-card { border: 1px solid #e2e8f0; border-radius: 8px; padding: 10px; }
.detail-card h4 { margin: 0 0 8px; font-size: 12px; }
.bm-log { background: #0f172a; color: #e2e8f0; padding: 10px; border-radius: 6px; font-family: ui-monospace, monospace; font-size: 11px; line-height: 1.5; max-height: 50vh; overflow: auto; white-space: pre-wrap; word-break: break-word; margin: 8px 0; }
.bm-toast { position: fixed; bottom: 20px; right: 20px; padding: 10px 16px; border-radius: 8px; color: white; font-size: 12px; max-width: 340px; box-shadow: 0 8px 24px rgba(0,0,0,0.2); z-index: 2000; white-space: pre-wrap; }
.bm-toast.t-success { background: #16a34a; }
.bm-toast.t-error { background: #dc2626; }
.bm-toast.t-info { background: #0ea5e9; }
.fade-enter-active, .fade-leave-active { transition: opacity .3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
.spinning { animation: spin 1s linear infinite; }
@keyframes spin { from { transform: rotate(0); } to { transform: rotate(360deg); } }

@media (max-width: 720px) {
  .detail-grid, .settings-grid { grid-template-columns: 1fr; }
  .bm-tabs { width: 100%; overflow-x: auto; }
}
</style>
