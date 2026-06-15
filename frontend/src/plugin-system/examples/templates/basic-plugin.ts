/**
 * Basic Plugin Template
 * Copy this template to create your own plugin
 */

import type { PluginContext } from '../../types/plugin';

// Plugin manifest (required)
export const manifest = {
  id: 'my-basic-plugin',
  name: 'My Basic Plugin',
  version: '1.0.0',
  description: 'A basic plugin template',
  author: 'Your Name',
  main: 'plugin.ts',

  // Permissions your plugin needs
  permissions: ['storage', 'ui'] as any,

  // Dependencies on other plugins
  dependencies: [] as string[],

  // Optional metadata
  category: 'utilities',
  keywords: ['basic', 'template', 'example'],
  icon: '',

  // Version requirements
  requires: {
    api: '1.0.0',
    app: '1.0.0'
  }
};

// Plugin factory function (required)
export default function createPlugin(context: PluginContext) {
  const { api, storage, logger, utils, pluginId } = context;

  logger.info('Plugin initialized:', pluginId);

  return {
    // Lifecycle hooks

    async onInstall() {
      logger.info('Plugin installed');
      // Run installation logic
      await storage.set('initialized', true);
    },

    async onEnable() {
      logger.info('Plugin enabled');
      // Run enable logic
    },

    async onActivate() {
      logger.info('Plugin activated');
      // Register UI components
      // api.ui.registerComponent(MyComponent);
    },

    async onDeactivate() {
      logger.info('Plugin deactivated');
      // Unregister UI components
      // api.ui.unregisterComponent('my-component');
    },

    async onDisable() {
      logger.info('Plugin disabled');
      // Run disable logic
    },

    async onUninstall() {
      logger.info('Plugin uninstalled');
      // Cleanup
      await storage.clear();
    },

    async onUpdate(fromVersion: string, toVersion: string) {
      logger.info(`Updating from ${fromVersion} to ${toVersion}`);
      // Run migration logic
    },

    // Optional: Configuration change handler
    async onConfigChange(newConfig: any) {
      logger.info('Configuration changed:', newConfig);
    },

    // Your plugin's public API
    doSomething() {
      logger.info('Doing something');
    },

    async getData() {
      return await storage.get('myData');
    },

    async setData(data: any) {
      await storage.set('myData', data);
    }
  };
}
