<template>
  <div class="sparkline-container">
    <svg
      :width="width"
      :height="height"
      :viewBox="`0 0 ${width} ${height}`"
      preserveAspectRatio="none"
    >
      <path
        :d="pathData"
        :stroke="color"
        :stroke-width="strokeWidth"
        fill="none"
        vector-effect="non-scaling-stroke"
      />
    </svg>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  data: number[]
  color?: string
  width?: number
  height?: number
  strokeWidth?: number
}

const props = withDefaults(defineProps<Props>(), {
  color: '#3b82f6',
  width: 100,
  height: 30,
  strokeWidth: 2
})

const pathData = computed(() => {
  if (!props.data || props.data.length === 0) return ''

  const max = Math.max(...props.data)
  const min = Math.min(...props.data)
  const range = max - min || 1

  const stepX = props.width / (props.data.length - 1)
  const scaleY = (props.height - props.strokeWidth) / range

  let path = `M 0 ${props.height - ((props.data[0] - min) * scaleY)}`

  for (let i = 1; i < props.data.length; i++) {
    const x = i * stepX
    const y = props.height - ((props.data[i] - min) * scaleY)
    path += ` L ${x} ${y}`
  }

  return path
})
</script>

<style scoped>
.sparkline-container {
  display: inline-block;
}

svg {
  display: block;
}
</style>