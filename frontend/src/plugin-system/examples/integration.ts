/**
 * Complete Plugin System Integration Example
 * Demonstrates how to integrate the plugin system into your application
 */

import { createApp } from 'vue';
import { createPinia } from 'pinia';
import App from './App.vue';

// Import plugin system
import {
  getPluginLoader,
  getPluginManager,
  getPluginMarketplace
} from './plugin-system';

/**
 * Plugin System Integration
 */
async function integratePluginSystem() {
  console.log('🔌 Initializing Plugin System...');

  // 1. Get instances
  const loader = getPluginLoader();
  const manager = getPluginManager();
  const marketplace = getPluginMarketplace();

  // 2. Load plugins from storage
  await loadInstalledPlugins(loader);

  // 3. Setup plugin management UI
  setupPluginUI(manager, marketplace);

  // 4. Setup hot reload for development
  if (import.meta.env.DEV) {
    setupHotReload(loader);
  }

  console.log('✅ Plugin System Initialized');
}

/**
 * Load installed plugins from localStorage
 */
async function loadInstalledPlugins(loader: any) {
  try {
    const installedPlugins = JSON.parse(
      localStorage.getItem('installed-plugins') || '[]'
    );

    for (const pluginInfo of installedPlugins) {
      try {
        // Load plugin code
        const response = await fetch(pluginInfo.url);
        const code = await response.text();

        // Evaluate plugin
        const module = evaluatePlugin(code);
        const manifest = module.manifest;

        // Load plugin
        await loader.loadPlugin(module.default, manifest);

        // Enable if previously enabled
        if (pluginInfo.enabled) {
          await loader.enablePlugin(manifest.id);
          if (pluginInfo.active) {
            await loader.activatePlugin(manifest.id);
          }
        }

        console.log(`✅ Loaded plugin: ${manifest.name}`);
      } catch (error) {
        console.error(`❌ Failed to load plugin:`, error);
      }
    }
  } catch (error) {
    console.error('Failed to load installed plugins:', error);
  }
}

/**
 * Setup plugin management UI
 */
function setupPluginUI(manager: any, marketplace: any) {
  // This would integrate with your Vue app
  // Example: Add plugin management routes, components, etc.

  console.log('🎨 Plugin UI configured');
}

/**
 * Setup hot reload for development
 */
function setupHotReload(loader: any) {
  // Watch for file changes in development
  if (typeof window !== 'undefined') {
    (window as any).__pluginLoader = loader;

    console.log('🔥 Hot reload enabled');
  }
}

/**
 * Evaluate plugin code safely
 */
function evaluatePlugin(code: string) {
  try {
    // Create module wrapper
    const wrappedCode = `
      ${(code)}
      return { manifest, default: typeof createPlugin !== 'undefined' ? createPlugin : null };
    `;

    // Create function
    const factory = new Function(wrappedCode);

    // Execute
    return factory();
  } catch (error) {
    console.error('Failed to evaluate plugin code:', error);
    throw error;
  }
}

/**
 * Vue Composable for plugin system
 */
export function usePluginSystem() {
  const loader = getPluginLoader();
  const manager = getPluginManager();
  const marketplace = getPluginMarketplace();

  return {
    loader,
    manager,
    marketplace,

    // Convenience methods
    async installPlugin(url: string) {
      const response = await fetch(url);
      const code = await response.text();
      const module = evaluatePlugin(code);

      const result = await manager.installPlugin(module.default, module.manifest);

      if (result.success) {
        // Save to installed plugins
        const installedPlugins = JSON.parse(
          localStorage.getItem('installed-plugins') || '[]'
        );
        installedPlugins.push({
          id: module.manifest.id,
          url,
          enabled: true,
          active: false,
          installedAt: new Date().toISOString()
        });
        localStorage.setItem('installed-plugins', JSON.stringify(installedPlugins));
      }

      return result;
    },

    async uninstallPlugin(pluginId: string) {
      const result = await manager.uninstallPlugin(pluginId);

      if (result.success) {
        // Remove from installed plugins
        const installedPlugins = JSON.parse(
          localStorage.getItem('installed-plugins') || '[]'
        );
        const filtered = installedPlugins.filter((p: any) => p.id !== pluginId);
        localStorage.setItem('installed-plugins', JSON.stringify(filtered));
      }

      return result;
    },

    getInstalledPlugins() {
      return JSON.parse(localStorage.getItem('installed-plugins') || '[]');
    },

    getAllPlugins() {
      return loader.getAllPlugins();
    },

    getActivePlugins() {
      return loader.getActivePlugins();
    }
  };
}

/**
 * Initialize app with plugin system
 */
export async function initializeApp() {
  // Create Vue app
  const app = createApp(App);
  const pinia = createPinia();

  app.use(pinia);

  // Integrate plugin system
  await integratePluginSystem();

  // Mount app
  app.mount('#app');

  console.log('🚀 Application started with plugin system');
}

