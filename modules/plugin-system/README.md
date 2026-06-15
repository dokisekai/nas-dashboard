# 插件系统模块

## 概述

插件系统模块提供完整的插件架构，支持第三方开发者扩展系统功能，包括前端插件、后端插件和完整应用包的开发和分发。

## 功能特性

### 插件架构
- **前端插件系统** - Vue 3组件插件
- **后端插件API** - Go服务插件
- **应用包系统** - 完整应用分发
- **插件市场** - 插件发现和安装

### 插件管理
- 插件安装和卸载
- 插件版本管理
- 插件依赖处理
- 插件权限控制
- 插件状态监控

### 开发工具
- 插件开发模板
- SDK和API文档
- 调试和测试工具
- 打包和发布工具

## 后端API实现

### 主要API文件
- `backend/internal/api/plugin_management.go` - 插件管理API
- `backend/pkg/application/` - 应用包管理
- `frontend/src/plugin-system/` - 前端插件系统

### 前端插件系统
- `frontend/src/plugin-system/README.md` - 插件系统说明
- `frontend/src/plugin-system/docs/API_REFERENCE.md` - API参考文档
- `frontend/src/plugin-system/sdk/` - 插件SDK
- `frontend/src/plugin-system/core/` - 插件核心功能
- `frontend/src/plugin-system/marketplace/` - 插件市场

## API端点

### 插件管理
- `GET /api/plugins` - 获取插件列表
- `POST /api/plugins/install` - 安装插件
- `DELETE /api/plugins/:id` - 卸载插件
- `PUT /api/plugins/:id/config` - 配置插件
- `GET /api/plugins/:id/status` - 获取插件状态
- `POST /api/plugins/:id/enable` - 启用插件
- `POST /api/plugins/:id/disable` - 禁用插件

### 应用包管理
- `GET /api/apps` - 获取应用列表
- `POST /api/apps/upload` - 上传应用包
- `POST /api/apps/:id/install` - 安装应用
- `DELETE /api/apps/:id` - 卸载应用
- `GET /api/apps/available` - 获取可用应用

## 插件架构

### 前端插件系统
```typescript
// 插件结构
interface Plugin {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  icon?: string;
  permissions: string[];
  dependencies: string[];
  component: Component;
  config?: PluginConfig;
}

// 插件加载器
class PluginLoader {
  async load(plugin: Plugin): Promise<void>
  async unload(pluginId: string): Promise<void>
  async getPlugin(pluginId: string): Promise<Plugin>
}

// 插件上下文
interface PluginContext {
  storage: PluginStorage;
  api: PluginAPI;
  logger: PluginLogger;
  utils: PluginUtils;
}
```

### 后端插件API
```go
// 插件接口
type Plugin interface {
    ID() string
    Name() string
    Version() string
    Init(ctx context.Context, config map[string]interface{}) error
    Start() error
    Stop() error
    Status() PluginStatus
    Execute(action string, params map[string]interface{}) (interface{}, error)
}

// 插件管理器
type PluginManager struct {
    plugins map[string]Plugin
    configs map[string]PluginConfig
    logger  *log.Logger
}
```

## 插件开发

### 前端插件开发
```typescript
// plugin.ts
import { definePlugin } from '@nas/plugin-system';

export default definePlugin({
  id: 'my-plugin',
  name: 'My Plugin',
  version: '1.0.0',
  description: 'A sample plugin',
  
  // 插件组件
  component: () => import('./MyComponent.vue'),
  
  // 插件配置
  config: {
    settings: [
      {
        key: 'apiKey',
        type: 'text',
        label: 'API Key',
        required: true
      }
    ]
  },
  
  // 插件初始化
  init: async (context) => {
    const { api, storage, logger } = context;
    
    // 存储初始化
    await storage.setItem('initialized', true);
    
    logger.info('Plugin initialized');
  },
  
  // 插件生命周期
  mount: () => {
    console.log('Plugin mounted');
  },
  
  unmount: () => {
    console.log('Plugin unmounted');
  }
});
```

