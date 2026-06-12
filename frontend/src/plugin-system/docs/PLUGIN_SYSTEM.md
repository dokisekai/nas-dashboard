# Plugin System Documentation

## Overview

The NAS Dashboard features a comprehensive plugin system that allows developers to extend functionality with custom features, widgets, data sources, and integrations.

## Architecture

### Core Components

1. **Plugin Loader** - Handles dynamic loading and lifecycle management
2. **Plugin SDK** - Provides APIs and utilities for plugin development
3. **Plugin Manager** - Manages installation, updates, and configuration
4. **Plugin Marketplace** - Handles plugin discovery and distribution

### Plugin Lifecycle

```
Install → Enable → Activate → [Active] → Deactivate → Disable → Uninstall
            ↓                              ↓
         Update (optional)              Update (optional)
```

## Quick Start

### Creating Your First Plugin

1. **Create a plugin file:**

```typescript
// my-plugin.ts
import type { PluginContext } from './plugin-system/types/plugin';

export const manifest = {
  id: 'my-plugin',
  name: 'My Plugin',
  version: '1.0.0',
  description: 'My awesome plugin',
  author: 'Your Name',
  main: 'my-plugin.ts',
  permissions: ['storage'] as const
};

export default function createPlugin(context: PluginContext) {
  const { storage, logger } = context;

  return {
    async onInstall() {
      logger.info('Plugin installed');
      await storage.set('initialized', true);
    },

    async onEnable() {
      logger.info('Plugin enabled');
    },

    async onActivate() {
      logger.info('Plugin active');
      // Your plugin logic here
    }
  };
}
```

2. **Load the plugin:**

```typescript
import { getPluginLoader } from './plugin-system';

const loader = getPluginLoader();
await loader.loadPlugin('/path/to/my-plugin.ts', manifest);
await loader.enablePlugin('my-plugin');
await loader.activatePlugin('my-plugin');
```

## Plugin SDK

### Context API

The plugin context provides access to:

#### `api` - Application Interface

```typescript
// Navigation
api.app.navigate('/dashboard');

// State Management
const state = api.app.getState();
api.app.setState('key', value);

// UI Components
api.ui.registerComponent(myComponent);
api.ui.unregisterComponent('component-id');
api.ui.showNotification('Hello!', 'success');

// Network Requests
const data = await api.network.get('/api/data');
await api.network.post('/api/data', { key: 'value' });

// WebSocket
const ws = api.websocket.connect('ws://localhost:8080');
api.websocket.send(ws, { type: 'message' });
api.websocket.on(ws, 'message', handler);
```

#### `storage` - Persistent Storage

```typescript
// Key-value storage isolated to your plugin
await storage.set('key', { any: 'data' });
const value = await storage.get('key');
await storage.remove('key');
await storage.clear();
const keys = await storage.keys();
```

#### `logger` - Logging

```typescript
logger.debug('Detailed info');
logger.info('General info');
logger.warn('Warning message');
logger.error('Error occurred');
```

#### `utils` - Utilities

```typescript
// Function utilities
const debouncedFn = utils.debounce(myFunction, 300);
const throttledFn = utils.throttle(myFunction, 1000);

// Object utilities
const merged = utils.deepMerge(target, source);
const cloned = utils.clone(original);

// Validation
utils.validate.email('test@example.com');
utils.validate.url('https://example.com');
utils.validate.required(value);

// Event Emitter
const emitter = new utils.EventEmitter();
emitter.on('event', handler);
emitter.emit('event', data);
emitter.off('event', handler);
```

### Lifecycle Hooks

