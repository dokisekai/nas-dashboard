# NAS Dashboard Plugin System

A comprehensive, production-ready plugin system for the NAS Dashboard that enables dynamic feature extension through modular, sandboxed plugins.

## Features

### Core Capabilities

- **Dynamic Plugin Loading** - Load plugins at runtime without app restart
- **Lifecycle Management** - Complete plugin lifecycle control (install, enable, activate, deactivate, disable, uninstall)
- **Dependency Resolution** - Automatic dependency management and resolution
- **Error Isolation** - Plugin errors don't crash the main application
- **Permission System** - Granular permission control for security
- **Isolated Storage** - Each plugin has its own storage namespace
- **Hot Reload** - Load, update, and unload plugins without restart

### Developer Experience

- **Type-Safe SDK** - Full TypeScript support with comprehensive types
- **Rich API** - Access to navigation, UI, network, WebSocket, and storage
- **Utilities** - Built-in helpers for common operations
- **Logging** - Plugin-specific logging with debug support
- **Testing Tools** - Comprehensive testing utilities
- **Templates** - Ready-to-use plugin templates
- **Examples** - Complete working examples

### Marketplace Integration

- **Plugin Discovery** - Search and browse available plugins
- **Ratings & Reviews** - Community feedback system
- **Version Management** - Multiple version support
- **Developer Profiles** - Track plugin creators
- **Category Browsing** - Organized by functionality

## Quick Start

### Installation

The plugin system is included in the NAS Dashboard frontend.

```typescript
import {
  getPluginLoader,
  getPluginManager,
  getPluginMarketplace
} from './plugin-system';
```

### Creating Your First Plugin

```typescript
// my-plugin.ts
import type { PluginContext } from './plugin-system/types/plugin';

export const manifest = {
  id: 'my-first-plugin',
  name: 'My First Plugin',
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
      logger.info('Plugin installed!');
      await storage.set('initialized', true);
    },

    async onActivate() {
      logger.info('Plugin activated!');
      // Your plugin logic here
    }
  };
}
```

### Loading Your Plugin

```typescript
import { getPluginLoader } from './plugin-system';

const loader = getPluginLoader();

// Load from URL
await loader.loadPlugin('https://example.com/plugin.js', manifest);

// Or load from object
await loader.loadPlugin(pluginCode, manifest);

// Enable and activate
await loader.enablePlugin('my-first-plugin');
await loader.activatePlugin('my-first-plugin');
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    NAS Dashboard                          │
│                                                             │
│  ┌──────────────────────────────────────────────────┐   │
│  │              Plugin System Core                   │   │
│  │                                                    │   │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │   │
│  │  │   Loader   │  │  Manager   │  │Marketplace │  │   │
│  │  └────────────┘  └────────────┘  └────────────┘  │   │
│  └──────────────────────────────────────────────────┘   │
│                           │                               │
│  ┌──────────────────────────────────────────────────┐   │
│  │                 Plugin SDK                        │   │
│  │                                                    │   │
│  │  API • Storage • Logger • Utils                  │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                           │
┌─────────────────────────────────────────────────────────┐
│                  Plugin Sandbox                         │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐           │
│  │ Plugin A │  │ Plugin B │  │ Plugin C │  ...      │
│  └──────────┘  └──────────┘  └──────────┘           │
└─────────────────────────────────────────────────────────┘
```

## Plugin Lifecycle

```
┌──────────────────────────────────────────────────────────────────┐
│                            Plugin Lifecycle                        │
└──────────────────────────────────────────────────────────────────┘

  Install → Enable → Activate → [Active] → Deactivate → Disable → Uninstall
              ↓                              ↓
           Update                        Update

  States:
  - Installed: Plugin files present but not loaded
  - Loaded: Plugin code loaded into memory
  - Enabled: Plugin ready to be activated
  - Active: Plugin functional and running
  - Disabled: Plugin disabled but installed
  - Uninstalled: Plugin removed from system
```

## Documentation

- [Plugin System Guide](./docs/PLUGIN_SYSTEM.md) - Comprehensive guide
- [API Reference](./docs/API_REFERENCE.md) - Complete API documentation
- [Examples](./examples/) - Plugin examples and templates

## Plugin Examples

### Basic Plugin
Simple template showing core concepts:
```bash
examples/templates/basic-plugin.ts
```

### Widget Plugin
Dashboard widget template:
```bash
examples/templates/widget-plugin.ts
```

### Data Source Plugin
Data provider template:
```bash
examples/templates/data-source-plugin.ts
```

### Complete Examples

**Notification Plugin**
- Custom sounds and filters
- Desktop notifications
- History management
```bash
examples/plugins/notification-plugin.ts
```

**Analytics Plugin**
- Event tracking
- Dashboard visualization
- Report generation
```bash
examples/plugins/analytics-plugin.ts
```