### 后端插件开发
```go
package main

import (
    "context"
    "github.com/nas-dashboard/plugin-sdk"
)

type MyPlugin struct {
    config map[string]interface{}
    logger *log.Logger
}

func (p *MyPlugin) ID() string {
    return "my-plugin"
}

func (p *MyPlugin) Name() string {
    return "My Plugin"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Init(ctx context.Context, config map[string]interface{}) error {
    p.config = config
    p.logger = ctx.Value("logger").(*log.Logger)
    return nil
}

func (p *MyPlugin) Start() error {
    p.logger.Info("Plugin started")
    return nil
}

func (p *MyPlugin) Stop() error {
    p.logger.Info("Plugin stopped")
    return nil
}

func (p *MyPlugin) Status() plugin.Status {
    return plugin.Status{
        Running: true,
        Health:  "healthy",
    }
}

func (p *MyPlugin) Execute(action string, params map[string]interface{}) (interface{}, error) {
    switch action {
    case "getData":
        return p.getData(params)
    default:
        return nil, fmt.Errorf("unknown action: %s", action)
    }
}

// 导出插件
func New() plugin.Plugin {
    return &MyPlugin{}
}
```

## 插件打包

### 前端插件包
```
my-plugin.spk
├── manifest.json       # 插件清单
├── icon.png           # 插件图标
├── plugin.js          # 插件代码
├── styles.css         # 样式文件
└── assets/           # 资源文件
```

### 完整应用包
```
my-app.spk
├── INFO.txt           # 应用信息
├── manifest.json      # 应用清单
├── installer.sh       # 安装脚本
├── app/               # 应用文件
│   ├── backend/      # 后端代码
│   ├── frontend/     # 前端代码
│   └── config/       # 配置文件
└── scripts/          # 管理脚本
    ├── start.sh
    ├── stop.sh
    └── status.sh
```

### manifest.json示例
```json
{
  "id": "my-plugin",
  "name": "My Plugin",
  "version": "1.0.0",
  "description": "A sample plugin for NAS Dashboard",
  "author": "Your Name",
  "license": "MIT",
  "type": "frontend",
  "entry": "plugin.js",
  "icon": "icon.png",
  "permissions": [
    "storage.read",
    "storage.write",
    "network.request"
  ],
  "dependencies": {
    "nas-dashboard": ">=1.0.0"
  },
  "config": {
    "settings": [
      {
        "key": "apiKey",
        "type": "text",
        "label": "API Key",
        "required": true,
        "encrypted": true
      },
      {
        "key": "updateInterval",
        "type": "number",
        "label": "Update Interval",
        "default": 60,
        "min": 10,
        "max": 3600
      }
    ]
  }
}
```

## SDK使用

### 前端SDK
```typescript
import { usePluginContext } from '@nas/plugin-system';

// 在插件组件中使用
export default {
  setup() {
    const { api, storage, logger, utils } = usePluginContext();
    
    // API调用
    const fetchData = async () => {
      const data = await api.get('/api/data');
      return data;
    };
    
    // 存储操作
    const saveData = async (key: string, value: any) => {
      await storage.setItem(key, value);
    };
    
    // 日志记录
    logger.info('Data saved successfully');
    
    // 工具函数
    const formatted = utils.formatDate(new Date());
    
    return {
      fetchData,
      saveData,
      formatted
    };
  }
};
```

### 后端SDK
```go
package main

import "github.com/nas-dashboard/plugin-sdk"

// 使用插件SDK
type MyPlugin struct {
    sdk *plugin.SDK
}

func (p *MyPlugin) Init(ctx context.Context, config map[string]interface{}) error {
    p.sdk = plugin.NewSDK(ctx, config)
    
    // 使用SDK功能
    p.sdk.Log().Info("Initializing plugin")
    
    // 存储访问
    if err := p.sdk.Storage().Set("key", "value"); err != nil {
        return err
    }
    
    // API调用
    result, err := p.sdk.API().Get("/api/system/info")
    if err != nil {
        return err
    }
    
    p.sdk.Log().Info("System info: %v", result)
    return nil
}
```

## 插件市场

### 市场功能
- 插件浏览和搜索
- 分类和标签过滤
- 评分和评论
- 安装统计
- 版本历史