```typescript
export default function createPlugin(context: PluginContext) {
  return {
    // Installation
    async onInstall() {
      // Run once when plugin is installed
      await context.storage.set('installed', true);
    },

    async onUninstall() {
      // Cleanup when plugin is removed
      await context.storage.clear();
    },

    // Updates
    async onUpdate(fromVersion: string, toVersion: string) {
      // Handle version migrations
      if (fromVersion === '1.0.0' && toVersion === '2.0.0') {
        await migrateData();
      }
    },

    // Enable/Disable
    async onEnable() {
      // Plugin is enabled but not active
    },

    async onDisable() {
      // Plugin is disabled
    },

    // Activate/Deactivate
    async onActivate() {
      // Plugin is active and functional
      // Register UI components, start timers, etc.
    },

    async onDeactivate() {
      // Plugin is deactivated
      // Unregister components, stop timers, etc.
    },

    // Optional: Configuration changes
    async onConfigChange(newConfig: any) {
      // Handle configuration updates
    }
  };
}
```

## Plugin Manifest

Required fields:

```typescript
{
  id: string;           // Unique identifier (lowercase, letters, numbers, hyphens)
  name: string;         // Human-readable name
  version: string;      // Semver version (x.y.z)
  description: string;  // Plugin description
  author: string;       // Author name
  main: string;         // Entry point file
}
```

Optional fields:

```typescript
{
  permissions?: PluginPermission[];  // Required permissions
  dependencies?: string[];          // Plugin dependencies
  category?: string;                // Plugin category
  keywords?: string[];              // Search keywords
  icon?: string;                    // Plugin icon URL
  screenshots?: string[];          // Screenshots
  homepage?: string;                // Homepage URL
  repository?: string;              // Repository URL
  license?: string;                 // License
  requires?: {
    api: string;                    // Required API version
    app: string;                    // Required app version
  };
}
```

## Permissions

Plugins must request permissions for sensitive operations:

- **`storage`** - Access persistent storage
- **`network`** - Make network requests
- **`ui`** - Register UI components
- **`settings`** - Modify app settings
- **`notifications`** - Show notifications
- **`websocket`** - Use WebSocket connections
- **`custom`** - Custom permissions

## Plugin Manager

### Installation

```typescript
import { getPluginManager } from './plugin-system';

const manager = getPluginManager();

// Install from URL
await manager.installPlugin(
  'https://example.com/plugin.js',
  manifest
);

// Install from object
await manager.installPlugin(
  pluginCode,
  manifest
);
```

### Updates

```typescript
await manager.updatePlugin(
  'plugin-id',
  newSource,
  newManifest
);
```

### Uninstallation

```typescript
await manager.uninstallPlugin('plugin-id');
```

### Configuration

```typescript
// Get configuration
const config = await manager.getPluginConfig('plugin-id');

// Update configuration
await manager.updatePluginConfig('plugin-id', {
  setting1: 'value1',
  setting2: 'value2'
});

// Reset to defaults
await manager.resetPluginConfig('plugin-id');

// Export/Import
const json = manager.exportPluginConfig('plugin-id');
await manager.importPluginConfig('plugin-id', json);
```

### Permissions

```typescript
// Get current permissions
const permissions = manager.getPluginPermissions('plugin-id');

// Request additional permissions
const result = await manager.requestPermissions(
  'plugin-id',
  ['network', 'websocket']
);
```

## Plugin Marketplace

### Discovery

```typescript
import { getPluginMarketplace } from './plugin-system';

const marketplace = getPluginMarketplace();

// Get all plugins
const plugins = await marketplace.getAllPlugins();

// Search plugins
const results = await marketplace.searchPlugins('analytics', {
  category: 'data',
  minRating: 4,
  sortBy: 'rating'
});

// Get featured plugins
const featured = await marketplace.getFeaturedPlugins();

// Get by category
const widgets = await marketplace.getPluginsByCategory('widgets');
```

### Installation from Marketplace

```typescript
// Get plugin details
const plugin = await marketplace.getPlugin('plugin-id');

// Get download URL
const downloadUrl = await marketplace.getDownloadUrl('plugin-id', '1.0.0');

// Install
const manager = getPluginManager();
const response = await fetch(downloadUrl);
const code = await response.text();
await manager.installPlugin(code, plugin.manifest);
```

### Reviews and Ratings

