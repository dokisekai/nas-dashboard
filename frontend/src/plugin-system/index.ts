/**
 * Plugin System - Main Entry Point
 * Complete plugin architecture for the NAS Dashboard
 */

// Core
export { getPluginLoader, PluginLoader } from './core/PluginLoader';

// SDK
export {
  createPluginContext,
  createPluginAPI,
  createPluginStorage,
  createPluginLogger,
  createPluginUtils
} from './sdk';

// Manager
export { getPluginManager, PluginManager } from './manager';

// Marketplace
export {
  getPluginMarketplace,
  configurePluginMarketplace,
  PluginMarketplace
} from './marketplace';

// Types
export type {
  PluginManifest,
  PluginContext,
  PluginAPI,
  PluginStorage,
  PluginLogger,
  PluginUtils,
  PluginHooks,
  PluginLoadResult,
  PluginState,
  PluginError,
  PluginPermission,
  PluginRegistry,
  PluginMarketplaceInfo
} from './types/plugin';
