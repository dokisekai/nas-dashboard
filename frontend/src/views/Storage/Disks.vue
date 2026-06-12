<template>
  <div class="space-y-6">
    <!-- 磁盘列表 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">磁盘列表</h3>
          <p class="text-gray-500 text-sm">管理和配置存储设备</p>
        </div>
        <button
          @click="loadDisks"
          :disabled="loading.disks"
          class="flex items-center gap-2 px-4 py-2 bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-400 rounded-xl transition-all border border-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg
            class="w-4 h-4"
            :class="{ 'animate-spin': loading.disks }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          刷新
        </button>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading.disks && disks.length === 0" class="flex items-center justify-center py-12">
        <div class="flex items-center gap-3 text-gray-400">
          <svg class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>加载磁盘列表中...</span>
        </div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="error.disks" class="flex items-center justify-center py-12">
        <div class="text-center">
          <div class="w-12 h-12 mx-auto mb-4 bg-red-500/10 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <p class="text-red-400 mb-2">{{ error.disks }}</p>
          <button @click="loadDisks" class="text-indigo-400 hover:text-indigo-300 text-sm">重试</button>
        </div>
      </div>

      <!-- 磁盘列表 -->
      <div v-else-if="disks.length > 0" class="space-y-3">
        <div
          v-for="disk in disks"
          :key="disk.name"
          class="bg-gray-900/50 rounded-xl p-5 border border-gray-800 hover:border-gray-700 transition-all"
        >
          <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-3 mb-3">
                <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center">
                  <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                  </svg>
                </div>
                <div class="min-w-0">
                  <div class="flex items-center gap-2">
                    <span class="text-white font-medium truncate">{{ disk.name }}</span>
                    <span
                      class="px-2 py-0.5 text-xs rounded-full flex-shrink-0"
                      :class="disk.mounted ? 'bg-green-500/10 text-green-400' : 'bg-gray-500/10 text-gray-400'"
                    >
                      {{ disk.mounted ? '已挂载' : '未挂载' }}
                    </span>
                  </div>
                  <div class="flex items-center gap-4 mt-1 text-sm text-gray-400">
                    <span>{{ disk.type }}</span>
                    <span v-if="disk.label" class="truncate">{{ disk.label }}</span>
                  </div>
                </div>
              </div>

              <!-- 容量信息 -->
              <div class="space-y-2">
                <div class="flex items-center justify-between text-sm">
                  <span class="text-gray-400">容量: {{ formatBytes(disk.size) }}</span>
                  <span v-if="disk.used !== undefined" class="text-gray-300">
                    {{ formatBytes(disk.used) }} / {{ formatBytes(disk.size) }}
                  </span>
                </div>
                <!-- 使用率进度条 -->
                <div v-if="disk.usagePercent !== undefined" class="h-2 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full rounded-full transition-all duration-500"
                    :class="getUsageColorClass(disk.usagePercent)"
                    :style="{ width: `${disk.usagePercent}%` }"
                  />
                </div>
                <div v-if="disk.usagePercent !== undefined" class="flex items-center justify-between text-xs">
                  <span class="text-gray-500">使用率</span>
                  <span :class="getUsageTextColorClass(disk.usagePercent)">{{ disk.usagePercent }}%</span>
                </div>
              </div>

              <!-- 挂载点信息 -->
              <div v-if="disk.mountPoint" class="mt-3 flex items-center gap-2 text-sm text-gray-400">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                </svg>
                <span class="font-mono">{{ disk.mountPoint }}</span>
              </div>
            </div>

            <div class="flex flex-wrap gap-2">
              <button
                v-if="!disk.mounted"
                @click="openMountDialog(disk)"
                :disabled="loading.mount"
                class="flex items-center gap-1.5 px-3 py-2 bg-green-500/10 hover:bg-green-500/20 text-green-400 rounded-lg text-sm transition-all border border-green-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                </svg>
                挂载
              </button>
              <button
                v-if="disk.mounted"
                @click="openUnmountDialog(disk)"
                :disabled="loading.unmount"
                class="flex items-center gap-1.5 px-3 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg text-sm transition-all border border-red-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                卸载
              </button>
              <button
                @click="openFormatDialog(disk)"
                class="flex items-center gap-1.5 px-3 py-2 bg-orange-500/10 hover:bg-orange-500/20 text-orange-400 rounded-lg text-sm transition-all border border-orange-500/20"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                格式化
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="flex flex-col items-center justify-center py-12">
        <div class="w-16 h-16 mb-4 bg-gray-700/50 rounded-full flex items-center justify-center">
          <svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
          </svg>
        </div>
        <p class="text-gray-500 mb-4">暂无可用磁盘</p>
        <button @click="loadDisks" class="text-indigo-400 hover:text-indigo-300 text-sm">刷新列表</button>
      </div>
    </div>

    <!-- SMB 共享 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
      <div class="flex items-center justify-between mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">SMB 共享</h3>
          <p class="text-gray-500 text-sm">Windows 文件共享配置</p>
        </div>
        <button
          @click="openCreateShareDialog"
          class="flex items-center gap-2 px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-xl transition-all shadow-lg shadow-indigo-500/25"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          创建共享
        </button>
      </div>

      <!-- 加载状态 -->
      <div v-if="loading.shares && shares.length === 0" class="flex items-center justify-center py-12">
        <div class="flex items-center gap-3 text-gray-400">
          <svg class="w-5 h-5 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <span>加载共享列表中...</span>
        </div>
      </div>

      <!-- 共享列表 -->
      <div v-else-if="shares.length > 0" class="overflow-hidden rounded-xl border border-gray-800">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-900/50">
              <tr class="text-gray-400 text-sm">
                <th class="text-left py-4 px-4 font-medium">名称</th>
                <th class="text-left py-4 px-4 font-medium">路径</th>
                <th class="text-left py-4 px-4 font-medium hidden md:table-cell">描述</th>
                <th class="text-left py-4 px-4 font-medium">状态</th>
                <th class="text-right py-4 px-4 font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="share in shares"
                :key="share.name"
                class="border-t border-gray-800 hover:bg-gray-900/30 transition-colors"
              >
                <td class="py-4 px-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 bg-blue-500/10 rounded-lg flex items-center justify-center">
                      <svg class="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2" />
                      </svg>
                    </div>
                    <span class="text-white font-medium">{{ share.name }}</span>
                  </div>
                </td>
                <td class="py-4 px-4 text-gray-400 font-mono text-sm">{{ share.path }}</td>
                <td class="py-4 px-4 text-gray-500 text-sm hidden md:table-cell">{{ share.description || '-' }}</td>
                <td class="py-4 px-4">
                  <span class="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs rounded-full bg-green-500/10 text-green-400">
                    <div class="w-1.5 h-1.5 bg-green-400 rounded-full animate-pulse"></div>
                    启用
                  </span>
                </td>
                <td class="py-4 px-4 text-right">
                  <button
                    @click="openEditShareDialog(share)"
                    class="text-indigo-400 hover:text-indigo-300 transition-colors mr-3"
                  >
                    编辑
                  </button>
                  <button
                    @click="openDeleteShareDialog(share)"
                    class="text-red-400 hover:text-red-300 transition-colors"
                  >
                    删除
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="flex flex-col items-center justify-center py-12">
        <div class="w-16 h-16 mb-4 bg-gray-700/50 rounded-full flex items-center justify-center">
          <svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7v8a2 2 0 002 2h6M8 7V5a2 2 0 012-2h4.586a1 1 0 01.707.293l4.414 4.414a1 1 0 01.293.707V15a2 2 0 01-2 2h-2M8 7H6a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2v-2" />
          </svg>
        </div>
        <p class="text-gray-500 mb-4">暂无 SMB 共享</p>
        <button @click="openCreateShareDialog" class="text-indigo-400 hover:text-indigo-300 text-sm">创建第一个共享</button>
      </div>
    </div>

    <!-- 挂载对话框 -->
    <BaseModal v-model:show="dialogs.mount" title="挂载磁盘" @confirm="handleMount" @cancel="closeMountDialog">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-gray-300 mb-2">设备</label>
          <div class="px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white">
            {{ selectedDisk?.name }}
          </div>
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">挂载点</label>
          <input
            v-model="mountForm.mountPoint"
            type="text"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="/mnt/data"
          />
          <p class="mt-2 text-xs text-gray-500">默认挂载点: /mnt/{{ selectedDisk?.name?.split('/').pop() }}</p>
        </div>
        <div v-if="selectedDisk?.type">
          <label class="block text-sm text-gray-300 mb-2">文件系统类型</label>
          <div class="px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white">
            {{ selectedDisk.type }}
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- 卸载确认对话框 -->
    <ConfirmModal
      v-model:show="dialogs.unmount"
      title="确认卸载"
      message="确定要卸载此磁盘吗？卸载后无法访问磁盘上的数据。"
      :detail="`设备: ${selectedDisk?.name}\n挂载点: ${selectedDisk?.mountPoint}`"
      @confirm="handleUnmount"
    />

    <!-- 格式化确认对话框 -->
    <ConfirmModal
      v-model:show="dialogs.format"
      title="格式化磁盘"
      message="警告：格式化将删除磁盘上的所有数据！此操作不可逆。"
      type="danger"
      :detail="`设备: ${selectedDisk?.name}\n容量: ${selectedDisk?.size ? formatBytes(selectedDisk.size) : '未知'}`"
      @confirm="handleFormat"
    >
      <div class="space-y-4 mt-4">
        <div>
          <label class="block text-sm text-gray-300 mb-2">文件系统</label>
          <select
            v-model="formatForm.fsType"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            <option value="ext4">ext4 (推荐)</option>
            <option value="xfs">XFS</option>
            <option value="btrfs">BTRFS</option>
            <option value="ntfs">NTFS</option>
            <option value="exfat">exFAT</option>
          </select>
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">
            <input
              v-model="formatForm.confirm"
              type="checkbox"
              class="mr-2 rounded bg-gray-900 border-gray-700 text-red-500 focus:ring-red-500"
            />
            我了解此操作将删除所有数据
          </label>
        </div>
      </div>
    </ConfirmModal>

    <!-- 创建共享对话框 -->
    <BaseModal v-model:show="dialogs.createShare" title="创建 SMB 共享" @confirm="handleCreateShare" @cancel="closeCreateShareDialog">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-gray-300 mb-2">共享名称</label>
          <input
            v-model="shareForm.name"
            type="text"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="例如: public"
          />
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">共享路径</label>
          <input
            v-model="shareForm.path"
            type="text"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="/mnt/data"
          />
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">描述</label>
          <input
            v-model="shareForm.description"
            type="text"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="共享描述（可选）"
          />
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">访问权限</label>
          <div class="space-y-2">
            <label class="flex items-center gap-2">
              <input
                v-model="shareForm.readOnly"
                type="checkbox"
                class="rounded bg-gray-900 border-gray-700 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-gray-400">只读访问</span>
            </label>
            <label class="flex items-center gap-2">
              <input
                v-model="shareForm.guest"
                type="checkbox"
                class="rounded bg-gray-900 border-gray-700 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-gray-400">允许访客访问（无需密码）</span>
            </label>
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- 编辑共享对话框 -->
    <BaseModal v-model:show="dialogs.editShare" title="编辑 SMB 共享" @confirm="handleEditShare" @cancel="closeEditShareDialog">
      <div class="space-y-4">
        <div>
          <label class="block text-sm text-gray-300 mb-2">共享名称</label>
          <div class="px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-gray-500">
            {{ selectedShare?.name }}
          </div>
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">共享路径</label>
          <input
            v-model="shareForm.path"
            type="text"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          />
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">描述</label>
          <input
            v-model="shareForm.description"
            type="text"
            class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            placeholder="共享描述（可选）"
          />
        </div>
        <div>
          <label class="block text-sm text-gray-300 mb-2">访问权限</label>
          <div class="space-y-2">
            <label class="flex items-center gap-2">
              <input
                v-model="shareForm.readOnly"
                type="checkbox"
                class="rounded bg-gray-900 border-gray-700 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-gray-400">只读访问</span>
            </label>
            <label class="flex items-center gap-2">
              <input
                v-model="shareForm.guest"
                type="checkbox"
                class="rounded bg-gray-900 border-gray-700 text-indigo-500 focus:ring-indigo-500"
              />
              <span class="text-gray-400">允许访客访问</span>
            </label>
          </div>
        </div>
      </div>
    </BaseModal>

    <!-- 删除共享确认对话框 -->
    <ConfirmModal
      v-model:show="dialogs.deleteShare"
      title="删除 SMB 共享"
      message="确定要删除此 SMB 共享吗？"
      :detail="`共享: ${selectedShare?.name}\n路径: ${selectedShare?.path}`"
      @confirm="handleDeleteShare"
    />

    <!-- Toast 通知 -->
    <Toast
      v-model:show="toast.show"
      :message="toast.message"
      :type="toast.type"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { storageApi } from '../../api'