```typescript
// Get reviews
const reviews = await marketplace.getPluginReviews('plugin-id');

// Submit review
await marketplace.submitReview('plugin-id', {
  userId: 'user-123',
  userName: 'John Doe',
  rating: 5,
  title: 'Great plugin',
  content: 'Works perfectly!'
});

// Rate plugin
await marketplace.ratePlugin('plugin-id', 5);
```

## Error Handling

Plugins should handle errors gracefully:

```typescript
export default function createPlugin(context: PluginContext) {
  const { logger } = context;

  return {
    async onActivate() {
      try {
        // Your plugin logic
        await riskyOperation();
      } catch (error) {
        logger.error('Operation failed:', error);
        // Handle error gracefully
        context.api.ui.showNotification(
          'An error occurred',
          'error'
        );
      }
    }
  };
}
```

## Best Practices

### 1. Use the SDK properly

```typescript
// ✅ Good
const data = await context.api.network.get('/api/data');

// ❌ Bad - bypasses SDK
const data = await fetch('/api/data');
```

### 2. Handle async operations

```typescript
// ✅ Good
async onActivate() {
  await this.initialize();
  await this.loadData();
}

// ❌ Bad - missing await
async onActivate() {
  this.initialize();
  this.loadData();
}
```

### 3. Clean up resources

```typescript
// ✅ Good
async onDeactivate() {
  if (this.interval) {
    clearInterval(this.interval);
  }
  if (this.websocket) {
    this.websocket.close();
  }
}

// ❌ Bad - resource leak
async onDeactivate() {
  // No cleanup
}
```

### 4. Use appropriate logging

```typescript
// ✅ Good
logger.debug('Detailed debug info');
logger.info('General information');
logger.warn('Warning message');
logger.error('Error occurred');

// ❌ Bad - using console
console.log('This bypasses the plugin logger');
```

### 5. Validate input

```typescript
// ✅ Good
async setData(key: string, value: any) {
  if (!context.utils.validate.required(key)) {
    throw new Error('Key is required');
  }
  await context.storage.set(key, value);
}

// ❌ Bad - no validation
async setData(key: string, value: any) {
  await context.storage.set(key, value);
}
```

## Troubleshooting

### Plugin won't load

1. Check manifest format is correct
2. Verify all dependencies are installed
3. Check browser console for errors
4. Ensure permissions are declared

### Lifecycle hooks not firing

1. Verify hook is exported in plugin
2. Check for errors in previous hooks
3. Ensure plugin is in correct state
4. Check browser console for errors

### Storage issues

1. Ensure you're using async/await correctly
2. Check storage quota isn't exceeded
3. Handle errors with try-catch
4. Use correct storage methods

### Performance issues

1. Debounce/throttle expensive operations
2. Cache data appropriately
3. Clean up event listeners
4. Use efficient data structures

## Security Considerations

1. **Never expose sensitive data** - Don't log passwords or tokens
2. **Validate all input** - Check user input and API responses
3. **Use HTTPS** - Always use secure connections
4. **Sanitize data** - Clean data before displaying
5. **Limit permissions** - Only request what you need
6. **Handle errors** - Don't expose stack traces

## Testing Your Plugin

```typescript
// test-plugin.ts
import { getPluginLoader } from './plugin-system';

async function testPlugin() {
  const loader = getPluginLoader();

  // Load plugin
  const result = await loader.loadPlugin(
    '/path/to/plugin.ts',
    manifest
  );

  if (!result.success) {
    console.error('Failed to load:', result.error);
    return;
  }

  // Test lifecycle
  await loader.enablePlugin(manifest.id);
  await loader.activatePlugin(manifest.id);

  // Test functionality
  const plugin = loader.getPluginState(manifest.id);
  // Test plugin methods...

  // Cleanup
  await loader.deactivatePlugin(manifest.id);
  await loader.unloadPlugin(manifest.id);
}

testPlugin();
```

## Distribution

### Packaging

Create a distributable package:

```
my-plugin/
├── package.json
├── manifest.json
├── dist/
│   └── plugin.js
└── README.md
```

### Publishing

1. Submit to marketplace
2. Provide documentation
3. Include examples
4. Support users

## Support

For issues, questions, or contributions, please refer to the main project repository.
