<template>
  <div ref="chartContainer" class="w-full h-full"></div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, watch, ref } from 'vue'
import ApexCharts, { type ApexOptions } from 'apexcharts'

interface ChartOptions {
  chart?: {
    type?: string
    height?: number
    animations?: any
    toolbar?: { show?: boolean }
    zoom?: { enabled?: boolean }
  }
  colors?: string[]
  fill?: any
  dataLabels?: { enabled?: boolean }
  stroke?: {
    show?: boolean
    curve?: string | string[]
    width?: number | number[]
    colors?: string[]
    dashArray?: number | number[]
  }
  xaxis?: {
    categories?: string[]
    labels?: any
    axisBorder?: any
    axisTicks?: any
  }
  yaxis?: any
  grid?: any
  tooltip?: any
  series?: any[]
  labels?: string[]
  plotOptions?: any
  legend?: any
}

interface Props {
  type?: string
  options?: ChartOptions
  series?: any[]
  height?: number
}

const props = withDefaults(defineProps<Props>(), {
  type: 'line',
  height: 300,
})

const chartContainer = ref<HTMLElement>()
let chart: ApexCharts | null = null

const initChart = () => {
  if (!chartContainer.value) return

  const defaultOptions: ChartOptions = {
    chart: {
      type: props.type,
      height: props.height,
      toolbar: { show: false },
      zoom: { enabled: false },
    },
    grid: {
      borderColor: '#374151',
      strokeDashArray: 4,
    },
    tooltip: { theme: 'dark' },
  }

  const mergedOptions: ApexOptions = {
    ...defaultOptions,
    ...props.options,
    chart: {
      ...defaultOptions.chart,
      ...(props.options?.chart as any),
    },
  }

  chart = new ApexCharts(chartContainer.value, {
    ...mergedOptions,
    series: props.series || [],
  })

  chart.render()
}

// 监听 series 变化
watch(() => props.series, (newSeries) => {
  if (chart && newSeries) {
    chart.updateSeries(newSeries)
  }
}, { deep: true })

// 监听 options 变化
watch(() => props.options, (newOptions) => {
  if (chart && newOptions) {
    chart.updateOptions(newOptions as ApexOptions)
  }
}, { deep: true })

onMounted(() => {
  initChart()
})

onUnmounted(() => {
  if (chart) {
    chart.destroy()
    chart = null
  }
})
</script>
