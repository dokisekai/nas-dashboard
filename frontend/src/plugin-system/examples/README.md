# Plugin Examples

This directory contains example plugins to help you get started with plugin development.

## Templates

### Basic Plugin (`templates/basic-plugin.ts`)
A minimal plugin template showing:
- Plugin manifest structure
- Lifecycle hooks
- Basic storage operations
- Logging

### Widget Plugin (`templates/widget-plugin.ts`)
Template for creating dashboard widgets:
- UI component registration
- Widget rendering
- Dynamic content
- Styling

### Data Source Plugin (`templates/data-source-plugin.ts`)
Template for data provider plugins:
- Network API integration
- WebSocket real-time data
- Data caching
- Debouncing utilities

## Complete Examples

### Notification Plugin (`plugins/notification-plugin.ts`)
A fully-featured notification system:
- Custom sounds and filters
- Desktop notifications
- Notification history
- Settings management
- UI panel

### Analytics Plugin (`plugins/analytics-plugin.ts`)
Advanced analytics and reporting:
- Event tracking
- Batch processing
- Dashboard visualization
- Metrics queries
- Report generation

## How to Use

1. **Choose a template** that matches your plugin type
2. **Copy the template** to your project
3. **Update the manifest** with your plugin information
4. **Implement your functionality** using the SDK
5. **Test locally** using the plugin manager
6. **Package and distribute** through the marketplace

## Plugin Development Best Practices

1. **Always use the SDK** - Don't bypass the plugin API
2. **Handle errors gracefully** - Use try-catch blocks
3. **Log important events** - Use the provided logger
4. **Clean up resources** - Implement onDeactivate properly
5. **Version your plugin** - Follow semantic versioning
6. **Test lifecycle hooks** - Ensure install/update/uninstall work
7. **Document your API** - Help other developers integrate
8. **Use permissions appropriately** - Only request what you need

## SDK Reference

The plugin SDK provides:

- `api` - Application API (navigation, UI, network, WebSocket)
- `storage` - Isolated key-value storage
- `logger` - Plugin-specific logging
- `utils` - Common utilities (debounce, throttle, validation)
- `hooks` - Lifecycle hooks

## Testing Your Plugin

```typescript
// Load your plugin for testing
import { getPluginLoader } from './plugin-system/core/PluginLoader';

const loader = getPluginLoader();
const result = await loader.loadPlugin(
  '/path/to/your/plugin.ts',
  manifest
);

if (result.success) {
  await loader.enablePlugin(manifest.id);
  await loader.activatePlugin(manifest.id);
}
```

## Troubleshooting

### Plugin won't load
- Check manifest syntax
- Verify dependencies are installed
- Check browser console for errors

### Lifecycle hooks not firing
- Ensure hook is defined in plugin export
- Check for errors in previous hooks
- Verify plugin state in registry

### Storage issues
- Use correct storage methods (get/set/remove)
- Handle errors with try-catch
- Check storage quota

## Additional Resources

- [Plugin System Documentation](../docs/PLUGIN_SYSTEM.md)
- [API Reference](../docs/API_REFERENCE.md)
- [Marketplace Guidelines](../docs/MARKETPLACE.md)
