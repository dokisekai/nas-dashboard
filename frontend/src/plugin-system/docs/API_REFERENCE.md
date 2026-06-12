# Plugin API Reference

Complete API reference for the NAS Dashboard Plugin System.

## Table of Contents

- [Plugin Loader](#plugin-loader)
- [Plugin Manager](#plugin-manager)
- [Plugin Marketplace](#plugin-marketplace)
- [Plugin SDK](#plugin-sdk)
- [Types](#types)

## Plugin Loader

### `getPluginLoader()`

Get the singleton PluginLoader instance.

```typescript
const loader = getPluginLoader();
```

### `loadPlugin(source, manifest)`

Load a plugin from a URL or code object.

**Parameters:**
- `source: string | object` - Plugin source URL or code object
- `manifest: PluginManifest` - Plugin manifest object

**Returns:** `Promise<PluginLoadResult>`

```typescript
const result = await loader.loadPlugin(
  'https://example.com/plugin.js',
  manifest
);
```

### `enablePlugin(pluginId)`

Enable a loaded plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<void>`

```typescript
await loader.enablePlugin('my-plugin');
```

### `disablePlugin(pluginId)`

Disable an enabled plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<void>`

```typescript
await loader.disablePlugin('my-plugin');
```

### `activatePlugin(pluginId)`

Activate an enabled plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<void>`

```typescript
await loader.activatePlugin('my-plugin');
```

### `deactivatePlugin(pluginId)`

Deactivate an active plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<void>`

```typescript
await loader.deactivatePlugin('my-plugin');
```

### `unloadPlugin(pluginId)`

Unload and remove a plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<void>`

```typescript
await loader.unloadPlugin('my-plugin');
```

### `getPluginState(pluginId)`

Get plugin state.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `PluginState | undefined`

```typescript
const state = loader.getPluginState('my-plugin');
```

### `getAllPlugins()`

Get all loaded plugins.

**Returns:** `PluginState[]`

```typescript
const plugins = loader.getAllPlugins();
```

### `getEnabledPlugins()`

Get all enabled plugins.

**Returns:** `PluginState[]`

```typescript
const enabled = loader.getEnabledPlugins();
```

### `getActivePlugins()`

Get all active plugins.

**Returns:** `PluginState[]`

```typescript
const active = loader.getActivePlugins();
```

## Plugin Manager

### `getPluginManager()`

Get the singleton PluginManager instance.

```typescript
const manager = getPluginManager();
```

### `installPlugin(pluginSource, manifest)`

Install a new plugin.

**Parameters:**
- `pluginSource: string | object` - Plugin source
- `manifest: PluginManifest` - Plugin manifest

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await manager.installPlugin(
  'https://example.com/plugin.js',
  manifest
);
```

### `uninstallPlugin(pluginId)`

Uninstall a plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await manager.uninstallPlugin('my-plugin');
```

### `updatePlugin(pluginId, newSource, newManifest)`

Update a plugin to a new version.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `newSource: string | object` - New plugin source
- `newManifest: PluginManifest` - New manifest

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await manager.updatePlugin(
  'my-plugin',
  newSource,
  newManifest
);
```

### `getPluginConfig(pluginId)`

Get plugin configuration.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<any>`

```typescript
const config = await manager.getPluginConfig('my-plugin');
```

### `updatePluginConfig(pluginId, config)`

Update plugin configuration.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `config: any` - Configuration object

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await manager.updatePluginConfig('my-plugin', {
  setting: 'value'
});
```

### `resetPluginConfig(pluginId)`

Reset plugin configuration to defaults.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await manager.resetPluginConfig('my-plugin');
```

### `exportPluginConfig(pluginId)`

Export plugin configuration as JSON string.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `string`

```typescript
const json = manager.exportPluginConfig('my-plugin');
```

### `importPluginConfig(pluginId, configJson)`

Import plugin configuration from JSON string.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `configJson: string` - JSON configuration string

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await manager.importPluginConfig('my-plugin', json);
```

### `getPluginPermissions(pluginId)`

Get plugin permissions.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `string[]`

```typescript
const permissions = manager.getPluginPermissions('my-plugin');
```

### `requestPermissions(pluginId, permissions)`

Request additional permissions.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `permissions: string[]` - Requested permissions

**Returns:** `Promise<{ granted: boolean; error?: string }>`

```typescript
const result = await manager.requestPermissions('my-plugin', ['network']);
```

### `getPluginStats(pluginId)`

Get plugin statistics.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `PluginStats | null`

```typescript
const stats = manager.getPluginStats('my-plugin');
```

### `validateCompatibility(manifest)`

Validate plugin compatibility.

**Parameters:**
- `manifest: PluginManifest` - Plugin manifest

**Returns:** `{ compatible: boolean; issues: string[] }`

```typescript
const validation = manager.validateCompatibility(manifest);
```

## Plugin Marketplace

### `getPluginMarketplace()`

Get the singleton PluginMarketplace instance.

```typescript
const marketplace = getPluginMarketplace();
```

### `configurePluginMarketplace(config)`

Configure marketplace with custom settings.

**Parameters:**
- `config: MarketplaceConfig` - Marketplace configuration

**Returns:** `PluginMarketplace`

```typescript
const marketplace = configurePluginMarketplace({
  apiEndpoint: 'https://api.example.com/marketplace',
  cacheTimeout: 60000
});
```

### `getAllPlugins()`

Get all available plugins.

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const plugins = await marketplace.getAllPlugins();
```

### `getPlugin(pluginId)`

Get plugin details.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<PluginMarketplaceInfo | null>`

```typescript
const plugin = await marketplace.getPlugin('my-plugin');
```

### `searchPlugins(query, filters)`

Search for plugins.

**Parameters:**
- `query: string` - Search query
- `filters?: SearchFilters` - Optional filters

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const results = await marketplace.searchPlugins('analytics', {
  category: 'data',
  minRating: 4,
  sortBy: 'rating'
});
```

### `getPluginsByCategory(category)`

Get plugins by category.

**Parameters:**
- `category: string` - Category name

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const plugins = await marketplace.getPluginsByCategory('widgets');
```

### `getCategories()`

Get all plugin categories.

**Returns:** `Promise<string[]>`

```typescript
const categories = await marketplace.getCategories();
```

### `getFeaturedPlugins()`

Get featured plugins.

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const featured = await marketplace.getFeaturedPlugins();
```

### `getPopularPlugins(limit)`

Get popular plugins.

**Parameters:**
- `limit: number` - Number of plugins to return

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const popular = await marketplace.getPopularPlugins(10);
```

### `getRecentlyUpdated(limit)`

Get recently updated plugins.

**Parameters:**
- `limit: number` - Number of plugins to return

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const recent = await marketplace.getRecentlyUpdated(10);
```

### `getPluginReviews(pluginId)`

Get plugin reviews.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<PluginReview[]>`

```typescript
const reviews = await marketplace.getPluginReviews('my-plugin');
```

### `submitReview(pluginId, review)`

Submit a plugin review.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `review: Omit<PluginReview, 'id' | 'date'>` - Review data

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await marketplace.submitReview('my-plugin', {
  userId: 'user-123',
  userName: 'John Doe',
  rating: 5,
  title: 'Great plugin',
  content: 'Works perfectly!'
});
```

### `ratePlugin(pluginId, rating)`

Rate a plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `rating: number` - Rating (1-5)

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await marketplace.ratePlugin('my-plugin', 5);
```

### `getDownloadUrl(pluginId, version)`

Get plugin download URL.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `version?: string` - Optional version

**Returns:** `Promise<string>`

```typescript
const url = await marketplace.getDownloadUrl('my-plugin', '1.0.0');
```

### `getVersions(pluginId)`

Get plugin versions.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `Promise<PluginVersion[]>`

```typescript
const versions = await marketplace.getVersions('my-plugin');
```

### `reportPlugin(pluginId, reason, details)`

Report a plugin.

**Parameters:**
- `pluginId: string` - Plugin identifier
- `reason: string` - Report reason
- `details?: string` - Additional details

**Returns:** `Promise<{ success: boolean; error?: string }>`

```typescript
const result = await marketplace.reportPlugin('my-plugin', 'Broken', 'Does not work');
```

### `getDeveloperInfo(developerId)`

Get developer information.

**Parameters:**
- `developerId: string` - Developer identifier

**Returns:** `Promise<DeveloperInfo | null>`

```typescript
const info = await marketplace.getDeveloperInfo('dev-123');
```

### `getPluginsByDeveloper(developerId)`

Get plugins by developer.

**Parameters:**
- `developerId: string` - Developer identifier

**Returns:** `Promise<PluginMarketplaceInfo[]>`

```typescript
const plugins = await marketplace.getPluginsByDeveloper('dev-123');
```

## Plugin SDK

### `createPluginContext(manifest)`

Create a plugin context.

**Parameters:**
- `manifest: PluginManifest` - Plugin manifest

**Returns:** `PluginContext`

```typescript
const context = createPluginContext(manifest);
```

### `createPluginAPI(manifest)`

Create plugin API object.

**Parameters:**
- `manifest: PluginManifest` - Plugin manifest

**Returns:** `PluginAPI`

```typescript
const api = createPluginAPI(manifest);
```

### `createPluginStorage(pluginId)`

Create plugin storage object.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `PluginStorage`

```typescript
const storage = createPluginStorage('my-plugin');
```

### `createPluginLogger(pluginId)`

Create plugin logger object.

**Parameters:**
- `pluginId: string` - Plugin identifier

**Returns:** `PluginLogger`

```typescript
const logger = createPluginLogger('my-plugin');
```

### `createPluginUtils()`

Create plugin utilities object.

**Returns:** `PluginUtils`

```typescript
const utils = createPluginUtils();
```

## Types

### `PluginManifest`

```typescript
interface PluginManifest {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  license?: string;
  homepage?: string;
  repository?: string;
  dependencies?: string[];
  peerDependencies?: Record<string, string>;
  permissions?: PluginPermission[];
  main: string;
  styles?: string;
  keywords?: string[];
  category?: string;
  icon?: string;
  screenshots?: string[];
  requires?: {
    api: string;
    app: string;
  };
}
```

### `PluginContext`

```typescript
interface PluginContext {
  pluginId: string;
  version: string;
  permissions: PluginPermission[];
  api: PluginAPI;
  storage: PluginStorage;
  hooks: PluginHooks;
  utils: PluginUtils;
  logger: PluginLogger;
}
```

### `PluginAPI`

```typescript
interface PluginAPI {
  app: {
    navigate: (path: string) => void;
    getState: () => any;
    setState: (key: string, value: any) => void;
  };
  ui: {
    registerComponent: (component: any) => void;
    unregisterComponent: (componentId: string) => void;
    showNotification: (message: string, type?: string) => void;
  };
  network: {
    request: (config: any) => Promise<any>;
    get: (url: string, config?: any) => Promise<any>;
    post: (url: string, data?: any, config?: any) => Promise<any>;
    put: (url: string, data?: any, config?: any) => Promise<any>;
    delete: (url: string, config?: any) => Promise<any>;
  };
  websocket: {
    connect: (url: string) => any;
    send: (ws: any, data: any) => void;
    on: (ws: any, event: string, handler: Function) => void;
    off: (ws: any, event: string, handler: Function) => void;
  };
}
```

### `PluginStorage`

```typescript
interface PluginStorage {
  get: (key: string) => Promise<any>;
  set: (key: string, value: any) => Promise<void>;
  remove: (key: string) => Promise<void>;
  clear: () => Promise<void>;
  keys: () => Promise<string[]>;
}
```

### `PluginLogger`

```typescript
interface PluginLogger {
  debug: (...args: any[]) => void;
  info: (...args: any[]) => void;
  warn: (...args: any[]) => void;
  error: (...args: any[]) => void;
}
```

### `PluginUtils`

```typescript
interface PluginUtils {
  debounce: (func: Function, wait: number) => Function;
  throttle: (func: Function, limit: number) => Function;
  deepMerge: (target: any, source: any) => any;
  clone: <T>(obj: T) => T;
  EventEmitter: class;
  validate: {
    email: (email: string) => boolean;
    url: (url: string) => boolean;
    required: (value: any) => boolean;
  };
}
```

### `PluginHooks`

```typescript
interface PluginHooks {
  onInstall: () => void | Promise<void>;
  onUninstall: () => void | Promise<void>;
  onUpdate: (fromVersion: string, toVersion: string) => void | Promise<void>;
  onEnable: () => void | Promise<void>;
  onDisable: () => void | Promise<void>;
  onActivate: () => void | Promise<void>;
  onDeactivate: () => void | Promise<void>;
}
```

### `PluginState`

```typescript
interface PluginState {
  id: string;
  manifest: PluginManifest;
  loaded: boolean;
  enabled: boolean;
  active: boolean;
  error: Error | null;
  version: string;
  instance?: any;
}
```

### `PluginLoadResult`

```typescript
interface PluginLoadResult {
  success: boolean;
  pluginId: string;
  error?: Error;
  warnings?: string[];
}
```

### `PluginMarketplaceInfo`

```typescript
interface PluginMarketplaceInfo {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  downloads: number;
  rating: number;
  reviews: number;
  lastUpdated: string;
  homepage?: string;
  repository?: string;
  icon?: string;
  screenshots?: string[];
  tags: string[];
  category: string;
  price?: number;
}
```

### `PluginPermission`

```typescript
type PluginPermission =
  | 'storage'
  | 'network'
  | 'ui'
  | 'settings'
  | 'notifications'
  | 'websocket'
  | 'custom';
```

## Error Handling

All async methods can throw errors. Always handle errors appropriately:

```typescript
try {
  await loader.loadPlugin(source, manifest);
} catch (error) {
  console.error('Failed to load plugin:', error);
  // Handle error
}
```

For methods that return result objects, check the success flag:

```typescript
const result = await manager.installPlugin(source, manifest);
if (!result.success) {
  console.error('Installation failed:', result.error);
  // Handle error
}
```

## Events and Callbacks

Plugin lifecycle hooks are called automatically during state transitions:

1. **onInstall** - Called once when plugin is first installed
2. **onUninstall** - Called when plugin is removed
3. **onUpdate** - Called when plugin is updated
4. **onEnable** - Called when plugin is enabled
5. **onDisable** - Called when plugin is disabled
6. **onActivate** - Called when plugin is activated
7. **onDeactivate** - Called when plugin is deactivated

All hooks are optional and can be sync or async.

## Best Practices

1. **Always handle errors** - Use try-catch blocks
2. **Check return values** - Verify success flags
3. **Clean up resources** - Implement proper cleanup in onDeactivate
4. **Use appropriate permissions** - Only request what you need
5. **Log important events** - Use the provided logger
6. **Validate input** - Check parameters before use
7. **Handle edge cases** - Consider null/undefined values
8. **Test thoroughly** - Test all lifecycle hooks
