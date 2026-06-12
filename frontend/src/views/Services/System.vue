<template>
  <div class="space-y-6">
    <!-- 系统服务 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-4 sm:p-6 border border-gray-700/50">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white flex items-center gap-2">
            <svg class="w-5 h-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            系统服务
          </h3>
          <p class="text-gray-500 text-sm">管理系统 systemd 服务</p>
        </div>
        <button
          @click="loadServices"
          :disabled="servicesLoading"
          class="flex items-center justify-center gap-2 px-4 py-2 bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-400 rounded-xl transition-all border border-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg
            class="w-4 h-4"
            :class="{ 'animate-spin': servicesLoading }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          刷新
        </button>
      </div>

      <!-- 服务加载状态 -->
      <div v-if="servicesLoading" class="flex items-center justify-center py-12">
        <div class="flex flex-col items-center gap-3">
          <svg class="w-8 h-8 text-indigo-400 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span class="text-gray-400">加载服务列表...</span>
        </div>
      </div>

      <!-- 服务错误状态 -->
      <div v-else-if="servicesError" class="bg-red-500/10 border border-red-500/20 rounded-xl p-6">
        <div class="flex items-center gap-3">
          <svg class="w-6 h-6 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <p class="text-red-400 font-medium">加载失败</p>
            <p class="text-red-400/70 text-sm">{{ servicesError }}</p>
          </div>
          <button @click="loadServices" class="ml-auto px-3 py-1.5 bg-red-500/20 hover:bg-red-500/30 text-red-400 rounded-lg text-sm transition-all">
            重试
          </button>
        </div>
      </div>

      <!-- 服务列表 -->
      <div v-else class="space-y-3">
        <div v-if="services.length === 0" class="text-center py-12 text-gray-500">
          <svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <p>暂无系统服务</p>
        </div>
        <div
          v-for="service in services"
          :key="service.name"
          class="bg-gray-900/50 rounded-xl p-4 sm:p-5 border border-gray-800 hover:border-gray-700 transition-all"
        >
          <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-3 mb-2 flex-wrap">
                <span class="text-white font-medium truncate">{{ service.name }}</span>
                <span
                  class="px-2.5 py-1 text-xs font-medium rounded-full flex items-center gap-1.5 flex-shrink-0"
                  :class="getServiceStatusClass(service.status, service.enabled)"
                >
                  <span
                    class="w-2 h-2 rounded-full animate-pulse"
                    :class="getServiceStatusDotClass(service.status)"
                  ></span>
                  {{ getServiceStatusText(service.status, service.enabled) }}
                </span>
                <span
                  v-if="service.enabled"
                  class="px-2 py-0.5 text-xs rounded bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 flex-shrink-0"
                  title="已启用开机自启"
                >
                  已启用
                </span>
              </div>
              <p class="text-gray-500 text-sm truncate">{{ service.description || '无描述' }}</p>
              <div class="flex items-center gap-4 mt-2 text-xs text-gray-600">
                <span v-if="service.loaded">Loaded: {{ service.loaded }}</span>
                <span v-if="service.active">Active: {{ service.active }}</span>
                <span v-if="service.sub">Sub: {{ service.sub }}</span>
              </div>
            </div>
            <div class="flex items-center gap-2 flex-shrink-0">
              <button
                v-if="service.status !== 'running'"
                @click="startService(service.name)"
                :disabled="serviceActionLoading[service.name]"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-green-500/10 hover:bg-green-500/20 text-green-400 rounded-lg text-sm transition-all border border-green-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                title="启动服务"
              >
                <svg
                  class="w-4 h-4"
                  :class="{ 'animate-spin': serviceActionLoading[service.name] }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="hidden sm:inline">启动</span>
              </button>
              <button
                v-if="service.status === 'running'"
                @click="stopService(service.name)"
                :disabled="serviceActionLoading[service.name]"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg text-sm transition-all border border-red-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                title="停止服务"
              >
                <svg
                  class="w-4 h-4"
                  :class="{ 'animate-spin': serviceActionLoading[service.name] }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 10a1 1 0 011-1h4a1 1 0 011 1v4a1 1 0 01-1 1h-4a1 1 0 01-1-1v-4z" />
                </svg>
                <span class="hidden sm:inline">停止</span>
              </button>
              <button
                @click="restartService(service.name)"
                :disabled="serviceActionLoading[service.name]"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-yellow-500/10 hover:bg-yellow-500/20 text-yellow-400 rounded-lg text-sm transition-all border border-yellow-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                title="重启服务"
              >
                <svg
                  class="w-4 h-4"
                  :class="{ 'animate-spin': serviceActionLoading[service.name] }"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                <span class="hidden sm:inline">重启</span>
              </button>
              <button
                v-if="!service.enabled"
                @click="enableService(service.name)"
                :disabled="serviceActionLoading[service.name]"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-indigo-500/10 hover:bg-indigo-500/20 text-indigo-400 rounded-lg text-sm transition-all border border-indigo-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                title="启用开机自启"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
              </button>
              <button
                v-if="service.enabled"
                @click="disableService(service.name)"
                :disabled="serviceActionLoading[service.name]"
                class="flex items-center gap-1.5 px-3 py-1.5 bg-gray-500/10 hover:bg-gray-500/20 text-gray-400 rounded-lg text-sm transition-all border border-gray-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
                title="禁用开机自启"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Docker 容器 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-4 sm:p-6 border border-gray-700/50">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white flex items-center gap-2">
            <svg class="w-5 h-5 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
              <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
            </svg>
            Docker 容器
          </h3>
          <p class="text-gray-500 text-sm">管理 Docker 容器</p>
        </div>
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
          刷新
        </button>
      </div>

      <!-- 容器加载状态 -->
      <div v-if="containersLoading" class="flex items-center justify-center py-12">
        <div class="flex flex-col items-center gap-3">
          <svg class="w-8 h-8 text-indigo-400 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span class="text-gray-400">加载容器列表...</span>
        </div>
      </div>

      <!-- 容器错误状态 -->
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

      <!-- 容器列表 -->
      <div v-else>
        <div v-if="containers.length === 0" class="text-center py-12 text-gray-500">
          <svg class="w-12 h-12 mx-auto mb-3 opacity-50" fill="currentColor" viewBox="0 0 24 24">
            <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
          </svg>
          <p>暂无 Docker 容器</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="container in containers"
            :key="container.id"
            class="bg-gray-900/50 rounded-xl p-5 border border-gray-800 hover:border-gray-700 transition-all hover:shadow-lg hover:shadow-gray-900/20"
          >
            <div class="flex items-start justify-between mb-4">
              <div class="flex items-center gap-3 min-w-0">
                <div class="w-10 h-10 bg-blue-500/20 rounded-lg flex items-center justify-center flex-shrink-0">
                  <svg class="w-5 h-5 text-blue-400" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M13.975 2.287a.75.75 0 00-1.357-.47l-1.96.548a.75.75 0 00-.517.735l.525 2.252a.75.75 0 01-.06.448l-.882 3.222a.75.75 0 01-.295.392l-3.228.882a.75.75 0 01-.447-.06l-2.252-.525a.75.75 0 00-.735.517l-.548 1.96a.75.75 0 00.47 1.357l2.252.548a.75.75 0 01.06.447l-.882 3.222a.75.75 0 01-.392.295l-3.222.882a.75.75 0 01.06.447l-.525 2.252a.75.75 0 00.517.735l1.96.548a.75.75 0 001.357-.47l.548-1.96a.75.75 0 01.517-.735l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.295-.392l3.222-.882a.75.75 0 01.447.06l2.252.525a.75.75 0 00.735-.517l.525-2.252a.75.75 0 01.735-.517l1.96-.548a.75.75 0 00.47-1.357l-2.252-.548a.75.75 0 01-.06-.447l.882-3.222a.75.75 0 01.392-.295l3.222-.882a.75.75 0 01-.06-.447l.525-2.252a.75.75 0 00-.517-.735z" />
                  </svg>
                </div>
                <div class="min-w-0">
                  <p class="text-white font-medium truncate" :title="container.name">{{ container.name }}</p>
                  <p class="text-gray-500 text-xs truncate" :title="container.image">{{ container.image }}</p>
                </div>
              </div>
              <span
                class="px-2.5 py-1 text-xs font-medium rounded-full flex items-center gap-1.5 flex-shrink-0 ml-2"
                :class="getContainerStateClass(container.state)"
              >
                <span
                  class="w-2 h-2 rounded-full"
                  :class="getContainerStateDotClass(container.state)"
                ></span>
                {{ container.state === 'running' ? '运行中' : '已停止' }}
              </span>
            </div>

            <!-- 容器信息 -->
            <div class="space-y-2 mb-4 text-sm">
              <div class="flex items-center justify-between text-gray-400">
                <span>端口:</span>
                <span class="text-gray-300">{{ formatPorts(container.ports) }}</span>
              </div>
              <div v-if="container.command" class="flex items-center justify-between text-gray-400">
                <span>命令:</span>
                <span class="text-gray-300 truncate ml-2" :title="container.command">{{ container.command }}</span>
              </div>
            </div>

            <!-- 容器操作按钮 -->
            <div class="flex gap-2">
              <button
                v-if="container.state !== 'running'"
                @click="startContainer(container.id)"
                :disabled="containerActionLoading[container.id]"
                class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 bg-green-500/10 hover:bg-green-500/20 text-green-400 rounded-lg text-sm transition-all border border-green-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
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
                启动
              </button>
              <button
                v-if="container.state === 'running'"
                @click="stopContainer(container.id)"
                :disabled="containerActionLoading[container.id]"
                class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-400 rounded-lg text-sm transition-all border border-red-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
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
                停止
              </button>
              <button
                @click="restartContainer(container.id)"
                :disabled="containerActionLoading[container.id]"
                class="flex-1 flex items-center justify-center gap-1.5 px-3 py-2 bg-yellow-500/10 hover:bg-yellow-500/20 text-yellow-400 rounded-lg text-sm transition-all border border-yellow-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
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
                重启
              </button>
            </div>
          </div>
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
import { ref, onMounted, reactive } from 'vue'
import { serviceApi } from '../../api'

