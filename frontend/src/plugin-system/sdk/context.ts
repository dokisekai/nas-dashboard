/**
 * Plugin Context Factory - Creates plugin execution context
 */

import type { PluginContext, PluginManifest } from '../types/plugin';
import { createPluginAPI } from './api';
import { createPluginStorage } from './storage';
import { createPluginLogger } from './logger';
import { createPluginUtils } from './utils';

export function createPluginContext(manifest: PluginManifest): PluginContext {
  return {
    pluginId: manifest.id,
    version: manifest.version,
    permissions: manifest.permissions || [],

    // API access
    api: createPluginAPI(manifest),

    // Plugin storage
    storage: createPluginStorage(manifest.id),

    // Utilities
    utils: createPluginUtils(),

    // Logger
    logger: createPluginLogger(manifest.id),

    // Lifecycle hooks (empty by default)
    hooks: {
      onInstall: () => {},
      onUninstall: () => {},
      onUpdate: () => {},
      onEnable: () => {},
      onDisable: () => {},
      onActivate: () => {},
      onDeactivate: () => {}
    }
  };
}
