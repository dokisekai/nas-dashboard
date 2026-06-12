<template>
  <div class="space-y-6">
    <!-- Docker 概览统计 -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-gray-800/50 backdrop-blur rounded-xl p-4 border border-gray-700/50">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-blue-500/20 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
            </svg>
          </div>
          <div>
            <p class="text-gray-500 text-sm">总容器数</p>
            <p class="text-2xl font-bold text-white">{{ stats.total }}</p>
          </div>
        </div>
      </div>
      <div class="bg-gray-800/50 backdrop-blur rounded-xl p-4 border border-gray-700/50">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-green-500/20 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <p class="text-gray-500 text-sm">运行中</p>
            <p class="text-2xl font-bold text-white">{{ stats.running }}</p>
          </div>
        </div>
      </div>
      <div class="bg-gray-800/50 backdrop-blur rounded-xl p-4 border border-gray-700/50">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-red-500/20 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div>
            <p class="text-gray-500 text-sm">已停止</p>
            <p class="text-2xl font-bold text-white">{{ stats.stopped }}</p>
          </div>
        </div>
      </div>
      <div class="bg-gray-800/50 backdrop-blur rounded-xl p-4 border border-gray-700/50">
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-purple-500/20 rounded-xl flex items-center justify-center">
            <svg class="w-6 h-6 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
            </svg>
          </div>
          <div>
            <p class="text-gray-500 text-sm">镜像数</p>
            <p class="text-2xl font-bold text-white">{{ stats.images }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Docker 容器列表 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-4 sm:p-6 border border-gray-700/50">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white flex items-center gap-2">
            <svg class="w-5 h-5 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
            </svg>
            容器列表
          </h3>
          <p class="text-gray-500 text-sm">管理 Docker 容器</p>
        </div>
        <div class="flex items-center gap-3">
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索容器..."
              class="w-48 sm:w-64 pl-10 pr-4 py-2 bg-gray-900/50 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500/50 transition-colors text-sm"
            />
            <svg class="w-4 h-4 text-gray-500 absolute left-3 top-1/2 -translate-y-1/2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <select
            v-model="filterState"
            class="px-3 py-2 bg-gray-900/50 border border-gray-700 rounded-xl text-white text-sm focus:outline-none focus:border-indigo-500/50 transition-colors"
          >
            <option value="all">全部状态</option>
            <option value="running">运行中</option>
            <option value="stopped">已停止</option>
          </select>
          <button
            @click="loadContainers"
            :disabled="containersLoading"
            class="flex items-center justify-center gap-2 px-4 py-2 bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-400 rounded-xl transition-all border border-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg
              class="w-4 h-4"
              :class="{ 'animate-spin': containersLoading }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span class="hidden sm:inline">刷新</span>
          </button>
        </div>
      </div>

      <!-- 容器表格 -->
      <div v-if="containersLoading" class="flex items-center justify-center py-12">
        <div class="flex flex-col items-center gap-3">
          <svg class="w-8 h-8 text-indigo-400 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span class="text-gray-400">加载容器列表...</span>
        </div>
      </div>

      <div v-else-if="containersError" class="bg-red-500/10 border border-red-500/20 rounded-xl p-6">
        <div class="flex items-center gap-3">
          <svg class="w-6 h-6 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <p class="text-red-400 font-medium">加载失败</p>
            <p class="text-red-400/70 text-sm">{{ containersError }}</p>
          </div>
          <button @click="loadContainers" class="ml-auto px-3 py-1.5 bg-red-500/20 hover:bg-red-500/30 text-red-400 rounded-lg text-sm transition-all">
            重试
          </button>
        </div>
      </div>

      <div v-else class="overflow-x-auto">
        <div v-if="filteredContainers.length === 0" class="text-center py-12 text-gray-500">
          <svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="currentColor" viewBox="0 0 24 24">
            <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
          </svg>
          <p>暂无容器</p>
        </div>
        <table v-else class="w-full">
          <thead>
            <tr class="border-b border-gray-700/50">
              <th class="text-left py-3 px-4 text-gray-400 font-medium text-sm">容器名称</th>
              <th class="text-left py-3 px-4 text-gray-400 font-medium text-sm">镜像</th>
              <th class="text-left py-3 px-4 text-gray-400 font-medium text-sm">状态</th>
              <th class="text-left py-3 px-4 text-gray-400 font-medium text-sm hidden sm:table-cell">端口</th>
              <th class="text-right py-3 px-4 text-gray-400 font-medium text-sm">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="container in filteredContainers"
              :key="container.id"
              class="border-b border-gray-700/30 hover:bg-gray-700/20 transition-colors"
            >
              <td class="py-3 px-4">
                <div class="flex items-center gap-2">
                  <svg class="w-4 h-4 text-blue-400 flex-shrink-0" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
                  </svg>
                  <span class="text-white font-medium truncate" :title="container.name">{{ container.name }}</span>
                </div>
              </td>
              <td class="py-3 px-4">
                <span class="text-gray-400 text-sm truncate" :title="container.image">{{ container.image }}</span>
              </td>
              <td class="py-3 px-4">
                <span
                  class="px-2.5 py-1 text-xs font-medium rounded-full inline-flex items-center gap-1.5"
                  :class="getContainerStateClass(container.state)"
                >
                  <span
                    class="w-1.5 h-1.5 rounded-full"
                    :class="getContainerStateDotClass(container.state)"
                  ></span>
                  {{ container.state === 'running' ? '运行中' : '已停止' }}
                </span>
              </td>
              <td class="py-3 px-4 hidden sm:table-cell">
                <span class="text-gray-400 text-sm">{{ formatPorts(container.ports) }}</span>
              </td>
              <td class="py-3 px-4">
                <div class="flex items-center justify-end gap-2">
                  <button
                    v-if="container.state !== 'running'"
                    @click="startContainer(container.id)"
                    :disabled="containerActionLoading[container.id]"
                    class="p-2 bg-green-500/10 hover:bg-green-500/20 text-green-400 rounded-lg transition-all border border-green-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                    title="启动"
                  >
                    <svg
                      class="w-4 h-4"
                      :class="{ 'animate-spin': containerActionLoading[container.id] }"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </button>
                  <button
                    v-if="container.state === 'running'"
                    @click="stopContainer(container.id)"
                    :disabled="containerActionLoading[container.id]"
                    class="p-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg transition-all border border-red-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                    title="停止"
                  >
                    <svg
                      class="w-4 h-4"
                      :class="{ 'animate-spin': containerActionLoading[container.id] }"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
                    </svg>
                  </button>
                  <button
                    @click="restartContainer(container.id)"
                    :disabled="containerActionLoading[container.id]"
                    class="p-2 bg-yellow-500/10 hover:bg-yellow-500/20 text-yellow-400 rounded-lg transition-all border border-yellow-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                    title="重启"
                  >
                    <svg
                      class="w-4 h-4"
                      :class="{ 'animate-spin': containerActionLoading[container.id] }"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                  </button>
                  <button
                    @click="viewLogs(container.id, container.name)"
                    class="p-2 bg-gray-700/50 hover:bg-gray-700 text-gray-300 rounded-lg transition-all"
                    title="查看日志"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </button>
                  <button
                    @click="confirmRemove(container.id, container.name)"
                    class="p-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg transition-all border border-red-500/20"
                    title="删除"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 镜像管理 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-4 sm:p-6 border border-gray-700/50">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white flex items-center gap-2">
            <svg class="w-5 h-5 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            镜像列表
          </h3>
          <p class="text-gray-500 text-sm">管理 Docker 镜像</p>
        </div>
        <div class="flex items-center gap-3">
          <button
            @click="showPullImageModal = true"
            class="flex items-center gap-2 px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-xl transition-all"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            拉取镜像
          </button>
          <button
            @click="loadImages"
            :disabled="imagesLoading"
            class="flex items-center justify-center gap-2 px-4 py-2 bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-400 rounded-xl transition-all border border-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <svg
              class="w-4 h-4"
              :class="{ 'animate-spin': imagesLoading }"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span class="hidden sm:inline">刷新</span>
          </button>
        </div>
      </div>

      <!-- 镜像列表 -->
      <div v-if="imagesLoading" class="flex items-center justify-center py-12">
        <div class="flex flex-col items-center gap-3">
          <svg class="w-8 h-8 text-indigo-400 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span class="text-gray-400">加载镜像列表...</span>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-if="images.length === 0" class="col-span-full text-center py-12 text-gray-500">
          <svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <p>暂无镜像</p>
        </div>
        <div
          v-for="image in images"
          :key="image.id"
          class="bg-gray-900/50 rounded-xl p-4 border border-gray-800 hover:border-gray-700 transition-all"
        >
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1 min-w-0">
              <p class="text-white font-medium truncate" :title="image.repoTags?.[0] || image.id">{{ image.repoTags?.[0] || image.id.substring(0, 12) }}</p>
              <p class="text-gray-500 text-xs mt-1">{{ formatSize(image.size) }}</p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button
              @click="confirmRemoveImage(image.id, image.repoTags?.[0] || image.id)"
              class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg text-sm transition-all border border-red-500/20"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 拉取镜像弹窗 -->
    <div v-if="showPullImageModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="showPullImageModal = false"></div>
      <div class="relative bg-gray-800 rounded-2xl p-6 w-full max-w-md border border-gray-700">
        <h3 class="text-lg font-semibold text-white mb-4">拉取镜像</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm text-gray-400 mb-2">镜像名称</label>
            <input
              v-model="pullImageName"
              type="text"
              placeholder="例如: nginx:latest"
              class="w-full px-4 py-2 bg-gray-900/50 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:border-indigo-500/50 transition-colors"
              @keypress.enter="pullImage"
            />
          </div>
          <div class="flex gap-3">
            <button
              @click="showPullImageModal = false"
              class="flex-1 px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-xl transition-all"
            >
              取消
            </button>
            <button
              @click="pullImage"
              :disabled="pullingImage || !pullImageName"
              class="flex-1 px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-xl transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ pullingImage ? '拉取中...' : '拉取' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 日志弹窗 -->
    <div v-if="showLogsModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="showLogsModal = false"></div>
      <div class="relative bg-gray-800 rounded-2xl p-6 w-full max-w-4xl h-96 flex flex-col border border-gray-700">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-white">容器日志 - {{ logsContainerName }}</h3>
          <button @click="showLogsModal = false" class="text-gray-400 hover:text-white transition-colors">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="flex-1 overflow-auto bg-gray-900 rounded-xl p-4">
          <pre class="text-sm text-gray-300 whitespace-pre-wrap font-mono">{{ logs }}</pre>
        </div>
        <div class="flex items-center justify-end gap-3 mt-4">
          <button
            @click="loadLogs"
            :disabled="loadingLogs"
            class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-lg text-sm transition-all"
          >
            <svg class="w-4 h-4 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
          <button
            @click="showLogsModal = false"
            class="px-4 py-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-lg text-sm transition-all"
          >
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- Toast 通知 -->
    <Transition
      enter-active-class="transform ease-out duration-300 transition"
      enter-from-class="translate-y-2 opacity-0 sm:translate-y-0 sm:translate-x-2"
      enter-to-class="translate-y-0 opacity-100 sm:translate-x-0"
      leave-active-class="transition ease-in duration-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="toast.show"
        class="fixed bottom-4 right-4 z-50 max-w-sm w-full bg-gray-800 border rounded-xl shadow-lg p-4 flex items-center gap-3"
        :class="toast.type === 'success' ? 'border-green-500/30' : 'border-red-500/30'"
      >
        <svg
          v-if="toast.type === 'success'"
          class="w-6 h-6 text-green-400 flex-shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <svg
          v-else
          class="w-6 h-6 text-red-400 flex-shrink-0"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="flex-1">
          <p class="text-white" :class="toast.type === 'success' ? 'text-green-400' : 'text-red-400'">
            {{ toast.message }}
          </p>
        </div>
        <button @click="toast.show = false" class="text-gray-400 hover:text-white transition-colors">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue'
import { serviceApi } from '../../api'

// 类型定义
interface Container {
  id: string
  name: string
  image: string
  state: string
  ports?: string
  command?: string
}

interface Image {
  id: string
  repoTags?: string[]
  size: number
  created: number
}

interface Stats {
  total: number
  running: number
  stopped: number
  images: number
}

// 状态数据
const containers = ref<Container[]>([])
const images = ref<Image[]>([])
const stats = ref<Stats>({ total: 0, running: 0, stopped: 0, images: 0 })

// 加载状态
const containersLoading = ref(false)
const imagesLoading = ref(false)
const containerActionLoading = reactive<Record<string, boolean>>({})

// 错误状态
const containersError = ref('')

// 搜索和过滤
const searchQuery = ref('')
const filterState = ref<'all' | 'running' | 'stopped'>('all')

// 弹窗状态
const showPullImageModal = ref(false)
const showLogsModal = ref(false)
const pullImageName = ref('')
const pullingImage = ref(false)

// 日志相关
const logs = ref('')
const logsContainerId = ref('')
const logsContainerName = ref('')
const loadingLogs = ref(false)

// Toast 通知
const toast = reactive({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error'
})

// 过滤后的容器列表
const filteredContainers = computed(() => {
  let result = containers.value

  if (filterState.value !== 'all') {
    result = result.filter(c => c.state === filterState.value)
  }

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(c =>
      c.name.toLowerCase().includes(query) ||
      c.image.toLowerCase().includes(query)
    )
  }

  return result
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.message = message
  toast.type = type
  toast.show = true
  setTimeout(() => {
    toast.show = false
  }, 3000)
}

// 容器状态样式
const getContainerStateClass = (state: string) => {
  if (state === 'running') {
    return 'bg-green-500/10 text-green-400 border border-green-500/20'
  }
  return 'bg-gray-500/10 text-gray-400 border border-gray-500/20'
}

const getContainerStateDotClass = (state: string) => {
  if (state === 'running') {
    return 'bg-green-400 animate-pulse'
  }
  return 'bg-gray-400'
}

// 格式化端口显示
const formatPorts = (ports?: string) => {
  if (!ports) return '-'
  const portList = ports.split(',').slice(0, 2)
  return portList.join(', ') + (ports.split(',').length > 2 ? '...' : '')
}

// 格式化大小
const formatSize = (bytes: number) => {
  const units = ['B', 'KB', 'MB', 'GB']
  let size = bytes
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(1)} ${units[unitIndex]}`
}

// 加载容器列表
const loadContainers = async () => {
  containersLoading.value = true
  containersError.value = ''
  try {
    const data = await serviceApi.getContainers()
    containers.value = data.containers || []
    // 更新统计
    stats.value.total = containers.value.length
    stats.value.running = containers.value.filter(c => c.state === 'running').length
    stats.value.stopped = containers.value.filter(c => c.state !== 'running').length
  } catch (error: any) {
    containersError.value = error.response?.data?.message || '获取容器列表失败'
    console.error('获取容器列表失败:', error)
  } finally {
    containersLoading.value = false
  }
}

// 加载镜像列表
const loadImages = async () => {
  imagesLoading.value = true
  try {
    const data = await serviceApi.getImages()
    images.value = data.images || []
    stats.value.images = images.value.length
  } catch (error: any) {
    showToast(error.response?.data?.message || '获取镜像列表失败', 'error')
    console.error('获取镜像列表失败:', error)
  } finally {
    imagesLoading.value = false
  }
}

// 启动容器
const startContainer = async (id: string) => {
  containerActionLoading[id] = true
  try {
    await serviceApi.startContainer(id)
    showToast('容器已启动')
    await loadContainers()
  } catch (error: any) {
    showToast(error.response?.data?.message || '启动容器失败', 'error')
  } finally {
    containerActionLoading[id] = false
  }
}

// 停止容器
const stopContainer = async (id: string) => {
  containerActionLoading[id] = true
  try {
    await serviceApi.stopContainer(id)
    showToast('容器已停止')
    await loadContainers()
  } catch (error: any) {
    showToast(error.response?.data?.message || '停止容器失败', 'error')
  } finally {
    containerActionLoading[id] = false
  }
}

// 重启容器
const restartContainer = async (id: string) => {
  containerActionLoading[id] = true
  try {
    await serviceApi.restartContainer(id)
    showToast('容器已重启')
    await loadContainers()
  } catch (error: any) {
    showToast(error.response?.data?.message || '重启容器失败', 'error')
  } finally {
    containerActionLoading[id] = false
  }
}

// 删除容器
const removeContainer = async (id: string) => {
  containerActionLoading[id] = true
  try {
    await serviceApi.removeContainer(id)
    showToast('容器已删除')
    await loadContainers()
  } catch (error: any) {
    showToast(error.response?.data?.message || '删除容器失败', 'error')
  } finally {
    containerActionLoading[id] = false
  }
}

// 确认删除容器
const confirmRemove = (id: string, name: string) => {
  if (confirm(`确定要删除容器 "${name}" 吗？`)) {
    removeContainer(id)
  }
}

// 删除镜像
const removeImage = async (id: string) => {
  try {
    await serviceApi.removeImage(id)
    showToast('镜像已删除')
    await loadImages()
  } catch (error: any) {
    showToast(error.response?.data?.message || '删除镜像失败', 'error')
  }
}

// 确认删除镜像
const confirmRemoveImage = (id: string, name: string) => {
  if (confirm(`确定要删除镜像 "${name}" 吗？`)) {
    removeImage(id)
  }
}

// 拉取镜像
const pullImage = async () => {
  if (!pullImageName.value) return
  pullingImage.value = true
  try {
    await serviceApi.pullImage(pullImageName.value)
    showToast('镜像拉取成功')
    showPullImageModal.value = false
    pullImageName.value = ''
    await loadImages()
  } catch (error: any) {
    showToast(error.response?.data?.message || '拉取镜像失败', 'error')
  } finally {
    pullingImage.value = false
  }
}

// 查看日志
const viewLogs = (id: string, name: string) => {
  logsContainerId.value = id
  logsContainerName.value = name
  showLogsModal.value = true
  loadLogs()
}

// 加载日志
const loadLogs = async () => {
  loadingLogs.value = true
  try {
    const data = await serviceApi.getContainerLogs(logsContainerId.value)
    logs.value = data.logs || '无日志'
  } catch (error: any) {
    logs.value = error.response?.data?.message || '获取日志失败'
  } finally {
    loadingLogs.value = false
  }
}

// 初始化加载
onMounted(() => {
  loadContainers()
  loadImages()
})
</script>