// 类型定义
interface Service {
  name: string
  status: string
  description: string
  enabled?: boolean
  loaded?: string
  active?: string
  sub?: string
}

interface Container {
  id: string
  name: string
  image: string
  state: string
  ports?: string
  command?: string
}

// 状态数据
const services = ref<Service[]>([])
const containers = ref<Container[]>([])

// 加载状态
const servicesLoading = ref(false)
const containersLoading = ref(false)
const serviceActionLoading = reactive<Record<string, boolean>>({})
const containerActionLoading = reactive<Record<string, boolean>>({})

// 错误状态
const servicesError = ref('')
const containersError = ref('')

// Toast 通知
const toast = reactive({
  show: false,
  message: '',
  type: 'success' as 'success' | 'error'
})

const showToast = (message: string, type: 'success' | 'error' = 'success') => {
  toast.message = message
  toast.type = type
  toast.show = true
  setTimeout(() => {
    toast.show = false
  }, 3000)
}

// 服务状态样式
const getServiceStatusClass = (status: string, enabled?: boolean) => {
  if (status === 'running') {
    return 'bg-green-500/10 text-green-400 border border-green-500/20'
  } else if (status === 'stopped' || status === 'dead') {
    return 'bg-red-500/10 text-red-400 border border-red-500/20'
  } else if (status === 'failed') {
    return 'bg-red-500/10 text-red-400 border border-red-500/20'
  }
  return 'bg-gray-500/10 text-gray-400 border border-gray-500/20'
}