## Testing

Comprehensive testing utilities included:

```typescript
import { runPluginTests, quickTest } from './plugin-system/examples/testing';

// Run full test suite
await runPluginTests(pluginCode, manifest);

// Quick test
const success = await quickTest(pluginCode, manifest);
```

## API Overview

### Plugin Loader

```typescript
// Load plugin
await loader.loadPlugin(source, manifest);

// Lifecycle control
await loader.enablePlugin(pluginId);
await loader.activatePlugin(pluginId);
await loader.deactivatePlugin(pluginId);
await loader.disablePlugin(pluginId);
await loader.unloadPlugin(pluginId);

// Query state
const state = loader.getPluginState(pluginId);
const allPlugins = loader.getAllPlugins();
```

### Plugin Manager

```typescript
// Installation
await manager.installPlugin(source, manifest);
await manager.uninstallPlugin(pluginId);

// Updates
await manager.updatePlugin(pluginId, newSource, newManifest);

// Configuration
const config = await manager.getPluginConfig(pluginId);
await manager.updatePluginConfig(pluginId, config);
await manager.resetPluginConfig(pluginId);

// Permissions
const perms = manager.getPluginPermissions(pluginId);
await manager.requestPermissions(pluginId, ['network']);
```

### Plugin Marketplace

```typescript
// Discovery
const plugins = await marketplace.getAllPlugins();
const results = await marketplace.searchPlugins('analytics');

// Categories
const categories = await marketplace.getCategories();
const widgets = await marketplace.getPluginsByCategory('widgets');

// Reviews
const reviews = await marketplace.getPluginReviews(pluginId);
await marketplace.submitReview(pluginId, review);
```

## Plugin SDK

### Context API

```typescript
// Navigation
api.app.navigate('/dashboard');

// UI Components
api.ui.registerComponent(component);
api.ui.showNotification('Message', 'success');

// Network
const data = await api.network.get('/api/data');

// WebSocket
const ws = api.websocket.connect('ws://localhost:8080');
```

### Storage

```typescript
// Key-value storage
await storage.set('key', data);
const value = await storage.get('key');
await storage.remove('key');
await storage.clear();
```

### Utilities

```typescript
// Function utilities
const debounced = utils.debounce(fn, 300);
const throttled = utils.throttle(fn, 1000);

// Object utilities
const merged = utils.deepMerge(target, source);
const cloned = utils.clone(original);

// Validation
utils.validate.email('test@example.com');
utils.validate.url('https://example.com');

// Event Emitter
const emitter = new utils.EventEmitter();
emitter.on('event', handler);
emitter.emit('event', data);
```

### Logger

```typescript
logger.debug('Detailed info');
logger.info('General info');
logger.warn('Warning');
logger.error('Error');
```

## Security

### Permission System

Plugins must declare required permissions:

```typescript
permissions: ['storage', 'network', 'ui', 'settings', 'notifications', 'websocket']
```

### Sandboxing

- Isolated execution context
- No direct DOM access
- No direct localStorage access
- Controlled network access
- Limited WebSocket access

### Best Practices

1. **Never expose sensitive data**
2. **Validate all input**
3. **Use HTTPS only**
4. **Sanitize user input**
5. **Handle errors gracefully**
6. **Limit permissions**
7. **Clean up resources**

## Development

### Project Structure

```
plugin-system/
├── core/              # Core plugin loading
│   └── PluginLoader.ts
├── sdk/               # Plugin SDK
│   ├── context.ts
│   ├── api.ts
│   ├── storage.ts
│   ├── logger.ts
│   └── utils.ts
├── manager/           # Plugin management
│   └── PluginManager.ts
├── marketplace/       # Marketplace API
│   └── PluginMarketplace.ts
├── types/             # TypeScript types
│   └── plugin.ts
├── examples/          # Examples and templates
│   ├── templates/
│   ├── plugins/
│   └── testing/
├── docs/             # Documentation
│   ├── PLUGIN_SYSTEM.md
│   └── API_REFERENCE.md
└── index.ts          # Main entry point
```

### Building

```bash
# Install dependencies
npm install

# Build for development
npm run dev

# Build for production
npm run build

# Run tests
npm test
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Write tests for your changes
4. Ensure all tests pass
5. Submit a pull request

## Support

For issues, questions, or contributions:

- GitHub Issues: [Project Issues]
- Documentation: [Plugin System Guide](./docs/PLUGIN_SYSTEM.md)
- Examples: [Plugin Examples](./examples/)

## License

MIT License - see LICENSE file for details

## Credits

Developed as part of the NAS Dashboard project.

---

**Ready to build amazing plugins?** Start with our [Basic Plugin Template](./examples/templates/basic-plugin.ts) or explore the [Complete Examples](./examples/plugins/)!
