# 桌面系统集成指南

## 快速开始

### 1. 基本使用

在主应用中引入桌面系统：

```vue
<template>
  <DSMDesktop />
</template>

<script setup>
import DSMDesktop from './components/Desktop/DSMDesktop.vue'
</script>
```

### 2. 状态管理集成

如果需要与应用状态集成：

```javascript
import { useAppStore } from './stores/app'

const appStore = useAppStore()

// 桌面会自动同步应用状态
// 包括：窗口状态、小部件配置、主题设置等
```

### 3. 路由集成

如需在多页面应用中使用：

```javascript
// router/index.js
{
  path: '/desktop',
  name: 'Desktop',
  component: () => import('@/components/Desktop/DSMDesktop.vue')
}
```

## 配置项详解

### 桌面全局配置

```javascript
// src/config/desktop.config.js
export default {
  // 默认小部件
  defaultWidgets: [
    {
      id: 'system-monitor',
      type: 'system-monitor',
      position: { x: 20, y: 20 },
      size: 'large',
      config: { showCpu: true, showMemory: true, showDisk: true }
    }
  ],

  // 默认Dock项
  defaultDockItems: [
    {
      id: 'app-center',
      label: '应用中心',
      icon: 'GridIcon',
      appId: 'app-center'
    }
  ],

  // 主题设置
  defaultTheme: {
    background: {
      id: 'gradient1',
      type: 'builtin'
    },
    theme: {
      id: 'light',
      mode: 'auto'
    }
  },

  // 功能开关
  features: {
    windowSnap: true,        // 窗口吸附
    dockMagnification: true, // Dock放大效果
    widgetDrag: true,        // 小部件拖拽
    keyboardShortcuts: true  // 键盘快捷键
  }
}
```

## 自定义扩展

### 1. 添加自定义小部件

```vue
<!-- src/components/Desktop/widgets/CustomWidget.vue -->
<template>
  <div class="custom-widget">
    <!-- 你的小部件内容 -->
  </div>
</template>

<script setup>
// 定义小部件配置
const props = defineProps({
  config: Object,
  size: String
})

// 小部件逻辑
</script>

<style scoped>
.custom-widget {
  /* 小部件样式 */
}
</style>
```

在小部件库中注册：

```javascript
// src/components/Desktop/WidgetLibrary.vue
const widgetTemplates = [
  // ... 现有小部件
  {
    id: 'custom-widget',
    type: 'custom-widget',
    name: '自定义小部件',
    description: '我的自定义小部件',
    category: 'custom',
    defaultSize: 'medium',
    icon: 'CustomIcon'
  }
]
```

### 2. 添加自定义应用

```javascript
// src/config/apps.config.js
export const customApps = [
  {
    id: 'my-app',
    name: '我的应用',
    component: 'MyAppComponent',
    icon: 'MyAppIcon',
    windowConfig: {
      width: 800,
      height: 600,
      resizable: true,
      minimizable: true,
      maximizable: true
    }
  }
]
```

### 3. 自定义主题

```javascript
// src/config/themes.config.js
export const customThemes = [
  {
    id: 'my-theme',
    name: '我的主题',
    description: '自定义主题',
    colors: ['#primary', '#secondary', '#accent', '#background'],
    css: `
      :root {
        --primary-color: #primary;
        --secondary-color: #secondary;
      }
    `
  }
]
```

## API 集成

### 1. 系统监控 API

```javascript
// src/api/system.js
export const systemApi = {
  // 获取系统信息
  async getSystemInfo() {
    const response = await fetch('/api/system/info')
    return response.json()
  },

  // 获取实时监控数据
  async getMonitorData() {
    const response = await fetch('/api/system/monitor')
    return response.json()
  }
}

// 在小部件中使用
import { systemApi } from '@/api/system'

const updateMonitor = async () => {
  const data = await systemApi.getMonitorData()
  // 更新小部件显示
}
```

### 2. 天气 API 集成

```javascript
// src/api/weather.js
export const weatherApi = {
  async getWeather(location = '北京') {
    const response = await fetch(`/api/weather?location=${location}`)
    return response.json()
  }
}
```

### 3. 存储状态 API

```javascript
// src/api/storage.js
export const storageApi = {
  async getStorageInfo() {
    const response = await fetch('/api/storage/info')
    return response.json()
  },

  async getVolumeStats(volume = '/') {
    const response = await fetch(`/api/storage/volume/${volume}`)
    return response.json()
  }
}
```

## 事件系统

### 桌面事件监听

```javascript
import { eventBus } from '@/utils/eventBus'

// 监听窗口事件
eventBus.on('window:opened', (window) => {
  console.log('窗口已打开:', window)
})

eventBus.on('window:closed', (windowId) => {
  console.log('窗口已关闭:', windowId)
})

eventBus.on('widget:added', (widget) => {
  console.log('小部件已添加:', widget)
})

eventBus.on('theme:changed', (theme) => {
  console.log('主题已更改:', theme)
})
```

### 自定义事件

```javascript
// 触发自定义事件
eventBus.emit('custom:event', { data: 'value' })

// 监听自定义事件
eventBus.on('custom:event', (data) => {
  console.log('自定义事件:', data)
})
```

## 数据持久化

### LocalStorage 配置

