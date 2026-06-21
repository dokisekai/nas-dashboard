import { useNotificationStore } from '../stores/notification'
import { systemApi } from '../api'

export function usePowerOperations() {
  const notificationStore = useNotificationStore()

  const notify = (type: 'warning' | 'error', title: string, message: string) =>
    notificationStore.add({ type, title, message })

  const reboot = async () => {
    if (!confirm('确定要重启系统吗？所有未保存的工作将丢失！')) return
    try {
      await notify('warning', '系统重启', '系统正在准备重启，请稍候...')
      await systemApi.restart()
    } catch (error) {
      await notify('error', '重启失败', `无法触发系统重启：${(error as Error).message}`)
    }
  }

  const shutdown = async () => {
    if (!confirm('确定要关机吗？')) return
    try {
      await notify('warning', '系统关机', '系统正在准备关机，请稍候...')
      await systemApi.shutdown()
    } catch (error) {
      await notify('error', '关机失败', `无法触发系统关机：${(error as Error).message}`)
    }
  }

  return { reboot, shutdown }
}
