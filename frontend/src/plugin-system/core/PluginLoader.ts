/**
 * Plugin Loader - Handles dynamic plugin loading and lifecycle management
 */

import type {
  PluginManifest,
  PluginContext,
  PluginLoadResult,
  PluginState,
  PluginError,
  PluginHooks
} from '../types/plugin';

export class PluginLoader {
  private registry: Map<string, PluginState> = new Map();
  private dependencyGraph: Map<string, Set<string>> = new Map();
  private loadQueue: string[] = [];
  private isProcessingQueue = false;

  constructor() {
    this.initializeStorage();
  }

  /**
   * Load a plugin from a URL or code
   */
  async loadPlugin(source: string | object, manifest: PluginManifest): Promise<PluginLoadResult> {
    try {
      // Validate manifest
      this.validateManifest(manifest);

      // Check if plugin already exists
      if (this.registry.has(manifest.id)) {
        throw new Error(`Plugin ${manifest.id} is already loaded`);
      }

      // Resolve dependencies
      await this.resolveDependencies(manifest);

      // Load plugin code
      const pluginCode = await this.loadPluginCode(source, manifest);

      // Create plugin sandbox
      const sandbox = await this.createSandbox(manifest, pluginCode);

      // Create plugin state
      const state: PluginState = {
        id: manifest.id,
        manifest,
        loaded: true,
        enabled: false,
        active: false,
        error: null,
        version: manifest.version,
        instance: sandbox
      };

      // Register plugin
      this.registry.set(manifest.id, state);

      // Initialize plugin
      await this.initializePlugin(state);

      return {
        success: true,
        pluginId: manifest.id
      };
    } catch (error) {
      const pluginError = this.createPluginError(manifest.id, 'LOAD_FAILED', error as Error);
      return {
        success: false,
        pluginId: manifest.id,
        error: pluginError
      };
    }
  }

  /**
   * Enable a plugin
   */
  async enablePlugin(pluginId: string): Promise<void> {
    const state = this.registry.get(pluginId);
    if (!state) {
      throw new Error(`Plugin ${pluginId} not found`);
    }

    if (state.enabled) {
      return;
    }

    // Enable dependencies first
    const deps = this.dependencyGraph.get(pluginId) || new Set();
    for (const dep of deps) {
      await this.enablePlugin(dep);
    }

    // Enable plugin
    await this.executeLifecycleHook(state, 'onEnable');
    state.enabled = true;

    // Save state
    await this.savePluginState(state);
  }

  /**
   * Disable a plugin
   */
  async disablePlugin(pluginId: string): Promise<void> {
    const state = this.registry.get(pluginId);
    if (!state) {
      throw new Error(`Plugin ${pluginId} not found`);
    }

    if (!state.enabled) {
      return;
    }

    // Disable dependent plugins first
    const dependents = this.getDependents(pluginId);
    for (const dependent of dependents) {
      await this.disablePlugin(dependent);
    }

    // Disable plugin
    await this.executeLifecycleHook(state, 'onDisable');
    state.enabled = false;
    state.active = false;

    // Save state
    await this.savePluginState(state);
  }

  /**
   * Activate a plugin (make it functional)
   */
  async activatePlugin(pluginId: string): Promise<void> {
    const state = this.registry.get(pluginId);
    if (!state) {
      throw new Error(`Plugin ${pluginId} not found`);
    }

    if (!state.enabled) {
      throw new Error(`Plugin ${pluginId} is not enabled`);
    }

    if (state.active) {
      return;
    }

    // Activate dependencies first
    const deps = this.dependencyGraph.get(pluginId) || new Set();
    for (const dep of deps) {
      await this.activatePlugin(dep);
    }

    // Activate plugin
    await this.executeLifecycleHook(state, 'onActivate');
    state.active = true;
  }

  /**
   * Deactivate a plugin
   */
  async deactivatePlugin(pluginId: string): Promise<void> {
    const state = this.registry.get(pluginId);
    if (!state) {
      throw new Error(`Plugin ${pluginId} not found`);
    }

    if (!state.active) {
      return;
    }

    // Deactivate dependent plugins
    const dependents = this.getDependents(pluginId);
    for (const dependent of dependents) {
      await this.deactivatePlugin(dependent);
    }

    // Deactivate plugin
    await this.executeLifecycleHook(state, 'onDeactivate');
    state.active = false;
  }

  /**
   * Unload a plugin
   */
  async unloadPlugin(pluginId: string): Promise<void> {
    const state = this.registry.get(pluginId);
    if (!state) {
      throw new Error(`Plugin ${pluginId} not found`);
    }

    // Deactivate if active
    if (state.active) {
      await this.deactivatePlugin(pluginId);
    }

    // Disable if enabled
    if (state.enabled) {
      await this.disablePlugin(pluginId);
    }

    // Run uninstall hook
    await this.executeLifecycleHook(state, 'onUninstall');

    // Remove from registry
    this.registry.delete(pluginId);

    // Clean up dependencies
    this.dependencyGraph.delete(pluginId);

    // Remove from storage
    await this.removePluginState(pluginId);
  }

  /**
   * Get plugin state
   */
  getPluginState(pluginId: string): PluginState | undefined {
    return this.registry.get(pluginId);
  }

  /**
   * Get all plugins
   */
  getAllPlugins(): PluginState[] {
    return Array.from(this.registry.values());
  }