```javascript
// src/utils/storage.js
export const desktopStorage = {
  // 保存桌面配置
  saveConfig(config) {
    localStorage.setItem('desktop-config', JSON.stringify(config))
  },

  // 加载桌面配置
  loadConfig() {
    const config = localStorage.getItem('desktop-config')
    return config ? JSON.parse(config) : null
  },

  // 导出配置
  exportConfig() {
    const config = this.loadConfig()
    return JSON.stringify(config, null, 2)
  },

  // 导入配置
  importConfig(configString) {
    try {
      const config = JSON.parse(configString)
      this.saveConfig(config)
      return true
    } catch (error) {
      console.error('配置导入失败:', error)
      return false
    }
  }
}
```

### 服务器同步（可选）

```javascript
// src/api/sync.js
export const syncApi = {
  // 上传配置到服务器
  async uploadConfig(config) {
    const response = await fetch('/api/sync/config', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(config)
    })
    return response.json()
  },

  // 从服务器下载配置
  async downloadConfig() {
    const response = await fetch('/api/sync/config')
    return response.json()
  }
}
```

## 性能优化

### 1. 懒加载配置

```javascript
// 路由懒加载
const DesktopSettings = () => import('./components/Desktop/DesktopSettings.vue')

// 组件懒加载
const WidgetLibrary = defineAsyncComponent(() =>
  import('./components/Desktop/WidgetLibrary.vue')
)
```

### 2. 虚拟化长列表

```javascript
// 使用虚拟滚动优化大量小部件
import { RecycleScroller } from 'vue-virtual-scroller'

<RecycleScroller
  :items="widgets"
  :item-size="200"
  key-field="id"
>
  <template #default="{ item }">
    <DesktopWidget :widget="item" />
  </template>
</RecycleScroller>
```

### 3. 防抖和节流

```javascript
// src/utils/performance.js
export const debounce = (fn, delay) => {
  let timeoutId
  return (...args) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

export const throttle = (fn, limit) => {
  let inThrottle
  return (...args) => {
    if (!inThrottle) {
      fn(...args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

// 使用示例
const handleResize = debounce(() => {
  // 处理窗口大小变化
}, 200)
```

## 测试指南

### 单元测试示例

```javascript
// tests/components/DesktopWidget.spec.js
import { mount } from '@vue/test-utils'
import DesktopWidget from '@/components/Desktop/DesktopWidget.vue'

describe('DesktopWidget', () => {
  it('renders widget correctly', () => {
    const wrapper = mount(DesktopWidget, {
      props: {
        widget: {
          id: 'test-widget',
          type: 'clock',
          position: { x: 100, y: 100 },
          size: 'medium',
          config: {}
        }
      }
    })
    expect(wrapper.find('.desktop-widget').exists()).toBe(true)
  })

  it('emits drag-start event on drag', async () => {
    const wrapper = mount(DesktopWidget, {
      props: { widget: mockWidget }
    })

    await wrapper.trigger('dragstart')
    expect(wrapper.emitted('drag-start')).toBeTruthy()
  })
})
```

### E2E 测试示例

```javascript
// tests/e2e/desktop.spec.js
describe('Desktop E2E', () => {
  it('should open application from dock', () => {
    cy.visit('/')
    cy.get('.dock-item').first().click()
    cy.get('.desktop-window').should('be.visible')
  })

  it('should add widget from library', () => {
    cy.visit('/')
    cy.get('body').rightclick()
    cy.contains('添加小部件').click()
    cy.get('.widget-card').first().find('.add-btn').click()
    cy.get('.widgets-grid').children().should('have.length.greaterThan', 0)
  })
})
```

## 故障排除

### 常见问题

1. **小部件不显示**
   - 检查小部件组件是否正确导入
   - 确认小部件配置格式正确
   - 查看浏览器控制台错误信息

2. **窗口拖拽不工作**
   - 确认没有其他元素覆盖窗口
   - 检查 z-index 设置
   - 验证事件监听器是否正常

3. **Dock图标不显示**
   - 检查图标组件是否正确导入
   - 确认 Heroicons 图标库已安装
   - 验证图标名称拼写正确

4. **主题不生效**
   - 清除浏览器缓存
   - 检查 LocalStorage 中的配置
   - 确认 CSS 变量正确应用

### 调试模式

```javascript
// 启用调试模式
localStorage.setItem('desktop-debug', 'true')

// 查看当前配置
console.log(JSON.parse(localStorage.getItem('desktop-config')))

// 重置桌面配置
localStorage.removeItem('desktop-config')
location.reload()
```

## 部署建议

### 1. 生产环境配置

```javascript
// vue.config.js
module.exports = {
  productionSourceMap: false,
  configureWebpack: {
    optimization: {
      splitChunks: {
        chunks: 'all',
        maxSize: 244 * 1024 // 244KB
      }
    }
  }
}
```

### 2. 性能监控

```javascript
// src/utils/monitoring.js
export const performanceMonitor = {
  init() {
    if (process.env.NODE_ENV === 'production') {
      // 监控页面加载性能
      window.addEventListener('load', () => {
        const perfData = performance.timing
        const pageLoadTime = perfData.loadEventEnd - perfData.navigationStart
        console.log('页面加载时间:', pageLoadTime)
      })
    }
  }
}
```

### 3. 错误处理

```javascript
// src/utils/errorHandler.js
export const errorHandler = {
  init() {
    window.addEventListener('error', (event) => {
      console.error('全局错误:', event.error)
      // 发送错误报告到服务器
    })

    window.addEventListener('unhandledrejection', (event) => {
      console.error('未处理的Promise拒绝:', event.reason)
      // 发送错误报告到服务器
    })
  }
}
```

这个完整的桌面系统已经可以投入使用，提供了完整的桌面体验和丰富的自定义选项。