### 市场API
```typescript
// 获取插件列表
const plugins = await marketplace.getPlugins({
    category: 'monitoring',
    query: 'temperature',
    sort: 'popularity'
});

// 获取插件详情
const plugin = await marketplace.getPluginDetail('temp-monitor');

// 安装插件
await marketplace.installPlugin('temp-monitor', '1.0.0');

// 获取插件评分
const rating = await marketplace.getPluginRating('temp-monitor');
```

## 配置文件

### 插件系统配置
```json
{
  "enabled": true,
  "pluginDirectory": "/var/lib/nas-plugins",
  "tempDirectory": "/tmp/nas-plugins",
  "maxPlugins": 100,
  "permissions": {
    "storage.read": {
      "description": "Read storage information",
      "dangerous": false
    },
    "storage.write": {
      "description": "Write to storage",
      "dangerous": true
    },
    "network.request": {
      "description": "Make network requests",
      "dangerous": false
    },
    "system.config": {
      "description": "Modify system configuration",
      "dangerous": true
    }
  },
  "marketplace": {
    "enabled": true,
    "url": "https://plugins.nas-dashboard.com",
    "autoUpdate": true,
    "updateInterval": 86400
  }
}
```

## 安全考虑

### 插件权限
- 最小权限原则
- 权限显式声明
- 运行时权限检查
- 危险操作确认

### 沙箱隔离
- 前端iframe隔离
- 后端容器隔离
- 资源限制配置
- 网络访问控制

### 代码签名
- 插件包签名验证
- 作者身份验证
- 完整性校验
- 恶意代码检测

## 监控和日志

### 插件监控
- 插件运行状态
- 资源使用统计
- 错误率监控
- 性能指标收集

### 插件日志
```typescript
// 插件日志API
logger.debug('Debug message');
logger.info('Info message');
logger.warn('Warning message');
logger.error('Error message');

// 结构化日志
logger.info('Plugin action', {
    action: 'user_login',
    userId: 123,
    timestamp: new Date()
});
```

## 开发示例

### 数据展示插件
```typescript
// temp-display.ts
import { definePlugin, ref, onMounted } from '@nas/plugin-system';

export default definePlugin({
    id: 'temp-display',
    name: 'Temperature Display',
    version: '1.0.0',
    
    component: {
        setup() {
            const { api, storage, logger } = usePluginContext();
            const temperature = ref(0);
            
            const fetchTemperature = async () => {
                try {
                    const data = await api.get('/api/monitor/temperature');
                    temperature.value = data.cpu;
                } catch (error) {
                    logger.error('Failed to fetch temperature', error);
                }
            };
            
            onMounted(() => {
                fetchTemperature();
                setInterval(fetchTemperature, 5000);
            });
            
            return { temperature };
        },
        
        template: `
            <div class="temp-display">
                <h3>Temperature</h3>
                <p>{{ temperature }}°C</p>
            </div>
        `
    }
});
```

## 故障排除

### 插件加载失败
1. 检查插件清单：`cat manifest.json`
2. 验证依赖项：`npm list`
3. 查看控制台错误：浏览器开发者工具
4. 检查插件权限配置

### 插件权限问题
1. 验证权限声明：`manifest.json`
2. 检查权限授予：插件设置页面
3. 查看权限日志：系统日志
4. 重新授权插件

### 插件性能问题
1. 监控资源使用：浏览器任务管理器
2. 检查内存泄漏：Chrome Memory Profiler
3. 优化代码逻辑：代码审查
4. 限制插件权限：权限控制

## 相关文档

- [插件开发指南](../../docs/PLUGIN_DEVELOPMENT.md)
- [前端插件系统](../../frontend/src/plugin-system/README.md)
- [插件API参考](../../frontend/src/plugin-system/docs/API_REFERENCE.md)
- [应用系统模块](../application-system/)

## DSM对标功能

本模块对标Synology DSM的插件和包管理功能：
- ✅ 插件架构设计
- ✅ 应用包管理
- ✅ 插件市场概念
- ✅ SDK和开发工具
- ✅ 权限管理系统

功能完整度: 85% ✅