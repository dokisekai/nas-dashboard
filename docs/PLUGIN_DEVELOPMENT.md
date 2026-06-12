# NAS Dashboard - Plugin Development Guide

Complete guide for developing plugins for the NAS Dashboard, including architecture, SDK reference, examples, and best practices.

## 📋 Table of Contents

1. [Plugin System Overview](#plugin-system-overview)
2. [Getting Started](#getting-started)
3. [Plugin Architecture](#plugin-architecture)
4. [SDK Reference](#sdk-reference)
5. [Plugin Development](#plugin-development)
6. [Plugin Lifecycle](#plugin-lifecycle)
7. [Permission System](#permission-system)
8. [Examples](#examples)
9. [Testing Plugins](#testing-plugins)
10. [Publishing Plugins](#publishing-plugins)
11. [Best Practices](#best-practices)

---

## Plugin System Overview

### What is a Plugin?

A plugin is a self-contained module that extends the functionality of the NAS Dashboard without modifying the core codebase. Plugins can:

- Add new desktop widgets
- Create custom applications
- Integrate external services
- Add new monitoring capabilities
- Extend the UI with custom components
- Provide custom data sources

### Plugin Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    NAS Dashboard Core                       │
│  ┌──────────────────────────────────────────────────────┐ │
│  │              Plugin Loader                             │ │
│  │  - Dynamic Loading                                     │ │
│  │  - Lifecycle Management                               │ │
│  │  - Dependency Resolution                              │ │
│  └──────────────────────────────────────────────────────┘ │
│  ┌──────────────────────────────────────────────────────┐ │
│  │              Plugin SDK                                │ │
│  │  - API Access                                         │ │
│  │  - Storage System                                     │ │
│  │  - UI Integration                                     │ │
│  │  - WebSocket Communication                            │ │
│  └──────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ Plugin Interface
                            │
┌─────────────────────────────────────────────────────────────┐
│                     Plugin Ecosystem                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Widget     │  │ Application  │  │   Service     │     │
│  │   Plugin     │  │   Plugin     │  │   Plugin      │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   Monitor    │  │ Integration  │  │   Custom     │     │
│  │   Plugin     │  │   Plugin     │  │   Plugin      │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

- **Dynamic Loading**: Load plugins at runtime without restart
- **Sandboxed Execution**: Isolated execution context for security
- **Lifecycle Management**: Install, enable, disable, uninstall plugins
- **Permission System**: Granular control over plugin capabilities
- **Hot Reload**: Development-friendly hot reload support
- **Type Safety**: Full TypeScript support
- **SDK**: Comprehensive API for plugin development

---

## Getting Started

### Prerequisites

- Node.js 20+ installed
- TypeScript knowledge
- Vue 3 understanding
- Basic JavaScript/ES6 knowledge

### Quick Start Example

```typescript
// my-first-plugin.ts
import { definePlugin } from '@nas-dashboard/sdk';

export const manifest = {
  id: 'my-first-plugin',
  name: 'My First Plugin',
  version: '1.0.0',
  description: 'A simple example plugin',
  author: 'Your Name',
  permissions: ['storage'] as const,
};

export default definePlugin((context) => {
  const { storage, logger } = context;

  return {
    async onInstall() {
      logger.info('Plugin installed!');
      await storage.set('initialized', true);
    },

    async onActivate() {
      logger.info('Plugin activated!');
      const data = await storage.get('initialized');
      logger.info('Stored data:', data);
    },

    async onDeactivate() {
      logger.info('Plugin deactivated!');
    }
  };
});
```

### Plugin Structure

```
my-plugin/
├── src/
│   ├── index.ts          # Main plugin file
│   ├── components/       # Vue components
│   ├── widgets/          # Widget components
│   ├── utils/           # Helper functions
│   └── types/           # TypeScript types
├── public/              # Static assets
├── package.json         # NPM configuration
├── tsconfig.json        # TypeScript configuration
└── README.md            # Plugin documentation
```

---

## Plugin Architecture

### Plugin Components

#### 1. Manifest

The manifest defines plugin metadata:

```typescript
export const manifest = {
  id: 'unique-plugin-id',           // Required: Unique identifier
  name: 'Plugin Name',              // Required: Display name
  version: '1.0.0',                // Required: Semantic version
  description: 'Plugin description', // Required: What it does
  author: 'Author Name',            // Optional: Author information
  license: 'MIT',                   // Optional: License type
  permissions: [                    // Required: Requested permissions
    'storage',
    'ui',
    'network'
  ] as const,
  dependencies: [                   // Optional: Plugin dependencies
    'other-plugin-id'
  ],
  engine: '>=0.1.0'                // Optional: Required engine version
};
```

#### 2. Plugin Function

The main plugin function exports the plugin implementation:

```typescript
export default definePlugin((context: PluginContext) => {
  // Access SDK
  const { storage, logger, api, ui } = context;

  // Return plugin interface
  return {
    // Lifecycle hooks
    async onInstall() { /* ... */ },
    async onActivate() { /* ... */ },
    async onDeactivate() { /* ... */ },
    async onUninstall() { /* ... */ },
    
    // Optional: UI registration
    widgets: [/* widget definitions */],
    apps: [/* app definitions */],
    
    // Optional: Custom API
    api: { /* custom API endpoints */ }
  };
});
```

#### 3. SDK Context

The SDK context provides access to platform capabilities:

```typescript
interface PluginContext {
  // Storage
  storage: PluginStorage;
  
  // Logging
  logger: PluginLogger;
  
  // API Access
  api: {
    navigate: (path: string) => void;
    request: (config: RequestConfig) => Promise<Response>;
    websocket: WebSocketClient;
  };
  
  // UI Integration
  ui: {
    registerWidget: (widget: WidgetDefinition) => void;
    registerApp: (app: AppDefinition) => void;
    notify: (message: string, type: NotificationType) => void;
  };
  
  // Utilities
  utils: {
    debounce: (fn: Function, delay: number) => Function;
    throttle: (fn: Function, limit: number) => Function;
    merge: (...objects: any[]) => any;
    clone: <T>(obj: T) => T;
  };
  
  // Plugin Info
  plugin: {
    id: string;
    version: string;
    dataPath: string;
  };
}
```

---

## SDK Reference

### Storage API

The storage API provides persistent key-value storage:

```typescript
interface PluginStorage {
  // Store a value
  set(key: string, value: any): Promise<void>;
  
  // Retrieve a value
  get(key: string): Promise<any>;
  
  // Check if key exists
  has(key: string): Promise<boolean>;
  
  // Delete a value
  delete(key: string): Promise<void>;
  
  // Clear all storage
  clear(): Promise<void>;
  
  // Get all keys
  keys(): Promise<string[]>;
  
  // Get storage size
  size(): Promise<number>;
}
```

**Example:**

```typescript
export default definePlugin(({ storage }) => {
  return {
    async onActivate() {
      // Store data
      await storage.set('user-preferences', {
        theme: 'dark',
        language: 'en'
      });
      
      // Retrieve data
      const prefs = await storage.get('user-preferences');
      console.log('Theme:', prefs.theme);
      
      // Check existence
      const hasPrefs = await storage.has('user-preferences');
      
      // Get all keys
      const keys = await storage.keys();
      
      // Get storage size
      const size = await storage.size();
      
      // Delete specific key
      await storage.delete('user-preferences');
      
      // Clear all
      await storage.clear();
    }
  };
});
```

### Logger API

The logger API provides plugin-specific logging:

```typescript
interface PluginLogger {
  debug(message: string, ...args: any[]): void;
  info(message: string, ...args: any[]): void;
  warn(message: string, ...args: any[]): void;
  error(message: string, ...args: any[]): void;
}
```

**Example:**

```typescript
export default definePlugin(({ logger }) => {
  return {
    async onActivate() {
      logger.debug('Debug message with details:', { some: 'data' });
      logger.info('Plugin activated successfully!');
      logger.warn('This is a warning message');
      logger.error('Something went wrong:', new Error('details'));
    }
  };
});
```

### Navigation API

The navigation API allows programmatic navigation:

```typescript
export default definePlugin(({ api }) => {
  return {
    async onActivate() {
      // Navigate to a route
      api.navigate('/dashboard');
      
      // Navigate with parameters
      api.navigate('/storage?path=/mnt/data');
    }
  };
});
```

### Request API

The request API allows making HTTP requests:

```typescript
export default definePlugin(({ api }) => {
  return {
    async onActivate() {
      try {
        const response = await api.request({
          url: 'https://api.example.com/data',
          method: 'GET',
          headers: {
            'Authorization': 'Bearer token'
          }
        });
        
        const data = await response.json();
        console.log('API Response:', data);
      } catch (error) {
        console.error('Request failed:', error);
      }
    }
  };
});
```

### UI Integration API

The UI API allows registering widgets and apps:

```typescript
export default definePlugin(({ ui }) => {
  return {
    widgets: [
      {
        id: 'custom-widget',
        name: 'Custom Widget',
        component: () => import('./components/CustomWidget.vue'),
        icon: 'chart-bar',
        category: 'monitoring',
        defaultConfig: {
          refreshInterval: 5000
        }
      }
    ],
    
    apps: [
      {
        id: 'custom-app',
        name: 'Custom App',
        component: () => import('./components/CustomApp.vue'),
        icon: 'cog',
        category: 'utilities'
      }
    ]
  };
});
```

### Notification API

The notification API allows showing notifications:

```typescript
export default definePlugin(({ ui }) => {
  return {
    async onActivate() {
      ui.notify('Plugin activated!', 'success');
      ui.notify('Warning message', 'warning');
      ui.notify('Error occurred', 'error');
      ui.notify('Info message', 'info');
    }
  };
});
```

---

## Plugin Development

### Creating a Widget Plugin

```typescript
// custom-widget-plugin.ts
import { definePlugin } from '@nas-dashboard/sdk';
import CustomWidget from './components/CustomWidget.vue';

export const manifest = {
  id: 'custom-widget-plugin',
  name: 'Custom Widget Plugin',
  version: '1.0.0',
  description: 'Adds a custom widget to the desktop',
  permissions: ['storage', 'ui'] as const,
};

export default definePlugin(({ ui, storage }) => {
  return {
    async onInstall() {
      // Initialize widget settings
      await storage.set('widget-settings', {
        title: 'My Custom Widget',
        refreshInterval: 10000
      });
    },

    widgets: [
      {
        id: 'custom-widget',
        name: 'Custom Widget',
        component: CustomWidget,
        icon: 'chart-bar',
        category: 'monitoring',
        description: 'A custom monitoring widget',
        defaultConfig: {
          title: 'Custom Widget',
          refreshInterval: 10000,
          showChart: true
        }
      }
    ]
  };
});
```

**Widget Component (CustomWidget.vue):**

```vue
<template>
  <div class="custom-widget">
    <h3>{{ config.title }}</h3>
    <div class="widget-content">
      <div class="stat">{{ value }}</div>
      <button @click="refresh">Refresh</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

interface WidgetConfig {
  title: string;
  refreshInterval: number;
  showChart: boolean;
}

interface Props {
  config: WidgetConfig;
}

const props = defineProps<Props>();
const value = ref(0);

const refresh = async () => {
  // Simulate data fetch
  value.value = Math.floor(Math.random() * 100);
};

onMounted(() => {
  refresh();
  // Set up refresh interval
  setInterval(refresh, props.config.refreshInterval);
});
</script>

<style scoped>
.custom-widget {
  padding: 1rem;
  background: var(--widget-bg);
  border-radius: 8px;
}

.widget-content {
  margin-top: 1rem;
}

.stat {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 1rem;
}
</style>
```

### Creating an Application Plugin

```typescript
// custom-app-plugin.ts
import { definePlugin } from '@nas-dashboard/sdk';
import CustomApp from './components/CustomApp.vue';

export const manifest = {
  id: 'custom-app-plugin',
  name: 'Custom App Plugin',
  version: '1.0.0',
  description: 'Adds a custom application',
  permissions: ['storage', 'ui', 'network'] as const,
};

export default definePlugin(({ ui, storage, api }) => {
  return {
    apps: [
      {
        id: 'custom-app',
        name: 'Custom App',
        component: CustomApp,
        icon: 'cog',
        category: 'utilities',
        description: 'A custom utility application',
        route: '/custom-app'
      }
    ],

    async onActivate() {
      // Initialize app data
      await storage.set('app-data', {
        items: [],
        settings: {}
      });
    }
  };
});
```

**Application Component (CustomApp.vue):**

```vue
<template>
  <div class="custom-app">
    <header>
      <h2>Custom Application</h2>
      <button @click="addItem">Add Item</button>
    </header>

    <main>
      <div class="items-list">
        <div v-for="item in items" :key="item.id" class="item">
          <span>{{ item.name }}</span>
          <button @click="removeItem(item.id)">Remove</button>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useStorage } from '@nas-dashboard/sdk';

const items = ref<any[]>([]);

const addItem = () => {
  const newItem = {
    id: Date.now(),
    name: `Item ${items.value.length + 1}`
  };
  items.value.push(newItem);
  saveItems();
};

const removeItem = (id: number) => {
  items.value = items.value.filter(item => item.id !== id);
  saveItems();
};

const saveItems = () => {
  // Save to plugin storage
  useStorage().set('app-items', items.value);
};

onMounted(async () => {
  // Load from plugin storage
  const saved = await useStorage().get('app-items');
  if (saved) {
    items.value = saved;
  }
});
</script>

<style scoped>
.custom-app {
  padding: 2rem;
  height: 100%;
  display: flex;
  flex-direction: column;
}

header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.items-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background: var(--item-bg);
  border-radius: 4px;
}
</style>
```

---

## Plugin Lifecycle

### Lifecycle States

```
┌──────────┐
│ Installed │
└─────┬────┘
      │ onInstall()
      ▼
┌──────────┐
│ Disabled  │◄───────────┐
└─────┬────┘              │
      │ enable()          │ disable()
      ▼                   │
┌──────────┐              │
│  Enabled  │──────────────┘
└─────┬────┘
      │ onActivate()
      ▼
┌──────────┐
│  Active  │
└─────┬────┘
      │ onDeactivate()
      ▼
┌──────────┐
│  Enabled  │
└──────────┘

onUninstall()
      ▼
┌──────────┐
│ Removed   │
└──────────┘
```

### Lifecycle Hooks

#### onInstall()

Called when the plugin is first installed:

```typescript
async onInstall() {
  // Initialize storage
  await storage.set('version', '1.0.0');
  await storage.set('settings', defaultSettings);
  
  // Register resources
  logger.info('Plugin installed successfully');
}
```

#### onActivate()

Called when the plugin is activated:

```typescript
async onActivate() {
  // Start services
  startMonitoring();
  
  // Register UI components
  registerWidgets();
  registerApps();
  
  logger.info('Plugin activated');
}
```

#### onDeactivate()

Called when the plugin is deactivated:

```typescript
async onDeactivate() {
  // Stop services
  stopMonitoring();
  
  // Clean up resources
  clearTimers();
  closeConnections();
  
  logger.info('Plugin deactivated');
}
```

#### onUninstall()

Called when the plugin is uninstalled:

```typescript
async onUninstall() {
  // Clean up all data
  await storage.clear();
  
  // Remove registered components
  unregisterWidgets();
  unregisterApps();
  
  logger.info('Plugin uninstalled');
}
```

---

## Permission System

### Available Permissions

```typescript
type Permission =
  | 'storage'      // Access to plugin storage
  | 'ui'           // Register widgets and apps
  | 'network'      // Make HTTP requests
  | 'websocket'    // Access WebSocket
  | 'navigation'   // Navigate routes
  | 'notification' // Show notifications
  | 'system'       // Access system information
  | 'docker'       // Manage Docker containers
  | 'file'         // File system access
  | 'process';     // Process management
```

### Requesting Permissions

```typescript
export const manifest = {
  id: 'my-plugin',
  name: 'My Plugin',
  permissions: [
    'storage',
    'ui',
    'network'
  ] as const
};
```

### Permission Check

The system automatically checks permissions before granting access:

```typescript
export default definePlugin(({ api, logger }) => {
  return {
    async onActivate() {
      try {
        // This will fail if 'network' permission not granted
        const response = await api.request({
          url: 'https://api.example.com/data',
          method: 'GET'
        });
      } catch (error) {
        logger.error('Permission denied or request failed:', error);
      }
    }
  };
});
```

---

## Examples

### Example 1: Weather Widget Plugin

```typescript
// weather-widget-plugin.ts
import { definePlugin } from '@nas-dashboard/sdk';
import WeatherWidget from './components/WeatherWidget.vue';

export const manifest = {
  id: 'weather-widget-plugin',
  name: 'Weather Widget',
  version: '1.0.0',
  description: 'Shows current weather information',
  permissions: ['storage', 'ui', 'network'] as const,
};

export default definePlugin(({ ui, storage, api }) => {
  return {
    async onInstall() {
      await storage.set('weather-settings', {
        location: 'London',
        units: 'metric',
        apiKey: '' // User will provide
      });
    },

    widgets: [
      {
        id: 'weather-widget',
        name: 'Weather Widget',
        component: WeatherWidget,
        icon: 'cloud-sun',
        category: 'information',
        description: 'Current weather conditions',
        defaultConfig: {
          location: 'London',
          units: 'metric'
        }
      }
    ]
  };
});
```

### Example 2: Notification Plugin

```typescript
// notification-plugin.ts
import { definePlugin } from '@nas-dashboard/sdk';

export const manifest = {
  id: 'notification-plugin',
  name: 'Notification Plugin',
  version: '1.0.0',
  description: 'Advanced notification system',
  permissions: ['storage', 'ui', 'websocket'] as const,
};

export default definePlugin(({ storage, websocket }) => {
  let notificationCheck: NodeJS.Timeout | null = null;

  return {
    async onActivate() {
      // Start notification check
      notificationCheck = setInterval(async () => {
        const notifications = await storage.get('notifications');
        if (notifications && notifications.length > 0) {
          // Process notifications
          notifications.forEach((notif: any) => {
            ui.notify(notif.message, notif.type);
          });
        }
      }, 60000); // Check every minute
    },

    async onDeactivate() {
      // Stop notification check
      if (notificationCheck) {
        clearInterval(notificationCheck);
        notificationCheck = null;
      }
    }
  };
});
```

### Example 3: System Monitor Plugin

```typescript
// system-monitor-plugin.ts
import { definePlugin } from '@nas-dashboard/sdk';
import SystemMonitorWidget from './components/SystemMonitorWidget.vue';

export const manifest = {
  id: 'system-monitor-plugin',
  name: 'System Monitor',
  version: '1.0.0',
  description: 'Advanced system monitoring',
  permissions: ['storage', 'ui', 'websocket', 'system'] as const,
};

export default definePlugin(({ ui, websocket, logger }) => {
  let wsConnection: WebSocket | null = null;

  return {
    async onActivate() {
      // Connect to WebSocket for real-time data
      wsConnection = websocket.connect('/monitor/ws');
      
      wsConnection.onmessage = (event) => {
        const data = JSON.parse(event.data);
        logger.debug('Received monitoring data:', data);
        // Process and store data
      };

      logger.info('System monitor activated');
    },

    async onDeactivate() {
      // Close WebSocket connection
      if (wsConnection) {
        wsConnection.close();
        wsConnection = null;
      }
      logger.info('System monitor deactivated');
    },

    widgets: [
      {
        id: 'system-monitor-widget',
        name: 'System Monitor',
        component: SystemMonitorWidget,
        icon: 'cpu',
        category: 'monitoring',
        description: 'Real-time system monitoring'
      }
    ]
  };
});
```

---

## Testing Plugins

### Local Testing

1. **Create a test plugin directory:**

```bash
mkdir -p test-plugin/src
cd test-plugin
npm init -y
```

2. **Install SDK:**

```bash
npm install @nas-dashboard/sdk --save-dev
```

3. **Create plugin file:**

```typescript
// src/index.ts
import { definePlugin } from '@nas-dashboard/sdk';

export const manifest = {
  id: 'test-plugin',
  name: 'Test Plugin',
  version: '1.0.0',
  description: 'A test plugin',
  permissions: ['storage'] as const,
};

export default definePlugin(({ logger, storage }) => {
  return {
    async onInstall() {
      logger.info('Test plugin installed');
      await storage.set('test', 'data');
    },

    async onActivate() {
      logger.info('Test plugin activated');
      const data = await storage.get('test');
      logger.info('Retrieved data:', data);
    }
  };
});
```

4. **Build plugin:**

```bash
npm run build
```

5. **Load plugin in dashboard:**

- Navigate to Plugin Store
- Click "Install from File"
- Select built plugin file

### Unit Testing

```typescript
// test-plugin.spec.ts
import { describe, it, expect, beforeEach } from 'vitest';
import { createPluginContext } from '@nas-dashboard/sdk/test-utils';
import plugin from './src/index';

describe('Test Plugin', () => {
  let context: any;

  beforeEach(() => {
    context = createPluginContext();
  });

  it('should install successfully', async () => {
    await plugin.onInstall();
    expect(context.logger.info).toHaveBeenCalledWith('Test plugin installed');
    expect(context.storage.set).toHaveBeenCalledWith('test', 'data');
  });

  it('should activate successfully', async () => {
    await plugin.onActivate();
    expect(context.logger.info).toHaveBeenCalledWith('Test plugin activated');
  });
});
```

---

## Publishing Plugins

### Package Plugin

1. **Build for production:**

```bash
npm run build
```

2. **Create package.json:**

```json
{
  "name": "my-nas-dashboard-plugin",
  "version": "1.0.0",
  "description": "My NAS Dashboard plugin",
  "main": "dist/index.js",
  "files": [
    "dist"
  ],
  "keywords": [
    "nas-dashboard",
    "plugin"
  ]
}
```

3. **Publish to npm:**

```bash
npm publish
```

### Submit to Marketplace

1. **Create plugin metadata:**

```json
{
  "id": "my-plugin",
  "name": "My Plugin",
  "version": "1.0.0",
  "description": "Plugin description",
  "author": "Your Name",
  "license": "MIT",
  "repository": "https://github.com/username/my-plugin",
  "homepage": "https://github.com/username/my-plugin#readme",
  "keywords": ["nas-dashboard", "plugin"],
  "category": "monitoring",
  "screenshot": "https://example.com/screenshot.png",
  "dependencies": {}
}
```

2. **Submit to plugin marketplace:**

- Go to Plugin Store in dashboard
- Click "Submit Plugin"
- Fill in plugin details
- Upload plugin package
- Wait for review

---

## Best Practices

### 1. Performance

```typescript
// Good: Debounce expensive operations
const updateData = debounce(async () => {
  const data = await fetchData();
  renderData(data);
}, 300);

// Bad: No debouncing
const updateData = async () => {
  const data = await fetchData();
  renderData(data);
};
```

### 2. Error Handling

```typescript
// Good: Proper error handling
export default definePlugin(({ logger }) => {
  return {
    async onActivate() {
      try {
        await initializePlugin();
      } catch (error) {
        logger.error('Failed to activate:', error);
        // Handle error gracefully
      }
    }
  };
});

// Bad: No error handling
export default definePlugin(() => {
  return {
    async onActivate() {
      await initializePlugin(); // Might throw
    }
  };
});
```

### 3. Storage Management

```typescript
// Good: Clean up storage on uninstall
export default definePlugin(({ storage }) => {
  return {
    async onUninstall() {
      await storage.clear();
    }
  };
});

// Bad: Leave data behind
export default definePlugin(() => {
  return {
    async onUninstall() {
      // No cleanup
    }
  };
});
```

### 4. Permission Requests

```typescript
// Good: Request only needed permissions
export const manifest = {
  permissions: ['storage', 'ui'] as const
};

// Bad: Request all permissions
export const manifest = {
  permissions: ['storage', 'ui', 'network', 'websocket', 'system'] as const
};
```

### 5. Resource Management

```typescript
// Good: Clean up resources
let timer: NodeJS.Timeout;

export default definePlugin(() => {
  return {
    async onActivate() {
      timer = setInterval(doSomething, 1000);
    },

    async onDeactivate() {
      if (timer) {
        clearInterval(timer);
      }
    }
  };
});

// Bad: Leave resources running
export default definePlugin(() => {
  return {
    async onActivate() {
      setInterval(doSomething, 1000); // Never cleared
    }
  };
});
```

---

## Conclusion

Plugin development for NAS Dashboard provides a powerful way to extend functionality while maintaining security and performance. Follow these guidelines and examples to create robust, maintainable plugins that enhance the NAS Dashboard experience.

### Next Steps

1. **Explore Examples**: Check out the example plugins in the repository
2. **Read API Docs**: Detailed API reference available in API.md
3. **Join Community**: Connect with other plugin developers
4. **Publish Plugin**: Share your plugin with the community

---

**Last Updated**: 2026-06-12  
**Plugin SDK Version**: 0.1.0  
**Status**: Active Development
