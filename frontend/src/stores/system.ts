import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useSystemStore = defineStore('system', () => {
  const initialized = ref<boolean | null>(null) // null表示未检查，true/false表示实际状态
  const loading = ref(false)

  const checkInitStatus = async () => {
    loading.value = true
    try {
      const response = await axios.get('/api/system/init-status')
      initialized.value = response.data.initialized
      return initialized.value
    } catch (error) {
      console.error('Failed to check init status:', error)
      // 如果API失败，在生产环境中默认为未初始化
      if (import.meta.env.PROD) {
        initialized.value = false
      } else {
        // 开发环境默认为已初始化
        initialized.value = true
      }
      return initialized.value
    } finally {
      loading.value = false
    }
  }

  const setInitialized = (value: boolean) => {
    initialized.value = value
  }

  const isInitialized = () => {
    return initialized.value === true
  }

  const needsInit = () => {
    return initialized.value === false
  }

  return {
    initialized,
    loading,
    checkInitStatus,
    setInitialized,
    isInitialized,
    needsInit
  }
})