// Auto-initialize if this is the main entry
if (typeof window !== 'undefined') {
  initializeApp().catch(console.error);
}

/**
 * Example: Plugin Management Component
 */
export const PluginManagementComponent = {
  name: 'PluginManagement',

  setup() {
    const { manager, marketplace, getInstalledPlugins } = usePluginSystem();

    const installedPlugins = ref(getInstalledPlugins());
    const availablePlugins = ref([]);
    const loading = ref(false);
    const searchQuery = ref('');

    // Load available plugins
    async function loadAvailablePlugins() {
      loading.value = true;
      try {
        availablePlugins.value = await marketplace.getAllPlugins();
      } catch (error) {
        console.error('Failed to load plugins:', error);
      } finally {
        loading.value = false;
      }
    }

    // Install plugin
    async function installPlugin(plugin: any) {
      const result = await manager.installPlugin(plugin.url, plugin.manifest);
      if (result.success) {
        installedPlugins.value = getInstalledPlugins();
      }
      return result;
    }

    // Uninstall plugin
    async function uninstallPlugin(pluginId: string) {
      const result = await manager.uninstallPlugin(pluginId);
      if (result.success) {
        installedPlugins.value = getInstalledPlugins();
      }
      return result;
    }

    // Search plugins
    const filteredPlugins = computed(() => {
      if (!searchQuery.value) return availablePlugins.value;

      const query = searchQuery.value.toLowerCase();
      return availablePlugins.value.filter((p: any) =>
        p.name.toLowerCase().includes(query) ||
        p.description.toLowerCase().includes(query)
      );
    });

    // Load plugins on mount
    onMounted(() => {
      loadAvailablePlugins();
    });

    return {
      installedPlugins,
      availablePlugins,
      filteredPlugins,
      loading,
      searchQuery,
      installPlugin,
      uninstallPlugin
    };
  },

  template: `
    <div class="plugin-management">
      <h2>Plugin Management</h2>

      <!-- Installed Plugins -->
      <section>
        <h3>Installed Plugins</h3>
        <div class="plugin-list">
          <div v-for="plugin in installedPlugins" :key="plugin.id" class="plugin-item">
            <h4>{{ plugin.name }}</h4>
            <p>{{ plugin.description }}</p>
            <button @click="uninstallPlugin(plugin.id)">Uninstall</button>
          </div>
        </div>
      </section>

      <!-- Available Plugins -->
      <section>
        <h3>Available Plugins</h3>
        <input v-model="searchQuery" placeholder="Search plugins..." />
        <div v-if="loading">Loading...</div>
        <div class="plugin-list">
          <div v-for="plugin in filteredPlugins" :key="plugin.id" class="plugin-item">
            <h4>{{ plugin.name }}</h4>
            <p>{{ plugin.description }}</p>
            <div class="plugin-meta">
              <span>⭐ {{ plugin.rating }}</span>
              <span>📥 {{ plugin.downloads }}</span>
            </div>
            <button @click="installPlugin(plugin)">Install</button>
          </div>
        </div>
      </section>
    </div>
  `
};

/**
 * Example: Plugin Settings Component
 */
export const PluginSettingsComponent = {
  name: 'PluginSettings',

  props: {
    pluginId: {
      type: String,
      required: true
    }
  },

  setup(props: any) {
    const { manager } = usePluginSystem();

    const config = ref(null);
    const loading = ref(false);
    const saving = ref(false);

    // Load configuration
    async function loadConfig() {
      loading.value = true;
      try {
        config.value = await manager.getPluginConfig(props.pluginId);
      } catch (error) {
        console.error('Failed to load config:', error);
      } finally {
        loading.value = false;
      }
    }

    // Save configuration
    async function saveConfig() {
      saving.value = true;
      try {
        await manager.updatePluginConfig(props.pluginId, config.value);
      } catch (error) {
        console.error('Failed to save config:', error);
      } finally {
        saving.value = false;
      }
    }

    // Reset configuration
    async function resetConfig() {
      try {
        await manager.resetPluginConfig(props.pluginId);
        await loadConfig();
      } catch (error) {
        console.error('Failed to reset config:', error);
      }
    }

    onMounted(() => {
      loadConfig();
    });

    return {
      config,
      loading,
      saving,
      saveConfig,
      resetConfig
    };
  },

  template: `
    <div class="plugin-settings">
      <h3>Plugin Settings</h3>

      <div v-if="loading">Loading settings...</div>

      <div v-else-if="config">
        <form @submit.prevent="saveConfig">
          <div v-for="(value, key) in config" :key="key" class="setting-item">
            <label :for="key">{{ key }}</label>
            <input
              :id="key"
              v-model="config[key]"
              :type="typeof value === 'number' ? 'number' : 'text'"
            />
          </div>

          <button type="submit" :disabled="saving">
            {{ saving ? 'Saving...' : 'Save Settings' }}
          </button>

          <button type="button" @click="resetConfig">Reset to Defaults</button>
        </form>
      </div>
    </div>
  `
};

export default integratePluginSystem;