  /**
   * Get enabled plugins
   */
  getEnabledPlugins(): PluginState[] {
    return this.getAllPlugins().filter(p => p.enabled);
  }

  /**
   * Get active plugins
   */
  getActivePlugins(): PluginState[] {
    return this.getAllPlugins().filter(p => p.active);
  }

  /**
   * Private methods
   */

  private validateManifest(manifest: PluginManifest): void {
    if (!manifest.id || !manifest.name || !manifest.version || !manifest.main) {
      throw new Error('Invalid manifest: missing required fields');
    }

    if (!/^[a-z0-9-]+$/.test(manifest.id)) {
      throw new Error('Invalid plugin ID: must be lowercase letters, numbers, and hyphens only');
    }

    if (!/^\d+\.\d+\.\d+$/.test(manifest.version)) {
      throw new Error('Invalid version: must be semver format (x.y.z)');
    }
  }

  private async resolveDependencies(manifest: PluginManifest): Promise<void> {
    const dependencies = manifest.dependencies || [];

    for (const depId of dependencies) {
      const dep = this.registry.get(depId);
      if (!dep) {
        throw new Error(`Dependency ${depId} not found`);
      }

      // Add to dependency graph
      if (!this.dependencyGraph.has(manifest.id)) {
        this.dependencyGraph.set(manifest.id, new Set());
      }
      this.dependencyGraph.get(manifest.id)!.add(depId);
    }
  }

  private async loadPluginCode(source: string | object, manifest: PluginManifest): Promise<any> {
    if (typeof source === 'string') {
      // Load from URL
      const response = await fetch(source);
      if (!response.ok) {
        throw new Error(`Failed to load plugin from ${source}`);
      }
      const code = await response.text();
      return this.evaluatePluginCode(code, manifest);
    } else {
      // Load from object
      return source;
    }
  }

  private evaluatePluginCode(code: string, manifest: PluginManifest): any {
    try {
      // Create isolated function
      const factory = new Function('module', 'exports', 'require', code);

      // Create module exports
      const module = { exports: {} };
      const exports = module.exports;

      // Execute in sandbox
      factory.call({}, module, exports, this.createRequire(manifest));

      return module.exports;
    } catch (error) {
      throw new Error(`Failed to evaluate plugin code: ${error}`);
    }
  }

  private createRequire(manifest: PluginManifest): (id: string) => any {
    return (id: string) => {
      const dep = this.registry.get(id);
      if (!dep) {
        throw new Error(`Cannot find module '${id}'`);
      }
      return dep.instance;
    };
  }

  private async createSandbox(manifest: PluginManifest, pluginCode: any): Promise<any> {
    // Import SDK
    const { createPluginContext } = await import('../sdk/context');

    // Create plugin context
    const context = createPluginContext(manifest);

    // Call plugin factory with context
    if (typeof pluginCode === 'function') {
      return pluginCode(context);
    } else if (pluginCode.default) {
      return pluginCode.default(context);
    } else {
      return pluginCode;
    }
  }

  private async initializePlugin(state: PluginState): Promise<void> {
    // Run install hook if this is a fresh install
    const savedState = await this.loadPluginState(state.id);
    if (!savedState) {
      await this.executeLifecycleHook(state, 'onInstall');
    }

    // Save initial state
    await this.savePluginState(state);
  }

  private async executeLifecycleHook(state: PluginState, hookName: keyof PluginHooks): Promise<void> {
    if (!state.instance) return;

    const hook = state.instance[hookName];
    if (typeof hook === 'function') {
      try {
        await hook.call(state.instance);
      } catch (error) {
        console.error(`Plugin ${state.id} ${hookName} hook failed:`, error);
        throw error;
      }
    }
  }

  private getDependents(pluginId: string): string[] {
    const dependents: string[] = [];

    for (const [id, deps] of this.dependencyGraph) {
      if (deps.has(pluginId)) {
        dependents.push(id);
      }
    }

    return dependents;
  }

  private createPluginError(pluginId: string, code: string, error: Error): PluginError {
    const pluginError = error as PluginError;
    pluginError.pluginId = pluginId;
    pluginError.code = code;
    return pluginError;
  }

  private async initializeStorage(): Promise<void> {
    const key = 'plugin-states';
    if (!localStorage.getItem(key)) {
      localStorage.setItem(key, JSON.stringify({}));
    }
  }

  private async savePluginState(state: PluginState): Promise<void> {
    const key = 'plugin-states';
    const states = JSON.parse(localStorage.getItem(key) || '{}');
    states[state.id] = {
      enabled: state.enabled,
      version: state.version
    };
    localStorage.setItem(key, JSON.stringify(states));
  }

  private async loadPluginState(pluginId: string): Promise<any> {
    const key = 'plugin-states';
    const states = JSON.parse(localStorage.getItem(key) || '{}');
    return states[pluginId];
  }

  private async removePluginState(pluginId: string): Promise<void> {
    const key = 'plugin-states';
    const states = JSON.parse(localStorage.getItem(key) || '{}');
    delete states[pluginId];
    localStorage.setItem(key, JSON.stringify(states));
  }
}

// Singleton instance
let loaderInstance: PluginLoader | null = null;

export function getPluginLoader(): PluginLoader {
  if (!loaderInstance) {
    loaderInstance = new PluginLoader();
  }
  return loaderInstance;
}