import BaseModal from '../../components/Common/BaseModal.vue'
import ConfirmModal from '../../components/Common/ConfirmModal.vue'
import Toast from '../../components/Common/Toast.vue'

// 磁盘接口
interface Disk {
  name: string
  type: string
  size: number
  used?: number
  usagePercent?: number
  mountPoint?: string
  mounted: boolean
  label?: string
}

// 共享接口
interface Share {
  name: string
  path: string
  description?: string
  readOnly?: boolean
  guest?: boolean
}

// 状态管理
const disks = ref<Disk[]>([])
const shares = ref<Share[]>([])
const selectedDisk = ref<Disk | null>(null)
const selectedShare = ref<Share | null>(null)

// 加载和错误状态
const loading = reactive({
  disks: false,
  shares: false,
  mount: false,
  unmount: false,
  format: false
})

const error = reactive({
  disks: '',
  shares: ''
})

// 对话框状态
const dialogs = reactive({
  mount: false,
  unmount: false,
  format: false,
  createShare: false,
  editShare: false,
  deleteShare: false
})

// 表单数据
const mountForm = reactive({
  mountPoint: ''
})

const formatForm = reactive({
  fsType: 'ext4',
  confirm: false
})

const shareForm = reactive({
  name: '',
  path: '',
  description: '',
  readOnly: false,
  guest: false
})

