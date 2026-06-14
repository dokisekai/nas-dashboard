<template>
  <div class="space-y-6 p-6">
    <!-- 用户列表 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">用户列表</h3>
          <p class="text-gray-500 text-sm">管理系统用户</p>
        </div>
        <button
          @click="openCreateUserModal"
          :disabled="usersLoading"
          class="flex items-center justify-center gap-2 px-4 py-2 bg-indigo-500 hover:bg-indigo-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-xl transition-all shadow-lg shadow-indigo-500/25"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          添加用户
        </button>
      </div>

      <!-- 加载状态 -->
      <div v-if="usersLoading" class="flex items-center justify-center py-12">
        <div class="flex flex-col items-center gap-3">
          <div class="w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
          <p class="text-gray-400 text-sm">加载用户列表...</p>
        </div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="usersError" class="flex flex-col items-center justify-center py-12">
        <div class="w-12 h-12 bg-red-500/20 rounded-full flex items-center justify-center mb-3">
          <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <p class="text-red-400 mb-2">{{ usersError }}</p>
        <button @click="loadUsers" class="text-indigo-400 hover:text-indigo-300 text-sm">重试</button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="users.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="w-16 h-16 bg-gray-700/50 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
        </div>
        <p class="text-gray-400 mb-1">暂无用户</p>
        <p class="text-gray-500 text-sm">点击上方按钮添加第一个用户</p>
      </div>

      <!-- 用户列表 -->
      <div v-else class="overflow-hidden rounded-xl border border-gray-800">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-900/50">
              <tr class="text-gray-400 text-sm">
                <th class="text-left py-4 px-4 font-medium">用户名</th>
                <th class="text-left py-4 px-4 font-medium">UID</th>
                <th class="text-left py-4 px-4 font-medium hidden md:table-cell">组</th>
                <th class="text-left py-4 px-4 font-medium hidden lg:table-cell">家目录</th>
                <th class="text-left py-4 px-4 font-medium hidden sm:table-cell">Shell</th>
                <th class="text-right py-4 px-4 font-medium">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="user in users"
                :key="user.username"
                class="border-t border-gray-800 hover:bg-gray-900/30 transition-colors"
              >
                <td class="py-4 px-4">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center flex-shrink-0">
                      <span class="text-white font-medium">{{ user.username.charAt(0).toUpperCase() }}</span>
                    </div>
                    <span class="text-white font-medium truncate">{{ user.username }}</span>
                  </div>
                </td>
                <td class="py-4 px-4 text-gray-400 font-mono text-sm">{{ user.uid }}</td>
                <td class="py-4 px-4 text-gray-400 hidden md:table-cell">{{ user.group }}</td>
                <td class="py-4 px-4 text-gray-400 font-mono text-sm hidden lg:table-cell">{{ user.home }}</td>
                <td class="py-4 px-4 text-gray-400 font-mono text-sm hidden sm:table-cell">{{ user.shell }}</td>
                <td class="py-4 px-4 text-right">
                  <div class="flex items-center justify-end gap-2">
                    <button
                      @click="openEditUserModal(user)"
                      class="p-2 text-indigo-400 hover:text-indigo-300 hover:bg-indigo-500/10 rounded-lg transition-colors"
                      title="编辑用户"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <button
                      @click="confirmDeleteUser(user)"
                      class="p-2 text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-colors"
                      title="删除用户"
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
    </div>

    <!-- SSH 密钥管理 -->
    <div class="bg-gray-800/50 backdrop-blur rounded-2xl p-6 border border-gray-700/50">
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-6">
        <div>
          <h3 class="text-lg font-semibold text-white">SSH 密钥管理</h3>
          <p class="text-gray-500 text-sm">管理 SSH 公钥</p>
        </div>
        <button
          @click="openAddKeyModal"
          :disabled="keysLoading"
          class="flex items-center justify-center gap-2 px-4 py-2 bg-emerald-500 hover:bg-emerald-600 disabled:bg-gray-600 disabled:cursor-not-allowed text-white rounded-xl transition-all shadow-lg shadow-emerald-500/25"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
          添加密钥
        </button>
      </div>

      <!-- 加载状态 -->
      <div v-if="keysLoading" class="flex items-center justify-center py-12">
        <div class="flex flex-col items-center gap-3">
          <div class="w-8 h-8 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin"></div>
          <p class="text-gray-400 text-sm">加载 SSH 密钥...</p>
        </div>
      </div>

      <!-- 错误状态 -->
      <div v-else-if="keysError" class="flex flex-col items-center justify-center py-12">
        <div class="w-12 h-12 bg-red-500/20 rounded-full flex items-center justify-center mb-3">
          <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <p class="text-red-400 mb-2">{{ keysError }}</p>
        <button @click="loadSSHKeys" class="text-indigo-400 hover:text-indigo-300 text-sm">重试</button>
      </div>

      <!-- 空状态 -->
      <div v-else-if="sshKeys.length === 0" class="flex flex-col items-center justify-center py-12">
        <div class="w-16 h-16 bg-gray-700/50 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
        </div>
        <p class="text-gray-400 mb-1">暂无 SSH 密钥</p>
        <p class="text-gray-500 text-sm">添加 SSH 公钥以便安全登录</p>
      </div>

      <!-- SSH 密钥列表 -->
      <div v-else class="space-y-3">
        <div
          v-for="key in sshKeys"
          :key="key.id"
          class="bg-gray-900/50 rounded-xl p-4 sm:p-5 border border-gray-800 hover:border-gray-700 transition-all"
        >
          <div class="flex items-start justify-between gap-4">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-3 mb-2">
                <div class="w-10 h-10 bg-emerald-500/20 rounded-lg flex items-center justify-center flex-shrink-0">
                  <svg class="w-5 h-5 text-emerald-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                  </svg>
                </div>
                <div class="min-w-0 flex-1">
                  <p class="text-white font-medium truncate">{{ key.name }}</p>
                  <p class="text-gray-500 text-xs">{{ key.addedAt }}</p>
                </div>
              </div>
              <p class="text-gray-400 text-sm font-mono break-all">{{ key.fingerprint }}</p>
            </div>
            <button
              @click="confirmDeleteKey(key)"
              class="p-2 text-red-400 hover:text-red-300 hover:bg-red-500/10 rounded-lg transition-colors flex-shrink-0"
              title="删除密钥"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 用户表单模态框 -->
    <Transition name="modal">
      <div
        v-if="showUserModal"
        class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="closeUserModal"
      >
        <div class="bg-gray-800 rounded-2xl p-6 w-full max-w-md border border-gray-700 shadow-2xl max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-white">
              {{ isEditing ? '编辑用户' : '添加用户' }}
            </h3>
            <button
              @click="closeUserModal"
              class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="saveUser" class="space-y-4">
            <!-- 用户名 -->
            <div>
              <label class="block text-sm text-gray-300 mb-2">
                用户名 <span class="text-red-400">*</span>
              </label>
              <input
                v-model="userForm.username"
                type="text"
                :disabled="isEditing"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 disabled:bg-gray-800 disabled:cursor-not-allowed rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
                placeholder="username"
                :class="{ 'border-red-500': userErrors.username }"
              />
              <p v-if="userErrors.username" class="text-red-400 text-sm mt-1">{{ userErrors.username }}</p>
            </div>

            <!-- 密码 (仅新建时显示) -->
            <div v-if="!isEditing">
              <label class="block text-sm text-gray-300 mb-2">
                密码 <span class="text-red-400">*</span>
              </label>
              <div class="relative">
                <input
                  v-model="userForm.password"
                  :type="showPassword ? 'text' : 'password'"
                  class="w-full px-4 py-3 pr-12 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
                  placeholder="•••••••"
                  :class="{ 'border-red-500': userErrors.password }"
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-300"
                >
                  <svg v-if="showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                  </svg>
                </button>
              </div>
              <p v-if="userErrors.password" class="text-red-400 text-sm mt-1">{{ userErrors.password }}</p>
              <div class="flex items-center gap-2 mt-2">
                <div class="flex-1 h-1 bg-gray-700 rounded-full overflow-hidden">
                  <div
                    class="h-full transition-all duration-300"
                    :class="passwordStrengthColor"
                    :style="{ width: passwordStrengthPercent }"
                  ></div>
                </div>
                <span class="text-xs" :class="passwordStrengthTextColor">{{ passwordStrengthText }}</span>
              </div>
            </div>

            <!-- 组 -->
            <div>
              <label class="block text-sm text-gray-300 mb-2">
                用户组 <span class="text-red-400">*</span>
              </label>
              <input
                v-model="userForm.group"
                type="text"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
                placeholder="group"
                :class="{ 'border-red-500': userErrors.group }"
              />
              <p v-if="userErrors.group" class="text-red-400 text-sm mt-1">{{ userErrors.group }}</p>
            </div>

            <!-- Shell -->
            <div>
              <label class="block text-sm text-gray-300 mb-2">Shell</label>
              <select
                v-model="userForm.shell"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
              >
                <option value="/bin/bash">/bin/bash</option>
                <option value="/bin/zsh">/bin/zsh</option>
                <option value="/bin/sh">/bin/sh</option>
                <option value="/usr/sbin/nologin">/usr/sbin/nologin (无登录权限)</option>
              </select>
            </div>

            <!-- 操作按钮 -->
            <div class="flex justify-end gap-3 pt-4">
              <button
                type="button"
                @click="closeUserModal"
                class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-white rounded-xl transition-colors"
              >
                取消
              </button>
              <button
                type="submit"
                :disabled="userSaving"
                class="px-5 py-2.5 bg-indigo-500 hover:bg-indigo-600 disabled:bg-indigo-700 disabled:cursor-not-allowed text-white rounded-xl transition-colors flex items-center gap-2"
              >
                <div v-if="userSaving" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                {{ userSaving ? '保存中...' : '保存' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- SSH 密钥表单模态框 -->
    <Transition name="modal">
      <div
        v-if="showKeyModal"
        class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="closeKeyModal"
      >
        <div class="bg-gray-800 rounded-2xl p-6 w-full max-w-lg border border-gray-700 shadow-2xl max-h-[90vh] overflow-y-auto">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-white">添加 SSH 密钥</h3>
            <button
              @click="closeKeyModal"
              class="p-2 text-gray-400 hover:text-white hover:bg-gray-700 rounded-lg transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="saveKey" class="space-y-4">
            <!-- 名称 -->
            <div>
              <label class="block text-sm text-gray-300 mb-2">
                名称 <span class="text-red-400">*</span>
              </label>
              <input
                v-model="keyForm.name"
                type="text"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-colors"
                placeholder="例如: MacBook Pro"
                :class="{ 'border-red-500': keyErrors.name }"
              />
              <p v-if="keyErrors.name" class="text-red-400 text-sm mt-1">{{ keyErrors.name }}</p>
            </div>

            <!-- 公钥内容 -->
            <div>
              <label class="block text-sm text-gray-300 mb-2">
                公钥内容 <span class="text-red-400">*</span>
              </label>
              <textarea
                v-model="keyForm.content"
                rows="6"
                class="w-full px-4 py-3 bg-gray-900 border border-gray-700 rounded-xl text-white placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition-colors text-sm font-mono resize-none"
                placeholder="ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..."
                :class="{ 'border-red-500': keyErrors.content }"
              ></textarea>
              <p v-if="keyErrors.content" class="text-red-400 text-sm mt-1">{{ keyErrors.content }}</p>
              <p class="text-gray-500 text-xs mt-2">
                通常位于 ~/.ssh/id_rsa.pub 或 ~/.ssh/id_ed25519.pub
              </p>
            </div>

            <!-- 操作按钮 -->
            <div class="flex justify-end gap-3 pt-4">
              <button
                type="button"
                @click="closeKeyModal"
                class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-white rounded-xl transition-colors"
              >
                取消
              </button>
              <button
                type="submit"
                :disabled="keySaving"
                class="px-5 py-2.5 bg-emerald-500 hover:bg-emerald-600 disabled:bg-emerald-700 disabled:cursor-not-allowed text-white rounded-xl transition-colors flex items-center gap-2"
              >
                <div v-if="keySaving" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                {{ keySaving ? '添加中...' : '添加' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>

    <!-- 删除确认对话框 -->
    <Transition name="modal">
      <div
        v-if="showDeleteModal"
        class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4"
        @click.self="closeDeleteModal"
      >
        <div class="bg-gray-800 rounded-2xl p-6 w-full max-w-sm border border-gray-700 shadow-2xl">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-12 h-12 bg-red-500/20 rounded-full flex items-center justify-center flex-shrink-0">
              <svg class="w-6 h-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-white">确认删除</h3>
              <p class="text-gray-400 text-sm">此操作不可撤销</p>
            </div>
          </div>
          <p class="text-gray-300 mb-6">
            {{ deleteTarget.type === 'user'
              ? `确定要删除用户 "${deleteTarget.name}" 吗？该用户的所有数据将被永久删除。`
              : '确定要删除此 SSH 密钥吗？' }}
          </p>
          <div class="flex justify-end gap-3">
            <button
              @click="closeDeleteModal"
              class="px-5 py-2.5 bg-gray-700 hover:bg-gray-600 text-white rounded-xl transition-colors"
            >
              取消
            </button>
            <button
              @click="executeDelete"
              :disabled="deleting"
              class="px-5 py-2.5 bg-red-500 hover:bg-red-600 disabled:bg-red-700 disabled:cursor-not-allowed text-white rounded-xl transition-colors flex items-center gap-2"
            >
              <div v-if="deleting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              {{ deleting ? '删除中...' : '确认删除' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Toast 通知 -->
    <Transition name="toast">
      <div
        v-if="toast.show"
        class="fixed bottom-6 right-6 z-50 flex items-center gap-3 px-4 py-3 rounded-xl shadow-lg border"
        :class="{
          'bg-green-500/90 border-green-400': toast.type === 'success',
          'bg-red-500/90 border-red-400': toast.type === 'error',
          'bg-blue-500/90 border-blue-400': toast.type === 'info'
        }"
      >
        <svg v-if="toast.type === 'success'" class="w-5 h-5 text-white flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <svg v-else-if="toast.type === 'error'" class="w-5 h-5 text-white flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        <svg v-else class="w-5 h-5 text-white flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span class="text-white">{{ toast.message }}</span>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { userApi } from '../../api'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()

// 用户列表相关
const users = ref<any[]>([])
const usersLoading = ref(true)
const usersError = ref('')

// SSH 密钥相关
const sshKeys = ref<any[]>([])
const keysLoading = ref(true)
const keysError = ref('')

// 用户表单相关
const showUserModal = ref(false)
const isEditing = ref(false)
const userSaving = ref(false)
const showPassword = ref(false)
const userForm = ref({
  username: '',
  password: '',
  group: '',
  shell: '/bin/bash'
})
const userErrors = ref<Record<string, string>>({})

// SSH 密钥表单相关
const showKeyModal = ref(false)
const keySaving = ref(false)
const keyForm = ref({
  name: '',
  content: ''
})
const keyErrors = ref<Record<string, string>>({})

// 删除确认相关
const showDeleteModal = ref(false)
const deleteTarget = ref<{ type: 'user' | 'key', name: string, id?: string }>({ type: 'user', name: '' })
const deleting = ref(false)

// Toast 通知
const toast = ref<{ show: boolean, message: string, type: 'success' | 'error' | 'info' }>({
  show: false,
  message: '',
  type: 'info'
})

// 密码强度计算
const passwordStrength = computed(() => {
  const password = userForm.value.password
  if (!password) return 0

  let strength = 0
  if (password.length >= 8) strength += 25
  if (password.length >= 12) strength += 15
  if (/[a-z]/.test(password)) strength += 20
  if (/[A-Z]/.test(password)) strength += 20
  if (/[0-9]/.test(password)) strength += 10
  if (/[^a-zA-Z0-9]/.test(password)) strength += 10

  return Math.min(100, strength)
})

const passwordStrengthPercent = computed(() => `${passwordStrength.value}%`)

const passwordStrengthColor = computed(() => {
  const strength = passwordStrength.value
  if (strength < 40) return 'bg-red-500'
  if (strength < 70) return 'bg-yellow-500'
  return 'bg-green-500'
})

const passwordStrengthTextColor = computed(() => {
  const strength = passwordStrength.value
  if (strength < 40) return 'text-red-400'
  if (strength < 70) return 'text-yellow-400'
  return 'text-green-400'
})

const passwordStrengthText = computed(() => {
  const strength = passwordStrength.value
  if (strength < 40) return '弱'
  if (strength < 70) return '中等'
  return '强'
})

// 显示 Toast 通知
const showToast = (message: string, type: 'success' | 'error' | 'info' = 'info') => {
  toast.value = { show: true, message, type }
  setTimeout(() => {
    toast.value.show = false
  }, 3000)
}

// 加载用户列表
const loadUsers = async () => {
  usersLoading.value = true
  usersError.value = ''
  try {
    const data = await userApi.getUsers() as any
    users.value = data.users || []
  } catch (error: any) {
    console.error('获取用户列表失败:', error)
    usersError.value = error?.message || '获取用户列表失败'
  } finally {
    usersLoading.value = false
  }
}

// 加载 SSH 密钥
const loadSSHKeys = async () => {
  keysLoading.value = true
  keysError.value = ''
  try {
    const data = await userApi.getSSHKeys() as any
    sshKeys.value = data.keys || []
  } catch (error: any) {
    console.error('获取SSH密钥失败:', error)
    keysError.value = error?.message || '获取SSH密钥失败'
  } finally {
    keysLoading.value = false
  }
}

// 验证用户表单
const validateUserForm = () => {
  userErrors.value = {}

  if (!userForm.value.username.trim()) {
    userErrors.value.username = '请输入用户名'
  } else if (!/^[a-z_][a-z0-9_-]{0,31}$/.test(userForm.value.username)) {
    userErrors.value.username = '用户名只能包含小写字母、数字、下划线和连字符，且必须以字母或下划线开头'
  }

  if (!isEditing.value) {
    if (!userForm.value.password) {
      userErrors.value.password = '请输入密码'
    } else if (userForm.value.password.length < 6) {
      userErrors.value.password = '密码长度至少为6位'
    }
  }

  if (!userForm.value.group.trim()) {
    userErrors.value.group = '请输入用户组'
  }

  return Object.keys(userErrors.value).length === 0
}

// 验证密钥表单
const validateKeyForm = () => {
  keyErrors.value = {}

  if (!keyForm.value.name.trim()) {
    keyErrors.value.name = '请输入密钥名称'
  }

  if (!keyForm.value.content.trim()) {
    keyErrors.value.content = '请输入公钥内容'
  } else if (!keyForm.value.content.startsWith('ssh-')) {
    keyErrors.value.content = '请输入有效的 SSH 公钥格式'
  }

  return Object.keys(keyErrors.value).length === 0
}

// 打开创建用户模态框
const openCreateUserModal = () => {
  isEditing.value = false
  userForm.value = { username: '', password: '', group: '', shell: '/bin/bash' }
  userErrors.value = {}
  showPassword.value = false
  showUserModal.value = true
}

// 打开编辑用户模态框
const openEditUserModal = (user: any) => {
  isEditing.value = true
  userForm.value = {
    username: user.username,
    password: '',
    group: user.group,
    shell: user.shell
  }
  userErrors.value = {}
  showPassword.value = false
  showUserModal.value = true
}

// 关闭用户模态框
const closeUserModal = () => {
  showUserModal.value = false
  userForm.value = { username: '', password: '', group: '', shell: '/bin/bash' }
  userErrors.value = {}
}

// 保存用户
const saveUser = async () => {
  if (!validateUserForm()) return

  userSaving.value = true
  try {
    if (isEditing.value) {
      await userApi.updateUser(userForm.value.username, {
        group: userForm.value.group,
        shell: userForm.value.shell
      })
      showToast('用户更新成功', 'success')
    } else {
      await userApi.createUser({
        username: userForm.value.username,
        password: userForm.value.password,
        group: userForm.value.group,
        shell: userForm.value.shell
      })
      showToast('用户创建成功', 'success')
    }
    await loadUsers()
    closeUserModal()
  } catch (error: any) {
    console.error('保存用户失败:', error)
    showToast(error?.message || (isEditing.value ? '更新用户失败' : '创建用户失败'), 'error')
  } finally {
    userSaving.value = false
  }
}

// 确认删除用户
const confirmDeleteUser = (user: any) => {
  deleteTarget.value = { type: 'user', name: user.username }
  showDeleteModal.value = true
}

// 确认删除密钥
const confirmDeleteKey = (key: any) => {
  deleteTarget.value = { type: 'key', name: key.name, id: key.id }
  showDeleteModal.value = true
}

// 关闭删除确认模态框
const closeDeleteModal = () => {
  showDeleteModal.value = false
  deleteTarget.value = { type: 'user', name: '' }
}

// 执行删除
const executeDelete = async () => {
  deleting.value = true
  try {
    if (deleteTarget.value.type === 'user') {
      await userApi.deleteUser(deleteTarget.value.name)
      showToast(`用户 "${deleteTarget.value.name}" 已删除`, 'success')
      await loadUsers()
    } else {
      await userApi.deleteKey(deleteTarget.value.id!)
      showToast('SSH 密钥已删除', 'success')
      await loadSSHKeys()
    }
    closeDeleteModal()
  } catch (error: any) {
    console.error('删除失败:', error)
    showToast(error?.message || '删除失败', 'error')
  } finally {
    deleting.value = false
  }
}

// 打开添加密钥模态框
const openAddKeyModal = () => {
  keyForm.value = { name: '', content: '' }
  keyErrors.value = {}
  showKeyModal.value = true
}

// 关闭密钥模态框
const closeKeyModal = () => {
  showKeyModal.value = false
  keyForm.value = { name: '', content: '' }
  keyErrors.value = {}
}

// 保存密钥
const saveKey = async () => {
  if (!validateKeyForm()) return

  keySaving.value = true
  try {
    const username = authStore.user?.username || 'admin'
    await userApi.addKey({ ...keyForm.value, user: username })
    showToast('SSH 密钥添加成功', 'success')
    await loadSSHKeys()
    closeKeyModal()
  } catch (error: any) {
    console.error('添加密钥失败:', error)
    showToast(error?.message || '添加密钥失败', 'error')
  } finally {
    keySaving.value = false
  }
}

onMounted(() => {
  loadUsers()
  loadSSHKeys()
})
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.2s ease, opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.95);
  opacity: 0;
}

.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(20px);
}
</style>
