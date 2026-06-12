/**
 * Plugin System Type Definitions
 */

export interface PluginManifest {
  id: string;
  name: string;
  version: string;
  description: string;
  author: string;
  license?: string;
  homepage?: string;
  repository?: string;

  // Plugin requirements
  dependencies?: string[];
  peerDependencies?: Record<string, string>;

  // Plugin capabilities
  permissions?: PluginPermission[];

  // Plugin entry points
  main: string;
  styles?: string;

  // Plugin configuration
  config?: Record<string, any>;

  // Plugin metadata
  keywords?: string[];
  category?: string;
  icon?: string;
  screenshots?: string[];

  // Compatibility
  requires?: {
    api: string;
    app: string;
  };
}

export type PluginPermission =
  | 'storage'
  | 'network'
  | 'ui'
  | 'settings'
  | 'notifications'
  | 'websocket'
  | 'custom';

export interface PluginContext {
  pluginId: string;
  version: string;
  permissions: PluginPermission[];

  // API access
  api: PluginAPI;

  // Plugin storage
  storage: PluginStorage;

  // Lifecycle hooks
  hooks: PluginHooks;

  // Utilities
  utils: PluginUtils;

  // Logger
  logger: PluginLogger;
}

export interface PluginAPI {
  // Application API
  app: {
    navigate: (path: string) => void;
    getState: () => any;
    setState: (key: string, value: any) => void;
  };

  // UI API
  ui: {
    registerComponent: (component: any) => void;
    unregisterComponent: (componentId: string) => void;
    showNotification: (message: string, type?: 'success' | 'error' | 'warning' | 'info') => void;
  };

  // Network API
  network: {
    request: (config: any) => Promise<any>;
    get: (url: string, config?: any) => Promise<any>;
    post: (url: string, data?: any, config?: any) => Promise<any>;
    put: (url: string, data?: any, config?: any) => Promise<any>;
    delete: (url: string, config?: any) => Promise<any>;
  };

  // WebSocket API
  websocket: {
    connect: (url: string) => any;
    send: (data: any) => void;
    on: (event: string, handler: Function) => void;
    off: (event: string, handler: Function) => void;
  };
}

export interface PluginStorage {
  get: (key: string) => Promise<any>;
  set: (key: string, value: any) => Promise<void>;
  remove: (key: string) => Promise<void>;
  clear: () => Promise<void>;
  keys: () => Promise<string[]>;
}

export interface PluginHooks {
  onInstall: () => void | Promise<void>;
  onUninstall: () => void | Promise<void>;
  onUpdate: (fromVersion: string, toVersion: string) => void | Promise<void>;
  onEnable: () => void | Promise<void>;
  onDisable: () => void | Promise<void>;
  onActivate: () => void | Promise<void>;
  onDeactivate: () => void | Promise<void>;
}

export interface PluginUtils {
  // Utilities
  debounce: (func: Function, wait: number) => Function;
  throttle: (func: Function, limit: number) => Function;
  deepMerge: (target: any, source: any) => any;
  clone: <T>(obj: T) => T;

  // Event emitter
  EventEmitter: class;

  // Validation
  validate: {
    email: (email: string) => boolean;
    url: (url: string) => boolean;
    required: (value: any) => boolean;
  };
}

export interface PluginLogger {
  debug: (...args: any[]) => void;
  info: (...args: any[]) => void;
  warn: (...args: any[]) => void;
  error: (...args: any[]) => void;
}

export interface PluginLoadResult {
  success: boolean;
  pluginId: string;
  error?: Error;
  warnings?: string[];
}

export interface PluginState {
  id: string;
  manifest: PluginManifest;
  loaded: boolean;
  enabled: boolean;
  active: boolean;
  error: Error | null;
  version: string;
  instance?: any;
}

export interface PluginError extends Error {
  pluginId: string;
  code: string;
  details?: any;
}

export interface PluginRegistry {
  [pluginId: string]: PluginState;
}

export interface PluginMarketplaceInfo {
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
