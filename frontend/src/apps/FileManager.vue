<template>
  <div class="file-manager">
    <!-- 工具栏 -->
    <div class="fm-toolbar">
      <div class="toolbar-left">
        <!-- 导航按钮 -->
        <div class="nav-buttons">
          <button class="nav-btn" @click="goBack" :disabled="!history.canGoBack">
            <ArrowLeftIcon class="w-4 h-4" />
          </button>
          <button class="nav-btn" @click="goForward" :disabled="!history.canGoForward">
            <ArrowRightIcon class="w-4 h-4" />
          </button>
          <button class="nav-btn" @click="goUp" :disabled="currentPath === '/'">
            <ArrowUpIcon class="w-4 h-4" />
          </button>
          <button class="nav-btn" @click="refreshFiles">
            <ArrowPathIcon class="w-4 h-4" />
          </button>
        </div>

        <!-- 路径导航 -->
        <div class="path-navigation">
          <div class="path-breadcrumbs">
            <span
              v-for="(segment, index) in pathSegments"
              :key="index"
              class="path-segment"
              @click="navigateToBreadcrumb(index)"
            >
              {{ segment }}
              <ChevronRightIcon v-if="index < pathSegments.length - 1" class="w-4 h-4" />
            </span>
          </div>
        </div>
      </div>

      <div class="toolbar-right">
        <!-- 搜索框 -->
        <div class="search-box">
          <MagnifyingGlassIcon class="w-4 h-4" />
          <input
            v-model="searchQuery"
            type="text"
            placeholder="搜索文件..."
            @input="filterFiles"
          />
        </div>

        <!-- 视图切换 -->
        <div class="view-toggle">
          <button
            :class="['view-btn', { active: viewMode === 'grid' }]"
            @click="viewMode = 'grid'"
            title="网格视图"
          >
            <Squares2X2Icon class="w-4 h-4" />
          </button>
          <button
            :class="['view-btn', { active: viewMode === 'list' }]"
            @click="viewMode = 'list'"
            title="列表视图"
          >
            <ListBulletIcon class="w-4 h-4" />
          </button>
        </div>

        <!-- 操作按钮 -->
        <button class="action-btn primary" @click="showUploadModal = true">
          <ArrowUpTrayIcon class="w-4 h-4" />
          上传
        </button>
        <button class="action-btn" @click="createFolder">
          <FolderPlusIcon class="w-4 h-4" />
          新建文件夹
        </button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="fm-content">
      <!-- 侧边栏 -->
      <div class="fm-sidebar">
        <div class="sidebar-section">
          <h3>位置</h3>
          <div class="sidebar-items">
            <div
              class="sidebar-item"
              :class="{ active: currentLocation === 'home' }"
              @click="navigateToHome"
            >
              <HomeIcon class="w-4 h-4" />
              <span>主目录</span>
            </div>
            <div
              class="sidebar-item"
              :class="{ active: currentLocation === 'shared' }"
              @click="navigateToShared"
            >
              <FolderIcon class="w-4 h-4" />
              <span>共享文件夹</span>
            </div>
            <div
              class="sidebar-item"
              :class="{ active: currentLocation === 'external' }"
              @click="navigateToExternal"
            >
              <ExternalLinkIcon class="w-4 h-4" />
              <span>外接设备</span>
            </div>
          </div>
        </div>

        <div class="sidebar-section">
          <h3>收藏</h3>
          <div class="sidebar-items">
            <div
              v-for="bookmark in bookmarks"
              :key="bookmark.id"
              class="sidebar-item"
              @click="navigateTo(bookmark.path)"
            >
              <StarIcon class="w-4 h-4" />
              <span>{{ bookmark.name }}</span>
            </div>
          </div>
        </div>

        <div class="sidebar-section">
          <h3>存储空间</h3>
          <div class="storage-info">
            <div class="storage-item">
              <div class="storage-label">系统盘</div>
              <div class="storage-bar">
                <div class="storage-fill" :style="{ width: '45%' }"></div>
              </div>
              <div class="storage-text">450GB / 1TB</div>
            </div>
            <div class="storage-item">
              <div class="storage-label">存储池1</div>
              <div class="storage-bar">
                <div class="storage-fill" :style="{ width: '78%' }"></div>
              </div>
              <div class="storage-text">7.8TB / 10TB</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 文件列表区 -->
      <div class="fm-main">
        <!-- 状态信息 -->
        <div class="fm-status">
          <span>当前路径: {{ currentPath }}</span>
          <span>{{ filteredFiles.length }} 个项目</span>
        </div>

        <!-- 文件网格/列表 -->
        <div v-if="viewMode === 'grid'" class="files-grid">
          <div
            v-for="file in filteredFiles"
            :key="file.id"
            class="file-card"
            :class="{ selected: selectedFiles.includes(file.id) }"
            @click="selectFile(file)"
            @dblclick="openFile(file)"
            @contextmenu.prevent="showContextMenu($event, file)"
          >
            <div class="file-icon">
              <FolderIcon v-if="file.type === 'folder'" class="w-12 h-12" />
              <DocumentIcon v-else class="w-12 h-12" />
            </div>
            <div class="file-info">
              <div class="file-name">{{ file.name }}</div>
              <div class="file-meta">{{ formatFileSize(file.size) }}</div>
            </div>
          </div>
        </div>

        <div v-else class="files-list">
          <table class="files-table">
            <thead>
              <tr>
                <th @click="sortBy('name')">
                  名称
                  <ChevronUpDownIcon class="w-4 h-4 sort-icon" />
                </th>
                <th @click="sortBy('size')">
                  大小
                  <ChevronUpDownIcon class="w-4 h-4 sort-icon" />
                </th>
                <th @click="sortBy('modified')">
                  修改时间
                  <ChevronUpDownIcon class="w-4 h-4 sort-icon" />
                </th>
                <th>权限</th>
                <th>所有者</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="file in filteredFiles"
                :key="file.id"
                :class="{ selected: selectedFiles.includes(file.id) }"
                @click="selectFile(file)"
                @dblclick="openFile(file)"
                @contextmenu.prevent="showContextMenu($event, file)"
              >
                <td>
                  <div class="file-name-cell">
                    <FolderIcon v-if="file.type === 'folder'" class="w-5 h-5" />
                    <DocumentIcon v-else class="w-5 h-5" />
                    <span>{{ file.name }}</span>
                  </div>
                </td>
                <td>{{ file.type === 'folder' ? '-' : formatFileSize(file.size) }}</td>
                <td>{{ formatDate(file.modified) }}</td>
                <td>{{ file.permissions || 'rw-rw-r--' }}</td>
                <td>{{ file.owner || 'admin' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- 上传模态框 -->
    <Transition name="fade">
      <div v-if="showUploadModal" class="upload-modal" @click.self="showUploadModal = false">
        <div class="upload-content">
          <div class="upload-header">
            <h3>上传文件</h3>
            <button @click="showUploadModal = false">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="upload-area" @drop="handleFileDrop" @dragover.prevent @dragenter.prevent>
            <ArrowUpTrayIcon class="w-12 h-12" />
            <p>拖拽文件到这里上传</p>
            <p>或者</p>
            <button class="browse-btn">选择文件</button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- 右键菜单 -->
    <Transition name="fade">
      <div
        v-if="contextMenu.visible"
        class="context-menu"
        :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
        @click="hideContextMenu"
      >
        <div class="context-item" @click="downloadFile(contextMenu.file)">
          <ArrowDownTrayIcon class="w-4 h-4" />
          下载
        </div>
        <div class="context-item" @click="renameFile(contextMenu.file)">
          <PencilIcon class="w-4 h-4" />
          重命名
        </div>
        <div class="context-item" @click="copyFile(contextMenu.file)">
          <ClipboardIcon class="w-4 h-4" />
          复制
        </div>
        <div class="context-item" @click="moveFile(contextMenu.file)">
          <FolderOpenIcon class="w-4 h-4" />
          移动到...
        </div>
        <div class="context-divider"></div>
        <div class="context-item" @click="shareFile(contextMenu.file)">
          <ShareIcon class="w-4 h-4" />
          分享
        </div>
        <div class="context-item" @click="addToBookmarks(contextMenu.file)">
          <StarIcon class="w-4 h-4" />
          添加到收藏
        </div>
        <div class="context-divider"></div>
        <div class="context-item danger" @click="deleteFile(contextMenu.file)">
          <TrashIcon class="w-4 h-4" />
          删除
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import {
  ArrowLeftIcon,
  ArrowRightIcon,
  ArrowUpIcon,
  ArrowPathIcon,
  MagnifyingGlassIcon,
  Squares2X2Icon,
  ListBulletIcon,
  ArrowUpTrayIcon,
  FolderPlusIcon,
  HomeIcon,
  FolderIcon,
  ArrowTopRightOnSquareIcon,
  StarIcon,
  DocumentIcon,
  ChevronRightIcon,
  ChevronUpDownIcon,
  XMarkIcon,
  ArrowDownTrayIcon,
  PencilIcon,
  ClipboardIcon,
  FolderOpenIcon,
  ShareIcon,
  TrashIcon
} from '@heroicons/vue/24/outline'

interface File {
  id: string
  name: string
  type: 'file' | 'folder'
  size: number
  modified: Date
  permissions?: string
  owner?: string
}

interface History {
  canGoBack: boolean
  canGoForward: boolean
}

// 状态管理
const currentPath = ref('/')
const viewMode = ref<'grid' | 'list'>('grid')
const searchQuery = ref('')
const selectedFiles = ref<string[]>([])
const showUploadModal = ref(false)
const currentLocation = ref('home')

const history = ref<History>({
  canGoBack: false,
  canGoForward: false
})

const contextMenu = ref({
  visible: false,
  x: 0,
  y: 0,
  file: null as File | null
})

// 模拟文件数据
const files = ref<File[]>([
  {
    id: '1',
    name: 'Documents',
    type: 'folder',
    size: 0,
    modified: new Date('2024-06-10'),
    permissions: 'rwxr-xr-x',
    owner: 'admin'
  },
  {
    id: '2',
    name: 'Photos',
    type: 'folder',
    size: 0,
    modified: new Date('2024-06-09'),
    permissions: 'rwxr-xr-x',
    owner: 'admin'
  },
  {
    id: '3',
    name: 'Videos',
    type: 'folder',
    size: 0,
    modified: new Date('2024-06-08'),
    permissions: 'rwxr-xr-x',
    owner: 'admin'
  },
  {
    id: '4',
    name: 'report.pdf',
    type: 'file',
    size: 2048576,
    modified: new Date('2024-06-11'),
    permissions: 'rw-r--r--',
    owner: 'admin'
  },
  {
    id: '5',
    name: 'backup.zip',
    type: 'file',
    size: 1073741824,
    modified: new Date('2024-06-07'),
    permissions: 'rw-------',
    owner: 'admin'
  }
])

const bookmarks = ref([
  { id: '1', name: '工作文档', path: '/home/documents/work' },
  { id: '2', name: '照片备份', path: '/home/photos' },
  { id: '3', name: '项目文件', path: '/home/projects' }
])

// 计算属性
const pathSegments = computed(() => {
  return currentPath.value.split('/').filter(Boolean)
})

const filteredFiles = computed(() => {
  if (!searchQuery.value) return files.value
  const query = searchQuery.value.toLowerCase()
  return files.value.filter(file =>
    file.name.toLowerCase().includes(query)
  )
})

// 方法
const navigateTo = (path: string) => {
  currentPath.value = path
}

const navigateToBreadcrumb = (index?: number) => {
  if (index === undefined) {
    currentPath.value = '/'
    return
  }

  const segments = pathSegments.value.slice(0, index + 1)
  currentPath.value = '/' + segments.join('/')
}

const navigateToHome = () => {
  currentPath.value = '/home/admin'
  currentLocation.value = 'home'
}

const navigateToShared = () => {
  currentPath.value = '/shared'
  currentLocation.value = 'shared'
}

const navigateToExternal = () => {
  currentPath.value = '/external'
  currentLocation.value = 'external'
}

const goBack = () => {
  // 实现后退逻辑
  console.log('Go back')
}

const goForward = () => {
  // 实现前进逻辑
  console.log('Go forward')
}

const goUp = () => {
  const segments = pathSegments.value.slice(0, -1)
  if (segments.length === 0) {
    currentPath.value = '/'
  } else {
    currentPath.value = '/' + segments.join('/')
  }
}

const refreshFiles = () => {
  // 实现刷新逻辑
  console.log('Refresh files')
}

const selectFile = (file: File) => {
  const index = selectedFiles.value.indexOf(file.id)
  if (index > -1) {
    selectedFiles.value.splice(index, 1)
  } else {
    selectedFiles.value.push(file.id)
  }
}

const openFile = (file: File) => {
  if (file.type === 'folder') {
    currentPath.value = currentPath.value + '/' + file.name
  } else {
    // 打开文件
    console.log('Open file:', file.name)
  }
}

const createFolder = () => {
  const folderName = prompt('请输入文件夹名称:')
  if (folderName) {
    files.value.push({
      id: Date.now().toString(),
      name: folderName,
      type: 'folder',
      size: 0,
      modified: new Date()
    })
  }
}

const showContextMenu = (event: MouseEvent, file: File) => {
  contextMenu.value = {
    visible: true,
    x: event.clientX,
    y: event.clientY,
    file
  }
}

const hideContextMenu = () => {
  contextMenu.value.visible = false
}

const downloadFile = (file: File | null) => {
  if (file) console.log('Download:', file.name)
}

const renameFile = (file: File | null) => {
  if (file) {
    const newName = prompt('请输入新名称:', file.name)
    if (newName) {
      file.name = newName
    }
  }
}

const copyFile = (file: File | null) => {
  if (file) console.log('Copy:', file.name)
}

const moveFile = (file: File | null) => {
  if (file) console.log('Move:', file.name)
}

const shareFile = (file: File | null) => {
  if (file) console.log('Share:', file.name)
}

const addToBookmarks = (file: File | null) => {
  if (file) {
    bookmarks.value.push({
      id: Date.now().toString(),
      name: file.name,
      path: currentPath.value + '/' + file.name
    })
  }
}

const deleteFile = (file: File | null) => {
  if (file && confirm(`确定要删除 "${file.name}" 吗?`)) {
    const index = files.value.findIndex(f => f.id === file.id)
    if (index > -1) {
      files.value.splice(index, 1)
    }
  }
}

const handleFileDrop = (event: DragEvent) => {
  event.preventDefault()
  const droppedFiles = event.dataTransfer?.files
  if (droppedFiles) {
    console.log('Files dropped:', droppedFiles)
    // 处理文件上传
  }
}

const filterFiles = () => {
  // 搜索逻辑已在 filteredFiles 计算属性中处理
}

const sortBy = (field: string) => {
  // 实现排序逻辑
  console.log('Sort by:', field)
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

// 点击外部关闭右键菜单
const handleClickOutside = () => {
  hideContextMenu()
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.file-manager {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #f9fafb;
}

.fm-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  gap: 20px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 20px;
  flex: 1;
}

.nav-buttons {
  display: flex;
  gap: 4px;
}

.nav-btn {
  padding: 8px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.nav-btn:hover:not(:disabled) {
  background: #f3f4f6;
}

.nav-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.path-navigation {
  flex: 1;
  max-width: 600px;
}

.path-breadcrumbs {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.path-segment {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #6b7280;
  cursor: pointer;
  font-size: 14px;
}

.path-segment:hover {
  color: #3b82f6;
}

.toolbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
}

.search-box input {
  border: none;
  background: transparent;
  outline: none;
  font-size: 14px;
  min-width: 200px;
}

.view-toggle {
  display: flex;
  gap: 4px;
}

.view-btn {
  padding: 8px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 6px;
  cursor: pointer;
}

.view-btn.active {
  background: #eff6ff;
  border-color: #3b82f6;
  color: #3b82f6;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: 1px solid #e5e7eb;
  background: white;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f3f4f6;
}

.action-btn.primary {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

.action-btn.primary:hover {
  background: #2563eb;
}

.fm-content {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.fm-sidebar {
  width: 260px;
  background: white;
  border-right: 1px solid #e5e7eb;
  overflow-y: auto;
  padding: 16px 0;
}

.sidebar-section {
  margin-bottom: 24px;
}

.sidebar-section h3 {
  padding: 8px 20px;
  font-size: 12px;
  font-weight: 600;
  color: #9ca3af;
  text-transform: uppercase;
  margin-bottom: 8px;
}

.sidebar-items {
  display: flex;
  flex-direction: column;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 20px;
  cursor: pointer;
  transition: all 0.2s;
  color: #4b5563;
  font-size: 14px;
}

.sidebar-item:hover {
  background: #f3f4f6;
  color: #1f2937;
}

.sidebar-item.active {
  background: #eff6ff;
  color: #3b82f6;
}

.storage-info {
  padding: 0 20px;
}

.storage-item {
  margin-bottom: 16px;
}

.storage-label {
  font-size: 12px;
  color: #6b7280;
  margin-bottom: 4px;
}

.storage-bar {
  height: 6px;
  background: #e5e7eb;
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 4px;
}

.storage-fill {
  height: 100%;
  background: #3b82f6;
  transition: width 0.3s;
}

.storage-text {
  font-size: 11px;
  color: #9ca3af;
}

.fm-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.fm-status {
  display: flex;
  justify-content: space-between;
  padding: 12px 20px;
  font-size: 12px;
  color: #6b7280;
  background: white;
  border-bottom: 1px solid #e5e7eb;
}

.files-grid {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;
  padding: 20px;
  overflow-y: auto;
  align-content: start;
}

.file-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.file-card:hover {
  border-color: #3b82f6;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.1);
}

.file-card.selected {
  border-color: #3b82f6;
  background: #eff6ff;
}

.file-icon {
  color: #6b7280;
  margin-bottom: 12px;
}

.file-info {
  text-align: center;
}

.file-name {
  font-size: 13px;
  color: #1f2937;
  margin-bottom: 4px;
  word-break: break-word;
}

.file-meta {
  font-size: 11px;
  color: #9ca3af;
}

.files-list {
  flex: 1;
  overflow-y: auto;
}

.files-table {
  width: 100%;
  border-collapse: collapse;
}

.files-table thead {
  background: #f9fafb;
  position: sticky;
  top: 0;
}

.files-table th {
  padding: 12px 20px;
  text-align: left;
  font-size: 12px;
  font-weight: 600;
  color: #6b7280;
  cursor: pointer;
  user-select: none;
}

.files-table th:hover {
  color: #1f2937;
}

.sort-icon {
  display: inline-block;
  margin-left: 4px;
  opacity: 0.5;
}

.files-table tbody tr {
  border-bottom: 1px solid #e5e7eb;
  cursor: pointer;
  transition: all 0.2s;
}

.files-table tbody tr:hover {
  background: #f9fafb;
}

.files-table tbody tr.selected {
  background: #eff6ff;
}

.files-table td {
  padding: 12px 20px;
  font-size: 14px;
  color: #4b5563;
}

.file-name-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.upload-modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.2);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.upload-content {
  width: 90%;
  max-width: 600px;
  background: white;
  border-radius: 16px;
  overflow: hidden;
}

.upload-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #e5e7eb;
}

.upload-header h3 {
  font-size: 18px;
  font-weight: 600;
}

.upload-header button {
  padding: 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 6px;
}

.upload-header button:hover {
  background: #f3f4f6;
}

.upload-area {
  padding: 60px 20px;
  text-align: center;
  border: 2px dashed #e5e7eb;
  margin: 20px;
  border-radius: 12px;
  transition: all 0.2s;
}

.upload-area:hover {
  border-color: #3b82f6;
  background: #f9fafb;
}

.upload-area svg {
  color: #9ca3af;
  margin-bottom: 16px;
}

.upload-area p {
  color: #6b7280;
  margin-bottom: 8px;
}

.browse-btn {
  padding: 10px 20px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
}

.context-menu {
  position: fixed;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid #e5e7eb;
  min-width: 180px;
  z-index: 2000;
  overflow: hidden;
}

.context-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
}

.context-item:hover {
  background: #f3f4f6;
}

.context-item.danger {
  color: #ef4444;
}

.context-item.danger:hover {
  background: #fef2f2;
}

.context-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 4px 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>