import { onUnmounted } from 'vue'

export function useClickOutside(
  containerRef: () => HTMLElement | null,
  onClose: (event: MouseEvent) => void,
) {
  const handler = (event: MouseEvent) => {
    const el = containerRef()
    if (el && !el.contains(event.target as Node)) {
      onClose(event)
      document.removeEventListener('click', handler)
    }
  }

  const start = () => document.addEventListener('click', handler)
  const stop = () => document.removeEventListener('click', handler)

  onUnmounted(stop)

  return { start, stop, handler }
}