const getServiceStatusDotClass = (status: string) => {
  if (status === 'running') {
    return 'bg-green-400'
  }
  return 'bg-red-400'
}

const getServiceStatusText = (status: string, enabled?: boolean) => {
  if (status === 'running') return '运行中'
  if (status === 'stopped') return '已停止'
  if (status === 'failed') return '失败'
  return status
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

// 加载服务列表
const loadServices = async () => {
  servicesLoading.value = true
  servicesError.value = ''
  try {
    const data = await serviceApi.getServices()
    services.value = data.services || []
  } catch (error: any) {
    servicesError.value = error.response?.data?.message || '获取服务列表失败'
    console.error('获取服务列表失败:', error)
  } finally {
    servicesLoading.value = false
  }
}

// 加载容器列表
const loadContainers = async () => {
  containersLoading.value = true
  containersError.value = ''
  try {
    const data = await serviceApi.getContainers()
    containers.value = data.containers || []
  } catch (error: any) {
    containersError.value = error.response?.data?.message || '获取容器列表失败'
    console.error('获取容器列表失败:', error)
  } finally {
    containersLoading.value = false
  }
}

// 启动服务
const startService = async (name: string) => {
  serviceActionLoading[name] = true
  try {
    await serviceApi.startService(name)
    showToast(`服务 ${name} 已启动`)
    await loadServices()
  } catch (error: any) {
    showToast(error.response?.data?.message || `启动服务 ${name} 失败`, 'error')
  } finally {
    serviceActionLoading[name] = false
  }
}

// 停止服务
const stopService = async (name: string) => {
  serviceActionLoading[name] = true
  try {
    await serviceApi.stopService(name)
    showToast(`服务 ${name} 已停止`)
    await loadServices()
  } catch (error: any) {
    showToast(error.response?.data?.message || `停止服务 ${name} 失败`, 'error')
  } finally {
    serviceActionLoading[name] = false
  }
}

// 重启服务
const restartService = async (name: string) => {
  serviceActionLoading[name] = true
  try {
    await serviceApi.restartService(name)
    showToast(`服务 ${name} 已重启`)
    await loadServices()
  } catch (error: any) {
    showToast(error.response?.data?.message || `重启服务 ${name} 失败`, 'error')
  } finally {
    serviceActionLoading[name] = false
  }
}

// 启用服务
const enableService = async (name: string) => {
  serviceActionLoading[name] = true
  try {
    await serviceApi.enableService(name)
    showToast(`服务 ${name} 已启用开机自启`)
    await loadServices()
  } catch (error: any) {
    showToast(error.response?.data?.message || `启用服务 ${name} 失败`, 'error')
  } finally {
    serviceActionLoading[name] = false
  }
}

// 禁用服务
const disableService = async (name: string) => {
  serviceActionLoading[name] = true
  try {
    await serviceApi.disableService(name)
    showToast(`服务 ${name} 已禁用开机自启`)
    await loadServices()
  } catch (error: any) {
    showToast(error.response?.data?.message || `禁用服务 ${name} 失败`, 'error')
  } finally {
    serviceActionLoading[name] = false
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

// 初始化加载
onMounted(() => {
  loadServices()
  loadContainers()
})
</script>
