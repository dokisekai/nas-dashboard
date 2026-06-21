import { ref, onMounted, onUnmounted } from 'vue'

export interface ClockOptions {
  intervalMs?: number
  locale?: string
  options?: Intl.DateTimeFormatOptions
}

export function useClock(options: ClockOptions = {}) {
  const {
    intervalMs = 1000,
    locale = 'zh-CN',
    options: formatOptions = { hour: '2-digit', minute: '2-digit' },
  } = options

  const currentTime = ref('')
  let timer: number | undefined

  const update = () => {
    currentTime.value = new Date().toLocaleTimeString(locale, formatOptions)
  }

  onMounted(() => {
    update()
    timer = window.setInterval(update, intervalMs)
  })

  onUnmounted(() => {
    if (timer) clearInterval(timer)
  })

  return { currentTime }
}