// Toast 通知
const toast = reactive({
  show: false,
  message: '',
  type: 'info' as 'info' | 'success' | 'error' | 'warning'
})

// 显示 Toast
const showToast = (message: string, type: 'info' | 'success' | 'error' | 'warning' = 'info') => {
  toast.message = message
  toast.type = type
  toast.show = true
  setTimeout(() => {
    toast.show = false
  }, 3000)
}

// 格式化字节
const formatBytes = (bytes: number): string => {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

// 获取使用率颜色类
const getUsageColorClass = (percent: number): string => {
  if (percent >= 90) return 'bg-red-500'
  if (percent >= 75) return 'bg-orange-500'
  if (percent >= 50) return 'bg-yellow-500'
  return 'bg-green-500'
}

// 获取使用率文本颜色类
const getUsageTextColorClass = (percent: number): string => {
  if (percent >= 90) return 'text-red-400'
  if (percent >= 75) return 'text-orange-400'
  if (percent >= 50) return 'text-yellow-400'
  return 'text-green-400'
}

// 加载磁盘列表
const loadDisks = async () => {
  loading.disks = true
  error.disks = ''
  try {
    const response = await storageApi.getDisks()
    disks.value = (response.data?.disks || []).map((disk: any) => ({
      ...disk,
      usagePercent: disk.used && disk.size ? Math.round((disk.used / disk.size) * 100) : undefined
    }))
  } catch (err: any) {
    console.error('获取磁盘列表失败:', err)
    error.disks = err.response?.data?.message || err.message || '获取磁盘列表失败'
    showToast(error.disks, 'error')
  } finally {
    loading.disks = false
  }
}

// 加载共享列表
const loadShares = async () => {
  loading.shares = true
  error.shares = ''
  try {
    const response = await storageApi.getSMBShares()
    shares.value = response.data?.shares || []
  } catch (err: any) {
    console.error('获取共享列表失败:', err)
    error.shares = err.response?.data?.message || err.message || '获取共享列表失败'
    showToast(error.shares, 'error')
  } finally {
    loading.shares = false
  }
}

// 挂载相关
const openMountDialog = (disk: Disk) => {
  selectedDisk.value = disk
  mountForm.mountPoint = `/mnt/${disk.name.split('/').pop()}`
  dialogs.mount = true
}

const closeMountDialog = () => {
  dialogs.mount = false
  selectedDisk.value = null
  mountForm.mountPoint = ''
}

const handleMount = async () => {
  if (!selectedDisk.value || !mountForm.mountPoint) return

  loading.mount = true
  try {
    await storageApi.mount(selectedDisk.value.name, mountForm.mountPoint)
    // 更新本地状态
    const disk = disks.value.find(d => d.name === selectedDisk.value?.name)
    if (disk) {
      disk.mounted = true
      disk.mountPoint = mountForm.mountPoint
    }
    showToast(`成功挂载 ${selectedDisk.value.name} 到 ${mountForm.mountPoint}`, 'success')
    closeMountDialog()
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || '挂载失败'
    showToast(msg, 'error')
  } finally {
    loading.mount = false
  }
}

// 卸载相关
const openUnmountDialog = (disk: Disk) => {
  selectedDisk.value = disk
  dialogs.unmount = true
}

const handleUnmount = async () => {
  if (!selectedDisk.value?.mountPoint) return

  loading.unmount = true
  try {
    await storageApi.umount(selectedDisk.value.mountPoint)
    // 更新本地状态
    const disk = disks.value.find(d => d.name === selectedDisk.value?.name)
    if (disk) {
      disk.mounted = false
      disk.mountPoint = ''
    }
    showToast(`成功卸载 ${selectedDisk.value.name}`, 'success')
    dialogs.unmount = false
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || '卸载失败'
    showToast(msg, 'error')
  } finally {
    loading.unmount = false
  }
}

// 格式化相关
const openFormatDialog = (disk: Disk) => {
  selectedDisk.value = disk
  formatForm.fsType = 'ext4'
  formatForm.confirm = false
  dialogs.format = true
}

const handleFormat = async () => {
  if (!selectedDisk.value || !formatForm.confirm) return

  loading.format = true
  try {
    await storageApi.formatDisk(selectedDisk.value.name, formatForm.fsType)
    showToast(`成功格式化 ${selectedDisk.value.name}`, 'success')
    dialogs.format = false
    // 重新加载磁盘列表
    await loadDisks()
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || '格式化失败'
    showToast(msg, 'error')
  } finally {
    loading.format = false
  }
}

// 共享管理 - 创建
const openCreateShareDialog = () => {
  shareForm.name = ''
  shareForm.path = ''
  shareForm.description = ''
  shareForm.readOnly = false
  shareForm.guest = false
  dialogs.createShare = true
}

const closeCreateShareDialog = () => {
  dialogs.createShare = false
}

const handleCreateShare = async () => {
  if (!shareForm.name || !shareForm.path) {
    showToast('请填写共享名称和路径', 'warning')
    return
  }

  try {
    await storageApi.createSMBShare(
      shareForm.name,
      shareForm.path,
      shareForm.description,
      shareForm.readOnly,
      shareForm.guest
    )
    shares.value.push({
      name: shareForm.name,
      path: shareForm.path,
      description: shareForm.description,
      readOnly: shareForm.readOnly,
      guest: shareForm.guest
    })
    showToast(`成功创建共享 ${shareForm.name}`, 'success')
    closeCreateShareDialog()
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || '创建共享失败'
    showToast(msg, 'error')
  }
}

// 共享管理 - 编辑
const openEditShareDialog = (share: Share) => {
  selectedShare.value = share
  shareForm.name = share.name
  shareForm.path = share.path
  shareForm.description = share.description || ''
  shareForm.readOnly = share.readOnly || false
  shareForm.guest = share.guest || false
  dialogs.editShare = true
}

const closeEditShareDialog = () => {
  dialogs.editShare = false
  selectedShare.value = null
}

const handleEditShare = async () => {
  if (!selectedShare.value) return

  try {
    await storageApi.updateSMBShare(
      selectedShare.value.name,
      shareForm.path,
      shareForm.description,
      shareForm.readOnly,
      shareForm.guest
    )
    // 更新本地状态
    const index = shares.value.findIndex(s => s.name === selectedShare.value?.name)
    if (index !== -1) {
      shares.value[index] = {
        name: selectedShare.value.name,
        path: shareForm.path,
        description: shareForm.description,
        readOnly: shareForm.readOnly,
        guest: shareForm.guest
      }
    }
    showToast(`成功更新共享 ${selectedShare.value.name}`, 'success')
    closeEditShareDialog()
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || '更新共享失败'
    showToast(msg, 'error')
  }
}

// 共享管理 - 删除
const openDeleteShareDialog = (share: Share) => {
  selectedShare.value = share
  dialogs.deleteShare = true
}

const handleDeleteShare = async () => {
  if (!selectedShare.value) return

  try {
    await storageApi.deleteSMBShare(selectedShare.value.name)
    shares.value = shares.value.filter(s => s.name !== selectedShare.value?.name)
    showToast(`成功删除共享 ${selectedShare.value.name}`, 'success')
    dialogs.deleteShare = false
    selectedShare.value = null
  } catch (err: any) {
    const msg = err.response?.data?.message || err.message || '删除共享失败'
    showToast(msg, 'error')
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadDisks()
  loadShares()
})
</script>
