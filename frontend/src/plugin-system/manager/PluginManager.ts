/**
 * Plugin Manager - Handles plugin installation, updates, and configuration
 */

import type {
  PluginManifest,
  PluginState,
  PluginMarketplaceInfo
} from '../types/plugin';
import { getPluginLoader } from '../core/PluginLoader';

export class PluginManager {
  private loader = getPluginLoader();
  private configCache: Map<string, any> = new Map();

  /**
   * Install a plugin
   */
  async installPlugin(
    pluginSource: string | object,
    manifest: PluginManifest
  ): Promise<{ success: boolean; error?: string }> {
    try {
      // Check if plugin is already installed
      const existing = this.loader.getPluginState(manifest.id);
      if (existing) {
        return {
          success: false,
          error: `Plugin ${manifest.id} is already installed`
        };
      }

      // Load plugin
      const result = await this.loader.loadPlugin(pluginSource, manifest);

      if (!result.success) {
        return {
          success: false,
          error: result.error?.message
        };
      }

      // Enable plugin by default
      await this.loader.enablePlugin(manifest.id);

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Uninstall a plugin
   */
  async uninstallPlugin(pluginId: string): Promise<{ success: boolean; error?: string }> {
    try {
      const state = this.loader.getPluginState(pluginId);
      if (!state) {
        return {
          success: false,
          error: `Plugin ${pluginId} is not installed`
        };
      }

      // Unload plugin
      await this.loader.unloadPlugin(pluginId);

      // Clear configuration
      this.configCache.delete(pluginId);

      // Clear storage
      const { clearPluginStorage } = await import('../sdk/storage');
      clearPluginStorage(pluginId);

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Update a plugin
   */
  async updatePlugin(
    pluginId: string,
    newSource: string | object,
    newManifest: PluginManifest
  ): Promise<{ success: boolean; error?: string }> {
    try {
      const state = this.loader.getPluginState(pluginId);
      if (!state) {
        return {
          success: false,
          error: `Plugin ${pluginId} is not installed`
        };
      }

      const oldVersion = state.version;

      // Disable plugin
      await this.loader.disablePlugin(pluginId);

      // Unload old version
      await this.loader.unloadPlugin(pluginId);

      // Load new version
      const result = await this.loader.loadPlugin(newSource, newManifest);

      if (!result.success) {
        return {
          success: false,
          error: result.error?.message
        };
      }

      // Call update hook
      const newState = this.loader.getPluginState(pluginId);
      if (newState && newState.instance) {
        const onUpdate = newState.instance.onUpdate;
        if (typeof onUpdate === 'function') {
          await onUpdate(oldVersion, newManifest.version);
        }
      }

      // Re-enable plugin
      await this.loader.enablePlugin(pluginId);

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Get plugin configuration
   */
  async getPluginConfig(pluginId: string): Promise<any> {
    if (this.configCache.has(pluginId)) {
      return this.configCache.get(pluginId);
    }

    try {
      const key = `plugin-config:${pluginId}`;
      const stored = localStorage.getItem(key);
      const config = stored ? JSON.parse(stored) : {};
      this.configCache.set(pluginId, config);
      return config;
    } catch (error) {
      console.error(`Failed to get config for plugin ${pluginId}:`, error);
      return {};
    }
  }

  /**
   * Update plugin configuration
   */
  async updatePluginConfig(
    pluginId: string,
    config: any
  ): Promise<{ success: boolean; error?: string }> {
    try {
      const key = `plugin-config:${pluginId}`;
      localStorage.setItem(key, JSON.stringify(config));
      this.configCache.set(pluginId, config);

      // Notify plugin of config change
      const state = this.loader.getPluginState(pluginId);
      if (state && state.instance && typeof state.instance.onConfigChange === 'function') {
        await state.instance.onConfigChange(config);
      }

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Reset plugin configuration
   */
  async resetPluginConfig(pluginId: string): Promise<{ success: boolean; error?: string }> {
    try {
      const key = `plugin-config:${pluginId}`;
      localStorage.removeItem(key);
      this.configCache.delete(pluginId);

      // Get default config from plugin manifest
      const state = this.loader.getPluginState(pluginId);
      if (state && state.manifest.config) {
        await this.updatePluginConfig(pluginId, state.manifest.config);
      }

      return { success: true };
    } catch (error) {
      return {
        success: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Export plugin configuration
   */
  exportPluginConfig(pluginId: string): string {
    const config = this.configCache.get(pluginId) || {};
    return JSON.stringify(config, null, 2);
  }

  /**
   * Import plugin configuration
   */
  async importPluginConfig(
    pluginId: string,
    configJson: string
  ): Promise<{ success: boolean; error?: string }> {
    try {
      const config = JSON.parse(configJson);
      return await this.updatePluginConfig(pluginId, config);
    } catch (error) {
      return {
        success: false,
        error: 'Invalid JSON format'
      };
    }
  }

  /**
   * Get plugin permissions
   */
  getPluginPermissions(pluginId: string): string[] {
    const state = this.loader.getPluginState(pluginId);
    return state?.manifest.permissions || [];
  }

  /**
   * Request additional permissions
   */
  async requestPermissions(
    pluginId: string,
    permissions: string[]
  ): Promise<{ granted: boolean; error?: string }> {
    try {
      const state = this.loader.getPluginState(pluginId);
      if (!state) {
        return {
          granted: false,
          error: `Plugin ${pluginId} not found`
        };
      }

      // Check if permissions are already granted
      const currentPermissions = state.manifest.permissions || [];
      const newPermissions = permissions.filter(p => !currentPermissions.includes(p as PluginPermission));

      if (newPermissions.length === 0) {
        return { granted: true };
      }

      // In a real app, this would show a user dialog
      const confirmed = confirm(
        `Plugin ${pluginId} is requesting additional permissions:\n\n${newPermissions.join('\n')}\n\nDo you want to grant these permissions?`
      );

      if (!confirmed) {
        return {
          granted: false,
          error: 'Permission request denied by user'
        };
      }

      // Update plugin manifest
      state.manifest.permissions = [...currentPermissions, ...newPermissions] as PluginPermission[];

      return { granted: true };
    } catch (error) {
      return {
        granted: false,
        error: (error as Error).message
      };
    }
  }

  /**
   * Get plugin statistics
   */
  getPluginStats(pluginId: string): {
    installDate: string;
    lastActive: string;
    version: string;
    enabled: boolean;
    active: boolean;
  } | null {
    const state = this.loader.getPluginState(pluginId);
    if (!state) {
      return null;
    }

    return {
      installDate: new Date().toISOString(), // This should be stored during install
      lastActive: new Date().toISOString(),
      version: state.version,
      enabled: state.enabled,
      active: state.active
    };
  }

  /**
   * Validate plugin compatibility
   */
  validateCompatibility(manifest: PluginManifest): {
    compatible: boolean;
    issues: string[];
  } {
    const issues: string[] = [];

    // Check required API version
    if (manifest.requires?.api) {
      const currentApiVersion = '1.0.0'; // This should be from app config
      if (!this.isVersionCompatible(manifest.requires.api, currentApiVersion)) {
        issues.push(`API version ${manifest.requires.api} required, but ${currentApiVersion} is installed`);
      }
    }

    // Check required app version
    if (manifest.requires?.app) {
      const currentAppVersion = '1.0.0'; // This should be from app config
      if (!this.isVersionCompatible(manifest.requires.app, currentAppVersion)) {
        issues.push(`App version ${manifest.requires.app} required, but ${currentAppVersion} is installed`);
      }
    }

    return {
      compatible: issues.length === 0,
      issues
    };
  }

  private isVersionCompatible(required: string, current: string): boolean {
    // Simple semver comparison (can be enhanced)
    const requiredParts = required.split('.').map(Number);
    const currentParts = current.split('.').map(Number);

    for (let i = 0; i < Math.min(requiredParts.length, currentParts.length); i++) {
      if (currentParts[i] < requiredParts[i]) {
        return false;
      }
      if (currentParts[i] > requiredParts[i]) {
        return true;
      }
    }

    return true;
  }
}

// Singleton instance
let managerInstance: PluginManager | null = null;

export function getPluginManager(): PluginManager {
  if (!managerInstance) {
    managerInstance = new PluginManager();
  }
  return managerInstance;
}
